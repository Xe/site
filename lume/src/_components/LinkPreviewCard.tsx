export interface LinkPreviewCardProps {
  url: string;
  title: string;
  description?: string;
  /** Optional preview/OG image URL. When omitted the card renders text-only. */
  image?: string;
  /** Override the source label. Defaults to the URL's hostname. */
  siteName?: string;
}

export default function LinkPreviewCard({
  url,
  title,
  description,
  image,
  siteName,
}: LinkPreviewCardProps) {
  const host = siteName ?? new URL(url).hostname.replace(/^www\./, "");

  return (
    <a
      href={url}
      target="_blank"
      rel="noopener noreferrer"
      class="group not-prose no-underline visited:no-underline flex flex-col sm:flex-row max-w-xl mx-auto my-8 overflow-hidden rounded-xl border border-bg-3 dark:border-bgDark-3 bg-bg-0 dark:bg-bgDark-0 shadow-sm transition duration-150 hover:-translate-y-0.5 hover:shadow-md"
    >
      {image && (
        <div class="sm:w-52 sm:shrink-0 aspect-[1.91/1] sm:aspect-auto overflow-hidden bg-bg-2 dark:bg-bgDark-2">
          <img
            src={image}
            alt=""
            loading="lazy"
            class="h-full w-full object-cover transition-transform duration-200 group-hover:scale-105"
          />
        </div>
      )}
      <div class="flex flex-col justify-center gap-1 p-4">
        <span class="font-mono text-xs uppercase tracking-wide text-orange-light dark:text-orangeDark-light group-hover:text-bg-0 dark:group-hover:text-white">
          {host}
        </span>
        <span class="font-serif text-lg font-semibold leading-snug text-fg-0 dark:text-fgDark-0 group-hover:text-bg-0 dark:group-hover:text-white">
          {title}
        </span>
        {description && (
          <span class="text-sm leading-snug text-fg-3 dark:text-fgDark-3 line-clamp-3 group-hover:text-bg-1 dark:group-hover:text-bg-1">
            {description}
          </span>
        )}
      </div>
    </a>
  );
}
