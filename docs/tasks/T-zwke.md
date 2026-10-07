# T-zwke — the build tasks of M2a

Issue: [#124](https://github.com/pharzam/layup/issues/124), opened by `T-zck8`
(#123) by the Operator's decision O-164 (b). Serves `F-0003#51` and the
requirements of `M2a`. Base `add60ef` (the merge of #125). Author: Claude Opus 5.5
on Claude Code. Evidence: [`runs/T-zwke/`](../../runs/T-zwke/).

## Plan and plan review

The plan (R12, comment 6035858403), its review and the author's answer
(6036058952) are comments on #124. The plan review (Claude Fable 5.1, effort
`xhigh`, on the Claude Code CLI with stream output, a fresh read-only session in
a clone at `add60ef`, 12 min 1 s; comment 6036058630) gave
`approve-with-conditions`: Budget maximum 650 lines added plus removed over 10
files against `add60ef`, close-out inside; Cycle cap 1; no panel. Its seven
conditions are applied, and its notes except 9 (the traceability rows: each build
task adds its own) and 12 (one row for `internal/run`, as its steps and the lease
share one integration test).

## What was done

1. **D1:** [The tasks of M2a](../plan/README.md#the-tasks-of-m2a), rows 21 to 27,
   one per package boundary of the table of `M2a`, in the order of its imports,
   with one outcome per demo; the stale sentences of the plan, the milestone row
   and the host rows. [`parts.tsv`](../../runs/T-zwke/parts.tsv) puts each part of
   the specification in one row, or names it with a reason;
   [`parts.sh`](../../runs/T-zwke/parts.sh) failed before the table and passes
   after.
2. **D2:** issues #126 to #132, opened by the App through
   [`gen-issues.py`](../../runs/T-zwke/gen-issues.py), which also renders the
   table from the same rows; row 27's issue lists the Operator's inputs (R6).
   [`issues.sh`](../../runs/T-zwke/issues.sh) failed before the issues and passes
   after (it reads the forge, so it runs by hand).
3. **The backlog** has the seven lines, each with its issue.

**The rejected alternatives:** one task for all of `M2a`; a task per section of
`run.md`; the uat inside row 26; the copy and decision rules in `internal/records`
(the specification gives them to `internal/run`); the change of
`TestPackageRules` in row 24 (row 22 makes the first package of the table).
