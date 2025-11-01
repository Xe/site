# Models Package

This package contains GORM models for the GitHub Sponsors webhook processing system, based on the requirements defined in [SITE-10](https://linear.app/xeiaso/issue/SITE-10/github-sponsors-webhooks).

## Models

### Account
Represents a GitHub account (user or organization) that can be either a sponsor or sponsoree.

**Fields:**
- `GitHubID` - GitHub's internal ID (unique)
- `NodeID` - GitHub's global node ID
- `Login` - GitHub username
- `AvatarURL` - Profile avatar URL
- `URL` - GitHub profile URL
- `Type` - "User" or "Organization"
- `SiteAdmin` - Whether the account is a site admin

### Tier
Represents a sponsorship tier with pricing and benefits.

**Fields:**
- `GitHubID` - GitHub's internal ID (unique)
- `Name` - Tier name
- `MonthlyPriceInCents` - Price in cents
- `Description` - Tier description
- `Benefits` - JSON array of tier benefits
- `IsOneTime` - Whether this is a one-time payment
- `IsCustomAmount` - Whether this accepts custom amounts
- `Published` - Whether this tier is publicly visible
- `SelectedTier` - Whether this is the selected tier for the sponsor
- `SponsorshipCount` - Number of active sponsorships

### Sponsorship
Represents a sponsorship relationship between a sponsor and sponsoree.

**Fields:**
- `GitHubID` - GitHub's internal ID (unique)
- `PrivacyLevel` - "public" or "private"
- `Variant` - "recurring" or "one_time"
- `SponsorshipType` - "user" or "organization"
- `TierID` - Foreign key to Tier
- `SponsorID` - Foreign key to Account (sponsor)
- `SponsoreeID` - Foreign key to Account (sponsoree)
- `GitHubCreatedAt` - Creation timestamp from GitHub
- `GitHubUpdatedAt` - Last update timestamp from GitHub
- `GitHubCancelledAt` - Cancellation timestamp from GitHub
- `Metadata` - JSON metadata map

### WebhookEvent
Tracks incoming GitHub Sponsors webhook events for auditing and debugging.

**Fields:**
- `GitHubID` - GitHub's delivery ID (unique)
- `Action` - Event action ("created", "edited", "cancelled", etc.)
- `SenderID` - Account that triggered the event
- `SponsorshipID` - Related sponsorship
- `EventType` - Always "sponsorship" for our use case
- `ProcessedAt` - When we processed the event
- `Success` - Whether processing was successful
- `ErrorMessage` - Error message if processing failed
- `Payload` - Raw webhook payload (JSON)
- `RemoteAddr` - Client IP address
- `UserAgent` - Client user agent

## Usage

### Database Setup

```go
import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "xeiaso.net/v4/internal/models"
)

func main() {
    dsn := "host=localhost user=gorm dbname=gorm port=9920 sslmode=disable TimeZone=UTC"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // Auto-migrate all models
    err = models.SetupDatabase(db)
    if err != nil {
        panic("failed to migrate database")
    }
}
```

### Creating Records

```go
// Create a sponsor
sponsor := models.Account{
    GitHubID:  12345,
    NodeID:    "U_kgDOBQ...",
    Login:     "example_sponsor",
    AvatarURL: "https://avatars.githubusercontent.com/u/12345?v=4",
    URL:       "https://github.com/example_sponsor",
    Type:      models.AccountTypeUser,
}

db.Create(&sponsor)

// Create a tier
tier := models.Tier{
    GitHubID:            67890,
    Name:                "Gold Sponsor",
    MonthlyPriceInCents: 1000, // $10.00
    Description:         "Gold level sponsorship with benefits",
}

tier.SetBenefits([]string{
    "Logo on website",
    "Discord role",
    "Monthly shoutout",
})

db.Create(&tier)

// Create a sponsorship
sponsorship := models.Sponsorship{
    GitHubID:       11111,
    PrivacyLevel:   models.PrivacyLevelPublic,
    Variant:        models.VariantRecurring,
    SponsorshipType: models.SponsorshipTypeUser,
    TierID:         tier.ID,
    SponsorID:      sponsor.ID,
    SponsoreeID:    sponsoree.ID,
}

db.Create(&sponsorship)
```

### Querying Records

```go
// Find active sponsorships for a specific sponsor
var sponsorships []models.Sponsorship
db.Where("sponsor_id = ? AND github_cancelled_at IS NULL", sponsorID).
  Preload("Tier").
  Preload("Sponsor").
  Preload("Sponsoree").
  Find(&sponsorships)

// Count webhooks by action type
var eventCounts []struct {
    Action  string
    Count   int64
}

db.Model(&models.WebhookEvent{}).
  Select("action, count(*) as count").
  Group("action").
  Scan(&eventCounts)
```

## Validation

The package provides validation functions for common fields:

```go
models.IsValidSponsorshipAction("created")  // true
models.IsValidPrivacyLevel("public")       // true
models.IsValidVariant("recurring")        // true
models.IsValidSponsorshipType("user")      // true
models.IsValidAccountType("Organization")  // true
```

## Constraints

- `MaxLoginLength`: 255 characters
- `MaxNameLength`: 255 characters
- `MaxNodeIDLength`: 255 characters
- `MaxAvatarURLLength`: 500 characters
- `MaxURLLength`: 500 characters
- `MaxDescriptionLength`: 1000 characters
- `MaxMetadataSize`: 10000 bytes

## Webhook Integration

The models are designed to work with GitHub Sponsors webhooks. See the [github-sponsor-webhook](../../cmd/github-sponsor-webhook) service for the webhook processing implementation.