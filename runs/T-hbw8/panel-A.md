# T-hbw8 — panel member A (orchestration, state in Git, and the forge), word for word

Model and harness: GPT-6 Sol (`gpt-6-sol-high`) on Devin (`devin -p --permission-mode dangerous`). Time: 2026-09-27T18:00:58Z to 18:07:52Z (6 min 54 s). Brief: `COMMON.md` and `DOMAIN-A.md` (in `panel-brief.md`), the tree at `6721419`, and the Jev pages of `jev-sources.md`. One pass; the member saw no other member's output (ADR-0006; plan of #72, D2).

The output follows between `--- begin` and `--- end`.

--- begin
# Panel A: orchestration, Git state, and the forge

This is a comparison, not a selection or a vote. The target design in O-67 to O-71 is the boundary for both options in each section. A target contains its adapted Armature kit, its own records, its own stack gates, and its product. It contains no LAYUP program or gate implementation. A proposed file or command below is an architecture option, not a claim that it exists today. All limits and thresholds need evidence before use; no number below is an assumed pilot result.

## 1. `layup run TARGET`: phase and process control [B1, B2; C 3.5]

### Option 1A — One foreground coordinator
- **Mechanism:** A foreground `layup run TARGET` owns a serial state machine: Intake, Scaffold, Design, Architect, Plan, Implement, Accept, Retrospective, then the next milestone. Before delivery, the project records its planned approval points and named approvers. Each transition needs a valid handoff, its exit check, and any mirrored approval; Architect needs the Operator's issue comment, and Accept ends only when every Must requirement has the idea owner's acceptance record.
- **Mechanism:** For each role, a registered Go adapter starts a harness process with an argument list and a scoped checkout. It captures the exit code, output, child process group, elapsed time, and handoff path. It stops the group on cancellation, timeout, or budget stop; a normal exit without a valid handoff fails.
- **Mechanism:** The coordinator commits and pushes transition and attempt events to a protected target Git ref before it starts the next role. On restart it fetches that ref, checks the last committed transition and Git head, and reconciles an interrupted attempt before it starts new work. An attempt ID prevents a duplicate transition.
- **Mechanism:** A clock compares now with the last *durable* progress event, not a stream of text. The adapter records changed artifact hashes, successful gate runs, and visible failed-command signatures. Equal signatures or no new pushed Git head are not progress. If a harness hides inner tool calls, its session deadline is the fallback; do not claim to see its private retry loop.
- **Mechanism:** A bounded retry and a fresh examiner precede a stall package with the attempt history, exit codes, artifacts, diagnosis, and outcome. The CLI shows phase, role, elapsed time, and last progress. Loss of the foreground process needs a new invocation; that restart is recorded and counted as unplanned human input if a human had to cause it.
- **Rests on:** `F-0003#45`, `#49`, `#61`, `#73`; `F-0001#1`, `#14`, `#28`; O-67, O-73.
- **Strains:** `F-0003#70`: manual restarts can make the Operator the process manager. A local event not pushed cannot be recovered from a clean clone.
- **Falsifiable:** Kill the coordinator between two handoffs; a clean clone must resume the correct role once, or mark the interrupted attempt, without a lost decision or duplicate delivery.
- **Cost:** One long-lived local process and a Git push per durable event; harness sessions, examiner calls, and operator restarts when the host stops.

### Option 1B — Short foreground runs with a Git work queue
- **Mechanism:** The same phase transitions are expressed as queued records in a protected Git ref. `layup run TARGET` runs in the foreground for one bounded unit, starts and waits for one harness process, validates its handoff, pushes the outcome, then exits. A new invocation takes the next queued unit; only one active attempt ID per task can hold the Git lease.
- **Mechanism:** A GitHub Actions worker can invoke the same command on each new queue event and on a scheduled sweep. This is a runner, not a daemon or a second state store. It fetches the target ref on each run. If no worker is configured, the CLI reports `not-active`, not autonomous progress.
- **Mechanism:** A sweep reads persisted UTC deadlines and last artifact/head hashes. It detects no pushed progress, a repeated visible failure, an expired worker lease, or disagreement. It stops a still-running child through its adapter where possible; after worker loss, it marks the attempt interrupted, does not assert that the old process died, and blocks a duplicate until that fact is checked.
- **Mechanism:** Retry, fresh-context examination, panel, and Operator routes have finite recorded limits. A human comment that answers a stalled task enters through the mirror in section 4, not through a prompt. Each worker run reports step and heartbeat; a missed sweep is visible as `not-active` until a worker runs again.
- **Rests on:** `F-0003#49`, `#52`, `#61`; `F-0001#1`, `#14`; O-67, O-73.
- **Strains:** `F-0003#71`, `#73`: queue and schedule delays can raise turnaround. A GitHub worker adds service cost and can lose a process while Git still shows its lease.
- **Falsifiable:** Stop a worker without a handoff, restart on a clean runner, and check that the old attempt cannot produce a second accepted transition and the overdue unit gets a stall record.
- **Cost:** More Git and worker starts than 1A; no permanently running host. Hosted job time and recovery checks replace manual foreground restarts.

## 2. Named Jev decisions [A3, B8; C 3.1, 3.2]

### Option 2A — One typed decision request per event
- **Mechanism:** `layup decide` runs exact checks first: schema, gate status, identity, clock, money, band, eligible harnesses, and state transitions. Only the decision component calls pinned `jev-1.13.0` on the remaining semantic question. An event may batch independent questions over a short text state: four separate Noul axes for business-forking impact; Choice for ambiguity kind, stall action, or eligible model/harness; Score only for a stated ordered rubric.
- **Mechanism:** For panel synthesis, role agents write the compared text and evidence; Jev selects a defined disposition or flags conflict. Jev does not generate a panel summary or explanation. A fresh reasoning role writes the synthesis, and a human alone makes an intent or outside-band decision. Unknown kinds and missing evidence stop.
- **Mechanism:** Each call adds a `runs/<task>/decisions.tsv` row: event ID, decision kind, source Git commit/hash, question-set version, model requested/returned, full typed result and probabilities by referenced safe file, threshold source, selected action, fallback reason, UTC times, token usage, and price source. No credential enters Git. Thresholds differ by consequence and need labeled pilot evidence; an unavailable or below-threshold Jev answer stops for a human, never `pass`.
- **Rests on:** `F-0001#1`, `#6`, `#13`, `#26`; `F-0003#46`, `#48`; O-67, O-71, O-72, O-73.
- **Strains:** `F-0001#13`: a confident false negative on an impact axis could let an agent make a business decision. Jev does not prove semantic correctness.
- **Falsifiable:** A seeded legal, budget, intent, or strategic fork reaches implementation without an idea-owner answer, or a low-confidence/failed call continues as `pass`.
- **Cost:** One billed text request per decision event plus a human for uncertain cases; $0.042 per million input tokens is the published Jev 1.13 price, not a fixed future price.

### Option 2B — Evidence gate, then a narrow Jev question
- **Mechanism:** A rule table first maps explicit budget changes, funding requests, legal flags, scope changes, and failed deterministic checks to `stop`. For text with no explicit flag, Jev answers the four business-impact Nouls before work. A second Choice request is made only when code needs a semantic route among responsible roles, stall actions, panel dispositions, or eligible harnesses.
- **Mechanism:** A negative on all four impact axes is not proof of safety: near-boundary values, missing facts, and any Jev error route to the human under O-71. The four Nouls have no separate confidence field; compare their probabilities to separately pre-registered yes/no bands. Choice and Score use their own returned confidence. Code, not Jev, compares dates, prices, counts, and budgets.
- **Mechanism:** The same decision record stores both requests and the exact check that ran first. The resulting route is advisory to the state machine until its normal handoff and gate checks pass. ADR-0011 decision 8 keeps the rule that engine *checks* make no model call; its blanket prohibition on engine model calls and PRD NFR-005's matching claim need a successor for this separate decision component.
- **Rests on:** `F-0001#1`, `#6`, `#13`; `F-0003#48`, `#57`; O-67, O-71, O-72.
- **Strains:** `F-0003#71`: two sequential requests can lengthen a clarification; the hard-coded floor can miss an indirect business fork.
- **Falsifiable:** An indirect strategic trade-off that a pilot labels business-forking bypasses both the floor and the four axes before any work starts.
- **Cost:** Fewer calls for explicit floor matches, but two calls for a routed semantic case; human review of uncertain cases and threshold evidence.

## 3. Squads, roles, and counterpart harnesses [B5; C 2.2]

### Option 3A — Seven named specialist roles
- **Mechanism:** `docs/harnesses.tsv` registers harness ID, executable adapter, supported models, tested capabilities, usage-meter capability, and evidence of each probe. `docs/routing.tsv` maps phase, role, complexity, model tier, eligible harness, and distinct verifier harness. Registry entries with no successful probe are ineligible.
- **Mechanism:** Product Owner, Domain Expert, Systems Architect, Software Architect, Software Engineer, Software Developer, and QA Engineer mirror PSB §2. Domain questions go to Domain Expert; boundaries to Systems Architect; contracts to Software Architect; environment to Software Engineer; business forks to the idea owner. The asking role records whether it accepted the answer.
- **Mechanism:** Before any change advances, a fresh QA or specialist session on a *different harness product* reads a clean clone, runs the target gates, and writes a verification record bound to the change SHA. The routing validator refuses a verifier with the author's harness ID, even if the model differs.
- **Rests on:** `F-0001#9`, `#18`, `#21`; `F-0003#46`, `#52`, `#66`; O-67, O-73.
- **Strains:** Many role handoffs can raise cost and slow a small task (`F-0003#35`).
- **Falsifiable:** A change made on harness A passes with only another model on A, or a domain question always goes to the Operator.
- **Cost:** More sessions and role records; no new database; adapter and capability tests for each harness.

### Option 3B — Four owner groups with a separate verifier
- **Mechanism:** The register uses the same fields, but routing maps Product and Domain, Architecture and Contracts, Build and Environment, and QA to four role groups. Each group has one named answer owner per ambiguity kind, recorded in the target kit. The idea owner still owns business forks.
- **Mechanism:** Code filters candidates by required capability, tier, evidence, and counterpart harness. Jev can choose among the remaining model/harness pairs; no candidate means `not-active`, not a same-harness fallback. A separate QA harness verifies every change from a clean clone against the same target-native rules.
- **Mechanism:** The route and verifier identities, Git heads, and signed-off handoff IDs live in the task record. Harness-specific rule files are only links or generated views of the neutral kit; they are not independent rule sources.
- **Rests on:** `F-0001#1`, `#9`, `#21`; `F-0003#46`, `#52`, `#66`; O-67.
- **Strains:** A combined owner can approve its own ambiguous contract unless a separate asking role accepts the answer.
- **Falsifiable:** A reviewer on a clean clone cannot identify the contract owner or run the same gates without the first harness.
- **Cost:** Fewer launches than 3A, but more load on each role and on route validation.

## 4. Records in Git, including forge input [B3, B6, B7; C 3.4]

### Option 4A — One protected project event ledger
- **Mechanism:** A protected target ref such as `refs/heads/layup-records` holds append-only `runs/<task>/events.tsv` and typed Markdown payloads. Rows have event ID, kind, requirement ID, actor/role, source URL or comment ID, source content hash, UTC event time, Git head, payload path, and predecessor ID. Only a trusted record writer pushes this ref; a second harness fetches it by its documented name without LAYUP.
- **Mechanism:** Handoff payloads name from/to role, required artifacts and their hashes, preconditions, and receipt. A validator rejects empty, out-of-order, unmet, or missing handoffs. `acceptance.tsv` holds requirement ID, delivery SHA, first submission time, idea-owner accept/reject, review sequence, comment provenance, and decision time. Rejection sends the requirement to another Implement cycle; only acceptance can close it.
- **Mechanism:** `questions.tsv` holds task, kind, assigned owner, asked time, severity, expected response, answer time, answer source, accepted time, and later reversal link. `telemetry.tsv` has action/attempt, task, requirement ID(s), harness/model, input/output tokens, latency, duration, price per token with source and effective date, amount, and `not reported` where a harness gives no count. A shared action records a cost share per requirement; sums cannot count it twice. Incomplete values fail the completeness check; cost is not silently zero.
- **Mechanism:** The GitHub App reads issue comments and PR reviews by authenticated event ID, verifies the author and current body, copies safe content and the source hash into the ledger, pushes, then releases the waiting transition. Replay by event ID is idempotent. An edit or deletion gets a new correction event and pauses dependent work. A comment cannot authorize work while it is only on GitHub or only in a local checkout.
- **Mechanism:** Git is canonical; GitHub is the input and review surface. The runner pins checks to a PR head and ledger commit so a new ledger answer cannot make a stale check green. If a comment contains a secret, stop and ask for a safe replacement; do not copy or summarize the secret into Git.
- **Rests on:** `F-0001#1`, `#12`, `#31`, `#38`; `F-0003#45`, `#50`, `#60`, `#68`–`#75`; O-73.
- **Strains:** A single ref needs a trusted writer and serialized pushes. Without another authorized account or App, an agent that shares the Operator's credentials can forge an approval.
- **Falsifiable:** Disconnect GitHub after a decision was mirrored; a clean clone must show the decision, its approver and source event, and a fresh harness must continue without the old session.
- **Cost:** Git commits and one authenticated relay; no database. More pushes and protected-ref administration.

### Option 4B — Per-task protected record refs
- **Mechanism:** Each task has a protected `refs/heads/records/<task>` and a `docs/tasks/record-refs.tsv` index in the target. Its append-only `handoffs.tsv`, `questions.tsv`, `acceptance.tsv`, `telemetry.tsv`, and `decisions.tsv` use the fields of 4A. The task ref owns its events; milestone acceptance and totals are derived from committed task refs and copied into the next milestone plan.
- **Mechanism:** A GitHub event relay validates actor, comment/review ID, content hash, and target task. It appends the source body as a safe Markdown payload and pushes the task ref before it acknowledges the answer. Corrections are new events; a replacement review does not silently change an earlier decision. The task branch or PR cites the exact task-ref SHA used for its result.
- **Mechanism:** Agents can push code branches but cannot write record refs. A clean clone fetches the refs listed in the index; a validation check rejects a referenced ref or event that is missing, unpushed, stale, or from an untrusted writer. A new task ref is registered before its first role starts; the register is also updated under a trusted writer.
- **Rests on:** `F-0001#1`, `#2`, `#12`; `F-0003#45`, `#50`, `#52`, `#60`; O-73.
- **Strains:** `F-0001#2`: a default-branch-only clone misses the task refs unless the index and fetch rule are explicit; many refs make milestone aggregation harder.
- **Falsifiable:** Two tasks receive comments at once; a missing or misassigned ref must block rather than lose an answer or attach it to the wrong task.
- **Cost:** More refs, permissions, and fetches than 4A; less single-writer contention, with the same relay and GitHub API load.

## 5. Stack gates and rule protection [A1, A2, B9]

### Option 5A — Stack modules in the adapted target kit
- **Mechanism:** A trusted gate owner, not a governed coding agent, writes target-native `docs/gates/stack.tsv`, check scripts, tests with known-bad fixtures, and a CI workflow as part of the target's adapted kit. Each row states stack, gate kind (layout, boundary, contract, test quality), command, expected result, rule source, and evidence that the check runs. The target stack and evidence select the rows; the owner tests at least two stacks.
- **Mechanism:** Target CI runs baseline plus selected stack gates with no LAYUP binary or service. `layup gate TARGET` runs those *same target commands* from outside, reports pass/fail/not-active, and adds no second gate implementation. This supersedes only ADR-0011 decision 4's outside-only stack-gate execution; it keeps the outside pilot run in ADR-0012 part 6 and the no-LAYUP-code rule in ADR-0011 decisions 3 and 4.
- **Mechanism:** A trusted GitHub App runner checks the target PR head against protected baseline gate definitions in an isolated, unprivileged checkout; it publishes a required check for that exact head. Agent edits to rule/gate paths fail a separate trusted path check. If CI does not run, `not-active` blocks merge. A later PR must be checked with the then-current trusted gates, not gates it proposes to weaken.
- **Mechanism:** At the planned post-milestone retrospective, a distinct gate owner writes and approves one batch of rule-path changes with a separate trusted identity. An out-of-window rule change is unplanned input even if somebody approves a PR. A routine PR approval is never reclassified as a planned approval. With today's shared account (ADR-0011 O-9), identity separation is not yet proved; until it is, this gate is `not-active`.
- **Rests on:** `F-0001#2`, `#3`, `#5`, `#7`, `#25`, `#28`; `F-0003#43`, `#44`, `#47`, `#64`–`#67`; O-69, O-67, O-10, O-11.
- **Strains:** Native check scripts in a target can be changed through a PR unless the separate trusted rule-path control really blocks that change; shared credentials defeat a name-only check.
- **Falsifiable:** Remove LAYUP, inject one known-bad change per gate kind in each pilot stack, and try to merge it; any passing or inactive check, or a successful agent rule edit, disproves the option.
- **Cost:** Target-local gate code and two-stack fixture maintenance; GitHub App checks, CI minutes, and a distinct trusted gate owner.

### Option 5B — Stack-native tools and checked-in policy files
- **Mechanism:** The trusted gate owner adapts the kit with configuration and tests native to the target stack: package/module layout checks, dependency boundary checks, interface-contract tests, and test-quality checks. `docs/gates/stack.tsv` names each exact target command, its source and evidence, and a fixture with a violation. The target owns all definitions; LAYUP only discovers and invokes them from outside.
- **Mechanism:** Target CI itself runs every selected command. A separate trusted runner checks at the PR SHA that the protected base workflow and rule files supplied the commands, runs the same commands without privileged credentials, and publishes a required GitHub App verdict. A missing command or unsupported stack is `not-active`, not a reduced gate set.
- **Mechanism:** An isolated trusted identity guards workflow/config/test-rule paths and writes the single retrospective batch. A proposed gate change outside that point counts as unplanned input. The test that removes LAYUP and runs target CI remains part of the pilot; the outside `layup gate` run remains the ADR-0012 exit condition.
- **Rests on:** `F-0001#2`, `#3`, `#5`, `#7`, `#28`; `F-0003#44`, `#47`, `#65`, `#67`; O-69, O-10, O-11.
- **Strains:** A stack's native tools may not express one of the four gate kinds; claiming tool presence is a pass would break `F-0001#5`. A test rule embedded in product tests can be changed by the governed agent unless protected.
- **Falsifiable:** A seeded boundary or test-quality defect passes target CI with LAYUP absent, or an agent edits a protected test rule and the required check stays green.
- **Cost:** Less custom check code than 5A if tools already fit; separate tool setup, two-stack proof, trusted CI and protected-path administration.

## 6. Early escalation and cost stop [B7, B8]

### Option 6A — Live metered sessions
- **Mechanism:** Before a role starts, a decision record names the proposed action, requirement, approved budget and tolerance band, and four business-impact axes. Code stops on a known fork; Jev judges only uncertain semantics. Outside the band or on a legal/intent/strategic fork, the idea owner answers on the issue, and the answer is mirrored into Git before work starts. Inside the band Jev may select an action, but cannot set funding.
- **Mechanism:** The adapter reserves a maximum spend for the session, meters reported usage during it, and stops the process before the available allowance is exhausted. Code uses the price record and measured tokens; a model does not do arithmetic. If a harness cannot report usage soon enough or enforce a spend ceiling, that harness is not eligible for autonomous paid work.
- **Mechanism:** Before each remote model call, the adapter checks the remaining allowance against an enforceable call cap. The record separates reserved amount, reported spend, and unused amount, so a later reader can audit why the next action did or did not start.
- **Rests on:** `F-0001#4`, `#13`, `#26`; `F-0003#48`, `#50`; O-68, O-71.
- **Strains:** Provider token reports can arrive after the spend; a local process kill cannot revoke a remote call already billed.
- **Falsifiable:** A seeded retry loop spends beyond its allowance before the adapter stops it, or an outside-band action starts without the idea owner's recorded answer.
- **Cost:** Continuous metering, adapter work, and possibly lost partial sessions; fewer expensive retries.

### Option 6B — Prepaid bounded actions
- **Mechanism:** Divide a task into short actions. Before each launch, the runner reserves the action's *enforced maximum* charge against the remaining budget and band in `budget.tsv`; a provider-side cap or a verified harness cap must make that maximum real. Reconcile actual usage in `telemetry.tsv` after each action; no next action starts if the next reservation would cross the limit.
- **Mechanism:** The same pre-work four-axis escalation check runs before reservation. Jev can act only inside the idea owner's pre-written band; failures, uncertain axes and outside-band reservations stop for a human. A reservation with no enforceable bound is refused, never recorded as free or safe.
- **Mechanism:** `budget.tsv` records the approving comment ID, price evidence, reservation ID, cap evidence, and result. A new action checks the committed balance and holds a lease so concurrent workers cannot reserve the same funds twice.
- **Rests on:** `F-0001#1`, `#13`, `#26`; `F-0003#22`, `#48`, `#50`; O-68, O-71.
- **Strains:** Small actions cause more starts and more unfinished work; some harnesses cannot offer a hard cap.
- **Falsifiable:** Start a paid action with no proved cap, or let an action spend past its reserved maximum without a stop and a recorded budget breach.
- **Cost:** Extra process launches and per-action records; less need for a live token stream than 6A, but a cap-capable provider is required.

## 7. Learning at the retrospective [C 2.3, 3.3]

### Option 7A — Audited score table
- **Mechanism:** `runs/<milestone>/reward.tsv` joins task/requirement, accepted result, gate first-pass result, counterpart verification, reversal, stall, cost, and duration. A documented, pre-registered formula scores a route only on observed results; missing usage is missing evidence, not zero cost. A route that breaks a hard invariant is ineligible regardless of score.
- **Mechanism:** At the retrospective, a role proposes one `docs/routing.tsv` update from observed route outcomes, with a before/after comparison and sample count. The Operator approves the planned batch of rule-path changes there; code applies the approved table for the next milestone only. No weights or model parameters change.
- **Mechanism:** `runs/<milestone>/lessons.md` links each proposed change to its source events and outcome. A check binds the approved retrospective comment to the changed routing-table hash before the new milestone begins.
- **Rests on:** `F-0001#1`, `#3`, `#12`; `F-0003#39`, `#50`, `#52`; O-67, O-69.
- **Strains:** A small or biased sample can promote a weak route; a weighted sum can hide a gate failure unless the invariant veto comes first.
- **Falsifiable:** A failing gate raises a route's eligibility, or a routing change affects a current milestone before its retrospective approval.
- **Cost:** Little infrastructure; audit and human batch approval each milestone.

### Option 7B — Bounded route trials
- **Mechanism:** `routing.tsv` records safe eligible harness/model pairs and per-pair counts of accepted, rejected, stalled, and costly tasks. A deterministic bounded trial rule chooses between eligible pairs at the next milestone; Jev may label semantic complexity first, but code enforces the trial limit and counterpart-harness requirement.
- **Mechanism:** The retrospective records rewards and the candidates' observed outcomes, then updates only future route weights or counts and any approved rule batch. No weight is changed mid-milestone and no LLM is trained. Human acceptance and later reversals correct the reward record by new Git events.
- **Mechanism:** A recorded seed, eligible set, rule version, and input count make each trial route replayable. The check keeps a counterpart harness available before it admits a pair, even when that reduces the trial set to one.
- **Rests on:** `F-0001#1`, `#9`, `#12`; `F-0003#39`, `#66`, `#72`; O-67, O-69.
- **Strains:** Trials can spend more money or send a task to a worse eligible route; eligibility and budget caps must apply before exploration.
- **Falsifiable:** The route selector tries an ineligible harness, changes a route before retrospective approval, or cannot replay its choice from the recorded counts.
- **Cost:** More comparison sessions and reward data than 7A; no training service.

## 8. Specification synthesis [B4; C 2.1]

### Option 8A — Draft first, check each trace
- **Mechanism:** Design roles read the approved PSB and vision, draft a preliminary design review (PDR) and a PRD with requirement IDs and acceptance criteria, then draft a technical specification in Architect and a milestone/task plan in Plan. A generator role writes the prose; Jev does not. Each document cites a source file hash and exact PSB byte span or quote for each requirement; a task maps to its requirement ID.
- **Mechanism:** `layup spec check TARGET` checks source hashes, exact spans, unique IDs, required columns, acceptance criteria, task coverage, and missing/extra trace links without a model. A separate reasoning role checks whether each paraphrase and technical specification truly agrees with the cited text; planned architecture approval remains with the named human.
- **Mechanism:** The check saves the source hash, the failing requirement ID, and the span in a Git record. The next phase cannot use a draft that fails, even if the prose looks complete to a reader.
- **Rests on:** `F-0001#4`, `#6`, `#31`, `#39`; `F-0003#51`, `#62`; O-67, O-73.
- **Strains:** A byte-exact quote can support a wrong interpretation; a syntactically complete trace is not semantic proof.
- **Falsifiable:** A delivered requirement with no real PSB span or no acceptance criterion passes the check; or a pilot review finds a false trace that no reviewer noticed.
- **Cost:** One generation pass plus deterministic validation and semantic review.

### Option 8B — Source-span ledger before drafting
- **Mechanism:** Design first creates `docs/spec/source-spans.tsv`: PSB hash, byte start/end, quote, proposed requirement ID, and disposition for each in-scope need. Code checks that every cited span exists and that each in-scope item has an explicit disposition before any PRD draft. A generative role then writes a preliminary design review and the PRD from that ledger, the Architect role writes the specification, and Plan links milestone exit tests to requirement IDs.
- **Mechanism:** A mechanical trace check follows links in both directions: PSB span to requirement, requirement to acceptance criterion, specification to requirement, and task to specification. A different role judges semantic agreement and records disagreements; an unapproved new intent goes to the idea owner, never to a text generator or Jev.
- **Mechanism:** The approved ledger and its source hash are pinned at Design exit. Architect and Plan read that version; a new source version reopens the trace check before any phase transition.
- **Rests on:** `F-0001#6`, `#10`, `#31`, `#39`; `F-0003#41`, `#51`, `#62`; O-67.
- **Strains:** The span ledger can give false confidence about complete *meaning* if one long sentence carries several needs.
- **Falsifiable:** A known in-scope PSB need has no ledger disposition, yet Design advances, or a task loses its requirement link without failing Plan.
- **Cost:** One extra ledger and review before prose; less rework from missing links later.

## Forge limit

GitHub-only for the pilot is an accepted, stated limit in O-67 to O-71's target design. I give no competing forge option. In both record options, a GitHub comment or review is input, not the sole lasting decision. The target's native gates still work when LAYUP and its App are absent; a fresh harness can fetch the Git records and continue.

## Questions for the Operator

1. **Trusted identity for rule paths and approvals (ADR-0011 O-9 remains open).** Which separate authority can block an agent that currently uses the Operator's GitHub account: (a) a separately credentialed human gate owner plus protected rule paths and CODEOWNERS; or (b) a separately credentialed trusted App that writes rule changes only after a separately authenticated human approves the retrospective batch? If neither can be made separate, do not claim `F-0001#3` passes.
2. **Evidence for stall limits (PRD §11 question 5 remains open).** Should the initial retry count and no-progress clock limit be (a) set from measured pilot harness/task runs before autonomous use, with their source recorded; or (b) set by the Operator before the pilot as explicit safety ceilings, then calibrated from pilot evidence? Both keep a fixed cap during each milestone; neither takes a number from the closed PR #69.
3. **Enforceable cost ceilings.** Which capability is available for each proposed pilot harness: (a) a live usage stream and an enforceable stop before a provider call, or (b) an enforced provider/harness maximum for each bounded action? If neither is available, autonomous paid work must wait rather than claim a budget stop that only reports cost afterward.
--- end
