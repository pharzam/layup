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

## Review rounds

The records are comments on #225. Round 1 (Claude Fable 5.1, effort `xhigh`, a
fresh read-only session in a clone at `56dd2d8`, cycle 0): `material`, one
finding: with `budget.tsv` not yet written, the pair rule of `B` and `U` was
skipped when another field was still unparsed. The fix `2d7299b` (cycle 1)
holds the pair rule whatever the state of the hours, red first, and keeps a
`budget.tsv` with no `B` or `U` as it was (note 2). The Fixes reply first named
the old head, as the fix's commit had been refused by `git diff --check`; a
correction names `2d7299b`, and a round 2 started on the old head was stopped
before it wrote anything. Round 2 (the same reviewer and effort, a fresh
session at `2d7299b`, cycle 1, the cap): `nothing material in scope`, three
notes, which are #260.

## Verdict

Delivered: the step `follow-up-answers`, the last step of Intake: the answers
to the follow-up recorded, `budget.tsv` written or a value of it replaced with
the follow-up as its source, a field still unparsed `open`, and the end of
Intake. The review ended by decay at cycle 1 of cap 1. The diff against
`863c3e1`, the branch's base, is inside 1,100 lines over 18 files. CI runs the
three levels on the landing head. Next: row 48.

## Resource record

Recorded, not budgeted (ADR-0007). UTC; reviewers' tokens are `modelUsage`;
the author's are not reported. The stopped round, about a minute of a Fable
session, wrote no record and reported no figure.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 2026-10-11 12:20 to 12:21 (the claim to the post) |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 800,447 (USD 5.39) | 7 min 30 s |
| The code, test first | execution | Claude Opus 5.5 | not reported | not reported | 12:29 to 12:46 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 929,526 (USD 4.62) | 8 min 48 s |
| The fix of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 12:55 to 12:58 |
| Round 2 | reasoning | Claude Fable 5.1 | `xhigh` | 759,967 (USD 4.32) | 7 min 29 s |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 13:06 to 13:06 |
