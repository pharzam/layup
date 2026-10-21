# T-xyws — the forge calls that close a pull request and an issue

Issue [#423](https://github.com/pharzam/layup/issues/423), row 66w, a size child of row 66k
([#367](https://github.com/pharzam/layup/issues/367)). Serves REQ-005, NFR-001 (F-0003#45). Plan 6252962455
(Opus 5.5); plan review 6253052723 (GPT-6.1 Sol, six conditions, all met); the Operator's decisions, comment
6253151243: goal count 2, and `Closes #373` by whichever of 66q (#413) and 66w lands second. Maximum 800 lines
over 18 files, cap 1. Evidence: [sources](../../runs/T-xyws/api-sources.md), [test runs](../../runs/T-xyws/test-runs.md),
[mutants](../../runs/T-xyws/mutants.md).

## Delivered and boundaries

`ClosePullRequest(ctx, n) (PullRequest, error)` is `PATCH /pulls/{n}` with `{"state":"closed"}` alone;
`CloseIssue(ctx, n) error` is `PATCH /issues/{n}` with `{"state":"closed","state_reason":"not_planned"}` alone.
Both refuse `n < 1` before any request (the adapter makes none, not even a token) and read the answer: a state
that is not `closed` is a forge error of status 200, and so is an issue answer that holds `pull_request`
(after the effect, a known limit). `forge.Fake` and the test `standIn` give both; `Fake.IssueState(n)` shows
what a close changed. The R10 text is in `forge.md` ("The closes", each choice marked decided here), with the rows
in `packages.md` and a pointer in `loop.md`. Not delivered: the redo that calls them (66x), `EditBody` (66q).

**Decided here (condition 2):** a merged pull request is already closed, so a close of it is a success that gives it
with `Merged` true, in the adapter and in the `Fake` alike; the caller reads `Merged`. GitHub's page is silent, so the
`Fake` makes no `422` that the forge was never documented to give.

## Tests mapped to the criteria

| Criterion | Test |
| --- | --- |
| A close, observed in the stand-in; a repeat; a merged pull request | `TestFakeClosesAPullRequestAndReadsItClosedByNumberAndByBranch`, `TestFakeCloseOfAMergedPullRequestChangesNothingAndGivesItMerged` |
| An issue close, a second close, `IssueState` of open and unknown | `TestFakeClosesAnIssueAsNotPlannedAndKeepsItAfterASecondClose` |
| The refusal before the failure; 404, 502 | `TestFakeRefusesACloseBelow1BeforeItsFailureAndGivesTheForgeErrors` |
| The requests and the answers | `TestClosePullRequestPatchesTheStateAloneAndGivesThePullRequest`, `TestClosePullRequestGivesAMergedPullRequestAsTheForgeAnswersIt`, `TestCloseIssuePatchesTheStateAndTheReasonAlone` |
| Refusal with no request at all; errors; shared 401 and rate limit | `TestTheClosesRefuseANumberBelow1BeforeAnyRequest`, `TestTheClosesGiveAForgeErrorForAnyOtherStatusOrAnAnswerThatIsNotClosed`, `TestTheClosesShareTheHandlingOfA401AndOfARateLimit` |

## Alternatives

A reason parameter on `CloseIssue` (rejected: no second caller). A `GET /issues/{n}` before the close (rejected: one
more request, and the race between the check and the `PATCH` stays). A `422` for a merged pull request in the `Fake`
(rejected: undocumented, and it would diverge from the adapter). An effect-then-error switch in `Fake` (left to 66x).

## Verification and verdict

Own test-first tests (stubs committed alone at `6b06b11`, red recorded verbatim), 16 compiled mutants, all killed
(none survived; statuses and hashes per mutant in `mutants.md`). At head `eecbfce` (2026-10-18, 18:16Z to 18:23Z) `go build`, the three vets, `gofmt -l`, `go test ./...`, the integration level of `./internal/forge/...` and `./cmd/layup/`, the e2e level of `./cmd/layup/`, the five text lints, `nested-checkout-check` and `git diff --check` all exit 0 (direct statuses). The full integration level of `internal/run` was not run: no file of it changed. Final exact-head CI remains a landing condition.

The merge of `main` at `06f274e` (`1c90e0e`) resolved index lines only (O-219): the plan rows and the REQ-005
cells of PRD-0001 §12. It also removed a blank line that split the PRD change log before the row of `T-c6sx`.
That resolution dropped this task's NFR-001 cells, which round 1 found; they were restored at `75fd8b8`.

Round 1 (GPT-6.1 Sol, xhigh, on Codex), comment 6253669149: **material**, three findings, all in the documents
and the evidence: the NFR-001 Task and Test cells lost in the merge, the source record still naming the rejected
`422` of a merged pull request, and two misread lines of the red run (one green at the stub, one unreached).
Fixed at `75fd8b8`; the raw log is unchanged. Round 2, comment 6253810961: **nothing material in scope**, each
finding settled against its falsifier. **Verdict: mergeable**, at 609 lines over 15 files before this close-out,
inside the maximum of 800 over 18. Harness exception: ADR-0026 names `claude`, `devin` and `opencode`; Codex
was used on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens are as the harness reported them.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 | not recorded | 168,637 | about 5 min, ending 17:32Z |
| Plan review | reasoning | GPT-6.1 Sol | xhigh | 177,339 | 2026-10-18T17:32:59Z to 17:41:22Z |
| Implementation | execution | Claude Sonnet 5.5 | not recorded | 169,435 | 32 min, ending 18:25Z |
| Merge of `main` | reasoning | Claude Opus 5.5 | not recorded | not separately recorded | about 10 min |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 191,811 | 18:31:46Z to 18:41:11Z |
| Fix of round 1 | reasoning | Claude Opus 5.5 | not recorded | not separately recorded | about 6 min |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 179,922 | 18:47:46Z to 18:55:04Z |
| Close-out | reasoning | Claude Opus 5.5 | not recorded | not separately recorded | about 5 min |
| Orchestration | reasoning | Claude Opus 5.5 | not recorded | not separately recorded | not separately recorded |
