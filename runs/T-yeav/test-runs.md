# T-yeav test runs

Condition 4 of the plan review: the compiling stubs were committed alone (cc9625a), the tests were written and run against them (red, below), and the tests were committed with the green code (596f4b0). Times are from `date -u` read after each step. Output of each run was redirected to a file and its exit status read directly (no pipe).

## Red against the stubs (cc9625a)

Stubs: `spec.HasSection` false, `rules.RecordRendered` always `not equal`, `rules.RenderingEqual` `(false, nil)`, `specSections` false, `specFilesAt` empty, no wiring. Every failure below is an assertion on a value (a control that should be true, or a caller's `true` read as `valid`), never an undefined symbol; the build and `go vet` (both tags) were exit 0 before the run. A refusal row that is green at the stub is paired in its test with a control that is red, so it is not credited alone.

### unit ./internal/spec (09:11:28Z, exit 1)

```
--- FAIL: TestHasSectionNeedsANonEmptySectionWhoseHeadingHoldsTheID (0.00s)
    full_test.go:705: the control, a heading and a body line, is no section
--- FAIL: TestHasSectionReadsTheIDAtItsBoundaries (0.00s)
    full_test.go:733: "## REQ-00 — x" holds "REQ-00": false, want true
    full_test.go:733: "## (REQ-005)" holds "REQ-005": false, want true
    full_test.go:733: "## REQ-005" holds "REQ-005": false, want true
--- FAIL: TestHasSectionOfNoIDOrNoFilesIsFalse (0.00s)
    full_test.go:744: the control, the real ID, is no section
```

### unit ./internal/rules (09:11:28Z, exit 1)

```
--- FAIL: TestRecordRenderedIsTheOneRuleOfJudgeAndOfTheItem (0.00s)
    rendering_test.go:76: an equal facts record: RecordRendered = "not equal", want ""
    rendering_test.go:76: a PRD of its pre-image: RecordRendered = "not equal", want ""
    rendering_test.go:76: a link: RecordRendered = "not equal", want "not a regular file"
    rendering_test.go:76: an executable: RecordRendered = "not equal", want "not a regular file"
    rendering_test.go:76: a gitlink: RecordRendered = "not equal", want "not a regular file"
    rendering_test.go:76: a rendering with a finding: RecordRendered = "not equal", want "no rendering"
--- FAIL: TestRenderingEqualIsTrueOnlyForTheRenderingOfARegularFile (0.00s)
    rendering_test.go:90: the control, a facts record and a PRD of its pre-image: false, <nil>
--- FAIL: TestRenderingEqualRefusesEachBadInputAndItsControlPasses (0.00s)
    rendering_test.go:147: no artifact: the control is false, <nil>
--- FAIL: TestRenderingEqualTakesThePreImageFromTheNamedRevisionNeverTheHead (0.00s)
    rendering_test.go:183: the restored byte: false, <nil>, read []; want true, reading the pre-image at p
--- FAIL: TestRenderingEqualGivesAnErrorForEachFailedCallOfGit (0.00s)
    rendering_test.go:205: the head does not resolve: false, <nil>; want false and the error
    rendering_test.go:205: the PreRev does not resolve: false, <nil>; want false and the error
    rendering_test.go:205: the rev-parse fails: false, <nil>; want false and the error
    rendering_test.go:205: the ancestry fails: false, <nil>; want false and the error
    rendering_test.go:205: the tree fails: false, <nil>; want false and the error
    rendering_test.go:205: the show fails: false, <nil>; want false and the error
```

### unit ./internal/run (09:11:29Z, exit 1)

```
--- FAIL: TestSectionsOfNoRequirementIDAreFalse (0.00s)
    treeitems_test.go:28: one ID, present (the control): false, want true
    treeitems_test.go:28: two IDs, both present: false, want true
--- FAIL: TestATreeItemWithNoHeadIsFalseAndNotAnError (0.00s)
    treeitems_test.go:71: a session with no bound head, the caller's Head and literal: {Valid:true Detail:valid}, <nil>; want invalid spec-sections
    treeitems_test.go:74: the last event ["4" "handoff" "1" "S-1a2b3c4d" "" "" "valid" "2026-10-09T12:00:00Z"]; want the handoff invalid spec-sections
    treeitems_test.go:79: code-written with no Head, the caller's literal true: {Valid:true Detail:valid}, <nil>; want invalid rendering-equal
--- FAIL: TestTheTreeItemsAreComputedOnlyWhereTheRowListsThem (0.00s)
    treeitems_test.go:114: a replaced row that lists the item: {Valid:true Detail:valid}, <nil>; want invalid spec-sections
```

### integration ./internal/rules (09:11:31Z, exit 1)

```
--- FAIL: TestRenderingEqualReadsARealTreeWithLinksAndModes (0.69s)
    basic_integration_test.go:398: a regular file (the control): false, <nil>; want true
    basic_integration_test.go:410: a PRD of its earlier revision: false, <nil>; want true
    basic_integration_test.go:419: a PreRev that is no revision: false, <nil>; want false and an error
    basic_integration_test.go:419: a head that is no revision: false, <nil>; want false and an error
```

### integration ./internal/run (09:10:30Z to 09:11:22Z, exit 1)

```
--- FAIL: TestTheSpecFilesAreTheRegularMdBlobsOfDocsSpecAtTheHead (0.69s)
    treeitems_integration_test.go:157: a regular .md file (the control): sections false, want true (files [])
    treeitems_integration_test.go:163: files [], <nil>; want only a.md with its bytes
    treeitems_integration_test.go:166: a head the clone lacks reads as no file
--- FAIL: TestAHandoffOfTheTreeItemsIsValidOrInvalidByTheTreeNotByTheCaller (15.77s)
    treeitems_integration_test.go:180: the section is at the bound head, the caller says false: invalid spec-sections; want valid
    treeitems_integration_test.go:180: a heading with no body, the caller says true: valid; want invalid spec-sections
    treeitems_integration_test.go:191: the rendering is at the head, the caller says false: invalid rendering-equal; want valid
    treeitems_integration_test.go:191: a byte of the head is forged, the caller says true: valid; want invalid rendering-equal
--- FAIL: TestASessionHandoffIgnoresASuppliedHeadAndTheCloneHead (9.01s)
    treeitems_integration_test.go:205: the section at the bound head, a supplied Head without it: invalid spec-sections; want valid
    treeitems_integration_test.go:212: the section at the clone's HEAD only: valid; want invalid spec-sections
    treeitems_integration_test.go:216: the section at the supplied Head only: valid; want invalid spec-sections
--- FAIL: TestARenderedTaskHandoffNeedsAPreImageThatIsNotTheHead (3.36s)
    treeitems_integration_test.go:223: the head as its own pre-image: valid; want invalid rendering-equal
    treeitems_integration_test.go:226: no PreRev: valid; want invalid rendering-equal
    treeitems_integration_test.go:229: the control, the earlier revision: invalid rendering-equal; want valid
--- FAIL: TestAFailedReadOfTheTreeIsAnErrorAndWritesNoEvent (7.39s)
    treeitems_integration_test.go:239: a head the clone lacks: {Valid:true Detail:valid}, <nil>, 1 new events; want an error and none
    treeitems_integration_test.go:243: a PreRev that is no revision: {Valid:true Detail:valid}, <nil>, 1 new events; want an error and none
    treeitems_integration_test.go:251: a register that fails: {Valid:true Detail:valid}, <nil>, 1 new events; want an error and none
    treeitems_integration_test.go:255: the control, the store read again: invalid spec-sections; want valid
```

## Green (596f4b0 and the test commit 1a2f92d)

| Run | Exit | Time (UTC) |
| --- | --- | --- |
| `go test ./internal/spec ./internal/rules ./internal/run` (unit) | 0 | 09:32:41Z (rerun at the final head) |
| `go test -tags=integration ./internal/rules` | 0 | 7.0 s |
| named integration of `internal/run` (`TestTheSpecFilesAre`, `TestAHandoffOfTheTreeItems`, `TestASessionHandoffIgnores`, `TestARenderedTaskHandoffNeeds`, `TestAFailedReadOfTheTree`) | 0 | 09:13:52Z, 34.3 s |

Characterization: the existing `TestJudge*`, `TestBasic*`, `TestEachResultIsOneOfTheClosedSet` and the whole of `internal/spec` stayed green and unedited through the extraction of `HasSection`, `RecordRendered`, `recordInputOK` and `fileAt`. One test was retired with its written reason (loop.md point 6 now says a caller's value for these two items is never read): `TestTheProgramItemsOfTheCallerAreHandedToTheCheck`; its artifact-hash case is `TestAnArtifactWhoseHashDiffersIsInvalidArtifacts`.

## Mutants

Applied on a copy of the tree (`scratchpad/mut/tree`, no `.git`), each from the committed file, `go build ./...` and `go vet` (with the tags the test needs) first, then the named tests with the exit read directly; the file was then copied back from the worktree and its SHA-256 (12 hex) compared: every restore matched (`b0bdf4dfcfe4` treeitems.go, `7544a5ea2383` basic.go, `e9c07c98a4d9` full.go). A mutant that did not build or vet is not run, never killed.

| Mutant | Result | Killed by |
| --- | --- | --- |
| H1 `HasSection` ignores the ID boundary | killed | `TestHasSectionReadsTheIDAtItsBoundaries` |
| H2 a heading with no body is a section | killed | `TestHasSectionNeedsANonEmptySection...`, `TestSectionsJudgesOnlyTheDelivered` |
| H3 `HasSection("")` is not refused | killed | `TestHasSectionOfNoIDOrNoFilesIsFalse` |
| R1 the pre-image may be the head | killed | `TestRenderingEqualTakesThePreImage...`, `TestBasicPreImageAndBaseChains` |
| R2 the pre-image is read at the head | killed | `TestRenderingEqualTakesThePreImage...` |
| R3 `RecordRendered` does not require mode 100644 | killed | `TestRecordRenderedIsTheOneRule...`, `TestJudgeARecordEntryThatIsNoRegularFileFails` |
| R4 `RecordRendered` ignores a finding | killed | `TestRecordRenderedIsTheOneRule...`, `TestJudgeARenderingThatFailsIsNoPass` |
| R5 no artifact and no record is true | killed | `TestRenderingEqualRefusesEachBadInput...` |
| R6 a record with no artifact is accepted | killed | the same |
| R7 a repeated artifact path is accepted | killed | the same |
| R8 a path `ValidPath` refuses is judged | killed | the same (the no-call assertion), `TestJudgeBadInputIsAnErrorNeverASuccess` |
| R9 a git error becomes `(false, nil)` | killed | `TestRenderingEqualGivesAnErrorForEachFailedCallOfGit` |
| R10 `Judge` does not apply `RecordRendered` | killed | four `TestJudge*` |
| R11 the head is not resolved before the pre-image comparison | killed | four `TestRenderingEqual*` |
| R12 an unresolvable `PreRev` is false, not an error | killed | `TestRenderingEqualGivesAnError...` |
| R13 a failed read of the head's tree is false | **survived**, then killed | first sweep: the error rows all failed at the pre-image read; two facts-only rows added (`...the tree of the head fails, with no PRD`), now killed |
| S1 no depth filter on `docs/spec` | killed | `TestTheSpecFilesAreTheRegularMdBlobs...` |
| S2 the mode filter is dropped | killed | the same |
| S3 the `.md` filter is dropped | killed | the same (the first form did not build: `strings` unused; not counted) |
| S3b the type filter is dropped | equivalent | a gitlink has mode 160000, so the mode filter already excludes it |
| S4 no requirement ID is true | killed | `TestSectionsOfNoRequirementIDAreFalse` |
| S5 any one ID suffices | killed | the same |
| S6 the head is the clone's `HEAD` | killed | `TestAHandoffOfTheTreeItems...`, `TestASessionHandoffIgnores...` |
| W1 the caller's literal survives when there is no head | killed | `TestATreeItemWithNoHeadIsFalse...`, `TestTheTreeItemsAreComputedOnlyWhere...` |
| W1b, W2 the caller's literal ORed with the computed value | equivalent | `treeItems` resets both items first; the real cuts are W1c and W2c |
| W1c, W2c the reset removed and the literal ORed | killed (both) | `TestAHandoffOfTheTreeItemsIsValidOrInvalidByTheTree...` |
| W3 a failed read of the tree is swallowed | killed | `TestAFailedReadOfTheTreeIsAnErrorAndWritesNoEvent` |
| W3b a failed read of the register is swallowed | killed | the same |
| W5 the items are computed although the row does not list them | killed | `TestTheTreeItemsAreComputedOnlyWhereTheRowListsThem` |
| W6 no head is an error | killed | the handoff tests of `handoffpost_test.go` (many) |

W4 (the event written before the compute) has no clean cut that builds and keeps the order of the calls readable, so it is not run; the order is asserted by `TestAFailedReadOfTheTreeIsAnErrorAndWritesNoEvent` (events unchanged). 29 valid mutants killed (one after a test fix), 3 equivalent, 0 survivors.

## Gates at the final head

Exit status read directly; the whole integration level ran in the background with `TMPDIR` under `layup-session/tmp-green`.

| Gate | Exit | Time (UTC) |
| --- | --- | --- |
| `go build ./...`, `go vet ./...`, `gofmt -l .` (empty) | 0, 0, 0 | 09:29:28Z to 09:29:30Z |
| `go test ./...` (unit) | 0 | 09:29:46Z |
| `go vet -tags=e2e ./cmd/layup/` | 0 | 09:29:48Z |
| `go test -count=1 -tags=integration -timeout 30m ./...` (whole level; `internal/run` 1644.6 s, `internal/rules` 11.7 s) | 0 | 09:29:21Z to 09:58:08Z |
| `adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests`, `setup-check`, `nested-checkout-check` | 0 each | 09:29:52Z to 09:31:55Z (re-run after the evidence files, below) |
