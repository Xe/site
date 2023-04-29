declare module "@xeserv/xeact" {
/**
 * Creates a DOM element, assigns the properties of `data` to it, and appends all `children`.
 *
 * @type{function(string|Function, Object=, Node|Array.<Node|string>=)}
 */
export const h: (arg0: string | Function, arg1: any | undefined, arg2: (Node | Array<Node | string>) | undefined) => any;
/**
 * Create a text node.
 *
 * Equivalent to `document.createTextNode(text)`
 *
 * @type{function(string): Text}
 */
export const t: (arg0: string) => Text;
/**
 * Remove all child nodes from a DOM element.
 *
 * @type{function(Node)}
 */
export const x: (arg0: Node) => any;
/**
 * Get all elements with the given ID.
 *
 * Equivalent to `document.getElementById(name)`
 *
 * @type{function(string): HTMLElement}
 */
export const g: (arg0: string) => HTMLElement;
/**
 * Get all elements with the given class name.
 *
 * Equivalent to `document.getElementsByClassName(name)`
 *
 * @type{function(string): HTMLCollectionOf.<Element>}
 */
export const c: (arg0: string) => HTMLCollectionOf<Element>;
/** @type{function(string): HTMLCollectionOf.<Element>} */
export const n: (arg0: string) => HTMLCollectionOf<Element>;
/**
 * Generate a relative URL from `url`, appending all key-value pairs from `params` as URL-encoded parameters.
 *
 * @type{function(string=, Object=): string}
 */
export const u: (arg0: string | undefined, arg1: any | undefined) => string;
/**
 * Get all elements matching the given HTML selector.
 *
 * Matches selectors with `document.querySelectorAll(selector)`
 *
 * @type{function(string): Array.<HTMLElement>}
 */
export const s: (arg0: string) => Array<HTMLElement>;
/**
 * Takes a callback to run when all DOM content is loaded.
 *
 * Equivalent to `window.addEventListener('DOMContentLoaded', callback)`
 *
 * @type{function(function())}
 */
export const r: (arg0: () => any) => any;
/**
 * Allows a stateful value to be tracked by consumers.
 *
 * This is the Xeact version of the React useState hook.
 *
 * @type{function(any): [function(): any, function(any): void]}
 */
export const useState: (arg0: any) => [() => any, (arg0: any) => void];
/**
 * Debounce an action for up to ms milliseconds.
 *
 * @type{function(number): function(function(any): void)}
 */
  export const d: (arg0: number) => (arg0: (arg0: any) => void) => any;
}
