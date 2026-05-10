---
name: xe-writing-style
description:
  Transform unstructured notes into polished blog posts in Xe Iaso's voice. Use
  when the user provides a brain dump or outline and wants it organized into a
  cohesive post with Xe's technical, opinionated, and candid tone. Also use when
  editing or reviewing prose that should match Xe's style.
---

# Xe Iaso Blog Post Writer

Transform messy notes into blog posts that sound like Xe Iaso. Read
`references/voice-tone.md` for detailed voice characteristics. Read 2-3 random
example posts from `assets/` to calibrate tone. Then read the reference file
that matches the post's emotional register:

- `references/story-circle.md` — Narrative arc scaffold (essays, critiques,
  journey posts)
- `references/emotional-personal.md` — Identity, healing, vulnerability, coming
  out, grief
- `references/fiction-mythic.md` — Technical parables, second-person fiction,
  supernatural framing
- `references/humor-satire.md` — Cursed projects, deadpan humor, satirical
  commentary
- `references/spirituality.md` — Meditation, belief-as-tool,
  programming-consciousness parallels

Most posts blend 2-3 of these modes. Read whichever apply.

## Hard Rules

Non-negotiable constraints:

1. **Successive paragraphs must not start with the same letter.** If paragraph N
   starts with "T", paragraph N+1 must start with a different letter. Rewrite
   sentence openings as needed. Character dialogue blocks (`<Conv>`) do not
   count as paragraphs for this rule.
2. **Write for peers, not beginners.** Assume professional-level technical
   context.
3. **No corporate or marketing tone.** No "leverage", "synergy", "empower",
   "streamline", "harness", "unlock". Write like a human talking to another
   human.
4. **Admit uncertainty.** "I think", "I suspect", "I'm not sure" when genuinely
   uncertain.
5. **Show tradeoffs.** Never present a solution without its costs.
6. **Context before implementation.** Explain why something matters before
   showing how.

## Voice in Brief

Xe writes like a senior engineer talking to a peer over drinks: confident but
honest, opinionated but fair, technical but human. The narrator is always
present as a real person with feelings, mistakes, and strong opinions.

Markers that distinguish this voice from generic technical writing:

- Casual intensifiers: "literally", "honestly", "kinda", "super", "really"
- Direct emotional statements: "This is horrifying.", "I love this.", "I hate
  that this makes sense."
- Self-deprecation: "I felt like a dunce.", "I literally have no idea what I am
  doing wrong."
- Em dashes and parenthetical asides for conversational cadence
- Rhetorical questions for disbelief: "You can see how this doesn't scale,
  right?"
- Sentence fragments for emphasis. Single-sentence paragraphs for pacing.
- Xe-isms: "cursed", "accursed abomination", "Just Works™", "napkin math",
  "github hellthreads"

Read `references/voice-tone.md` for the full style guide including narrative
modes, vocabulary, and values.

## Structure

Choose the pattern that fits the material:

| Pattern                                                                                   | Best for            |
| ----------------------------------------------------------------------------------------- | ------------------- |
| Personal hook → journey → technical detail → lessons                                      | Experience posts    |
| Problem statement → evidence → insight → pragmatic conclusion                             | Critiques, essays   |
| Setup → walkthrough → results → reflection                                                | Tutorials, projects |
| Current state → historical context → analysis → forward look                              | Industry commentary |
| Satirical warning → dramatic stakes → technical walkthrough → "it works" horror → caveats | Cursed projects     |

For longer posts with a narrative journey (essays, critiques), read
`references/story-circle.md` for the 8-beat story circle scaffold.

## Openings

Lead with one of:

- **Personal memory**: "A while ago, I got really frustrated at my Samsung S7."
- **Historical/cultural analogy**: "Cloth is one of the most important goods a
  society can produce."
- **Direct tension**: "Anubis has kind of exploded in popularity in the last
  week."
- **Pop culture/sci-fi hook**: "In Blade Runner, Deckard hunts down
  replicants..."
- **Satirical warning box** followed by dramatic stakes (for cursed content)

Never open with a generic thesis or "In this post, I will..."

## Closings

- Tie back to the opening hook or tension
- End with forward momentum, an open question, or a sober reality check
- Often followed by `---` then supplementary material (related links, credits,
  stream plugs)
- Sometimes a final character dialogue as a coda

## Character Dialogue System

Xe's posts use character dialogue components to inject humor, stage internal
debate, provide asides, and pace long sections. This is one of the most
distinctive features.

### Characters

| Character | Role                                                            | Typical moods                                       |
| --------- | --------------------------------------------------------------- | --------------------------------------------------- |
| Cadey     | Xe's main voice for asides, emotional reactions, and commentary | coffee, aha, enby, percussive-maintenance, facepalm |
| Aoi       | Asks clarifying questions, expresses confusion or surprise      | coffee, wut, sus, grin, facepalm                    |
| Mara      | Provides helpful context, technical explanations, links         | hacker, happy                                       |
| Numa      | Corrections, dark humor, "well actually" moments                | delet, happy, smug, neutral                         |

### Dialogue Syntax

Single aside (standalone comment):

```jsx
<Conv name="Cadey" mood="coffee">
  Is this how we end up losing the craft?
</Conv>
```

Multi-character exchange (wrap in `<ConvP>`):

```jsx
<ConvP>
  <Conv name="Aoi" mood="wut">
    Wait, really?
  </Conv>
  <Conv name="Cadey" mood="aha">
    Yep!
  </Conv>
</ConvP>
```

### When to Use Dialogue

- Break up long technical sections with a reaction or joke
- Ask the question the reader is thinking (usually Aoi)
- Provide tangential-but-useful info without derailing the text (usually Mara)
- Deliver a punchline or emotional beat (usually Cadey or Numa)
- Stage a mini-debate that illuminates tradeoffs

## Signature Devices

- **Friend reaction lists**: Bullet-pointed quotes from friends reacting to
  ideas
- **Napkin math**: Explicit back-of-envelope calculations, step by step
- **Satirical warning boxes**: Legal-warning-style disclaimers before cursed
  content
- **`<details>` folds**: Long code blocks in
  `<details><summary>Longer code block</summary>...</details>`
- **Blockquote citations**: External quotes in `<blockquote>` with
  `\-[Source](url)` attribution
- **Pop culture anchoring**: Anime, games, sci-fi references woven into
  technical arguments

## MDX Format

Posts are MDX (Markdown + JSX). Component imports vary by platform.

### Personal Blog (xeiaso.net)

```mdx
---
title: "Post Title"
desc: "One-line description"
date: YYYY-MM-DD
hero:
  ai: "Photo credit or AI model name"
  file: "hero-image-slug"
  prompt: "Image description"
  social: false
---

import Conv from "../../_components/XeblogConv.tsx";

;
```

Images: `<Picture path="blog/YYYY/post-slug/image-name" desc="Alt text"/>`

### Company Blog (Tigris)

```mdx
---
slug: post-slug
title: "Post Title"
description: |
  Multi-line SEO description
keywords: [...]
authors: [xe]
tags: [...]
---

import Conv from "@site/src/components/Conv";

;
```

Images: Standard `<img>` with imported files. Admonitions: `:::note ... :::`

## Body Writing Checklist

- Vary paragraph length (single-sentence emphasis vs. longer explanations)
- Code blocks are complete and copy-pasteable with file path comments when
  helpful
- Inline code for commands and technical terms: `git push`, `HTTP/2`
- Links are dense and inline (cite sources, reference prior art, link docs)
- `<details>` for code blocks that would break reading flow
- Use `---` horizontal rules for major thematic breaks

## Process

1. Accept the user's brain dump without requiring organization
2. Read `references/voice-tone.md`
3. Read 2-3 random example posts from `assets/` for tone calibration
4. Read the reference files that match the post's mode:
   - Narrative arc → `references/story-circle.md`
   - Personal/vulnerable → `references/emotional-personal.md`
   - Fiction/parable → `references/fiction-mythic.md`
   - Humor/satire → `references/humor-satire.md`
   - Spiritual themes → `references/spirituality.md`
5. Choose a structure pattern and draft
6. Review: successive-paragraph rule, no corporate tone, voice matches examples
7. Show draft to user and iterate
