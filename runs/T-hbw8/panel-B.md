# T-hbw8 — panel member B (decisions, routing, and learning), word for word

Model and harness: Grok 4.7 (`grok-4-7-high`) on Devin (`devin -p --permission-mode dangerous`). Time: 2026-09-27T18:18:10Z to 18:22:38Z (4 min 28 s). Brief: `COMMON.md` and `DOMAIN-B.md` (in `panel-brief.md`), the tree at `6721419`, and the Jev pages of `jev-sources.md`. One pass; the member saw no other member's output (ADR-0006; plan of #72, D2).

The first run, Grok 4.7 (`grok-4.7#xhigh`) on OpenCode (`opencode run --auto`), started at 18:00:58Z and gave no output at all in 15 minutes; it was stopped and skipped (Bootstrap mode rule 4, applied to a panel member by the plan). The same brief then ran on Devin with the same model, in a fresh clone.

The output follows between `--- begin` and `--- end`.

--- begin
# Panel member B — decisions, routing, and learning

Domain: which calls need a semantic judgement, and which a rule settles. Code owns the workflow. Jev returns a typed answer. It does not write text, code, or a reason (System One page). Pin `jev-1.13.0`. An alias can move, and a tuned threshold then moves with it (Models page). Text only. Context is 64k per request, 32k for the state plus the longest question. Price is $0.042 per million input tokens. Output is free. Rate limits can change without notice. A 429 is a retry with backoff, not a new decision. Known weak points: numbers, dates, counting, large state, literal reading, generation (jaggedness page). Keep arithmetic, the clock, and the spend sum in Go.

Split used below. A deterministic check decides first (Invariant 6, `F-0001#6`). Jev is called only when the answer is a snap judgement over a closed set, a yes, or a level. A human decides only at a Human Decision Point (`F-0001#10` to `#14`). If Jev is down, over the rate limit after backoff, over the context limit, or below the threshold, the call goes to a human. It never becomes "pass" (O-71). Start thresholds are conservative and have no pilot evidence yet (Confidence page; `F-0001#4`). The pilot calibrates them. A threshold with no citation is a defect.

**Target-design strain (separate).** The accepted design asks one Noul for each of the four PSB axes to select an escalation. That puts arithmetic and an implied conjunction on Jev. The jaggedness page says keep arithmetic in code, and a Noul must be one literal condition. The options in question 2 replace that shape. They do not replace O-68 or O-71.

**What ADR-0011 decision 8 keeps.** Keep: the engine checks make no model call; their interface is plain text and files, so any harness can run them (Invariant 9); no daemon and no database. Replace only the sentence "the engine makes no model call". The successor is O-72: checks make no model call; only the decision component calls Jev, and only at the named points.

## 1. The orchestrator

### 1A. Foreground state machine, append-only phase log

**Mechanism.** `layup run TARGET` is a foreground Go process. It reads `runs/<task>/phase.tsv` (columns: `seq`, `phase`, `entered_at`, `exit`, `session_id`) and advances only on a legal exit. Phases are Intake, Scaffold, Design, Architect, Plan, Implement, Accept, Retrospective, then the next milestone. It starts a harness session from `docs/harnesses.tsv`, passes the handoff path, and stops the session when the handoff row validates or the clock fires. The process holds the only clock. A stall is: no new handoff or progress row for T, or N retries with no new commit. T and N are cited start values, compared in Go, not by Jev.

**Rests on.** O-67; B1; B2; `F-0003#45`, `#49`; `F-0001#37`, `#1`.

**Strains.** A dead process cannot see its own stall (`F-0003#17`). A second terminal that restarts it is unplanned input (`F-0001#28`) unless the restart is a recorded planned action.

**Falsifiable.** In the pilot, a session that retries and never pushes still has no stall row after T.

**Cost.** One long-lived process per run. No model call for the clock. Human load only on a restart.

### 1B. Same machine, wake on a forge event plus a lease file

**Mechanism.** No always-on process. A GitHub event starts `layup run --once`. Each session writes `runs/<task>/lease.tsv` (`session_id`, `phase`, `deadline`). The runner starts a watch when `deadline` is nearer than the event gap. The watch compares the clock in Go. Same stall rule as 1A.

**Rests on.** O-67; B1; B2; ADR-0011 decision 1 (no daemon).

**Strains.** Between events the clock is off, so B2 returns unless the runner itself polls. That poll is a daemon under another name.

**Falsifiable.** A stall that starts just after an event is unseen until the next event, and that gap is longer than T.

**Cost.** Less idle time. More runner logic. Same human restart risk if the runner dies.

## 2. The decision component

### 2A. Fixed question pack, code composes, three paths

**Mechanism.** One Go package, `internal/decide`. It never runs inside a gate. At each named point it builds a small text state (no raw diffs, no logs — jaggedness: filter first), sends one request with independent questions, and writes `runs/<task>/decisions.tsv`: `seq`, `point`, `model` (the response field, not the alias), `question_id`, `answer`, `probability_or_noul`, `confidence`, `threshold`, `path` (`act`, `hold`, `human`), `input_tokens`. Paths follow the Confidence page: high acts, medium holds for a confirm, low goes to a human. A missing field, a 429 after backoff, or a context overflow is `path=human` (O-71).

Named packs, all literal:

- Escalation, before work (Decision Point 4). Code checks the floor first: the task text names a budget figure, a legal term, a change to approved intent, or an option outside the approved goals (`F-0001#26`). If the floor hits, stop. Do not call Jev. If the floor is silent, one request, four Nouls, each one condition: "Does this task change the approved budget?", "Does it change the legal or compliance position?", "Does it change the approved problem or success criteria?", "Does it trade one approved goal against another, such as scope against date?" An architectural trade-off inside the approved intent is a `false` criterion (`F-0001#26`). Code ORs the four. Any Noul above its start threshold escalates. Any Noul in the middle band escalates too: a missed escalation costs more than an extra human question (`F-0003#57` wants 0).
- Question routing. One Choice, options `domain`, `architecture`, `interface`, `environment`, `human`, `other`. `other` and low confidence go to a human. A table, not Jev, maps the first four kinds to a role. Jev must not infer the owner (`F-0003#9`).
- Stall action. One Choice: `retry`, `examiner`, `panel`, `operator`. Code forbids `retry` when the retry count is already N. Dates and counts stay in the state as labels that code computed.
- Panel synthesis. Members write options. Jev does not generate them (generation is a weak point). One Score per option on "how far this option cites a fact that the other options do not contradict", levels written by the Operator. Code ranks. Below threshold, no synthesis: the pack goes to the Operator.
- Routing. See question 3. Score on fit is allowed. Jev does not see a price number and does not pick a cheaper model by arithmetic.
- Over budget. Code computes spend against the band (O-68). Jev is called only inside the band, one Noul: "Does the remaining work still match the approved task goal?" Outside the band, escalate. Jev never approves money (`F-0001#26`).

**Rests on.** O-71, O-72, O-68, O-67; A3; B8; `F-0001#6`, `#13`, `#26`; `F-0003#9`, `#46`, `#48`.

**Strains.** Start thresholds have no evidence until the pilot (`F-0001#4`). A literal Noul can miss a business fork that the text does not state. That miss is the falsifier, not a reason to skip the floor.

**Falsifiable.** An audited agent decision is business-forking and has no escalation row (`F-0003#57`). Or a Jev call exists whose point is not in the named list.

**Cost.** One small request per point. Input tokens only, at the stated price. Human load on low confidence and on every floor hit.

### 2B. Cascade: cheap model proposes, Jev only verifies

**Mechanism.** A generative model writes a candidate class or a candidate route. Jev asks one Noul per candidate: "Is this class absent from the source text?" Code escalates when any Noul fires. Same `decisions.tsv`. Same human fallback.

**Rests on.** Cascade cookbook (verify, then escalate). O-71. `F-0003#9`.

**Strains.** The proposer is a model call outside the six named points, so it breaks the O-72 successor unless the Operator adds it. It also adds a harness call on the hot path of every question (`F-0003#71`, 120 seconds).

**Falsifiable.** Turnaround at the 95th percentile exceeds 120 seconds only on the propose step, or a proposal that Jev did not verify still moves the phase.

**Cost.** One generative call plus one Jev request per point. Higher human load when the verifier fires often.

### 2C. No Jev on the hot path; Jev only at retrospective

**Mechanism.** Routing, stall action, and question class use tables only. Jev runs once per milestone to score recorded decisions for the learning loop. A live semantic point goes straight to a human.

**Rests on.** O-71 as the default for every live point. Invariant 6.

**Strains.** A3 stays open: every ambiguous question hits a human (`F-0003#1`, `#9`). O-67 puts Jev on the live path. This option rejects that part of the accepted design. The fact that shows the cost is `F-0003#70` (unplanned input).

**Falsifiable.** Task Intervention Rate stays above the start value, and the extra inputs are question routes that a closed Choice could have classed.

**Cost.** Almost no Jev spend. High human load. Conflicts with O-67.

## 3. Squads and routing

### 3A. Two tables, code filters, one Score for fit

**Mechanism.** `docs/harnesses.tsv`: `harness_id`, `invoke` (a command the Operator cites), `models` (a list), `can_verify` (`yes` or `no`). `docs/routing.tsv`: `role`, `tier` (`execute` or `judge`), `harness_id`, `model`, `weight` (an integer), `evidence`. Roles for the pilot, because `F-0001#21` does not set the list: `clarifier`, `designer`, `architect`, `planner`, `implementer`, `reviewer`, `examiner`. The four ambiguity kinds map by a third small table `docs/ambiguity.tsv` (`kind`, `role`): domain to clarifier, architecture to architect, interface to designer, environment to examiner. `human` is not a row. Code admits a row only if the harness is registered, the model is listed, the role matches, and — for a verification — `harness_id` is not the harness that made the change (Invariant 9, `F-0003#66`). Among the admitted rows, one Jev Score asks which admitted pair fits the task text. Levels: `poor`, `adequate`, `strong`. Code picks the highest score above threshold, and uses `weight` only as a tie break. Low confidence uses the highest `weight` and writes `path=hold`. A verification with no other harness registered does not run: the check is `not-active`, which is not a pass (`F-0001#5`).

**Rests on.** O-67; B5; `F-0003#52`, `#66`, `#7`; `F-0001#9`, `#18`, `#21`; vision `F-0002` §2.2.

**Strains.** Two harnesses must be real before any change merges. A table edit is a content change, not a rule-path change, but a bad `invoke` silently stops routing.

**Falsifiable.** A merged change has a verification row whose `harness_id` equals the implementer's, or a routing pick whose pair was not in the admitted set.

**Cost.** One Score per task start. Human load to register the second harness and to cite each `invoke`.

### 3B. Fixed pair, no Jev on routing

**Mechanism.** The Operator writes one execute row and one verify row per role. Code checks the counterpart harness and does nothing else. Weights do not change until the retrospective, and a human edits them there.

**Rests on.** Invariant 9; Invariant 6; O-69 (changes at the retrospective).

**Strains.** Vision §2.2 (allocate by complexity) is unmet. O-67 names that item as in scope. Say so: this option keeps the safety check and drops the semantic pick.

**Falsifiable.** Tasks of very different kinds always get the same pair, and First-Review Acceptance is worse on the kind the fixed pair does not fit.

**Cost.** No routing call. Human load at each retrospective.

### 3C. Jev Choice over the full register

**Mechanism.** One Choice whose options are every `harness_id|model|role`. No code filter before the call.

**Rests on.** Choice primitive. O-67.

**Strains.** The model can pick the same harness for verify and implement. Invariant 9 then depends on a judgement. It can also pick a pair that is not registered if `other` is omitted, or hide behind `other` if it is present. Structural identity belongs in code (jaggedness, row 8).

**Falsifiable.** A verification row names the same harness as the change, and the Choice confidence was above the act threshold.

**Cost.** One Choice. The cost that matters is a false pass on Invariant 9.

## 4. Records in Git

### 4A. Append-only tables, forge text copied by a command

**Mechanism.** Handoffs: `runs/<task>/handoff.tsv` (`from_role`, `to_role`, `artifact`, `schema_ok`). Questions: `runs/<task>/questions.tsv` (`kind`, `role`, `asked_at`, `answered_at`, `human` yes or no). Acceptance: `runs/<task>/acceptance.tsv` (`requirement_id`, `approver`, `verdict`, `comment_id`) at Decision Point 3, before the next milestone. Telemetry: `runs/<task>/telemetry.tsv` (`action`, `requirement_id`, `tokens`, `latency_ms`, `wall_s`, `price_usd` or `not reported`). `not reported` is incomplete (`F-0003#60`). A forge comment is not state. `layup record import ISSUE` copies the comment body into the matching table and stores the comment id. The orchestrator will not advance until that row exists (O-73, Invariant 1).

**Rests on.** B3, B6, B7; `F-0001#1`, `#12`, `#38`; `F-0003#50`, `#59`; O-73.

**Strains.** Import can lag. A decision that lives only in the comment still breaks Invariant 1 until the command runs.

**Falsifiable.** A phase advances with a cited comment id that has no table row, or a telemetry row has no `requirement_id`.

**Cost.** No model call. One import command per human comment. Low human load if the command is the only door.

### 4B. Commit the comment as a blob, table holds the path

**Mechanism.** The import writes `runs/<task>/forge/<comment_id>.md` and a one-line table row that points at it. Same advance rule.

**Rests on.** Same as 4A. Easier audit of wording.

**Strains.** Two copies can diverge if a later comment edit is not imported again. GitHub is the editor; Git must win (`F-0001#1`).

**Falsifiable.** The blob hash and the comment body differ at audit, and the table still says current.

**Cost.** More files. Same human step.

## 5. Gates and rule protection

### 5A. Gates live in the target; LAYUP only runs them from outside

**Mechanism.** Scaffold writes the stack gate scripts into the target, as kit content, not as LAYUP code (O-10, O-11, ADR-0011 decisions 3 and 4). `docs/stack.tsv` in the target lists `gate_id`, `command`, `author` (`scaffold` or `human`), `evidence`. The author is the setup role, recorded, not LAYUP at merge time. The target CI runs those commands. The LAYUP runner also runs them from outside and posts a check whose name is pinned to an account the agents do not use. A required check that did not run blocks the merge (`F-0001#5`). Rule paths (the gate scripts, `docs/*.md` rules, the workflow file) accept an agent write only in the retrospective commit, and only with an approval from that account (O-69, `F-0001#3`). Any other write is unplanned input (`F-0001#28`).

**Rests on.** A1, A2, B9; Invariant 2; ADR-0012 part 6; O-69; `F-0003#43`, `#44`, `#64`, `#65`.

**Strains.** Who writes a gate for a new stack is still a human citation. Generality needs two stacks (`F-0003#67`). The pilot must show one external gate run (ADR-0012 part 6). This option does not invent the second stack.

**Falsifiable.** A target with LAYUP removed fails a gate that `stack.tsv` marks active. Or an agent commit outside the retrospective touches a rule path and still merges.

**Cost.** No model call. Human approval once per retrospective batch. Runner must be up, or the external check is `not-active` and blocks.

### 5B. LAYUP holds the gate commands and the runner is the only enforcer

**Mechanism.** `stack.tsv` lives in LAYUP. The App runs the commands. The target CI has no stack gates.

**Rests on.** ADR-0011 decision 4, read alone.

**Strains.** Invariant 2 and A2. The review already names this as a contradiction. Stated because it is the design that must not return.

**Falsifiable.** Remove LAYUP and the target has no stack-gate command left.

**Cost.** Lower target setup. Fails the independence check.

## 6. Escalation and budget

### 6A. Floor before the session starts; band checked before each action

**Mechanism.** Before Implement, the orchestrator runs the escalation pack in 2A on the task text and the approved intent. A hit stops the session. No harness starts. `budget.md` holds `task_id`, `cap_usd`, `band_low`, `band_high`, written by the idea owner before delivery (O-68). Before each new action, code sums `price_usd` from telemetry. Above `band_high`, stop and escalate. Inside the band, the one Noul in 2A. A row with `not reported` counts as at the cap for the stop, and incomplete for the measure (B7). The stop is a Go comparison. Jev never sees the sum.

**Rests on.** O-68; B7; B8; `F-0001#13`, `#26`; `F-0003#22`, `#48`, `#57`.

**Strains.** A session that spends inside one action can cross the cap before the next check. The band does not stop that action. Say so as a limit.

**Falsifiable.** A task whose summed price exceeds `band_high` has a later action row with no escalation row between them.

**Cost.** One Noul inside the band only. Human decision outside it. No spend on the floor path.

### 6B. Session estimates a price; Jev judges the estimate

**Mechanism.** The harness writes an expected price. Jev Scores whether the estimate is "safe".

**Rests on.** O-68, weakly.

**Strains.** Numbers are a known weak point. A missing estimate recreates B7. The model would judge money. `F-0001#26` forbids that outside the band, and inside the band the useful check is still the sum.

**Falsifiable.** An estimate Jev called safe is followed by a spend past `band_high` with no stop.

**Cost.** A Score per action, and a worse stop.

## 7. Learning loop and retrospective

### 7A. Reward in Go from records; weights change only at the retrospective

**Mechanism.** At retrospective, `layup learn` reads only Git tables. The reward for a routing row is an integer formula, in code: plus one for a first-review accept of the requirement, minus one for a reversal, minus one for a stall that needed the Operator, minus one for an unplanned input, zero when telemetry is incomplete. No model weights change (`F-0003#53`). The command writes `runs/<milestone>/reward.tsv` and a proposed `routing.tsv` where each `weight` moves by a bounded step (start: at most 1). A proposed row with no evidence cell is rejected (`F-0001#4`). The batch of rule-path edits is a separate commit in the same retrospective, approved by the account the agents do not use (O-69). Jev may Score each lesson text once, "is this lesson actionable and tied to a record?", and code drops lessons below threshold. Jev does not edit the weight. The next milestone reads the new weights only after that commit merges. Between milestones the tables are frozen.

**Rests on.** O-69; O-67; vision `F-0002` §2.3 and §3.3; `F-0003#39`, `#53`; `F-0001#28`, `#3`.

**Strains.** A small pilot can move a weight on one lucky task. The bound of 1 limits that. It does not remove it. The formula is a start value and needs a citation from the Operator before the first run (`F-0001#4`).

**Falsifiable.** A `routing.tsv` weight changes in a commit that is not the retrospective merge, or a weight changes with no `reward.tsv` row.

**Cost.** One Score per lesson, optional. Human approval of one batch. No training job.

### 7B. Jev proposes the new weight

**Mechanism.** The state is the reward table. One Score per row, levels `lower`, `keep`, `raise`. Code applies the level.

**Rests on.** O-69; Score primitive.

**Strains.** The Score is not a magnitude (jaggedness: do not interpolate). The policy is then a judgement with no formula to audit. `F-0003#53` is safe only if code still refuses a weight edit that is not a level step.

**Falsifiable.** Two milestones with the same reward table produce different weight tables, and no human edited either.

**Cost.** One Score per routing row per milestone. Harder audit.

### 7C. Lessons recorded, weights never move in the pilot

**Mechanism.** The retrospective writes lessons and the rule batch only. `routing.tsv` stays at the Operator's initial weights until the ADR that ends bootstrap mode.

**Rests on.** ADR-0012 part 6 (the pilot is one external gate run). O-69 for the rule batch. Invariant 4.

**Strains.** Vision §3.3 and O-67 ask for a learning loop. This option defers it. The fact that supports the deferral is `F-0001#4`: a weight moved with no baseline is a guess.

**Falsifiable.** Not applicable as a failure of learning. It fails O-67 if the pilot ends with no reward row at all.

**Cost.** Lowest. Human writes lessons. No routing change to regress.

## 8. Specification synthesis

### 8A. Role agents draft; a deterministic trace check gates the phase

**Mechanism.** Designer and architect sessions write the PRD, the specification, and the phased plan. They are harness work, not Jev (Jev does not generate). `layup trace` checks each requirement row: an id, a byte offset into the PSB file, a criterion id. The offset must land on a fact line. No model call. A row that fails blocks Plan (Invariant 6, `F-0003#51`, `#62`). Jev may Noul, once, "Does the criterion sentence restate the cited fact line?" only as a hold signal. A fail of that Noul does not pass the trace. A low score sends the row to a human.

**Rests on.** B4; `F-0003#51`, `#62`, `#24`; `F-0001#39`, `#6`; vision §2.1.

**Strains.** A byte trace can pass while the words drift. The optional Noul catches some of that and can also false-hold. The gate that blocks is the offset check, not the Noul.

**Falsifiable.** A requirement in Plan has no offset, or the offset does not land on PSB text, and Plan still starts.

**Cost.** Role-agent calls to draft. One optional Noul per requirement. Human only on hold.

### 8B. Jev extracts requirements from the PSB

**Mechanism.** Choice or Score over PSB paragraphs to name requirements.

**Rests on.** `F-0003#51`, weakly.

**Strains.** Generation and counting are weak points. Identifiers would be guessed. `F-0001#39` wants a trace a machine can validate. A model label is not that trace.

**Falsifiable.** Two runs on the same PSB emit different requirement ids, or an id has no offset.

**Cost.** One request per section. Fails the trace measure if used as the gate.

## Forge

GitHub as the only forge for the pilot is a planned known limit. No option. The record rule in question 4 does not need a second forge: state is the Git table, and the comment is only a source the import command reads.

## Questions for the Operator

1. Stall start values T and N. Options: T = 15 minutes and N = 3 (the old review figure, now with a real clock); or the Operator cites a different pair before the first pilot run. A guess is a defect (`F-0001#4`).
2. Reward step. Options: weights move by at most 1 at each retrospective (7A); or weights do not move until the ADR that ends bootstrap mode (7C).
3. Verification when only one harness is registered. Options: the change does not merge (`not-active` is not a pass); or the pilot may merge with the gap recorded as a known limit until the second harness exists.
--- end
