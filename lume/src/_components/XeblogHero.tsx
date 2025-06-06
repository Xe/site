export interface XeblogHeroProps {
  ai: string;
  file: string;
  prompt: string;
}

export default function XeblogHero({ ai, file, prompt }: XeblogHeroProps) {
  return (
    <>
      <figure className="hero not-prose w-full mx-auto">
        <picture>
          <source
            type="image/avif"
            srcset={`https://files.xeiaso.net/hero/${file}.avif`}
          />
          <source
            type="image/webp"
            srcset={`https://files.xeiaso.net/hero/${file}.webp`}
          />
          <img
            alt={`An image of ${prompt}`}
            className="hero-image"
            loading="lazy"
            src={`https://files.xeiaso.net/hero/${file}.jpg`}
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
