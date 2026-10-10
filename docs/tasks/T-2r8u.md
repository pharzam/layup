# T-2r8u — the build tasks of M2c

Issue: [#211](https://github.com/pharzam/layup/issues/211), the slicing task that
`T-g49f` (#206) named, under the Operator's decision O-164 and directive after
`M2b` (comment 6099589342 of #199). Serves `F-0003#51` and the requirements of
`M2c`. Base `b54748b` (the merge of #208). Author: Claude Opus 5.5 on Claude Code.
Evidence: [`runs/T-2r8u/`](../../runs/T-2r8u/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #211. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `b54748b`, 13 min 3 s) gave `approve-with-conditions`: Budget maximum
1,000 lines added plus removed over 17 files, close-out inside; Cycle cap 1; no
panel. Its five conditions are applied: the goal count of each row went to the
Operator before any issue opened; row 42a holds the session of a phase only, and
each row of a step adds its own name to `run-steps`; row 41a owns `spec-check`
and a second Go value of `spans`; `intake.md` has the section of the review of
the release of `M2c`; row 48 is after 42b.

**The Operator's decision:** O-196 (comment 6103234258): the 17 rows, 40 to 49b,
with the count of each final.

## What was done

1. **T:** [`parts.tsv`](../../runs/T-2r8u/parts.tsv) puts each part of the
   specification of `M2c`, and each clause of an acceptance row that the slicing
   cuts, in one row, or names it with a reason;
   [`parts.sh`](../../runs/T-2r8u/parts.sh) failed before the table and passes
   after.
2. **D1:** [The tasks of M2c](../plan/README.md#the-tasks-of-m2c), rows 40 to 49b,
   with the release of `M2c` defined; the stale sentences of the plan, the
   milestone row and the hosts of the three items.
3. **D2:** the six notes of #210, in `intake.md`, `records.md` and
   `psb-check.md`, each a "decided here" with its reason; and the section of the
   review of the release of `M2c` in `intake.md` (condition 4, decision 8 of
   #206), with the sentence on the enum of `run-steps` (condition 2).
4. **D3:** issues #212 to #228, opened by the App through
   [`gen-issues.py`](../../runs/T-2r8u/gen-issues.py), copied from `T-fdaq`'s and
   changed for `M2c`; row 49b's issue lists the Operator's inputs (R6).
   [`issues.sh`](../../runs/T-2r8u/issues.sh) failed before the issues and passes
   after (it reads the forge, so it runs by hand).
5. **The backlog** has the 17 lines, each with its issue.

**The rejected alternatives:** one task per package with no split (the rows of
`internal/run` would be one task of about 4,000 lines); one task per session
kind; the e2e inside each row; the names of every step of Intake in one row
before their code (the restart's order cannot run a step with no code).

## Review rounds

The record is a comment on #211. Round 1 (Claude Fable 5.1, effort `xhigh`, a
fresh read-only session in a clone at `87c19ea`, cycle 0, the cap): `nothing
material in scope`, seven notes. The close-out may change only the task indexes
and this file, so the notes on the specification (1, 2 and 7) are #230, for the
build rows that edit those files (41a, 44, 47a, 47c). Notes 3 to 6 are recorded
here and change nothing: the red record shows three of the five rules of
`parts.sh` (rule 3 is the precedent's code, unchanged); the five steps landed as
one commit, as the hook refused a commit of the red step on a link to the table
that did not exist yet; `gen-issues.py` says it was copied from itself (it was
copied from `runs/T-fdaq/gen-issues.py`) and `issues.sh` calls the issues "D2"
(they are D3 of the plan); and "What was done" item 3 says each note of #210 is
a "decided here", where notes 2 and 4 are plain edits.

## Verdict

Delivered: the build tasks of `M2c`, rows 40 to 49b of the plan, each with its
issue (#212 to #228), the six notes of #210, and the section of the review of
the release of `M2c`. The review ended by decay at cycle 0 of cap 1. The diff
against `b54748b`, the branch's base, is inside 1,000 lines over 17 files. Next:
the build of `M2c`, rows 40, 41a and 43a first.

## Resource record

Recorded, not budgeted (ADR-0007). UTC; reviewers' tokens are `modelUsage`;
the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 2026-10-10 20:12 to 20:16 |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,636,271 (USD 7.74) | 13 min 3 s |
| The answer, the counts to the Operator, the notes of #210 | execution | Claude Opus 5.5 | not reported | not reported | 20:29 to 20:31 |
| T, D1 to D3 and the issues, after O-196 | execution | Claude Opus 5.5 | not reported | not reported | 23:17 to 23:23 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 1,216,279 (USD 6.91) | 13 min 5 s |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 23:37 to 23:37 |
