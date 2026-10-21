# T-08rj — the edit of a pull request's body

Issue [#413](https://github.com/pharzam/layup/issues/413), row 66q, a size child of row 66g
([#363](https://github.com/pharzam/layup/issues/363)). Serves REQ-007, NFR-004 (F-0003#47, F-0003#58). Plan 6250160997
(Opus 5.5); plan review 6250241345 (GPT-6.1 Sol, six conditions, all met); maximum 1,500 lines over 24 files, cap 1, raised
to 2 by the Operator at 6263439587. Evidence: [sources](../../runs/T-08rj/api-sources.md),
[test runs](../../runs/T-08rj/test-runs.md), [mutants](../../runs/T-08rj/mutants.md).

## Delivered and boundaries

**Step A.** `EditBody(ctx, n, body) (PullRequest, error)` is `PATCH /pulls/{n}` with `{"body": ...}` alone; it refuses `n < 1`
before any request. `forge.PullRequest` gains `Body` (a null body is the empty string), kept by `forge.Fake` (which also edits a
closed or merged pull request, as the forge does, and has `FailAfter` for an edit that takes effect and gives an error) and
given by the stand-in. The R10 text is in `forge.md`, merged with row 66w's closes (nine calls), with `packages.md`.

**Step B.** `(*Sessions).LinkReviewRecord(ctx, task, verifier) (edited bool, err error)` (`internal/run/reviewlink.go`) replaces
the body of the task's pull request with `PullBody` and the line `The review record: <web>/<owner>/<name>/issues/<n>#issuecomment-<id>`.
Every input is from the records: the verifier's start row (task, role, newest attempt), its **latest** handoff, which must be
`valid comment <ID>-<seen>`, the copy as a review record of the verifier's base, the pull request read by the attempt's branch
(the register's, open, at that head), and the opening records commit read from the body (`bodyRecords`, the inverse of `PullBody`)
and checked against the sessions that bound the head. It reads before it writes, reads the lease just before the edit, reads the
answer, and writes no record. The decisions are in `loop.md` ("The link of the review record in the body", marked decided here).
Not delivered: the step that calls it (66k), opening and merging (66g), closing (66w and 66x).

**Known limit (ruled by the Operator, 6264437496).** The edit accepts the records commit of the body it replaces when a session that
bound the reviewed head started at it. A person's edit that names another genuine records commit of a session that bound the same
head, or puts such a line before the opening one, is accepted and that commit is published; it is still a real commit of a session
of the head. A durable association (the snapshot written at the registration of the pull request, row 66g's code) is not built in M2e.

## Tests mapped to the criteria

| Criterion | Test |
| --- | --- |
| `EditBody` in the stand-in: the body kept, the refusal first, 404, 502, an edit that takes effect and errors, a merged pull request | `TestFakeKeepsTheBodyOfADraftAndEditsIt`, `TestFakeRefusesAnEditBelow1BeforeItsFailureAndGivesTheForgeErrors`, `TestFakeCanTakeAnEditAndStillGiveAnError` |
| The adapter: the request, the answer, other statuses, the refusal with no request, the body by number and by branch, 401 and rate limit | `TestEditBody*`, `TestPullRequestGivesTheBodyByNumberAndByBranchAndANullBodyIsEmpty` |
| The bytes of the body; one edit, then none | `TestTheBodyWithAReviewRecordEndsInTheLinkOfItsComment`, `TestTheBodyIsEditedOnceToLinkTheReviewRecord`, `TestAnEqualBodyIsNotEdited` |
| The opening snapshot through a close-out, a restart, a second verifier and two developers; an unresolved association | `TestTheRecordsLinkIsTheProducingSessionsNotTheCloseOutsNorTheVerifiers`, `TestTheEditKeepsTheSnapshotThatTheDraftWasOpenedWith` |
| Each refusal of the records, with a one-token control; the latest handoff after a valid one | `TestEachRecordThatFailsRefusesTheLinkBeforeAnyCall` |
| The attempt's pull request, the older one in the window, the head | `TestANewerAttemptWithTheRegisterOnTheOlderPullRequestEditsNothing`, `TestAPullRequestWhoseHeadMovedIsNotEdited` |
| The lease, a closed or merged pull request, a failed edit and one that took effect, a failed read, a failed store read, another body in the answer | `TestALostLeaseEditsNothing`, `TestAClosedOrMergedPullRequestIsNotEdited`, `TestAFailedEditIsAnErrorAndAnEditThatTookEffectIsNotRepeated`, `TestAFailedReadOfThePullRequestEditsNothing`, `TestAFailedStoreReadIsAnErrorWithNoForgeCall`, `TestAnAnswerWithAnotherBodyIsAnError` |
| The demo, through the production constructor, the adapter and real git | `TestTheReviewRecordIsLinkedFromTheBodyAfterItsPost`, `TestTheLinkSurvivesTheCloseOutAndARestart`, `TestAnEditRefusedByTheForgeIsAForgeErrorOfEditBody` |

## Alternatives

The first, the last, the last-before-the-verifier session and the developer role as the source of the snapshot (each refuted by
a review round: a verifier or a close-out binds the same head again). A new record or event for the opening snapshot (not built:
row 66g's registration is where it belongs). A `GET` of the pull request by number (rejected: the branch lookup also checks the attempt).

## Verification and verdict

Own test-first tests (stub committed alone at `f30ba1c`), 17 compiled mutants, all killed, each with a restored-test control. At the
merge head `94934f9` the full integration level ran green (2026-10-19T08:18:37Z to 08:42:40Z, exit 0); after the fix of round 2 the
unit level and the `internal/run` integration tests of the link and the pull request ran green (09:27:32Z to 09:29:19Z). The text of
the known limit and this close-out change documents only, so no new integration run was made. Final exact-head CI remains a landing
condition. gitleaks exit 0 before each push.

Round 1 (GPT-6.1 Sol, xhigh, Codex), comment 6262452176: **material**, two findings (the edit could replace the opening snapshot; the
Step B commands held placeholders), fixed at `e5ad2a5`. Round 2, comment 6263425521: verdict edited to **material** on the Operator's
end 1; cap raised to 2 at 6263439587; finding 1 (a second verifier moved the snapshot) fixed at `bb7d3f3`. Round 3, comment 6264271909:
**not mergeable, findings recorded** (finding 1: a hand edit naming another genuine commit; finding 2: the budget). The Operator's end 2,
comment 6264437496: finding 1 is the known limit above, and finding 2 is settled by one approval of the exact landing figure.
**Verdict: not mergeable until the Operator approves the exact landing figure.** Harness exception: ADR-0026 names `claude`,
`devin` and `opencode`; Codex was used on the Operator's message of 2026-10-13.

**Issue #373.** This task is the second lander of 66q and 66w, so its pull request says `Closes #413` and `Closes #373`; `forge.md`
holds the merged text of both calls.

## Resource record

ADR-0007, recorded, not budgeted. Tokens are as the harness reported them.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 (Plan agent) | not recorded | not kept | not kept |
| Plan review | reasoning | GPT-6.1 Sol | xhigh | 239,590 | 2026-10-18T12:54:42Z to 13:02:28Z |
| Step A implementation | execution | Claude Sonnet 5.5 | not recorded | 139,100 | 16 min |
| Merge of `main` and Step B | execution | Claude Sonnet 5.5 | not recorded | 290,774 | 81 min |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 239,228 | 2026-10-19T07:16:03Z to 07:26:15Z |
| Fix of round 1 and the merge of `d9852d2` | execution | Claude Sonnet 5.5 | not recorded | 327,510 | 47 min |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 197,989 | 08:44:21Z to 08:52:24Z |
| Fix of round 2 | execution | Claude Sonnet 5.5 | not recorded | 367,211 | 37 min |
| Review round 3 | reasoning | GPT-6.1 Sol | xhigh | 213,451 | 09:31:57Z to 09:40:19Z |
| Known limit text and close-out | execution | Claude Sonnet 5.5 | not recorded | not reported by the harness | ending 09:51Z |

The plan agent's figures were not read before the context that ran it ended. Codex is not in ADR-0026's routing list; used on the
Operator's message of 2026-10-13.
