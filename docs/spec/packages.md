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
   `os/exec`; no other package imports a Git library or starts `git`.
4. Only the packages that the table below marks "starts a program" import
   `os/exec`.
5. No package of phase 1 imports `net`, `net/http` or `crypto/tls`: the engine
   checks open no connection (`NFR-005`, [`gate.md`](gate.md#nfr-005--no-model-call-in-the-engine-checks)).
   The one later exception is `internal/smartif` (phase 3), the smart-if
   client, and the forge adapter (phase 2).

### The table of phase 1

"May import" lists the packages of this module that a package may import; any
package of the standard library is allowed, except as rules 3 to 5 say. An
import that the table does not allow is a defect, and the boundary rule of the
future Go gate of LAYUP reads this table.

| Package | Job | May import | Starts a program |
| ------- | --- | ---------- | ---------------- |
| `cmd/layup` | `main`: passes the arguments to `internal/cli` and exits with its code | `internal/cli` | no |
| `internal/cli` | parses the arguments, runs one command, maps its result to an exit code ([`README.md`](README.md#commands)) | `internal/psb`, `internal/setup`, `internal/verify`, `internal/gate` | no |
| `internal/tsv` | reads and writes a record: checks the header row against a schema, the field count of each row, the types and the key; parses the `tsv-schema` blocks of `docs/spec/`, and compares a block with the Go schema of its record ([`README.md`](README.md#the-schema-block)) | — | no |
| `internal/git` | the one caller of the `git` program: clone, `ls-remote`, `init`, commit, `rev-parse`, `worktree`, `show`, `diff --name-only`, `apply` | — | `git` |
| `internal/psb` | the rules G1 to G5 and the gap table ([`psb-check.md`](psb-check.md)) | `internal/tsv` | no |
| `internal/catalog` | the stack catalog, embedded with `embed` ([`setup.md`](setup.md#the-stack-catalog)) | `internal/tsv` | no |
| `internal/gate` | runs the kinds of a gate manifest on a head ([`gate.md`](gate.md)) | `internal/tsv`, `internal/git` | the gate commands, with `sh -c` |
| `internal/setup` | the step runner of `layup setup` ([`setup.md`](setup.md)) | `internal/tsv`, `internal/git`, `internal/catalog` | no |
| `internal/verify` | the checks of `layup setup verify` ([`setup.md`](setup.md)) | `internal/tsv`, `internal/git`, `internal/catalog`, `internal/gate` | the baseline's own check scripts, with `sh` |

`internal/psb` today imports no package of this module and writes its table
itself; it moves to `internal/tsv` with task `T-5zmw`, row 6 of the
[plan](../plan/README.md#the-tasks-of-phase-1).

**Decided here:** the split of `internal/setup`, `internal/verify` and
`internal/gate` (ADR-0011 decision 1 names "setup, gates, … state files, Git
access" as concerns and gives no names). Reason: `layup setup verify` runs the
gates of a kind on a clean tree and on its known-bad fixture
([`setup.md`](setup.md)), so it imports `internal/gate`; `internal/setup` writes
a tree and must not depend on the checks that judge it, so `internal/cli`
runs the check of each step from `internal/verify` after `internal/setup` did
the step.

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
