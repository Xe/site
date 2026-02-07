# Sponsor Panel Service Specification

## Overview

The **sponsor-panel** service is a web application that provides GitHub sponsors with management capabilities based on their sponsorship tier. The service will be exposed at `sponsors.xeiaso.net` and integrates with GitHub OAuth, GitHub Sponsors API, and the TecharoHQ organization.

**Tech Stack:**

- **Backend:** Go with `net/http`
- **Frontend:** Templ + HTMX + Tailwind CSS
- **Database:** PostgreSQL with Gorm
- **Storage:** Tigris (via AWS SDK v2)
- **GitHub Integration:** go-github v82, GraphQL for sponsors data

---

## Features

### 1. Discord Invite Link

- Display a Discord invite link configured via flag/envvar
- Pattern: `DISCORD_INVITE` environment variable

### 2. Team Invitation (for $50+ monthly sponsors)

- Users donating $50+ per month can invite people to the `botstopper-customers` team in `TecharoHQ` org
- Organization-level sponsors at $50+ can also invite members

### 3. Logo Submission

- Form to upload company logo and link
- Image resizing and optimization
- Automatic GitHub issue creation in `TecharoHQ/anubis` with instructions for README update

---

## Database Schema

### Models (Gorm)

```go
// User represents an authenticated user from GitHub
type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    GitHubID  int64          `json:"github_id" gorm:"uniqueIndex;not null;index:idx_github_id"`
    Login     string         `json:"login" gorm:"not null;size:255;index:idx_login"`
    AvatarURL string         `json:"avatar_url" gorm:"size:500"`
    Name      string         `json:"name" gorm:"size:255"`
    Email     string         `json:"email" gorm:"size:255;index:idx_email"`
    Company   string         `json:"company" gorm:"size:255"`
    BlogURL   string         `json:"blog_url" gorm:"size:500"`
    Type      string         `json:"type" gorm:"not null;size:50;default:User"`

    // GitHub auth token (encrypted at application level)
    GitHubTokenEncrypted string `json:"-" gorm:"type:text;not null"`
    GitHubTokenScopes    string `json:"-" gorm:"size:500"`

    // Role-based access control
    IsAdmin   bool   `json:"is_admin" gorm:"default:false"`
    Role      string `json:"role" gorm:"size:50;default:user"`

    // Timestamps
    LastSeenAt    time.Time      `json:"last_seen_at" gorm:"index"`
    GitHubTokenAt *time.Time     `json:"-" gorm:"index"`
    CreatedAt     time.Time      `json:"created_at"`
    UpdatedAt     time.Time      `json:"updated_at"`
    DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

    // Relationships
    SponsorshipCaches       []SponsorshipCache `json:"-" gorm:"foreignKey:UserID"`
    TeamInvitationsSent     []TeamInvitation   `json:"-" gorm:"foreignKey:InviterID"`
    TeamInvitationsReceived []TeamInvitation   `json:"-" gorm:"foreignKey:InviteeID"`
    LogoSubmissions         []LogoSubmission    `json:"-" gorm:"foreignKey:SubmittedByUserID"`
}

// SponsorshipCache represents cached sponsorship data from GitHub API
type SponsorshipCache struct {
    ID     uint `json:"id" gorm:"primaryKey"`
    UserID uint `json:"user_id" gorm:"not null;index:idx_user_id;index:idx_user_sponsor,unique"`
    User   User `json:"user" gorm:"foreignKey:UserID"`

    // Sponsor identification
    SponsorID        uint   `json:"sponsor_id" gorm:"not null;index:idx_sponsor_id"`
    SponsorLogin     string `json:"sponsor_login" gorm:"not null;size:255"`
    SponsorName      string `json:"sponsor_name" gorm:"size:255"`
    SponsorAvatarURL string `json:"sponsor_avatar_url" gorm:"size:500"`
    SponsorType      string `json:"sponsor_type" gorm:"size:50;not null"`

    // Organization sponsor details
    SponsorOrganizationID *uint        `json:"sponsor_organization_id"`
    SponsorOrganization   *Organization `json:"sponsor_organization" gorm:"foreignKey:SponsorOrganizationID"`

    // Sponsorship details
    TierID        uint   `json:"tier_id" gorm:"not null;index:idx_tier_id;index:idx_user_tier,unique"`
    TierName      string `json:"tier_name" gorm:"not null;size:255"`
    AmountInCents int    `json:"amount_in_cents" gorm:"not null;default:0"`
    MonthlyAmount int    `json:"monthly_amount" gorm:"not null;default:0"`

    // Privacy and status
    PrivacyLevel string `json:"privacy_level" gorm:"not null;size:20;default:public"`
    IsActive     bool   `json:"is_active" gorm:"not null;default:true;index:idx_active"`

    // Cache management
    LastFetchedAt time.Time  `json:"last_fetched_at" gorm:"not null;index:idx_last_fetched"`
    ExpiresAt     time.Time  `json:"expires_at" gorm:"not null;index:idx_expires"`
    GitHubUpdatedAt *time.Time `json:"github_updated_at"`

    // Timestamps
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Organization represents GitHub organizations for org-level sponsorships
type Organization struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    GitHubID  int64          `json:"github_id" gorm:"uniqueIndex;not null;index:idx_org_github"`
    Login     string         `json:"login" gorm:"not null;uniqueIndex;size:255"`
    Name      string         `json:"name" gorm:"size:255"`
    AvatarURL string         `json:"avatar_url" gorm:"size:500"`
    Email     string         `json:"email" gorm:"size:255"`
    IsVerified bool          `json:"is_verified" gorm:"default:false"`

    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

    SponsorshipCaches []SponsorshipCache `json:"-" gorm:"foreignKey:SponsorOrganizationID"`
}

// TeamInvitation represents an invitation to join the GitHub team
type TeamInvitation struct {
    ID    uint `json:"id" gorm:"primaryKey"`

    // Who is being invited
    InviteeID    uint   `json:"invitee_id" gorm:"not null;index:idx_invitee;index:idx_invitee_status"`
    InviteeLogin string `json:"invitee_login" gorm:"not null;size:255;index:idx_invitee_login"`
    InviteeEmail string `json:"invitee_email" gorm:"size:255"`

    // Who sent the invitation
    InviterID    uint   `json:"inviter_id" gorm:"not null;index:idx_inviter"`
    InviterLogin string `json:"inviter_login" gorm:"not null;size:255"`

    // GitHub team details
    GitHubTeamID   int64  `json:"github_team_id" gorm:"not null;index:idx_github_team"`
    GitHubTeamName string `json:"github_team_name" gorm:"not null;size:255"`
    GitHubTeamSlug string `json:"github_team_slug" gorm:"not null;size:255;index:idx_team_slug"`

    // GitHub invitation details
    GitHubInvitationID  *int64 `json:"github_invitation_id" gorm:"uniqueIndex"`
    GitHubInvitationURL string `json:"github_invitation_url" gorm:"size:500"`

    // Status tracking
    Status        string `json:"status" gorm:"not null;size:50;default:pending;index:idx_status"`
    FailureReason string `json:"failure_reason" gorm:"size:500"`

    // Timestamps
    SentAt      time.Time  `json:"sent_at" gorm:"not null"`
    RespondedAt *time.Time `json:"responded_at" gorm:"index"`
    ExpiresAt   *time.Time `json:"expires_at" gorm:"index"`

    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// LogoSubmission represents a logo submitted by a sponsor
type LogoSubmission struct {
    ID    uint `json:"id" gorm:"primaryKey"`

    // Who submitted
    SubmittedByUserID uint `json:"submitted_by_user_id" gorm:"not null;index:idx_submitted_by"`

    // Company/organization details
    CompanyName    string `json:"company_name" gorm:"not null;size:255;index:idx_company"`
    CompanyWebsite string `json:"company_website" gorm:"not null;size:500"`

    // Logo details
    LogoURL      string `json:"logo_url" gorm:"not null;size:500"`
    LogoAltText  string `json:"logo_alt_text" gorm:"size:255"`
    LogoWidth    int    `json:"logo_width"`
    LogoHeight   int    `json:"logo_height"`
    LogoFormat   string `json:"logo_format" gorm:"size:50"`
    LogoFileSize int64  `json:"logo_file_size"`

    // GitHub issue tracking
    GitHubIssueURL    string `json:"github_issue_url" gorm:"size:500;index:idx_issue_url"`
    GitHubIssueNumber int    `json:"github_issue_number" gorm:"index"`
    GitHubIssueState  string `json:"github_issue_state" gorm:"size:50;default:open"`

    // Approval workflow
    Status          string     `json:"status" gorm:"not null;size:50;default:pending;index:idx_status"`
    ReviewedByUserID *uint     `json:"reviewed_by_user_id"`
    ReviewedAt      *time.Time `json:"reviewed_at"`
    RejectionReason string     `json:"rejection_reason" gorm:"size:1000"`
    AdminNotes      string     `json:"admin_notes" gorm:"size:1000"`

    // Timestamps
    PublishedAt *time.Time `json:"published_at" gorm:"index"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Session represents a gorilla/sessions session record
type Session struct {
    ID        string    `json:"id" gorm:"primaryKey;size:255"`
    UserID    uint      `json:"user_id" gorm:"not null;index:idx_session_user"`
    Data      string    `json:"-" gorm:"type:text"`
    ExpiresAt time.Time `json:"expires_at" gorm:"not null;index:idx_session_expires"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
```

---

## Configuration Patterns

### Environment Variables & Flags

```go
_ "github.com/joho/godotenv/autoload"

var (
    bind         = flag.String("bind", ":4823", "Port to listen on")
    databaseURL  = flag.String("database-url", "", "Database URL")
    githubToken  = flag.String("github-token", "", "GitHub token (GITHUB_TOKEN)")
    discordInvite = flag.String("discord-invite", "", "Discord invite link (DISCORD_INVITE)")

    // OAuth config
    githubClientID     = flag.String("github-client-id", "", "GitHub OAuth Client ID (GITHUB_CLIENT_ID)")
    githubClientSecret = flag.String("github-client-secret", "", "GitHub OAuth Client Secret (GITHUB_CLIENT_SECRET)")
    oauthRedirectURL   = flag.String("oauth-redirect-url", "", "OAuth redirect URL (OAUTH_REDIRECT_URL)")

    // Tigris/S3
    tigrisBucket = flag.String("tigris-bucket", "sponsor-panel-logos", "Tigris bucket name")
)

func main() {
    flagenv.Parse()
    flag.Parse()

    // Validate required flags
    if *databaseURL == "" {
        slog.Error("database-url is required")
        os.Exit(1)
    }
    // ... other validations
}
```

### Database Initialization

```go
db, err := gorm.Open(postgres.Open(*databaseURL), &gorm.Config{})
if err != nil {
    slog.Error("can't connect to database", "err", err)
    os.Exit(1)
}

// Test connection
if err := db.Exec("SELECT 1 + 1").Error; err != nil {
    slog.Error("can't ping database", "err", err)
    os.Exit(1)
}

// Auto-migrate models
if err := models.SetupDatabase(db); err != nil {
    slog.Error("database setup error", "err", err)
    os.Exit(1)
}
```

---

## GitHub OAuth Integration

### OAuth Flow

```text
User -> /login -> /login/oauth/authorize (GitHub) -> callback -> /callback -> exchange token -> fetch user info -> session
```

### Implementation

```go
import (
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/github"
)

var oauthConfig = &oauth2.Config{
    ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
    ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
    Scopes:       []string{"read:user", "user:email", "read:org", "read:sponsors"},
    Endpoint:     github.Endpoint,
    RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
}

// Login handler - redirect to GitHub
func loginHandler(w http.ResponseWriter, r *http.Request) {
    state := generateRandomState()
    session.Values["oauth_state"] = state
    session.Save(r, w)

    // With PKCE (recommended July 2025+)
    verifier := oauth2.GenerateVerifier()
    challenge := oauth2.S256ChallengeFromVerifier(verifier)
    session.Values["pkce_verifier"] = verifier
    session.Save(r, w)

    url := oauthConfig.AuthCodeURL(state,
        oauth2.AccessTypeOffline,
        oauth2.S256ChallengeOption(challenge),
    )

    http.Redirect(w, r, url, http.StatusFound)
}

// Callback handler - exchange code for token
func callbackHandler(w http.ResponseWriter, r *http.Request) {
    // Verify state
    if r.URL.Query().Get("state") != session.Values["oauth_state"] {
        http.Error(w, "state mismatch", http.StatusBadRequest)
        return
    }

    code := r.URL.Query().Get("code")

    // Exchange code for token with PKCE
    verifier := session.Values["pkce_verifier"].(string)
    token, err := oauthConfig.Exchange(ctx, code,
        oauth2.VerifierOption(verifier),
    )
    if err != nil {
        http.Error(w, "token exchange failed", http.StatusBadRequest)
        return
    }

    // Fetch user info
    client := oauthConfig.Client(ctx, token)
    resp, err := client.Get("https://api.github.com/user")
    if err != nil {
        http.Error(w, "fetch user failed", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    var githubUser GitHubUser
    json.NewDecoder(resp.Body).Decode(&githubUser)

    // Create/update user in database
    user := createOrUpdateUser(db, githubUser, token.AccessToken)

    // Set session
    session.Values["user_id"] = user.ID
    session.Save(r, w)

    http.Redirect(w, r, "/", http.StatusFound)
}
```

### Required Scopes

| Scope           | Purpose                                    |
| --------------- | ------------------------------------------ |
| `read:user`     | User profile (login, ID, avatar)           |
| `user:email`    | Email access                               |
| `read:org`      | Organization membership                    |
| `read:sponsors` | GraphQL Sponsors API (required since 2023) |

---

## Session Management

```go
import "github.com/gorilla/sessions"

var (
    store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
)

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        session, err := store.Get(r, "sponsor-panel")
        if err != nil || session.Values["user_id"] == nil {
            http.Redirect(w, r, "/login", http.StatusFound)
            return
        }

        // Load user from database
        var user User
        if err := db.First(&user, session.Values["user_id"]).Error; err != nil {
            http.Redirect(w, r, "/login", http.StatusFound)
            return
        }

        // Add user to context
        ctx := context.WithValue(r.Context(), "user", &user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

---

## GitHub Sponsors API

**Important:** GitHub Sponsors is **GraphQL-only**. The `google/go-github` library does not support sponsors.

### GraphQL Endpoint

```text
https://api.github.com/graphql
```

### Key Queries

#### Get Your Sponsors

```graphql
query {
  viewer {
    sponsorshipsAsMaintainer(first: 100, includePrivate: true) {
      nodes {
        sponsorEntity {
          ... on User {
            login
            name
            avatarUrl
          }
          ... on Organization {
            login
            name
            avatarUrl
          }
        }
        tier {
          id
          name
          monthlyPriceInCents
          isOneTime
        }
        privacyLevel
        createdAt
        isActive
      }
    }
  }
}
```

#### Check if User Sponsors You

```graphql
query ($login: String!) {
  user(login: $login) {
    ... on Sponsorable {
      sponsorshipForViewerAsSponsorable {
        isActive
        tier {
          name
          monthlyPriceInCents
        }
      }
    }
  }
}
```

### Go Implementation

```go
type GraphQLClient struct {
    client *http.Client
    token  string
}

func (c *GraphQLClient) Do(query string, variables map[string]interface{}, result interface{}) error {
    payload := map[string]interface{}{
        "query":     query,
        "variables": variables,
    }

    body, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewReader(body))
    req.Header.Set("Authorization", "Bearer "+c.token)
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return json.NewDecoder(resp.Body).Decode(result)
}
```

### Rate Limits

- **5,000 points per hour** for authenticated users
- Monitor headers: `x-ratelimit-remaining`, `x-ratelimit-reset`
- Cache with 15-30 minute TTL

---

## Templ + HTMX Integration

### HTMX Detection

```go
import "within.website/x/web/htmx"

func handler(w http.ResponseWriter, r *http.Request) {
    if htmx.Is(r) {
        // Return partial HTML for HTMX
        components.Partial().Render(r.Context(), w)
    } else {
        // Return full page
        components.FullPage().Render(r.Context(), w)
    }
}
```

### Common HTMX Patterns

```templ
<!-- Live search with debounce -->
<input hx-get="/search" hx-trigger="keyup changed delay:500ms" hx-target="#results" />

<!-- Inline edit -->
<button hx-get={"/edit/" + id} hx-target={"#field-" + id} hx-swap="outerHTML">Edit</button>

<!-- Delete with confirmation -->
<button hx-delete={"/items/" + itemID} hx-confirm="Are you sure?">Delete</button>

<!-- Out-of-band updates -->
<div id="cart-btn" hx-swap-oob="true">
    @CartButton(getCartCount())
</div>

<!-- Form with HTMX -->
<form hx-post="/submit" hx-target="#result" hx-swap="outerHTML">
    <input type="text" name="value" />
    <button type="submit">Submit</button>
</form>
```

---

## Image Upload & GitHub Issues

### Image Processing

```go
import "github.com/disintegration/imaging"

func processLogo(input io.Reader) (original, thumbnail, large []byte, err error) {
    src, err := imaging.Decode(input)
    if err != nil {
        return nil, nil, nil, err
    }

    // Generate sizes
    originalBuf := new(bytes.Buffer)
    imaging.Encode(originalBuf, src, imaging.PNG)

    thumbnail := imaging.Thumbnail(src, 300, 300, imaging.Lanczos)
    thumbBuf := new(bytes.Buffer)
    imaging.Encode(thumbBuf, thumbnail, imaging.PNG)

    large := imaging.Resize(src, 800, 0, imaging.Lanczos)
    largeBuf := new(bytes.Buffer)
    imaging.Encode(largeBuf, large, imaging.PNG)

    return originalBuf.Bytes(), thumbBuf.Bytes(), largeBuf.Bytes(), nil
}
```

### Tigris/S3 Storage

```go
import "github.com/aws/aws-sdk-go-v2/service/s3/presign"

// Generate presigned URL for upload
presignClient := presign.NewPresignClient(s3Client)
putReq := &s3.PutObjectInput{
    Bucket: aws.String("logos"),
    Key:    aws.String("submissions/" + filename),
}
presignedResult, err := presignClient.PresignPutObject(ctx, putReq,
    s3.WithPresignExpires(15*time.Minute),
)
```

### GitHub Issue Creation

```go
import "github.com/google/go-github/v60/github"

client := github.NewClient(tokenSource)

issueBody := fmt.Sprintf(`# Logo Submission: %s

**Submitted by:** @%s

## Logo Asset

![%s Logo](%s)

- **Original:** [Download](%s)
- **Thumbnail:** [Download](%s)
- **Format:** %s
- **Dimensions:** %dx%d

## Details

- **Website:** %s
- **Description:** %s

## Next Steps

1. Review logo meets requirements
2. Add to Anubis README
3. Update logos/index
4. Close this issue

/label needs-review
`, companyName, submitter, companyName, publicURL, originalURL, thumbURL, format, width, height, website, description)

issueReq := &github.IssueRequest{
    Title:  github.String("Logo Submission: " + companyName),
    Body:   github.String(issueBody),
    Labels: &[]string{"logo-submission", "needs-review"},
}

issue, _, err := client.Issues.Create(ctx, "TecharoHQ", "anubis", issueReq)
```

---

## Server Structure

```go
func main() {
    // ... setup code ...

    mux := http.NewServeMux()

    // OAuth routes
    mux.HandleFunc("/login", loginHandler)
    mux.HandleFunc("/callback", callbackHandler)
    mux.HandleFunc("/logout", logoutHandler)

    // Protected routes
    mux.Handle("/", authMiddleware(indexHandler))
    mux.Handle("/invite", authMiddleware(inviteHandler))
    mux.Handle("/logo/submit", authMiddleware(logoSubmitHandler))

    // HTMX endpoints
    mux.Handle("/api/check-sponsorship", authMiddleware(checkSponsorshipHandler))
    mux.Handle("/api/upload-logo-presign", authMiddleware(uploadLogoPresignHandler))

    // Static files
    mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServerFS(assets)))

    // Health check
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"ok"}`))
    })
    mux.Handle("/metrics", promhttp.Handler())

    // Middleware chain
    var h http.Handler = mux
    h = internal.CacheHeader(h)
    h = internal.AcceptEncodingMiddleware(h)
    h = internal.RefererMiddleware(h)

    slog.Info("sponsor-panel service ready", "bind", *bind)
    log.Fatal(http.Serve(ln, h))
}
```

---

## Project Structure

```
cmd/sponsor-panel/
├── main.go              # Entry point, server setup
├── handlers.go          # HTTP handlers
├── oauth.go             # OAuth flow handlers
├── github_client.go     # GraphQL client, go-github wrapper
├── image_processing.go  # Image upload and processing
└── templates/
    └── *.templ          # Templ components

internal/sponsorpanel/
├── models.go            # Gorm models (shared)
├── database.go          # Database setup
└── middleware.go        # Auth middleware
```

---

## Dependencies

```go
require (
    github.com/google/go-github/v82 v82.0.0
    github.com/gorilla/sessions v1.2.2
    github.com/joho/godotenv v1.5.1
    github.com/disintegration/imaging v1.6.2
    github.com/aws/aws-sdk-go-v2 v1.24.0
    github.com/aws/aws-sdk-go-v2/service/s3 v1.48.0
    github.com/aws/aws-sdk-go-v2/service/s3/presign v1.48.0
    golang.org/x/oauth2 v0.34.0
    gorm.io/gorm v1.25.5
    gorm.io/driver/postgres v1.5.4
    within.website/x/web/htmx v0.0.0
)
```

---

## Security Considerations

1. **PKCE for OAuth** - Required for security (July 2025+)
2. **State parameter** - CSRF protection for OAuth
3. **Encrypted tokens** - Store GitHub tokens encrypted at application layer
4. **HTTPS only** - All cookies and OAuth flows
5. **Session expiration** - Reasonable TTL with refresh capability
6. **Rate limiting** - Respect GitHub API limits
7. **Input validation** - Sanitize all user inputs

---

## GitHub Team Management API

**Dependency:** `github.com/google/go-github/v82` (add with `go get`)

### API Methods for Team Operations

#### Check Team Membership

```go
import "github.com/google/go-github/v82/github"

func (s *TeamsService) GetTeamMembershipBySlug(
    ctx context.Context,
    org, slug, user string,
) (*Membership, *Response, error)
```

- **org**: "TecharoHQ"
- **slug**: "botstopper-customers"
- **user**: GitHub username to check
- **Returns**: `*Membership` with State ("active" or "pending") and Role ("member" or "maintainer")

#### Add User to Team

```go
func (s *TeamsService) AddTeamMembershipBySlug(
    ctx context.Context,
    org, slug, user string,
    opts *TeamAddTeamMembershipOptions,
) (*Membership, *Response, error)
```

**Options:**

- `Role`: "member" (default) or "maintainer"

**Important:** User must already be an organization member before adding to team.

#### Remove User from Team

```go
func (s *TeamsService) RemoveTeamMembershipBySlug(
    ctx context.Context,
    org, slug, user string,
) (*Response, error)
```

### Organization Membership

Users must be organization members before being added to teams.

#### Check if User is in Organization

```go
func (s *OrganizationsService) IsMember(
    ctx context.Context,
    org, user string,
) (bool, *Response, error)
```

#### Invite User to Organization

```go
func (s *OrganizationsService) CreateOrgInvitation(
    ctx context.Context,
    org string,
    opts *CreateOrgInvitationOptions,
) (*Invitation, *Response, error)
```

**Options:**

- `InviteeID`: GitHub user ID
- `Email`: Email address (alternative to InviteeID)
- `Role`: "admin", "direct_member" (default), or "billing_manager"
- `TeamID`: []int64 - list of team IDs to auto-add upon acceptance

### Required GitHub App Permissions

| Permission                       | Scope        | Purpose                                       |
| -------------------------------- | ------------ | --------------------------------------------- |
| **Members: Read & Write**        | Organization | Check and manage org membership, invite users |
| **Administration: Read & Write** | Organization | Team management operations                    |

### Error Handling

| Status  | Meaning              | Action                                                               |
| ------- | -------------------- | -------------------------------------------------------------------- |
| **200** | Success              | Operation completed                                                  |
| **202** | Accepted             | Invitation sent, pending acceptance                                  |
| **403** | Forbidden            | Insufficient permissions                                             |
| **404** | Not Found            | Resource doesn't exist OR insufficient permissions (GitHub security) |
| **422** | Unprocessable Entity | Validation error (e.g., user already in team)                        |

### Recommended Workflow

```text
1. Check if user is in organization
   └─ NO: Invite to organization with CreateOrgInvitation() (include TeamID)
   └─ YES: Continue

2. Check team membership
   └─ Active: User already has access, done
   └─ Not found: Add to team with AddTeamMembershipBySlug()

3. Log invitation to database for audit trail
```

### Implementation Example

```go
func inviteToTeam(ctx context.Context, client *github.Client, username, org, teamSlug string) error {
    // 1. Check org membership
    isMember, _, err := client.Organizations.IsMember(ctx, org, username)
    if err != nil {
        return fmt.Errorf("check org membership: %w", err)
    }

    if !isMember {
        // Invite to org (with team auto-add)
        invite := &github.CreateOrgInvitationOptions{
            InviteeID: github.Int64(getUserID(username)),
            Role:      github.String("direct_member"),
        }
        _, _, err = client.Organizations.CreateOrgInvitation(ctx, org, invite)
        if err != nil {
            return fmt.Errorf("invite to org: %w", err)
        }
        return nil // Invitation pending
    }

    // 2. Check team membership
    membership, _, err := client.Teams.GetTeamMembershipBySlug(ctx, org, teamSlug, username)
    if err != nil {
        // 404 means not in team, add them
        if strings.Contains(err.Error(), "404") {
            membership, _, err = client.Teams.AddTeamMembershipBySlug(ctx, org, teamSlug, username, nil)
            if err != nil {
                return fmt.Errorf("add to team: %w", err)
            }
        } else {
            return fmt.Errorf("check team membership: %w", err)
        }
    }

    if membership != nil && membership.State == "active" {
        return nil // Already a member
    }

    return nil
}
```

---

## Sources

- GitHub OAuth Docs: https://docs.github.com/en/apps/oauth-apps
- GitHub Sponsors GraphQL: https://docs.github.com/en/sponsors
- Go OAuth2: https://golang.org/x/oauth2
- go-github: https://github.com/google/go-github
- Templ: https://templ.guide
- HTMX: https://htmx.org
- Gorm: https://gorm.io/docs/
