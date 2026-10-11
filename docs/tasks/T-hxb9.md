# T-hxb9 — row 47a of the plan, the answers

Issue: [#223](https://github.com/pharzam/layup/issues/223), row 47a of [the tasks
of M2c](../plan/README.md#the-tasks-of-m2c), opened by `T-2r8u` (#211). Serves
`F-0003#41` through REQ-001 and NFR-003. Base `4c90530` (the merge of #253, row
46); the budget is read against it, the base of this branch's own work. Author:
Claude Opus 5.5 on Claude Code. Evidence: [`runs/T-hxb9/`](../../runs/T-hxb9/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #223. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `4c90530`, 6 min 57 s) gave `approve-with-conditions`: Budget maximum
1,400 lines added plus removed over 24 files, close-out inside; Cycle cap 1; no
panel. Its three conditions are applied: an answer comes from a first copy that
`Decision` accepts for the role of its question's block; the App's own comments
are not copied; the rule of an answer line is written in `intake.md`. The goal
count is 2, final by O-196: the wait and the answer lines; the parse of the
intent form and `budget.tsv`.

## What was done

1. **The step `answers`** after `batch`: the comments of the Intake issue
   copied (but the App's own); the answer lines of the approvers in the role of
   each question's block, from first copies after the batch, the last line of
   an ID its answer; "waiting" while a block that asks has no such comment.
2. **The records**, in one commit: `intake/answers.tsv` in the order of the
   questions, the answered states of `questions.tsv`, and `budget.tsv` when the
   intent form parses (decimals, spaces around them trimmed, `B` below `U` as
   numbers).
3. **"Waiting" stops the restart:** a step `done` with a detail that starts with
   "waiting" is the last row; the run releases the lease and exits 0.
4. **The documents:** `intake.md` (the decided-here points), `run.md` (the
   block; "waiting"), the Job cell of `internal/run`, `traceability.md`, the
   Task and Test cells of REQ-001 and NFR-003 and a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-hxb9/test-runs.md): red before the code;
17 of 17 mutations caught; green at the three levels.
