# T-c06a: the test runs

The runs of `T-c06a` (#91). Host: macOS, `go1.27.1`, `git` 2.54.0, 2026-10-02.
"…" marks a cut part of an output; `\u2039` and `\u203a` are the two angle
quotes of a marker, written as escapes so that check `markers` reads no marker
in this file.

## The measurements of the plan

On a scratch module (`go 1.26`, one badly formatted file `bad.go`), with the
two forms of the static check of `gofmt`:

| Form | `gofmt` on `PATH`, a bad file | no `gofmt` on `PATH` | a file that does not parse |
| ---- | ----------------------------- | -------------------- | -------------------------- |
| `test -z "$(gofmt -l .)"` | exit 1 | **exit 0** (a pass on nothing) | exit 1 |
| `out=$(gofmt -l .) && test -z "$out"` | exit 1 | exit 127 | exit 2 |

A package with only one failing test file: `go test -count=1 ./...` exits 1.
A module with no package: `go vet ./...` exits 1 (`no packages to vet`, D9).

The evidence URLs, read on 2026-10-02 at about 18:58 UTC: `https://go.dev/doc/go1.26`
(the release notes of Go 1.26, February 2026); `https://pkg.go.dev/cmd/go@go1.26.0`
(the `go` command at `go1.26.0`: `go vet`, `go test -count=1`, and "Fmt runs the
command 'gofmt -l -w'"); `https://pkg.go.dev/cmd/gofmt@go1.26.0` (`-l`);
`https://github.com/actions/checkout/tree/v4` (`fetch-depth`: "0 indicates all
history"); `https://github.com/actions/setup-go/tree/v5` (`go-version-file`: a
minor version gives the latest patch); and, for the test entry,
`https://pubs.opengroup.org/onlinepubs/9799919799/utilities/sh.html` (IEEE Std
1003.1-2024).

The two fixtures, made with `git diff --cached` on a scratch repository and
applied to its clean commit: with `static.patch`, the static command exits 1
and the test command 0; with `test.patch`, the static command exits 0 and the
test command 1.

## The red runs

1. **The unit tests, on a skeleton** (`Gaps` gives nothing, `Embedded` an
   empty root): ten tests fail, each for its rule.

   ```text
   $ go test -count=1 ./internal/catalog/ ./internal/verify/
   --- FAIL: TestReadRefusesAnEntryThatBreaksARule      (16 cases: each rule of a kind and of a gap: "error <nil>; want one that holds …")
   --- FAIL: TestFilesReplaceTheModuleAndNothingElse    (docs/floor.txt "{{gap:the floor}}" want "\u2039the floor\u203a"; the module paths "a b", "a\nb", "a\tb": no error)
   --- FAIL: TestHasNamesTheFilesOfTheEntryByTheirPathInIt   (Has("go/gaps.tsv") = false, want true)
   --- FAIL: TestGapsGiveEachGapAtItsLine
   --- FAIL: TestTheEntriesOfTheBinary                  (the binary embeds no entry: [], <nil>)
   --- FAIL: TestTheKindsOfTheGoEntry, TestTheFilesOfTheGoEntry, TestTheWorkflowOfTheGoEntry, TestTheJobScriptOfTheGoEntry
                                                        (catalog: open go/kinds.tsv: file does not exist)
   --- FAIL: TestACatalogRefOfTheBinary                 (catalogHas("go", "go/kinds.tsv") = false, want true, …)
   ```

2. **The block tests**, with `catalog-gaps` listed and no block in
   `docs/spec/`: `TestTheSchemaBlocks` ("docs/spec/ has no block
   catalog-gaps") and `TestEverySchemaBlockIsBuiltOrNotYetBuilt` ("catalog-gaps
   is listed, but no block of docs/spec/ has that name").

3. **Mutations,** each run against the integration tests (and M1, M2 and M9
   against the e2e test), each file restored after its run:

   | Mutation | Result |
   | -------- | ------ |
   | M1 the static fixture is formatted | fail: "the fixture of static: pass; want fail"; e2e: exit 1, but no row `static fail` (the pending kinds fail on the new `.go` file, so the test checks the row) |
   | M2 the test of the test fixture passes (`t.Log`) | fail: "the fixture of test: pass; want fail"; e2e the same |
   | M3 the static fixture changes a file that is not there | fail: "the fixture of static does not apply to the clean commit" |
   | M4 the job passes each pending kind | fail: three cases, "the job pend clear …; layup gate fail pending: product path changed: src/a/b.txt" |
   | M5 the job does not look up the tool | fail: "the job notool pass; layup gate not-active tool not found: no-such-tool-of-91" |
   | M6 the job does not check the form of the manifest | fail: three cases; with a scope pattern of another form, the job passed with "no product path" |
   | M7 a name that is only the extension matches | fail: "the job pend fail … src/.txt; layup gate clear" |
   | M8 the command of the test kind fails on each tree | fail: "good code: test fail exit 1; want pass" |
   | M9 the binary embeds no entry (e2e) | fail: "the binary embeds no entry: [], <nil>" |

## The green runs

On the tree of the commit that adds this file.

| Command | Result |
| ------- | ------ |
| The local checks of `AGENTS.md`; `go build ./...`; `go vet` with no tag, `-tags=integration` and `-tags=e2e`; `gofmt -l .` | exit 0; no file |
| `go test -count=1 ./...` | `ok` × 11 packages (`internal/catalog` 0.7 s) |
| `go test -count=1 -tags=integration ./...` | `ok` × 11 (`internal/catalog` 5.5 s: the fixtures of each entry, good code, and the job script against `layup gate`) |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | `ok` × 11 (`cmd/layup` 6.5 s; `TestGateOnEachEntryOfTheCatalog` about 2 s) |
| `go test -race -count=1` on `internal/catalog`, `internal/verify`, `internal/work` and `internal/standin`, and on `internal/catalog` and `internal/verify` with `-tags=integration` | `ok` |
| `sh docs/setup/tests/run.sh`; `sh docs/tests/run-discipline-tests.sh` | 44 passed, 0 failed; 81 passed, 0 failed |

## The red runs of the fix of review round 1

Round 1 (`b00d5c9`) gave one material finding (the job passed on forms of a
manifest that `layup gate` refuses) and notes 2 to 5. The new cases of
`TestTheJobScriptOfEachEntry`, in their final form, with each `sh` of the host,
on the script of `b00d5c9`:

```text
$ go test -count=1 -tags=integration -run TestTheJobScriptOfEachEntry ./internal/catalog/
    go: an empty field: the job "ok\tpass\t—\n" (passed true), layup gate …: line 2, column "command": an empty field …
    go: an empty field of config: the job "ok\tpass\t—\n" (passed true), layup gate …: line 2, column "config": an empty field …
    go: a config value that is not a path: the job "ok\tpass\t—\n" (passed true), layup gate …: "../x" is not a value of list(path)
    go: a scope of two spaces: the job "ok\tpass\t—\n" (passed true), layup gate …: "./*.txt  ./*.md" is not a value of list(text)
    go: a carriage return in a row: the job "ok\tpass\t—\n" (passed true), layup gate …: a carriage return …
    go: an empty line: the job "ok\tfail\tthe manifest: line 3 has 0 fields, not 6\n" …; want a failed job that names "an empty line"
    go: no line feed at the end: the job "ok\tpass\t—\n" (passed true), layup gate …: no line feed after the last line
    go: a byte that is not UTF-8: the job with [bash --posix] gives "ok\tclear\tno product path\n"; with [sh] "ok\tpass\t—\n"
    go: a builtin as the tool: the job "ok\tpass\t—\n" (passed true); layup gate not-active "tool not found: :"
    go: a base that is a tree: the job "ok\tpass\t—\n" (passed true), layup gate the revision "49dabb5…" is not a commit …
```

A defect found during the fix: on a pending case, `bash` 5.3 gave `clear` in 4
runs of 20, and `dash` gave `fail` in 20 of 20. The author took the trap on
`EXIT` for its cause, removed the trap, and measured 30 of 30 right; that
reading was wrong (below).

## The late run of round 2, and the fix of its defect

The first run of round 2 (Claude Fable 5.1, on `b9e1b05`) wrote its record at
15 min 9 s, so it is skipped (Bootstrap mode rule 4). Its text reported a
defect, which the author measured on this host before the round that counts:

```text
$ for i in $(seq 1 100); do … bash --posix .github/gates.sh pend; done | sort | uniq -c    (the script of b9e1b05, bash 5.3.9)
      5 pend clear pending: no product path
     95 pend fail pending: product path changed: src/a/b.txt
```

A subshell of the script crashed (the late record names a segmentation fault
of the pipeline segment `first`); its cause is not known. The script mapped
each failure of that segment to `clear`, a pass with the command not run. The
fix: the scope check reads a file, not a pipeline, and gives three answers
(`path <path>`, `none`, or a failure, which is `not-active`); each read of a
field checks its status; the tool is found as `exec.LookPath` finds it (a
`dash` builtin that is also a program, such as `true`, is found again); a last
byte NUL is no line feed. With the new script, each of three cases (a pending
kind, two active kinds) gave the right row in 100 of 100 runs with `sh`,
`dash` and `bash` 5.3 (900 runs). The new cases of the parity test, on the
script of `b9e1b05`: "a program that is a builtin of dash: the job ok
not-active tool not found: true; layup gate pass" and "a last byte NUL: the
job … line 7 is an empty line". Then `go test -count=3 -tags=integration` of
the three integration tests of the catalog: `ok`.
