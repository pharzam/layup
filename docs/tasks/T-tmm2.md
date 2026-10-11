# T-tmm2 — row 47c of the plan, the answers to the follow-up

Issue: [#225](https://github.com/pharzam/layup/issues/225), row 47c of [the tasks
of M2c](../plan/README.md#the-tasks-of-m2c), opened by `T-2r8u` (#211). Serves
`F-0003#41` through REQ-001. Base `863c3e1` (the merge of #257, row 47b); the
budget is read against it, the base of this branch's own work. Author: Claude
Opus 5.5 on Claude Code. Evidence: [`runs/T-tmm2/`](../../runs/T-tmm2/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #225. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `863c3e1`, 7 min 30 s) gave `approve-with-conditions`: Budget maximum
1,100 lines added plus removed over 18 files, close-out inside; Cycle cap 1; no
panel. Its two conditions are applied: the integration test runs the
replacement of a written value of `budget.tsv`; `records.md` says that
`intake/answers.tsv` holds the answers to the follow-up too. The goal count is
1, final by O-196.

## What was done

1. **The step `follow-up-answers`** after `follow-up`: no follow-up row, Intake
   is done; "done before" when a follow-up has an answer; else the comments
   copied (but the App's own), the answers to the follow-up by the rule of the
   step `answers`, the wait.
2. **The records**, in one commit: the answers to the follow-up in
   `intake/answers.tsv`, the follow-up rows `answered`, the intent questions
   still missing or unparsed `open`, and `budget.tsv` written for the first time
   or a value of it replaced, its source the follow-up; the detail says that
   Intake is done.
3. **The documents:** `intake.md`, `records.md` (the layout row of
   `intake/answers.tsv`), `run.md` (the block), the Job cell of `internal/run`,
   `traceability.md`, the Task and Test cells of REQ-001 and a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-tmm2/test-runs.md): red before the code;
13 of 13 mutations caught; green at the three levels.
