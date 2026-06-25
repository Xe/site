---
name: no-way-to-prevent-this-memory-safety
description: Use when asked to generate a "no way to prevent this" / no-way-to-prevent-this satire post for a memory-safety CVE (out-of-bounds write/read, buffer/heap overflow, use-after-free, integer overflow) in a C or C++ project, given a CVE number or NVD link.
---

# Generate a "no way to prevent this" memory-safety post

## Overview

`cmd/no-way-to-prevent-this` is a Go generator that writes an Onion-style satire post
("'No way to prevent this' say users of only language where this regularly happens")
for a memory-safety CVE. It fills a template from CLI flags and writes the result to
`lume/src/shitposts/no-way-to-prevent-this/<template>/<CVE>.md`.

## Workflow

1. **Read the flags.** Run `go run ./cmd/no-way-to-prevent-this --help` to confirm the
   current flag set before composing the command.
2. **Look up the CVE.** Fetch the NVD page (or the link the user gave) to extract:
   - the affected **project** name
   - the **vulnerability type** and a one-line technical description (function, root cause)
   - whether the project is **C or C++** (C is the default; pass `-c++` only for C++)
3. **Find the real project homepage.** Web-search for the official homepage — do NOT
   guess or construct the URL. Use the canonical site (e.g. `https://libssh2.org/`),
   not a package mirror or the GitHub repo unless that is genuinely the home.
4. **Run the generator** with the flags filled in (see below).
5. **Read the generated file** to confirm it reads correctly, then report the flag
   values you used and the output path.

## Flags

| Flag | Value |
|------|-------|
| `-cve` | CVE id, e.g. `CVE-2026-55200` (also names the output file) |
| `-cve-link` | the NVD/advisory URL |
| `-project` | affected project name |
| `-project-link` | the **web-searched** official homepage |
| `-summary` | concise technical description of the flaw and its impact; flows into the sentence "...to fix `<summary>`." Write it to read naturally there. |
| `-c++` | add only if the project is C++ (omit for C — it is the default) |
| `-date` | defaults to today; override only if backdating |
| `-template` | defaults to `memory-safety`; use `supply-chain` for that variant |

## Example

```bash
go run ./cmd/no-way-to-prevent-this \
  -cve "CVE-2026-55200" \
  -cve-link "https://nvd.nist.gov/vuln/detail/CVE-2026-55200" \
  -project "libssh2" \
  -project-link "https://libssh2.org/" \
  -summary "an out-of-bounds write in ssh2_transport_read() due to a missing upper bound check on the packet_length field, resulting in heap corruption and potential remote code execution"
```

Writes `lume/src/shitposts/no-way-to-prevent-this/memory-safety/CVE-2026-55200.md`.

## Common mistakes

- **Guessing `-project-link`.** Always web-search for the homepage; a wrong/constructed URL ships a broken link.
- **Passing `-c++` for a C project.** C is the default; the flag changes the post's wording. Only set it when the affected code is actually C++.
- **A summary that doesn't fit the sentence.** It is appended after "to fix " — read it back in context so the grammar works.
- **Expecting stdout.** The command is silent on success; verify by checking the new file under `lume/src/shitposts/no-way-to-prevent-this/`.
