# 0013. Orchestrate the lifecycle with an external `layup run`

Date: 2026-09-28

## Status

Proposed

## Context

The review of the previous architecture (PR #69, closed without a merge) found
that no component started or moved the role agents (finding B1), so a human
stayed the message bus of PSB Problem 1 (`F-0003#1`), and that no process watched
the clock, so a stall was found by chance (finding B2, `F-0003#17`). The
Operator then decided that LAYUP is a deterministic orchestrator of the whole
lifecycle (O-67), that LAYUP is only an external binary and a target never
depends on it (O-76), and that every limit and policy of a run is a parameter
that the Operator sets per target (O-78 to O-83). The decisions are copied word
for word in [`runs/T-hbw8/`](../../runs/T-hbw8/operator-decisions.md).

A panel of three members generated the options (ADR-0006; O-70); the compared
set is [`runs/T-hbw8/selection.md`](../../runs/T-hbw8/selection.md), section 3.
The search for a public solution found agent orchestrators (for example
Bernstein, Agent Orchestrator, Claude Squad) and durable workflow engines
(Temporal, `durabletask-go`); none keeps its state only in Git, and none runs the
phases of the PSB (section 1 of the selection).

## Decision

We will build the orchestrator as the command `layup run TARGET` of the one
`layup` binary ([ADR-0011](0011-structure-the-core-engine-as-a-go-cli-over-repository-files.md),
decision 1):

1. **External.** `layup run` runs outside the target, in the foreground, on the
   Operator's host. It reads and writes the target through Git and the forge
   API. Nothing of LAYUP goes into the target's product, stack or tools (O-76).
2. **The phase loop.** A state machine with the phases Intake, Scaffold, Design,
   Architect, Plan, Implement, Accept and Retrospective; after the
   Retrospective, the next milestone starts at Design, until every `Must`
   requirement is accepted (O-67). Each phase has an exit check; a phase moves
   only on a valid handoff record, a passed exit check and, at a planned
   approval point, the human decision copied into Git
   ([ADR-0016](0016-keep-the-run-records-in-git-and-tell-agent-from-human.md)).
3. **Role sessions.** For each step, `layup run` starts one harness session as a
   child process through a harness adapter (a command line from the harness
   register, [ADR-0015](0015-route-role-sessions-over-registered-harnesses.md)),
   in its own worktree, with a prompt file that names the records to read and
   the handoff to write. It stops the process group on a timeout, a budget stop
   or a cancel. A session that exits with no valid handoff fails.
4. **Durable events.** Each transition and each attempt is a row on the
   target's records branch, pushed before the next role starts. An attempt ID
   stops a duplicate transition. On a restart, `layup run` fetches the records,
   finds the last pushed transition and marks an interrupted attempt before it
   starts new work. A restart that a human causes is recorded as unplanned
   input (`F-0001#28`).
5. **The clock.** A progress event is a new durable record: a new artifact hash,
   a passed gate run, a handoff row. An equal failure signature and no new
   pushed head are not progress. When no progress event comes for `T` minutes,
   or after `N` retries with no progress, `layup run` opens a stall record.
   The next action is one of: a retry (only while fewer than `N` retries were
   made), a fresh examiner on a model and a harness that did not do the stalled
   work, a panel, or the Operator with a stall package (`F-0003#49`); which of
   the allowed actions is taken is a decision point of
   [ADR-0014](0014-decide-at-named-points-with-a-pluggable-smart-if-provider.md).
   A harness that hides its inner tool calls is judged by the same clock.
6. **The dead-man check.** `layup run` pushes a heartbeat row every `H`
   minutes. A scheduled job in LAYUP's own repository reads the heartbeat of
   each active run; a heartbeat older than `2H` opens a stall record for the
   orchestrator itself and notifies the Operator with a stall package, not by
   chance.
7. **The parameter register.** Every value that O-78 to O-83 make a parameter
   is a row of `parameters.tsv` on the records branch: `name`, `value`,
   `default`, `source`, `set_at`. `layup run` reads it at the start of each
   step; a changed row takes effect at the next step and is itself a record.
   LAYUP ships the defaults, each with its source (for example `stall.T` = 10
   minutes and `stall.N` = 1 from O-82). A parameter can tune when and how a
   rule of the PSB applies; it cannot switch off a PSB invariant (O-84): a
   business-forking decision goes to the idea owner (`F-0001#13`), a check that
   did not run is not a pass (`F-0001#5`), and a failed or uncertain smart-if
   answer goes to a human (O-71).
8. **No LAYUP service.** The only process of LAYUP is `layup run` while it runs,
   plus the scheduled dead-man job. There is no daemon, no database and no
   queue (ADR-0011, decision 1).

We reject: short `layup run --once` steps from a Git work queue or from forge
events (the clock is off between events, and the harness sessions and their
credentials move onto the forge runner); and a durable workflow engine with its
own store (the state leaves Git, so a second harness on a clean clone cannot
read it; Invariants 1 and 9).

## Consequences

- A human no longer starts, moves or restarts the role agents in the normal
  flow; each restart that a human causes is counted (Task Intervention Rate,
  `F-0003#70`).
- The host that runs `layup run` must stay on during delivery. Its loss is
  detected by the dead-man check within `2H` plus one schedule period, not
  prevented.
- The records branch, the harness register and the parameter register are new
  files that the implementation plan (#42) specifies with their headers; the
  architecture (`docs/architecture.md`) names them.
- `T`, `N` and `H` are parameters with defaults, not calibrated values; the
  pilot records the measurements that calibrate them (`F-0001#4`).
