# T-qk66 — test runs

The step `rule-gaps` (#217, row 43a).

## Red (04:23Z, 2026-10-11)

The tests before the step. `go vet ./internal/run/`, before `restartNames`:

```
vet: internal/run/command_test.go:38:12: undefined: restartNames
```

With `restartNames` only (the names, not the step),
`go test -tags=integration -run TestTheStepRuleGaps ./internal/run/`:

```
run_integration_test.go:536: the restart: the steps [forge clone version lease probe phase], want [forge clone version lease probe rule-gaps]
run_integration_test.go:538: the progress ["1/6 forge" "2/6 clone" "3/6 version" "4/6 lease" "5/6 probe" "6/6 phase"], want ["1/6 forge" "2/6 clone" "3/6 version" "4/6 le
run_integration_test.go:544: git show refs/heads/layup-records:intake/psb.tsv: exit status 128
```

## Mutations (2026-10-11, 04:26Z to 04:29Z)

Each rule mutated alone, the integration tests of the restart run: 5 of 5
caught. A first run caught 4: the condition "the watch is written" of Start
done had no test; `TestARestartWhoseWatchIsNotWrittenRunsPhase` was added,
and the second run caught it.

```
== done-before: caught run_integration_test.go:554: the second restart's last row {Name:rule-gaps Result:fail Detail:git commit -m layup run: the rule gaps of the problem statement: exit status
== start-done-watch: caught run_integration_test.go:569: the restart: the steps [forge clone version lease probe rule-gaps], want [forge clone version lease probe phase]
== release-after: caught run_integration_test.go:549: the lease is not released after the last row
== progress-count: caught run_integration_test.go:538: the progress ["1/6 forge" "2/6 clone" "3/6 version" "4/6 lease" "5/6 probe" "6/6 phase"], want ["1/6 forge" "2/6 clone" "3/6 version" "4/6 le
== the-table: caught run_integration_test.go:546: intake/psb.tsv:
5 of 5
```

## Green (2026-10-11, 04:31Z)

Each with exit 0: `go build ./...`, `go vet ./...`, `gofmt -l internal cmd`
(empty), `go test ./...`, `go test -tags=integration ./...` (with
`TestTheTextBlocksOfRunMd`, `run-steps` with `rule-gaps`),
`go test -tags=e2e ./...` (the restart's `rule-gaps`, and "done before" on the
third restart); `adr-lint`, `prd-lint`, `link-lint`, `setup-check`,
`run-discipline-tests`, `git diff --check`.
