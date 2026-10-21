# T-xyws: the test runs

Times are UTC, read from `date -u` around each run. Statuses are direct (`> file 2>&1; echo $?`), with no pipe.

## S3: the red run against the stubs of `6b06b11`

The new tests were appended and `gofmt -l` and `go vet` of both packages exited 0 first. Not committed at this point.

### Unit, 2026-10-18T18:06:14Z to 18:06:16Z, exit 1

`go test -count=1 ./internal/forge/`

```
--- FAIL: TestFakeClosesAPullRequestAndReadsItClosedByNumberAndByBranch (0.00s)
    forge_test.go:314: ClosePullRequest(1) = {Number:0 NodeID: State: Draft:false Merged:false HeadSHA: MergeSHA:}, <nil>; want Number 1, State closed, Merged false
--- FAIL: TestFakeCloseOfAMergedPullRequestChangesNothingAndGivesItMerged (0.00s)
    forge_test.go:345: ClosePullRequest of a merged pull request = {Number:0 NodeID: State: Draft:false Merged:false HeadSHA: MergeSHA:}, <nil>; want Number 1, Merged true, State closed, MergeSHA "0000000000000000000000000000000000abc001", no error
--- FAIL: TestFakeClosesAnIssueAsNotPlannedAndKeepsItAfterASecondClose (0.00s)
    forge_test.go:358: IssueState of an open issue = "", "", false; want open, an empty reason, true
    forge_test.go:364: IssueState after the close = "", "", false; want closed, not_planned, true
    forge_test.go:370: IssueState after a second close = "", "", false; want closed, not_planned, true
--- FAIL: TestFakeRefusesACloseBelow1BeforeItsFailureAndGivesTheForgeErrors (0.00s)
    forge_test.go:395: ClosePullRequest(0) with Fail = <nil>, want a refusal of that call with the status 0
    forge_test.go:395: ClosePullRequest(-1) with Fail = <nil>, want a refusal of that call with the status 0
    forge_test.go:399: ClosePullRequest(1) with Fail = <nil>, want the forge error of status 502
    forge_test.go:403: ClosePullRequest(99) = <nil>, want a forge error of status 404
    forge_test.go:395: CloseIssue(0) with Fail = <nil>, want a refusal of that call with the status 0
    forge_test.go:395: CloseIssue(-1) with Fail = <nil>, want a refusal of that call with the status 0
    forge_test.go:399: CloseIssue(1) with Fail = <nil>, want the forge error of status 502
    forge_test.go:403: CloseIssue(99) = <nil>, want a forge error of status 404
FAIL
FAIL	github.com/pharzam/layup/internal/forge	0.370s
FAIL
```

### Integration, 2026-10-18T18:06:16Z to 18:06:19Z, exit 1

`go test -count=1 -tags=integration ./internal/forge/github/`

```
--- FAIL: TestClosePullRequestPatchesTheStateAloneAndGivesThePullRequest (0.00s)
    github_integration_test.go:1285: PATCH /repos/acme/target/pulls/5 requests = 0; want 1 (err <nil>)
--- FAIL: TestClosePullRequestGivesAMergedPullRequestAsTheForgeAnswersIt (0.00s)
    github_integration_test.go:1303: ClosePullRequest of a merged one, call 1 = {Number:0 NodeID: State: Draft:false Merged:false HeadSHA: MergeSHA:}, <nil>; want Merged true, State closed, MergeSHA 3333333333333333333333333333333333333333, no error
    github_integration_test.go:1303: ClosePullRequest of a merged one, call 2 = {Number:0 NodeID: State: Draft:false Merged:false HeadSHA: MergeSHA:}, <nil>; want Merged true, State closed, MergeSHA 3333333333333333333333333333333333333333, no error
--- FAIL: TestCloseIssuePatchesTheStateAndTheReasonAlone (0.00s)
    github_integration_test.go:1317: PATCH /repos/acme/target/issues/7 requests = 0; want 1 (err <nil>)
--- FAIL: TestTheClosesRefuseANumberBelow1BeforeAnyRequest (0.00s)
    github_integration_test.go:1337: ClosePullRequest(0) = <nil>, want a refusal of the call ClosePullRequest with the status 0
    github_integration_test.go:1341: CloseIssue(0) = <nil>, want a refusal of the call CloseIssue with the status 0
    github_integration_test.go:1337: ClosePullRequest(-1) = <nil>, want a refusal of the call ClosePullRequest with the status 0
    github_integration_test.go:1341: CloseIssue(-1) = <nil>, want a refusal of the call CloseIssue with the status 0
--- FAIL: TestTheClosesGiveAForgeErrorForAnyOtherStatusOrAnAnswerThatIsNotClosed (0.00s)
    github_integration_test.go:1390: pull 422: err = <nil>; want the call ClosePullRequest, the status 422 and the message "rejected by the test"
    github_integration_test.go:1390: pull 404: err = <nil>; want the call ClosePullRequest, the status 404 and the message "rejected by the test"
    github_integration_test.go:1390: pull 410: err = <nil>; want the call ClosePullRequest, the status 410 and the message "rejected by the test"
    github_integration_test.go:1390: pull 403 with no reset time: err = <nil>; want the call ClosePullRequest, the status 403 and the message "rejected by the test"
    github_integration_test.go:1390: issue 422: err = <nil>; want the call CloseIssue, the status 422 and the message "rejected by the test"
    github_integration_test.go:1390: issue 404: err = <nil>; want the call CloseIssue, the status 404 and the message "rejected by the test"
    github_integration_test.go:1390: issue 410: err = <nil>; want the call CloseIssue, the status 410 and the message "rejected by the test"
    github_integration_test.go:1390: issue 403 with no reset time: err = <nil>; want the call CloseIssue, the status 403 and the message "rejected by the test"
    github_integration_test.go:1390: a pull request answer that is open: err = <nil>; want the call ClosePullRequest, the status 200 and the message "not closed"
    github_integration_test.go:1390: an issue answer that is open: err = <nil>; want the call CloseIssue, the status 200 and the message "not closed"
    github_integration_test.go:1390: an issue answer that is a pull request: err = <nil>; want the call CloseIssue, the status 200 and the message "the number is a pull request"
--- FAIL: TestTheClosesShareTheHandlingOfA401AndOfARateLimit (0.00s)
    github_integration_test.go:1420: ClosePullRequest: tokens made = 0, want 2
    github_integration_test.go:1428: ClosePullRequest: the wait: sleeps [], lines [], want three waits of at most ten seconds
    github_integration_test.go:1420: CloseIssue: tokens made = 0, want 2
    github_integration_test.go:1428: CloseIssue: the wait: sleeps [], lines [], want three waits of at most ten seconds
FAIL
FAIL	github.com/pharzam/layup/internal/forge/github	0.975s
FAIL
```

### Reading

Every failing line is a semantic failure of the stub (a zero value and a nil error, no request), not a compile error. The
cases of condition 3 of the plan review are red: the second issue close and `IssueState` of an open issue
(`TestFakeClosesAnIssueAsNotPlannedAndKeepsItAfterASecondClose`), the merged pull request
(`...ChangesNothingAndGivesItMerged`, `...GivesAMergedPullRequestAsTheForgeAnswersIt`), the statuses 200 of the
semantic errors (rows "answer that is open" and "is a pull request"), and the 401 and the rate limit with a token primed
first.

**Green negatives** (the stub satisfies them, so they prove nothing alone; each has its red control above): the
"made 0 requests" half of `TestTheClosesRefuseANumberBelow1BeforeAnyRequest` (control: the refusal lines of the same test,
which are red); the second `CloseIssue` returning no error in F2 (control: the `IssueState` lines); `IssueState` of an unknown
number in F2, since the stub's `"", "", false` is what the case expects (its mutation evidence is CM14 of `mutants.md`).

**Unreached in this run:** the refusal of `Merge` of a closed pull request in `TestFakeClosesAPullRequest...`. F1 stops
at `t.Fatalf` on its line 314 before it, so the run proves nothing of it, red or green.

*Corrected after round 1 (findings 2 and 3 of comment 6253669149): the raw log above is unchanged; this reading is.*

The line numbers of the integration run are those before `gofmt -w` wrapped two lines of the added tests.

## S4: green at the implementation

Unit, 2026-10-18T18:07:15Z to 18:07:17Z, exit 0: `go test -count=1 ./internal/forge/`
```
ok  	github.com/pharzam/layup/internal/forge	0.606s
```
Integration, 18:07:17Z to 18:07:22Z, exit 0: `go test -count=1 -tags=integration ./internal/forge/github/`
```
ok  	github.com/pharzam/layup/internal/forge/github	2.383s
```
