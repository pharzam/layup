# `layup gate` and the gate manifest

Conventions: [`README.md`](README.md). The stack gates live in the target, as
its own CI jobs; `layup gate` runs the same commands from outside
([`architecture.md`](../architecture.md) §6, ADR-0016).

## REQ-004 — The stack gates of a target, run from outside

Requirement: "`layup gate` runs the stack-dependent gates of a target
(repository layout, interface boundaries, contract checks, test quality) from
outside the target and reports `pass`, `fail` or `not-active` per gate; the
gates add rules and weaken none." Derives from `architecture.md` §6 (The stack
gates), ADR-0016 and ADR-0011 decision 4.

### The gate manifest

The setup writes it into the target ([`setup.md`](setup.md#the-steps), step
S12), from the stack catalog. One row per gate kind; the target's CI has one
job per row, with the kind as the job's name.

```tsv-schema gate-manifest target:docs/gates.tsv
kind     id(<word>)           key  the gate kind: `static`, `layout`, `boundary`, `contract`, `test`, or a further kind of the stack's catalog entry; lowercase letters and `-`
state    enum(active|pending)  -   `pending` until the activation batch of the first bet makes it `active` (architecture §6)
tool     text                 -    the program that the command needs, as `exec.LookPath` finds it, for example `go`
command  text                 -    the command, run with `sh -c` at the root of the tree, for example `test -z "$(gofmt -l .)"`
scope    list(text)           -    the paths in scope, each a pattern (below)
config   list(path)           -    the gate files of this kind: its configuration and its tests; `—` when none
```

Two packages hold this schema as a Go value, and each one's test compares it
with this block (task `T-3jpx`, #81): `internal/catalog`, which writes the
manifest of a stack's entry, and `internal/gate`, which reads it (task `T-5sgt`,
row 5 of the plan). The block is the one home of the form.

**Decided here:** the columns `tool` and `config`. §6 names "gate kind, command,
scope, state". `layup gate` must know which program a kind needs, to tell "did
not run" from "failed" (`NFR-004`), and which files are the kind's gate files,
to take them from the base (below).

**A scope pattern** (decided here: §6 says "the paths in scope" and gives no
form). A pattern `P/*.E` matches each file whose path is under the directory
`P` (any depth) and whose name ends with `.E`; `./*.E` matches such a file
anywhere. A pattern with no `*` matches the file at that path, and each file
under it as a directory: `internal` matches `internal/x.go`, not
`internal2/x.go` (decided here, task `T-5sgt`, #82: a string prefix would let
`internal` match `internal2/`). `P`, and a pattern with no `*`, have the form of
the type `path`, and `E` is not empty and holds no `/` or `*`; any other form,
for example `*.go`, `P/*`, `P/*/x.E`, `./internal`, `../x/*.go` or `/x/*.go`, is
an input error of the manifest (exit 2), so a pattern never matches nothing in
silence. Example for Go: `./*.go`. A **product path** of a kind is a path that
its scope matches.

**The column `tool`** names the one program that `layup gate` can look up.
Decided here (task `T-5sgt`): a command that needs a second program that is
not found exits 127 in `sh`, which is `fail`, not `not-active`. So a kind's
command needs no program that its `tool` does not imply (a rule for the
catalog entries, task `T-c06a`, row 14 of the plan).

### The command

```text
layup gate REPO --base REV --head REV
```

- `REPO`: a local Git clone of the target. `layup gate` changes nothing in it
  except a scratch work tree, which it removes before it exits.
- `--base REV`: the base of the change (for a pull request, its base branch's
  commit). `--head REV`: the head to judge. Each is any revision that `git
  rev-parse` resolves in `REPO` to a commit (decided here, task `T-5sgt`: each
  is peeled with `^{commit}`, so a tag gives its commit, and a tree or a blob is
  an input error).
- Exit codes: 0 when each row is `pass` or `clear`; 1 when a row is `fail` or
  `not-active`; 2 on a usage error, a revision that does not resolve to a
  commit, or a manifest at the base that is missing, does not match its schema
  or has no row. **Decided here** (task `T-5sgt`): also 2, with no table, when
  `git` is not found or is older than 2.32, or `sh` is not found; and 2 after a
  complete table when the scratch work tree cannot be removed, with a diagnostic
  that names its path, because the run then changed `REPO`. `internal/gate`
  gives each of these as its own input error; `internal/cli` never maps a `git`
  error ([`packages.md`](packages.md#the-calls-of-internalgit)). Reason: code 2
  is what the user fixes in the input, on the host or in `REPO`; a check that
  could run and did not is `not-active` (`NFR-004`).

**The run.**

1. Read `docs/gates.tsv` at the base (`git show <base>:docs/gates.tsv`).
2. Make a scratch work tree of the head (`git worktree add --detach`). Write
   into it, from the base, the manifest and each path of each row's `config`;
   a path that the base does not have is removed from the scratch tree. So the
   head's own gate files never judge the head (FT4).
3. For each row, in the order of the manifest, give one result by the table
   below: the first line of the table that matches the row decides. A command runs with the scratch tree as its working directory.
4. Remove the scratch work tree. Print the table.

**Decided here** (task `T-5sgt`, #82):

- **The scratch tree** is in a new temporary directory outside `REPO`. Every
  call of `internal/git` has the hooks off (`core.hooksPath=/dev/null`), so no
  hook of `REPO` runs. A `config` path is put as the base has it: a file, or
  each file under a directory, with its mode (`LsTree` of
  [`packages.md`](packages.md#the-calls-of-internalgit)); a `config` path that
  is a symbolic link or a submodule at the base is an input error. Each write
  goes through an `os.Root` of the tree, so a symbolic link of the head cannot
  send a write out of the tree; such a write fails the overlay. The tree is
  removed on every exit path of a run. **Known limit:** a run that a signal
  kills leaves the tree; `git worktree prune` in `REPO` removes its record. The
  run does not prune at its start, because that could remove another stale
  record of `REPO`.
- **The product paths of an `active` kind** are the files of the scratch tree
  after the overlay, without `.git`; a `pending` kind reads
  `git diff --name-only --no-renames -z`, so a renamed path counts at both ends,
  and `<first path>` is the first in `git`'s order that the scope matches.
- **The command** runs with `sh -c` and with the environment of `layup`, which
  a target's tool needs (`HOME`, its caches). A setting of the host, for
  example `GOFLAGS`, can change a verdict: a case of known limit L-A1
  ([`architecture.md`](../architecture.md#15-known-limits)); the target's own
  CI job runs the same command on a clean runner. `<name>` of a signal is the
  name that Go gives it (`syscall.Signal.String()`, for example `terminated`).
  **Known limit:** no timeout in phase 1; a command that hangs blocks the run,
  its progress lines show that it is alive
  ([`README.md`](README.md#commands), Progress), and the target's CI job has its
  own timeout.
- **A failure of the scratch tree or of the diff** gives the rows that need
  the failed part a fixed reason, with no part of the error, because an error
  of `git` can name a scratch path, which a table never holds (the repeat rule
  of `REQ-007`); the error itself goes to standard error. A `pending` row needs
  no scratch tree, so it keeps its result from the diff: a result is given
  wherever its input exists.

| The failure | The rows | Result | Reason |
| ----------- | -------- | ------ | ------ |
| the scratch tree cannot be made (`git worktree add`, the temporary directory) | each `active` row | `not-active` | `scratch tree: add failed` |
| the overlay or the reading of the tree fails | each `active` row | `not-active` | `scratch tree: overlay failed` |
| `git diff` fails | each `pending` row | `not-active` | `diff failed` |

| The row | Result | Reason |
| ------- | ------ | ------ |
| `pending`, and `git diff --name-only <base> <head>` names no product path of the kind | `clear` | `pending: no product path` |
| `pending`, and the diff names a product path of the kind | `fail` | `pending: product path changed: <first path>` |
| `active`, and its `tool` is not found | `not-active` | `tool not found: <tool>` |
| `active`, and the scratch tree has no product path of the kind | `clear` | `no product path` |
| `active`, the command exits 0 | `pass` | `—` |
| `active`, the command exits with another code, or is killed by a signal | `fail` | `exit <code>`, or `signal <name>` |

A `pending` kind's command never runs: `clear` is the pass of the rule "no
product path may change while this kind is pending", not of the kind's gate
(the reading of `REQ-004` in `PRD-0001` §7.1).

### The table

```tsv-schema gate-result stdout
base    sha1                                 -    the base commit
head    sha1                                 -    the head commit
kind    id(<word>)                           key  the kind, as in the manifest
state   enum(active|pending)                 -    as in the manifest at the base
result  enum(pass|fail|not-active|clear)     -    by the table of the run
reason  text                                 -    by the table of the run; `—` for `pass`
```

The command's own output goes to standard error, one block per kind, headed by
the kind, so that a human sees why a kind failed; it is not part of the table.
**Decided here** (task `T-5sgt`): the heading of a block is the step line of the
kind ([`README.md`](README.md#commands), Progress); the command's standard output
and standard error go through one shared writer, so their lines keep the order
in which the command wrote them, and they are held and printed after the step
ends, so no progress line comes inside them.

### Not in phase 1

- **The commit status `layup/gates`** that `layup run` posts from this table
  (§6): phase 2, with `layup run` and the forge adapter.
- **A rule batch's own gate files** on its head, and the run with each recorded
  known-bad patch and the approval hash (§6 Activation): phase 2, with the
  rule batches of `REQ-003`.
- **The known-bad fixtures at setup** are run by `layup setup verify`, not by
  this command ([`setup.md`](setup.md#the-checks-of-layup-setup-verify)).

## REQ-007 — Each change passes the gates before review or merge

Requirement: "Each change passes the gates for repository layout, interface
boundaries and testing pyramids before it reaches human review or merges; a
gate is deterministic where a rule can be checked mechanically." Derives from
`architecture.md` §6 and §8 (the task loop, step 5), ADR-0016 and ADR-0019.

**In phase 1:**

1. **The target's own jobs.** The setup writes one CI job per gate kind
   ([`setup.md`](setup.md#the-steps), step S12), each a required check of the
   default branch's ruleset, which the Operator applies (S13). A pull request
   with a failed job cannot merge, with or without LAYUP.
2. **The same verdict from outside.** `layup gate` gives the verdict of each
   kind by the rules of `REQ-004`.
3. **The repeat rule.** Two runs of `layup gate` with the same `REPO` content,
   base and head print the same bytes. The table holds no time and no scratch
   path, its rows follow the manifest, and a gate command is expected to be
   deterministic; a kind whose command gives two verdicts on one input is a
   defect of its catalog entry.

**Not in phase 1:** the status `layup/gates` as a required check, and the rule
that no pull request is marked ready or merged while a selected kind did not
run (§8 steps 5 and 8): phase 2, with `layup run`.

## NFR-004 — A check that is not active is not passed

Requirement: "A check that is not active does not count as passed." Derives
from `architecture.md` §6 (the table of `layup/gates`), ADR-0016 and ADR-0011
decision 4. It holds for `layup gate` and `layup setup verify` alike.

1. A result is `pass` only when the check ran and passed. A check that did not
   run is `not-active`, and `not-active` gives exit code 1.
2. `clear` counts as a pass only in the two cases of the table of `REQ-004`
   (`layup gate`) and in the cases that [`setup.md`](setup.md#the-checks-of-layup-setup-verify)
   names (`layup setup verify`); each `clear` row carries its reason.
3. A missing toolchain of an `active` kind is `not-active` (`tool not
   found`), never `pass` or `clear` (known limit L-B1). **Decided here** (K19 of
   the [defect register](../plan/README.md#the-defect-register), task
   `T-5sgt`): a `pending` kind's command never runs, so its tool is not looked
   up; a `pending` kind with a missing tool gives `clear` or `fail` by the two
   `pending` lines of the table of the run, never `pass`.
4. A missing or malformed manifest, and a manifest with no row, is an input
   error (exit 2), never an empty table with exit 0.
5. The fixture test that the criterion asks for ("a check that does not run
   cannot produce `pass`") is `TestGateNeverPassesACheckThatDidNotRun` of
   `cmd/layup` (task `T-5sgt`): a manifest whose `tool` does not exist gives
   `not-active` and exit 1.

## NFR-005 — No model call in the engine checks

Requirement: "A deterministic check is preferred to an LLM judgement wherever
a rule can be checked mechanically; the engine checks make no model call, and
the `layup` process calls a model only through the smart-if provider." Derives
from `architecture.md` §4, ADR-0015.

1. No command of phase 1 opens a connection to a model service or starts a
   model process. Its only network use is `git` to the baseline's repository
   (S02) and the gate commands of a target. The import rule of [`packages.md`](packages.md#nfr-007--go-the-standard-library-only-and-git-as-the-git-program)
   (no `net`, `net/http` or `crypto/tls` in phase 1) is the mechanical check.
   A gate command of a target may use the network (for example `go` that
   fetches modules); that is the target's tool, not a model call of `layup`.
2. Each verdict of `layup gate` and `layup setup verify` is reproducible by the
   repeat rule of `REQ-007`.

**Not in phase 1:** the smart-if client and `decisions.tsv`, one row per call
([`records.md`](records.md#the-layout-of-the-records-branch)): phase 3.
