# T-gkvi — test runs

The probe of the forge check on GitHub's real refusal (#323). The plan, its review and the author's answer are comments on #323;
the answer's condition 1 sets the order below: every regression test is written first, in the order unit, job, adapter, world and
e2e, and run red before any production change. Times are UTC, 2026-10-14, from `date -u`; each red run is read for its cause.

## The skeleton (12:00:24Z)

Declared before the first test, and containing no fix: the three fields of `forge.Rule` (`RulesetSourceType`, `RulesetSource`,
`RulesetID`: data that nothing reads), a function `probeAttribution` whose body is `return nil` (no fault is ever found), a function
`rulesetID` whose body is `return 0, nil` (no ID is ever read), and the third parameter of `probeVerdict`, which its old body ignores
and which still requires the ruleset's name. Their purpose is that each test below compiles and fails by an assertion, not by a build
error. The old tests pass on it (`go test ./internal/run/ ./internal/forge/...`, 12:00:24Z: `ok`).

## Red, before any production code

| # | Level | Time | Command | Result |
| - | ----- | ---- | ------- | ------ |
| 1 | unit and job | 12:01:55Z to 12:01:59Z | `go test -count=1 ./internal/run/` | 7 tests fail, each by an assertion |
| 2 | adapter | 12:02:26Z to 12:02:28Z | `go test -count=1 -tags=integration ./internal/forge/github/` | 3 tests fail |
| 3 | world | 12:03:09Z to 12:03:57Z | `go test -count=1 -tags=integration -run 'TestTheProbeIsAPassOnlyForGitHubsRefusalOfTheProbeRef|TestTheStepForgeCheckChecksTheForgeAndCommitsItsTableAtEachRestart' ./internal/run/` | 2 tests fail |
| 4 | e2e | 12:04:02Z to 12:04:34Z | `go test -count=1 -tags=e2e -timeout 10m -run TestTheStepsOfTheScaffold ./cmd/layup/` | fails at stage 4 |

1. **Unit and job.** The verdict test fails in 4 of its 18 cases: GitHub's refusal of the probe ref (the real form, which names no
   ruleset), the line without its final period, the line that ends with a carriage return, and the refusal that the attribution does not
   explain: the old rule says "probe not run" for each. The other 14 cases pass on arrival and are **guards** of what stays (another ref, no
   `GH013`, the accepted push, the remote that cannot be reached, a failure that is no refusal, a call that did not start, no standard
   error, and the old invented text that holds the name). The attribution test fails in 8 of 12 cases (the four that expect no fault pass
   on the stub: guards), and the test of the ID in 11 of 11. In the job's world the default refusal is now GitHub's real form, so the old
   rule fails the probe in every case: `TestTheStepForgeCheckCommitsEightRowsAndIsDoneWhenEachIsAPass` (1 assertion),
   `TestEachCheckThatFailsIsItsRowAndTheTableStillHasEightRows` (20 assertions: the probe row fails where the case expects a pass or
   the attribution's detail) and `TestADetailIsTheFirstLineOfItsReasonInValidUTF8WithDIRForTheHost` (3). The test of the one read fails in 3
   of its 5 cases (the work area, the missing ruleset file and the missing `setup/tree.tsv` read the rules of `layup-probe` 0 times); the
   other two pass on arrival (guards).
2. **Adapter.** The old adapter reads no ruleset field: `TestEffectiveRulesGivesEachActiveRuleOfABranchOverEveryPage`,
   `TestEffectiveRulesOfABranchThatDoesNotExistYetGivesTheRulesThatNameIt` (their rules gain the three fields) and the new
   `TestEffectiveRulesGivesTheRulesetOfEachRuleAndZeroValuesWhenTheAnswerLacksIt` (a `ruleset_id` that is a string gives no error) fail.
3. **World.** `TestTheProbeIsAPassOnlyForGitHubsRefusalOfTheProbeRef` and
   `TestTheStepForgeCheckChecksTheForgeAndCommitsItsTableAtEachRestart`: the bare repository refuses as GitHub does, and the old rule
   gives "probe not run: remote: error: GH013: Repository rule violations found for refs/heads/layup-probe."
4. **e2e.** The hook refuses in GitHub's form and the fake forge gives each rule its ruleset: the restart of stage 4 ends `fail` at
   `forge-check` (exit 1).

## Green (the production code written after the last red run, 12:05Z)

- `go test -count=1 ./...` (12:05:40Z): every package `ok`.
- The forge integration tests, `go test -count=1 -tags=integration ./internal/forge/...`: `ok`.
- The world, 12:05Z to 12:07Z: the first run failed in the test's own sequence, not in the code (the block that adds another ruleset's
  `creation` rule left a table with a failing row, so the next restart is of the first-start chain, not of the later one: the test
  now runs a restart that passes after it, and counts five tables). Run again, 12:07:31Z to 12:08:50Z: `ok` (73.7 s).
- The e2e, 12:08:55Z to 12:09:30Z: `ok` (33.0 s).

## The mutants (12:10:16Z to 12:15:31Z)

Each alone on a copy of the tree with no `.git` (`tar --exclude=.git`), built first (a build that failed would be "not run"), with the
tests that can see it, run until one fails: **18 caught, 0 survived, 0 not run, of 18.**

| Mutant | What it changes | Caught by |
| ------ | --------------- | --------- |
| `attribution-always-nil` | `probeAttribution` returns nil at its first line | unit: `TestAProbeIsAttributedOnlyWhenEachCreationRuleOfTheRefIsTheDefaultRulesets`; world: `TestTheStepForgeCheckChecksTheForgeAndCommitsItsTableAtEachRestart` |
| `ref-not-required` | the refusal pattern is `GH013` alone | unit: `TestTheVerdictOfTheProbeIsAPassOnlyForGitHubsRefusalOfTheProbeRefThatARulesetExplains` |
| `ref-as-prefix` | the pattern has no end of line after the ref | unit: `TestTheVerdictOfTheProbeIsAPassOnlyForGitHubsRefusalOfTheProbeRefThatARulesetExplains` |
| `id-compared-with-zero` | a creation rule is the default's when its ID is not 0, not when it is the default's ID | unit: `TestAProbeIsAttributedOnlyWhenEachCreationRuleOfTheRefIsTheDefaultRulesets` |
| `every-type-is-a-creation` | every rule type, not `creation` alone, is read as a creation rule | unit: `TestAProbeIsAttributedOnlyWhenEachCreationRuleOfTheRefIsTheDefaultRulesets` |
| `no-creation-accepted` | no `creation` rule is no fault | unit: `TestAProbeIsAttributedOnlyWhenEachCreationRuleOfTheRefIsTheDefaultRulesets` |
| `source-type-ignored` | the source type of a rule is not read | unit: `TestAProbeIsAttributedOnlyWhenEachCreationRuleOfTheRefIsTheDefaultRulesets` |
| `probe-rules-read-twice` | the row `rules:layup-probe` reads the rules again | unit: `TestTheRulesOfLayupProbeAreReadOnceWhateverTheWorkAreaAndTheRulesetFile` |
| `old-name-required-again` | the refusal must also hold `layup: the default branch` | unit: `TestTheVerdictOfTheProbeIsAPassOnlyForGitHubsRefusalOfTheProbeRefThatARulesetExplains`; world: `TestTheProbeIsAPassOnlyForGitHubsRefusalOfTheProbeRef` |
| `id-of-the-records-output` | the ID is read from `records.json` | unit: `TestTheStepForgeCheckCommitsEightRowsAndIsDoneWhenEachIsAPass` |
| `probe-rules-read-after-the-area` | the rules of `layup-probe` are read only when the work area is sure | unit: `TestEachCheckThatFailsIsItsRowAndTheTableStillHasEightRows` |
| `fault-of-the-attribution-ignored` | `probeVerdict` ignores the fault of the attribution | unit: `TestTheVerdictOfTheProbeIsAPassOnlyForGitHubsRefusalOfTheProbeRefThatARulesetExplains`; world: `TestTheStepForgeCheckChecksTheForgeAndCommitsItsTableAtEachRestart` |
| `id-of-zero-accepted` | an `id` of 0 is an ID | unit: `TestTheIDOfTheDefaultRulesetIsThePositiveIntegerIDOfItsApplyOutput` |
| `id-read-as-a-float` | the `id` is decoded as a float (no `UseNumber`) | unit: `TestTheIDOfTheDefaultRulesetIsThePositiveIntegerIDOfItsApplyOutput` |
| `rules-not-read-when-the-attribution-fails-open` | a failed read of the rules is no fault of the attribution | unit: `TestEachCheckThatFailsIsItsRowAndTheTableStillHasEightRows` |
| `adapter-drops-the-id` | the adapter leaves out `RulesetID` | adapter: `TestEffectiveRulesGivesEachActiveRuleOfABranchOverEveryPage`; world: `TestTheStepForgeCheckChecksTheForgeAndCommitsItsTableAtEachRestart` |
| `adapter-drops-the-source-type` | the adapter leaves out `RulesetSourceType` | adapter: `TestEffectiveRulesGivesEachActiveRuleOfABranchOverEveryPage` |
| `adapter-drops-the-source` | the adapter leaves out `RulesetSource` | adapter: `TestEffectiveRulesGivesEachActiveRuleOfABranchOverEveryPage` |
