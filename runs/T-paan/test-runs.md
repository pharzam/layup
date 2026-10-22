# T-paan test runs

## Red, before any edit of the blocks, the constants or the goldens

Command: `go test -count=1 ./internal/run/ -run 'TestEachStepOfASpecificationTaskStatesTheBetsCellsAndRecord'`
Start 2026-10-22T14:42:38Z, end 2026-10-22T14:42:40Z, exit 1 (the missing-text failures: the three summary sentences in each of the prompts and
constants of plan-of-a-task, plan-review and verification, and the two headings of the contract in loop.md; the specification control passes).

```
--- FAIL: TestEachStepOfASpecificationTaskStatesTheBetsCellsAndRecord (0.00s)
    instructions_test.go:657: the prompt of plan-of-a-task for a task of kind specification does not hold "The initial task of kind specification leaves the MoSCoW and Phase cells of the PRD blank (not —) until the first bet and writes no numbered facts record, because code supplies both at the rendered commit"
    instructions_test.go:660: the constant of plan-of-a-task does not hold "The initial task of kind specification leaves the MoSCoW and Phase cells of the PRD blank (not —) until the first bet and writes no numbered facts record, because code supplies both at the rendered commit"
    instructions_test.go:657: the prompt of plan-of-a-task for a task of kind specification does not hold "never ask that task to supply them, correct a plan that promises them, and treat cells it filled or a facts record it wrote as defects"
    instructions_test.go:660: the constant of plan-of-a-task does not hold "never ask that task to supply them, correct a plan that promises them, and treat cells it filled or a facts record it wrote as defects"
    instructions_test.go:657: the prompt of plan-of-a-task for a task of kind specification does not hold "A reopened task of kind specification keeps the cells, lines and record that already exist."
    instructions_test.go:660: the constant of plan-of-a-task does not hold "A reopened task of kind specification keeps the cells, lines and record that already exist."
    instructions_test.go:657: the prompt of plan-review for a task of kind specification does not hold "The initial task of kind specification leaves the MoSCoW and Phase cells of the PRD blank (not —) until the first bet and writes no numbered facts record, because code supplies both at the rendered commit"
    instructions_test.go:660: the constant of plan-review does not hold "The initial task of kind specification leaves the MoSCoW and Phase cells of the PRD blank (not —) until the first bet and writes no numbered facts record, because code supplies both at the rendered commit"
    instructions_test.go:657: the prompt of plan-review for a task of kind specification does not hold "never ask that task to supply them, correct a plan that promises them, and treat cells it filled or a facts record it wrote as defects"
    instructions_test.go:660: the constant of plan-review does not hold "never ask that task to supply them, correct a plan that promises them, and treat cells it filled or a facts record it wrote as defects"
    instructions_test.go:657: the prompt of plan-review for a task of kind specification does not hold "A reopened task of kind specification keeps the cells, lines and record that already exist."
    instructions_test.go:660: the constant of plan-review does not hold "A reopened task of kind specification keeps the cells, lines and record that already exist."
    instructions_test.go:657: the prompt of verification for a task of kind specification does not hold "The initial task of kind specification leaves the MoSCoW and Phase cells of the PRD blank (not —) until the first bet and writes no numbered facts record, because code supplies both at the rendered commit"
    instructions_test.go:660: the constant of verification does not hold "The initial task of kind specification leaves the MoSCoW and Phase cells of the PRD blank (not —) until the first bet and writes no numbered facts record, because code supplies both at the rendered commit"
    instructions_test.go:657: the prompt of verification for a task of kind specification does not hold "never ask that task to supply them, correct a plan that promises them, and treat cells it filled or a facts record it wrote as defects"
    instructions_test.go:660: the constant of verification does not hold "never ask that task to supply them, correct a plan that promises them, and treat cells it filled or a facts record it wrote as defects"
    instructions_test.go:657: the prompt of verification for a task of kind specification does not hold "A reopened task of kind specification keeps the cells, lines and record that already exist."
    instructions_test.go:660: the constant of verification does not hold "A reopened task of kind specification keeps the cells, lines and record that already exist."
    instructions_test.go:678: loop.md does not hold "The contract of the cells and the record"
    instructions_test.go:678: loop.md does not hold "The review scope of the contract"
FAIL
FAIL	github.com/pharzam/layup/internal/run	0.004s
FAIL
```

## Green

Command: `go test -count=1 ./internal/run/ -run 'TestEachStepOfASpecificationTaskStatesTheBetsCellsAndRecord'`
Start 2026-10-22T14:45:28Z, end 2026-10-22T14:45:28Z, exit 0: ok  	github.com/pharzam/layup/internal/run	0.003s

Command: `go test -count=1 -tags=integration ./internal/run/ -run 'Instruction|Prompt|Golden'`
Start 2026-10-22T14:45:28Z, end 2026-10-22T14:45:36Z, exit 0: ok  	github.com/pharzam/layup/internal/run	7.395s

## Round 1 (M1): the document check moves to integration level

Falsifier: `docs/spec/loop.md` moved away from the tree, `go test -count=1 ./internal/run/` still passes (the unit level reads no document); the file was put back.

Command: `go test -count=1 -tags=integration ./internal/run/ -run 'Instruction|Prompt|Golden'` (includes TestTheInstructionContractOfTheCellsAndTheRecordHasItsHomeInLoopMd). Start 2026-10-22T14:56:54Z, end 2026-10-22T14:57:04Z, exit 0.
