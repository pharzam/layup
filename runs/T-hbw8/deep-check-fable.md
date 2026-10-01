# T-hbw8 — deep check by Claude Fable 5.1

Reviewer: Claude Fable 5.1 (`claude-fable-5-1`) on Claude Code (`claude -p`), a fresh session, 2026-09-28T07:24:33Z to 07:42:23Z. Head `6368b57`.

The output follows word for word inside a fenced block (so that its links, which are relative to the reviewed files, are not read as links of this file); each U+2039 and U+203A character is written by name.

~~~~text
# Deep check of the LAYUP architecture — T-hbw8 at 6368b57

Reviewer: Claude Fable 5.1 (not the author's model). Read: `docs/architecture.md`, ADR-0011, ADR-0013 to ADR-0020, the PSB and `F-0001`/`F-0003`/`F-0004`, the vision brief, O-66 to O-84, `PRD-0001`, `jev-sources.md`, `internal/psb/check.go`, the kit CI in the target copy. GitHub facts come from docs.github.com (fetched 2026-09-28).

**Verdict: not mergeable, findings recorded.** 23 material findings, 15 notes. The defect class of `layup psb check` is present in more places: a mechanism that cannot see what the text says it checks (M3, M4, M7, M14, M21), and PSB items with a table row but no mechanism (M14, M15, M16, M18).

## Material findings

### M1. One App gives the agents every right of the orchestrator — Severity: material
- **Where:** architecture §2: "Owned by the Operator's account: its user access token for the agents (the badge), its installation token for the orchestrator (the only writer of the records branch)".
- **Source:** GitHub: with a user access token "The app can only access resources that the user has access to" and "that it has permission to access". The Operator owns the target, so the agent token has all App permissions. That is the same set that the installation token has. Invariant 3 (`F-0001#3`).
- **Why it matters:** An agent session calls the merge API on its own pull request. If the App has `administration` (to create the ruleset at setup), the agent can also delete the ruleset. GitHub accepts both calls. Invariant 3 then rests on the prompt.
- **Fix:** Use two Apps under the Operator's account. The agent App gets contents, pull requests and issues, but no `administration` and no merge path. The orchestrator App is the only bypass actor.

### M2. §10 says GitHub cannot prevent what ADR-0016 prevents — Severity: material
- **Where:** §10: "GitHub cannot stop an actor with the Operator's rights from merging or from changing a rule path on a repository owned by one user; the merge rule and `layup audit` detect it, they do not prevent it." ADR-0016 decision 2: "A ruleset on the branch blocks every other update and lists only that App as a bypass actor … a push with the agents' token must fail."
- **Source:** The probe of ADR-0016 itself assumes that a ruleset refuses the Operator-account token. The same ruleset on the default branch refuses an agent merge. Invariant 3 says "cannot change"; REQ-003 criterion: "an attempted rule change by an agent without the control is refused".
- **Why it matters:** The architecture picks detection where it already uses prevention. REQ-003 fails in the pilot, because nothing refuses the change.
- **Fix:** Put the default branch under a ruleset whose only bypass is the orchestrator App (with M1). Keep the audit as the complement. Say how the Operator lifts the ruleset when LAYUP is absent (Invariant 2).

### M3. The App marker does not identify reviews, commits or the records writer — Severity: material
- **Where:** ADR-0016 decision 3: "A human decision is a comment or review by the Operator's or the idea owner's account … with `performed_via_github_app` empty." Decision 2: "a deterministic audit of the branch history (every commit by the App's bot) is the complement."
- **Source:** The REST schemas of a pull-request review and of a review comment have no `performed_via_github_app` field. Issue comments have it. A Git commit has only an author and a committer, and the client writes that text.
- **Why it matters:** An agent review through the App token has no field, so it reads as a human decision. An agent can set its commit author to the App's bot, so the fallback audit passes on a forged record. "0 agent writes to rule paths" (`F-0003#64`) cannot be counted from commits.
- **Fix:** Accept a human decision only as an issue comment. Make records commits through the API as the App, and audit the GitHub `verification` of each commit, not the author text.

### M4. A disagreement between two roles is never a stall — Severity: material
- **Where:** ADR-0013 decision 5: "A progress event is a new durable record: a new artifact hash, a passed gate run, a handoff row."
- **Source:** `F-0003#17`: "Two role agents give different answers, and neither one gives way." `F-0001#37` names this as a stall.
- **Why it matters:** An author and a verifier loop: each round pushes a new head and writes a handoff row. Each round is "progress", so the clock never opens a stall. The loop spends the budget, and the Stall Rate reads 0.
- **Fix:** Count the rounds between the same two roles on the same finding or requirement, with a limit parameter. A repeated rejection is not progress.

### M5. The clock opens a stall at each wait for a human — Severity: material
- **Where:** ADR-0013 decision 5: "When no progress event comes for `T` minutes … `layup run` opens a stall record"; §8: `stall.T` = 10 minutes.
- **Source:** Intake answers, the architecture approval, each escalation and the Accept phase wait for a person. `F-0001#37`: a stall is disagreement or a step that repeats, not a wait.
- **Why it matters:** Each planned wait longer than 10 minutes becomes a stall and a stall package to the Operator. The Stall Rate and the Operator's load go up with no real stall.
- **Fix:** Stop the clock while a step waits on a Human Decision Point. Record the wait as a phase event.

### M6. The stall path can reach the Operator with no diagnosis — Severity: material
- **Where:** ADR-0013 decision 5: "The next action is one of: a retry …, a fresh examiner on a model and a harness that did not do the stalled work, a panel, or the Operator with a stall package". Decision 6: the dead-man job "notifies the Operator with a stall package".
- **Source:** `F-0001#14`: "the agent stops, packages the evidence and the diagnosis, and sends them to the Operator"; `F-0003#61`: "Every stall has a recorded diagnosis". The PSB asks for "a fresh context" (`F-0003#49`), not a third harness.
- **Why it matters:** The choice "the Operator" skips the examiner, and the dead-man path has no examiner. With two harnesses, a stall between the author harness and the verifier harness leaves no eligible examiner. Stall Diagnosis fails.
- **Fix:** Make the examiner's diagnosis a step before the Operator branch. Require a fresh session, not a third harness. Give the dead-man stall an examiner at the next start.

### M7. The escalation screen reads the agent's own list, and shadow mode is not defined — Severity: material
- **Where:** ADR-0018 decision 1: "at each handoff on the handoff's list of the choices the role made"; the floor: "a number of the budget file, a line of the approved problem statement, or a `Must` requirement". ADR-0014 decision 5: "`shadow` (… the deterministic default or a human decides)".
- **Source:** `F-0001#13`: "An agent never makes a business-forking decision." `F-0003#57`: Missed Escalations = 0, a Layer 1 check that a pilot does not calibrate. Target design: the screen "replaces the agent's own declaration".
- **Why it matters:** An agent that leaves a scope-against-date choice off its list is never screened. The floor matches only file edits that agents cannot make. In milestone 1 every point is `shadow`: either each screened choice goes to the idea owner (each unconfirmed one is unplanned input, `F-0001#28`) or none does (a missed escalation).
- **Fix:** Screen the diff and the handoff artifact, not only the list. Define the shadow default of each point (escalation: escalate), and state its cost in the Task Intervention Rate.

### M8. Question routing has no "needs a human" branch and trusts the asker's label — Severity: material
- **Where:** §5: "Question routing | a question that names its kind | one choice: domain, architecture boundary, interface contract, environment, none".
- **Source:** `F-0003#9`: "An agent cannot tell which kind of ambiguity it has … or if the question needs a human decision."
- **Why it matters:** The "deterministic check" accepts the kind that the asking agent wrote, which the PSB says it cannot know. A scope question with the label "domain" goes to the Domain Expert agent, and that agent makes a business-forking decision.
- **Fix:** Give the asker's label to the provider as input. Put each question and each answer through the four escalation questions. Add a branch to the idea owner.

### M9. The target's own Armature CI rejects every agent pull request — Severity: material
- **Where:** §4 Implement: "Per task: the screen, the role session, the target's gates, the verification on a counterpart harness, the merge by the orchestrator".
- **Source:** O-67: "should strictly follow Armature principles"; target design: "Implement (the Armature gate per task …)". The target keeps the kit workflows (`steps.tsv` S12) and requires "one required check per job" (S13). `review-record-lint.sh` "passes only when the record a task leaves behind can be parsed"; `pr-link-lint.sh` needs a linked issue.
- **Why it matters:** No actor opens the task issue, writes the plan and its review, or records the review rounds. Each agent pull request fails a required check, so `layup run` never merges. Or the Operator removes the jobs, which is "a change to a gate" (`F-0001#28`).
- **Fix:** Map each step of the target's quality gate to a phase step and an actor (issue, plan, plan review, rounds, close-out). Or record which kit jobs a target keeps, and why.

### M10. Intake needs a target, records and harnesses that exist only after Scaffold — Severity: material
- **Where:** §4 Intake: "a review session on a harness … a session on a counterpart harness …; `approvers.tsv`, `budget.tsv` and the list of planned approval points are written". Scaffold, the next phase: "the records branch, the App installed". ADR-0015: "A harness enters the register only after a recorded probe run."
- **Source:** Invariant 1: "No project state and no decision is kept only outside the project repository." O-73: a person writes an issue comment.
- **Why it matters:** Decision Points 1 and 2 have no repository, no issue to answer on and no records branch. The Intake sessions cannot be admitted, because the harness register does not exist.
- **Fix:** Create the repository, install the App and probe the harnesses before Intake. Or name where Intake records live and when `layup setup` copies them in.

### M11. The budget file asks for milestones and tasks before they exist, from a writer who cannot write — Severity: material
- **Where:** ADR-0018 decision 2: "At Intake the idea owner writes `budget.tsv` on the records branch: per milestone, the budget `B` and the upper edge of the band `U`, and per task an estimate."
- **Source:** Milestones and tasks come from the Plan phase. ADR-0016 decision 2: "Only `layup run` writes the records branch". O-73: the person writes an issue comment.
- **Why it matters:** The idea owner has no milestone list at Intake, and the ruleset refuses their push. A task estimate that a role writes later "changes a number of the budget file", so the floor escalates each estimate. The cost stop has no base for its projection.
- **Fix:** At Intake, take a total `B` and `U` from an issue comment. In Plan, write the per-milestone and per-task estimates as derived records inside that total.

### M12. "One writer" has six writers — Severity: material
- **Where:** ADR-0016 decision 2: "Only `layup run` writes the records branch". Against: §3 `layup gate` "writes the result to the records branch"; ADR-0019: "The retrospective's role session writes `runs/<milestone>/lessons.tsv`"; `layup learn` writes `reward.tsv`; the dead-man job "opens a stall record"; M11's idea owner.
- **Source:** ADR-0016 against ADR-0013, ADR-0017, ADR-0018, ADR-0019.
- **Why it matters:** The implementation must open the ruleset to these writers (then agents write records), or their pushes fail.
- **Fix:** State that each other producer gives its output to `layup run`, which commits it. Name where the dead-man job's record goes when `layup run` is down.

### M13. Two smart-if calls at points that ADR-0014 forbids — Severity: material
- **Where:** ADR-0019 decision 4: "The smart-if provider may score how far a lesson is tied to a record"; ADR-0020 decision 4: "The smart-if provider may flag a row for the verifier".
- **Source:** ADR-0014 decision 3: "A call at any other point is a defect." §5 lists six points. NFR-005 criterion: "a row whose point is not a named decision point, fails the criterion."
- **Why it matters:** An implementation that follows ADR-0019 or ADR-0020 fails NFR-005.
- **Fix:** Add the two points to ADR-0014 decision 3 and §5, each with its deterministic check and fallback. Or remove them.

### M14. Handoff validation exists only as a column — Severity: material
- **Where:** §9.1: "`F-0003#45` Role Handoffs | `handoffs.tsv` and its check (ADR-0016, ADR-0013)"; §7 column `valid`.
- **Source:** `F-0003#45`: "a form that a machine can validate, based on Armature conventions"; `F-0003#59`: "100% schema-validated state artifacts"; REQ-005: "a schema that accepts everything does not pass". No ADR gives a schema, an allowed-transition table or a check command.
- **Why it matters:** A row with a path and a hash "validates" any artifact. The metric reads 100% on empty content.
- **Fix:** Name the schema of each handoff kind (fields from the kit's task and review records), where the transition table lives in the target, and the command that checks it.

### M15. Three PSB measures have no record and no actor — Severity: material
- **Where:** ADR-0016 Consequences: "the measures of PSB §7.2 can be computed from the records". ADR-0019 decision 2: "a reversal … count[s] against it". ADR-0016 decision 4: "Before any action reads a comment or a review, `layup run` copies it".
- **Source:** `F-0003#72` (overturned answers, "a random sample and a window of 30 days"); `F-0003#57` (an audit sample of agent decisions); `F-0003#70` and `F-0001#28` ("an answer, a correction, a restart, or a change to a gate").
- **Why it matters:** No table has an "overturned" column, and no step draws a sample. `layup audit` checks only "the records branch history, the merges and the rule-path changes". A human comment that no action reads is never copied, and a harness session reads the forge directly with its token. Reversal Rate, Missed Escalations and Task Intervention Rate cannot be computed, and the reward term "reversal" reads nothing.
- **Fix:** Add an audit step (actor, time, sample rule) and a column for overturned answers. Copy every human comment on the target. Count parameter edits and restarts as input.

### M16. No actor for the role of each step, the task class or the prompt — Severity: material
- **Where:** §8: "`roles.matrix` | the seven PSB §2 functions and the owner map of ADR-0015"; §7 `routing.tsv` is keyed by `task_class`; ADR-0013 decision 3: "a prompt file that names the records to read". §9.4: "model allocation by complexity | the routing table and the fit score"; "context: the prompt file names the records a session reads".
- **Source:** Vision 2.2 ("based on task complexity") and 3.1 ("Selectively scope, prune, and route minimal necessary context windows"), kept by O-67.
- **Why it matters:** The default gives kind owners but no phase per role. `layup run` cannot pick the role that writes code, the one that reviews, or the one that writes a synthesis or a diagnosis. Nothing sets a task's class. Nobody writes the prompt file. The two vision items exist in name only.
- **Fix:** Add a default role-per-step table, a rule that sets the task class at Plan, and prompt templates as target content with the records of each role.

### M17. Learning does not reach routing once a point is `delegate` — Severity: material
- **Where:** ADR-0015 decision 4: "code takes the best pair above the threshold and uses the weight to break a tie." ADR-0019 changes only the weights.
- **Source:** `jev-sources.md`: "Jev is not fine-tuned or LoRA-adapted with customer data … the same weights serve every account." Vision 3.3: "tune routing".
- **Why it matters:** After promotion, a score that never sees the pilot's records picks the pair. The learned weight acts only on an exact tie of two probabilities. The learning loop acts only while the point stays `shadow`.
- **Fix:** Let code decide with the weight among the pairs above the score threshold, or combine the score and the weight by a stated formula.

### M18. Blind panels are covered by a sentence that no ADR decides — Severity: material
- **Where:** §9.4: "a panel as a stall action and in Design (ADR-0006); members see no other member's output; the examiner gets the evidence, not the performers' reasoning".
- **Source:** Vision 3.2 (fresh isolated panel, hypothesis posture, synthesis), kept by O-67. ADR-0006 in a target: panel-worthy decisions are named by the project, and "which domains sit on one is the adopter's `<U+2039>…<U+203A>` marker".
- **Why it matters:** No ADR from 0013 to 0020 gives the size, the member domains, the blind rule, the posture, the author of the synthesis or the trigger in Design. `layup run` has nothing to execute for "a panel".
- **Fix:** Add one decision with the panel's composition, blind rule, posture prompt, synthesis author and trigger.

### M19. A wrong gate stops the milestone with no exit — Severity: material
- **Where:** ADR-0017 Consequences: "A wrong gate that blocks all work in a milestone waits for the retrospective through a stall".
- **Source:** ADR-0019 decision 1: the Retrospective runs "After the Accept phase of each milestone". Accept needs delivered work.
- **Why it matters:** A blocked milestone never reaches Accept, so it never reaches the Retrospective. The only exit is a rule change outside the batch, which is unplanned input.
- **Fix:** Allow an early Retrospective, listed at Intake as a planned point, when a stall names a gate as its cause.

### M20. PRD REQ-001 keeps the overclaim that the Operator found — Severity: material
- **Where:** `PRD-0001` §6: "REQ-001 | `layup psb check` finds the gaps of a problem statement before delivery starts". §13 of this branch changes NFR-005 only.
- **Source:** `check.go`: rules G1–G5 only; architecture §3: "It does not find a gap of meaning".
- **Why it matters:** The requirement still gives meaning gaps to a deterministic check, and its criterion (the golden file) passes with no gap of meaning found.
- **Fix:** Change the REQ-001 statement and criterion: rule gaps by `layup psb check`, meaning gaps by the Intake review, one batch.

### M21. `layup spec check` knows the In-Scope facts only from a PSB form that no text states — Severity: material
- **Where:** ADR-0020 decision 1: "a record that numbers each clause as a fact"; decision 3: "an In-Scope fact has no `Must` row".
- **Source:** Code can number list lines and table rows (the rule of `F-0001`: "For a list line, the fact is the line without its list number"). It cannot find "clauses", and it knows "In-Scope" only from a heading. An agent numbered LAYUP's own facts. No command in §3 does the numbering.
- **Why it matters:** A second pilot PSB with a "Goals" heading gives zero In-Scope facts. The check then passes with any set of requirements. This is the class of the `psb check` defect.
- **Fix:** State the PSB form that a target needs, add a `psb check` gap when it is absent, and name the command that numbers.

### M22. The Operator cannot set a parameter or move a stalled task — Severity: material
- **Where:** ADR-0013 decision 7: "a changed row takes effect at the next step"; only `layup run` writes the register. ADR-0019 decision 3: "Between two retrospectives the routing table does not change."
- **Source:** O-79: "operator-configurable parameters set when running LAYUP against the target"; O-82: "editable per run"; `F-0001#14`: "can give the task to a different harness agent".
- **Why it matters:** No command, flag or comment form sets a value. The Operator's Decision Point 5 choice conflicts with the frozen routing table.
- **Fix:** Name the interface (a flag or an issue-comment form that is copied into Git). Say if a change during delivery is planned or unplanned. Add a per-task harness override.

### M23. The host's user-level harness setup is a second rule source — Severity: material
- **Where:** ADR-0015 decision 6: "The target's rules live only in `AGENTS.md` and `docs/`. Each harness entry file … holds only a pointer"; ADR-0013: sessions run "on the Operator's host".
- **Source:** `F-0003#27`: "Rules, context, and task state are kept in the form of one harness agent". This review session got host rules that no repository file holds: a global `CLAUDE.md`, a command-rewrite hook and a plugin's session-start text.
- **Why it matters:** In a pilot, Claude sessions obey host rules that Codex sessions do not get. The entry-file check still passes.
- **Fix:** Start each session with an empty user-level setup (a clean `HOME`, or the O-77 sandbox). Make the probe record that the harness loaded no user rules.

## Notes

- **N1 (note).** Where: §1 "the idea owner decides the intent, the budget and each escalation (`F-0001#10`–`#14`)". Source: `#14` gives stalls to the Operator. Why: a wrong owner for a reader. Fix: cite `#10`–`#13`.
- **N2 (note).** Where: §9.2 "only the smart-if component calls one"; NFR-005 "A `layup run` opens a connection to a model service only from the smart-if component". Source: role sessions are children of `layup run` and call models. Why: the criterion fails if a test counts child processes. Fix: say "the `layup` process itself".
- **N3 (note).** Where: ADR-0018 "Before each new action"; §9.4 "telemetry per action". Source: ADR-0013: a harness "hides its inner tool calls"; `layup run` sees one session. Why: "action" has no definition, so latency and the cost stop have no unit. Fix: define an action as one session, or as one model call where the harness streams usage.
- **N4 (note).** Where: §8 `harness.<id>.paid`. Source: O-80 makes telemetry incomplete; `F-0003#60` is pass/fail. Why: each task on such a harness fails a Layer 1 check. Fix: state this in §10.
- **N5 (note).** Where: §2 "A scheduled job in LAYUP's own repository". Why: another Operator cannot run it in `pharzam/layup`. It needs a credential for each target and a list of active runs that no record holds, and a GitHub scheduled run can start late. Fix: a job in a repository that the Operator owns, with a named list of runs.
- **N6 (note).** Where: §8 `heartbeat.H` "set at the first run, with its reason". Source: Invariant 4 asks for evidence. Fix: "with its evidence".
- **N7 (note).** Where: ADR-0014 decision 5 `cautious` "may only take the safer branch". Why: model and harness fit and panel disposition have no safer branch, and promotion by "recorded agreement" has no ground truth for them. Fix: state the allowed levels and the agreement measure for each point.
- **N8 (note).** Where: ADR-0020 decision 3 "a quote is not a byte-exact substring of the stored problem statement". Why: a requirement from an answer (as NFR-007 from `F-0004#1`) fails. Fix: accept the answers fact as a trace source.
- **N9 (note).** Where: ADR-0017 decision 4 "its `docs/` rules". Why: the list has no files. The kit's step 7 writes a lesson into `guardrails.md` §2 in the same PR, so the merge is refused. Fix: list the files, and send step-7 lessons to the batch.
- **N10 (note).** Where: ADR-0017 "a gate kind with neither is `not-active`". Why: in the target's CI with LAYUP absent, a missing job looks like a pass, and a failing job blocks every merge. Fix: state the CI form of a not-active gate.
- **N11 (note).** Where: §10. Source: on GitHub Free, protected branches and rulesets work on public repositories only; S01 lets the Operator choose a private target. Fix: state the plan that a private target needs.
- **N12 (note).** Where: §7 `runs/<task>/questions.tsv`. Why: the Intake batch is in no records table, so Early Question Share (`F-0003#75`) has no numerator. Fix: a project-level questions table.
- **N13 (note).** Where: ADR-0019 "that path is work for the implementation plan". Source: it cites `F-0003#39` "the next project gets no lesson". Fix: name the component, or list it in §10.
- **N14 (note).** Where: §3 "per `setup/steps.tsv`". Source: S12 adds the LAYUP "setup-check job" to the target (O-76 forbids it); S02 still uses `npx degit` (ADR-0011 decision 7). Fix: state the S02 and S12 change in §11.
- **N15 (note).** Where: §7 `questions.tsv` column `accepted`. Source: `F-0003#71` ends at "the accepted answer". Why: no actor accepts. Fix: name the actor.

## Coverage: PSB items checked and found answered

- Invariant 2 (`F-0001#2`): native gates in the target, ADR-0017 decision 1.
- Invariant 4 (`#4`): each parameter and price has a source; `layup setup verify` (see N6).
- Invariant 5 (`#5`): `not-active` is not a pass; a smart-if failure goes to a human (see N10).
- Invariant 6 (`#6`): check, then smart-if, then human; code does the arithmetic (see M8).
- Invariant 7 (`#7`): stack gates add rules; the baseline stays.
- Invariant 8 (`#8`): the pin in each target.
- Invariant 9 (`#9`): admission in code, verifier harness is not the author harness (see M23).
- `F-0003#41`: `psb check` plus the Intake review, checked on a counterpart, one batch (see M10, M20).
- `#42`: `layup setup`, `layup setup verify`, the gate recipe.
- `#44`: native gates per recipe, a known-bad fixture per gate kind.
- `#47`: gates and counterpart verification before the orchestrator merges (see M9).
- `#50`: telemetry per task with requirement and price (see N3, N4).
- `#51`: numbered facts, `spec draft`, `spec check`, counterpart verifier (see M21).
- `#52`: records readable from a clean clone; two harnesses (see M23).
- `#53`–`#56`: no weight change; intent with the idea owner; no cloud change; gates only add.
- Decision Point 1 (`#10`): intent at Intake (see M10, M11 for the order).
- Decision Point 2 (`#11`): one batch; the answers stored as a fact.
- Decision Point 3 (`#12`): `approvers.tsv`; `acceptance.tsv` closes a requirement only on accept.
- Structural Conformance (`#58`): only the orchestrator merges, after the gates.
- Telemetry Completeness (`#60`): `not reported` counts as incomplete (see N4).
- Specification Traceability (`#62`): byte-exact trace in `spec check`.
- Setup Correctness (`#63`): `setup verify` and the discipline tests.
- Independence (`#65`): the target's gates run with LAYUP absent.
- Harness-Agent Neutrality (`#66`): admission; `not-active` with fewer than two harnesses.
- Generality (`#67`): two recipes, two pilot targets.
- Delivery Lead Time (`#68`): telemetry times and acceptance time.
- First-Review Acceptance (`#69`): `review_count` in `acceptance.tsv`.
- Clarification Turnaround (`#71`): times in `questions.tsv` (see N15).
- Stall Rate and Resolution (`#73`): outcome in `stall.md` (see M5).
- Cost per Requirement (`#74`): price, cost share, `prices.tsv`.
- Problems 2, 3, 5, 6 (`F-0003#2`, `#3`, `#5`, `#6`): native gates; Intake batch and setup; telemetry and the band; ADR-0020.
- Vision 2.1, 2.3, 3.4, 3.5: phases Design to Plan; the Retrospective; the forge copy and telemetry; stall limits and package (see M6, M15, M19).

Checked and not answered: Invariant 1 (M10, M15); Invariant 3 and `#43` (M1–M3); `#45` (M14); `#46` (M8); `#48` and Decision Point 4 (M7); `#49` (M4–M6); Decision Point 5 (M6, M22); Missed Escalations (M7, M15); Inter-Role Format (M14); Stall Diagnosis (M6); Gate Integrity, "0 agent writes" (M3); Task Intervention Rate and Reversal Rate (M15); Early Question Share (N12); Problems 1, 4, 7 (M8, M4–M6, M23); vision 2.2, 3.1, 3.2, 3.3 (M16, M17, M18).
~~~~
