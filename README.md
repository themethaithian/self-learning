# self-learning

> ย้ายเครื่อง / เปิด session ใหม่? อ่าน [`docs/HANDOFF.md`](docs/HANDOFF.md) ก่อน

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

เปิด **http://localhost:3000** แล้ววาง token ที่หน้า `/token`:

- **ไม่มีไฟล์ `.env`** (ค่า default ตอนรันครั้งแรก) → token คือ `local-dev-token`
- **มีไฟล์ `.env`** (เช่น copy จาก `.env.example` มา override พอร์ต) →
  `.env.example` เป็น template เดียวกับที่ prod ใช้ด้วย จึงตั้ง `API_BEARER_TOKEN`
  เป็น placeholder แบบเดียวกับรหัสผ่าน DB (`changeme-bearer-token`) ไม่ใช่ค่า
  local-dev — ต้องเช็คค่าจริงที่ compose จะใช้ด้วยคำสั่งนี้ แล้ว copy ค่าที่ได้
  ไปวางแทน:
  ```powershell
  docker compose config | Select-String API_BEARER_TOKEN
  ```

ค่าไหนก็ตาม เป็น dev-only ทั้งหมด ของจริงบน prod ไม่มีทางใช้ค่าพวกนี้

คำสั่งอื่นที่ใช้บ่อย:

```powershell
make logs
```
ดู log ของทุก service แบบ tail ต่อเนื่อง (ไม่มี `make`: `docker compose logs -f`)

```powershell
make dev-down
```
หยุดทุก service (ข้อมูลใน MySQL ยังอยู่) (ไม่มี `make`: `docker compose down`)

```powershell
make dev-reset
```
หยุด + ลบ MySQL volume ทิ้งทั้งหมด (เริ่มนับหนึ่งใหม่) (ไม่มี `make`: `docker compose down -v`)

```powershell
make seed
```
รัน import ใหม่อีกครั้ง หลังเพิ่ม lesson JSON ใหม่ใน `content/` — ปลอดภัย, idempotent
(upsert by slug ซ้ำได้ไม่ซ้ำข้อมูล) (ไม่มี `make`: `docker compose run --rm seed`)

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
