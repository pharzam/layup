# T-fdaq — the build tasks of M2b

Issue: [#152](https://github.com/pharzam/layup/issues/152), the slicing task that
`T-ywk7` (#147) named, under the Operator's decision O-164. Serves `F-0003#51`
and the requirements of `M2b`. Base `faeb1dd` (the merge of #149). Author: Claude
Opus 5.5 on Claude Code. Evidence: [`runs/T-fdaq/`](../../runs/T-fdaq/).

## Plan and plan review

The plan (R12, comment 6079120704), its review and the author's answer
(6079317854) are comments on #152. The plan review (Claude Fable 5.1, effort
`xhigh`, on the Claude Code CLI with stream output, a fresh read-only session in
a clone at `faeb1dd`, 10 min 54 s; comment 6079310864) gave
`approve-with-conditions`: Budget maximum 750 lines added plus removed over 14
files against `faeb1dd`, close-out inside; Cycle cap 1; no panel. Its three
conditions are applied: the honest goal count of each row went to the Operator
again; `parts.tsv` holds each clause of an acceptance row that the slicing cuts,
and `parts.sh` checks the rows of "The acceptance tests of M2b"; the duty of the
build tasks for `release-check.sh` is dropped. Its notes are applied as the
answer says.

**The Operator's decisions:** O-186 (comment 6079171344) split rows 30, 33, 34
and 36; O-187 (6079376251), after condition 1, split rows 37 and 39 too, and made
the goal count of each row final.

## What was done

1. **D1:** [The tasks of M2b](../plan/README.md#the-tasks-of-m2b), rows 28 to
   39b, one per package boundary of the specification of `M2b`, in the order of
   its imports, with the goal count of each row in its issue; the stale
   sentences of the plan, the milestone row, the item hosts and the row of #148.
   [`parts.tsv`](../../runs/T-fdaq/parts.tsv) puts each part of the
   specification, and each clause of an acceptance row that the slicing cuts, in
   one row, or names it with a reason; [`parts.sh`](../../runs/T-fdaq/parts.sh)
   failed before the table and passes after.
2. **Two moves in the specification:** the extension of `TestPackageRules` to
   the table of `M2b` is row 29's, before any package of that table exists (as
   O-170 did for `M2a`), in `session.md`, `packages.md` and `gate.md`; the review
   of the release of `M2b` is row 39a's, before the demo (O-187), in `session.md`.
3. **D2:** issues #153 to #170, opened by the App through
   [`gen-issues.py`](../../runs/T-fdaq/gen-issues.py), copied from `T-zwke`'s and
   changed for `M2b` (it reads the parts of each row from `parts.tsv`, and
   renders the table from the same rows); row 39b's issue lists the Operator's
   inputs (R6), and row 39a's takes #148.
   [`issues.sh`](../../runs/T-fdaq/issues.sh) failed before the issues and passes
   after (it reads the forge, so it runs by hand).
4. **The backlog** has the 18 lines, each with its issue; **the lesson** "A
   question whose summary its own table contradicts" is in
   [`guardrails.md`](../guardrails.md) §2.

**The rejected alternatives:** one task per package without the splits (rows 33
to 38 as two tasks, of about 1,800 and 2,400 lines); the package rules inside
the first task of `internal/session`, as the specification said; the step
`probe` inside the records of a task session; each row cut to one class (about
34 rows, each with its own gate); a task of its own for #148 now, before the
release exists.

## Review rounds

The records and the Fixes reply are comments on #152. Round 1 (Claude Fable
5.1; `95962fc`, cycle 0): `material`, three findings (the review of the release
in two rows, `setup-check` failing on an orphan task file, row 34b not after 28),
fixed in `43dd1b6` with notes 4 and 7 to 10; note 5 declined (a fifteenth file);
#159, #162, #164 and #169 written again. Round 2 (Fable; `43dd1b6`, cycle 1):
`nothing material in scope`, five notes: note 2 (which refusals row 36a builds)
goes to that row's plan review; note 3 (a Fact cell that names one In-Scope fact
of the row, not each) is declined for the budget, as the cell drives no order or
test; notes 1, 4 and 5 need no change.

## Verdict

Delivered: [the tasks of M2b](../plan/README.md#the-tasks-of-m2b), rows 28 to
39b, each with its issue (#153 to #170), and the 18 backlog lines. The review
ended by decay at cycle 1, the cap. The diff against `faeb1dd` is inside 750
lines over 14 files. Next: rows 28, 29, 30a and 31 (#153, #154, #155, #157).

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC, 2026-10-09; the reviewers'
tokens are `modelUsage` of the CLI's `result` event; the author's are not
reported (the author's effort setting is not readable from the session).

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan, the goal count, O-186 | reasoning | Claude Opus 5.5 | not reported | not reported | 10:20 to 10:40 |
| The plan review | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,015,815 (USD 5.81) | 10 min 54 s, from 10:33 |
| The answer, O-187; the checks, test first; the issues; the table | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | not reported | not reported | 10:44 to 10:59 |
| Round 1 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,645,805 (USD 7.54) | 13 min 27 s, from 10:59 |
| The fix of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 11:13 to 11:17 |
| Round 2 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,079,291 (USD 6.66) | 10 min 45 s, from 11:17 |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 11:29 to 11:35 |
