# T-dgy7 — the stall record: the schema of `stalls.tsv` in code

Issue: [#95](https://github.com/pharzam/layup/issues/95), row 18 of the
[implementation plan](../plan/README.md); child of the core engine (#29).
Serves `F-0003#49`. Base `d36340b` (the merge of row 17, #117). Author: Claude
Opus 5.5 on Claude Code. Evidence: [`runs/T-dgy7/`](../../runs/T-dgy7/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #95. The
reviewer order is the Operator's (Devin, then OpenCode, then Claude). Devin
(its usage quota) and OpenCode (Grok 4.7: no output in five minutes) gave no
record. The plan review (Claude Fable 5.1, effort `xhigh`, on the Claude Code
CLI, a fresh read-only session in a clone at `d36340b`) gave
`approve-with-conditions`: Budget maximum 800 lines added plus removed over 16
files against the base, close-out inside; Cycle cap 1. Its two conditions are
applied: the block line of `note` has its clause for `—` (1), and the sentence
"Each stall has three rows" now speaks of a closed stall, with the open stall
after it (2). Its notes are applied or named known limits: `examiner` is
required on a `diagnosis` row (1); the rows of one stall name one `task` (2);
the cross-column rules of §11 are the writer's (3); the one-row cases of each
per-kind column are in the tests (4); the scheme of a task ID is open (5); the
time order is per stall (6); the order check adds only what the key does not
hold (8); the backlog line moves at the close-out (9).

## What was done

1. **`internal/records`**, the home of row 17 (K23, D1): the Go schema of the
   block `stalls`, the reader `ReadStalls`, the rule of a row `CheckStall`,
   which a writer of phase 3 calls before it writes a row, and the rule of the
   order `CheckStalls`.
2. **The rules of a row** that the block gives in words (D2): `time` and
   `task` never hold `—`; each kind has its columns and `—` in the others;
   `examiner` is a session ID, required on a `diagnosis` row, and a session ID
   or `—` on a `diagnosis-failed` row; `task` is `project` exactly when the
   trigger is `orchestrator`. An error names the line and the column.
3. **The rules of the order** (D3): the stall IDs from `ST-001` with no gap;
   each row after the `stall` row of its ID, with its task; one second row of
   a stall; the `outcome` after it; the time order inside one stall. A stall
   with no second row or no `outcome` is open, and valid.
4. **The tests:** the unit tests of closed, open and orchestrator stalls and
   of each rule (30 cases, each with its line and its column), and the block
   test; the name `stalls` leaves the list `notYetBuilt` of `internal/tsv`,
   which is now empty.
5. **The documents:** `packages.md` (the Job cell of `internal/records`),
   `records.md` (the block lines of `kind`, `examiner` and `note`, the sentence
   of the closed and the open stall, D2 and D3 as "decided here", with the
   known limits), the traceability rows (the planned handle replaced, and the
   block test's row covers `REQ-009` too), and the `PRD-0001` cell of REQ-009
   with a §13 row.

**The rejected alternatives:** a second package for the stall record,
`internal/stall` (phase 3), as `packages.md` gives `internal/records` the
schemas of every record kind (row 17, K23); a stall that must have its three
rows (a writer appends a row when its step ends, so a file with an open stall
would be refused, and Stall Diagnosis could not count a missing second row);
the lenient `examiner` of the block on a `diagnosis` row (§11 step 2: an
examiner session writes it); a time order over the whole file (the block gives
no rule across stalls).

**Known limits:** the writer (phase 3) holds the rules of §11 that no block
sentence gives: the outcome after a `diagnosis-failed` row, and the outcome
and the `rung` of an orchestrator stall; `task` is `project` or the target's
own task ID, whose scheme the validator leaves open; that an `evidence` hash
names a file of `payloads/` is a check of two places, the writer's.

**Lessons:** none new for `guardrails.md` §2. The two order cases that the key
refused first (the evidence) are the lesson "A known-bad fixture that fails
for another reason", which the check of the line and the column held.
