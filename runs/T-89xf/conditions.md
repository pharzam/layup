# T-89xf: how each condition is met

Plan review 6262560771 (conditions 1 to 8), the Operator's decisions 6262587821, and conditions 1, 5 and 8 of the plan review of #367 (6252684822).

| Condition | Evidence |
| --- | --- |
| 1. The merge proof names the sets T and P, checks the counterexample, and checks the produced merge before it is pushed | `internal/run/closeout.go:290-320`, `TestAMergeThatTouchesATaskPathOrLeavesTheDefaultIsNotCarried`, `internal/run/baseup.go:181-196`, `TestAHeadThatDoesNotCarryTheVerdictIsNotMerged` |
| 2. The push of the merge is announced before it is made; the sessionless consumers are audited | `internal/run/baseup.go:197-225`, `TestAPushedMergeCommitIsAdoptedOnARestart`, `TestAMergeAnnouncedAndNotPushedIsMadeAgainWhenTheBaseMovesAgain`, `internal/run/repair.go:56`, `docs/spec/loop.md` |
| 3. The intent record is taken; the pull request is read before any close; the lease is read before every write of the forge; the duplicate-issue window is stated | `internal/run/redo.go:169-250`, `TestARestartAfterEachActOfTheRedoCompletesItOnce` (with its driver entries), `TestALostLeaseOpensNoIssueAndWritesNothingAfterIt`, `TestAPullRequestMergedAfterItsRedoStops` |
| 4. The replacement keeps the binding and the pins, composes the issue from kept inputs, and the driver uses the rebased base | `TestTheFirstAttemptOfARedoneTaskStartsAtTheRebasedBase`, `TestARestartAfterEachActOfTheRedoCompletesItOnce`, `internal/run/redo.go:130-166`, `TestTheNewIssueIsComposedFromTheInputsTheRecordsKeep`, `TestTheRedoIsOneCommitOfTheRegisterTheTwoEventsAndTheBindings` |
| 5. The older attempt's pull request is closed only when the register names the newer one | `internal/run/redo.go:256-295`, `TestTheOlderAttemptsPullRequestIsClosedOnce` |
| 6. A merge failure that is no conflict and a failed read are told from a conflict; the commit that holds the register and both events is observed; the fixtures of U1 are valid IDs | `TestMergeTellsAConflictFromEveryOtherFailure`, `TestAnOverlappingMovedBaseRedoesTheTask`, `TestAFailedGitReadIsAnErrorNeverARedo`, `internal/records/session_test.go:164-165` |
| 7. The eligible heads and the overlap test are specified; the rendered commit is a stop (the Operator's Q14) | `internal/run/closeout.go:156-195`, `TestTheCloseOutHeadIsMergedWhenTheBaseLeavesItsFilesAndRedoneWhenItDoesNot`, `TestARedoThatWouldDropTheRenderedCommitIsAStop`, `docs/spec/loop.md` |
| 8. The sources are described as landed; the size is recounted inside 2,800 lines and 30 files | `runs/T-89xf/test-runs.md:1-30`, `docs/plan/README.md` |
| Operator: the goal count is 3 (carry, redo, the superseded pull request) | `TestADisjointMovedBaseIsMergedAndCarriesTheVerdict`, `TestAnOverlappingMovedBaseRedoesTheTask`, `TestTheOlderAttemptsPullRequestIsClosedOnce` |
| Parent 1: a reviewed plan supplies no callable API; the landed names are refreshed | `runs/T-89xf/test-runs.md:1-30` |
| Parent 5: a clean merge precedes the carry; a failed git read is no evidence for a redo; recovery after each close and open | `TestAFailedGitReadIsAnErrorNeverARedo`, `TestARestartAfterEachActOfTheRedoCompletesItOnce` (with its driver entries), `internal/run/baseup.go:64-178` |
| Parent 8: the detail of closed and rebased is a task ID, with red cases | `internal/records/session.go:272`, `TestTheDetailOfClosedAndRebasedIsATaskID` |
