# T-56dj — row 52c of the plan, the commands of S13

Issue: [#275](https://github.com/pharzam/layup/issues/275), row 52c of [the tasks
of M2d](../plan/README.md#the-tasks-of-m2d), opened by `T-bng8` (#267). Serves
`F-0003#43` through REQ-003. Base `38ee78e`, read after the merges of `main` against
`4687b39`. Author: Claude Sonnet 5.5 on Claude Code. Evidence: [`runs/T-56dj/`](../../runs/T-56dj/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #275. The plan
review (GPT-6.1 Sol on Codex CLI, `high`, 5 min 19 s) gave `approve-with-conditions`:
800 lines over 16 files against `38ee78e`, close-out inside; cap 1; no panel. Its four
conditions are applied: the shell form of the head (each `git` read its own assignment,
every operand quoted); two more cases (an empty remote head, both outputs present) under
`sh` and `bash`, one through `StepsPhase2`; four entries in phase 2, each with its own
comment; "decided here" for the choices the specification does not rule, and the documents
and the merge before the validation. The goal count is 1, final by O-209 of #267.

## What was done

`commandsText` (`setup.go`) writes `set -eu` first, in both phases; with no command the
file stays empty. `commandsS13Phase2` (`final.go`) gives four entries: the reads, the
guard on the head of the default branch and the push of `layup-setup` once (order 2); the
directory of the apply outputs; and the apply of each ruleset when its output file is
missing, through a temporary file that `mv` renames when `gh` exits 0 (order 4). A head
that is neither the root commit nor the head of `layup-setup`, an empty one included, stops
the file with exit 1 and names the three. `Phase2` gains `Host` and `DefaultBranch`, and S13
is a `fail` with no file for a host that is not absolute or an empty branch (the plan did
not name them). Documents: `scaffold.md`, `packages.md`, the traceability row, `PRD-0001`.

**Tests:** [`test-runs.md`](../../runs/T-56dj/test-runs.md): red because the package did
not build; the file under `sh` and `bash` with stand-ins of `git` and `gh`; each rule
mutated alone, 41 of 41 caught after two runs that left two alive (the quoting); the
first full run found an e2e test that did not expect the first line (fixed). No lesson is
written back to `guardrails.md`: the one trap met, a test of quoting with a branch that `sh`
prints as it is, was caught by the mutation run that every task makes.

**Known limits:** the head is read from the work area when the file runs, and a failed
apply leaves its temporary file for the next run (both "decided here"). No code calls
`StepsPhase2` yet (row 54, `T-eahr`) or prints `commands.sh` (row 56, `T-hfn1`). The
reviewer's sandbox has no local socket: the integration level is the author's run and CI's.

## Review rounds

Round 1 (`6f637f1`, cycle 0; comment 6143727150 on #275), by GPT-6.1 Sol on Codex CLI
(`xhigh`; a harness that ADR-0026 does not name, allowed by the Operator on 2026-10-13):
`nothing material in scope`, two notes, applied in `5107cfd`. The merge of `main` after it
(`8485b86`, #312) was clean and changed no file of the branch, so it needs no round.

## Verdict

Delivered: `commands.sh` starts with `set -eu`, and in phase 2 the commands of S13 push
once and apply each ruleset only when its output file is missing, so the file can be run
again after any failure. The review ended by decay at cycle 0. The diff against `4687b39`
is inside 800 lines over 16 files (786 lines added plus removed over 16 files, measured before the commit). Next: row 54, `T-eahr`.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-13, from the clock; GPT-6.1 Sol tokens
are the "tokens used" of Codex CLI; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning (not reached) | Claude Sonnet 5.5 | `max` | not reported | 09:20 to 09:23 |
| The plan review | reasoning | GPT-6.1 Sol, Codex CLI | `high` | 120,605 | 5 min 19 s, 09:22:57 |
| The tests, the code, the mutations, the documents, the green runs, the record | execution | Claude Sonnet 5.5 | `max` | not reported | 09:31 to 10:28 |
| Round 1 | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 174,470 | 10 min 29 s, 10:28:26 |
| The notes, the merge and the close-out | execution | Claude Sonnet 5.5 | `max` | not reported | 10:40 to 10:45 |
