# Design — Self-Improve Web (Phase 0, approved 2026-07-21)

เอกสารนี้คือ design ฉบับเต็มที่คุยกันไว้ — ภาพใหญ่อ่านที่ [`roadmap.md`](roadmap.md),
ticket รายชิ้นอยู่ [`tickets/`](tickets/)

## Decision log

| วันที่ | ตัดสินใจ |
|---|---|
| 2026-07-21 | Auth = static bearer token (เก็บใน `.env`, frontend จำใน localStorage) |
| 2026-07-21 | Deploy = DigitalOcean VPS Singapore ~$6/เดือน — **ไม่ host บน AWS** (AWS เป็นเนื้อหาเรียนเท่านั้น + free-tier sandbox แยก) |
| 2026-07-21 | DSA problems = generate offline เป็น bank ต่อ pattern, runtime ใช้แค่รีวิวคำตอบ |
| 2026-07-21 | LLM client เขียนเองด้วย net/http — ไม่ใช้ SDK / agent framework (พิจารณา Hermes Agent แล้ว ตัดออก) |
| 2026-07-21 | คุยกันภาษาไทย, docs ภาษาไทย, UI copy อังกฤษ, lesson ไทย |
| 2026-07-22 | Lesson มีภาพประกอบเป็น mermaid diagram ใน body_md (เรนเดอร์ฝั่งเว็บ) — animation เต็มรูปแบบเป็น v2 |
| 2026-07-22 | Lesson มี `references` 2–4 แหล่ง (primary sources จริง ห้ามแต่ง URL); เพิ่มบทเรียนภายหลังได้เสมอ (pipeline เป็น additive) |
| 2026-07-22 | ขึ้น VPS ตั้งแต่สัปดาห์ 2 + ทุก ticket เป็น PR เข้า develop, รีวิวผ่าน GitHub mobile ได้, merge = auto-deploy |
| 2026-07-22 | Vision: v1 เป็นหนังสือเรียนส่วนตัวใช้คนเดียว → Phase 2 publish เป็น portfolio (ไม่หาเงิน) — ดู §8 |

## 1. DDD Tactical Design

Modular monolith, Go binary เดียว, 6 bounded contexts — contexts คุยกันผ่าน
**in-process domain events** และ shared ID value objects เท่านั้น (ห้าม import
domain type ข้าม context)

### Curriculum (reference data — write path มีแค่ importers)
- Aggregates: `Topic` (root → `Chapter` → `Concept`), `Lesson` (root → `RecallCheck`)
  — 1 lesson ต่อ 1 concept = 1 "chunk" อ่าน 5–10 นาที
- VOs: `Slug`, `Track` (ddd|distsys|aws|go|dsa), `Position`, `EstMinutes`
- Events: `LessonImported`

### Learning (session runner + การอ่าน)
- Aggregates: `Session` (state machine: Planned → Active → Completed/Abandoned,
  invariant: active ได้ทีละ 1), `LessonProgress` (invariant: lesson ถัดไปปลดล็อก
  เมื่อ lesson ก่อนหน้าในบทผ่าน recall ครบ — คือกติกา gating), `RecallAttempt`
- VOs: `SessionType` (read_recall|drill|build|mock), `TimerSpec` (30/30/45/45),
  `Grade` (0–5), `ChunkState` (locked|in_progress|passed)
- Events: `SessionStarted`, `SessionCompleted`, `RecallAnswered{grade}`,
  `RecallFailed`, `LessonCompleted`

### Practice (DSA)
- Aggregates: `Problem` (import จาก offline bank), `Attempt` (approach + code + `ReviewResult`)
- VOs: `PatternRef`, `Difficulty`, `ReviewResult` (verdict, time/space complexity, idiomatic notes, score 0–5)
- Events: `AttemptReviewed{pattern, score}`

### Review (SRS)
- Aggregate: `ReviewCard` — `Apply(grade)` รัน SM-2 คืน state ใหม่ (pure function,
  เป้าหมายหลักของ table-driven tests)
- VOs: `SM2State` (easiness factor, interval, repetitions, due date), `CardSource`
- Events: `CardScheduled`, `CardReviewed`

### Assessment (pre-test / retest)
- Aggregates: `QuestionBank` (ต่อ topic, **immutable หลัง import** — retest ใช้ชุดเดิม
  เพื่อเทียบ delta ได้), `TestAttempt`
- Events: `TestCompleted{topic, score}`

### Progress (stats — read models ล้วน ไม่มี aggregate)
- Projections จาก events: `Measurement` (1 จุดกราฟต่อ 1 การวัด), `DailyActivity`
  (streaks), `PatternStats`, `SkillRadar`; weekly retro = query จาก event log

**Event mechanism**: dispatcher แบบ synchronous in-process ใน application layer +
append ลงตาราง `events` (outbox-lite) → projections rebuild ได้ และเป็นแหล่งข้อมูล retro

## 2. MySQL Schema

```
-- curriculum
topics          (id PK, track ENUM('ddd','distsys','aws','go','dsa'), slug UNIQ, title, position)
chapters        (id PK, topic_id FK, slug, title, position, UNIQUE(topic_id, slug))
concepts        (id PK, chapter_id FK, slug, title, outline TEXT, position, UNIQUE(chapter_id, slug))
lessons         (id PK, concept_id FK UNIQ, version INT, title_en, est_minutes, body_md MEDIUMTEXT,
                 refs JSON, imported_at)   -- body_md มี ```mermaid``` fences ได้; refs = "references" JSON (คำสงวนใน MySQL)
recall_checks   (id PK, lesson_id FK, position, type ENUM('short_answer','mcq'),
                 question TEXT, expected_answer TEXT, options JSON NULL)

-- learning
sessions        (id PK, type ENUM('read_recall','drill','build','mock'), planned_minutes,
                 status ENUM('active','completed','abandoned'), started_at, ended_at NULL, summary JSON NULL)
lesson_progress (id PK, lesson_id FK UNIQ, state ENUM('locked','in_progress','passed'),
                 first_passed_at NULL, last_read_at NULL)
recall_attempts (id PK, recall_check_id FK, session_id FK NULL, answer TEXT, grade TINYINT,
                 passed BOOL, feedback TEXT, graded_by ENUM('llm','self'), created_at)

-- practice (DSA)
dsa_problems    (id PK, concept_id FK, slug UNIQ, difficulty ENUM('easy','medium','hard'),
                 prompt_md MEDIUMTEXT, examples JSON, hints JSON NULL, reference_approach_md TEXT, imported_at)
dsa_attempts    (id PK, problem_id FK, session_id FK NULL, approach_md TEXT, code TEXT,
                 verdict ENUM('correct','partially_correct','incorrect'),
                 complexity_time, complexity_space, score TINYINT, review_md TEXT, created_at)

-- review (SRS)
review_cards    (id PK, source_type ENUM('recall_check','dsa_pattern'), source_id BIGINT,
                 concept_id FK, ef DECIMAL(3,2) DEFAULT 2.50, interval_days, repetitions,
                 due_on DATE, suspended BOOL, UNIQUE(source_type, source_id))
review_logs     (id PK, card_id FK, grade TINYINT, prev_interval, next_interval, ef_after, reviewed_at)

-- assessment
diag_questions  (id PK, topic_id FK, position, question_md, options JSON, correct_index TINYINT, explanation_md)
test_attempts   (id PK, topic_id FK, kind ENUM('baseline','retest'), score, total, answers JSON, taken_at)

-- progress + platform
events          (id PK, name, payload JSON, occurred_at)
measurements    (id PK, metric VARCHAR(50), topic_id FK NULL, value DECIMAL(8,2), taken_at)
daily_activity  (day DATE PK, sessions INT, minutes INT, recalls INT, drills INT)
tickets         (id PK, week TINYINT, position, title, description, agent, est_minutes,
                 status ENUM('todo','doing','done','skipped'), done_at NULL)
schema_migrations (version PK, applied_at)
```

## 3. API Endpoints

ทั้งหมดอยู่ใต้ `/api/v1` หลัง bearer-token middleware ยกเว้น `/healthz`

| ส่วน | Endpoints |
|---|---|
| Curriculum | `GET /curriculum` (tree + progress), `GET /lessons/{id}` (ตัด expected answers ออก) |
| Sessions | `POST /sessions`, `GET /sessions/active`, `POST /sessions/{id}/complete`, `POST /sessions/{id}/abandon` |
| Read & Recall | `GET /read/next`, `POST /recall-checks/{id}/attempts` → Haiku ตรวจ → `{grade, passed, feedback, next_unlocked}` |
| Drill | `GET /drill/queue?limit=`, `POST /review-cards/{id}/review {grade}` |
| DSA | `GET /dsa/patterns`, `GET /dsa/patterns/{slug}/next-problem`, `GET /dsa/problems/{id}`, `POST /dsa/problems/{id}/attempts` |
| Pre-test | `GET /tests/{topicSlug}`, `POST /tests/{topicSlug}/attempts` |
| Stats | `GET /stats/dashboard`, `GET /stats/trends?metric=&topic=`, `GET /stats/retro?week=` |
| Build | `GET /tickets/next`, `GET /tickets?week=`, `PATCH /tickets/{id}` |
| Platform | `GET /healthz` (public) |

## 4. Repo Structure

```
self-learning/
├── CLAUDE.md · go.mod · Makefile · Dockerfile · docker-compose.yml
├── cmd/{api, import-curriculum, import-lessons, import-dsa, import-tests}/
├── internal/
│   ├── curriculum/ learning/ practice/ review/ assessment/   # แต่ละอัน: domain/ app/ infra/
│   ├── stats/                          # app/ + infra/ (projections)
│   └── platform/{httpserver, middleware, mysql, llm, events, config}/
├── migrations/*.sql
├── web/                                # Next.js static export + Tailwind
│   ├── app/{dashboard,read,drill,dsa,test,tickets}/
│   ├── components/charts/
│   └── lib/api.ts
├── content/
│   ├── curriculum/{ddd,distsys,aws,go,dsa}.json
│   ├── lessons/<topic>/<concept>.json
│   ├── dsa-problems/<pattern>/<slug>.json
│   ├── tests/<topic>.json
│   └── aws-docs/*.md                   # excerpt จาก official docs ไว้ ground บทเรียน AWS
├── deploy/{docker-compose.prod.yml, Caddyfile, setup-vps.sh, backup.sh}
├── docs/{design.md, roadmap.md, tickets/}
└── .claude/{agents/, skills/}
```

กติกา layer: domain import ได้แค่ stdlib + package ตัวเอง — MySQL/Anthropic client
อยู่ `platform/` inject ผ่าน interface ที่ application layer ประกาศ (ports & adapters)

## 5. Deployment — DigitalOcean VPS (Singapore)

```
[Browser] ──HTTPS──> [Caddy]  ── static ──> Next.js export (ไฟล์ในเครื่อง)
                        │
                        └── reverse proxy /api ──> [Go API :8080] ──> [MySQL :3306]
                                                        │ (Docker network เดียวกัน,
                                                        │  MySQL ไม่ expose ออกนอกเครื่อง)
                                                        └──HTTPS──> Anthropic API (Haiku)
```

- VPS $6/เดือน (1GB RAM) + domain ~$10/ปี → รวม **~$7/เดือน** (เทียบ AWS ~$27)
- Secrets ใน `.env` (chmod 600): `ANTHROPIC_API_KEY`, `DB_PASSWORD`, `API_BEARER_TOKEN`
- Security: SSH key only, UFW เปิด 22/80/443, fail2ban
- Backup: cron `mysqldump` ทุกคืน → rclone ไป Cloudflare R2 / Backblaze B2 (free tier)
- Deploy: merge PR เข้า `develop` → GitHub Actions → build image → push GHCR →
  SSH `docker compose up -d` → smoke test `/healthz` (ขึ้น VPS ตั้งแต่สัปดาห์ 2 —
  ทุก ticket ที่ merge เห็นบนเว็บจริงทันที ใช้/รีวิวจากมือถือได้)
- เหตุผลที่ไม่ใช้ AWS: แพงกว่า ~4 เท่าสำหรับ single user; SAA-C03 เรียนผ่าน
  curriculum + hands-on lab ใน free-tier sandbox แยกซึ่งสร้าง/ลบได้อิสระ

## 6. Content Pipeline

Format tree: `content/curriculum/<track>.json`
```json
{ "track": "ddd", "topics": [{ "slug": "…", "title": "…", "chapters": [
  { "slug": "…", "title": "…", "concepts": [
    { "slug": "…", "title": "…", "outline": "bullet 3–5 ข้อที่ lesson ต้องครอบคลุม" } ] } ] }] }
```

ต่อ 1 batch (= **1 chapter, 3–8 concepts** ต่อ 1 Claude Code session เพื่อประหยัด quota):
1. orchestrator เลือก chapter ถัดไป → spawn **lesson-writer (opus)** ทีละ concept
   (AWS concepts แนบ excerpt จาก `content/aws-docs/`)
2. ทุก draft ผ่าน **lesson-verifier (sonnet)** → PASS/FAIL + issues
3. FAIL → regenerate พร้อมแนบ issues (สูงสุด 2 รอบ, ไม่ผ่าน = พัก + รายงาน)
4. `go run ./cmd/import-lessons -dir content/lessons/<topic>` (idempotent upsert)
5. รายงาน batch สั้น ๆ: เขียน / ตก / import แล้ว

Pipeline เดียวกันใช้กับ **DSA bank** (batch = 1 pattern × 5 ข้อ; verifier เช็คเพิ่ม:
ห้ามเหมือนโจทย์ LeetCode + reference approach แก้โจทย์ได้จริง) และ **pre-test bank**
(batch = 1 topic × 15–20 MCQ)

Lesson ทุกบท: มี mermaid diagram เมื่อ concept มี flow/สถาปัตยกรรม/state machine
และมี `references` 2–4 แหล่ง (หนังสือระบุ chapter, official docs) — verifier
เช็คทั้งคู่ (URL ที่ดูแต่งขึ้น = FAIL)

Cadence: just-in-time — generate ล่วงหน้า 1 chapter ก่อนถึงคิวอ่านเสมอ ไม่ generate ทิ้งไว้ทั้ง tree

**เพิ่มบทเรียนในอนาคต** (ทำได้เสมอ ไม่ต้องแก้โค้ด): เพิ่ม concept ลง
`content/curriculum/<track>.json` → `import-curriculum` → รัน batch writer/verifier
→ `import-lessons` — ทุกขั้น idempotent/additive จะเพิ่ม track ใหม่ทั้ง track ก็ได้

## 7. Curriculum Tree v1 (~175 concepts)

### (a) Domain-Driven Design — 5 บท, 27 concepts
1. **Model-Driven Foundations**: ubiquitous-language, model-driven-design, knowledge-crunching, hands-on-modelers
2. **Building Blocks**: layered-architecture, entities, value-objects, domain-services, modules, aggregates, aggregate-design-rules, factories, repositories, domain-events
3. **Supple Design & Refactoring**: intention-revealing-interfaces, side-effect-free-functions, assertions, specification-pattern, making-implicit-concepts-explicit, refactoring-toward-deeper-insight
4. **Strategic Design**: bounded-context, context-mapping, shared-kernel, customer-supplier-conformist, anticorruption-layer, open-host-service-published-language, core-domain-distillation, generic-subdomains, large-scale-structure
5. **DDD in Go**: ddd-go-project-layout, persistence-without-orm, in-process-domain-events, testing-the-domain-layer

### (b) Distributed Systems — 8 บท, 37 concepts
1. **Foundations**: why-distributed, transparency-goals, scalability-dimensions, fallacies-of-distributed-computing
2. **Architectures**: client-server, multi-tier-layered, peer-to-peer, microservices-vs-monolith, event-driven-architecture
3. **Communication**: rpc-fundamentals, message-queues, publish-subscribe, rest-vs-grpc, serialization-formats
4. **Naming & Discovery**: flat-naming, structured-naming-dns, service-discovery
5. **Coordination & Time**: physical-clocks-ntp, lamport-clocks, vector-clocks, distributed-mutex, leader-election, gossip-protocols
6. **Consistency & Replication**: replication-motivation, linearizability-sequential, causal-consistency, eventual-consistency, client-centric-models, quorum-protocols, crdt-intro
7. **Fault Tolerance**: failure-models, failure-detection, process-resilience, consensus-raft, two-phase-commit, sagas-compensation, recovery-checkpointing
8. **Interview Patterns**: cap-pacelc, partitioning-sharding, idempotency-retries, backpressure-rate-limiting, outbox-pattern, distributed-caching

### (c) AWS SAA-C03 — 4 domains, 37 concepts
1. **Secure Architectures (30%)**: iam-users-roles-policies, iam-policy-evaluation, organizations-scp, cognito, kms-encryption, secrets-vs-parameter-store, sg-vs-nacl, vpc-endpoints-privatelink, waf-shield, s3-security
2. **Resilient Architectures (26%)**: regions-az-edge, elb-types, auto-scaling-groups, rds-multi-az-read-replicas, aurora-ha, sqs-sns-decoupling, eventbridge, route53-routing-policies, dr-strategies, backup-strategies
3. **High-Performing (24%)**: ec2-families-purchasing, ebs-efs-instance-store, s3-performance, cloudfront, elasticache, dynamodb-fundamentals, dynamodb-advanced, rds-performance, kinesis, athena-glue, lambda-performance, ecs-eks-fargate
4. **Cost-Optimized (20%)**: pricing-models-ri-sp-spot, s3-storage-classes-lifecycle, compute-cost-optimization, data-transfer-costs, cost-tools-budgets

### (d) Deep Go — 6 บท, 33 concepts
1. **Runtime & Scheduler**: gmp-model, goroutine-stacks-growth, preemption, netpoller, sysmon
2. **Memory & GC**: escape-analysis, allocator-layout, gc-tricolor-write-barriers, gc-pacing-gogc, stack-vs-heap-performance
3. **Concurrency Internals**: channel-internals-hchan, select-implementation, mutex-internals, waitgroup-once-cond, atomics-memory-model, context-propagation, common-concurrency-bugs
4. **Types & Generics**: interface-internals-itab, nil-interface-pitfalls, embedding-composition, generics-implementation, reflection-cost
5. **Performance & Tooling**: pprof-cpu-heap, execution-tracer, benchmarking-methodology, race-detector, compiler-optimizations-pgo
6. **Stdlib Internals**: net-http-server, http-client-transport, database-sql-pooling, encoding-json, errors-wrapping, slices-maps-internals

### (e) DSA — NeetCode 150 patterns — 18 patterns, ~41 concepts
| Pattern | Concepts |
|---|---|
| arrays-hashing | hash-frequency, prefix-sums, two-sum-family |
| two-pointers | converging, in-place-partition |
| sliding-window | fixed, variable-shrink |
| stack | monotonic-stack, matching-pairs |
| binary-search | on-index, on-answer-space, rotated-arrays |
| linked-list | reversal, fast-slow-cycle, merge |
| trees | dfs-patterns, bfs-level-order, bst-properties, lca, serialization |
| tries | build-search, word-search-with-trie |
| heap | top-k, two-heaps-median, k-way-merge |
| backtracking | subsets-permutations, constraint-pruning |
| graphs | representation-traversal, islands-components, topological-sort, union-find |
| advanced-graphs | dijkstra, mst, bellman-ford |
| dp-1d | memo-vs-tabulation, house-robber-family, coin-change, lis |
| dp-2d | grid-paths, lcs-edit-distance, knapsack-01 |
| greedy | exchange-argument, jump-gas |
| intervals | sort-merge, sweep-line-rooms |
| math-geometry | matrix-ops, number-theory-basics |
| bit-manipulation | bit-tricks, xor-patterns |

## 8. Phase 2 — Publish เป็น portfolio (หลัง v1 พิสูจน์ตัวเองแล้ว)

เงื่อนไขเริ่ม: ใช้ v1 เรียนเองจนรู้สึกว่า "ได้จริง" ไม่ใช่ตามปฏิทิน

ขอบเขตคร่าว ๆ (ยังไม่ design ละเอียด — ไว้ทำตอนถึงเวลา):
- Public read-only mode: บทเรียน + curriculum เปิดอ่านได้โดยไม่ต้องมี token,
  ข้อมูลส่วนตัว (stats, attempts, streak) ยังอยู่หลัง auth เสมอ
- Open-source ตัว repo เป็นผลงานหลัก (โค้ด DDD + hand-written middleware +
  LLM client คือจุดขายต่อ recruiter มากกว่าตัวเนื้อหา)
- ก่อน publish เนื้อหา: ตรวจ licensing อีกรอบ — บทเรียนสรุปจากความรู้ทั่วไป
  ไม่มีข้อความจากหนังสือ (กติกาเดิม) แต่การเผยแพร่สาธารณะต้อง audit ซ้ำ
  + ใส่ attribution ชัดเจนว่าโครงอิงหนังสือ/docs อะไร
- ไม่หาเงิน ไม่มี multi-user, ไม่มี billing — ตัด scope พวกนี้ทิ้งได้เลย

ผลต่อ v1 ตอนนี้: แทบไม่มี — แค่ (1) อย่า hardcode อะไรที่ผูกกับความเป็น
single-user ลึกเกินถอน (auth middleware แยกชั้นอยู่แล้ว), (2) เก็บกติกา
"ห้าม copy ข้อความหนังสือ" เข้มตั้งแต่วันแรก
