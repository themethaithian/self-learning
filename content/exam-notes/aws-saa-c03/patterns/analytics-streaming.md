# Patterns — Analytics & Streaming

## KDS → Lambda → DynamoDB: anonymize PII in transit ก่อนลง storage

**Signal:** real-time data จากหลาย sources · มี PII ต้อง anonymize **ก่อน** landing
ใน storage · ปลายทางเป็น NoSQL

**คำตอบ:** Kinesis Data Streams รับ stream → Lambda แปลง/anonymize ระหว่างทาง
(in transit) → เก็บผลลง DynamoDB

**เหตุผล:** KDS เป็น real-time streaming ที่ scale ได้มหาศาลและต่อกับ Lambda ได้ตรง ๆ
จุดชนะคือข้อมูลถูกแปลง "ระหว่างทาง" — PII ดิบไม่เคย landing ใน data store ปลายทาง
ตรงกับ requirement เป๊ะ ส่วน DynamoDB ปิด requirement ข้อ NoSQL

**Trap ในข้อนี้:**

- ลง S3 ก่อนแล้วใช้ S3 trigger + Lambda แปลง → PII ดิบ landing แล้ว
  ([store-then-transform](../traps.md#store-then-transform))
- ลง DynamoDB แล้วใช้ DynamoDB Streams แปลง item ที่เขียนแล้ว → ดิบ landing
  เหมือนกัน แถมแจก FullAccess ([fullaccess-red-flag](../traps.md#fullaccess-red-flag))
- Firehose → Redshift → Redshift เป็น relational DW ไม่ใช่ NoSQL
  ([service-type-mismatch](../traps.md#service-type-mismatch))

**เก็บเพิ่ม:** ตัว KDS stream เอง retain record ดิบไว้ (default 24 ชม. ขยายได้ถึง
365 วัน) — โจทย์สาย encryption จะถามเรื่องเปิด server-side encryption (KMS) ที่ stream
เพื่อคุ้มครอง at rest ตรงนี้ · Firehose ก็มี built-in Lambda transformation แต่ปลายทาง
หลักคือ S3 / Redshift / OpenSearch / Splunk (+ HTTP endpoint / SaaS partner) —
**ส่งเข้า DynamoDB ตรง ๆ ไม่ได้** ปลายทางต้องเป็น DynamoDB ให้ใช้ KDS + consumer เขียนเอง

**ศัพท์:** PII · anonymization · in transit / at rest · least privilege · data lake ·
data warehouse · NoSQL → [glossary](../glossary.md)

<sub>ที่มา: TD-20260730-kds-pii · refs: aws.amazon.com/kinesis/data-streams/ ·
docs.aws.amazon.com/lambda/latest/dg/with-kinesis.html · 2026-07-30</sub>
