package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"gorm.io/gorm"
)

// PanelUser represents an authenticated user with cached sponsorship data.
type PanelUser struct {
	ID uint `json:"id" gorm:"primaryKey"`
	// The column tag is mandatory. GORM derives git_hub_id from the field
	// name, which does not match the column the schema has always used.
	GitHubID             *int64    `json:"github_id" gorm:"column:github_id;uniqueIndex:users_github_id_key"`
	PatreonID            *string   `json:"patreon_id" gorm:"uniqueIndex:users_patreon_id_key"`
	Provider             string    `json:"provider" gorm:"not null;default:'github';uniqueIndex:idx_users_provider_login"`
	Login                string    `json:"login" gorm:"not null;uniqueIndex:idx_users_provider_login"`
	AvatarURL            string    `json:"avatar_url"`
	Name                 string    `json:"name"`
	Email                string    `json:"email"`
	ThothUserID          *string   `json:"thoth_user_id" gorm:"column:thoth_user_id"`
	SponsorshipData      string    `json:"-" gorm:"type:jsonb"`
	LastSponsorshipCheck time.Time `json:"last_sponsorship_check"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func (PanelUser) TableName() string { return "users" }

// SponsorshipData represents the cached GraphQL response.
type SponsorshipData struct {
	IsActive      bool   `json:"is_active"`
	MonthlyAmount int    `json:"monthly_amount_cents"`
	TierName      string `json:"tier_name"`
	PrivacyLevel  string `json:"privacy_level"`
}

// IsSponsorAtTier returns true if user sponsors at or above the given amount (in cents).
func (u *PanelUser) IsSponsorAtTier(minCents int) bool {
	if u.SponsorshipData == "" {
		return false
	}

	var data SponsorshipData
	if err := json.Unmarshal([]byte(u.SponsorshipData), &data); err != nil {
		slog.Error("IsSponsorAtTier: failed to parse sponsorship data", "user_id", u.ID, "err", err)
		return false
	}

	return data.IsActive && data.MonthlyAmount >= minCents
}

// LogoSubmission represents a logo submission.
type LogoSubmission struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	UserID            uint      `json:"user_id" gorm:"not null"`
	CompanyName       string    `json:"company_name" gorm:"not null"`
	Website           string    `json:"website" gorm:"not null"`
	LogoURL           string    `json:"logo_url"`
	GitHubIssueURL    string    `json:"github_issue_url"`
	GitHubIssueNumber int       `json:"github_issue_number"`
	SubmittedAt       time.Time `json:"submitted_at" gorm:"autoCreateTime"`
}

func (LogoSubmission) TableName() string { return "logo_submissions" }

// SponsorUsername represents a synced sponsor username (user or org).
type SponsorUsername struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	Username           string    `json:"username" gorm:"uniqueIndex:github_sponsor_usernames_username_key;not null"`
	EntityType         string    `json:"entity_type" gorm:"not null"`
	MonthlyAmountCents int       `json:"monthly_amount_cents" gorm:"default:0"`
	TierName           string    `json:"tier_name"`
	IsActive           bool      `json:"is_active" gorm:"default:true;index"`
	SyncedAt           time.Time `json:"synced_at"`
	CreatedAt          time.Time `json:"created_at"`
}

func (SponsorUsername) TableName() string { return "github_sponsor_usernames" }

// PanelModels returns all sponsor-panel models for AutoMigrate.
func PanelModels() []any {
	return []any{
		&PanelUser{},
		&LogoSubmission{},
		&SponsorUsername{},
	}
}

// consolidateGitHubIDColumn merges the legacy users.git_hub_id column back
// into users.github_id.
//
// The schema has always called this column github_id. GORM derives git_hub_id
// from the field name GitHubID, so an early AutoMigrate added a second column
// and writes went to one column while reads went to the other. Move the values
// across, then take the original name back.
//
// Run this before AutoMigrate. It is a no-op once the columns are merged.
func consolidateGitHubIDColumn(db *gorm.DB) error {
	m := db.Migrator()

	if !m.HasTable(&PanelUser{}) {
		return nil
	}

	if !m.HasColumn(&PanelUser{}, "git_hub_id") {
		return nil
	}

	hasLegacy := m.HasColumn(&PanelUser{}, "github_id")

	return db.Transaction(func(tx *gorm.DB) error {
		if hasLegacy {
			res := tx.Exec(`UPDATE users SET git_hub_id = github_id WHERE git_hub_id IS NULL AND github_id IS NOT NULL`)
			if res.Error != nil {
				return fmt.Errorf("cannot backfill users.git_hub_id: %w", res.Error)
			}
			slog.Info("consolidateGitHubIDColumn: backfilled GitHub IDs", "rows", res.RowsAffected)

			if err := tx.Exec(`ALTER TABLE users DROP COLUMN github_id`).Error; err != nil {
				return fmt.Errorf("cannot drop the legacy users.github_id column: %w", err)
			}
		}

		if err := tx.Exec(`ALTER TABLE users RENAME COLUMN git_hub_id TO github_id`).Error; err != nil {
			return fmt.Errorf("cannot rename users.git_hub_id to github_id: %w", err)
		}

		slog.Info("consolidateGitHubIDColumn: users.git_hub_id is now users.github_id")
		return nil
	})
}

// emailUniqueIndex describes a legacy unique index on users(email).
type emailUniqueIndex struct {
	IndexName      string
	ConstraintName string
}

// dropEmailUniqueIndex removes any legacy unique index on users(email).
//
// AutoMigrate never drops indexes, so a unique index created by an older
// schema stays in the database forever. Look the index up in the catalog
// because the name depends on how the old schema declared it.
func dropEmailUniqueIndex(db *gorm.DB) error {
	const findQuery = `
SELECT i.relname AS index_name, COALESCE(c.conname, '') AS constraint_name
FROM pg_index x
JOIN pg_class i ON i.oid = x.indexrelid
JOIN pg_class t ON t.oid = x.indrelid
JOIN pg_namespace n ON n.oid = t.relnamespace
LEFT JOIN pg_constraint c ON c.conindid = x.indexrelid
WHERE t.relname = 'users'
  AND n.nspname = current_schema()
  AND x.indisunique
  AND x.indnatts = 1
  AND (SELECT a.attname FROM pg_attribute a WHERE a.attrelid = t.oid AND a.attnum = x.indkey[0]) = 'email'`

	var found []emailUniqueIndex
	if err := db.Raw(findQuery).Scan(&found).Error; err != nil {
		return fmt.Errorf("cannot look up unique indexes on users(email): %w", err)
	}

	for _, idx := range found {
		var stmt string
		switch {
		case idx.ConstraintName != "":
			stmt = fmt.Sprintf("ALTER TABLE users DROP CONSTRAINT %s", quoteIdent(idx.ConstraintName))
		default:
			stmt = fmt.Sprintf("DROP INDEX %s", quoteIdent(idx.IndexName))
		}

		if err := db.Exec(stmt).Error; err != nil {
			return fmt.Errorf("cannot drop unique index %s on users(email): %w", idx.IndexName, err)
		}

		slog.Info("dropEmailUniqueIndex: dropped unique index on users(email)", "index", idx.IndexName, "constraint", idx.ConstraintName)
	}

	return nil
}

// quoteIdent quotes a SQL identifier for PostgreSQL.
func quoteIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// --- DB helper functions (GORM) ---

// getUserByID retrieves a user by ID from the database.
func getUserByID(db *gorm.DB, userID int) (*PanelUser, error) {
	var user PanelUser
	if err := db.First(&user, userID).Error; err != nil {
		slog.Error("getUserByID: user not found", "user_id", userID, "err", err)
		return nil, err
	}
	return &user, nil
}

// upsertUser creates or updates a GitHub user in the database.
func upsertUser(db *gorm.DB, user *PanelUser) error {
	if user.GitHubID == nil {
		return fmt.Errorf("upsertUser: user %q has no GitHub ID", user.Login)
	}

	var existing PanelUser
	// Match on the struct field, not a literal column name. The field decides
	// which column both this query and Create use, so the two cannot drift.
	result := db.Where(&PanelUser{GitHubID: user.GitHubID}).First(&existing)
	if result.Error == nil {
		// Update existing
		existing.Login = user.Login
		existing.AvatarURL = user.AvatarURL
		existing.Name = user.Name
		existing.Email = user.Email
		existing.SponsorshipData = user.SponsorshipData
		existing.LastSponsorshipCheck = time.Now()
		if err := db.Save(&existing).Error; err != nil {
			return err
		}
		*user = existing
		return nil
	}

	// Insert new
	user.Provider = "github"
	return db.Create(user).Error
}

// upsertPatreonUser creates or updates a Patreon user in the database.
func upsertPatreonUser(db *gorm.DB, user *PanelUser) error {
	var existing PanelUser
	result := db.Where("patreon_id = ?", user.PatreonID).First(&existing)
	if result.Error == nil {
		existing.Login = user.Login
		existing.AvatarURL = user.AvatarURL
		existing.Name = user.Name
		existing.Email = user.Email
		existing.SponsorshipData = user.SponsorshipData
		existing.LastSponsorshipCheck = time.Now()
		if err := db.Save(&existing).Error; err != nil {
			return err
		}
		*user = existing
		return nil
	}

	user.Provider = "patreon"
	return db.Create(user).Error
}

// createLogoSubmission creates a logo submission in the database.
func createLogoSubmission(db *gorm.DB, submission *LogoSubmission) error {
	return db.Create(submission).Error
}

// getActiveSponsorsByUsernames returns active sponsors matching any of the given usernames.
func getActiveSponsorsByUsernames(db *gorm.DB, usernames []string) ([]*SponsorUsername, error) {
	if len(usernames) == 0 {
		return nil, nil
	}

	var sponsors []*SponsorUsername
	err := db.Where("username IN ? AND is_active = ?", usernames, true).
		Order("monthly_amount_cents DESC").
		Find(&sponsors).Error
	return sponsors, err
}

// upsertSponsorUsername inserts or updates a sponsor username.
func upsertSponsorUsername(db *gorm.DB, sponsor *SponsorUsername) error {
	var existing SponsorUsername
	result := db.Where("username = ?", sponsor.Username).First(&existing)
	if result.Error == nil {
		existing.EntityType = sponsor.EntityType
		existing.MonthlyAmountCents = sponsor.MonthlyAmountCents
		existing.TierName = sponsor.TierName
		existing.IsActive = sponsor.IsActive
		existing.SyncedAt = time.Now()
		if err := db.Save(&existing).Error; err != nil {
			return err
		}
		sponsor.ID = existing.ID
		sponsor.CreatedAt = existing.CreatedAt
		return nil
	}

	return db.Create(sponsor).Error
}

// markInactiveSponsorsNotIn marks all sponsors as inactive that are not in the given usernames list.
func markInactiveSponsorsNotIn(db *gorm.DB, usernames []string) (int64, error) {
	var result *gorm.DB
	if len(usernames) == 0 {
		result = db.Model(&SponsorUsername{}).Where("is_active = ?", true).Update("is_active", false)
	} else {
		result = db.Model(&SponsorUsername{}).
			Where("username NOT IN ? AND is_active = ?", usernames, true).
			Update("is_active", false)
	}
	return result.RowsAffected, result.Error
}
