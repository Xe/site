package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"xeiaso.net/v4/internal/github"
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
	// No lume.FS dependency for the standalone service
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

	// Handle different sponsorship event types
	switch event.Action {
	case github.SponsorsEventCreated:
		gsh.handleSponsorshipCreated(event)
	case github.SponsorsEventEdited:
		gsh.handleSponsorshipEdited(event)
	case github.SponsorsEventCancelled:
		gsh.handleSponsorshipCancelled(event)
	case github.SponsorsEventPendingTierChange:
		gsh.handlePendingTierChange(event)
	case github.SponsorsEventPendingCancellation:
		gsh.handlePendingCancellation(event)
	default:
		slog.Info("unhandled GitHub Sponsors event type", "action", event.Action)
		sponsorsErrorCount.WithLabelValues("unhandled_action").Inc()
	}

	// Respond with success
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "OK")
}

// handleSponsorshipCreated processes new sponsorships
func (gsh *GitHubSponsorsWebhook) handleSponsorshipCreated(event github.SponsorsEvent) {
	slog.Info("New sponsorship created",
		"sponsor", event.Sponsorship.Sponsor.Login,
		"sponsor_id", event.Sponsorship.Sponsor.ID,
		"tier", event.Sponsorship.Tier.Name,
		"monthly_price", event.Sponsorship.Tier.MonthlyPriceInDollars,
		"privacy_level", event.Sponsorship.PrivacyLevel,
	)
	// TODO: Add sponsorship creation logic here
	// This could include updating database records, sending notifications, etc.
}

// handleSponsorshipEdited processes sponsorship changes
func (gsh *GitHubSponsorsWebhook) handleSponsorshipEdited(event github.SponsorsEvent) {
	slog.Info("Sponsorship edited",
		"sponsor", event.Sponsorship.Sponsor.Login,
		"sponsor_id", event.Sponsorship.Sponsor.ID,
		"tier", event.Sponsorship.Tier.Name,
		"monthly_price", event.Sponsorship.Tier.MonthlyPriceInDollars,
		"privacy_level", event.Sponsorship.PrivacyLevel,
	)
	// TODO: Add sponsorship edit logic here
	// This could include updating records, changing access levels, etc.
}

// handleSponsorshipCancelled processes cancelled sponsorships
func (gsh *GitHubSponsorsWebhook) handleSponsorshipCancelled(event github.SponsorsEvent) {
	slog.Info("Sponsorship cancelled",
		"sponsor", event.Sponsorship.Sponsor.Login,
		"sponsor_id", event.Sponsorship.Sponsor.ID,
		"tier", event.Sponsorship.Tier.Name,
	)
	// TODO: Add sponsorship cancellation logic here
	// This could include revoking access, updating records, etc.
}

// handlePendingTierChange processes pending tier changes
func (gsh *GitHubSponsorsWebhook) handlePendingTierChange(event github.SponsorsEvent) {
	slog.Info("Pending sponsorship tier change",
		"sponsor", event.Sponsorship.Sponsor.Login,
		"sponsor_id", event.Sponsorship.Sponsor.ID,
		"tier", event.Sponsorship.Tier.Name,
		"monthly_price", event.Sponsorship.Tier.MonthlyPriceInDollars,
	)
	// TODO: Add pending tier change logic here
}

// handlePendingCancellation processes pending cancellations
func (gsh *GitHubSponsorsWebhook) handlePendingCancellation(event github.SponsorsEvent) {
	slog.Info("Pending sponsorship cancellation",
		"sponsor", event.Sponsorship.Sponsor.Login,
		"sponsor_id", event.Sponsorship.Sponsor.ID,
		"tier", event.Sponsorship.Tier.Name,
	)
	// TODO: Add pending cancellation logic here
}