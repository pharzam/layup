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

**A change of the check before its green run (plan step 3):** check 6 also accepts the task IDs of the table "The parent tasks" of the plan, because `PRD-0001` §12 names the parent tasks `T-vk3k` (`layup gate`) and `T-b97r` (`layup setup`), which are not rows of the task table (O-121).

**A fix of the check (plan step 3):** check 8 read the rows of the area-summary table of the inventory as items; it now reads only the rows whose first cell is an item key.

## Green (plan step 4, after the plan, the registration and the fill of the Host column)

```text
info  inventory: 103 items, 41 conflicts
PASS  plan holds no marker
info  plan: 15 milestones, 20 tasks, 35 other hosts, 41 register rows
PASS  phase-1 requirements have a task
PASS  each task names a requirement
PASS  predecessors resolve with no cycle
PASS  each inventory item has one host
PASS  each conflict has a host
PASS  PRD-0001 section 12 Task cells
PASS  traceability table
PASS  inventory Host column agrees with the plan
exit 0
```
