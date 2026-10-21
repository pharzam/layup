# T-r7b9 — a question during the work, in the loop

Issue [#426](https://github.com/pharzam/layup/issues/426), row 66z (size child K5 of row 66k, [#367](https://github.com/pharzam/layup/issues/367)). Serves REQ-005 (F-0003#45).
Plan 6262327790 (Opus 5.5); plan review 6262463499 (GPT-6.1 Sol, seven conditions); the Operator's decisions 6262588622 (goal count 4, kept at the maximum) and 6289096719 (round 1 findings 2 and 4); cap-2 record 6289990772; maximum 2,500 lines over 27 files, cap 1, then 2.
Evidence: [test runs](../../runs/T-r7b9/test-runs.md) and [mutants](../../runs/T-r7b9/mutants.md).

## Delivered and boundaries

- [x] **The question stages in the loop.** When a session ends `needs_context`, `AdvanceShapeTask` takes `posQuestion` (`questionPosition`, `askStage`, `internal/run/questionloop.go`): records the asking handoff, posts and copies each question via `QuestionComment` and `postBody`, starts the read-only `answer-to-a-question` session on the owner role at the last asking session's base, posts and copies each answer for its column `n` question via `AnswerComment`, records the answer handoff and starts the next attempt; invalid handoffs, non-completed owner sessions, missing or excess answers and uncomposable rows stay repeatable stops.
- [x] **The human and mixed-owner waits.** A question routed to a human or questions with more than one owner role return `waiting for the Operator: <task>: question <id>` with no answer session.
- [x] **The answer prompt and citations.** The blocks `instruction-answer-to-a-question` and `instruction-answers` in `loop.md` and `instructions.go`, the golden file `testdata/prompts/answer-to-a-question.txt`, the asking results in `handoffs-and-answers` and `readOnlyStep("answer-to-a-question")`.
- [x] **Material findings in `Position`.** `openMaterial` (`internal/run/tier.go`) selects the last verifier by its `bound` event and attempt, then qualifies its `completed` result, valid comment handoff, kept file hash and reviewed head; an open, missing or unqualified record keeps `Position` uphill.
- [x] **Docs.** `loop.md`, `session.md`, `packages.md`, row 66z (and 66k's After cell) in `docs/plan/README.md`, PRD-0001 §12 cells, `traceability.md`, one entry of `guardrails.md` §2.
  Not delivered: the production tier consumer of `Position` → `TierOf` and its semantic test (released to row 66k, #367, by Operator decision 6289096719 on finding 2 of round 1), human/mixed-owner question resolution, and frozen-test amendment execution (M2g).

## Verification and verdict

Test first (inert stubs at `103a318`), 30 + 10 + 3 compiled mutants killed with hashes and controls (`runs/T-r7b9/mutants.md`). Full integration level green on the merge head (2026-10-20T10:14:14Z to 10:30:52Z, exit 0), after round 1 (11:03:15Z to 11:19:46Z, exit 0) and after round 2 (11:56:27Z to 12:03:22Z, exit 0); e2e and lints green; gitleaks exit 0 before each push.

Round 1 (GPT-6.1 Sol, xhigh, Codex), comment 6289084862: **material** (F1–F4; F1 and F3 fixed, F2 deferred to 66k and the F4 R8 exception accepted by 6289096719). Round 2, comment 6289957559: **material** (F5–F8 at cap 1; the Operator chose one more cycle, cap 2 in 6289990772, and the fix was made test-first). Round 3, comment 6290566991: settled F5–F8, **nothing material in scope**.
**Verdict: nothing material in scope.** The Operator's goal count is 4. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 | not recorded | not kept | not kept (posted with comment 6262327790) |
| Plan review | reasoning | GPT-6.1 Sol | xhigh | not kept | not kept (posted with comment 6262463499) |
| Implementation and the host-move gates | execution | Claude Sonnet 5.5 | not recorded | not kept | 2026-10-20T04:59:23Z to 11:31:11Z, commits f64315d to d218424, with the host move inside |
| Fix of round 1 | execution | Claude Sonnet 5.5 | not recorded | not kept | commits 494d36b to 9bb845e, 10:51:13Z to 11:20:07Z |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 233,520 | 2026-10-20T10:32:54Z to 10:42:28Z |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 203,534 | 2026-10-20T11:26:12Z to 11:33:41Z |
| Fix of round 2 | execution | Gemini 4 Argon High (`gemini-4-argon-high`) on Devin CLI | high | not kept | 2026-10-20T11:39:00Z to 12:06:05Z |
| Review round 3 | reasoning | GPT-6.1 Sol | xhigh | 206,007 | 2026-10-20T12:06:54Z to 12:12:27Z |
| Close-out | execution | Gemini 4 Argon High (`gemini-4-argon-high`) on Devin CLI | high | not kept at writing | 2026-10-20T12:13:29Z to 2026-10-20T12:15:29Z |
