# T-xvvd test runs

The stubs were committed alone (68c2e94: `PullBody` returns "", `checkTitle` nil, `OpenPullRequest`, `Ready`, `MarkReady`, `MergeTask`, `RecordMerged`, `PushCloseOut`, `CheckRendered` and `CarryVerify` return their zero values, `withPR` and `explainHead` an error, `mergeMaps` an empty read; the behaviour-preserving extractions `pushChecked` and `setup.ActionsApp` keep the behaviour, and `Sessions.Web` is an inert field). The integration level of `internal/run` was run on that commit (the push and tree characterization tests execute; the log is below). The tests were then written and run against the stubs, red, below. The code and the tests are committed together.

## Red against the stubs

Each line below is an assertion of a test, none a build error and none a panic. The unit run is `go test -count=1 -v ./internal/run` at 16:20:30Z to 16:20:44Z (exit 1, 26 failing tests, every one of them new; the tests of the earlier rows all pass). The integration run is `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -v -tags=integration -timeout 40m ./internal/run -run '<the 14 new tests>'` at 16:20:44Z to 16:24:24Z (exit 1, 15 failing tests; a 16th line of its summary is an earlier test, `TestAVerifierOfAnotherTaskOrRoleOrNoneIsAnError`, that the pattern also matched, and it passes).

### Unit

```
--- FAIL: TestTheHeadIsExplainedOnlyByAChainOfMerges
--- FAIL: TestTheTreeRowAcceptsAMergeNotYetRecordedWithItsLine
--- FAIL: TestTheTreeRowFailsAHeadThatIsNoRecordedMerge
--- FAIL: TestTheTreeRowWithNoMergeKnownIsTheSetupCheckAndAReadOfTheMergesThatFailsIsAFail
--- FAIL: TestTheMergesAreReadFromTheEventsAndFromTheForgeForATaskWithNoEvent
--- FAIL: TestTheTitleIsTheFormOfThePRTitleCheck
--- FAIL: TestThePullBodyIsTheTwoLinesAndTakesFurtherLines
--- FAIL: TestTheDraftIsOpenedOnceWithItsTitleAndBody
--- FAIL: TestAStopAfterTheDraftOpensNoSecondDraft
--- FAIL: TestTheRecordsCommitIsTheBoundSessionsAndNotALaterOnes
--- FAIL: TestATitleOffItsFormIsRefusedBeforeAnyCall
--- FAIL: TestAnAttemptWithNoBoundHeadAndATaskOutsideTheRegisterAreErrorsBeforeAnyCall
--- FAIL: TestAPullRequestOfTheBranchIsReadBeforeAnythingIsOpened
--- FAIL: TestWithPRSetsTheNumberOfOneTaskAndKeepsTheOthers
--- FAIL: TestTheRequiredJobsAreTheActionsEntriesOfTheRequiredChecksRules
--- FAIL: TestTheReadyRuleNamesTheFirstContextOrJobThatIsNotSuccess
--- FAIL: TestARulesetWithNoActionsJobIsNotReadyAndAMalformedOneIsAnError
--- FAIL: TestTheReadRuleFailsClosedOnAnErrorOfTheForgeAndOnABadInput
--- FAIL: TestMarkReadyMarksOnlyAReadyDraftByTheRead
--- FAIL: TestMarkReadyRefusesAMovedHeadAndAPullRequestThatIsNotOpenOrNotRegistered
--- FAIL: TestTheMergeIsAtTheHeadItReadAndWritesTheMergedEvent
--- FAIL: TestAHeadThatMovedBeforeTheMergeIsNotMerged
--- FAIL: TestAMergeRefusedForAMovedHeadReadsThePullRequestAndTheNewStatusesAndMergesNoMore
--- FAIL: TestAFailedMergeWritesNoEvent
--- FAIL: TestALostLeaseOpensMarksAndMergesNothing
--- FAIL: TestAMergedPullRequestIsRecordedOnceAndTheAttemptIsNeverGuessed
```

### Integration

```
--- FAIL: TestTheDraftPullRequestOfABoundHead
--- FAIL: TestTheReadyAndMergeRuleBlocksEachCaseAndNothingIsMarkedOrMerged
--- FAIL: TestTheFourSuccessesMarkTheDraftReadyAndMergeAtTheHeadSHA
--- FAIL: TestAStatusThatChangesAfterTheMarkBlocksTheMerge
--- FAIL: TestAHeadThatMovesAfterTheLastReadIsRefusedAndTheNewHeadIsReadBeforeAnotherCall
--- FAIL: TestAMergeNotYetRecordedPassesTheTreeCheck
--- FAIL: TestAPullRequestFromDraftToMergeAndTheTreeAfterIt
--- FAIL: TestTheCloseOutCommitCarriesTheVerdict
--- FAIL: TestACloseOutThatChangesMoreThanTheTwoFilesIsRefusedAsScope
--- FAIL: TestACloseOutOnAParentThatDoesNotCarryTheVerdictIsRefusedAsVerdict
--- FAIL: TestAVerifierOfAnotherTaskIsAnErrorAndNothingIsPushed
--- FAIL: TestCarryVerifyChecksTheWholeEdgeAndPostsNothingForAHeadThatIsNoCloseOut
--- FAIL: TestTheRenderedCommitCarriesAndTheCloseOutAfterItIsAccepted
--- FAIL: TestARenderedCommitThatIsNotAPermittedChildOfTheVerifiedHeadDoesNotCarry
--- FAIL: TestARenderedCommitOffTheRecordsIsNotAccepted
```

### Added after the first red run

`TestAVerifiedHeadThatIsAMergeCommitCarriesTheVerdict` (condition 4: a verified head that is a merge commit carries the verdict) was written after the run above and run against the same stubs, in a copy of the tree of 68c2e94 with the two inert declarations, at about 16:50Z, exit 1:

```
--- FAIL: TestAVerifiedHeadThatIsAMergeCommitCarriesTheVerdict
    ship_integration_test.go:359: the close-out on a merge commit: bound "", refused ""; want it accepted
```

## Green negatives beside red controls

A test that passes at the stub is a green negative; it is listed, never counted as a red run, and each stands beside the red control of its own table.

- `TestTheTitleIsTheFormOfThePRTitleCheck`: the rows "the control, a title of the form", "each type of the list" and "three characters" are green at the stub (a stub that accepts every title); the twelve rows that must be refused are the red ones above.
- `TestTheRequiredJobsAreTheActionsEntriesOfTheRequiredChecksRules`: the rows "no rule of the type" and "no rule" are green at the stub (it returns none); the four rows that name a job are red.
- `TestTheReadyRuleNamesTheFirstContextOrJobThatIsNotSuccess`: its control, "the control, four successes of the App and the job", is RED at the stub (`Ready = false; want true`), as is the row "an older failure of another source under the App's success"; the other rows are red controls of the reasons.
- `TestTheTreeRowWithNoMergeKnownIsTheSetupCheckAndAReadOfTheMergesThatFailsIsAFail`: its first part, the setup head with no merge known (the old rule, byte for byte), is green at the stub and stays green; the failing line is the read of the merges. In the same test the chain that ends at a commit that the setup check refuses fails at the stub (a head that no walk explains) and is green for the right reason once the walk exists, beside the red control of the code-owner commit as the bottom.
- `TestTheTreeRowFailsAHeadThatIsNoRecordedMerge`: the row "a head above the setup head with no merge known" is a green control (the old words of the setup check); the two rows that expect `<sha> is not a recorded merge` are red.
- `TestACloseOutThatChangesMoreThanTheTwoFilesIsRefusedAsScope` and `TestTheCloseOutCommitCarriesTheVerdict` stand beside each other: the control, a valid close-out, is red at the stub because the stub's push hook never binds, and the refused rows are red because nothing is refused.

## The integration level on the stubs commit (S2)

The command was `TMPDIR=/home/layup/.cache/layup-session/tmp-green <the test binary of 68c2e94> -test.count=1 -test.timeout 45m` in `internal/run` (the binary built with `go test -c -tags=integration`), so that the push and tree characterization tests (`push_integration_test.go`, `forgecheck_integration_test.go`) execute on the extraction `pushChecked`.

It started at 16:01:47Z and ended at 16:24:28Z (22 minutes 41 seconds) with `PASS` and exit 0: the characterization of the extraction holds.
