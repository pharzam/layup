# T-h65e — the brief of the first bet, its session, inputs and checks

Issue [#418](https://github.com/pharzam/layup/issues/418), row 66u (size child of row 66j, [#366](https://github.com/pharzam/layup/issues/366)). Serves REQ-012 (F-0003#51, F-0003#54).
Plan 6250830299 (Opus 5.5); plan review 6250925715 (GPT-6.1 Sol, seven conditions, all met); the Operator's goal count 6250932388 (1, four failure classes beneath it); cap-2 record 6290111962; maximum 2,600 lines over 26 files, cap 1, then 2.
Evidence: [test runs](../../runs/T-h65e/test-runs.md) and [mutants](../../runs/T-h65e/mutants.md).

## Delivered and boundaries

- [x] **The inputs and the session.** `betBrief` (`internal/run/betbrief.go`) checks the milestone and `first` kind, returns an already-kept brief before fresh-build checks, reads the verified specification and architecture heads (`verifiedHead`), compares registered requirement IDs with `finalIDs(draft)`, reads verified PRD statements, criteria and metrics (`spec.ContentOfPRDs`, refusing absent or blank statements), extracts `stillOpen` (`internal/run/shapequestions.go`) and `answerFactClass`, checks `brief.lines` before any session, and runs one `bet-brief` Product Owner phase session.
- [x] **The assembled brief and its checks.** Code assembles `bets/1/brief.md` around the session's `problem.md`, `solution.md` and `rabbit-holes.md`, presenting the verified heads, survivors with unmarked two-survivor tables, dropped options, facts with class and text, and open questions; refuses banned prefixes in the three files, heading/appetite/architecture/comment-size findings, option words outside code's part, links outside the verified heads' PRs and changed files, and candidates missing the current fixed parts.
- [x] **The appetite and the promotion.** Extracts `ledger.KnownMoney` (`internal/ledger/budget.go`); judges `spec.CheckAppetite` after the brief session's own telemetry row is committed; promotes the candidate to `bets/1/brief.md` with `bets/1/basis.tsv` (`bet-basis`, `internal/run/betbasis.go`, read back by `briefBasis`) in one fenced records commit.
- [x] **The retry and the recovery.** Retries a refused candidate with its finding in the next prompt; stops as `no-result` on a missing output file and starts afresh on the next run; recovers an unpromoted candidate across a stop after checking its fixed parts and appetite.
- [x] **Docs.** `loop.md`, `records.md`, `session.md`, `packages.md`, row 66u (and 66j's After cell) in `docs/plan/README.md`, PRD-0001 §12 cells, `traceability.md`, one entry of `guardrails.md` §2.
  Not delivered: posting the brief, taking the approver's answer, `bets.tsv`, the inventory and rendered commit, and `Restart` wiring (row 66j, #366); the later-bet caller (moved to row 66t, #417, in M2g; F4 of round 1).

## Verification and verdict

Inert stubs at `4dd75e4`, 33 + 6 + 2 compiled mutants killed with hashes and controls (`runs/T-h65e/mutants.md`). Chronology exceptions (the gate was not followed without exception): `TestAShapeTaskWithNoVerifiedHeadIsAnErrorBeforeAnySession` was added after its behavior and then run red against the stub (`runs/T-h65e/test-runs.md:114`); one round-1 red commit (`8cc825e`) was made with `--no-verify` (`runs/T-h65e/test-runs.md:203`); and F5 is a regression test of an existing guard whose red control is the compiled removal mutant M40 (`runs/T-h65e/test-runs.md:218`). Full integration level green on the merge head `0246d9c` (10:10:37Z to 10:28:41Z, exit 0), after round 1 (10:53:20Z to 11:11:24Z, exit 0), after round 2 (12:04:21Z to 12:11:58Z, exit 0) and on the merge head `2866a3f` (12:34:23Z to 12:42:27Z, exit 0); e2e and lints green; gitleaks exit 0 before each push.

Round 1 (GPT-6.1 Sol, xhigh, Codex), comment 6289116981: **material** (F1–F4, all settled in round 2). Round 2, comment 6289787413: **material** (F5–F6 at cap 1; the Operator chose one more cycle, cap 2 in 6290111962, and the fix was made test-first for F6 with compiled mutant M40 for F5). Round 3, comment 6290759372: settled F5–F6, **nothing material in scope**.
**Verdict: nothing material in scope.** The Operator's goal count is 1. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 | not recorded | not kept | not kept (posted with comment 6250830299) |
| Plan review | reasoning | GPT-6.1 Sol | xhigh | not kept | not kept (posted with comment 6250925715) |
| Implementation and the host-move gates | execution | Claude Sonnet 5.5 | not recorded | not kept | 2026-10-20T04:57:02Z to 10:32:06Z, commits 11aa82b to 645123a, with the host move inside |
| Fix of round 1 | execution | Claude Sonnet 5.5 | not recorded | not kept | commits 8cc825e to 5d5f564, 10:47:55Z to 11:14:57Z |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 410,225 | 2026-10-20T10:35:21Z to 10:44:36Z |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 229,952 | 2026-10-20T11:16:46Z to 11:26:00Z |
| Fix of round 2 | execution | Gemini 4 Argon High (`gemini-4-argon-high`) on Devin CLI | high | not kept | 2026-10-20T11:46:05Z to 12:16:22Z |
| Review round 3 | reasoning | GPT-6.1 Sol | xhigh | 207,602 | 2026-10-20T12:16:55Z to 12:23:43Z |
| Close-out | execution | Gemini 4 Argon High (`gemini-4-argon-high`) on Devin CLI | high | not kept at writing | 2026-10-20T12:30:04Z to 2026-10-20T12:48:00Z |
