# T-ax3r — row 25b of the plan, internal/run: Start and the restart

Issue: [#142](https://github.com/pharzam/layup/issues/142), opened by `T-trej`
(#130) when O-173 (a) split row 25 of [the tasks of
M2a](../plan/README.md#the-tasks-of-m2a) into 25a and 25b. Serves `F-0003#42`
through `NFR-001`, `NFR-002`, `NFR-006` and `REQ-002`. Base `dc3b094` (the
merge of #143). Author: Claude Opus 5.5 on Claude Code. Started under O-169.
Evidence: [`runs/T-ax3r/`](../../runs/T-ax3r/).

## Plan and plan review

The plan (R12, comment 6057786107) kept the scope of O-173 (a) and named six
decisions of the specification. Its review (Claude Fable 5.1, effort `xhigh`,
on the Claude Code CLI with stream output, a fresh read-only session in a clone
at `dc3b094`, 10 min 8 s; comment 6057962144): `approve-with-conditions`, no
wider scope; six conditions (the base of a restart's first records write; the
heartbeat stopped before the release, and a lost beat that stops the step; the
bot's ID on the restart; a call that gives the branches; the period of step 8;
a Start that finds its own leftovers) and twelve notes. The author's answer
(6057963866) took the six conditions and notes 1, 2, 3, 5, 8, 9 and 11; notes 4
and 6 went to row 26 (#131, comment 6057964158); note 7 is a known limit.
Budget maximum 1,700 lines added plus removed over 20 files against `dc3b094`,
close-out inside; Cycle cap 1; no panel.

## What was done

1. **`internal/forge`, `internal/forge/github`:** `Repository` gives the names
   of the branches, every page (condition 4).
2. **`internal/git`:** `Clone` takes `Auth`, as `Fetch` and `Push` (decision
   5); `internal/setup` passes `git.Auth{}`.
3. **`internal/run`:** `Config`, `Start` (the nine steps) and `Restart` (`forge`,
   `clone`, `version`, `lease`, `phase`), each giving the rows of `run-steps`
   (`RunStepsSchema`, moved to `built`); the records store on `internal/git`
   (each records commit on the run's own last pushed commit in a scratch work
   tree, or before its first push on the commit of its last read; one lock over
   each read and each commit with its push); the heartbeat in its own task,
   ended before the release and before any exit, and a lost beat that stops the
   step; the texts of the README and of the two issues.
4. **A defect of row 25a, found on this row's path:** `Copy` and `Decision` used
   `—` for "no App" in memory, while `tsv.Read` gives `""`; a decision read back
   from `copies.tsv` would not count on a restart. Fixed test first.
5. **The documents:** `run.md` (the six decisions; conditions 1, 2, 3, 5, 6;
   the known limit of note 7; the blocks of the two issues), `forge.md`
   (`Repository`; the client and the progress function of the adapter),
   `packages.md` (`Clone` with a token), the traceability, the Test cells of
   `NFR-001`, `NFR-002`, `NFR-006`, `REQ-002` and a §13 line.
   [`docs.sh`](../../runs/T-ax3r/docs.sh) checks them.

**Tests:** [`test-runs.md`](../../runs/T-ax3r/test-runs.md): red, then green,
for each part; three defects found by the integration tests; two mutations for
the tests of conditions 1 and 2, each caught; one flaky test helper, fixed, then
eight clean runs with `-race`.

**The rejected alternatives:** the steps in `internal/cli`; a reuse of
`internal/setup` (not in the May import cell); a stand-in forge with no HTTP for
the integration test.

## Review rounds

Round 1 (Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI with stream
output, a fresh read-only session in a clone at `57de0de`, 10 min 37 s; comment
6058465795): `material`. Finding 1: the restart's step 8 ran for `watch.T` of
the command, which a restart does not have (0: one read, then
`not-confirmed`). Fixed in `eb45c56` (comment 6058528738): `clone` takes
`lease.H` and `watch.T` from `start.tsv`, test first; notes 2, 3 and 6 applied
(one antecedent in the fencing sentence; the read into the run's clone, not a
scratch repository, recorded on the issue; no second announcement of
`opening`).

Round 2 (the same set-up, at `eb45c56`, 9 min 44 s; comment 6058692643):
`nothing material in scope`, six notes. Note 2 (no progress line at the start
of each step) is written on #131 for row 26 (comment 6058693030). No other note
changes the reviewed head.

**Known limits** (notes of rounds 1 and 2, not applied):

- An account whose default branch is not `main`: step 4 waits for that branch,
  while the printed command pushes `main` (round 1, note 4).
- A relative `--host` turns each `.` of a detail into `DIR`, and a transport
  failure of a records push can put the ID of a records commit in a detail
  (round 1, note 5); row 26 passes an absolute `DIR`.
- Rows 10 and 12 of the input states have no test in `internal/run`; row 10 is
  tested by `records.ReadLease` through the same reader (round 1, note 7).
- A `Show` of `copies.tsv` that fails for a reason other than a missing file is
  read as no copies (round 1, note 8).
- Two readings of the new sentence of `forge.md` on the client and the progress
  function, and restart step 2 does not name `copies.tsv` (round 2, notes 1
  and 3).
- The test of finding 1 proves `watch.T` of `start.tsv`, not `lease.H`
  (round 2, note 5).
- A token renewal that the end of the heartbeat cancels reads as a lost beat,
  so step 9 can fail with "context canceled" (round 2, note 6; one renewal an
  hour against one beat each `lease.H`).

## Verdict

Delivered: against the `httptest` forge and a local bare repository, Start
makes the first records commit with the pin, the briefs, `approvers.tsv`,
`start.tsv` and the lease row, then the two issues, the watch and the release;
the restart rebuilds the clone, checks the version, takes the lease (over a
stopped run too) and runs the steps of Start that are not done. The review
ended by decay at cycle 1. The diff against `dc3b094` is inside 2,100 lines
over 26 files (O-174 a).

Next: row 26 (`T-mqty`, #131), the command `layup run`, After 25b.

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC on 2026-10-08; tokens are the
`result` event of the Claude Code CLI (stream runs); `not reported` otherwise.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The reading of the specification; the plan | reasoning | Claude Opus 5.5 | max | not reported | 10:15 to 10:22 |
| Its plan review | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,551,743 (USD 6.91) | 10 min 8 s, from 10:22 |
| The answer; the work, test first; the documents | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 10:32 to 10:50 |
| O-174 | reasoning | Claude Opus 5.5 | max | not reported | 10:50 to 10:52 |
| Round 1 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 999,933 (USD 5.56) | 10 min 37 s, from 10:53 |
| The fix of round 1 | execution | Claude Opus 5.5 | max | not reported | 11:04 to 11:07 |
| Round 2 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,926,754 (USD 5.98) | 9 min 44 s, from 11:07 |
| The close-out | execution | Claude Opus 5.5 | max | not reported | 11:17 to 11:20 |
