# T-0p5p — one clock for each test world of `internal/run`

Issue: [#293](https://github.com/pharzam/layup/issues/293), cause 1 of the flaky `tests`
job of CI ([#291](https://github.com/pharzam/layup/issues/291)). Serves `F-0003#47`
(Verification on Every Change). `T-d6q5` (#92) met a failure of a test helper in CI alone
too (the maintenance that `git` runs in the background, the entry "A background
maintenance of git" of `docs/guardrails.md`); that was not a clock. Base `1c167eb` (the
merge of #286). Author: Claude Sonnet 5.5 on Claude Code. Evidence:
[`runs/T-0p5p/`](../../runs/T-0p5p/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #293. The plan review
(Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at `1c167eb`, 9 min
45 s) gave `approve-with-conditions`: Budget maximum 260 lines added plus removed over 7
files, close-out inside; Cycle cap 1; no panel. Its three conditions are applied: "with the
CPUs busy" is one busy loop for each CPU of the host (2), the CPU count beside the counts;
the risk of the hours cap is bounded by numbers (8.0 hours at 600 times is 48 s of real
time, and the evidence file gives the time of each test of the suite, idle and beside the
loops, against it); and the entry of `docs/guardrails.md` section 2 and the rows of
`docs/tests/traceability.md` are written. The notes are applied or carried: one pair of
restarts can collide, not two (the issue body is corrected); the window of about 28
seconds and the empty commit are readings; the mutant runs on the final file.

## What was done

1. **T:** `TestTheClockOfAWorldDoesNotGoBack` reads the clock of a first config after 20 ms
   (12 s of the clock) and asserts that the clock of a second config is not before it. Red
   on `1c167eb`: the second reads 12:00:00.02, the first 12:00:12.27.
2. **D1:** the world has one clock (`world.clock`), made in `newWorld`, and `config` reads
   it. The test is green, and so are 100 runs of `TestAFailedSpanCheckRunsNumberingAgain`
   and the whole integration suite of `internal/run`, idle and beside a busy loop for each
   CPU; the longest test of the suite takes 11.04 s beside the loops, against the 48 s.
3. **M:** one mutant on a copy of the final tree: `config` makes its clock again, and the
   new test fails with the line of the red run.
4. **Step 7:** the entry "A fake clock made anew at each call of a world's helper" in the
   testing pitfalls of `docs/guardrails.md`; and in `docs/tests/traceability.md` a row for
   the new test, with `—` in Covers and Fact, and `T-0p5p` in the row of the flaky test.

**The cause, and what is a reading.** Each call of `config` made a clock that began at
12:00:00, and each restart kept the run ID. A restart takes the lease over after 15
minutes of its own clock (three times `LeaseH`, 5 minutes), and writes the row with
`Started` cut to the second. By the plan review only one pair can collide: the first
restart finds the row that Start released and claims it at once; the second waits for the
15 minutes; the third waits for 15 minutes of a clock that began again, so its `Started` can
fall in the same second as the second's. With no beat between them the two rows are equal,
and `git commit` exits 1 with nothing to commit (`internal/git/git.go:168-174` passes no
`--allow-empty`). The window of about 28 seconds that the issue gives, and the empty commit
as the reason for the exit, are readings of the code, not measured. With one clock for the
world the third restart reads the second's row after it was written, so its wait ends at
least 15 minutes of the clock after that row's `Started` (`lease.go:66-67, 81`), and a
released row differs by `State`: the two rows cannot be equal.

## Review rounds

Round 1 (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at
`08506d5`, cycle 0, the cap of 1; its record is a comment on #293): `nothing material in
scope`, three notes. The `tests` job of CI passed on that head.

- **Applied** in `63c419e`, a commit that no round read (documents only; the linters exit
  0): note 2 (the guardrails entry says "as it does when nothing changed", and the record
  of the runs says that the reason of the exit is read from the mechanism: git writes it to
  its standard output, which `internal/git` does not keep).
- **Carried by this file:** note 1 (the close-out is inside the room that the review
  counted) and the two readings of note 2 ("The cause, and what is a reading").
- **Revealed, opened as an issue:** note 3 is [#296](https://github.com/pharzam/layup/issues/296):
  `FailedError` drops the standard output of a failed call of `internal/git`, so the lease
  commit of the flake read `exit status 1` alone. A change of product code, off the path of
  this task.

## Verdict

Delivered: the test world of `internal/run` has one clock, made once in `newWorld` and
read by `config`, so the clock of a restart does not go back, and
`TestTheClockOfAWorldDoesNotGoBack` fails on a clock made at each call and passes on the
clock of the world. The old test failed 1 of 200 runs here with the message of CI, the new
test 0 of 200: the counts alone cannot show the fix, and the red test and the sources do.
The whole integration suite passes idle and beside the loops, its longest test, 11.04 s,
being 23% of the 48 s of the hours cap. The review ended by decay at cycle 0 of cap 1. The
diff against `origin/main` (`1c167eb`, the merge of #286) is 256 lines over 6 files with the close-out,
inside the Budget maximum of 260 lines over 7 files. No round read the notes commit or this
close-out.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-12. The reviewers' tokens are `modelUsage`,
with Claude Code's estimate of the money; the author's are not reported. The Operator chose
Claude Sonnet 5.5 for the session; Model tiers names `Claude Sonnet 5` for the execution
tier, so the plan and its answer are reasoning parts on it. The counts of this task and of
`T-6rgp` ran one after another on one host, so the span of the fourth row holds both.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Sonnet 5.5 | `max` | not reported | 13:52 to 13:54 |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 812,860 (USD 5.11) | 9 min 45 s, 13:54 to 14:04 |
| The answer to the plan review | reasoning | Claude Sonnet 5.5 | `max` | not reported | 14:05 to 14:08 |
| T, D1, the counts and the mutant | execution | Claude Sonnet 5.5 | `max` | not reported | 14:09 to 15:36 |
| The record, the checks and the pull request | execution | Claude Sonnet 5.5 | `max` | not reported | 15:36 to 15:44 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 954,301 (USD 5.15) | 11 min 21 s, 15:45 to 15:56 |
| The notes and #296 | execution | Claude Sonnet 5.5 | `max` | not reported | 15:58 to 16:00 |
| The close-out | — | Claude Sonnet 5.5 | `max` | not reported | 16:00 to 16:02 |
