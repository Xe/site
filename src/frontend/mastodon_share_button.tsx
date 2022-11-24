import { g, r, u, x } from "xeact";

r(() => {
  const root = g("mastodon_share_button");

  let defaultURL = localStorage["mastodon_instance"];

  const title = document.querySelectorAll('meta[property="og:title"]')[0]
    .getAttribute("content");
  let series = g("mastodon_share_series").innerText;
  if (series != "") {
    series = `#${series} `;
  }
  const tags = g("mastodon_share_tags");
  const articleURL = u();

  const tootTemplate = `${title}

${articleURL}

${series}${tags.innerText}@cadey@pony.social`;

  const instanceBox = (
    <input type="text" placeholder="https://pony.social" value={defaultURL} />
  );
  const tootBox = <textarea rows="6" cols="40">{tootTemplate}</textarea>;

  const doShare = () => {
    const instanceURL = instanceBox.value;
    localStorage["mastodon_instance"] = instanceURL;
    const text = tootBox.value;
    const mastodon_url = u(instanceURL + "/share", { text });
    console.log({ text, mastodon_url });
    window.open(mastodon_url, "_blank");
  };

  const shareButton = <button onclick={doShare}>Share</button>;

  x(root);

  root.appendChild(
    <div>
      <details>
        <summary>Share on Mastodon</summary>
        <span>Instance URL (https://mastodon.example)</span>
        <br />
        {instanceBox}
        <br />
        {tootBox}
        <br />
        {shareButton}
      </details>
    </div>,
  );
});
