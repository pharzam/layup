# T-hbw8 — panel member C (governance, invariants, and human decision points), word for word

Model and harness: Claude Opus 5.5 (`claude-opus-5-5`) on Claude Code (`claude -p`), a fresh session. Time: 2026-09-27T18:00:58Z to 18:05:38Z (4 min 40 s). Brief: `COMMON.md` and `DOMAIN-C.md` (in `panel-brief.md`), the tree at `6721419`, and the Jev pages of `jev-sources.md`. One pass; the member saw no other member's output (ADR-0006; plan of #72, D2).

The output follows between `--- begin` and `--- end`.

--- begin
# Panel member C: governance, invariants, and human decision points

Task T-hbw8, issue #72. One pass. Inputs: F-0001, F-0003, F-0004, F-0002, PRD-0001, ADR-0010 to ADR-0012, `runs/T-hbw8/inputs-from-pr-69.md` (O-66 to O-75, the review), `runs/T-hbw8/jev-sources.md` and `.panel-brief/jev/`. This record generates and compares options. It does not recommend one.

The lens: for each option, can a human become the message bus again (`F-0003#1`, `#9`, `#10`), and does each human input land at a Human Decision Point (`F-0001#10`–`#14`) or count as unplanned (`F-0001#28`, Task Intervention Rate `F-0003#70`)?

## Points where the target design is at risk (separate from the options)

- **P1. Jev cannot write a panel synthesis.** The target design names "the synthesis of a panel" as a Jev decision point. The `concepts/system-one` page says System One models "do not write replies, produce code, or generate explanations of their reasoning". Jev can only select (Choice) or score (Score) among panel outputs. A reasoning-model session must write the synthesis text.
- **P2. A human comment cannot be told from an agent comment today.** ADR-0011 records the fact of O-9: "the agents and the Operator push with the same GitHub account today". O-73 makes an issue comment the form of every human decision. If an agent can write that comment, an agent can approve its own work (`F-0001#3`) and make a business-forking decision (`F-0001#13`, Missed Escalations `F-0003#57`). Every option in Q4 to Q6 below assumes a separate human account. Without it, no option is sound.
- **P3. A Jev threshold with no data is a guess.** The target design says "each threshold is a start value with evidence". The `confidence` page gives no number: "Start with conservative thresholds, test with your own data, and adjust." So before the first pilot data, any number breaks `F-0001#4`. Option Q2-B below is one answer.
- **P4. Jev is weak at the numbers that the budget decision needs.** The `model-jaggedness/jev-1.13` page lists "math and numbers", "counting", "date and time comparison" as weak points: "Keep the arithmetic in code"; "compare in code". So code computes spend, the band, and the clock. Jev gets only the semantic question.
- **P5. The retrospective approval is planned only if it is listed before delivery.** `F-0001#25`: a planned approval point "is recorded in the project repository before delivery starts". O-69 is sound only if Intake writes the retrospective approval of each milestone into the list of planned points, with its approver (`F-0001#12`). A rule change at any other time stays unplanned (`F-0001#28`), as O-69 says.
- **P6. O-10 and O-11 contradict A2.** O-11 says stack gates "run outside the target". `F-0003#65` says "The project repository passes its gates without the automation." The target design puts the stack gates into the target. That is a change to O-11, and it needs the Operator's word (see Questions).

## Q1. The orchestrator `layup run TARGET` [B1, B2]

**Q1-A. Foreground loop with Git state and a clock watchdog.**
- Mechanism: one Go process. It reads `runs/phase.tsv` (columns: ts, milestone, phase, event, actor, sha) from the target, and appends one row per transition. It starts each role session as a child process with a headless harness command from `docs/harnesses.tsv`, and gives it a prompt file that names the handoff record to read. It reads `runs/<task>/handoffs.tsv` after the session exits. A watchdog goroutine records a progress event on each new commit, handoff row or telemetry row; no progress in T, or N identical failing gate fingerprints, opens a stall. A session that burns tokens with no push is caught by the same clock.
- Rests on: `F-0003#1`, `#17`, `#18`, `#49`; `F-0001#1`, `#14`; O-67.
- Strains: the orchestrator itself can stall or crash, and a human then restarts it: that is a restart (`F-0001#28`) for each open task. The option must record each restart as unplanned input, or interventions hide (`F-0001#14`).
- Falsifiable: in the pilot, a seeded session that loops with no push is not stopped within T plus one watchdog tick.
- Cost: one host that stays on during delivery; no model call; the Operator restarts on a crash.

**Q1-B. Event-driven steps on the forge runner, with a scheduled watchdog.**
- Mechanism: `layup step TARGET` does one transition and exits. The forge runner calls it on each push, comment and review event, and on a schedule (for example every T/2) for the stall check. Each role session runs as a runner job. State is only `runs/phase.tsv` in Git.
- Rests on: `F-0001#1`, `#2`; `F-0003#38`; O-73.
- Strains: GitHub lock-in grows from a check runner to the whole loop (`F-0003#40`); runner minutes and job time limits cap a session.
- Falsifiable: a scheduled watchdog run is delayed or dropped by the runner, and a stall is seen later than T by more than one period.
- Cost: runner minutes for each role session; no host; harness credentials on the runner.

**Q1-C. Foreground loop plus an independent dead-man check.**
- Mechanism: Q1-A, and the loop writes a heartbeat row to `runs/heartbeat.tsv` and pushes it every H. A scheduled forge job only reads the last heartbeat. If it is older than 2H, it opens a stall record for the orchestrator and notifies the Operator as a stall package (`F-0001#14`), not as a chat message.
- Rests on: `F-0003#17` ("A human finds the stall by chance") applied to the orchestrator itself.
- Strains: two moving parts; the heartbeat push adds commits to Git history.
- Falsifiable: a killed orchestrator is not reported within 2H plus one schedule period.
- Cost: one host, and a small scheduled job.

## Q2. The decision component [A3, B8]

For all options: a file `docs/decision-points.tsv` (columns: point, deterministic_precheck, primitive, question, threshold, threshold_evidence, fallback) lists each named point. The deterministic pre-check runs first; Jev is called only when it returns `undecided` (`F-0001#6`). Each call appends to `runs/<task>/decisions.tsv` (ts, point, state_sha256, question, model_version as the response reports it, answer, probabilities, confidence, threshold, outcome act|human, latency, tokens, price). Jev is pinned to `jev-1.13.0`, never `jev-latest`, because thresholds are tuned against a version (`models` page). A 429, an error, or a result below threshold goes to a human (O-71), never to pass.

**Q2-A. Direct authority at each point.** Jev decides when confidence (Choice) or probability (Noul) is above the per-point threshold. Rests on O-67, O-71. Strains `F-0001#4` until thresholds have data (P3). Falsifiable: the Reversal Rate audit (`F-0003#72`) finds more than the start value of Jev decisions overturned. Cost: one HTTP call per point; lowest human load.

**Q2-B. Shadow first, then promote per point.** For milestone 1 of each pilot, Jev runs and records at each point, but the fallback path (the human or the deterministic default) decides. At the first retrospective, a point is promoted to Q2-A only if its agreement with the recorded human decisions supports a threshold; the records are the threshold's evidence. Rests on `F-0001#4`, P3, O-69. Strains the Task Intervention Rate of milestone 1 (all semantic points go to a human), and the Early Question Share (`F-0003#75`). Falsifiable: after promotion, a point's live agreement drops below its shadow agreement. Cost: more human answers in milestone 1; the same Jev calls.

**Q2-C. Asymmetric authority.** Jev acts alone only in the cautious direction (escalate, stop, open a panel, route a question to a human). In the permissive direction (do not escalate, continue over budget, accept an agent answer), Jev needs the threshold and a deterministic corroboration (for example, a gate pass, or the cost projection inside the band). Rests on `F-0003#57` (target 0 missed), `F-0001#13`. Strains the Task Intervention Rate: each cautious false positive is human input, and an escalation that the idea owner does not confirm is unplanned (`F-0001#28`). Falsifiable: the idea owner confirms fewer than half of the escalations in a milestone. Cost: more human input than Q2-A.

**ADR-0011 decision 8, the text to keep:** "The engine makes no model call" stays for every engine check (`layup psb check`, `layup setup verify`, `layup gate`, `layup spec check`). "Its interface is plain text and files, so any harness agent can run it (Invariant 9)" stays whole. Replaced: "judgement stays with the harness agents that call it" becomes "judgement stays with the harness agents, and at the named decision points of `layup run`, with the decision component". A target never needs Jev to pass its gates (`F-0001#2`).

## Q3. Squads and routing [B5]

For all options: `docs/harnesses.tsv` (harness, headless_command, version, models, reports_tokens yes|no, status, evidence) and `docs/routing.tsv` (role, task_class, author_harness, author_model, verifier_harness, verifier_model, evidence, set_at_retro). A deterministic check fails a row where verifier_harness equals author_harness (`F-0003#66`), and fails a merged change whose verification record names the author's harness.

**Q3-A. Roles by the PSB §2 functions.** Seven roles as PSB §2 names them. Each ambiguity kind maps to a role that owns its answer: domain to Domain Expert, architecture boundary to Software Architect, interface contract to Systems Architect, environment to Software Engineer. Rests on `F-0003#9`, `#46`. Strains: seven roles need seven prompts and more handoffs. Falsifiable: in the pilot, more than half of the questions go to one role. Cost: more sessions per task.

**Q3-B. Minimal roles plus answerers.** Author, verifier, and four answerer roles, one per ambiguity kind. The answerer is a fresh session with the specification in context. Rests on `F-0003#9`, `#46`, `#71`. Strains: an answerer with no stake can answer wrongly, and the Reversal Rate (`F-0003#72`) shows it. Falsifiable: the Reversal Rate of answerer answers is above its start value. Cost: fewer sessions than Q3-A.

**Q3-C. One neutral rule source, harness files are pointers.** (Combines with A or B.) The rules are only in `AGENTS.md` and `docs/`. Each harness entry file (`CLAUDE.md` and others) holds only a link; a deterministic check fails any other content. Rests on `F-0003#26`, `#27`, `F-0001#9`. Strains: a harness that does not follow the link reads no rules. Falsifiable: a verifier on a second harness passes a change that breaks a rule in `docs/`. Cost: none at run time.

Risk common to all: ADR-0012 records 18 failed external harness runs. If only one harness works, no change can get a counterpart verification. `F-0001#5` then blocks the merge, and the task stalls to the Operator. That is correct, but it makes the Operator the message bus for harness access. The pilot must prove two working harnesses before delivery starts.

## Q4. The records in Git [B3, B6, B7]

Records in the target: `runs/<task>/handoffs.tsv`; `runs/<task>/questions.tsv` (qid, asked_by, kind, kind_source precheck|jev|human, routed_to, answer_ref, accepted_by, asked_at, answered_at, human yes|no, planned yes|no); `docs/acceptance.tsv` (req, milestone, delivered_sha, evidence, review_number, decision accept|reject, approver, source_ref, copied_at); `runs/<task>/telemetry.tsv` per action (ts, task, req, role, harness, model, action, tokens_in, tokens_out, latency_ms, wall_ms, price, price_source). A value that a harness does not give is `not reported` and counts as incomplete (`F-0003#60`). Prices come from `docs/prices.tsv` with a cited source (`F-0001#4`).

**Q4-A. The orchestrator copies, then acts only on the copy.** `layup run` polls the forge, copies each comment or review verbatim into `runs/<task>/forge/<id>.md` with author, URL, time and SHA-256, commits it, and only then reads it. It never acts on a forge object that is not in Git. Rests on `F-0001#1`, O-73. Strains: nothing is copied while the orchestrator is down. Falsifiable: an audit finds a decision in `decisions.tsv` or `acceptance.tsv` whose source_ref has no Git copy. Cost: forge API calls; no human load.

**Q4-B. A human decision is a commit.** The human writes the decision as a row in a PR to the target (for example `acceptance.tsv`), and the comment is only a notice. Rests on `F-0001#1`. Strains: human load and Git skill for the idea owner; the PR approval itself is unplanned input by the PSB text ("a required approval of a PR is unplanned human input"). Falsifiable: the idea owner writes a decision as a comment and not as a commit in the pilot. Cost: high human load.

**Q4-C. A workflow in the target copies the comment.** A forge workflow in the target commits each comment from a listed approver on `issue_comment`. The target keeps its records without LAYUP (`F-0001#2`). Strains O-10 and O-11 (LAYUP writes a workflow into the target). Falsifiable: a comment edited after the copy differs from the Git copy and nothing detects it. Cost: runner minutes.

For all: the approver of each point is in `docs/approvers.tsv`, written at Intake (`F-0001#12`). A comment from any other account is monitoring, not input. P2 applies.

## Q5. Gates and rule protection [A1, A2, B9]

**Q5-A. Stack gates as target code, written once at Scaffold, owned by humans.**
- Mechanism: for each stack, LAYUP holds a gate recipe (`stacks/<stack>.tsv`: gate kind, public tool, config file, evidence). At Scaffold, `layup setup` writes the tool configuration and a CI job into the target, in the target's own stack (for Go, for example `go vet` and a layout test). The target's CI runs them, and branch protection makes each a required check. `layup gate` also runs them from outside, so ADR-0012 part 6 stays true. The rule paths are in `CODEOWNERS` in the target, with a human account as owner.
- Rests on: `F-0003#65`, `#44`, `#43`, `F-0001#2`, `#3`, `#7`; A2.
- Strains: O-10 and O-11 as written (P6). The recipe values need evidence (`F-0001#4`) for each new stack (B9).
- Falsifiable: a Gate Integrity run (`F-0003#64`) on the target with LAYUP removed misses a known-bad commit that `layup gate` catches.
- Cost: one recipe per stack, written and reviewed once; runner minutes.

**Q5-B. Declarative rules in the target, public tools only.**
- Mechanism: the target holds rules as data, in one neutral form (`docs/rules/layout.tsv`, `boundaries.tsv`, `contracts.tsv`, `test-levels.tsv`). The CI job calls public stack tools (a dependency linter, a contract validator) that read these files. No LAYUP code goes into the target: only data and the tool pins. `layup gate` reads the same files from outside.
- Rests on: `F-0003#27` ("one neutral form in the project repository"), `F-0001#2`, `#7`; O-10 read as "no LAYUP code".
- Strains: a stack with no public tool for a gate kind gives `not-active` (`F-0001#5`), and the target then cannot merge; Generality (`F-0003#67`) can fail on the second stack.
- Falsifiable: for one pilot stack, a gate kind has no public tool, and the gate stays `not-active`.
- Cost: a tool survey per stack; no LAYUP runtime in the target.

**Q5-C. Gates stay outside, and a release exports them.**
- Mechanism: ADR-0011 decision 4 as written. The LAYUP App runs `layup gate` on each PR. At each retrospective and at the end, `layup gate export` writes the gates into the target.
- Rests on: O-10, O-11, ADR-0012 part 6.
- Strains: `F-0003#65` fails for every commit between two exports; that is finding A2 again.
- Falsifiable: an Independence run (`F-0003#65`) at mid-milestone fails.
- Cost: lowest change to ADR-0011.

**Rule-path changes only at the retrospective (O-69), for all options.** The orchestrator puts each proposed rule or gate change into `runs/retro/<milestone>/rule-changes.md` during the milestone, and does not apply it. At the retrospective, one PR holds the batch. `CODEOWNERS` needs the approval of the human account in `docs/approvers.tsv`. A deterministic check in `layup run` fails any merged commit that touches a rule path outside a retrospective PR, and records it as unplanned input (`F-0001#28`), not as an approval. The retrospective is in the list of planned points before delivery (P5). A1 is then closed because the approval is a planned project-level point, not a step inside a task (`F-0001#25`).

Strain common to all: a gate that is wrong and blocks all work mid-milestone can wait for the retrospective only by a stall. That is a planned stall (`F-0001#14`), and it shows in the Stall Rate.

## Q6. Escalation and budget [B7, B8]

The four business-forking axes (`F-0001#26`): budget, legal or compliance position, approved intent, strategic trade-off between approved goals.

**Q6-A. Screen at the plan, before the work.**
- Mechanism: in the Plan phase, each task has a plan record with its decisions listed. The deterministic floor runs first: any decision that changes a budget number, a line of the approved PSB, or a `Must` row is business-forking with no model call. Then one Noul per axis on each other decision, framed so that `true` means escalate, with a max gate (any axis fires). The task does not start until each selected decision has the idea owner's answer, copied into Git (Q4).
- Rests on: `F-0001#13` ("the agent stops and the idea owner decides"), `F-0003#48`, B8; the `cookbooks/sde_cascade` pattern.
- Strains: a decision that appears during Implement is not in the plan. The implementer must then open a decision row and stop, which is the agent's own declaration again.
- Falsifiable: the Missed Escalations audit (`F-0003#57`) finds a business-forking decision made inside Implement.
- Cost: four Noul per planned decision; one idea-owner answer per escalation.

**Q6-B. Screen at the plan and at each handoff.**
- Mechanism: Q6-A, and the handoff schema has a required field `decisions` (text of each choice the role made). The orchestrator screens each handoff before the next role starts. An unscreened handoff fails the schema check.
- Rests on: `F-0003#45`, `#57`; B8.
- Strains: a decision is screened after the role made it, but before its result goes on; this is later than "before the work" for that role.
- Falsifiable: a seeded business-forking decision inside a role session reaches a merged commit.
- Cost: four Noul per handoff decision.

**Q6-C. Budget: band plus a cost stop in code.**
- Mechanism: `budget.md` (O-68) holds, per milestone, the budget B and the band upper edge U, written by the idea owner at Intake (`F-0001#10`). Code computes spend from `telemetry.tsv` and a projection (spend so far plus the median cost of the remaining tasks from the routing record). Below B: continue. Between B and U: Jev decides (Noul "is the remaining work likely to reach the goal with this approach"), with the arithmetic in code (P4). At U, or with projected spend above U: hard stop and escalation, with no Jev call (`F-0001#26`). Per task, a stop at k times the task estimate kills the session (B7).
- Rests on: `F-0003#22`, `#39`, O-68, P4.
- Strains: a harness that does not report tokens (`not reported`) blinds the stop. Then the stop must use wall-clock as a proxy, or the task must not run on that harness.
- Falsifiable: in the pilot, a milestone's recorded spend passes U before the stop fires.
- Cost: no model call below B; one Noul per check in the band.

A known limit for all: the idea owner confirms each escalation (`F-0001#27`). An escalation that is not confirmed is unplanned input (`F-0001#28`). So a sensitive screen lowers Missed Escalations and raises the Task Intervention Rate. The pilot must read the two together.

## Q7. The learning loop and the retrospective [table C 2.3, 3.3]

The reward, for all: a formula in code over the records of the milestone, per routing row: requirement accepted at first review (`F-0003#69`), verifier findings, stall (`F-0003#73`), unplanned input (`F-0003#70`), cost (`F-0003#74`). It does not include the first-attempt gate pass rate, because the PSB says it "is monitored only, so that no one weakens a rule" (`F-0003#58`). No model weights change (`F-0003#53`).

**Q7-A. Deterministic update with a margin.** At the retrospective, `layup retro` computes the reward per row. A row changes only if the evidence count is at least n and the difference passes a margin. The change is a PR in the retrospective batch, approved at the planned point. Rests on `F-0001#4`, O-69. Strains: slow; few tasks per milestone give little evidence. Falsifiable: after three retrospectives, no routing row changed. Cost: no model call.

**Q7-B. Bounded exploration.** A share e of tasks uses a different routing row (a bandit), inside the budget band only. Rests on F-0002 §3.3, O-67. Strains the budget (`F-0003#22`) and the Stall Rate. Falsifiable: the explored tasks cost more than the band allows. Cost: more tasks on weaker routes.

**Q7-C. Lessons as data.** The retrospective session writes lessons as rows (`runs/retro/<milestone>/lessons.tsv`: lesson, evidence_ref, proposed change, target file). Jev Score ranks them; the reward still decides routing; rule changes go to the O-69 batch. Rests on F-0002 §2.3, `F-0003#39` ("the next project gets no lesson"). Strains: Jev ranks prose it did not see run. Falsifiable: a lesson ranked high is reverted at the next retrospective. Cost: one Score per lesson.

## Q8. Specification synthesis [B4]

**Q8-A. Fact IDs and a coverage check.**
- Mechanism: at Intake, `layup facts number` numbers each clause of the PSB, as F-0003 does here. In Design, a role session writes PRD rows that cite fact IDs. `layup spec check` fails if an ID does not resolve, if a cited fact is not a byte-exact substring of the PSB, if an In-Scope fact has no `Must` row, or if a row has no acceptance criterion. A verifier on a counterpart harness checks that the row's meaning agrees with the fact.
- Rests on: `F-0003#51`, `#62`, `#24`, `F-0001#6`, `#39`.
- Strains: an ID trace does not stop a paraphrase that changes the need (`F-0003#25`).
- Falsifiable: at Accept, the idea owner rejects a requirement because its text is not what the cited fact says.
- Cost: one session and one verifier per PRD.

**Q8-B. A verbatim quote in each requirement.**
- Mechanism: Q8-A, and each requirement row has a column `quote`, a byte-exact span of the PSB. The check fails if the span is not a substring. Then one Jev Noul per row: "does this requirement state a need that the quote does not state?" A `true` above threshold sends the row to the verifier, not to a human.
- Rests on: `F-0003#25`, `#62`; P3 for the threshold.
- Strains: Jev literal reading is a known weak point (`model-jaggedness`); a false flag costs a verifier session, not human input.
- Falsifiable: the Noul flags fewer drifted rows than the verifier finds.
- Cost: one Noul per row.

**Q8-C. A generated skeleton, and agents only add.**
- Mechanism: `layup spec draft` writes one requirement row per In-Scope fact, with the fact text as the requirement text and an empty acceptance criterion. The role session only adds acceptance criteria, splits rows (each child keeps the parent's fact ID), and writes the specification and the phased plan. `layup spec check` fails if a fact text was changed.
- Rests on: `F-0003#24` ("Nothing converts it into identified requirements"), `F-0001#6`.
- Strains: a PSB In-Scope list in prose that is too coarse gives large requirements; Delivery Lead Time per requirement (`F-0003#68`) grows.
- Falsifiable: in the pilot, the role session had to rewrite more than a quarter of the generated rows.
- Cost: lowest model load.

For all: the trace continues down. Each task names its requirement; each telemetry row names it (Cost per Requirement `F-0003#74`); `acceptance.tsv` closes it (`F-0001#12`). The Architect phase ends at a planned approval point (O-67) listed before delivery (P5).

## The forge question

I give no option. GitHub only for the pilot does not break an invariant if Q4 puts each forge decision into Git before use (`F-0001#1`). The larger risk is P2, the shared account, not the forge.

## Questions for the Operator

1. **P6, O-11 and Independence.** Stack gates in the target conflict with O-11 as written. Options: (a) change O-11 to allow stack-gate configuration and a CI job in the target, as data and public tools (Q5-A or Q5-B); (b) keep O-11 and record `F-0003#65` as a known limit of the pilot (Q5-C).
2. **P2, the human account.** Options: (a) a separate GitHub account for the Operator and the idea owner before the pilot, and agents never use it; (b) a signed commit by the human for each decision (Q4-B); (c) keep one account, and record that no human decision point can be proven in the pilot.
3. **Jev authority at the start.** Options: Q2-A (direct, thresholds from a start value), Q2-B (shadow for milestone 1), Q2-C (cautious direction only). No fact settles it: `F-0001#4` points to B, and the Task Intervention Rate points to A.
4. **Escalation timing.** Options: Q6-A (plan only) or Q6-B (plan and each handoff). The trade-off is Missed Escalations against Jev calls.
5. **Tasks on a harness that does not report tokens.** Options: (a) not allowed, because the cost stop is blind; (b) allowed, with a wall-clock stop as the proxy and the task counted as incomplete telemetry.
6. **The list of roles** (PRD §11 question 3). Options: Q3-A (PSB §2 functions) or Q3-B (author, verifier and four answerers).
--- end
