---
name: templ-components
description: Create reusable templ UI components with props, children, and composition patterns. Use when building UI components, creating component libraries, mentions 'button component', 'card component', or 'reusable templ components'.
---

# Templ Components

## Overview

Build reusable, type-safe UI components with templ. Components accept strongly-typed props and can be composed together to create complex UIs.

## When to Use This Skill

Use when:

- Creating reusable UI components
- Building component libraries
- User mentions "button", "card", "form", "modal" components
- Designing component APIs
- Working on design systems

## Core Component Patterns

### Basic Component

```templ
package components

templ Button(text string) {
    <button class="btn">
        { text }
    </button>
}
```

### Component with Multiple Props

```templ
templ Button(text string, variant string, disabled bool) {
    <button
        class={ "btn btn-" + variant }
        disabled?={ disabled }
    >
        { text }
    </button>
}
```

Usage:

```go
components.Button("Submit", "primary", false)
```

### Component with Children

```templ
templ Card(title string) {
    <div class="card">
        <div class="card-header">
            <h3>{ title }</h3>
        </div>
        <div class="card-body">
            { children... }
        </div>
    </div>
}
```

Usage:

```templ
@Card("User Profile") {
    <p>Name: John Doe</p>
    <p>Email: john@example.com</p>
}
```

### Component with Struct Props

```templ
package components

type ButtonProps struct {
    Text     string
    Variant  string
    Disabled bool
    OnClick  string
}

templ Button(props ButtonProps) {
    <button
        class={ "btn btn-" + props.Variant }
        disabled?={ props.Disabled }
        onclick={ props.OnClick }
    >
        { props.Text }
    </button>
}
```

## Common Components

### Button Variants

```templ
templ PrimaryButton(text string) {
    <button class="btn btn-primary">{ text }</button>
}

templ SecondaryButton(text string) {
    <button class="btn btn-secondary">{ text }</button>
}

templ DangerButton(text string, onclick string) {
    <button class="btn btn-danger" onclick={ onclick }>
        { text }
    </button>
}
```

### Card Component

```templ
templ Card(title string, footer string) {
    <div class="card">
        if title != "" {
            <div class="card-header">
                <h3>{ title }</h3>
            </div>
        }
        <div class="card-body">
            { children... }
        </div>
        if footer != "" {
            <div class="card-footer">
                { footer }
            </div>
        }
    </div>
}
```

### List Component

```templ
type ListItem struct {
    ID   string
    Text string
}

templ List(items []ListItem) {
    <ul class="list">
        for _, item := range items {
            <li data-id={ item.ID }>
                { item.Text }
            </li>
        }
    </ul>
}
```

### Modal Component

```templ
templ Modal(id string, title string, isOpen bool) {
    <div
        class={ "modal", templ.KV("modal-open", isOpen) }
        id={ id }
    >
        <div class="modal-backdrop"></div>
        <div class="modal-content">
            <div class="modal-header">
                <h2>{ title }</h2>
                <button class="modal-close">&times;</button>
            </div>
            <div class="modal-body">
                { children... }
            </div>
        </div>
    </div>
}
```

### Form Components

```templ
templ Input(name string, label string, value string) {
    <div class="form-group">
        <label for={ name }>{ label }</label>
        <input
            type="text"
            id={ name }
            name={ name }
            value={ value }
            class="form-control"
        />
    </div>
}

templ TextArea(name string, label string, rows int) {
    <div class="form-group">
        <label for={ name }>{ label }</label>
        <textarea
            id={ name }
            name={ name }
            rows={ strconv.Itoa(rows) }
            class="form-control"
        >
            { children... }
        </textarea>
    </div>
}

templ Select(name string, label string, options []string) {
    <div class="form-group">
        <label for={ name }>{ label }</label>
        <select id={ name } name={ name } class="form-control">
            for _, option := range options {
                <option value={ option }>{ option }</option>
            }
        </select>
    </div>
}
```

## Layout Components

### Container

```templ
templ Container(fluid bool) {
    <div class={ templ.KV("container", !fluid), templ.KV("container-fluid", fluid) }>
        { children... }
    </div>
}
```

### Grid

```templ
templ Row() {
    <div class="row">
        { children... }
    </div>
}

templ Col(size int) {
    <div class={ "col-" + strconv.Itoa(size) }>
        { children... }
    </div>
}
```

Usage:

```templ
@Container(false) {
    @Row() {
        @Col(6) {
            <p>Left column</p>
        }
        @Col(6) {
            <p>Right column</p>
        }
    }
}
```

### Navigation

```templ
type NavItem struct {
    Text   string
    Href   string
    Active bool
}

templ Nav(items []NavItem) {
    <nav class="navbar">
        <ul class="nav">
            for _, item := range items {
                <li class={ templ.KV("nav-item-active", item.Active) }>
                    <a href={ templ.URL(item.Href) }>
                        { item.Text }
                    </a>
                </li>
            }
        </ul>
    </nav>
}
```

## Composition Patterns

### Slots Pattern

```templ
templ Layout(title string, headerContent templ.Component, footerContent templ.Component) {
    <!DOCTYPE html>
    <html>
        <head>
            <title>{ title }</title>
        </head>
        <body>
            <header>
                @headerContent
            </header>
            <main>
                { children... }
            </main>
            <footer>
                @footerContent
            </footer>
        </body>
    </html>
}
```

### Render Props Pattern

```templ
templ DataTable(headers []string, renderRow func(int) templ.Component) {
    <table>
        <thead>
            <tr>
                for _, header := range headers {
                    <th>{ header }</th>
                }
            </tr>
        </thead>
        <tbody>
            for i := 0; i < 10; i++ {
                @renderRow(i)
            }
        </tbody>
    </table>
}
```

### Wrapper Components

```templ
templ WithLoading(isLoading bool) {
    if isLoading {
        <div class="spinner">Loading...</div>
    } else {
        { children... }
    }
}

templ WithError(err error) {
    if err != nil {
        <div class="alert alert-error">
            { err.Error() }
        </div>
    } else {
        { children... }
    }
}
```

## Best Practices

### 1. Single Responsibility

```templ
// ✅ Good: One purpose
templ Avatar(src string, alt string) {
    <img src={ src } alt={ alt } class="avatar" />
}

// ❌ Bad: Too many responsibilities
templ UserSection(user User, posts []Post, comments []Comment) {
    // Too much in one component
}
```

### 2. Type-Safe Props

```templ
// ✅ Good: Strongly typed
type ButtonProps struct {
    Text    string
    Variant ButtonVariant
}

type ButtonVariant string

const (
    Primary   ButtonVariant = "primary"
    Secondary ButtonVariant = "secondary"
)

templ Button(props ButtonProps) {
    <button class={ "btn btn-" + string(props.Variant) }>
        { props.Text }
    </button>
}
```

### 3. Composition Over Complexity

```templ
// ✅ Good: Compose small components
templ UserCard(user User) {
    @Card(user.Name) {
        @Avatar(user.AvatarURL, user.Name)
        @UserInfo(user)
        @UserActions(user.ID)
    }
}

// Each sub-component is simple and reusable
```

### 4. Conditional Rendering

```templ
// ✅ Good: Clear conditions
templ Message(text string, isError bool) {
    if isError {
        <div class="alert-error">{ text }</div>
    } else {
        <div class="alert-info">{ text }</div>
    }
}
```

### 5. Default Props

```go
// In Go code
type ButtonProps struct {
    Text     string
    Variant  string
    Disabled bool
}

func NewButton(text string) ButtonProps {
    return ButtonProps{
        Text:     text,
        Variant:  "primary",
        Disabled: false,
    }
}
```

```templ
templ Button(props ButtonProps) {
    <button
        class={ "btn btn-" + props.Variant }
        disabled?={ props.Disabled }
    >
        { props.Text }
    </button>
}
```

### Package Organization

```templ
// components/ui/button.templ
package ui

templ Button(text string) {
    <button class="btn">{ text }</button>
}
```

```templ
// components/layout/container.templ
package layout

templ Container() {
    <div class="container">
        { children... }
    </div>
}
```

Usage:

```go
import (
    "myapp/components/ui"
    "myapp/components/layout"
)

layout.Container() {
    ui.Button("Click me")
}
```

## Testing Components

```go
func TestButton(t *testing.T) {
    var buf bytes.Buffer
    props := ButtonProps{
        Text:    "Submit",
        Variant: "primary",
    }

    err := Button(props).Render(context.Background(), &buf)
    if err != nil {
        t.Fatal(err)
    }

    html := buf.String()
    if !strings.Contains(html, "Submit") {
        t.Error("Button text not found")
    }
    if !strings.Contains(html, "btn-primary") {
        t.Error("Primary class not found")
    }
}
```

## Next Steps

- **Connect to server** → Use `templ-http` skill
- **Add interactivity** → Use `templ-htmx` skill
- **Style components** → Use `templ-css` skill
