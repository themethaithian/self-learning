# Self-Improve Web — Project Memory

Personal learning web app (single user: me) to train 4 skills for a job change:
system design (DDD + Distributed Systems concepts), AWS (SAA-C03), deep Golang,
DSA interview prep. The app itself is part of the curriculum: I am learning by
building it.

## Working with the user

- Converse in THAI; keep technical terms in English (no transliterated
  jargon). Project docs are Thai, app UI copy is English, lessons are Thai.

## Orchestration rules (IMPORTANT)

- The main session (you) is the ORCHESTRATOR ONLY. Delegate ALL work:
  - content generation → `lesson-writer`, then ALWAYS `lesson-verifier`
  - code implementation → `go-implementer`
  - reviewing diffs → `code-reviewer`
- You may implement something yourself ONLY when a subagent has failed twice
  on the same task, AND you must first ask the user for explicit permission,
  stating what failed and why.
- Model policy: Fable 5 runs the main session as ORCHESTRATOR ONLY — it
  plans, delegates, and reviews subagent results, never writes lessons or
  code itself. All worker subagents run on opus, sonnet, or haiku via the
  `model` field in their `.claude/agents/*.md` frontmatter
  (lesson-writer: opus, lesson-verifier: sonnet, go-implementer: sonnet,
  code-reviewer: opus; use haiku for cheap mechanical tasks). Never launch
  a subagent on Fable 5.
- Never set CLAUDE_CODE_SUBAGENT_MODEL (it would override per-agent models).
- Work in small tickets. One ticket at a time. Show a short summary after each.
- Git workflow: every ticket on branch `ticket/<id>-<slug>` → PR to
  `develop` using `.github/pull_request_template.md`. The PR body must be
  self-sufficient for PHONE review: Thai walkthrough (what/why per file),
  the ticket's Review focus checklist, vet/test output, code-reviewer
  verdict. The user reviews and merges (often via GitHub mobile). Merging
  to develop auto-deploys to the VPS. Never start the next ticket before
  the current PR is merged.
- The plan lives in `docs/roadmap.md` (big picture + ticket index) and
  `docs/tickets/week-*.md` (small per-ticket detail with Review focus).
  Full design: `docs/design.md`. Keep ticket Status lines and roadmap
  checkboxes updated as work completes.
- Skills (MANDATORY): every ticket touching `web/` follows
  `.claude/skills/frontend-design/SKILL.md`; the orchestrator and ALL
  subagents follow `.claude/skills/token-efficiency/SKILL.md`.

## Tech constraints (non-negotiable)

- Backend: Go 1.26+, standard library net/http ONLY (no web framework).
  Hand-written middleware (logging, recovery, auth) — this is a learning goal.
- Architecture: DDD layering — domain / application / infrastructure.
  The domain layer imports NO vendor SDKs (no AWS SDK, no Anthropic SDK).
- DB: MySQL 8 in Docker (both locally and on the VPS).
- Frontend: Next.js + Tailwind, static export (`output: 'export'`),
  served by Caddy on the VPS, calling the Go API. UI copy in ENGLISH.
- Deploy: DigitalOcean VPS (Singapore, ~$6/mo), Docker Compose
  (mysql + api + caddy with auto-HTTPS). Secrets in `.env` on the server
  (chmod 600). Deploy via GitHub Actions over SSH. Nightly mysqldump backup
  to object storage. NO AWS hosting — AWS (SAA-C03) is curriculum content
  only; hands-on labs happen in a separate free-tier sandbox account.
- LLM integration: hand-written Anthropic API client using net/http in
  `internal/platform/llm` (no vendor SDK, no agent framework).
- Runtime LLM: Anthropic API, `claude-haiku-4-5` for grading/quizzes.
  Context assembled per request from MySQL (notes, answer history).

## Content model

- No book scanning. A curriculum tree (topic → chapter → concept) based on the
  tables of contents of "Domain-Driven Design" and "Distributed Systems",
  plus AWS SAA-C03 domains, plus Go internals topics, plus DSA patterns.
- Lessons are generated OFFLINE by the Claude Code pipeline (subscription),
  saved as JSON under `content/lessons/<topic>/<concept>.json`,
  imported into MySQL by `cmd/import-lessons`.
- Lesson language: THAI content, English technical terms kept. 5–10 min read
  per chunk. Every chunk ends with recall-check questions.
- Lessons include mermaid diagrams (inside body_md) whenever a flow,
  architecture, or state machine helps, plus 2–4 `references` to primary
  sources (official docs, book chapters). Never invent URLs.
- Adding lessons later is always possible and fully additive: add concepts
  to `content/curriculum/*.json` → `import-curriculum` → run a
  writer/verifier batch → `import-lessons`.
- Never reproduce book text verbatim; lessons come from general knowledge of
  these well-known works. AWS lessons must be grounded in official AWS docs
  excerpts provided in context.

## Code style (IMPORTANT)

- Minimal comments. Code must be self-documenting: clear names, small
  functions, obvious control flow. A comment is allowed ONLY for a
  constraint the code cannot express (a non-obvious WHY, an invariant,
  a spec reference like the SM-2 formula). Never comment WHAT the code
  does, never add section-banner comments, never leave TODO noise.

## Quality rules

- Every lesson passes lesson-verifier before being saved. FAIL → regenerate.
- Domain layer requires table-driven tests. `go vet` + `go test ./...` must
  pass before a ticket is done.
- Generate lessons in small batches (one topic at a time) to respect
  subscription usage limits.
