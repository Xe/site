// FIG 09 — S3 requests to clone each repo, git packfiles vs objgit .bin/.cue
// Standalone React component. Inline styles only, React is the only dependency.
// Static figure: same panel and palette as the other figures in this post.
//
// Generated art. Do not hand-edit the columns — regenerate with
// `node gen-requests.mjs clone` from the ascii-diagrams skill. Same layout and
// same log10 axis as FIG 07 so the push and clone charts can be read against
// each other: every tick is 10x the last, every bar starts at column 14.

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
    ["red", "█████████████████████████"],
    "                    ",
    ["red", "323"],
  ],
  [
    "    ",
    ["gray", "new"],
    "       ",
    ["green", "████████████"],
    "                                  ",
    ["green", "17"],
  ],
  [],
  ["  ", "Xe/x"],
  [
    "    ",
    ["gray", "old"],
    "       ",
    ["red", "██████████████████████████████████████"],
    "     ",
    ["red", "6,428"],
  ],
  [
    "    ",
    ["gray", "new"],
    "       ",
    ["green", "████████████"],
    "                                  ",
    ["green", "17"],
  ],
  [],
  ["  ", "tigris-blog"],
  [
    "    ",
    ["gray", "old"],
    "       ",
    ["red", "████████████████████████████████████"],
    "       ",
    ["red", "3,675"],
  ],
  [
    "    ",
    ["gray", "new"],
    "       ",
    ["green", "██████████████████████"],
    "                       ",
    ["green", "158"],
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

export default function CloneRequestCounts({
  label = "FIG 09",
  title = "S3 requests to clone, before and after the format change (log scale)",
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
