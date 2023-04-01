import * as esbuild from "@esbuild";
import { denoPlugin } from "@esbuild/deno";

const result = await esbuild.build({
  plugins: [denoPlugin({
    importMapURL: new URL("./import_map.json", import.meta.url),
  })],
  entryPoints: [
    "./components/ConvSnippet.tsx",
    "./components/MastodonShareButton.tsx",
    "./components/Video.tsx",
    "./components/WASITerm.tsx",
  ],
  outdir: Deno.env.get("WRITE_TO")
    ? Deno.env.get("WRITE_TO")
    : "../../static/xeact",
  bundle: true,
  splitting: true,
  format: "esm",
  minifyWhitespace: !!Deno.env.get("MINIFY"),
  inject: ["xeact"],
  jsxFactory: "h",
});
console.log(result.outputFiles);

esbuild.stop();
