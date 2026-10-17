# T-b3s3 — compiled mutants (#401)

Each mutant is a real edit of `internal/run/answer.go` or `tier.go` by one `sed` expression, built and run against the
named selector, then restored from a copy; the restore is checked by the SHA-256 of both files and `git diff --exit-code`
(`restored byte-exact: yes`). Unit selector `Answer|OpenQuestions|Ledgers|OwnerRole`; integration selector the four names
of `test-runs.md`. "Killed" is a direct exit 1 with a failing assertion, or a compile error where stated.

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

M6 (drop a separate human check) was an equivalent mutant: a human question's route role is `operator`, which M8
covers, so the separate check was removed from the code instead of kept; M1 first failed to build (an unused variable)
and was re-cut to compile. The first run of all mutants, `/tmp` only, is not kept; this table is the second.
