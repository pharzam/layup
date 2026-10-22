# T-dtw6 — a task session that edits a rule path to keep a document current is told, and its refusal gets one new attempt

Issue [#491](https://github.com/pharzam/layup/issues/491), a child of row 68b ([#370](https://github.com/pharzam/layup/issues/370)), plan row 68b4. It serves `F-0003#45`, through REQ-005 (the task loop) and REQ-003 (rule protection).
Budget maximum 800 lines over 28 files, close-out inside, under ADR-0028; cycle cap 1. Goal count 1.
Evidence: [test runs](../../runs/T-dtw6/test-runs.md), [conditions](../../runs/T-dtw6/conditions.md).

## Plans and reviews

- Plan [6346188891](https://github.com/pharzam/layup/issues/491#issuecomment-6346188891).
- The Operator's decisions [6346348137](https://github.com/pharzam/layup/issues/491#issuecomment-6346348137): all as recommended.
- Plan review [6346542416](https://github.com/pharzam/layup/issues/491#issuecomment-6346542416) (approve-with-conditions, five conditions, all met: `runs/T-dtw6/conditions.md`).
- Round 1 [6347312700](https://github.com/pharzam/layup/issues/491#issuecomment-6347312700): nothing material in scope.

## What is delivered

- The instruction of each step that commits (`specification`, the reopened `specification`, `shape-architecture`, `implement`) names the rule paths of the records' register, with the exception of `docs/guardrails.md`, and says that a change to one refuses the whole session. The driver of `implement` is outside this task (a build task is not run in M2e); its rendered instruction is tested. Six golden files and four blocks of `docs/spec/loop.md` carry the sentence.
- A `done` session with the status `completed` and exactly one refusal whose reason is exactly `rule-path` is an eligible end of its stage with its own retry (once per stage and class: at most twelve retries and fifteen added sessions per task). `rule-paths`, `workflow`, `base` and `push-refused` stay stops. The refused head is never pushed or bound.
- `rules.PayloadPaths` reads the paths of the refused payload (`payloads/<sha>`, `git diff --binary --no-renames`) and fails closed; `rulePathFindings` keeps the paths the current register matches; the developer of the new attempt is told them as `rule-path <path>:` findings. A payload that is missing, unreadable or unparseable is the run's own error, raised before the new attempt is recorded.
- The rule and the contract of the parser in `docs/spec/loop.md`, a pointer in `docs/spec/session.md`, plan row 68b4, the REQ-005 and REQ-003 cells of the PRD, traceability, and one lesson in `guardrails.md` §2.

## Findings of round 1, and their closing

Round 1 found nothing material in scope and no raw finding. Its notes, recorded with no change:

- The demo test uses reconstructed events and a synthetic payload of the original's shape, so it does not replay the original payload `c844bf5e…`.
- The reviewer reproduced the unit reds against the base; the historical order of work and the full integration and e2e results stay author reports.
- The reviewer's selected integration run aborted in the `httptest` fixture because its sandbox denies sockets, so it gives no behavioral verdict.
- The close-out was pending at review time.

## Verification and verdict

Red first, then green, with `date -u` times in the test runs (`runs/T-dtw6/test-runs.md`; the hook blocks a red commit, so each slice's tests and code are one commit). Verified by running: build, vet (three tags), gofmt, unit, the full integration level (exit 0), the e2e level (exit 0), the six lints, `git diff --check`, and the freeze check (K1 to K6).

Round 1, comment 6347312700 (GPT-6.1 Sol on Codex): **nothing material in scope.**

**Verdict: nothing material in scope.** Goal count 1. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used for the plan review and for round 1 on the Operator's message of 2026-10-13.

## Acceptance

- [x] Each step that commits names the rule paths and the refusal, in four blocks and six goldens.
- [x] A `rule-path` refusal alone gets one new attempt per stage, bounded per stage and class.
- [x] The developer of the new attempt is told the paths, read from the payload; the demo's records give attempt 2 at `d03c865`.
- [x] A payload that cannot be read is the run's error before another session starts; the refused head is never pushed or bound.
- [x] The rule, the parser contract, the PRD, traceability, plan row 68b4 and the lesson are written.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them; build times are commit times.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 (a Plan agent) | not recorded | not kept | not kept |
| Plan review | reasoning | GPT-6.1 Sol (Codex) | xhigh | 216,447 | 2026-10-22T16:56:10Z to 17:04:57Z |
| Build | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-22T17:13:05Z to 17:31:51Z (first to last commit) |
| Review round 1 | reasoning | GPT-6.1 Sol (Codex) | xhigh | 175,447 | 2026-10-22T17:33:53Z to 17:43:13Z |
| Close-out | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | not kept; its commit follows c27db1f at 2026-10-22T17:31:51Z |
| **Total** | | | | 391,894 reported (the two numeric entries); the complete total is not reported | not reported |
