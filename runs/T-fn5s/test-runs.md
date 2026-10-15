# T-fn5s — the runs of the tests

Times from the clock of the host, UTC, 2026-10-15.

## Step T: the tests first, red

- 17:26:10, `go test ./internal/records/` with the new valid rows and the 32 new refusal cases of `session_test.go` and the old code: exit 1, 65 lines of failure; each is "`handoff` is not a value of enum(attempt|...)" or the like: the enums of the Go schemas had none of the six new kinds.
- 17:26:18, after the two enums only (`EventsSchema`, `ResultSchema`), before any rule: exit 1, 30 failures, each for the right reason: a valid new row is refused by the old rules of `CheckEvent` ("session is the empty value exactly for attempt, closed, rebased and a refused whose reason is pair"; "sha256 is given exactly for an artifact") and each refusal case of a new rule reads as accepted or is refused in another column.

## Green

- After `CheckEvent`, `CheckResultRow` and `CheckResult`: `go test ./internal/records/` ok. The valid rows now include a question with an unknown label and a `needs_human` word `maybe`, a lesson with records, a handoff with a session and with `—`, a valid and an invalid handoff, a comment suffix, and a merged event.

## The nine commands of acceptance criterion 3, before the freeze

`go build ./...`, `go vet ./...`, `go test ./...`: exit 0 each. `go test -count=1 -tags=integration ./...` (TMPDIR on the disk): exit 0, 17:31:28 to 17:50:43Z (19 min 15 s, the host loaded by two other tasks); the block test of `internal/tsv` and the comparison of the two blocks with their Go values (`internal/records`) among them. `adr-lint`, `prd-lint`, `link-lint`, `setup-check` and `run-discipline-tests`: exit 0 each. `git diff --check`: exit 0.

## After round 1, 2026-10-15T18:10Z

Round 1 found three rules with no failing case (the numbers of the decisions and the lessons; `result/..` and `result/.`; comment ID 0), the item of a handoff checked by form and not against the table (`pair` added for a pair with no row), and "one space between" not checked. Each is fixed with its case. **Mutants** on a copy of `session.go`, `go test ./internal/records/` exit 1 for each: the numbers of `decision` and `lesson` dropped from the loop; `isPath` dropped for a comment; the comment ID loosened to `[0-9]*`; the item list not read (`|| true`); the one-space check dropped for a question; the number form of `merged` loosened to allow a leading zero. After the fixes: the five checks exit 0; `go test ./internal/records/ ./internal/session/ ./internal/run/ ./internal/tsv/` ok; `go test -tags=integration ./internal/records/ ./internal/tsv/ ./internal/session/` ok (the full integration run before the freeze is above and its code is unchanged apart from `internal/records`); `git diff --check` exit 0; budget 232 lines over 7 files before the close-out.
