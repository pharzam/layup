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

## Verdict

Delivered: the frame of every command in `internal/cli` (one command table with
the parser of the argument rules and the usage, the exit map, and the progress
lines) and the end-to-end harness of `cmd/layup`, with the scenarios
`TestVersion` and `TestUsageErrors` (the demo: each usage error of the built
binary gives exit code 2) and `TestInputRule`; D1 to D7 and the condition of the
plan review in `docs/spec/README.md` ("Commands"), which settle K29 and K34;
seven traceability rows; the `PRD-0001` §12 Test cells of `REQ-001`, `NFR-004`
and `NFR-005`, with a §13 row.

The plan review (Claude Fable 5.1) gave `approve-with-conditions`; the author
applied its condition and its eight notes. Review round 1 (Claude Fable 5.1,
`a5e2552`, cycle 0; [the record](../../runs/T-2yw7/review-round-1.md)) gave
`nothing material in scope`, with 13 notes. Note 6 is applied: the evidence names
the two files whose line numbers come from an earlier state. Notes 2, 5, 7, 9,
10, 12 and 13 record checks with no defect. The others need no change: note 1
(an empty positional argument reaches the command's own input check, which gives
exit 2, as "Commands" says); note 3 (the scan finds the readers that D7 names,
and a program that `layup` starts is outside the rule, note 7 of the plan
review); note 4 (only `newProgress` makes the progress lines, with ten seconds);
note 8 (§12 names the tests that prove a requirement at the integration and
end-to-end levels, as for `TestPackageRules` of row 2; the unit tests of the
helpers are rows of `traceability.md`); note 11 (the first line of the usage is
not a documented value). At the head, `go build`, `go vet` with each tag,
`gofmt`, the three test levels, `go test -race` on `internal/cli` and all local
checks pass, and `review-record-lint` passes on the comments of #80. The diff
against `origin/main` is 1,190 lines over 23 files with the close-out, inside the
Budget maximum of 1,500 lines over 24 files.

Next: row 4 of the plan (`T-3jpx`, #81), the stack catalog package.

## Resource record

Recorded, not budgeted (ADR-0007). Times are 2026-10-02, UTC. Token counts are
`not reported` where the harness does not give them; the review sessions ran
with `--output-format stream-json`, which gives them.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | max | not reported | a draft within 06:55 to 07:01 (beside the review round of `T-4wrw`); 07:07 to 07:08 |
| The plan review | reasoning | Claude Fable 5.1 | the default of `claude -p` | 769,197 (input 258, cache write 113,625, cache read 630,553, output 24,761 of which thinking 14,541); USD 3.67 at list price | 5 min 6 s, 07:08:45 to 07:13:51 |
| The answer to the plan review | reasoning | Claude Opus 5.5 | max | not reported | 07:14 |
| The tests, the code, the specification and the records (steps 1 to 5); the freeze | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | drafts within 07:09 to 07:13; 07:14 to 07:22 |
| Review round 1 | reasoning | Claude Fable 5.1 | the default of `claude -p` | 864,340 (input 200, cache write 105,290, cache read 731,962, output 26,888 of which thinking 12,876); USD 3.64 at list price | 5 min 26 s, 07:22:58 to 07:28:24 |
| The notes of round 1 and the close-out | execution | Claude Opus 5.5 | max | not reported | 07:28 to 07:33 |
