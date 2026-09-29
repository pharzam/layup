# W-09 — Stall Resolution

**Item.** S9, Stall Resolution (`F-0003#49`): "A disagreement between role
agents, or a step that repeats without progress, gets a procedure with a limit.
An independent examination with a fresh context gives a diagnosis. The project
repository keeps the diagnosis and the outcome." Also Decision Point 5
(`F-0001#14`), Stall Diagnosis (`F-0003#61`), `REQ-009`, `REQ-010`, O-82, and
vision 3.2 and 3.5.

**The cases.** Part 1: in task `T-9`, the verifier on Codex rejects the change
because the invoice total excludes tax; the developer on Claude Code changes the
test, not the code, and the verifier rejects it again for the same finding.
Part 2: in task `T-11`, a developer session stops writing any output.

Sections are those of [`architecture.md`](../architecture.md).

## Part 1 — a disagreement

| # | Actor | Input | Mechanism | Record | Tag | Where |
| - | ----- | ----- | --------- | ------ | --- | ----- |
| 1 | the verifier | the head | rejects: a finding at `invoice/total.go:40` | a review record | `model` | §8 |
| 2 | the developer session | the finding | changes the test; commits | a result | `model` | §8 |
| 3 | the verifier | the new head | rejects again, the same finding | a review record | `model` | §8 |
| 4 | `layup run` | the task's records | open | open | `code` | open |
| 5 | `layup run` | the stall | open | open | `code` | open |
| 6 | an examiner session | the evidence | open | open | `model` | open |
| 7 | the smart-if (P3) | the diagnosis | open | open | `model` | open |
| 8 | panel sessions | the sealed input | open | open | `model` | open |
| 9 | `layup run` | the panel's outputs | open | open | `code` | open |
| 10 | `layup run` | the resolution | open | open | `code` | open |

## Part 2 — a session that hangs

| # | Actor | Input | Mechanism | Record | Tag | Where |
| - | ----- | ----- | --------- | ------ | --- | ----- |
| 11 | `layup run` | the running session | open | open | `code` | open |
| 12 | an examiner session | the evidence | open | open | `model` | open |
| 13 | `layup run` | the second hang | open | open | `code` | open |
| 14 | the Operator | the stall package | open | open | `human` | open |
| 15 | `layup run` | the Operator's answer | open | open | `code` | open |
