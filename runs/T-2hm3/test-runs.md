# T-2hm3 — test runs

The bare repositories that tests push into, set to start no maintenance, as task `T-2hm3`
(#297), cause 2 of the flaky `tests` job of CI (#291). UTC, 2026-10-12. Host: Linux 6.12
(amd64), 2 CPUs, `go1.26.9`. Three gits ran, and each run names its own: **2.47.3**
(`/usr/bin/git`); **2.55.0**, the git of the CI image (the package
`git_2.55.0-0ppa1~ubuntu24.04.2_amd64.deb` of the git-core PPA, unpacked with `dpkg-deb -x`
into the scratchpad, run through a wrapper that sets `GIT_EXEC_PATH` and stands first on
`PATH`); **2.56.0**, which the Operator installed at 16:57 (`~/.local/bin/git`, a link to
miniforge3), the git of the plain `PATH` since then. The base is `a8b3180` (the merge of
#295); the change is `ecca1ee`, and the documents follow it. Each time is from `date -u`, or
from the start and end files of a count.

## What a push starts in the target

One commit pushed into a bare repository, with `GIT_TRACE2_EVENT` set: the events of each
process of git, the receiving one included. The receiving git started `git maintenance run
--auto --quiet --detach` in each of the six cases (17:57Z), whether the pusher passed
`-c maintenance.auto=false` or not:

| git | the call of the pusher | started in the target |
| --- | ---------------------- | --------------------- |
| 2.47.3, 2.55.0, 2.56.0 | `git -c maintenance.auto=false push` | `git maintenance run --auto --quiet --detach` |
| 2.47.3, 2.55.0, 2.56.0 | `git push` | the same |

The flag of the pusher does not arrive: a local transport clears `GIT_CONFIG_PARAMETERS`
before it starts `git receive-pack` (`environment.c`, `local_repo_env`, in the source of git
2.55.0, which the tree does not hold). The process takes `objects/maintenance.lock`,
detaches, and runs its tasks after the push has returned.

## The made setup

[`repro.go`](repro.go) (run with `go run`; it is no package of the module): a bare
`acme/target.git` with two loose objects in `objects/17`, a push into it, and `os.RemoveAll`
of the directory at once, 100 times. The repository's own settings are `receive.autogc` and
`maintenance.auto`, both `false` with `-fix`.

| git | the target's settings | removals that failed | The runs, UTC |
| --- | --------------------- | -------------------- | ------------- |
| 2.55.0 | none | 16 of 100 | 17:19:04 to 17:23:15, 4 min 11 s |
| 2.55.0 | the two, false | 0 of 100 | 17:23:15 to 17:27:27, 4 min 12 s |
| 2.56.0 | none | 0 of 100 | 17:27:27 to 17:35:28, 8 min 1 s |
| 2.56.0 | the two, false | 0 of 100 | 17:35:28 to 17:43:22, 7 min 54 s |
| 2.47.3 | none | 0 of 100 | 17:43:22 to 17:46:27, 3 min 5 s |

The failures of git 2.55.0 have the message of CI: `unlinkat …/acme/target.git/objects:
directory not empty`. Which task ran, from a trace of 6 runs of the same setup with every
git process traced (17:57Z):

```
== git 2.55.0: runs/T-2hm3/repro.go -n 6, every git process traced (17:57:24Z)
fix=false: 2 of 6 removals failed
exit status 1
  maintenance --auto runs: 6 | regions entered: {('detach', 'geometric-repack'): 6}
  repack processes: {'repack -d -l --cruft --cruft-expiration=2.weeks.ago --quiet --write-midx': 6}
== git 2.56.0: runs/T-2hm3/repro.go -n 6, every git process traced (17:57:43Z)
git version 2.56.0
fix=false: 0 of 6 removals failed
  maintenance --auto runs: 6 | regions entered: {('detach',): 6}
  repack processes: {}
```

So on git 2.55.0 the task `geometric-repack` runs in the background of the push, and writes
a cruft pack and a multi-pack index while the directory is removed. Git 2.56.0 enters
`detach` and runs nothing (the condition moved into the object database layer, `odb.c`, and
two loose objects in `objects/17` no longer cross it). The function of git 2.54.0 is that of
2.55.0 (read, not run), so the setup applies from 2.54.

## Red (T) and green (D1)

The new test, `TestAPushIntoTheTargetOfAWorldStartsNoMaintenance`
(`go test -tags=integration -count=1 -v ./internal/run`), with each git: red on the tree with
no setting in `newWorld`, green with the two settings. Each run holds the positive control
(the event of `receive-pack`), so the trace is that of the receiving git; the red runs fail
twice, for the argv and for the `cmd_name` of the maintenance.

| git | red, UTC | green, UTC |
| --- | -------- | ---------- |
| 2.47.3 | FAIL in 0.17 s, 17:16:08Z | PASS in 0.31 s, 17:16:39Z |
| 2.55.0 | FAIL in 0.14 s, 17:16:12Z | PASS in 0.29 s, 17:16:43Z |
| 2.56.0 | FAIL in 0.36 s, 17:16:13Z | PASS in 0.28 s, 17:16:44Z |

```
git version 2.55.0
    run_integration_test.go:348: a push into the target started the maintenance of git: "maintenance","run" is in the trace
    run_integration_test.go:348: a push into the target started the maintenance of git: "name":"maintenance" is in the trace
--- FAIL: TestAPushIntoTheTargetOfAWorldStartsNoMaintenance (0.14s)
FAIL	github.com/pharzam/layup/internal/run	0.149s
exit 1
```

## The counts of 100

The three tests of the failures (`TestATaskSessionProbesItsVersionFirst`,
`TestARefusedPushBindsNothing`, `TestAnArtifactThatTheHeadLacksIsRefused`) in one run of
`go test -tags=integration -count=100 -timeout 90m -v` of `internal/run`, with git 2.55.0
first on `PATH`: 300 test runs in each count. The load is one busy loop (`while :; do :; done`
in `sh`) for each of the 2 CPUs, started two seconds before the count and ended after it;
"idle" is no process of mine beside the test. The old rows read the tree at `a8b3180`.

| Fixture | Head | Load | Runs | Failed | A test took (s) | The count, UTC | Load average before |
| ------- | ---- | ---- | ---- | ------ | --------------- | -------------- | ------------------- |
| old: no setting | a8b3180 | idle | 300 | 0 | 1.25 to 3.06 | 16:48:50 to 16:59:09, 10 min 19 s | 0.50 0.79 1.02 |
| old: no setting | a8b3180 | 2 loops | 300 | 0 | 2.12 to 5.67 | 16:59:12 to 17:15:26, 16 min 14 s | 1.97 1.89 1.49 |
| new: the two settings | ecca1ee | idle | 300 | 0 | 1.32 to 3.07 | 17:46:27 to 17:56:28, 10 min 1 s | 2.17 2.00 2.16 |
| new: the two settings | ecca1ee | 2 loops | 300 | 0 | 2.09 to 4.95 | 17:56:30 to 18:12:25, 15 min 55 s | 1.44 1.64 1.91 |

**What the counts show, and what they do not.** The three tests failed 3 times in the 1,820
runs of the 26 tests that use `sessionWorld` in 70 runs of CI, about 1 in 600. The old
fixture failed 0 of 600 here, and the new one 0 of 600: with a rate of 1 in 600 the chance of 600 runs with no
failure is 37%, so the counts show only that the change breaks nothing. The evidence is the
test of the trace (red on each git, green after) and the made setup (16 of 100 failures on
git 2.55.0, none with the settings).

## The suites

`go test -count=1 -v` of each package, whole, with each git first on `PATH`, idle (the
integration level for `internal/run` and `internal/git`, the end-to-end level for
`cmd/layup`, as CI runs them):

| git | package | tests | failed | skipped | The run, UTC |
| --- | ------- | ----- | ------ | ------- | ------------ |
| 2.47.3 | `internal/run` (integration) | 105 | 0 | 0 | 18:12:37 to 18:15:18, 2 min 41 s |
| 2.47.3 | `internal/git` (integration) | 34 | 0 | 0 | 18:15:18 to 18:15:21, 0 min 3 s |
| 2.47.3 | `cmd/layup` (e2e) | 34 | 0 | 0 | 18:15:21 to 18:16:34, 1 min 13 s |
| 2.55.0 | `internal/run` (integration) | 105 | 0 | 0 | 18:16:34 to 18:19:37, 3 min 3 s |
| 2.55.0 | `internal/git` (integration) | 34 | 0 | 0 | 18:19:37 to 18:19:40, 0 min 3 s |
| 2.55.0 | `cmd/layup` (e2e) | 34 | 0 | 0 | 18:19:41 to 18:20:57, 1 min 16 s |
| 2.56.0 | `internal/run` (integration) | 105 | 0 | 0 | 18:20:57 to 18:25:03, 4 min 6 s |
| 2.56.0 | `internal/git` (integration) | 34 | 0 | 0 | 18:25:03 to 18:25:07, 0 min 4 s |
| 2.56.0 | `cmd/layup` (e2e) | 34 | 0 | 0 | 18:25:07 to 18:26:30, 1 min 23 s |

## The mutants

On a copy of the tree at `ecca1ee`, never on the worktree: the file of `internal/run` taken
from the head, the lines of `newWorld` replaced, the test run with git 2.55.0 and 2.56.0. The
file was written back from the head and compared with it afterwards.

```
174,175d173
< 	gitOut(t, w.bare, "config", "receive.autogc", "false")
< 	gitOut(t, w.bare, "config", "maintenance.auto", "false")
17:59:02Z
=== M1-none, git 2.55.0
    run_integration_test.go:351: a push into the target started the maintenance of git: "maintenance","run" is in the trace
    run_integration_test.go:351: a push into the target started the maintenance of git: "name":"maintenance" is in the trace
--- FAIL: TestAPushIntoTheTargetOfAWorldStartsNoMaintenance (0.54s)
FAIL
FAIL	github.com/pharzam/layup/internal/run	0.556s
FAIL
17:59:12Z
=== M1-none, git 2.56.0
    run_integration_test.go:351: a push into the target started the maintenance of git: "maintenance","run" is in the trace
    run_integration_test.go:351: a push into the target started the maintenance of git: "name":"maintenance" is in the trace
--- FAIL: TestAPushIntoTheTargetOfAWorldStartsNoMaintenance (0.94s)
FAIL
FAIL	github.com/pharzam/layup/internal/run	0.949s
FAIL
175d174
< 	gitOut(t, w.bare, "config", "maintenance.auto", "false")
17:59:15Z
=== M2-autogc-only, git 2.55.0
--- PASS: TestAPushIntoTheTargetOfAWorldStartsNoMaintenance (1.06s)
ok  	github.com/pharzam/layup/internal/run	1.078s
17:59:23Z
=== M2-autogc-only, git 2.56.0
--- PASS: TestAPushIntoTheTargetOfAWorldStartsNoMaintenance (0.37s)
ok  	github.com/pharzam/layup/internal/run	0.381s
174d173
< 	gitOut(t, w.bare, "config", "receive.autogc", "false")
17:59:26Z
=== M3-maintenance-only, git 2.55.0
--- PASS: TestAPushIntoTheTargetOfAWorldStartsNoMaintenance (0.61s)
ok  	github.com/pharzam/layup/internal/run	0.625s
17:59:33Z
=== M3-maintenance-only, git 2.56.0
--- PASS: TestAPushIntoTheTargetOfAWorldStartsNoMaintenance (0.52s)
ok  	github.com/pharzam/layup/internal/run	0.533s
```

Each key alone stops the receiving git (M2, M3), so the test cannot tell the two apart; the
entry of the guardrails says to set both. With neither (M1) the test fails as in the red run.

## The checks

`go build ./...`, `go vet ./...`, `go vet -tags=integration ./internal/run ./internal/git`,
`go vet -tags=e2e ./cmd/layup` and `gofmt -l .` exit 0, with no output of `gofmt`; `go list
./...` names no package under `runs/`, as `repro.go` is build-ignored. After `git add` of the
new files: `sh docs/adr/adr-lint.sh`, `sh docs/prd/prd-lint.sh`, `sh docs/links/link-lint.sh`
(2,358 links), `sh docs/setup/setup-check.sh` (each check OK) and
`sh docs/tests/run-discipline-tests.sh` (127 passed) exit 0, and so does `git diff --cached
--check` (18:26Z to 18:27Z). The task file and its line of the completed log come with the
close-out, so `kit-history` finds no task file without its line.
