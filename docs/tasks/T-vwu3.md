# T-vwu3 — row 47b of the plan, the review of the answers and the follow-up

Issue: [#224](https://github.com/pharzam/layup/issues/224), row 47b of [the tasks
of M2c](../plan/README.md#the-tasks-of-m2c), opened by `T-2r8u` (#211). Serves
`F-0003#41` through REQ-001. Base `48d8023` (the merge of #255, row 47a); the
budget is read against it, the base of this branch's own work. Author: Claude
Opus 5.5 on Claude Code. Evidence: [`runs/T-vwu3/`](../../runs/T-vwu3/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #224. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `48d8023`, 6 min 16 s) gave `approve-with-conditions`: Budget maximum
1,600 lines added plus removed over 26 files, close-out inside; Cycle cap 1; no
panel. Its three conditions are applied: an `open` intent field is asked again
whether or not `budget.tsv` is written, and this step replaces no value of it;
the review's input is `answers.tsv`; the column `start` of `questions` names the
follow-up among the rows with no quote. The goal count is 1, final by O-196.

## What was done

1. **The step `answers-review`** after `answers`: one session of
   `systems-architect` on `questions.tsv` and `answers.tsv`, its file the block
   `answers-review` without `session`, each answered question once; the record
   and the states of `questions.tsv` in one commit; the record with no row when
   no question is answered.
2. **The step `follow-up`** after `answers-review`: one row per question asked
   again (still `asked`, `open`, or an intent field that does not parse), one
   records commit, then one comment of its own header in the form of the batch;
   "done before", or the comment posted from the record; nothing written when
   nothing is left to ask.
3. **The comment writer** of the batch is shared with the follow-up.
4. **The documents:** `intake.md`, `records.md` (`start` of `questions`),
   `run.md` (the block), the Job cell of `internal/run`, `traceability.md`, the
   Task and Test cells of REQ-001 and a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-vwu3/test-runs.md): red before the code;
12 of 12 mutations caught; green at the three levels.
