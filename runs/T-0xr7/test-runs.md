# T-0xr7 test runs

## Red, before any edit of the instructions, the blocks, the goldens, the driver or the documents

The hook refuses a red commit (guardrails.md section 2), so the tests below are committed with their code; this is the red, run on 3efd52b8 plus the tests only.

### The three rendered prompts (unit)

Command: `go test -count=1 -run TestEachStepOfAnArchitectureTaskStatesThatTheSurvivorsAreTheBets ./internal/run/`
Start 2026-10-22T18:41:00Z, end 2026-10-22T18:41:00Z.

```
--- FAIL: TestEachStepOfAnArchitectureTaskStatesThatTheSurvivorsAreTheBets (0.00s)
    instructions_test.go:728: the prompt of plan-of-a-task for a task of kind architecture does not hold "A task of kind architecture writes as decided each decision that has one survivor among the survivors listed above"
    instructions_test.go:731: the constant of plan-of-a-task does not hold "A task of kind architecture writes as decided each decision that has one survivor among the survivors listed above"
    instructions_test.go:728: the prompt of plan-of-a-task for a task of kind architecture does not hold "writes each decision that has more than one as a section \"Decision pending the bet\" that lists its survivors and decides nothing, because the bet chooses among them"
    instructions_test.go:731: the constant of plan-of-a-task does not hold "writes each decision that has more than one as a section \"Decision pending the bet\" that lists its survivors and decides nothing, because the bet chooses among them"
    instructions_test.go:728: the prompt of plan-of-a-task for a task of kind architecture does not hold "correct a plan that chooses among survivors, never ask that task to choose among them, and never treat the deferral itself (no option chosen for a decision with more than one survivor) as a defect"
    instructions_test.go:731: the constant of plan-of-a-task does not hold "correct a plan that chooses among survivors, never ask that task to choose among them, and never treat the deferral itself (no option chosen for a decision with more than one survivor) as a defect"
    instructions_test.go:728: the prompt of plan-of-a-task for a task of kind architecture does not hold "check a decision of one survivor, a survivor that is wrong, a choice hidden in a pending section, a contradiction, and the rest of the work, as any other"
    instructions_test.go:731: the constant of plan-of-a-task does not hold "check a decision of one survivor, a survivor that is wrong, a choice hidden in a pending section, a contradiction, and the rest of the work, as any other"
    instructions_test.go:728: the prompt of plan-review for a task of kind architecture does not hold "A task of kind architecture writes as decided each decision that has one survivor among the survivors listed above"
    instructions_test.go:731: the constant of plan-review does not hold "A task of kind architecture writes as decided each decision that has one survivor among the survivors listed above"
    instructions_test.go:728: the prompt of plan-review for a task of kind architecture does not hold "writes each decision that has more than one as a section \"Decision pending the bet\" that lists its survivors and decides nothing, because the bet chooses among them"
    instructions_test.go:731: the constant of plan-review does not hold "writes each decision that has more than one as a section \"Decision pending the bet\" that lists its survivors and decides nothing, because the bet chooses among them"
    instructions_test.go:728: the prompt of plan-review for a task of kind architecture does not hold "correct a plan that chooses among survivors, never ask that task to choose among them, and never treat the deferral itself (no option chosen for a decision with more than one survivor) as a defect"
    instructions_test.go:731: the constant of plan-review does not hold "correct a plan that chooses among survivors, never ask that task to choose among them, and never treat the deferral itself (no option chosen for a decision with more than one survivor) as a defect"
    instructions_test.go:728: the prompt of plan-review for a task of kind architecture does not hold "check a decision of one survivor, a survivor that is wrong, a choice hidden in a pending section, a contradiction, and the rest of the work, as any other"
    instructions_test.go:731: the constant of plan-review does not hold "check a decision of one survivor, a survivor that is wrong, a choice hidden in a pending section, a contradiction, and the rest of the work, as any other"
    instructions_test.go:728: the prompt of verification for a task of kind architecture does not hold "A task of kind architecture writes as decided each decision that has one survivor among the survivors listed above"
    instructions_test.go:731: the constant of verification does not hold "A task of kind architecture writes as decided each decision that has one survivor among the survivors listed above"
    instructions_test.go:728: the prompt of verification for a task of kind architecture does not hold "writes each decision that has more than one as a section \"Decision pending the bet\" that lists its survivors and decides nothing, because the bet chooses among them"
    instructions_test.go:731: the constant of verification does not hold "writes each decision that has more than one as a section \"Decision pending the bet\" that lists its survivors and decides nothing, because the bet chooses among them"
    instructions_test.go:728: the prompt of verification for a task of kind architecture does not hold "correct a plan that chooses among survivors, never ask that task to choose among them, and never treat the deferral itself (no option chosen for a decision with more than one survivor) as a defect"
    instructions_test.go:731: the constant of verification does not hold "correct a plan that chooses among survivors, never ask that task to choose among them, and never treat the deferral itself (no option chosen for a decision with more than one survivor) as a defect"
    instructions_test.go:728: the prompt of verification for a task of kind architecture does not hold "check a decision of one survivor, a survivor that is wrong, a choice hidden in a pending section, a contradiction, and the rest of the work, as any other"
    instructions_test.go:731: the constant of verification does not hold "check a decision of one survivor, a survivor that is wrong, a choice hidden in a pending section, a contradiction, and the rest of the work, as any other"
FAIL
FAIL	github.com/pharzam/layup/internal/run	0.004s
FAIL
```

### The copied command at the cap, the wait message, and the homes (integration)

Command: `go test -count=1 -tags=integration -run TestAGrantOnTheControlIssue|TestAMaterialRoundStartsOneAttemptAndTheSecondWaitsAtTheCap|TestTheInstructionContractOfTheSurvivors|TestTheCycleCommandHasItsHome ./internal/run/`
Start 2026-10-22T18:41:00Z, end 2026-10-22T18:41:03Z.

```
--- FAIL: TestTheInstructionContractOfTheSurvivorsHasItsHomeInLoopMd (0.00s)
    promptfile_integration_test.go:541: loop.md does not hold "The contract of the survivors"
    promptfile_integration_test.go:541: loop.md does not hold "The review scope of the survivors"
--- FAIL: TestTheCycleCommandHasItsHomeInLoopMdAndRecordsMd (0.01s)
    shapeloop_integration_test.go:1077: ../../docs/spec/loop.md does not hold "One more cycle at the cap"
    shapeloop_integration_test.go:1077: ../../docs/spec/loop.md does not hold "cycle <task> because <reason>"
    shapeloop_integration_test.go:1077: ../../docs/spec/records.md does not hold "A valid copy of a `cycle` command is the grant"
--- FAIL: TestAMaterialRoundStartsOneAttemptAndTheSecondWaitsAtTheCap (2.88s)
    shapeloop_integration_test.go:172: AdvanceShapeTask = {Task:T-izn4 Verified:false Merged:false Head: Attempt:0 Wait:waiting for the Operator: T-izn4: the review-cycle cap 1 is reached with findings open Stop: Replaced:}, <nil>; want the wait "waiting for the Operator: T-izn4: the review-cycle cap 1 is reached with findings o
--- FAIL: TestAGrantOnTheControlIssueIsCopiedAndTheTaskAtTheCapRunsOneMoreAttempt (3.03s)
    shapeloop_integration_test.go:1047: AdvanceShapeTask after the grant = {Task:T-74dk Verified:false Merged:false Head: Attempt:0 Wait:waiting for the Operator: T-74dk: the review-cycle cap 1 is reached with findings open Stop: Replaced:}, <nil>; want the verified head
FAIL
FAIL	github.com/pharzam/layup/internal/run	3.054s
FAIL
```

### The binary: the Operator's command after the wait (e2e)

Command: `go test -count=1 -tags=e2e -run TestTheSpecificationTaskEndToEnd ./cmd/layup/`
Start 2026-10-22T18:41:03Z, end 2026-10-22T18:41:22Z.

```
--- FAIL: TestTheSpecificationTaskEndToEnd (18.50s)
    loop_e2e_test.go:455: stage 5a: the shape row ["done" "waiting for the Operator: T-bsuu: the review-cycle cap 1 is reached with findings open"]; want the wait "waiting for the Operator: T-bsuu: the review-cycle cap 1 is reached with findings open; the Operator may grant one more cycle with `cycle T-bsuu because <reason>` on
        step	result	detail
        forge	done	the installation token, the six capabilities and the permissions of M2a; a public repository, with layup-records
        clone	done	start.tsv, approvers.tsv and lease.tsv read
        version	done	LAYUP 0.1.0-dev
```

## Green, on 7da43c52

- Build, vet (no tag, integration, e2e), gofmt, unit level: exit 0.
- Full integration level: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...`, output in `/home/layup/.cache/layup-session/T-0xr7-full-integration.txt`. Start 2026-10-22T19:00:37Z, end 2026-10-22T19:04:54Z, exit 0 (21 packages ok).
- E2e level: `go test -count=1 -tags=e2e ./cmd/layup/`. Start 2026-10-22T19:05:10Z, end 2026-10-22T19:06:41Z, exit 0 (`ok github.com/pharzam/layup/cmd/layup 89.786s`).
- Lints: adr-lint OK, prd-lint OK, link-lint OK (3165 links), run-discipline-tests 138 passed, nested-checkout-check OK, setup-check OK.
