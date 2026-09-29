# 0015. Keep model calls out of the engine checks

Date: 2026-09-29

## Status

Proposed

## Context

[ADR-0011](0011-structure-the-core-engine-as-a-go-cli-over-repository-files.md)
decision 8 says: "The engine makes no model call". The Operator decided that
decision 8 gets a successor: "the engine checks make no model call, and only the
decision component calls Jev" (O-72), and that the smart-if provider is
pluggable (O-78). `layup run` also starts role sessions, which are harness
processes that call models; a criterion that counts the child processes of
`layup` then fails (Fable-N2).

The deep check found that a session on the host gets rules that no repository
file holds: a global instruction file, a command-rewrite hook and a plugin's
start text (Fable-M23). In the evaluation, three tools wrote into the Operator's
home directory on their first calls, and one tried to change branch protection
on GitHub (`runs/T-hbw8/evaluation/summary.md`, "Incidents"). The options
compared are in [`runs/T-hbw8/selection-v2.md`](../../runs/T-hbw8/selection-v2.md),
row 0015.

## Decision

We will keep model calls out of the engine checks, and start each role session
with only what LAYUP gives it:

1. **The engine checks** start no model process and open no connection to a
   model service.
2. **The `layup` process** opens a connection to a model service only in the
   smart-if provider client, at a named decision point, and writes one row to
   `decisions.tsv` for each call.
3. **A role session** is a harness process that `layup run` starts. Its model
   calls belong to the harness; the cost ledger records the session.
4. **A session starts with** its own clone of the target at the base commit (not
   a linked work tree of LAYUP's clone), an empty home directory, a prompt file
   that code builds, and an environment from a named list: the harness's own
   credential, as a variable or as a file that the register row names, copied
   into the home directory; no forge credential and no SSH agent. The command
   line comes from the harness's row in the harness register. Before the start,
   code refuses a rule file of that harness in any directory from the session
   directory up to the root, and records each system-wide policy file that the
   row lists; the session's configuration directory is in its empty home.
5. **A session ends** when its process exits. `layup run` reads its typed result
   file, fetches its branch by commit SHA into LAYUP's own clone with hooks
   turned off, checks that it descends from the base commit and changes nothing
   under `.github/workflows/` (the App has no workflows permission, O-92), and
   pushes that SHA with the App's token. A session is not resumed across
   attempts.

This replaces ADR-0011 decision 8. The `NFR-005` criterion of `PRD-0001` changes
with it: "the engine checks start no model process and open no connection to a
model service; the `layup` process opens one only from the smart-if client, and
each such call writes one row to `decisions.tsv`".

We reject: a harness SDK linked into LAYUP (LAYUP uses the standard library only,
`F-0004#1`); a session that uses the Operator's home directory (the incidents of
the evaluation, Fable-M23).

## Consequences

- A second harness gets the same rules as the first, from the same files
  (Invariant 9).
- A harness whose login is a file in its configuration directory gets a copy of
  that file in the session's home directory; the harness register names it.
- A process on the host can still read what the Operator's user can read; the
  separation holds against a session that uses its prompt and its environment
  (known limit L-A1 of [`architecture.md`](../architecture.md)).
