# Sponsor Panel - Simplified UX Specification

**Design philosophy:** Minimal UI for a minimal backend. No over-engineering.

---

## Overview

This UX specification matches the simplified backend:
- **2 tables** (users, logo_submissions)
- **5 endpoints** (/login, /callback, /, /invite, /logo)
- **No external storage** (GitHub attachments only)
- **No approval workflows**
- **No audit trails**

The UI is simple, functional, and matches what the backend actually does.

---

## Pages (3 Total)

### 1. Login Page (`GET /`)

Public page that redirects to GitHub OAuth.

**Layout:**
```html
<h1>Sponsor Panel</h1>
<p>Manage your sponsorship benefits</p>
<a href="/login" class="btn">Login with GitHub</a>

<div class="benefits">
  <h3>Benefits by tier:</h3>
  <ul>
    <li>✓ All sponsors: Discord access</li>
    <li>✓ $50+/month: Team invitations</li>
    <li>✓ All sponsors: Logo submission</li>
  </ul>
</div>
```

**Behavior:**
- If authenticated: redirect to `/`
- If not authenticated: show login page

---

### 2. Dashboard (`GET /`)

Main dashboard for authenticated sponsors. All sections on one page.

**Layout:**
```html
<nav>
  <img src="{avatar}" alt="">
  <span>{username}</span>
  <a href="/logout">Logout</a>
</nav>

<main>
  <!-- Section 1: Discord (always visible) -->
  <section id="discord">
    <h2>Discord Community</h2>
    <p>Join our Discord server to chat with other sponsors.</p>
    <a href="{discordInvite}" target="_blank" class="btn">Join Discord</a>
  </section>

  <!-- Section 2: Sponsorship Status -->
  <section id="sponsorship">
    <h2>Your Sponsorship</h2>
    {if isSponsor}
      <p class="success">You sponsor at ${amount}/month - thank you!</p>
    {else}
      <p>Not an active sponsor.</p>
      <a href="https://github.com/sponsors/Xe" target="_blank">Become a Sponsor</a>
    {/if}
  </section>

  <!-- Section 3: Team Invitation ($50+ sponsors only) -->
  {if isFiftyPlus}
  <section id="invite">
    <h2>Team Invitation</h2>
    <p>Invite team members to the TecharoHQ organization.</p>
    <form hx-post="/invite" hx-target="#invite-result">
      <input type="text" name="username" placeholder="GitHub username" required>
      <button type="submit">Invite to Team</button>
    </form>
    <div id="invite-result"></div>
  </section>
  {/if}

  <!-- Section 4: Logo Submission (all sponsors) -->
  {if isSponsor}
  <section id="logo">
    <h2>Logo Submission</h2>
    <p>Submit your company logo for the Anubis project README.</p>
    <form hx-post="/logo" hx-encoding="multipart/form-data" hx-target="#logo-result">
      <input type="text" name="company" placeholder="Company Name" required>
      <input type="url" name="website" placeholder="Website URL" required>
      <input type="file" name="logo" accept="image/png,image/jpeg,image/svg+xml" required>
      <button type="submit">Submit Logo</button>
    </form>
    <div id="logo-result"></div>
  </section>
  {/if}
</main>
```

**Behavior:**
- If not authenticated: redirect to `/login`
- If authenticated: render dashboard with sections based on sponsorship tier

---

### 3. OAuth Error (`/callback` with error)

Simple error page when OAuth fails.

**Layout:**
```html
<h1>Authentication Failed</h1>
<p>{errorMessage}</p>
<a href="/login" class="btn">Try Again</a>
```

**Error messages:**
- "Invalid OAuth state" - CSRF mismatch
- "Failed to exchange token" - GitHub API error
- "Failed to fetch user" - GitHub API error

---

## Templ Components (5 Total)

### base.templ

Layout wrapper with head, styles, and body structure.

```templ
templ base(title string, content templ.Component) {
    <!DOCTYPE html>
    <html lang="en">
        <head>
            <meta charset="UTF-8"/>
            <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
            <title>{ title }</title>
            <script src="https://unpkg.com/htmx.org@1.9.10"></script>
            <script src="https://cdn.tailwindcss.com"></script>
        </head>
        <body class="bg-gray-50 dark:bg-gray-900">
            @content
        </body>
    </html>
}
```

### Login.templ

Login page with GitHub OAuth button.

```templ
templ Login(discordInvite string) {
    @base("Sponsor Panel - Login") {
        <div class="min-h-screen flex items-center justify-center">
            <div class="max-w-md w-full bg-white dark:bg-gray-800 rounded-xl p-8 shadow-lg">
                <h1 class="text-2xl font-bold mb-2">Sponsor Panel</h1>
                <p class="text-gray-600 dark:text-gray-400 mb-6">Manage your sponsorship benefits</p>

                <a href="/login" class="block w-full bg-gray-900 hover:bg-gray-800 text-white text-center py-3 rounded-lg mb-6">
                    Login with GitHub
                </a>

                <div class="border-t pt-4">
                    <h3 class="font-semibold mb-2">Benefits:</h3>
                    <ul class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
                        <li>✓ All sponsors: Discord access</li>
                        <li>✓ $50+/month: Team invitations</li>
                        <li>✓ All sponsors: Logo submission</li>
                    </ul>
                </div>
            </div>
        </div>
    }
}
```

### Dashboard.templ

Main dashboard with conditional sections.

```templ
templ Dashboard(props DashboardProps) {
    @base("Sponsor Panel - Dashboard") {
        @Navbar(props.User.Login, props.User.AvatarURL)

        <main class="max-w-4xl mx-auto px-4 py-8">
            <h1 class="text-2xl font-bold mb-6">Welcome, @props.User.Login!</h1>

            <div class="grid md:grid-cols-2 gap-6">
                @DiscordCard(props.DiscordInvite)
                @SponsorshipCard(props.IsSponsor, props.SponsorAmount, props.SponsorTier)
            </div>

            <div class="grid md:grid-cols-2 gap-6 mt-6">
                if props.IsFiftyPlus {
                    @TeamInviteCard()
                }
                if props.IsSponsor {
                    @LogoSubmitCard()
                }
            </div>
        </main>
    }
}

templ Navbar(login, avatarURL string) {
    <nav class="bg-white dark:bg-gray-800 border-b px-4 py-3">
        <div class="max-w-4xl mx-auto flex items-center justify-between">
            <div class="flex items-center gap-3">
                <img src={ avatarURL } class="w-8 h-8 rounded-full" alt="">
                <span class="font-semibold">{ login }</span>
            </div>
            <a href="/logout" class="text-sm text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white">
                Logout
            </a>
        </div>
    </nav>
}

templ DiscordCard(inviteURL string) {
    <div class="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
        <h2 class="font-semibold mb-2">Discord Community</h2>
        <p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
            Join our Discord server to chat with other sponsors.
        </p>
        <a href={ inviteURL } target="_blank" class="inline-block bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-lg text-sm">
            Join Discord
        </a>
    </div>
}

templ SponsorshipCard(isSponsor bool, amount int, tier string) {
    <div class="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
        <h2 class="font-semibold mb-2">Your Sponsorship</h2>
        if isSponsor {
            <p class="text-green-600 dark:text-green-400">
                ${ fmt.Sprintf("%.2f", float64(amount) / 100) }/month - { tier }
            </p>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">Thank you for your support!</p>
        } else {
            <p class="text-gray-600 dark:text-gray-400">Not an active sponsor.</p>
            <a href="https://github.com/sponsors/Xe" target="_blank" class="inline-block bg-pink-600 hover:bg-pink-700 text-white px-4 py-2 rounded-lg text-sm mt-2">
                Become a Sponsor
            </a>
        }
    </div>
}

templ TeamInviteCard() {
    <div class="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
        <h2 class="font-semibold mb-2">Team Invitation</h2>
        <p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
            Invite team members to TecharoHQ.
        </p>
        <form hx-post="/invite" hx-target="#invite-result" class="space-y-3">
            <input type="text" name="username" placeholder="GitHub username" required
                class="w-full border rounded-lg px-3 py-2 dark:bg-gray-700 dark:border-gray-600">
            <button type="submit" class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm w-full">
                Invite to Team
            </button>
        </form>
        <div id="invite-result"></div>
    </div>
}

templ LogoSubmitCard() {
    <div class="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
        <h2 class="font-semibold mb-2">Logo Submission</h2>
        <p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
            Submit your company logo for the Anubis README.
        </p>
        <form hx-post="/logo" hx-encoding="multipart/form-data" hx-target="#logo-result" class="space-y-3">
            <input type="text" name="company" placeholder="Company Name" required
                class="w-full border rounded-lg px-3 py-2 dark:bg-gray-700 dark:border-gray-600">
            <input type="url" name="website" placeholder="Website URL" required
                class="w-full border rounded-lg px-3 py-2 dark:bg-gray-700 dark:border-gray-600">
            <input type="file" name="logo" accept="image/png,image/jpeg,image/svg+xml" required
                class="w-full border rounded-lg px-3 py-2 dark:bg-gray-700 dark:border-gray-600">
            <button type="submit" class="bg-orange-600 hover:bg-orange-700 text-white px-4 py-2 rounded-lg text-sm w-full">
                Submit Logo
            </button>
        </form>
        <div id="logo-result"></div>
    </div>
}
```

### OAuthError.templ

Error page for OAuth failures.

```templ
templ OAuthError(message string) {
    @base("Authentication Error") {
        <div class="min-h-screen flex items-center justify-center">
            <div class="max-w-md w-full bg-red-50 dark:bg-red-900/20 rounded-xl p-8 text-center">
                <h1 class="text-xl font-bold text-red-900 dark:text-red-100 mb-2">Authentication Failed</h1>
                <p class="text-red-700 dark:text-red-300 mb-6">{ message }</p>
                <a href="/login" class="inline-block bg-red-600 hover:bg-red-700 text-white px-6 py-2 rounded-lg">
                    Try Again
                </a>
            </div>
        </div>
    }
}
```

### FormResult.templ

Generic success/error message for form submissions.

```templ
templ FormResult(message string, isSuccess bool) {
    if isSuccess {
        <div class="bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-300 p-4 rounded-lg">
            { message }
        </div>
    } else {
        <div class="bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-300 p-4 rounded-lg">
            { message }
        </div>
    }
}
```

---

## HTMX Interactions (5 Patterns)

### 1. Login Redirect

```
User clicks "Login with GitHub"
  → GET /login
  → Redirect to GitHub OAuth
```

### 2. OAuth Callback

```
GitHub redirects to /callback
  → Exchange code for token
  → Fetch user + sponsorship
  → Create session cookie
  → Redirect to /
```

### 3. Dashboard Load

```
GET / (with session cookie)
  → Validate session
  → Fetch user from database
  → Render Dashboard.templ
```

### 4. Team Invitation

```
User submits form (POST /invite)
  → Validate sponsorship ($50+)
  → Call GitHub Teams API
  → Return FormResult component
  → HTMX swaps #invite-result
```

### 5. Logo Submission

```
User submits form (POST /logo)
  → Validate form fields
  → Create GitHub issue with attachment
  → Return FormResult component with issue link
  → HTMX swaps #logo-result
```

---

## Error Handling (3 Patterns)

### Pattern 1: OAuth Error

```go
// In callback handler
if err != nil {
    components.OAuthError("Authentication failed: " + err.Error()).Render(r.Context(), w)
    return
}
```

### Pattern 2: Form Inline Error

```go
// In invite/logo handler
if err != nil {
    components.FormResult("Failed: " + err.Error(), false).Render(r.Context(), w)
    return
}
```

### Pattern 3: Not Authenticated

```go
// In dashboard handler
if session == nil {
    http.Redirect(w, r, "/login", http.StatusFound)
    return
}
```

---

## Success Responses (2 Patterns)

### Team Invitation Success

```templ
templ FormResult("Invitation sent to @" + username + "!", true)
```

### Logo Submission Success

```templ
templ LogoSubmissionSuccess(companyName, issueURL string, issueNumber int) {
    <div class="bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-300 p-4 rounded-lg">
        <p class="font-semibold">Logo submitted!</p>
        <p class="text-sm mt-1">
            View issue: <a href={ issueURL } target="_blank" class="underline">#{ issueNumber }</a>
        </p>
    </div>
}
```

---

## State Management (2 States)

### 1. Authentication State

```go
// Checked via session cookie
session, _ := r.Cookie("session")
if session == nil {
    // Not authenticated
}
```

### 2. Sponsorship State

```go
// Fetched once at login, stored in database
user.IsSponsorAtTier(5000) // $50+
user.IsSponsorAtTier(100)  // $1+
```

**No other state.** No loading spinners, no skeleton screens, no background polling.

---

## Responsive Design

Use Tailwind's responsive classes:

```html
<!-- Stack on mobile, grid on desktop -->
<div class="grid md:grid-cols-2 gap-6">
    <!-- Cards stack on mobile, 2 columns on desktop -->
</div>
```

No custom breakpoints. No hamburger menus. Just Tailwind utilities.

---

## File Structure

```
cmd/sponsor-panel/
├── main.go
├── templates/
│   ├── base.templ
│   ├── Login.templ
│   ├── Dashboard.templ
│   ├── OAuthError.templ
│   └── FormResult.templ
└── static/
    └── (none - using CDN for HTMX/Tailwind)
```

---

## Summary

**Original UX_FLOWS.md:**
- 1,562 lines
- 20+ components
- 10+ HTMX endpoints
- Complex state management
- Multiple pages with nested flows

**Simplified UX:**
- 3 pages
- 5 components
- 5 endpoints (matching backend)
- 2 states (auth + sponsorship)
- Simple forms with inline feedback

This matches the simplified backend: **minimal, functional, no over-engineering**.

---

_Simplified UX specification, 2026-02-07_
