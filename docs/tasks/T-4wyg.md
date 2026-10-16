# T-4wyg — row 62c of the plan of M2e, the full spec check

Issue: [#351](https://github.com/pharzam/layup/issues/351), row 62c of [the tasks of
M2e](../plan/README.md#the-tasks-of-m2e). Serves `F-0003#51`, `F-0003#62`
(`REQ-012`) and `F-0001#5` (`NFR-004`). Branched off `e16568c`; `main` was
merged once (`9d500fd`, merge-base `cfcfc1c`). Author of the
plan, of the answer and of record: Devin SWE-2 Max (execution tier; the routing
note is on the issue). Evidence: [`runs/T-4wyg/`](../../runs/T-4wyg/).

## Plan and plan review

The plan, its review and the author's complete answer are comments on #351
(6215504530, 6215968646, 6216005137). The independent review (R12, a model that
differs from the author) gave `approve-with-conditions`: budget maximum 2,100
lines added plus removed over 20 files, against `e16568c`, evidence and
close-out inside; cycle cap 1; no panel; goal count 2 under O-214 (ten-rule and
pre-bet checking; the expected `prd-lint` failure). Its thirteen conditions are
applied under the review's own numbering: the classifier is non-vacuous over
all non-empty lines of both streams (C1), both streams are exercised through
the real `sh` with the anchored two-form regex fed by the real captured lines
(C2), the `—` MoSCoW refuses while a `—` Phase reads as empty and an empty
directory exits 0 (C3), the whole-table `clear` applies only before the first
bet and only to a head that touches no spec path (C4), `classes` agrees with a
requirement that waits for its next bet (C5), copy lines are scoped to
REQ/NFR heads by the one shared requirement-line grammar (C6), fact tokens
bind to their record over the caller's two problem-statement records (C7),
`measures` reads only the `Fact` column under a `## … Success metrics` heading
(C8), `waits for a bet` matches case-insensitively at any level to the next
heading with the citation token forms (C9), a `clear` row carries its reason
and `records.md` is amended to match (C10), the second `tasks` value is named
in `records.md` and `packages.md` (C11), `ran` is `ranRow` over the nine rows
and every `did not run` row is `fail` (C12), and the post-bet universe is the
latest confirmed version, an unreadable one a `did not run` (C13).

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
