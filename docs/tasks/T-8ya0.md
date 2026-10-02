# T-8ya0 — the check `adapted`

Issue: [#88](https://github.com/pharzam/layup/issues/88), row 11 of the
[implementation plan](../plan/README.md); child of `layup setup` (#77). Serves
`F-0003#42`. Base `b0cb95a` (the merge of row 8, #108). Author: Claude Opus 5.5
on Claude Code. Evidence: [`runs/T-8ya0/`](../../runs/T-8ya0/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #88. The
reviewer order is the Operator's instruction of 2026-10-02: Devin, then
OpenCode, then Claude. Devin failed at once ("Unknown model", with no model in
its list, after its weekly quota ran out) and OpenCode stopped with "Go usage
limit exceeded"; both are skipped (Bootstrap mode rule 4). The plan review
(Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI, a fresh session)
gave `approve-with-conditions`: Budget maximum 1,500 lines added plus removed
over 20 files against the base, close-out inside; Cycle cap 1. Its one
condition (D3 states the removal of a space at the start and at the end of a
line) and its eight notes are applied.

## What was done

1. **`internal/verify`:** the core of `check_adapted` on the bytes of the
   tracked `.md` files (D2 to D5): the 16 patterns on every start position, a
   hit at the first letter of its word, the lower case of `A` to `Z` only, a
   marker as one unit, the findings once each in byte order; `AD_EXCLUDE` and
   `ad_allowed` as they are (O-123); `Flagged`, the list of the flagged files
   for the prose step of row 13 (K12); the check row `adapted`; the call
   `Files` of the git part of a check.
2. **The tests:** the unit tests of the patterns, the text and the files; the
   harness runs the group `adapted` and the frame case; the drift guard of the
   two lists; `Flagged` on fixture repositories; the rows of the built checks
   in the verify tests and the e2e test.
3. **The evidence:** red runs on a skeleton, on the base and on five
   mutations (two tests first passed for the wrong reason and were changed);
   and the comparisons with the sh function, among them the same 441 lines on
   the unchanged baseline of LAYUP's root commit.
4. **The documents:** the decided rules of check `adapted` and its table row in
   `docs/spec/setup.md`, with the known limit of the sh function in a UTF-8
   locale; the traceability rows; the `PRD-0001` §12 Test cell of `REQ-002`,
   with a §13 row; two lessons in `guardrails.md` §2 (a red run on a skeleton
   does not prove each rule; an `awk` check that reads the locale).

**The rejected alternatives:** a target's own lists in place of LAYUP's (O-123
keeps the rule as it is); a change of the sh function to `LC_ALL=C` (D8: no
change of a check script in this row; the fault is a known limit); the sh
function's word position (one byte after the match), which a boundary of more
than one byte moves; `strings.ToLower`, which changes the length of a text that
is not ASCII; a plain find-all of the matches, which loses a hit that shares a
boundary with the one before it.
