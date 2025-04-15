export interface XeblogSlideProps {
  name: string;
  essential?: boolean;
  desc?: string;
}

export default function XeblogSlide({
  name,
  essential,
  desc,
}: XeblogSlideProps) {
  return (
    <figure
      class={essential ? "xeblog-sides-essential" : "xeblog-slides-fluff"}
    >
      <picture>
        <source
          type="image/avif"
          srcset={`https://files.xeiaso.net/talks/${name}.avif`}
        />
        <source
          type="image/webp"
          srcset={`https://files.xeiaso.net/talks/${name}.webp`}
        />
        <img
          alt={desc || `Slide ${name}`}
          loading="lazy"
          src={`hhttps://files.xeiaso.net/talks/${name}.jpg`}
        />
      </picture>
      {desc && <figcaption>{desc}</figcaption>}
    </figure>
  );
}
