# T-hbw8 — the selection for the rewrite (plan v2, step 2)

Written by the author (Claude Opus 5.5) on 2026-09-29, before the design text of
slice A (plan-review conditions 6 and 7). It replaces [`selection.md`](selection.md),
which selected for the withdrawn first draft. The method is Solution selection
(`docs/engineering-discipline.md`): the public search first, then the
considerations. The inputs: the panel (`panel-A.md`, `panel-B.md`, `panel-C.md`,
option names as in `selection.md`, for example A-1A), the evaluation
([`evaluation/summary.md`](evaluation/summary.md), `evN` = `evaluation/NN-*.md`),
the searches of section 1, and O-66 to O-92.

## 1. The searches for the parts with no evaluated pattern

Five read-only search sessions (Claude Opus 5.5, Claude Code subagents, one per
part, 2026-09-29). Each read the `keep` entries of
[`search-v2/searcher-A.md`](search-v2/searcher-A.md) and
[`search-v2/searcher-B.md`](search-v2/searcher-B.md) in full and ran 8 to 11 web
queries in the words of the PSB or the vision brief. Each row below is condensed
from its output; the author opened no URL that is not named.

### 1.1 The escalation screen and the smart-if points

- **Kept entries.** Cordum (policy before the job, `REQUIRE_APPROVAL`), numbat
  (a monitor or enforce switch per rule), Squelette (refuses a change outside a
  recorded human decision), Ouroboros (a model score with a fixed threshold): keep
  as patterns. Rejected: sprix-sage-router (human gates on its roadmap only),
  humanlayer (deprecated), Tenuo, HOL Guard, TraceFold, ActPlane, NeMo Guardrails
  and the like (they govern tool calls, not business decisions), and 18 list rows
  whose gates sit at fixed steps or are raised by the agent itself. Paperclip and
  AgentPlane: evaluated (ev03, ev11), no new reason.
- **Queries** (source): "escalation rule selects business-forking decisions agent
  stops human decides budget legal compliance" (PSB §6 In Scope); "classify agent
  decision business impact escalate to human approval policy" (§8 Business-forking
  decision); "coding agent must stop and ask human before scope trade-off or budget
  decision" (§8, "scope against date"); "audit sample … missed escalation" (§7.1);
  "AI decision point shadow mode then delegate authority confidence threshold human
  fallback" (vision 3.1, 3.5); "Laya decision model" and "TypeSafe Jev System One"
  (vision 3.1).
- **Patterns.** A deterministic floor from the policy, and the model may raise an
  escalation but never lower it (Governed APA, https://doi.org/10.3390/a19080627);
  code decides, a model only compiles the rule once (Policy-Approval-Agent); typed
  questions, ordered tiers, fail-closed escalation, recorded threshold overrides
  (jevals https://github.com/openlayer-ai/jevals, JevLoop
  https://github.com/Ylr9933/JevLoop, jevscript); frozen thresholds and a "who
  decided" trace (layaAgent); a monitor or enforce switch per rule (numbat).
- **Result.** (a) A rule that selects business-forking *decisions* before the
  agent acts: **no public pattern** — every public gate works on tool calls or
  paths. (b) A pluggable classifier with thresholds and a human fallback:
  **public pattern** (jevals, JevLoop); a named authority level per point only as
  numbat's per-rule switch. (c) **Laya exists**: Convai Innovations' open-weight
  decision model (https://github.com/NandhaKishorM/laya), a Jev-compatible
  `systemone` API; weak zero-shot, over-confident as shipped, a confidence that is
  not comparable with Jev's.

### 1.2 Question routing by ambiguity kind

- **Kept entries.** None kept for this part. Rejected: AgentBridge, hcom, CLI
  Agent Orchestrator, Concord MCP (transports or mailboxes; the sender picks the
  peer), clideck (autopilot removed), sprix-sage-router (routes work, not
  questions), Ouroboros (a spec gate, Problem 3), A2A and MCP (transports),
  humanlayer, Calyx, Taskuary (ask-human only), Atlas, vibey, 5dive (spec-time
  interviews), and 13 chat or spawn tools.
- **Queries** (source): "agent clarifying question routed to specialist agent
  instead of human" (§1 item 1); "which agent owns the answer escalate human" (§3
  Problem 1 Root Cause); "coding agent asks question peer agent answers resume
  session" (Problem 1 Downstream Friction); "classify ambiguity domain business
  architecture interface contract environment" (Root Cause, the four kinds); "RACI
  matrix AI agent team … who answers" (the owner per kind); "clarification
  turnaround … question broker" (§7.2).
- **Patterns.** Answer within delegated authority, else escalate, with a written
  list (fveracoechea/operator issue #21, a plan); a conductor session that answers
  a waiting worker and escalates by confidence (agent-deck); a kind taxonomy split
  between a supervisor and domain experts (MAC, https://arxiv.org/html/2512.13154v1);
  the answer resumes the same session id (buildd PR #2526).
- **Result.** (a) Classifying a question into the four PSB kinds: **no public
  pattern**. (b) A default owner per kind with evidence: **no public pattern**
  (MAC is for dialogue, and the user answers). (c) A needs-a-human branch: a
  written escalation list (operator #21) and a confidence rule (agent-deck).

### 1.3 The records that the PSB measures need

- **Kept entries.** aGiTrack (tokens, turns, human prompts in Git notes),
  agenttrace (reads the logs of about 15 harnesses): keep. Worth a look, not
  opened: tu, ax, Vibestrate, scion, Cordum, Even, myc, clu, pro-workflow, crit.
  Rejected as the store: Langfuse, LangSmith, Phoenix and similar (a server
  database); Helicone and screen-only cost tools; budget-stop rows (covered by
  ev03).
- **Queries** (source): "tasks with unplanned human input rate" (§7.2 Task
  Intervention Rate); "human intervention rate AI coding agent" (§8 Unplanned
  human input); "OpenTelemetry semantic conventions gen_ai agent spans" (§7.1
  Telemetry Completeness); "audit random sample of agent decisions overturned"
  (§7.2 Reversal Rate); "missed escalation rate audit sample" (§7.1); "first review
  acceptance rate" (§7.2); "question to accepted answer time p95" (§7.2); "share
  of questions before development" (§7.2 Early Question Share).
- **Patterns.** OTel GenAI field names for tokens and duration (status
  "Development", no cost field); proposed task and human-review spans with
  `requester.type` and `intervention_type` (semantic-conventions #2665,
  OpenLLMetry #3460); a reproducible audit sample with a recorded seed and
  population hash (kla.digital, PCAOB AS 2315); lead time from task assignment
  (agent-dora, Warp).
- **Result.** Telemetry and lead time: **public pattern** (field names only; ev03
  and ev11 already give the ledger). The split of planned from unplanned human
  input, Clarification Turnaround to an *accepted* answer, Early Question Share,
  First-Review Acceptance per requirement, and a baseline stored as data: **no
  public pattern**. Audit sampling: a public method, no schema.

### 1.4 Learning that changes routing

- **Kept entries.** Hivelore (a mistake becomes a guard in Git hooks; one
  repository): keep as a lesson pattern; model-watchdog (rollback only): pattern.
  Rejected: Claudexor (routes by quota headroom, learns nothing), pro-workflow
  (rules in SQLite), Mneme (ADRs to guardrails, no learning), pi-reflect and
  OpenProgram (self-edit rules without approval), memory stores (SAGE, Caura and
  similar), and five thin list rows that state no learning signal.
- **Queries** (source): "route tasks to models by complexity learn from outcomes"
  (vision 2.2); "multi-armed bandit LLM model selection coding tasks" (vision
  3.3); "reward from test pass rate verification token cost routing" (vision 3.3,
  the reward bullet); "lessons learned from one project applied to next project"
  (PSB §5, "the next project gets no lesson"); "automated retrospective after
  milestone agents update model assignment" (vision 2.3); "cross-project defaults
  routing policy" (PSB §5, §6 Out of Scope).
- **Patterns.** Bounded weight proposals that a human must adopt, with rollback
  (guangyang1206/adaptive-model-router); clamped offsets moved by approvals and
  rejections (Neil0619/adaptive-model-router); only observed evidence counts, an
  agent's own outcome report does not (openshard PR #350); unknown tokens stay
  unknown (EricT1230/ModelRouter); reviewed lessons in a Git repository, synced to
  every project (smirnowld/agent-practices PR #16).
- **Result.** (a) Bounded outcome-based routing outside model weights: **public
  pattern in parts**; none keeps its policy in files a clone carries. (b) A reward
  from delivery records: **no public pattern**. (c) Lessons to the next project:
  public for lessons (smirnowld), **none for a learned routing default**.

### 1.5 The order of Intake and setup

- **Kept entries.** Squelette (human decisions and a charter as Git files, checked
  before work), agents-cli and Harness Starter Kit (pinned versions, a drift
  check): keep. Spec Kitty, BMAD, Shape Up: evaluated (ev06, ev07, ev19); new fact
  only: Spec Kitty and BMAD make the repository first. Rejected: MetaGPT (fills a
  gap with a guess), Ouroboros (one question at a time, SQLite), AgentPlane,
  Kiro, Caliber, fab-kit, vibey, Atlas, headless-claude-template.
- **Queries** (source): "check problem statement for gaps before development
  starts" (§3 Problem 3 Root Cause); "ask all clarifying questions in one batch"
  (§6 Decision Point 2); "copier answers file … reproducible template" (§6
  Reproducible Discipline Setup, I8); "record answers … evidence" (I4); "create
  GitHub repository … branch protection as code … pinned template version" (§6,
  I8); "idea owner sets appetite success criteria before project starts recorded
  in git" (Decision Point 1); "bootstrap … then create repository scaffold CI
  rulesets" (I1).
- **Patterns.** Make the repository with the baseline first, record the intake as
  commits in it, and apply the forge settings (rulesets, CI, credentials) as a
  later step (blue126/agent-project-bootstrap, Spec Kit, agent-project-starter);
  an answers file with the pinned template commit (copier `.copier-answers.yml`);
  the appetite and the bet as one file in Git (openproj).
- **Result.** (a) An order that keeps the answers in Git: **public pattern**,
  repository first. (b) A deterministic gap check in one batch: **no public
  pattern** (only research phrase lists; every tool uses a model). (c) Evidence per
  setup value: **no public pattern** (copier and blue126 record the value and the
  pin, not the source).

## 2. The selection per ADR

The ADR list is the one of plan v2 as amended (slice A 3, B 2, C 1, D 2, E 2, F 1,
G 1, H 1). That list has thirteen entries; the plan's word "twelve" miscounts it
(reported on #72). "Rejected" names what a later reader must not reopen without
new evidence.

| ADR | Slice | Compared (panel; evaluation; search) | Selected | Rejected, and the tradeoff |
| --- | ----- | ------------------------------------ | -------- | -------------------------- |
| 0013 Orchestrate a target from outside | A | A-1A, B-1A, C-Q1-A (a foreground run); A-1B, B-1B, C-Q1-B (`--once` on the forge); ev01, ev04 (state in a town or SQLite) | One foreground `layup run` per target on the LAYUP host; the engine checks are pure commands; only the three things of FT5 enter a target | A forge-runner loop (harness credentials on the forge, lock-in); a daemon or database (FT6); a candidate as a runtime (summary: none should run next to LAYUP). Cost: one host must be up during delivery |
| 0014 Keep the records in the target's Git, one writer, the actor from a credential | A | A-4A, A-4B, C-Q4-A/B/C, B-4A/4B; ev11 (typed result, one committer, state hash), ev03 (verified actor, typed decision), ev10 and ev02 (the commit author, as the failure) | A records branch of the target; TSV and Markdown; only `layup run` commits, with the LAYUP App's installation token; a human decision is an issue comment by a named human account with no App; reads of the rulesets and the push actors fail closed | Records on the default branch (a commit per event blocks open PRs); per-task refs (a clone misses them); the commit author as proof (FT3); a PR review as a human decision (no App field) |
| 0015 Keep model calls out of the engine checks | A | C's decision-8 text in `selection.md` §4; C-Q3-C (one neutral rule source), A-3A (a probe); ev04 (a harness table), ev11 (a clean start), ev01 (`prime`) | Successor of ADR-0011 decision 8: the engine checks start no model; the `layup` process calls a model only through the smart-if provider; a role session is a harness process with a home, an environment and a work tree that LAYUP gives it, and no forge credential | A harness SDK inside LAYUP (standard library only, `F-0004#1`); a session with the Operator's home (Fable-M23; the evaluation incidents) |
| 0016 Put the native stack gates in the target | B | A-5A/5B, B-5A, C-Q5-A/B (in the target), C-Q5-C (export), B-5B (runner only); ev14 (a catalog; the PR's own copy as the failure), ev17 (test-quality rules), ev13 (rules from the base) | A stack catalog in LAYUP with a known-bad fixture per gate kind; setup writes the tools, their configuration and one CI job per gate kind; `layup gate` runs the base branch's definitions from outside; a gate that is not active fails | Gates only outside (finding A2, O-76); the PR's own copy as the only gate (FT4). Cost: the LAYUP host needs each stack's toolchain (K17) |
| 0017 Prevent rule changes by agents | B | O-69; C-P6; ev14 (a ruleset applied and read back), ev11 (a hash of the rule files) | Sessions hold no forge credential; rulesets with no bypass on the default branch and only the LAYUP App on the records branch; required checks pinned to the LAYUP App; a rule path changes only in a batch that a human approves at a planned point | Code owners with a second human account (O-77: no second account); detection by audit only (Sol-1). Limit: a process on the host can reach the Operator's own credentials |
| 0018 Derive the specification from numbered source lines | C | A-8A/8B, B-8A, C-Q8-A/B/C; B-8B (Jev extracts); ev07 (`covers`; a trace to a missing line as the failure), ev06, ev08 (a guess as the failure), ev09; ev19 (scope hammering) | A session numbers every line; code checks that every line is a fact or is listed with a reason, and that each quote is byte-exact; the idea owner confirms the needs and sets each priority; a counterpart session reviews completeness | Code numbering clauses (K19); code setting `Must` (K21); Jev extracting requirements (it writes no text) |
| 0019 Run the lifecycle as a phase loop with one bet per milestone | D | Target design of O-67; ev11 (a transition table), ev13 (a draft PR until verified), ev04 (merge preconditions), ev19 (appetite, the bet, the pitch) | Intake, Scaffold, Shape, Bet, Build, Accept, Retrospective, then the next Bet; a typed handoff checked by code; a draft PR until gates and verification pass; merge at the verified head | An architecture approval in each milestone as its own step (K41, human load); a review request before the gates (K28) |
| 0020 Route role sessions over registered harnesses | D | B-3A, A-3B, C shared tables; B-3B (fixed pair); B-3C (one Choice); ev06, ev08, ev05 (a step table), ev04, ev18 (a fresh read-only verifier); §1.2 (no pattern for kind owners) | The seven PSB §2 roles and a step table as defaults (O-81); the owner map confirmed by the Operator at Intake; admission in code; the verifier's harness differs from the author's by the ledger; the tier from the computed hill position | A fixed pair (drops vision 2.2); a judgement deciding Invariant 9; an owner map with no evidence (Sol-32) |
| 0021 Branch at named points through a smart-if provider | E | B-2A, A-2B (a floor, literal packs, code composes); A-2A; B-2B (a model proposes); B-2C (no provider); §1.1 (jevals, JevLoop, numbat) | A provider interface (Jev first; Laya through the same API); five named points; the authority per point `off`, `shadow`, `cautious` or `delegate`, set by the Operator per target; a safe branch per point; O-71 fallback | Mandatory `shadow` (Sol-6, O-78); a point not on the list (Fable-M13); a generative model on the live path |
| 0022 Screen for business-forking decisions before the work | E | B-2A/6A, C-Q6-A (before the task), C-Q6-B (at each handoff); §1.1 (Governed APA: the model may only raise) | The session stops on a candidate decision; a floor over proposed diffs; four questions to the provider that may only raise; before each task and at each handoff (a parameter, O-79) | The agent's own list as the only input (Fable-M7); a floor that claims to read prose (Author-4) |
| 0023 Stop a stall at a limit and diagnose it with a fresh context | F | A-1A (progress), C-Q1-C (dead-man); ev03, ev01 (round and attempt limits), ev04 (wait states), ev18, ev15 (a deterministic aggregator), ev19 (the circuit breaker, the computed hill) | Five triggers; progress is a computed hill move, not a new record; a fresh diagnosis before any human; a blind panel with an aggregator; one package; the Operator's reroute | A clock that fires on waits (Fable-M5); a new record as progress (Sol-22); a panel with no exit (Sol-24) |
| 0024 Record every session's cost and stop at the milestone cap | G | C-Q6-C, B-6A (a band); A-6A/6B (per-action caps); ev03 (a ledger, budget scopes, a hard stop; zero for a subscription as the failure), ev11, ev10; ev19 (appetite); §1.3 (field names) | A ledger row per session, tokens `observed`, `partial` or `unavailable`; the appetite and the band at Intake; a cap per bet; a stop that kills the session; an unknown cost is never zero | Per-task estimates at Intake (K52); a harness's own price estimate judged by the provider (B-6B) |
| 0025 Learn routing from the records at each retrospective | H | B-7A, C-Q7-A/C, A-7A; A-7B (exploration); B-7B; B-7C; §1.4 | A reward computed in code; a bounded step; a proposal that the retrospective batch approves (the default), or applied within bounds (O-83); the learned weight ranks the admitted pairs; LAYUP's defaults change by a LAYUP change for the next project | The weight as a tie-break only (Fable-M17); exploration trials before a baseline; the provider proposing weights |

## 3. Plan-review note 7: isitdone and the Spec Kitty packet

The correct statement is the plan's: no candidate runs next to LAYUP, and no
candidate file enters a target. `evaluation/summary.md` says "usable" for two
of them; this selection does not use either as a program or a file. LAYUP
borrows the patterns: isitdone's test-quality findings (a skipped test, a
deleted assertion) become rules of the Go gate catalog and lines of the
verifier's checklist (ADR-0016, ADR-0020); the fields of Spec Kitty's handoff
packet v1 inform the handoff schema (ADR-0019). The reason: LAYUP is Go with the
standard library only (`F-0004#1`), and O-76 keeps every LAYUP-only file out of
a target.
