export interface XeblogPicture {
  path: string;
  desc?: string;
}

export default function XeblogPicture({ path, desc }: XeblogPicture) {
  return (
    <figure>
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
          alt={`An image of ${prompt}`}
          loading="lazy"
          src={`https://cdn.xeiaso.net/file/christine-static/${path}.jpg`}
        />
      </picture>
      {desc && <figcaption>{desc}</figcaption>}
    </figure>
  );
}
