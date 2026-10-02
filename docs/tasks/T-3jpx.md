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

## Verdict

Delivered: `internal/catalog`, the reader of a catalog entry by its rules (the
schema `catalog-kinds`, the fixture of each `active` kind, `—` for a `pending`
kind, no stray fixture, the `.tmpl` form of `files/`), which gives the kinds,
the files with their target paths and the module path, the fixture of a kind,
the bytes of the manifest, the config paths of a kind, the lookup of a `catalog`
ref, and the list of the entries; the embedded test entry with its `go.mod` and
`.github/` files (the demo); `catalog-kinds` and `gate-manifest` built; the
bytes of the catalog held by `.gitattributes`; K35 and the rules of an entry in
`docs/spec/setup.md`; five traceability rows; the `PRD-0001` §12 Test cells of
`REQ-002`, `REQ-004` and `NFR-003`, with a §13 row.

The plan review (Claude Fable 5.1) gave `approve-with-conditions`; the author
applied its condition and its eight notes. Review round 1 (Claude Fable 5.1,
`5a923bf`, cycle 0; [the record](../../runs/T-3jpx/review-round-1.md)) gave
`nothing material in scope`, with five notes. Notes 1, 2 and 5 are applied: the
reason of D1 now says that `embed` skips a directory with a `go.mod` with no
error (measured again, and the specification asks each entry's test to check its
files); the header of `.gitattributes` counts the third block; an empty fixture
shows as `—` in the error. Notes 3 and 4 need no change: `embed` cannot make a
`files` that is a file, or a `kinds.tsv` that is a directory, and the error of
`Read` names the file in both cases. At the head, `go build`, `go vet` with each
tag, `gofmt`, the three test levels, `go test -race` on `internal/catalog` and
all local checks pass, and `review-record-lint` passes on the comments of #81.
The diff against `origin/main` is 921 lines over 18 files with the close-out,
inside the Budget maximum of 1,300 lines over 22 files.

Next: row 5 of the plan (`T-5sgt`, #82), `layup gate`.

## Resource record

Recorded, not budgeted (ADR-0007). Times are 2026-10-02, UTC. Token counts are
`not reported` where the harness does not give them; the review sessions ran
with `--output-format stream-json`, which gives them.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan, with the two measurements of D1 | reasoning | Claude Opus 5.5 | max | not reported | notes within 07:23 to 07:26 (beside the review round of `T-2yw7`); 07:32 to 07:34 |
| The plan review | reasoning | Claude Fable 5.1 | the default of `claude -p` | 1,154,092 (input 290, cache write 147,402, cache read 976,217, output 30,183 of which thinking 17,094); USD 4.70 at list price | 7 min 41 s, 07:34:59 to 07:42:40 |
| The answer to the plan review | reasoning | Claude Opus 5.5 | max | not reported | 07:43 |
| The tests, the code, the specification and the records; the freeze | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | drafts within 07:35 to 07:41; 07:43 to 07:47 |
| Review round 1 | reasoning | Claude Fable 5.1 | the default of `claude -p` | 1,018,316 (input 290, cache write 134,422, cache read 847,239, output 36,365 of which thinking 18,649); USD 4.72 at list price | 9 min 7 s, 07:47:37 to 07:56:44 |
| The notes of round 1 and the close-out | execution | Claude Opus 5.5 | max | not reported | 07:57 to 08:01 |

**The size class** (note 7 of the plan review): the plan's table gives row 4 the
class `small`, "a diff under about 400 lines"; this task landed at 921 lines,
about 330 of them the evidence and the records. No row of the plan so far landed
under 400 lines. An input for the task that next edits the plan.
