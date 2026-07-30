# Keyword Map — เห็นคำนี้ในโจทย์ ให้นึกถึงอะไร

ตารางกวาดตาเร็ว: signal (คำ/วลีในโจทย์) → สิ่งที่มักเป็นคำตอบ → ข้อควรระวัง
รายละเอียดเต็มอยู่ใน [patterns/](patterns/)

## Analytics & Streaming

| Signal ในโจทย์ | นึกถึง | ระวัง |
|---|---|---|
| real-time ingest / streaming จากหลายพัน sources | Kinesis Data Streams (KDS) | Firehose เป็น near real-time (มี buffer) ไม่ใช่ real-time แท้ |
| transform/anonymize ข้อมูล**ระหว่างทาง**ก่อนลง storage | KDS หรือ Firehose + Lambda (แปลง in transit) | option ที่เก็บดิบก่อนค่อยแปลง = ผิด ([trap](traps.md#store-then-transform)) |
| stream ปลายทาง S3 / Redshift / OpenSearch / Splunk (+ HTTP endpoint / SaaS partner) แบบ managed | Amazon Data Firehose | Firehose ส่งเข้า DynamoDB ตรง ๆ ไม่ได้ |

## Database

| Signal ในโจทย์ | นึกถึง | ระวัง |
|---|---|---|
| NoSQL / key-value / latency หลัก ms | DynamoDB | Redshift/RDS/Aurora เป็น relational — ผิดทันทีถ้าโจทย์สั่ง NoSQL ([trap](traps.md#service-type-mismatch)) |
