# T-tetz — row 50 of the plan, the records of the Scaffold

Issue: [#270](https://github.com/pharzam/layup/issues/270), row 50 of [the tasks
of M2d](../plan/README.md#the-tasks-of-m2d), opened by `T-bng8` (#267). Serves
`F-0001#1` and `F-0003#42` through NFR-001, REQ-002 and REQ-003. Base `9b4860e`
(the merge of #301); the budget is read against it. Author: Claude Sonnet 5.5 on
Claude Code. Evidence: [`runs/T-tetz/`](../../runs/T-tetz/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #270. The plan
review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at
`9b4860e`, 7 min 29 s) gave `approve-with-conditions`: Budget maximum 1,000 lines
added plus removed over 18 files, close-out inside; Cycle cap 1; no panel. Its two
conditions are applied: every rule that the plan reads beyond the words of a block is
written into the block (seven, in `records.md`), and the form of a run ID has one home
(`records.ValidRunID`, beside `CheckLease`). The goal count is 2, final by O-209 of
#267.

## What was done

1. **The schemas:** `MarkersSchema` in `internal/records/scaffold.go`, and
   `SetupTreeSchema` and `ForgeCheckSchema` in `internal/run/scaffold.go`, compared
   with their blocks by `TestTheSchemasEqualTheirBlocks` and
   `TestTheSchemasOfTheScaffoldEqualTheirBlocks`; the three names moved from
   `notYetBuilt` to `built`, which is now empty.
2. **The rules:** a check of one row and a check of the table for each block, with a
   reader that refuses a row or a table that breaks a rule, naming its line, its column
   and its reason. `markers`: the form of the `id`, a line from 1, the two angle quotes,
   `start` and `end` exactly when `value`, `end` above `start`, the session's form, and
   the order by file and then line. `setup-tree`: `root` exactly for order 0, else
   `S04` to `S14`, and the order 0, 1, 2, … from the root. `forge-check`: `detail` by
   the result, the form of the run, and the eight checks in their order.
3. **The documents:** the seven readings in the words of the blocks (`records.md`), the
   two Job cells of `packages.md`, the package comment of `internal/records`, the rows
   of `traceability.md`, the Task and Test cells of NFR-001, REQ-002 and REQ-003 and a
   §13 line of `PRD-0001`, and two lessons in `guardrails.md` §2.

**Tests:** [`test-runs.md`](../../runs/T-tetz/test-runs.md): red before the code
(`undefined`); each rule of a row or a table mutated alone, 46 of 46 caught (a first run
of 37 left 5 alive, the empty value of a column whose form refuses it too, and the tests
now assert the reason); every check of the gate exits 0 on the code.

**Known limits:** the `id` of `markers` is read by its form and not computed again,
because the row holds the marker after the field rule (a tab or a carriage return is one
space); row 53 computes it from the scan with `setup.MarkerID`. The readers read the
committed form: the file that the session reads, `input/markers.tsv`, has `session` as
`—`, and is row 53's. What reads another file or the forge is the writer's: that the
last row of `setup-tree` is the head of `layup-setup` and that each commit has its tree
(row 55), that the time of `forge-check` rises and that each row has the run of the table
(row 57a), and the quote of a proposal in the problem statement (row 53). The step of a
setup commit is `S04` to `S14` as the block says, though S10 makes no commit; a tighter
list is a change of the specification for row 55.

**The rejected alternatives:** a reader that derives its rules from the words of its
block; the three blocks in `internal/records`; the hash of the `id` in the row; a
`forge-check` table that may hold a prefix of its eight rows; an order of the steps in
`setup-tree`; a second rule for the name of a check in the row (an unknown name was
refused twice at one column, so a mutant of either lived); a second expression of the
run form in `internal/run`.

## Review rounds

Round 1 (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at
`89917d2`, cycle 0, the cap of 1; its record is comment 6138526938 on #270): `nothing
material in scope`, four notes, each in the change. The `tests` job of CI passed on that
head.

- **Applied** in `c0329f1`, a commit that no round read (a comment, a cell of the PRD,
  two rows of `traceability.md`, a Job cell and the green record): note 1 (the tests of
  `setup-tree` are traced to REQ-002 too), note 2 (the comment of
  `internal/run/scaffold.go` names rows 55 and 57a as the writers), note 3 (the Job cell
  of `internal/records` names the specification task `T-msrd`) and note 4 (the green
  record says that its runs read the working tree whose files are those of `a89030b`).

## Verdict

Delivered: the records of the Scaffold of `M2d`. The Go schemas of `markers`
(`internal/records`), `setup-tree` and `forge-check` (`internal/run`) equal their blocks,
and each reader refuses a row or a table that breaks a rule of its block, with seven
readings written into the words of the blocks. The tests were red first, each rule was
mutated alone and caught (46 of 46), and every check of the gate exits 0. The review
ended by decay at cycle 0 of cap 1. The diff against `9b4860e` is inside 1,000 lines
over 18 files (889 lines added plus removed, measured with the close-out before the
commit). Next: row 53, `T-4c5c`, writes the step `markers` on these records, and rows
51a, 51b and 52a start at once.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-13; reviewers' tokens are `modelUsage`;
the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Sonnet 5.5 | `max` | not reported | 03:53 to 04:05 (the claim to the post) |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 922,528 (USD 5.15) | 7 min 29 s, 04:05 to 04:13 |
| The tests and the code, test first | execution | Claude Sonnet 5.5 | `max` | not reported | 04:15 to 04:35 (the answer to the pull request) |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 1,069,683 (USD 5.56) | 15 min 19 s, 04:38 to 04:54 |
| The notes | execution | Claude Sonnet 5.5 | `max` | not reported | 04:54 to 04:55 |
| The close-out | execution | Claude Sonnet 5.5 | `max` | not reported | 04:56, to the commit of this file |
