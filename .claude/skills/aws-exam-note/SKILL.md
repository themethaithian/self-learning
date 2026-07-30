---
name: aws-exam-note
description: Ingest AWS SAA-C03 practice-exam questions + answer explanations into content/exam-notes/. Use whenever the user pastes a practice exam question with its answer or "Overall explanation" (Tutorials Dojo style), or asks to จดเทคนิคสอบ AWS / บันทึกโจทย์ / เก็บเทคนิคจากเฉลย. Must work identically in any session, old or new. After each summary the user may ask follow-ups — those must be recorded as re-read points.
---

# AWS Exam Note — จดเทคนิคสอบจาก practice exam

Input: โจทย์ + เฉลย ("Overall explanation") หนึ่งข้อหรือหลายข้อที่ user วางมา
Output: ไฟล์โน้ตใต้ `content/exam-notes/aws-saa-c03/` อัปเดตให้อยู่ในรูปแบบเดียวกันเสมอ
จบด้วยสรุปสั้น + quiz · หลังสรุป user ถามต่อได้เสมอ และคำถาม follow-up หรือ quiz
ที่ตอบผิด = จุดที่ยังไม่เข้าใจ **ต้องจดลง `review-again.md` ทันที** (ขั้นตอน 8)

งาน ingest ทำใน main session ได้เลย — เป็นงานจัดเก็บโน้ต ไม่ใช่ lesson/code
(ข้อยกเว้น orchestrator-only ตกลงไว้ใน ticket AWS-EN1: หัวข้อ "AWS-EN" ใน
`docs/tickets/aws-cert.md` และบรรทัด exception ใน CLAUDE.md) โหลด
`token-efficiency` skill ก่อนรอบแรกของ session

## ไฟล์ (ใต้ `content/exam-notes/aws-saa-c03/`)

| ไฟล์ | เก็บอะไร | อ่านตอน ingest |
|---|---|---|
| `index.md` | สารบัญ + สถิติ + ลิงก์หมวด | ทุกรอบ |
| `keyword-map.md` | ตาราง signal → answer รายหมวด | **เต็มไฟล์ ทุกรอบ** |
| `traps.md` | trap patterns + ตัวนับ | **เต็มไฟล์ ทุกรอบ** |
| `glossary.md` | ศัพท์ A→Z คำอธิบายไทย 1 บรรทัด | **เต็มไฟล์ ทุกรอบ** |
| `patterns/<category>.md` | entry เต็มรายหมวด | Grep เฉพาะไฟล์หมวดที่เกี่ยว |
| `review-again.md` | จุดที่ user เคยไม่เข้าใจ ไว้อ่านซ้ำ | เมื่อเข้าขั้นตอน 8 |

ทำไม 3 ไฟล์กลางต้องอ่านเต็มไฟล์: trap/ศัพท์/keyword ถูก generalize ให้ไม่ผูกกับชื่อ
service โดยเจตนา — grep ด้วยชื่อ service จะหาของเดิมไม่เจอ แล้วจะสร้างซ้ำจนตัวนับเพี้ยน
ไฟล์เหล่านี้เล็กและ bounded อ่านเต็มไม่ขัด token-efficiency · grep จำกัดไว้กับ
`patterns/*.md` เท่านั้น

หมวด canonical (สร้างไฟล์เมื่อมี entry แรก แล้ว**เพิ่มลิงก์ใน index.md ทันที**):
compute · storage · database · networking · security-iam · analytics-streaming ·
integration-messaging · containers-serverless · migration-transfer ·
management-governance · ml-ai · cost-optimization
โจทย์คาบหลายหมวด → ยึดหมวดของ service ที่**เป็นตัวตัดสิน**คำตอบ ห้ามลง entry ซ้ำสองไฟล์

## ขั้นตอนต่อ 1 รอบ ingest

1. แยก input เป็นรายข้อ (user อาจวางหลายข้อติดกัน) · ตั้ง id ต่อข้อ: ใช้เลขชุด/ข้อ
   ถ้า user แนบมา (เช่น `TD set3·Q17`) ไม่มีก็ใช้ `YYYYMMDD-<slug สั้น>` — ห้ามหยุดถาม
2. อ่าน `index.md` + `keyword-map.md` + `traps.md` + `glossary.md` เต็มไฟล์
   แล้ว Grep `patterns/*.md` ด้วยชื่อ service / คำ signal ของข้อนั้น
3. ต่อข้อ สกัด 4 อย่าง:
   - **Pattern** — signal → architecture ที่ถูก + เหตุผลไทย 2–4 บรรทัด
   - **Traps** — ทำไม option อื่นผิด, generalize เป็นรูปแบบใช้ซ้ำได้ ไม่ผูกกับข้อเดียว
   - **Keyword rows** — แถวใหม่/แถวเดิมที่ควรขยาย
   - **ศัพท์** — เฉพาะคำที่ยังไม่มีใน glossary
4. Merge ตามกติกาด้านล่าง (merge ก่อน สร้างทีหลัง)
5. **Currency check** — เทียบ service ในคำตอบ+keyword rows กับ
   `docs/tickets/aws-currency-checklist.md` (ban list + §9 CONFLICT): ถ้าเฉลยใช้
   service ที่ AWS ปิดรับ/เปลี่ยนชื่อ/แทนที่แล้ว ให้จดตามเฉลย แต่เพิ่มบรรทัด
   `⚠ currency: ข้อสอบเฉลย X — ปัจจุบัน … ของจริงวันนี้คือ Y` ใน entry
6. เขียนตาม template แล้วอัปเดต `index.md`: สถิติ (recompute — ห้าม +1), วันที่วันนี้,
   ลิงก์ไฟล์หมวดใหม่ถ้าเพิ่งสร้าง
7. ตอบ user เป็นไทย: bullet สั้น ๆ ต่อข้อว่าอะไรใหม่/อะไร merge เก็บไว้ที่ไหน แล้ว
   **จบด้วย quiz MCQ 2–4 ข้อ** (ก/ข/ค, distractor ต้อง plausible) ทวนสิ่งที่เพิ่งเก็บ
8. **Follow-up loop (ห้ามข้าม)** — หลังสรุป/เฉลย quiz:
   - user ถามต่อเรื่องไหน → ตอบให้เข้าใจ แล้วถือว่านั่นคือจุดที่ยังไม่เข้าใจ
     **จดลง `review-again.md` ทันทีโดยไม่ต้องรอสั่ง** (template ด้านล่าง)
   - quiz ตอบผิดข้อไหน → เฉลยพร้อมอธิบายว่าตัวเลือกที่ผิด ผิดเพราะอะไรทุกตัว
     แล้วจดจุดนั้นลง `review-again.md` เหมือนกัน
   - จุดเดิมสะดุดซ้ำ → อัปเดตรายการเดิม (+วันที่ล่าสุด) ไม่สร้างรายการซ้ำ

## กติกา merge / กันซ้ำ (ผลลัพธ์ต้อง idempotent)

- **Pattern** — มี entry ความหมายเดียวกันอยู่แล้ว → เสริมของเดิม (signal variant /
  trap / ref / เพิ่ม id ใน `ที่มา`) ห้ามสร้าง entry ใหม่
- **Trap** — เทียบด้วย**ความหมาย** ไม่ใช่ชื่อ slug · trap เดิมโผล่ในข้อใหม่ → บวก
  `**เจอ:** N` · เรียง trap ในไฟล์ตามตัวนับ มาก→น้อย (เท่ากันคงลำดับเดิม)
- **Keyword row** — signal ซ้ำ/ใกล้เคียงแถวเดิม → ขยายแถวเดิม ห้ามเพิ่มแถวซ้ำ
- **Glossary** — เพิ่มเฉพาะคำใหม่ แทรกตามลำดับ A→Z
- **โจทย์เดิมวางซ้ำ** (id ตรง หรือ pattern+traps เหมือนเดิมทั้งชุด) → merge อย่างเดียว
  **ห้ามบวกตัวนับและสถิติ**
- **สถิติใน index.md ห้าม increment** — recompute จากไฟล์จริงทุกรอบ:
  โจทย์ = จำนวน id ไม่ซ้ำในบรรทัด `ที่มา` รวมทุก patterns/ · pattern entries =
  จำนวน `## ` ใน patterns/*.md · traps = จำนวน `## ` ใน traps.md · ศัพท์ = จำนวนแถว
  ตาราง glossary · จุดอ่านซ้ำ = จำนวน `## ` ใน review-again.md
  (วางโจทย์ซ้ำหรือรอบก่อนนับพลาด ก็ converge ที่เลขเดียวกัน)

## Templates

### pattern entry (`patterns/<category>.md`)

```markdown
## <ชื่อ pattern: flow หรือหลักการ สั้น ๆ>

**Signal:** <คำ/วลีในโจทย์ที่บ่งบอก คั่นด้วย ·>

**คำตอบ:** <architecture ที่ถูก เขียนเป็น flow A → B → C>

**เหตุผล:** <ไทย 2–4 บรรทัด อ่านรู้เรื่องโดยไม่ต้องเห็นโจทย์เดิม>

**Trap ในข้อนี้:**

- <option ผิดแบบย่อ → ทำไมผิด> ([trap-slug](../traps.md#trap-slug))

**เก็บเพิ่ม:** <fact ข้างเคียงที่มีประโยชน์ ถ้ามี — ไม่มีให้ตัดหัวข้อนี้ทิ้ง>

**ศัพท์:** <ทุกคำที่ entry นี้เพิ่มเข้า glossary คั่นด้วย ·> → [glossary](../glossary.md)

<sub>ที่มา: <id> · refs: <URL จากเฉลย หรือ official AWS docs ที่แน่ใจว่ามีจริงเท่านั้น
— ห้ามเดา> · <YYYY-MM-DD></sub>
```

### trap (`traps.md`)

```markdown
## <trap-slug-english>

<อธิบายไทย 2–4 บรรทัด: โจทย์แบบไหน ตัวเลือกหลอกหน้าตายังไง ทำไมผิด>

**เจอ:** N ครั้ง
```

### keyword-map row — ตาราง 3 คอลัมน์ใต้ heading หมวด

```markdown
| Signal ในโจทย์ | นึกถึง | ระวัง |
```

### glossary row — ตาราง 2 คอลัมน์ `| ศัพท์ | คำอธิบาย |`

### review-again item (`review-again.md`)

```markdown
## <หัวข้อสั้นของจุดที่สะดุด>

- **เมื่อ:** <YYYY-MM-DD> · **entry:** [<ชื่อ pattern>](patterns/<file>.md)
- **ที่สงสัย:** <คำถามของ user / ข้อ quiz ที่พลาด แบบ paraphrase>
- **คำตอบสั้น:** <ไทย 2–4 บรรทัด พออ่านซ้ำแล้วเคลียร์>
```

### index.md บรรทัดสถิติ (recompute เสมอ)

```markdown
- โจทย์ที่เก็บแล้ว: **N** · pattern entries: **N** · traps: **N** · ศัพท์: **N** · จุดอ่านซ้ำ: **N**
- อัปเดตล่าสุด: YYYY-MM-DD
```

## กติกาการเขียน

- ไทย + ศัพท์เทคนิคอังกฤษ · entry ละ ≤ ~15 บรรทัด · อ่านเข้าใจโดยไม่ต้องเห็นโจทย์เดิม
- **PARAPHRASE เท่านั้น** — ห้าม copy ข้อความโจทย์/เฉลยของ practice exam ลงไฟล์ตรง ๆ
  และห้ามเก็บตัวโจทย์เต็ม ๆ (ลิขสิทธิ์) เก็บเฉพาะเทคนิคที่สรุปใหม่แล้ว
- URL ใช้เฉพาะที่แนบมากับเฉลย หรือ official AWS docs ที่แน่ใจว่ามีจริง — ห้ามเดา URL
  (กติกาเดียวกับ template ด้านบน)
- ชื่อ service ไม่แน่ใจ → เช็ค docs.aws.amazon.com ประกอบ
  `docs/tickets/aws-currency-checklist.md`

## เมื่อไหร่ต้อง escalate (fable)

เฉลยดูขัดแย้งกันเอง / สงสัยว่าเฉลยผิด / เทคนิคซับซ้อนจนอธิบายไม่มั่นใจ / user ขอ
audit โน้ตทั้งชุด → ทำตาม escalation path ใน CLAUDE.md: **ขอ permission จาก user
ก่อนทุกครั้ง ไม่ใช่ routine** แล้วค่อย spawn one-off subagent `model: fable` พร้อม
โจทย์+เฉลยเต็ม ๆ เป็น context · ระหว่างรอ อย่าบันทึกสิ่งที่ยังสงสัยเป็นโน้ตปกติ —
ติดป้าย `⚠ ต้อง verify` ไว้ใน entry

## Git

- ระหว่าง session แก้ไฟล์ไปเรื่อย ๆ **ไม่ commit ต่อข้อ**
- เมื่อ user บอกจบรอบ ("พอแค่นี้" / "commit" / "ปิดรอบ") หรือก่อนจบ session:
  branch `notes/aws-exam-YYYYMMDD` → PR เข้า `develop` ตาม template ปกติ,
  Review level 🟢 (โน้ต additive ไม่มีโค้ด), body สรุปรายการ entry ที่เพิ่ม/merge
  + Review focus เป็นคำถามตามกติกา CLAUDE.md
