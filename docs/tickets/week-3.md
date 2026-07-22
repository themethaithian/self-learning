# Week 3 — Read & Recall จบ end-to-end + กันข้อมูลหาย

**เป้าหมาย**: อ่าน lesson จริง (DDD บท 1) บนเว็บจริงจากมือถือได้ → ตอบ recall → Haiku ตรวจ →
ปลดล็อกบทถัดไป พร้อม timer และ backup ทำงานก่อนมีข้อมูลจริงสะสม

---

## T12 — Session domain + endpoints `[go-implementer]` ~45 นาที
- **Goal**: session runner: เริ่ม/จบ/ทิ้ง พร้อม invariant "active ทีละ 1"
- **Scope**: `internal/learning/domain` (Session state machine), app service, infra handler+repo, routes `POST /sessions`, `GET /sessions/active`, `POST /sessions/{id}/complete|abandon`
- **Acceptance**: เปิด session ซ้อน → 409, state ผิด transition → error, tests ครอบ state machine ทุกเส้น
- **Review focus**: race ตอนเปิด session พร้อมกัน (unique constraint หรือ `SELECT … FOR UPDATE`), เวลาใช้ UTC ใน DB
- **ห้าม copy T5 แบบหลับตา** (จาก deep review ของ T5): T5 เป็น reference data จึงมีแต่ constructor + getter ถ้าลอกทั้งดุ้นจะได้ anemic domain — context นี้ต้องมีเพิ่ม 3 อย่าง (1) state transition เป็น method ที่คืนค่าใหม่ `func (s Session) Complete(now time.Time) (Session, error)` ไม่ใช่ setter และไม่ใช่ logic ที่ app service (2) **surrogate identity VO** `SessionID` — "ไม่มี id ใน domain" ของ T5 เป็น choice เฉพาะ reference data ที่มี natural key ไม่ใช่กฎของโปรเจกต์ route `POST /sessions/{id}/complete` ต้องการ id จริง (3) invariant "active ทีละ 1" span หลาย instance → constructor บังคับไม่ได้ ต้องอยู่ app layer + DB constraint
- Status: `todo`

## T13 — Recall attempt endpoint + gating `[go-implementer]` ~45 นาที
- **Goal**: เส้นครบวงจร: ตอบ → Haiku ตรวจ → บันทึก → ปลดล็อก lesson ถัดไป
- **Scope**: `POST /recall-checks/{id}/attempts` — เรียก Grader port, เก็บ recall_attempts, อัป lesson_progress, บันทึก events `RecallAnswered`/`RecallFailed` ลงตาราง events ตรง ๆ (dispatcher มาใน T17 สัปดาห์ 4)
- **Acceptance**: grade ≥ 3 = passed + คืน `next_unlocked`, LLM ล่ม → 502 พร้อมข้อความให้ retry (คำตอบ user ไม่หาย), tests mock Grader
- **Review focus**: ทั้งก้อนอยู่ใน DB transaction เดียวไหม, idempotency ถ้ายิงซ้ำ, คำตอบ user ต้อง save ก่อนเรียก LLM
- Status: `todo`

## T14 — Read UI: lesson + recall gate `[go-implementer]` ~45 นาที
- **Goal**: หน้าอ่านหลัก — อ่าน markdown ไทย + diagram → ตอบ recall → เห็น feedback → ไปต่อ
- **Scope**: `web/app/read/` — lesson viewer (render markdown + **mermaid diagram** ใน code fence), section "อ่านต่อของจริง" ท้ายบท (references: title + source + why), recall form (short answer + mcq), แสดง grade/feedback, ปุ่มไป lesson ถัดไป (disabled จนกว่าผ่าน)
- **Acceptance**: flow ครบด้วยข้อมูลจริงจาก C1, diagram เรนเดอร์สวยทั้งจอเล็ก/ใหญ่, `npm run build` ผ่าน
- **Review focus**: mermaid lazy-load (อย่าให้ bundle บวมทุกหน้า), ฟอนต์ไทยตาม skill `frontend-design`, สถานะ "กำลังตรวจ..." ระหว่างรอ LLM, Ctrl+Enter submit
- Status: `todo`

## T15 — Session timer UI + summary `[go-implementer]` ~40 นาที
- **Goal**: timer 30 นาทีที่เห็นตลอด + สรุปเมื่อจบ session
- **Scope**: timer component (นับถอยหลัง, เตือนตอน 5 นาทีสุดท้าย, ครบเวลาแจ้งแต่ไม่บังคับหยุด), summary screen (lessons อ่าน, recall ผ่าน/ตก)
- **Acceptance**: timer ไม่ reset ตอน navigate ภายใน session, refresh หน้า timer ยังถูกต้อง (คำนวณจาก started_at ของ API)
- **Review focus**: อย่าเก็บ state เวลาไว้ฝั่ง client อย่างเดียว (source of truth คือ `started_at`), timer สงบไม่กดดัน (ตาม skill)
- Status: `todo`

## T33 — Backup cron + restore test `[go-implementer]` ~40 นาที
- **Goal**: มี backup ก่อนข้อมูลเรียนจริงเริ่มสะสม (ย้ายขึ้นมาจากสัปดาห์ 6 เดิม)
- **Scope**: `deploy/backup.sh` (mysqldump → gzip → rclone ไป R2/B2, เก็บ 14 วัน), cron ทุกคืน 03:00 ICT, เพิ่มวิธี restore ใน runbook
- **Acceptance**: **restore จริง 1 ครั้ง** ลง MySQL เปล่าแล้วข้อมูลครบ — backup ที่ไม่เคย restore = ไม่มี backup
- **Review focus**: dump ใช้ `--single-transaction`, ไฟล์ backup ไม่ world-readable, ลบไฟล์เก่าจริง
- Status: `todo`

## C1 — Content: DDD บท 1 (4 lessons) `[lesson-writer → lesson-verifier]`
- Concepts: ubiquitous-language, model-driven-design, knowledge-crunching, hands-on-modelers
- Batch แรกที่ใช้ format ใหม่: มี mermaid diagram (เมื่อเหมาะ) + references 2–4 แหล่ง
- Flow: writer → verifier → (FAIL ≤ 2 รอบ) → import → รายงาน batch
- **คุณรีวิว**: อ่าน 1 lesson เต็ม ๆ เช็คภาษา/ความยาก/diagram/references ก่อนอนุมัติ format นี้ยาว ๆ
- Status: `todo`
