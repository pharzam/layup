# T-b3s3 — test runs and compiled mutants (#401)

Base `9d8bb9a` (main with #361). Author Sonnet 5.5. Times are from `date -u` after each step.

## Red, on the unchanged helper (before any production edit)

The three wrong-role fixtures and the `Position` pair were written first; the red run used a scratch copy of the
base plus those test files (the one test that calls the new API was given the old three-argument call there), so
no production file was edited. Command:
`go test -tags=integration ./internal/run -run 'TestAnAnswerOfADeveloperSessionIsNotAcceptedFromTheRecords|TestAnAnswerOfAWrongRoleSessionLeavesThePositionUphill|TestThePlanAndTheRecordsGiveTheTierOfARetry'`
Direct exit 1, 61.55 s. The failures are the wanted ones (the helper accepts the wrong role); the owner-role arms
and the existing test pass, which are the controls against over-refusal:

```
--- FAIL: TestAnAnswerOfADeveloperSessionIsNotAcceptedFromTheRecords (7.13s)
    handoffpost_integration_test.go:404: the answer of a developer session is accepted by S-bbd75b5c
--- FAIL: TestAnAnswerOfAWrongRoleSessionLeavesThePositionUphill (45.65s)
    promptfile_integration_test.go:209: a developer: downhill, <nil>; want uphill
    promptfile_integration_test.go:209: the owner of another kind: downhill, <nil>; want uphill
    promptfile_integration_test.go:209: one of two owners of one session's questions: downhill, <nil>; want uphill
FAIL
FAIL	github.com/pharzam/layup/internal/run	61.550s
FAIL
```

## Green, on the change (focused selectors, direct exit)

- `go build ./...`, `go vet ./...`, `go vet -tags=integration ./internal/run`, `go test ./... -count=1`: exit 0.
- `go test -tags=integration ./internal/run -run '<the four names>'` (the two new tests, `TestThePlanAndTheRecordsGiveTheTierOfARetry`
  with a real `software-architect` owner session, `TestTheQuestionOfASessionIsRoutedAndItsAnswerAcceptedFromTheRecords`): exit 0, 42.1 s.
- Unit: `TestAnAnswerIsAcceptedByTheNextAttemptThatCitesIt`, `TestAnAnswerIsNotAcceptedUnlessEachConditionHolds` (its 20 cases unchanged),
  `TestAnAnswerOfASessionOutsideTheOwnerRoleIsNotAccepted` (20 cases: the mixed shapes across two sessions and inside one, with the
  same-owner positive controls), `TestTwoLedgersThatDifferOnlyInTheRoleOfTheAnsweringSessionAreToldApart`,
  `TestTheOwnerRoleIsReadFromTheQuestionsAndTheOwnerMapGiven`, `TestOpenQuestionsAreTheQuestionsNoAcceptedAnswerClosed`,
  `TestAWrongRoleAnswerClosesNoOpenQuestion`: exit 0.
- Each commit passed the hook (adr-lint, prd-lint, link-lint, discipline tests, setup-check, `go test ./...`).

No full-module or unfiltered integration run was made (the root serializes them); CI is the authority for the whole ladder.

## Compiled mutants

Each mutant is a real edit of `internal/run/answer.go` or `tier.go` by one `sed` expression, built and run against the
named selector, then restored from a copy; the restore is checked by the SHA-256 of both files and `git diff --exit-code`
(`restored byte-exact: yes`). Unit selector `Answer|OpenQuestions|Ledgers|OwnerRole`; integration selector the four names
of the runs above. "Killed" is a direct exit 1 with a failing assertion, or a compile error where stated.

| # | Mutant | Result |
| - | ------ | ------ |
| M1 | the role comparison made always true | killed: the outside-role matrix, the pair test, `TestAWrongRoleAnswerClosesNoOpenQuestion` |
| M2 | compare the display name, not `RoutingRole` | killed: the acceptance tests (the owner role refused) |
| M3 | read the role under the answer ID, not the session | killed: the acceptance tests |
| M4 | route with `DefaultOwners()`, not the map given | killed: the replaced-map cases, `TestTheOwnerRoleIsReadFromTheQuestionsAndTheOwnerMapGiven` |
| M5 | accept questions of more than one owner | killed: the mixed shapes (across sessions and inside one session) |
| M7 | accept an attempt with no question | killed: `TestTheOwnerRoleIsReadFromTheQuestionsAndTheOwnerMapGiven` |
| M8 | accept the role `operator` (a human question, a map to the operator) | killed: the human case answered by a session of role `operator` |
| M9 | read the questions of every attempt | killed: `TestTheOwnerRoleIsReadFromTheQuestionsAndTheOwnerMapGiven` |
| M10 | read the first question only | killed: the mixed shapes |
| M11 | swap `roles` and `owners` at the call in `OpenQuestions` | killed at compile time: the named types `Roles` and `Owners` differ |
| M12 | `Position` feeds an empty roles map | killed: the `Position` pair (the owner arm), `TestThePlanAndTheRecordsGiveTheTierOfARetry` |
| M13 | `sessionRoles` reads the task column | killed: `TestTheOwnerRoleIsReadFromTheQuestionsAndTheOwnerMapGiven` |
| M14 | the same, run through the records | killed: the `Position` pair, the two existing integration tests |

**Credit, exactly.** The table lists 13 rows: **12 valid compiled mutants killed by an assertion** (M1 to M5, M7 to M10,
M12 to M14) and **1 deliberate type-check rejection** (M11, which does not compile, so no assertion ran; it shows that
the named types do their job and is not counted as an assertion kill).

**Attempts disclosed separately, not credited.** M6 (drop the separate human check) was an equivalent mutant: a human
question's route role is `operator`, which M8 covers, so the separate check was removed from the code instead of kept.
M1 first failed to build (an unused variable), an invalid mutant, and was re-cut to compile; M12 failed to build the
same way in the first run and was re-cut. The first run of all mutants (in `/tmp`, before the human check was removed)
is not kept; this table is the second.
