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
