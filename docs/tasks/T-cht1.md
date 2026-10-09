# T-cht1 — row 30b of the plan, admission and the routing order

Issue: [#156](https://github.com/pharzam/layup/issues/156), row 30b of [the tasks
of M2b](../plan/README.md#the-tasks-of-m2b), opened by `T-fdaq` (#152). Serves
`F-0003#52` through `REQ-013`. Base `5fbe7e5` (the merge of #174; the plan was
reviewed at `5cd8a18`). Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-cht1/`](../../runs/T-cht1/).

## Plan and plan review

The plan (R12, comment 6081666992), its review (6082020793) and the author's
answer (6082021438) are comments on #156. The plan review (Claude Fable 5.1,
effort `xhigh`, a fresh read-only session in a clone at `5cd8a18`, 6 min 55 s)
gave `approve-with-conditions`: Budget maximum 600 lines added plus removed
over 12 files, close-out inside; Cycle cap 1; no panel. Its three conditions
are applied: the `—` session of a refused `pair` is cited from row 28's cases;
the rule of a passed probe has one home, Admission, and step 2 and the step
`probe` point to it; where the caller gets a harness's version is row 36a's.
The goal count is 2, final by O-187 of #152.

## What was done

1. **D1:** `Probed` (the last row of the harness at the version, in the order
   of the file, passed) and `Admitted` (and the model's row has `use` `yes`),
   in `internal/route/admit.go`.
2. **D2:** `Pair` (the first admitted pair of the role's list for the tier, by
   `position`) and `ErrNoPair`, the refusal `pair`.
3. **D3:** `session.md` (Admission the one home of the rule; step 2 and the step
   `probe` point to it), `packages.md`, the package comment, `traceability.md`,
   the Test cell of `REQ-013` and a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-cht1/test-runs.md): red before the
functions (no compile); a mutation of each rule, each caught (one case added
when a mutation first passed); green after.

**The rejected alternatives:** admission by the last row of the harness at any
version; a closed list of roles here; the version check in `internal/route`.
