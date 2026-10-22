# T-c3zv: test runs

Task T-c3zv, issue #473: the Shape options prompts name the empty evidence of a fail as U+2014, not "a dash". Times are from `date -u` after each run. Go 1.26.10 from mise.

## 1. Red (R8), against the current prompts

Command: `go test -count=1 -run '^TestEachOptionsPromptNamesTheEmptyEvidenceOfAFail$' ./internal/run/`, 05:36:10Z to 05:36:12Z, exit 1.

First failing assertion: `the options prompt names the empty evidence in words (a dash)`; then `the options prompt does not say that a hyphen or other text is refused`; the rewrite prompt fails the same way. The control (`-` in a fail's evidence is refused by `spec.ReadOptions`) passes today, as it must.

The red could not be committed alone: the pre-commit hook runs the unit level and refused it. The test is committed together with the green change (no `t.Skip`, no `--no-verify`); the run above is the red.

## 2. Green

The two sentences of `internal/run/options.go` (`optionsPromptText`, `rewritePromptText`) now say `the single character "—" (U+2014), the empty value, and nothing else; a hyphen - or any other text is refused`; the pass clause no longer says "never a dash".

- `go test -count=1 -run '^TestEachOptionsPromptNamesTheEmptyEvidenceOfAFail$|TestThePromptsList' -v ./internal/run/`: exit 0, both pass.
- Unit level, `go test -count=1 ./...`, 05:37:06Z, exit 0.
- `go build ./...`, `go vet ./...` (no tag, `integration`, `e2e`), `gofmt -l .`: clean.
- `go test -count=1 -tags=integration ./internal/run/ -run 'Options|Shape'`, 05:37:24Z, exit 0 (8.3 s).
- The six lints and `sh docs/tests/freeze-check/run-all.sh`: exit 0 each.
