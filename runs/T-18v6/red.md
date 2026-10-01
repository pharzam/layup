# T-18v6: the red runs

The red evidence of the test-first steps of task `T-18v6` (#78), gate step 3.
The `pre-commit` hook refuses a commit with a red unit test (guardrails §2, "The
hook refuses a red commit"), so each step ran its new tests before its code, and
one commit holds the tests and the code of a step. In the red run of a unit
step, a stub held the new functions with their signatures and empty bodies (a
zero value and no error). Go 1.27.1 on darwin/arm64, base `b48764f`. The output
is trimmed to the relevant lines; `[...]` marks where like lines are cut.

## Step 1: the types, the writer and the reader

2026-10-01T16:50:24Z, `go test -count=1 ./internal/tsv/`, exit 1:

```text
--- FAIL: TestWriteAppliesTheFieldRule (0.00s)
    tsv_test.go:59: got ""
        want "id\trule\tline\texcerpt\tquestion\nQ-001\tG1\t0\t—\tWhich stack?\nQ-002\tG4\t12\ta b c d  e\t—\n"
--- FAIL: TestTheNoHeaderFormHasRowsOnly (0.00s)
--- FAIL: TestWriteRefusesARowThatReadWouldRefuse (0.00s)
    --- FAIL: TestWriteRefusesARowThatReadWouldRefuse/invalid_UTF-8_(D5) (0.00s)
        tsv_test.go:105: error <nil>; want an *Error at line 2
[...] the same for the other 9 cases
--- FAIL: TestReadRefusesARecordThatDoesNotMatch (0.00s)
    --- FAIL: TestReadRefusesARecordThatDoesNotMatch/a_carriage_return_before_a_line_feed (0.00s)
        tsv_test.go:147: error <nil>; want an *Error at line 2
[...] the same for the other 17 cases
--- FAIL: TestReadTakesAKeyOfTwoColumnsAsOneTuple (0.00s)
    tsv_test.go:156: rows [], error <nil>; want 4 rows
--- FAIL: TestReadGivesTheEmptyMarkAsTheEmptyValue (0.00s)
--- FAIL: TestAListValueHoldsNoSpaceAndAnEmptyListIsTheEmptyMark (0.00s)
    tsv_test.go:178: JoinList with "docs/my file.md": no error; want one
--- FAIL: TestWriteThenReadGivesTheSameRowsAndBytes (0.00s)
--- FAIL: TestWriteAndReadRefuseASchemaThatIsNotValid (0.00s)
--- FAIL: TestParseTypeRefusesATypeOffTheList (0.00s)
    types_test.go:29: parseType("float"): no error; want one
--- FAIL: TestTypeCheckTakesTheValuesOfItsType (0.00s)
    types_test.go:75: decimal: check("3"): no error; want one
    types_test.go:75: id(Q-NNN): check("Q-01"): no error; want one
[...] 71 lines "no error; want one" in this test
FAIL	github.com/pharzam/layup/internal/tsv	0.746s
```

Each test fails for the right reason: the stub writes no byte, reads no row,
refuses no row, input, list value or schema, and takes every type and every
value, so each assertion of a behaviour of the real code meets an empty result
or a missing error. `TestParseTypeTakesEachTypeOfTheClosedList` passed on the
stub; it guards the other failure, a parser that refuses a valid type.

## Step 2: the block parser and the comparer

2026-10-01T16:51:36Z, `go test -count=1 ./internal/tsv/`, exit 1:

```text
--- FAIL: TestParseBlocksReadsTheFormOfTheREADME (0.00s)
    block_test.go:38: got [], error <nil>
--- FAIL: TestParseBlocksRefusesABlockThatDoesNotHaveTheForm (0.00s)
    block_test.go:63: a tab in a column line: error <nil>; want "line 2: a tab" in it
[...] the same for the other 15 blocks
--- FAIL: TestReadBlocksRefusesTwoBlocksWithOneName (0.00s)
    block_test.go:79: got map[], error <nil>; want the blocks a and b
--- FAIL: TestCompareNamesEachDifference (0.00s)
    block_test.go:119: the order: error <nil>; want "column 2 (rule): the name: the block has \"rule\"; the Go schema has \"line\"" in it
[...] the same for the other 6 changes
FAIL	github.com/pharzam/layup/internal/tsv	0.245s
```

Each test fails for the right reason: the stub parser finds no block and
refuses none, so the two blocks outside the longer fences do not come back and
none of the 16 blocks without the form is refused; the stub comparer names none
of the 7 changes (its two checks of no difference pass on the stub).

## Step 3: the integration test of the blocks of `docs/spec/`

2026-10-01T16:52:14Z, with both lists empty,
`go test -count=1 -tags=integration ./internal/tsv/`, exit 1:

```text
--- FAIL: TestEverySchemaBlockIsBuiltOrNotYetBuilt (0.00s)
    blocks_integration_test.go:40: the block catalog-kinds is listed 0 times in built and notYetBuilt; want 1
[...] the same for the other 13 blocks of docs/spec/
FAIL	github.com/pharzam/layup/internal/tsv	0.263s
```

It fails for the right reason: `ReadBlocks` parsed the 14 blocks of
`docs/spec/` with no error, and no block was on a list. At 16:52:22Z, with
`prices` in both lists and `psb-gaps` written as `psb-gap`, each other check
failed on its own case:

```text
    blocks_integration_test.go:39: psb-gap is listed, but no block of docs/spec/ has that name
    blocks_integration_test.go:44: the block prices is listed 2 times in built and notYetBuilt; want 1
    blocks_integration_test.go:44: the block psb-gaps is listed 0 times in built and notYetBuilt; want 1
```
