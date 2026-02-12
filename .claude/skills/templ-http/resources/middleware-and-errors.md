# Templ HTTP Middleware and Errors

## Middleware Pattern

```go
func logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}
```

Use middleware for shared concerns: auth, logging, tracing, request IDs.

## Component Error Pages

```go
func renderError(w http.ResponseWriter, r *http.Request, code int, msg string) {
    w.WriteHeader(code)
    _ = components.ErrorPage(code, msg).Render(r.Context(), w)
}
```

## Context and Middleware

- Templates can read request context, so middleware can attach request-scoped values.
- Be careful with type assertions from context values; missing keys can panic.

## CSP Nonce Pattern

When using strict CSP, generate a nonce per request, set the CSP header, and pass it with `templ.WithNonce(ctx, nonce)` before rendering.

## Buffered vs Streaming

- Default templ handler behavior buffers output before writing, which improves status/error handling safety.
- `templ.WithStreaming()` enables progressive output and `@templ.Flush()`, but headers/status are harder to change once bytes are sent.

## Status and Render Discipline

- Set status codes before rendering the response body.
- Return early on handler errors.
- Keep handler bodies thin; call services/helpers for business logic.

## Sources

- https://templ.guide/syntax-and-usage/context
- https://templ.guide/security/content-security-policy
- https://templ.guide/server-side-rendering/streaming
