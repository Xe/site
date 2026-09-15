// Fig 06 — racing the beam: four ranged GETs while the container is 2 MiB in
// Standalone React component. Inline styles only, React is the only dependency.
// Same panel and palette as the other figures in this post.
//
// Generated art. The bar is 48 cells over a 128 MiB container, so one cell is
// 2.67 MiB. The four wanted objects sit at cells 6, 15, 29 and 40, which is
// where their "16 MiB in", "40 MiB in", "77 MiB in" and "107 MiB in" come
// from: offset and column are derived from the same cell number, so they
// cannot disagree. Each object gets its own colour, and its ranged GET draws a
// progress bar at the same columns as the object, so a read is visibly a read
// of those bytes and nothing else.
//
// The animation is the race. Five requests go out together: four ranged GETs
// and the background download of the whole file. A and C land first. When the
// whole file lands, B and D are still in flight and get cancelled, because
// their bytes are now on local disk.

import { AnimatedAsciiFigure } from "./AnimatedAsciiFigure.jsx";

// Container geometry. Everything below is derived from these four numbers.
const W = 48; // cells across the container
const OBJ_W = 5; // cells an object occupies
const LEFT = 14; // columns before the left rail: "  " + a 12-wide row label
const CELL0 = LEFT + 1; // the rail itself sits at LEFT, so cell 0 is one right
const OBJ = [
  { id: "A", cell: 6, at: " 16 MiB", color: "cyan" },
  { id: "B", cell: 15, at: " 40 MiB", color: "violet" },
  { id: "C", cell: 29, at: " 77 MiB", color: "amber" },
  { id: "D", cell: 40, at: "107 MiB", color: "red" },
];

const bar = (n, width) => "█".repeat(n) + "░".repeat(width - n);
const pad = (n) => " ".repeat(n);

// One marker per object, centred over its block.
function markers(glyph) {
  const segs = [];
  let col = 0;
  for (const o of OBJ) {
    const at = CELL0 + o.cell + Math.floor((OBJ_W - 1) / 2);
    segs.push(pad(at - col));
    segs.push([o.color, glyph || o.id]);
    col = at + 1;
  }
  return segs;
}

// The container itself: dotted green, with each object drawn in its own colour
// at its own cells.
function container() {
  const segs = [];
  let cell = 0;
  for (const o of OBJ) {
    if (o.cell > cell) segs.push(["green", "░".repeat(o.cell - cell)]);
    segs.push([o.color, "█".repeat(OBJ_W)]);
    cell = o.cell + OBJ_W;
  }
  if (cell < W) segs.push(["green", "░".repeat(W - cell)]);
  return segs;
}

// A ranged GET: an empty track with a progress bar parked at the object's own
// columns, so the bar can only ever fill the bytes the request asked for.
function request(o) {
  return [
    "  ",
    [`#row${o.id}`, `${o.id}  ${o.at}  `, o.color],
    "│",
    pad(o.cell),
    [`#prog${o.id}`, bar(0, OBJ_W), o.color],
    pad(W - o.cell - OBJ_W),
    "│",
    "  ",
    [`#stat${o.id}`, "in flight", o.color],
  ];
}

const LINES = [
  ["  ", ["gray", "git wants four objects out of packs/019a7f3c..bin"]],
  [],
  markers(null),
  markers("▼"),
  [pad(LEFT), "┌", "─".repeat(W), "┐"],
  ["  ", "container   ", "│", ...container(), "│"],
  [pad(LEFT), "└", "─".repeat(W), "┘"],
  [],
  ["  ", ["#hdr", "five requests go out at once:"]],
  [],
  ...OBJ.map(request),
  [
    "  ",
    ["#rowWhole", "whole file  ", "green"],
    "│",
    ["#progWhole", bar(0, W), "green"],
    "│",
    "  ",
    ["#statWhole", "in flight", "green"],
  ],
  [],
  [
    "  ",
    ["gray", "four ranged GETs race the background download of the whole"],
  ],
  ["  ", ["gray", "file. whatever the download reaches stops needing a GET."]],
];

const tagsFor = (name) => [`row${name}`, `prog${name}`, `stat${name}`];
const NAMES = [...OBJ.map((o) => o.id), "Whole"];
const LIVE = NAMES.flatMap(tagsFor);
const EVERYTHING = [...LIVE, "hdr"];

// fills: A, B, C, D out of OBJ_W, then the whole file out of W.
function at(fills, statuses) {
  const text = {};
  OBJ.forEach((o, i) => {
    text[`prog${o.id}`] = bar(fills[i], OBJ_W);
    text[`stat${o.id}`] = statuses[i];
  });
  text.progWhole = bar(fills[4], W);
  text.statWhole = statuses[4];
  return text;
}

const FLYING = [
  "in flight",
  "in flight",
  "in flight",
  "in flight",
  "in flight",
];

const STEPS = [
  {
    show: [],
    focus: [],
    caption: "128 MiB of packfile; git wants A, B, C and D out of it",
    ms: 1600,
  },
  {
    show: EVERYTHING,
    focus: LIVE,
    text: at([0, 0, 0, 0, 1], FLYING),
    caption: "four ranged GETs, plus the whole file, 2 MiB in",
    ms: 1800,
  },
  {
    focus: LIVE,
    text: at([2, 1, 2, 1, 12], FLYING),
    caption: "every GET is racing the same background download",
    ms: 900,
  },
  {
    focus: LIVE,
    text: at([3, 2, 3, 2, 22], FLYING),
    caption: "every GET is racing the same background download",
    ms: 900,
  },
  {
    focus: LIVE,
    text: at(
      [5, 3, 4, 2, 31],
      ["done", "in flight", "in flight", "in flight", "in flight"]
    ),
    caption: "A lands. the whole file is still coming",
    ms: 1400,
  },
  {
    focus: LIVE,
    text: at(
      [5, 4, 5, 3, 40],
      ["done", "in flight", "done", "in flight", "in flight"]
    ),
    caption: "C lands. B and D are still in flight",
    ms: 1400,
  },
  {
    focus: [...tagsFor("A"), ...tagsFor("C"), ...tagsFor("Whole")],
    text: at(
      [5, 4, 5, 3, 48],
      ["done", "cancelled", "done", "cancelled", "done"]
    ),
    caption: "the whole file lands, so B and D are cancelled mid-flight",
    ms: 2600,
  },
];

// xeiaso.net has no client-side JS, so this copy rests on the finished frame
// instead of step 0. No-JS readers get the complete race, not the starting gun.
const REST = {
  focus: LIVE,
  text: at(
    [5, 4, 5, 3, 48],
    ["done", "cancelled", "done", "cancelled", "done"]
  ),
  caption: "the whole file lands, so B and D are cancelled mid-flight",
};

export default function RacingTheBeam({
  label = "FIG 06",
  title = "Racing the beam: four ranged GETs while the container is 2 MiB in",
  fontSize,
}) {
  return (
    <AnimatedAsciiFigure
      label={label}
      title={title}
      lines={LINES}
      steps={STEPS}
      rest={REST}
      fontSize={fontSize}
    />
  );
}
