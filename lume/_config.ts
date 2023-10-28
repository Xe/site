import lume from "lume/mod.ts";
import jsx_preact from "lume/plugins/jsx_preact.ts";
import attributes from "lume/plugins/attributes.ts";
import date from "lume/plugins/date.ts";
import esbuild from "lume/plugins/esbuild.ts";
import feed from "lume/plugins/feed.ts";
import mdx from "lume/plugins/mdx.ts";
import pagefind from "lume/plugins/pagefind.ts";
import tailwindcss from "lume/plugins/tailwindcss.ts";
import postcss from "lume/plugins/postcss.ts";
import sitemap from "lume/plugins/sitemap.ts";

import tailwindOptions from "./tailwind.config.js";

import BlockQuote from "./src/_components/BlockQuote.jsx";
import Figure from "./src/_components/Figure.tsx";
import XeblogConv from "./src/_components/XeblogConv.tsx";
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
})

site.use(jsx_preact());
site.use(attributes());
site.use(date());
site.use(esbuild({ esm: true }));
site.use(feed({
  output: ["/blog.rss", "/blog.json"],
  query: "index=true",
  info: {
    title: "Xe Iaso's blog",
    description: "Thoughts and musings from Xe Iaso",
  },
  items: {
    title: "=title",
    description: "=excerpt",
  },
}));
site.use(mdx({
  components: {
    "BlockQuote": BlockQuote,
    "Figure": Figure,
    "Image": Figure,
    "XeblogConv": XeblogConv,
    "XesiteConv": XeblogConv,
    "XeblogHero": XeblogHero,
    "XeblogPicture": XeblogPicture,
    "XeblogSlide": XeblogSlide,
    "XeblogSticker": XeblogSticker,
    "XeblogToot": XeblogToot,
    "XeblogVideo": XeblogVideo,
  },
  rehypePlugins: [
    rehypePrism,
  ],
}));
site.use(pagefind({
  indexing: {
    bundleDirectory: "_pagefind",
    glob: "**/*.html",
    rootObject: "article",
  },
}));
site.use(tailwindcss({
  extensions: [".mdx", ".jsx", ".tsx", ".md", ".html", ".njx"],
  options: tailwindOptions,
}));
site.use(postcss());
site.use(sitemap({
  query: "index=true",
}));

export default site;
