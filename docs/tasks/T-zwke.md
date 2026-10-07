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
files against `add60ef`, close-out inside; Cycle cap 1; no panel. **O-166** (a)
raised the cap to 2 (comment 6036688427), and the Operator approved the budget
of 700 lines over 12 files ("budget 700/12 ok", comment 6036730040). Its seven
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

The records, the Fixes replies and the decisions are comments on #124. Round 1
(Claude Fable 5.1; `7704895`, cycle 0): `material`, one finding (row 11 of "Input
states" was in no row), fixed with notes 2 to 7; the bodies of #126, #130, #131
and #132 were edited by the App to match the generator. Round 2 (Fable; `bc011e7`,
cycle 1): no material finding, four notes, applied in commit `6ce4573` (#126 and
#131 edited again). Then the CI job `security` of #134 failed: `gitleaks` read the
line `env = dict(os.environ, GH_TOKEN=token, …)` of `592c398` as a key, a false
positive. **O-166** (a): one more cycle, cap 2; the verdict of round 2 is edited to
`material`, with the reason. The fix: one line of `.gitleaksignore` with the
fingerprint, the safe form of the line, the pitfall in `docs/guardrails.md` §2;
in a clone without the file, `gitleaks` finds exactly that one leak. Round 3
(Fable; `2f1c810`, cycle 2): `nothing material in scope`, four notes: note 1 (this
record) and note 4 (the plan says "a merge of `origin/main`", not "rebases") are
applied here; note 2 came from the author's copy of the issues of round 2 in the
brief, while the bodies of #126 and #131 hold the new text; note 3 is the budget,
which the Operator approved. `review-record-lint` on the comments of #124 gives
`OK  11 comments; 3 round(s); cap 2`.

## Verdict

Delivered: [the tasks of M2a](../plan/README.md#the-tasks-of-m2a), rows 21 to
27, each with its issue (#126 to #132), and the seven backlog lines. The review
ended by decay at cycle 2 of a cap that the Operator raised once (O-166). The diff
against `add60ef` is inside 700 lines over 12 files (the Operator's approval).

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
| The first close-out; the CI finding; O-166 asked | execution | Claude Opus 5.5 | max | not reported | 10:55 to 11:00 |
| The fix of O-166 (the gitleaks runs) | execution | Claude Opus 5.5 | max | not reported | 11:09 to 11:13 |
| Round 3 | reasoning | the same as round 1 | `xhigh` | 1,707,735 (USD 6.75) | 11 min 38 s, from 11:13 |
| The close-out, with notes 1 and 4 of round 3 | execution | Claude Opus 5.5 | max | not reported | 11:25 to 11:30 |
