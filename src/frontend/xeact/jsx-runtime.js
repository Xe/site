import { h } from './xeact.ts';

/**
 * Create a DOM element, assign the properties of `data` to it, and append all `data.children`.
 *
 * @type{function(string, Object=): HTMLElement}
 */
export const jsx = (tag, data) => {
  let children = data.children;
  delete data.children;
  const result = h(tag, data, children);
  result.classList.value = result.class;
  return result;
};
export const jsxs = jsx;
export const jsxDEV = jsx;
