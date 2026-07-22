---
name: frontend-design
description: Design system + UX rules for the Self-Improve Web frontend. MUST be followed for every ticket that touches web/. Keeps all pages looking like one coherent, calm, beautiful learning app.
---

# Frontend Design System — Self-Improve Web

A calm, focused learning app used in long study sessions (often at night).
Every frontend ticket follows this file. When a rule conflicts with a ticket,
stop and ask the orchestrator.

## Identity

- Mood: quiet focus — like a good reading app, not a dashboard-y SaaS.
- Dark-first: default dark theme. Base = Tailwind `slate` (bg `slate-950`,
  surface `slate-900`, borders `slate-800`, text `slate-200`/`slate-400`).
- One accent color only: `indigo-400/500` (actions, active nav, chart primary).
- Semantic colors: success `emerald-400`, warning `amber-400`, danger `rose-400`.
  Used for meaning only (grades, pass/fail, streak) — never decoration.

## Typography

- Thai lesson text: `Noto Sans Thai Looped` via `next/font` — body 17–18px,
  line-height 1.8, reading column `max-w-[65ch]`. Never uppercase-transform Thai.
- UI / English: `Inter`. Code: `JetBrains Mono` (DSA editor, code in lessons).
- Hierarchy by size + weight, not color variety: page title `text-xl font-semibold`,
  section `text-sm font-medium text-slate-400 uppercase tracking-wide`.

## Layout

- Fixed left sidebar nav: Dashboard / Read / Drill / DSA / Test / Tickets.
  Collapses to bottom bar on mobile (desktop is primary, mobile must not break).
- Page shell: consistent header (page title left, session timer right when a
  session is active), content `max-w-5xl mx-auto px-6`.
- Spacing scale: stick to 4/6/8 Tailwind steps; cards `rounded-xl border
  border-slate-800 bg-slate-900 p-6`. No shadows in dark theme — borders separate.

## Components (reuse, don't reinvent per page)

- `Card`, `StatTile` (big number + label + optional delta arrow),
  `TrendChart`, `RadarChart`, `Timer`, `Button` (primary/ghost/danger),
  `EmptyState` (icon + one sentence + one action). New variant → extend the
  component, never fork a page-local copy.
- Charts: chart primary = indigo, comparison series = slate-500, grid lines
  `slate-800`, no chart borders. Every measurement view shows current value +
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

- [ ] Contrast AA on all text (slate-400 on slate-950 is the minimum allowed)
- [ ] Focus rings visible (`focus-visible:ring-2 ring-indigo-500`)
- [ ] Works at 375px width and 1440px; no horizontal scroll
- [ ] All four async states implemented
- [ ] Thai text renders with correct font + line-height
- [ ] `npm run build` (static export) passes
