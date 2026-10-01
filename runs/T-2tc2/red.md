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
