# Plan: Add Patreon OAuth Login to Sponsor Panel

## Context

The sponsor-panel (`cmd/sponsor-panel/`) currently only supports GitHub OAuth login. Patreon patrons have no way to access sponsor benefits (Discord invite, team invitations, logo submissions). This change adds Patreon as a second authentication provider so patrons get feature parity with GitHub Sponsors.

**Design decisions:**

- User-token based verification (no saasproxy dependency)
- Separate identities (no account linking between GitHub and Patreon)
- Patreon API v2 via direct HTTP calls; `patreon-go` library used only for OAuth URL constants
- OAuth scopes: `identity`, `identity[email]`, `campaigns.members`
- Team invite: Patreon $50+ users see the same card and input a GitHub username

---

## Step 1: Database Migration

**File:** `cmd/sponsor-panel/migrations.go`

Add `migration002` constant with idempotent DDL:

- `ALTER TABLE users ALTER COLUMN github_id DROP NOT NULL`
- `ADD COLUMN IF NOT EXISTS patreon_id TEXT UNIQUE`
- `ADD COLUMN IF NOT EXISTS provider TEXT NOT NULL DEFAULT 'github'`
- Drop `users_login_key` unique constraint, replace with `UNIQUE INDEX (provider, login)`
- Add `idx_users_patreon_id` index

Update `runMigrations()` to execute `migration002` after `migrationSchema`.

---

## Step 2: Update User Model

**File:** `cmd/sponsor-panel/models.go`

- Change `User.GitHubID` from `int64` to `*int64` (nullable)
- Add `PatreonID *string` and `Provider string` fields
- Update all `SELECT`/`Scan` calls in `getUserByID`, `upsertUser` to include new columns
- Add `upsertPatreonUser(ctx, pool, user)` — same pattern as `upsertUser` but upserts by `patreon_id`, sets `provider='patreon'`

---

## Step 3: Add Patreon OAuth Config

**File:** `cmd/sponsor-panel/main.go`

New flags (all optional — service still works GitHub-only if omitted):

- `--patreon-client-id`
- `--patreon-client-secret`
- `--patreon-redirect-url`
- `--patreon-campaign-id` (to match pledge against)

Add to `Server` struct:

```go
patreonOAuth      *oauth2.Config  // nil if not configured
patreonCampaignID string
```

In `main()`, conditionally create `oauth2.Config` using `patreon.AuthorizationURL` and `patreon.AccessTokenURL` from `gopkg.in/mxpv/patreon-go.v1`.

Register new routes:

```
/login/patreon    → server.patreonLoginHandler
/callback/patreon → server.patreonCallbackHandler
```

---

## Step 4: Patreon OAuth Handlers (new file)

**File:** `cmd/sponsor-panel/patreon_oauth.go` (new)

### `patreonLoginHandler`

Mirrors `loginHandler` in `oauth.go`: generate state, set CSRF cookie, redirect to `s.patreonOAuth.AuthCodeURL(state)`. Returns 404 if `s.patreonOAuth == nil`.

### `patreonCallbackHandler`

1. Validate state cookie (same CSRF pattern as GitHub)
2. Exchange code via `s.patreonOAuth.Exchange(ctx, code)`
3. Call Patreon API v2 identity endpoint:
   ```
   GET https://www.patreon.com/api/oauth2/v2/identity
     ?include=memberships.campaign
     &fields[user]=full_name,vanity,email,image_url
     &fields[member]=patron_status,currently_entitled_amount_cents
   ```
4. Parse JSON:API response to extract user identity and membership data
5. Find membership matching `s.patreonCampaignID`
6. Build `SponsorshipData` JSON: `{is_active, monthly_amount_cents, tier_name}`
7. Call `upsertPatreonUser()` with `provider="patreon"`, login = vanity or full_name
8. Create session (same `user_id` in gorilla/sessions cookie)
9. Redirect to `/`

### Response types (define in same file):

- `patreonIdentityResponse` — JSON:API envelope with user data + included memberships
- `patreonMember` — membership attributes (patron_status, currently_entitled_amount_cents) + campaign relationship

---

## Step 5: Update Login Template

**File:** `cmd/sponsor-panel/templates/login.templ`

Change signature to `templ Login(patreonEnabled bool)`. Add a "Login with Patreon" button (with Patreon SVG icon) conditionally rendered when `patreonEnabled` is true, linking to `/login/patreon`.

---

## Step 6: Update Dashboard for Provider Awareness

**File:** `cmd/sponsor-panel/templates/dashboard.templ`

- Add `Provider string` to `UserProps`
- In `SponsorshipCard`, when user is not a sponsor: show Patreon link for `provider == "patreon"`, GitHub Sponsors link otherwise

**File:** `cmd/sponsor-panel/dashboard.go`

- `loginPageHandler`: pass `s.patreonOAuth != nil` to `templates.Login()`
- `dashboardHandler`: set `UserProps.Provider` from `user.Provider`

---

## Step 7: Generate & Build

1. `go tool templ generate` (regenerate `*_templ.go` files)
2. `go build ./cmd/sponsor-panel/`
3. `npm test` (`go test ./...`)

---

## Files Modified

| File                                          | Change                                          |
| --------------------------------------------- | ----------------------------------------------- |
| `cmd/sponsor-panel/migrations.go`             | Add migration002                                |
| `cmd/sponsor-panel/models.go`                 | Update User struct, add upsertPatreonUser       |
| `cmd/sponsor-panel/main.go`                   | Add flags, Server fields, routes                |
| `cmd/sponsor-panel/patreon_oauth.go`          | **New** — Patreon login/callback handlers       |
| `cmd/sponsor-panel/oauth.go`                  | No changes (existing GitHub flow untouched)     |
| `cmd/sponsor-panel/dashboard.go`              | Pass patreonEnabled and Provider                |
| `cmd/sponsor-panel/templates/login.templ`     | Add Patreon button                              |
| `cmd/sponsor-panel/templates/dashboard.templ` | Add Provider to props, conditional sponsor link |

## Existing Code to Reuse

- `generateState()` in `oauth.go:20` — reuse for Patreon CSRF state
- `SponsorshipData` struct in `models.go:27` — same JSON format for both providers
- `User.IsSponsorAtTier()` in `models.go:35` — works provider-agnostically
- Session management via `getSessionUser()` in `oauth.go:634` — no changes needed
- `patreon.AuthorizationURL` / `patreon.AccessTokenURL` from `gopkg.in/mxpv/patreon-go.v1`
- All feature handlers (`inviteHandler`, `logoHandler`) — work unchanged since they only check `IsSponsorAtTier()`

## Verification

1. **Build**: `go build ./cmd/sponsor-panel/` compiles without errors
2. **Tests**: `npm test` passes
3. **GitHub flow unchanged**: Login with GitHub still works identically
4. **Patreon login**: With Patreon OAuth credentials set, clicking "Login with Patreon" redirects to Patreon, callback creates user with `provider=patreon`
5. **Tier gating**: A Patreon user pledging $50+/month sees team invite card; $1+ sees logo submission and Discord
6. **No Patreon config**: When Patreon flags are omitted, login page only shows GitHub button, `/login/patreon` returns 404
