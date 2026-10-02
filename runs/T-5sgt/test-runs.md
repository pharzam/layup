# T-5sgt: the test runs

The runs of `T-5sgt` (#82). Host: macOS, `go1.27.1`, `git` 2.54.0, 2026-10-02.
Outputs are cut to the relevant lines ("…").

## A measurement before the plan

`git ls-tree -r -z --full-tree --end-of-options REV -- PATH` on a test
repository: a directory gives each file under it, a file gives its mode
(`100644`, `100755`), a symbolic link gives `120000`, and a path that the
commit does not have gives no line and exit 0. So one call gives the files and
the modes of a `config` path at the base (D3).

## The red runs

Each test ran before the code that makes it green: the tests of `LsTree` on a
skeleton of the call (it starts no `git`); the unit and integration tests of
`internal/gate` on a skeleton of the package (its names, with no behaviour); the
tests of the command in `internal/cli` before its row; the end-to-end
scenarios on the binary of the base `f394433`, which has no command `gate`. The
integration and end-to-end runs used a scratch clone of the base with the new
test files. So each test fails on its own assertion.

```text
$ go test -count=1 -run 'TestEachCallRunsItsVerb|TestLsTreeReadsEachEntry' ./internal/git/     # run 1
    --- FAIL: TestEachCallRunsItsVerb/ls-tree
        git_test.go:82: 0 starts of git, want 1
--- FAIL: TestLsTreeReadsEachEntry
    git_test.go:210: LsTree = [], <nil>; want [{"100644" "blob" … "layout/a test.go"} {"100755" "blob" … "run.sh"}]
$ go test -count=1 -tags=integration -run TestLsTree ./internal/git/                            # run 2
--- FAIL: TestLsTree
    git_integration_test.go:368: LsTree("layout") = [], <nil>; want ["100644 blob layout/a_test.go" "100644 blob layout/b/c.txt"]
    … the same for run.sh (100755), link.sh (120000) and layout/a_test.go
$ go test -count=1 ./internal/gate/                                                            # run 3
--- FAIL: TestScopePatterns
    scope_test.go:30: "./*.go" matches "x.go": false, want true          … and 6 more
--- FAIL: TestOtherFormsOfAPatternAreErrors
    scope_test.go:38: parsePattern("*.go"): no error                       … and the other forms
--- FAIL: TestReadManifestGivesTheKindsInTheirOrder
    manifest_test.go:19: the kinds []
--- FAIL: TestReadManifestRefusesEachMalformedForm
    manifest_test.go:38: a wrong header: no error                          … and the other 9 forms
--- FAIL: TestDecideActiveByTheFirstLineThatMatches
    result_test.go:40: a missing tool:  "" with 0 runs; want not-active "tool not found: go" with 0
    … the same for the other 5 lines
--- FAIL: TestDecidePendingNamesTheFirstChangedProductPath
    result_test.go:48: a changed product path:  ""
--- FAIL: TestACheckThatDidNotRunNeverPasses
    result_test.go:79: active, files ["x.go"]: ; want pass
--- FAIL: TestRunGivesOneRowPerKindInTheOrderOfTheManifest
    run_test.go:170: the table {Base: Head: Rows:[]}
--- FAIL: TestRunGivesAnInputErrorAndNoTable
    run_test.go:215: git not found: 0 rows, error <nil>, 0 scratch trees; want an *InputError, no row and no scratch tree
    … the same for the other 8 inputs
    run_test.go:222: sh not found: <nil>; want an *InputError
--- FAIL: TestAFailureOfTheScratchTreeOrTheDiffGivesNotActiveRows
    run_test.go:248: the error of git is not on standard error: []
--- FAIL: TestACleanupFailureComesWithTheWholeTable
    run_test.go:268: 0 rows, <nil>; want 6 rows and a *CleanupError
--- FAIL: TestTheOverlayPutsTheGateFilesOfTheBase
    run_test.go:306: layout/a_test.go in the scratch tree: ""; want "package layout // base\n"   … and 5 more
--- FAIL: TestTheOverlayDoesNotWriteOutOfTheTree
    run_test.go:322: 0 rows, <nil>; want 6 rows
$ go test -count=1 -tags=integration ./internal/gate/   # on the skeleton                      # run 4
--- FAIL: TestTheSchemaBlocks
    gate_integration_test.go:32: schema gate-manifest: the name: the block has "gate-manifest"; the Go schema has "" …
--- FAIL: TestRunOnARealRepository
    gate_integration_test.go:125: the table {Base: Head: Rows:[]}
--- FAIL: TestAHookOfTheRepositoryDoesNotRun
    gate_integration_test.go:160: [], <nil>; want the row ok, pass
--- FAIL: TestTheRevisionsAndTheManifestOnARealRepository
    gate_integration_test.go:175: the head of a tag: , <nil>; want the commit 23aa338e…
--- FAIL: TestARenameAndADeleteChangeAProductPath
    gate_integration_test.go:198: [], <nil>; want fail on internal/a.go, the first path that git names
--- FAIL: TestAFailedScratchTreeGivesAReasonWithNoPath
    gate_integration_test.go:211: [], <nil>; want not-active, scratch tree: add failed
$ go test -count=1 ./internal/cli/   # before the row of gate                                  # run 5
--- FAIL: TestGateUsageErrors
    gate_test.go:48: ["gate"]: exit 2, stdout "", stderr "layup: unknown command \"gate\"\n\nusage: …"   … and the other 4
--- FAIL: TestGatePrintsTheTableAndGivesItsExitCode
    gate_test.go:69: ["pass" "clear"]: exit 2, stdout ""; want 0 and the table   … and the other 2
--- FAIL: TestGateGivesTwoOnAnInputErrorOrALeftoverScratchTree
--- FAIL: TestTheUsageListsGate
$ go test -count=1 -tags=e2e -run TestGate ./cmd/layup/   # the binary of the base              # run 6
--- FAIL: TestGateOnAGoRepository
    gate_e2e_test.go:92: each kind passes or is clear: exit 2, stdout: …
--- FAIL: TestGateNeverPassesACheckThatDidNotRun
    gate_e2e_test.go:121: a missing tool: exit 2, stdout ""; want 1 and not-active
```

**Two tests made strict after a first red run.** `TestTheOverlayDoesNotWriteOutOfTheTree`
and `TestAHookOfTheRepositoryDoesNotRun` passed on the skeleton, because the
skeleton makes no row and runs nothing. The first now also asks for the six
rows; the second runs a control first (a plain `git worktree add` in the same
repository runs the hook, so the test can fail) and asks for the row of the run.
Both then fail on the skeleton (runs 3 and 4 above show them). The exit-2 cases
of `TestGateNeverPassesACheckThatDidNotRun` also refuse a usage text, so an
unknown command cannot pass them.

**One test expectation was wrong.** `TestRunOnARealRepository` first wanted the
reason `pending: product path changed: y.go`; `git`'s order names
`layout/a_test.go` first, which the scope `./*.go` also matches. The test now
wants that path; the code did not change.

## The green runs

On the tree of the commit that adds this file.

| Command | Result |
| ------- | ------ |
| `go build ./...` | exit 0 |
| `go vet ./...`; `go vet -tags=integration ./...`; `go vet -tags=e2e ./...` | exit 0 each |
| `gofmt -l .` | no file |
| `go test -count=1 ./...` | `ok` × 7 packages |
| `go test -count=1 -tags=integration ./...` | `ok` × 7; the integration tests of `internal/gate` and `TestLsTree`; `TestPackageRules` holds `internal/gate` and its `sh -c`; `TestInputRule` passes |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | `ok` × 7; `TestGateOnAGoRepository` and `TestGateNeverPassesACheckThatDidNotRun` pass in about 3 s with the host's `GOCACHE` |
| `go test -race -count=1 ./internal/gate/ ./internal/cli/` (also with `-tags=integration` for `internal/gate`) | `ok` |
