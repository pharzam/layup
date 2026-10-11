# T-ybd5 — row 42b of the plan, the budget of milestone 0

Issue: [#216](https://github.com/pharzam/layup/issues/216), row 42b of [the tasks
of M2c](../plan/README.md#the-tasks-of-m2c), opened by `T-2r8u` (#211). Serves
`F-0003#50` through REQ-011 and NFR-001. Base `cafebd3` (the merge of #237, row
42a); the budget is read against it. Author: Claude Opus 5.5 on Claude Code.
Evidence: [`runs/T-ybd5/`](../../runs/T-ybd5/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #216. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `cafebd3`, 8 min 47 s) gave `approve-with-conditions`: Budget maximum
700 lines added plus removed over 20 files, close-out inside; Cycle cap 1; no
panel. Its three conditions are applied: the rule is in `internal/ledger`, which
`internal/run` calls; `MONEY` is US dollars, and a known row in another currency
refuses on the safe side; `CommitTime` reads `%ct`. The goal count is 1, final
by O-196.

## What was done

1. **`internal/ledger`:** `MilestoneCap` and `ErrMilestoneCap`: the known money
   (the rows of `telemetry.tsv` whose money is `reported` or `computed`, the
   probes' included, in exact decimals) plus the new session's cap against
   `MONEY`, and the time since Start against `HOURS`; a sum or a time passes
   when it is above its cap; a known row not in `USD` refuses.
2. **`internal/git`:** `CommitTime`, `git log -1 --format=%ct`, read as seconds
   since the epoch.
3. **`internal/run`:** `SessionStore.Started`, the committer time of the root
   commit of the run's records base; `PhaseSession` asks the rule before it
   draws a task ID, with `intake.cap` of `start.tsv`, `telemetry.tsv`, the time
   of Start, `Now` and the pair's cap, so a refused session leaves no record.
4. **The documents:** `intake.md` (the decisions of the task, and a row of Input
   states for a known row not in `USD`), `packages.md` (the Job cells of
   `internal/ledger` and `internal/run`, the ledger's row of the later phases,
   and the row of `CommitTime` in the calls of `internal/git`),
   `traceability.md`, the Task and Test cells of REQ-011 and NFR-001 and a §13
   line.

**Tests:** [`test-runs.md`](../../runs/T-ybd5/test-runs.md): red before the code;
12 of 12 mutations caught; green at the three levels.

**The demo** is shown at the integration level by
`TestTheCapOfMilestoneZeroStartsNoSession` on `PhaseSession`: no step of the
restart calls `PhaseSession` yet (rows 43b to 47b do). The error's first line is
the fixed text "the cap of milestone 0", so each later step's `fail(name, err)`
gives row 7's detail with no mapping of its own.

## Known limits

- **A probe after the rule** (note 4 of the plan review): the session's own
  version check may run a probe after the rule read the known money; its spend
  counts from the next session, as `intake.md` now says.
