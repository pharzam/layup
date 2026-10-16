# T-0pt8 — row 62a of M2e, the Shape filter

Issue: [#349](https://github.com/pharzam/layup/issues/349), row 62a of [the tasks of
M2e](../plan/README.md#the-tasks-of-m2e), opened by `T-57vx` (#346). Serves `F-0003#51`
(`REQ-012`). Base `48f07ca`; the budget is read against the current `main`. Author: Claude
Sonnet 5.5 on Claude Code. Evidence: [`runs/T-0pt8/`](../../runs/T-0pt8/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #349. The plan review
(GPT-6.1 Sol on Codex CLI, effort `xhigh`, a fresh sandboxed session in a clone at
`48f07ca`, 4 min 25 s; the harness is outside the routing rule of ADR-0026 and was allowed
by the Operator on 2026-10-13) gave `approve-with-conditions`: Budget maximum 600 lines
added plus removed over 12 files, close-out inside; Cycle cap 1; no panel. Its five
conditions are applied: the sources of the constraints (the draft or the PRD, the facts,
the invariants, the records) and the owner rules of `options`, the explicit test cases and
a mutant for each rule, the validation step, and the fact `F-0003#51` in the plan.

## What was done

1. **T:** `internal/spec/shape_test.go`, red before the code (no `Filter`) and again for each
   later case; one mutant for each rule fails its own test
   ([`test-runs.md`](../../runs/T-0pt8/test-runs.md)).
2. **C:** `internal/spec/shape.go`: `OptionsSchema` (the block `options`, moved to `built`),
   `ReadOptions` and its owner rules, `Standing` (the constraints that apply, from the
   records), `Filter` and `Chosen`.
3. **D:** a sentence of `loop.md` (R10: an amended record applies, a fact of the answers keeps
   its `Q-NNN`), the rows of `traceability.md` and the cells of `REQ-012` in `PRD-0001` §12.

**For row 66h (#364):** the readers of the PRD, the invariants and the decision records, the
source of the facts, and two rules that no code checks, are in a comment on #364.

## Review rounds

The records are comments on #349. Both rounds were Claude Opus 5.5, effort `xhigh`, a fresh
read-only `claude -p` session, a model that wrote none of the change.

| Round | Head | Cycle | Cap | Material | Then |
| ----- | ---- | ----- | --- | -------- | ---- |
| 1 | `7b2e5c3` | 0 | 1 | 2 | fixes (an amended record applies; a fact of the answers keeps its ID) |
| 2 | `798c116` | 1 | 1 | 0 | `nothing material in scope`, four notes, applied in `0eb4915` |

Main moved twice after round 2 (rows 60 and 61). The merges of `origin/main` changed shared
index lines (the cells of `PRD-0001` §12, `traceability.md`, and the lists of
`blocks_integration_test.go`, the last resolved by keeping both sides' moves) and no other
file of this branch, so by **O-217** and the standing **O-219** they need no scoped round and
round 2's verdict stands. The verdict of round 1 was edited to `material`, as RR8 asks.

## Verdict

Delivered: the Shape filter of `M2e` in `internal/spec` and the block `options`, ended by
decay at cycle 1 of the cap 1 (`nothing material in scope`). The diff against the current
`main` is inside 600 lines over 12 files.

## Resource record

Recorded, not budgeted (ADR-0007). UTC; the author's tokens are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Sonnet 5.5 | not reported | not reported | 2026-10-15 17:22 to 17:25 |
| The plan review | reasoning | GPT-6.1 Sol | `xhigh` | 140,263 | 4 min 25 s |
| T and C, the validation | execution | Claude Sonnet 5.5 | not reported | not reported | 17:25 to 17:52 |
| Round 1 | reasoning | Claude Opus 5.5 | `xhigh` | 4,086,414 (USD 3.08) | 11 min 1 s |
| The fix of round 1 | execution | Claude Sonnet 5.5 | not reported | not reported | 18:03 to 18:25 |
| Round 2 | reasoning | Claude Opus 5.5 | `xhigh` | 2,841,236 (USD 2.51) | 9 min 57 s |
| The notes, the merges, the close-out | execution | Claude Sonnet 5.5 | not reported | not reported | 18:35 to 2026-10-16 05:45 (the waits for the Operator's rulings are in it) |
