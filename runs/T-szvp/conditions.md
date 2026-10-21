# T-szvp plan-review conditions

Binding sources: the supplied approve-with-conditions plan review and Plan
amendment 1, posted as issue #452 comment 6304288973. There is one condition;
the review's notes do not add conditions.

| Condition | Evidence |
| --- | --- |
| 1. Keep one assertion of the restored step, not two. | `runs/T-c826/checks.sh:3-43` removes only the PR-number interface, keep/number/environment assertions, pinned skip body and exception-comment assertion. Its margin, unchanged-runner, nine-path, source and default-branch assertions remain. The sole complete Restore-step comparison stays in the unchanged `runs/T-szvp/checks.sh:33-40`. `docs/tests/README.md:67-71` names both commands without arguments. Both task commands exit 0 in Stage B: the T-c826 check passes all six remaining assertions and the unchanged T-szvp check passes its complete comparison. `runs/T-szvp/test-runs.md` records their times, exits and output. |
