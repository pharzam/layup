# The packages

The Go packages of `layup`, their jobs, and the import rules between them. The
boundary gate of a Go target is "configured from the package table of the
approved architecture" ([`architecture.md`](../architecture.md) §6); this is
that table for LAYUP itself. Conventions: [`README.md`](README.md).

## NFR-007 — Go, the standard library only, and Git as the `git` program

Requirement: "LAYUP is written in Go with the standard library only, and calls
Git as the `git` program." Derives from `architecture.md` §1 and §2, ADR-0010,
ADR-0011 decision 1 and ADR-0013.

1. The module is `github.com/pharzam/layup` (the present `go.mod`), with one
   binary, `cmd/layup`, and packages under `internal/`.
2. `go.mod` has no `require` line. `go list -deps ./...` names only packages
   of the standard library and of this module (the criterion of `NFR-007`).
3. Only the package `internal/git` starts the `git` program, with
   `os/exec`; no other package imports a Git library or starts `git`. Its
   calls are [the calls of `internal/git`](#the-calls-of-internalgit).
4. Only the packages that the table below marks "starts a program" import
   `os/exec`.
5. No package of phase 1 depends on `net`, `net/http` or `crypto/tls`, by its
   own imports or through another package: the engine checks open no
   connection (`NFR-005`, [`gate.md`](gate.md#nfr-005--no-model-call-in-the-engine-checks)).
   The specification task of milestone `M2a` gives the forge adapter (phase 2)
   its rule, and that of `M3a` gives `internal/smartif` (phase 3) its rule
   ([the milestones](../plan/README.md#milestones)).

**Decided here** (D6 of #79, K8): the rules bind the non-test Go files only. A
test file (`_test.go`) can import any package of the standard library and of
this module, and can start any program. Reason: `go list -deps -test` names
`os/exec` in each test binary, through `testing`, and the tests start `git` and
`go`.

**Decided here** (D8 of #79, K39): rule 5 counts each package that
`go list -deps` names for a package, not only its own imports. Reason: a
connection that a dependency opens is a connection of the package that imports
it.

### The table of phase 1

"May import" lists the packages of this module that a package may import; any
package of the standard library is allowed, except as rules 3 to 5 say. An
import that the table does not allow is a defect, and the boundary rule of the
future Go gate of LAYUP reads this table.

**The form of a cell** (decided here, D7 of #79): "Package" is one code span.
"May import" is one or more code spans with a comma between two; for none, the
cell is the character — and no code span. "Starts a program" is exactly one
code span at the start of the cell, whose first word is the program, and text
with no code span after it; for none, the cell is the word no and no code
span. Reason: the [test of the package rules](#the-test-of-the-package-rules)
reads this table, and holds no copy of it, which could differ from it.

| Package | Job | May import | Starts a program |
| ------- | --- | ---------- | ---------------- |
| `cmd/layup` | `main`: passes the arguments to `internal/cli` and exits with its code | `internal/cli` | no |
| `internal/cli` | parses the arguments, runs one command, maps its result to an exit code ([`README.md`](README.md#commands)) | `internal/psb`, `internal/setup`, `internal/verify`, `internal/gate` | no |
| `internal/tsv` | reads and writes a record: checks the header row against a schema, the field count of each row, the types and the key; parses the `tsv-schema` blocks of `docs/spec/`, and compares a block with the Go schema of its record ([`README.md`](README.md#the-schema-block)) | — | no |
| `internal/git` | the one caller of the `git` program: [its calls](#the-calls-of-internalgit) | — | `git` |
| `internal/psb` | the rules G1 to G5 and the gap table ([`psb-check.md`](psb-check.md)) | `internal/tsv` | no |
| `internal/catalog` | the stack catalog, embedded with `embed` ([`setup.md`](setup.md#the-stack-catalog)) | `internal/tsv` | no |
| `internal/work` | the work area of a target: its paths, and the schemas and the readers of `answers.tsv` and `record.tsv` ([`setup.md`](setup.md#the-answers)), which `internal/setup` and `internal/verify` share (K9) | `internal/tsv` | no |
| `internal/gate` | runs the kinds of a gate manifest on a head ([`gate.md`](gate.md)) | `internal/tsv`, `internal/git` | `sh -c`: the gate commands |
| `internal/setup` | the step runner of `layup setup` ([`setup.md`](setup.md)) | `internal/tsv`, `internal/git`, `internal/catalog`, `internal/work` | no |
| `internal/verify` | the checks of `layup setup verify` ([`setup.md`](setup.md)) | `internal/tsv`, `internal/git`, `internal/catalog`, `internal/gate`, `internal/work` | `sh`: the baseline's own check scripts |
| `internal/standin` | for the tests only: a stand-in of the pinned baseline, and a work area set up from it, built at test time (K11; [`setup.md`](setup.md#the-checks-of-layup-setup-verify)); only test files import it | `internal/tsv`, `internal/git`, `internal/work` | no |

**Decided here:** the split of `internal/setup`, `internal/verify` and
`internal/gate` (ADR-0011 decision 1 names "setup, gates, … state files, Git
access" as concerns and gives no names). Reason: `layup setup verify` runs the
gates of a kind on a clean tree and on its known-bad fixture
([`setup.md`](setup.md)), so it imports `internal/gate`; `internal/setup` writes
a tree and must not depend on the checks that judge it, so `internal/cli`
runs the check of each step from `internal/verify` after `internal/setup` did
the step.

### The test of the package rules

**Decided here** (D7 and D9 of #79):

- The test is `TestPackageRules` in `cmd/layup/rules_integration_test.go`, an
  integration test: it starts `go` and reads files. The checker and its unit
  tests are untagged test files of package `main`: the binary holds no checker
  code, and the hook runs the unit part.
- It fails on a cell that it cannot read, on a missing heading or column, and
  on a table with no row. A package that the table names and that does not
  exist yet is not an error; a package that exists and has no row is.
- Rules 1, 2, 4, 5 and "May import" come from `go mod edit -json` and
  `go list -deps -json ./...` at the module root, and from the imports of each
  non-test Go file. `go list` gives no import of a file behind a build
  constraint, so the test lists each import that only such a file has, with
  `go list -e -deps -json`.
- Rule 3 and "Starts a program" come from a `go/ast` scan of each non-test Go
  file. A program starts with `exec.Command` or `exec.CommandContext` and a
  string literal that is the program of the row. Each other start that the
  scan finds is a defect: `os.StartProcess`, `syscall.Exec`,
  `syscall.ForkExec`, `syscall.StartProcess`, and an `exec.Cmd` that the code
  makes itself, whose program the scan cannot read. The scan does not read
  cgo code or a raw system call.
- The same checker must find the breaches of rule 5, and no other, in the
  fixture module `cmd/layup/testdata/netimport`: an import of `net/http`, and
  an import of `net/smtp` in a file behind a build constraint.

Reason: `cmd/layup` is the entry of the module, and its tests already start
programs; a new package for the test needs a row of its own, and the root
`tests/` holds end-to-end fixtures only.

### The calls of `internal/git`

The calls of phase 1 as task `T-2tc2` (#79) leaves them; a later task adds
each call that it needs here first. **Decided here** (D1 of #79, K7): the calls
that the steps, the checks and `layup gate` name.

| Call | The command, after the `-c` values below | Used by |
| ---- | ---------------------------------------- | ------- |
| `Version` | `git --version` | the minimum version (below) |
| `LsRemote` | `git ls-remote --exit-code -- URL REF` | S02 |
| `Clone` | `git clone --no-checkout -- URL DIR` | S02 |
| `CheckoutDetach` | `git checkout --detach COMMIT` | S02 |
| `Init` | `git init -b main -- DIR` | S03 |
| `Add` | `git add --all -- PATH…`; no path is the whole tree | S03 to S15; a fixture run of `gate:<kind>` |
| `Commit` | `git commit -m MESSAGE` | S03 to S15; a fixture run |
| `SwitchCreate` | `git switch -c BRANCH COMMIT` | S04: the branch `layup-setup` |
| `SwitchOrphan` | `git switch --orphan BRANCH` | S15: the branch `layup-records` |
| `RevParse` | `git rev-parse --verify --end-of-options REV` | S02, S03: the tree of a commit; `layup gate`: `--base`, `--head` |
| `RootCommits` | `git rev-list --max-parents=0 --end-of-options REV --` | check `pin` |
| `LsFiles` | `git ls-files -z` | S10; checks `markers` and `adapted` |
| `WorktreeAdd` | `git worktree add --detach -- PATH REV` | `layup gate`, step 2 of the run; a fixture run |
| `WorktreeRemove` | `git worktree remove --force -- PATH` | `layup gate`, step 4 of the run; a fixture run |
| `Show` | `git show --end-of-options REV:PATH --` | `layup gate`, steps 1 and 2 of the run |
| `DiffNames` | `git diff --name-only --no-renames -z --end-of-options BASE HEAD --` | `layup gate`: a `pending` kind |
| `Apply` | `git apply -- PATCH` | check `gate:<kind>`: the known-bad fixture |
| `LsTree` | `git ls-tree -r -z --full-tree --end-of-options REV -- PATH` | `layup gate`, step 2 of the run: the files of a `config` path at the base, with their modes (task `T-5sgt`) |

- `--end-of-options` or `--` comes before each revision, URL and path, so an
  input is never an option (`layup gate` takes revisions from its arguments).
  `checkout` and `switch` may read `--end-of-options` as a revision before
  `git` 2.44 (a reading of git's option parser; not measured, the LAYUP host
  has 2.54.0 only). So `CheckoutDetach` and `SwitchCreate` get none: `COMMIT`
  is a full object ID (40 or 64 hexadecimal characters), and the call refuses
  any other text before `git` starts.
- `DiffNames` names a renamed path at both ends, so a renamed product path
  counts as changed; `-z` gives each path unchanged.
- The orphan commit of S15 (task `T-d6q5`, #92) needs a scratch work tree, or
  a new call: `switch --orphan` removes the tracked files from the work tree.

**No default identity** (decided here, D2 of #79). `Commit` takes a name, an
e-mail address and a time: the author and the committer, and both dates
(`GIT_AUTHOR_DATE`, `GIT_COMMITTER_DATE`). Reason: one input then gives one
commit ID (`NFR-005`), and a CI runner has no `user.name`. The identity of the
setup's commits is an external input, a decision of the Operator: the plan of
the step runner (#85) asks the Operator for it, or writes a marker.

**A call reads no configuration of the host** (decided here, D3 of #79, with
condition 1 of its plan review and findings 1 and 7 of its verification).
Reason: no file and no variable of the host
may change a branch, a byte or a mode of a tree, a file list, a hook, an author
or a signature. `TestAHostileHostChangesNothing` of `internal/git` seeds each
input that the plan review measured with `git` 2.54.0, and first shows that it
changes a plain `git` run. `Commit` sets its identity after the fixed list, so
the host's `GIT_AUTHOR_NAME` cannot change a commit in that test: the unit test
of the environment, `TestTheEnvironmentIsAFixedList`, proves that it does not
pass.

- Each call starts with `-c core.hooksPath=/dev/null`,
  `-c core.attributesFile=/dev/null`, `-c core.excludesFile=/dev/null`,
  `-c core.autocrlf=false`, `-c core.precomposeUnicode=false`,
  `-c commit.gpgsign=false` and `-c http.emptyAuth=false`. Without the second
  and the third, the per-user attributes file changed the bytes of a staged
  file, and the per-user ignore file dropped a file.
- On macOS, `git init` writes `core.precomposeunicode = true` into the
  repository; `Add` then stored a decomposed (NFD) file name composed (NFC),
  so a tree with such a path differed on a macOS host only. With `false`,
  "file names are handled fully transparent by Git" (`git-config(1)`), and on
  APFS both forms kept their bytes (measured with `git` 2.54.0). A file system
  that changes a name itself ("the unicode decomposition of filenames done by
  Mac OS", the same manual; not measured) still changes the tree, and the S03
  check `root tree = pin.tree` stops the setup.
- The environment of a call is a fixed list: `PATH` and `TMPDIR` of the host
  when they are set; `LC_ALL=C`, `HOME=/dev/null`, `GIT_CONFIG_NOSYSTEM=1`,
  `GIT_CONFIG_GLOBAL=/dev/null`, `GIT_ATTR_NOSYSTEM=1`, `GIT_TERMINAL_PROMPT=0`,
  `GIT_ALLOW_PROTOCOL=file:git:http:https`; for a commit, the six variables of
  its identity. No other `GIT_*` variable, no `GIT_ASKPASS` and no
  `SSH_ASKPASS` of the host reaches `git`: `GIT_CONFIG_COUNT` set a value,
  `GIT_DIR` sent a commit to another repository, and `GIT_AUTHOR_NAME` changed
  the author. A call has no standard input.
- `GIT_ALLOW_PROTOCOL` allows `file` (the stand-in baselines of the tests),
  `https` (a public baseline), `http` (a baseline on a plain host, and the
  loopback server of the test of the no-prompt rule) and `git` (the
  unauthenticated `git` protocol). It refuses `ssh` and each remote helper.
- The list also leaves out the host's proxy and certificate variables
  (`http_proxy`, `https_proxy`, `no_proxy` and their upper-case forms,
  `SSL_CERT_FILE`, `SSL_CERT_DIR`, `GIT_SSL_CAINFO`), because a proxy address can
  carry a credential. So a host behind a proxy, or with a private certificate
  authority, cannot clone the baseline: part of the known limit
  [L-A7](../architecture.md#15-known-limits).
- The `git(1)` manual names `/dev/null` for `GIT_CONFIG_GLOBAL`, and the
  `git-config(1)` manual names `core.hooksPath=/dev/null`; the plan gave
  `core.hooksPath` an empty value, which no manual names. `GIT_ATTR_NOSYSTEM`
  is not in the manual of 2.54.0. Measured on the LAYUP host: its `git` has a
  system attributes file, which `GIT_CONFIG_NOSYSTEM=1` and
  `core.attributesFile=/dev/null` do not skip, and `GIT_ATTR_NOSYSTEM=1` does.
- No credential of the host reaches `git` (K31). With no configuration, no
  credential helper runs. libcurl, under `git`, reads `.netrc` in `HOME`, or in
  the home of the password database when `HOME` is not set; under
  `HOME=/dev/null` no file can be. `GIT_ALLOW_PROTOCOL` (`git(1)`) refuses
  `ssh` and each remote helper, so neither starts: `ssh` finds the user's home
  through the password database, not through `HOME`, and reads its keys there.
  `http.emptyAuth=false` stops GSS-Negotiate with no user name
  (`git-config(1)`; not measured, the LAYUP host has no Kerberos ticket).
  `TestNoCallUsesACredentialOfTheHost` seeds a `.netrc`, an `ssh` and a remote
  helper; under the host's `HOME` and with no `GIT_ALLOW_PROTOCOL`, the server
  got an `Authorization` header and both programs started
  ([`runs/T-2tc2/red.md`](../../runs/T-2tc2/red.md)).
- So `git` asks no question and uses no credential: a clone that needs one
  fails at once. The baseline's repository is public; a private baseline is
  known limit [L-A7](../architecture.md#15-known-limits).

**Two kinds of error** (decided here, D4 of #79, with condition 2 of its plan
review), as Go types for `errors.As`. `NotFoundError`: `git` is not on the
`PATH`. `FailedError`: each other failure, with the arguments, the exit code
and the standard error; the code is -1 when `git` did not start (also when
the call refused its input) or was stopped, and 0 when `git` exited 0 with an
output that the call cannot read.
The packages that may import `internal/git` (`internal/gate`, `internal/setup`,
`internal/verify`) tell the kinds apart and give `internal/cli` a result;
`internal/cli` maps results, never a `git` error, to exit codes. Reason: a
check that could not run is `not-active`, never `fail` (`NFR-004`).

**The minimum version of `git` is 2.32.0** (decided here, D5 of #79): the
first version with `GIT_CONFIG_GLOBAL` (the release notes of git 2.32.0, in
[`Documentation/RelNotes`](https://github.com/git/git/tree/master/Documentation/RelNotes)).
It also covers `switch` (2.23), `init -b` (2.28) and `--end-of-options`
(2.24; for `rev-parse`, 2.30). `internal/git` gives the
version (`Version`) and its test (`Supported`); the packages that import it
check it. The LAYUP host has 2.54.0; no test runs 2.32.0, so the minimum rests
on the release notes.

### The components of phase 1

[`architecture.md`](../architecture.md) §2 has a row per component and no row
for `layup setup`, which is neither an engine check (it writes files) nor
`layup run`. **Decided here:**

| Component | What it is | Where it runs | What it may write |
| --------- | ---------- | ------------- | ----------------- |
| The setup runner | `layup setup`: the steps of a target's setup that need no judgement ([`setup.md`](setup.md)) | the LAYUP host | its work area: the baseline copy and its setup branch, the files it prints commands for; and its table on standard output |

### The packages of later phases

Name, job and phase only. Each later milestone gives a package its row in the
table above, with its import rules, in the same change that adds its
requirement's section.

| Package | Job | Phase |
| ------- | --- | ----- |
| `internal/records` | the records branch: one writer, the schemas of every record kind ([`records.md`](records.md)) | 2 |
| `internal/forge` | the forge interface: the six capabilities of `architecture.md` §1 | 2 |
| `internal/forge/github` | the GitHub adapter of the forge interface | 2 |
| `internal/run` | `layup run`: Start, the phase loop, the lease, fencing | 2 |
| `internal/session` | a role session: its directory, its start, its result | 2 |
| `internal/handoff` | the transition table and the check of a handoff | 2 |
| `internal/spec` | `layup spec check` | 2 |
| `internal/rules` | the rule-path register, the check `layup/rules`, rule batches | 2 |
| `internal/audit` | `layup audit` | 2 |
| `internal/ledger` | the cost ledger (`telemetry.tsv`) and the budget | 2 |
| `internal/smartif` | the smart-if client: the only package that connects to a model service | 3 |
| `internal/escalate` | the escalation screen | 3 |
| `internal/stall` | progress, the stall triggers and the procedure | 3 |
| `internal/route` | the harness register, the probe, admission and routing | 2 |
| `internal/report` | `layup report`: the measures | 4 |
| `internal/learn` | `layup learn`: the reward and the routing update | 4 |

**Decided here:** the names and the phases. The phase of a package is the
earliest `PRD-0001` phase whose requirement needs it: `layup run` and the
ledger are first needed for the handoffs of `REQ-005` (phase 2), which start
role sessions.
