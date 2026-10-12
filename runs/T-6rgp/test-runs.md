# T-6rgp — test runs

The test `TestEachEndGetsItsClass` (`internal/session/end_integration_test.go`) with a limit
for each case and a late crash case, as task `T-6rgp` (#292), cause 3 of the flaky `tests`
job of CI (#291). UTC, 2026-10-12. Host: Linux 6.12 (amd64), 2 CPUs, `go1.26.9`, `git`
2.47.3. The base is `1c167eb` (the merge of #286); the change is `d6851cb`, and the
documents of the task follow it. Each time is from `date -u`, or from the start and end
files of a count.

## Red (14:43Z)

On `1c167eb` with the one new line of the case "crash by a signal of its own, late" in the
table, and the limit of 300 ms still shared by each case:

```
$ go test -tags=integration -count=1 -v -run '^TestEachEndGetsItsClass$' ./internal/session
=== RUN   TestEachEndGetsItsClass
    end_integration_test.go:40: crash by a signal of its own, late: "wall" ({… Signal:interrupt StoppedBy:wall Sent:[interrupt]}, <nil>, <nil>), want "crash"
--- FAIL: TestEachEndGetsItsClass (0.62s)
FAIL
exit 1
```

One line, for the case that was added: the process ran 0.303 s (`Start` to `End`), was stopped
by the wall limit and sent an interrupt, before the shell's `sleep 0.5` ended and it killed
itself. The seven other cases pass. This is the failure that CI showed for the case
without the `sleep` (0.31 s and 0.368 s after the start), made certain: 0.5 s is more
than 0.3 s, and the timer of the wall limit starts just after the process does
(`process.go:153,202`).

## Green (14:43Z)

With the limit chosen by the case (300 ms for the case that expects `wall`, 5 s for each
other case): `--- PASS: TestEachEndGetsItsClass (0.82s)`, exit 0.

## The counts of 100

The load is one busy loop (`while :; do :; done` in `sh`) for each CPU of the host, so 2,
started two seconds before the count and ended after it. A count is one run of
`go test -tags=integration -count=100 -timeout 90m -v -run '^TestEachEndGetsItsClass$' ./internal/session`
in a directory at the head named, and the table counts the lines `--- PASS` and `--- FAIL`
of its output. "Idle" is no process of mine beside the test; the counts ran one after
another, so the load average before a count can still hold the end of the one before it
(the first count of the old test began as the loops of a count of `T-0p5p` ended).

| Test | Head | Load | Runs | Failed | A run took (s) | The count, UTC | Load average before |
| ---- | ---- | ---- | ---- | ------ | -------------- | -------------- | ------------------- |
| old: seven cases, 300 ms each | `1c167eb` | idle | 100 | 0 | 0.31 to 0.33 | 14:41:17 to 14:41:51 | 4.40 4.16 3.40 |
| old | `1c167eb` | 2 loops | 100 | 0 | 0.32 to 0.35 | 14:41:53 to 14:42:27 | 2.79 3.80 3.30 |
| new: eight cases, a limit for each | `d6851cb` | idle | 100 | 0 | 0.82 to 0.85 | 14:59:03 to 15:00:26 | 1.95 1.96 2.31 |
| new | `d6851cb` | 2 loops | 100 | 0 | 0.82 to 0.85 | 15:03:15 to 15:04:40 | 1.87 1.87 2.17 |

The two idle counts were run again from a quiet start, as the first of them began just after
the loops of a count of `T-0p5p` ended (the load average before it, 4.40, still held that
count): the old test, 100 runs, 0 failed (15:34:50 to 15:35:23, load average before 0.72
2.00 2.92); the new test, 100 runs, 0 failed (15:35:23 to 15:36:46, 0.49 1.83 2.83).

**What the counts show, and what they do not.** The old test did not fail once in 300 runs
on this host, idle or loaded: a shell here starts in a few milliseconds, even beside a
loop for each CPU, and the failure of CI needs a start of more than 300 ms. So the counts
do not show the flake, and 0 of 300 runs of the new test is not evidence of the fix. The
evidence is the late case: it fails on the shared limit on any host (the red run) and
passes on the limit of its own (the green run and the 300 runs of the new test). The new
test takes 0.5 s more for each run, the sleep of that case.

## The mutants (15:05Z)

On copies, never on the worktree: the tree of `d6851cb` taken out with `git archive` into a
directory of the scratchpad, and one file of each mutant made from the text of the head by
one replacement, its diff read before the run. After each run the file was written back from
the head and compared with it.

- **A1. The limit of 300 ms shared by each case again** (the block of the choice replaced by
  the line of `1c167eb`). It fails, with the one line of the red run: the late case,
  `"wall" … StoppedBy:wall Sent:[interrupt] … want "crash"`, exit 1. This repeats the red
  of step T on the final file, and shows one of the two rules: that a case with no limit
  of its own is stopped by the short one.
- **A2. The case `wall` with 5 s as well** (the `if` that gives it 300 ms removed). It
  passes, in 5.52 s where the test takes 0.82 s: the case is still classed `wall`, and only
  the run is longer. This is a **known limit** of the test, as the plan review asked: it
  cannot tell the short limit from the long one by the class. An upper bound on the time of
  the test would bring back the dependency on the clock of the host, which this task
  removes, so the test has none.

## The other tests of the package

`go test -tags=integration -count=1 -v ./internal/session` on `d6851cb`, once idle
(15:03Z) and once beside the 2 loops (15:06Z): 33 tests, 33 pass each time, in 7.6 s and
7.9 s. `TestProcessStopsTheWholeGroup` took 1.15 s and 1.16 s, and
`TestProcessEndsWhenAProgramLeftTheGroup` 1.10 s and 1.11 s: the two tests that also race
a shell's start against a limit of 300 ms (note 6 of the plan review) did not fail, so no
issue is opened for them.

## The checks (15:41Z to 15:42Z)

`go build ./...`, `go vet ./...`, `go vet -tags=integration ./internal/session` and `gofmt -l .` exit 0 with no output of
`gofmt` (15:06Z). After `git add` of the new files: `sh docs/adr/adr-lint.sh`,
`sh docs/prd/prd-lint.sh`, `sh docs/links/link-lint.sh` (2,353 links),
`sh docs/setup/setup-check.sh` (each of its checks OK) and `sh docs/tests/run-discipline-tests.sh`
(127 passed) exit 0, and so does `git diff --cached --check`. A first run of
`setup-check.sh`, with the task file staged, exited 1 on `kit-history`: the file had no line
in `docs/tasks/backlog.md` or `docs/tasks/completed.md`. The task file and its line of the
completed log are written together at the close-out, so the head under review has neither,
and the check passes.
