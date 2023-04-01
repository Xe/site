// @jsxImportSource xeact
// @jsxRuntime automatic

import { t, x } from "xeact";
import Terminal from "@xterm";
import * as fitAdd from "@xterm/addon-fit";
import { Fd, File, PreopenDirectory, WASI } from "@bjorn3/browser_wasi_shim";

class XtermStdio extends Fd {
  term: Terminal;

  constructor(term: Terminal) {
    super();
    this.term = term;
  }

  fd_write(view8: Uint8Array, iovs: any) {
    let nwritten = 0;
    for (let iovec of iovs) {
      console.log(
        iovec.buf_len,
        iovec.buf_len,
        view8.slice(iovec.buf, iovec.buf + iovec.buf_len),
      );
      let buffer = view8.slice(iovec.buf, iovec.buf + iovec.buf_len);
      this.term.write(buffer);
      nwritten += iovec.buf_len;
    }
    return { ret: 0, nwritten };
  }
}

const loadExternalFile = async (path: string) => {
  return new File(await (await (await fetch(path)).blob()).arrayBuffer());
};

export interface WASITermProps {
  href: string;
  env: string[];
  args: string[];
}

export default function WASITerm({ href, env, args }: WASITermProps) {
  const root = <div style="max-width:80ch;max-height:20ch"></div>;

  const term = new Terminal({
    convertEol: true,
    fontFamily: "Iosevka Curly Iaso",
  });

  const fit = new fitAdd.default();
  term.loadAddon(fit);
  fit.fit();

  return (
    <div>
      {root}
      <button
        onClick={async () => {
          term.writeln(`\x1B[93mfetching ${href}...\x1B[0m`);
          const wasm = await WebAssembly.compileStreaming(fetch(href));
          term.writeln("\x1B[93mdone, instantiating\x1B[0m");

          const fds = [
            new XtermStdio(term),
            new XtermStdio(term),
            new XtermStdio(term),
            new PreopenDirectory("/tmp", {}),
          ];

          let wasi = new WASI(args, env, fds);
          term.clear();
          // let inst = await WebAssembly.instantiate(wasm, {
          //   "wasi_snapshot_preview1": wasi.wasiImport,
          // });
          // wasi.start(inst);
        }}
      >
        Run
      </button>
    </div>
  );
}
