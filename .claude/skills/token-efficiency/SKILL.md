---
name: token-efficiency
description: Rules for spending tokens efficiently — applies to the orchestrator, ALL subagents, and the app's runtime LLM calls. Load at the start of multi-step work in this repo.
---

# Token Efficiency Rules

Two budgets to protect: (1) Claude Code subscription quota (orchestrator +
subagents), (2) Anthropic API spend at app runtime (Haiku). Rules for both.

## Orchestrator (main session)

- Pass subagents the MINIMUM context: the ticket text + relevant file paths +
  acceptance criteria. Never paste whole design docs — point to
  `docs/design.md` sections and let the agent read what it needs.
- Don't re-read files you already have in context; don't re-verify what a
  tool result already confirmed.
- One ticket at a time. Don't preload future tickets into context.
- Match model tier to task: haiku for mechanical transforms, sonnet for
  implementation, opus for review/lesson content, fable ONLY via the
  escalation path in CLAUDE.md. Never use a bigger model where a smaller
  one passes the same acceptance criteria.
- Once PR CI exists, treat it as the source of truth for vet/test/build —
  don't re-run locally what CI already proves, and don't paste long
  passing-test output into reports (link/summarize).
- Prefer summary + quiz over exhaustive prose for the user (see CLAUDE.md
  Active learning) — it's fewer tokens AND better retention.
- Content batches: one chapter/pattern per session, ~4–8 items max.
- If a subagent fails twice on the same task, STOP (per CLAUDE.md) — do not
  burn a third attempt without the user's decision.

## All subagents

- Search before reading: Glob/Grep to locate, then Read only the relevant
  ranges of the files you need — not whole directories.
- Report back summaries + diffs, never full file contents the orchestrator
  already has or can read from disk.
- Run quiet commands: `go test ./...` without `-v` (add `-run X -v` only on
  the failing test), no verbose npm output.
- Don't regenerate unchanged code — edit in place with the smallest diff.
- code-reviewer: review `git diff` of the ticket, not the whole repo; read
  surrounding files only when the diff genuinely requires it.
- lesson-writer/verifier: one concept per invocation; on FAIL-retry, send
  only the issue list + the failed lesson, not the whole batch history.

## App runtime (Haiku API — costs real money per request)

- Model: `claude-haiku-4-5` only. Set an explicit low `max_tokens` per use
  case (recall grading ~500, DSA review ~1500).
- Prompts include only what the task needs: the ONE lesson/problem at hand +
  the user's answer. Never dump full lesson history or curriculum context.
- Ask for JSON-only output (no prose preamble) — smaller responses, no
  parsing retries.
- Reuse stable system prompts verbatim across requests (enables prompt
  caching server-side; changing one word breaks the cache).
- Never call the LLM in a loop per item when one batched call can grade the
  same payload; never retry non-retryable errors (400s).
