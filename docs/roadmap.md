# Roadmap — Self-Improve Web (v1, 6 สัปดาห์)

> **ภาพใหญ่ทั้งหมดอยู่ไฟล์นี้ไฟล์เดียว** — รายละเอียด ticket แตกเป็นชิ้นเล็กอยู่ใน
> [`docs/tickets/week-1.md`](tickets/week-1.md) … [`week-6.md`](tickets/week-6.md)
> ส่วน design ฉบับเต็ม (DDD, schema, API, VPS) อยู่ที่ [`docs/design.md`](design.md)

## เป้าหมาย v1

เว็บฝึก 4 ทักษะเพื่อเปลี่ยนงาน: System Design (DDD + Distributed Systems),
AWS SAA-C03, Deep Go, DSA — ใช้เวลา ~10 ชม./สัปดาห์ ship ใน 4–6 สัปดาห์
โดยตัวแอปเองคือแบบฝึกหัด (เรียนไปสร้างไป) — เป็น "หนังสือเรียนส่วนตัว" ใช้เองก่อน
Phase 2 ค่อย publish โชว์เป็นผลงาน (ดู design.md §8)

## สถาปัตยกรรมสรุป (5 บรรทัด)

- Go 1.26 stdlib `net/http` เท่านั้น + MySQL 8 (Docker) + Next.js static export
- DDD modular monolith 6 bounded contexts: curriculum / learning / practice / review / assessment / stats
- Deploy: DigitalOcean VPS Singapore (~$6/เดือน) — Docker Compose: mysql + api + caddy (auto-HTTPS)
- LLM runtime: `claude-haiku-4-5` ผ่าน client ที่เขียนเอง (ตรวจ recall, รีวิวโค้ด DSA)
- เนื้อหา: generate offline (writer → verifier) มี mermaid diagram + references → import เข้า MySQL

## แผนรายสัปดาห์

> เลข ticket อ้างจาก plan เดิม (T28–T30, T33 ถูกเลื่อนขึ้นมา) — เลขไม่เรียงแต่คงที่
> เพื่อให้อ้างอิงตรงกันทุกไฟล์

| สัปดาห์ | เป้าหมาย | Definition of Done | Tickets |
|---|---|---|---|
| 1 | Walking skeleton (local) | `docker compose up` + เว็บโชว์ curriculum tree 5 track, repo ขึ้น GitHub + PR template | [T1–T8](tickets/week-1.md) |
| 2 | **ขึ้น VPS** + แกนของระบบอ่าน | แอปรันบน HTTPS จริง, merge → auto-deploy, importer + LLM client พร้อม | [T28–T30, T9–T11](tickets/week-2.md) |
| 3 | Read & Recall จบ + กันข้อมูลหาย | อ่าน lesson จริง (มือถือได้), recall → Haiku ตรวจ → ปลดล็อก, backup ทำงาน | [T12–T15, T33, C1](tickets/week-3.md) |
| 4 | SRS + กราฟแรก | Drill จาก SM-2 queue, dashboard มี streak + trend | [T16–T21, C2](tickets/week-4.md) |
| 5 | DSA ครบวงจร | โจทย์ → approach + โค้ด → LLM รีวิว → สถิติ pattern | [T22–T25, C3–C4](tickets/week-5.md) |
| 6 | Pre-test + stats ครบ + buffer | radar baseline, weekly retro, ปิดงานค้าง | [T26–T27, T31–T32, T34, C5–C6](tickets/week-6.md) |

## Ticket index (ติ๊กเมื่อ merge แล้ว)

- สัปดาห์ 1: [x] T1 [x] T2 [x] T3 [x] T4 [x] T35 [x] T5 [x] T6 [x] T7 [x] T7b [x] T8
- สัปดาห์ 2 (reading-first): [x] T9 [x] T10 [~] C-content (batch 1: DDD บท 1 = 4 lessons) [x] T-read → **reading slice ครบวง: คลิก concept → อ่านบทเรียนจริง (mermaid + references + recall self-grade)** → แล้วค่อย [ ] T28 [ ] T29 [ ] T30 (T11 เลื่อน). content เดินต่อ batch ละ 1 บท
- สัปดาห์ 3: [ ] T12 [ ] T13 [ ] T14 [ ] T15 [ ] T33 [ ] C1
- สัปดาห์ 4: [ ] T16 [ ] T17 [ ] T18 [ ] T19 [ ] T20 [ ] T21 [ ] C2
- สัปดาห์ 5: [ ] T22 [ ] T23 [ ] T24 [ ] T25 [ ] C3 [ ] C4
- สัปดาห์ 6: [ ] T26 [ ] T27 [ ] T31 [ ] T32 [ ] T34 [ ] C5 [ ] C6
- DDIA track (additive, [รายละเอียด](tickets/ddia.md)): [x] T-ddia-track [x] T-ddia-curriculum · **lessons 12/12 บท ครบทั้ง track (61 บทเรียน)** — reader: [Study Reader artifact](https://claude.ai/code/artifact/b67a1cad-4b9e-4b2e-9cf1-755128e1ebb5)

## กติกาการทำงาน (ทุก ticket)

1. ticket ละ ≤ 45 นาที (1 Build session) — ใหญ่กว่านั้นต้องแตกก่อนเริ่ม
2. ลำดับเสมอ: **go-implementer** ทำบน branch `ticket/<id>-<slug>` →
   `go vet` + `go test ./...` ผ่าน → **code-reviewer** รีวิว → เปิด **PR เข้า develop**
   (body มี walkthrough ไทย + Review focus + ผล test + verdict ของ reviewer) →
   **คุณรีวิว + merge** (จาก GitHub mobile ได้ — diff เล็กพออ่านบนจอมือถือ) →
   merge = **auto-deploy ขึ้น VPS** → ติ๊กใน index ข้างบน
3. ห้ามเริ่ม ticket ถัดไปก่อน PR ปัจจุบันถูก merge
4. **Review tiers**: ทุก PR ระบุ Review level — 🟢 skim (scaffolding/config/docs;
   CI เขียว = merge จากสรุปได้เลย) / 🟡 normal / 🔴 careful (domain logic, auth,
   migration, SQL, LLM spend — อ่าน diff จริง) พร้อมเหตุผล 1 บรรทัด
5. **Review-as-quiz**: Review focus ใน PR เป็นคำถาม 2-4 ข้อให้ตอบระหว่างอ่าน
   (เฉลยพับไว้ท้าย PR) — และ design/คำอธิบายยาว ๆ ในแชทจะจบด้วย quiz เสมอ
6. **Model policy**: main session = Opus (orchestrator), workers = opus/sonnet/haiku,
   **Fable = escalation เท่านั้น** (worker พลาด 2 ครั้ง → ขออนุญาต → Fable subagent)
4. Content batch (C1–C6): lesson-writer → lesson-verifier → FAIL เกิน 2 รอบ = พัก concept แล้วรายงาน
5. Frontend ticket ทุกตัวใช้ skill `frontend-design`, ทุก agent ใช้ skill `token-efficiency`

## สิ่งที่คุณต้องเตรียมเอง

- [ ] สร้าง GitHub repo + push (ระหว่างสัปดาห์ 1 — ต้องมีก่อน T30)
- [ ] ติดตั้ง GitHub mobile app + เปิด notification ของ repo (ไว้รีวิว PR จากมือถือ)
- [ ] ต่อ Claude GitHub App + ลงแอป Claude บนมือถือ (Code tab) — ไว้สั่งงาน
  cloud session จากมือถือเวลาไม่อยู่หน้าคอม (vibe coding ต่างจังหวัดได้)
- [ ] บัญชี DigitalOcean + domain ~$10/ปี หรือ DuckDNS ฟรี (ก่อนสัปดาห์ 2)
- [ ] Anthropic API key (ก่อน T11 สัปดาห์ 2)
- [ ] บัญชี Cloudflare R2 หรือ Backblaze B2 สำหรับ backup (ก่อน T33 สัปดาห์ 3)
