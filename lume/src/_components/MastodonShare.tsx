import { useState } from "npm:preact/hooks";

const u = (url = "", params = {}) => {
    let result = new URL(url, window.location.href);
    Object.entries(params).forEach((kv) => {
      let [k, v] = kv;
      result.searchParams.set(k, v as string);
    });
    return result.toString();
  };

export interface MastodonShareButtonProps {
  title: string;
  url: string;
  series?: string;
  tags: string;
}

export default function MastodonShareButton(
  { title, url, series, tags }: MastodonShareButtonProps,
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

  const [getURL, setURL] = useState(defaultURL);
  const [getToot, setToot] = useState(tootTemplate);

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
          oninput={(e) => setURL(e.target.value)}
        />
        <br />
        <textarea
          rows={6}
          cols={40}
          oninput={(e) => setToot(e.target.value)}
        >
          {getToot()}
        </textarea>
        <br />
        <button
          onclick={() => {
            let instanceURL = getURL;

            if (!instanceURL.startsWith("https://")) {
              instanceURL = `https://${instanceURL}`;
            }

            localStorage["mastodon_instance"] = instanceURL;
            const text = getToot;
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
