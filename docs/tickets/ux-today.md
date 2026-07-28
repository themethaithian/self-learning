# /today page (additive) — tickets

Epic เล็ก: หน้า `/today` ตอบ "วันนี้อ่านอะไรต่อ" โดยไล่จาก track ที่ผู้ใช้เลือกเป็น
**focus** ก่อน แล้วค่อย fallback track อื่น. ผู้ใช้เลือก focus track เอง (ไม่ hardcode
priority), เก็บที่ API/DB เพราะอ่านสลับมือถือ/desktop.

## UX-1 — curriculum tree: expose has_lesson/est_minutes `[go-implementer]` (merged, PR #40)

- `GET /api/v1/curriculum` แต่ละ concept มี `has_lesson` + `est_minutes` ให้ frontend
  รู้ว่า concept ไหนมี lesson พร้อมอ่านโดยไม่ต้อง round-trip แยก
- บทเรียนสำหรับ ticket ถัดไป: stub SQL driver ใน test ไม่พาร์ส query จริง — repository
  test ที่พึ่ง fixture อย่างเดียวเขียวได้แม้ SQL upsert พัง ต้อง assert query text
  + มือ mutation-test เอง

## UX-1b — API: focus track preference `[go-implementer]`

- **Scope**: `migrations/005_prefs.sql` (ตาราง name/value ทั่วไป `prefs`),
  bounded context ใหม่ `internal/prefs/` (domain/app/infra), endpoint
  `GET|PUT /api/v1/prefs/focus-track` (validate เฉพาะ slug format, ไม่แตะ
  `internal/curriculum`), wire ใน `cmd/api/main.go` หลัง auth เหมือน endpoint อื่น
- **การตัดสินใจหลัก**: `prefs` เป็น name/value เพื่อ preference ตัวถัดไปไม่ต้อง
  migrate ใหม่; clear focus = **ลบแถว** ไม่ใช่เก็บ NULL (ทำให้ "ไม่เคยตั้ง" กับ
  "เลิกโฟกัสแล้ว" เป็น code path เดียวกัน, `value` คอลัมน์เก็บ `NOT NULL` ได้)
- **รอบ review**: code-reviewer เจอบั๊กจริง (JSON decode แยก "ส่ง null" กับ
  "พิมพ์ field ผิด" ไม่ได้ → ลบ pref เงียบ ๆ), CORS `PUT` หลุดจาก allowed-methods
  (endpoint ใช้จากเว็บไม่ได้เลยถ้าไม่แก้), และ mutation-testing เจอ comment/test
  หลายจุดที่ไม่ได้ล็อก invariant ไว้จริง (persistence key, table name, `IF NOT
  EXISTS`) — แก้ครบตามรอบ 2/3 ก่อน merge
- **Review focus**:
  - ทำไม validate แค่ slug format ห้ามเช็คว่า track มีจริงใน curriculum?
  - ทำไม clear focus ต้องลบแถวแทนที่จะเก็บ `value = NULL`?
- Status: `in review` (ดูสรุปรอบ 2/3 fix ใน PR)

## UX-2 — `/today` page (ยังไม่เริ่ม)

- ไล่ priority: focus track (จาก UX-1b) → concept ถัดไปที่ยังไม่ passed ตาม
  position → fallback track อื่นถ้า focus track เคลียร์หมดแล้ว
- รอ UX-1b merge ก่อนเริ่ม
- Status: `todo`
