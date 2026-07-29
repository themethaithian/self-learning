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

## Q-2b — frontend wiring (not started)

เรียก `POST /api/v1/progress/{topic}/{concept}/attempts` (Q-2a) จาก
`RecallCheckCard`/`lesson/page.tsx` จริง — **นี่คือ ticket ที่ทำให้ confidence/
selected option/outcome ของ Q-1 ที่อยู่ใน memory เฉย ๆ persist จริงในที่สุด
แทนที่จะหายตอน refresh** (หนี้ที่ Q-1 บันทึกไว้ตรง ๆ ว่ายังไม่ทำ)

## Q-2c — review_cards/review_logs + SM-2 scheduling (not started)

- Schema เพิ่ม: `review_cards` / `review_logs` (คนละตารางกับ `recall_attempts`
  ที่ Q-2a สร้างไว้แล้ว)
- อ่านข้อมูลจาก `recall_attempts` (Q-2a) เป็น input ของ SM-2 quality score —
  ไม่ใช่ schema ใหม่ที่ไม่เกี่ยวกับของเดิม
- Query pattern `WHERE check_key = ? ORDER BY created_at DESC` ที่ index
  `idx_recall_attempts_check_key_created_at` ใน `006_recall.sql` เตรียมไว้ให้แล้ว
