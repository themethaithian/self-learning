# Week 6 — Pre-test + stats ครบ + buffer

**เป้าหมาย**: ทำ pre-test เห็น radar baseline, weekly retro ใช้งานได้, ปิดงานค้างจาก review
(ถ้าแน่นเกิน: T31/T34 เลื่อนเป็นสัปดาห์ 7 ได้โดยไม่กระทบอย่างอื่น)

---

## T26 — Assessment domain + import + endpoints `[go-implementer]` ~45 นาที
- **Scope**: `internal/assessment` ครบ 3 layers, `migrations/005_assessment.sql`, `cmd/import-tests`, `GET /tests/{topicSlug}` (ไม่ส่งเฉลย), `POST /tests/{topicSlug}/attempts` → ตรวจฝั่ง server → score + เฉลย + event `TestCompleted`
- **Acceptance**: bank immutable (import ซ้ำไม่แก้ข้อสอบเดิม), attempt เก็บ kind=baseline/retest อัตโนมัติ (มี baseline แล้ว = retest), tests
- **Review focus**: เฉลยหลุดใน response GET ไหม (สำคัญ), การนับ delta baseline vs retest
- Status: `todo`

## T27 — Pre-test UI + radar chart `[go-implementer]` ~45 นาที
- **Scope**: `web/app/test/` — เลือก topic → ทำ 15–20 ข้อ → หน้าผล (คะแนน + เฉลย + อธิบาย), `components/charts/RadarChart.tsx` ใน dashboard (5 แกน = 5 track, ซ้อน baseline vs ล่าสุด)
- **Acceptance**: ทำครบ flow, radar โชว์ 2 ชั้นเมื่อมี retest, build ผ่าน
- **Review focus**: radar ตาม skill charts, สถานะ "ยังไม่ได้ทำ pre-test" ของ track ที่ยังว่าง
- Status: `todo`

## T31 — Weekly retro `[go-implementer]` ~45 นาที
- **Scope**: `GET /stats/retro?week=` — สรุปจาก event log: นาทีรวม/ประเภท session, recall accuracy, การ์ดทวน, DSA ต่อ pattern, เทียบสัปดาห์ก่อน; `web/app/dashboard/retro`
- **Acceptance**: เลขตรงกับข้อมูลจริงของสัปดาห์ที่ผ่านมา, สัปดาห์ว่างมี empty state
- **Review focus**: นิยาม "สัปดาห์" (จันทร์เริ่ม, Asia/Bangkok) consistent กับ streak ใน T20
- Status: `todo`

## T32 — Tickets table + Build session UI `[go-implementer]` ~40 นาที
- **Goal**: Build session ในแอปดึง ticket ถัดไปจาก plan นี้มาแสดงพร้อม timer 45 นาที
- **Scope**: `migrations/006_tickets.sql`, script แปลง `docs/tickets/week-*.md` → seed ตาราง tickets, `GET /tickets/next`, `PATCH /tickets/{id}`, `web/app/tickets/`
- **Acceptance**: หน้า Build โชว์ ticket ถัดไป + ติ๊กเสร็จได้, sync สถานะกลับ roadmap ทำมือ (พอ)
- **Review focus**: parser markdown → tickets ไม่ fragile เกิน (ยึด format หัวข้อ `## T<n>`)
- Status: `todo`

## T34 — Buffer: ปิดงานค้าง `[go-implementer]` ~45 นาที
- **Scope**: เก็บ `suggested` items ที่ค้างจาก code-reviewer ตลอด 5 สัปดาห์ + บั๊กที่เจอระหว่างใช้จริง — orchestrator รวบรวม list ให้ก่อนเริ่ม
- **ค้างจาก T4** (code-reviewer, non-blocking):
  - `config.go` CORS origin validation ยังหลุด query/fragment/userinfo/scheme-case — ใช้ round-trip `(&url.URL{Scheme,Host}).String() != raw` จับได้ทีเดียวหมด (เคส userinfo ทำให้ password หลุดใน error message)
  - `cors.go` ไม่มี `Access-Control-Expose-Headers` — จะกัดตอน ticket แรกที่เพิ่ม pagination/request-id header (frontend อ่านค่าไม่ได้ เงียบ ๆ)
  - `TestCORS` ใช้ `Header().Get("Vary")` (ได้ค่าแรกค่าเดียว) → เปลี่ยนเป็น `Values()` + `slices.Contains`
  - `TestBearerAuthPanicsOnEmptyToken` ไม่ได้ assert ข้อความ panic
  - `TrimSpace` ให้ `API_BEARER_TOKEN` ตอนโหลด config (token ที่มีเว้นวรรคนำ = 401 ทุก request แบบงง ๆ)
- **ค้างจาก T5** (code-reviewer, non-blocking):
  - `ErrInvalidTitle`/`ErrInvalidOutline` แยกตามแกน field แต่รวม empty กับ too-long ไว้ด้วยกัน — ถ้า UI อยากแยก "กรอกด้วย" กับ "สั้นลงหน่อย" ต้องกลับมาแตะ (ตัดสินใจก่อน copy pattern ไป context อื่น)
  - error message ซ้อน context ซ้ำ (`... title: empty: invalid title`) — ให้ `validateBounded` คืน sentinel ตรง ๆ แล้วให้ caller ใส่ context ชั้นเดียว
  - `Slug` เป็น type เดียวใช้ทั้ง 3 ชั้น — compiler จับไม่ได้ถ้า T6 ส่ง chapter slug ไปตำแหน่งที่ต้องการ concept slug
  - `Track` hardcode 5 ค่าใน domain + ENUM ใน migration 001 ขัดกับ design §6 ที่บอกว่าเพิ่ม track ได้โดยไม่แก้โค้ด — ตัดสินใจว่ายอมรับ (เพิ่ม track = migration + แก้ 2 จุด) หรือเปลี่ยนเป็น lookup table
  - truncate ค่า raw ที่ echo ใน error message (`slug.go` ใส่ `%q` ของ input ดิบ — input 300 ตัวอักษรได้ error ยาว 333 ตัวอักษรไหลเข้า logging middleware)
- **ค้างจาก T6** (code-reviewer, non-blocking):
  - `GET /curriculum` ดึง `co.outline` มาแล้วทิ้ง — วัดจริงบน tree 175 concepts (outline ไทย 400 runes) = ดึงจาก MySQL 222 KB ในนั้นเป็น outline 210 KB (94%) เพื่อตอบ 13 KB **ขยาย ~17 เท่า** สาเหตุคือ rehydrate ผ่าน `NewConcept` ที่บังคับ outline ไม่ว่าง → ตัดสินใจว่ายอมรับ (single user) หรือทำ read projection แยกจาก domain (CQRS-lite)
  - schema ถือข้อมูลที่ domain ไม่ยอมรับได้: `position INT` รับ 0/ติดลบ, `title VARCHAR(255)` รับ `''`, `slug` รับตัวพิมพ์ใหญ่ — และเพราะ assemble เป็น all-or-nothing แถวเสียแถวเดียว 500 ทั้ง 175 concepts → พิจารณา `CHECK` constraint (MySQL 8.0.16+) + `UNIQUE KEY (track, position)` บน topics
- **ค้างจาก T8** (code-reviewer + visual check, non-blocking):
  - `CurriculumTree.tsx` `ChapterRow` — ที่จอ 375px ถ้า title ยาวจนวรรค (เช่น "Supple Design & Refactoring") badge "N concepts" จะตกลงมาชิดซ้ายแทนที่จะชิดขวา (flex 2 children + `justify-between` พอวรรคแล้ว child ที่สองไปชิดซ้าย) — cosmetic, แก้ด้วยการจัด layout ใหม่ (เช่น badge เป็น shrink-0 หรือใช้ grid)
  - `ci.yml` web job รันทุก PR ไม่มี path filter (จงใจ — path-filtered required check อาจค้าง "expected but never ran" บล็อก merge) ถ้าอยาก optimize ใช้ `dorny/paths-filter` ไม่ใช่ top-level `paths:`
- **ค้างจาก T9** (code-reviewer, non-blocking / watch-item):
  - `LessonProgress` ไม่ถือ `first_passed_at`/`last_read_at` — วันนี้ถูกต้อง (ไม่มี invariant ไหนอ่านมัน, repository stamp เอง) แต่**ถ้าวันหน้ามีกฎที่ใช้ "passed เมื่อไหร่" เป็น input ตัดสินใจ** (streak, "passed ใน session", retro delta) ต้องย้าย timestamp เข้า aggregate ไม่งั้น logic "first pass พิเศษ" จะรั่วไป app/infra
- **ค้างจาก T10** (code-reviewer, non-blocking):
  - `cmd/import-lessons` ไม่มี app-layer `ImportService` (orchestration load→save→summary อยู่ใน `run()`) → two-level walk / zero-files / first-failure / skip stray file ไม่มี unit test (verify live แล้วถูกทั้งหมด) — แยก service + test `lessonFiles` กับ `t.TempDir()` จะปิด gap (เหมือน import-curriculum)
- **Acceptance**: list ที่ตกลงกันไว้เคลียร์หมดหรือมีเหตุผลที่ข้าม
- Status: `todo`

## C5 — Content: Diagnostic banks DDD + Go (15–20 ข้อ/topic) `[lesson-writer → lesson-verifier]`
- MCQ พร้อมคำอธิบายเฉลย ครอบทุกบทของ track — verifier เช็คเฉลยถูกจริง + ตัวลวงสมเหตุผล
- Status: `todo`

## C6 — Content: Diagnostic banks DistSys + AWS + DSA, DistSys บท 2 `[lesson-writer → lesson-verifier]`
- ปิด pre-test ให้ครบ 5 track + เดินหน้า lesson ตาม cadence ปกติ
- Status: `todo`

---

## หลังสัปดาห์ 6 (steady state)

- Content pipeline เดินต่อ batch ละ 1 chapter/pattern ตามความเร็วอ่านจริง (นำหน้า 1 chapter เสมอ)
- Retest ทุก 4–6 สัปดาห์ → radar delta
- Mock session = v2
- เมื่อรู้สึกว่าเนื้อหา "ใช้ได้จริง" แล้ว → Phase 2: publish เป็น portfolio (design.md §8)
