import { crypto } from "@std/crypto";
import { encodeHex } from "@std/encoding/hex";

export interface XeblogTootProps {
  url: string;
}

const sha256Hex = (data: string): string =>
  encodeHex(crypto.subtle.digestSync("SHA-256", new TextEncoder().encode(data)));

export default function XeblogToot({ url }: XeblogTootProps) {
  const tootHash = sha256Hex(url + ".json");
  const tootJSON = (new TextDecoder("utf-8")).decode(
    Deno.readFileSync(`./src/_data/toots/${tootHash}.json`),
  );
  const toot = JSON.parse(tootJSON);

  const userHash = sha256Hex(toot.attributedTo + ".json");
  const userJSON = (new TextDecoder("utf-8")).decode(
    Deno.readFileSync(`./src/_data/users/${userHash}.json`),
  );
  const user = JSON.parse(userJSON);

  return (
    <>
      <div class="bg-bg-soft dark:bg-bgDark-soft rounded-xl m-2 shadow-md max-w-lg">
        <div class="items-center flex flex-row text-xl px-4 font-bold max-h-[4rem]">
          <img class="rounded-full w-8 h-8" src={user.icon.url} />
          <span class="pl-2">{user.name}</span>
        </div>
        <div class="flex flex-row items-center px-4 py-2">
          <div class="flex flex-wrap px-5">
            <div class="px-2 py-1 m-1 bg-bg-2 dark:bg-bgDark-2 rounded-lg">
              {toot.published}
              <div
                dangerouslySetInnerHTML={{ __html: toot.content }}
              >
              </div>
            </div>
          </div>
        </div>
        <div class={`grid grid-cols-${toot.attachment.length > 1 ? 2 : 1} px-4`}>
          {toot.attachment.map((attachment) => {
            if (attachment.mediaType.startsWith("image/")) {
              return (
                <div class="flex flex-row items-center justify-center m-1 max-w-xs">
                  <a href={attachment.url} target="_blank">
                    <img src={attachment.url} />
                  </a>
                </div>
              );
            } else {
              return (
                <div class="flex flex-row items-center justify-center m-1 max-w-xs">
                  <a href={attachment.url}>
                    {attachment.name}
                  </a>
                </div>
              );
            }
          })}
        </div>
        <a href={toot.url} className="pb-4 px-4" target="_blank">Link</a>
      </div>
    </>
  );
}
