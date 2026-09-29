# 0013. Orchestrate a target from outside with `layup run`

Date: 2026-09-29

## Status

Proposed

## Context

The Operator decided that LAYUP is "a deterministic orchestrator driving the
entire lifecycle" of a new product (O-67), and that it "exists for the target
strictly as an external binary"; the target's stack and tools "have zero
dependency on LAYUP" (O-76). [ADR-0011](0011-structure-the-core-engine-as-a-go-cli-over-repository-files.md)
made LAYUP one Go command-line program over repository files, with no daemon,
no database and no network service. The review of #69 found that nothing starts
or moves the role agents (B1), and the deep check found that a design which let
several commands write records had no single writer (Sol-3, Fable-M12).

The hands-on evaluation of eighteen public candidates found no base and no
program that should run next to LAYUP: each keeps state that a Git clone does not
carry, or writes into the target or the host
([`runs/T-hbw8/evaluation/summary.md`](../../runs/T-hbw8/evaluation/summary.md)).
The options compared are in
[`runs/T-hbw8/selection-v2.md`](../../runs/T-hbw8/selection-v2.md), row 0013.

## Decision

We will orchestrate each target from outside it with one foreground process,
`layup run TARGET`:

1. **One process per target.** `layup run` runs on the LAYUP host in the
   foreground and drives the phase loop of one target. It holds a lease row on
   the target's records branch; a second run that finds a live lease, or whose
   push is refused, stops.
2. **The engine checks are pure.** `layup psb check`, `layup setup verify`,
   `layup gate`, `layup spec check`, `layup audit` and `layup report` read files
   and print a typed table. They write nothing else. `layup run` calls them in
   its own process and commits their results.
3. **Only three things go into a target:** the setup output that O-76 and
   `F-0003#42` name; the work of the role sessions, through pull requests; and
   LAYUP's records, on the records branch
   ([ADR-0014](0014-keep-the-records-in-the-target-with-one-writer.md)). No file
   goes in that the target needs LAYUP to build, test or pass its gates.
4. **The forge is GitHub for the pilot,** behind one package with a named set of
   calls (issues and comments, pull requests, commit statuses, the effective
   rules of a branch, the repository activity).

This amends ADR-0011 decision 1 (`layup run` is a long-running foreground
process; it is still not a daemon, and LAYUP still has no database and no network
service of its own) and decision 3 (what LAYUP writes into a target).

We reject: a loop of short runs on the forge's runner (the harness credentials
and sessions would live on the forge, and the whole loop would depend on it); a
daemon or a database (state that a clone does not carry, FT6); any evaluated
candidate as a runtime next to LAYUP (each fails Invariant 1, 2 or 5 in its own
code).

## Consequences

- One host must be up while a target is in delivery; when it is down, nothing
  moves, and the stall procedure finds it (a later slice of this task).
- The target stays independent: it holds the setup output, its product and
  plain-text records.
- A second forge needs a second implementation of the forge package; the pilot
  does not have one (known limit L-A2 of [`architecture.md`](../architecture.md)).
