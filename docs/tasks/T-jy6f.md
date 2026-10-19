# T-jy6f — the Shape tasks through the task loop

Issue [#415](https://github.com/pharzam/layup/issues/415), row 66r, a size child of row 66h ([#364](https://github.com/pharzam/layup/issues/364)).
Serves REQ-005, REQ-012 (F-0003#45, F-0003#51). Plan 6250184301 (Opus 5.5); plan review 6250262119 (GPT-6.1 Sol, seven conditions, all met);
amendment 1 (the base of a session, 6264868822); scope addition 6262589736; goal count 3 (6250443025); maximum 3,100 lines over 25 files, cap 1, then 2.
Evidence: [test runs and mutants](../../runs/T-jy6f/test-runs.md).

## Delivered and boundaries

- [x] **Opening.** `(*state).openShapeTask` opens the specification and the architecture task once: permissions, draft, drawn ID, dry run of the commit, lease, issue, row
  (`shapetask.go`); `finalID` is the one home of the requirement mapping; the issue text is two blocks of `loop.md`.
- [x] **Driver.** `(*Sessions).AdvanceShapeTask` calls `RepairTask`, then acts on the pure `loopPosition` of the records: plan, plan review, developer, pull request, four statuses,
  waiting-spec acceptance, verifier, verdict; ended stages are finished from their kept files; cycle cap 1 (`shapeloop.go`). Not delivered: `shape` (66h), scheduling and the
  repair of a push (66k), the close-out (66k children).
- [x] **Sessions that read have no push**, and `TaskSession` refuses a prompt file of another task or base before the sweep (`session.go`).
- [x] **The base of a verifier** is made reachable by `git.UpdateRef` and a scratch branch `layup/base/<task>/<attempt>` that the session and its admission probe clone.
- [x] **Docs.** `loop.md`, `session.md`, `packages.md`, row 66r of the plan, PRD §12 cells, traceability, two entries of `guardrails.md` §2.

## Verification and verdict

Test first (stub at `b538f03`), 34 compiled mutants killed with hashes and controls, one equivalent survivor (SM19b). Full integration level green on the final code head
(2026-10-19T14:56:17Z to 15:23:56Z, exit 0); gitleaks exit 0 before each push. Finding off the path: the architecture task cannot pass `layup/spec` on a target with no PRD,
issue [#432](https://github.com/pharzam/layup/issues/432), row 66h.

Round 1 (GPT-6.1 Sol, xhigh, Codex), comment 6267186962: **not mergeable** (admission probe lost the scratch branch; recovery tests; SM18; budget), fixed. Round 2, comment 6268303209:
**not mergeable** (a bound outside the developer role; mutant hashes; budget); the Operator raised the cap to 2 (6268912691), fixed. Round 3, comment 6269963495:
**not mergeable** (the status refusals untested; budget). The Operator's end 3 (6269982316): the tests of finding 2 landed with no further round, so **no round read the end-3 fix**.
**Verdict: not mergeable until the Operator approves the exact landing figure** (26 files before this close-out, over 25). Harness exception: ADR-0026's routing rule names
`claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 | not recorded | not kept | not kept (posted 2026-10-18T12:57:03Z) |
| Plan review | reasoning | GPT-6.1 Sol | xhigh | 183,423 | 2026-10-18T12:57:12Z to 13:04:43Z |
| Code, three fix cycles and the end 3 fix | execution | Claude Sonnet 5.5 | not recorded | not reported by the harness | 2026-10-19T10:12Z (first commit) to about 15:30Z |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 265,594 | 2026-10-19T12:15:56Z to 12:26:51Z |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 224,559 | 13:13:07Z to 13:24:38Z |
| Review round 3 | reasoning | GPT-6.1 Sol | xhigh | 224,455 | 14:35:49Z to 14:43:24Z |
