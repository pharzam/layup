# T-xvvd — the pull request and the merge

Issue [#363](https://github.com/pharzam/layup/issues/363), row 66g of M2e. Serves REQ-007, NFR-004, NFR-001
(F-0003#47, F-0003#58). Plan 6249038254 (Opus 5.5); plan review 6249142236 (GPT-6.1 Sol, approve-with-conditions;
condition 1 became the size child T-08rj, #413, row 66q; conditions 2 to 9 met). Maximum set by the review:
3,000 lines over 31 files, cap 1. Evidence: [test runs](../../runs/T-xvvd/test-runs.md) and [mutants](../../runs/T-xvvd/mutants.md).

## Delivered and boundaries
`internal/run` opens the draft pull request of a bound head (`OpenPullRequest`, `PullBody`; the title is the form
of the checked-in `pr-title` check, without the task ID), reads the ready rule (`Ready`: the newest status of each of the
four contexts of any creator is the App's `success`; each required job of GitHub Actions has runs that are
`completed` and `success`), marks ready and merges at the head SHA it read (`MarkReady`, `MergeTask`: the statuses are
read again before the merge, and a 409 reads the pull request and the new head's statuses and sends no second merge),
writes the event `merged` (`RecordMerged`), checks the close-out commit (`PushCloseOut`: `scope`, `verdict`),
the rendered commit (`CheckRendered`) and carries `layup/verify` over them (`CarryVerify`), and `forge-check`
accepts a merge not yet recorded (`explainHead`, `mergeMaps`, `treeSeams`). `setup.ActionsApp` is exported.
Not here: the edit of the body (66q), the step loop, the redo and the clean-merge carry-over (66k), the push of the
rendered commit (66j), `docs/tasks/backlog.md` in the close-out (the Operator's decision).

## Tests mapped to the criteria

| Criterion | Test |
| --- | --- |
| The pull request, title, body, register | `TestTheTitleIsTheFormOfThePRTitleCheck`, `TestThePullBodyIsTheTwoLinesAndTakesFurtherLines`, `TestTheDraftIsOpenedOnceWithItsTitleAndBody`, `TestAStopAfterTheDraftOpensNoSecondDraft`, `TestTheRecordsCommitIsTheBoundSessionsAndNotALaterOnes`, `TestTheDraftPullRequestOfABoundHead` |
| The ready and merge rule | `TestTheReadyRuleNamesTheFirstContextOrJobThatIsNotSuccess`, `TestTheRequiredJobsAreTheActionsEntriesOfTheRequiredChecksRules`, `TestTheReadyAndMergeRuleBlocksEachCaseAndNothingIsMarkedOrMerged`, `TestAStatusThatChangesAfterTheMarkBlocksTheMerge`, `TestAHeadThatMovesAfterTheLastReadIsRefusedAndTheNewHeadIsReadBeforeAnotherCall` |
| The close-out and the carry-over | `TestTheCloseOutCommitCarriesTheVerdict`, `TestACloseOutThatChangesMoreThanTheTwoFilesIsRefusedAsScope`, `TestACloseOutOnAParentThatDoesNotCarryTheVerdictIsRefusedAsVerdict`, `TestCarryVerifyChecksTheWholeEdgeAndPostsNothingForAHeadThatIsNoCloseOut`, `TestTheRenderedCommitCarriesAndTheCloseOutAfterItIsAccepted`, `TestAVerifiedHeadThatIsAMergeCommitCarriesTheVerdict` |
| The merge and the event | `TestTheMergeIsAtTheHeadItReadAndWritesTheMergedEvent`, `TestAMergedPullRequestIsRecordedOnceAndTheAttemptIsNeverGuessed`, `TestALostLeaseOpensMarksAndMergesNothing` |
| The tree rule | `TestTheHeadIsExplainedOnlyByAChainOfMerges`, `TestTheTreeRowAcceptsAMergeNotYetRecordedWithItsLine`, `TestTheMergesAreReadFromTheEventsAndFromTheForgeForATaskWithNoEvent`, `TestAMergeNotYetRecordedPassesTheTreeCheck`, `TestAPullRequestFromDraftToMergeAndTheTreeAfterIt` (the demo) |

## Alternatives
| Decision | Selected | Rejected, and why |
| --- | --- | --- |
| Readiness of a context | the newest status of any creator must be the App's success | the App's newest only (a newer foreign failure would not block) |
| The carry-over | the verifier named, VerifyAt, the pre-image pinned to the verified head | a parent check alone (an ordinary child would carry) |
| The tree rule | a walk of first parents through known merges; the setup head only as the bottom | any descendant (admits a push); the setup head alone after a merge (a rollback) |
| The title | the form of the `pr-title` check, no task ID | `<type>: <task ID> <description>` (refused by the check) |

## Verification and verdict
Test-first (stubs alone in `68c2e94`, the integration level of `internal/run` run on them: exit 0 in 22 min 41 s; red runs in
the test-runs file) with 30 compiled mutants, all killed. One exception to the order: the merge-head control (`TestAVerifiedHeadThatIsAMergeCommitCarriesTheVerdict`) was added after the implementation; its stub replay proves sensitivity, not order. The two main red blocks (the interrupted register write, the records link) are kept as a labelled replay of 2026-10-18, not the original output. Round 1's head-moved check (finding 2) was written test-first, red in bc98321.

Rounds. Plan review (comment 6249142236, GPT-6.1 Sol): approve-with-conditions. Round 1 (comment 6252777914, GPT-6.1 Sol,
xhigh, Codex): **material**, recorded with the author's note: finding 1 was the budget, finding 2 a pull request head that was
not the bound head (fixed test first, red in `bc98321`, fix in `2e7b323`), finding 3 a stale architecture title, finding 4 a
missing dependency of 66h and 66k on 66q, finding 5 missing failing assertions (kept as labelled reconstructed evidence).
Round 2 (comment 6253345723, same reviewer, head `10b5bfb`): **not mergeable, findings recorded**, for the budget alone;
findings 2 to 5 were settled. Its one note, the wording of the task ID (the task ID is in the commit subjects, the issue
reference and the records link in the body), is fixed in `bec37c0` (`docs/architecture.md`, `docs/spec/loop.md`). The cap of 1
was reached with that round; the budget goes to the Operator's approval of the exact landing figure.
Harness exception: ADR-0026 names `claude`, `devin` and `opencode`; Codex was used on the Operator's message of 2026-10-13.

Merges of `origin/main`. First `06f274e` (#420, row 66s), merge `b36f9ab`: four index-only conflicts (the plan rows 66q
and 66s, the PRD change log, `packages.md` prose of `internal/run`, the traceability rows), resolved by keeping both sides
(O-219, no round); the plan rows 67, 68a and 68b are `main`'s, which add 66s to their After cells. Then `86f2c51` (#429, row
66w), merge `1518e3c`: one file, `PRD-0001-layup.md`, the Task and Test cells of NFR-001 and the change log, both sides
kept (O-219); each added line of this branch was checked present, and the change log has no blank line.

Gates after the merges, on `1518e3c`: `go build ./...`, `go vet ./...`, `go vet -tags=integration ./...`,
`go vet -tags=e2e ./cmd/layup/`, `gofmt -l .` (no file), `go test -count=1 ./...`, the five text checks,
`nested-checkout-check` and `git diff --check origin/main HEAD`: exit 0 each. The whole integration level
(`TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 90m ./...`): exit 0, every
package `ok` (`internal/run` 1987 s), 2026-10-18T19:42:19Z to 20:16:36Z. Two earlier runs on `b36f9ab` are recorded as
they were. The first, with the 45 m timeout, **failed**: `TestACoverageFileThatNamesNoFactIsRefused` failed with a cap
message and `internal/run` hit the timeout (2,700 s) while another session's integration run loaded the host (load average about 10);
the test passed alone (12.7 s), and the second run of `internal/run` with a 90 m timeout passed (2,113 s,
2026-10-18T19:03:24Z to 19:38:40Z). The exact-head CI is the Operator's merge gate.

Landing figure before this task file and the completed line: 3,168 lines (3,129 added, 39 removed) over 31 files, against a
maximum of 3,000/31 set by the plan review; the merge waits for the Operator's one approval of the exact figure (round 2's finding).

## Resource record
ADR-0007, recorded, not budgeted.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 | not recorded | not separately recorded | not separately recorded |
| Plan review | reasoning | GPT-6.1 Sol | xhigh | 174,204 | 2026-10-18T11:01:24Z to 11:08:23Z |
| Implementation and fix of round 1 | execution | Claude Sonnet 5.5 | not recorded | not recorded | commits `c0de25a` 15:58:09Z to `2e7b323` 17:19:58Z |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | not recorded | not recorded |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 225,817 | 2026-10-18T18:00:38Z to 18:09:26Z |
| Note fix, two merges, gates and close-out | execution | Claude Sonnet 5.5 | not recorded | not recorded | from 18:10Z; the integration runs 18:15:47Z to 19:02:27Z (failed), 19:03:24Z to 19:38:40Z (`internal/run` alone), 19:42:19Z to 20:16:36Z (green) |
| Orchestration | reasoning | Claude Opus 5.5 | not recorded | not separately recorded | not separately recorded |
