# T-fn5s — the runs of the tests

Times from the clock of the host, UTC, 2026-10-15.

## Step T: the tests first, red

- 17:26:10, `go test ./internal/records/` with the new valid rows and the 32 new refusal cases of `session_test.go` and the old code: exit 1, 65 lines of failure; each is "`handoff` is not a value of enum(attempt|...)" or the like: the enums of the Go schemas had none of the six new kinds.
- 17:26:18, after the two enums only (`EventsSchema`, `ResultSchema`), before any rule: exit 1, 30 failures, each for the right reason: a valid new row is refused by the old rules of `CheckEvent` ("session is the empty value exactly for attempt, closed, rebased and a refused whose reason is pair"; "sha256 is given exactly for an artifact") and each refusal case of a new rule reads as accepted or is refused in another column.

## Green

- After `CheckEvent`, `CheckResultRow` and `CheckResult`: `go test ./internal/records/` ok. The valid rows now include a question with an unknown label and a `needs_human` word `maybe`, a lesson with records, a handoff with a session and with `—`, a valid and an invalid handoff, a comment suffix, and a merged event.

## The nine commands of acceptance criterion 3, before the freeze

`go build ./...`, `go vet ./...`, `go test ./...`: exit 0 each. `go test -count=1 -tags=integration ./...` (TMPDIR on the disk): exit 0, 17:31:28 to 17:50:43Z (19 min 15 s, the host loaded by two other tasks); the block test of `internal/tsv` and the comparison of the two blocks with their Go values (`internal/records`) among them. `adr-lint`, `prd-lint`, `link-lint`, `setup-check` and `run-discipline-tests`: exit 0 each. `git diff --check`: exit 0.
