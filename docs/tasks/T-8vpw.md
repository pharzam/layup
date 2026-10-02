# T-8vpw — the checks `markers`, `sources`, `discipline-tests` and `link-lint`

Issue: [#87](https://github.com/pharzam/layup/issues/87), row 10 of the
[implementation plan](../plan/README.md); child of `layup setup` (#77). Serves
`F-0003#42`. Base `2ea523b` (the merge of row 12, #110). Author: Claude Opus 5.5
on Claude Code. Evidence: [`runs/T-8vpw/`](../../runs/T-8vpw/). It fixes items
1 to 5 of #21.

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #87. The
reviewer order is the session's to name (Bootstrap mode rule 4); this session
follows the Operator's instruction given in the session on 2026-10-02 (Devin,
then OpenCode, then Claude). Devin (its usage quota) and OpenCode ("Go usage
limit exceeded") gave no record and are skipped. The plan review (Claude Fable
5.1, effort `xhigh`, on the Claude Code CLI, a fresh session) gave
`approve-with-conditions`: Budget maximum 2,400 lines added plus removed over 36
files against the base, close-out inside; Cycle cap 2 (the task changes
`setup-check.sh`). Its two conditions (the call of K10 takes a checked-out tree
and reads standard error; #21 closes only by the Operator's decision, asked on
#21) and its nine notes are applied.

## What was done

1. **#21:** `check_markers` lists the files with `git ls-files -z` and refuses
   a blank question or the empty mark, with the fixtures
   `markers/bad-quoted-name` (a DEL character in the name) and
   `markers/bad-blank-question` (items 1 and 5); of the stale sentences of
   item 2, `T-t8qp`, `T-9mmm` and `T-745n` rewrote all but one, and this task
   rewrites the last (`guardrails.md`, "In plain terms": a script fails on each
   gap with no open question, not on each gap; the plan said that later tasks
   rewrote all of them, and the check before the freeze found this one); the
   record of the setup says "summarised" (item 3); a lesson in
   `guardrails.md` §2 (item 4); the added scope of the comment on #21 is a
   known limit in `docs/setup/README.md`, and its decision waits on #21.
2. **`internal/verify`:** one scanner of markers for check `markers` and S10;
   check `markers` with the open gaps read by tabs; check `sources` for each
   value row of the record; the checks `discipline-tests` and `link-lint`, the
   baseline's own scripts run as a gate command; the call of K10 for S05.
3. **`internal/work`:** the schema of `open-gaps`, whose block moves to built.
4. **The tests:** the unit tests of the scanner, the open gaps, the sources,
   the scripts and the call; the harness on the group `markers` and the drift
   guard of `MK_EXEMPT`; the scripts on a tree and LAYUP's own scripts on a
   clone of its `HEAD`; the stand-in has stub scripts; an e2e case of a missing
   script.
5. **The documents:** the decided rules in `docs/spec/setup.md`, the lists that
   `internal/cli` hands to `internal/setup` in `packages.md`, the traceability
   rows, and the `PRD-0001` cells of `REQ-002`, `NFR-003` and `NFR-004` with a
   §13 row.

**The rejected alternatives:** the strict reader of the schema for
`open-gaps.tsv` (it refuses a row of two fields, which `check_markers` reads as
an empty question: `markers/bad-stale`); a fixture name with `"` (a Windows
checkout refuses it; DEL is quoted by git and accepted by Windows); a Go link
checker for S05 (the baseline's own `link-lint.sh` is the rule of check
`link-lint`); a hash of a `computed` ref by the check (the step that wrote the
hash checks it); closing #21 with its added scope undone (the Operator's
decision).

## Review round 1 and its fix (cycle 1)

Devin (its usage quota) and OpenCode ("Go usage limit exceeded") gave no
record. Round 1 (Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI, on
`44a0766`; the record is on #87) gave `material`, with one finding and ten
notes. The fix has its red run ([`test-runs.md`](../../runs/T-8vpw/test-runs.md)):

1. **A name with `\` in `check_markers`:** `awk -v` read the `\` as an escape,
   so the sh function gave `unlisted: docs/bx.md …` and
   `stale: docs/b\x.md …` for a listed marker of `docs/b\x.md`, where the
   engine passes. The `awk` now gets the name through `ENVIRON`. No fixture can
   hold the name (a Windows checkout refuses it), so
   `TestTheShAndTheGoFormOfMarkersAgree` makes a scratch repository at test time
   and runs both forms.

The notes: note 3 (macOS `awk` in a UTF-8 locale misses a marker at the end of
a line) is on the path of D3, the same lines from the two forms, and is fixed
in the same line: the `awk` of `check_markers` runs with `LC_ALL=C`, and the
same test runs the sh function in two UTF-8 locales. Notes 4 and 5 are fixed as
text in `setup.md` (the progress line names the check; an `sh` that does not
start is `fail`, reason `exit -1`, as for a gate command). Note 2 needs no
change: the field rule of `tsv.Write` makes a tab or a line feed of a cell a
space. Notes 6 to 11 confirm the change. One more gap, found by the author
during the fix: the planned row `T-8vpw/e2e/verify-not-active` of
`traceability.md` was still `planned`; it now names
`TestSetupVerifyWithNoBaselineScript`, `green`, and the second row of that test
is gone. A lesson in `guardrails.md` §2: "A path in `awk -v`".

## Review round 2

Devin (its usage quota) and OpenCode ("Go usage limit exceeded") gave no record
again. Round 2 (Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI, a
fresh session with another prompt, on `77dab2b`; the record is on #87) gave
`nothing material in scope`, with eight notes. It ran the fixed sh function on
nine names with `\`, `"`, a space, `%` or `é`, and on the cases of note 3 of
round 1, in three locales: the same lines as the engine; the script of
`44a0766` gave 12 wrong lines. Note 1 is applied as text at the close-out: a
marker in a file whose name holds a tab cannot be listed in `open-gaps.tsv` (a
cell of a TSV file holds no tab), a known limit of both forms in `setup.md`.
**Known limit** (note 3): the locale half of
`TestTheShAndTheGoFormOfMarkersAgree` finds a defect on a macOS host only; a
host with no UTF-8 locale, or an `awk` that reads bytes (as in CI), runs that
half as the C locale. The other notes need no change.

## Verdict

Delivered: the checks `markers`, `sources`, `discipline-tests` and `link-lint`
of `layup setup verify`: one scanner of markers for check `markers` and for
S10, the open gaps read by tabs, check `sources` for each value row of the
setup record, the baseline's two scripts run by the rule of a gate command, and
the call that gives S05 the files whose links break (K10 settled). Items 1 to 5
of #21 are fixed in `check_markers` and the documents; the added scope of the
comment on #21 (`check_ci`) waits for the Operator's decision there, so the
pull request says `Refs #21`.

The plan review (Claude Fable 5.1) gave `approve-with-conditions`, with two
conditions, applied; round 1 (`44a0766`) gave `material`, one finding, fixed
with notes 3 to 5 of that round; round 2 (`77dab2b`) gave `nothing material in
scope`. The records are on #87. At `77dab2b`, `go build`, `go vet` with each
tag, `gofmt`, the three test levels, `go test -race` on `internal/verify`,
`internal/work` and `internal/standin`, `run.sh` and the discipline tests pass;
at the head, all local checks pass, and `review-record-lint` passes on the
comments of #87 (2 rounds, cap 2). The diff against `origin/main` is
1,318 lines over 35 files with the close-out, inside the Budget maximum of
2,400 lines over 36 files.

Next: row 9 of the plan (`T-7s0y`, #86), the lowest row whose predecessors
have merged.

## Resource record

Recorded, not budgeted (ADR-0007). Times are 2026-10-02, UTC. A token count is
the `result` event of the Claude Code CLI (input, output, cache creation and
cache read tokens, and its cost) where that harness gave one; `not reported`
where the harness or the author's session gives none.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan, with the four inventory items, #21 and K10 | reasoning | Claude Opus 5.5 | max | not reported | 16:16 to 16:20 |
| The plan review, first and second harness: skipped (Devin's usage quota; "Go usage limit exceeded") | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `low` | not reported | 72 s, 16:20:11 to 16:21:23 |
| The plan review | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 645,998 (USD 3.09) | 9 min 11 s, 16:21:43 to 16:30:54 |
| The fix of items 1 and 5 of #21, test first, while the plan review ran | execution | Claude Opus 5.5 | max | not reported | 16:22 to 16:25 |
| The answer to the plan review, and the question on #21 | reasoning | Claude Opus 5.5 | max | not reported | 16:31 to 16:32 |
| The tests first, the package, the harness, the documents, the mutations and the evidence | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 16:32 to 16:48 |
| The freeze: the ladder, the last stale sentence of item 2 of #21, the ladder again | execution | Claude Opus 5.5 | max | not reported | 16:48 to 16:54 |
| Review round 1, first and second harness: skipped (Devin's usage quota; "Go usage limit exceeded") | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `low` | not reported | 73 s, 16:54:36 to 16:55:49 |
| Review round 1 | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 1,042,347 (USD 4.12) | 7 min 52 s, 16:55:55 to 17:03:47 |
| The fix of round 1 (finding 1, notes 3 to 5, the planned row), test first, and the freeze | execution | Claude Opus 5.5 | max | not reported | 17:04 to 17:14 |
| Review round 2, first and second harness: skipped (the same) | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `low` | not reported | 71 s, 17:14:13 to 17:15:24 |
| Review round 2 | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 952,709 (USD 3.00) | 6 min 42 s, 17:15:29 to 17:22:11 |
| The close-out, with note 1 of round 2 | execution | Claude Opus 5.5 | max | not reported | 17:22 to 17:24 |
