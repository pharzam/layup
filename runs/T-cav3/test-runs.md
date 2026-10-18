# T-cav3 test runs (red, then green)

Recorded at 2026-10-18T11:04Z from the clock after the step. Commands run in the worktree
`.worktree/T-cav3`, `-count=1`; the integration level with
`TMPDIR=/home/layup/.cache/layup-session/tmp-green`.

## Red: the compiling stubs, every named test

The stubs are `specInput`, `specRecords` and `specTree` returning the zero input and
`nil`, `isNumberedFacts` returning false and `psbRecordsOf` returning `nil, nil`. The
`filesAt` and `readRequired` extractions were already in (the commit of the stubs), with
the characterization tests of the tree items run green on them first (integration,
`-run 'TreeItems|SpecSections|RenderingEqual|Spec|Batch|Answers'`, 95.9 s, ok).

### Unit level (the final test file; the verdict of each test)

```
--- FAIL: TestEachRecordFieldOfTheFullInputIsItsRecord (0.00s)
--- FAIL: TestAMissingRecordOfTheFullInputIsAnError (0.00s)
--- FAIL: TestAFailedReadOfARecordIsAnErrorEvenWhereItsAbsenceIsNot (0.00s)
--- FAIL: TestARecordThatItsReaderRefusesIsAnError (0.00s)
--- PASS: TestNoInventoryIndexIsBeforeTheFirstBet (0.00s)
--- FAIL: TestEachVersionOfTheIndexIsItsExactBytes (0.00s)
--- FAIL: TestMergedIsTheTasksWithAMergedEvent (0.00s)
--- FAIL: TestTheRegisterIsItsBytesAndNotJudged (0.00s)
--- FAIL: TestAnAnswersCitationIsAQuestionIDNotARecordLine (0.00s)
--- FAIL: TestANonFullSHAIsRefusedBeforeAnyRead (0.00s)
--- FAIL: TestOnlyTheNumberedFactsFileNamesAProblemStatementRecord (0.00s)
FAIL
FAIL	github.com/pharzam/layup/internal/run	0.075s
FAIL
```

The lines of the first red run, with the same stub (the fixture of the questions was then
corrected, which changes no red here), show the reasons:

```
--- FAIL: TestEachRecordFieldOfTheFullInputIsItsRecord (0.00s)
    the input: PSB 0 bytes, Spans [], Questions 0, Answers [], Register 0 bytes, Merged []; want each its record
--- FAIL: TestAMissingRecordOfTheFullInputIsAnError (0.00s)
    start/problem-statement.md: <nil>; want a SourceError with the source "records:start/problem-statement.md"
    spec/spans.tsv: <nil>; want a SourceError with the source "records:spec/spans.tsv"
    questions.tsv: <nil>; want a SourceError with the source "records:questions.tsv"
    intake/answers.tsv: <nil>; want a SourceError with the source "records:intake/answers.tsv"
    tasks.tsv: <nil>; want a SourceError with the source "records:tasks.tsv"
--- FAIL: TestAFailedReadOfARecordIsAnErrorEvenWhereItsAbsenceIsNot (0.00s)
    start/problem-statement.md: <nil>; want a SourceError with the source "records:start/problem-statement.md"
    spec/spans.tsv: <nil>; want a SourceError with the source "records:spec/spans.tsv"
    tasks.tsv: <nil>; want a SourceError with the source "records:tasks.tsv"
    spec/inventory.tsv: <nil>; want a SourceError with the source "records:spec/inventory.tsv"
    spec/inventory-2.tsv: <nil>; want a SourceError with the source "records:spec/inventory-2.tsv"
    tasks/T-bbbb/events.tsv: <nil>; want a SourceError with the source "records:tasks/T-bbbb/events.tsv"
--- FAIL: TestARecordThatItsReaderRefusesIsAnError (0.00s)
    spec/spans.tsv: <nil>; want a SourceError with the source "records:spec/spans.tsv"
    questions.tsv: <nil>; want a SourceError with the source "records:questions.tsv"
    intake/answers.tsv: <nil>; want a SourceError with the source "records:intake/answers.tsv"
    tasks.tsv: <nil>; want a SourceError with the source "records:tasks.tsv"
    spec/inventory.tsv: <nil>; want a SourceError with the source "records:spec/inventory.tsv"
    tasks/T-aaaa/events.tsv: <nil>; want a SourceError with the source "records:tasks/T-aaaa/events.tsv"
--- FAIL: TestEachVersionOfTheIndexIsItsExactBytes (0.00s)
    Index [], Versions 0, <nil>; want two rows and the exact bytes of both versions
--- FAIL: TestMergedIsTheTasksWithAMergedEvent (0.00s)
    Merged [], <nil>; want [T-aaaa] (T-bbbb has only a session event, T-cccc no events file)
    Merged []; want [T-aaaa T-cccc] in the order of the register
--- FAIL: TestTheRegisterIsItsBytesAndNotJudged (0.00s)
    Register "", <nil>; want the bytes and no error
--- FAIL: TestAnAnswersCitationIsAQuestionIDNotARecordLine (0.00s)
    covers for Q-002: [covers fail did not run: spec/spans.tsv]; want pass
    covers for F-0001#2: [covers fail did not run: spec/spans.tsv]; want fail
    covers for F-0001-setup-answers#2: [covers fail did not run: spec/spans.tsv]; want fail
    the control, F-0001 named as the record of the problem statement: [covers fail did not run: spec/spans.tsv]; want covers pass on its fact 2
--- FAIL: TestANonFullSHAIsRefusedBeforeAnyRead (0.00s)
    base: <nil>; want a SourceError with the source "base"
    head: <nil>; want a SourceError with the source "head"
    upper case: <nil>; want a SourceError with the source "head"
--- FAIL: TestOnlyTheNumberedFactsFileNamesAProblemStatementRecord (0.00s)
    isNumberedFacts("F-0007-numbered-facts.md") = false; want true
    psbRecordsOf([F-0007-numbered-facts.md]) = [], <nil>; want [F-0007], error false
    psbRecordsOf([F-0007-numbered-facts.md F-0009-numbered-facts.md]) = [], <nil>; want [], error true
FAIL
FAIL	github.com/pharzam/layup/internal/run	0.043s
FAIL
```

Two tests are **green at the stub on purpose**, each beside its red control, and are not
counted as reds: `TestNoInventoryIndexIsBeforeTheFirstBet` (U4; its control is
`TestEachVersionOfTheIndexIsItsExactBytes`, which reads an index) and the clear-table case
of the demo (I2, below; its control is the judged table of the same test). The case `covers`
of `TestAnAnswersCitationIsAQuestionIDNotARecordLine` is red at the stub for the reason of
the review (condition 3): a zero input lacks `Spans`, so `covers` is
`did not run: spec/spans.tsv`; the controls (bare `Q-002` passes, `F-0001` named as a
record passes on its fact 2) make the case distinguish the wrong answers mapping, which is
killed by mutant M1.

### Integration level, the tree half

```
--- FAIL: TestTheTreeFieldsOfTheFullInputAreTheRegularFilesOfTheHead (0.28s)
--- FAIL: TestTheNumberedFactsFileAtTheHeadIsTheProblemStatementRecord (0.53s)
--- FAIL: TestAFailedReadOfTheTreeIsASourceErrorOfTheTarget (0.53s)
FAIL
FAIL	github.com/pharzam/layup/internal/run	1.365s
FAIL
```

Reasons (first red run): `SpecFiles []`, `PRDs []`, `Changed []` for a head with the
files; `PSBRecords []; want [F-0007]`; two files `<nil>; want a SourceError with the
source target:docs/facts@<sha>`; a head the clone lacks `<nil>; want a SourceError`;
and each of the three later blobs that cannot be read (after the earlier reads
succeed) `<nil>; want a SourceError`.

### Integration level, the demo (`TestTheFullInputOfAHeadAfterIntakeAndTheScaffold`)

```
--- FAIL: TestTheFullInputOfAHeadAfterIntakeAndTheScaffold (80.78s)
    Changed []; want [docs/facts/F-0007-numbered-facts.md docs/prd/PRD-0001-x.md docs/spec/a.md]
    SpecFiles 0, PRDs 0, PSBRecords [], AnswersRecord ""; want a.md, PRD-0001-x.md, [F-0007] and the empty answers record
    PSB 0 bytes, Spans 0, Questions 0, Answers 0; want each the record of the run
    Register 0 bytes, Merged [], Index [], Copy ""; want the register, [T-aaaa], none and none
    covers on the head: [covers clear no specification yet]; want pass (the PRD cites F-0007#1, the span fact)
    no copy: Copy "", Index [], Versions 0, <nil>; want the caller's copy, the index and the exact version
FAIL
FAIL	github.com/pharzam/layup/internal/run	80.788s
FAIL
exit 1
```

The case of the head that changes no specification (`base..base`) is green at the stub (the
zero input has no change, so the table is the clear one) and is the reason that the same
test also asserts the judged table of the head that changes the specification: its control.

## Green: the code

```
go test -count=1 ./internal/run/                     ok
go test -count=1 -tags=integration (tree half, 3 tests)   ok  1.0 s
go test -count=1 -tags=integration -run TheFullInputOfAHeadAfterIntake ./internal/run/   ok  56.3 s
```

The full gates (build, vet, gofmt, the whole unit and integration levels, the linters) are
in the section "Gates" below.
