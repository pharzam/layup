# T-3xvc: the test runs

Task `T-3xvc` (#368), row 67 of M2e: the e2e of the specification task. Times are `date -u` readings (or the modification times of the output files, in UTC) taken after each step; an exit status is the direct one of the command, never of a pipe. Go is 1.26.10 (`GOROOT` and `PATH` of mise, in every command). The raw outputs are in `/home/layup/.cache/layup-session/t3/` (`base.txt`, `red1.txt`, `red-copies.txt`, `red-base.txt`, `mutants-out.txt`).

## S0. The reconciliation with 66k

`origin/main` fetched over HTTPS: `30e2c0d48b5a0f8c9e6132e9232ce81bf3779a49` (the merge of #465, T-9nbf). The branch stood at the claim `c0f691ef` with no commit of its own, so `git merge origin/main` was a fast-forward and made no merge commit.

| Cited by the plan (6312881896) | At `30e2c0d` | Divergence |
| --- | --- | --- |
| `afterScaffold` `steps.go:152` | `steps.go:152` gives `shape`, `bet`, `loop` wrapped by `atBoundary` | none |
| `restartNames` `steps.go:211-220` | `steps.go:211-220`; with the Scaffold done: forge clone version lease forge-check + Intake + setup verify scaffold-push + shape bet loop | none |
| `runLoop` `loop.go:147`, `done before` l.186, `the loop merged` l.188, `actOn` l.330 | same lines | none |
| the e2e fake's loop permissions `run_e2e_test.go:114` | `run_e2e_test.go:115` (contents, issues, metadata, pull_requests, statuses write; checks read) | one line |
| `Restart` chooses its tail before the forge check runs (review condition 1, `steps.go:135`) | `steps.go:135-139` | none: the restart after `commands.sh` ends at `forge-check` |

No landed form makes a step impossible. Baseline: `go test -count=1 -tags=e2e ./cmd/layup/` green, `ok github.com/pharzam/layup/cmd/layup 64.407s` (09:51:35 start, 09:53:33 end, exit 0).

## S4. The first run of `TestTheSpecificationTaskEndToEnd`, before any product fix (condition 5)

`go test -count=1 -tags=e2e -run 'TestTheSpecificationTaskEndToEnd' ./cmd/layup/`, 09:58:28, exit 1 (`red1.txt`). Stage 5 failed at the row `shape`: `waiting for the Operator: shape: no survivor`. This is a **fixture** failure and proves nothing of the demo: the world's constraints were the invariant `Inv-1` and also the accepted decision record `ADR-0001` of the baseline (`docs/adr/0001-pin-the-baseline.md`), for which the option `pg` had no row. The fixture was fixed (a row `db pg ADR-0001 pass docs/adr/0001-pin-the-baseline.md`).

Three more fixture failures followed, each fixed in the test only: the success-criterion row that the setup answers put into the draft (`R-002`, covering `Q-005`) had no row in the generated PRD (the harness now reads the draft that the prompt lists); the fake's comment IDs of two issues collided (`the last copy of the comment 1001 is not of these bytes`); and the fake's `merge_method` field had no JSON tag.

The first run that reached a **product** assertion, `stage 5`, exit 1: `bet fail T-8hnw: Shape has no verified head for it`. Cause, read from the records branch (`git log --stat -- copies.tsv`): the commit `records: the copy of the comment 904` (the Operator's `set brief.lines`, read by the parameter boundary of `bet` for the first time, since `parameters.tsv` is skipped before the first task) wrote `copies.tsv` from the rows read at the start of the run: 1 insertion, 6 deletions, the copies `1002` to `1014` of the Shape tasks' comments lost, so the verifier's review record had no kept copy. A defect on the demo path, fixed in `internal/run/steps.go` (`copy` reads the file again before it adds a row).

The next run, after that fix: stage 6 exit 0 but `loop done waiting: <T>: not ready at 08d8cddc...: layup/rules: failure`. The status description at the architecture task's close-out head: `rule path: docs/facts/F-0003-numbered-facts.md`. Cause: the architecture task starts at the specification's verified head; after the specification task merged (rendered commit and close-out), `BringUpToDate` merged the default branch into it, and the four checks judged the head against the base of the attempt, so the facts record that the merge brought in was a rule-path change of the task. A defect on the demo path, fixed in `internal/run/finish.go` (`FinishTask` judges a head that carries a merge of the default branch against the second parent of the merge).

Two more assertion fixes in the test (first-parent merges only, since the architecture branch holds a merge of `main` too; `git diff --name-status` sorts the added facts record first).

## The red evidence of the two product fixes

Each fix reverted alone (`git apply -R` of its patch, `fix-copies.patch` sha256 `84d138b0...`, `fix-base.patch` sha256 `55c3e85b...`), the test run, the patch applied again:

| Fix reverted | Time | Result |
| --- | --- | --- |
| `steps.go` (copies) | 10:07:20 | exit 1, `bet fail T-9ty7: Shape has no verified head for it` (`red-copies.txt`) |
| `finish.go` (base) | 10:07:41 | exit 1, `loop done waiting: T-wjz9: not ready at a494afd3...: layup/rules: failure` (`red-base.txt`) |

## Green

| Run | Result |
| --- | --- |
| `TestTheSpecificationTaskEndToEnd` alone, both fixes in | `ok ... 21.371s`, exit 0 |
| control before the mutants and after them (`mutants-out.txt`) | `ok ... 11.437s` and `ok ... 21.422s` |

## S7. The checks

| Check | Time | Exit |
| --- | --- | --- |
| `go build ./...` | 10:13:06 | 0 |
| `go vet ./...`, `-tags=integration`, `-tags=e2e` | 10:13:06 to 10:13:07 | 0, 0, 0 |
| `gofmt -l .` | 10:13:07 | 0, no file listed |
| `go test -count=1 ./...` (unit) | 10:13:08 | 0 |
| `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...` (output only in `/home/layup/.cache/layup-session/T-3xvc-full-integration.txt`) | 10:13:14 to 10:17:31 | 0, every package `ok` |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | 10:17:3x to 10:19:01 | 0; the package `cmd/layup` **85.720s** (target under 4 min, hard limit 10m; the baseline before this task was 64.407s) |
| `adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests`, `nested-checkout-check`, `setup-check`, `freeze-check/run-all` | after 10:19:01 | 0 each |
| `git diff --check origin/main HEAD` | after 10:19:01 | 0 |
| `docs/tasks/task-state.sh` | after 10:19:01 | `T-3xvc #368 in progress` |

## Round 1 of #368 (material: two test-oracle findings)

**Finding 1, the statuses at the merge.** The fake now keeps, at each merge request, a snapshot of the statuses at the requested SHA (`atMerge`), and `loopMerged` asserts the four newest App statuses from that snapshot through `notGreen`. Red first: the probe test `TestTheMergeOracleRefusesStatusesPostedAfterTheMergeRequest` (a snapshot with only `layup/verify`, the other three posted after the request) against a stub `notGreen` that accepted everything: exit 1, `a snapshot with layup/verify alone: []; want the other three contexts` and `a newer failure and a newer status of another creator: []` (`red-oracle.txt`). Green: with `notGreen` written, the probe passes, and so does the real run (`ok 21.407s`, both tests together).

**Finding 2, C1.** Stage 6 now asserts, before the comparison of the endpoint, that the log holds two merge requests. C1 rerun with the existing method (control green, build ok, vet ok, restore hash equal, control green after): killed by `loop_e2e_test.go:444: stage 6: 0 merge requests in the log; want 2, one for each task of Shape`, before the row comparison at `:447` (`mutants.md`). E2 is now also killed there (`1 merge requests`).

Checks on the head `dd02f40a` plus this record: build, vet (default, integration, e2e), gofmt, unit, all exit 0 (10:35); `go test -count=1 -tags=e2e -timeout 10m ./...` exit 0 at 10:36:52, package `cmd/layup` **85.678s**; the six lints, `freeze-check/run-all` and `git diff --check origin/main HEAD` exit 0. No product code changed, so the integration level was not rerun.

## The pre-freeze check

`sh docs/tests/freeze-check.sh --task T-3xvc --base 30e2c0d48b5a0f8c9e6132e9232ce81bf3779a49 --budget 850/12 --evidence runs/T-3xvc/test-runs.md --evidence runs/T-3xvc/mutants.md --evidence runs/T-3xvc/conditions.md` on the clean committed head `c826afa9` (main `30e2c0d` already in it), 10:21:25 UTC, exit 0. The first run failed K3 and K6 on the form of the evidence (selectors unquoted, code spans that were not paths); fixed in the evidence files, not in code.

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
ok K5: projection: figure 634/10, maximum 850/12, tolerance 170/4, limit 1020/16, margin 386/6 (the figure at landing is still measured at landing)
ok K6: conditions table
```
