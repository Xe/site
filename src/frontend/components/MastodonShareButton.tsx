import { u } from "@xeserv/xeact";
import { useState } from 'preact/hooks';

export interface MastodonShareButtonProps {
  title: string;
  url: string;
  series?: string;
  tags: string[];
}

export default function MastodonShareButton(
  { title, url = u("", {}), series, tags }: MastodonShareButtonProps,
) {
  let defaultURL = localStorage["mastodon_instance"];

  if (defaultURL == undefined) {
    defaultURL = "";
  }

  const tootTemplate = `${title}

${url}

${series ? "#" + series + " " : ""}${
    tags ? tags.map((x) => "#" + x).join(" ") : ""
  } @cadey@pony.social`;

  const [theURL, setURL] = useState(defaultURL);
  const [toot, setToot] = useState(tootTemplate);

  return (
    <div>
      <details>
        <summary>Share on Mastodon</summary>
        <span>{"Instance URL (https://mastodon.example)"}</span>
        <br />
        <input
          type="text"
          placeholder="https://pony.social"
          value={defaultURL}
            onInput={(e) => {
                const target = e.target as HTMLInputElement;
                setURL(target.value);
            }}
        />
        <br />
        <textarea
          rows={6}
          cols={40}
            onInput={(e) => {
                const target = e.target as HTMLTextAreaElement;
                setToot(target.value)
            }}
        >
          {toot}
        </textarea>
        <br />
        <button
          onClick={() => {
            let instanceURL = theURL;

            if (!instanceURL.startsWith("https://")) {
              instanceURL = `https://${instanceURL}`;
            }

            localStorage["mastodon_instance"] = instanceURL;
            const text = toot;
            const mastodonURL = u(instanceURL + "/share", {
              text,
              visibility: "public",
            });
            console.log({ text, mastodonURL });
            window.open(mastodonURL, "_blank");
          }}
        >
          Share
        </button>
      </details>
    </div>
  );
}
