# T-mo9w — row 48 of the plan, the e2e of the steps of Intake

Issue: [#226](https://github.com/pharzam/layup/issues/226), row 48 of [the tasks
of M2c](../plan/README.md#the-tasks-of-m2c), opened by `T-2r8u` (#211). Serves
`F-0003#41` and `F-0003#51` through REQ-001, REQ-012 and NFR-001. Base `121d6a6`
(the merge of #259, row 47c); the budget is read against it, the base of this
branch's own work. Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-mo9w/`](../../runs/T-mo9w/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #226. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `121d6a6`, 8 min 20 s) gave `approve-with-conditions`: Budget maximum
700 lines added plus removed over 11 files, close-out inside; Cycle cap 1; no
panel. Its condition is applied: each scripted comment takes its time from the
`asked` of the record plus a minute. The goal count is 1, final by O-196.

## What was done

1. **`TestTheStepsOfIntake`** (`cmd/layup`, e2e): the built binary, the fake
   forge and the world's scripted harnesses, no secret: a restart runs Intake
   to the batch (four questions) and waits; after a comment of the idea owner
   that leaves Q-001 unanswered, the next restart reviews the answers, posts one
   follow-up for Q-001 and waits; after the answer to the follow-up, the next
   restart records it and ends Intake ("budget.tsv kept; Intake is done").
2. **The e2e world:** its fake forge serves a person's comment on the Intake
   issue beside the App's; its harness script answers the review prompt.
3. **The documents:** `traceability.md`, the Task and Test cells of REQ-001,
   REQ-012 and NFR-001 and a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-mo9w/test-runs.md): red before the
world was extended; 2 of 2 mutations of the product caught; green at the three
levels. No defect of the product was found on the path.
