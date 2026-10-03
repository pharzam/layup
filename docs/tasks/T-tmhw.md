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
