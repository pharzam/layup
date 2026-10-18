# T-yeav — the tree items of the task loop

Issue [#404](https://github.com/pharzam/layup/issues/404). Serves REQ-005/NFR-001/NFR-004/F-0003#45.
Group (b) of the Operator's three-group ruling (comment 6235098876): the comment items are
`T-4pje` (#396, landed) and the composers `T-gyvo` (#405). Author: Sonnet 5.5.
Plan comment 6235235730; plan review 6246370488 (GPT-6.1 Sol, approve-with-conditions, five
conditions all met); first review 6233740566 supplies C1-C12. Maximum 1,600 changed
lines/21 files against main at landing, cycle cap 1, no panel.
Evidence: [runs](../../runs/T-yeav/test-runs.md) (stub red, green, mutants, gates).

## Delivered and boundaries

`spec-sections` and `rendering-equal` are computed by `RecordHandoff` from the tree at a head
that code holds, in the row of the pair that lists them, and a value of the caller's for either
is never read. The head is the SHA of the session's `bound` event (`boundHead`), or the `Head`
that code sets on a handoff that code wrote and that is read for no other; none is false.
`spec.HasSection` is the one section predicate (the full check calls it too). `rules.RecordRendered`
is the one per-record rule (`Judge` and `rules.RenderingEqual` both call it). A failed call of
`git` (the head, a `PreRev` that does not resolve, a tree or a blob), a failed register read, is an
error and no event; every refusal of the input or of the tree is false. Nothing is wired into the
step loop (row 66k). The contract is in `loop.md` (decided there before code, own commit).

## Tests mapped to the criteria

| Criterion | Test |
| --- | --- |
| One section predicate, boundaries, no ID | `TestHasSectionNeedsANonEmptySectionWhoseHeadingHoldsTheID`, `TestHasSectionReadsTheIDAtItsBoundaries`, `TestHasSectionOfNoIDOrNoFilesIsFalse`, `TestSectionsOfNoRequirementIDAreFalse` |
| Spec files are regular `.md` blobs at depth one | `TestTheSpecFilesAreTheRegularMdBlobsOfDocsSpecAtTheHead` |
| Shared per-record rule (C5b, condition 3) | `TestRecordRenderedIsTheOneRuleOfJudgeAndOfTheItem` and the existing `TestJudge*` |
| Refusals, static first, with controls (condition 3) | `TestRenderingEqualRefusesEachBadInputAndItsControlPasses` |
| Pre-image from the named revision, forged byte outside the rendered fields (C5a, condition 5) | `TestRenderingEqualTakesThePreImageFromTheNamedRevisionNeverTheHead`, `TestARenderedTaskHandoffNeedsAPreImageThatIsNotTheHead` |
| One error contract (condition 2) | `TestRenderingEqualGivesAnErrorForEachFailedCallOfGit`, `TestAFailedReadOfTheTreeIsAnErrorAndWritesNoEvent`, `TestRenderingEqualReadsARealTreeWithLinksAndModes` |
| Head is the bound one, a supplied Head ignored (condition 1) | `TestATreeItemWithNoHeadIsFalseAndNotAnError`, `TestASessionHandoffIgnoresASuppliedHeadAndTheCloneHead` |
| Both literals overwritten, only where the row lists the item (C6) | `TestTheTreeItemsAreComputedOnlyWhereTheRowListsThem` |
| Demo, eight outcomes on a real tree | `TestAHandoffOfTheTreeItemsIsValidOrInvalidByTheTreeNotByTheCaller` |

## Alternatives

| Decision | Selected | Rejected, and why |
| --- | --- | --- |
| Section predicate | one exported `spec.HasSection` called by `sections` | exporting the method of `fullCheck`: a second copy of the regexp |
| Per-record rule | exported `rules.RecordRendered` called by `Judge` and by the item | a private `renderedReason`: `Judge` would not call the exported name (condition 3) |
| Error contract | a wrapper `Repo` keeps the first git error, so `preImages`' mapping of an unresolvable `PreRev` to a reason does not hide it | trusting `ReasonGit` only: the case the review named |
| Pre-image base | an empty base skips the base link of `preImages`; the pre-image is only the head's ancestor, never the head | inventing a pull-request base for a handoff that has none |
| Files of `docs/spec` | recursive listing filtered on depth, type, mode and extension | `blobAt`: checks the type, not the mode |
| Which handoffs compute | only where the row lists the item | every handoff: unrelated handoffs would depend on a clone |
| Head of a session | its `bound` event | a field of the caller or the clone's `HEAD` |
| Evidence files | one `runs/T-yeav/test-runs.md` | two files: counted against the file maximum |

## Verification and verdict

Own test-first tests (stubs committed alone, red recorded, tests committed with the code), the
integration demo and 30 valid compiled mutants (all killed, one after a test fix; three equivalent
recorded as such; build or vet failures were never counted) cover the predicate, the files,
the shared rule, the refusals, the error contract and the wiring. The whole integration level
passed (exit 0, 09:29:21Z to 09:58:08Z); build, vet, unit, the five text checks and the
nested-checkout check pass. Final exact-head CI remains a landing condition, with its result on the issue.

Round 1 (GPT-6.1 Sol, xhigh, on Codex), record 6248663528: **material**, two findings. F1 an empty
`PreRev` of a PRD record was not refused before git; F2 the error of reading a spec blob had no
effective test. Both were fixed at 5263ab5 (F1 test first, red recorded; F2 needed no production
change and its mutant is killed). Round 2 (GPT-6.1 Sol, xhigh): **nothing material in scope**.
Main (#410) was merged cleanly; only index documents conflicted, and both sides were kept.
Harness exception: ADR-0026 names `claude`, `devin` and `opencode`; Codex was used on the
Operator's message of 2026-10-13, for both rounds.

Complete landing figure against main is recorded on the issue at the merge.

## Resource record

ADR-0007, recorded, not budgeted. Tokens are as the harness reported them.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Sonnet 5.5 | not recorded | not recorded | not recorded |
| Plan review | reasoning | GPT-6.1 Sol | xhigh | 168,081 | 2026-10-18T05:58:41Z to 06:03:40Z |
| Implementation | execution | Claude Sonnet 5.5 | not recorded | 258,221 | 59 min |
| Merge of main | execution | Claude Sonnet 5.5 | not recorded | not recorded | about 3 min |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 177,421 | 10:08:40Z to 10:17:03Z |
| Fix of round 1 | execution | Claude Sonnet 5.5 | not recorded | not recorded | about 5 min |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 159,339 | 10:22:12Z to 10:29:10Z |
| Close-out | execution | Claude Sonnet 5.5 | not recorded | not recorded | see the clock of the commit |
| Orchestration | reasoning | Claude Opus 5.5 | not recorded | not separately recorded | not separately recorded |
