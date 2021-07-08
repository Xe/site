---
title: "ReConLangMo 5: Sentence Structure"
date: 2020-05-18
series: reconlangmo
tags:
 - conlang
 - lewa
---

The last post in [this series][reconlangmo] was more of a grammar dump with few
concrete examples or much details about things (mostly because of a lack of
vocabulary to make examples with). I'll fix this in the future, but for now
let's continue on with sentence structure goodness. This is a response to [this
prompt][rclm5].

[reconlangmo]: /blog/series/reconlangmo
[rclm5]: https://www.reddit.com/r/conlangs/comments/gmbwb5/reconlangmo_5_sentence_structure/

## Independent Clause Structure

Most of the time L'ewa sentences have only one clause. This can be anything from
a single verb to a subject, verb and object. However, sometimes more information
is needed. Consider this sentence:

```
The dog which is blue is large.
```

This kind of a relative clause would be denoted using `hoi`, which would make
the sentence roughly the following in L'ewa:

```
le wufra hoi blanu xi brado.
```

The particle `xi` is needed here in order to make it explicit that the subject
noun-phrase has ended.

Similarly, an incidental relative clause is done with with `joi`:

```
le  wufra  joi              blanu    ke brado
the dog,   which by the way is blue,    is big.
```

## Questions

There are a few ways to ask questions in L'ewa. They correlate to the different
kinds of things that the speaker could want to know. 

### `ma`

`ma` is the particle used to fill in a missing/unknown noun phrase. Consider
these sentences:

```
ma   blanu?
what is blue?
```

```
ro  qa madsa   ma?
you are eating what?
```

### `no`

`no` is the particle used to fill in a missing/unknown verb. Consider these
sentences:

```
ro no?
How are you doing?
```

```
le wufra xi no?
The dog did what?
```

### `so`

`so` is the particle used to ask questions about numbers, similar to the "how
many" construct in English.

```
ro madsa so spalo?
You ate how many apples?
```

```
le so zasko xi qa'te glowa
How many plants grow quickly?
```

## Color Words

L'ewa uses a RGB color system like computers. The basic colors are red, green
and blue, with some other basic ones for convenience:

| English  | L'ewa  |
| -------  | ------ |
| blue     | blanu  |
| red      | delja  |
| green    | qalno  |
| yellow   | yeplo  |
| teal     | te'ra  |
| pink     | hetlo  |
| black    | xekri  |
| white    | pu'ro  |
| 50% gray | flego  |

Colors will be mixed by creating compound words between base colors. Compound
words still need to be fleshed out, but generally all CVCCV words will have
wordparts made out of the first, second and fifth letter, unless the vowel pair
is illegal and all CCVCV words are the first, second and fifth letter unless
this otherwise violates the morphology rules. Like I said though, this really
needs to be fleshed out and this is only a preview for now.

For example a light green would be `puoqa'o` (`pu'lo qalno`, white-green).

---

I hit a snag while hacking at the tooling for making word creation and the like
easier. I am still working on it, but most of my word creation is manual and
requires me to keep a phonology information document up on my monitor while I
sound things out. As part of writing this article I had to add the letters `f`
and `r` to L'ewa for the word `wufra`.

I am documenting my work for this language [here](https://tulpa.dev/cadey/lewa).
This repo will build the grammar book PDF, website and eBook. This will also be
the home of the word generation, similarity calculation, dictionary and
(eventually) automatic translation tools. I am also documenting each of the
words in the language in their own files that will feed into the grammar book
generation. More on this when I have more of a coherent product!
