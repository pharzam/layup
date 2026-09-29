# 0024. Record every session's cost and stop at the milestone cap

Date: 2026-09-29

## Status

Proposed

## Context

The PSB asks that each task get a record of "its token count, its latency, and its
wall-clock duration" in the repository (`F-0003#50`, `#60`), names the budget a
business-forking decision of the idea owner (`F-0001#26`), and warns that "a
retry loop can consume a budget before a human sees it" (Problem 5). The Operator
decided that the idea owner writes a band before delivery, that the smart-if
decides inside it, and that outside it escalates (O-68); and that a harness with
no token report or spend cap may do paid work under a wall-clock limit, with its
telemetry counted as incomplete (O-80).

The deep check found that the money stop could not sum a cost with no tokens
(Sol-11); that O-80's default fails Telemetry Completeness by design (Sol-12,
Fable-N4); that the budget asked for milestone and task numbers before they exist
(Fable-M11, Author-7); and that "an action" had no unit (Fable-N3). The review of
#69 found no stop before the spend (B7). The evaluation gave a ledger where an
unknown cost is never zero, budgets on scopes that exist, and a hard stop that
kills the session (Paperclip, AgentPlane, GNAP), and Shape Up's appetite. The
options compared are in
[`runs/T-hbw8/selection-v2.md`](../../runs/T-hbw8/selection-v2.md), row 0024.

## Decision

We will record every session's cost and stop at the milestone cap:

1. **One ledger row per session** (the unit of an action): harness, model,
   billing type, start, latency, duration, tokens by class with a status
   (`observed`, `partial`, `unavailable`), and money with a status (`reported`,
   `computed` from a cited price, `unknown`). An unknown value is never zero.
2. **The appetite and the band** (`B`, `U`) at Intake, from the idea owner, in
   `budget.tsv`; no per-task budget.
3. **A cap per milestone,** in money and wall-clock, fixed by its bet; the caps
   stay within `U`.
4. **Before each session start,** code sums the known spend and the spend caps of
   running sessions: past the milestone cap, the circuit breaker stops the
   milestone and kills the running sessions; between `B` and `U` with no unknown
   part, smart-if P5 may continue; an unknown part never counts as below `B`, so
   the idea owner decides; at `U`, a hard stop.
5. **A session with no token report and no cap** runs under a wall-clock limit,
   and its task's telemetry is incomplete (O-80).

We reject: per-task estimates at Intake; a cost that is zero because it is not
known; a harness's own price estimate judged by the smart-if.

## Consequences

- Telemetry Completeness fails for each task that a harness without a token report
  runs; the pilot can avoid it only by routing (known limit L-G1 of
  [`architecture.md`](../architecture.md)).
- Spend inside one session stops only where the harness enforces its cap; without
  one, the wall-clock limit is the only stop.
- The measures of PSB §7 read the ledger with the records of ADR-0014.
