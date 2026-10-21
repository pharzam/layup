# T-ym0k — row 68a of the plan, the review of the release of M2e

Issue [#369](https://github.com/pharzam/layup/issues/369), row 68a of [the tasks of M2e](../plan/README.md#the-tasks-of-m2e), opened by `T-57vx` (#346). Serves `F-0003#51` through REQ-015 and REQ-017. Goal count 1 (O-214). Evidence: [`runs/T-ym0k/`](../../runs/T-ym0k/).
Plan 6314742187 (Claude Opus 5.5, Plan agent); its plan review 6314912248 (GPT-6.1 Sol on Codex CLI: `approve-with-conditions`, seven conditions, D1 to D5, Q1 to Q4, 400 lines over 12 files against `main` at landing, close-out inside, ADR-0028, cap 2, no panel).

## The release

`01d7da3082333efadb227afda14fd98715edf5c4`, the merge of **PR #466** for **issue #368** (row 67, `T-3xvc`, the last build row): 149 paths and 38,708 lines (140 non-test Go files of `cmd/` and `internal/` out of `testdata/`, 38,297 lines; `go.mod`, 3; the eight files of `internal/catalog/go/`, 408). Since M2d the embedded files changed by +112/-1 (`gates.sh.tmpl` +73/-0, `gates.yml.tmpl` +39/-1). Every run names the full ID.

## Decisions (R10), with their reasons

- **D1.** The release is fixed at the full ID: check (1) of the script passes any name, such as `HEAD`.
- **D2.** Check (3) allows the exact call `exec.Command("sh", script, prdDir)` of `internal/spec/prdlint.go` (the pinned prd-lint, which `packages.md` allows as a process seam), not the file: a `-c` form (m1) passes a file-wide allowance and fails this one.
- **D3.** `RemoteHead` joins the regex and the list of check (7): `BranchHead` calls it with an installation token, and it was a silent omission from the listing, not a failure.
- **D4.** Check (6) counts every `.do(` occurrence and the rows that parse, lists the parsed rows only, and fails on a difference or on no occurrence (condition 2). What it does not prove is stated in the header: the form parses, not that the operation or method fits the call.
- **D5.** The comparator `runs/T-ym0k/coverage-check.sh` stays in the tree (counted in the budget) with its test, so a reader can rerun it; the precedent kept it in scratch. It builds both halves from COMMIT (a temporary index made by `git read-tree`), which is the answer to condition 3. It is a shell script, not a new `layup` command.

## Delivered and boundaries

- [x] **The review of the release of `M2e`** (the uat): [`release-review.md`](../../runs/T-ym0k/release-review.md), posted on #369 as comment 6316749856. Claude Fable 5.1 on OpenCode (variant `max`; chosen by the selector, [`selection.txt`](../../runs/T-ym0k/selection.txt) is the selector's `selection.json` byte for byte; 12:30:04Z to 12:53:29Z, USD 41.26, 106,552,552 tokens with cache), a model that wrote none of the release, read the 149 paths in full and recorded `holds` for REQ-015 and REQ-017 with 11 findings, none a defect. The Operator's decisions on the reviewer, its spend and its time of 60 minutes are comment 6315293477 of #369, recorded before the run.
- [x] **The check** (D2 to D4, `runs/T-efmy/release-check.sh`), at the release `01d7da3082333efadb227afda14fd98715edf5c4`: exit 0, `== PASS` ([`release-check.txt`](../../runs/T-ym0k/release-check.txt)); nine mutants beyond the control, each red against the control where the fix is new and refused by the final script on its own line (`test-runs.md`).
- [x] **The comparator** (D5, `runs/T-ym0k/coverage-check.sh` with `coverage-check-test.sh`, 19 cases): the Coverage table of the record matches the commit, 149 rows, each `in full`, exit 0 on the unedited record.
- [x] **S8, by finding** (all are `holds` or `note`; no defect, so no fix in this PR and no issue): findings 1 to 5 and 11 are `holds`; 9 is `holds`, and its one observation, that the two actions are pinned by tag and not by commit SHA, is off the lens of REQ-015 and REQ-017 and is noted here, not opened as an issue; findings 6 (the environment of started programs), 7 (the probe push), 8 (`unsafe`, `syscall`) and 10 (URLs and stand-ins) are notes closed with "No change" and the reviewer's reason.
- [x] **Docs.** `loop.md` (D1 to D5 with their reasons), the plan (the release paragraph and row 68a), `traceability.md` (the uat row, `green`), the PRD-0001 §12 Task and Test cells and a §13 line, and one pitfall of `guardrails.md` §2 (the author set of a release), one clause added to another (an instruction whose shown form the parser does not read).
  Not delivered: row 68b (`T-8s13`, #370), the demo of `M2e` on a real GitHub target, whose inputs the Operator gives; the batch part of `layup/rules` (`M2f`).
- Corrections to the plan's statements (condition 7): the embedded diff is +112/-1; the trailerless blame group also holds `T-3q0b`; Gemini's attempt at `T-2vva` did no work; in m3 the `grep` omits the variable-operation call.

## Verification and verdict

Tests: [`test-runs.md`](../../runs/T-ym0k/test-runs.md); the conditions of the plan review and their evidence: [`conditions.md`](../../runs/T-ym0k/conditions.md). The branch changes no Go file, so the build, vet (three tags), gofmt and the unit level are those of `main`; the six lints, `git diff --check`, the coverage test and the comparator on the record, and the freeze check (K1 to K6, budget 400/12) pass at the head. Not run: the integration and e2e levels (no Go file changed). Known limit: the check matches lines and literals and does not prove the behaviour of an external program; a `.do(` in the form it parses is judged by the reviewer, not by the script.

Round 1 (comment 6315399127, GPT-6.1 Sol, xhigh, Codex): **material**, two holes in `release-check.sh`: check (3) exempted the whole line that held the allowed call, and check (6) matched `a.do` without a boundary (`aa.do`); fixed with red mutants m7 to m9. Round 2 (6315693833, the same reviewer): **material**, the receiver `αa` passed the ASCII boundary; fixed with m10 and a positive boundary on bytes. Round 3 (6315924484, the last at cap 2): the record said `nothing material in scope`, and its verdict was edited to material for the cap raise; the Operator raised the cap to 3 (record 6316752070) so that the comparator, which refused the review's plain `in full`, could be fixed and read. Round 4 (6317380713, Claude Fable 5.1 on OpenCode, `max`, cycle 3 of the cap of 3, 13:14:13Z to 13:32:49Z, USD 9.29, 13,210,116 tokens with cache): **nothing material in scope**, two notes of the records and the rest classifications, applied as documents only (three citations of the conditions table, the Operator's allowance cited by its comment, the link text, "rounds 1 to 3", the doubled path of the traceability row, and the second occurrence in the guardrails entry).
**Verdict: nothing material in scope.** Goal count 1. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13 (the plan review and rounds 1 to 3).

## Resource record

ADR-0007, recorded, not budgeted. UTC, 2026-10-21, from the clock, `start.txt` and `end.txt` of each session; tokens as the harness reported them (Codex "tokens used", Claude Code and OpenCode the run's totals with cache).

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 (Plan agent) | not recorded | not kept | not kept |
| Plan review | reasoning | GPT-6.1 Sol, Codex CLI | xhigh | 199,127 | 11:25:24Z to 11:34:09Z |
| S0 to S5: the red slice, the check, the comparator, the brief | execution | Claude Sonnet 5.5 (subagent) | not recorded | not reported | 11:36Z to 11:46Z |
| Round 1 | reasoning | GPT-6.1 Sol, Codex CLI | xhigh | 160,979 | 11:51:04Z to 12:01:28Z |
| Fix of round 1 | execution | Claude Sonnet 5.5 (subagent) | not recorded | not reported | 12:02Z to 12:05Z |
| Round 2 | reasoning | GPT-6.1 Sol, Codex CLI | xhigh | 149,941 | 12:08:12Z to 12:16:04Z |
| Fix of round 2 | execution | Claude Sonnet 5.5 (subagent) | not recorded | not reported | 12:17Z to 12:19Z |
| Round 3 | reasoning | GPT-6.1 Sol, Codex CLI | xhigh | 177,028 | 12:20:25Z to 12:28:32Z |
| The release review (the uat) | reasoning | Claude Fable 5.1, OpenCode | max | 106,552,552 (USD 41.26) | 12:30:04Z to 12:53:29Z |
| The comparator fix, the review record, S8 and S9 | execution | Claude Sonnet 5.5 (subagent) | not recorded | not reported | 13:05Z to 13:14Z |
| Round 4 | reasoning | Claude Fable 5.1, OpenCode | max | 13,210,116 (USD 9.29) | 13:14:13Z to 13:32:49Z |
| The notes of round 4 and the close-out | execution | Claude Sonnet 5.5 (subagent) | not recorded | not reported | 13:33Z to 13:35Z |
| **Total** | | | | 120,449,743 reported (the numeric entries: 199,127 + 160,979 + 149,941 + 177,028 + 106,552,552 + 13,210,116); the author-tier rows and the plan have no figure | 11:25Z to 13:35Z, with gaps |
