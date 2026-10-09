# The test runs of T-cht1

## Red 1: the tests before the functions (2026-10-09T13:36Z)

`go vet ./internal/route/` does not compile:

```
# github.com/pharzam/layup/internal/route
# [github.com/pharzam/layup/internal/route]
vet: internal/route/admit_test.go:44:13: undefined: Probed
```

## Red 2: a mutation of each rule, in full (the fix of round 1)

On a copy of `admit.go` (put back after), each of the fourteen rules broken alone, each line saying what was broken; each makes its own cases fail. Each line of output is cut at 110 characters. The mutation of the harness in `Probed` first passed, so the case "a later failed probe of another harness at the same version" was added before round 1; round 1 named four rules with no mutation (the first admitted pair, `ErrNoPair`, no row at the version, the harness of a models row), now here. `TestThePairOfASession` (note 3) shows the demo in one run.

```
== last: Probed: any passed row at the version, not the last
admit_test.go:47: a later failed probe at the same version: true, want false
admit_test.go:133: the pair: claude claude-opus-5-5 <nil>, want devin swe-2-high
== version: Probed: the rows of any version
admit_test.go:47: a passed probe at another version only: true, want false
admit_test.go:47: a later probe at another version is not this one: false, want true
admit_test.go:67: a version with no passed probe: true, want false
== harness: Probed: the rows of any harness
admit_test.go:47: a later failed probe of another harness at the same version: false, want true
== noRow: Probed: no row at the version is true
admit_test.go:47: a passed probe at another version only: true, want false
admit_test.go:47: a harness with no row: true, want false
admit_test.go:47: no row at all: true, want false
admit_test.go:67: a version with no passed probe: true, want false
== probed: Admitted: no probe needed
admit_test.go:67: a version with no passed probe: true, want false
admit_test.go:133: the pair: claude claude-opus-5-5 <nil>, want devin swe-2-high
== use: Admitted: use no admitted
admit_test.go:67: a model of use no: true, want false
== useyes: Admitted: a model with no row admitted
admit_test.go:67: a model with no row of models.tsv: true, want false
admit_test.go:67: a model of another harness: true, want false
== model: Admitted: the row of any model of the harness
admit_test.go:67: a model of use no: true, want false
admit_test.go:67: a model with no row of models.tsv: true, want false
admit_test.go:67: a model of another harness: true, want false
== modelharness: Admitted: the row of the model under any harness
admit_test.go:67: a model of another harness: true, want false
== order: Pair: the file's order, not position
admit_test.go:101: the first by position: devin swe-2-high <nil>, want claude claude-opus-5-5
== role: Pair: the pairs of any role
admit_test.go:101: the role's own list: claude claude-opus-5-5 <nil>, want codex gpt-6
admit_test.go:101: the second when the first is not admitted: codex gpt-6 <nil>, want devin swe-2-high
admit_test.go:113: a role of no list: <nil>, want ErrNoPair
admit_test.go:113: a tier of no list: <nil>, want ErrNoPair
== tier: Pair: the pairs of any tier
admit_test.go:101: the tier's own list: claude claude-opus-5-5 <nil>, want claude claude-fable-5-1
admit_test.go:113: a tier of no list: <nil>, want ErrNoPair
== first: Pair: the first pair, admitted or not
admit_test.go:101: the second when the first is not admitted: claude claude-opus-5-5 <nil>, want devin swe-2-h
admit_test.go:101: the third when the first two are not: claude claude-opus-5-5 <nil>, want codex gpt-6
admit_test.go:113: no admitted pair: <nil>, want ErrNoPair
admit_test.go:133: the pair: claude claude-opus-5-5 <nil>, want devin swe-2-high
== nopair: Pair: no admitted pair is no error
admit_test.go:113: no admitted pair: <nil>, want ErrNoPair
admit_test.go:113: a role of no list: <nil>, want ErrNoPair
admit_test.go:113: a tier of no list: <nil>, want ErrNoPair
```

## Green (2026-10-09T13:38Z)

Each with exit 0 on the tree of the commit `feat: T-cht1 …`: `go build ./...`, `go vet ./...`, `gofmt -l internal` (empty), `go test ./...`, `go test -tags=integration ./...` (with `TestPackageRules`); `adr-lint`, `prd-lint`, `link-lint`, `setup-check`, `run-discipline-tests`, `git diff --check`.

## Green of the fix of round 1, and the notes of round 2

On `9c5d344` (the fix of round 1, with `TestThePairOfASession`), each with exit 0: `go build ./...`, `go vet ./...`, `go test ./...`, `go test -tags=integration ./...`; `prd-lint`, `link-lint`, `setup-check`, `git diff --check` (round 2 ran them too, each exit 0). The close-out applies the three notes of round 2.
