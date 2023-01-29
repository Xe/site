---
title: "The Next-Generation Universal Hlang compiler"
date: 2022-12-31
series: h
tags:
 - hlang
 - wasm
vod:
  twitch: https://www.twitch.tv/videos/1693936831
  youtube: https://youtu.be/QY1O2n4tOhE
---

In a world where simple tasks have hundreds of dependencies and most of them are
not documented, everything falls to chaos. The monolithigarchy dictates that
your build times must be slow so that They (the dependocracy) can win over your
hearts and minds with video games that you play during your compile times. One
person gets mad about their string padding library being used by corporations
without paying and then the entire internet explodes for a few days. This is
unsustainable.

hlang is the sledgehammer that will break down this complexity and deliver you a
truly uncompromised development experience.

<xeblog-conv name="Numa" mood="delet">You can't spell _sledgehammer_ without
_h_!</xeblog-conv>

If none of this is making any sense, please read [the rest of the
series](https://xeiaso.net/blog/series/h). This will hopefully help something
make sense.

<xeblog-conv name="Numa" mood="delet">If you need even more context, check [this
page](https://pkg.go.dev/context) for more information.</xeblog-conv>

There was one major flaw with hlang in the past though. It was a hollow shell of
itself and had rot to the slains and arrows of time. The playground stopped
working, so people could not understand the sheer might of hlang by playing with
it.

Lo, behold, a new compiler was born. In this article, I will describe the nguh
compiler and how it revolutionizes the ways that you use hlang for both
professional and personal uses.

<xeblog-conv name="Mara" mood="wat">Wait, what, there _were_ professional users
of hlang???</xeblog-conv>

<xeblog-conv name="Numa" mood="delet">Having 2 years of hlang on your resume
will let you get hired by Google!</xeblog-conv>

## The Old Compiler

The old compiler was a HACK. The main way it worked was by feeding the program
source code as a string to this [Go template](https://pkg.go.dev/text/template):

```
(module
 (import "h" "h" (func $h (param i32)))
 (func $h_main
       (local i32 i32 i32)
       (local.set 0 (i32.const 10))
       (local.set 1 (i32.const 104))
       (local.set 2 (i32.const 39))
       {{ range . -}}
       {{ if eq . 32 -}}
       (call $h (get_local 0))
       {{ end -}}
       {{ if eq . 104 -}}
       (call $h (get_local 1))
       {{ end -}}
       {{ if eq . 39 -}}
       (call $h (get_local 2))
       {{ end -}}
       {{ end -}}
       (call $h (get_local 0))
 )
 (export "h" (func $h_main))
)
```

This template worked by taking the program input _as a string_ and looping over
each character to decide what to do. If it was a space, it would print a
newline. If it was an `h`, it would print `h`. If it was a `'`, it would print a
`'`. Anything else is ignored.

However, this means that the parser was mostly ignored. And the parser spec
compiles to 117 bytes when gzipped, which means that it can fit on a tshirt.

<xeblog-conv name="Numa" mood="delet">That's a savings of 0.8475%!</xeblog-conv>

Additionally, this would then use the command
[`wat2wasm`](https://developer.mozilla.org/en-US/docs/WebAssembly/Text_format_to_wasm)
to compile it to a WebAssembly file instead of doing it directly. This combined
with the fact that the `get_local` instruction was renamed to `local.get` in the
text format some time in the last 2 years means that not only was my compiler
hacky, it didn't work anymore.

<xeblog-conv name="Mara" mood="hacker">Apparently that was renamed before WASM
hit 1.0 and the legacy name was an alias they planned to remove. Guess who
didn't get the memo!</xeblog-conv>

Needless to say, this could be fixed by doing a simple
`s/get_local/local\.get/g` on the source file, but that's not fun. You know
what's really fun? Reverse-engineering a binary file on stream and reassembling
an identical replica in code. That's fun.

## The nguh compiler

On December 31st, 2022, I wrote the nguh compiler [on
stream](https://www.twitch.tv/princessxen). The nguh (nguh gives u hlang or
Next-Generation Universal Hlang compiler, whichever you prefer) compiler outputs
WebAssembly bytecode directly instead of using `wat2wasm` as a middleman.

<xeblog-conv name="Mara" mood="happy">This means that hlang has even fewer
dependencies!</xeblog-conv>

nguh is supposed to be pronounced with the final sound of `-ing` and `uh`
smashed together. It is not phonetically valid in English. It will take some
practice to say it correctly. I'm not sorry. If you can read IPA, it's
pronounced /ŋə/. The name comes from the youtuber [Agma
Schwa](https://www.youtube.com/@AgmaSchwa)'s show about conlangs named /ŋə/.

To help you understand the architecture of nguh, it will be helpful to get some
context about how WebAssembly files work.

## How WebAssembly files work

<details>
  <summary>What is WebAssembly?</summary>
  
WebAssembly is a standard that specifies a way to run programs on arbitrary
hardware in a sandboxed way. It is used mainly in web browsers to power things
like YouTube's player component, Twitch stream viewing, and by developers any
time they need to put a block of code into a website without having to rewrite
it in JavaScript.

I'm part of a slowly growing group of developers that want to run WebAssembly
code on the server so that you can take the same `.wasm` file and run it on any
hardware without having to have the source code and a working compiler setup.

hlang is compiled to WebAssembly for no reason in particular.
</details>

At a high level, a WebAssembly module has a bunch of sections in it. Each
section contains information for things like what functions the module exports,
the types of imported fuctions, how much memory the module needs, what should be
in memory by default, and the function bodies for your code. Here's an annotated
disassembly of a hlang binary:

```
0x00, 0x61, 0x73, 0x6d, // \0asm wasm magic number
0x01, 0x00, 0x00, 0x00, // version 1

0x01, // type section
0x08, // 8 bytes long
0x02, // 2 entries
0x60, 0x01, 0x7f, 0x00, // function type 0, 1 i32 param, 0 return
0x60, 0x00, 0x00, // function type 1, 0 param, 0 return

0x02, // import section
0x07, // 7 bytes long
0x01, // 1 entry
0x01, 0x68, // module h
0x01, 0x68, // name h
0x00, // type index
0x00, // function number

0x03, // func section
0x02, // 2 bytes long
0x01, // function 1
0x01, // type 1

0x07, // export section
0x05, // 5 bytes long
0x01, // 1 entry
0x01, 0x68, // "h"
0x00, 0x01, // function 1

0x0a, // code section
0x1b, // 27 bytes long
0x01, // 1 entry
0x19, // 25 bytes long
0x01, // 1 local declaration
0x03, 0x7f, // 3 i32 values - (local i32 i32 i32)
0x41, 0x0a, // i32.const 10 (newline)
0x21, 0x00, // local.set 0
0x41, 0xe8, 0x00, // i32.const 104 (h)
0x21, 0x01, // local.set 1
0x41, 0x27, // i32.const 39 (')
0x21, 0x02, // local.set 2
0x20, 0x01, // local.get 1 push h
0x10, 0x00, // call 0 (putchar)
0x20, 0x00, // local.get 0 push newline
0x10, 0x00, // call 0 (putchar)
0x0b // end of function
```

At a high level, nguh just takes all the needed sections and [puts them in the
target
binary](https://github.com/Xe/x/blob/2fe527950512b97a544d2d59539026514ad59544/cmd/hlang/nguh/compile.go#L53).
Most of the sections are copied verbatim from that disassembly I pasted above
because they don't need any modification for the binary to work.

The exciting part happens when the individual nodes in the hlang syntax tree get
compiled to WebAssembly bytecode. Each node in the tree has maybe its character
to print and maybe a list of child nodes. A syntax tree for hlang could look
like this if it has one character in the program:

```
input: h
H("h")
```

Or it could look like this if there are multiple characters in the program:

```
input: h h h
H{
	"h",
	"h",
	"h",
}
```

This means I need something like this:

```go
// compile AST to wasm
if len(tree.Kids) == 0 {
    if err := compileOneNode(funcBuf, tree); err != nil {
        return nil, err
    }
} else {
    for _, node := range tree.Kids {
        if err := compileOneNode(funcBuf, node); err != nil {
            return nil, err
        }
    }
}
```

This will either read from the root of the tree or all of the tree's children in
order to compile the entire program. The `compileOneNode` function will turn the
text associated with the node into the correlating WASM bytecode (pushing the
relevant character to the stack and then calling the `h.h` (`putchar`) function).

Finally it will generate the end of the function including a trailing newline
and end the `.wasm` file.

<xeblog-conv name="Mara" mood="hacker">Fun fact: the generated binary for a
hlang program that only prints `h` is 69 bytes.</xeblog-conv>

<xeblog-conv name="Numa" mood="delet">NICE!</xeblog-conv>

Here is a base-64 encoded hlang binary in case you find this interesting:

```
AGFzbQEAAAABCAJgAX8AYAAAAgcBAWgBaAAAAwIB
AQcFAQFoAAEKHQEbAQN/QQohAEHoACEBQSchAiAB
EAAgABAAAQEL
```

---

If you want to play with hlang, head to its new home at
[h.within.lgbt](https://h.within.lgbt). If you want to witness things such as
this being created live, follow me [on twitch](https://www.twitch.tv/princessxen) or
on my VTuber business account at [@xe@vt.social](https://vt.social/@xe).

<xeblog-conv name="Cadey" mood="enby">Happy new year to those that
celebrate!</xeblog-conv>
