# T-3xvc — the e2e of the specification task

Issue [#368](https://github.com/pharzam/layup/issues/368), row 67 of M2e. Serves REQ-005, REQ-007, REQ-012 and NFR-004 (F-0003#45, F-0003#58). Goal count 1.
Plan 6312881896 (Opus 5.5, Plan agent); its plan review 6313077503 (GPT-6.1 Sol: seven conditions, D1 to D3, Q1 to Q3, budget 850 lines over 12 files, cap 1); round 1 6313804500; round 2 6314116308; the Operator's end 3 6314134298.
Evidence: [test runs](../../runs/T-3xvc/test-runs.md), [mutants](../../runs/T-3xvc/mutants.md) and [conditions](../../runs/T-3xvc/conditions.md).

## Delivered and boundaries

- [x] **`TestTheSpecificationTaskEndToEnd`** (`cmd/layup/loop_e2e_test.go`): the built binary, a fake forge with the M2e routes (the e2e copy of the integration fakes: comments of any issue, statuses, check runs of the Actions app, pull requests with a real merge commit), local bare repositories and scripted harnesses, from `layup run --new` through Shape (one survivor, one drop), the first bet and the loop to the merge of both Shape tasks. Seven runs of the binary, each counted; the repeat restart makes no forge write. Per task: one merge request at the head H after a ready mark, the four App statuses at H as they stood at the request, the merge commit M with the second parent H, the event `merged` at M with the pull request number, every handoff valid; the rendered commit changes only the two cells of the requirement's row, compared on the PRD bytes.
- [x] **The oracle probes** `TestTheMergeOracleRefusesStatusesPostedAfterTheMergeRequest` and `TestTheCountOfMergeRequestsIsExact` (the control C1 expects exactly zero merge requests, the scenario exactly two).
- [x] **Two defects on the path, fixed (R10).** `internal/run/steps.go` (`copy`): it wrote `copies.tsv` from the rows read at the start of the run, so the first read of the Operator's `set` (after the first task) dropped the copies of the Shape tasks' comments, and `bet` found no verified head; red: `bet fail … Shape has no verified head`. `internal/run/finish.go` (`FinishTask`): the four checks of a task brought current by a merge of the default branch judged it against the base of the attempt, so the numbered facts record that the merge brought in failed `layup/rules` for the architecture task; red: `loop … layup/rules: failure`. Both reds are in `runs/T-3xvc/test-runs.md`; each was reproduced by reverting its fix alone.
- [x] **Docs.** `loop.md` (the paragraph "The specification task, end to end", the demo item 4), `traceability.md`, the PRD-0001 §12 Test cells, one entry of `guardrails.md` §2.
  Not delivered: 68a (`T-ym0k`, the release review of M2e), 68b (`T-8s13`, the demo on a real GitHub target with the code-owner refusal, the real CI jobs and the Operator's inputs), and M2g's items (build tasks, the milestone plan, the rendered task 66t, a redo of the rendered specification, frozen-test amendments).

## Verification and verdict

Tests and code were committed in steps (the hook runs `go test`). The chronology, stated honestly: the first run of the scenario failed on fixtures, not on the demo (a constraint of the baseline's accepted decision record, the success row of the draft, colliding comment IDs of the fake); the first product red came after those; the product fixes followed it. Verified by running: build, vet with the three tags, gofmt and the unit level; the full integration level (exit 0, at `5ad1405f`; no product code changed since); the e2e level (exit 0, package `cmd/layup` 85.7 s, baseline 64.4 s); the six lints, `git diff --check`; the freeze check K1 to K6 (base `30e2c0d`); five compiled mutants (E1, E2, E3, C1a, C1b), all restored by hash. Not verified: the exact-head CI run (the Operator's end 3 leaves it to CI); the reviewers could not run the socket tests in their sandbox and relied on the recorded runs. Known limit: the fake forge is a second copy of the integration one, kept in step by hand (D1).

Round 1, comment 6313804500 (GPT-6.1 Sol, xhigh, Codex): **material** (the merge oracle read the final status map, not the statuses at the merge request; C1 lacked a no-merge assertion), fixed. Round 2, comment 6314116308 (the same reviewer): **not mergeable, findings recorded**, one finding at cap 1 (C1 expected two merge requests, so zero and one failed alike). The Operator chose end 3 (6314134298, 2026-10-21T10:50:12Z): the finding was fixed with no further round, by an explicit zero count for C1, a negative probe, and a rerun of C1 killed by the zero assertion.
**Verdict: not mergeable, findings recorded; the finding was fixed under the Operator's end 3.** Goal count 1. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 (Plan agent) | not recorded | 169,169 | 2026-10-21T09:33Z to 09:41Z |
| Plan review | reasoning | GPT-6.1 Sol | xhigh | 229,749 | 2026-10-21T09:41:13Z to 09:49:47Z |
| Implementation and the fix of round 1 | execution | Claude Sonnet 5.5 (subagent) | not recorded | 395,027 (cumulative, reported at the round-1 fix) | 2026-10-21T09:52Z to 10:39Z |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 172,071 | 2026-10-21T10:22:45Z to 10:30:58Z |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 174,950 | 2026-10-21T10:40:35Z to 10:48:53Z |
| Fix of round 2 and close-out | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept at writing | 2026-10-21T10:50:37Z to 2026-10-21T10:58:40Z |
| **Total** | | | | 1,140,966 reported (the five numeric entries: 169,169 + 229,749 + 395,027 + 172,071 + 174,950); the complete total is not reported: the fix of round 2 and the close-out have no figure | not reported (the rows have no common start and end) |
