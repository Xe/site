// @jsxImportSource xeact
// @jsxRuntime automatic

export interface ConvSnippetProps {
  name: string;
  mood: string;
  children: HTMLElement[];
}

const ConvSnippet = ({ name, mood, children }: ConvSnippetProps) => {
  const nameLower = name.toLowerCase();
  name = name.replace(" ", "_");

  return (
    <div class="conversation">
      <div class="conversation-standalone">
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
            style="max-height:4.5rem"
            alt={`${name} is ${mood}`}
            loading="lazy"
            src={`https://cdn.xeiaso.net/file/christine-static/stickers/${nameLower}/${mood}.png`}
          />
        </picture>
      </div>
      <div className="conversation-chat">
        {"<"}
        <a href={`/characters#${nameLower}`}>
          <b>{name}</b>
        </a>
        {">"} {children}
      </div>
    </div>
  );
};

export default ConvSnippet;
