package main

import (
	"encoding/json"
	"expvar"
	"fmt"
	"log/slog"
	"net/http"

	"tailscale.com/metrics"
	"xeiaso.net/v4/internal/github"
	"xeiaso.net/v4/internal/lume"
)

var (
	sponsorsWebhookCount = metrics.LabelMap{Label: "action"}
	sponsorsErrorCount   = metrics.LabelMap{Label: "error_type"}
)

func init() {
	expvar.Publish("gauge_xesite_sponsors_webhook_count", &sponsorsWebhookCount)
	expvar.Publish("gauge_xesite_sponsors_error_count", &sponsorsErrorCount)
}

// GitHubSponsorsWebhook handles GitHub Sponsors webhook events.
type GitHubSponsorsWebhook struct {
	fs *lume.FS
}

// ServeHTTP processes incoming GitHub Sponsors webhook events.
func (gsh *GitHubSponsorsWebhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check for GitHub Sponsors event header
	eventType := r.Header.Get("X-GitHub-Event")
	if eventType != "sponsorship" {
		slog.Info("not a sponsorship event", "event", eventType)
		sponsorsErrorCount.Add("invalid_event_type", 1)
		http.Error(w, "Invalid event type", http.StatusBadRequest)
		return
	}

	// Decode the webhook payload
	var event github.SponsorsEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		slog.Error("error decoding GitHub Sponsors event", "error", err)
		sponsorsErrorCount.Add("decode_error", 1)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Increment the specific action counter
	sponsorsWebhookCount.Add(event.Action, 1)

	// Log the sponsorship event
	slog.Info("GitHub Sponsors webhook received",
		"action", event.Action,
		"sponsor", event.Sponsorship.Sponsor.Login,
		"sponsorable", event.Sponsorship.Sponsorable.Login,
		"tier", event.Sponsorship.Tier.Name,
		"monthly_price", event.Sponsorship.Tier.MonthlyPriceInDollars,
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
		sponsorsErrorCount.Add("unhandled_action", 1)
	}

	// Respond with success
	fmt.Fprintln(w, "OK")
}

// handleSponsorshipCreated processes new sponsorships
func (gsh *GitHubSponsorsWebhook) handleSponsorshipCreated(event github.SponsorsEvent) {
	slog.Info("New sponsorship created",
		"sponsor", event.Sponsorship.Sponsor.Login,
		"tier", event.Sponsorship.Tier.Name,
		"monthly_price", event.Sponsorship.Tier.MonthlyPriceInDollars,
	)
	// TODO: Add sponsorship creation logic here
	// This could include updating database records, sending notifications, etc.
}

// handleSponsorshipEdited processes sponsorship changes
func (gsh *GitHubSponsorsWebhook) handleSponsorshipEdited(event github.SponsorsEvent) {
	slog.Info("Sponsorship edited",
		"sponsor", event.Sponsorship.Sponsor.Login,
		"tier", event.Sponsorship.Tier.Name,
		"monthly_price", event.Sponsorship.Tier.MonthlyPriceInDollars,
	)
	// TODO: Add sponsorship edit logic here
	// This could include updating records, changing access levels, etc.
}

// handleSponsorshipCancelled processes cancelled sponsorships
func (gsh *GitHubSponsorsWebhook) handleSponsorshipCancelled(event github.SponsorsEvent) {
	slog.Info("Sponsorship cancelled",
		"sponsor", event.Sponsorship.Sponsor.Login,
		"tier", event.Sponsorship.Tier.Name,
	)
	// TODO: Add sponsorship cancellation logic here
	// This could include revoking access, updating records, etc.
}

// handlePendingTierChange processes pending tier changes
func (gsh *GitHubSponsorsWebhook) handlePendingTierChange(event github.SponsorsEvent) {
	slog.Info("Pending sponsorship tier change",
		"sponsor", event.Sponsorship.Sponsor.Login,
		"tier", event.Sponsorship.Tier.Name,
		"monthly_price", event.Sponsorship.Tier.MonthlyPriceInDollars,
	)
	// TODO: Add pending tier change logic here
}

// handlePendingCancellation processes pending cancellations
func (gsh *GitHubSponsorsWebhook) handlePendingCancellation(event github.SponsorsEvent) {
	slog.Info("Pending sponsorship cancellation",
		"sponsor", event.Sponsorship.Sponsor.Login,
		"tier", event.Sponsorship.Tier.Name,
	)
	// TODO: Add pending cancellation logic here
}
