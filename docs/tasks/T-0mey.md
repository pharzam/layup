# T-0mey — row 45b of the plan, the completeness review

Issue: [#221](https://github.com/pharzam/layup/issues/221), row 45b of [the tasks
of M2c](../plan/README.md#the-tasks-of-m2c), opened by `T-2r8u` (#211). Serves
`F-0003#51` and `F-0003#52` through REQ-012 and REQ-013. Base `1b36954` (the
merge of #248, row 45a); the budget is read against it, the base of this
branch's own work. Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-0mey/`](../../runs/T-0mey/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #221. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `1b36954`, 6 min 56 s) gave `approve-with-conditions`: Budget maximum
1,000 lines added plus removed over 24 files, close-out inside; Cycle cap 1; no
panel. Its condition is applied: one detail, "not run: no admitted pair on
another harness", for every case of the step with no pair, written in
`intake.md`. The goal count is 1, final by O-196.

## What was done

1. **`PhasePair`** takes a harness to pass over (`""` for none).
2. **The step `completeness`** after `draft`: the harness of the draft's session
   is passed over; the first admitted pair of `qa-engineer` on another harness
   runs one session on the problem statement, the span table and the draft;
   with none, `fail` "not run: no admitted pair on another harness"; "done
   before" with the record. `run-steps` gains `completeness`.
3. **`Keep` of the step:** the block `coverage` without `id` and `session`, each
   `fact` a fact of the span table and each `requirement` a row of the draft,
   `C-NNN` and the session, as `intake/coverage.tsv`.
4. **The test worlds:** a second scripted harness in each (`devin` in the world
   of `internal/run`, `spare` in the e2e world), with the routes of
   `qa-engineer`, the first on the draft's harness.
5. **The documents:** `intake.md` (the harness passed over, one detail with no
   pair, the form and the references of the review's file, its row of Input
   states), `run.md` (the block), the Job cell of `internal/run`,
   `traceability.md`, the Task and Test cells of REQ-012 and REQ-013 and a §13
   line.

**Tests:** [`test-runs.md`](../../runs/T-0mey/test-runs.md): red before the code;
9 of 9 mutations caught; green at the three levels.

## Known limits

- **The harness of the numbering** is not read: the review passes over the
  draft's harness, which is the numbering's by the routing rule; a numbering
  rerun under a changed register could run on another harness (note 2 of the
  plan review).
- **A draft of no row** (note 2 of round 1): with no `need`, `constraint` or
  `measure` fact, the draft is committed with no row, so the draft names no
  session and the step `completeness` fails on each restart; the state is #251's.
