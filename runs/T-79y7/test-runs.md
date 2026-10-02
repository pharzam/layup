# T-79y7: the test runs

The runs of `T-79y7` (#85). Host: macOS, `go1.27.1`, `git` 2.54.0, 2026-10-02.
A path of a temporary directory shows as `<tmp>/`; "…" marks a cut part of an
output.

## The red runs

Each red run is of a test in its final form, on a skeleton of the code (the
names, with no behaviour), or on a copy of the branch or of the base.

1. **`internal/setup`, unit** (the skeleton: `Run` gives an empty result, and
   each helper gives nothing): each of the 14 tests fails. Two tests were
   changed before this run: one panicked on the empty result (it now checks the
   length first), and one wanted no error for a changed `M-` row when S10 is
   done, which reads it (now an error).

   ```text
   $ go test -count=1 ./internal/setup/
   --- FAIL: TestARunResumesAfterItsDoneSteps
       setup_test.go:107: the calls [], want ["S05" "S06" "S07" "S08" "S09" "S14" "S10" "S11" "S12" "S13" "S15"]
   --- FAIL: TestTheOutcomes
       setup_test.go:153: a fail: the rows [], want ["fail root tree differs" "not-active not run: S03 did not pass" …]
   --- FAIL: TestTheHandOff
       setup_test.go:168: the done row of S13: "", false
   … and TestTheProseGroupStopsOnce, TestTheStopTableOrder, TestTheTables, TestTheAnswersRule,
     TestCheckAsked, TestTheAnswersOfADoneStep, TestMarkerID, TestTheCommandsFile,
     TestTheCommitOfAStep, TestTheInputs, TestTheStubs
   ```

   The case of a commit that fails (`TestTheCommitOfAStep`) was added after
   this run, before the code: it decides that the step fails with no `done`
   row, so a rerun does it again.

2. **`internal/setup`, integration**, on a copy of `34cae00` where `Run` gives
   an empty result and the two Go schemas have no column:

   ```text
   $ go test -count=1 -tags=integration ./internal/setup/
       setup_integration_test.go:30: schema setup-steps: the location: the block has "stdout"; the Go schema has ""
           the number of columns: the block has "4"; the Go schema has "0"
       … the same for setup-stop
       setup_integration_test.go:91: the results ""
       setup_integration_test.go:94: 0 new commits on layup-setup, want 1
   ```

3. **The command row**, with the seams `setupRun` and `setupSteps` in place and
   no row `setup` in the command table:

   ```text
   $ go test -count=1 -run TestSetupCommand ./internal/cli/
       setup_test.go:40: ["setup"]: exit 2, stdout "", stderr "layup: unknown command \"setup\"…"; want 2, nothing and "missing argument WORK" with the usage
       setup_test.go:63: done and a hand-off: exit 2, stdout "", WORK ""; want 0 and the table of w
       setup_test.go:63: a stop: exit 2, stdout "", WORK ""; want 3 and the table of w
       setup_test.go:74: the usage has no setup WORK beside setup verify WORK: …
   $ go test -count=1 -tags=integration -run TestSetupExitCodesOnAWorkArea ./internal/cli/
       setup_integration_test.go:46: a stop: exit 2, stdout ""; want 3 and only the stop table
       setup_integration_test.go:59: a fail: exit 2, stdout ""; want 1
       … the same for the resumed run and the record of another form
   ```

4. **The end-to-end scenarios**, on the binary of a copy of the base
   `1ef71b4`, which has no command `setup`:

   ```text
   $ go test -count=1 -tags=e2e -run 'TestSetup$' ./cmd/layup/
       setup_e2e_test.go:27: layup ["setup"]: exit 2, stdout "", stderr "layup: unknown command \"setup\"…"
       … the same for a flag and a missing work area
       setup_e2e_test.go:33: a new work area: exit 2, stdout "", stderr "layup: unknown command \"setup <tmp>/…\"…"
   ```

The pre-commit hook refused the first commit of the runner: check `markers`
found the marker of `TestMarkerID` as two characters in the Go source; it is
now written as escapes (K11).

## The red runs of the fixes of review round 1

Round 1 (`53b3a99`) gave five material findings. First, a move with no change
of behaviour, under green tests: the hash of the answers went from
`internal/setup` to `work.AnswersHash`, so `internal/standin` can give its
stand-in record the row; the skeletons `git.Branch` (no `git` start) and the
seam `branch` of the runner were added. Then each test below, in its final
form, failed on that tree:

```text
$ go test -count=1 ./internal/setup/
    setup_test.go:218: S13 with no command: <nil>, done true, commits ["chore: setup S13 by …"], …    (finding 3)
    setup_test.go:218: an outcome operator with an empty list: <nil>, done true, commits [], …          (finding 3)
    setup_test.go:312: an S01- question that S01 does not ask (finding 4 of round 1): <nil>, calls ["S01" … "S15"]; want an input error …
    setup_test.go:363: the row S01 answers.sha256 [… "computed" "sha256 of the rows S01- Q- of inputs/answers.tsv"]; want … ref sha256 inputs/answers.tsv S01- Q-   (finding 5)
    setup_test.go:393: a done S01 with no hash: <nil>, calls ["S02" … "S15"]; want an input error before any step   (finding 1)
    setup_test.go:393: a done S10 with no hash: <nil>, calls []; want an input error before any step
    setup_test.go:440: a run with no command: <nil>, commands.sh "echo stale\n"; want it written again, empty         (finding 3)
    setup_test.go:476: a target on refs/heads/main: <nil>, done true, commits ["chore: setup S05 by …"], …            (finding 2)
    setup_test.go:476: a target on detached: <nil>, done true, commits ["chore: setup S05 by …"], …
$ go test -count=1 ./internal/git/
    --- FAIL: TestEachCallRunsItsVerb/symbolic-ref_--quiet
        git_test.go:84: 0 starts of git, want 1
$ go test -count=1 -tags=integration -run 'TestBranch$' ./internal/git/
    git_integration_test.go:187: after init: Branch "", <nil>; want "refs/heads/main"
    … the same after switch -c, and no error for a detached HEAD
$ go test -count=1 -tags=integration ./internal/setup/
    setup_integration_test.go:137: a target on main: <nil>, the rows [… {"S05" "layup-setup" "done" "x.md"} …]; want S05 fail
    setup_integration_test.go:140: the branches moved: "1bd7b8c…\na6a1bea…", then "3ca5c9f…\na6a1bea…"
$ go test -count=1 -tags=integration -run TestSetupExitCodesOnAWorkArea ./internal/cli/
    setup_integration_test.go:78: a done S01 with no hash: exit 1, stdout "step\tactor\tresult\tevidence\nS01\tlayup-setup\tdone\t…"
$ go test -count=1 -tags=e2e -run 'TestSetup$' ./cmd/layup/
    setup_e2e_test.go:37: layup ["setup" "<tmp>/…"]: exit 1, stdout "…S01\tlayup-setup\tnot-active\tnot built yet\n…"
```

The second line of the integration run of `internal/setup` is finding 2 itself:
the commit of S05 moved `main`. The tests of the hand-off and of the runs that
reach S13 now give S13 a command (`commands` of `setup_test.go`, `doneSteps`,
`stubSteps`), and the stand-in record has the row `S01 answers.sha256`.

## The green runs

On the tree of the commit that adds the fixes of round 1.

| Command | Result |
| ------- | ------ |
| `go build ./...`; `go vet` with no tag, `-tags=integration` and `-tags=e2e`; `gofmt -l .` | exit 0; no file |
| `go test -count=1 ./...` | `ok` × 11 packages |
| `go test -count=1 -tags=integration ./...` | `ok` × 11; the two blocks, the run on a real work area (one commit, none on a rerun, a hostile host, `sh -n`), no commit off `layup-setup`, `TestBranch`, the exit codes 0 to 3, `TestPackageRules` and `TestInputRule` pass |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | `ok` × 11 |
| `go test -race -count=1 ./internal/setup/ ./internal/cli/ ./internal/git/`; `./internal/setup/ ./internal/cli/` with `-tags=integration` | `ok` |
