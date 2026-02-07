package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"xeiaso.net/v4/cmd/sponsor-panel/templates"
)

// loginPageHandler renders the login page.
func (s *Server) loginPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	slog.Debug("loginPageHandler: rendering login page")

	// If already authenticated, redirect to dashboard
	if user, err := s.getSessionUser(r); err == nil {
		slog.Debug("loginPageHandler: user already authenticated, redirecting to dashboard", "user_id", user.ID, "login", user.Login)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	templ.Handler(
		templates.Base("Login", templates.Login(s.discordInvite)),
	).ServeHTTP(w, r)
}

// dashboardHandler renders the main dashboard for authenticated sponsors.
func (s *Server) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check authentication
	user, err := s.getSessionUser(r)
	if err != nil {
		slog.Debug("dashboardHandler: unauthenticated user, showing login page", "err", err)
		// Show login page instead of redirecting
		s.loginPageHandler(w, r)
		return
	}

	slog.Debug("dashboardHandler: authenticated user", "user_id", user.ID, "login", user.Login)

	// Check sponsorship tiers
	isFiftyPlus := user.IsSponsorAtTier(5000) // $50 = 5000 cents
	isSponsor := user.IsSponsorAtTier(100)    // $1 = 100 cents

	slog.Debug("dashboardHandler: sponsorship tier check",
		"user_id", user.ID,
		"is_sponsor", isSponsor,
		"is_fifty_plus", isFiftyPlus)

	// Parse sponsorship data for display
	monthlyAmount := 0
	tierName := "Sponsor"
	if user.SponsorshipData != "" {
		var data SponsorshipData
		if err := json.Unmarshal([]byte(user.SponsorshipData), &data); err == nil {
			if data.IsActive {
				monthlyAmount = data.MonthlyAmount
				tierName = data.TierName
				if tierName == "" {
					tierName = "Sponsor"
				}
			}
			slog.Debug("dashboardHandler: parsed sponsorship data",
				"user_id", user.ID,
				"is_active", data.IsActive,
				"monthly_amount", monthlyAmount,
				"tier_name", tierName)
		} else {
			slog.Error("dashboardHandler: failed to parse sponsorship data", "err", err, "user_id", user.ID, "raw_data", user.SponsorshipData)
		}
	} else {
		slog.Debug("dashboardHandler: no sponsorship data", "user_id", user.ID)
	}

	props := templates.DashboardProps{
		User: templates.UserProps{
			Login:     user.Login,
			AvatarURL: user.AvatarURL,
		},
		IsSponsor:     isSponsor,
		SponsorAmount: monthlyAmount,
		SponsorTier:   tierName,
		IsFiftyPlus:   isFiftyPlus,
		DiscordInvite: s.discordInvite,
	}

	slog.Debug("dashboardHandler: rendering dashboard", "user_id", user.ID, "login", user.Login)

	w.Header().Set("Content-Type", "text/html")

	templ.Handler(
		templates.Base("Dashboard", templates.Dashboard(props)),
	).ServeHTTP(w, r)
}
