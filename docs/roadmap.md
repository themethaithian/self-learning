# Roadmap — Self-Improve Web

> **ภาพใหญ่ทั้งหมดอยู่ไฟล์นี้ไฟล์เดียว** — รายละเอียด ticket อยู่ใน [`docs/tickets/`](tickets/)
> ส่วน design ฉบับเต็ม (DDD, schema, API, VPS) อยู่ที่ [`docs/design.md`](design.md)
>
> อัปเดตล่าสุด **2026-07-29** (หลัง PR #46)

## เป้าหมายตอนนี้

เว็บ "หนังสือเรียนส่วนตัว" ที่**เรียนแล้ววัดผลได้** เพื่อเตรียมสัมภาษณ์ตำแหน่ง
Backend / AI-CRM Platform — โหมด **learn-first**: ให้ความสำคัญกับ *เนื้อหา + วงจร
อ่าน→ทดสอบ→วัดผล* ก่อน feature และก่อน deploy

ทักษะเป้าหมาย 7 track: DDD · **DDIA (Distributed Data)** · **AI & LLM Systems** ·
Distributed Systems · AWS SAA-C03 · Go · DSA
(สองตัวหนา = เขียนบทเรียนครบแล้ว, ที่เหลือมี tree แต่ยังไม่มีบทเรียน)

## สถาปัตยกรรมสรุป

- Go 1.26 stdlib `net/http` เท่านั้น + MySQL 8 (Docker) + Next.js static export
- DDD modular monolith — bounded context ที่ **มีจริงแล้ว**: `curriculum` · `learning` · `prefs`
  (design.md ยังพูดถึง practice / review / assessment / stats ซึ่งยังไม่ได้สร้าง)
- เนื้อหา: generate offline (lesson-writer → lesson-verifier) มี mermaid + references → import เข้า MySQL
- Grading: **self-graded** (`graded_by='self'`) — ยังไม่มี LLM runtime, ไม่มี API key, ไม่มีค่าใช้จ่าย
- รันเอง: `make dev` (mysql + api + web + seed) — **ยังไม่ deploy ขึ้น VPS**

## แผนตาม phase

> **หมายเหตุสำคัญ**: แผน "6 สัปดาห์ T1–T35" ฉบับเดิม**เลิกใช้แล้ว** ตั้งแต่เป้าหมายแคบลง
> เป็นตำแหน่ง Backend/AI-CRM — เลข T เดิมยังอ้างอิงได้ใน `tickets/week-*.md` แต่ลำดับ
> การทำงานจริงคือ phase ข้างล่างนี้

| Phase | เป้าหมาย | สถานะ |
|---|---|---|
| 0. Walking skeleton + reading slice | คลิก concept → อ่านบทเรียนจริง (mermaid + references + recall) | ✅ จบ (PR #1–#17) |
| 1. Content | เขียนบทเรียน 2 track ที่ตรงกับตำแหน่งที่สุด | ✅ จบ (PR #18–#37) |
| 2. Quality + รันเองได้ | MCQ ที่เดาไม่ได้ + `make dev` คำสั่งเดียว | ✅ จบ (PR #38–#43) |
| 3. Guided learning path | รู้ว่า "วันนี้อ่านอะไรต่อ" และอ่านถึงไหนแล้ว | 🔄 กำลังทำ (UX-1…UX-7) |
| 4. Measurable recall | quiz กดเลือกได้ + เก็บประวัติ + SRS ที่วัดผลได้ | ⏭ ถัดไป (Q-1…Q-3) |
| 5. Visual simulation | บทเรียนที่เห็นภาพและโต้ตอบได้ | ⏭ หลัง phase 4 (SIM-0…SIM-3) |
| — | Deploy ขึ้น VPS + track ที่เหลือ | ⏸ **พักไว้โดยตั้งใจ** |

## Ticket index (ติ๊กเมื่อ merge แล้ว)

### Phase 0 — walking skeleton + reading slice ✅
- [x] T1 [x] T2 [x] T3 [x] T4 [x] T35 [x] T5 [x] T6 [x] T7 [x] T7b [x] T8 ([week-1](tickets/week-1.md))
- [x] T9 [x] T10 [x] T-read — คลิก concept → อ่านบทเรียนจริง ([week-2](tickets/week-2.md))
- [x] C-content — DDD บท 1 (4 บทเรียน) + trade-off sections
- [x] T-reading-comfort — mermaid อ่านออกบนมือถือ + พื้นหลัง cream

### Phase 1 — content ✅
- **DDIA** ([รายละเอียด](tickets/ddia.md)): [x] T-ddia-track [x] T-ddia-curriculum ·
  **61 บทเรียน / 12 บท ครบทั้ง track** — [Study Reader artifact](https://claude.ai/code/artifact/b67a1cad-4b9e-4b2e-9cf1-755128e1ebb5)
- **AI & LLM Systems** ([รายละเอียด](tickets/ai-systems.md)): [x] T-ai-track (tree 11 บท / 53 concept) ·
  **53 บทเรียน / 11 บท ครบทั้ง track** — เจาะตำแหน่ง Backend/AI-CRM ([🎯 roadmap](https://claude.ai/code/artifact/a6d332fd-de7b-44bb-91db-37549984c7a7))
- **DDD**: 4 บทเรียน (บท 1 เท่านั้น) — ยังไม่ครบ track
- `distsys` / `aws` / `go` / `dsa`: มี curriculum tree แต่ **0 บทเรียน**

### Phase 2 — quality + รันเองได้ ✅
- [x] T-local-docker ([รายละเอียด](tickets/local-docker.md)) — `make dev` = mysql + api + web + seed
- [x] C-mcq-sweep, [x] C-mcq-balance ([รายละเอียด](tickets/mcq-quality.md)) —
  MCQ ทั้ง 356 ข้อ เดาด้วย heuristic ได้ **33.7%** (เดิม 89.6%) + กติกาการแก้ distractor

### Phase 3 — guided learning path 🔄 ([รายละเอียด](tickets/ux-today.md))
- [x] UX-1 (#40) has_lesson/est_minutes บน curriculum API
- [x] UX-1b (#41) focus track preference API + bounded context `prefs`
- [x] UX-2 (#42) หน้า `/learn` — track card + focus picker
- [x] UX-3 (#44) track view — concept ที่ยังไม่มีบทเรียนกดไม่ได้ + `n/N ready`
- [x] UX-4 (#45) 🔴 progress API (`GET`/`PUT /api/v1/progress`) + bounded context `learning`
- [x] UX-5 (#46) reader loop — breadcrumb + Finish + Next (+ vitest ตัวแรกของ `web/`)
- [ ] UX-6 — หน้า Today + IA switch: "วันนี้อ่านอะไรต่อ" ไล่จาก focus track ก่อน
- [ ] UX-7 — หน้า Progress + ตัวบอกสถานะทั้งแอป

### Phase 4 — measurable recall ⏭
- [ ] Q-1 — quiz **recall-first 3 stage**: ตอบในใจ → เลือกความมั่นใจ → เฉลย · shuffle ตัวเลือกตอน render
- [ ] Q-2 — schema `recall_attempts` / `review_cards` / `review_logs`
  · key ด้วย `check_key = SHA256(topic/concept/question)` **ไม่ใช้ FK ไป `recall_checks.id`**
  เพราะ importer ลบแล้ว insert ใหม่ทุกครั้ง (id ของ `lessons` เท่านั้นที่คงที่)
  · **ปลด debt ของ UX-5**: rating รายข้อยังไม่ถูก persist — จะกลายเป็น data loss ทันทีที่ ticket นี้ขึ้น
- [ ] Q-3 — เพิ่ม `explanation` ให้ MCQ ครบทั้ง 356 ข้อ (รอบเดียวกับที่ขัดเกลา distractor ที่หลุดธีม 5 จุด)

### Phase 5 — visual simulation ⏭
- [ ] SIM-0 — predict-before-reveal (0 KB, ไม่ใช้ library) · งานวิจัยชี้ว่า animation แบบดูเฉย ๆ
  **ทำให้คนที่พื้นฐานน้อยเรียนแย่ลง** (Kaushal & Panda 2019, d = −0.16) การทำนายก่อนเฉลยถึงจะได้ผล
- [ ] SIM-1 latency-percentiles (slider) · [ ] SIM-2 quorum W+R>N (slider) ·
  [ ] SIM-3 2PC (step + ปุ่ม crash coordinator) → **หยุดประเมินก่อนทำเพิ่ม**

### พักไว้โดยตั้งใจ (ไม่ได้ยกเลิก)
- **Deploy**: [x] T28 (prod Dockerfile + compose + Caddyfile, #25) · [ ] T29 [ ] T30 (GitHub Actions → VPS)
  — พักตั้งแต่เปลี่ยนเป็นโหมด learn-first
- **Feature จากแผนเดิม**: SRS drill + dashboard ([week-4](tickets/week-4.md)) · DSA loop ([week-5](tickets/week-5.md)) ·
  pre-test + radar + stats ([week-6](tickets/week-6.md)) · T11 (LLM client), T33 (backup), T12–T15
  — บางส่วนจะถูกแทนที่ด้วย phase 4 ที่ออกแบบใหม่แล้ว
- **Track ที่เหลือ**: DDD บท 2–5, distsys, aws, go, dsa · และ track ที่ยังไม่มี:
  event-driven/Kafka, gRPC, observability, Kubernetes

## กติกาการทำงาน (ทุก ticket)

1. ticket ละ ≤ 45 นาที (1 Build session) — ใหญ่กว่านั้นต้องแตกก่อนเริ่ม
2. ลำดับเสมอ: **go-implementer** ทำบน branch `ticket/<id>-<slug>` →
   `go vet` + `go test ./...` (+ `npm test` ถ้าแตะ `web/`) ผ่าน → **code-reviewer** รีวิว →
   แก้จนไม่มี required → เปิด **PR เข้า develop** (body มี walkthrough ไทย + Review focus +
   ผล test + verdict ของ reviewer) → **คุณรีวิว + merge** (จาก GitHub mobile ได้)
3. ห้ามเริ่ม ticket ถัดไปก่อน PR ปัจจุบันถูก merge
4. **Review tiers**: ทุก PR ระบุ Review level — 🟢 skim (scaffolding/config/docs;
   CI เขียว = merge จากสรุปได้เลย) / 🟡 normal / 🔴 careful (domain logic, auth,
   migration, SQL, LLM spend — อ่าน diff จริง) พร้อมเหตุผล 1 บรรทัด
5. **Review-as-quiz**: Review focus ใน PR เป็นคำถาม 2-4 ข้อให้ตอบระหว่างอ่าน
   (เฉลยพับไว้ท้าย PR) — และคำอธิบายยาว ๆ ในแชทจะจบด้วย quiz เสมอ
6. **Model policy**: main session = Opus (orchestrator ล้วน ไม่เขียนเอง),
   workers = opus/sonnet/haiku, **Fable = escalation เท่านั้น**
7. **Content batch**: lesson-writer → lesson-verifier → FAIL เกิน 2 รอบ = พัก concept แล้วรายงาน
8. Frontend ticket ทุกตัวใช้ skill `frontend-design`, ทุก agent ใช้ skill `token-efficiency`
9. **Mutation-test ตัว test เอง**: repo นี้เคยมี test เขียวทั้งที่โค้ดพัง 4 ครั้ง
   (CORS test เทียบกับ constant ตัวเอง · stub SQL driver ที่ dispatch ด้วย query identity
   ทำให้ `AND`→`OR` รอด · fixture frontend ที่มี track เดียวทำให้ "ห้ามข้าม track" ไม่ถูกทดสอบ ·
   วัด "กฎที่ตั้งไว้ถูกละเมิดไหม" แทน "เดาถูกกี่ %") → **ต้องพังของจริงแล้วดูว่า test ไหนจับได้**
10. **คำกล่าวอ้างเรื่อง DOM / a11y / contrast ต้องวัดจาก browser จริง** (Playwright + headless
    Chromium ใช้ได้ในเครื่องนี้) ไม่ใช่การอ่านโค้ด — และ rebuild container ก่อนวัดสี

## สิ่งที่คุณต้องเตรียมเอง

- [x] GitHub repo + PR template + CI
- [x] GitHub mobile app (ไว้รีวิว PR จากมือถือ)
- [ ] บัญชี DigitalOcean + domain — **ยังไม่ต้อง** จนกว่าจะกลับมาทำ T29/T30
- [ ] Anthropic API key — **ยังไม่ต้อง** (v1 เป็น self-graded, ดู design.md §LLM)
- [ ] Cloudflare R2 / Backblaze B2 สำหรับ backup — คู่กับ T33 ซึ่งพักไว้พร้อม deploy
