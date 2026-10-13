# T-56dj — test runs

The commands of S13 (#275, row 52c): `commands.sh` starts with `set -eu` in both phases, and in
phase 2 the commands of S13 read the head, push once and apply each ruleset when its output file is
missing, in `internal/setup`. All times are UTC, 2026-10-13, read from the clock after each step.

## The order of the gate

The plan was posted at 09:22:52Z, its review (GPT-6.1 Sol on Codex CLI, effort `high`) ran from
09:22:57Z to 09:28:16Z and was posted at 09:30:39Z, and the author's answer at 09:30:54Z. The tests
were written after that, and run red at 09:33:15Z.

## The base (09:22:46Z to 09:23:18Z)

At `38ee78e`, each with exit 0 (`go test` of `internal/setup` and `internal/cli`, then with
`-tags=integration`):

```
go build ./...: exit 0 (09:22:52Z)
go vet ./...: exit 0 (09:22:54Z)
go test ./internal/setup/ ./internal/cli/: exit 0 (09:23:01Z)
go test -tags=integration ./internal/setup/ ./internal/cli/: exit 0 (09:23:18Z)
```

## Red (09:33:15Z to 09:33:24Z)

The tests before the code: `go vet` of `internal/setup` with and without the tag `integration`
stops at the first error, and `go test` lists the others, as the package does not build.

```
== go vet ./internal/setup/
# github.com/pharzam/layup/internal/setup
# [github.com/pharzam/layup/internal/setup]
vet: internal/setup/commands_test.go:66:4: p.Host undefined (type Phase2 has no field or method Host)
== go vet -tags=integration ./internal/setup/
# github.com/pharzam/layup/internal/setup
# [github.com/pharzam/layup/internal/setup]
vet: internal/setup/commands_test.go:66:4: p.Host undefined (type Phase2 has no field or method Host)
== go test ./internal/setup/ (the compile errors)
# github.com/pharzam/layup/internal/setup [github.com/pharzam/layup/internal/setup.test]
internal/setup/commands_test.go:66:4: p.Host undefined (type Phase2 has no field or method Host)
internal/setup/commands_test.go:66:12: p.DefaultBranch undefined (type Phase2 has no field or method DefaultBranch)
internal/setup/commands_test.go:160:76: undefined: commandsText
FAIL	github.com/pharzam/layup/internal/setup [build failed]
FAIL
```

Two more red steps came after the first green, and are dated by the order of this record:

- **The two tests that the new rule makes false** were changed after the code, each failing as it
  should: `TestTheCommandsFile` (the file now starts with `set -eu`) and the test of S13 in phase 2
  of `T-9v92`, which said that its commands were phase 1's and is now
  `TestPhase2S13KeepsItsActorAndHandOff`.
- **The inputs of the commands** (a host directory that is not an absolute path, an empty default
  branch) were not in the plan; the test `TestPhase2S13RefusesInputsThatMakeNoCommands` was written
  first and failed for the right reason (before the second run of the mutations, at 09:37:55Z), each case with `done` where `fail` was wanted:

```
--- FAIL: TestPhase2S13RefusesInputsThatMakeNoCommands (0.00s)
    commands_test.go:427: a host directory that is relative: done "the ruleset files, and their commands in commands.sh", 3 files; want fail "S13: the host directory of the commands, host/dir, is not an absolute path" and no file
    commands_test.go:427: no host directory: done "the ruleset files, and their commands in commands.sh", 3 files; want fail "S13: the host directory of the commands, , is not an absolute path" and no file
    commands_test.go:427: no default branch: done "the ruleset files, and their commands in commands.sh", 3 files; want fail "S13: the default branch of the forge is empty" and no file
```

## Mutations

Each rule of the writer, of the commands of S13 in phase 2, of the checks of its inputs and of the
wiring mutated alone, on a copy of the tree and never on the worktree: a condition inverted, a
name, a flag, a path, a field or an order changed, so each mutant builds; then `go test
./internal/setup/` (the unit level, which runs the shell text with stand-ins of `git` and `gh`
under `sh` and under `bash`). A mutant that does not build would be counted as not run, never as
caught.

- **The first run** (09:38:16Z to 09:41:23Z): 41 mutants, 39 caught, **2 survived**:
  `message-branch-quoting`, because a double-quoted word of `sh` prints a space and a single quote
  as they are, so the branch of the test (`we ird'x`) needed no quoting (and a space is not valid in
  a branch name of git); and `unsafe-chars`, the same branch. The test now has the branches `a$b`
  and `a$b"c`, which are valid in git and are read by `sh`.
- **The second run** (09:42:52Z to 09:46:34Z): 41 mutants, 40 caught, **1 survived**,
  `unsafe-chars` (a mutant that adds `$` to the characters that need no quotes): `a$b"c` is
  unsafe by its `"` too, so only a branch with the `$` alone tells the mutant; `a$b` was added.
- **The final run** (09:46:56Z to 09:51:09Z): 41 of 41 caught, each with its first failing
  line (each line cut at 170 characters).

```
== header-set-e-only: caught commands_test.go:191: commands.sh of phase 2:
== header-dropped: caught commands_test.go:191: commands.sh of phase 2:
== empty-with-no-command: caught commands_test.go:216: no done step has a command: commands.sh is "set -eu\n"; want it empty
== order-reversed: caught commands_test.go:191: commands.sh of phase 2:
== read-remote: caught commands_test.go:191: commands.sh of phase 2:
== read-ref-main: caught commands_test.go:397: sh, a plain name: exit 0, the read "git -C ~/host/work/acme/widget/target ls-remote origin refs/heads/main", the push "git -C ~/host/work/acme/widge
== head-field: caught commands_test.go:191: commands.sh of phase 2:
== head-read-in-one-pipe: caught commands_test.go:191: commands.sh of phase 2:
== root-read: caught commands_test.go:191: commands.sh of phase 2:
== setup-read: caught commands_test.go:191: commands.sh of phase 2:
== push-when-root: caught commands_test.go:191: commands.sh of phase 2:
== stop-when-setup: caught commands_test.go:191: commands.sh of phase 2:
== push-refspec: caught commands_test.go:397: sh, a plain name: exit 0, the read "git -C ~/host/work/acme/widget/target ls-remote origin refs/heads/trunk", the push "git -C ~/host/work/acme/widg
== push-remote: caught commands_test.go:191: commands.sh of phase 2:
== message-text: caught commands_test.go:191: commands.sh of phase 2:
== message-head: caught commands_test.go:191: commands.sh of phase 2:
== message-stderr: caught commands_test.go:191: commands.sh of phase 2:
== message-branch-quoting: caught commands_test.go:402: sh, a name with a dollar: exit 2, stderr "/home/layup/.cache/layup-session/tmp-T-56dj/TestTheCommandsNameTheDefaultBranchOfTheForge124467133/006/com
== exit-status: caught commands_test.go:191: commands.sh of phase 2:
== host-absolute: caught commands_test.go:427: a host directory that is relative: done "the ruleset files, and their commands in commands.sh", 3 files; want fail "S13: the host directory of the c
== host-reason: caught commands_test.go:427: a host directory that is relative: fail "S13: the host directory is not an absolute path", 0 files; want fail "S13: the host directory of the comman
== branch-empty: caught commands_test.go:427: no default branch: done "the ruleset files, and their commands in commands.sh", 3 files; want fail "S13: the default branch of the forge is empty" a
== branch-reason: caught commands_test.go:427: no default branch: fail "S13: no default branch", 0 files; want fail "S13: the default branch of the forge is empty" and no file
== unsafe-chars: caught commands_test.go:397: sh, a name with a dollar: exit 2, the read "git -C ~/host/work/acme/widget/target ls-remote origin refs/heads/a$b", the push ""; want the branch "a$
== applies-mkdir-p: caught commands_test.go:191: commands.sh of phase 2:
== apply-guard: caught commands_test.go:191: commands.sh of phase 2:
== apply-method: caught commands_test.go:191: commands.sh of phase 2:
== apply-repo-path: caught commands_test.go:191: commands.sh of phase 2:
== apply-input-swapped: caught commands_test.go:191: commands.sh of phase 2:
== apply-records-input: caught commands_test.go:191: commands.sh of phase 2:
== apply-no-temp-file: caught commands_test.go:191: commands.sh of phase 2:
== apply-no-rename: caught commands_test.go:191: commands.sh of phase 2:
== apply-rename-order: caught commands_test.go:191: commands.sh of phase 2:
== apply-order-swapped: caught commands_test.go:191: commands.sh of phase 2:
== apply-dir: caught commands_test.go:191: commands.sh of phase 2:
== work-dir: caught commands_test.go:191: commands.sh of phase 2:
== target-dir: caught commands_test.go:191: commands.sh of phase 2:
== order-of-head: caught commands_test.go:202: the orders of the commands of S13 in phase 2 are [4 4 4 4]; want [2 4 4 4]
== order-of-mkdir: caught commands_test.go:202: the orders of the commands of S13 in phase 2 are [2 2 4 4]; want [2 4 4 4]
== comment-dropped: caught commands_test.go:191: commands.sh of phase 2:
== wiring-commands: caught commands_test.go:191: commands.sh of phase 2:
41 caught, 0 survived, 0 not run, of 41
```

## Green

A first run (09:54:24Z to 10:04:15Z, on `8a0fed5`) had every line at exit 0 but `go test -tags=e2e ./...`, which
exited 1: `TestTheRecordsAreInTheTargetsGit` (`cmd/layup`, e2e) lists the commands of `commands.sh`
and did not expect the first line `set -eu`. The unit and the integration levels do not run that
test, so only the run of the e2e level found it; `702e47e` changes its expectation. The
branch then took `main` by a merge (`6734a23`, the merge of #307, `T-4c5c`; the conflicts were the
cells of the PRD, kept from both tasks), and the whole run was made again.

**The final run on `ad2201b`** (10:08:52Z to 10:18:54Z), each with exit 0, in this order; the time of a line is when its command ended. The three runs of
`go test` are with `-count=1`, then with `-tags=integration` and with `-tags=e2e`, `TMPDIR` is on the
disk, and `git diff --check` reads the commits against `origin/main`:

```
go build ./...: exit 0 (10:08:54Z)
go vet ./...: exit 0 (10:08:55Z)
gofmt -l internal cmd (empty): exit 0 (10:08:55Z)
go test ./...: exit 0 (10:09:04Z)
go test -tags=integration ./...: exit 0 (10:16:34Z)
go test -tags=e2e ./...: exit 0 (10:18:03Z)
adr-lint: exit 0 (10:18:05Z)
prd-lint: exit 0 (10:18:05Z)
link-lint: exit 0 (10:18:19Z)
run-discipline-tests: exit 0 (10:18:26Z)
nested-checkout-check: exit 0 (10:18:27Z)
setup-check: exit 0 (10:18:54Z)
git diff --check origin/main HEAD: exit 0 (10:18:54Z)
```
