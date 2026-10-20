# T-c826 plan-review conditions

Binding sources: the supplied plan review 6301028661, amendment 1 and the
Operator's recorded decisions, received before S1. These rows record the supplied
authorization; this stage does not claim a fresh forge verification.

| Condition | Evidence |
| --- | --- |
| 1. Narrow the exception to the checker for one numbered PR; retain every restoration path. | Adopted before S1; `runs/T-c826/checks.sh` asserts the narrowed form. Workflow implementation and the actual PR number remain for S2. Removal issue: #452, T-szvp. No pin, blob comparison or restore replay. |
| 2. Keep both runners unchanged; drop mode and symlink rules. | `runs/T-c826/checks.sh` compares both runners with d27b331. S1 uses only existing .in overlays and ABSENT.in; there are no executable-mode or symlink cases. |
| 3. Obtain the Operator's CI-route and goal-count decisions before S1. | The binding input supplies the Operator's second R4 approval, alongside reviewer 6301028661, and the final goal count of 2 kept in one task. Reading 2 applies; removal issue #452 already exists. `runs/T-c826/test-runs.md` records S1 only; the coordinating session retains issue-thread and later CI evidence. |
