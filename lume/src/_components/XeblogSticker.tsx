export interface XeblogStickerProps {
  name: string;
  mood: string;
}

export default function XeblogSticker({ name, mood }: XeblogStickerProps) {
  const nameLower = name.toLowerCase();
  name = name.replace(" ", "_");

  return (
    <>
      <picture>
        <source
          type="image/avif"
          srcset={`https://cdn.xeiaso.net/file/christine-static/stickers/${nameLower}/${mood}.avif`}
        />
        <source
          type="image/webp"
          srcset={`https://cdn.xeiaso.net/file/christine-static/stickers/${nameLower}/${mood}.webp`}
        />
        <img
          alt={`${name} is ${mood}`}
          loading="lazy"
          src={`https://cdn.xeiaso.net/file/christine-static/stickers/${nameLower}/${mood}.png`}
        />
      </picture>
    </>
  );
}
