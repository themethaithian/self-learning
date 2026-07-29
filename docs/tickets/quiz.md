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

## Q-2 — persist recall attempts + SRS scheduling (not started)

Scope moved here from `docs/roadmap.md` so the Q-series has one home —
roadmap keeps only a one-line pointer.

- Schema: `recall_attempts` / `review_cards` / `review_logs`
- Key ด้วย `check_key = SHA256(topic/concept/question)` **ไม่ใช้ FK ไป
  `recall_checks.id`** เพราะ importer ลบแล้ว insert ใหม่ทุกครั้ง (`lessons`
  เท่านั้นที่ id คงที่ผ่าน `LAST_INSERT_ID(id)`)
- **ปลดหนี้ของ UX-5/Q-1**: state ต่อข้อ (rating, selected option, confidence)
  ยังไม่ถูก persist เลย — อยู่ใน memory ของ `RecallCheckCard`/`lesson/page.tsx`
  เท่านั้น จะกลายเป็น data loss ทันทีที่ ticket นี้ขึ้นถ้าไม่ออกแบบให้ครอบคลุม
  ทั้งสามค่า ไม่ใช่แค่ pass/fail เดิม
- จุดนี้คือที่ที่ `graded_by='self'` จะเกิดขึ้นจริงครั้งแรก
- **เมื่อ Q-2 ออกแบบ endpoint จริง**: ถ้า `expected_answer` ยังจำเป็นต้องอยู่ที่
  client ต่อ (ตามการตัดสินใจของ Q-1) หรือย้ายไป server-side ทั้งหมด ให้ตัดสินใจ
  ตอนนั้นจากรูปร่าง request/response จริงของ endpoint นี้ ไม่ใช่เดาไว้ล่วงหน้า
