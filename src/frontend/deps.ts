import * as wasi from "https://deno.land/x/wasm@v1.2.2/wasi.ts";
import * as xeact from "xeact";

await wasi.init();

export { wasi, xeact };
