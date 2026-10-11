# T-gmey — row 42a of the plan, the session of a phase

Issue: [#215](https://github.com/pharzam/layup/issues/215), row 42a of [the tasks
of M2c](../plan/README.md#the-tasks-of-m2c), opened by `T-2r8u` (#211). Serves
`F-0003#41` and `F-0003#51` through REQ-001, REQ-012 and NFR-001. Base `ad332a7`
(the merge of #234, row 43a); the budget is read against it. Author: Claude Opus
5.5 on Claude Code. Evidence: [`runs/T-gmey/`](../../runs/T-gmey/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #215. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `ad332a7`, 7 min 3 s) gave `approve-with-conditions`: Budget maximum
1,300 lines added plus removed over 20 files, close-out inside; Cycle cap 1; no
panel. Its two conditions are applied: `input/` is written by `internal/session`
(`MakeInput`), and `Keep` takes the session ID. The goal count is 1, final by
O-196.

## What was done

1. **`internal/session`:** `Dir.Input`; `MakeInput` writes `input/` right after
   the directory, each file 0600, byte for byte, a name that is not one file
   name refused; `ReadOutput` reads `result/<name>` by the rules of a result
   file, which `ResultFile` now shares (`readResult`), so the rules have one
   home.
2. **`internal/run`:** a session of a phase is a third kind of the one call
   `TaskSession`, beside a task session and a probe: `PhaseSession` draws its
   task ID as a probe's, starts with the sweep and the probe of a new version,
   writes its start row alone (no attempt, no event, no push), builds its prompt
   from the absolute paths of its input and output files, and at its end reads
   its output file, runs the step's `Keep`, and makes one records commit of the
   kept records and its telemetry row (the telemetry row alone for a class that
   is not `done` or a refused file), then its comment. `PhasePair` gives the
   first admitted pair of the role and the tier `reasoning`; with none,
   `route.ErrNoPair`. A refused start, a class that is not `done` and a refused
   file are a `*PhaseError` that names the session.
3. **The documents:** `intake.md` (one reading of `input/`; the version of the
   pair; the telemetry row alone; the comment; a row of Input states for a
   refused start), `session.md` (a pointer), the Job cells of `internal/run` and
   `internal/session`, `traceability.md`, the Task and Test cells of REQ-001,
   REQ-012 and NFR-001 and a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-gmey/test-runs.md): red before the code;
16 of 16 mutations caught; green at the three levels.

## Known limits

- **The version of a pair** (note 6 of the plan review): `PhasePair` reads a
  harness at the version of its last row of `harnesses.tsv`. A last row that is
  a refused start (version `—`) refuses a pair that the step would otherwise
  take, until the next step `probe` writes a row of that harness; a failed probe
  at a version the harness has since left is corrected by the session's own
  version check, which probes the new version.
