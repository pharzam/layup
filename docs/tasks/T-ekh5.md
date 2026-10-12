# T-ekh5 — row 49b of the plan, the demo of M2c

Issue: [#228](https://github.com/pharzam/layup/issues/228), row 49b of [the tasks
of M2c](../plan/README.md#the-tasks-of-m2c), opened by `T-2r8u` (#211). Serves
`F-0003#41` and `F-0003#51` through REQ-001 and REQ-012. Base `4ccbb12` (the merge of
#263, row 49a); the budget is read against `origin/main` at landing. Authors: Claude
Opus 5.5 on Claude Code (the plan, the registers, Start, runs 1 to 5) and Claude
Sonnet 5.5 on Claude Code (runs 6 and 7, the evidence, the documents, the close-out).
Evidence: [`runs/T-ekh5/`](../../runs/T-ekh5/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #228. The plan
review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at
`4ccbb12`, 10 min 7 s) gave `approve-with-conditions`: Budget maximum 2400 lines added
plus removed over 24 files, close-out inside; Cycle cap 1; no panel. Its five
conditions are applied: `--lease-h 10`, named before Start with the wait of
`3 × lease.H` after a `fail`; the release commit named `0b881fa`, with the diff of the
binary's Go code against it recorded; one reading of "holds the rule gaps and the gaps
of meaning" in the rules; the uat scenario and its sign-off; the body of the batch from
the forge, with `questions.tsv` as its record. The goal count is 1, final by O-196.

## What was done

1. **The rules and the scenario** (T), written before any run (`3f6e910`).
2. **The registers and the binary** (D1, D2): `claude-opus-5-5` in the models, the routes
   of three roles, the cap of `claude` at 5.00 (O-201) and its wall at 30 minutes (O-203);
   `layup 0.1.0-dev`, built on this branch.
3. **Start** (D3), by the Operator, on 2026-10-11 at 20:33Z, with `--lease-h 10`.
4. **The restarts** (D4): seven, runs 1 to 7, after Start, run 0 ([the
   README](../../runs/T-ekh5/README.md#the-runs)). Runs 1 to 3 were refused at `numbering`
   by a sentence of the prompt that read two ways (#268, fixed by `T-pzgw`); run 4 ended
   `crash` (the weekly limit of `claude`); run 5 posted the batch of 34 questions; runs 6
   and 7 (D6, not a criterion) ran the follow-up and ended Intake.
5. **The evidence** (D5): the batch, `questions.tsv`, the span table and its check, the
   actors and the telemetry; the four rules hold. The finding of the demo is #300.
6. **The documents:** the planned uat row of `traceability.md` in its place, the cells of
   REQ-001 and REQ-012 and a §13 line of `PRD-0001`, a pitfall in `guardrails.md` §2.
7. **The sign-off:** The Operator accepted the scenario on 2026-10-12 at 19:20Z, by
   Telegram, on reading `runs/T-ekh5/` at `52b3311`; the comment 6131023042 of #228
   records it, and the README takes its `Accepted:` line from it.

**Tests:** no Go file changed; the README of `runs/T-ekh5/` is the record of the uat.

## Review rounds

Round 1 (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at
`6a88c70`, cycle 0, the cap of 1; its record is comment 6131238615 on #228): `nothing
material in scope`, five notes. The `tests` job of CI passed on that head.

- **Applied** in `8d9fc59`, a commit that no round read (the README and one entry of
  `guardrails.md`): note 1 (the README says that 30 of the 34 questions carry a quote and
  why Q-001 and the three questions of the intent form do not), note 2 (the wait of run 5
  is about 30 minutes, by the times that the directory holds), note 3 (the check of the
  pitfall says that the edit repairs a quote, not a step) and note 5 (the release check at
  the code of the demo, exit 0 with the output of row 49a, and a comment on #268, whose
  task ran none).
- **Carried by this file:** note 4 (the close-out is inside the room that the review
  counted: 24 files).

## Verdict

Delivered: the demo of `M2c` on a pilot problem statement. Restarts of `layup run` on
`pharzam/layup-uat-m2c`, with the brief PSB-CHAT-001 and two real harnesses, posted one
comment of the App with 34 questions (1 rule, 21 meaning, 3 intent, 9 coverage); `layup
spec check --facts` exits 0 on its 800 numbered spans; each write of LAYUP is by the App's
bot; and Intake ended after one follow-up. The four rules written before any run hold, and
the Operator accepted the scenario on 2026-10-12 at 19:20Z (comment 6131023042 of #228).
Not shown: the money of the cap of milestone 0 (the billing is a subscription, so it
counts none) and the clause "each push of a task branch" (Intake pushes none). The review
ended by decay at cycle 0 of cap 1. The diff against `7378f6d` is inside 2,400 lines over
24 files (1,551 lines added plus removed, measured against `origin/main` before the
commit). The demo revealed #300, left open; #268, fixed by `T-pzgw`, had caused runs 1 to
3 to fail. Next: row 50 of `M2d`, `T-tetz` (#270).

## Resource record

Recorded, not budgeted (ADR-0007). UTC; reviewers' tokens are `modelUsage`; the authors'
are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 2026-10-11 15:20 to 15:23 (the claim to the post) |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,543,620 (USD 7.78) | 10 min 7 s |
| The registers, the binary and the rules | execution | Claude Opus 5.5 | not reported | not reported | 2026-10-11 15:34 to 15:54 |
| Start, runs 1 to 5 | execution | Claude Opus 5.5 | not reported | not reported | the Operator's Start 20:26 to 20:34; the runs in the README |
| The sessions of Intake in the runs (the demo's subject) | reasoning | Claude Opus 5.5 (8 sessions), Devin `swe-2-high` (1 session, 1 probe), Claude Sonnet 5.5 (1 probe) | not reported | in 162, out 106,495, cache 4,347,006 (observed by Claude Code; Devin none; money unknown) | 22 min 25 s in all |
| Runs 6 and 7, the evidence and the documents | execution | Claude Sonnet 5.5 | `max` | not reported | 2026-10-12 18:07 to 19:15 (runs 6 and 7 to 18:44; the pull request opened at 19:15) |
| The sign-off | — | the Operator | — | — | 19:15 to 19:20 (asked at 19:15, answered at 19:20:12) |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 1,225,763 (USD 5.36) | 9 min 43 s, 19:21 to 19:31 |
| The notes | execution | Claude Sonnet 5.5 | `max` | not reported | 19:33 to 19:35 |
| The close-out | execution | Claude Sonnet 5.5 | `max` | not reported | 19:36, to the commit of this file |
