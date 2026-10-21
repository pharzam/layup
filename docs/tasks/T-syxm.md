# T-syxm — the repairs of a task session at a restart

Issue [#422](https://github.com/pharzam/layup/issues/422), row 66v, a size child of row 66k
([#367](https://github.com/pharzam/layup/issues/367)). Serves REQ-005, NFR-001 (F-0003#45). Plan 6252986032
(Opus 5.5); plan review 6253071393 (GPT-6.1 Sol, five conditions, all met); the Operator's goal count 2 (comment
6253150947). Maximum 1,650 lines over 28 files, cap 1. Evidence: [test runs](../../runs/T-syxm/test-runs.md),
[mutants](../../runs/T-syxm/mutants.md).

## Delivered and boundaries

`(*Sessions).RepairTask(ctx, task)` repairs a task's sessions from `events.tsv`, `sessions.tsv`, `tasks.tsv` and the
forge, never `bets.tsv`, so it works during Shape; each driver of a task calls it before it reads the position (the
text is in `loop.md`; the call from 66r's `AdvanceShapeTask` and from 66k's act is their own change).
- **Class 1, the push.** `RepairPushes` reads the branch of each `push` with no `bound` and no `refused`
  (`SessionStore.BranchHead`, over the new `git.RemoteHead`, an authenticated `ls-remote --exit-code` that keeps
  only the exact line `refs/heads/BRANCH`). A head equal to the SHA, byte for byte, is `bound`; any other result
  (absent, another SHA, a descendant) is `refused push-lost`; a failed read is an error and writes nothing. A push
  of an older attempt gets the same event and no attempt. The replacement rule: an open attempt (the existing
  `openAttempt`, so a `closed` or `rebased` latest attempt counts as closed) that holds a `push-lost` refusal gets
  one `StartAttempt` at its base, which also covers a restart between the two commits.
- **Class 2, the comment.** The comment of a task session moves from the control issue to the task's issue
  (`taskIssue`; no register, no row or a refused register is an error, no fallback). It reads the issue's comments
  before the post (the App's own comment with the same bytes is not posted again) and the lease right before it;
  `RepostComments` posts the missing line of each session of the task from the events. Probes and phase sessions
  keep the control issue, and a session with no `result` event has no comment, before any read of the register.
- Not delivered: the driver's call (66r, 66k), the question and answer comments (66z), the spec and rules retries
  (66aa), `git.Merge` (66x).

## Tests mapped to the criteria

| Criterion | Test |
| --- | --- |
| A Shape push repaired before any bet | `TestARestartRepairsAShapePushBeforeAnyBet` |
| Bound on an exact head; push-lost on an absent or another head; a failed read is an error | `TestAPushWhoseBranchHoldsItsSHAIsBound`, `TestAPushWhoseBranchIsAbsentIsPushLostAndStartsOneAttempt`, `TestABranchHeadThatDiffersFromThePushIsPushLost`, `TestAMovedBranchHeadIsPushLost`, `TestAFailedBranchReadIsAnErrorAndWritesNothing` |
| Exactly one new attempt, also between the two commits; only for an open latest attempt | `TestAnOpenAttemptHoldingPushLostGetsExactlyOneAttempt`, `TestARestartBetweenPushLostAndTheAttemptStartsExactlyOne`, `TestSeveralLostPushesOfTheOpenAttemptGiveOneAttempt`, `TestAClosedOrRebasedLatestAttemptGetsNoReplacement`, `TestAPushLostInAnOlderAttemptStartsNoAttempt`, `TestTheNewAttemptStartsAtTheLostAttemptsBase`, `TestAPushThatIsAlreadyBoundOrRefusedIsNotReadAgain` |
| The read of a branch | `TestRemoteHeadRunsItsVerb`, `TestRemoteHeadRefusesABranchBeforeGitStarts`, `TestRemoteHeadReadsOnlyTheExactRefLine`, `TestRemoteHeadOfABareRepository`, `TestRemoteHeadSendsTheTokenAsAHeaderOnly` |
| The comment on the task's issue, read before the post, reposted from the records | `TestATaskSessionsCommentIsOnItsIssue`, `TestASessionCommentThatTheAppPostedIsNotPostedAgain`, `TestAMissingSessionCommentIsPostedAgainFromTheRecords`, `TestTheSameTextFromAPersonOrAnotherAppDoesNotSuppressTheSummary` |
| No register, no row or a refused register; a failed read of the comments or the register; a lost lease, on both paths | `TestNoRegisterRowOrABrokenRegisterPostsNowhere`, `TestTheGuardsOfThePostHoldOnBothPaths` |

## Alternatives

One atomic commit of `refused` and `attempt`: rejected, it would remove only that window and not find a refusal that
is already recorded. Binding on "found" alone: rejected (Q7), code also advances task branches. Session comments
copied into `copies.tsv`: rejected, they are built from the events and are not artifacts.

## Verification and verdict

Own test-first tests (the stubs committed alone, the red run recorded with its observed outcomes, one of the plan's
green negatives was red), 22 compiled mutants, all killed (M21 was rebuilt once because it did not build; per-mutant
statuses and hashes in `mutants.md`). The first full run of the integration level failed on the lease of the session
world (26 tests), which the plan had not moved with the register; fixed in ac05f6a, and the full level then passed at
that head (exit 0, 19:23:40Z to 19:57:06Z). Build, vet (plain, integration, e2e, uat), gofmt, unit, e2e and the text
checks pass. The full integration level at the merge head 9260f17 failed on two load flakes of this host (the `fastClock` of #291 and a push under load), which passed alone at 7f1ba6d; the runs are in the test record.

Rounds (GPT-6.1 Sol, xhigh, on Codex; cap 1):
- Round 1, comment 6254891498: **material**, three findings (`Repair.Attempt` set after a failed `StartAttempt`; the handoff tests selecting the artifact by exclusion or index; no test of a refused `sessions.tsv`). Fixed at 7f1ba6d, test first (mutants M23 to M25).
- Round 2, comment 6255109180, on 7f1ba6d: **not mergeable, findings recorded**. The three round 1 findings are settled. Finding 1 (material): no test of a refused `events.tsv` in `RepairPushes` and `RepostComments`. Finding 2 (material): the close-out did not fit the landing budget.
- The cap allows no third round. The Operator chose **end 3** (comment 6260569452, 2026-10-19T05:44Z): finding 1 fixed with no further round, then this close-out, and finding 2 settled by one Operator approval of the exact landing figure.
- The fix of finding 1, with no round: `TestARefusedEventsTableStopsTheRepairAndTheRepost`, written first, green on the code and red on its mutants M26 and M27 (a third variant, the first read of `RepairPushes` alone, is not a kill: its second read returns the same error). No production code changed. No reviewer has read this commit.
- The merge of `main` at c735788 (4a2119d, row 66g): index lines and the fields of one test type. Its full integration run failed one test of row 66g that opens a task session of an unregistered task, which this branch refuses; the test now registers that task (99572ca, test only, no round). The runs are in [`test-runs.md`](../../runs/T-syxm/test-runs.md).

**Verdict: not mergeable until the Operator approves the exact landing figure.** Finding 1 is fixed and unread by any round (end 3); finding 2, the budget, is open and is not the author's to close: the maximum is 1,650 lines over 28 files with the close-out inside, and the figure against the current `main` goes to the Operator on the issue. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex was used for the plan review and both rounds on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens are as the harness reported them.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 (Plan agent) | not recorded | 185,513 | not recorded |
| Plan review | reasoning | GPT-6.1 Sol (Codex CLI) | xhigh | 219,396 | 2026-10-18T17:35:15Z to 17:42:56Z |
| Implementation | execution | Claude Sonnet 5.5 | not recorded | 288,491 | not recorded |
| Review round 1 | reasoning | GPT-6.1 Sol (Codex CLI) | xhigh | not kept: its run directory was removed before they were read | not kept |
| Fix of round 1 | execution | Claude Sonnet 5.5 | not recorded | 96,158 | not recorded |
| Review round 2 | reasoning | GPT-6.1 Sol (Codex CLI) | xhigh | 183,806 | 2026-10-18T20:54:21Z to 21:02:09Z |
| Fix of end 3 and close-out | execution | Claude Sonnet 5.5 | not recorded | 87,972 | 38 min, 2026-10-19T05:45Z to 06:23Z, the full integration level included |
| Merge of `main` and its test fix | reasoning | Claude Opus 5.5 (the coordinator) | not recorded | not reported by the harness | 2026-10-19T06:24Z to 07:15Z, the full integration level included |
