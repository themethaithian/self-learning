---
name: frontend-design
description: Design system + UX rules for the Self-Improve Web frontend. MUST be followed for every ticket that touches web/. Keeps all pages looking like one coherent, calm, beautiful learning app.
---

# Frontend Design System — Self-Improve Web

A calm, focused learning app used in long study sessions. Modern minimal,
light, generous whitespace, soft rounded shapes, one accent colour.
Every frontend ticket follows this file. When a rule conflicts with a ticket,
stop and ask the orchestrator.

## Identity

- Mood: quiet focus — like a good reading app, not a dashboard-y SaaS.
- Light theme only. Page background is a warm cream off-white `#faf7f2`
  (not pure white, not cool zinc-50 — pure white glares in long night
  sessions, and a cool/blue-ish near-white reads as harsh; warm tone reduces
  eye strain), surfaces a warm-nudged near-white `#fffdfa`, hairline borders
  a warm-nudged `#e6e2dc`, text `zinc-900` (headings) / `zinc-600`
  (secondary). Body copy is `zinc-800`, never pure black.
- One accent colour only: `violet-600` (actions, active nav, chart primary),
  `violet-700` on hover, `violet-50` for tinted surfaces and selected rows.
- Semantic colours: success `emerald-600`, warning `amber-600`, danger
  `rose-600` — meaning only (grades, pass/fail, streak), never decoration.
- **Every colour goes through a Tailwind theme token, never a raw class in a
  component** (`bg-surface`, `text-muted`, `border-subtle`, `bg-accent`).
  A dark theme is out of scope for v1 but must stay a token swap, not a
  rewrite of every component.

## Shape & depth

- Radius: cards and panels `rounded-2xl`, buttons/inputs/selects `rounded-xl`,
  chips and badges `rounded-full`. Nothing square except table cells.
- Depth is subtle: `shadow-sm` plus a hairline `border-zinc-200`. Never
  `shadow-lg` or larger, no gradients, no glow. Elevation is for overlays
  (modal, dropdown) only — flat surfaces on the page itself.
- Cards: `rounded-2xl border border-zinc-200 bg-white p-6 shadow-sm`.

## Typography

- Thai lesson text: `Noto Sans Thai Looped` via `next/font` — body 17–18px,
  line-height 1.8, reading column `max-w-[68ch]`. Never uppercase-transform Thai.
- UI / English: `Inter`. Code: `JetBrains Mono` (DSA editor, code in lessons).
- Hierarchy by size + weight, not colour variety: page title
  `text-xl font-semibold`, section label `text-xs font-medium text-zinc-500
  tracking-wide` (uppercase is allowed for English labels only).

## Layout

- Fixed left sidebar nav: Dashboard / Read / Drill / DSA / Test / Tickets.
  Active item = `violet-50` pill with `violet-700` text. Collapses to a bottom
  bar on mobile (desktop is primary, mobile must not break).
- Page shell: consistent header (page title left, session timer right when a
  session is active), content `max-w-5xl mx-auto px-6`.
- Spacing scale: stick to 4/6/8 Tailwind steps. Prefer whitespace over
  dividers — add a border only when whitespace alone fails to separate.

## Components (reuse, don't reinvent per page)

- `Card`, `StatTile` (big number + label + optional delta arrow),
  `TrendChart`, `RadarChart`, `Timer`, `Button` (primary/ghost/danger),
  `EmptyState` (icon + one sentence + one action). New variant → extend the
  component, never fork a page-local copy.
- Charts: primary series violet, comparison series `zinc-400`, grid lines
  `zinc-200`, no chart borders. Every measurement view shows current value +
  trend. Empty/1-point data gets a designed EmptyState, not a broken chart.

## Interaction & states

- Every async view designs all four states: loading (skeleton, not spinner-only),
  empty, error (with retry), success.
- LLM grading takes seconds: show a calm "กำลังตรวจ..." state, disable
  double-submit, never lose the user's typed answer on error.
- Keyboard-first: Ctrl/Cmd+Enter submits forms; Drill uses 1–5 for grades,
  Space to reveal. Show hints inline (`kbd` styling).
- Motion: 150ms ease transitions on hover/expand only. Timer is visible but
  calm — no red flashing; last-5-min = subtle amber text change.

## Accessibility & quality gate (check before ticket done)

- [ ] Contrast AA on all text (`zinc-500` on the cream page `#faf7f2` is the
      lightest text allowed — verified ≈4.5:1, AA; violet text on white must
      be `violet-600` or darker)
- [ ] Focus rings visible (`focus-visible:ring-2 ring-violet-500 ring-offset-2`)
- [ ] Colour is never the only signal (pass/fail also carries an icon or label)
- [ ] Works at 375px width and 1440px; no horizontal scroll
- [ ] All four async states implemented
- [ ] Thai text renders with correct font + line-height
- [ ] No raw colour class in a component — theme tokens only
- [ ] `npm run build` (static export) passes
