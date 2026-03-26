package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/webhook"
)

// processedEvents tracks webhook event IDs to detect duplicates.
// Events expire after 24 hours.
var processedEvents = struct {
	sync.Mutex
	m map[string]time.Time
}{m: make(map[string]time.Time)}

func isEventProcessed(eventID string) bool {
	processedEvents.Lock()
	defer processedEvents.Unlock()

	// Clean up old entries
	cutoff := time.Now().Add(-24 * time.Hour)
	for id, t := range processedEvents.m {
		if t.Before(cutoff) {
			delete(processedEvents.m, id)
		}
	}

	if _, ok := processedEvents.m[eventID]; ok {
		return true
	}
	processedEvents.m[eventID] = time.Now()
	return false
}

// stripeWebhookHandler handles incoming Stripe webhook events.
func (s *Server) stripeWebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	const maxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("stripeWebhookHandler: failed to read body", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), s.stripeWebhookSecret)
	if err != nil {
		slog.Error("stripeWebhookHandler: signature verification failed", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if isEventProcessed(event.ID) {
		slog.Warn("stripeWebhookHandler: duplicate event", "event_id", event.ID, "type", event.Type)
		w.WriteHeader(http.StatusOK)
		return
	}

	slog.Info("stripeWebhookHandler: processing event", "event_id", event.ID, "type", event.Type)

	switch event.Type {
	case "invoice.paid":
		s.handleInvoicePaid(event)
	case "invoice.payment_failed":
		s.handleInvoicePaymentFailed(event)
	case "customer.updated":
		slog.Info("stripeWebhookHandler: customer updated", "event_id", event.ID)
	default:
		slog.Debug("stripeWebhookHandler: unhandled event type", "type", event.Type)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"received": true}`))
}

func (s *Server) handleInvoicePaid(event stripe.Event) {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		slog.Error("handleInvoicePaid: failed to unmarshal invoice", "err", err, "event_id", event.ID)
		return
	}

	if invoice.Customer == nil {
		slog.Error("handleInvoicePaid: invoice has no customer", "invoice_id", invoice.ID)
		return
	}

	ctx := context.Background()
	user, err := getUserByStripeCustomerID(ctx, s.pool, invoice.Customer.ID)
	if err != nil {
		slog.Error("handleInvoicePaid: user not found for stripe customer", "stripe_customer_id", invoice.Customer.ID, "err", err)
		return
	}

	data := SponsorshipData{
		IsActive:      true,
		MonthlyAmount: int(invoice.AmountPaid),
		TierName:      "Stripe Sponsor",
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		slog.Error("handleInvoicePaid: failed to marshal sponsorship data", "err", err)
		return
	}

	_, err = s.pool.Exec(ctx, `
		UPDATE users SET sponsorship_data=$1, last_sponsorship_check=NOW(), updated_at=NOW() WHERE id=$2
	`, string(dataJSON), user.ID)
	if err != nil {
		slog.Error("handleInvoicePaid: failed to update user", "err", err, "user_id", user.ID)
		return
	}

	slog.Info("handleInvoicePaid: updated sponsorship",
		"user_id", user.ID,
		"login", user.Login,
		"amount_paid", invoice.AmountPaid,
		"invoice_id", invoice.ID)
}

func (s *Server) handleInvoicePaymentFailed(event stripe.Event) {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		slog.Error("handleInvoicePaymentFailed: failed to unmarshal invoice", "err", err, "event_id", event.ID)
		return
	}

	if invoice.Customer == nil {
		slog.Error("handleInvoicePaymentFailed: invoice has no customer", "invoice_id", invoice.ID)
		return
	}

	ctx := context.Background()
	user, err := getUserByStripeCustomerID(ctx, s.pool, invoice.Customer.ID)
	if err != nil {
		slog.Error("handleInvoicePaymentFailed: user not found for stripe customer", "stripe_customer_id", invoice.Customer.ID, "err", err)
		return
	}

	data := SponsorshipData{
		IsActive: false,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		slog.Error("handleInvoicePaymentFailed: failed to marshal sponsorship data", "err", err)
		return
	}

	_, err = s.pool.Exec(ctx, `
		UPDATE users SET sponsorship_data=$1, last_sponsorship_check=NOW(), updated_at=NOW() WHERE id=$2
	`, string(dataJSON), user.ID)
	if err != nil {
		slog.Error("handleInvoicePaymentFailed: failed to update user", "err", err, "user_id", user.ID)
		return
	}

	slog.Info("handleInvoicePaymentFailed: marked sponsorship inactive",
		"user_id", user.ID,
		"login", user.Login,
		"invoice_id", invoice.ID)
}

// billingPortalHandler redirects authenticated users with a stripe_customer_id
// to the Stripe Billing Portal for self-service management.
func (s *Server) billingPortalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := s.getSessionUser(r)
	if err != nil {
		http.Redirect(w, r, "/login-page", http.StatusFound)
		return
	}

	if user.StripeCustomerID == nil || *user.StripeCustomerID == "" {
		http.Error(w, "No Stripe billing account linked", http.StatusBadRequest)
		return
	}

	params := &stripe.BillingPortalSessionCreateParams{
		Customer:  stripe.String(*user.StripeCustomerID),
		ReturnURL: stripe.String(fmt.Sprintf("%s/", s.baseURL)),
	}
	if s.stripePortalConfigID != "" {
		params.Configuration = stripe.String(s.stripePortalConfigID)
	}

	session, err := s.stripeClient.V1BillingPortalSessions.Create(r.Context(), params)
	if err != nil {
		slog.Error("billingPortalHandler: failed to create portal session", "err", err, "user_id", user.ID)
		http.Error(w, "Failed to create billing portal session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, session.URL, http.StatusSeeOther)
}
