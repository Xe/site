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

	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
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
		Secure:   s.cookieSecure,
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

// sponsorshipInfo represents the sponsorship tier info returned from GraphQL.
type sponsorshipInfo struct {
	Tier struct {
		MonthlyPriceInCents int    `json:"monthlyPriceInCents"`
		Name                string `json:"name"`
	} `json:"tier"`
	PrivacyLevel string `json:"privacyLevel"`
	IsActive     bool   `json:"isActive"`
	CreatedAt    string `json:"createdAt"`
}

// organizationSponsorship represents an organization with its sponsorship status.
type organizationSponsorship struct {
	Login                string           `json:"login"`
	Name                 string           `json:"name"`
	SponsorshipForViewer *sponsorshipInfo `json:"sponsorshipForViewer"`
}

// graphqlUserOrganizationsResponse represents the GraphQL response for user organizations with sponsorship info.
type graphqlUserOrganizationsResponse struct {
	Data struct {
		User struct {
			Login         string `json:"login"`
			Organizations struct {
				Nodes []organizationSponsorship `json:"nodes"`
			} `json:"organizations"`
		} `json:"user"`
	} `json:"data"`
}

// fetchUserOrganizations fetches the list of organizations the user belongs to via REST API.
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

	slog.Debug("fetchUserOrganizations: found organizations", "count", len(orgMap), "orgs", orgMap)
	return orgMap, nil
}

// fetchUserOrganizationsWithSponsorship fetches the user's organizations with their sponsorship status via GraphQL.
// This implements the query from the implementation guide.
func fetchUserOrganizationsWithSponsorship(ctx context.Context, token string, userLogin string) (map[string]*sponsorshipInfo, error) {
	slog.Debug("fetchUserOrganizationsWithSponsorship: fetching user organizations with sponsorship info", "user", userLogin)

	// Build the request body as a map and marshal to JSON
	reqBody := map[string]any{
		"query": `query ($userLogin: String!) { user(login: $userLogin) { organizations(first: 20) { nodes { login name sponsorshipForViewer { isActive tier { name monthlyPriceInCents } } } } } }`,
		"variables": map[string]string{"userLogin": userLogin},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.github.com/graphql", strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GraphQL API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result graphqlUserOrganizationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	orgMap := make(map[string]*sponsorshipInfo)
	for _, org := range result.Data.User.Organizations.Nodes {
		sponsorship := org.SponsorshipForViewer
		if sponsorship != nil && sponsorship.IsActive {
			slog.Debug("fetchUserOrganizationsWithSponsorship: found sponsoring org",
				"org", org.Login,
				"tier_name", sponsorship.Tier.Name,
				"monthly_amount_cents", sponsorship.Tier.MonthlyPriceInCents)
		}
		orgMap[org.Login] = sponsorship
	}

	slog.Debug("fetchUserOrganizationsWithSponsorship: found organizations", "count", len(orgMap))
	return orgMap, nil
}

// fetchSponsorshipForEntity checks if a specific entity (user or org) sponsors the viewer.
// Returns the sponsorship tier info if active, nil otherwise.
func fetchSponsorshipForEntity(ctx context.Context, token string, entityType string, entityLogin string) (*sponsorshipInfo, error) {
	// For checking if someone sponsors the viewer, we need to use viewer.sponsorshipsAsSponsor
	// and filter by the entity login
	var queryStr string
	if entityType == "user" {
		queryStr = `query {
			viewer {
				sponsorshipsAsSponsor(first: 100) {
					nodes {
						sponsorEntity {
							... on User {
								login
							}
						}
						tier {
							monthlyPriceInCents
							name
						}
						isActive
					}
				}
			}
		}`
	} else {
		queryStr = `query {
			viewer {
				sponsorshipsAsSponsor(first: 100) {
					nodes {
						sponsorEntity {
							... on Organization {
								login
							}
						}
						tier {
							monthlyPriceInCents
							name
						}
						isActive
					}
				}
			}
		}`
	}

	// Build the request body as a map and marshal to JSON
	reqBody := map[string]any{
		"query": queryStr,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.github.com/graphql", strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GraphQL API returned status %d: %s", resp.StatusCode, string(body))
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response - need a new response type for viewer.sponsorshipsAsSponsor
	type viewerSponsorshipsResponse struct {
		Data struct {
			Viewer struct {
				SponsorshipsAsSponsor struct {
					Nodes []struct {
						SponsorEntity struct {
							Login string `json:"login"`
						} `json:"sponsorEntity"`
						Tier struct {
							MonthlyPriceInCents int    `json:"monthlyPriceInCents"`
							Name                string `json:"name"`
						} `json:"tier"`
						IsActive bool `json:"isActive"`
					} `json:"nodes"`
				} `json:"sponsorshipsAsSponsor"`
			} `json:"viewer"`
		} `json:"data"`
	}

	var result viewerSponsorshipsResponse
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, err
	}

	// Find the matching sponsor
	for _, sponsorship := range result.Data.Viewer.SponsorshipsAsSponsor.Nodes {
		if sponsorship.SponsorEntity.Login == entityLogin && sponsorship.IsActive {
			slog.Debug("fetchSponsorshipForEntity: found active sponsorship",
				"entity", entityLogin,
				"tier_name", sponsorship.Tier.Name,
				"monthly_amount_cents", sponsorship.Tier.MonthlyPriceInCents)
			return &sponsorshipInfo{
				Tier: struct {
					MonthlyPriceInCents int    `json:"monthlyPriceInCents"`
					Name                string `json:"name"`
				}{
					MonthlyPriceInCents: sponsorship.Tier.MonthlyPriceInCents,
					Name:                sponsorship.Tier.Name,
				},
				IsActive: sponsorship.IsActive,
			}, nil
		}
	}

	slog.Debug("fetchSponsorshipForEntity: no active sponsorship found",
		"entity_type", entityType,
		"entity_login", entityLogin)
	return nil, nil
}

// fetchSponsorship fetches sponsorship data from GitHub GraphQL API.
// It checks the explicit allowlist first, then the synced sponsor table, then direct user sponsorship, then organizational membership.
func fetchSponsorship(ctx context.Context, pool *pgxpool.Pool, token string, userLogin string, userOrgs map[string]bool, userOrgsWithSponsorship map[string]*sponsorshipInfo, fiftyPlusSponsors map[string]bool) (string, error) {
	slog.Debug("fetchSponsorship: checking sponsorship", "user", userLogin)

	// Check if user is in the fifty-plus sponsors list first (highest priority)
	if fiftyPlusSponsors[userLogin] {
		slog.Info("fetchSponsorship: user in fifty-plus sponsors list", "user", userLogin)
		resultJSON, _ := json.Marshal(map[string]any{
			"is_active":            true,
			"monthly_amount_cents": 5000,
			"tier_name":            "Fifty Plus Sponsor",
		})
		return string(resultJSON), nil
	}

	// Check if any of the user's organizations are in the fifty-plus sponsors list (using REST API org list)
	for org := range userOrgs {
		if fiftyPlusSponsors[org] {
			slog.Info("fetchSponsorship: org in fifty-plus sponsors list", "org", org)
			resultJSON, _ := json.Marshal(map[string]any{
				"is_active":            true,
				"monthly_amount_cents": 5000,
				"tier_name":            "Fifty Plus Sponsor (via " + org + ")",
			})
			return string(resultJSON), nil
		}
	}

	// Check synced sponsor table for user or their orgs
	usernames := make([]string, 0, 1+len(userOrgs))
	usernames = append(usernames, userLogin)
	for org := range userOrgs {
		usernames = append(usernames, org)
	}

	syncedSponsors, err := getActiveSponsorsByUsernames(ctx, pool, usernames)
	if err != nil {
		slog.Warn("fetchSponsorship: failed to check synced sponsors table", "err", err)
	} else if len(syncedSponsors) > 0 {
		// Return the highest tier sponsor (already sorted by monthly_amount_cents DESC)
		sponsor := syncedSponsors[0]
		tierName := sponsor.TierName
		if sponsor.Username != userLogin {
			tierName = fmt.Sprintf("%s (via %s)", sponsor.TierName, sponsor.Username)
		}
		slog.Info("fetchSponsorship: found synced sponsor",
			"user", userLogin,
			"sponsor_username", sponsor.Username,
			"tier_name", tierName,
			"monthly_amount_cents", sponsor.MonthlyAmountCents)

		resultJSON, _ := json.Marshal(map[string]any{
			"is_active":            true,
			"monthly_amount_cents": sponsor.MonthlyAmountCents,
			"tier_name":            tierName,
		})
		return string(resultJSON), nil
	}

	// Check direct user sponsorship
	userSponsorship, err := fetchSponsorshipForEntity(ctx, token, "user", userLogin)
	if err != nil {
		slog.Warn("fetchSponsorship: failed to check user sponsorship", "err", err)
	} else if userSponsorship != nil {
		slog.Info("fetchSponsorship: found active user sponsorship",
			"user", userLogin,
			"tier_name", userSponsorship.Tier.Name,
			"monthly_amount_cents", userSponsorship.Tier.MonthlyPriceInCents)

		resultJSON, _ := json.Marshal(map[string]any{
			"is_active":            true,
			"monthly_amount_cents": userSponsorship.Tier.MonthlyPriceInCents,
			"tier_name":            userSponsorship.Tier.Name,
		})
		return string(resultJSON), nil
	}

	// Check organizational sponsorships using the GraphQL results
	for org, sponsorship := range userOrgsWithSponsorship {
		if sponsorship != nil && sponsorship.IsActive {
			slog.Info("fetchSponsorship: found active org sponsorship via GraphQL",
				"org", org,
				"tier_name", sponsorship.Tier.Name,
				"monthly_amount_cents", sponsorship.Tier.MonthlyPriceInCents)

			resultJSON, _ := json.Marshal(map[string]any{
				"is_active":            true,
				"monthly_amount_cents": sponsorship.Tier.MonthlyPriceInCents,
				"tier_name":            sponsorship.Tier.Name,
			})
			return string(resultJSON), nil
		}
	}

	slog.Debug("fetchSponsorship: no active sponsorship found", "user", userLogin)
	return `{"is_active": false}`, nil
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
		Secure:   s.cookieSecure,
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

	// Fetch user's organizations via REST API (for allowlist checking)
	userOrgs, err := fetchUserOrganizations(r.Context(), token.AccessToken)
	if err != nil {
		slog.Error("callbackHandler: failed to fetch user organizations", "err", err)
		// Non-fatal: continue with empty org map
		userOrgs = make(map[string]bool)
	}

	// Fetch user's organizations with sponsorship info via GraphQL (implementation guide query)
	userOrgsWithSponsorship, err := fetchUserOrganizationsWithSponsorship(r.Context(), token.AccessToken, ghUser.Login)
	if err != nil {
		slog.Error("callbackHandler: failed to fetch user organizations with sponsorship", "err", err)
		// Non-fatal: continue with empty org map
		userOrgsWithSponsorship = make(map[string]*sponsorshipInfo)
	}

	// Fetch sponsorship data (checks allowlist, synced table, then user, then org sponsorships)
	sponsorData, err := fetchSponsorship(r.Context(), s.pool, token.AccessToken, ghUser.Login, userOrgs, userOrgsWithSponsorship, s.fiftyPlusSponsors)
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

	if err := upsertUser(r.Context(), s.pool, user); err != nil {
		slog.Error("callbackHandler: failed to upsert user", "err", err, "github_id", ghUser.ID)
		renderOAuthError(w, "Failed to create user")
		return
	}

	slog.Debug("callbackHandler: user upserted successfully", "user_id", user.ID, "github_id", ghUser.ID)

	// Create session with user ID
	// Try to get existing session first, but if we get a decode error (old cookie format), create new one
	session, err := s.sessionStore.Get(r, "session")
	if err != nil {
		// Failed to decode existing session (probably old format), create a fresh one
		slog.Debug("callbackHandler: failed to decode existing session, creating new one", "err", err)
		session = sessions.NewSession(s.sessionStore, "session")
	}
	session.Values["user_id"] = user.ID
	if err := s.sessionStore.Save(r, w, session); err != nil {
		slog.Error("callbackHandler: failed to save session", "err", err)
		renderOAuthError(w, "Failed to save session")
		return
	}

	slog.Info("callbackHandler: user logged in successfully", "user_id", user.ID, "login", ghUser.Login)

	// Redirect to dashboard
	http.Redirect(w, r, "/", http.StatusFound)
}

// logoutHandler logs the user out by clearing the session.
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

	// Clear the session (ignore decode errors from old cookie format)
	session, err := s.sessionStore.Get(r, "session")
	if err == nil {
		session.Values["user_id"] = nil
		session.Options.MaxAge = -1
		s.sessionStore.Save(r, w, session)
	} else {
		// If we can't decode the old session, just clear the cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// getSessionUser retrieves the user from the session.
func (s *Server) getSessionUser(r *http.Request) (*User, error) {
	session, err := s.sessionStore.Get(r, "session")
	if err != nil {
		// Failed to decode session - might be old format, try to read raw cookie
		slog.Debug("getSessionUser: failed to get session, trying old format", "err", err)
		cookie, err := r.Cookie("session")
		if err != nil {
			return nil, fmt.Errorf("no session cookie")
		}
		var userID int
		if _, err := fmt.Sscanf(cookie.Value, "%d", &userID); err != nil {
			return nil, fmt.Errorf("invalid session format")
		}
		if userID == 0 {
			return nil, fmt.Errorf("invalid user id in session")
		}
		slog.Debug("getSessionUser: fetched user from old session format", "user_id", userID)
		return getUserByID(r.Context(), s.pool, userID)
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok || userID == 0 {
		slog.Debug("getSessionUser: no user_id in session")
		return nil, fmt.Errorf("no user_id in session")
	}

	slog.Debug("getSessionUser: fetching user from session", "user_id", userID)
	return getUserByID(r.Context(), s.pool, userID)
}

// renderOAuthError renders an OAuth error page.
func renderOAuthError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusBadRequest)
	templates.Base("OAuth Error", templates.OAuthError(message)).
		Render(context.Background(), w)
}
