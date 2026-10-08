# The test runs of T-esfe

## Red 1: the new tests on the checker of `main` (2026-10-08T06:15Z)

```
# github.com/pharzam/layup/cmd/layup
# [github.com/pharzam/layup/cmd/layup]
vet: cmd/layup/rules_test.go:77:64: unknown field connects in struct literal of type row
```

## Red 2: a stub that compiles (`connects`, `readAdapter` that reads nothing, the new signature of `checkRules`)

```
--- FAIL: TestReadTable (0.00s)
    rules_test.go:86: rows map[cmd/layup:{imports:[internal/cli] program: connects:[]} internal/cli:{imports:[internal/psb internal/gate] program: connects:[]} internal/gate:{imports:[internal/git] program:sh connects:[]} internal/git:{imports:[] program:git connects:[]} internal/psb:{imports:[] program: connects:[]} internal/setup:{imports:[internal/git] program: connects:[]}], <nil>
        want map[cmd/layup:{imports:[internal/cli] program: connects:[net/http]} internal/cli:{imports:[internal/psb internal/gate] program: connects:[net/http]} internal/forge/github:{imports:[] program: connects:[net net/http crypto/tls]} internal/gate:{imports:[internal/git] program:sh connects:[]} internal/git:{imports:[] program:git connects:[]} internal/psb:{imports:[] program: connects:[]} internal/setup:{imports:[internal/git] program: connects:[]}]
--- FAIL: TestReadTableRefusesWhatItCannotRead (0.00s)
    rules_test.go:117: no column Connects: no error, 6 rows; want an error
    rules_test.go:117: a row of phase 1 that phase 1 lacks: no error, 6 rows; want an error
    rules_test.go:117: a row that mixes the two kinds of cell: no error, 6 rows; want an error
    rules_test.go:117: a Connects cell in prose: no error, 6 rows; want an error
    rules_test.go:117: a full row of M2a for a package of phase 1: no error, 6 rows; want an error
    rules_test.go:117: no table of M2a: no error, 6 rows; want an error
--- FAIL: TestReadAdapter (0.00s)
    rules_test.go:126: readAdapter = "", <nil>; want internal/forge/github
--- FAIL: TestEachRuleAndColumnFindsItsBreach (0.00s)
    --- FAIL: TestEachRuleAndColumnFindsItsBreach/rule_5,_a_dependency_with_no_Connects (0.00s)
        rules_test.go:242: findings
            rule 5: internal/cli depends on net/http
            rule 5: internal/psb depends on net/http
            want exactly
            rule 5: internal/psb depends on net/http
    --- FAIL: TestEachRuleAndColumnFindsItsBreach/rule_5,_a_dependency_that_Connects_names (0.00s)
        rules_test.go:242: findings
            rule 5: cmd/layup depends on net/http
            want exactly
    --- FAIL: TestEachRuleAndColumnFindsItsBreach/rule_5,_an_own_import_outside_the_adapter (0.00s)
        rules_test.go:242: findings
            rule 5: cmd/layup depends on net/http
            rule 5: internal/cli depends on net/http
            want exactly
            rule 5: internal/cli imports net/http; only internal/forge/github imports them
    --- FAIL: TestEachRuleAndColumnFindsItsBreach/rule_5,_the_adapter_imports_it (0.00s)
        rules_test.go:242: findings
            rule 5: internal/forge/github depends on crypto/tls
            rule 5: internal/forge/github depends on net
            rule 5: internal/forge/github depends on net/http
            table: internal/forge/github has no row
            want exactly
FAIL
FAIL	github.com/pharzam/layup/cmd/layup	0.134s
FAIL
```

The cases "rule 5, a dependency that Connects does not name" and the fixture are regression guards: they pass on the checker of `main` too.

## Green: the checker, then `packages.md` (2026-10-08T06:16Z)

With the real checker, the unit tests of the checker pass; `TestPackageRules` then failed with `0 lines start with "The one package that imports them:", want 1` until `packages.md` held the line of rule 5; with it, `go test -tags=integration -run TestPackageRules ./cmd/layup/` passes, and the fixture gives `rule 5: cmd/layup imports net/http; only internal/forge/github imports them` and `rule 5: internal/psb depends on net`.
