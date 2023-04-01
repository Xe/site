// @jsxImportSource xeact
// @jsxRuntime automatic

import { c } from "xeact";

const onclick = () => {
  Array.from(c("xeblog-slides-fluff")).forEach((el) =>
    el.classList.toggle("hidden")
  );
};

export default function NoFunAllowed() {
  const button = (
    <button
      class=""
      onclick={() => onclick()}
    >
      No fun allowed
    </button>
  );
  return button;
}
