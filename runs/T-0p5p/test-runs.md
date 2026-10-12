# T-0p5p — test runs

The test world of `internal/run` (`internal/run/run_integration_test.go`) with one clock for
each world, and the test `TestTheClockOfAWorldDoesNotGoBack`, as task `T-0p5p` (#293),
cause 1 of the flaky `tests` job of CI (#291). UTC, 2026-10-12. Host: Linux 6.12 (amd64),
2 CPUs, `go1.26.9`, `git` 2.47.3. The base is `1c167eb` (the merge of #286); the change is
`1df6801`, and the documents of the task follow it. Each time is from `date -u`, or from
the start and end files of a count.

## Red (14:43Z)

On `1c167eb` with the new test added and the clock still made at each call of `config`:

```
$ go test -tags=integration -count=1 -v -run '^TestTheClockOfAWorldDoesNotGoBack$' ./internal/run
=== RUN   TestTheClockOfAWorldDoesNotGoBack
    run_integration_test.go:314: the clock of a second config reads 2026-10-08T12:00:00.016731Z, before the 2026-10-08T12:00:12.267489Z that the first read
--- FAIL: TestTheClockOfAWorldDoesNotGoBack (0.18s)
FAIL
exit 1
```

The first config's clock read 12 s after its start, after the 20 ms of the test (600 times
20 ms); a second config began again at 12:00:00. This is the cause, and the red is certain:
a clock that is made at the call begins at 12:00:00 whenever it is made, and the first one
has run on for 20 ms.

## Green (14:43Z)

With the clock made once in `newWorld` and read by `config`:
`--- PASS: TestTheClockOfAWorldDoesNotGoBack (0.18s)` and
`--- PASS: TestAFailedSpanCheckRunsNumberingAgain (7.92s)`, exit 0.

## The counts of 100

The load is one busy loop (`while :; do :; done` in `sh`) for each CPU of the host, so 2,
started two seconds before the count and ended after it. A count is one run of
`go test -tags=integration -count=100 -timeout 90m -v -run '^TestAFailedSpanCheckRunsNumberingAgain$' ./internal/run`
in a directory at the head named, and the table counts the lines `--- PASS` and `--- FAIL`
of its output. "Idle" is no process of mine beside the test; the counts ran one after
another, so the load average before a count can still hold the end of the one before it.
(The default timeout of 10 minutes is not enough for 100 runs of 8 s: a first count, with
it, stopped at 600 s, and was thrown away; it had shown one failure of the message below.)

| Test | Head | Load | Runs | Failed | A run took (s) | The count, UTC | Load average before |
| ---- | ---- | ---- | ---- | ------ | -------------- | -------------- | ------------------- |
| old: a clock made at each config | `1c167eb` | idle | 100 | 1 (the 10th) | 7.75 to 8.67 | 14:09:00 to 14:22:31 | 0.16 0.63 0.85 |
| old | `1c167eb` | 2 loops | 100 | 0 | 10.41 to 12.92 | 14:22:33 to 14:41:17 | 1.92 1.85 1.50 |
| new: one clock for the world | `1df6801` | idle | 100 | 0 | 7.47 to 8.49 | 14:45:44 to 14:59:03 | 1.82 2.83 3.01 |
| new | `1df6801` | 2 loops | 100 | 0 | 10.29 to 11.87 | 15:06:25 to 15:24:45 | 1.55 1.81 2.10 |

The failed run of the old test (the 10th of the idle count, 6.15 s) shows the message of CI:

```
=== RUN   TestAFailedSpanCheckRunsNumberingAgain
    run_integration_test.go:805: the restart after the fail: the step lease is fail: git commit -m lease: the run 0123456789abcdef takes the lease over from the run 0123456789abcdef: exit status 1
    run_integration_test.go:805: the restart after the fail: the steps [forge clone version lease], want [forge clone version lease probe rule-gaps meaning … ]
    …
--- FAIL: TestAFailedSpanCheckRunsNumberingAgain (6.15s)
```

**What the counts show, and what they do not.** The old test failed once in 200 runs on
this host, with that message: a rate of 0.5%, below the 5 in 70 of CI (7%) and below the one
in 28 (3.6%) that the plan review reads from the mechanism; a count below one in 28 does not
refute that reading. The new test did not fail in 200 runs. At a rate of 0.5% the chance of
200 runs with no failure is 37% (0.995 to the power 200), at 1 in 28 it is 0.07%, so the
counts alone cannot show the fix. The fix is shown by the red test, which names the cause
(a clock that goes back), and by the sources: with one clock for the world, the third
restart reads the row of the second after it was written, so its wait ends at least 15
minutes of the clock after that row's `Started` (`lease.go:66-67, 81`), and the two rows
cannot be equal.

## The suite, and the hours cap

The hours cap of Intake compares the clock of the world with the commit time of the root
records commit. The cap of the world is `intake.cap` `50.0,8.0` (`run_integration_test.go:243`),
and 8.0 hours of a clock at 600 times is 48 s of real time, so a world breaks the cap only
when more than 48 s of real time pass between its Start and its last session of Intake. The
time of a whole test is an upper bound on that. `go test -tags=integration -count=1 -v
./internal/run` on `1df6801`, the 104 tests of the package, all pass:

| Load | The run, UTC | Sum of the test times | The longest tests (s) |
| ---- | ------------ | --------------------- | --------------------- |
| idle | 15:00:26 to 15:03:05 | 157.4 s | `TestAFailedSpanCheckRunsNumberingAgain` 7.68; `TestAnOutputFileThatIsRefusedIsNoRecord` 6.95; `TestTheAnswersToTheFollowUpReplaceAValueOfTheBudget` 6.05 |
| 2 loops | 15:24:47 to 15:29:04 | 254.7 s | `TestAnOutputFileThatIsRefusedIsNoRecord` 11.04; `TestAFailedSpanCheckRunsNumberingAgain` 10.97; `TestTheAnswersToTheFollowUpReplaceAValueOfTheBudget` 9.55 |

The longest test takes 11.04 s beside the loops, 23% of the 48 s: a world would need more
than four times that to break the cap. `TestTheAnswersToTheFollowUpReplaceAValueOfTheBudget`, the
world with four restarts after Start, takes 6.05 s and 9.55 s, and the new test 0.16 s and
0.26 s.

## The mutant (15:05Z)

On a copy, never on the worktree: the tree of `1df6801` taken out with `git archive` into a
directory of the scratchpad, the line of `config` read back as it was before the change (the
clock made at each call), its diff read before the run. After the run the file was written
back from the head and compared with it. The mutant is made from the final file, and with
that one line it is the file of `1c167eb` in effect, so its failure repeats the red run of
step T on the final file, as the plan review expected: nothing else in the file keeps the
clock of a second config from going back.

```
$ go test -tags=integration -count=1 -v -run '^TestTheClockOfAWorldDoesNotGoBack$' ./internal/run
    run_integration_test.go:316: the clock of a second config reads 2026-10-08T12:00:00.0118968Z, before the 2026-10-08T12:00:12.2048652Z that the first read
--- FAIL: TestTheClockOfAWorldDoesNotGoBack (0.11s)
exit 1
```

## The checks (15:40Z to 15:41Z)

`go build ./...`, `go vet ./...`, `go vet -tags=integration ./internal/run` and `gofmt -l .` exit 0 with no output of
`gofmt` (15:06Z). After `git add` of the new files: `sh docs/adr/adr-lint.sh`,
`sh docs/prd/prd-lint.sh`, `sh docs/links/link-lint.sh` (2,353 links),
`sh docs/setup/setup-check.sh` (each of its checks OK) and `sh docs/tests/run-discipline-tests.sh`
(127 passed) exit 0, and so does `git diff --cached --check`. A first run of
`setup-check.sh`, with the task file staged, exited 1 on `kit-history`: the file had no line
in `docs/tasks/backlog.md` or `docs/tasks/completed.md`. The task file and its line of the
completed log are written together at the close-out, so the head under review has neither,
and the check passes.
