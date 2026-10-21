# T-4naq — the specification task that a changed class reopens

Issue [#428](https://github.com/pharzam/layup/issues/428), row 66ab of M2e (K7, a size child of row 66k, [#367](https://github.com/pharzam/layup/issues/367)). Serves REQ-012 and REQ-005 (F-0003#51, F-0003#45).
Plan 6262367609 (Opus 5.5); plan review 6262491026 (GPT-6.1 Sol, seven conditions, Q1 to Q10); the Operator's decisions 6262589337 (goal count 3, the exemption of the full check is this row's, the restart criterion narrowed); the allocation Q15 6250705338; the registration of the children 6252887262.
Evidence: [test runs](../../runs/T-4naq/test-runs.md), [mutants](../../runs/T-4naq/mutants.md) and [conditions](../../runs/T-4naq/conditions.md).

## Delivered and boundaries

- [x] **`openReopenings(ctx, spec)`** (`internal/run/reopen.go`) opens, after the event `merged` of the specification task, one specification task for each fact that the first bet re-classed into `need`, `constraint` or `measure` (the version of the inventory against the class at Intake, a need changed to a constraint included), through the opening core `openTask` with the reserved ID, the specification as predecessor, no reuse by kind and no opening pin. Nothing is opened, written or asked of the forge before the event; a second version of the inventory fails closed.
- [x] **The record `spec/reopenings.tsv`** (the block `reopenings`, `ReopeningsSchema`, `ReadReopenings`): the fact, class, token (`F-NNNN#n` of the numbered facts record at the merge commit, or `Q-NNN`), the requirement ID of one REQ/NFR number space, the reserved task and the copy, all rows of a version in one fenced records commit before any issue. Each call binds the rows to the selection (version, facts, classes, copy, token), the empty selection included, and a registered task must agree with its reservation; draws avoid unopened reservations and their own batch.
- [x] **Recovery.** A stop after the allocation, after the issue or after the registration is finished by the next call under the same IDs, with one registered task for each row; the second issue of a stop between the issue and the row is the known limit the Operator accepted (6262589337).
- [x] **The pending prose and the prompt.** The requirement is one line under the heading "Waits for a bet", never a table row; the block `instruction-specification-reopened` (need, constraint and measure, three golden files) and the inputs (the whole PRD, the row, the pinned numbered record) reach every step of a reopened task.
- [x] **The exemption of the full check** (`FullInput.Reopening`, `FullInput.Pending`): a re-classed fact is exempt at every head but its own task's until that task has merged; the head of a reopened task must hold its pending line by the stored identity (fact and allocated ID), and the line of a merged reopening stays required at every later head, though a requirement of the table covers the fact.
- [x] **Records and docs.** `loop.md` (the new section, the full check, the instruction, the acceptance row), `records.md`, `architecture.md` section 7, `packages.md`, row 66ab (and 66ab in the After cell of 66k), PRD-0001 §12 cells (the earlier task IDs kept), `traceability.md`, one entry of `guardrails.md` §2.
  Not delivered: the drive of a reopened task by `FinishTask` and by the driver of a Shape task, which take every task of the kind `specification` for the first one (`internal/run/finish.go:87`, `internal/run/shapeloop.go:743`), is row 66k's, which tells the two apart by `reopenedTasks` (note 6304631026); the redo of a reopened task (it has no opening pin, so `redoIssue` refuses it before any write); a later bet's reopening (`M2g`).

## Verification and verdict

Tests and code were committed together (the hook runs `go test`); the red transcripts are in `runs/T-4naq/test-runs.md` before each code commit: those of the first version, of F1 to F3 and of F6 and F7. The compiled mutants (hashes, controls) are in `runs/T-4naq/mutants.md`: 66 of 66 killed on the final code head, each on a copy of the tree. The chronology, stated honestly: four checks (the numbered-record test, the reservation-conflict test, the register case of the not-reused test and the inputs assertion of the harness test) were written after the code, shown by their mutants and not by an order of writing, the exception the Operator accepted (finding F5, 6307600855); the harness test and the two-task test stop at their first step at the stub, so their own assertions had no red of their own; M58 is killed by a panic of the reader, which is also the red of F2; F8 asked for tests only, so its red is the reviewer's mutant. The full integration level was green on the final code head (exit 0, 21 packages), and so were e2e, the six lints, the freeze check and gitleaks.

Round 1 (GPT-6.1 Sol, xhigh, Codex), comment 6304746108: **material** (F1 to F5: the pending output of a constraint that the table covers, an empty class that panicked, the version of a reservation, the Task cells of the PRD, the order of four late checks), fixed except F5. Round 2, comment 6305218412: F6 to F8 (a merged reopening's line, an empty selection that skipped the binding, the constraint prompt) at cap 1; the Operator chose end 1, cap 2, and accepted F5's R8 exception (6307600855); fixed. Round 3, comment 6308023513: F9 (the text of `Pending`) at cap 2; the Operator chose end 1, cap 3 (6308037476); fixed by a text-only commit. Round 4, comment 6308157140: nothing material in scope.
**Verdict: nothing material in scope** (round 4). The Operator's goal count is 3. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 | not recorded | not kept | not kept (6262367609) |
| Plan review | reasoning | GPT-6.1 Sol | xhigh | not kept | not kept (6262491026) |
| Implementation and the fixes of rounds 1–3 | execution | Claude Sonnet 5.5 (subagent) | not recorded | 719,166 (cumulative, reported at the round-3 fix) | 2026-10-20T23:15Z to 2026-10-21T04:29Z |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 216,120 | 2026-10-21T00:07:04Z to 00:14:32Z |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 264,326 | 2026-10-21T00:40:28Z to 00:49:33Z |
| Review round 3 | reasoning | GPT-6.1 Sol | xhigh | 205,856 | 2026-10-21T04:16:48Z to 04:25:12Z |
| Review round 4 | reasoning | GPT-6.1 Sol | xhigh | 183,775 | 2026-10-21T04:30:33Z to 04:36:37Z |
| Close-out | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept at writing | 2026-10-21T04:37:18Z to 2026-10-21T04:37:59Z |
| **Total** | | | | 1,589,243 reported (the five numeric entries: 719,166 + 216,120 + 264,326 + 205,856 + 183,775); the complete total is not reported: the plan, the plan review and the close-out have no figure | not reported (the rows have no common start and end) |
