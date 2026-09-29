# T-hbw8 — the compared set and the selection

**The selection for the withdrawn first draft.** This record selected for the
first draft (ADR-0013 to ADR-0020 and `docs/architecture.md`), which plan v2
removed in step 0. It predates O-76 to O-92 and the evaluation. The selection
for the rewrite is `selection-v2.md`.

The author (Claude Opus 5.5) selects from the options of the three panel members: `panel-A.md` (orchestration, Git state, the forge; GPT-6 Sol), `panel-B.md` (decisions, routing, learning; Grok 4.7) and `panel-C.md` (governance, invariants, human decision points; Claude Opus 5.5). An option is named by member and number, for example A-1A. The method is Solution selection (`docs/engineering-discipline.md`): first the search for a public solution, then the applicable considerations, which are not a scoring formula. The invariants `F-0001#1`–`#9` remove a candidate; they do not rank one. A choice that no fact and no O-decision settles is a question for the Operator (section 11); the selection names the option it depends on and does not decide it.

Marks: ✓ serves, ✗ strains or fails, — neutral.

## 1. The search for a public solution (Solution selection, first step)

Searched on 2026-09-27 (web search; the sources are listed below). The question: does a public solution already do what the orchestrator, the squad manager or the learning loop must do, so that LAYUP should use it and not build its own?

| Need | Public solutions found | What they do | Why LAYUP does not adopt one as the component |
| ---- | ---------------------- | ------------ | --------------------------------------------- |
| Orchestrator of coding-agent CLIs | Bernstein (Python, Apache-2.0; plan → execute → verify → merge; 49 harness adapters; state in a `.sdd/` directory); Agent Orchestrator (TypeScript, Apache-2.0; 26 harnesses; desktop app and daemon); Claude Squad (Rust, AGPL-3.0; worktrees and tmux); Microsoft Conductor (Python; YAML workflows); crew, ccswarm, Archon, NXTG-Forge, loki-mode, Claudexor (cross-family review, quota-aware rotation) | Start several harness CLIs, each in its own git worktree, and collect their pull requests; some add test and lint gates before a merge; one (Claudexor) routes across harnesses and reviews across model families | None keeps its state only in the project repository in Git (Invariant 1); none has the Armature gate, the PSB Human Decision Points or the escalation rule; none runs the lifecycle from a problem statement to acceptance and retrospective; most are Python, TypeScript or Rust, where ADR-0010 sets Go and O-10/O-11 keep LAYUP's code out of a target. **Kept as patterns:** a worktree per role session; a harness adapter per CLI (a command line, not an SDK); a verify step before the merge. |
| Durable workflow engine for the phase state machine | Temporal; Microsoft `durabletask-go` (Go, embeddable, SQLite storage); DBOS (Postgres) | Durable execution of a workflow that survives a crash, with its history in a database | Each keeps the state in its own store (a server or a database), not in Git, so a second harness on a clean clone cannot read it (Invariants 1 and 9). **Kept as a pattern:** the state machine replays from a recorded event history; in LAYUP that history is the TSV records in Git. |
| Learning router for models and harnesses | RouteLLM (preference-data routers between a strong and a weak model); LLMRouter (a library of routers); MetaLLM and OrcaRouter (multi-armed bandits, LinUCB); ChainRoute's BanditRouter (epsilon-greedy over success and failure) | Choose a model per request from observed quality and cost | They route single model calls at request time on large volumes; LAYUP routes a few role sessions per task, where the reward (a review verdict, an acceptance, a cost) comes days later. **Kept as a pattern:** a bandit-style update of a routing table from recorded outcomes, with no change to model weights (O-67, O-69). |

Sources:
- https://www.augmentcode.com/tools/open-source-agent-orchestrators (fetched 2026-09-27)
- https://github.com/andyrewlee/awesome-agent-orchestrators (fetched 2026-09-27)
- https://github.com/sipyourdrink-ltd/bernstein ; https://github.com/nwiizo/ccswarm ; https://github.com/razzant/claudexor
- https://github.com/microsoft/durabletask-go ; https://temporal.io
- https://github.com/ulab-uiuc/LLMRouter ; https://arxiv.org/html/2605.30736v1 (OrcaRouter) ; https://github.com/aaravbaphna/chainroute/pull/4

## 2. Points where the members say the target design is at risk

| Point | Members | Result |
| ----- | ------- | ------ |
| Jev cannot write a panel synthesis: System One models "do not write replies, produce code, or generate explanations of their reasoning" (`jev-sources.md`) | C (P1), B (2A), A (2A) | **Accepted.** A reasoning-model session writes a synthesis; Jev may only select or score among written options. The decision point is renamed "panel disposition". |
| One Noul "for each of the four PSB axes" hides a conjunction and arithmetic | B, A (2B), C (P4) | **Accepted.** A deterministic floor first; then four literal Nouls, one condition each, combined in code with OR; code does each sum, date and count. |
| A human comment cannot be told from an agent comment: agents and the Operator push with one GitHub account (O-9, ADR-0011) | C (P2), A (5A, question 1), B (5A) | **A question for the Operator** (section 11, Q-2). No option below proves a human decision while this holds. |
| A Jev threshold with no data breaks `F-0001#4` | C (P3), B, A | **Accepted.** Each threshold starts as a cited conservative start value and is calibrated from pilot records; how much authority Jev gets before that is Q-3. |
| The retrospective approval is planned only if it is listed before delivery (`F-0001#25`) | C (P5) | **Accepted.** Intake writes each milestone's retrospective approval, with its approver, into the list of planned approval points. |
| Stack gates in the target change O-10 and O-11 as written | C (P6), plan review condition 1 | **A question for the Operator** (Q-1). |

## 3. Question 1 — the orchestrator (B1, B2)

| Option | Inv-1 state in Git | Inv-2 target alone | Stall seen by the clock (B2) | Crash of the orchestrator seen | Infrastructure | Selected |
| ------ | ------------------ | ------------------ | ---------------------------- | ------------------------------ | -------------- | -------- |
| One foreground `layup run`, events pushed to Git before the next role, attempt IDs, reconcile on restart (A-1A, B-1A, C-Q1-A) | ✓ | ✓ (the target does not need it to pass its gates) | ✓ a clock in Go against the last durable progress event | ✗ alone: a dead process sees nothing | one host during delivery | **yes, with the next row** |
| A heartbeat row pushed every H, and a scheduled dead-man job that opens a stall record when the heartbeat is older than 2H (C-Q1-C) | ✓ | ✓ | — | ✓ the stall reaches the Operator as a stall package, not by chance (`F-0003#17`) | a small scheduled job on the runner | **yes, added to the row above** |
| Short `--once` runs from a Git work queue or forge events (A-1B, B-1B, C-Q1-B) | ✓ | ✓ | ✗ the clock is off between events unless the runner polls (B) | ✓ | harness credentials and sessions on the forge runner; lock-in grows from checks to the whole loop | no |

Kept from the rejected rows: the attempt ID and the "one active attempt per task" lease (A-1B). A progress event is a new durable record: a changed artifact hash, a successful gate run, a handoff row; equal failure signatures and no new pushed head are not progress (A-1A). A restart that a human causes is recorded as unplanned input (`F-0001#28`). The start values T and N are Q-7.

## 4. Question 2 — the decision component (A3, B8)

| Option | Inv-6 check first | O-71 fallback | Generation needed | Human load | Selected |
| ------ | ----------------- | ------------- | ----------------- | ---------- | -------- |
| A deterministic floor, then fixed literal question packs per named point, one request each, code composes; `decisions.tsv` per call; pinned `jev-1.13.0` (B-2A, A-2B, C shared text) | ✓ | ✓ error, 429 after backoff, context overflow or below threshold → human | none | low to medium | **yes** |
| One typed request per event, all questions batched (A-2A) | ✓ | ✓ | none | low | kept as the request shape inside the row above |
| A generative model proposes, Jev verifies (B-2B) | ✓ | ✓ | a model call outside the named points | medium | no: it adds a model call that O-72's successor does not name |
| No Jev on the live path (B-2C) | ✓ | ✓ | none | high | no: rejects O-67 |

The named points: escalation selection, question routing by ambiguity kind, stall action, panel disposition, model and harness fit, over budget inside the band. Jev never sees a sum, a date comparison or a count (`jev-sources.md`, jaggedness page). Jev's authority before the thresholds have evidence (direct, shadow for milestone 1, or cautious direction only: C-Q2-A/B/C) is Q-3. **ADR-0011 decision 8:** the engine checks (`layup psb check`, `layup setup verify`, `layup gate`, `layup spec check`) make no model call, and their interface is plain text and files; only the decision component of `layup run` calls Jev, and a target never needs Jev to pass its gates (C, B).

## 5. Question 3 — squads and routing (B5)

| Option | Inv-9 counterpart harness | Vision 2.2 (allocate by complexity) | Structural rule in code | Selected |
| ------ | ------------------------- | ----------------------------------- | ----------------------- | -------- |
| `harnesses.tsv` and `routing.tsv`; code admits a pair (registered, probed, role match, verifier harness ≠ author harness); one Jev Score for fit among the admitted pairs; weight as the tie break (B-3A, A-3B, C shared tables) | ✓ enforced by code | ✓ | ✓ | **yes** |
| A fixed pair per role, no Jev (B-3B) | ✓ | ✗ | ✓ | no: drops the vision item O-67 keeps |
| One Jev Choice over the whole register (B-3C) | ✗ a judgement decides Invariant 9 | ✓ | ✗ | no |
| One neutral rule source; each harness entry file is only a pointer, checked (C-Q3-C, A-3B) | ✓ | — | ✓ | **yes, added** |

A harness enters `harnesses.tsv` only after a recorded probe (A-3A). With fewer than two working harnesses, a verification is `not-active` and the change does not merge: Invariant 5 settles this (`F-0001#5`), so member B's question 3 is not asked. The list of roles and the owner of each ambiguity kind (PRD §11 question 3; A-3A, A-3B, B-3A, C-Q3-A, C-Q3-B) is Q-6.

## 6. Question 4 — the records in Git (B3, B6, B7)

| Option | Inv-1 | Forged record by an agent | Clean-clone reading (Inv-9) | Selected |
| ------ | ----- | ------------------------- | --------------------------- | -------- |
| One protected records branch of the target, append-only TSV tables and typed payloads, written only by the orchestrator's identity (A-4A) | ✓ | ✓ agents cannot push it (needs Q-2) | ✓ fetch one named branch | **yes** |
| A protected ref per task with an index (A-4B) | ✓ | ✓ | ✗ a default clone misses task refs unless the index is read | no |
| Records on the default branch through agent commits (C tables, B-4A) | ✓ | ✗ an agent can write a row | ✓ | tables kept, location no |
| A forge comment is copied into Git (body, author, id, SHA-256) before any action reads it; an edit is a new correction event (C-Q4-A, B-4B, A-4A) | ✓ | — | ✓ | **yes** |
| A human decision as a commit in a PR (C-Q4-B) | ✓ | ✓ | ✓ | no: human load; a PR approval is unplanned input by the PSB |
| A workflow in the target copies comments (C-Q4-C) | ✓ | — | ✓ | no: it writes a LAYUP workflow into the target (O-10) |

The tables (union of A, B and C): `handoffs.tsv`, `questions.tsv`, `acceptance.tsv` (Decision Point 3), `telemetry.tsv` per action with requirement ID and price source, `decisions.tsv`, `approvers.tsv` written at Intake, `prices.tsv` with a cited source. `not reported` is incomplete, never zero (`F-0003#60`).

## 7. Question 5 — gates and rule protection (A1, A2, B9)

| Option | Inv-2 target alone | O-10/O-11 as written | ADR-0012 part 6 | Selected |
| ------ | ------------------ | -------------------- | --------------- | -------- |
| Stack gates in the target as kit content: a gate recipe per stack in LAYUP; Scaffold writes `stack.tsv`, the tool configuration and a CI job into the target; public stack tools first, a target-local check script only where no tool fits, `not-active` otherwise; `layup gate` runs the same commands from outside (A-5A, A-5B, B-5A, C-Q5-A, C-Q5-B) | ✓ | ✗ needs Q-1 | ✓ | **yes, if Q-1 allows** |
| Gates outside, exported at each retrospective (C-Q5-C) | ✗ between exports | ✓ | ✓ | the fallback if Q-1 keeps O-10/O-11; then A2 is a known limit |
| LAYUP holds the gates, the runner is the only enforcer (B-5B) | ✗ | ✓ | ✓ | no: this is finding A2 |

Rule paths for all: a proposed rule or gate change is collected during the milestone and lands in one retrospective batch, approved at that planned point by the approver of `approvers.tsv` (O-69; C, A, B). A rule-path commit at any other time fails a check and counts as unplanned input (`F-0001#28`); it is never reclassified as a planned approval (A1). Who writes a gate recipe for a new stack (B9): a gate owner through a LAYUP change with a known-bad fixture per gate kind and evidence for two stacks (`F-0003#67`), not a governed coding agent (A-5A).

## 8. Question 6 — escalation and budget (B7, B8)

| Option | Before the work (DP 4) | Arithmetic in code | Cost stop before the spend | Selected |
| ------ | ---------------------- | ------------------ | -------------------------- | -------- |
| The floor, then four literal Nouls, OR in code, the middle band escalates, before a session starts (B-2A/6A, C-Q6-A) | ✓ for planned decisions | ✓ | — | **yes** |
| The same screen also on each handoff's `decisions` field (C-Q6-B) | ✓ at each role boundary | ✓ | — | Q-4 |
| Band in the budget file (B and U per milestone, written at Intake); code sums telemetry and projects; below B continue, between B and U one Jev Noul on the semantic question, at U or projected above U a hard stop with no Jev call (C-Q6-C, B-6A) | — | ✓ | ✓ between actions | **yes** |
| Prepaid bounded actions with an enforced maximum (A-6B), or a live metered session (A-6A) | — | ✓ | ✓ inside an action | the per-action cap depends on what each harness can enforce: Q-5 |
| The harness estimates a price and Jev judges it (B-6B) | — | ✗ | ✗ | no |

## 9. Question 7 — the learning loop (table C 2.3, 3.3)

| Option | Audit | `F-0001#4` | Vision 3.3 | Selected |
| ------ | ----- | ---------- | ---------- | -------- |
| A reward formula in Go over the Git tables at the retrospective (accepted at first review +, reversal −, stall needing the Operator −, unplanned input −; incomplete telemetry 0; not the first-attempt gate pass rate, which `F-0003#58` keeps as monitoring only); a bounded step, a minimum evidence count; applied only in the approved retrospective batch (B-7A, C-Q7-A, A-7A) | ✓ | ✓ with a cited start formula | ✓ | **yes** |
| Lessons as rows with their evidence; an optional Jev Score ranks them; the reward still decides (C-Q7-C, B-7A) | ✓ | ✓ | ✓ | **yes, added** |
| Bounded exploration trials (A-7B, C-Q7-B) | ✓ replayable | — | ✓ | no for the pilot: it spends budget on weaker routes before a baseline exists |
| Jev proposes weights (B-7B) | ✗ | ✗ | ✓ | no |
| Weights frozen in the pilot (B-7C) | ✓ | ✓ | ✗ | Q-8 |

## 10. Question 8 — specification synthesis (B4, table C 2.1)

| Option | Deterministic trace (`F-0003#62`) | Drift of meaning | Selected |
| ------ | --------------------------------- | ---------------- | -------- |
| Intake numbers the PSB clauses as facts (as `F-0001`/`F-0003` here, with `layup psb check`); a generated skeleton, one requirement row per In-Scope fact; role sessions add criteria, split rows (a child keeps the fact ID), write the specification and the phased plan; each row carries a byte-exact PSB quote; `layup spec check` checks IDs, quotes, coverage of every In-Scope fact, criteria and task links, with no model call; a counterpart-harness verifier judges agreement of meaning (C-Q8-A/B/C, A-8A/8B, B-8A) | ✓ | checked by the verifier; Jev at most as a hold signal | **yes** |
| Jev extracts requirements (B-8B) | ✗ | ✗ | no |

## 11. Questions for the Operator

Asked on #72 in one comment (plan D4). Each names the options that the members gave.

- **Q-1** Stack gates and their CI job in the target as kit content (no LAYUP code or binary): does this supersede O-10 and O-11, as O-69 superseded O-61? (condition 1; C P6)
- **Q-2** The human identity: a separate account for the Operator and the idea owner before the pilot; a signed commit per human decision; or one account, with the record that no human decision point is provable in the pilot. (C P2, A question 1)
- **Q-3** Jev's authority before its thresholds have evidence: direct from a start value; shadow for milestone 1, then promote per point; or the cautious direction only. (C question 3)
- **Q-4** Escalation screening: before each task only, or also at each handoff. (C question 4)
- **Q-5** A harness that does not report tokens or cannot enforce a spend cap: not eligible for paid autonomous work; or eligible with a wall-clock proxy and incomplete telemetry. (C question 5, A question 3)
- **Q-6** The roles and the owner of each ambiguity kind: the seven PSB §2 functions; the author, the verifier and four answerers; or the pipeline roles of B-3A. (PRD §11 question 3)
- **Q-7** The stall start values T and N: set by the Operator now; or measured in a calibration run before autonomous use. (A question 2, B question 1)
- **Q-8** Routing weights in the pilot: move by a bounded step at each retrospective; or stay frozen until the ADR that ends bootstrap mode. (B question 2)
