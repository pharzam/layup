# T-tmhw — the telemetry record: the schemas of `telemetry.tsv` and `prices.tsv` in code

Issue: [#94](https://github.com/pharzam/layup/issues/94), row 17 of the
[implementation plan](../plan/README.md); child of the core engine (#29).
Serves `F-0003#50`. Base `a2f9753` (the merge of row 16, #116). Author: Claude
Opus 5.5 on Claude Code. Evidence: [`runs/T-tmhw/`](../../runs/T-tmhw/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #94. The
reviewer order is the Operator's (Devin, then OpenCode, then Claude). Devin
(its usage quota) and OpenCode (Grok 4.7: no output in five minutes) gave no
record. The plan review (Claude Fable 5.1, effort `xhigh`, on the Claude Code
CLI, a fresh read-only session in a clone at `a2f9753`) gave
`approve-with-conditions`: Budget maximum 800 lines added plus removed over 16
files against the base, close-out inside; Cycle cap 1. Its one condition is
applied (the block test is red first, on a skeleton with no column), and its
notes are applied or confirmed: the two rows of `internal/records` in
`packages.md` read one way (2); `source` is an `http` or `https` URL, with no
`https` rule (3); the row rules are functions of a row, which the readers call
(4).

## What was done

1. **`internal/records`**, a new row of the table of phase 1 (K23): the Go
   schemas of the blocks `telemetry` and `prices`, the readers
   `ReadTelemetry` and `ReadPrices`, and the row rules `CheckTelemetry` and
   `CheckPrice`, which a writer of phase 2 calls before it writes a row.
2. **The row rules** that the blocks give in words (D2, D3): the form of a
   session ID; `latency_s` with `first_output`; `start ≤ first_output ≤ end`
   and the two durations recomputed in whole seconds; the tokens of each
   status and their reason; a cost and its currency with `money_status`; a
   price row for a computed cost; a subscription that is never `reported`; the
   form of a currency; an `http` or `https` source. An error names the line and
   the column.
3. **The tests:** the unit tests of each status and of each rule, and the
   block test; the names `telemetry` and `prices` leave the list
   `notYetBuilt` of `internal/tsv`.
4. **The documents:** `packages.md` (the row of `internal/records`, its line
   under the later phases, and K23 as "decided here"), `records.md` (the row
   rules as "decided here", with the known limits), the traceability rows that
   replace the planned handle, and the `PRD-0001` cell of REQ-011 with a §13
   row.

**The rejected alternatives:** the home `internal/ledger` (phase 2), and a
second package for the stall record, `internal/stall` (phase 3): "The packages
of later phases" gives `internal/records` the schemas of every record kind, and
one home serves rows 17 and 18; the rules inside the readers only (a writer
could not check a row before it writes it, note 4); a Go struct per row (the
other records of LAYUP are rows of strings that `internal/tsv` reads).

**Known limits:** the list of the codes of ISO 4217 is not checked, only their
form; that a `price` ID is a row of the host's `prices.tsv` is the writer's
check (phase 2); and the start row of a session records no spend cap and no
wall-clock limit, which architecture §12 names, as the block has no such
columns: the specification task of `M2b` gives them their columns (note 7 of
the plan review).

## Review round 1 and its fix (cycle 1)

Devin (its usage quota) and OpenCode (no output in five minutes) gave no
record. CI of PR #117 on `9beb171` passed its job `tests` first. Round 1
(Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI, a fresh read-only
session in a clone at `9beb171`, its record at 6 min 24 s; the record is on
#94) gave `material`, with two findings and eight notes. It ran 34 rows of its
own through a scratch copy of the package.

1. **A row with no start and no end passed.** `seconds` gave whether both
   values were times, and each caller dropped it, so `duration_s` 0 passed
   with no `start` and no `end`. Fixed: the result is used, and `seconds`
   reads the form of the type `time` (note 3: an offset passed a direct call of
   `CheckTelemetry`).
2. **A column whose rule forbids the empty value held it.** `docs/spec/README.md`
   gives the owner of a record the check of a column whose rule forbids `—`,
   and the rules named none. Fixed: in `telemetry.tsv` the ten columns with no
   clause for `—`, and in `prices.tsv` each column, never hold `—`, with a test
   row for each, as a "decided here" of `records.md`; `requirements` may be an
   empty list. The lesson is in `guardrails.md` §2, for row 18's stall record.

The notes: note 4 (the clause of a subscription read as a rule or as an
explanation) is applied as text, as an explanation; notes 3 and 5 to 10 need
no change (note 5: a source with a space or a scheme in capitals passes, a
lenient reading of "the URL"; note 6: the case of a session ID that is not
hexadecimal is the stronger one; note 7: the red runs are consistent; note 9:
a row of another length does not panic).

## Review round 2

Devin and OpenCode gave no record again. Round 2 (Claude Fable 5.1, a fresh
read-only session with another lens, the rule against the blocks and the edges
of each row rule, on `ed8ac6b`, its record at 5 min 5 s) gave `nothing
material in scope`. It checked each column of the two blocks against the code
and the "decided here", ran 44 rows of its own, and reproduced the 14 red cases
of the fix on the code of `9beb171`. Its notes 2 (`requirements` is `—` by its
type, not its rule) and 5 (the key of `prices.tsv`) are applied as text; note 3
(an error of the form of `end` names the column `start`, on a direct call
only) needs no change, as a writer's row passes the type first.

## Verdict

Delivered: the schemas of `telemetry.tsv` and `prices.tsv` in `internal/records`,
a new row of the table of phase 1 (K23), with the readers and the row rules
that the blocks give in words, and the rule that a column whose block rule or
type has no clause for `—` never holds it. The plan review (Claude Fable 5.1)
gave `approve-with-conditions`, with one condition, applied; round 1
(`9beb171`) gave `material`, two findings, fixed in cycle 1 (`ed8ac6b`); round
2 (`ed8ac6b`) gave `nothing material in scope`. The records are on #94;
`review-record-lint` passes (2 rounds, cap 1). At `ed8ac6b` the ladder passes
and CI passed its job `tests`. The close-out commit changes text only. The diff
against `origin/main` is 716 lines over 13 files, inside 800 over 16.

Next: row 18 (`T-dgy7`, #95), the stall record, the last ready row of phase 1
before row 19.

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC, 2026-10-03; the token count
is the `result` event of the Claude Code CLI; `not reported` otherwise.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan (D1 to D6) | reasoning | Claude Opus 5.5 | max | not reported | 02:29 to 02:31 |
| The plan review, first and second harness: skipped | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `xhigh` | not reported | 02:31:17 to 02:37:36 |
| The plan review | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 455,726 (USD 2.75) | 4 min 19 s, from 02:37:41 |
| The answer, the tests first, the code, the documents; the freeze, the ladder and the pull request | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 02:42 to 02:50 |
| Review round 1, first and second harness: skipped | reasoning | the same | `xhigh` | not reported | 02:50:55 to 02:57:15 |
| Review round 1 | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 792,046 (USD 3.37) | 6 min 40 s, 02:57:22 to 03:04:02 |
| The fix of round 1, test first; the freeze and the ladder | execution | Claude Opus 5.5 | max | not reported | 03:04 to 03:09 |
| Review round 2, first and second harness: skipped | reasoning | the same | `xhigh` | not reported | 03:09:21 to 03:15:39 |
| Review round 2 | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 447,936 (USD 2.97) | 5 min 16 s, 03:15:47 to 03:21:03 |
| The close-out | execution | Claude Opus 5.5 | max | not reported | 03:21 to 03:25 |
