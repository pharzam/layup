# The test runs of T-ysph

## Red 1: the tests before the readers (2026-10-09T12:30Z)

`go vet ./internal/route/` does not compile:

```
# github.com/pharzam/layup/internal/route
# [github.com/pharzam/layup/internal/route]
vet: internal/route/models_test.go:22:15: undefined: ReadModels
```

## Red 2: `layup run` before it reads the registers of `M2b` (2026-10-09T12:32Z)

With `internal/route` built and `hostInputs` not yet changed, `TestRunNewChecksEachInputBeforeTheFirstStep` fails for the four new cases, each a step that ran where exit 2 was wanted; the other cases pass:

```
    run_test.go:195: no models.tsv: exit 1, stdout "step\tresult\tdetail\n", a step ran true, stderr ""; want exit 2 and an input error that
    run_test.go:195: no routing.tsv: exit 1, stdout "step\tresult\tdetail\n", a step ran true, stderr ""; want exit 2 and an input error tha
    run_test.go:195: a routing row of a harness that the register lacks: exit 1, stdout "step\tresult\tdetail\n", a step ran true, stderr ""
    run_test.go:195: a credential of mode 0644: exit 1, stdout "step\tresult\tdetail\n", a step ran true, stderr ""; want exit 2 and an inpu
```

## Green (2026-10-09T12:34Z)

Each with exit 0 on the tree of the commit `feat: T-ysph …`: `go build ./...`, `go vet ./...`, `gofmt -l internal cmd` (empty), `go test ./...`, `go test -tags=integration ./...` (with the e2e of `layup run` on the registers of `M2b`); `adr-lint`, `prd-lint`, `link-lint`, `setup-check`, `run-discipline-tests`, `git diff --check`.
