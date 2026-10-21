# T-wqrk: the conditions of the plan review of #462 and of round 1, and their evidence

One row per condition of the plan review (comment 6314145194), the Operator's decisions (6315292816) and the findings of round 1 of #462. This table is check K6's input; K6 proves that each cited file exists, not that it meets the condition.

| Condition | Evidence |
| --- | --- |
| 6314145194 1: a two-home plant, its ground truth frozen | `runs/T-wqrk/trial.md` |
| 6314145194 2: a fixed judge, run order and invalid-run rule, checked on a positive and a mere-citation negative | `runs/T-wqrk/trial.md`, `runs/T-wqrk/harness.md` |
| 6314145194 3: named preparation cases, red on the old preparer, green on the new | `runs/T-wqrk/prep-test.sh`, `runs/T-wqrk/test-runs.md` |
| 6314145194 4: no earlier verdict in the trial inputs or the fix-round path | `runs/T-wqrk/trial.md`, `runs/T-wqrk/prep-test.sh` |
| 6314145194 5: a search with old and new wording and variants; exclusions that hide no operative file | `runs/T-wqrk/harness.md`, `runs/T-wqrk/prep-test.sh` |
| 6314145194 6: the shared-harness protocol for the template and preparer | `runs/T-wqrk/harness.md` |
| 6314145194 7: the producing code reviewed before the six runs; versions preserved | `runs/T-wqrk/harness.md` |
| 6314145194 8: the baseline claims corrected | `runs/T-wqrk/trial.md`, `docs/engineering-discipline.md` |
| 6315292816 Q1 to Q3: Sol at high, 3 and 3 runs one at a time; "not red" stays "not red"; the limit accepted | `runs/T-wqrk/trial.md` |
| Round 1 M1: the judge ties an affirmative diagnosis to this deadline; frozen with its sha256 | `runs/T-wqrk/trial.md`, `runs/T-wqrk/harness.md` |
| Round 1 M2: the `earlier_record` opt-in removed; a fix-round config with `round1.md` is refused | `runs/T-wqrk/prep-test.sh`, `runs/T-wqrk/test-runs.md` |
| Round 1 M3: each refusal case is otherwise valid and asserts its message; each guard removed fails its own case | `runs/T-wqrk/prep-test.sh`, `runs/T-wqrk/test-runs.md` |
| Round 1 M4: exclusions are named records, never by suffix; operative Markdown stays visible | `runs/T-wqrk/prep-test.sh`, `runs/T-wqrk/harness.md` |
| Round 1 M5: one exact heading test and one validity check for both arms | `runs/T-wqrk/trial.md`, `runs/T-wqrk/harness.md` |
| Round 2 F1: the affirmative clause names ADR-0026 and its record deadline, with no denial or hedge; both controls give MISS; the judge frozen with its sha256 | `runs/T-wqrk/trial.md`, `runs/T-wqrk/harness.md` |
| Round 2 F2: the consumer check reads full paths and same-directory names only; the suite exits 0 on the committed head in a clean clone; the guards still fail their own cases | `runs/T-wqrk/prep-test.sh`, `runs/T-wqrk/test-runs.md` |
