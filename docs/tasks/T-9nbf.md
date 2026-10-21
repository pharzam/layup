# T-9nbf — the step loop and the restart

Issue [#367](https://github.com/pharzam/layup/issues/367), row 66k of M2e (the composition of the seven size children 66v to 66ab). Serves REQ-005 and NFR-001 (F-0003#45).
Plan revision 1 6304795183 (Opus 5.5, it replaces the plan 6251067521); its plan review 6304893937 (GPT-6.1 Sol, six conditions, D1 to D11, Q1 to Q3); the Operator's decisions on Q1 and Q2 6307601453; the original review 6252684822 and decisions 6252712313 (goal count 2); the registration of the children 6252887262; the notes 6289097222 (the tier consumer) and 6304631026 (the reopened task).
Evidence: [test runs](../../runs/T-9nbf/test-runs.md), [mutants](../../runs/T-9nbf/mutants.md) and [conditions](../../runs/T-9nbf/conditions.md).

## Delivered and boundaries

- [x] **The step `loop`** (`internal/run/loop.go`): the merge order (specification, the tasks of `spec/reopenings.tsv`, the architecture task of Shape, then the bet's), one task acted on at a time with its predecessors merged, the register read again after a replacement and a reopening, `openReopenings` whenever the specification's `merged` event exists, the row (`done before`, `the loop merged …`, a wait as the last row, a stop as a `fail`, by the fields of `Advance`), the guard against an act that makes no progress, and a kind M2e does not run failing closed (rendered names 66t).
- [x] **The act on a task**, from the records and the forge: a merged or closed pull request is `FinishTask`'s before any repair or base movement; the pushes of a verified task are resolved, then `BringUpToDate`, then `FinishTask`; any other task goes to the driver, given the first copy of the answer and its options (`bindingOptions` for the bet's architecture tasks, the bound option the only survivor; `optionResults` for the architecture task of Shape and its redo).
- [x] **The restart.** `shape` and `bet` are `done before` once `bets.tsv` holds a bet (D1), so a restart after the first close-out, or after the architecture task of Shape was redone, reaches `loop` (F1, F2); the three steps run in `afterScaffold`, each after the check of the loop permissions (`forge: …`, the row of the step) and the parameter boundary, which is also called between two acts; an absent `parameters.tsv` is skipped only before the first task (F3). `run-steps` gains `loop`.
- [x] **A reopened task is not the original** (D9, note 6304631026): only the original specification task is asked for the rendered commit and for the specification that waits for its bet.
- [x] **Q2.** A specification attempt after the bet with no rendered commit is the stop `waiting for the Operator: <task>: an attempt after the bet has no rendered commit` (F4), before any driver. **Q1.** I12 is the consumer test of the tier from `Position`.
- [x] **Round 1.** A close-out push that the forge took and whose answer was lost is bound before the base moves; G1 and G5 get tests that kill them.
- [x] **Docs.** `loop.md` (The step loop, the known limit of the tier consumer, Not in M2e), `run.md`, `forge.md`, `packages.md`, row 66k and its After cell, PRD-0001 §12 cells, `traceability.md`, one entry of `guardrails.md` §2.
  Not delivered: rows 67 (the e2e of the specification task), 68a and 68b; M2g's items (the open-material variant of the tier test, whose first `by-class` step with a frozen list is `implement`; a later bet; row 66t, with a redo that holds the rendered commit and the follow-on of Q2).

## Verification and verdict

Tests and code were committed together (the hook runs `go test`); the red transcripts, with the exact commands, are in `runs/T-9nbf/test-runs.md` before the code commit: 13 of 13 unit and 17 of 17 integration tests red at the stubs, F1 to F3 reproduced and F4's premise, and the red of finding 1 of round 1 before its fix. The chronology, stated honestly: F4 is evidenced more narrowly (with no pull request for attempt 2, `FinishTask` fails at the read of the pull request, before the rendered-head error; the test asserts that error, and the stop of Q2 is the loop's); G1 and G5 were first recorded as equivalent and not reached, and were late coverage repairs after round 1, killed by a test through `afterScaffold` and by a test with a failing pinned `prd-lint` (their redness is the mutant, as those tests pass on the code of the freeze); G11 (the reread of the view after the pushes are resolved) is not killed and is kept as the specified order. The compiled mutants (hashes, controls) are in `runs/T-9nbf/mutants.md`: 29 of 30 killed, each on a copy of the tree. The full integration level, e2e, the six lints, the freeze check (K1 to K6) and gitleaks passed on the final code head.

Round 1, comment 6312207004 (GPT-6.1 Sol, xhigh, Codex): **material** (finding 1, a successful unresolved close-out push when the base moved; finding 2, G1; finding 3, G5), fixed. Round 2 (the same reviewer, record posted on #367): nothing material in scope.
**Verdict: nothing material in scope** (round 2). The Operator's goal count is 2. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Original plan | reasoning | Claude Opus 5.5 | not recorded | not kept | not kept (6251067521) |
| Original plan review | reasoning | GPT-6.1 Sol | xhigh | not kept | not kept (6252684822) |
| Plan revision 1 | reasoning | Claude Opus 5.5 (Plan agent) | not recorded | 233,565 | 2026-10-21T00:08Z to 00:18Z |
| Plan review of revision 1 | reasoning | GPT-6.1 Sol | xhigh | 207,812 | 2026-10-21T00:19:08Z to 00:25:45Z |
| Implementation and the fix of round 1 | execution | Claude Sonnet 5.5 (subagent) | not recorded | 586,672 (cumulative, reported at the round-1 fix) | 2026-10-21T07:50Z to 09:16Z |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 291,653 | 2026-10-21T08:45:20Z to 09:02:43Z |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 273,414 | 2026-10-21T09:19:39Z to 09:30:39Z |
| Close-out | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept at writing | 2026-10-21T09:32:05Z to 2026-10-21T09:32:32Z |
| **Total** | | | | 1,593,116 reported (the five numeric entries: 233,565 + 207,812 + 586,672 + 291,653 + 273,414 = 1,593,116); the complete total is not reported: the original plan, its review and the close-out have no figure | not reported (the rows have no common start and end) |
