# T-m1dx — row 32 of the plan, the rule-path check before a push

Issue: [#158](https://github.com/pharzam/layup/issues/158), row 32 of [the tasks
of M2b](../plan/README.md#the-tasks-of-m2b), opened by `T-fdaq` (#152). Serves
`F-0003#43` through `REQ-003`. Base `ad8774d` (the merge of #173). Author:
Claude Opus 5.5 on Claude Code. Evidence: [`runs/T-m1dx/`](../../runs/T-m1dx/).

## Plan and plan review

The plan (R12, comment 6081078540), its review (6081326544) and the author's
answer (6081327134) are comments on #158. The plan review (Claude Fable 5.1,
effort `xhigh`, on the Claude Code CLI with stream output, a fresh read-only
session in a clone at `ad8774d`, 6 min 36 s; a first run at the wrong base was
stopped after one minute) gave `approve-with-conditions`: Budget maximum 700
lines added plus removed over 14 files against `ad8774d`, close-out inside;
Cycle cap 1; no panel. Its two conditions are applied: the exception holds only
for a diff read as hunks with an added line, "between" excluding both
boundaries, written in `session.md`; the exception is bound to its one
pattern. The goal count is 2, final by O-187 of #152.

## What was done

1. **D1:** `internal/rules`: `ReadRegister` (the one exception), `Match` (a
   directory pattern, an exact pattern, and any `.sh` path), and `Check`: first
   `workflow` over every path, then `rule-path`, in the order of the names.
2. **D2:** `GuardrailsAdditions`: the hunks of a `-U0` diff, no removed line, at
   least one added line, each strictly inside §2 of the head.
3. **D3:** a second Go value of the block `rule-paths` (its owner is
   `internal/setup`), compared with the block by `TestTheSchemaEqualsItsBlock`.
4. **D4:** `session.md` (the reading of the exception), `packages.md` and
   `setup.md` (the second Go value), `traceability.md`, the Test cell of
   `REQ-003` and a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-m1dx/test-runs.md): red before the
package (no compile); a mutation of each of the eleven rules, each caught by its
own cases; green after.

**The rejected alternatives:** moving the schema to `internal/records` (two
rows of phase 1 would change); a reader that takes the rows from the caller; a
`.sh` match only when the register holds a `.sh` entry.
