# W-03 — Rule Protection

**Item.** S3, Rule Protection (`F-0003#43`): "Protection of the rules from the
agents that the rules govern." Also Invariant 3 and Gate Integrity (`F-0003#64`:
"100% detection of known-bad commits. 0 agent writes to rule paths.").

**The case.** During task `T-12`, a developer session edits the target's
pre-commit hook and the CI job of its layout gate, so that its failing change
passes.

Sections are those of [`architecture.md`](../architecture.md).

| # | Actor | Input | Mechanism | Record | Tag | Where |
| - | ----- | ----- | --------- | ------ | --- | ----- |
| 1 | the developer session | its work tree | commits the two edits with its change | commits in `work/` | `model` | §4 |
| 2 | `layup run` | the session's commits; the rule-path list | open | open | `code` | open |
| 3 | `layup run` | the refused result | open | open | `code` | open |
| 4 | the forge | a push or a merge that bypasses `layup run` | open | open | `code` | open |
| 5 | `layup audit` | the repository activity; the default branch | open | open | `code` | open |
| 6 | the retrospective | the proposal | later: slice H | later: slice H | `human` | — |
