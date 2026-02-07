# Sponsor Panel Service - Simplified Specification

**Service:** sponsor-panel
**Exposure:** sponsors.xeiaso.net
**Design Philosophy:** Minimal, working service without over-engineering

---

## Overview

This is a **simplified version** of the sponsor-panel specification. The original spec (see `SPEC.md` and `UX_FLOWS.md`) contained significant over-engineering. This version removes unnecessary complexity while delivering the same core value.

**Core principle:** Build the simplest thing that works. Add complexity only when proven necessary.

---

## Features (Simplified)

| Feature             | $0  | $1-49 | $50+ | Implementation                             |
| ------------------- | :-: | :---: | :--: | ------------------------------------------ |
| Discord invite link | ✅  |  ✅   |  ✅  | Environment variable, display on dashboard |
| Sponsorship status  |  -  |  ✅   |  ✅  | GraphQL on login, store in session         |
| Team invitation     |  -  |   -   |  ✅  | Direct GitHub API call, no audit trail     |
| Logo submission     |  -  |  ✅   |  ✅  | Create GitHub issue with attachment        |

**Removed from original spec:**

- ❌ Organization-specific team invitations
- ❌ Invitation history/audit trail
- ❌ Logo approval workflow
- ❌ Multiple image sizes
- ❌ Image resizing/optimization
- ❌ Real-time username validation
- ❌ Activity feed
- ❌ Admin panel

---

## Tech Stack (Simplified)

| Layer        | Technology                  | Notes                 |
| ------------ | --------------------------- | --------------------- |
| **Backend**  | Go with `net/http`          | Standard library only |
| **Frontend** | Templ + HTMX + Tailwind CSS | Unchanged             |
| **Database** | PostgreSQL with sqlx        | No ORM                |
| **Storage**  | GitHub issue attachments    | No external storage   |
| **Auth**     | GitHub OAuth 2.0            | No PKCE (server-side) |
| **Sessions** | Cookie-only                 | No database backing   |

### Removed Dependencies

```
- github.com/gorilla/sessions  → Use simple cookie store
- github.com/disintegration/imaging → No image processing
- github.com/aws/aws-sdk-go-v2 → Use GitHub attachments
- gorm.io/gorm → Use sqlx
- gorm.io/driver/postgres → Use lib/pq directly
```

---

## Data Model (Simplified)

### 2 Tables Instead of 6

```sql
-- Users table: GitHub accounts + sponsorship data
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    github_id BIGINT UNIQUE NOT NULL,
    login TEXT NOT NULL UNIQUE,
    avatar_url TEXT,
    name TEXT,
    email TEXT,

    -- Sponsorship data from GraphQL (cached)
    sponsorship_data JSONB,
    last_sponsorship_check TIMESTAMP DEFAULT NOW(),

    -- Session data (encrypted cookie, no DB table)
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Logo submissions: Simple tracking only
CREATE TABLE logo_submissions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,

    company_name TEXT NOT NULL,
    website TEXT NOT NULL,
    logo_url TEXT,              -- GitHub issue attachment URL
    github_issue_url TEXT,
    github_issue_number INTEGER,

    submitted_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for common queries
CREATE INDEX idx_users_github_id ON users(github_id);
CREATE INDEX idx_users_login ON users(login);
CREATE INDEX idx_logo_user_id ON logo_submissions(user_id);
```

### Go Models

```go
package models

import "time"

// User represents a GitHub user with cached sponsorship data
type User struct {
    ID                      int       `json:"id" db:"id"`
    GitHubID                int64     `json:"github_id" db:"github_id"`
    Login                   string    `json:"login" db:"login"`
    AvatarURL               string    `json:"avatar_url" db:"avatar_url"`
    Name                    string    `json:"name" db:"name"`
    Email                   string    `json:"email" db:"email"`
    SponsorshipData         string    `json:"-" db:"sponsorship_data"` // JSON blob
    LastSponsorshipCheck    time.Time `json:"last_sponsorship_check" db:"last_sponsorship_check"`
    CreatedAt               time.Time `json:"created_at" db:"created_at"`
    UpdatedAt               time.Time `json:"updated_at" db:"updated_at"`
}

// SponsorshipData represents the cached GraphQL response
type SponsorshipData struct {
    IsActive      bool    `json:"is_active"`
    MonthlyAmount int     `json:"monthly_amount_cents"`
    TierName      string  `json:"tier_name"`
    PrivacyLevel  string  `json:"privacy_level"`
}

// IsSponsorAtTier returns true if user sponsors at or above the given amount (in cents)
func (u *User) IsSponsorAtTier(minCents int) bool {
    if u.SponsorshipData == "" {
        return false
    }

    var data SponsorshipData
    if err := json.Unmarshal([]byte(u.SponsorshipData), &data); err != nil {
        return false
    }

    return data.IsActive && data.MonthlyAmount >= minCents
}

// LogoSubmission represents a logo submission
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
```

---

## Architecture (Simplified)

```
┌─────────────────────────────────────────────────────────────────────┐
│                     USER BROWSER (HTMX)                             │
└─────────────────────────────┬───────────────────────────────────────┘
                              │
                              ▼
┌───────────────────────────────────────────────────────────────────┐
│                       sponsor-panel Service                       │
├───────────────────────────────────────────────────────────────────┤
│  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐   │
│  │   OAuth    │  │   Sponsor  │  │   Team     │  │    Logo    │   │
│  │  Handler   │  │  Checker   │  │  Inviter   │  │  Submitter │   │
│  └──────┬─────┘  └──────┬─────┘  └──────┬─────┘  └──────┬─────┘   │
└─────────┼───────────────┼───────────────┼───────────────┼─────────┘
          │               │               │               │
          ▼               ▼               ▼               ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     sqlx (PostgreSQL)                               │
│  • users (2 tables)                                                 │
│  • logo_submissions                                                 │
└─────────────────────────────────────────────────────────────────────┘
          │                │                │
          ▼                ▼                ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     GitHub APIs                                     │
│  • OAuth 2.0                                                        │
│  • GraphQL (Sponsors)                                               │
│  • REST (Teams, Issues)                                             │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Endpoints (Simplified)

### 5 Endpoints Total

| Method | Endpoint    | Auth | Purpose                        |
| ------ | ----------- | ---- | ------------------------------ |
| GET    | `/login`    | No   | Initiate GitHub OAuth          |
| GET    | `/callback` | No   | OAuth callback, create session |
| GET    | `/`         | Yes  | Dashboard                      |
| POST   | `/invite`   | Yes  | Invite user to team            |
| POST   | `/logo`     | Yes  | Submit logo (create issue)     |

### Endpoint Details

#### `GET /login`

Initiate GitHub OAuth flow.

```go
func loginHandler(w http.ResponseWriter, r *http.Request) {
    // Generate random state
    state := generateState()

    // Set state in cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "oauth_state",
        Value:    state,
        Path:     "/",
        Secure:   true,
        HttpOnly: true,
        SameSite: http.SameSiteStrictMode,
    })

    // Redirect to GitHub
    url := fmt.Sprintf(
        "https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=read:user%%2Cuser:email%%2Cread:org%%2Cread:sponsors&state=%s",
        *clientID, *redirectURL, state,
    )
    http.Redirect(w, r, url, http.StatusFound)
}
```

#### `GET /callback`

Handle OAuth callback, fetch user + sponsorship data, create session.

```go
func callbackHandler(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Verify state
        stateCookie, _ := r.Cookie("oauth_state")
        if r.URL.Query().Get("state") != stateCookie.Value {
            http.Error(w, "Invalid state", http.StatusBadRequest)
            return
        }

        // Exchange code for token
        token, err := exchangeCode(r.URL.Query().Get("code"))
        if err != nil {
            http.Error(w, "Token exchange failed", http.StatusBadRequest)
            return
        }

        // Fetch user from GitHub
        ghUser, err := fetchGitHubUser(token)
        if err != nil {
            http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
            return
        }

        // Fetch sponsorship data from GraphQL
        sponsorData, err := fetchSponsorship(token)
        if err != nil {
            // Non-fatal: log but continue
            slog.Error("failed to fetch sponsorship", "err", err)
            sponsorData = "{}"
        }

        // Upsert user in database
        var user User
        err = db.Get(&user, "SELECT * FROM users WHERE github_id = $1", ghUser.ID)
        if err == sql.ErrNoRows {
            err = db.QueryRow(
                `INSERT INTO users (github_id, login, avatar_url, name, email, sponsorship_data)
                 VALUES ($1, $2, $3, $4, $5, $6)
                 RETURNING *`,
                ghUser.ID, ghUser.Login, ghUser.AvatarURL, ghUser.Name, ghUser.Email, sponsorData,
            ).Scan(&user.ID, &user.GitHubID, &user.Login, &user.AvatarURL,
                &user.Name, &user.Email, &user.SponsorshipData, &user.LastSponsorshipCheck,
                &user.CreatedAt, &user.UpdatedAt)
        } else {
            db.QueryRow(
                `UPDATE users SET login=$1, avatar_url=$2, name=$3, email=$4, sponsorship_data=$5,
                 last_sponsorship_check=NOW(), updated_at=NOW() WHERE github_id=$6 RETURNING *`,
                ghUser.Login, ghUser.AvatarURL, ghUser.Name, ghUser.Email, sponsorData, ghUser.ID,
            ).Scan(&user.ID, &user.GitHubID, &user.Login, &user.AvatarURL,
                &user.Name, &user.Email, &user.SponsorshipData, &user.LastSponsorshipCheck,
                &user.CreatedAt, &user.UpdatedAt)
        }

        // Create encrypted session cookie
        sessionData := fmt.Sprintf("%d|%s", user.ID, user.Login)
        encrypted := encrypt(sessionData, *sessionKey)

        http.SetCookie(w, &http.Cookie{
            Name:     "session",
            Value:    encrypted,
            Path:     "/",
            Secure:   true,
            HttpOnly: true,
            SameSite: http.SameSiteStrictMode,
            MaxAge:   30 * 24 * 3600, // 30 days
        })

        http.Redirect(w, r, "/", http.StatusFound)
    }
}
```

#### `GET /` (Dashboard)

Render dashboard based on sponsorship tier.

```go
func dashboardHandler(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Get user from session
        session, _ := r.Cookie("session")
        if session == nil {
            http.Redirect(w, r, "/login", http.StatusFound)
            return
        }

        decrypted := decrypt(session.Value, *sessionKey)
        parts := strings.Split(decrypted, "|")
        userID := parts[0]

        // Fetch user
        var user User
        if err := db.Get(&user, "SELECT * FROM users WHERE id = $1", userID); err != nil {
            http.Redirect(w, r, "/login", http.StatusFound)
            return
        }

        // Render dashboard
        isFiftyPlus := user.IsSponsorAtTier(5000) // $50 = 5000 cents
        isSponsor := user.IsSponsorAtTier(100)   // $1 = 100 cents

        components.Dashboard(components.DashboardProps{
            User:         user,
            IsSponsor:    isSponsor,
            IsFiftyPlus:  isFiftyPlus,
            DiscordInvite: *discordInvite,
        }).Render(r.Context(), w)
    }
}
```

#### `POST /invite` (Team Invitation)

Invite a GitHub user to `botstopper-customers` team.

```go
func inviteHandler(db *sqlx.DB, ghClient *github.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Get user from session (verify $50+ sponsor)
        user := getUserFromSession(r, db)
        if !user.IsSponsorAtTier(5000) {
            http.Error(w, "Requires $50+/month sponsorship", http.StatusForbidden)
            return
        }

        // Parse form
        username := r.FormValue("username")
        if username == "" {
            http.Error(w, "Username required", http.StatusBadRequest)
            return
        }

        // Strip @ if present
        username = strings.TrimPrefix(username, "@")

        // Invite to team (direct call, no audit trail)
        membership, _, err := ghClient.Teams.AddTeamMembershipBySlug(
            r.Context(), "TecharoHQ", "botstopper-customers", username, nil,
        )

        if err != nil {
            // Check if already member or invited
            if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "422") {
                components.InviteError("User already invited or not found").Render(r.Context(), w)
                return
            }
            http.Error(w, "Failed to invite: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // Show success (no database storage)
        state := "pending"
        if membership != nil && membership.State == "active" {
            state = "active"
        }

        components.InviteSuccess(username, state).Render(r.Context(), w)
    }
}
```

#### `POST /logo` (Logo Submission)

Create GitHub issue with logo attachment.

```go
func logoHandler(db *sqlx.DB, ghClient *github.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Get user from session
        user := getUserFromSession(r, db)

        // Parse form
        companyName := r.FormValue("company_name")
        website := r.FormValue("website")

        // Upload file to GitHub issue as attachment
        file, header, err := r.FormFile("logo")
        if err != nil {
            http.Error(w, "Logo file required", http.StatusBadRequest)
            return
        }
        defer file.Close()

        // Validate size (max 5MB)
        if header.Size > 5*1024*1024 {
            http.Error(w, "File too large (max 5MB)", http.StatusBadRequest)
            return
        }

        // Create issue first
        issueBody := fmt.Sprintf(`# Logo Submission: %s

**Submitted by:** @%s

## Details

- **Company:** %s
- **Website:** %s
- **File:** %s

## Next Steps

1. Review logo
2. Add to Anubis README
3. Close this issue

/label logo-submission
/label needs-review
`, companyName, user.Login, companyName, website, header.Filename)

        issue, _, err := ghClient.Issues.Create(r.Context(), "TecharoHQ", "anubis", &github.IssueRequest{
            Title:  github.String("Logo Submission: " + companyName),
            Body:   github.String(issueBody),
            Labels: &[]string{"logo-submission", "needs-review"},
        })

        if err != nil {
            http.Error(w, "Failed to create issue: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // TODO: Upload file as issue attachment
        // Note: GitHub's REST API doesn't support direct file uploads to issues
        // Alternative: Upload to repo assets, reference in issue
        // For now: create issue with instructions to email logo

        // Store submission record
        var submissionID int
        err = db.QueryRow(
            `INSERT INTO logo_submissions (user_id, company_name, website, github_issue_url, github_issue_number)
             VALUES ($1, $2, $3, $4, $5) RETURNING id`,
            user.ID, companyName, website, issue.GetHTMLURL(), issue.GetNumber(),
        ).Scan(&submissionID)

        if err != nil {
            slog.Error("failed to store submission", "err", err)
        }

        // Show success
        components.LogoSuccess(companyName, issue.GetHTMLURL(), issue.GetNumber()).Render(r.Context(), w)
    }
}
```

---

## Configuration (Simplified)

### Environment Variables

| Variable               | Purpose                             | Required |
| ---------------------- | ----------------------------------- | -------- |
| `DATABASE_URL`         | PostgreSQL connection               | Yes      |
| `GITHUB_CLIENT_ID`     | OAuth app ID                        | Yes      |
| `GITHUB_CLIENT_SECRET` | OAuth app secret                    | Yes      |
| `GITHUB_TOKEN`         | GitHub App token for team/issue ops | Yes      |
| `SESSION_KEY`          | Session encryption key              | Yes      |
| `DISCORD_INVITE`       | Discord invite link                 | Yes      |
| `OAUTH_REDIRECT_URL`   | OAuth callback URL                  | Yes      |

### Flags

```go
var (
    bind         = flag.String("bind", ":4823", "Port to listen on")
    databaseURL  = flag.String("database-url", "", "Database URL")
    clientID     = flag.String("github-client-id", "", "GitHub OAuth Client ID")
    clientSecret = flag.String("github-client-secret", "", "GitHub OAuth Client Secret")
    githubToken  = flag.String("github-token", "", "GitHub token for operations")
    sessionKey   = flag.String("session-key", "", "Session encryption key")
    discordInvite = flag.String("discord-invite", "", "Discord invite link")
    redirectURL  = flag.String("oauth-redirect-url", "", "OAuth redirect URL")
)
```

---

## GraphQL Queries (Simplified)

### Single Query on Login

```graphql
query CheckSponsorship {
  viewer {
    sponsorshipsAsMaintainer(first: 100, includePrivate: true) {
      nodes {
        sponsorEntity {
          ... on User {
            login
          }
          ... on Organization {
            login
          }
        }
        tier {
          monthlyPriceInCents
          name
        }
        privacyLevel
        isActive
      }
    }
  }
}
```

### Go Implementation

```go
func fetchSponsorship(token string) (string, error) {
    query := `{"query": "query { viewer { sponsorshipsAsMaintainer(first: 100, includePrivate: true) { nodes { sponsorEntity { ... on User { login } ... on Organization { login } } tier { monthlyPriceInCents name } privacyLevel isActive } } } }"}`

    req, _ := http.NewRequest("POST", "https://api.github.com/graphql", strings.NewReader(query))
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var result struct {
        Data struct {
            Viewer struct {
                Sponsorships struct {
                    Nodes []struct {
                        SponsorEntity struct {
                            Login string `json:"login"`
                        } `json:"sponsorEntity"`
                        Tier struct {
                            MonthlyPriceInCents int `json:"monthlyPriceInCents"`
                            Name                string `json:"name"`
                        } `json:"tier"`
                        IsActive     bool   `json:"isActive"`
                        PrivacyLevel string `json:"privacyLevel"`
                    } `json:"nodes"`
                } `json:"sponsorshipsAsMaintainer"`
            } `json:"viewer"`
        } `json:"data"`
    }

    json.NewDecoder(resp.Body).Decode(&result)

    // Return highest tier sponsorship as JSON
    if len(result.Data.Viewer.Sponsorships.Nodes) == 0 {
        return `{"is_active": false}`, nil
    }

    // Find highest active tier
    var highest *struct {
        MonthlyPriceInCents int
        Name                string
    }

    for _, node := range result.Data.Viewer.Sponsorships.Nodes {
        if node.IsActive && (highest == nil || node.Tier.MonthlyPriceInCents > highest.MonthlyPriceInCents) {
            highest = &node.Tier
        }
    }

    if highest == nil {
        return `{"is_active": false}`, nil
    }

    resultJSON, _ := json.Marshal(map[string]interface{}{
        "is_active":            true,
        "monthly_amount_cents": highest.MonthlyPriceInCents,
        "tier_name":            highest.Name,
    })

    return string(resultJSON), nil
}
```

---

## Session Management (Simplified)

### Cookie-Only Sessions

```go
// Simple encryption using crypto/aes
func encrypt(plaintext, key string) string {
    block, _ := aes.NewCipher([]byte(key))
    gcm, _ := cipher.NewGCM(block)
    nonce := make([]byte, gcm.NonceSize())
    io.ReadFull(rand.Reader, nonce)
    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    return base64.StdEncoding.EncodeToString(ciphertext)
}

func decrypt(ciphertext, key string) string {
    data, _ := base64.StdEncoding.DecodeString(ciphertext)
    block, _ := aes.NewCipher([]byte(key))
    gcm, _ := cipher.NewGCM(block)
    nonceSize := gcm.NonceSize()
    nonce, cipherData := data[:nonceSize], data[nonceSize:]
    plaintext, _ := gcm.Open(nil, nonce, cipherData, nil)
    return string(plaintext)
}

func getUserFromSession(r *http.Request, db *sqlx.DB) User {
    session, _ := r.Cookie("session")
    if session == nil {
        return User{}
    }

    decrypted := decrypt(session.Value, *sessionKey)
    parts := strings.Split(decrypted, "|")
    if len(parts) != 2 {
        return User{}
    }

    userID := parts[0]
    var user User
    if err := db.Get(&user, "SELECT * FROM users WHERE id = $1", userID); err != nil {
        return User{}
    }

    return user
}
```

---

## Templ Components (Simplified)

### Dashboard Structure

```templ
templ Dashboard(props DashboardProps) {
    @base("Sponsor Panel") {
        @Navbar(props.User.Login)
        <main class="max-w-4xl mx-auto px-4 py-8">
            <h1 class="text-2xl font-bold mb-6">Welcome, @props.User.Login!</h1>

            <div class="grid md:grid-cols-2 gap-6">
                @DiscordCard(props.DiscordInvite)
                @SponsorshipCard(props.User, props.IsSponsor, props.IsFiftyPlus)
                @TeamInviteCard(props.IsFiftyPlus)
                @LogoSubmitCard(props.IsSponsor)
            </div>
        </main>
    }
}

templ DiscordCard(inviteURL string) {
    <div class="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
        <h2 class="font-semibold mb-2">Discord Community</h2>
        <p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
            Join our Discord server to chat with other sponsors.
        </p>
        <a href={ inviteURL } target="_blank" class="inline-block bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-lg">
            Join Discord
        </a>
    </div>
}

templ SponsorshipCard(user User, isSponsor, isFiftyPlus bool) {
    <div class="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
        <h2 class="font-semibold mb-2">Your Sponsorship</h2>
        if !isSponsor {
            <p class="text-sm text-gray-600 dark:text-gray-400">
                No active sponsorship found.
            </p>
            <a href="https://github.com/sponsors/Xe" target="_blank" class="inline-block bg-pink-600 hover:bg-pink-700 text-white px-4 py-2 rounded-lg mt-2">
                Become a Sponsor
            </a>
        } else if isFiftyPlus {
            <p class="text-green-600 dark:text-green-400 font-semibold">
                💎 Premium Sponsor ($50+/month)
            </p>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
                Full access to all benefits!
            </p>
        } else {
            <p class="text-blue-600 dark:text-blue-400 font-semibold">
                ❤️ Supporter ($1-49/month)
            </p>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
                Thank you for your support!
            </p>
        }
    </div>
}

templ TeamInviteCard(isFiftyPlus bool) {
    <div class={ "bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm", !isFiftyPlus: "opacity-50" }>
        <h2 class="font-semibold mb-2">Team Invitation</h2>
        if !isFiftyPlus {
            <p class="text-sm text-gray-600 dark:text-gray-400">
                Requires $50+/month sponsorship.
            </p>
        } else {
            <form hx-post="/invite" hx-target="#invite-result" class="space-y-3">
                <input type="text" name="username" placeholder="GitHub username" required
                    class="w-full border rounded-lg px-3 py-2 dark:bg-gray-700 dark:border-gray-600" />
                <button type="submit" class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg">
                    Invite to Team
                </button>
            </form>
            <div id="invite-result"></div>
        }
    </div>
}

templ LogoSubmitCard(isSponsor bool) {
    <div class="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
        <h2 class="font-semibold mb-2">Logo Submission</h2>
        <p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
            Submit your company logo for the Anubis README.
        </p>
        if !isSponsor {
            <a href="https://github.com/sponsors/Xe" target="_blank" class="inline-block bg-orange-600 hover:bg-orange-700 text-white px-4 py-2 rounded-lg">
                Become a Sponsor
            </a>
        } else {
            <a href="/logo" class="inline-block bg-orange-600 hover:bg-orange-700 text-white px-4 py-2 rounded-lg">
                Submit Logo
            </a>
        }
    </div>
}
```

### Simple Logo Form

```templ
templ LogoSubmitPage() {
    @base("Submit Logo") {
        @Navbar("")
        <main class="max-w-2xl mx-auto px-4 py-8">
            <h1 class="text-2xl font-bold mb-6">Submit Your Logo</h1>

            <form hx-post="/logo" hx-encoding="multipart/form-data" hx-target="#result" class="space-y-4 bg-white dark:bg-gray-800 rounded-xl p-6">
                <div>
                    <label class="block text-sm font-medium mb-1">Company Name</label>
                    <input type="text" name="company_name" required class="w-full border rounded-lg px-3 py-2 dark:bg-gray-700 dark:border-gray-600" />
                </div>

                <div>
                    <label class="block text-sm font-medium mb-1">Website</label>
                    <input type="url" name="website" required class="w-full border rounded-lg px-3 py-2 dark:bg-gray-700 dark:border-gray-600" />
                </div>

                <div>
                    <label class="block text-sm font-medium mb-1">Logo (PNG, SVG, or JPG, max 5MB)</label>
                    <input type="file" name="logo" accept="image/png,image/svg+xml,image/jpeg" required class="w-full border rounded-lg px-3 py-2 dark:bg-gray-700 dark:border-gray-600" />
                </div>

                <button type="submit" class="bg-orange-600 hover:bg-orange-700 text-white px-6 py-2 rounded-lg">
                    Submit Logo
                </button>
            </form>

            <div id="result"></div>
        </main>
    }
}
```

---

## Error Handling (Simplified)

### Simple Error Responses

```go
// Return HTML for HTMX or redirect for full page requests
func handleError(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
    if htmx.Is(r) {
        w.Header().Set("Content-Type", "text/html")
        w.WriteHeader(statusCode)
        components.ErrorMessage(message).Render(r.Context(), w)
    } else {
        w.Header().Set("Content-Type", "text/html")
        w.WriteHeader(statusCode)
        components.ErrorPage(message).Render(r.Context(), w)
    }
}

templ ErrorMessage(message string) {
    <div class="bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-300 p-4 rounded-lg">
        { message }
    </div>
}
```

---

## Security Considerations

| Concern           | Mitigation                          |
| ----------------- | ----------------------------------- |
| CSRF on OAuth     | State parameter in cookie           |
| Session hijacking | Secure, HttpOnly cookies            |
| Token storage     | Encrypted at rest, in memory only   |
| SQL injection     | Use sqlx with parameterized queries |
| File uploads      | Size limit (5MB), type validation   |

---

## Deployment

### Environment

```yaml
# Kubernetes deployment example
apiVersion: v1
kind: Secret
metadata:
  name: sponsor-panel-secrets
stringData:
  DATABASE_URL: "postgres://..."
  GITHUB_CLIENT_ID: "..."
  GITHUB_CLIENT_SECRET: "..."
  GITHUB_TOKEN: "..."
  SESSION_KEY: "32-byte-random-key"
  DISCORD_INVITE: "https://discord.gg/..."
  OAUTH_REDIRECT_URL: "https://sponsors.xeiaso.net/callback"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sponsor-panel
spec:
  replicas: 1
  template:
    spec:
      containers:
        - name: sponsor-panel
          image: ghcr.io/xeiaso/sponsor-panel:latest
          ports:
            - containerPort: 4823
          envFrom:
            - secretRef:
                name: sponsor-panel-secrets
```

---

## Implementation Checklist (Simplified)

### Phase 1: Foundation

- [ ] Create `internal/models/` with User and LogoSubmission structs
- [ ] Set up sqlx with PostgreSQL
- [ ] Create migration for 2 tables
- [ ] Implement cookie encryption/decryption

### Phase 2: Authentication

- [ ] Implement `/login` OAuth redirect
- [ ] Implement `/callback` with user creation
- [ ] Add session cookie handling
- [ ] Create logout handler

### Phase 3: Sponsorship

- [ ] Implement GraphQL client
- [ ] Cache sponsorship data in user table
- [ ] Add `IsSponsorAtTier()` helper method

### Phase 4: Dashboard & Features

- [ ] Build Dashboard.templ with 4 cards
- [ ] Implement `/invite` with GitHub API
- [ ] Implement `/logo` with GitHub issue creation
- [ ] Add error handling components

### Phase 5: Testing & Deploy

- [ ] Unit tests for handlers
- [ ] Integration tests with test database
- [ ] Deploy to Kubernetes

---

## Migration from Original Spec

### What Changed

| Aspect            | Original                    | Simplified          | Change |
| ----------------- | --------------------------- | ------------------- | ------ |
| Tables            | 6                           | 2                   | -67%   |
| Dependencies      | 15+                         | 6                   | -60%   |
| Endpoints         | 10+                         | 5                   | -50%   |
| External services | 3 (GitHub, Tigris, Discord) | 2 (GitHub, Discord) | -33%   |
| Code estimate     | ~4000 lines                 | ~1500 lines         | -62%   |

### Features Removed

- Organization-specific teams
- Invitation audit trail
- Logo approval workflow
- Image processing/resizing
- Real-time validation
- Activity feeds
- Admin panel
- Complex session management

### What's Still There

- ✅ GitHub OAuth login
- ✅ Sponsorship tier checking
- ✅ Discord invite link
- ✅ Team invitations ($50+)
- ✅ Logo submission (all sponsors)

---

## Summary

This simplified specification delivers the same core value with **60-70% less complexity**:

- **2 tables** instead of 6
- **5 endpoints** instead of 10+
- **No ORM** (use sqlx)
- **No external storage** (use GitHub attachments)
- **No image processing** (upload as-is)
- **No audit trails** (fire and forget)
- **Cookie-only sessions** (no database backing)

Build the simplest thing that works. Add complexity only when you have proven need.

---

**Simplified specification created:** 2026-02-07
**Based on original spec:** `SPEC.md`, `UX_FLOWS.md`
**Ruthless review by:** Cynical Engineer Agent
