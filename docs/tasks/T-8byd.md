# T-8byd — row 49a of the plan, the review of the release of M2c

Issue: [#227](https://github.com/pharzam/layup/issues/227), row 49a of [the tasks
of M2c](../plan/README.md#the-tasks-of-m2c), opened by `T-2r8u` (#211). Serves
`F-0003#51` through REQ-015 and REQ-017. Base `0b881fa` (the merge of #261, row
48, the release of `M2c`); the budget is read against it. Author: Claude Opus
5.5 on Claude Code. Evidence: [`runs/T-8byd/`](../../runs/T-8byd/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #227. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `0b881fa`, 6 min 37 s) gave `approve-with-conditions`: Budget maximum
500 lines added plus removed over 11 files, close-out inside; Cycle cap 1; no
panel. Its three conditions are applied: the task serves `F-0003#51`; the
reading of criterion 2 is `T-4tjy`'s (a finding that names a defect is fixed on
the path or opens an issue; a note that its reviewer closes with "no change"
names no defect, so it is neither); criterion 3 is run on the frozen head. The
goal count is 2, final by O-196.

## What was done

1. **The release of `M2c`** (D1): the code of `main` that the demo runs, at
   `0b881fa`, the merge of row 48; it holds phase 1, `M2a` and `M2b`. The plan's
   paragraph of row 49a names the commit.
2. **The check** (D3): `runs/T-efmy/release-check.sh` passes at `0b881fa` with
   no adaptation, as `M2c` starts no new program and adds no verb of `git` that
   reaches a remote: [`release-check.txt`](../../runs/T-8byd/release-check.txt).
   Two red runs on mutated scratch clones show that it still bites.
3. **The release review** (D4, the uat):
   [`release-review.md`](../../runs/T-8byd/release-review.md), posted on #227: a
   fresh session of Claude Fable 5.1, a model whose sessions wrote no commit of
   the release, records `holds` for REQ-015 and for REQ-017, with eighteen
   findings; none names a defect, so none is fixed or opens an issue (the
   reading of criterion 2).
4. **The documents** (D5): the plan (the release of `M2c`), the Task and Test
   cells of REQ-015 and REQ-017 and a §13 line of `PRD-0001`, a uat row of
   `traceability.md`.

**Tests:** [`test-runs.md`](../../runs/T-8byd/test-runs.md).

## Review rounds

The record is a comment on #227. Round 1 (Claude Fable 5.1, effort `xhigh`, a
fresh read-only session in a clone at `15130c2`, cycle 0, the cap): `nothing
material in scope`, three notes. Notes 1 and 3 (a line of the plan not wrapped;
the §12 cells naming no output of the check for `M2c`) are #264; note 2 (the
green run came before the red runs) is in the test record, with its reason.

## Verdict

Delivered: the review of the release of `M2c`: `holds` for REQ-015 and REQ-017
over the code of `main` at `0b881fa` (phase 1, `M2a`, `M2b` and `M2c`), by the
check, unchanged, and by a review of a model that wrote none of it. The review
ended by decay at cycle 0 of cap 1. The diff against `0b881fa` is inside 500
lines over 11 files. Next: row 49b, the demo, whose inputs the Operator gives
(comment 6109884732 on #228).

## Resource record

Recorded, not budgeted (ADR-0007). UTC; reviewers' tokens are `modelUsage`;
the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 2026-10-11 13:45 to 13:46 (the claim to the post) |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 790,608 (USD 4.50) | 6 min 37 s |
| The check and the red runs | execution | Claude Opus 5.5 | not reported | not reported | 13:46 to 13:54 |
| The release review (the uat) | reasoning | Claude Fable 5.1 | `xhigh` | 2,185,023 (USD 10.31) | 8 min 39 s |
| The documents | execution | Claude Opus 5.5 | not reported | not reported | 14:03 to 14:08 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 611,564 (USD 4.11) | 9 min 10 s |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 14:18 to 14:18 |
