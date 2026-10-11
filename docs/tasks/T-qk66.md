# T-qk66 — row 43a of the plan, the step rule-gaps

Issue: [#217](https://github.com/pharzam/layup/issues/217), row 43a of [the tasks
of M2c](../plan/README.md#the-tasks-of-m2c), opened by `T-2r8u` (#211). Serves
`F-0003#41` through REQ-001. Base `6f3adae` (the merge of #231, row 40); the
budget is read against it. Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-qk66/`](../../runs/T-qk66/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #217. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `6f3adae`, 6 min 13 s) gave `approve-with-conditions`: Budget maximum
650 lines added plus removed over 15 files, close-out inside; Cycle cap 1; no
panel. Its three conditions are applied: each test that expected `phase` after a
complete Start changes with the code, the e2e among them; the rows after `clone`
come from the records that `clone` read; the second restart compares the blob of
`intake/psb.tsv`. The goal count is 1, final by O-196.

## What was done

1. **The restart** runs `forge` and `clone`, then the rows that the records
   choose: with Start not done, `version`, `lease`, `probe`, `phase`, as before;
   with Start done (the two issues and the watch written), the steps of Intake in
   place of `phase`, and the lease released after the last row (a release that
   fails makes that row `fail`, its `detail` naming the lease).
2. **The step `rule-gaps`:** the table of `layup psb check` on
   `start/problem-statement.md` (`psb.Check`, `psb.WriteTSV`, what the command
   runs), committed as `intake/psb.tsv`; "done before" when the record exists.
3. **`run-steps`** gains `rule-gaps`, in the block and its Go schema.
4. **The documents:** `run.md` (the block, step 6, the progress of a restart),
   the Job cell of `internal/run`, the comment of the uat test,
   `traceability.md`, the Task and Test cells of REQ-001 and a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-qk66/test-runs.md): red before the code;
5 of 5 mutations caught; green at the three levels.

**Decided here:** the progress of `forge` and `clone` counts the rows of a
restart with Start not done, as the run plans its rows when it prints them.

**The rejected alternatives:** the steps of Start and Intake in one restart; a
row `lease` at the end of a restart.
