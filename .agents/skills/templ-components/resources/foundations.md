# Templ Components Foundations

Use this for component API design and composition rules.

templ components compile to Go functions returning `templ.Component`.

## Typed Props

```templ
type ButtonProps struct {
    Text     string
    Variant  string
    Disabled bool
}

templ Button(props ButtonProps) {
    <button class={ "btn btn-" + props.Variant } disabled?={ props.Disabled }>
        { props.Text }
    </button>
}
```

## Core Rules

- Keep components single-purpose.
- Prefer struct props once call sites exceed 2-3 arguments.
- Support composition with `{ children... }` for wrapper components.
- Avoid embedding fetch/mutation logic in templ components.
- Prefer idempotent render behavior: pass data in, render out.
- Export reusable components with capitalized names for cross-package use.

## Code-Only Components

You can implement components in Go with `templ.ComponentFunc`, but then escaping and HTML safety are your responsibility.

## Sources

- https://templ.guide/core-concepts/components
- https://templ.guide/syntax-and-usage/template-composition
