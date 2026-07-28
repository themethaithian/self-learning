# self-learning

## รันเองบนเครื่อง

ต้องมีแค่ Docker Desktop (Windows/Mac/Linux) — ไม่ต้องลง Go, Node, หรือ MySQL เอง

รันคำสั่งเดียวจาก repo root:

```powershell
make dev
```

ถ้าเครื่องไม่มี `make` ใช้คำสั่งดิบแทนได้:

```powershell
docker compose up -d --build
```

คำสั่งนี้จะ build image ของ `api` (Go) และ `web` (Next.js static export + Caddy),
สตาร์ท MySQL, รอจน MySQL/API healthy แล้วรัน `seed` (one-shot import curriculum +
lessons) ให้อัตโนมัติ — ไม่ต้องรันคำสั่ง import เองเลย

รอบแรก (build จาก image ที่ base image ยังไม่เคย pull มาก่อน + MySQL init
volume ใหม่) วัดจริงบนเครื่องนี้ได้ประมาณ **2-3 นาที** (ส่วนใหญ่หมดไปกับ
MySQL first-boot ที่ใช้เวลาราว 100 วินาที ก่อน API/seed จะเริ่มทำงานได้ —
เป็นเหตุผลที่ mysql healthcheck ใน `docker-compose.yml` ตั้ง `start_period`
ไว้ยาวถึง 90s). รอบถัดไป (มี volume/cache แล้ว) เร็วกว่านี้มาก

เปิด **http://localhost:3000** แล้ววาง token `local-dev-token` ที่หน้า `/token`
— ค่านี้ใช้ได้ทั้งตอนไม่มีไฟล์ `.env` เลย (compose ตั้ง default ไว้เป็นค่านี้) และตอน
copy `.env.example` มาเป็น `.env` (ค่าใน `.env.example` ก็เป็น `local-dev-token`
เหมือนกัน) เว้นแต่คุณจะไปแก้ `API_BEARER_TOKEN` ในไฟล์ `.env` ของตัวเองเป็นค่าอื่น
ก็ต้องใช้ค่านั้นแทน — นี่คือ dev-only default เท่านั้น ของจริงบน prod ไม่มีทางใช้ค่านี้

คำสั่งอื่นที่ใช้บ่อย:

```powershell
make logs
```
ดู log ของทุก service แบบ tail ต่อเนื่อง

```powershell
make dev-down
```
หยุดทุก service (ข้อมูลใน MySQL ยังอยู่)

```powershell
make dev-reset
```
หยุด + ลบ MySQL volume ทิ้งทั้งหมด (เริ่มนับหนึ่งใหม่)

```powershell
make seed
```
รัน import ใหม่อีกครั้ง หลังเพิ่ม lesson JSON ใหม่ใน `content/` — ปลอดภัย, idempotent
(upsert by slug ซ้ำได้ไม่ซ้ำข้อมูล)

### แก้ปัญหาที่เจอบ่อย

**Port ชนกับโปรแกรมอื่นในเครื่อง** — สร้างไฟล์ `.env` ที่ repo root (copy จาก
`.env.example`) แล้ว override `DB_PORT`, `API_PORT`, หรือ `WEB_PORT` เป็น port
อื่นที่ว่าง

**Docker Desktop ยังไม่ได้เปิด** — `docker compose up` จะ error ประมาณ
`error during connect... The system cannot find the file specified` หรือ
`pipe/dockerDesktopLinuxEngine` — เปิด Docker Desktop ก่อนแล้วรอจน tray icon
ขึ้นสถานะ running ค่อยรันคำสั่งใหม่

**หน้าเว็บขึ้นแต่ไม่มี curriculum/lessons โผล่มา** — เช็คว่า one-shot import
ทำงานจบจริงหรือยังด้วย:
```powershell
docker compose logs seed
```
ถ้า log จบด้วย `"lesson import complete"` แปลว่า import สำเร็จแล้ว ถ้าไม่เจอบรรทัดนี้
หรือเจอ error ให้รัน `make seed` ใหม่อีกรอบ
