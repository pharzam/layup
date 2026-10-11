# T-2et0 — row 45a of the plan, the draft

Issue: [#220](https://github.com/pharzam/layup/issues/220), row 45a of [the tasks
of M2c](../plan/README.md#the-tasks-of-m2c), opened by `T-2r8u` (#211). Serves
`F-0003#51` through REQ-012. Base `ba4d33e` (the merge of #246, row 44); the
budget is read against it, the base of this branch's own work. Author: Claude
Opus 5.5 on Claude Code. Evidence: [`runs/T-2et0/`](../../runs/T-2et0/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #220. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `ba4d33e`, 6 min 57 s) gave `approve-with-conditions`: Budget maximum
1,000 lines added plus removed over 18 files, close-out inside; Cycle cap 1; no
panel. Its two conditions are applied: `intake.md` gives the form of the
session's file as a decided-here, and its Input states row and acceptance row
name the three rules of the check. The goal count is 1, final by O-196.

## What was done

1. **The check of the draft:** each row's `covers` names at least one number,
   each a fact of `spec/spans.tsv`; each `need` and `constraint` fact is
   covered; each `measure` fact has a success criterion.
2. **`Keep` of the step:** the session's file (the block `draft` without `id`
   and `session`), checked, with `R-NNN` and the session ID, as
   `spec/draft.tsv`.
3. **The step `draft`** after `spans`: "done before" with the record; else one
   session of `product-owner` (the first admitted pair, as for the numbering) on
   the problem statement and the span table; `done` with "N rows, in
   spec/draft.tsv". `run-steps` gains `draft`.
4. **The test worlds:** the scripts answer the draft prompt with one
   requirement per `need` fact of the span table.
5. **The documents:** `intake.md` (R10: the input of the draft holds
   `spans.tsv`; the form of its file; the three rules in Input states and the
   acceptance row), `run.md` (the block), the Job cell of `internal/run`,
   `traceability.md`, the Task and Test cells of REQ-012 and a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-2et0/test-runs.md): red before the code;
11 of 11 mutations caught; green at the three levels.
