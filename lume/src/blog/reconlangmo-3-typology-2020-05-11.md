---
title: "ReConLangMo 3: Morphosyntactic Typology"
date: 2020-05-11
series: reconlangmo
tags:
 - conlang
 - lewa
---

In the last post of [this series][reconlangmoseries], we covered the sounds and
word patterns of L'ewa. This time we are covering morphosyntactic typology, or
how words and sentences are formed out of root words, details about sentences,
word order and those kinds of patterns. I'll split each of these into their own
headings so it's a bit easier to grok. This is a response to [this
prompt][rclm3].

[reconlangmoseries]: /blog/series/reconlangmo
[rclm3]: https://www.reddit.com/r/conlangs/comments/ghvo48/reconlangmo_3_morphosyntactic_typology/

## Word Order

L'ewa is normally a Subject-Verb-Object (SVO) language like English. However,
the word order of a sentence can be changed if it is important to specify some
part of the sentence in particular.

I haven't completely finalized the particles for this, but I'd like to use `ka` to
denote the subject, `ke` to denote the verb and `ku` to denote the object. For
example if the input sentence is something like:

```
/mi/ /mad.sa/ /lo/ /spa.lo/
mi   madsa    lo   spalo
 I   eat      an   apple
```

You could emphasize the eating with:

```
/kɛ/ /mad.sa/ /ka/ /mi/ /lo/ /spa.lo/
[ke] madsa    ka   mi   lo   spalo
V    eat      S    I    an   apple
```

(the `ke` is in square brackets here because it is technically not required, but
it can make more sense to be explicit in some cases)

or the apple with:

```
/ku/ /lo/ /spalo/ /kɛ/ /mad.sa/ /mi
ku   lo   spalo   ke   madsa    mi
O    an   apple   V    eat      I
```

L'ewa doesn't really have adjectives or adverbs in the normal indo-european
sense, but it does have a way to analytically combine meanings together. For
example if `qa'te` is the word for `is fast/quick/rapid in rate`, then saying
you are quickly eating (or wolfing food down) would be something like:

```
/qaʔ.tɛ/          /mad.sa/
qa'te             madsa
is fast [kind of] eat
```

These are assumed to be metaphorical by default. It's not always clear what
someone would mean by a fast kind of language (would they be referencing
[Speedtalk][speedtalk]?)

[speedtalk]: https://en.wikipedia.org/wiki/Speedtalk

L'ewa doesn't always require a subject or object if it can be figured out from
context. You can just say "rain" instead of "it's raining". By default, the
first word in a sentence without an article is the verb. The ka/ke/ku series
needs to be used if the word order deviates from Subject-Verb-Object (it
functions a lot like the selma'o FA from Lojban).

## Morphological Typology

L'ewa is a analytic language. Every single word has only one form and particles
are used to modify the meaning or significance of words. There are only two word
classes: content and particles.

### Alignment

L'ewa is a nominative-accusative language. Other particles may be introduced in
the future to help denote the relations that exist in other alignments, but I
don't need them yet.

### Word Classes

As said before, L'ewa only has two word classes, content (or verbs) and
particles to modify the significance or relations between content. There is also
a hard limit of two arguments per verb, which should help avoid the problems
that Lojban has with its inconsistent usage of the x3, x4 and x5 places.

As the content words are all technically verbs, there is no real need for a
copula. The ka/ke/ku series can also help to break out of other things that
modify "noun-phrases" (when those things exist). There are also no nouns,
adjectives or adverbs, because analytically combining words completely replaces
the need for them.

Nouns and verbs do not inflect for numbers. If numbers are needed they can be
provided, otherwise the default is to assume "one or more".

## Conscript

I am still working on the finer details of the conscript for L'ewa, but here is
a sneak preview of the letter forms I am playing with (this image below might
not render properly in light mode):

![The letters in the L'ewa
conscript](https://pbs.twimg.com/media/EXwr2rIWAAE95co?format=png&name=4096x4096)

My inspirations for this script were [zbalermorna][zbalermorna], Hangul, Hanzi,
Katakana, Greek, international computer symbols, traditional Japanese art and
the [International Phonetic Alphabet][ipa].

[zbalermorna]: https://mw.lojban.org/images/b/b3/ZLM4_Writeup_v2.pdf
[ipa]: https://en.wikipedia.org/wiki/International_Phonetic_Alphabet

This script is very decorative, and is primarily intended to be used in
spellcraft and other artistic uses. It will probably show up in my art from time
to time, and will definitely show up in any experimental video production that I
work on in the future. I will go into more detail about this in the future, but
here is my prototype. Please do let me know what you think about it.

---

As a side note, the words `madsa`, `spalo` and `qa'te` are now official L'ewa
words, I guess. The entire vocabulary of the language can now be listed below:

**Content Words**

| L'ewa word | IPA        | English                     |
| ---------- | ---        | -------                     |
| `l'ewa`    | `/lʔ.ɛwa/` | is a language               |
| `madsa`    | `/mad.sa/` | eats/is eating              |
| `qa'te`    | `/qaʔ.tɛ/` | is fast/quick/rapid in rate |
| `zasko`    | `/ʒa.sko/` | is a plant/is vegetation    |
| `spalo`    | `/spa.lo/` | is an apple                 |

**Particles**

| L'ewa word | IPA  | English                   |
| ---------- | ---  | -------                   |
| lo         | /lo/ | a, an, indefinite article |
| le         | /lɛ/ | the, definite article     |
| ka         | /ka/ | subject marker            |
| ke         | /kɛ/ | verb marker               |
| ku         | /ku/ | object marker             |
| mi         | /mi/ | the current speaker       |
