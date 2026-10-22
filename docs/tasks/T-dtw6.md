# T-dtw6 — a task session that edits a rule path to keep a document current is told, and its refusal gets one new attempt

Issue [#491](https://github.com/pharzam/layup/issues/491), a child of row 68b ([#370](https://github.com/pharzam/layup/issues/370)), plan row 68b4. It serves `F-0003#45`, through REQ-005 (the task loop) and REQ-003 (rule protection).
Budget maximum 800 lines over 28 files, close-out inside, under ADR-0028; cycle cap 1. Goal count 1.
Evidence: [test runs](../../runs/T-dtw6/test-runs.md), [conditions](../../runs/T-dtw6/conditions.md).

## Plans and reviews

- Plan [6346188891](https://github.com/pharzam/layup/issues/491#issuecomment-6346188891).
- The Operator's decisions [6346348137](https://github.com/pharzam/layup/issues/491#issuecomment-6346348137): all as recommended.
- Plan review [6346542416](https://github.com/pharzam/layup/issues/491#issuecomment-6346542416) (approve-with-conditions, five conditions: `runs/T-dtw6/conditions.md`).

## What is delivered

- The instruction of each step that commits (`specification`, the reopened `specification`, `shape-architecture`, `implement`) names the rule paths of the records' register, with the exception of `docs/guardrails.md`, and says that a change to one refuses the whole session. The driver of `implement` is outside this task (a build task is not run in M2e); its rendered instruction is tested. Six golden files and four blocks of `docs/spec/loop.md` carry the sentence.
- A `done` session with the status `completed` and exactly one refusal whose reason is exactly `rule-path` is an eligible end of its stage with its own retry (once per stage and class: at most twelve retries and fifteen added sessions per task). `rule-paths`, `workflow`, `base` and `push-refused` stay stops. The refused head is never pushed or bound.
- `rules.PayloadPaths` reads the paths of the refused payload (`payloads/<sha>`, `git diff --binary --no-renames`) and fails closed; `rulePathFindings` keeps the paths the current register matches; the developer of the new attempt is told them as `rule-path <path>:` findings. A payload that is missing, unreadable or unparseable is the run's own error, raised before the new attempt is recorded.
- The rule and the contract of the parser in `docs/spec/loop.md`, a pointer in `docs/spec/session.md`, plan row 68b4, the REQ-005 and REQ-003 cells of the PRD, traceability, and one lesson in `guardrails.md` §2.

## Verification

Red first, then green, with `date -u` times in the test runs; see `runs/T-dtw6/test-runs.md`.
