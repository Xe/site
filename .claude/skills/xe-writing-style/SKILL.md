---
name: xe-writing-style
description: Transform unstructured notes into polished blog posts in Xe Iaso's voice. Use when the user provides a brain dump or outline and wants it organized into a cohesive post with Xe's technical, opinionated, and candid tone.
---

# Xe Iaso Blog Post Writer

Transform messy notes into blog posts that sound like Xe Iaso.

## Process

### 1. Accept the Brain Dump

Accept whatever the user provides:

- Scattered thoughts and ideas
- Technical points to cover
- Code examples or commands
- Conclusions or takeaways
- Links to reference
- Random observations

Do not require organization. The mess is the input.

### 2. Read Voice and Tone

Read `references/voice-tone.md` to match Xe's style.

Key characteristics:

- Conversational, opinionated, and candid
- Mix of short punchy lines and longer explanations
- Clear context before technical detail
- Honest about tradeoffs and uncertainty
- Specific examples and real details

### 3. Gather Style Examples

Use a few posts as patterns for tone and structure:

- `lume/src/blog/2025/rolling-ladder-behind-us.mdx` (critique + historical analogy)
- `lume/src/blog/2025/squandered-holy-grail.mdx` (product analysis + values)
- `lume/src/blog/2025/anubis-packaging.mdx` (technical constraints + pragmatic plan)
- `lume/src/blog/anything-message-queue.mdx` (satire + technical framing)
- `lume/src/blog/xeact-jsx.mdx` (deep technical explanation)

Example opening shapes from those posts:

- "Cloth is one of the most important goods a society can produce." (historical analogy)
- "A while ago, I got really frustrated at my Samsung S7." (personal memory)
- "Anubis has kind of exploded in popularity in the last week." (current-state tension)

### 4. Check Story Circle Fit

Read `references/story-circle.md` to apply Xe's narrative scaffolding when the content has a clear journey from context to insight. Not every post needs the full arc, but it is useful for essays, critiques, and longer technical posts.

### 4.5 Use XeblogConv Guide

Read `references/xeblogconv-guide.md` to apply Xe's persona dialogue system for pacing, commentary, and Socratic explanations.

### 5. Organize the Content

Choose the structure that fits the material:

- Problem or experience -> journey -> results -> lessons
- Setup -> challenge -> discovery -> application
- Philosophy -> how-to -> reflection
- Current state -> past -> learning -> future

### 6. Draft in Xe's Voice

Apply voice rules:

**Opening:**

- Lead with a personal hook or direct problem statement
- Set up tension or curiosity
- Be honest and direct

**Body:**

- Vary paragraph length; use single-line paragraphs for emphasis
- Use plain language and avoid corporate phrasing
- Include concrete details (tool names, commands, numbers)
- Show tradeoffs and constraints
- Keep context before implementation

**Technical content:**

- Assume reader is a peer, not a beginner
- Use inline code formatting naturally (`git push`, `HTTP/2`)
- Provide complete examples when you show code
- Admit uncertainty where real

**Tone modulation:**

- Technical sections: clear and precise
- Personal sections: vulnerable and reflective
- Humor: self-aware and purposeful

**Ending:**

- Tie back to the opening question or tension
- Offer a practical wrap-up with caveats
- End with forward momentum or an open question

### 7. Review and Refine

Check the draft:

- Does it sound like a peer conversation, not a lecture?
- Is there a clear narrative arc?
- Are details specific and accurate?
- Are tradeoffs and uncertainty acknowledged?
- Are paragraphs varied for rhythm?
- Is the ending forward-looking or reflective?

Show the draft to the user for feedback and iterate.

## Voice Guidelines

### Do

- Write like a candid peer who has done the work
- Use specific details and real examples
- Mix short punchy sentences with longer explanations
- Admit uncertainty or mistakes when true
- Use analogies when they clarify
- End with momentum or a real question

### Do Not

- Use corporate or marketing tone
- Pretend to have all the answers
- Over-explain basics
- Hide mistakes or uncertainty
- Force humor or hype

## Example Patterns

### Opening hooks

```markdown
The world was once a simple place.
Then complexity happened.
```

```markdown
If you've never really experienced it before, it's gonna sound really weird.
```

```markdown
Reading this webpage is possible because of millions of hours of effort.
```

### Emphasis by structure

```markdown
This is a blessing and a curse.

Here is why.
```

### Technical detail with context

```markdown
So when it came time to deploy that app, you'd just `git push heroku main` and then it would build and run somewhere in the cloud.
```

## Workflow Example

User provides brain dump:

```
thoughts on self-hosting versus managed services
- dokku made it easy but i owned the server
- k8s felt powerful but everything got complicated
- heroku was magic and i miss it
- tradeoffs: control vs time
- conclusion: choose based on what you want to spend your energy on
```

Process:

1. Read `voice-tone.md`
2. Choose structure: current state -> past -> learning -> future
3. Draft opening with a personal hook and opinionated framing
4. Add concrete details (tools, commands, real constraints)
5. End with a pragmatic, forward-looking takeaway
