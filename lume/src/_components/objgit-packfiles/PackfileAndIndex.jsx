// Fig 02 — eleven million objects, one packfile, one index
// Standalone React component. Inline styles only, React is the only dependency.
// Same panel and palette as the other figures in this post.
//
// Generated art. The figure cycles one .idx row at a time: the row lights up
// and so does the run of bytes its offset points at, while everything else in
// the index and the packfile stays gray. Three rows, one second each, forever.
//
// It starts at rest with all three pairings lit at once, which is the figure as
// a still. Play walks them one at a time; pause puts all three back up.

import { AnimatedAsciiFigure } from "./AnimatedAsciiFigure.jsx";

const LINES = [
  ["  ", "$ git count-objects -v", "          ", ".git/objects/pack/"],
  [
    "  ",
    "count:            0",
    "             ",
    "└── pack-45986f41..c5786.pack   3.7 GiB",
  ],
  ["  ", "in-pack:   11827138"],
  [
    "  ",
    "packs:            1",
    "             ",
    ["gray", "eleven million loose files would be"],
  ],
  [
    "  ",
    "size-pack:  3876775",
    "             ",
    ["gray", "eleven million inodes. so: one file."],
  ],
  [],
  ["  ", " .idx", "                            ", " .pack"],
  [
    "  ",
    "┌───────────────────────┐",
    "       ",
    "┌──────────────────────────────────────┐",
  ],
  [
    "  ",
    "│ ",
    ["#rowViolet", "0002ff4c..  0x0000c", "violet"],
    "   │",
    "       ",
    "│",
    ["#runCyan", "███", "cyan"],
    " ",
    ["#rest1", "██"],
    " ",
    ["#rest2", "████"],
    " ",
    ["#rest3", "█"],
    " ",
    ["#runGreen", "███████", "green"],
    " ",
    ["#rest4", "██"],
    " ",
    ["#rest5", "█"],
    " ",
    ["#runViolet", "████", "violet"],
    " ",
    ["#rest6", "██████"],
    "│",
  ],
  [
    "  ",
    "│ ",
    ["#rowCyan", "0031ab90..  0x0a13f", "cyan"],
    "   │",
    "       ",
    "└──────────────────────────────────────┘",
  ],
  ["  ", "│ ", ["#rowGreen", "1c7a26a9..  0x1f3a4", "green"], "   │"],
  ["  ", "│ ", ["#restRows", "..."], "                   │"],
  ["  ", "└───────────────────────┘"],
  [],
  [
    "  ",
    [
      "gray",
      "each row's colour is the run of bytes its offset points at, so a",
    ],
  ],
  ["  ", ["gray", "read is a seek to an offset inside one very big file"]],
];

// Everything is on screen from the first frame; only the focus moves.
const ALL = [
  "rowViolet",
  "rowCyan",
  "rowGreen",
  "restRows",
  "runViolet",
  "runCyan",
  "runGreen",
  "rest1",
  "rest2",
  "rest3",
  "rest4",
  "rest5",
  "rest6",
];

const STEPS = [
  {
    show: ALL,
    focus: ["rowViolet", "runViolet"],
    caption: "0002ff4c.. → seek to 0x0000c",
  },
  {
    focus: ["rowCyan", "runCyan"],
    caption: "0031ab90.. → seek to 0x0a13f",
  },
  {
    focus: ["rowGreen", "runGreen"],
    caption: "1c7a26a9.. → seek to 0x1f3a4",
  },
];

// Where the figure sits when it is not cycling, and where it starts: all three
// pairings lit at once. The filler rows and runs stay gray, because they are
// the haystack, not the needles.
const REST = {
  focus: [
    "rowViolet",
    "runViolet",
    "rowCyan",
    "runCyan",
    "rowGreen",
    "runGreen",
  ],
  caption: "three rows, three offsets, three runs of bytes",
};

export default function PackfileAndIndex({
  label = "FIG 02",
  title = "Eleven million objects, one packfile, one index",
  fontSize,
}) {
  return (
    <AnimatedAsciiFigure
      label={label}
      title={title}
      lines={LINES}
      steps={STEPS}
      loop
      rest={REST}
      stepMs={1000}
      fontSize={fontSize}
    />
  );
}
