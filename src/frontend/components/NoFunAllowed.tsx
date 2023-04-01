// @jsxImportSource xeact
// @jsxRuntime automatic

import { c } from "xeact";

export default function NoFunAllowed() {
  return (
    <button
      onclick={(() => {
        Array.from(c("xeblog-slides-fluff")).forEach((el) =>
          el.classList.toggle("hidden")
        );
      })()}
    >
      No fun allowed
    </button>
  );
}
