# T-3jpx — the stack catalog package

Issue: [#81](https://github.com/pharzam/layup/issues/81), row 4 of the
[implementation plan](../plan/README.md); child of `layup gate` (#33). Serves
`F-0003#44`. Base `a0dd19c` (the merge of row 3, #103). Author: Claude Opus 5.5
on Claude Code. Evidence: [`runs/T-3jpx/`](../../runs/T-3jpx/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #81. The plan
review (Claude Fable 5.1 on Claude Code, a fresh session) gave
`approve-with-conditions`: Budget maximum 1,300 lines added plus removed over 22
files against the base, close-out inside; Cycle cap 1. The author applied the
condition (each rule of an entry that the reader checks is written into
`docs/spec/setup.md`) and the eight notes.

## What was done

1. **Test first** ([`test-runs.md`](../../runs/T-3jpx/test-runs.md)): the two
   measurements of `embed` that D1 rests on; the unit tests of the reader on
   in-memory entries and the integration tests of the embedded test entry and of
   the two schema blocks, each red on a skeleton of the package. The first red
   run found a test that could not fail (the schema test looked a block up by
   the Go schema's own name); it now looks each block up by its expected name.
2. **`internal/catalog`:** `Read` reads an entry by the schema `catalog-kinds`
   and checks its rules; an entry gives its kinds, its files with their target
   paths and `{{module}}` replaced, the fixture of a kind, the bytes of the
   manifest, the config paths of a kind, and whether a `catalog` ref names one
   of its files; `Stacks` lists the entries. The test entry
   `internal/catalog/testdata/test/` has a `go.mod` and a `.github/` workflow
   (both `.tmpl`), an `active` and a `pending` kind, and the fixture of the
   `active` one.
3. **The registration:** `catalog-kinds` and `gate-manifest` move to `built` in
   the block test of `internal/tsv`; `.gitattributes` holds the bytes of the
   catalog (`-text`), with its measured reason.
4. **`docs/spec/setup.md`, "The stack catalog":** the row of `files/` with
   `.tmpl`; D1 (K35), the rules of an entry, the writer of the manifest, the
   form of a `catalog` ref, the test entry and the bytes as values decided
   here. **`docs/spec/gate.md`:** one sentence beside the block
   `gate-manifest`, which two packages hold as a Go value.
5. **The tests' rows:** five rows in `docs/tests/traceability.md`; the
   `PRD-0001` §12 Test cells of `REQ-002`, `REQ-004` and `NFR-003`, with a §13
   row. `NFR-007` keeps its cell: `TestPackageRules` now holds
   `internal/catalog` too.

**The rejected alternatives:** a suffix only for the files that `embed` or the
Go tools trip on (`go.mod`, `*.go`): two rules where one does; one archive file
per entry: the files could not be read or reviewed one by one; a test entry in
the binary: S01 would accept a stack named `test`; a reader that needs no fixture
file: a kind whose fixture is missing would reach a target's setup.
