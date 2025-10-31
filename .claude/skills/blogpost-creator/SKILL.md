---
name: blogpost-creator
description: Create new posts using the hydrate command.
allowed-tools: Read, Grep, Glob, WebFetch
---

# Blog Post Creator

Create new blog posts for Xe's site using the hydrate command.

## Usage

Invoke this skill when you want to create a new blog post. The skill will prompt you for:

- **Post type**: blog, linkpost, note, talk, or xecast
- **Title**: The title of your post (for linkposts, can auto-fetch from URL)
- **Description**: A brief description/summary (for linkposts, can auto-fetch from URL)
- **Link URL**: (required for linkposts) The URL this post should redirect to
- **Publication date**: (optional, for linkposts can auto-extract from the page, defaults to today if not found)
- **Custom slug**: (optional, for linkposts will auto-extract from URL basename if not provided)

## What it does

1. **For linkposts**: Extracts slug from URL basename if no custom slug provided (e.g., https://example.com/blog/my-post becomes "my-post")
2. **For other post types**: Generates a URL-friendly slug from your title (or uses your custom slug)
3. **For linkposts**: Can auto-fetch title, description, and publication date from the provided URL
4. Validates that linkposts have a URL provided (prompts if missing)
5. Runs `go run ./cmd/hydrate <kind> <slug>` with the appropriate parameters
6. For linkposts: Updates the `redirect_to` field in the frontmatter with the provided URL
7. **For linkposts**: Adds fetched summary to the post body if available
8. Opens the created file in VS Code for editing
9. Shows you the file location for reference

## File structure

Blog posts are created in:

- `lume/src/blog/<year>/<slug>.mdx` for blog and linkpost
- `lume/src/notes/<year>/<slug>.mdx` for notes
- `lume/src/talks/<year>/<slug>.mdx` for talks
- `lume/src/xecast/<year>/<slug>.mdx` for xecast

## Frontmatter templates

Each post type has its own frontmatter template:

**Blog posts** include hero image configuration:

```yaml
---
title: ""
desc: ""
date: YYYY-MM-DD
hero:
  ai: ""
  file: ""
  prompt: ""
  social: false
---
```

**Link posts** include a redirect URL:

```yaml
---
title: ""
date: YYYY-MM-DD
redirect_to: "https://example.com"
---
```

**Other types** (notes, talks, xecast) have simpler frontmatter:

```yaml
---
title: ""
desc: ""
date: YYYY-MM-DD
---
```

## Linkpost Special Features

Linkposts have enhanced functionality:

1. **Automatic slug extraction**: If no custom slug is provided, the skill will extract the basename from the URL (e.g., `https://example.com/blog/my-post` becomes `my-post`)

2. **Auto-fetching content**: The skill can automatically fetch the webpage to extract:

   - The page title (used as the post title)
   - A summary/description (added to the post body)
   - The publication date (used as the post date, defaults to today if not found)
   - This saves time and ensures accurate representation of the linked content

Use the extract-meta.js file in this folder to extract meta-information from webpages:

```bash
node extract-meta.js <url>
```

3. **URL handling**: The skill handles redirects and will follow them to get the final content for title/description extraction

## Example Linkpost Workflow

When creating a linkpost with a URL like `https://anubis.techaro.lol/blog/2025/file-abuse-reports`:

1. Skill detects it's a linkpost with a URL
2. Extracts slug "file-abuse-reports" from URL basename
3. Fetches the webpage to get:
   - The actual title: "Taking steps to end traffic from abusive cloud providers"
   - A summary of the content for the post body
   - The publication date from the page (e.g., "2025-01-15")
4. Creates the post with auto-generated slug, fetched title, and extracted date
5. Adds the summary to the post body for context
6. Updates the redirect_to field with the provided URL

## Date Extraction Details

The skill will look for publication dates in various formats:

- **Meta tags**: `<meta property="article:published_time" content="2025-01-15">`
- **JSON-LD structured data**: `"datePublished": "2025-01-15"`
- **HTML5 semantic elements**: `<time datetime="2025-01-15">`
- **Common date patterns in the page content**
- **URL patterns**: Extracts date from URL structure like `/blog/2025/my-post`

If no date is found, it defaults to today's date.
