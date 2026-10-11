# T-g49f — the technical specification of M2c, Intake

Issue: [#206](https://github.com/pharzam/layup/issues/206). Milestone `M2c` of the
[implementation plan](../plan/README.md#milestones), which starts with its
specification task (O-114), taken by the Operator's directive after `M2b`
(comment 6099589342 of #199). Serves `F-0003#51` (`REQ-012`) for `REQ-001`,
`REQ-012` and `NFR-003`. Base `b6b7a8a` (the merge of #204). Author: Claude
Opus 5.5 on Claude Code. Evidence: [`runs/T-g49f/`](../../runs/T-g49f/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #206. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `b6b7a8a`, 9 min 52 s) gave `approve-with-conditions`: Budget maximum
3,200 lines added plus removed over 28 files, close-out inside; Cycle cap 1; no
panel. Its five conditions are applied: the test step T before D1; decision 2
rewritten on `records.md` as it is (K38 settled `questions.tsv` and `budget.tsv`
already); the form of `layup spec check --facts` and the steps of Intake;
D3 completed (the cells of `internal/cli` and `internal/run` changed in place,
`internal/spec` and the engine checks); gap check step 4 and the two blocks of
step 5 as decisions 9 and 10.

## What was done

1. **S1, the public-solution search** ([`runs/T-g49f/search/`](../../runs/T-g49f/search/),
   committed before S3): five curated lists and 20 GitHub searches in the words of
   the facts that `M2c`'s requirements cite; two blind searchers, SWE-2 High on
   Devin (its free model, O-179, as its paid quota was exhausted) and Claude
   Sonnet 5.5 on Claude Code (OpenCode's login failed). Both conclude that no
   public solution can be the base of `M2c`; the summary gives what each part
   takes (the standoff span with its text, the smell and ambiguity taxonomies).
   Each searcher's first try failed on the author's permission setting, and its
   file says how.
2. **S2:** the survey of `T-w73g` has no section `M2c`, so no row is taken or
   rejected.
3. **S3, the decisions** (comment 6100128462 of #206): the second harness of the
   completeness review; `questions.tsv` and `budget.tsv` (K38); the budget of
   milestone 0; F-21; the sessions of a phase; the form of the batch; the pilot
   problem statement and the release review as rows of the slicing task; the
   setup questions and the markers deferred to `M2d`; the fields of the intent
   form; the command and the steps of Intake.
4. **T, the tests first:** [`sections.tsv`](../../runs/T-g49f/sections.tsv) and
   [`sections.sh`](../../runs/T-g49f/sections.sh), and the block test with the
   nine new names, red before the files and green after
   ([`test-runs.md`](../../runs/T-g49f/test-runs.md)).
5. **D1, [`docs/spec/intake.md`](../spec/intake.md)** (new): the sessions of
   Intake (a session of a phase, its `input/`), the budget of milestone 0, the
   rule gaps, the review of meaning and its quote check, the numbering, the span
   check and `layup spec check --facts SPANS --psb FILE`, the draft, the
   completeness review, the merge, the comment of the batch, the answers, the
   review of the answers and one follow-up, the raw fact of the answers, the
   steps of Intake, NFR-003 at Intake, the input states, "Not in M2c" and the
   acceptance tests.
6. **D2, [`docs/spec/records.md`](../spec/records.md#nfr-001--the-records-of-intake):**
   nine blocks (`meaning`, `spans`, `spec-check`, `draft`, `coverage`,
   `questions`, `answers`, `answers-review`, `budget`) in `notYetBuilt`, and the
   rows of the layout table.
7. **D3, [`docs/spec/packages.md`](../spec/packages.md#the-table-of-m2c):** the
   table of `M2c` (`internal/spec`), and the cells of `internal/cli` and
   `internal/run` changed in place.
8. **D4:** `README.md` (the file and the milestone), `psb-check.md` (Not in phase
   1), `setup.md` (the boundary rows of the gap check), `run.md` (Intake in the
   restart), `session.md` (the directory's `input/`, the budget), the plan (the
   `M2c` row, K38, F-21), the Task cells of `PRD-0001` §12 and a §13 line.

**The open questions of `later-p2-intake-spec`**, settled: the second harness
(decision 1, `intake.md`, The completeness review); the phase of
`questions.tsv` (K38, settled by `T-zck8`; its block here).

## Review rounds

The records are comments on #206; each round is Claude Fable 5.1, effort
`xhigh`, a fresh read-only session in a clone at its head. Round 1 (`cfb94e8`,
cycle 0, the cap of 1): seven material findings (the review of the answers
waited for every answer; the pair of the completeness review; the band as the
architecture gives it; no step read the follow-up's answers; the README; a
`git diff --check` that compared nothing, and five trailing spaces; the clock of
milestone 0). **O-193** (a): cap 2. Round 2 (`f3a321e`, cycle 1): the seven
fixed; one material finding that the fix opened (`budget.tsv` from the
follow-up). **O-194** (a): cap 3. Round 3 (`358f74e`, cycle 2): that fixed; two
material findings on the paths the fixes made explicit (a failed span check;
`B` not below `U`). **O-195** (a): cap 4. Round 4 (`3ff521a`, cycle 3, the cap):
`nothing material in scope`, six notes. Each round's verdict at its cap was
edited to `material` after the Operator raised it. The six notes are #210: the
close-out may change only the task indexes and this file, and the cap is
reached, so they land in the next change of `intake.md`, read by its round.

## Verdict

Delivered: the specification of `M2c`, Intake, in `docs/spec/intake.md`, with
the nine blocks of its records, the table of `M2c`, and the documents it
changed, after a public-solution search by two blind searchers. The review ended
by decay at cycle 3 of a cap that the Operator raised three times (O-193, O-194,
O-195). The diff against `b6b7a8a`, the branch's base, is inside 3,200 lines
over 28 files. Next: the slicing task of `M2c`, which also takes #210.

## Resource record

Recorded, not budgeted (ADR-0007). UTC; reviewers' and searchers' tokens are as
each harness reports them; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 2026-10-10 16:35 to 16:41 |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,247,562 (USD 6.76) | 9 min 52 s |
| S1, the sources | execution | Claude Opus 5.5 | not reported | not reported | 16:51 to 16:55 |
| Searcher A (two tries) | reasoning | SWE-2 High (execution tier, O-179) | — | not reported by Devin | 3 min 47 s; 12 min 13 s |
| Searcher B (two tries) | reasoning | Claude Sonnet 5.5 | — | USD 1.47; 2,562,338 (USD 2.08) | 3 min 9 s; 3 min 26 s |
| The summary, S3, T, D1 to D4 | execution | Claude Opus 5.5 | not reported | not reported | 17:12 to 17:23 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 1,231,180 (USD 6.28) | 12 min 17 s |
| The fix of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 17:37 to 19:14 (with O-193) |
| Round 2 | reasoning | Claude Fable 5.1 | `xhigh` | 1,303,804 (USD 6.13) | 11 min 54 s |
| The fix of round 2 | execution | Claude Opus 5.5 | not reported | not reported | 19:27 to 19:29 |
| Round 3 | reasoning | Claude Fable 5.1 | `xhigh` | 1,302,337 (USD 5.74) | 9 min 16 s |
| The fix of round 3 | execution | Claude Opus 5.5 | not reported | not reported | 19:41 to 19:54 (with O-195) |
| Round 4 | reasoning | Claude Fable 5.1 | `xhigh` | 1,784,378 (USD 6.62) | 12 min 16 s |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 20:09 to 20:08 |
