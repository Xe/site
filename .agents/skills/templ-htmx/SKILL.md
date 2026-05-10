---
name: templ-htmx
description: Build interactive hypermedia-driven applications with templ and HTMX. Use when creating dynamic UIs, real-time updates, AJAX interactions, mentions 'HTMX', 'dynamic content', or 'interactive templ app'.
---

# Templ + HTMX Integration

Use progressive disclosure: first make one interaction work, then scale to advanced behaviors.

## Level 1: First Working Flow

Use this skill for server-driven interactivity without a JS framework.

1. Mount HTMX assets in server setup.
2. Include HTMX script in the layout.
3. Add `hx-*` attributes to a component.
4. Return a partial component from the handler.
5. Branch full-page vs fragment responses with HTMX request detection.

```go
import "xeiaso.net/v4/web/htmx"

func main() {
    mux := http.NewServeMux()
    htmx.Mount(mux)
}
```

```templ
import "xeiaso.net/v4/web/htmx"

templ Layout() {
    <html>
        <head>@htmx.Use()</head>
        <body>{ children... }</body>
    </html>
}
```

## Level 2: Core HTMX Controls

- `hx-get` / `hx-post`: trigger requests.
- `hx-target`: pick where response lands.
- `hx-swap`: choose replacement strategy (`innerHTML`, `outerHTML`, `beforeend`).
- `hx-trigger`: control event timing (`click`, `change`, `every 5s`, etc).
- `hx-indicator`: show loading state.

## Level 3: Advanced Server Patterns

- Detect HTMX requests with `htmx.Is(r)` and return fragments.
- Use out-of-band updates for multi-region refreshes.
- Use response headers (`HX-Trigger`, `HX-Redirect`) for client behavior.
- Preserve progressive enhancement: endpoints should still work without JS.

```go
func profileHandler(w http.ResponseWriter, r *http.Request) {
    if htmx.Is(r) {
        _ = components.ProfilePanel().Render(r.Context(), w)
        return
    }
    _ = components.ProfilePage().Render(r.Context(), w)
}
```

## Escalate to Other Skills

- Need handler/routing structure: use `templ-http`.
- Need reusable component APIs: use `templ-components`.
- Need template syntax help: use `templ-syntax`.

## References

- Quick start: `resources/quick-start.md`
- Interaction patterns: `resources/interaction-patterns.md`
- Advanced responses: `resources/advanced-responses.md`
- HTMX docs: https://htmx.org/docs/
- Hypermedia Systems: https://hypermedia.systems/
