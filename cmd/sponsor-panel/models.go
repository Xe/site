package main

import (
	"encoding/json"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

// PanelUser represents an authenticated user with cached sponsorship data.
type PanelUser struct {
	ID                   uint      `json:"id" gorm:"primaryKey"`
	GitHubID             *int64    `json:"github_id" gorm:"uniqueIndex:users_github_id_key"`
	PatreonID            *string   `json:"patreon_id" gorm:"uniqueIndex:users_patreon_id_key"`
	Provider             string    `json:"provider" gorm:"not null;default:'github';uniqueIndex:idx_users_provider_login"`
	Login                string    `json:"login" gorm:"not null;uniqueIndex:idx_users_provider_login"`
	AvatarURL            string    `json:"avatar_url"`
	Name                 string    `json:"name"`
	Email                string    `json:"email"`
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
func PanelModels() []interface{} {
	return []interface{}{
		&PanelUser{},
		&LogoSubmission{},
		&SponsorUsername{},
	}
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
	var existing PanelUser
	result := db.Where("github_id = ?", user.GitHubID).First(&existing)
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
