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

### Root cause — mechanism unconfirmed, two hypotheses
Chrome had cached an HSTS policy for the host `localhost`, so every later
`fetch('/api/v1/...')` got force-upgraded to `https://localhost:3000`,
where nothing speaks TLS — `fetch` fails at the TLS layer, before
`web/lib/api.ts` ever sees a response. Incognito worked because Chrome's
HSTS store is per-profile and in-memory there, so it started clean. That
part is solid. **How the policy got installed in the first place is not
confirmed** — a code-reviewer pass on the first version of this fix showed
the original story doesn't hold up, so both candidates are recorded here
rather than one confident-sounding guess:

- **H1 (the original theory, now doubted).** `deploy/Caddyfile` sent
  `Strict-Transport-Security` unconditionally, including over the plain
  HTTP that local dev (`DOMAIN=:80`) actually serves. Problem: RFC 6797
  requires a compliant UA to ignore an STS header received over non-secure
  transport, and Chrome does gate HSTS handling on the connection actually
  being TLS. It's also not supported by this repo's history — `git log`
  shows the local dev stack's `DOMAIN` was always an explicitly-safe value
  (`http://localhost`, then `:80` since commit `5fb8ab6`), never a bare
  hostname on the port the browser was actually hitting. A spec-compliant
  Chrome should not have installed HSTS from that traffic at all.
- **H2 (spec-consistent, better supported).** `deploy/docker-compose.prod.yml`
  defaulted the `caddy` service's `DOMAIN` to bare `localhost` when unset
  (fixed below, see "DOMAIN default hardening"). Running that prod compose
  file locally without an explicit `DOMAIN` — plausible while working on
  the prod-compose ticket — makes Caddy treat `localhost` as a real
  hostname eligible for automatic HTTPS, issue an internal-CA cert, and
  send a **spec-legitimate** HSTS header over that real TLS connection.
  Chrome caching that is completely standard, and `includeSubDomains` would
  poison every `*.localhost` in the profile along with it. Under this
  hypothesis, the plain-HTTP dev stack on port 3000 never installed
  anything; what actually unblocked the browser afterward was the user's
  own `chrome://net-internals` deletion below, not this Caddyfile change.

Which one actually happened is not known. Both are written down so the
next person who hits this searches the repo and finds the honest state,
not a tidy story that doesn't survive a fact-check.

### Why deleting the HSTS entry alone did not help
Under the original incident, clearing the entry in
`chrome://net-internals/#hsts` did not fix the browser. The likely
explanation now (H2) is that the deletion happened while a real TLS source
was still capable of re-installing the policy — for instance the
prod-compose `caddy` container still running in the background — rather
than plain-HTTP responses "reinstating" a policy no compliant browser
should have accepted from them in the first place (H1's claim, now
doubted; see above).

### Fix
`deploy/Caddyfile` scopes the header behind an `@https protocol https`
matcher so it is only ever sent over a connection that is actually TLS.
`DOMAIN=:80` (local) never matches and stays HSTS-free; a real domain still
gets the header once Caddy auto-negotiates HTTPS. This stays correct and
necessary regardless of which hypothesis above explains the past incident:
RFC 6797 requires HSTS to be scoped to secure transport, so the plain-HTTP
dev stack must never send it, whether or not it's what caused this
particular lockout.

### DOMAIN default hardening (R2)
`deploy/docker-compose.prod.yml`'s `caddy` service used to default
`DOMAIN` to `localhost` when unset — precisely the H2 footgun above. It now
has no fallback (`${DOMAIN:?set DOMAIN in .env to the real deploy
hostname}`), so `docker compose -f deploy/docker-compose.prod.yml up` with
`DOMAIN` unset fails immediately with that message instead of silently
starting a local HTTPS-with-HSTS server. The Caddyfile's own comment about
`{$DOMAIN}` was also wrong in the same way — it claimed `DOMAIN=localhost`
serves plain HTTP like `:80` does; it does not, Caddy still treats
`localhost` as a real hostname for automatic HTTPS. The comment now says so
directly and points here.

### Manual recovery for a browser that already cached the policy
`chrome://net-internals/#hsts` → "Delete domain security policies" → enter
`localhost` → Delete, then hard-reload. Works and is needed under either
hypothesis above. Local dev (`DOMAIN=:80`) never sends the header after
this fix either way; the prod-compose `DOMAIN` fix above closes the H2
path, so nothing should be able to re-install it against `localhost` again.

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

`/srv/404.html` is generated by `next build` but is never actually served —
the SPA fallback always serves `/index.html` for an unmatched app route, so
that file is dead output. Noted here so nobody assumes it's wired up.

### F2 — no `Cache-Control` on any response
Nothing set `Cache-Control`, so a rebuilt app could keep serving
browser-cached assets from an old build indefinitely. `/_next/static/*`
filenames are content-hashed by Next, so they get `public,
max-age=31536000, immutable`; everything else (the SPA-fallback HTML/other
unhashed files, and anything under `/_next/` outside `/_next/static/`)
gets `no-cache` so a rebuild is picked up on the next request instead of
being blanket `no-store`d or left to browser-heuristic caching.

The first version of this fix set that `Cache-Control` eagerly (a plain
`header` "set", no deletion), so it survived into a missing-asset's 404
response — the exact bug class F1 exists to kill, just with a cache header
instead of a content-type mismatch: a chunk that later reappears with the
same hash (rebuild without it, then rebuild with it again) would serve a
cached 404 for a year. `X-Content-Type-Options`/`X-Frame-Options`/
`Referrer-Policy`/the `-Server` deletion were also silently absent from
that same 404, for a related reason — they're deferred (any header
deletion forces Caddy to defer the whole block until `WriteHeader` time),
and `file_server`'s 404 unwinds as an error before that deferred write
ever happens, so `handle`-level headers never applied to it either way.
Fixed with a `handle_errors` block that deletes `Cache-Control` and
restates the security headers, since deletions are exactly what forces
deferred application — they run when the error response actually gets
written, not before.

### CI
`.github/workflows/ci.yml` gained a `caddy` job that runs `caddy validate`
against `deploy/Caddyfile` for both `DOMAIN=:80` and a real-looking domain,
using the plain `caddy:2` image with the file bind-mounted the same way
compose does it. Previously nothing in CI or the image build ever parsed
this file — a syntax error would only surface after a deploy.

- Status: `done` — verified against the running local stack: `curl` shows
  HSTS absent / security headers present / correct `Cache-Control` on the
  root document, on a missing `/_next/static/*` asset (404, no long-lived
  cache directive), and on a real one (200, `immutable`); the API still
  answers with a bearer token; headless Chromium lands on `/learn` with all
  requests 200 and no console output; `caddy validate` passes for both
  `DOMAIN=:80` and a real hostname; the new CI `caddy` job was proven to
  fail against a deliberately broken Caddyfile and pass again once
  restored; `docker compose -f deploy/docker-compose.prod.yml config` was
  proven to fail fast with `DOMAIN` unset and succeed with it set.
