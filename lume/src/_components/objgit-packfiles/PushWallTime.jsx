// FIG 08 — wall time to push each repo, git packfiles vs objgit .bin/.cue
// Standalone React component. Inline styles only, React is the only dependency.
// Static figure: same panel and palette as the other figures in this post.
//
// Generated art. Do not hand-edit the columns — regenerate with
// `node gen-time.mjs push` from the ascii-diagrams skill. The axis is log10:
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
    ["red", "███████████"],
    "                              ",
    ["red", "8.7s"],
  ],
  [
    "    ",
    ["gray", "new"],
    "       ",
    ["green", "████"],
    "                                     ",
    ["green", "2.2s"],
    "   ",
    ["green", "4.0x"],
  ],
  [],
  ["  ", "Xe/x"],
  [
    "    ",
    ["gray", "old"],
    "       ",
    ["red", "████████████████████████████"],
    "          ",
    ["red", "3m29.4s"],
  ],
  [
    "    ",
    ["gray", "new"],
    "       ",
    ["green", "██████████████"],
    "                          ",
    ["green", "14.3s"],
    "  ",
    ["green", "14.6x"],
  ],
  [],
  ["  ", "tigris-blog"],
  [
    "    ",
    ["gray", "old"],
    "       ",
    ["red", "██████████████████████████"],
    "            ",
    ["red", "2m13.4s"],
  ],
  [
    "    ",
    ["gray", "new"],
    "       ",
    ["green", "█████████████████"],
    "                       ",
    ["green", "26.5s"],
    "   ",
    ["green", "5.0x"],
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

export default function PushWallTime({
  label = "FIG 08",
  title = "Wall time to push, log scale, speedup on the right",
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
