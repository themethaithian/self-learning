# AWS Certified Solutions Architect – Associate (SAA-C03) — แผนสอบ

> **Priority อันดับ 1 ปัจจุบัน** — สอบใน 1–2 เดือน, อ่านวันละ 1–2 ชม., ยังไม่เคยจับ AWS จริง
> (รู้แค่ชื่อ service หลัก ๆ) · Q-2d / Q-3 พักไว้บางส่วนระหว่างนี้ (ดูหัวข้อ blocker ด้านล่าง — Q-3
> ถูก unpause บางส่วนเพราะ schema change (2) คือ schema half ของ Q-3 เอง) · SIM-\* พักเต็ม
> ([roadmap.md](../roadmap.md))

## Ticket ที่เกี่ยวข้อง

- **AWS-0** (เอกสารนี้) — รีบาลานซ์ `content/curriculum/aws.json` ให้ตรง exam guide จริง + pin แผนนี้
  (round 2: แก้ currency issues + content-model blocker หลัง code review)
- **AWS-S1** (code review round 2 แก้แล้ว, PR pending — ดูหัวข้อ "AWS-S1 — explanation + raised
  check ceiling" ด้านล่าง) — `maxRecallChecks` 5→15 + `explanation` end-to-end คนละ ticket จาก
  AWS-0 เพราะแตะ `internal/` ซึ่ง AWS-0 ไม่แตะ
- **AWS-S2** (ยังไม่เริ่ม) — multiple-response support (second correct answer, "Select TWO",
  all-or-nothing scoring) — ดูหัวข้อ blocker ข้อ 3 ด้านล่าง
- **AWS-S3** (implemented, PR pending — ดูหัวข้อ "AWS-S3 — guessability measurement tool" ด้านล่าง) —
  guessability measurement tool ที่ parameterize baseline ได้ต่อ corpus, gate คลังข้อสอบจริงใช้งานได้แล้ว
- **AWS-C0** (implemented, pending code review) — landed the 2 pilot lessons
  (`vpc-fundamentals`, `s3-security`) that lock the lesson-writer template and fix the
  count contradiction below (~11/550 → 4–6/~250) — ดูหัวข้อ "Template" และ "AWS-C0" Status
  ด้านล่าง
- **AWS-C1..C4** (ยังไม่เริ่ม, gated บน S1/S2/S3 + AWS-C0's template) — คลังข้อสอบจริง ~250 ข้อ
  ทีละ domain (4–6 ข้อ/concept)

## Blueprint ข้อสอบ (verified โดยตรงจาก docs.aws.amazon.com ระหว่าง round 2 ของ ticket นี้)

**SAA-C03 คือเวอร์ชันปัจจุบัน ไม่มี SAA-C04** — หน้าเว็บบุคคลที่สามที่อ้างว่ามี SAA-C04 เป็นข้อมูลผิด
ยืนยันจาก https://docs.aws.amazon.com/aws-certification/latest/solutions-architect-associate-03/solutions-architect-associate-03.html

- 65 ข้อ (50 ข้อคิดคะแนน + 15 ข้อ unscored ปนอยู่โดยไม่บอกว่าข้อไหน)
- 130 นาที · ผ่านที่ 720/1000 (scaled score, 100–1,000)
- **Compensatory scoring** — คะแนนรวมต้องผ่าน ไม่ได้ตัดเป็นรายหมวด (คำต่อคำจากคู่มือ: "You need to
  pass only the overall exam")
- **ไม่มีการหักคะแนนจากการเดา** ("Unanswered questions are scored as incorrect; there is no
  penalty for guessing")
- รูปแบบคำถาม: **Multiple choice** (1 คำตอบถูก + 3 distractor) และ **Multiple response** (2+
  คำตอบถูกจาก 5+ ตัวเลือก) — stem จะบอกจำนวนที่ต้องเลือกเสมอ เช่น "(Select TWO.)"

แหล่งอ้างอิงชื่อ service ที่ถูกต้องคือ **docs.aws.amazon.com** ไม่ใช่ PDF จาก `d1.awsstatic.com`
(PDF v1.1 ล้าสมัยอย่างน้อย 2 จุด: ยังเขียน "AWS IAM Identity Center [AWS Single Sign-On]" และ
เขียน QuickSight ในจุดที่ docs ปัจจุบันใช้ "Amazon Quick"/"Amazon QuickSuite") — ถ้าต้องเช็คชื่อ
service ให้ fetch docs HTML ไม่ใช่เชื่อ PDF ทุก batch เขียนข้อสอบ**ต้อง**เปิด
[`docs/tickets/aws-currency-checklist.md`](aws-currency-checklist.md) (branch
`docs/aws-currency-checklist`) เป็น context บังคับด้วย — ไฟล์นั้นมีรายการ service ที่ปิดรับลูกค้าใหม่/
เปลี่ยนชื่อ/เปลี่ยนราคา ที่โมเดลเขียนคำตอบผิดถ้าไม่เช็ค

**Sequencing ที่ต้องรู้**: ลิงก์ด้านบนเป็น relative link ไปยังไฟล์ที่ยังอยู่แค่บน branch
`docs/aws-currency-checklist` (PR #54) เท่านั้น — บน `develop` ตอนนี้ลิงก์นี้ตายเพราะไฟล์ยังไม่ merge
**PR #54 ต้อง merge เข้า develop ก่อน AWS-S1 เริ่มงานเสมอ** ไม่ใช่แค่ก่อนเขียนคำถามจริง

### Domain, น้ำหนัก และ task statement ทั้ง 14 ข้อ (verbatim จาก official exam guide)

| Domain | น้ำหนัก | Task statements (verbatim) |
|---|---|---|
| 1. Design Secure Architectures | **30%** | 1.1 *Design secure access to AWS resources* · 1.2 *Design secure workloads and applications* · 1.3 *Determine appropriate data security controls* |
| 2. Design Resilient Architectures | **26%** | 2.1 *Design scalable and loosely coupled architectures* · 2.2 *Design highly available and/or fault-tolerant architectures* |
| 3. Design High-Performing Architectures | **24%** | 3.1 *Determine high-performing and/or scalable storage solutions* · 3.2 *Design high-performing and elastic compute solutions* · 3.3 *Determine high-performing database solutions* · 3.4 *Determine high-performing and/or scalable network architectures* · 3.5 *Determine high-performing data ingestion and transformation solutions* |
| 4. Design Cost-Optimized Architectures | **20%** | 4.1 *Design cost-optimized storage solutions* · 4.2 *Design cost-optimized compute solutions* · 4.3 *Design cost-optimized database solutions* · 4.4 *Design cost-optimized network architectures* |

แหล่งที่มาต่อ domain (verified 2026-07-30 ด้วย WebFetch ตรง ไม่ใช่จำจาก training data):
[Domain 1](https://docs.aws.amazon.com/aws-certification/latest/solutions-architect-associate-03/solutions-architect-associate-03-domain1.md) ·
[Domain 2](https://docs.aws.amazon.com/aws-certification/latest/solutions-architect-associate-03/solutions-architect-associate-03-domain2.md) ·
[Domain 3](https://docs.aws.amazon.com/aws-certification/latest/solutions-architect-associate-03/solutions-architect-associate-03-domain3.md) ·
[Domain 4](https://docs.aws.amazon.com/aws-certification/latest/solutions-architect-associate-03/solutions-architect-associate-03-domain4.md)

**Note สำหรับ session ถัดไป (est_minutes)**: schema ของ `content/curriculum/*.json`
(`conceptFileDTO` ใน `internal/curriculum/infra/contentfile.go`) มีแค่ `slug`/`title`/`outline`/
`position` — **ไม่มีฟิลด์ `est_minutes` ในไฟล์ curriculum เลย** (มันอยู่ในฟิลด์ `Lesson` ของไฟล์
`content/lessons/*.json` คนละ schema กัน) และ `LoadTopic`/`LoadLesson` ทั้งคู่ใช้
`dec.DisallowUnknownFields()` เพิ่ม field ใหม่แบบเงียบ ๆ ไม่ได้ — เป็นแค่ข้อเท็จจริงที่ต้องรู้ ไม่ใช่
บั๊กที่ต้องแก้

### สองข้อสรุปเชิงกลยุทธ์ที่ต้องจำ

1. **ไม่มีการหักคะแนนจากการเดา → ห้ามเว้นข้อว่างเด็ดขาด** ถ้าไม่แน่ใจให้เดาตัวที่ตัดตัวเลือกผิด
   ออกได้มากที่สุด ดีกว่าเว้นว่างเสมอ (เว้นว่าง = ผิดชัวร์, เดา = มีโอกาสถูก)
2. **Compensatory scoring → มองภาพรวม ไม่ต้องกลัวหมวดใดหมวดหนึ่งอ่อน** สิ่งที่ต้องบริหารคือ
   คะแนนรวม 720/1000 ไม่ใช่การผ่านแยกรายหมวด — เวลาอ่านจึงควรกระจายตามน้ำหนัก (30/26/24/20)
   ไม่ใช่ทุ่มเวลาให้หมวดที่ถนัดจนหมวดอื่นได้เวลาน้อยเกินไป

## ⚠️ Blocker: แผนเดิม (550 ข้อ, batch ระดับ domain) รันจริงกับ content model ตอนนี้ไม่ได้

**นี่คือ planning error ของรอบก่อน ไม่ใช่ของ AWS-0's curriculum tree** — tree เองไม่มีปัญหา

`internal/curriculum/domain/lesson.go` บังคับ `minRecallChecks = 3`, `maxRecallChecks = 5` ใน
`NewLesson` คำถามหนึ่งข้อมีที่อยู่ได้ทางเดียวคือเป็น `RecallCheck` ข้างใน `Lesson` หนึ่งใบ และ
lesson หนึ่งใบผูกกับ concept เดียว — เพดานตามจริงคือ **50 concept × 5 = 250 ข้อ** ไม่ใช่ 550
กำแพงที่สองและสาม: `recallCheckFileDTO` ไม่มีฟิลด์ `explanation` และ `LoadLesson` ใช้
`DisallowUnknownFields()` เหมือนกัน ไฟล์ที่มี field นี้เข้ามาจะ**ถูก reject ทั้งไฟล์ตอน import**
ไม่ใช่แค่ field นั้นถูก strip; และ `RecallCheck.expectedAnswer` เป็น `string` เดี่ยวที่ validate ด้วย
`slices.Contains(options, answer)` — multiple-response (คำตอบถูกได้มากกว่า 1) จึงเป็นไปไม่ได้ในเชิง
โครงสร้างด้วย

**การตัดสินใจที่ปักไว้ตรงนี้**: เก็บ `recall_checks` ไว้แล้วขยาย ไม่สร้าง question store คู่ขนาน —
`recall_attempts` มี FK ไป `lessons.id` อยู่แล้ว และ Q-2a/Q-2b/Q-2c (persist attempt, wire หน้า quiz,
SM-2 scheduling) ทั้งหมดสร้างบน path นี้ ถ้าเปิด store ใหม่ต้องสร้างทั้งสามอย่างใหม่

**รูปทรงใหม่**: **1 concept = 1 lesson (concept note ภาษาไทยสั้น ๆ 5–10 นาที ตาม CLAUDE.md — นี่คือ
"บทเรียน" ที่ผู้ใช้อยากได้ด้วย) + recall_checks ต่อ concept** (คำถามของ concept นั้น)

**อัปเดตเป้า (AWS-C0, 2026-07-30) — 4–6 ข้อ/concept, รวม ~250 ข้อ ไม่ใช่ ~11/550 อีกต่อไป**: ตัวเลข
เดิม (~11 ข้อ/concept, 550 ข้อรวม) มาจากตอนที่แผนนี้ยังคิดว่าแอปต้อง **จำลองข้อสอบเต็มรูปแบบ** เหมือน
คลังข้อสอบเชิงพาณิชย์ — ข้อเท็จจริงตอนนี้คือ **แอปนี้ไม่ใช่ exam simulator**: ผู้ใช้ทำข้อสอบแนว SAA จริง
อยู่แล้วบน third-party platform ที่ซื้อไว้ต่างหาก งานของแอปนี้คือคอร์สภาษาไทย + recall loop แบบ SM-2
คำถามสั้นแบบ decision-rule ที่ SM-2 หมุนซ้ำได้เร็ว มีค่ามากกว่า scenario ยาวที่แข่งกับ product ที่ผู้ใช้
มีอยู่แล้ว บันทึกไว้ตรงนี้ **ทำไมเลขถึงย้าย** เพื่อไม่ให้ใครเอา ~11/550 กลับมาใช้ทีหลังโดยไม่รู้ตัวว่าเป็น
เป้าที่ตกไปแล้ว — รายละเอียดการตัดสินใจอื่นอยู่ที่หัวข้อ "Template" ด้านล่าง

| Domain | concepts | × 4–6 | จำนวนข้อ |
|---|---|---|---|
| 1. Secure Architectures | 15 | ×4–6 | **60–90** |
| 2. Resilient Architectures | 13 | ×4–6 | **52–78** |
| 3. High-Performing Architectures | 12 | ×4–6 | **48–72** |
| 4. Cost-Optimized Architectures | 10 | ×4–6 | **40–60** |
| **รวม** | 50 | ×4–6 | **200–300 (เป้ากลาง ~250 ที่ 5/concept)** |

(ตัวเลข ×11/550 เดิมยังโผล่ในหัวข้อ "AWS-S1"/"AWS-S3" ด้านล่างเป็นบริบทประวัติของการตัดสินใจที่ทำตอนนั้น
เช่น ทำไม `maxRecallChecks` ถูกยกจาก 5 เป็น 15 — เก็บไว้ตามที่เกิดขึ้นจริง ไม่ใช่เป้าปัจจุบัน)

Mock exam คือ **sampling policy** เหนือคำถามทั้งหมด ถ่วงน้ำหนัก 30/26/24/20 ตอน sample ไม่ใช่
question store แยก

### 3 schema change ที่ต้องทำก่อน (นี่คือ AWS-S1's งาน ไม่ใช่ AWS-0 — AWS-0 ไม่แตะ `internal/`)

1. ~~**เพิ่ม `maxRecallChecks`** จาก 5 เป็น ~15 (เผื่อ headroom เหนือเป้า 11/concept)~~ — **เสร็จแล้ว
   ใน AWS-S1** (5 → 15) ดูหัวข้อ "AWS-S1 — explanation + raised check ceiling" ด้านล่าง
2. ~~**เพิ่มฟิลด์ `explanation`**: domain (`RecallCheck`) → `recallCheckFileDTO` → importer → API
   response → หน้า quiz ต้อง render~~ — **เสร็จแล้วใน AWS-S1** (schema half ของ Q-3 เท่านั้น — การ
   เขียน explanation ให้ 356 ข้อเดิมยังพักไว้ตามเดิมจนกว่าจะว่าง)
3. **เพิ่ม multiple-response support**: ต้องมี second correct-answer representation (เช่น
   `expectedAnswers []string` หรือ MCQ กับ multiple-response แยก type) + UI ฝั่ง quiz ให้เลือกได้
   มากกว่า 1 ตัวเลือก — multi-response เป็น ~15–20% ของข้อสอบจริงและเป็นประเภทที่คนเสียคะแนนมากที่สุด
   ทิ้งไปไม่ได้ — **ยังไม่ทำ นี่คือ AWS-S2**

**ข้อจำกัดที่ตามมาจนกว่า (3) จะเสร็จ**: multiple-response ยังเขียนไม่ได้เชิงโครงสร้าง
(`expected_answer` ยังเป็น string เดี่ยว) — ห้ามเขียนคำถามที่ต้องการ "(Select TWO.)" ก่อน AWS-S2
landed (explanation ตอนนี้เขียนได้แล้วหลัง AWS-S1 — ข้อจำกัดเดิมเรื่อง `DisallowUnknownFields()`
reject ทั้งไฟล์หมดไปแล้ว)

## Concept → Task statement mapping (tree ที่รีบาลานซ์แล้ว, 50 concept)

Schema ปัจจุบันของ `content/curriculum/aws.json` มีแค่ `slug`/`title`/`outline`/`position` — ไม่มี
ช่องสำหรับ task statement และตามกติกาของ ticket นี้ **ห้ามเพิ่ม field ใหม่ในไฟล์นั้น** เพราะจะ
กระทบ parser และทุก content file อื่น ๆ ที่ใช้ schema เดียวกัน mapping จึงอยู่ในตารางนี้แทน

**Caveat ที่ยอมรับตรง ๆ**: chapter ที่ concept สังกัดคือ domain สำหรับนับสัดส่วน 30/26/24/20 แต่
**ไม่ใช่ตัวบอก task statement ที่ concept นั้น serve เสมอไป** — `vpc-fundamentals` และ
`hybrid-cross-vpc-connectivity` อยู่ใน chapter Domain 1 แต่ตารางด้านล่างให้ทั้งคู่ทำหน้าที่ 3.4/4.4
ด้วย เพราะเนื้อหาเดียวกัน (VPC/Direct Connect/VPN mechanics) ต้องใช้ตอบทั้งสามโดเมน และการย้ายไป
chapter อื่นจะทำให้สัดส่วน concept-count ต่อ domain เพี้ยน ดังนั้น **per-chapter concept count
ตรงกับสัดส่วนน้ำหนักโดเมนพอดี (30/26/24/20) แต่ไม่ได้แปลว่าเวลาอ่านต่อ task statement กระจายตาม
สัดส่วนเดียวกันเป๊ะ ๆ** — ใช้ตารางนี้ (ไม่ใช่แค่ chapter) เป็นตัวตัดสินว่าข้อสอบข้อไหน tag task ไหน

### Domain 1 — Secure Architectures (15 concepts)

| Concept slug | Task |
|---|---|
| iam-users-roles-policies | 1.1 |
| iam-policy-evaluation | 1.1 |
| organizations-scp | 1.1 |
| cognito | 1.2 (verified: Task 1.2's Knowledge ระบุ "Security services with appropriate use cases (for example, AWS Cognito, AWS GuardDuty, AWS Macie)" — Cognito เป็นชื่อแรกในวงเล็บเดียวกับ GuardDuty/Macie ที่ย้ายไปแล้ว แก้จาก 1.1 เดิม) |
| kms-encryption | 1.3 |
| secrets-vs-parameter-store | 1.2 |
| sg-vs-nacl | 1.2 |
| vpc-endpoints-privatelink | 1.2, 3.4 (Task 3.4's Knowledge ระบุ "Network connection options (for example, AWS VPN, AWS Direct Connect, AWS PrivateLink)" ตรงตัว) |
| waf-shield | 1.2 |
| s3-security | 1.3 |
| vpc-fundamentals | 1.2 (prerequisite ให้ 3.4, 4.4 ด้วย) |
| hybrid-cross-vpc-connectivity | 1.2, 3.4 |
| security-detection-governance | 1.2 (verified: GuardDuty/Macie ระบุชื่อตรงใน Task 1.2's Knowledge — แก้จาก 1.1 เดิมหลัง verify) |
| multi-account-access-governance | 1.1 |
| acm-tls-in-transit-encryption | 1.3 |

### Domain 2 — Resilient Architectures (13 concepts)

| Concept slug | Task |
|---|---|
| regions-az-edge | 2.2 |
| elb-types | 2.1 |
| auto-scaling-groups | 2.1, 2.2 (Task 2.2's Knowledge ระบุ "Immutable infrastructure" ตรงตัว — outline ของ concept นี้ fold เรื่องนี้ไว้แล้ว) |
| rds-multi-az-read-replicas | 2.2, 2.1 (Task 2.1's Skills ระบุ "When to use read replicas" ตรงตัว) |
| aurora-ha | 2.2 |
| sqs-sns-decoupling | 2.1 |
| eventbridge | 2.1 |
| route53-routing-policies | 2.2 |
| dr-strategies | 2.2 |
| backup-strategies | 2.2 |
| observability-cloudwatch-cloudtrail | 2.2 |
| api-gateway-step-functions | 2.1 |
| managed-ai-services-overview | 2.2 (verified ตรงกับคู่มือ: Task 2.2's Knowledge ระบุ "AWS Managed Services (AMS) with appropriate use cases (for example, Amazon Comprehend, Amazon Polly)" ตรงตัว — ดูหมายเหตุ round-2 ด้านล่าง) |

### Domain 3 — High-Performing Architectures (12 concepts)

| Concept slug | Task |
|---|---|
| ec2-families-purchasing | 3.2 |
| ebs-efs-instance-store | 3.1 |
| s3-performance | 3.1 |
| cloudfront | 3.4 |
| elasticache | 3.3 |
| dynamodb-fundamentals | 3.3 |
| dynamodb-advanced | 3.3 |
| rds-performance | 3.3 (RDS Proxy บรรทัดเสริมยัง serve 2.2 ด้วย — คู่มือระบุ "Proxy concepts (for example, Amazon RDS Proxy)" ใน Task 2.2 ไม่ใช่ 3.3) |
| kinesis | 3.5 |
| athena-glue | 3.5 |
| lambda-performance | 3.2 |
| ecs-eks-fargate | 3.2 |

Domain 3 เดิมมี 12 concept กระจุกที่ 3.1/3.2/3.3 อยู่แล้ว (9 จาก 12) — รอบนี้**ไม่เพิ่มจำนวน**
แต่ enrich outline ของ concept ที่มีอยู่แทน เพื่อดันน้ำหนักไปทาง 3.4/3.5 โดยไม่ทำให้ Domain 3
เกิน 24% ของ tree (12/50 = 24% พอดีอยู่แล้ว ตรงข้ามกับก่อนรีบาลานซ์ที่ 12/37 = 32%) — round 3 เติม
AWS Lake Formation และ Amazon QuickSuite (Task 3.5's Knowledge ระบุตรงตัว: "Data analytics and
visualization services with appropriate use cases (for example, Amazon Athena, AWS Lake
Formation, Amazon QuickSuite)") เข้า `athena-glue`'s outline แทนที่จะแยก concept ใหม่ เพราะ Task
3.5 มีแค่ 2 concept (`kinesis`, `athena-glue`) คุม ~22 ข้อ และก่อนหน้านี้ไม่มี concept ไหนแตะ data
lake governance หรือ visualization เลย

### Domain 4 — Cost-Optimized Architectures (10 concepts)

| Concept slug | Task |
|---|---|
| pricing-models-ri-sp-spot | 4.2 |
| s3-storage-classes-lifecycle | 4.1 |
| compute-cost-optimization | 4.2 |
| data-transfer-costs | 4.4 |
| cost-tools-budgets | 4.1–4.4 (cross-cutting) |
| cost-optimized-databases-capacity | 4.3 |
| cost-optimized-databases-storage-lifecycle | 4.3 |
| cost-optimized-networking | 4.4 |
| migration-and-transfer-services | 4.1, 3.5 (DMS component ยัง serve 4.3 — Task 4.3's Skills ระบุ "Migrating database schemas and data to different locations and/or different database engines" ตรงตัว) |
| storage-cost-optimization-beyond-s3 | 4.1 |

### Small named items folded into existing concepts' outline (ไม่ได้แยก concept ใหม่)

| Item | Folded into |
|---|---|
| AWS Resource Access Manager (RAM) | `multi-account-access-governance` (ย้ายจาก `organizations-scp` ใน round 3 เพื่อลด outline bloat) |
| Resource Control Policies (RCPs, ใหม่ 2024-11-13) | `organizations-scp` |
| IAM Permissions Boundary (แยกจาก SCP) | `iam-policy-evaluation` (ย้ายจาก `organizations-scp` ใน round 3 — เข้าธีม policy-layer evaluation ตรงกว่า) |
| CloudFormation / immutable infrastructure, Elastic Beanstalk | `auto-scaling-groups` |
| AWS Outposts | `regions-az-edge` (ยัง serve **4.2** ด้วย — Task 4.2's Knowledge ระบุ "Hybrid compute options (for example, AWS Outposts)" ตรงตัว) |
| Amazon MQ | `sqs-sns-decoupling` |
| AWS Systems Manager Session Manager | `secrets-vs-parameter-store` |
| S3 Versioning / MFA Delete / Object Lock | `s3-security` |
| S3 Intelligent-Tiering | `s3-storage-classes-lifecycle` |
| EC2 hibernation | `compute-cost-optimization` (ย้ายจาก `ec2-families-purchasing` ใน round 2 — verified ว่าคู่มือระบุ hibernation ใน Task 4.2 ไม่ใช่ 3.2) |
| Amazon FSx (all types) | `ebs-efs-instance-store` |
| RDS Proxy | `rds-performance` |
| AWS Batch, Amazon ECR | `ecs-eks-fargate` |
| Amazon EMR | `athena-glue` |
| AWS Lake Formation | `athena-glue` (ใหม่ round 3 — Task 3.5) |
| Amazon QuickSuite (เดิมชื่อ QuickSight ก่อน 2025-10-09) | `athena-glue` (ใหม่ round 3 — Task 3.5) |
| DocumentDB, Neptune, Keyspaces (purpose-built DB) | `dynamodb-advanced` |
| AWS Trusted Advisor | `cost-tools-budgets` |
| AWS Cost and Usage Report (CUR) | `cost-tools-budgets` (ใหม่ round 3 — ระบุชื่อในทุก task 4.x, ปรากฏบ่อยกว่า named service อื่นใดใน Domain 4) |
| Multi-account billing (consolidated billing) | `cost-tools-budgets` (ใหม่ round 3 — ระบุชื่อในทุก task 4.x) |
| Requester Pays (S3) | `s3-storage-classes-lifecycle` (ใหม่ round 3 — Task 4.1's Knowledge ระบุตรงตัว) |
| public IPv4 hourly charge ($0.005/IP/ชม. ตั้งแต่ 2024-02-01) | `data-transfer-costs` (ย้ายจาก `cost-optimized-networking` ใน round 3 เพื่อลด outline bloat — ยังอยู่ task 4.4 เหมือนเดิม) |

## แผนคลังข้อสอบ (~250 ข้อ + 50 lesson)

จำนวนข้อต่อ domain มาจากตาราง "รูปทรงใหม่" ด้านบน (4–6 ข้อ/concept × concept ต่อ domain, เป้ากลาง
~250 — ปรับลดจาก ~11/550 ใน AWS-C0 ดูเหตุผลที่หัวข้อนั้น) ไม่ใช่การหารตามเปอร์เซ็นต์น้ำหนักตรง ๆ —
ตัวเลขบังเอิญใกล้เคียงสัดส่วนน้ำหนักเพราะ concept count ต่อ domain ถูกออกแบบให้ตรงสัดส่วนอยู่แล้ว
(15/13/12/10 = 30/26/24/20%)

**หน่วยของการทำงานคือ concept ไม่ใช่ domain** ตาม token-efficiency skill
(`.claude/skills/token-efficiency/SKILL.md`: "lesson-writer/verifier: one concept per invocation")
— แผนเดิมที่ให้ `lesson-writer` เขียนทีเดียว 132–165 ข้อ/domain ละเมิดกฎนี้ตรง ๆ และเป็นสาเหตุคลาสสิก
ของ explanation ที่ตื้น/restate นิยามซ้ำ (โมเดลล้าเมื่อ context ยาว):

- **1 invocation ต่อ 1 concept**: `lesson-writer` เขียน lesson (concept note) + recall_checks
  4–6 ข้อของ concept นั้น (เพดานหลัง schema change (1) คือ 15, เผื่อ headroom เท่านั้น ไม่ใช่เป้า)
  พร้อมกัน → `lesson-verifier`
  fact-check ทันที ตาม convention เดิมของ track อื่น (DDD/DDIA/AI-systems) — FAIL ต้อง regenerate
  เฉพาะ concept นั้นก่อนไปต่อ ไม่ retry ทั้ง batch
- รวม **50 invocation-pair** ทั้งหมด (ไม่ใช่ 4 batch ระดับ domain แบบแผนเดิม)
- จัด ticket AWS-C1..C4 เป็น **1 ticket ต่อ domain** เพื่อความสะดวกเวลารีวิว (ticket ละ 10–15
  invocation-pair) แต่หน่วยการทำงานจริงข้างในคือ concept — การจัด ticket ระดับ domain ไม่ใช่การ
  สั่งให้เขียนทีเดียวทั้ง domain

## กติกาการเขียนคำถาม (ticket เขียนคลังข้อสอบทุกตัวต้องทำตาม)

- **Stem แบบ scenario** สไตล์ข้อสอบจริง ไม่ใช่ถามนิยามตรง ๆ — ตั้งบริบทธุรกิจ/ข้อจำกัดก่อนถาม
- ใช้ **qualifier** ที่ทำให้ตัวเลือกที่ทำงานได้จริงหลายตัวกลายเป็นตัวเลือกที่ผิด เช่น **MOST
  cost-effective / LEAST operational overhead / MOST secure** — โจทย์ SAA มักมีคำตอบที่ "ใช้ได้"
  มากกว่า 1 ตัว แต่มีตัวเดียวที่ตรง qualifier
- **Distractor ต้องเป็น AWS service จริงที่ใช้ผิดบริบท** ไม่ใช่ตัวเลือกที่มองออกชัด ๆ ว่าไม่เกี่ยว
- **Multiple-response ต้องระบุจำนวนใน stem เสมอ** เช่น "(Select TWO.)" — ใช้ได้จริงหลัง schema
  change (3) เท่านั้น
- **ทุก batch ต้องเปิด [`aws-currency-checklist.md`](aws-currency-checklist.md) ก่อนเขียน** และ
  ห้ามใช้ service/ชื่อใดที่อยู่ใน "Ban list" ของไฟล์นั้นเป็นคำตอบที่ถูก (S3 Select, S3 Object
  Lambda, Elastic Transcoder, Snowball/Snowcone, Spot blocks, NAT instance, OAI, Aurora
  Serverless v1, Glacier vault, FSx File Gateway, Simple AD, WAF Classic ฯลฯ)
- **Guessability baseline สำหรับ corpus นี้คือ 25%** (1/4 ตัวเลือก) ไม่ใช่ 33% — corpus เดิมใน
  [`mcq-quality.md`](mcq-quality.md) มี 3 ตัวเลือกจึง baseline 33% เป็นค่าที่ถูกต้องสำหรับ corpus
  นั้นเท่านั้น เอาไปใช้ตรง ๆ กับ corpus 4-option ของ AWS จะเป็น false-green แบบเดียวกับที่ #38 เคย
  รายงาน `VIOLATIONS: NONE` ทั้งที่เดาถูก 89.6% — วัดผิด baseline คือ วัดผิดโจทย์แบบเดียวกัน (ดู
  mcq-quality.md's "บทเรียนที่ 1") กติกา distractor length/position เดิมยังใช้ต่อ แค่เทียบกับ 25%
  แทน 33% — **เกทนี้ enforceable แล้ว** (AWS-S3 เสร็จแล้ว — `go run ./cmd/mcq-guessability -dir
  content/lessons/aws -max-excess <threshold>` ดู "AWS-S3 — guessability measurement tool"
  ด้านล่างสำหรับรายละเอียด tool + ตัวเลขที่วัดได้จากคลังเดิม) — ยังต้องเลือกค่า `<threshold>` ที่แน่นอน
  ก่อนเริ่ม AWS-C1 จริง (spec นี้กำหนดแค่ baseline 25%, ไม่ได้ตรึงตัวเลข excess ที่ยอมรับได้)

## AWS-S3 — guessability measurement tool (spec ที่นี่เท่านั้น — ticket แยกต่างหาก ไม่ทำใน AWS-S1)

**สถานะ**: **implemented, PR pending** — ดูหัวข้อ "AWS-S3 — guessability measurement tool" ท้ายเอกสาร
นี้สำหรับสิ่งที่ deliver จริง (`cmd/mcq-guessability` + `internal/mcqguess`), การตัดสินใจ, mutation
table, และตัวเลขที่วัดได้จากคลังเดิม — สเปกด้านล่างนี้คือของเดิมที่ปักไว้ล่วงหน้า (เก็บไว้อ้างอิง
ไม่ได้แก้ตามการ implement จริง)

**สิ่งที่ AWS-C1..C4 (คลังข้อสอบ) ต้องรอ**: AWS-S3 ต้องเสร็จก่อน batch เขียนคำถามจริงชุดแรกเริ่ม —
ไม่ใช่แค่ก่อนตรวจ ไม่งั้นจะเขียนคำถามเป็นร้อยข้อโดยไม่มีทางรู้ว่าเดาได้ง่ายไปหรือเปล่าจนสายเกินแก้

### Spec ขั้นต่ำ

1. **Input**: อ่าน `content/lessons/*/*.json` ทุกไฟล์ (หรือรับ path/glob เป็น flag เพื่อจำกัด
   scope ต่อ track ได้)
2. **Filter**: ใช้เฉพาะ `recall_checks` ที่ `type == "mcq"` เท่านั้น (`short_answer` ไม่มี
   ตัวเลือกให้เดา ไม่เกี่ยว)
3. **Baseline ต่อข้อ**: `1/len(options)` เสมอ — **ห้าม hardcode 33% หรือ 25%** คำนวณจากจำนวน
   ตัวเลือกจริงของแต่ละข้อ (corpus ที่ผสมข้อ 3 กับ 4 ตัวเลือกกันจะได้ baseline คนละค่าต่อข้อ ซึ่งถูกต้อง)
4. **Heuristic ที่ต้องรายงาน hit rate เทียบ baseline**:
   - **ความยาว**: เดาว่าตัวเลือกที่ยาวที่สุด (หรือสั้นที่สุด) คือคำตอบ — วัด hit rate จริงเทียบ
     baseline ของแต่ละข้อ (ไม่ใช่ baseline คงที่ตัวเดียวทั้ง corpus ถ้า corpus ผสมจำนวนตัวเลือก)
   - **ตำแหน่ง**: เดาตำแหน่งเดิมเสมอ (เช่น index 0) ทุกข้อ — วัด hit rate เดียวกัน
5. **Output**: hit rate ของแต่ละ heuristic เทียบ baseline เฉลี่ยของ corpus ที่วัด (เช่น
   "length heuristic: 33.7% actual vs 33.3% baseline avg — PASS/FAIL ตามเกณฑ์ที่ตั้งไว้")
6. **ต้องใช้ซ้ำได้กับหลาย corpus**: รันแยกได้ต่อ track/ต่อไฟล์ ไม่ใช่ผูกกับ DDIA/AI-systems
   corpus เดิมตายตัว — เหตุผลทั้งหมดของ ticket นี้คือ corpus AWS (4 ตัวเลือก) กับ corpus เดิม
   (3 ตัวเลือก) ต้องวัดแยกกันได้โดยไม่ต้องแก้โค้ด

**ไม่ใช่ scope ของ AWS-S3**: การแก้ distractor ที่วัดแล้วมีปัญหา (เป็นงานของ batch เขียนคำถามเอง
ตามกฎใน mcq-quality.md) — AWS-S3 เป็นแค่เครื่องวัด ไม่ใช่เครื่องแก้

## รูปแบบ explanation (ข้อกำหนดสำคัญที่สุดจากผู้ใช้)

ทุกข้อต้องมี explanation ที่มีครบ 3 ส่วนนี้เสมอ (Thai prose, technical term เป็นอังกฤษ) —
**แต่ห้ามเขียนไฟล์คำถามที่มี field นี้ก่อน schema change (2) landed จริง** (ดูหัวข้อ blocker):

1. **ทำไมคำตอบที่ถูกถึงถูก** — ผูกกับ constraint เฉพาะใน stem โดยตรง ไม่ใช่คำอธิบาย service
   ทั่วไป
2. **ทำไม distractor แต่ละตัวถึงผิด** — ต้องระบุ constraint ใน stem ที่ distractor ตัวนั้นละเมิด
   (ทำครบทุกตัวเลือกที่ผิด ไม่ใช่แค่ตัวเดียว)
3. **Decision rule / discriminator ที่นำกลับไปใช้ซ้ำได้** — สรุปเป็นกฎสั้น ๆ ที่ใช้ตัดสินใจข้อ
   อื่นที่มี pattern เดียวกันได้ ไม่ใช่แค่สรุปข้อนี้ข้อเดียว

**เหตุผลของกติกานี้**: user profile คือ "รู้ชื่อ service แต่ไม่เคยจับจริง" — สิ่งที่ข้อสอบ SAA วัด
คือ **การเลือกระหว่าง service ภายใต้ constraint** ไม่ใช่ความรู้นิยาม ดังนั้น explanation ต้องสอน
วิธี **เลือก** ไม่ใช่แค่บอกว่า service แต่ละตัว "คืออะไร" ซ้ำกับที่ user รู้อยู่แล้ว — ถ้า explanation
อ่านแล้วรู้สึกเหมือนอ่าน definition ซ้ำ แปลว่าเขียนผิดจุด

**อัปเดต (AWS-C0 round 2, ปักรูปแบบเป็น bold header บังคับ)**: สามส่วนข้างบนยังเป็นเนื้อหาที่ต้องมี
ครบเหมือนเดิม แต่ตอนนี้ต้องขึ้นต้นแต่ละส่วนด้วย bold Thai header ตายตัว ไม่ใช่ร้อยแก้วไหลต่อกันแบบ
`s3-security`'s draft แรกที่ใช้ inline (ก)(ข)(ค) — ตามลำดับนี้เป๊ะ:

**โจทย์ถามว่า** → **ทำไมข้อที่ถูกถึงถูก** → **ทำไมตัวอื่นผิด** (ตามด้วย bullet list `- *ชื่อตัวเลือก* —
เหตุผล` ทีละตัวเลือกที่ผิด) → **Decision rule**

รายละเอียดเต็ม + ตัวอย่างจริงอยู่ที่ `.claude/agents/lesson-writer.md`'s "## AWS track" — สอง
เอกสารนี้ต้องพูดตรงกันเสมอ (ก่อนหน้านี้เอกสารนี้ยังเขียนว่า "Thai prose สามส่วน" เฉย ๆ ไม่ระบุชื่อ
header เลย ในขณะที่ `lesson-writer.md` เขียน header ไว้แล้ว — ทำให้ `lesson-verifier` (ที่อ้างอิง
หัวข้อนี้) ไม่มีทางเช็ค header ได้เลย `s3-security`'s ทั้ง 4 ข้อผ่าน verifier ทั้งที่ไม่มี header สักข้อ
เป็นหลักฐานตรงว่าช่องว่างนี้เคยมีจริง)

## หมายเหตุ round 2 (ความเห็นต่างที่ verify แล้ว — ไม่ได้ตามการแก้ทุกจุดแบบไม่เช็ค)

Code review รอบนี้ระบุว่า `managed-ai-services-overview` (Comprehend/Polly/ฯลฯ) "have nothing to do
with HA" และแนะให้ re-map ออกจาก 2.2 — ตอน fix ได้ fetch
[Domain 2's task page](https://docs.aws.amazon.com/aws-certification/latest/solutions-architect-associate-03/solutions-architect-associate-03-domain2.md)
ตรง ๆ และพบว่า Task 2.2's "Knowledge of" list เขียนไว้ตรงตัวว่า **"AWS Managed Services (AMS) with
appropriate use cases (for example, Amazon Comprehend, Amazon Polly)"** — ยืนยันว่า mapping เดิม
(2.2) ถูกต้องสำหรับ Comprehend/Polly อย่างน้อยสองตัว ไม่ใช่ "ไม่เกี่ยวกับ HA" ตามที่ review ระบุ
**การแก้จริงที่ทำ**: ไม่ย้าย task mapping (ยังเป็น 2.2) แต่แก้ outline ให้ foreground เฉพาะ
Comprehend/Polly ที่ verified และลดน้ำหนักของ service อื่น (Rekognition/Lex/Kendra/ฯลฯ) ให้ชัดว่า
"อยู่ใน scope แต่ไม่ได้ระบุชื่อตรงใน task ไหนเป็นการเฉพาะ" — บันทึกไว้ตรงนี้เพื่อความโปร่งใส เผื่อ
ผู้รีวิวมีแหล่งข้อมูลอื่นที่ขัดกับสิ่งที่ fetch ได้ตอนนี้

## Status

**In review (round 3)** — round 2 ได้ REQUEST_CHANGES แบบแคบ (4 gap ที่ verify แล้วกับคู่มือ + dual
mapping ที่ขาด + 2 เรื่องเล็ก) round 3 แก้ครบ: `cognito` remap 1.1→1.2, เติม AWS Lake
Formation/Amazon QuickSuite เข้า `athena-glue` (Task 3.5 เดิมมีแค่ 2 concept คุมทั้ง task),
เติม AWS Cost and Usage Report + multi-account billing เข้า `cost-tools-budgets` และ Requester
Pays เข้า `s3-storage-classes-lifecycle`, บันทึก dual-mapping ที่ verify แล้ว 5 จุด
(`vpc-endpoints-privatelink`→+3.4, `rds-multi-az-read-replicas`→+2.1, `auto-scaling-groups`→+2.2,
migration-and-transfer-services's DMS→+4.3, regions-az-edge's Outposts→+4.2), แก้ title
`ec2-families-purchasing` ที่ยังขัดกับ outline ของตัวเองหลัง round 2 ลดขอบเขตแล้ว, ระบุ sequencing
ว่า PR #54 ต้อง merge ก่อน AWS-S1 เริ่ม — **`aws.json` ยังไม่มีการเพิ่ม/ลบ/ย้าย concept ในรอบนี้เช่นกัน
(ยังคง 15/13/12/10 = 50) จึงไม่ต้องแตะ `contentfile_aws_test.go`**

**Outline bloat**: แก้ด้วยการ redistribute เนื้อหาระหว่าง concept ที่มีอยู่ (ไม่ใช่แยก concept ใหม่ —
รอบนี้ห้ามแตะ Go) — ย้าย IAM Permissions Boundary distinction จาก `organizations-scp` ไป
`iam-policy-evaluation`, ย้าย AWS RAM จาก `organizations-scp` ไป `multi-account-access-governance`,
ย้าย public IPv4 charge จาก `cost-optimized-networking` ไป `data-transfer-costs`, และตัดคำฟุ่มเฟือย
ใน `migration-and-transfer-services`/`managed-ai-services-overview`/`cost-optimized-networking`
ผล: max outline **1846 → 1136 bytes** (`athena-glue`, หลังเติม Lake Formation/QuickSuite),
`organizations-scp` (offender ที่ระบุชื่อตรง) **1462 → 866 bytes**. ยังไม่เท่า ddia's 811 เป๊ะ
(1136 ≈ 1.4×) เพราะข้อจำกัดรอบนี้ (ห้ามแยก concept ใหม่) ทำให้ทำได้แค่ redistribute+trim ไม่ใช่ split —
ถ้าต้องการลดต่ำกว่านี้อีกต้องแยก concept ซึ่งกระทบ golden test และเป็นงานของ ticket ที่แตะ Go

## Review focus (round 3)

<details>
<summary>คำถามสำหรับรีวิว diff รอบนี้ (เฉลยพับไว้ด้านล่าง)</summary>

1. ทำไม AWS RAM ถึงย้ายจาก `organizations-scp` ไป `multi-account-access-governance` แทนที่จะ
   ตัดทิ้งไปเลยเพื่อลด byte เร็วกว่า?
2. ทำไม max outline byte ยังไม่เท่ากับ track อื่น (ddia 811 bytes) ทั้งที่ลดจาก 1846 มาเยอะแล้ว?
3. ทำไม dual-mapping ของ `rds-multi-az-read-replicas` (+2.1) ถึงไม่ทำให้ต้องย้าย concept นี้ไปอยู่
   chapter อื่นหรือปรับ concept count ของ Domain 2?

<details>
<summary>เฉลย</summary>

1. เพราะเนื้อหายังต้องมีที่อยู่ — AWS RAM เป็น item ที่ยืนยันแล้วว่าอยู่ใน scope ข้อสอบ (ticket AWS-0
   เดิมเรียกมันว่า item ที่ "ต้องมีแต่ปริมาณคำถามน้อย ไม่คุ้มแยก concept ใหม่") การตัดทิ้งจะทำให้
   coverage หายไปจริง ไม่ใช่แค่ byte หาย ส่วน `multi-account-access-governance` เป็นบ้านที่เหมาะกว่า
   `organizations-scp` อยู่แล้วเพราะ RAM คือกลไก resource-sharing ระดับ multi-account เหมือนเนื้อหา
   อื่นในนั้น (IAM Identity Center, Control Tower, STS cross-account)
2. เพราะรอบนี้ห้ามแตะ Go ไฟล์ใด ๆ (`contentfile_aws_test.go` ต้องคงเดิม) การลด byte จึงทำได้แค่สอง
   ทาง: redistribute เนื้อหาไปหา concept ที่มีที่ว่าง (ทำไปแล้วกับ organizations-scp,
   cost-optimized-networking) กับตัดคำฟุ่มเฟือยในประโยค (ทำไปแล้วเช่นกัน) — การจะลดต่ำกว่านี้อีก
   ต้องแยก concept ที่มีหลาย topic ออกเป็นสองใบจริง ๆ ซึ่งเปลี่ยนจำนวน concept และต้องแก้ golden
   test ทันที เป็นการเปลี่ยนแปลงเชิงโครงสร้างที่ควรอยู่ใน ticket ที่แตะ `internal/` ได้ (เช่น AWS-S1)
   ไม่ใช่ ticket นี้ที่ scope จำกัดไว้แค่ content + docs
3. เพราะ dual-mapping ที่บันทึกในตารางเป็นแค่ **ข้อมูลสำหรับ tag คำถามตอนเขียนคลังข้อสอบ** (ข้อไหน
   cover task ไหนบ้าง) ไม่ใช่ตัวกำหนดว่า concept ต้องอยู่ chapter ไหน — `rds-multi-az-read-replicas`
   ยังคงเป็นเนื้อหาแกนของ Domain 2 (high availability ผ่าน Multi-AZ) การที่มันมีบรรทัดหนึ่งที่ตอบ
   Task 2.1 ได้ด้วย (read replica ช่วย scale) ไม่ได้แปลว่ามันควรย้ายไป Domain นั้น เหมือนกับ
   `vpc-fundamentals`/`hybrid-cross-vpc-connectivity` ที่ยอมรับไว้แล้วตั้งแต่ round 1 ว่า
   chapter-count ตรงสัดส่วนน้ำหนัก แต่ task-serving ข้ามโดเมนได้

</details>
</details>

## AWS-S1 — explanation + raised check ceiling

**สิ่งที่ทำ**: แก้สอง schema blocker ที่ปักไว้ในหัวข้อ "3 schema change ที่ต้องทำก่อน" ด้านบน (ข้อ 1
และ 2) end-to-end — content file → domain → MySQL → API → หน้า quiz — โดยไม่แตะ multiple-response
(AWS-S2) และไม่แยก concept ใดใน `aws.json`

### สิ่งที่ deliver

- Migration `migrations/008_recall-check-explanation.sql` — adds `recall_checks.explanation TEXT
  NULL`, guarded via `information_schema.COLUMNS` + `PREPARE`/`EXECUTE` (not a bare `ALTER TABLE`)
  so a restart after a crash between the DDL and the `schema_migrations` INSERT converges instead
  of looping (see "Round 2" below — this guard is a code-review fix, not part of the original cut)
- `internal/platform/mysql/migrate.go`'s `applyOne` now runs one migration's statements over a
  single acquired `*sql.Conn` instead of the ambient pool, so the guard's `SET`/`PREPARE`/`EXECUTE`
  sequence can't land on different pooled connections mid-migration (also a Round 2 fix)
- Domain: `RecallCheck` มีฟิลด์ `explanation` + accessor `Explanation()`, `NewRecallCheck` รับ
  parameter ใหม่, `maxRecallChecks` 5 → 15, `validateRecallExplanation` (bound ใหม่)
- Infra: `recallCheckFileDTO` มีฟิลด์ `explanation,omitempty`, `lessonwriter.go`/`lessonreader.go`
  thread ค่าผ่าน INSERT/SELECT ใหม่, API `recallCheckDTO` มีฟิลด์ `explanation,omitempty`
- Frontend: `RecallCheck.explanation?: string` ใน `web/lib/api.ts`, `RecallCheckCard.tsx` render
  ใน stage 3 (reveal) เท่านั้น
- `docs/tickets/mcq-quality.md` บันทึกตรง ๆ ว่าไม่มี measurement script commit ไว้ในโปรเจกต์

### การตัดสินใจ (พร้อมเหตุผล)

- **`maxRecallChecks` 5 → 15**: แผน AWS ตอนนั้นต้องการ ~11 ข้อ/concept (docs/tickets/aws-cert.md)
  เพดาน 15 เผื่อ headroom เหนือเป้านั้นโดยไม่เปิดให้ quiz ต่อ lesson ยาวไม่จำกัด — **เป้าย้ายเป็น 4–6
  ข้อ/concept ใน AWS-C0 ภายหลัง** (ดูหัวข้อ "รูปทรงใหม่" ต้นเอกสาร) แต่เพดาน 15 ยังเผื่อ headroom
  พอเหมือนเดิม ไม่ต้องแก้
- **`minRecallChecks` คงที่ 3 ไม่แก้**: blocker เดิมคือเพดานบน (5×50=250 ไม่พอ 550) ไม่ใช่พื้นล่าง —
  **แก้ไขหลัง code review รอบ 2**: ข้อความรอบแรกที่นี่อ้างว่า "118 lesson เดิม บางใบมีพอดี 3 ข้อ"
  ผิด — นับจริงจาก `content/lessons/*/*.json` ทั้ง 118 ไฟล์ ได้ distribution {4 checks: 4, 5 checks:
  114} **ต่ำสุดคือ 4 ไม่มีใบไหนมี 3 เลย** พื้นนี้จึงไม่เคย binding กับเนื้อหาจริงตั้งแต่ต้น
  การตัดสินใจ (คงพื้นไว้ที่ 3) ยังถูกต้องเหมือนเดิม แค่เหตุผลเดิมที่อ้างหลักฐานผิดต้องแก้ — เหตุผลจริง
  คือไม่มีอะไรใน scope ของ ticket นี้เรียกร้องให้เปลี่ยน contract ของพื้นนี้ ไม่ใช่เพราะมีใบไหนอยู่ติด
  ขอบ 3 พอดี
- **Length bound ของ explanation = 4000 runes**: ใกล้เคียง `maxOutlineRunes` (โครงสร้างระดับ
  "ย่อหน้าหลายส่วน" เหมือนกัน) มากกว่า `maxRecallAnswerRunes` (2000, คำตอบสั้นข้อเดียว) แต่ยังห่างจาก
  `maxBodyMdRunes` (50000, เนื้อหาบทเรียนเต็ม) มาก — เนื้อหาจริงตามสเปค (3 ส่วน หลายประโยค) น่าจะอยู่
  ราว 500–1500 ตัวอักษร ตัวเลข 4000 จึงเผื่อเกิน ~3 เท่าโดยยังจับค่าที่หลุดขอบเขตจริง ๆ ได้ (เช่น
  paste เอกสารทั้งฉบับผิดช่อง)
- **Empty string == absent (ไม่ error)**: `recallCheckFileDTO.Explanation` เป็น `string` ธรรมดา
  (ไม่ใช่ `*string`) ตาม convention เดิมของไฟล์นี้ (`Q`, `ExpectedAnswer` ก็เป็น `string` เดี่ยว) — ผล
  คือ JSON ที่ไม่มี key `explanation` เลย กับ JSON ที่มี `"explanation": ""` **decode ออกมาเหมือนกัน
  ทุกประการ** (ทั้งคู่ได้ `""`) ดังนั้นการแยกสองเคสนี้เป็นไปไม่ได้อยู่แล้วในระดับ decode เว้นแต่เปลี่ยน
  type เป็น pointer ซึ่งเกินความจำเป็น — `validateRecallExplanation` จึง treat `""` (หลัง trim) เป็น
  "ไม่มี explanation" เสมอ ไม่ error เด็ดขาด เพราะ 118 lesson เดิมทุกใบไม่มี key นี้และต้อง import
  ผ่านเหมือนเดิม
- **DB representation: `NULL` ไม่ใช่ `''`**: ใช้ convention เดียวกับคอลัมน์ `options` ที่มีอยู่แล้วใน
  ตารางเดียวกัน (ค่า optional → SQL NULL)
- **Explanation render ที่ไหน**: อยู่ใน stage 3 (reveal) เท่านั้น — ก้อนเดียวกับที่ `expected_answer`
  โผล่ (กฎเดียวกันทุกประการ: ต้องไม่ถึง DOM ก่อน stage 3) วางไว้หลัง breakdown ตัวเลือก (mcq) /
  แถบคำตอบ ก่อนปุ่ม Pass/Not yet (short_answer) ใช้ label "Why" (text-xs uppercase text-faint)
  + ข้อความ (font-thai text-sm leading-[1.8] text-body, `whitespace-pre-wrap` เผื่อผู้เขียนคั่น 3
  ส่วนด้วยบรรทัดว่าง) ในกล่อง `border-subtle bg-page rounded-xl` — **ไม่มีสีใหม่เลย** ใช้ token เดิม
  ทั้งหมด (`text-faint`, `text-body`, `bg-page`, `border-subtle`) ที่ verify contrast ไว้แล้วใน
  frontend-design skill

### Mutation table (13 mutation หลัง round 2, ทุกตัว revert กลับหลัง confirm แล้ว)

| # | Mutation | Killed by |
|---|---|---|
| 1 | explanation หายที่ hop file→domain (`lessonfile.go`'s `toDomain` ส่ง `""` แทน `rc.Explanation`) | `TestLoadLesson_ManyChecksWithExplanation` |
| 2 | explanation หายที่ hop domain→DB เขียน (`lessonwriter.go`'s `insertRecallCheck` ไม่ set ค่า) | `TestRepositorySaveLesson_RecallChecksReplacedInOrder`, `TestRepositorySaveLesson_RoundTripManyChecksExplanationNotMixedUp` |
| 3 | explanation หายที่ hop DB→domain อ่าน (`lessonreader.go`'s `toRecallCheck` รับ `""` แทน `explanation.String`) | `TestRepositoryLessonByConcept_HappyPath`, `TestRepositorySaveLesson_RoundTripManyChecksExplanationNotMixedUp` |
| 4 | explanation หายที่ hop domain→API DTO (`handler.go`'s `toRecallCheckDTOs` ไม่ set `Explanation`) | `TestHandlerGetLesson_Success` |
| 5 | explanation หายที่ hop API→UI (`RecallCheckCard.tsx` ไม่ render `check.explanation` เลย) | 3 เคสใน `describe("RecallCheckCard — explanation (AWS-S1)")` |
| 6 | explanation render ก่อน stage 3 (ย้ายบล็อกออกไปนอก `stage === "reveal"`) | เคสเดียวกับ #5 ("never puts explanation... before stage 3") |
| 7 | optionality หลุด (`validateRecallExplanation` reject string ว่าง) | `TestNewRecallCheck` (เคส valid, whitespace-only) + `TestLoadLesson_Valid` (fixture จริงไม่มี explanation) |
| 8 | `maxRecallChecks` revert กลับ 5 | `TestNewLesson` ("fifteen recall checks"), `TestLoadLesson_ManyChecksWithExplanation`, `TestRepositorySaveLesson_RoundTripManyChecksExplanationNotMixedUp` |
| 9 | `minRecallChecks` bound หลุด | `TestNewLesson` ("two recall checks"), `TestLoadLesson_Errors` ("two recall checks") |
| 10 | length bound หลุด | `TestNewRecallCheck` ("explanation exceeding max runes") |
| 11 | explanation ผูกผิด check (index mixup ที่ file→domain, ใช้ `d.RecallChecks[0].Explanation` แทน `rc.Explanation` ทุกตัว) | `TestLoadLesson_ManyChecksWithExplanation` (fixture 8 ข้อ, เกิน 5 เดิม) |
| 12 | (round 2, code review S1) `RecallCheckCard.tsx`'s `check.explanation &&` → `check.explanation !== undefined &&` — เดิม**รอด**เพราะไม่มีเทสต์ไหน supply `explanation: ""` | `"renders no explanation section when explanation is an empty string, not just when it is absent"` (เทสต์ใหม่) |
| 13 | (round 2, code review S1) `recallCheckDTO.Explanation`'s JSON tag `omitempty` ถูกลบ — เดิม**รอด**เพราะเทสต์เดิม decode เป็น struct ก่อนเทียบ ซึ่ง `""` จาก key หาย กับ `""` จาก key ว่างเปล่า decode ออกมาเหมือนกันเป๊ะ | `TestHandlerGetLesson_OmitsExplanationKeyWhenAbsent` (เทสต์ใหม่ เช็ค raw response body string) |

**Survivor**: ไม่มี — ทุก mutation ตายตามที่คาด (รวม 13 มิวเทชันหลัง round 2 คือ 11 ตัวเดิม + 2
ตัวที่ code review พบว่ารอดรอบแรกแล้วมีเทสต์ใหม่มา kill)

**หมายเหตุ scope ของ #11**: mutate จริงทำที่ hop file→domain เท่านั้น (ที่เดียวที่ explanation กับ
question ถูก zip มาจาก array คนละตัว/index) เพราะ hop อื่น (DB write/read, API DTO) ส่งต่อ
`RecallCheck` เป็น value เดียวที่ explanation เป็นฟิลด์ติดอยู่กับ struct เสมอ ไม่มีทางแยก array ให้
mix up ได้ในเชิงโครงสร้าง — `TestRepositorySaveLesson_RoundTripManyChecksExplanationNotMixedUp`
(DB round-trip 8 ข้อ) และ `TestNewLessonSortKeepsExplanationWithItsOwnCheck` (domain sort) ยังคง
เป็น regression net แต่ไม่ได้ผ่านการ mutate จริงเพราะไม่มี mutation ที่สมจริงจะ decouple ได้ที่ hop
นั้น

### Migration re-run evidence (R1, code review round 2)

Reproduce ทั้ง before และ after บน stack แยก `docker compose -p aws1fix` (ไม่แตะ
`self-learning_mysql_data`, ปิดท้ายด้วย `docker compose down` เปล่า ๆ ไม่มี `-v`):

- **ก่อนแก้ (bare `ALTER TABLE`, จำลอง crash)**: สร้างตารางแยก, รัน `ALTER TABLE ... ADD COLUMN
  explanation` สำเร็จ, รันซ้ำจำลอง "restart หลัง crash ก่อนบันทึก version row" → **`ERROR 1060
  (42S21): Duplicate column name 'explanation'`** ตรงตามที่ reviewer อธิบาย ยืนยันว่า bug จริง
  ไม่ใช่ทฤษฎี
- **จำลอง crash จริงกับ API binary จริง**: boot stack เต็ม (mysql + api) ครั้งแรก → migration 8
  ผ่าน, `schema_migrations` มี version 1-8, คอลัมน์ `explanation` มีจริง → `DELETE FROM
  schema_migrations WHERE version = 8` (จำลอง "ALTER สำเร็จแต่ crash ก่อนบันทึก version") →
  `docker compose restart api`
- **หลังแก้ (guarded ALTER + `applyOne` ใช้ connection เดียว)**: container **ไม่ crash-loop**
  (`RestartCount = 0`, health check ผ่านต่อเนื่องหลัง restart), log สะอาดไม่มี error 1060,
  `schema_migrations` มี version 8 กลับมาเหมือนเดิม (migration รันซ้ำแล้ว no-op ผ่านการ guard),
  คอลัมน์ `explanation` ยังอยู่ครบ ไม่มีการ error หรือ column ซ้ำ
- **Static guard**: `TestAllMigrationsNonIdempotentAlterIsGuarded` ใน
  `internal/platform/mysql/migrate_embedded_test.go` — **แก้ 2 รอบแล้ว** รอบแรก
  (`TestAllMigrationsAlterTableAddColumnIsGuarded`) ทำ whole-file `strings.Contains` 4 จุด และ
  claim ผิดว่า "เป็น static-text limit แบบเดียวกับเทสต์ CREATE TABLE" — **claim นั้นเท็จ**:
  `TestAllMigrationsCreateTableIsIdempotent` เดิมเดินทีละ occurrence ของ `CREATE TABLE` จริง
  (positional, ทุกจุด) ส่วน whole-file Contains ไม่ใช่แบบเดียวกันเลย และ reviewer เขียนไฟล์ probe 5
  แบบที่ผ่านเทสต์เดิมหมด: (A) ALTER guarded ถูกต้อง + bare ALTER ต่อท้ายในไฟล์เดียวกัน, (B) SQL
  lowercase ไม่มี guard, (C) bare ALTER + comment ที่ quote คำ marker ของ guard ไว้ข้าง ๆ, (D) bare
  `CREATE INDEX`, (E) bare `ALTER TABLE ... DROP COLUMN` — **รอบ 3 เขียนใหม่เป็น per-statement**:
  parse ไฟล์ผ่าน `splitStatements` ตัวเดียวกับที่ `applyOne` ใช้จริง (production path) หลังผ่าน
  `stripLineComments` ใหม่ (ตัด `-- ...` ทิ้งก่อน split กัน probe C), แล้วเช็คทุก statement ว่าไม่มี
  ตัวไหนขึ้นต้นด้วย `ALTER TABLE` ที่มี `ADD COLUMN`/`DROP COLUMN`, หรือขึ้นต้นด้วย `CREATE INDEX`
  แบบ bare — **verify ซ้ำทั้ง 5 probe shape ของ reviewer + 2 shape ที่คิดเพิ่มเอง (multi-line ALTER
  แยกบรรทัด, bare ALTER + inline trailing comment บนบรรทัดเดียวกัน) ตายครบทั้ง 7** เทสต์นี้พิสูจน์แค่
  ว่า production parser (`splitStatements`) เห็น statement นั้นเป็น bare non-idempotent DDL จริง
  (**ไม่ใช่**เทสต์แบบเดียวกับ CREATE TABLE เป๊ะ ๆ — ยอมรับตรง ๆ ว่าเป็นคนละเครื่องมือ ไม่ claim parity
  อีกต่อไป) ส่วน**พฤติกรรม**จริงของ MySQL พิสูจน์ด้วย live re-run ข้างบนแทน ไม่ใช่ `go test`
- **เหตุผลที่ต้องแก้ `applyOne` ด้วย ไม่ใช่แค่ SQL**: `SET @var`/`PREPARE`/`EXECUTE` เป็น
  session-scoped state ของ MySQL connection เดียว — `database/sql` ไม่การันตีว่า
  `db.ExecContext` สองครั้งติดกันจะได้ connection เดิมจาก pool เสมอ (ทดสอบ stress จริงด้วย 9
  goroutine ยิง query พร้อมกันระหว่าง sequence 5 statement ไม่พบการสลับ connection เลยใน 15 รอบ
  แต่นั่นเป็นพฤติกรรมที่ไม่มีสัญญาเป็นเอกสาร ไม่ใช่ contract ที่พึ่งได้) `applyOne` จึงเปลี่ยนไปใช้
  `db.Conn(ctx)` ตัวเดียวตลอดทั้ง migration file แทนที่จะพึ่งพฤติกรรมที่ verify แต่ไม่รับประกัน

### Import-parity evidence

Stack แยก `docker compose -p aws1check` (ไม่แตะ `self-learning_mysql_data`) เทียบกับ baseline
`docker compose -p aws1baseline` ที่ build จาก `develop` (`git worktree`) — รัน `import-curriculum`
+ `import-lessons` จริงทั้งคู่:

- จำนวนไฟล์ import: **118 ไฟล์** ทั้งสองฝั่ง (DDIA 61 + AI-systems 53 + DDD ch.1 4 lesson)
- `lessons` table: **118 แถวเท่ากันทั้งสองฝั่ง**
- `recall_checks` table: **586 แถวเท่ากันทั้งสองฝั่ง**, `explanation IS NOT NULL` = **0** (ยังไม่มี
  content file ใดใช้ field นี้จริง — ตามสโคปที่ตั้งใจ)
- MD5 hash ของ `(id, title_en, body_md, refs)` ทุกแถว lessons และของ
  `(lesson_id, position, type, question, expected_answer, options)` ทุกแถว recall_checks
  **เท่ากันตัวต่อตัว** ระหว่าง develop กับ branch นี้ — ยืนยัน byte-identical import จริง ไม่ใช่แค่
  นับจำนวนแถว

### Playwright evidence (headless Chromium, stack `aws1check` เดิม)

ใช้ lesson `domain-driven-design/ubiquitous-language` (5 recall_checks จริงในฐานข้อมูล) — UPDATE
`explanation` ตรงให้ check #2 (mcq, ข้อความจริง 3 ส่วนยาว 752 ตัวอักษร) และ #5 (mcq, 300 ตัวอักษร)
ผ่าน SQL โดยตรง (ไม่มี content file ไหนใช้ field นี้จริงตามสโคป):

- Stage 1 (recall) และ stage 2 (commit) ของ serialized `innerHTML`: **ไม่มี**ข้อความ explanation
  ปนอยู่เลย — `explanation_in_stage1 = false`, `explanation_in_stage2 = false`
- Stage 3 (reveal): explanation โผล่จริง พร้อม label "Why" — `explanation_in_stage3 = true`,
  `why_label_present = true`
- Check ที่ไม่มี explanation (position 1) render เหมือนเดิมทุกประการที่ stage 3: **ไม่มี** label
  "Why" โผล่ (`no_explanation_card_has_why_label = false`)
- Overflow ที่ viewport 375px: วัด 3 ระดับ ไม่มีจุดไหนล้น — explanation `<p>` เอง
  `scrollWidth = clientWidth = 251`, การ์ดทั้งใบ `scrollWidth = clientWidth = 325`, ระดับ document
  `scrollWidth = clientWidth = 375` (ตรงกับเงื่อนไขที่ ticket ระบุเป๊ะ)
- Contrast วัดจริงจาก computed style: label ("Why", `text-faint` บน `bg-page`) = **4.52:1** (ผ่าน AA
  4.5:1 พอดี ตรงกับตัวเลขที่ frontend-design skill เคย verify ไว้), body text (`text-body` บน
  `bg-page`) = **13.94:1** — **ไม่มีสีใหม่เลยในทั้งสองกรณี** ใช้ token เดิมของระบบ
- `console` error และ `pageerror` = **0 ทั้งคู่** ตลอดการทดสอบ

### Test count

- Go: `go vet ./...` clean, `gofmt -l .` clean, `go test -count=1 ./...` **296 (develop) → 299
  (round 1) → 301 (round 2) → 301 (round 3)** (นับจาก `go test -v` ผ่าน `grep -c "^--- PASS"`;
  round 2 เพิ่ม 2 เทสต์ใหม่; round 3 **แทนที่** `TestAllMigrationsAlterTableAddColumnIsGuarded` ด้วย
  `TestAllMigrationsNonIdempotentAlterIsGuarded` แบบ 1 ต่อ 1 — จำนวนสุทธิเท่าเดิมเพราะเป็นการเขียน
  เทสต์เดิมใหม่ ไม่ใช่เพิ่มเทสต์คู่ขนาน)
- Frontend: `npx tsc --noEmit` clean, `npx eslint .` clean, `npm run build` clean,
  `npx vitest run` **169 (develop) → 175 (round 1) → 176 (round 2) → 176 (round 3, ไม่แตะ
  frontend)**

### สิ่งที่ตั้งใจไม่ทำในรอบนี้

- ไม่เขียน explanation ให้ MCQ 356 ข้อเดิม (Q-3's เนื้อหาส่วนที่เหลือ ยังพักตามแผน)
- ไม่สร้าง AWS-S3 (guessability measurement tool) เอง — spec เขียนไว้แล้วในหัวข้อ "AWS-S3" ด้านบน
  ตามที่ reviewer ขอ รอ ticket แยก
- ไม่แตะ multiple-response (AWS-S2) และไม่แยก concept ใน `aws.json`
- ไม่ refresh `docs/design.md` ทั้งไฟล์ (เพิ่มแค่ 1 บรรทัดในหนี้ที่รู้ตัวเดิมของ roadmap.md ว่า schema
  block ขาด `explanation` — ตามที่ code review บอกว่าไม่ต้องทำ refresh เต็มในรอบนี้)

### Round 2 (code review) — สรุปสิ่งที่แก้

REQUEST_CHANGES รอบแรกพบ 7 ปัญหาหลัก (R1–R7) + 3 ปัญหาเสริม (S1, S2, S4) แก้ครบทุกข้อ:

- **R1 (แก้ที่ใหญ่ที่สุด)**: migration 008 เดิมเป็น bare `ALTER TABLE ADD COLUMN` ซึ่งไม่มี
  `IF NOT EXISTS` ใน MySQL 8 — restart หลัง crash ระหว่าง DDL กับการบันทึก version row จะวน error
  1060 ตลอดไป (compose ตั้ง `restart: unless-stopped`) แก้ด้วย guard ผ่าน
  `information_schema.COLUMNS` + `PREPARE`/`EXECUTE`, และแก้ `applyOne` ให้ใช้ connection เดียว
  ตลอด migration file (เหตุผลเต็มอยู่ที่ "Migration re-run evidence" ด้านบน) — พิสูจน์ด้วย live
  crash-recovery replay จริงบน `docker compose -p aws1fix`
- **R2**: comment ของ `minRecallChecks` อ้างหลักฐานผิด ("118 lesson บางใบมีพอดี 3 ข้อ") นับจริงแล้ว
  ต่ำสุดคือ 4 ไม่มีใบไหนมี 3 — แก้ comment ให้ตรงข้อเท็จจริง (การตัดสินใจเดิมยังถูกต้อง)
- **R3**: ตัวเลข "356" ผิดที่ 2 จุด (migration comment, review-focus เฉลยข้อ 3) ที่จริงคือ 586
  (ทุกแถว ไม่ใช่แค่ mcq) และ "114 lesson" ผิดที่ 2 จุด (text.go comment, vitest test name) ที่จริง
  คือ 118 — แก้ครบทุกจุด
- **R4**: ลบ WHAT-comment บน `Explanation()` ที่ restate signature ซ้ำกับ struct doc
- **R5**: `.claude/agents/lesson-writer.md`/`lesson-verifier.md` ยังไม่รู้จัก `explanation` และยัง
  ระบุ "3-5 items" — แก้ schema ให้มี `explanation` (REQUIRED เฉพาะ AWS track), เพิ่ม item count
  เป็น ~11 สำหรับ AWS track, เพิ่ม verifier gate ข้อ 8 ที่ FAIL ถ้า AWS-track question ไม่มี
  3-part explanation ครบ
- **R6**: บันทึกตรง ๆ ว่า guessability baseline 33.7%/33.4%/89.6% ใน roadmap.md **รันซ้ำไม่ได้**
  (ไม่มี script), 25% gate ของ AWS corpus **unenforceable จนกว่า AWS-S3 จะมี**, เขียน spec ขั้นต่ำ
  ของ AWS-S3 ไว้ใน aws-cert.md
- **R7**: ชื่อ ticket ชนกัน ("AWS-1" หมายถึงทั้ง ticket นี้และ placeholder เดิมของคลังข้อสอบ) — เปลี่ยน
  เป็น AWS-S1 (ticket นี้) / AWS-S2 (multiple-response) / AWS-S3 (measurement tool) / AWS-C1..C4
  (คลังข้อสอบต่อ domain) ทั้งใน aws-cert.md และ roadmap.md — **ไม่เปลี่ยนชื่อ branch**
  (`ticket/aws-1-explanation` เกิดก่อนเปลี่ยนชื่อ ตามที่ reviewer สั่งให้บันทึกไว้แทนการ force-push)
- **S1**: `check.explanation &&` ใน `RecallCheckCard.tsx` เปลี่ยนเป็น `!== undefined` แล้ว**รอด**
  ทุกเทสต์เดิม (ไม่มีเทสต์ไหน supply `explanation: ""`) — เพิ่มเทสต์ frontend ใหม่ + เทสต์ backend
  ใหม่ที่เช็ค raw JSON body ไม่มี substring `"explanation"` เมื่อไม่มีค่า (ดู mutation #12, #13)
- **S2**: เพิ่ม `explanation` เข้า comment ของ `lessonDTO` ที่บันทึกการตัดสินใจส่งข้อมูลเฉลยให้ client
- **S4**: เปลี่ยนชื่อเทสต์ที่อ้างว่า "two cards mounted side by side" ทั้งที่จริง sequential
  (unmount การ์ดแรกก่อน mount การ์ดที่สอง)

### Status: implemented, code review round 2 fixes applied, PR pending

### Review focus (round 2 — เน้นจุดที่แก้ตาม code review)

<details>
<summary>คำถามสำหรับรีวิว diff รอบนี้ (เฉลยพับไว้ด้านล่าง)</summary>

1. ทำไมแค่แก้ SQL ของ migration 008 ให้มี guard อย่างเดียวไม่พอ ต้องแก้ `applyOne` ใน `migrate.go`
   ด้วย?
2. `TestAllMigrationsAlterTableAddColumnIsGuarded` พิสูจน์อะไร และ**ไม่**พิสูจน์อะไร?
3. ทำไม `minRecallChecks` ยังคงค่า 3 เหมือนเดิมหลัง round 2 ทั้งที่ comment ที่ใช้เป็นหลักฐานผิด?

<details>
<summary>เฉลย</summary>

1. เพราะ guard ใหม่ (`SET @var` → `PREPARE`/`EXECUTE`) พึ่ง MySQL user variable ซึ่งเป็น
   session-scoped state ของ connection เดียว ในขณะที่ `applyOne` เดิมเรียก `db.ExecContext` บน
   `*sql.DB` (pool) ตรง ๆ ทีละ statement — `database/sql` ไม่การันตีว่าสอง `Exec` ติดกันจะได้
   connection เดิมจาก pool เสมอ (ไม่มีสัญญาเอกสารข้อนี้) ถ้าเผอิญได้ connection คนละตัว
   `@add_explanation_ddl` จะเป็น NULL และ `PREPARE stmt FROM NULL` จะ error 1064 ทันที — แก้แค่ SQL
   แต่ไม่ปักหมุด connection คือแก้ crash-loop เดิม (R1) แต่เปิดความเสี่ยงใหม่ที่พิสูจน์ไม่ได้ว่าไม่เกิด
   จึงต้องแก้ `applyOne` ให้ใช้ `db.Conn(ctx)` ตัวเดียวตลอด statement ของ migration file นั้นด้วย
2. พิสูจน์แค่ว่าไฟล์ migration ที่มี `ADD COLUMN` มี **marker ของ guard pattern ครบ** (text-based,
   เหมือนเทสต์ `CREATE TABLE`/`IF NOT EXISTS` เดิม) — **ไม่**พิสูจน์ว่า guard **ทำงานถูกจริง** บน
   MySQL จริง (เช่น syntax ผิดเล็กน้อยที่ยังมี marker ครบแต่รันไม่ผ่านจริงจะไม่ถูกจับ) พฤติกรรมจริง
   พิสูจน์ด้วย live crash-recovery replay บน `docker compose -p aws1fix` แทน ไม่ใช่ `go test`
   **[แก้ไขใน round 3: คำตอบข้อนี้เองก็ผิด — "เหมือนเทสต์ CREATE TABLE เดิม" ไม่จริง เป็น whole-file
   Contains ที่ผ่านง่ายกว่ามาก reviewer เขียน probe 5 แบบผ่านหมด ดูหัวข้อ "Round 3" ด้านล่าง]**
3. เพราะ**ข้อสรุป**ของ decision (คงพื้นไว้ที่ 3) ยังถูกต้องอยู่ — สิ่งที่ผิดคือ**หลักฐาน**ที่ยกมาอ้าง
   (อ้างว่ามีใบที่มีพอดี 3 ข้อ ซึ่งนับจริงแล้วไม่มี ต่ำสุดคือ 4) ไม่ใช่ตัวการตัดสินใจเอง — ไม่มีอะไรใน
   scope ของ ticket นี้เรียกร้องให้เปลี่ยน floor จาก 3 เป็นค่าอื่น (blocker เดิมคือเพดานบนเท่านั้น)
   ดังนั้นแก้แค่ comment ให้ตรงข้อเท็จจริง ไม่ต้องเปลี่ยนค่า constant

</details>
</details>

## Round 3 (code review) — regression guard ตัวเองอ่อนกว่าที่ comment อ้าง

**NO-SHIP หนึ่งเดียว**: `TestAllMigrationsAlterTableAddColumnIsGuarded` (เขียนใน round 2) ทำ
whole-file `strings.Contains` 4 จุด และ**ตัวเทสต์เองที่มี comment claim ว่า "เป็น static-text limit
แบบเดียวกับเทสต์ CREATE TABLE"** — claim นั้นเท็จ (`TestAllMigrationsCreateTableIsIdempotent` เดินทีละ
occurrence จริง ไม่ใช่ whole-file membership) reviewer เขียนไฟล์ probe 5 แบบผ่านเทสต์เดิมหมด:

| Probe | รูปแบบ | ผ่านเทสต์เดิม (round 2) เพราะอะไร | ผลหลังแก้ (round 3) |
|---|---|---|---|
| A | ALTER guarded ถูกต้อง + bare `ALTER TABLE lessons ADD COLUMN b INT NULL;` ต่อท้ายในไฟล์เดียวกัน | whole-file Contains เจอ marker ครบจากส่วน guarded แล้วไม่มองต่อว่ามี statement เปล่าเพิ่มมา | **FAIL** ✅ |
| B | `alter table lessons add column b int null;` (lowercase, ไม่มี guard เลย) | `Contains(content, "ADD COLUMN")` (ตัวพิมพ์ใหญ่) เป็น false เทสต์เลย early-return | **FAIL** ✅ |
| C | bare ALTER + `--` comment ที่ quote คำ `information_schema.COLUMNS`/`PREPARE`/`EXECUTE`/`DEALLOCATE PREPARE` ไว้ข้าง ๆ | whole-file Contains เจอคำพวกนี้ใน comment เฉย ๆ ไม่สนว่าเป็น comment หรือโค้ดจริง | **FAIL** ✅ |
| D | bare `CREATE INDEX idx_probe ON lessons (version);` | เทสต์เดิมเช็คแค่ `ADD COLUMN` ไม่รู้จัก `CREATE INDEX` เลย | **FAIL** ✅ |
| E | bare `ALTER TABLE lessons DROP COLUMN version;` | เทสต์เดิมเช็คแค่ `ADD COLUMN` ไม่รู้จัก `DROP COLUMN` เลย | **FAIL** ✅ |
| F (คิดเพิ่มเอง) | `ALTER TABLE`/`ADD`/`COLUMN` แยกคนละบรรทัด มี whitespace คั่น | ไม่เคย test กับเทสต์เดิม แต่ทดสอบแล้วก็จะรอดเหมือนกันเพราะ Contains ยังเจอ substring "ADD COLUMN" ปกติ (probe นี้พิสูจน์ฝั่ง**เทสต์ใหม่**ว่า normalize whitespace ถูกต้อง ไม่ได้พิสูจน์ว่าเทสต์เก่าพัง) | **FAIL** ✅ |
| G (คิดเพิ่มเอง) | bare ALTER + inline trailing `-- comment` บนบรรทัดเดียวกัน | เช่นเดียวกับ F | **FAIL** ✅ |

ทั้ง 7 probe (5 ของ reviewer + 2 ที่คิดเพิ่ม) ทดสอบจริงโดยสร้างไฟล์ `.sql` ชั่วคราวใต้ `migrations/`
รัน `go test`, ลบไฟล์ทิ้งทันที — `git status` สะอาดหลังทำเสร็จทุกรอบ

**แก้จริง**: เขียน `TestAllMigrationsAlterTableAddColumnIsGuarded` ใหม่ทั้งหมดเป็น
`TestAllMigrationsNonIdempotentAlterIsGuarded` — parse ไฟล์ผ่าน `splitStatements` **ตัวเดียวกับที่
`applyOne` ใช้รันจริง** (ไม่ใช่ parser แยกต่างหากที่พิสูจน์กันคนละหน่วย) หลังผ่าน `stripLineComments`
ใหม่ (ตัด `-- ...` ทิ้งก่อน split กัน probe C) แล้วเช็ค**ทุก statement เดี่ยว ๆ** ว่าไม่มีตัวไหนขึ้นต้น
`ALTER TABLE` ที่มี `ADD COLUMN`/`DROP COLUMN`, หรือขึ้นต้น `CREATE INDEX` แบบ bare — ครอบคลุมทั้ง 3
รูปแบบ DDL ที่ไม่มี `IF NOT EXISTS`/`IF EXISTS` form ใน MySQL 8 ตามกฎของ `migrate.go`'s doc comment
เอง ("every migration file must therefore be written to converge on re-run") ไม่ใช่แค่ `ADD COLUMN`

**ยืนยัน 8 migration ที่มีอยู่จริงทั้งหมดยังผ่านเทสต์ใหม่** (`go test -run
TestAllMigrationsNonIdempotentAlterIsGuarded` เขียว) — 003/004's `ALTER TABLE topics MODIFY COLUMN
...` ไม่ตรงเงื่อนไขไหนเลย (ไม่ใช่ ADD/DROP COLUMN หรือ CREATE INDEX) จึงผ่านถูกต้อง, 008's guarded
ADD COLUMN อยู่ใน string literal ของ statement `SET @add_explanation_ddl = IF(...)` ซึ่งขึ้นต้นด้วย
`SET` ไม่ใช่ `ALTER TABLE` จึงไม่ถูกจับผิด

**S8**: เพิ่ม 1 บรรทัดใน `applyOne`'s doc comment — `go-sql-driver/mysql` (เช็คจริงที่ v1.10.0)'s
`ResetSession` ทำแค่ liveness check ไม่ส่ง `COM_RESET_CONNECTION` ดังนั้น user variable ของ MySQL
(`@has_explanation_col`/`@add_explanation_ddl`) รอดอยู่บน connection ที่ pool คืนกลับไปได้แม้เรียก
`conn.Close()` แล้ว — ไม่มีปัญหาวันนี้เพราะ 008 `SET` ค่าเองก่อนอ่านเสมอ แต่ migration guarded ตัวถัดไป
ในอนาคตที่ก๊อปปี้ pattern นี้แล้ว**อ่าน**ตัวแปรโดยไม่ `SET` เองก่อน มีความเสี่ยงรับค่าเก่าที่ค้างจาก
migration อื่นแบบเงียบ ๆ — เขียนกฎนี้ไว้เป็น doc comment ให้คนเขียน migration ถัดไปเห็น

### Test count (round 3)

- Go: **301 → 301** (แทนที่เทสต์เดิม 1 ตัวด้วยเทสต์ใหม่ 1 ตัว ไม่มีการเพิ่มจำนวนสุทธิ)
- Frontend: **176 → 176** (ไม่แตะ frontend ในรอบนี้)

### Status: implemented, code review round 3 fixes applied, PR pending

### Review focus (round 3)

<details>
<summary>คำถามสำหรับรีวิว diff รอบนี้ (เฉลยพับไว้ด้านล่าง)</summary>

1. ทำไม probe C (bare ALTER + comment ที่ quote คำ marker) ถึงหลอกเทสต์ round 2 ได้ แต่หลอกเทสต์
   round 3 ไม่ได้?
2. ทำไม `stripLineComments` ต้องทำงาน**ก่อน** `splitStatements` ไม่ใช่หลัง?
3. ทำไมเทสต์ใหม่ต้องเช็ค `DROP COLUMN` กับ `CREATE INDEX` ด้วย ทั้งที่ ticket นี้ไม่มี migration ไหน
   ใช้สองอย่างนี้เลย?

<details>
<summary>เฉลย</summary>

1. เพราะเทสต์ round 2 เช็คแค่ "คำเหล่านี้ปรากฏอยู่ที่ไหนสักแห่งในไฟล์ทั้งไฟล์หรือเปล่า"
   (`strings.Contains(content, ...)`) ไม่สนว่าคำนั้นอยู่ใน comment, string literal, หรือ statement
   จริง — comment ที่พิมพ์คำ marker ตรง ๆ จึงทำให้เช็คผ่านได้ง่าย ๆ โดยไม่มี guard จริงเลย เทสต์ round
   3 ไม่เช็คคำที่ปรากฏในไฟล์อีกต่อไป แต่ parse ไฟล์เป็น statement ก่อน (ผ่าน `stripLineComments` +
   `splitStatements`) แล้วเช็คแต่ละ statement เป็นหน่วย — comment ถูกตัดทิ้งไปตั้งแต่ก่อน parse จึงไม่มี
   ทางเข้าไปปนกับเนื้อหา statement จริงได้อีก
2. เพราะ `splitStatements` แบ่งไฟล์ด้วย `;` และ comment ในไฟล์นี้ (สไตล์ `-- ...`) ไม่ได้ถูก parse
   เป็น token พิเศษโดย `splitStatements` เอง (มันไม่รู้จัก syntax comment เลย เห็นแค่ตัวอักษร) — ถ้า
   split ก่อนแล้วค่อยตัด comment ทีหลัง, comment ที่บังเอิญมี `;` อยู่ข้างในจะทำให้ statement นับผิดไป
   แล้ว (ตัด comment ในสิ่งที่คิดว่าเป็น "statement" หนึ่งอันซึ่งจริง ๆ ถูกตัดครึ่งไปแล้วจาก `;` ใน
   comment) การตัด comment ก่อนจึงรับประกันว่า `splitStatements` เห็นแต่ SQL จริงล้วน ๆ ไม่มี
   comment ปนมาทำให้แบ่ง statement ผิดที่
3. เพราะกฎที่ `migrate.go`'s doc comment เขียนไว้เองคือ "**every** migration file must ... converge
   on re-run" ไม่ได้จำกัดแค่ `ADD COLUMN` — `DROP COLUMN` และ `CREATE INDEX` มีข้อจำกัดเดียวกันเป๊ะ ๆ
   ใน MySQL 8 (ไม่มี `IF EXISTS`/`IF NOT EXISTS` form) เขียนเทสต์ให้ครอบคลุมกฎที่ตั้งไว้เองทั้งหมด
   ตั้งแต่ตอนนี้ถูกกว่าการรอให้มี migration จริงมาเจอปัญหาเดียวกันซ้ำแล้วค่อยแก้เทสต์เพิ่มทีหลัง

</details>
</details>

## AWS-S3 — guessability measurement tool

**สิ่งที่ทำ**: สร้าง `cmd/mcq-guessability/internal/mcqguess` (pure logic, ไม่แตะ DB/HTTP, unit-tested
หนัก) + `cmd/mcq-guessability` (thin CLI wrapper ตาม convention เดียวกับ `cmd/import-lessons`/
`cmd/import-curriculum`) ที่วัด guessability ของ `recall_checks` ประเภท mcq จาก
`content/lessons/*/*.json` โดยตรง (ไม่ผ่าน curriculum domain) — baseline คำนวณต่อข้อจาก
`1/len(options)` เสมอ ไม่ hardcode ทั้งเก่า (kept) ไม่แตะ Docker/DB (ตาม ticket ระบุ ไม่ต้องมี DB)

**Round 2 (code review) เพิ่ม**: gate ต่อ track แทน pooled (`EvaluateTrackGates`), ขั้นต่ำจำนวนตัวอย่าง
ก่อน gate ตัดสิน (`MinSampleSize`), ขั้นต่ำ 2 ตัวเลือกต่อคำถาม (`minMeasurableOptions`) — รายละเอียด
เหตุผลทั้งหมดอยู่ที่หัวข้อ "Round 2 (code review)" ท้ายเอกสารนี้

**Round 3 (code review) เพิ่ม**: shortest/middle-option heuristic (`Report.Shortest`/`Middle`,
`shortestOptionIndex`/`middleOptionIndex`), `GateResult.Requested`/`Judged` (แยก "ตัดสินแล้วผ่าน" กับ
"ไม่ได้ตัดสินอะไรเลย"), ย้าย package ไป `cmd/mcq-guessability/internal/mcqguess` จริง — รายละเอียด
เหตุผลทั้งหมดอยู่ที่หัวข้อ "Round 3 (code review)" ท้ายเอกสารนี้

### สิ่งที่ deliver

- `.../internal/mcqguess/question.go` — `Question` (Track/Source/Options/ExpectedAnswer),
  `Baseline()` (`1/len(options)`, guard คืน 0 แทน +Inf ถ้า options ว่าง), `ExpectedIndex()`
- `.../internal/mcqguess/heuristics.go` — `longestOptionIndex`/`shortestOptionIndex`/
  `middleOptionIndex` (unexported, rune-based, deterministic tie-break), `HeuristicResult`
  (Hits/Total/baselineSum → `HitRate()`/`AvgBaseline()`/`ExcessRatio()`)
- `.../internal/mcqguess/report.go` — `Report` (per-track หรือ overall, มี `Longest`/`Shortest`/
  `Middle`/`Position`), `Measure(questions)` คืน overall + `map[string]Report` ต่อ track, abort
  ทันทีถ้า mcq มี options น้อยกว่า 2 (`minMeasurableOptions`) หรือ `expected_answer` ไม่อยู่ใน options
  ของตัวเอง
- `.../internal/mcqguess/gate.go` — `EvaluateGate(report, maxExcessRatio)` (per-report, ข้าม
  heuristic ที่ n < `MinSampleSize` ไปเป็น `Insufficient` แทนที่จะตัดสิน, นับ `Judged`),
  `EvaluateTrackGates(byTrack, maxExcessRatio)` (fail ถ้า track ใด track หนึ่ง fail), `ExitCode(gate)`
  (คืน 1 ถ้า requested แล้ว `Failed` หรือ `Judged==0`)
- `.../internal/mcqguess/load.go` — `LoadQuestions(root)` เดิน `filepath.WalkDir` หา `*.json` ทุกไฟล์
  ใต้ root แบบ recursive, decode DTO ของตัวเอง (ไม่ใช้ `curriculuminfra.LoadLesson`), filter เฉพาะ
  `type == "mcq"`
- `cmd/mcq-guessability/main.go` — flag `-dir` (default `content/lessons`), `-max-excess` (default
  `-1` = ปิด gate), พิมพ์รายงานต่อ track + overall (context only, ไม่ถูก gate) แล้วจบด้วยผล gate
  ต่อ track (`PASS`/`FAIL`/`NOT JUDGED`) + รายการ heuristic ที่ n ไม่พอ
- เทสต์: `.../internal/mcqguess/{question,heuristics,report,gate,load}_test.go` (31 top-level
  test) + `cmd/mcq-guessability/main_test.go` (5 top-level test, end-to-end ผ่าน `run()` จริง)

### การตัดสินใจ (พร้อมเหตุผล)

- **Tie-break ของ longest-option heuristic**: เลือก index **ต่ำสุด**ในกลุ่มที่ยาวเท่ากัน
  (`l > bestLen` ไม่ใช่ `>=`, ดังนั้นค่าที่เจอก่อนชนะเสมอ) เหตุผลสองข้อ: (1) จำลอง test-taker จริงที่
  กวาดสายตาบนลงล่างแล้วหยุดที่ตัวแรกที่ "ดูยาวที่สุด" ไม่ใช่สุ่มเลือกในกลุ่มที่เสมอกัน (2) deterministic
  — ไม่มี RNG เลยทั้ง tool นี้ ผลลัพธ์จึงรันซ้ำได้ byte-ต่อ-byte ทุกครั้ง ซึ่งจำเป็นถ้าจะใช้เป็น gate ใน
  CI ในอนาคต (RNG ในเกทจะทำให้ build เขียว/แดงสลับกันได้โดยไม่มีอะไรเปลี่ยนจริง)
- **Excess metric = ratio `(actual-baseline)/baseline` ไม่ใช่ percentage points สัมบูรณ์**:
  percentage points ไม่ portable ข้าม corpus ที่ baseline ต่างกัน — 5pp เหนือ baseline 25% คือ
  relative jump 20%, แต่ 5pp เดียวกันเหนือ baseline 33.3% คือแค่ 15% — ใช้ pp ตรง ๆ จะสร้างบั๊กแบบ
  เดียวกับที่ ticket นี้มีไว้แก้ (raw hit rate ที่เทียบข้าม corpus กันไม่ได้) ขึ้นมาใหม่อีกชั้นหนึ่งที่ตัว
  metric ของ excess เอง แทนที่จะแก้ที่ตัว baseline
- **Gate เช็คทั้ง length heuristic และทุก position-index heuristic** ไม่ใช่แค่ length: แม้ Q-1
  (`web/lib/shuffle.ts`) จะ shuffle ตัวเลือกตอน render ทำให้ position bias ที่เก็บใน JSON ไม่ถึงผู้ใช้
  จริง (เป็นแค่ content-quality signal) แต่ position bias ยังบอกว่ากระบวนการเขียน distractor มีรูรั่ว
  เชิงระบบ — เกทควรจับไว้ตั้งแต่ต้น ไม่ใช่รอจน length heuristic แสดงอาการก่อน (การรายงานแยกสองป้ายกำกับ
  ชัดเจน "reaches the user" vs "content-quality signal only" ใน output แทน เพื่อไม่ให้ตกใจกับตัวเลข
  position ที่จริง ๆ ไม่กระทบผู้ใช้ และไม่ชะล่าใจกับตัวเลข length ที่กระทบจริง)
- **`expected_answer` ไม่อยู่ใน options ของตัวเอง → abort (error) ไม่ silent-skip**: การข้ามคำถามที่
  เสียหายไปเงียบ ๆ จะลดขนาด N โดยไม่มีใครรู้ตัว — เป็นความผิดพลาดแบบเดียวกับที่ ticket นี้มีไว้ป้องกัน
  (วัดผิดโจทย์แบบไม่รู้ตัว) เพียงแค่คนละจุดในไปป์ไลน์ tool นี้อ่านไฟล์ตรง ไม่ผ่าน domain
  (`RecallCheck`'s `slices.Contains(options, answer)` check ไม่ได้อยู่ในเส้นทางนี้เลย) จึงต้องเช็คเอง
  แล้ว abort พร้อม path + expected_answer ที่ error message เพื่อ debug ได้ทันที
- **ไม่ reuse `curriculuminfra.LoadLesson`**: loader นั้น decode ด้วย `DisallowUnknownFields()` และ
  สร้าง domain `Lesson` เต็มรูปแบบ ซึ่งจะ reject ไฟล์ที่เสียหายก่อนที่ tool นี้จะเห็นด้วยซ้ำ (ขัดกับ
  เหตุผลข้อก่อนหน้า) แถมต้อง resolve topic/concept_id กับ concept row จริงใน MySQL — tool นี้ไม่ต้องมี
  DB เลยตามที่ ticket ระบุ จึง decode เฉพาะฟิลด์ที่ต้องใช้ (`topic`, `recall_checks[].{type,
  expected_answer, options}`) ด้วย DTO ของตัวเอง
- **`-max-excess` default = `-1` (ปิด gate)**: รัน tool เฉยๆ (ไม่ใส่ flag) ต้องไม่ fail อะไรเลย เป็น
  ค่าเริ่มต้นที่ปลอดภัยสำหรับตอนนี้ (ยังไม่ wire เข้า CI) — ticket ในอนาคตที่ wire เข้า CI/Makefile
  ต้องเลือกค่า threshold เองตอนนั้น
- **Track = ค่า `"topic"` field ในไฟล์ JSON โดยตรง** ไม่ใช่ชื่อโฟลเดอร์: ทำให้ synthetic corpus ใน
  เทสต์สร้าง `Question` struct ตรง ๆ ได้โดยไม่ต้องพึ่ง filesystem เลย (`Measure`/`Report` ไม่รู้จัก
  concept "โฟลเดอร์" เลย) — ของจริงสอง field นี้ตรงกันเสมอเพราะ `curriculuminfra.LoadLesson` บังคับไว้
  ตอน import แต่ tool นี้ไม่ได้พึ่งการบังคับนั้น (อ่านไฟล์ตรง)

**การตัดสินใจเพิ่มจาก round 2 (code review):**

- **Gate ต่อ track ไม่ใช่ pooled overall (`EvaluateTrackGates`)**: code review ชี้ตัวเลขจริงว่า track
  4-option ขนาด 550 ข้อที่ excess +60% เมื่อ pool รวมกับ corpus 3-option เดิมที่ fair (356 ข้อ) จะเหลือ
  แค่ **+32.3% pooled** — เจือจางจนหลุดเกทที่ตั้งไว้เข้มกว่านั้นได้ ทวนด้วยตัวเลขเล็กกว่าที่ derive เอง
  (fair track 300 ข้อ excess 0% + bad track 100 ข้อ excess 60% → pooled 15%) และยืนยันด้วยเทสต์
  `TestEvaluateTrackGates_DilutionExample_PerTrackCatchesWhatPooledWouldMiss` — เกทจึงประเมินทีละ
  track (fail ถ้า track ใด track หนึ่ง fail) ไม่ใช่ตัวเลข pooled เดียว ซึ่งตรงกับเจตนาที่ AWS-C1..C4
  เป็น track ใหม่ (4-option) ที่ต้อง gate แยกจาก track เดิม (3-option) อยู่แล้ว รายงาน "Overall" ยัง
  พิมพ์ไว้เพื่อดูภาพรวม แต่ label ชัดว่า "context only" ไม่ใช่ตัวตัดสิน pass/fail
- **`MinSampleSize = 100`**: heuristic ที่ n ต่ำเกินไปมี false-positive rate สูงจาก sampling noise
  ล้วน ๆ ไม่ใช่ signal จริง — derive เอง (exact one-sided binomial tail, threshold=baseline×
  (1+max-excess), ที่ max-excess=0.25) ที่ **p=baseline=0.25 (corpus 4 ตัวเลือก) ตลอดทั้งตาราง**:
  heuristic ที่ fair จริง (hit rate == baseline) ยัง false-FAIL **~19.7% ที่ n=30**, **~6.9% ที่
  n=100**, **~0.05% ที่ n=550** (คำนวณเองสองวิธี — log-gamma กับ iterative pmf recurrence — ตรงกันเป๊ะ)
  — **แก้ round 3**: เลข n=30 ที่ reviewer อ้างไว้ตอน round 2 (23.9%) ตอนแรกดูเหมือนตรวจสอบซ้ำไม่ได้
  round 3 หาสาเหตุเจอ: 23.9% คำนวณที่ **p=0.20 (1/5)** ไม่ใช่ p=0.25 (1/4) — ตั้งใจแสดงตัวอย่างของ
  `Position[4]` บน corpus 5 ตัวเลือก (multi-response) ไม่ใช่ baseline ของตารางหลัก และตารางของ reviewer
  เอง**ผสม baseline ข้ามแถวโดยไม่ระบุ** (p=0.20 ที่ n=5/n=30, p=0.25 ที่แถวอื่น) จึงไม่ใช่โมเดลเดียวที่
  self-consistent แบบที่ค่าคงที่ตัวเดียว (`MinSampleSize`) ต้องมี — ตารางในเอกสารนี้ยังเป็น p=0.25
  ตลอดทั้งตาราง (self-consistent) ตามเดิม เลือก 100 เป็นจุดกึ่งกลางที่จงใจ ไม่ใช่ค่าที่ปลอดภัยที่สุด
  (n=550 ปลอดภัยกว่ามากแต่จะกันบางช่วง n ไม่ให้ถูก gate ได้เลยตลอดไป — ดู `MinSampleSize`'s doc comment
  ใน `gate.go` สำหรับตัวอย่าง track จริงที่ n=100 ยังกันได้แต่ n=550 จะกันไม่ได้)
- **ขั้นต่ำ 2 ตัวเลือกต่อ mcq (`minMeasurableOptions`)**: mirror `internal/curriculum/domain
  /recallcheck.go`'s `minMCQOptions = 2` — คำถามที่มีแค่ 1 ตัวเลือกจะผ่านเช็ค "expected_answer อยู่ใน
  options" แบบไม่มีความหมาย (มีตัวเดียวให้ match พอดี) แล้ว `Baseline()` จะได้ 1.0 ซึ่งดึง `AvgBaseline`
  ของทั้ง corpus ขึ้นแบบเงียบ ๆ ทำให้ excess ที่คำนวณดูต่ำกว่าจริง — เป็น failure mode แบบเดียวกับ
  `expected_answer` ที่ไม่อยู่ใน options (ลด N แบบไม่มีใครรู้) แค่คนละกลไก จึง abort ด้วย error รูปแบบ
  เดียวกัน (path + option count)
- **`longestOptionIndex` unexport**: ฟังก์ชันนี้ panic ถ้าได้ slice ว่าง และมีแค่ `accumulate` เป็น caller
  เดียวซึ่งตอนนี้การันตีแล้วว่า options มีอย่างน้อย 2 ตัวเสมอ (ผ่านเช็คข้างบน) ก่อนเรียก — unexport
  ตัดโอกาสที่ code ภายนอก package จะเรียกตรง ๆ ด้วย slice ว่างแล้ว panic ออกไปทั้งหมด (แทนที่จะเพิ่ม guard
  แล้วยังคง export ไว้ ซึ่งเพิ่มพื้นที่ผิวของ public API โดยไม่มี caller ที่ถูกต้องคนไหนต้องการมันจริง)

### Synthetic-corpus test table (ทุกกรณีที่ ticket ระบุไว้ขั้นต่ำ)

| # | Corpus | Assertion | Test |
|---|---|---|---|
| 1 | ทุกข้อคำตอบถูกคือตัวเลือกที่ยาวที่สุด (ไม่มีเสมอ) | longest heuristic = 100% | `TestMeasure_CorrectAnswerAlwaysLongest_LongestHeuristicIsHundredPercent` |
| 2 | ตัวเลือกทุกตัวยาวเท่ากัน (เสมอทั้งหมด), คำตอบถูกวนสม่ำเสมอ index 0/1/2 (100 ข้อ/index จาก 300 ข้อ) | longest heuristic = baseline พอดี (ภายใน 1e-9) | `TestMeasure_AllOptionsEqualLength_LongestHeuristicNearBaseline` |
| 3 | คำตอบถูกอยู่ index 1 เสมอ (40 ข้อ, 4 ตัวเลือก) | position[1] = 100%, position[0,2,3] = 0% พอดี | `TestMeasure_CorrectAnswerAlwaysIndex1_PositionHeuristic` |
| 4 | ผสม track 3 ตัวเลือก (3 ข้อ) กับ track 4 ตัวเลือก (3 ข้อ) | baseline ต่อ track ต่างกัน (1/3 vs 1/4), overall baseline = ค่าเฉลี่ยถ่วงน้ำหนัก 0.29166... ไม่ตรงกับ track ใด track หนึ่งเป๊ะ | `TestMeasure_MixedOptionCounts_PerTrackBaselinesDifferAndOverallIsNotEither` |
| 5 | corpus สร้างให้ longest heuristic excess = 0.20 พอดี (30/100 hit, baseline 0.25) — options สร้างให้ index 3 เป็นตัวยาวที่สุดของทุกข้อ ดังนั้น "guess longest" กับ "guess index 3" คือกลยุทธ์เดียวกันในคอร์ปัสนี้ ทั้งสอง heuristic จึงข้าม threshold พร้อมกันเสมอ | gate ผ่านที่ max=0.21, gate ไม่ผ่านที่ max=0.19 พร้อม violation **2 รายการ** (longest-option + position index 3) | `TestMeasure_CorpusStraddlingGate_ExitCodeFlips` (internal, assert จำนวน violation ด้วย), `TestRun_GateStraddlingCorpus_ExitCodeFlips` (end-to-end ผ่าน CLI `run()` จริง) |
| — | corpus 3-ตัวเลือกล้วน / 4-ตัวเลือกล้วน (แยกกัน) | baseline = 1/3 พอดี / baseline = 1/4 พอดี (ฆ่า hardcode ทั้งสองทิศทางแยกกัน) | `TestMeasure_BaselineForThreeOptionCorpusIsOneThird`, `TestMeasure_BaselineForFourOptionCorpusIsOneQuarter` |
| — | Thai option (4 runes/12 bytes) vs English option (8 runes/8 bytes) | `longestOptionIndex` เลือกตาม rune count ไม่ใช่ byte count | `TestLongestOptionIndex/rune_count,_not_byte_count` |
| — | ไฟล์จริงผสม mcq + short_answer | `LoadQuestions` คืนเฉพาะ mcq | `TestLoadQuestions_FiltersShortAnswer` |
| — | คำถามที่ `expected_answer` ไม่อยู่ใน options | `Measure` abort พร้อม error ระบุ source file + answer | `TestMeasure_ExpectedAnswerNotInOptionsAborts` |
| — (round 2) | คำถามมีแค่ 1 ตัวเลือก (`expected_answer` ตรงกับตัวเดียวนั้นพอดี) | `Measure` abort พร้อม error ระบุ source file + จำนวน option จริง — ไม่ใช่ผ่านแบบเงียบ ๆ | `TestMeasure_TooFewOptionsAborts` |
| — (round 2) | ไฟล์ JSON พังจริง (truncated, ไม่ปิด bracket) | `LoadQuestions` คืน error ระบุชื่อไฟล์ — ไม่ใช่ตัดไฟล์ทิ้งเงียบ ๆ แล้วนับ N น้อยลงโดยไม่บอกใคร | `TestLoadQuestions_MalformedJSONErrors` |
| — (round 2) | Report ที่ longest heuristic สะอาด (excess 0) แต่ position[2] แย่ (excess 28% เทียบ threshold 20%) | `EvaluateGate` ต้อง fail จาก position เพียงอย่างเดียว | `TestEvaluateGate_PositionAloneCanFailWithLongestClean` |
| — (round 2) | heuristic ที่ n=99 (ต่ำกว่า `MinSampleSize`=100) hit rate 100% | `EvaluateGate` ต้อง**ไม่** fail แค่รายงาน insufficient | `TestEvaluateGate_BelowMinSampleSizeIsInsufficientNotFailedNorJudgedPass` (round 3: เพิ่ม assert `Judged==0` + `ExitCode==1`) |
| — (round 2) | track ใหญ่ fair (300 ข้อ, excess 0%) pool กับ track เล็กแย่ (100 ข้อ, excess 60%) → pooled 15% | `EvaluateTrackGates` fail (จาก track แย่) ทั้งที่ `EvaluateGate` บน pooled report จะผ่าน | `TestEvaluateTrackGates_DilutionExample_PerTrackCatchesWhatPooledWouldMiss` (internal), `TestRun_PerTrackGateCatchesWhatPooledWouldMiss` (end-to-end ผ่าน CLI จริง) |
| — (round 3, N2) | ทุกข้อคำตอบถูกคือตัวเลือกที่สั้นที่สุด (ไม่มีเสมอ) | shortest heuristic = 100% | `TestMeasure_CorrectAnswerAlwaysShortest_ShortestHeuristicIsHundredPercent` |
| — (round 3, N2) | 3 ตัวเลือกความยาวต่างกัน 3 ค่า (สั้น/กลาง/ยาว ภาษาไทย), คำตอบถูกคือตัวกลางเสมอ | middle heuristic = 100%, longest/shortest = 0% พอดี (คำตอบไม่เคยเป็นตัวสุดขั้วเลย) | `TestMeasure_CorrectAnswerAlwaysMiddle_MiddleHeuristicIsHundredPercent` |
| — (round 3, N2) | boundary: 2 ตัวเลือก (ไม่มีตัวกลางเชิงโครงสร้าง) + 3 ตัวเลือกที่ยาวชนกันเหลือ 2 ค่า (ไม่มีตัวกลางเพราะเสมอ) | `Middle.Total == NumMCQs` เสมอ (**นับเป็น miss ไม่ exclude** — จุดตัดสินใจที่ทำให้ 89.6%/33.7% reproduce ได้พอดี ดูหัวข้อ "Finding") | `TestMeasure_MiddleOptionIndex_BoundaryCases` |
| — (round 3, N2) | หน่วย `middleOptionIndex` เอง: 3 ตัวเลือก (1 ตัวกลาง), 4 ตัวเลือก (2 ตัวกลาง เสมอเลือก index ต่ำสุด), 2 ตัวเลือก, 3 ตัวเลือกยาวชนกัน 2 ค่า, ทุกตัวยาวเท่ากัน, 5 ตัวเลือก (3 ตัวกลางเสมอกัน) | ok/index ตรงตามนิยามทุกกรณี | `TestMiddleOptionIndex` (table-driven 6 เคส) |
| — (round 3, N1) | Report ที่ทุก heuristic มี n ต่ำกว่า `MinSampleSize` (ไม่มีอันไหนถูกตัดสินเลย) | `Judged == 0`, `ExitCode == 1` (ไม่ใช่ 0) แม้ `Failed == false` | `TestEvaluateGate_BelowMinSampleSizeIsInsufficientNotFailedNorJudgedPass`, `TestExitCode` (เคส "requested, nothing judged") |
| — (round 3, N1) | corpus จริง 50 ข้อ track เดียว คำตอบถูกเป็นตัวยาวสุด 100% (guessable ที่สุดเท่าที่สร้างได้) แต่ n=50<100 | `run()` พิมพ์ `Gate: NOT JUDGED` (ไม่ใช่ `PASS`) และ exit 1 | `TestRun_NotJudgedWhenEveryHeuristicIsBelowMinSampleSize` (end-to-end ผ่าน CLI จริง) |

**Tolerance ที่ใช้ (`floatEps = 1e-9`) และทำไมไม่ flaky**: ทุก fixture ในไฟล์เทสต์เหล่านี้สร้างด้วยมือ
แบบ deterministic ล้วน ๆ ไม่มีการสุ่มเลยสักจุดเดียว (เช่น กรณี #2 ข้างบน คำตอบถูกถูก "วน" index
0/1/2/0/1/2/... ด้วย `i % 3` ไม่ใช่สุ่ม) ดังนั้น "tolerance" ในที่นี้ทำหน้าที่ดูดซับแค่ float64 rounding
error (เช่น `100.0/300.0` อาจต่างจาก `1.0/3.0` ไปสัก 1 ULP แม้ค่าเท่ากันทางคณิตศาสตร์) ไม่ใช่ดูดซับความ
แปรปรวนทางสถิติ — เพราะไม่มี sampling ในเทสต์เหล่านี้เลย จึงไม่มีทางที่ tolerance ตัวนี้จะ flake ไม่ว่า
จะรันกี่รอบ ค่า `1e-9` ยังห่างจากบั๊กจริงที่ต้องจับได้ (เช่น off-by-one ที่ทำให้ hit นับผิดไปแค่ 1 ข้อ
จาก 300 ≈ 0.33 percentage point) อยู่ประมาณ 6 อันดับ (order of magnitude) จึงแคบพอที่จะจับบั๊กจริงแต่ไม่
แคบจนจับ float rounding ผิดเป็นบั๊ก

### Mutation table (17 mutation — 9 เดิม + 2 ที่ code review round 1 พบว่ารอด + 6 ใหม่จาก N1/N2 round 3 — ทุกตัวตายจริง ไม่มี survivor)

| # | Mutation | Killed by |
|---|---|---|
| 1 | baseline hardcode เป็น 1/3 (`Question.Baseline()` คืน `1.0/3.0` เสมอ) | `TestQuestionBaseline`, `TestMeasure_BaselineForFourOptionCorpusIsOneQuarter`, `TestMeasure_MixedOptionCounts_PerTrackBaselinesDifferAndOverallIsNotEither`, `TestMeasure_CorpusStraddlingGate_ExitCodeFlips`, `TestRun_GateStraddlingCorpus_ExitCodeFlips` |
| 2 | baseline hardcode เป็น 1/4 (`Question.Baseline()` คืน `0.25` เสมอ) | `TestQuestionBaseline`, `TestMeasure_AllOptionsEqualLength_LongestHeuristicNearBaseline`, `TestMeasure_MixedOptionCounts_PerTrackBaselinesDifferAndOverallIsNotEither`, `TestMeasure_BaselineForThreeOptionCorpusIsOneThird` |
| 3 | วัดความยาวเป็น byte แทน rune (`longestOptionIndex` ใช้ `len()` แทน `utf8.RuneCountInString`) | `TestLongestOptionIndex` (เคส Thai vs English) |
| 4 | `longestOptionIndex` เลือกตัวสั้นที่สุดแทนยาวที่สุด (`l < bestLen` แทน `l > bestLen`) | `TestLongestOptionIndex`, `TestMeasure_CorrectAnswerAlwaysLongest_LongestHeuristicIsHundredPercent` |
| 5 | position heuristic off-by-one (`pos == expectedIdx+1` แทน `pos == expectedIdx`) | `TestMeasure_CorrectAnswerAlwaysIndex1_PositionHeuristic` |
| 6 | `short_answer` หลุดเข้ามาในตัวหาร (ลบ `if rc.Type != "mcq" { continue }` ใน `loadFile`) | `TestLoadQuestions_FiltersShortAnswer` |
| 7 | per-track aggregation ยุบเป็น global baseline เดียว (`byTrack[q.Track] = overall` แทน `= track`) | `TestMeasure_MixedOptionCounts_PerTrackBaselinesDifferAndOverallIsNotEither`, `TestEvaluateTrackGates_DilutionExample_PerTrackCatchesWhatPooledWouldMiss` |
| 8 | gate comparison กลับด้าน (`excess < maxExcessRatio` แทน `excess > maxExcessRatio`) | `TestEvaluateGate_PassesBelowThreshold`, `TestEvaluateGate_FailsAboveThreshold`, `TestEvaluateGate_ExactlyAtThresholdPasses`, `TestEvaluateGate_PositionAloneCanFailWithLongestClean`, `TestEvaluateTrackGates_PassesWhenAllTracksPass`, `TestEvaluateTrackGates_FailsWhenAnyTrackFails`, `TestEvaluateTrackGates_DilutionExample_PerTrackCatchesWhatPooledWouldMiss`, `TestMeasure_CorpusStraddlingGate_ExitCodeFlips` |
| 9 | exit code เป็น 0 เสมอ (`ExitCode` คืน `0` ไม่เช็ค `g.Failed`) | `TestExitCode`, `TestMeasure_CorpusStraddlingGate_ExitCodeFlips` (internal), `TestRun_GateStraddlingCorpus_ExitCodeFlips`, `TestRun_PerTrackGateCatchesWhatPooledWouldMiss` (ผ่าน CLI `run()` จริงทั้งคู่ — ปิดช่องที่ mutation อาจรอดถ้าทดสอบแค่ internal function) |
| 10 (round 2, code review พบว่ารอด) | ลบ loop เช็ค position heuristic ทั้งก้อนออกจาก `EvaluateGate` (เหลือเช็คแค่ longest/shortest/middle) — เดิม**รอด**เพราะทุกเทสต์ pass/fail ก่อนหน้าใช้ `reportWithLongestResult` ที่ปล่อย `Position` ว่างเปล่า จึงไม่มีทางเห็นความแตกต่าง | `TestEvaluateGate_PositionAloneCanFailWithLongestClean` (Report ที่ longest สะอาดแต่ position[2] แย่, ต้อง fail ได้จาก position อย่างเดียว) |
| 11 (round 2, code review พบว่ารอด) | `loadFile`'s decode-error path เปลี่ยนเป็น `return nil, nil` (กลืน error, ทำเหมือนไฟล์ไม่มี mcq เลย) — เดิม**รอด**เพราะไม่มีเทสต์ไหนป้อนไฟล์ JSON ที่พังจริง | `TestLoadQuestions_MalformedJSONErrors` |
| 12 (round 3, N2) | `shortestOptionIndex` เลือกตัวยาวที่สุดแทนสั้นที่สุด (`l > bestLen` แทน `l < bestLen`) | `TestShortestOptionIndex`, `TestMeasure_CorrectAnswerAlwaysShortest_ShortestHeuristicIsHundredPercent`, `TestMeasure_CorpusStraddlingGate_ExitCodeFlips` (ผลข้างเคียง: shortest กลายเป็นเหมือน longest บนฟิกซ์เจอร์นั้น ทำให้จำนวน violation ผิดไปด้วย) |
| 13 (round 3, N2) | `middleOptionIndex`'s classification ผิด (`l >= minLen && l <= maxLen` แทน `l > minLen && l < maxLen`, ทำให้ทุกตัวเลือก "ผ่าน" เงื่อนไข) | `TestMiddleOptionIndex`, `TestMeasure_CorrectAnswerAlwaysMiddle_MiddleHeuristicIsHundredPercent`, `TestMeasure_MiddleOptionIndex_BoundaryCases` |
| 14 (round 3, N2/N3) | `Middle`'s "count as miss" ถูก revert กลับเป็น "exclude" (`if middleOK { r.Middle.add(...) }` แทน `r.Middle.add(middleOK && ..., baseline)`) — **นี่คือ mutation ที่สำคัญที่สุดของรอบนี้**: ถ้ารอด ตัวเลข 89.6%/33.7% จะ reproduce ไม่ได้อีก (จะได้ 36.6%/29.5% ที่ denominator 328 แทน) | `TestMeasure_MiddleOptionIndex_BoundaryCases` (assert `Middle.Total == NumMCQs` ตรง ๆ) |
| 15 (round 3, N1) | `ExitCode`'s `g.Judged == 0` เงื่อนไขถูกลบ (เหลือเช็คแค่ `g.Failed`) | `TestExitCode` (เคส "requested, nothing judged"), `TestEvaluateGate_BelowMinSampleSizeIsInsufficientNotFailedNorJudgedPass`, `TestRun_NotJudgedWhenEveryHeuristicIsBelowMinSampleSize` (ผ่าน CLI จริง) |
| 16 (round 3, N1) | `EvaluateGate`'s `check()` ลืม `g.Judged++` (heuristic ที่ถูกตัดสินจริงไม่ถูกนับ) | `TestEvaluateGate_PassesBelowThreshold` (assert `Judged == 1` ตรง ๆ — เพิ่ม assertion นี้ใน round 3 เพื่อปิดช่องนี้โดยเฉพาะ), `TestRun_GateStraddlingCorpus_ExitCodeFlips`, `TestMeasure_CorpusStraddlingGate_ExitCodeFlips` |
| 17 | mcq มี 1 ตัวเลือก (`minMeasurableOptions` ลดจาก 2 เหลือ 1) | `TestMeasure_TooFewOptionsAborts` |

**Survivor**: ไม่มี — ทั้ง 17 mutation ตายจริงทุกตัว (revert กลับหลัง confirm แล้วทุกจุด, `git status`
สะอาดหลังทำเสร็จ — ยืนยันซ้ำทั้ง 11 mutation เดิมจาก round 1-2 ด้วยว่ายังตายอยู่หลังการรีแฟกเตอร์ของ round 3
ทั้งการย้าย package และการเพิ่ม shortest/middle/tri-state gate)

**หมายเหตุ scope**: mutation "exit code always zero" ทดสอบทั้งที่ `internal/mcqguess.ExitCode`
(unit test ตรง ๆ) และที่ `cmd/mcq-guessability`'s `run()` (end-to-end ผ่าน CLI จริง) — ปิดช่องว่างที่
เคยมีในโปรเจกต์นี้ (ไม่มี `cmd/*/main.go` ไหนมี `_test.go` มาก่อนเลยทั้ง `import-lessons`,
`import-curriculum`, `api`) เพราะ mutation table เรียกร้องตรง ๆ ว่าต้องมีเทสต์ฆ่า ไม่ใช่แค่ verify มือ
ส่วน `main()` เองที่แค่เรียก `os.Exit(run(...))` (3 บรรทัด) ยังไม่มี test ตรง (ตาม convention เดิมของ
repo ที่ `main()` ไม่ถูกเทสต์ตรง ๆ) — แต่ `run()` ที่ `main()` เรียกครอบคลุมด้วย `main_test.go` เต็มแล้ว

### ตัวเลขที่วัดได้จาก corpus จริง

คำสั่ง: `go run ./cmd/mcq-guessability -dir content/lessons`

corpus ปัจจุบันมี **356 MCQ ทุกข้อมี 3 ตัวเลือก** (baseline 33.3% เท่ากันทุกข้อ, ยืนยันตรงกับตัวเลข
586 recall_checks รวม = mcq 356 + short_answer 230 ที่บันทึกไว้ใน roadmap.md):

| Heuristic | Actual | Baseline avg | Excess (relative) |
|---|---|---|---|
| Longest option | 32.6% | 33.3% | −2.2% |
| Shortest option | 31.5% | 33.3% | −5.6% |
| Middle option | **33.7%** | 33.3% | +1.1% |
| Position index 0 (guess "A" เสมอ) | 33.4% | 33.3% | +0.3% |
| Position index 1 (guess "B" เสมอ) | 33.4% | 33.3% | +0.3% |
| Position index 2 (guess "C" เสมอ) | 33.1% | 33.3% | −0.6% |

ทุกตัวเลขอยู่ใกล้ baseline มาก (excess ทุกตัว < 6% แบบ relative) — ไม่มี tell ที่มีนัยสำคัญใน corpus
ปัจจุบัน ทั้ง longest/shortest/middle สอดคล้องกับที่ C-mcq-balance (#43) บันทึกไว้ว่าแก้ปัญหาเดิม (89.6%)
แล้วจริง ไม่ใช่แค่บาง heuristic

รันด้วย `-max-excess 0.15` (ตัวอย่างการใช้เป็นเกท): `Gate: PASS` (ทุกตัวเลขข้างบนอยู่ใต้ 15% หมด)

### Finding — 89.6% และ 33.7% (เดิม) คือค่า middle-length heuristic วัดผิดชื่อว่า "length" (round 3 คลี่ปม)

**Round 1 ของ code review จับได้ว่าข้อความเดิมของหัวข้อนี้ผิดสองจุด**: (R1) อ้างว่า 89.6% "วัดซ้ำไม่ได้
โดยหลักการเพราะ corpus ถูก rewrite" ซึ่งเป็น**ข้อสรุปที่ประดิษฐ์ขึ้นเอง**; (R2) อ้างว่า tie-break rule
"เป็นไปได้มากที่สุด" ที่อธิบายช่องว่าง 33.7% vs 32.6% ทั้งที่ไม่เคยคำนวณตัวเลขจริงมารองรับ ทั้งสองจุดคือ
**การเชื่อคำกล่าวอ้างโดยไม่ re-derive** แก้ด้วยการนับ tie จริง (13 ข้อเสมอ, 12 ชนะได้) แล้วคำนวณขอบเขต
ของทุก correctness-blind tie-break rule ได้ **30.90%–34.27%** สรุปว่า 33.7% (120 hits, ต้องถูก 10/12
tie, p≈0.85% ถ้าสุ่ม) อยู่นอกช่วงที่ tie-break อธิบายได้ — ถูกต้องแล้วที่ปฏิเสธ tie-break เป็นคำตอบ แต่
round 2 หยุดที่ "unexplained" ทั้งที่คำตอบจริงอยู่ในเอกสารของโปรเจกต์เองมาตลอด

**Round 3 พบคำตอบ**: `docs/tickets/mcq-quality.md`'s "บทเรียนที่ 1" บันทึกไว้ตรง ๆ อยู่แล้วว่า
"รอบแรกตั้งกฎว่า 'คำตอบที่ถูกห้ามยาวที่สุดและห้ามสั้นที่สุด' ... เมื่อมี 3 ตัวเลือก กฎนี้บังคับให้คำตอบเป็น
'ตัวกลาง' เสมอ → เดาถูก 89.6%" — ประโยคนี้บอกตรง ๆ ว่า heuristic ที่แท้จริงคือ **middle-length** (เดาว่า
คำตอบคือตัวเลือกที่ไม่ยาวสุดไม่สั้นสุด) ไม่ใช่ longest ที่ tool นี้วัดมาตลอดจนถึง round 2 — เพิ่ม
`Report.Middle`/`middleOptionIndex` (ดูหัวข้อ "N2" ด้านล่าง) แล้ววัดซ้ำ **ทั้งสองตัวเลขที่บันทึกไว้
reproduce แม่นเป๊ะ**:

| corpus (commit) | middle-length rate | บันทึกไว้เดิม | longest-length rate (ที่ tool วัดมาตลอด round 1-2) |
|---|---|---|---|
| `7dfb0ba` (หลัง #38 "C-mcq-sweep") | **319/356 = 89.6%** | 89.6% | 2.0% |
| `057ec9a` / ปัจจุบัน (หลัง #43 "C-mcq-balance") | **120/356 = 33.7%** | 33.7% | 32.6% |

ที่ `7dfb0ba` โดยเฉพาะ: ในกลุ่มคำถามที่มี "ตัวกลาง" จริง (319 จาก 356 ข้อ, อีก 37 ข้อความยาวชนกันเหลือแค่
2 ค่าจึงไม่มีตัวกลาง) **compliance กับกฎ "ห้ามยาวสุด/สั้นสุด" คือ 100% พอดี** (319/319) — ตรงกับที่
mcq-quality.md's "บทเรียนที่ 1" อธิบายไว้ว่า "รอบแรกวัดว่ากฎถูกละเมิดไหม (ผ่าน 100%)" เป๊ะ ทุกตัวเลขต่อกัน
สนิท ไม่มีจุดไหนต้องเดาอีกต่อไป

**นี่คือคำอธิบายว่าทำไม 41.6% (position) ตรงเป๊ะตั้งแต่ round 1 แต่ 89.6% (length) ไม่ตรงจนกระทั่งตอนนี้**:
position heuristic ของทั้งสอง tool (เดิมกับใหม่) วัดสิ่งเดียวกันมาตลอด (เดาตำแหน่งเดิมเสมอ) จึงตรงกันตั้งแต่
ต้น — แต่ "length heuristic" ของ tool เดิมกับของ tool นี้ (round 1-2) **วัดคนละอย่าง**: tool นี้วัด
longest-length, tool เดิมวัด middle-length เครื่องมือวัดถูกต้องทั้งคู่ (calibrated ถูก) แค่ชี้วัดคนละ
เป้าหมาย — ห้ามเทียบตัวเลข "length" ข้ามสอง tool กันอีกต่อไปโดยไม่ระบุว่าเป็น longest หรือ middle

**ทำไม 6b1a765 (คอมมิตก่อนหน้า 7dfb0ba) ไม่ reproduce ตัวเลขไหนเลย**: `6b1a765` คือคอมมิตก่อน `7dfb0ba`
(#38) หนึ่งขั้น — วัดได้ longest **36.8%**, shortest **28.4%**, middle **29.5%** (position index 0
**41.6%** เท่ากับที่ `7dfb0ba` เป๊ะ เพราะ #38 แก้แค่**เนื้อหา**ตัวเลือกผิด ไม่ได้ย้ายตำแหน่งคำตอบ) —
ไม่มีตัวไหนใกล้ 89.6% เลย ยืนยันว่า defect 89.6% (compliance กับกฎ "ห้ามยาวสุด/สั้นสุด") **เพิ่งเกิดขึ้น
จากการแก้ของ #38 เอง** ไม่ใช่มีอยู่ก่อนแล้ว — #38 แก้ปัญหาเดิม (อะไรก็ตามที่ทำให้ longest สูงที่ 6b1a765)
แต่สร้างปัญหาใหม่ (middle-length compliance 100% ในข้อที่ตัดสินได้) ซึ่งเป็นเหตุผลที่ C-mcq-balance (#43,
4 commit ถัดไป) ต้องตามมาแก้อีกที

**สิ่งที่ยังไม่รู้ (ตรงไปตรงมา)**: อัลกอริทึมที่แท้จริงของ script เดิม (เช่น นิยาม tie-break ของมันสำหรับ
"ตัวกลาง" ตอน 37 ข้อความยาวชนกัน) ไม่มีให้ตรวจสอบ เพราะไม่เคย commit ไว้ในโปรเจกต์ — แต่คำถามหลักที่ round
1-2 ค้างไว้ ("length" เดิมวัดอะไรกันแน่) ได้คำตอบแล้ว: มันวัด middle ไม่ใช่ longest

### N2 — ทำไมต้องมี middle-length heuristic (ไม่ใช่แค่ longest)

`docs/tickets/mcq-quality.md`'s ตัวชี้วัดกำหนดไว้ตั้งแต่ต้นว่าเดาด้วยกฎ "ยาว/สั้น/กลาง" ต้องได้ ≈baseline
ทั้งสามแบบ — tool นี้จนถึง round 2 implement แค่ longest ตัวเดียว ซึ่ง**ตาบอด**ต่อ corpus ที่มีข้อบกพร่อง
แบบ "ห้ามคำตอบยาวสุด/สั้นสุด" (กฎที่ **ฟังดูเป็น best practice** และมีโอกาสสูงที่คนเขียนคำถาม AWS ในอนาคต
จะเผลอใช้ตรง ๆ) — บน corpus `7dfb0ba` ที่ middle-length อ่านได้ 89.6%, longest heuristic (ที่ tool วัดมา
ตลอด) อ่านได้แค่ **2.0%** เกทที่เช็คแค่ longest จะรายงานว่า corpus นี้ "สะอาดผิดปกติ" ทั้งที่จริงเดาถูก
เกือบทุกข้อ — เป็นการพลาดแบบเดียวกับที่ #38 เคยพลาดมาก่อน เพียงแค่คนละเครื่องมือ

**สิ่งที่ทำ**: เพิ่ม `Report.Longest`/`Shortest`/`Middle` (จากเดิมมีแค่ `Length`), `shortestOptionIndex`
(mirror ของ `longestOptionIndex` เป๊ะ ๆ), `middleOptionIndex` (ใหม่) — เกทเช็คทั้งสามตัวพร้อม position
ทุกตัวรายงานด้วย framing เดียวกับ longest เดิม ("survives shuffle — reaches the user")

**การตัดสินใจเรื่อง middle สำหรับ 3 ตัวเลือก vs 4+ ตัวเลือก (ตามที่ ticket ขอให้ตัดสินใจ + ทดสอบ boundary)**:

- **นิยาม**: ตัวเลือกที่ความยาว (rune count) อยู่**ระหว่าง**ค่าต่ำสุดกับสูงสุดของคำถามนั้นอย่างเคร่งครัด
  (`l > minLen && l < maxLen`) — ไม่ใช่ตัวสุดขั้วทั้งสองด้าน สำหรับ 3 ตัวเลือกไม่มีเสมอกัน มีตัวกลางแบบนี้
  พอดี 1 ตัวเสมอ; สำหรับ 4+ ตัวเลือกอาจมีมากกว่า 1 ตัว (เช่น ความยาว 3,6,7,9 → 6 กับ 7 ทั้งคู่เข้าเงื่อนไข)
  — เสมอกันแก้ด้วย index ต่ำสุด กติกาเดียวกับ longest/shortest (deterministic, ไม่มี RNG)
- **คำถามที่ไม่มีตัวกลางเลย (2 ตัวเลือก, หรือ 3+ ตัวเลือกที่ความยาวชนเหลือแค่ 2 ค่า)**: `middleOptionIndex`
  คืน `ok=false` — **แต่ `Measure` ยังนับเป็น miss ใน `Middle.Total` ไม่ใช่ exclude ออกจากตัวหาร** นี่คือ
  จุดตัดสินใจที่สำคัญที่สุด: ลองแบบ exclude ก่อน (เหมือนที่ `Position[k]` ทำกับคำถามที่มีตัวเลือกน้อยกว่า
  k+1) แล้วพบว่า**ไม่ reproduce ตัวเลขที่บันทึกไว้** (`120/328=36.6%` ไม่ใช่ `120/356=33.7%`) — เปลี่ยนเป็น
  "count as miss" แล้ว reproduce ตรงเป๊ะทั้ง 89.6% และ 33.7% ดูหัวข้อ "Finding" ด้านบน เหตุผลที่ต่างจาก
  `Position`: "guess index 3" ไม่มีความหมายเลยสำหรับคำถาม 2 ตัวเลือก (ไม่มี index 3 ให้ถูกหรือผิด) แต่
  "เดาตัวเลือกที่ไม่สุดขั้ว" มีความหมายชัดเจนเสมอ แค่บังเอิญเดาไม่ถูกแน่นอนเมื่อทุกตัวเลือกเป็นตัวสุดขั้วหมด
  — การ exclude แทนที่จะนับ miss จะเปิดช่องให้ corpus ซ่อน tell จริงไว้หลังตัวหารที่หดตัวได้ ซึ่งเป็นกลไก
  เดียวกับที่ทำให้เหตุการณ์ 89.6% เดิมหลุดรอดมาได้ตั้งแต่ต้น

### N1 — Gate ที่ไม่ได้ตัดสินอะไรเลยต้องไม่รายงาน PASS

**พบโดย code review round 2**: corpus 50 ข้อ track เดียว คำตอบถูกเป็นตัวยาวสุด 100% ของทุกข้อ (guessable
ที่สุดเท่าที่สร้างได้) แต่ n=50 < `MinSampleSize`=100 → ทุก heuristic ถูกข้าม (insufficient) → เดิมพิมพ์
`Gate: PASS` และ exit code **0** — เพราะ `GateResult{Failed: false}` เป็นค่าเดียวกันทั้งกรณี "ตัดสินแล้ว
ผ่าน" และ "ไม่ได้ตัดสินอะไรเลย" ทั้งที่สองสถานะนี้ควรต่างกัน

บล็อกเพราะ `aws-cert.md` เขียนแผนไว้ว่าคลังข้อสอบเขียนทีละ domain — เมื่อ `MinSampleSize=100` batch แรก
ของทุก domain (ก่อนมีคำถามสะสมถึง 100 ข้อ) จะเข้าเงื่อนไข exempt นี้พอดี ซึ่งเป็นจังหวะที่ถูกที่สุดและ
มีค่าที่สุดที่จะจับ bias เชิงระบบตั้งแต่ต้น — เกทที่เงียบตอนนั้นพอดีคือเกทที่ใช้ป้องกันอะไรไม่ได้เลย

**สิ่งที่ทำ**: `GateResult` เพิ่ม `Requested bool` (true เมื่อ maxExcessRatio>=0) และ `Judged int`
(จำนวน heuristic ที่ผ่านเกณฑ์ MinSampleSize จนถูกตัดสินจริง) — `ExitCode` คืน 1 เมื่อ `Requested &&
(Failed || Judged==0)`, คืน 0 เมื่อไม่ได้ requested เลย (`maxExcessRatio<0`) `run()` พิมพ์
`Gate: NOT JUDGED` แยกจาก `Gate: PASS`/`Gate: FAIL` เมื่อ `Judged==0` — หลักการเดียวกับที่ `run()` คืน
exit 1 อยู่แล้วเมื่อไม่พบ mcq เลยในไดเรกทอรี ("เกทที่วัดอะไรไม่ได้เลยต้องไม่รายงานว่าผ่าน")

### สิ่งที่ตั้งใจไม่ทำในรอบนี้

- ไม่ wire เกทเข้า CI/Makefile จริง (ticket ระบุแค่ "a future ticket can wire it into CI or a
  Makefile target" — รอบนี้แค่ทำให้ flag/exit code ใช้งานได้)
- ไม่เลือกค่า `-max-excess` threshold ที่แน่นอนสำหรับ AWS-C1..C4 (ticket นี้ทำแค่เครื่องมือวัด ไม่ใช่
  ตรึงเกณฑ์ผ่าน/ไม่ผ่านของคลังข้อสอบที่ยังไม่มีเนื้อหาจริง)
- ไม่แก้ distractor ใด ๆ ในคลังเดิม (ไม่ใช่ scope ของ AWS-S3 ตามที่ spec ระบุไว้ตั้งแต่ต้น)
- ไม่แตะ Docker/DB เลย (ticket ระบุไม่ต้องมี DB — ยืนยันด้วยว่า
  `go test ./cmd/mcq-guessability/...` รันได้โดยไม่มี `docker compose up` ใด ๆ)
- ไม่เพิ่ม multiple-answer support ใน `Question`/`ExpectedIndex` (AWS-S2's scope) — tool นี้จึงวัด
  guessability ได้แค่ single-answer mcq เท่านั้นตอนนี้ ดู decision bullet เรื่อง scope ด้านบนและ
  `cmd/mcq-guessability/main.go`'s package doc
- (round 3) ไม่ตามหานิยาม tie-break ที่แท้จริงของ script เดิมสำหรับ 37 ข้อที่ความยาวชนกันเหลือ 2 ค่า
  (ที่ commit `7dfb0ba`) — คำถามหลักที่ค้างจาก round 1-2 ("length เดิมวัดอะไรกันแน่") มีคำตอบแล้ว (วัด
  middle ไม่ใช่ longest) ส่วนรายละเอียด tie-break ของ script ที่ไม่มีให้ตรวจสอบยังคงเป็นเช่นนั้นต่อไป
  ไม่ใช่สิ่งที่บล็อกอะไรอีกแล้ว

### Status: implemented, code review round 3 fixes applied, PR pending

### Review focus (round 1)

<details>
<summary>คำถามสำหรับรีวิว diff รอบนี้ (เฉลยพับไว้ด้านล่าง)</summary>

1. ทำไม `EvaluateGate` ต้องเช็คทั้ง length heuristic และทุก position-index heuristic ทั้งที่ Q-1
   shuffle ตัวเลือกตอน render ทำให้ position bias ไม่ถึงผู้ใช้จริงอยู่แล้ว?
2. ทำไม `internal/mcqguess` ถึงไม่ใช้ `curriculuminfra.LoadLesson` ทั้งที่มันมีอยู่แล้วและ decode
   ไฟล์เดียวกันเป๊ะ?
3. ทำไม excess ที่รายงานถึงเป็น ratio (`(actual-baseline)/baseline`) แทนที่จะเป็น percentage points
   สัมบูรณ์ (`actual-baseline`) ตรง ๆ?

<details>
<summary>เฉลย</summary>

1. เพราะ shuffle ที่ Q-1 แก้แค่**อาการ**ที่ผู้ใช้เห็น (ตัวเลือกสลับตำแหน่งก่อน render) ไม่ได้แก้**ต้นตอ**
   ในกระบวนการเขียนคำถาม — ถ้า distractor ถูกเขียนด้วยรูปแบบที่ทำให้คำตอบถูกกระจุกอยู่ index เดิมซ้ำ ๆ
   นั่นคือสัญญาณว่าคนเขียน (หรือ batch เขียนคำถามอัตโนมัติ) มีอคติเชิงระบบบางอย่าง ซึ่งอาจไปโผล่เป็นปัญหา
   อื่นที่ shuffle ช่วยไม่ได้ (เช่น distractor ที่ตำแหน่งเดิมมักจะสั้นกว่า/ยาวกว่าเสมอ ซึ่งกลายเป็น length
   bias ที่ shuffle ช่วยไม่ได้เลย) เกทจึงเช็คไว้ก่อนแม้ position เองจะไม่กระทบผู้ใช้โดยตรงในแอปนี้ —
   การรายงานแยกป้ายกำกับ "reaches the user" vs "content-quality signal only" ในเอาต์พุตคือกลไกที่กัน
   ไม่ให้คนอ่าน over-react กับตัวเลข position (ไม่กระทบผู้ใช้จริง) หรือ under-react กับ length (กระทบจริง)
2. เพราะ `LoadLesson` decode ด้วย `DisallowUnknownFields()` แล้วสร้าง domain `Lesson` เต็มรูปแบบทันที
   ซึ่งจะ reject ไฟล์ที่ `expected_answer` ไม่อยู่ใน options ของตัวเองไปตั้งแต่ decode (ผ่าน
   `NewRecallCheck`'s `slices.Contains` check) — ก่อนที่ `internal/mcqguess` จะได้เห็นคำถามนั้นด้วยซ้ำ
   ทำให้ tool นี้ไม่มีทางเลือกที่จะ "abort พร้อม error ที่ระบุ path + คำตอบ" ตามที่ตัดสินใจไว้ (จะได้แค่
   error จาก domain แทน ซึ่งไม่ใช่ scope ของ tool วัด) แถม `LoadLesson` ยังต้องการ topic/concept_id ที่
   ตรงกับโฟลเดอร์/ชื่อไฟล์ (เช็คสำหรับ import จริงที่ต้อง resolve concept row ใน MySQL) ซึ่ง tool วัดนี้
   ไม่ต้องมี DB เลยตาม ticket ระบุ
3. เพราะ percentage points ไม่ portable ข้าม corpus ที่ baseline ต่างกัน — excess 5pp เหนือ baseline
   25% (corpus AWS 4 ตัวเลือก) คือ relative jump 20% แต่ 5pp เดียวกันเหนือ baseline 33.3% (corpus เดิม
   3 ตัวเลือก) คือแค่ relative jump 15% เท่านั้น ถ้าใช้ threshold เดียวกันเป็น pp ตรง ๆ ข้ามสอง corpus
   นี้ เกทจะเข้มกว่าจริงกับ corpus AWS โดยไม่ได้ตั้งใจ (หรือหลวมกว่าจริงกับ corpus เดิม) ซึ่งเป็นบั๊ก
   ประเภทเดียวกับที่ ticket ทั้งใบนี้มีไว้แก้ (เทียบ raw hit rate ข้าม corpus ที่ baseline ไม่เท่ากัน)
   เพียงแต่เกิดขึ้นที่ตัว metric ของ excess เอง ไม่ใช่ที่ตัว baseline

</details>
</details>

### Round 2 (code review) — สรุปสิ่งที่แก้

Round 1 ของ code review เป็น REQUEST_CHANGES กว้าง (R1–R7 + 3 "Also required") — สาระสำคัญของทุกจุด
คือ**เชื่อคำกล่าวอ้างโดยไม่ re-derive** (สองจุดในเอกสาร) กับ**เทสต์ที่พิสูจน์น้อยกว่าที่ comment อ้าง**
(หลายจุดในโค้ด) แก้ครบทุกข้อ:

- **R1**: 89.6% เป็นค่าวัดซ้ำได้จริง (ไม่ใช่ "วัดไม่ได้โดยหลักการ" อย่างที่เขียนไว้เดิม) — กู้คืน corpus
  ก่อน PR #43 ด้วย `git archive 6b1a765`, วัดจริงได้ 36.8% (longest) ไม่ใช่ 89.6% **(round 3 พบว่า
  89.6% คือค่า middle-length ไม่ใช่ longest — ดูหัวข้อ "Finding" ด้านบนสำหรับคำตอบเต็ม)**
- **R2**: เหตุผล "tie-break rule เป็นไปได้มากที่สุด" สำหรับช่องว่าง 33.7%→32.6% ถูกแทนที่ด้วยขอบเขต
  ที่คำนวณจริง (30.90%–34.27%) และ derive ว่า 33.7% ต้องการ 10/12 tie ชนะ (p≈0.85%) — สรุปว่า tie-break
  ไม่ใช่คำตอบ **(round 3 พบคำตอบจริง: 33.7% คือค่า middle-length เช่นกัน ดูหัวข้อ "Finding" ด้านบน)**
- **R3**: `EvaluateGate`'s position-heuristic loop ไม่มีเทสต์ที่พิสูจน์ว่ามันทำงานจริง (ลบทั้ง loop
  ออกแล้ว `go test` ยังเขียว) — เพิ่ม `TestEvaluateGate_PositionAloneCanFailWithLongestClean`
- **R4**: comment ของ `TestMeasure_CorpusStraddlingGate_ExitCodeFlips` (และที่ echo ใน
  `cmd/mcq-guessability/main_test.go` กับตาราง synthetic-corpus ด้านบน) อ้างผิดว่า position heuristic
  "อยู่ใต้ threshold สบาย ๆ" ทั้งที่ position index 3 ข้าม threshold พร้อม longest heuristic พอดี
  (เพราะ options ถูกสร้างให้ index 3 ยาวที่สุดเสมอ — "guess longest" กับ "guess index 3" จึงเป็นกลยุทธ์
  เดียวกันในฟิกซ์เจอร์นี้) แก้ comment ให้ตรงความจริง + เพิ่ม assertion นับจำนวน violation (ต้องได้ 2
  ไม่ใช่ 1) ทั้งสามจุด
- **R5**: mcq ที่มี 1 ตัวเลือกผ่านเช็คทั้งหมดแบบไม่มีความหมาย แล้วดึง `AvgBaseline` ขึ้นแบบเงียบ ๆ —
  เพิ่ม `minMeasurableOptions = 2` ใน `Measure` (mirror domain's `minMCQOptions`) + abort พร้อม error
  เหมือน `expected_answer` ไม่อยู่ใน options
- **R6**: `LoadQuestions` ไม่มีเทสต์ที่พิสูจน์ว่า decode error จริงจะ surface เป็น error (เปลี่ยนเป็น
  `return nil, nil` แล้ว `go test` ยังเขียว) — เพิ่ม `TestLoadQuestions_MalformedJSONErrors`
- **R7**: ลบ WHAT-comment ที่ restate signature/return ซ้ำใน `HitRate`, `ExpectedIndex`, `Baseline`,
  `ExitCode` — เก็บเฉพาะ WHY (invariant, เหตุผลของค่าคงที่, ทำไม unreachable case ถึง return 0)
- **Per-track gating (`EvaluateTrackGates`)**: gate เดิมเช็คแค่ pooled overall ซึ่ง dilute ได้ (track
  แย่เล็ก + track ดีใหญ่ = ตัวเลข pooled ที่ดูดีกว่าจริง) — เปลี่ยนเป็น fail ถ้า track ใด track หนึ่ง
  fail ยืนยันด้วยตัวเลขที่ derive เอง (ไม่ใช่แค่ก็อปของ reviewer): fair 300 ข้อ excess 0% + bad 100 ข้อ
  excess 60% → pooled 15% (ผ่าน threshold 20% ทั้งที่ track แย่ควร fail)
- **`MinSampleSize = 100`**: เพิ่มขั้นต่ำจำนวนตัวอย่างก่อน gate ตัดสิน (ต่ำกว่านี้รายงานเป็น
  "insufficient" แทนที่จะ fail/pass) — derive false-positive rate เอง ที่ p=0.25 ตลอดตาราง
  (19.7%/6.9%/0.05% ที่ n=30/100/550, สองวิธีคำนวณตรงกัน) **ต่างจากตัวเลขที่ reviewer อ้างไว้ที่ n=30
  (23.9%)** ตอนแรกดูเหมือนตรวจสอบซ้ำไม่ได้ **(round 3 หาสาเหตุเจอ: 23.9% คือ p=0.20/1/5 ไม่ใช่ p=0.25 —
  ดู "Round 3" ด้านล่าง — เก็บตารางของเอกสารนี้ไว้ที่ p=0.25 ตลอดตามเดิม เพราะเป็นโมเดลเดียวที่
  self-consistent สำหรับค่าคงที่ตัวเดียว)**
- **Scope ของเกท 25% แคบลงเหลือ single-answer mcq**: `Question.ExpectedAnswer` เป็น string เดี่ยว —
  AWS-S2's "Select TWO" ยังวัดไม่ได้ด้วย tool นี้ (ไม่ error แต่ก็ไม่ถูกต้อง) บันทึกไว้ใน package doc
  ของ `cmd/mcq-guessability/main.go` ว่าต้องขยาย `Question`/`ExpectedIndex` ก่อนเกทจะครอบคลุม
  AWS-C1..C4 เต็มรูปแบบ
- **`longestOptionIndex` unexport + `Baseline()` guard คืน 0**: ปิดช่องที่ caller ภายนอกเรียกด้วย slice
  ว่างแล้ว panic (`longestOptionIndex`) หรือได้ +Inf แบบเงียบ ๆ (`Baseline()`)

### Test count (round 2)

- Go: **301 (develop) → 328 (round 1) → 337 (round 2)** (นับจาก `go test -v` ผ่าน `grep -c "^--- PASS"`
  — round 2 เพิ่ม 9 เทสต์ใหม่: 1 ที่ R3, 1 ที่ R5, 1 ที่ R6, 2 ที่ MinSampleSize, 4 ที่
  `EvaluateTrackGates`/per-track dilution)
- `go vet ./...`, `gofmt -l .` clean ทั้งคู่
- Mutation: 11 mutation (9 เดิม + 2 ที่ round 1 ของ review พบว่ารอด) ตายจริงทุกตัวหลัง round 2 — ตาราง
  เต็มอยู่ด้านบน

### Review focus (round 2 — เน้นจุดที่แก้ตาม code review)

<details>
<summary>คำถามสำหรับรีวิว diff รอบนี้ (เฉลยพับไว้ด้านล่าง)</summary>

1. ทำไม `TestEvaluateGate_PositionAloneCanFailWithLongestClean` ถึงจำเป็น ทั้งที่มี
   `TestMeasure_CorpusStraddlingGate_ExitCodeFlips` ที่ทำให้ position heuristic fail อยู่แล้ว?
2. ทำไม `EvaluateTrackGates` ถึงไม่รับ overall `Report` (ตัว pooled) เป็น input เลย แทนที่จะรับแล้ว
   เลือกไม่ใช้มันตัดสิน?
3. ทำไมตัวเลข false-positive rate ของ `MinSampleSize` ที่บันทึกในเอกสารนี้ (19.7% ที่ n=30) ถึงต่างจาก
   ตัวเลขที่ reviewer อ้างตอน request changes (23.9%) — และทำไมเอกสารนี้ถึงยังไม่เปลี่ยนไปใช้ 23.9%?

<details>
<summary>เฉลย</summary>

1. เพราะฟิกซ์เจอร์ของ `TestMeasure_CorpusStraddlingGate_ExitCodeFlips` สร้าง options ให้ index 3 ยาว
   ที่สุดเสมอ ทำให้ "guess longest" (longest-option heuristic) กับ "guess index 3" (position
   heuristic) เป็น**กลยุทธ์เดียวกันเป๊ะ** ในฟิกซ์เจอร์นั้น — longest-option heuristic กับ position[3]
   จึงข้าม threshold พร้อมกันเสมอโดยไม่มีทางแยกว่า mutation ที่ลบ position-loop ทั้งก้อนออกจะถูกจับได้
   จาก longest-option heuristic เพียงอย่างเดียวหรือเปล่า (คือถ้าลบ position loop ออกจริง เทสต์นี้ก็ยัง
   fail ได้จาก longest ตามปกติ ซึ่งพิสูจน์ไม่ได้ว่า position loop มีส่วนจริง) ต้องมีฟิกซ์เจอร์ที่ longest
   "สะอาด" (excess 0) แต่ position ตัวเดียว "แย่" โดยเฉพาะ ถึงจะพิสูจน์ได้ว่า position loop ทำงานเป็น
   อิสระจาก longest จริง ๆ
2. เพราะการมี parameter ให้เลือก "จะ gate ด้วย pooled หรือไม่" จะเปิดช่องให้ future caller เลือกกลับไปใช้
   pooled ได้ ซึ่งเป็นบั๊กเดิมที่ round นี้กำลังแก้ (dilution) — การไม่รับ `Report` ตัว pooled เข้ามาเลย
   (รับแค่ `map[string]Report`) ทำให้การ gate บน pooled figure เป็นไปไม่ได้เชิงโครงสร้างที่ signature
   ระดับ type ไม่ใช่แค่ระดับ convention/comment ที่ลืมทำตามได้
3. **(แก้ตอน round 3)** เพราะ 19.7% กับ 23.9% เป็นคำตอบที่ถูกทั้งคู่ของ**คนละคำถาม**: 19.7% คือ n=30
   ที่ p=baseline=0.25 (corpus 4 ตัวเลือก, ตรงกับตารางหลักของเอกสารนี้ตลอด) ส่วน 23.9% คือ n=30 ที่
   p=0.20 (1/5) ซึ่ง reviewer ตั้งใจใช้แสดงตัวอย่างของ `Position[4]` บน corpus 5 ตัวเลือก
   (multi-response) — ตารางต้นฉบับของ reviewer เองผสม baseline ข้ามแถวโดยไม่ระบุ (p=0.20 ที่ n=5/n=30,
   p=0.25 ที่แถวอื่น) จึงดูเหมือนขัดกันจนกว่าจะสืบจนเจอว่าเป็นคนละ p เอกสารนี้ยังคงใช้ p=0.25 ตลอดทั้ง
   ตาราง เพราะ `MinSampleSize` เป็นค่าคงที่ตัวเดียวที่ใช้กับทุก heuristic (longest/shortest/middle/
   position ทุก index) จึงต้องมีโมเดลเดียวที่ self-consistent ไม่ใช่ผสม p หลายค่าตามบริบทของแต่ละแถว —
   บทเรียนคือตัวเลขสองตัวที่ดูขัดกันไม่ได้แปลว่าตัวใดตัวหนึ่งผิดเสมอไป บางครั้งทั้งคู่ถูกแต่ตอบคนละคำถาม
   ต้อง re-derive จนเจอว่าต่างกันตรงไหนจริง ๆ ก่อนสรุปว่า "ไม่มีคำอธิบาย"

</details>
</details>

### Round 3 (code review) — คลี่ปม 89.6%/33.7% + เพิ่ม shortest/middle heuristic + ปิดช่อง gate ที่ไม่ได้ตัดสิน

Round 2 ของ code review เป็น **NO-SHIP สามข้อ** (N1–N3): N3 คือตัวที่คลี่ปมทั้งสองข้อที่ round 1-2 เขียน
ไว้ว่า "unexplained"/"provenance unknown" — คำตอบอยู่ใน `docs/tickets/mcq-quality.md`'s "บทเรียนที่ 1"
มาตั้งแต่ต้น เพียงแค่ไม่เคยเอามาทดสอบกับ tool จริง แก้ครบทุกข้อ:

- **N3**: 89.6% และ 33.7% (เดิม) คือค่า **middle-length** heuristic ไม่ใช่ longest ที่ tool วัดมาตลอด
  round 1-2 — เพิ่ม `Report.Middle`/`middleOptionIndex` แล้ววัดซ้ำ reproduce ตรงเป๊ะทั้งคู่ (319/356 ที่
  `7dfb0ba`, 120/356 ที่ปัจจุบัน) รายละเอียดเต็มอยู่ที่หัวข้อ "Finding" ด้านบน (แทนที่ข้อความเดิมทั้งหมด
  ไม่ใช่แค่เพิ่มเติม)
- **N2**: เพิ่ม `shortestOptionIndex` (mirror ของ `longestOptionIndex`) และ `middleOptionIndex` (ใหม่)
  พร้อม synthetic test ครบ (100%, boundary 2 ตัวเลือก, ตัวเลือกยาวชนกันเหลือ 2 ค่า, 4+ ตัวเลือกมีตัวกลาง
  หลายตัว) — เกทเช็คทั้งสามตัวพร้อม position รายละเอียดการตัดสินใจเรื่อง "นับ miss vs exclude" อยู่ที่
  หัวข้อ "N2" ด้านบน (จุดตัดสินใจนี้คือสิ่งที่ทำให้ reproduce เลข 89.6%/33.7% ได้ตรงเป๊ะ ลองแบบ exclude
  ก่อนแล้วไม่ reproduce)
- **N1**: `GateResult` เพิ่ม `Requested`/`Judged` ทำให้แยก "ตัดสินแล้วผ่าน" กับ "ไม่ได้ตัดสินอะไรเลย" ได้
  — `ExitCode` คืน 1 เมื่อ requested แต่ `Judged==0`, `run()` พิมพ์ `Gate: NOT JUDGED` แยกจาก `PASS`/
  `FAIL` รายละเอียดอยู่ที่หัวข้อ "N1" ด้านบน
- **Also fix — `MinSampleSize`'s doc comment ตัวอย่างผิด**: comment เดิมอ้างว่า n=550 "จะกัน
  domain-driven-design (11 mcqs) ไม่ให้ถูก gate" แต่ n=100 ก็กัน DDD เหมือนกัน (11 < ทั้งคู่) — ไม่ใช่
  ตัวอย่างที่แยกความต่างระหว่าง 100 กับ 550 ได้จริง แก้เป็นตัวอย่างที่ถูกต้อง: `ai-and-llm-systems`
  (158 mcqs) และ `designing-data-intensive-applications` (187 mcqs) ผ่าน 100 วันนี้แต่จะไม่ผ่าน 550 —
  นี่คือคู่ที่แสดงความต่างจริงระหว่างสองค่า
- **Also fix — ประโยคที่ขัดกันเองในวงเล็บ**: ข้อความเดิม ("เกินขอบบนของทุก correctness-blind rule (122
  คือ ceiling จริง แต่ 120 อยู่ในช่วง [floor, ceiling] พอดี ไม่ได้เกินขอบ...")` เปิดด้วยข้อความหนึ่งแล้ว
  ถอนคำพูดในวงเล็บถัดไปทันที — เขียนใหม่เป็นประโยคเดียวที่สอดคล้องกันในหัวข้อ "Finding" ด้านบน (ไม่มีการ
  ขัดแย้งกันเองอีกแล้ว)
- **Also fix — `ExitCode`'s comment**: comment เดิม ("maps a gate verdict to a process exit code, 1 on
  failure and 0 otherwise") restate ตัว function 5 บรรทัดตรง ๆ ไม่มี WHY เลย (WHY จริงถูกย้ายไปที่
  `-max-excess`'s flag help ตั้งแต่ round 2 แล้ว) — ลบทิ้ง เหตุผลของ tri-state (`Requested`/`Judged`)
  ย้ายไปอยู่ที่ `GateResult`'s doc comment แทน ซึ่งเป็นที่ที่เหมาะสมกว่า (อธิบาย field ที่ตัดสินใจ ไม่ใช่
  ฟังก์ชันที่แค่ map ผลไปเป็นตัวเลข)
- **Also fix — package ย้ายจริง**: round 2 ระบุใน "Also required" ว่าต้องย้าย `internal/mcqguess` ไป
  `cmd/mcq-guessability/internal/mcqguess` แต่รายงานสรุปตอนนั้นไม่ได้ระบุชัดว่าทำหรือไม่ทำ (บอกแค่ "all
  implemented" แบบรวม ๆ) — round 3 ตรวจสอบแล้วพบว่า**ยังไม่ได้ทำจริง** ย้ายด้วย `git mv` (คง history) จริง
  ในรอบนี้ อัปเดต import path ใน `cmd/mcq-guessability/main.go` และ reference ในเอกสารทุกจุด

### Test count (round 3)

- Go: **301 (develop) → 328 (round 1) → 337 (round 2) → 343 (round 3)** (นับจาก `go test -v` ผ่าน
  `grep -c "^--- PASS"` — round 3 เพิ่ม 6 เทสต์ใหม่: `TestShortestOptionIndex`, `TestMiddleOptionIndex`,
  `TestMeasure_CorrectAnswerAlwaysShortest_ShortestHeuristicIsHundredPercent`,
  `TestMeasure_CorrectAnswerAlwaysMiddle_MiddleHeuristicIsHundredPercent`,
  `TestMeasure_MiddleOptionIndex_BoundaryCases`,
  `TestRun_NotJudgedWhenEveryHeuristicIsBelowMinSampleSize`)
- `go vet ./...`, `gofmt -l .` clean ทั้งคู่
- Mutation: 17 mutation (11 เดิม + 6 ใหม่จาก N1/N2) ตายจริงทุกตัวหลัง round 3 — ตารางเต็มอยู่ด้านบน,
  ยืนยันซ้ำทั้ง 11 เดิมด้วยว่ายังตายอยู่หลังการรีแฟกเตอร์ของรอบนี้

### สี่คอร์ปัส สามเฮอริสติก (สำหรับ reviewer ตรวจตรง)

คำสั่ง reproduce: current = `go run ./cmd/mcq-guessability -dir content/lessons`; historical =
`git archive <commit> content/lessons | tar -x -C <tmp> && go run ./cmd/mcq-guessability -dir
<tmp>/content/lessons`

| Commit | Longest | Shortest | Middle | Position (0/1/2) |
|---|---|---|---|---|
| `6b1a765` (ก่อน #38 หนึ่งขั้น) | 36.8% | 28.4% | 29.5% | 41.6% / 30.3% / 28.1% |
| `7dfb0ba` (หลัง #38 "C-mcq-sweep") | 2.0% | 3.7% | **89.6%** | 41.6% / 30.3% / 28.1% |
| ปัจจุบัน / `057ec9a` (หลัง #43 "C-mcq-balance") | 32.6% | 31.5% | **33.7%** | 33.4% / 33.4% / 33.1% |

(baseline 33.3% ทุก cell, corpus 3-option ทั้งหมด ทั้งสามคอมมิต — position ที่ `6b1a765` กับ `7dfb0ba`
เท่ากันเป๊ะเพราะ #38 แก้แค่เนื้อหาตัวเลือก ไม่ได้ย้ายตำแหน่งคำตอบ)

### Review focus (round 3)

<details>
<summary>คำถามสำหรับรีวิว diff รอบนี้ (เฉลยพับไว้ด้านล่าง)</summary>

1. ทำไมการเปลี่ยน `Middle` จาก "exclude คำถามที่ไม่มีตัวกลาง" เป็น "นับเป็น miss" ถึงเป็นจุดตัดสินใจที่
   สำคัญที่สุดของรอบนี้ ทั้งที่ดูเหมือนรายละเอียดเล็ก ๆ?
2. ทำไม 41.6% (position) ตรงกับตัวเลขเดิมตั้งแต่ round 1 แต่ 89.6% (length) ต้องรอถึง round 3 ถึงจะ
   reproduce ได้?
3. ทำไม `ExitCode` ต้องคืน 1 เมื่อ `Judged == 0` ทั้งที่ `Failed` ยังเป็น `false` อยู่?

<details>
<summary>เฉลย</summary>

1. เพราะมันคือความต่างระหว่าง "89.6%/33.7% reproduce ได้" กับ "ยังเป็นปริศนาต่อไป" ล้วน ๆ — ลองวัดแบบ
   exclude ก่อน (เหมือน `Position[k]` ทำกับคำถามที่มีตัวเลือกไม่พอ) ได้ 120/328=36.6% ไม่ตรงกับ 33.7%
   ที่บันทึกไว้เลย เปลี่ยนเป็นนับ miss (Total เท่ากับ NumMCQs เสมอ เหมือน Longest/Shortest) ได้
   120/356=33.7% ตรงเป๊ะทันที เหตุผลเชิงแนวคิดก็สอดคล้อง: "guess index 3" ไม่มีความหมายสำหรับคำถาม
   2 ตัวเลือก (ไม่มี index 3 อยู่จริง) แต่ "เดาตัวเลือกที่ไม่สุดขั้ว" มีความหมายเสมอ แค่บังเอิญเดาไม่ถูก
   แน่นอนเมื่อทุกตัวเลือกเป็นสุดขั้วหมด — การ exclude แทนที่จะนับ miss จะซ่อน tell จริงไว้หลังตัวหารที่
   หดตัวได้ ซึ่งเป็นกลไกเดียวกับที่ทำให้เหตุการณ์ 89.6% เดิมหลุดรอดมาได้ตั้งแต่ต้น (วัดผิดจุด ไม่ใช่วัดผิด
   วิธี)
2. เพราะ position heuristic ของ tool เดิมกับ tool นี้วัดสิ่งเดียวกันมาตลอด (เดาตำแหน่งเดิมเสมอ ไม่มี
   ทางตีความสองแบบ) จึงตรงกันตั้งแต่ round 1 แต่ "length heuristic" ของ tool เดิม (script ที่สูญหายไป
   แล้ว) จริง ๆ วัด **middle**-length ไม่ใช่ **longest**-length ที่ tool นี้ implement มาตลอดจนถึง
   round 2 — เครื่องมือทั้งสองวัดถูกต้องในสิ่งที่แต่ละตัวถูกออกแบบมาให้วัด (ไม่มีบั๊กในทั้งคู่) แค่ตั้งชื่อ
   heuristic เดียวกัน ("length") ให้กับสองสิ่งที่ต่างกัน — ปมนี้แก้ได้ก็ต่อเมื่อมี middle heuristic ให้
   เทียบเท่านั้น ซึ่งเป็นสิ่งที่ N2 เพิ่มเข้ามาในรอบเดียวกันพอดี
3. เพราะ `Judged == 0` แปลว่าทุก heuristic ถูกข้ามเพราะ n ต่ำกว่า `MinSampleSize` — ไม่มีการเปรียบเทียบ
   กับ threshold เกิดขึ้นเลยสักครั้ง `Failed == false` ในสถานะนี้ไม่ได้แปลว่า "ตรวจแล้วผ่าน" แต่แปลว่า
   "ยังไม่ได้ตรวจ" ถ้า `ExitCode` คืน 0 ในสถานะนี้ CI จะอ่านเป็น PASS ทั้งที่ไม่มีการวัดอะไรเกิดขึ้นจริง —
   อันตรายที่สุดคือ batch แรกของทุก domain ในแผนเขียนคำถามทีละ domain จะตกอยู่ในสถานะนี้พอดี (ยังไม่ถึง
   100 ข้อ) ซึ่งเป็นจังหวะที่ควรจับ bias เชิงระบบให้ได้มากที่สุด ไม่ใช่จังหวะที่เกทเงียบไปเฉย ๆ

</details>
</details>

## Template — decisions pinned in `.claude/agents/lesson-writer.md` (AWS-C0)

รายละเอียดเต็มอยู่ใน `.claude/agents/lesson-writer.md`'s "## AWS track" section (ปักเป็น agent spec
โดยตรง ไม่ใช่แค่บันทึกไว้ที่นี่) — สรุป 9 การตัดสินใจที่มาจากการรีวิว pilot 2 บทเรียน
(`vpc-fundamentals`, `s3-security`, ทั้งคู่ PASS lesson-verifier) **หลังแก้ไข round 2 ของ code review**
(ตัวเลขและการอ้างอิงในหัวข้อนี้แก้ตามที่ round 2 ชี้ว่าผิด — ดูหัวข้อ "AWS-C0" ด้านล่างสำหรับรายละเอียด
ของทุกจุดที่แก้):

1. **Prose-rune ceiling ≤ 950 runes/min** (ไม่ใช่ band 1,150–1,200/min แบบ round 1) — "prose runes" คือ
   `body_md` หลังตัด fenced block ทุกชนิด (รวม mermaid) และแถวตาราง markdown (บรรทัดที่ขึ้นต้นด้วย `|`)
   ออก แล้วหารด้วย `est_minutes` — **บรรทัดว่างนับรวมอยู่ด้วย** (นี่คือจุดเดียวที่กำกวม เลือก convention
   นี้ทุกครั้งที่ re-derive เพื่อให้ตัวเลขทำซ้ำได้) **Round 1 วัดผิดตัวชี้วัด**: นับ rune ทั้งไฟล์รวม
   mermaid/table/JSON ทำให้ตั้ง band ที่ไม่มีบทเรียนไหนในคลัง 120 ใบเดิมเข้าเกณฑ์เลยนอกจาก 2 pilot เอง
   (ซึ่งทั้งคู่ยังเกินค่าสูงสุดของ 118 ใบที่เหลือด้วยซ้ำ) — วัดใหม่แบบ prose-only แล้ว derive ceiling จาก
   corpus จริงทั้ง **120 ใบ** (ไม่ใช่ 118 ใบที่ไม่รวม AWS — round 2 ผสม population สองชุดผิด): mean **681**
   + 3 standard deviation (**78**) = **915** ซึ่งอยู่ใกล้ค่าสูงสุดจริงในคลังมาก
   (`ai-and-llm-systems/online-eval-and-ab-testing` วัดได้ **918** ด้วย convention เดียวกันนี้ — ใกล้กัน
   ไม่ใช่ตัวเลขเดียวกันเป๊ะ ไม่ได้บังคับให้เท่ากัน) — ปัดขึ้นเป็น 950 ให้มี margin เหนือทั้งสองค่า ไม่ใช่
   ตัดคลังเดิมทิ้ง ที่ ceiling นี้ ทั้งสอง pilot วัดได้ **s3-security 830/min, vpc-fundamentals 838/min**
   (หลัง retrofit จุด "." ท้ายประโยคตาม Template ข้อ 10 — ดู round 3) ไม่ใช่เหนือค่าสูงสุดของคลังอีกต่อไป
   — **ทั้งคู่ไม่เคยเกินจริง ไม่ต้องตัด/แยกไฟล์**
   **เกทนี้ reformat แล้วโกงไม่ได้**: ตาราง (รวม cue section) ยกเว้นจากการนับเพราะเป็นเนื้อหาที่ผู้อ่าน
   *scan* ไม่ได้อ่านเรียงบรรทัด ไม่ใช่ช่องโหว่ให้แปลงร้อยแก้วเป็น pseudo-table row เพื่อลดตัวเลข — การแปลง
   เนื้อหาจริงเป็นตารางปลอมเพื่อซื้อ budget ทำลายจุดประสงค์ของเกทนี้ ต่อให้ตัวเลขผ่าน
2. **Explanation format บังคับใช้ bold header ภาษาไทย** — **โจทย์ถามว่า** / **ทำไมข้อที่ถูกถึงถูก** /
   **ทำไมตัวอื่นผิด** (ตามด้วย bullet `- *ชื่อตัวเลือก* — เหตุผล` ทีละตัว) / **Decision rule** — ตอนนี้
   ใช้ครบทั้ง 4 ข้อของทั้งสอง lesson แล้ว (round 1 ปล่อยให้ `s3-security` ทั้ง 4 ข้อยังเป็นร้อยแก้ว
   inline (ก)(ข)(ค) ทั้งที่กฎนี้ระบุไว้แล้ว — round 2 แปลงให้ครบ)
3. **"คำในโจทย์ → คำตอบที่ต้องมองก่อน" เป็น section บังคับ และปักรูปแบบเป็นตาราง**
   (`| คำในโจทย์ | คำตอบที่ต้องมองก่อน |`) ไม่ใช่ bullet list — round 1 ปล่อยให้ `vpc-fundamentals` ใช้
   bullet list ใต้หัว "## มุมข้อสอบ" คนละชื่อคนละรูปแบบกับ `s3-security` round 2 แปลง
   `vpc-fundamentals` เป็นตารางหัวเดียวกันแล้ว
4. **Distractor policy**: distractor ใช้ real AWS service ที่บทเรียนไม่ได้สอนได้ ถ้าคำตอบถูกยังเลือก
   ได้ครบจากเนื้อหาบทเรียนอย่างเดียว (pilot ใช้ S3 Transfer Acceleration, IAM Access Analyzer,
   presigned URL เป็น distractor ทั้งที่ไม่ได้สอนในบท) — ตรงกับกฎ "real service misapplied" ใน
   `aws-cert.md` เดิม บันทึกไว้ให้ชัดกันคนรีวิวถัดไป flag ซ้ำ
5. **อย่างมากแค่ 1 ใน 4–6 คำตอบถูกของบทเรียนหนึ่งใบเป็นตัวเลือกที่ยาวที่สุด** — วัดจริงจาก pilot คือ
   **2 ใน 8 ข้อ** (ข้อละ 1 ต่อ lesson) ไม่ใช่ 3 ใน 8 ตามที่ round 1 เขียนผิด ทั้งสองข้อยาวเพราะต้องระบุ
   เงื่อนไข compound ให้ครบ ("compliance mode ... including the root user") ซึ่งเป็น soft tell ที่
   guessability metric รวมไม่จับที่ n ต่ำ — กฎนี้เช็คได้ทีละ lesson ตอนเขียน (นับก่อน save) ไม่ใช่แค่
   "บางครั้ง" แบบ round 1 ที่ไม่มีทางวัด
6. **Recall checks: 4–6 ข้อ/concept** ไม่ใช่ ~11 — ดูหัวข้อ "รูปทรงใหม่" ต้นเอกสารสำหรับเหตุผลเต็ม
7. **ทุก fact ต้องมี URL ที่ fetch จริง** ข้อที่ verify ไม่ได้ = ตัดทิ้งแล้วรายงาน ไม่ใช่เดา
8. **Quote ต้อง verbatim หรือไม่ใส่เครื่องหมายคำพูดเลย** — ถ้า quote เป็นบางส่วนของประโยคยาว ต้องขึ้นต้น
   ด้วย `...` ให้รู้ว่าเป็น fragment และคงตัวเน้น (bold) ของต้นฉบับไว้ในคำพูด (round 1 ตัดคำว่า **and**
   ที่ตัวหนาออกทั้งที่เป็นคำที่รับน้ำหนักประโยค — round 2 คืนกลับพร้อมใส่ `...` นำหน้า)
9. **Position balance (longest/shortest/middle/index) เป็นตัวเลขที่วัดหลังเขียน ไม่ใช่เป้าที่ต้องเขียน
   ให้ตรง** — pilot บังเอิญลง 2/2/2/2 ทุก index แต่ไม่ได้ตั้งใจ `cmd/mcq-guessability` เองระบุว่า
   position เป็น "content-quality signal only, not user-facing" เพราะ `web/lib/shuffle.ts` สลับ
   ตำแหน่งใหม่ตอน render อยู่แล้ว
10. **Thai sentence จบด้วย "." เสมอ** — วัดจริงจาก `designing-data-intensive-applications` และ
    `ai-and-llm-systems` ทั้งคู่ลง "." ในสัดส่วนข้างมาก (รวมกันกว่า 80%) — นี่คือ track ที่กติกาข้อนี้
    มาจากจริง ๆ ส่วน `domain-driven-design` **โหวตสวนทาง** (ใช้ "." เป็นส่วนน้อย) ไม่ใช่ที่มาของ
    convention — round 2 เคยอ้าง `domain-driven-design` เป็นแหล่งด้วยผิด แก้แล้วในรอบนี้ pilot ทั้งสองใบ
    อยู่ต่ำกว่า majority convention ชัดเจนก่อนรอบนี้ (`s3-security` ไม่มีเลย, `vpc-fundamentals` มีแค่
    ส่วนน้อย) — **retrofit ใส่ "." ย้อนหลังให้ครบทั้งสองไฟล์แล้วในรอบนี้** (ดูหัวข้อ "AWS-C0" round 3)

## AWS-C0 — land the pilot lessons and lock the lesson template

**สถานะ**: implemented, code review round 3 fixes applied, pending re-review

**สิ่งที่ทำ (round 1)**: trim `s3-security.json`'s `body_md` ตัดเนื้อหาซ้ำ (S3 Bucket Keys
re-explained ใน fintech example, 2 bullet ซ้ำใน trade-off list), แก้ quote ที่ verifier หาต้นฉบับไม่
เจอ, เพิ่ม 7 currency-checklist item, ปักเทมเพลตใน `lesson-writer.md`, แก้ count contradiction
(~11/550 → 4–6/~250) — รายละเอียดเต็มของ round 1 อยู่ในประวัติ commit (`git log -p` ของ commit แรก
ของ ticket นี้), ไม่ทวนซ้ำที่นี่เพราะ round 2 แก้ตัวเลขและเนื้อหาหลายจุดที่ round 1 อ้างไว้ผิด

### Round 2 (code review) — 7 findings (R1–R7), ทุกข้อแก้แล้ว

Round 1 ได้ REQUEST_CHANGES: R2 (rune budget) เป็นต้นเหตุของ R3 (การตัดเนื้อหาตามเป้าที่ผิดทำให้
ประโยคเสียหาย 3 จุด, จุดหนึ่งสอนผิด) — แก้ R2 ก่อนเพราะเป็นราก แล้วค่อยคืนเนื้อหาที่ R3 เสียหาย:

- **R2 — rune budget วัดผิดตัวชี้วัดและใช้ band ที่ unsatisfiable**: นับ rune ทั้งไฟล์ (รวม mermaid,
  table, JSON fence) ทำให้ band 1,150–1,200/min ตัดคลังเดิมทั้ง 118 ใบทิ้งหมด (มีแค่ 2 pilot เข้าเกณฑ์
  และทั้งคู่ยังเกิน max ของ 118 ใบที่เหลือด้วยซ้ำ) — แก้เป็น prose-only (ตัด fenced block + table row
  ออกก่อนนับ, **นับบรรทัดว่างรวมด้วย** — เดิมเป็นจุดกำกวมที่ไม่ได้ระบุ) และเป็น **ceiling 950/min** ไม่ใช่
  band, derive จาก corpus จริงทั้ง **120 ใบ**: mean **681** + 3σ (**78**) = **915** ใกล้ค่าสูงสุดจริงของ
  คลังมาก (`ai-and-llm-systems/online-eval-and-ab-testing` วัดได้ **918** ด้วย convention เดียวกัน —
  ใกล้กันไม่ใช่เท่ากันเป๊ะ) — ที่ ceiling ใหม่นี้ `s3-security` วัดได้ 830/min, `vpc-fundamentals` 838/min
  (หลัง retrofit period ท้ายประโยคใน round 3) ทั้งคู่**ไม่เคยเกินจริง** ไม่ต้องแยกไฟล์ — ดูสูตรเต็มที่
  หัวข้อ "Template" ข้อ 1 ด้านบน
  (**round 3 แก้ arithmetic**: round 2 เขียน mean 678 + 3σ(78) ≈ 909 ผสม population สองชุดผิด — 678
  คือ mean ของ 118 ใบที่ไม่รวม AWS ด้วย convention ตัดบรรทัดว่างทิ้ง, 78 คือ σ ของ 120 ใบ และตัวเลขสองตัว
  นี้เองก็รวมกันได้ 912 ไม่ใช่ 909 อยู่ดี — เป็นทั้งการผสม population ผิดและ arithmetic ผิดซ้อนกันสองชั้น
  วิธีที่ทำให้ arithmetic กลับมาถูก (mean 681 + 3σ(78) = 915) คือเปลี่ยนไปนับบรรทัดว่างรวมด้วย ไม่ใช่ตัด
  ทิ้ง — เพราะเหตุนี้ convention ของบรรทัดว่างจึงต้องระบุชัดและใช้ให้เหมือนกันทั้งตอนวัด lesson เดี่ยว ๆ
  และตอน derive corpus statistic)
- **R3 — การตัดตามเป้าที่ผิดทำให้ 3 ประโยคเสียหาย คืนกลับแล้วทั้งหมด**:
  - BPA scope ใน section 2: ประโยคเดิมหลัง trim อ่านกำกวมจนสอนผิดว่า BPA ระดับ account ก็เปิด-ปิดได้
    แค่ยกชุด 4 setting เหมือน organization — **ความจริง: ระดับ account ยังเลือกทีละ setting ได้ปกติ
    มีแค่ระดับ organization เท่านั้นที่ยกชุดอย่างเดียว** เขียนประโยคใหม่แยกสองระดับให้ชัด
  - fintech worked example's punchline: ชี้ไปหัวข้อ 3 ที่ไม่เคยพูดถึง Object Lock/retention เลย —
    เขียน connective chain ใหม่ครบ (retention 7 ปี → object สะสมมาก → SSE-KMS ไม่เปิด Bucket Keys
    ต้นทุนพุ่ง → เปิด Bucket Keys ลดต้นทุนแต่ CloudTrail เห็น event หยาบขึ้น → ขัดกับข้อ 4 ที่ต้อง audit
    ละเอียด)
  - Versioning ใน section 4: "กู้คืนได้" ห้อยลอยไม่มีประธาน แก้เป็น "ของเดิมยังอยู่ครบและกู้คืนได้"
  - เพิ่มเติม: เจอ duplicate จริงอีกจุดที่ round 1 พลาด (`Versioning ทำให้บิลโตเงียบ ๆ` ซ้ำกับ section 4's
    "(ค) ทุกเวอร์ชันคิดเงินเต็มใบ") ตัดออก เหลือ 4 bullet ใน trade-off list (จาก 5)
  - opener ของ "แก่นของเรื่อง": คืน causal connector ที่ trim ตัดทิ้ง ("จำตารางนี้ให้ได้ เพราะ...")
- **R4 — SSE-C fact ถูกตัดให้แคบกว่าความจริงจนกลายเป็นข้อมูลผิด**: ข้อความเดิมพูดแค่ครึ่งเดียวของ fact
  ที่ยืนยันแล้ว (general purpose bucket **ใหม่** ปิด SSE-C) ตัดครึ่งที่สองทิ้ง (bucket **เดิม** ของ
  account ที่ไม่มี object เข้ารหัสด้วย SSE-C อยู่เลยก็ถูกปิดให้ด้วย — ขอบเขตเป็นระดับ **account** ไม่ใช่
  ราย bucket) ทำให้ผู้อ่านสรุปผิดว่า bucket เดิมยังใช้ SSE-C ได้เสมอ — คืน fact ทั้งสองครึ่งพร้อม scope
  ที่ถูกต้อง ทั้งใน `body_md` §3's table และใน `recall_checks[1].explanation`
- **R1 — ทั้งสอง pilot ยังไม่ทำตามเทมเพลตที่ ticket เดียวกันปักไว้**: `s3-security` ทั้ง 4 explanation
  ยังเป็นร้อยแก้ว inline (ก)(ข)(ค) ไม่มี bold header เลย, `vpc-fundamentals` ไม่มี "คำในโจทย์ →
  คำตอบที่ต้องมองก่อน" section เลย (มีแค่ "## มุมข้อสอบ" แบบ bullet list) — แปลง `s3-security`'s 4
  explanation เป็น bold header ครบ, เพิ่ม section ให้ `vpc-fundamentals` เป็นตาราง (ปักรูปแบบตารางไว้
  เป็นมาตรฐานเดียว ไม่ใช่ bullet — ดู Template ข้อ 3)
- **R5 — lesson-verifier ไม่มีทางจับข้อบกพร่องข้างบนได้เลยสักข้อ**: `s3-security` ผ่าน verifier ทั้งที่
  ไม่มี bold header เลยสักข้อ เป็นหลักฐานตรงว่า verifier เช็คได้แค่ "มี 3 ส่วน" ไม่เช็ครูปแบบ — เพิ่ม
  เกทใน `lesson-verifier.md` (bold header ครบ 4, cue section เป็นตารางและมีจริง, prose-rune ≤ 950,
  recall check count 4–6, correct-answer-is-longest ไม่เกิน 1 ข้อ, quote verbatim) และแก้
  `aws-cert.md`'s "## รูปแบบ explanation" ให้ระบุชื่อ header ทั้ง 4 ตรงกับ `lesson-writer.md` (เดิมสอง
  เอกสารนี้พูดไม่ตรงกัน — เอกสารที่ verifier อ้างอิงไม่เคยรู้จัก header เลย)
- **R6 — "3 ใน 8" นับผิด ที่ถูกคือ "2 ใน 8"** ขัดกับตัวเลข guessability ในเอกสารเดียวกันเอง (longest
  heuristic 25.0%/+0.0% ที่ n=8 — ถ้าคำตอบถูกเป็นตัวยาวสุด 3 ใน 8 ข้อ ตัวเลขจะไม่ใช่ baseline พอดี)
  แก้ทั้งใน `lesson-writer.md` และ Template ข้อ 5 ด้านบน พร้อมเขียนกฎใหม่ให้วัดได้จริงต่อ lesson
  ("อย่างมากแค่ 1 ใน 4–6")
- **R7 — เลข ~550/×11 เดิมยังหลงเหลือใน 2 เอกสารที่ใช้งานจริง**: `docs/HANDOFF.md:128` (ตารางเดิม
  165/143/132/110 แบบคำต่อคำ — เอกสารนี้คือที่ session ใหม่อ่านก่อนเริ่มงาน จึงเป็นจุดเสี่ยงที่สุดที่เลข
  ตายจะถูกหยิบกลับมาใช้) และ `docs/roadmap.md:16` (ใน priority block ที่ยังใช้งานอยู่) — แก้ทั้งคู่เป็น
  4–6/~250 คงหมายเหตุประวัติที่ `roadmap.md:76` (ใต้ AWS-0) และ `aws-cert.md`'s AWS-S1/AWS-S3 section
  ไว้ตามเดิมเพราะเป็นบันทึกประวัติของการตัดสินใจตอนนั้นจริง ไม่ใช่เป้าปัจจุบัน

Also fixed (suggested, cheap): quote ทั้งสองที่ (body_md section 1 และ recall_checks[3]) ตอนนี้ขึ้นต้น
ด้วย `...` บอกว่าเป็น fragment และคง **and** ตัวหนาของต้นฉบับไว้; Verification prose แก้ "track aws" เป็น
"track aws-saa-c03" (ชื่อ track จริงคือ slug เต็ม ไม่ใช่ชื่อไฟล์ curriculum).

**Pilot results (หลัง round 3)**: ทั้งสองบทเรียน (`vpc-fundamentals`, `s3-security`) **PASS
lesson-verifier** — ทุก reference URL fetch แล้วยืนยันจริง, ทุก AWS claim เช็คกับ live docs, ไม่มี
คำตอบที่ชื่อ service อยู่ใน ban list, explanation ครบ 3 ส่วนทุกข้อภายใต้ bold header ครบทั้ง 4 ทุกข้อ,
guessability อยู่ที่ baseline พอดี — **25.0% ทั้ง 7 heuristic ที่ n=8** (insufficient sample, gate
รายงาน NOT JUDGED ไม่ใช่ PASS ตามที่ `MinSampleSize=100` ออกแบบไว้ — ดู "Verification" ด้านล่างสำหรับ
output เต็มจาก `cmd/mcq-guessability`)

**สิ่งที่ตั้งใจไม่ทำในรอบนี้**:

- ไม่แก้ status ของ AWS-S1/AWS-S3 ในเอกสารนี้ (ทั้งคู่ merge แล้วจริง — PR #58, #59 — แต่ ticket นี้
  ระบุให้ tick แค่ AWS-S3 ใน `roadmap.md` เท่านั้น ไม่ได้ระบุให้แก้ status ใน `aws-cert.md`)
- ไม่เขียน AWS-C1..C4 เอง (ยัง gated เหมือนเดิม, ticket นี้แค่ปักเทมเพลตที่ AWS-C1..C4 ต้องตามให้ถูก)

### Round 3 (code review) — 4 findings, ทุกข้อแก้แล้ว

Round 2 ได้ NO-SHIP แคบ 4 ข้อ (root cause ของ round 2's rune-budget derivation ยังไม่ reproduce +
เอกสารอีก 2 จุดยังไม่ sync + retrofit ที่ round 2 เลื่อนออกไปกลับกลายเป็นปัญหาจริง):

- **แก้ arithmetic ของ derivation (ดูหัวข้อ "R2" ด้านบนที่แก้ในรอบนี้)**: round 2 เขียน mean 678 +
  3σ(78) ≈ 909 ผสม population 118/120 ผิดและบวกเลขผิดด้วย (678+3×78 = 912 ไม่ใช่ 909) — round 3 วัดซ้ำ
  เองด้วย convention "นับบรรทัดว่างรวม" ได้ mean **681**, σ **78**, mean+3σ = **915** ตรงกับที่ควรจะเป็น
  ระบุ population เป็น all 120 ชัดเจน และระบุ blank-line convention ที่ขาดไปเดิม (ไม่ระบุ = ใครมาวัดซ้ำ
  จะได้คนละตัวเลข)
- **`aws-currency-checklist.md` §1.11 ยังไม่ได้แก้**: บทเรียนแก้เป็น scope ระดับ account ถูกแล้ว
  (R4 ของ round 2) แต่ checklist ที่ `lesson-writer.md` ประกาศว่า authoritative ยังเขียน per-bucket
  เดิม — คนเขียนบทเรียนถัดไปที่เปิด checklist นี้จะกลับไปเขียนผิดแบบเดียวกับที่ R4 เพิ่งแก้ แก้แล้ว
  ให้ตรงกับบทเรียน
- **Retrofit จุด "." ท้ายประโยคจริงในรอบนี้** (ไม่เลื่อนต่อแบบ round 2): `s3-security` 34 บรรทัด,
  `vpc-fundamentals` 30 บรรทัด ได้รับ "." ท้ายบรรทัดตาม majority convention ของคลัง — ก่อนแก้ทั้งสอง
  pilot เป็นบทเรียนที่ conform กับ convention นี้น้อยที่สุดในคลังทั้งหมด ขัดกับที่ round 2 อ้างว่า
  "ทั้งสอง pilot เห็นต่างกันเอง" ราวกับเป็นเรื่องที่ยุติแล้ว ทั้งที่จริงทั้งคู่ต่ำกว่า norm พร้อมกัน —
  แก้ attribution ของกฎด้วย: ตัวอย่างที่แท้จริงมาจาก `designing-data-intensive-applications`/
  `ai-and-llm-systems` ไม่ใช่ `domain-driven-design` ซึ่งจริง ๆ โหวตสวนทาง — ปรับประโยค BPA scope ใน
  section 2 ให้มี "." คั่นหลัง "bucket-level" ด้วย เพื่อตัดการอ่านเป็นวลีเดียวกับ "ระดับ organization"
  ที่กำกวม (เนื้อหาถูกอยู่แล้ว เป็นแค่การอ่านสะดุด)
- **คืนเนื้อหาที่หายไปพร้อมกับ bullet ที่ตัดทิ้งจริง**: ตอนตัด `Versioning ทำให้บิลโตเงียบ ๆ` ทิ้ง
  (duplicate จริงตามที่ round 2 พบ) ประโยคย่อยที่บอกว่า **workload แบบไหน** เสี่ยงสะสมค่าใช้จ่าย
  (เขียนทับไฟล์เดิมบ่อย) หายไปด้วยทั้งที่ไม่ใช่ duplicate — เพิ่มกลับเข้าไปในย่อหน้า Versioning section 4
  จุด (ค) แทน
- **`docs/roadmap.md:76`**: จุดสุดท้ายที่เหลือของ `~550` ที่ไม่ได้ label ว่าเป็นตัวเลขประวัติ — เพิ่ม
  label แล้ว

### Verification (round 3)

- `go vet ./...`: clean. `go test -count=1 ./...`: all packages `ok` (ไม่มีการแก้ไฟล์
  `.go`/`migrations`/`web` เลยในรอบนี้ — `git diff --name-only develop... -- '*.go' 'web/*'
  'migrations/*'` ว่างเปล่า)
- Import evidence (stack แยก `docker compose -p awsc0fix2`, ไม่แตะ `self-learning_mysql_data`, ปิดท้าย
  ด้วย `docker compose down` เปล่า ๆ ไม่มี `-v`): ดูผลจริงในรายงานสรุป PR
- `go run ./cmd/mcq-guessability -dir content/lessons`: track `aws-saa-c03` (n=8) ยังรายงาน 25.0%/+0.0%
  ทุก heuristic เหมือนเดิม (round 3 ไม่แตะ `options`/`expected_answer` ของ recall_checks เลย มีแค่
  `body_md` กับ `explanation` ที่แก้) — ดูผลจริงในรายงานสรุป PR
