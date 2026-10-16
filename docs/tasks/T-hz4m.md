# T-hz4m — row 65a of the plan of M2e, the gate jobs read the base

Issue: [#355](https://github.com/pharzam/layup/issues/355), row 65a of [the tasks of
M2e](../plan/README.md#the-tasks-of-m2e), opened by `T-57vx` (#346). Serves `F-0003#47`
(`REQ-007`, F-22). Base `6bc1051`; `main` was merged three times, at `2aafa64`, `d24fae3`
and `f12da6a`. Pull request [#381](https://github.com/pharzam/layup/pull/381). Author:
Claude Sonnet 5.5 on Claude Code. Evidence: [`runs/T-hz4m/`](../../runs/T-hz4m/).

## Plan and plan review

The plan, its review and the author's answer are comments on #355 (6208797865, 6209093081,
6209117686). The plan review (GPT-6 Sol on Devin CLI) gave `approve-with-conditions`: Budget
maximum 650 lines over 14 files, close-out inside; Cycle cap 2 (a change to a gate script and
its workflow); no panel. Its three conditions are applied: the documents of the slice include
`architecture.md` L-B3; the cases of a directory `config` path and of a link at or above a
destination are in the tests; the match of `batch/<point>` is one segment, decided in
`loop.md` (R13).

## What was done

The step of each job of the catalog's workflow and the mode `--overlay` of the job script,
which puts the base's manifest, script and `config` paths into the checkout and the index
(after checking every destination), and takes the head's own files only for `batch/` and one
segment; `TestTheJobsReadTheBase` (25 cases, in `sh`, `dash` and `bash --posix`), red twice
before the code, 22 mutants all caught; `setup.md`, `loop.md`, `architecture.md` L-B3,
`traceability.md`, `PRD-0001` §12 and one pitfall in `guardrails.md` §2. **Known limits,
written in `setup.md`:** the job refuses every link of the head at or above a destination
where `layup gate` follows one into the tree; the exception looks at the branch name, not the
repository (#345, row 65b). **Off the path:** none opened by this task.

## Review rounds

| Round | Head | Cycle | Cap | Material | Then |
| ----- | ---- | ----- | --- | -------- | ---- |
| 1 | `511cc4e` | 0 | 2 | 0 | `nothing material in scope`, three notes |

Round 1 (comment 6210141901) was Claude Sonnet 5.5 in a fresh `claude -p` session. **It holds
context, method and execution independence only; model independence is not held**, as the
author is Claude Sonnet 5.5 too. The Operator's mandate for the task was Claude Sonnet 5.5
and Devin SWE-2 Max only. The Devin SWE-2 Max attempt reached the 15-minute limit with no
record (exit 124) and is discarded, with no verdict credited; the OpenCode attempts failed on
authentication or gave no assessment. The ladder stopped there (comment 6210070305). The
notes: **1** the table of budget in `test-runs.md` says 490 lines and the head was 510 — kept
as the evidence of its moment, superseded by the figure below, because a correction would
change the reviewed head; **2** the base of the budget is cited below; **3** a long line of
`setup.md` — declined (comment 6210166565). The merge of `main` at `f12da6a` came after the
round: one shared index line, the row of REQ-007 in `PRD-0001` §12, kept on both sides
(O-219); the other shared files (`loop.md`, `traceability.md`) merged in disjoint hunks.

## Verdict

Delivered, as code and test result: the workflow step, the overlay mode and their tests. The
full integration run at `4e016b7` exited 0 (06:25:05Z to 06:37:46Z, 21 packages ok,
[`test-runs.md`](../../runs/T-hz4m/test-runs.md)); it is not repeated after the merge of
`f12da6a`, which brings no Go that this change touches. At the merged head `go build`, `go
vet`, `go test ./...` and `go test -tags=integration ./internal/catalog/` exit 0, and the
five checks exit 0. **Budget:** the maximum is 650 lines over 14 files, read against current
`main` at landing (the Operator's answer, comment 6209416363, recorded as 6209446202 on #355
and 6209446809 on #354, to the question 6209355485). The diff against `main` (`f12da6a`),
close-out included, is 587 lines over 14 files; at `511cc4e` it was 510 over 11. Against
`6bc1051` the branch is larger, as it carries rows 62a, 63b and 64.

## Resource record

Recorded, not budgeted (ADR-0007). UTC; "not reported" is a value the harness did not give.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Sonnet 5.5 | not reported | not reported | start not reported, to 05:30 on 2026-10-16 |
| The plan review | reasoning | GPT-6 Sol on Devin CLI | not reported | not reported | start not reported, to 05:54 |
| The implementation (T and D) | execution | Claude Sonnet 5.5 | `high` | input 126, output 65,893 (thinking 30,338 of it), cache read 7,388,824, cache creation 153,484; cost 2.0120004 USD | 1,272,732 ms (21 min 13 s), from 05:56 |
| The merge of `d24fae3`, the full integration run, the evidence | execution | Claude Sonnet 5.5 | `high` | not reported | from 06:19 to the comment of 06:39 |
| The discarded Devin SWE-2 Max review attempt | reasoning | SWE-2 Max on Devin CLI | not reported | not reported | 15 min, exit 124, no record |
| Round 1 | reasoning | Claude Sonnet 5.5 (model independence not held) | not reported | not reported | not reported, ends 07:11 |
| The answer to round 1, the merge of `f12da6a` and this close-out | execution | Claude Sonnet 5.5 | `high` | not reported | from about 07:11 on 2026-10-16 |
