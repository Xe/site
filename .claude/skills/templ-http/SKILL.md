---
name: templ-http
description: Integrate templ components with Go HTTP server using net/http. Use when connecting templ to web server, creating HTTP handlers, mentions 'templ server', 'HTTP routes', or 'serve templ components'.
---

# Templ HTTP Integration

Use progressive disclosure: begin with the handler pattern, then add routing and middleware only as needed.

## Level 1: Minimal Integration

Use this skill when serving templ from `net/http`.

1. Parse request input.
2. Build view model/data.
3. Render templ component with `r.Context()`.
4. Handle render errors and status codes.

```go
func homeHandler(w http.ResponseWriter, r *http.Request) {
    if err := components.HomePage().Render(r.Context(), w); err != nil {
        http.Error(w, "render failed", http.StatusInternalServerError)
        return
    }
}
```

## Level 2: Routing and Request Data

- **Routes:** start with `http.NewServeMux()` and explicit handlers.
- **Methods:** branch on `r.Method` when endpoint serves multiple actions.
- **Input sources:** query (`r.URL.Query()`), form (`r.ParseForm()` + `FormValue`), path segments.
- **Status discipline:** set status before rendering non-200 responses.
- **Post/Redirect/Get:** redirect after successful form POSTs.

## Level 3: Production Patterns

- Add middleware for auth/logging/request IDs.
- Return component-based error pages for consistent UX.
- Serve static files with `http.FileServer` and a dedicated prefix.
- Keep handler functions thin; move business logic into services.

```go
func usersHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        users := listUsers(r.Context())
        _ = components.UserList(users).Render(r.Context(), w)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}
```

## Escalate to Other Skills

- Need reusable view APIs: use `templ-components`.
- Need HTMX partial responses: use `templ-htmx`.
- Need syntax corrections in `.templ`: use `templ-syntax`.

## References

- Handlers and routing: `resources/handlers-and-routing.md`
- Middleware and errors: `resources/middleware-and-errors.md`
- Request data and response patterns: `resources/request-and-response-patterns.md`
- Go `net/http`: https://pkg.go.dev/net/http
- templ server-side rendering: https://templ.guide/server-side-rendering/
