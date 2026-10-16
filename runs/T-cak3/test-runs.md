# T-cak3 — test runs (row 66a of M2e, the task register)

Author: Claude Sonnet 5.5 on Claude Code. Base `f12da6a`. Every time is UTC, 2026-10-16, taken with `date -u` after the command it follows (the start of a red or green run from a file written before it); every exit code is that of the command itself, run alone with its output in a file, never through a pipe. The plan is comment 6210010288 on #357, its review 6210141165 (`approve`), the author's answer 6210153718.

## Step T0 and T: red (07:14 to 07:16Z)

The stub `internal/run/register.go` (signatures, bodies that refuse with `not built`, and a `loopPermissionsStep` that says `done`) and the field `granted` of the state (never set), with the 11 unit and 6 integration tests written against it. No commit holds the red state: the `pre-commit` hook runs the unit level and refuses it (guardrails §2, "The hook refuses a red commit"), so the red runs are these.

| Run | Command | Started | Ended | Exit |
| --- | ------- | ------- | ----- | ---- |
| unit, red | `go test ./internal/run -run '<the 11 unit tests>'` | 07:14:46 | 07:14:57 | 1 |
| integration, red | `TMPDIR=… go test -tags=integration ./internal/run -run '<the 6 tests>'` | 07:15:00 | 07:15:50 | 1 |

All 11 unit tests failed by assertion: the register tests with `not built` where a table or a refusal at a line and column was wanted; `TestTheStepForgeOfStartAndOfTheRestartRecordsTheLoopPermissionsAndDoesNotFailOnThem` with the grant empty (`map[]`) where the installation's permissions were wanted (red because `granted` was not yet set, not green on arrival as the plan expected); `TestTheLoopPermissionsStepNamesEachMissingPermission` with the stub's `done` for seven of nine cases. The six integration tests failed after a successful Start and restart clone of the world (the world is not the cause): the schema (`the name: the block has "tasks"; the Go schema has ""`), the permissions step (`done`, want `fail`), the refusal (`not built` does not name `pull_requests: write`), and `not built` at each write. Every red is for the reason of its row in the plan.

## Step C: green (07:17 to 07:20Z)

| Run | Started | Ended | Exit | Note |
| --- | ------- | ----- | ---- | ---- |
| unit, `go test ./internal/run` | 07:17:48 | 07:17:55 | 1 | two tests: a list read back as an empty, non-nil slice where nil was wanted (`strings.Fields("")`); fixed with `words` |
| unit | 07:17:55 | 07:18:07 | 0 | |
| focused integration | 07:18:10 | 07:18:57 | 1 | `TestTheRegisterWriteIsFencedByTheLease`: a defect of the test, a run ID of 16 hexadecimal characters was needed for the lease row |
| focused integration | 07:18:57 | 07:20:02 | 0 | |
| `go test -tags=integration ./internal/tsv/` (`tasks` moved to `built`) | 07:20:08 | 07:20:25 | 0 | |
| `go test -tags=integration ./cmd/layup/ -run TestPackageRules` | 07:20:08 | 07:20:25 | 0 | green before and after the code; the May import cell of `internal/run` already held `internal/handoff` |

## Step M: 30 mutants (07:23 to 07:28Z)

Each mutant is applied alone to one file, keeps the code building (a condition joined with `&& false`, never deleted; three first versions did not build, M03, M21 and M22, and were rewritten: a "build failed" is counted as not run, guardrails §2), runs the tests of its row, and is undone. The script is a one-off, not kept. A control run of the unit level before them was exit 0. All 30 exit 1 with `built=True` and failing tests by assertion:

| Mutant | Caught by |
| ------ | --------- |
| M01 two columns of a row swapped | `TestARegisterIsWrittenAndReadBackWithItsDashes`, `TestARegisterThatReadWouldRefuseIsNotWritten`, `TestAnExistingTransitionTableIsKeptAndAnInvalidOneRefusesTheTask`, `TestTheFilesOfTheFirstTaskAreTheRegisterAndTheDefaultTable` |
| M02 no — for a pull request that is not opened | `TestARegisterIsWrittenAndReadBackWithItsDashes` |
| M03 the writer skips the read-back | `TestACycleOfPredecessorsIsRefusedAndADiamondIsNot`, `TestARegisterThatReadWouldRefuseIsNotWritten` |
| M04 requirements may be none | `TestARegisterThatReadWouldRefuseIsNotWritten`, `TestARowWithNoRequirementIsRefusedNamingItsLineAndColumn` |
| M05 a predecessor may be no task | `TestAPredecessorThatIsNotATaskIsRefusedWhereverTheTaskIsInTheTable` |
| M06 a predecessor must be an earlier row | `TestACycleOfPredecessorsIsRefusedAndADiamondIsNot`, `TestAPredecessorThatIsNotATaskIsRefusedWhereverTheTaskIsInTheTable`, `TestAnExistingRegisterThatReadRefusesStopsTheWrite` |
| M07 no cycle check | `TestACycleOfPredecessorsIsRefusedAndADiamondIsNot`, `TestAnExistingRegisterThatReadRefusesStopsTheWrite` |
| M08 only a task that is its own predecessor is a cycle | `TestACycleOfPredecessorsIsRefusedAndADiamondIsNot`, `TestAnExistingRegisterThatReadRefusesStopsTheWrite` |
| M09 a diamond is a cycle (a control) | `TestACycleOfPredecessorsIsRefusedAndADiamondIsNot` |
| M10 AddTask appends to the table given | `TestAddTaskLeavesTheTableItWasGivenAlone` |
| M11 a task may take the ID of a row | `TestAddTaskLeavesTheTableItWasGivenAlone` |
| M12 the default table is never written | `TestTheFilesOfTheFirstTaskAreTheRegisterAndTheDefaultTable` |
| M13 the table is always written from the default | `TestAnExistingTransitionTableIsKeptAndAnInvalidOneRefusesTheTask`, `TestTheFilesOfTheFirstTaskAreTheRegisterAndTheDefaultTable` |
| M14 CheckTable is skipped | `TestAnExistingTransitionTableIsKeptAndAnInvalidOneRefusesTheTask` |
| M15 an existing register is not read | `TestAnExistingRegisterThatReadRefusesStopsTheWrite`, `TestTheFilesOfTheFirstTaskAreTheRegisterAndTheDefaultTable` |
| M16 the permissions step ignores the result | `TestTheLoopPermissionsStepNamesEachMissingPermission` |
| M17 the permissions step names only the first | `TestTheLoopPermissionsStepNamesEachMissingPermission` |
| M18 a read grant is taken for a write one | `TestTheLoopPermissionsStepNamesEachMissingPermission` |
| M19 a state that read no installation passes | `TestTheLoopPermissionsStepNamesEachMissingPermission` |
| M20 step forge fails on a loop permission | `TestTheStepForgeOfStartAndOfTheRestartRecordsTheLoopPermissionsAndDoesNotFailOnThem` |
| M21 the grant shares the map of the installation | `TestTheStepForgeOfStartAndOfTheRestartRecordsTheLoopPermissionsAndDoesNotFailOnThem` |
| M22 the register write does not check the permissions | `TestTheLoopWritesNothingWhenThePermissionsAreMissing` |
| M23 the register write is not fenced | `TestTheRegisterWriteIsFencedByTheLease` |
| M24 each file in its own commit | `TestTheFirstTaskIsOneRecordsCommitOfTheRegisterAndTheDefaultTable` |
| M25 another subject of the commit | `TestTheFirstTaskIsOneRecordsCommitOfTheRegisterAndTheDefaultTable` |
| M26 opened is not stamped from the clock | `TestTheFirstTaskIsOneRecordsCommitOfTheRegisterAndTheDefaultTable` |
| M27 the Go schema has two columns swapped | `TestTheSchemaOfTasksEqualsItsBlock` |
| M28 the table of the Operator is rewritten (integration) | `TestASecondTaskWritesOnlyTheRegisterAndKeepsTheOperatorsTable` |
| M29 the existing register is not read (integration) | `TestARefusedTaskWritesNothing`, `TestASecondTaskWritesOnlyTheRegisterAndKeepsTheOperatorsTable` |
| M30 internal/handoff is not in the May import cell of internal/run | `TestPackageRules` |

The last row, M30, is a document mutant: removing `internal/handoff` from the May import cell of `internal/run` makes `TestPackageRules` red, which shows the check reads the cell.

## Step V: the checks before the first push (07:28 to 07:30Z)

All exit 0, each alone: `gofmt -l cmd internal` (no output) 07:28:38; `go build ./...` 07:28:39; `go vet ./...` 07:28:41; `go test ./...` 07:28:49; `go vet -tags=integration ./...` 07:28:52; the focused integration run of `internal/run` (the six tests of the task register) 07:29:12; `go test -tags=integration ./internal/tsv/` 07:29:12; `go test -tags=integration ./cmd/layup/ -run TestPackageRules` 07:29:15; `sh docs/adr/adr-lint.sh` 07:29:18; `sh docs/prd/prd-lint.sh` 07:29:18; `sh docs/links/link-lint.sh` 07:29:31; `sh docs/tests/run-discipline-tests.sh` 07:29:38; `sh docs/setup/setup-check.sh` 07:30:10 (after `git add` of the new files); `git diff --check` and `git diff --cached --check` 07:30:10. The full integration level (`go test -tags=integration ./...`) has **not** been run: it waits for the coordinator's go-ahead (it is serialised with another row's run on this host).
The focused integration run was repeated with `-count=1` at 07:30:51, exit 0 (21 s).
