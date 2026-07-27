# AI & LLM Systems track (7th track) — tickets

Track ที่ 7: **AI & LLM Systems** (track slug `ai-systems`, topic slug
`ai-and-llm-systems`, position 7). เป้าหมาย: เรียน AI/LLM engineering **จากศูนย์ →
ลึกระดับสัมภาษณ์** เจาะตำแหน่ง Backend/AI-CRM (LINE MAN Wongnai). radar กลายเป็น 7 แกน.
Slug = contract แช่แข็ง (เหมือน DDD/DDIA).

## T-ai-track — enable track + curriculum tree + loader test (PR open)

- domain `Track` array + `Tracks()` ได้ `{value:"ai-systems"}` (ต่อท้ายสุด) + `track_test.go`
- `migrations/004_ai-systems-track.sql` ขยาย ENUM `topics.track`
- `content/curriculum/ai-systems.json` — 11 บท / 53 concept
- `internal/curriculum/infra/contentfile_ai_test.go` — loader test ล็อก slug
- go vet + test เขียว, code-reviewer APPROVE (migration order ตรง enum, AI claims ถูกต้อง)

## T-ai-lessons-* — lesson batches (batch = 1 บท)

lesson-writer → lesson-verifier ทีละ concept. ไฟล์ลง
`content/lessons/ai-and-llm-systems/<concept>.json`. ทุก lesson: Thai + mermaid +
section trade-off/"when not to use" + recall (mcq **เฉลยห้ามเป็น option ยาวสุด**,
สลับตำแหน่ง incl. middle) + refs. เน้น**ความแม่นยำ** (ใช้ claude-api skill ground
เมื่อแตะ Claude/Anthropic specifics; log-prob เป็นตัวอย่างทั่วไป — Anthropic ไม่ expose).

**ลำดับบท (ตาม tree ใน `content/curriculum/ai-systems.json`):**
1. llm-foundations · 2. prompting-and-context · 3. tool-and-function-calling ·
4. retrieval-augmented-generation · 5. streaming-and-realtime (4) ·
6. ai-agents-and-orchestration · 7. conversational-and-agent-assist ·
8. guardrails-and-safety · 9. llm-evaluation · 10. production-llm-systems ·
11. crm-and-support-ai (4)

**เรียน→test→วัดผล:** อ่าน lesson → recall check → SM-2 SRS → radar แกน AI เต็ม =
พร้อมตอบลึก. Roadmap ภาพใหญ่: [🎯 artifact](https://claude.ai/code/artifact/a6d332fd-de7b-44bb-91db-37549984c7a7)
