# T-08rj - the test runs of Step A (the forge call)

(Trailing blanks of the blank lines inside the quoted output are stripped; nothing else is changed.)

Every command was run directly in `.worktree/T-08rj`, never piped; the exit status is the one `go` gave
and the times are `date -u` read around the run. The red runs are against the compiling stubs of commit
`32118c3` (`Body` declared and set by nothing, `EditBody` returning the zero pull request with no request),
plus the inert field `Fake.FailAfter`, which the red tests need to compile.

## Red, the stand-in (F1, F2, F3)

Command: `go test -count=1 ./internal/forge/`

```
start 2026-10-18T13:09:15Z
--- FAIL: TestFakeKeepsTheBodyOfADraftAndEditsIt (0.00s)
    forge_test.go:309: OpenDraft = {Number:1 NodeID:PR_fake1 State:open Draft:true Merged:false HeadSHA:1111111111111111111111111111111111111111 MergeSHA: Body:}, <nil>, want the body it was given
--- FAIL: TestFakeRefusesAnEditBelow1BeforeItsFailureAndGivesTheForgeErrors (0.00s)
    forge_test.go:361: EditBody(0) with Fail = <nil>; want a refusal of the call EditBody with the status 0 (a 404 or a 502 is not one)
    forge_test.go:361: EditBody(-1) with Fail = <nil>; want a refusal of the call EditBody with the status 0 (a 404 or a 502 is not one)
    forge_test.go:367: EditBody(99) = <nil>; want a forge error of status 404
    forge_test.go:372: EditBody(1) with Fail = <nil>; want the forge error of status 502
--- FAIL: TestFakeCanTakeAnEditAndStillGiveAnError (0.00s)
    forge_test.go:387: EditBody with FailAfter = <nil>; want the forge error of status 502
    forge_test.go:390: after the error: Body = ""; want "new" (the edit took effect)
    forge_test.go:394: a refusal below 1 gave no error
FAIL
FAIL	github.com/pharzam/layup/internal/forge	0.059s
FAIL
exit=1
end 2026-10-18T13:09:17Z
```

F1 stops at its first assertion (a `t.Fatalf`), so its later lines are not reached at the stub. F3 is the
case of the review's condition 4, an edit that takes effect but returns an error (`Fake.FailAfter`).

## Red, the adapter (A1 to A5)

Command: `go test -count=1 -tags=integration ./internal/forge/github/`

A5 primes no token: the first request makes the first token and the 401 makes the second, so a stub
makes zero tokens (condition 4). Green negatives at the stub, each with its red control: the "0 requests"
part of A3 (control: its refusal line), the `body: null` rows of A4 (control: the by-number row), and the
unchanged `prJSON` comparisons of the existing tests.

```
start 2026-10-18T13:09:20Z
--- FAIL: TestEditBodyPatchesTheBodyAloneAndGivesThePullRequest (0.00s)
    github_integration_test.go:1275: the request body = map[]; want map[body:Closes #7

        The records commit: abc] (the body alone: the title, the state and the base cannot change)
    github_integration_test.go:1278: PATCH /repos/acme/target/pulls/5 requests = 0; want 1
    github_integration_test.go:1281: EditBody = {Number:0 NodeID: State: Draft:false Merged:false HeadSHA: MergeSHA: Body:}; want {Number:5 NodeID:PR_node5 State:open Draft:false Merged:false HeadSHA:1111111111111111111111111111111111111111 MergeSHA: Body:Closes #7

        The records commit: abc} (the answer is decoded, with its body)
--- FAIL: TestEditBodyGivesAForgeErrorForAnyOtherStatus (0.00s)
    github_integration_test.go:1295: EditBody with 422 = <nil>; want the call EditBody, the status 422 and the first line
    github_integration_test.go:1295: EditBody with 403 = <nil>; want the call EditBody, the status 403 and the first line
    github_integration_test.go:1295: EditBody with 404 = <nil>; want the call EditBody, the status 404 and the first line
--- FAIL: TestEditBodyRefusesANumberBelow1BeforeAnyRequest (0.00s)
    github_integration_test.go:1309: EditBody(0) = <nil>; want a refusal of the call EditBody with the status 0
    github_integration_test.go:1309: EditBody(-1) = <nil>; want a refusal of the call EditBody with the status 0
--- FAIL: TestPullRequestGivesTheBodyByNumberAndByBranchAndANullBodyIsEmpty (0.02s)
    github_integration_test.go:1335: by number: Body = "", <nil>; want "Closes #7\n\nThe records commit: abc"
    github_integration_test.go:1338: by branch (the list): Body = "", <nil>; want "Closes #7\n\nThe records commit: abc"
--- FAIL: TestEditBodySharesTheHandlingOfA401AndOfARateLimit (0.00s)
    github_integration_test.go:1356: EditBody after a 401 = {Number:0 NodeID: State: Draft:false Merged:false HeadSHA: MergeSHA: Body:}, <nil>; want one new token and one retry
    github_integration_test.go:1359: PATCH requests = 0; want 2 (one 401, one retry)
    github_integration_test.go:1362: tokens made = 0; want 2 (the first request, and the new one after the 401)
    github_integration_test.go:1370: the wait: sleeps [], lines []; want three waits of at most ten seconds, a line of the call each
    github_integration_test.go:1375: a 500 gave no error
FAIL
FAIL	github.com/pharzam/layup/internal/forge/github	1.220s
FAIL
exit=1
end 2026-10-18T13:09:23Z
```

## Green

Command: `go test -count=1 ./internal/forge/`

```
start 2026-10-18T13:10:04Z
ok  	github.com/pharzam/layup/internal/forge	0.452s
exit=0
end 2026-10-18T13:10:06Z
```

Command: `go test -count=1 -tags=integration ./internal/forge/github/`

```
start 2026-10-18T13:10:06Z
ok  	github.com/pharzam/layup/internal/forge/github	2.684s
exit=0
end 2026-10-18T13:10:11Z
```

# Step B - the link of the review record (row 66q)

The first record of Step B (2026-10-19, 06:02 to 07:08Z) wrote placeholders in place of its test selectors, so its commands cannot
be reproduced; round 1 (finding 2) found it. It is replaced, not corrected: the runs below are new, with their own `date -u`
times, taken after the fix of finding 1 and the merge of `origin/main` at `d9852d2`. None of the older times is reused.
The selectors are exact, and the exit status is the one `go` gave.

- Unit selector U: `^(TestTheBodyWithAReviewRecordEndsInTheLinkOfItsComment|TestTheBodyIsEditedOnceToLinkTheReviewRecord|TestAnEqualBodyIsNotEdited|TestTheRecordsLinkIsTheProducingSessionsNotTheCloseOutsNorTheVerifiers|TestTheEditKeepsTheSnapshotOfTheLastBoundSessionBeforeTheVerifier|TestEachRecordThatFailsRefusesTheLinkBeforeAnyCall|TestALostLeaseEditsNothing|TestAClosedOrMergedPullRequestIsNotEdited|TestAFailedEditIsAnErrorAndAnEditThatTookEffectIsNotRepeated|TestAFailedReadOfThePullRequestEditsNothing|TestAFailedStoreReadIsAnErrorWithNoForgeCall|TestAnAnswerWithAnotherBodyIsAnError|TestANewerAttemptWithTheRegisterOnTheOlderPullRequestEditsNothing|TestAPullRequestWhoseHeadMovedIsNotEdited)$`
- Integration selector I: `^(TestTheReviewRecordIsLinkedFromTheBodyAfterItsPost|TestTheLinkSurvivesTheCloseOutAndARestart|TestAnEditRefusedByTheForgeIsAForgeErrorOfEditBody)$`

Differences from the plan, found at S6 (the reconciliation against the landed 66g, PR #421, merge `c735788`):

- `PullBody(issue, web, repo, records, more ...string)` already takes further lines, so there is no `review` parameter: the link is
  `more = "", reviewLine(...)`.
- `Task.PR` is written by `registerPR`, `holdsLease` is in `handoffpost.go`, `Sessions.Web` is set by `state.sessions()`
  (`internal/run/meaning.go`); row 66g already excludes the body edit, and row 66q and the After cells of 66h and 66k are on main.
- The snapshot is the one of the opening: the last bound session before the verifier's own event `session`, through the
  `lastBound` that `OpenPullRequest` uses (round 1, finding 1; the plan's "first session" and the review's condition 1 both
  gave way to it). The latest handoff is chosen and then required valid; the pull request is read by the attempt's branch.

## Finding 1, red on the code before the fix

Command: `go test -count=1 -run '^TestTheEditKeepsTheSnapshotOfTheLastBoundSessionBeforeTheVerifier$' ./internal/run/` on the code of `363bbfc` (the first-bound
selection, the head reviewed in round 1) with the new test added to the tree:

```
start 2026-10-19T07:56:46Z
--- FAIL: TestTheEditKeepsTheSnapshotOfTheLastBoundSessionBeforeTheVerifier (0.00s)
    reviewlink_test.go:267: LinkReviewRecord = true, <nil> with the body "Closes #3\n\nThe records commit this task started from: https://github.example/acme/target/commit/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n\nThe review record: https://github.example/acme/target/issues/3#issuecomment-601"; wa
FAIL
FAIL	github.com/pharzam/layup/internal/run	0.011s
FAIL
exit=1
end 2026-10-19T07:57:01Z
```

(The second part of that test, a verifier with no event `session`, was added to the same test in the same step.)

## Red at the stub, the unit tests

The stub is the `reviewlink.go` of commit `f30ba1c` (`LinkReviewRecord` returns false, nil; `reviewLine` returns ""), put in place
by `-overlay` so that the tree is not touched (`$SP/stub/overlay.json` maps `internal/run/reviewlink.go` to `git show f30ba1c:internal/run/reviewlink.go`).

Command: `go test -count=1 -overlay $SP/stub/overlay.json -run '^(TestTheBodyWithAReviewRecordEndsInTheLinkOfItsComment|TestTheBodyIsEditedOnceToLinkTheReviewRecord|TestAnEqualBodyIsNotEdited|TestTheRecordsLinkIsTheProducingSessionsNotTheCloseOutsNorTheVerifiers|TestTheEditKeepsTheSnapshotOfTheLastBoundSessionBeforeTheVerifier|TestEachRecordThatFailsRefusesTheLinkBeforeAnyCall|TestALostLeaseEditsNothing|TestAClosedOrMergedPullRequestIsNotEdited|TestAFailedEditIsAnErrorAndAnEditThatTookEffectIsNotRepeated|TestAFailedReadOfThePullRequestEditsNothing|TestAFailedStoreReadIsAnErrorWithNoForgeCall|TestAnAnswerWithAnotherBodyIsAnError|TestANewerAttemptWithTheRegisterOnTheOlderPullRequestEditsNothing|TestAPullRequestWhoseHeadMovedIsNotEdited)$' ./internal/run/` (U above)

```
start 2026-10-19T08:00:55Z
--- FAIL: TestTheBodyWithAReviewRecordEndsInTheLinkOfItsComment
    reviewlink_test.go:182: the body = "Closes #3\n\nThe records commit this task started from: https://github.example/acme/target/commit/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n\n"; want "Clo
--- FAIL: TestTheBodyIsEditedOnceToLinkTheReviewRecord
    reviewlink_test.go:194: LinkReviewRecord = false, <nil>; want true, nil
--- FAIL: TestAnEqualBodyIsNotEdited
    reviewlink_test.go:216: LinkReviewRecord = false, <nil> with the calls []; want false, nil and one PullRequest
--- FAIL: TestTheRecordsLinkIsTheProducingSessionsNotTheCloseOutsNorTheVerifiers
    reviewlink_test.go:243: after a person's edit: false, <nil>, body "a person's text"; want the whole body replaced
--- FAIL: TestTheEditKeepsTheSnapshotOfTheLastBoundSessionBeforeTheVerifier
    reviewlink_test.go:268: LinkReviewRecord = false, <nil> with the body "Closes #3\n\nThe records commit this task started from: https://github.example/acme/target/commit/aaaaaaaaaaaaaaaaaaaaa
--- FAIL: TestEachRecordThatFailsRefusesTheLinkBeforeAnyCall
    reviewlink_test.go:333: the session is of another task: err = <nil>; want one that says "not of T-ab12"
--- FAIL: TestALostLeaseEditsNothing
    reviewlink_test.go:352: LinkReviewRecord = false, <nil> with the calls []; want *LostError and no EditBody
--- FAIL: TestAClosedOrMergedPullRequestIsNotEdited
    reviewlink_test.go:363: a closed pull request: err = <nil>; want one that says "is closed"
--- FAIL: TestAFailedEditIsAnErrorAndAnEditThatTookEffectIsNotRepeated
    reviewlink_test.go:385: a failed edit: err = <nil>; want one that says "EditBody: status 502"
--- FAIL: TestAFailedReadOfThePullRequestEditsNothing
    reviewlink_test.go:412: a failed read: err = <nil>; want one that says "PullRequest: status 502"
--- FAIL: TestAFailedStoreReadIsAnErrorWithNoForgeCall
    reviewlink_test.go:425: a failed read of tasks.tsv: false, <nil> with the calls []; want the store's error and no call
--- FAIL: TestAnAnswerWithAnotherBodyIsAnError
    reviewlink_test.go:445: another body: err = <nil>; want one that says "gave another body"
--- FAIL: TestANewerAttemptWithTheRegisterOnTheOlderPullRequestEditsNothing
    reviewlink_test.go:477: no pull request of the branch: err = <nil>; want one that says "no pull request"
--- FAIL: TestAPullRequestWhoseHeadMovedIsNotEdited
    reviewlink_test.go:511: a moved head: err = <nil>; want one that says "not the reviewed"
FAIL
FAIL	github.com/pharzam/layup/internal/run	0.031s
FAIL
exit=1
end 2026-10-19T08:01:06Z
```

14 of 14 fail; the first failing line of each is kept.

## Red at the stub, the integration tests

Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -overlay $SP/stub/overlay.json -run '^(TestTheReviewRecordIsLinkedFromTheBodyAfterItsPost|TestTheLinkSurvivesTheCloseOutAndARestart|TestAnEditRefusedByTheForgeIsAForgeErrorOfEditBody)$' -timeout 30m ./internal/run/` (I above)

```
start 2026-10-19T08:01:06Z
--- FAIL: TestTheReviewRecordIsLinkedFromTheBodyAfterItsPost
    reviewlink_integration_test.go:94: LinkReviewRecord = false, <nil>; want true, nil
--- FAIL: TestTheLinkSurvivesTheCloseOutAndARestart
    reviewlink_integration_test.go:133: after the close-out and a restart: false, <nil> with 0 PATCH; want false, nil and the one
--- FAIL: TestAnEditRefusedByTheForgeIsAForgeErrorOfEditBody
    reviewlink_integration_test.go:147: LinkReviewRecord = false, <nil>; want EditBody: status 403
FAIL
FAIL	github.com/pharzam/layup/internal/run	13.933s
FAIL
exit=1
end 2026-10-19T08:01:37Z
```

## Green

Command: `go test -count=1 -v -run '^(TestTheBodyWithAReviewRecordEndsInTheLinkOfItsComment|TestTheBodyIsEditedOnceToLinkTheReviewRecord|TestAnEqualBodyIsNotEdited|TestTheRecordsLinkIsTheProducingSessionsNotTheCloseOutsNorTheVerifiers|TestTheEditKeepsTheSnapshotOfTheLastBoundSessionBeforeTheVerifier|TestEachRecordThatFailsRefusesTheLinkBeforeAnyCall|TestALostLeaseEditsNothing|TestAClosedOrMergedPullRequestIsNotEdited|TestAFailedEditIsAnErrorAndAnEditThatTookEffectIsNotRepeated|TestAFailedReadOfThePullRequestEditsNothing|TestAFailedStoreReadIsAnErrorWithNoForgeCall|TestAnAnswerWithAnotherBodyIsAnError|TestANewerAttemptWithTheRegisterOnTheOlderPullRequestEditsNothing|TestAPullRequestWhoseHeadMovedIsNotEdited)$' ./internal/run/`

```
start 2026-10-19T08:01:37Z
--- PASS: TestTheBodyWithAReviewRecordEndsInTheLinkOfItsComment (0.00s)
--- PASS: TestTheBodyIsEditedOnceToLinkTheReviewRecord (0.00s)
--- PASS: TestAnEqualBodyIsNotEdited (0.00s)
--- PASS: TestTheRecordsLinkIsTheProducingSessionsNotTheCloseOutsNorTheVerifiers (0.00s)
--- PASS: TestTheEditKeepsTheSnapshotOfTheLastBoundSessionBeforeTheVerifier (0.00s)
--- PASS: TestEachRecordThatFailsRefusesTheLinkBeforeAnyCall (0.02s)
--- PASS: TestALostLeaseEditsNothing (0.00s)
--- PASS: TestAClosedOrMergedPullRequestIsNotEdited (0.00s)
--- PASS: TestAFailedEditIsAnErrorAndAnEditThatTookEffectIsNotRepeated (0.00s)
--- PASS: TestAFailedReadOfThePullRequestEditsNothing (0.00s)
--- PASS: TestAFailedStoreReadIsAnErrorWithNoForgeCall (0.00s)
--- PASS: TestAnAnswerWithAnotherBodyIsAnError (0.00s)
--- PASS: TestANewerAttemptWithTheRegisterOnTheOlderPullRequestEditsNothing (0.00s)
--- PASS: TestAPullRequestWhoseHeadMovedIsNotEdited (0.00s)
PASS
ok  	github.com/pharzam/layup/internal/run	0.047s
exit=0
end 2026-10-19T08:01:38Z
```

Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -v -tags=integration -run '^(TestTheReviewRecordIsLinkedFromTheBodyAfterItsPost|TestTheLinkSurvivesTheCloseOutAndARestart|TestAnEditRefusedByTheForgeIsAForgeErrorOfEditBody)$' -timeout 30m ./internal/run/`

```
start 2026-10-19T08:01:38Z
--- PASS: TestTheReviewRecordIsLinkedFromTheBodyAfterItsPost (4.44s)
--- PASS: TestTheLinkSurvivesTheCloseOutAndARestart (5.67s)
--- PASS: TestAnEditRefusedByTheForgeIsAForgeErrorOfEditBody (4.13s)
PASS
ok  	github.com/pharzam/layup/internal/run	14.273s
exit=0
end 2026-10-19T08:02:09Z
```
