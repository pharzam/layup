# T-q0kr — the four statuses

Issue [#362](https://github.com/pharzam/layup/issues/362), row 66f of M2e. Serves REQ-007, REQ-003, NFR-004
(F-0003#47, F-0003#58). Plan 6248671869 (Opus 5.5); plan review 6248738516 (GPT-6.1 Sol, approve-with-conditions;
condition 1 became the size child T-cav3, #412, merged as 68a8a74; conditions 2 to 8 met). Maximum set by the
review: 2,200 lines over 26 files, cap 1. Landing figure: see the verdict. Evidence: [test runs and mutants](../../runs/T-q0kr/test-runs.md).

## Delivered and boundaries
`internal/run` computes and posts the four statuses (`GatesAt`, `SpecAt`, `RulesAt`, `VerifyAt`, `PostStatuses`;
`PostStatus` is the one rule: post only when the App's newest status differs in state or description). The tables
`gates.tsv` (block `task-gates`, now built) and `spec-check.tsv` are written once before the post and checked
for completeness before a status is made from them. A `target:` source error is a posted `error`, a `records:`
one a Go error; the verifier is checked for independence and its status is posted on the head it reviewed.
Not here: the step `loop` (66k), ready and merge (66g), the bet's copy (66j). Decisions are in `loop.md`,
`records.md` and `packages.md`.

## Tests mapped to the criteria

| Criterion | Test |
| --- | --- |
| Gates by row; `not-active` and a changed product path fail; saved tables, completeness, the block | `TestEachRowCombinationMapsToTheGatesStatus`, `TestAPendingKindWithAProductPathAndANotActiveKindAreFailures`, `TestAGatesTableIsReadByItsRulesAndRefusesEachBrokenOne`, `TestASavedTableIsCompleteOnlyWhenItsKindsAreTheManifestOfItsBase`, `TestASpecTableIsReadByItsRules`, `TestTheGatesBlockEqualsItsGoSchema` |
| Post once, newest App status, other sources, lease, failures, fit, the bot | `TestAStatusIsPostedOnceForAnUnchangedStateAndAgainWhenItChanges`, `TestOnlyTheNewestStatusOfTheAppCountsAndOtherSourcesNeverDo`, `TestALostLeasePostsNoStatus`, `TestADescriptionIsCutAtOneHundredFortyCharacters`, `TestTheSessionsOfTheRunPostAsTheApp` |
| Record before post; recovery; output and progress; cleanup error | `TestAFailedPostIsRecoveredFromTheSavedTableWithoutRunningTheGateAgain`, `TestTheDefaultProgressOfAKindBeatsAndStopsBeforeTheOutput`, `TestACleanupErrorOfTheGatesStopsNoPostAndIsReturnedByPostStatuses`, `TestTheFourStatusesOfAHead` (the demo), `TestAFailedPostOfOneContextLeavesTheEarlierPostsAndTheNextCallPostsTheRest` |
| Spec, rules and verify statuses and their sources | `TestTheSpecStatusIsSuccessOnlyWithNoFailRowAndACompleteTable`, `TestATargetSourceOfTheFullCheckIsAPostedErrorAndARecordsSourceIsAGoError`, `TestTheRegisterCommitOfAHeadIsTheRecordsOfItsBoundSessionOrTheCallers`, `TestVerifyIsSuccessOnlyForNothingMaterialAtTheReviewedHeadWithAResultThatIsDone`, `TestAVerifierThatIsAnAuthorOrNotAdmittedIsNotIndependent` |

## Alternatives
| Decision | Selected | Rejected, and why |
| --- | --- | --- |
| Whether a status was posted | the forge's newest App status | a `statuses.tsv` record (trusts our own post); any creator |
| Register of `layup/rules` | the bound session's `records` column, or the caller's full SHA | the head's register, the run's changing base |
| The verifier | the caller names it; code checks task, role, harness, model | the records alone (plan review and verifier share a role) |

## Verification and verdict
Test-first (stubs alone in `e17c09d`, `1cc2b3f`, `8e85e68`; red runs recorded) with 38 compiled mutants, all
killed. After the fix `2b6182c`: build, vets, `go test ./...`, the 13 named integration tests and the five text
checks pass; the full integration level passed on `c546014` (`internal/run` alone after a timeout under load).

Round 1 (GPT-6.1 Sol, xhigh, Codex): **material**. F1 an absent result of a saved spec row gave success; F2 no
ten-second heartbeat in the default progress; F3 a malformed result file lost its Go error; F4 `SpecAt` read saved
rows before checking its SHAs; F5 three caller checks and the production constructor untested; F6 the budget.
F1 to F5 fixed test first in `2b6182c`. Round 2 (same reviewer, comment 6251314024): **not mergeable, findings
recorded**, F6 the only finding. Harness exception: ADR-0026 names `claude`, `devin` and `opencode`; Codex was used
on the Operator's message of 2026-10-13. Landing figure: 2,534 lines (2,474 added, 60 removed) over 23 files, before
this task file and the two indexes, against a maximum of 2,200/26; the merge waits for the Operator's one approval.

## Resource record

ADR-0007, recorded, not budgeted.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 | not recorded | not separately recorded | not separately recorded |
| Plan review | reasoning | GPT-6.1 Sol | xhigh | 201,460 | 2026-10-18T10:18:40Z to 10:25:42Z |
| Implementation, phase A | execution | Claude Sonnet 5.5 | not recorded | not recorded | about 75 min, ending 11:26Z |
| Implementation, phase B | execution | Claude Sonnet 5.5 | not recorded | not recorded | about 110 min, ending 14:12Z |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 218,219 | 14:12:54Z to 14:26:47Z |
| Fix of round 1 | execution | Claude Sonnet 5.5 | not recorded | not recorded | about 15 min, ending 14:41Z |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 206,519 | 14:42:23Z to 14:51:22Z |
| Close-out | execution | Claude Sonnet 5.5 | not recorded | not recorded | about 5 min, ending 14:56Z |
| Orchestration | reasoning | Claude Opus 5.5 | not recorded | not separately recorded | not separately recorded |
