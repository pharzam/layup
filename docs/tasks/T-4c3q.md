# T-4c3q — row 35 of the plan, the writer of the telemetry record

Issue: [#163](https://github.com/pharzam/layup/issues/163), row 35 of [the tasks
of M2b](../plan/README.md#the-tasks-of-m2b), opened by `T-fdaq` (#152). Serves
`F-0003#50` through `REQ-011`. Base `ad8774d` (the merge of #173). Author:
Claude Opus 5.5 on Claude Code. Evidence: [`runs/T-4c3q/`](../../runs/T-4c3q/).

## Plan and plan review

The plan (R12, comment 6081096088), its review (6081332272) and the author's
answer (6081332778) are comments on #163. The plan review (Claude Fable 5.1,
effort `xhigh`, on the Claude Code CLI with stream output, a fresh read-only
session in a clone at `ad8774d`, 8 min 5 s) gave `approve-with-conditions`:
Budget maximum 650 lines added plus removed over 11 files against `ad8774d`,
close-out inside; Cycle cap 1; no panel. Its three conditions are applied: the
price rows of the session's harness and model, the newest date among them,
with a case of another harness; the two decisions written in `session.md`; a
reported cost as the text of the JSON number, written in the form of the type.
The goal count is 1, final by O-187 of #152.

## What was done

1. **D1:** `internal/ledger`: `Session`, `Usage` and `Row`: the row of
   `telemetry.tsv` of one session, its times to the second, its tokens, and its
   money `reported`, `computed` (exact, with `math/big`) or `unknown`, checked
   by `records.CheckTelemetry` before it is given. `Row` reads no file and no
   clock; the records commit of the session's end is row 36b's, and the probe's
   row row 38's.
2. **D2:** `session.md` (the price rows, the reasons of `unknown`, the exact
   decimal, the form of a reported cost), `packages.md`, `traceability.md`, the
   Test cell of `REQ-011` and a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-4c3q/test-runs.md): red before the
package (no compile); a mutation of each rule, each caught by its own cases
(one equivalent, said there); green after.

**Decided here** (note 6 of the plan review): one model in the report that is
not the start row's gives `unknown`.

**The rejected alternatives:** the writer in `internal/run`; a money computed
from partial tokens; a rounding rule.
