export interface XeblogConvProps {
  name: string;
  mood: string;
  children: HTMLElement[];
  standalone?: boolean;
}

const ConvSnippet = ({ name, mood, children, standalone }: XeblogConvProps) => {
  const nameLower = name.toLowerCase();
  name = name.replace(" ", "_");
  const size = standalone ? 128 : 64;

  return (
    <>
      <div className="my-4 flex space-x-4 rounded-md border border-solid border-fg-4 bg-bg-2 p-3 dark:border-fgDark-4 dark:bg-bgDark-2 max-w-full">
        <div className="flex max-h-16 max-w-16 shrink-0 items-center justify-center self-center overflow-hidden rounded-lg bg-gray-200 dark:bg-gray-700">
            <img
              style="max-height:4.5rem"
              alt={`${name} is ${mood}`}
              loading="lazy"
              src={`https://cdn.xeiaso.net/sticker/${nameLower}/${mood}/${size}`}
            />
        </div>
        <div className="convsnippet min-w-0 self-center">
          {"<"}
          <a href={`/characters#${nameLower}`}>
            <b>{name}</b>
          </a>
          {">"} {children}
        </div>
      </div>
    </>
  );
};

export default ConvSnippet;
