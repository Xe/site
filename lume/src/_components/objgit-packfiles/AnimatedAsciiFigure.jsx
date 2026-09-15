// Shared scaffolding for the *animated* ASCII diagrams.
// Companion to AsciiFigure: same type, same panel, same palette. The only
// difference is that a figure declares steps, and a step decides which of the
// tagged segments are lit, dimmed, or not there yet.
//
// Standalone React helpers. Inline styles only, React is the only dependency.
//
// The whole diagram is always in the DOM. A segment that has not been revealed
// yet is drawn at opacity 0, never replaced by spaces, which means:
//   - columns never shift, so hand-aligned art stays aligned;
//   - the server-rendered HTML holds the complete figure, so a reader without
//     JavaScript (or a crawler, or a screen reader) still gets all of it;
//   - reveals can cross-fade instead of snapping.

import {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "npm:preact/hooks";

const FONT_STACK = "ui-monospace, SFMono-Regular, Menlo, Consolas, monospace";

const PRE_STYLE = {
  margin: "0",
  // The blog layout caps .prose pre at 85ch; the art sets its own width.
  maxWidth: "none",
  background: "#0f172a",
  border: "1px solid #16202f",
  borderRadius: "6px",
  padding: "32px",
  overflowX: "auto",
  lineHeight: "1.3",
  color: "#cbd5e1",
  fontFamily: FONT_STACK,
  fontSize: "13px",
  whiteSpace: "pre",
  textAlign: "left",
};

// The three AsciiFigure colors, plus three more for figures that need to tell
// more than two things apart. All Tailwind 400-weight, so they sit at the same
// brightness as the originals against #0f172a.
export const COLORS = {
  green: "#4ade80",
  red: "#f87171",
  gray: "#64748b",
  amber: "#fbbf24",
  cyan: "#38bdf8",
  violet: "#a78bfa",
  base: "#cbd5e1",
};

// Light enough to stay readable once a step has moved on, dark enough that the
// focused segments are obviously the ones being talked about.
const DIM = "#5c6b7f";
const TRANSITION = "color 300ms ease, opacity 300ms ease";

// A line is an array of segments. A segment is one of:
//   "text"                     plain structure: always visible, never dimmed
//   ["green", "text"]          statically colored, always visible
//   ["#id", "text", "green"]   tagged: visibility and color follow the steps
// A tag is any key beginning with "#". The third element is the color the
// segment takes when a step focuses it, defaulting to base.
function segmentParts(seg) {
  if (typeof seg === "string") {
    return { tag: null, text: seg, color: null, plain: true };
  }
  const [head, text, color] = seg;
  if (typeof head === "string" && head.startsWith("#")) {
    // Steps name ids without the "#", which is only there to tell a tag apart
    // from a static color key.
    return { tag: head.slice(1), text, color: color || "base", plain: false };
  }
  return { tag: null, text, color: head, plain: true };
}

// Per-step text substitutions (progress bars, counters, byte ranges) are padded
// or clipped to the original segment's width so the art cannot reflow.
function fit(text, width) {
  if (text.length === width) return text;
  if (text.length > width) return text.slice(0, width);
  return text + " ".repeat(width - text.length);
}

function useReducedMotion() {
  const [reduced, setReduced] = useState(false);
  useEffect(() => {
    if (typeof window === "undefined" || !window.matchMedia) return undefined;
    const mq = window.matchMedia("(prefers-reduced-motion: reduce)");
    const sync = () => setReduced(mq.matches);
    sync();
    mq.addEventListener("change", sync);
    return () => mq.removeEventListener("change", sync);
  }, []);
  return reduced;
}

export function AnimatedAsciiFigure({
  label,
  title,
  lines,
  steps,
  fontSize,
  // Ids revealed in an earlier step stay revealed. Every figure here reads as
  // an accumulation, so this is the default; pass false for a figure whose
  // steps are alternatives rather than a sequence.
  cumulative = true,
  // A looping figure never finishes: it wraps from the last step back to the
  // first and keeps going. The control row turns into play/pause, because
  // "replay" means nothing to something that never stopped.
  loop = false,
  // The frame a paused figure rests on: step-shaped ({ focus, text, caption }),
  // with everything the steps ever reveal already visible. A cycling figure
  // shows one thing at a time, so its resting frame is the whole picture, which
  // is also what a reader who hits pause is asking to see. It renders on the
  // server too, so no-JS readers and crawlers get the complete figure.
  rest = null,
  // Figures that rest start there and wait to be played.
  autoPlay = rest === null,
  stepMs = 2200,
}) {
  const reduced = useReducedMotion();
  const lastStep = steps.length - 1;

  const containerRef = useRef(null);
  const [step, setStep] = useState(0);
  const [playing, setPlaying] = useState(false);
  const [finished, setFinished] = useState(false);
  // Resting is its own state, not just "not playing": clicking a dot pauses on
  // that specific step, which is the opposite of resting on all of them.
  const [resting, setResting] = useState(rest !== null);

  // Resolve a step to the set of ids that are visible and the set that is lit.
  const resolved = useMemo(() => {
    const shown = new Set();
    const out = [];
    for (const s of steps) {
      if (!cumulative) shown.clear();
      for (const id of s.show || []) shown.add(id);
      const focus = new Set(s.focus || s.show || []);
      out.push({
        shown: new Set(shown),
        focus,
        caption: s.caption || "",
        text: s.text || null,
        ms: s.ms,
      });
    }
    return out;
  }, [steps, cumulative]);

  // Everything the steps ever reveal, lit unless the resting frame is choosier.
  const restFrame = useMemo(() => {
    if (!rest) return null;
    const shown = new Set();
    for (const s of steps) for (const id of s.show || []) shown.add(id);
    return {
      shown,
      focus: new Set(rest.focus || shown),
      caption: rest.caption || "",
      text: rest.text || null,
    };
  }, [rest, steps]);

  const play = useCallback(() => {
    setStep(0);
    setResting(false);
    setFinished(false);
    setPlaying(true);
  }, []);

  // Looping figures resume in place instead of restarting, and pausing one
  // drops it back on its resting frame.
  const toggle = useCallback(() => {
    if (playing) {
      setPlaying(false);
      setResting(rest !== null);
    } else {
      setPlaying(true);
      setResting(false);
    }
  }, [playing, rest]);

  // Start once, when enough of the figure is on screen to be worth watching.
  const startedRef = useRef(false);
  useEffect(() => {
    const node = containerRef.current;
    if (!node || typeof IntersectionObserver === "undefined") return undefined;
    if (!autoPlay) return undefined;
    if (reduced) {
      // No motion: hand the reader the finished figure and a button.
      setStep(lastStep);
      setFinished(true);
      return undefined;
    }
    const io = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          if (entry.isIntersecting && !startedRef.current) {
            startedRef.current = true;
            play();
          }
        }
      },
      { threshold: 0.35 }
    );
    io.observe(node);
    return () => io.disconnect();
  }, [play, reduced, lastStep, autoPlay]);

  // Advance.
  useEffect(() => {
    if (!playing) return undefined;
    if (!loop && step >= lastStep) {
      setPlaying(false);
      setFinished(true);
      return undefined;
    }
    const hold = resolved[step]?.ms ?? stepMs;
    const t = setTimeout(
      () => setStep((s) => (s >= lastStep ? 0 : s + 1)),
      hold
    );
    return () => clearTimeout(t);
  }, [playing, step, lastStep, resolved, stepMs, loop]);

  const atRest = restFrame !== null && resting && !playing;
  const frame = atRest ? restFrame : resolved[Math.min(step, lastStep)];

  const rendered = useMemo(
    () =>
      lines.map((segments, i) => (
        <div key={i}>
          {segments.length === 0
            ? " "
            : segments.map((seg, j) => {
                const { tag, text, color, plain } = segmentParts(seg);
                if (plain) {
                  return (
                    <span
                      key={j}
                      style={{
                        color: color ? COLORS[color] : undefined,
                        transition: TRANSITION,
                      }}
                    >
                      {text}
                    </span>
                  );
                }
                const visible = frame.shown.has(tag);
                const lit = frame.focus.has(tag);
                const substituted = frame.text && frame.text[tag];
                return (
                  <span
                    key={j}
                    style={{
                      color: lit ? COLORS[color] : DIM,
                      opacity: visible ? 1 : 0,
                      transition: TRANSITION,
                    }}
                  >
                    {substituted ? fit(substituted, text.length) : text}
                  </span>
                );
              })}
        </div>
      )),
    [lines, frame]
  );

  const preStyle = fontSize ? { ...PRE_STYLE, fontSize } : PRE_STYLE;
  const stepNo = String(Math.min(step, lastStep) + 1).padStart(2, "0");
  const stepTotal = String(steps.length).padStart(2, "0");

  return (
    <figure
      ref={containerRef}
      style={{
        // Size to the art with no cap, then center in the prose column.
        margin: "0 auto 1rem",
        width: "max-content",
        display: "flex",
        flexDirection: "column",
        gap: 10,
        textAlign: "left",
      }}
    >
      <figcaption
        style={{
          display: "flex",
          alignItems: "baseline",
          gap: 14,
          fontFamily: FONT_STACK,
        }}
      >
        <span
          style={{ fontSize: 11, color: "#475569", letterSpacing: "0.14em" }}
        >
          {label}
        </span>
        <span style={{ fontSize: 12, color: "#94a3b8" }}>{title}</span>
      </figcaption>

      <pre style={preStyle}>{rendered}</pre>

      <div
        style={{
          display: "flex",
          alignItems: "center",
          gap: 12,
          fontFamily: FONT_STACK,
          fontSize: 12,
          color: "#94a3b8",
          minHeight: "1.4em",
        }}
      >
        <span
          style={{
            color: "#475569",
            fontSize: 11,
            letterSpacing: "0.14em",
            flexShrink: 0,
          }}
        >
          {atRest ? `ALL/${stepTotal}` : `${stepNo}/${stepTotal}`}
        </span>

        <span
          style={{ display: "flex", gap: 4, flexShrink: 0 }}
          aria-hidden="true"
        >
          {steps.map((s, i) => (
            <button
              key={i}
              type="button"
              onClick={() => {
                setPlaying(false);
                setResting(false);
                setFinished(true);
                setStep(i);
              }}
              title={s.caption}
              style={{
                width: 14,
                height: 4,
                padding: 0,
                border: "none",
                borderRadius: 2,
                cursor: "pointer",
                background: (atRest ? true : loop ? i === step : i <= step)
                  ? "#4ade80"
                  : "#243147",
                transition: "background 300ms ease",
              }}
            />
          ))}
        </span>

        <span style={{ flex: 1, transition: "opacity 200ms ease" }}>
          {frame.caption}
        </span>

        <button
          type="button"
          onClick={loop ? toggle : play}
          style={{
            flexShrink: 0,
            background: "transparent",
            border: "1px solid #243147",
            borderRadius: 4,
            color: finished || !playing ? "#94a3b8" : "#475569",
            fontFamily: FONT_STACK,
            fontSize: 11,
            padding: "2px 8px",
            cursor: "pointer",
          }}
        >
          {loop ? (playing ? "❚❚ pause" : "▶ play") : "↺ replay"}
        </button>
      </div>
    </figure>
  );
}

export default AnimatedAsciiFigure;
