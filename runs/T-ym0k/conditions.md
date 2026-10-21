# T-ym0k: the conditions of the plan review of #369, and their evidence

One row per condition of the plan review (comment 6314912248), per decision D1 to D5 and per question Q1 to Q4. This table is the input of the conditions check of the freeze. Lines are those of the head that carries this file. Condition 4 and Q1 to Q3 wait for the review (S6), which this head has not run; the row says what is still to do.

| Condition | Evidence |
| --- | --- |
| 1: the red slice first: the control, the mutants m1 to m4 (m2 and m4 labelled regression checks), the doctored records and both empty-list cases red against an incomplete comparator, the `Read` refusal, the D3 output assertion | `runs/T-ym0k/test-runs.md:9` (Red 1, with the control and the assertion of D3), `runs/T-ym0k/test-runs.md:31` (Red 2: two incomplete comparators), `runs/T-ym0k/coverage-check-test.sh:52-69` (the cases) |
| 2: the claim of D4 bounded to what the guard proves: count the occurrences and the parsed rows, reject the rest, two empty counts are no coverage; tested before its code | `runs/T-efmy/release-check.sh:173-189`, `runs/T-efmy/release-check.sh:74-81` (the header states the claim and the retained limit), `runs/T-ym0k/test-runs.md:40` (m3, m5, m6 and the empty adapter, red against the control first) |
| 3: both halves of the comparison bound to COMMIT, with a case that exposes an index/COMMIT mismatch | `runs/T-ym0k/coverage-check.sh:12-29` (a temporary index of git read-tree), `runs/T-ym0k/coverage-check-test.sh:64-68` (HEAD at B with a staged file out of B; the records of A and B for the other commit) |
| 4: the changed check and the comparator are reviewed by an independent pass before the paid review; the 60-minute allowance needs the Operator's recorded exception | Not done by this head: the review of the producing code, the Operator's approval of the spend and the exception are the next steps (`docs/tasks/T-ym0k.md:20`) |
| 5: the scope reworded (a new `layup` command is out, the shell comparator stays); D1 to D5 with their reasons in the specification; the close-out in this PR | `docs/tasks/T-ym0k.md:16` (D5: a script, not a command), `docs/spec/loop.md:2932-2943` (D1 to D5 with reasons); the verdict, the completed log, the backlog and the resource record are written at the close-out |
| 6: the size classified: the review, the output and the run log are evidence; the scripts, the test of the comparator and the task file count; scratch inputs outside the branch | `runs/T-ym0k/test-runs.md`, `runs/T-ym0k/release-check.txt` (declared evidence); `runs/T-ym0k/coverage-check.sh`, `runs/T-ym0k/coverage-check-test.sh`, `runs/T-efmy/release-check.sh` (counted); the mutants, the control, the brief and the clone are in ~/.cache/layup-session/t-ym0k and ~/.cache/layup-session/reviews/issue369-release |
| 7: the factual statements repaired (PR #466 for #368; embedded +112/-1; T-3q0b in the trailerless group; T-2vva did no work; m3: the grep omits the call) | `docs/tasks/T-ym0k.md:8`, `docs/tasks/T-ym0k.md:20`, `docs/plan/README.md:371-380`, `runs/T-ym0k/test-runs.md:9` |
| D1: the full ID | `docs/plan/README.md:371-380`, `docs/spec/loop.md:2935` |
| D2: the exact pinned-linter call | `runs/T-efmy/release-check.sh:118`, `runs/T-ym0k/test-runs.md:40` (m1) |
| D3: `RemoteHead` in (7) | `runs/T-efmy/release-check.sh:71`, `runs/T-efmy/release-check.sh:201`, `runs/T-ym0k/release-check.txt:132` |
| D4: the omission alarm of (6) | `runs/T-efmy/release-check.sh:173-189` |
| D5: the comparator in the tree | `runs/T-ym0k/coverage-check.sh:1-20`, `runs/T-ym0k/coverage-check-test.sh:1-8` |
| Q1: the reviewer, Claude Fable 5.1, then the eligible fallback | Not run by this head (S6) (`docs/tasks/T-ym0k.md:20`) |
| Q2: the effort and the spend | Not run by this head (S6): the Operator's approval of the spend is still to come (`docs/tasks/T-ym0k.md:20`) |
| Q3: the time limit | Not run by this head (S6): the exception is the Operator's (`docs/tasks/T-ym0k.md:20`) |
| Q4: bounded fixes in this PR | Not applicable yet: no finding exists before the review (`docs/tasks/T-ym0k.md:20`) |
