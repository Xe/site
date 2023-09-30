---
title: "pa'i: hello world!"
date: 2020-02-22
series: olin
tags:
  - rust
  - wasm
  - dhall
---

It's been a while since I gave an update on the Olin ecosystem (which now
exists, apparently). Not much has really gone on with it for the last few
months. However, recently I've decided to tackle one of the core problems of
Olin's implementation in Go: execution speed.

Originally I was going to try and handle this with
["hyperjit"](https://innative.dev), but support for linking C++ programs into Go
is always questionable at best. All of the WebAssembly compiling and
running tooling has been written in Rust, and as far as I know I was the only
holdout still using Go. This left me kinda stranded and on my own, seeing as the
libraries that I was using were starting to die.

I have been following the [wasmer][wasmer] project for a while and thanks to
their recent [custom ABI sample][wasmercustomabisample], I was able to start
re-implementing the Olin API in it. Wasmer uses a JIT for handling WebAssembly,
so I'm able to completely destroy the original Go implementation in terms of
performance. I call this newer, faster runtime pa'i (/pa.hi/, paw-hee), which
is a [Lojban][lojban] [rafsi][rafsi] for the word prami which means love.

[wasmer]: https://wasmer.io
[wasmercustomabisample]: https://github.com/wasmerio/wasmer-rust-customabi-example
[lojban]: https://mw.lojban.org/papri/Lojban
[rafsi]: https://lojban.org/publications/cll/cll_v1.1_xhtml-section-chunks/section-rafsi.html

[pa'i][pahi] is written in [Rust][rust]. It is built with [Nix][nix]. It
requires a nightly version of Rust because the WebAssembly code it compiles
requires it. However, because it is built with Nix, this quickly becomes a
non-issue. You can build pa'i by doing the following:

[pahi]: https://github.com/Xe/pahi
[rust]: https://www.rust-lang.org
[nix]: https://nixos.org/nix/

```console
$ git clone git@github.com:Xe/pahi
$ cd pahi
$ nix-build
```

and then `nix-build` will take care of:

- downloading the pinned nightly version of the rust compiler
- building the reference Olin interpreter
- building the pa'i runtime
- building a small suite of sample programs
- building the documentation from [dhall][dhall] files
- building a small test runner

[dhall]: https://dhall-lang.org

If you want to try this out in a more predictable environment, you can also
`nix-build docker.nix`. This will create a Docker image as the result of the Nix
build. This docker image includes [the pa'i composite package][pahidefaultnix],
bash, coreutils and `dhall-to-json` (which is required by the test runner).

[pahidefaultnix]: https://github.com/Xe/pahi/blob/master/default.nix

I'm actually really proud of how the documentation generation works. The
[cwa-spec folder in Olin][cwaspecolin] was done very ad-hoc and was only
consistent because there was a template. This time functions, types, errors,
namespaces and the underlying WebAssembly types they boil down to are all
implemented as Dhall records. For example, here's the definition of a
[namespace][cwans] [in Dhall][nsdhall]:

[cwaspecolin]: https://github.com/Xe/olin/tree/master/docs/cwa-spec
[cwans]: https://github.com/Xe/pahi/tree/master/olin-spec#namespaces
[nsdhall]: https://github.com/Xe/pahi/blob/5ea1184c09df4e657524f9d5e77941cda5560d9a/olin-spec/types/ns.dhall

```
let func = ./func.dhall

in  { Type = { name : Text, desc : Text, funcs : List func.Type }
    , default =
        { name = "unknown"
        , desc = "please fill in the desc field"
        , funcs = [] : List func.Type
        }
    }
```

which gets rendered to [Markdown][markdown] using
[`renderNSToMD.dhall`][shownsasmd]:

[markdown]: https://github.github.com/gfm/
[shownsasmd]: https://github.com/Xe/pahi/blob/5ea1184c09df4e657524f9d5e77941cda5560d9a/olin-spec/types/renderNSToMD.dhall

```
let ns = ./ns.dhall

let func = ./func.dhall

let type = ./type.dhall

let showFunc = ./renderFuncToMD.dhall

let Prelude = ../Prelude.dhall

let toList = Prelude.Text.concatMapSep "\n" func.Type showFunc

let show
    : ns.Type → Text
    =   λ(namespace : ns.Type)
      → ''
        # ${namespace.name}

        ${namespace.desc}

        ${toList namespace.funcs}
        ''

in  show
```

This would render [the logging namespace][logns] as [this markdown][lognsmd].

[logns]: https://github.com/Xe/pahi/blob/5ea1184c09df4e657524f9d5e77941cda5560d9a/olin-spec/ns/log.dhall
[lognsmd]: https://github.com/Xe/pahi/blob/5ea1184c09df4e657524f9d5e77941cda5560d9a/olin-spec/ns/log.md

It seems like overkill to document things like this (and at some level it is),
but I plan to take advantage of this later when I need to do things like
generate C/Rust/Go/TinyGo bindings for the entire specification at once. I also
have always wanted to document something so precisely like this, and now I get
the chance.

pa'i is just over a week old at this point, and as such it is NOT
[feature-complete with the reference Olin interpreter][compattodo]. I'm working
on it though. I'm kinda burnt out from work, and even though working on this
project helps me relax (don't ask me how, I don't understand either) I have
limits and will take this slowly and carefully to ensure that it stays
compatible with all of the code I have already written in Olin's repo. Thanks to
[go-flag][goflags], I might actually be able to get it mostly flag-compatible.
We'll see though.

[compattodo]: https://github.com/Xe/pahi/issues/1
[goflags]: https://crates.io/crates/go-flag

I have also designed a placeholder logo for pa'i. Here it is:

![the logo for pa'i](/static/blog/pahi-logo.png)

It might be changed in the future, but this is what I am going with for now. The
circuit traces all spell out messages of love (inspired from the Senzar runes of
the [WingMakers][wingmakers]). The text on top of the microprocessor reads pa'i
in [zbalermorna][zbalermorna], a constructed writing script for Lojban. The text
on the side probably needs to be revised, but it says something along the lines
of "a future after programs".

[wingmakers]: https://www.wingmakers.us/wingmakersorig/wingmakers/ancient_arrow_project.shtml
[zbalermorna]: https://mw.lojban.org/images/b/b3/ZLM4_Writeup_v2.pdf

pa'i is chugging along. When I have closed the [compatibility todo
list][compattodo] for all of the Olin API calls, I'll write more. For now, pa'i
is a very complicated tool that lets you print "Hello, world" in new and
exciting ways (this will change once I get resource calls into it), but it's
getting there.

I hope this was interesting. Be well.
