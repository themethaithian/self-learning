# DDIA track (6th track) — tickets

Track ที่ 6: **Designing Data-Intensive Applications** (Kleppmann). เพิ่มแบบ
additive หลัง v1 — radar กลายเป็น 6 แกน. โครงสร้าง tree ที่ user อนุมัติแล้ว
(12 บท / 61 concept) เป็น contract: slug ห้ามเปลี่ยนหลัง import (เหมือน T7/T7b).

## T-ddia-track — enable track (DONE, PR #18 merged 2026-07-23)

- domain `Track` array + `Tracks()` ได้ `{value:"ddia"}` (ต่อท้ายสุด)
- `migrations/003_ddia-track.sql` ขยาย ENUM `topics.track`
- verify บน MySQL จริงแล้ว track `ddia` live บน develop

## T-ddia-curriculum — content/curriculum/ddia.json + loader test (DONE, PR #19 merged)

**Status:** DONE — PR #19 merged. ddia.json (12 บท/61 concept) + loader test เข้า develop แล้ว.

สร้างไฟล์ tree ของ track `ddia` ให้ import-curriculum โหลดได้ + ล็อก slug ด้วย
loader test เหมือนทุก track.

Acceptance:
- `content/curriculum/ddia.json`: track `ddia`, topic slug
  `designing-data-intensive-applications`, title `Designing Data-Intensive Applications`,
  position 6. โครงสร้าง 12 บท / 61 concept ตาม tree (slug ตรงเป๊ะ, position 1-based
  เรียงตามลำดับใน tree). แต่ละ concept มี `outline` ภาษาไทย 2–3 บรรทัด (สไตล์เดียว
  กับ `ddd.json`: บอกว่าบทเรียนควรสอนอะไร + เชื่อมโยง trade-off / "when not to use").
  ชื่อบทใช้ชื่อบทจริงในหนังสือ (เช่น ch.8 "The Trouble with Distributed Systems",
  ch.4 "Encoding and Evolution").
- `internal/curriculum/infra/contentfile_ddia_test.go`: mirror
  `TestLoadTopic_DistsysContentFile` — `wantDDIAChapters` ระบุ slug ทุกบท/concept
  ตาม tree, assert track/topic slug/จำนวนบท/จำนวน concept/position.
- `go vet ./...` + `go test ./...` เขียว. (ไม่รัน import-curriculum ที่นี่ — ต้องมี DB;
  loader test ครอบคลุมความถูกต้องของไฟล์แล้ว)

**Review focus (ตอบระหว่างอ่าน diff):**
1. ทำไม loader test ถึงจำเป็นทั้งที่ JSON ก็ parse ได้อยู่แล้ว — มันกันความผิดพลาด
   ประเภทไหนที่ compiler จับไม่ได้?
2. ถ้าอนาคตต้องลบ concept หนึ่งออกจาก ddia.json ตอน import-curriculum จะเกิดอะไรขึ้น
   ถ้า concept นั้นมี lesson อ้างอยู่?
3. ทำไม slug ถึงถือเป็น contract ที่เปลี่ยนไม่ได้หลัง import (คิดถึง URL ของหน้า
   lesson + progress ที่ผูกกับ concept)?

## T-ddia-lessons — lesson batches (batch = 1 บท)

lesson-writer → lesson-verifier ทีละ concept. ทุก lesson มี section trade-off /
"when not to use". ไฟล์ลง `content/lessons/designing-data-intensive-applications/<concept>.json`
(topic slug = folder name, บังคับโดย `LoadLesson`), import ด้วย `cmd/import-lessons`.

- **batch 1 — ch.5 Replication (5 lessons):** single-leader, replication-lag-problems,
  multi-leader, leaderless-dynamo-style, quorum-consistency. verifier PASS ครบ 5.
  → **branch `content/ddia-replication`, PR รอ merge.**
- **batch 2 — ch.7 Transactions (5 lessons):** acid-meaning, read-committed,
  snapshot-isolation-mvcc, lost-updates-and-write-skew, serializability. (NEXT)
- ที่เหลืออีก 10 บท ตาม tree ใน `content/curriculum/ddia.json`.
