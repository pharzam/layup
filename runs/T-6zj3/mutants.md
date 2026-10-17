# T-6zj3 — the compiled mutants (#361)

A mutant keeps `go build ./...`, `go vet ./internal/run` and `go vet -tags=integration ./internal/run` green; a mutant that does not build is **not run** and is never counted as caught. A kill is a test named by its own `--- FAIL` line in the mutant's output. The file is written back from its bytes (no `git checkout`) and its SHA-256 is compared with the hash taken before the change (the column "restored").

| Run | Time (UTC) | Mutant | File | Result | Level | Killed by | File hash before (12) | Restored |
| --- | ---------- | ------ | ---- | ------ | ----- | --------- | --------------------- | -------- |
| 1 | 16:13:27 | F1 passes-at-head not required | `frozen.go` | killed | unit | TestAListedTestThatDoesNotFailAtTheBaseMakesTheHandoffInvalid | fec9653ee4ec | True |
| 1 | 16:13:38 | F2 fails-at-base not required | `frozen.go` | killed | unit | TestAListedTestThatDoesNotFailAtTheBaseMakesTheHandoffInvalid | fec9653ee4ec | True |
| 1 | 16:13:50 | F3 empty list is red-green | `frozen.go` | killed | unit | TestAnEmptyListFreezesNothingAndIsInvalid | fec9653ee4ec | True |
| 1 | 16:14:00 | F4 caller-independent: tests-red-green forced true | `handoffpost.go` | killed | integration | TestAChangedFrozenSourceIsAFindingAndTheHandoffIsInvalid, TestAListedTestThatIsNotRedThenGreenMakesTheHandoffInvalidAndFreezesNothing | 547fa32c28f9 | True |
| 1 | 16:14:37 | F5 hash read at the clone HEAD, not the bound head | `frozen.go` | killed | integration | TestTheFrozenListHoldsTheHashesOfTheFilesAtTheBoundHead | fec9653ee4ec | True |
| 1 | 16:15:15 | F6 amendment review need not start after the plan | `frozen.go` | killed | unit | TestAStaleVerifiersRecordIsNoReviewOfTheNewPlan | fec9653ee4ec | True |
| 1 | 16:15:39 | F7 amendment: plan role not required for the plan | `frozen.go` | SURVIVED |  |  | fec9653ee4ec | True |
| 1 | 16:15:58 | F8 amendment: review role not required for the review | `frozen.go` | SURVIVED |  |  | fec9653ee4ec | True |
| 1 | 16:16:13 | F9 refreeze with no amendment | `frozen.go` | build failed: not run |  |  | fec9653ee4ec | True |
| 1 | 16:16:14 | F10 harness of an unchanged test not kept | `frozen.go` | killed | unit | TestTheHarnessOfAnUnchangedTestIsKept | fec9653ee4ec | True |
| 1 | 16:16:29 | F11 amendment never consumed | `frozen.go` | killed | unit | TestAnAmendmentIsANewPlanThenANewPlanReview | fec9653ee4ec | True |
| 1 | 16:16:44 | F12 freezes at a handoff to any role | `frozen.go` | killed | unit | TestTheListFreezesOnlyAtAValidHandoffToTheVerifierThatGaveOne | fec9653ee4ec | True |
| 1 | 16:16:57 | F13 frozen files not committed | `session.go` | build failed: not run |  |  | 1164fa41417b | True |
| 1 | 16:16:58 | F14 frozen bytes not checked against the hash | `frozen.go` | killed | integration | TestAFrozenListWhoseHashNoHeadHoldsIsRefused | fec9653ee4ec | True |
| 1 | 16:17:26 | F15 a session with no bound head is frozen | `frozen.go` | SURVIVED |  |  | fec9653ee4ec | True |
| 1 | 16:17:54 | F16 list written outside Fenced | `handoffpost.go` | killed | integration | TestALostLeaseWritesNoFrozenList | 547fa32c28f9 | True |
| 1 | 16:18:23 | F17 a missing source is not a change | `frozen.go` | killed | unit | TestAChangedSourceIsTheFindingTestChanged | fec9653ee4ec | True |
| 1 | 16:18:37 | A1 every bound session is an author | `authors.go` | killed | unit | TestASessionBoundToACommitOutsideTheDiffIsNoAuthor, TestTheAuthorsOfAVerificationAreTheHarnessesBoundToACommitOfTheDiff | e28cf89e2a42 | True |
| 1 | 16:18:51 | A2 a commit of the base is an author | `authors.go` | build failed: not run |  |  | e28cf89e2a42 | True |
| 1 | 16:18:52 | A3 the writer of the frozen tests is no author | `authors.go` | killed | unit | TestTheWriterOfTheFrozenTestsIsAnAuthor | e28cf89e2a42 | True |
| 1 | 16:19:07 | A4 PairFor passes over no author | `authors.go` | killed | unit | TestNoHarnessOutsideTheAuthorsIsTheErrorOfThePair, TestThePairIsTheFirstAdmittedOneOutsideTheAuthors | e28cf89e2a42 | True |
| 1 | 16:19:22 | A5 no harness outside is not wrapped | `authors.go` | killed | integration | TestTheReviewersPairComesFromTheRoutingOfTheRecordsOutsideTheAuthors | e28cf89e2a42 | True |
| 1 | 16:19:53 | A6 plan authors of any role | `authors.go` | killed | unit | TestTheAuthorsOfAPlanReviewAreTheHarnessesOfThePlanSessions | e28cf89e2a42 | True |
| 1 | 16:20:08 | A7 plan authors of any task | `authors.go` | killed | unit | TestTheAuthorsOfAPlanReviewAreTheHarnessesOfThePlanSessions | e28cf89e2a42 | True |
| 1 | 16:20:23 | T1 not-active counted as a pass | `attempt.go` | killed | unit | TestAGateFailureIsTheRowsThatDidNotPass | 8d0ad6a85b9d | True |
| 1 | 16:20:38 | T2 an attempt with no failed row | `attempt.go` | killed | unit | TestNoFailureStartsNoAttempt | 8d0ad6a85b9d | True |
| 1 | 16:20:52 | R1 open questions not in the tier | `tier.go` | killed | unit | TestTheTierOfEachPlanClass | 4966b76508d7 | True |
| 1 | 16:21:08 | R2 new interfaces not in the tier | `tier.go` | killed | unit | TestTheTierOfEachPlanClass | 4966b76508d7 | True |
| 1 | 16:21:23 | R3 new dependencies not in the tier | `tier.go` | killed | unit | TestTheTierOfEachPlanClass | 4966b76508d7 | True |
| 1 | 16:21:37 | R4 only large is reasoning | `tier.go` | killed | unit | TestTheTierOfEachPlanClass | 4966b76508d7 | True |
| 1 | 16:21:51 | R5 the position ignores the unknowns | `tier.go` | killed | integration | TestThePlanAndTheRecordsGiveTheTierOfARetry | 4966b76508d7 | True |
| 1 | 16:22:22 | R6 by-class left to the row | `tier.go` | killed | unit | TestEveryRowOfTheStepTableResolvesToATier | 4966b76508d7 | True |
| 1 | 16:22:38 | R7 owner role not turned into the routing key | `tier.go` | killed | unit | TestOwnerOfItsKindIsResolvedByTheOwnerMap | 4966b76508d7 | True |
| 1 | 16:22:53 | R8 uphill leaves the plan tier | `tier.go` | killed | unit | TestTheProgressPositionMovesTheTier | 4966b76508d7 | True |
| 1 | 16:23:07 | R9 an accepted answer closes no question | `tier.go` | killed | unit | TestOpenQuestionsAreTheQuestionsNoAcceptedAnswerClosed | 4966b76508d7 | True |
| 1 | 16:23:22 | R10 downhill leaves the plan tier | `tier.go` | killed | unit | TestTheProgressPositionMovesTheTier | 4966b76508d7 | True |
| 1 | 16:23:38 | P1 kinds in the reverse order | `promptfile.go` | killed | unit | TestThePromptFileListsTheFiveKindsInTheirOrder | 8ab3d68544c7 | True |
| 1 | 16:23:56 | P2 the hash of another text | `promptfile.go` | build failed: not run |  |  | 8ab3d68544c7 | True |
| 1 | 16:23:57 | P3 the records copy is not in the input | `promptfile.go` | killed | unit | TestEachItemHasItsPathAndTheSHA256OfItsBytes | 8ab3d68544c7 | True |
| 1 | 16:24:15 | P4 a neighbour ID holds the ID | `promptfile.go` | build failed: not run |  |  | 8ab3d68544c7 | True |
| 1 | 16:24:16 | P5 a section runs to the end of the file | `promptfile.go` | killed | unit | TestADecisionRecordIsOneThatASectionNames | 8ab3d68544c7 | True |
| 1 | 16:24:30 | P6 every file is a rule file | `promptfile.go` | killed | unit | TestTheIssueAndPlanArePartOneAndAMissingPlanIsAnEmptyPartAndNotAnError, TestTheRuleFilesAreTheFilesOfTheRulePathRegister | 8ab3d68544c7 | True |
| 1 | 16:24:47 | P7 two records of one name are accepted | `promptfile.go` | killed | unit | TestTwoRecordsOfOneFileNameAreRefused | 8ab3d68544c7 | True |
| 1 | 16:25:10 | P8 prose is a PRD row | `promptfile.go` | killed | unit | TestTheRequirementItemsAreTheFilesWhoseRowOrHeadingHoldTheID | 8ab3d68544c7 | True |
| 1 | 16:25:33 | P9 the findings are left out of the text | `promptfile.go` | killed | unit | TestTheFindingsAreInThePromptAndThePromptHoldsNoDiff | 8ab3d68544c7 | True |
| 1 | 16:25:46 | P10 the input files are not made | `session.go` | killed | integration | TestAPromptFileIsBuiltFromTheRecordsOfARealTaskAndTheSessionReadsTheFiles | 1164fa41417b | True |
| 1 | 16:26:08 | P11 the prompt is not the context text | `session.go` | killed | integration | TestAGateFailureStartsTheNextAttemptFromTheBaseWithItsFindings, TestAPromptFileIsBuiltFromTheRecordsOfARealTaskAndTheSessionReadsTheFiles | 1164fa41417b | True |
| 1 | 16:26:35 | P12 a missing ADR is skipped | `promptfile.go` | build failed: not run |  |  | 8ab3d68544c7 | True |
| 1 | 16:26:40 | A8 session of another task is an author | `authors.go` | killed | unit | TestASessionOfAnotherTaskIsNoAuthor | e28cf89e2a42 | True |
| 1 | 16:26:50 | F18 the plan of a handoff to the verifier ignored To | `handoffpost.go` | SURVIVED |  |  | 547fa32c28f9 | True |
| 2 | 16:43:20 | F7 amendment: plan role not required for the plan | `frozen.go` | killed | unit | TestAnAmendmentIsNotAReviewWithoutAPlanOrInTheWrongOrder | fec9653ee4ec | True |
| 2 | 16:43:29 | F8 amendment: review role not required for the review | `frozen.go` | killed | unit | TestAnAmendmentIsNotAReviewWithoutAPlanOrInTheWrongOrder | fec9653ee4ec | True |
| 2 | 16:43:38 | F9 refreeze with no amendment | `frozen.go` | killed | integration | TestAChangedFrozenSourceIsAFindingAndTheHandoffIsInvalid, TestWithoutAnAmendmentAChangedTestIsNotFrozenAgain | fec9653ee4ec | True |
| 2 | 16:44:24 | F13 frozen files not committed | `session.go` | killed | integration | TestTheFrozenFilesAreInTheNextAttemptAtTheirHashes | 1164fa41417b | True |
| 2 | 16:44:49 | F15 a session with no bound head is frozen | `frozen.go` | killed | integration | TestASessionWithNoBoundHeadFreezesNothing | fec9653ee4ec | True |
| 2 | 16:45:13 | A2 a commit of the base is an author | `authors.go` | killed | integration | TestTheAuthorsAreReadFromRealCommitsAndRealRecords | e28cf89e2a42 | True |
| 2 | 16:45:50 | P2 the hash of another text | `promptfile.go` | killed | unit | TestEachItemHasItsPathAndTheSHA256OfItsBytes | 8ab3d68544c7 | True |
| 2 | 16:46:14 | P4 a neighbour ID holds the ID | `promptfile.go` | SURVIVED |  |  | 8ab3d68544c7 | True |
| 2 | 16:46:27 | P12 a missing ADR is skipped | `promptfile.go` | killed | unit | TestADecisionRecordIsOneThatASectionNames | 8ab3d68544c7 | True |
| 2 | 16:46:37 | F18 the plan of a handoff to the verifier ignored To | `handoffpost.go` | killed | unit | TestAListGivenToAnotherHandoffIsIgnored | 547fa32c28f9 | True |
| 3 | 16:47:14 | P4 a neighbour ID holds the ID | `promptfile.go` | killed | unit | TestTheRequirementItemsAreTheFilesWhoseRowOrHeadingHoldTheID | 8ab3d68544c7 | True |
| 4 | 16:55:58 | L1 the file itself may be a link | `session.go` | killed | integration | TestAFrozenFileThatIsALinkIsRefusedAndWritesNothingOutside, TestAFrozenFileThroughALinkInsideTheCheckoutIsRefused | 27bb8e3ad523 | True |
| 4 | 16:56:35 | L2 a parent may be a link | `session.go` | killed | integration | TestAFrozenFileThroughAParentLinkIsRefusedAndWritesNothingOutside, TestARefusedFrozenListWritesNoFile | 27bb8e3ad523 | True |
| 4 | 16:57:07 | L3 the check is not called | `session.go` | killed | integration | TestAFrozenFileThatIsALinkIsRefusedAndWritesNothingOutside, TestAFrozenFileThroughALinkInsideTheCheckoutIsRefused, TestAFrozenFileThroughAParentLinkIsRefusedAndWritesNothingOutside | 27bb8e3ad523 | True |
| 4 | 16:57:33 | L4 the walk follows links | `session.go` | killed | integration | TestAFrozenFileThatIsALinkIsRefusedAndWritesNothingOutside, TestAFrozenFileThroughALinkInsideTheCheckoutIsRefused, TestAFrozenFileThroughAParentLinkIsRefusedAndWritesNothingOutside | 27bb8e3ad523 | True |
| 4 | 16:58:00 | L5 the check comes with each write, not before all | `session.go` | killed | integration | TestARefusedFrozenListWritesNoFile | 27bb8e3ad523 | True |

Run 1 has 50 rows: 40 killed, 4 survived, 6 did not build (not run). Run 2 redid the 10 that needed it, with a valid variant of each of the six that did not build (F9, F13, A2, P2, P4, P12) and the survivors F7, F8, F15, F18, after the author added the tests that attack them (`TestAnAmendmentIsNotAReviewWithoutAPlanOrInTheWrongOrder` gained two cases, the no-head test asserts `ErrNoBoundHead`, `TestAListGivenToAnotherHandoffIsIgnored`, a base-is-the-head case in the authors test): 9 killed, P4 survived. Run 3 reran P4 after `TestTheRequirementItemsAreTheFilesWhoseRowOrHeadingHoldTheID` gained the longer-ID and the prefixed-ID cases: killed. Run 4 is the five mutants of the link repair: 5 killed. Final account: 55 distinct mutants, each killed by a named assertion, and none counted from a build failure.

The second-run P4 survivor is reported as it was: the first test of the neighbour ID checked only one direction.
