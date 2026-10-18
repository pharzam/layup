# T-c6sx: the red run of every named test

The stubs of commits aaaf899 and 8a82d77 (inert declarations: the empty `stepInstructions`, the empty
`resultInstruction`, and `inputsWanted`, `checkStepKind`, `optionsKept`, `labelList` and `fill` that return zero
values) with the finished tests and the golden files. Each failing line is copied from the run and cut at 300
characters; a line that a table test repeats is shown twice at most.

## The unit level

Command, in the worktree: `go test -count=1 ./internal/run` (exit status 1).

```
--- FAIL: TestThePromptFileOfEachStepIsItsGoldenFile (0.00s)
    --- FAIL: TestThePromptFileOfEachStepIsItsGoldenFile/plan-of-a-task (0.00s)
        instructions_test.go:87: Text differs from testdata/prompts/plan-of-a-task.txt: got 1404 bytes, want 4649; first difference at byte 465: got "\nrequirements:\n- docs/prd/PRD-0001.md 93", want "- /s/S-aaaaaaaa/input/tests.tsv 11d5bbdb"
    --- FAIL: TestThePromptFileOfEachStepIsItsGoldenFile/plan-review (0.00s)
        instructions_test.go:87: Text differs from testdata/prompts/plan-review.txt: got 1401 bytes, want 4649; first difference at byte 462: got "\nrequirements:\n- docs/prd/PRD-0001.md 93", want "- /s/S-aaaaaaaa/input/tests.tsv 11d5bbdb"
    --- FAIL: TestThePromptFileOfEachStepIsItsGoldenFile/specification (0.00s)
    --- FAIL: TestThePromptFileOfEachStepIsItsGoldenFile/shape-architecture (0.00s)
    --- FAIL: TestThePromptFileOfEachStepIsItsGoldenFile/implement (0.00s)
    --- FAIL: TestThePromptFileOfEachStepIsItsGoldenFile/verification (0.00s)
    --- FAIL: TestThePromptFileOfEachStepIsItsGoldenFile/close-out (0.00s)
--- FAIL: TestTheVerifierIsToldTheHeadThatNothingMaterialReads (0.00s)
    instructions_test.go:110: the prompt of the verifier does not hold "\nhead: 0123456789abcdef0123456789abcdef01234567\n"
    instructions_test.go:110: the prompt of the verifier does not hold " from fedcba9876543210fedcba9876543210fedcba98 to 0123456789abcdef0123456789abcdef01234567,"
--- FAIL: TestTheLabelsOfTheResultInstructionAreTheOwnerMaps (0.00s)
    instructions_test.go:117: labelList = ""; want the sorted keys, joined by a comma and a space
    instructions_test.go:120: labelList(DefaultOwners()) = ""
    instructions_test.go:123: the result instruction does not name the labels of the owner map:
--- FAIL: TestTheFormShownInEachInstructionParses (0.00s)
    instructions_test.go:166: the form of the plan: []; want its four fields, tests: and one bullet
--- FAIL: TestTheResultRowsThatTheInstructionQuotesAreReadByTheResultReader (0.00s)
    instructions_test.go:231: plan-of-a-task: the result that follows the instruction: line 1: an empty line ([])
    instructions_test.go:236: plan-of-a-task: a stopping question that follows the instruction: line 1: an empty line, []
    instructions_test.go:231: plan-review: the result that follows the instruction: line 1: an empty line ([])
    instructions_test.go:236: plan-review: a stopping question that follows the instruction: line 1: an empty line, []
--- FAIL: TestTheStepAndTheKindOfTheTaskAgree (0.00s)
    instructions_test.go:253: checkStepKind(specification, architecture) = <nil>; want ok false
    instructions_test.go:253: checkStepKind(specification, build) = <nil>; want ok false
--- FAIL: TestAShapeTaskReceivesTheRevisedDraftAndItsFactsAndAnArchitectureTaskItsOptions (0.00s)
    instructions_test.go:268: inputsWanted(plan-of-a-task, specification) = false, false, false; want true, false, false
    instructions_test.go:268: inputsWanted(plan-of-a-task, architecture) = false, false, false; want true, true, false
    instructions_test.go:275: plan-of-a-task of a specification task: requirements = [docs/prd/PRD-0001.md docs/spec/loop.md]
    instructions_test.go:275: plan-review of a architecture task: requirements = [docs/prd/PRD-0001.md docs/spec/loop.md]
--- FAIL: TestTheNewInputsAreInTheirKindsAfterTheExistingItems (0.00s)
    instructions_test.go:295: requirements = [docs/prd/PRD-0001.md docs/spec/loop.md]; want the four Intake records last, in order
    instructions_test.go:298: decision-records = [docs/adr/0019-phase-loop.md docs/adr/0020-route.md docs/adr/0021-smart-if.md]; want the ADRs then the options table
    instructions_test.go:301: the survivors are not the text under the decision records
    instructions_test.go:304: issue-and-plan = [tasks.tsv copies/701-1.md]; want the frozen list after the plan
    instructions_test.go:307: handoffs-and-answers = [copies/702-1.md tasks/T-ab12/results/S-bbbbbbbb.tsv]; want the ledger after the answers
--- FAIL: TestTwoNewInputsOfOneBaseNameAreRefused (0.00s)
    instructions_test.go:328: an options table named draft.tsv: <nil>; want the refusal naming both
    instructions_test.go:328: an options table named spans.tsv: <nil>; want the refusal naming both
    instructions_test.go:334: a handoff named sessions.tsv beside the ledger: <nil>
--- FAIL: TestTheResultsOfTheFilterAreCheckedAgainstTheirTables (0.00s)
    instructions_test.go:362: a result that agrees with its table: [], [], <nil>
--- FAIL: TestTheSpecimenOfTheSpecificationInstructionPassesTheFullCheck (0.00s)
    instructions_test.go:439: the specification instruction shows 0 forms; want the requirements, the criteria, the success metrics and the matrix
exit 1
```

Green at the stub, beside their red controls (not red runs):

- `TestNoPlaceholderIsLeftInAnyPrompt`: no text at the stub, so no placeholder; its red controls are the golden tests.
- The existing prompt tests (`TestThePromptFileListsTheFiveKindsInTheirOrder` and the others of `promptfile_test.go`).

## The integration level

Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green timeout 900 go test -count=1 -tags=integration -timeout 20m ./internal/run -run 'TestEachStepsPromptFileFromARealRecordsTree|TestPromptFileRefusalsHaveTheirControls|TestAHarnessThatFollowsTheInstructionGivesValidHandoffs|TestTheInstructionBlocksOfLoopMd|TestTheSpecimenOfTheSpecificationInstructionIsAcceptedByThePinnedPrdLint'` (exit status 1).

```
--- FAIL: TestEachStepsPromptFileFromARealRecordsTree (19.17s)
    promptfile_integration_test.go:308: the prompt of plan-of-a-task holds no line "The work of this step:" followed by its instruction
    promptfile_integration_test.go:316: plan-of-a-task of T-sp01: the draft is in the prompt: false
    promptfile_integration_test.go:308: the prompt of plan-review holds no line "The work of this step:" followed by its instruction
    promptfile_integration_test.go:316: plan-review of T-ar01: the draft is in the prompt: false
--- FAIL: TestPromptFileRefusalsHaveTheirControls (12.17s)
    promptfile_integration_test.go:346: the step shape-activation-batch: git ls-tree -r -z --full-tree --end-of-options 0000000000000000000000000000000000000000 -- .: exit status 128: fatal: not a tree object
    promptfile_integration_test.go:346: the step plan-of-a-milestone: git ls-tree -r -z --full-tree --end-of-options 0000000000000000000000000000000000000000 -- .: exit status 128: fatal: not a tree object
    promptfile_integration_test.go:354: specification of a build task: <nil>; want an error holding "the step specification is not a step of a task of kind build"
    promptfile_integration_test.go:355: implement of a specification task: <nil>; want an error holding "the step implement is not a step of a task of kind specification"
    promptfile_integration_test.go:357: no options for the architecture: <nil>; want an error holding "the options of the task T-ar01"
    promptfile_integration_test.go:358: options for a build task: <nil>; want an error holding "the options of the task T-ab12"
    promptfile_integration_test.go:359: options for the close-out: <nil>; want an error holding "the options of the task T-ar01"
    promptfile_integration_test.go:362: the missing spec/draft.tsv: <nil>; want an error holding "the records have no spec/draft.tsv"
    promptfile_integration_test.go:362: the missing spec/spans.tsv: <nil>; want an error holding "the records have no spec/spans.tsv"
    promptfile_integration_test.go:367: a draft with a bad kind: <nil>; want an error holding "spec/draft.tsv: "
    promptfile_integration_test.go:370: a missing table: <nil>; want an error holding "the records have no spec/options/db.tsv"
    promptfile_integration_test.go:373: a renamed option: <nil>; want an error holding "decision db: the options of the result are"
    promptfile_integration_test.go:374: a decision named like a record: <nil>; want an error holding "two records are named draft.tsv"
    promptfile_integration_test.go:377: a malformed frozen list: <nil>; want an error holding "tasks/T-sp01/tests.tsv: "
    promptfile_integration_test.go:381: a good frozen list: <nil>, &{[{issue-and-plan records tasks.tsv 807578c491a21195a4f674cfd3186a4805d459a978196bf94e24841ef5d3561b} {requirements target docs/prd/PRD-0001.md 218ab42144fb9be6d46440652e37bd89f0ebd134bd2ad7a1df8174f943dc1db3} {requirements target d
    promptfile_integration_test.go:384: verification with no attempt: <nil>; want an error holding "the task T-sp01 has no attempt"
    promptfile_integration_test.go:392: close-out with no ledger: <nil>; want an error holding "the records have no sessions.tsv"
    promptfile_integration_test.go:401: a refused telemetry table: <nil>; want an error holding "telemetry.tsv: "
--- FAIL: TestAHarnessThatFollowsTheInstructionGivesValidHandoffs (15.73s)
    --- FAIL: TestAHarnessThatFollowsTheInstructionGivesValidHandoffs/plan-of-a-task (6.25s)
        promptfile_integration_test.go:426: the prompt of plan-of-a-task does not name result/plan.md
        promptfile_integration_test.go:440: PostComment: the records have no kept file tasks/T-ab12/results/S-9d7a0c60/plan.md
    --- FAIL: TestAHarnessThatFollowsTheInstructionGivesValidHandoffs/plan-review (5.21s)
        promptfile_integration_test.go:426: the prompt of plan-review does not name result/plan-review.md
        promptfile_integration_test.go:440: PostComment: the records have no kept file tasks/T-ab12/results/S-0bb5da8e/plan-review.md
    --- FAIL: TestAHarnessThatFollowsTheInstructionGivesValidHandoffs/verification (4.26s)
--- FAIL: TestTheInstructionBlocksOfLoopMd (0.00s)
    promptfile_integration_test.go:465: 1 texts in the code; want the seven steps and the shared result text
    promptfile_integration_test.go:471: loop.md has the block instruction-result: true; it is equal to the code: false
--- FAIL: TestTheSpecimenOfTheSpecificationInstructionIsAcceptedByThePinnedPrdLint (0.27s)
    promptfile_integration_test.go:495: the specimen: accepted false, [FAIL  PRD-0001-x.md: no REQ/NFR requirement rows found (an unfilled PRD skeleton?)], <nil>; want true
FAIL
FAIL	github.com/pharzam/layup/internal/run	47.485s
FAIL
exit 1
```

Green controls inside `TestPromptFileRefusalsHaveTheirControls` at the stub: the stage-0 control (`implement` with an
unresolvable base fails by the later `git ls-tree`, not by "has no instruction", which is what makes the early
refusal observable: at the stub the three refused steps fail by that same `git ls-tree` line), and `plan-review` of a
build task.

## The gates

Head `39e02ef` (the code is that of `b98ccc9`), the host clock in UTC; every exit status read directly.

| Command, as run | Exit | Time |
| --- | --- | --- |
| `go build ./...` | 0 | 13:22:30Z |
| `go vet ./...` | 0 | 13:22:30Z |
| `gofmt -l .` | no file | 13:22:30Z |
| `go vet -tags=e2e ./cmd/layup/` | 0 | 13:22:30Z |
| `go vet -tags=integration ./...` | 0 | 13:22:30Z |
| `go test -count=1 ./...` | 0 | about 1.5 min |
| `sh docs/adr/adr-lint.sh`, `sh docs/prd/prd-lint.sh`, `sh docs/links/link-lint.sh`, `sh docs/tests/run-discipline-tests.sh`, `sh docs/setup/setup-check.sh` (after `git add`), `sh docs/tests/nested-checkout-check.sh`, `git diff --check origin/main HEAD` | 0 each | about 13:25Z |
| `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...` | 1: every package `ok`, and `internal/run` timed out at 30m on a loaded host (load average 5 to 10, other sessions' tests) | 13:22:42Z to 13:54:20Z |
| the same, `./internal/run` alone, `-timeout 45m -v` | 1: `TestTheAnswersToTheFollowUpWriteTheBudget` (51.8 s) and `TestAnUnparsedAnswerToTheFollowUpLeavesItsQuestionOpen` (60.2 s, then a panic) failed under load, tests of Intake that this task does not touch; the two pass alone: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 20m ./internal/run -run 'TestTheAnswersToTheFollowUpWriteTheBudget\|TestAnUnparsedAnswerToTheFollowUpLeavesItsQuestionOpen'` exit 0, 59.6 s | 13:55Z to 14:16:15Z |
| the same, `./internal/run` alone, `-timeout 60m -v` | 0, 1498.5 s, no test failed | 14:17:48Z to 14:42:48Z |
