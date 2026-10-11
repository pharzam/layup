# T-4ywt — row 41a of the plan, the span check and the package rules of M2c

Issue: [#213](https://github.com/pharzam/layup/issues/213), row 41a of [the tasks
of M2c](../plan/README.md#the-tasks-of-m2c), opened by `T-2r8u` (#211). Serves
`F-0003#51` through REQ-012, NFR-005 and NFR-007; takes note 2 of #230. Base
`e689d12` (the merge of #229), with `main` at `6f3adae` merged in before the
work; the budget is read against `e689d12`. Author: Claude Opus 5.5 on Claude
Code. Evidence: [`runs/T-4ywt/`](../../runs/T-4ywt/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #213. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `e689d12`, 8 min 19 s) gave `approve-with-conditions`: Budget maximum
800 lines added plus removed over 18 files, close-out inside; Cycle cap 1; no
panel. Its two conditions are applied: the rule `classes` fails each broken
pairing of the block `spans` (a fact with a reason among them), and
`intake.md` says so; the Task and Test cells of NFR-007 too. The goal count is
2, final by O-196.

## What was done

1. **The span check:** `internal/spec`, `Check` (the five rules, each with its
   first broken place, in the order of the block) and `Read` (a span table by a
   second Go value of `spans`); `CheckSchema`, the block `spec-check`, moved to
   `built`. **Decided here:** the rule Text skips a span that breaks Bounds, and
   compares a field of `—` as the text `—`.
2. **The package rules:** `TestPackageRules` reads "The table of M2c" in the
   form of the table of `M2b` (`readOwnTable`), and `internal/spec` is on the
   line of the engine checks.
3. **The documents:** `intake.md` (rule 5, the decided-here of the rule Text),
   `packages.md` (the owner sentence, #230 note 2; the line of the engine
   checks; the row of later phases), `records.md` (the owner sentence of
   Intake), `traceability.md`, the Task and Test cells of REQ-012, NFR-005 and
   NFR-007 and a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-4ywt/test-runs.md): D1 and D2 red for the
right reason; the engine line red with the package aside; 15 of 15 mutations
caught; green after.

**A deviation:** the schemas' integration test was written after the schemas;
its red is a mutation, caught.

**The rejected alternatives:** the check in `internal/records` (the
specification puts it in `internal/spec`); the row rules of `spans` in
`internal/spec` (a second home of row 40's rules).
