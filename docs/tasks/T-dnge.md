# T-dnge — the acceptance row for the review of the release of M2d

Issue: [#322](https://github.com/pharzam/layup/issues/322), a defect that round 1 of `T-5xjc` (#283) revealed, off its path, so a task of its own (gate step 5). Serves `F-0003#43` through
REQ-015 and REQ-017. Base `b7f3f7d`. Author: Claude Sonnet 5.5 on Claude Code. Evidence: [`runs/T-dnge/`](../../runs/T-dnge/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #322. The plan review (GPT-6.1 Sol on Codex CLI, `high`, 2 min 23 s) gave `approve-with-conditions`: 100 lines over 6 files, cap 1, no
panel. Its two conditions are applied: the line of the coverage map `runs/T-bng8/parts.tsv` and its check, and the explicit green repeat of the count command (exactly one match). The check
failed on the base for two input-state rows that `T-gkvi` reworded: a defect off this path, #328, and the criterion became "no new line". The goal count is 1. The plan is reasoning-tier work that an
execution-tier model wrote under the Operator's routing message of 2026-10-13: a routing mismatch by choice.

## What was done

The table "The acceptance tests of M2d" (`docs/spec/scaffold.md`) has the row "The review of the release" (uat, `REQ-015` and `REQ-017`, `runs/T-efmy/release-check.sh` on the release commit, before the
demo), as the tables of `M2b` and `M2c` have it. The Source cell of row 59a of the plan names it, and `runs/T-bng8/parts.tsv` maps it to 59a. **Tests:** [`test-runs.md`](../../runs/T-dnge/test-runs.md):
the count command red then green, and the coverage check red a second time and back to the two lines of the base. **Known limit:** the coverage check still exits 1 on its two lines of #328.

## Review rounds

Round 1 (GPT-6.1 Sol on Codex CLI, `high`, `83cdfe5`, cycle 0 of the cap of 1; comment 6181821077): `nothing material in scope`, two raw findings. One note, applied: the opening of the test record and
the two diagnostics, kept in full. One finding of the class `material`, revealed and outside scope (the two lines of #328), which the reviewer accepted as classified.

## Verdict

Delivered: the acceptance table of `M2d` lists the review of its release, registered in the plan and the map. The review ended at `nothing material in scope` in one round.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-15, from the clock and the comments of #322; the tokens of GPT-6.1 Sol are the "tokens used" of Codex CLI; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning (a routing mismatch) | Claude Sonnet 5.5 | `max` | not reported | 00:41 |
| The plan review | reasoning | GPT-6.1 Sol, Codex CLI | `high` | 127,653 | 2 min 23 s, 00:41:55 |
| The answer, the tests, the row, the gate | execution | Claude Sonnet 5.5 | `max` | not reported | 00:44 to 00:48 |
| Round 1 | reasoning | GPT-6.1 Sol, Codex CLI | `high` | 81,919 | 2 min 56 s, 00:48:35 |
| The note and the close-out | execution | Claude Sonnet 5.5 | `max` | not reported | 00:52 |
