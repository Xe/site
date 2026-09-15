// Shared scaffolding for the ASCII diagrams in this post.
// Standalone React helpers. Inline styles only, React is the only dependency.

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
  fontFamily: "ui-monospace, SFMono-Regular, Menlo, Consolas, monospace",
  fontSize: "13px",
  whiteSpace: "pre",
  textAlign: "left",
};

export const COLORS = {
  green: "#4ade80",
  red: "#f87171",
  gray: "#64748b",
  // Added for figures that need to tell more than two things apart. All
  // Tailwind 400-weight like the three above, so they sit at the same
  // brightness against #0f172a, and the same set AnimatedAsciiFigure uses.
  amber: "#fbbf24",
  cyan: "#38bdf8",
  violet: "#a78bfa",
  base: "#cbd5e1",
};

// A line is an array of segments; a segment is either a plain string or a
// [colorKey, text] pair. Concatenated segment text preserves the exact
// monospace alignment of the source art.
export function AsciiFigure({ label, title, lines, fontSize }) {
  const preStyle = fontSize ? { ...PRE_STYLE, fontSize } : PRE_STYLE;
  return (
    <figure
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
          fontFamily: PRE_STYLE.fontFamily,
        }}
      >
        <span
          style={{ fontSize: 11, color: "#475569", letterSpacing: "0.14em" }}
        >
          {label}
        </span>
        <span style={{ fontSize: 12, color: "#94a3b8" }}>{title}</span>
      </figcaption>
      <pre style={preStyle}>
        {lines.map((segments, i) => (
          <div key={i}>
            {segments.length === 0
              ? " "
              : segments.map((seg, j) =>
                  typeof seg === "string" ? (
                    seg
                  ) : (
                    <span key={j} style={{ color: COLORS[seg[0]] }}>
                      {seg[1]}
                    </span>
                  )
                )}
          </div>
        ))}
      </pre>
    </figure>
  );
}
