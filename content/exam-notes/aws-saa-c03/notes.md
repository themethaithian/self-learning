# โน้ตเตรียมสอบ AWS SAA-C03 — เทคนิคอ่านโจทย์ + ศัพท์สำคัญ

<!-- ids: TD-20260730-kds-pii, TD-20260730-asg-scaling, TD-20260730-s3-hotlink, TD-20260730-fsx-sharepoint, TD-20260730-fsx-ontap, TD-20260730-saml-federation -->

จัดการโดย skill `/aws-exam-note` — วางโจทย์ + เฉลย แล้ว Claude แปล สรุป
และผสานเข้าไฟล์นี้ให้ · จุดที่เคยถามซ้ำ/ไม่เข้าใจอยู่ที่ [review-again.md](review-again.md)

**สรุปจากการทำโจทย์ 6 ข้อ:** Kinesis anonymization · Auto Scaling · S3 hotlinking ·
FSx SharePoint · FSx ONTAP trading app · SAML federation — อัปเดตล่าสุด 2026-07-30

## ส่วนที่ 1 — เทคนิคอ่านโจทย์

### 1.1 กฎเหล็ก: อ่านคำบรรยาย architecture ก่อนดูตัวเลือก

ข้อสอบชอบซ่อน **disqualifier** (ตัวตัดสิทธิ์) ไว้ในประโยคบรรยายระบบตอนต้นโจทย์
ไม่ใช่ในคำถาม

ตัวอย่างจริง: โจทย์ Auto Scaling เขียนว่า "instances with **different instance types
and sizes**" — ประโยคนี้ตัด Predictive scaling ทิ้งทันที เพราะ Predictive
ตั้งสมมติฐานว่า ASG เป็น homogeneous (ทุกเครื่องขีดความสามารถเท่ากัน)
ถ้าเครื่องต่างขนาดปนกัน forecast จะเพี้ยน

วิธีทำ: อ่านโจทย์รอบแรก ขีดเส้นใต้ทุกคุณสมบัติของระบบ ก่อนมองตัวเลือกเลย

- multi-AZ / cross-region
- different instance types
- Windows / Linux
- stateless / stateful
- encrypted at rest
- on-premises + direct connection

แต่ละอันมักผูกกับข้อจำกัดของบริการใดบริการหนึ่งเสมอ

### 1.2 จับ "คำชี้ขาด" ในคำถาม

| คำ | ความหมายที่ AWS ต้องการ |
|---|---|
| MOST operationally efficient | แก้ปัญหาได้ด้วยความซับซ้อนน้อยที่สุด ไม่ใช่ "ฉลาด/อัตโนมัติที่สุด" |
| MOST cost-effective | ถูกที่สุดที่ยังตอบโจทย์ครบ |
| MOST secure | least privilege, encryption, no public access |
| LEAST operational overhead | เลือก managed/serverless |
| MOST effective | แก้ที่ต้นเหตุจริง ไม่ใช่แก้ปลายเหตุ |

บทเรียนจากที่พลาด: เห็น "MOST operationally efficient" แล้วรีบตอบ Predictive
เพราะดูฉลาดกว่า — แต่ถ้าโจทย์ให้เวลามาเป๊ะ ๆ แล้ว การเอา ML มาทายสิ่งที่เรารู้อยู่แล้ว
คือ over-engineering

### 1.3 ระวังคำที่สลับกัน (ตัวลวงคลาสสิก)

| ในโจทย์ขอ | ตัวเลือกลวงเขียน |
|---|---|
| cross-AZ | cross-region replication |
| NoSQL | Redshift (data warehouse) |
| block storage | file storage / object storage |
| Windows / SMB | NFS file share |

เห็นคำพวกนี้สลับกันเมื่อไหร่ ให้สงสัยไว้ก่อนว่าเป็นตัวลวง

### 1.4 หลักการที่ตัดตัวเลือกได้เร็ว

- **transform-before-store** — โจทย์บอกว่าต้องแปลง/ปกปิดข้อมูล **ก่อน** เก็บ
  ให้เอาตัวประมวลผล (Lambda) ไปวางระหว่าง stream กับ storage ถ้าข้อไหนแปลง
  **หลังจาก** เขียนลง storage ไปแล้ว (เช่น DynamoDB Streams, S3 trigger) →
  ตัดทิ้งทันที ไม่ว่าจะดูดีแค่ไหน
- **reactive vs proactive** — ถ้าอาการคือ "ช้าตอนเริ่ม แล้วดีขึ้นทีหลัง" แปลว่า
  scaling มันไล่ตามโหลดอยู่ → ต้องแก้ด้วยวิธี proactive (เตรียมล่วงหน้า)
  ไม่ใช่ dynamic scaling ที่รอ metric พุ่งก่อน
- **least privilege เป็นธง** — เห็น `AmazonDynamoDBFullAccess` หรือ managed policy
  แบบ Full Access ในตัวเลือก = สัญญาณว่าข้อนั้นน่าจะผิด
- **เปลี่ยน IP ง่าย** — วิธีที่พึ่ง IP (NACL, Security Group) หรือ referrer header
  จะแพ้เสมอ เพราะปลอม/เปลี่ยนได้ง่าย ต้องแก้ที่ระดับ authentication แทน

## ส่วนที่ 2 — Storage (เรื่องที่ออกสอบเยอะที่สุด)

### 2.1 สามประเภทหลัก

| | Block Storage | File Storage | Object Storage |
|---|---|---|---|
| เปรียบเหมือน | ฮาร์ดดิสก์เปล่า | ตู้เอกสารมีลิ้นชัก | บริการรับฝากของ |
| หน่วยข้อมูล | block (ข้อมูลดิบ) | ไฟล์ในโฟลเดอร์ | object + metadata |
| เข้าถึงยังไง | ต่อเป็นไดรฟ์ (C:, /dev/sda) | mount เป็น network drive | HTTP API |
| ใครจัดโครงสร้าง | OS ของเรา | ตัว storage เอง | ไม่มีโครงสร้าง (flat) |
| แชร์หลายเครื่อง | ปกติไม่ได้ | ได้ (จุดขายหลัก) | ได้ ผ่าน HTTP |
| แก้ไขบางส่วน | ได้ | ได้ | ไม่ได้ ต้องเขียนทับทั้งชิ้น |

ความสัมพันธ์แบบชั้น: Object ทับ File ทับ Block — ยิ่งลงล่าง = เร็วกว่า
ควบคุมละเอียดกว่า แต่แชร์ยาก แพงกว่า · ยิ่งขึ้นบน = ถูกกว่า แชร์ง่าย ขยายง่าย
แต่ช้ากว่า ยืดหยุ่นน้อยกว่า

### 2.2 ตารางเลือกบริการ

| โจทย์เขียนว่า | ประเภท | บริการ | โปรโตคอล |
|---|---|---|---|
| database, boot volume, high IOPS | Block | EBS | — |
| block storage + HA ข้าม AZ | Block | FSx for NetApp ONTAP | iSCSI |
| shared file, Linux, NFS | File | EFS | NFS |
| Windows, SMB, AD, SharePoint | File | FSx for Windows File Server | SMB |
| HPC, machine learning | File | FSx for Lustre | Lustre |
| รูป วิดีโอ backup static site | Object | S3 / Glacier | HTTP |

จุดตัดสิน: เห็นคำว่า **block + ข้าม AZ → ONTAP ตัวเดียว** เพราะเป็น FSx เจ้าเดียว
ที่พูด iSCSI ได้ · EBS = block ก็จริง แต่อยู่ได้แค่ AZ เดียว และปกติต่อเครื่องเดียว
(ยกเว้น Multi-Attach ของ io1/io2 ต่อได้ถึง 16 เครื่อง — แต่ก็ยังอยู่ใน AZ เดียวอยู่ดี)

### 2.3 SMB vs NFS

| | SMB | NFS |
|---|---|---|
| โลกของ | Windows | Linux / Unix |
| path | `\\server\share` | `server:/export/path` |
| สิทธิ์ | Active Directory / Windows ACL | UID/GID ของ Unix |

เหมือนปลั๊กไฟคนละมาตรฐาน — ทำหน้าที่เดียวกันแต่เสียบข้ามกันไม่ได้

### 2.4 ทำไม S3 แทน file share ไม่ได้

- S3 เข้าถึงผ่าน HTTP API เท่านั้น mount เป็นไดรฟ์ตรง ๆ ไม่ได้
- แก้ไขบางส่วนไม่ได้ — อยากเปลี่ยนตัวอักษรเดียวในไฟล์ 1 GB ต้องอัปโหลดใหม่ทั้ง
  1 GB → database รันบน S3 ไม่ได้
- ระบบสิทธิ์เป็น IAM policy ผสานกับ Active Directory ไม่ได้
- แอปเก่าอย่าง SharePoint คาดหวัง path `\\server\share` → ใช้ S3 ไม่ได้

### 2.5 S3 hotlinking — ป้องกันคนอื่นดึงรูปไปใช้

- bucket policy กรอง referrer → ไม่พอ เพราะ header ปลอมได้
- คำตอบคือ **ปิด public read + ใช้ pre-signed URL ที่มีวันหมดอายุ**
- pre-signed URL = URL ที่เซ็นด้วย credentials ของเจ้าของ ให้สิทธิ์ดาวน์โหลด
  จำกัดเวลา ต้องระบุ bucket, object key, HTTP method, และเวลาหมดอายุ

## ส่วนที่ 3 — Auto Scaling

| โจทย์บอกว่า | ตอบ |
|---|---|
| ระบุเวลา/ตารางชัดเจน ("9 to 5", "ทุกวันจันทร์", "สิ้นเดือน") | Scheduled |
| มี pattern ซ้ำ ๆ แต่ไม่ระบุเวลา/ปริมาณ ("recurring", "cyclical") | Predictive |
| โหลดคาดเดาไม่ได้ / spike กะทันหัน | Dynamic (target tracking) |

**Scheduled scaling** — สร้าง scheduled action ระบุเวลาเริ่ม + ค่า min/max/desired
ใหม่ ทำได้ทั้งครั้งเดียวและแบบ recurring (cron)

**ข้อจำกัดของ Predictive ที่ต้องจำ:**

- ต้องมี historical data ≥ 24 ชม. (แนะนำ 14 วัน)
- scale out อย่างเดียว ไม่ scale in → ต้องใช้คู่กับ dynamic เสมอ
- สมมติว่า ASG เป็น homogeneous → เจอ "different instance types and sizes" คือตัดทิ้ง
  (⚠ ข้อนี้เป็นการตีความจากเฉลย TD ไม่ใช่ข้อจำกัดที่ AWS ประกาศ — ใช้เป็น heuristic
  ประกอบคำชี้ขาดใน 1.2 ไม่ใช่กฎเด็ดขาด)

**Dynamic ตาม memory** — จำไว้ว่า memory utilization ไม่ใช่ metric default ของ
CloudWatch ต้องติดตั้ง CloudWatch agent เพิ่ม = ops overhead มากขึ้น

## ส่วนที่ 4 — Identity & Access

### 4.1 Active Directory (AD)

ระบบจัดการ identity ขององค์กรของ Microsoft รันบน Windows Server เก็บ users,
groups, computers, printers, file shares

ศัพท์ที่ต้องรู้:

- **Domain** — ขอบเขตการปกครองของ AD เช่น `company.local`
- **Domain Controller (DC)** — เซิร์ฟเวอร์ที่รัน AD
- **LDAP** — โปรโตคอล query ข้อมูลจาก AD (ในไดอะแกรมเขียนว่า "LDAP identity store")
- **Kerberos** — โปรโตคอลยืนยันตัวตนภายในองค์กร
- **AD FS** — ส่วนเสริมที่ทำให้ AD คุย SAML กับระบบภายนอกได้

จุดสำคัญ: AD พูด Kerberos/LDAP ซึ่ง AWS ไม่เข้าใจ เลยต้องมี AD FS มาแปลงเป็น SAML

### 4.2 IdP vs SP

- **IdP (Identity Provider)** — ตอบคำถาม "คุณคือใคร" → AD FS, Okta, Google, Facebook
- **SP (Service Provider)** — ตอบคำถาม "คุณทำอะไรได้" → AWS

เปรียบเทียบ: IdP = กรมการปกครองที่ออกบัตรประชาชน, SP = ธนาคารที่ขอดูบัตร

### 4.3 SAML

มาตรฐานเปิดที่กำหนดรูปแบบการคุยกันระหว่าง IdP กับ SP ตัวมันคือ XML ที่เซ็น
ลายเซ็นดิจิทัล เรียกว่า **SAML assertion**

แปลเป็นภาษาคน: "ผมคือ AD FS ของบริษัท ขอรับรองว่าคนนี้ชื่อสมชาย ผ่านการยืนยัน
ตัวตนแล้ว ควรได้ role Developer — นี่ลายเซ็นผม"

หัวใจ: AWS ไม่ต้องรู้ password เลย แค่ตรวจว่าลายเซ็นตรงกับ certificate ของ IdP
ที่ลงทะเบียนไว้

### 4.4 Flow การ login (ตามไดอะแกรม 7 ขั้น)

| ขั้น | เกิดอะไร |
|---|---|
| 1 | ผู้ใช้เปิดเบราว์เซอร์เข้าหน้า portal ของ IdP |
| 2 | IdP ตรวจตัวตนกับ LDAP identity store (= AD) |
| 3 | IdP ส่ง SAML assertion กลับมาให้เบราว์เซอร์ |
| 4 | เบราว์เซอร์ POST assertion ไปที่ SAML endpoint ของ AWS (`signin.aws.amazon.com/saml`) |
| 5 | AWS ตรวจลายเซ็น → เรียก STS ด้วย `AssumeRoleWithSAML` |
| 6 | endpoint ส่ง redirect กลับมา |
| 7 | เข้า AWS Management Console |

สามจุดที่ต้องจำ:

- password ไม่เคยเดินทางมาถึง AWS
- AWS ไม่มี IAM user ของคนนั้น มีแต่ IAM Role ที่รอให้มา assume
- ได้ temporary credentials จาก STS มีวันหมดอายุ (ปกติ 1 ชม. ตั้งได้ถึง 12 ชม.)

### 4.5 IAM

องค์ประกอบ: Users / Groups / Policies / Roles

Policy = JSON ระบุ Effect (Allow/Deny), Action (`service:operation`),
Resource (ARN), Condition

กฎการตัดสินสิทธิ์:

1. Default = DENY
2. มี Allow ที่ตรง → อนุญาต
3. มี explicit Deny ที่ไหนก็ตาม → ห้ามทันที ชนะทุกอย่าง

User vs Role:

| | IAM User | IAM Role |
|---|---|---|
| ผูกกับ | คนหนึ่งคน | ไม่ผูกกับใคร |
| Credentials | ถาวร | ชั่วคราว มีวันหมดอายุ |
| ใช้ยังไง | login ตรง ๆ | ต้อง assume ก่อน |

User = บัตรพนักงานของคุณ / Role = เสื้อกาวน์แขวนหน้าห้องแล็บ ใครมีสิทธิ์ก็หยิบมาใส่ได้

Role มีสองส่วน: **Trust policy** (ใครสวมได้) + **Permission policy** (สวมแล้วทำอะไรได้)

ห้ามฝัง access key ในโค้ด — ใช้ Role แทนเสมอ:

- EC2 → Instance Profile
- Lambda → Execution Role
- คนจาก AD → assume role ผ่าน SAML
- ข้าม account → cross-account role

**STS** = บริการออก temporary credentials — เห็นคำว่า "temporary credentials"
ในข้อสอบให้นึกถึง STS ทันที

### 4.6 เลือกวิธี authentication

| ผู้ใช้เป็นใคร | มีอะไรอยู่แล้ว | ตอบ |
|---|---|---|
| พนักงานในองค์กร | Active Directory | SAML 2.0 + AD FS |
| พนักงาน + หลาย AWS account | Active Directory | IAM Identity Center (ชื่อเดิม AWS SSO) |
| ลูกค้าทั่วไปของแอป | ไม่มี | Cognito User Pool |
| ลูกค้า login ผ่าน Google/FB | ไม่มี | Cognito Identity Pool / Web Identity Federation |
| server, script (ไม่ใช่คน) | — | IAM Role |

Cognito สองส่วนที่คนสับสน:

| | User Pool | Identity Pool |
|---|---|---|
| ทำหน้าที่ | authentication | authorization เข้า AWS |
| คือ | ฐานข้อมูลผู้ใช้ + หน้า login | ตัวแลก token → AWS credentials |
| ผลลัพธ์ | JWT token | AWS temporary credentials |

Cognito ใช้กับ **customer** (ผู้ใช้ภายนอก) ไม่ใช่ **employee** (พนักงานภายใน)

ศัพท์ที่ต้องแยกให้ออก:

- **Authentication (AuthN)** = คุณคือใคร → หน้าที่ของ IdP
- **Authorization (AuthZ)** = คุณทำอะไรได้ → หน้าที่ของ IAM Role/Policy

## ศัพท์เบ็ดเตล็ดที่เจอ

| ศัพท์ | ความหมาย |
|---|---|
| On-premises | ระบบที่ตั้งในสถานที่ของบริษัทเอง ตรงข้ามกับ cloud |
| Direct Connect | สายเชื่อมต่อเฉพาะระหว่างออฟฟิศกับ AWS ไม่ผ่านอินเทอร์เน็ตสาธารณะ |
| SharePoint | ซอฟต์แวร์ Microsoft ทำระบบเอกสาร/intranet องค์กร รันบน Windows เท่านั้น |
| File share | โฟลเดอร์ที่แชร์ให้เครื่องอื่นในเครือข่ายใช้ร่วมกัน |
| ARN | Amazon Resource Name — รหัสประจำตัวของทุกสิ่งใน AWS |
| Kinesis Data Streams | รับ stream แบบเรียลไทม์จากหลาย producer (ในตัว stream ข้อมูลถูก retain ไว้ default 24 ชม. ถึง 365 วัน — เข้ารหัสด้วย SSE/KMS ได้) |
| Firehose | ส่ง stream เข้าปลายทาง (S3/Redshift/OpenSearch/Splunk + HTTP endpoint) มี buffer → near-real-time · ส่งเข้า DynamoDB ตรง ๆ ไม่ได้ |
| Redshift | data warehouse ไม่ใช่ NoSQL |
| CloudTrail | บันทึกว่าใครทำอะไรเมื่อไหร่ ตรวจย้อนหลังได้ |

## Checklist ก่อนตอบทุกข้อ

1. อ่านคำบรรยาย architecture ให้ครบทุกวลี ขีดเส้นใต้คุณสมบัติของระบบ
2. หาคำชี้ขาดในคำถาม (MOST xxx)
3. เช็คว่ามี disqualifier ซ่อนอยู่ไหม (different types, cross-AZ vs region,
   Windows vs Linux)
4. ตัดตัวเลือกที่แปลงข้อมูลหลังเก็บ / reactive / Full Access policy /
   พึ่ง IP-referrer
5. เลือกอันที่ **พอดี** กับ signal ไม่ใช่อันที่ฉลาดที่สุด
