# Templ + HTMX Advanced Responses

## Fragment-Based Partial Responses

Use `htmx.Is(r)` to decide whether to return a full page or an HTMX-targeted partial.

```go
func profileHandler(w http.ResponseWriter, r *http.Request) {
    ctx := templ.WithFragments(r.Context(), "profile-panel")
    _ = components.ProfilePage().Render(ctx, w)
}
```

- Define fragment blocks in templates with `@templ.Fragment("profile-panel") { ... }`.
- Nested fragments can be selectively rendered.
- Important: fragment filtering limits output, but non-fragment logic still runs.

## Rendering Fragments Outside HTTP

- Use `templ.RenderFragments(...)` when writing to buffers/files.

## Streaming for Progressive Delivery

- `templ.WithStreaming()` enables progressive output.
- `@templ.Flush()` can improve perceived latency for long renders.
- Tradeoff: after bytes are written, header/status changes are constrained.

## External HTMX Concepts

- Headers like `HX-Trigger`, `HX-Redirect`, and `HX-Refresh` are HTMX features, but they are not documented as templ-specific patterns in templ.guide.

## Sources

- https://templ.guide/syntax-and-usage/fragments
- https://templ.guide/server-side-rendering/streaming
