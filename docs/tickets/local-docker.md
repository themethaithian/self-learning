# Infra — Local Docker dev stack

Ad hoc infra ticket (not part of the week-1..6 content plan): the project
owner needs to open the app on their own machine with one command, without
installing Go/Node/MySQL locally.

---

## T-local-docker — One-command local dev stack `[go-implementer]`
- **Scope**: root `docker-compose.yml` grows from a bare `mysql` service to
  a full local stack (`mysql` + `api` + `web` + one-shot `seed`), reusing
  the existing prod `Dockerfile`/`web/Dockerfile`/`deploy/Caddyfile`
  unmodified (Caddyfile is bind-mounted read-only, not edited); `Makefile`
  gains `dev`/`dev-down`/`dev-reset`/`seed`/`logs`; `README.md` gains a
  "รันเองบนเครื่อง" section; two accidentally-committed scratch files
  removed.
- **Acceptance**: `docker compose up -d --build` from a clean clone (no
  `.env` needed) ends with curriculum + lessons queryable through the
  browser-facing URL; rerunning is idempotent (no duplicate/corrupted
  data); `deploy/docker-compose.prod.yml` and `deploy/Caddyfile` untouched
  and still parse.
- **Review focus**:
  1. ทำไม `web` service ต้อง mount `deploy/Caddyfile` แทนที่จะ bake
     `NEXT_PUBLIC_API_BASE` เป็น absolute URL ตรง ๆ?
  2. ทำไม service `seed` ต้อง bind-mount `./content` แทนที่จะพึ่ง
     `COPY content/` ที่มีอยู่แล้วใน `Dockerfile`?
- Status: `implemented, PR pending` — three code-reviewer passes on branch
  `ticket/local-docker-dev-stack`, vet/test green; PR not yet opened

---

## T-local-access — HSTS lockout + stale-asset caching hotfix `[go-implementer]`

### Symptom
Normal (non-incognito) Chrome could not reach the app at all: every page
loaded, but every `fetch('/api/v1/...')` from `web/lib/api.ts` failed with
`ApiError(0, "network request failed — is the API reachable?")` before any
HTTP response came back. Incognito worked fine. The browser never reached
`/token` even though there was no session — a 401 from the API would have
redirected there, so never getting that far proves the request died on the
client, before it ever reached the API.

### Root cause
`deploy/Caddyfile` sent `Strict-Transport-Security` unconditionally, even
when Caddy was serving plain HTTP on `DOMAIN=:80` for local dev. Chrome
honors HSTS on first response and caches it per-host: once it saw that
header from `http://localhost:3000`, it force-upgraded every subsequent
same-origin request to `https://localhost:3000`, where nothing speaks TLS.
`fetch` fails at the TCP/TLS layer, never producing a response for
`web/lib/api.ts` to read. Incognito worked because Chrome's HSTS store is
per-profile and in-memory there, so it started clean.

### Why deleting the HSTS entry alone did not help
Every response — including the one serving the still-broken page — kept
re-sending `Strict-Transport-Security: max-age=31536000; includeSubDomains`.
Clearing the entry in `chrome://net-internals/#hsts` only stuck once the
server stopped sending the header; otherwise the very next request
reinstated the policy.

### Fix
Scope the header behind an `@https protocol https` matcher so it is only
ever sent over a connection that is actually TLS. `DOMAIN=:80` (local) never
matches and stays HSTS-free; a real domain still gets the header once Caddy
auto-negotiates HTTPS.

### Manual recovery for a browser that already cached the policy
`chrome://net-internals/#hsts` → "Delete domain security policies" → enter
`localhost` → Delete, then hard-reload. Needed once per affected browser
profile; the header itself no longer gets re-sent after this fix.

### F1 — missing `/_next/*` assets returned `200 text/html` instead of 404
The SPA fallback (`try_files {path} ... /index.html`) also caught requests
under Next's build output, so a stale `index.html` referencing a chunk that
no longer exists (after a rebuild) served a 200 `text/html` response for a
`.js` request instead of a 404. Combined with `X-Content-Type-Options:
nosniff`, Chrome silently refuses to execute it — a dead app with nothing
in the network tab to point at. Fixed by giving `/_next/*` its own `handle`
block with plain `file_server` (no `try_files`), so a missing asset 404s
like it should; real app routes (`/today`, `/learn`, `/lesson`, `/token`,
...) still fall through the unchanged SPA-fallback `handle` below it.

### F2 — no `Cache-Control` on any response
Nothing set `Cache-Control`, so a rebuilt app could keep serving
browser-cached assets from an old build indefinitely. `/_next/static/*`
filenames are content-hashed by Next, so they get `public,
max-age=31536000, immutable`; everything served by the SPA-fallback handle
(HTML documents and other unhashed files) gets `no-cache` so a rebuild is
picked up on the next request instead of being blanket `no-store`d.

- Status: `done` — verified against the running local stack: `curl` shows
  HSTS absent / security headers present / correct `Cache-Control` on both
  a missing and a real `/_next/static/*` asset (404 vs 200), the API still
  answers with a bearer token, headless Chromium lands on `/learn` with all
  requests 200 and no console output, and `caddy validate` passes for both
  `DOMAIN=:80` and a real hostname.
