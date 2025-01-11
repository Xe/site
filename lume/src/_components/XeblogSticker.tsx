export interface XeblogStickerProps {
  name: string;
  mood: string;
  maxHeight?: string | null;
}

export default function XeblogSticker({
  name,
  mood,
  maxHeight = null,
}: XeblogStickerProps) {
  const nameLower = name.toLowerCase();
  name = name.replace(" ", "_");
  if (maxHeight === null) {
    maxHeight = "12rem";
  }

  return (
    <>
      <img
        style={`max-height:${maxHeight}`}
        alt={`${name} is ${mood}`}
        loading="lazy"
        src={`https://stickers.xeiaso.net/sticker/${nameLower}/${mood}`}
      />
    </>
  );
}
