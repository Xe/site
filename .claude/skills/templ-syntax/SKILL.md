---
name: templ-syntax
description: Learn and write templ component syntax including expressions, conditionals, loops, and Go integration. Use when writing .templ files, learning templ syntax, or mentions 'templ component', 'templ expressions', or '.templ file syntax'.
---

# Templ Syntax

## Overview

Templ syntax combines Go code with HTML markup in `.templ` files. Components are compiled to Go functions, giving you type safety and IDE support.

## When to Use This Skill

Use when:

- Writing `.templ` component files
- Learning templ syntax and expressions
- User mentions "templ syntax", "component definition", ".templ file"
- Converting HTML to templ components

## Basic Syntax

### Component Definition

```templ
package components

templ ComponentName(param1 type1, param2 type2) {
    <div>Content</div>
}
```

**Rules:**

- Start with `package` declaration
- Use `templ` keyword for component
- PascalCase for exported components
- camelCase for internal components

### Expressions

Output Go variables with `{ }`:

```templ
templ Greeting(name string, age int) {
    <p>Hello, { name }!</p>
    <p>You are { strconv.Itoa(age) } years old</p>
}
```

**Expression rules:**

- `{ variable }` outputs text (auto-escaped)
- Can call Go functions: `{ strings.ToUpper(name) }`
- Must return string or implement `fmt.Stringer`

### HTML Elements

Write HTML directly:

```templ
templ Card(title string) {
    <div class="card">
        <h2>{ title }</h2>
        <p>Content goes here</p>
    </div>
}
```

### Attributes

Static attributes:

```templ
<div class="container" id="main">
```

Dynamic attributes:

```templ
templ Button(id string, disabled bool) {
    <button
        id={ id }
        disabled?={ disabled }
        class={ getButtonClass() }
    >
        Click
    </button>
}
```

**Attribute syntax:**

- `attr={ value }` - dynamic value
- `attr?={ bool }` - conditional attribute
- Use quotes for static: `class="btn"`

### Conditional Rendering

If/else:

```templ
templ Alert(message string, isError bool) {
    if isError {
        <div class="alert-error">{ message }</div>
    } else {
        <div class="alert-info">{ message }</div>
    }
}
```

Multiple conditions:

```templ
templ Status(code int) {
    if code < 300 {
        <span class="success">Success</span>
    } else if code < 400 {
        <span class="redirect">Redirect</span>
    } else {
        <span class="error">Error</span>
    }
}
```

### Loops

For loop:

```templ
templ List(items []string) {
    <ul>
        for _, item := range items {
            <li>{ item }</li>
        }
    </ul>
}
```

With index:

```templ
templ NumberedList(items []string) {
    <ol>
        for i, item := range items {
            <li>{ strconv.Itoa(i+1) }. { item }</li>
        }
    </ol>
}
```

### Switch Statements

```templ
templ Badge(status string) {
    switch status {
        case "active":
            <span class="badge-green">Active</span>
        case "pending":
            <span class="badge-yellow">Pending</span>
        case "inactive":
            <span class="badge-gray">Inactive</span>
        default:
            <span class="badge-default">Unknown</span>
    }
}
```

### Component Composition

Call other components with `@`:

```templ
templ Layout(title string) {
    <!DOCTYPE html>
    <html>
        <head>
            <title>{ title }</title>
        </head>
        <body>
            @Header()
            <main>
                { children... }
            </main>
            @Footer()
        </body>
    </html>
}

templ Header() {
    <header>
        <h1>My Site</h1>
    </header>
}

templ Footer() {
    <footer>
        <p>&copy; 2024</p>
    </footer>
}
```

**Usage:**

```templ
templ HomePage() {
    @Layout("Home") {
        <p>Welcome to home page</p>
    }
}
```

### Children

Accept child content:

```templ
templ Card(title string) {
    <div class="card">
        <h3>{ title }</h3>
        <div class="card-body">
            { children... }
        </div>
    </div>
}
```

Use:

```templ
templ Profile() {
    @Card("User Profile") {
        <p>Name: John Doe</p>
        <p>Email: john@example.com</p>
    }
}
```

### CSS Blocks

Inline styles:

```templ
templ StyledComponent() {
    <style>
        .custom-class {
            color: blue;
            font-size: 16px;
        }
    </style>
    <div class="custom-class">Styled content</div>
}
```

### Script Blocks

Inline JavaScript:

```templ
templ Interactive() {
    <button id="myBtn">Click me</button>

    <script>
        document.getElementById('myBtn').addEventListener('click', function() {
            alert('Clicked!');
        });
    </script>
}
```

### Comments

Go-style comments:

```templ
package components

// Card component displays a card
templ Card(title string) {
    // Main card container
    <div class="card">
        { /* This is a block comment */ }
        <h3>{ title }</h3>
    </div>
}
```

## Quick Reference

| Syntax            | Purpose               | Example                          |
| ----------------- | --------------------- | -------------------------------- |
| `{ expr }`        | Output expression     | `{ name }`                       |
| `attr={ val }`    | Dynamic attribute     | `id={ userId }`                  |
| `attr?={ bool }`  | Conditional attribute | `disabled?={ isDisabled }`       |
| `@Component()`    | Call component        | `@Header()`                      |
| `{ children... }` | Accept children       | `{ children... }`                |
| `if/else`         | Conditional           | `if isAdmin { }`                 |
| `for range`       | Loop                  | `for _, item := range items { }` |
| `switch`          | Switch statement      | `switch status { case "ok": }`   |

## Common Patterns

### Pattern 1: Conditional Classes

```templ
templ Button(text string, isPrimary bool) {
    <button class={ templ.KV("btn-primary", isPrimary) }>
        { text }
    </button>
}
```

### Pattern 2: List with Empty State

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

### Pattern 3: Data Attributes

```templ
templ DataCard(id string, value int) {
    <div
        data-id={ id }
        data-value={ strconv.Itoa(value) }
    >
        Content
    </div>
}
```

## Best Practices

1. **Keep components small** - One component, one purpose
2. **Use type-safe params** - Leverage Go's type system
3. **Avoid complex logic** - Move to Go functions
4. **Consistent naming** - PascalCase for exports
5. **Escape when needed** - Use `{ }` for auto-escaping

## Next Steps

- **Build components** → Use `templ-components` skill
- **Connect to HTTP** → Use `templ-http` skill
- **Add interactivity** → Use `templ-htmx` skill
