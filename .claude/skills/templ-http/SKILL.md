---
name: templ-http
description: Integrate templ components with Go HTTP server using net/http. Use when connecting templ to web server, creating HTTP handlers, mentions 'templ server', 'HTTP routes', or 'serve templ components'.
---

# Templ HTTP Integration

## Overview

Connect templ components to Go's `net/http` server. Render components in HTTP handlers and serve dynamic HTML pages.

## When to Use This Skill

Use when:

- Setting up HTTP server with templ
- Creating route handlers
- User mentions "serve templ", "HTTP server", "web server"
- Connecting components to routes
- Rendering templ in handlers

## Basic Integration

### Simple Handler

```go
package main

import (
    "net/http"
    "myapp/components"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
    components.HomePage().Render(r.Context(), w)
}

func main() {
    http.HandleFunc("/", homeHandler)
    http.ListenAndServe(":8080", nil)
}
```

### Handler with Data

```go
func userHandler(w http.ResponseWriter, r *http.Request) {
    user := getUserFromDB(r.URL.Query().Get("id"))

    components.UserProfile(user).Render(r.Context(), w)
}
```

## Rendering Patterns

### Pattern 1: Direct Render

```go
func handler(w http.ResponseWriter, r *http.Request) {
    component := components.Page("Title")
    component.Render(r.Context(), w)
}
```

### Pattern 2: Error Handling

```go
func handler(w http.ResponseWriter, r *http.Request) {
    err := components.Page("Title").Render(r.Context(), w)
    if err != nil {
        http.Error(w, "Render failed", http.StatusInternalServerError)
        log.Printf("Render error: %v", err)
    }
}
```

### Pattern 3: With Layout

```go
func pageHandler(w http.ResponseWriter, r *http.Request) {
    content := components.PageContent()

    components.Layout("Page Title", content).Render(r.Context(), w)
}
```

## Routing

### ServeMux

```go
func main() {
    mux := http.NewServeMux()

    // Static pages
    mux.HandleFunc("/", homeHandler)
    mux.HandleFunc("/about", aboutHandler)

    // Dynamic routes
    mux.HandleFunc("/user/", userHandler)
    mux.HandleFunc("/post/", postHandler)

    // Static files
    fs := http.FileServer(http.Dir("static"))
    mux.Handle("/static/", http.StripPrefix("/static/", fs))

    http.ListenAndServe(":8080", mux)
}
```

### RESTful Routes

```go
func usersHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        users := getAllUsers()
        components.UserList(users).Render(r.Context(), w)

    case "POST":
        // Handle create
        user := createUser(r)
        components.UserCard(user).Render(r.Context(), w)

    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}
```

## Request Data

### Query Parameters

```go
func searchHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    page := r.URL.Query().Get("page")

    results := search(query, page)

    components.SearchResults(query, results).Render(r.Context(), w)
}
```

### Form Data

```go
func loginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        components.LoginForm().Render(r.Context(), w)
        return
    }

    // POST
    r.ParseForm()
    email := r.FormValue("email")
    password := r.FormValue("password")

    if authenticate(email, password) {
        http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
    } else {
        components.LoginForm("Invalid credentials").Render(r.Context(), w)
    }
}
```

### Path Parameters

```go
// /user/123
func userHandler(w http.ResponseWriter, r *http.Request) {
    // Extract ID from path
    path := strings.TrimPrefix(r.URL.Path, "/user/")
    userID := path

    user := getUserByID(userID)
    if user == nil {
        http.NotFound(w, r)
        return
    }

    components.UserProfile(user).Render(r.Context(), w)
}
```

## Middleware

### Logging Middleware

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", homeHandler)

    http.ListenAndServe(":8080", loggingMiddleware(mux))
}
```

### Auth Middleware

```go
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        session := getSession(r)
        if session == nil {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }

        next(w, r)
    }
}

// Usage
http.HandleFunc("/dashboard", authMiddleware(dashboardHandler))
```

## Error Handling

### Custom Error Pages

```go
func errorHandler(w http.ResponseWriter, r *http.Request, status int, message string) {
    w.WriteHeader(status)
    components.ErrorPage(status, message).Render(r.Context(), w)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
    user, err := getUserByID(r.URL.Query().Get("id"))
    if err != nil {
        errorHandler(w, r, 500, "Failed to load user")
        return
    }
    if user == nil {
        errorHandler(w, r, 404, "User not found")
        return
    }

    components.UserProfile(user).Render(r.Context(), w)
}
```

### Error Component

```templ
// components/error.templ
package components

templ ErrorPage(code int, message string) {
    @Layout("Error") {
        <div class="error-page">
            <h1>{ strconv.Itoa(code) }</h1>
            <p>{ message }</p>
            <a href="/">Go Home</a>
        </div>
    }
}
```

## Static Files

```go
func main() {
    mux := http.NewServeMux()

    // Serve static files
    fs := http.FileServer(http.Dir("static"))
    mux.Handle("/static/", http.StripPrefix("/static/", fs))

    // Routes
    mux.HandleFunc("/", homeHandler)

    http.ListenAndServe(":8080", mux)
}
```

## Context Usage

### Passing Data via Context

```go
func contextMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := context.WithValue(r.Context(), "userID", "123")
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func handler(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(string)

    components.Page(userID).Render(r.Context(), w)
}
```

## Full Example

```go
package main

import (
    "log"
    "net/http"
    "myapp/components"
)

func main() {
    mux := http.NewServeMux()

    // Static files
    fs := http.FileServer(http.Dir("static"))
    mux.Handle("/static/", http.StripPrefix("/static/", fs))

    // Routes
    mux.HandleFunc("/", homeHandler)
    mux.HandleFunc("/about", aboutHandler)
    mux.HandleFunc("/contact", contactHandler)

    // Start server
    addr := ":8080"
    log.Printf("Server starting on %s", addr)
    if err := http.ListenAndServe(addr, mux); err != nil {
        log.Fatal(err)
    }
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    components.HomePage().Render(r.Context(), w)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
    components.AboutPage().Render(r.Context(), w)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        components.ContactForm().Render(r.Context(), w)
        return
    }

    // POST
    r.ParseForm()
    name := r.FormValue("name")
    email := r.FormValue("email")
    message := r.FormValue("message")

    // Send email...

    components.ContactSuccess(name).Render(r.Context(), w)
}
```

## Best Practices

1. **Always use r.Context()** when rendering
2. **Handle errors** from Render()
3. **Set appropriate status codes** before rendering
4. **Use middleware** for common functionality
5. **Separate routes** from handler logic
6. **Return early** on errors

## Common Patterns

### Pattern: Redirect After Post

```go
func formHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        components.Form().Render(r.Context(), w)
        return
    }

    // POST - process form
    processForm(r)

    // Redirect
    http.Redirect(w, r, "/success", http.StatusSeeOther)
}
```

### Pattern: JSON API + HTML

```go
func usersHandler(w http.ResponseWriter, r *http.Request) {
    users := getUsers()

    // Check Accept header
    if r.Header.Get("Accept") == "application/json" {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users)
        return
    }

    // Default: HTML
    components.UserList(users).Render(r.Context(), w)
}
```

## Next Steps

- **Add interactivity** → Use `templ-htmx` skill
