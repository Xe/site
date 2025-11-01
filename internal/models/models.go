package models

import (
	"gorm.io/gorm"
)

// AllModels returns all models that should be auto-migrated.
func AllModels() []interface{} {
	return []interface{}{
		&Account{},
		&Tier{},
		&Sponsorship{},
		&WebhookEvent{},
	}
}

// AutoMigrate runs database auto-migration for all models.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(AllModels()...)
}

// SetupDatabase creates all necessary tables and sets up constraints.
func SetupDatabase(db *gorm.DB) error {
	// Run auto-migration
	if err := AutoMigrate(db); err != nil {
		return err
	}

	// Add any additional constraints or indexes here
	// For example, we might want to add indexes for common queries

	return nil
}

// Common validation constraints
const (
	MaxLoginLength    = 255
	MaxNameLength     = 255
	MaxNodeIDLength   = 255
	MaxAvatarURLLength = 500
	MaxURLLength      = 500
	MaxDescriptionLength = 1000
	MaxMetadataSize   = 10000 // Maximum size for JSON metadata in bytes
)

// Sponsorship actions from GitHub Sponsors webhooks
const (
	SponsorshipActionCreated            = "created"
	SponsorshipActionEdited             = "edited"
	SponsorshipActionCancelled          = "cancelled"
	SponsorshipActionPendingTierChange  = "pending_tier_change"
	SponsorshipActionPendingCancellation = "pending_cancellation"
)

// Privacy levels for sponsorships
const (
	PrivacyLevelPublic  = "public"
	PrivacyLevelPrivate = "private"
)

// Sponsorship variants
const (
	VariantRecurring = "recurring"
	VariantOneTime   = "one_time"
)

// Sponsorship types
const (
	SponsorshipTypeUser         = "user"
	SponsorshipTypeOrganization = "organization"
)

// Account types
const (
	AccountTypeUser         = "User"
	AccountTypeOrganization = "Organization"
)

// IsValidSponsorshipAction checks if the action is valid.
func IsValidSponsorshipAction(action string) bool {
	validActions := []string{
		SponsorshipActionCreated,
		SponsorshipActionEdited,
		SponsorshipActionCancelled,
		SponsorshipActionPendingTierChange,
		SponsorshipActionPendingCancellation,
	}

	for _, validAction := range validActions {
		if action == validAction {
			return true
		}
	}
	return false
}

// IsValidPrivacyLevel checks if the privacy level is valid.
func IsValidPrivacyLevel(level string) bool {
	return level == PrivacyLevelPublic || level == PrivacyLevelPrivate
}

// IsValidVariant checks if the variant is valid.
func IsValidVariant(variant string) bool {
	return variant == VariantRecurring || variant == VariantOneTime
}

// IsValidSponsorshipType checks if the sponsorship type is valid.
func IsValidSponsorshipType(sponsorshipType string) bool {
	return sponsorshipType == SponsorshipTypeUser || sponsorshipType == SponsorshipTypeOrganization
}

// IsValidAccountType checks if the account type is valid.
func IsValidAccountType(accountType string) bool {
	return accountType == AccountTypeUser || accountType == AccountTypeOrganization
}