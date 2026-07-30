# Handoff — ย้ายเครื่อง Windows → Mac (3 วัน)

เขียนไฟล์นี้เพราะ memory files ของ Claude Code (`~/.claude/projects/…`) อยู่บนเครื่อง
Windows เท่านั้น ไม่ตามไป Mac ด้วย — session ใหม่บน Mac จะไม่มีความจำอะไรเลยนอกจาก
สิ่งที่อยู่ใน git repo นี้ อ่านไฟล์นี้ไฟล์เดียวควร resume งานต่อได้ภายใน 10 นาที

## 1. เริ่มตรงนี้

**Priority เปลี่ยนแล้ว**: เป้าหมายอันดับหนึ่งตอนนี้คือสอบ **AWS Certified Solutions
Architect – Associate (SAA-C03)** ภายใน **1-2 เดือน**ข้างหน้า อ่านวันละ **1-2
ชั่วโมง** โดยที่ผู้ใช้ยัง**ไม่เคยแตะ AWS จริง** (รู้แค่ชื่อ service หลักๆ)

งานอื่นทั้งหมด — Q-2d, Q-3, SIM phase, deploy — **พักไว้ก่อน** แผนงาน AWS
cert อยู่ที่ `docs/tickets/aws-cert.md` (กำลังเขียนอยู่บน branch
`ticket/aws-0-saa-tree` — ดูสถานะจริงในหัวข้อ 4 ว่ารอดจาก handoff มาหรือเปล่า)

## 2. รันบน macOS

Prereqs: Docker Desktop, Go **1.26.5** (ตาม `go.mod`), Node + npm, `make`, `gh`
CLI ที่ login แล้ว

**`make dev` ไม่ต้องมีไฟล์ `.env` เลย** — ตรวจสอบแล้วว่า `docker-compose.yml`
ใส่ default ของทุก env var ไว้ในไฟล์เอง (`${VAR:-default}`): `API_BEARER_TOKEN`
default เป็น `local-dev-token`, `ANTHROPIC_API_KEY` เป็น placeholder (ไม่มีจุดไหน
ในโค้ดเรียก Anthropic จริง ไม่ต้องใช้ key จริงเลย), DB creds เป็น `app`/`changeme`
`.env*` ถูก gitignore ยกเว้น `.env.example` — และ `.env.example` เป็น template
สำหรับ **prod** เท่านั้น อย่า copy มาเป็น `.env` แล้วคาดหวัง dev defaults
(ค่าใน `.env.example` เป็นคนละชุดกับ default ใน compose)

Make targets: `up` (mysql อย่างเดียว) · `dev` (`docker compose up -d --build`) ·
`dev-down` · `seed` · `logs` · `test` · `vet` · `run`

**⚠️ `make dev-reset` = `docker compose down -v` — ลบ MySQL volume ทิ้งทั้งหมด**
เป็น target เดียวที่ห้ามรันเล่นๆ เด็ดขาด เพราะมันลบ reading progress กับ recall
attempts ที่เป็นข้อมูลจริงของผู้ใช้ ไม่ใช่ fixture ที่สร้างใหม่ได้ฟรี

Ports: web `127.0.0.1:3000`, MySQL `127.0.0.1:3306`, API `8080`

เปิดแอปที่ **`http://localhost:3000` ไม่ใช่ `http://127.0.0.1:3000`** — สอง
origin นี้ต่างกัน token เก็บแยกตาม origin ถ้าเปิดผิด origin จะดูเหมือน
"ยังไม่ login" ทั้งที่ login ไปแล้ว

รันครั้งแรกแอปจะถามหา token — ใส่ **`local-dev-token`**

Frontend tests: `cd web && npm test` (vitest, jsdom + React Testing Library)
Backend: `go vet ./... && go test ./...`

## 3. อะไรไม่เดินทางไปด้วย

- **Memory files** (`~/.claude/projects/C--Users-60010375-Desktop-self-learning/memory/`)
  อยู่บนเครื่อง Windows เท่านั้น อะไรที่ session บน Mac ต้องรู้ ต้องอยู่ใน repo
- **Docker volume** — reading progress, recall attempts, review cards ทุกอย่าง
  เริ่มจากศูนย์บน Mac curriculum กับ lessons re-seed อัตโนมัติผ่าน `make dev`
  อยู่แล้ว เป็นเรื่องคาดหวัง ไม่ใช่ bug
- `.env` — ไม่ได้ commit ไว้ใน git และไม่จำเป็นสำหรับ dev

## 4. สถานะ ณ ตอนส่งมอบ

Merged บน `develop` แล้ว: guided-learning epic UX-1…UX-7 (PR #40–#50), ต่อด้วย
**Q-1** recall-first 3-stage quiz (#51), **Q-2a** `recall_attempts` +
`POST /api/v1/progress/{topic}/{concept}/attempts` (#52), **Q-2b** frontend
wiring พร้อม debounced submit, unload flush, visit-based identity guard (#53)
จำนวน test บน `develop`: **Go 275, web 169**

**Branch ที่กำลังทำค้างอยู่ — เช็คสถานะจริงด้วย
`git log --oneline develop..<branch>` และ `gh pr list` ก่อนทำอะไรต่อ:**

- **`docs/aws-currency-checklist`** (push แล้ว, 1 commit) — `docs/tickets/aws-currency-checklist.md`
  รายการ AWS ที่เปลี่ยนไปจนคำตอบเก่ากลายเป็นผิด ทุก item มี URL ทางการ
  **นี่คือไฟล์ที่มีค่าที่สุดของ track นี้** ต้องเป็น context บังคับของทั้ง
  `lesson-writer` และ Fable ทุก batch · merge ก่อนเริ่มเขียนข้อสอบ
- **`docs/session-handoff`** (push แล้ว, 1 commit) — ไฟล์ที่คุณกำลังอ่านอยู่นี้
- **`ticket/q-2c-sm2-domain`** (push แล้ว, **2 commits**: `e08abb4` แล้ว
  `4d4c416` "Q-2c round 2: DATETIME due dates, exact-integer interval, honest
  deviations") — รอบแก้ลงแล้วแต่ **ยังไม่ได้ผ่าน re-review** ประเด็นที่รอบแรก
  เจอและรอบสองอ้างว่าแก้แล้ว: `due_at TIMESTAMP` เก็บวันที่ไกล ๆ ของ SM-2 ไม่ได้
  (MySQL สูงสุด 2038-01-19 แต่ q=5 ติดกัน 8 ครั้งได้ due date ราวปี 2042 →
  `ERROR 1292`) · float path ทำ interval เพี้ยนไปหนึ่งวัน (`95 × 2.30 = 218.5`
  ได้ 218 แทน 219, กวาดทั้งช่วงเจอ 77 คู่) · เบี่ยงจาก published SM-2 สองจุด
  ที่เอกสารเดิมอ้างผิดว่าเป็น "การแก้ความกำกวม" ทั้งที่ step 6 เขียนชัด
  **ต้อง re-review ก่อน merge** — รายละเอียดครบใน `docs/tickets/quiz.md`
  · หมายเหตุ: branch นี้มี `git stash` ค้างอยู่หนึ่งอัน ที่ระบุว่าเป็น artifact
  ของ autocrlf ไม่ใช่งานจริง — `git stash list` แล้วทิ้งได้ถ้าไม่มี diff จริง
- **`ticket/aws-0-saa-tree`** — rebalance `content/curriculum/aws.json` จาก 37
  เป็นราว 50 concept + เขียน `docs/tickets/aws-cert.md`
  **ตอน handoff branch นี้มี 0 commit** และถูกสั่งให้ commit แบบ WIP ทันที
  → เช็ก `git log --oneline develop..ticket/aws-0-saa-tree` ก่อน ถ้ายังว่าง
  แปลว่างานหายไปกับ session ต้องเริ่มใหม่ ซึ่งไม่แพงเพราะ spec ทั้งหมด
  (blueprint + gap analysis) อยู่ในหัวข้อ 5 ของไฟล์นี้แล้ว
- ถ้า branch ไหนไม่มี commit เลย (เท่ากับ `develop` เป๊ะ) แปลว่างานหายไปพร้อม
  session ต้องเริ่มใหม่จาก ticket doc ที่เกี่ยวข้อง

## 5. แผน AWS — blueprint ที่ยืนยันแล้ว

ข้อมูลด้านล่างตรวจสอบกับแหล่งทางการของ AWS แล้ว **ห้าม derive ใหม่จากความจำ
ของโมเดล** เพราะเป็นจุดที่โมเดลมักตอบผิด:

- **SAA-C03 คือ exam version ปัจจุบัน ไม่มี SAA-C04** มีเว็บ third-party
  หลายที่อ้างว่ามี SAA-C04 (บางที่บอกว่า "released March 2024") — เท็จ
  exam-guide index ของ AWS เองมีแค่ SAA-C03 และ index นั้นยัง maintain อยู่จริง
  (มีวันที่ของ SOA-C03 เดือนกันยายน 2025 อยู่ในหน้าเดียวกัน) เจอ resource ไหน
  พูดถึง "SAA-C04" ให้เพิกเฉย
- 65 ข้อ (**50 ข้อนับคะแนน + 15 ข้อไม่นับคะแนนแต่ไม่บอกว่าข้อไหน**), 130 นาที,
  ผ่านที่ **720/1000** (scaled score)
- **Compensatory scoring** — ผ่านที่คะแนนรวมเท่านั้น ไม่ต้องผ่านทีละ domain
- **ตอบผิดไม่โดนหักคะแนน** และข้อที่เว้นว่างนับเป็นผิด → ห้ามเว้นข้อไหนว่างไว้
  เด็ดขาด
- รูปแบบข้อสอบ: multiple choice (เลือก 1 จาก 4) และ multiple response
  (เลือก 2+ จาก 5+ ตัวเลือก) โดยโจทย์จะบอกจำนวนที่ต้องเลือกตรงๆ เช่น
  "(Select TWO.)"
- Domains: **1 Secure 30% · 2 Resilient 26% · 3 High-Performing 24% ·
  4 Cost-Optimized 20%** รวม 14 task statement
- ชื่อ service ที่ถูกต้องอิงจาก **docs.aws.amazon.com** เท่านั้น ไม่ใช่ PDF จาก
  `d1.awsstatic.com` — PDF (v1.1) เก่ากว่าเว็บอย่างน้อย 2 จุดที่เจอแล้ว

## 6. Pipeline การผลิตข้อสอบ (ส่วนที่เสี่ยงหายมากที่สุด)

เป้าหมาย ≈ **550 ข้อ** แบ่งตามน้ำหนัก domain (≈165 / 143 / 132 / 110 ข้อ)
ผลิตทีละ domain มี 4 gate:

1. `lesson-writer` (opus) ร่างโจทย์ + answer key + explanation โดย**ทุก
   fact ต้องมี URL ของ AWS official docs กำกับ**
2. Grounding pass — **fetch URL เหล่านั้นจริง** เพื่อยืนยันว่า claim นั้นอยู่
   ในหน้าเว็บจริงๆ
3. `lesson-verifier` (sonnet) เช็คตาม `docs/tickets/mcq-quality.md` —
   guessability จาก length/position ต้องใกล้ 33%
4. **Fable ตัดสิน answer key เท่านั้น** แบบ adversarial: *"พิสูจน์ว่าคำตอบที่
   มาร์กไว้ผิด หรือมีตัวเลือกอื่นที่ถูกเท่ากันหรือถูกกว่า"* ใช้ Fable แค่ gate
   นี้ gate เดียว เพราะจุดที่ผิดพลาดแล้วแพงที่สุดคือ answer key ผิด และใน AWS
   ส่วนใหญ่มักเป็นเคส "สองตัวเลือก defensible ทั้งคู่ แต่เลือกผิดตัว" ซึ่ง
   ต้องใช้โมเดลที่แข็งที่สุดถึงจะจับได้ — ข้อยกเว้นนี้ผู้ใช้ authorize เอง
   นอกเหนือจากกฎปกติของ CLAUDE.md ที่สงวน Fable ไว้ให้ escalation อย่างเดียว

**กฎ canary — เขียนเด่นๆ ไว้ตรงนี้เลย** ทุก audit batch จะแอบใส่ข้อที่ทำผิด
ตั้งใจไว้โดยไม่บอกผู้ตรวจ: answer key ผิด, ตัวเลขโควตาผิด, service ที่ deprecated
แล้วแต่นำเสนอเหมือนยังใช้อยู่, และข้อที่มีสองตัวเลือกถูกเท่ากัน **ถ้าผู้ตรวจพลาด
canary ข้อไหนแม้ข้อเดียว ผลตรวจของ batch นั้นถือเป็นโมฆะทันที ไม่ว่าจะรายงาน
ผลว่าอย่างไร** เหตุผลที่ต้องเข้มขนาดนี้: 3 ticket ที่ผ่านมาของโปรเจกต์นี้เคย
ส่งผลตรวจว่า "11/11 mutations killed" (ที่จริงมี 4 ตัวรอด), "19/19 ไม่มี
coverage gap" (ที่จริงมี 2 ตัวรอด), และ debt entry ที่กุขึ้นจาก grep ที่ทำผิด
audit ที่ผ่านทุกอย่างแยกไม่ออกจากการไม่มี audit เลย นอกจากจะพิสูจน์ได้ว่ามันจับ
ผิดจริงได้บ้าง

**รูปแบบ explanation** (requirement อันดับหนึ่งที่ผู้ใช้ระบุไว้ตรงๆ): ทุกข้อ
ต้องอธิบาย (a) ทำไมคำตอบที่ถูกถึงถูก ผูกกับ constraint เฉพาะในโจทย์ (b) ทำไม
**แต่ละ** distractor ถึงผิด ระบุ constraint ในโจทย์ที่มันละเมิด (c) decision
rule ที่เอาไปใช้ซ้ำได้ เหตุผล: ผู้ใช้รู้จักชื่อ service แต่ไม่เคยใช้งานจริง
explanation จึงต้องสอนวิธี**เลือกระหว่าง service ภายใต้ constraint**ไม่ใช่ทบทวน
คำนิยาม — เพราะนั่นคือสิ่งที่ข้อสอบจริงวัด

**Currency checklist** ของการเปลี่ยนแปลง AWS ล่าสุดที่ทำให้คำตอบที่โมเดลจำไว้
ผิดพลาด กำลังค้นคว้าอยู่ตอน handoff ถ้า `docs/tickets/aws-cert.md` ไม่มีหัวข้อนี้
แปลว่าหายไปพร้อม session ต้องทำใหม่: ค้นการเปลี่ยนแปลง/เปลี่ยนชื่อ/deprecate/
default ใหม่ของ AWS ที่ตัดกับ syllabus ของ SAA-C03 โดย source มาจาก
`docs.aws.amazon.com` และ `aws.amazon.com/about-aws/whats-new/` เท่านั้น แต่ละ
item ต้องมี URL กำกับ และแยกหัวข้อสำหรับเคสที่ **คำตอบที่ดีที่สุดทางเทคนิค
วันนี้ ต่างจากคำตอบที่ข้อสอบคาดหวัง** (topic กลุ่มนี้ให้หลีกเลี่ยงในข้อสอบ
แทนที่จะเดา)

## 7. กติกาที่ต้องอยู่รอดข้ามเครื่อง

- Main session เป็น **orchestrator อย่างเดียว** — วางแผน, delegate, review
  ผลลัพธ์ของ subagent ไม่เขียน lesson หรือโค้ดเอง
- **ทุก subagent prompt ต้องเปิดด้วยการบอกว่าตัวเองเป็น sole worker ห้าม
  spawn agent อื่นหรือรับคำสั่งจาก agent อื่น** session ที่ผ่านมามี agent ตัวหนึ่ง
  อ่าน orchestration rules ใน CLAUDE.md แล้วสรุปว่าตัวเองเป็น orchestrator เอง
  เลย dispatch reviewer ของตัวเอง จนมี 3 agent มาแก้ไฟล์เดียวกันพร้อมกันแล้ว
  revert การแก้ mutation test ของกันและกัน ต้องตั้งเป็นกฎตายตัวด้วยเหตุผลนี้
- **ห้าม `docker compose down -v` และห้ามลบ Docker volume เด็ดขาด**
  code-reviewer ตัวหนึ่งเคยรันคำสั่งนี้แล้วทำ reading-progress rows ของผู้ใช้
  บนเครื่อง local หายหมด
- Mutation-test ตัว test เองทุก ticket และ **verify ตัว verifier ด้วย** —
  หลาย claim ที่ "ยืนยันแล้ว" ใน session ที่ผ่านมากลับเป็น self-confirming:
  test ที่ kill mutation ได้เพราะ mutation ดันไปใช้ string เดียวกับที่ test
  ใช้อยู่แล้ว, distribution test ที่ไม่เคยรัน production code path จริงเพราะ
  ทุก caller inject fake RNG
- Claim ที่เห็นผลบน browser จริง (DOM, a11y, contrast) ต้องวัดจาก headless
  Chromium จริงผ่าน Playwright ไม่ใช่อ่านจากโค้ดเฉยๆ
- ทำทีละ ticket, PR merge ก่อนเริ่ม ticket ถัดไป **มีข้อยกเว้นเดียว**: งาน
  AWS content (`content/`) กับ AWS engine (`internal/`, `web/`) ทำขนานกันได้
  เพราะ question bank เป็นคอขวดที่ใช้เวลานานที่สุด และสองสายนี้ไม่แตะไฟล์ร่วมกัน

## 8. สามคำสั่งแรกบน Mac

1. Clone/pull repo แล้ว `make dev` เปิด `http://localhost:3000` วาง token
   `local-dev-token`
2. `git log --oneline develop..ticket/q-2c-sm2-domain` และ
   `git log --oneline develop..ticket/aws-0-saa-tree` บวกกับ `gh pr list`
   เพื่อดูว่างานอะไรรอดมาบ้าง
3. อ่าน `docs/tickets/aws-cert.md`
