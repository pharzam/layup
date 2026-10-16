# T-cr8j — the runs of the tests

Times from the clock of the host, UTC, 2026-10-15.

## Step T: red, before the code (18:36:23Z)

The types, the interface, the stand-in and the adapter's six methods were stubs that return
the zero value, so that the tests compile and fail by their own assertions, not by a build
error. `go test -tags=integration ./internal/forge/...`: exit 1, **26 tests failed**:

- `internal/forge` (9): `TestCheckLoopPermissionsNamesEachMissingPermission` (the stub gave no
  name), `TestValidStatusRefusesAWrongStateOrATooLongDescription` (the stub accepted `ok` and 141
  characters), and the seven `TestFake*` (no capability, no pull request, no status).
- `internal/forge/github` (17): for example `TestOpenDraftPostsADraftPullRequestAndGivesIt`:
  "the request body = map[], want map[base:main body:Closes #7 draft:true head:task/T-aaaa/1
  title:feat: T-aaaa x]", `TestMergeSendsTheHeadSHAAndAMergeCommitAndGivesTheMergeCommit`:
  `Merge = "", <nil>`, `TestGraphqlURLIs…` (no `/api/graphql`).
- A guard that must pass in both states, `TestTheLoopPermissionsLeaveTheM2aOnesAlone`, passed.

## Green, at the head before the freeze

- `go build ./...`, `go vet ./...`, `go vet -tags=integration ./...`: exit 0.
- `go test -count=1 ./...`: exit 0. `go test -count=1 -tags=integration ./...` (TMPDIR on the
  disk): exit 0, 18:45:08Z to 18:58:29Z (13 min 21 s, both runs).
- `adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests`, `setup-check`: exit 0 each.
  `git diff --check`: exit 0. `gofmt -l internal`: no file.

## Mutant runs (`go test -tags=integration ./internal/forge/...` on a changed copy of the code)

17 mutants, each restored after its run; **16 were caught at the first run, one was alive**:
M9 (`Merge` ignores `merged` false) passed because the test's answer had an empty `sha`, which the
code refuses anyway. The test now also answers `merged` false with a SHA, and M9 is caught. The
17: `OpenDraft` sends `draft` false; `MarkReady` puts the node ID in the text of the mutation; the
state check skipped; the description counted in bytes; `filter=latest` dropped; the next page of
the check runs unchecked; a branch's pull request taken from the end of the list; `Merge` without
`sha`; `Merge` ignoring `merged`; `Merge` with `squash`; the loop check reading the M2a list; the
stand-in ignoring a moved head; `graphqlURL` ignoring `/api/v3`; the stand-in appending a status
at the end; a branch's list without `state=all`; `MarkReady` ignoring an `errors` array;
`SetStatus` skipping the SHA check.

## After round 1 (the fixes, 19:25Z to 19:38Z)

Round 1 (comment 6200399675) had three material findings and ten notes. The fixes, test first where they
add a rule: the fixtures now follow the two shapes of a pull request (`prJSON`, a list item with
`merged_at` and no `merged`), a test merge commit before a merge, the mutation's `state` and
`headRefOid`; the refusals of a SHA in `Statuses` and of a negative number in `PullRequest` have tests
that count the requests; `forge.Fake` refuses what the adapter refuses and locks in every call; the
rows of `internal/forge` and `internal/forge/github` of the table of `M2a`, the cell of `NFR-001`, the
cell of `MarkReady` and the section "The test of the adapter" of `forge.md` are changed. Seven more
mutants (`merged_at` ignored; a test merge commit given as the merge commit; the state of `MarkReady` a
constant; the SHA check of `Statuses` skipped; a negative number accepted; the stand-in merging an empty
head; the head of `MarkReady` dropped): two were alive at first (the state constant and the negative
number) and the tests were tightened; **all 24 mutants are now caught**. `go build`, `go vet` (both tags),
`go test -count=1 ./...` and `go test -count=1 -tags=integration ./...`: exit 0 (19:25:53Z to
19:37:33Z, 11 min 40 s); the five checks: exit 0. Finding 12 of round 1 (five `internal/run` tests fail
at the step `read-back` when the temporary directory is inside the checkout) did not reproduce here:
`TMPDIR=<worktree>/.tmp go test -tags=integration ./internal/run -run 'TestARefusedPushBindsNothing|TestTheStepMarkersWithNoMarker'`
exit 0 (12.2 s).

## After round 2 (the fixes that wait for the Operator's answer), 2026-10-15T19:55Z

Round 2 (comment 6200984600) ended `not mergeable, findings recorded` at the cap of 1 (one material finding, eight notes). The fixes are a local commit, not pushed until the Operator chooses: the refusals of `Merge` are tested by a test that counts the requests (and whose handler would succeed); `forge.Fake` refuses before it fails, as the adapter refuses before any request, and its tests assert the status 0 of a refusal and the status 502 of `Fail`; `MarkReady` is tested on the state `MERGED`, `data` null and `pullRequest` null; the list fixture gives a merged pull request the state `closed`; the cells of the PRD and a traceability row. Five more mutants (a refusal of `Merge` removed; the stand-in ignoring `Fail` in `Merge` and in `SetStatus`; the merged state of `MarkReady` dropped; the stand-in accepting a negative number): **all 29 mutants are caught**. `go vet` and `go test -count=1 -tags=integration ./internal/forge/...`: exit 0; the five checks: exit 0. The full integration run is to be repeated before the push.

## After round 2: the full run and the merges of main (2026-10-16, UTC)

At `7896b20` (the fixes of round 2, with `origin/main` `2aafa64` merged; one conflict, a cell of `PRD-0001` §12, kept on both sides under O-219): the five checks, `git diff --check`, `go build ./...`, `go vet ./...` (both tags), `go test ./...`: exit 0; `TMPDIR=$HOME/.cache/layup-session/tmp-green-cr8j go test -tags=integration -timeout 25m ./...`: **exit 0, 06:10:23Z to 06:23:15Z (12 min 52 s)**, 21 packages ok (`internal/forge` 0.218 s, `internal/forge/github` 1.139 s, `internal/run` 721.095 s); gitleaks alone on `origin/main..HEAD`: exit 0. This replaces the sentence above, "the full integration run is to be repeated", which was true when written.
At `b91afe8` (`origin/main` `d24fae3` merged, no conflict; `PRD-0001`, `loop.md`, `packages.md` and `traceability.md` merged with disjoint hunks): `setup-check`, `git diff --check`, `go build`, `go vet` (both tags), `go test ./...`: exit 0; `go test -tags=integration -timeout 25m ./internal/forge/... ./internal/session/... ./internal/handoff/... ./internal/spec/... ./cmd/layup/...`: exit 0; gitleaks alone: exit 0; the full integration run was not repeated; the CI on the full SHA of the landing head is the check of that head (the job `tests` of `b91afe8` was still running when it was read, so no result of it is recorded here).
