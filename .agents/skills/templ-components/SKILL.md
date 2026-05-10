---
name: templ-components
description: Create reusable templ UI components with props, children, and composition patterns. Use when building UI components, creating component libraries, mentions 'button component', 'card component', or 'reusable templ components'.
---

# Templ Components

Use progressive disclosure: solve with Level 1 first, then pull deeper guidance only if complexity requires it.

## Level 1: Component Checklist

Use this skill for reusable templ UI components.

1. Define a small, single-purpose component.
2. Prefer typed props (structs for complex APIs).
3. Support composition with `{ children... }`.
4. Keep variants explicit and predictable.
5. Extract shared layout wrappers instead of duplicating markup.

```templ
type ButtonProps struct {
    Label    string
    Variant  string
    Disabled bool
}

templ Button(props ButtonProps) {
    <button class={ "btn btn-" + props.Variant } disabled?={ props.Disabled }>
        { props.Label }
    </button>
}
```

## Level 2: API Design Rules

- **Small surface area:** avoid giant prop lists.
- **Typed options:** enums/constants beat free-form strings for critical variants.
- **Composition first:** prefer wrappers (`Card`, `Modal`, `Layout`) over monolith components.
- **No hidden side effects:** components render; handlers/loaders do data work.
- **Stable contracts:** avoid breaking prop shape without migration.

## Level 3: Reusable Patterns

```templ
templ Card(title string) {
    <article class="card">
        <header><h3>{ title }</h3></header>
        <div class="card-body">{ children... }</div>
    </article>
}

templ WithError(err error) {
    if err != nil {
        <p class="error">{ err.Error() }</p>
    } else {
        { children... }
    }
}
```

## Escalate to Other Skills

- Need syntax details: use `templ-syntax`.
- Need HTTP routing/rendering: use `templ-http`.
- Need partial updates/interactions: use `templ-htmx`.

## References

- Foundations: `resources/foundations.md`
- Patterns: `resources/patterns.md`
- Testing and API notes: `resources/testing-and-api.md`
- templ component guide: https://templ.guide/component-composition/
- templ API docs: https://pkg.go.dev/github.com/a-h/templ
