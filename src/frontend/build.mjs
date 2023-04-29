import * as esbuild from "esbuild";

const result = await esbuild.build({
  entryPoints: process.argv,
  outdir: process.env.WRITE_TO
    ? process.env.WRITE_TO
    : "../../static/xeact",
  bundle: true,
  splitting: true,
  format: "esm",
  minifyWhitespace: !!process.env.MINIFY,
});
