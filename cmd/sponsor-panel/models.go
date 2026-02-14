package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// User represents a GitHub user with cached sponsorship data.
type User struct {
	ID                   int       `json:"id" db:"id"`
	GitHubID             int64     `json:"github_id" db:"github_id"`
	Login                string    `json:"login" db:"login"`
	AvatarURL            string    `json:"avatar_url" db:"avatar_url"`
	Name                 string    `json:"name" db:"name"`
	Email                string    `json:"email" db:"email"`
	SponsorshipData      string    `json:"-" db:"sponsorship_data"` // JSON blob
	LastSponsorshipCheck time.Time `json:"last_sponsorship_check" db:"last_sponsorship_check"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// SponsorshipData represents the cached GraphQL response.
type SponsorshipData struct {
	IsActive      bool   `json:"is_active"`
	MonthlyAmount int    `json:"monthly_amount_cents"`
	TierName      string `json:"tier_name"`
	PrivacyLevel  string `json:"privacy_level"`
}

// IsSponsorAtTier returns true if user sponsors at or above the given amount (in cents).
func (u *User) IsSponsorAtTier(minCents int) bool {
	if u.SponsorshipData == "" {
		return false
	}

	var data SponsorshipData
	if err := json.Unmarshal([]byte(u.SponsorshipData), &data); err != nil {
		slog.Error("IsSponsorAtTier: failed to parse sponsorship data", "user_id", u.ID, "err", err, "raw_data", u.SponsorshipData)
		return false
	}

	result := data.IsActive && data.MonthlyAmount >= minCents
	slog.Debug("IsSponsorAtTier: tier check",
		"user_id", u.ID,
		"login", u.Login,
		"min_cents", minCents,
		"actual_cents", data.MonthlyAmount,
		"is_active", data.IsActive,
		"result", result)

	return result
}

// LogoSubmission represents a logo submission.
type LogoSubmission struct {
	ID                int       `json:"id" db:"id"`
	UserID            int       `json:"user_id" db:"user_id"`
	CompanyName       string    `json:"company_name" db:"company_name"`
	Website           string    `json:"website" db:"website"`
	LogoURL           string    `json:"logo_url" db:"logo_url"`
	GitHubIssueURL    string    `json:"github_issue_url" db:"github_issue_url"`
	GitHubIssueNumber int       `json:"github_issue_number" db:"github_issue_number"`
	SubmittedAt       time.Time `json:"submitted_at" db:"submitted_at"`
}

// SponsorUsername represents a synced sponsor username (user or org).
type SponsorUsername struct {
	ID                  int       `json:"id" db:"id"`
	Username            string    `json:"username" db:"username"`
	EntityType          string    `json:"entity_type" db:"entity_type"`
	MonthlyAmountCents  int       `json:"monthly_amount_cents" db:"monthly_amount_cents"`
	TierName            string    `json:"tier_name" db:"tier_name"`
	IsActive            bool      `json:"is_active" db:"is_active"`
	SyncedAt            time.Time `json:"synced_at" db:"synced_at"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}

// getUserByID retrieves a user by ID from the database.
func getUserByID(ctx context.Context, pool *pgxpool.Pool, userID int) (*User, error) {
	slog.Debug("getUserByID: querying user", "user_id", userID)

	var user User
	err := pool.QueryRow(ctx, `
		SELECT id, github_id, login, avatar_url, name, email,
		       sponsorship_data, last_sponsorship_check, created_at, updated_at
		FROM users WHERE id = $1
	`, userID).Scan(
		&user.ID, &user.GitHubID, &user.Login, &user.AvatarURL,
		&user.Name, &user.Email, &user.SponsorshipData,
		&user.LastSponsorshipCheck, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		slog.Error("getUserByID: user not found", "user_id", userID, "err", err)
		return nil, err
	}

	slog.Debug("getUserByID: user found", "user_id", userID, "login", user.Login)
	return &user, nil
}

// upsertUser creates or updates a user in the database.
func upsertUser(ctx context.Context, pool *pgxpool.Pool, user *User) error {
	slog.Debug("upsertUser: attempting upsert", "github_id", user.GitHubID, "login", user.Login)

	// Try update first
	tag, err := pool.Exec(ctx, `
		UPDATE users
		SET login=$1, avatar_url=$2, name=$3, email=$4,
		    sponsorship_data=$5, last_sponsorship_check=NOW(), updated_at=NOW()
		WHERE github_id=$6
	`, user.Login, user.AvatarURL, user.Name, user.Email,
		user.SponsorshipData, user.GitHubID)
	if err != nil {
		slog.Error("upsertUser: update failed", "err", err, "github_id", user.GitHubID)
		return err
	}

	if tag.RowsAffected() > 0 {
		slog.Debug("upsertUser: updated existing user", "github_id", user.GitHubID, "rows_affected", tag.RowsAffected())
		return pool.QueryRow(ctx, `
			SELECT id, github_id, login, avatar_url, name, email,
			       sponsorship_data, last_sponsorship_check, created_at, updated_at
			FROM users WHERE github_id = $1
		`, user.GitHubID).Scan(
			&user.ID, &user.GitHubID, &user.Login, &user.AvatarURL,
			&user.Name, &user.Email, &user.SponsorshipData,
			&user.LastSponsorshipCheck, &user.CreatedAt, &user.UpdatedAt,
		)
	}

	slog.Debug("upsertUser: inserting new user", "github_id", user.GitHubID, "login", user.Login)

	return pool.QueryRow(ctx, `
		INSERT INTO users (github_id, login, avatar_url, name, email, sponsorship_data)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, github_id, login, avatar_url, name, email,
		          sponsorship_data, last_sponsorship_check, created_at, updated_at
	`, user.GitHubID, user.Login, user.AvatarURL, user.Name, user.Email,
		user.SponsorshipData).Scan(
		&user.ID, &user.GitHubID, &user.Login, &user.AvatarURL,
		&user.Name, &user.Email, &user.SponsorshipData,
		&user.LastSponsorshipCheck, &user.CreatedAt, &user.UpdatedAt,
	)
}

// createLogoSubmission creates a logo submission in the database.
func createLogoSubmission(ctx context.Context, pool *pgxpool.Pool, submission *LogoSubmission) error {
	slog.Debug("createLogoSubmission: inserting submission",
		"user_id", submission.UserID,
		"company", submission.CompanyName,
		"issue_number", submission.GitHubIssueNumber)

	err := pool.QueryRow(ctx, `
		INSERT INTO logo_submissions (user_id, company_name, website, logo_url, github_issue_url, github_issue_number)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, submitted_at
	`, submission.UserID, submission.CompanyName, submission.Website,
		submission.LogoURL, submission.GitHubIssueURL, submission.GitHubIssueNumber,
	).Scan(&submission.ID, &submission.SubmittedAt)

	if err != nil {
		slog.Error("createLogoSubmission: failed to insert", "err", err, "user_id", submission.UserID)
		return err
	}

	slog.Debug("createLogoSubmission: submission inserted",
		"user_id", submission.UserID,
		"submission_id", submission.ID,
		"submitted_at", submission.SubmittedAt)

	return nil
}

// getActiveSponsorsByUsernames returns active sponsors matching any of the given usernames.
func getActiveSponsorsByUsernames(ctx context.Context, pool *pgxpool.Pool, usernames []string) ([]*SponsorUsername, error) {
	if len(usernames) == 0 {
		return nil, nil
	}

	slog.Debug("getActiveSponsorsByUsernames: querying sponsors", "usernames", usernames)

	rows, err := pool.Query(ctx, `
		SELECT id, username, entity_type, monthly_amount_cents, tier_name, is_active, synced_at, created_at
		FROM github_sponsor_usernames
		WHERE username = ANY($1) AND is_active = TRUE
		ORDER BY monthly_amount_cents DESC
	`, usernames)
	if err != nil {
		slog.Error("getActiveSponsorsByUsernames: query failed", "err", err)
		return nil, err
	}
	defer rows.Close()

	var sponsors []*SponsorUsername
	for rows.Next() {
		s := &SponsorUsername{}
		if err := rows.Scan(&s.ID, &s.Username, &s.EntityType, &s.MonthlyAmountCents, &s.TierName, &s.IsActive, &s.SyncedAt, &s.CreatedAt); err != nil {
			slog.Error("getActiveSponsorsByUsernames: scan failed", "err", err)
			return nil, err
		}
		sponsors = append(sponsors, s)
	}

	slog.Debug("getActiveSponsorsByUsernames: found sponsors", "count", len(sponsors))
	return sponsors, nil
}

// upsertSponsorUsername inserts or updates a sponsor username.
func upsertSponsorUsername(ctx context.Context, pool *pgxpool.Pool, sponsor *SponsorUsername) error {
	slog.Debug("upsertSponsorUsername: upserting sponsor",
		"username", sponsor.Username,
		"entity_type", sponsor.EntityType,
		"monthly_amount_cents", sponsor.MonthlyAmountCents,
		"tier_name", sponsor.TierName)

	tag, err := pool.Exec(ctx, `
		UPDATE github_sponsor_usernames
		SET entity_type=$1, monthly_amount_cents=$2, tier_name=$3, is_active=$4, synced_at=NOW()
		WHERE username=$5
	`, sponsor.EntityType, sponsor.MonthlyAmountCents, sponsor.TierName, sponsor.IsActive, sponsor.Username)
	if err != nil {
		slog.Error("upsertSponsorUsername: update failed", "err", err, "username", sponsor.Username)
		return err
	}

	if tag.RowsAffected() > 0 {
		slog.Debug("upsertSponsorUsername: updated existing sponsor", "username", sponsor.Username)
		return pool.QueryRow(ctx, `SELECT id, created_at FROM github_sponsor_usernames WHERE username = $1`, sponsor.Username).Scan(&sponsor.ID, &sponsor.CreatedAt)
	}

	err = pool.QueryRow(ctx, `
		INSERT INTO github_sponsor_usernames (username, entity_type, monthly_amount_cents, tier_name, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, synced_at, created_at
	`, sponsor.Username, sponsor.EntityType, sponsor.MonthlyAmountCents, sponsor.TierName, sponsor.IsActive).Scan(&sponsor.ID, &sponsor.SyncedAt, &sponsor.CreatedAt)
	if err != nil {
		slog.Error("upsertSponsorUsername: insert failed", "err", err, "username", sponsor.Username)
		return err
	}

	slog.Debug("upsertSponsorUsername: inserted new sponsor", "username", sponsor.Username, "id", sponsor.ID)
	return nil
}

// markInactiveSponsorsNotIn marks all sponsors as inactive that are not in the given usernames list.
func markInactiveSponsorsNotIn(ctx context.Context, pool *pgxpool.Pool, usernames []string) (int64, error) {
	if len(usernames) == 0 {
		tag, err := pool.Exec(ctx, `UPDATE github_sponsor_usernames SET is_active = FALSE WHERE is_active = TRUE`)
		if err != nil {
			return 0, err
		}
		slog.Debug("markInactiveSponsorsNotIn: marked all sponsors inactive", "count", tag.RowsAffected())
		return tag.RowsAffected(), nil
	}

	tag, err := pool.Exec(ctx, `
		UPDATE github_sponsor_usernames
		SET is_active = FALSE
		WHERE NOT (username = ANY($1)) AND is_active = TRUE
	`, usernames)
	if err != nil {
		slog.Error("markInactiveSponsorsNotIn: update failed", "err", err)
		return 0, err
	}

	slog.Debug("markInactiveSponsorsNotIn: marked sponsors inactive", "count", tag.RowsAffected())
	return tag.RowsAffected(), nil
}
