# Templ Syntax Foundations

Use this when you need baseline templ syntax behavior.

## Component Shape

```templ
package components

templ Greeting(name string, age int) {
    <p>Hello, { name }!</p>
    <p>You are { strconv.Itoa(age) } years old</p>
}
```

- Start with `package` and imports, then `templ Name(args...)` blocks.
- Components compile to Go functions returning `templ.Component`.
- Keep component inputs typed and narrow.

## Elements and Expressions

- Close tags in templ source (`</tag>` or `/>`), even for void elements.
- `{ expr }` renders escaped output with context-aware escaping.
- `attr={ value }` sets dynamic attrs; `attr?={ bool }` toggles boolean attrs.

## Control Flow and Raw Go

- Use normal Go flow directly: `if`, `switch`, `for`.
- Use raw Go blocks with `{{ ... }}` for local variables or pre-compute steps.
- Parser gotcha: raw text starting with `if`, `for`, or `switch` may parse as statements; rewrite text or wrap it in an expression.

Keep heavy logic in Go helpers/services, not inline template logic.

## Sources

- https://templ.guide/syntax-and-usage/basic-syntax
- https://templ.guide/syntax-and-usage/elements
- https://templ.guide/syntax-and-usage/expressions
- https://templ.guide/syntax-and-usage/statements
- https://templ.guide/syntax-and-usage/raw-go
