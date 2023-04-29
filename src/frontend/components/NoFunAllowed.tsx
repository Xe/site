import { c } from "@xeserv/xeact";

const onclick = () => {
  Array.from(c("xeblog-slides-fluff")).forEach((el) =>
    el.classList.toggle("hidden")
  );
};

export default function NoFunAllowed() {
  const button = (
    <button
      class=""
      onClick={() => onclick()}
    >
      No fun allowed
    </button>
  );
  return button;
}
