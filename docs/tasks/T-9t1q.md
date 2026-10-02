# T-9t1q — the checks `facts`, `onboarding`, `glossary` and `guardrails` in a target's form

Issue: [#89](https://github.com/pharzam/layup/issues/89), row 12 of the
[implementation plan](../plan/README.md); child of `layup setup` (#77). Serves
`F-0003#42`. Base `bea865f` (the merge of row 11, #109). Author: Claude Opus 5.5
on Claude Code. Evidence: [`runs/T-9t1q/`](../../runs/T-9t1q/). It also fixes
#48 and #61.

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #89. The
reviewer order is the session's to name (Bootstrap mode rule 4); this session
follows the Operator's instruction given in the session on 2026-10-02 (Devin,
then OpenCode, then Claude), which no comment holds. Devin ("Unknown model",
with no model in its list) and OpenCode ("Go usage limit exceeded") gave no
record and are skipped. The plan review (Claude Fable 5.1, effort `xhigh`, on
the Claude Code CLI, a fresh session) gave `approve-with-conditions`: Budget
maximum 2,600 lines added plus removed over 40 files against the base,
close-out inside; Cycle cap 2 (the task changes `setup-check.sh`). Its three
conditions (the exit rule of the harness for the cases of LAYUP's form; what
check `facts` reads when its step's `done` row is not there yet; the test of
#48 at the discipline level) and its ten notes are applied.

## What was done

1. **#61:** `setup-check.sh` refuses a numbered fact of `F-0001` or `F-0003`
   whose text is a tab, with the fixture `facts/bad-blank-tab`.
2. **#48:** `.gitattributes` keeps the bytes of `docs/facts/**` (`-text`) and
   the line feeds of `docs/setup/setup-check.sh` and of the list
   `docs/setup/facts.sha256` (`text eol=lf`); the list was a third file that the
   test found. The mode `autocrlf` of `docs/setup/tests/run.sh` clones `HEAD`
   with `core.autocrlf=true`, and the case `facts/good-autocrlf` runs check
   `facts` on the clone.
3. **`internal/verify`:** the checks `facts`, `onboarding`, `glossary` and
   `guardrails` in a target's form, with the fact resolver; the core of each
   gives the lines of its sh function that a target keeps, and the harness
   compares those kinds of lines and names the cases of LAYUP's form.
4. **`internal/standin`:** the setup by hand also writes the answers record of
   S04 with its index row and its line in `facts.sha256`, the brief of S06 and
   the files of S07 to S09; the option `Done` adds the `done` rows of S03 to
   S14.
5. **The documents:** the target's form of the four checks in
   `docs/spec/setup.md` (the form of an answers record, K15; no marker in a
   record, K17; what check `facts` reads, K18; the briefs, the resolver, the
   entries of guardrails, the shared fixtures), `records.md` (S04, S11), the
   traceability row, the `PRD-0001` cells of `REQ-001`, `REQ-002` and `NFR-003`
   with a §13 row, and a lesson of #48 in `guardrails.md` §2.

**The rejected alternatives:** the `done` rows as the only switch of check
`facts` (the evidence call of a step comes before its `done` row: condition 2);
a Go test for #48 (the acceptance criterion names the discipline level, and
`run.sh` already runs in the job `setup-check`); a new check script for #48
(a new CI job and a list of checks in many documents); an answers record that
quotes a marker in a code span (check `markers` skips only the mention
`` `‹` ``); every line with `Check:` as an entry of guardrails (a `Check:` of
another paragraph would count).

## Review round 1

Devin (its usage quota) and OpenCode ("Go usage limit exceeded") gave no record.
A first run of Claude Fable 5.1 stopped on API errors (`api_retry`, with no
status) and gave no record in fifteen minutes; its process did not stop on the
signal of `timeout`, so the harness stopped it at about 15:58. Both are skipped
(Bootstrap mode rule 4). The order was tried again: Devin and OpenCode failed
again, and a second Fable run (effort `xhigh`, on the Claude Code CLI) reviewed
`7d03a87` and gave `nothing material in scope`, with six notes. It compared the
engine with the sh functions on a scratch repository of edge cases: the lines
are the same, except where `setup.md` names a difference. Notes 1 to 4 are
applied as text at the close-out, with no change of code: the reason of
`facts/good-autocrlf` in the list of the cases that the harness does not
compare; "the code span that opens the text after `N. `"; the path of a
`Check:` value is a file of the tree; the labels of the scratch commits in the
evidence. Notes 5 and 6 need no change.

## Verdict

Delivered: the checks `facts`, `onboarding`, `glossary` and `guardrails` of
`layup setup verify` in a target's form, with the fact resolver. The form of an
answers record is fixed for its writers (rows 9 and 13): K15 (the question ID,
the answer, `by`, `source` and the question text), K17 (no marker in a record)
and K18 (what exists is read; the `done` rows decide what must exist). The
fixes of #48 (`docs/facts/**` keeps its bytes; `setup-check.sh` and the list of
hashes keep their line feeds; a case of `run.sh` clones `HEAD` with
`core.autocrlf=true`) and #61 (a blank fact of one tab). **Known limit:** CI
restores `setup-check.sh` and `docs/setup/tests/` from the default branch, so
the fix of #61 and the case `facts/good-autocrlf` run in CI only after the
merge (note 6).

The plan review (Claude Fable 5.1) gave `approve-with-conditions`, with three
conditions, applied; round 1 (`7d03a87`, Claude Fable 5.1) gave `nothing
material in scope`. The records are on #89. At `7d03a87`, `go build`, `go vet`
with each tag, `gofmt`, the three test levels, `go test -race` on
`internal/verify` and `internal/standin`, `run.sh` and the discipline tests
pass; at the head, all local checks pass, and `review-record-lint` passes on
the comments of #89 (1 round, cap 2). The diff against `origin/main` is
1,130 lines over 24 files with the close-out, inside the Budget
maximum of 2,600 lines over 40 files.

Next: row 10 of the plan (`T-8vpw`, #87), the lowest row whose predecessors
have merged; row 9 waits for it.

## Resource record

Recorded, not budgeted (ADR-0007). Times are 2026-10-02, UTC. A token count is
the `result` event of the Claude Code CLI (input, output, cache creation and
cache read tokens, and its cost) where that harness gave one; `not reported`
where the harness or the author's session gives none.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan, with the five inventory items, the four sh functions and their fixtures | reasoning | Claude Opus 5.5 | max | not reported | 14:24 to 14:30 |
| The plan review, first and second harness: skipped ("Unknown model"; "Go usage limit exceeded") | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `low` | not reported | 7 s, 14:30:25 to 14:30:32 |
| The plan review | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 1,015,012 (USD 5.42) | 8 min 43 s, 14:30:55 to 14:39:38 |
| The fixes of #61 and #48, test first, while the plan review ran | execution | Claude Opus 5.5 | max | not reported | 14:31 to 14:34 |
| The answer to the plan review | reasoning | Claude Opus 5.5 | max | not reported | 14:39 to 14:41 |
| The test of #48 at the discipline level, the four checks, the stand-in, the harness, the specification and the records; the freeze | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 14:42 to 14:59 |
| Review round 1, first and second harness: skipped (Devin's usage quota; "Go usage limit exceeded") | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `low` | not reported | 12 s, 14:59:38 to 14:59:50 |
| Review round 1, a first Fable run: skipped (API errors, no record in fifteen minutes) | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | not reported | 15:00:19 to about 15:58 |
| Review round 1, the order tried again: Devin and OpenCode skipped | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `low` | not reported | 15:59:27 to 16:00:36 |
| Review round 1 | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 1,564,356 (USD 5.46) | 9 min 21 s, 16:00:56 to 16:10:17 |
| The close-out, with notes 1 to 4 of round 1 | execution | Claude Opus 5.5 | max | not reported | 16:10 to 16:11 |
