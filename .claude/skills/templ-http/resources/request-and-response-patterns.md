# Templ HTTP Request and Response Patterns

## Request Data Sources

- Query params: `r.URL.Query().Get("q")`
- Form data: `r.ParseForm(); r.FormValue("email")`
- Path segments: trim or route extraction helpers

## Post/Redirect/Get

```go
func contactHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        _ = components.ContactForm().Render(r.Context(), w)
        return
    }
    _ = r.ParseForm()
    // process form
    http.Redirect(w, r, "/contact/success", http.StatusSeeOther)
}
```

## Form Handling Flow

Typical templ form flow:

1. `ParseForm` and decode into a request model.
2. Validate.
3. On error, re-render with user input + validation messages.
4. On success, perform action and redirect.

Use view models to keep templates focused on display state instead of transport/domain objects.

## Partial Responses

For targeted updates, render only named fragments with `templ.WithFragments(...)`.

## Mixed HTML/JSON Endpoint

- If `Accept: application/json`, encode JSON and return.
- Otherwise render templ component for HTML.

## Static Files

- Serve with `http.FileServer` behind `/static/` prefix.

## Sources

- https://templ.guide/syntax-and-usage/forms
- https://templ.guide/core-concepts/view-models
- https://templ.guide/syntax-and-usage/fragments
