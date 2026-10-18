# T-08rj - the test runs of Step A (the forge call)

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
