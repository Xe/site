---
title: "ReConLangMo 6: Lexicon"
date: 2020-05-22
series: reconlangmo
tags:
 - conlang
 - lewa
---

Previously in [this series][reconlangmo], we've covered a lot of details about
how sentences work, tenses get marked and how words work in general; however
this doesn't really make L'ewa a _language_. Most of the difficulty in making a
language like this is the vocabulary. In this post I'll be describing how I am
making the vocabulary for L'ewa and I'll include an entire table of the
dictionary words. This answers [this
prompt](https://www.reddit.com/r/conlangs/comments/gojncp/reconlangmo_6_lexicon/).

[reconlangmo]: https://xeiaso.net/blog/series/reconlangmo

## Word Distinctions

L'ewa is intended to be a logical language. One of the side effects of L'ewa
being a logical language is that each word should have as minimal and exact of a
meaning/function as possible. English has lots of words that cover large
semantic spaces (like go, set, run, take, get, turn, good, etc.) without much of a
pattern to it. I don't want this in L'ewa. 

Let's take the word "good" as an example. Off the top of my head, good can mean
any of the following things:

- beneficial
- aesthetically pleasing
- favorful taste
- saintly (coincidentally this is the source of the idiom "God is good")
- healthy

I'm fairly sure there are more "senses" of the word good, but let's break these
into their own words:

| L'ewa | Definition                         |
|-------|------------------------------------|
| firgu | is beneficial/nice to              |
| n'ixu | is aesthetically pleasing to       |
| flawo | is tasty/has a pleasant flavor to  |
| spiro | is saintly/holy/morally good to    |
| qanro | is healthy/fit/well/in good health |

Each of these words has a very distinct and fine-grained meaning, even though
the range is a bit larger than it would be in English. These words also differ
from a lot of the other words in the L'ewa dictionary so far because they can
take an object. Most of the words so far are adjective-like because it doesn't
make sense for there to be an object attached to the color blue.

By default, if a word that can take an object doesn't have one, it's assumed to
be obvious from context. For example, consider the following set of sentences:

```
mi qa madsa lo spalo. ti flawo!

I am eating an apple. It's delicious!
```

I am working at creating more words using a [Swaedish list][swaedish207].

[swaedish207]: https://tulpa.dev/cadey/lewa/src/branch/master/words/swaedish207.csv

## Family Words

Family words are a huge part of a language because it encodes a lot about the
culture behind that language. L'ewa isn't really intended to have much of a
culture behind it, but the one place I want to take a cultural stance is here.
The major kinship word is kirta, or "is an infinite slice of an even greater
infinite". This is one of the few literal words in L'ewa that is defined using a
metaphor, as there is really no good analog for this in English.

There are also words for other major family terms in English:

| L'ewa | Definition              |
|-------|-------------------------|
| brota | is the/a brother of     |
| sistu | is the/a sister of      |
| mamta | is the/a mother of      |
| patfu | is the/a father of      |
| grafa | is the/a grandfather of |
| grama | is the/a grandmother of |
| wanto | is the/a aunt of        |
| tunke | is the/a uncle of       |

Cousins are all called brother/sister. None of these words are inherently
gendered and `brota` can refer to a female or nonbinary person. The words are
separate because I feel it flows better, for now at least.

## Idioms

L'ewa strives to have as few idioms as possible. If something is meant
non-literally (or as a [conceptual metaphor][cmet]), the particle ke'a can be used:

[cmet]: https://en.wikipedia.org/wiki/Conceptual_metaphor

```
ti firgu
This is beneificial

ti ke'a firgu
This is metaphorically/non-literally beneficial
```

---

I have been documenting L'ewa and all of its words/grammar in a [git
repo][lewarepo]. The layout of this repo is as follows:

| Folder   | Purpose                                                                                                                |
|----------|------------------------------------------------------------------------------------------------------------------------|
| `book`   | The source files and build scripts for the L'ewa book (this book may end up being published)                           |
| `nix`    | [Nix][nix] crud, custom packages for the eBook render and development tools                                            |
| `script` | Where experiments for the written form of L'ewa live                                                                   |
| `tools`  | Tools for hacking at L'ewa in Rust/Typescript (none published yet, this is where the dictionary server code will live) |
| `words`   | Where the definitions of each word are defined in [Dhall][dhall], this will be fed into the dictionary server code       |

I also have the entire process of building and testing everything (from the
eBook to the unit tests of the tools) automated with [Drone][droneci].
Eventually this will be automatically deployed to my Kubernetes cluster
and the book will be a subpath/subdomain of `lewa.within.website`.

I have created a system of defining words that allows you to focus on each word
at once, but then fit it back into the greater whole of the language. For
example here is `kirta.dhall`:

```dhall
-- kirta.dhall
let ContentWord = ../types/ContentWord.dhall

in  ContentWord::{
    , word = "kirta"
    , gloss = "Creator"
    , definition =
        "is an infinite slice of an even greater infinite/our Creator/a Creator"
    }
```

This is put in `words/roots` because it is a root (or uncombined) word. Then it
is added to the `dictionary.dhall`:

```dhall
-- dictionary.dhall
let ContentWord = ./types/ContentWord.dhall

let ParticleWord = ./types/ParticleWord.dhall

in  { rootWords =
      [ -- ...
      ./roots/kirta.dhall
      -- ...
      ]
    , particles [ -- ...
    ]
```

And then the build process will automatically generate the new dictionary from
all of these definitions. Downside of this is that each new kind of word needs
subtle adjustments to the build process of the dictionary and that
removals/changes to lots of words requires a larger-scale refactor of the
language, but I feel the tradeoff is worth the effort. I will undoubtedly end up
creating a few tools to help with this.

I will keep working on additional vocabulary on my own, but [here][vocab] is the
list of vocabulary that has been written up so far.

[vocab]: https://git.io/JfaeF

Be well.

[lewarepo]: https://tulpa.dev/cadey/lewa
[nix]: https://nixos.org/nix/
[dhall]: https://dhall-lang.org/
[droneci]: https://drone.io
