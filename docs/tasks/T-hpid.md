# T-hpid — row 41b of the plan, the command layup spec check --facts

Issue: [#214](https://github.com/pharzam/layup/issues/214), row 41b of [the tasks
of M2c](../plan/README.md#the-tasks-of-m2c), opened by `T-2r8u` (#211). Serves
`F-0003#51` through REQ-012. Base `bec9093` (the merge of #239, row 42b); the
budget is read against it. Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-hpid/`](../../runs/T-hpid/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #214. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `bec9093`, 5 min 38 s) gave `approve-with-conditions`: Budget maximum
600 lines added plus removed over 14 files, close-out inside; Cycle cap 1; no
panel. Its two conditions are applied: each exit-2 case asserts its reason, so
the red is `unknown command "spec"`; a table that cannot be written has a unit
test. The goal count is 1, final by O-196.

## What was done

1. **`internal/cli`:** the row `spec check --facts SPANS --psb FILE` of the
   command table, and `specCheck`: both files read, a file missing, unreadable
   or not valid UTF-8 is exit 2 naming it; a span table that its block refuses
   is exit 2, the reason prefixed with the file; else the table of
   `spec.Check`, and the exit code of its column `result` (0 or 1); a table
   that cannot be written is exit 2.
2. **`internal/spec`:** `WriteTSV`, the writer of the table by `CheckSchema`,
   as `internal/cli` may not import `internal/tsv` (`TestPackageRules` refused
   the first form).
3. **The documents:** `intake.md` (the decided-here on a table that cannot be
   written), `README.md` (Progress now names `layup spec check` among the
   commands that print none; the exit code 2 row names it), `traceability.md`,
   the Task and Test cells of REQ-012 and a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-hpid/test-runs.md): red before the code;
7 of 7 mutations caught; green at the three levels.
