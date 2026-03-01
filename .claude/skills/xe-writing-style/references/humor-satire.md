# Humor & Satire Patterns

How Xe writes comedy, cursed projects, and satirical technical content.

## When to Read This

Read when the post involves: a deliberately absurd technical project, satirical
commentary, deadpan humor, "cursed" builds, or content that's meant to be funny
while still being technically rigorous.

## Core Comedic Voice

Xe's humor is deadpan expertise applied to absurd premises. The comedy comes
from treating something ridiculous with complete technical seriousness. The
author never winks at the camera or signals "this is a joke" — the gap between
the earnest delivery and the absurd subject IS the joke.

### The Cursed Project Formula

1. **Satirical warning** — A warning box or legal disclaimer that's obviously
   overblown for the content
2. **Dramatic stakes** — Frame the problem as if civilization depends on solving
   it
3. **Legitimate technical walkthrough** — Actually build the absurd thing with
   real engineering
4. **"It works" horror** — Express genuine surprise/dismay that the cursed thing
   functions
5. **Philosophical reflection** — End by questioning what it means that this is
   possible

Example (from "Anything can be a message queue"):

- Warning: legal disclaimer about not using IPv6 over S3 in production
- Stakes: AWS NAT Gateway costs framed as existential threat
- Walkthrough: Real S3 API calls, encoding schemes, polling architecture
- Horror: "This code legitimately works and I don't know how to feel about it"
- Reflection: What does it mean that cloud APIs are so flexible they enable
  abuse?

### Friend Reaction Lists

A signature device: bullet-pointed quotes from friends reacting to the absurd
idea, presented without commentary.

```markdown
- "You have entered the land of partially specified problems."
- "You need to be studied."
- "Did you just reinvent COBOL?"
- "I think something is either wrong with you, or wrong with me."
```

These serve dual purposes: social proof that the idea is funny AND providing the
reader permission to laugh.

### The Turing-Incomplete Bit

In "The h Programming Language" — an entire formal specification for a language
that only outputs the letter 'h'. Written with complete academic rigor: syntax,
semantics, implementation, and prior art. The humor is entirely structural: the
apparatus is real but the subject is trivially absurd.

Pattern: Take a format reserved for serious work (language spec, RFC, research
paper) and fill it with trivially absurd content.

## Satirical Commentary

When satire targets industry problems rather than just being funny for its own
sake:

### "Markdownlang" Pattern

1. Open with a cultural/literary reference that provides the moral framework
   (Blade Runner, replicants)
2. Transition to the real industry trend being satirized (AI replacing
   programmers)
3. Present the satirical creation with genuine technical detail
4. Let the reader sit with the discomfort of it actually working
5. Close by naming the real horror: not the technology, but its deployment

The satire is never mean-spirited toward individuals. It targets systems,
incentives, and the gap between what technology CAN do and what we SHOULD do
with it.

## Humor Mechanics

### Understatement

- "As you can imagine, the possibilities here are truly endless." (after showing
  FizzBuzz in markdown)
- "It is now ten times faster." (after changing one constant)

### Escalating Absurdity

Each section raises the stakes or adds another layer of wrongness, building on
the previous absurdity rather than resetting.

### Technical Precision as Comedy

The funniest parts are often the most technically accurate. Real error messages,
actual benchmarks, working code that does something ridiculous.

### Parenthetical Asides

Quick jokes delivered in parentheses that the reader might miss on first read:

- "(I'm assuming someone was inspired by my satirical post where I fixed the
  'strawberry' problem with AI models)"

### Character Dialogue for Punchlines

```jsx
<Conv name="Numa" mood="delet">
  This is why we can't have nice things.
</Conv>
```

Characters deliver reactions the author can't say in their own voice without
breaking the deadpan.

## What NOT to Do

- Don't signal that it's a joke. The delivery must be completely straight-faced.
- Don't sacrifice technical accuracy for humor. The real code must work.
- Don't punch down. Satire targets systems and powerful entities, not
  individuals or beginners.
- Don't force humor into posts that don't need it. Not every post is funny. The
  humor posts work BECAUSE other posts are sincere.
- Don't use "lol" or "haha" in satirical posts. The deadpan is sacred.
