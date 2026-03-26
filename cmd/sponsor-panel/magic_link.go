package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"net/mail"
	"time"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"xeiaso.net/v4/cmd/sponsor-panel/templates"
)

// magicLinkRequestHandler handles POST /login/email.
func (s *Server) magicLinkRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")

	// Validate email format
	if _, err := mail.ParseAddress(email); err != nil {
		slog.Debug("magicLinkRequestHandler: invalid email", "email", email)
		templ.Handler(templates.MagicLinkSent()).ServeHTTP(w, r)
		return
	}

	// Rate limit: max 3 per email per hour
	count, err := countRecentMagicLinks(r.Context(), s.pool, email)
	if err != nil {
		slog.Error("magicLinkRequestHandler: failed to count recent links", "err", err)
		templ.Handler(templates.MagicLinkSent()).ServeHTTP(w, r)
		return
	}
	if count >= 3 {
		slog.Warn("magicLinkRequestHandler: rate limited", "email", email, "count", count)
		templ.Handler(templates.MagicLinkSent()).ServeHTTP(w, r)
		return
	}

	// Generate token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		slog.Error("magicLinkRequestHandler: failed to generate token", "err", err)
		templ.Handler(templates.MagicLinkSent()).ServeHTTP(w, r)
		return
	}
	token := base64.RawURLEncoding.EncodeToString(tokenBytes)

	// Store SHA-256 hash
	hash := sha256.Sum256([]byte(token))
	tokenHash := base64.RawURLEncoding.EncodeToString(hash[:])
	expiresAt := time.Now().Add(15 * time.Minute)

	if err := createMagicLinkToken(r.Context(), s.pool, email, tokenHash, expiresAt); err != nil {
		slog.Error("magicLinkRequestHandler: failed to store token", "err", err)
		templ.Handler(templates.MagicLinkSent()).ServeHTTP(w, r)
		return
	}

	// Send email
	link := fmt.Sprintf("%s/login/email/verify?token=%s", s.baseURL, token)
	htmlBody := fmt.Sprintf(`<html><body>
<h2>Sign in to Sponsor Panel</h2>
<p>Click the link below to sign in. This link expires in 15 minutes.</p>
<p><a href="%s">Sign in to Sponsor Panel</a></p>
<p>If you didn't request this, you can safely ignore this email.</p>
</body></html>`, link)

	if err := s.emailSender.SendEmail(r.Context(), email, "Sign in to Sponsor Panel", htmlBody); err != nil {
		slog.Error("magicLinkRequestHandler: failed to send email", "err", err, "to", email)
	} else {
		slog.Info("magicLinkRequestHandler: magic link sent", "email", email)
	}

	// Always return success to prevent email enumeration
	templ.Handler(templates.MagicLinkSent()).ServeHTTP(w, r)
}

// magicLinkVerifyHandler handles GET /login/email/verify.
func (s *Server) magicLinkVerifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		renderOAuthError(w, "Missing token")
		return
	}

	// Hash the token and look it up
	hash := sha256.Sum256([]byte(token))
	tokenHash := base64.RawURLEncoding.EncodeToString(hash[:])

	magicToken, err := consumeMagicLinkToken(r.Context(), s.pool, tokenHash)
	if err != nil {
		slog.Error("magicLinkVerifyHandler: invalid or expired token", "err", err)
		renderOAuthError(w, "Invalid or expired login link")
		return
	}

	// Upsert user
	user := &User{
		Provider:        "email",
		Login:           magicToken.Email,
		Email:           magicToken.Email,
		SponsorshipData: `{"is_active": false}`,
	}

	if err := upsertEmailUser(r.Context(), s.pool, user); err != nil {
		slog.Error("magicLinkVerifyHandler: failed to upsert user", "err", err)
		renderOAuthError(w, "Failed to create user")
		return
	}

	// Create session
	session, err := s.sessionStore.Get(r, "session")
	if err != nil {
		slog.Debug("magicLinkVerifyHandler: failed to decode existing session, creating new one", "err", err)
		session = sessions.NewSession(s.sessionStore, "session")
	}
	session.Values["user_id"] = user.ID
	if err := s.sessionStore.Save(r, w, session); err != nil {
		slog.Error("magicLinkVerifyHandler: failed to save session", "err", err)
		renderOAuthError(w, "Failed to save session")
		return
	}

	slog.Info("magicLinkVerifyHandler: user logged in", "user_id", user.ID, "email", magicToken.Email)

	http.Redirect(w, r, "/", http.StatusFound)
}
