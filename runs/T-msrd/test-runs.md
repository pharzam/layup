# T-msrd — test runs

The specification of `M2d` (#265), step T of the plan with condition 1 of the
plan review of #206 applied from the start: the section check and the block
test, red before the files of D1 and D2, green after.

## Red (2026-10-11, 14:48Z)

`sh runs/T-msrd/sections.sh`, before `docs/spec/scaffold.md` existed and before
the records and package headings of `M2d`: exit 1, each of the 31 headings
`FAIL`, the first lines:

```
FAIL docs/spec/scaffold.md: ## REQ-002 — The setup questions in the one batch (the setup questions of S01 and S10 in the one batch (deferred by decision 9 of #206; decision 5; F-21))
FAIL docs/spec/scaffold.md: ### The step markers (the proposed marker sources with their byte-exact quote and the code check (decision 5))
FAIL docs/spec/scaffold.md: ### The answers of the setup (the answers of the setup written from the copied Intake answers (later-p2-scaffold; decision 5))
```

`go test -tags=integration -run TestEverySchemaBlockIsBuiltOrNotYetBuilt ./internal/tsv/`,
with the three new names in `notYetBuilt` and no block yet:

```
blocks_integration_test.go:36: forge-check is listed, but no block of docs/spec/ has that name
blocks_integration_test.go:36: markers is listed, but no block of docs/spec/ has that name
blocks_integration_test.go:36: setup-tree is listed, but no block of docs/spec/ has that name
```

## Green (2026-10-11, 14:58Z)

`sh runs/T-msrd/sections.sh`: exit 0, each of the 31 headings `ok` (the row
of the setup under `layup run` named its heading without the backquotes of
the file at first, and was corrected to the heading as written). The block test
passes with the three blocks. Each with exit 0: `go build ./...`,
`go vet ./...`, `go test -tags=integration ./internal/tsv/ ./internal/records/
./internal/run/ ./internal/setup/`, `go test -tags=integration -run
TestPackageRules ./cmd/layup/` (which reads the changed cell of
`internal/run`), `adr-lint`, `prd-lint`, `link-lint`,
`run-discipline-tests` and `setup-check` (whose check `adapted` first
found the word "optional" in `scaffold.md`, reworded).
