# T-5sgt: the review records

The three review rounds of `T-5sgt` (#82), copied from the issue as they stand there (the verdict of round 2 was edited under O-130).

## Review record — round 1

| Field | Value |
| ----- | ----- |
| Commit reviewed | `a0435a62f780d24d79b365345249ba7c0e72bd2d` |
| Reviewer | Claude Fable 5.1 on Claude Code |
| Lens | correctness and acceptance criteria |
| Briefed on | `.review-in/brief.md`, `.review-in/issue-body.md` (#82), `.review-in/issue-comments.md` (the plan, the plan review, the author's answer), `.review-in/inventory.md`; `docs/plan/README.md` (row 5, K19, K29); the whole diff `f394433..HEAD` (29 files); `docs/spec/gate.md` in full; `docs/spec/README.md` (Commands, The types); `docs/spec/packages.md` (the calls of `internal/git`); `docs/spec/setup.md` (the changed part); `docs/architecture.md` §6, §15 (L-A1, L-B1); `docs/tests/test-levels.md`; `docs/tests/traceability.md`; `docs/prd/PRD-0001-layup.md` §12, §13; `docs/engineering-discipline.md` (Bootstrap mode, Material has a test); `docs/guardrails.md` §2 (the new lesson); `runs/T-5sgt/test-runs.md`; `docs/tasks/T-5sgt.md`; the Go sources and tests of `internal/gate`, `internal/git`, `internal/cli`, `cmd/layup` |
| Barred from | nothing; the comments of #82 are the plan and its record |
| Independence claimed | a fresh session with no memory of the author's work, read-only, on a clone at the reviewed commit; the model (Claude Fable 5.1) differs from the author's (Claude Opus 5.5); no file of the repository was changed; no network |
| Cycle | 0 |
| Verdict | `material` |

### Raw findings

1. **Material.** `internal/gate/scratch.go:104`, `for c, under := range files`: the overlay walks the `config` paths in the order of a Go map, which is random. When two `config` paths overlap and one of them is a file at the base or at the head, one order succeeds and the other fails with `RemoveAll a/b: not a directory`, so two runs on the same input give two verdicts. **Basis (a test):** a base with the file `a` and the manifest row `k active sh "cat a" ./*.go "a a/b"`, a head where `a` is a directory with `a/b`; 40 runs of `layup gate` on the same `REPO`, base and head gave 36 rows `not-active scratch tree: overlay failed` (exit 1) and 4 rows `pass —` (exit 0). That contradicts `docs/spec/gate.md:193-195` (the repeat rule: two runs with the same input print the same bytes) and `NFR-005` item 2. A sorted walk of the paths (shortest first, or a fixed order) gives one result; `TestTheOverlayPutsTheGateFilesOfTheBase` holds no overlapping path, so no test catches it.
2. **Material.** `internal/gate/scratch.go:105` with `docs/spec/gate.md:88-89` ("a path that the base does not have is removed from the scratch tree"): the `config` path `.git` is a value of the type `path` (`docs/spec/README.md:251`; `internal/tsv/types.go:36` accepts it), the base never has it, so the overlay runs `RemoveAll(".git")` on the scratch tree and removes the `.git` file of the linked work tree. The command then runs, `git worktree remove` fails (`validation failed, cannot remove working tree: '.../tree/.git' does not exist`), the exit code is 2, and `REPO` keeps a prunable work-tree record: `git worktree list` in `REPO` shows the scratch path. **Basis (a test):** the manifest row `k active sh true ./*.go .git` on a one-commit repository; the table printed `pass`, then `layup: the scratch tree ... is not removed`, exit 2, and `git worktree list` showed the tree as `prunable`. That contradicts `docs/spec/gate.md:65-66` ("changes nothing in it except a scratch work tree, which it removes before it exits") for an input that the schema accepts; the spec names only a symbolic link and a submodule at the base as refused `config` paths (`gate.md:101-102`). A `config` path whose first part is `.git` is an input error, or the overlay never removes `.git`.
3. **Material.** `docs/prd/PRD-0001-layup.md:242` (the `NFR-005` row of §12) and `docs/prd/PRD-0001-layup.md:257` (the §13 row): the plan's step 4 and the author's answer name "the `PRD-0001` §12 Test cells of `REQ-004`, `REQ-007`, `NFR-004` and `NFR-005`"; the diff fills three cells and leaves `NFR-005` as it was (`TestVersion` for the repeat rule), while `docs/tests/traceability.md:47` maps `TestGateOnAGoRepository` to `NFR-005`, and `docs/tasks/T-5sgt.md:48-49` records three cells with no reason for the fourth. **Basis (a cited clause):** acceptance criterion 4 of #82, "`PRD-0001` §12 names its tests in the Test column", and the contract of step 4. One cell and one word of the §13 row fix it, or the task record says why `NFR-005` stays.
4. **Note.** `internal/gate/run_test.go` (every `Run` test, `TestTheOverlayPutsTheGateFilesOfTheBase`, `TestTheOverlayDoesNotWriteOutOfTheTree`): the tests of the unit level write real files and a real symbolic link under `TMPDIR` (`use` sets it to `t.TempDir()`; `fakeGit.WorktreeAdd` writes the head; `newScratch` calls `os.MkdirTemp` and `os.OpenRoot`). `docs/tests/test-levels.md:43` says "Unit tests touch no file"; the inventory's unit test of `gate-scratch-tree` says "The unit test touches no real file". The tests are deterministic and pass under `-race`; the level label is the point. Record the exception in the traceability row, or move the two overlay tests to the integration file.
5. **Note.** `docs/spec/gate.md:102-104`: "a symbolic link of the head cannot send a write out of the tree; such a write fails the overlay". A link of the head that points inside the tree fails the overlay too: a head with `docs -> realdocs` gave `mkdirat docs: file exists`, `not-active scratch tree: overlay failed`, exit 1. A head where the parent of a `config` path is a file (`cfg` a file, `config` `cfg/x`) gave `RemoveAll cfg/x: not a directory` and the same row. Both results are safe (never `pass`), but the sentence names only the first case; one clause ("a symbolic link or a file of the head at or above a `config` path or `docs` fails the overlay") makes the spec say what the code does.
6. **Note.** `internal/gate/scratch.go:183`: when `sh` cannot start after `lookPath` found it, the outcome is `exit -1`, a reason that `docs/spec/gate.md:143` does not name. Not reached in a test.
7. **Note.** `internal/cli/args.go:58-63` (row 3's frame, not this task's code): a revision that starts with `-` cannot be given as `--base -x` (`flag --base has no value`, exit 2); `--base=-x` reaches `git rev-parse --end-of-options` and gives the input error. Correct and safe; `gate.md:67-69` could say the `=` form for such a revision.
8. **Note.** `docs/spec/gate.md:26`, the column `scope` is `list(text)` and accepts `—`: an `active` kind with no pattern is always `clear` (`no product path`), exit 0. By the table that is correct; the catalog rule of row 14 may want "a kind has at least one pattern".
9. **Note.** `runs/T-5sgt/test-runs.md:116` reports the end-to-end run "in about 3 s with the host's `GOCACHE`"; this clone's run took 2.5 s to 3.0 s. The red runs are plausible and their line numbers agree with the files; no run is contradicted by what this clone shows.

Checked and in agreement (no finding): D1 (eighteen odd patterns, a kind with a digit and an uppercase kind: each an input error with exit 2, or the documented match); D2 (a tag, a tree, a blob, an unknown revision, an empty revision, a base with no manifest, a manifest with no row, a manifest with a carriage return: exit 2, no table); D3 (a `config` path that is a directory at the base and a file at the head, and the reverse: the base's files with their modes; a `config` directory removed at the head: restored; a path absent at the base: removed; `docs` removed at the head: written; the hooks off, proved by the control run of the integration test); the `os.Root` decision (a link out of the tree: the write fails, nothing written outside); D4 (a scratch tree that cannot be removed, by a locked file: the whole table, the diagnostic with the path, exit 2, `REPO` is clean); D5 (K19 in `NFR-004` item 3); D6 (`signal killed`, `signal hangup`, `signal terminated`; a missing second program: `exit 127`, `fail`; a command that reads standard input gets end-of-file, not `layup`'s input; 3 MB of output held and printed; `HOME` of the host reaches the command); D7 (one block per kind under its step line, no beat inside); D8 (`REPO` as a linked work tree of another repository: works; full commit IDs; the frame's exit map); D9; base equal to head; two runs: the same bytes; `gate-result` in `built`.

### Acceptance criteria

- Each inventory item delivered as its specification sections say: **not met** on one point (finding 1, the repeat rule of `gate-run`/`gate-command` for an overlapping `config`; finding 2, `gate-scratch-tree` and `.git`); every other point met.
- The tests of the plan pass (unit, integration, e2e): met.
- K19 settled in `docs/spec/gate.md` with its reason: met.
- Traceability rows `green` and `PRD-0001` §12 names its tests: **not met** for the `NFR-005` cell (finding 3); met for the rows and the other three cells.
- Tests cover the change and pass (R8): met, with the gap of finding 1 (no test of overlapping `config` paths).
- Docs updated in the same PR: met.
- Budget (3,000 lines over 32 files): met, 1,945 lines over 29 files.

### Runs

All with `HOME=/tmp/rv-T-5sgt/home` and the `GOCACHE` of the host; scratch repositories under `/tmp/rv-T-5sgt/`, removed after the runs.

| Command | Result |
| ------- | ------ |
| `git diff --shortstat f394433 HEAD` | `29 files changed, 1912 insertions(+), 33 deletions(-)` |
| `gofmt -l .` | no file |
| `go vet ./...`; `-tags=integration`; `-tags=e2e` | exit 0 each |
| `go test -count=1 ./...` | `ok` × 7 packages |
| `go test -count=1 -tags=integration ./...` | `ok` × 7 |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | `ok` × 7; `TestGateOnAGoRepository` 1.36 s, `TestGateNeverPassesACheckThatDidNotRun` 0.66 s (both ran, `-v`) |
| `go test -count=1 -race ./internal/gate ./internal/cli` | `ok` × 2 |
| `go test -count=1 -race -tags=integration ./internal/gate ./internal/cli` | `ok` × 2 |
| `go build -o /tmp/rv-T-5sgt/layup ./cmd/layup` | exit 0 |
| `config` dir at base, file at head, and the reverse | `pass`, the base's files, exit 0 |
| head: `docs -> /tmp/rv-T-5sgt/outside` | `not-active scratch tree: overlay failed`, exit 1, nothing written outside |
| head: `docs -> realdocs` (inside); `cfg -> cfg2` | `not-active scratch tree: overlay failed`, exit 1 (note 5) |
| head: `cfg` a file, `config` `cfg/x` | `not-active scratch tree: overlay failed`, exit 1 (note 5) |
| 18 odd scope patterns, a kind `k1`, a kind `K`, `config` `../x`, `/etc` | exit 2 with the reason, or the documented match |
| `config` `.git` | `pass`, then exit 2 `not removed`; a prunable work-tree record in `REPO` (finding 2) |
| `config` `a a/b`, `a` a file at base, a directory at head, 40 runs | 36 × `not-active overlay failed` (exit 1), 4 × `pass` (exit 0) (finding 1) |
| 3 MB output; `cat; read x`; `kill -KILL $$`; `kill -HUP $$`; a missing second program; `echo $HOME; pwd` | `pass`; `pass` (end-of-file on standard input); `fail signal killed`; `fail signal hangup`; `fail exit 127`; `pass`, the host's `HOME`; exit 1 |
| a locked file in the scratch tree (`chflags uchg`) | the whole table, `layup: the scratch tree <path> is not removed: …`, exit 2 |
| `--base -q`; `--base=-q`; `--base ''` | exit 2 `flag --base has no value` (usage); exit 2 `the revision "-q" is not a commit`; exit 2 (usage) |
| `--base X --head X` | the active kinds run, the pending kinds `clear`, exit 0 |
| `REPO` a linked work tree of another repository | the same table as in the main tree; the work-tree list unchanged |
| `REPO` not a directory | exit 2, the reason on standard error, no table |
| a manifest with CRLF | exit 2, the reason of `internal/tsv` |
| `layup gate` twice on one input (three repositories) | the same bytes, no scratch path in the table |
| `git worktree list` and `git status --porcelain` in each scratch `REPO` after each run | only the main tree; unchanged (except finding 2) |

*Posted for a reviewer session (Claude Fable 5.1 on Claude Code, `claude -p`, a fresh read-only session in a clone at `a0435a6`, from 11:24 +03, 7 min 52 s) through the `layup-agent` App. The text is the session's file, unchanged.*

## Review record — round 2

| Field | Value |
| ----- | ----- |
| Commit reviewed | `bab6759de0a26f29b9721107e2b5d82a9d14c9d0` |
| Reviewer | Claude Fable 5.1 on Claude Code |
| Lens | the fixes of round 1, and correctness of what they touch |
| Briefed on | `.review-in/brief.md`, `.review-in/issue-body.md` (#82), `.review-in/contract.md` (the plan, the plan review, the author's answer), `.review-in/round-1.md`, `.review-in/inventory.md`; the fix diff `a0435a6..HEAD` in full (9 files); the whole diff `f394433..HEAD` where the fix touches it; `docs/spec/gate.md` (the manifest, the command, the run, the table); `docs/spec/README.md` (Commands); `docs/tests/test-levels.md`; `docs/tests/traceability.md`; `docs/prd/PRD-0001-layup.md` §12, §13; `runs/T-5sgt/test-runs.md`; `docs/tasks/T-5sgt.md`; `internal/gate/scratch.go`, `manifest.go`, `manifest_test.go`, `run_test.go`; `internal/tsv/types.go` (the type `path`); the `kinds.tsv` files of the repository |
| Barred from | the comments of #82 after the author's answer |
| Independence claimed | a fresh session with no memory of the author's work or of round 1, read-only, on a clone at the reviewed commit; the model (Claude Fable 5.1) differs from the author's (Claude Opus 5.5); no file of the repository was changed; no `gh`, no `git push`, no `git fetch`, no network; scratch files under `/tmp/rv2-T-5sgt/` only, with `HOME=/tmp/rv2-T-5sgt/home` |
| Cycle | 1 |
| Verdict | `material` (first posted with the last-round verdict, because it was the last round under the cap of that time; O-130 raised the cap to 2 and a fix followed, so this is an intermediate round; edited by the author under O-130) |

### Round-1 findings

1. **Fixed.** `internal/gate/scratch.go:108` walks the `config` paths in `slices.Sorted(maps.Keys(files))`, and `scratch.go:145-150` treats `ENOTDIR` of `RemoveAll` as "absent". The test of round 1 on a real repository (`config` `a a/b`, `a` a file at the base and a directory at the head): 30 runs, 30 × `pass —`, exit 0 each, the same bytes. The neighbours hold: `a/b a` in the manifest (20 runs, 20 × `pass`); `a a/b/c` (10 × `pass`); `a` a directory at the base and a file at the head (10 × `pass`, the base's `a/b`); `a/b a` with `a/b/c` at the base and an extra `a/x` at the head (10 × `pass`, `a/x` removed, the base's `a/b/c`). `TestOverlappingConfigPathsGiveOneResult` (`run_test.go:342`) runs the case 20 times and passes under `-race`; `runs/T-5sgt/test-runs.md:107-120` shows its red run on `a0435a6`.
2. **Fixed for the spelling `.git`; a neighbour stays open (new finding 1).** `internal/gate/manifest.go:69-73` refuses a `config` path equal to `.git` or with the prefix `.git/`. On a real repository: `.git` and `.git/x` give exit 2, no table, the reason `line 2, column "config": .git is in .git` on standard error, and `git worktree list` in `REPO` shows only the main tree; `./.git` is refused by the type `path` (`internal/tsv/types.go:36`, `fs.ValidPath`); `.gitignore` and `.gitx` are not refused and give `pass`; `a .git` gives exit 2. `TestReadManifestRefusesEachMalformedForm` holds `.git` and `layout .git/hooks/x`. But `.GIT` on this case-insensitive filesystem repeats the symptom of the finding: see new finding 1.
3. **Fixed.** `docs/prd/PRD-0001-layup.md:241` names `TestGateOnAGoRepository (cmd/layup, e2e: two runs of layup gate give the same bytes)` in the `NFR-005` cell; `PRD-0001-layup.md:257` (the §13 row) names the four cells; `docs/tests/traceability.md:47` agrees (`TestGateOnAGoRepository` → `NFR-005`); `docs/tasks/T-5sgt.md:60-65` records the fix.
4. **Applied.** `internal/gate/run_test.go:1` has `//go:build integration` with the reason in a comment; `go test -count=1 ./...` (no tag) still builds `internal/gate` (the untagged tests `scope_test.go`, `manifest_test.go`, `result_test.go` hold their own helpers; `manifestHeader` is in `manifest_test.go:8`); `docs/tests/traceability.md:51-52` has the unit row without `run_test.go` and a new integration row for it, with the reason "it writes temporary files"; `runs/T-5sgt/test-runs.md:131` says that the untagged run no longer holds these tests.
5. **Applied, with one sentence now too wide (new finding 2).** `docs/spec/gate.md:108-110` names "a symbolic link or a file of the head at or above a `config` path or `docs`". The fix of finding 1 changed the behaviour of one of its cases: see new finding 2.
6. **Applied.** `docs/spec/gate.md:126-127`: "`sh` that is found and then does not start gives `exit -1`"; `internal/gate/scratch.go:196` gives `outcome{code: -1}`, which the table's line "exits with another code" reads as `fail`, `exit -1`.
7. **Applied.** `docs/spec/gate.md:73-75` gives the form `--base=-x` and links to `README.md` (Commands), whose lines 57-62 define `--name=VALUE` and "a word that starts with `-`" as no value.
8. **Applied.** `docs/spec/gate.md:50-52`: a kind with no scope pattern is an input error; `internal/gate/manifest.go:66-68` gives `line N, column "scope": the kind K has no scope pattern`; on a real repository, scope `—` gives exit 2 and no table. No golden, fixture or test of the repository holds a manifest with an empty scope: the one `kinds.tsv` (`internal/catalog/testdata/test/kinds.tsv`) has `./*.txt` on each row; every manifest in `internal/gate/gate_integration_test.go`, `internal/gate/run_test.go` and `cmd/layup/gate_e2e_test.go` has a pattern; the `—` of `internal/tsv/tsv_test.go:187` is in a `config` column. All three test levels pass.
9. **Not a change; agreed.** `TestGateOnAGoRepository` took 1.47 s and `TestGateNeverPassesACheckThatDidNotRun` 0.66 s in this clone with the host's `GOCACHE`.

### Raw findings

1. **Material.** `internal/gate/manifest.go:70`, `if c == ".git" || strings.HasPrefix(c, ".git/")`, with `docs/spec/gate.md:106-108` ("a `config` path that is `.git` or under it, is an input error") and `gate.md:65-66` ("changes nothing in it except a scratch work tree, which it removes before it exits"). The check compares strings; the filesystem of the scratch tree compares names. On a case-insensitive filesystem (the default of macOS, where this clone and the demo run; also Windows) the `config` path `.GIT` is the file `.git` of the scratch tree, the type `path` accepts it (`internal/tsv/types.go:36` allows upper case), and the overlay runs `RemoveAll(".GIT")`. **Basis (a test):** the manifest row `k active sh true ./*.go .GIT` on a one-commit repository, two runs: each printed the table with `pass —`, then `layup: the scratch tree /var/folders/…/layup-gate-…: is not removed: … fatal: validation failed, cannot remove working tree: '…/tree/.git' does not exist`, exit 2; `git worktree list` in `REPO` then shows two `prunable` records. That is the symptom of round-1 finding 2 (`pass` printed, exit 2, `REPO` changed) under a spelling the fix does not cover, so the fix of finding 2 does not hold on the operator's platform. One line fixes it: refuse a first path segment that `strings.EqualFold`s `.git` (and say "in any case" in `gate.md:107`); or, in `overlay`, refuse to remove a path whose `Lstat` through the root is `os.SameFile` with `.git`, which covers each alias a filesystem can make.
2. **Note.** `docs/spec/gate.md:108-110`, "a symbolic link or a file of the head at or above a `config` path or `docs` fails the overlay", beside `gate.md:110-112`, "a path under a file of the tree counts as absent". The two sentences disagree for one input, and the code follows the second: `config` `a/b`, the head has the file `a`, the base has no `a/b`: `pass`, exit 0, the head's `a` kept (round 1, note 5, saw `overlay failed` for this input on `a0435a6`; the `ENOTDIR` tolerance of the fix of finding 1 changed it). With `a/b` at the base, the same head gives `scratch tree: overlay failed`, exit 1, as the sentence says. The result is safe and deterministic, and correct by step 2 of the run (the base has no gate file at `a/b`). One clause makes the sentence true: "fails a write of the overlay". A second input of the same sentence: the head has `a -> d` (a link inside the tree) and `d/b`, `config` `a/b`, the base has no `a/b`: `RemoveAll("a/b")` follows the link and removes `d/b`, and the kind gives `fail exit 1` for `test -f d/b` where the sentence says `not-active`. The exit code is 1 either way, the behaviour precedes the fix (`os.Root` follows a link inside the root), and the input is a head that aliases a `config` parent; recorded here for the sentence, not for the verdict.
3. **Note.** `docs/spec/gate.md:50-52`: a kind with no scope pattern is now an input error of `layup gate`, but the catalog reader (`internal/catalog/catalog.go:34,49,101`) still accepts a `kinds.tsv` row with scope `—` and would write it into a manifest that `layup gate` refuses with exit 2. No entry of the repository has such a row. The sentence could end as the `tool` sentence does (`gate.md:55-58`): "a rule for the catalog entries, task `T-c06a`, row 14", so the check moves to the source once.

### Acceptance criteria

- Each inventory item delivered as its specification sections say: **not met** on one point (new finding 1, `gate-scratch-tree` and `.GIT` on a case-insensitive filesystem); the three findings of round 1 are fixed for their own inputs.
- The tests of the plan pass (unit, integration, e2e): met.
- K19 settled in `docs/spec/gate.md`: met (unchanged since round 1).
- Traceability rows `green` and `PRD-0001` §12 names its tests: met (finding 3 of round 1 fixed).
- Tests cover the change and pass (R8): met; the new tests of the fixes ran red before them (`test-runs.md:107-120`).
- Docs updated in the same PR: met, with note 2 (one sentence too wide).
- Budget (3,000 lines over 32 files): met, 2,052 lines over 29 files.

### Runs

All with `HOME=/tmp/rv2-T-5sgt/home`, `GOTOOLCHAIN=local` and the `GOCACHE` of the host; scratch repositories under `/tmp/rv2-T-5sgt/scen/`, with a fixed identity.

| Command | Result |
| ------- | ------ |
| `git diff --shortstat f394433 HEAD` | `29 files changed, 2018 insertions(+), 34 deletions(-)` |
| `git diff a0435a6 HEAD --stat` | 9 files, 124 insertions, 19 deletions |
| `go build ./...`; `go build -o … ./cmd/layup` | exit 0 each |
| `gofmt -l .` | no file |
| `go vet ./...`; `-tags=integration`; `-tags=e2e` | exit 0 each |
| `go test -count=1 ./...` | `ok` × 7 packages |
| `go test -count=1 -tags=integration ./...` | `ok` × 7 |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | `ok` × 7; `TestGateOnAGoRepository` 1.47 s, `TestGateNeverPassesACheckThatDidNotRun` 0.66 s (`-v`) |
| `go test -count=1 -race -tags=integration ./internal/gate ./internal/cli` | `ok` × 2 |
| `go test -count=1 -race ./internal/gate ./internal/cli` | `ok` × 2 |
| `go test -count=1 -tags=integration -v -run 'TestOverlappingConfigPathsGiveOneResult\|TestReadManifestRefusesEachMalformedForm' ./internal/gate` | both `PASS` |
| `sh docs/adr/adr-lint.sh` | `adr-lint: OK` |
| `sh docs/prd/prd-lint.sh` | `prd-lint: OK` |
| `sh docs/links/link-lint.sh` | `link-lint: OK  1314 links resolved` |
| `sh docs/tests/run-discipline-tests.sh` | `81 passed, 0 failed` |
| `sh docs/setup/setup-check.sh` | each check `OK` |
| `git diff --check f394433 HEAD` | exit 0 |
| `config` `a a/b`, `a` a file at the base and a directory at the head, 30 runs | 30 × `pass —`, exit 0 |
| `config` `a/b a` (the reverse order), 20 runs | 20 × `pass —`, exit 0 |
| `config` `a a/b/c`, `a` a file at the base, 10 runs | 10 × `pass —`, exit 0 |
| `config` `a a/b`, `a` a directory at the base and a file at the head, 10 runs | 10 × `pass —`, the base's `a/b`, exit 0 |
| `config` `a/b a`, `a/b/c` at the base, `a/x` added at the head, 10 runs | 10 × `pass —`, `a/x` removed, exit 0 |
| `config` `a/b`, the head has the file `a`, no `a/b` at the base | `pass —`, exit 0 (note 2) |
| `config` `a/b`, the head has the file `a`, `a/b` at the base | `not-active scratch tree: overlay failed`, exit 1 |
| `config` `a/b`, the head has `a -> d` and `d/b`, no `a/b` at the base | `fail exit 1`: `d/b` removed through the link (note 2) |
| `config` `.git`; `.git/x`; `a .git` | exit 2, no table, `… is in .git`; `REPO` unchanged |
| `config` `./.git` | exit 2, `"./.git" is not a value of list(path)` |
| `config` `.gitignore`; `.gitx` | `pass —`, exit 0 |
| `config` `.GIT`, two runs, on a case-insensitive filesystem | `pass —` printed, then `the scratch tree … is not removed`, exit 2; two `prunable` records in `REPO` (finding 1) |
| scope `—` | exit 2, `the kind k has no scope pattern`, no table |
| `git worktree list` and `git status --porcelain` in each scratch `REPO` after each run | only the main tree, clean (except `.GIT`) |

*Posted for a reviewer session (Claude Fable 5.1 on Claude Code, `claude -p`, a fresh read-only session in a clone at `bab6759`, from 11:36 +03, 6 min 52 s) through the `layup-agent` App. The text is the session's file, unchanged.*

## Review record — round 3

| Field | Value |
| ----- | ----- |
| Commit reviewed | `51755cdf9b160382d184cf8fcdec23d615d60282` |
| Reviewer | Claude Fable 5.1 on Claude Code |
| Lens | an adversarial check of the scratch tree and the overlay; the fix of round 2 |
| Briefed on | `.review-in/brief.md`, `.review-in/issue-body.md` (#82), `.review-in/contract.md` (the plan, the plan review, the author's answer), `.review-in/round-1.md`, `.review-in/round-2.md`, `.review-in/inventory.md`; the fix diff `bab6759..HEAD` in full (8 files); the diff `a0435a6..HEAD` (10 files); `docs/spec/gate.md` (the manifest, the command, the run, the decided rules, `REQ-007`); `docs/tests/traceability.md` (the rows of `T-5sgt`); `runs/T-5sgt/test-runs.md`; `docs/tasks/T-5sgt.md`; `internal/gate/scratch.go`, `manifest.go` in full; `internal/git/git.go` (`LsTree`, `Show`, `WorktreeAdd`, `WorktreeRemove`, `environ`); `internal/tsv/types.go` (the type `path`) |
| Barred from | the comments of #82 after the author's answer |
| Independence claimed | a fresh session with no memory of the author's work or of rounds 1 and 2, read-only, on a clone at the reviewed commit; the model (Claude Fable 5.1) differs from the author's (Claude Opus 5.5); no file of the repository was changed (`git status --porcelain` is empty after the runs); no `gh`, no `git push`, no `git fetch`, no network; scratch files under `/tmp/rv3-T-5sgt/` only, with `HOME=/tmp/rv3-T-5sgt/home` for every command |
| Cycle | 2 |
| Verdict | `nothing material in scope` |

### Earlier findings

**Round 2, finding 1 (`.GIT` on a file system that folds case): fixed.** Two parts. (a) `internal/gate/manifest.go:71` refuses a `config` path whose first segment `strings.EqualFold`s `.git`. On a real repository on this host (macOS, APFS, case-insensitive): `.GIT`, `.Git/x`, `.gIt` and `.git` each give exit 2, no table, the reason `line 2, column "config": <path> is in .git` on standard error; `git worktree list --porcelain` in `REPO` shows one work tree (the main tree) and `git status --porcelain` is empty after each run. (b) `internal/gate/scratch.go:148-152`: `removeAll` refuses a path whose `Lstat` through the root is `os.SameFile` with `.git`. The guard holds on each path the overlay removes or writes, because `write` (`scratch.go:134`) calls `removeAll` before `WriteFile`, and `overlay` (`scratch.go:109`) calls it for each `config` path. A head cannot give a hard link (git checks out no hard link), so the real alias on this host is a path that resolves to `.git` through a symbolic link of the head: `d -> .` with `config` `d/.GIT`, and with `config` `d/.git` (the first segment `d` passes the manifest check): each gives `not-active scratch tree: overlay failed`, exit 1, the diagnostic `the scratch tree: d/.GIT is the file .git of the scratch tree`, one work tree in `REPO`, `REPO` clean. The hard-link case is `TestTheOverlayNeverRemovesTheGitFileOfTheTree` (`run_test.go:374`) with the fake, which passes under `-race`. The guard refuses nothing it should accept: `config` `.gitignore .github/x git`, each a file at the base and changed at the head, gives `pass`, exit 0, and the command printed the base's content of all three; a head with the symbolic link `alias -> .git` and `config` `alias` (a file at the base) gives `pass`, the base's `alias` content, and `.git` of the tree stays (`test -f .git` passed in the command); `TestReadManifestTakesAConfigPathThatOnlyStartsWithGit` (`manifest_test.go:49`) holds `.gitignore .github/x.yml`. The integration test `TestTheRevisionsAndTheManifestOnARealRepository` (`gate_integration_test.go:182-188`) holds `.GIT` and asserts one work tree; `runs/T-5sgt/test-runs.md:130-153` shows its red run on `bab6759`.

**Round 2, note 2 (the sentence about a symbolic link of the head): applied, and the sentence is still too wide for one input (new finding 1).** `docs/spec/gate.md:111-114` now says "a removal follows a symbolic link that stays inside the tree"; `config` `d/x` with `d -> .` at the head (`d/x` a file at the base) gives `pass` and the base's `d/x` content, 5 runs, the same bytes, exit 0.

**Round 2, note 3 (a kind with no scope pattern, a rule for the catalog): applied.** `docs/spec/gate.md:50-52` ends "a rule for the catalog entries too (task `T-c06a`, row 14 of the plan)". Scope `—` on a real repository: exit 2, `the kind k has no scope pattern`, no table.

**Round 1, finding 1 (the order of the `config` paths): still fixed.** `scratch.go:108` walks `slices.Sorted(maps.Keys(files))`; `config` `a a/b`, `a` a file at the base and a directory at the head: 20 runs, 20 × `pass —`, exit 0, the same bytes. `TestOverlappingConfigPathsGiveOneResult` passes under `-race`.

**Round 1, finding 2 (`.git` as a `config` path): still fixed.** `.git` in lower case: exit 2, no table, `REPO` unchanged (above).

**Round 1, finding 3 (the `NFR-005` cell): still fixed.** `docs/prd/PRD-0001-layup.md` names `TestGateOnAGoRepository` in the `NFR-005` cell and the §13 row names four cells (unchanged since round 2); `docs/tests/traceability.md:47` agrees.

### Raw findings

1. **Note.** `docs/spec/gate.md:111-114`, "a symbolic link or a file of the head at or above a path that the overlay writes fails the overlay, and a removal follows a symbolic link that stays inside the tree". A write follows such a link too, so the first clause is false for a link of the head that points to a directory inside the tree. **Basis (tests on a real repository):** the head has `docs -> real` with `real/keep` a file: the overlay writes `docs/gates.tsv` through the link, the command listed `real` as `gates.tsv keep`, `pass`, exit 0, where the sentence says `overlay failed`. The head has `docs -> nowhere` (a dangling link): `mkdirat docs: file exists`, `not-active scratch tree: overlay failed`, exit 1, as the sentence says. A second input of the same mechanism: `config` `a/x b/x`, both files at the base (`A` and `B`), the head has `a -> b`: the sorted walk writes `a/x` through the link, then removes and writes `b/x`, so the command read `B` at `a/x` and at `b/x`, `pass`, exit 0: a head link that aliases two `config` paths makes one base gate file stand for another. Each result is safe (no write out of the tree, `REPO` clean, one work tree), deterministic (the same bytes on each run), and inside the trust of `FT4` (the file that judges is a base file). The same class as round 2, note 2, which rated this sentence a note; the author applied it in this cycle, and the clause for a write stays wide. One clause makes the sentence true: "a symbolic link of the head that points out of the tree or to no file, or a file of the head, at or above a path that the overlay writes fails the overlay; a link that points to a directory inside the tree is followed by a removal and by a write". A stricter overlay (an `Lstat` of each parent of a `config` path and of `docs`, a link refused) would close the second input, a choice for the author or a later task, not a condition of this verdict.
2. **Note.** `internal/gate/scratch.go:153` with `docs/spec/gate.md:88-89` ("a path that the base does not have is removed from the scratch tree"): the removal goes by name, and on a file system that folds case a `config` path that differs from a head path only by case removes that head file. **Basis (tests on a real repository, this host):** `config` `X.GO`, the head has `x.go`, scope `./*.go`, command `test -f x.go`: `clear no product path`, exit 0, so the kind did not run on `x.go`; on a file system that does not fold case the same input gives `pass` (the control `config` `x.go` gives `pass` here). `config` `DOCS`, command `test -f docs/gates.tsv`: `fail exit 1`, because the removal took `docs` with the manifest the overlay had written. Not material: the manifest is the base's own, each result is one result on one host (the repeat rule of `gate.md:200-205` is for one `REPO`), no sentence of `gate.md` is contradicted, `REPO` is unchanged, and the `.git` guard of the fix is the one name the spec protects. A sentence beside the `.GIT` sentence of `gate.md:115-117` ("on such a file system a `config` path is also the head's file of that name in any case") would warn the author of a manifest.

Checked and in agreement (no finding): `git ls-tree -r --full-tree -- PATH` takes `PATH` as a literal prefix, not a pattern, so `config` `*`, `x*` and `[x].go` match nothing at the base, remove nothing at the head (`os.Root.RemoveAll` is literal) and give the head's `x.go` unchanged (`grep -q head x.go` passed), where the control `config` `x.go` wrote the base's file (`fail exit 1`); a `config` path `a/.git` or `d/.GIT` with no link at the head: nothing at the base, nothing at the head, no removal; `.git` of the scratch tree is a file (`git worktree add`), so `files()` (`scratch.go:166-169`) skips the one name; `environ()` (`git.go:64-73`) keeps `TMPDIR` and `PATH` only, so a `git` of `REPO`'s own configuration cannot reach the calls; `WorktreeRemove` runs on every exit path of the run and the directory of `os.MkdirTemp` is removed after it; the git version is 2.54.0; after every scenario `git worktree list --porcelain` in `REPO` shows one tree and `git status --porcelain` is empty.

### Acceptance criteria

- Each inventory item delivered as its specification sections say: met; round 2's finding 1 is fixed for `.GIT`, `.Git/x`, `.gIt`, `.git`, and for a path that resolves to `.git` through a link of the head.
- The tests of the plan pass (unit, integration, e2e): met.
- K19 settled in `docs/spec/gate.md`: met (unchanged since round 1).
- Traceability rows `green` and `PRD-0001` §12 names its tests: met.
- Tests cover the change and pass (R8): met; the tests of the fix ran red on `bab6759` (`test-runs.md:130-153`) and pass here, with `-race`.
- Docs updated in the same PR: met, with notes 1 and 2 (two sentences of `gate.md`).
- Budget (3,000 lines over 32 files): met, 2,152 lines over 29 files.

### Runs

All with `HOME=/tmp/rv3-T-5sgt/home` and the `GOCACHE` of the host (`/Users/farzam/Library/Caches/go-build`); scratch repositories under `/tmp/rv3-T-5sgt/scen/`, with a fixed identity; the binary `go build -o /tmp/rv3-T-5sgt/layup ./cmd/layup`.

| Command | Result |
| ------- | ------ |
| `git diff --shortstat f394433 HEAD` | `29 files changed, 2118 insertions(+), 34 deletions(-)` (2,152 lines over 29 files) |
| `git diff --stat bab6759 HEAD`; `git diff --stat a0435a6 HEAD` | 8 files, 114 insertions, 14 deletions; 10 files, 225 insertions, 20 deletions |
| `go build -o /tmp/rv3-T-5sgt/layup ./cmd/layup` | exit 0 |
| `go test -count=1 ./...` | `ok` × 7 packages, exit 0 |
| `go test -count=1 -tags=integration ./...` | `ok` × 7, exit 0 (`internal/gate` 2.5 s) |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | `ok` × 7, exit 0 (`cmd/layup` 2.7 s) |
| `go test -count=1 -race -tags=integration ./internal/gate ./internal/cli` | `ok` × 2, exit 0 |
| `go vet ./...`; `-tags=integration`; `-tags=e2e` | exit 0 each |
| `gofmt -l .` | no file, exit 0 |
| `sh docs/adr/adr-lint.sh` | `adr-lint: OK` |
| `sh docs/prd/prd-lint.sh` | `prd-lint: OK` |
| `sh docs/links/link-lint.sh` | `link-lint: OK  1314 links resolved` |
| `sh docs/tests/run-discipline-tests.sh` | `81 passed, 0 failed` |
| `sh docs/setup/setup-check.sh` | each check `OK` |
| `git diff --check f394433 HEAD` | exit 0 |
| `config` `.GIT`; `.Git/x`; `.gIt`; `.git` | exit 2 each, no table, `<path> is in .git`; one work tree in `REPO`, clean |
| `config` `.gitignore .github/x git`, each changed at the head | `pass —`, exit 0, the base's three files |
| scope `—` | exit 2, `the kind k has no scope pattern`, no table |
| `config` `a a/b`, `a` a file at the base and a directory at the head, 20 runs | 20 × `pass —`, exit 0, the same bytes |
| head `alias -> .git`, `config` `alias` (a file at the base) | `pass —`, exit 0, `.git` of the tree stays, `REPO` clean |
| head `d -> .`, `config` `d/.GIT`; `d/.git` | `not-active scratch tree: overlay failed`, exit 1, `… is the file .git of the scratch tree`; one work tree, `REPO` clean |
| head `d -> .`, `config` `d/x` (a file at the base), 5 runs | 5 × `pass —`, exit 0, the base's content through the link |
| head `docs -> .`; `docs -> real` (an existing directory) | `pass —`, exit 0, `gates.tsv` written through the link (note 1) |
| head `docs -> nowhere` (a dangling link) | `not-active scratch tree: overlay failed`, exit 1, `mkdirat docs: file exists` |
| head `a -> b`, `config` `a/x b/x` (both files at the base) | `pass —`, exit 0, `B` at `a/x` and `b/x` (note 1) |
| `config` `*`; `x*`; `[x].go`, the head's `x.go` changed | `pass —`, exit 0 each: the head's `x.go` kept (a literal prefix in `git ls-tree`); the control `config` `x.go`: `fail exit 1`, the base's file |
| `config` `X.GO` (a case alias of the head's `x.go`) | `clear no product path`, exit 0 (note 2) |
| `config` `DOCS` (a case alias of `docs`) | `fail exit 1`: `docs/gates.tsv` gone from the tree (note 2) |
| `git worktree list --porcelain` and `git status --porcelain` in each scratch `REPO` after each run | one tree; empty |
| `git status --porcelain` in this clone after the review | empty |

*Posted for a reviewer session (Claude Fable 5.1 on Claude Code, `claude -p`, a fresh read-only session in a clone at `51755cd`, from 11:49 +03, 7 min 41 s) through the `layup-agent` App. The text is the session's file, unchanged.*
