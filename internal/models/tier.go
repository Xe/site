package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Tier represents a sponsorship tier.
// Based on the Tier type from SITE-10.
type Tier struct {
	ID                  uint            `json:"id" gorm:"primaryKey"`
	GitHubID            int64           `json:"github_id" gorm:"uniqueIndex;not null"` // GitHub's internal ID
	Name                string          `json:"name" gorm:"not null"`                   // Tier name
	MonthlyPriceInCents int             `json:"monthly_price_in_cents" gorm:"not null"` // Price in cents
	Description         string          `json:"description"`                           // Tier description
	IsOneTime           bool            `json:"is_one_time" gorm:"default:false"`       // Whether this is a one-time payment
	IsCustomAmount      bool            `json:"is_custom_amount" gorm:"default:false"`  // Whether this accepts custom amounts
	Published           bool            `json:"published" gorm:"default:false"`         // Whether this tier is publicly visible
	SelectedTier        bool            `json:"selected_tier" gorm:"default:false"`     // Whether this is the selected tier for the sponsor
	SponsorshipCount    int             `json:"sponsorship_count" gorm:"default:0"`      // Number of active sponsorships

	// Benefits stored as JSON array
	BenefitsJSON string `json:"-" gorm:"type:text"` // JSON encoded benefits

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Sponsorships []Sponsorship `json:"-" gorm:"foreignKey:TierID"`
}

// TableName specifies the table name for the Tier model.
func (Tier) TableName() string {
	return "tiers"
}

// Benefits returns the benefits as a slice of strings.
func (t *Tier) Benefits() []string {
	if t.BenefitsJSON == "" {
		return []string{}
	}

	var benefits []string
	if err := json.Unmarshal([]byte(t.BenefitsJSON), &benefits); err != nil {
		return []string{}
	}
	return benefits
}

// SetBenefits sets the benefits from a slice of strings.
func (t *Tier) SetBenefits(benefits []string) error {
	if benefits == nil || len(benefits) == 0 {
		t.BenefitsJSON = ""
		return nil
	}

	data, err := json.Marshal(benefits)
	if err != nil {
		return err
	}
	t.BenefitsJSON = string(data)
	return nil
}

// BeforeSave marshals the benefits before saving.
func (t *Tier) BeforeSave(tx *gorm.DB) error {
	return t.SetBenefits(t.Benefits())
}

// BeforeCreate ensures data consistency before creating records.
func (t *Tier) BeforeCreate(tx *gorm.DB) error {
	if t.Name == "" {
		return gorm.ErrInvalidField
	}
	if t.MonthlyPriceInCents < 0 {
		return gorm.ErrInvalidField
	}
	return nil
}