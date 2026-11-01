// Fig 04 — a CD cue sheet and an objgit cue sheet solve the same problem
// Standalone React component. Inline styles only, React is the only dependency.
// Static figure: same panel and palette as the animated ones in this post.

import { AsciiFigure } from "./AsciiFigure.jsx";

const LINES = [
  ["CD AUDIO DISC                                   OBJGIT PACKFILE"],
  [
    "──────────────────────────────────────────      ──────────────────────────────────────────",
  ],
  [
    "U0008_0000002_0000147.WAV    ",
    ["gray", "one big file"],
    "       objects.bin                  ",
    ["gray", "one big file"],
  ],
  [
    "┌──────────┬─────────────┬─────────────────┐    ┌──────────┬─────────────┬─────────────────┐",
  ],
  [
    "│ track 01 │  track 02   │       ...       │    │  blob A  │   tree B    │       ...       │",
  ],
  [
    "└──────────┴─────────────┴─────────────────┘    └──────────┴─────────────┴─────────────────┘",
  ],
  ["  ▲          ▲                                    ▲          ▲"],
  [
    "  │",
    ["green", " 00:00"],
    "    │",
    ["green", " 00:13"],
    "                              │",
    ["green", " @0"],
    "       │",
    ["green", " @142"],
  ],
  [],
  [
    "FILE.cue     ",
    ["gray", "plain text, parsed top-down"],
    "        objects.cue  ",
    ["gray", "packed records, fixed width"],
  ],
  [
    "┌──────────────────────────────────────────┐    ┌──────────────────────────────────────────┐",
  ],
  [
    '│ FILE "U0008_0000002_0000147.WAV" WAVE    │    │ HEADER                              16 B │',
  ],
  [
    "│                                          │    │   magic  version  rec_size  record_count │",
  ],
  [
    "│ TRACK 01 AUDIO                           │    ├──────────────────────────────────────────┤",
  ],
  [
    "│ INDEX 01 ",
    ["green", "00:00"],
    "                           │    │ RECORD 0                            58 B │",
  ],
  [
    "│ TITLE ",
    ["red", '"U0008_0000002_0000147_0001"'],
    "       │    │   ",
    ["red", "hash"],
    " 3a7f…c21        type blob         │",
  ],
  [
    "│                                          │    │   comp zstd            size 311          │",
  ],
  [
    "│ TRACK 02 AUDIO                           │    │   ",
    ["green", "bin_offset"],
    " 0         ",
    ["green", "bin_length"],
    " 142    │",
  ],
  [
    "│ INDEX 01 ",
    ["green", "00:13"],
    "                           │    │   ",
    ["red", "delta_base"],
    " 0000…0000                   │",
  ],
  [
    "│ TITLE ",
    ["red", '"U0008_0000002_0000147_0002"'],
    "       │    ├──────────────────────────────────────────┤",
  ],
  [
    "└──────────────────────────────────────────┘    │ RECORD 1                            58 B │",
  ],
  [
    "                                                └──────────────────────────────────────────┘",
  ],
  [["gray", "to reach track 2 a player parses every"]],
  [
    ["gray", "line above it. it already has the disc."],
    "         ",
    ["gray", "record N sits at 16 + N*58. one seek,"],
  ],
  [
    "                                                ",
    ["gray", "then one ranged GET into objects.bin."],
  ],
];

export default function CueSheetComparison({
  label = "FIG 04",
  title = "A CD cue sheet and an objgit cue sheet solve the same problem",
  fontSize = "12px",
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
