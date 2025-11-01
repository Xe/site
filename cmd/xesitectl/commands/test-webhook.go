package commands

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/subcommands"
	"xeiaso.net/v4/internal/github"
)

const (
	defaultWebhookURL = "http://localhost:4823/.within/hook/github-sponsors"
	defaultSecret     = "test-secret"
)

// TestWebhookCmd implements the test webhook command using subcommands
type TestWebhookCmd struct {
	action      string
	sponsor     string
	sponsorable string
	tier        string
	price       int
	webhookURL  string
	secret      string
	verbose     bool
}

// Name returns the command name
func (*TestWebhookCmd) Name() string { return "test-webhook" }

// Synopsis returns a short command description
func (*TestWebhookCmd) Synopsis() string { return "Test GitHub Sponsors webhook handler" }

// Usage returns detailed command usage
func (*TestWebhookCmd) Usage() string {
	return `test-webhook [-action] [-sponsor] [-sponsorable] [-tier] [-price] [-url] [-secret] [-verbose]

Test the GitHub Sponsors webhook handler by sending mock sponsorship events.

Examples:
  xesitectl test-webhook -action created -sponsor testsponsor -tier "Pro Tier" -price 10
  xesitectl test-webhook -action cancelled -sponsor oldsponsor -url https://myapp.com/hook
  xesitectl test-webhook -action edited -sponsor testsponsor -price 20 -verbose

Flags:
  -action      Event action (created, edited, cancelled, pending_tier_change, pending_cancellation)
  -sponsor     Sponsor login name (default: "testsponsor")
  -sponsorable Sponsorable login name (default: "xeiaso")
  -tier        Tier name (default: "Test Tier")
  -price       Monthly price in dollars (default: 10)
  -url         Webhook URL (default: http://localhost:3000/.within/hook/github-sponsors)
  -secret      Webhook secret for HMAC signature (default: "test-secret")
  -verbose     Enable verbose logging
`
}

// SetFlags defines the command-line flags
func (cmd *TestWebhookCmd) SetFlags(f *flag.FlagSet) {
	// Get defaults from environment variables if available
	defaultSecret := os.Getenv("GITHUB_SPONSORS_SECRET")
	if defaultSecret == "" {
		defaultSecret = "test-secret"
	}

	defaultURL := os.Getenv("WEBHOOK_URL")
	if defaultURL == "" {
		defaultURL = defaultWebhookURL
	}

	defaultSponsor := os.Getenv("TEST_SPONSOR")
	if defaultSponsor == "" {
		defaultSponsor = "testsponsor"
	}

	defaultTier := os.Getenv("TEST_TIER")
	if defaultTier == "" {
		defaultTier = "Test Tier"
	}

	defaultPrice := 10
	if priceStr := os.Getenv("TEST_PRICE"); priceStr != "" {
		if price, err := strconv.Atoi(priceStr); err == nil {
			defaultPrice = price
		}
	}

	f.StringVar(&cmd.action, "action", "created", "Sponsorship event action")
	f.StringVar(&cmd.sponsor, "sponsor", defaultSponsor, "Sponsor login name")
	f.StringVar(&cmd.sponsorable, "sponsorable", "xeiaso", "Sponsorable login name")
	f.StringVar(&cmd.tier, "tier", defaultTier, "Tier name")
	f.IntVar(&cmd.price, "price", defaultPrice, "Monthly price in dollars")
	f.StringVar(&cmd.webhookURL, "url", defaultURL, "Webhook URL")
	f.StringVar(&cmd.secret, "secret", defaultSecret, "Webhook secret")
	f.BoolVar(&cmd.verbose, "verbose", false, "Enable verbose logging")
}

// Execute runs the command
func (cmd *TestWebhookCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "Unexpected arguments: %v\n", f.Args())
		return subcommands.ExitUsageError
	}

	// Configure logging
	logLevel := slog.LevelInfo
	if cmd.verbose {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	})))

	// Validate action
	validActions := map[string]bool{
		github.SponsorsEventCreated:             true,
		github.SponsorsEventEdited:              true,
		github.SponsorsEventCancelled:           true,
		github.SponsorsEventPendingTierChange:   true,
		github.SponsorsEventPendingCancellation: true,
	}

	if !validActions[cmd.action] {
		fmt.Fprintf(os.Stderr, "Invalid action: %s\nValid actions: created, edited, cancelled, pending_tier_change, pending_cancellation\n", cmd.action)
		return subcommands.ExitUsageError
	}

	slog.Info("Testing GitHub Sponsors webhook",
		"action", cmd.action,
		"sponsor", cmd.sponsor,
		"sponsorable", cmd.sponsorable,
		"tier", cmd.tier,
		"price", cmd.price,
		"url", cmd.webhookURL,
	)

	// Create test sponsorship event
	event := cmd.createTestEvent()

	// Send webhook request
	if err := cmd.sendWebhook(ctx, event); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to send webhook: %v\n", err)
		return subcommands.ExitFailure
	}

	fmt.Printf("✅ Successfully sent %s webhook event for sponsor %s\n", cmd.action, cmd.sponsor)
	return subcommands.ExitSuccess
}

// createTestEvent creates a mock GitHub Sponsors event
func (cmd *TestWebhookCmd) createTestEvent() github.SponsorsEvent {
	now := time.Now()

	event := github.SponsorsEvent{
		Action: cmd.action,
		Sponsorship: github.Sponsorship{
			ID:        12345,
			NodeID:    "SPONSORSHIP_12345",
			CreatedAt: github.Time{Time: now},
			Sponsor: github.Sponsor{
				Login:     cmd.sponsor,
				ID:        67890,
				NodeID:    "USER_67890",
				AvatarURL: fmt.Sprintf("https://avatars.githubusercontent.com/u/67890?v=4"),
				HTMLURL:   fmt.Sprintf("https://github.com/%s", cmd.sponsor),
				Type:      "User",
			},
			Sponsorable: github.Sponsor{
				Login:     cmd.sponsorable,
				ID:        54321,
				NodeID:    "USER_54321",
				AvatarURL: fmt.Sprintf("https://avatars.githubusercontent.com/u/54321?v=4"),
				HTMLURL:   fmt.Sprintf("https://github.com/%s", cmd.sponsorable),
				Type:      "User",
			},
			Tier: github.Tier{
				ID:                    98765,
				NodeID:                "TIER_98765",
				CreatedAt:             github.Time{Time: now},
				Description:           fmt.Sprintf("Test sponsorship tier for %s", cmd.tier),
				MonthlyPriceInDollars: cmd.price,
				IsOneTime:             false,
				IsCustomAmount:        false,
				Name:                  cmd.tier,
				Published:             true,
				SelectedTier:          true,
				SponsorshipCount:      1,
			},
			PrivacyLevel:    "public",
			Variant:         "recurring",
			SponsorshipType: "user",
		},
		Sender: github.User{
			Login:     cmd.sponsor,
			ID:        67890,
			AvatarURL: fmt.Sprintf("https://avatars.githubusercontent.com/u/67890?v=4"),
			HTMLURL:   fmt.Sprintf("https://github.com/%s", cmd.sponsor),
			Type:      "User",
		},
	}

	return event
}

// sendWebhook sends the event as an HTTP request to the webhook endpoint
func (cmd *TestWebhookCmd) sendWebhook(ctx context.Context, event github.SponsorsEvent) error {
	// Marshal the event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if cmd.verbose {
		slog.Debug("Webhook payload", "json", string(eventJSON))
	}

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "POST", cmd.webhookURL, bytes.NewBuffer(eventJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", "sponsorship")
	req.Header.Set("X-GitHub-Delivery", fmt.Sprintf("%d", time.Now().UnixNano()))
	req.Header.Set("User-Agent", "Xesitectl-Webhook-Test/1.0")

	// Add HMAC signature
	if cmd.secret != "" {
		signature := cmd.calculateHMACSignature(eventJSON)
		req.Header.Set("X-Hub-Signature-256", signature)
		slog.Debug("Added HMAC signature", "signature", signature)
	}

	slog.Info("Sending webhook request",
		"url", cmd.webhookURL,
		"action", event.Action,
		"sponsor", event.Sponsorship.Sponsor.Login,
		"content_length", len(eventJSON),
	)

	// Send request with timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook returned status %d %s", resp.StatusCode, resp.Status)
	}

	slog.Info("Webhook test completed successfully",
		"action", event.Action,
		"sponsor", event.Sponsorship.Sponsor.Login,
		"status", resp.StatusCode,
	)

	return nil
}

// calculateHMACSignature calculates the HMAC-SHA256 signature for the webhook payload
func (cmd *TestWebhookCmd) calculateHMACSignature(payload []byte) string {
	mac := hmac.New(sha256.New, []byte(cmd.secret))
	mac.Write(payload)
	signature := mac.Sum(nil)
	return "sha256=" + hex.EncodeToString(signature)
}
