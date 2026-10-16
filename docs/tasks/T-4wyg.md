# T-4wyg — row 62c of the plan of M2e, the full spec check

Issue: [#351](https://github.com/pharzam/layup/issues/351), row 62c of [the tasks of
M2e](../plan/README.md#the-tasks-of-m2e). Serves `F-0003#51`, `F-0003#62`
(`REQ-012`) and `F-0001#5` (`NFR-004`). Branched off `e16568c`. Author of the
plan, of the answer and of record: Devin SWE-2 Max (execution tier; the routing
note is on the issue). Evidence: [`runs/T-4wyg/`](../../runs/T-4wyg/).

## Plan and plan review

The plan, its review and the author's complete answer are comments on #351
(6215504530, 6215968646, 6216005137). The independent review (R12, a model that
differs from the author) gave `approve-with-conditions`: budget maximum 2,100
lines added plus removed over 20 files, against `e16568c`, evidence and
close-out inside; cycle cap 1; no panel; goal count 2 under O-214 (ten-rule and
pre-bet checking; the expected `prd-lint` failure). Its thirteen conditions are
applied: the strict two-form classifier over both streams (C1–C4), the
source-bound mode precedence (C5), the shared requirement-line grammar (C6),
record-bound fact tokens over the caller's two problem-statement records (C7),
success-criterion-only `measures` (C8), the both-directions `classes`
comparison with `—` Phase read as empty (C9), the clear-row reasons and the
whole-table emission (C10), the hash-verification boundary (C11), no CLI/loop
or forge wiring claimed (C12), and this task's docs, traceability and
known-pitfall writeback (C13).

## What was done

`internal/spec/full.go` — `CheckFull` (the ten rule rows in the block's order,
every mode returning ten rows), `TaskSpecCheckSchema`, the second Go value of
`tasks`, `WriteSpecCheck`, and the nine judged rules plus `ran`;
`internal/spec/prdlint.go` — `RunPRDLint` (the one permitted `sh` seam, streams
kept apart), `WaitingSpecAccepted`, `waitingVerdict` and `emptyCellFailure`;
`internal/spec/answer.go` — `parseReqLine`, the one requirement-line grammar
shared by the answer and `classes`; the unit and integration tests; the block
move of `task-spec-check` to `built`; `loop.md`, `records.md`, `guardrails.md`
§2 and the §12 cells of `PRD-0001`. The pinned `docs/prd/prd-lint.sh` is not
modified; `internal/cli` is not touched; no `layup run` wiring is claimed.

## Review rounds

None yet — this file records the implementation before the frozen-head rounds.
