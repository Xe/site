# Templ + HTMX Interaction Patterns

Project standard: use `import "xeiaso.net/v4/web/htmx"` with `@htmx.Use()` in layout, rather than manually adding script tags.

## Core Attributes

- `hx-get`, `hx-post`, `hx-delete`
- `hx-target`
- `hx-swap` (`innerHTML`, `outerHTML`, `beforeend`)
- `hx-trigger` (`click`, `blur`, `every 5s`, etc)
- `hx-indicator`

Use `templ.URL(...)` for URL-valued custom attrs such as `hx-get` and `hx-post` when values are dynamic.

## Common Patterns

- Live search with delayed keyup trigger.
- Infinite scroll with `hx-trigger="revealed"`.
- Inline edit by swapping display row with form row.
- Form validation via field-level HTMX requests.

## Boosted Navigation and Forms

- `hx-boost="true"` upgrades regular links/forms to HTMX behavior.
- Keep standard links/forms semantics so direct navigation still works.

## Event Handlers in Attributes

- For `hx-on:*` or `on*`, use static handler strings for simple cases.
- Use `templ.JSFuncCall(...)` for safer dynamic JS argument encoding.

## Progressive Enhancement

Always include normal `action`/`method` on forms and server endpoints that work without JS.

## Sources

- https://templ.guide/server-side-rendering/htmx
- https://templ.guide/syntax-and-usage/attributes
- https://templ.guide/syntax-and-usage/forms
- https://templ.guide/syntax-and-usage/script-templates
