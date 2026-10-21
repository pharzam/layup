# T-wqrk — a fix round's brief reads the fix diff and every home of the same rule

Issue [#462](https://github.com/pharzam/layup/issues/462), B4 of the M2e throughput investigation [#438](https://github.com/pharzam/layup/issues/438). It serves `F-0003#47`.
Plan [6314041425](https://github.com/pharzam/layup/issues/462#issuecomment-6314041425); plan review [6314145194](https://github.com/pharzam/layup/issues/462#issuecomment-6314145194) (approve-with-conditions); the Operator's decisions Q1 to Q3 [6315292816](https://github.com/pharzam/layup/issues/462#issuecomment-6315292816).
Budget maximum 200 lines over 6 files under ADR-0028; cycle cap 1, raised to 2 ([6316538920](https://github.com/pharzam/layup/issues/462#issuecomment-6316538920)).
Goal count 1. Round 1: [6315934925](https://github.com/pharzam/layup/issues/462#issuecomment-6315934925). Round 2: [6316255898](https://github.com/pharzam/layup/issues/462#issuecomment-6316255898). Round 3: [6316980681](https://github.com/pharzam/layup/issues/462#issuecomment-6316980681), not mergeable at cap 2. The Operator's end 3: [6316981250](https://github.com/pharzam/layup/issues/462#issuecomment-6316981250). The decision on the trial, option b: [6317693581](https://github.com/pharzam/layup/issues/462#issuecomment-6317693581).
Evidence: [the trial](../../runs/T-wqrk/trial.md), [test runs](../../runs/T-wqrk/test-runs.md), [harness snapshots](../../runs/T-wqrk/harness.md), [conditions](../../runs/T-wqrk/conditions.md), [the preparation tests](../../runs/T-wqrk/prep-test.sh).

## What is delivered

- **The fix-round preparer and brief**, installed in the shared harness by the coordinator at 2026-10-21T13:18:59Z (backups `*.bak-20261021T131859Z`): a round 2 or later reads the fix diff and searches for every home of the rule the fix touches, and the preparer refuses an earlier verdict and a malformed fix base.
- **The preparation tests** (`runs/T-wqrk/prep-test.sh`, 14 cases), red on the old preparer (12 of 14) and green on the new, with each of the eight guards removed in turn failing its own case.
- **The trial and its honest result** ([the trial](../../runs/T-wqrk/trial.md) sections 5 and 9): six runs of GPT-6.1 Sol on a planted two-home fix. Old brief 1 of 3 finds, new brief 0 of 3, on the frozen judge, so **red holds and green fails: no gain is shown**. The coordinator's hand reading found the planted disagreement in all six records; under term 2 of the end 3 a judge MISS stays a miss, so that reading is reported and not counted. The plant was too easy: the old brief finds it as well.

## Not delivered

- **The gate rule** (the Method row of "Who may review" and the `Briefed on` clause), dropped on the Operator's option b because no gain was shown. `docs/engineering-discipline.md` is unchanged by this task.
- A second, harder plant. None was chosen after the result (section 6 of the trial).

## Verification and verdict

Red first, then green, with `date -u` times: the preparation tests, and the judge (four controls, red on the third judge at 13:15:18Z, green on the fourth at 13:16:21Z). The three rounds are on the issue. Round 3's finding F1 (the judge scored two notes that diagnose nothing) was fixed under the end 3 **with no round on the fix**, the Operator's terms 1 to 3: the repaired judge, the hand confirmation of every FOUND, and no spend before both were pushed. The six runs followed.

**Verdict: not mergeable, findings recorded; the finding was fixed under the Operator's end 3, and the judge fix has had no round.** Goal count 1. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them; times from the clock (`date -u`) or the commit time in UTC.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 (Plan agent) | not recorded | not kept | 2026-10-21T10:40:11Z (claim) to 10:44:47Z (the comment) |
| Plan review | reasoning | GPT-6.1 Sol | not recorded | not kept | to 2026-10-21T10:50:52Z (the comment) |
| Build, phase 1 (driver, judge, preparer, tests) | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-21T12:00:52Z to 12:06:15Z (commits) |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 155,087 | 2026-10-21T12:18:07Z to 12:28:54Z |
| Fix of round 1 | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-21T12:34:50Z to 12:35:12Z (commits) |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 156,431 | 2026-10-21T12:37:09Z to 12:44:54Z |
| Cap-raise plan review | reasoning | GPT-6.1 Sol | not recorded | not kept | to 2026-10-21T12:55:05Z (the comment) |
| Fix of round 2 | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-21T12:56:17Z to 12:59:01Z (commits) |
| Review round 3 | reasoning | GPT-6.1 Sol | xhigh | 143,558 | 2026-10-21T13:03:43Z to 13:11:07Z |
| Judge repair (end 3) | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-21T13:15:55Z to 13:16:45Z (commits; the judge runs 13:15:18Z to 13:16:21Z) |
| Trial, six runs | reasoning | GPT-6.1 Sol | high | 607,091 (66,503; 128,623; 80,483; 111,817; 88,764; 130,901) | 2026-10-21T13:19:27Z to 13:40:23Z (20 min 56 s) |
| Results, drop of the gate text, close-out | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-21T13:50:18Z to 2026-10-21T13:50:24Z |
| **Total** | | | | 1,062,167 reported (the four numeric rows); the complete total is not reported | not reported (the rows have no common start and end) |
