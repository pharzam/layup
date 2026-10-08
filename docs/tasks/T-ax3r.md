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
