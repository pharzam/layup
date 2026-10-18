# T-c6sx — the instructions and the inputs of the task steps

Issue [#416](https://github.com/pharzam/layup/issues/416). Serves REQ-005/REQ-012/F-0003#45/F-0003#51.
Row 66s of the plan of M2e, a child of row 66h (#364) found by Q15 and condition 4 of its plan review
(comment 6249724875). Author: Sonnet 5.5 (execution tier; the plan is Opus 5.5's).
Plan comment 6249846962; plan review 6249930540 (GPT-6.1 Sol, approve-with-conditions; conditions 2 to
7 met, condition 1 settled by the Operator at goal count 2, comment 6249938546, the base guard moved to #415).
Maximum 2,100 changed lines/27 files against main at landing, cycle cap 1, no panel.
Evidence: [runs](../../runs/T-c6sx/test-runs.md) (red), [mutants](../../runs/T-c6sx/mutants.md).

## Delivered and boundaries

For each of the seven task steps (`plan-of-a-task`, `plan-review`, `specification`, `shape-architecture`,
`implement`, `verification`, `close-out`) the prompt file of `PromptFile` is the existing ordered, hashed list
of the five kinds and, after it, the step's fixed instruction (`The work of this step:`), then the shared
result instruction. The eight texts are blocks of `loop.md` ("The instruction of a task step"), checked byte
for byte against the constants of `internal/run/instructions.go`. New inputs, inside the five kinds: the frozen
list (`issue-and-plan`), the revised draft and its facts (`requirements`), the options with the survivors
(`decision-records`) and the ledger (`handoffs-and-answers`). `shape-activation-batch`, `plan-of-a-milestone` and
`answer-to-a-question` are refused by name before every read. The step loop, the driver, the filter and the base
guard of `TaskSession` are rows 66k, 66r and 66h.

## Tests mapped to the criteria

| Criterion | Test |
| --- | --- |
| The bytes of each step's prompt (work, outputs, result file, form, inputs) | `TestThePromptFileOfEachStepIsItsGoldenFile`, `TestNoPlaceholderIsLeftInAnyPrompt` |
| The shown forms are read by the real parsers (condition 2) | `TestTheFormShownInEachInstructionParses`, `TestTheResultRowsThatTheInstructionQuotesAreReadByTheResultReader`, `TestAHarnessThatFollowsTheInstructionGivesValidHandoffs` |
| The specification output passes the existing checks (condition 3) | `TestTheSpecimenOfTheSpecificationInstructionPassesTheFullCheck`, `TestTheSpecimenOfTheSpecificationInstructionIsAcceptedByThePinnedPrdLint` |
| The options contract and the collisions (condition 4) | `TestTheResultsOfTheFilterAreCheckedAgainstTheirTables`, `TestTwoNewInputsOfOneBaseNameAreRefused`, `TestAShapeTaskReceivesTheRevisedDraftAndItsFactsAndAnArchitectureTaskItsOptions` |
| The refusals with their controls, the early stage 0 (condition 6) | `TestPromptFileRefusalsHaveTheirControls`, `TestTheStepAndTheKindOfTheTaskAgree` |
| One home for the texts | `TestTheInstructionBlocksOfLoopMd` |
| A real records tree for each step (the demo) | `TestEachStepsPromptFileFromARealRecordsTree` |
| The base of the last attempt | `TestTheChangeOfTheVerifierStartsAtTheBaseOfTheLastAttempt` |

## Alternatives

| Decision | Selected | Rejected, and why |
| --- | --- | --- |
| Home of the texts | blocks of `loop.md`, constants checked by a test | `go:embed` of `docs/` (the package cannot reach it); texts in code only (an unread second copy) |
| Position of the instruction | after the lists | before them (the order test of the existing lists, and a reader reads the links first) |
| Kinds of the new inputs | the five kinds | a sixth kind (changes `ContextKinds` and `session.md`) |
| The filtered options | a map of `spec.Result` from the caller, checked against each table | computing the constraints in `PromptFile` (needs row 66h's filter) |
| The survivors after the bet | the same options and the chosen option as the only survivor | a second input |
| Facts citations | `F-{facts}#n` and the existing `F-{record}` | omitting citations (refused by the pinned `prd-lint`) |
| Oracle of the bytes | embedded golden files, no update flag | a flag that rewrites them |

## Acceptance

- [x] Each of the seven task steps has a prompt file from a real records tree: the ordered, hashed list, the step's work, its outputs, the result file and the form (the golden files and the real-tree demo).
- [x] The shown forms are read by the real parsers; a scripted harness that follows the text writes valid handoffs.
- [x] The specification output passes the full check and the pinned `prd-lint` through the production input (Contract A: bare existing records, the predicted ID before the bet), also for an executable or a linked `F-NNNN` entry (finding 1 of round 2).
- [x] Required inputs and outputs are refused with controls; the three steps outside M2e are refused by name; the instruction texts are the blocks of `loop.md`, byte for byte.

## Verification and verdict

Own test-first tests (stubs committed alone, the red of every named test recorded in
`runs/T-c6sx/test-runs.md`, tests committed with the code), and 36 compiled mutants (M1 to M37 less M13; all that built were
killed; M28 survives and is equivalent; the first M7 did not build and the first M22 was a faulty mutant, both rerun;
the plan's M13 moved to #415) cover the texts, the inputs, the refusals and the wiring.

Rounds. Plan review (GPT-6.1 Sol): approve-with-conditions. Round 1 (comment 6251434844, GPT-6.1 Sol, xhigh,
Codex): **material**: the specification head could not pass both `covers` and `prd-lint` before the bet (decided as
Contract A, comment 6251477575, an R10 fix in this branch), a required output of the close-out was not a hashed
artifact of the result, and the refusal table lacked controls for a missing brief and ledger; all fixed test first.
Round 2 (comment 6252675475, same reviewer): **not mergeable, findings recorded**: finding 1 (the prompt and the full
check predicted different facts records for an executable or a linked entry) and finding 2 (the budget). Harness
exception: ADR-0026 names `claude`, `devin` and `opencode`; Codex was used on the Operator's message of 2026-10-13.
The cap was reached; the Operator chose **end 3** (comment 6252711795, 2026-10-18T17:08:30Z): fix finding 1 with no
further round, then close out. Round 2's verdict stands as the last round's record.

The fix of finding 1 is **not reviewed**, by the Operator's end 3. It is test first: the extraction `targetTreeAt`
(a refactor, green), `TestThePromptAndTheFullInputPredictTheSameFactsRecords` red for the executable file and the link
(the control and the gitlink pass), then `factsNames`, the one home of the selection (regular files, mode 100644,
directly in `docs/facts/`), which `factsRecords` (the prompt) and `specTree` (the full input) both call; green, with
two mutants (M36, M37) killed. Finding 2 is the Operator's single approval of the landing figure, below.

Gates after the merge of `origin/main` (1945bc6, no conflict): `go build ./...`, `go vet ./...`, `go vet -tags=integration ./...`,
`go vet -tags=e2e ./cmd/layup/`, `gofmt -l .` (no file), `go test -count=1 ./...`, the five text checks,
`nested-checkout-check` and `git diff --check origin/main HEAD`: exit 0 each. The whole integration level
(`TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 45m ./...`) on the merge
head `0d8a715`: exit 0, every package `ok` (`internal/run` 1967 s), 2026-10-18T17:18:04Z to 17:52:15Z. The one commit after
it changes the plan row 66s only (its After cell is `66l, 66p`, which `task-state.sh` reads; the words about 66h's R10
decisions moved to the Task cell); the text checks passed on it. The exact-head CI is the Operator's merge gate.

Landing figure before this task file and the completed line: 2,229 lines (2,187 added, 42 removed) over 26 files, against
a maximum of 2,100/27; the merge waits for the Operator's one approval of the exact figure.

## Resource record

ADR-0007, recorded, not budgeted.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 | not recorded | not separately recorded | not separately recorded |
| Plan review | reasoning | GPT-6.1 Sol | xhigh | 242,572 | 2026-10-18T12:20:58Z to 12:29:24Z |
| Implementation and fixes of round 1 | execution | Claude Sonnet 5.5 | not recorded | not recorded | not recorded |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 378,334 | 2026-10-18T14:44:59Z to 14:55:26Z |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 186,152 | 2026-10-18T15:52:16Z to 16:02:08Z |
| Fix of finding 1, merge, gates and close-out | execution | Claude Sonnet 5.5 | not recorded | not recorded | from about 17:10Z to the close-out commit; the integration run 17:18:04Z to 17:52:15Z |
| Orchestration | reasoning | Claude Opus 5.5 | not recorded | not separately recorded | not separately recorded |
