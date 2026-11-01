// Fig 01 — a sea of objects, and a few names into it
// Standalone React component. Inline styles only, React is the only dependency.
// Static figure: the loose objects on the left are colour-matched to the graph
// nodes they spell out on the right.

import { AsciiFigure } from "./AsciiFigure.jsx";

const LINES = [
  ["  .git/objects/                    ", ["amber", "refs/heads/main"]],
  ["                                     ", ["amber", "│"]],
  [
    "  ",
    ["violet", "├── 1c/7a26a901..ec7966"],
    "  ",
    ["gray", "─────▶"],
    "  ",
    ["violet", "commit 1c7a26a"],
  ],
  ["  ", ["violet", "│"], "                                  │"],
  [
    "  ",
    ["cyan", "├── 8e/67afbb2e..857bd3"],
    "  ",
    ["gray", "─────▶"],
    "    ",
    ["cyan", "tree 8e67afb"],
  ],
  ["  ", ["cyan", "│"], "                                  │  hello.txt"],
  [
    "  ",
    ["green", "└── 9c/c9867337..09fe26"],
    "  ",
    ["gray", "─────▶"],
    "    ",
    ["green", "blob 9cc9867"],
  ],
  ["                                          ", ["green", '"Hello, blog!"']],
  [],
  [
    "  ",
    ["gray", "the filename is the sha1 of the bytes in the file, so the same"],
  ],
  ["  ", ["gray", "content is always, everywhere, the very same object"]],
];

export default function SeaOfObjects({
  label = "FIG 01",
  title = "A sea of objects, and a few names into it",
  fontSize,
}) {
  return (
    <AsciiFigure
      label={label}
      title={title}
      lines={LINES}
      fontSize={fontSize}
    />
  );
}
