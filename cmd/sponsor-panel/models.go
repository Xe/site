package main

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"time"
)

// User represents a GitHub user with cached sponsorship data.
// This is the simplified model from SPEC.md using sqlx (not GORM).
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

// getUserByID retrieves a user by ID from the database.
func getUserByID(db *sql.DB, userID int) (*User, error) {
	slog.Debug("getUserByID: querying user", "user_id", userID)

	var user User
	err := db.QueryRow(`
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

// getUserByGitHubID retrieves a user by GitHub ID from the database.
func getUserByGitHubID(db *sql.DB, githubID int64) (*User, error) {
	var user User
	err := db.QueryRow(`
		SELECT id, github_id, login, avatar_url, name, email,
		       sponsorship_data, last_sponsorship_check, created_at, updated_at
		FROM users WHERE github_id = $1
	`, githubID).Scan(
		&user.ID, &user.GitHubID, &user.Login, &user.AvatarURL,
		&user.Name, &user.Email, &user.SponsorshipData,
		&user.LastSponsorshipCheck, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// upsertUser creates or updates a user in the database.
func upsertUser(db *sql.DB, user *User) error {
	slog.Debug("upsertUser: attempting upsert", "github_id", user.GitHubID, "login", user.Login)

	// Try update first
	result, err := db.Exec(`
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

	rows, _ := result.RowsAffected()
	if rows > 0 {
		slog.Debug("upsertUser: updated existing user", "github_id", user.GitHubID, "rows_affected", rows)
		// Fetch the updated user
		return db.QueryRow(`
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

	// Insert new user
	return db.QueryRow(`
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
func createLogoSubmission(db *sql.DB, submission *LogoSubmission) error {
	slog.Debug("createLogoSubmission: inserting submission",
		"user_id", submission.UserID,
		"company", submission.CompanyName,
		"issue_number", submission.GitHubIssueNumber)

	err := db.QueryRow(`
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
