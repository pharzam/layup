# The test runs of T-m1dx

## Red 1: the tests before the package (2026-10-09T12:58Z)

`go vet ./internal/rules/` does not compile:

```
# github.com/pharzam/layup/internal/rules
# [github.com/pharzam/layup/internal/rules]
vet: internal/rules/rules_test.go:20:40: undefined: Entry
```

## Red 2: a mutation of each rule

On a copy of `rules.go` (put back after), each rule broken alone; each makes its own cases fail (two mutations were written again with `&& false` or `|| true`, as their first form left a variable unused and did not compile):

```
== the one exception
rules_test.go:42: an exception of another path: read, want an error
rules_test.go:42: another exception of guardrails: read, want an error
== a directory pattern
rules_test.go:72: .github/workflows/ci.yml matches {Pattern: Exception:} (false), want .github/
rules_test.go:72: docs/ci/review-record-lint.sh matches {Pattern:*.sh Exception:} (true), want docs/ci/
rules_test.go:134: a rule path: "" "", want "rule-path" "docs/tests/test-levels.md"
== an exact pattern
rules_test.go:72: docs/guardrails.md matches {Pattern: Exception:} (false), want docs/guardrails.md
rules_test.go:72: AGENTS.md matches {Pattern: Exception:} (false), want AGENTS.md
rules_test.go:72: go.mod matches {Pattern: Exception:} (false), want go.mod
== the kind .sh
rules_test.go:70: internal/tool/run.sh: no match, want the kind .sh
rules_test.go:70: build.sh: no match, want the kind .sh
rules_test.go:134: a script anywhere: "" "", want "rule-path" "tools/gen.sh"
== no ## 2.
rules_test.go:111: a head with no ## 2.: true, want false
== the next ## line
rules_test.go:107: a line added in §3: true, want false
rules_test.go:107: an added line at the line of ## 3.: true, want false
rules_test.go:134: an addition in §3 of the guardrails: "" "", want "rule-path" "docs/guardrails.md"
== a removed line
rules_test.go:107: a changed line: true, want false
== at least one addition
rules_test.go:107: an empty diff: true, want false
== the workflow first
rules_test.go:134: a workflow before a rule path: "rule-path" "AGENTS.md", want "workflow" ".github/workflows/
== inside section 2
rules_test.go:107: a line added in §1: true, want false
rules_test.go:107: a line added in §3: true, want false
rules_test.go:107: an added line at the line of ## 3.: true, want false
== the exception in Check
rules_test.go:141: the guardrails with no exception in the register: "", want rule-path
```

## Green (2026-10-09T13:00Z)

Each with exit 0 on the tree of the commit `feat: T-m1dx …`: `go build ./...`, `go vet ./...`, `gofmt -l internal` (empty), `go test ./...`, `go test -tags=integration ./...` (with `TestPackageRules` and the new package); `adr-lint`, `prd-lint`, `link-lint`, `setup-check`, `run-discipline-tests`, `git diff --check`.

## The notes of round 1 (2026-10-09T13:12Z)

Note 4: the case "a removed line" had no added line, so the rule of an addition refused it whichever rule stood; it is now a removed line beside an addition, and with the rule of a removed line broken alone (a copy, put back after) it fails: `rules_test.go:108: a removed line beside an addition: true, want false`, with "a changed line". Note 3: a case of a mode change with an addition in §2, refused. Note 5: `section2` ends one past the last line of a head that ends with a line feed. `go test ./internal/rules/` passes.
