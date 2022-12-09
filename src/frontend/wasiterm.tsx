import { wasi as wasiMod, xeact } from "./deps.ts";

const { x, t } = xeact;
//const { Terminal } = xterm;
const { WASI } = wasiMod;

const init = async (rootElem: Element, wasmURL: string) => {
    const termElem = <code></code>;

    //const term = new Terminal();
    //term.open(termElem);

    const runProgram = async () => {
        await wasiMod.init(new URL("https://cdn.xeiaso.net/file/christine-static/wasm/5410143de81b20061e9750d1cf80aceef56d2938ab949e30dd7b13fa699307ad.wasm"));

        termElem.appendChild(t(`loading ${wasmURL}`));

        const wasi = new WASI({
            env: {
                "HOSTNAME": "pyra",
            },
        })

        const moduleBytes = fetch(wasmURL);
        const module = await WebAssembly.compileStreaming(moduleBytes);
        await wasi.instantiate(module, {});
        //term.writeln("executing");
        termElem.appendChild(t("executing"));
        let exitCode = wasi.start();
        let stdout = wasi.getStdoutString();
        console.log(`${stdout}\n\n(exit code: ${exitCode})`);
        termElem.appendChild(t(`${stdout}\n\n(exit code: ${exitCode})`));
    };

    const runButton = <button onclick={runProgram}>Run</button>;

    const root = <div>
        <link rel="stylesheet" href="https://cdn.xeiaso.net/file/christine-static/wasm/xterm/xterm-a36f07105014cc9220cae423d97c30d1a59fdb0230da8e53bb74bb0faade4310.css" type="text/css" />
        {runButton}
        <pre>{termElem}</pre>
    </div>;

    x(rootElem);
    rootElem.appendChild(root);

    //term.writeln(`loading ${wasmURL}`);
    //term.writeln(`${stdout}\n\n(exit code: ${exitCode})`);
}

export { init };
