gopreload
=========

An emulation of the linux libc `LD_PRELOAD` except for use with Go plugins for
the addition of instrumentation and debugging utilities.

## Pluginizer

`pluginizer` is a bit of glue that makes it easier to turn underscore imports
into plugins:

```console
$ go get github.com/Xe/gopreload/cmd/pluginizer
$ pluginizer -help
Usage of pluginizer:
  -dest string
        destination package to generate
  -pkg string
        package to underscore import
$ pluginizer -pkg github.com/lib/pq -dest github.com/Xe/gopreload/database/postgres
To build this plugin:
  $ go build -buildmode plugin -o /path/to/output.so github.com/Xe/gopreload/database/postgres
```

### Database drivers

I have included plugin boilerplate autogenned versions of the sqlite, postgres
and mysql database drivers.

## Manhole

[`manhole`][manhole] is an example of debugging and introspection tooling that has
been useful when debugging past issues with long-running server processes.

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

---

[manhole]: https://github.com/Xe/gopreload/tree/master/manhole
[ld-preload]: https://rafalcieslak.wordpress.com/2013/04/02/dynamic-linker-tricks-using-ld_preload-to-cheat-inject-features-and-investigate-programs/
