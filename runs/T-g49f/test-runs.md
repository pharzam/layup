# T-g49f — test runs

The specification of `M2c` (#206), step T of the plan with condition 1 of its
review: the section check and the block test, red before the files of D1 and
D2, green after.

## Red (2026-10-10, 17:16Z)

`sh runs/T-g49f/sections.sh`, before `docs/spec/intake.md` existed: exit 1,
each of the 20 headings `FAIL`, the first lines:

```
FAIL docs/spec/intake.md: ### The sessions of Intake (Intake, what M2c adds to layup run: the sessions of a phase, their roles and routing (decision 5))
FAIL docs/spec/intake.md: ### The budget of Intake (the cap of milestone 0 before each Intake session (K38, decision 3))
FAIL docs/spec/intake.md: ### The rule gaps (the rule gaps of layup psb check, committed (later-p2-intake-spec; F-21, decision 4))
```

`go test -tags=integration -run TestEverySchemaBlockIsBuiltOrNotYetBuilt ./internal/tsv/`,
with the nine new names in `notYetBuilt` and no block yet:

```
blocks_integration_test.go:36: answers is listed, but no block of docs/spec/ has that name
blocks_integration_test.go:36: answers-review is listed, but no block of docs/spec/ has that name
blocks_integration_test.go:36: budget is listed, but no block of docs/spec/ has that name
blocks_integration_test.go:36: coverage is listed, but no block of docs/spec/ has that name
blocks_integration_test.go:36: draft is listed, but no block of docs/spec/ has that name
blocks_integration_test.go:36: meaning is listed, but no block of docs/spec/ has that name
blocks_integration_test.go:36: questions is listed, but no block of docs/spec/ has that name
blocks_integration_test.go:36: spans is listed, but no block of docs/spec/ has that name
blocks_integration_test.go:36: spec-check is listed, but no block of docs/spec/ has that name
```

## Green (2026-10-10, 17:22Z)

`sh runs/T-g49f/sections.sh`: exit 0, each of the 20 headings `ok` (the row of
the table of `M2c` named its heading `##` at first and was corrected to the
heading's real level, `###`). The block test passes with the nine blocks. Each
with exit 0: `go build ./...`, `go vet ./...`, `go test ./...`,
`go test -tags=integration ./...` (with `TestEverySchemaBlockIsBuiltOrNotYetBuilt`
and `TestPackageRules`, which reads the changed cells of `internal/cli` and
`internal/run`), `adr-lint`, `prd-lint`, `link-lint`, `setup-check`,
`run-discipline-tests`, `git diff --check`.
