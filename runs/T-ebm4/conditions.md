# T-ebm4 plan-review conditions

Stage A records the amended test contract, not completion of the prose changes.
The unchanged documents must fail the new assertions. The same check will run
after the prose changes in the next stage.

| Condition | Evidence |
| --- | --- |
| 1. Require a new issue comment naming the figure or end; edits post nothing. The posting time selects the window. Prose implementation pending. | `runs/T-ebm4/checks.sh` asserts C1-new-comment-wait, C1-retired-wait and C1-one-shared-wait; `runs/T-ebm4/test-runs.md` records the red run. |
| 2. Drop decision 8 and preserve the validation rule. No change to that section is planned. | `runs/T-ebm4/checks.sh` compares the complete validation section against the baseline in C2-validation-preserved; `runs/T-ebm4/test-runs.md` records the result. |
| 3. Bound Records finding by "and changes nothing else" and add a complete UTC row. Prose implementation pending. | `runs/T-ebm4/checks.sh` asserts C3-records-bound and the UTC term, abbreviation, description, collision and example; `runs/T-ebm4/test-runs.md` records the red run. |

These records are evidence, not execution inputs. The task check reads neither
this table nor the test-run record as configuration or expected output. K6 will
audit the table and its citations at the committed head; it remains enabled.
