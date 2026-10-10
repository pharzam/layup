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
