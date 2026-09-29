# Slice H review record

Reviewer: Claude Fable 5.1 (`claude-fable-5-1`), a fresh `claude -p` session in a disposable clone, read-only tools, no forge token, `--setting-sources project,local`. Not a gate round.

## Round 1

2026-09-29, head `95017d9`; 5 min 48 s; 25,971 output tokens; USD 3.19 at list price. Word for word:

~~~~text
# Slice H review — round 1
Reviewer: claude-fable-5-1. Head: `95017d9`. Verdict: `material findings: 8`.

## Material findings

### M1. After the first update, the learning stops moving the weights
- Where: `docs/architecture.md:1233-1234`: "the old weight plus `learn.step` times the route's reward minus the mean". `§9:795`: "the learned weight of §13 ranks them". ADR-0025 rejects "exploration trials before a baseline exists".
- Source: vision 3.3 ("tune routing"), kept in scope by O-67. Checklist P20, C5.
- Why it matters: in W-13 step 9, route Y ranks first, so every developer task of milestone 3 goes to Y. X gets no tasks, so it never changes. Y is the only route with tasks, so its reward equals the mean and its step is 0. If a new model version makes Y worse, no term can lower Y's weight. The loop learns one time and then stops.
- Fix: add an exploration share (a parameter, for example a fixed share of tasks given to the next pair) or compare each route with a stored baseline, not with the milestone's mean. Otherwise record this as a known limit.

### M2. `layup learn` has no rule that gives a record to a route
- Where: `architecture.md:1219-1230`: "computes, per route (a role, a tier, a harness and a model) … plus a requirement accepted at its first review … minus each material finding, each stall, each unplanned human input … and the cost per task (money, or wall-clock where the money is unknown), each against the milestone's median … A term whose input is missing … is left out for that route". W-13 step 2 is tagged `code`.
- Source: lens 1 (a `code` step names its input and its rule). FT2. O-80, which allows paid work without tokens.
- Why it matters:
  - A requirement spans several tasks and routes (plan, developer, verifier). A verifier's finding is about the developer's route. The text does not say which route gets each term.
  - It does not say whether a count is per task. A route with more tasks collects more minus points.
  - It does not say which routes make up "the mean". Mixing a verifier route with a developer route has no meaning.
  - Money and wall-clock are in different units against one median.
  - A term that is left out acts as a value at the median. So a harness that reports no tokens escapes the cost term, and an expensive route that reports its cost ranks below it.
- Fix: state the attribution of each term to a route. Normalize each term per task. Take the mean within one role and tier. Use one median per unit. Give an unknown cost a penalty, or bar that route from any upward step.

### M3. The walkthrough reads the audit before it runs, and its approver is not in the one comment
- Where: W-13 step 2's input is "… the audit", and step 3 reads "the audit". But step 6 "runs the audit (§12)" after the approval. `§12:1181-1192`: the sample holds only tasks "merged at least 30 days before", and "The idea owner confirms each positive at the retrospective". `§13:1252`: "The approver answers by one comment" (the Operator by default).
- Source: ADR-0025 decision 1 ("confirmed reversals … subtract"). `F-0001#12` (each planned point names its approver).
- Why it matters: at milestone 2, the reversal term can only read an earlier audit about earlier milestones' tasks. Nothing says which milestone's route pays for it. The idea owner's confirmation is a second answer at the same point. It is neither in the brief nor in W-13.
- Fix: run the audit, and get the idea owner's confirmation, before `layup learn`. Say which milestone's routes a confirmed reversal charges. Add the idea owner's comment as a W-13 step.

### M4. W-13 steps 4 and 6 skip the batch task and break O-93
- Where: W-13 step 4 (`code`): "the rule batch (with W-03's refused change to the gate, as a verified batch task, §6)". Step 6 (`code`): "`layup run` … merges the approved batch".
- Source: O-93. §6:463-465 ("by its approver, who also pushes it"). §8:614 ("each rule batch … goes through the task loop"). Lens 1.
- Why it matters: W-03's change is to the CI job of the layout gate, which is under `.github/workflows/`. The App has no workflows permission, so `layup run` can neither push nor merge this batch. Also, a "verified batch task" needs plan, plan-review, developer and verifier sessions. Code cannot make one from a stored proposal, and a batch that a lesson proposes cannot exist before step 4.
- Fix: add the batch task's `model` steps between step 3 and step 4. Add the approver's push and merge as `human` steps.

### M5. A retrospective batch that changes a gate is never tested against its known-bad patch
- Where: §6:441-443 runs the known-bad patches only for "the activation batch". §13:1249-1257 and W-13 step 6 merge a gate-changing batch with no such run.
- Source: `F-0003#64` ("100% detection of known-bad commits"). §11:1084 (an early retrospective's batch "may fix the gate").
- Why it matters: W-03's change was made by an agent so that a failing change would pass. If the approver approves it (W-13 step 5), the weakened gate lands. No check shows that it still fails its known-bad patch.
- Fix: apply the activation rule to every batch that touches a gate kind. After approval, run the batch's own gate files on its head and on each known-bad patch of that kind, and refuse the merge when a patch passes.

### M6. The trigger policy is not a parameter
- Where: `§10:972` lists "the reward's term weights; `learn.step`, `learn.min_weight`, `learn.max_weight`, `learn.min`; the learning authority". The old row named "trigger". §13:1215 fixes "When".
- Source: O-83: "The weight adjustments, their bounds and the trigger policy of the learning loop are parameters." Checklist D14.
- Fix: add a trigger parameter (for example: each retrospective, every N milestones, or also the early retrospective). Or state which listed parameter is the trigger and why.

### M7. The early retrospective counts the same tasks twice, or its scope is undefined
- Where: §13:1215: "After the Accept phase of each milestone, and early when a stall's diagnosis names a gate". §13:1219: `layup learn` "reads the milestone's records".
- Why it matters: an early retrospective in the middle of milestone 3 runs `layup learn` on part of the milestone. The retrospective at the end then reads the same tasks again, so the weights move twice on the same evidence. The early medians come from a small population.
- Fix: say whether `layup learn` runs at an early retrospective. If it does, limit each run to records that no earlier learning run read.

### M8. The retrospective session has no row in the step table
- Where: §13:1243: "A retrospective session reads the records of the milestone …". W-13 step 3 (`model`). The step table (§9:758-767) has no retrospective row.
- Source: §9 ("Routing is … for each role and tier"). ADR-0020.
- Why it matters: code cannot pick the role, the tier, the harness or the prompt kinds for this session. The session judges stalls and findings, so being outside their authors matters, and nothing states that either.
- Fix: add a row, for example Systems Architect on the reasoning tier, on a harness outside the milestone's authors where one exists, with its prompt record kinds.

## Notes
- N1. W-13 step 7: an installation token covers only the repositories where the App is installed. An Operator who does not own LAYUP's repository cannot open the issue. Name the credential. Also, a lesson that quotes a client target goes to LAYUP's public repository; add a rule against exposing the target's data.
- N2. §5 Intake does not record the LAYUP version, but §13:1263 and W-13 step 9 depend on it. Also say whether a running target changes its LAYUP version when the host upgrades.
- N3. There is no rule for the first "old weight", and no rule for ranking a mix of pairs with and without weights (§9). The evaluated pattern has a rollback (adaptive-model-router); ADR-0025 has none.
- N4. P4 is asked only on a tie, so its calibration record (§10) gets almost no rows, and a `delegate` P4 almost never acts. Say that this is accepted.
- N5. D03: inside a target, a lesson never changes routing. O-69's "learned … as RL" is met only through the records-based reward and a later LAYUP release. State this reading.
- N6. Invariant 4 requires evidence for the term weights, but no fact can support a reward weight. Say that the evidence is the Operator's comment, as the harness caps already do.
- N7. The brief should mark where each rule proposal comes from (a refused session edit or a lesson). Then the approver can see that W-03's change was an agent's attempt to pass a failing change.

## Checklist rows
- K59: answered, through the lesson issue and a LAYUP release (see N1 and N2).
- K60: answered: the weight ranks at every authority level (see M1 for how the weights move).
- K61: the terms are answered and the `#58` conflict is stated, but the rule that gives each term to a route is missing (M2).
- P01: answered: approval at the retrospective, a planned point (O-69).
- P17: answered, subject to M3 and M7.
- P20: not answered in effect: the learning freezes after one update (M1).
- C5: answered in parts; M1, M2 and M6 are open.
- D03: answered; the lesson-to-routing path is a reading (N5).
- D14: not answered: the trigger is not a parameter (M6).
- FT6: answered: the weights, the lessons and the brief copy are on the records branch; the priors are in LAYUP's Git.

## Existing solutions
- The design repeats Hivelore and smirnowld/agent-practices (reviewed lessons kept in Git) and the bounded, human-adopted weights of guangyang1206/adaptive-model-router, but leaves out that router's rollback.
- It ignores the standard answer to M1: bandit routing with exploration (epsilon-greedy, UCB, Thompson sampling) and a share of exploration inside the budget.
- GitHub: an App needs the `workflows` permission to write workflow files (docs.github.com, "Permissions required for GitHub Apps"). This supports O-93 and M4.~~~~

### The author's answer to round 1

All eight material findings and notes N1 to N7 are applied in one commit (§13 rewritten).

| Finding | Fix |
| ------- | --- |
| M1 | An exploration share `learn.explore` (a parameter); with zero, known limit L-H1. |
| M2 | Terms are charged per task to its implementing route; the mean within one role and tier; money and wall-clock each against its own median; a route with unknown money gets no upward step; other roles' routes: L-H2. |
| M3 | The order is the audit and the idea owner's confirmations, then `layup learn`; a confirmed reversal is charged at the retrospective that confirms it; W-13 step 2. |
| M4 | The batch is a task with its sessions before the brief (W-13 step 5); the Operator pushes and merges a workflows batch. |
| M5 | Every rule batch that changes a gate kind must fail on each known-bad patch of that kind before it merges (§6, ADR-0016 d4). |
| M6 | `learn.trigger` is a parameter. |
| M7 | An early retrospective carries only a rule batch; each `layup learn` run reads only unread records. |
| M8 | The step table has a retrospective row. |
| N1 to N7 | The issue is opened where the App is installed, else recorded for the Operator, with no target content; Intake records the LAYUP version; the first weight, ranking with and without weights, and a revert; P4's few calibration rows are accepted; the O-69 reading is stated; term weights' evidence is the Operator's comment; the brief marks each proposal's source. |

## Round 2

2026-09-29, head `dbed09a`; 3 min 26 s; 16,613 output tokens; USD 2.00 at list price. Same reviewer and setup. Word for word:

~~~~text
# Slice H review — round 2
Reviewer: claude-fable-5-1. Head: `dbed09a`. Verdict: `material findings: 6`.

Round 1: M3, M4, M6, M7 and M8 are closed. M1, M2 and M5 are closed in their main case, but their fixes add the new defects M1 to M5 below.

## Material findings

### M1. A task with unknown money still pulls the mean, so known-cost routes step down
- Where: `docs/architecture.md:1241-1243`: "Cost is compared by unit: money against the median of the tasks with known money, wall-clock against the median wall-clock. A route with any task of unknown money gets no upward step (FT2)." ADR-0025 decision 2 says the same.
- Source: FT2. O-80 (paid work without tokens is allowed by default). Lens 1: a `code` step names its rule.
- Why it matters: The fix removed the old rule "a missing term is left out". Nothing replaced it, so the value of the money term for a task with unknown money is undefined. Take two routes on one tier: route X reports its money and is cheap; route Z runs on a harness with no token report. If Z's money term is left out, Z's value is higher, and so is the mean "over the routes of the same role and tier". X then falls below the mean and steps down. The unknown cost acts as zero against X. Z is also barred from ever stepping up, but it can still step down.
- Fix: state the money term of a task with unknown money. Keep routes with unknown money out of the mean that other routes are compared with, or give that term a stated penalty.

### M2. "Reads only unread records" contradicts "tasks since its last update"
- Where: `§13:1226`: "each run reads only records that no earlier run read". `§13:1245-1246`: "A route with fewer than `learn.min` tasks since its last update gets no change." ADR-0025 decision 1.
- Why it matters: Take `learn.min` = 5. Exploration gives route X 3 tasks in milestone 2 and 3 in milestone 3. Under the first sentence, the run at milestone 3 cannot see the first 3 tasks, so X never reaches 5 and never updates. Under the second, X has 6 and updates. Exploration is the fix for round 1's M1, and small exploration routes are exactly the ones this blocks. It is also not said whether "last update" means the proposal or its adoption under `propose`.
- Fix: a run reads, for each route, the records since that route's last adopted update, and a route's records carry over while it gets no change.

### M3. A route with no weight has no update rule, and §9 and §13 disagree on how it ranks
- Where: `§13:1248-1252`: "the old weight plus `learn.step` times …"; "The first weight of a route is LAYUP's prior for it, or none". `§13:1256-1257` and ADR-0025 decision 4: "the fit point is asked only on a tie". `§9:804-806`: "with no weight yet, or a tie, the smart-if's fit point (§10) or the table's order picks one".
- Source: lens 1 (`code` rule). A contradiction inside the design.
- Why it matters: A harness added to the register with no LAYUP prior has weight "none". The formula has no old weight for it, so it can never get a weight. Because "a pair with a weight ranks before one without", it gets tasks only if it happens to be the "next" pair for exploration. At the first milestone with no priors, §9 asks P4 but ADR-0025 does not.
- Fix: give a route's first weight a stated value (for example the mean of the weighted routes of its role and tier, or the middle of the bounds). Make §9 and ADR-0025 decision 4 say the same thing about pairs with no weight.

### M4. §6 still says the activation batch is the only exception to FT4
- Where: `§6:429-431`: "It never runs the head's gate files (FT4), with one exception: the activation batch below". `§6:436-439` (new): "Each rule batch that adds or changes a gate kind … `layup gate` runs the batch's own gate files on its head".
- Source: FT4. A contradiction within §6, caused by this fix.
- Why it matters: An implementer who follows line 429 refuses to run a retrospective batch's own gate files, so the check that round 1's M5 fix added never runs.
- Fix: change the exception to "an approved rule batch, whose rule files are the approved ones".

### M5. Which known-bad patches run is ambiguous, and the batch can write its own easy patch
- Where: `§6:436-439`: "carries one known-bad patch per kind it touches; … on the head with each patch". ADR-0016 decision 4: "carries one known-bad patch per kind it touches … and each must fail". But `§13:1277`: "on each known-bad patch of that kind (each must fail)".
- Source: `F-0003#64` (100% detection of known-bad commits). FT1.
- Why it matters: Under §6, only the new patch that the batch's own model sessions wrote is run. Suppose a batch weakens the layout test so it skips one directory, and brings a patch with a violation somewhere else. That patch fails, the check passes, and the activation patch for the skipped directory is never run again. Also, an older patch may not apply to a later head, and nothing says that such a patch refuses the merge rather than being skipped.
- Fix: run every recorded known-bad patch of each touched kind, plus the batch's new one. A patch that does not apply refuses the merge (FT1). Make §6, §13 and ADR-0016 decision 4 say the same thing.

### M6. "A target keeps its LAYUP version until a bet changes it" has no mechanism
- Where: `§13:1288-1290`: "Intake records the LAYUP version (§5), and a target keeps it until a bet changes it". `§5:283`: "records the LAYUP version". W-13 step 10.
- Source: lens 1. No step in §8 changes a version at a bet (`grep -n version` finds none in §8).
- Why it matters: The host upgrades `layup` from v1 to v2 in the middle of a target recorded at v1. Nothing says whether v2 refuses to run, runs with v1's defaults (which it would have to carry), or uses its own new prompts. So "keeps it" cannot be implemented, and the new plan-review prompt may reach a running target without any approval.
- Fix: name the record, the check at each `layup run` start (for example: stop on a mismatch), and the bet step that adopts a new version.

## Notes
- N1. Exploration goes to "the next admitted pair in the table's order". It is unclear which pair that is when the leader is not first in the table. A third pair is never explored. Exploring plan, review and verification roles spends budget, but those routes never learn (L-H2).
- N2. Round 1's N4 ("P4's few calibration rows are accepted") is listed as applied, but no sentence in §10 or §15 says it.
- N3. W-03 step 6 is tagged `human`, but its mechanism now includes the batch task's model sessions. Tag it `model`, `human`, or point only to W-13 steps 5 to 7.
- N4. W-13 step 5 builds a gate-changing batch but does not write the known-bad patch that §6 requires.
- N5. The idea owner confirms the audit in W-13 step 2 (tagged `human`), but is not in that step's actor column.
- N6. `§12:1186` runs the audit "At each retrospective", but `§13:1223` says an early retrospective "carries only a rule batch".
- N7. The developer's route is charged for a reversal of an owner role's answer. The reversals confirmed at the final audit, 30 days after delivery, are never charged. Say both.
- N8. When the Operator merges a batch that changes workflows, `layup/rules` checks only the hash. No forge check blocks a merge after a known-bad patch has passed. Nothing records what happens to an approved batch whose merge is refused.
- N9. `§5:283` has a stray 6-space indent. A LAYUP issue that links a target's records still exposes the target's name.

## Checklist rows
- K59: answered, subject to M6.
- K60: answered, subject to M3.
- K61: answered; the unknown-money case is open (M1).
- P01: answered.
- P17: answered.
- P20: answered through `learn.explore`, but M2 can still freeze a small route.
- C5: answered in parts; M1, M2 and M3 are open.
- D03: answered; the reading of O-69 is stated.
- D14: answered (`learn.trigger`).
- FT6: answered (the routing register, the patches and the version are on the records branch; the priors are in LAYUP's Git).
- L-H1: known limit (acceptable).
- L-H2: known limit (acceptable).

## Existing solutions
- `learn.explore` is a fixed-share, epsilon-greedy bandit, but it is deterministic and only reaches the "next" pair. UCB or Thompson sampling would explore every arm and would solve the cold start in M3 (untried arms get an optimistic value).
- Model routers (RouteLLM, Not Diamond, the adaptive-model-router of `selection-v2.md`) learn from outcomes per task. The bounded, human-adopted weights, now with a revert, follow the last one.
- Re-running every recorded known-bad patch (M5) is mutation-testing regression (go-mutesting, Stryker): old mutants are kept and must stay caught.
- GitHub: an installation token reaches only the repositories where the App is installed, and writing workflow files needs the App's `workflows` permission. Both match the design (§13, O-93).~~~~

### The author's answer to round 2

The six new material findings and notes N1 to N9 are applied in the next commit. Under O-98, these fixes are a named item of the whole-design review.

| Finding | Fix |
| ------- | --- |
| M1 | A route with any unknown money gets no step and is left out of the mean (L-H2). |
| M2 | Each run reads, per route, the records since its last adopted update; records carry over. |
| M3 | A route's first weight is the prior, else the mean of its role and tier's weights, else the middle of the bounds; §9, §13 and ADR-0025 d4 agree: the fit point or the table's order only with no weight or on a tie. |
| M4 | §6's exception is an approved rule batch. |
| M5 | Every recorded known-bad patch of each touched kind and the batch's new one must fail; a patch that no longer applies refuses the merge (§6, §13, ADR-0016 d4). |
| M6 | `layup run` stops when its version differs from the recorded one; a bet adopts a new version by one line. |
| N1 to N9 | Exploration goes to the other pairs in turn, implementing roles only; P4's few rows accepted (L-E2); W-03 step 6 tagged `model`, `human`; W-13 step 5 writes the patch; the idea owner is an actor of W-13 step 2; the audit runs at each end-of-milestone retrospective; reversals charge the answer's task route, the final audit's none; a refused batch merge goes back as a task; indentation; links only for a public target. |
