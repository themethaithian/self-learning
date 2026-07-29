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
- **จงใจไม่ทำในรอบแรก**: ไม่ได้ทดสอบ `"all-done"` กับ live DB จริง (ต้อง pass
  lesson ทั้ง 61+53 บทของ ddia/ai-systems ผ่าน curl ซึ่งไม่คุ้มเวลา) —
  ครอบคลุมด้วย unit test แทน (`"returns all-done when every lesson in every
  track is passed"`)

### รอบ code-reviewer (REQUEST_CHANGES → แก้ครบ)

- **R1 — progress พังแล้วหน้าโกหกตัวเลข**: `getProgress()` เดิม catch แล้ว
  คืน `[]` เงียบ ๆ — `/today` เอา `[]` ไปใช้ราวกับเป็นข้อมูลจริง ทำให้ track
  ที่อ่านไปเยอะแล้วโชว์ `0/61 done` เมื่อ progress fetch ล้มเหลว (ผู้ใช้แยก
  "ศูนย์จริง" กับ "โหลดไม่สำเร็จ" ไม่ออก — เป็นบั๊กเดียวกับที่ทั้ง epic นี้
  ปฏิเสธไม่ยอมทำ progress bar ที่เต็มเสมอ) แก้โดยเปลี่ยน `getProgress()` คืน
  discriminated result `{ok: true, entries} | {ok: false}` (`web/lib/api.ts`)
  `/today` เก็บ `progressAvailable` แยกจาก `progressByKey`: เมื่อ `ok: false`
  เลขที่ยังไม่รู้ค่าจริงจะโชว์ `"{total} lessons"` แทน `n/N done` (ไม่โกหกว่า
  "ศูนย์"), และมี `role="alert"` banner บอกตรง ๆ ว่าโหลด progress ไม่สำเร็จ
  (เหมือน banner ของ focus-track error ที่มีอยู่แล้ว) — `pickNextUp` เอง
  ไม่ต้องแก้ (รับ empty map เป็น "ทุกอย่าง not_started" ต่อไปตามเดิม)
- **R2 — สาม exported function ใหม่ไม่มี mutation-resistant test**: เพิ่ม
  test 3 ชุด (`web/lib/curriculum.test.ts`, `web/lib/trackMeta.test.ts`)
  แล้วพัง mutation จริงยืนยัน:
  - `indexProgress` คืน `{}` เดิมเสมอ → พัง 2 เทสต์ ("indexes entries by
    topic and concept, preserving state", "keys by topic then concept, not
    the reverse")
  - `indexProgress` สลับ argument เป็น `progressKey(entry.concept,
    entry.topic)` → พังเทสต์เดิม 2 ชุดเดียวกัน (ทั้งคู่จับได้)
  - `computeTrackProgress` เปลี่ยน `=== "passed"` เป็น `!== "not_started"`
    (นับ `in_progress` เป็น done) → พัง 1 เทสต์ ("does not count an
    in_progress lesson as done")
  - `sortTracksByDisplayOrder` เปลี่ยน `Number.MAX_SAFE_INTEGER` เป็น `-1`
    (track ที่ไม่รู้จักไปอยู่หัวแถวแทนท้ายแถว) → พัง 1 เทสต์ ("appends a
    track absent from TRACK_DISPLAY_ORDER to the end, not the start")
- **R3 — `getProgress()` ยังคืน `undefined` ได้จริงถ้า response ผิดรูป**:
  `res.concepts` จาก payload ที่ไม่มี field นี้จะเป็น `undefined` แล้ว
  `indexProgress(undefined)` throw กลาง `load()`'s try ทำให้ curriculum ที่
  โหลดสำเร็จกลายเป็น error ทั้งหน้าไปด้วย แก้ด้วย
  `Array.isArray(res.concepts) ? res.concepts : []` ก่อนส่งต่อ — `{ok:
  true, entries}` การันตีว่า `entries` เป็น array เสมอ
- **S1** — "Other tracks" กรอง track ที่เป็น Next up ออก (`nextUp.kind ===
  "next" && stats.track === nextUp.track` ถูกตัดจากลิสต์) กัน track เดียวกัน
  โผล่ซ้ำสองที่พร้อมกัน
- **S2** — `/today`'s "Other tracks" ใช้ `pinFocusFirst` ต่อจาก
  `sortTracksByDisplayOrder` เหมือน `/learn` แล้ว (เดิมขาด `pinFocusFirst`
  ทำให้สอง page เรียง "focus" คนละแบบ)
- **S3** — `pickNextUp`'s return object ใช้ `state,` ตรง ๆ แทน ternary
  `state === "in_progress" ? "in_progress" : "not_started"` — TypeScript
  narrow `state` เหลือ `"not_started" | "in_progress"` ไปแล้วจาก `if (state
  === "passed") continue;` บรรทัดก่อนหน้า ternary เดิมสื่อผิดว่ามี state ที่
  สี่อยู่
- **S4** — ดึง optimistic update/revert/pending/error state ทั้งก้อนออกเป็น
  `web/lib/useFocusTrack.ts` (คืน `{focusTrack, setFocusTrackValue,
  pendingTrack, error, setFocus, dismiss}`) ใช้ทั้ง `/learn` และ `/today` —
  ก่อนหน้านี้ extraction หยุดแค่ปุ่ม (`FocusToggleButton`) แต่ตัว state
  machine ~35 บรรทัดยังก็อปสองที่เหมือนกันเป๊ะ
- **S5/S6** — ลบ comment ที่จะ rot เร็ว: `FocusToggleButton.tsx` ตัด call-site
  inventory ("currently TrackCard and /today") ออก เก็บแค่ WHY ของ
  `aria-disabled`; `curriculum.ts`'s `pickNextUp` comment ตัดประโยคท้ายที่พูดซ้ำ
  กับ `NextLessonResult`'s union ที่นิยามอยู่ในไฟล์เดียวกันแล้ว;
  `TrackCard.tsx` แก้ comment ที่อ้างว่า "reading progress isn't tracked
  (UX-4)" ทั้งที่ UX-4 ทำเสร็จแล้วและ ticket นี้เองก็ใช้ progress อยู่ —
  เปลี่ยนเป็นอธิบายว่าการ์ดนี้นับ availability ไม่ใช่ completion
- **S8** — `docs/roadmap.md` เปลี่ยน `[~]` เป็น `[ ]` + "(PR pending)"
  เพราะ GitHub เรอนเดอร์ `[~]` เป็น raw text ไม่ใช่ checkbox (รีวิวจากมือถือ
  จะเห็นเป็นตัวหนังสือเปล่า ๆ)
- **S9** — `clearToken()` เพิ่ม `typeof window` guard + try/catch: ถ้า
  localStorage ถูกบล็อก (extension/private mode — gotcha ที่โปรเจกต์นี้เคยเจอ
  จริงกับ `getToken`) `removeItem` throw จะทำให้ `apiFetch`'s 401 handler
  ไม่ทัน `throw new UnauthorizedError()` เลย กลาย เป็น error ธรรมดาที่
  `getProgress`/`getFocusTrack` กลืนทิ้งเงียบ ๆ แทนที่จะ redirect ไป `/token`
- **Take-if-cheap ที่ทำ**: S7 (`progressKey`'s `::` separator) — เพิ่มเทสต์
  `progressKey("ab","c") !== progressKey("a","bc")` ล็อก invariant กันการต่อ
  string ตรง ๆ โดยไม่มี separator
- **Re-verify**: `npm test` 35/35 ผ่าน (เพิ่มจาก 27 → 35: +3 `indexProgress`
  +1 `computeTrackProgress` +1 `progressKey` +3 `sortTracksByDisplayOrder`
  ใน `web/lib/trackMeta.test.ts` ไฟล์ใหม่), `tsc --noEmit`/`eslint`/`npm run
  build` สะอาด, `go vet`/`go test` cached ผ่าน (ยืนยันด้วย `git diff
  --name-only develop...` ไม่มีไฟล์ Go), mutation ทั้ง 7 จุด (4 เดิม + 3
  ใหม่จาก R2) พังตามที่ระบุไว้ข้างบนทุกจุดแล้ว revert กลับก่อน commit
- **Live-stack re-verify**: `docker compose up -d --build web` แล้ว
  Playwright จริง — normal load: Next up = ddia's "Latency and Percentiles"
  (`in_progress` จริงจาก DB) ปุ่ม **Continue**, Other tracks โชว์ DDD
  `4/4 done` + AI & LLM Systems `0/53 done` (ไม่นับ 2 concept ที่
  `in_progress` เป็น done) และ **ไม่โชว์ ddia ซ้ำ** (เพราะมันคือ Next up
  track อยู่แล้ว, S1); mock `GET /api/v1/progress` → 500: banner
  "Could not load your reading progress..." ขึ้นจริง, Next up ตกกลับไป
  `ddd`'s "Ubiquitous Language" ปุ่ม **Start** (เพราะไม่รู้ progress จริง),
  Other tracks โชว์ `"61 lessons"` / `"53 lessons"` **ไม่มี `0/N` ปลอม**;
  กด "Set focus" จาก `/today` แล้ว `curl GET /api/v1/prefs/focus-track` ยืนยัน
  persist จริง; `/learn` เองก็ยังใช้งานได้ปกติหลัง refactor เป็น
  `useFocusTrack()` (focus toggle persist ตรวจแล้ว); 375px:
  `scrollWidth === clientWidth === 375`, console error = 0 ทุกเคสยกเว้น
  mock-500 ที่มี browser's built-in "Failed to load resource: 500" log
  เดียวเหมือนรอบแรก
- **จงใจไม่ทำในรอบนี้เช่นกัน**: ไม่ได้ reset `ddd`/ddia/ai-systems's progress
  rows กลับเป็นค่าว่างหลังทดสอบ (ไม่มี DELETE endpoint ให้ใช้) — ไม่กระทบ
  prod เพราะเป็น local dev DB, focus track ถูก reset กลับเป็น `null` แล้ว
  ทุกครั้งหลังทดสอบ

### รอบ code-reviewer ที่ 3 (REQUEST_CHANGES → แก้ครบ)

- **A — R3's fix เปิดช่องให้ R1 กลับมาทางประตูหลัง**: `Array.isArray(res.concepts)
  ? res.concepts : []` เดิม คืน `{ok: true, entries: []}` เมื่อ response 200
  แต่ shape ผิด (เช่น key เปลี่ยนชื่อเป็น `entries`, หรือ error payload คนละ
  รูปแบบ) — เท่ากับ "โหลดสำเร็จ ประวัติว่างเปล่า" ปลอม ทำให้ banner หาย และ
  `/today` มั่นใจโชว์ `0/N done` + **Start** ทั้งที่ progress ไม่ได้โหลดจริง
  (บั๊กเดียวกับ R1 แค่ trigger เปลี่ยนจาก 500 เป็น shape drift — และ PR #48
  เพิ่ง handle_errors บน Caddy ทำให้พื้นที่ "200 ที่ shape ไม่ตรง" กว้างขึ้น)
  แก้เป็น `Array.isArray(res?.concepts) ? {kind:"ok", entries: res.concepts}
  : {kind:"error"}` — shape ผิดตอนนี้จัดเป็น error เสมอ ไม่ใช่ ok เปล่า ๆ;
  comment เดิมที่บอกว่า "รับประกันไม่มีทาง entries เป็น undefined" ก็แก้เป็น
  พูดถึงสิ่งที่ต้องกันจริง (false ok) แทน
- **B — `useFocusTrack` กลืน `UnauthorizedError` แล้วบอกให้ retry ทั้งที่ retry
  ไม่มีทางสำเร็จ**: token หมดอายุ → `apiFetch` เจอ 401 → `clearToken()` →
  throw `UnauthorizedError` → hook เดิมจับด้วย bare `catch` เงียบ ๆ → banner
  ขึ้น "Could not save focus track — please try again" → ไม่ redirect, หน้า
  ยังดูเหมือน login อยู่ → กด retry ยิง request **ไม่มี Authorization header
  เลย** → 401 วนซ้ำตลอดไป (ขัด policy ของ `api.ts` เองที่บอกว่า "expired
  token คือความล้มเหลวเดียวที่ caller ต้องเห็น" ซึ่ง `lesson/page.tsx` ทำตาม
  แต่ hook นี้ทำลาย policy ให้พังพร้อมกันสองหน้า) แก้โดยให้ hook เช็ค
  `err instanceof UnauthorizedError` แล้ว `router.replace("/token")` เอง
  (ใช้ `useRouter` จาก `next/navigation` ภายใน hook) แทนที่จะ rethrow ให้ทั้ง
  สองหน้าไปดักเอง — จุดเดียวแก้ครบทั้ง `/today` และ `/learn`
- **S10 (promoted) — โค้ดใหม่ของรอบนี้ทั้งหมดไม่มี test เลย, 11 mutation รอด
  หมด**: ปิดช่องที่ปิดได้ด้วย unit test (ไม่ต้อง jsdom):
  - เพิ่ม `web/lib/api.test.ts` (ใช้ `vi.stubGlobal("fetch", ...)` บน
    `environment: "node"` เดิม ไม่ต้องพึ่ง dependency ใหม่) คุม 4 mutation:
    ลบ `Array.isArray` guard, catch คืน `{kind:"ok", entries:[]}` แทน
    `{kind:"error"}`, `getProgress`/`getFocusTrack` กลืน `UnauthorizedError`
  - ดึง dedupe filter (ตัด track ที่เป็น Next Up ออกจาก "Other tracks")
    ออกจาก JSX ไปเป็น pure function `buildOtherTracks(tracks, progressByKey,
    focusTrack, nextUp)` ใน `curriculum.ts` แล้วมี unit test ปิด mutation
    "ลบ dedupe filter" ได้ตรง ๆ (จาก 11 mutation ที่รอด ปิดได้ 5:
    Array.isArray guard, catch คืนผิด, getProgress/getFocusTrack กลืน
    Unauthorized, dedupe filter หาย)
  - **ยังไม่ปิดในรอบนี้ (ตั้งใจ)**: `progressAvailable` force-true, ลบ
    banner ทิ้ง, โชว์ `n/N done` ตอน progress ไม่มา, hook's optimistic
    revert หาย, `pendingTrack` หาย (Saving.../busy พัง), hook ไม่โชว์
    error เลย — ทั้ง 6 จุดนี้เป็น React component/hook state ต้องมี jsdom +
    `.tsx` test include ซึ่งยังไม่ได้ตั้งในรอบนี้ (`vitest.config.ts` ยัง
    `environment: "node"`, `include: ["**/*.test.ts"]` เท่านั้น) — บันทึกไว้
    เป็น debt ที่ยอมรับได้ ตรวจแทนด้วย Playwright ทุกครั้งที่แก้หน้านี้
- **S12/S13 — ข้อความ banner เอง**: เดิมบอกว่าตัวเลข "below" กับ label
  "above" ทั้งที่ banner เป็น first child ของ `space-y-8` stack ทำให้ทั้งคู่
  อยู่ "below" จริง ๆ (ไม่มีอะไรอยู่เหนือ banner นอกจาก `<h1>Today</h1>`) และ
  เดิมไม่ได้เตือนเรื่อง **คำแนะนำเอง** เลย — ทั้งที่ progress ว่างทำให้
  `pickNextUp` ชี้ไปที่ concept แรกของ track เสมอ (คนอ่าน DDIA ไป 40 บทแล้ว
  จะเจอ "Next up: บทที่ 1") แก้ข้อความเป็น "the pick below, its
  Continue/Start label, and every count under Other tracks may be wrong
  until this loads" ครอบคลุมทั้งสามจุด โดยไม่ซ่อนการ์ด Next Up ทิ้ง
- **S11** — ย้าย concurrency guard (`if (pendingTrack !== null) return`)
  เข้าไปอยู่ใน `useFocusTrack.setFocus` เอง แทนที่จะพึ่ง `busy` prop ที่หน้า
  คำนวณแล้วส่งให้ `FocusToggleButton` (caller ที่ข้าม prop นี้ หรือ toggle
  สองครั้งพร้อมกัน จะทำให้ `previous` ที่ capture ไว้ผิดตัว แล้ว revert เพี้ยน)
- **S14** — `ProgressResult` เปลี่ยนจาก `{ok: boolean}` เป็น `{kind: "ok" |
  "error"}` ให้ตรงกับ `NextUpResult`/`NextLessonResult` ในไฟล์เดียวกันที่ใช้
  `kind` discriminant อยู่แล้ว (ยังไม่เพิ่ม failure reason ตามที่ ticket บอก
  ว่าเป็นของอนาคต)
- **S15** — `today/page.tsx` เรียก `plural()` ที่ export จาก
  `TrackCard.tsx` แทนที่จะ inline ternary ซ้ำเอง
- **S16** — ซ่อน section "Other tracks" ทั้งก้อนเมื่อ `nextUp.kind !==
  "next"` (all-done/no-lessons) กัน UI โชว์ "You've read everything" ตามด้วย
  list ของทุก track ซ้ำซ้อนไม่มีจุดหมาย
- **S17** — ลบ `"use client"` ออกจาก `useFocusTrack.ts` (ไม่ใช่ component
  boundary เอง แค่ hook ที่ import จาก client component เท่านั้น)
- **Re-verify**: `npm test` 44/44 ผ่าน (35→44: +6 `api.test.ts` ใหม่ +3
  `buildOtherTracks`), `tsc --noEmit`/`eslint`/`npm run build` สะอาด,
  `go vet`/`go test` cached ผ่าน ยืนยันด้วย `git diff --name-only
  develop...` ไม่มีไฟล์ Go — mutation ทั้งหมดพังตามคาด แล้ว revert กลับ:
  - 4 mutation ใหม่ของ `api.ts` (ลบ `Array.isArray` guard, catch คืนผิด,
    `getProgress`/`getFocusTrack` กลืน Unauthorized) → พังแต่ละอันตรงตาม
    เทสต์ใน `api.test.ts` ที่เขียนไว้เฉพาะ
  - dedupe filter หายจาก `buildOtherTracks` → พัง "excludes the track
    pickNextUp is currently pointing at"
  - 7 mutation เดิมจากรอบ 1/2 (`pickNextUp` ignore focus/treat in_progress
    as done/drop has-lesson filter, `indexProgress` return-empty/swap-args,
    `computeTrackProgress` count in_progress as done,
    `sortTracksByDisplayOrder` unknown-track-first) → ยังพังทุกจุดเหมือนเดิม
    ไม่มีจุดไหน regress
- **Live-stack re-verify**: `docker compose up -d --build` แล้ว Playwright —
  malformed-200 (mock `GET /api/v1/progress` → 200 กับ `{"entries":[]}`,
  key ผิดแทนที่จะเป็น 500): banner ขึ้นจริง, Other tracks โชว์
  `"61 lessons"`/`"53 lessons"` ไม่ใช่ `0/N`, console error = 0 (ไม่มีแม้แต่
  browser's network log เพราะ status เป็น 200 จริง ๆ); expired-token-on-toggle
  (mock แค่ PUT ของ `/api/v1/prefs/focus-track` ให้เป็น 401 ส่วน GET ผ่านปกติ):
  คลิก "Set focus" แล้ว URL เปลี่ยนเป็น `/token` จริง, token ใน localStorage
  เป็น `null`, ไม่มีข้อความ "please try again" ค้างอยู่เลย; all-done +
  single-track mock: ซ่อน "Other tracks" section ทั้งก้อนจริงตาม S16; normal
  path (ไม่ mock อะไร): ผลเหมือนรอบ 2 เป๊ะ (ddia's "Latency and Percentiles"
  + Continue, ddd `4/4 done`, ai-systems `0/53 done`); กด "Set focus" ปกติ
  ยัง persist ผ่าน curl ยืนยัน; `/learn` ยังทำงานปกติหลัง hook เปลี่ยน (โชว์
  focus state ที่ set จาก `/today` ถูกต้องข้ามหน้า); 375px:
  `scrollWidth === clientWidth === 375`, console error = 0
- Status: `implemented, PR pending (review round 3 addressed)`

## UX-7 — reading-status indicators across /learn `[go-implementer]`

- **Scope ที่ตัดจริง**: เดิม UX-7 ตั้งใจทำหน้า Progress รวมแยกต่างหาก แต่ทำแค่
  indicator บน `/learn` (list + `/learn?track=`) ก็ตอบโจทย์ "อ่านไปถึงไหนแล้ว"
  ได้เกือบหมดโดยไม่ต้องมีหน้าใหม่ — หน้า Progress รวมเลื่อนเป็น **UX-8**
  (ทำเมื่อใช้ ticket นี้แล้วยังรู้สึกขาดภาพรวมจริง ๆ เท่านั้น) ไม่แตะ `/lesson`
  reader เลย (จะต้องมี data source ที่สามในหน้านั้น — deferred เหมือนกัน) และ
  ไม่แตะ Go เลยสักไฟล์ (`git diff --name-only develop...` ยืนยันว่างเปล่า)
- **`/learn` list view (`TrackCard.tsx`)**: การ์ดที่มี lesson อย่างน้อย 1 บท
  โหลด progress มาด้วย (`Promise.all([getCurriculum(), getProgress(),
  getFocusTrack()])` — reuse discipline เดียวกับ `/today` เป๊ะ) แล้วโชว์
  `ProgressBar` + ตัวเลข `{read}/{available} lessons read` แทนที่บรรทัด
  `{lessonsReady} lessons ready` เดิม (บรรทัด `{chapterCount} chapters ·
  {totalConcepts} concepts planned` ยังอยู่เหมือนเดิมเพื่อให้เห็น gap ระหว่าง
  "มีอยู่" กับ "วางแผนไว้") — ตัวส่วนของบาร์คือ **lessons available
  (`has_lesson`)** เท่านั้น ไม่ใช่ `totalConcepts`, เพราะ DDD วันนี้มี 4 lesson
  พร้อมอ่านจาก 27 concept ที่วางแผนไว้ — บาร์เต็ม "100%" ต้องแปลว่า "อ่านครบทุก
  lesson ที่มีอยู่ตอนนี้" ไม่ใช่ "จบหลักสูตรแล้ว" (กฎเดิมจาก UX-2 ที่ตอนนั้นยังไม่มี
  progress data จริงเลยไม่มีบาร์เลยสักใบ) เมื่อ `progressAvailable === false`
  การ์ดเรนเดอร์เหมือนเดิมทุกอย่าง (บรรทัด "ready" กลับมา, **ไม่มีบาร์, ไม่มี
  `0/N`**) — comment เดิมใน `ProgressBar.tsx` ที่บอกว่า "not wired into any
  page yet" ตอนนี้เท็จแล้ว เปลี่ยนเป็นอธิบาย invariant ของตัวเศษแทน
- **`/learn?track=` detail view (`TrackTopics.tsx`)**: ต่อ concept แสดง
  read-state marker เมื่อ progress โหลดสำเร็จเท่านั้น — `passed` = ✓
  (`text-success-strong`), `in_progress` = วงกลมมีจุดตรงกลาง
  (`text-warning-strong`, icon ใหม่ `InProgressIcon`), `not_started`
  **ไม่มี marker เลย** (การไม่มีคือสัญญาณเอง — ป้องกันไม่ให้ track ที่มี 61 แถว
  ต้องแบก badge "not started" ซ้ำ 55+ อัน) concept ที่ `has_lesson === false`
  ก็ไม่มี marker เหมือนกัน (ยังคง treatment เดิมจาก UX-3 ทั้งหมด — `<div>` กดไม่ได้,
  `text-muted`, "No lesson yet") ทั้งสอง state ที่มี marker แยกกันด้วย **shape**
  (เครื่องหมายถูก vs วงกลม-จุด ไม่ใช่แค่สี ผ่านทั้ง greyscale และ screen reader)
  บวก accessible name ผ่าน `aria-label="Read"` / `"In progress"` บน `<span>`
  ที่ห่อ icon (icon เองเป็น `aria-hidden`) — ตรง WCAG 1.4.1 ที่ห้ามใช้สีเป็น
  สัญญาณเดียว; ต่อ chapter header เดิมโชว์ `{available}/{total} ready` อย่างเดียว
  ตอนนี้เพิ่ม `{read}/{available} read` เป็นบรรทัดหลักเมื่อ progress พร้อม แล้วโชว์
  `{available}/{total} ready` เป็นบรรทัดรองต่อเมื่อ `available < total`
  เท่านั้น (ไม่ใช่โชว์คู่กันเสมอ) — กติกาที่ยึดคือ **ห้ามโชว์สองอัตราส่วนที่ตัวส่วน
  ต่างกันพร้อมกันโดยไม่มี label แยกให้ชัด**; `read`/`available`/`ready` เป็นสาม
  label ที่ต่างกันชัดเจนพอ เลยไม่ใช่การ "stack สามอัตราส่วน" อย่างที่ ticket เตือนไว้
  (แค่สองบรรทัดสูงสุด ไม่เคยสามพร้อมกัน) เคส `available === 0` (chapter ที่ยังไม่มี
  lesson เลยสักบท เช่น DDD's Building Blocks) ตั้งใจ**ไม่**โชว์บรรทัด "0/0 read"
  ที่ไม่มีความหมาย — fallback ไปโชว์แค่ `{available}/{total} ready` เหมือนตอน
  progress ไม่พร้อม (เจอ bug นี้จาก Playwright จริง ไม่ใช่ตอนออกแบบ — ดูหัวข้อ
  live-stack ด้านล่าง) เมื่อ `progressAvailable === false` ทั้ง accordion
  เรนเดอร์เหมือนเดิมทุกจุด (ไม่มี marker, ไม่มีบรรทัด "read" เลย)
- **`--color-warning-strong` token ใหม่**: `--color-warning` เดิม (amber-600,
  `#d97706`) วัด contrast บน `bg-surface` ได้แค่ ~2.8:1 — ต่ำกว่าเกณฑ์ non-text
  1.4.11 ด้วยซ้ำ (3:1) เพราะไม่เคยถูกใช้เป็นสี icon/text จริงมาก่อน (มีแค่ตัวแปร
  เผื่อไว้) เพิ่ม amber-700 (`#b45309`, ~4.70:1 ตามการคำนวณมือ, วัดจริงจาก
  browser ได้ 4.95:1) เป็น `--color-warning-strong` ตาม pattern เดิมของ
  `success-strong`/`danger-strong` ที่มีอยู่แล้ว แล้วให้ marker `in_progress`
  ใช้ตัวนี้แทน — ถ้าไม่แก้ตอนนี้ ticket นี้จะเป็นจุดแรกที่ใช้ `text-warning` เป็นสี
  จริงแล้ว fail AA/1.4.11 ทันที
- **Logic ทั้งหมดอยู่ใน `web/lib/curriculum.ts`** (ตามบทเรียนจาก UX-6 ที่ 11
  mutation รอดเพราะ logic ใหม่อยู่ใน JSX): `trackReadStats(track,
  progressByKey, progressAvailable)` และ `chapterReadStats(topicSlug,
  concepts, progressByKey, progressAvailable)` คืน **`null` เมื่อ
  `progressAvailable === false`** (ไม่ใช่ object ที่ zero ทุกฟิลด์) — หลักการ
  เดียวกับ `getProgress()`'s `{kind:"error"}` ทุกจุด: เลขที่ไม่รู้ค่าจริงต้องไม่
  แสดงเป็นศูนย์ที่ดูมั่นใจ; `conceptReadMarker(hasLesson, state)` คืน
  `"passed" | "in_progress" | "none"` โดย `hasLesson=false` ชนะเสมอไม่ว่า
  state จะเป็นอะไร คอมโพเนนต์ (`TrackCard`, `TrackTopics`) ทำแค่ map ผลลัพธ์
  เป็น markup เท่านั้น
- **Mutation-testing** (แก้ source จริง รัน `npm test` แล้ว revert ทุกครั้ง;
  fixture ใช้ 2 chapters + ผสมครบสาม state ตามที่ ticket บังคับ, ดู
  `readStatsTrack`/`readStatsProgress` ใน `curriculum.test.ts`):
  1. **(บังคับ) นับ `in_progress` เป็น read** — `chapterReadStats`:
     `state === "passed"` → `state === "passed" || state === "in_progress"`
     → พัง **"counts only passed concepts as read, out of has_lesson
     concepts in the chapter"**
  2. **(บังคับ) สลับตัวส่วนของ chapter** — `chapterReadStats`: ลบ
     `if (!concept.has_lesson) continue;` ออก (นับทุก concept เป็น
     available) → พัง **"excludes a lesson-less concept from available
     while still counting it in total"**
  3. **(บังคับ) invert `progressAvailable` (ฝั่ง track card)** —
     `trackReadStats`: `if (!progressAvailable) return null;` →
     `if (progressAvailable) return null;` → พังทั้ง **"counts passed
     concepts across every chapter as read, out of lessons available"**
     และ **"returns null when progress is unavailable, not a zeroed-out
     object"**
  4. **(บังคับ) ลบ has_lesson filter ออกจาก marker** —
     `conceptReadMarker`: ลบ `if (!hasLesson) return "none";` → พัง
     **"returns none for a concept without a lesson, regardless of its
     recorded state"**
  5. invert `progressAvailable` ฝั่ง chapter ด้วย (คนละฟังก์ชันจาก #3,
     คุมคนละจุด) — `chapterReadStats`: same flip → พังพร้อมกัน 3 เทสต์
     ("counts only passed...", "excludes a lesson-less concept...",
     "returns null when progress is unavailable...")
  6. `conceptReadMarker` ยุบ `in_progress` ไปเป็น `"passed"` (marker สอง
     state ไม่ต่างกันอีกต่อไป — ขัด requirement "distinct markers" ตรง ๆ) →
     พัง **"returns a distinct in_progress marker rather than collapsing it
     into passed"**
  7. (bonus) ย้าย `total += 1` เข้าไปอยู่หลัง `if (!has_lesson) continue`
     (ทำให้ `total === available` เสมอ, ไม่นับ concept ที่ไม่มี lesson ใน
     total ด้วย) → พังเทสต์เดียวกับ #2 (**"excludes a lesson-less
     concept..."**) เพราะเทสต์นั้นเช็คทั้ง `available` และ `total` แยกกัน —
     ไม่มี mutation ไหนรอดเลยจากทั้ง 7 จุด
- **ทดสอบกับ live stack จริง** (`docker compose up -d --build`, DB local มี
  progress จริงจากการทดสอบ ticket ก่อนหน้าอยู่แล้ว: ddd ผ่านครบ 4/4, ddia ผ่าน
  4 บทจาก 61 + ผสม in_progress, ai-systems มี 2 concept `in_progress` จาก 53
  — ยืนยันตัวเลขทุกจุดด้วย `curl /api/v1/progress` + `/api/v1/curriculum`
  ก่อนเทียบกับ Playwright):
  - **Normal path, list view**: การ์ดโชว์ `4/4 lessons read` (ddd, บาร์เต็ม),
    `4/61 lessons read` (ddia, บาร์ ~7%), `0/53 lessons read` (ai-systems,
    บาร์ว่าง) ตรงกับตัวเลขจาก curl เป๊ะ, มี `role="progressbar"` 3 อัน
    (เท่าจำนวนการ์ดที่มี lesson), console error = 0
  - **Normal path, detail view**: `/learn?track=ddd` chapter "Model-Driven
    Foundations" (4/4 lesson ผ่านหมด) header โชว์ `"4/4 read"` เฉย ๆ (ไม่มี
    บรรทัด ready เพราะ available===total) มี marker "Read" (aria-label)
    ครบ 4 อัน; `/learn?track=ai-systems` chapter "LLM Foundations" header
    `"0/5 read"` มี marker "In progress" 2 อัน ตรงกับ 2 concept ที่ backend
    ส่งมาเป๊ะ, marker "Read" = 0 อันในหน้านี้ (ยังไม่มี concept ไหนผ่านเลยใน
    track นี้); `/learn?track=ddd` chapter "Building Blocks" (0 lesson พร้อม
    จาก 10 concept) header โชว์ `"0/10 ready"` อย่างเดียว **ไม่มี** `"0/0
    read"` โผล่มา — นี่คือ edge case ที่เจอจาก Playwright จริง ไม่ได้คิดไว้
    ตอนออกแบบแรก ต้องเพิ่มเงื่อนไข `readStats.available > 0` ก่อนเข้าโหมด
    read-line ถึงจะปิดจุดนี้ได้
  - **Mock `GET /api/v1/progress` → 500**: banner `role="alert"` ขึ้นจริง 1
    อัน (แยกจาก Next.js's `__next-route-announcer__` ที่มี `role="alert"`
    ติดตัวเองอยู่แล้วทุกหน้า — ต้อง exclude ID นี้ตอนนับ banner ไม่งั้นนับผิด),
    `role="progressbar"` = 0 ทั้งหน้า, ไม่มี string `"0/N"` โผล่ที่ไหนเลยในหน้า
    (regex กวาดทั้ง `<main>`), console error มีแค่ 1 บรรทัด "Failed to load
    resource: 500" ซึ่งเป็น browser network log อัตโนมัติ ไม่ใช่โค้ดแอป
  - **Mock `GET /api/v1/progress` → 200 กับ `{"entries":[]}`** (malformed —
    key ผิด `concepts`): banner ขึ้นจริง 1 อัน ที่ `/learn?track=ddia`, ไม่มี
    marker (`aria-label="Read"`/`"In progress"`) เลยสักอันในหน้า, chapter
    "Replication" fallback กลับไปโชว์ `"5/5 ready"` แบบเดิม ไม่ใช่ read line,
    console error = 0 (status เป็น 200 จริง ไม่มี network log)
  - **375px**: `scrollWidth === clientWidth === 375` ที่ทั้ง list และ detail
    view ที่ chapter เปิดอยู่
  - **Contrast วัดจริงจาก browser** (`getComputedStyle` + คำนวณ relative
    luminance เอง ไม่ใช่อ่านจาก palette): marker `passed`
    (`text-success-strong`, `rgb(4,120,87)`) บน card surface
    (`rgb(255,253,250)`) = **5.40:1**; marker `in_progress`
    (`text-warning-strong`, `rgb(180,83,9)`) = **4.95:1**; ตัวหนังสือ
    `{n}/{N} lessons read` บน card (`text-muted`, `rgb(82,82,91)`) =
    **7.61:1** — ทั้งสามผ่าน AA (≥4.5:1) จริง ไม่ใช่แค่คำนวณมือ
  - **Regression spot-check**: กด "Set focus" จากการ์ด DDD หลัง refactor
    `TrackCard`/`learn/page.tsx` แล้ว `curl GET
    /api/v1/prefs/focus-track` ยืนยัน persist จริง (`{"track":"ddd"}`),
    ไม่มี `pageerror` เกิดขึ้น — reset กลับเป็น `null` หลังทดสอบเสร็จ
- **จงใจไม่ทำในรอบนี้**: หน้า `/lesson` reader ไม่มี indicator ใด ๆ เพิ่ม (นอก
  scope ตามที่ ticket ระบุ — เลื่อนไป UX-8 หรือ ticket แยกถ้าจำเป็น); ไม่ได้เพิ่ม
  jsdom/`.tsx` test ให้ `TrackCard`/`TrackTopics` เอง (ทั้งคู่ยังเป็น debt เดิม
  จาก UX-6's S10 — coverage ของ ticket นี้อยู่ที่ pure function ใน
  `curriculum.ts` ทั้งหมด ตรวจ component ผ่าน Playwright แทน)
- **Review focus**:
  - ทำไม `trackReadStats`/`chapterReadStats` ต้องคืน `null` เมื่อ
    `progressAvailable === false` แทนที่จะคืน `{read: 0, available: 0}`?
  - ทำไม chapter ที่ `available === 0` (ยังไม่มี lesson เลยสักบท) ถึงต้อง
    fallback ไปโชว์ `"0/10 ready"` แทนที่จะโชว์ `"0/0 read"` ทั้งที่ทั้งคู่คำนวณ
    ถูกต้องทางคณิตศาสตร์เหมือนกัน?
  - ทำไม concept marker ถึงต้องแยกเป็นทั้ง shape (เครื่องหมายถูก vs
    วงกลม-จุด) และ `aria-label` พร้อมกัน ไม่ใช้แค่สีคนละสีก็พอ?
- Status: `implemented, PR pending`
