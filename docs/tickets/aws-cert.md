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
- **AWS-S3** (ยังไม่เริ่ม, spec ปักไว้แล้ว) — guessability measurement tool ที่ parameterize baseline
  ได้ต่อ corpus — ดูหัวข้อ "AWS-S3" ด้านล่าง gate คลังข้อสอบจริง
- **AWS-C1..C4** (ยังไม่เริ่ม, gated บน S1/S2/S3) — คลังข้อสอบจริง ~550 ข้อ ทีละ domain

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
"บทเรียน" ที่ผู้ใช้อยากได้ด้วย) + ~11 recall_checks** (คำถามของ concept นั้น) — 50 × 11 = 550 พอดี
และยอดต่อ domain ก็ตกลงมาเองจากจำนวน concept ต่อ domain × 11:

| Domain | concepts | × 11 | จำนวนข้อ |
|---|---|---|---|
| 1. Secure Architectures | 15 | ×11 | **165** |
| 2. Resilient Architectures | 13 | ×11 | **143** |
| 3. High-Performing Architectures | 12 | ×11 | **132** |
| 4. Cost-Optimized Architectures | 10 | ×11 | **110** |
| **รวม** | 50 | ×11 | **550** |

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

## แผนคลังข้อสอบ (~550 ข้อ + 50 lesson)

จำนวนข้อต่อ domain มาจากตาราง "รูปทรงใหม่" ด้านบน (11 ข้อ/concept × concept ต่อ domain) ไม่ใช่การ
หาร 550 ตามเปอร์เซ็นต์น้ำหนักตรง ๆ อีกต่อไป — ตัวเลขบังเอิญออกมาเท่ากันเพราะ concept count ต่อ
domain ถูกออกแบบให้ตรงสัดส่วนน้ำหนักอยู่แล้ว (15/13/12/10 = 30/26/24/20%)

**หน่วยของการทำงานคือ concept ไม่ใช่ domain** ตาม token-efficiency skill
(`.claude/skills/token-efficiency/SKILL.md`: "lesson-writer/verifier: one concept per invocation")
— แผนเดิมที่ให้ `lesson-writer` เขียนทีเดียว 132–165 ข้อ/domain ละเมิดกฎนี้ตรง ๆ และเป็นสาเหตุคลาสสิก
ของ explanation ที่ตื้น/restate นิยามซ้ำ (โมเดลล้าเมื่อ context ยาว):

- **1 invocation ต่อ 1 concept**: `lesson-writer` เขียน lesson (concept note) + recall_checks
  ~11 ข้อของ concept นั้น (เพดานหลัง schema change (1) คือ 15) พร้อมกัน → `lesson-verifier`
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
  แทน 33% — **แต่เกทนี้ unenforceable จนกว่า AWS-S3 จะสร้างเครื่องมือวัดจริง** (ดู "AWS-S3 —
  guessability measurement tool" ด้านล่างสำหรับ spec) เพราะไม่มี script วัดใด ๆ commit ไว้ใน
  โปรเจกต์นี้เลย (verified: ไม่มีทั้งใน working tree และ `git log --all --diff-filter=A`) — **ห้าม
  เริ่มเขียนคำถามชุดจริง (AWS-C1..C4) ก่อน AWS-S3 เสร็จ** เพราะไม่มีทางวัดว่าชุดที่เขียนผ่านเกท 25%
  จริงหรือไม่จนกว่าเครื่องมือจะมี

## AWS-S3 — guessability measurement tool (spec ที่นี่เท่านั้น — ticket แยกต่างหาก ไม่ทำใน AWS-S1)

**สถานะ**: ยังไม่เริ่ม — เอกสารนี้กำหนด spec ขั้นต่ำไว้ล่วงหน้าเพื่อให้ ticket ถัดไปสร้างได้เลย
ไม่ต้องออกแบบใหม่ gate ทั้งเก่า (33% บน corpus 3 ตัวเลือก) และใหม่ (25% บน corpus AWS 4 ตัวเลือก)
**unenforceable จนกว่า tool นี้จะมี** เพราะไม่มี measurement script ใด commit ไว้ในโปรเจกต์

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

- **`maxRecallChecks` 5 → 15**: แผน AWS ต้องการ ~11 ข้อ/concept (docs/tickets/aws-cert.md) เพดาน 15
  เผื่อ headroom เหนือเป้านั้นโดยไม่เปิดให้ quiz ต่อ lesson ยาวไม่จำกัด
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
- **Static guard**: `TestAllMigrationsAlterTableAddColumnIsGuarded` (ใหม่ ใน
  `internal/platform/mysql/migrate_embedded_test.go`) เช็คว่าไฟล์ migration ใดก็ตามที่มี
  `ADD COLUMN` ต้องมี marker ของ guard pattern ครบ (`information_schema.COLUMNS`, `PREPARE `,
  `EXECUTE `, `DEALLOCATE PREPARE`) — **ไม่ใช่การ extend `TestAllMigrationsCreateTableIsIdempotent`
  ตรง ๆ** เพราะ `ADD COLUMN` ไม่มี token เดียวแบบ `IF NOT EXISTS` ให้เช็คต่อท้ายเหมือน `CREATE TABLE`
  จึงเขียนเป็นเทสต์คู่ขนานแยกต่างหากที่เช็ค marker set ของ pattern นี้โดยเฉพาะ เทสต์นี้พิสูจน์แค่
  **syntax** ของ guard ว่ามีอยู่ (ข้อจำกัดแบบเดียวกับเทสต์ CREATE TABLE เดิม) ส่วน**พฤติกรรม**จริงของ
  MySQL พิสูจน์ด้วย live re-run ข้างบนแทน ไม่ใช่ `go test`
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
  (round 1) → 301 (round 2)** (นับจาก `go test -v` ผ่าน `grep -c "^--- PASS"`; round 2 เพิ่ม
  `TestAllMigrationsAlterTableAddColumnIsGuarded` + `TestHandlerGetLesson_OmitsExplanationKeyWhenAbsent`)
- Frontend: `npx tsc --noEmit` clean, `npx eslint .` clean, `npm run build` clean,
  `npx vitest run` **169 (develop) → 175 (round 1) → 176 (round 2)** (เพิ่ม explanation=""
  test case)

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
3. เพราะ**ข้อสรุป**ของ decision (คงพื้นไว้ที่ 3) ยังถูกต้องอยู่ — สิ่งที่ผิดคือ**หลักฐาน**ที่ยกมาอ้าง
   (อ้างว่ามีใบที่มีพอดี 3 ข้อ ซึ่งนับจริงแล้วไม่มี ต่ำสุดคือ 4) ไม่ใช่ตัวการตัดสินใจเอง — ไม่มีอะไรใน
   scope ของ ticket นี้เรียกร้องให้เปลี่ยน floor จาก 3 เป็นค่าอื่น (blocker เดิมคือเพดานบนเท่านั้น)
   ดังนั้นแก้แค่ comment ให้ตรงข้อเท็จจริง ไม่ต้องเปลี่ยนค่า constant

</details>
</details>
