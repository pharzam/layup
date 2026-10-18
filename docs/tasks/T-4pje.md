# T-4pje — the comment items of the task loop

Issue [#396](https://github.com/pharzam/layup/issues/396). Serves REQ-005/NFR-001/NFR-004/F-0003#45.
Group (a) of the Operator's three-group ruling (comment 6235098876): the tree items are
`T-yeav` (#404) and the composers `T-gyvo` (#405). Author: Sonnet 5.5.
Plan comment 6246373422 (the repost of 6235223916); plan review 6246378861 (GPT-6.1 Sol,
approve-with-conditions); first review 6233740566 supplies C1-C12. Maximum 1,600 changed
lines/20 files against main at landing, cycle cap 1, no panel.
Evidence: [runs](../../runs/T-4pje/test-runs.md) (red, green, demo, mutants, gates).

## Delivered and boundaries

`plan-fields`, `verdict-field`, `findings-located` and `nothing-material` are computed by
`RecordHandoff` from the kept copy of the comment, never from the caller's value (either
literal). The kept bytes are parsed only when their SHA-256 is the hash of the result's
comment row, which `artifacts` verifies. `nothing-material` compares the review record's
`head:` with the `base` of the verifier session's start row (C2, settled in `loop.md`
before code). The line reader and the strict whole-comment parse of the three forms are in
`internal/run/commentform.go`; `ParsePlanComment`, `ParsePlanReview`, `ParseReviewRecord`
and the fail-uphill `OpenMaterial` are the exposures that #367 consumes. `rules.ValidPath`
is exported for the path rule. Nothing is wired into `Position` or the step loop (row 66k).

## Tests mapped to the criteria

| Criterion | Red test (failed on a semantic assertion against the stubs) |
| --- | --- |
| Pinned grammar, each rule with a control | `TestThePlanCommentParsesInThePinnedForm`, `TestAPlanFieldTwiceOrAbsentIsNoPlan`, `TestTheTestsMarkerOnceWithBulletsOfTwoWords`, `TestAPlanListFieldIsOneLineOfText`, `TestAFieldValueIgnoresOneCarriageReturn`, `TestAMarkerOrBulletWithACarriageReturnIsNotRead`, `TestAFieldWithNothingAfterTheColonIsMalformed`, `TestTheLineReaderIgnoresAFencedBlock` |
| Whole-comment strict, invalid input | `TestACommentThatFailsItsFormIsFalseEverywhere`, `TestAnEmptyOrInvalidCommentParsesToNothing` |
| The four items | `TestVerdictFieldIsApproveOrApproveWithConditions`, `TestFindingsLocatedNeedsALocatedFinding`, `TestNothingMaterialNeedsTheVerdictTheHeadAndAnEmptyBlock` |
| Fail-uphill exposure (C9) | `TestOpenMaterialFailsUphill` |
| Both literals overwritten, kind dispatch (C6) | `TestTheCommentItemsAreComputedFromTheKeptCopy`, `TestAnItemOfAnotherCommentKindIsFalse` |
| Verifier head is the start row's `base`, no `bound` event (C2) | `TestNothingMaterialReadsTheBaseOfTheVerifierStartRow` |
| Refusals: code-written, absent, empty, invalid kept copy, store error | `TestACodeWrittenHandoffHasNoCommentItem`, `TestAKeptCopyTheStoreLacksOrCannotReadMakesTheItemsFalse`, `TestAStoreErrorReturnsAnErrorAndWritesNoHandoff` |
| Parsed bytes are the verified bytes (round 1, F1) | `TestTheKeptFileThatIsNotTheNamedHashIsNoComment` |
| Demo (integration, real git and the `httptest` forge) | `TestTheCommentChecksOfTheTaskLoop` |

## Alternatives

| Decision | Selected | Rejected, and why |
| --- | --- | --- |
| Verifier's reviewed head (C2) | the `base` of its start row: present from the start row, written by code before the session runs, and the commit its clone is made at | its `bound` event: exists only after a passed push and equals the reviewed head only because the verifier commits nothing |
| Parse scope | whole-comment strict: a comment failing any part of its form makes each of its items false | per-item parse: a bad `head:` would leave `findings-located` true |
| Bullet with `\r` | malformed, pinned in the grammar | implied by the path rule: `ValidPath` does not reject `\r` |
| Tests bullet path | exactly two words; a space in a path is malformed | rest-of-line path: lets `- t1 p extra` pass silently |
| List fields | `—` is the empty list, any other text is one element | a comma split: invents a separator that the tier does not need |
| Dispatch | by the handoff's kind; an item of another kind is false | computing every item for every copy |
| Missing kept copy | the four items false, the event still written; a failing read is an error with no event | treating the absence as valid, or swallowing the error |
| Evidence files | one `runs/T-4pje/test-runs.md` | four files: counted against the 20-file maximum |

## Verification and verdict

Own test-first unit tests, the integration demo and 35 valid compiled mutants (all
killed; build or vet failures were excluded, never counted as kills) cover the
pinned grammar, the four items, both caller literals and the refusals. The whole
integration level passed as root at the fix head (exit 0, 07:52:47Z to 08:22:04Z on
a loaded host); build, vet, unit, the five text checks and the nested-checkout check
pass. Final exact-head CI remains a landing condition, with its result on the issue.

Round 1 on 284b8c4 (GPT-6.1 Sol, xhigh, on Codex), record 6247306813: **material**,
five findings. F1 the parsed bytes were not tied to the hash `artifacts` verified; F2 a
colon in a finding rule had no effective test; F3 the PRD cells were put in the wrong
tables; F4 contradictory ownership clauses; F5 this file was absent. All five were fixed
at 67f1d1c (F1 and F2 test first, red recorded) and the head was frozen again.
Round 2 on b8a549b, record 6247681138: **nothing material in scope**
(the re-read that a fix requires, after the cap-1 round). Note N1, an acceptance row that read a parse as a pass although a rejected
plan review parses and is false, was fixed in a notes commit with no round. Main
(#365) was merged cleanly with no conflict. Harness exception: ADR-0026 names
`claude`, `devin` and `opencode`; Codex was used on the Operator's message of
2026-10-13, for both rounds.

Complete landing figure against main is recorded on the issue at the merge.

## Resource record

ADR-0007, recorded, not budgeted. Tokens are as the harness reported them.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan (revised) | reasoning | SWE-2 Max | not recorded | not recorded | not recorded |
| Plan review | reasoning | GPT-6.1 Sol | xhigh | 184,567 | 2026-10-18T05:58:41Z to 06:04:44Z |
| Implementation | execution | Claude Sonnet 5.5 | not recorded | 250,379 | 92 min (to 07:35Z) |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 195,163 | 07:39:01Z to 07:48:21Z |
| Fix of round 1 | execution | Claude Sonnet 5.5 | not recorded | not recorded | about 35 min |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 170,929 | 08:24:30Z to 08:33:08Z |
| Close-out | execution | Claude Sonnet 5.5 | not recorded | not recorded | about 5 min |
| Orchestration | reasoning | Claude Opus 5.5 | not recorded | not separately recorded | not separately recorded |
| Whole integration level (root) | deterministic | | | | 29 min (fix head) |
