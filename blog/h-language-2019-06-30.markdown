---
title: The h Programming Language
date: 2019-06-30
series: h
tags:
 - wasm
 - release
---

[h](https://h.christine.website) is a project of mine that I have released
recently. It is a single-paradigm, multi-tenant friendly, turing-incomplete
programming language that does nothing but print one of two things:

- the letter h
- a single quote (the Lojbanic "h")

It does this via [WebAssembly](https://webassembly.org). This may sound like a
pointless complication, but actually this ends up making things _a lot simpler_.
WebAssembly is a virtual machine (fake computer that only exists in code) intended
for browsers, but I've been using it for server-side tasks.

I have written more about/with WebAssembly in the past in these posts:

- https://xeiaso.net/talks/webassembly-on-the-server-system-calls-2019-05-31
- https://xeiaso.net/blog/olin-1-why-09-1-2018
- https://xeiaso.net/blog/olin-2-the-future-09-5-2018
- https://xeiaso.net/blog/land-1-syscalls-file-io-2018-06-18
- https://xeiaso.net/blog/templeos-2-god-the-rng-2019-05-30

This is a continuation of the following two posts:

- https://xeiaso.net/blog/the-origin-of-h-2015-12-14
- https://xeiaso.net/blog/formal-grammar-of-h-2019-05-19

All of the relevant code for h is [here](https://github.com/Xe/x/tree/v1.1.7/cmd/h).

h is a somewhat standard three-phase compiler. Each of the phases is as follows:

## Parsing the Grammar

As mentioned in a prior post, h has a formal grammar defined in [Parsing Expression Grammar](https://en.wikipedia.org/wiki/Parsing_expression_grammar).
I took this [grammar](https://github.com/Xe/x/blob/v1.1.7/h/h.peg) (with some
minor modifications) and fed it into a tool called [peggy](https://github.com/eaburns/peggy)
to generate a Go source [version of the parser](https://github.com/Xe/x/blob/v1.1.7/h/h_gen.go).
This parser has some minimal [wrappers](https://github.com/Xe/x/blob/v1.1.7/h/parser.go)
around it, mostly to simplify the output and remove unneeded nodes from the tree.
This simplifies the later compilation phases.

The input to h looks something like this:

```
h
```

The output syntax tree pretty-prints to something like this:

```
H("h")
```

This is also represented using a tree of nodes that looks something like this:

```
&peg.Node{
    Name: "H",
    Text: "h",
    Kids: nil,
}
```

A more complicated program will look something like this:

```
&peg.Node{
    Name: "H",
    Text: "h h h",
    Kids: {
        &peg.Node{
            Name: "",
            Text: "h",
            Kids: nil,
        },
        &peg.Node{
            Name: "",
            Text: "h",
            Kids: nil,
        },
        &peg.Node{
            Name: "",
            Text: "h",
            Kids: nil,
        },
    },
}
```

Now that we have this syntax tree, it's easy to go to the next phase of
compilation: generating the WebAssembly Text Format.

## WebAssembly Text Format

[WebAssembly Text Format](https://developer.mozilla.org/en-US/docs/WebAssembly/Understanding_the_text_format)
is a human-editable and understandable version of WebAssembly. It is pretty low
level, but it is actually fairly simple. Let's take an example of the h compiler
output and break it down:

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

Fundamentally, WebAssembly binary files are also called modules. Each .wasm file
can have only one module defined in it. Modules can have sections that contain the
following information:

- External function imports
- Function definitions
- Memory information
- Named function exports
- Global variable definitions
- Other custom data that may be vendor-specific

h only uses external function imports, function definitions and named function
exports.

`import` imports a function from the surrounding runtime with two fields: module
and function name. Because this is an obfuscated language, the function `h` from
module `h` is imported as `$h`. This function works somewhat like the C library
function [putchar()](https://www.tutorialspoint.com/c_standard_library/c_function_putchar.htm).

`func` creates a function. In this case we are creating a function named `$h_main`.
This will be the entrypoint for the h program.

Inside the function `$h_main`, there are three local variables created: `0`, `1` and `2`.
They correlate to the following values:

| Local Number | Explanation       | Integer Value |
| :----------- | :---------------- | :------------ |
|            0 | Newline character |            10 |
|            1 | Lowercase h       |           104 |
|            2 | Single quote      |            39 |

As such, this program prints a single lowercase h and then a newline.

`export` lets consumers of this WebAssembly module get a name for a function,
linear memory or global value. As we only need one function in this module,
we export `$h_main` as `"h"`.

## Compiling this to a Binary

The next phase of compiling is to turn this WebAssembly Text Format into a binary.
For simplicity, the tool `wat2wasm` from the [WebAssembly Binary Toolkit](https://github.com/WebAssembly/wabt)
is used. This tool creates a WebAssembly binary out of WebAssembly Text Format.

Usage is simple (assuming you have the WebAssembly Text Format file above saved as `h.wat`):

```
wat2wasm h.wat -o h.wasm
```

And you will create `h.wasm` with the following sha256 sum:

```
sha256sum h.wasm
8457720ae0dd2deee38761a9d7b305eabe30cba731b1148a5bbc5399bf82401a  h.wasm
```

Now that the final binary is created, we can move to the runtime phase.

## Runtime

The h [runtime](https://github.com/Xe/x/blob/v1.1.7/cmd/h/run.go) is incredibly
simple. It provides the `h.h` putchar-like function and executes the `h`
function from the binary you feed it. It also times execution as well as keeps
track of the number of instructions the program runs. This is called "gas" for
historical reasons involving [blockchains](https://blockgeeks.com/guides/ethereum-gas/).

I use [Perlin Network's life](https://github.com/perlin-network/life) as the
implementation of WebAssembly in h. I have experience with it from [Olin](https://github.com/Xe/olin).

## The Playground

As part of this project, I wanted to create an [interactive playground](https://h.christine.website/play).
This allows users to run arbitrary h programs on my server. As the only system
call is putchar, this is safe. The playground also has some limitations on how
big of a program it can run. The playground server works like this:

- The user program is sent over HTTP with Content-Type [text/plain](https://github.com/Xe/x/blob/v1.1.7/cmd/h/http.go#L402-L413)
- The program is [limited to 75 bytes on the server](https://github.com/Xe/x/blob/v1.1.7/cmd/h/http.go#L44) (though this is [configurable](https://github.com/Xe/x/blob/v1.1.7/cmd/h/http.go#L15) via flags or envvars)
- The program is [compiled](https://github.com/Xe/x/blob/v1.1.7/cmd/h/http.go#L53)
- The program is [run](https://github.com/Xe/x/blob/v1.1.7/cmd/h/http.go#L59)
- The output is [returned via JSON](https://github.com/Xe/x/blob/v1.1.7/cmd/h/http.go#L65-L72)
- This output is then put [into the playground page with JavaScript](https://github.com/Xe/x/blob/v1.1.7/cmd/h/http.go#L389-L394)

The output of this call looks something like this:

```
curl -H "Content-Type: text/plain" --data "h" https://h.christine.website/api/playground | jq
{
  "prog": {
    "src": "h",
    "wat": "(module\n (import \"h\" \"h\" (func $h (param i32)))\n (func $h_main\n       (local i32 i32 i32)\n       (local.set 0 (i32.const 10))\n       (local.set 1 (i32.const 104))\n       (local.set 2 (i32.const 39))\n       (call $h (get_local 1))\n       (call $h (get_local 0))\n )\n (export \"h\" (func $h_main))\n)",
    "bin": "AGFzbQEAAAABCAJgAX8AYAAAAgcBAWgBaAAAAwIBAQcFAQFoAAEKGwEZAQN/QQohAEHoACEBQSchAiABEAAgABAACw==",
    "ast": "H(\"h\")"
  },
  "res": {
    "out": "h\n",
    "gas": 11,
    "exec_duration": 12345
  }
}
```

The execution duration is in [nanoseconds](https://pkg.go.dev/time#Duration), as
it is just directly a Go standard library time duration.

## Bugs h has Found

This will be updated in the future, but h has already found a bug in [Innative](https://innative.dev).
There was a bug in how Innative handled C name mangling of binaries. Output of
the h compiler is now [a test case in Innative](https://github.com/innative-sdk/innative/commit/6353d59d611164ce38b938840dd4f3f1ea894e1b#diff-dc4a79872612bb26927f9639df223856R1).
I consider this a success for the project. It is such a little thing, but it
means a lot to me for some reason. My shitpost created a test case in a project
I tried to integrate it with.

That's just awesome to me in ways I have trouble explaining.

As such, h programs _do_ work with Innative. Here's how to do it:

First, install the h compiler and runtime with the following command:

```
go get within.website/x/cmd/h
```

This will install the `h` binary to your `$GOPATH/bin`, so ensure that is part
of your path (if it is not already):

```
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

Then create a h binary like this:

```
h -p "h h" -o hh.wasm
```

Now we need to provide Innative the `h.h` system call implementation, so open
`h.c` and enter in the following:

```
#include <stdio.h>

void h_WASM_h(char data) {
  putchar(data);
}
```

Then build it to an object file:

```
gcc -c -o h.o h.c
```

Then pack it into a static library `.ar` file:

```
ar rsv libh.a h.o
```

Then create the shared object with Innative:

```
innative-cmd -l ./libh.a hh.wasm
```

This should create `hh.so` in the current working directory.

Now create the following [Nim](https://nim-lang.org) wrapper at `h.nim`:

```
proc hh_WASM_h() {. importc, dynlib: "./hh.so" .}

hh_WASM_h()
```

and build it:

```
nim c h.nim
```

then run it:

```
./h
h
```

And congrats, you have now compiled h to a native shared object.

## Why

Now, something you might be asking yourself as you read through this post is
something like: "Why the heck are you doing this?" That's honestly a good
question. One of the things I want to do with computers is to create art for the
sake of art. h is one of these such projects. h is not a productive tool. You
cannot create anything useful with h. This is an exercise in creating a compiler
and runtime from scratch, based on my past experiences with parsing lojban,
WebAssembly on the server and frustrating marketing around programming tools. I
wanted to create something that deliberately pokes at all of the common ways
that programming languages and tooling are advertised. I wanted to make it a
fully secure tool as well, with an arbitrary limitation of having no memory
usage. Everything is fully functional. There are a few grammar bugs that I'm
calling features.
