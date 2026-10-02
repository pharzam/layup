# T-5zmw: the test runs

The runs of `T-5zmw` (#83). Host: macOS, `go1.27.1`, 2026-10-02. A path of a
temporary directory shows as `<tmp>/`, and a raw carriage return as `<CR>`;
"…" marks a part of an output that is cut.

## What is red, and what is not

Note 2 of the plan review: the goldens `edge` and `values`, `TestEdgeCases`,
`TestCRLFGivesTheSameTable`, `TestGoldenRealPSB`, and the end-to-end runs on the
real problem statement, on a statement with no gap, on a missing file and on a
directory pin the behaviour of the base. They pass on the base, and this file
claims no red run for them. The contents of the golden tables are records of
the block on the base too, but the test that reads them so,
`TestEveryGoldenIsARecordOfTheBlock`, cannot run there: it needs `GapsSchema`
(red run 2; note 4 of review round 1). The red runs below
are of the behaviours that change: the lone carriage return, G1 (O-131), the
error of the output (through `WriteTSV`, `internal/cli` and the binary), the
UTF-8 check (K32), and the Go schema of `psb-gaps`.

## The red runs

1. **`internal/psb`, unit**, with the tests in their final form, the G1
   expression of the base, and `WriteTSV` as on the base (it writes the rows
   itself and gives `nil`):

   ```text
   $ go test -count=1 ./internal/psb/
   --- FAIL: TestALoneCarriageReturnBecomesASpace (0.00s)
       check_test.go:55: a lone CR: the table differs
           --- got ---
           id	rule	line	excerpt	question
           Q-001	G4	2	The API<CR>is fast.	Which number or threshold does "fast" stand for here?
           --- want ---
           id	rule	line	excerpt	question
           Q-001	G4	2	The API is fast.	Which number or threshold does "fast" stand for here?
   --- FAIL: TestG1ReadsTheValueOfAStack (0.00s)
       check_test.go:73: **Technology stack:**: the table differs
           --- got ---
           id	rule	line	excerpt	question
           --- want ---
           id	rule	line	excerpt	question
           Q-001	G1	0	—	Which technology stack does the product use (languages, frameworks, tools)?
   --- FAIL: TestWriteTSVGivesTheErrorOfItsOutput (0.00s)
       check_test.go:107: WriteTSV to an output that refuses each write: no error
   FAIL
   ```

2. **`internal/psb`, integration**, at the same point: the test of the block
   does not compile, because the Go schema does not exist yet.

   ```text
   $ go vet -tags=integration ./internal/psb/
   vet: internal/psb/check_integration_test.go:48:31: undefined: GapsSchema
   ```

3. **`internal/cli`**, with the seam `readFile` in place and nothing else
   changed (`WriteTSV` already gave its error; `psbCheck` dropped it):

   ```text
   $ go test -count=1 ./internal/cli/
   --- FAIL: TestPSBCheckRefusesAFileThatIsNotUTF8 (0.00s)
       psb_test.go:36: "\xff": exit 1, stdout "id\trule\tline\texcerpt\tquestion\nQ-001\tG1\t0\t—\tWhich technology stack does the product use (languages, frameworks, tools)?\n", stderr ""; want 2, nothing and "layup: in.md: line 1 is not valid UTF-8\n"
       psb_test.go:36: "Technology stack: Go\n\nThe café is f\xe9st.\n": exit 0, stdout "id\trule\tline\texcerpt\tquestion\n", stderr ""; want 2, nothing and "layup: in.md: line 3 is not valid UTF-8\n"
       psb_test.go:36: "Technology stack: Go\nthe file ends inside a character \xe2\x80": exit 0, stdout "id\trule\tline\texcerpt\tquestion\n", stderr ""; want 2, nothing and "layup: in.md: line 2 is not valid UTF-8\n"
       psb_test.go:36: "a surrogate half \xed\xa0\x80\n": exit 1, … want 2, nothing and "layup: in.md: line 1 is not valid UTF-8\n"
       psb_test.go:36: "one\ntwo\n\x80\n\xff\n": exit 1, … want 2, nothing and "layup: in.md: line 3 is not valid UTF-8\n"
   --- FAIL: TestPSBCheckGivesTwoWhenItCannotWriteTheTable (0.00s)
       psb_test.go:55: "No stack is named here.\n": exit 1, stderr ""; want 2 and the error of the output
       psb_test.go:55: "Technology stack: Go\n": exit 0, stderr ""; want 2 and the error of the output
   FAIL

   $ go test -count=1 -tags=integration -run TestPSBCheckExitCodes ./internal/cli/
   --- FAIL: TestPSBCheckExitCodes (0.00s)
       psb_integration_test.go:44: a file that is not valid UTF-8: exit 0, stdout "id\trule\tline\texcerpt\tquestion\n", stderr ""; want 2, nothing, and the reason "<tmp>/TestPSBCheckExitCodes329426185/001/bad.md: line 2 is not valid UTF-8" with no usage
   FAIL
   ```

4. **`cmd/layup`, end to end**, on the binary of the base: a copy of `4b47982`
   (`git archive`) with the harness and the scenarios of this branch. The write
   scenario had the name `TestPSBCheckGivesTwoWhenItCannotWriteTheTable` then;
   it is `TestPSBCheckWithAReadOnlyStandardOutput` now, so that its name differs
   from the unit test of `internal/cli`.

   ```text
   $ go test -count=1 -tags=e2e -run TestPSBCheck ./cmd/layup/
   --- FAIL: TestPSBCheckInputErrors (0.02s)
       psb_e2e_test.go:67: a file that is not valid UTF-8: exit 0, stdout "id\trule\tline\texcerpt\tquestion\n", stderr ""; want 2, nothing, and the reason "layup: <tmp>/TestPSBCheckInputErrors1724932604/001/bad.md: line 3 is not valid UTF-8\n" with no usage
   --- FAIL: TestPSBCheckGivesTwoWhenItCannotWriteTheTable (0.00s)
       psb_e2e_test.go:93: exit status 1, stderr ""; want exit 2 and the error of the write
   FAIL
   ```

## One measurement: a closed pipe and a read-only output

Condition 2 of the plan review asks for the result of a table that the command
cannot write. The binary of this branch, with no environment variable, on the
real problem statement, from a Python harness (recorded once, not a test):

```text
closed pipe (the read end closed before the start): returncode -13 stderr b''
read-only standard output: returncode 2 stderr b'layup: write /dev/stdout: bad file descriptor\n'
```

So a closed pipe stops the binary with `SIGPIPE`, as Go does for a write to a
broken pipe on standard output, and gives no code; the end-to-end scenario uses
a standard output that is open for reading only, where each write fails, and
`docs/spec/psb-check.md` says both.

## Review round 1: the fixes (cycle 1)

Review round 1 (`edbeb60`) gave `material` with three findings. The runs of the
fixes:

1. **Finding 1, G1 and white space.** The six new cases of
   `TestG1ReadsTheValueOfAStack` on the code of `edbeb60`: the run stops at the
   first failure, because `golden` calls `t.Fatalf`.

   ```text
   $ go test -count=1 -run TestG1ReadsTheValueOfAStack ./internal/psb/
   --- FAIL: TestG1ReadsTheValueOfAStack (0.00s)
       check_test.go:82: "Technology stack: \u00a0\n": the table differs
           --- got ---
           id	rule	line	excerpt	question
           --- want ---
           id	rule	line	excerpt	question
           Q-001	G1	0	—	Which technology stack does the product use (languages, frameworks, tools)?
   FAIL
   ```

   The binary of `edbeb60` on the three values that are only white space (a
   file each, `Technology stack: ` and then the character):

   ```text
   U+00A0: 0 gap rows, exit 0
   U+2003: 0 gap rows, exit 0
   vertical tab: 0 gap rows, exit 0
   ```

   The three other new cases (`Technology stack:` with U+00A0 or a vertical tab
   and then `Go`, no gap; with a tab and then `Go`, a gap) pass on `edbeb60`:
   they pin results that the fix keeps. After the fix, two checks of the new
   expression, each a small Go program outside the repository (recorded once,
   not a test): its two classes equal the Unicode property `White_Space` (by
   `unicode.IsSpace`) for each rune, with 0 differences; and on 465,010 lines
   (`Technology stack`, then one of 7 prefixes, then `:`, then each string of
   up to 5 characters from a space, `*`, a tab, a vertical tab, a form feed, a
   carriage return, U+00A0, U+2003 and `G`) it names a stack on no line where
   the expression of the base names none, and each of the 129,420 lines where
   the two differ is a line where the base read `*` or a white space character
   as the value.

2. **Finding 2, the edge golden**, and **finding 3, `Fast` in mixed case:** the
   new rows of `edge.md` and `edge.tsv` and the changed line of `values.md` pass
   at once. They pin the behaviour of the base (note 2): the code had it, and
   the goldens did not hold it.

## The green runs

On the tree of the commit that adds this table, and again on the head of
cycle 1.

| Command | Result |
| ------- | ------ |
| `go build ./...` | exit 0 |
| `go vet ./...`; `go vet -tags=integration ./...`; `go vet -tags=e2e ./...` | exit 0 each |
| `gofmt -l .` | no file |
| `go test -count=1 ./...` | `ok` × 7 packages |
| `go test -count=1 -tags=integration ./...` | `ok` × 7; `TestGoldenRealPSB`, `TestEveryGoldenIsARecordOfTheBlock`, `TestPSBCheckExitCodes`, `TestEverySchemaBlockIsBuiltOrNotYetBuilt` (with `psb-gaps` in `built`), `TestPackageRules` (with the import of `internal/tsv` by `internal/psb`) and `TestInputRule` pass |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | `ok` × 7; the four scenarios of `cmd/layup/psb_e2e_test.go` pass |
| `go test -race -count=1 ./internal/psb/ ./internal/cli/`, also with `-tags=integration` | `ok` |
| `git diff 4b47982 HEAD -- internal/psb/testdata/psb.tsv docs/facts/` | empty: the golden of the real problem statement and `F-0004` do not change |
