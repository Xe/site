package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
)

// googleUserInfo represents the Google OAuth2 userinfo response.
type googleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// googleLoginHandler initiates the Google OAuth flow.
func (s *Server) googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.googleOAuth == nil {
		http.NotFound(w, r)
		return
	}

	slog.Debug("googleLoginHandler: initiating Google OAuth flow")

	state, err := generateState()
	if err != nil {
		slog.Error("googleLoginHandler: failed to generate state", "err", err)
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	// Set state in cookie for CSRF protection
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   s.cookieSecure,
	})

	url := s.googleOAuth.AuthCodeURL(state)
	slog.Debug("googleLoginHandler: redirecting to Google OAuth", "url", url)
	http.Redirect(w, r, url, http.StatusFound)
}

// googleCallbackHandler handles the OAuth callback from Google.
func (s *Server) googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.googleOAuth == nil {
		http.NotFound(w, r)
		return
	}

	slog.Debug("googleCallbackHandler: received OAuth callback")

	// Verify state for CSRF protection
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		slog.Error("googleCallbackHandler: missing oauth_state cookie")
		renderOAuthError(w, "Invalid OAuth state")
		return
	}

	state := r.URL.Query().Get("state")
	if state != stateCookie.Value {
		slog.Error("googleCallbackHandler: oauth state mismatch")
		renderOAuthError(w, "Invalid OAuth state")
		return
	}

	// Clear the state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   s.cookieSecure,
	})

	// Exchange code for token
	code := r.URL.Query().Get("code")
	if code == "" {
		slog.Error("googleCallbackHandler: missing authorization code")
		renderOAuthError(w, "Missing authorization code")
		return
	}

	token, err := s.googleOAuth.Exchange(r.Context(), code)
	if err != nil {
		slog.Error("googleCallbackHandler: failed to exchange token", "err", err)
		renderOAuthError(w, "Failed to exchange token")
		return
	}

	slog.Debug("googleCallbackHandler: token exchange successful")

	// Fetch user info from Google
	userInfo, err := fetchGoogleUserInfo(r.Context(), token.AccessToken)
	if err != nil {
		slog.Error("googleCallbackHandler: failed to fetch user info", "err", err)
		renderOAuthError(w, "Failed to fetch Google user info")
		return
	}

	slog.Debug("googleCallbackHandler: fetched Google user info",
		"google_id", userInfo.ID,
		"name", userInfo.Name,
		"email", userInfo.Email)

	// No sponsorship check for Google login
	sponsorData := `{"is_active": false}`

	// Upsert user in database
	googleID := userInfo.ID
	login := userInfo.Name
	if login == "" {
		login = userInfo.Email
	}

	user := &User{
		GoogleID:        &googleID,
		Provider:        "google",
		Login:           login,
		AvatarURL:       userInfo.Picture,
		Name:            userInfo.Name,
		Email:           userInfo.Email,
		SponsorshipData: sponsorData,
	}

	if err := upsertGoogleUser(r.Context(), s.pool, user); err != nil {
		slog.Error("googleCallbackHandler: failed to upsert user", "err", err, "google_id", googleID)
		renderOAuthError(w, "Failed to create user")
		return
	}

	slog.Debug("googleCallbackHandler: user upserted successfully", "user_id", user.ID, "google_id", googleID)

	// Create session with user ID
	session, err := s.sessionStore.Get(r, "session")
	if err != nil {
		slog.Debug("googleCallbackHandler: failed to decode existing session, creating new one", "err", err)
		session = sessions.NewSession(s.sessionStore, "session")
	}
	session.Values["user_id"] = user.ID
	if err := s.sessionStore.Save(r, w, session); err != nil {
		slog.Error("googleCallbackHandler: failed to save session", "err", err)
		renderOAuthError(w, "Failed to save session")
		return
	}

	slog.Info("googleCallbackHandler: user logged in successfully", "user_id", user.ID, "login", login)

	http.Redirect(w, r, "/", http.StatusFound)
}

// fetchGoogleUserInfo calls the Google OAuth2 userinfo endpoint.
func fetchGoogleUserInfo(ctx context.Context, accessToken string) (*googleUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("google userinfo request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google API returned status %d: %s", resp.StatusCode, string(body))
	}

	var userInfo googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode google userinfo: %w", err)
	}

	return &userInfo, nil
}
