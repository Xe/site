import { g, h, x } from "./xeact.min.js";
import { div, span } from "./xeact-html.min.js";

export const mkConversation = (who, mood, message, extraClasses = "") =>
      h("div", {className: "conversation gruvbox-dark " + extraClasses}, [
          h("div", {className: "conversation-picture conversation-smol"}, [
              h("picture", {}, [
                  h("source", {type: "image/avif", srcset: `https://cdn.xeiaso.net/file/christine-static/stickers/${who.toLowerCase()}/${mood}.avif`}),
                  h("source", {type: "image/webp", srcset: `https://cdn.xeiaso.net/file/christine-static/stickers/${who.toLowerCase()}/${mood}.webp`}),
                  h("img", {alt: `${who} is ${mood}`, src: `https://cdn.xeiaso.net/file/christine-static/stickers/${who.toLowerCase()}/${mood}.png`})
              ])
          ]),
          h("div", {className: "conversation-chat"}, [
              h("span", {innerText: "<"}),
              h("b", {innerText: who}),
              h("span", {innerText: "> "}),
              span({}, Array.from(message))
          ])
      ]);

export class Conversation extends HTMLElement {
    constructor() {
        super();

        let root = this.attachShadow({mode: "open"});
        let who = this.getAttribute("name");
        let mood = this.getAttribute("mood");

        root.appendChild(h("link", {rel: "stylesheet", href: "/css/hack.css"}));
        root.appendChild(h("link", {rel: "stylesheet", href: "/css/gruvbox-dark.css"}));
        root.appendChild(h("link", {rel: "stylesheet", href: "/css/shim.css"}));
        root.appendChild(mkConversation(who, mood, this.childNodes));
    }
}

window.customElements.define("xeblog-conv", Conversation);
