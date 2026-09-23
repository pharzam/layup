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

The Operator decided, on 2026-09-23, after the panel's questions. O-12 and O-13
are copied word for word; O-9 to O-11 summarise the answers, whose full text is on
the issue of this decision:

| ID | Decision |
|----|----------|
| O-9 | Rule protection for Invariant 3 (`F-0001#3`): **decide later.** Invariant 3 stays an open gap, and this ADR names the missing control. (Fact recorded with it: the agents and the Operator push with the same GitHub account today.) |
| O-10 | Nothing from LAYUP goes into a target repository beyond the Armature kit itself. |
| O-11 | Strict reading of O-10: LAYUP writes nothing into a target except the Armature kit copy and the facts; stack gates run outside the target, in LAYUP's own runner. |
| O-12 | What LAYUP is, in the Operator's words: "The layup project final result is the product that uses armature as a kit to get a new problem statement brief and architectural vision brief from the idea owner and orchestrate the paths to deliver the solution for that problem statement. It's not a scaffold or template for obeying or adapting to the Product that should be created" |
| O-13 | The Operator's clarification of O-12, word for word: "The layup project final result is the product that uses armature as a kit to get a new problem statement brief and architectural vision brief from the idea owner and orchestrate the paths to deliver the solution for that problem statement. It's not a scaffold or template for obeying lay up help the idea honors idea owners to idea owners or anyone has product statement brief or vision brief and want to solve it and deliver the product. And lay up here help this role for using armature as a baseline. With adapting it in new repository. And machinery it as a principle and discipline in the armature kit. Or based on armature kit rules that adapted in new product repository. For going forward and deliver the final product." |

This record reads O-10 to O-13 together: LAYUP is an **orchestrator product**. It
helps anyone who has a problem statement and a vision brief to deliver the
product. A **target** is the repository of that new product. LAYUP creates it,
adapts Armature into it as the baseline, and then drives delivery by the adapted
kit's own rules. So the target holds the adapted kit, the facts, the product that
the orchestrated role agents deliver, and the records that the kit's rules
require — task files, decisions, and resource records of tokens and time
([ADR-0007](0007-record-task-resource-use.md)). "The kit copy" in O-11 is read as
the kit with those records. LAYUP's own machinery — its engine code, scripts, and
binary — is never written into a target.

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
   An event table (telemetry, gate results) only adds rows; a register (for
   example `open-gaps.tsv`) is edited in place, and Git history is its log.
   Per-task records go under an evidence store, `runs/<task-id>/` in LAYUP's own
   repository (Operator decision O-2) and the value the setup of a target sets
   for it; setup state stays under `docs/setup/`. A deterministic check
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
   against a checkout of it; they add rules and weaken none (`F-0001#7`). Every
   gate verdict is deterministic (`F-0001#6`), and each gate reports `pass`,
   `fail` or `not-active`; `not-active` never counts as passed (`F-0001#5`).
5. **LAYUP's own checks.** In this repository, `docs/setup/setup-check.sh` stays
   the POSIX `sh` gate for LAYUP's own setup, and the Go code is tested with
   `go test` under the Go gates. The engine does not replace the kit's checks.
6. **Steps.** `layup setup` follows `setup/steps.tsv`. It stops at each row whose
   `human_decision` is `yes`, asks its questions in one batch, and writes the
   answers to Git before the next step (Decision Point 2, `F-0001#11`). A `no`
   row needs no human decision, but it is not always a program step: the rows
   that write prose (S07 onboarding, S08 glossary, S09 guardrails, S14 identity)
   are done by a role agent that a harness agent runs, and the engine then runs
   the row's deterministic check. The engine itself performs the other `no`
   rows.
7. **Copy without Node.** The engine copies the kit with `git clone` of the
   pinned commit and removes the `.git` directory; it does not call `npx degit`
   (ADR-0010 rejected a Node runtime). The manual run used `npx degit`, and
   step S02 of `setup/steps.tsv` still says so. Likewise, step S12 adds a CI job
   that runs LAYUP's `setup-check.sh`, which cannot go into a target (decision 3):
   in a target, CI runs the kit's own jobs only, and `layup setup verify` checks
   the setup from outside (decision 4). The task that builds `layup setup`
   changes S02 and S12.
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
  gates need a LAYUP runner that reports on the target's pull requests; the
  stack-gate child of the Step 2 parent issue decides that runner.
- The question Q-1 (may LAYUP write state records into a target?) is answered by
  O-13: the target follows the adapted kit's rules, and those rules keep the
  records in the project repository. If the Operator corrects the reading,
  decision 3 changes with it, in a new ADR.
- Invariant 3 has no check (O-9). **The missing control** is a guard before the
  merge: a code-owners rule on the rule paths (`docs/*.md` rules, `.github/**`,
  the check scripts, and the gate code under `cmd/layup` and `internal/`) that
  needs an approval from an account that the agents do not use. It waits for O-9. Until then, a deterministic post-merge list of each
  change to a rule path is the complement.
- Two languages hold checks: `sh` for LAYUP's own setup, Go for the engine. A
  check that exists in both can drift apart, the risk that `guardrails.md` §2
  names for a hand-mirrored check set; the engine's setup verification must run
  the same fixtures as `setup/tests/run.sh`.
- The next ADR is `0012`.
