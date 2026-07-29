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
- **รอบ review ที่ 3 (เพิ่มเติม)**:
  - `selectProgressByLessonIDSQL` (อ่าน progress หลังล็อก) ไม่มี shape test มาก่อน — ลอง mutate
    `WHERE lesson_id = ?` → `WHERE id = ?` (คอลัมน์ทั้งคู่เป็น `BIGINT UNSIGNED` บนตารางเดียวกัน)
    behavioral test ทั้งชุดเขียวผ่านหมด มีแค่ `TestSelectProgressByLessonIDSQLShape` (เพิ่มใหม่รอบนี้)
    เท่านั้นที่จับได้ — บั๊ก class เดียวกับรอบ 1's R2 เป๊ะ
  - `FOR UPDATE` เดี่ยว ๆ ล็อกทุกตารางใน JOIN (ไม่ใช่แค่ `lessons`) — เพราะ `concepts.slug` ไม่มี
    index เดี่ยว (unique key คือ `(chapter_id, slug)`) แปลว่า query น่าจะไล่ `topics→chapters→concepts`
    แล้ว X-lock ทุกแถว `chapters` ของ topic นั้นทั้งหมด บวก gap lock ใน `concepts` — เปลี่ยนเป็น
    `FOR UPDATE OF l` (MySQL 8.0.1+, 8.4 ที่ pin ไว้รองรับ) ให้ล็อกเฉพาะแถว `lessons` จริง ๆ
    ยืนยันแล้วว่า parse ผ่านและยัง serialize ได้จริงบน MySQL 8.4 ของ stack (ดู e2e ด้านล่าง)
  - `lockLesson` เปลี่ยนจาก `QueryRow` เป็น `Query` + เช็คจำนวนแถว — ถ้า concept slug ซ้ำกันข้าม
    chapter ในหนึ่ง topic (ซึ่ง schema อนุญาตจริง ๆ) จะ error ทันทีแทนที่จะเงียบ ๆ เลือกแถวใดแถวหนึ่ง
  - `decide` callback ตัด `lessonExists bool` ออก — infra คืน `ErrLessonNotFound` ตรง ๆ เมื่อ
    `lockLesson` หา lesson ไม่เจอ ไม่ต้องเรียก `decide` เพื่อ "ถาม" error string อีกต่อไป (ของเดิมมี
    branch `if decideErr == nil` ที่ unreachable และถ้าวันหน้ามี writer ตัวที่สอง (`Grader` port)
    ที่ decide คืน `nil` เผลอ จะกลายเป็น 500 ที่ handler map เป็น 404 ไม่ได้)
  - shape check ของ slug ดึงออกมาเป็น `domain.IsValidSlugShape` (exported, มี test ของตัวเอง)
    แทนที่จะยืม `domain.NewLessonRef` (concept-specific ตาม doc comment) มาเช็ค topic slug แล้วทิ้งค่า
  - stub driver: เช็ค `c.tx == nil` **ก่อน** acquire lock เสมอ (ไม่ใช่หลัง) — กัน regression ที่ลบ
    transaction ทิ้งจากที่จะ hang การ test 10 นาทีแทนที่จะ fail ทันที
  - concurrency test (`TestRepositoryTransition_ConcurrentRaceNeverContradicts`) แก้ 2 จุด: (1)
    comment เดิมเรียกตัวเองว่า "R3's proof" เกินจริง — ที่จริงพิสูจน์แค่ข้อความมีเงื่อนไข: **ถ้า** MySQL
    serialize ที่ `FOR UPDATE` จริง (พิสูจน์ด้วย reasoning จาก InnoDB semantics ไม่ใช่ test ไหนเลย
    รวมถึง curl 10 คู่ใน e2e ที่ timing สั้นเกินจะ interleave จริง) **แล้ว** logic ถึงจะ converge ที่
    `passed`; (2) เดิม assert แค่ "ไม่ contradictory" ทำให้ถ้า `passed` goroutine error เงียบ ๆ ทุกครั้ง
    (เหลือ `state=in_progress, first_passed_at=nil`) test จะยังผ่าน — แก้เป็น capture error ทั้งสอง
    goroutine + assert `state == "passed"` ตรง ๆ
- **Known debt (บันทึกไว้ ยังไม่แก้ในรอบนี้)**:
  - "`first_passed_at` write-once" เป็น invariant ที่มีอยู่แค่ใน infra (`decision.State.IsPassed()`
    + `COALESCE` ใน SQL) — `domain.LessonProgress` และ `app.ProgressEntry` ไม่ได้ model concept นี้
    ไว้เลย คำตอบที่ตรงกับความจริงของ "domain เป็นเจ้าของกฎทุกข้อไหม": forward-only — ใช่;
    first-passed-at write-once — ไม่ใช่ (infra เป็นเจ้าของ); fresh-row path — bypass domain
    transitions ไปเลยตรง ๆ (ไม่ผ่าน `Unlock`/`MarkPassed`)
  - `Repository.Transition` ยังรับ raw string (`topicSlug, conceptSlug string`) ไม่ใช่ value object —
    caller ใหม่ในอนาคตที่เรียก `Transition` ตรง ๆ (ข้าม `Service.SetProgress`) จะข้าม slug validation
    ไปเลยโดยไม่รู้ตัว — `TestRepositoryTransition_CanonicalSlugsReturned` สาธิตพฤติกรรมนี้อยู่แล้ว
    (เรียก `repo.Transition` ตรง ๆ ด้วย `"DDIA"`/`"B-Trees"` แล้วผ่าน เพราะ shape validation อยู่ที่
    `Service` เท่านั้น) — การแก้แบบเต็ม (thread VO ผ่าน `Transition`) เป็นงานใหญ่กว่าที่ตัดสินใจไม่ทำรอบนี้
  - `concepts` unique key คือ `(chapter_id, slug)` ไม่ใช่ต่อ topic ที่ระดับ schema จริง — **แก้ไขหลัง
    UX-5's code review**: สรุปเดิมของ note นี้ผิด บอกว่าสองบทใน topic เดียวกันที่ใช้ concept slug
    ซ้ำกันจะทำให้ UX-4's PUT progress "เขียนแถวผิดเงียบ ๆ" ได้ — ที่จริงไปไม่ถึงจุดนั้นเลย เพราะ
    `domain.NewTopic` (`internal/curriculum/domain/topic.go:64-70`, เรียกจาก curriculum's read path
    เองใน `internal/curriculum/infra/repository.go`) บังคับ concept slug ไม่ให้ซ้ำกัน **ทั้ง topic**
    ตอน assemble แล้ว — ถ้ามีข้อมูลซ้ำแบบนี้จริงใน DB, **`GET /api/v1/curriculum` จะ error (500) ทั้ง
    topic นั้นทันที** ก่อนที่ PUT progress จะมีโอกาสเจอแถวซ้ำด้วยซ้ำ — failure mode ที่แท้จริงคือ
    "curriculum tree พังทั้งก้อนแบบเห็นชัด" ไม่ใช่ "เขียนข้อมูลผิดแบบเงียบ ๆ" `lockLesson`'s
    Query-not-QueryRow ยังเป็น defense-in-depth ที่ดีอยู่ แต่กันไว้สำหรับสถานการณ์ที่ curriculum's
    domain layer เองก็ไม่ปล่อยให้ถึง live system อยู่แล้ว — schema/unique constraint จริงยังไม่ได้แก้
    (debt เดิมยังอยู่ แค่ความเสี่ยงต่ำกว่าที่ note เดิมประเมินไว้มาก)
  - lock-wait timeout / deadlock (MySQL error 1205/1213) จาก `FOR UPDATE OF l` ยังไม่ map เป็น
    status ที่วินิจฉัยได้ (เช่น 409/503) — ตอนนี้ตกไปที่ 500 ทั่วไปเหมือน error อื่น ๆ; เคสที่จะเจอจริง
    คือรัน `import-lessons`/`import-curriculum` พร้อม API รับ traffic (lock wait default 50s ใกล้
    `writeTimeout = 60s` ของ server พอสมควร)
- Status: `merged, PR #45`

## UX-5 — Reader loop: breadcrumb + Finish + Next `[go-implementer]`

- **Scope**: `web/app/(app)/lesson/page.tsx` ได้ breadcrumb (`Track › Chapter ›
  Concept`), mark-in-progress on load, Finish button ที่ยิง `PUT
  /api/v1/progress/{topic}/{concept}` (จาก UX-4) จริง, และ Next ชี้ concept ถัดไปที่
  มี lesson ในอ่าน order เดียวกัน; `web/lib/api.ts` เพิ่ม `getProgress`/`setProgress`;
  `web/lib/curriculum.ts` เพิ่ม pure function สองตัว (`locateLessonBreadcrumb`,
  `findNextLesson`) และย้าย `pinFocusFirst` มาจาก `learn/page.tsx` (เดิม unexported,
  ทดสอบไม่ได้) — ทั้งสามใช้ `flattenTrack` ตัวเดียวกันภายในไฟล์
- **Setup เพิ่ม**: ติดตั้ง **vitest** (`web/vitest.config.ts`, `npm test` script,
  wired เข้า `.github/workflows/ci.yml` ก่อน `npm run build`) — `web/` ไม่เคยมี test
  runner มาก่อน สอง code review ก่อนหน้าเตือนเรื่องนี้ไว้แล้ว
- **การตัดสินใจหลัก**:
  - **Soft-guide ไม่ gate**: ไม่มี lesson ไหนถูกล็อกหรือแสดงเป็น unavailable เพราะบทก่อนหน้า
    ยังไม่จบ — Next แค่ "แนะนำ" ไม่ "บังคับ"
  - **"Finished" = รีวิวครบทุกข้อ ไม่ใช่ "ผ่านทุกข้อ"**: `allRated` เช็คว่าทุก
    `recall_check.position` มี rating (`pass` หรือ `fail`) อยู่ใน state — บทที่ตอบ
    `Not yet` ทั้งหมดก็ finish ได้ปกติ, บทที่ไม่มี recall check เลย (`totalChecks===0`)
    ก็ finish ได้ทันที (`0 === 0`)
  - **หนึ่ง round trip พอสำหรับทั้ง mark-visit และเช็คว่า finish แล้วหรือยัง**: PUT
    `in_progress` ทุกครั้งที่ lesson โหลดสำเร็จ (ไม่ใช่ก่อนโหลด — URL พังต้องยังเป็น 404
    สะอาด ไม่ใช่ side effect เงียบ ๆ) แล้วอ่าน **response** — forward-only clamp ของ
    UX-4 ทำให้ concept ที่ passed แล้วตอบกลับ `{"state":"passed"}` เสมอ ไม่ต้องยิง
    `GET /api/v1/progress` แยกอีกเส้น
  - **`aria-disabled` ไม่ใช่ `disabled` (ซ้ำกับ UX-2)**: `disabled` attribute ทำให้
    browser blur ปุ่มที่เพิ่งกดทันที (บั๊กที่ UX-2 เจอ) — Finish ใช้ `aria-disabled` +
    hint นับความคืบหน้าจริง (`Rate all 4 checks to finish (2/4 rated)`) และกด "activate"
    ตอนยังไม่ครบจะย้าย focus ไปที่ recall check ใบแรกที่ยังไม่ rate (ผ่าน ref map +
    `RecallCheckCard` เป็น `forwardRef`) แทนที่จะเฉย ๆ
  - **รายงาน error จริง ไม่ fake success**: Finish มี state `idle/saving/saved/error`
    ชัดเจน — ถ้า PUT fail โชว์ error banner + ปุ่ม Retry, ไม่เปลี่ยนเป็น "Lesson finished"
    จนกว่า PUT จะสำเร็จจริง
  - **curriculum tree กับ lesson แยก failure กัน (ซ้ำกับบทเรียนจาก UX-2)**: หน้านี้เรียก
    `getCurriculum()` เพิ่มจาก `getLesson()` แต่เป็นคนละ `useEffect` — breadcrumb/Next
    fallback เป็น `Learn › {lesson.title_en}` เฉย ๆ ถ้า tree โหลดไม่ทัน/พัง, ตัว lesson
    เองยัง render ปกติเสมอ
  - **Traversal ข้าม chapter/topic แต่ไม่ข้าม track**: `findNextLesson` flatten
    topics→chapters→concepts ของ **track เดียวกับ concept ปัจจุบัน** ตาม position ที่
    backend sort มาให้แล้ว (ไม่ re-sort ฝั่ง frontend) แล้ว `.find()` ตัวแรกที่
    `has_lesson`, ข้าม concept ที่ไม่มี lesson ไปเรื่อย ๆ รวมถึงข้าม topic boundary —
    คืน `null` ทั้งกรณี "หมด track แล้ว" และ "position ปัจจุบันไม่อยู่ใน tree เลย" (ทั้งคู่
    แปลว่า "ไม่มีอะไรให้ชี้ต่อ" เหมือนกันจากมุมมอง UI)
  - **breadcrumb เลือก wrap ไม่ truncate ที่ 375px**: หัวข้อ chapter ภาษาไทยยาว ๆ ถ้าตัด
    กลางคำจะเสียความหมาย, หน้า reader มีที่ว่างแนวตั้งเหลือพอให้ breadcrumb ขึ้นบรรทัดใหม่ได้
    โดยไม่กระทบ layout อื่น
- **Review focus**:
  - ทำไม mark-in-progress ต้องยิง **หลัง** `getLesson()` สำเร็จเท่านั้น ไม่ใช่ก่อนหน้านั้น?
  - ทำไมอ่าน state จาก **response ของ PUT in_progress** พอ ไม่ต้องมี `GET
    /api/v1/progress` แยกอีกเส้น?
  - ทำไม "Finished" ต้องนับจาก "รีวิวครบทุกข้อ" ไม่ใช่ "ผ่านทุกข้อ"?
  - ทำไม `findNextLesson` ต้องคืน discriminated union (`"next" | "end-of-track" | "not-found"`)
    แทนที่จะคืน `null` เฉย ๆ สำหรับทั้งสองกรณีที่ไม่มี next?
- **Known debt จาก code review (บันทึกไว้ ตั้งใจไม่แก้ใน ticket นี้)**:
  - **S12 — `ratings` (recall check pass/fail ต่อข้อ) เป็น page-local state ไม่เคย persist**:
    rate 3 จาก 4 ข้อแล้วออกจากหน้าไปโดยไม่กด Finish หายเงียบ ๆ — วันนี้ไม่มีผลอะไร (ยังไม่มีที่ไหน
    เก็บ per-check rating ลง DB เลย, มีแค่ aggregate "passed" ทั้ง lesson ที่ Finish เขียน) แต่พอ
    SM-2 / recall-attempt write ลงจริง (สัปดาห์ 4, Drill tickets T16–T21) จุดนี้จะกลายเป็น data
    loss ของจริงทันที — ต้อง revisit ตอนนั้น ไม่ใช่ ticket นี้
  - **S15 — `Button`'s `"danger"` variant ไม่มีใครเรียกใช้เลยทั้งแอป, `text-danger` (rose-600) บน
    cream วัดได้ ~4.39:1 ซึ่ง**จะ fail AA**ทันทีที่มีคนเอาไปใช้จริง (14px ต้องการ 4.5:1) — บันทึกไว้
    เป็น debt เฉย ๆ ไม่เพิ่ม caller ปลอมขึ้นมาเพื่อ "justify" การมีอยู่ของ variant นี้
- Status: `merged, PR #46`

## UX-6 — Today page + IA switch `[go-implementer]`

- **Scope**: หน้าใหม่ `/today` ตอบ "วันนี้อ่านอะไรต่อ" — **Next up** card (track ·
  chapter, ชื่อ concept, `~N min`, ปุ่ม **Continue**/**Start** ตาม state) +
  **Other tracks** section (ทุก track ที่มี lesson อย่างน้อย 1 บท โชว์
  `n/N done`, ปุ่มสลับ focus, ลิงก์ `/learn?track=`); IA ตัดเหลือ **Today ·
  Learn** ใน nav (`web/lib/nav.ts`) เพราะอีก 5 หน้า (Dashboard/Drill/DSA/
  Test/Tickets) เป็น `ComingSoon` placeholder ล้วน — **ไฟล์ page ยังอยู่บน
  disk เหมือนเดิม แค่ไม่มีลิงก์ใน nav ชี้ไปหา** (รอ phase 4/5 ค่อยกลับมาใส่
  ใน nav ใหม่ตอนมีของจริง); `/` และ `/token` (หลัง login) ชี้ไป `/today`
  แทน `/learn` เดิม — ส่วน `/read` ยังคงชี้ไป `/learn` เหมือนเดิม (มันคือ
  URL เก่าของหน้า Learn เอง ไม่ใช่ "ทางเข้าแอป" ที่ ticket นี้ต้องย้าย)
- **`web/lib/api.ts`**: เพิ่ม `getProgress()` กลับมา (ถูกลบเป็น dead code
  ตอน UX-5 เพราะตอนนั้นไม่มีใครเรียก) — **bug ที่จับได้จากการรัน Playwright
  จริงกับ live stack**: `GET /api/v1/progress` ห่อ array ไว้ใน
  `{"concepts": [...]}` ไม่ใช่ array เปล่า ๆ ตามที่ผมสมมุติไว้ตอนแรกจาก
  UX-4's ticket text เฉย ๆ — ถ้าไม่ทดสอบกับ live stack จริง โค้ดจะ compile
  ผ่านทุกอย่าง (TypeScript ไม่เช็ค runtime shape) แต่ `indexProgress` จะ
  throw ตอน `for...of` วน object ที่ไม่ใช่ array จริง กลายเป็น curriculum
  โหลดสำเร็จแต่ทั้งหน้า error เพราะ progress พังแทน — ตรงข้ามกับ requirement
  "ต้อง degrade เฉพาะส่วนที่พัง" ของ ticket นี้เป๊ะ ๆ. `getProgress()` เลย
  unwrap `res.concepts` แล้ว catch เหมือน `getFocusTrack()` เดิม (non-401
  → คืน `[]`, 401 → rethrow ให้ redirect ไป `/token`)
- **`web/lib/curriculum.ts`**: `pinFocusFirst` เปลี่ยนเป็น generic
  `<T extends {track: string}>` (เดิมรับแค่ `TrackStats[]`) เพื่อใช้ซ้ำกับ
  `Track[]` ใน `pickNextUp` โดยไม่ต้องมี copy ที่สอง; ย้าย sort-by-display-
  order ออกจาก `learn/page.tsx`'s local `trackSortIndex` ไปเป็น
  `sortTracksByDisplayOrder` (generic เหมือนกัน) ใน `web/lib/trackMeta.ts`
  แล้วให้ทั้ง `/learn` และ `pickNextUp` เรียกตัวเดียวกัน — เหตุผลเดียวกับที่
  ticket บอกให้ reuse `flattenTrack`: มี logic เดียวกันสองที่วันหนึ่งมันจะ
  drift
- **ฟังก์ชันหลัก (pure, `web/lib/curriculum.ts`)**:
  `pickNextUp(tracks: Track[], progressByKey: ProgressByKey, focusTrack:
  string | null): NextUpResult` โดย
  `ProgressByKey = Record<string, ProgressState>` (คีย์จาก
  `progressKey(topic, concept)` เท่านั้น ห้าม hardcode separator เอง) และ
  ```
  type NextUpResult =
    | { kind: "next"; track; topic; concept; chapterTitle; title;
        estMinutes: number | null; state: "not_started" | "in_progress" }
    | { kind: "all-done" }
    | { kind: "no-lessons" };
  ```
  แยก `"all-done"` (ทุก concept ที่มี lesson ทั้งแอป passed หมด — ยินดีด้วย
  จบจริง) ออกจาก `"no-lessons"` (ไม่มี track ไหนมี lesson เลยสักบท — ยังไม่มี
  อะไรให้อ่าน) เพราะสอง state นี้ว่างเหมือนกันแต่ความหมายตรงข้ามกัน —
  ถ้ายุบเป็น `null` เดียวกัน UI จะพูดผิดว่า "เก่งมาก อ่านจบหมดแล้ว" ทั้งที่
  ยังไม่มีเนื้อหาเลยสักบท (ใช้หลักการเดียวกับ UX-5's `findNextLesson` ที่
  แยก `"end-of-track"` ออกจาก `"not-found"`)
  - **ลำดับ track**: `pinFocusFirst(sortTracksByDisplayOrder(tracks),
    focusTrack)` — sort ตาม `TRACK_DISPLAY_ORDER` ก่อน (เป็นแค่ sort key,
    track ที่ backend ส่งมาแต่ไม่อยู่ใน list จะถูกต่อท้าย ไม่หาย) แล้วค่อยดัน
    focus track ไปหน้าสุด ถ้ามีและยังอยู่ในลิสต์
  - **ลำดับ concept ในแต่ละ track**: เดินตาม `flattenTrack`'s reading order
    (topic → chapter → concept ตาม position ที่ backend sort มาแล้ว), ข้าม
    concept ที่ `has_lesson=false`, คืนตัวแรกที่ state **ไม่ใช่** `passed`
    (`in_progress` หรือไม่มี row เลยก็ได้ทั้งคู่) — บทที่อ่านค้างมาก่อนเลย
    โผล่มาก่อนบทที่ยังไม่แตะเลยโดยอัตโนมัติ ไม่ต้องมี priority พิเศษแยก
    ต่างหาก
  - **Soft-guide**: ไม่มี state `"locked"` ปรากฏใน result เลย (`domain.Gate`
    ยังไม่ถูก wire — ตรงตาม constraint ของ ticket นี้)
- **`FocusToggleButton` component ใหม่** (`web/components/
  FocusToggleButton.tsx`): ดึงปุ่ม toggle focus (aria-pressed/aria-disabled
  + optimistic label) ออกจาก `TrackCard.tsx` มาเป็น component แชร์ได้ ใช้ทั้ง
  ใน `TrackCard` (หน้า Learn) และแถว "Other tracks" ของ `/today` — ป้องกัน
  behavior (โดยเฉพาะ `aria-disabled` แทน `disabled`) แยกกันสองที่แล้ว drift
  ไปคนละแบบ
- **Failure decoupling (R1 จาก UX-2 อีกครั้ง)**: `/today` ยิงสามคำขอพร้อมกัน
  `Promise.all([getCurriculum(), getProgress(), getFocusTrack()])` —
  `getProgress`/`getFocusTrack` catch ทุก error ที่ไม่ใช่ 401 ไว้ในตัวเอง
  แล้วคืนค่า default (`[]` / `{track: null}`) ดังนั้น `Promise.all` reject
  ได้จากสาเหตุเดียวคือ `getCurriculum()` พังจริง — progress/focus-track พัง
  ไม่มีทางทำให้ทั้งหน้าเป็น error state ได้เลย
- **ไม่มี progress bar**: "Other tracks" โชว์ตัวเลข `n/N done` ตรง ๆ
  (constraint เดิมจาก UX-2: บาร์ที่ตัวเศษ=ตัวส่วนเสมอจะสื่อว่า "จบแล้ว" ทั้ง
  ที่ยังไม่ได้ progress-track จริง — บาร์มาทีหลังใน UX-7)
- **Review focus**:
  - ทำไม `pickNextUp` ต้องคืน `"all-done"` กับ `"no-lessons"` แยกกันสองแบบ
    แทนที่จะคืน `null`/`undefined` ตัวเดียวสำหรับทั้งสองกรณี?
  - ทำไม `Promise.all` ของ `/today` ถึงไม่มีทาง reject จากแค่ progress หรือ
    focus-track พังเพียงอย่างเดียว?
  - `pinFocusFirst` เปลี่ยนจากรับ `TrackStats[]` เป็น generic
    `<T extends {track: string}>` แล้ว behavior ของที่เรียกใช้เดิม (`/learn`)
    เปลี่ยนไปหรือเปล่า เพราะอะไร?
  - ทำไม bug ของ `getProgress()` (ลืม unwrap `{concepts: [...]}`) ถึงรอดจาก
    `npx tsc --noEmit` และ `npm test` ผ่านหมด แต่จับได้จากการรัน Playwright
    กับ live stack เท่านั้น?
- **Mutation-testing** (ทำจริงตามที่ ticket ขอ แก้โค้ดชั่วคราว รัน `npm test`
  แล้ว revert ทุกครั้ง):
  - ลบ `pinFocusFirst(...)` ออก (เท่ากับเมิน focus track) → เจ๊ง 1 เทสต์:
    "returns the focus track's next concept ahead of an earlier-ranked track"
  - เปลี่ยน `if (state === "passed")` เป็น `if (state === "passed" ||
    state === "in_progress")` (นับ in_progress ว่าจบแล้ว) → เจ๊ง 1 เทสต์:
    "returns an in_progress concept ahead of a later not_started one in the
    same track"
  - ลบ `if (!item.hasLesson) continue;` ออก (ไม่กรอง has-lesson) → เจ๊ง 2
    เทสต์: "skips concepts without a lesson" และ "returns no-lessons when no
    track has any lesson at all"
  - fixture ของทุกเทสต์ใช้อย่างน้อย 2 track เสมอ (ตาม lesson จาก UX-5 ที่ตอน
    แรก fixture track เดียวทำให้ 2 ใน 3 mutation รอด) — ยืนยันด้วยผลข้างต้น
    ว่าแต่ละ mutation โดนจับตรงจุดจริง ไม่ใช่บังเอิญ
- **ทดสอบกับ live stack จริง** (`docker compose up -d --build`, bearer token
  จาก `.env` local, curl ตั้งค่าล่วงหน้า): focus=`ddd` + concept แรก
  `in_progress` → Next up โชว์ "Ubiquitous Language" ปุ่ม **Continue** ลิงก์
  `/lesson?topic=domain-driven-design&concept=ubiquitous-language` ถูกต้อง;
  pass ทั้ง 4 concept ของ `ddd` แล้ว Next up ย้ายไป track อื่นตาม
  `TRACK_DISPLAY_ORDER` (ข้าม distsys/aws/go/dsa ที่ยังไม่มี lesson เลย
  ไปเจอ ddia); mock `GET /api/v1/progress` เป็น 500 → หน้ายัง render Next up
  ปกติ (degrade เป็น "ทุกอย่างยังไม่เริ่ม" แทนที่จะ error ทั้งหน้า), ไม่มี
  error banner ปลอม; กด "Set focus" จากหน้า `/today` แล้ว
  `curl GET /api/v1/prefs/focus-track` ยืนยันค่าถูก persist จริง; nav บน
  `/today`/`/learn` highlight ถูกหน้า, เหลือแค่ 2 item; 375px:
  `scrollWidth === clientWidth === 375`, console error = 0 ทุกกรณียกเว้น
  เคส mock-500 ที่มี browser's built-in "Failed to load resource: 500" log
  เดียว (เป็น network log อัตโนมัติของ browser เอง ไม่ใช่โค้ดแอป log เพิ่ม —
  ไม่ใช่ error storm)
- **จงใจไม่ทำในรอบนี้**: ไม่ได้ทดสอบ `"all-done"` กับ live DB จริง (ต้อง pass
  lesson ทั้ง 61+53 บทของ ddia/ai-systems ผ่าน curl ซึ่งไม่คุ้มเวลา) —
  ครอบคลุมด้วย unit test แทน (`"returns all-done when every lesson in every
  track is passed"`); ไม่ได้ reset ddd's progress rows กลับเป็นค่าว่างหลัง
  ทดสอบ (ไม่มี DELETE endpoint ให้ใช้) เหลือ `ddd` ติด `passed` ทั้ง 4
  concept ใน local dev DB — ไม่กระทบ prod เพราะเป็น local DB ของเครื่องพัฒนา
  เอง, focus track ถูก reset กลับเป็น `null` แล้ว
- Status: `implemented, PR pending`

## UX-7 — Progress page + chapter/track indicators

- Progress indicator ที่เห็นคร่าว ๆ: ของ chapter (แต่ละบทเรียนเป็นไหนแล้ว) + ของ track (overview),
  โชว์บน layout ทั่วแอป (progress bar, % เลยน้อย ๆ ที่หน้า reader หรือ learn)
- Status: `ยังไม่เริ่ม`
