# T-hiv9 — row 44 of the plan, the steps numbering and spans

Issue: [#219](https://github.com/pharzam/layup/issues/219), row 44 of [the tasks
of M2c](../plan/README.md#the-tasks-of-m2c), opened by `T-2r8u` (#211). Serves
`F-0003#51` through REQ-012. Base `6d531db` (the merge of #243, row 43b); the
budget is read against it, the base of this branch's own work. Author: Claude
Opus 5.5 on Claude Code. Evidence: [`runs/T-hiv9/`](../../runs/T-hiv9/).

## Plan and plan review

The plan (R12), its correction, its review and the author's answer are
comments on #219. The plan review (Claude Fable 5.1, effort `xhigh`, a fresh
read-only session in a clone at `6d531db`, 8 min 11 s) gave
`approve-with-conditions`: Budget maximum 1,400 lines added plus removed over
24 files, close-out inside; Cycle cap 1; no panel. Its four conditions are
applied: the goal count is the issue's 2; a span with no byte in the file and a
repeated `start` are refused by the block; D3 asserts that `spans` `fail` is the
last row; the writer of `spec/check.tsv` is `tsv.Write` by `spec.CheckSchema`.
The goal count is 2, final by O-196: the two steps and their records; the rerun
after a failed check.

## What was done

1. **The step `numbering`** after `meaning`: one session of `product-owner`
   whose file is the block `spans` without `text`; code adds each span's text
   from the file; "done before" unless the last `spec/check.tsv` has a `fail`
   row; with no pair, `fail` (`pair`).
2. **The step `spans`** after `numbering`: `spec.Check` in process on the
   records, its table committed as `spec/check.tsv` (kept as it is when equal to
   the last); `fail` with the first `fail` row, so no step after it runs.
3. **`run-steps`** gains `numbering` and `spans`, in the block and its Go schema.
4. **The test worlds:** the scripts of the world of `internal/run` and of the
   e2e world compute a span table that covers the problem statement, and route
   `product-owner`.
5. **The documents:** `intake.md` (the decided-here on the session's file, the
   details of `spans`, an equal table kept), `run.md` (the block), the Job cell
   of `internal/run`, `traceability.md`, the Task and Test cells of REQ-012 and
   a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-hiv9/test-runs.md): red before the code;
12 of 12 mutations caught; green at the three levels.

## Known limits

- **A span whose text is exactly `—`** reads back as the empty value, which the
  block refuses, so such a file is refused at `numbering` (note 5 of the plan
  review; row 40's limit, as row 43b met it for a quote).
- **A spare session** (note 9): a run that stops between the rerun of
  `numbering` and `spans` leaves the `fail` row, so the next restart runs
  `numbering` once more.
