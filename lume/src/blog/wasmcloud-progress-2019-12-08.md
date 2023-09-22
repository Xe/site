---
title: "Trisiel Progress: Hello, World!"
date: 2019-12-08
series: olin
tags:
  - wasm
  - faas
---

I have been working off and on over the years and have finally created the base
of a functions as a service backend for [WebAssembly][wasm] code. I'm code-naming this
wasmcloud. [Trisiel][wasmcloud] is a pre-alpha prototype and is currently very much work in
progress. However, it's far enough along that I would like to explain what I
have been doing for the last few years and what it's all built up to.

Here is a high level view of all of the parts that make up wasmcloud and how
they correlate:

![wasmcloud graphviz dependency map](/static/blog/wasmcloud-grid.png)

## Land: The Beginning

A little bit after I found WebAssembly I started to play with it. It seemed like
it was too good to be true. A completely free and open source VM format that
would run on almost any platform? Sounds like the kind of black magick
witchcraft you hear about on Star Trek.

However, I kept at it and continued experimenting. I eventually came up with
[Land][land]. This was a very simple thing and was really used to help me invent
Dagger.

Dagger was an attempt at an incredible amount of minimalism. I based it on an
extreme interpretation of the Unix philosophy (everything is a file ->
everything is a bytestream) combined with some Plan 9 for flavor. It had only 5
system calls:

- `open()` - opens a stream by URL, returning a stream descriptor
- `close()` - closes a stream descriptor
- `read()` - reads from a stream
- `write()` - writes to a stream
- `flush()` - flushes intermediate data and turns async behavior into syncronous
  behavior

And yet this was enough to implement a HTTP client.

The core guiding idea was that a cloud-native OS API should expose internet
resources as easily as it exposes native resources. It should be as easy to use
WebSockets as it is to use normal sockets. Additionally, all of the details
should be abstracted away from the WebAssembly module. DNS resolution is not its
job. TLS configuration is not its job. Its job is to run your code. Everything
else should just be provided by the system.

I wrote a
[blogpost](https://xeiaso.net/blog/land-1-syscalls-file-io-2018-06-18)
about this work and even did a
[talk at GoCon
Canada](https://xeiaso.net/talks/webassembly-on-the-server-system-calls-2019-05-31)
about it.

And this worked for several months as I learned WebAssembly and started to
experiment with bigger and better things.

## Olin: Phase 2

Land taught me a lot. I started to quickly run into the limits of Dagger though.
I ended up needing calls like non-cryptographic entropy, environment variables,
command-line arguments and getting the current time. After doing some research
(and trying/failing to implement my own such API based on [newlib][newlib]) I
found a library and specification called [CommonWA][cwa]. This claimed to offer
a lot of what I was looking for. Namely URLs as filenames and all of the host
interop support I could hope for. I named this platform Olin, or the One
Language Intelligent Network.

However the specification was somewhat dead. The author of it had largely moved
on to more ferrous pastures and I became one of the few users of it. I ended up
[forking the specification][olincwa] and implementing my view of what it should
be.

I ended up implementing a [Rust implementation][olincwarust] of the guest ->
host API for the Webassembly side of things. I forked some of the existing Rust
code for this and gradually started adding more and more things. The [test
harness][olincwatest] is the biggest wasm program I've written for a while.
Seriously, there's a lot going on there. It tests every single function exposed
in the CWA spec as well as all of the schemes I had implemented.

Over time I ended up testing Olin in more and more places and on more and more
hardware. As a side effect of all of this being pure go, it was very easy to
cross compile for PowerPC, 32 bit arm (including a $9 arm board that lives under
my desk) and even other targets that gccgo supports. I even ended up porting
[part of TempleOS to Olin][olintempleos] as a proof of concept, but have more
plans in the future for porting other parts of its kernel as a way to help
people understand low-level operating system development.

I've even written a few blogposts about Olin:

- [Olin: Why](https://xeiaso.net/blog/olin-1-why-09-1-2018)
- [Olin: The Future](https://xeiaso.net/blog/olin-2-the-future-09-5-2018)

But, this was great for running stuff interactively and via the command line. It
left me wanting more. I wanted to have that mythical functions as a service
backend that I've been dreaming of. So, I created [Trisiel][wasmcloud].

## h

As an interlude, I also created the [h programming language][hlang] during this
time as a satirical parody of [V][vlang]. This ended up helping me test a lot of
the core functionality that I had built up with Olin. Here's an example of a
program in h:

```
h
```

And this compiles to:

```
(module
 (import "h" "h" (func $h (param i32)))
 (func $h_main
       (local i32 i32 i32)
       (local.set 0 (i32.const 10))
       (local.set 1 (i32.const 104))
       (local.set 2 (i32.const 39))
       (call $h (get_local 1))
       (call $h (get_local 0))
 )
 (export "h" (func $h_main))
)
```

This ends up printing:

```
h
```

I think this is the smallest (if not one of the smallest) quine generator in the
world. I even got this program running on bare metal:

![](/static/blog/xeos_h.png)

[hlang]: https://h.christine.website
[vlang]: https://vlang.io

## Trisiel

[Trisiel][wasmcloud] is the culmination of all of this work. The goal of
wasmcloud is to create a functions as a service backend for running people's
code in an isolated server-side environment.

Users can use the `wasmcloud` command line tool to do everything at the moment:

```
$ wasmcloud
Usage: wasmcloud <flags> <subcommand> <subcommand args>

Subcommands:
        commands         list all command names
        flags            describe all known top-level flags
        help             describe subcommands and their syntax

Subcommands for api:
        login            logs into wasmcloud
        whoami           show information about currently logged in user

Subcommands for handlers:
        create           create a new handler
        logs             shows logs for a handler

Subcommands for utils:
        namegen          show information about currently logged in user
        run              run a webassembly file with the same environment as production servers


Top-level flags (use "wasmcloud flags" for a full list):
  -api-server=http://wasmcloud.kahless.cetacean.club:3002: default API server
  -config=/home/cadey/.wasmc.json: default config location
```

This tool lets you do a few basic things:

- Authenticate with the Trisiel server
- Create handlers from WebAssembly files that meet the CommonWA API as realized
  by Olin
- Get logs for individual handler invocations
- Run WebAssembly modules locally like they would get run on Trisiel

Nearly all of the complexity is abstracted away from users as much as possible.

## Future Steps

In the future I hope to do the following things:

- Support updating handlers to new versions of the code
- Support live-streaming of logs
- Support handler deletion
- Support bulk queue export
- Support [wasi](https://wasi.dev) for easier interoperability
- Support more resource types such as websockets
- Investigate porting the wasmcloud executor to Rust
- Documentation/a book on how to use wasmcloud
- Create an easier way to create accounts that can make handlers
- Deploy to production somewhere

## GReeTZ

Every single one of these people was immeasurably helpful in this research over
the years.

- A. Wilcox
- acln
- as
- bb010g
- [dalias](https://twitter.com/RichFelker)
- [jaddr2line](https://twitter.com/jaddr2line)
- [neelance](https://github.com/neelance)

And many more I can't remember because it's been so many.

---

If you want to support my work, please do so via
[Patreon](https://www.patreon.com/cadey). It really means a lot to me and helps
to keep the dream alive!

[wasm]: https://webassembly.org
[land]: https://tulpa.dev/cadey/land
[newlib]: https://wiki.osdev.org/Porting_Newlib
[cwa]: https://github.com/CommonWA
[olincwa]: https://github.com/Xe/olin/tree/master/docs/cwa-spec
[olincwarust]: https://github.com/Xe/olin/tree/53746b195a6fb302e968d76ffa01b49ad7505330/cwa/olin
[olincwatest]: https://github.com/Xe/olin/blob/53746b195a6fb302e968d76ffa01b49ad7505330/cwa/tests/src/main.rs
[olintempleos]: https://xeiaso.net/blog/templeos-2-god-the-rng-2019-05-30
[wasmcloud]: https://tulpa.dev/within/wasmcloud
