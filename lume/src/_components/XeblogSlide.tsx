export interface XeblogSlideProps {
  name: string;
  essential?: boolean;
  desc?: string;
}

export default function XeblogSlide({ name, essential, desc }: XeblogSlideProps) {
  return (
    <figure class={essential ? "xeblog-sides-essential" : "xeblog-slides-fluff"}>
      <picture>
        <source
          type="image/avif"
          srcset={`https://cdn.xeiaso.net/file/christine-static/talks/${name}.avif`}
        />
        <source
          type="image/webp"
          srcset={`https://cdn.xeiaso.net/file/christine-static/talks/${name}.webp`}
        />
        <img
          alt={desc || `Slide ${name}`}
          loading="lazy"
          src={`https://cdn.xeiaso.net/file/christine-static/talks/${name}.jpg`}
        />
      </picture>
      {desc && <figcaption>{desc}</figcaption>}
    </figure>
  );
}
