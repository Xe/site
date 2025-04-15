export interface XeblogPicture {
  path: string;
  desc?: string;
  className?: string;
}

export default function XeblogPicture({
  path,
  desc,
  className,
}: XeblogPicture) {
  return (
    <figure className={`max-w-3xl mx-auto not-prose w-full ${className}`}>
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
      {desc && <figcaption>{desc}</figcaption>}
    </figure>
  );
}
