# 0011. Structure the core engine as a Go CLI over repository files

Date: 2026-09-23

## Status

Accepted

## Context

Step 2 builds the LAYUP core engine with no UI. PSB Invariant 1 says Git is the
system of record (`F-0001#1`); Invariant 2 says a repository works without the
automation that made it (`F-0001#2`); Invariant 6 prefers a deterministic check to
an LLM judgement (`F-0001#6`); Invariant 9 says a harness agent is replaceable
(`F-0001#9`). The stack is Go ([ADR-0010](0010-use-go-as-the-technology-stack.md)).
The manual setup procedure that the engine automates is
[`setup/README.md`](../setup/README.md) with [`setup/steps.tsv`](../setup/steps.tsv).
The architectural vision brief (`F-0002`) was read as input only; none of its
solutions is a requirement.

A panel of three members that differ in domain (CLI and Go architecture;
Git-native state; governance and the invariants) generated options, one pass each,
with no vote (ADR-0006, Operator decision O-4). The compared set and the recorded
disagreement are on the issue of this decision.

The Operator decided, on 2026-09-23, after the panel's questions:

| ID | Decision |
|----|----------|
| O-9 | Rule protection for Invariant 3 is decided later. The agents and the Operator push with the same GitHub account today, so no forge setting can tell them apart. |
| O-10 | Nothing from LAYUP goes into a target repository beyond the Armature kit itself. |
| O-11 | Strict reading of O-10: LAYUP writes nothing into a target except the Armature kit copy and the facts; stack gates run outside the target, in LAYUP's own runner. |
| O-12 | What LAYUP is, in the Operator's words: "The layup project final result is the product that uses armature as a kit to get a new problem statement brief and architectural vision brief from the idea owner and orchestrate the paths to deliver the solution for that problem statement. It's not a scaffold or template for obeying or adapting to the Product that should be created" |

This record reads O-10 to O-12 together as follows, and the Operator can correct
the reading: LAYUP is an **orchestrator product**. A **target** is the repository
of the new product that LAYUP delivers from the idea owner's problem statement and
vision brief. The target holds the adapted Armature kit (its discipline system),
the facts, the product code that the orchestrated role agents deliver, and the
project's own state records, because Invariant 1 puts project state in the
project repository. LAYUP's own machinery — its engine code, scripts, and binary —
is never written into a target.

## Decision

We will build the core engine as **one Go command-line program, `layup`**, with
these properties:

1. **Form.** One Go module, one binary `layup` under `cmd/layup`, and packages
   under `internal/` (one per concern: problem-statement check, setup, gates,
   telemetry, state files, Git access). The standard library only; Git is called
   as the `git` program. No daemon, no database, no network service.
2. **State.** All state is plain files in a repository, read and written as
   tab-separated tables with a fixed header row (one file per record kind, one
   row per event, rows only added), beside Markdown records in the kit's shape.
   Per-task records go under the evidence store `runs/<task-id>/` (Operator
   decision O-2); setup state stays under `docs/setup/`. A deterministic check
   validates each table's header and columns. A human reads every state file
   with no tool.
3. **Where it writes.** Into a target, the engine writes the adapted Armature kit
   copy (the steps of `setup/steps.tsv`), the facts (the problem statement and
   the vision brief), and the project's state records under `runs/` and
   `docs/`. The product code is written by the role agents that LAYUP
   orchestrates. The engine writes no LAYUP code, script, or binary into a target
   (O-10, O-11).
4. **Gates.** A target's own gates are the Armature kit's gates, so the target
   passes them without LAYUP (Invariant 2). LAYUP's stack-dependent gates and its
   setup verification are `layup` commands that run **outside** the target,
   against a checkout of it; they add rules and weaken none (Invariant 7). Every
   gate verdict is deterministic, and each gate reports `pass`, `fail` or
   `not-active`; `not-active` never counts as passed (Invariant 5).
5. **LAYUP's own checks.** In this repository, `docs/setup/setup-check.sh` stays
   the POSIX `sh` gate for LAYUP's own setup, and the Go code is tested with
   `go test` under the Go gates. The engine does not replace the kit's checks.
6. **Steps.** `layup setup` follows `setup/steps.tsv`: it performs each row whose
   `human_decision` is `no`, and it stops at each `yes` row to ask its questions
   in one batch and write the answers to Git before the next step.
7. **Copy without Node.** The engine copies the kit with `git clone` of the
   pinned commit and removes the `.git` directory; it does not call `npx degit`
   (ADR-0010 rejected a Node runtime).
8. **No LLM in the engine.** The engine makes no model call; judgement stays with
   the harness agents that call it. Its interface is plain text and files, so any
   harness agent can run it (Invariant 9).

We reject: a pinned `layup` binary in the target's CI, and LAYUP check source
copied into the target (both forbidden by O-10 and O-11); append-only JSON events
with derived state (a human needs a tool to read the state); state in `git notes`
or an orphan branch (outside the tree a reviewer reads); and rules in a separate
repository (a target then fails its gates without it).

## Consequences

- A target passes its own gates without LAYUP by construction. LAYUP's stack
  gates need a LAYUP runner that reports on the target's pull requests; the task
  that builds the stack gates (`T-vk3k`) decides that runner.
- The question Q-1 (may LAYUP write state records into a target?) is answered by
  the reading of O-12 above: state records are project data, not LAYUP
  machinery, so they go into the target. If the Operator corrects the reading,
  decision 3 changes with it, in a new ADR.
- Invariant 3 has no check (O-9). A deterministic post-merge list of each change
  to a rule path is the complement until O-9 is decided.
- Two languages hold checks: `sh` for LAYUP's own setup, Go for the engine. A
  check that exists in both is a hand mirror (`guardrails.md` §2); the engine's
  setup verification must run the same fixtures as `setup/tests/run.sh`.
- The next ADR is `0012`.
