# AWS Certified Solutions Architect – Associate (SAA-C03) — แผนสอบ

> **Priority อันดับ 1 ปัจจุบัน** — สอบใน 1–2 เดือน, อ่านวันละ 1–2 ชม., ยังไม่เคยจับ AWS จริง
> (รู้แค่ชื่อ service หลัก ๆ) · Q-2d / Q-3 พักไว้บางส่วนระหว่างนี้ (ดูหัวข้อ blocker ด้านล่าง — Q-3
> ถูก unpause บางส่วนเพราะ schema change (2) คือ schema half ของ Q-3 เอง) · SIM-\* พักเต็ม
> ([roadmap.md](../roadmap.md))

## Ticket ที่เกี่ยวข้อง

- **AWS-0** (เอกสารนี้) — รีบาลานซ์ `content/curriculum/aws.json` ให้ตรง exam guide จริง + pin แผนนี้
  (round 2: แก้ currency issues + content-model blocker หลัง code review)
- **AWS-1** — schema prerequisites (ดูหัวข้อ blocker) + เขียนคลังข้อสอบตามแผนด้านล่าง คนละ ticket
  จาก AWS-0 เพราะแตะ `internal/` ซึ่ง AWS-0 ไม่แตะ

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
**PR #54 ต้อง merge เข้า develop ก่อน AWS-1 เริ่มงานเสมอ** ไม่ใช่แค่ก่อนเขียนคำถามจริง

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

### 3 schema change ที่ต้องทำก่อน (นี่คือ AWS-1's งาน ไม่ใช่ AWS-0 — AWS-0 ไม่แตะ `internal/`)

1. **เพิ่ม `maxRecallChecks`** จาก 5 เป็น ~15 (เผื่อ headroom เหนือเป้า 11/concept)
2. **เพิ่มฟิลด์ `explanation`**: domain (`RecallCheck`) → `recallCheckFileDTO` → importer → API
   response → หน้า quiz ต้อง render — **นี่คือ schema half ของ Q-3** (Q-3 เดิมคือ "เพิ่ม explanation
   ให้ MCQ 356 ข้อที่มีอยู่") ดังนั้น Q-3 **ถูก unpause บางส่วน**: schema change ต้องทำ (ทำที่ AWS-1)
   แต่การเขียน explanation ให้ 356 ข้อเดิมยังพักไว้ได้จนกว่าจะว่าง
3. **เพิ่ม multiple-response support**: ต้องมี second correct-answer representation (เช่น
   `expectedAnswers []string` หรือ MCQ กับ multiple-response แยก type) + UI ฝั่ง quiz ให้เลือกได้
   มากกว่า 1 ตัวเลือก — multi-response เป็น ~15–20% ของข้อสอบจริงและเป็นประเภทที่คนเสียคะแนนมากที่สุด
   ทิ้งไปไม่ได้

**ข้อจำกัดที่ตามมาจนกว่า (2) จะเสร็จ: ไฟล์คำถามชุดไหนก็ตามที่มี `explanation` จะถูก
`DisallowUnknownFields()` reject ทั้งไฟล์ตอน import — ห้ามเขียนคำถามที่มี explanation ก่อน schema
change (2) landed จริงใน MySQL + importer**

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
- จัด ticket AWS-1..N เป็น **1 ticket ต่อ domain** เพื่อความสะดวกเวลารีวิว (ticket ละ 10–15
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
  แทน 33%

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
ว่า PR #54 ต้อง merge ก่อน AWS-1 เริ่ม — **`aws.json` ยังไม่มีการเพิ่ม/ลบ/ย้าย concept ในรอบนี้เช่นกัน
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
   test ทันที เป็นการเปลี่ยนแปลงเชิงโครงสร้างที่ควรอยู่ใน ticket ที่แตะ `internal/` ได้ (เช่น AWS-1)
   ไม่ใช่ ticket นี้ที่ scope จำกัดไว้แค่ content + docs
3. เพราะ dual-mapping ที่บันทึกในตารางเป็นแค่ **ข้อมูลสำหรับ tag คำถามตอนเขียนคลังข้อสอบ** (ข้อไหน
   cover task ไหนบ้าง) ไม่ใช่ตัวกำหนดว่า concept ต้องอยู่ chapter ไหน — `rds-multi-az-read-replicas`
   ยังคงเป็นเนื้อหาแกนของ Domain 2 (high availability ผ่าน Multi-AZ) การที่มันมีบรรทัดหนึ่งที่ตอบ
   Task 2.1 ได้ด้วย (read replica ช่วย scale) ไม่ได้แปลว่ามันควรย้ายไป Domain นั้น เหมือนกับ
   `vpc-fundamentals`/`hybrid-cross-vpc-connectivity` ที่ยอมรับไว้แล้วตั้งแต่ round 1 ว่า
   chapter-count ตรงสัดส่วนน้ำหนัก แต่ task-serving ข้ามโดเมนได้

</details>
</details>
