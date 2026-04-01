# Thoth Token Issuance Card

## Context

Sponsors need Thoth API tokens to access Thoth services. Currently there is no
self-service way to get these tokens -- they must be manually provisioned. This
feature adds a card to the sponsor panel dashboard that lets any $1+/month
sponsor generate a Thoth token with one click.

The Thoth gRPC client is already wired into the sponsor panel (`server.thothClient`)
with `AdminUsers` and `AuthJWT` service clients. The missing pieces are: a
`ThothUserID` field on the user model, a handler to orchestrate user creation
and JWT minting, a template card, and a success response template.

## Approach: Lazy Thoth User Creation

Create the Thoth user on-demand when the sponsor first generates a token. Store
the Thoth user ID on the `PanelUser` model. Subsequent requests skip creation
and go straight to `MakeJWT`.

## Changes

### 1. Data Model (`models.go`)

Add field to `PanelUser`:

```go
ThothUserID *string `json:"thoth_user_id" gorm:"column:thoth_user_id"`
```

Nullable -- existing users get NULL. GORM AutoMigrate adds the column on next
startup.

### 2. Handler (`handlers.go`)

New function: `func (s *Server) thothTokenHandler(w http.ResponseWriter, r *http.Request)`

Flow:

1. Require POST method
2. Get session user via `s.getSessionUser(r)`
3. Check `user.IsSponsorAtTier(100)` (any $1+ sponsor)
4. If `user.ThothUserID == nil`:
   - Call `s.thothClient.AdminUsers.Create(ctx, &adminv1.UsersServiceCreateRequest{...})`
     - `EmailAddress`: `user.Email`
     - `Name`: `user.Login`
     - `CustomerId`: `user.Provider + ":" + user.Login`
   - Save `resp.User.Id` to `user.ThothUserID` in the database
5. Call `s.thothClient.AdminUsers.MakeJWT(ctx, &adminv1.UsersServiceMakeJWTRequest{...})`
   - `UserId`: `*user.ThothUserID`
   - `Comment`: `"sponsor-panel token for " + user.Login`
6. Render `templates.ThothTokenSuccess(resp.TokenInfo.Jwt)` on success
7. Use `renderError()` for all error paths (same pattern as other handlers)

### 3. Templates (`templates/dashboard.templ`)

New card component:

```
templ ThothTokenCard()
```

- Default card style (no color variant, like LogoSubmitCard)
- Key/token icon in header
- Title: "Thoth API Token"
- Description: "Generate an API token for Thoth services."
- Single button: "Generate Token"
- `hx-post="/thoth-token"` `hx-target="#thoth-result"`
- Result div: `<div id="thoth-result"></div>`

New success component (in `formresult.templ` or a new file):

```
templ ThothTokenSuccess(token string)
```

- Renders inside an `alert-success` div
- Contains a `<pre>` block with:
  ```
  THOTH_URL=passthrough:///thoth.techaro.lol:443
  THOTH_TOKEN=<token>
  ```
- Monospace styling so users can copy-paste

### 4. Dashboard Layout (`templates/dashboard.templ`)

Add `ThothTokenCard()` to the second row grid, shown when `props.IsSponsor` is
true:

```templ
<div class="grid md:grid-cols-2 gap-8 mt-8">
    if props.IsFiftyPlus {
        @TeamInviteCard()
    }
    if props.IsSponsor {
        @LogoSubmitCard()
    }
    if props.IsSponsor {
        @ThothTokenCard()
    }
</div>
```

If a $50+ sponsor sees all three cards (TeamInvite, LogoSubmit, ThothToken),
add a third row for ThothTokenCard to keep the 2-column grid clean.

### 5. Routing (`main.go`)

Add route:

```go
mux.HandleFunc("/thoth-token", server.thothTokenHandler)
```

Add to the debug routes list.

## Files Modified

- `cmd/sponsor-panel/models.go` -- add ThothUserID field to PanelUser
- `cmd/sponsor-panel/handlers.go` -- add thothTokenHandler + renderThothSuccess
- `cmd/sponsor-panel/templates/dashboard.templ` -- add ThothTokenCard, wire into layout
- `cmd/sponsor-panel/templates/formresult.templ` -- add ThothTokenSuccess component
- `cmd/sponsor-panel/main.go` -- add /thoth-token route

## Existing Code to Reuse

- `server.getSessionUser(r)` -- session auth (`oauth.go:634`)
- `user.IsSponsorAtTier(100)` -- tier check (`models.go:38`)
- `renderError(w, msg, code)` -- error rendering (`handlers.go:302`)
- `templates.FormResult(msg, bool)` -- generic result component (`formresult.templ`)
- `server.thothClient.AdminUsers` -- gRPC client (`internal/thoth/thoth.go`)
- Generated proto types from `xeiaso.net/v4/gen/techaro/thoth/auth/admin/v1`

## Verification

1. `go build ./cmd/sponsor-panel` -- compiles without errors
2. `npm test` -- all tests pass
3. `npm run dev:sponsor-panel` -- start dev server
4. Log in as a sponsor, verify the Thoth card appears
5. Click "Generate Token", verify credentials are displayed
6. Click again, verify it works without creating a duplicate Thoth user
7. Check database: `ThothUserID` column populated after first generation
