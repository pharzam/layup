# T-0p5p — one clock for each test world of `internal/run`

Issue: [#293](https://github.com/pharzam/layup/issues/293), cause 1 of the flaky `tests`
job of CI ([#291](https://github.com/pharzam/layup/issues/291)). Serves `F-0003#47`
(Verification on Every Change). `T-d6q5` (#92) met a CI-only failure of a test helper, not
a clock (the entry "A background maintenance of git", `docs/guardrails.md`). Base `1c167eb`
(the merge of #286). Author: Claude Sonnet 5.5 on Claude Code. Evidence:
[`runs/T-0p5p/`](../../runs/T-0p5p/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #293. The plan review
(Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at `1c167eb`, 9 min
45 s) gave `approve-with-conditions`: Budget maximum 260 lines added plus removed over 7
files, close-out inside; Cycle cap 1; no panel. Its three conditions are applied: the load
is one busy loop for each CPU of the host (2), the CPU count beside the counts; the hours
cap (8.0 hours at 600 times is 48 s of real time) is bounded by the time of each test of
the suite, idle and loaded; the guardrails entry and the traceability rows are written.
Its notes are applied or carried below; the issue body is corrected (one pair, not two).

## What was done

1. **T:** `TestTheClockOfAWorldDoesNotGoBack`: the clock of a second config is not before
   what a first read 20 ms (12 s of the clock) later. Red on `1c167eb`: 12:00:00.02 against
   12:00:12.27.
2. **D1:** the world has one clock (`world.clock`), made in `newWorld`, and `config` reads
   it. The test is green, and so are 100 runs of `TestAFailedSpanCheckRunsNumberingAgain`
   and the whole integration suite of `internal/run`, idle and beside a busy loop for each
   CPU; its longest test takes 11.04 s beside the loops, against the 48 s.
3. **M:** one mutant on a copy of the final tree: `config` makes its clock again, and the
   new test fails as in the red run.
4. **Step 7:** the entry "A fake clock made anew at each call of a world's helper" in
   `docs/guardrails.md`; a row for the new test in `docs/tests/traceability.md`, and
   `T-0p5p` in the row of the flaky test.

**The cause, and what is a reading.** Each call of `config` made a clock that began at
12:00:00, and each restart kept the run ID. A restart takes the lease over after 15
minutes of its own clock (three times `LeaseH`, 5 minutes) and writes the row with
`Started` cut to the second. Only one pair can collide: the first restart finds the row
that Start released and claims it at once; the second and the third each wait 15 minutes
of a clock that began again, so their `Started` can fall in the same second, and the two
rows are equal when no beat came between them: `git commit` exits 1 with nothing to
commit (`internal/git/git.go:168-174`). The window of about 28 seconds and the empty
commit as the reason of the exit are readings of the code, not measured. With one clock
the third restart's wait ends at least 15 minutes of the clock after the second's
`Started` (`lease.go:66-67, 81`): the rows cannot be equal.

## Review rounds

Round 1 (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at
`08506d5`, cycle 0, the cap of 1; its record is a comment on #293): `nothing material in
scope`, three notes. The `tests` job of CI passed on that head.

- **Applied** in `63c419e` (documents only; the linters exit 0): note 2 (the guardrails
  entry says "as it does when nothing changed", and the record of the runs says that the
  reason of the exit is read from the mechanism: git writes it to its standard output,
  which `internal/git` does not keep).
- **Carried by this file:** note 1 (the close-out is inside the room) and the two readings
  of note 2 ("The cause, and what is a reading").
- **Revealed, opened as an issue:** note 3 is [#296](https://github.com/pharzam/layup/issues/296):
  `FailedError` drops the standard output of a failed call of `internal/git`, so the lease
  commit of the flake read `exit status 1` alone; product code, off this path.

Round 2 (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at
`1a31d0d`, scoped to the merge of `origin/main` that brought #294 (`T-6rgp`) and changed
three files of this branch; it consumes no cycle, though the lint counts cycle 1):
`nothing material in scope`, one note, declined: the header of `docs/tasks/completed.md`
says "most recent first" and the log is appended in the order of landing, older than this
branch. Round 1's `Verdict` row was edited to `material` for the lint, as its comment
says. The `tests` job of CI passed on that head.

## Verdict

Delivered: the test world of `internal/run` has one clock, made once in `newWorld` and
read by `config`, so the clock of a restart does not go back, and
`TestTheClockOfAWorldDoesNotGoBack` fails on a clock made at each call and passes on the
clock of the world. The old test failed 1 of 200 runs here with the message of CI, the new
test 0 of 200: the counts alone cannot show the fix, and the red test and the sources do.
The whole integration suite passes idle and beside the loops, its longest test, 11.04 s,
being 23% of the 48 s of the hours cap. The review ended by decay at cycle 0 of cap 1. The
diff against `origin/main` (`a27b87e`, after the merge of #294) is 258 lines over 6 files with the
close-out, inside the Budget maximum of 260 lines over 7 files. Round 2 read the three
shared files at the merge; no round read the rest of the notes commit or the close-out.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-12. The reviewers' tokens are `modelUsage`,
with Claude Code's estimate of the money; the author's are not reported. The Operator chose
Claude Sonnet 5.5 for the session; Model tiers names `Claude Sonnet 5` for the execution
tier, so the plan and its answer are reasoning parts on it. The counts of this task and of
`T-6rgp` ran one after another on one host: the span of the fourth row holds both.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Sonnet 5.5 | `max` | not reported | 13:52 to 13:54 |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 812,860 (USD 5.11) | 9 min 45 s, 13:54 to 14:04 |
| The answer to the plan review | reasoning | Claude Sonnet 5.5 | `max` | not reported | 14:05 to 14:08 |
| T, D1, the counts and the mutant | execution | Claude Sonnet 5.5 | `max` | not reported | 14:09 to 15:36 |
| The record, the checks and the pull request | execution | Claude Sonnet 5.5 | `max` | not reported | 15:36 to 15:44 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 954,301 (USD 5.15) | 11 min 21 s, 15:45 to 15:56 |
| The notes, #296 and the close-out | execution | Claude Sonnet 5.5 | `max` | not reported | 15:58 to 16:02 |
| Round 2, scoped | reasoning | Claude Fable 5.1 | `xhigh` | 553,754 (USD 3.42) | 6 min 1 s, 16:05 to 16:11 |
| The merge of `origin/main`, round 2 and its close-out | — | Claude Sonnet 5.5 | `max` | not reported | 16:02 to 16:15 |
