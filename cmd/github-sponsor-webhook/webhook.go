package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"xeiaso.net/v4/internal/github"
	"xeiaso.net/v4/internal/models"
)

var (
	sponsorsWebhookCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "github_sponsors_webhook_total",
			Help: "Total number of GitHub Sponsors webhook events processed, by action",
		},
		[]string{"action"},
	)

	sponsorsErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "github_sponsors_webhook_errors_total",
			Help: "Total number of GitHub Sponsors webhook errors, by error type",
		},
		[]string{"error_type"},
	)
)

func init() {
	prometheus.MustRegister(sponsorsWebhookCount)
	prometheus.MustRegister(sponsorsErrorCount)
}

// GitHubSponsorsWebhook handles GitHub Sponsors webhook events.
type GitHubSponsorsWebhook struct {
	DB *gorm.DB // Database connection for persisting webhook data
}

// ServeHTTP processes incoming GitHub Sponsors webhook events.
func (gsh *GitHubSponsorsWebhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for web requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-GitHub-Event, X-Hub-Signature-256")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		slog.Info("method not allowed", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check for GitHub Sponsors event header
	eventType := r.Header.Get("X-GitHub-Event")
	if eventType != "sponsorship" {
		slog.Info("not a sponsorship event", "event", eventType)
		sponsorsErrorCount.WithLabelValues("invalid_event_type").Inc()
		http.Error(w, "Invalid event type", http.StatusBadRequest)
		return
	}

	// Decode the webhook payload
	var event github.SponsorsEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		slog.Error("error decoding GitHub Sponsors event", "error", err)
		sponsorsErrorCount.WithLabelValues("decode_error").Inc()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create webhook event record for tracking
	webhookEvent := &models.WebhookEvent{
		GitHubID:     r.Header.Get("X-GitHub-Delivery"),
		Action:       event.Action,
		EventType:    eventType,
		ProcessedAt:  time.Now(),
		RemoteAddr:   r.RemoteAddr,
		UserAgent:    r.UserAgent(),
		Timestamp:    time.Now(),
	}

	// Store the raw payload for debugging
	eventBytes, _ := json.Marshal(event)
	webhookEvent.SetPayload(map[string]interface{}{
		"raw_payload": string(eventBytes),
		"action":      event.Action,
	})

	// Increment the specific action counter
	sponsorsWebhookCount.WithLabelValues(event.Action).Inc()

	// Log the sponsorship event
	slog.Info("GitHub Sponsors webhook received",
		"action", event.Action,
		"sponsor", event.Sponsorship.Sponsor.Login,
		"sponsorable", event.Sponsorship.Sponsorable.Login,
		"tier", event.Sponsorship.Tier.Name,
		"monthly_price", event.Sponsorship.Tier.MonthlyPriceInDollars,
		"sender", event.Sender.Login,
	)

	// Handle different sponsorship event types with database operations
	var err error
	switch event.Action {
	case github.SponsorsEventCreated:
		err = gsh.handleSponsorshipCreated(event, webhookEvent)
	case github.SponsorsEventEdited:
		err = gsh.handleSponsorshipEdited(event, webhookEvent)
	case github.SponsorsEventCancelled:
		err = gsh.handleSponsorshipCancelled(event, webhookEvent)
	case github.SponsorsEventPendingTierChange:
		err = gsh.handlePendingTierChange(event, webhookEvent)
	case github.SponsorsEventPendingCancellation:
		err = gsh.handlePendingCancellation(event, webhookEvent)
	default:
		err = fmt.Errorf("unhandled GitHub Sponsors event type: %s", event.Action)
		sponsorsErrorCount.WithLabelValues("unhandled_action").Inc()
	}

	// Record webhook event processing result
	if err != nil {
		slog.Error("error processing webhook", "error", err, "action", event.Action)
		webhookEvent.Success = false
		webhookEvent.ErrorMessage = err.Error()
		sponsorsErrorCount.WithLabelValues("processing_error").Inc()
	} else {
		webhookEvent.Success = true
	}

	// Save webhook event to database
	if dbErr := gsh.DB.Create(webhookEvent).Error; dbErr != nil {
		slog.Error("failed to save webhook event", "error", dbErr)
		// Don't fail the request if we can't save the webhook event
	}

	// Respond with success
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "OK")
}

// handleSponsorshipCreated processes new sponsorships
func (gsh *GitHubSponsorsWebhook) handleSponsorshipCreated(event github.SponsorsEvent, webhookEvent *models.WebhookEvent) error {
	slog.Info("New sponsorship created",
		"sponsor", event.Sponsorship.Sponsor.Login,
		"sponsor_id", event.Sponsorship.Sponsor.ID,
		"tier", event.Sponsorship.Tier.Name,
		"monthly_price", event.Sponsorship.Tier.MonthlyPriceInDollars,
		"privacy_level", event.Sponsorship.PrivacyLevel,
	)

	// Create or find the sponsor account
	sponsor := models.Account{
		GitHubID:  int64(event.Sponsorship.Sponsor.ID),
		NodeID:    event.Sponsorship.Sponsor.NodeID,
		Login:     event.Sponsorship.Sponsor.Login,
		AvatarURL: event.Sponsorship.Sponsor.AvatarURL,
		URL:       event.Sponsorship.Sponsor.HTMLURL,
		Type:      event.Sponsorship.Sponsor.Type,
	}

	if err := gsh.DB.Where("github_id = ?", sponsor.GitHubID).FirstOrCreate(&sponsor).Error; err != nil {
		return fmt.Errorf("failed to create/find sponsor: %w", err)
	}

	// Create or find the sponsoree account
	sponsoree := models.Account{
		GitHubID:  int64(event.Sponsorship.Sponsorable.ID),
		NodeID:    event.Sponsorship.Sponsorable.NodeID,
		Login:     event.Sponsorship.Sponsorable.Login,
		AvatarURL: event.Sponsorship.Sponsorable.AvatarURL,
		URL:       event.Sponsorship.Sponsorable.HTMLURL,
		Type:      event.Sponsorship.Sponsorable.Type,
	}

	if err := gsh.DB.Where("github_id = ?", sponsoree.GitHubID).FirstOrCreate(&sponsoree).Error; err != nil {
		return fmt.Errorf("failed to create/find sponsoree: %w", err)
	}

	// Create or find the tier
	tier := models.Tier{
		GitHubID:            int64(event.Sponsorship.Tier.ID),
		Name:                event.Sponsorship.Tier.Name,
		MonthlyPriceInCents: event.Sponsorship.Tier.MonthlyPriceInDollars * 100, // Convert to cents
		Description:         event.Sponsorship.Tier.Description,
		SelectedTier:        event.Sponsorship.Tier.SelectedTier,
		SponsorshipCount:    event.Sponsorship.Tier.SponsorshipCount,
	}

	if err := gsh.DB.Where("github_id = ?", tier.GitHubID).FirstOrCreate(&tier).Error; err != nil {
		return fmt.Errorf("failed to create/find tier: %w", err)
	}

	// Convert GitHub timestamps
	var createdAt *time.Time
	if event.Sponsorship.CreatedAt.Time.Year() > 1 {
		createdAt = &event.Sponsorship.CreatedAt.Time
	}

	// Create the sponsorship
	sponsorship := models.Sponsorship{
		GitHubID:          int64(event.Sponsorship.ID),
		PrivacyLevel:      event.Sponsorship.PrivacyLevel,
		Variant:           event.Sponsorship.Variant,
		SponsorshipType:   event.Sponsorship.SponsorshipType,
		TierID:            tier.ID,
		SponsorID:         sponsor.ID,
		SponsoreeID:       sponsoree.ID,
		GitHubCreatedAt:   createdAt,
	}

	// Set metadata
	sponsorship.SetMetadata(map[string]interface{}{
		"github_id": event.Sponsorship.ID,
		"node_id":   event.Sponsorship.NodeID,
	})

	if err := gsh.DB.Create(&sponsorship).Error; err != nil {
		return fmt.Errorf("failed to create sponsorship: %w", err)
	}

	// Update webhook event with sponsorship ID
	webhookEvent.SponsorshipID = &sponsorship.ID

	slog.Info("Successfully created sponsorship in database",
		"sponsorship_id", sponsorship.ID,
		"sponsor", sponsor.Login,
		"sponsoree", sponsoree.Login,
		"tier", tier.Name,
	)

	return nil
}

// handleSponsorshipEdited processes sponsorship changes
func (gsh *GitHubSponsorsWebhook) handleSponsorshipEdited(event github.SponsorsEvent, webhookEvent *models.WebhookEvent) error {
	slog.Info("Sponsorship edited",
		"sponsor", event.Sponsorship.Sponsor.Login,
		"sponsor_id", event.Sponsorship.Sponsor.ID,
		"tier", event.Sponsorship.Tier.Name,
		"monthly_price", event.Sponsorship.Tier.MonthlyPriceInDollars,
		"privacy_level", event.Sponsorship.PrivacyLevel,
	)

	// Find existing sponsorship
	var sponsorship models.Sponsorship
	if err := gsh.DB.Where("github_id = ?", event.Sponsorship.ID).First(&sponsorship).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// If sponsorship doesn't exist, treat it as a new sponsorship
			return gsh.handleSponsorshipCreated(event, webhookEvent)
		}
		return fmt.Errorf("failed to find sponsorship: %w", err)
	}

	// Update tier if changed
	if event.Sponsorship.Tier.ID != 0 {
		var tier models.Tier
		if err := gsh.DB.Where("github_id = ?", event.Sponsorship.Tier.ID).First(&tier).Error; err != nil {
			return fmt.Errorf("failed to find tier: %w", err)
		}
		sponsorship.TierID = tier.ID
	}

	// Update other fields
	sponsorship.PrivacyLevel = event.Sponsorship.PrivacyLevel
	sponsorship.Variant = event.Sponsorship.Variant
	sponsorship.SponsorshipType = event.Sponsorship.SponsorshipType

	// Update GitHub timestamps
	if event.Sponsorship.CreatedAt.Time.Year() > 1 {
		sponsorship.GitHubCreatedAt = &event.Sponsorship.CreatedAt.Time
	}
	// Note: UpdatedAt is not available in the GitHub Sponsors webhook payload
	now := time.Now()
	sponsorship.GitHubUpdatedAt = &now

	// Update metadata
	sponsorship.SetMetadata(map[string]interface{}{
		"github_id": event.Sponsorship.ID,
		"node_id":   event.Sponsorship.NodeID,
		"edited_at": time.Now(),
	})

	if err := gsh.DB.Save(&sponsorship).Error; err != nil {
		return fmt.Errorf("failed to update sponsorship: %w", err)
	}

	// Update webhook event with sponsorship ID
	webhookEvent.SponsorshipID = &sponsorship.ID

	slog.Info("Successfully updated sponsorship in database",
		"sponsorship_id", sponsorship.ID,
		"sponsor", event.Sponsorship.Sponsor.Login,
	)

	return nil
}

// handleSponsorshipCancelled processes cancelled sponsorships
func (gsh *GitHubSponsorsWebhook) handleSponsorshipCancelled(event github.SponsorsEvent, webhookEvent *models.WebhookEvent) error {
	slog.Info("Sponsorship cancelled",
		"sponsor", event.Sponsorship.Sponsor.Login,
		"sponsor_id", event.Sponsorship.Sponsor.ID,
		"tier", event.Sponsorship.Tier.Name,
	)

	// Find existing sponsorship
	var sponsorship models.Sponsorship
	if err := gsh.DB.Where("github_id = ?", event.Sponsorship.ID).First(&sponsorship).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			slog.Info("Sponsorship not found for cancellation, possibly already cancelled",
				"github_id", event.Sponsorship.ID,
				"sponsor", event.Sponsorship.Sponsor.Login,
			)
			return nil // Don't treat as error, just log and continue
		}
		return fmt.Errorf("failed to find sponsorship: %w", err)
	}

	// Mark as cancelled (soft delete)
	now := time.Now()
	sponsorship.GitHubCancelledAt = &now

	// Update metadata to reflect cancellation
	metadata := sponsorship.Metadata()
	metadata["cancelled_at"] = now
	metadata["cancellation_reason"] = "github_sponsors_webhook"
	sponsorship.SetMetadata(metadata)

	if err := gsh.DB.Save(&sponsorship).Error; err != nil {
		return fmt.Errorf("failed to cancel sponsorship: %w", err)
	}

	// Update webhook event with sponsorship ID
	webhookEvent.SponsorshipID = &sponsorship.ID

	slog.Info("Successfully cancelled sponsorship in database",
		"sponsorship_id", sponsorship.ID,
		"sponsor", event.Sponsorship.Sponsor.Login,
		"cancelled_at", now,
	)

	return nil
}

// handlePendingTierChange processes pending tier changes
func (gsh *GitHubSponsorsWebhook) handlePendingTierChange(event github.SponsorsEvent, webhookEvent *models.WebhookEvent) error {
	slog.Info("Pending sponsorship tier change",
		"sponsor", event.Sponsorship.Sponsor.Login,
		"sponsor_id", event.Sponsorship.Sponsor.ID,
		"tier", event.Sponsorship.Tier.Name,
		"monthly_price", event.Sponsorship.Tier.MonthlyPriceInDollars,
	)

	// Find existing sponsorship
	var sponsorship models.Sponsorship
	if err := gsh.DB.Where("github_id = ?", event.Sponsorship.ID).First(&sponsorship).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// If sponsorship doesn't exist, treat it as a new sponsorship
			return gsh.handleSponsorshipCreated(event, webhookEvent)
		}
		return fmt.Errorf("failed to find sponsorship: %w", err)
	}

	// Update tier if changed
	if event.Sponsorship.Tier.ID != 0 {
		var tier models.Tier
		if err := gsh.DB.Where("github_id = ?", event.Sponsorship.Tier.ID).First(&tier).Error; err != nil {
			return fmt.Errorf("failed to find new tier: %w", err)
		}
		sponsorship.TierID = tier.ID
	}

	// Update metadata to reflect pending tier change
	metadata := sponsorship.Metadata()
	metadata["pending_tier_change"] = map[string]interface{}{
		"new_tier_id":   event.Sponsorship.Tier.ID,
		"new_tier_name": event.Sponsorship.Tier.Name,
		"notified_at":   time.Now(),
	}
	sponsorship.SetMetadata(metadata)

	if err := gsh.DB.Save(&sponsorship).Error; err != nil {
		return fmt.Errorf("failed to update sponsorship for pending tier change: %w", err)
	}

	// Update webhook event with sponsorship ID
	webhookEvent.SponsorshipID = &sponsorship.ID

	slog.Info("Successfully recorded pending tier change in database",
		"sponsorship_id", sponsorship.ID,
		"sponsor", event.Sponsorship.Sponsor.Login,
		"new_tier", event.Sponsorship.Tier.Name,
	)

	return nil
}

// handlePendingCancellation processes pending cancellations
func (gsh *GitHubSponsorsWebhook) handlePendingCancellation(event github.SponsorsEvent, webhookEvent *models.WebhookEvent) error {
	slog.Info("Pending sponsorship cancellation",
		"sponsor", event.Sponsorship.Sponsor.Login,
		"sponsor_id", event.Sponsorship.Sponsor.ID,
		"tier", event.Sponsorship.Tier.Name,
	)

	// Find existing sponsorship
	var sponsorship models.Sponsorship
	if err := gsh.DB.Where("github_id = ?", event.Sponsorship.ID).First(&sponsorship).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			slog.Info("Sponsorship not found for pending cancellation",
				"github_id", event.Sponsorship.ID,
				"sponsor", event.Sponsorship.Sponsor.Login,
			)
			return nil // Don't treat as error, just log and continue
		}
		return fmt.Errorf("failed to find sponsorship: %w", err)
	}

	// Update metadata to reflect pending cancellation
	metadata := sponsorship.Metadata()
	metadata["pending_cancellation"] = map[string]interface{}{
		"notified_at": time.Now(),
		"tier":       event.Sponsorship.Tier.Name,
	}
	sponsorship.SetMetadata(metadata)

	if err := gsh.DB.Save(&sponsorship).Error; err != nil {
		return fmt.Errorf("failed to update sponsorship for pending cancellation: %w", err)
	}

	// Update webhook event with sponsorship ID
	webhookEvent.SponsorshipID = &sponsorship.ID

	slog.Info("Successfully recorded pending cancellation in database",
		"sponsorship_id", sponsorship.ID,
		"sponsor", event.Sponsorship.Sponsor.Login,
	)

	return nil
}