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
    tab ไปไม่ถึง, `cursor-default`) หัวข้อใช้ `text-muted` (7.61:1 ต่อ AA สูงกว่า `text-body`
    ปกติของแถวที่มี lesson เพื่อให้เห็นความต่างแม้บนมือถือที่ไม่มี hover) โชว์ข้อความ
    `No lesson yet` แทนเวลาอ่าน ด้วย `text-faint` (contrast 4.76:1 ผ่าน AA)
  - header ของ chapter เปลี่ยนจาก `{concepts.length} concepts` เป็น `{n}/{N} ready`
    (`n` = concept ที่มี lesson จริง, เลี่ยงคำว่า "lessons" เพราะ `N` คือจำนวน concept
    ไม่ใช่จำนวน lesson จริง) รวมถึง edge case `concepts: []` → `0/0 ready` ไม่ crash และ
    accordion ยังกดขยายได้ปกติแม้ `n === 0`; `<ul>` ของ concept list render อยู่เสมอ
    (toggle ด้วย Tailwind `hidden` แทนที่จะไม่ render เลยตอนปิด) เพื่อให้ `aria-controls`
    ของปุ่ม toggle ชี้ไปยัง element ที่มีอยู่จริงเสมอ
  - แต่ละแถว concept โชว์เลขลำดับ (`index + 1` ในรายการที่ backend sort ตาม `position`
    มาแล้ว) แทนที่ `#{position} · {slug}` เดิม — เลิกโชว์ raw slug (debug info) บนหน้านี้;
    เวลาอ่านโชว์ `~{est_minutes} min` เฉพาะตอนเป็น number ที่ `> 0` เท่านั้น — กัน `null`
    (ไม่มี estimate) และ `0` (ค่าที่ domain layer ไม่ตั้งใจเขียนแต่ schema ไม่ได้ห้ามไว้)
    แยกจากกันทั้งคู่ ไม่ render ข้อความเวลาเลยในทั้งสองกรณี
  - เพิ่ม pure function `countAvailableLessons(concepts): { available, total }` ใน
    `web/lib/curriculum.ts` ใหม่ (แยกจาก `web/lib/api.ts` ที่โฟกัส fetch/types), เรียกใช้
    ทั้งใน `TrackTopics.tsx` (header ระดับ chapter) และ `learn/page.tsx` (`computeStats`
    ระดับ track) เพื่อไม่ให้ definition ของ "lesson ready" ซ้ำกันสองที่
- **Review focus**:
  - ทำไม concept ที่ `has_lesson=false` ต้อง render เป็น `<div>` เฉย ๆ แทนที่จะเป็น
    `<a aria-disabled="true">`?
  - ทำไมต้องกันทั้งกรณี `est_minutes = null` และ `est_minutes = 0` แยกจากกัน ทั้งที่
    domain layer ไม่เคยตั้งใจเขียน 0 ลง DB?
  - ทำไมเลขลำดับหน้าแต่ละ concept ใช้ array index (`index + 1`) แทนที่จะใช้ `concept.position`
    ตรง ๆ จาก API?
  - ทำไม `<ul>` ของ concept list ต้อง render อยู่เสมอ (toggle ด้วย `hidden` class) แทนที่จะ
    conditional-render แบบ `{open && <ul>...}` เหมือนเดิม?
- Status: `merged, PR #44`

## UX-4 — API: learning progress (read + write) `[go-implementer]`

- **Scope**: bounded context ใหม่ `internal/learning/` (app/infra ต่อยอดจาก domain ที่มีอยู่แล้ว —
  `ChunkState`, `LessonProgress`, `LessonRef`, `Gate`), สอง endpoint คีย์ด้วย **topic slug +
  concept slug** เหมือน `/lesson?topic=&concept=` ไม่ใช่ `lesson_id` (curriculum API ไม่เคย
  expose lesson id):
  - `GET /api/v1/progress` — คืนเฉพาะ concept ที่มี progress row จริง (ไม่มี row = "ยังไม่เริ่ม",
    frontend เดาเอง); ว่างต้อง marshal เป็น `[]` ไม่ใช่ `null`
  - `PUT /api/v1/progress/{topic}/{concept}` body `{"state":"in_progress"|"passed"}`
    (Go 1.22+ `ServeMux` path wildcard)
  - reuse `migrations/002_learning.sql` เดิมทั้งหมด (ตาราง `lesson_progress` มีครบทุกคอลัมน์ที่
    ต้องใช้อยู่แล้ว) — **ไม่เพิ่ม migration ใหม่**
- **การตัดสินใจหลัก**:
  - **Soft-guide ไม่ gate**: API นี้ไม่มีทาง store/return `"locked"` ได้เลย — ค่านี้ยังอยู่ใน ENUM
    ของ migration 002 (เผื่อ UX-6 เอา `domain.Gate` มาต่อยอดทีหลัง) แต่ application layer reject
    ด้วย 400 ถ้า client ส่ง `"locked"` มา (เช่นเดียวกับค่าอื่นที่ไม่รู้จัก)
  - **Slug shape validate ก่อน DB round trip ใด ๆ (R1a จาก code review รอบแรก)**: MySQL
    `utf8mb4_0900_ai_ci` ทำให้ `co.slug = 'B-Trees'` แมตช์แถวที่เก็บจริงเป็น `b-trees` — ถ้าปล่อยให้
    ถึง `domain.NewLessonRef` ก่อน (ซึ่งรับแค่ lowercase) จะพังเป็น 500 ไม่ใช่ 400 ทั้งที่ต้นเหตุคือ
    client ส่ง input ผิดรูป แก้โดยเช็ค shape ของทั้ง `{topic}` และ `{concept}` (ยืม
    `domain.NewLessonRef`'s shape check มาใช้แม้ตัวมันตั้งใจแทน concept slug อย่างเดียว) **ก่อน**
    เรียก repository เลย — คืน `ErrInvalidSlug` → 400
  - **Idempotency เป็นงานของ application layer ไม่ใช่ domain**: `domain.LessonProgress.Unlock()`/
    `MarkPassed()` ยัง error เหมือนเดิมทุกอย่างเมื่อเรียกซ้ำ (ไม่ได้ไปอ่อน invariant เพื่อความสะดวก
    ของ HTTP) — `decideTransition` เรียก domain method จริง ๆ แล้วจับ
    `ErrAlreadyUnlocked`/`ErrAlreadyPassed` แปลงเป็น no-op ที่ refresh `last_read_at` อย่างเดียว
    **ข้อแก้ไขจากรอบ review แรก**: รายงานตอนแรกอธิบายผิดว่า forward-only guarantee "ได้มาฟรี" จาก
    `Unlock()` ครอบคลุมทุก path — ที่จริง path "concept ที่ยังไม่เคยมี progress row" (fresh) **ข้าม
    domain transition ไปเลย** (`ProgressDecision{State: requested}` ตรง ๆ ไม่เรียก `NewLessonProgress`
    ด้วยซ้ำ เพราะไม่มี state เดิมให้ transition จาก) ส่วน `Unlock()`'s "ไม่ใช่ locked" property อธิบาย
    ได้แค่กรณี **passed + PUT in_progress ไม่ downgrade** เท่านั้น (Unlock error เหมือนกันไม่ว่า
    current จะเป็น in_progress หรือ passed)
  - **`locked` row + PUT passed = unlock-then-pass (R1b)**: ถ้ามี row สถานะ `locked` อยู่แล้ว (เคสของ
    UX-6 ในอนาคต ที่ ticket นี้เองไม่เคยเขียน) แล้วมี PUT `passed` มา, `MarkPassed()` เดี่ยว ๆ จะ reject
    (`ErrLessonLocked`) — เพราะ soft-guide คือ "ไม่เคยปฏิเสธการอ่าน/finish" (ไม่ใช่ 409) เลยแก้เป็น
    เรียก `Unlock()` ก่อน (เพิกเฉย `ErrAlreadyUnlocked`) แล้วค่อย `MarkPassed()` — เป็นทางเดียวที่
    locked row จะไปถึง passed ได้ และเป็นครั้งแรกที่ `Unlock()`'s locked→in_progress edge ถูกใช้จริง
  - **Race ระหว่าง PUT in_progress (เปิด lesson) กับ PUT passed (Finish) — UX-5 จะยิงคู่นี้จริง (R3)**:
    ออกแบบเดิม (อ่าน progress แยก 1 query แล้วค่อยเขียนทีหลัง) มี race window ที่ทำให้ state จบที่
    `in_progress` พร้อม `first_passed_at` ไม่ว่าง (ขัดแย้งกันเอง) ถ้า in_progress เขียนทับหลัง passed
    แก้ด้วย `SELECT ... FOR UPDATE` ล็อกแถว `lessons` ไว้ตลอด 1 transaction (`Repository.Transition`)
    แล้วให้ `decide` callback (มาจาก `Service`) ตัดสินใจ **หลัง** ได้อ่านค่าที่ล็อกแล้วเท่านั้น — ไม่ใช้
    ทางลัด `state = IF(...)` ฝั่ง SQL เพราะจะย้าย domain invariant (forward-only) ไปฝังใน infra
    ซึ่งขัด DDD layering ตรง ๆ และจะมี writer ตัวที่สอง (LLM grading / UX-6) เข้ามาอีกในอนาคต —
    `domain.LessonProgress` ต้องเป็นเจ้าของกฎนี้คนเดียวเสมอ
  - **คืน slug จาก DB ไม่ใช่จาก caller (S1)**: `SELECT l.id, t.slug, co.slug ... FOR UPDATE` คืน slug
    ตามที่ DB เก็บจริงเสมอ (เป็นผลพลอยได้จากการล็อกแถว lesson ข้างต้น) — ป้องกัน mismatch ถ้า
    ai_ci collation ทำให้ query แมตช์ input คนละ case กับที่เก็บจริง (แม้ตอนนี้ R1a จะ block input
    ผิด shape ไปตั้งแต่ต้นแล้วก็ตาม กันไว้อีกชั้นที่ infra เพราะ `Repository.Transition` เองไม่ควร
    พึ่ง caller ส่ง case ถูกเสมอ)
  - `first_passed_at` ตั้งครั้งเดียว ไม่ถูกเขียนทับตอน re-finish — บังคับด้วย SQL
    `COALESCE(lesson_progress.first_passed_at, new.first_passed_at)` ใน `ON DUPLICATE KEY UPDATE`
    (ไม่ใช่แค่ logic ฝั่ง Go — กันไว้สองชั้น)
- **บทเรียนจาก mutation-testing (2 รอบ)**: stub driver dispatch ด้วย query-string identity (เทียบ
  constant กับตัวเอง เพราะ production กับ test import constant เดียวกัน) ผ่านเสมอไม่ว่า SQL text
  จะพังแค่ไหน — พังจริงทั้งหมด 5 จุดข้ามสองรอบ (ลบ `COALESCE` ออกจาก upsert, เติม `state = '...'`
  เข้าไปใน refresh's SET clause, เปลี่ยน `AND` เป็น `OR` ใน `WHERE t.slug = ? AND co.slug = ?`,
  เปลี่ยน `JOIN` เป็น `LEFT JOIN` ใน GET's `selectAllProgressSQL`) ทุกจุด behavioral test (ที่ stub
  เขียน logic เองแยกจาก SQL text จริง) เขียวผ่านหมด มีแค่ SQL-shape test ที่เทียบ **exact literal**
  ทั้งก้อน (ไม่ใช่ `strings.Contains` แยกท่อนแบบรอบแรก ซึ่งเช็ค `t.slug = ?` กับ `co.slug = ?` แยกกัน
  จับ `AND`→`OR` ไม่ได้) เท่านั้นที่จับได้ — แก้แล้ว restore กลับก่อน commit ทุกครั้ง
- **Review focus**:
  - `Service.SetProgress` มี 3 path: fresh (ไม่เคยมี row), unlock-then-pass (จาก locked),
    idempotent/forward-only no-op — path ไหนเรียก `domain.NewLessonProgress`/`Unlock`/`MarkPassed`
    จริง และ path ไหน bypass ไปเลย เพราะอะไร?
  - ทำไมการล็อกแถว `lessons` (ไม่ใช่ `lesson_progress`) ด้วย `FOR UPDATE` ถึงพอป้องกัน race แม้ตอน
    concept ยังไม่เคยมี `lesson_progress` row เลย?
  - ทำไม repository test ที่ seed ข้อมูลผ่าน stub แล้วอ่านกลับ (behavioral) ถึงจับบั๊ก SQL text
    ไม่ได้ ต้องเทียบ SQL string ทั้งก้อนแบบ exact literal แทน `strings.Contains` แยกท่อน?
  - ทำไม `lesson_progress.state` ENUM ยังเก็บค่า `'locked'` ไว้ ทั้งที่ตอนนี้ ticket นี้เขียนมันได้แล้ว
    (ผ่าน unlock-then-pass) แต่ก็ยังไม่มี path ไหนของ ticket นี้ที่ **สร้าง** row สถานะ `locked` เอง?
- **Known debt (บันทึกไว้ ยังไม่แก้ในรอบนี้)**:
  - "`first_passed_at` write-once" เป็น invariant ที่มีอยู่แค่ใน SQL (`COALESCE` ใน
    `ON DUPLICATE KEY UPDATE`) — `domain.LessonProgress` และ `app.ProgressEntry` ไม่ได้ model
    concept นี้ไว้เลย ถ้าวันหนึ่งมี writer อื่นที่ไม่ผ่าน SQL statement นี้ (เช่น import/backfill
    script) กฎนี้จะหายไปเงียบ ๆ
  - `concepts` unique key คือ `(chapter_id, slug)` ไม่ใช่ต่อ topic — สองบทใน topic เดียวกันที่ใช้
    concept slug ซ้ำกันจะทำให้ query ของ ticket นี้ (JOIN topics→chapters→concepts บน
    `t.slug + co.slug`) แมตช์ได้มากกว่า 1 แถว ตอนนี้ยังไม่มี concept slug ซ้ำแบบนี้ใน
    `content/curriculum/*.json` จริง และ identity model นี้สืบทอดมาจาก `lessonreader.go` เดิม
    (curriculum's read path) แต่ UX-4 เปลี่ยนมันจาก read-only ไปเป็น **write path** แล้ว — ถ้าจะแก้
    ต้องคิดเรื่อง unique constraint ใหม่ทั้งระบบ ไม่ใช่แค่ ticket นี้
- Status: `implemented, PR pending`

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
