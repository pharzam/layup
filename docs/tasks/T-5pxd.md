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
