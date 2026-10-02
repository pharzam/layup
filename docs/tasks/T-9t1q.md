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
