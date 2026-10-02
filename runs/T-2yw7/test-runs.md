# T-2yw7: the test runs

The red step of each test of `T-2yw7` (#80), before the code that makes it
green (gate step 3). The tests ran on skeletons of the new code: the names of
`parse`, `usage`, `exitCode`, `newProgress`, `step`, `end`, `checkInputs` and
`sameBytes`, with no behaviour (each one gives its zero value or does nothing),
and on the command code of the base `0506c13`. So each test fails on its own
assertion. Outputs are cut to the relevant lines ("…"). Host: macOS, `go1.27.1`,
2026-10-02. The runs used the first layout of the test files: `exit_test.go`,
`input_rule_test.go` and `input_rule_checker_test.go`. Their tests and the scan
moved, with no change, into `cli_test.go`, `rules_test.go` and
`rules_checker_test.go`, so that the diff stays inside the file count of the
budget. The green runs are at the end.

Tests that pass on the skeleton guard behaviour that the base already has:
`TestAnInputErrorPrintsTheReasonAndNoUsage`, `TestTheUsageOfLayupListsItsCommands`,
`TestProgressEndWithNoStepPrintsNothing`, the good cases of `TestSameBytes` and
of `TestExitCode`, and the end-to-end `TestVersion` (the binary of the base
prints its version).

```text
$ go test -count=1 ./internal/cli/                                        # run 1
--- FAIL: TestParseTakesTheLongestMatchOfTheCommandWords
    args_test.go:33: ["setup" "verify" "w"]: command [], arguments [], error <nil>; want "setup verify", ["w"]
    … the same for the other 4 cases
--- FAIL: TestParseTakesAFlagBeforeOrAfterThePositionalArguments
    args_test.go:46: ["gate" "r" "--base" "b" "--head" "h"]: arguments [], flags map[], error <nil>
    … the same for the other 2 cases
--- FAIL: TestParseGivesTheReasonOfEachUsageError
    args_test.go:75: []: error <nil>, want "no command"
    args_test.go:75: ["version" "extra"]: error <nil>, want "extra argument \"extra\""
    args_test.go:75: ["gate" "r" "--base=" "--head" "h"]: error <nil>, want "flag --base has no value"
    … the same for the other 14 cases
--- FAIL: TestTheUsageListsEachCommandWithItsArgumentsAndTheExitCodes
    args_test.go:85: the usage has no "usage: layup <command>":
    … the same for the other 5 parts
--- FAIL: TestVersionWithAnArgumentIsAUsageError
    cli_test.go:76: exit 0, stdout "layup 0.1.0-dev\n"; want 2 and nothing
--- FAIL: TestAUsageErrorPrintsTheReasonAndTheUsageOnStandardErrorOnly
    cli_test.go:86: stderr "layup: usage: layup psb check FILE\n", want it to start with "layup: extra argument \"b.md\"\n\nusage: layup "
--- FAIL: TestExitCode
    exit_test.go:24: exitCode(["pass" "fail"]) = 0, want 1
    exit_test.go:24: exitCode([]) = 0, want 1
    … the same for the other 5 cases of 1
--- FAIL: TestProgressPrintsALinePerStepAndABeatEveryTenSeconds
    progress_test.go:38: 0 tickers started, want a ticker for step 1
--- FAIL: TestNewProgressBeatsEveryTenSeconds
    progress_test.go:69: every 0s, command ""; want 10s and setup
$ go test -count=1 ./cmd/layup/                                           # run 2
--- FAIL: TestSameBytes
    bytes_test.go:19: sameBytes("layup 0.1.0-dev\n", "layup 0.1.1-dev\n") = <nil>, want same false
    bytes_test.go:19: sameBytes("layup\n", "layup") = <nil>, want same false
--- FAIL: TestCheckInputsFindsEachReadOfTheEnvironmentOrTheStandardInput
    input_rule_test.go:31: the findings:
        want:
        input rule: internal/cli/a.go:5 reads os.Getenv
        … and the other 6 findings
$ go test -count=1 -tags=integration -run TestInputRule ./cmd/layup/      # run 3
--- FAIL: TestInputRule
    rules_integration_test.go:141: the fixture gives
        want only the finding input rule: internal/cli/cli.go:9 reads os.Getenv
$ go test -count=1 -tags=e2e -run 'TestUsageErrors|TestVersion' -v ./cmd/layup/   # run 4
--- PASS: TestVersion
    usage_e2e_test.go:35: layup []: exit 2, stdout "", stderr "usage: layup <command>\n…"; want 2, nothing, and the reason with the usage
    usage_e2e_test.go:35: layup ["version" "extra"]: exit 0, stdout "layup 0.1.0-dev\n", stderr ""; want 2, nothing, and the reason with the usage
    usage_e2e_test.go:35: layup ["version" "--x"]: exit 0, stdout "layup 0.1.0-dev\n", stderr ""; want 2, …
    usage_e2e_test.go:35: layup ["psb"]: exit 2, stdout "", stderr "layup: usage: layup psb check FILE\n"; want 2, …
    … the same for ["psb" "check"], ["psb" "chek" FILE], ["psb" "check" FILE FILE], ["psb" "check" "--x" FILE]
--- FAIL: TestUsageErrors
```

Run 5 is on the real scan, with the allowance of `environ` of `internal/git`
taken out for the run. The test then fails on the real module, so it reads the
real files of `internal/git`; with the allowance back, it passes.

```text
$ go test -count=1 -tags=integration -run TestInputRule ./cmd/layup/      # run 5
--- FAIL: TestInputRule
    rules_integration_test.go:137: the module breaks the input rule:
        input rule: internal/git/git.go:68 reads os.LookupEnv
```

## The green runs

The runs on the tree of the commit that adds this section, with the code of
`T-2yw7` and the merged test files.

| Command | Result |
| ------- | ------ |
| `go build ./...` | exit 0 |
| `go vet ./...`; `go vet -tags=integration ./...`; `go vet -tags=e2e ./...` | exit 0 each |
| `gofmt -l .` | no file |
| `go test -count=1 ./...` | `ok` × 5 packages |
| `go test -count=1 -tags=integration ./...` | `ok` × 5; `TestInputRule` and `TestPackageRules` in `cmd/layup` |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | `ok` × 5; `TestVersion` and `TestUsageErrors` pass (`-v`), with one build of the binary in `TestMain` |
| `go test -race -count=1 ./internal/cli/` | `ok`: the beats of the progress lines print from their own goroutine with no data race |
