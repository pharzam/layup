# T-gyvo — test runs and compiled mutants (#405)

Base `0d85e86` (main with #401 and #403). Author Sonnet 5.5. Times are `date -u` after each step (2026-10-18, UTC); every status is the direct exit of the command, never through a pipe.

## Step 1: compiling stubs, committed alone (46ffe25)

`QuestionComment` and `AnswerComment` return `(nil, errNotBuilt)`; `firstLinePrefix` returns `""` and is not yet called. `go build ./...`, `go vet ./...`, `go vet -tags=integration ./internal/run` exit 0; the pre-commit hook (`go test ./...`, lints) exit 0, 07:44:17.

## Step 2: the tests, red on the stubs

Unit, `go test -count=1 ./internal/run -run '<the eight names>'`: exit 1 at 07:45:58. All eight tests fail, 72 subtests of `TestAMalformedRowGivesAnErrorAndNoBytes` among them; each fails on a semantic assertion (wanted bytes, or the valid control of the subtest does not compose), none on an undefined symbol. No refusal-only subtest passes on the stub, because each carries its valid control. Excerpt:

```
--- FAIL: TestTheQuestionCommentIsItsFirstLineABlankLineAndTheText (0.00s)
    qacomment_test.go:64: got "", not built; want "S-0a1b2c3d: developer session of T-ab12, attempt 2: question S-0a1b2c3d.3\n\nWhich owner?"
--- FAIL: TestTheFirstLinePrefixIsTheOneOfTheEndAndTheComposers (0.00s)
    qacomment_test.go:104: firstLinePrefix ""; want "S-0a1b2c3d: developer session of T-ab12, attempt 2: "
--- FAIL: TestAnAnswersQuestionIDIsCheckedByItsShapeAndNeverResolved (0.00s)
    qacomment_test.go:123: question ID "S-0a1b2c3d.3": "", not built; want it composed
--- FAIL: TestAMalformedRowGivesAnErrorAndNoBytes/kind_of_question_answer (0.00s)
    qacomment_test.go:181: the control "question" does not compose: not built
--- FAIL: TestTheSameRowsGiveTheSameBytesAndTheInputsAreNotChanged / TestABodyThat... / TestTheIDUsesTheColumnN... (same form)
```

Integration, `go test -count=1 -tags=integration ./internal/run -run 'TestTheQuestionAndAnswerComments|TestARowOfARealResult'`: exit 1 at 07:46:41, 17.5 s. Both fail on the semantic assertion (the demo: four byte comparisons against literals written with the real session ID; the refusal test: the control question row does not compose).

One test was corrected after its red run: the expectation of the field-by-field change in `TestTheFirstLinePrefix...` compared a replaced substring ("2" occurs in the session ID too); it now builds each expected line from its four fields. The red of that test (`firstLinePrefix ""`) was on its first assertion, which did not change.

## Step 3: green (2d9d3fb)

- `go test -count=1 ./internal/run` exit 0 at 07:47:39 (unit, after the correction above).
- Integration, the two new tests and the five baseline tests of the plan (`TestTheResultIsCommittedWithItsTelemetryRow`, `TestTheCommentNamesARefusalOfThePushHook`, `TestTheCommentOfASecondSessionIsItsOwn`, `TestAnArtifactThatDiffersIsRefusedAndKept`, `TestACommentBodyIsRefusedBeforeThePost`): exit 0 at 07:48:09, 27.7 s. The baseline five are unedited and green.
- The commit passed the hook, 07:48:49.

## Step 4: compiled mutants, on a copy of the tree

The harness copies `internal/`, `cmd/` and `go.mod` to a scratch directory, applies one textual mutant to the copy (the worktree is never mutated), requires `go build ./...`, `go vet ./internal/run` and `go vet -tags=integration ./internal/run` to exit 0 before the run (a mutant that does not build is "not run", never "killed"), runs the unit tests (or, for the `end.go` mutant, the end tests and the demo at the integration level), reads the `--- FAIL` lines, restores the file and compares its SHA-256 with the file of the commit. The control after the last restore is `go test -count=1 ./internal/run`, exit 0. 21 mutants, 21 built, 21 killed, 0 survivors, 21 restored hash-equal. Run 07:49 to 07:57.

| Id | Time | Mutant | File | Result | Level | Killed by | Restore |
| -- | ---- | ------ | ---- | ------ | ----- | --------- | ------- |
| M1 | 07:49:29 | ID from the text row[2] | internal/run/qacomment.go | killed | unit | TestABodyThatCheckCommentBodyRefusesGivesAnErrorAndNoBytes, TestAMalformedRowGivesAnErrorAndNoBytes, TestAnAnswersQuestionIDIsCheckedByItsShapeAndNeverResolved, TestTheAnswerCommentNamesTheAnswerAndTheQuestionItAnswers, TestTheFirstLinePrefixIsTheOneOfTheEndAndTheComposers, TestTheIDUsesTheColumnNNeverTheTextNorTheReason, TestTheQuestionCommentIsItsFirstLineABlankLineAndTheText, TestTheSameRowsGiveTheSameBytesAndTheInputsAreNotChanged | restored, hash equal |
| M2 | 07:49:49 | ID from the reason row[4] | internal/run/qacomment.go | killed | unit | TestABodyThatCheckCommentBodyRefusesGivesAnErrorAndNoBytes, TestAMalformedRowGivesAnErrorAndNoBytes, TestAnAnswersQuestionIDIsCheckedByItsShapeAndNeverResolved, TestTheAnswerCommentNamesTheAnswerAndTheQuestionItAnswers, TestTheFirstLinePrefixIsTheOneOfTheEndAndTheComposers, TestTheIDUsesTheColumnNNeverTheTextNorTheReason, TestTheQuestionCommentIsItsFirstLineABlankLineAndTheText, TestTheSameRowsGiveTheSameBytesAndTheInputsAreNotChanged | restored, hash equal |
| M3 | 07:50:09 | separator \n | internal/run/qacomment.go | killed | unit | TestTheAnswerCommentNamesTheAnswerAndTheQuestionItAnswers, TestTheIDUsesTheColumnNNeverTheTextNorTheReason, TestTheQuestionCommentIsItsFirstLineABlankLineAndTheText, TestTheSameRowsGiveTheSameBytesAndTheInputsAreNotChanged | restored, hash equal |
| M4 | 07:50:31 | final \n appended | internal/run/qacomment.go | killed | unit | TestTheAnswerCommentNamesTheAnswerAndTheQuestionItAnswers, TestTheIDUsesTheColumnNNeverTheTextNorTheReason, TestTheQuestionCommentIsItsFirstLineABlankLineAndTheText, TestTheSameRowsGiveTheSameBytesAndTheInputsAreNotChanged | restored, hash equal |
| M5 | 07:50:47 | composer builds its own line, task and role swapped | internal/run/qacomment.go | killed | unit | TestTheAnswerCommentNamesTheAnswerAndTheQuestionItAnswers, TestTheFirstLinePrefixIsTheOneOfTheEndAndTheComposers, TestTheQuestionCommentIsItsFirstLineABlankLineAndTheText, TestTheSameRowsGiveTheSameBytesAndTheInputsAreNotChanged | restored, hash equal |
| M6 | 07:51:26 | End comment drifts from the home | internal/run/end.go | killed | integration | TestAnArtifactThatDiffersIsRefusedAndKept, TestTheCommentNamesARefusalOfThePushHook, TestTheCommentOfASecondSessionIsItsOwn, TestTheQuestionAndAnswerCommentsAreComposedFromTheRowsOfARealEnd, TestTheResultIsCommittedWithItsTelemetryRow | restored, hash equal |
| M7 | 07:51:46 | CheckCommentBody not called | internal/run/qacomment.go | killed | unit | TestABodyThatCheckCommentBodyRefusesGivesAnErrorAndNoBytes | restored, hash equal |
| M8 | 07:52:05 | limit in bytes | internal/run/qacomment.go | killed | unit | TestABodyThatCheckCommentBodyRefusesGivesAnErrorAndNoBytes | restored, hash equal |
| M9 | 07:52:20 | questionID shape check dropped | internal/run/qacomment.go | killed | unit | TestAnAnswersQuestionIDIsCheckedByItsShapeAndNeverResolved | restored, hash equal |
| M10 | 07:52:34 | questionID checked by prefix S- only | internal/run/qacomment.go | killed | unit | TestAnAnswersQuestionIDIsCheckedByItsShapeAndNeverResolved | restored, hash equal |
| M11 | 07:52:48 | questionID from the first ID word of the reason, parameter ignored | internal/run/qacomment.go | killed | unit | TestTheAnswerCommentNamesTheAnswerAndTheQuestionItAnswers, TestTheIDUsesTheColumnNNeverTheTextNorTheReason | restored, hash equal |
| M12 | 07:53:17 | CheckResultRow not called | internal/run/qacomment.go | killed | unit | TestAMalformedRowGivesAnErrorAndNoBytes | restored, hash equal |
| M13 | 07:53:49 | CheckSession not called | internal/run/qacomment.go | killed | unit | TestAMalformedRowGivesAnErrorAndNoBytes | restored, hash equal |
| M14 | 07:54:20 | composed-ID answerForm check dropped | internal/run/qacomment.go | killed | unit | TestAMalformedRowGivesAnErrorAndNoBytes | restored, hash equal |
| M15 | 07:54:51 | row kind not required to be the composer's | internal/run/qacomment.go | killed | unit | TestAMalformedRowGivesAnErrorAndNoBytes | restored, hash equal |
| M16 | 07:55:23 | reader schema and type check dropped | internal/run/qacomment.go | killed | unit | TestAMalformedRowGivesAnErrorAndNoBytes | restored, hash equal |
| M17 | 07:55:55 | error returned with the bytes | internal/run/qacomment.go | killed | unit | TestABodyThatCheckCommentBodyRefusesGivesAnErrorAndNoBytes | restored, hash equal |
| M18 | 07:56:21 | answer IDs in the other order | internal/run/qacomment.go | killed | unit | TestAnAnswersQuestionIDIsCheckedByItsShapeAndNeverResolved, TestTheAnswerCommentNamesTheAnswerAndTheQuestionItAnswers, TestTheIDUsesTheColumnNNeverTheTextNorTheReason | restored, hash equal |
| M19 | 07:56:45 | no-repair equality dropped | internal/run/qacomment.go | killed | unit | TestAMalformedRowGivesAnErrorAndNoBytes | restored, hash equal |
| M20 | 07:57:10 | HandoffRefusal rewrapped as a plain error | internal/run/qacomment.go | killed | unit | TestABodyThatCheckCommentBodyRefusesGivesAnErrorAndNoBytes | restored, hash equal |
| M21 | 07:57:34 | column count check dropped | internal/run/qacomment.go | killed | unit | TestAMalformedRowGivesAnErrorAndNoBytes | restored, hash equal |

The reader-boundary mutants (M12, M13, M14, M15, M16, M19, M21) and the body mutants (M7, M8, M17, M20) are each killed by the one test that names the rule they remove; M1, M2, M5 and M18 by several. No survivor, so no test was added after the sweep.

## Step 5: gates

All statuses are direct. Times are `date -u` after each command.

- `go build ./...`, `go vet ./...`, `go vet -tags=e2e ./cmd/layup/`, `go vet -tags=integration ./...`: exit 0 (08:02). `gofmt -l .`: empty.
- `go test ./...` in the hook of each commit: exit 0.
- `go test -count=1 -tags=integration ./cmd/layup/` (`TestPackageRules` after the edit of `packages.md`): exit 0, 08:02:40.
- The whole integration level, `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...`, started 07:57:59 on the head 2d9d3fb (the document commit fc5ab73 changes no Go file; `cmd/layup` was re-run above after it): exit 0 at 08:24:19, about 26 min on the loaded host; no package FAIL.
- `sh docs/adr/adr-lint.sh`, `sh docs/prd/prd-lint.sh`, `sh docs/links/link-lint.sh`, `sh docs/tests/run-discipline-tests.sh` (138 passed), `sh docs/setup/setup-check.sh` (after `git add`), `git diff --check origin/main HEAD`: exit 0 (08:03).
- E2E stays CI-only (`go vet -tags=e2e ./cmd/layup/` exit 0 is the local proof it compiles).
