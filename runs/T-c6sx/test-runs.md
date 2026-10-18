# T-c6sx: the red run of every named test

The stubs of commits aaaf899 and 8a82d77 (empty `stepInstructions` and `resultInstruction`; `inputsWanted`,
`checkStepKind`, `optionsKept`, `labelList` and `fill` returning zero values) with the finished tests and the golden
files. Each line is a failing test and its first assertion (cut at 230 characters), copied from the run; a subtest is
not repeated. Green at the stub, beside their red controls: `TestNoPlaceholderIsLeftInAnyPrompt`, the existing prompt tests,
the stage-0 control and the `plan-review` control of `TestPromptFileRefusalsHaveTheirControls` (at the stub the three
refused steps fail by the later `git ls-tree` line, which is what makes the early refusal observable).

Unit, `go test -count=1 ./internal/run` (exit 1):

```
TestThePromptFileOfEachStepIsItsGoldenFile
TestTheVerifierIsToldTheHeadThatNothingMaterialReads -> instructions_test.go:110: the prompt of the verifier does not hold "\nhead: 0123456789abcdef0123456789abcdef01234567\n"
TestTheLabelsOfTheResultInstructionAreTheOwnerMaps -> instructions_test.go:117: labelList = ""; want the sorted keys, joined by a comma and a space
TestTheFormShownInEachInstructionParses -> instructions_test.go:166: the form of the plan: []; want its four fields, tests: and one bullet
TestTheResultRowsThatTheInstructionQuotesAreReadByTheResultReader -> instructions_test.go:231: plan-of-a-task: the result that follows the instruction: line 1: an empty line ([])
TestTheStepAndTheKindOfTheTaskAgree -> instructions_test.go:253: checkStepKind(specification, architecture) = <nil>; want ok false
TestAShapeTaskReceivesTheRevisedDraftAndItsFactsAndAnArchitectureTaskItsOptions -> instructions_test.go:268: inputsWanted(plan-of-a-task, specification) = false, false, false; want true, false, false
TestTheNewInputsAreInTheirKindsAfterTheExistingItems -> instructions_test.go:295: requirements = [docs/prd/PRD-0001.md docs/spec/loop.md]; want the four Intake records last, in order
TestTwoNewInputsOfOneBaseNameAreRefused -> instructions_test.go:328: an options table named draft.tsv: <nil>; want the refusal naming both
TestTheResultsOfTheFilterAreCheckedAgainstTheirTables -> instructions_test.go:362: a result that agrees with its table: [], [], <nil>
TestTheSpecimenOfTheSpecificationInstructionPassesTheFullCheck -> instructions_test.go:439: the specification instruction shows 0 forms; want the requirements, the criteria, the success metrics and the matrix
exit 1
```

Integration, `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 20m ./internal/run -run 'TestEachStepsPromptFileFromARealRecordsTree|TestPromptFileRefusalsHaveTheirControls|TestAHarnessThatFollowsTheInstructionGivesValidHandoffs|TestTheInstructionBlocksOfLoopMd|TestTheSpecimenOfTheSpecificationInstructionIsAcceptedByThePinnedPrdLint'` (exit 1):

```
TestEachStepsPromptFileFromARealRecordsTree -> promptfile_integration_test.go:308: the prompt of plan-of-a-task holds no line "The work of this step:" followed by its instruction
TestPromptFileRefusalsHaveTheirControls -> promptfile_integration_test.go:346: the step shape-activation-batch: git ls-tree -r -z --full-tree --end-of-options 0000000000000000000000000000000000000000 -- .: exit status 128: fatal: not a tree object
TestAHarnessThatFollowsTheInstructionGivesValidHandoffs
TestTheInstructionBlocksOfLoopMd -> promptfile_integration_test.go:465: 1 texts in the code; want the seven steps and the shared result text
TestTheSpecimenOfTheSpecificationInstructionIsAcceptedByThePinnedPrdLint -> promptfile_integration_test.go:495: the specimen: accepted false, [FAIL  PRD-0001-x.md: no REQ/NFR requirement rows found (an unfilled PRD skeleton?)], <nil>; want true
exit 1
```

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

## Round 1 (comment 6251434844): findings 1 to 3, red before code

Commands as run, on the head `cbb40b0` plus the new tests (for finding 1 also an inert field `FullInput.Records` and three
inert functions). Finding 3's new cases pass on the code that was right: their red is the mutants M23, M24, M26, M27 and
M28 of `mutants.md`, which pass every test of `cbb40b0`.

Finding 2, `go test -count=1 ./internal/run`, then `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 20m ./internal/run -run 'TestAHarnessThatFollowsTheInstructionGivesValidResultsForTheOtherStepsAndAStoppingQuestion|TestPromptFileRefusalsHaveTheirControls'`:

```
TestTheResultNamesEveryRequiredOutputWithItsHash -> instructions_test.go:515: close-out: the required outputs of the text are []; want [docs/tasks/T-ab12.md docs/tasks/completed.md]
exit 1
```
```
TestAHarnessThatFollowsTheInstructionGivesValidResultsForTheOtherStepsAndAStoppingQuestion
exit 1
```

Finding 1, `go test -count=1 ./internal/spec ./internal/run`, then `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 20m ./internal/run -run 'TestTheTreeFieldsOfTheFullInputAreTheRegularFilesOfTheHead|TestTheNumberedFactsFileAtTheHeadIsTheProblemStatementRecord|TestTheSpecimenPassesTheFullCheckAndThePrdLintThroughTheProductionInput'`:

```
TestCoversBindsEachTokenToItsRecord -> full_test.go:324: the rule covers: result "fail", want "pass" (detail "the Facts cell of REQ-001 cites F-0005, which is no fact of the inventory or the answers")
TestTheSpecimenOfTheSpecificationInstructionPassesTheFullCheck -> instructions_test.go:481: the rule covers fails on the specimen: the Facts cell of REQ-001 cites F-0003, which is no fact of the inventory or the answers
TestTheFactsRecordsAreThoseOfPrdLintAndTheNextIDIsOneMore -> instructions_test.go:503: factRecordIDs = []; want [F-0001 F-0002 F-0003 F-0004]
exit 1
```
```
TestTheTreeFieldsOfTheFullInputAreTheRegularFilesOfTheHead -> specinput_integration_test.go:61: PSBRecords [], Records []; want the predicted ID F-0004 (no numbered facts file) and the setup record F-0003
TestTheNumberedFactsFileAtTheHeadIsTheProblemStatementRecord -> specinput_integration_test.go:93: none: PSBRecords [], <nil>; want [F-0001]
TestTheSpecimenPassesTheFullCheckAndThePrdLintThroughTheProductionInput -> specinput_integration_test.go:270: the rule covers fails on the specimen through the production input: the Facts cell of REQ-001 cites F-0004#1, which is no fact of the inventory or the answers
exit 1
```

Gates of the fixes of round 1, head `640d0c6`: `go build ./...`, `go vet ./...`, `go vet -tags=integration ./...`, `go vet -tags=e2e ./cmd/layup/`, `gofmt -l .` (no file), `go test -count=1 ./...`, the five text checks, `nested-checkout-check` and `git diff --check origin/main HEAD`: exit 0 each. `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 45m ./...`: exit 0, 15:30:12Z to 15:50:51Z, every package `ok`.

Finding 1 of round 2 (the prompt and the full input predict different facts records for an executable or a linked `F-NNNN` entry), red before the fix, head `b4c3fb4` plus the extraction `targetTreeAt` and the test, started 2026-10-18T17:11:42Z: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 20m ./internal/run -run 'TestThePromptAndTheFullInputPredictTheSameFactsRecords'`:

```
an_executable_file -> specinput_integration_test.go:314: the prompt: record F-0008, next F-0009; want F-0003 and F-0004
a_link -> specinput_integration_test.go:314: the prompt: record F-0008, next F-0009; want F-0003 and F-0004
(also line 321: the full input: PSBRecords [F-0004], Records [F-0003]; and line 327: the rule covers fails: the Facts cell of REQ-001 cites F-0009#1, which is no fact of the inventory or the answers)
the control (regular files only) and the gitlink pass (a gitlink is a commit entry, not a blob, to both)
exit 1
```
