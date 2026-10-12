# T-bgvq — the guardrails lessons of the review rounds of T-bng8

Issue: [#288](https://github.com/pharzam/layup/issues/288), the step-7 lesson of
`T-bng8` (#267), which the frozen head of that task kept out of its branch: a
documentation fix that a task left owed. Serves `F-0003#51`, through #267, as `T-bng8`
does. Base `45f546c` (the merge of #285). Author: Claude Sonnet 5.5 on Claude Code.
Evidence: [`runs/T-bgvq/`](../../runs/T-bgvq/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #288. The plan
review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at
`45f546c`, 8 min 13 s) gave `approve-with-conditions`: Budget maximum 250 lines added
plus removed over 6 files, close-out inside; Cycle cap 1; no panel. Its three conditions
are applied: the plan names the fact it serves (this header); both entries go at the end
of the first list of section 2, not under "Gate pitfalls", as neither trap lets the gate
report success and the sibling entry is in that list; and the check of the first entry is
the local form that the workflow builds, since the check of the pull request runs on a
pull-request event and shows the gap at the next push. The notes are applied: the count of
1 is by the rule of a task of one artifact; the entries carry no time and no count of
pushes; the second names the rules as a harness script's, outside the repository (the
three of the issue and of round 7's record; the plan named two); `check.sh` has the
header of `parts.sh`.

## What was done

1. **T:** [`check.sh`](../../runs/T-bgvq/check.sh), red on `main`, with one line for each
   missing entry. Its first form printed five lines and was fixed before the record.
2. **D1:** the two entries, the last of the first list of section 2 of
   `docs/guardrails.md` (30 lines): "A review record without its heading" and "An allow
   rule whose glob reaches a script with a write mode". `check.sh` exits 0, and exits 1 on
   five mutants of the file; the linters exit 0.

## Review rounds

Round 1 (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at
`601cee6`, cycle 0, the cap of 1; its record is a comment on #288): `nothing material in
scope`, seven notes.

- **Applied** in `181e18e`, a commit that no round read (`check.sh` and the five mutants
  give the lines they gave, and the linters exit 0): note 3 (the second entry gave a cause
  for the widening of the rules that no source names), note 4 (the red check lasted until
  the heading was added) and note 7 (#288 is the issue of this task, #267 that of `T-bng8`).
- **Declined:** note 5 (`check.sh` names the first missing part of an entry only: a
  one-off evidence command, each mutant has one fault, and the exit is 1 either way) and
  note 6 (no glossary row for "the App": the glossary's rule is "no undefined
  abbreviation", and "App" is not one; two rows use the term the same way).
- **Carried by this file:** note 1 (the close-out is inside the room that the review
  counted) and note 2 (condition 1, in the header).

The `tests` job of CI failed on the first push, on an integration test that this change
does not touch; #291 records that flake and its three causes.

## Verdict

Delivered: the two lessons of the review rounds of `T-bng8` in section 2 of the
guardrails, under a check that was red before them and fails on five mutants of the
file. The review ended by decay at cycle 0 of cap 1. The diff against `origin/main` (`b460d3f`, after the merge of #289) is 249
lines over 5 files with the close-out, inside the Budget maximum of 250 lines over 6
files. No round read the merge of `origin/main` (clean, with none of the 3 files of the
branch), the notes commit or this close-out.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-12. The reviewers' tokens are
`modelUsage`, with Claude Code's estimate of the money; the author's are not reported.
The Operator chose Claude Sonnet 5.5 for the session; Model tiers names `Claude Sonnet 5`
for the execution tier, so the plan and its answer are reasoning parts on it.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Sonnet 5.5 | `max` | not reported | 12:02 to 12:07 |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 713,942 (USD 5.48) | 8 min 13 s, 12:08 to 12:16 |
| The answer to the plan review | reasoning | Claude Sonnet 5.5 | `max` | not reported | 12:18 |
| T, D1 and the mutants | execution | Claude Sonnet 5.5 | `max` | not reported | 12:22 to 12:26 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 822,124 (USD 4.83) | 10 min 2 s, 12:30 to 12:40 |
| The merge of `origin/main` and the notes | execution | Claude Sonnet 5.5 | `max` | not reported | 12:41 to 12:44 |
| #291, the flaky job that CI showed | — | Claude Sonnet 5.5 | `max` | not reported | 12:45 to 12:47 |
| The close-out | — | Claude Sonnet 5.5 | `max` | not reported | 12:47 to 12:47 |
