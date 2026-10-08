# T-ax3r — the test runs

Evidence for row 25b (#142): each test red first, for the right reason, then green.

## `forge.Repository` with the branch names, `git.Clone` with `Auth` (2026-10-08T10:40Z)

Red: the adapter's test asks `Branches` (one page, then 101 branches over two pages, then none), and `TestEachCallRunsItsVerb` gets the row "clone with a token" (the three `GIT_CONFIG_*` values): neither compiled (`unknown field Branches`, `too many arguments in call to Clone`). Then `Repository` reads every page of the branches (at the register's `api` only), `Clone` takes `Auth` through `doAuth`, and `internal/setup` passes `git.Auth{}`: `go test` and `go test -tags=integration` of `internal/git`, `internal/forge/...` and `internal/setup` pass.

## `internal/run`: red 1 (2026-10-08T10:36Z)

`steps_test.go` (a stand-in forge: the `forge` step of Start and of the restart with each input state of rows 6, 7 and 11; the plan check on each plan and visibility; `run-steps`) and `run_integration_test.go` (a fake GitHub on `httptest` whose branches are those of a bare target, the real adapter, a bare baseline, a stand-in Operator who pushes the root commit, one clock 600 times faster than real time for the run and the adapter) before the code: `go test ./internal/run/` did not compile (`undefined: Config`, `undefined: state`).

## Green, with three defects found on the way (2026-10-08T10:55Z)

The code (`steps.go`, `store.go`, `texts.go`): the unit tests passed at once. The integration tests (`-race`) found, in order:
- step 6 refused `harness.devin.cap` of `—`: `internal/tsv` holds the empty value as `""` and writes it as `—`; the run put `—` in memory. The same defect was in `Copy` and `Decision` of row 25a: a row read back from `copies.tsv` holds `""`, so a real decision would not count on a restart. Red first: `rules_test.go` was changed to the empty value of `tsv.Read`, and four cases of `Decision` and the first copy failed (`Decision = false, want true`); then `""` in `rules.go` and `steps.go`, and both pass.
- a restart after a stop at `opening` wrote `opening` again, so `git commit` had nothing to commit (exit 1); the step now skips the announcement that a restart finds made.
- `beat` made a context that no step used, so a lost beat stopped nothing (found by reading before the run; the test of condition 2 below proves the fix).

Then the demo, the restart with another version, the restart after a stop at `opening` (a takeover from the stopped run), and rows 6 and 8 of the input states pass with `-race`; `TestTheTextBlocksOfRunMd` fails until `run.md` has the two issue blocks (the documents step).

## The tests of conditions 1 and 2, and two mutations (2026-10-08T11:02Z)

`TestATakeoverAfterACommitThatTheCloneDidNotSee` and `TestALostBeatEndsTheWatch` passed at once on the code, so each was shown to fail on a mutation of a backup copy, put back after (`cmp` equal):
- `ReadLease` keeps no commit of its read → `the restart: the step lease is fail: the records push was refused; the lease names the run 0123456789abcdef`;
- a lost beat does not end the run's context → `the watch read 316 times after the beat was lost; want it ended at the loss` (the check of the number of reads was added first, as the detail of a failed write names the other run too).

## One flaky test, fixed (2026-10-08T11:15Z)

In six runs of `go test -race -count=1 -tags=integration ./internal/run/`, one failed: `TestALostBeatEndsTheWatch`, `git push -q origin layup-records: exit status 1`. The cause is the test's helper that plays the other run: a beat of the run under test came between the helper's clone and its push, so git refused the helper's push. A real other run reads again and tries again, so the helper now tries up to five times. Then eight runs in a row passed.

## The documents (2026-10-08T11:12Z)

Red: `runs/T-ax3r/docs.sh` on a work tree of the code commit `a949c65` gave 21 `FAIL` lines of 21. After the edits of `run.md`, `forge.md`, `packages.md`, the traceability and the PRD (two rules of the check were fixed on the way: a phrase that the text wraps, and the PRD rows that a guard of the edit skipped), 21 `ok`, exit 0.
