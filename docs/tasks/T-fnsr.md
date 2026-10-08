# T-fnsr — row 27 of the plan, the demo of M2a (uat)

Issue: [#132](https://github.com/pharzam/layup/issues/132), row 27 of [the tasks
of M2a](../plan/README.md#the-tasks-of-m2a), opened by `T-zwke` (#124). Serves
`F-0003#42` through `NFR-001`, `NFR-002`, `NFR-006` and `REQ-002`. Base
`76d52a1` (the merge of #145). Author: Claude Opus 5.5 on Claude Code.
Evidence: [`runs/T-fnsr/`](../../runs/T-fnsr/).

## Inputs, plan and plan review

The Operator's inputs (comment 6060791998): the target `pharzam/layup-uat`
(public, empty), the App's installation, the key file, the two host registers,
and the values of the command. The plan (R12, 6060858558) counted one goal: the
evidence of one real Start. Its review (Claude Fable 5.1, effort `xhigh`, on
the Claude Code CLI with stream output, a fresh read-only session in a clone at
`76d52a1`, 8 min 44 s; comment 6061031750): `approve-with-conditions`, one goal;
four conditions (no check of `pin.commit` against `armature.pin`, as Start pins
the baseline's head; a red for each check; the source of the bot's ID; the
Operator's read of `main` and the issues) and thirteen notes; Budget maximum 600
lines over 12 files against `76d52a1`; Cycle cap 2 (`check.sh` decides the
result, the reading of T-evad). The author's answer (6061036376) took all of
them. The Operator's confirmation (6061142551): heartbeat commits checked apart
from the seven; the negative tests offline; consent to publish the host name,
the login, the ID and the key path; no growth of the budget without a decision.

## What was done

1. [`check.sh`](../../runs/T-fnsr/check.sh) `WEB API OWNER/NAME`: a plain clone
   and unauthenticated reads; it checks the author and committer of the root
   commit and of each records commit, the seven commits of Start in order, the
   five files, the README of `run.md`, the brief byte for byte, each value of
   `start.tsv`, `approvers.tsv`, `lease.tsv`, and the author of the two issues.
2. [`check-test.sh`](../../runs/T-fnsr/check-test.sh): 36 offline cases on a
   local bare target and saved API responses; each `FAIL` line seen once, and a
   heartbeat commit passes. It found a defect of `check.sh` (`[bot]` as a
   bracket expression in `grep`), fixed before the run.
3. The Operator's first Start stopped with exit 2 before its first step: the
   author's proposal `--intake-cap 5,1` is not of the type `decimal`; the value
   became `5.0,1.0`, the same amounts (comment 6065525521).
4. The Operator's second Start: exit 0, nine `done` rows; the root push with
   the Operator's login; the Operator's read with a plain `git clone`;
   `check.sh` green on the real target. Each write of LAYUP on the target is by
   `layup-agent[bot]`: the root commit, the seven records commits, issues #1 and
   #2.
5. The documents: the traceability (a uat row), the Test cells of `PRD-0001`
   §12 for REQ-002, NFR-001, NFR-002 and NFR-006, and a §13 row.

**The rejected alternatives:** the author runs Start and the Operator only the
push (the uat row names the Operator); a restart on the target (the e2e of row
26 covers it); a private target (the plan is `free`).

**Known limits:** the two of the progress lines of `T-mqty` were seen in this
run (the step name twice in a wait line; a beat every 150 s between the wait
lines). `check.sh` reads the issues with no token, so it needs a public target.

## Review rounds

## Verdict

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC on 2026-10-08; tokens are the
`result` event of the Claude Code CLI (stream runs); `not reported` otherwise.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The inputs; the reading; the goal count; the plan | reasoning | Claude Opus 5.5 | max | not reported | to 13:24 (the inputs comment 13:20, the plan 13:24) |
| Its plan review | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 784,545 (USD 5.75) | 8 min 44 s, from 13:24 |
| The answer; the checks, test first | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 13:34 to 13:47 (`cd45a1a`) |
| The correction of the intake cap; the Operator's runs (the root commit 17:34:57, the records 17:44:51) and read; the evidence; the documents | execution | Claude Opus 5.5 | max | not reported | 17:30 to 17:55 |
