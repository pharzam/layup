# The test runs of T-cht1

## Red 1: the tests before the functions (2026-10-09T13:36Z)

`go vet ./internal/route/` does not compile:

```
# github.com/pharzam/layup/internal/route
# [github.com/pharzam/layup/internal/route]
vet: internal/route/admit_test.go:44:13: undefined: Probed
```

## Red 2: a mutation of each rule, in full

On a copy of `admit.go` (put back after), each rule broken alone. The mutation of the harness in `Probed` first passed: no case had another harness at the same version string, so the case "a later failed probe of another harness at the same version" was added, and it fails then (the last group).

```
== last
admit_test.go:45: a later failed probe at the same version: true, want false
== version
admit_test.go:45: a passed probe at another version only: true, want false
admit_test.go:45: a later probe at another version is not this one: false, want true
admit_test.go:65: a version with no passed probe: true, want false
== harness
== probed
admit_test.go:65: a version with no passed probe: true, want false
== use
admit_test.go:65: a model of use no: true, want false
== useyes
admit_test.go:65: a model with no row of models.tsv: true, want false
admit_test.go:65: a model of another harness: true, want false
== model
admit_test.go:65: a model of use no: true, want false
admit_test.go:65: a model with no row of models.tsv: true, want false
admit_test.go:65: a model of another harness: true, want false
== order
admit_test.go:99: the first by position: devin swe-2-high <nil>, want claude claude-opus-5-5
== role
admit_test.go:99: the role's own list: claude claude-opus-5-5 <nil>, want codex gpt-6
admit_test.go:99: the second when the first is not admitted: codex gpt-6 <nil>, want devin swe-2-high
admit_test.go:111: a tier of no list: <nil>, want ErrNoPair
admit_test.go:111: a role of no list: <nil>, want ErrNoPair
== tier
admit_test.go:99: the tier's own list: claude claude-opus-5-5 <nil>, want claude claude-fable-5-1
admit_test.go:111: a tier of no list: <nil>, want ErrNoPair
== harness (with the case of another harness at the same version)
admit_test.go:47: a later failed probe of another harness at the same version: false, want true
```

## Green (2026-10-09T13:38Z)

Each with exit 0 on the tree of the commit `feat: T-cht1 …`: `go build ./...`, `go vet ./...`, `gofmt -l internal` (empty), `go test ./...`, `go test -tags=integration ./...` (with `TestPackageRules`); `adr-lint`, `prd-lint`, `link-lint`, `setup-check`, `run-discipline-tests`, `git diff --check`.
