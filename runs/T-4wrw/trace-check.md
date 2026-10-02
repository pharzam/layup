# T-4wrw — the one-off check of the PDR

The check is [`trace-check.py`](trace-check.py), run by hand from the repository root. No hook and no CI job runs it (Bootstrap mode rule 1). Its header says what it proves and what it cannot prove.

## Red 1: a copy of PRD-0001 with one seeded defect per part (plan-review condition 3)

The copy is `PRD-seeded.md` in a temporary directory: the Facts cell of `REQ-001` (§6) is empty, and the Facts cell of `REQ-002` (§6) also cites `F-0003#76`, `F-0002#1` and `F-0099#1`. `docs/prd/` does not change.

```text
info  numbered facts: F-0001 39, F-0002 0, F-0003 75, F-0004 19
FAIL  each requirement of §6 and §7 cites a fact (25 rows)
      REQ-001: no fact in its Facts cell
FAIL  each fact token resolves (102 distinct tokens)
      F-0002#1: F-0002 has no numbered fact 1 (0 numbered facts)
      F-0003#76: F-0003 has no numbered fact 76 (75 numbered facts)
      F-0099#1: no record F-0099 in docs/facts/
FAIL  each §12 Facts cell equals its §6 or §7 cell
      REQ-001: §12 cites ['F-0003#41', 'F-0003#14', 'F-0003#15', 'F-0001#11'], §6/§7 cite []
      REQ-002: §12 cites ['F-0003#42', 'F-0003#15', 'F-0003#37', 'F-0001#4', 'F-0001#8'], §6/§7 cite ['F-0003#42', 'F-0003#15', 'F-0003#37', 'F-0001#4', 'F-0001#8', 'F-0003#76', 'F-0002#1', 'F-0099#1']
FAIL  the record names each document by a commit of main that holds it
      docs/pdr/PDR-0001.md is missing
FAIL  the record holds the approval
      no record
FAIL  PRD-0001 has the Status `Accepted`
      Status is Draft
exit 1
```

## Red 2: the real tree at `b48764f`, before the record exists

```text
info  numbered facts: F-0001 39, F-0002 0, F-0003 75, F-0004 19
PASS  each requirement of §6 and §7 cites a fact (25 rows)
PASS  each fact token resolves (99 distinct tokens)
PASS  each §12 Facts cell equals its §6 or §7 cell
FAIL  the record names each document by a commit of main that holds it
      docs/pdr/PDR-0001.md is missing
FAIL  the record holds the approval
      no record
FAIL  PRD-0001 has the Status `Accepted`
      Status is Draft
exit 1
```

## Green: the real tree after the approval (plan step 4)

Run on 2026-10-02 on the tree of the commit that adds this section: the record
holds the approval (O-129), and `PRD-0001` has the Status `Accepted`.
`origin/main` was `3fb44d3` (the merge of row 2, #101).

```text
info  numbered facts: F-0001 39, F-0002 0, F-0003 75, F-0004 19
PASS  each requirement of §6 and §7 cites a fact (25 rows)
PASS  each fact token resolves (99 distinct tokens)
PASS  each §12 Facts cell equals its §6 or §7 cell
PASS  the record names each document by a commit of main that holds it
PASS  the record holds the approval
PASS  PRD-0001 has the Status `Accepted`
exit 0
```
