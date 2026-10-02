# T-3jpx: the test runs

The runs of `T-3jpx` (#81). Host: macOS, `go1.27.1`, 2026-10-02.

## The two measurements of D1 (K35)

A scratch module outside the repository (`example.com/p`), with a test file
that embeds a directory of the form of an entry. Recorded once, not a test of
LAYUP (note 3 of the plan review).

```text
$ go version
go version go1.27.1 darwin/arm64
# testdata/a/files/ holds go.mod.tmpl and .github/workflows/g.yml
//go:embed testdata/a        -> [testdata/a/files/go.mod.tmpl]
//go:embed all:testdata/a    -> [testdata/a/files/.github/workflows/g.yml testdata/a/files/go.mod.tmpl]
# testdata/b/files/ holds a plain go.mod, and nothing else
//go:embed testdata/b
cat/d_test.go:5:12: pattern testdata/b: cannot embed directory testdata/b: contains no embeddable files
# testdata/c holds kinds.tsv, files/go.mod, files/other.txt and fixtures/k.patch
//go:embed all:testdata/c   -> [testdata/c/fixtures/k.patch testdata/c/kinds.tsv]
```

So `embed` leaves out `.github/` without the prefix `all:`, and skips a
directory that holds a `go.mod`, with every file in it: with no error when the
pattern still matches other files (`testdata/c`: `files/other.txt` is lost
too), and with the error above only when nothing else is left (`testdata/b`).
The third measurement came from review round 1, note 1; the first conclusion
said "refuses", which was one case wide.

## The red runs

The texts of three errors changed after review round 1 (note 5): an empty
fixture now shows as `—`, and a fixture path shows with no quotes; the runs below
show the texts of that time. The tests ran on a skeleton of `internal/catalog`: its names, with no behaviour
(each function gives its zero value), and the empty Go schemas. So each test
fails on its own assertion. The first red run showed that `TestTheSchemaBlocks`
looked a block up by the Go schema's own name, so an empty schema compared
equal to no block; the test now looks each block up by its expected name, and
the run below is of that test. Outputs are cut to the relevant lines ("…").

```text
$ go test -count=1 ./internal/catalog/                                    # run 1
--- FAIL: TestReadGivesTheKindsInTheOrderOfTheFile
    catalog_test.go:52: the kinds [], want static and layout
--- FAIL: TestReadRefusesAnEntryThatBreaksARule
    catalog_test.go:89: no kinds.tsv: error <nil>; want one that holds "go/kinds.tsv"
    catalog_test.go:89: a bad state: error <nil>; want one that holds "go/kinds.tsv: line 2, column \"state\""
    catalog_test.go:89: a fixture of no kind: error <nil>; want one that holds "go/fixtures/extra.patch: the file is the fixture of no active kind (no active kind extra)"
    catalog_test.go:89: a file of files/ without the suffix: error <nil>; want one that holds "go/files/go.mod: a file of files/ ends with .tmpl"
    … the same for the other 7 cases
--- FAIL: TestReadRefusesANameThatIsNotAnEntry
    catalog_test.go:97: Read(""): no error
    … the same for ".", "go/files", "../go", "rust"
--- FAIL: TestFilesReplaceTheModuleAndNothingElse
    catalog_test.go:113: 0 files, want 3: []
--- FAIL: TestManifestIsTheKindsWithoutVersionFixtureAndEvidence
    catalog_test.go:134: the manifest: … want: kind	state	tool	command	scope	config …
--- FAIL: TestFixtureAndConfigOfAKind
    catalog_test.go:141: Fixture(static) = "", <nil>
    catalog_test.go:144: Fixture(layout), a pending kind: no error
    catalog_test.go:147: Config(layout) = [], false
--- FAIL: TestHasNamesTheFilesOfTheEntryByTheirPathInIt
    catalog_test.go:166: Has("go/kinds.tsv") = false, want true
    … the same for the other 3 files of the entry
--- FAIL: TestStacksListsTheDirectoriesWithAKindsFile
    catalog_test.go:179: Stacks = [], <nil>; want go and rust
$ go test -count=1 -tags=integration -run 'TestTheEmbeddedTestEntry|TestTheSchemaBlocks' ./internal/catalog/   # run 2
--- FAIL: TestTheEmbeddedTestEntry
    catalog_integration_test.go:34: Stacks = [], <nil>; want test
--- FAIL: TestTheSchemaBlocks
    catalog_integration_test.go:78: schema catalog-kinds: the name: the block has "catalog-kinds"; the Go schema has "" …
    catalog_integration_test.go:78: schema gate-manifest: the name: the block has "gate-manifest"; the Go schema has "" …
    catalog_integration_test.go:92: the manifest by the block gate-manifest: 0 rows, line 1: no header row
```

Run 3 is on the real package, with the pattern of the test entry changed to
`//go:embed testdata/test` (no `all:`) for the run. The test then loses the
`.github/` file of the entry; with `all:` back, it passes.

```text
$ go test -count=1 -tags=integration -run TestTheEmbeddedTestEntry ./internal/catalog/   # run 3
--- FAIL: TestTheEmbeddedTestEntry
    catalog_integration_test.go:52: the files ["go.mod"]; want .github/workflows/gates.yml and go.mod
```

## The green runs

On the tree of the commit that adds this file.

| Command | Result |
| ------- | ------ |
| `go build ./...` | exit 0 |
| `go vet ./...`; `go vet -tags=integration ./...`; `go vet -tags=e2e ./...` | exit 0 each |
| `gofmt -l .` | no file |
| `go test -count=1 ./...` | `ok` × 6 packages |
| `go test -count=1 -tags=integration ./...` | `ok` × 6; `TestTheEmbeddedTestEntry`, `TestTheSchemaBlocks`, `TestEverySchemaBlockIsBuiltOrNotYetBuilt` (with `catalog-kinds` and `gate-manifest` built) and `TestPackageRules` (which now holds `internal/catalog`) pass |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | `ok` × 6 |
| `go test -race -count=1 ./internal/catalog/` | `ok` |
| `git check-attr text` on a file of the test entry, and on `internal/catalog/catalog.go` | `unset`; `unspecified` |
