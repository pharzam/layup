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
   `markers/bad-blank-question`; the record of the setup says "summarised";
   a lesson in `guardrails.md` §2; the added scope of the comment on #21 is a
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
