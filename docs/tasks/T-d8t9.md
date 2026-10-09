# T-d8t9 — row 36a of the plan, the start of a task session and its refusals

Issue: [#164](https://github.com/pharzam/layup/issues/164), row 36a of [the tasks
of M2b](../plan/README.md#the-tasks-of-m2b), opened by `T-fdaq` (#152). Serves
`F-0003#52` through `REQ-013` and `NFR-001`. Base `53f1bef` (the merge of #188,
row 34a); the budget is read against it. Author: Claude Opus 5.5 on Claude
Code. Evidence: [`runs/T-d8t9/`](../../runs/T-d8t9/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #164. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session at
`53f1bef`, 6 min 39 s) gave `approve-with-conditions`: Budget maximum 1,200
lines added plus removed over 16 files, close-out inside; Cycle cap 1; no
panel. Its three conditions are applied: the sweep before the ID, the
directory right after it, the version check in `repo/` with no credential, and
a refused start's directory removed at once; the attempt number, the empty
table and the one commit of the start row written in `session.md`; the start
row and the refusal through `Fenced`, with one refusal retried and two a
`LostError`. The goal count is 2, final by O-187 of #152.

## What was done

1. **D1:** `Sessions.StartAttempt`: the event `attempt`, its number one more
   than the largest attempt of an event `attempt`.
2. **D2:** `Sessions.TaskSession`, `TaskSpec`, `Pair`: the hooks of
   `session.Call` for a task session; the start row and the event `session`
   in one records commit; a refused start, one event `refused` and no start
   row; `SessionStore`, which the records store of row 25a now meets
   (`ReadFile`, `Commit`).
3. **D3:** the demo at integration: a fake harness sees its start row on
   `layup-records` of the target.
4. **D4:** `session.md`, `packages.md`, `traceability.md`, the Test cells of
   `REQ-013` and `NFR-001` and a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-d8t9/test-runs.md): red before the
functions; a mutation of each of the twenty rules, each caught, two of them
added by the fix of round 1; green after.

**A deviation** (note 3 of round 1): the integration test of the demo was
written after `session.go`, so it was not red before the code; its red is the
mutation `startrow-order`, which shows that it fails on its own rule.

**For the first caller** (note 5 of the plan review): `TaskSession` takes the
pair as an input. The read of each harness's version for `route.Pair`
(`internal/route/admit.go`) lands with the first caller of `TaskSession`, and
the admission hook (`Admit`) is row 38's to fill with the probe. No step of
`layup run` builds a `Sessions` yet: rows 38 and 39b wire it to the run's
records store.

**The rejected alternatives:** two commits for the start row and the event;
the refusal written from inside `internal/session` (one writer: the records
store is `internal/run`'s, which `internal/session` cannot import).
