# The test runs of T-4c3q

## Red 1: the tests before the package (2026-10-09T13:04Z)

`go vet ./internal/ledger/` does not compile:

```
# github.com/pharzam/layup/internal/ledger
# [github.com/pharzam/layup/internal/ledger]
vet: internal/ledger/ledger_test.go:25:16: undefined: Session
```

## Red 2: a mutation of each rule

On a copy of `ledger.go` (put back after), each rule broken alone. One mutation is equivalent: "unknown when the tokens are not observed" alone leaves the nil checks of the counts, and `CheckTelemetry` makes `observed` exactly the three counts, so no input tells them apart. The mutation of the session's harness first passed: no case had a row of the model under another harness on the newest date, so the case of condition 1 was added, and it fails then. Three mutations were written again (their first form did not compile, or the shell broke their text):

```
== reported only for api
ledger_test.go:120: the telemetry row of S-1a2b3c4d: column money_status: a subscription session is computed from the li
== unknown for two models
ledger_test.go:193: a report that names two models: "computed" "11.5" "USD" "P-004", want unknown and no money, currency
== unknown for another model
ledger_test.go:193: a report that names one model, not the session's: "computed" "11.5" "USD" "P-004", want unknown and
== the newest date
ledger_test.go:122: a subscription with a cost is computed: "34.5" "USD" "computed" "P-001", want "11.5" "USD" "computed
ledger_test.go:122: computed from the newest rows: "34.5" "USD" "computed" "P-001", want "11.5" "USD" "computed" "P-004"
== two rows of a class
ledger_test.go:193: a class with two rows of the newest date: "computed" "12.5" "USD" "P-004", want unknown and no money
== one currency
ledger_test.go:193: rows of two currencies: "computed" "11.5" "USD" "P-004", want unknown and no money, currency or pric
== the price of in
ledger_test.go:122: a subscription with a cost is computed: "11.5" "USD" "computed" "P-005", want "11.5" "USD" "computed
ledger_test.go:122: computed from the newest rows: "11.5" "USD" "computed" "P-005", want "11.5" "USD" "computed" "P-004"
== per million
ledger_test.go:122: computed from the newest rows: "11500.0" "USD" "computed" "P-004", want "11.5" "USD" "computed" "P-0
ledger_test.go:122: a subscription with a cost is computed: "11500.0" "USD" "computed" "P-004", want "11.5" "USD" "compu
== one decimal digit
ledger_test.go:120: line 2, column "money": "6" is not a value of decimal
== harness
ledger_test.go:129: a row of another harness on that date: "" "" "unknown" "", want "11.5" "USD" "computed" "P-004"
== cost
panic: runtime error: invalid memory address or nil pointer dereference [recovered, repanicked]
panic({0x595aa0?, 0x72ed20?})
== check
ledger_test.go:209: an end before the start: a row, want an error
```

## Green (2026-10-09T13:06Z)

Each with exit 0 on the tree of the commit `feat: T-4c3q …`: `go build ./...`, `go vet ./...`, `gofmt -l internal` (empty), `go test ./...`, `go test -tags=integration ./...` (with `TestPackageRules` and the new package); `adr-lint`, `prd-lint`, `link-lint`, `setup-check`, `run-discipline-tests`, `git diff --check`. The case of `3 × 0.1` per million writes `0.0000003`, where a float64 gives `3.0000000000000004e-07`.

## The fix of round 1 (2026-10-09T13:35Z)

Finding 1: the copied line of the mutation record kept a trailing space; it is gone, and `git diff --check ad8774d` exits 0 on this tree. Finding 2: the §6 row of `REQ-011` is back to its text at `ad8774d`; the tests are in its §12 Test cell only, now with the fifth test (note 3). Note 5: `Row` refuses a price row of another width than the block's, with a case. Each check of the green run passes again on the fix tree.
