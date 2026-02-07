# Sponsor Panel Service - Executive Summary

**Service:** sponsor-panel
**Exposure:** sponsors.xeiaso.net
**Status:** Specification Complete (Simplified)

---

## Purpose

The sponsor-panel service provides GitHub sponsors with a self-service dashboard to manage their sponsorship benefits. The service verifies sponsorship status via GitHub Sponsors GraphQL API and unlocks features based on contribution tier.

**Design philosophy:** Build the simplest thing that works. Both technical and UX specifications have been ruthlessly simplified from original over-engineered versions.

---

## Feature Matrix

| Feature                                   | $0  | $1-49 | $50+ |
| ----------------------------------------- | :-: | :---: | :--: |
| Discord invite link                       | ✅  |  ✅   |  ✅  |
| Sponsorship status display                |  -  |  ✅   |  ✅  |
| Team invitation to `botstopper-customers` |  -  |   -   |  ✅  |
| Logo submission for Anubis README         |  -  |  ✅   |  ✅  |

**Removed from original specs:**

- Organization-specific teams, Invitation audit trail, Logo approval workflow, Image processing, Real-time validation, Activity feeds, Multiple pages, Toast notifications, Locked/unlocked states

---

## Technical Stack (Simplified)

| Layer        | Technology                  | Notes                 |
| ------------ | --------------------------- | --------------------- |
| **Backend**  | Go with `net/http`          | Standard library only |
| **Frontend** | Templ + HTMX + Tailwind CSS | Unchanged             |
| **Database** | PostgreSQL with sqlx        | No ORM                |
| **Storage**  | GitHub issue attachments    | No external storage   |
| **Auth**     | GitHub OAuth 2.0            | No PKCE (server-side) |
| **Sessions** | Cookie-only                 | No database backing   |

**Removed dependencies:** gorilla/sessions, disintegration/imaging, aws-sdk-go-v2, gorm

---

## Complexity Reduction

| Aspect                | Original | Simplified | Change |
| --------------------- | -------- | ---------- | ------ |
| **Tables**            | 6        | 2          | -67%   |
| **Dependencies**      | 15+      | 6          | -60%   |
| **Endpoints**         | 10+      | 5          | -50%   |
| **Pages**             | 7+       | 3          | -57%   |
| **Components**        | 20+      | 5          | -75%   |
| **External services** | 3        | 2          | -33%   |
| **Code estimate**     | ~4000    | ~1500      | -62%   |

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                     USER BROWSER (HTMX)                             │
└─────────────────────────────┬───────────────────────────────────────┘
                              │
                              ▼
┌───────────────────────────────────────────────────────────────────┐
│                       sponsor-panel Service                       │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐   │
│  │   OAuth    │  │   Sponsor  │  │   Team     │  │    Logo    │   │
│  │  Handler   │  │  Checker   │  │  Inviter   │  │  Submitter │   │
│  └──────┬─────┘  └──────┬─────┘  └──────┬─────┘  └──────┬─────┘   │
└─────────┼───────────────┼───────────────┼───────────────┼─────────┘
          │               │               │               │
          ▼               ▼               ▼               ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     sqlx (PostgreSQL)                               │
│  • users              • logo_submissions                            │
└─────────────────────────────────────────────────────────────────────┘
          │               │               │
          ▼               ▼               ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     GitHub APIs                                     │
│  • OAuth 2.0  • GraphQL (Sponsors)  • REST (Teams, Issues)          │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Data Model (Simplified)

| Model            | Purpose                                                |
| ---------------- | ------------------------------------------------------ |
| `User`           | GitHub OAuth users with JSON-embedded sponsorship data |
| `LogoSubmission` | Logo submissions with GitHub issue tracking            |

**Removed tables:** Organization, SponsorshipCache, TeamInvitation, Session, OAuthState

---

## Key Integrations (Simplified)

### GitHub OAuth

- Scopes: `read:user`, `user:email`, `read:org`, `read:sponsors`
- Session: Encrypted cookie only (no database)

### GitHub Sponsors (GraphQL)

- Single query on login, cached in user table
- No separate cache table or expiration logic

### Team Management

- Target: `TecharoHQ/botstopper-customers` only
- No audit trail or invitation history

### Logo Submission

- Create GitHub issue in `TecharoHQ/anubis`
- Attach logo directly to issue (no external storage)

---

## Endpoints (Simplified)

| Endpoint    | Method | Auth | Purpose                        |
| ----------- | ------ | ---- | ------------------------------ |
| `/login`    | GET    | No   | Initiate GitHub OAuth          |
| `/callback` | GET    | No   | OAuth callback, create session |
| `/`         | GET    | Yes  | Dashboard                      |
| `/invite`   | POST   | Yes  | Invite user to team            |
| `/logo`     | POST   | Yes  | Submit logo (create issue)     |

---

## UX Overview (Simplified)

### Pages (3 total)

1. **Login page** - GitHub OAuth button
2. **Dashboard** - All sections on one page
3. **OAuth error** - Simple error page

### Components (5 total)

- `base.templ` - Layout wrapper
- `Login.templ` - Login page
- `Dashboard.templ` - Main dashboard
- `OAuthError.templ` - Error page
- `FormResult.templ` - Form success/error messages

### Dashboard Sections

- Discord invite (always visible)
- Sponsorship status
- Team invitation ($50+ sponsors only)
- Logo submission (all sponsors)

---

## Configuration (Simplified)

### Required Environment Variables

| Variable               | Purpose                         |
| ---------------------- | ------------------------------- |
| `DATABASE_URL`         | PostgreSQL connection           |
| `GITHUB_CLIENT_ID`     | OAuth app ID                    |
| `GITHUB_CLIENT_SECRET` | OAuth app secret                |
| `GITHUB_TOKEN`         | GitHub token for team/issue ops |
| `SESSION_KEY`          | Session encryption key          |
| `DISCORD_INVITE`       | Discord invite link             |
| `OAUTH_REDIRECT_URL`   | OAuth callback URL              |

---

## Documentation

| Document                 | Lines | Description                                                     |
| ------------------------ | ----- | --------------------------------------------------------------- |
| **SPEC.md**              | 975   | **(Simplified)** Technical: 2 tables, 5 endpoints, sqlx, no ORM |
| **UX_FLOWS.md**          | 522   | **(Simplified)** UX: 3 pages, 5 components, simple forms        |
| **ORIGINAL_SPEC.md**     | 909   | Original technical spec (over-engineered reference)             |
| **ORIGINAL_UX_FLOWS.md** | 1,562 | Original UX spec (over-engineered reference)                    |

---

## Implementation Checklist

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

- [ ] Build 5 Templ components
- [ ] Implement `/invite` with GitHub API
- [ ] Implement `/logo` with GitHub issue creation
- [ ] Add error handling components

### Phase 5: Testing & Deploy

- [ ] Unit tests for handlers
- [ ] Integration tests with test database
- [ ] Deploy to Kubernetes

---

## Next Steps

1. **Review the simplified specs:**
   - `SPEC.md` - Technical implementation
   - `UX_FLOWS.md` - UI/UX components and flows

2. **Set up development:**
   - Configure PostgreSQL database
   - Create GitHub OAuth app
   - No storage setup needed (uses GitHub)

3. **Begin implementation:**
   - Start with Phase 1 (Foundation)
   - Follow the simplified specs
   - Build incrementally

---

_Simplified specifications based on ruthless code review, 2026-02-07_
