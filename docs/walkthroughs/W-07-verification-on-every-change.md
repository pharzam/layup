# W-07 — Verification on Every Change

**Item.** S7, Verification on Every Change (`F-0003#47`): "Each change passes the
gates for repository layout, interface boundaries, and testing pyramids before it
reaches human review or merges." Also Structural Conformance (`F-0003#58`),
Harness-Agent Neutrality (`#66`, the independent verification), Invariants 5 and
9, and `REQ-007`.

**The case.** Task `T-9` of W-06, made on Claude Code. Its first push fails the
boundary gate (W-04). The second passes the gates and is verified on Codex.

Sections are those of [`architecture.md`](../architecture.md).

| # | Actor | Input | Mechanism | Record | Tag | Where |
| - | ----- | ----- | --------- | ------ | --- | ----- |
| 1 | `layup run` | the task | open | open | `code` | open |
| 2 | a plan session, a plan-review session | the task | open | open | `model` | open |
| 3 | `layup run` | the developer's result | open | open | `code` | open |
| 4 | the target's CI, `layup gate` | the head | a failure (W-04 steps 8, 9) | the statuses | `code` | §6 |
| 5 | `layup run` | the failed statuses | open | open | `code` | open |
| 6 | `layup run` | the second result | open | open | `code` | open |
| 7 | a verifier session on Codex | the change | open | open | `model` | open |
| 8 | `layup run` | the verdict | open | open | `code` | open |
| 9 | `layup run` | the green head | open | open | `code` | open |
