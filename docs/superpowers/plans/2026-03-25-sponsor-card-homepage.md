# Replace GitHub Sponsor iframe with Styled Sponsor Card

## Context

The homepage (`lume/src/index.jsx` line 40) currently embeds a GitHub Sponsors iframe. This should be replaced with a custom-styled card linking to Patreon, GitHub Sponsors, and the Sponsor Panel (sponsors.xeiaso.net). The card design should echo the sponsor-panel app's warm Gruvbox aesthetic.

## File to modify

- `lume/src/index.jsx` — the only file that needs changes

## Design reference (from `cmd/sponsor-panel`)

- Cards: `rounded-2xl`, border, surface bg, `overflow-hidden`, subtle shadow
- Gradient accent line: 2px at top, orange-to-pink for warm variant
- Card titles: `font-serif`, semibold
- Buttons: `rounded-xl`, colored bg, hover lift (`hover:-translate-y-px`), white text

## Plan

### 1. Add small icon components before the default export

Copy SVG paths from `lume/src/donate.jsx` (GitHubIcon, PatreonIcon) scaled to 20x20. Add a star icon for the Sponsor Panel link.

### 2. Add `SponsorCard` component

Structure:

```
<div>  — card container (rounded-2xl, border, bg-bg-2, dark:bg-bgDark-2, shadow-sm, overflow-hidden, max-w-xl, mx-auto, my-6)
  <div> — gradient accent line (2px, orange→purple, with dark mode variant via dark:hidden/hidden dark:block)
  <div> — content area (px-6 py-5 text-center)
    <h3> — "Support My Work" with heart SVG icon
    <p>  — brief description (text-sm, muted color)
    <div> — flex row of 3 button-links (flex-wrap, gap-3, centered)
      <a> Patreon      — bg-orange-light / dark:bg-orangeDark-light
      <a> GitHub Sponsors — bg-fg-0 / dark:bg-purpleDark-light
      <a> Sponsor Panel — bg-purple-light / dark:bg-blueDark-light
```

### 3. Replace the iframe (line 40) with `<SponsorCard />`

## Key implementation details

**Gradient accent line:** Use two `<div>`s with `dark:hidden` / `hidden dark:block` to swap light/dark gradients, since inline `style` can't respond to `prefers-color-scheme`:

- Light: `linear-gradient(90deg, #d65d0e, #b16286)` (orange→purple)
- Dark: `linear-gradient(90deg, #fe8019, #d3869b)` (bright orange→pink)

**Button link style overrides:** The site's base `<a>` styles (in `lume/src/styles.css` line 61-63 and `hack.css` line 65-71) apply link colors, underline, border-bottom, and visited styles. Each button `<a>` needs:

- `no-underline border-0` to remove text decoration and bottom border
- `text-white hover:text-white hover:bg-[color]` to override link/hover colors
- `visited:text-white visited:hover:text-white visited:hover:bg-[color]` to override visited states

**Color tokens available** (from `lume/tailwind.config.js`):

- `bg-orange-light`, `dark:bg-orangeDark-light`, `hover:bg-orange-dark`
- `bg-purple-light`, `dark:bg-blueDark-light`, `hover:bg-purple-dark`
- `bg-fg-0`, `dark:bg-purpleDark-light`, `hover:bg-fg-1`
- `text-fg-0`, `dark:text-fgDark-1`, `text-fg-3`, `dark:text-fgDark-3`

## Verification

1. Run `npm run dev` and check the homepage in both light and dark mode
2. Verify all three links open correctly in new tabs
3. Confirm the card is responsive (shrinks gracefully on mobile)
4. Check that button hover states work (lift effect, color change, no link underline artifacts)
