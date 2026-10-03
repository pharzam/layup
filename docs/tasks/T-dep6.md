# T-dep6 — a whole setup, end to end, with no network

Issue: [#93](https://github.com/pharzam/layup/issues/93), row 16 of the
[implementation plan](../plan/README.md); child of `layup setup` (#77). Serves
`F-0003#42`. Base `9a16ff8` (the merge of row 15, #115). Author: Claude Opus
5.5 on Claude Code. Evidence: [`runs/T-dep6/`](../../runs/T-dep6/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #93. The
reviewer order is the Operator's (Devin, then OpenCode, then Claude). Devin
(its usage quota) and OpenCode (Grok 4.7: no output in five minutes) gave no
record. The plan review (Claude Fable 5.1, effort `xhigh`, on the Claude Code
CLI, a fresh read-only session in a clone at `9a16ff8`) gave
`approve-with-conditions`: Budget maximum 2,600 lines added plus removed over
24 files against the base, close-out inside; Cycle cap 1 (the branch changes
no CI file, no hook and no check script). Its one condition is applied, with
one change of its word rule (below), and its ten notes are applied or
confirmed. The demo is the restated sentence of the issue: on a stand-in
baseline with no network, a whole setup ends with `layup setup verify` at exit
code 0; the records in Git and the gate with LAYUP absent are acceptance
criteria of the same run.

**The condition** (the rules of the test of `NFR-002` item 2): no tracked file
named `layup` or `setup-check.sh`, or with the magic bytes of a program; no
word `layup` or `setup-check` in a line under `.github/` with its comment
removed; no `layup/` context in the two JSON files of S13. **The change of the
word rule:** a word is delimited by a character that is not a letter, a digit,
`-` or `_`. The condition also kept `.` and `/` inside a word, so `./layup gate`
and a URL of LAYUP's repository would be one longer word and pass; with this
rule both are refused, and `layup-records` is not.

## What was done

1. **One shared whole-setup run** (D1, K27): the first scenario makes it, in
   the directory of `TestMain`: the stand-in baseline with a second commit of
   the cases of a whole run (LAYUP's own `link-lint.sh`, a note of the history,
   a guide whose link S05 breaks, three markers, a flagged file, the
   baseline's own workflow), a problem statement with one gap, each run of
   the binary with the answers and the input files of each stop, `layup setup
   verify` into `out/verify.tsv`, and the last run (S15). It keeps a snapshot
   at the stop `O-verify`. A scenario copies a work area before it changes it.
2. **The scenarios** of `cmd/layup/whole_e2e_test.go`: the stops and the end,
   each stop the same bytes on two runs, and the pin resolved once (D2, D3,
   `NFR-005`, `NFR-006`); `layup setup verify` at exit 0 with the exact table,
   four seeded defects that each fail their own row with exit 1, and a missing
   baseline script that is `not-active`, whose table S15 refuses (D4,
   `NFR-004`); the records in the target's Git (D5, `NFR-001`); the pushes of
   `commands.sh` to a bare repository with a stub of `gh` (D7); and the target
   with LAYUP absent (D6, `NFR-002`).
3. **Test first:** 11 mutations of the code that the scenarios guard, each
   detected by its scenario's own assertion.
4. **The documents:** `setup.md` (K27, K28 and the pushes, linked to
   `test-levels.md` and to Bootstrap mode, with no copied value), `records.md`
   (the test of `NFR-001` in phase 1, and the mechanical test of `NFR-002`
   item 2), the traceability rows that replace the three planned handles, the
   `PRD-0001` cells with a §13 row, and two lessons in `guardrails.md` §2.

No code of the product changes, and no CI file (K28).

**The rejected alternatives:** a run per scenario (five whole setups, and each
verify runs the gate twice per active kind); a fixture in a `t.TempDir()` of
the first scenario (Go removes it when that scenario ends, note 5); a rewrite of
the `remote add` line of `commands.sh` (the test would not run the file as the
Operator runs it; `url.<bare>.insteadOf` in the environment of the test keeps
the file and the target's own `origin`); a YAML reader for the run lines (the
test reads the run line `sh .github/gates.sh <kind>` of each kind as text, as
the catalog writes it).

**Known limits:** the scenarios run on the stand-in baseline only, not on the
real baseline at LAYUP's pin, whose root commit a shallow checkout of CI does
not hold (row 20, the first pilot, uses it); the apply of the ruleset reaches a
stub of `gh`, not GitHub (row 20); the UAT parts of `verify-acceptance` and the
`NFR-003` audit are row 20's and row 19's.
