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

**The demo** is shown at the unit level, as the Tests cell and the acceptance
row say, by `TestTheCapOfMilestoneZeroDrawsNoSession` on `PhaseSession` with
the stand-in store (no session ID drawn, no commit), and at the integration
level by `TestTheCapOfMilestoneZeroStartsNoSession` on the real records branch:
no step of the restart calls `PhaseSession` yet (rows 43b to 47b do). The error's first line is
the fixed text "the cap of milestone 0", and its detail holds no clock reading
(round 1 of #216), so each later step's `fail(name, err)` gives row 7's detail,
the same bytes on a repeat, with no mapping of its own.

## Known limits

- **The evidence of the unit test** (note 1 of round 2): with the check of
  `PhaseSession` cut, `TestTheCapOfMilestoneZeroDrawsNoSession` fails by a
  panic at a function of `Sessions` that the stand-in leaves unset, not at its
  own assertions on the session IDs drawn and the commits; a change that sets
  those functions in the stand-in should run that mutation again.

- **A probe after the rule** (note 4 of the plan review): the session's own
  version check may run a probe after the rule read the known money; its spend
  counts from the next session, as `intake.md` now says.

## Review rounds

The records are comments on #216. Round 1 (Claude Fable 5.1, effort `xhigh`, a
fresh read-only session in a clone at `a8adfaf`, cycle 0): `material`, two
findings: the hours refusal wrote a clock reading into the step's detail, and
the demo had no unit test, which the Tests cell names. The fix `0072ab5` (cycle
1) made the detail fixed text and added D2 at the unit level, red first, with a
lesson of §2 (note 3) and the doc comment of note 4. Round 2 (the same reviewer
and effort, a fresh session at `0072ab5`, cycle 1, the cap): `nothing material
in scope`, four notes: note 1 is a known limit above; note 2 is the run below;
note 3, a quoted text of the lesson, is #240; note 4, the headroom, holds at
landing.

## Verdict

Delivered: the budget of milestone 0, asked before each session of Intake: the
known money since Start plus the new session's cap against `MONEY` in US
dollars, and the hours since Start against `HOURS`; past either, no session
starts, with a fixed detail. The review ended by decay at cycle 1 of cap 1. The
diff against `cafebd3`, the branch's base, is inside 700 lines over 20 files.
At `0072ab5` (06:28Z to 06:29Z): `go build ./...`, `go vet ./...`, `go test
./...`, `go test -tags=integration ./...` and `go test -count=1 -tags=e2e ./...`
pass, every package `ok`. Next: row 41b.

## Resource record

Recorded, not budgeted (ADR-0007). UTC; reviewers' tokens are `modelUsage`;
the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 2026-10-11 05:39 to 05:40 (the claim to the post) |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,467,128 (USD 6.80) | 8 min 47 s |
| The code, test first | execution | Claude Opus 5.5 | not reported | not reported | 05:49 to 06:00 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 1,462,491 (USD 5.30) | 12 min 58 s |
| The fix of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 06:13 to 06:18 |
| Round 2 | reasoning | Claude Fable 5.1 | `xhigh` | 937,953 (USD 4.87) | 10 min 10 s |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 06:28 to 06:30 |
