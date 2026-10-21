# T-8s13: the conditions of the plan review of #370 and the Operator's decisions, and their evidence (step T)

One row per condition of the plan review (comment 6321481585) and per decision of the Operator (comment 6321600725). This table is check K6's input; K6 proves that each cited file exists, not that it meets the condition.
The live run (R0 to item 10) is not part of step T; its rows are added when it has run.

| Condition | Evidence |
| --- | --- |
| 6321481585 1: the exact formula with E, `stillOpen`'s rule, G1 at 79 (76 plus one decision of two options), G2's 20-line floor, `brief.lines` 100 | `runs/T-8s13/README.md:48`, `runs/T-8s13/README.md:56`, `runs/T-8s13/README.md:61`, `runs/T-8s13/brief-bound.sh:1`, `internal/run/betbrief.go:364`, `internal/run/shapequestions.go:25` |
| 6321481585 2: the red step explicit and reproducible: the test before the counter, an assertion failure, then green; the open-question case; the 80/81 boundary; the M2d comparison as a supplement; fixtures counted as inputs | `runs/T-8s13/brief-bound-test.sh:1`, `runs/T-8s13/test-runs.md:8`, `runs/T-8s13/test-runs.md:17`, `runs/T-8s13/test-runs.md:52` |
| 6321481585 3 and 6321600725 2: S1 is the specification's set before the first bet, answer-derived `Q-NNN` constraints not checked, nothing repaired | `runs/T-8s13/README.md:106`, `docs/spec/loop.md:1741`, `internal/run/options.go:73` |
| 6321481585 4 and 6321600725 1: AC2 reads "if a code-owner review occurs it is the Operator's"; none occurs; every pull request's reviews are `[]`; the wording synchronized | `runs/T-8s13/README.md:115`, `runs/T-8s13/README.md:120`, `docs/spec/loop.md:2950` |
| 6321481585 5: a status or check counts only strictly before the event; a tie is inconclusive and fails unless independent ordering evidence exists, which is named | `runs/T-8s13/README.md:68`, `runs/T-8s13/README.md:89` |
| 6321481585 6: rule 1 reconciles the actual transitions of each merged task with its events and the transition table; a missing row fails as does an invalid one | `runs/T-8s13/README.md:74`, `docs/spec/loop.md:704`, `docs/architecture.md:836` |
| 6321481585 note 6: the stale citations corrected (the late-status control at `cmd/layup/loop_e2e_test.go:588`, Probe's call at `internal/run/probe.go:224`) | `runs/T-8s13/README.md:96`, `runs/T-8s13/README.md:143`, `cmd/layup/loop_e2e_test.go:588`, `internal/run/probe.go:224` |
| 6321600725 3: the codex register row's context size is a named open input, not invented | `runs/T-8s13/README.md:140`, `runs/T-8s13/README.md:150`, `runs/T-8s13/README.md:169` |
| 6321600725 4: the SSH alias is `LayupZAM` | `runs/T-8s13/README.md:146` |
| Plan section 8, R10 (one home each): the target wording, the e2e repair rule not applying, the permission wording of `forge.md`, row 68b and the apply of the code-owner rule | `docs/spec/loop.md:2939`, `docs/spec/loop.md:2950`, `docs/spec/loop.md:2953`, `docs/plan/README.md:383`, `docs/plan/README.md:432` |
