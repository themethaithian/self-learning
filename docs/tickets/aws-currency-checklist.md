# SAA-C03 Currency Checklist

> **ทำไมไฟล์นี้ถึงมีอยู่**: โมเดลที่เขียนข้อสอบมีวันหมดอายุของความรู้ และ AWS เปลี่ยนเร็วมาก
> ไฟล์นี้คือรายการ "สิ่งที่โมเดลน่าจะจำมาผิด" ที่ตรวจกับเอกสาร AWS ทางการแล้วทุกบรรทัด
> **ทุก item มี URL** — ถ้าไม่มี URL แปลว่ายังไม่ได้ยืนยัน และอยู่ใน §10
>
> ใช้เป็น **context บังคับ** ของทั้ง `lesson-writer` และ Fable adjudicator ทุก batch
> ค้นเมื่อ 2026-07-30 · แหล่งที่ยอมรับ: `docs.aws.amazon.com` และ `aws.amazon.com` เท่านั้น

---

## §0 สิ่งที่ต้องเปิดซ้ำทุก batch ก่อนเริ่มเขียน

สองหน้านี้มีค่ามากกว่าข้อเท็จจริงรายตัว เพราะ **นี่คือที่ที่คำตอบคลาสสิกของข้อสอบไปตาย**

| หน้า | คืออะไร | URL |
|---|---|---|
| **Services in Maintenance** | รายการ service/feature ที่ **ปิดรับลูกค้าใหม่** | https://docs.aws.amazon.com/general/latest/gr/maintenance_services.html |
| **Services in Sunset** | รายการที่มี **วันสิ้นสุดการ support** ชัดเจน | https://docs.aws.amazon.com/general/latest/gr/sunset_services.html |
| นิยาม lifecycle stage | Maintenance / Sunset / Full Shutdown | https://docs.aws.amazon.com/general/latest/gr/service-lifecycle.html |

ระหว่าง ต.ค. 2025 – มี.ค. 2026 ปีเดียว AWS ย้าย S3 Object Lambda, Snowball Edge, App Runner,
Audit Manager, CloudTrail Lake, Kendra และ Simple AD เข้า maintenance — **รายการนี้จะขยับอีกแน่นอน**

**ยืนยันแล้ว: ข้อสอบยังเป็น SAA-C03 ไม่มี SAA-C04**
65 ข้อ (50 คิดคะแนน + 15 ไม่คิด) · 130 นาที · 720/1000 · โดเมน 30/26/24/20
https://docs.aws.amazon.com/aws-certification/latest/solutions-architect-associate-03/solutions-architect-associate-03.html

**สำคัญ: AWS Certification อัปเดตรายชื่อ service ใน scope จริง** — รายการปัจจุบันใช้ชื่อหลัง rename แล้ว
(`Amazon Data Firehose`, `Amazon Quick`, `Amazon SageMaker AI`, `AWS IAM Identity Center`, `Amazon OpenSearch Service`)
→ **อย่าคิดว่าข้อสอบยังใช้ชื่อเก่า ให้ใช้ชื่อปัจจุบัน**
- in-scope: https://docs.aws.amazon.com/aws-certification/latest/solutions-architect-associate-03/saa-03-in-scope-services.html
- out-of-scope: https://docs.aws.amazon.com/aws-certification/latest/solutions-architect-associate-03/saa-03-out-of-scope-services.html
- ใช้ short name: https://docs.aws.amazon.com/aws-certification/latest/solutions-architect-associate-03/saa-service-mentions.html

แต่คู่มือ **ตามความจริงไม่ทัน** — จุดที่ตามไม่ทันคือ §9

---

## §1 Tier 1 — โผล่บ่อยและโมเดลเก่าตอบผิดแน่

| # | เปลี่ยนอะไร | เมื่อไหร่ | คำตอบเก่า → คำตอบที่ถูก | URL |
|---|---|---|---|---|
| 1.1 | **public IPv4 คิดเงินแล้ว** $0.005/IP/ชม. ทุกตัว ต่อหรือไม่ต่อก็คิด รวม EC2/RDS/EKS/ELB | 2024-02-01 | "ฟรีถ้า EIP ต่ออยู่" → **ทุก public IPv4 มีค่าใช้จ่าย**; เลี่ยงด้วย IPv6 / private subnet + NAT / BYOIP | https://aws.amazon.com/blogs/aws/new-aws-public-ipv4-address-charge-public-ip-insights/ |
| 1.2 | **S3 เข้ารหัส object ใหม่ทุกชิ้นด้วย SSE-S3 โดยดีฟอลต์** | 2023-01-05 | "ต้องไปเปิด default encryption เอง" → **เข้ารหัสอยู่แล้ว** คำถามต้องถามว่า *key แบบไหน* (SSE-S3/SSE-KMS/DSSE-KMS/SSE-C) ไม่ใช่ว่าเข้ารหัสหรือยัง | https://aws.amazon.com/about-aws/whats-new/2023/01/amazon-s3-automatically-encrypts-new-objects/ |
| 1.3 | **bucket ใหม่: Block Public Access เปิด + ACL ปิด โดยดีฟอลต์** ทุกช่องทางสร้าง | 2023-04-27 | "bucket เปิด public จนกว่าจะบล็อก" / "ใช้ ACL ให้สิทธิ์" → **ใช้ bucket policy / IAM เท่านั้น** คำตอบที่อิง ACL ผิดโดยดีฟอลต์ | https://aws.amazon.com/about-aws/whats-new/2023/04/amazon-s3-security-best-practices-buckets-default/ |
| 1.4 | **S3 Select ปิดรับลูกค้าใหม่** (รวม Glacier Select) | 2024-07-25 | "ใช้ S3 Select ดึงบางส่วนของไฟล์" → **Athena** · ⚠️ ดู §9 C3 | https://docs.aws.amazon.com/AmazonS3/latest/userguide/selecting-content-from-objects.html |
| 1.5 | **S3 Object Lambda ปิดรับลูกค้าใหม่** | 2025-11-07 | "ใช้ Object Lambda แปลง/ปกปิดข้อมูลตอน GET" → CloudFront+Lambda · ⚠️ §9 C4 | https://docs.aws.amazon.com/AmazonS3/latest/userguide/amazons3-ol-change.html |
| 1.6 | **Snowball Edge ปิดรับลูกค้าใหม่** — ไม่มีอุปกรณ์ Snow ตัวไหนสั่งได้อีกแล้ว | ประกาศ 2025-10-07 มีผล 2025-11-07 | "ส่งข้อมูล 100 TB ผ่านลิงก์ช้า → Snowball" → **DataSync**, AWS Data Transfer Terminal, Outposts · ⚠️ §9 C1 | https://docs.aws.amazon.com/snowball/latest/developer-guide/snowball-edge-availability-change.html |
| 1.7 | **IMDSv2 เป็นดีฟอลต์** — instance type ใหม่ ๆ เป็น IMDSv2-only | 2022→2024 | "IMDSv1 ใช้ได้ ยิง 169.254.169.254 ได้เลย" → **ต้องมี token** | https://aws.amazon.com/about-aws/whats-new/2024/03/set-imdsv2-default-new-instance-launches/ |
| 1.8 | **root MFA บังคับทุกประเภทบัญชี** (ทยอยตั้งแต่ 2024 ครบ 2025-06) + centralized root access management ลบ root credential ของ member account ได้ | 2024–2025 | "เปิด root MFA เป็น best practice" → **บังคับแล้ว** | https://aws.amazon.com/about-aws/whats-new/2025/06/aws-iam-mfa-root-users-across-all-account-types/ |
| 1.9 | **S3 strongly consistent ทุก operation ทุก Region ฟรี** | 2020-12-01 | "S3 eventually consistent ต้องใช้ DynamoDB ทำ index" → **ห้ามออกข้อสอบแนวนี้เด็ดขาด มันไม่มีอยู่แล้ว** | https://aws.amazon.com/about-aws/whats-new/2020/12/amazon-s3-now-delivers-strong-read-after-write-consistency-automatically-for-all-applications/ |
| 1.10 | **Aurora Serverless v1 ตายแล้ว** (EOL 2025-03-31) และ **v2 scale ลง 0 ACU ได้แล้ว** | 2024-11 / 2025 | ตัวแยกแยะเดิม "v1 pause ได้ v2 ไม่ได้" **ตายสนิท ห้ามใช้** | https://aws.amazon.com/about-aws/whats-new/2024/11/amazon-aurora-serverless-v2-scaling-zero-capacity/ |

---

## §2 การเปลี่ยนชื่อ service ใน scope

ชื่อเก่าคือ error ที่มองเห็นได้ทันที — ไล่เช็กทุก draft

| ชื่อเก่า | ชื่อปัจจุบัน | เมื่อไหร่ | URL |
|---|---|---|---|
| AWS Single Sign-On (AWS SSO) | **AWS IAM Identity Center** | 2022-07-26 | https://aws.amazon.com/about-aws/whats-new/2022/07/aws-single-sign-on-aws-sso-now-aws-iam-identity-center/ |
| Amazon Elasticsearch Service | **Amazon OpenSearch Service** | 2021-09-08 | https://docs.aws.amazon.com/opensearch-service/latest/developerguide/rename.html |
| Kinesis Data Firehose | **Amazon Data Firehose** | 2024-02 | https://aws.amazon.com/about-aws/whats-new/2024/02/amazon-data-firehose-formerly-kinesis-data-firehose/ |
| Kinesis Data Analytics | **Amazon Managed Service for Apache Flink** | 2023-08 | https://aws.amazon.com/about-aws/whats-new/2023/08/amazon-managed-service-apache-flink/ |
| Amazon QuickSight | **Amazon Quick Suite** (คู่มือสอบเขียน "Amazon Quick") | 2025-10-09 | https://aws.amazon.com/blogs/business-intelligence/reimagine-business-intelligence-amazon-quicksight-evolves-to-amazon-quick-suite/ |
| Route 53 Application Recovery Controller | **Amazon Application Recovery Controller (ARC)** | (วันที่ไม่ระบุในหน้า) | https://docs.aws.amazon.com/r53recovery/latest/dg/what-is-route53-recovery.html |
| AWS Security Hub (ตัวเดิม) | **AWS Security Hub CSPM** — และชื่อ "AWS Security Hub" ตอนนี้หมายถึง service **คนละตัว** | 2025 | https://docs.aws.amazon.com/securityhub/latest/userguide/what-are-securityhub-services.html |
| Amazon S3 Glacier (storage class) | **S3 Glacier Flexible Retrieval** | 2021-11-30 | https://aws.amazon.com/about-aws/whats-new/2021/11/amazon-s3-glacier-instant-retrieval-storage-class/ |
| Amazon SageMaker | **Amazon SageMaker AI** | — | in-scope list §0 |
| ElastiCache for Redis | **ElastiCache for Redis OSS** / **for Valkey** (Valkey ถูกกว่าและเป็นตัวแนะนำ) | 2024-10 | https://aws.amazon.com/about-aws/whats-new/2024/10/amazon-elasticache-valkey/ |

⚠️ **Security Hub คือกับดักที่ร้ายที่สุด** — วันนี้ "AWS Security Hub" เป็นผลิตภัณฑ์ใหม่คนละตัว
ถ้าโจทย์หมายถึง "รวม finding + เช็ก best practice" ชื่อที่ถูกคือ **Security Hub CSPM**
คู่มือสอบยังเขียน "AWS Security Hub" → **เลี่ยงข้อที่ต้องแยกสองตัวนี้**

---

## §3 เลิกใช้ / ปิดรับลูกค้าใหม่ — **หมวดเสี่ยงที่สุด**

อะไรที่อยู่ในนี้แล้วโมเดลยังแนะนำ = เฉลยผิดแบบมั่นใจ

### 3a อยู่ใน scope ของ SAA-C03 **และ** ถูกลดสถานะ

| service/feature | สถานะ | วันที่ | URL |
|---|---|---|---|
| **AWS Snowball Edge** (ทั้ง Storage/Compute Optimized) | maintenance | 2025-11-07 | https://docs.aws.amazon.com/snowball/latest/developer-guide/snowball-edge-availability-change.html |
| **Amazon S3 Select**, **S3 Glacier Select** | maintenance | 2024-07-25 | https://docs.aws.amazon.com/AmazonS3/latest/userguide/selecting-content-from-objects.html |
| **S3 Object Lambda** | maintenance | 2025-11-07 | https://docs.aws.amazon.com/AmazonS3/latest/userguide/amazons3-ol-change.html |
| **Amazon Glacier (vault service เดิม)** | maintenance | 2025-11-07 | https://docs.aws.amazon.com/general/latest/gr/maintenance_services.html |
| **FSx File Gateway** | maintenance → ต่อ FSx for Windows File Server ตรง ๆ | 2024-10-28 | https://docs.aws.amazon.com/general/latest/gr/maintenance_services.html |
| **Storage Gateway Hardware Appliance** | ไม่ขายแล้ว support ถึง 2028-05 | 2025-05-12 | https://docs.aws.amazon.com/filegateway/latest/files3/hardware-appliance.html |
| **Amazon Elastic Transcoder** | **สิ้นสุด support — ไม่มี service แล้ว** | 2025-11-13 | https://aws.amazon.com/blogs/media/support-for-amazon-elastic-transcoder-ending-soon/ |
| **AWS Audit Manager** | maintenance | มีผล 2026-04-30 | https://docs.aws.amazon.com/audit-manager/latest/userguide/audit-manager-availability-change.html |
| **Amazon Kendra** | maintenance | 2026-06-30 | https://docs.aws.amazon.com/general/latest/gr/maintenance_services.html |
| **Directory Service – Simple AD** | maintenance | 2026-06-30 | https://docs.aws.amazon.com/general/latest/gr/maintenance_services.html |
| **SSM Incident Manager / Change Manager** | maintenance | 2025-11-07 | https://docs.aws.amazon.com/general/latest/gr/maintenance_services.html |
| **AWS WAF Classic** | **sunset — EOS 2026-10-07** | ประกาศ 2025-10-07 | https://docs.aws.amazon.com/general/latest/gr/sunset_services.html |
| **CloudTrail Lake** | ปิดรับลูกค้าใหม่ 2026-05-31 → ใช้ CloudWatch | ประกาศ 2026-03-31 | https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-lake-service-availability-change.html |

### 3b นอก scope หรือใกล้เคียง แต่โมเดลอาจยังแนะนำ

`AWS App Runner` (maintenance 2026-04-30 → ECS Express Mode) ·
`Inspector Classic` (EOS 2026-05-20) · `Macie Classic` (เลิกแล้ว) ·
`Kinesis Data Analytics for SQL` (ลบทิ้ง 2026-01-27) · `EC2-Classic` (เลิก 2022-08-15) ·
`Spot blocks` (เลิก 2022-12-31) · `NAT instance` (AMI EOL 2023-12-31 AWS แนะนำให้ย้าย) ·
`S3 RRS` (เอกสารเขียนว่า "ไม่แนะนำให้ใช้") · `CloudSearch` · `Cloud9` ·
`Migration Hub` / `Application Discovery Service` · `Timestream for LiveAnalytics` ·
`AWS Proton` / `App Mesh` (sunset 2026) · `Snowcone` (เลิก 2024-11-12)

---

## §4 Storage

- **storage class ปัจจุบันของ S3**: Standard · Intelligent-Tiering · Standard-IA · One Zone-IA · **Express One Zone** · **Glacier Instant Retrieval** · **Glacier Flexible Retrieval** · Glacier Deep Archive · Outposts · RRS(ไม่แนะนำ) — ข้อไหนเขียนว่า "S3 Glacier" เฉย ๆ = เก่า https://docs.aws.amazon.com/AmazonS3/latest/userguide/storage-class-intro.html
- **Glacier Instant Retrieval** (2021-11): archive ที่ดึงได้ระดับ **มิลลิวินาที** min 90 วัน — ทำลายสมมติฐาน "archive = ต้องรอเป็นนาที/ชม."
- **Intelligent-Tiering**: Archive Instant Access ที่ 90 วัน **auto ไม่ต้อง opt-in** (ทับบน Frequent 0d / Infrequent 30d) ส่วน Archive Access 90d และ Deep Archive Access 180d ยัง opt-in + async · **ไม่มี minimum duration** และ object < 128 KB **ไม่เสียค่า monitoring**
- **S3 Express One Zone** + directory bucket (2023-11): single-AZ, latency หลักมิลลิวินาทีเดียว, ถึง 2M req/s — คำตอบใหม่ของ "object storage latency ต่ำสุด วางคู่ compute ใน AZ เดียวกัน" https://docs.aws.amazon.com/AmazonS3/latest/userguide/s3-express-one-zone.html
- **S3 Tables** (2024-12): bucket ชนิดที่สาม, Iceberg แบบ managed
- **gp3** (2020-12): แยก IOPS/throughput ออกจากขนาด, ถูกกว่า gp2 ~20%, baseline 3000 IOPS / 125 MB/s → **สูตร "3 IOPS ต่อ GB" ของ gp2 เป็นความรู้ legacy แล้ว**
- **io2 Block Express** (2021-07): ถึง 256,000 IOPS / 4,000 MB/s, durability **99.999%**, ratio 1000 IOPS/GB — ไม่ใช่ io1 ที่ 99.8–99.9% และ 50 IOPS/GB อีกต่อไป
- **EBS Multi-Attach**: **io2** ทุก Region ที่มี io2 ถึง 16 instance (io1 มีแค่ 3 Region) → "Multi-Attach คือ io1" ผิดแล้ว
- **EFS**: class คือ Standard / **IA** / **Archive** (Archive มา 2023-11) ส่วน Regional กับ One Zone เป็น *ชนิดของ file system* ไม่ใช่ class · lifecycle ดีฟอลต์ →IA 30 วัน →Archive 90 วัน
- **FSx มี 4 ตัวพอดี**: Windows File Server · Lustre · **NetApp ONTAP** · **OpenZFS** — โมเดลที่ตอบแค่สองตัวแรกคือเก่า https://aws.amazon.com/fsx/when-to-choose-fsx/
- **FSx Intelligent-Tiering** (2025) สำหรับ Lustre และ OpenZFS

---

## §5 Databases

- **RDS Multi-AZ DB *cluster*** (2022-03): 1 writer + **2 standby ที่อ่านได้** ข้าม 3 AZ, failover **< 35 วิ** — ต่างจาก Multi-AZ *instance* (standby อ่านไม่ได้, 60–120 วิ) → "standby ของ Multi-AZ อ่านไม่ได้" **จริงเฉพาะแบบ instance** https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/multi-az-db-clusters-concepts.html
- **DynamoDB on-demand ลดราคา 50%** (มีผล 2024-11-01) → ตรรกะเดิม "provisioned + auto scaling ถูกกว่ามาก" อ่อนลงเยอะ · **ห้ามอ้างอัตราส่วนราคา**
- **DynamoDB Standard-IA table class** (2021-12): storage ถูกลง 60% แลกกับ throughput แพงขึ้น 20%+
- **RDS Extended Support เป็นค่าใช้จ่ายที่ auto-enroll** (MySQL 5.7 / PostgreSQL 11 ตั้งแต่ 2024-03-01) → "รันเวอร์ชันเก่าฟรีแค่ไม่ support" **ผิด มันคิดเงินอัตโนมัติ** https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/extended-support.html
- **Redshift Serverless** GA 2022-07 · **Aurora DSQL** GA 2025-05 · **Aurora PostgreSQL Limitless** GA 2024-10
- **ElastiCache for Valkey** (2024-10): Serverless ถูกกว่า 33% / node-based 20% → "ElastiCache = Redis หรือ Memcached" ไม่ครบแล้ว
- คู่มือสอบ **ไม่มี** DAX, MemoryDB, Timestream ในรายการ Database

---

## §6 Networking

- **Origin Access Control (OAC) แทน OAI** (2022-08) — OAC รองรับ **SSE-KMS**, SigV4, POST และทุก Region · **OAI ทำ SSE-KMS origin ไม่ได้** https://aws.amazon.com/about-aws/whats-new/2022/08/amazon-cloudfront-origin-access-control/
- **NLB รองรับ security group แล้ว** (2023-08-10) → "มีแต่ ALB ที่มี SG ส่วน NLB ต้องกรองที่ instance" **ผิดแล้ว** https://aws.amazon.com/about-aws/whats-new/2023/08/network-load-balancer-supports-security-groups/
- **NAT gateway คือสิ่งที่เอกสารแนะนำ** · NAT instance AMI EOL 2023-12-31 → ไม่ใช่แค่ trade-off แต่เป็นภาระการดูแล
- **gateway endpoint (S3, DynamoDB) ไม่มีค่าใช้จ่ายเพิ่ม** ส่วน interface endpoint คิดรายชั่วโมง + ต่อ GB — ยังเป็นตัวแยกแยะด้านราคาที่ใช้ได้ https://docs.aws.amazon.com/vpc/latest/privatelink/gateway-endpoints.html
- **PrivateLink เข้าถึง VPC *resource* ได้แล้ว** (2024-12) ผ่าน resource endpoint / service-network endpoint
- **AWS Cloud WAN** GA 2022-07 — ตัวเลือกที่สามนอกจาก TGW กับ peering สำหรับ global/multi-Region
- **CloudFront free tier: 1 TB/เดือน ถาวร** (ไม่จำกัด 12 เดือน) ตั้งแต่ 2021-12-01
- **ACM public cert export ได้แล้ว** (2025-06, $15/FQDN) — cert ที่ออกก่อน 2025-06-17 export ไม่ได้ → "ACM public cert เอาออกไปใช้บน EC2 ไม่ได้" ไม่จริงอีกแล้ว

---

## §7 Compute / serverless / containers

- **Lambda: memory 10,240 MB, /tmp 10,240 MB** → "512 MB /tmp, 3 GB memory" เก่า
- **SnapStart**: Java (2022-11) + **Python 3.12+ และ .NET 8+ (2024-11)** → "SnapStart มีแต่ Java" ผิด
- **Lambda บน Graviton2 (arm64)** price-performance ดีขึ้นถึง 34%
- **EKS Auto Mode** (2024-12) เป็นแนวทางที่ AWS แนะนำ **เหนือกว่า** EKS Fargate → "serverless Kubernetes = EKS + Fargate" ไม่ใช่คำตอบเดียวแล้ว
- **ECS Express Mode** (2025-11) — ตัวแทนของ App Runner
- **AWS แนะนำ Savings Plans เหนือ RI** (RI ยังอยู่: Standard 72% / Convertible 66%) · Compute SP 66% ครอบ EC2+Fargate+Lambda · EC2 Instance SP 72% ผูก family+Region
- **generation ปัจจุบันคือ Graviton4 รุ่นที่ 8** (C8g/M8g/R8g, 2024-09) → โมเดลที่ตอบ m5/c5/r5 ว่า "current" คือเก่า
- **Spot blocks ไม่มีแล้ว** → งาน 4 ชม. ห้ามขาด ต้อง On-Demand

---

## §8 Identity / security / governance

- **Resource control policies (RCPs)** ใน Organizations (2024-11-13) — จำกัดสิทธิ์ที่ **resource** ต่างจาก SCP ที่จำกัดที่ **principal** (ครอบ S3, STS, KMS, SQS, Secrets Manager) → "SCP คือ guardrail เชิงป้องกันตัวเดียวขององค์กร" ไม่จริงแล้ว https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_policies_rcps.html
- **GuardDuty Malware Protection for S3** GA 2024 → "GuardDuty สแกน object ใน S3 ไม่ได้" ผิด
- **Cognito feature plan: Lite / Essentials / Plus** (2024-11) → "advanced security features เป็นสวิตช์เปิดปิด" ไม่ใช่แล้ว
- **OpenSearch Serverless** GA 2023-01
- **Free Tier เปลี่ยนโครงสร้าง** (2025-07-15): บัญชีใหม่ได้เครดิต $200 + แผนฟรี 6 เดือน ส่วนโมเดล 12 เดือน/always-free ใช้กับบัญชีเก่าเท่านั้น

---

## §9 ⚠️ CONFLICT — จุดที่ "คำตอบที่ถูกวันนี้" ≠ "คำตอบที่ข้อสอบคาดหวัง"

**ออกข้อสอบ *รอบ ๆ* จุดเหล่านี้ ห้ามออก *ตรง* จุด**

| # | ข้อสอบคาดหวัง | ความจริงวันนี้ | ตัดสิน |
|---|---|---|---|
| **C1 Snow Family** | อยู่ใน scope ชัดเจน · คำตอบคลาสสิกของ "ย้าย 80 TB ผ่านลิงก์ 100 Mbps" | Snowball Edge ปิดรับลูกค้าใหม่ 2025-11-07 | **ข้อสอบยังเฉลย Snowball** → เลี่ยงหัวข้อขนถ่ายข้อมูลออฟไลน์ไปเลย หรือออกเป็นโจทย์คำนวณ bandwidth ซึ่งไม่ขึ้นกับเวอร์ชัน |
| **C2 Elastic Transcoder** | ยังอยู่ใน scope | **service หายไปแล้ว** 2025-11-13 | **ห้ามออกเด็ดขาด** เลี่ยงหมวด media ทั้งหมด |
| **C3 S3 Select** | คำตอบคลาสสิกของ "ดึงบางส่วนจากไฟล์ใหญ่ ลด data transfer" | ปิดรับลูกค้าใหม่ · AWS ชี้ไป Athena | **ห้ามออกข้อนี้** — เป็นจุดที่จะสอนของผิดมากที่สุด · ถ้าจะทดสอบ "query in place" ให้ใช้ Athena บน dataset ที่ partition แล้ว (ซึ่ง S3 Select ทำไม่ได้อยู่แล้ว) |
| **C4 S3 Object Lambda** | คำตอบคลาสสิกของ "ปกปิด PII ตอน retrieve" | ปิดรับลูกค้าใหม่ | เลี่ยง · ทดสอบแนวคิดเดียวกันด้วย CloudFront + Lambda@Edge |
| **C5 Aurora Serverless v1 vs v2** | "v1 pause ได้ เลือก v1 สำหรับ dev/test" | v1 ตาย · **v2 scale ลง 0 ได้แล้ว** | **ห้ามใช้ "scale to zero" เป็นตัวแยกแยะ v1/v2** |
| **C6 Security Hub** | ชื่อเดิม = รวม finding + เช็ก best practice | ตัวนั้นคือ **Security Hub CSPM** แล้ว | อธิบายด้วยหน้าที่ อย่าอิงชื่อ · ใส่ทั้งสองชื่อในเฉลย |
| **C7 Kendra / Simple AD / Audit Manager** | อยู่ใน scope ทั้งสามตัว | เข้า maintenance ปี 2026 | ใช้เป็น **distractor** ได้ แต่ **อย่าให้เป็นเฉลย** |
| **C8 Glacier vault vs storage class** | ใช้คำว่า "Glacier" ลอย ๆ | vault service เข้า maintenance · **storage class ยังใช้ได้ปกติ** | เขียนเต็มเสมอ: "S3 Glacier Flexible/Instant Retrieval / Deep Archive" ห้ามเขียน "archive ไป Glacier" |
| **C9 DynamoDB on-demand cost** | "provisioned ถูกกว่ามากสำหรับ traffic คงที่" | on-demand ถูกลง 50% ตั้งแต่ 2024-11 | ให้เหตุผลเชิงคุณภาพเท่านั้น (spiky → on-demand, steady → provisioned) **ห้ามอ้างตัวเลขอัตราส่วน** |
| **C10 public IPv4 charge** | ข้อเก่าอาจถือว่าฟรี | คิดเงินตั้งแต่ 2024-02-01 | **ตัดไปทางความจริงปัจจุบัน — ออกได้และควรออก** cost คือ 20% ของข้อสอบ |
| **C11 FSx File Gateway** | คำตอบคลาสสิกของ "SMB cache on-prem หลังบ้านเป็น FSx" | maintenance 2024-10-28 | ใช้ S3 File Gateway / Volume / Tape Gateway แทน |

---

## §10 ยังไม่ได้ยืนยัน — **ห้ามถือเป็นข้อเท็จจริง**

1. **เลขเวอร์ชันและวันที่ของคู่มือ SAA-C03** — เห็น `Version 1.1` ใน title ของผลค้นหาแต่เปิด PDF ไม่ได้ · มีค่าเพราะบอก snapshot ที่ข้อสอบอิง
2. วันที่ rename Route 53 ARC → Amazon ARC
3. **Aurora PostgreSQL Limitless ถูกยกเลิกหรือยัง** — ค้นแล้ว **ไม่พบว่ายกเลิก** พบแต่การอัปเดตต่อเนื่องถึง 2025-09 → ถือว่ายัง active
4. **Elastic Beanstalk**: ไม่อยู่ใน maintenance/sunset (ยังมีชีวิต) แต่ platform branch ที่อิง Amazon Linux 2 จะ retire 2026-06-30 — ยังไม่ได้ยืนยันจากหน้าหลัก
5. ตัวเลข request cost ของ S3 Express One Zone (เอกสารบอก 50% หน้าโฆษณาบอก 80%) → ใช้เลขจากเอกสาร หรือเลี่ยงตัวเลข
6. EC2 hibernation / Capacity Reservations / Capacity Blocks — **ยังไม่ได้ค้น**
7. Athena / Glue / Lake Formation — ยังไม่ได้ค้นนอกจากเรื่อง S3 Select
8. TGW vs VPC peering — ไม่พบว่าเปลี่ยน แต่ยังไม่ได้ยืนยันโดยตรง
9. KMS key types / IAM Access Analyzer / Detective — ยังไม่ได้ค้น

---

## กติกาการเขียนข้อสอบที่ตกผลึกจากทั้งหมดนี้

1. **Ban list — ห้ามเป็นเฉลย**: S3 Select · S3 Glacier Select · S3 Object Lambda · Elastic Transcoder · Snowball/Snowcone · Spot blocks · NAT instance · OAI · Aurora Serverless v1 · Glacier vault · FSx File Gateway · Simple AD · WAF Classic · Inspector Classic · Macie Classic · Kinesis Data Analytics for SQL · CloudSearch · EC2-Classic
2. **Ban topic**: S3 eventual consistency — ไม่มีอยู่จริงตั้งแต่ 2020-12-01
3. **วินัยเรื่องชื่อ**: ไล่ทุก draft กับ §2 ใช้ short name ทางการ
4. **ห้ามอ้างอัตราส่วนราคา** — §1.1, §4 (gp3), §5 (DynamoDB) ขยับหมด ใช้เหตุผลเชิงคุณภาพ
5. **เปิดสองหน้าใน §0 ซ้ำทุก batch** — รายการ maintenance ขยับตลอด
