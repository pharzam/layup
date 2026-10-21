# T-5ve2: the conditions of the plan review of #461, and their evidence

One row per condition of the plan review (comment 6314252349) and per Operator decision Q1 to Q4 (comment 6315292072). This table is the input of the conditions check of the freeze. Lines are those of the head that carries this file.

| Condition | Evidence |
| --- | --- |
| 1: the round exception is narrow: only the two union indexes, only new records from main; existing records preserved, no conflicting versions, headings and instructions unchanged, no mode or path change, the numstat endpoints recorded | `docs/engineering-discipline.md:305-315` |
| 2: planned-row delivery and corrections stay ordinary edits outside the exception; no blanket "add, don't edit"; the `—` and `F-0003#47` row is justified by K2 and not by the frame precedent | `docs/tests/traceability.md:27-41`, `docs/tests/traceability.md:45` |
| 3: step 3 (the `built` reformat) is removed, and the comparison says why union and fragments are rejected for it | `docs/tasks/T-5ve2.md:23` |
| 4: the test reuses `isolate`, `plain` and `gitIn`; has the PRD fixture; gives `git diff --check` explicit endpoints; asserts B before A in both files; the red is the merge conflict | `internal/git/indexmerge_integration_test.go:94-97`, `internal/git/indexmerge_integration_test.go:149`, `internal/git/indexmerge_integration_test.go:126-129`, `runs/T-5ve2/test-runs.md` |
| 5: the base is the current main merged first; evidence and bookkeeping leave the count; the forge experiment is not part of the core validation | `runs/T-5ve2/test-runs.md`, `runs/T-5ve2/forge-merge.md` |
| Q1: the PRD section 12 rule is a follow-up, not here; the PRD takes no attribute | `.gitattributes:284-285`, `internal/git/indexmerge_integration_test.go:167-195` |
| Q2: no widening of the exception to guardrails.md | `docs/engineering-discipline.md:305-315`, `internal/git/indexmerge_integration_test.go:169-174` |
| Q3: new rows go directly below the separator; the test asserts B before A in both files | `docs/tests/traceability.md:33-41`, `docs/tasks/completed.md:17-24`, `internal/git/indexmerge_integration_test.go:126-129` |
| Q4: the forge experiment is the coordinator's; the record is left for it | `runs/T-5ve2/forge-merge.md` |
