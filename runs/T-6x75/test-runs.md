# T-6x75: the test runs

The runs of `T-6x75` (#84). Host: macOS, `go1.27.1`, `git` 2.54.0, 2026-10-02.
A path of a temporary directory shows as `<tmp>/`; "…" marks a cut part of an
output.

## A measurement of the plan

For the harness (D6 of #84): on the base, each case of `docs/setup/tests/pin/`,
`kit-history/` and `identity/` was built as `run.sh` builds it, and
`setup-check.sh --only <check>` ran on it. For each of these 10 cases, and for
the checks `pin`, `kit-history` and `identity` of `frame/good-passthrough`,
the lines of `EXPECT` for a check are the whole output of the sh check for it,
as a set (for `pin/bad-commit`, `EXPECT` and the sh output have the two lines in
another order: condition 3 of the plan review). So the harness compares sets.

## The red runs

Each red run is of a test in its final form, on a skeleton of the code (the
names, with no behaviour), or on a copy of the branch where the code under test
is a skeleton.

1. **`IsShallow`** (skeleton: it gives `false`):

   ```text
   $ go test -count=1 ./internal/git/
   --- FAIL: TestEachCallRunsItsVerb/rev-parse_--is-shallow-repository
       git_test.go:83: 0 starts of git, want 1
   --- FAIL: TestTheOutputIsRead
       git_test.go:149: is-shallow: a shallow clone: %!q(bool=false), <nil>; want %!q(bool=true)
       git_test.go:159: IsShallow with the output yes: false, <nil>; want false and a *FailedError with code 0
   $ go test -count=1 -tags=integration -run TestIsShallow ./internal/git/
       git_integration_test.go:389: IsShallow(<tmp>/…/shallow): false, <nil>; want true
       git_integration_test.go:389: IsShallow(<tmp>/…/tree): false, <nil>; want true
   ```

2. **`internal/work`** (skeleton: schemas with no column, readers that read
   nothing). The block test was then changed to look each block up by its
   expected name; its new form was run on the same skeleton in a copy, with the
   same two failures.

   ```text
   $ go test -count=1 ./internal/work/
       work_test.go:35: Value(S01, name) = "", false; want "acme", true
       work_test.go:75: a source of no kind: <nil>; want a *tsv.Error at line 2
       … the same for the other four forms
   $ go test -count=1 -tags=integration ./internal/work/
       work_integration_test.go:28: schema setup-answers: the location: the block has "host:<work>/inputs/answers.tsv"; the Go schema has ""
           the number of columns: the block has "5"; the Go schema has "0"
       … the same for setup-record
       work_integration_test.go:59: ReadRecord of a file of another form: <nil>; want an error that starts with out/record.tsv: line 1
       work_integration_test.go:62: ReadAnswers with no file: <nil>; want an error that starts with inputs/answers.tsv:
   ```

3. **`internal/verify`, unit** (skeleton: no core finds anything, `Run` and
   `Check` give an empty table): each of the 11 tests fails. The first run
   stopped at a panic of `TestTheInputErrors`, which read the text of a nil
   error; the test now checks the error first, and the run below is of that
   form.

   ```text
   $ go test -count=1 ./internal/verify/
   --- FAIL: TestPinFindings … TestTheTable (11 tests)
       checks_test.go:78: a key twice, the first commit short:
            got []
           want ["key: commit appears 2 times, expected 1" "commit: not 40 hexadecimal characters: a959655"]
       verify_test.go:125: the rows
            got []
           want [{"discipline-tests" "not-active" "not built yet"} {"pin" "pass" ""} {"kit-history" "fail" "orphan: x"} …]
       verify_test.go:208: an old git: <nil>, 0 rows; want an input error that starts with "git 2.32.0 or newer is needed", and no row
       verify_test.go:256: a tree that cannot be removed: <nil>, 0 rows; want a *CleanupError for /stand-in/scratch and the whole table
   ```

4. **`internal/standin`** (skeleton: it builds nothing):

   ```text
   $ go test -count=1 ./internal/standin/
       standin_test.go:16: pinText: (empty); want source=file:///b … date=2026-10-02
       standin_test.go:26: pin.commit: "", false; want the changed value c2
   $ go test -count=1 -tags=integration -run TestAStandIn ./internal/standin/
       standin_integration_test.go:33: ls-remote  HEAD: "", … fatal: bad repository ''; want the commit
       standin_integration_test.go:63: the root commits of layup-setup: [], … chdir target: no such file or directory; want one
   ```

5. **The harness**, on a copy of `7cc82e6` where the three cores find
   nothing: each bad case fails, with its lines and its exit. The good cases
   pass there; a core that always fails is caught by them (`kit-history/good`,
   and the three checks of `frame/good-passthrough`). A second copy takes
   `facts` out of the lists and lists a group that does not exist:

   ```text
   $ go test -count=1 -tags=integration -run TestTheFixturesOfSetupCheck ./internal/verify/
       harness_integration_test.go:104: pin/bad-commit: the lines of pin
            got ["setup-check: pin OK"]
           want ["setup-check: pin FAIL commit: not 40 hexadecimal characters: a959655" "setup-check: pin FAIL key: commit appears 2 times, expected 1"]
       harness_integration_test.go:108: pin/bad-commit: exit 0; want ["1"]
       … the same for the other 8 bad cases of pin, kit-history and identity
   (the second copy)
       harness_integration_test.go:51: the group facts is in 0 lists; want 1
       harness_integration_test.go:56: the listed group jobs has no directory under docs/setup/tests/
   ```

6. **The frame on stand-in work areas**, on a copy where `Check` gives an
   empty table and the Go schema of `setup-verify` has no column:

   ```text
   $ go test -count=1 -tags=integration -run 'TestRunOnAStandInWorkArea|TestEachFindingOfATarget|TestTheInputErrorsOfAWorkArea|TestTheTableSchemaEqualsItsBlock' ./internal/verify/
       verify_integration_test.go:79: the rows []; want ["pin pass " "kit-history pass " "identity pass "], …
       verify_integration_test.go:124: a decision of the kit: [], <nil>; want kit-history fail, decisions: docs/decisions/ exists (kit step 4 deletes it)
       … the same for the other 5 findings
       verify_integration_test.go:149: no answers: <nil>; want an input error that starts with "inputs/answers.tsv: "
       verify_integration_test.go:173: schema setup-verify: the location: the block has "stdout"; the Go schema has ""
   ```

7. **The command row** (skeleton: no row `setup verify` in the command table).
   The three usage cases first passed there, for the wrong reason: `setup` was
   an unknown command. Each now checks its own reason, and the run below is of
   that form; the same change was made to the end-to-end scenario.

   ```text
   $ go test -count=1 -run TestSetupVerify ./internal/cli/
       verify_test.go:43: ["setup" "verify"]: exit 2, stdout "", stderr "layup: unknown command \"setup\"…"; want 2, nothing and "missing argument WORK" with the usage
       verify_test.go:54: ["pass" "clear"]: exit 2, stdout "", WORK ""; want 0 and the table of w
       verify_test.go:62: an input error: exit 2, stdout "", stderr "layup: unknown command \"setup\"…"
       verify_test.go:69: the usage has no setup verify: …
   ```

8. **The end-to-end scenarios**, on the binary of a copy of `7cc82e6`, which
   has no row `setup verify`:

   ```text
   $ go test -count=1 -tags=e2e -run TestSetupVerify ./cmd/layup/
       verify_e2e_test.go:38: exit 2, stdout: (empty)
       verify_e2e_test.go:53: no target: exit 2, stdout "", stderr "layup: unknown command \"setup\"…"
       verify_e2e_test.go:58: layup ["setup" "verify"]: exit 2, … "layup: unknown command \"setup\"…"
       … the same for the other two usage errors
   ```

## Review round 1: the fixes (cycle 1)

Review round 1 (`b31cc70`) gave `material` with five findings. The red runs of
the fixes, each on the code of `b31cc70`:

```text
$ go test -count=1 ./internal/verify/
    checks_test.go:125: a repository whose name only starts with the same text: got ["kit-link: docs/tasks/backlog.md links github.com/pharzam/armature (a kit task or note)"]; want []
    checks_test.go:221: kitLink("https://github.com/pharzam/armature") = "github.com/pharzam/armature/", want "github.com/pharzam/armature"
    … (and each case of the core whose message names the repository: the old core prints the text that it gets, with no "/" added)
    verify_test.go:269: a failed add that leaves its directory: <nil>, 15 rows; want a *CleanupError and the whole table
    verify_test.go:280: TMPDIR in the work area: <nil>, calls ["rev-parse …" "show …" "worktree add …"]; want an input error before any scratch tree
    verify_test.go:300: the calls […]; want the worktree add inside the step of pin
    verify_test.go:308: Check(jobs, gates): <nil>, calls […"worktree add …"]; want no scratch tree
$ go test -count=1 -tags=integration -run 'TestEachFindingOfATarget|TestTheInputErrorsOfAWorkArea' ./internal/verify/
    verify_integration_test.go:126: a link to the repository of the baseline (round 1, finding 1): [{"kit-history" "pass" ""}], <nil>; want kit-history fail, …
    verify_integration_test.go:172: TMPDIR in WORK/out: <nil>; want an input error
```

Finding 4: in a copy of the branch, the pin of the overlay of
`frame/bad-missing-linter` names the tree `0000…`. The harness of `b31cc70`
passes on it (the reviewer's run); the fixed harness fails:

```text
    harness_integration_test.go:107: frame/bad-missing-linter: the lines of pin
         got ["setup-check: pin FAIL tree: pin names 0000000000000000000000000000000000000000, root commit has cb89dd90fb277f5d5767688e3e56bb2872814180"]
        want ["setup-check: pin OK"]
```

The first run of round 2, on `dba5ff0`, was stopped by the author after 30 s,
before it wrote a record: the fix of finding 2 failed for a `WORK` given as a
relative path, because `filepath.EvalSymlinks` keeps a relative path relative.
The red run of the new case on the code of `dba5ff0`:

```text
$ go test -count=1 -tags=integration -run TestTheInputErrorsOfAWorkArea ./internal/verify/
    verify_integration_test.go:184: TMPDIR in WORK/out, WORK given as the relative path ../../…/work: <nil>; want an input error
```

## Review round 2: the fixes (cycle 2, O-133)

Review round 2 (`ff76d69`) found the link rule too wide, and two tests with no
traceability row. The red run of the new cases, on the code of `ff76d69` (the
cases of `.git` and of the end of a sentence pass there too: they pin what the
fix keeps):

```text
$ go test -count=1 -run TestKitHistoryFindings ./internal/verify/
    checks_test.go:131: other repositories (round 2, finding 1):
         got ["kit-link: docs/tasks/backlog.md links github.com/pharzam/armature/ (a kit task or note)"]
        want []
```

## Review round 3: the fix (cycle 3, O-134)

Review round 3 (`ff06ff1`) found that a `.` before the repository text still
counted as a boundary, so another host that ends with the baseline's host
gave a link. The red run of the new case, on the code of `ff06ff1`:

```text
$ go test -count=1 -run TestKitHistoryFindings ./internal/verify/
    checks_test.go:133: another host that ends with the host of the baseline (round 3, finding 1):
         got ["kit-link: docs/tasks/backlog.md links github.com/pharzam/armature/ (a kit task or note)"]
        want []
```

A probe of the final expression (a small Go program outside the repository,
recorded once, not a test) gives the wanted result on 22 links to the baseline
and to its neighbours, and on 4 `file://` links.

## The green runs

On the tree of the commit that adds this file, and again on the heads of
cycles 1, 2 and 3.

| Command | Result |
| ------- | ------ |
| `go build ./...`; `go vet` with no tag, `-tags=integration` and `-tags=e2e`; `gofmt -l .` | exit 0; no file |
| `go test -count=1 ./...` | `ok` × 10 packages |
| `go test -count=1 -tags=integration ./...` | `ok` × 10; the harness, the frame on stand-in work areas, the three blocks, `TestPackageRules` (with `internal/work` and `internal/standin`) and `TestInputRule` pass |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | `ok` × 10 |
| `go test -race -count=1` on `internal/verify`, `internal/cli`, `internal/standin`, `internal/work`; and with `-tags=integration` on `internal/verify`, `internal/standin` | `ok` |
