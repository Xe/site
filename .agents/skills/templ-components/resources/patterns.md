# Templ Components Patterns

## Wrapper Pattern

```templ
templ Card(title string) {
    <article class="card">
        <header><h3>{ title }</h3></header>
        <section>{ children... }</section>
    </article>
}
```

## Guard/Decorator Pattern

```templ
templ WithError(err error) {
    if err != nil {
        <div class="alert alert-error">{ err.Error() }</div>
    } else {
        { children... }
    }
}
```

## Layout Composition

- Build pages by nesting small components.
- Prefer `Layout -> Section -> Leaf` composition over one large component.
- Pass render fragments/components when you need slot-like APIs.

## Slot API Choices

- Use `{ children... }` when call-site readability matters most.
- Use explicit `templ.Component` parameters when slots need names (`header`, `footer`, `actions`).
- Use `templ.Join(...)` in Go to stitch optional components in a predictable order.

## View Model Pattern

When domain models do not match display needs, shape a view model first and pass that into components. This keeps templates simple and easier to test.

## Sources

- https://templ.guide/syntax-and-usage/template-composition
- https://templ.guide/core-concepts/view-models
