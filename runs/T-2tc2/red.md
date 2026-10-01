# T-2tc2: the red runs

The red step of each test of `T-2tc2` (#79), on the final test files, before
the code that makes it green (gate step 3). Outputs are cut to the relevant
lines ("…"). Host: macOS, `go1.27.1`, `git` 2.54.0, 2026-10-01.

Runs 1 and 2 are on a skeleton of `internal/git`: its exported names, and no
behaviour. Each test fails on its own assertion, because no call starts `git`
or the stub: no repository, commit, server request, output or error kind.

```text
$ go test -count=1 ./internal/git/                                   # run 1
        git_test.go:78: 0 starts of git, want 1      (TestEachCallRunsItsVerb, each of the 19 calls)
    git_test.go:100: 0 starts of git, want 1         (TestTheEnvironmentIsAFixedList)
    git_test.go:141: ls-remote: the line of exactly HEAD: "", <nil>; want "2222"   … 5 more, and line 147
    git_test.go:155: Supported("2.32.0") = false, want true                       … and 4 more
    git_test.go:164: git missing: <nil>; want a *NotFoundError                    (TestTheTwoErrorKinds)
$ go test -count=1 -tags=integration ./internal/git/                 # run 2
    git_integration_test.go:101: commit "" (<nil>): want 40 hexadecimal characters
        … the same in TestCloneAndCheckout, TestBranchesAndAnOrphan, TestAScratchTree, TestAHostileHostChangesNothing
    git_integration_test.go:296: clone: <nil>; want a *FailedError that says terminal prompts disabled
    git_integration_test.go:302: git did not call the server, so the test proves nothing
    git_integration_test.go:313: git failed: <nil>; want a *FailedError with code 128 and the standard error
    git_integration_test.go:324: an empty PATH: <nil>; want a *NotFoundError
```

Runs 3 and 3b are on the real package with the isolation of condition 1 taken
out for the run. Run 3 has the settings of the plan before its review (the
host's environment, `GIT_CONFIG_NOSYSTEM=1`, an empty `GIT_CONFIG_GLOBAL`,
`GIT_TERMINAL_PROMPT=0`, no `core.attributesFile` or `core.excludesFile`): the
host's `GIT_DIR` moves the repository. Run 3b has the fixed environment and no
`core.attributesFile` or `core.excludesFile`: the per-user ignore file drops
`doc.md`, and the commit differs. In both, the control subtest passes, so each
seeded input has an effect.

```text
$ go test -count=1 -tags=integration -run TestAHostileHostChangesNothing ./internal/git/
    git_integration_test.go:262: Init made no repository in …/005: GIT_DIR of the host sent it to …/004/decoy.git  # run 3
    git_integration_test.go:270: commit "9b053301…", want 1923606c…: the host changed the tree or the author      # run 3b
         crlf.txt | 2 ++
         run.sh   | 1 +
```

Runs 4 and 5 are on a skeleton of the checker that reads no row and gives no
finding; run 6 is the real checker on the table of the base; run 7 is the real
checker with rule 5 off. The skeleton refuses nothing and finds nothing (the
case of a good module passes on it, so the breach cases prove that its pass is
not empty); the test stops when it reads no row; the program cell of the base
does not start with its code span; and the fixture breaks rule 5 only, so with
rule 5 off the same checker and `go list` find nothing in it.

```text
$ go test -count=1 ./cmd/layup/                                      # run 4
    rules_test.go:68: rows map[], <nil>                                (TestReadTable)
    rules_test.go:93: zero rows: no error, 0 rows; want an error       … and the other 13 tables
        rules_test.go:170: findings … want exactly rule 5: internal/psb depends on net/http   … and the other 17 breaches
$ go test -count=1 -tags=integration -run TestPackageRules ./cmd/layup/
    rules_integration_test.go:25: the table of phase 1: 0 rows, <nil>                                     # run 5
    rules_integration_test.go:25: the table of phase 1: 0 rows, internal/gate: cannot read
        Starts a program "the gate commands, with `sh -c`"                                                # run 6
    rules_integration_test.go:36: the fixture that imports net/http gives … want the finding
        rule 5: cmd/layup depends on net/http                                                             # run 7
```

## The fixes of the verification

A fresh verifier checked the head `9866f48` and gave twelve findings. Each red
run below is on the new test and the code of that head, before the fix.

Run 8 (finding 1): the host's `HOME` passed to `git`, and `ssh` and a remote
helper could start. The unit test fails because the list holds the host's
`HOME` and `GIT_SSH_COMMAND`. The integration test fails because the server
got an `Authorization` header from the seeded `.netrc`, and the fake `ssh` and
the fake remote helper on the `PATH` both started.

```text
$ go test -count=1 ./internal/git/                                   # run 8
    git_test.go:80: dir "", args [… "-c" "commit.gpgsign=false" "--version"]
        want dir "", args [… "-c" "commit.gpgsign=false" "-c" "http.emptyAuth=false" "--version"]   … each call
    git_test.go:102: env
         got ["LC_ALL=C" "GIT_CONFIG_NOSYSTEM=1" … "GIT_SSH_COMMAND=ssh -o BatchMode=yes" "PATH=/stub/bin" "HOME=/stub/home" "TMPDIR=/stub/tmp"]
        want ["LC_ALL=C" "HOME=/dev/null" "GIT_CONFIG_NOSYSTEM=1" … "GIT_ALLOW_PROTOCOL=file:git:http:https" "PATH=/stub/bin" "TMPDIR=/stub/tmp"]
$ go test -count=1 -tags=integration -run TestNoCallUsesACredentialOfTheHost ./internal/git/
    git_integration_test.go:308: clone of ssh://git@127.0.0.1/baseline.git: … fatal: Could not read from remote
        repository. …; want a *FailedError that says transport 'ssh' not allowed
    git_integration_test.go:308: clone of layuptest::baseline: … fatal: remote helper 'layuptest' aborted session;
        want a *FailedError that says transport 'layuptest' not allowed
    git_integration_test.go:315: the server got 2 requests, 1 with an Authorization header; want one or more, none with it
    git_integration_test.go:318: 2 programs of the PATH started, first git-remote-layuptest; want none
```
