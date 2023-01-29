---
title: "ln - The Natural Log Function"
date: 2020-10-17
tags:
 - golang
 - go
---

One of the most essential things in software is a good interface for logging
data to places. Logging is a surprisingly hard problem and there are many
approaches to doing it. This time, we're going to talk about my favorite logging
library in Go that uses my favorite function I've ever written in Go.

Today we're talking about [ln](https://github.com/Xe/ln), the natural log
function. ln works with key value pairs and logs them to somewhere. By default
it logs things to standard out. Here is how you use it:

```go
package main

import (
  "context"

  "within.website/ln"
)

func main() {
  ctx := context.Background()
  ln.Log(ctx, ln.Fmt("hello %s", "world"), ln.F{"demo": "usage"})
}
```

ln works with key value pairs called
[F](https://pkg.go.dev/within.website/ln#F). This type allows you to
log just about _anything_ you want, including custom data types with an
[Fer](https://pkg.go.dev/within.website/ln#Fer). This will let
you annotate your data types so that you can automatically extract the important
information into your logs while automatically filtering out passwords or other
secret data. Here's an example:

```go
type User struct {
  ID       int
  Username string
  Password []byte
}

func (u User) F() ln.F {
	return ln.F{
		"user_id":       u.ID,
		"user_name": u.Username,
	}
}
```

Then if you create that user somehow, you can log the ID and username without
logging the password on accident:

```go
var theDude User = abides()

ln.Log(ctx, ln.Info("created new user"), theDude)
```

This will create a log line that looks something like this:

```
level=info msg="created new user" user_name="The Dude" user_id=1337
```

[You can also put values in contexts! See <a
href="https://github.com/Xe/ln/blob/master/ex/http.go#L21">here</a> for more
detail on how this works.](conversation://Mara/hacker)

The way this is all glued together is that F itself is an Fer, meaning that the
Log/Error functions take a variadic set of Fers. This is where my favorite Go
function comes into play, it is the implementation of the Fer interface for F.
Here is that function verbatim:

```go
// F makes F an Fer
func (f F) F() F {
	return f
}
```

I love how this function looks like some kind of abstract art. This function
holds this library together.

If you end up using ln for your projects in the future, please let me know what
your experience is like. I would love to make this library the best it can
possibly be. It is not a nanosecond scale zero allocation library (I think those
kind of things are a bit of a waste of time, because most of the time your
logging library is NOT going to be your bottleneck), but it is designed to have
very usable defaults and solve the problem good enough that you shouldn't need
to care. There are a few useful tools in the
[ex](https://pkg.go.dev/within.website/ln/ex) package nested in ln. The biggest
thing is the HTTP middleware, which has saved me a lot of effort when writing
web services in Go.
