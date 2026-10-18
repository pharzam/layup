# T-cav3 — the input of the full spec check at a head

Issue [#412](https://github.com/pharzam/layup/issues/412), row 66p, a size child of row 66f
([#362](https://github.com/pharzam/layup/issues/362)). Serves REQ-012, REQ-007, NFR-004 (F-0003#47,
F-0003#51). Plan 6248854300 (Opus 5.5); plan review 6248921962 (GPT-6.1 Sol, six conditions, all met).
Maximum 1,100 lines over 17 files, cap 1. Evidence: [test runs](../../runs/T-cav3/test-runs.md),
[mutants](../../runs/T-cav3/mutants.md).

## Delivered and boundaries

`(*Sessions).specInput(ctx, base, head, copy)` builds the `spec.FullInput` of `spec.CheckFull` for a
head, each field from one named source (the table is in `loop.md`, "The full spec check", decided before
code): the records at the run's base at the call, the tree at the head in the clone, and the caller's copy.
A missing required record, a refused record, a failed read, a failed `git` call or a non-full SHA is a
`SourceError` with the zero input and no write; no `spec/inventory.tsv` is "before the first bet".
`AnswersRecord` is `""` (the setup's answers record numbers `S01-stack` first, so its fact 2 is not `Q-002`).
`PSBRecords` is the one file `docs/facts/F-NNNN-numbered-facts.md` at the head, which row 66j produces at the
first bet (none: empty, one: its ID, two: a `target:` error). Not delivered: the post of `layup/spec` (66f),
the bet's copy (66j), the step loop (66k).

## Tests mapped to the criteria

| Criterion | Test |
| --- | --- |
| Each field from its record | `TestEachRecordFieldOfTheFullInputIsItsRecord`, `TestMergedIsTheTasksWithAMergedEvent`, `TestTheRegisterIsItsBytesAndNotJudged` |
| A missing, failed or refused record, or a bad SHA, is an error with the zero input | `TestAMissingRecordOfTheFullInputIsAnError`, `TestAFailedReadOfARecordIsAnErrorEvenWhereItsAbsenceIsNot`, `TestARecordThatItsReaderRefusesIsAnError`, `TestANonFullSHAIsRefusedBeforeAnyRead` |
| Before the first bet; the index versions; the answers mapping | `TestNoInventoryIndexIsBeforeTheFirstBet`, `TestEachVersionOfTheIndexIsItsExactBytes`, `TestAnAnswersCitationIsAQuestionIDNotARecordLine` |
| The numbered facts file and the tree fields | `TestOnlyTheNumberedFactsFileNamesAProblemStatementRecord`, `TestTheNumberedFactsFileAtTheHeadIsTheProblemStatementRecord`, `TestTheTreeFieldsOfTheFullInputAreTheRegularFilesOfTheHead`, `TestAFailedReadOfTheTreeIsASourceErrorOfTheTarget` |
| The extractions; the demo and the copy | `TestReadRequiredKeepsTheStoreErrorTheMissingDiagnosticAndTheReaderPath`, `TestTheFullInputOfAHeadAfterIntakeAndTheScaffold` |

## Alternatives

`PSBRecords`: the named file at the head, empty until 66j writes it; rejected: every `F-NNNN-*.md` (it counts
a baseline's own records) and always empty (a permanent limit). `AnswersRecord`: `""`; rejected: the S04
record, whose fact N is not `Q-00N`, so a wrong citation would pass.

## Verification and verdict

Own test-first tests (stubs committed alone, red recorded), the integration demo on a target after Intake
and the Scaffold, and 20 compiled mutants, all killed (M18 and M19 survived the first run and led to a
test; statuses and hashes per mutant are in `mutants.md`). The whole integration level passed (exit 0,
11:27:57Z to 11:55:39Z); build, vet, unit and the text checks pass. Final exact-head CI remains a landing
condition.

Round 1 (GPT-6.1 Sol, xhigh, on Codex), comment 6249722495: **material**, two findings: the PRD selector's
regex missed a name that holds a line feed (fixed test first), and the mutant record lacked per-mutant
evidence (rerun). Round 2, comment 6249923251: one finding, the integration command recorded unquoted
(rerun quoted). The Operator chose end 3 (a fix with no further round) on Telegram at 2026-10-18T12:29:56Z
(comment 6249932698); the fix of the record is 1f456df. Harness exception: ADR-0026 names `claude`, `devin`
and `opencode`; Codex was used on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens are as the harness reported them.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 | not recorded | not separately recorded | not separately recorded |
| Plan review | reasoning | GPT-6.1 Sol | xhigh | 185,969 | 2026-10-18T10:40:01Z to 10:47:16Z |
| Implementation | execution | Claude Sonnet 5.5 | not recorded | 240,517 | 69 min |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 193,398 | 11:57:40Z to 12:06:47Z |
| Fix of round 1 | execution | Claude Sonnet 5.5 | not recorded | not recorded | about 15 min |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 149,061 | 12:22:26Z to 12:28:32Z |
| Fix of round 2 | execution | Claude Sonnet 5.5 | not recorded | not recorded | about 5 min |
| Close-out | execution | Claude Sonnet 5.5 | not recorded | not recorded | about 5 min, ending 12:32Z |
| Orchestration | reasoning | Claude Opus 5.5 | not recorded | not separately recorded | not separately recorded |
