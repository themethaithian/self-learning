# Quiz (Q-series) — tickets

## Q-1 — recall-first 3-stage quiz `[go-implementer]`

- **Scope**: `web/components/RecallCheckCard.tsx` เขียนใหม่ทั้งไฟล์ — จากเดิม 2
  stage (Reveal answer → self-rate Pass/Not yet) เป็น 3 stage ต่อ recall check:
  **Recall** (เห็นแค่คำถาม, mcq ยังไม่โชว์ตัวเลือก) → **Commit** (mcq: เลือก
  ตัวเลือก + confidence, short_answer: เลือก confidence อย่างเดียว) →
  **Reveal** (โชว์ `expected_answer`, mcq mark ถูก/ผิดจากการเทียบข้อความ ไม่ใช่
  self-report, short_answer ยังใช้ Pass/Not yet self-rate เหมือนเดิมเพราะ
  free-text เทียบเองไม่ได้). ไฟล์ใหม่ `web/lib/shuffle.ts` (Fisher-Yates,
  `rng` inject ได้), `web/lib/shuffle.test.ts`,
  `web/components/RecallCheckCard.test.tsx`,
  `web/app/(app)/lesson/page.test.tsx`. แก้ `web/app/(app)/lesson/page.tsx`
  ให้ผูกกับ API ใหม่ของการ์ด (`onFinishedChange` แทน `rating`/`onRate`) และ
  `web/components/icons.tsx` เพิ่ม `XIcon` (สัญลักษณ์ผิด คู่กับ `CheckIcon` เดิม).
  **Frontend ล้วน** — `git diff --name-only develop... -- '*.go'` ว่างเปล่า
  จริง, `go vet`/`go test` รันผ่านเพราะไม่มีอะไรให้กระทบ

### การตัดสินใจหลัก

- **ทำไม 3 stage**: จุดของ ticket นี้คือบังคับ *free recall* ก่อน
  *recognition* — mcq ที่โชว์ตัวเลือกตั้งแต่แรกกลายเป็นแค่ "จำหน้าตัวเลือกได้"
  ไม่ใช่ "นึกคำตอบได้เอง" ซึ่งเป็นสัญญาณ recall ที่ SRS (Q-2) ต้องการ.
  Confidence (Guessed/Unsure/Confident) ถูกเก็บ **ก่อน** เห็นเฉลยเสมอ เพราะ
  "มั่นใจแต่ผิด" คือสัญญาณที่มีค่าที่สุดสำหรับ SRS และเก็บย้อนหลังไม่ได้ถ้าถาม
  หลังเฉลย
- **`expected_answer` ยังอยู่ใน payload ต่อไป — ตัดสินใจแล้ว ไม่ต้อง revisit**:
  `docs/roadmap.md` เดิมเขียนว่าต้องตัด `expected_answer` ออกจาก `GET
  /lessons/{topic}/{concept}` ก่อนทำ ticket นี้ (อ้าง `design.md` §API ที่พูดถึง
  โมเดล LLM-graded ในอนาคต ที่คำตอบฝั่ง client จะไป corrupt server-side
  scoring) — แต่ v1 เป็น **self-graded** (คนตอบ = คนให้คะแนนเอง) ทำ endpoint
  เฉลยแยกไม่ได้อะไรเพิ่ม มีแต่เสีย round trip + failure state ใหม่ (network
  error ตอนกด reveal, token หมดอายุตอนกด reveal ฯลฯ) ค่อยย้าย server-side ตอน
  Q-2 ที่เริ่ม submit attempt จริงและออกแบบ endpoint จาก requirement จริง.
  กฎเดียวที่ยังผูกอยู่: **`expected_answer` ต้องไม่ถึง DOM ก่อน stage 3**
  (component เดิมทำถูกอยู่แล้ว ตอนนี้พิสูจน์ด้วยเทสต์ตรง ๆ — ดูตาราง mutation)
- **State ของแต่ละการ์ดเก็บใน `RecallCheckCard` เอง ไม่ยกขึ้น parent**: stage,
  ตัวเลือกที่เลือก, confidence, short-answer rating เป็น `useState` ภายใน
  component — `lesson/page.tsx` รู้แค่ boolean `finished` ต่อ `position`
  (ผ่าน `onFinishedChange(position, finished)`) เพราะไม่มีใครต้องอ่าน
  รายละเอียดเหล่านั้นใน ticket นี้ (Q-2 ค่อย persist จริง). "finished"
  ต่อ type ถูกนิยามที่**จุดเดียว** — `isRecallCheckFinished(checkType, stage,
  shortAnswerRating)` (exported, มี table-driven test ตรง ๆ ไม่ต้อง mount
  component เพื่อเทสต์ logic นี้)
- **Shuffle: lazy `useState` ไม่ใช่ `useMemo`**: `useState(() =>
  shuffleOptions(...))` รันครั้งเดียวจริงตาม React ยืนยัน (`useMemo` เป็นแค่
  cache ที่ React มีสิทธิ์ทิ้งได้) — reshuffle กลางอากาศจะทำให้ตัวเลือกขยับใต้
  นิ้วผู้ใช้ระหว่างแตะกับยืนยัน. `rng` เป็น optional param (default
  `Math.random`) ไหลจาก component prop → `shuffleOptions` โดยตรง ไม่มี branch
  แยกสำหรับ test environment
- **ระบุตัวถูกด้วย index ของ shuffled array ไม่ใช่ข้อความ**: `selectedIndex`
  (ไม่ใช่ `selectedOption` string) กัน bug ตัวเลือกซ้ำข้อความกันสองตัวถูก
  select พร้อมกัน; ส่วนความถูก/ผิดยัง derive จากการเทียบ `option ===
  expected_answer` เสมอ (ตาม spec — ต้อง "รอดจากการ shuffle" ได้)
- **A11y**: options และ confidence เป็น `<fieldset><legend>` + native
  `<input type="radio">` จริง (ไม่ใช่ custom `role="radiogroup"` + roving
  tabindex) — arrow-key navigation และการประกาศ "group, N of M" ได้มาฟรีจาก
  browser, ไม่ต้องเขียน JS เอง. confidence ใช้ chip look (`sr-only` input +
  label ที่ style เป็น pill) ยังเป็น radio จริงใต้ผิว. Stage transition ย้าย
  focus ไปที่ region ใหม่เสมอ (`tabIndex={-1}` + `.focus()` ใน `useEffect`
  คีย์ด้วย `stage` — pattern เดียวกับ `finishedBannerRef` ใน UX-5) และ stage 3
  ห่อด้วย `role="status"` ให้ screen reader ได้ยินผลจริง ไม่ใช่โผล่มาเงียบ ๆ.
  ถูก/ผิดมีทั้งไอคอน (`CheckIcon`/`XIcon`) และข้อความ ("Correct answer"/"Your
  answer — incorrect") ไม่ใช่สีอย่างเดียว (WCAG 1.4.1)

### สิ่งที่ค้นพบระหว่าง mutation testing — per-check reset ไม่ได้พึ่ง `key` อย่างเดียว

ตั้งใจจะ key การ์ดด้วย `` `${lesson.topic}/${lesson.concept}/${check.position}` ``
เพื่อบังคับ remount ตอนเปลี่ยน lesson แต่พอลองมูเทตจริง (revert กลับเป็น
`key={check.position}` เฉย ๆ) เทสต์ **ไม่ตาย** — เหตุผลคือ `lesson/page.tsx`
เดินผ่าน `{status: "loading"}` (แสดง `LessonSkeleton`) ทุกครั้งที่ topic/concept
เปลี่ยน **ก่อน** จะ render success state ใหม่ ซึ่งเป็นคนละ JSX tree กับตอน
success ทั้งก้อน — React unmount ทุกอย่างใต้ success branch (รวม
`RecallCheckCard` ทุกใบ) แล้ว mount ใหม่เมื่อ lesson ถัดไปโหลดเสร็จ, ไม่ว่า key
จะเป็นอะไรก็ตาม. ยืนยันด้วยการมูเทตจริง: comment
`setState({status:"loading"})` ทิ้ง (จำลอง "query-param navigation ไม่
remount หน้า" ตามที่ ticket เตือนไว้ตรง ๆ) → เทสต์ตายทันที (เห็นคำตอบเก่าเลค
เข้ามาจริง) แล้ว revert กลับ. สรุป: **`key={check.position}` เดิมพอแล้ว**
(ไม่ต้องเปลี่ยน) เพราะ reset มาจาก architecture ของ loading-branch, ไม่ใช่จาก
key — แต่ `setFinishedChecks({})` ใน effect เดียวกับที่ reset lesson ยังคงไว้
เป็น defense-in-depth (เผื่อ refactor ในอนาคตที่เอา loading flash ออก เช่น
"keep previous content visible while loading" ซึ่งจะทำให้ key กลับมาจำเป็น
จริง ๆ)

### Mutation table (11/11 killed)

แก้ source ทีละจุด รัน `npm test` แล้ว revert ทุกครั้ง:

| # | Mutation | ผลลัพธ์ |
|---|---|---|
| M1 | `canReveal` ตัด `selectedIndex != null` ออก (mcq reveal ได้โดยไม่เลือกตัวเลือก) | killed — "does not reach reveal (mcq) without an option selected, even with confidence set" |
| M2 | `canReveal` ตัด `confidence != null` ออก (reveal ได้โดยไม่มี confidence) | killed — "does not reach reveal (mcq) without a confidence..." + "...(short_answer) without a confidence" |
| M3 | render ตัวเลือก mcq ตั้งแต่ stage "recall" | killed — "shows only the question, no options and no expected_answer, for mcq" |
| M4 | render กล่อง `expected_answer` แบบไม่มีเงื่อนไข (นอก stage 3) | killed — "never puts expected_answer in the serialized DOM before reveal" + 2 เทสต์อื่นที่เช็ค DOM ระหว่าง stage 2 |
| M5 | shuffle ถูกแทนด้วย identity (`.slice()` เฉย ๆ) | killed — "is stable across re-renders" (rng ไม่ถูกเรียกเลย) + "actually reorders the source array" + "marks the correct option by matching expected_answer, not by its original index" |
| M6 | `shuffleOptions` เปลี่ยนเป็น `sort(() => rng() - 0.5)` | killed — **distribution test** ใน `shuffle.test.ts` ("distributes every option to every position near-uniformly") จับได้จริง (cell นึงเบี่ยงไปถึง ~7171 จาก expected 5000, tolerance ±750) + เทสต์ correctness-survives-shuffle 2 ตัวใน component |
| M7 | ย้าย shuffle ออกจาก lazy `useState` initializer มาคำนวณตรง ๆ ใน render body | killed — "is stable across re-renders (lazy useState, not recomputed every render)" (rng call count เพิ่มขึ้นทุก rerender) |
| M8 | ตรวจตัวถูกด้วย index เดิมก่อน shuffle (`i === options.indexOf(expected_answer)`) แทนการเทียบข้อความ | killed — "marks the correct option by matching expected_answer, not by its original index" |
| M9 | `isRecallCheckFinished` คืน `true` ตั้งแต่ stage "commit" | killed — table-driven test 2 แถว + "reports mcq finished only once stage 3 is reached" |
| M10 | `isRecallCheckFinished` ไม่เช็ค `shortAnswerRating` เลย (คืน `true` ทันทีที่ stage เป็น reveal) | killed — table-driven test 1 แถว + "does not report short_answer finished at stage 3 until Pass/Not yet is chosen" |
| M11 | ลบ `setState({status:"loading"})` ออกจาก effect เปลี่ยน lesson (จำลอง "ไม่ remount" ตรงตามที่ ticket เตือน) | killed — "does not leak the previous lesson's revealed answer or finished state into the next lesson's same-position check" (เห็น `expected_answer` ของ lesson เก่าเลคเข้ามาจริงตอน mutate) |

**11/11 killed** — ไม่มี survivor

### Re-verify เต็มชุด

- `npm test`: **87 → 117** (+30: 4 `shuffle.test.ts` + 25
  `RecallCheckCard.test.tsx` + 1 `lesson/page.test.tsx`)
- `npx tsc --noEmit`, `npx eslint .`, `npm run build` — สะอาดทั้งสามคำสั่ง
- `go vet ./...` + `go test ./...` — ผ่านทั้งหมด (ไม่มีไฟล์ `.go` เปลี่ยนเลย
  ใน ticket นี้)

### Playwright evidence (live stack, `docker compose up -d --build`, ปิดด้วย `docker compose down` เปล่า ๆ — ไม่มี `-v`)

ทดสอบกับ `domain-driven-design/ubiquitous-language` (มีทั้ง `short_answer`
และ `mcq` checks จริงจาก DB):

- **Stage 1**: ดึง `expected_answer` ทุกตัวจาก `GET
  /api/v1/lessons/.../ubiquitous-language` มาเทียบกับ `page.content()`
  (serialized HTML เต็ม ไม่ใช่แค่ visible text) — **0 ตัวเลค**
- **Stage 2**: กด "Show options" แล้วเช็คว่ากล่อง "Answer" (reveal box)
  **ยังไม่ปรากฏ** (count = 0); กด "Reveal answer" ทั้งที่ยังไม่ commit อะไร
  เลย (`{force:true}` ข้าม visual disabled) → `role="status"` count ยังเป็น
  **0** (guard ในโค้ดกันไว้จริง ไม่ใช่แค่ `aria-disabled` มองไม่เห็น)
- **Keyboard-only end-to-end**: focus ตัวเลือกแรกด้วย `.focus()` แล้วกด
  `ArrowDown` → มี option ถูก check จริง (count 1); ทำนองเดียวกันกับ
  confidence ด้วย `ArrowRight`; focus ปุ่ม Reveal แล้วกด `Enter` →
  `role="status"` โผล่จริง ยืนยัน "Reveal reached via keyboard-only
  interaction: OK"
- **หลัง reveal**: `expected_answer` โผล่ใน HTML จริง; ข้อความใน
  `role="status"` มีทั้ง "Not quite"/"Correct" และ tag ต่อ option
  ("Correct answer" / "Your answer — incorrect") ที่**เป็นข้อความจริงในต้นไม้
  accessibility ไม่ใช่แค่สี** — พิสูจน์ด้วย `ariaSnapshot()` ที่โชว์ tag พวกนี้
  ตรง ๆ ใน `listitem`; `document.activeElement` เป็น element เดียวกับ
  `role="status"` (focus ย้ายไปที่ผลลัพธ์จริง)
- **AX tree ที่ stage 2** (`ariaSnapshot()`): เห็น `group "Choose one"` และ
  `group "How confident are you?"` ห่อ `radio` ที่มี accessible name เป็น
  ข้อความตัวเลือกจริง (มาจาก `fieldset`/`legend` โดย browser แม็ปให้เอง) —
  ยืนยัน radiogroup semantics โดยไม่ต้องเขียน ARIA เอง
- **Contrast วัดจริงจาก browser** (`getComputedStyle` + composite ของทุก
  background layer ที่ element สืบทอด แปลง `oklab()` ที่ Tailwind v4's `/10`
  opacity utility serialize กลับมาเป็น sRGB ด้วยสูตร CSS Color 4 ก่อนคำนวณ
  relative luminance — อ่าน raw `rgba()`/`oklab()` เฉย ๆ โดยไม่ composite จะ
  ได้ค่าผิด/ทึบเกินจริง): banner ผลลัพธ์ (fg `rgb(190,18,60)` บน bg
  `rgb(252,231,232)`) = **5.31:1**, แถวตัวเลือกถูก (fg `rgb(4,120,87)` บน bg
  `rgb(230,243,236)`) = **4.80:1**, แถวตัวเลือกผิดที่ผู้ใช้เลือก = **5.31:1**
  — **ทั้งสามผ่าน AA (≥4.5:1)**
- **375px**: `scrollWidth === clientWidth === 375` จริง ไม่มี horizontal
  scroll
- **axe-core (รันทุก rule ไม่ scope เฉพาะ wcag2a/wcag2aa)**: 1 violation —
  `scrollable-region-focusable` บน `.my-4.overflow-x-auto` ซึ่งเป็นโค้ดบล็อก
  ของ `LessonBody.tsx` (**ไฟล์นี้ ticket นี้ไม่ได้แตะเลย** — `git diff
  --name-only develop -- '*.tsx' '*.ts'` ยืนยันว่าไม่อยู่ในรายการไฟล์ที่แก้ —
  pre-existing บน `develop` ก่อน ticket นี้, นอก scope) — **0 violations**
  จาก `RecallCheckCard` เอง
- Console error / page error = 0 ตลอดการทดสอบ

### Review focus

- ทำไม `selectedIndex` ต้องเป็น index ของ shuffled array แทนที่จะเก็บข้อความ
  ตัวเลือกตรง ๆ เหมือน `rating` เดิม?
- ทำไม confidence ต้องถูกเก็บ**ก่อน**กด Reveal เสมอ แทนที่จะถามพร้อมกับ
  Pass/Not yet ตอน stage 3?
- ทำไมการมูเทต `key={check.position}` กลับเป็นค่าเดิม (ไม่ใส่ lesson identity
  เข้าไป) ถึงไม่ทำให้เทสต์ตาย ทั้งที่ ticket เตือนเรื่อง per-check state reset
  ไว้ตรง ๆ?

### จงใจไม่ทำในรอบนี้

- ไม่ persist อะไรเพิ่มจากเดิม — `selectedIndex`/`confidence`/
  `shortAnswerRating` อยู่ใน memory ของ `RecallCheckCard` เท่านั้น เหมือนที่
  `ratings` เดิมก็อยู่แค่ใน `lesson/page.tsx` — Q-2 คือ ticket ที่ทำให้ทุกอย่าง
  persist จริง
- ไม่แก้ `scrollable-region-focusable` ของ `LessonBody.tsx` — pre-existing,
  ไม่เกี่ยวกับไฟล์ที่ ticket นี้แตะ
- ไม่ตัด `expected_answer` ออกจาก payload — ตัดสินใจแล้วว่าไม่ทำใน v1 (ดู
  หัวข้อการตัดสินใจด้านบน)

Status: implemented, PR pending
