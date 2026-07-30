# Roadmap — Self-Improve Web

> 📍 **ย้ายเครื่อง / กลับมาทำงานต่อ?** อ่าน [`docs/HANDOFF.md`](HANDOFF.md) ก่อน —
> priority ปัจจุบันคือสอบ AWS SAA-C03 งานอื่นด้านล่างพักไว้ชั่วคราว
>
> **ภาพใหญ่ทั้งหมดอยู่ไฟล์นี้ไฟล์เดียว** — รายละเอียด ticket อยู่ใน [`docs/tickets/`](tickets/)
> ส่วน design ฉบับเต็ม (DDD, schema, API, VPS) อยู่ที่ [`docs/design.md`](design.md)
> ⚠️ `design.md` **เก่ากว่าไฟล์นี้** — ดูหัวข้อ "หนี้ที่รู้ตัว" ท้ายไฟล์
>
> อัปเดตล่าสุด **2026-07-30** (AWS-0 — AWS certification กลายเป็น priority อันดับ 1)

## 🔥 Priority ปัจจุบัน: AWS Certified Solutions Architect – Associate (SAA-C03)

สอบจริงใน **1–2 เดือน**, อ่านวันละ 1–2 ชม., ยังไม่เคยจับ AWS จริง (รู้แค่ชื่อ service) —
เป้าหมายนี้แทรกหน้าคิวชั่วคราว รายละเอียดเต็มอยู่ที่ [`docs/tickets/aws-cert.md`](tickets/aws-cert.md)
(blueprint ข้อสอบ, concept → task statement mapping, แผนคลังข้อสอบ ~550 ข้อ, กติกาการเขียน
คำถาม+เฉลย) — **Q-2d, Q-3, SIM-\* พักไว้ระหว่างนี้** (ดูหัวข้อ "พักไว้โดยตั้งใจ" ด้านล่าง)

## เป้าหมายตอนนี้

เว็บ "หนังสือเรียนส่วนตัว" ที่**เรียนแล้ววัดผลได้** เพื่อเตรียมสัมภาษณ์ตำแหน่ง
Backend / AI-CRM Platform — โหมด **learn-first**: ให้ความสำคัญกับ *เนื้อหา + วงจร
อ่าน→ทดสอบ→วัดผล* ก่อน feature และก่อน deploy

ทักษะเป้าหมาย 7 track: DDD · **DDIA (Distributed Data)** · **AI & LLM Systems** ·
Distributed Systems · **AWS SAA-C03 (priority ปัจจุบัน)** · Go · DSA
(สองตัวหนา = เขียนบทเรียนครบแล้ว, ที่เหลือมี tree แต่ยังไม่มีบทเรียน)

## สถาปัตยกรรมสรุป

- Go 1.26 stdlib `net/http` เท่านั้น + MySQL 8 (Docker) + Next.js static export
- DDD modular monolith — bounded context ที่ **มีจริงแล้ว**: `curriculum` · `learning` · `prefs`
  (`internal/platform` = infrastructure ไม่ใช่ context) · design.md ยังพูดถึง
  practice / review / assessment / stats ซึ่ง**ยังไม่ได้สร้าง**
- เนื้อหา: generate offline (lesson-writer → lesson-verifier) มี mermaid + references → import เข้า MySQL
- **ยังไม่มีการให้คะแนนใด ๆ ในระบบ** — ไม่มีตาราง attempts/cards/logs, ไม่มี LLM call path
  (`internal/platform/llm` ยังไม่มีจริง) · แผน self-graded (`graded_by='self'`) อยู่ใน phase 4
- รันเอง: `make dev` — **ยังไม่ deploy ขึ้น VPS**

## เริ่มงานใน session ใหม่ต้องรู้

- `make dev` = `docker compose up -d --build` → mysql + api + web + seed (seed รันอัตโนมัติ)
  · `make dev-reset` ล้าง volume · `make seed` import ใหม่หลังแก้ JSON ใน `content/`
- **ทุก endpoint อยู่หลัง bearer token** — เว็บจะว่างเปล่าถ้าไม่ใส่ token (`API_BEARER_TOKEN`)
- env ที่ **บังคับ** (`internal/platform/config`): `DB_HOST` `DB_PORT` `DB_NAME` `DB_USER`
  `DB_PASSWORD` `API_BEARER_TOKEN` `CORS_ALLOWED_ORIGIN` (ห้ามเป็น `*` และห้ามมี path ต่อท้าย)
  และ **`ANTHROPIC_API_KEY`** — บังคับให้มีค่า แต่ **ไม่ต้องเป็น key จริง** เพราะยังไม่มีโค้ดที่เรียก
  Anthropic เลย · `docker-compose.yml` ใส่ placeholder ให้แล้ว ส่วน `make run` (รัน Go ตรง ๆ)
  จะ error ถ้าไม่ตั้งเอง

## แผนตาม phase

> **หมายเหตุสำคัญ**: แผน "6 สัปดาห์ T1–T35" ฉบับเดิม**เลิกใช้แล้ว** ตั้งแต่เป้าหมายแคบลง
> เป็นตำแหน่ง Backend/AI-CRM — เลข T เดิมยังอ้างอิงได้ใน `tickets/week-*.md` แต่ลำดับ
> การทำงานจริงคือ phase ข้างล่างนี้
>
> phase **ไม่ใช่ช่วงเลข PR ที่ต่อเนื่องกัน** — งาน content, quality และ UX คาบเกี่ยวกันจริง
> (เช่น 2026-07-28 วันเดียวมีทั้ง #40 UX-1 และ #43 MCQ) ให้ดู ticket index เป็นหลัก

| Phase | เป้าหมาย | สถานะ |
|---|---|---|
| **AWS-cert** | รีบาลานซ์ curriculum tree ให้ตรง exam guide + คลังข้อสอบ ~550 ข้อ ก่อนสอบ SAA-C03 | 🔥 **priority ปัจจุบัน** ([aws-cert.md](tickets/aws-cert.md)) |
| 0. Walking skeleton + reading slice | คลิก concept → อ่านบทเรียนจริง (mermaid + references + recall) | ✅ จบ |
| 1. Content | เขียนบทเรียน 2 track ที่ตรงกับตำแหน่งที่สุด | ✅ จบ |
| 2. Quality + รันเองได้ | MCQ ที่เดาไม่ได้ + `make dev` คำสั่งเดียว | ✅ จบ |
| 3. Guided learning path | รู้ว่า "วันนี้อ่านอะไรต่อ" และอ่านถึงไหนแล้ว | ✅ จบ (UX-8 เป็น optional, ทำเมื่อรู้สึกขาดจริง) |
| 4. Measurable recall | quiz กดเลือกได้ + เก็บประวัติ + SRS ที่วัดผลได้ | 🟡 บางส่วนเดินต่ออยู่: Q-2b (PR pending) + Q-2c (branch `ticket/q-2c-sm2-domain` มีอยู่แล้วบน origin, ไม่ได้พัก) — เฉพาะ Q-2d/Q-3 พักไว้ระหว่าง AWS-cert |
| 5. Visual simulation | บทเรียนที่เห็นภาพและโต้ตอบได้ | ⏸ พักไว้ระหว่าง AWS-cert (SIM-0…SIM-3) |
| — | Deploy ขึ้น VPS + track ที่เหลือ | ⏸ **พักไว้โดยตั้งใจ** |

## Ticket index (ติ๊กเมื่อ merge แล้ว)

### AWS-cert — priority ปัจจุบัน 🔥 ([รายละเอียด](tickets/aws-cert.md))
- [ ] AWS-0 (in review) — รีบาลานซ์ curriculum tree จาก 37 → **50 concept** ให้ตรง 30/26/24/20
  ตาม domain weight จริงของ SAA-C03 exam guide (เดิม 10/10/12/5 เอียงไปทาง Domain 3) +
  pin blueprint ข้อสอบ, concept → task statement mapping, แผนคลังข้อสอบ ~550 ข้อ ไว้ที่
  [`docs/tickets/aws-cert.md`](tickets/aws-cert.md)
- [ ] AWS-S1 (in review, round 2) — `maxRecallChecks` 5→15 + `explanation` end-to-end
  (recall_checks schema, ไม่แตะ multiple-response) — branch `ticket/aws-1-explanation`
  (ชื่อ branch เกิดก่อน ticket ID นี้ถูกตั้งใหม่เป็น AWS-S1/AWS-S2/AWS-S3/AWS-C1..C4 — ไม่ force-push
  เปลี่ยนชื่อ branch, ดูรายละเอียดที่ aws-cert.md)
- [ ] AWS-S2 (ยังไม่เริ่ม) — multiple-response support (second correct answer, "Select TWO",
  all-or-nothing scoring)
- [ ] AWS-S3 (ยังไม่เริ่ม, spec ปักไว้แล้วใน aws-cert.md) — guessability measurement tool
  ที่ parameterize ได้ (`1/len(options)` ต่อ corpus) — **gate คลังข้อสอบจริง (AWS-C1..C4) ต้องรอ
  ticket นี้ก่อน**
- [ ] AWS-C1..C4 (ยังไม่ตัดชื่อ, gated บน AWS-S1/S2/S3) — เขียนคลังข้อสอบทีละ domain ตามแผนใน
  aws-cert.md

### Phase 0 — walking skeleton + reading slice ✅
- [x] T1 [x] T2 [x] T3 [x] T4 [x] T35 [x] T5 [x] T6 [x] T7 [x] T7b [x] T8 ([week-1](tickets/week-1.md))
- [x] T9 [x] T10 [x] T-read — คลิก concept → อ่านบทเรียนจริง ([week-2](tickets/week-2.md))
- [x] C-content — DDD บท 1 (4 บทเรียน) + trade-off sections
- [x] design system เป็น light modern-minimal (#9) · [x] T-reading-comfort (#21) — mermaid อ่านออกบนมือถือ + พื้นหลัง cream

### Phase 1 — content ✅
- **DDIA** ([รายละเอียด](tickets/ddia.md)): [x] T-ddia-track [x] T-ddia-curriculum ·
  **61 บทเรียน / 12 บท ครบทั้ง track** — [Study Reader artifact](https://claude.ai/code/artifact/b67a1cad-4b9e-4b2e-9cf1-755128e1ebb5)
- **AI & LLM Systems** ([รายละเอียด](tickets/ai-systems.md)): [x] T-ai-track (tree 11 บท / 53 concept) ·
  **53 บทเรียน / 11 บท ครบทั้ง track** — เจาะตำแหน่ง Backend/AI-CRM ([🎯 roadmap](https://claude.ai/code/artifact/a6d332fd-de7b-44bb-91db-37549984c7a7))
- **DDD**: 4 บทเรียน (บท 1 จาก 5) — ยังไม่ครบ track
- `distsys` / `aws` / `go` / `dsa`: มี curriculum tree แต่ **0 บทเรียน** (aws รีบาลานซ์เป็น
  50 concept ใน AWS-0 แล้ว แต่ยังไม่มีบทเรียน — คลังข้อสอบคือ deliverable ถัดไป ไม่ใช่บทเรียนอ่าน)
- รวม concept ทั้ง 7 tree = **322** (หลัง AWS-0) · recall check ทั้งหมด **586** ข้อ (mcq 356 / short_answer 230)

### Phase 2 — quality + รันเองได้ ✅
- [x] T-local-docker (#39, [รายละเอียด](tickets/local-docker.md)) — `make dev` = mysql + api + web + seed
- [x] C-mcq-sweep (#38) · [x] C-mcq-balance (#43) ([รายละเอียด](tickets/mcq-quality.md)) —
  **MCQ 356 ข้อ** เดาด้วย heuristic ความยาวได้ **33.7%** (เดิม 89.6%) เดาด้วยตำแหน่ง 33.4% (เดิม 41.6%)
  \+ กติกาการแก้ distractor ที่ grep ตรวจได้ — **หมายเหตุ (ยืนยันแล้วตอน AWS-S1)**: ตัวเลขเหล่านี้
  วัดด้วย script ที่รันแบบ ad hoc นอก version control **ไม่มี script นี้ commit ไว้ในโปรเจกต์เลย**
  (เช็คแล้วทั้ง working tree และ `git log --all --diff-filter=A`) ตัวเลขจึง **รันซ้ำไม่ได้ตอนนี้** —
  ต้องรอ AWS-S3 (สร้าง measurement tool ใหม่) ก่อนถึงจะ verify ซ้ำหรือรันกับ corpus อื่นได้

### Phase 3 — guided learning path ✅ ([รายละเอียด](tickets/ux-today.md))
- [x] UX-1 (#40) `has_lesson`/`est_minutes` บน curriculum API
- [x] UX-1b (#41) focus track preference API + bounded context `prefs`
- [x] UX-2 (#42) หน้า `/learn` — track card + focus picker
- [x] UX-3 (#44) track view — concept ที่ยังไม่มีบทเรียนกดไม่ได้ + `n/N ready`
- [x] UX-4 (#45) 🔴 progress API — `GET /api/v1/progress` + `PUT /api/v1/progress/{topic}/{concept}`
  + bounded context `learning`
- [x] UX-5 (#46) reader loop — breadcrumb + Finish + Next (+ vitest ตัวแรกของ `web/`)
- [x] UX-6 (#49) หน้า `/today` + IA switch (nav เหลือ Today · Learn):
  "วันนี้อ่านอะไรต่อ" ไล่จาก focus track ก่อน แล้ว fallback ตาม
  `TRACK_DISPLAY_ORDER` · `domain.Gate` ยังไม่ถูก wire (soft-guide เท่านั้น)
- [x] UX-7 (#50) — reading-status indicator บน `/learn` (list +
  track detail) เท่านั้น ไม่ใช่หน้า Progress แยก: `ProgressBar` + `n/N lessons
  read` ต่อ track card (ตัวส่วน = lesson ที่มีจริง ไม่ใช่ concept ที่วางแผนไว้),
  read-state marker (`passed`/`in_progress`, shape+aria-label ไม่ใช่สีอย่างเดียว)
  ต่อ concept, read count ต่อ chapter header · ไม่แตะ `/lesson` reader
- [ ] UX-8 — หน้า `/progress` รวม (ทำเมื่อใช้ UX-7 แล้วยังรู้สึกขาดภาพรวมเท่านั้น)

### Phase 4 — measurable recall ⏭
- [ ] Q-1 (PR pending) — quiz **recall-first 3 stage**: Recall (เห็นคำถามอย่างเดียว, mcq
  ไม่โชว์ตัวเลือก) → Commit (mcq: เลือกตัวเลือก + ระดับความมั่นใจ Guessed/Unsure/Confident
  ก่อนเห็นเฉลยเสมอ; short_answer: เลือกแค่ระดับความมั่นใจ) → Reveal (mcq: mark ถูก/ผิดจากการ
  เทียบ `expected_answer` ไม่ใช่ self-report; **short_answer ยังเป็น self-report Pass/Not
  yet เหมือนเดิม** เพราะ free-text เทียบเองไม่ได้) · shuffle ตัวเลือกด้วย Fisher-Yates ตอน
  render (เสถียรตลอดอายุการ์ดผ่าน lazy `useState`, `rng` inject ได้เพื่อเทสต์) ·
  รายละเอียดเต็มใน [`docs/tickets/quiz.md`](tickets/quiz.md)
  · **ตัดสินใจแล้ว ไม่ต้อง revisit**: `expected_answer` **ยังอยู่**ใน payload ต่อไปใน v1 — v1
  self-graded (คนตอบ = คนให้คะแนนเอง) ทำ endpoint เฉลยแยกไม่ได้อะไรเพิ่ม มีแต่เสีย round trip +
  failure state ใหม่ ค่อยย้าย server-side ตอน Q-2 ที่เริ่ม submit attempt จริงและออกแบบ endpoint
  จาก requirement จริง (design.md §API ที่บอกว่าต้องตัดเป็นข้อความล้าสมัย ไม่ใช่ bug ของโค้ด)
- [x] Q-2a (**PR #52**, merged) — persist recall attempts: schema `recall_attempts`
  (append-only, `check_key = SHA256(topic/concept/trimmed-question)` ไม่มี FK ไป
  `recall_checks.id`), domain VOs (`CheckKey`/`Confidence`/`AttemptOutcome`/`GradedBy`/
  `RecallAttempt`), `POST /api/v1/progress/{topic}/{concept}/attempts` — `graded_by='self'`
  เกิดจริงครั้งแรกที่นี่ · รายละเอียดเต็ม + เหตุผลการตัดสินใจ + mutation table อยู่ที่
  [`docs/tickets/quiz.md`](tickets/quiz.md)
- [x] Q-2b (**PR #53**, merged) — wire the quiz to
  the attempts API: **ticket ที่ทำให้ข้อมูล in-memory ของ Q-1 (confidence ที่เลือก, ตัวเลือกที่กด,
  ผลถูก/ผิด) persist จริงในที่สุด แทนที่จะหายตอน refresh** — `web/lib/api.ts`'s `postAttempt`
  เรียก endpoint ของ Q-2a จาก `RecallCheckCard` (ยิงเมื่อ attempt ครบจริง: mcq ที่ stage 3 ทันที;
  short_answer debounce 1s ก่อน submit — กัน flip-flop ของ Pass/Not yet เขียนหลายแถวโดยไม่ตั้งใจ,
  flush ทันทีตอน Finish/navigate away/Retry/`pagehide`+`visibilitychange` ด้วย `keepalive:true`)
  ผ่าน `lesson/page.tsx` (identity guard คนละกลไกกับ `finish()`'s slug-based `stillCurrent()` เดิม
  — ใช้ `loadGenerationRef` นับ "visit" แทน slug equality, จำเป็นเพราะการกลับมาที่ lesson เดิมมี
  slug ซ้ำกับ visit ก่อนหน้า; per-check "Not saved" + Retry indicator, aggregate banner ทั้ง
  "saving" และ "error" เหนือปุ่ม Finish) · รายละเอียดเต็ม + เหตุผลการตัดสินใจ + mutation table +
  Playwright/SQL evidence อยู่ที่ [`docs/tickets/quiz.md`](tickets/quiz.md)
  · **หนี้ที่รู้ตัวแล้ว (follow-up เล็ก ๆ ทำเมื่อสะดวก, ไม่ใช่ ticket แยก)**: `lesson/page.tsx`
  มี identity mechanism สองแบบข้าง ๆ กัน — `submitAttempt` ใช้ `loadGenerationRef` (นับ visit)
  แต่ `finish()` ยังใช้ `currentIdentityRef` (slug equality) เดิมจาก UX-5 — ย้าย `finish()` มาใช้
  `loadGenerationRef` เหมือนกัน (แค่เปลี่ยน `stillCurrent()`'s เงื่อนไข 2 บรรทัด) แล้วลบ
  `currentIdentityRef`/slug guard ทิ้งไปเลย ไม่ใช่เพราะ `finish()` มีบั๊กจริงตอนนี้ (blast radius
  ของมันเล็กกว่า attempt case มาก — ดูเหตุผลเต็มที่ quiz.md's หัวข้อ "finish()") แต่เพราะมีสอง
  identity mechanism ซ้อนกันอยู่ในไฟล์เดียวเป็นกับดักสำหรับคนอ่านโค้ดครั้งถัดไป
- [ ] Q-2c (implemented, PR pending) — SM-2 scheduling: schema + pure domain. Migration
  `migrations/007_review.sql` (`review_cards` one row per `check_key`, `review_logs`
  append-only audit trail FK'd to `recall_attempts.id`) + pure domain (`ReviewQuality`,
  `EaseFactor`, `ReviewCard.Advance`) — **ไม่มี repository/write path เลย ตั้งใจแยกจาก Q-2d**
  · รายละเอียดเต็ม + เหตุผลการตัดสินใจ + mutation table + migration evidence อยู่ที่
  [`docs/tickets/quiz.md`](tickets/quiz.md)
- [ ] Q-2d — wire SM-2 into the attempt path: transactional `review_cards` update + `review_logs`
  write ต่อ `recall_attempts` row ที่บันทึกจริง (ตามสัญญา "ทุก attempt = 1 advance" ที่ Q-2c
  ตัดสินใจไว้), due-cards read endpoint · **ข้อกำหนดจาก Q-2b's code review**: query ที่อ่าน
  "แถวล่าสุดของ check_key นี้" ต้อง `ORDER BY created_at DESC, id DESC` ไม่ใช่แค่
  `created_at DESC` เฉย ๆ — `recall_attempts.created_at` เป็น `TIMESTAMP` (second precision)
  พิสูจน์แล้วว่าสองแถวที่ submit ห่างกันจริงในเวลาปกติ (ไม่ใช่ race condition) ตกอยู่วินาทีเดียวกัน
  ได้จริง ทำให้ `created_at DESC` เดี่ยว ๆ เรียงลำดับ "ล่าสุด" ผิดได้ (Pass→Not yet ในวินาทีเดียวกัน
  อาจอ่านกลับมาเป็น Pass) — `id` เป็น `AUTO_INCREMENT` (monotonic เสมอ) จึงต้องเป็น tie-breaker
- [ ] Q-3 — เพิ่ม `explanation` ให้ MCQ ครบทั้ง 356 ข้อ (ตอนนี้ 0 ข้อมี) — รอบเดียวกับที่ขัดเกลา
  distractor ที่หลุดธีม **3 concept** (`b-trees`, `column-oriented-storage`, `process-pauses`)

### Phase 5 — visual simulation ⏭
- [ ] SIM-0 — predict-before-reveal (0 KB, ไม่ใช้ library)
  · Kaushal & Panda (2019) วัด **การสอนด้วย animation เทียบกับไม่ใช้ animation**: กลุ่มพื้นฐานน้อย
  ได้ effect size **−0.16** (กลุ่มพื้นฐานสูง +0.49) — งานวิจัยนี้**ไม่ได้**เทียบ passive กับ interactive
  · "ต้องให้ทำนายก่อนเฉลย" เป็น**ข้อสรุปของเราเอง** จากงานวิจัยชุดอื่น ไม่ใช่ผลของ paper นี้
- [ ] SIM-1 latency-percentiles (slider) · [ ] SIM-2 quorum W+R>N (slider) ·
  [ ] SIM-3 2PC (step + ปุ่ม crash coordinator) → **หยุดประเมินก่อนทำเพิ่ม**

### พักไว้โดยตั้งใจ (ไม่ได้ยกเลิก)
- **Q-2d, Q-3, SIM-0…SIM-3**: พักตั้งแต่ AWS certification กลายเป็น priority อันดับ 1
  (สอบใน 1–2 เดือน) — กลับมาทำต่อหลังคลังข้อสอบ AWS ในแผน [`aws-cert.md`](tickets/aws-cert.md) เสร็จ
- **Deploy**: [x] T28 (prod Dockerfile + compose + Caddyfile, #25) · [ ] T29 [ ] T30
  (GitHub Actions → VPS) — พักตั้งแต่เปลี่ยนเป็นโหมด learn-first
- **Feature จากแผนเดิม** — บางส่วนจะถูกแทนที่ด้วย phase 4/5 ที่ออกแบบใหม่แล้ว:
  - [ ] T11 (LLM client) [ ] T12–T15 [ ] T33 (backup) [ ] C1 ([week-3](tickets/week-3.md))
  - [ ] T16–T21 SRS drill + dashboard [ ] C2 ([week-4](tickets/week-4.md))
  - [ ] T22–T25 DSA loop [ ] C3 [ ] C4 ([week-5](tickets/week-5.md))
  - [ ] T26 T27 T31 T32 [ ] T34 (debt backlog จาก T8/T9/T10/T-read) [ ] C5 [ ] C6 ([week-6](tickets/week-6.md))
- **Track ที่เหลือ**: DDD บท 2–5, distsys, aws, go, dsa · และ track ที่ยังไม่มี:
  event-driven/Kafka, gRPC, observability, Kubernetes

## หนี้ที่รู้ตัว (ยังไม่แก้ ตั้งใจปล่อย)

- **`docs/design.md` ล้าสมัยหลายจุด** — §7 ยังเขียนว่า "195 concepts" (จริง 309), layout ยังเป็น
  5 track (ก่อนมี ddia/ai-systems), 6 bounded context ที่ไม่มี `prefs`, §API ยังบอกว่า
  `GET /lessons/{id}` ต้องตัด expected answer ออก (Q-1 ตัดสินใจแล้วว่า v1 self-graded ไม่ตัด
  เหตุผลเต็มอยู่ที่บรรทัด Q-1 ด้านบน) และ **§2 MySQL Schema's `recall_checks` block (บรรทัด ~77)
  ยังไม่มีคอลัมน์ `explanation`** ที่เพิ่มใน AWS-S1 (migration 008) — ทุกจุดนี้ **ล้าสมัย ไม่ใช่
  code ผิด**, ตั้งใจไม่ refresh design.md ทั้งไฟล์ตอนนี้ (ทำทีเดียวตอนว่างจริง ๆ ตามที่บันทึกไว้)
- **`cors_test.go` ยังเทียบกับ constant ตัวเองบางส่วน** — `wantMethods` ถูกแก้เป็น literal แล้วใน #41
  แต่ `wantHeaders`/`wantMaxAge` ยังอ้าง `corsAllowedHeaders`/`corsMaxAge` = mutation ไม่มีทางจับได้
- rating รายข้อในหน้า reader ไม่ถูก persist (ปลดใน Q-2b — endpoint พร้อมแล้วจาก Q-2a) ·
  `Button` variant `"danger"` ไม่มีใครใช้
  และจะตก AA (~4.39:1) วันที่มีคนใช้
- **UX-7 (`/learn`): `TrackCard`'s `<h3>` sits under the page `<h1>` with no `<h2>` between**
  — axe's `heading-order` (best-practice, ไม่ใช่ WCAG) เจอทั้งบน `/learn` และหน้าอื่นที่มีอยู่แล้วบน
  `develop` ก่อน ticket นี้ ไม่ได้แก้ในรอบนี้เพราะเป็น pattern ที่ใช้ทั้งแอป ต้องแก้พร้อมกันทีเดียว
  ไม่ใช่แก้เฉพาะการ์ดเดียว

## กติกาการทำงาน (ทุก ticket)

1. ticket ละ ≤ 45 นาที (1 Build session) — ใหญ่กว่านั้นต้องแตกก่อนเริ่ม
2. ลำดับเสมอ: **go-implementer** ทำบน branch `ticket/<id>-<slug>` →
   `go vet` + `go test ./...` (+ `npm test` ถ้าแตะ `web/`) ผ่าน → **code-reviewer** รีวิว →
   แก้จนไม่มี required → เปิด **PR เข้า develop** (body มี walkthrough ไทย + Review focus +
   ผล test + verdict ของ reviewer) → **คุณรีวิว + merge** (จาก GitHub mobile ได้)
3. ห้ามเริ่ม ticket ถัดไปก่อน PR ปัจจุบันถูก merge
4. **Review tiers**: ทุก PR ระบุ Review level — 🟢 skim (scaffolding/config/docs;
   CI เขียว = merge จากสรุปได้เลย) / 🟡 normal / 🔴 careful (domain logic, auth,
   migration, SQL, LLM spend — อ่าน diff จริง) พร้อมเหตุผล 1 บรรทัด
5. **Review-as-quiz**: Review focus ใน PR เป็นคำถาม 2-4 ข้อให้ตอบระหว่างอ่าน
   (เฉลยพับไว้ท้าย PR) — และคำอธิบายยาว ๆ ในแชทจะจบด้วย quiz เสมอ
6. **Model policy**: main session = Opus (orchestrator ล้วน ไม่เขียนเอง),
   workers = opus/sonnet/haiku, **Fable = escalation เท่านั้น**
7. **Content batch**: lesson-writer → lesson-verifier → FAIL เกิน 2 รอบ = พัก concept แล้วรายงาน
8. Frontend ticket ทุกตัวใช้ skill `frontend-design`, ทุก agent ใช้ skill `token-efficiency`
9. **Mutation-test ตัว test เอง — เกือบทุก ticket ที่แตะ SQL หรือ fixture มี mutation ที่ทั้ง
   suite จับไม่ได้** (พบใน ≥10 PR: #7 fixture ASCII ล้วนซ่อน `utf8.RuneCountInString`→`len` ·
   #8 `Repository.Topics` ไม่มี test เลย 8/8 mutation รอด · #10 รอด 15/32 · #40 มี test ที่
   เทียบกับตัวเอง · #41 รอด 5/11 · #45 รอด 5 · #46 รอด 2 · UX-7 (ก่อนเปิด PR) รอด 11/18 —
   ทุกจุดที่รอดอยู่ใน component (JSX) ไม่ใช่ `curriculum.ts` เพราะ `web/` ยังไม่มี jsdom
   ตอนนั้น มี `@testing-library/react`/`jsdom` ทีหลังตอนแก้ ปิดครบ 20/20 รวม bonus 2 จุด)
   ส่วนใหญ่ code-reviewer จับก่อน merge — **มี 1 ครั้งที่หลุดขึ้น develop จริง**: #38 รายงาน
   `VIOLATIONS: NONE` ขณะที่ MCQ เดาถูก 89.6% เพราะวัด *"กฎที่ตั้งไว้ถูกละเมิดไหม"* แทน
   *"เดาด้วย heuristic ง่าย ๆ แล้วถูกกี่ %"* (แก้ใน #43)
   → **ต้องพัง invariant ของจริงแล้วบอกให้ได้ว่า test ตัวไหนจับ** ถ้าไม่มีตัวไหนจับ = test นั้นเป็นของประดับ
10. **คำกล่าวอ้างเรื่อง DOM / a11y / contrast ต้องวัดจาก browser จริง** (Playwright + headless
    Chromium ใช้ได้ในเครื่องนี้) ไม่ใช่การอ่านโค้ด — และ rebuild container ก่อนวัดสี
    (เคยเกือบ false pass เพราะอ่านจาก container เก่า)

## สิ่งที่คุณต้องเตรียมเอง

- [x] GitHub repo + PR template + CI
- [x] GitHub mobile app (ไว้รีวิว PR จากมือถือ)
- [ ] บัญชี DigitalOcean + domain — **ยังไม่ต้อง** จนกว่าจะกลับมาทำ T29/T30
- [ ] Anthropic API key **ตัวจริง** — **ยังไม่ต้อง** (ใส่ค่าอะไรก็ได้ให้ config ผ่าน;
      ยังไม่มีโค้ดที่เรียก Anthropic)
- [ ] Cloudflare R2 / Backblaze B2 สำหรับ backup — คู่กับ T33 ซึ่งพักไว้พร้อม deploy
