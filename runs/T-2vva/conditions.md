# T-2vva: the conditions of the plan review of #366 and of round 1, and their evidence

One row per condition of the plan review (comment 6250697855), the plan amendment 1 (6293953072) and the findings of round 1 (6295520491). This table is check K6's input; K6 proves that each cited span exists, not that it meets the condition.

| Condition | Evidence |
| --- | --- |
| 6250697855 1: the splits and the claim | `docs/plan/README.md`, `runs/T-2vva/test-runs.md`, `3d97225` |
| 6250697855 2: the predecessors, the shared contracts | `internal/run/bet.go`, `internal/run/shapetask.go`, `internal/run/rendered.go`, `TestTheStepBet` |
| 6250697855 3: the brief contract (66u's) and the post | `internal/run/bet.go`, `internal/run/handoffpost.go`, `TestTheStepBet` |
| 6250697855 4: the appetite (66u's) | `docs/spec/loop.md`, `internal/run/betbrief.go` |
| 6250697855 5: recovery across the version and the registration | `TestTheBetStepIsAtomicToItsLastAct`, `TestShapeStaysDoneAfterTheRenderedPush`, `internal/run/bet.go` |
| 6250697855 6: the settled decisions and the opening | `TestTheSettledDecisionsAreTheOptionLines`, `TestTheRequirementsOfAnOptionAreItsMappedRows`, `TestADecisionWhoseOptionNamesNoRequirementIsRefusedBeforeAnyWrite`, `internal/run/archtasks.go`, `internal/run/shapetask.go` |
| 6250697855 7: the allocation of Q15 to row 66k | `docs/plan/README.md`, `docs/spec/loop.md` |
| 6250697855 8: the red slice and the evidence | `TestABetsTableRefusesEachBrokenRule`, `TestTheNextVersionOfTheInventory`, `runs/T-2vva/test-runs.md`, `runs/T-2vva/mutants.md`, `423bf25` |
| 6293953072 1: the phases in the specification instruction | `internal/run/instructions.go`, `internal/run/testdata/prompts/specification.txt`, `docs/spec/loop.md`, `TestThePhaseTagsAreTheFirstColumnOfSection9` |
| 6293953072 2: the guards of a bound of code | `TestACodeBoundWithNoSessionIsSkippedByItsReaders`, `TestAPushAndABoundOfCodeHaveNoSession`, `TestShapeStaysDoneAfterTheRenderedPush` |
| 6295520491 1: recovery after the rendered push, before its bound | `TestTheBetStepIsAtomicToItsLastAct`, `TestACodeBoundWithNoSessionIsSkippedByItsReaders`, `internal/run/shapeloop.go`, `427ff06` |
| 6295520491 2: coverage (G6, G7, the reasons, the close-out) | `TestShapeStaysDoneAfterTheRenderedPush`, `TestABetsTableRefusesEachBrokenRule`, `TestTheStepBet`, `runs/T-2vva/mutants.md` |
| 6295520491 3: malformed phase tables | `TestThePhaseTagsAreTheFirstColumnOfSection9`, `internal/run/bet.go` |
| 6295520491 4: the branch-tip refusal | `TestARenderedPushToAnotherBranchTipIsRefused`, `internal/run/bet.go` |
