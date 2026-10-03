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

## Review round 1

Devin (its usage quota) and OpenCode (no output in five minutes) gave no
record. CI of PR #118 on `6e3fa70` passed its job `tests` first. Round 1
(Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI, a fresh read-only
session in a clone at `6e3fa70`, its record at 5 min 36 s; the record is on
#95) gave `nothing material in scope`, with three notes. It checked each rule
of the "decided here" against the code in both directions, and ran edge cases
of its own on a scratch copy of the module: two rows of one stall at one time,
times across a day, `ST-1000` after `ST-999`, an `outcome` after an `outcome`,
a row after a closed stall, rows of another length, an empty file, a file with
only the header, CRLF lines.

The three notes need no change. Note 1 (the Fact and the ADR cells of the row
of the block test name those of `REQ-011` only) is the rule of the table: "A
row that covers more than one requirement gives the fact and the ADR of the
first". Note 2 (a row with more fields than columns passes a direct call of
`CheckStalls`) is outside its contract, "as tsv.Read gives them", as
`tsv.Read` refuses such a row. Note 3 says that the red runs of the evidence
agree with the code.

## Verdict

Delivered: the schema of `stalls.tsv` in `internal/records` (the home of row
17, K23), with the rules of a row (`CheckStall`) and of the order
(`CheckStalls`) that the block gives in words, an open stall valid in the
file, and the block lines of `kind`, `examiner` and `note` and the sentence of
the closed and the open stall amended to one reading. The plan review (Claude
Fable 5.1) gave `approve-with-conditions`, with two conditions, applied; round
1 (`6e3fa70`) gave `nothing material in scope`. The records are on #95;
`review-record-lint` reads them in CI. At `6e3fa70` the ladder passes (18
steps, each exit 0) and CI passed its job `tests`. The close-out commit
changes text only. The diff against `origin/main` is 500 lines over 12 files,
inside 800 over 16.

Next: row 19 (`T-efmy`, #96), the release review of phase 1, whose After cell
(rows 16, 17 and 18) is now met.

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC, 2026-10-03; the token count
is the `result` event of the Claude Code CLI; `not reported` otherwise.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan (D1 to D6) | reasoning | Claude Opus 5.5 | max | not reported | 03:24 to 03:25 |
| The plan review, first and second harness: skipped | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `xhigh` | not reported | 03:25:51 to 03:32:09 |
| The plan review | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 477,004 (USD 2.61) | 3 min 55 s, from 03:32:14 |
| The answer, the tests first, the code, the documents; the freeze, the ladder and the pull request | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 03:36 to 03:46 |
| Review round 1, first and second harness: skipped | reasoning | the same | `xhigh` | not reported | 03:46:33 to 03:52:53 |
| Review round 1 | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 1,123,470 (USD 2.86) | 5 min 55 s, 03:53:07 to 03:59:03 |
| The close-out | execution | Claude Opus 5.5 | max | not reported | 03:59 to 04:01 |
