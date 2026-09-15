// Fig 03 — one ranged GET pulls 366 bytes out of the middle of a 128 MiB file
// Standalone React component. Inline styles only, React is the only dependency.
// Static figure: same panel and palette as the other figures in this post.
//
// Generated art. The ██ in the container, the ▲ under it, the │ under that and
// the "└── 366 bytes" pointer all start at column 38, so the callout reads as
// one vertical down from the object's first byte.

import { AsciiFigure } from "./AsciiFigure.jsx";

const LINES = [
  ["  ", "packs/019a7f3c..bin"],
  ["  ", "┌────────────────────────────────────────────────┐"],
  [
    "  ",
    "│",
    "                                   ",
    ["amber", "██"],
    "           ",
    "│",
  ],
  [
    "  ",
    "└───────────────────────────────────",
    ["amber", "▲"],
    "────────────┘",
  ],
  [
    "  ",
    ["gray", " 0"],
    "                                  ",
    ["amber", "│"],
    ["gray", "      128 MiB"],
  ],
  [
    "                                      ",
    ["amber", "└── 366 bytes, right here"],
  ],
  [],
  ["  ", ["gray", "you know where it starts and how long it is, so ask for"]],
  ["  ", ["gray", "exactly that span and nothing else:"]],
  [],
  ["  ", "GET /packs/019a7f3c..bin"],
  ["  ", ["green", "Range: bytes=100663296-100663661"]],
  [],
  ["  ", ["gray", "and that is all the bucket sends back:"]],
  [],
  ["  ", ["green", "206 Partial Content"]],
  ["  ", "Content-Range: bytes 100663296-100663661/134217728"],
  ["  ", "┌────┐"],
  [
    "  ",
    "│",
    ["amber", "████"],
    "│",
    "  ",
    ["gray", "366 B on the wire, not 128 MiB"],
  ],
  ["  ", "└────┘"],
];

export default function RangeRequestGap({
  label = "FIG 03",
  title = "One ranged GET, 366 bytes out of the middle of 128 MiB",
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
