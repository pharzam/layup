# T-2yw7 — the command frame and the end-to-end harness

Issue: [#80](https://github.com/pharzam/layup/issues/80), row 3 of the
[implementation plan](../plan/README.md); child of `layup gate` (#33). Serves
`F-0003#44`. Base `0506c13` (the merge of the PDR, #102). Author: Claude Opus
5.5 on Claude Code. Evidence: [`runs/T-2yw7/`](../../runs/T-2yw7/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #80. The plan
review (Claude Fable 5.1 on Claude Code, a fresh session) gave
`approve-with-conditions`: Budget maximum 1,500 lines added plus removed over 24
files against the base, close-out inside; Cycle cap 1. The author applied the
condition (the repeat helper compares the exit code and the standard output,
never the standard error) and the eight notes.

**A correction of the answer:** note 1 of the plan review, and the answer, said
that row 5 decides a manifest with no row under K19. K19 is a pending kind with a
missing tool, so `docs/spec/README.md` names row 5 for that question and no
defect number.

## What was done

1. **Test first** ([`test-runs.md`](../../runs/T-2yw7/test-runs.md)): the unit
   tests of the parser, the usage, the exit map and the progress lines; the byte
   compare of the harness and the scan of the input rule; the end-to-end
   scenarios. Each one failed on its own assertion against skeletons of the new
   names.
2. **`internal/cli`:** one command table, with the parser of the argument rules
   and the usage that come from it (`args.go`); the exit codes and the exit map
   (`cli.go`); the progress lines with their beats (`progress.go`).
   `layup version extra` is now a usage error (K34).
3. **`cmd/layup`:** the end-to-end harness (`harness_e2e_test.go`): `TestMain`
   builds the binary once, a run helper with a fixed environment that a test can
   change by entry, and a repeat helper; the scenarios `TestVersion` and
   `TestUsageErrors` (`usage_e2e_test.go`), which replace `main_e2e_test.go`; the
   scan of the input rule (in `rules_checker_test.go`) with its unit cases (in
   `rules_test.go`) and `TestInputRule` on the real module and on the fixture
   module `testdata/envread`.
4. **`docs/spec/README.md`, "Commands":** D1 to D7 and the condition as values
   decided here; K29 (the progress lines) and K34 are settled.
5. **The tests' rows:** `docs/tests/traceability.md` (the tests of the frame
   that prove no requirement have `—` in Covers, note 4); the `PRD-0001` §12
   Test cells of `REQ-001`, `NFR-004` and `NFR-005`, with a §13 row.

**What this task does not hold** (D9): the helper that makes a repository from
a fixture (row 7, `verify-test-baseline`), the scenario of a `not-active` row
with a changed `PATH` (row 5), and the scenarios of `layup psb check` (row 6).
`tests/` gets no fixture, so `tests/README.md` does not change.

**The rejected alternatives:** Go's `flag` package (it stops at the first
positional argument, and `layup gate REPO --base REV --head REV` puts flags after
it); a status line that redraws itself on standard error (it does not work in a
CI log); a second usage text per command (it can differ from the dispatch); a
fixture repository helper in this task (no scenario of this task needs one, and
the inventory gives it to row 7).
