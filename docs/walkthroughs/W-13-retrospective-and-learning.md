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
| 1 | `layup run` | the milestone's acceptances | every requirement of milestone 2 is accepted or rejected: the Retrospective phase starts (a planned point that Intake listed, approver: the Operator) | a phase row | `code` | §8, §13 |
| 2 | `layup learn` | the milestone's records: telemetry, acceptances, findings, stalls, human inputs, the audit | the reward per route from the term weights; route Y is above the mean and route X below; both have at least `learn.min` tasks; proposes Y's weight up and X's down by `learn.step`, inside the bounds | a reward table and a routing proposal | `code` | §13; ADR-0025 |
| 3 | a retrospective session | the stalls, diagnoses, findings, escalations and the audit | writes two lessons with their records: `target` ("a plan's test contradicted the specification": a §2 pitfall) and `layup` ("the plan review's checklist compares each test with its specification section") | the lessons | `model` | §13 |
| 4 | `layup run` | the proposals | builds one brief: the reward table, the routing proposal, the rule batch (with W-03's refused change to the gate, as a verified batch task, §6), the lessons | the brief on the milestone's issue | `code` | §13; ADR-0025 |
| 5 | the approver | the retrospective brief | answers by one comment: adopt the routing proposal; approve the rule batch by its head and hash, or not; keep both lessons (O-69) | the comment | `human` | §13 |
| 6 | `layup run` | the approval | copies it (§3); writes the new weights to the routing register; merges the approved batch; adds the §2 pitfall as a task (added lines, no batch); runs the audit (§12) | the routing register; the merges | `code` | §6, §13 |
| 7 | `layup run` | the lesson for LAYUP | opens an issue on LAYUP's repository with the lesson and the links to its records | the issue | `code` | §13 |
| 8 | the Operator and a LAYUP task | the issue on LAYUP's repository | the LAYUP task, under LAYUP's gate, changes the plan-review prompt of the default step table; a LAYUP release carries it | the LAYUP release | `human`, `model` | §13; ADR-0025 |
| 9 | `layup run` | the next milestone's first task | ranks the admitted pairs by the new weights: route Y first; the next target's Intake records the new LAYUP version, so its plan reviews use the new prompt | a routing row | `code` | §9, §13 |

## Checklist rows

K59, K60, K61, P01, P17, P20, C5, D03, D14; FT6. No known limit of its own.
