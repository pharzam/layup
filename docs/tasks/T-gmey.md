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
  take, until the next step `probe` writes a row of that harness, and so does a
  last row `failed` at a version the harness has since left. A last row
  `passed` at a version the harness has since left is corrected by the
  session's own version check, which probes the version it reads (round 1 of
  #215, note 2).

## Review rounds

The records are comments on #215. Round 1 (Claude Fable 5.1, effort `xhigh`, a
fresh read-only session in a clone at `a226a6c`, cycle 0): `material`, one
finding: the sentence of `intake.md` that placed `input/` "before the prompt"
read two ways, and under one the tree made it false. The fix `3031c8f` (cycle 1)
rewrote it to one reading in words, corrected the known limit (note 2) and the
version read (note 4), and made a failed comment keep the session's reason
(note 5), red first. Round 2 (the same reviewer and effort, a fresh session at
`3031c8f`, cycle 1, the cap): `nothing material in scope`, two notes on doc
comments that the fix left in the old words, now #238.

## Verdict

Delivered: the session of a phase, a role session of Intake with no attempt and
no event, run from a directory that holds `input/`, its output file read by the
rules of a result file and kept by its step, its records and its telemetry row
in one records commit, then its comment; the pair of a role, with none no
session. The review ended by decay at cycle 1 of cap 1. The diff against
`ad332a7`, the branch's base, is inside 1,300 lines over 20 files. Next: row
42b.

## Resource record

Recorded, not budgeted (ADR-0007). UTC; reviewers' tokens are `modelUsage`;
the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 2026-10-11 04:45 to 04:50 |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,028,558 (USD 5.28) | 7 min 3 s |
| The code, test first | execution | Claude Opus 5.5 | not reported | not reported | 04:57 to 05:11 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 1,543,696 (USD 5.75) | 10 min 26 s |
| The fix of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 05:22 to 05:26 |
| Round 2 | reasoning | Claude Fable 5.1 | `xhigh` | 789,359 (USD 5.02) | 8 min 6 s |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 05:35 to 05:35 |
