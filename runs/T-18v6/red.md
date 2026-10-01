# T-18v6: the red runs

The red evidence of the test-first steps of task `T-18v6` (#78), gate step 3.
The `pre-commit` hook runs the unit level, so it refuses a commit with a red
unit test (guardrails §2, "The hook refuses a red commit"). So each step ran its
new tests before its code, and one commit holds the tests and the code of a
step. In a red run of a unit step, a stub held the new functions with their
signatures and empty bodies (a zero value and no error), so the tests compiled
and failed on their assertions. Go 1.27.1 on darwin/arm64, base `b48764f`. The
output is trimmed to the relevant lines; `[...]` marks a cut.

## Step 1: the types, the writer and the reader

2026-10-01T16:50:24Z. `go test -count=1 ./internal/tsv/` exits 1:

```text
--- FAIL: TestWriteAppliesTheFieldRule (0.00s)
    tsv_test.go:59: got ""
        want "id\trule\tline\texcerpt\tquestion\nQ-001\tG1\t0\t—\tWhich stack?\nQ-002\tG4\t12\ta b c d  e\t—\n"
--- FAIL: TestTheNoHeaderFormHasRowsOnly (0.00s)
    tsv_test.go:70: got ""; want "docs/a.md\tm1\tWhich value?\ndocs/b.md\tm1\t—\n"
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
    tsv_test.go:170: rows [], error <nil>; want ["P-001" "" "" "" "" "" "" "" "" "" "" ""]
--- FAIL: TestAListValueHoldsNoSpaceAndAnEmptyListIsTheEmptyMark (0.00s)
    tsv_test.go:178: JoinList with "docs/my file.md": no error; want one
[...] the same for the other 3 values
    tsv_test.go:185: got "", error <nil>; want "kind\tconfig\nstatic\t—\nlint\t.golangci.yml docs/gates.tsv\n"
--- FAIL: TestWriteThenReadGivesTheSameRowsAndBytes (0.00s)
    tsv_test.go:203: read back [], error <nil>; want [["Q-001" "G1" "0" "" "Which technology stack?"] [...]]
--- FAIL: TestWriteAndReadRefuseASchemaThatIsNotValid (0.00s)
    tsv_test.go:225: Write with columns []: no error; want one
    tsv_test.go:228: Read with columns []: no error; want one
[...] the same for the other 5 schemas
--- FAIL: TestParseTypeRefusesATypeOffTheList (0.00s)
    types_test.go:29: parseType(""): no error; want one
[...] the same for the other 18 texts
--- FAIL: TestTypeCheckTakesTheValuesOfItsType (0.00s)
    types_test.go:75: int: check("07"): no error; want one
    types_test.go:75: decimal: check("3"): no error; want one
    types_test.go:75: time: check("2026-10-01T08:09:21.5Z"): no error; want one
    types_test.go:75: path: check("./a"): no error; want one
    types_test.go:75: id(Q-NNN): check("Q-01"): no error; want one
    types_test.go:75: list(text): check("a  b"): no error; want one
[...] 71 lines "no error; want one" in this test
FAIL	github.com/pharzam/layup/internal/tsv	0.746s
```

| Test | Why it fails, and why that is the right reason |
| ---- | ---------------------------------------------- |
| `TestWriteAppliesTheFieldRule` | The stub writes no byte: no header row, no `—`, no spaces of the field rule, no line feeds. |
| `TestTheNoHeaderFormHasRowsOnly` | The stub writes no row and reads no row. |
| `TestWriteRefusesARowThatReadWouldRefuse` | The stub writer refuses none of the 10 rows that the reader refuses. |
| `TestReadRefusesARecordThatDoesNotMatch` | The stub reader refuses none of the 18 inputs that do not match the schema. |
| `TestReadTakesAKeyOfTwoColumnsAsOneTuple` | The stub reads no row, so the 4 rows, three of them for one stall, do not come back. |
| `TestReadGivesTheEmptyMarkAsTheEmptyValue` | The stub reads no row, so `—` does not come back as the empty value of each type. |
| `TestAListValueHoldsNoSpaceAndAnEmptyListIsTheEmptyMark` | The stub `JoinList` refuses no value with a space, and joins no value. |
| `TestWriteThenReadGivesTheSameRowsAndBytes` | The stub reads no row back. |
| `TestWriteAndReadRefuseASchemaThatIsNotValid` | The stub takes each schema that is not valid. |
| `TestParseTypeRefusesATypeOffTheList` | The stub parser takes every text as a type. |
| `TestTypeCheckTakesTheValuesOfItsType` | The stub `check` takes every value. |

`TestParseTypeTakesEachTypeOfTheClosedList` passed on the stub, because the stub
takes every text. It guards the other failure, a parser that refuses a valid
type.
