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
   From milestone `M2a`, only the packages that the column "Connects" of
   [the table of M2a](#the-table-of-m2a) names depend on them (by D8 below,
   through the adapter too), and one package alone imports them itself; `M3a`
   gives `internal/smartif` (phase 3) its rule
   ([the milestones](../plan/README.md#milestones)).
   The one package that imports them: `internal/forge/github`.

**Decided here** (task `T-esfe`, #127, row 22a of the plan): the test of the
package rules reads the one package of rule 5 from the line that starts with
"The one package that imports them:", its one code span, so the test holds no
list of its own (D7 of #79); no such line, two such lines, or a line with no
code span or with two, is an error. The three packages of rule 5 (`net`,
`net/http`, `crypto/tls`) stay a constant of the checker, as before.

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
with no code span after it (a register row's cell is its own form, [below](#the-table-of-m2b)); for none, the cell is the word no and no code
span. Reason: the [test of the package rules](#the-test-of-the-package-rules)
reads this table, and holds no copy of it, which could differ from it.

| Package | Job | May import | Starts a program |
| ------- | --- | ---------- | ---------------- |
| `cmd/layup` | `main`: passes the arguments to `internal/cli` and exits with its code | `internal/cli` | no |
| `internal/cli` | parses the arguments, runs one command, maps its result to an exit code ([`README.md`](README.md#commands)); for `layup run`, reads the two host registers and the key file, builds the adapter and hands it to `internal/run` (task `T-mqty`), and from `M2b` reads `models.tsv` and `routing.tsv` and checks each harness credential (task `T-ysph`) | `internal/psb`, `internal/setup`, `internal/verify`, `internal/gate`, `internal/run`, `internal/forge`, `internal/forge/github`, `internal/route` | no |
| `internal/tsv` | reads and writes a record: checks the header row against a schema, the field count of each row, the types and the key; parses the `tsv-schema` blocks of `docs/spec/`, and compares a block with the Go schema of its record ([`README.md`](README.md#the-schema-block)) | — | no |
| `internal/git` | the one caller of the `git` program: [its calls](#the-calls-of-internalgit) | — | `git` |
| `internal/psb` | the rules G1 to G5 and the gap table ([`psb-check.md`](psb-check.md)) | `internal/tsv` | no |
| `internal/catalog` | the stack catalog, embedded with `embed` ([`setup.md`](setup.md#the-stack-catalog)) | `internal/tsv` | no |
| `internal/work` | the work area of a target: its paths, and the schemas and the readers of `answers.tsv` and `record.tsv` ([`setup.md`](setup.md#the-answers)), which `internal/setup` and `internal/verify` share (K9); the record row `answers.sha256`, which `internal/setup` writes and `internal/standin` gives its stand-in; the texts of the four questions of S01, which the stop table, the answers record of S04 and the stand-in use (task `T-7s0y`); the schema of the table of `layup setup verify`, which `internal/verify` writes and S15 reads (task `T-d6q5`) | `internal/tsv` | no |
| `internal/gate` | runs the kinds of a gate manifest on a head ([`gate.md`](gate.md)) | `internal/tsv`, `internal/git` | `sh -c`: the gate commands |
| `internal/setup` | the step runner of `layup setup` ([`setup.md`](setup.md)) | `internal/tsv`, `internal/git`, `internal/catalog`, `internal/work` | no |
| `internal/verify` | the checks of `layup setup verify` ([`setup.md`](setup.md)) | `internal/tsv`, `internal/git`, `internal/catalog`, `internal/gate`, `internal/work` | `sh`: the baseline's own check scripts |
| `internal/records` | the schemas and the row rules of the record kinds that phase 1 defines and later phases write: `telemetry.tsv` and `prices.tsv` of `REQ-011` ([`records.md`](records.md#req-011--the-telemetry-record), task `T-tmhw`), and `stalls.tsv` of `REQ-009` ([`records.md`](records.md#req-009--the-stall-record), task `T-dgy7`); and the records of Start of `M2a`, `start.tsv`, `approvers.tsv`, `lease.tsv` and `copies.tsv` ([`records.md`](records.md#nfr-001--the-records-of-start), task `T-8kqn`), which `internal/run` writes; and the records of a session of `M2b`, `sessions.tsv`, `harnesses.tsv`, `routing.tsv`, the events and the result of a task ([`records.md`](records.md#nfr-001--the-records-of-a-session), task `T-ywk7`; built by task `T-3py1`) | `internal/tsv` | no |
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

### The table of M2a

Milestone `M2a` (task `T-zck8`, #123; [`run.md`](run.md), [`forge.md`](forge.md)).
The columns of the table of phase 1, and one more: **"Connects"** is the
character — for a package that may not depend on a package of rule 5, or one
code span per package of rule 5 that it may depend on, by its own imports or
through another package (D8). `cmd/layup` and `internal/cli` depend on them
through the adapter, so they have a row here too, with only their "Connects"
cell; their other cells stay in the table of phase 1. **Decided
here** (O-170 of #127; task `T-esfe`, row 22a of the plan): `TestPackageRules`
reads this table too, and rule 5 from this column, before any package of this
table exists, so the test holds no list of its own (D7 of #79). A cell of this
table that starts with "(the row of phase 1" stands for that package's cell of
the table of phase 1; a row whose cells Job, May import and Starts a program
each start so adds only its cell Connects to that row of phase 1.

| Package | Job | May import | Starts a program | Connects |
| ------- | --- | ---------- | ---------------- | -------- |
| `internal/run` | `layup run`: the steps of Start, the restart, the lease and fencing, the copy of a comment and the rule of a decision ([`run.md`](run.md)); the Go schema of the table `run-steps`; the checks of `git` and of the values of the flags, which `internal/cli` calls before the first step (task `T-mqty`); in `M2b`, the step `probe`, the call that starts an attempt (the event `attempt`), the records of a session, the checks before a push, the push and the bind ([`session.md`](session.md)) | `internal/tsv`, `internal/git`, `internal/records`, `internal/forge`, `internal/route`, `internal/session`, `internal/ledger`, `internal/rules` | no | — |
| `internal/forge` | the forge interface: the six capabilities and their types ([`forge.md`](forge.md)), `Missing` and `CheckPermissions` (task `T-6bq5`); the reader of the forge register (`host:registers/forge.tsv`), with the Go schema of its block; the check of the key file (task `T-1g1q`) | `internal/tsv` | no | — |
| `internal/forge/github` | the GitHub adapter: the JWT, the installation token, the calls of [`forge.md`](forge.md#the-calls-of-m2a) | `internal/forge` | no | `net`, `net/http`, `crypto/tls` |
| `internal/route` | the reader of the harness register (`host:registers/harnesses.tsv`), with the Go schema of its block; in `M2b`, the columns that `M2b` adds, the readers of `models.tsv` and `routing.tsv`, the checks across the three and of a credential file (task `T-ysph`), the version that needs a probe, admission and the order of the routing register ([`session.md`](session.md#req-013--the-probe-admission-and-routing)) | `internal/tsv`, `internal/records` | no | — |
| `cmd/layup` | (the row of phase 1) | (the row of phase 1) | (the row of phase 1) | `net`, `net/http`, `crypto/tls` |
| `internal/cli` | (the row of phase 1, with the change below) | (the row of phase 1, with the change below) | (the row of phase 1) | `net`, `net/http`, `crypto/tls` |

In `M2a` two rows of the table of phase 1 change, in the build task that needs
each: `internal/cli` may also import `internal/run`, `internal/forge`, `internal/forge/github` and
`internal/route` (it reads the forge register through `internal/forge`, builds
the adapter from it, and hands it to `internal/run`); `internal/records` also holds the schemas of the records of
Start ([`records.md`](records.md#nfr-001--the-records-of-start)), and still
imports `internal/tsv` only. `internal/run` commits and pushes the records with
`internal/git`; `internal/records` gives the rows. The later table below keeps
`internal/route` for its jobs of `M2b`.

### The table of M2b

Milestone `M2b` (task `T-ywk7`, #147; [`session.md`](session.md)). The columns
of [the table of M2a](#the-table-of-m2a). **Decided here** (condition 2 of the
plan review of #147): a package that a table already has keeps its one row,
changed in place (`internal/route` and `internal/run` in the table of `M2a`,
`internal/records` in the table of phase 1, and the calls of `internal/git`
[below](#the-calls-of-internalgit)); this table holds only the packages that no
table has. `TestPackageRules` reads this table (task `T-y10b`, row 29 of [the
plan](../plan/README.md#the-tasks-of-m2b)), and the line of the engine checks
below it ([`session.md`](session.md#nfr-005--no-harness-in-the-engine-checks)).

| Package | Job | May import | Starts a program | Connects |
| ------- | --- | ---------- | ---------------- | -------- |
| `internal/session` | a role session: its directory, its environment and credential, the rule-file check, the start, the limits and the stop, its end, the result file and the usage report ([`session.md`](session.md#req-013--a-role-session)) | `internal/tsv`, `internal/git`, `internal/records` | `harness` (the command of a harness register row) | — |
| `internal/ledger` | the writer of `telemetry.tsv`: one row per session, from the usage report and `host:prices.tsv` ([`session.md`](session.md#req-011--the-writer-of-the-telemetry-record)) | `internal/tsv`, `internal/records` | no | — |
| `internal/rules` | in `M2b`, the reader of the rule-path register and the check of a diff before a push ([`session.md`](session.md#a-workflow-or-rule-path-change)), with a second Go value of the block `rule-paths`, whose owner is `internal/setup`, compared with the block by its own test (task `T-m1dx`); the check `layup/rules` and rule batches come in `M2f` | `internal/tsv`, `internal/records` | no | — |

The packages of the engine checks: `internal/psb`, `internal/verify`, `internal/gate`.

**The rule of the engine checks** (NFR-005, task `T-y10b`): no package of that
line depends on `internal/session` or on a package of rule 5, by any path; the
rule reads no cell, so a cell Connects cannot lift it. A later engine check adds
its package to the line. No such line, two such lines, a line with no code span,
or a span that names no package of the module is an error of the test, so the
rule never passes with a package left out.

**A register row.** A cell "Starts a program" that is exactly one code span, one
space and the words "(the command of a harness register row)" makes the row a
register row: the span names the register, not a program. Its package may
import `os/exec` (rule 4), and may start, by a call of `exec.Command` or
`exec.CommandContext`, a program that is not a string literal: the command of a
harness register row. A start by a string literal, a use of `exec.Command` that
is not a call, and each other start are findings there, as elsewhere. A cell that
holds those words and any other text cannot be read. **Known limit:** in a
register row the scan cannot tell `git` from a harness, so rule 3 rests on review
in that package.

### The test of the package rules

**Decided here** (D7 and D9 of #79):

- The test is `TestPackageRules` in `cmd/layup/rules_integration_test.go`, an
  integration test: it starts `go` and reads files. The checker and its unit
  tests are untagged test files of package `main`: the binary holds no checker
  code, and the hook runs the unit part.
- It reads the table of phase 1, [the table of M2a](#the-table-of-m2a) and the
  line of rule 5 (task `T-esfe`), and [the table of M2b](#the-table-of-m2b) and
  the line of the engine checks (task `T-y10b`). A row of the table of `M2b` is a
  package of its own: a row whose cell starts with "(the row of phase 1", or
  whose package another table has, is an error. It fails on a cell that it cannot read, on a
  missing heading or column, on a table with no row, on a row of the table of
  `M2a` that mixes cells of phase 1 and cells of its own, on such a row of
  phase 1 that the table of phase 1 lacks, and on a package in two rows. A
  package that a table names and that does not exist yet is not an error; a
  package that exists and has no row is.
- Rule 5 holds two checks: an own import of a package of rule 5 in a package
  other than the one of its line, and a dependency on one that the cell
  Connects of the package does not name.
- Rules 1, 2, 4, 5 and "May import" come from `go mod edit -json` and
  `go list -deps -json ./...` at the module root, and from the imports of each
  non-test Go file. `go list` gives no import of a file behind a build
  constraint, so the test lists each import that only such a file has, with
  `go list -e -deps -json`.
- Rule 3 and "Starts a program" come from a `go/ast` scan of each non-test Go
  file. A program starts with `exec.Command` or `exec.CommandContext` and a
  string literal that is the program of the row, or, in a register row, a
  program that is not a string literal ([the table of M2b](#the-table-of-m2b)).
  Each other start that the scan finds is a defect: `os.StartProcess`, `syscall.Exec`,
  `syscall.ForkExec`, `syscall.StartProcess`, and an `exec.Cmd` that the code
  makes itself, whose program the scan cannot read. The scan does not read
  cgo code or a raw system call.
- The same checker must find the breaches of rule 5 and of the rule of the
  engine checks, and no other, in the fixture module
  `cmd/layup/testdata/netimport`: an import of `net/http` in `cmd/layup`,
  outside the one package of rule 5, and an import of `net/smtp` in a file
  behind a build constraint, whose dependency on `net` no cell Connects of
  `internal/psb` names, and which breaks the rule of the engine checks in
  `internal/psb` too. In the fixture module `cmd/layup/testdata/enginedep`, it
  must find that `internal/gate`, a package of the engine checks, depends on
  `internal/session` (and imports it, which its row does not allow), and no other
  breach: the start of `internal/session`, by a variable, is its register row's.

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
| `LsRemote` | `git ls-remote --exit-code -- URL REF` | S02; `layup run`, step 3 |
| `Clone` | `git clone --no-checkout -- URL DIR`; with a token, as `Fetch` and `Push` (task `T-ax3r`) | S02; `layup run`, step 3 and the restart (a private target needs the token) |
| `CheckoutDetach` | `git checkout --detach COMMIT` | S02 |
| `Init` | `git init -b main -- DIR` | S03; `layup run`, step 5 |
| `Add` | `git add --all -- PATH…`; no path is the whole tree | S03 to S15; a fixture run of `gate:<kind>` |
| `Commit` | `git commit -m MESSAGE` | S03 to S15; a fixture run |
| `SwitchCreate` | `git switch -c BRANCH COMMIT` | S04: the branch `layup-setup`; `M2b`: the branch of a session's clone ([`session.md`](session.md#the-session-directory)) |
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
| `DiffNames` | `git diff --name-only --no-renames -z --end-of-options BASE HEAD --` | `layup gate`: a `pending` kind; `M2b`: the check before a push ([`session.md`](session.md#a-workflow-or-rule-path-change)) |
| `Apply` | `git apply -- PATCH` | check `gate:<kind>`: the known-bad fixture |
| `LsTree` | `git ls-tree -r -z --full-tree --end-of-options REV -- PATH` | `layup gate`, step 2 of the run: the files of a `config` path at the base, with their modes (task `T-5sgt`); S04: the records of `docs/adr/` and `docs/facts/` at the root commit (task `T-7s0y`); S05: the history at the head; S11: the records of `docs/facts/` at the head (task `T-b3r1`); S12: the tree of the head; S15: the `.sh` files of `layup-setup`, and the tree of the records commit of a run that stopped (task `T-d6q5`) |

**The calls of `M2a`** (task `T-zck8`; [`run.md`](run.md)). Each takes the
token of [`forge.md`](forge.md#the-app-identity) in its environment when the
remote needs it; no other value enters the fixed list.

| Call | The command, after the `-c` values below | Used by |
| ---- | ---------------------------------------- | ------- |
| `Fetch` | `git fetch --no-tags --update-head-ok -- URL REF:REF` | `layup run`: the read-back of the root commit (step 5 of [`run.md`](run.md#the-steps-of-layup-run---new)) |
| `Push` | `git push --porcelain -- URL COMMIT:refs/heads/BRANCH`; never `--force` | `layup run`: each records commit (fencing: a push that is not a fast-forward is refused); `M2b`: the push of a session's SHA ([`session.md`](session.md#the-push-and-the-bind)) |

- `Fetch` and `Push` refuse their input before `git` starts unless `REF`
  starts with `refs/` and holds no `:`, `COMMIT` is a full object ID and
  `BRANCH` is not empty and holds no `:`, so no `+` forces an update, and,
  with a token, unless `web` is not empty and has no final slash, as the key
  `http.<web>/.extraHeader` of [`forge.md`](forge.md#the-app-identity) needs
  (task `T-xhgz`). `--update-head-ok` lets `Fetch` set the branch of a repository that
  `Init` made: without it, `git` 2.54.0 refuses to fetch into the branch that is
  checked out (exit 128); the result is the state that `Clone --no-checkout`
  gives. A push that `git` refuses (not a fast-forward, or refused by the remote)
  is a `FailedError` with `Code` 1; its line `!` is on the standard output of
  `--porcelain`, which the error does not keep. A remote that cannot be reached
  gives another code (128, measured with `git` 2.54.0). `Code` 1 is also a
  commit that is not in the local repository, which a records commit that the
  run has just made never is.

**The calls of `M2b`** (task `T-ywk7`; [`session.md`](session.md)). A session's
clone and `layup run`'s clone are both local directories, so no call of `M2b`
takes a token.

| Call | The command, after the `-c` values below | Used by |
| ---- | ---------------------------------------- | ------- |
| `CloneLocal` | `git clone --no-local --no-checkout --single-branch --no-tags --branch BRANCH -- SRC DIR`, then `git -C DIR remote remove origin` | the clone of a session, `repo/` |
| `FetchSession` | `git init --bare -- TMP`; the object directory of the session's `.git` as the one line of `TMP/objects/info/alternates`; `git -C TMP cat-file -t SHA`, which must print `commit`, else the result is refused (`branch`); `git -C TMP update-ref refs/heads/session SHA`; `git -c core.hooksPath=EMPTY fetch --no-tags --no-write-fetch-head -- TMP refs/heads/session:DST`, where `EMPTY` is an empty directory; `TMP` removed | the fetch of a session's head into `layup run`'s clone, with no command in the session's clone |
| `IsAncestor` | `git merge-base --is-ancestor BASE HEAD`; exit 1 is "no", not an error | the check that a session's head descends from its base |
| `DiffFile` | `git diff -U0 --no-color --no-renames --end-of-options BASE HEAD -- PATH` | the exception of `docs/guardrails.md` before a push |
| `DiffBinary` | `git diff --binary --no-renames --end-of-options BASE HEAD --` | the payload of a refused diff |

- `CloneLocal` copies the objects and shares no ref, configuration or hook with
  `SRC` (`--no-local`, §4); `IsAncestor` takes full object IDs only, as
  `CheckoutDetach` does, and refuses any other text before `git` starts.
- `FetchSession` takes the SHA that `internal/session` read from the files of
  the session's `.git`, a full object ID; `git upload-pack` runs in `TMP`, whose
  configuration is `layup`'s, so no configuration of the session's clone is read
  ([`session.md`](session.md#the-fetch-by-sha)).

Start makes its clone with `Init` and `Fetch`, and the restart with `Clone`. The
first records commit is an orphan commit in a scratch work tree, as S15 makes it
(`WorktreeAdd`, `SwitchOrphan`, `Add`, `Commit`, `RevParse`, `WorktreeRemove`);
each later one uses the same calls on the last records commit, with no
`SwitchOrphan`.

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
  its identity; with a token, `Fetch` and `Push` add `GIT_CONFIG_COUNT`,
  `GIT_CONFIG_KEY_0` and `GIT_CONFIG_VALUE_0` of
  [`forge.md`](forge.md#the-app-identity), with LAYUP's values. No other
  `GIT_*` variable, no `GIT_ASKPASS` and no `SSH_ASKPASS` of the host reaches
  `git`: `GIT_CONFIG_COUNT` set a value,
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
| `internal/records` | the schemas and rows of the records ([`records.md`](records.md)); the package has its row in the table of phase 1 since task `T-tmhw` (K23); `internal/run` commits the records ([the table of M2a](#the-table-of-m2a)) | 2 |
| `internal/forge` | the forge interface: the six capabilities of `architecture.md` §1; its row is in [the table of M2a](#the-table-of-m2a) | 2 |
| `internal/forge/github` | the GitHub adapter of the forge interface; its row is in [the table of M2a](#the-table-of-m2a) | 2 |
| `internal/run` | `layup run`: Start, the phase loop, the lease, fencing; its row is in [the table of M2a](#the-table-of-m2a) | 2 |
| `internal/session` | a role session: its directory, its start, its result; its row is in [the table of M2b](#the-table-of-m2b) | 2 |
| `internal/handoff` | the transition table and the check of a handoff | 2 |
| `internal/spec` | `layup spec check` | 2 |
| `internal/rules` | the rule-path register, the check `layup/rules`, rule batches; its part before a push is in [the table of M2b](#the-table-of-m2b) | 2 |
| `internal/audit` | `layup audit` | 2 |
| `internal/ledger` | the cost ledger (`telemetry.tsv`) and the budget; the writer is in [the table of M2b](#the-table-of-m2b), the budget comes later | 2 |
| `internal/smartif` | the smart-if client: the only package that connects to a model service | 3 |
| `internal/escalate` | the escalation screen | 3 |
| `internal/stall` | progress, the stall triggers and the procedure | 3 |
| `internal/route` | the harness register, the probe, admission and routing; its row is in [the table of M2a](#the-table-of-m2a), with its jobs of `M2b` | 2 |
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
