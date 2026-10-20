# T-05f4 — test runs

Task `T-05f4` (#439): the integration tests of `internal/run` run in parallel.
Top-level tests in `internal/run` with `-tags=integration` (`go test -tags=integration -list . ./internal/run/ | grep -c ^Test`): **705** (433 unit tests in `*_test.go` and 272 integration tests in `*_integration_test.go`, of which 257 call `t.Parallel()` and 15 stay serial).

## 1. Red evidence before Commit B

### S1a — `t.Parallel()` without moving `HOME` to `TestMain` (scratch copy of merged head)

- Command: `nice -n 19 env TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/`
- Start: `2026-10-20T13:41:28Z`
- End: `2026-10-20T13:41:41Z`
- Wall-clock seconds: `13`
- Load average before run (`cut -d' ' -f1-3 /proc/loadavg`): `0.31 1.96 4.72`
- Top-level tests in package: `705`
- Exit status: `1`

```
panic: testing: test using t.Setenv, t.Chdir, or cryptotest.SetGlobalRandom can not use t.Parallel [recovered, repanicked]

goroutine 22 [running]:
testing.tRunner.func1.2({0xb9fca0, 0xceaf30})
--
panic({0xb9fca0?, 0xceaf30?})
	/home/layup/.local/share/mise/installs/go/1.26.9/src/runtime/panic.go:860 +0x13a
testing.(*T).checkParallel(...)
	/home/layup/.local/share/mise/installs/go/1.26.9/src/testing/testing.go:1838
testing.(*T).Setenv(0x39e81ebe0d88?, {0xc56c78?, 0x40?}, {0x39e81ee8dd00?, 0x4ab56d?})
```

### S1a, reproduced against a named tree (the fix of finding F1 of round 1)

The tree of the run above was a scratch copy that was not kept, so its base and its patch cannot be shown. This run reproduces it on a named tree. It does not stand in for the run above.

- Tree: commit A, `9ac5aada65842f6a1107f3415015147c9b9e5c12`, with the HOME change undone: `newWorldWith` sets `t.Setenv("HOME", t.TempDir())` again (`internal/run/run_integration_test.go:232`), and `internal/run/main_integration_test.go` is removed. Every `t.Parallel()` of commit A stays. It was made in a detached worktree of `9ac5aad` with `git rm internal/run/main_integration_test.go` and one line inserted after `func newWorldWith(t *testing.T, files map[string]string) *world {`.
- Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green nice -n 19 go test -count=1 -tags=integration -timeout 30m ./internal/run/`
- Start: `2026-10-20T14:28:27Z`
- End: `2026-10-20T14:28:37Z`
- Exit status: `1`

The panic, with the caller that names the HOME assignment of the world:

```
panic: testing: test using t.Setenv, t.Chdir, or cryptotest.SetGlobalRandom can not use t.Parallel [recovered, repanicked]
testing.(*T).Setenv(0xf3da5012008?, {0xc56c78?, 0x500000005?}, {0xf3da4f42120?, 0x3900000030?})
	/home/layup/.local/share/mise/installs/go/1.26.9/src/testing/testing.go:1852 +0x65
github.com/pharzam/layup/internal/run.newWorldWith(0xf3da5012008, 0x0)
	/home/layup/.cache/layup-session/tmp-green/t05f4-s1a/internal/run/run_integration_test.go:232 +0x5b
github.com/pharzam/layup/internal/run.TestTheDraftPullRequestOfABoundHead(0xf3da5012008)
	/home/layup/.cache/layup-session/tmp-green/t05f4-s1a/internal/run/pullmerge_integration_test.go:125 +0x39
```

The frames between `newWorldWith` and the test (`sessionWorldWith` and the helpers above it) are cut. So the conflict is the per-test `HOME` of the world under `t.Parallel()`, which commit A moves into `TestMain`.

### S1b — Commit A (`9ac5aad`): `verifyRun` tests running in parallel (3 of 3 red)

#### Commit A, run 1

- Command: `nice -n 19 env TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/`
- Start: `2026-10-20T13:43:06Z`
- End: `2026-10-20T13:44:15Z`
- Wall-clock seconds: `69`
- Load average before run: `2.91 2.14 4.50`
- Top-level tests in package: `705`
- Exit status: `1`

```
--- FAIL: TestTheStepSetupRunsS01ToS14 (10.06s)
    --- FAIL: TestTheStepSetupRunsS01ToS14/a_second_restart_before_verify_rebuilds_the_same_head_and_starts_no_session (0.50s)
        setup_integration_test.go:243: the row of the step setup of a second restart: {Name:setup Result:done Detail:done before}; want setup done
--- FAIL: TestTheStepForgeCheckChecksTheForgeAndCommitsItsTableAtEachRestart (12.59s)
    forgecheck_integration_test.go:256: the step forge-check: {Name: Result: Detail:}; want {Name:forge-check Result:done Detail:the forge check passed}
```

#### Commit A, run 2

- Command: `nice -n 19 env TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/`
- Start: `2026-10-20T13:44:15Z`
- End: `2026-10-20T13:45:29Z`
- Wall-clock seconds: `74`
- Load average before run: `11.61 5.13 5.37`
- Top-level tests in package: `705`
- Exit status: `1`

```
--- FAIL: TestTheStepSetupRunsS01ToS14 (12.88s)
    --- FAIL: TestTheStepSetupRunsS01ToS14/a_second_restart_before_verify_rebuilds_the_same_head_and_starts_no_session (0.50s)
        setup_integration_test.go:243: the row of the step setup of a second restart: {Name:setup Result:done Detail:done before}; want setup done
```

#### Commit A, run 3

- Command: `nice -n 19 env TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/`
- Start: `2026-10-20T13:45:29Z`
- End: `2026-10-20T13:46:38Z`
- Wall-clock seconds: `69`
- Load average before run: `13.92 7.52 6.20`
- Top-level tests in package: `705`
- Exit status: `1`

```
--- FAIL: TestTheStepSetupRunsS01ToS14 (12.73s)
    --- FAIL: TestTheStepSetupRunsS01ToS14/a_second_restart_before_verify_rebuilds_the_same_head_and_starts_no_session (0.65s)
        setup_integration_test.go:243: the row of the step setup of a second restart: {Name:setup Result:done Detail:done before}; want setup done
--- FAIL: TestTheFullInputOfAHeadAfterIntakeAndTheScaffold (10.06s)
    specinput_integration_test.go:152: the Scaffold is not done: {Name: Result: Detail:}
```

## 2. Matched serial baseline on Commit B (`87bfbcb`)

- Command: `nice -n 19 env TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m -parallel 1 ./internal/run/`
- Start: `2026-10-20T13:48:13Z`
- End: `2026-10-20T13:54:49Z`
- Wall-clock seconds: `396` (`ok github.com/pharzam/layup/internal/run 394.614s`)
- Load average before run: `4.09 7.38 6.48`
- Top-level tests in package: `705`
- Exit status: `0`

## 3. Five parallel runs on Commit B (`87bfbcb`)

| Run | Command | Start (UTC) | End (UTC) | Wall-clock (s) | Package time | Load avg before | Top-level tests | Exit status |
| --- | ------- | ----------- | --------- | -------------- | ------------ | --------------- | --------------- | ----------- |
| 1 | `nice -n 19 env TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/` | `2026-10-20T13:54:49Z` | `2026-10-20T13:56:54Z` | `125` | `124.739s` | `0.96 2.64 4.54` | `705` | `0` |
| 2 | `nice -n 19 env TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/` | `2026-10-20T13:56:54Z` | `2026-10-20T13:59:03Z` | `129` | `128.254s` | `11.42 5.20 5.17` | `705` | `0` |
| 3 | `nice -n 19 env TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/` | `2026-10-20T13:59:03Z` | `2026-10-20T14:01:21Z` | `138` | `138.441s` | `14.25 7.40 5.90` | `705` | `0` |
| 4 | `nice -n 19 env TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/` | `2026-10-20T14:01:21Z` | `2026-10-20T14:03:25Z` | `124` | `122.923s` | `17.84 11.46 7.64` | `705` | `0` |
| 5 | `nice -n 19 env TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/` | `2026-10-20T14:03:25Z` | `2026-10-20T14:05:26Z` | `121` | `121.503s` | `13.73 11.05 7.89` | `705` | `0` |

## 4. Run with `-race` on Commit B (`87bfbcb`)

- Command: `nice -n 19 env TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -race -timeout 40m ./internal/run/`
- Start: `2026-10-20T14:08:19Z`
- End: `2026-10-20T14:10:46Z`
- Wall-clock seconds: `147` (`ok github.com/pharzam/layup/internal/run 133.289s`)
- Load average before run: `7.01 8.45 7.61`
- Top-level tests in package: `705`
- Exit status: `0`

## 5. Probe with `-cpu 4 -parallel 4` on Commit B (`87bfbcb`)

- Command: `nice -n 19 env TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m -cpu 4 -parallel 4 ./internal/run/`
- Start: `2026-10-20T14:05:46Z`
- End: `2026-10-20T14:08:19Z`
- Wall-clock seconds: `153` (`ok github.com/pharzam/layup/internal/run 153.544s`)
- Load average before run: `10.23 10.27 8.00`
- Top-level tests in package: `705`
- Exit status: `0`

## 6. Verdict

- Baseline (`-parallel 1`): `396 s` (`394.614s`)
- Five parallel runs (wall-clock): `125 s`, `129 s`, `138 s`, `124 s`, `121 s` — median `125 s` (`124.739s`)
- Ratio (`median(5 parallel) / baseline`): `125 / 396 = 0.32` (`124.739 / 394.614 = 0.32`)
- Rule (`pass <= 0.5`, `investigate 0.5-0.8`, `fail > 0.8`): **pass** (`0.32 <= 0.5`)

## 7. Other levels and text checks on Commit B (`87bfbcb`)

| Command | Start (UTC) | End (UTC) | Wall-clock (s) | Load avg before | Top-level tests in `internal/run` | Exit status |
| ------- | ----------- | --------- | -------------- | --------------- | --------------------------------- | ----------- |
| `go build ./...` | `2026-10-20T14:10:46Z` | `2026-10-20T14:10:46Z` | `0` | `14.99 9.89 8.18` | `705` | `0` |
| `go vet ./...` | `2026-10-20T14:10:46Z` | `2026-10-20T14:10:46Z` | `0` | `14.99 9.89 8.18` | `705` | `0` |
| `go vet -tags=integration ./...` | `2026-10-20T14:10:46Z` | `2026-10-20T14:10:47Z` | `1` | `14.99 9.89 8.18` | `705` | `0` |
| `go vet -tags=e2e ./...` | `2026-10-20T14:10:47Z` | `2026-10-20T14:10:47Z` | `0` | `14.99 9.89 8.18` | `705` | `0` |
| `nice -n 19 env TMPDIR=/home/layup/.cache/layup-session/tmp-green go test ./...` | `2026-10-20T14:10:47Z` | `2026-10-20T14:10:48Z` | `1` | `14.99 9.89 8.18` | `705` | `0` |
| `nice -n 19 env TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...` | `2026-10-20T14:10:48Z` | `2026-10-20T14:12:54Z` | `126` | `14.99 9.89 8.18` | `705` | `0` |
| `sh docs/adr/adr-lint.sh` | `2026-10-20T14:12:54Z` | `2026-10-20T14:12:55Z` | `1` | `12.97 10.51 8.60` | `705` | `0` |
| `sh docs/prd/prd-lint.sh` | `2026-10-20T14:12:55Z` | `2026-10-20T14:12:55Z` | `0` | `12.97 10.51 8.60` | `705` | `0` |
| `sh docs/links/link-lint.sh` | `2026-10-20T14:12:55Z` | `2026-10-20T14:13:07Z` | `12` | `12.97 10.51 8.60` | `705` | `0` |
| `sh docs/tests/run-discipline-tests.sh` | `2026-10-20T14:13:07Z` | `2026-10-20T14:13:10Z` | `3` | `10.97 10.16 8.51` | `705` | `0` |
| `sh docs/setup/setup-check.sh` | `2026-10-20T14:13:10Z` | `2026-10-20T14:13:31Z` | `21` | `10.18 10.01 8.47` | `705` | `0` |
| `git diff --check origin/main HEAD` | `2026-10-20T14:13:31Z` | `2026-10-20T14:13:31Z` | `0` | `7.42 9.39 8.30` | `705` | `0` |
