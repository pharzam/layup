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

## Step 4: the fixes of the verification of `0db1aa1`

A fresh verifier found the findings below in the head `0db1aa1`. The first
run is the new tests on that code. 2026-10-01T17:30:22Z,
`go test -count=1 ./internal/tsv/`, exit 1:

```text
--- FAIL: TestWriteRefusesARowThatReadWouldRefuse (0.00s)
    --- FAIL: TestWriteRefusesARowThatReadWouldRefuse/no_header_row:_the_column_names_as_the_first_row (0.00s)
        tsv_test.go:106: error <nil>; want an *Error at line 1
--- FAIL: TestReadRefusesARecordThatDoesNotMatch (0.00s)
    --- FAIL: TestReadRefusesARecordThatDoesNotMatch/a_carriage_return_before_a_line_feed (0.00s)
        tsv_test.go:150: error "line 2: a carriage return; save the file with line-feed endings only, and with no carriage return in a field" (line 2, column ""); want line 2, column "question", and "line-feed endings" in the text
[...] the same for a lone carriage return, with the column "excerpt"
    --- FAIL: TestReadRefusesARecordThatDoesNotMatch/no_header_row:_the_column_names_as_the_first_line (0.00s)
        tsv_test.go:148: rows [["file" "marker" "question"] ["docs/a.md" "m" "q"]]; want none
        tsv_test.go:150: error <nil>; want an *Error at line 1
--- FAIL: TestTypeCheckTakesTheValuesOfItsType (0.00s)
    types_test.go:76: id(Q-NNN): check("Q-0001"): no error; want one
    types_test.go:76: id(Q-NNN): check("Q-01000"): no error; want one
    types_test.go:76: id(SNN): check("S001"): no error; want one
FAIL	github.com/pharzam/layup/internal/tsv	0.863s
```

Each test fails for the right reason: the writer writes the column names as the
first row of a record with no header row, and the reader takes that line as a
row; the error for a carriage return names no column; and a run of `N` takes a
leading zero in a number that is longer than the run, so one number has two
forms.

Two changed tests pass on that code, because they guard against a mutant of the
verification: the schema test now wants the error of the schema, and the list
test also gives a carriage return to `JoinList`. On a copy of the fixed code,
each fails on its mutant (`go test -count=1 -run <the test>`, exit 1):

```text
M1, Read checks the types of the schema, not the column names:
    tsv_test.go:233: Read with columns []: error line 1: the header row has 1 fields; the schema has 0 columns; want the error of the schema
[...] the same for the other 4 schemas with a bad column name
M3, JoinList takes a carriage return:
    tsv_test.go:181: JoinList with "a\rb": no error; want one
```

The block parser, 2026-10-01T17:33:17Z, `go test -count=1 ./internal/tsv/`,
exit 1:

```text
--- FAIL: TestParseBlocksRefusesABlockThatDoesNotHaveTheForm (0.00s)
    block_test.go:68: a tab in the closing fence: error <nil>; want "line 3: a tab" in it
    block_test.go:68: a block in a blockquote: error <nil>; want "line 1: a tsv-schema fence after" in it
    block_test.go:68: a block after a list marker: error <nil>; want "line 1: a tsv-schema fence after" in it
    block_test.go:68: a block in a list item, four spaces in: error <nil>; want "line 3: a tsv-schema fence after" in it
--- FAIL: TestReadBlocksRefusesADirectoryThatItCannotReadOrWithNoBlock (0.00s)
    block_test.go:107: a directory that cannot be read: 0 blocks, error <nil>; want "permission denied" in it
    block_test.go:107: no block: 0 blocks, error <nil>; want "no tsv-schema block" in it
FAIL	github.com/pharzam/layup/internal/tsv	0.239s
```

Each test fails for the right reason: the parser takes a tab in a closing
fence, and it does not see a block after a blockquote mark, a list marker or
four spaces, so such a block escapes the test of the blocks; and `fs.Glob`
drops the error of a directory that cannot be read, so a wrong directory gives
no block and no error. `TestParseBlocksFollowsTheFenceRules` passes on that
code; on a copy of the fixed code it fails on each fence mutant of the
verification (`-run TestParseBlocksFollowsTheFenceRules`, exit 1):

```text
M4, a closing fence of the same length only:
    block_test.go:86: a closing fence longer than the opening fence: 0 blocks, error line 3: "`````": a column line has a name, a type, key or -, and a rule; want 1 blocks
M7, a fence line with an info string closes a fence:
    block_test.go:86: a fence line with an info string closes no fence: 1 blocks, error <nil>; want 0 blocks
M19, a backtick in the info string of a backtick fence is allowed:
    block_test.go:86: a backtick in the info string of a backtick fence: 0 blocks, error <nil>; want 1 blocks
M20, no limit on the spaces before a fence:
    block_test.go:86: four spaces before a fence: 0 blocks, error <nil>; want 1 blocks
```
