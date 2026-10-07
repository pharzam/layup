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
