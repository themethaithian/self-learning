# Glossary — ศัพท์ที่เจอในโจทย์/เฉลย

เรียง A→Z · คำอธิบายไทยสั้น ๆ พอเข้าใจในบริบทข้อสอบ

| ศัพท์ | คำอธิบาย |
|---|---|
| anonymization | การแปลง/ปกปิดข้อมูลให้ระบุตัวบุคคลไม่ได้ (ลบชื่อ, mask เลขบัตร ฯลฯ) |
| data lake | ที่เก็บข้อมูลดิบทุกรูปแบบรวมกันขนาดใหญ่ — ในข้อสอบมักหมายถึง S3 |
| data warehouse | ฐานข้อมูลเชิงวิเคราะห์สำหรับ query ข้อมูลมหาศาล (Redshift) — เป็น relational ไม่ใช่ NoSQL |
| in transit / at rest | ข้อมูล "ระหว่างเดินทาง" (ในสาย/ใน stream) vs "ตอนถูกเก็บ" (ในดิสก์/DB) — โจทย์ security ชอบแยกสองจังหวะนี้ |
| least privilege | หลักให้สิทธิ์น้อยที่สุดเท่าที่งานต้องใช้ — เห็น FullAccess ในตัวเลือกให้สงสัยไว้ก่อน |
| NoSQL | ฐานข้อมูลที่ไม่ใช่ตาราง relational เช่น key-value (DynamoDB), document — เด่นเรื่อง scale และ latency ต่ำ |
| PII (Personally Identifiable Information) | ข้อมูลที่ระบุตัวบุคคลได้ เช่น ชื่อ ที่อยู่ เลขบัตร — โจทย์ compliance/healthcare ใช้บ่อย |
