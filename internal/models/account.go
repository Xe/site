package models

import (
	"time"

	"gorm.io/gorm"
)

// Account represents a sponsoring or sponsored account (user or org).
// Based on the Account type from SITE-10.
type Account struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	GitHubID  int64          `json:"github_id" gorm:"uniqueIndex;not null"` // GitHub's internal ID
	NodeID    string         `json:"node_id" gorm:"not null"`                // GitHub's global node ID
	Login     string         `json:"login" gorm:"not null"`                  // GitHub username
	AvatarURL string         `json:"avatar_url"`                             // Avatar URL
	URL       string         `json:"url"`                                    // GitHub profile URL
	Type      string         `json:"type" gorm:"not null"`                   // "User" or "Organization"
	SiteAdmin bool           `json:"site_admin"`                             // Whether the account is a site admin

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	SponsorshipsAsSponsor []Sponsorship `json:"-" gorm:"foreignKey:SponsorID"`
	SponsorshipsAsSponsoree []Sponsorship `json:"-" gorm:"foreignKey:SponsoreeID"`
}

// TableName specifies the table name for the Account model.
func (Account) TableName() string {
	return "accounts"
}

// BeforeCreate ensures data consistency before creating records.
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	// Ensure login is not empty
	if a.Login == "" {
		return gorm.ErrInvalidField
	}
	return nil
}