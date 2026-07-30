# Trap Patterns — ตัวเลือกหลอกที่ออกซ้ำ ๆ

รู้จัก trap แล้วตัด choice ได้เร็วขึ้น — เรียงตามตัวนับ "เจอ" มาก→น้อย (ยิ่งบน = ยิ่งออกบ่อย)

## store-then-transform

โจทย์สั่งว่าข้อมูลอ่อนไหวต้องถูกแปลง/ปกปิด **ก่อน** ลง storage แต่ตัวเลือกเก็บ
ข้อมูลดิบลงก่อนแล้วค่อยแปลงทีหลัง เช่น ลง S3 แล้วใช้ S3 event trigger Lambda,
หรือลง DynamoDB แล้วให้ DynamoDB Streams แปลง item ที่เขียนไปแล้ว — ผิดเสมอ
เพราะข้อมูลดิบ "landing" ไปแล้วหนึ่งจุด ต่อให้ชั่วคราวก็เพิ่มความเสี่ยงข้อมูลรั่ว

**เจอ:** 1 ครั้ง

## fullaccess-red-flag

ตัวเลือกที่แจกสิทธิ์กว้าง เช่น `AmazonDynamoDBFullAccess` ทั้งที่งานใช้แค่บาง action —
ขัดหลัก least privilege (ดู [glossary](glossary.md)) ข้อสอบชอบใส่คำว่า FullAccess
เป็นตัวบอกใบ้ว่า option นั้นผิด แม้ส่วนอื่นของ option จะฟังดูเข้าท่า

**เจอ:** 1 ครั้ง

## service-type-mismatch

โจทย์ระบุประเภทระบบตรงตัว (เช่น "NoSQL database") แต่ตัวเลือกจบที่ service
คนละประเภท (เช่น Redshift ซึ่งเป็น relational data warehouse) — architecture
ช่วงต้นจะสวยแค่ไหนก็ผิด เช็คปลายทางของทุก option กับ requirement ก่อนเสมอ

**เจอ:** 1 ครั้ง
