export interface XeblogHeroProps {
  ai: string;
  file: string;
  prompt: string;
}

export default function XeblogHero({ ai, file, prompt }: XeblogHeroProps) {
  return (
    <>
      <figure className="hero not-prose w-full">
        <picture>
          <source
            type="image/avif"
            srcset={`https://cdn.xeiaso.net/file/christine-static/hero/${file}.avif`}
          />
          <source
            type="image/webp"
            srcset={`https://cdn.xeiaso.net/file/christine-static/hero/${file}.webp`}
          />
          <img
            alt={`An image of ${prompt}`}
            className="hero-image"
            loading="lazy"
            src={`https://cdn.xeiaso.net/file/christine-static/hero/${file}.jpg`}
          />
        </picture>
        {ai !== undefined ? (
          <></>
        ) : (
          <figcaption className="mx-2 my-1 text-center text-gray-600 dark:text-gray-400">
            {prompt}
          </figcaption>
        )}
      </figure>
    </>
  );
}
