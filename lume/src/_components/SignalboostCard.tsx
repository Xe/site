export type Link = {
  url: string;
  title: string;
  description?: string;
};

export interface SignalBoostCardProps {
  name: string;
  tags: string[];
  links: Link[];
}

export default function Card(
  {
    name,
    tags,
    links,
  }: SignalBoostCardProps,
) {
  return (
    <div id={name} class="bg-bg-1 dark:bg-bgDark-1 rounded-xl m-2 shadow-md max-w-[29rem]">
      <div class="items-center text-xl pl-4 pt-4 font-bold">
        {name}
      </div>
      <div class="flex flex-row items-center px-4 py-2">
        <div class="flex flex-wrap items-start justify-center p-5">
          {tags.map((tag) => (
            <div class="px-2 py-1 m-1 bg-bg-2 dark:bg-bgDark-2 rounded-lg">
              {tag}
            </div>
          ))}
          {links.map((link) => (
            <div class="px-2 py-1 m-1 bg-bg-2 dark:bg-bgDark-2 rounded-lg">
              <a class="flex flex-row items-center" href={link.url}>
                {link.title}
              </a>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
