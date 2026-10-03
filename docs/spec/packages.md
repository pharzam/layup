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
| `internal/work` | the work area of a target: its paths, and the schemas and the readers of `answers.tsv` and `record.tsv` ([`setup.md`](setup.md#the-answers)), which `internal/setup` and `internal/verify` share (K9); the record row `answers.sha256`, which `internal/setup` writes and `internal/standin` gives its stand-in; the texts of the four questions of S01, which the stop table, the answers record of S04 and the stand-in use (task `T-7s0y`); the schema of the table of `layup setup verify`, which `internal/verify` writes and S15 reads (task `T-d6q5`) | `internal/tsv` | no |
| `internal/gate` | runs the kinds of a gate manifest on a head ([`gate.md`](gate.md)) | `internal/tsv`, `internal/git` | `sh -c`: the gate commands |
| `internal/setup` | the step runner of `layup setup` ([`setup.md`](setup.md)) | `internal/tsv`, `internal/git`, `internal/catalog`, `internal/work` | no |
| `internal/verify` | the checks of `layup setup verify` ([`setup.md`](setup.md)) | `internal/tsv`, `internal/git`, `internal/catalog`, `internal/gate`, `internal/work` | `sh`: the baseline's own check scripts |
| `internal/records` | the schemas and the row rules of the record kinds that phase 1 defines and later phases write: `telemetry.tsv` and `prices.tsv` of `REQ-011` ([`records.md`](records.md#req-011--the-telemetry-record), task `T-tmhw`), and `stalls.tsv` of `REQ-009` ([`records.md`](records.md#req-009--the-stall-record), task `T-dgy7`); the one writer of the records branch is phase 2's | `internal/tsv` | no |
| `internal/standin` | for the tests only: a stand-in of the pinned baseline, and a work area set up from it, built at test time (K11; [`setup.md`](setup.md#the-checks-of-layup-setup-verify)); only test files import it | `internal/tsv`, `internal/git`, `internal/work` | no |

**Decided here:** the split of `internal/setup`, `internal/verify` and
`internal/gate` (ADR-0011 decision 1 names "setup, gates, … state files, Git
access" as concerns and gives no names). Reason: `layup setup verify` runs the
gates of a kind on a clean tree and on its known-bad fixture
([`setup.md`](setup.md)), so it imports `internal/gate`; `internal/setup` writes
a tree and must not depend on the checks that judge it, so `internal/cli`
runs the check of each step from `internal/verify` after `internal/setup` did
the step. For the same reason `internal/cli` hands `internal/setup` the lists
that `internal/verify` gives from a tree: the flagged files of the prose step
(task `T-8ya0`), the markers of S10 and the files whose links break at S05
(K10, task `T-8vpw`; [`setup.md`](setup.md#the-checks-of-layup-setup-verify)),
and the evidence call of a step: a function that runs the one-check call of
`internal/verify` on the work area (D12 of #86, task `T-7s0y`). Since task
`T-b3r1` (D9 of #90) these are one value, `setup.Calls`, which `internal/cli`
fills from `internal/verify`: the one-check call, the markers of a tree (each
with the byte column of its open quote, so S11 fills a marker where the scanner
found it, D7), the files whose links break (S05), the files that check
`adapted` flags (S14), and the link rule of check `kit-history` (S05).
`internal/cli` also refuses a brief that holds a marker, by the same scanner. In the same way
`internal/cli` reads the problem statement once per run, runs `internal/psb`,
and hands `internal/setup` the gap table, as `layup psb check` writes it, and
the SHA-256 of the file; `internal/setup` reads the table by its own Go value
of the block `psb-gaps`, which its block test compares with the block (D5 of
#86), so the IDs `Q-NNN` keep one home and `internal/setup` imports no
`internal/psb`.

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
| `Branch` | `git symbolic-ref --quiet HEAD` | the step runner: a commit of S04 to S14 only on `layup-setup` (task `T-79y7`) |
| `RevParse` | `git rev-parse --verify --end-of-options REV` | S02, S03: the tree of a commit; S04: the root commit and the branch `layup-setup`; the step runner: the head of `layup-setup` (task `T-7s0y`); `layup gate`: `--base`, `--head`; S15: the heads of `layup-setup` and `layup-records`; a fixture run: its commit (task `T-d6q5`) |
| `ResetSoft` | `git reset --soft COMMIT` | the step runner: the undo of a step whose evidence fails (task `T-7s0y`, D12 of #86) |
| `ResetHard` | `git reset --hard --quiet HEAD` | the step runner: before each step from S04 to S14, the target back to its head (task `T-b3r1`, D11 of #90) |
| `Staged` | `git diff --cached --quiet --exit-code`; exit 1 is a staged change | the step runner: a commit only of a staged change (task `T-b3r1`, D12 of #90) |
| `RootCommits` | `git rev-list --max-parents=0 --end-of-options REV --` | check `pin`; S03: the root commit of a run that stopped (task `T-7s0y`); S15: the records commit of a run that stopped (task `T-d6q5`) |
| `Message` | `git log -1 --format=%B --end-of-options REV --` | S03: the message of the root commit of a run that stopped (task `T-7s0y`); S15: the same for the records commit (task `T-d6q5`) |
| `IsShallow` | `git rev-parse --is-shallow-repository` | check `pin` (task `T-6x75`) |
| `LsFiles` | `git ls-files -z` | S10; checks `markers` and `adapted` |
| `WorktreeAdd` | `git worktree add --detach -- PATH REV` | `layup gate`, step 2 of the run; `layup setup verify`: the scratch tree; a fixture run; S15: the scratch tree of the records commit (task `T-d6q5`) |
| `WorktreeRemove` | `git worktree remove --force -- PATH` | `layup gate`, step 4 of the run; `layup setup verify`; a fixture run; S15 |
| `Show` | `git show --end-of-options REV:PATH --` | `layup gate`, steps 1 and 2 of the run; `layup setup verify`: the manifest at the head of `layup-setup`; S04: the two index files and `docs/setup/facts.sha256` of the root commit (task `T-7s0y`); S05, S06 and S11: the task indexes, the facts index, `facts.sha256` and `open-gaps.tsv` of the head (task `T-b3r1`); S12: the paths of the entry and `open-gaps.tsv` of the head; S13 and S15: the manifest of the head; S15: the files of the records commit of a run that stopped (task `T-d6q5`) |
| `DiffNames` | `git diff --name-only --no-renames -z --end-of-options BASE HEAD --` | `layup gate`: a `pending` kind |
| `Apply` | `git apply -- PATCH` | check `gate:<kind>`: the known-bad fixture |
| `LsTree` | `git ls-tree -r -z --full-tree --end-of-options REV -- PATH` | `layup gate`, step 2 of the run: the files of a `config` path at the base, with their modes (task `T-5sgt`); S04: the records of `docs/adr/` and `docs/facts/` at the root commit (task `T-7s0y`); S05: the history at the head; S11: the records of `docs/facts/` at the head (task `T-b3r1`); S12: the tree of the head; S15: the `.sh` files of `layup-setup`, and the tree of the records commit of a run that stopped (task `T-d6q5`) |

- `--end-of-options` or `--` comes before each revision, URL and path, so an
  input is never an option (`layup gate` takes revisions from its arguments).
  `checkout` and `switch` may read `--end-of-options` as a revision before
  `git` 2.44 (a reading of git's option parser; not measured, the LAYUP host
  has 2.54.0 only). So `CheckoutDetach`, `SwitchCreate` and `ResetSoft` get
  none: `COMMIT` is a full object ID (40 or 64 hexadecimal characters), and the
  call refuses any other text before `git` starts.
- `DiffNames` names a renamed path at both ends, so a renamed product path
  counts as changed; `-z` gives each path unchanged.
- S15 makes the orphan commit in a scratch work tree (task `T-d6q5`, #92), as
  `switch --orphan` removes the tracked files from the work tree.

**No default identity** (decided here, D2 of #79). `Commit` takes a name, an
e-mail address and a time: the author and the committer, and both dates
(`GIT_AUTHOR_DATE`, `GIT_COMMITTER_DATE`). Reason: one input then gives one
commit ID (`NFR-005`), and a CI runner has no `user.name`. The identity of the
setup's commits is a decision of the Operator: O-136 (#85) makes it the LAYUP
App's bot, at the date `pin.time` ([`setup.md`](setup.md#the-steps)).

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
  `-c commit.gpgsign=false`, `-c http.emptyAuth=false` and
  `-c maintenance.auto=false`. Without the second and the third, the per-user
  attributes file changed the bytes of a staged file, and the per-user ignore
  file dropped a file. Without the last, a commit runs `git maintenance run
  --auto`, whose tasks go on in the background after the call ends (since
  `git` 2.47); in CI of the pull request #115 (`git` 2.55.0) such a task
  repacked into `.git/objects` of a work area while a test removed it (task
  `T-d6q5`). Which version first repacks after a commit is not shown: `git`
  2.54 made the "geometric" strategy the default of manual maintenance, and a
  hand run on 2.54.0 saw no repack after a commit (review round 2 of #92). So
  no call leaves a process of `git` in a repository of LAYUP.
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
| `internal/records` | the one writer of the records branch ([`records.md`](records.md)); the package has its row in the table of phase 1 since task `T-tmhw` (K23), with the schemas of the record kinds that phase 1 defines | 2 |
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

**Decided here** (task `T-tmhw`, #94, K23): `internal/records` has its row in
the table of phase 1, as the section of `REQ-011` is in phase 1 and this table
gives the package "the schemas of every record kind"; one home serves the
telemetry record and its price table (row 17) and the stall record (row 18),
and `internal/ledger` and `internal/stall` import the schemas from it in their
phases.

**Decided here:** the names and the phases. The phase of a package is the
earliest `PRD-0001` phase whose requirement needs it: `layup run` and the
ledger are first needed for the handoffs of `REQ-005` (phase 2), which start
role sessions.
