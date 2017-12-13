# manhole

An opinionated HTTP manhole into Go processes.

## Assumptions This Package Makes

- Make each server instance have a unique HTTP port that is randomized by default.
  This makes it very hard to accidentally route this manhole to the outside world.
  If more assurance is required I personally suggest using [yubikey totp][yktotp],
  but do research.
- Application code does not touch [`http.DefaultServeMux`][default-servemux]'. This is so that
  administative control rods can be dynamically flipped in the case they are
  needed.
- [pprof][pprof] endpoints added to `http.DefaultServeMux`. This allows easy
  access to [pprof runtime tracing][pprof-tracing] to debug issues on long-running
  applications like HTTP services.
- Make the manhole slightly inconvenient to put into place in production. This
  helps make sure that this tool remains a debugging tool and not a part of a
  long-term production rollout.

## Usage

Compile this as a plugin:

```console
$ go get -d github.com/Xe/gopreload/manhole
$ go build -buildmode plugin -o manhole.so github.com/Xe/gopreload/manhole
```

Then add [`gopreload`][gopreload] to your application:

```go
// gopreload.go
package main

/*
    This file is separate to make it very easy to both add into an application, but
    also very easy to remove.
*/

import _ "github.com/Xe/gopreload"
```

And at runtime add the `manhole.so` file you created earlier to the target system
somehow and add the following environment variable to its run configuration:

```sh
GO_PRELOAD=/path/to/manhole.so
```

---

[pprof]: https://godoc.org/net/http/pprof
[default-servemux]: https://godoc.org/net/http#pkg-variables
[yktotp]: https://github.com/GeertJohan/yubigo
[gopreload]: https://github.com/Xe/gopreload
