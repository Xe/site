---
title: Olin Improvements
date: 2019-12-14
series: olin
tags:
 - wasm
 - wasmcloud
 - rust
 - zig
 - cgi
---

Over the last week or so I've been doing a _lot_ of improvements to [Olin][olin] in order to make it ready to be the kernel for the minimum viable product of [wasmcloud][wasmcloud-hello-world]. Here's an overview of the big things that have happened from version [0.1.1][olin-0.1.1] to version [0.4.0][olin-0.4.0].

[olin]: https://github.com/Xe/olin
[wasmcloud-hello-world]: https://xeiaso.net/blog/wasmcloud-progress-2019-12-08
[olin-0.1.1]: https://github.com/Xe/olin/releases/tag/v0.1.1
[olin-0.4.0]: https://github.com/Xe/olin/releases/tag/v0.4.0

## What is Olin?

[Olin][olin] is a userspace kernel designed for multi-tenant secure computing. It provides isolation via WebAssembly to limit the attack scope of malicious user input, resource accounting via its runtime statistics, and a familiar Unix-like API. It is the core that you can build a functions as a service platform on top of.

[olin]: https://github.com/Xe/olin

As Olin is just a kernel, it needs some work in order to really shine as a true child of the cloud. That work is incoming during the next weeks and months.

## Announcements

Here is what has been done since the [last Olin post][last-olin-post]:

[last-olin-post]: https://xeiaso.net/blog/olin-2-the-future-09-5-2018

* An official, automated build of the example Olin components has been published to the Docker Hub
* The Go ABI has been deprecated for the moment
* The entrypoint of Olin programs has changed to _start
* The beginning of support in the Zig standard library
* Official binfmt_misc rigging has been created for experimentation

### Official Docker Hub Build

The Docker Hub repo [xena/olin][docker-hub] now is automatically built off of the latest master release of Olin.

[docker-hub]: https://hub.docker.com/r/xena/olin

To use this image, run the following commands:

```console
$ docker pull xena/olin:latest
$ docker run --rm -it xena/olin:latest sh
```

Then you can use the `cwa` tool to run programs in `/wasm`. See `cwa -help` for more information.

### Deprecation of Go Support

For the moment, I am deprecating support for [Go][golang] in `GOOS=js GOARCH=wasm`. The ABI for the Go compiler in this mode is too unstable for me right now. If other people want to fix [`abi/wasmgo`][abi-wasmgo] to support Go 1.13 and newer, I would be _very_ welcome to the patches.

[golang]: https://go.dev/
[abi-wasmgo]: https://github.com/Xe/olin/tree/master/abi/wasmgo

### The Entrypoint is Now `_start()`

Early on in the experiments that make up Olin, I have made a mistake in my fundamental understanding of how operating systems run programs. I thought that the main function would return the exit code of the program. This is not the case. There is a small shim that wraps the main function of your language and passes the result of it to `exit()`. Olin now copies this behavior. In order to return a value to the Olin runtime, you can either call `runtime_exit()` or return from the `_start()` function to exit with 0. Many thanks to Andrew Kelly for helping me realize this error.

This behavior is copied in the [Olin rust package][olin-rust-entrypoint], which now has a fancy macro to automate the creation of the `_start()` function.

[olin-rust-entrypoint]: https://github.com/Xe/olin/blob/ffc4ec5d436b6536d8b3917990ac6c53650f4297/rust/olin/src/lib.rs#L424

I am waiting on Zig to release a new nightly version in order to enable it, but the [bring-your-own-OS package][bring-your-own-os] support in Zig means that the Zig standard library is starting to be exposed into Olin programs. Here's an example based on the [example program][zig-example-program]:

[bring-your-own-os]: https://github.com/ziglang/zig/commit/b375f6e027a159616e80906aa05e253fbe8cc9df
[zig-example-program]: https://github.com/ziglang/zig/blob/b375f6e027a159616e80906aa05e253fbe8cc9df/lib/std/special/init-exe/src/main.zig

```zig
pub const os = @import("./olin/olin.zig");
const std = @import("std");

pub fn main() anyerror!void {
    std.debug.warn("All your base are belong to us.\n", .{});
}
```

### `binfmt_misc` Rigging

For a while I've had a [binfmt_misc][binfmt-misc] configuration floating around in the Olin repo. Here's how to use [it][olin-binfmt]:

[binfmt-misc]: https://en.wikipedia.org/wiki/Binfmt_misc
[olin-binfmt]: https://github.com/Xe/olin/blob/master/run/binfmt_misc/cwa.cfg

First, install Olin's `cmd/cwa` to `/usr/local/bin`:

```console
$ cd cmd/cwa
$ go build
$ sudo mv cwa /usr/local/bin
```

Then activate the binfmt_misc configuration:

```console
$ cd ../../run/binfmt_misc
$ cat cwa.cfg | sudo tee /proc/sys/fs/binfmt_misc/register
```

Then you can run Olin programs without calling `cwa`:

```console
$ ./olinfetch.wasm
```

## Features

### Policy Support

Olin now has a declarative policy engine for accessing external resources. This is inspired from OpenBSD `pledge()` and macOS sandboxing. These policies allow setting the following attributes:

* Resources an Olin program can access, matched by regular expressions
* Resources an Olin program CANNOT access, matched by regular expressions 
* The maximum amount of memory an Olin program can use
* The maximum number of WebAssembly instructions a WebAssembly program can execute

Here's an example policy file intended to help with relaying webhooks:

```
## This is an example policy, the ## signifies this line is a comment.

## These are the URL patterns that this handler can open:
allow (
  ^https://tulpa.dev
  ^https://discordapp.com/api/webhooks/
  ^random://$
)

## These are the URL patterns that this handler cannot open:
disallow (
  ^https://tulpa.dev/admin.*$
)

## This is the ram limit in pages (64k each):
ram-page-limit 128

## This is the gas limit in instructions:
gas-limit 1048576
```

This would allow a WebAssembly program to open a HTTP socket to https://tulpa.dev (my git server) and Discord, but disallows any administrative API calls to my git server. It also allows the Olin program to use up to 128 pages of memory (about 8MB, which goes surprisingly far) and 1.04 million instructions. If the handler tries to open any resource that is not explicitly allowed, it is killed. If the handler tries to open a resource that is explicitly forbidden, it is killed. If the handler uses too much ram or too many instructions, it is killed.

This allows handlers to safely process user controlled input and even use that as part of the call to the open function.

When policies are violated, the error thrown is a [vibe check failure][vibe-check]:

[vibe-check]: https://www.urbandictionary.com/define.php?term=Vibe%20Check

```console
$ cwa -policy ../policy/testdata/gitea.policy httptest.wasm
httptest.wasm: 2019/12/13 13:16:15 info: making request 
  to https://xena.greedo.xeserv.us/files/hello_olin.txt
httptest.wasm: 2019/12/13 13:16:15 vibe check failed: 
  https://xena.greedo.xeserv.us/files/hello_olin.txt 
  forbidden by policy
2019/12/13 13:16:15 httptest.wasm: exit status -1
```

### `runtime_exit()` System Call

Along with making `_start()` the entrypoint, there comes a new problem: exiting. I fixed this by adding a [`runtime_exit()`][runtime-exit] system call in Olin. When you call this function with the status code you want to return, execution of the Olin program instantly ends, uncleanly stopping everything and closing all files the program has open. This is similar to Linux's `exit()` system call.

[runtime-exit]: https://github.com/Xe/olin/commit/0036ee8620abe8a25b24c5b7feb76caefba35a8f

It's probably best to save this call for cases where the program _really can't/shouldn't_ continue executing, like for panic handlers.

### Generic CGI Support

Previously there was a half-baked idea I called cwagi in Olin's codebase. The idea was to emulate part of how CGI worked in order to let Olin programs handle HTTP easily. I realize this was a mistake, so now it [just supports normal CGI][cgi-patch], conforming to [RFC 3875][rfc3875].

[cgi-patch]: https://github.com/Xe/olin/commit/92e703fcb2571e1f32e0bf1ba4f17bb45c1d6408
[rfc3875]: https://tools.ietf.org/html/rfc3875

### End of File Error

One of the more common errors in operating systems is the "end of file" error. It is raised when a file has no more data in it. Olin didn't have this, but [now it does][eof-patch].

[eof-patch]: https://github.com/Xe/olin/commit/beb19fd9c6ee2de11f61c6f93fc1813f5f317aff

## Wasmcloud Features

Thanks to these improvements, the following wasmcloud features have been implemented:

* Updating handlers (`wasmcloud update`) [link](https://tulpa.dev/within/wasmcloud/issues/11)
* Deleting handlers (`wasmcloud delete`) [link](https://tulpa.dev/within/wasmcloud/issues/21)
* Listing deleted handlers (`wasmcloud list -show-deleted`)
* Brand outgoing HTTP requests [link](https://tulpa.dev/within/wasmcloud/commit/3024971fdf3d437a2bda95206f7fb123be1a8df5)

As of the writing of this post, wasmcloud is currently 75% through the [MVP development cycle][mvp-milestone]. Here are the remaining issues:

[mvp-milestone]: https://tulpa.dev/within/wasmcloud/milestone/1

* CGI support for handlers [link](https://tulpa.dev/within/wasmcloud/issues/16)
* Policy support for handlers [link](https://tulpa.dev/within/wasmcloud/issues/12)
* Configuration variables for handlers [link](https://tulpa.dev/within/wasmcloud/issues/6)

---

Overall, this project is fun. Here's to 1.0 happening soon! Be well.
