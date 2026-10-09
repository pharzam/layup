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

Each with exit 0 on the tree of the commit `feat: T-ysph …`: `go build ./...`, `go vet ./...`, `gofmt -l internal cmd` (empty), `go test ./...`, `go test -tags=integration ./...`; `adr-lint`, `prd-lint`, `link-lint`, `setup-check`, `run-discipline-tests`, `git diff --check`.

## The fix of round 1 (2026-10-09T12:54Z)

Finding 1: `ReadRoutingRegister` refuses a position below 1 at its line; the two cases of the record (a `0` beside a `2`, and a lone `0`) are in the test. Finding 2: a mutation of each rule of `internal/route` alone, on a copy of `register.go` and `models.go` (put back after; four mutations were written again with `&& false`, as their first form left a variable unused and did not compile): each makes its own cases fail at their column. The empty-value cases run on `plainRow`, which has no cap, so the rule of `{cap}` cannot refuse an empty command for them.

```
== the columns that never hold the empty value
register_test.go:87: the empty value in wall: read, want an error in column wall
register_test.go:87: the empty value in command: read, want an error in column command
register_test.go:87: the empty value in prompt: read, want an error in column prompt
register_test.go:87: the empty value in version: read, want an error in column version
== wall 0
register_test.go:118: a wall of 0: read, want an error in column wall
== {cap} exactly when cap
register_test.go:118: a cap and no {cap}: read, want an error in column command
register_test.go:118: {cap} and no cap: read, want an error in column command
== one space between words
register_test.go:118: two spaces in command: read, want an error in column command
== the version's words
register_test.go:118: a space at the end of version: read, want an error in column version
== an absolute credential
register_test.go:118: a credential that is not absolute: read, want an error in column credential
== credential_to exactly when credential
register_test.go:118: a credential and no credential_to: read, want an error in column credential_to
register_test.go:118: credential_to and no credential: read, want an error in column credential_to
== the form and list of var:NAME
register_test.go:118: var: with a name of another form: read, want an error in column credential_to
register_test.go:118: var: with a name of the named list: read, want an error in column credential_to
register_test.go:118: var: with a forge credential's name: read, want an error in column credential_to
== the form of credential_to
register_test.go:118: credential_to of another form: read, want an error in column credential_to
== an absolute policy
register_test.go:118: a policy that is not absolute: read, want an error in column policy
== the empty value of models
models_test.go:33: the empty value in context: read, want an error in column context
models_test.go:33: the empty value in date: read, want an error in column date
models_test.go:33: the empty value in use: line 2, column "reason": reason is the empty value exactly when use
== reason exactly when use is no
models_test.go:42: a model in use with a reason: read, want an error in column reason
models_test.go:42: a model not to use with no reason: read, want an error in column reason
== a position 1 or more
models_test.go:70: a position 0 beside a position 2: read, want an error in column position
== the positions 1 to k
models_test.go:68: a gap in the positions: read, want an error in column position
== the harness of a models row
models_test.go:95: a models row of a harness that the register lacks: <nil>; want an error of models.tsv
== the harness of a routing row
models_test.go:98: a routing row of a harness that the register lacks: routing.tsv: line 4, column "model": gp
== the model of a routing row
models_test.go:95: a routing row whose model has no row of its harness: <nil>; want an error of routing.tsv
== the mode of a credential
models_test.go:128: mode 0644: <nil>; want an error that names the file
models_test.go:128: mode 0400: <nil>; want an error that names the file
models_test.go:128: a directory: <nil>; want an error that names the file
== the owner of a credential
models_test.go:128: another user's file: <nil>; want an error that names the file
== the PATH of file:
register_test.go:118: file: with an absolute PATH: read, want an error in column credential_to
register_test.go:118: file: with a part ..: read, want an error in column credential_to
register_test.go:118: file: with .. at its end: read, want an error in column credential_to
register_test.go:118: file: with an empty PATH: read, want an error in column credential_to
== NAME=VALUE of vars
register_test.go:118: a fixed variable that is not NAME=VALUE: read, want an error in column vars
register_test.go:118: a fixed variable whose name is of another form: read, want an error in column vars
== the refused names of vars
register_test.go:118: a fixed variable of the named list: read, want an error in column vars
register_test.go:118: a fixed variable that is the credential's: read, want an error in column vars
register_test.go:118: a fixed variable GH_TOKEN: read, want an error in column vars
register_test.go:118: a fixed variable GITHUB_TOKEN: read, want an error in column vars
== an http or https source
models_test.go:42: a source that is not a URL: read, want an error in column source
models_test.go:42: a source of another scheme: read, want an error in column source
```

Note 4: `go test -tags=e2e ./...` exits 0 on this tree (the e2e of `layup run` writes the four registers of `M2b`); the green record above no longer says that the integration run held it.
