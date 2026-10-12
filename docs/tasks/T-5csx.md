# T-5csx — the acceptance test of the merge into the batch and the comment

Issue: [#287](https://github.com/pharzam/layup/issues/287), a gap that round 6 of
`T-bng8` (#267) revealed, off its path (gate step 5). Serves `F-0003#51` through
`REQ-012`, as `T-bng8` does. Base `45f546c` (the merge of #285). Author: Claude Sonnet
5.5 on Claude Code. Evidence: [`runs/T-5csx/`](../../runs/T-5csx/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #287. The plan
review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at
`45f546c`, 6 min 27 s) gave `approve-with-conditions`: Budget maximum 160 lines added
plus removed over 7 files, close-out inside; Cycle cap 1; no panel. Its one condition
is applied: the plan names the fact it serves, `F-0003#51` through `REQ-012`, and this
header does. The notes are applied: the plan takes the second reading of criterion 2;
the mutants are recorded as they came (one line, then two); D1 changes the Task cell of
row 53 alone; the row is in the words of the section, and `parts.sh` is unchanged.

## What was done

1. **T:** the row "The setup questions in the batch and the comment" in "The acceptance
   tests of M2d" (`scaffold.md`), at the level `integration` that row 53 gives its Tests
   cell, in the words of "The answers of the setup". `sh runs/T-bng8/parts.sh` went red
   with the one line of check 4.
2. **D1:** the line in [`parts.tsv`](../../runs/T-bng8/parts.tsv), the test as a part of
   row 53, and the Task cell of row 53 in the plan, from `gen-issues.py --table`. The
   check reads the map and the plan, and is green after them.
3. **D2:** the body of #276, written again from the map by `gen-issues.py --edit 53`: one
   line differs, the list of parts. No check reads it.
4. **M:** two mutants on copies of the tree: the line dropped from the map gives the line
   of check 4; the part dropped from the Task cell gives the two lines of check 6.
5. **The lesson (gate step 7):** none owed. The trap that this task closes, a clause of
   the map that no test of the specification holds, is round 6's finding of #267, and the
   lessons of `T-bng8` are #288.

**Criterion 2 of #287** ("the coverage map ... is updated") is read as the map and what
is written from it: the Task cell of row 53 and the body of #276, the two derived edits.

## Review rounds

Round 1 (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at
`90d21db`, cycle 0, the cap of 1; its record is a comment on #287): `nothing material in
scope`, two notes. Note 1, on the header of `test-runs.md` ("green after the map, the
plan and the issue of row 53 name it"), is declined: the close-out may change only the
task indexes and this file, and the sentence is true, as the check is green once the map
and the plan name the test; this file says that the check reads those two and that #276 is
written from the map. Note 2 asks the close-out to carry the fact, the two derived edits,
the lesson sentence and the figure at landing; this file does.

## Verdict

Delivered: the acceptance test of the merge of the setup questions into the batch and of
the comment, in "The acceptance tests of M2d", registered as a part of row 53 in the map,
the plan and the body of #276. The review ended by decay at cycle 0 of cap 1. The diff
against `45f546c` is 144 lines over 6 files with the close-out, inside the Budget
maximum of 160 lines over 7 files. Next: the build of row 53 (`T-4c5c`, #276) delivers the
test.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-12. The reviewers' tokens are
`modelUsage`, with Claude Code's own estimate of the money; the author's are not
reported. Claude Sonnet 5.5 is the model that the Operator chose for the session
(`/model`); Model tiers names `Claude Sonnet 5` for the execution tier, so the plan, its
answer and the close-out are parts of that tier or of none, and the plan is a reasoning
part on it.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Sonnet 5.5 | `max` | not reported | 12:02 (the issue) to 12:07 (the plan) |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,005,364 (USD 4.83) | 6 min 27 s, 12:08 to 12:14 |
| The answer to the plan review | reasoning | Claude Sonnet 5.5 | `max` | not reported | 12:15 |
| T, D1, D2 and M | execution | Claude Sonnet 5.5 | `max` | not reported | 12:16 to 12:20 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 762,848 (USD 4.20) | 8 min 41 s, 12:21 to 12:30 |
| The close-out | — | Claude Sonnet 5.5 | `max` | not reported | 12:32 to 12:33 |
