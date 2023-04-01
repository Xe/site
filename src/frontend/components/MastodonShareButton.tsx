// @jsxImportSource xeact
// @jsxRuntime automatic

import { u } from "xeact";

export interface MastodonShareButtonProps {
  title: string;
  url: string;
  series?: string;
  tags: string;
}

export default function MastodonShareButton(
  { title, url = u(), series, tags }: MastodonShareButtonProps,
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

  const instanceBox = (
    <input
      type="text"
      placeholder="https://pony.social"
      value={defaultURL}
    />
  );
  const tootBox = (
    <textarea rows={6} cols={40}>
      {tootTemplate}
    </textarea>
  );

  return (
    <div>
      <details>
        <summary>Share on Mastodon</summary>
        <span>{"Instance URL (https://mastodon.example)"}</span>
        <br />
        {instanceBox}
        <br />
        {tootBox}
        <br />
        <button
          onClick={() => {
            let instanceURL = instanceBox.value;

            if (!instanceURL.startsWith("https://")) {
              instanceURL = `https://${instanceURL}`;
            }

            localStorage["mastodon_instance"] = instanceURL;
            const text = tootBox.value;
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
