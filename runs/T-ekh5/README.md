# T-ekh5 — the demo of M2c (uat)

Row 49b of [the tasks of M2c](../../docs/plan/README.md#the-tasks-of-m2c), issue
[#228](https://github.com/pharzam/layup/issues/228). The inputs are the
Operator's O-201 (comment 6110518260 of #228). The plan, its review and the
author's answer are comments on #228.

## The rules, written before any run

Written on 2026-10-11, before the registers changed and before Start (step T of
the plan, with conditions 3 and 4 of its review). The demo passes when each of
the four conditions holds; it fails when one is false, or when a step of the
restart is `fail` after the batch is posted.

1. The App posts the comment of the batch on the Intake issue, and its first line
   is `LAYUP Intake: the batch of questions`. It holds at least one question of
   the source `rule` and at least one of the source `meaning` (`questions.tsv`
   of the records). A source with no row means the demo is not shown, and the
   record says so (one reading, condition 3).
2. `layup spec check --facts spans.tsv --psb problem-statement.md` exits 0 on
   `spec/spans.tsv` and `start/problem-statement.md` of the records.
3. `spec/check.tsv` of the records has no `fail` row.
4. Each commit of `layup-records`, the Intake issue, the control issue and each
   comment of LAYUP has the App's bot, `layup-agent[bot]`, as its actor; the push
   of the root commit of `main` is the Operator's, by `root-push`. Intake pushes
   no task branch, so no push of one is read.

## The scenario (uat)

```text
Scenario: The demo of M2c on a pilot problem statement

Given the target pharzam/layup-uat-m2c, made by Start from the host with the pilot brief PSB-CHAT-001, and two registered harnesses, Claude Code and Devin, routed for the roles of Intake
When the restarts of layup run run the steps of Intake until the batch waits for its answers
Then one comment of the App on the Intake issue holds the rule gaps and the gaps of meaning of the brief, each with its question and its quote
And layup spec check --facts passes on the numbered spans that the records hold
And each write of LAYUP on the target shows the App's bot as its author

Accepted: not yet; the Operator's sign-off, with its date, is written here after the run
Covers REQ-001, REQ-012
```

The scenario rides the e2e path of row 48 (`TestTheStepsOfIntake`,
[`runs/T-mo9w/test-runs.md`](../T-mo9w/test-runs.md)). The Operator signs it
off, as for `M2b` ([`runs/T-x7cs/README.md`](../T-x7cs/README.md)).
