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

// Patreon API v2 response types (JSON:API format)

type patreonIdentityResponse struct {
	Data struct {
		ID         string `json:"id"`
		Attributes struct {
			FullName string `json:"full_name"`
			Vanity   string `json:"vanity"`
			Email    string `json:"email"`
			ImageURL string `json:"image_url"`
		} `json:"attributes"`
		Relationships struct {
			Memberships struct {
				Data []struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"memberships"`
		} `json:"relationships"`
	} `json:"data"`
	Included []json.RawMessage `json:"included"`
}

type patreonMember struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes struct {
		PatronStatus                string `json:"patron_status"`
		CurrentlyEntitledAmountCents int    `json:"currently_entitled_amount_cents"`
	} `json:"attributes"`
	Relationships struct {
		Campaign struct {
			Data struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"campaign"`
	} `json:"relationships"`
}

// patreonLoginHandler initiates the Patreon OAuth flow.
func (s *Server) patreonLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.patreonOAuth == nil {
		http.NotFound(w, r)
		return
	}

	slog.Debug("patreonLoginHandler: initiating Patreon OAuth flow")

	state, err := generateState()
	if err != nil {
		slog.Error("patreonLoginHandler: failed to generate state", "err", err)
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

	url := s.patreonOAuth.AuthCodeURL(state)
	slog.Debug("patreonLoginHandler: redirecting to Patreon OAuth", "url", url)
	http.Redirect(w, r, url, http.StatusFound)
}

// patreonCallbackHandler handles the OAuth callback from Patreon.
func (s *Server) patreonCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.patreonOAuth == nil {
		http.NotFound(w, r)
		return
	}

	slog.Debug("patreonCallbackHandler: received OAuth callback")

	// Verify state for CSRF protection
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		slog.Error("patreonCallbackHandler: missing oauth_state cookie")
		renderOAuthError(w, "Invalid OAuth state")
		return
	}

	state := r.URL.Query().Get("state")
	if state != stateCookie.Value {
		slog.Error("patreonCallbackHandler: oauth state mismatch")
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
		slog.Error("patreonCallbackHandler: missing authorization code")
		renderOAuthError(w, "Missing authorization code")
		return
	}

	token, err := s.patreonOAuth.Exchange(r.Context(), code)
	if err != nil {
		oauthTotal.WithLabelValues("patreon", "error_token_exchange").Inc()
		slog.Error("patreonCallbackHandler: failed to exchange token", "err", err)
		renderOAuthError(w, "Failed to exchange token")
		return
	}

	slog.Debug("patreonCallbackHandler: token exchange successful")

	// Fetch identity from Patreon API v2
	identity, err := fetchPatreonIdentity(r.Context(), token.AccessToken)
	if err != nil {
		slog.Error("patreonCallbackHandler: failed to fetch identity", "err", err)
		renderOAuthError(w, "Failed to fetch Patreon identity")
		return
	}

	slog.Debug("patreonCallbackHandler: fetched Patreon identity",
		"patreon_id", identity.Data.ID,
		"full_name", identity.Data.Attributes.FullName,
		"vanity", identity.Data.Attributes.Vanity)

	// Determine login name (prefer vanity URL, fall back to full name)
	login := identity.Data.Attributes.Vanity
	if login == "" {
		login = identity.Data.Attributes.FullName
	}

	// Check if user is in the Patreon fifty-plus sponsors list
	slog.Debug("patreonCallbackHandler: checking fifty-plus allowlist",
		"login", login,
		"vanity", identity.Data.Attributes.Vanity,
		"full_name", identity.Data.Attributes.FullName,
		"allowlist", s.patreonFiftyPlusSpons)
	sponsorData := `{"is_active": false}`
	if s.patreonFiftyPlusSpons[login] {
		slog.Info("patreonCallbackHandler: user in patreon fifty-plus list", "login", login)
		resultJSON, _ := json.Marshal(map[string]any{
			"is_active":            true,
			"monthly_amount_cents": 5000,
			"tier_name":            "Fifty Plus Sponsor",
		})
		sponsorData = string(resultJSON)
	}

	// Find membership matching our campaign (skip if already blessed via allowlist)
	slog.Debug("patreonCallbackHandler: included resources count", "count", len(identity.Included))
	if sponsorData == `{"is_active": false}` {
		for i, raw := range identity.Included {
			slog.Debug("patreonCallbackHandler: included resource", "index", i, "raw", string(raw))
			var member patreonMember
			if err := json.Unmarshal(raw, &member); err != nil {
				continue
			}
			if member.Type != "member" {
				continue
			}
			if s.patreonCampaignID != "" && member.Relationships.Campaign.Data.ID != s.patreonCampaignID {
				continue
			}

			isActive := member.Attributes.PatronStatus == "active_patron"
			amountCents := member.Attributes.CurrentlyEntitledAmountCents

			slog.Info("patreonCallbackHandler: found matching membership",
				"campaign_id", member.Relationships.Campaign.Data.ID,
				"patron_status", member.Attributes.PatronStatus,
				"amount_cents", amountCents)

			tierName := "Patreon Supporter"
			if amountCents >= 5000 {
				tierName = "Patreon Premium"
			}

			resultJSON, _ := json.Marshal(map[string]any{
				"is_active":            isActive,
				"monthly_amount_cents": amountCents,
				"tier_name":            tierName,
			})
			sponsorData = string(resultJSON)
			break
		}
	}

	slog.Debug("patreonCallbackHandler: sponsorship data", "data", sponsorData)

	// Upsert user in database
	patreonID := identity.Data.ID
	user := &User{
		PatreonID:       &patreonID,
		Provider:        "patreon",
		Login:           login,
		AvatarURL:       identity.Data.Attributes.ImageURL,
		Name:            identity.Data.Attributes.FullName,
		Email:           identity.Data.Attributes.Email,
		SponsorshipData: sponsorData,
	}

	if err := upsertPatreonUser(r.Context(), s.pool, user); err != nil {
		oauthTotal.WithLabelValues("patreon", "error_upsert").Inc()
		slog.Error("patreonCallbackHandler: failed to upsert user", "err", err, "patreon_id", patreonID)
		renderOAuthError(w, "Failed to create user")
		return
	}

	slog.Debug("patreonCallbackHandler: user upserted successfully", "user_id", user.ID, "patreon_id", patreonID)

	// Create session with user ID
	session, err := s.sessionStore.Get(r, "session")
	if err != nil {
		slog.Debug("patreonCallbackHandler: failed to decode existing session, creating new one", "err", err)
		session = sessions.NewSession(s.sessionStore, "session")
	}
	session.Values["user_id"] = user.ID
	if err := s.sessionStore.Save(r, w, session); err != nil {
		slog.Error("patreonCallbackHandler: failed to save session", "err", err)
		renderOAuthError(w, "Failed to save session")
		return
	}

	oauthTotal.WithLabelValues("patreon", "success").Inc()
	slog.Info("patreonCallbackHandler: user logged in successfully", "user_id", user.ID, "login", login)

	http.Redirect(w, r, "/", http.StatusFound)
}

// fetchPatreonIdentity calls the Patreon API v2 identity endpoint.
func fetchPatreonIdentity(ctx context.Context, accessToken string) (*patreonIdentityResponse, error) {
	url := "https://www.patreon.com/api/oauth2/v2/identity" +
		"?include=memberships.campaign" +
		"&fields[user]=full_name,vanity,email,image_url" +
		"&fields[member]=patron_status,currently_entitled_amount_cents"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("patreon identity request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("patreon API returned status %d: %s", resp.StatusCode, string(body))
	}

	var identity patreonIdentityResponse
	if err := json.NewDecoder(resp.Body).Decode(&identity); err != nil {
		return nil, fmt.Errorf("failed to decode patreon identity: %w", err)
	}

	return &identity, nil
}
