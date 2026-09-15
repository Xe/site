// Fig 05 — objects.cue on the wire: one header, then fixed-width records
// Standalone React component. Inline styles only, React is the only dependency.
// Static figure: same panel and palette as the animated ones in this post.

import { AsciiFigure } from "./AsciiFigure.jsx";

const LINES = [
  ["objects.cue HEADER   ", ["gray", "16 bytes, once at the top of the file"]],
  ["┌──────────────┬─────────┬──────────┬────────────────────┐"],
  ['│ magic "OGCU" │ version │ rec_size │ record_count       │'],
  ["│          4 B │ u16 2 B │ u16  2 B │ u64            8 B │"],
  ["└──────────────┴─────────┴──────────┴────────────────────┘"],
  [["gray", "0              4         6          8                   16"]],
  [],
  ["objects.cue RECORD   ", ["gray", "58 bytes, repeated record_count times"]],
  [
    "┌──────────────┬──────┬──────┬────────────┬────────────┬────────┬──────────────┐",
  ],
  [
    "│ ",
    ["red", "hash"],
    "         │ type │ comp │ ",
    ["green", "bin_offset"],
    " │ ",
    ["green", "bin_length"],
    " │ size   │ ",
    ["red", "delta_base"],
    "   │",
  ],
  [
    "│ ",
    ["red", "sha1  20 B"],
    "   │ u8 1B│ u8 1B│ ",
    ["green", "u64   8 B"],
    "  │ ",
    ["green", "u32   4 B"],
    "  │ u32 4B │ ",
    ["red", "sha1  20 B"],
    "   │",
  ],
  [
    "└──────────────┴──────┴──────┴────────────┴────────────┴────────┴──────────────┘",
  ],
  [
    [
      "gray",
      "0              20     21     22           30           34       38            58",
    ],
  ],
  [],
  [
    ["gray", "identity"],
    " ",
    ["red", "hash"],
    ", ",
    ["red", "delta_base"],
    "   ",
    ["gray", "location"],
    " ",
    ["green", "bin_offset"],
    ", ",
    ["green", "bin_length"],
    "   ",
    ["gray", "decode"],
    " type, comp, size",
  ],
];

export default function ObjgitCueFormat({
  label = "FIG 05",
  title = "objects.cue on the wire: one header, then fixed-width records",
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
