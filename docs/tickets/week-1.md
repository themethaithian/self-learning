# Week 1 — Walking skeleton (local)

**เป้าหมาย**: `docker compose up` แล้วเปิด http://localhost:3000 เห็น curriculum tree ครบ 5 track + repo ขึ้น GitHub พร้อม PR workflow
**กติกา**: ทุก ticket → branch `ticket/<id>-<slug>` → go-implementer ทำ → vet/test ผ่าน → code-reviewer → เปิด PR เข้า develop → คุณรีวิว+merge (มือถือได้) → ติ๊กใน roadmap.md

---

## T1 — Repo scaffold + docker-compose MySQL `[go-implementer]` ~30 นาที
- **Goal**: โครง repo ตาม design §4 + MySQL 8 รันใน Docker ได้
- **Scope**: โฟลเดอร์เปล่าตามโครง, `Makefile` (targets: `up`, `test`, `vet`, `run`), `docker-compose.yml` (mysql:8 + volume + healthcheck), `.gitignore`, `.env.example`, `.github/pull_request_template.md` (หัวข้อ: สรุปไทย / walkthrough รายไฟล์ / Review focus checklist / ผล vet+test / verdict ของ code-reviewer)
- **Acceptance**: `make up` แล้ว MySQL พร้อมใช้, `mysql -h127.0.0.1` ต่อได้, PR แรกของ repo ใช้ template นี้
- **Review focus**: healthcheck ถูกต้องไหม, `.env` ไม่หลุดเข้า git
- Status: `todo`

## T2 — Config loader + migration runner + migration 001 `[go-implementer]` ~45 นาที
- **Goal**: อ่าน config จาก env + migration runner เขียนเองด้วย stdlib (learning goal — ไม่ใช้ golang-migrate)
- **Scope**: `internal/platform/config`, `internal/platform/mysql` (pool + runner อ่าน `migrations/*.sql` เรียง version, บันทึกลง `schema_migrations`), `migrations/001_curriculum.sql` (ตาราง topics/chapters/concepts/lessons/recall_checks)
- **Acceptance**: รัน API ครั้งแรกแล้ว migrate อัตโนมัติ, รันซ้ำ = no-op, มี table-driven test ของ runner (เรียงลำดับ, ข้าม version ที่ apply แล้ว)
- **Review focus**: transaction ตอน apply migration, SQL injection ใน runner (ห้าม interpolate ชื่อไฟล์ลง query), pool settings (`SetMaxOpenConns` ฯลฯ)
- Status: `todo`

## T3 — HTTP server + healthz + logging/recovery middleware `[go-implementer]` ~45 นาที
- **Goal**: `net/http` server + middleware chain เขียนเอง (learning goal หลักของโปรเจกต์)
- **Scope**: `cmd/api`, `internal/platform/httpserver` (mux, graceful shutdown), `internal/platform/middleware` (logging: method/path/status/duration, recovery: panic → 500 + stack log), `GET /healthz` (เช็ค DB ping)
- **Acceptance**: `curl /healthz` → 200 + JSON, panic ใน handler ไม่ทำ server ตาย, log ออกครบ, มี test ของ middleware ทั้งสอง
- **Review focus**: middleware chaining pattern (func(http.Handler) http.Handler), `http.Server` timeouts (Read/Write/Idle) ต้องตั้ง, graceful shutdown ด้วย context
- Status: `todo`

## T4 — Bearer auth middleware + CORS `[go-implementer]` ~30 นาที
- **Goal**: ทุก endpoint ยกเว้น `/healthz` ต้องมี `Authorization: Bearer <token>`
- **Scope**: `internal/platform/middleware/auth.go`, `cors.go`; token จาก config
- **Acceptance**: ไม่มี token → 401, token ผิด → 401, ถูก → ผ่าน; CORS ตอบ preflight ให้ origin ของ frontend; tests ครบ 3 กรณี
- **Review focus**: ใช้ `crypto/subtle.ConstantTimeCompare` เทียบ token, CORS ไม่เปิด `*` ตอน production
- Status: `todo`

## T5 — Curriculum domain + tests `[go-implementer]` ~45 นาที
- **Goal**: domain layer แรก — เป็นแม่แบบ DDD ให้ context อื่นทั้งหมด
- **Scope**: `internal/curriculum/domain` — Topic/Chapter/Concept entities, VOs (Slug, Track, Position) validate ใน constructor, ห้าม import อะไรนอก stdlib
- **Acceptance**: table-driven tests ครอบ constructor ทุกตัว (valid/invalid), `go vet` ผ่าน
- **Review focus**: **นี่คือ ticket ที่ควรรีวิวละเอียดสุดของสัปดาห์** — VO เป็น immutable ไหม, error เป็น sentinel/wrapped ถูกแบบไหม, ไม่มี DB/JSON tag ใน domain struct
- Status: `todo`

## T6 — Curriculum repo + GET /curriculum `[go-implementer]` ~45 นาที
- **Goal**: เส้นแรกที่ต่อครบ 3 layers: handler → app service → repo
- **Scope**: `internal/curriculum/app` (service + repository interface), `internal/curriculum/infra` (MySQL repo + HTTP handler), route `GET /api/v1/curriculum` คืน tree ซ้อน 3 ชั้น
- **Acceptance**: curl ได้ tree JSON, repo มี test ต่อ MySQL จริง (ใน docker) หรือ interface test
- **Review focus**: repository interface ประกาศฝั่ง app (ไม่ใช่ infra), query N+1 (ควร join/รวม query), DTO แปลงที่ infra ไม่ใช่ domain
- Status: `todo`

## T7 — Curriculum JSONs + cmd/import-curriculum `[go-implementer]` ~45 นาที
- **Goal**: tree ทั้ง 5 track (จาก design §7) ลง MySQL
- **Scope**: `content/curriculum/{ddd,distsys,aws,go,dsa}.json` (แปลงจาก design §7 ตรง ๆ), `cmd/import-curriculum` (idempotent upsert ตาม slug)
- **Acceptance**: import 2 รอบ ได้ผลเท่าเดิม (no duplicate), `GET /curriculum` เห็นครบ ~175 concepts
- **Review focus**: upsert logic (`INSERT … ON DUPLICATE KEY UPDATE` หรือ select-then-update ใน tx), slug ใน JSON ตรงกับ design
- Status: `todo`

## T8 — Next.js scaffold + curriculum tree page `[go-implementer]` ~45 นาที
- **Goal**: frontend แรก — เห็น tree จริงจาก API
- **Scope**: `web/` (Next.js, `output: 'export'`, Tailwind), `lib/api.ts` (fetch wrapper + bearer จาก localStorage), หน้า token entry, หน้า curriculum tree (expand/collapse)
- **Acceptance**: `npm run build` ผ่าน (static export), เปิดหน้าเว็บ ใส่ token เห็น tree ครบ
- **Review focus**: **ใช้ skill `frontend-design` เป็นครั้งแรก — เช็คว่า layout/สี/ฟอนต์ตาม skill**, token ไม่โผล่ใน URL, มี loading/error state
- Status: `todo`
