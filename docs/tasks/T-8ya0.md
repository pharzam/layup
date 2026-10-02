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

## Review round 1

Devin and OpenCode failed at once (Devin's usage quota; "Go usage limit
exceeded") and are skipped. Claude Fable 5.1 on the Claude Code CLI (effort
`xhigh`) reviewed `c6118e0` and gave `nothing material in scope`, with five
notes. It compared the core with the sh function on 46 tracked paths of edge
cases: the lines are the same, except the NUL bytes (note 1) and the paths that
git quotes (D2). Notes 1 and 2 are applied as text in `setup.md` at the
close-out, with no change of code: a NUL byte is one more byte for the engine,
where macOS `awk` ends the line; the exemption is of the first pattern of rule 1
only, as the code has it. Note 3 (`TestFlagged` has no case `bad-rule-3`) is not
applied: the harness runs the core on that case and `Flagged` reads the same
files, so the reviewed code stays as it is. Notes 4 and 5 need no change.

## Verdict

Delivered: check `adapted` of `layup setup verify`. `internal/verify` ports
`check_adapted` on the bytes of the tracked `.md` files, with LAYUP's
`AD_EXCLUDE` and `ad_allowed` as they are (O-123): the 16 patterns on every
start position, a hit at the first letter of its word, the lower case of `A` to
`Z` only, a marker as one unit, the findings once each in byte order.
`Flagged` gives the list of the flagged files for the prose step of row 13
(K12). On LAYUP's unchanged baseline the engine and the sh function give the
same 441 lines; the sh function's fault in a UTF-8 locale of macOS `awk` is a
known limit in `setup.md`.

The plan review (Claude Fable 5.1) gave `approve-with-conditions`, with one
condition, applied; round 1 (`c6118e0`, Claude Fable 5.1) gave `nothing
material in scope`. The records are on #88. At `c6118e0`, `go build`, `go vet`
with each tag, `gofmt`, the three test levels and `go test -race` on
`internal/verify` pass; at the head, all local checks pass, and
`review-record-lint` passes on the comments of #88 (1 round, cap 1). The diff
against `origin/main` is 712 lines over 17 files with the
close-out, inside the Budget maximum of 1,500 lines over 20 files.

Next: row 12 of the plan (`T-9t1q`, #89), the lowest row whose predecessors
have merged; rows 9 and 10 wait for it.

## Resource record

Recorded, not budgeted (ADR-0007). Times are 2026-10-02, UTC. A token count is
the `result` event of the Claude Code CLI (input, output, cache creation and
cache read tokens, and its cost) where that harness gave one; `not reported`
where the harness or the author's session gives none.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan, with the inventory item and the shell function | reasoning | Claude Opus 5.5 | max | not reported | 13:43 to 13:48 |
| The plan review, first harness: skipped ("Unknown model", with no model in its list) | reasoning | GPT-6 Sol on the Devin CLI | `xhigh` | not reported | 5 s, 13:48:26 to 13:48:31 |
| The plan review, second harness: skipped ("Go usage limit exceeded") | reasoning | Grok 4.7 on the OpenCode CLI | `xhigh` | not reported | 6 min 11 s, 13:48:38 to 13:54:49 |
| The plan review | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 943,608 (USD 4.38) | 8 min 1 s, 13:55:24 to 14:03:25 |
| The answer to the plan review | reasoning | Claude Opus 5.5 | max | not reported | 14:03 to 14:04 |
| The core, the tests, the comparisons with the sh function, the specification and the records; the freeze | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 13:50 to 14:11; the core drafted 13:50 to 13:55, while the plan review ran |
| Review round 1, first and second harness: skipped (Devin's usage quota; "Go usage limit exceeded") | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `low` | not reported | 7 s, 14:11:36 to 14:11:43 |
| Review round 1 | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 796,558 (USD 3.54) | 7 min 19 s, 14:12:07 to 14:19:26 |
| The close-out, with notes 1 and 2 of round 1 | execution | Claude Opus 5.5 | max | not reported | 14:20 to 14:20 |
