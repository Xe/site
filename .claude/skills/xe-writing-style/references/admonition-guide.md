# Admonition Component Guide

Use `<Admonition>` for short, high-signal callouts that should stand apart from the main narrative. Keep them concise and avoid interrupting dense technical explanation mid-thought.

## Import Required

When using `<Admonition>` in MDX files, you must import it at the top of your file after the frontmatter:

```jsx
---
title: "Your post title"
---
import Admonition from "../../_components/Admonition.jsx";

Your content here...
```

The path may vary depending on where your file is located in the directory structure.

## Usage

```jsx
<Admonition type="note">
  This is a short callout that adds context without derailing the flow.
</Admonition>
```

## Type Guidance

- `note`: Clarifications, context, or low-stakes caveats.
- `warning`: Safety, security, or risk-related guidance.
- `tip`: Practical advice or shortcuts.

## Tone Rules

- Keep it brief and direct.
- Avoid repeating the main text.
- Use when it helps pacing or emphasizes stakes.
