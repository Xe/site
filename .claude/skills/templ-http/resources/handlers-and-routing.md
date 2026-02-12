# Templ HTTP Handlers and Routing

Use `templ.Handler(...)` for simple fixed-render endpoints, and custom handlers with `.Render(...)` when request data drives output.

## templ.Handler Shortcut

```go
http.Handle("/", templ.Handler(
    components.HomePage(),
    templ.WithStatus(http.StatusOK),
    templ.WithContentType("text/html; charset=utf-8"),
))
```

Common options include `templ.WithStatus`, `templ.WithContentType`, and `templ.WithErrorHandler`.

## Minimal Handler

```go
func homeHandler(w http.ResponseWriter, r *http.Request) {
    if err := components.HomePage().Render(r.Context(), w); err != nil {
        http.Error(w, "render failed", http.StatusInternalServerError)
        return
    }
}
```

## Routing Setup

```go
func routes() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/", homeHandler)
    mux.HandleFunc("/users", usersHandler)
    return mux
}
```

## Method Dispatch

- Branch on `r.Method` where one path supports multiple actions.
- Return `http.StatusMethodNotAllowed` for unsupported methods.
- Keep route registration explicit for readability.

## Fragment Responses

For partial UI updates, set `ctx := templ.WithFragments(r.Context(), "fragment-name")` before rendering.

## Sources

- https://templ.guide/server-side-rendering/creating-an-http-server-with-templ
- https://templ.guide/syntax-and-usage/fragments
