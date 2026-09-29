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
