# T-5pxd — row 34a of the plan, the process of a session and its stop

Issue: [#161](https://github.com/pharzam/layup/issues/161), row 34a of [the tasks
of M2b](../plan/README.md#the-tasks-of-m2b), opened by `T-fdaq` (#152). Serves
`F-0003#52` through `REQ-013`. Base `db4651a` (the merge of #187, row 33b).
Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-5pxd/`](../../runs/T-5pxd/).

## Plan and plan review

The plan (R12, comment 6082571847), its review (6085502029) and the author's
answer (6085503142) are comments on #161. The plan review (Claude Fable 5.1,
effort `xhigh`, a fresh read-only session at row 33a's head `0e8f15b`) gave
`approve-with-conditions`: Budget maximum 1,050 lines added plus removed over
15 files, close-out inside; Cycle cap 1; no panel. Its four conditions are
applied: the work waited for row 33b's `Refusal`, and adds the reason
`version`; the frame of the call has a hook after the end, for a task session
only; `{prompt}` of a `stdin` row is the path of `prompt.md`, written in
`session.md` (a `file` or `arg` command with no `{prompt}` is #186); a fake
for each rule of the limit. The goal count is 2, final by O-187 of #152.

## What was done

1. **D1:** `Version` (the first line, trimmed; a non-zero exit or an empty line
   `Refusal{version}`), `Words` (the placeholders in one pass), `Prompt` (the
   value of `{prompt}` by the row's mode).
2. **D2:** `Spec`, `NewSpec` (the real limits: 64 MiB, 8 MiB, 10 s), `Run`,
   `Process`: a process group of its own, the input by the mode, the stop at
   `wall` or an output cap (`SIGINT`, `SIGTERM`, `SIGKILL`, each only while the
   process has not ended), the outputs cut at their caps, the first byte of
   `stdout` timed, `ErrStart` for a program that cannot start.
3. **D3:** `Steps` and `Call`: the one call's order, each refusal ending it, the
   push hook after the end for a task session.
4. **D4:** `session.md` (the `stdin` reading, the end of a process, the line of
   the line cap, the version check moved to the integration row), `packages.md`,
   the plan's row 34a, `traceability.md`, the Test cell of `REQ-013` and a §13
   line, a lesson in `guardrails.md` §2.

**Tests:** [`test-runs.md`](../../runs/T-5pxd/test-runs.md): red before the
functions; a mutation of each of the twenty-five rules, each caught, two of them
added by the fix of round 1; green after.

**Not as planned:** `Process` takes no context: a cancel of `layup run` is no
stop that the specification names, so the caller's only stop is `wall`.
`Version` keeps its context. Note 3 of the plan review: the version check
starts a program, so its acceptance clause moved to the row at integration.

**The rejected alternatives:** `exec.CommandContext` for the stop (`SIGKILL`
to one process); a shell to start the command; the waits as constants.

**Linux only.** `exitedUnreaped` calls `waitid` with `WNOWAIT` through
`syscall.Syscall6`, with Linux's constants and a `siginfo_t` of 128 bytes: the
module's first `unsafe` (one `unsafe.Pointer`, to that buffer), for the release
review of row 39a to weigh. With it, and `TestProcessEndsWhenAProgramLeftTheGroup`
reading `/proc`, this row runs on Linux only; the module was Unix-only before
(`Setpgid`, `syscall.Stat_t`), and CI runs on `ubuntu-latest`.

## Review rounds

The records are comments on #161. Round 1 (Claude Fable 5.1, effort `xhigh`, a
fresh read-only session in a clone at `179646f`, cycle 0): `material`. Finding
1, material: the close of the outputs after `SIGKILL` had no case. Fixed in
`a7e1bf0` with `TestProcessEndsWhenAProgramLeftTheGroup`, with notes 2 to 5:
the leader reaped only at the end (`waitid` with `WNOWAIT`), the comment of
`Version`, the **decided here** markers and row 9, the column `end` of
`records.md`. Round 2 (the same model and effort, a fresh session at
`a7e1bf0`, cycle 1, the cap): `nothing material in scope`, four notes. Applied
in the close-out: note 1 (the `end` of `session.md` points to `records.md`, its
one home) and note 4 (the paragraph Linux only above). Notes 2 (the tail of an
output at the kill) and 3 (a limit of the version check) are for the rows that
call `internal/session`: #189.

## Verdict

Delivered: the process of a session and its stop in `internal/session`:
`Version`, `Words`, `Prompt`, `NewSpec`, `Process` and `Call`. The review ended
by decay at cycle 1 of cap 1. The diff against `db4651a`, the branch's base, is
inside 1,050 lines over 15 files. Next: row 34b.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-09; reviewers' tokens are
`modelUsage`; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | to 14:07 |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 703,040 (USD 4.75) | 6 min 40 s |
| The code, test first | execution | Claude Opus 5.5 | not reported | not reported | 17:25 to 17:41 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 1,317,183 (USD 6.19) | 13 min 34 s |
| The fix of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 17:55 to 18:00 |
| Round 2 | reasoning | Claude Fable 5.1 | `xhigh` | 900,900 (USD 4.83) | 9 min 30 s |
| The close-out, with notes 1 and 4 | execution | Claude Opus 5.5 | not reported | not reported | 18:11 to 18:14 |
