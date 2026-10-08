# T-esfe — row 22a of the plan, the package rules of M2a

Issue: [#127](https://github.com/pharzam/layup/issues/127), row 22 of [the tasks
of M2a](../plan/README.md#the-tasks-of-m2a), opened by `T-zwke` (#124), split
into 22a (this task) and 22b (`T-1g1q`, #137) by the Operator's decision O-170
(b). Serves `F-0003#42` through `NFR-007`. Base `2c40ac2`. Author: Claude Opus
5.5 on Claude Code. Evidence: [`runs/T-esfe/`](../../runs/T-esfe/).

**The scope of #127 before the split** (its body until O-170; the body now holds
the scope of 22a):

> ## Scope
>
> - **The parts of the specification:** `forge.md`: "The six capabilities" (the interface; the set of permissions of `M2a` and its check), the key-file checks of "The App identity"; `records.md`: the blocks `forge-register`, `harness-register`; `run.md`: rows 1 to 5 of "Input states"; `packages.md`: "The table of M2a", rule 5 by "Connects", read by `TestPackageRules` (the checker, its unit tests, the fixture `testdata/netimport`).
> - **Packages:** `internal/forge`, `internal/route`, `cmd/layup` (the test). **Requirements:** NFR-001, NFR-007. **Items:** `later-p2-run-start`.
> - **Tests:** unit; integration (the package rules); the block test (the two names to `built`).
> - **Size:** large, about 600 lines (an estimate; the plan review sets the budget). **After:** —. **Expected cap:** 1 (2 if its plan review reads the test as a gate).

## Plan and plan review

The first plan (comment 6040876369) had two plan reviews: Claude Fable 5.1 (from
15:14 UTC, 2026-10-07) stopped on a network error of the host with no record, and
was skipped (rule 4); GPT-6 Sol on the Devin CLI (comment 6045278041) gave
`reject` on the goal count (R11), with five more conditions. **O-170** (b) (the
Operator's answer in the session, comment 6053495230): row 22 splits into 22a and
22b; the Operator added two points, the specification and the plan for the new
sequence, and the goal classes of 22b. The revised plan (comment 6053505660) and
its review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `2c40ac2`, 8 min 6 s; comment 6053631596): `approve-with-conditions`,
Budget maximum 700 lines added plus removed over 14 files against `2c40ac2`,
close-out inside; Cycle cap 1 (a test file is none of the four kinds of gate); no
panel. Its four conditions and its notes are applied (the author's answer,
6053631896).

## What was done

1. **The checker** (`cmd/layup/rules_checker_test.go`, test code): `readTable`
   reads the table of phase 1 and the table of `M2a`; a row of `M2a` whose cells
   start with "(the row of phase 1" adds its "Connects" to the row of phase 1;
   `readAdapter` reads the one package of rule 5 from its line; rule 5 holds two
   checks, an own import outside that package, and a dependency that "Connects"
   does not name.
2. **The specification** (`packages.md`): the line "The one package that imports
   them: `internal/forge/github`." in rule 5, with its form; the sentence of the
   table of `M2a` names row 22a and the form of a cell "(the row of phase 1";
   the bullets of "The test of the package rules".
3. **The plan:** rows 22a and 22b; 22b after 22a, 24 after 22b, 25 after 21,
   22b, 23 and 24; the rows that can start at once. **Issue #137** (`T-1g1q`,
   row 22b) lists its goal classes, for its plan review to count.
4. **The documents:** the traceability rows, the Task and Test cells of
   `NFR-007` in `PRD-0001` §12 and a §13 line; [`docs.sh`](../../runs/T-esfe/docs.sh)
   checks them.

**Tests:** [`test-runs.md`](../../runs/T-esfe/test-runs.md): red 1, the new
tests did not compile on the checker of `main`; red 2, a stub that compiled
failed each new case for its reason (the two regression guards pass on `main`
too); green with the checker; `TestPackageRules` red on the real module until
`packages.md` held the line of rule 5, green after. `docs.sh`: nine of its ten
rules red before the edits (rule 1 was green since step 2 changed
`packages.md`), all green after.

**Known limit:** "Connects" lets a package depend on a network package by any
path; the words "through the adapter too" bind no path (note 4 of the plan
review).

**The rejected alternatives:** a list of the adapter in the checker (D7 of #79);
a new form of the cell Connects (condition 6 of the review of comment 6045278041).

## Review rounds

Round 1 (Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI with stream
output, a fresh read-only session in a clone at `637c2be`, 7 min 54 s; comment
6053992309): `nothing material in scope`, seven notes. A first run of the round
(from 06:20 UTC) stopped with no result when the Operator restarted Claude Code.
In the close-out: note 3, the traceability row of the unit tests of the checker
names `NFR-007` only, as condition 4 asked; note 5, the verdicts of `T-zwke` and
`T-8kqn` were written before O-170, which split row 22 into 22a and 22b (comment
6053495230 of #127). Declined: note 1 (the answer comment named the other case as
the guard; `test-runs.md` names the right one); note 2 (the answer said "in the
section of `NFR-007`"; the specification and the checker read the whole file, and
no result differs); note 4 (row 22a is `small` with 500 lines; the budget of the
plan review bounds it). Note 6: #137 lists the three goal classes (checked by the
author). Note 7: no change. `review-record-lint` on the comments of #127 gives
`OK  9 comments; 1 round(s); cap 1`.

## Verdict

Delivered: `TestPackageRules` reads the table of `M2a` and rule 5 by its line and
"Connects", and refuses a direct import of a network package outside the package
that rule 5 names. The review ended by decay at cycle 0. The diff against
`2c40ac2` is inside 700 lines over 14 files.

Next: row 22b (`T-1g1q`, #137), whose After cell is 22a; then row 23 (`T-xhgz`,
#128).

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC; tokens are the `result` event of
the Claude Code CLI (stream runs); `not reported` otherwise.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The first plan (row 22) | reasoning | Claude Opus 5.5 | max | not reported | 10-07 15:10 to 15:14 |
| Its plan review, first harness: skipped (a network error of the host) | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,069,120 (USD 5.08), no record | 10-07 15:14 to 15:22 |
| Its plan review | reasoning | GPT-6 Sol, Devin CLI | `xhigh` | not reported | 5 min 24 s, from 10-07 19:23 |
| The answer; O-170; the revised plan of 22a | reasoning | Claude Opus 5.5 | max | not reported | 10-07 19:29 to 19:35; 10-08 06:00 to 06:05 |
| The plan review of 22a | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,005,708 (USD 4.89) | 8 min 6 s, from 06:05 |
| The answer; the work, test first; issue #137 | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 06:14 to 06:20 |
| Round 1, first run: stopped by the restart of Claude Code | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | not reported | 06:20 to the restart |
| Round 1 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,295,487 (USD 5.42) | 7 min 54 s, from 06:28 |
| The close-out, with notes 3 and 5 | execution | Claude Opus 5.5 | max | not reported | 06:37 to 06:45 |

From 10-07 15:22 to 19:22 the host's network was down; from 10-07 19:35 to 10-08
06:00 the task waited for O-170.
