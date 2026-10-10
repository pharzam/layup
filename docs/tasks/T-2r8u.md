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
