# Technical SEO Checklist

Essential technical SEO elements for modern web applications.

## Metadata

### Title Tags

- Unique per page
- 50-60 characters
- Primary keyword near the beginning
- Brand name at the end (optional)

### Meta Descriptions

- Unique per page
- 150-160 characters
- Include call-to-action
- Contain relevant keywords

### Open Graph

```html
<meta property="og:title" content="Page Title" />
<meta property="og:description" content="Description" />
<meta property="og:image" content="https://example.com/image.jpg" />
<meta property="og:url" content="https://example.com/page" />
<meta property="og:type" content="article" />
```

### Sanity + Next.js Implementation

```typescript
export async function generateMetadata({ params }): Promise<Metadata> {
  const { data } = await sanityFetch({
    query: PAGE_QUERY,
    stega: false, // Critical: no stega in metadata
  });

  return {
    title: data.seo?.title || data.title,
    description: data.seo?.description,
    openGraph: {
      images: data.seo?.image
        ? [
            {
              url: urlFor(data.seo.image).width(1200).height(630).url(),
              width: 1200,
              height: 630,
            },
          ]
        : [],
    },
    robots: data.seo?.noIndex ? "noindex" : undefined,
  };
}
```

## Sitemaps

Dynamic sitemap from CMS content:

```typescript
// app/sitemap.ts
import { MetadataRoute } from "next";

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const pages = await client.fetch(`
    *[_type in ["page", "post"] && defined(slug.current) && seo.noIndex != true]{
      "url": select(
        _type == "page" => "/" + slug.current,
        _type == "post" => "/blog/" + slug.current
      ),
      _updatedAt
    }
  `);

  return pages.map((page) => ({
    url: `https://example.com${page.url}`,
    lastModified: new Date(page._updatedAt),
    changeFrequency: "weekly",
    priority: 0.8,
  }));
}
```

## Canonical URLs

Prevent duplicate content issues:

```typescript
export async function generateMetadata({ params }): Promise<Metadata> {
  return {
    alternates: {
      canonical: `https://example.com/${params.slug}`,
    },
  };
}
```

## Redirects

CMS-managed redirects:

```typescript
// next.config.ts
async redirects() {
  const redirects = await client.fetch(`
    *[_type == "redirect" && isEnabled == true]{
      source,
      destination,
      permanent
    }
  `)
  return redirects
}
```

## Performance

[Core Web Vitals](https://web.dev/articles/defining-core-web-vitals-thresholds)
impact rankings:

- **LCP (Largest Contentful Paint):** < 2.5s
- **INP (Interaction to Next Paint):** < 200ms
- **CLS (Cumulative Layout Shift):** < 0.1

### Image Optimization (Next.js example)

- Use `next/image` with Sanity URL builder
- Serve WebP/AVIF formats
- Implement LQIP blur placeholders
- Set explicit dimensions

### Font Loading (Next.js example)

```typescript
// Prevent layout shift
import { Inter } from "next/font/google";
const inter = Inter({ subsets: ["latin"], display: "swap" });
```

## Robots.txt

```
# public/robots.txt
User-agent: *
Allow: /
Disallow: /api/
Disallow: /studio/

Sitemap: https://example.com/sitemap.xml
```
