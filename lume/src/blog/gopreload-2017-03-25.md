---
title: "gopreload: LD_PRELOAD for the Gopher crowd"
date: "2017-03-25"
---

A common pattern in Go libraries is to take advantage of [init functions][initf]
to do things like settings up defaults in loggers, automatic metrics instrumentation,
flag values, [debugging tools][manhole] or database drivers. With monorepo culture
prevalent in larger microservices based projects, this can lead to a few easily
preventable problems:

- Forgetting to set up a logger default or metrics submission, making operations
  teams blind to the performance of the app and developer teams blind to errors
  that come up during execution.
- The requirement to make code changes to add things like metrics or HTTP routing
  extensions.

There is an environment variable in Linux libc's called `LD_PRELOAD` that will
load arbitrary shared objects into ram before anything else is started. This
has been used for [good][good-ld-preload] and [evil][evil-ld-preload], but the
behavior is the same basic idea as [underscore imports in Go][underscore-import].

My solution for this is [gopreload][gopreload]. It emulates the behavior of
`LD_PRELOAD` but with [Go plugins][go-plugins]. This allows users to explicitly
automatically load arbitrary Go code into ram while the process starts.

## Usage

To use this, add `gopreload` to your application's imports:

```go
// gopreload.go
package main

/*
    This file is separate to make it very easy to both add into an application, but
    also very easy to remove.
*/

import _ "github.com/Xe/gopreload"
```

and then compile `manhole`:

```console
$ go get -d github.com/Xe/gopreload/manhole
$ go build -buildmode plugin -o $GOPATH/manhole.so github.com/Xe/gopreload/manhole
```

then run your program with `GO_PRELOAD` set to the path of `manhole.so`:

```console
$ export GO_PRELOAD=$GOPATH/manhole.so
$ go run *.go
2017/03/25 10:56:22 gopreload: trying to open: /home/xena/go/manhole.so
2017/03/25 10:56:22 manhole: Now listening on http://127.0.0.2:37588
```

That endpoint has pprof and a few other [fun tools][manhole-tools] set up, making
it a good stopgap "manhole" into the performance of a service.

## Security Implications

This package assumes that programs run using it are never started with environment
variables that are set by unauthenticated users. Any errors in loading the plugins
will be logged using the standard library logger `log` and ignored.

This has about the same security implications as [`LD_PRELOAD`][ld-preload] does in most
Linux distributions, but the risk is minimal compared to the massive benefit for
being able to have arbitrary background services all be able to be dug into using
the same tooling or being able to have metric submission be completely separated
from the backend metric creation. Common logging setup processes can be _always_
loaded, making the default logger settings into the correct settings.

## Feedback

To give feedback about gopreload, please contact me on [twitter][twitter-addr] or
on the Gophers slack (I'm `@xena` there). For issues with gopreload please file
[an issue on Github][gopreload-issues].

[initf]: https://go.dev/doc/effective_go#init
[manhole]: https://github.com/Xe/gopreload/tree/master/manhole
[good-ld-preload]: http://www.logix.cz/michal/devel/faketime/
[evil-ld-preload]: https://rafalcieslak.wordpress.com/2013/04/02/dynamic-linker-tricks-using-ld_preload-to-cheat-inject-features-and-investigate-programs/
[underscore-import]: https://go.dev/doc/effective_go#blank
[gopreload]: https://github.com/Xe/gopreload
[go-plugins]: https://pkg.go.dev/plugin
[manhole-tools]: https://github.com/Xe/gopreload/blob/master/manhole/server.go
[ld-preload]: https://rafalcieslak.wordpress.com/2013/04/02/dynamic-linker-tricks-using-ld_preload-to-cheat-inject-features-and-investigate-programs/
[twitter-addr]: https://twitter.com/theprincessxena
[gopreload-issues]: https://github.com/Xe/gopreload/issues/new
