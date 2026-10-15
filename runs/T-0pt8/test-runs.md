# T-0pt8 — the runs of the tests

Times from the clock of the host, UTC, 2026-10-15; commands as the shell ran them.

## Step T: red for the right reason

- **2026-10-15T17:25:58Z**, before `internal/spec/shape.go` existed: `go test ./internal/spec/` fails to build, with `undefined: Option`, `undefined: Filter`, `undefined: Fact` and the like, in `shape_test.go` (the package has none of `Filter`, `ReadOptions`, `Standing`, `FactsOfSpans` or `OptionsSchema`). The first version of the tests.
- **2026-10-15T17:30:44Z**, after the plan review (condition 3), with the tests expanded and the code still of the first version: `go vet ./internal/spec` fails with `undefined: Constraint` (and `Requirement`, `ADR`, `RequirementsOfDraft`): the new cases for each source of the constraints, the exclusion of the other kinds and of a record that is not Accepted, the replacement of the span facts by the inventory's, the missing row told between two options, an extra row with a pass and no evidence, and an empty result. Then the code, then green.

## Each rule's test fails when that rule alone is broken (condition 3)

One mutant on a copy of `shape.go` at a time, restored after each:

```
a failed row does not drop | caught by TestAnOptionWithAFailedRowIsDropped TestAnOptionWithSeveralBreaksNamesEachInTheOrderOfTheConstraints
a missing row does not drop | caught by TestAnOptionWithAMissingRowIsDropped
a pass with no evidence survives | caught by TestAPassWithNoEvidenceDropsTheOption TestAnOptionWithSeveralBreaksNamesEachInTheOrderOfTheConstraints TestAnExtraRowWithAPassAndNoEvidenceDropsItsOption
the first survivor is chosen when several | caught by TestMoreThanOneSurvivorChoosesNone
a survivor is chosen when none | caught by TestEveryOptionDroppedLeavesNoSurvivor
every kind of requirement applies | caught by TestTheNonFunctionalRequirementsOfTheDraftApplyWithTheirCriterion TestTheNonFunctionalRequirementsOfThePRDApplyToo TestTheFourKindsComeInOneOrder
every fact applies | caught by TestTheConstraintFactsOfTheSpansApplyBeforeTheFirstBet TestTheConfirmedInventoryReplacesTheFactsOfTheSpans
a record that is not Accepted applies | caught by TestOnlyTheAcceptedDecisionRecordsApply
the invariants do not apply | caught by TestTheInvariantsApply TestTheFourKindsComeInOneOrder
the criterion is dropped | caught by TestTheNonFunctionalRequirementsOfTheDraftApplyWithTheirCriterion TestTheNonFunctionalRequirementsOfThePRDApplyToo
an extra constraint is not judged | caught by TestAnExtraRowWithAPassAndNoEvidenceDropsItsOption TestAnExtraConstraintIsJudgedButNotDemanded
an option name is unchecked | caught by TestAnOptionsTableThatBreaksAnOwnerRuleIsRefused
a fail with evidence is accepted | caught by TestAnOptionsTableThatBreaksAnOwnerRuleIsRefused
an empty result is accepted | caught by TestAnOptionsTableThatBreaksAnOwnerRuleIsRefused
a non-fact span is a fact | caught by TestTheFactsOfASpanTableAreItsFactRows
an extra constraint is demanded of every option | caught by TestAnExtraRowWithAPassAndNoEvidenceDropsItsOption TestAnExtraConstraintIsJudgedButNotDemanded
```

## Green and the validation (condition 4), before the freeze

```
go build ./... → exit 0\ngo vet ./... → exit 0\ngo test -count=1 ./... → exit 0\ngo test -count=1 -tags=integration ./internal/spec/ ./internal/tsv/ ./internal/cli/ ./cmd/layup/ → exit 0\nsh docs/adr/adr-lint.sh → exit 0\nsh docs/prd/prd-lint.sh → exit 0\nsh docs/links/link-lint.sh → exit 0\nsh docs/setup/setup-check.sh → exit 0\nsh docs/tests/run-discipline-tests.sh → exit 0\ngit diff --check --cached → exit 0\ngofmt -l . → exit 0\n
```

`go test -count=1 -tags=integration ./...` (the whole module, TMPDIR on the disk, under the load of two other forks): exit 0, 17:28:23 to 17:47:07Z (18 min 44 s). That run began before the last edits of `shape.go` and its tests (17:30 and 17:31), so the four packages that read `internal/spec` or the block test (`internal/spec`, `internal/tsv`, `internal/cli`, `cmd/layup`) were run again at 17:48:06 to 17:48:31Z: exit 0, shown above.
