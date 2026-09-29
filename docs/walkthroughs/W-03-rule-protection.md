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
| 2 | `layup run` | the session's commits; the rule-path register of the base branch | compares the changed paths with the register: `.githooks/` and the CI job are rule paths, and no approved batch covers them; refuses the result before any push | a refusal row; the edits stored as a proposal for the next retrospective | `code` | §6; ADR-0017 |
| 3 | `layup run` | the refused result | gives the task the finding "rule path changed"; the next attempt starts from the base commit | a finding row for `T-12` | `code` | §6 |
| 4 | the forge | a push or a merge that does not come from `layup run` | a role session holds no credential; the default branch's ruleset requires a pull request and the `layup/rules` check from the App, which fails on a rule path outside an approved batch | — | `code` | §6; ADR-0017 |
| 5 | `layup audit` | the repository activity; the default branch | lists each rule-path change with its merge actor and its approval; a change outside an approved batch counts against `F-0003#64` | an audit table, committed by `layup run` | `code` | §3, §6; ADR-0017 |
| 6 | the retrospective | the proposal | later: slice H | later: slice H | `human` | — |

## Checklist rows

S3, I3, K09, K13, K18, D03 (the batch), D07; FT3, FT4.
