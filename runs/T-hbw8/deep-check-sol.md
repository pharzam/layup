# T-hbw8 — deep check by GPT-6 Sol

Reviewer: GPT-6 Sol (`gpt-6-sol-high`) on Devin, a fresh session, 2026-09-28T07:24:33Z to 07:31:42Z. Head `6368b57`.

The output follows word for word inside a fenced block (so that its links, which are relative to the reviewed files, are not read as links of this file); each U+2039 and U+203A character is written by name.

~~~~text
# Independent deep check — `T-hbw8`, head `6368b57`

Scope: proposed design, not a test of unbuilt commands. The fixed rules in `internal/psb/check.go` are correctly limited in `docs/architecture.md:64` and ADR-0020: they do not judge meaning. The findings below test the *other* claims by the same standard. A stated limit does not satisfy a PSB pass/fail criterion.

## Material findings

1. **Rule protection is detection, not protection.**
   - Where: `docs/adr/0017-put-native-stack-gates-in-the-target.md:80-84`: “LAYUP's merge rule and the audit carry Invariant 3 (`F-0001#3`), with the limit that [ADR-0016](0016-keep-the-run-records-in-git-and-tell-agent-from-human.md) states.”
   - Source: `F-0001#3`: “The agents that do the work cannot change the rules or the gates that check the work.” `F-0003#64` requires “0 agent writes to rule paths.” ADR-0016:102-104 admits GitHub cannot prevent a merge by the owner. The agent's user token uses that owner's rights.
   - Why it matters: A role agent can use the owner's GitHub permissions to merge a changed gate. A later audit reports the breach but cannot undo the invalid acceptance. REQ-003 cannot pass.
   - Fix: Do not claim Invariant 3 is met by an audit. Require a preventive control that the work token cannot bypass, or record this as an unsatisfied invariant and test the limit.
   - Severity: `material`.

2. **A pull-request gate cannot prevent human review.**
   - Where: `docs/architecture.md:199`: “the target's CI gates required before the merge; the orchestrator merges only a verified change”. `docs/adr/0017-put-native-stack-gates-in-the-target.md:29-35`: “run by the target's own CI job on every pull request, and required by the target's branch protection.”
   - Source: `F-0003#58`: “No PR with a failed gate reaches human review or merges.” GitHub's required checks control merging, not when a person reads or reviews an open PR.
   - Why it matters: A pilot reviewer opens a PR while CI is pending or failing. The merge is blocked, but the required zero failed-gate PRs reaching human review is already false.
   - Fix: Run gates before opening a PR; request human review only after its CI passes. State the limit for a human who reviews without a request.
   - Severity: `material`.

3. **The stated single writer conflicts with the commands that write records.**
   - Where: `docs/adr/0016-keep-the-run-records-in-git-and-tell-agent-from-human.md:42-43`: “Only `layup run` writes the records branch, with the installation token of the LAYUP GitHub App (its author is the App's bot).”
   - Source: `docs/adr/0017-put-native-stack-gates-in-the-target.md:47-49`: “`layup gate TARGET` runs the same native commands against a checkout of the target and writes the result to the records branch”. `docs/adr/0019-learn-routing-from-the-records-at-each-retrospective.md:30-39` also makes `layup learn` produce `reward.tsv`.
   - Why it matters: A standalone gate command cannot record its result under the only-writer rule; giving it the token makes the rule false.
   - Fix: Make `layup run` commit the outputs of subordinate commands, or define one shared writer with explicit ownership and authentication.
   - Severity: `material`.

4. **The dead-man job cannot both recover a dead run and obey the records rule.**
   - Where: `docs/adr/0013-orchestrate-the-lifecycle-with-an-external-layup-run.md:66-70`: “A scheduled job in LAYUP's own repository reads the heartbeat of each active run; a heartbeat older than `2H` opens a stall record for the orchestrator itself and notifies the Operator with a stall package, not by chance.”
   - Source: ADR-0016:42-43: “Only `layup run` writes the records branch”; ADR-0013:81-83: “The only process of LAYUP is `layup run` while it runs, plus the scheduled dead-man job.” No active-target register or independent records-writer authority is defined.
   - Why it matters: When the host dies, `layup run` cannot write a diagnosis. The scheduled job has no list of targets to inspect and no authorized path to store or send its stall package.
   - Fix: Define the durable active-run index, job credentials, authorized write path, notification path and takeover rule; reconcile them with the single-writer claim.
   - Severity: `material`.

5. **A Git commit's author does not prove who pushed it.**
   - Where: `docs/adr/0016-keep-the-run-records-in-git-and-tell-agent-from-human.md:46-48`: “If the forge cannot give this, the run records say so, and a deterministic audit of the branch history (every commit by the App's bot) is the complement.”
   - Source: The preceding decision requires “Only `layup run` writes the records branch”. A Git commit author is caller-supplied metadata, not the authenticated pusher; the proposed audit checks authors, not push credentials.
   - Why it matters: If the ruleset probe fails, an agent that can push can set its commit author to the App bot. The audit will accept a forged decision or cost record.
   - Fix: Fail closed when the branch cannot be protected, or audit forge-authenticated push events and verify their coverage; do not use commit author as proof.
   - Severity: `material`.

6. **Mandatory shadow mode blocks the Operator's per-target authority choice.**
   - Where: `docs/adr/0014-decide-at-named-points-with-a-pluggable-smart-if-provider.md:64-67`: “The default of every point is `shadow` until a retrospective promotes it with the recorded agreement as its evidence (`F-0001#4`).”
   - Source: O-78: “The authority level (for example shadow or autonomous) is a parameter per target and per decision point.” O-79: “The thresholds, the escalation rules and the check frequency (before each task only, or also at each handoff) are parameters that the Operator sets when running LAYUP against a target, not fixed in the architecture.” O-84 approves these readings. The first retrospective is after a milestone.
   - Why it matters: A pilot Operator who selects delegated decisions at setup still gets a whole milestone of shadow decisions and human calls. The run does not use the approved policy.
   - Fix: Keep shadow as the *default*, but allow an evidenced, explicit per-target override at setup. State what valid authority applies before the first retrospective.
   - Severity: `material`.

7. **Shadow and provider fallback put routine decisions on a human.**
   - Where: `docs/adr/0014-decide-at-named-points-with-a-pluggable-smart-if-provider.md:60-63`: “`shadow` (the provider answers and is recorded, and the deterministic default or a human decides)”. `docs/architecture.md:98-102` gives the Operator as fallback for question routing, stall action and panel disposition, and the idea owner for over-budget choices.
   - Source: `F-0001#24`: “Humans monitor delivery and give input only at the Human Decision Points (§6). No human approves each task.” `F-0003#46` calls for answers without a human where no human decision is needed. O-71 authorizes a human fallback when a provider *cannot decide*, not for every routine shadow result.
   - Why it matters: In the first milestone, each ambiguous routing or panel choice can block on the Operator even though the provider answered. The Operator is again the message bus.
   - Fix: Specify a safe deterministic branch for shadow at each point; count and classify any human call outside a listed Human Decision Point.
   - Severity: `material`.

8. **The escalation screen can find a business fork after the agent makes it.**
   - Where: `docs/adr/0018-escalate-before-the-work-and-stop-cost-at-the-band.md:70-72`: “A business-forking decision is found before the work that depends on it; one that appears inside a role session is found at the next handoff, not before (a stated limit of the default).”
   - Source: `F-0001#13`: “When the escalation rule selects a decision as business-forking, the agent stops and the idea owner decides. An agent never makes a business-forking decision.” `F-0003#57` requires zero missed escalations in the audit sample.
   - Why it matters: A developer chooses a legal or scope change inside a session and codes it. A later handoff detects the choice, but it cannot make the earlier agent decision unmade.
   - Fix: Require sessions to stop and submit a proposed choice *before* they take a business-forking action; use the handoff screen as a second check, not the only check.
   - Severity: `material`.

9. **Question classification is not an accepted answer or a handoff.**
   - Where: `docs/adr/0014-decide-at-named-points-with-a-pluggable-smart-if-provider.md:47-49`: “question routing by ambiguity kind (domain, architecture boundary, interface contract, environment, or none; a table maps the kind to its owner role”.
   - Source: `F-0003#46`: “A question that does not need a human decision gets an accepted answer from the responsible role agent, without a human.” ADR-0016:64-67 lists `asked, answered, accepted` columns but no answer producer, acceptance check or procedure that resumes the blocked role.
   - Why it matters: A pilot question reaches the Domain Expert's label, but nobody starts an answer session, checks its answer or gives it back to the blocked developer. A human routes it by hand.
   - Fix: Define issue creation, owner-session launch, answer record, acceptance authority, timeout and resume handoff; test the full question-to-accepted-answer path.
   - Severity: `material`.

10. **Agent questions do not use the required issue thread.**
   - Where: `docs/architecture.md:245`: “issues carry the human decisions; the forge copy puts them into Git (ADR-0016)”.
   - Source: O-73: “Any commiunication should be through git not agent prompt, answer option”; its recorded reading: “All questions, decisions and approvals go through Git and issue comments.” The vision brief §3.4 says: “All inter-agent handoffs, status reports, human prompts, and code reviews must take place exclusively through structured repository issues and threaded comments.” `docs/issue-workflow.md` R6 requires an agent question as an issue comment.
   - Why it matters: The routed `questions.tsv` row does not notify a peer on the issue or let a fresh session read the required thread. A peer can wait while the row says it was routed.
   - Fix: State the issue-thread protocol for role questions and answers; copy the thread to Git before using it, with a link from each record.
   - Severity: `material`.

11. **Missing token data makes the money stop unable to compute a total.**
   - Where: `docs/adr/0018-escalate-before-the-work-and-stop-cost-at-the-band.md:49-54`: “Before each new action, code sums the prices of `telemetry.tsv` and adds a projection of the remaining work.”
   - Source: ADR-0018:57-62 permits a harness with no token report, and ADR-0016:70-73 says its token value is `not reported`. O-80 explicitly permits this incomplete telemetry with a wall-clock proxy. No price-per-minute source or conversion from elapsed time to `B` and `U` is defined.
   - Why it matters: A paid session reports no tokens. Code cannot sum its cost against the money band. It can continue past `U` while its wall-clock timer has not expired.
   - Fix: Use an evidenced money upper bound or a metered cap for each action; treat unknown spend as unknown and stop or request a decision before the next action.
   - Severity: `material`.

12. **The allowed incomplete telemetry cannot pass the PSB completeness check.**
   - Where: `docs/adr/0015-route-role-sessions-over-registered-harnesses.md:63-66`: “By default a harness that does not report tokens or cannot enforce a spend cap may do paid work with a wall-clock limit as the proxy, and its telemetry counts as incomplete (O-80)”.
   - Source: `F-0003#60`: “Every task has a token count, a latency, and a wall-clock duration in the project repository.” `PRD-0001` REQ-011 acceptance: “tasks with a complete telemetry record divided by all tasks equals 1”. O-80 permits incomplete telemetry; it does not change that PSB success check.
   - Why it matters: The default route uses a paid harness without token reporting. Its pilot task is incomplete by definition, so the required 100% is impossible.
   - Fix: State that such pilot work fails the success check, or ask the idea owner to revise the criterion; use a token-reporting harness for tasks whose pilot result must pass.
   - Severity: `material`.

13. **The reversal reward has no reversal record.**
   - Where: `docs/adr/0019-learn-routing-from-the-records-at-each-retrospective.md:30-39`: “`layup learn` reads only the records branch” and “a requirement accepted at its first review counts for the route, a reversal, a stall that needed the Operator and an unplanned input count against it, and a task with incomplete telemetry counts as missing evidence.”
   - Source: `F-0003#72`: “The audit uses a random sample and a window of 30 days after the merge.” `docs/architecture.md:137-146` lists no audited-answer sample, overturned-answer event or 30-day audit result in the records schema; `questions.tsv` has only `accepted`.
   - Why it matters: An answer accepted today is reversed after a merge. The retrospective cannot find that reversal in Git, so it gives the route a false reward.
   - Fix: Record answer IDs, sampled audit results, reversal events, audit dates and the task/route join; do not compute that reward term without them.
   - Severity: `material`.

14. **The claimed Task Intervention Rate omits human actions.**
   - Where: `docs/architecture.md:200`: “the escalation screen and the human-input accounting (`questions.tsv` human and planned columns)”.
   - Source: `F-0001#28`: “Human input to a task that is not at a Human Decision Point (§6), for example an answer, a correction, a restart, or a change to a gate.” `F-0003#70` measures tasks with at least one such input. A question row does not record a human correction, restart, PR approval or gate change; ADR-0013:54-55 records a restart but gives no join to this accounting.
   - Why it matters: A pilot task takes a human restart and a PR correction but asks no question. It is counted as autonomous, so the intervention rate is too low.
   - Fix: Record every input type with task, actor, time and planned/unplanned reason; compute the rate from that complete event set.
   - Severity: `material`.

15. **Early Question Share has no project-wide question population.**
   - Where: `docs/architecture.md:77`: “the idea owner answers the whole batch once, and the answers are stored as a fact”; `docs/architecture.md:142`: “`runs/<task>/questions.tsv`”.
   - Source: `F-0003#75`: “Questions before delivery, divided by all human questions for the project.” The intake batch occurs before the first task; the given question table is per task.
   - Why it matters: Intake questions have no row in the per-task question table. A pilot cannot count the numerator from the stated data or reproduce the reported share.
   - Fix: Give intake and delivery questions a shared project-wide index, with a time, human flag and delivery-start boundary.
   - Severity: `material`.

16. **The start values and baseline measurements have no run path.**
   - Where: `docs/architecture.md:86-87`: “After a Retrospective, the next milestone starts at Design, until every `Must` requirement is accepted.” The phase loop has no current-process baseline or approval step for the PSB §7.2 start values.
   - Source: `F-0004#11`: “The first pilot measures the baseline with the current process; then the idea owner sets each start value in one batch, records it in the pilot's PRD, and the pilot's report compares the measured trend with it.” `F-0003#68`–`#75` name eight distinct measures.
   - Why it matters: A pilot can compute its new-process values but has no current-process runs, approved batch or report to compare against. The lead-time and cost targets have no denominator.
   - Fix: Add a baseline phase and its records for all eight measures; require the owner's batch of start values in the pilot PRD before evaluation.
   - Severity: `material`.

17. **A fact-numbering command cannot find semantic clauses in arbitrary prose.**
   - Where: `docs/adr/0020-derive-the-specification-from-numbered-facts.md:34-36`: “`layup` stores the problem statement byte-identical in the target's `docs/facts/`, with a record that numbers each clause as a fact (`F-NNNN#n`), as LAYUP's own facts are stored.”
   - Source: `F-0003#24`: “A problem statement is prose. Nothing converts it into identified requirements that a machine can validate and trace back to the text that states them.” ADR-0020:20-23 says `layup psb check` detects only G1–G5 and the smart-if cannot generate text.
   - Why it matters: A pilot brief puts two needs in one paragraph and has no “In Scope” heading. A deterministic formatter cannot know where the clauses or In-Scope facts are; the draft has no reliable input.
   - Fix: Assign a role session to extract and classify candidate verbatim spans; have the idea owner approve the fact/need inventory, then check its byte-exact references.
   - Severity: `material`.

18. **A trace check does not find needs omitted from the fact inventory.**
   - Where: `docs/adr/0020-derive-the-specification-from-numbered-facts.md:44-49`: “`layup spec check TARGET` fails when a fact ID does not resolve, a quote is not a byte-exact substring of the stored problem statement, an In-Scope fact has no `Must` row, a requirement has no acceptance criterion, a task names no requirement, or a specification section names no requirement.”
   - Source: `F-0003#51`: “Derivation of identified requirements and technical specifications from the approved problem statement.” `F-0003#62` requires a trace for “Every delivered requirement”. The listed checks are over facts already numbered; none checks unnumbered needs in the original text.
   - Why it matters: Intake misses a need. All remaining fact IDs resolve and all rows have quotes, so the trace passes while a product need is never planned or delivered.
   - Fix: Add an independent completeness review of original brief versus fact inventory and requirements; state that the byte-exact check proves links, not full coverage.
   - Severity: `material`.

19. **The generated skeleton assigns `Must` priority without an intent decision.**
   - Where: `docs/adr/0020-derive-the-specification-from-numbered-facts.md:37-43`: “`layup spec draft` writes one requirement row per In-Scope fact, with the fact's words as its text, the fact ID, and an empty acceptance criterion.” ADR-0020:44-48 then requires “an In-Scope fact” to have a “`Must` row”; ADR-0013:37-40 runs until every `Must` is accepted.
   - Source: `F-0001#10`: “Before delivery starts, the idea owner sets the problem, the success criteria, and the funding.” `F-0003#54` leaves decisions that set intent to the idea owner. An In-Scope clause in another product's brief does not itself state MoSCoW priority.
   - Why it matters: One optional feature gets a `Must` row and blocks completion of the pilot product even though the owner never made it mandatory.
   - Fix: Make priority an explicit idea-owner decision at Intake or Design; let the draft leave it unset and block approval until set.
   - Severity: `material`.

20. **The specification check does not require a specification for every requirement.**
   - Where: `docs/adr/0020-derive-the-specification-from-numbered-facts.md:47-49`: “a task names no requirement, or a specification section names no requirement.”
   - Source: `F-0003#51`: “Derivation of identified requirements and technical specifications from the approved problem statement.” `PRD-0001` REQ-012 acceptance requires “each delivered requirement has a technical specification in the target that names the requirement it implements”. The stated check only checks existing sections, not missing ones.
   - Why it matters: A pilot has a traced delivered requirement but no design section for it. `layup spec check` passes; the requirement's technical specification is absent.
   - Fix: Check each delivered requirement for a linked, nonempty specification; keep semantic agreement with the independent verifier.
   - Severity: `material`.

21. **The handoff's `valid` field is not a defined validation rule.**
   - Where: `docs/adr/0016-keep-the-run-records-in-git-and-tell-agent-from-human.md:64-65`: “`handoffs.tsv` (from role, to role, kind, artifact path and hash, condition, validity)”. ADR-0013:40-43 says a phase moves “only on a valid handoff record”.
   - Source: `F-0003#45`: “Information that passes between role agents has a form that a machine can validate, based on Armature conventions.” `F-0003#59` requires “schema-validated state artifacts across all role transitions”. No transition schema, allowed kinds, condition evaluator or check of the actual artifact is stated in these decisions.
   - Why it matters: A role writes `valid=yes` for a nonexistent artifact or an unmet condition. The phase can advance if code trusts the column; a header-only check cannot reject it.
   - Fix: Define the schema and transition/condition checks from cited baseline rules; compute validity from the artifact and state, never from an agent's assertion.
   - Severity: `material`.

22. **A new record can reset the stall clock without real progress.**
   - Where: `docs/adr/0013-orchestrate-the-lifecycle-with-an-external-layup-run.md:56-59`: “A progress event is a new durable record: a new artifact hash, a passed gate run, a handoff row. An equal failure signature and no new pushed head are not progress.”
   - Source: `F-0003#49`: “a step that repeats without progress, gets a procedure with a limit.” New hashes and new handoff rows do not prove that the repeated step came closer to its goal.
   - Why it matters: A retry loop changes one log file and writes a new handoff each time. `T` resets forever and `N` does not trigger if each attempt has a new signature.
   - Fix: Tie progress to an exit condition or a bounded task state, and put an absolute attempt/time cap on repeated work even when artifacts change.
   - Severity: `material`.

23. **The architecture claims model allocation by complexity but has no complexity input.**
   - Where: `docs/architecture.md:239`: “2.2 squads, model allocation by complexity | the routing table and the fit score among admitted pairs (ADR-0015)”.
   - Source: The vision brief §2.2: “allocating frontier and high-reasoning models based on task complexity.” ADR-0015:43-55 lists `role, task class, author harness and model, verifier harness and model, weight ...`; it defines no task-complexity class, assessment or tier constraint.
   - Why it matters: A high-risk architectural task and a simple edit can receive the same cheap pair because admission and fit have no complexity value to use.
   - Fix: Define a task-complexity assessment with evidence and a routing policy that uses it before fit scoring.
   - Severity: `material`.

24. **The panel claim lacks the panel execution path.**
   - Where: `docs/architecture.md:243`: “a panel as a stall action and in Design (ADR-0006); members see no other member's output; the examiner gets the evidence, not the performers' reasoning”.
   - Source: Vision brief §3.2: “The system must spawn a fresh, context-isolated panel of high-reasoning "Blind Agents"” and “The panel delivers a synthesized, evidence-backed resolution path”. ADR-0014:50-52 only assigns synthesis to “a reasoning session”; ADR-0013:44-49 starts “one harness session” for each step. No panel fan-out, sealed inputs or independent synthesis handoff is defined.
   - Why it matters: At a pilot deadlock, the stall action selects “panel”, but no process starts isolated members and collects their options. The next role has no resolution to act on.
   - Fix: Define panel convening, blind inputs, member isolation, option records, synthesis actor and exit rule; do not equate the option label with a working panel.
   - Severity: `material`.

25. **The proposed learning result does not reach the next project.**
   - Where: `docs/adr/0019-learn-routing-from-the-records-at-each-retrospective.md:66-68`: “The lessons of one target reach the next project only through LAYUP's own default parameters and recipes, which change under LAYUP's gate; that path is work for the implementation plan.”
   - Source: `F-0003#39`: “the next project gets no lesson from it.” Vision brief §2.3: “committing actionable lessons learned to iteratively refine subsequent execution cycles”; O-67 keeps the learning loop in scope. The proposed mechanism changes only the *current target's* `routing.tsv`.
   - Why it matters: A second pilot target starts with unchanged LAYUP defaults and repeats the first target's costly route, despite an approved lesson.
   - Fix: Define who reviews, promotes and versions a cross-target lesson in LAYUP, and when a new target uses it; keep each source and approval.
   - Severity: `material`.

26. **ADR-0011 still states the gate location that O-76 rejected.**
   - Where: `docs/adr/0011-structure-the-core-engine-as-a-go-cli-over-repository-files.md:74-77`: “LAYUP's stack-dependent gates and its setup verification are `layup` commands that run **outside** the target, against a checkout of it”.
   - Source: O-76: “Anything in an ADR that contradicts this is removed.” ADR-0017:64-69 explicitly amends decision 4: “the stack gates are the target's own and run inside it”. The old operative decision and consequences at ADR-0011:111-113 remain in the tree.
   - Why it matters: An implementer follows accepted ADR-0011 and builds stack enforcement only into LAYUP; a target then fails the independence test without LAYUP.
   - Fix: Replace the superseded operative clauses in ADR-0011 with a short pointer to ADR-0017; keep the historical reason as context, not an active instruction.
   - Severity: `material`.

27. **The phase loop has no enforceable pull-request creation and review step.**
   - Where: `docs/architecture.md:82`: “Per task: the screen, the role session, the target's gates, the verification on a counterpart harness, the merge by the orchestrator”.
   - Source: `docs/issue-workflow.md` R1: “A pull request with no linked issue does not merge.” `docs/engineering-discipline.md:243-245`: “The change then lands through a pull request whose body links that issue”. ADR-0017:29-35 runs stack CI “on every pull request”. ADR-0016:74-77 names the merger, but neither decision assigns issue/PR creation, links, verified head, reviewer verdict or how the orchestrator uses them.
   - Why it matters: A task agent pushes a branch and a verifier passes it; `layup run` has no specified PR to open or to merge and no required review record to inspect. The CI trigger may never run.
   - Fix: Add issue/PR lifecycle, agent author, independent review record, head binding and merge preconditions to Implement; reject a task with no linked PR.
   - Severity: `material`.

28. **Clarification Turnaround has no time for the accepted answer.**
   - Where: `docs/architecture.md:142`: “`runs/<task>/questions.tsv` | events | qid, asked_by, kind, kind_source, owner, asked, answered, accepted, human, planned”.
   - Source: `F-0003#71`: “Time from the question to the accepted answer, for questions that are resolved without a human.” The stated columns contain no asked-at or accepted-at time; a row can say accepted without recording when.
   - Why it matters: A pilot gives a peer answer after hours, but the system cannot calculate its 95th-percentile turnaround or test the 120-second start value.
   - Fix: Add time-stamped question, answer and acceptance events linked by ID; compute time from first ask to accepted answer.
   - Severity: `material`.

29. **A prompt-file list is not minimal context routing.**
   - Where: `docs/architecture.md:242`: “context: the prompt file names the records a session reads (ADR-0013)”.
   - Source: Vision brief §3.1: “Selectively scope, prune, and route minimal necessary context windows per agent invocation.” ADR-0013:44-49 says only that a prompt file names records and a handoff; it states no selector, pruning criterion or context-size check.
   - Why it matters: A pilot role receives a full history that exceeds its model's context, or misses the one relevant prior decision despite a prompt-file name. The claimed routing does not choose or bound context.
   - Fix: Define how a role selects the minimal cited records, preserves required constraints and refuses a session whose context cannot fit.
   - Severity: `material`.

## Notes

30. **A review is not an issue comment.**
   - Where: `docs/adr/0016-keep-the-run-records-in-git-and-tell-agent-from-human.md:58-60`: “Before any action reads a comment or a review, `layup run` copies it into the records branch: the body word for word, the author, the comment ID, the `performed_via_github_app` value, the time and a SHA-256 of the body.”
   - Source: GitHub REST documents `performed_via_github_app` for [issue comments](https://docs.github.com/en/rest/issues/comments), but its [pull-request review](https://docs.github.com/en/rest/pulls/reviews) response schema lists a review ID, state, body and submit time, not that field. Review and issue-comment APIs are different.
   - Why it matters: A pilot copies a PR approval from the review API and cannot populate the promised agent/human field or “comment ID” without a defined source.
   - Fix: Specify and probe a documented timeline-event join for reviews, or restrict decisions to issue comments and mark unclassifiable reviews as such.
   - Severity: `note`.

31. **A first-review count needs an explicit review boundary.**
   - Where: `docs/adr/0016-keep-the-run-records-in-git-and-tell-agent-from-human.md:67-69`: “`acceptance.tsv` (requirement, milestone, delivered commit, the idea owner's accept or reject, the comment ID, the review count, the time)”.
   - Source: `F-0003#69`: “Requirements that the idea owner accepts at the first review, divided by all delivered requirements.” A count field by itself does not say which delivery or rejection starts and ends one review.
   - Why it matters: A rejected requirement resubmitted in the same milestone can be marked review 1 twice; First-Review Acceptance is then too high.
   - Fix: Give each delivery and review an immutable ID and an ordered acceptance/rejection event; compute review number from history.
   - Severity: `note`.

32. **The default ambiguity-owner map has no stated evidence.**
   - Where: `docs/adr/0015-route-role-sessions-over-registered-harnesses.md:39-42`: “this map is a default with no fact behind it, and the panel members gave two different maps (`panel-A.md`, `panel-C.md`). The Operator can replace the matrix per target (O-81)”.
   - Source: `F-0001#4`: “No configuration value without evidence. A value that comes from a guess is a defect.” O-81 approves the seven default functions but does not choose an owner for each of the four kinds.
   - Why it matters: An environment question can go to the Software Engineer when the target needs another owner. The owner may answer without the right expertise.
   - Fix: Ask the Operator to accept or edit the owner map per target, with a recorded source; treat an unapproved default as unset.
   - Severity: `note`.

## PSB items checked and answered

All seven problems, twelve In-Scope items, nine invariants, five Human Decision Points, eleven pass/fail measures and eight calibrated measures were checked. Findings above name the gaps. These items have a stated answer:

- `F-0003#41` — Intake combines fixed text rules with a review of meaning before the idea owner's batch.
- `F-0003#42` — setup follows the pinned baseline, records values and runs setup verification.
- `F-0003#44` — stack recipes supply native layout, boundary, contract and test gates in the target.
- `F-0003#52` — neutral rules, a named records branch and a second harness support replacement.
- `F-0001#5` — a selected gate that did not run is `not-active`, not passed.
- `F-0001#6` — mechanical gates take priority over model judgement.
- `F-0001#7` — stack recipes add gates rather than remove baseline rules.
- `F-0001#8` — Armature has a recorded pinned version.
- `F-0001#11` — the idea owner answers the Intake questions in one batch.
- `F-0001#12` — Intake records approvers and Accept records a decision per delivered requirement.
- `F-0003#61` — a live stall gets a diagnosis and outcome record.
- `F-0003#65` — native gates run on a target without LAYUP.
- `F-0003#66` — a distinct harness must verify each change; fewer than two blocks merge.
- `F-0003#67` — two problem statements on different stacks are in the pilot plan.
~~~~
