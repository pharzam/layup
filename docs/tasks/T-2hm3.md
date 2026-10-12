# T-2hm3 — a bare repository that receives a push asks for no maintenance

Issue: [#297](https://github.com/pharzam/layup/issues/297), cause 2 of the flaky `tests`
job of CI ([#291](https://github.com/pharzam/layup/issues/291)). Serves `F-0003#47`
(Verification on Every Change). `T-d6q5` (#92) started each git of the product and of a
helper that commits with `-c maintenance.auto=false` and left the receiving side open (the
entry "A background maintenance of git", `docs/guardrails.md`). Base `a8b3180` (the merge
of #295). Author: Claude Sonnet 5.5 on Claude Code. Evidence:
[`runs/T-2hm3/`](../../runs/T-2hm3/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #297. The plan review
(Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at `a8b3180`, 10
min 18 s) gave `approve-with-conditions`: Budget maximum 480 lines added plus removed over
11 files, close-out inside; Cycle cap 1; no panel. Its three conditions are applied: the
new test pushes with `gitOut`, because a call of `internal/git` drops `GIT_TRACE2_EVENT`,
and its positive control is the `cmd_name` event of `receive-pack`; the issue has a
criterion for the three other fixtures; the setup is `runs/T-2hm3/repro.go`, a Go program
with `//go:build ignore`, no package of the module. The notes are applied or carried: the
entry says that the receiving repository of a local push is not one of those that
`docs/spec/packages.md` speaks of; the Task cell of six rows and one new row of the
traceability table; each key alone is a mutant that passes; each count has its own
timeout. The one line of `gitOut` of `internal/run` (no `-c maintenance.auto=false`) is
not in this branch: it is issue #298.

## What was done

1. **The cause,** before the issue: the logs of the three runs, the tests, the trace of a
   push, the made setup, and the source of git 2.54.0, 2.55.0 and 2.56.0 (outside the
   repository). The plan named two gits; the Operator installed 2.56.0 at 16:57 and it is
   on `PATH`, so three ran.
2. **T:** `TestAPushIntoTheTargetOfAWorldStartsNoMaintenance`, red with each git: the
   trace of the push holds `"maintenance","run"` and a `cmd_name` of `maintenance`.
3. **D1:** `receive.autogc` and `maintenance.auto` are `false` in the bare repository of
   `newWorld`, `newRunWorld`, `pushed` and `TestPushAndFetchOfABareRepository`. The test is
   green with each git, and the three suites pass with each.
4. **M:** the two settings removed from `newWorld`: the test fails as in the red run; each
   key alone passes, as the plan review said, so the entry says to set both.
5. **Step 7:** the entry "A bare repository that receives a push starts its own
   maintenance" in `docs/guardrails.md`; in `docs/tests/traceability.md` a row for the new
   test and `T-2hm3` in six rows.

**What the work found that the plan had not.** The geometric strategy became the default of
unscheduled maintenance in git 2.54, not before (2.52 and 2.53 have it but keep `gc`), and
git 2.56.0 starts no repack in the made setup: the condition moved into the object database
layer. So the made setup fails on 2.55.0 (16 of 100) and not on 2.47.3 or 2.56.0 (0 of 100
each); the fix is the same for every version, as each key stops the receiving git.

**The counts.** The three tests failed 3 times in 1,820 runs in CI, about 1 in 600. Here
the old fixture failed 0 of 600 and the new one 0 of 600, so the counts show only that the
change breaks nothing; the evidence is the test of the trace and the made setup.

## Review rounds

Round 1 (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at
`c3387d1`, cycle 0, the cap of 1; its record is a comment on #297): `nothing material in
scope`, five notes. The `tests` job of CI passed on that head.

- **Applied** in `c03d155`, a commit that no round read (documents, and the comment of the
  new test): note 1 (the record has a section "Where the readings come from": the source
  files by tag and the three CI runs), note 2 (the entry names the whole argv, `--auto
  --quiet --detach`) and note 3 (the entry and the comment say that git 2.54 is read from
  its source and 2.55.0 measured, and that 2.47.3 and 2.56.0 failed none of 100 removals).
- **Carried by this file:** note 4 (the close-out is inside the room that the review
  counted) and note 5 (#298, opened for the plan review's note 5).

## Verdict

Delivered: the four bare repositories that tests push into set `receive.autogc` and
`maintenance.auto` to `false`, so the git that receives a push starts no maintenance in
them, and `TestAPushIntoTheTargetOfAWorldStartsNoMaintenance` is red with git 2.47.3, 2.55.0
and 2.56.0 and green after. The made setup, which fails 16 of 100 removals on git 2.55.0,
fails none of 100 with the settings. The three tests of the failures failed 0 of 600 before
and 0 of 600 after, so the counts show only that the change breaks nothing, and the three
suites pass with each git. The review ended by decay at cycle 0 of cap 1. The diff against
`origin/main` (`a8b3180`, after the merge of #295) is 466 lines over 10 files with the close-out, inside the
Budget maximum of 480 lines over 11 files. No round read the notes commit or this
close-out.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-12. The reviewers' tokens are `modelUsage`,
with Claude Code's estimate of the money; the author's are not reported. The Operator chose
Claude Sonnet 5.5 for the session; Model tiers names `Claude Sonnet 5` for the execution
tier, so the plan and its answer are reasoning parts on it.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The cause: the logs, the traces, the made setup | reasoning | Claude Sonnet 5.5 | `max` | not reported | 16:24 to 16:46 |
| The plan | reasoning | Claude Sonnet 5.5 | `max` | not reported | 16:46 to 16:48 |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 860,747 (USD 5.97) | 10 min 18 s, 16:48 to 16:58 |
| The answer to the plan review | reasoning | Claude Sonnet 5.5 | `max` | not reported | 16:59 |
| T, D1, the mutants, the counts and the made setup | execution | Claude Sonnet 5.5 | `max` | not reported | 17:16 to 18:12 |
| The suites, the record, the checks and the pull request | execution | Claude Sonnet 5.5 | `max` | not reported | 18:12 to 18:29 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 902,182 (USD 5.32) | 11 min 8 s, 18:29 to 18:40 |
| The notes | execution | Claude Sonnet 5.5 | `max` | not reported | 18:40 to 18:42 |
| The close-out | — | Claude Sonnet 5.5 | `max` | not reported | 18:42 to 18:43 |
