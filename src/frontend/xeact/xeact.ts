export * from "./xeact.js";

declare global {
    export namespace JSX {
        interface IntrinsicElements {
            [elemName: string]: any;
        }
    }
}
