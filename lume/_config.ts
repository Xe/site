import lume from "lume/mod.ts";
import jsx_preact from "lume/plugins/jsx_preact.ts";
import attributes from "lume/plugins/attributes.ts";
import nunjucks from "lume/plugins/nunjucks.ts";
import date from "lume/plugins/date.ts";
import esbuild from "lume/plugins/esbuild.ts";
import mdx from "lume/plugins/mdx.ts";
import tailwindcss from "lume/plugins/tailwindcss.ts";
import postcss from "lume/plugins/postcss.ts";
import sitemap from "lume/plugins/sitemap.ts";
import readInfo from "lume/plugins/reading_info.ts";

import annotateYear from "./plugins/annotate_year.ts";
import feed from "./plugins/feed.ts";
import podcast_feed from "./plugins/podcast_feed.ts";

//import pagefind from "lume/plugins/pagefind.ts";
//import _ from "npm:@pagefind/linux-x64";

import tailwindOptions from "./tailwind.config.js";

import BlockQuote from "./src/_components/BlockQuote.jsx";
import ChatFrame from "./src/_components/ChatFrame.jsx";
import ChatBubble from "./src/_components/ChatBubble.jsx";
import Figure from "./src/_components/Figure.tsx";
import LoadingSpinner from "./src/_components/LoadingSpinner.jsx";
import IntercomButton from "./src/_components/IntercomButton.jsx";
import PullQuote from "./src/_components/PullQuote.jsx";
import TecharoDisclaimer from "./src/_components/TecharoDisclaimer.jsx";
import Reflection from "./src/_components/Reflection.jsx";
import XeblogConv from "./src/_components/XeblogConv.tsx";
import XeblogConvParent from "./src/_components/XeblogConvParent.jsx";
import XeblogHero from "./src/_components/XeblogHero.tsx";
import XeblogPicture from "./src/_components/XeblogPicture.tsx";
import XeblogSlide from "./src/_components/XeblogSlide.tsx";
import XeblogSticker from "./src/_components/XeblogSticker.tsx";
import XeblogToot from "./src/_components/XeblogToot.tsx";
import XeblogVideo from "./src/_components/XeblogVideo.tsx";

import rehypePrism from "npm:rehype-prism-plus/all";

const site = lume({
  src: "./src",
  emptyDest: false,
});

site.copy("static");
site.copy("favicon.ico");
site.copy("static/font/inter/inter.css");
site.copy("static/img");
site.copy("src/static", "static");

site.data("getYear", () => {
  return new Date().getFullYear();
});


site.use(nunjucks());
site.use(jsx_preact());
site.use(attributes());
site.use(date({
  formats: {
    "DATE_US": "MM/dd/yyyy",
  },
}));
site.use(esbuild({
  extensions: [".ts", ".js", ".tsx"],
  options: {
    bundle: true,
    format: "esm",
    minify: true,
    keepNames: true,
    platform: "browser",
    target: "esnext",
    treeShaking: true,
  },
}));
site.use(feed({
  output: ["/blog.rss", "/blog.json"],
  query: "index=true",
  info: {
    title: "Xe Iaso's blog",
    description: "Thoughts and musings from Xe Iaso",
    published: new Date(),
    lang: "en",
  },
  items: {
    title: "=title",
    description: "=desc",
    /*image: (data) => {
      if (data.hero && data.hero.file) {
        return `https://cdn.xeiaso.net/file/christine-static/hero/${data.hero.file}.jpg`;
      }

      if (data.image) {
        return `https://cdn.xeiaso.net/file/christine-static/${data.image}.jpg`;
      }

      return undefined;
    },*/
  },
}));

site.use(podcast_feed({
  output: ["/xecast.rss"],
  query: "is_xecast=true",
  info: {
    title: "Xecast",
    description: "Thoughts and musings from Xe Iaso, now in podcast form",
    published: new Date(),
    lang: "en",
  },
  items: {
    title: "=title",
    description: "=desc",
    podcast: "=podcast",
  },
}))

site.use(mdx({
  components: {
    "BlockQuote": BlockQuote,
    "Blockquote": BlockQuote,
    "ChatFrame": ChatFrame,
    "ChatBubble": ChatBubble,
    "Figure": Figure,
    "Image": Figure,
    "IntercomButton": IntercomButton,
    "LoadingSpinner": LoadingSpinner,
    "TecharoDisclaimer": TecharoDisclaimer,
    "Conv": XeblogConv,
    "XeblogConv": XeblogConv,
    "XesiteConv": XeblogConv,
    "ConvP": XeblogConvParent,
    "Reflection": Reflection,
    "Hero": XeblogHero,
    "XeblogHero": XeblogHero,
    "Picture": XeblogPicture,
    "XeblogPicture": XeblogPicture,
    "PullQuote": PullQuote,
    "Slide": XeblogSlide,
    "XeblogSlide": XeblogSlide,
    "Sticker": XeblogSticker,
    "XeblogSticker": XeblogSticker,
    "Toot": XeblogToot,
    "XeblogToot": XeblogToot,
    "Video": XeblogVideo,
    "XeblogVideo": XeblogVideo,
  },
  rehypePlugins: [
    rehypePrism,
  ],
}));
site.use(tailwindcss({
  extensions: [".mdx", ".jsx", ".tsx", ".md", ".html", ".njx"],
  options: tailwindOptions,
}));
site.use(postcss());
site.use(sitemap({
  query: "",
}));
site.use(readInfo({
  extensions: [".md", ".mdx"],
}));
site.use(annotateYear());

export default site;
