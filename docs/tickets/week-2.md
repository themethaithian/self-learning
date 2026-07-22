# Week 2 — ขึ้น VPS + แกนของระบบอ่าน

**เป้าหมาย**: skeleton จากสัปดาห์ 1 รันจริงบน HTTPS ที่ DigitalOcean, merge → auto-deploy,
แล้ววางฐาน Read & Recall (schema + importer + LLM client)
**เตรียมก่อน**: GitHub repo (จากสัปดาห์ 1), บัญชี DigitalOcean, domain/DuckDNS, Anthropic API key (T11)

---

## T28 — Production Dockerfile + compose + Caddyfile `[go-implementer]` ~40 นาที
- **Scope**: `Dockerfile` (multi-stage: build Go + build Next export → runtime image เล็ก), `deploy/docker-compose.prod.yml` (mysql + api + caddy, MySQL ไม่ expose port), `deploy/Caddyfile` (domain → static files + reverse proxy `/api` → api:8080)
- **Acceptance**: รัน prod compose บนเครื่อง local แล้วใช้งานผ่าน Caddy ได้ครบ
- **Review focus**: image ไม่มี source/secret ค้าง, MySQL volume + healthcheck, Caddy security headers
- Status: `todo`

## T29 — VPS setup script + runbook `[go-implementer]` ~45 นาที
- **Scope**: `deploy/setup-vps.sh` (สร้าง user non-root, SSH hardening: key-only + no root login, UFW 22/80/443, fail2ban, ติดตั้ง Docker, วาง `.env` chmod 600), `docs/runbook.md` (ขั้นตอน manual: สร้าง droplet SGP $6, ชี้ DNS, รัน script, first deploy, วิธี restore)
- **Acceptance**: รัน script บน droplet ใหม่จริง 1 ครั้งแล้วได้เครื่องพร้อมใช้ (**ticket นี้คุณต้องลงมือรันเอง** — implementer เขียน script/runbook ให้)
- **Review focus**: script idempotent (รันซ้ำไม่พัง), ไม่มี secret hardcode ใน script
- Status: `todo`

## T30 — GitHub Actions deploy + smoke test `[go-implementer]` ~45 นาที
- **Scope**: `.github/workflows/deploy.yml` — **merge เข้า develop** → vet/test → build image → push GHCR → SSH เข้า VPS → `docker compose pull && up -d` → curl `/healthz` เช็คผล (นี่คือปลายทางของ PR workflow: คุณ merge จากมือถือ = deploy เอง)
- **Acceptance**: merge PR 1 ครั้งแล้ว deploy จริงสำเร็จ, smoke test fail = workflow แดง
- **Review focus**: secrets ผ่าน GitHub Secrets เท่านั้น, SSH key เป็น deploy key แยก (ไม่ใช่ key ส่วนตัว), ไม่ log ค่า env
- Status: `todo`

## T9 — Lesson domain + migration 002 `[go-implementer]` ~45 นาที
- **Goal**: Lesson/RecallCheck entities + กติกา gating ใน `LessonProgress`
- **Scope**: `internal/curriculum/domain` (Lesson, RecallCheck, References), `internal/learning/domain` (LessonProgress + ChunkState), `migrations/002_lessons_learning.sql` — ตาราง lessons มีคอลัมน์ `refs JSON` (แหล่งอ่านต่อของจริง, ชื่อคอลัมน์เลี่ยงคำสงวน REFERENCES)
- **Acceptance**: table-driven tests ของ gating: lesson แรกของบท = ปลดล็อกเสมอ, ถัดไปปลดเมื่อก่อนหน้า passed, ข้ามบทไม่ได้
- **Review focus**: กติกา gating อยู่ใน domain method ไม่ใช่ SQL/handler, state transition ผิด order ต้อง error
- Status: `todo`

## T10 — cmd/import-lessons `[go-implementer]` ~30 นาที
- **Goal**: importer ตาม schema JSON ของ lesson-writer (ดู `.claude/agents/lesson-writer.md`)
- **Scope**: `cmd/import-lessons` — อ่าน dir, validate JSON (รวม `references` 2–4 รายการ), upsert ตาม concept slug + version
- **Acceptance**: import ซ้ำ idempotent, JSON ผิด schema → error ชัดเจนพร้อมชื่อไฟล์, ไม่ import ทับ version เดิมโดยไม่ bump
- **Review focus**: validation ครบไหม (recall_checks 3–5 ข้อ, est_minutes 5–10, references มี title+source), error message บอกไฟล์+field
- Status: `todo`

## T11 — LLM port + Anthropic client (เขียนเอง) `[go-implementer]` ~45 นาที
- **Goal**: **ticket สำคัญของสัปดาห์** — Anthropic Messages API client ด้วย `net/http` ล้วน
- **Scope**: port `Grader` interface ฝั่ง `internal/learning/app`; adapter `internal/platform/llm` — `claude-haiku-4-5`, retry + exponential backoff เมื่อ 429/5xx, timeout ผ่าน context, max_tokens cap, JSON-mode prompt สำหรับตรวจ recall (grade 0–5 + feedback ไทย)
- **Acceptance**: unit test ด้วย `httptest.Server` mock (ตอบปกติ / 429 / timeout), ไม่มี dependency ใหม่ใน go.mod
- **Review focus**: API key ไม่หลุดใน log, `resp.Body.Close()` ครบทาง, retry ไม่ retry ตัว non-idempotent error (400), context cancellation ไหลถึง http.Client จริง
- Status: `todo`
