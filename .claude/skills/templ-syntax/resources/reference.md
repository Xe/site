# Templ Syntax Reference

## Quick Syntax Table

| Syntax | Purpose | Example |
| --- | --- | --- |
| `{ expr }` | Output expression | `{ name }` |
| `attr={ val }` | Dynamic attr | `id={ userID }` |
| `attr?={ bool }` | Conditional attr | `disabled?={ busy }` |
| `{ attrs... }` | Spread attrs | `{ attrs... }` |
| `@Component()` | Component call | `@Header()` |
| `{ children... }` | Child slot | `{ children... }` |
| `{{ ... }}` | Raw Go block | `{{ n := len(items) }}` |
| `@templ.Fragment(key)` | Named fragment | `@templ.Fragment("row") { ... }` |
| `templ.WithFragments` | Select output fragments | `ctx = templ.WithFragments(ctx, "row")` |

## Best Practices

1. Keep components small and focused.
2. Move complex branching/formatting into Go helpers.
3. Use strongly typed inputs over ad-hoc maps.
4. Keep rendering deterministic and side-effect free.

## Safety APIs

- `templ.URL(...)`: safe URL type for custom URL attrs (for example `hx-get`).
- `templ.SafeURL(...)`: bypass URL sanitization only for trusted values.
- `templ.Raw(...)`: inject trusted raw HTML only.
- `templ.JSFuncCall(...)`: safe JS call construction for handler attrs.
- `templ.JSONString(...)` and `templ.JSONScript(...)`: safer JSON embedding patterns.

## Cross-Skill Routing

- Reusable component APIs: `templ-components`
- HTTP integration: `templ-http`
- Interactivity with partial swaps: `templ-htmx`

## Sources

- https://templ.guide/syntax-and-usage/attributes
- https://templ.guide/syntax-and-usage/raw-go
- https://templ.guide/syntax-and-usage/fragments
- https://templ.guide/syntax-and-usage/rendering-raw-html
- https://templ.guide/syntax-and-usage/script-templates
