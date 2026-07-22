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
