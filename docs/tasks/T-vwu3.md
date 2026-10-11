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
   `systems-architect` on the problem statement, `questions.tsv` and
   `answers.tsv` (the problem statement by the fix of round 1), its file the block
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

## Review rounds

The records are comments on #224. Round 1 (Claude Fable 5.1, effort `xhigh`, a
fresh read-only session in a clone at `b6f1161`, cycle 0): `material`, one
finding: the session of the review got no `problem-statement.md`, which The
sessions of Intake give each step. The fix `bf4663f` (cycle 1) gave it the file,
red first by a world's script that refuses a session of a phase with none, and
made `follow-up` fail on a missing problem statement (note 4). Round 2 (the
same reviewer and effort, a fresh session at `bf4663f`, cycle 1, the cap):
`nothing material in scope`, three notes: note 1 is applied above; notes 2 and
3, with notes 2 and 5 of round 1, are #258.

## Verdict

Delivered: the steps `answers-review` and `follow-up`: a session of
`systems-architect` closes or opens each answered question, and one follow-up
comment asks again only the missing, the unparsed and the open questions. The
review ended by decay at cycle 1 of cap 1. The diff against `48d8023`, the
branch's base, is inside 1,600 lines over 26 files. CI runs the three levels on
the landing head. Next: row 47c.

## Resource record

Recorded, not budgeted (ADR-0007). UTC; reviewers' tokens are `modelUsage`;
the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 2026-10-11 11:22 to 11:23 (the claim to the post) |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,205,699 (USD 5.96) | 6 min 16 s |
| The code, test first | execution | Claude Opus 5.5 | not reported | not reported | 11:29 to 11:47 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 1,118,099 (USD 5.66) | 12 min 45 s |
| The fix of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 12:00 to 12:06 |
| Round 2 | reasoning | Claude Fable 5.1 | `xhigh` | 1,085,657 (USD 4.83) | 9 min 19 s |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 12:16 to 12:16 |
