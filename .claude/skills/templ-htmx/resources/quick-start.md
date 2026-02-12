# Templ + HTMX Quick Start

## Required Integration Pattern

This project requires:

- `import "xeiaso.net/v4/web/htmx"`
- `htmx.Mount(mux)` in server setup
- `@htmx.Use()` in the layout `<head>`
- `htmx.Is(r)` checks in handlers to branch full-page vs fragment responses

## Server Setup

```go
import "xeiaso.net/v4/web/htmx"

func main() {
    mux := http.NewServeMux()
    htmx.Mount(mux)
    // register app routes
}
```

## Layout Setup

```templ
import "xeiaso.net/v4/web/htmx"

templ Layout() {
    <html>
        <head>
            @htmx.Use()
        </head>
        <body>{ children... }</body>
    </html>
}
```

## First Interactive Action

- Add `hx-post` or `hx-get` to a button/input.
- Use `hx-target` and `hx-swap` to control replacement.
- Use `hx-select` when only part of response HTML should be swapped.
- Keep normal form `action`/`method` for no-JS fallback.

## Request Detection With `htmx.Is`

```go
import "xeiaso.net/v4/web/htmx"

func profileHandler(w http.ResponseWriter, r *http.Request) {
    if htmx.Is(r) {
        _ = components.ProfilePanel().Render(r.Context(), w)
        return
    }
    _ = components.ProfilePage().Render(r.Context(), w)
}
```

## Fragment Optimization

When templ responses are large, combine HTMX with templ fragments (`templ.WithFragments`) to return only the needed section.

## Sources

- https://templ.guide/server-side-rendering/htmx
- https://templ.guide/syntax-and-usage/fragments
