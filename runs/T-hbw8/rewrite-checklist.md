# T-hbw8 — the checklist of the rewrite (plan v2, step 1)

Written by the author (Claude Opus 5.5) on 2026-09-29, before the design text of
any slice. It is the evidence for the coverage tables and the known limits of
`docs/architecture.md` (plan-review condition 8); those tables live there, not
here.

**Sources.** `Sol-N`: [`deep-check-sol.md`](deep-check-sol.md) (29 material, 3
notes). `Fable-MN` and `Fable-NN`: [`deep-check-fable.md`](deep-check-fable.md)
(23 material, 15 notes). `Author-N`: [`deep-check-author.md`](deep-check-author.md)
(10 material, 9 notes). `#69-A1` … `#69-C10`: the review of #69 in
[`inputs-from-pr-69.md`](inputs-from-pr-69.md) §1 (table C rows numbered 1 to 10 in
source order). `S1` … `S12`: the In-Scope items `F-0003#41` … `#52`. `I1` …
`I9`: the invariants `F-0001#1` … `#9`. `O-NN`: [`operator-decisions.md`](operator-decisions.md)
and `inputs-from-pr-69.md` §2. `FT1` … `FT6`: the six failure tests of plan v2.
Every note of every deep check has a row (plan-review note 6).

**Pattern.** `evNN` is `evaluation/NN-*.md`; `none yet` means that slice's
search in `selection-v2.md` decides it.

**Status.** `open` until the slice closes it; then `answered: W-nn step k`, or
`known limit` with the reason (it then goes into the known-limits section of
`docs/architecture.md`).

**One change from the plan.** Plan v2 puts the fix of `PRD-0001` REQ-001 in slice C. Row K12 puts it in slice B, because W-01 (slice B) walks the gap check that REQ-001 states.

## 1. Deep-check findings, merged

| Row | Sources | The defect, in one line | Slice | Pattern | Status |
| --- | ------- | ----------------------- | ----- | ------- | ------ |
| K01 | Sol-3, Fable-M12 | "Only `layup run` writes the records" while gate, learn, retrospective, dead-man job and idea owner also write them | A | ev11 typed result, one committer | open |
| K02 | Sol-5, Fable-M3 | The actor is read from a commit author, and a PR review has no `performed_via_github_app` field | A | ev03, ev10 (failure), ev11 | open |
| K03 | Fable-M1 | One App gives agent sessions every right of the orchestrator (merge, ruleset) | A | ev03 | open |
| K04 | Fable-M23 | The host's user-level harness setup is a second rule source that one harness gets and another does not | A | ev04, ev11 clean start | open |
| K05 | Sol-30 (note) | A PR review is not an issue comment; the agent/human field cannot be filled for it | A | none yet | open |
| K06 | Fable-N2 (note) | "Only the smart-if calls a model" fails if role sessions, children of `layup run`, are counted | A | — | open |
| K07 | Author-14 (note) | Vision 3.4 "exclusively through issues" is a deviation (Git records), stated as "answered" | A | — | open |
| K08 | Fable-M10, Author-6 | Intake writes records, runs sessions and needs an issue before the target, its records and its harness register exist | B | none yet | open |
| K09 | Sol-1, Fable-M2 | Rule protection is detection by audit; a ruleset that could prevent it is not used on the default branch | B | ev14, ev11 | open |
| K10 | Sol-26 | ADR-0011 decision 4 still says the stack gates run outside the target (against O-76) | B | — | open |
| K11 | Sol-16 | No run path for the baseline measurement of PSB §7.2 and the idea owner's batch of start values (`F-0004#11`) | B | none yet | open |
| K12 | Fable-M20 | `PRD-0001` REQ-001 still says `layup psb check` finds the gaps of a problem statement (only rules G1–G5) | B | — | open |
| K13 | Fable-N9 (note) | The rule-path list has no files; a step-7 lesson in the same PR is refused | B | ev14 | open |
| K14 | Fable-N10 (note) | The CI form of a `not-active` gate is not stated (a missing job looks like a pass) | B | ev14, ev17 | open |
| K15 | Fable-N11 (note) | Rulesets on GitHub Free work on public repositories only; a private target needs a plan | B | — | open |
| K16 | Fable-N14 (note) | Setup steps S02 (`npx degit`) and S12 (LAYUP's `setup-check` job in the target) contradict ADR-0011 decision 7 and O-76 | B | — | open |
| K17 | Author-17 (note) | `layup gate` from outside needs each target stack's toolchain on the LAYUP host | B | ev14 | open |
| K18 | Author-15 (note) | "0 agent writes to rule paths" against a retrospective batch that an agent writes | B | ev11 | open |
| K19 | Sol-17, Fable-M21, Author-5 | Code cannot number the clauses of prose or find "In Scope" without a fixed form | C | ev07, ev06 | open |
| K20 | Sol-18 | A byte-exact trace proves links, not that no need was missed | C | ev07 | open |
| K21 | Sol-19, Author-18 | Code sets `Must` for each In-Scope fact; priority is the idea owner's intent | C | ev19 scope hammering | open |
| K22 | Sol-20 | No check that each delivered requirement has a specification | C | ev07, ev09 | open |
| K23 | Fable-N8 (note) | A requirement traced to an answers fact (like NFR-007 to `F-0004#1`) fails a check against the problem statement only | C | — | open |
| K24 | Sol-21, Fable-M14, Author-8 | Handoff validity is a column an agent writes; no schema, transition table or check | D | ev11, ev12, ev09 | open |
| K25 | Sol-9, Author-2, Fable-N15 | No path from a question to an accepted answer back in the asking session; no actor accepts | D | ev04, ev05, ev13 | open |
| K26 | Sol-10 | Role questions do not use the issue thread (O-73, vision 3.4, R6) | D | ev04 | open |
| K27 | Sol-27, Fable-M9 | No issue and PR lifecycle in Implement; the target's own gate jobs reject every agent PR | D | ev13, ev04 | open |
| K28 | Sol-2 | Required checks gate the merge, not a human's review of an open PR | D | ev13 draft PR | open |
| K29 | Fable-M16, Sol-23 | No role per step, no task class or complexity input, no author of the prompt | D | ev06, ev08, ev05, ev19 | open |
| K30 | Sol-29, Author-11 (note) | Context routing (vision 3.1) is a prompt file name with no selector or size bound | D | ev04 prime | open |
| K31 | Sol-32 (note) | The owner of each ambiguity kind has no evidence | D | none yet | open |
| K32 | Sol-8, Fable-M7, Author-4 | The escalation screen runs after the work, on the agent's own list, and its floor claims to read prose | E | none yet | open |
| K33 | Sol-6 | Mandatory `shadow` blocks the Operator's per-target authority (O-78) | E | — | open |
| K34 | Sol-7 | `shadow` and the provider fallback put routine decisions on a human | E | — | open |
| K35 | Fable-M8, Author-3 | Question routing trusts the asker's label and has no "needs a human" branch | E | none yet | open |
| K36 | Fable-M13 | ADR-0019 and ADR-0020 call the smart-if at points that ADR-0014 does not name | E | — | open |
| K37 | Fable-M22 | No interface for the Operator to set a parameter during a run | E | ev18, ev03 | open |
| K38 | Fable-N7 (note) | `cautious` has no safer branch at some points; promotion has no ground truth | E | — | open |
| K39 | Fable-N1 (note) | `F-0001#14` gives stalls to the Operator, not to the idea owner | E | — | open |
| K40 | Author-13 (note) | Vision 3.1 names Jev and Laya; Laya is never researched | E | none yet | open |
| K41 | Author-19 (note) | An architecture approval in each milestone adds human load | E | ev19 bet per milestone | open |
| K42 | Sol-22, Fable-M4, Author-1 | A new record resets the stall clock; a two-role disagreement is never a stall | F | ev03, ev01, ev19 hill | open |
| K43 | Fable-M5 | The clock opens a stall at each wait for a human | F | ev04 wait states | open |
| K44 | Fable-M6 | The stall path can reach the Operator with no diagnosis; the examiner needs a third harness | F | ev18 | open |
| K45 | Sol-24, Fable-M18 | "A panel" has no convening, sealed input, isolation, synthesis or exit | F | ev15 | open |
| K46 | Sol-4, Author-10, Fable-N5 (note) | The dead-man job has no run list, no credential, no write path, no package | F | ev03 leases | open |
| K47 | Fable-M19 | A wrong gate blocks the milestone, so the retrospective that could fix it never comes | F | ev19 circuit breaker | open |
| K48 | Author-16 (note) | No mechanism for the Operator to give a stalled task to another harness | F | ev04 | open |
| K49 | Fable-N6 (note) | `heartbeat.H` "with its reason", not "with its evidence" | F | — | open |
| K50 | Sol-11 | The money stop cannot sum a cost when a harness gives no tokens | G | ev03, ev11 | open |
| K51 | Sol-12, Fable-N4 (note) | The O-80 default makes Telemetry Completeness (`F-0003#60`) fail by design | G | ev03 | open |
| K52 | Fable-M11, Author-7 | Per-milestone and per-task budgets at Intake, before they exist, written by a person who cannot push | G | ev19 appetite, ev03 | open |
| K53 | Sol-14, Fable-M15 | Task Intervention Rate misses restarts, corrections, PR approvals and gate changes | G | none yet | open |
| K54 | Sol-15, Fable-N12 (note) | Early Question Share has no project-wide question population | G | none yet | open |
| K55 | Sol-28 | Clarification Turnaround has no asked-at and accepted-at times | G | ev04 | open |
| K56 | Sol-13, Fable-M15, Author-9 | Reversal has no record and no audit sample; Missed Escalations has no sample either | G | none yet | open |
| K57 | Sol-31 (note) | First-Review Acceptance has no review boundary | G | — | open |
| K58 | Fable-N3 (note) | "Action" has no definition, so latency and the cost stop have no unit | G | ev03 | open |
| K59 | Sol-25, Fable-N13 (note) | A lesson does not reach the next project | H | ev06, ev07 | open |
| K60 | Fable-M17 | Once a point is `delegate`, the learned weight acts only on a tie | H | none yet | open |
| K61 | Author-12 (note) | The reward has no cost or verifier-finding term; "test pass rates" conflict with `F-0003#58` | H | none yet | open |

## 2. The findings of the #69 review

| Row | Source | The item | Slice | Status |
| --- | ------ | -------- | ----- | ------ |
| P01 | #69-A1 | A rule-path approval counted as a planned point | H (O-69) | open |
| P02 | #69-A2 | The target cannot pass its stack gates without LAYUP | B (O-76) | open |
| P03 | #69-A3 | Every question goes to the Operator | D, E | open |
| P04 | #69-B1 | Nothing starts or moves the role agents | A | open |
| P05 | #69-B2 | The time-limited stall cannot be detected | F | open |
| P06 | #69-B3 | Requirement acceptance has no record | G | open |
| P07 | #69-B4 | Specification synthesis has no component | C | open |
| P08 | #69-B5 | Verification by a different harness is not enforced; rule copies drift | D | open |
| P09 | #69-B6 | Some decisions stay on the forge | A | open |
| P10 | #69-B7 | Cost has no stop before the budget is spent | G | open |
| P11 | #69-B8 | An escalation can come after the work | E | open |
| P12 | #69-B9 | The source of the stack gates is not stated | B | open |
| P13 | #69 forge question | Is a forge lock-in acceptable? | A | open |
| P14 | #69-C1 (2.1) | Derivation of PDR and PRD, phased plan | C | open |
| P15 | #69-C2 (2.2) | Squads, model allocation by complexity | D | open |
| P16 | #69-C3 (2.2) | Cross-verification by a counterpart harness | D | open |
| P17 | #69-C4 (2.3) | Milestones, retrospectives, lessons to the next project | H | open |
| P18 | #69-C5 (3.1) | Model, context and solution routing (Jev, Laya) | D, E | open |
| P19 | #69-C6 (3.2) | Blind panel with a hypothesis posture | F | open |
| P20 | #69-C7 (3.3) | Reinforcement learning on routing | H | open |
| P21 | #69-C8 (3.4) | Communication through issues | A | open |
| P22 | #69-C9 (3.4) | Telemetry per action | G | open |
| P23 | #69-C10 (3.5) | Circuit breaker, diagnostic package, external answers | F | open |

## 3. The In-Scope items, the phase-1 requirements and the invariants

| Row | Item | Walkthrough | Slice | Status |
| --- | ---- | ----------- | ----- | ------ |
| S1 | Problem Statement Quality (`F-0003#41`) | W-01 | B | open |
| S2 | Reproducible Discipline Setup (`#42`) | W-02 | B | open |
| S3 | Rule Protection (`#43`) | W-03 | B | open |
| S4 | Stack-Dependent Gates (`#44`) | W-04 | B | open |
| S5 | Role Handoffs (`#45`) | W-05 | D | open |
| S6 | Autonomous Clarification (`#46`) | W-06 | D | open |
| S7 | Verification on Every Change (`#47`) | W-07 | D | open |
| S8 | Human-on-the-Loop (`#48`) | W-08 | E | open |
| S9 | Stall Resolution (`#49`) | W-09 | F | open |
| S10 | Cost Visibility (`#50`) | W-10 | G | open |
| S11 | Specification Synthesis (`#51`) | W-11 | C | open |
| S12 | Harness-Agent Neutrality (`#52`) | W-12 | A | open |
| R01 | REQ-001 `layup psb check` (phase 1) | W-01 | B | open |
| R02 | REQ-002 `layup setup`, `layup setup verify` (phase 1) | W-02 | B | open |
| R03 | REQ-004 `layup gate` (phase 1) | W-04 | B | open |
| R04 | REQ-007 gates before review or merge (phase 1) | W-07 | D | open |
| R05 | REQ-009 the stall record (phase 1) | W-09 | F | open |
| R06 | REQ-011 the telemetry record (phase 1) | W-10 | G | open |
| R07 | NFR-001 Git is the system of record | W-12 | A | open |
| R08 | NFR-002 the target is independent | W-12 | A | open |
| R09 | NFR-003 no value without evidence | W-02 | B | open |
| R10 | NFR-004 a check not active is not passed | W-04 | B | open |
| R11 | NFR-005 deterministic first; which process calls a model | W-12 | A | open |
| R12 | NFR-006 the pinned baseline | W-02 | B | open |
| R13 | NFR-007 Go, standard library, `git` program | — | A | open |
| I1 | Git is the system of record | W-12 | A | open |
| I2 | The project repository is independent | W-12 | A | open |
| I3 | Agents cannot change the rules or gates | W-03 | B | open |
| I4 | No value without evidence | W-02, W-08 | B, E | open |
| I5 | A check not active is not passed | W-04, W-07 | B, D | open |
| I6 | Deterministic before a model | W-08 | E | open |
| I7 | The domain changes content, not rules | W-04 | B | open |
| I8 | The pinned baseline | W-02 | B | open |
| I9 | A harness is replaceable | W-12, W-07 | A, D | open |

## 4. Acceptance criterion 1: the components and the rules of order

| Row | Item | Slice | Status |
| --- | ---- | ----- | ------ |
| C1 | The engine checks (`layup psb check`, `setup verify`, `gate`, `spec check`) | A, B, C | open |
| C2 | The orchestrator `layup run` | A | open |
| C3 | The decision component (the smart-if) | E | open |
| C4 | The squad manager | D | open |
| C5 | The learning loop | H | open |
| C6 | The runner (the GitHub App and the host) | A | open |
| C7 | The target | A, B | open |
| C8 | The phase loop, from Intake to Retrospective | D | open |
| C9 | The order of decisions: a check, then the smart-if, then a human only at a decision point | E | open |
| C10 | The limits on the smart-if (O-68, O-71) | E | open |

## 5. The Operator's decisions and ADR-0012

| Row | Decision | What the design must keep | Slice | Status |
| --- | -------- | ------------------------- | ----- | ------ |
| D01 | O-67 | LAYUP a deterministic orchestrator of the whole lifecycle; squads; smart-if for escalations, panels, over budget | A, D | open |
| D02 | O-68 | The idea owner writes a band before delivery; the smart-if decides inside it; outside it escalates | G | open |
| D03 | O-69 | Rule-path changes only at the retrospective, approved there; lessons feed learning | H | open |
| D04 | O-71 | Smart-if unavailable or below threshold: a human, never "pass" | E | open |
| D05 | O-72 | The successor of ADR-0011 decision 8: engine checks make no model call; only the decision component calls the provider | A | open |
| D06 | O-73 (design rule) | At each human decision point a person writes an issue comment; the orchestrator copies it into Git | A | open |
| D07 | O-76 | LAYUP strictly external; native stack gates in the target, written at setup; nothing the target needs from LAYUP | A, B | open |
| D08 | O-77 | Agents under the Operator's account through an App user token; the badge tells agent from human; agent sessions get only the App token | A | open |
| D09 | O-78 | A provider interface; the provider and the authority level per target and per point | E | open |
| D10 | O-79 | Thresholds, escalation rules and check frequency are parameters set per run | E | open |
| D11 | O-80 | Default (b): a harness without tokens or a cap may do paid work under a wall-clock limit; telemetry incomplete; per-harness parameter | G | open |
| D12 | O-81 | Default: the seven PSB §2 functions, each ambiguity kind owned by one; the matrix replaceable per project | D | open |
| D13 | O-82 | `T` and `N` parameters; defaults `T` = 10 minutes (maximum), `N` = 1 | F | open |
| D14 | O-83 | Learning adjustments, bounds and triggers are parameters | H | open |
| D15 | O-84 | Parameters tune the PSB rules, never switch one off; smart-if is a conditional branch | E | open |
| D16 | O-90 | Shape Up ideas as options (appetite, circuit breaker, computed hill, bet per milestone) | E, F, G | open |
| D17 | O-92 | Every agent write to GitHub through the App; separation by convention while sessions reach the Operator's login | A | open |
| D18 | ADR-0012 part 6 | The first pilot: `layup` sets up a target and runs its gate from outside; an amendment of ADR-0011 decisions 3 and 4 keeps this true | B | open |

## 6. The six failure tests (every slice)

| Row | Test | Status |
| --- | ---- | ------ |
| FT1 | No check that did not run counts as a pass | open |
| FT2 | Unknown cost is never zero | open |
| FT3 | The actor never comes from text or a commit author | open |
| FT4 | A gate never runs the copy under review | open |
| FT5 | Only three things go into a target: the setup output that O-76 and `F-0003#42` name; the work of the role sessions, through pull requests; and LAYUP's records, at the place that slice A decides. No file goes in that the target needs LAYUP to build, test or pass its gates | open |
| FT6 | No state lives where a Git clone does not carry it | open |
