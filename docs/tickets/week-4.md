# Week 4 — SRS + กราฟแรก

**เป้าหมาย**: Drill session ใช้งานได้จริงจาก SM-2 queue, dashboard มี streak + trend graph แรก

---

## T16 — Review domain: SM-2 `[go-implementer]` ~45 นาที
- **Goal**: หัวใจของ SRS — `ReviewCard.Apply(grade)` เป็น pure function
- **Scope**: `internal/review/domain` — SM2State VO (EF, interval, repetitions, due), สูตร SM-2 เต็ม (grade < 3 reset repetitions, EF ต่ำสุด 1.3), `migrations/003_review.sql`
- **Acceptance**: **table-driven test ละเอียดสุดในโปรเจกต์** — ทุก grade 0–5 × สถานะ (การ์ดใหม่/ทวนแล้ว n รอบ), ค่า expected คำนวณมือแนบใน test
- **Review focus**: เทียบสูตรกับ SM-2 spec ตรง ๆ, ปัดเศษ interval, EF ไม่ต่ำกว่า 1.3, due date ใช้ DATE ไม่ใช่ DATETIME
- Status: `todo`

## T17 — Event dispatcher + RecallFailed → ReviewCard `[go-implementer]` ~45 นาที
- **Goal**: กลไก in-process events ที่ทุก context ใช้ร่วม + subscriber แรก
- **Scope**: `internal/platform/events` (sync dispatcher + append ตาราง events), subscriber: `RecallFailed` → upsert ReviewCard (source=recall_check), แก้ T13 ให้ยิงผ่าน dispatcher
- **Acceptance**: recall ตก → การ์ดโผล่ใน queue พรุ่งนี้, ตกซ้ำ → ไม่สร้างการ์ดซ้ำ (unique source), tests ของ dispatcher + subscriber
- **Review focus**: dispatcher error ของ subscriber หนึ่งไม่ทำให้ request หลัก fail (log แล้วไปต่อ? หรือ tx เดียว? — ดู decision ใน code แล้วถามว่าทำไม), ordering ของ subscribers
- Status: `todo`

## T18 — Drill queue + review endpoints `[go-implementer]` ~40 นาที
- **Scope**: `GET /drill/queue?limit=` (การ์ด due วันนี้ + เนื้อหา join มาให้), `POST /review-cards/{id}/review {grade}` → Apply + log + คืนการ์ดถัดไป
- **Acceptance**: การ์ดยังไม่ due ไม่โผล่, review แล้ว due ขยับตาม SM-2, tests
- **Review focus**: query ORDER BY (due เก่าสุดก่อน), join ข้าม context ทำที่ infra layer (read model) ไม่ใช่ลาก domain ข้าม context
- Status: `todo`

## T19 — Drill UI `[go-implementer]` ~45 นาที
- **Scope**: `web/app/drill/` — โชว์คำถาม → คิด → เปิดเฉลย → ให้คะแนนตัวเอง 0–5 (สำหรับการ์ด recall ใช้คำตอบเดิมเป็นเฉลย) → การ์ดถัดไป, ใช้ timer 30 นาทีร่วมกับ T15
- **Acceptance**: flow ลื่น keyboard-first (1–5 = grade, space = เปิดเฉลย), จบ queue มีหน้าสรุป
- **Review focus**: keyboard shortcuts ตาม skill, ปุ่ม grade อธิบายความหมาย (0=จำไม่ได้เลย … 5=ง่ายมาก), บนมือถือมีปุ่มกดแทนคีย์บอร์ด
- Status: `todo`

## T20 — Measurements + daily_activity projections `[go-implementer]` ~45 นาที
- **Goal**: ทุกการวัด = 1 จุดกราฟ (กติกา "graph after every measurement")
- **Scope**: subscribers: RecallAnswered → measurement(recall_accuracy), CardReviewed → measurement(retention), SessionCompleted → daily_activity upsert; endpoints `GET /stats/trends?metric=&topic=`, `GET /stats/dashboard` (streak + สรุปสัปดาห์)
- **Acceptance**: ทำ activity แล้วจุดกราฟ/streak ขยับจริง, streak นับตาม timezone Asia/Bangkok, tests ของ projection logic
- **Review focus**: timezone (DB เก็บ UTC แต่ตัด "วัน" ที่ Asia/Bangkok — ดูว่าตัดตรงไหน), projection idempotent ไหมถ้า event ซ้ำ
- Status: `todo`

## T21 — Dashboard v1: streak + trend chart `[go-implementer]` ~45 นาที
- **Scope**: `web/app/dashboard/` + `components/charts/TrendChart.tsx` (reusable — ทุกหน้าที่มีการวัดจะใช้ตัวนี้), streak card, สรุปนาที/สัปดาห์
- **Acceptance**: chart แสดง score + trend line, reusable รับ props (metric, topic), build ผ่าน
- **Review focus**: **chart ตาม skill `frontend-design` หมวด charts** (palette เดียวกันทุกกราฟ, empty state สวยตอนข้อมูลน้อย), ไม่ลาก chart lib ใหญ่เกินจำเป็น
- Status: `todo`

## C2 — Content: DDD บท 2 (5 จาก 10) + DistSys บท 1 (4) `[lesson-writer → lesson-verifier]`
- DDD building-blocks ครึ่งแรก: layered-architecture, entities, value-objects, domain-services, modules
- DistSys foundations: why-distributed, transparency-goals, scalability-dimensions, fallacies-of-distributed-computing
- Status: `todo`
