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

## Review rounds

The records and the Fixes reply are comments on #156. Round 1 (Claude Fable
5.1; `655c841`, cycle 0): `material`, one finding (the red record held no
mutation for four rules), fixed in `9c5d344` (a described mutation of each of
the fourteen rules; notes 2 and 3: the glossary, `TestThePairOfASession`).
Round 2 (Fable; `9c5d344`, cycle 1): `nothing material in scope`, three notes,
applied in the close-out: the record says its lines are cut and has a green of
the fix; the Test cell names the fourth test; "or that has no probe at its
version" in `session.md` and the glossary.

## Verdict

Delivered: admission and the order of the routing register in
`internal/route`. The review ended by decay at cycle 1, the cap. The diff
against `5fbe7e5`, the branch's base, is inside 600 lines over 12 files (note 4
of round 1). Next: row 36a calls `Pair` and `Admitted`.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-09; reviewers' tokens are
`modelUsage`; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 13:15 to 13:17 |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 728,374 (USD 4.19) | 6 min 55 s |
| The code, test first | execution | Claude Opus 5.5 | not reported | not reported | 13:34 to 13:39 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 879,416 (USD 4.21) | 8 min 2 s |
| The fix of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 13:50 to 13:53 |
| Round 2 | reasoning | Claude Fable 5.1 | `xhigh` | 600,640 (USD 3.44) | 8 min 9 s |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 14:03 to 14:06 |
