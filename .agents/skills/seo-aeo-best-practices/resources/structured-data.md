# Structured Data (JSON-LD)

Structured data helps search engines and AI understand your content. JSON-LD is
the recommended format.

## Why Structured Data Matters

- **Rich snippets:** Enhanced search result appearance
- **Knowledge panels:** Featured information boxes
- **AI training:** Better content understanding
- **Voice search:** Answer selection for voice queries

## Common Schema Types

### Article / Blog Post

```typescript
import { Article, WithContext } from "schema-dts";

const articleSchema: WithContext<Article> = {
  "@context": "https://schema.org",
  "@type": "Article",
  headline: post.title,
  description: post.excerpt,
  image: post.image?.url,
  datePublished: post.publishedAt,
  dateModified: post.updatedAt,
  author: {
    "@type": "Person",
    name: post.author.name,
    url: post.author.url,
  },
  publisher: {
    "@type": "Organization",
    name: "Your Company",
    logo: {
      "@type": "ImageObject",
      url: "https://example.com/logo.png",
    },
  },
};
```

### FAQ Page

```typescript
import { FAQPage, WithContext } from "schema-dts";

const faqSchema: WithContext<FAQPage> = {
  "@context": "https://schema.org",
  "@type": "FAQPage",
  mainEntity: faqs.map((faq) => ({
    "@type": "Question",
    name: faq.question,
    acceptedAnswer: {
      "@type": "Answer",
      text: faq.answer, // Plain text, use pt::text() in GROQ
    },
  })),
};
```

### Organization

```typescript
import { Organization, WithContext } from "schema-dts";

const orgSchema: WithContext<Organization> = {
  "@context": "https://schema.org",
  "@type": "Organization",
  name: "Your Company",
  url: "https://example.com",
  logo: "https://example.com/logo.png",
  sameAs: [
    "https://twitter.com/company",
    "https://linkedin.com/company/company",
  ],
  contactPoint: {
    "@type": "ContactPoint",
    telephone: "+1-555-555-5555",
    contactType: "customer service",
  },
};
```

### Product

```typescript
import { Product, WithContext } from "schema-dts";

const productSchema: WithContext<Product> = {
  "@context": "https://schema.org",
  "@type": "Product",
  name: product.name,
  description: product.description,
  image: product.images,
  offers: {
    "@type": "Offer",
    price: product.price,
    priceCurrency: "USD",
    availability: "https://schema.org/InStock",
  },
  aggregateRating: product.rating
    ? {
        "@type": "AggregateRating",
        ratingValue: product.rating.average,
        reviewCount: product.rating.count,
      }
    : undefined,
};
```

## Implementation in Next.js

```typescript
// Component to render JSON-LD
function JsonLd({ data }: { data: WithContext<Thing> }) {
  return (
    <script
      type="application/ld+json"
      dangerouslySetInnerHTML={{ __html: JSON.stringify(data) }}
    />
  )
}

// Usage in page
export default function PostPage({ post }) {
  return (
    <>
      <JsonLd data={generateArticleSchema(post)} />
      <article>...</article>
    </>
  )
}
```

## GROQ for Plain Text

Structured data often needs plain text, not rich text:

```groq
*[_type == "faq"]{
  question,
  "answer": pt::text(answerRichText)  // Convert Portable Text to plain string
}
```

## Testing Tools

- [Google Rich Results Test](https://search.google.com/test/rich-results)
- [Schema.org Validator](https://validator.schema.org/)
