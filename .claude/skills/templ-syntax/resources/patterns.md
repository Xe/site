# Templ Syntax Patterns

## Composition and Children

```templ
templ Layout(title string) {
    <html>
        <head><title>{ title }</title></head>
        <body>
            @Header()
            <main>{ children... }</main>
            @Footer()
        </body>
    </html>
}
```

- Call components with `@Component(...)`.
- Accept nested content with `{ children... }`.

Code-side composition helpers:

- `templ.WithChildren(ctx, child)` injects children for downstream rendering.
- `templ.GetChildren(ctx)` reads current children.
- `templ.ClearChildren(ctx)` prevents accidental child forwarding.

## Component Parameters and Join

```templ
templ Panel(header templ.Component) {
    <section>
        @header
        <div>{ children... }</div>
    </section>
}
```

- Use `templ.Component` params for explicit slots.
- Use `templ.Join(...)` in Go for rendering multiple components in order.

## Conditional List Rendering

```templ
templ ItemList(items []string) {
    if len(items) == 0 {
        <p>No items found</p>
    } else {
        <ul>
            for _, item := range items {
                <li>{ item }</li>
            }
        </ul>
    }
}
```

## Conditional Class Pattern

```templ
templ Button(text string, isPrimary bool) {
    <button class={ templ.KV("btn-primary", isPrimary) }>{ text }</button>
}
```

## Fragment Pattern

```templ
templ UserPage(user User) {
    @templ.Fragment("user-card") {
        <article>{ user.Name }</article>
    }
}
```

- Render selected fragments with `templ.WithFragments(ctx, "user-card")`.
- Use `templ.RenderFragments` when writing outside normal HTTP handlers.
- Fragment filtering trims output only; non-fragment template logic still executes.

## Sources

- https://templ.guide/syntax-and-usage/template-composition
- https://templ.guide/syntax-and-usage/fragments
- https://templ.guide/syntax-and-usage/attributes
