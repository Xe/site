---
title: "ReConLangMo 7: Discourse"
date: 2020-05-25
series: reconlangmo
tags:
 - conlang
 - lewa
---

Previously on [ReConLangMo][reconlangmo], we covered a lot of new words for the
lexicon of L'ewa. This helps to flesh out a lot of what can be said, but
conversations themselves can be entirely different from formal sentences.
Conversations flow and ebb based on the needs/wants of the interlocutors. This
post will start to cover a lot of the softer skills behind L'ewa as well as
cover some other changes I'm making under the hood. This is a response to [this
prompt][rclm7].

[reconlangmo]: https://xeiaso.net/blog/series/reconlangmo
[rclm7]: https://www.reddit.com/r/conlangs/comments/gqo8jn/reconlangmo_7_discourse/

## Information Structure

L'ewa doesn't have any particular structure for marking previously known
information, as normal sentences should suffice in most cases. Consider this
paragraph:

```
I saw you eat an apple. Was it tasty?
```

Since `an apple` was the last thing mentioned in the paragraph, the vague "it"
pronoun in the second sentence can be interpreted as "the apple".

L'ewa doesn't have a way to mark the topic of a sentence, that should be obvious
from context (additional clauses to describe things will help here). In most
cases the subject should be equivalent to the topic of a sentence.

L'ewa doesn't directly offer ways to emphasize parts of sentences with phonemic
stress like English does (eg: "I THOUGHT you ate an apple" vs "I thought you ATE
an apple"), but emotion words can be used to help indicate feelings about
things, which should suffice as far as emphasis goes.

## Discourse Structure

Conversationally, a lot of things in L'ewa grammar get dropped unless it's
ambiguous. The I/yous that get tacked on in English are completely unneeded. A
completely valid conversation could look something like this:

```
<Mai> xoi
<Cadey> xoi
<Mai> xoi madsa?
<Cadey> lo spalo
```

And it would roughly equate to:

```
<Mai> Hi
<Cadey> Hi, you doing okay?
<Mai> Yes, have you eaten?
<Cadey> Yes, I ate an apple
```

People know when they can speak after a sufficient pause between utterances.
Interrupting is not common but not a social faux-pas, and can be used to stop a
false assumption from being said.

## Utterances

An utterance in L'ewa is anything from a single content word all the way up to
an entire paragraph of sentences. An emotion particle can be a complete
utterance. A question particle can be a complete utterance, anything can be an
utterance. A speaker may want to choose more succinct options when the other
detail is already contextually known or simply not relevant to the listener.

L'ewa has a few discourse particles, here are a few of the more significant
ones:

| L'ewa | Function                                             |
|-------|------------------------------------------------------|
| xi    | signals that the verb of the sentence is coming next |
| ko    | ends a noun phrase                                   |
| ka    | marks something as the subject of the sentence       |
| ke    | marks something as the verb of the sentence          |
| ku    | marks something as the object of the sentence        |

## Formality

The informal dialect of L'ewa drops everything it can. The formal dialect
retains everything it can, to the point where it includes noun phrase endings,
the verb signaler, ka/ke/ku and every single optional particle in the language.
The formal dialect will end up sounding rather wordy compared to informal slangy
speech. Consider the differences between informal and formal versions of "I eat
an apple":

```
mi madsa lo spalo.
```

```
ka mi ko xi ke madsa ku lo spalo ko.
```

Nearly all of those particles are not required in informal speech (you could
even get away with `madsa lo spalo` depending on context), but are required in
formal speech to ensure there is as little contextual confusion as possible.
Things like laws or legal rulings would be written out in the formal register.

## Greetings and Farewell

"Hello" in L'ewa is said using `xoi`. It can also be used as a reply to hello
similar to «ça va» in French. It is possible to have an entire conversation with
just `xoi`:

```
<Mai> xoi
<Cadey> xoi
<Mai> xoi
```

The other implications of `xoi` are "how are you?" "I am good, you?", "I am
good", etc. If more detail is needed beyond this, then it can be supplied
instead of replying with `xoi`.

"Goodbye" is said using `xei`. Like `xoi` it can be used as a reply to another
goodbye and can form a mini-conversation:

```
<Cadey> xei
<Mai> xei
<Cadey> xei
```

## Emotion Words

Feelings in L'ewa are marked with a family of particles called "UI". These can
also be modified with other particles. Here are the emotional markers:

| L'ewa | English        |
|-------|----------------|
| `a'a` | attentive      |
| `a'e` | alertness      |
| `ai`  | intent         |
| `a'i` | effort         |
| `a'o` | hope           |
| `au`  | desire         |
| `a'u` | interest       |
| `e'a` | permission     |
| `e'e` | competence     |
| `ei`  | obligation     |
| `e'i` | constraint     |
| `e'o` | request        |
| `e'u` | suggestion     |
| `ia`  | belief         |
| `i'a` | acceptance     |
| `ie`  | agreement      |
| `i'e` | approval       |
| `ii`  | fear           |
| `i'i` | togetherness   |
| `io`  | respect        |
| `i'o` | appreciation   |
| `iu`  | love           |
| `i'u` | familiarity    |
| `o'a` | pride          |
| `o'e` | closeness      |
| `oi`  | complaint/pain |
| `o'i` | caution        |
| `o'o` | patience       |
| `o'u` | relaxation     |
| `ua`  | discovery      |
| `u'a` | gain           |
| `ue`  | surprise       |
| `u'e` | wonder         |
| `ui`  | happiness      |
| `u'i` | amusement      |
| `uo`  | completion     |
| `u'o` | courage        |
| `uu`  | pity           |
| `u'u` | repentant      |

If an emotion is unknown in a conversation, you can ask with `kei`:

```
<Mai> xoi, so kei?
      hi,  what-verb what-feeling?

<Cadey> madsa ui
        eating :D
```

This system is wholesale stolen from [Lojban](https://lojban.github.io/cll/13/1/).

## Connectives

Connectives exist to link noun phrases and verbs together into larger
noun phrases and verbs. They can also be used to link together sentences. There
are four simple connectives: `fa` (OR), `fe` (AND), `fi` (connective question),
`fo` (if-and-only-if) and `fu` (whether-or-not).

### OR

```
ro au madsa lo spalo fa lo hafto?
Do you want to eat an apple or an egg?
```

### AND

```
ro au madsa lo spalo fe lo hafto?
Do you want to eat an apple and an egg?
```

### If and Only If

```
ro 'amwo mi fo mi madsa hafto?
Do you love me if I eat eggs?
```

### Whether or Not

```
mi 'amwo ro. fu ro madsa hafto.
I love you, whether or not you eat eggs.
```

### Connective Question

```
ro au madsa lo spalo fi lo hafto?
Do you want to eat apples and/or eggs?
```

## Changes Being Made to L'ewa

Early on, I mentioned that family terms were gendered. This also ended up with
me making some gendered terms for people. I have since refactored out all of the
gendered terms in favor of more universal terms. Here is a table of some of the
terms that have been replaced:

| English                 | L'ewa term  | L'ewa word |
|-------------------------|-------------|------------|
| brother/sister          | sibling     | xinga      |
| mother/father           | parent      | pa'ma      |
| grandfather/grandmother | grandparent | gra'u      |
| aunt/uncle              | parent      | pa'ma      |
| cousin                  | sibling     | xinga      |
| man/woman               | Creator     | kirta      |
| man/woman               | human       | renma      |

In some senses, gender exists. In other senses, gender does not. With L'ewa I
want to explore what is possible with language. It would be interesting to
create a language where gender can be discussed as it is, not as the categories
that it has historically fit into. Consider colors. There are millions of
colors, all sightly different but many follow general patterns. No one or two
colors can be thought of as the "default" color, yet we can have long and
meaningful conversations about what color is and what separates colors from
eachother.

I aim to have the same kind of granularity in L'ewa. As a goal of the language,
I should be able to point to any content word in the dictionary and be able to
say "that's my gender" in the same way I can describe color or music with that
tree. These will implicitly be metaphors (which does detract a bit from the
logical stance L'ewa normally takes) because gender is almost always a metaphor
in practice. L'ewa will not have binary gender.

Issue [number two](https://tulpa.dev/cadey/lewa/issues/2) on the L'ewa repo will
help track the creation and implementation of a truly non-binary "gender" system
for L'ewa.

---

I've been chugging through the Swaedish list more and more to build up more of
L'ewa's vocabulary in preparation for starting to translate sentences more
complicated than simple "I eat an apple" or "Do you like eating plants?". One of
the first things I want to translate is the classic [tower of babel
story][babel].

[babel]: https://en.wikipedia.org/wiki/Tower_of_Babel

Be well.
