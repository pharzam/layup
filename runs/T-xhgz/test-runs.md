# T-xhgz — the test runs

Evidence for row 23 (#128): each test red first, for the right reason, then green. `git` 2.54.0 on the host.

## The measurements of the plan review (2026-10-08T09:12Z)

- `git init -b main dst`, then `git -C dst fetch --no-tags -- ../bare.git refs/heads/main:refs/heads/main`: exit 128, `fatal: refusing to fetch into branch 'refs/heads/main' checked out at …` (condition 2). With `--update-head-ok`: exit 0, `refs/heads/main` the fetched commit.
- A push that is not a fast-forward: exit 1, standard output `!	<commit>:refs/heads/main	[rejected] (fetch first)`, standard error `error: failed to push some refs`; a push to a remote that does not exist: exit 128 (condition 1).

## Red 1 (2026-10-08T09:13Z)

The rows of `Fetch` and `Push` (with and without a token, the environment of each row checked), `TestFetchAndPushRefuseTheirInputBeforeGitStarts` and `TestNoErrorHoldsTheToken`, before the calls: `go test ./internal/git/` did not compile (`undefined: Fetch`, `undefined: Auth`).

## Red 2 (2026-10-08T09:16Z)

With stubs that start no `git` and give no error: the four rows `0 starts of git, want 1`; each of the eight inputs `<nil>, 0 starts of git; want a *FailedError of code -1`; `TestNoErrorHoldsTheToken` `the error <nil> …`; the integration tests `a push that is not a fast-forward: <nil>; want a *FailedError of code 1`, `the fetch from a server that has no repository: <nil>`, `the server got Authorization "", want Basic …`. One row failed for another reason: the row of `Commit` now checks the environment, which holds the identity; the row got its six variables.

## Green, and two mutations (2026-10-08T09:20Z)

`Auth`, `Fetch`, `Push`: `go test ./internal/git/` and `go test -tags=integration ./internal/git/` pass. Mutations on a backup copy of `git.go`, each put back (`cmp` equal):
- `Fetch` with no `--update-head-ok` → `git fetch --no-tags -- file://…/target.git refs/heads/main:refs/heads/main: exit status 128` in `TestPushAndFetchOfABareRepository`;
- `Push` with `--force` → `a push that is not a fast-forward: <nil>; want a *FailedError of code 1`.

## The documents (2026-10-08T09:28Z)

Red before the edits: `sh runs/T-xhgz/docs.sh` gave ten `FAIL` lines, exit 1. After the edits of `packages.md`, `forge.md`, `guardrails.md`, the traceability and the PRD: ten `ok`, exit 0.

## Fixes of round 1 (2026-10-08T09:33Z)

Finding 1: a new rule of `docs.sh` ("the rule of web with a token") failed on the `packages.md` of the frozen head (put in place from `git show HEAD:` and then put back), and passed after the clause; the code and its unit test rows ("a token with no web", "a web that ends in a slash") were already there, so no Go test changed. `docs.sh`: eleven `ok`, exit 0.
