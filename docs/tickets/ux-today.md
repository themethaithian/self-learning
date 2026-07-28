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
- Status: `merged, PR #41`

## UX-2 — Learn page: track cards + focus picker `[go-implementer]`

- **Scope**: หน้า `/learn` ใหม่แทน `/read` — 7 `TrackCard` (ชื่อเต็มจาก `trackMeta.ts`,
  จำนวน lesson ready/chapter/concept, focus chip/ปุ่ม `Set focus`), track ที่ยังไม่มี
  lesson เลยไปอยู่ section "Coming soon" (ไม่ใช่ card, กดไม่ได้), focus picker ผูกกับ
  `GET|PUT /api/v1/prefs/focus-track` (จาก UX-1b) แบบ optimistic update + revert เมื่อ
  PUT fail, `/read` และ root/`/token` redirect ไป `/learn` (client-side redirect —
  static export ไม่มี server redirect ให้ใช้)
- **บทเรียนสำหรับ ticket ถัดไป**:
  - progress bar ที่ตัวเศษ = ตัวส่วนเสมอ (บาร์เต็ม 100% ทุกใบ) แย่กว่าไม่มีบาร์เลย เพราะสื่อว่า
    "เรียนจบแล้ว" ทั้งที่ยังไม่มี progress tracking จริง (v1 เลยไม่ render bar เลย โชว์
    ตัวเลขตรง ๆ แทน จนกว่า progress API มา)
  - วนจาก allowlist ที่ hardcode (เช่น track display order) ไปหา data ทำให้ track/ข้อมูล
    ที่ backend ส่งมาแต่ยังไม่อยู่ใน list **หายไปเงียบ ๆ ไม่มี error** — ต้องวนจาก data จริง
    เสมอ แล้วใช้ list เป็นแค่ sort key เท่านั้น
- **Review focus**:
  - ทำไม focus toggle บน `TrackCard` ต้องใช้ `aria-disabled` แทน `disabled` attribute ตรง ๆ?
  - ทำไม `ProgressBar.tsx` ถึงยังเก็บไว้ในโค้ดทั้งที่ไม่มีหน้าไหนเรียกใช้ตอนนี้?
- Status: `merged, PR #42`

## UX-3 — Track view: chapter accordion + availability

- หน้า track detail (`/learn?track=`) มี gap: concept ที่ `has_lesson=false` ยังเป็นลิงก์อยู่
  (กดแล้วเจอ "no lesson yet" ที่ `/lesson`) และ `ChapterRow` โชว์จำนวน concept รวมของ chapter
  แทนจำนวนที่มี lesson จริง — ticket นี้ปิด gap ทั้งสองจุด:
  - concept ที่ `has_lesson=false` render เป็น `<div>` เฉย ๆ (ไม่ใช่ `<Link>`, ไม่มี `href`,
    tab ไปไม่ถึง) โชว์ข้อความ `No lesson yet` แทนเวลาอ่าน, ใช้ `text-faint` เดิม (contrast
    4.76:1 ผ่าน AA อยู่แล้ว ไม่ต้องเพิ่ม token ใหม่)
  - header ของ chapter เปลี่ยนจาก `{concepts.length} concepts` เป็น `{n}/{N} lessons`
    (`n` = concept ที่มี lesson จริง) รวมถึง edge case `concepts: []` → `0/0 lessons`
    ไม่ crash และ accordion ยังกดขยายได้ปกติแม้ `n === 0`
  - แต่ละแถว concept โชว์เลขลำดับ (`index + 1` ในรายการที่ backend sort ตาม `position`
    มาแล้ว) แทนที่ `#{position} · {slug}` เดิม — เลิกโชว์ raw slug (debug info) บนหน้านี้;
    เวลาอ่านโชว์ `~{est_minutes} min` เฉพาะตอนเป็น number จริง, เป็น `null` ไม่ render
    ข้อความเวลาเลย (ไม่ใช่ `~0 min`)
  - เพิ่ม pure function `countAvailableLessons(concepts): { available, total }` ใน
    `web/lib/curriculum.ts` ใหม่ (แยกจาก `web/lib/api.ts` ที่โฟกัส fetch/types)
- Status: `implemented, PR pending`

## UX-4 — API: learning progress (read + write)

- Endpoint `GET /api/v1/progress` (ดึง history ของ concept ที่ผู้ใช้เคยอ่านเพื่อรู้ "อ่านตัวนี้แล้วหรือยัง"),
  `POST /api/v1/progress` (save recall grade ตัวต่อตัวหลังจากปล่อย answer)
- Status: `ยังไม่เริ่ม`

## UX-5 — Reader loop: breadcrumb + finish + next

- Implement reader flow: breadcrumb (back to chapter), finish + next button (ปะปนกับ progress save),
  next-up link/pill ชี้ concept ถัดไป
- Status: `ยังไม่เริ่ม`

## UX-6 — Today page + IA switch

- ไล่ priority: focus track (จาก UX-1b) → concept ถัดไปที่ยังไม่ passed ตาม
  position → fallback track อื่นถ้า focus track เคลียร์หมดแล้ว
- IA switch: ปุ่มสลับ focus track + visibility ของ track อื่น ๆ เพื่อ "เช็คสิ่งที่ยังค้างตามหลัง"
- รอ UX-1b merge ก่อนเริ่ม
- Status: `ยังไม่เริ่ม`

## UX-7 — Progress page + chapter/track indicators

- Progress indicator ที่เห็นคร่าว ๆ: ของ chapter (แต่ละบทเรียนเป็นไหนแล้ว) + ของ track (overview),
  โชว์บน layout ทั่วแอป (progress bar, % เลยน้อย ๆ ที่หน้า reader หรือ learn)
- Status: `ยังไม่เริ่ม`
