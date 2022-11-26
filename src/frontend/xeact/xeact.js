/**
 * Creates a DOM element, assigns the properties of `data` to it, and appends all `children`.
 *  
 * @type{function(string|Function, Object=, Node|Array.<Node|string>=)}
 */
const h = (name, data = {}, children = []) => {
    const result = typeof name == "function" ? name(data) : Object.assign(document.createElement(name), data);
    if (!Array.isArray(children)) {
        children = [children];
    }
    result.append(...children);
    return result;
};

/**
 * Create a text node.
 * 
 * Equivalent to `document.createTextNode(text)`
 * 
 * @type{function(string): Text}
 */
const t = (text) => document.createTextNode(text);

/**
 * Remove all child nodes from a DOM element.
 * 
 * @type{function(Node)} 
 */
const x = (elem) => {
    while (elem.lastChild) {
        elem.removeChild(elem.lastChild);
    }
};

/**
 * Get all elements with the given ID.
 * 
 * Equivalent to `document.getElementById(name)`
 * 
 * @type{function(string): HTMLElement}
 */
const g = (name) => document.getElementById(name);

/**
 * Get all elements with the given class name.
 * 
 * Equivalent to `document.getElementsByClassName(name)`
 * 
 * @type{function(string): HTMLCollectionOf.<Element>} 
 */
const c = (name) => document.getElementsByClassName(name);

/** @type{function(string): HTMLCollectionOf.<Element>} */
const n = (name) => document.getElementsByName(name);

/**
 * Get all elements matching the given HTML selector.
 * 
 * Matches selectors with `document.querySelectorAll(selector)`
 * 
 * @type{function(string): Array.<HTMLElement>}
 */
const s = (selector) => Array.from(document.querySelectorAll(selector));

/**
 * Generate a relative URL from `url`, appending all key-value pairs from `params` as URL-encoded parameters.
 * 
 * @type{function(string=, Object=): string}
 */
const u = (url = "", params = {}) => {
    let result = new URL(url, window.location.href);
    Object.entries(params).forEach((kv) => {
        let [k, v] = kv;
        result.searchParams.set(k, v);
    });
    return result.toString();
};

/**
 * Takes a callback to run when all DOM content is loaded.
 * 
 * Equivalent to `window.addEventListener('DOMContentLoaded', callback)`
 * 
 * @type{function(function())}
 */
const r = (callback) => window.addEventListener('DOMContentLoaded', callback);

export { h, t, x, g, c, n, u, s, r };
