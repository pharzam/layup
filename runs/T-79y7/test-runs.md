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

## The green runs

On the tree of the commit that adds this file.

| Command | Result |
| ------- | ------ |
| `go build ./...`; `go vet` with no tag, `-tags=integration` and `-tags=e2e`; `gofmt -l .` | exit 0; no file |
| `go test -count=1 ./...` | `ok` × 11 packages |
| `go test -count=1 -tags=integration ./...` | `ok` × 11; the two blocks, the run on a real work area (one commit, none on a rerun, a hostile host, `sh -n`), the exit codes 0 to 3, `TestPackageRules` and `TestInputRule` pass |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | `ok` × 11 |
| `go test -race -count=1 ./internal/setup/ ./internal/cli/`, also with `-tags=integration` | `ok` |
