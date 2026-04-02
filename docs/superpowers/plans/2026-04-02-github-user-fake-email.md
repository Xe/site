# Plan: Fake email shim + HTMX error display for sponsor-panel

## Context

GitHub users without a public email address hit an "Email address required. Please update your profile." error when requesting a Thoth token. Additionally, error responses from HTMX POST handlers never appear in the card because HTMX is configured to not swap 4xx/5xx responses (`responseHandling: [{ code: "[45]..", swap: false, error: true }]` in htmx.js:64).

## Changes

### 1. Generate fake email at OAuth callback (`cmd/sponsor-panel/oauth.go`)

After line 531 (after `ghUser` is fetched successfully), add:

```go
if ghUser.Email == "" {
    ghUser.Email = ghUser.Login + "@fake-address.github"
    slog.Info("callbackHandler: generated fake email for user", "login", ghUser.Login, "email", ghUser.Email)
}
```

This prevents new users from ever storing an empty email.

### 2. Generate fake email at token issuance (`cmd/sponsor-panel/handlers.go`, lines 352-356)

Replace the error block:

```go
if user.Email == "" {
    slog.Error(...)
    renderError(w, "Email address required...", http.StatusBadRequest)
    return
}
```

With fake email generation + DB save:

```go
if user.Email == "" {
    user.Email = user.Login + "@fake-address.github"
    slog.Info("thothTokenHandler: generated fake email for user", "user_id", user.ID, "login", user.Login, "email", user.Email)
    if err := s.db.Save(user).Error; err != nil {
        slog.Error("thothTokenHandler: failed to save fake email", "err", err, "user_id", user.ID)
        renderError(w, "Failed to save user email")
        return
    }
}
```

This covers existing users who already have empty emails in the DB.

### 3. Fix `renderError` to return 200 (`cmd/sponsor-panel/handlers.go`, lines 303-308)

Change `renderError` to always return HTTP 200. HTMX won't swap non-2xx responses, so the current 4xx/5xx status codes mean error messages never appear in the `#thoth-result`, `#invite-result`, or `#logo-result` divs. The error styling is already handled by `FormResult(message, false)` rendering a red alert box.

```go
func renderError(w http.ResponseWriter, message string) {
    w.Header().Set("Content-Type", "text/html")
    w.WriteHeader(http.StatusOK)
    templates.FormResult(message, false).Render(context.Background(), w)
}
```

Remove the `statusCode` parameter from all call sites since it's no longer used.

## Files to modify

- `cmd/sponsor-panel/oauth.go` -- add fake email after line 531
- `cmd/sponsor-panel/handlers.go` -- fake email at lines 352-356, fix `renderError` signature + all call sites

## Verification

1. `go build ./cmd/sponsor-panel/` compiles
2. `go vet ./cmd/sponsor-panel/`
3. `go test ./cmd/sponsor-panel/...`
