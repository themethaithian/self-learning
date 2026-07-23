# Week 1 — Walking skeleton (local)

**เป้าหมาย**: `docker compose up` แล้วเปิด http://localhost:3000 เห็น curriculum tree ครบ 5 track + repo ขึ้น GitHub พร้อม PR workflow
**กติกา**: ทุก ticket → branch `ticket/<id>-<slug>` → go-implementer ทำ → vet/test ผ่าน → code-reviewer → เปิด PR เข้า develop → คุณรีวิว+merge (มือถือได้) → ติ๊กใน roadmap.md

---

## T1 — Repo scaffold + docker-compose MySQL `[go-implementer]` ~30 นาที
- **Goal**: โครง repo ตาม design §4 + MySQL 8 รันใน Docker ได้
- **Scope**: โฟลเดอร์เปล่าตามโครง, `Makefile` (targets: `up`, `test`, `vet`, `run`), `docker-compose.yml` (mysql:8 + volume + healthcheck), `.gitignore`, `.env.example`, `.github/pull_request_template.md` (หัวข้อ: สรุปไทย / walkthrough รายไฟล์ / Review focus checklist / ผล vet+test / verdict ของ code-reviewer)
- **Acceptance**: `make up` แล้ว MySQL พร้อมใช้, `mysql -h127.0.0.1` ต่อได้, PR แรกของ repo ใช้ template นี้
- **Review focus**: healthcheck ถูกต้องไหม, `.env` ไม่หลุดเข้า git
- Status: `done`

## T2 — Config loader + migration runner + migration 001 `[go-implementer]` ~45 นาที
- **Goal**: อ่าน config จาก env + migration runner เขียนเองด้วย stdlib (learning goal — ไม่ใช้ golang-migrate)
- **Scope**: `internal/platform/config`, `internal/platform/mysql` (pool + runner อ่าน `migrations/*.sql` เรียง version, บันทึกลง `schema_migrations`), `migrations/001_curriculum.sql` (ตาราง topics/chapters/concepts/lessons/recall_checks)
- **Acceptance**: รัน API ครั้งแรกแล้ว migrate อัตโนมัติ, รันซ้ำ = no-op, มี table-driven test ของ runner (เรียงลำดับ, ข้าม version ที่ apply แล้ว)
- **Review focus**: transaction ตอน apply migration, SQL injection ใน runner (ห้าม interpolate ชื่อไฟล์ลง query), pool settings (`SetMaxOpenConns` ฯลฯ)
- Status: `done`

## T3 — HTTP server + healthz + logging/recovery middleware `[go-implementer]` ~45 นาที
- **Goal**: `net/http` server + middleware chain เขียนเอง (learning goal หลักของโปรเจกต์)
- **Scope**: `cmd/api`, `internal/platform/httpserver` (mux, graceful shutdown), `internal/platform/middleware` (logging: method/path/status/duration, recovery: panic → 500 + stack log), `GET /healthz` (เช็ค DB ping)
- **Acceptance**: `curl /healthz` → 200 + JSON, panic ใน handler ไม่ทำ server ตาย, log ออกครบ, มี test ของ middleware ทั้งสอง
- **Review focus**: middleware chaining pattern (func(http.Handler) http.Handler), `http.Server` timeouts (Read/Write/Idle) ต้องตั้ง, graceful shutdown ด้วย context
- Status: `done`

## T4 — Bearer auth middleware + CORS `[go-implementer]` ~30 นาที
- **Goal**: ทุก endpoint ยกเว้น `/healthz` ต้องมี `Authorization: Bearer <token>`
- **Scope**: `internal/platform/middleware/auth.go`, `cors.go`; token จาก config
- **Acceptance**: ไม่มี token → 401, token ผิด → 401, ถูก → ผ่าน; CORS ตอบ preflight ให้ origin ของ frontend; tests ครบ 3 กรณี
- **Review focus**: ใช้ `crypto/subtle.ConstantTimeCompare` เทียบ token, CORS ไม่เปิด `*` ตอน production
- Status: `done`

## T35 — CI: vet + test + build ทุก PR `[go-implementer]` ~25 นาที
- **Goal**: PR ที่ CI เขียว = ผล test เชื่อถือได้โดยไม่ต้องเชื่อ output ที่ paste มา — เป็นเงื่อนไขให้ Review level 🟢 merge จากสรุปได้
- **Scope**: `.github/workflows/ci.yml` — trigger `pull_request` + push `develop`: `go vet ./...`, `go test ./...`, `go build ./...` พร้อม Go module cache (เพิ่ม `npm run build` เมื่อ `web/` เกิดใน T8)
- **Acceptance**: PR ถัดไปโชว์ CI เขียว, ลองทำ test พังใน branch ทดสอบแล้ว CI แดงจริง
- **Review focus**: ไม่มี secret ใน workflow, permissions ของ GITHUB_TOKEN แคบสุด (contents: read)
- Status: `done`

## T5 — Curriculum domain + tests `[go-implementer]` ~45 นาที
- **Goal**: domain layer แรก — เป็นแม่แบบ DDD ให้ context อื่นทั้งหมด
- **Scope**: `internal/curriculum/domain` — Topic/Chapter/Concept entities, VOs (Slug, Track, Position) validate ใน constructor, ห้าม import อะไรนอก stdlib
- **Acceptance**: table-driven tests ครอบ constructor ทุกตัว (valid/invalid), `go vet` ผ่าน
- **Review focus**: **นี่คือ ticket ที่ควรรีวิวละเอียดสุดของสัปดาห์** — VO เป็น immutable ไหม, error เป็น sentinel/wrapped ถูกแบบไหม, ไม่มี DB/JSON tag ใน domain struct
- Status: `done`

## T6 — Curriculum repo + GET /curriculum `[go-implementer]` ~45 นาที
- **Goal**: เส้นแรกที่ต่อครบ 3 layers: handler → app service → repo
- **Scope**: `internal/curriculum/app` (service + repository interface), `internal/curriculum/infra` (MySQL repo + HTTP handler), route `GET /api/v1/curriculum` คืน tree ซ้อน 3 ชั้น
- **Acceptance**: curl ได้ tree JSON, repo มี test ต่อ MySQL จริง (ใน docker) หรือ interface test
- **Review focus**: repository interface ประกาศฝั่ง app (ไม่ใช่ infra), query N+1 (ควร join/รวม query), DTO แปลงที่ infra ไม่ใช่ domain
- Status: `done` — response เปลี่ยนเป็น group ตาม track (ดู design.md decision log 2026-07-23)

## T7 — JSON format + cmd/import-curriculum + ddd.json `[go-implementer]` ~45 นาที
- **Goal**: write path template + track แรกลง MySQL (แยกจาก T7b เพราะ 195 concepts เกิน 45 นาทีแน่)
- **Scope**: `content/curriculum/ddd.json` (1 topic / 5 chapters / 33 concepts), loader ใน infra, `Writer` port ใน app, `SaveTopic` (1 transaction ต่อ 1 topic) — idempotent upsert ตาม slug **และ** reconcile ลบ chapter/concept ที่หายไปจากไฟล์ (orphan) ภายใน transaction เดียวกัน, `cmd/import-curriculum` (T7b จะ copy format นี้ไปอีก 4 ไฟล์ ต้องรู้ไว้ว่า import ไม่ใช่ additive-only ล้วน ๆ)
- **Acceptance**: import 2 รอบ ได้ผลเท่าเดิม (no duplicate), มี test ที่โหลดไฟล์ content จริงเพื่อจับ typo
- **Review focus**: upsert logic (`INSERT … ON DUPLICATE KEY UPDATE` หรือ select-then-update ใน tx), slug ใน JSON ตรงกับ design
- **ต้อง import ใน transaction เสมอ** (จาก deep review ของ T5): `NewTopic` เป็นทางเดียวที่สร้าง Topic ได้ และมันบังคับ `ErrNoChildren` ดังนั้น topic ที่ถูกเขียนลง DB ค้างไว้แบบยังไม่มี chapter จะทำให้ `GET /curriculum` (T6 reconstitute ผ่าน constructor) ล้มทั้งเส้น ไม่ใช่แค่ topic เดียว
- Status: `done`

## T7b — Curriculum JSON อีก 4 track `[go-implementer × 4 ขนาน]` ~40 นาที
- **Goal**: distsys (43) + aws (37) + go (33) + dsa (49) = 162 concepts ที่เหลือ
- **Scope**: `content/curriculum/{distsys,aws,go,dsa}.json` ตาม format ของ T7 — slug ลอกจาก design §7 ตรงตัวอักษร (เป็น contract กับ `content/lessons/<topic>/<concept>.json`), title อังกฤษ, outline ไทย 2–4 bullet
- **Acceptance**: loader test ของแต่ละไฟล์ผ่าน + จำนวน chapter/concept ตรงกับ design §7, `GET /curriculum` เห็นครบ 195 concepts
- **Review focus**: slug ตรง design เป๊ะไหม (typo = lesson file ไปคนละที่), outline สั้นพอที่จะเป็น guidance ไม่ใช่บทเรียนย่อ
- Status: `done` — 4 ไฟล์ยิงขนาน 4 agent (คนละไฟล์ ไม่ชนกัน), reviewer เขียนสคริปต์ diff slug กับ design §7 เชิงกลไก

## T8 — Next.js scaffold + curriculum tree page `[go-implementer]` ~45 นาที
- **Goal**: frontend แรก — เห็น tree จริงจาก API
- **Scope**: `web/` (Next.js, `output: 'export'`, Tailwind), `lib/api.ts` (fetch wrapper + bearer จาก localStorage), หน้า token entry, หน้า curriculum tree (expand/collapse)
- **Acceptance**: `npm run build` ผ่าน (static export), เปิดหน้าเว็บ ใส่ token เห็น tree ครบ
- **Review focus**: **ใช้ skill `frontend-design` เป็นครั้งแรก — เช็คว่า layout/สี/ฟอนต์ตาม skill**, token ไม่โผล่ใน URL, มี loading/error state
- Status: `done` — walking skeleton ครบวง: เปิดเว็บ ใส่ token เห็น tree 195 concepts จาก MySQL จริง (screenshot ยืนยันด้วย puppeteer headless); backend T6/T7/T7b พิสูจน์กับ MySQL 8.4 จริงแล้วในตัว
