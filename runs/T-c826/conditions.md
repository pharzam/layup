# T-c826 plan-review conditions

Binding sources: the supplied plan review 6301028661, amendment 1 and the
Operator's recorded decisions, received before S1. These rows record the supplied
authorization; this stage does not claim a fresh forge verification.

| Condition | Evidence |
| --- | --- |
| 1. Narrow the exception to the checker for one numbered PR; retain every restoration path. | Implemented in S2 for supplied PR #453: event pull_request, literal number 453 and only docs/tests/freeze-check.sh. All nine restore-list paths remain; both runners and every other listed script restore for every event. The unchanged task check passes all CI assertions with argument 453; full setup-check exits 0, including ci. R4 approvals 6301028661 and 6301084598 and removal issue #452 (T-szvp) are named in the workflow comment. Hosted red CI evidence remains for the coordinating session before S3. No pin, blob comparison or restore replay. |
| 2. Keep both runners unchanged; drop mode and symlink rules. | `runs/T-c826/checks.sh` compares both runners with d27b331. S1 uses only existing .in overlays and ABSENT.in; there are no executable-mode or symlink cases. |
| 3. Obtain the Operator's CI-route and goal-count decisions before S1. | The binding input supplies the Operator's second R4 approval 6301084598, alongside reviewer 6301028661, and the final goal count of 2 kept in one task. Reading 2 applies; removal issue #452 already exists. `runs/T-c826/test-runs.md` records S1 and S2; it distinguishes the supplied PR identity from the failed network verification. The coordinating session retains issue-thread evidence and must capture red CI before checker implementation. |
