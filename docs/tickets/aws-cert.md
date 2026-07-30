# AWS Certified Solutions Architect – Associate (SAA-C03) — แผนสอบ

> **Priority อันดับ 1 ปัจจุบัน** — สอบใน 1–2 เดือน, อ่านวันละ 1–2 ชม., ยังไม่เคยจับ AWS จริง
> (รู้แค่ชื่อ service หลัก ๆ) · Q-2d / Q-3 / SIM-\* พักไว้ระหว่างนี้ ([roadmap.md](../roadmap.md))

## Ticket ที่เกี่ยวข้อง

- **AWS-0** (เอกสารนี้) — รีบาลานซ์ `content/curriculum/aws.json` ให้ตรง exam guide จริง + pin แผนนี้
- **AWS-1..N** (ยังไม่ตัดชื่อ ticket ละ 1 domain) — เขียนคลังข้อสอบตามแผนด้านล่าง

## Blueprint ข้อสอบ (verified จาก official source — ไม่ต้อง re-verify)

**SAA-C03 คือเวอร์ชันปัจจุบัน ไม่มี SAA-C04** — หน้าเว็บบุคคลที่สามที่อ้างว่ามี SAA-C04 เป็นข้อมูลผิด

- 65 ข้อ (50 ข้อคิดคะแนน + 15 ข้อ unscored ปนอยู่โดยไม่บอกว่าข้อไหน)
- 130 นาที · ผ่านที่ 720/1000 (scaled score)
- **Compensatory scoring** — คะแนนรวมต้องผ่าน ไม่ได้ตัดเป็นรายหมวด
- **ไม่มีการหักคะแนนจากการเดา**
- รูปแบบคำถาม: multiple choice (เลือก 1 จาก 4) และ multiple response (เลือก 2+ จาก 5+ ตัวเลือก
  โดย stem จะบอกจำนวนที่ต้องเลือกเสมอ เช่น "(Select TWO.)")

แหล่งอ้างอิงชื่อ service ที่ถูกต้องคือ **docs.aws.amazon.com** ไม่ใช่ PDF จาก `d1.awsstatic.com`
(PDF v1.1 ล้าสมัยอย่างน้อย 2 จุด: ยังเขียน "AWS IAM Identity Center [AWS Single Sign-On]" และ
เขียน QuickSight ในจุดที่ docs ปัจจุบันใช้ "Amazon Quick"/"Amazon QuickSuite") — ถ้าต้องเช็คชื่อ
service ให้ fetch docs HTML ไม่ใช่เชื่อ PDF

### Domain และน้ำหนัก (ใช้ตัวเลขนี้เป๊ะ ๆ)

| Domain | น้ำหนัก | Task statements |
|---|---|---|
| 1. Design Secure Architectures | **30%** | 1.1 Secure access to AWS resources · 1.2 Secure workloads and applications · 1.3 Appropriate data security controls |
| 2. Design Resilient Architectures | **26%** | 2.1 Scalable and loosely coupled architectures · 2.2 Highly available and/or fault-tolerant architectures |
| 3. Design High-Performing Architectures | **24%** | 3.1 Storage · 3.2 Compute · 3.3 Database · 3.4 Network architectures · 3.5 Data ingestion and transformation |
| 4. Design Cost-Optimized Architectures | **20%** | 4.1 Storage · 4.2 Compute · 4.3 Database · 4.4 Network architectures |

### สองข้อสรุปเชิงกลยุทธ์ที่ต้องจำ

1. **ไม่มีการหักคะแนนจากการเดา → ห้ามเว้นข้อว่างเด็ดขาด** ถ้าไม่แน่ใจให้เดาตัวที่ตัดตัวเลือกผิด
   ออกได้มากที่สุด ดีกว่าเว้นว่างเสมอ (เว้นว่าง = ผิดชัวร์, เดา = มีโอกาสถูก)
2. **Compensatory scoring → มองภาพรวม ไม่ต้องกลัวหมวดใดหมวดหนึ่งอ่อน** สิ่งที่ต้องบริหารคือ
   คะแนนรวม 720/1000 ไม่ใช่การผ่านแยกรายหมวด — เวลาอ่านจึงควรกระจายตามน้ำหนัก (30/26/24/20)
   ไม่ใช่ทุ่มเวลาให้หมวดที่ถนัดจนหมวดอื่นได้เวลาน้อยเกินไป

## Concept → Task statement mapping (tree ที่รีบาลานซ์แล้ว, 50 concept)

Schema ปัจจุบันของ `content/curriculum/aws.json` (`internal/curriculum/infra/contentfile.go`'s
`conceptFileDTO`) มีแค่ `slug`/`title`/`outline`/`position` — ไม่มีช่องสำหรับ task statement และ
ตามกติกาของ ticket นี้ **ห้ามเพิ่ม field ใหม่ในไฟล์นั้น** เพราะจะกระทบ parser และทุก content
file อื่น ๆ ที่ใช้ schema เดียวกัน (ddd/ddia/distsys/dsa/go/ai-systems) mapping จึงอยู่ในตารางนี้
แทน — ใช้ตอนเขียนคลังข้อสอบเพื่อ tag ว่าแต่ละข้อ cover task ไหน

### Domain 1 — Secure Architectures (15 concepts)

| Concept slug | Task |
|---|---|
| iam-users-roles-policies | 1.1 |
| iam-policy-evaluation | 1.1 |
| organizations-scp | 1.1 |
| cognito | 1.1 |
| kms-encryption | 1.3 |
| secrets-vs-parameter-store | 1.2 |
| sg-vs-nacl | 1.2 |
| vpc-endpoints-privatelink | 1.2 |
| waf-shield | 1.2 |
| s3-security | 1.3 |
| vpc-fundamentals | 1.2 (prerequisite ให้ 3.4, 4.4 ด้วย) |
| hybrid-cross-vpc-connectivity | 1.2, 3.4 |
| security-detection-governance | 1.1 |
| multi-account-access-governance | 1.1 |
| acm-tls-in-transit-encryption | 1.3 |

### Domain 2 — Resilient Architectures (13 concepts)

| Concept slug | Task |
|---|---|
| regions-az-edge | 2.2 |
| elb-types | 2.1 |
| auto-scaling-groups | 2.1 |
| rds-multi-az-read-replicas | 2.2 |
| aurora-ha | 2.2 |
| sqs-sns-decoupling | 2.1 |
| eventbridge | 2.1 |
| route53-routing-policies | 2.2 |
| dr-strategies | 2.2 |
| backup-strategies | 2.2 |
| observability-cloudwatch-cloudtrail | 2.2 |
| api-gateway-step-functions | 2.1 |
| managed-ai-services-overview | 2.2 |

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
| rds-performance | 3.3 |
| kinesis | 3.5 |
| athena-glue | 3.5 |
| lambda-performance | 3.2 |
| ecs-eks-fargate | 3.2 |

Domain 3 เดิมมี 12 concept กระจุกที่ 3.1/3.2/3.3 อยู่แล้ว (9 จาก 12) — รอบนี้**ไม่เพิ่มจำนวน**
แต่ enrich outline ของ concept ที่มีอยู่แทน (ดูหัวข้อถัดไป) เพื่อดันน้ำหนักไปทาง 3.4/3.5 โดยไม่ทำให้
Domain 3 เกิน 24% ของ tree

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
| migration-and-transfer-services | 4.1, 3.5 |
| storage-cost-optimization-beyond-s3 | 4.1 |

### Small named items folded into existing concepts' outline (ไม่ได้แยก concept ใหม่)

รายการ AWS-0 ไล่จาก exam guide ที่ปริมาณคำถามน้อยแต่ต้องมี ไม่คุ้มแยก concept ใหม่ — เพิ่มเป็น
บรรทัดในบทเรียน (outline) ของ concept ที่ใกล้เคียงที่สุดแทน:

| Item | Folded into |
|---|---|
| AWS Resource Access Manager (RAM) | `organizations-scp` |
| CloudFormation / immutable infrastructure, Elastic Beanstalk | `auto-scaling-groups` |
| AWS Outposts | `regions-az-edge` |
| Amazon MQ | `sqs-sns-decoupling` |
| EC2 hibernation | `ec2-families-purchasing` |
| Amazon FSx (all types) | `ebs-efs-instance-store` |
| RDS Proxy | `rds-performance` |
| AWS Batch, Amazon ECR | `ecs-eks-fargate` |
| Amazon EMR | `athena-glue` |
| DocumentDB, Neptune, Keyspaces (purpose-built DB) | `dynamodb-advanced` |
| AWS Trusted Advisor | `cost-tools-budgets` |

## แผนคลังข้อสอบ (~550 ข้อ)

จัดสรรตามน้ำหนัก domain (30/26/24/20) ของ 550 ข้อ:

| Domain | น้ำหนัก | จำนวนข้อ |
|---|---|---|
| 1. Secure Architectures | 30% | **165** |
| 2. Resilient Architectures | 26% | **143** |
| 3. High-Performing Architectures | 24% | **132** |
| 4. Cost-Optimized Architectures | 20% | **110** |
| **รวม** | 100% | **550** |

เขียนทีละ domain (batch = 1 domain, ไม่ใช่ 1 concept) โดย `lesson-writer` เขียนแล้วส่งต่อ
`lesson-verifier` fact-check ทุก batch เสมอ ตามกติกา orchestration ปกติของโปรเจกต์ — FAIL ต้อง
regenerate ก่อนเก็บเป็น final

## กติกาการเขียนคำถาม (ticket เขียนคลังข้อสอบทุกตัวต้องทำตาม)

- **Stem แบบ scenario** สไตล์ข้อสอบจริง ไม่ใช่ถามนิยามตรง ๆ — ตั้งบริบทธุรกิจ/ข้อจำกัดก่อนถาม
- ใช้ **qualifier** ที่ทำให้ตัวเลือกที่ทำงานได้จริงหลายตัวกลายเป็นตัวเลือกที่ผิด เช่น **MOST
  cost-effective / LEAST operational overhead / MOST secure** — โจทย์ SAA มักมีคำตอบที่ "ใช้ได้"
  มากกว่า 1 ตัว แต่มีตัวเดียวที่ตรง qualifier
- **Distractor ต้องเป็น AWS service จริงที่ใช้ผิดบริบท** ไม่ใช่ตัวเลือกที่มองออกชัด ๆ ว่าไม่เกี่ยว
- **Multiple-response ต้องระบุจำนวนใน stem เสมอ** เช่น "(Select TWO.)"
- กติกา guessability เดิมใน [`docs/tickets/mcq-quality.md`](mcq-quality.md) ยังใช้ต่อ (ห้าม
  distractor ยาว/สั้นผิดปกติจนเดาได้จาก pattern ความยาวหรือตำแหน่ง)

## รูปแบบ explanation (ข้อกำหนดสำคัญที่สุดจากผู้ใช้)

ทุกข้อต้องมี explanation ที่มีครบ 3 ส่วนนี้เสมอ (Thai prose, technical term เป็นอังกฤษ):

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

## Status

**In review** — AWS-0 (curriculum rebalance + เอกสารนี้) รอ code review

## Review focus

<details>
<summary>คำถามสำหรับรีวิว diff นี้ (เฉลยพับไว้ด้านล่าง)</summary>

1. ทำไม concept ใหม่ 13 ตัวถึงถูกเติมท้ายแต่ละ chapter (position ต่อจากตัวสุดท้ายเดิม) แทนที่จะ
   แทรกกลางลำดับตามความสัมพันธ์เชิงเนื้อหา (เช่น `vpc-fundamentals` ควรมาก่อน `sg-vs-nacl`)?
2. ทำไม `internal/curriculum/infra/contentfile_aws_test.go` ต้องแก้ทั้งที่ ticket บอกว่า scope
   คือ `content/curriculum/aws.json` + docs เท่านั้น ไม่แตะ `internal/`?
3. ทำไม Domain 3 (high-performing) ยังคงมี 12 concept เท่าเดิมทั้งที่ ranked gap list มีทั้ง
   "hybrid connectivity" (เกี่ยวกับ Task 3.4) และ "purpose-built databases" (Task 3.3) เป็นช่อง
   ว่างของ Domain 3 โดยตรง?

<details>
<summary>เฉลย</summary>

1. เพื่อคุม diff ให้เล็กที่สุด — ถ้าแทรกกลางลำดับ ต้อง renumber `position` ของ concept เดิมที่อยู่
   ถัดจากจุดแทรกทั้งหมด ทำให้ diff บวมโดยไม่จำเป็น (การเรียงอ่านตามลำดับที่ "ดีที่สุด" ไม่ใช่
   ข้อกำหนดของ ticket นี้ แค่ position ไม่ชนกันและ unique พอ) ส่วนลำดับการอ่านจริงเป็นเรื่องของ
   ticket เขียนบทเรียน/คลังข้อสอบทีหลัง ไม่ใช่ของ tree นี้
2. `contentfile_aws_test.go` เป็น golden test ที่ hardcode slug list ของ `aws.json` เป๊ะ ๆ
   (comment ในไฟล์บอกตรง ๆ ว่า "guards content/curriculum/aws.json itself") มันไม่ใช่ business
   logic แต่เป็น snapshot ของไฟล์ data — เมื่อ data เปลี่ยนโดยตั้งใจ (ตามที่ ticket สั่ง) ต้อง sync
   golden data นี้ด้วยไม่งั้น `go test ./...` จะแดงถาวร ซึ่งขัดกับ verification requirement ของ
   ticket เอง (`go vet` + `go test ./...` ต้องผ่าน) — เป็นข้อยกเว้นแคบ ๆ ที่จำเป็นทางกลไก ไม่ใช่การ
   แตะ business logic ของ `internal/`
3. เพราะ Domain 3 (12/50 = 24%) **ตรงสัดส่วนเป้าหมายพอดีอยู่แล้ว** เทียบกับ total tree ใหม่ 50
   concept — ต่างจากตอนก่อนรีบาลานซ์ที่ 12/37 = 32% (เกินสัดส่วนจริง) การเพิ่ม concept ใหม่เข้า
   Domain 3 อีกจะทำให้เกิน 24% จึงต้อง"redistribute"แทน"grow": เนื้อหาของ hybrid connectivity
   (Task 3.4, แต่ mechanic เดียวกับที่ต้องใช้ตอบ Task 1.2 ด้วย) ถูกวางไว้ที่ Domain 1 แทน (ซึ่ง
   ต้องโต 10→15 อยู่แล้ว), ส่วน purpose-built databases (DocumentDB/Neptune/Keyspaces) ถูก fold
   เป็นบรรทัดเพิ่มใน outline ของ `dynamodb-advanced` ที่มีอยู่แล้วแทนการแยก concept ใหม่ — mapping
   table ด้านบนยังบันทึกไว้ครบว่า concept ไหน cover task ไหนบ้าง แม้จะข้ามไปอยู่คนละ chapter

</details>
</details>
