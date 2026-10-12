# T-pzgw — the prompt of the numbering, one clause per column

Issue: [#268](https://github.com/pharzam/layup/issues/268), a defect that the demo
of `M2c` revealed (row 49b, `T-ekh5`, #228), off its path (gate step 5). Serves
`F-0003#51` through REQ-012. Base `347d845` (the merge of #266). Author: Claude
Opus 5.5 on Claude Code. Evidence: [`runs/T-pzgw/`](../../runs/T-pzgw/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #268. The
plan review (GPT-6 Sol on Devin, a fresh session in a clone at `347d845`,
1 min 0 s; Claude Fable 5.1, first in the order, was at its usage limit) gave
`approve-with-conditions`: Budget maximum 160 lines added plus removed over 6
files, close-out inside; Cycle cap 1; no panel. Its three conditions are applied:

1. the directions of the numbers and the classes are kept;
2. the lesson is in the change under review;
3. no note on #228 is a step of this task.

## What was done

1. **The cause** ([`runs/T-pzgw/README.md`](../../runs/T-pzgw/README.md)):
   - the file of session `S-445dee66`, copied before its directory was removed:
     its 75 spans that are not facts each had `—` as their reason;
   - the sentence of `numberingPrompt` that reads two ways.
2. **T:** `TestTheNumberingPromptGivesEachColumnOnce`, red on `347d845` with
   both new sentences missing and the old one present.
3. **D1:** `numberingPrompt` gives a fact and a span that is not a fact one
   sentence each, with each column named once. The reason of a non-fact is a
   few words, never `—`.
4. **D2:** the other prompts of Intake were read for the same fault; none reads
   two ways (the README).
5. **The lesson:** `docs/guardrails.md` §2, "A prompt that lists the columns of
   a row in one sentence".

## Review rounds

Round 1 (GPT-6 Sol on Devin, `91743cf`, cycle 0, the cap of 1):
`not mergeable, findings recorded`.

- **The material finding:** the record of D2, which the plan had put in the task
  file, was not in the change under review.
- **The note:** the reviewer's integration run met a full `/tmp`, which the
  author's review clones had filled; they were removed.

The fix (`ed813d1`) added the record to the README. By **O-208** (b), the
Operator settled the finding as fixed by `ed813d1`, as O-207 did for
`T-msrd`, and no round reads that commit.

## Verdict

Delivered: the prompt of the numbering gives each column of each kind of span
once, under a test of its text, with the lesson in `guardrails.md`. The review
ended by the Operator's decision at the cap of 1 (O-208). The diff against
`347d845` is inside 160 lines over 6 files. Next: the demo of `M2c` runs again
on a binary built from this merge.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-11 unless named. Devin reports no
usage; the author's tokens are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The cause and the plan | reasoning | Claude Opus 5.5 | not reported | not reported | 21:44 (the issue) to 21:45 (the plan) |
| The plan review | reasoning | GPT-6 Sol | — | not reported by Devin | 1 min 0 s |
| T, D1, D2, the lesson | execution | Claude Opus 5.5 | not reported | not reported | 21:46 to 21:47 |
| Round 1 | reasoning | GPT-6 Sol | — | not reported by Devin | 4 min 37 s |
| The fix of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 21:56 to 21:57 |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 2026-10-12 05:53 to 05:53 |
