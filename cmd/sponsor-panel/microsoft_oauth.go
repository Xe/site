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

// microsoftUserInfo represents the Microsoft Graph /me response.
type microsoftUserInfo struct {
	ID                string `json:"id"`
	DisplayName       string `json:"displayName"`
	Mail              string `json:"mail"`
	UserPrincipalName string `json:"userPrincipalName"`
}

// microsoftLoginHandler initiates the Microsoft OAuth flow.
func (s *Server) microsoftLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.microsoftOAuth == nil {
		http.NotFound(w, r)
		return
	}

	slog.Debug("microsoftLoginHandler: initiating Microsoft OAuth flow")

	state, err := generateState()
	if err != nil {
		slog.Error("microsoftLoginHandler: failed to generate state", "err", err)
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

	url := s.microsoftOAuth.AuthCodeURL(state)
	slog.Debug("microsoftLoginHandler: redirecting to Microsoft OAuth", "url", url)
	http.Redirect(w, r, url, http.StatusFound)
}

// microsoftCallbackHandler handles the OAuth callback from Microsoft.
func (s *Server) microsoftCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.microsoftOAuth == nil {
		http.NotFound(w, r)
		return
	}

	slog.Debug("microsoftCallbackHandler: received OAuth callback")

	// Verify state for CSRF protection
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		slog.Error("microsoftCallbackHandler: missing oauth_state cookie")
		renderOAuthError(w, "Invalid OAuth state")
		return
	}

	state := r.URL.Query().Get("state")
	if state != stateCookie.Value {
		slog.Error("microsoftCallbackHandler: oauth state mismatch")
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
		slog.Error("microsoftCallbackHandler: missing authorization code")
		renderOAuthError(w, "Missing authorization code")
		return
	}

	token, err := s.microsoftOAuth.Exchange(r.Context(), code)
	if err != nil {
		slog.Error("microsoftCallbackHandler: failed to exchange token", "err", err)
		renderOAuthError(w, "Failed to exchange token")
		return
	}

	slog.Debug("microsoftCallbackHandler: token exchange successful")

	// Fetch user info from Microsoft Graph
	userInfo, err := fetchMicrosoftUser(r.Context(), token.AccessToken)
	if err != nil {
		slog.Error("microsoftCallbackHandler: failed to fetch user info", "err", err)
		renderOAuthError(w, "Failed to fetch Microsoft user info")
		return
	}

	slog.Debug("microsoftCallbackHandler: fetched Microsoft user",
		"microsoft_id", userInfo.ID,
		"display_name", userInfo.DisplayName)

	// Determine login (prefer mail, fall back to userPrincipalName)
	login := userInfo.Mail
	if login == "" {
		login = userInfo.UserPrincipalName
	}

	email := userInfo.Mail
	if email == "" {
		email = userInfo.UserPrincipalName
	}

	// No sponsorship check for Microsoft users
	sponsorData := `{"is_active": false}`

	// Upsert user in database
	microsoftID := userInfo.ID
	user := &User{
		MicrosoftID:     &microsoftID,
		Provider:        "microsoft",
		Login:           login,
		Name:            userInfo.DisplayName,
		Email:           email,
		SponsorshipData: sponsorData,
	}

	if err := upsertMicrosoftUser(r.Context(), s.pool, user); err != nil {
		slog.Error("microsoftCallbackHandler: failed to upsert user", "err", err, "microsoft_id", microsoftID)
		renderOAuthError(w, "Failed to create user")
		return
	}

	slog.Debug("microsoftCallbackHandler: user upserted successfully", "user_id", user.ID, "microsoft_id", microsoftID)

	// Create session with user ID
	session, err := s.sessionStore.Get(r, "session")
	if err != nil {
		slog.Debug("microsoftCallbackHandler: failed to decode existing session, creating new one", "err", err)
		session = sessions.NewSession(s.sessionStore, "session")
	}
	session.Values["user_id"] = user.ID
	if err := s.sessionStore.Save(r, w, session); err != nil {
		slog.Error("microsoftCallbackHandler: failed to save session", "err", err)
		renderOAuthError(w, "Failed to save session")
		return
	}

	slog.Info("microsoftCallbackHandler: user logged in successfully", "user_id", user.ID, "login", login)

	http.Redirect(w, r, "/", http.StatusFound)
}

// fetchMicrosoftUser calls the Microsoft Graph /me endpoint.
func fetchMicrosoftUser(ctx context.Context, accessToken string) (*microsoftUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://graph.microsoft.com/v1.0/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("microsoft graph request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("microsoft graph API returned status %d: %s", resp.StatusCode, string(body))
	}

	var userInfo microsoftUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode microsoft user info: %w", err)
	}

	return &userInfo, nil
}
