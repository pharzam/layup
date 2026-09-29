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
| 1 | `layup run` | the task register | opens the issue of `T-9` with its requirement IDs and Definition of Done | the issue | `code` | §8; ADR-0019 |
| 2 | a plan session; a plan-review session on another harness | the task | writes the plan; reviews it; `layup run` posts both in the forms the target's `review-record-lint` parses, and reads the verdict field | two issue comments | `model` | §8 |
| 3 | `layup run` | the developer's result | the handoff check; the rule-path and workflow checks; pushes the branch; opens a draft pull request that links the issue; requests no review | the draft pull request | `code` | §8; ADR-0019 |
| 4 | the target's CI, `layup gate` | the head | a failure (W-04 steps 8, 9) | the statuses | `code` | §6 |
| 5 | `layup run` | the failed statuses | gives the developer a new attempt with the findings; the pull request stays a draft | an attempt row | `code` | §8 |
| 6 | `layup run` | the second result | pushes; the native jobs, `layup/gates` and `layup/spec` pass | the statuses | `code` | §6, §7 |
| 7 | a verifier session on Codex | the head; the checklist | admitted because its harness differs from the author's (Claude Code) by the ledger rows; fresh, one turn, read-only, told to refute "done"; verdict `nothing material` | its record | `model` | §8, §9; ADR-0020 |
| 8 | `layup run` | the verdict | posts the record as a pull-request comment in the target's form; sets `layup/verify` at this head SHA | the comment; the status | `code` | §8 |
| 9 | `layup run` | the green head | marks the pull request ready and merges it at that head SHA; requests no human review | the merge | `code` | §8; ADR-0019 |

## Checklist rows

S7, R04, I5 (in part), I9 (in part), K27, K28, P08, P16; FT1, FT3, FT4.
