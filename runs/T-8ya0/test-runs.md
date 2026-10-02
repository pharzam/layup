# T-8ya0: the test runs

The runs of `T-8ya0` (#88). Host: macOS, `go1.27.1`, `git` 2.54.0, macOS `awk`
20200816, 2026-10-02. A path of a temporary directory shows as `<tmp>/`; "…"
marks a cut part of an output.

## The order of the work

The core (`internal/verify/adapted.go`) was drafted while the plan review ran,
before its unit tests. So each red run below is of a test in its final form on
a **skeleton** of the core (the names, with no behaviour: no finding, an empty
`AD_EXCLUDE`, no `ad_allowed`), on the **base** `b0cb95a`, or on a **mutation**
of the core that breaks one rule; the comparisons with the sh function are the
evidence that each expected line is the line of `check_adapted`.

## The red runs

1. **The unit tests, on the skeleton:** the three tests fail, 50 lines.

   ```text
   $ go test -count=1 ./internal/verify/
   --- FAIL: TestTheHitsOfAdapted
       adapted_test.go:69: "This kit is copied.": got [] want ["rule-1 kit: docs/a.md:1"]
   … each case with a hit
   --- FAIL: TestTheTextOfAdapted
   … each case with a hit
   --- FAIL: TestTheFilesOfAdapted
       adapted_test.go:120: the files: got [] want ["rule-1 kit: AGENTS.md:1" "rule-1 kit: README.md:1" …]
       adapted_test.go:124: docs/adr/0001-a.md is not excluded; the test lists it as excluded
   … each excluded path, the directory and the missing file, the failed list
   ```

2. **The integration tests, on the skeleton:**

   ```text
   $ go test -count=1 -tags=integration -run 'TestTheFixturesOfSetupCheck|TestTheListsOfAdaptedEqualTheSh|TestFlagged' ./internal/verify/
       harness_integration_test.go:109: adapted/bad-rule-1: the lines of adapted …
       harness_integration_test.go:113: adapted/bad-rule-1: exit 0; want ["1"]
   … the same for bad-rule-2 and bad-rule-3
       harness_integration_test.go:146: AD_EXCLUDE of setup-check.sh …
       harness_integration_test.go:149: ad_allowed of setup-check.sh …
       harness_integration_test.go:168: bad-rule-1: [], <nil>; want ["docs/a.md" "docs/adr/0013-new.md" "docs/tests/t.md"]
       harness_integration_test.go:168: bad-rule-2: [], <nil>; want ["docs/b.md"]
       harness_integration_test.go:172: a directory that is not a repository: [], no error
   ```

3. **The rows of the built checks, on the base `b0cb95a`** (the changed test
   files on the code of the base, where `adapted` is not built):

   ```text
   $ go test -count=1 -run 'TestRunGivesEachRowInTheOrderOfTheTable|TestTheScratchTree' ./internal/verify/
       verify_test.go:127: the rows … {"adapted" "not-active" "not built yet"} …
       verify_test.go:244: the row {"adapted" "not-active" "not built yet"}; want not-active, scratch tree: add failed
   $ go test -count=1 -tags=integration -run TestRunOnAStandInWorkArea ./internal/verify/
       verify_integration_test.go:79: the rows …; want ["pin pass " "kit-history pass " "adapted pass " "identity pass "] …
   $ go test -count=1 -tags=e2e -run TestSetupVerifyOnAStandInWorkArea ./cmd/layup/
       verify_e2e_test.go:38: exit 1, stdout: … adapted	not-active	not built yet …
   ```

4. **Mutations of the core,** each run with `go test -run Adapted`:

   | Mutation | Result |
   | -------- | ------ |
   | `strings.ToLower` in place of the lower case of `A` to `Z` | first run: **pass**, a test that passed for the wrong reason (`İ` is shorter in lower case, so the hit stayed on its line); with the cases of a character that is shorter and one that is longer in lower case: fail, both cases |
   | the word one byte after the start of the match, as the sh function moves it | first run: **pass** (the extra finding of `ékit-history` had the text of the real `ékit` finding on that line); with `ékit-history` on its own line: fail |
   | no one-time rule (`slices.Compact` removed) | fail, each case with a hit |
   | no space removed at the start and end of a line (condition 1 of the plan review) | fail: the three cases of a space or a tab at a line end and a space at a line start |
   | `FindAllStringSubmatchIndex` on the whole paragraph in place of the loop of the sh function (note 1) | fail: two hits that share one space give one finding |

## The comparisons with the sh function

`sh docs/setup/setup-check.sh --only adapted` in the C locale (`LC_ALL=C`), and
the core on the same tree:

| Tree | Result |
| ---- | ------ |
| LAYUP's own tree at `b0cb95a` | both: no finding (`adapted OK`) |
| LAYUP's root commit `d2516fd`, the unchanged baseline | the same 441 lines in 38 files; the sh function in the host's UTF-8 locale gives the same lines |
| each of the 50 text cases of `TestTheHitsOfAdapted` and `TestTheTextOfAdapted`, one file each | the same 36 findings |
| a scratch repository: `ékit`, `ékit-history`, `İİİİ kit`, a CRLF line, a tab, and a file `docs/q"x.md` | the same four lines; only the core reads `docs/q"x.md` (D2: git quotes that path for the sh function) |
| the same scratch repository, in the host's UTF-8 locale | macOS `awk` 20200816 stops with `towc: multibyte conversion failure`, and the sh function prints `adapted OK` (the known limit in `setup.md`) |

## The green runs

On the tree of the commit that adds this file.

| Command | Result |
| ------- | ------ |
| `go build ./...`; `go vet` with no tag, `-tags=integration` and `-tags=e2e`; `gofmt -l .` | exit 0; no file |
| `go test -count=1 ./...` | `ok` × 11 packages |
| `go test -count=1 -tags=integration ./...` | `ok` × 11; the group `adapted` and the frame case through the harness, the drift guard, `Flagged`, the stand-in row `adapted pass` |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | `ok` × 11 |
| `go test -race -count=1 ./internal/verify/`, also with `-tags=integration` | `ok` |
