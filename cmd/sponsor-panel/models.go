package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// User represents an authenticated user with cached sponsorship data.
type User struct {
	ID                   int       `json:"id" db:"id"`
	GitHubID             *int64    `json:"github_id" db:"github_id"`
	PatreonID            *string   `json:"patreon_id" db:"patreon_id"`
	GoogleID             *string   `json:"google_id" db:"google_id"`
	MicrosoftID          *string   `json:"microsoft_id" db:"microsoft_id"`
	StripeCustomerID     *string   `json:"stripe_customer_id" db:"stripe_customer_id"`
	Provider             string    `json:"provider" db:"provider"` // "github", "patreon", "google", "microsoft", or "email"
	Login                string    `json:"login" db:"login"`
	AvatarURL            string    `json:"avatar_url" db:"avatar_url"`
	Name                 string    `json:"name" db:"name"`
	Email                string    `json:"email" db:"email"`
	SponsorshipData      string    `json:"-" db:"sponsorship_data"` // JSON blob
	LastSponsorshipCheck time.Time `json:"last_sponsorship_check" db:"last_sponsorship_check"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// MagicLinkToken represents a passwordless email login token.
type MagicLinkToken struct {
	ID        int        `json:"id" db:"id"`
	Email     string     `json:"email" db:"email"`
	TokenHash string     `json:"token_hash" db:"token_hash"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	UsedAt    *time.Time `json:"used_at" db:"used_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
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

// scanUser scans a user row into a User struct. The query must SELECT columns
// in the order returned by userQuery().
func scanUser(row interface{ Scan(...any) error }) (*User, error) {
	var user User
	err := row.Scan(
		&user.ID, &user.GitHubID, &user.PatreonID, &user.GoogleID, &user.MicrosoftID, &user.StripeCustomerID,
		&user.Provider, &user.Login, &user.AvatarURL, &user.Name, &user.Email,
		&user.SponsorshipData, &user.LastSponsorshipCheck, &user.CreatedAt, &user.UpdatedAt,
	)
	return &user, err
}

// getUserByID retrieves a user by ID from the database.
func getUserByID(ctx context.Context, pool *pgxpool.Pool, userID int) (*User, error) {
	slog.Debug("getUserByID: querying user", "user_id", userID)

	user, err := scanUser(pool.QueryRow(ctx, `
		SELECT id, github_id, patreon_id, google_id, microsoft_id, stripe_customer_id,
		       provider, login, avatar_url, name, email,
		       sponsorship_data, last_sponsorship_check, created_at, updated_at
		FROM users WHERE id = $1
	`, userID))
	if err != nil {
		slog.Error("getUserByID: user not found", "user_id", userID, "err", err)
		return nil, err
	}

	slog.Debug("getUserByID: user found", "user_id", userID, "login", user.Login)
	return user, nil
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
		found, scanErr := scanUser(pool.QueryRow(ctx, `
			SELECT id, github_id, patreon_id, google_id, microsoft_id, stripe_customer_id,
			       provider, login, avatar_url, name, email,
			       sponsorship_data, last_sponsorship_check, created_at, updated_at
			FROM users WHERE github_id = $1
		`, user.GitHubID))
		if scanErr != nil {
			return scanErr
		}
		*user = *found
		return nil
	}

	slog.Debug("upsertUser: inserting new user", "github_id", user.GitHubID, "login", user.Login)

	found, scanErr := scanUser(pool.QueryRow(ctx, `
		INSERT INTO users (github_id, provider, login, avatar_url, name, email, sponsorship_data)
		VALUES ($1, 'github', $2, $3, $4, $5, $6)
		RETURNING id, github_id, patreon_id, google_id, microsoft_id, stripe_customer_id,
		          provider, login, avatar_url, name, email,
		          sponsorship_data, last_sponsorship_check, created_at, updated_at
	`, user.GitHubID, user.Login, user.AvatarURL, user.Name, user.Email,
		user.SponsorshipData))
	if scanErr != nil {
		return scanErr
	}
	*user = *found
	return nil
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

// upsertPatreonUser creates or updates a Patreon user in the database.
func upsertPatreonUser(ctx context.Context, pool *pgxpool.Pool, user *User) error {
	slog.Debug("upsertPatreonUser: attempting upsert", "patreon_id", user.PatreonID, "login", user.Login)

	// Try update first
	tag, err := pool.Exec(ctx, `
		UPDATE users
		SET login=$1, avatar_url=$2, name=$3, email=$4,
		    sponsorship_data=$5, last_sponsorship_check=NOW(), updated_at=NOW()
		WHERE patreon_id=$6
	`, user.Login, user.AvatarURL, user.Name, user.Email,
		user.SponsorshipData, user.PatreonID)
	if err != nil {
		slog.Error("upsertPatreonUser: update failed", "err", err, "patreon_id", user.PatreonID)
		return err
	}

	if tag.RowsAffected() > 0 {
		slog.Debug("upsertPatreonUser: updated existing user", "patreon_id", user.PatreonID, "rows_affected", tag.RowsAffected())
		found, scanErr := scanUser(pool.QueryRow(ctx, `
			SELECT id, github_id, patreon_id, google_id, microsoft_id, stripe_customer_id,
			       provider, login, avatar_url, name, email,
			       sponsorship_data, last_sponsorship_check, created_at, updated_at
			FROM users WHERE patreon_id = $1
		`, user.PatreonID))
		if scanErr != nil {
			return scanErr
		}
		*user = *found
		return nil
	}

	slog.Debug("upsertPatreonUser: inserting new user", "patreon_id", user.PatreonID, "login", user.Login)

	found, scanErr := scanUser(pool.QueryRow(ctx, `
		INSERT INTO users (patreon_id, provider, login, avatar_url, name, email, sponsorship_data)
		VALUES ($1, 'patreon', $2, $3, $4, $5, $6)
		RETURNING id, github_id, patreon_id, google_id, microsoft_id, stripe_customer_id,
		          provider, login, avatar_url, name, email,
		          sponsorship_data, last_sponsorship_check, created_at, updated_at
	`, user.PatreonID, user.Login, user.AvatarURL, user.Name, user.Email,
		user.SponsorshipData))
	if scanErr != nil {
		return scanErr
	}
	*user = *found
	return nil
}

// upsertGoogleUser creates or updates a Google OAuth user in the database.
func upsertGoogleUser(ctx context.Context, pool *pgxpool.Pool, user *User) error {
	slog.Debug("upsertGoogleUser: attempting upsert", "google_id", user.GoogleID, "login", user.Login)

	tag, err := pool.Exec(ctx, `
		UPDATE users
		SET login=$1, avatar_url=$2, name=$3, email=$4,
		    sponsorship_data=$5, last_sponsorship_check=NOW(), updated_at=NOW()
		WHERE google_id=$6
	`, user.Login, user.AvatarURL, user.Name, user.Email,
		user.SponsorshipData, user.GoogleID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() > 0 {
		found, scanErr := scanUser(pool.QueryRow(ctx, `
			SELECT id, github_id, patreon_id, google_id, microsoft_id, stripe_customer_id,
			       provider, login, avatar_url, name, email,
			       sponsorship_data, last_sponsorship_check, created_at, updated_at
			FROM users WHERE google_id = $1
		`, user.GoogleID))
		if scanErr != nil {
			return scanErr
		}
		*user = *found
		return nil
	}

	found, scanErr := scanUser(pool.QueryRow(ctx, `
		INSERT INTO users (google_id, provider, login, avatar_url, name, email, sponsorship_data)
		VALUES ($1, 'google', $2, $3, $4, $5, $6)
		RETURNING id, github_id, patreon_id, google_id, microsoft_id, stripe_customer_id,
		          provider, login, avatar_url, name, email,
		          sponsorship_data, last_sponsorship_check, created_at, updated_at
	`, user.GoogleID, user.Login, user.AvatarURL, user.Name, user.Email,
		user.SponsorshipData))
	if scanErr != nil {
		return scanErr
	}
	*user = *found
	return nil
}

// upsertMicrosoftUser creates or updates a Microsoft OAuth user in the database.
func upsertMicrosoftUser(ctx context.Context, pool *pgxpool.Pool, user *User) error {
	slog.Debug("upsertMicrosoftUser: attempting upsert", "microsoft_id", user.MicrosoftID, "login", user.Login)

	tag, err := pool.Exec(ctx, `
		UPDATE users
		SET login=$1, avatar_url=$2, name=$3, email=$4,
		    sponsorship_data=$5, last_sponsorship_check=NOW(), updated_at=NOW()
		WHERE microsoft_id=$6
	`, user.Login, user.AvatarURL, user.Name, user.Email,
		user.SponsorshipData, user.MicrosoftID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() > 0 {
		found, scanErr := scanUser(pool.QueryRow(ctx, `
			SELECT id, github_id, patreon_id, google_id, microsoft_id, stripe_customer_id,
			       provider, login, avatar_url, name, email,
			       sponsorship_data, last_sponsorship_check, created_at, updated_at
			FROM users WHERE microsoft_id = $1
		`, user.MicrosoftID))
		if scanErr != nil {
			return scanErr
		}
		*user = *found
		return nil
	}

	found, scanErr := scanUser(pool.QueryRow(ctx, `
		INSERT INTO users (microsoft_id, provider, login, avatar_url, name, email, sponsorship_data)
		VALUES ($1, 'microsoft', $2, $3, $4, $5, $6)
		RETURNING id, github_id, patreon_id, google_id, microsoft_id, stripe_customer_id,
		          provider, login, avatar_url, name, email,
		          sponsorship_data, last_sponsorship_check, created_at, updated_at
	`, user.MicrosoftID, user.Login, user.AvatarURL, user.Name, user.Email,
		user.SponsorshipData))
	if scanErr != nil {
		return scanErr
	}
	*user = *found
	return nil
}

// upsertEmailUser creates or updates an email magic link user in the database.
func upsertEmailUser(ctx context.Context, pool *pgxpool.Pool, user *User) error {
	slog.Debug("upsertEmailUser: attempting upsert", "email", user.Email, "login", user.Login)

	tag, err := pool.Exec(ctx, `
		UPDATE users
		SET avatar_url=$1, name=$2,
		    last_sponsorship_check=NOW(), updated_at=NOW()
		WHERE provider='email' AND login=$3
	`, user.AvatarURL, user.Name, user.Login)
	if err != nil {
		return err
	}

	if tag.RowsAffected() > 0 {
		found, scanErr := scanUser(pool.QueryRow(ctx, `
			SELECT id, github_id, patreon_id, google_id, microsoft_id, stripe_customer_id,
			       provider, login, avatar_url, name, email,
			       sponsorship_data, last_sponsorship_check, created_at, updated_at
			FROM users WHERE provider='email' AND login = $1
		`, user.Login))
		if scanErr != nil {
			return scanErr
		}
		*user = *found
		return nil
	}

	found, scanErr := scanUser(pool.QueryRow(ctx, `
		INSERT INTO users (provider, login, avatar_url, name, email, sponsorship_data)
		VALUES ('email', $1, $2, $3, $4, $5)
		RETURNING id, github_id, patreon_id, google_id, microsoft_id, stripe_customer_id,
		          provider, login, avatar_url, name, email,
		          sponsorship_data, last_sponsorship_check, created_at, updated_at
	`, user.Login, user.AvatarURL, user.Name, user.Email,
		user.SponsorshipData))
	if scanErr != nil {
		return scanErr
	}
	*user = *found
	return nil
}

// getUserByStripeCustomerID retrieves a user by their Stripe customer ID.
func getUserByStripeCustomerID(ctx context.Context, pool *pgxpool.Pool, stripeID string) (*User, error) {
	user, err := scanUser(pool.QueryRow(ctx, `
		SELECT id, github_id, patreon_id, google_id, microsoft_id, stripe_customer_id,
		       provider, login, avatar_url, name, email,
		       sponsorship_data, last_sponsorship_check, created_at, updated_at
		FROM users WHERE stripe_customer_id = $1
	`, stripeID))
	if err != nil {
		return nil, err
	}
	return user, nil
}

// setStripeCustomerID links a Stripe customer ID to a user.
func setStripeCustomerID(ctx context.Context, pool *pgxpool.Pool, userID int, stripeID string) error {
	_, err := pool.Exec(ctx, `UPDATE users SET stripe_customer_id=$1, updated_at=NOW() WHERE id=$2`, stripeID, userID)
	return err
}

// createMagicLinkToken stores a magic link token hash in the database.
func createMagicLinkToken(ctx context.Context, pool *pgxpool.Pool, email, tokenHash string, expiresAt time.Time) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO magic_link_tokens (email, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, email, tokenHash, expiresAt)
	return err
}

// consumeMagicLinkToken marks a token as used and returns it if valid and unexpired.
func consumeMagicLinkToken(ctx context.Context, pool *pgxpool.Pool, tokenHash string) (*MagicLinkToken, error) {
	var token MagicLinkToken
	err := pool.QueryRow(ctx, `
		UPDATE magic_link_tokens
		SET used_at = NOW()
		WHERE token_hash = $1 AND used_at IS NULL AND expires_at > NOW()
		RETURNING id, email, token_hash, expires_at, used_at, created_at
	`, tokenHash).Scan(&token.ID, &token.Email, &token.TokenHash, &token.ExpiresAt, &token.UsedAt, &token.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// countRecentMagicLinks returns the number of magic link tokens created for an email in the last hour.
func countRecentMagicLinks(ctx context.Context, pool *pgxpool.Pool, email string) (int, error) {
	var count int
	err := pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM magic_link_tokens
		WHERE email = $1 AND created_at > NOW() - INTERVAL '1 hour'
	`, email).Scan(&count)
	return count, err
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
