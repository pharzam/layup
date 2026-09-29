# W-13 — Retrospective and learning

**Item.** Vision 2.3 ("automated retrospectives at the close of every milestone,
committing actionable lessons learned to iteratively refine subsequent execution
cycles") and 3.3 (reward-driven routing), kept in scope by O-67; O-69 (rule
changes at the retrospective, approved there, fed to learning) and O-83 (the
learning parameters); PSB §5 ("the next project gets no lesson"); and the PSB's
Out of Scope (no change to model weights).

**The case.** Milestone 2 of the invoicing target is accepted. Its developer
tasks ran on two routes: Codex with model X (4 tasks: 3 requirements accepted at
the first review, 1 stall, higher cost) and Claude Code with model Y (5 tasks: all
accepted at the first review, lower cost). W-09 found that a plan's test
contradicted its specification. W-03 left a refused change to a gate as a
proposal.

Sections are those of [`architecture.md`](../architecture.md).

| # | Actor | Input | Mechanism | Record | Tag | Where |
| - | ----- | ----- | --------- | ------ | --- | ----- |
| 1 | `layup run` | the milestone's acceptances | every requirement of milestone 2 is accepted or rejected: the Retrospective phase starts (a planned point that Intake listed, approvers: the idea owner for the audit, the Operator for the batch) | a phase row | `code` | §8, §13 |
| 2 | `layup run`, an auditor session, the idea owner | the items merged at least 30 days before | the audit of §12; the idea owner confirms its positives by one comment (none here) | `audit.tsv` | `code`, `model`, `human` | §12, §13 |
| 3 | `layup learn` | the records since each route's last adopted update: telemetry, acceptances, findings, stalls, human inputs, the audit | per task of each implementing route, the terms; route Y's mean is above the mean of the developer routes on its tier, and X's below; X has a stall; both have `learn.min` tasks; proposes Y up and X down by `learn.step`, inside the bounds | a reward table and a routing proposal | `code` | §13; ADR-0025 |
| 4 | a retrospective session (Systems Architect, reasoning tier, on a harness outside the milestone's authors) | the stalls, diagnoses, findings, escalations and the audit | writes two lessons with their records: `target` ("a plan's test contradicted the specification": a §2 pitfall) and `layup` ("the plan review's checklist compares each test with its specification section") | the lessons | `model` | §13 |
| 5 | batch sessions | W-03's refused change (marked: a session's attempt to pass a failing change) | a rule batch task: plan, plan review, change with a new known-bad patch for the layout kind, verification (§8) | the batch's pull request, head and hash | `model` | §6, §8 |
| 6 | `layup run` | the proposals | one brief to the Operator: the reward table, the routing proposal, the batch with its source, the lessons | the brief on the milestone's issue | `code` | §13; ADR-0025 |
| 7 | the Operator | the brief | adopts the routing proposal; rejects the batch (its source is an attempt to pass a failing change); keeps both lessons (O-69) | the comment | `human` | §13 |
| 8 | `layup run` | the answer | copies it (§3); writes the new weights to the routing register; the rejected batch closes; the §2 pitfall goes in as a normal task | the routing register | `code` | §13 |
| 9 | `layup run`; the Operator and a LAYUP task | the lesson for LAYUP | opens an issue on LAYUP's repository (the App is installed there) with the lesson and links, no target content; a LAYUP task, under LAYUP's gate, changes the default plan-review prompt; a LAYUP release carries it | the issue; the release | `code`, `human`, `model` | §13; ADR-0025 |
| 10 | `layup run` | the next milestone's tasks | ranks the admitted pairs by the new weights: route Y first, with `learn.explore` of the tasks on the other pairs in turn; a new target's Intake records the new LAYUP version, so its plan reviews use the new prompt | routing rows | `code` | §9, §13 |

## Checklist rows

K59, K60, K61, P01, P17, P20, C5, D03, D14; FT6. Known limits L-H1, L-H2.
