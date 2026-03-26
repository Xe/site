# Sponsor Panel GORM Migration Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace raw SQL (`pgxpool` + hand-written queries) in `cmd/sponsor-panel` with GORM, matching the pattern already used by `cmd/github-sponsor-webhook` and `internal/models`.

**Architecture:** Define three new GORM models (`PanelUser`, `LogoSubmission`, `SponsorUsername`) in `cmd/sponsor-panel/models.go`. Replace every `pgxpool` call site with GORM equivalents. Swap the `*pgxpool.Pool` on the `Server` struct for `*gorm.DB`. Delete `migrations.go` entirely since GORM AutoMigrate handles schema. The existing `internal/models` package is for the webhook service's domain (accounts, tiers, sponsorships, webhook events) and is unrelated -- we do NOT reuse those models here.

**Tech Stack:** Go, GORM (`gorm.io/gorm` + `gorm.io/driver/postgres`) -- both already in `go.mod`.

---

## File Structure

| Action | File | Responsibility |
|--------|------|----------------|
| Rewrite | `cmd/sponsor-panel/models.go` | GORM model structs + all DB helper functions |
| Delete | `cmd/sponsor-panel/migrations.go` | Replaced by GORM AutoMigrate |
| Modify | `cmd/sponsor-panel/main.go` | Swap `pgxpool.Pool` for `gorm.DB`, call AutoMigrate, update `Server` |
| Modify | `cmd/sponsor-panel/oauth.go` | Update all `pool` calls to use `*gorm.DB` |
| Modify | `cmd/sponsor-panel/patreon_oauth.go` | Update `upsertPatreonUser` call to GORM |
| Modify | `cmd/sponsor-panel/handlers.go` | Update `createLogoSubmission` call to GORM |
| Modify | `cmd/sponsor-panel/dashboard.go` | No direct DB calls, but `getSessionUser` flows through `oauth.go` |
| Modify | `cmd/sponsor-panel/sync_sponsors.go` | Update `upsertSponsorUsername`, `markInactiveSponsorsNotIn` to GORM |

---

### Task 1: Define GORM Models

Replace the raw structs and SQL helpers in `models.go` with GORM-tagged models and GORM-based DB functions.

**Files:**
- Rewrite: `cmd/sponsor-panel/models.go`

- [ ] **Step 1: Write the failing test**

Create `cmd/sponsor-panel/models_test.go` with a compilation test that imports the new model types and calls a GORM-specific method:

```go
package main

import (
	"testing"

	"gorm.io/gorm"
)

// TestModelsCompile verifies GORM model structs exist and have the expected shape.
func TestModelsCompile(t *testing.T) {
	// These should compile once models.go is rewritten
	var _ gorm.Model

	u := PanelUser{}
	if u.TableName() != "users" {
		t.Errorf("PanelUser.TableName() = %q, want %q", u.TableName(), "users")
	}

	ls := LogoSubmission{}
	if ls.TableName() != "logo_submissions" {
		t.Errorf("LogoSubmission.TableName() = %q, want %q", ls.TableName(), "logo_submissions")
	}

	su := SponsorUsername{}
	if su.TableName() != "github_sponsor_usernames" {
		t.Errorf("SponsorUsername.TableName() = %q, want %q", su.TableName(), "github_sponsor_usernames")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd cmd/sponsor-panel && go test -run TestModelsCompile -count=1 ./...`
Expected: compilation error (old `models.go` still has `pgxpool` imports, no `TableName()` methods)

- [ ] **Step 3: Rewrite `models.go` with GORM models**

Replace the entire file with GORM-tagged structs. Key changes:
- `User` -> `PanelUser` (avoids collision with `internal/models.Account`; table stays `users`)
- `db:"..."` tags -> `gorm:"..."` tags
- `pgxpool` helper functions -> GORM helper functions
- Keep `SponsorshipData` struct and `IsSponsorAtTier` method unchanged (they're pure logic)

```go
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
	GitHubID             *int64    `json:"github_id" gorm:"uniqueIndex"`
	PatreonID            *string   `json:"patreon_id" gorm:"uniqueIndex"`
	Provider             string    `json:"provider" gorm:"not null;default:'github'"`
	Login                string    `json:"login" gorm:"not null"`
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
	Username           string    `json:"username" gorm:"uniqueIndex;not null"`
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
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd cmd/sponsor-panel && go test -run TestModelsCompile -count=1 ./...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/sponsor-panel/models.go cmd/sponsor-panel/models_test.go
git commit -m "refactor(sponsor-panel): rewrite models with GORM tags and helpers"
```

---

### Task 2: Delete Raw SQL Migrations

Since GORM AutoMigrate handles schema creation, the hand-written `migrations.go` is no longer needed.

> **Dependency:** Task 1 must be committed first. Task 1 rewrites `models.go` which defines the DB helper functions that replaced the ones previously in `models.go`. Deleting `migrations.go` before Task 1 is committed will cause additional unresolvable compilation errors.

**Files:**
- Delete: `cmd/sponsor-panel/migrations.go`

- [ ] **Step 1: Delete `migrations.go`**

```bash
git rm cmd/sponsor-panel/migrations.go
```

- [ ] **Step 2: Verify the build still compiles**

Run: `cd cmd/sponsor-panel && go build ./...`
Expected: compilation errors -- `runMigrations` is called in `main.go` but no longer defined. This is expected and will be fixed in Task 3.

- [ ] **Step 3: Commit (allow build break, fixed in next task)**

```bash
git commit -m "refactor(sponsor-panel): remove raw SQL migrations (replaced by GORM AutoMigrate)"
```

---

### Task 3: Swap `pgxpool.Pool` for `gorm.DB` in `main.go`

Replace database connection setup, migration call, and `Server` struct to use GORM.

**Files:**
- Modify: `cmd/sponsor-panel/main.go`

- [ ] **Step 1: Update the `Server` struct**

Change `pool *pgxpool.Pool` to `db *gorm.DB`.

In the import block:
- Remove: `"github.com/jackc/pgx/v5/pgxpool"`
- Add: `"gorm.io/driver/postgres"` and `"gorm.io/gorm"`

```go
// Server holds the application dependencies.
type Server struct {
	db                *gorm.DB
	ghClient          *gh.Client
	oauth             *oauth2.Config
	patreonOAuth          *oauth2.Config
	patreonCampaignID     string
	patreonFiftyPlusSpons map[string]bool
	discordInvite     string
	fiftyPlusSponsors map[string]bool
	sessionStore      *sessions.CookieStore
	cookieSecure      bool
	bucketName        string
	s3Client          *s3.Client
}
```

- [ ] **Step 2: Replace database connection + migration in `main()`**

Replace lines 148-169 (the `pgxpool.New` / `pool.Ping` / `runMigrations` block) with:

```go
	// Connect to database via GORM
	slog.Debug("main: connecting to database via GORM")
	db, err := gorm.Open(postgres.Open(*databaseURL), &gorm.Config{})
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	slog.Info("main: database connection established")

	// Run GORM AutoMigrate
	slog.Debug("main: running GORM auto-migration")
	if err := db.AutoMigrate(PanelModels()...); err != nil {
		slog.Error("failed to auto-migrate", "err", err)
		os.Exit(1)
	}
	slog.Info("main: auto-migration completed")
```

- [ ] **Step 3: Update sync loop call**

Change:
```go
go startSyncLoop(syncCtx, pool, *githubToken)
```
To:
```go
go startSyncLoop(syncCtx, db, *githubToken)
```

- [ ] **Step 4: Update `Server` initialization**

Change `pool: pool,` to `db: db,`

Remove the `defer pool.Close()` line (GORM manages its own connection pool; if you want explicit close, use `sqlDB, _ := db.DB(); defer sqlDB.Close()`).

- [ ] **Step 5: Verify build compiles**

Run: `cd cmd/sponsor-panel && go build ./...`
Expected: compilation errors in `oauth.go`, `patreon_oauth.go`, `handlers.go`, `sync_sponsors.go` -- they still reference `s.pool` and old function signatures. This is expected and fixed in Tasks 4-6.

- [ ] **Step 6: Commit**

```bash
git add cmd/sponsor-panel/main.go
git commit -m "refactor(sponsor-panel): swap pgxpool for gorm.DB in Server and main()"
```

---

### Task 4: Update OAuth Handlers to Use GORM

Replace all `pool`-based DB calls in `oauth.go` with the new GORM helper functions.

**Files:**
- Modify: `cmd/sponsor-panel/oauth.go`

- [ ] **Step 1: Update imports**

Remove: `"github.com/jackc/pgx/v5/pgxpool"`

- [ ] **Step 2: Update `fetchSponsorship` signature**

Change:
```go
func fetchSponsorship(ctx context.Context, pool *pgxpool.Pool, token string, ...) (string, error) {
```
To:
```go
func fetchSponsorship(ctx context.Context, db *gorm.DB, token string, ...) (string, error) {
```

**Important:** Keep `ctx context.Context` as the first parameter. Although the GORM DB helpers no longer need it, `fetchSponsorship` internally calls `fetchSponsorshipForEntity(ctx, ...)` and `fetchUserOrganizationsWithSponsorship(ctx, ...)` which make raw HTTP requests and require a context for cancellation/timeout.

Update the `getActiveSponsorsByUsernames` call inside:
```go
// Old:
syncedSponsors, err := getActiveSponsorsByUsernames(ctx, pool, usernames)
// New:
syncedSponsors, err := getActiveSponsorsByUsernames(db, usernames)
```

- [ ] **Step 3: Update `callbackHandler`**

Change all `s.pool` references:
```go
// Old:
sponsorData, err := fetchSponsorship(r.Context(), s.pool, token.AccessToken, ...)
// New:
sponsorData, err := fetchSponsorship(r.Context(), s.db, token.AccessToken, ...)

// Old:
if err := upsertUser(r.Context(), s.pool, user); err != nil {
// New:
if err := upsertUser(s.db, user); err != nil {
```

**Critical: Cast `user.ID` to `int` when storing in session.** Since `PanelUser.ID` is `uint` but gorilla/sessions gob-encodes the value, and `getSessionUser` reads it back with `.(int)`, the type assertion will fail at runtime if `uint` is stored. Fix in `callbackHandler`:
```go
// Old:
session.Values["user_id"] = user.ID
// New:
session.Values["user_id"] = int(user.ID)
```

- [ ] **Step 4: Update `getSessionUser`**

Change:
```go
// Old:
return getUserByID(r.Context(), s.pool, userID)
// New (both call sites):
return getUserByID(s.db, userID)
```

- [ ] **Step 5: Update `User` references to `PanelUser`**

Throughout `oauth.go`, change:
- `user := &User{` -> `user := &PanelUser{`
- Any type annotations referencing `*User` -> `*PanelUser`

- [ ] **Step 6: Verify file compiles in isolation**

Run: `cd cmd/sponsor-panel && go vet ./...`

- [ ] **Step 7: Commit**

```bash
git add cmd/sponsor-panel/oauth.go
git commit -m "refactor(sponsor-panel): update oauth.go to use GORM helpers"
```

---

### Task 5: Update Patreon OAuth and Handlers

**Files:**
- Modify: `cmd/sponsor-panel/patreon_oauth.go`
- Modify: `cmd/sponsor-panel/handlers.go`

- [ ] **Step 1: Update `patreon_oauth.go`**

Change:
```go
// Old:
if err := upsertPatreonUser(r.Context(), s.pool, user); err != nil {
// New:
if err := upsertPatreonUser(s.db, user); err != nil {
```

Change `user := &User{` to `user := &PanelUser{`.

**Critical: Cast `user.ID` to `int` when storing in session** (same fix as `callbackHandler` in Task 4):
```go
// Old:
session.Values["user_id"] = user.ID
// New:
session.Values["user_id"] = int(user.ID)
```

- [ ] **Step 2: Update `handlers.go`**

Change:
```go
// Old:
if err := createLogoSubmission(r.Context(), s.pool, submission); err != nil {
// New:
if err := createLogoSubmission(s.db, submission); err != nil {
```

Update `LogoSubmission` struct field `UserID` from `int` to `uint` if needed (GORM primaryKey is `uint`). Check the `submission` initialization in `logoHandler` -- `user.ID` is now `uint`, so the assignment should work directly.

- [ ] **Step 3: Verify full build**

Run: `cd cmd/sponsor-panel && go build ./...`
Expected: PASS (compiles cleanly)

- [ ] **Step 4: Commit**

```bash
git add cmd/sponsor-panel/patreon_oauth.go cmd/sponsor-panel/handlers.go
git commit -m "refactor(sponsor-panel): update patreon_oauth and handlers to use GORM"
```

---

### Task 6: Update Sync Sponsors to Use GORM

**Files:**
- Modify: `cmd/sponsor-panel/sync_sponsors.go`

- [ ] **Step 1: Update function signatures**

Change:
```go
// Old:
func syncSponsors(ctx context.Context, pool *pgxpool.Pool, ghToken string) error {
// New:
func syncSponsors(ctx context.Context, db *gorm.DB, ghToken string) error {
```

```go
// Old:
func startSyncLoop(ctx context.Context, pool *pgxpool.Pool, ghToken string) {
// New:
func startSyncLoop(ctx context.Context, db *gorm.DB, ghToken string) {
```

- [ ] **Step 2: Update DB call sites**

```go
// Old:
if err := upsertSponsorUsername(ctx, pool, sponsor); err != nil {
// New:
if err := upsertSponsorUsername(db, sponsor); err != nil {

// Old:
inactiveCount, err := markInactiveSponsorsNotIn(ctx, pool, allSponsors)
// New:
inactiveCount, err := markInactiveSponsorsNotIn(db, allSponsors)
```

- [ ] **Step 3: Update imports**

Remove: `"github.com/jackc/pgx/v5/pgxpool"`
Add: `"gorm.io/gorm"` (if not already imported)

- [ ] **Step 4: Update `startSyncLoop` internal calls**

```go
// Old:
if err := syncSponsors(ctx, pool, ghToken); err != nil {
// New:
if err := syncSponsors(ctx, db, ghToken); err != nil {
```

- [ ] **Step 5: Full build + vet**

Run: `cd cmd/sponsor-panel && go build ./... && go vet ./...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add cmd/sponsor-panel/sync_sponsors.go
git commit -m "refactor(sponsor-panel): update sync_sponsors to use GORM"
```

---

### Task 7: Update Dashboard (Type Rename)

`dashboard.go` doesn't call DB directly but references the `User` type via `getSessionUser` and `SponsorshipData`.

**Files:**
- Modify: `cmd/sponsor-panel/dashboard.go`

- [ ] **Step 1: Update type references**

If `getSessionUser` now returns `*PanelUser`, update any explicit type annotations in `dashboard.go`. The `SponsorshipData` struct is unchanged so JSON unmarshaling still works. The `user.ID` field is now `uint` instead of `int` -- check that `session.Values["user_id"]` cast works (it stores `int`; you may need `int(user.ID)` when saving to session and `uint(userID)` when reading).

- [ ] **Step 2: Verify the session ID type handling**

In `oauth.go`'s `callbackHandler`, the session stores `user.ID`:
```go
session.Values["user_id"] = int(user.ID) // cast uint -> int for session compatibility
```

In `getSessionUser`:
```go
userID, ok := session.Values["user_id"].(int)
// ...
return getUserByID(s.db, userID) // getUserByID takes int, casts internally
```

This works because `getUserByID` takes `int` and GORM's `First(&user, userID)` handles the type.

- [ ] **Step 3: Full build and test**

Run: `cd cmd/sponsor-panel && go build ./... && go test ./...`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add cmd/sponsor-panel/dashboard.go cmd/sponsor-panel/oauth.go
git commit -m "refactor(sponsor-panel): fix type references for PanelUser rename"
```

---

### Task 8: Add GORM AutoMigrate Unique Index for Provider+Login

The old migration002 created a unique index on `(provider, login)`. GORM AutoMigrate won't create composite unique indexes from struct tags alone. Add an explicit index.

**Files:**
- Modify: `cmd/sponsor-panel/models.go`

- [ ] **Step 1: Add composite unique index to PanelUser**

Add a GORM hook or use the `gorm:"uniqueIndex:idx_users_provider_login"` tag on both fields:

```go
type PanelUser struct {
	ID                   uint      `json:"id" gorm:"primaryKey"`
	GitHubID             *int64    `json:"github_id" gorm:"uniqueIndex"`
	PatreonID            *string   `json:"patreon_id" gorm:"uniqueIndex"`
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
```

- [ ] **Step 2: Verify build**

Run: `cd cmd/sponsor-panel && go build ./...`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add cmd/sponsor-panel/models.go
git commit -m "refactor(sponsor-panel): add composite unique index for provider+login"
```

---

### Task 9: Clean Up Unused pgx Dependency

After all files are converted, `pgx` should no longer be imported by `cmd/sponsor-panel`.

**Files:**
- Verify: `cmd/sponsor-panel/*.go`

- [ ] **Step 1: Search for lingering pgx references**

Run: `grep -r "pgx\|pgxpool" cmd/sponsor-panel/`
Expected: no matches

- [ ] **Step 2: Run go mod tidy**

Run: `go mod tidy`

Note: `pgx` may still be in `go.mod` if other packages use it -- that's fine. We just want it removed from `cmd/sponsor-panel`'s imports.

- [ ] **Step 3: Full test suite**

Run: `npm test`
Expected: PASS (runs `go test ./...` for the whole project)

- [ ] **Step 4: Commit**

```bash
git add go.mod go.sum
git commit -m "chore: go mod tidy after sponsor-panel GORM migration"
```

---

### Task 10: Manual Integration Test

Since there's no test database fixture setup, verify the app starts and works against a real Postgres instance.

- [ ] **Step 1: Start the dev server**

Run: `npm run dev:sponsor-panel` (or equivalent with required env vars)

Check logs for:
- "database connection established"
- "auto-migration completed"
- No errors about missing columns or tables

- [ ] **Step 2: Test login flow**

Visit the app in a browser, complete OAuth login, verify:
- User record is created in `users` table
- Dashboard renders with correct sponsorship data
- Session persistence works across page loads

- [ ] **Step 3: Verify sponsor sync**

Check logs for:
- "syncSponsors: sync completed" within ~1 minute of startup
- `github_sponsor_usernames` table is populated
