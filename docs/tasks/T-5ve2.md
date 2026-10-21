# T-5ve2 — the shared index files merge without conflicts

Issue [#461](https://github.com/pharzam/layup/issues/461), B3 of the M2e throughput investigation [#438](https://github.com/pharzam/layup/issues/438). It serves `F-0003#47`.
Plan [6314092825](https://github.com/pharzam/layup/issues/461#issuecomment-6314092825); plan review [6314252349](https://github.com/pharzam/layup/issues/461#issuecomment-6314252349) (approve-with-conditions); the Operator's decisions Q1 to Q4 [6315292072](https://github.com/pharzam/layup/issues/461#issuecomment-6315292072).
Budget maximum 400 lines over 10 files, close-out inside, under ADR-0028; cycle cap 1.
Evidence: [test runs](../../runs/T-5ve2/test-runs.md), [forge merge](../../runs/T-5ve2/forge-merge.md).

## What is delivered

- `merge=union` for `docs/tasks/completed.md` and `docs/tests/traceability.md` only, in [`.gitattributes`](../../.gitattributes), with its reason and limits.
- The insertion rule in each file's head: a line directly under `## Log`, a row directly below the table's separator, so the later-landing branch is on top in both.
- `TestTheSharedIndexesMergeWithNoConflictInLandingOrder` (`internal/git/indexmerge_integration_test.go`), red before the attribute and green after.
- The narrow review exception in [Reviewing until findings decay](../engineering-discipline.md#reviewing-until-findings-decay), and one lesson in `guardrails.md` §2.

## The comparison, per file

| File | Choice | Why not union or fragments |
| --- | --- | --- |
| `completed.md`, `traceability.md` | union, with the insertion rule | Readers check presence, never order. Fragments would change `kit-history` (sh and Go), `task-state.sh`, `task-labels.sh`, K2 and the targets' close-out. |
| `backlog.md` | none | Edits are deletions: adjacent deletions come back under union. |
| PRD §12 and §13 | none; follow-up (Q1) | Two edits of one matrix row become two rows, which `prd-lint` passes. |
| `guardrails.md` | none (Q2) | An added lesson can contradict the branch's text. |
| `built` in `internal/tsv/blocks_integration_test.go` | unchanged (review condition 3) | Union reaches code hunks of a `.go` file, and a duplicated `var` line does not compile. Fragments, one file per name, need a loader and a schema change for a test that checks each block once and passes today. Reformatting it is another deliverable. |

## Not in scope

The PRD §12 rule and any K2 change, reordering existing `completed.md` lines, a custom merge driver, and the reformat of `built`.
