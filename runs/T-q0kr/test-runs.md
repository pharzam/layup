# T-q0kr test runs

The stubs were committed alone (e17c09d, 1cc2b3f, inert: every mapper returns `Status{}`, `PostStatus` `false, nil`, `ReadGates` `nil, nil`, `GatesSchema` empty, `gatesComplete` nil; the lease and kept-copy extractions keep their behaviour). The tests were then written and run against them, red, below (10:59Z to 11:02Z); the code and the tests are committed together.

## Red against the stubs

Each line below is an assertion of a test, none a build error. Unit `go test ./internal/run/` (exit 1, 22 failing tests), then the integration tests of this task (exit 1, 7 failing tests).

```
--- FAIL: TestEachRowCombinationMapsToTheGatesStatus
    gates_test.go:46: every kind passes: GatesStatus = {Context: State: Description:}, want {Context:layup/gates State:success Description:pass}
--- FAIL: TestAGatesTableIsReadByItsRulesAndRefusesEachBrokenOne
    gates_test.go:56: the control, a good table of two heads: 0 rows, <nil>
--- FAIL: TestASavedTableIsCompleteOnlyWhenItsKindsAreTheManifestOfItsBase
    gates_test.go:98: a state that differs from the manifest: the table was complete
--- FAIL: TestTheGatesTableOfAHeadRebuildsTheFindings
    gates_test.go:120: GatesTable = {Base: Head: Rows:[]}, false, <nil>; want the two rows of the head
--- FAIL: TestTheSpecStatusIsSuccessOnlyWithNoFailRowAndACompleteTable
    gates_test.go:170: every rule passes: SpecStatus = {Context: State: Description:}, want {Context:layup/spec State:success Description:pass}
--- FAIL: TestASpecTableIsReadByItsRules
    gates_test.go:195: SpecTable = [], false, <nil>; want the ten rows (rule, result, detail) of the head
--- FAIL: TestADescriptionIsCutAtOneHundredFortyCharacters
    statuses_test.go:48: 141 are cut to 139 and an ellipsis: fit = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
--- FAIL: TestOnlyTheNewestStatusOfTheAppCountsAndOtherSourcesNeverDo
    statuses_test.go:61: appStatus = {Context: State: Description: CreatorID:0 CreatorLogin:}, false; want the newest status of the App (failure, new), not the other sour
--- FAIL: TestAStatusIsPostedOnceForAnUnchangedStateAndAgainWhenItChanges
    statuses_test.go:79: the first post: false, <nil>; want posted
--- FAIL: TestAStatusOfAnotherSourceDoesNotCountAsPosted
    statuses_test.go:101: a status of another creator: false, <nil>; want the App's posted
--- FAIL: TestADescriptionOverTheBoundIsFitAndPostedOnce
    statuses_test.go:113: call 1: false, <nil>; want posted=true and no refusal
--- FAIL: TestAStatusIsRefusedBeforeAnyRequestWhenItsInputIsBad
    statuses_test.go:134: a context of no status of ours: false, <nil>; want a refusal before the request
--- FAIL: TestALostLeasePostsNoStatus
    statuses_test.go:154: a lost lease: false, <nil>, 0 statuses; want a LostError and nothing posted
--- FAIL: TestAFailedReadOrPostIsAnErrorAndAnotherCallPostsAfterIt
    statuses_test.go:168: a failing Statuses: false, <nil>; want an error
--- FAIL: TestTheRulesStatusIsTheResultOfTheBasicRule
    statuses_test.go:195: RulesStatus({State:success Reason:no rule path Path:}) = {Context: State: Description:}, want {Context:layup/rules State:success Description:no 
--- FAIL: TestTheSessionsOfTheRunKnowTheBotOfTheApp
    statuses_test.go:211: state.sessions().Bot = 0, want 9001 (the ID that step forge read)
--- FAIL: TestTheRegisterCommitOfAHeadIsTheRecordsOfItsBoundSessionOrTheCallers
    statuses_test.go:244: the session bound to the head: "", <nil>; want "5555555555555555555555555555555555555555", ok=true
--- FAIL: TestNoHarnessOutsideTheAuthorsIsFailureNotActive
    verifystatus_test.go:91: VerifyAt = , {Context: State: Description:}; want failure, not-active: no harness outside the authors, at the head
--- FAIL: TestAWaitingVerifierIsPendingAndARunningOneIsToo
    verifystatus_test.go:107: a waiting verifier: , {Context: State: Description:}
--- FAIL: TestVerifyIsSuccessOnlyForNothingMaterialAtTheReviewedHeadWithAResultThatIsDone
    verifystatus_test.go:148: nothing material at the reviewed head: , {Context: State: Description:}; want {Context:layup/verify State:success Description:nothing materi
--- FAIL: TestAVerifierThatIsAnAuthorOrNotAdmittedIsNotIndependent
    verifystatus_test.go:158: a clean record of an author: , {Context: State: Description:}; want {Context:layup/verify State:failure Description:not-active: the verifier
--- FAIL: TestAVerifierOfAnotherTaskOrRoleOrNoneIsAnError
    verifystatus_test.go:197: a session of another task: , {Context: State: Description:}, no error

--- FAIL: TestTheGatesBlockEqualsItsGoSchema
    statuses_integration_test.go:29: schema task-gates: the name: the block has "task-gates"; the Go schema has ""
--- FAIL: TestTheGatesStatusIsPostedOnceFromTheRecordedTableAndItsOutputIsKept
    statuses_integration_test.go:88: GatesAt = {Context: State: Description:}, <nil>; want {Context:layup/gates State:success Description:clear: pending: no product path}
--- FAIL: TestAFailedPostIsRecoveredFromTheSavedTableWithoutRunningTheGateAgain
    statuses_integration_test.go:122: GatesAt = {Context: State: Description:}, <nil>; want {Context:layup/gates State:failure Description:bad: fail: exit 3}
--- FAIL: TestAPendingKindWithAProductPathAndANotActiveKindAreFailures
    statuses_integration_test.go:159: a pending kind and a product path that changed: {Context: State: Description:}, <nil>; want {Context:layup/gates State:failure Descr
--- FAIL: TestABaseWithNoManifestIsAPostedErrorAndWritesNoRow
    statuses_integration_test.go:172: GatesAt = {Context: State: Description:}, <nil>; want error, error: docs/gates.tsv at the base …
--- FAIL: TestTheRulesStatusReadsTheRegisterOfTheRecordsCommit
    statuses_integration_test.go:189: RulesAt = {Context: State: Description:}, <nil>; want {Context:layup/rules State:failure Description:rule path: AGENTS.md} (the regi
--- FAIL: TestTheVerifyStatusFollowsTheVerifierFromTheRecords
    statuses_integration_test.go:223: VerifyAt = , {Context: State: Description:}, <nil>; want {Context:layup/verify State:success Description:nothing material in scope} 
```

## Green

After the code: `go test ./internal/run/` exit 0; the seven named integration tests exit 0 (about 56 s); `TestTheGatesBlockEqualsItsGoSchema` moved `task-gates` from notYetBuilt to built.

## Compiled mutants

Each built and vetted (a build failure is "not run" and was redone), was run against the unit level or the named integration tests, and was restored byte for byte (checked with a read-back compare by the script) before the next. All 23 are killed; none is claimed equivalent.

| # | Mutant | Killed by |
| - | ------ | --------- |
| M1 | clear counts as a failure | TestEachRowCombinationMapsToTheGatesStatus |
| M2 | not-active counts as a pass | TestEachRowCombinationMapsToTheGatesStatus |
| M3 | the compare ignores the description | TestAStatusIsPostedOnceForAnUnchangedStateAndAgainWhenItChanges |
| M4 | any creator counts | TestOnlyTheNewestStatusOfTheAppCountsAndOtherSourcesNeverDo, TestAStatusOfAnotherSourceDoesNotCountAsPosted |
| M5 | the oldest App status is compared | TestOnlyTheNewestStatusOfTheAppCountsAndOtherSourcesNeverDo |
| M6 | the lease read is skipped | TestALostLeasePostsNoStatus |
| M7 | the gate runs again for a recorded head | TestTheGatesStatusIsPostedOnceFromTheRecordedTableAndItsOutputIsKept, TestAFailedPostIsRecoveredFromTheSavedTableWithoutRunningTheGateAgain |
| M8 | the table is not committed before the status | the same two (the recovery test) |
| M9 | the spec table may hold the rules in any order | TestTheSpecStatusIsSuccessOnlyWithNoFailRowAndACompleteTable |
| M10 | no harness outside the authors gives pending | TestNoHarnessOutsideTheAuthorsIsFailureNotActive |
| M11 | success without the reviewed head | TestVerifyIsSuccessOnlyForNothingMaterialAtTheReviewedHeadWithAResultThatIsDone |
| M12 | the count of findings plus one | the same |
| M13 | the register is read at the head | TestTheRulesStatusReadsTheRegisterOfTheRecordsCommit |
| M14 | fit cuts at 141 | TestADescriptionIsCutAtOneHundredFortyCharacters |
| M15 | ReadGates accepts heads that are apart | TestAGatesTableIsReadByItsRulesAndRefusesEachBrokenOne |
| M16 | the author check of the verifier is dropped | TestAVerifierThatIsAnAuthorOrNotAdmittedIsNotIndependent |
| M17 | the admission check of the verifier is dropped | the same |
| M18 | the constructor forgets the bot | TestTheSessionsOfTheRunKnowTheBotOfTheApp |
| M19 | the restart forgets the bot | TestAFailedPostIsRecoveredFromTheSavedTableWithoutRunningTheGateAgain |
| M20 | a table with a kind omitted is complete | TestASavedTableIsCompleteOnlyWhenItsKindsAreTheManifestOfItsBase |
| M21 | a pass may carry a reason | TestAGatesTableIsReadByItsRulesAndRefusesEachBrokenOne |
| M22 | a verifier of another task is accepted | TestAVerifierOfAnotherTaskOrRoleOrNoneIsAnError |
| M23 | the first fail is not the first row | TestTheSpecStatusIsSuccessOnlyWithNoFailRowAndACompleteTable |
