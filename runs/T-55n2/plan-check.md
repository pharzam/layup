# T-55n2 — the one-off check of the plan

The check is [`plan-check.py`](plan-check.py), run by hand from the repository root with `python3 runs/T-55n2/plan-check.py`. No hook and no CI job runs it (Bootstrap mode rule 1).

## Red (plan step 1, at the base `7cdd346` plus the inventory, before the plan)

```text
info  inventory: 103 items, 41 conflicts
FAIL  plan exists
      docs/plan/README.md is missing
FAIL  phase-1 requirements have a task
      no plan
FAIL  each task names a requirement
      no plan
FAIL  predecessors resolve with no cycle
      no plan
FAIL  each inventory item has one host
      no plan
FAIL  each conflict has a host
      no plan
FAIL  PRD-0001 section 12 Task cells
      no plan
FAIL  traceability table
      no plan
FAIL  inventory Host column agrees with the plan
      no plan
exit 1
```
