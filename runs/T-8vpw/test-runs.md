# T-8vpw: the test runs

The runs of `T-8vpw` (#87). Host: macOS, `go1.27.1`, `git` 2.54.0, 2026-10-02.
A path of a temporary directory shows as `<tmp>/`; "…" marks a cut part of an
output; `^?` is the DEL character of a file name, and `\u2039` and `\u203a`
are the two angle quotes of a marker, written as escapes so that check
`markers` reads no marker in this file.

## The red runs

Each red run is of a test in its final form, on the base `2ea523b` or on a
skeleton of the code (the names, with no behaviour), or on a mutation of the
code that breaks one rule.

1. **#21, the two fixtures first,** on the script of the base:

   ```text
   $ sh docs/setup/tests/run.sh
   FAIL  markers/bad-blank-question: exit 0, want 1
   FAIL  markers/bad-blank-question: no output line: setup-check: markers FAIL question: docs/setup/open-gaps.tsv line 1 has no question
   FAIL  markers/bad-quoted-name: exit 0, want 1
   FAIL  markers/bad-quoted-name: no output line: setup-check: markers FAIL unlisted: docs/q^?x.md \u2039port\u203a
   ```

   Then `git ls-files -z | tr '\0' '\n'` and the blank question in
   `check_markers` (`f32c647`): 44 passed.

2. **The unit tests, on the skeleton:** six tests fail (35 lines):
   `TestTheMarkersOfALine`, `TestScanMarkers`, `TestTheOpenGaps`,
   `TestTheBaselineScripts`, `TestTheBrokenLinks`, `TestTheSourcesOfARecord`.

3. **The integration tests, on the skeleton:** the harness on the group
   `markers` (each of its five cases with a finding: "the lines of markers",
   "exit 0; want 1"), the drift guard of `MK_EXEMPT`, the block `open-gaps`
   (`TestTheSchemasEqualTheirBlocks` of `internal/work`), the scripts on a
   tree and on a clone of LAYUP, and the rows of the stand-in fail.

4. **The e2e tests, on the skeleton:** `TestSetupVerifyOnAStandInWorkArea`
   (the rows of the four checks) and `TestSetupVerifyWithNoBaselineScript`
   (`exit 1`, but no row `not-active` with `missing: …`) fail.

5. **Mutations:**

   | Mutation | Result |
   | -------- | ------ |
   | a backtick before an open quote alone makes a mention | fail: `` `\u2039`\u2039x\u203a `` and `` two lines: `\u2039State one `` |
   | a question of only blanks, or the empty mark, is a question | fail: "blank questions and the empty mark" |
   | a line `FAIL` of link-lint that names no file is skipped | first run: **pass**, a test that passed for the wrong reason (each such case exited 1 with no other line, so the rule "a failed exit with no file" gave the error); with a case where that line comes beside a line that names a file: fail |

## The comparison with the sh function

`sh docs/setup/setup-check.sh --only markers` in the C locale, and the Go core
(`markersFindings`) on the same tree:

| Tree | Result |
| ---- | ------ |
| LAYUP's own tree at the head | both: no finding (`markers OK`) |
| LAYUP's root commit `d2516fd`, the unchanged baseline | the same 154 lines |
| each case of `docs/setup/tests/markers/` | the lines of `EXPECT`; on `bad-stale` the sh function prints one more line than its `EXPECT` held (`stale: docs/y.md \u2039baz\u203a …`), which was added, so the harness's comparison as a set holds |

## The green runs

On the tree of the commit that adds this file.

| Command | Result |
| ------- | ------ |
| `go build ./...`; `go vet` with no tag, `-tags=integration` and `-tags=e2e`; `gofmt -l .` | exit 0; no file |
| `go test -count=1 ./...` | `ok` × 11 packages |
| `go test -count=1 -tags=integration ./...` | `ok` × 11; the group `markers` through the harness, the drift guard, the block `open-gaps`, the scripts on a tree and LAYUP's own scripts on a clone of its `HEAD` |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | `ok` × 11 |
| `go test -race -count=1 ./internal/verify/ ./internal/work/ ./internal/standin/`, and `./internal/verify/` with `-tags=integration` | `ok` |
| `sh docs/setup/tests/run.sh`; `sh docs/tests/run-discipline-tests.sh` | 44 passed, 0 failed; 81 passed, 0 failed |
