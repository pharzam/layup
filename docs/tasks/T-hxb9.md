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

## Review rounds

The record is a comment on #223. Round 1 (Claude Fable 5.1, effort `xhigh`, a
fresh read-only session in a clone at `86e40b8`, cycle 0, the cap): `nothing
material in scope`, six notes. Notes 1, 2, 4, 5 and 6 (the class of the spaces
trimmed, the tab rule in `records.md`, the detail of a budget not written, the
form of `decimal`, a lesson of a shared fixture) are #256; note 3, a unit test
written after its code, is in the test record.

## Verdict

Delivered: the step `answers`, the copies of the comments of the Intake issue,
the answers of the approvers from first copies after the batch, the wait, and
`budget.tsv` from the intent form; a step that waits ends the restart. The
review ended by decay at cycle 0 of cap 1. The diff against `4c90530`, the
branch's base, is inside 1,400 lines over 24 files. CI runs the three levels on
the landing head. Next: row 47b.

## Resource record

Recorded, not budgeted (ADR-0007). UTC; reviewers' tokens are `modelUsage`;
the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 2026-10-11 10:40 to 10:41 (the claim to the post) |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,734,197 (USD 7.41) | 6 min 57 s |
| The code, test first | execution | Claude Opus 5.5 | not reported | not reported | 10:49 to 11:06 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 713,823 (USD 4.39) | 10 min 44 s |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 11:17 to 11:18 |
