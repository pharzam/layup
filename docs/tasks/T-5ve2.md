# T-5ve2 — the shared index files merge without conflicts

Issue [#461](https://github.com/pharzam/layup/issues/461), B3 of the M2e throughput investigation [#438](https://github.com/pharzam/layup/issues/438). It serves `F-0003#47`.
Plan [6314092825](https://github.com/pharzam/layup/issues/461#issuecomment-6314092825); plan review [6314252349](https://github.com/pharzam/layup/issues/461#issuecomment-6314252349) (approve-with-conditions); the Operator's decisions Q1 to Q4 [6315292072](https://github.com/pharzam/layup/issues/461#issuecomment-6315292072).
Budget maximum 400 lines over 10 files, close-out inside, under ADR-0028; cycle cap 1.
Goal count 1. Round 1: comment 6315695160.
Evidence: [test runs](../../runs/T-5ve2/test-runs.md), [forge merge](../../runs/T-5ve2/forge-merge.md), [conditions](../../runs/T-5ve2/conditions.md).

## What is delivered

- `merge=union` for `docs/tasks/completed.md` and `docs/tests/traceability.md` only, in [`.gitattributes`](../../.gitattributes), with its reason and limits.
- The insertion rule in each file's head: a line directly under `## Log`, a row directly below the table's separator, so the later-landing branch is on top in both.
- `TestTheSharedIndexesMergeWithNoConflictInLandingOrder` (`internal/git/indexmerge_integration_test.go`), red before the attribute and green after.
- The forge result: GitHub's own merge and "Update branch" do not honour `merge=union`, so a branch takes main by a local `git merge origin/main` ([the record](../../runs/T-5ve2/forge-merge.md)).
- The narrow review exception in [Reviewing until findings decay](../engineering-discipline.md#reviewing-until-findings-decay), and one lesson in `guardrails.md` §2.

## The comparison, per file

| File | Choice | Why not union or fragments |
| --- | --- | --- |
| `completed.md`, `traceability.md` | union, with the insertion rule | Readers check presence, never order. Fragments would change `kit-history` (sh and Go), `task-state.sh`, `task-labels.sh`, K2 and the targets' close-out. |
| `backlog.md` | none | Edits are deletions: adjacent deletions come back under union. |
| PRD §12 and §13 | none; follow-up (Q1) | Two edits of one matrix row become two rows, which `prd-lint` passes. |
| `guardrails.md` | none (Q2) | An added lesson can contradict the branch's text. |
| `built` in `internal/tsv/blocks_integration_test.go` | unchanged (review condition 3) | Union reaches code hunks of a `.go` file, and a duplicated `var` line does not compile. Fragments, one file per name, need a loader and a schema change for a test that checks each block once and passes today. Reformatting it is another deliverable. |

## Not delivered

The PRD section 12 freeze (K2 reading traceability alone), a follow-up issue under Q1 at cap 2, which the coordinator opens.

## Not in scope

The PRD §12 rule and any K2 change, reordering existing `completed.md` lines, a custom merge driver, and the reformat of `built`.

## Verification and verdict

Red first, then green, recorded in the test runs with `date -u` times. Verified by running: build, vet (three tags), gofmt, unit, the full integration level and e2e (exit 0), the six lints, `git diff --check`, and the freeze check (K1 to K6).

Round 1, comment 6315695160 (GPT-6.1 Sol on Codex): **nothing material in scope.** Its notes were applied in a commit of documents only (the forge result and one-home links in the guardrails lesson), which needs no round.

**Verdict: nothing material in scope.** Goal count 1. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 (Plan agent) | not recorded | 129,946 | 2026-10-21T10:50Z to 10:58Z |
| Plan review | reasoning | GPT-6.1 Sol | xhigh | 186,290 | 2026-10-21T10:47:55Z to 10:57:00Z |
| Implementation | execution | Claude Sonnet 5.5 (subagent) | not recorded | 111,697 | 2026-10-21T11:48Z to 11:59Z |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 161,782 | 2026-10-21T12:08:37Z to 12:16:50Z |
| Forge experiment | reasoning | Claude Opus 5.5 (coordinator) | not recorded | not kept | 2026-10-21T13:08Z to 13:12Z |
| Notes and close-out | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept at writing | 2026-10-21T13:09:38Z to 2026-10-21T13:10:01Z |
| **Total** | | | | 589,715 reported (the four numeric entries); the complete total is not reported | not reported (the rows have no common start and end) |
