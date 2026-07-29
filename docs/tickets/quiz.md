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
  `web/app/(app)/lesson/page.test.tsx`,
  `web/app/(app)/lesson/page.finishedReset.test.tsx` (round 2, mock
  `RecallCheckCard` เพื่อ pin `setFinishedChecks({})` แยกจาก self-correction
  ของ component จริง). แก้ `web/app/(app)/lesson/page.tsx` ให้ผูกกับ API ใหม่
  ของการ์ด (`onFinishedChange` แทน `rating`/`onRate`), `web/components/icons.tsx`
  เพิ่ม `XIcon` (สัญลักษณ์ผิด คู่กับ `CheckIcon` เดิม), และ (round 2)
  `web/components/Button.tsx` เพิ่ม `type="button"` default (ดูรอบ
  code-reviewer). **round 2 ยังมีสองจุดที่ไม่ได้อยู่ใน REQ table ของ reviewer**
  (เป็น cleanup ที่ทำไปพร้อมกัน): stage-1 label เปลี่ยนจากเช็ค `check.type ===
  "mcq"` มาเป็นเช็ค `hasOptions` แทน (**เป็นการเปลี่ยนพฤติกรรมจริง**: mcq ที่
  `options` ว่าง/ไม่มีจะขึ้น "I've answered" แทน "Show options" ตอนนี้ — **ไม่มี
  เทสต์ pin จุดนี้เลย ยังเป็น unpinned survivor อยู่**, ยืนยันด้วยการ revert
  กลับเป็น `check.type` แล้ว `npm test` ยังเขียว 124/124); และตัด `export`
  ออกจาก `RecallRating`/`Confidence`/`RecallStage` (ไม่มีที่ไหนใช้จากนอกไฟล์แล้ว
  หลัง `page.tsx` เลิก import `RecallRating`). **Frontend ล้วน** — `git diff
  --name-only develop... -- '*.go'` ว่างเปล่าจริง, `go vet`/`go test` รันผ่าน
  เพราะไม่มีอะไรให้กระทบ

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
  **กฎที่ผูกอยู่จริงไม่ใช่ "expected_answer ห้ามถึง DOM ก่อน stage 3"** — สำหรับ
  mcq ข้อความคำตอบ**คือ**หนึ่งในตัวเลือกที่โชว์ตั้งแต่ stage 2 อยู่แล้ว (พิสูจน์
  ด้วย Playwright จริง) ประโยคนี้เป็น wording ที่ผิดมาจาก ticket brief เอง.
  กฎที่ถูกต้องและพิสูจน์แล้วจริงคือ: **ไม่มีอะไร mark ว่าตัวเลือกไหนถูกก่อน
  stage 3** (โชว์ตัวเลือกเฉย ๆ ไม่บอกว่าอันไหนถูก), **shuffle ตอน render ทำให้
  ตำแหน่งไม่ใบ้คำตอบ**, และ **คำตอบของ check อื่นที่ยังไม่ถูกเปิดไม่โผล่ที่ไหน
  เลย** — ดูตาราง mutation
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
  คีย์ด้วย `stage` — pattern เดียวกับ `finishedBannerRef` ใน UX-5) ทำให้ screen
  reader ได้ยินการเปลี่ยน stage ทุกครั้งโดยไม่ต้องมี live region ห่อทั้ง panel.
  `role="status"` ใช้แค่กับบรรทัดผลลัพธ์ถูก/ผิดของ mcq เท่านั้น (แก้ใน round 2
  หลัง code-reviewer ชี้ว่าห่อทั้ง panel รวม Pass/Not yet buttons ของ
  short_answer จะทำให้ live region re-announce ทุกครั้งที่ `aria-pressed`
  เปลี่ยน — ดูรอบ code-reviewer ด้านล่าง). ถูก/ผิดมีทั้งไอคอน (`CheckIcon`/
  `XIcon`) และข้อความ ("Correct answer"/"Your answer — incorrect") ไม่ใช่สี
  อย่างเดียว (WCAG 1.4.1)

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

### Mutation table — รอบแรก (11/11 mutants tried, 11 killed)

แก้ source ทีละจุด รัน `npm test` แล้ว revert ทุกครั้ง. **"11/11 killed" หมายถึง
มูเทชัน 11 จุดที่ลองในรอบแรกเท่านั้น ไม่ใช่ข้อสรุปว่าไม่มี gap เหลือ** — รอบ
code-reviewer ด้านล่างเจอ gap จริงที่ชุดนี้ไม่ครอบคลุม:

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

### รอบ code-reviewer (REQUEST_CHANGES → แก้ครบ)

Reviewer รัน mutation set ของตัวเองเพิ่มอีก **30 จุด** เจาะเฉพาะจุดที่รอบแรกไม่
ได้ทดสอบตรง ๆ — **ยืนยันอิสระว่า gate/shuffle/leak ที่รอบแรกอ้างว่า kill แล้ว
ยัง kill จริง**: reveal gate กัน force-click/keyboard Enter/keyboard
Space/programmatic `.click()` ได้ทั้งหมด, distribution test จับ
`sort(()=>rng()-0.5)`, swap-first-two-only, naive Fisher-Yates ที่ `j` วิ่ง
เต็ม range, และ off-by-one bound ได้ที่ 12.2σ tolerance
(P(spurious failure) ≈ 1e-33), 40 mount จริงได้ 6 ตำแหน่งครบ (ไม่ pin),
edge case ทุกแบบ (options ว่าง/ไม่มี/1 ตัว/ซ้ำ, `expected_answer` ไม่ match
อะไรเลย, `recall_checks` ว่าง) ไม่ crash, contrast/375px/axe ที่ 1440px ยืนยัน
ซ้ำได้, และ loading-state unmount ที่รอบแรกอ้างว่าเป็นกลไก reset จริง
ถูกยืนยันอิสระผ่าน real client-side navigation + browser Back ด้วย.

ในจำนวนนั้นมี **6 มูเทชันที่เป็นการมูเทตจริง** (REQ-1, REQ-2, REQ-5, REQ-6,
REQ-7a, REQ-7b) กับ **2 จุดที่เป็น defect จริงที่พบนอกกระบวนการมูเทต ไม่ใช่
mutant เลย** (REQ-3, REQ-4 — พบจากการอ่านโค้ด/ทดสอบ browser ตรง ๆ):

จาก 6 มูเทชัน: **5 เป็น survivor ที่เป็น coverage gap จริง** (REQ-1, REQ-2,
REQ-5, REQ-7a, REQ-7b — แก้แล้วด้วยเทสต์ใหม่ด้านล่าง, ตายหมดหลังแก้) และ
**1 เป็น equivalent mutant ที่ยอมรับได้** (REQ-6 — พฤติกรรมเหมือนเดิมจริง
ไม่ต้อง kill, อธิบายในแถวของมัน). ส่วน REQ-3/REQ-4 แก้ตรงที่โค้ดโดยตรง
(ไม่ใช่แก้ด้วยเทสต์ที่ kill มูเทชัน เพราะไม่มีมูเทชันให้ kill ตั้งแต่แรก):

| # | Mutation | ผลลัพธ์ | แก้อย่างไร |
|---|---|---|---|
| REQ-1 | Production path (ไม่ pass `rng` prop) ไม่เคยถูกเทสต์ตรง ๆ — เปลี่ยน `shuffleOptions(check.options, rng)` เป็น `shuffleOptions(check.options, rng ?? (() => 0))` | **survived** — 117/117 เขียวเดิม, 40 mount ได้ order เดียว (`BCA`) ทุกครั้ง | เพิ่มเทสต์ mount 40 ครั้งแบบไม่ pass `rng` เลย เก็บ order ทั้งหมดใส่ `Set`, assert `size > 1` — kill แล้ว |
| REQ-2 | `onChange={() => setConfidence(opt.value)}` → `setConfidence("guessed")` แข็ง | **survived** — ไม่มีเทสต์ไหน assert `.checked` ของ confidence radio เลย, คลิก "Confident" จริงได้ "Guessed" ติ๊กแทน | เพิ่ม `it.each` 3 ค่า (Guessed/Unsure/Confident) คลิกแล้ว assert เหลือแค่ตัวที่คลิก `.checked` — kill แล้ว |
| REQ-3 | `role="status"` ห่อทั้ง reveal panel รวม Pass/Not yet buttons ของ short_answer | **ไม่ใช่ mutation ที่ทดสอบ — เป็น defect ที่มีอยู่จริงในโค้ด** live region ห่อ interactive control จะ re-announce ทุกครั้งที่ `aria-pressed` เปลี่ยน | ย้าย `role="status"` มาห่อเฉพาะบรรทัด Correct/Not quite (mcq เท่านั้น); short_answer stage 3 ไม่มี live region เลย (focus ที่ container ทำงานพอ) — เพิ่มเทสต์ "does not wrap the short_answer Pass/Not yet buttons in a live region" (ยืนยันด้วยการ revert `role="status"` กลับไปห่อทั้ง panel เหมือนรอบแรกเป๊ะ ๆ → **เทสต์ตายพอดี 1 ตัว** คือเทสต์ตัวนี้เอง, `Tests 1 failed | 123 passed`) |
| REQ-4 | `Button.tsx` ไม่ตั้ง `type` default → `<button>` ในการ์ดกลายเป็น `type="submit"` โดยไม่ตั้งใจ | **ไม่ใช่ mutation — เป็น defect จริง** ตรวจแล้ว: repo มี `<form>` เดียว (`app/token/page.tsx`) ซึ่งตั้ง `type="submit"` ชัดเจนอยู่แล้ว จึงไม่มีที่ไหนพึ่งพฤติกรรม submit โดยปริยาย | เพิ่ม `type="button"` เป็น default ใน `Button.tsx` (วางก่อน `{...props}` เพื่อให้ caller override ได้) |
| REQ-5 | `isYourAnswer = i === selectedIndex` → `option === selectedOptionText` (เทียบข้อความแทน index) | **survived ครึ่งเดียว** — เทสต์ duplicate เดิมเช็คแค่ stage 2 (`.checked`), ไม่ได้เช็ค stage 3 ("your answer" tag) | ขยายเทสต์เดียวกันให้ทำต่อถึง stage 3: fixture `["Same","Same","Different"]`, เลือกตัวที่ 2 ("Same" ตัวหลัง), assert แถวแรก (ข้อความเหมือนกันแต่คนละ index) **ไม่มี** tag "Your answer" — kill แล้ว |
| REQ-6 | `key={i}` → `key={option}` บน label ของตัวเลือก mcq | **survived แต่เป็น equivalent mutant ที่ยอมรับได้** — DOM เหมือนกันทุก byte, ไม่มี React key-warning เพราะเทสต์ duplicate ไม่มี string ซ้ำแบบที่ key จะชนกันจริงในทางที่สังเกตได้จาก DOM/behavior | **ไม่แก้โค้ด** (เก็บ `key={i}` ไว้ตามเดิม — ยังถูกกว่าในหลักการ) แค่เปลี่ยนชื่อเทสต์จาก "keys options by shuffled position, not text" (อ้างว่าเทส key ทั้งที่เทสแค่ selection independence) เป็นชื่อที่ตรงกับสิ่งที่วัดจริง |
| REQ-7a | ลบ `setFinishedChecks({})` ออกจาก effect เปลี่ยน lesson | **survived** — RecallCheckCard จริงเรียก `onFinishedChange(position, false)` เองตอน mount (stage เริ่มที่ "recall" เสมอ) ซึ่ง self-correct ค่าเก่าทันทีอยู่แล้ว ทำให้การลบบรรทัดนี้ไม่มีผลสังเกตได้ผ่าน component จริง | เพิ่มไฟล์ `page.finishedReset.test.tsx` ที่ mock `RecallCheckCard` เป็น fake ที่**ไม่** self-report ตอน mount (มีแค่ปุ่มกด "Mark finished" ตรง ๆ) — แยก concern การ reset ของ parent ออกจาก self-correction ของ child จริง แล้ว assert ว่า Finish gate กลับมา disabled ตอนเปลี่ยน lesson — kill แล้ว |
| REQ-7b | ลบ `if (!stillCurrent()) return;` ออกจาก `finish()` (guard ของ UX-5, **มีอยู่ก่อน ticket นี้ ไม่ใช่ regression จาก Q-1**) | **survived** — ไม่เคยมีเทสต์ pin guard นี้เลยตั้งแต่ UX-5 | เพิ่มเทสต์ใน `page.test.tsx`: deferred promise ควบคุมการ resolve ของ Finish PUT, กด Finish lesson A แล้วสลับไป lesson B ก่อน PUT resolve, resolve ทีหลัง (stale response) → assert lesson B ไม่ถูก mark finished — kill แล้ว |

**สรุป: 6 มูเทชันจากรอบ reviewer (REQ-1/2/5/6/7a/7b) — 5 survivor เป็น gap
จริงที่แก้แล้วด้วยเทสต์ใหม่และตายหมด (REQ-1/2/5/7a/7b), 1 เป็น equivalent
mutant ที่ตั้งใจไม่ kill (REQ-6). แยกต่างหาก: 2 defect จริงที่พบนอก
มูเทชัน (REQ-3, REQ-4) แก้ที่โค้ดโดยตรง — REQ-3 มีเทสต์ใหม่ pin ไว้ด้วย
(kill ได้จริงเมื่อ revert กลับไปเป็นรูปแบบรอบแรก, ดูแถว REQ-3 ด้านบน) แต่
**REQ-4 ไม่มีเทสต์ pin เลย** — ไม่มี test ไหนใน repo assert `type` ของปุ่มเลย
สักตัว, การแก้ยืนยันได้แค่จาก browser จริง (ดูหัวข้อ Playwright evidence)
ไม่ใช่จาก `npm test`
— รายละเอียดการแก้แต่ละจุดอยู่ในตารางด้านบน (ดูหัวข้อ "Re-verify เต็มชุด"
ด้านล่างสำหรับตัวเลขรวม)

### Re-verify เต็มชุด (หลังแก้ครบ round 2)

- `npm test`: **87 → 117 (รอบแรก) → 124 (หลังแก้ round 2)** — +7 จาก round 2:
  1 (REQ-1) + 3 (REQ-2, `it.each`) + 1 (REQ-3) + 1 (`page.test.tsx`, REQ-7b) +
  1 ไฟล์ใหม่ `page.finishedReset.test.tsx` (REQ-7a); REQ-5/REQ-6 ขยาย/เปลี่ยน
  ชื่อเทสต์เดิมโดยไม่เพิ่มจำนวน
- `npx tsc --noEmit`, `npx eslint .`, `npm run build` — สะอาดทั้งสามคำสั่งหลัง
  แก้ครบ
- `go vet ./...` + `go test ./...` — ผ่านทั้งหมด (ไม่มีไฟล์ `.go` เปลี่ยนเลย
  ใน ticket นี้)
- Re-run มูเทชันทั้ง 11 จุดจากรอบแรก **บนโค้ดหลังแก้** — ตายครบ 11/11 เหมือนเดิม
  (ยืนยันว่าแก้ REQ ต่าง ๆ ไม่ทำให้ coverage เดิมถอยหลัง) และ re-run REQ-1/2/5/7a/7b
  ทุกจุด — ตายครบตามตารางด้านบน; REQ-6 คงสถานะ equivalent (ไม่ตาย โดยตั้งใจ)

### Playwright evidence — round 2 (rebuilt stack after all REQ fixes)

`docker compose up -d --build web` (mysql/api ไม่ได้แก้ ข้อมูลเดิมยังอยู่ผ่าน
volume), ปิดท้ายด้วย `docker compose down` เปล่า ๆ — ไม่มี `-v`. Viewport
1440×900 เป็นค่าเริ่มต้นของ context (ตามที่ reviewer เรียกร้องให้วัด axe ที่
1440px ด้วย ไม่ใช่แค่ 375px), ทดสอบกับ `domain-driven-design/ubiquitous-language`
(มีทั้ง `short_answer` และ `mcq` checks จริงจาก DB):

- **Stage 1**: ดึง `expected_answer` ทุกตัวจาก `GET
  /api/v1/lessons/.../ubiquitous-language` มาเทียบกับ `page.content()`
  (serialized HTML เต็ม ไม่ใช่แค่ visible text) — **0 ตัวเลค**
- **Stage 2**: กด "Show options" แล้วเช็คว่ากล่อง "Answer" (reveal box)
  **ยังไม่ปรากฏ** (count = 0); ปุ่ม "Show options" และ "Reveal answer" ทั้งคู่
  มี `type="button"` จริง (REQ-4 ยืนยันจาก browser ไม่ใช่แค่อ่านโค้ด); กด
  "Reveal answer" ทั้งที่ยังไม่ commit อะไรเลย (`{force:true}` ข้าม visual
  disabled) → `role="status"` count ยังเป็น **0** (guard ในโค้ดกันไว้จริง
  ไม่ใช่แค่ `aria-disabled` มองไม่เห็น)
- **Keyboard-only end-to-end**: focus ตัวเลือกแรกด้วย `.focus()` แล้วกด
  `ArrowDown` → มี option ถูก check จริง (count 1); ทำนองเดียวกันกับ
  confidence ด้วย `ArrowRight`; focus ปุ่ม Reveal แล้วกด `Enter` → stage 3
  โผล่จริง ยืนยัน "Reveal reached via keyboard-only interaction: OK"
- **หลัง reveal (mcq) — role=status ที่ scope แล้ว (REQ-3)**: `expected_answer`
  โผล่ใน HTML จริง; `role="status"` count = **1 พอดี** (ไม่ใช่ทั้ง panel) และ
  `.innerText()` ของมันมีแค่ **"Correct"** คำเดียว (รอบที่สุ่มได้ตัวเลือกถูก) —
  ไม่รวม Answer box/option list; `document.activeElement` ตรวจแล้วเป็น
  **container ด้านนอก** (`data-testid="stage-reveal"`) ไม่ใช่ status line
  (ยืนยันว่า focus ทำหน้าที่ประกาศ reveal โดยรวม ส่วน status แคบแค่ผลลัพธ์)
  — `ariaSnapshot()` ยืนยันโครงสร้าง: `status: Correct` เป็น node แยกจาก
  `paragraph: Answer` และ `list` ของตัวเลือก ซึ่งเป็น sibling กัน ไม่ใช่ลูกของ
  status
- **short_answer stage 3 (REQ-3 อีกกรณี)**: `role="status"` count = **0**
  พอดี — ปุ่ม Pass/Not yet เป็น sibling ธรรมดา ไม่ถูกห่อด้วย live region เลย
  (`ariaSnapshot()` ยืนยัน `button "Pass"` / `button "Not yet"` เป็น node
  แยก ไม่อยู่ใต้ status ใด ๆ)
- **AX tree ที่ stage 2** (`ariaSnapshot()`): เห็น `group "Choose one"` และ
  `group "How confident are you?"` ห่อ `radio` ที่มี accessible name เป็น
  ข้อความตัวเลือกจริง (มาจาก `fieldset`/`legend` โดย browser แม็ปให้เอง) —
  ยืนยัน radiogroup semantics โดยไม่ต้องเขียน ARIA เอง
- **Contrast วัดจริงจาก browser** (`getComputedStyle` + composite ของทุก
  background layer ที่ element สืบทอด แปลง `oklab()` ที่ Tailwind v4's `/10`
  opacity utility serialize กลับมาเป็น sRGB ด้วยสูตร CSS Color 4 ก่อนคำนวณ
  relative luminance): รอบแรกสุ่มได้ตัวเลือกผิด → banner **"Not quite"**
  (fg `rgb(190,18,60)` บน bg `rgb(252,231,232)`) = **5.31:1**; รอบสอง
  (หลัง rebuild) สุ่มได้ตัวเลือกถูก → banner **"Correct"** (fg
  `rgb(4,120,87)` บน bg `rgb(230,243,236)`) = **4.80:1** (ค่าเดียวกับแถว
  ตัวเลือกถูกใน list เพราะใช้สีชุดเดียวกัน) — **ทั้งสองกรณีผ่าน AA (≥4.5:1)**
- **375px**: `scrollWidth === clientWidth === 375` จริง ไม่มี horizontal
  scroll
- **axe-core ที่ 1440px และ 375px ทั้งคู่ (รันทุก rule ไม่ scope เฉพาะ
  wcag2a/wcag2aa)**: **1 violation เท่ากันทั้งสอง viewport** —
  `scrollable-region-focusable` บน `.my-4.overflow-x-auto` ซึ่งเป็น
  `web/components/Mermaid.tsx:68` (**ไม่ใช่ `LessonBody.tsx`** ตามที่รอบแรก
  เขียนผิด — `LessonBody.tsx` ใช้ class `mt-4` ไม่ใช่ `my-4`; **ไฟล์นี้ ticket
  นี้ไม่ได้แตะเลย** — `git diff --name-only develop -- '*.tsx' '*.ts'` ยืนยัน
  ว่าไม่อยู่ในรายการไฟล์ที่แก้ — pre-existing บน `develop` ก่อน ticket นี้, นอก
  scope) — **0 violations** จาก `RecallCheckCard` เอง
- **`Button.tsx`'s `type="button"` default ไม่กระทบหน้าอื่น**: smoke-test
  `/today`, `/learn`, `/learn?track=domain-driven-design` — ปุ่มทุกตัวได้
  `type="button"` จริง, คลิกปุ่ม "Set focus" บน `/learn` แล้วไม่ crash;
  grep `<form` ทั้ง repo เจอที่เดียวคือ `app/token/page.tsx` ซึ่งตั้ง
  `type="submit"` ชัดเจนอยู่แล้ว (override default ได้ตามที่ออกแบบไว้)
- Console error / page error = 0 ตลอดการทดสอบทั้งสองรอบ

### Review focus

- ทำไม `selectedIndex` ต้องเป็น index ของ shuffled array แทนที่จะเก็บข้อความ
  ตัวเลือกตรง ๆ เหมือน `rating` เดิม?
- ทำไม confidence ต้องถูกเก็บ**ก่อน**กด Reveal เสมอ แทนที่จะถามพร้อมกับ
  Pass/Not yet ตอน stage 3?
- ทำไมการมูเทต `key={check.position}` กลับเป็นค่าเดิม (ไม่ใส่ lesson identity
  เข้าไป) ถึงไม่ทำให้เทสต์ตาย ทั้งที่ ticket เตือนเรื่อง per-check state reset
  ไว้ตรง ๆ?
- ทำไม `role="status"` ต้องห่อแค่บรรทัด Correct/Not quite เดียว แทนที่จะห่อ
  ทั้ง reveal panel เหมือนตอนแรก ทั้งที่ short_answer ก็ต้องการให้ผลลัพธ์ถูก
  ประกาศเหมือนกัน?

### จงใจไม่ทำในรอบนี้

- ไม่ persist อะไรเพิ่มจากเดิม — `selectedIndex`/`confidence`/
  `shortAnswerRating` อยู่ใน memory ของ `RecallCheckCard` เท่านั้น เหมือนที่
  `ratings` เดิมก็อยู่แค่ใน `lesson/page.tsx` — Q-2 คือ ticket ที่ทำให้ทุกอย่าง
  persist จริง
- ไม่แก้ `scrollable-region-focusable` ของ `web/components/Mermaid.tsx` —
  pre-existing, ไม่เกี่ยวกับไฟล์ที่ ticket นี้แตะ
- ไม่ตัด `expected_answer` ออกจาก payload — ตัดสินใจแล้วว่าไม่ทำใน v1 (ดู
  หัวข้อการตัดสินใจด้านบน)

Status: implemented, round 2 fixes applied post code-review, PR pending

## Q-2a — persist recall attempts `[go-implementer]`

- **Scope**: schema ใหม่ `migrations/006_recall.sql` (`recall_attempts` เท่านั้น —
  `review_cards`/`review_logs` เป็น Q-2c, ไม่ใช่รอบนี้), domain VOs 5 ไฟล์ใหม่ใน
  `internal/learning/domain/` (`checkkey.go`, `confidence.go`,
  `attemptoutcome.go`, `gradedby.go`, `recallattempt.go` — ไฟล์หลังมี `CheckKind`
  อยู่ด้วย), `internal/learning/app/service.go` เพิ่ม `RecordAttempt` +
  `AttemptRecord` DTO (มี field `question` คืน canonical text ด้วย) +
  `ErrCheckNotInLesson`/`ErrInvalidAttempt` + `Repository.RecordAttempt` (resolve
  ก่อนแล้วเรียก `build` closure ให้ app layer สร้าง `domain.RecallAttempt` จริง —
  ดูเหตุผลที่หัวข้อ R1 ด้านล่าง), `internal/learning/infra/repository.go` เพิ่ม SQL
  3 statement + `resolveRecallCheck` (มิเรอร์ `lockLesson`'s ambiguous-match guard)
  + `Repository.RecordAttempt`, `internal/learning/infra/handler.go` เพิ่ม route
  `POST /api/v1/progress/{topic}/{concept}/attempts`. Backend ล้วน —
  `git diff --name-only develop... -- 'web/*'` ว่างเปล่าจริง, ไม่แตะ
  `web/components/RecallCheckCard.tsx` เลย (นั่นคือ Q-2b). **รอบ 2 (หลัง code
  review)**: เพิ่มไฟล์ที่ 6 `canonicalquestion.go` (`domain.CanonicalQuestion`,
  R3), `kind` ถูกตัดออกจาก request body ทั้งหมด (ตอนนี้ resolve จาก
  `recall_checks.type` เสมอ, R1), เพิ่มเทสต์ exhaustive-set 4 ตัวสำหรับ enum
  ทั้งสี่ (R2), และ `stub_driver_test.go` โมเดล `recall_checks.type` +
  ENUM validity ของ `recall_attempts` เพิ่ม (R4) — รายละเอียดทั้งหมดอยู่ที่หัวข้อ
  "รอบ code-reviewer 2" ด้านล่าง

### การตัดสินใจหลัก

- **`check_key = SHA256(topic + "/" + concept + "/" + trimmed-question)` hex
  lowercase, ไม่ใช้ FK ไป `recall_checks.id`, ไม่ใช้ `(lesson_id, position)`** —
  ทั้งสามทางเลือกถูกพิจารณาจริง ไม่ใช่แค่หยิบ content-addressing มาเฉย ๆ:
  - **ทำไมไม่ใช้ FK ไป `recall_checks.id`**:
    `internal/curriculum/infra/lessonwriter.go` ทำ `DELETE FROM recall_checks
    WHERE lesson_id = ?` แล้ว insert ใหม่ทั้งหมดทุกครั้งที่ import เนื้อหาซ้ำ
    (`deleteRecallChecksSQL`/`insertRecallCheckSQL`) — id พวกนี้**ไม่เสถียร**ข้าม
    การ import แต่ละรอบ ต่างจาก `lessons.id` ที่เสถียรจริงผ่าน
    `ON DUPLICATE KEY UPDATE id = LAST_INSERT_ID(id)` ใน `upsertLessonSQL`
  - **ทำไมไม่ใช้ `(lesson_id, position)`**: การ reorder คำถามของ lesson เดียวกัน
    ตอน re-import (เช่น สลับข้อ 2 กับข้อ 3) จะทำให้ position เดิมชี้ไปคำถามคนละ
    ข้อ — ประวัติ attempt เก่าจะไปติดกับคำถามผิดข้อแบบเงียบ ๆ โดยไม่มีใครรู้
  - **content-addressing แก้ปัญหาแบบ fail-safe**: แก้คำถามแม้แต่ตัวเดียว = key
    ใหม่ = ประวัติเริ่มใหม่ (เสียประวัติเก่าไปแต่ไม่ผิดที่ผิดทาง) ดีกว่า
    misattach ประวัติไปคำถามอื่นแบบไม่รู้ตัว
  - **hash คำถามที่ trim แล้วเท่านั้น ไม่ทำ normalization อื่นเพิ่ม**:
    `internal/curriculum/domain/recallcheck.go`'s `NewRecallCheck` เรียก
    `validateRecallQuestion` ที่ trim คำถามไว้ก่อนเก็บลง DB อยู่แล้ว ดังนั้น
    `recall_checks.question` ใน DB คือ trimmed เสมอ — hash รูปแบบ trimmed ให้
    ตรงกับสิ่งที่ DB เก็บจริงไบต์ต่อไบต์ **ตัดสินใจแล้ว: trim เท่านั้น ไม่ case-fold
    ไม่ unicode-normalize เพิ่ม** เพราะกฎ normalization เพิ่มเติมใด ๆ จะกลาย
    เป็น source of truth ที่สองที่ drift จากสิ่งที่ curriculum เก็บจริงได้ โดยไม่ได้
    ประโยชน์อะไรเพิ่ม (คำถามเป็น authored content ไม่ใช่ user input ที่พิมพ์เพี้ยน
    case/space ได้ง่าย)
- **`check_key` collation: `CHAR(64) CHARACTER SET ascii COLLATE ascii_bin`**
  ไม่ใช้ table default (`utf8mb4_0900_ai_ci`) — column นี้เก็บ lowercase hex ที่
  ฝั่ง Go เป็นคนสร้างเองเสมอ (ไม่ใช่ user-facing text ที่ต้องการ accent/case
  folding) การเทียบแบบ byte-exact คือสิ่งที่ hash column ต้องการจริง และ ascii
  ใช้พื้นที่ครึ่งเดียวของ utf8mb4 ต่อ index ที่จะถูก query ทุกครั้งที่ Q-2c ทำ SRS
  review
- **Index `(check_key, created_at DESC)`**: ตรงกับ query pattern ของ Q-2c
  (`WHERE check_key = ? ORDER BY created_at DESC`) แม้ query นี้ยังไม่ถูกเขียนใน
  รอบนี้ก็ตาม — เตรียม index ไว้ตอนสร้างตาราง ถูกกว่าเพิ่มทีหลัง
- **`kind` เป็น field ที่ server resolve จาก `recall_checks.type` เอง — ไม่ใช่
  client ส่งมา (แก้ในรอบ 2, ดู R1 ในหัวข้อ "รอบ code-reviewer 2" ด้านล่าง;
  ย่อหน้านี้เคยเขียนตรงข้ามในรอบแรกและถูกพิสูจน์ว่าผิดจริง)**: request body มีแค่
  `question`/`confidence`/`outcome`/`selected_option` — server เชื่อ client
  เฉพาะเนื้อหาที่เป็น**การตัดสินใจของ user จริง** (`outcome` self-graded,
  `confidence` self-reported, `selected_option` คือสิ่งที่กดจริง) ส่วน `kind`
  เป็น curriculum content ที่ server เป็นเจ้าของอยู่แล้วผ่าน `recall_checks.type`
  ไม่ใช่ user judgement เหมือน `outcome` เลย — เชื่อ client เรื่องนี้จะเปิดช่องให้
  `type='mcq'` กับ `selected_option` ปลอมไปเก็บใต้คำถาม `short_answer` จริงได้
  (พิสูจน์แล้วจริงในรอบ 2). server ต้อง**ยืนยันเอง**สองเรื่อง: "คำถามนี้เป็นของ
  lesson นี้จริงไหม" (identity/scoping) และ "คำถามนี้ kind อะไรจริง" (content,
  ไม่ใช่ identity แต่ก็ไม่ใช่ user judgement)
- **Grading ยังเป็น self-graded**: `graded_by` เก็บต่อแถว default เป็น `'self'`
  — `domain.GradedBySelf` เป็นค่าเดียวที่ `Service.RecordAttempt` ส่งเข้า
  `domain.NewRecallAttempt` จริง ๆ (enum เองมีทั้ง `self`/`llm` เผื่ออนาคต แต่
  `'llm'` **เข้าไม่ถึงได้เลยผ่าน endpoint นี้** — ไม่ใช่ analogy เดียวกับ
  `RecallKind` ที่ทั้งสองค่าถูกใช้งานจริง). Server **ไม่**คำนวณถูก/ผิดเองจาก
  `curriculum`'s `expected_answer` — client ส่ง `outcome` มาตรง ๆ ตามที่ Q-1
  ออกแบบไว้ (mcq เทียบข้อความฝั่ง client, short_answer self-rate) และ
  `internal/learning` **ไม่ import** `internal/curriculum`'s domain types เพื่อ
  re-derive อะไรเลย (ผิด bounded-context boundary ตามที่
  `internal/learning/domain/lessonref.go`'s doc comment อธิบายไว้)
- **"คำถามไม่ใช่ของ lesson นี้" → 400 ไม่ใช่ 404**: ตาม precedent เดิมใน
  `internal/learning/app/service.go` — `ErrLessonNotFound` (404) สงวนไว้กับ
  "ไม่รู้จัก topic/concept นี้เลย" (ตัว resource หลักไม่มีอยู่จริง),
  `ErrInvalidSlug`/`ErrInvalidState` (400) คือ "input ที่ client ส่งมาผิดรูป"
  คำถามที่ไม่ match `recall_checks` ของ lesson **มีอยู่จริง** (lesson หาเจอ) แค่
  client อ้างอิงผิด (typo/ cache เก่า) — เข้าเงื่อนไข "malformed reference" มากกว่า
  "resource ไม่มีอยู่" จึงเลือก 400 (`ErrCheckNotInLesson`) สอดคล้องกับ
  `ErrInvalidSlug` มากกว่า `ErrLessonNotFound`
- **Endpoint path: `POST /api/v1/progress/{topic}/{concept}/attempts`** — nest
  ใต้ resource `progress` เดิมแทนที่จะเป็น flat resource ใหม่ เพราะ attempt
  ผูกกับ lesson หนึ่งเดียวเสมอ (topic+concept คือ identity ของมัน เหมือนที่
  ticket บังคับให้ "คำถามต้องเป็นของ lesson นี้") — path param แบบเดียวกับ
  `PUT /api/v1/progress/{topic}/{concept}` ที่มีอยู่แล้ว และ `recall_attempts`
  ก็เป็น user-state ของ learning bounded context เหมือน `lesson_progress` ไม่ใช่
  content ที่ `curriculum` เป็นเจ้าของ
- **`RecordAttempt` เป็น resolve-then-insert ไม่ใช่ locked read-modify-write
  แบบ `Transition`**: `recall_attempts` เป็น append-only ไม่มี state ให้ race
  แก้ไขซ้ำ (ต่างจาก `lesson_progress` ที่ writer สองตัวแข่งกันแก้ค่าเดียวกันได้)
  — งานเดียวที่ต้องพึ่ง DB จริง ๆ คือ resolve ว่าคำถามเป็นของ lesson นี้ก่อน
  insert จึงไม่ต้องเปิด transaction/lock แบบ `selectLessonForUpdateSQL`

### Mutation table

รัน `go test ./internal/learning/...` แบบ quiet หลังแก้แต่ละจุด แล้ว revert ทุกครั้ง:

| # | Mutation | ผลลัพธ์ |
|---|---|---|
| 1 | สลับลำดับ hash เป็น `concept + "/" + topic + "/" + question` | killed — `TestNewCheckKey` (เทียบ hash จริงจาก `crypto/sha256` คำนวณแยก) + `TestNewCheckKey_HashOrderMatters` |
| 2 | hash คำถามแบบไม่ trim (ใช้ `question` ดิบแทน `trimmed`) | killed — `TestNewCanonicalQuestion` (`trims leading/trailing whitespace` + `whitespace-only` cases) — trimming ย้ายมาอยู่ที่ `NewCanonicalQuestion` ในรอบ 2 (R3), ไม่ใช่ `NewCheckKey` อีกต่อไป; `TestNewCheckKey_TrimOnly` เดิมถูกแทนที่ด้วย `TestNewCheckKey_UsesQuestionVerbatim` (เทสต์คนละเรื่อง — ยืนยันว่า `NewCheckKey` ไม่ normalize ซ้ำ ไม่ใช่เรื่อง trim) |
| 3 | ตัดเช็ค "คำถามเป็นของ lesson นี้" ทิ้ง (ใช้ `selectLessonExistsSQL` แทน `selectRecallCheckExistsSQL` ใน `RecordAttempt`) | killed — `TestRepositoryRecordAttempt_CheckNotInLesson` + `TestRepositoryRecordAttempt_ScopedByBothSlugs` |
| 4 | `AND` → `OR` ใน `selectRecallCheckExistsSQL`'s WHERE clause | killed — `TestSelectRecallCheckExistsSQLShape` (literal-string pin; stub dispatch ด้วย query-string identity เอง detect ไม่ได้ ตาม precedent เดิมของไฟล์นี้) |
| 5 | เติม `ON DUPLICATE KEY UPDATE created_at = VALUES(created_at)` ใน `insertRecallAttemptSQL` | killed — `TestInsertRecallAttemptSQLShape` (literal-string pin เท่านั้น — behavioral test `TestRepositoryRecordAttempt_AppendOnly` จับไม่ได้เพราะ stub append เข้า slice เสมอไม่สนใจ SQL text จริง ยืนยันด้วย live MySQL evidence ด้านล่างแทน) |
| 6 | hardcode `graded_by` เป็น `"self"` ใน exec call แทนที่จะใช้ `attempt.GradedBy().String()` | killed **ที่ repository layer เท่านั้น** — `TestRepositoryRecordAttempt_GradedByThreadsThrough` (สร้าง `domain.RecallAttempt` ด้วย `GradedBy=llm` ตรง ๆ ข้าม app layer) — ที่ app/HTTP layer มูเทชันนี้ **equivalent จริง** เพราะ `Service.RecordAttempt` ส่งแค่ `domain.GradedBySelf` เข้ามาเสมอวันนี้ (ยืนยัน: app-layer tests ทั้งชุดยังเขียวหลังมูเทต) |
| 7a | ลบ `"confident"` ออกจาก `confidences` array | killed — `TestNewConfidence` + `TestNewRecallAttempt` + `TestServiceRecordAttempt_EchoesSubmittedValues` + `TestHandlerPostAttempt_Success` + `TestRepositoryRecordAttempt_Success` |
| 7b | เติมค่าปลอม `"partial"` เข้า `attemptOutcomes` array | killed — `TestNewAttemptOutcome` (`unknown value` case ที่คาด error แต่ไม่ error) + `TestServiceRecordAttempt_InvalidFields` |
| 7c | ลบ `"mcq"` ออกจาก `checkKinds` array | killed — `TestNewCheckKind` + `TestNewRecallAttempt` + `TestServiceRecordAttempt_EchoesSubmittedValues`/`_EmptySelectedOptionTreatedAsAbsent` + `TestHandlerPostAttempt_Success` + `TestRepositoryRecordAttempt_Success` |
| 8 | hardcode confidence เป็น `domain.NewConfidence("guessed")` ใน `Service.RecordAttempt` แทนที่จะใช้ `rawConfidence` ที่ submit มา | killed — `TestServiceRecordAttempt_EchoesSubmittedValues` (assert ค่าที่ submit ตรงกับที่ได้กลับมา ตรงจุดที่ Q-1's frontend review เจอ bug class นี้) + `TestServiceRecordAttempt_InvalidFields` |
| 9 | ลบเช็ค `selectedOption != nil && !kind.IsMCQ()` ออกจาก `NewRecallAttempt` | killed — `TestNewRecallAttempt` (`selected option on short_answer is rejected` case) + `TestServiceRecordAttempt_SelectedOptionOnShortAnswerRejected` |
| 10 | ลบ ambiguous-match guard ออกจาก `resolveRecallCheck` (`QueryContext`+loop → `QueryRowContext` ตัวเดียว, ตาม R1's code review round) | killed — `TestRepositoryRecordAttempt_AmbiguousCaseVariantRejected` (seed 2 recall_checks ที่ fold เป็นคำถามเดียวกันใน lesson เดียวกัน ยืนยันว่า error ไม่ใช่แค่เลือกแถวแรกเงียบ ๆ) |
| 11 | ลบ `.Truncate(time.Second)` ออกจาก `now` ก่อน insert (R2's fix) | killed — `TestRepositoryRecordAttempt_CreatedAtTruncatedToSeconds` (assert `CreatedAt.Nanosecond() == 0`) |
| 12 | `Repository.RecordAttempt` ส่ง kind ที่ hardcode (`"short_answer"`) เข้า `build` แทน kind ที่ resolve จาก `rc.type` จริง (round 2, R1's kind fix) | killed — `TestRepositoryRecordAttempt_KindResolvedFromRecallChecksType` (ใหม่ในรอบ 2, seed สอง lesson คนละ kind ยืนยันว่า kind ที่ได้กลับมาตรงกับที่ seed เสมอ) |

**สรุปรอบแรก (retracted บางส่วนในรอบ 2 — ดูด้านล่าง): เดิมเขียนว่า "11/11 มูเทชันตายหมด
ไม่มี survivor ที่เป็น coverage gap จริง"** — ประโยคนี้**เกินจริง** พบใน code review
รอบ 2 (นี่เป็น ticket ที่สองแล้วที่ claim ความครบของ mutation testing เกินจริงกว่าที่
ทดสอบจริง ดู Q-1's ticket สำหรับครั้งแรก): #7a/#7b/#7c เดิมทดสอบด้วยค่าปลอมที่
**บังเอิญตรงกับ literal ที่ test table ใช้เป็น "unknown value" case อยู่แล้ว**
(`"partial"` สำหรับ outcome, ทำนองเดียวกันกับ confidence/kind/graded_by) — เทสต์
เดิมไม่เคย assert ว่า **enum ทั้งชุด** ตรงกับ literal list ที่กำหนดไว้ แค่ assert
ว่าค่าที่ทดลองสามค่าถูกปฏิเสธ. ทดสอบด้วยค่าปลอม**คนละตัว**ที่ไม่ตรงกับ literal เดิม
(`"maybe"`/`"essay"`/`"skipped"`/`"robot"` — เพิ่มเข้า `confidences`/`checkKinds`/
`attemptOutcomes`/`gradedBys` ทีละตัว) ยืนยันว่า **ทั้งสี่ค่าปลอมนี้ survive เทสต์
เดิมทั้ง 266 ข้อในรอบนั้น** — แก้แล้วในรอบ 2 (ดู R2 ด้านล่าง) ด้วยเทสต์ใหม่ที่เทียบ
array ทั้งชุดกับ literal slice ตรง ๆ (`TestConfidenceAcceptedSetIsExactly` และเทสต์
คู่กันอีกสามตัว) — ทดสอบซ้ำด้วยค่าปลอมชุดเดียวกันหลังแก้: **ตายครบทั้งสี่**.
มูเทชัน #6 (graded_by hardcode) ยังเป็น equivalent mutant ที่ **ตั้งใจ**ปล่อยไว้ที่
app layer เหมือนเดิม (เพราะ `'llm'` เข้าไม่ถึงได้จริงผ่าน endpoint วันนี้ตามการ
ออกแบบ) แต่มี test จริงที่ repository layer กัน regression ไว้ล่วงหน้าสำหรับตอนที่
Q-2c/LLM adapter มาเสียบจริง.

**หมายเหตุเรื่อง #4/#5 (SQL literal-pin mutations)**: ทั้งสอง kill ได้เพราะ
`TestSelectRecallCheckExistsSQLShape`/`TestInsertRecallAttemptSQLShape` เทียบ
literal string ตรง ๆ — **ถ้าใครแก้ SQL source พร้อมแก้ `want` literal ในเทสต์ให้
ตรงกันไปด้วย (ไม่ใช่แค่แก้ source อย่างเดียว) มูเทชันนี้จะไม่ตายเลย** ตารางด้านบน
บันทึกจุดนี้ไว้ตรง ๆ แทนที่จะให้อ่านว่าเป็น kill แบบไม่มีเงื่อนไข — ป้องกันด้วย
review เท่านั้น (คนแก้ SQL + เทสต์คู่กันในคอมมิตเดียวควรโดนจับตอน diff review)
ไม่มีกลไกอัตโนมัติกันจุดนี้ได้

**ยืนยันซ้ำหลังรอบ 2 (R1's kind fix + CanonicalQuestion refactor)**: มูเทชัน
#1/#2/#3/#4/#5/#6/#9/#10/#11 (ทุกจุดที่โค้ดที่เกี่ยวข้องถูก refactor ในรอบ 2) ถูก
re-run ทีละจุดอีกครั้งบนโค้ดหลังแก้ครบ — ตายเหมือนเดิมทุกจุด ไม่มี regression จาก
การ refactor เป็น `CanonicalQuestion` type และการตัด `kind` ออกจาก request

### รอบ code-reviewer (CHANGES NEEDED → แก้ครบ)

Reviewer proved บน `mysql:8.4` จริงว่า `recall_checks.question` (TEXT ใต้
`utf8mb4_0900_ai_ci`) match แบบ case/accent-insensitive — `what is a b-tree?`
กับ Thai case-variant ต่างก็ match แถวเดิมที่เก็บเป็น `What is a B-tree?`/
`Ubiquitous Language...` ได้จริง แต่ก่อนแก้ code hash คำถามที่ **client ส่งมา**
(`rawQuestion`) ไม่ใช่คำถาม**ที่ query จริง ๆ match ได้** — ทำให้คำถามที่ผ่านเช็ค
"เป็นของ lesson นี้" ได้ (เพราะ collation ยอมให้ match) แต่ hash ออกมาคนละค่ากับ
canonical text จริง = fork ประวัติ SRS ของคำถามเดียวกันแบบเงียบ ๆ — เป็นความ
บกพร่องที่ขัดกับเหตุผลหลักที่ `checkkey.go` มีอยู่ (content-addressing ต้องชี้
กลับไปที่เนื้อหาเดียวกันเสมอ ไม่ว่า caller จะพิมพ์ด้วย casing ไหน)

- **R1 (blocking, แก้แล้ว) — ต้อง hash คำถาม canonical จาก DB ไม่ใช่จาก client**:
  `selectRecallCheckExistsSQL` เปลี่ยนจาก `SELECT l.id` เป็น
  `SELECT l.id, rc.question` — คืนคำถามที่ query จริง ๆ match ได้กลับมาด้วย. ผล
  คือ `Repository.RecordAttempt` เปลี่ยนจากรับ `domain.RecallAttempt` สำเร็จรูป
  มาเป็นรับ **`build func(canonicalQuestion string) (domain.RecallAttempt,
  error)`** แทน — resolve lesson+question ก่อน แล้วค่อยเรียก `build` ด้วย
  canonical text ที่ query คืนมา (ไม่ใช่ resolve-then-key แบบเดิมที่คำนวณ
  `CheckKey` ตั้งแต่ก่อนรู้ว่า DB match ได้ด้วยข้อความอะไร) — เลือก pattern นี้
  (ไม่ใช่ two-step port แยก resolve/insert) เพราะตรงกับ precedent ของ
  `Transition`'s `decide` callback ที่มีอยู่แล้วในไฟล์เดียวกัน: infra resolve ให้
  ก่อน แล้วให้ app layer logic (build ที่ปิดด้วย closure) ตัดสินใจ/สร้าง object
  จริงโดยใช้ข้อมูลที่ resolve มาได้ ไม่ต้องเปิด round trip ที่สอง. `checkkey.go`'s
  doc comment ที่เคยอ้างว่า "hashing the trimmed form matches what MySQL
  actually stores byte for byte" เดิมเป็นเท็จ (โค้ดไม่ได้บังคับ invariant นี้จริง)
  แก้ให้บอกตรง ๆ ว่า NewCheckKey เองพิสูจน์ไม่ได้ว่า input เป็น canonical —
  เป็นหน้าที่ของ caller (`Service.RecordAttempt`) ที่ต้องส่ง canonical text จาก
  repository เท่านั้น ไม่ใช่ raw client input. Test fake เดิม (`stub_driver_test.go`,
  `service_test.go`) เข้มกว่า MySQL จริงตรงจุดนี้พอดี (exact-match lookup ทั้งที่
  slug ใช้ `foldKey` case-fold อยู่แล้ว) — แก้ให้ fold คำถามด้วย (`foldQuestion`)
  แล้วเพิ่มเทสต์ยืนยันตรง ๆ ว่า case-variant submission ได้ `check_key` **เดียวกัน**
  กับ canonical (`TestRepositoryRecordAttempt_CaseVariantResolvesToCanonicalCheckKey`
  ที่ repository layer, `TestServiceRecordAttempt_CheckKeyDerivedFromCanonicalQuestion`
  ที่ service layer — สองระดับ ไม่ใช่แค่ระดับเดียว). ผลข้างเคียงที่ได้มาฟรี:
  `AttemptRecord`/response DTO เพิ่ม field `question` (canonical text ที่ resolve
  ได้ ไม่ใช่ข้อความดิบที่ client ส่งมา) — client ที่ส่งคำถามแบบ case/accent-variant
  จะเห็นว่าถูก normalize แล้วจาก response ตรง ๆ แทนที่จะเดาเอาว่า `check_key`
  ตรงกับ bytes ของตัวเองหรือเปล่า
- **R2 (blocking, แก้แล้ว) — `created_at` ที่ echo กลับอาจไม่ตรงกับแถวจริงใน DB**:
  `time.Now().UTC()` มี sub-second precision, MySQL's `TIMESTAMP` (fsp 0) ปัดครึ่ง
  ขึ้น (round half-up) ตอนเก็บ, แต่ `handler.go` format ด้วย RFC3339 ซึ่ง**ตัดทิ้ง**
  (truncate) — สำหรับ request ที่มาถึงตอน ≥.500s ค่าที่ API ตอบจะช้ากว่าแถวจริงใน
  DB 1 วินาที. แก้ด้วยวิธีที่ diff เล็กที่สุด: `now :=
  time.Now().UTC().Truncate(time.Second)` **ก่อน**ส่งเข้า `insertRecallAttemptSQL`
  (ไม่ใช่แค่ตอน format คืน) — ทำให้ค่าที่ Go เก็บกับค่าที่ MySQL เก็บตรงกันโดย
  โครงสร้าง ไม่ต้องพึ่งการปัดของ MySQL เลย. เพิ่มเทสต์
  `TestRepositoryRecordAttempt_CreatedAtTruncatedToSeconds` ยืนยันทั้ง
  `CreatedAt.Nanosecond() == 0` และค่าที่ echo กลับตรงกับแถวที่เก็บจริงเป๊ะ
- **R3 (blocking, แก้แล้ว) — WHAT-comment**: ลบ comment บน
  `RecallAttempt.SelectedOption()` ที่แค่พูดซ้ำ signature ("returns the mcq
  option the user picked, or nil for a short_answer attempt") ทิ้ง — เหตุผลที่
  ไม่ชัดจากชื่อ (ทำไม nil ถึง valid เฉพาะ short_answer) มีอยู่แล้วที่ constructor
- **S1 (แก้แล้ว) — immutability ของ `selectedOption`**: `*string` เดิมแชร์
  pointer กับ caller/getter ตรง ๆ — reviewer พิสูจน์ด้วย probe จริงว่า mutate ผ่าน
  pointer เดิมได้ ทำให้ aggregate ไม่ immutable เหมือน VO อื่นในแพ็กเกจนี้ แก้ด้วย
  `copyStringPtr` (copy ทั้งขาเข้าตอน construct และขาออกตอน getter) เพิ่มเทสต์
  `TestRecallAttempt_SelectedOptionIsImmutable` ยืนยันทั้งสองทิศทาง
- **S2 (แก้เกินกว่าที่เสนอ — ปิด gap จริง ไม่ใช่แค่แก้คำอธิบาย)**: reviewer เสนอแค่
  "แก้เหตุผลที่เขียนผิดในเอกสาร" (ไม่บังคับแก้โค้ด เพราะ `check_key` ไม่ใช่
  `lesson_id` คือ identity จริงที่ SRS ใช้อยู่แล้ว) แต่ระหว่างแก้ R1's resolve step
  logic ก็ถูกดึงออกมาเป็นฟังก์ชันแยกชื่อ `resolveRecallCheck` (มิเรอร์ `lockLesson`
  เป๊ะ ๆ: `QueryContext` + loop ที่ error ทันทีถ้า match มากกว่าหนึ่งแถว แทนที่จะ
  เลือกแถวแรกเงียบ ๆ แบบ `QueryRowContext`) — ปิด gap ที่เอกสารรอบแรกอ้างว่า
  "จงใจไม่ทำ" ไปเลย ไม่ใช่แค่แก้คำอธิบายให้ตรงกับพฤติกรรมเดิม เพิ่มเทสต์
  `TestRepositoryRecordAttempt_AmbiguousCaseVariantRejected` ยืนยัน (mutation #10)
- **S4 (แก้แล้ว) — `maxAttemptBodyBytes` comment เกินจริง**: จาก 4096 พร้อม comment
  ที่ไม่ได้อ้างตัวเลขจริง เปลี่ยนเป็น 8192 พร้อม comment ที่อ้าง
  `maxRecallQuestionRunes`/`maxRecallOptionRunes` (1000/255 runes,
  `internal/curriculum/domain/text.go`) ตรง ๆ — unexported อ้างชื่อได้แต่ import
  ไม่ได้ (คนละ bounded context)
- **S5 (แก้แล้ว) — test gap เล็ก ๆ**: เพิ่ม
  `TestServiceRecordAttempt_EmptyQuestionRejected` (question ว่าง/whitespace →
  `ErrInvalidAttempt` ก่อนถึง repository); `checkkey_test.go`'s
  `expectedCheckKeyHex` เปลี่ยนจาก literal ตายตัวมาคำนวณจาก `tt.topic`/
  `tt.concept`/`tt.question` ของแต่ละแถวเอง; `repository_test.go`'s
  `mustRecallAttempt` เปลี่ยนมาเรียก `mustCheckKind`/`mustConfidence`/
  `mustAttemptOutcome` ของตัวเองแทนที่จะ inline ซ้ำ
- **S3 (ตั้งใจไม่แก้)**: reviewer เสนอให้ query เดียวกันที่แก้เพื่อ R1 select
  `rc.type` มาด้วยแล้ว cross-check กับ `kind` ที่ client ส่งมา — **ไม่ทำในรอบนี้**
  เพราะเป็นการ**พลิกการตัดสินใจที่ตั้งใจไว้แล้ว** ("kind เป็น field ที่ client ส่งมา
  ไม่ใช่ server อ่านจาก DB" ด้านบน) ไม่ใช่แค่ bug fix เล็ก ๆ — ถ้าจะทำจริงควรเป็น
  การตัดสินใจแยกต่างหาก ไม่ใช่ทำแทรกในรอบแก้ review

Mutation ที่ได้รับผลกระทบจาก R1's restructure (repository/service ทั้งคู่เปลี่ยน
signature เป็น `build` closure) ถูก re-run ครบหลังแก้ — รายละเอียดอยู่ที่ตาราง
ด้านบนแล้ว (อัปเดตแล้วให้ตรงกับโค้ดหลังแก้), ผลลัพธ์เดิมทั้งหมดยังตายเหมือนเดิม
ไม่มี regression จากการ refactor

### รอบ code-reviewer 2 (REQUEST_CHANGES → แก้ครบ)

Reviewer พิสูจน์สดว่า `kind` ที่ client ส่งมาเป็นช่องโหว่จริง ไม่ใช่การตัดสินใจ
ที่ปิดเรียบร้อยแล้วอย่างที่รอบแรกเขียนไว้ (S3 ด้านบน) — **S3's decision ถูก
พลิกในรอบนี้**: `recall_checks.id=5270` (`type='short_answer'`, `options`
NULL) ยิง POST ด้วย `kind:"mcq"` + `selected_option` ที่กุขึ้นเอง → **201**,
เก็บเป็น `type='mcq'` ใต้ `check_key` เดียวกับแถว `short_answer` ข้างเคียง —
invariant ที่ตรวจกับ input ที่ client คุมได้ไม่ใช่ invariant จริง (`kind` ผ่านเข้า
`NewRecallAttempt`'s เช็ค "selected option เฉพาะ mcq" แต่ `kind` เองมาจาก
client ที่โกหกได้)

- **R1 (blocking, แก้แล้ว) — ตัด `kind` ออกจาก request body ทั้งหมด, server
  resolve จาก `recall_checks.type`**: เหตุผลที่ "client ส่ง kind เอง" ใช้ไม่ได้
  อีกต่อไป — `outcome` เชื่อ client ได้เพราะ self-graded (คนตอบ = คนให้คะแนน)
  แต่ `kind` เป็น curriculum content ที่ server เป็นเจ้าของอยู่แล้ว ไม่ใช่การ
  ตัดสินใจของ user. แก้ตาม pattern เดียวกับที่ R1 (รอบแรก) แก้ให้ `question`:
  `selectRecallCheckExistsSQL` เพิ่ม `rc.type` เข้า SELECT (query เดียวกับที่คืน
  `rc.question` อยู่แล้ว ไม่มี round trip เพิ่ม), `resolveRecallCheck` คืน
  `kindRaw` เพิ่ม, `Repository.RecordAttempt`'s `build` closure เปลี่ยน
  signature เป็น `func(domain.CanonicalQuestion, domain.CheckKind)
  (domain.RecallAttempt, error)` — kind เป็น parameter ที่ resolve แล้วเหมือน
  question. `Service.RecordAttempt` ตัด `rawKind` param ออกทั้งหมด (6 params
  เหลือแทน 7), `recordAttemptRequest`'s `Kind` field ถูกลบ — client ที่ส่ง
  `"kind":"..."` มาจะโดน `DisallowUnknownFields()` ปฏิเสธเป็น 400 เหมือน field
  แปลกอื่น ๆ (ยืนยันด้วย `TestHandlerPostAttempt_KindFieldRejected` ใหม่, และ
  curl สดด้านล่าง) — ทำให้ "kind ไม่ตรงกับ recall_checks.type" **เป็นไปไม่ได้
  โดยโครงสร้าง** ไม่ใช่แค่ถูก validate แล้วปฏิเสธ
- **R2 (blocking, แก้แล้ว) — ไม่มีเทสต์ pin accepted-set ของ enum ทั้งชุด, claim
  "11/11 killed" เดิมเป็น coincidence**: ดูตาราง mutation ด้านบน (ย่อหน้า
  "สรุปรอบแรก") สำหรับรายละเอียดเต็ม — สรุปสั้น: เพิ่มเทสต์ใหม่ 1 ตัวต่อ enum
  (`TestConfidenceAcceptedSetIsExactly`, `TestCheckKindAcceptedSetIsExactly`,
  `TestAttemptOutcomeAcceptedSetIsExactly`, `TestGradedByAcceptedSetIsExactly`)
  เทียบ array ทั้งชุดกับ `[]string` literal ด้วย `slices.Equal` — เพิ่ม**หรือ**
  ลบสมาชิกกี่ตัวก็ fail ทันที ไม่ต้องเดาว่าเทสต์ table เดิมมีค่าที่ตรงกับ
  ค่าปลอมพอดีหรือเปล่า
- **R3 (แก้แล้ว, promoted จาก reviewer's S1) — แทนที่ caller invariant ด้วย
  type**: `checkkey.go`'s doc comment (29 บรรทัด) เป็น WHY ที่ถูกต้องและ
  reviewer ยืนยันทุกข้อเท็จจริงแล้ว **ไม่ได้ตัดทิ้งเพราะอยากสั้นลง** — ตัดทิ้ง
  เฉพาะประโยคที่ชดเชยด้วยคำอธิบายแทนที่จะบังคับด้วย type: "Callers MUST resolve
  and pass the DB's own canonical question, never raw client input" เพิ่ม type
  ใหม่ `domain.CanonicalQuestion` (`internal/learning/domain/canonicalquestion.go`)
  — `NewCanonicalQuestion(raw string)` trim แล้ว reject ค่าว่าง (ย้าย trim
  logic มาจาก `checkkey.go` เดิม). `NewCheckKey`'s พารามิเตอร์ที่สามเปลี่ยนจาก
  `string` เป็น `CanonicalQuestion` — `Service.RecordAttempt`'s `build`
  closure ไม่มีจุดไหนสร้าง `CanonicalQuestion` จาก `rawQuestion` เองเลย (รับ
  แค่ตัวที่ repository resolve มาให้เป็น parameter) วันนี้ — เป็นหลักการเดียวกับ
  ที่ R1 ใช้กับ `kind` ในรอบนี้เป๊ะ ๆ **แต่ต้องพูดให้แม่นยำ (แก้ตามที่ reviewer
  ชี้): นี่ไม่ใช่ "forge ไม่ได้" ในระดับ compile-time ทั้งหมด** —
  `NewCanonicalQuestion` เป็น exported function จาก `domain`, และ `app` import
  `domain` อยู่แล้ว ดังนั้น `Service.RecordAttempt` เขียนโค้ดหนึ่งบรรทัดสร้าง
  `domain.NewCanonicalQuestion(rawQuestion)` เองได้จริงถ้าใครแก้แบบนั้นในอนาคต
  (ตรงข้ามกับ `kind` ที่หายไปจาก request DTO เลย ซึ่งปิดกั้นได้จริงในระดับ HTTP
  layer). สิ่งที่ปิดกั้นได้จริงระดับ compile-time มีแค่: forge จาก**นอก** package
  `domain` ไม่ได้ (field `value` เป็น unexported, `NewCheckKey` reject zero
  value) สิ่งที่ปิดกั้นการ forge **ภายใน** `app` layer จริง ๆ คือ**เทสต์**
  `TestServiceRecordAttempt_CheckKeyDerivedFromCanonicalQuestion` (reviewer
  ยืนยันแล้วว่า kill โค้ดที่ forge แบบนี้ได้จริง) ไม่ใช่ตัว type เอง — สรุปคือ
  "ไม่มีโค้ดไหนในโปรเจกต์นี้ทำแบบนั้นวันนี้" ไม่ใช่ "ทำแบบนั้นไม่ได้เลย"
- **R4 (แก้แล้ว, promoted จาก reviewer's S2) — stub driver ตาบอดสองจุดที่ R1
  (ทั้งสองรอบ) อยู่พอดี**: `stub_driver_test.go` เดิมไม่มีคอลัมน์ `recall_checks.type`
  เลย และไม่ enforce ENUM ของ `recall_attempts` เลย (MySQL ภายใต้
  `STRICT_TRANS_TABLES` ปฏิเสธค่า out-of-ENUM ด้วย error 1265) — แก้ทั้งสองจุด:
  `seedRecallCheck` รับ `kind` เพิ่ม, `selectRecallCheckExistsSQL`'s stub
  dispatch คืน `rc.type` เป็นคอลัมน์ที่สาม (`recallCheckMatch`/`seededRecallCheck`
  เก็บ kind ไว้คู่กับ question); `insertRecallAttemptSQL`'s exec handler เพิ่ม
  `validateRecallAttemptEnums` เทียบ `type`/`confidence`/`outcome`/`graded_by`
  กับ literal set ที่ mirror `migrations/006_recall.sql`'s ENUM ทุกคอลัมน์ตรง ๆ
  คืน error จำลอง MySQL 1265 ถ้าไม่ match — เพิ่มเทสต์
  `TestStubInsertRecallAttemptRejectsInvalidEnumValue` (4 subtests, หนึ่งต่อ
  คอลัมน์) ยืนยันว่า fidelity fix เองทำงานจริง โดยเรียก `db.ExecContext` ตรง ๆ
  ข้าม domain-layer validation ไปเลย (ทดสอบตัว stub เอง ไม่ใช่ path ของ
  `RecordAttempt`)
- **R5 (แก้แล้ว) — จุดเล็กหลายจุด**:
  - `service.go`'s `RecordAttempt` เดิมส่ง `rawQuestion` **ไม่ trim** เข้า
    query resolve ทั้งที่ `NewCanonicalQuestion`/`NewCheckKey` trim อยู่แล้ว —
    ถ้า client ส่งคำถามที่มี trailing space จะได้ 400 "question not found"
    ทั้งที่จริงเป็นคำถามที่ถูกต้อง แก้ให้ trim ก่อนส่งเข้า `s.repo.RecordAttempt`
    (`trimmedQuestion` แทน `rawQuestion`) ให้ตรงกับสิ่งที่ hash จริง
  - `migrations/006_recall.sql`'s comment เขียนผิดว่า ascii "halves the byte
    width of utf8mb4" — จริง ๆ คือ **quarters** (ascii 1 byte/char เทียบ
    utf8mb4 สูงสุด 4 bytes/char) แก้ข้อความให้ตรง
  - `TestServiceRecordAttempt_EchoesSubmittedValues` เปลี่ยนเป็น table-driven
    ครบทั้งสาม `Confidence` (`guessed`/`unsure`/`confident`) แทนที่จะ assert
    แค่ `"confident"` ตัวเดียว
  - Cleanup แถว synthetic ใน dev DB: ลบ `id=4` (จาก stale-image test รอบก่อน)
    และ `id` 10–13 (เขียนโดย reviewer ระหว่างพิสูจน์ R1 — รวมถึงแถว 13 ที่**คือ
    R1's ตัวบั๊กเอง**, `type='mcq'` ใต้ `check_key` เดียวกับแถว `short_answer`)
    หลังแก้ R1 แล้วแถวพวกนี้เป็น invalid data ที่จะทำให้ Q-2b/Q-2c สับสน — ลบทิ้ง
    ยืนยันด้วย `SELECT` ก่อน/หลัง ด้านล่าง
- **R6 (แก้แล้ว) — เอกสารเกินจริง**: ดูย่อหน้า "สรุปรอบแรก" ในตาราง mutation
  ด้านบนสำหรับรายละเอียดตัวเลขจริง (18/21 killed จาก reviewer's independent set,
  4 survivor ทั้งหมดเป็น R2's class) แทนที่ "11/11 ไม่มี survivor"
- **จุดที่รู้ตัวแล้วไม่แก้ (by design)**: `RecordAttempt` เป็น resolve-then-insert
  ไม่มี transaction (ถูกแล้ว — append-only ไม่มี read-modify-write ให้ race) แต่
  `import-lessons` ทำ `DELETE FROM recall_checks` แล้ว insert ใหม่ — POST attempt
  ที่มาถึงพอดีในหน้าต่างนั้นจะได้ 400 (คำถามหาไม่เจอชั่วคราว) แอปนี้ single-user
  และ import เป็น manual step ไม่ใช่ automated job ที่รันพร้อมกับ traffic จริง
  จึงรับความเสี่ยงนี้ไว้โดยไม่ใส่ lock เพิ่ม — บันทึกไว้เป็น known window เฉย ๆ
- **หนี้ที่รู้ตัวแล้ว ไม่แก้รอบนี้ — `recallAttemptEnumColumns` ซ้ำมือกับ
  migration โดยไม่มีอะไร pin สองจุดนี้เข้าด้วยกัน**: `stub_driver_test.go`'s
  `recallAttemptEnumColumns` (R4) เป็น literal ที่คัดลอกมาจาก
  `migrations/006_recall.sql`'s ENUM ด้วยมือ ไม่มีเทสต์ไหนเทียบสองจุดนี้กันเองว่า
  ตรงกันจริง — ถ้าในอนาคตมีคนขยาย domain array, ขยาย `recallAttemptEnumColumns`,
  และขยาย `AcceptedSetIsExactly`'s literal ทั้งสามจุดพร้อมกัน (ต้องทำสามที่
  ให้ตรงกันเอง) แต่ลืมแก้ migration — ทุก test ในรอบนี้จะยัง**เขียวหมด** แล้วไป
  พังจริงตอน production (500 จาก MySQL ENUM truncation) พอดีเป็น drift ชนิด
  เดียวกับที่ R2 ยกขึ้นมา แค่ขยับไปอีกชั้นหนึ่ง (ระหว่าง stub กับ migration แทนที่
  จะเป็นระหว่าง domain array กับเทสต์) — ความเสี่ยงต่ำ (ต้องพลาดพร้อมกันสามจุด)
  แต่เป็น unpinned seam จุดสุดท้ายของ ticket นี้ที่ยังไม่มีกลไกอัตโนมัติกัน —
  บันทึกไว้เฉย ๆ ยังไม่แก้รอบนี้

### Live verification (docker + curl + MySQL)

Stack: `docker compose up -d --build` (ไม่ลบ volume, ไม่ใช้ `-v` ตอนปิด).
ใช้ lesson `domain-driven-design/ubiquitous-language` (มีทั้ง `short_answer`
และ `mcq` checks จริงจาก DB, ใช้ใน Q-1's evidence เดิมด้วย) และคำถามจริงจาก
`GET /api/v1/lessons/domain-driven-design/ubiquitous-language`.

- **POST สำเร็จ + เห็นแถวจริงใน MySQL**: ยืนยันด้วย `SELECT * FROM
  recall_attempts` หลังยิง POST — เห็นแถวที่มี `check_key`/`type`/`confidence`/
  `outcome`/`selected_option`/`graded_by='self'`/`created_at` ตรงกับที่ submit
- **submit คำถามเดิมซ้ำสองครั้ง → สองแถว**: พิสูจน์ด้วย `COUNT(*)` จริงจาก MySQL
  ก่อน/หลัง ไม่ใช่แค่ assert จาก unit test
- **คำถามของ lesson อื่น → reject 400**: ยิง POST ไปที่
  `.../ubiquitous-language/attempts` ด้วยคำถามที่จริง ๆ เป็นของ concept อื่น
- **confidence ผิด enum → reject 400**
- **`check_key` ตรวจข้ามวิธี**: คำนวณ SHA256 ของ `"topic/concept/question"`
  ด้วย PowerShell (`[System.Security.Cryptography.SHA256]`) แยกจาก Go แล้ว
  เทียบ hex ตรงกับค่าที่ API เก็บจริง — พิสูจน์ Go ไม่ได้ "เห็นด้วยกับตัวเอง" ฝ่ายเดียว
- **R1's fix พิสูจน์สดกับ MySQL จริง (case-variant → canonical check_key
  เดียวกัน)**: submit คำถามเดิมแบบตัวพิมพ์เล็กทั้งหมด (`ubiquitous language...`
  แทน `Ubiquitous Language...`) ไปที่ endpoint เดียวกัน — response's `check_key`
  ออกมา **ตรงกับ** check_key ของการ submit ด้วย casing ที่ถูกต้องเป๊ะทุกตัวอักษร
  (`6ce0ab5154b2e30501ddaf0d58934dd0b4579d06adc36b88d5fc2a2c47bba108`, ยืนยัน
  ด้วย PowerShell SHA256 อีกรอบจากข้อความ canonical โดยตรง) และ response's
  `question` field คืนข้อความ canonical จริง (`Ubiquitous Language...` ตัวใหญ่
  ตรงกับที่ query `recall_checks` ตรง ๆ ยืนยันแล้ว) ไม่ใช่ casing ที่ submit เข้ามา
  — **หมายเหตุตรง ๆ**: รอบแรกที่ทดสอบ (ก่อน rebuild image ครั้งสุดท้าย) ได้
  check_key ผิด (`523fdeca9c...`) เพราะ `docker compose up -d --build` ครั้งนั้น
  บังเอิญ build image จาก source ระหว่างที่กำลังทำ mutation testing อยู่พอดี
  (ไม่ใช่ source ที่ commit) — แก้โดย `docker compose build --no-cache api` แล้ว
  `docker compose up -d api` ใหม่ ยืนยัน source บน disk ถูกต้องด้วย `go build
  ./...`/`go vet ./...`/`go test ./... -count=1` ก่อน rebuild แล้วจึง retest ได้
  check_key ที่ถูกต้อง — แถวที่ผิด (id=4 ใน `recall_attempts`) ยังอยู่ใน DB จริง
  (local dev data ไม่ลบทิ้ง) เป็นหลักฐานของทั้งบั๊กที่เคยเกิดและการแก้ที่ยืนยันแล้ว

(รายละเอียด output/curl commands ทั้งหมดอยู่ใน PR description ของ
`ticket/q-2a-recall-attempts`)

### Live verification รอบ 2 (หลังแก้ R1's kind fix + CanonicalQuestion refactor)

Rebuild image จริง (`docker compose build api` แล้ว `docker compose up -d api`,
ไม่ใช้ `-v`) หลัง `go vet ./...`/`gofmt -l .`/`go test -count=1 ./...` เขียวหมด
บนโค้ดที่จะ commit เพื่อไม่ให้เกิดปัญหาเดิม (รอบก่อนหน้าเจอ image ที่ build จาก
source กลางอากาศระหว่างมูเทตพอดี — ดูหมายเหตุด้านบน):

- **ส่ง `kind` ใน request body → 400 ทันที**: `{"question":"...","kind":"mcq",...}`
  ไปที่ endpoint เดิม → `{"error":"invalid request body"}`, HTTP 400 —
  `recordAttemptRequest` ไม่มี field `Kind` อีกแล้ว `DisallowUnknownFields()`
  ปฏิเสธเหมือน field แปลกอื่น ๆ — พิสูจน์ว่า "kind ไม่ตรงกับ recall_checks.type"
  เป็นไปไม่ได้จริงในทางปฏิบัติ ไม่ใช่แค่ validate แล้วปฏิเสธ
- **case-variant → canonical check_key เดิม ยังคงอยู่หลัง refactor**: submit
  คำถาม mcq จริง (`ทำไมคำว่า Account...`) ด้วย casing ปกติ ได้
  `check_key=321730c1eb55d0131a5fe466611d192a4691eb0a908d008e9771d9da4994883c`,
  `kind:"mcq"` (resolve จาก DB เอง ไม่ใช่จาก request). submit คำถามเดิมแบบ
  `ACCOUNT` ตัวใหญ่ทั้งหมด → **`check_key` เดิมทุกตัวอักษร** และ response's
  `question` คืนข้อความ canonical (`Account` ตัวพิมพ์ปกติ) — ยืนยันว่า
  `CanonicalQuestion` refactor ไม่ทำให้พฤติกรรม R1's fix เปลี่ยนเลย
- **`check_key` ตรวจข้ามสามวิธี**: Go (จากค่าที่ API ตอบ), `sha256sum` ของ
  `printf 'domain-driven-design/ubiquitous-language/ทำไมคำว่า Account...'`,
  และ MySQL's **`SHA2(CONCAT(topic,'/',concept,'/',question), 256)`** เอง —
  ทั้งสามวิธีให้ hex เดียวกันเป๊ะ (`321730c1...`) พิสูจน์ว่าไม่ใช่แค่ Go เห็นด้วย
  กับตัวเอง
- **append-only ยังคงอยู่**: `SELECT id, type, confidence, outcome, created_at
  FROM recall_attempts WHERE check_key='321730c1...'` เห็น 6 แถวสำหรับคำถาม
  เดียวกัน (รวมสองแถวใหม่จากรอบนี้) ทุกแถว `type='mcq'` ตรงกัน (พิสูจน์ R1's kind
  fix แบบ positive: ไม่มีแถวไหนเป็น `type` อื่นใต้ `check_key` เดียวกันอีกแล้ว)
- **Cleanup**: ลบแถว synthetic `id=4` (stale-image bug รอบก่อน) และ `id`
  10–13 (รวมแถว 13 ที่เป็น R1's ตัวบั๊กเอง — `type='mcq'` ใต้ `check_key`
  เดียวกับแถว `short_answer`) ยืนยันด้วย `SELECT` ก่อนลบ (เห็นเนื้อหาตรงกับที่
  reviewer อธิบายไว้เป๊ะ) และหลังลบ (เหลือแถวที่ถูกต้องทั้งหมด, ไม่มี `type`
  ขัดแย้งกันใต้ `check_key` เดียวกันอีกเลย)

(รายละเอียด curl/SQL command ทั้งหมดของรอบ 2 อยู่ใน PR description)

### จงใจไม่ทำในรอบนี้

- **ป้องกัน concept slug ซ้ำข้ามสองบทของ topic เดียวกันแล้ว (แก้เพิ่มจากรอบ
  code-reviewer แรก)**: ตอนร่างเอกสารรอบแรกเข้าใจผิดว่าเรื่องนี้ยัง "จงใจไม่ทำ" —
  จริง ๆ แล้ว resolve step (`resolveRecallCheck` ใน `repository.go`) ใช้
  `QueryContext` + loop แบบเดียวกับ `lockLesson` ใน `Transition` เป๊ะ ๆ: ถ้า
  `selectRecallCheckExistsSQL` match มากกว่าหนึ่งแถว (สอง `recall_checks` ใน
  lesson เดียวกันที่ข้อความ fold เป็นคำเดียวกันภายใต้ case/accent-insensitive
  collation) จะ error ทันที ไม่ใช่เลือกแถวแรกเงียบ ๆ — pin ด้วย
  `TestRepositoryRecordAttempt_AmbiguousCaseVariantRejected` (mutation #10 ด้านบน)
- **ไม่มี `review_cards`/`review_logs`/SM-2** — Q-2c
- **ไม่แตะ frontend เลย** — `RecallCheckCard`'s confidence/selected
  option/outcome ยังอยู่ใน memory เหมือนเดิม จะหายตอน refresh เหมือนที่ Q-1 บันทึก
  ไว้ — Q-2b คือ ticket ที่เรียก endpoint นี้จริงแล้ว persist เป็นครั้งแรก
- **`kind` cross-check กับ `recall_checks.type` — ตัดสินใจแล้วในรอบ 2, ไม่ใช่
  "จงใจไม่ทำ" อีกต่อไป**: รอบแรกเขียนว่า "client บอก kind เอง, ไม่ cross-check"
  เป็นการตัดสินใจที่ตั้งใจ — code review รอบ 2 พิสูจน์ว่านั่นเป็นช่องโหว่จริง
  (ดู R1 ในหัวข้อ "รอบ code-reviewer 2" ด้านบน) แก้โดย**ตัด `kind` ออกจาก client
  ทั้งหมด** (แรงกว่าแค่ cross-check แล้วปฏิเสธ) — บรรทัดนี้เก็บไว้เป็นประวัติว่า
  การตัดสินใจรอบแรกเปลี่ยนไปยังไงและทำไม ไม่ใช่ debt ที่ยังค้างอยู่

### Review focus

- ทำไมการ hash คำถามที่ client ส่งมาโดยตรง (แทนคำถาม canonical ที่ resolve จาก
  DB) ถึงเป็นบั๊กจริง ทั้งที่เช็ค "คำถามเป็นของ lesson นี้" ผ่านแล้วก็ตาม? —
  แล้ว `kind` ในรอบ 2 เป็นบั๊กชนิดเดียวกันตรงไหน ต่างกันตรงไหน?
- ทำไมมูเทชัน hardcode `graded_by` เป็น `"self"` ถึง kill ที่ repository-layer
  test แต่ **ไม่** kill ที่ app-layer test ทั้งที่แก้โค้ด production จุดเดียวกัน —
  เป็น coverage gap จริงหรือเป็น equivalent mutant?
- ทำไม `check_key` ต้อง hash `topic + "/" + concept + "/" + canonical-question`
  แทนที่จะ FK ไปตรง ๆ ที่ `recall_checks.id` หรือใช้ `(lesson_id, position)`?
- ทำไมเทสต์แบบ "ลองค่าที่ถูก 3 ค่า + ค่าผิด 1 ค่า" ถึงไม่พอสำหรับปิด enum ทั้งชุด
  ทั้งที่ค่าที่ถูกทั้งหมดถูกทดสอบแล้วจริง ๆ?

Status: implemented, round 2 fixes applied post code-review, PR pending

## Q-2b — wire the quiz to the attempts API

- **Scope**: **Frontend ล้วน** — `git diff --name-only develop... -- '*.go' 'migrations/*'`
  ว่างเปล่าจริง. `web/lib/api.ts` เพิ่ม `postAttempt` + types (`AttemptConfidence`,
  `AttemptOutcome`, `AttemptInput`, `AttemptRecord`, `PostAttemptResult`).
  `web/components/RecallCheckCard.tsx` เพิ่ม prop `onAttemptReady` (ยิงเมื่อ attempt
  ครบจริง — export type `CompletedAttempt` ด้วย, รอบ 2), `saveStatus`/`onRetrySave`
  (แสดง "Saving…"/"Not saved" + Retry ต่อการ์ด), export `AttemptSaveStatus` ใหม่.
  `web/app/(app)/lesson/page.tsx` เพิ่ม `attemptStatus`/`pendingAttemptsRef`/
  `debounceTimersRef`/`loadGenerationRef` state, `submitAttempt`/`flushPendingAttempt`/
  `flushAllPendingAttempts`/`handleAttemptReady`/`handleRetrySave`, banner "Saving N
  recall attempt(s)…" + "Some recall attempts didn't save" เหนือปุ่ม Finish (รอบ 2
  แยกสองบรรทัด). **รอบ 2 (หลัง code review)**: `loadGenerationRef` แทน slug-based
  identity guard เดิม, short_answer debounce 1s ก่อน submit จริง (mcq ยังทันที),
  flush 3 เส้นทาง (quiet period/Finish/navigate away), แก้ dead branch + silent-drop
  path ใน `submitAttempt`/`handleAttemptReady` — รายละเอียดเต็มอยู่ที่หัวข้อ "รอบ
  code-reviewer" ด้านล่าง. **รอบ 3 (หลัง code review 2)**: `handleRetrySave`
  เปลี่ยนไปเรียก `flushPendingAttempt` แทน `submitAttempt` ตรง ๆ (กัน race กับ
  debounce timer ที่ยังไม่หมดเวลา), `postAttempt`/`submitAttempt` เพิ่ม
  `options?: {keepalive: boolean}`, effect ใหม่ฟัง `pagehide`/
  `visibilitychange` flush ด้วย `keepalive:true` — รายละเอียดเต็มอยู่ที่หัวข้อ
  "รอบ code-reviewer 3" ด้านล่าง. **รอบ 4 (หลัง code review 3, SHIP พร้อม
  follow-up)**: ไม่มีการแก้ production code เลย — เพิ่มเทสต์ 2 ตัวปักหมุด
  `delete debounceTimersRef.current[position]` และ `removeEventListener`
  cleanup ที่ไม่มีเทสต์คุ้มครองมาก่อน (ทั้งสองบรรทัดถูกอยู่แล้ว) — ดูหัวข้อ
  "รอบ code-reviewer 4" ด้านล่าง. เทสต์สุดท้าย: `api.test.ts` **6 → 12**
  (+6), `RecallCheckCard.test.tsx` **30 → 44** (+14), `page.test.tsx`
  **2 → 27** (+25, ไฟล์เดียวกับที่ UX-5 ปักไว้ ไม่ใช่ไฟล์แยก) — รวม
  `npm test` **124 → 169**

### การตัดสินใจหลัก

- **`postAttempt` คืน `{kind:"ok"|"error"}` ไม่ throw สำหรับ non-auth failure — ตาม
  precedent ของ `getProgress` ไม่ใช่ `setProgress`/`getLesson`**: caller (`lesson/
  page.tsx`) ต้องแยก "saved" กับ "not saved" ต่อ check เพื่อโชว์ indicator/retry ต่อ
  การ์ด — throw exception จะบังคับให้ต้อง try/catch กระจายอยู่ที่ทุกจุดเรียก และ
  ยากที่จะเก็บ "ยังไม่ได้ save" เป็น first-class state แบบที่ discriminated result
  ให้ได้ฟรี. `UnauthorizedError` ยัง throw เหมือนเดิม (ไม่ยุบรวมเข้า `{kind:
  "error"}`) เพราะเป็นเงื่อนไขระดับ session (bearer token ทั้งก้อนใช้ไม่ได้) ไม่ใช่
  เรื่องต่อ attempt — caller ทุกจุดใน api.ts จัดการ 401 แบบเดียวกันหมด (throw แล้ว
  redirect)
- **Submit ครั้งเดียวต่อ attempt ที่สมบูรณ์ — dedup ด้วย ref เก็บ "signature" ของ
  attempt ไม่ใช่ boolean "เคย submit หรือยัง"**: `lastReportedAttemptRef` เก็บ
  `` `${confidence}|${outcome}|${selectedOption}` `` — mcq freeze ค่าทั้งหมดตอน
  reveal จึง signature เดียวตลอด (submit ครั้งเดียวโดยธรรมชาติ), short_answer
  เปลี่ยน rating ได้หลัง reveal (Pass/Not yet ไม่ได้ล็อก) ทำให้ signature เปลี่ยน
  จริง — ถ้าใช้ boolean เดิมจะบล็อกการ resubmit ที่ตั้งใจให้เกิดขึ้น (ดูข้อถัดไป).
  Guard นี้กัน 2 เหตุการณ์: re-render ที่ parent ส่ง `onAttemptReady` reference
  ใหม่มา (dependency ของ effect เปลี่ยนแต่ค่า attempt ไม่เปลี่ยน) และคลิก
  Pass/Not yet ซ้ำค่าเดิม (React bail out การ set state ค่า primitive เดิม เอง
  — effect ไม่รันซ้ำด้วยซ้ำ). **แก้ไขจากที่เอกสารรอบแรกเคยเขียนผิด**: guard นี้
  **ไม่ได้**กัน React StrictMode's double-invoked effects — StrictMode
  double-invoke เฉพาะตอน mount เท่านั้น ส่วน effect นี้ early-return ตอน mount
  เสมอ (ยังไม่ครบเงื่อนไข finished ตอนนั้น) แล้วมารายงานจริงตอน state เปลี่ยน
  ทีหลังซึ่งเป็นคนละรอบ StrictMode ไม่ double-invoke ซ้ำให้ — พิสูจน์ด้วย
  mutation จริง: ปิด dedup check แล้ว StrictMode test **ยังผ่าน**, มีแค่เทสต์
  "re-render เปลี่ยน callback identity" เท่านั้นที่จับ mutation นี้ได้จริง
  (ดู M1 ในตารางมูเทชันด้านล่าง)
- **short_answer rating เปลี่ยน (Pass → Not yet) = resubmit ไม่ใช่ ignore หรือ
  "ค่าแรกชนะ", แต่ submit จริงหลัง debounce settle เท่านั้น (แก้ในรอบ 3 — ดู
  R5 ของ "รอบ code-reviewer" ด้านล่าง)**: ตัดสินใจแล้วว่าการแก้ไข self-rating
  คือข้อมูลใหม่ที่ SRS (Q-2c) ต้องการจริง ไม่ใช่ noise ที่ควรกรองทิ้ง — endpoint
  เป็น append-only อยู่แล้วตาม design ของ Q-2a (ไม่มี idempotency key) จึงรองรับ
  pattern นี้ได้โดยไม่ต้องแก้ backend เลย แต่ละแถวคือสิ่งที่ user เชื่อ ณ ขณะนั้น
  จริง ๆ (Pass ตอนแรกอาจเป็นการประเมินที่มั่นใจเกินจริง, แก้เป็น Not yet ทีหลัง
  คือสัญญาณที่มีค่ากว่าการทิ้งไป) — **ข้อสำคัญ: "resubmit" ในที่นี้หมายถึง
  ค่าที่ settle แล้วเท่านั้น** ไม่ใช่ทุกครั้งที่กด Pass/Not yet ทันที (นั่นคือ
  design รอบ 1/2 ที่พิสูจน์แล้วว่าเขียนแถวซ้ำไบต์ต่อไบต์ได้จาก flip-flop รัว ๆ
  — แก้ด้วย debounce ใน R5)
- **Retry ไม่ dedupe กับ request ที่ "จริง ๆ สำเร็จแต่ client เห็นเป็น fail" (เช่น
  timeout) — ยอมรับแถวซ้ำที่อาจเกิดขึ้น ไม่ปิดกั้น**: endpoint ไม่มี idempotency
  key (ตัดสินใจแล้วใน Q-2a) การจะกัน duplicate ที่มาจาก retry-after-phantom-success
  ต้องแก้ backend (ออก idempotency key ต่อ attempt) ซึ่งอยู่นอก scope ของ ticket นี้
  — ทางเลือกที่เหลือคือ "ไม่เสนอ retry เลยถ้าไม่ชัวร์ว่า fail จริง" ซึ่งขัดกับกฎข้อ 4
  ของ ticket นี้ตรง ๆ (ห้าม suppress ตัวบอกว่า fail) และจะทำให้ fail จริงกู้คืนไม่ได้
  เลือกยอมรับความเสี่ยง duplicate ไว้ (เหมือนกับ double-submit โดยตั้งใจ) —
  **แก้คำอธิบายในรอบ code review**: เอกสารรอบแรกเคยเขียนต่อท้ายว่า "เป็น data
  problem ที่ Q-2c อ่าน `created_at DESC` ล่าสุดอยู่แล้วไม่กระทบ" ซึ่งเป็นเท็จ
  สองชั้น — (1) อ้าง mitigation ของ ticket ที่ยังไม่เริ่ม (Q-2c ยังไม่ได้ทำ)
  (2) ต่อให้ Q-2c อ่าน `created_at DESC` จริง ก็ยังผิดอยู่ดีเพราะ `created_at`
  เป็น `TIMESTAMP` second-precision — สองแถวที่ submit ห่างกันจริงในเวลาปกติ
  (ไม่ใช่ race) ตกวินาทีเดียวกันได้จริง ทำให้เรียง "ล่าสุด" ผิดได้ (รายละเอียด
  เต็มอยู่ที่ R4/R6 ของ "รอบ code-reviewer" ด้านล่าง, ข้อกำหนดสำหรับ Q-2c
  บันทึกไว้ใน `docs/roadmap.md` แล้ว: `ORDER BY created_at DESC, id DESC`)
  — บันทึกไว้ตรง ๆ ว่านี่คือ **data-ordering risk ที่ยอมรับ** ไม่ใช่ risk ที่
  ไม่กระทบอะไรเลย
- **Identity guard: attempt state อยู่ที่ `LessonView` (parent) ไม่ใช่
  `RecallCheckCard`, และ `stillCurrent()` เทียบ "visit" ไม่ใช่ slug (แก้ในรอบ
  2 — ดู R1/R2 ของ "รอบ code-reviewer" ด้านล่างสำหรับหลักฐานและรายละเอียดเต็ม)**:
  attempt state (`attemptStatus`, `pendingAttemptsRef`) เป็นของ parent โดยตั้งใจ
  — ถ้าให้ `RecallCheckCard` เก็บ save-status เองในเครื่อง (local state) ความ
  ปลอดภัยจาก lesson A รั่วเข้า lesson B จะได้มาฟรีจากการที่การ์ดทั้งก้อน unmount
  ตอนเปลี่ยน lesson (ข้อสรุปเดิมจาก Q-1's M11) แต่ banner รวม "some attempts
  didn't save"/"saving" เหนือปุ่ม Finish ต้องอ่าน state ข้ามทุกการ์ดพร้อมกัน —
  บังคับให้ state ต้องอยู่ที่ parent ซึ่ง **ไม่** unmount ข้าม query-param
  navigation (ต่างจาก child cards) จึงต้องมี guard จริง ไม่ใช่ได้มาฟรีจาก
  unmount. **รอบแรกใช้ slug equality (`currentIdentityRef`) เหมือน `finish()`
  เดิม แล้วพบว่าไม่พอจริง**: slug เทียบได้แค่ "lesson นี้ตรงกับที่จอแสดงอยู่
  หรือเปล่า" ไม่ใช่ "response นี้เป็นของ visit ที่เริ่มมันขึ้นมาจริงหรือเปล่า" —
  กลับมาที่ lesson เดิมซ้ำ (A → B → กลับมา A) มี slug เดิมทุกตัวอักษรกับ visit
  ก่อนหน้า ทำให้ response ค้างของ visit แรกถูกเข้าใจผิดว่าเป็นของ visit ปัจจุบัน
  ได้ — แก้ด้วย `loadGenerationRef` (นับ visit จริง ไม่ใช่ slug). **หนี้ที่รู้ตัว
  แล้ว**: `lesson/page.tsx` ตอนนี้มี identity mechanism **สองแบบข้าง ๆ กัน** —
  `submitAttempt` ใช้ `loadGenerationRef` (นับ visit) แต่ `finish()` ที่อยู่
  ไม่กี่บรรทัดถัดไปในไฟล์เดียวกันยังใช้ `currentIdentityRef` (slug equality) —
  เป็นกับดักจริงสำหรับคนอ่านโค้ดครั้งถัดไปที่อาจก็อป pattern ผิดตัวไปใช้ที่อื่น
  โดยไม่รู้ว่ามีสองมาตรฐานซ้อนกันอยู่ (เหตุผลที่ยังไม่รวมเป็นอันเดียว ดูหัวข้อ
  "finish()" ท้าย ticket นี้)
- **Retry คือ flush ครั้งเดียว ไม่ใช่ submit คู่ขนานกับ timer ที่ยังไม่หมดเวลา
  (แก้ในรอบ 3 — ดู R... "Retry racing a live debounce timer" ด้านล่าง)**:
  `handleRetrySave` เรียก `flushPendingAttempt(position)` (เคลียร์ debounce
  timer ที่ค้างอยู่ก่อนเสมอ แล้วค่อย submit ค่าล่าสุดที่มี) แทนที่จะเรียก
  `submitAttempt` ตรง ๆ ด้วย payload เก่า — เพราะ Retry render ได้ตอน
  `saveStatus==="error"` ซึ่ง**ไม่เปลี่ยน**เวลาผู้ใช้แก้ไข rating ใหม่ (แก้ไข
  แค่ schedule timer ใหม่ ไม่แตะ `attemptStatus`) ทำให้ปุ่ม Retry ค้างอยู่บนจอ
  พร้อมกับ timer ที่กำลังนับถอยหลังอยู่พร้อมกันได้จริง — กด Retry ตอนนั้นโดยไม่
  เคลียร์ timer ก่อนจะได้ POST สองครั้งคนละค่ากันใต้ `check_key` เดียวกัน
- **Error handling ไม่แยกตาม HTTP status (400/404/413/500 ทั้งหมด =
  `{kind:"error"}` เดียวกัน)**: ตั้งใจไม่ทำ granular error taxonomy ต่อ ticket
  scope ("Not saved" + Retry พอสำหรับ v1) — ต่างจาก `finishState.error` เดิมที่โชว์
  `err.message` เพราะที่นั่น error มาจาก `ApiError`/network โดยตรงไม่ผ่าน
  discriminated result ที่ตั้งใจซ่อนรายละเอียดไว้แล้ว
- **Aggregate banner (Finish section) มีสองแบบคนละ role**: banner ความล้มเหลว
  ("Some recall attempts didn't save") ใช้ `role="alert"` เหมือน
  `finishState.error` เดิม — สอดคล้องกับ error banner ที่มีอยู่แล้วในไฟล์เดียวกัน
  ต่างจากที่ Q-1's REQ-3 เลี่ยง `role="status"` ห่อปุ่ม Pass/Not yet (ปัญหานั้นคือ
  live region ห่อ control ที่ toggle `aria-pressed` ถี่ ๆ) — banner นี้ไม่มี
  control ข้างในเลย เป็น one-shot mount เหมือน `finishState.error` เป๊ะ. banner
  ที่สอง ("Saving N recall attempt(s)…", เพิ่มในรอบ 2) ใช้ `role="status"` แทน
  เพราะเป็น**ข้อมูล**ไม่ใช่**ข้อผิดพลาด** (สีกลาง ไม่ใช่สี danger) และจำนวน
  N เปลี่ยนได้หลายครั้งระหว่างที่ยังมี attempt ค้างอยู่ (ไม่ใช่ one-shot เหมือน
  banner แรก) แต่ไม่มี control โต้ตอบข้างในเช่นกัน จึงไม่เข้าเงื่อนไขที่ REQ-3
  กังวลไว้

### Mutation table รอบ 1 (11/11 mutations ที่ทดลองในรอบนี้ตายหมด — ไม่ใช่ข้อสรุปว่าไม่มี gap เหลือ)

**หัวข้อนี้จำกัดเฉพาะ 11 มูเทชันที่ทดลองในรอบแรกเท่านั้น** — รอบ code-reviewer
2/3/4 ด้านล่างเจอ survivor จริงเพิ่มอีกหลายจุดที่ชุดนี้ไม่ครอบคลุม (identity
guard, retry-body, "saving" state, multi-check flush, unload listener
cleanup) ดูตารางของแต่ละรอบสำหรับรายละเอียด

รัน `npx vitest run components/RecallCheckCard.test.tsx "app/(app)/lesson/page.test.tsx"`
หลังแก้แต่ละจุด แล้ว revert ทุกครั้ง:

| # | Mutation | ผลลัพธ์ |
|---|---|---|
| M1 | ลบ `lastReportedAttemptRef` dedup check ออกจาก effect ใน `RecallCheckCard` | killed — **StrictMode test ไม่จับ** (effect ที่ report จริงไม่ได้รันตอน mount, รันตอน state เปลี่ยนทีหลัง ซึ่ง StrictMode double-invoke เฉพาะตอน mount) แต่ "does not re-fire on a re-render that leaves the completed attempt unchanged, even if onAttemptReady's own identity changes" (บังคับ re-render ด้วย `onAttemptReady` inline function ใหม่ทุกครั้ง) จับได้จริง (1 failed/43) — บันทึกไว้ตรง ๆ ว่า StrictMode test คนละแบบให้ coverage คนละมุม ไม่ใช่ redundant |
| M2 | สลับ `isMcqCorrect ? "correct" : "incorrect"` กลับด้าน | killed — 4 เทสต์ (2 ใน RecallCheckCard, 2 ใน page.test.tsx ที่เช็ค outcome ตรง ๆ) |
| M3 | hardcode `outcome` เป็น `"correct"` เสมอ | killed — 4 เทสต์ (mcq-incorrect ทั้งสองระดับ + short_answer resubmit ทั้งสองระดับ เพราะ hardcode ทำให้ signature ไม่เปลี่ยนตอน Pass→Not yet ด้วย เจอ dedup ผิดที่พ่วงมา) |
| M4 | hardcode `confidence` เป็น `"guessed"` ใน `onAttemptReady` call | killed — 8 เทสต์ (correct/incorrect x2, confidence table x2 x2 ระดับ, resubmit test) — จุดที่ Q-1's review เจอ bug class เดียวกันตรง ๆ |
| M5 | `selectedOption` ใช้ `shortAnswerRating` แทน `null` สำหรับ short_answer | killed — 3 เทสต์ |
| M6 | `handleAttemptReady` ส่ง `check.expected_answer` แทน `check.question` | killed — 2 เทสต์ (page-level เท่านั้น เพราะ bug อยู่ที่ page.tsx ไม่ใช่ RecallCheckCard) |
| M7 | ปิด error indicator (`{false && saveStatus === "error" && (...)}`) | killed — 2 เทสต์ |
| M8 | โชว์ retry affordance ตอน `saveStatus === "saved"` ด้วย | killed — 1 เทสต์ที่ RecallCheckCard level; **เทสต์ page-level เดิมของฉันเองมี race condition ที่ปิดบัง mutation นี้ได้ (ดูหัวข้อถัดไป) แก้แล้วก่อนสรุปตาราง** |
| M9 | ตัด `router.replace("/token")` ออกจาก catch ของ `submitAttempt` | killed — 1 เทสต์ |
| M10 | ลบ `stillCurrent()` guard ทั้งสองจุดใน `submitAttempt` | killed — **เทสต์แรกที่เขียนไว้ไม่จับ (false pass) แก้แล้วก่อนสรุปตาราง (ดูหัวข้อถัดไป)** |
| M11 | เปลี่ยนเงื่อนไข effect จาก `isRecallCheckFinished(...)` เป็น `stage === "recall"` (ยอมให้ยิงตั้งแต่ stage "commit") | killed — 4 เทสต์ |

**สรุป: 11/11 mutation ที่ ticket บังคับตายหมดจริง ไม่มี survivor** — แต่ระหว่างทำ
เจอว่าเทสต์ของตัวเอง 2 ตัว (สำหรับ M8, M10) เขียนพลาดจนปล่อยให้ mutation รอดในรอบ
แรก รายละเอียดสองจุดนี้อยู่หัวข้อถัดไป (ตรงตามที่ CLAUDE.md เรียกร้อง — "ต้องพัง
invariant ของจริงแล้วบอกให้ได้ว่า test ตัวไหนจับ" รวมถึงตอนที่เทสต์ตัวเองพลาดด้วย)

### สองจุดที่เทสต์ของตัวเองปล่อยให้มูเทชันรอดในรอบแรก (พบระหว่าง mutation testing เอง)

- **M8 (retry โชว์ตอน saved) — page-level test เดิมมี race**: เทสต์แรกเขียนแบบ
  `postAttempt.mockResolvedValueOnce(...)` (resolve ทันที) แล้ว
  `await waitFor(() => expect(screen.queryByText("Not saved")).toBeNull())` —
  ปัญหาคือระหว่าง retry, `saveStatus` เปลี่ยนผ่าน `"saving"` ก่อน (ซึ่งซ่อน "Not
  saved" อยู่แล้วโดยไม่เกี่ยวกับ mutation) ทำให้ `waitFor` เจอสภาวะที่ assertion
  ผ่าน**ชั่วคราว**ระหว่าง "saving" แล้วหยุดรอทันที ก่อนที่ resolve จริงเป็น "saved"
  (ซึ่งด้วย mutation จะทำให้ "Not saved" โผล่กลับมา) จะเกิดขึ้นด้วยซ้ำ — คลาสเดียว
  กับ false-green ที่ project นี้เจอมาก่อน (`docs/roadmap.md`'s "Review discipline"
  memory). แก้ด้วย deferred promise ควบคุมเองแทน `mockResolvedValueOnce` แล้ว
  `await act(async () => { resolveRetry(...); await retryPromise; })` ก่อน assert
  แบบ synchronous (ไม่ผ่าน `waitFor`) — ยืนยันว่า assert เกิด**หลัง**settle จริง ไม่ใช่
  จังหวะไหนก็ได้ที่บังเอิญผ่าน
- **M10 (ลบ `stillCurrent()` guard) — เทสต์แรก assert ผิดตัวบ่งชี้**: เทสต์แรก
  assert `screen.queryByText("Not saved")` (per-check indicator) เป็น null แต่
  per-check indicator อยู่ใน `stage === "reveal"` block เท่านั้น — lesson B's
  check ที่ยังไม่ถูกแตะยังอยู่ stage "recall" เสมอ ทำให้ text นี้เป็น null **ไม่ว่า
  guard จะทำงานหรือไม่** (structurally ไม่มีทางโผล่). ตัวบ่งชี้ที่ leak จริงคือ
  **aggregate banner** ("Some recall attempts didn't save") ที่ผูกกับ
  `attemptStatus` ระดับ `LessonView` ตรง ๆ ไม่ผ่าน stage ของ child เลย — เจอจาก
  `screen.debug()` ตอน mutation ยังไม่ถูก revert (เห็น banner โผล่จริงในเทสต์ที่
  "ผ่าน") แก้โดยเพิ่ม assertion บน aggregate banner text เป็นตัวหลัก

### Playwright + SQL evidence (docker compose up -d --build, ไม่ใช้ -v)

Stack: `docker compose up -d --build` แล้ว `docker compose down` เปล่า ๆ ท้ายสุด
(ไม่มี `-v`, ไม่มีการลบ volume). ใช้ lesson `domain-driven-design/
ubiquitous-language` เดิม (มีทั้ง mcq และ short_answer, ใช้ใน Q-1/Q-2a's evidence
ด้วย) กับ headless Chromium ผ่าน Playwright:

- **หนึ่ง mcq + หนึ่ง short_answer → สองแถวจริงใน MySQL, `check_key` ตรงกับที่คำนวณ
  อิสระด้วย Node's `crypto` (ไม่ใช่แค่ Go เห็นด้วยกับตัวเอง)**:
  ```
  id  type          confidence  outcome  selected_option                              check_key
  22  mcq           confident   correct  เพราะแต่ละความหมายอยู่คนละ bounded context   321730c1eb55d0131a5fe466611d192a4691eb0a908d008e9771d9da4994883c
  23  short_answer  guessed     correct  NULL                                          6ce0ab5154b2e30501ddaf0d58934dd0b4579d06adc36b88d5fc2a2c47bba108
  ```
  ทั้งสอง `check_key` ตรงกับ `SHA256("domain-driven-design/ubiquitous-language/" +
  question)` ที่คำนวณแยกด้วย Node เป๊ะทุกตัวอักษร — ยืนยันว่า `question` ที่ส่งไป
  ตรงกับคำถามจริงของ check นั้น (ไม่งั้น server จะตอบ 400 "question not found" ไป
  แล้ว ไม่มีแถวให้เห็นด้วยซ้ำ)
- **Reload หน้าแล้วตอบ mcq เดิมซ้ำ (option เดิม, confidence เปลี่ยนเป็น Unsure) →
  สองแถวจริงใต้ `check_key` เดียวกัน — นี่คือพฤติกรรมที่ตั้งใจ (append-only, ไม่มี
  idempotency key ตาม Q-2a)**: `id=22 (confident/correct)` และ `id=24
  (unsure/correct)` ทั้งคู่ `check_key=321730c1...` เหมือนกัน
- **Mocked 500 (Playwright route interception) → "Not saved" + Retry โผล่จริง,
  ไม่มี false "saved" state**: `role="alert"` บนการ์ดมีข้อความ "Not saved\nRetry"
  พอดี, `pageerror` count = **0** ตลอด (ไม่มี uncaught exception), `console`
  error ที่เห็นมีแค่ 1 บรรทัดคือ browser's network-level log ของ mocked 500 เอง
  ("Failed to load resource...") ไม่ใช่ error จาก JS runtime — ไม่มี error storm
- **Contrast วัดจริงจาก browser (`getComputedStyle` + แปลง Tailwind v4's
  `oklab()` compound-opacity background กลับเป็น sRGB ด้วยเมทริกซ์ CSS Color 4
  ก่อนคำนวณ, cross-check กับค่าที่ Q-1 เคยวัดสีชุดเดียวกันไว้)**: banner "Not
  saved" ใช้สีชุดเดียวกับ Finish's error banner เดิม (`bg-danger/10 text-
  danger-strong border-danger/30`) — fg `rgb(190,18,60)` บน bg
  `rgb(252,231,232)` = **5.30:1** (ใกล้เคียงค่า Q-1 เคยวัดไว้ 5.31:1 สำหรับสีชุด
  เดียวกัน ต่างกันแค่ rounding — ยืนยันว่าการแปลง oklab ในสคริปต์นี้ถูกต้อง) — ผ่าน
  AA (≥4.5:1)
- **401 (route interception) → redirect ไป `/token` จริง**: `page.url()` หลังกด
  Reveal ตรงกับ `http://localhost:3000/token` เป๊ะ
- **375px**: `document.documentElement.scrollWidth === clientWidth === 375` จริง
  ไม่มี horizontal scroll
- **Cleanup**: ลบเฉพาะแถวที่ script นี้สร้างเอง (`id > 18`, ทั้งหมด 6 แถวจาก 2 รอบ
  รัน — รอบแรก id 19-21 ก่อน crash ตอนแก้ contrast parser, รอบสุดท้าย id 22-24)
  ยืนยันด้วย `SELECT` ก่อน (`COUNT(*)=19, MAX(id)=24`) และหลังลบ (`COUNT(*)=13,
  MAX(id)=18` — กลับสู่สภาพเดิมของ Q-2a's dev data เป๊ะ ไม่แตะแถว 1-18 ที่มีอยู่
  ก่อนเลย)
- **หมายเหตุ**: ระหว่าง phase mocked-500 เจอว่า `lesson_progress` ของ lesson นี้
  ถูก mark `passed` ไว้แล้วจากการทดสอบ Q-1/Q-2a รอบก่อน ๆ (`state='passed'`,
  `first_passed_at='2026-07-29 09:05:49'`) — "Lesson finished" banner ที่เห็นบน
  หน้าจึงเป็น state เดิมที่ไม่เกี่ยวกับ ticket นี้เลย (ไม่ใช่ regression, ticket
  นี้ไม่ได้เรียก `finish()` เลยระหว่างทดสอบ) — ไม่ต้อง cleanup เพราะเป็น
  `lesson_progress` ที่ Q-1/Q-2a's evidence ทิ้งไว้ตั้งแต่ก่อนหน้านี้แล้ว

### Review focus

- ทำไม `attemptStatus`/`pendingAttemptsRef` ต้องอยู่ที่ `LessonView` (parent) แทนที่
  จะให้ `RecallCheckCard` เก็บ save-status ของตัวเองเหมือนที่ Q-1 เก็บ
  stage/selection ไว้ในการ์ดเอง? เกี่ยวอะไรกับ aggregate banner เหนือปุ่ม Finish?
- ทำไม dedup guard ของ `onAttemptReady` (`lastReportedAttemptRef`) ต้องเก็บ
  "signature" ของ attempt แทนที่จะเก็บ boolean "เคย submit แล้วหรือยัง" ตรง ๆ?

> คำถามข้อสองของรอบแรก ("ทำไมเปลี่ยน rating ถึงต้อง resubmit ทันที ไม่ใช่
> ignore/ค่าแรกชนะ") **ล้าสมัยบางส่วนหลังรอบ 3**: คำตอบเดิม (resubmit เป็น
> honest data point) ยังถูก แต่ "ทันที" ไม่ถูกต้องแล้ว — ตอนนี้ submit จริง
> เกิดหลัง debounce settle เท่านั้น ดู Review focus รอบ 3 ด้านล่างสำหรับ
> คำถามที่แทนที่

### จงใจไม่ทำในรอบนี้

- ไม่แก้ granular error message ต่อ HTTP status ของ attempts POST — `{kind:
  "error"}` เดียวพอสำหรับ "Not saved" + Retry ตาม scope ที่ ticket กำหนด
  (`400`/`404`/`413`/`500` แสดงผลเหมือนกันหมดจากมุมผู้ใช้)
- ไม่มี "your attempt history" view, ไม่แตะ `review_cards`/`review_logs`/SM-2 —
  Q-2c
- ไม่เปลี่ยน 3-stage flow เดิมของ Q-1 เลย
- ไม่แก้ `finish()`'s identity guard เดิม (UX-5) ให้ใช้ `loadGenerationRef`
  ด้วย — เหตุผลไม่ใช่ "reviewer ไม่ได้ชี้จุดนั้น" (นั่นเป็นแค่ข้อสังเกตเชิง
  กระบวนการ ไม่ใช่เหตุผลที่จะยังสมเหตุสมผลอีก 6 เดือนข้างหน้า) แต่เป็นเพราะ
  **blast radius ต่างกันจริง**: response ที่ค้างของ `finish()` ที่มาช้าเกี่ยวข้อง
  กับ**แค่ lesson เดียว**ที่มันเริ่มขึ้นมา (เขียน `passed` ให้ lesson นั้น) —
  `ok` ที่มาช้าทำให้เห็น "Lesson finished" ตอนกลับมาที่ lesson เดิม ซึ่ง**ตรงกับ
  truth ฝั่ง server อยู่แล้ว** (lesson นั้น pass จริง) แถมโหลดใหม่ก็จะเจอ banner
  เดิมผ่าน `alreadyPassed` อยู่ดี ส่วน `error` ที่มาช้าก็แค่โชว์ "Could not save"
  ที่ปุ่ม Retry ของมัน scroll ไปหา check แรกที่ยังไม่เสร็จแทนที่จะ retry จริง —
  งงแต่ไม่ทำข้อมูลเสียหาย ไม่มีอะไรเหมือนกรณี attempt เลยที่ stale `ok` ลบสถานะ
  fail จริงทิ้งแล้วซ่อนว่าแถวหายไป — เมื่อ blast radius ต่างกันขนาดนี้ การเลื่อน
  ไปทำทีหลังจึงสมเหตุสมผล ไม่ใช่แค่ "ยังไม่มีใครสั่งให้แก้". ผลข้างเคียงที่ต้อง
  จำไว้: `lesson/page.tsx` ตอนนี้มี identity mechanism สองแบบข้าง ๆ กันจริง (ดู
  หมายเหตุใน "การตัดสินใจหลัก" ด้านบน) — ติดตามที่ `docs/roadmap.md`

### รอบ code-reviewer (REQUEST_CHANGES → แก้ครบ)

Reviewer รัน mutation set ของตัวเอง **14 จุดตามที่ ticket เรียกร้อง ตายหมดจริง**
แต่เพิ่มอีก **5 จุดที่รอด** จาก edge case ที่รอบแรกไม่ได้ทดสอบ:

- **R1/R2 (blocking, แก้แล้ว) — identity guard เทียบ slug ไม่ใช่ "visit"**:
  `currentIdentityRef` (topic/concept equality) แยกไม่ออกระหว่าง "กลับมาที่
  lesson เดิมอีกครั้ง" กับ "ยังเป็น visit เดิม" — reviewer พิสูจน์สดด้วย A → B →
  กลับมา A: attempt ของ visit แรกค้างอยู่ (POST ช้า, `AbortSignal.timeout` 10
  วินาทีทำให้เรื่องนี้เป็นเรื่องปกติบนมือถือ ไม่ใช่ edge case แปลก) พอกลับมา A
  (visit ใหม่) แล้ว POST ของ visit ใหม่ fail จริง (เห็น "Not saved" + banner
  ถูกต้อง) แต่พอ visit แรกที่ค้างอยู่ resolve เป็น `{kind:"ok"}` ทีหลัง —
  `stillCurrent()` เดิม (เทียบ slug) คืน **true** เพราะ topic/concept ของสอง
  visit เหมือนกัน → เขียนทับสถานะ "error" จริงด้วย "saved" ปลอม (R1), และใน
  variant ที่ visit สองไม่แตะอะไรเลย stale error ของ visit แรกก็ leak เข้ามา
  โดยไม่มี Retry ให้กด เพราะ `pendingAttemptsRef` ถูก reset ไปแล้วตอนเปลี่ยน
  lesson (R2) — ตรงกับ bug class เดียวกับที่ UX-6 เคยเจอมาก่อน (slug equality
  ข้าม lifecycle boundary ไม่ใช่ identity ที่แท้จริง). แก้ด้วย `loadGenerationRef`
  (ref นับเลข bump ทุกครั้งที่ effect โหลด lesson รันใหม่ **รวมถึงการโหลด
  lesson เดิมซ้ำ**) — `submitAttempt` capture generation ตอนเริ่ม แล้วเทียบ
  `loadGenerationRef.current === generation` แทน slug equality เดิม. **ไม่แก้
  `finish()`'s guard เดิม (UX-5) ให้ใช้กลไกเดียวกัน** แม้จะมี bug class เดียวกัน
  ในทางทฤษฎี — เหตุผลคือ blast radius ต่างกันจริง ไม่ใช่แค่ "อยู่นอก scope"
  (ดูรายละเอียดเต็มที่หัวข้อ "finish()" ท้ายไฟล์นี้): response ค้างของ
  `finish()` เกี่ยวกับแค่ lesson เดียวที่เริ่มมันขึ้นมา และผลลัพธ์ที่มาช้าไม่ว่า
  จะ ok หรือ error ก็ไม่ทำข้อมูลเสียหายหรือซ่อนอะไรเงียบ ๆ เหมือนกรณี attempt —
  บันทึกไว้เป็นความเสี่ยงที่รู้ตัวแล้วในหัวข้อ "จงใจไม่ทำ" ด้านล่าง พร้อม
  follow-up ใน `docs/roadmap.md`
- **R3 (blocking, แก้แล้ว) — retry ไม่เคย assert request body**: มูเทต
  `handleRetrySave` ให้ resend `{...stored, confidence:"guessed"}` แทน `stored`
  ตรง ๆ แล้ว `npm test` ยังเขียว 161/161 — เทสต์เดิมเช็คแค่ "เรียก postAttempt
  2 ครั้ง" กับ "indicator หาย" ไม่เคยเทียบว่า call ที่สองส่งอะไรจริง ทั้งที่
  ตารางเป็น append-only ไม่มี idempotency key: retry ที่ผิด payload = แถวผิด
  ถาวร ไม่ใช่ glitch ชั่วคราว. เพิ่ม `expect(postAttempt.mock.calls[1]).toEqual
  (postAttempt.mock.calls[0])` เข้าไปในเทสต์เดิม
- **R4 (blocking, แก้แล้ว) — "saving" state ทั้งก้อนไม่มีเทสต์คุ้มครองเลย**:
  มูเทต 3 จุดตายทั้งหมดหลังแก้ (ก่อนแก้ **ไม่มีจุดไหนถูกจับเลยสักจุด**): (a) ลบ
  `setAttemptStatus(...,"saving")` ใน `submitAttempt`, (b) ลบ `"Saving…"` render
  ใน `RecallCheckCard`, (c) ลบ `if (attemptStatus[position]==="saving") return;`
  ใน `handleRetrySave`. (a)/(b) แก้ด้วยเทสต์ใหม่ที่ pin สถานะ "saving" ตรง ๆ ทั้ง
  ที่ระดับ `RecallCheckCard` เดี่ยว ๆ และที่ระดับ `LessonPage` เต็มระบบ (deferred
  promise คุม timing). **(c) พิสูจน์แล้วว่าเป็น equivalent mutant จริง ไม่ใช่
  survivor ที่ต้องแก้ — แต่เหตุผลที่เขียนไว้รอบแรกผิด แก้แล้วในรอบ 3**:
  รอบแรกอ้างว่า "React commit สถานะ saving ก่อน browser event ถัดไปจะ
  ประมวลผลได้เสมอ" โดยอ้าง probe ที่ใช้ `fireEvent.click` สองครั้งติดกัน — แต่
  `fireEvent.click` ของ RTL ห่อด้วย `act()` ซึ่ง flush + unmount ปุ่มให้เองก่อน
  คลิกที่สองจะไปถึง — นั่นพิสูจน์แค่พฤติกรรมของ **RTL's fireEvent** ไม่ใช่
  พฤติกรรมจริงของ browser. ทดสอบซ้ำด้วย raw `dispatchEvent` สองครั้งติดกัน
  **โดยไม่ห่อ** `act()`: ปุ่มยังอยู่ใน DOM ระหว่างสองคลิก และ**ทั้งสองคลิกไปถึง
  handler จริง** (`postAttempt` ถูกเรียก 3 ครั้ง) เพราะ `handleRetrySave`'s
  closure ที่ผูกกับปุ่มยัง capture ค่า `attemptStatus` เดิม (snapshot ของ
  render ก่อนหน้า) จนกว่า React จะ re-render จริง — **การ์ดด้วย closure เอง
  ไม่ airtight**. ข้อพิสูจน์ที่แน่นอนจริงคือ **render condition**: ปุ่ม Retry
  render ⟺ `saveStatus === "error"` เท่านั้น (เทสต์อยู่แล้ว: "does not show the
  Retry affordance when saveStatus='saved'") ประกอบกับข้อเท็จจริงที่ว่า browser
  จริงจะ flush แต่ละ discrete event (click) ให้ synchronous state update เสร็จ
  ก่อนประมวลผล event ถัดไปเสมอ (ไม่มีช่องให้สองคลิกจริงของผู้ใช้ไปถึง handler
  พร้อมกันในทางปฏิบัติ) — บทสรุป "equivalent mutant" ยังถูกต้องเหมือนเดิม แต่
  เหตุผลที่ถูกต้องคือ render condition + browser's discrete-event flushing
  ไม่ใช่ act()'s side effect ในเทสต์ตัวเดียว. เก็บ guard นี้ไว้เป็น
  defense-in-depth (เผื่ออนาคตเปลี่ยนให้ Retry โชว์ตอน "saving" ด้วย เหมือน
  pattern ของปุ่ม Finish ที่ใช้ `aria-disabled` ไม่ใช่ลบปุ่มทิ้ง) แต่ **ไม่มี
  เทสต์ไหน kill มันได้จริงในโครงสร้างปัจจุบัน** — บันทึกตรง ๆ
  แทนที่จะเสแสร้งว่ามี. เพิ่มเติม: "Finish อ่านว่าเสร็จทั้งที่ attempt ยังค้างอยู่"
  พิสูจน์สดว่าเป็นจริง (banner รวมเดิมเช็คแค่ `"error"` ไม่เช็ค `"saving"`) — แก้
  โดยเพิ่ม `savingAttemptCount` + banner แยก ("Saving N recall attempts…", role
  ="status", สีกลาง ไม่ใช่สี danger เพราะไม่ใช่ error) แสดงคู่กับ "Lesson
  finished" ได้พร้อมกัน ไม่ใช่แทนที่กัน — เพิ่มเทสต์ยืนยันทั้งสอง banner โชว์
  พร้อมกันจริง
- **R5 (blocking, แก้แล้ว) — flip-flop รัวๆ เขียนหลายแถว รวมถึงแถวซ้ำไบต์ต่อไบต์**:
  พิสูจน์สดว่า Pass/Not yet/Pass/Not yet บนการ์ดเดียวเขียน **4 แถว** ใต้
  `check_key` เดียวกัน โดยแถว 1&3 และ 2&4 เหมือนกันทุกไบต์ — เอกสารรอบแรกแก้ตัวว่า
  "การแก้ไขคือข้อมูลใหม่ที่มีค่า" ซึ่งใช้ได้กับการแก้ไข**ครั้งแรก**เท่านั้น ไม่ใช่
  แถว 3/4 ที่เป็นแค่ความลังเลของ UI. แก้ด้วย **debounce ก่อนส่ง ไม่ใช่กันการย้อน
  กลับไปค่าเดิม** (ทางเลือกหลังทำให้ "แถวล่าสุด" ไม่ใช่คำตอบล่าสุดของ user จริง
  ซึ่งแย่กว่า) — `handleAttemptReady` เลื่อนการ submit ของ short_answer ออกไป
  **1000ms** (mcq ไม่ต้องเพราะ input freeze ทันทีที่ reveal, มีค่าเดียวเสมอ) เคลียร์
  timer เดิมทุกครั้งที่มีการเปลี่ยนแปลงใหม่ (`clearTimeout` + `setTimeout` ใหม่)
  ทำให้มีแค่ค่าสุดท้ายที่ settle จริงเท่านั้นถูกส่ง. Flush ทันที (ข้าม debounce)
  ใน 3 จุดตามที่ ticket เรียกร้อง: (1) quiet period หมดเวลาเอง (2) `Finish` ถูก
  กด (`flushAllPendingAttempts()` ก่อนเรียก `finish()`) (3) navigate ออก (effect
  cleanup ของ per-navigation effect เรียก `flushAllPendingAttempts()` **ก่อน**
  effect ใหม่จะ null `currentIdentityRef`/bump generation ทำให้ยังส่งด้วย
  identity/generation ของ lesson เดิมได้ถูกต้อง) — ทดสอบครบทั้ง 3 เส้นทาง
  รวมถึง "ตอบแล้วออกจากหน้าเร็ว ๆ ก็ยังเขียนแค่แถวเดียว" ด้วย fake timers
  (`vi.useFakeTimers()` + `vi.advanceTimersByTimeAsync`) และพิสูจน์สดด้วย
  Playwright ว่า flip-flop จริงเขียนแค่ **1 แถว** พร้อมค่าที่ settle จริง (ดู
  Live evidence รอบ 2 ด้านล่าง). **ข้อควรระวัง (เขียนให้ตรงในรอบ 3): debounce
  ลดความเสี่ยงเรื่องลำดับ ไม่ได้ปิดมันสนิท** — สองครั้งที่แก้ rating ห่างกันจริง
  เกิน 1 วินาที (คนละหน้าต่าง debounce กันคนละครั้ง) โดยที่ request แรกช้ากว่า
  1 วินาที (เช่น network ช้า) ยังทำให้มีสอง POST ค้างอยู่พร้อมกันได้ ซึ่ง
  **ลำดับที่มาถึง (arrival order) ต่างหาก**ที่กำหนดทั้ง `created_at` และ `id`
  ไม่ใช่ลำดับที่ user กดจริง — แปลว่า `ORDER BY created_at DESC, id DESC` (ที่
  บันทึกเป็นข้อกำหนดของ Q-2c ใน `roadmap.md`) ให้ arrival order ไม่ใช่ user
  order เสมอไป. **flush ที่ล้มเหลว (navigate away/pagehide) ก็หายเงียบ ๆ โดย
  ไม่มี indicator หรือ Retry ให้กด**: cleanup ส่ง request ออกไปด้วย
  identity/generation ของ lesson เดิมถูกต้อง (ดูด้านบน) แต่พอผลลัพธ์กลับมา —
  ไม่ว่าจะสำเร็จหรือ fail — generation ใหม่ (จาก lesson ถัดไปที่โหลดไปแล้ว) ทำให้
  `stillCurrent()` ทิ้งผลลัพธ์นั้น และ effect ใหม่ก็ clear `pendingAttemptsRef`
  ของ lesson เดิมไปแล้วด้วย ไม่มีที่ให้ Retry อีก — เป็น trade-off ที่ยอมรับ
  (ไม่มี UI ให้ retry สำหรับ lesson ที่ไม่ได้อยู่บนจอแล้วจริง ๆ) ไม่ใช่บั๊กที่
  ไม่รู้ตัว
- **R6 (แก้แล้ว) — คำกล่าวอ้าง 3 จุดที่ไม่จริง**:
  - "ทั้งสอง POST ยัง insert แถวถูกต้องเสมอ กระทบแค่ transient UI display" —
    **เท็จ**: `created_at` เป็น `TIMESTAMP` second-precision
    (`migrations/006_recall.sql:19`) และ query ที่ Q-2c จะใช้คือ `WHERE
    check_key = ? ORDER BY created_at DESC` — reviewer พิสูจน์ว่าสองแถวที่
    submit ห่างกันจริงในเวลาปกติ (ไม่ใช่ race) ตกวินาทีเดียวกันได้จริง ทำให้
    "แถวล่าสุด" อ่านผิดได้ (Pass→Not yet ในวินาทีเดียวกันอ่านกลับมาเป็น Pass) —
    แก้คำอธิบายเป็น data-ordering risk ไม่ใช่ display glitch, และบันทึกเป็น
    ข้อกำหนดของ Q-2c ใน `docs/roadmap.md` (`ORDER BY created_at DESC, id DESC`
    — **ไม่แก้ migration ในรอบนี้**, R5's debounce ปิด root cause นี้ไปแล้วสำหรับ
    flip-flop โดยเฉพาะอยู่แล้วด้วย)
  - "Q-2c อ่าน `created_at DESC` ล่าสุดอยู่แล้วไม่กระทบ" — อ้างพฤติกรรมของ
    ticket ที่**ยังไม่เริ่ม** (`Q-2c — not started` อยู่ไม่กี่บรรทัดถัดจากนี้เอง
    ในไฟล์เดียวกัน) เป็น mitigation ที่มีอยู่จริง — ลบประโยคนี้ทิ้ง
  - `RecallCheckCard.tsx`'s comment เดิมอ้างว่า guard กัน "React StrictMode's
    double-invoked effects" — **StrictMode double-invoke เฉพาะตอน mount**, effect
    นี้ early-return ตอน mount (`stage==="recall"`) แล้วมารายงานจริงตอน state
    เปลี่ยนทีหลังซึ่ง StrictMode ไม่ double-invoke ซ้ำ — ยืนยันเชิงประจักษ์ว่า
    ปิด dedup check แล้ว StrictMode test **ยังผ่าน** (มีแค่เทสต์ "re-render
    เปลี่ยน callback identity" เท่านั้นที่จับได้จริง) — M1 row ของตารางมูเทชัน
    ด้านล่างยอมรับเรื่องนี้อยู่แล้ว แต่ comment/prose ยังพูดตรงข้าม แก้ให้ตรงกัน

  **หมายเหตุ (พบใน code review รอบ 3) — สองข้อบนนี้เขียนไว้ว่า "แก้แล้ว" ทั้งที่
  แก้ไม่ครบจริง**: ข้อ (1) แก้แค่โค้ดจริง (`postAttempt`'s comment/behavior) แต่
  ประโยค "เป็น data problem ที่ Q-2c อ่าน `created_at DESC` ล่าสุดอยู่แล้วไม่
  กระทบ" ยังค้างอยู่ verbatim ในหัวข้อ "การตัดสินใจหลัก" ด้านบน (ประโยคเดียวกับ
  ข้อ (2) ที่บอกว่า "ลบทิ้ง" แต่ไม่ได้ลบจริง). ข้อ (3) แก้แค่ comment ในโค้ด
  (`RecallCheckCard.tsx`) แต่**ไม่ได้แก้ prose** ในหัวข้อ "การตัดสินใจหลัก" ที่ยัง
  เขียนว่า "Guard นี้กัน 3 เหตุการณ์พร้อมกัน: React StrictMode double-invoke..."
  ตรงข้ามกับ M1's ผลจริง 12 บรรทัดถัดไปในตารางเดียวกัน — เป็นรูปแบบเดียวกับที่
  Q-1/Q-2a เจอมาก่อน (claim ว่าแก้แล้วในเอกสาร แต่แก้ไม่ครบทุกจุดที่พูดถึงเรื่อง
  เดียวกัน) **แก้ครบแล้วจริงในรอบ 3 นี้** — ทั้งสองจุดถูกแก้ที่ต้นตอในหัวข้อ
  "การตัดสินใจหลัก" ด้านบนแล้ว (ไม่ใช่แค่ที่นี่)
- **R7 (แก้แล้ว) — comment เกินจริง/ล้าสมัย**: `RecallCheckCard.tsx`'s comment
  บน `lastReportedAttemptRef` ตัดเหลือประโยคเดียว (resubmit policy) ตัด
  StrictMode parenthetical ที่ผิดทิ้ง; comment บน prop `onAttemptReady` ที่พูดว่า
  "fires again each time rating changes = new honest data point, save it"
  ตัดทิ้งทั้งหมด เพราะ R5 เปลี่ยนพฤติกรรมจริงแล้ว (ไม่ใช่ทุกครั้งที่เปลี่ยนจะกลาย
  เป็นแถวใหม่อีกต่อไป — เฉพาะค่าที่ settle) ชื่อ prop + type `CompletedAttempt`
  พอเป็น self-documenting โดยไม่ต้องมี comment ผิด ๆ
- **Also fixed (เล็ก ๆ ทั้งหมด)**: ลบ dead branch ใน `submitAttempt`'s catch
  (`postAttempt` ไม่ throw อะไรนอกจาก `UnauthorizedError` แล้ว จึง `if
  (!stillCurrent()) return; setAttemptStatus(...)` หลัง Unauthorized check
  เดิมเข้าไม่ถึงได้เลย); `submitAttempt`/`handleAttemptReady`'s silent-drop path
  (`!identity`/`!check`) เปลี่ยนจาก `return` เฉย ๆ เป็น mark `"error"` ก่อน
  return (ทั้งสอง edge case แทบเป็นไปไม่ได้ในทางปฏิบัติ วิเคราะห์แล้วว่า
  unreachable ผ่าน UI จริง จึง**ไม่มีเทสต์ pin** — บันทึกตรง ๆ); export
  `CompletedAttempt` จาก `RecallCheckCard` แทน inline type ซ้ำใน
  `handleAttemptReady` (กัน field ใหม่ถูก spread เข้า runtime แล้วเงียบ ๆ หายไป
  ที่ `submitAttempt`'s explicit field list); ลบ `type Confidence =
  AttemptConfidence` ที่เป็น pass-through เปล่า ๆ ใช้ `AttemptConfidence` ตรง ๆ;
  แก้ `api.ts`'s comment ที่อ้างว่า "callers here (RecallCheckCard/lesson page)"
  ทั้งที่ `RecallCheckCard` ไม่เคยเรียก `postAttempt` เอง; เพิ่มเทสต์
  slug-encoding ให้ `postAttempt`

### Mutation table รอบ 2 (14/14 required + 5 reviewer-found — 4 killed, 1 equivalent)

รัน `npx vitest run` เต็มชุดหลังแก้แต่ละจุด แล้ว revert ทุกครั้ง — รวม re-verify
มูเทชัน 11 จุดจากรอบแรก (M1-M11 เดิม, ดูตารางรอบแรกด้านบน) **ผ่านซ้ำทุกจุดหลัง
restructure** (เช็คจุดที่โค้ดย้ายที่จริง: M6/M9 ย้ายเข้า `submitAttempt`/
`handleAttemptReady` ใหม่ — kill ซ้ำได้เหมือนเดิม, M2 spot-check เพิ่มเติมใน
`RecallCheckCard` ที่ไม่ถูกแก้เลยในรอบนี้ — kill เหมือนเดิม):

| # | Mutation | ผลลัพธ์ |
|---|---|---|
| R1/R2 | `stillCurrent()` ใน `submitAttempt` เทียบ slug (`currentIdentityRef`) แทน `loadGenerationRef` | killed — 2 เทสต์ใหม่: "a stale OK response from an earlier visit..." และ "a stale failure response from an earlier visit..." (ทั้งคู่ A→B→A) |
| R3 | `handleRetrySave` ส่ง `{...stored, confidence:"guessed"}` แทน `stored` ตรง ๆ | killed — assertion ใหม่ `expect(postAttempt.mock.calls[1]).toEqual(postAttempt.mock.calls[0])` ในเทสต์ retry เดิม |
| R4a | ลบ `setAttemptStatus(...,"saving")` ใน `submitAttempt` | killed — 2 เทสต์ ("shows the pending-count banner..." + "surfaces a still-saving attempt...") |
| R4b | ลบ `{saveStatus === "saving" && <p>Saving…</p>}` ใน `RecallCheckCard` | killed — 2 เทสต์ (component-level ใหม่ "shows 'Saving…' on saveStatus='saving'" + page-level "shows the pending-count banner...") |
| R4c | ลบ `if (attemptStatus[position]==="saving") return;` ใน `handleRetrySave` | **equivalent mutant ยืนยันแล้ว** — เหตุผลที่ถูกต้อง (แก้จากรอบแรกที่อ้างผิด, ดู R4's ย่อหน้ายาวด้านบนสำหรับรายละเอียดเต็ม + การพิสูจน์ด้วย raw `dispatchEvent` ที่ไม่ห่อ `act()`): ปุ่ม Retry render ⟺ `saveStatus==="error"` (เทสต์อยู่แล้ว) + browser จริง flush แต่ละ discrete click ให้เสร็จก่อนประมวลผลอันถัดไปเสมอ — ไม่แก้โค้ด เก็บ guard ไว้เป็น defense-in-depth |
| R6c | comment เดิมอ้าง StrictMode ผิด (guard ไม่ได้มาจาก StrictMode จริง) | ไม่ใช่ mutation — แก้ comment/doc ให้ตรงกับ M1's ผลจริงที่มีอยู่แล้ว |
| M6 (re-verify) | `handleAttemptReady` ส่ง `check.expected_answer` แทน `check.question` (ย้ายที่แล้วยัง kill ได้) | killed — 4 เทสต์ (เพิ่มจาก 2 เดิม เพราะเทสต์ debounce ใหม่ก็แตะจุดเดียวกัน) |
| M9 (re-verify) | 401 redirect ใน `submitAttempt`'s catch (ย้ายที่แล้วยัง kill ได้) | killed — 1 เทสต์เหมือนเดิม |
| M2 (spot-check) | `isMcqCorrect` invert ใน `RecallCheckCard` (ไม่ถูกแก้ในรอบนี้เลย ยืนยันไม่ regress) | killed — 4 เทสต์เหมือนเดิม |

**สรุป: 14/14 มูเทชันที่ ticket รอบ 2 เรียกร้องตายหมด, 5 จุดที่ reviewer เจอเพิ่ม
4 killed จริง 1 equivalent (บันทึกไว้ตรง ๆ ไม่ใช่ซ่อน), 3 จุด re-verify จากรอบแรก
ไม่มี regression จากการ restructure**

### Live verification รอบ 2 (docker compose up -d --build, ไม่ใช้ -v)

Stack เดิม, lesson `domain-driven-design/ubiquitous-language` เดิม (5 recall
checks: short_answer/mcq/mcq/short_answer/mcq):

- **Flip-flop เขียนแค่ 1 แถว พร้อมค่าที่ settle จริง**: Pass → Not yet → Pass
  รัว ๆ บนการ์ดเดียว (`confidence=Guessed`) รอ 1.8 วินาที (เกิน debounce
  window 1s) → `SELECT` เห็นแค่ **1 แถวใหม่** (`guessed, correct` — ตรงกับ
  Pass ตัวสุดท้ายที่กด ไม่ใช่ตัวกลางที่เป็น Not yet)
- **A → B → A ผ่าน client-side navigation จริง (ปุ่ม Next แล้วกด Back ของ
  browser จริง ๆ — ไม่ใช่ `page.goto` ซึ่งจะ abort request ที่ค้างอยู่แทนที่จะ
  ทดสอบอะไร)**: ตอบ mcq ที่ visit แรก (POST ถูก intercept หน่วงเวลา 4 วินาที
  ผ่าน Playwright route), กด "Model-Driven Design →" (Next, client-side จริง
  — ยืนยัน URL เปลี่ยนเป็น `concept=model-driven-design`), กด browser Back
  (ยืนยัน URL กลับมา `concept=ubiquitous-language` — visit ใหม่, generation
  bump แล้ว), ไม่แตะอะไรใน visit สอง, รอให้ POST ของ visit แรกที่ค้างอยู่
  resolve (`{kind:"ok"}` จริง — แถวถูกเขียนจริงใน DB ยืนยันด้วย `SELECT`
  ทีหลัง) → **banner รวม "recall attempts didn't save" count = 0, "Not saved"
  count = 0** ทั้งคู่ — stale response ของ visit แรกไม่ leak เข้ามาใน visit
  สองที่ไม่ได้แตะอะไรเลยจริง (ตรงกับ R2's scenario เป๊ะ)
- **Finish อ่านว่าเสร็จพร้อมกับ attempt ที่ยังค้างอยู่ ยังโชว์ทั้งคู่พร้อมกัน**:
  lesson นี้ถูก mark `passed` ไว้แล้วจากการทดสอบ Q-1/Q-2a รอบก่อน — reset
  ชั่วคราวเป็น `in_progress` ด้วย SQL ตรง ๆ (บันทึกค่าดั้งเดิมไว้ก่อน) เพื่อทดสอบ
  เส้นทางกด Finish จริง, ตอบครบทั้ง 5 checks เร็ว ๆ (POST ทุกตัวถูกหน่วง 3
  วินาที), กด "Finish lesson" → banner "Lesson finished" โผล่จริง **พร้อมกับ**
  banner "Saving N recall attempt(s)…" ที่ยังอยู่ (count = 1 พอดี ตอนเช็ค) —
  ยืนยันว่า Finish ไม่ทำให้ attempt ที่ยังค้างอยู่หายไปจากสายตาผู้ใช้เงียบ ๆ
  อีกต่อไป — restore `lesson_progress` กลับเป็นค่าดั้งเดิมทันทีหลังทดสอบ
  (`state='passed'`, `first_passed_at` เดิมเป๊ะ, ยืนยันด้วย `SELECT` ก่อน/หลัง)
- **Cleanup**: ลบแถวที่ script นี้สร้างทั้งหมด (`id > 18`, รวม 7 แถวจากรอบสุดท้าย
  ที่รันสำเร็จ — บวกแถวจาก debug run ที่ครัชระหว่างพัฒนา script ซึ่งลบไปแล้ว
  ก่อนรันรอบสุดท้ายด้วย) ยืนยันด้วย `SELECT` ก่อน (`COUNT=20, MAX(id)=61` จาก
  debug runs) และหลัง (`COUNT=13, MAX(id)=18` — กลับสู่สภาพเดิมของ Q-2a's dev
  data เป๊ะ)

### Review focus (รอบ 2)

- ทำไม `submitAttempt`'s identity guard ต้องเปลี่ยนจากเทียบ slug
  (`currentIdentityRef`) เป็นนับ "visit" (`loadGenerationRef`) ทั้งที่
  `finish()` ข้างล่างยังใช้ slug comparison เดิมอยู่ไม่ได้แก้?
- ทำไมมูเทชัน "ลบ `if (attemptStatus[position]==="saving") return;`" ถึงไม่มี
  เทสต์ไหน kill ได้เลย ทั้งที่โค้ดบรรทัดนี้ยังคงอยู่ในไฟล์ — เป็น coverage gap
  จริงหรือเป็น equivalent mutant? อะไรคือหลักฐานที่แยกสองเรื่องนี้ออกจากกันได้?
- ทำไมการแก้ flip-flop (R5) ต้องเป็น debounce-ก่อนส่ง แทนที่จะเป็น "กันการ
  ย้อนกลับไปค่าเดิมที่เคย submit แล้ว"?

### รอบ code-reviewer 3 (NO-SHIP → แก้ครบ)

Reviewer ยืนยันว่า `loadGenerationRef`/A→B→A/StrictMode-bump-before-capture/
navigate-away-flush-before-bump/flip-flop-one-row/retry-body/ทั้งสอง "saving"
mutant **รอดการโจมตี 5 แบบที่ reviewer ลองเองแล้วจริง ไม่ regress** — แต่พบ
บั๊กใหม่ 1 จุดที่**รอบ 2 เป็นคนสร้างขึ้นเอง** (debounce) กับ test gap 1 จุด และ
window ใหม่ 1 จุดที่ debounce เปิดขึ้นมาโดยไม่ตั้งใจ:

- **R1 (blocking, แก้แล้ว) — Retry แข่งกับ debounce timer ที่ยังไม่หมดเวลา เขียน
  แถวซ้ำ**: `handleRetrySave` เดิมเรียก `submitAttempt` ตรง ๆ ด้วย
  `pendingAttemptsRef.current[position]` โดยไม่เคยดู `debounceTimersRef` เลย
  — Retry render ได้ตอน `saveStatus==="error"` ซึ่ง**ไม่เปลี่ยน**ตอนผู้ใช้แก้ไข
  rating ใหม่ (การแก้ไขแค่ schedule debounce timer ใหม่ ไม่แตะ `attemptStatus`)
  ทำให้ปุ่ม Retry ค้างอยู่บนจอพร้อมกับ timer ที่กำลังนับถอยหลังพร้อมกันได้จริง
  — พิสูจน์สด: Pass → 1.3s → POST#1 → 500 → "Not saved"+Retry, คลิก "Not yet"
  (schedule timer ใหม่, Retry ยังอยู่), คลิก "Retry" → POST#2 (incorrect),
  อีก 1s timer เดิมที่ยังไม่ถูกยกเลิกยิง → POST#3 (incorrect) — **สองแถว
  ไบต์ต่อไบต์เหมือนกันใต้ check_key เดียวกัน** ตรงกับ class ที่ debounce
  (R5 ของรอบ 2) ตั้งใจกำจัด แต่หลุดกลับมาทาง Retry แทน. แก้โดยดึง
  `clearTimeout`/`delete debounceTimersRef.current[position]` ออกมาเป็น
  `clearPendingTimer(position)` ใช้ร่วมกัน, `flushPendingAttempt` เรียก
  `clearPendingTimer` **แบบไม่มีเงื่อนไข** (เดิมมีเงื่อนไข "ต้องมี timer อยู่
  ก่อนถึงจะทำอะไร" ซึ่งทำให้ใช้กับ Retry ไม่ได้เพราะ Retry ไม่มี timer ของตัวเอง
  เสมอไป) แล้ว `handleRetrySave` เรียก `flushPendingAttempt(position)` แทนการ
  เรียก `submitAttempt` ตรง ๆ — "flush-then-submit-once" แทนที่จะปล่อยให้แข่งกับ
  timer ที่ยังไม่หมดเวลา
- **R2 (blocking, แก้แล้ว) — เทสต์ flush ทุกตัวใช้ lesson ที่มี check เดียว ซ่อน
  false-success mutant 2 จุด**: `flushPendingAttempt` ส่ง
  `Object.values(pendingAttemptsRef.current)[0]` แทน `[position]`, และ
  `flushAllPendingAttempts` flush แค่ timer แรก — ทั้งสองมูเทตแล้ว **รอด
  ทั้ง 161 เทสต์เดิม** เพราะไม่มีเทสต์ไหนมี 2 check ค้างพร้อมกันเลย. โค้ดจริง
  ถูกอยู่แล้ว (reviewer พิสูจน์ด้วย probe 2-check ว่าทั้งสอง flush path ส่งครบ
  ทั้งสอง check พร้อม `question` ของตัวเอง) แต่ถ้ามูเทชันนี้หลุดขึ้น production
  จริง จะเขียน**แถวซ้ำของ check #1 ซ้ำ, ไม่เขียน check #2 เลย, แต่ยัง mark
  check #2 ว่า "saved"** — false success พ่วง silent data loss ที่ไม่มีทางรู้
  ตัวเลย. แก้ด้วย fixture ใหม่ `twoShortAnswerLesson` (2 short_answer checks)
  — rate ทั้งสองภายใน debounce window เดียวกัน แล้วยืนยันว่า flush ทั้งทาง
  Finish และทาง navigate-away ส่งครบทั้งสอง `question` พร้อมค่าที่ถูกต้องของ
  ตัวเอง (ไม่ใช่แค่นับจำนวนครั้งที่เรียก)
- **R3 (blocking, แก้แล้ว) — reload/ปิดแท็บภายใน debounce window ทำ attempt
  หายเงียบ ๆ โดยไม่มี indicator เลย, และเป็น window ใหม่ที่ debounce สร้างขึ้น**:
  ก่อน R5 (รอบ 2) POST ยิงตอนคลิกทันที — การ reload เร็ว ๆ หลังคลิกไม่เคยเป็น
  ปัญหา. หลัง debounce, attempt ค้างอยู่ใน memory เป็นเวลาถึง 1 วินาทีก่อนจะยิง
  จริง — พิสูจน์สด: rate check แล้ว reload 150ms ทีหลัง → **0 แถวถูกเขียน**
  ไม่มี indicator ไม่มีอะไรเลย ทั้งที่ `quiz.md`'s "flush 3 เส้นทาง" (quiet
  period/Finish/navigate away) ไม่พูดถึง hard navigation เลย ทำให้คนอ่านสรุปผิด
  ว่า rating ปลอดภัยเสมอ. แก้ที่โค้ดจริง ไม่ใช่แค่บันทึกไว้: เพิ่ม effect ฟัง
  `pagehide` และ `visibilitychange` → `hidden` (คู่นี้เชื่อถือได้กว่า
  `beforeunload` เดี่ยว ๆ ซึ่งไม่น่าเชื่อถือบน mobile Safari — อุปกรณ์ที่ผู้ใช้
  รีวิวจริง) เรียก `flushAllPendingAttempts({keepalive:true})`; `postAttempt`
  รับ `options?: {keepalive?: boolean}` เพิ่ม forward เข้า `fetch`'s
  `keepalive` (ไม่ใช้ `navigator.sendBeacon` เพราะมันแนบ `Authorization`
  header ไม่ได้) — request ที่ยิงออกไปก่อนหน้าจะรอดแม้หน้าจะ unload ไปแล้ว.
  `submitAttempt` ส่ง `postAttempt` แบบ 3 argument เดิมทุกจุดที่ไม่ต้องการ
  keepalive (ไม่ใช่ 4 argument พร้อม `undefined`) เพื่อไม่ให้เทสต์เดิมที่
  assert exact call shape พังทั้งหมด. **หมายเหตุตรง ๆ (residual window ที่ยัง
  เหลืออยู่จริง ตามที่ต้องบันทึกไว้)**: `pagehide`/`visibilitychange` +
  `keepalive` ปิด window หลักได้ (reload, เปลี่ยนแท็บ, สลับแอปบนมือถือ) แต่
  **ไม่ปิด 100%** — process ถูก kill ทันที (OS kill, แบตหมดกะทันหัน, browser
  crash) เร็วกว่าที่ event handler จะรันได้เลยยังทำให้ attempt หายได้ทางทฤษฎี
  ไม่มีทางปิด window นี้ได้สนิทจาก client-side ฝั่งเดียว — เพิ่มเข้า "จงใจไม่ทำ"
  ด้านล่าง

### Mutation table รอบ 3

รัน `npx vitest run` เต็มชุดหลังแก้แต่ละจุด แล้ว revert ทุกครั้ง — รวม re-verify
R3 (retry sends altered payload ผ่าน `flushPendingAttempt` ใหม่ ยัง kill ได้
เหมือนเดิมหลัง `handleRetrySave` เปลี่ยนไปเรียก `flushPendingAttempt` แทน
`submitAttempt` ตรง ๆ):

| # | Mutation | ผลลัพธ์ |
|---|---|---|
| R1 | ลบ `clearTimeout(timer)` (ผ่าน `clearPendingTimer`) ออกจาก `flushPendingAttempt` | killed — "Retry cancels a live debounce timer instead of racing it..." |
| R2a | `flushPendingAttempt` ใช้ `Object.values(pendingAttemptsRef.current)[0]` แทน `[position]` | killed — 2 เทสต์ ("flushes BOTH pending short_answer submissions on Finish..." + "...when navigating away...") |
| R2b | `flushAllPendingAttempts` ใช้ `.slice(0,1)` (flush แค่ timer แรก) | killed — 2 เทสต์เดียวกับ R2a |
| R3a | `flushForUnload` เรียก `flushAllPendingAttempts()` โดยไม่ส่ง `{keepalive:true}` | killed — 2 เทสต์ (pagehide + visibilitychange) |
| R3b | ลบ `document.addEventListener("visibilitychange", ...)` | killed — เทสต์ visibilitychange |
| R3c | ลบ `window.addEventListener("pagehide", ...)` | killed — เทสต์ pagehide |
| re-verify R3 (รอบ 2) | `handleRetrySave` (ผ่าน `flushPendingAttempt` ใหม่) ส่ง `{...stored, confidence:"guessed"}` แทน `stored` ตรง ๆ | killed — เทสต์ retry-body เดิมเหมือนกัน |

**สรุป: 6/6 มูเทชันที่รอบ 3 เรียกร้องตายหมด, 1 re-verify จากรอบ 2 ไม่มี
regression จากการ refactor `handleRetrySave`**

### Live verification รอบ 3 (docker compose up -d --build, ไม่ใช้ -v)

Stack เดิม, lesson `domain-driven-design/ubiquitous-language` เดิม:

- **Retry แข่งกับ timer → เขียนแค่ 1 แถว**: Pass → mocked 500 → "Not saved" →
  คลิก "Not yet" (schedule timer ใหม่) → คลิก "Retry" (force) → รอ 1.8 วินาที
  (เกิน debounce window) → `SELECT` เห็นแค่ **1 แถวใหม่** (`guessed,
  incorrect` — ตรงกับค่าที่ Retry ส่งจริง ไม่มีแถวที่สองจาก timer เดิมที่ควร
  ถูกยกเลิกไปแล้ว)
- **Reload 150ms หลังตอบยังเขียนแถวได้จริง**: rate short_answer check →
  `page.reload()` หลัง 150ms → รอ request (`keepalive`) เข้า server → `SELECT`
  เห็น **1 แถวใหม่** (`confident, correct`) — ยืนยันว่า pagehide flush +
  keepalive ทำงานจริงกับการ reload จริงในเบราว์เซอร์ ไม่ใช่แค่ mock
- **สอง check ที่ rate ใน debounce window เดียวกัน flush ครบทั้งคู่ผ่าน
  navigate-away จริง (ปุ่ม Next)**: rate check #1 (Pass) และ check #4
  (Not yet) ของ lesson เดียวกันภายใน 1 วินาที แล้วกด "Model-Driven Design →"
  ทันที (ยังอยู่ใน debounce window ทั้งคู่) → `SELECT` เห็น **2 แถวใหม่**
  คนละ `check_key` กัน (`unsure/incorrect` สำหรับ check #4, `guessed/correct`
  สำหรับ check #1) — ไม่มีแถวไหนซ้ำหรือหายเลย
- **Cleanup**: ลบแถวที่ script นี้สร้างทั้งหมด (`id > 18`, รวม 4 แถวจากสามการ
  ทดสอบข้างต้น: 1 จาก retry-race + 1 จาก reload + 2 จาก two-check-flush)
  ยืนยันด้วย `SELECT` ก่อน (`COUNT=17, MAX(id)=75`) และหลัง (`COUNT=13,
  MAX(id)=18` — กลับสู่สภาพเดิมของ Q-2a's dev data เป๊ะ). `lesson_progress`
  ของ lesson นี้ไม่ถูกแตะเลยในรอบนี้ (ไม่ได้ reset/restore เหมือนรอบ 2 เพราะ
  ไม่มีการทดสอบ Finish รอบนี้) ยืนยันว่ายังเป็น `state='passed'`,
  `first_passed_at='2026-07-29 09:05:49'` เดิมทุกตัวอักษร (`last_read_at`
  ขยับตามการเปิดหน้าจริงเท่านั้น ซึ่งเป็นพฤติกรรมที่ถูกต้อง ไม่ใช่สิ่งที่ต้อง
  รักษาให้เดิม)

### Review focus (รอบ 3)

- ทำไม `flushPendingAttempt` (เดิมมีเงื่อนไข "ต้องมี timer อยู่ก่อน") ถึงต้อง
  เปลี่ยนเป็นเรียก `clearPendingTimer` **แบบไม่มีเงื่อนไข** เพื่อให้ `Retry`
  ใช้ฟังก์ชันเดียวกันได้? ถ้ายังใช้เงื่อนไขเดิมจะเกิดอะไรขึ้นตอน Retry ไม่มี
  timer ค้างอยู่เลย?
- ทำไมเทสต์ที่มี lesson แค่ 1 check ถึงไม่พอสำหรับพิสูจน์ว่า
  `flushAllPendingAttempts` ทำงานถูกสำหรับหลาย check พร้อมกัน ทั้งที่นับจำนวน
  ครั้งที่เรียก `postAttempt` ถูกอยู่แล้ว?
- ทำไม `keepalive: true` ถึงต้องส่งผ่าน `fetch` โดยตรง แทนที่จะใช้
  `navigator.sendBeacon` ซึ่งเป็นเครื่องมือมาตรฐานสำหรับ "ส่งข้อมูลก่อนหน้า
  unload"?

### จงใจไม่ทำในรอบนี้ (เพิ่มจากรอบ 3)

- **Residual window หลัง pagehide/visibilitychange + keepalive**: process ถูก
  kill ทันที (OS kill, แบตหมดกะทันหัน, browser crash) เร็วกว่า event handler
  จะรันได้เลย ยังทำให้ attempt ที่กำลัง debounce อยู่หายได้ในทางทฤษฎี — ไม่มี
  ทางปิดได้สนิทจาก client-side ฝั่งเดียว (ต้องมี mechanism ฝั่ง server เช่น
  ส่ง state บางส่วนไปเก็บไว้ก่อน ซึ่งอยู่นอก scope ของ ticket นี้)
- **ข้อควรระวังที่ต้องอ่านคู่กับ headline "flip-flop เขียนแค่ 1 แถว" ของรอบ 2
  (พบใน code review รอบ 4) — ประโยคนั้นเป็นจริงเฉพาะตอนแท็บยัง visible อยู่
  ตลอด**: ถ้าสลับแท็บ (`visibilitychange` → `hidden`) กลางคันระหว่างที่ค่า
  ยังไม่ settle (ยังอยู่ใน debounce window), flush ที่เกิดจาก
  pagehide/visibilitychange จะส่งค่า**ที่ยังไม่ settle**ออกไปทันที ไม่ใช่รอ
  ให้ user ตัดสินใจจบก่อน — ถ้ามีการแก้ไข rating อีกครั้งหลังจากนั้น (เช่น
  กลับมาที่แท็บแล้วเปลี่ยนใจ) จะได้แถวที่สองสำหรับ check เดียวกัน (พิสูจน์สด:
  `correct` ตามด้วย `incorrect` บน check เดียวกัน) — **นี่คือ trade-off ที่
  ถูกต้อง ไม่ใช่บั๊ก**: durability (ไม่เสียข้อมูลตอนปิดแท็บ) สำคัญกว่า dedup
  ในกรณีนี้ และทั้งสองแถวเป็น honest state ที่ user เชื่อจริง ณ ขณะนั้น (ค่าที่
  ถูก flush ก่อนสลับแท็บ + ค่าที่แก้ไขทีหลัง) เพียงแต่ headline "เขียนแค่ 1
  แถว" ต้องอ่านว่า "เขียนแค่ 1 แถว **ถ้าแท็บไม่ถูกสลับหรือปิดกลางคัน**" ไม่ใช่
  จริงเสมอไปทุกกรณี
- **`finish()`'s slug-based identity guard ยังไม่ย้ายมาใช้ `loadGenerationRef`**
  — ดูหัวข้อ "finish()" ท้ายไฟล์นี้สำหรับเหตุผลแบบ blast-radius (ไม่ใช่แค่
  "อยู่นอก scope") และ follow-up ใน `docs/roadmap.md`

### รอบ code-reviewer 4 (SHIP พร้อม 2 follow-up บังคับก่อน PR)

Reviewer ยืนยัน round 1-3 ทั้งหมดถูกต้องและพิสูจน์สดแล้วจริง (SHIP) แต่เจอ
survivor เพิ่ม 2 จุดที่เป็น regression class เดียวกับที่ ticket นี้ใช้เวลา
ทั้ง 3 รอบแก้ — โค้ดถูกอยู่แล้วทั้งคู่ ขาดแค่เทสต์ pin:

- **ปักหมุด (blocking, แก้แล้ว) — `delete debounceTimersRef.current[position]`
  ใน `clearPendingTimer` ไม่มีเทสต์คุ้มครองเลย**: ลบบรรทัดนี้ทิ้ง **รอดทั้งชุด
  169 เทสต์** — เป็นบรรทัดเดียวที่กันไม่ให้ลำดับเหตุการณ์จริงของ browser
  (`visibilitychange`→`hidden` ตามด้วย `pagehide` ติดกัน ซึ่งเกิดขึ้นจริงตอน
  ปิดแท็บ/reload) flush ซ้ำสองครั้ง: `visibilitychange` flush เคลียร์ timer
  ด้วย `clearTimeout` แต่ถ้าไม่ `delete` key ออกจาก `debounceTimersRef.current`
  ด้วย `pagehide`'s `flushAllPendingAttempts` ที่ iterate keys ทีหลังจะยังเห็น
  key เดิมอยู่แล้ว flush ซ้ำ — ทุก reload จริงจะเขียนแถวซ้ำเงียบ ๆ โดย suite
  ทั้งชุดยังเขียว. เพิ่มเทสต์ "the real browser sequence (visibilitychange:
  hidden, then pagehide) flushes exactly once per check, not twice" —
  dispatch ทั้งสอง event ติดกันในลำดับจริง แล้ว assert `postAttempt` ถูกเรียก
  พอดี **1 ครั้ง**
- **ปักหมุด (blocking, แก้แล้ว) — `removeEventListener` cleanup ของ
  pagehide/visibilitychange listener ไม่มีเทสต์คุ้มครองเลย**: ลบทั้งคู่ทิ้ง
  **รอดทั้งชุด 169 เทสต์** เช่นกัน — listener ที่ค้างอยู่หลัง unmount จะยัง
  fire flush ให้ lesson ที่ปิดไปแล้ว เป็น bug class เดียวกับ R1 (stale work
  ถูก attribute ผิดที่) แค่มาทาง listener แทนที่จะมาทาง response ที่ค้างอยู่.
  **เทสต์แรกที่ลองเขียน (เช็คว่า `postAttempt` ไม่ถูกเรียกหลัง unmount+
  pagehide) ใช้ไม่ได้จริง** — per-navigation effect's cleanup เอง (จาก R5
  รอบ 2) ก็ flush ตอน unmount เหมือนกัน ทำให้ debounce timer ถูกเคลียร์ไปแล้ว
  ก่อนที่ listener (ถ้ายังไม่ถูกลบ) จะมีอะไรให้ flush ซ้ำ — confound นี้จะบัง
  ไม่ให้เห็น mutation เลยไม่ว่าจะแก้โค้ดถูกหรือผิด. แก้โดยเปลี่ยนวิธีทดสอบ:
  spy ตรงที่ `window.addEventListener`/`removeEventListener` และ
  `document.addEventListener`/`removeEventListener` แล้ว assert ว่า unmount
  เรียก `removeEventListener` ด้วย handler reference **เดียวกัน**กับที่
  `addEventListener` ใช้ตอน mount — ปักหมุดที่ cleanup โดยตรง ไม่ผ่าน
  side-effect ที่ confound ได้
- **near-equivalent (ยืนยันแล้ว, ไม่แก้) — ตัด `document.visibilityState ===
  "hidden"` เช็คออกจาก `handleVisibilityChange`**: มูเทตแล้ว**รอดทั้ง 169
  เทสต์เหมือนกัน** — แต่เหตุผลต่างจากสองข้อบน: เทสต์ทุกตัวที่ dispatch
  `visibilitychange` ตั้ง `document.visibilityState = "hidden"` ไว้ก่อนเสมอ
  (จำลองสถานการณ์จริงที่มันจะเกิด) ทำให้ผลลัพธ์เหมือนกันไม่ว่าจะเช็คเงื่อนไข
  นี้หรือไม่ — ต่างจากสองข้อบนที่โค้ดถูกอยู่แล้วแค่ไม่มีเทสต์ ข้อนี้คือ
  "มูเทตแล้วสังเกตไม่ออกจากมุมที่ทดสอบอยู่จริง" ซึ่งอาจสังเกตออกได้ถ้าเพิ่ม
  เทสต์ที่ dispatch `visibilitychange` ตอน visibilityState เป็น `"visible"`
  (ยืนยันว่า flush ไม่ถูกเรียกตอนกลับมาเปิดแท็บ) — **ไม่เพิ่มเทสต์นั้นในรอบนี้**
  เพราะไม่ใช่ regression class ที่ ticket นี้กังวล (การ flush เกินความจำเป็น
  ตอนกลับมาเปิดแท็บไม่ทำข้อมูลเสียหาย แค่ทำงานถี่กว่าที่จำเป็น) — บันทึกไว้
  ตรง ๆ ว่าเป็น survivor ที่ยอมรับได้ ไม่ใช่ equivalent แบบสมบูรณ์

### Mutation table รอบ 4

| # | Mutation | ผลลัพธ์ |
|---|---|---|
| 1 | ลบ `delete debounceTimersRef.current[position];` ออกจาก `clearPendingTimer` | killed — "the real browser sequence (visibilitychange:hidden, then pagehide) flushes exactly once per check, not twice" |
| 2 | ลบ `removeEventListener` ทั้งสองบรรทัดออกจาก cleanup ของ pagehide/visibilitychange effect | killed — "removes the pagehide/visibilitychange listeners on unmount" |
| 3 | ตัด `document.visibilityState === "hidden"` เช็คออกจาก `handleVisibilityChange` | **survived (ยอมรับ, บันทึกไว้ตรง ๆ)** — ทุกเทสต์ตั้ง visibilityState เป็น "hidden" ก่อน dispatch เสมอ ไม่มีเทสต์ที่ dispatch ตอน "visible" เพื่อพิสูจน์ว่าไม่ flush เกินจำเป็น — ความเสี่ยงต่ำ (ไม่ทำข้อมูลเสียหาย แค่ flush ถี่กว่าที่จำเป็น) จึงไม่ปิดในรอบนี้ |

**สรุป: 2/2 มูเทชันที่รอบ 4 เรียกร้องตายหมด, 1 survivor ที่ยอมรับได้บันทึกไว้
ตรง ๆ (ไม่ใช่ equivalent สมบูรณ์)** — โค้ด production ไม่เปลี่ยนเลยในรอบนี้
(เพิ่มแค่เทสต์) จึงไม่ต้อง rebuild stack ใหม่สำหรับ live verification

Status: implemented, round 4 fixes applied post code-review, PR pending

## Q-2c — review_cards/review_logs + SM-2 scheduling (not started)

- Schema เพิ่ม: `review_cards` / `review_logs` (คนละตารางกับ `recall_attempts`
  ที่ Q-2a สร้างไว้แล้ว)
- อ่านข้อมูลจาก `recall_attempts` (Q-2a) เป็น input ของ SM-2 quality score —
  ไม่ใช่ schema ใหม่ที่ไม่เกี่ยวกับของเดิม
- Query pattern `WHERE check_key = ? ORDER BY created_at DESC` ที่ index
  `idx_recall_attempts_check_key_created_at` ใน `006_recall.sql` เตรียมไว้ให้แล้ว
