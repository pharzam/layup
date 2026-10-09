# T-3py1 — row 28 of the plan, the records of a session

Issue: [#153](https://github.com/pharzam/layup/issues/153), row 28 of [the tasks
of M2b](../plan/README.md#the-tasks-of-m2b), opened by `T-fdaq` (#152). Serves
`F-0003#52` and `F-0003#45` through `NFR-001`, `REQ-005` and `REQ-013`. Base
`8d2fa2e` (the merge of #171). Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-3py1/`](../../runs/T-3py1/).

## Plan and plan review

The plan (R12, comment 6080062383), its review (6080192815) and the author's
answer (6080193463) are comments on #153. The plan review (Claude Fable 5.1,
effort `xhigh`, on the Claude Code CLI with stream output, a fresh read-only
session in a clone at `8d2fa2e`, 7 min 51 s) gave `approve-with-conditions`:
Budget maximum 1,000 lines added plus removed over 13 files against `8d2fa2e`,
close-out inside; Cycle cap 1; no panel. Its six conditions are applied: a
wrong-header case per reader; the form of a session ID in the three records
that hold one; each "`—` when" and "`—` for" read as "exactly when"; the first
word of a failed probe's reason checked against the block's words and the
classes of the end; the form of the branch of `push` and `bound`, named in the
block in the same change; one empty-value case per column. The goal count is 2,
final by O-187 of #152.

## What was done

1. **D1:** `internal/records/session.go` holds `SessionsSchema`,
   `HarnessesSchema`, `RoutingSchema`, `EventsSchema` and `ResultSchema`;
   `TestTheSchemasEqualTheirBlocks` compares them with their blocks, and their
   names moved from `notYetBuilt` to `built`.
2. **D2:** `CheckSession`, `CheckHarness`, `CheckRouting` (the positions of
   each role and tier are 1 to k, each once, in any order of the file),
   `CheckEvent`, `CheckEvents`, `CheckEventsAppend` (no row is changed),
   `CheckResultRow` and `CheckResult` (the numbers of each kind are 1 to k, in
   any order of the file, as `position`; round 1, finding 1), run by the readers `ReadSessions`,
   `ReadHarnesses`, `ReadRouting`, `ReadEvents` and `ReadResult`.
   `session_test.go` has a valid file per record, a wrong header per reader,
   one case per column with no clause for `—`, and one case per rule; each
   refusal must name the column of its rule, so a case that fails for another
   reason fails the test.
3. **D3:** the rows of `docs/tests/traceability.md`; the Test cells of
   `REQ-005` and `REQ-013` (written) and of `NFR-001` (appended to) in
   `PRD-0001` §12, with a §13 line; the branch `task/<task>/<attempt>` in the
   block `events` of `records.md` (condition 5).

**Tests:** [`test-runs.md`](../../runs/T-3py1/test-runs.md): red with no Go
value (it did not compile); red against stubs (76 refusal cases and the five
block comparisons failed); green after.

**Known limits.** "A model of that harness in `models.tsv`" (`routing`) and
"a harness is admitted by its last row" (`harnesses`) are checks of two files:
the writer's and row 30b's, not rules of a row. `role` stays open, as the roles
are the Operator's matrix (`architecture.md`); the reason of a `refused` event
and the detail of `closed` and `rebased` stay open, as their block names no
words.

**The rejected alternatives:** the rules in `internal/run` (the specification
gives the schemas and the row rules to `internal/records`); a closed list of
the reasons of `refused` (its block gives none); one file per record.

## Review rounds

The records and the Fixes reply are comments on #153. Round 1 (Claude Fable
5.1; `c454cd8`, cycle 0): `material`, one finding (`CheckResult` held the rows
of a kind in the order of the file, which the block does not give), fixed in
`27edfc5` with notes 2 to 6. Round 2 (Fable; `27edfc5`, cycle 1): `nothing
material in scope`, two notes, applied in the close-out: the form of `NAME` of
`vars` in the block `sessions`, and `T-3py1` in the row of `internal/records`
of `packages.md`.

## Verdict

Delivered: the schemas, the rules and the readers of `sessions`, `harnesses`,
`routing`, `events` and `result` in `internal/records`, with their tests. The
review ended by decay at cycle 1, the cap. The diff against `8d2fa2e` is inside
1,000 lines over 13 files. Next: rows 30b, 34b, 36a and 36b read these records.

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC, 2026-10-09; the reviewers'
tokens are `modelUsage` of the CLI's `result` event; the author's are not
reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 11:35 to 11:36 |
| The plan review | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,121,548 (USD 5.48) | 7 min 51 s, from 11:36 |
| The answer; the tests, red; the code | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | not reported | not reported | 11:44 to 11:50 |
| Round 1 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,057,444 (USD 4.81) | 8 min 52 s, from 11:51 |
| The fix of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 12:00 to 12:02 |
| Round 2 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,181,317 (USD 4.91) | 8 min 32 s, from 12:02 |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 12:24 to 12:26 |
