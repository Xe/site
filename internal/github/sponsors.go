package github

// Sponsor represents a GitHub sponsor.
type Sponsor struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

// Sponsorship represents a GitHub sponsorship.
type Sponsorship struct {
	ID              int     `json:"id"`
	NodeID          string  `json:"node_id"`
	CreatedAt       Time    `json:"created_at"`
	Sponsor         Sponsor `json:"sponsor"`
	Sponsorable     Sponsor `json:"sponsorable"`
	Tier            Tier    `json:"tier"`
	PrivacyLevel    string  `json:"privacy_level"`
	Variant         string  `json:"variant"`
	SponsorshipType string  `json:"sponsorship_type"`
}

// Tier represents the sponsorship tier.
type Tier struct {
	ID                    int    `json:"id"`
	NodeID                string `json:"node_id"`
	CreatedAt             Time   `json:"created_at"`
	Description           string `json:"description"`
	MonthlyPriceInDollars int    `json:"monthly_price_in_dollars"`
	IsOneTime             bool   `json:"is_one_time"`
	IsCustomAmount        bool   `json:"is_custom_amount"`
	Name                  string `json:"name"`
	Published             bool   `json:"published"`
	SelectedTier          bool   `json:"selected_tier"`
	SponsorshipCount      int    `json:"sponsorship_count"`
}

// SponsorsEvent represents a GitHub Sponsors webhook event.
type SponsorsEvent struct {
	Action       string      `json:"action"`
	Sponsorship  Sponsorship `json:"sponsorship"`
	Sender       User        `json:"sender"`
	Repository   Repository  `json:"repository,omitempty"`
	Organization User        `json:"organization,omitempty"`
}

// GithubSponsorsWebhookEventTypes contains the possible GitHub Sponsors webhook event types.
const (
	SponsorsEventCreated             = "created"
	SponsorsEventEdited              = "edited"
	SponsorsEventCancelled           = "cancelled"
	SponsorsEventPendingTierChange   = "pending_tier_change"
	SponsorsEventPendingCancellation = "pending_cancellation"
)
