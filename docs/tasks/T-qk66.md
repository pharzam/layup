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

## Review rounds

The record is a comment on #217. Round 1 (Claude Fable 5.1, effort `xhigh`, a
fresh read-only session in a clone at `4826c2f`, cycle 0, the cap): `nothing
material in scope`, five notes. Notes 1 to 3 (the sentence of `M2a` in step 6 of
`run.md`; "before `phase`" in `session.md` and in the comment of `probeStep`)
are #236, for the next change of those files. Note 4: the red of the schema of
`run-steps` was seen and not recorded in `test-runs.md`: after the code, before
the block of `run.md` changed, `TestTheTextBlocksOfRunMd` gave "schema run-steps:
column 1 (step): the type: the block has
"enum(forge|…|probe|phase)"; the Go schema has "enum(forge|…|probe|phase|rule-gaps)""
(04:24Z), then green with the block. Note 5 is a limit of the round (the e2e
level was not open to it); the author ran it green, and CI runs it.

## Verdict

Delivered: the step `rule-gaps`, the first step of Intake, in place of the row
`phase` when Start is done, with the rule "done before" and the release of the
lease after the last row. The review ended by decay at cycle 0 of cap 1. The diff
against `6f3adae`, the branch's base, is inside 650 lines over 15 files. Next:
row 42a.

## Resource record

Recorded, not budgeted (ADR-0007). UTC; reviewers' tokens are `modelUsage`;
the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 2026-10-11 04:12 to 04:15 |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,059,200 (USD 5.56) | 6 min 13 s |
| The code, test first | execution | Claude Opus 5.5 | not reported | not reported | 04:21 to 04:32 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 1,187,638 (USD 4.70) | 9 min 40 s |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 04:44 to 04:43 |
