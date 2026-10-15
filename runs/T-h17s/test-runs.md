# T-h17s — the runs of the tests

UTC, 2026-10-15; times from the clock of the host.

## Step T: the tests first, seen red (18:34 to 18:35Z)

With a stub of the package (the types and the functions returning zero values) so that each test fails on its assertion:

- `go test ./internal/handoff/`, each test alone: FAIL for `TestTheDefaultRowsAreTheThirteenOfTheSpecification`, `TestDefaultGivesACopy`, `TestTheListsAreTheWordsOfTheSpecification`, `TestEachItemIsAWordOfTheDetailOfAnEvent`, `TestWriteThenReadGivesTheDefaultTable`, `TestReadTableRefuses`, `TestCheckTable`, `TestACheckedHandoffIsValidOnlyByItsTable`, `TestTheFirstFailingItemNamesTheDetail`, `TestWhatAnAgentWritesNeverMakesAHandoffValid`, `TestCheckUsesTheTableItIsGiven`. `TestEachDefaultPairNeedsItsOwnItemsAndKinds` passed on the stub because `Default()` was empty; a guard (`len(Default()) == 13`) was added, so it fails there too.
- `go test -tags=integration ./internal/handoff/`: FAIL for `TestTheSchemaEqualsItsBlock`, `TestTheDefaultRowsAreTheTableOfTheSpecification` and `TestTheListsAreThoseOfTheSpecification`; `TestThePackageHasItsRowInTheTableOfM2e` passed (the specification already had the row).
- `go test ./cmd/layup/ -run TestReadTable` with a table of `M2e` in `goodTable`: `TestReadTable` FAIL (the row of `internal/handoff` is not read) and `TestReadTableRefusesWhatItCannotRead` FAIL for "no table of M2e", "a row of M2e for a package of another table" and "a row of M2e that stands for a row of phase 1" (no error, 9 rows).
- The block test (`internal/tsv`, `transitions` still in `notYetBuilt`) is green until the Go schema exists; its change is the move of the name, checked by `TestTheSchemaEqualsItsBlock`.

## The tests of conditions 1 to 3 of the plan review, red against the first version of `Check`

The first version read `schema` and `artifacts` only when the row listed them, and only for the row's kinds. Written after the plan review (comment 6199608521), run at 18:53:13Z with that version of `Check` put back in a copy of the file, then restored: `TestAReplacementRowCannotLetABadResultThrough` FAIL on six cases (a bad schema, a named artifact whose hash differs, an extra artifact that is bad, an extra artifact of an unknown kind, an extra comment with no comment, the same handoff written by a session). With the final `Check`: ok. (The order of the work was the change of `Check`, then these tests; this run is the red evidence that they bite.)

## Mutant runs (`go test -tags=integration ./internal/handoff/`), 18:39 to 18:50Z

26 mutants of `check.go` and `table.go`, each applied alone to the file and undone: all 26 caught by a test (a named hash not compared with the found one; the forge's body, or the session's file, not compared; a comment hash that need not be a SHA-256; a hash that need not be a full object name; the schema result ignored; a code-written handoff that needs a schema; `artifacts` taken from the caller's claim; an artifact of an unknown kind; a required kind that may be absent; an item with a false result that passes; a pair with no row that is valid; the comment's copy not named; a review record that is not a comment; a default pair that may be missing; a pair listed twice; a kind or an item outside the list; a row with no artifact or no item; `ReadTable` not checking the words; a default row that lacks an item or has a wrong one; `Kinds` that does not copy). The script is not kept in the tree (a one-off evidence command).

## Green, at the head before the freeze

- `go build ./...`, `go vet ./...`, `go test -count=1 ./...` and `go test -count=1 -tags=integration ./...` (TMPDIR on the disk): see the end of this file.
- At the head: `go build ./...` exit 0; `go vet ./...` exit 0; `go test -count=1 ./...` exit 0; `go test -count=1 -tags=integration ./...` 18:53:19 to 19:03:53Z: 20 packages ok and `internal/run` timed out at the default 10 minutes (the host runs other tasks in parallel; the memory of the host says 6 to 9 minutes alone). Rerun alone with `-timeout 1500s`, 19:04:15 to 19:15:03Z: ok, 644.763 s. A first full run at 18:41:41Z failed on `TestAGoodModuleKeepsTheRules` and `TestEachRuleAndColumnFindsItsBreach` (my fixture of `goodModule` added two packages that had no row: "table: internal/records has no row"); fixed in `cmd/layup/rules_test.go`, then the runs above.
- The five scripts: `adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests`, `setup-check`: exit 0 each, and `git diff --check`: exit 0 (19:16Z).
