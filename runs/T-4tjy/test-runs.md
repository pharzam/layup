# T-4tjy — test runs

The check is `runs/T-efmy/release-check.sh`, adapted (D3 of the plan of #169,
with conditions 3 and 5 and notes 3, 5 and 6 of its plan review). Each run is
from a copy of the script outside the clone, as its header asks.

## Red: the check of phase 1 at the base (2026-10-10, 07:16Z to 07:17Z)

The script before this task, at `9490de9`, the base: exit 1, on check (4), the
five literals of `Fetch`, `Push`, `CloneLocal` and `FetchSession` (#148):

```
FAIL: a git verb that reaches a remote
== FAIL
```

## A first draft that read no file (2026-10-10T07:27Z)

The first draft of (4) listed the files of `internal/git` with `git ls-files
'internal/git/*.go' ':!*_test.go'`, which gives no file at all, so (4) printed
`none` and passed. The run of the draft showed it; the files are now filtered
by name, and a list with no file fails, as does (3) with no call found.

## Red: the mutations (2026-10-10T07:29Z)

Each mutation is its own commit on a scratch clone of `9490de9` (condition 5),
and the check runs with that commit, so check (1) holds and each exit 1 has its
own reason. The `FAIL` lines of each run, in full:

```
== curl-in-gate
FAIL
FAIL	github.com/pharzam/layup/cmd/layup	1.262s
FAIL
FAIL: TestPackageRules
FAIL: a call that starts a program, out of the list: internal/gate/scratch.go:190:	cmd := exec.Command("curl", "-c", command)
== ssh-in-session
FAIL
FAIL	github.com/pharzam/layup/cmd/layup	1.283s
FAIL
FAIL: TestPackageRules
FAIL: a call that starts a program, out of the list: internal/session/process.go:132:	cmd := exec.Command("ssh", s.Words[1:]...)
== push-in-new-func
FAIL: a git verb that reaches a remote, out of its function: internal/git/git.go:451: push in PushAll
== fetch-in-clonelocal
FAIL: a git verb that reaches a remote, out of its function: internal/git/git.go:395: fetch in CloneLocal
== push-in-no-function
FAIL: a git verb that reaches a remote, out of its function: internal/git/git.go:451: push in no function
```

The mutations: `curl-in-gate`, `exec.Command("curl", ...)` in
`internal/gate/scratch.go`; `ssh-in-session`, `exec.Command("ssh", ...)` in
`internal/session/process.go`; `push-in-new-func`, a function `PushAll` of
`internal/git` with `"push"`; `fetch-in-clonelocal`, a `"fetch"` in
`CloneLocal`; `push-in-no-function`, a top-level `var` with `"push"` after a
function. Each is caught by its line of (3) or (4); the first two also by
`TestPackageRules` (check (2)), as its program rules say.

## Green: the release of M2b (2026-10-10T07:29Z)

The adapted check at `9490de9`, in a clean clone: exit 0, `== PASS`. Its
output is [`release-check.txt`](release-check.txt).
