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

## Phase A, full integration level

`TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...` on `9d66a70`: exit 0, 11:20:30Z to 11:51:29Z (`internal/run` 1740.9 s).

## Phase B, after the merge of main (68a8a74, T-cav3), red against the inert stubs (8e85e68)

SpecAt and PostStatuses were stubs. The first red run (12:46:36Z) failed in the setup of the two demo tests (`tasks.tsv` replaced by the records of Intake before the comment was posted); the setup was fixed and the run below is the red one. Exit 1, four failing tests, every line an assertion.

```
--- FAIL: TestTheSpecStatusOfAHeadThatTouchesNoSpecificationIsClearAndIsRebuiltFromItsTable (2.50s)
    statuses_integration_test.go:257: SpecAt = {Context: State: Description:}, <nil>; want {Context:layup/spec State:success Description:clear: no specification yet}
--- FAIL: TestATargetSourceOfTheFullCheckIsAPostedErrorAndARecordsSourceIsAGoError (2.48s)
    statuses_integration_test.go:286: a target source: {Context: State: Description:}, <nil>; want error, error: PSBRecords: …
    statuses_integration_test.go:299: a records source: {Context: State: Description:}, <nil>; want a Go SourceError of records: and no status
    statuses_integration_test.go:305: a base that is no full SHA: no error
--- FAIL: TestTheFourStatusesOfAHead (4.04s)
    statuses_integration_test.go:333: PostStatuses = [], <nil>; want [{SHA:aca1b097a0786210b5a6458b2ff84050b35c2f7a Status:{Context:layup/gates State:failure Description:lint: not-active: to
--- FAIL: TestAFailedPostOfOneContextLeavesTheEarlierPostsAndTheNextCallPostsTheRest (4.28s)
    statuses_integration_test.go:361: PostStatuses = [], <nil>; want the gates post made and an error at the spec post
```

Green: the same four tests, then the whole `internal/run` integration set, passed after the code (below, the full run).

## Compiled mutants (29)

Run by `mutants.py` in the session scratchpad, which, for each mutant, replaces one exact text once, runs `go build ./...` and `go vet -tags=integration .`, then the unit level or the integration tests of this task, and restores the file; the hash is the first 8 hexadecimal characters of the SHA-256 of the file before, mutated, and after the restore (the first and the last equal). Run 12:49:20Z to 13:08:43Z. All 29 are killed and none is claimed equivalent. M20 is killed by a panic of the test (an index out of range), not an assertion; M8 and M7 are killed by the second call of `GatesAt` failing, which shows the record is read.

| # | Mutant | File | Build | Vet | Test | First failing assertion | SHA-256 before, mutated, after |
| - | ------ | ---- | ----- | --- | ---- | ----------------------- | ------------------------------ |
| M1 | clear counts as a failure | gates.go | exit 0 | exit 0 | exit 1 (killed) | TestEachRowCombinationMapsToTheGatesStatus: gates_test.go:46: a clear kind counts as a pass, with its reason: GatesStatus = {Context:layup/gates State:fai | d8e887fe>44e48c7d>d8e887fe |
| M2 | not-active counts as a pass | gates.go | exit 0 | exit 0 | exit 1 (killed) | TestEachRowCombinationMapsToTheGatesStatus: gates_test.go:46: a not-active kind is a failure, never a success: GatesStatus = {Context:layup/gates State:su | d8e887fe>20aa83aa>d8e887fe |
| M3 | the compare ignores the description | statuses.go | exit 0 | exit 0 | exit 1 (killed) | TestAStatusIsPostedOnceForAnUnchangedStateAndAgainWhenItChanges: statuses_test.go:89: {Context:layup/gates State:success Description:clear: no product path} after a different  | 94f70cde>a785aa2e>94f70cde |
| M4 | any creator counts | statuses.go | exit 0 | exit 0 | exit 1 (killed) | TestOnlyTheNewestStatusOfTheAppCountsAndOtherSourcesNeverDo: statuses_test.go:61: appStatus = {Context:layup/gates State:success Description:pass CreatorID:7 CreatorLogin: | 94f70cde>1520acdc>94f70cde |
| M5 | the oldest App status is compared | statuses.go | exit 0 | exit 0 | exit 1 (killed) | TestOnlyTheNewestStatusOfTheAppCountsAndOtherSourcesNeverDo: statuses_test.go:61: appStatus = {Context:layup/gates State:success Description:old CreatorID:9001 CreatorLogi | 94f70cde>9090e169>94f70cde |
| M6 | the lease read is skipped | statuses.go | exit 0 | exit 0 | exit 1 (killed) | TestALostLeasePostsNoStatus: statuses_test.go:154: a lost lease: true, <nil>, 1 statuses; want a LostError and nothing posted | 94f70cde>dd00bad7>94f70cde |
| M7 | the gate runs again for a recorded head | gates.go | exit 0 | exit 0 | exit 1 (killed) | TestTheGatesStatusIsPostedOnceFromTheRecordedTableAndItsOutputIsKept: statuses_integration_test.go:100: the second GatesAt = {Context: State: Description:}, line 4: the key (head " | d8e887fe>4d80e46b>d8e887fe |
| M8 | the table is not committed before the status | gates.go | exit 0 | exit 0 | exit 1 (killed) | TestTheGatesStatusIsPostedOnceFromTheRecordedTableAndItsOutputIsKept: statuses_integration_test.go:100: the second GatesAt = {Context:layup/gates State:success Description:clear: p | d8e887fe>624ed985>d8e887fe |
| M9 | the spec table may hold the rules in any order | specstatus.go | exit 0 | exit 0 | exit 1 (killed) | TestTheSpecStatusIsSuccessOnlyWithNoFailRowAndACompleteTable: gates_test.go:170: the rules in the wrong order: SpecStatus = {Context:layup/spec State:success Description:pa | bfb14ac7>12236203>bfb14ac7 |
| M10 | no harness outside the authors gives pending | verifystatus.go | exit 0 | exit 0 | exit 1 (killed) | TestNoHarnessOutsideTheAuthorsIsFailureNotActive: verifystatus_test.go:91: VerifyAt = aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa, {Context:layup/verify State:pend | 9dd50b17>8f16f952>9dd50b17 |
| M11 | success without the reviewed head | verifystatus.go | exit 0 | exit 0 | exit 1 (killed) | TestVerifyIsSuccessOnlyForNothingMaterialAtTheReviewedHeadWithAResultThatIsDone: verifystatus_test.go:148: a record of another head: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa, {Context:layup/v | 9dd50b17>07122d48>9dd50b17 |
| M12 | the count of findings is the count of bullets plus one | verifystatus.go | exit 0 | exit 0 | exit 1 (killed) | TestVerifyIsSuccessOnlyForNothingMaterialAtTheReviewedHeadWithAResultThatIsDone: verifystatus_test.go:148: a finding: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa, {Context:layup/verify State:fai | 9dd50b17>622bc842>9dd50b17 |
| M13 | the register is read at the head | statuses.go | exit 0 | exit 0 | exit 1 (killed) | TestTheRulesStatusReadsTheRegisterOfTheRecordsCommit: statuses_integration_test.go:190: RulesAt = {Context:layup/rules State:error Description:register}, <nil>; wan | 94f70cde>1d071add>94f70cde |
| M14 | fit cuts at 141 | statuses.go | exit 0 | exit 0 | exit 1 (killed) | TestADescriptionIsCutAtOneHundredFortyCharacters: statuses_test.go:48: 141 are cut to 139 and an ellipsis: fit = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa | 94f70cde>aeb856e2>94f70cde |
| M15 | ReadGates accepts rows of a head that are apart | gates.go | exit 0 | exit 0 | exit 1 (killed) | TestAGatesTableIsReadByItsRulesAndRefusesEachBrokenOne: gates_test.go:71: two bases in a head: ReadGates accepted the table | d8e887fe>5b21ddd3>d8e887fe |
| M16 | the author check is dropped | verifystatus.go | exit 0 | exit 0 | exit 1 (killed) | TestAVerifierThatIsAnAuthorOrNotAdmittedIsNotIndependent: verifystatus_test.go:158: a clean record of an author: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa, {Context:layu | 9dd50b17>8c406301>9dd50b17 |
| M17 | the admission check is dropped | verifystatus.go | exit 0 | exit 0 | exit 1 (killed) | TestAVerifierThatIsAnAuthorOrNotAdmittedIsNotIndependent: verifystatus_test.go:176: a model of use no: {Context:layup/verify State:success Description:nothing material  | 9dd50b17>f695219c>9dd50b17 |
| M18 | the constructor forgets the bot | meaning.go | exit 0 | exit 0 | exit 1 (killed) | TestTheSessionsOfTheRunKnowTheBotOfTheApp: statuses_test.go:211: state.sessions().Bot = 0, want 9001 (the ID that step forge read) | 54dbdeed>46d35da2>54dbdeed |
| M19 | the restart forgets the bot | handoffpost_integration_test.go | exit 0 | exit 0 | exit 1 (killed) | TestAFailedPostIsRecoveredFromTheSavedTableWithoutRunningTheGateAgain: statuses_integration_test.go:134: the Sessions of a restart have the bot 0; want 9001 | fadb394d>683b7868>fadb394d |
| M20 | a gate table with a kind omitted is complete | gates.go | exit 0 | exit 0 | exit 1 (killed) | TestASavedTableIsCompleteOnlyWhenItsKindsAreTheManifestOfItsBase: panic: runtime error: index out of range [0] with length 0 [recovered, repanicked] | d8e887fe>135529de>d8e887fe |
| M21 | a pass may carry a reason | gates.go | exit 0 | exit 0 | exit 1 (killed) | TestAGatesTableIsReadByItsRulesAndRefusesEachBrokenOne: gates_test.go:71: a pass with a reason: ReadGates accepted the table | d8e887fe>6c8fdbc2>d8e887fe |
| M22 | a verifier of another task is accepted | verifystatus.go | exit 0 | exit 0 | exit 1 (killed) | TestAVerifierOfAnotherTaskOrRoleOrNoneIsAnError: verifystatus_test.go:197: a session of another task: 0123456789abcdef0123456789abcdef01234567, {Context:layup/ | 9dd50b17>c842fcf3>9dd50b17 |
| M23 | the fail row is not the first in the order | specstatus.go | exit 0 | exit 0 | exit 1 (killed) | TestTheSpecStatusIsSuccessOnlyWithNoFailRowAndACompleteTable: gates_test.go:170: a fail past a clear row: SpecStatus = {Context:layup/spec State:success Description:clear:  | bfb14ac7>66d7a084>bfb14ac7 |
| M24 | SpecAt checks again for a recorded head | specstatus.go | exit 0 | exit 0 | exit 1 (killed) | TestTheSpecStatusOfAHeadThatTouchesNoSpecificationIsClearAndIsRebuiltFromItsTable: statuses_integration_test.go:276: the second SpecAt = {Context: State: Description:}, the input of the full sp | bfb14ac7>6d898057>bfb14ac7 |
| M25 | a target source is a Go error, not a posted error | specstatus.go | exit 0 | exit 0 | exit 1 (killed) | TestATargetSourceOfTheFullCheckIsAPostedErrorAndARecordsSourceIsAGoError: statuses_integration_test.go:294: a target source: {Context: State: Description:}, the input of the full spec  | bfb14ac7>873877c4>bfb14ac7 |
| M26 | a records source is a posted error | specstatus.go | exit 0 | exit 0 | exit 1 (killed) | TestATargetSourceOfTheFullCheckIsAPostedErrorAndARecordsSourceIsAGoError: statuses_integration_test.go:307: a records source: {Context:layup/spec State:error Description:error: Questio | bfb14ac7>994a86d1>bfb14ac7 |
| M27 | the table of a second head replaces the first | specstatus.go | exit 0 | exit 0 | exit 1 (killed) | TestTheSpecStatusOfAHeadThatTouchesNoSpecificationIsClearAndIsRebuiltFromItsTable: statuses_integration_test.go:269: the first head after the second: 0 rows, false, <nil>; want its ten rows kep | bfb14ac7>e55a06c5>bfb14ac7 |
| M28 | the verify status is posted on the head, not the reviewed head | statuses.go | exit 0 | exit 0 | exit 1 (killed) | TestTheFourStatusesOfAHead: statuses_integration_test.go:341: PostStatuses = [{SHA:cacb6a6e08a8b3a3556617e8983fe8a4cd6cc355 Status:{Contex | 94f70cde>07c032d7>94f70cde |
| M29 | a failed post is swallowed | statuses.go | exit 0 | exit 0 | exit 1 (killed) | TestAFailedPostOfOneContextLeavesTheEarlierPostsAndTheNextCallPostsTheRest: statuses_integration_test.go:369: PostStatuses = [{SHA:d9327a3ac8a5fe2cd63826f933ffc21d6c7b8d07 Status:{Contex | 94f70cde>334c5afc>94f70cde |

## Phase B, gates

Commands exactly as run, each with its exit status read directly (all in the worktree):

```
go build ./...                                   exit 0
go vet ./...                                     exit 0
go vet -tags=integration ./...                   exit 0
go vet -tags=e2e ./cmd/layup/                    exit 0
gofmt -l .                                       no output
go test -count=1 ./...                           every package ok
TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...
```

The full integration level, 13:09:23Z to 13:41:08Z: every package ok but `internal/run`, which hit the 30 minute limit (`panic: test timed out after 30m0s`, in a world set-up of the new tests; the load average of the host was about 6, and the run took 1740.9 s in phase A). That is the host's timing, not a failure of an assertion, so `internal/run` was run alone with a longer limit:

```
TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 45m ./internal/run/
```

13:41:34Z to 14:10:53Z, `ok  github.com/pharzam/layup/internal/run 1756.847s`, exit 0. The other 21 packages passed in the first run.

The five text checks (`sh docs/adr/adr-lint.sh`, `sh docs/prd/prd-lint.sh`, `sh docs/links/link-lint.sh`, `sh docs/tests/run-discipline-tests.sh`, `sh docs/setup/setup-check.sh`), `sh docs/tests/nested-checkout-check.sh` and `git diff --check origin/main HEAD`: exit 0 each (the last after removing the trailing spaces of the red log).
