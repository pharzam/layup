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
`P` (any depth) and ends with `.E`; `./*.E` matches such a file anywhere. A
pattern with no `*` matches each file whose path starts with it. Example for
Go: `./*.go`. A **product path** of a kind is a path that its scope matches.

### The command

```text
layup gate REPO --base REV --head REV
```

- `REPO`: a local Git clone of the target. `layup gate` changes nothing in it
  except a scratch work tree, which it removes before it exits.
- `--base REV`: the base of the change (for a pull request, its base branch's
  commit). `--head REV`: the head to judge. Each is any revision that `git
  rev-parse` resolves in `REPO`.
- Exit codes: 0 when each row is `pass` or `clear`; 1 when a row is `fail` or
  `not-active`; 2 on a usage error, a revision that does not resolve, or a
  manifest at the base that is missing or does not match its schema.

**The run.**

1. Read `docs/gates.tsv` at the base (`git show <base>:docs/gates.tsv`).
2. Make a scratch work tree of the head (`git worktree add --detach`). Write
   into it, from the base, the manifest and each path of each row's `config`;
   a path that the base does not have is removed from the scratch tree. So the
   head's own gate files never judge the head (FT4).
3. For each row, in the order of the manifest, give one result by the table
   below: the first line of the table that matches the row decides. A command runs with the scratch tree as its working directory.
4. Remove the scratch work tree. Print the table.

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
3. A missing toolchain is `not-active` (`tool not found`), never `pass` or
   `clear` (known limit L-B1).
4. A missing or malformed manifest is an input error (exit 2), never an empty
   table with exit 0.
5. The fixture test that the criterion asks for ("a check that does not run
   cannot produce `pass`") comes with the code (#29): a manifest whose `tool`
   does not exist gives `not-active` and exit 1.

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
