# T-4pje test runs

## Red

The named tests against the compiling stubs (6fccd51): unit run of ./internal/run at 2026-10-18T06:19Z, integration run (-tags=integration) at 06:20:34Z. Each test failed on a semantic assertion, none on a compile error; the lines show the first two failing assertions of each.

unit FAIL TestThePlanCommentParsesInThePinnedForm (4 failing assertions)
    commentform_test.go:55: ParsePlanComment = {Size: OpenQuestions:[] NewInterfaces:[] NewDependencies:[]}, [], false; want small, no lists, [{t1 internal/a_test.go} {t2 internal/b_test.go}]
    commentform_test.go:60: a plan in another order = {Size: OpenQuestions:[] NewInterfaces:[] NewDependencies:[]}, [], false

unit FAIL TestAPlanFieldTwiceOrAbsentIsNoPlan (1 failing assertions)
    commentform_test.go:78: the control: false for "size: small\nopen questions: —\nnew interfaces: —\nnew dependencies: —\ntests:\n- t1 internal/a_test.go\n- t2 internal/b_test.go\n"; want true

unit FAIL TestTheTestsMarkerOnceWithBulletsOfTwoWords (5 failing assertions)
    commentform_test.go:83: two words: false for "size: small\nopen questions: —\nnew interfaces: —\nnew dependencies: —\ntests:\n- t1 a.go\n"; want true
    commentform_test.go:83: the same path twice: false for "size: small\nopen questions: —\nnew interfaces: —\nnew dependencies: —\ntests:\n- t1 a.go\n- t2 a.go\n"; want true

unit FAIL TestAPlanListFieldIsOneLineOfText (3 failing assertions)
    commentform_test.go:107: a text with a comma = [], false; want the one element "a, b"
    commentform_test.go:110: — = [], false; want the empty list and the execution tier

unit FAIL TestAFieldValueIgnoresOneCarriageReturn (4 failing assertions)
    commentform_test.go:118: one \r after the size: false for "size: small\r\nopen questions: —\nnew interfaces: —\nnew dependencies: —\ntests:\n- t1 internal/a_test.go\n- t2 internal/b_test.go\n"; want true
    commentform_test.go:118: trailing spaces after the size: false for "size: small  \nopen questions: —\nnew interfaces: —\nnew dependencies: —\ntests:\n- t1 internal/a_test.go\n- t2 internal/b_test.go\n"; want tr

unit FAIL TestAMarkerOrBulletWithACarriageReturnIsNotRead (2 failing assertions)
    commentform_test.go:132: a test bullet, the control: false for "size: small\nopen questions: —\nnew interfaces: —\nnew dependencies: —\ntests:\n- t1 internal/a_test.go\n- t2 internal/b_test.go\n"; want true
    commentform_test.go:137: a finding, the control: false for "head: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\nverdict: material\nfindings:\n- internal/a.go:1 r: text\n"; want true

unit FAIL TestAFieldWithNothingAfterTheColonIsMalformed (2 failing assertions)
    commentform_test.go:146: the control: false for "size: small\nopen questions: —\nnew interfaces: —\nnew dependencies: —\ntests:\n- t1 internal/a_test.go\n- t2 internal/b_test.go\n"; want true
    commentform_test.go:152: a verdict, the control: false for "verdict: approve\n"; want true

unit FAIL TestTheLineReaderIgnoresAFencedBlock (10 failing assertions)
    commentform_test.go:156: a second size in a backtick fence: false for "size: small\nopen questions: —\nnew interfaces: —\nnew dependencies: —\ntests:\n- t1 internal/a_test.go\n- t2 internal/b_test.go\n```\nsize
    commentform_test.go:156: a second size in a tilde fence: false for "size: small\nopen questions: —\nnew interfaces: —\nnew dependencies: —\ntests:\n- t1 internal/a_test.go\n- t2 internal/b_test.go\n~~~\nsize: l

unit FAIL TestACommentThatFailsItsFormIsFalseEverywhere (1 failing assertions)
    commentform_test.go:179: the control record is not valid

unit FAIL TestVerdictFieldIsApproveOrApproveWithConditions (5 failing assertions)
    commentform_test.go:197: approve: false for "verdict: approve\n"; want true
    commentform_test.go:197: approve-with-conditions: false for "x\nverdict: approve-with-conditions\n"; want true

unit FAIL TestFindingsLocatedNeedsALocatedFinding (5 failing assertions)
    commentform_test.go:216: one finding: false for "head: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\nverdict: material\nfindings:\n- internal/a.go:12 r-x: a text: with a colon\n"; want true
    commentform_test.go:216: line 1: false for "head: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\nverdict: material\nfindings:\n- p.go:1 r: t\n"; want true

unit FAIL TestNothingMaterialNeedsTheVerdictTheHeadAndAnEmptyBlock (2 failing assertions)
    commentform_test.go:268: the control: false; want true
    commentform_test.go:268: 64 hexadecimal characters: false; want true

unit FAIL TestOpenMaterialFailsUphill (8 failing assertions)
    commentform_test.go:289: nothing material and a finding: false; want true
    commentform_test.go:289: material and an empty block: false; want true

unit FAIL TestAnEmptyOrInvalidCommentParsesToNothing (1 failing assertions)
    commentform_test.go:314: the control plan is not read

unit FAIL TestTheCommentItemsAreComputedFromTheKeptCopy (16 failing assertions)
    handoffpost_test.go:628: plan-fields with "size: small\r\nopen questions: —\nnew interfaces: —\nnew dependencies: —\ntests:\n# The plan\r\n\r\nno final newline, trailing  \t" and map[]: {Valid:false Detail:inva
    handoffpost_test.go:628: plan-fields with "size: small\r\nopen questions: —\nnew interfaces: —\nnew dependencies: —\ntests:\n# The plan\r\n\r\nno final newline, trailing  \t" and map[plan-fields:false]: {Valid:

unit FAIL TestAnItemOfAnotherCommentKindIsFalse (1 failing assertions)
    handoffpost_test.go:666: a plan copy of the kind plan: {Valid:false Detail:invalid plan-fields comment 701-1}, <nil>; want valid comment 701-1

unit FAIL TestNothingMaterialReadsTheBaseOfTheVerifierStartRow (2 failing assertions)
    handoffpost_test.go:704: the head is the base: {Valid:false Detail:invalid nothing-material comment 701-1}, <nil>; want valid comment 701-1
    handoffpost_test.go:704: the head is not the base: {Valid:true Detail:valid comment 701-1}, <nil>; want invalid nothing-material comment 701-1

unit FAIL TestACodeWrittenHandoffHasNoCommentItem (4 failing assertions)
    handoffpost_test.go:736: plan-fields of a handoff that code wrote: {Valid:true Detail:valid}, <nil>; want invalid plan-fields
    handoffpost_test.go:736: verdict-field of a handoff that code wrote: {Valid:true Detail:valid}, <nil>; want invalid verdict-field

unit FAIL TestAKeptCopyTheStoreLacksOrCannotReadMakesTheItemsFalse (12 failing assertions)
    handoffpost_test.go:756: absent copy, plan-fields: {Valid:true Detail:valid comment 701-1}, <nil>; want invalid plan-fields and the event written
    handoffpost_test.go:756: absent copy, verdict-field: {Valid:true Detail:valid comment 701-1}, <nil>; want invalid verdict-field and the event written

unit FAIL TestAStoreErrorReturnsAnErrorAndWritesNoHandoff (1 failing assertions)
    handoffpost_test.go:781: a store that fails: {Valid:true Detail:valid comment 701-1}, <nil>; 1 events, 1 commits; want an error and no write

integration FAIL TestTheCommentChecksOfTheTaskLoop (1 failing assertions)
    handoffpost_integration_test.go:456: RecordHandoff = {Valid:true Detail:valid comment 602-1}, <nil>; want invalid verdict-field comment 602-1

## Green

Green run of the named unit tests on the implementation, 2026-10-18T06:27:03Z

```
--- PASS: TestThePlanCommentParsesInThePinnedForm (0.00s)
--- PASS: TestAPlanFieldTwiceOrAbsentIsNoPlan (0.00s)
--- PASS: TestTheTestsMarkerOnceWithBulletsOfTwoWords (0.00s)
--- PASS: TestAPlanListFieldIsOneLineOfText (0.00s)
--- PASS: TestAFieldValueIgnoresOneCarriageReturn (0.00s)
--- PASS: TestAMarkerOrBulletWithACarriageReturnIsNotRead (0.00s)
--- PASS: TestAFieldWithNothingAfterTheColonIsMalformed (0.00s)
--- PASS: TestTheLineReaderIgnoresAFencedBlock (0.00s)
--- PASS: TestACommentThatFailsItsFormIsFalseEverywhere (0.00s)
--- PASS: TestVerdictFieldIsApproveOrApproveWithConditions (0.00s)
--- PASS: TestFindingsLocatedNeedsALocatedFinding (0.00s)
--- PASS: TestNothingMaterialNeedsTheVerdictTheHeadAndAnEmptyBlock (0.00s)
--- PASS: TestOpenMaterialFailsUphill (0.00s)
--- PASS: TestAnEmptyOrInvalidCommentParsesToNothing (0.00s)
--- PASS: TestTheCommentItemsAreComputedFromTheKeptCopy (0.04s)
--- PASS: TestAnItemOfAnotherCommentKindIsFalse (0.01s)
--- PASS: TestNothingMaterialReadsTheBaseOfTheVerifierStartRow (0.00s)
--- PASS: TestACodeWrittenHandoffHasNoCommentItem (0.00s)
--- PASS: TestAKeptCopyTheStoreLacksOrCannotReadMakesTheItemsFalse (0.03s)
--- PASS: TestAStoreErrorReturnsAnErrorAndWritesNoHandoff (0.00s)
PASS
ok  	github.com/pharzam/layup/internal/run	0.112s
```

## Demo

Demo of row 66l: TestTheCommentChecksOfTheTaskLoop (integration, real git, the httptest forge), 2026-10-18T06:23Z

```
--- PASS: TestTheCommentChecksOfTheTaskLoop (34.22s)
    --- PASS: TestTheCommentChecksOfTheTaskLoop/nothing_material_at_the_verifier's_base,_the_caller_says_false (4.87s)
    --- PASS: TestTheCommentChecksOfTheTaskLoop/nothing_material_at_another_head,_the_caller_says_true (5.88s)
    --- PASS: TestTheCommentChecksOfTheTaskLoop/nothing_material_with_a_finding,_the_caller_says_true (5.13s)
    --- PASS: TestTheCommentChecksOfTheTaskLoop/findings_located,_the_caller_says_false (4.63s)
    --- PASS: TestTheCommentChecksOfTheTaskLoop/a_plan_of_the_form,_the_caller_says_false (4.42s)
    --- PASS: TestTheCommentChecksOfTheTaskLoop/a_plan_that_is_not_of_the_form,_the_caller_says_true (4.68s)
    --- PASS: TestTheCommentChecksOfTheTaskLoop/a_plan_review_that_rejects,_the_caller_says_true (4.61s)
ok  	github.com/pharzam/layup/internal/run	61.716s
```

## Mutants

Each mutant changes one expression of `internal/run/commentform.go` or `handoffpost.go` in a copy, must pass `go build ./...` and `go vet ./internal/run` (a mutant that fails either is not run and is excluded), and is killed when `go test ./internal/run` fails a test. The copy is restored byte-exact after each (checked) and the control is green before and after.

| Mutant | Result | Failing tests (first four) |
| --- | --- | --- |
| field-\r unstripped | killed | TestAFieldValueIgnoresOneCarriageReturn, TestAnEditedCommentIsANewCopyAndTheFirstHandoffIsInvalid, TestAnItemOfAnotherCommentKindIsFalse, TestTheCommentItemsAreComputedFromTheKeptCopy +1 |
| trailing spaces unstripped | killed | TestAFieldValueIgnoresOneCarriageReturn, TestAFieldWithNothingAfterTheColonIsMalformed |
| tests:\r accepted as marker | killed | TestAMarkerOrBulletWithACarriageReturnIsNotRead |
| tests: x accepted as marker | killed | TestAMarkerOrBulletWithACarriageReturnIsNotRead, TestTheTestsMarkerOnceWithBulletsOfTwoWords |
| plan bullet \r accepted | killed | TestAMarkerOrBulletWithACarriageReturnIsNotRead |
| finding bullet \r accepted | killed | TestAMarkerOrBulletWithACarriageReturnIsNotRead |
| fence skip dropped | killed | TestTheLineReaderIgnoresAFencedBlock, TestVerdictFieldIsApproveOrApproveWithConditions |
| CRLF closer accepted | killed | TestTheLineReaderIgnoresAFencedBlock |
| closer of any count | killed | TestTheLineReaderIgnoresAFencedBlock |
| closer of any character | killed | TestTheLineReaderIgnoresAFencedBlock |
| verdict-field accepts reject | killed | TestTheCommentItemsAreComputedFromTheKeptCopy, TestVerdictFieldIsApproveOrApproveWithConditions |
| line >= 0 | killed | TestFindingsLocatedNeedsALocatedFinding |
| tests bullet of two or more words | killed | TestTheTestsMarkerOnceWithBulletsOfTwoWords |
| duplicate test IDs | killed | TestTheTestsMarkerOnceWithBulletsOfTwoWords |
| test path not checked | killed | TestTheTestsMarkerOnceWithBulletsOfTwoWords |
| empty field value read | killed | TestAFieldWithNothingAfterTheColonIsMalformed |
| field twice read | killed | TestAPlanFieldTwiceOrAbsentIsNoPlan |
| marker twice read | killed | TestAPlanFieldTwiceOrAbsentIsNoPlan, TestFindingsLocatedNeedsALocatedFinding |
| lists split on commas | killed | TestAPlanListFieldIsOneLineOfText |
| size unchecked | killed | TestAFieldValueIgnoresOneCarriageReturn, TestThePlanCommentParsesInThePinnedForm |
| invalid UTF-8 read | killed | TestAnEmptyOrInvalidCommentParsesToNothing |
| head not checked (non-strict parse) | killed | TestACommentThatFailsItsFormIsFalseEverywhere, TestFindingsLocatedNeedsALocatedFinding, TestNothingMaterialNeedsTheVerdictTheHeadAndAnEmptyBlock |
| finding rule of two words | killed | TestFindingsLocatedNeedsALocatedFinding |
| nothing-material ignores the head | killed | TestNothingMaterialNeedsTheVerdictTheHeadAndAnEmptyBlock, TestNothingMaterialReadsTheBaseOfTheVerifierStartRow, TestTheCommentItemsAreComputedFromTheKeptCopy |
| nothing-material ignores the block | killed | TestNothingMaterialNeedsTheVerdictTheHeadAndAnEmptyBlock |
| findings-located accepts an empty block | killed | TestFindingsLocatedNeedsALocatedFinding, TestTheCommentItemsAreComputedFromTheKeptCopy |
| OpenMaterial = findings only | killed | TestOpenMaterialFailsUphill |
| head compared to another column | killed | TestNothingMaterialReadsTheBaseOfTheVerifierStartRow |
| kind dispatch removed | killed | TestAnItemOfAnotherCommentKindIsFalse |
| caller literal honored (reset removed) | killed | TestACodeWrittenHandoffHasNoCommentItem, TestAKeptCopyTheStoreLacksOrCannotReadMakesTheItemsFalse |
| missing kept copy read as valid | killed | TestAKeptCopyTheStoreLacksOrCannotReadMakesTheItemsFalse |
| store error swallowed | killed | TestAStoreErrorReturnsAnErrorAndWritesNoHandoff |
| code-written handoff reads a copy | killed | TestACodeWrittenHandoffHasNoCommentItem |
| caller literal false honored (only the true literal reaches plan-fields) | killed | TestTheCommentItemsAreComputedFromTheKeptCopy, TestAnItemOfAnotherCommentKindIsFalse |

34 mutants, 34 killed, 0 excluded by a build or vet failure. Control after restore: green, files byte-identical: True.

## Gates at the head 8b403ff (each exit status read directly)

| Gate | Exit | Time (UTC, 2026-10-18) |
| --- | --- | --- |
| `go build ./...`, `go vet ./...`, `gofmt -l .`, `go vet -tags=e2e ./cmd/layup/` | 0, 0, none listed, 0 | 07:16 |
| `go test -count=1 ./...` (unit) | 0 | 07:35:30 to 07:35:40 |
| `go test -count=1 -tags=integration -timeout 30m ./...` (whole level, root) | 0 (`internal/run` 1,146 s on a loaded host) | 07:15:21 to 07:35:19 |
| the five text checks, `nested-checkout-check.sh`, `git diff --check` | 0 | 07:16 |

An earlier whole-level run at the head cfd4519 failed: the tests of the authors, the frozen list and the prompt file wrote comments of no form (fixed in 8b403ff, then rerun green above), and its 20-minute timeout fell on a host under load.
