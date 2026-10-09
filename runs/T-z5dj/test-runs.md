# The test runs of T-z5dj

The host's `git` is 2.47.3 (`git --version`), not the 2.54.0 that the plan named; the list of keys is that of `git-config(1)` of 2.47.3.

## Red 1: the tests before the calls (2026-10-09T12:1xZ)

`go vet ./internal/git/` did not compile: `undefined: IsAncestor` (`git_test.go`), `undefined: CloneLocal` (`git_integration_test.go`).

## The control of the hostile configuration

The control is a test of the list, and it failed three times before it was live, each time naming the keys that did not fire, which led to one driver per key: `include.path` first set `core.fsmonitor` and hid it (it sets an `alias` now); a process filter that fails aborted `git diff` and `git add` (the process filter has its own attribute now); `git diff` stops at the first external diff that fails, and `textconv` is read only by the built-in diff (one `git diff` per path, `--no-ext-diff` for `textconv`); `core.alternateRefsCommand` is read in a repository that has an alternate, on a fetch that negotiates (the control's copy gets an alternate and fetches a repository with a commit it lacks). A first run hung for five minutes: a marker program that read its standard input deadlocked the process filter and `upload-pack`; the program no longer reads it. At the end, the control fires each of the twelve live keys, and the test lists the ten that no plain read fires.

## Red 2: a FetchSession that runs git in the session (2026-10-09T12:20Z)

On a copy of `git.go` (put back after), `FetchSession` ran `git status` in the session's work tree and fetched from the session's clone:

```
--- FAIL: TestFetchSessionRunsNoGitInTheSession (0.00s)
--- FAIL: TestFetchSessionRefusesASHAOfNoCommit (0.00s)
--- FAIL: TestTheCallsOfM2bRefuseTheirInputBeforeGitStarts (0.00s)
--- FAIL: TestFetchSession (0.19s)
--- FAIL: TestAHostileSessionRunsNothing (0.57s)
    git_integration_test.go:783: FetchSession ran the session's programs: [core.fsmonitor filter.y.clean]
```

## Green (2026-10-09T12:20Z)

Each with exit 0 on the tree of the commit `feat: T-z5dj …`: `go build ./...`, `go vet ./...`, `gofmt -l internal cmd` (empty), `go test ./...`, `go test -tags=integration ./...`; `adr-lint`, `prd-lint`, `link-lint`, `setup-check`, `run-discipline-tests`, `git diff --check`.

## The fix of round 1 (2026-10-09T12:38Z)

Finding 1: the list of keys is now the keys that `git help --config` of git 2.47.3 lists and that name a program or a shell command, 45 in all (a key that only picks a tool reaches a program through a key of the list); the control fires 15 of them on this host (with `git y`, `git z` for `includeIf.path`, and `git fetch gp` for `core.gitProxy`), and `FetchSession` fires none. `includeIf` needed the pattern `gitdir:<work tree>/`; `gitdir:<work tree>/.git/` did not match on 2.47.3, tried by hand.

Note 3, a mutation red for each other call (on a copy of `git.go`, put back after):

```
== noremove
--- FAIL: TestCloneLocal (0.04s)
--- FAIL: TestCloneLocalRemovesItsRemote (0.00s)
== local
--- FAIL: TestCloneLocalRemovesItsRemote (0.00s)
== context
--- FAIL: TestEachCallRunsItsVerb (0.00s)
    --- FAIL: TestEachCallRunsItsVerb/diff_-U0_of_one_file (0.00s)
== binary
--- FAIL: TestEachCallRunsItsVerb (0.00s)
    --- FAIL: TestEachCallRunsItsVerb/diff_--binary (0.00s)
--- FAIL: TestTheReadsOfM2b (0.07s)
== isancestor (exit 1 read as an error)
--- FAIL: TestIsAncestorReadsTheExitCode (0.00s)
--- FAIL: TestTheReadsOfM2b (0.11s)
```

The mutation `local` (no `--no-local`) is caught by the unit test only: a local clone hard-links its objects, so `TestCloneLocal` reads the log after the source is gone either way.

## The fix of round 2, under O-188 (2026-10-09T13:01Z)

`git help --config` of 2.47.3 (903 keys) lists `difftool.<tool>.path` and `mergetool.<tool>.path`, which the list of 45 lacked; they join it (47, of which 32 are asserted and 15 shown live by the control). It does not list `tar.<format>.command`, `trailer.<keyAlias>.cmd` or `.command`, or `sendemail.sendmailCmd`, which round 2 named from the manual of a later version; the test's comment says so. `TestAHostileSessionRunsNothing` passes: the control fires the same 15, and `FetchSession` fires none of the 47.
