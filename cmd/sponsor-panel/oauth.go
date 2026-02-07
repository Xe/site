package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"xeiaso.net/v4/cmd/sponsor-panel/templates"
)

// generateState generates a random OAuth state parameter.
func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// loginHandler initiates the GitHub OAuth flow.
func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	slog.Debug("loginHandler: initiating GitHub OAuth flow")

	state, err := generateState()
	if err != nil {
		slog.Error("failed to generate state", "err", err)
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	slog.Debug("loginHandler: generated OAuth state", "state", state[:8]+"...") // Log prefix only for security

	// Set state in cookie for CSRF protection
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	})

	// Redirect to GitHub
	url := s.oauth.AuthCodeURL(state)
	slog.Debug("loginHandler: redirecting to GitHub OAuth", "url", url)
	http.Redirect(w, r, url, http.StatusFound)
}

// githubUser represents the GitHub user API response.
type githubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}

// fetchGitHubUser fetches the user from GitHub using the access token.
func fetchGitHubUser(ctx context.Context, token string) (*githubUser, error) {
	slog.Debug("fetchGitHubUser: fetching user from GitHub API")
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	tokenPrefix := token[:8] + "..."
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	slog.Debug("fetchGitHubUser: sending request to GitHub", "token_prefix", tokenPrefix)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("fetchGitHubUser: request failed", "err", err)
		return nil, err
	}
	defer resp.Body.Close()

	slog.Debug("fetchGitHubUser: received response", "status", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var user githubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		slog.Error("fetchGitHubUser: failed to decode response", "err", err)
		return nil, err
	}

	slog.Debug("fetchGitHubUser: successfully fetched user",
		"github_id", user.ID,
		"login", user.Login,
		"name", user.Name,
		"email", user.Email)

	return &user, nil
}

// graphqlSponsorshipResponse represents the GraphQL sponsorship response.
type graphqlSponsorshipResponse struct {
	Data struct {
		Viewer struct {
			SponsorshipsAsMaintainer struct {
				Nodes []struct {
					SponsorEntity struct {
						Login string `json:"login"`
					} `json:"sponsorEntity"`
					Tier struct {
						MonthlyPriceInCents int    `json:"monthlyPriceInCents"`
						Name                string `json:"name"`
					} `json:"tier"`
					IsActive     bool   `json:"isActive"`
					PrivacyLevel string `json:"privacyLevel"`
				} `json:"nodes"`
			} `json:"sponsorshipsAsMaintainer"`
		} `json:"viewer"`
	} `json:"data"`
}

// fetchUserOrganizations fetches the list of organizations the user belongs to.
func fetchUserOrganizations(ctx context.Context, token string) (map[string]bool, error) {
	slog.Debug("fetchUserOrganizations: fetching user organizations")
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user/orgs", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var orgs []struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&orgs); err != nil {
		return nil, err
	}

	orgMap := make(map[string]bool)
	for _, org := range orgs {
		orgMap[org.Login] = true
	}

	slog.Debug("fetchUserOrganizations: found organizations", "count", len(orgMap))
	return orgMap, nil
}

// fetchSponsorship fetches sponsorship data from GitHub GraphQL API.
// It checks both direct sponsorship and organizational membership.
func fetchSponsorship(ctx context.Context, token string, userLogin string, userOrgs map[string]bool, fiftyPlusSponsors map[string]bool) (string, error) {
	slog.Debug("fetchSponsorship: fetching sponsorship data from GitHub GraphQL")
	query := `{
		"query": "query { viewer { sponsorshipsAsMaintainer(first: 100, includePrivate: true) { nodes { sponsorEntity { ... on User { login } ... on Organization { login } } tier { monthlyPriceInCents name } privacyLevel isActive } } } }"
	}`

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.github.com/graphql", strings.NewReader(query))
	if err != nil {
		return "", err
	}

	tokenPrefix := token[:8] + "..."
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	slog.Debug("fetchSponsorship: sending GraphQL request", "token_prefix", tokenPrefix)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("fetchSponsorship: request failed", "err", err)
		return "", err
	}
	defer resp.Body.Close()

	slog.Debug("fetchSponsorship: received response", "status", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("fetchSponsorship: GraphQL API error",
			"status", resp.StatusCode,
			"response", string(body))
		return "", fmt.Errorf("GraphQL API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result graphqlSponsorshipResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Error("fetchSponsorship: failed to decode response", "err", err)
		return "", err
	}

	nodes := result.Data.Viewer.SponsorshipsAsMaintainer.Nodes
	slog.Debug("fetchSponsorship: received sponsorship nodes", "node_count", len(nodes))

	// Log all sponsorship nodes for debugging
	for i, node := range nodes {
		slog.Debug("fetchSponsorship: sponsorship node",
			"index", i,
			"sponsor_login", node.SponsorEntity.Login,
			"tier_name", node.Tier.Name,
			"tier_amount_cents", node.Tier.MonthlyPriceInCents,
			"is_active", node.IsActive,
			"privacy_level", node.PrivacyLevel)
	}

	// Find highest active tier
	var highestPriceInCents int
	var highestName string

	for _, node := range nodes {
		if !node.IsActive {
			continue
		}

		// Check if this sponsor is an organization the user belongs to
		isUsersOrg := userOrgs[node.SponsorEntity.Login]

		if isUsersOrg {
			slog.Debug("fetchSponsorship: found organizational sponsorship",
				"org", node.SponsorEntity.Login,
				"tier_amount_cents", node.Tier.MonthlyPriceInCents)
		}

		if node.Tier.MonthlyPriceInCents > highestPriceInCents {
			highestPriceInCents = node.Tier.MonthlyPriceInCents
			highestName = node.Tier.Name
		}
	}

	// Check if user or their orgs are in the fifty-plus sponsors list
	if fiftyPlusSponsors[userLogin] {
		slog.Info("fetchSponsorship: user in fifty-plus sponsors list", "user", userLogin)
		if highestPriceInCents < 5000 {
			highestPriceInCents = 5000
			highestName = "Fifty Plus Sponsor"
		}
	}
	for org := range userOrgs {
		if fiftyPlusSponsors[org] {
			slog.Info("fetchSponsorship: org in fifty-plus sponsors list", "org", org)
			if highestPriceInCents < 5000 {
				highestPriceInCents = 5000
				highestName = "Fifty Plus Sponsor (via " + org + ")"
			}
			break
		}
	}

	if highestPriceInCents == 0 {
		slog.Debug("fetchSponsorship: no active sponsorship found")
		return `{"is_active": false}`, nil
	}

	slog.Info("fetchSponsorship: found active sponsorship",
		"tier_name", highestName,
		"monthly_amount_cents", highestPriceInCents)

	resultJSON, _ := json.Marshal(map[string]interface{}{
		"is_active":            true,
		"monthly_amount_cents": highestPriceInCents,
		"tier_name":            highestName,
	})

	return string(resultJSON), nil
}

// callbackHandler handles the OAuth callback from GitHub.
func (s *Server) callbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	slog.Debug("callbackHandler: received OAuth callback")

	// Verify state for CSRF protection
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		slog.Error("callbackHandler: missing oauth_state cookie")
		renderOAuthError(w, "Invalid OAuth state")
		return
	}

	state := r.URL.Query().Get("state")
	if state != stateCookie.Value {
		slog.Error("callbackHandler: oauth state mismatch",
			"query_state", state[:8]+"...",
			"cookie_state", stateCookie.Value[:8]+"...")
		renderOAuthError(w, "Invalid OAuth state")
		return
	}

	slog.Debug("callbackHandler: oauth state verified successfully")

	// Clear the state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Exchange code for token
	code := r.URL.Query().Get("code")
	if code == "" {
		slog.Error("callbackHandler: missing authorization code")
		renderOAuthError(w, "Missing authorization code")
		return
	}

	slog.Debug("callbackHandler: exchanging code for token", "code_prefix", code[:8]+"...")

	token, err := s.oauth.Exchange(r.Context(), code)
	if err != nil {
		slog.Error("callbackHandler: failed to exchange token", "err", err)
		renderOAuthError(w, "Failed to exchange token")
		return
	}

	slog.Debug("callbackHandler: token exchange successful")

	// Fetch user from GitHub
	ghUser, err := fetchGitHubUser(r.Context(), token.AccessToken)
	if err != nil {
		slog.Error("callbackHandler: failed to fetch user", "err", err)
		renderOAuthError(w, "Failed to fetch user")
		return
	}

	slog.Debug("callbackHandler: fetched GitHub user", "github_id", ghUser.ID, "login", ghUser.Login)

	// Fetch user's organizations
	userOrgs, err := fetchUserOrganizations(r.Context(), token.AccessToken)
	if err != nil {
		slog.Error("callbackHandler: failed to fetch user organizations", "err", err)
		// Non-fatal: continue with empty org map
		userOrgs = make(map[string]bool)
	}

	// Fetch sponsorship data from GraphQL
	sponsorData, err := fetchSponsorship(r.Context(), token.AccessToken, ghUser.Login, userOrgs, s.fiftyPlusSponsors)
	if err != nil {
		slog.Error("callbackHandler: failed to fetch sponsorship", "err", err)
		// Non-fatal: continue with empty sponsorship data
		sponsorData = `{"is_active": false}`
	}

	slog.Debug("callbackHandler: sponsorship data", "data", sponsorData)

	// Upsert user in database
	user := &User{
		GitHubID:        ghUser.ID,
		Login:           ghUser.Login,
		AvatarURL:       ghUser.AvatarURL,
		Name:            ghUser.Name,
		Email:           ghUser.Email,
		SponsorshipData: sponsorData,
	}

	if err := upsertUser(s.db, user); err != nil {
		slog.Error("callbackHandler: failed to upsert user", "err", err, "github_id", ghUser.ID)
		renderOAuthError(w, "Failed to create user")
		return
	}

	slog.Debug("callbackHandler: user upserted successfully", "user_id", user.ID, "github_id", ghUser.ID)

	// Create session cookie with just user ID
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    fmt.Sprintf("%d", user.ID),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		MaxAge:   30 * 24 * 3600, // 30 days
	})

	slog.Info("callbackHandler: user logged in successfully", "user_id", user.ID, "login", ghUser.Login)

	// Redirect to dashboard
	http.Redirect(w, r, "/", http.StatusFound)
}

// logoutHandler logs the user out by clearing the session cookie.
func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user before logout for logging
	user, err := s.getSessionUser(r)
	if err == nil {
		slog.Info("logoutHandler: user logged out", "user_id", user.ID, "login", user.Login)
	} else {
		slog.Debug("logoutHandler: no active session to logout")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

// getSessionUser retrieves the user from the session cookie.
func (s *Server) getSessionUser(r *http.Request) (*User, error) {
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		slog.Debug("getSessionUser: no session cookie found")
		return nil, fmt.Errorf("no session cookie")
	}

	var userID int
	if _, err := fmt.Sscanf(sessionCookie.Value, "%d", &userID); err != nil {
		slog.Debug("getSessionUser: invalid user ID in session", "cookie_value", sessionCookie.Value)
		return nil, fmt.Errorf("invalid user ID in session")
	}

	slog.Debug("getSessionUser: fetching user from session", "user_id", userID)
	return getUserByID(s.db, userID)
}

// renderOAuthError renders an OAuth error page.
func renderOAuthError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusBadRequest)
	templates.Base("OAuth Error", templates.OAuthError(message)).
		Render(context.Background(), w)
}
