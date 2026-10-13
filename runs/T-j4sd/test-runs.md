# T-j4sd — test runs

The call `CommitTree` of `internal/git` (#272, row 51b), and `commitEnv`, the one home of the six
variables of the identity that `Commit` and `CommitTree` share. All times are UTC, 2026-10-13,
read from the clock after each step.

## The base

The plan was made against `9144e4e` (the merge of #305). The branch took `38ee78e` (the merge of
#306, `T-9v92`) by a fast-forward before it had a commit of its own, and the red run and the
mutation runs ran on that base with the change uncommitted; the mutants were made on a copy of the
tree whose `internal/git` equals the two commits that follow. Before the freeze the branch took
`6734a23` (the merge of #307, `T-4c5c`) by a merge, `da6cb27`, so the budget is read against
`6734a23`, as the comment "The base of the budget" on #272 says. The green run is on the merged head.

## Red (09:32:16Z to 09:32:19Z)

The tests before the code: `internal/git` does not build, for the call that the tests name:

```
== go vet ./internal/git/ and go vet -tags=integration ./internal/git/
vet: internal/git/git_test.go:459:13: undefined: CommitTree
vet: internal/git/git_integration_test.go:853:13: undefined: CommitTree
== go test ./internal/git/
internal/git/git_test.go:459:13: undefined: CommitTree
internal/git/git_test.go:478:15: undefined: CommitTree
internal/git/git_test.go:486:2: undefined: CommitTree
internal/git/git_test.go:501:15: undefined: CommitTree
internal/git/git_test.go:510:16: undefined: CommitTree
internal/git/git_test.go:521:14: undefined: CommitTree
```

## Mutations

Each rule of the new code mutated alone, on a copy of the tree and never on the worktree: a
condition, a value or an argument changed so that each mutant builds. The tests run in two
passes: the unit level of `internal/git` (09:34:25Z to 09:35:32Z), and the integration level with
`-tags=integration -run 'TestCommitTree|TestInitAddAndCommit'` (09:35:42Z to 09:36:27Z), which runs the unit tests of the
same names too. The line of a mutant is the first test that failed; a mutant that does not build
would be counted as not run, never as caught.

| Mutant | The change | Unit level | Integration level |
| ------ | ---------- | ---------- | ----------------- |
| `ct-refuse-tree` | the refusal of a `tree` that is no full object ID removed | caught: `TestCommitTreeTakesFullObjectIDsAndNoOtherText` | caught: `TestCommitTreeTakesFullObjectIDsAndNoOtherText` |
| `ct-refuse-parent` | the refusal of a `parent` that is no full object ID removed | caught: `TestCommitTreeTakesFullObjectIDsAndNoOtherText` | caught: `TestCommitTreeTakesFullObjectIDsAndNoOtherText` |
| `ct-refuse-code` | the refusal's `Code` -1 (refused, `git` not started) changed to 0 | caught: `TestCommitTreeTakesFullObjectIDsAndNoOtherText` | caught: `TestCommitTreeTakesFullObjectIDsAndNoOtherText` |
| `ct-args-order` | `tree` and `parent` swapped in the arguments | caught: `TestCommitTreeRunsItsVerb` | caught: `TestCommitTree` |
| `ct-args-parent` | `-p PARENT` dropped from the arguments | caught: `TestCommitTreeRunsItsVerb` | caught: `TestCommitTree` |
| `ct-args-message` | the message replaced by the text `message` | caught: `TestCommitTreeRunsItsVerb` | caught: `TestCommitTree` |
| `ct-env` | no identity: the call runs with the fixed list only | caught: `TestCommitTreeTakesTheIdentityOfCommit` | caught: `TestCommitTree` |
| `ct-output-check` | the output no longer read as an object ID | caught: `TestCommitTreeReadsItsOutput` | caught: `TestCommitTreeReadsItsOutput` |
| `ct-output-trim` | the line feed of the output kept | caught: `TestCommitTreeRunsItsVerb` | caught: `TestCommitTree` |
| `ct-output-code` | the `Code` of an unreadable output changed from 0 to 1 | caught: `TestCommitTreeReadsItsOutput` | caught: `TestCommitTreeReadsItsOutput` |
| `ct-git-error` | a failure of `git` ignored | caught: `TestCommitTreeReadsItsOutput` | caught: `TestCommitTree` |
| `env-committer-date` | the committer's date fixed to `@0 +0000` | caught: `TestEachCallRunsItsVerb` | caught: `TestInitAddAndCommit` |
| `env-author-email` | the author's email taken from the name | caught: `TestEachCallRunsItsVerb` | caught: `TestInitAddAndCommit` |
| `env-date-form` | the zone of both dates changed to `+0100` | caught: `TestEachCallRunsItsVerb` | caught: `TestInitAddAndCommit` |
| `commit-env` | `Commit` runs with no identity | caught: `TestEachCallRunsItsVerb` | caught: `TestInitAddAndCommit` |

Unit level: `15 caught, 0 survived, 0 not run, of 15`. Integration level: `15 caught, 0 survived, 0 not run, of 15`.

## Green (10:05:35Z to 10:16:17Z, head `da6cb27`)

Each with exit 0, in this order; the time of a line is when its command ended. The three runs of
`go test` are with `-count=1 -timeout 30m`, then with `-tags=integration` and with `-tags=e2e`, and
`git diff --check` reads the commits against `origin/main` (`6734a23`):

```
go build ./...: exit 0 (10:05:37Z)
go vet ./...: exit 0 (10:05:38Z)
gofmt -l internal cmd (empty): exit 0 (10:05:38Z)
go test ./...: exit 0 (10:05:48Z)
go test -tags=integration ./...: exit 0 (10:13:57Z)
go test -tags=e2e ./...: exit 0 (10:15:28Z)
adr-lint: exit 0 (10:15:29Z)
prd-lint: exit 0 (10:15:29Z)
link-lint: exit 0 (10:15:40Z)
run-discipline-tests: exit 0 (10:15:47Z)
nested-checkout-check: exit 0 (10:15:48Z)
setup-check: exit 0 (10:16:17Z)
git diff --check origin/main HEAD: exit 0 (10:16:17Z)
```

A first run of the same script began at 10:02:22Z on `a903ae5`, the head of the two commits before
the merge. It was stopped after `go build`, `go vet`, `gofmt` and `go test ./...` (each exit 0, to
10:02:34Z), while the integration level ran, because `origin/main` had moved to `6734a23` (the merge
of #307, `T-4c5c`), and the branch took it by a merge (`da6cb27`; one conflict, in `PRD-0001`, solved
by keeping each task's rows and both lines of section 13). The run above is the whole script again,
on the merged head.
