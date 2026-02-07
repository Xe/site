---
name: templ-htmx
description: Build interactive hypermedia-driven applications with templ and HTMX. Use when creating dynamic UIs, real-time updates, AJAX interactions, mentions 'HTMX', 'dynamic content', or 'interactive templ app'.
---

# Templ + HTMX Integration

## Overview

HTMX enables modern, interactive web applications with minimal JavaScript. Combined with templ's type-safe components, you get fast, reliable hypermedia-driven UIs.

**Key Benefits:**

- No JavaScript framework needed
- Server-side rendering
- Minimal client-side code
- Progressive enhancement
- Type-safe components

## When to Use This Skill

Use when:

- Building interactive UIs
- Creating dynamic content
- User mentions "HTMX", "dynamic updates", "real-time"
- Implementing AJAX-like behavior without JS
- Building SPAs without frameworks

## Quick Start

### 1. Import and Mount HTMX

First, import the htmx package and mount it in your server:

```go
import (
    "xeiaso.net/v4/web/htmx"
)

func main() {
    mux := http.NewServeMux()

    // Mount HTMX static files at /.within.website/x/htmx/
    htmx.Mount(mux)

    // ... other routes
}
```

### 2. Add HTMX to Layout

```templ
package components

import "xeiaso.net/v4/web/htmx"

templ Layout(title string) {
    <!DOCTYPE html>
    <html>
        <head>
            <title>{ title }</title>
            @htmx.Use()
        </head>
        <body>
            { children... }
        </body>
    </html>
}
```

The `htmx.Use()` component includes the core HTMX library. You can add extensions:

```templ
// Add SSE and path-params extensions
@htmx.Use("sse", "path-params")
```

Available extensions: `event-header`, `path-params`, `remove-me`, `websocket`.

### 3. Detect HTMX Requests

Use the `htmx.Is()` function to check if a request was made by HTMX:

```go
import "xeiaso.net/v4/web/htmx"

func handler(w http.ResponseWriter, r *http.Request) {
    if htmx.Is(r) {
        // Return partial HTML fragment for HTMX
        components.Partial().Render(r.Context(), w)
    } else {
        // Return full page for direct navigation
        components.FullPage().Render(r.Context(), w)
    }
}
```

### 4. Create Interactive Component

```templ
templ Counter(count int) {
    <div>
        <p>Count: { strconv.Itoa(count) }</p>
        <button
            hx-post="/counter/increment"
            hx-target="#counter"
            hx-swap="outerHTML"
        >
            Increment
        </button>
    </div>
}
```

### 5. Create Handler

```go
func incrementHandler(w http.ResponseWriter, r *http.Request) {
    count := getCount() + 1
    saveCount(count)

    components.Counter(count).Render(r.Context(), w)
}
```

## Core HTMX Attributes

### hx-get / hx-post

Trigger HTTP requests:

```templ
templ SearchBox() {
    <input
        type="text"
        name="q"
        hx-get="/search"
        hx-trigger="keyup changed delay:500ms"
        hx-target="#results"
    />
    <div id="results"></div>
}
```

Handler:

```go
func searchHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    results := search(query)

    components.SearchResults(results).Render(r.Context(), w)
}
```

### hx-target

Specify where to insert response:

```templ
templ LoadMore(page int) {
    <button
        hx-get={ "/posts?page=" + strconv.Itoa(page) }
        hx-target="#posts"
        hx-swap="beforeend"
    >
        Load More
    </button>
}
```

### hx-swap

Control how content is swapped:

```templ
// innerHTML (default)
hx-swap="innerHTML"

// outerHTML - replace element itself
hx-swap="outerHTML"

// beforeend - append inside
hx-swap="beforeend"

// afterend - insert after
hx-swap="afterend"
```

### hx-trigger

Control when requests fire:

```templ
// On click (default for buttons)
<button hx-get="/data">Click me</button>

// On change
<select hx-get="/filter" hx-trigger="change">

// On keyup with delay
<input hx-get="/search" hx-trigger="keyup changed delay:300ms">

// On page load
<div hx-get="/data" hx-trigger="load">

// Every 5 seconds
<div hx-get="/updates" hx-trigger="every 5s">
```

## Common Patterns

### Pattern 1: Live Search

Component:

```templ
templ SearchBox() {
    <div>
        <input
            type="text"
            name="q"
            placeholder="Search..."
            hx-get="/search"
            hx-trigger="keyup changed delay:500ms"
            hx-target="#search-results"
            hx-indicator="#spinner"
        />
        <span id="spinner" class="htmx-indicator">
            Searching...
        </span>
    </div>
    <div id="search-results"></div>
}

templ SearchResults(results []string) {
    <ul>
        for _, result := range results {
            <li>{ result }</li>
        }
    </ul>
}
```

Handler:

```go
func searchHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    results := performSearch(query)

    components.SearchResults(results).Render(r.Context(), w)
}
```

### Pattern 2: Infinite Scroll

```templ
templ PostList(posts []Post, page int) {
    <div id="posts">
        for _, post := range posts {
            @PostCard(post)
        }
    </div>

    if len(posts) > 0 {
        <div
            hx-get={ "/posts?page=" + strconv.Itoa(page+1) }
            hx-trigger="revealed"
            hx-swap="outerHTML"
        >
            Loading more...
        </div>
    }
}
```

### Pattern 3: Delete with Confirmation

```templ
templ DeleteButton(itemID string) {
    <button
        hx-delete={ "/items/" + itemID }
        hx-confirm="Are you sure?"
        hx-target="closest tr"
        hx-swap="outerHTML swap:1s"
    >
        Delete
    </button>
}
```

Handler:

```go
func deleteHandler(w http.ResponseWriter, r *http.Request) {
    itemID := strings.TrimPrefix(r.URL.Path, "/items/")
    deleteItem(itemID)

    // Return empty to remove element
    w.WriteHeader(http.StatusOK)
}
```

### Pattern 4: Inline Edit

```templ
templ EditableField(id string, value string) {
    <div id={ "field-" + id }>
        <span>{ value }</span>
        <button
            hx-get={ "/edit/" + id }
            hx-target={ "#field-" + id }
            hx-swap="outerHTML"
        >
            Edit
        </button>
    </div>
}

templ EditForm(id string, value string) {
    <form
        hx-post={ "/save/" + id }
        hx-target={ "#field-" + id }
        hx-swap="outerHTML"
    >
        <input type="text" name="value" value={ value } />
        <button type="submit">Save</button>
        <button
            hx-get={ "/cancel/" + id }
            hx-target={ "#field-" + id }
        >
            Cancel
        </button>
    </form>
}
```

### Pattern 5: Form Validation

```templ
templ SignupForm() {
    <form hx-post="/signup" hx-target="#form-errors">
        <div id="form-errors"></div>

        <input
            type="email"
            name="email"
            hx-post="/validate/email"
            hx-trigger="blur"
            hx-target="#email-error"
        />
        <div id="email-error"></div>

        <input type="password" name="password" />

        <button type="submit">Sign Up</button>
    </form>
}

templ ValidationError(message string) {
    <span class="error">{ message }</span>
}
```

### Pattern 6: Polling / Real-time Updates

```templ
templ LiveStats() {
    <div
        hx-get="/stats"
        hx-trigger="load, every 5s"
        hx-swap="innerHTML"
    >
        Loading stats...
    </div>
}

templ StatsDisplay(stats Stats) {
    <div>
        <p>Users online: { strconv.Itoa(stats.UsersOnline) }</p>
        <p>Active sessions: { strconv.Itoa(stats.Sessions) }</p>
    </div>
}
```

## Advanced Patterns

### Out-of-Band Updates (OOB)

Update multiple parts of page:

```templ
templ CartButton(count int) {
    <button id="cart-btn">
        Cart ({ strconv.Itoa(count) })
    </button>
}

templ AddToCartResponse(item Item) {
    // Main response
    <div class="notification">
        Added { item.Name } to cart!
    </div>

    // Update cart button (different part of page)
    <div id="cart-btn" hx-swap-oob="true">
        @CartButton(getCartCount())
    </div>
}
```

### Progressive Enhancement

```templ
templ Form() {
    <form
        action="/submit"
        method="POST"
        hx-post="/submit"
        hx-target="#result"
    >
        <input type="text" name="data" />
        <button type="submit">Submit</button>
    </form>
    <div id="result"></div>
}
```

Works without JavaScript, enhanced with HTMX.

### Loading States

```templ
templ DataTable() {
    <div
        hx-get="/data"
        hx-trigger="load"
        hx-indicator="#loading"
    >
        <div id="loading" class="htmx-indicator">
            Loading data...
        </div>
    </div>
}
```

CSS:

```css
.htmx-indicator {
  display: none;
}

.htmx-request .htmx-indicator {
  display: inline;
}

.htmx-request.htmx-indicator {
  display: inline;
}
```

## Response Headers

### HX-Trigger

Trigger client-side events:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Do work...

    // Trigger custom event
    w.Header().Set("HX-Trigger", "itemCreated")

    components.Success().Render(r.Context(), w)
}
```

Client side:

```javascript
document.body.addEventListener("itemCreated", function (evt) {
  console.log("Item created!");
});
```

### HX-Redirect

```go
func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("HX-Redirect", "/dashboard")
    w.WriteHeader(http.StatusOK)
}
```

### HX-Refresh

```go
func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("HX-Refresh", "true")
    w.WriteHeader(http.StatusOK)
}
```

### Special Status Code

Stop polling with status code 286:

```go
import "within.website/x/htmx"

func pollHandler(w http.ResponseWriter, r *http.Request) {
    if shouldStopPolling() {
        w.WriteHeader(htmx.StatusStopPolling)
        return
    }
    // ... return normal content
}
```

### Request Headers

Access HTMX request headers:

```go
import "within.website/x/htmx"

func handler(w http.ResponseWriter, r *http.Request) {
    // Check if this is an HTMX request
    if htmx.Is(r) {
        // Get user's response to hx-prompt
        promptResponse := r.Header.Get(htmx.HeaderPrompt)
    }
}
```

## Best Practices

1. **Keep handlers focused** - Return only the HTML fragment needed
2. **Use semantic HTML** - Works without JS
3. **Handle errors gracefully** - Return error components
4. **Optimize responses** - Send minimal HTML
5. **Use OOB for multi-updates** - Update multiple page sections
6. **Progressive enhancement** - Always provide fallback

## Full Example: Todo App

```templ
// components/todo.templ
package components

type Todo struct {
    ID        string
    Text      string
    Completed bool
}

templ TodoApp(todos []Todo) {
    @Layout("Todo App") {
        <div>
            <h1>My Todos</h1>

            @TodoForm()
            @TodoList(todos)
        </div>
    }
}

templ TodoForm() {
    <form
        hx-post="/todos"
        hx-target="#todo-list"
        hx-swap="beforeend"
        hx-on::after-request="this.reset()"
    >
        <input
            type="text"
            name="text"
            placeholder="New todo..."
            required
        />
        <button type="submit">Add</button>
    </form>
}

templ TodoList(todos []Todo) {
    <ul id="todo-list">
        for _, todo := range todos {
            @TodoItem(todo)
        }
    </ul>
}

templ TodoItem(todo Todo) {
    <li id={ "todo-" + todo.ID }>
        <input
            type="checkbox"
            checked?={ todo.Completed }
            hx-post={ "/todos/" + todo.ID + "/toggle" }
            hx-target={ "#todo-" + todo.ID }
            hx-swap="outerHTML"
        />
        <span class={ templ.KV("completed", todo.Completed) }>
            { todo.Text }
        </span>
        <button
            hx-delete={ "/todos/" + todo.ID }
            hx-target={ "#todo-" + todo.ID }
            hx-swap="outerHTML swap:500ms"
        >
            Delete
        </button>
    </li>
}
```

Handlers:

```go
func todosHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        todos := getAllTodos()
        components.TodoApp(todos).Render(r.Context(), w)

    case "POST":
        r.ParseForm()
        todo := createTodo(r.FormValue("text"))
        components.TodoItem(todo).Render(r.Context(), w)
    }
}

func todoToggleHandler(w http.ResponseWriter, r *http.Request) {
    id := extractID(r.URL.Path)
    todo := toggleTodo(id)
    components.TodoItem(todo).Render(r.Context(), w)
}

func todoDeleteHandler(w http.ResponseWriter, r *http.Request) {
    id := extractID(r.URL.Path)
    deleteTodo(id)
    w.WriteHeader(http.StatusOK) // Empty response removes element
}
```

## Resources

- [HTMX Documentation](https://htmx.org/docs/)
- [HTMX Examples](https://htmx.org/examples/)
- [Hypermedia Systems Book](https://hypermedia.systems/)
