// FIG 10 — wall time to clone each repo, git packfiles vs objgit .bin/.cue
// Standalone React component. Inline styles only, React is the only dependency.
// Static figure: same panel and palette as the other figures in this post.
//
// Generated art. Do not hand-edit the columns — regenerate with
// `node gen-time.mjs clone` from the ascii-diagrams skill. The axis is log10:
// every tick is 10x the last, and every bar starts at column 14 so the tick
// marks read straight down the bars. The right-hand column is old/new.

import { AsciiFigure } from "./AsciiFigure.jsx";

const LINES = [
  [
    "              ",
    ["gray", "1s"],
    "          ",
    ["gray", "10s"],
    "         ",
    ["gray", "100s"],
    "        ",
    ["gray", "1,000s"],
  ],
  ["              ", ["gray", "├───────────┼───────────┼───────────┤"]],
  [],
  ["  ", "objgit"],
  [
    "    ",
    ["gray", "old"],
    "       ",
    ["red", "█████████████"],
    "                           ",
    ["red", "11.8s"],
  ],
  [
    "    ",
    ["gray", "new"],
    "       ",
    ["green", "█████"],
    "                                    ",
    ["green", "2.6s"],
    "   ",
    ["green", "4.5x"],
  ],
  [],
  ["  ", "Xe/x"],
  [
    "    ",
    ["gray", "old"],
    "       ",
    ["red", "████████████████████████████"],
    "          ",
    ["red", "3m23.5s"],
  ],
  [
    "    ",
    ["gray", "new"],
    "       ",
    ["green", "█████████████████████"],
    "                   ",
    ["green", "54.4s"],
    "   ",
    ["green", "3.7x"],
  ],
  [],
  ["  ", "tigris-blog"],
  [
    "    ",
    ["gray", "old"],
    "       ",
    ["red", "██████████████████████████"],
    "            ",
    ["red", "2m23.6s"],
  ],
  [
    "    ",
    ["gray", "new"],
    "       ",
    ["green", "███████████████████████"],
    "               ",
    ["green", "1m22.0s"],
    "   ",
    ["green", "1.8x"],
  ],
  [],
  [
    "  ",
    ["red", "██"],
    " ",
    ["red", "old"],
    "  ",
    ["gray", "git packfiles behind a filesystem shim"],
  ],
  [
    "  ",
    ["green", "██"],
    " ",
    ["green", "new"],
    "  ",
    ["gray", ".bin/.cue columnar packfiles"],
  ],
];

export default function CloneWallTime({
  label = "FIG 10",
  title = "Wall time to clone, log scale, speedup on the right",
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
