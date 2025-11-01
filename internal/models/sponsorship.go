package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Sponsorship represents the sponsorship object.
// Based on the Sponsorship type from SITE-10.
type Sponsorship struct {
	ID           uint            `json:"id" gorm:"primaryKey"`
	GitHubID     int64           `json:"github_id" gorm:"uniqueIndex;not null"` // GitHub's internal ID
	PrivacyLevel string          `json:"privacy_level" gorm:"not null"`         // "public" or "private"
	Variant      string          `json:"variant" gorm:"not null"`              // "recurring" or "one_time"
	SponsorshipType string       `json:"sponsorship_type" gorm:"not null"`      // "user" or "organization"

	// Pricing and tier information
	TierID       uint   `json:"tier_id" gorm:"not null"`       // Foreign key to Tier
	Tier         Tier   `json:"tier" gorm:"foreignKey:TierID"` // Relationship to Tier
	SponsorID    uint   `json:"sponsor_id" gorm:"not null"`    // Foreign key to Account (sponsor)
	Sponsor      Account `json:"sponsor" gorm:"foreignKey:SponsorID"`
	SponsoreeID  uint   `json:"sponsoree_id" gorm:"not null"`  // Foreign key to Account (sponsoree)
	Sponsoree    Account `json:"sponsoree" gorm:"foreignKey:SponsoreeID"`

	// Timestamps from GitHub
	GitHubCreatedAt *time.Time `json:"github_created_at"`
	GitHubUpdatedAt *time.Time `json:"github_updated_at"`
	GitHubCancelledAt *time.Time `json:"github_cancelled_at"`

	// Metadata stored as JSON
	MetadataJSON string `json:"-" gorm:"type:text"` // JSON encoded metadata

	// Local timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships for webhook events
	WebhookEvents []WebhookEvent `json:"-" gorm:"foreignKey:SponsorshipID"`
}

// TableName specifies the table name for the Sponsorship model.
func (Sponsorship) TableName() string {
	return "sponsorships"
}

// Metadata returns the metadata as a map.
func (s *Sponsorship) Metadata() map[string]interface{} {
	if s.MetadataJSON == "" {
		return make(map[string]interface{})
	}

	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(s.MetadataJSON), &metadata); err != nil {
		return make(map[string]interface{})
	}
	return metadata
}

// SetMetadata sets the metadata from a map.
func (s *Sponsorship) SetMetadata(metadata map[string]interface{}) error {
	if metadata == nil {
		s.MetadataJSON = ""
		return nil
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	s.MetadataJSON = string(data)
	return nil
}

// BeforeSave marshals the metadata before saving.
func (s *Sponsorship) BeforeSave(tx *gorm.DB) error {
	return s.SetMetadata(s.Metadata())
}

// BeforeCreate ensures data consistency before creating records.
func (s *Sponsorship) BeforeCreate(tx *gorm.DB) error {
	if s.SponsorID == 0 || s.SponsoreeID == 0 {
		return gorm.ErrInvalidField
	}
	if s.TierID == 0 {
		return gorm.ErrInvalidField
	}
	return nil
}

// IsActive returns true if the sponsorship is active (not cancelled).
func (s *Sponsorship) IsActive() bool {
	return s.GitHubCancelledAt == nil && s.DeletedAt.Time.IsZero()
}

// GetMonthlyPriceInDollars returns the monthly price in dollars.
func (s *Sponsorship) GetMonthlyPriceInDollars() float64 {
	if s.Tier.MonthlyPriceInCents == 0 {
		return 0
	}
	return float64(s.Tier.MonthlyPriceInCents) / 100.0
}