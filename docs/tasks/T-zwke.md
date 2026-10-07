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

## Review rounds

The records and the Fixes reply are comments on #124. Round 1 (Claude Fable 5.1;
`7704895`, cycle 0): `material`, one finding (row 11 of "Input states" was in no
row), fixed with notes 2 to 7; the bodies of #126, #130, #131 and #132 were edited
by the App to match the generator. Round 2 (Fable; `bc011e7`, cycle 1): `nothing
material in scope`, four notes; notes 1, 2 and 4 are applied in the close-out
(the Job cell of `internal/records` in the Task cell of row 21, each flag of "The
command" in row 26, the hosts sentence), with #126 and #131 edited again; note 3
needs no change (the red record is dated and true for its list). The Python cache
of a dry run was committed by mistake and removed before the freeze.
`review-record-lint` on the comments of #124 gives `OK  6 comments; 2 round(s);
cap 1`.

## Verdict

Delivered: [the tasks of M2a](../plan/README.md#the-tasks-of-m2a), rows 21 to
27, each with its issue (#126 to #132), and the seven backlog lines. The review
ended by decay at cycle 1. The diff against `add60ef` is inside 650 lines over 10
files.

Next: rows 21, 22 and 23 (#126, #127, #128) can start at once. Row 27 needs the
Operator's inputs, which #132 lists.

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC, 2026-10-07; tokens are the
`result` event of the Claude Code CLI (stream runs); `not reported` otherwise.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | max | not reported | 10:14 to 10:17 |
| The plan review | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 952,422 (USD 6.13) | 12 min 1 s, from 10:18 |
| The answer; the parts list and the checks, test first; the issues; the table | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 10:30 to 10:35 |
| Round 1 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 880,223 (USD 5.49) | 8 min 21 s, from 10:35 |
| The fix of round 1 | execution | Claude Opus 5.5 | max | not reported | 10:44 to 10:46 |
| Round 2 | reasoning | the same as round 1 | `xhigh` | 933,447 (USD 5.30) | 7 min 42 s, from 10:46 |
| The close-out, with notes 1, 2 and 4 of round 2 | execution | Claude Opus 5.5 | max | not reported | 10:55 to 11:00 |
