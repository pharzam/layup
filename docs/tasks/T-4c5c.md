# T-4c5c — row 53 of the plan, the step markers

Issue: [#276](https://github.com/pharzam/layup/issues/276), row 53 of [the tasks
of M2d](../plan/README.md#the-tasks-of-m2d), opened by `T-bng8` (#267). Serves
`F-0003#42` and `F-0003#63` through REQ-002 and NFR-003. The plan was made against
`9144e4e` (the merge of #305); the branch took `38ee78e` (the merge of #306, `T-9v92`)
by a merge, and the budget is read against that, as the change of this task alone.
Author: Claude Sonnet 5.5 on Claude Code. Evidence: [`runs/T-4c5c/`](../../runs/T-4c5c/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #276. The plan
review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at
`9144e4e`, 8 min 10 s) gave `approve-with-conditions`: Budget maximum 1,800 lines
added plus removed over 30 files, close-out inside; Cycle cap 1; no panel. Its three
conditions are applied: (1) with `S01-stack` always asked, the Operator's block asks,
so the tests of rows 47a to 47c answer as the Operator too, at the integration level
and, as the answer added, in the end-to-end tests of `cmd/layup`; (2) the rows of the
session's file may be in any order, and the record is in the scan's order, written into
"The step markers"; (3) the follow-up's form is written into its own home in
`intake.md`. The answer named `internal/run/demo_uat_test.go` as a demo of M2c to read:
it is the demo of M2b and holds no step of Intake, so it did not change. The goal count
is 2, final by O-209 of #267.

## What was done

1. **The scan:** `setup.MarkerPlaces` (a thin export of S10's own grouping, so the rule
   "one place per file and key, at its first line, in the order of the stop table" has
   one home) and `pinnedPlaces` in `internal/run/markers.go`: the one root commit of
   the default branch, its tree equal to `pin.tree`, read in a scratch work tree by the
   scanner of S10 and removed after.
2. **The session and the record:** one session of `systems-architect` on
   `input/problem-statement.md` and `input/markers.tsv` (the list, in the form of the
   block with `value`, `start`, `end` and `session` as `—`); it writes `result/markers.tsv`
   (`id`, `value`, `quote`); `markersKeep` refuses a key that the list lacks, a key twice,
   an empty value or quote and a quote that `quoteAt` refuses, and gives
   `intake/markers.tsv`, one row per marker in the scan's order, with the offsets and the
   session. With no marker, the table has no row and no session starts.
3. **The batch and its comment:** after the intent form, `S01-stack` and one question for
   each row of `intake/markers.tsv`, in the Operator's block, with the source `setup`,
   never dropped as a repeat. `questionsComment` writes the two blocks, a marker's
   question with `proposed:` and its quote, the text of an empty block and the three
   forms of a marker's answer; the follow-up shares it, so it has the same form.
4. **The rows of a setup question:** `internal/records` gains `setup` in the enum of
   `source` and its row rules (the ref, the kind, the asker, the block, no quote for
   `S01-stack`).
5. **The restart:** the step in its order and in `restartNames`, and its name in the block
   `run-steps`.
6. **The tests of the earlier steps:** the lists of steps and the progress counts, and the
   answers of rows 47a to 47c and of the end-to-end tests, now answer the Operator's block.
7. **The documents:** `records.md` (the block `questions`), `intake.md` (the merge, the
   comment and the follow-up), `scaffold.md` (the readings of the step and of the setup
   questions, and the level of its acceptance test), `run.md`, `packages.md`,
   `traceability.md`, the cells of REQ-002 and NFR-003 and a line of §13 of `PRD-0001`,
   and two pitfalls in `guardrails.md` §2 (a step added to a list of steps, and the integration
   level of a package that needs minutes on a loaded host).

**Tests:** [`test-runs.md`](../../runs/T-4c5c/test-runs.md): red before the code (the
enum of `source` refuses the rows, and `internal/setup` and `internal/run` do not build);
each rule of the new code mutated alone on a copy of the tree, 61 of 61 caught (a first
run caught 54, and the other 7 mutants did not build, an unused variable or import; they
were repaired, and one of them showed that the scripted harness did not read the list
that the session is given, which it now does); every check of the gate exits 0 on the code.
On the merged head the integration level of `internal/run` was stopped by the 10 minutes that
`go test` allows (600.024 s, the host's load average about 6 on 2 CPUs, no test had failed); it
was run again alone with `-timeout 25m` and passed (486.611 s, 126 tests), and the other
packages' integration level passed (exit 0); CI ran the three levels in 3 min 22 s.

**Known limits:** a proposal's value is not checked for an angle quote, as the
specification names three refusals only; S11 refuses such a value when the setup takes
it (row 54). A batch that a run of this version posted before the step existed (a
development build) gets the record and no setup question. The marker's key keeps a tab,
which the record writes as a space, as `T-tetz` noted. The step reads the root commit's
tree in a scratch directory under `TMPDIR`. The answers of the setup (`accept`, a value
and `gap`, and the rows of `answers.tsv`) are row 54's.

**The rejected alternatives:** a second scan or grouping in `internal/run` (two homes for
the key rule); the scan over the working tree of the clone (it is `--no-checkout`) or over
the head of the default branch (after Scaffold 6 it is the setup head); one question for
all markers; the proposal's value in `questions.tsv` (it has no column); the setup
questions among the idea owner's, with no block of their own.

## Review rounds

Round 1 (GPT-6.1 Sol on Codex CLI, effort `xhigh`, a fresh session in a clone at `243418f`,
cycle 0 of the cap of 1; its record is comment 6142489140 on #276): `nothing material in
scope`, two notes. The harness `codex` is not one of the three that ADR-0026 names (`claude`,
`devin`, `opencode`): the Operator added Codex CLI to the host on 2026-10-13 and said it may be
used, and the posted records say so. A first start of the session (09:12:49Z to 09:13:55Z,
27,737 tokens) ended before it read a file, because its sandbox refuses a run directory that is
reached through a symlink; it gave no record and none was posted, and the second start ran on
the real path.

- **Applied** in `29bbcce`, a commit that no round read (the label of one red run in the test
  record, no code): note 2 (the red run of the list of the restart's steps names
  `TestStartAndRestartCallStepAtTheStartOfEachStep`, and says that the package did not build).
- **No change** for note 1: the reviewer's sandbox opened no socket and found the parent
  repository from its temporary directories, so four packages of the integration level and
  `TestFlagged` failed there, none of them by a change of this branch; the level exited 0 on
  this host (09:02Z to 09:10Z), and CI is the authority.

Round 2 (GPT-6.1 Sol on Codex CLI, effort `high`, a fresh session in a clone at `aa6c834`;
comment 6142686858): `origin/main` at `38ee78e`, the merge of `T-9v92`, was merged into the
branch after round 1, and it changed five files that this branch touched (`guardrails.md`,
`traceability.md`, `packages.md`, `scaffold.md` and `PRD-0001`, whose three conflicts were
resolved by keeping each task's rows and both lines of section 13). So a round scoped to those
files followed, as [the discipline](../engineering-discipline.md#reviewing-until-findings-decay)
says; it consumes no cycle, and its record has the cycle `1 — a round scoped to the merge of
origin/main`. The verdict row of round 1 was edited to `material`, with that reason, so that the
two records read in their positions (precedent: round 8 of `T-dm8m`). Round 2 found nothing: `nothing
material in scope`, no note.

## Verdict

Delivered: the step `markers` of Intake and the setup questions of the batch. A restart
scans the pinned tree, has one session of `systems-architect` propose values with quotes that
code finds once, and commits `intake/markers.tsv`; the batch asks `S01-stack` and each marker in
the Operator's block, and its comment and the follow-up show the proposals. The tests were red
first, each rule was mutated alone and caught (61 of 61), and every check of the gate exits 0.
The review ended by decay at cycle 0 of cap 1, and the round scoped to the merge of `main` found
nothing. The diff against `38ee78e` is inside 1,800 lines over 30 files (1,749 lines added plus
removed over 28 files, measured with the close-out inside). Next: row 54, `T-eahr`, writes
the step `setup` and the answers of the setup on these questions, once row 52c is done.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-13; reviewers' tokens are as the harness
reports them; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Sonnet 5.5 | `max` | not reported | 08:00 to 08:13 (the claim to the post) |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,164,041 (USD 5.93) | 8 min 10 s, 08:13 to 08:21 |
| The tests and the code, test first | execution | Claude Sonnet 5.5 | `max` | not reported | 08:23 to 09:12 (the answer to the pull request) |
| Round 1 | reasoning | GPT-6.1 Sol (Codex CLI, `codex exec`; the harness exception is above) | `xhigh` | 210,342 (as Codex reports; no USD) | 9 min 54 s, 09:14 to 09:24 |
| The note, then the merge of `main` | execution | Claude Sonnet 5.5 | `max` | not reported | 09:26 to 09:30 |
| Round 2, scoped to the merge | reasoning | GPT-6.1 Sol (Codex CLI, `codex exec`) | `high` | 174,461 (as Codex reports; no USD) | 6 min 11 s, 09:31 to 09:37 |
| The close-out | execution | Claude Sonnet 5.5 | `max` | not reported | 09:46 to 09:59 (the commit of this file) |
