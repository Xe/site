# MDX/JSX Tag Checker

A simple CLI tool to detect unclosed or mismatched JSX/MDX tags in your files.

## Usage

```bash
# Check a single file
./scripts/check-mdx-tags path/to/file.mdx

# Check all MDX files in a directory
./scripts/check-mdx-tags lume/src/blog/

# Or run directly with Go
go run scripts/check-mdx-tags.go path/to/file.mdx
```

## What it detects

- Unclosed opening tags (like `<ConvP>` without `</ConvP>`)
- Mismatched closing tags (like `<ConvP>` closed with `</Conv>`)
- Unexpected closing tags (like `</ConvP>` without an opening tag)

## What it ignores

- Self-closing tags (like `<br/>` or `<img ... />`)
- HTML tags (only checks components starting with capital letters)
- Tags inside code blocks or comments

## Example output

```
lume/src/blog/2025/example.mdx:42: Unclosed opening tag <ConvP>
lume/src/blog/2025/example.mdx:58: Mismatched closing tag </Conv>, expected </ConvP> (opened at line 42)
lume/src/blog/2025/good-file.mdx: ✓ All tags properly closed
```
