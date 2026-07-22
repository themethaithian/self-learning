# Week 5 — DSA ครบวงจร

**เป้าหมาย**: เลือก pattern → ทำโจทย์จาก bank → ส่ง approach + โค้ด Go → LLM รีวิว → สถิติรายแพทเทิร์นขยับ

---

## T22 — Practice domain + migration + cmd/import-dsa `[go-implementer]` ~45 นาที
- **Scope**: `internal/practice/domain` (Problem, Attempt, ReviewResult VO), `migrations/004_practice.sql`, `cmd/import-dsa` อ่าน `content/dsa-problems/<pattern>/*.json`
- **Acceptance**: import idempotent, tests ของ domain constructors
- **Review focus**: ReviewResult VO validate verdict/score สอดคล้องกัน (verdict=correct แต่ score=1 ควร error ไหม — ดูว่า implementer ตัดสินใจยังไง)
- Status: `todo`

## T23 — DSA review prompt + attempt endpoint `[go-implementer]` ~45 นาที
- **Goal**: `POST /dsa/problems/{id}/attempts` — Haiku รีวิว correctness / complexity / idiomatic Go
- **Scope**: ขยาย port LLM เพิ่ม `ReviewSolution`, prompt JSON-mode (verdict, complexity_time/space, score 0–5, review_md ภาษาไทย), ยิง event `AttemptReviewed` → ReviewCard(dsa_pattern) + measurement
- **Acceptance**: ส่งโค้ดผิดจริง ๆ แล้ว verdict ไม่ใช่ correct (ทดสอบมือ 2–3 เคส), LLM ล่มคำตอบ user ไม่หาย, tests mock LLM
- **Review focus**: prompt แนบ reference_approach_md ให้ Haiku เทียบไหม (ควรแนบ — ทำให้รีวิวแม่นขึ้นมาก), max_tokens พอสำหรับรีวิวยาว
- Status: `todo`

## T24 — DSA UI `[go-implementer]` ~45 นาที
- **Scope**: `web/app/dsa/` — หน้า pattern list (+สถิติ), หน้าโจทย์ (prompt + examples), form: approach (textarea) + โค้ด Go (textarea monospace), หน้ารีวิวผล
- **Acceptance**: flow ครบด้วยโจทย์จริงจาก C3, timer 30 นาทีใช้ร่วม
- **Review focus**: ลำดับ UX ต้องบังคับเขียน approach ก่อนถึงเปิดช่องโค้ด (กติกาการฝึกที่ตั้งไว้), code textarea ใช้ font mono + tab ทำงาน
- Status: `todo`

## T25 — Per-pattern stats `[go-implementer]` ~30 นาที
- **Scope**: `GET /dsa/patterns` รวม accuracy/attempts/avg score ต่อ pattern, UI heat-bar ในหน้า pattern list
- **Acceptance**: ตัวเลขตรงกับ attempts จริง, pattern ยังไม่แตะแสดง "ยังไม่เริ่ม"
- **Review focus**: aggregate query อ่านง่ายไหม (GROUP BY เดียวจบ), ไม่คำนวณฝั่ง Go ทีละ attempt
- Status: `todo`

## C3 — Content: DSA bank 3 patterns × 5 ข้อ `[lesson-writer → lesson-verifier]`
- arrays-hashing, two-pointers, sliding-window — pattern ละ easy×2 medium×2 hard×1
- Verifier เช็คเพิ่ม: ไม่ซ้ำ/เลียน LeetCode, reference approach แก้ได้จริง, มี examples ครบ
- **คุณรีวิว**: ลองทำ 1 ข้อจริง ๆ ก่อนอนุมัติ format
- Status: `todo`

## C4 — Content: AWS Domain 1 batch แรก (5 lessons) `[lesson-writer → lesson-verifier]`
- iam-users-roles-policies, iam-policy-evaluation, organizations-scp, cognito, kms-encryption
- **ก่อนเริ่ม**: ต้องเตรียม excerpt ลง `content/aws-docs/` ก่อน (orchestrator รวบรวมจาก official docs) — verifier จะ FAIL ทุกข้อเท็จจริงที่ไม่มีใน excerpt
- References ของ AWS lessons ชี้ไปหน้า official docs ที่ใช้ ground เสมอ
- Status: `todo`
