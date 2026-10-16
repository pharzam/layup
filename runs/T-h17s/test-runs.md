# T-h17s — the runs of the tests

UTC, 2026-10-15; times from the clock of the host.

## Step T, red (18:34 to 18:35Z)

With a stub of the package, so that each test fails on its assertion: all but one of the unit tests failed (`TestEachDefaultPairNeedsItsOwnItemsAndKinds` passed on an empty `Default()`; a guard `len(Default()) == 13` was added); in the integration tag `TestTheSchemaEqualsItsBlock`, `TestTheDefaultRowsAreTheTableOfTheSpecification` and `TestTheListsAreThoseOfTheSpecification` failed. With a table of `M2e` in `goodTable`, `TestReadTable` failed (the row of `internal/handoff` was not read) and `TestReadTableRefusesWhatItCannotRead` failed on "no table of M2e", "a row of M2e for a package of another table" and "a row of M2e that stands for a row of phase 1". These runs were made while the plan review ran (18:33 to 18:40Z); the tests of its conditions came after it.

## The code before its tests, and the mutants

Conditions 1 to 3 of the plan review, and finding 1 and notes 2 and 3 of round 1, were code first and tests second. The red evidence is a run against the earlier code (18:53Z: `TestAReplacementRowCannotLetABadResultThrough` failed on six cases) and 30 mutants of the final `check.go` and `table.go`, each applied alone and undone: all caught; after round 2, seven more (the range of each hexadecimal digit, the first comment named, the two checks of `WriteTable`) are caught by three added tests. The script is a one-off, not kept.

## Green

- 18:53:19 to 19:03:53Z: `go build`, `go vet`, `go test -count=1 ./...` exit 0; `go test -count=1 -tags=integration ./...`: 20 packages ok, `internal/run` timed out at the default 10 minutes (other tasks run on the host); alone with `-timeout 1500s` (19:04:15 to 19:15:03Z): ok, 644.763 s. A run at 18:41Z had failed on two tests of `cmd/layup` (my fixture added two packages with no row); fixed.
- After the fixes of round 1 and the notes of round 2: `go build`, `go vet`, `go test ./...` exit 0; `go test -tags=integration` of `internal/handoff`, `internal/tsv`, `cmd/layup`, `internal/records` ok; nothing outside `internal/handoff` and its fixtures changed since the full run.
- `adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests`, `setup-check`: exit 0 each; `git diff --check`: exit 0.

## Round 2's notes, left as they are

Note 3 (a replacement table that lists `schema` for the rendered task makes it impossible: a question for row 66b, posted on #358); note 5 (`ReadTable` returns a plain error that names the line and the column, as `records.md` asks, not `*tsv.Error`); note 6 (the order, above); note 7 (the budget).
