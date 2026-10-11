# T-hpid — test runs

The command `layup spec check --facts` (#214, row 41b).

## Red (06:44Z, 2026-10-11)

The tests before the row of the command: each case is red with `unknown
command "spec"`, as each exit-2 case asserts its own reason (condition 1 of the
plan review). The `pre-commit` hook runs the unit level, so the red is recorded
here and the tests land with the code. `go test -run
'TestSpecCheck|TestTheUsageOfLayupListsItsCommands' ./internal/cli/` and `go
test -count=1 -tags=e2e -run TestSpecCheck ./cmd/layup/` (the lines cut at 250
characters; the usage elided):

```
--- FAIL: TestTheUsageOfLayupListsItsCommands (0.00s)
    cli_test.go:78: the usage has no "\n  spec check --facts SPANS --psb FILE ":
--- FAIL: TestSpecCheckGivesTwoWhenItCannotWriteTheTable (0.00s)
    spec_test.go:24: "start\tend\tfact\tclass\treason\ttext\n0\t2\t1\tneed\t—\tAb\n": exit 2, stderr "layup: unknown command \"spec\"\n\nusage: layup <command> [<argument>...]\n\ncommands:\n  version                                                 
    spec_test.go:24: "start\tend\tfact\tclass\treason\ttext\n": exit 2, stderr "layup: unknown command \"spec\"\n\nusage: layup <command> [<argument>...]\n\ncommands:\n  version                                                                         
--- FAIL: TestSpecCheckPassesAndFails (0.03s)
    spec_e2e_test.go:41: a span table that covers the file: exit 2, stdout "", stderr "layup: unknown command \"spec\"\n\nusage: layup <command> [<argument>...]\n\ncommands:\n  version                                                                  
    spec_e2e_test.go:46: a byte in two spans: exit 2, stdout "", stderr "layup: unknown command \"spec\"\n\nusage: layup <command> [<argument>...]\n\ncommands:\n  version                                                                                
--- FAIL: TestSpecCheckInputErrors (0.03s)
    spec_e2e_test.go:76: a missing problem statement: exit 2, stdout "", stderr "layup: unknown command \"spec\"\n\nusage: layup <command> [<argument>...]\n\ncommands:\n  version                                                                        
    spec_e2e_test.go:76: a missing span table: exit 2, stdout "", stderr "layup: unknown command \"spec\"\n\nusage: layup <command> [<argument>...]\n\ncommands:\n  version                                                                               
    spec_e2e_test.go:76: a problem statement not valid UTF-8: exit 2, stdout "", stderr "layup: unknown command \"spec\"\n\nusage: layup <command> [<argument>...]\n\ncommands:\n  version                                                                
    spec_e2e_test.go:76: a span table not valid UTF-8: exit 2, stdout "", stderr "layup: unknown command \"spec\"\n\nusage: layup <command> [<argument>...]\n\ncommands:\n  version                                                                       
    spec_e2e_test.go:76: a span table that its block refuses: exit 2, stdout "", stderr "layup: unknown command \"spec\"\n\nusage: layup <command> [<argument>...]\n\ncommands:\n  version                                                                
--- FAIL: TestSpecCheckUsageErrors (0.04s)
    spec_e2e_test.go:96: layup ["spec" "check"]: exit 2, stdout "", stderr "layup: unknown command \"spec\"\n\nusage: layup <command> [<argument>...]\n\ncommands:\n  version                                                                             
    spec_e2e_test.go:96: layup ["spec" "check" "--facts" "/home/layup/.cache/gotmp/TestSpecCheckUsageErrors2044768212/001/spans.tsv"]: exit 2, stdout "", stderr "layup: unknown command \"spec\"\n\nusage: layup <command> [<argument>...]\n\ncommands:\n
    spec_e2e_test.go:96: layup ["spec" "check" "--facts" "/home/layup/.cache/gotmp/TestSpecCheckUsageErrors2044768212/001/spans.tsv" "--facts" "/home/layup/.cache/gotmp/TestSpecCheckUsageErrors2044768212/001/spans.tsv" "--psb" "/home/layup/.cache/got
    spec_e2e_test.go:96: layup ["spec" "check" "--facts=" "--psb" "/home/layup/.cache/gotmp/TestSpecCheckUsageErrors2044768212/001/psb.md"]: exit 2, stdout "", stderr "layup: unknown command \"spec\"\n\nusage: layup <command> [<argument>...]\n\ncomma
    spec_e2e_test.go:96: layup ["spec" "check" "--facts" "/home/layup/.cache/gotmp/TestSpecCheckUsageErrors2044768212/001/spans.tsv" "--psb" "/home/layup/.cache/gotmp/TestSpecCheckUsageErrors2044768212/001/psb.md" "extra"]: exit 2, stdout "", stderr 
```

## Mutations (2026-10-11, 06:49Z to 06:50Z)

Each rule of `specCheck` mutated alone, the unit and e2e tests run: 7 of 7
caught, each mutant built.

```
== exit-always-pass: caught spec_e2e_test.go:46: a byte in two spans: exit 0, stdout "rule\tresult\tdetail\ncover\tfail\tbyte 3 is in 2 spans\nbounds\tfail\tthe span at 3: out of order, or over the 
== utf8-cut: caught spec_e2e_test.go:76: a problem statement not valid UTF-8: exit 1, stdout "rule\tresult\tdetail\ncover\tpass\t—\nbounds\tpass\t—\ntext\tfail\tthe span at 3: its text diffe
== block-refusal-cut: caught spec_e2e_test.go:76: a span table that its block refuses: exit 1, stdout "rule\tresult\tdetail\ncover\tfail\tbyte 0 is in no span\nbounds\tpass\t—\ntext\tpass\t—\nids\tpa
== write-error-cut: caught spec_test.go:24: "start\tend\tfact\tclass\treason\ttext\n0\t2\t1\tneed\t—\tAb\n": exit 0, stderr ""; want 2 and the error of the output
== refusal-names-no-file: caught spec_e2e_test.go:76: a span table that its block refuses: exit 2, stdout "", stderr "layup: line 1, column \"fact\": the header row has no field for this column\n"; want 
== psb-read-cut: caught spec_e2e_test.go:76: a missing problem statement: exit 1, stdout "rule\tresult\tdetail\ncover\tpass\t—\nbounds\tfail\tthe span at 0: end 2 past the file of 0 bytes\ntext\
== facts-flag-is-psb: caught spec_test.go:24: "start\tend\tfact\tclass\treason\ttext\n0\t2\t1\tneed\t—\tAb\n": exit 2, stderr "layup: spans.tsv: line 1, column \"start\": the header row has \"Ab\" in
```

## The package rules

The first full run after the code (06:46Z to 06:49Z) failed at the integration
level: `TestPackageRules` gave "May import: internal/cli imports internal/tsv,
which its row does not allow". The writer of the table moved to
`internal/spec` (`spec.WriteTSV`, as `psb.WriteTSV` is `internal/psb`'s),
and the mutations above ran again on that code.

## Green


At the code with `spec.WriteTSV` (06:50Z to 06:51Z): `go build ./...`, `go vet
./...`, `go test ./...`, `go test -tags=integration ./...` and `go test
-count=1 -tags=e2e ./...` pass, every package `ok`.
