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
      <a href={`https://cdn.xeiaso.net/file/christine-static/${path}.jpg`}>
        <picture>
          <source
            type="image/avif"
            srcset={`https://cdn.xeiaso.net/file/christine-static/${path}.avif`}
          />
          <source
            type="image/webp"
            srcset={`https://cdn.xeiaso.net/file/christine-static/${path}.webp`}
          />
          <img
            alt={desc}
            className={className}
            loading="lazy"
            src={`https://cdn.xeiaso.net/file/christine-static/${path}.jpg`}
          />
        </picture>
      </a>
      {desc && <figcaption>{desc}</figcaption>}
    </figure>
  );
}
