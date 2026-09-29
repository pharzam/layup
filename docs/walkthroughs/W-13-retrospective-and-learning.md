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
| 1 | `layup run` | the milestone's acceptances | open | open | `code` | open |
| 2 | `layup learn` | the milestone's records | open | open | `code` | open |
| 3 | a retrospective session | the records | open | open | `model` | open |
| 4 | `layup run` | the proposals | open | open | `code` | open |
| 5 | the approver | the retrospective brief | open | open | `human` | open |
| 6 | `layup run` | the approval | open | open | `code` | open |
| 7 | `layup run` | the lesson for LAYUP | open | open | `code` | open |
| 8 | a LAYUP task | the issue on LAYUP's repository | open | open | `human`, `model` | open |
| 9 | `layup run` | the next milestone's first task | open | open | `code` | open |
