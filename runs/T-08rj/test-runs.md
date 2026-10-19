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

Every command was run directly in `.worktree/T-08rj`, with its exit status read directly; the times are `date -u` read
around each run. The red runs are against the compiling stub of commit `f30ba1c` (`LinkReviewRecord` returns false, nil and
reads nothing; `reviewLine` returns ""), with the final tests of Step B. Lines are cut at 230 characters. The unit tests U1
to U13 and the integration tests I1 to I3 are named in `docs/tests/traceability.md`. The integration red run was taken
before the one assertion of the web address that kills RM11 was added to I1.

Differences from the plan, found at S6 (the reconciliation against the landed 66g, PR #421, merge `c735788`):

- `PullBody(issue, web, repo, records, more ...string)` already takes further lines, so there is no `review` parameter to add:
  the link is `more = "", reviewLine(...)`, and U1 is red only through the stub of `reviewLine`.
- `Task.PR` is written by `registerPR` and `holdsLease` is in `handoffpost.go`; both are as planned. `Sessions.Web` is set by
  `state.sessions()` from `cfg.Register.Web` (`internal/run/meaning.go:139`).
- Row 66g's text (`loop.md`) already excludes the body edit and names row 66q as its owner; row 66q (After `64, 66g`) and the
  After cells of 66h and 66k (`66q`) are on main.
- The review's conditions 1 to 3 changed the design of Stage 0: the snapshot is the first session that bound the reviewed
  head (not "the session bound at that SHA"); the latest handoff is chosen and then required valid; the pull request is read by
  the attempt's branch and must be the register's, at the reviewed head.

## Red, the unit tests (U1 to U13)

Command: `go test -count=1 -run '<the 13 tests>' ./internal/run/`

```
start 2026-10-19T06:06:35Z
--- FAIL: TestTheBodyWithAReviewRecordEndsInTheLinkOfItsComment (0.00s)
    reviewlink_test.go:181: the body = "Closes #3\n\nThe records commit this task started from: https://github.example/acme/target/commit/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n\n"; want "Closes #3\n\nThe records commit this ta...
    reviewlink_test.go:184: reviewLine with a web address that ends in a slash = ""
--- FAIL: TestTheBodyIsEditedOnceToLinkTheReviewRecord (0.00s)
    reviewlink_test.go:193: LinkReviewRecord = false, <nil>; want true, nil
--- FAIL: TestAnEqualBodyIsNotEdited (0.00s)
    reviewlink_test.go:215: LinkReviewRecord = false, <nil> with the calls []; want false, nil and one PullRequest
--- FAIL: TestTheRecordsLinkIsTheProducingSessionsNotTheCloseOutsNorTheVerifiers (0.00s)
    reviewlink_test.go:242: after a person's edit: false, <nil>, body "a person's text"; want the whole body replaced
--- FAIL: TestEachRecordThatFailsRefusesTheLinkBeforeAnyCall (0.01s)
    reviewlink_test.go:313: the session is of another task: err = <nil>; want one that says "not of T-ab12"
    reviewlink_test.go:320: the session is of another task: the control with the one token repaired = false, <nil>; want true, nil
    reviewlink_test.go:313: the session is a developer's: err = <nil>; want one that says "not of the verification"
    reviewlink_test.go:320: the session is a developer's: the control with the one token repaired = false, <nil>; want true, nil
    reviewlink_test.go:313: a newer attempt exists: err = <nil>; want one that says "newest attempt"
    reviewlink_test.go:320: a newer attempt exists: the control with the one token repaired = false, <nil>; want true, nil
    reviewlink_test.go:313: the session has no handoff: err = <nil>; want one that says "has no handoff"
    reviewlink_test.go:320: the session has no handoff: the control with the one token repaired = false, <nil>; want true, nil
    reviewlink_test.go:313: the latest handoff is invalid after an earlier valid one: err = <nil>; want one that says "latest handoff"
    reviewlink_test.go:320: the latest handoff is invalid after an earlier valid one: the control with the one token repaired = false, <nil>; want true, nil
    reviewlink_test.go:313: the latest handoff names no comment: err = <nil>; want one that says "latest handoff"
    reviewlink_test.go:320: the latest handoff names no comment: the control with the one token repaired = false, <nil>; want true, nil
    reviewlink_test.go:313: the copy is not in copies.tsv: err = <nil>; want one that says "no copy of the comment"
    reviewlink_test.go:320: the copy is not in copies.tsv: the control with the one token repaired = false, <nil>; want true, nil
    reviewlink_test.go:313: the copy is of another issue: err = <nil>; want one that says "no copy of the comment"
    reviewlink_test.go:320: the copy is of another issue: the control with the one token repaired = false, <nil>; want true, nil
    reviewlink_test.go:313: the copy is a plan review: err = <nil>; want one that says "not a review record"
    reviewlink_test.go:320: the copy is a plan review: the control with the one token repaired = false, <nil>; want true, nil
    reviewlink_test.go:313: the record is of another head: err = <nil>; want one that says "is of the head"
    reviewlink_test.go:320: the record is of another head: the control with the one token repaired = false, <nil>; want true, nil
    reviewlink_test.go:313: the register names no pull request: err = <nil>; want one that says "no pull request in the register"
    reviewlink_test.go:320: the register names no pull request: the control with the one token repaired = false, <nil>; want true, nil
    reviewlink_test.go:313: no session bound the reviewed head: err = <nil>; want one that says "no event bound at the head"
    reviewlink_test.go:320: no session bound the reviewed head: the control with the one token repaired = false, <nil>; want true, nil
--- FAIL: TestALostLeaseEditsNothing (0.00s)
    reviewlink_test.go:332: LinkReviewRecord = false, <nil> with the calls []; want *LostError and no EditBody
--- FAIL: TestAClosedOrMergedPullRequestIsNotEdited (0.00s)
    reviewlink_test.go:343: a closed pull request: err = <nil>; want one that says "is closed"
    reviewlink_test.go:353: a merged pull request: err = <nil>; want one that says "is closed"
--- FAIL: TestAFailedEditIsAnErrorAndAnEditThatTookEffectIsNotRepeated (0.00s)
    reviewlink_test.go:365: a failed edit: err = <nil>; want one that says "EditBody: status 502"
    reviewlink_test.go:371: the next call = false, <nil> with the calls []; want one edit
    reviewlink_test.go:377: an edit that took effect: err = <nil>; want one that says "EditBody: status 502"
    reviewlink_test.go:379: an edit that took effect: edited false, body "Closes #3\n\nThe records commit this task started from: https://github.example/acme/target/commit/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"; want the body e...
--- FAIL: TestAFailedReadOfThePullRequestEditsNothing (0.00s)
    reviewlink_test.go:392: a failed read: err = <nil>; want one that says "PullRequest: status 502"
--- FAIL: TestAFailedStoreReadIsAnErrorWithNoForgeCall (0.00s)
    reviewlink_test.go:405: a failed read of tasks.tsv: false, <nil> with the calls []; want the store's error and no call
    reviewlink_test.go:405: a failed read of sessions.tsv: false, <nil> with the calls []; want the store's error and no call
    reviewlink_test.go:405: a failed read of tasks/T-ab12/events.tsv: false, <nil> with the calls []; want the store's error and no call
    reviewlink_test.go:405: a failed read of copies.tsv: false, <nil> with the calls []; want the store's error and no call
    reviewlink_test.go:405: a failed read of copies/601-1.md: false, <nil> with the calls []; want the store's error and no call
    reviewlink_test.go:411: a failed read of the lease: false, <nil> with the calls []; want the store's error and no edit
    reviewlink_test.go:416: a copy whose file is missing: err <nil>, calls []; want an error and no call
--- FAIL: TestAnAnswerWithAnotherBodyIsAnError (0.00s)
    reviewlink_test.go:425: another body: err = <nil>; want one that says "gave another body"
--- FAIL: TestANewerAttemptWithTheRegisterOnTheOlderPullRequestEditsNothing (0.00s)
    reviewlink_test.go:457: no pull request of the branch: err = <nil>; want one that says "no pull request"
    reviewlink_test.go:468: the register names the older pull request: err = <nil>; want one that says "not the register's"
    reviewlink_test.go:476: with the register on 2: false, <nil>; want true, nil
--- FAIL: TestAPullRequestWhoseHeadMovedIsNotEdited (0.00s)
    reviewlink_test.go:491: a moved head: err = <nil>; want one that says "not the reviewed"
FAIL
FAIL	github.com/pharzam/layup/internal/run	0.107s
FAIL
exit=1
end 2026-10-19T06:06:50Z
```

All 13 tests fail at the stub, each on the assertion of the behavior it names (`edited = false; want true`, `err = nil; want ...`,
and for the rows of U5 both the refusal and its one-token control). U3 and the control rows pass nothing at the stub.

## Red, the integration tests (I1 to I3)

Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -run '<the 3 tests>' -timeout 30m ./internal/run/`

```
start 2026-10-19T06:06:50Z
--- FAIL: TestTheReviewRecordIsLinkedFromTheBodyAfterItsPost (6.35s)
    reviewlink_integration_test.go:94: LinkReviewRecord = false, <nil>; want true, nil
--- FAIL: TestTheLinkSurvivesTheCloseOutAndARestart (7.18s)
    reviewlink_integration_test.go:130: after the close-out and a restart: false, <nil> with 0 PATCH; want false, nil and the one
--- FAIL: TestAnEditRefusedByTheForgeIsAForgeErrorOfEditBody (5.71s)
    reviewlink_integration_test.go:144: LinkReviewRecord = false, <nil>; want EditBody: status 403
FAIL
FAIL	github.com/pharzam/layup/internal/run	19.310s
FAIL
exit=1
end 2026-10-19T06:07:32Z
```

Two earlier red runs failed for the wrong reason and were fixed in the tests before this one: the verifier's session could not
start at the developer's head (a session clones the run's `main` alone; the test moves it onto the head, as the loop would), and
the verifier and the developer shared one records snapshot (the test now starts the verifier at the current records head).

## Green, the link tests

Command: `go test -count=1 -v -run '<the 13 tests>' ./internal/run/` (the pattern also runs three older tests whose names begin alike)

```
start 2026-10-19T06:07:50Z
--- PASS: TestALostLeaseStopsTheComment (0.00s)
--- PASS: TestALostLeaseOpensMarksAndMergesNothing (0.00s)
--- PASS: TestTheBodyWithAReviewRecordEndsInTheLinkOfItsComment (0.00s)
--- PASS: TestTheBodyIsEditedOnceToLinkTheReviewRecord (0.00s)
--- PASS: TestAnEqualBodyIsNotEdited (0.00s)
--- PASS: TestTheRecordsLinkIsTheProducingSessionsNotTheCloseOutsNorTheVerifiers (0.00s)
--- PASS: TestEachRecordThatFailsRefusesTheLinkBeforeAnyCall (0.02s)
--- PASS: TestALostLeaseEditsNothing (0.00s)
--- PASS: TestAClosedOrMergedPullRequestIsNotEdited (0.00s)
--- PASS: TestAFailedEditIsAnErrorAndAnEditThatTookEffectIsNotRepeated (0.00s)
--- PASS: TestAFailedReadOfThePullRequestEditsNothing (0.00s)
--- PASS: TestAFailedStoreReadIsAnErrorWithNoForgeCall (0.00s)
--- PASS: TestAnAnswerWithAnotherBodyIsAnError (0.00s)
--- PASS: TestANewerAttemptWithTheRegisterOnTheOlderPullRequestEditsNothing (0.00s)
--- PASS: TestAPullRequestWhoseHeadMovedIsNotEdited (0.00s)
--- PASS: TestALostLeaseStopsTheStart (0.00s)
--- PASS: TestALostLeasePostsNoStatus (0.00s)
PASS
ok  	github.com/pharzam/layup/internal/run	0.128s
exit=0
end 2026-10-19T06:07:51Z
```

Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -v -tags=integration -run '<the 3 tests>' -timeout 30m ./internal/run/`

```
start 2026-10-19T06:07:51Z
--- PASS: TestTheReviewRecordIsLinkedFromTheBodyAfterItsPost (6.51s)
--- PASS: TestTheLinkSurvivesTheCloseOutAndARestart (7.69s)
--- PASS: TestAnEditRefusedByTheForgeIsAForgeErrorOfEditBody (6.10s)
PASS
ok  	github.com/pharzam/layup/internal/run	20.368s
exit=0
end 2026-10-19T06:08:34Z
```
