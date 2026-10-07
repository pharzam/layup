# T-8kqn — row 21 of the plan, the records of Start

Issue: [#126](https://github.com/pharzam/layup/issues/126), row 21 of [the tasks
of M2a](../plan/README.md#the-tasks-of-m2a), opened by `T-zwke` (#124). Serves
`F-0003#42` through `NFR-001`. Base `6d256e0` (the merge of #134). Author: Claude
Opus 5.5 on Claude Code. Evidence: [`runs/T-8kqn/`](../../runs/T-8kqn/).

## Plan and plan review

The plan (R12, comment 6037445286) and its review are comments on #126. Claude
Fable 5.1 (from 12:00 UTC) stopped on a network error with no record and was
skipped (rule 4). The plan review (GPT-6 Sol, effort `xhigh`, on the Devin CLI, a
fresh read-only session in a clone at `6d256e0`, 5 min 15 s; comment 6040027027)
gave `reject`, on the goal count (R11) with seven more conditions. **O-167** (a)
(the Operator's answer in the session, comment 6040085862): row 21 is one goal;
the plan stands with conditions 2 to 8, which the author's answer (6040027455)
applies. Budget maximum 900 lines added plus removed over 13 files against
`6d256e0`, close-out inside; Cycle cap 1; no panel.

## What was done

1. **D1:** `internal/records/start.go` holds `StartSchema`, `ApproversSchema`,
   `LeaseSchema` and `CopiesSchema`; `TestTheSchemasEqualTheirBlocks` compares
   them with their blocks, and their names moved from `notYetBuilt` to `built`.
2. **D2:** `CheckStart` (with the IDs of the harness register: each fixed name
   once, in the block's order; each harness's two rows together, `cap` first, in
   the register's order; the form of each value; the empty value only where the
   block allows it; the source of each name), `CheckApprover`, `CheckLease` (no
   column holds the empty value; a run ID of 16 lowercase hexadecimal
   characters), `CheckCopy` and `CheckCopies` (`seen` from 1 with no gap; the
   body path); the readers `ReadStart`, `ReadApprovers`, `ReadLease` (exactly one
   row) and `ReadCopies` run them, as `ReadStalls` does. `start_test.go` has one
   case per rule, through the readers.
3. **D3 and D4:** the Job cell of `internal/records` in the table of phase 1, the
   rows of `docs/tests/traceability.md`, the Test cell of `NFR-001` in `PRD-0001`
   §12 and a §13 line; [`docs.sh`](../../runs/T-8kqn/docs.sh) checks the three.

**Tests:** [`test-runs.md`](../../runs/T-8kqn/test-runs.md): red with no Go value
(it did not compile); red against a stub with a wrong column and readers that
check no rule (each rule case and the comparison failed); red for the 21 rule
cases of `start` once the schema was right; green after. `docs.sh` failed for
each of its eight rules before the edits and passes after.

**The rejected alternatives:** the rules in `internal/run` (the specification
gives the schemas and the row rules to `internal/records`); a two-place rule for
`intake.cap` and a character rule for `app` (the blocks give neither).

## Review rounds

Round 1 (Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI with stream
output, a fresh read-only session in a clone at `26c026c`, 6 min 54 s; comment
6040324281): `nothing material in scope`, seven notes. In the close-out: notes 1
and 2 add two cases (a harness with no rows; a register ID that is not of the form
`<word>`, which fails when that check is off, `runs/T-8kqn/test-runs.md`); note 4
adds `TestAValidStartIsRead` to the Test cell of `NFR-001`; note 7 rewords the
comment of `intakeCap`. Declined: note 3 (a table-level error of `ReadLease` has
no column; line 1 is the header, and the reason names the count); note 5 (the red
of the 21 cases of `start` is written, not captured: the stub that gave it no
longer exists, and red 2 and the line numbers of the head show it). Note 6: the
close-out stays at 13 files. `review-record-lint` on the comments of #126 gives
`OK  5 comments; 1 round(s); cap 1`.

## Verdict

Delivered: the Go schemas of `start`, `approvers`, `lease` and `copies` equal
their blocks, with the rules of each block and the readers that run them. The
review ended by decay at cycle 0. The diff against `6d256e0` is inside 900 lines
over 13 files.

Next: rows 22 and 23 (#127, #128) can start; row 24 waits for row 22, and row 25
for rows 21 to 24.

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC, 2026-10-07; tokens are the
`result` event of the Claude Code CLI (stream runs); `not reported` otherwise.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | max | not reported | 11:56 to 12:00 |
| The plan review, first harness: skipped (a network error) | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 463,537 (USD 3.57), no record | 12:00 to 12:03 |
| The plan review | reasoning | GPT-6 Sol, Devin CLI | `xhigh` | not reported | 5 min 15 s, from 14:18 |
| The answer; O-167 | reasoning | Claude Opus 5.5 | max | not reported | 14:23 to 14:26 |
| The work, test first | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 14:26 to 14:31 |
| Round 1 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 788,744 (USD 3.70) | 6 min 54 s, from 14:32 |
| The close-out, with notes 1, 2, 4 and 7 | execution | Claude Opus 5.5 | max | not reported | 14:39 to 14:45 |

From 12:03 to 14:16 the task did no work: the network of the host was down.
