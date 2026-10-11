# T-4ywt — test runs

The span check and the package rules of `M2c` (#213, row 41a), in the order of
note 4 of the plan review: D1 red, the check; D2 red, the package rules; then
the documents.

## D1, red (04:06Z, 2026-10-11)

`internal/spec/spec_test.go` before `spec.go`, `go vet ./internal/spec/`:

```
vet: internal/spec/spec_test.go:42:11: undefined: Check
```

## D2, red (04:07Z)

The checker's fixture with a table of `M2c`, before the checker reads it,
`go test ./cmd/layup/ -run TestReadTable`:

```
rules_test.go:110: rows map[cmd/layup:{imports:[internal/cli] program: register:false connects:[net/http]} internal/cli:{imports:[internal/psb internal/gate] pr
rules_test.go:150: a row of M2c that stands for a row of phase 1: no error, 8 rows; want an error
rules_test.go:150: no table of M2c: no error, 8 rows; want an error
rules_test.go:150: a row of M2c for a package of another table: no error, 8 rows; want an error
```

The line of the engine checks with `internal/spec`, while the package was
moved aside (it does not exist), `go test -tags=integration -run
TestPackageRules ./cmd/layup/` at 04:07Z:

```
rules_integration_test.go:46: the line of the engine checks names internal/spec, which is no package of the module
```

With the package back, it passes.

## Mutations (2026-10-11, 04:08Z)

Each rule of `spec.go` mutated alone (its condition joined with `&& false`,
or its set narrowed), `go test ./internal/spec/`: 15 of 15 caught, none failed
to build. A first run caught 14: the replacement of the line feed in the rule
Text had no case; `TestASpanAcrossALineFeed` was added, and the second run
caught it. The schema's test is after the code; its red is a mutation (`detail`
made a key), caught: "column 3 (detail): key: the block has \"false\"".

```
== cover-none: caught spec_test.go:78: a byte in no span: cover fail "byte 30 is in 0 spans", want fail with "in no span"
== cover-count: caught spec_test.go:78: a byte in no span: cover pass "", want fail with "in no span"
== cover-blank: caught spec_test.go:78: a span of white space only: cover pass "", want fail with "white space only"
== cover-space-set: caught spec_test.go:47: cover: fail "byte 7 is in no span", want pass with no detail
== bounds-start-end: caught spec_test.go:78: a start not below the end: bounds pass "", want fail with "not below"
== bounds-size: caught /home/layup/projects/layup/.worktree/T-4ywt/internal/spec/spec_test.go:75 +0x15c7
== bounds-order: caught spec_test.go:78: two spans out of order: bounds pass "", want fail with "out of order"
== text-compare: caught spec_test.go:78: a text that differs from the file: text pass "", want fail with "differs"
== text-dash: caught spec_test.go:98: a span whose text is one em dash: fail "the span at 0: its text differs from the file's \"—\"", want pass
== text-norm: caught spec_test.go:119: a span across a tab and a line feed: fail "the span at 0: its text differs from the file's \"One fact across\\na line.\"", want pass
== ids-order: caught spec_test.go:78: a fact number that skips: ids pass "", want fail with "want 1"
== classes-class: caught spec_test.go:78: a fact with no class: classes pass "", want fail with "class"
== classes-fact-reason: caught spec_test.go:78: a fact with a reason: classes pass "", want fail with "reason"
== classes-nonfact-class: caught spec_test.go:78: a span that is not a fact, with a class: classes pass "", want fail with "class"
== classes-nonfact-reason: caught spec_test.go:78: a span that is not a fact, with no reason: classes pass "", want fail with "reason"
15 of 15
```

## Green (2026-10-11, 04:11Z)

Each with exit 0: `go build ./...`, `go vet ./...`, `gofmt -l internal cmd`
(empty), `go test ./...`, `go test -tags=integration ./...` (with
`TestPackageRules`, the block test with `spec-check` in `built`, and the
schemas of `internal/spec`), `go test -tags=e2e ./...`; `adr-lint`,
`prd-lint`, `link-lint`, `setup-check`, `run-discipline-tests`,
`git diff --check`.
