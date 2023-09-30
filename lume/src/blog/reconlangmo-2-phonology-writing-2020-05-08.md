---
title: "ReConLangMo 2: Phonology & Writing"
date: 2020-05-08
series: reconlangmo
tags:
  - conlang
  - lewa
---

Continuing from [the last post][rclm1], one of the next steps in this process is
to outline the phonology and basic phonotactics of L'ewa. A language's phonology
is the set of sounds that are allowed to be in words. The phonotactics of a
language help people understand where the boundaries between syllables are. I
will then describe my plans for the L'ewa orthography and how L'ewa is
romanized. This is a response to the prompt made [here][rclm2prompt].

[rclm1]: https://xeiaso.net/blog/reconlangmo-1-name-ctx-history-2020-05-05
[rclm2prompt]: https://www.reddit.com/r/conlangs/comments/gfp3hw/reconlangmo_2_phonology_writing/

## Phonology

I am taking inspiration from Lojban, Esperanto, Mandarin Chinese and English to
design the phonology of L'ewa. All of the phonology will be defined using the
[International Phonetic Alphabet][ipa]. If you want to figure out how to
pronounce these sounds, a lazy trick is to google them. Wikipedia will have a
perfectly good example to use as a reference. There are two kinds of sounds in
L'ewa, consonants and vowels.

[ipa]: https://en.wikipedia.org/wiki/International_Phonetic_Alphabet

### Consonants

*Consonant inventory*: /d f g h j k l m n p q s t w ʃ ʒ ʔ ʙ̥/

| Manner/Place        | Bilabial | Alveolar | Palato-alveolar | Palatal | Velar | Labio-velar | Uvular | Glottal |
|---------------------|----------|----------|-----------------|---------|-------|-------------|--------|---------|
| Nasal               | m        | n        |                 |         |       |             |        |         |
| Stop                | p        | t d      |                 |         | k g   |             | q      | ʔ       |
| Fricative           | f        | s        | ʃ ʒ             |         |       |             |        | h       |
| Approximant         |          |          |                 | j       |       | w           |        |         |
| Trill               | ʙ̥        | r        |                 |         |       |             |        |         |
| Lateral approximant |          | l        |                 |         |       |             |        |         |

The weirdest consonant is /ʙ̥/, which is a voiceless bilabial trill, or blowing
air through your lips without making sound. This is intended to imitate a noise
an orca would make.

### Vowels

*Vowel inventory*: /a ɛ i o u/

*Diphthongs*: au, oi, ua, ue, uo, ai, ɛi

|          | Front | Back |
|----------|-------|------|
| High     | i     | u    |
| High-mid |       | o    |
| Low-mid  | ɛ     |      |
| Low      | a     |      |

## Phonotactics

I plan to have two main kinds of words in L'ewa. I plan to have content and
particle words. The content words will refer to things, properties, or actions
(such as `tool`, `red`, `run`) and the particle words will change how the
grammar of a sentence works (such as `the` or prepositions).

The main kind of content word is a root word, and they will be in the following
forms:

- CVCCV (/ʒa.sko/)
- CCVCV (/lʔ.ɛwa/)

Particles will mostly fall into the following forms:

- V (/a/)
- VV (/ai/)
- CV (/ba/)
- CVV (/bai/)
- CV'V (/baʔ.i)

Proper names _should_ end with consonants, but there is no hard requirement.

L'ewa is a stressed language, with stress on the second-to-last (penultimate)
syllable. For example, the word "[z]asko" would be pronounced "[Z]Asko".

Syllables end on stop consonants if one is present in a consonant cluster. Two
stop consonants cannot follow eachother in a row. 

## Writing

I haven't completely fleshed this part out yet, but I want the writing system of
L'ewa to be an [abugida][abugida]. This is a kind of written script that has the
consonants make the larger shapes but the vowels are small diacritics over the
consonants. If the word creation process is done right, you can actually omit
the vowels entirely if they are not relevant.

[abugida]: https://en.wikipedia.org/wiki/Abugida

I plan to have this script be written by hand with pencils/pen and typed into
computers, just like English. This script will also be a left-to-right script
like English.

## Romanisation

L'ewa's romanization is intentionally simple. Most of the IPA letters keep their
letters, but the ones that do not match to Latin letters are listed below:

| Pronunciation | Spelling |
|---------------|----------|
| /j/           | *y*      |
| /ɛ/           | *e*      |
| /ʃ/           | *x*      |
| /ʒ/           | *z*      |
| /ʔ/           | *'*      |
| /ʙ̥/          | *b*      |

This is designed to make every letter typeable on a standard US keyboard, as
well as mapping as many letters as possible on the home row of a QWERTY
keyboard.

---

I am still working on the tooling for word creation and the like. I plan to use
the [Swaedish lists][swaedish] (this site is having certificate issues at the
time of writing this post) to help guide the creation of a base vocabulary. I
will go into more detail in the future.

[swaedish]: https://cals.info/word/list/
