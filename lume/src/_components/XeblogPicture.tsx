export interface XeblogPicture {
  path: string;
  desc?: string;
}

export default function XeblogPicture({ path, desc }: XeblogPicture) {
  return (
    <figure className="max-w-3xl mx-auto">
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
            loading="lazy"
            src={`https://cdn.xeiaso.net/file/christine-static/${path}.jpg`}
          />
        </picture>
      </a>
      {desc && <figcaption>{desc}</figcaption>}
    </figure>
  );
}
