// Fig 07 — S3 requests to push each repo, git packfiles vs objgit .bin/.cue
// Standalone React component. Inline styles only, React is the only dependency.
// Static figure: same panel and palette as the other figures in this post.
//
// Generated art. Do not hand-edit the columns — regenerate from the script in
// the ascii-diagrams skill. The axis is log10: every tick is 10x the last, and
// every bar starts at column 14 so the tick marks read straight down the bars.

import { AsciiFigure } from "./AsciiFigure.jsx";

const LINES = [
  [
    "              ",
    ["gray", "1"],
    "         ",
    ["gray", "10"],
    "        ",
    ["gray", "100"],
    "       ",
    ["gray", "1,000"],
    "     ",
    ["gray", "10,000"],
  ],
  ["              ", ["gray", "├─────────┼─────────┼─────────┼─────────┤"]],
  [],
  ["  ", "objgit"],
  [
    "    ",
    ["gray", "old"],
    "       ",
    ["red", "████████████████████████"],
    "                     ",
    ["red", "231"],
  ],
  [
    "    ",
    ["gray", "new"],
    "       ",
    ["green", "█████████████"],
    "                                 ",
    ["green", "18"],
  ],
  [],
  ["  ", "Xe/x"],
  [
    "    ",
    ["gray", "old"],
    "       ",
    ["red", "████████████████████████████████████████"],
    "   ",
    ["red", "9,236"],
  ],
  [
    "    ",
    ["gray", "new"],
    "       ",
    ["green", "███████████████"],
    "                               ",
    ["green", "30"],
  ],
  [],
  ["  ", "tigris-blog"],
  [
    "    ",
    ["gray", "old"],
    "       ",
    ["red", "███████████████████████████████████"],
    "        ",
    ["red", "3,324"],
  ],
  [
    "    ",
    ["gray", "new"],
    "       ",
    ["green", "█████████████████████"],
    "                        ",
    ["green", "136"],
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

export default function PushRequestCounts({
  label = "FIG 07",
  title = "S3 requests to push, before and after the format change (log scale)",
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
