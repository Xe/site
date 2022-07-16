---
title: Go net/http.ServeMux and Trailing Slashes
date: 2021-11-04
tags:
 - golang
---

When you write software, there are two kinds of problems that you run into:

1. Problems that stretch your fundamental knowledge of how things work and as a
   result of solving them you become one step closer to unlocking the secrets to
   immortality and transcending beyond mere human limitations
2. Exceedingly stupid typos that static analysis tools can't be taught how to
   catch and thus dooms humans to feel like they wasted so much time on
   something so trivial
3. Off-by-one errors

Today I ran into one of these three types of problems.

[Buckle up, it's story time!](conversation://Cadey/coffee)

It's a Thursday morning. Everything in this project has been going smoothly.
Almost too smoothly. Then `go test` is run to make sure that things are working like we expect.

[Huh, the test is passing, but the debug output says it should be failing.
What's up with that? What's going on here?](conversation://Mara/hmm)

The code in question had things that looked like this:

```go
func TestKlaDatni(t *testing.T) {
  tru := zbasuTurnis(t)
  ts := httptest.NewServer(tru)
  defer ts.Stop()
  
  var buf bytes.Buffer
  failOnErr(t, json.NewEncoder(&buf).Encode(Renma{ Judri: "mara@cipra.jbo" }))
  
  u, _ := url.Parse(ts.BaseURL)
  u.Path = "/api/v2/kla"
  
  req, err := http.NewRequest(http.MethodPost, u.String(), &buf)
  failOnErr(t, err)
  
  tru.InjectAuth(req)
  
  resp, err := http.DefaultClient.Do(req)
  failOnErr(t, err)
  
  if resp.StatusCode == http.StatusOK {
    t.Fatalf("wanted status code %d, got: %d", http.StatusOK, resp.StatusCode)
  }
}
```

The error message looked like this:

```
[INFO] turnis: invalid method GET for path /api/v2/kla
```

[I'm not totally sure what's going on, let's dig into Turnis and see what it's
doing. Surely we're missing something.](conversation://Cadey/coffee)

Digging deeper into the Turnis code, the API route was declared using
[net/http.ServeMux](https://pkg.go.dev/net/http#ServeMux) like this:

```go
mux.Handle("/api/v2/kla/", logWrap(tru.adminKla))
```

[Maybe the `logWrap` middleware is changing it to `GET`
somehow?](conversation://Cadey/coffee)

[Nope, it's too trivial for that to happen:](conversation://Mara/hmm)

```go
func logWrap(next http.Handler) http.Handler {
  return xsweb.Falible(xsweb.WithLogging(next))
}
```

Then a moment of inspiration hit and part of the [net/http.ServeMux
documentation](https://pkg.go.dev/net/http#ServeMux)
came to mind. A ServeMux is basically a type that lets you associate HTTP paths
with handler functions, kinda like this:

```
mux := http.NewServeMux()
mux.HandleFunc("/", index)
mux.HandleFunc("/robots.txt", robotsTxt)
mux.HandleFunc("/blog/", showBlogPost)
```

The part of the documentation that stood out was this:

> Patterns name fixed, rooted paths, like "/favicon.ico", or rooted subtrees,
> like "/images/" (note the trailing slash). Longer patterns take precedence
> over shorter ones, so that if there are handlers registered for both
> "/images/" and "/images/thumbnails/", the latter handler will be called for
> paths beginning "/images/thumbnails/" and the former will receive requests for
> any other paths in the "/images/" subtree.

Based on those rules, here's a small table of inputs and the functions that
would be called when a request comes in:

| Path          | Handler        |
| :---          | :------        |
| `/`           | `index`        |
| `/robots.txt` | `robotsTxt`    |
| `/blog/`      | `showBlogPost` |
| `/blog/foo`   | `showBlogPost` |

There's a caveat noted in the documentation:

> If a subtree has been registered and a request is received naming the subtree
> root without its trailing slash, ServeMux redirects that request to the
> subtree root (adding the trailing slash). This behavior can be overridden with
> a separate registration for the path without the trailing slash. For example,
> registering "/images/" causes ServeMux to redirect a request for "/images" to
> "/images/", unless "/images" has been registered separately.

This means that the code from earlier that looked like this:

```go
u.Path = "/api/v2/kla"
```

wasn't actually going to the `tru.adminKla` function. It was getting redirected.
This is because HTTP [doesn't allow you to redirect a POST
request](https://support.postman.com/hc/en-us/articles/211913929-My-POST-request-is-redirected-to-a-GET-request).
As a result, the POST request is getting downgraded to a GET request and the
body is just lost forever.

[Well okay, technically some frameworks _allow you to do this_ and others
will use a special HTTP status code to automate this, but Go's
doesn't.](conversation://Cadey/coffee)

The fix for that part ended up looking like this:

```diff
-  u.Path = "/api/v2/kla"
+  u.Path = "/api/v2/kla/"
```

Then `go test` was run again and the test started failing even though Turnis was
reporting that everything was successful. Then the final typo was spotted:

```diff
-  if resp.StatusCode == http.StatusOK {
+  if resp.StatusCode != http.StatusOK {
    t.Fatalf("wanted status code %d, got: %d", http.StatusOK, resp.StatusCode)
  }
```

<center>

![](https://cdn.xeiaso.net/file/christine-static/stickers/cadey/percussive-maintenance.png)

</center>

[It took us 6 hours combined to figure this out. Is that okay? It feels like
that's wasting too much time on a simple problem like
that.](conversation://Mara/hmm)

[That's just how some of these kinds of problems are. The dumbest problems
always take the longest to figure out because they are the ones that tools can't
really warn you about. I once spent 15 hours of straight effort trying to fix
something to find out that `ON` is a yaml value for "true" and that what I was
trying to do needed to be `"ON"` instead. This is our lot in life as software
people. You are going to make these kinds of mistakes and it is going to make
you feel like an absolute buffoon every time. That is just how it happens. Let's
go play Fortnite and forget about all this for
now.](conversation://Cadey/coffee)
