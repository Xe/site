export interface XeblogPicture {
  path: string;
  desc?: string;
  className?: string;
  fullBleed?: boolean;
}

export default function XeblogPicture({
  path,
  desc,
  className,
  fullBleed,
}: XeblogPicture) {
  const figClass = fullBleed
    ? `full-bleed not-prose my-8 ${className ?? ""}`
    : `max-w-3xl mx-auto not-prose w-full ${className ?? ""}`;

  return (
    <figure className={figClass}>
      <a href={`https://files.xeiaso.net/${path}.jpg`}>
        <picture>
          <source
            type="image/avif"
            srcset={`https://files.xeiaso.net/${path}.avif`}
          />
          <source
            type="image/webp"
            srcset={`https://files.xeiaso.net/${path}.webp`}
          />
          <img
            alt={desc}
            className={className}
            loading="lazy"
            src={`https://files.xeiaso.net/${path}.jpg`}
          />
        </picture>
      </a>
      {desc && <figcaption className="max-w-xl">{desc}</figcaption>}
    </figure>
  );
}
