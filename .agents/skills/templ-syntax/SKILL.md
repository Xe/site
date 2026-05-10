---
name: templ-syntax
description: Learn and write templ component syntax including expressions, conditionals, loops, and Go integration. Use when writing .templ files, learning templ syntax, or mentions 'templ component', 'templ expressions', or '.templ file syntax'.
---

# Templ Syntax

Use this skill with progressive disclosure: start at Level 1, then read deeper sections only when needed.

## Level 1: Fast Path

Use this skill when you need to write or fix `.templ` syntax.

1. Define a component with typed params.
2. Render values with `{ expr }`.
3. Compose components with `@OtherComponent(...)`.
4. Use Go control flow (`if`, `for`, `switch`) directly in markup.
5. Keep business logic in Go helpers, not inline in templates.

```templ
package components

templ Greeting(name string, isAdmin bool) {
    <h1>Hello, { name }</h1>
    if isAdmin {
        <p>Admin mode</p>
    }
}
```

## Level 2: Core Rules

- **Output:** `{ value }` auto-escapes text.
- **Dynamic attrs:** `class={ classes }`, `disabled?={ isDisabled }`.
- **Children:** accept with `{ children... }`, pass with block syntax.
- **Composition:** call components with `@Component(...)`.
- **Safety:** URLs should use safe helpers (for example `templ.URL(...)`) when appropriate.

## Level 3: Common Patterns

```templ
templ Card(title string) {
    <section class="card">
        <h2>{ title }</h2>
        <div>{ children... }</div>
    </section>
}

templ UserList(users []string) {
    if len(users) == 0 {
        <p>No users</p>
    } else {
        <ul>
            for _, user := range users {
                <li>{ user }</li>
            }
        </ul>
    }
}
```

## Escalate to Other Skills

- Building reusable UI APIs: use `templ-components`.
- Wiring templates to handlers/routes: use `templ-http`.
- Adding dynamic interactions: use `templ-htmx`.

## References

- Foundations: `resources/foundations.md`
- Patterns: `resources/patterns.md`
- Reference: `resources/reference.md`
- templ docs: https://templ.guide
- Syntax and expressions: https://templ.guide/syntax-and-usage/
