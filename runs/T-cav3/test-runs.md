# T-cav3 test runs (red, then green, then the gates)

Recorded from the clock (2026-10-18T11:28Z for this file's rewrite; the figures are those of the runs).
Commands run in `.worktree/T-cav3` with `-count=1`; the integration level with
`TMPDIR=/home/layup/.cache/layup-session/tmp-green`.

## Step S2: the extractions (stubs committed alone)

`filesAt` and `readRequired` were extracted with `specFilesAt` and `state.readRecord` as
one-line wrappers, with the compiling stubs of the builder. Characterization on that commit:
`go test ./...` ok, and the integration tests of the tree items and of the batch and the answers
(`-run 'TreeItems|SpecSections|RenderingEqual|Spec|Batch|Answers'`, 95.9 s) ok.

## Red: the stubs, every named test

The stubs return the zero input and `nil` (`isNumberedFacts` false, `psbRecordsOf` `nil, nil`).

Unit (`internal/run/specinput_test.go`), verdict at the stub:

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
```

Integration, tree half and demo:

```
--- FAIL: TestTheTreeFieldsOfTheFullInputAreTheRegularFilesOfTheHead (0.28s)
--- FAIL: TestTheNumberedFactsFileAtTheHeadIsTheProblemStatementRecord (0.53s)
--- FAIL: TestAFailedReadOfTheTreeIsASourceErrorOfTheTarget (0.53s)
--- FAIL: TestTheFullInputOfAHeadAfterIntakeAndTheScaffold (80.78s)
    Changed []; SpecFiles 0, PRDs 0, PSBRecords []; PSB 0 bytes, Spans 0; Register 0 bytes, Merged [];
    covers on the head: [covers clear no specification yet]; no copy: Index [], Versions 0
```

Reasons: every field is empty where its record has rows; a missing, failed or refused record
gives `<nil>; want a SourceError with the source records:<path>` (five required paths, six
failed reads, six refused records); a non-full SHA gives `<nil>`; `isNumberedFacts` is false for
`F-0007-numbered-facts.md`; a head with files gives `SpecFiles []`, `PRDs []`, `Changed []`; two numbered
facts files give `<nil>; want a SourceError target:docs/facts@<sha>`; a head the clone lacks and
each of three later unreadable blobs (after the earlier reads succeeded) give `<nil>`.

**Green at the stub, on purpose, beside their red controls (not counted as reds):**
`TestNoInventoryIndexIsBeforeTheFirstBet` (control `TestEachVersionOfTheIndexIsItsExactBytes`),
and the clear-table case of the demo (`base..base`; its control is the judged table of the same
test). `TestAnAnswersCitationIsAQuestionIDNotARecordLine` is red at the stub because a zero input
lacks `Spans` (`covers` is `did not run: spec/spans.tsv`), not because of the answers mapping; its
controls (bare `Q-002` passes, `F-0001` named passes on its fact 2) make it distinguish the wrong
mapping, which mutant M1 shows. `TestReadRequiredKeepsTheStoreError...` was added after mutants
M18 and M19 survived (see `mutants.md`); it was green on the extraction by construction, and
its two mutants are the red evidence.

## Green

`go test -count=1 ./internal/run/` ok; the tree half of the integration level ok (1.0 s); the demo
`TestTheFullInputOfAHeadAfterIntakeAndTheScaffold` ok (56.3 s).

## Gates
- `go build ./...`: exit 0 (2 s)
- `go vet ./...`: exit 0 (0 s)
- `go vet -tags=integration ./...`: exit 0 (4 s)
- `go vet -tags=e2e ./cmd/layup/`: exit 0 (0 s)
- `gofmt -l .`: 0 files listed (exit 0)
- `go test -count=1 ./...`: exit 0 (9 s)
- `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...`:
  exit 0, 11:27:57Z to 11:55:39Z (internal/run 1572.7 s, catalog 47.2 s, setup 71.2 s, verify 72.7 s), at the
  head 5716475 (the later commits change no Go file).
- Text checks at that head, exit 0 each: `adr-lint` (4 s), `prd-lint` (0 s), `link-lint` (40 s, 2825 links),
  `run-discipline-tests` (17 s, 138 passed), `nested-checkout-check` (1 s), `setup-check` (69 s, after `git add`),
  `git diff --check origin/main HEAD` (1 s).

## Round 1, finding 1 (a PRD name that holds a line feed)

Red, before the change (the regex `^PRD-.*\.md$` at head 32e1853, the case `docs/prd/PRD-0002-\nextra.md` added):
`specinput_integration_test.go:55: PRDs [PRD-0001-x.md]; want PRD-0001-x.md and the PRD whose name holds a line feed, with their bytes (...)`.
Green after `isPRDName` (HasPrefix, HasSuffix): the tree tests `-tags=integration` ok (6.6 s).

Gates of the fix:
- `go build ./...`: exit 0 (6 s)
- `go vet ./...`: exit 0 (1 s)
- `go vet -tags=integration ./internal/run/`: exit 0 (1 s)
- `go test -count=1 ./...`: exit 0 (18 s)
- `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 15m -run 'TheTreeFieldsOfTheFullInput|NumberedFactsFileAtTheHead|FailedReadOfTheTreeIsASourceError|TheSpecFilesAreTheRegular|TheFullInputOfAHeadAfterIntake|ReadRequired|AHandoffOfTheTreeItems' ./internal/run/`: exit 0 (73 s), rerun with the expression single-quoted (round 2).
