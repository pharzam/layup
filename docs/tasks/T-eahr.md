# T-eahr — row 54 of the plan, the step setup

Issue: [#277](https://github.com/pharzam/layup/issues/277), row 54 of [the tasks
of M2d](../plan/README.md#the-tasks-of-m2d), opened by `T-bng8` (#267). Serves
`F-0003#42` and `F-0001#1` through REQ-002, NFR-001 and NFR-003. Base `b28171b` (the
plan review's); the budget is read against `88a7a89`, the `main` that the branch took
last. Author: Claude Sonnet 5.5 on Claude Code. Evidence: [`runs/T-eahr/`](../../runs/T-eahr/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #277. The plan review
(GPT-6.1 Sol on Codex CLI, `high`, 8 min 54 s) gave `approve-with-conditions`: 2,200 lines
over 36 files against `b28171b`, close-out inside; cap 1; no panel. Its six conditions are
applied, each as the answer says. The Operator approved 2,850 lines over 34 files (comment 6145802484).

## What was done

The step `setup` of `layup run` in `internal/run`: the answers of the setup, the work area
`work/OWNER/NAME` built again from the records, and the loop of the step runner over its stops
up to the one stop `O-verify` of S15, with one `product-owner` session for the `F-` rows and the
removal of its files on a fail; `ReadOutputFiles` in `internal/session`, the removal and reading
of files in the records store, and `git.AddRemote`. The specification, the traceability and the
PRD cells are updated, and `guardrails.md` §2 has two pitfalls.
**Tests:** [`test-runs.md`](../../runs/T-eahr/test-runs.md): red for each slice and each fix; 118
mutants, each alone on a copy of the tree, all caught in the end; the whole gate on the head.
**Known limits:** the prose session is tested with stand-in harnesses; a real pair of
`product-owner` first answers a stop in the demo of row 59b (`T-sdtb`). A session whose records
commit failed gets a removal that finds nothing: the step is fail all the same, with both reasons.

## Review rounds

All rounds: GPT-6.1 Sol on Codex CLI (`xhigh`; ADR-0026 does not name it, the Operator allowed it).
Round 1 (`79850cd`, cycle 0; comment 6146226165): `material`, five findings, all in the change: a session
that ends after its commit left its files in the records; the output reader and two removals followed a
link; a filesystem test sat at the unit level; the records summary dropped an exception. All fixed.
Round 2 (`c70e60b`, cycle 1 of the cap of 1; comment 6147313832): the five settled; one new `material`
finding: the prompt forbade the markers that S14 requires an adapted file to keep. The Operator chose one more
cycle (comment 6147366674, cap 2) and the fix is `52bd8c7`.
Round 3 (`f2f6389`, cycle 2 of the cap of 2; comment 6147951880): `nothing material in scope`, no finding: the fix of round 2 and the five of round 1 hold.

## Verdict

Delivered: the step `setup` of `layup run`, with its tests, documents and runs. The review ended at
`nothing material in scope` after three rounds, at the cap of 2 that the Operator chose. The diff against
`88a7a89` is 2,837 lines over 32 files with the close-out inside, against the Operator's 2,850 over 34.
Next: `T-79zk` (the step verify), which runs the runner again from the stop `O-verify`.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-13, from the clock; the tokens of GPT-6.1 Sol are the
"tokens used" of Codex CLI; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning (not reached) | Claude Sonnet 5.5 | `max` | not reported | 10:55 to 11:04 |
| The plan review | reasoning | GPT-6.1 Sol, Codex CLI | `high` | 211,907 | 8 min 54 s, 11:04:21 |
| The tests, the code, the mutations, the documents | execution | Claude Sonnet 5.5 | `max` | not reported | 11:15 to 12:44 |
| Round 1 | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 247,150 | 17 min 37 s, 12:44:13 |
| The fixes of round 1, the runs | execution | Claude Sonnet 5.5 | `max` | not reported | 13:06 to 13:47 |
| Round 2 | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 271,595 | 12 min 48 s, 13:47:44 |
| The fix of round 2, the runs | execution | Claude Sonnet 5.5 | `max` | not reported | 14:04 to 14:21 |
| Round 3 | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 231,083 | 11 min 13 s, 14:21:14 |
| The notes and the close-out | execution | Claude Sonnet 5.5 | `max` | not reported | 14:33 to 14:33 |
