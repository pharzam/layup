# T-6rgp — a limit for each case of the end-class test

Issue: [#292](https://github.com/pharzam/layup/issues/292), cause 3 of the flaky `tests`
job of CI ([#291](https://github.com/pharzam/layup/issues/291)). Serves `F-0003#47`
(Verification on Every Change). `T-d6q5` (#92) met a failure of CI alone too (the entry "A
background maintenance of git" in `docs/guardrails.md`); that is a likeness of this flake,
not a precedent for the shape of this task. Base `1c167eb` (the merge of #286). Author:
Claude Sonnet 5.5 on Claude Code. Evidence: [`runs/T-6rgp/`](../../runs/T-6rgp/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #292. The plan review
(Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at `1c167eb`, 7 min
7 s) gave `approve-with-conditions`: Budget maximum 280 lines added plus removed over 6
files, close-out inside; Cycle cap 1; no panel. Its two conditions are applied: the busy
CPUs have numbers (the CPU count of the host, 2; one `sh` busy loop for each CPU, started
before a count and ended after it; one table of four counts of 100, in the evidence file);
and step M records the mutant that cannot fail (the case `wall` with 5 s) as a known limit
of the test, with no upper bound on its time. The notes are applied: a second failing line
in the red run would be the flake itself (the red run shows one line); a failure of
`TestProcessStopsTheWholeGroup` or `TestProcessEndsWhenAProgramLeftTheGroup` would open an
issue, and none failed; the default limit is explicit, as the next section says; and the 5 s
limit is the one that two tests of the package use, `TestProcessGivesTheInputOfEachMode` and
`TestProcessTimesTheFirstOutput`, not all of them.

## What was done

1. **T:** the case "crash by a signal of its own, late" (`sleep 0.5` before the
   signal), red on `1c167eb`: one line, class `wall` where `crash` is wanted, the process
   stopped by the limit at 0.303 s.
2. **D1:** the limit is chosen in the loop, by the class the case expects: 300 ms for the
   case that expects `wall`, 5 s for each other case. The plan wrote a field of the table
   with zero for 5 s; a field needs a fourth value in each of the eight literals, so the
   choice is made in the loop, with a comment that a `Wall` of zero stops a process at once
   (`process.go:202-207`). No zero reaches `Spec`. The test is green, and 100 runs, idle
   and beside a busy loop for each CPU, pass for the old test and for the new one.
3. **M:** two mutants on a copy of the tree. The shared 300 ms again fails with the line of
   the red run. The case `wall` with 5 s passes, in 5.5 s where the test takes 0.8 s: a
   known limit, recorded.
4. **Step 7:** the entry "A short time limit that a slow runner reaches", the last of the
   testing pitfalls of `docs/guardrails.md`; and `T-6rgp` in the Task cell of the row of
   the test in `docs/tests/traceability.md`. The plan review said that the row stays; its
   Task cell lists each task that changed the tests (`T-dq05, T-5zmw`), so it names this
   one.

**What the counts show.** The old test did not fail in 300 runs on this host, idle or
loaded: a shell starts in a few milliseconds here, and the flake of CI needs a start of
more than 300 ms. So 0 failures of the new test in 300 runs is no evidence of the fix. The
evidence is the late case, which fails on the shared limit on any host and passes on the
limit of its own.

## Review rounds

Round 1 (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at
`63b3bdd`, cycle 0, the cap of 1; its record is a comment on #292): `nothing material in
scope`, five notes. The `tests` job of CI passed on that head.

- **Applied** in `e4d4a1c`, a commit that no round read (documents only; the linters exit
  0): note 2 (the entry of the guardrails names its task before the check, as each
  neighbour does, and its lines run to 83 columns, within the 85 of the list).
- **Declined:** the last clause of the comment in the loop, which note 1 says explains a
  design the file never had ("a Wall of zero would stop it at once"). It warns against the
  other natural form, leaving the limit unset for the other cases: a `Wall` of zero stops a
  process at once. It stays.
- **Carried by this file:** note 1 (the form that landed is item 2 of "What was done": the
  limit is chosen in the loop by the class the case expects, in place of a field of the
  table with zero for 5 s), note 3 (the close-out is inside the room that the review
  counted), note 4 (the landing order: whichever of #292 and #293 lands second merges
  `origin/main` across the two files that both touched, `docs/guardrails.md` and
  `docs/tests/traceability.md`, and a round scoped to those files reads that merge) and note
  5 (what holds, checked: no action).

## Verdict

Delivered: `TestEachEndGetsItsClass` gives the short wall limit to the case that expects
`wall` and 5 s to each other case, and has a new case, "crash by a signal of its own,
late", that is red on the shared limit (class `wall`, the process stopped at 0.303 s) and
green on the limit of its own. The entry "A short time limit that a slow runner reaches" is
in section 2 of the guardrails. The old test did not fail in 300 runs here, idle or loaded,
so the counts do not show the flake: the late case is the evidence. The case `wall` with 5 s
still passes, a known limit that the record names. The review ended by decay at cycle 0 of
cap 1. The diff against `origin/main` (`1c167eb`, the merge of #286) is 230 lines over 6 files with the
close-out, inside the Budget maximum of 280 lines over 6 files. No round read the notes
commit or this close-out.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-12. The reviewers' tokens are `modelUsage`,
with Claude Code's estimate of the money; the author's are not reported. The Operator chose
Claude Sonnet 5.5 for the session; Model tiers names `Claude Sonnet 5` for the execution
tier, so the plan and its answer are reasoning parts on it. The counts of this task and of
`T-0p5p` ran one after another on one host, so the span of the third row holds both.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Sonnet 5.5 | `max` | not reported | 13:52 to 13:54 |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 883,632 (USD 4.94) | 7 min 7 s, 13:54 to 14:02 |
| The answer to the plan review | reasoning | Claude Sonnet 5.5 | `max` | not reported | 14:02 |
| T, D1, the counts and the mutants | execution | Claude Sonnet 5.5 | `max` | not reported | 14:09 to 15:36 |
| The record, the checks and the pull request | execution | Claude Sonnet 5.5 | `max` | not reported | 15:36 to 15:44 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 870,603 (USD 4.37) | 8 min 44 s, 15:45 to 15:54 |
| The notes | execution | Claude Sonnet 5.5 | `max` | not reported | 15:55 |
| The close-out | — | Claude Sonnet 5.5 | `max` | not reported | 15:56 to 15:57 |
