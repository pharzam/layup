# T-5zmw — `layup psb check` to its specification

Issue: [#83](https://github.com/pharzam/layup/issues/83), row 6 of the
[implementation plan](../plan/README.md); child of `layup setup` (#77). Serves
`F-0003#41`. Base `4b47982` (the merge of row 5, #105). Author: Claude Opus 5.5
on Claude Code. Evidence: [`runs/T-5zmw/`](../../runs/T-5zmw/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #83. The plan
review (GPT-6 Sol, effort `xhigh`, on the Devin CLI, a fresh session) gave
`approve-with-conditions`: Budget maximum 1,250 lines added plus removed over 28
files against the base, close-out inside; Cycle cap 1. Its four conditions and
five notes, and what the author did:

1. **The G1 case `**Technology stack:**` with no value** is the Operator's
   decision (the inventory names it as an external input). The author asked on
   #83; the Operator answered "My decision  is A" (**O-131**, 2026-10-02): such a
   line is a G1 gap, because the value of a named stack is a character that is
   neither white space nor `*`. The code of G1 changes in this task, and
   `TestG1ReadsTheValueOfAStack` pins it. The real problem statement has no such
   line, so `psb.tsv` and `F-0004` do not change.
2. **A table that the command cannot write** gives exit 2 and the error on
   standard error, decided in `docs/spec/psb-check.md` and named in the row of
   code 2 of `docs/spec/README.md`, with a failing output in a unit test of
   `WriteTSV`, a unit test of `internal/cli` and an end-to-end scenario.
3. **The record, not only its schema:** `TestEveryGoldenIsARecordOfTheBlock`
   reads each golden table by the block `psb-gaps` and checks the rules that
   code can check.
4. **An empty environment:** the end-to-end run on the real problem statement is
   repeated with no environment variable at all and no standard input.

The notes are applied: only G1 waited for the Operator; the evidence claims a
red run only for a changed behaviour; D7 changes only the sentence of
`psb-check.md`; the traceability rows name each moved or new test with its
level; `TestGoldenRealPSB` stays at the integration level.

## What was done

1. **Test first** ([`test-runs.md`](../../runs/T-5zmw/test-runs.md)): the red
   run of each changed behaviour before its code: the lone carriage return, G1,
   the error of the output (through `WriteTSV`, `internal/cli` and the binary of
   the base), the UTF-8 check, and the Go schema of `psb-gaps`; one measurement
   of a closed pipe and of a read-only output.
2. **`internal/psb`:** `GapsSchema`, compared with its block; `WriteTSV` writes
   through `tsv.Write`, so a tab, a line feed or a lone carriage return of a field
   is one space and the empty excerpt of line 0 is `—`, and it gives the error of
   its output (D1); G1 by O-131. The goldens are embedded, so the unit tests read
   no file (D5); the goldens `values` (D4) and `edge` (D3), with the cases that
   need a file of their own in `TestG1ReadsTheValueOfAStack` and
   `TestEdgeCases`, pin the present behaviour. `psb-gaps` moves to `built` in
   `internal/tsv/blocks_integration_test.go`.
3. **`internal/cli`:** a `FILE` that is not valid UTF-8 gives exit 2 and
   `layup: FILE: line <n> is not valid UTF-8`, with the line of the first byte
   that is not valid (D2, K32); a table that it cannot write gives exit 2. The
   read of `FILE` is a seam (`readFile`) that the unit tests replace, so they
   read no file; `TestPSBCheckExitCodes` moves to the integration level with real
   files, and gets a directory and a file that is not valid UTF-8 (D5).
4. **`cmd/layup`:** the end-to-end scenarios of `layup psb check` (D6): the real
   problem statement gives `psb.tsv` byte for byte with exit 1, twice, and the
   same with no environment variable; a statement with no gap gives the header
   only and exit 0; a missing file, a directory and a file that is not valid
   UTF-8 give exit 2 with no table and no usage; a read-only standard output
   gives exit 2. The harness gets `layupBare` and the shared `execute`.
5. **The specification:** in `docs/spec/psb-check.md` the two decided rules of
   the command, O-131 in the rule table, the edge cases of D3, the golden of
   D4, the schema and writer of D1, and S04 for the raw fact of the answers
   (D7, O-124); the row of code 2 and the K32 paragraph in `README.md`; the note
   of `packages.md` that `internal/psb` writes its table itself goes.
6. **The other documents:** the traceability rows (the planned row of this task
   becomes the rows of the real tests, and `TestPSBCheckExitCodes` is at the
   integration level); the `PRD-0001` §12 Test cells of `REQ-001`, `NFR-004`,
   `NFR-005` and `NFR-007`, with a §13 row.

**Two details inside the decisions of the plan:** the edge cases of D3 that
need a file of their own (a case of G1 or G3, which read the whole file) are
cases of `TestG1ReadsTheValueOfAStack` and `TestEdgeCases`, and the others are
rows of `edge.md` and `edge.tsv` (review round 1, finding 2); and the unit
tests of `internal/cli` read `FILE` through the seam `readFile`, so that the
failing output of condition 2 has a unit test there.

**The rejected alternatives:** a repair of bytes that are not valid UTF-8 (a
record would hold a value that the file does not hold); a closed pipe as the
end-to-end case of the write error (Go stops the program with `SIGPIPE`, so the
case gives no exit code); the write error as exit 1 (it is not a gap, and the
table is not complete).

## Review round 1 and its fixes (cycle 1)

Review round 1 (GPT-6 Sol, effort `xhigh`, on the Devin CLI; `edbeb60`, cycle
0) gave `material`, with three material findings and one note, all applied:

1. **G1 and white space.** A value of only U+00A0, U+2003 or a vertical tab
   named a stack, because Go's `\s` in the expression of O-131 has no such
   character, and the words of O-131 say "white space". The fix follows the
   words: white space is the Unicode property `White_Space`, and after the `:`
   the white space other than a tab, a form feed or a carriage return counts as
   a space, so a stack after U+00A0 is still named. A comparison on 465,010
   lines shows that a result of the base changes only where the base read `*`
   or a white space character as the value. The reading is in
   `psb-check.md`, as "decided here", with O-131; a lesson is in
   `docs/guardrails.md` §2.
2. **The edge golden.** The cases of D3 that a shared file can hold (the three
   cases of G2, and G4 inside a fence) are now rows of `edge.md` and
   `edge.tsv`; the cases of G1 and G3 stay in Go tests, each with a file of its
   own, and `psb-check.md` says where a later case goes.
3. **`Fast` in mixed case.** `values.md` holds `Fast`, so each word of G4 is in
   mixed case, as D4 says.

Note 4: the evidence no longer counts the record test of the goldens among the
tests that pass on the base.
