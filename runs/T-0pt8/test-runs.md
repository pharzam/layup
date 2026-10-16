# T-0pt8 — the runs of the tests

Times from the clock of the host, UTC, 2026-10-15.

## Red for the right reason

- 17:25:58Z, before `shape.go` existed: `go test ./internal/spec/` fails to build with `undefined: Option`, `Filter`, `Fact` (the first tests).
- 17:30:44Z, the tests expanded by condition 3 of the plan review (each source of the constraints, the other kinds, a record that is not Accepted, the facts replaced, the missing row between two options, an extra row, an empty result), the code still the first: `undefined: Constraint`, `Requirement`, `ADR`, `RequirementsOfDraft`. Then the code.
- 18:05:14Z, round 1's two fixes (an amended record applies; a fact of the answers keeps its `Q-NNN`), their cases on the old behaviour: `:135: the records: the accepted and the amended, not the others: ["ADR-0011"], want ["ADR-0011" "ADR-0007"]` and `:135: a fact of the answers by its own ID: ["F-Q-004" "F-6"], want ["Q-004" "F-6"]`.

## Each rule's test fails when that rule alone is broken

One mutant on a copy of `shape.go`, restored after each; each text change whole (a newline and a tab as a space), and the tests that fail; none builds badly and none is alive (re-run at 18:35Z, after round 2):

- a failed row does not drop: `case r[optResult] == "fail":` → `case false:`: TestEachBreakDropsOnlyItsOption TestAnOptionWithSeveralBreaksNamesEachInTheOrderOfTheConstraints
- a missing row does not drop: `case r == nil:` → `case false:`: TestEachBreakDropsOnlyItsOption
- a pass with no evidence survives: `case r[optEvidence] == "":` → `case false:`: TestEachBreakDropsOnlyItsOption TestAnOptionWithSeveralBreaksNamesEachInTheOrderOfTheConstraints
- several survivors give a choice: `if len(r.Survivors) == 1 {` → `if len(r.Survivors) >= 1 {`: TestTheSurvivorsAreChosenOnlyWhenThereIsOne
- no survivor gives a choice: `if len(r.Survivors) == 1 {` → `if len(r.Survivors) <= 1 {`: TestTheSurvivorsAreChosenOnlyWhenThereIsOne
- every requirement kind applies: `if r.Kind == "non-functional" {` → `if true {`: TestTheConstraintsThatApplyComeFromTheRecordsInOneOrder
- every fact applies: `if f.Class == "constraint" {` → `if true {`: TestTheConstraintsThatApplyComeFromTheRecordsInOneOrder TestTheCallerChoosesTheSourceOfTheFacts
- an unaccepted record applies: `if a.Status == "Accepted" || strings.HasPrefix(a.Status, "Accepted. Amended by ") {` → `if true {`: TestTheConstraintsThatApplyComeFromTheRecordsInOneOrder
- an amended record is left out: `if a.Status == "Accepted" || strings.HasPrefix(a.Status, "Accepted. Amended by ") {` → `if a.Status == "Accepted" || strings.HasPrefix("", "x") {`: TestTheConstraintsThatApplyComeFromTheRecordsInOneOrder
- a fact of the answers gets F-: `if strings.HasPrefix(f.Number, "Q-") {` → `if strings.HasPrefix("", "x") {`: TestTheConstraintsThatApplyComeFromTheRecordsInOneOrder
- the invariants do not apply: ` for _, v := range invariants { out = append(out, Constraint{v, ""}) }` → ` _ = invariants`: TestTheConstraintsThatApplyComeFromTheRecordsInOneOrder
- the criterion is dropped: `out = append(out, Constraint{r.ID, r.Criterion})` → `out = append(out, Constraint{r.ID, ""})`: TestTheConstraintsThatApplyComeFromTheRecordsInOneOrder
- an extra constraint is not judged: ` if !listed[c] { listed[c] = true judge(c, byOpt[o.Name][c]) }` → ` if !listed[c] { listed[c] = true }`: TestEachBreakDropsOnlyItsOption
- an extra constraint is demanded of every option: ` for _, c := range constraints { listed[c] = true` → ` all := append([]string(nil), constraints...) for _, r := range rows { all = append(all, r[optConstraint]) } for _, c := range all { listed[c] = true`: TestEachBreakDropsOnlyItsOption TestAnExtraConstraintIsJudgedButNotDemanded TestAMissingRowDropsOnlyTheOptionThatLacksIt TestAnOptionWithSeveralBreaksNamesEachInTheOrderOfTheConstraints
- a missing row is read from another option: ` listed[c] = true judge(c, byOpt[o.Name][c])` → ` listed[c] = true var any []string for _, m := range byOpt { if x, ok := m[c]; ok { any = x } } judge(c, any)`: TestEachBreakDropsOnlyItsOption TestAMissingRowDropsOnlyTheOptionThatLacksIt TestAnOptionWithSeveralBreaksNamesEachInTheOrderOfTheConstraints TestTheSurvivorsAreChosenOnlyWhenThereIsOne
- an option name is unchecked: `!optionName.MatchString(r[optName])` → `false`: TestAnOptionsTableThatBreaksAnOwnerRuleIsRefused
- a fail with evidence is accepted: `r[optResult] == "fail" && r[optEvidence] != ""` → `false`: TestAnOptionsTableThatBreaksAnOwnerRuleIsRefused
- an empty result is accepted: `r[optResult] == ""` → `false`: TestAnOptionsTableThatBreaksAnOwnerRuleIsRefused
- a non-fact span is a fact: ` if r[colFact] != "" {` → ` if true {`: TestTheFactsOfASpanTableAreItsFactRows

## Green and the validation, at the content of the head

```
go build ./... → exit 0
go vet ./... → exit 0
go test -count=1 ./... → exit 0
sh docs/adr/adr-lint.sh → exit 0
sh docs/prd/prd-lint.sh → exit 0
sh docs/links/link-lint.sh → exit 0
sh docs/setup/setup-check.sh → exit 0
sh docs/tests/run-discipline-tests.sh → exit 0
gofmt -l . → exit 0
git diff --check origin/main HEAD → exit 0
```

`go test -count=1 -tags=integration ./...` (the whole module, TMPDIR on the disk, under the load of two other forks): exit 0, 18:06:28 to 18:22:40Z (16 min 12 s), at the content of the head after the fixes of round 1. The first run, before them (17:28:23 to 17:47:07Z), was also exit 0.
