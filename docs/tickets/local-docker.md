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
