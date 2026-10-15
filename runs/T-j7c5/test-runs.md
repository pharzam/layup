# T-j7c5 — the test runs

Task `T-j7c5` (#330), base `8802e6c`. Times are UTC on 2026-10-15, read from the clock after each step. The plan review (comment 6187488408) and the answer (6187493580, 08:06:37Z) came first. **No declaration or stub holds a fix, and none was needed:** the first red tests fail by an assertion on the existing code. The code under test is `3aa9dd3`.

## The forge's answers (first hand)

The failed attempt, run 6 of the demo of `M2d` (`T-sdtb`, #284), 07:57:04Z to 07:57:17Z, on `pharzam/layup-uat-m2d` (not a fork); the ruleset 25145511 was applied at 07:56:50Z and the answers were read at 07:58:55Z.

| Answer | The rules of the records branch, as GitHub gave them |
| ------ | --------------------------------------------------- |
| The file sent (`out/ruleset-records.json`) | `[{"type":"update","parameters":{"update_allows_fetch_and_merge":false}},{"type":"deletion"},{"type":"non_fast_forward"}]` |
| The response of the apply, and `GET .../rulesets/25145511` | `[{"type":"update"},{"type":"deletion"},{"type":"non_fast_forward"}]` |
| `GET .../rules/branches/layup-records` | the same three rules, each with `ruleset_source_type` `Repository`, `ruleset_source` `pharzam/layup-uat-m2d` and `ruleset_id` 25145511, and no `parameters` |
| `setup/forge-check.tsv`, commit `4313def` of `layup-records` | `rules:layup-records fail update: update_allows_fetch_and_merge is missing`; the seven other rows `pass` |
| `GET .../rules/branches/main`, the rule `pull_request`, read again at 08:19:48Z | `{"allowed_merge_methods":["merge"],"dismiss_stale_reviews_on_push":false,"require_code_owner_review":false,"require_extra_approval_for_unattributed_changes":false,"require_last_push_approval":false,"required_approving_review_count":0,"required_review_thread_resolution":true,"required_reviewers":[]}`: the `false` parameters of the file are echoed |

## Step 3: the tests, red first, on the code of the base

| Level | Run | Result, read for its cause |
| ----- | --- | -------------------------- |
| unit and job (`internal/run`) | 08:08:41Z | 3 FAIL, each by an assertion: the read-back case "a parameter that GitHub omits at its default" (`false, "update: update_allows_fetch_and_merge is missing"; want true`); `TestTheStepForgeCheckCommitsEightRowsAndIsDoneWhenEachIsAPass` (the step is `fail` with that detail, want `done`); `TestEachCheckThatFailsIsItsRowAndTheTableStillHasEightRows`, 20 of its 28 cases, each by the row `rules:layup-records` that reads GitHub's form as a difference. **Guards that pass at red:** `TestTheFakeOfTheJobGivesTheRecordsRulesAsGitHubDoes` (it describes the fake), the seven other new read-back cases, and 8 cases of the fail table (the three new, and five that expect that row to fail already); the mutants below prove their assertions. |
| the world (integration) | 08:09:39Z | FAIL (42.9 s): the first restart's `forge-check` is `fail`, "rules:layup-records: update: update_allows_fetch_and_merge is missing", as run 6 had it; the assertion on the fake's answer passes (a guard) |
| the e2e | 08:10:06Z | FAIL (19.0 s): stage 4, exit 1, `forge-check fail rules:layup-records: update: update_allows_fetch_and_merge is missing` |

After the code (08:10:20Z) each level is green: unit and job 08:10:29Z, the world 08:11:50Z (75.5 s: it makes two restarts more), the e2e 08:12:25Z (32.7 s). A later case that names the omitted default after another parameter (`z_other`) was renamed before the commit (08:14:51Z), so that the mutant that ends the comparison at the first omitted parameter is caught.

## Mutants (each alone on a copy of the tree, a build failure counted as not run; the control is the green run above)

10 on `3aa9dd3`, from 08:15:50Z to 08:19:14Z, each caught, by the test that names it:

| Caught by | Mutants |
| --------- | ------- |
| `TestTheReadBackComparesEachRuleOfTheFileWithOneReturnedRule` | any-missing-false-passes, the-table-keyed-by-the-name-only, the-first-omitted-parameter-ends-the-comparison |
| that test and the job's two (`...IsDoneWhenEachIsAPass`, `...TheTableStillHasEightRows`) | the-table-entry-removed, the-default-is-true, the-type-is-not-passed |
| that test and `...TheTableStillHasEightRows` | an-omitted-true-passes |
| `TestTheFakeOfTheJobGivesTheRecordsRulesAsGitHubDoes`, and the fail table | the-job-fake-echoes-the-file |
| `TestTheStepForgeCheckChecksTheForgeAndCommitsItsTableAtEachRestart` (the assertion on the world's fake) | the-world-fake-echoes-the-file |
| `TestTheStepsOfTheScaffold` (the assertion on the e2e's fake) | the-e2e-fake-echoes-the-file |

## The gate on `3aa9dd3`, from 08:19:27Z (each exit code)

go build ./... 0; go vet ./... 0; gofmt -l internal cmd (empty) 0; go test ./... 0; go test -tags=integration ./... 1; go test -tags=e2e ./... 0; adr-lint 0; prd-lint 0; link-lint 0; run-discipline-tests 0; nested-checkout-check 0; setup-check 0; git diff --check origin/main HEAD 0. The integration level (`go test -tags=integration ./...`) exited 1 on its first run: `TestProcessTimesTheFirstOutput` of `internal/session`, a timing test that a parallel package's load broke by 0.23 ms ("the end 0.2 s later": 199.77 ms), the flake class of #291 (open); the test passes alone three times, and the whole level, run again from 08:34:12Z to 08:45:30Z, exited 0.
