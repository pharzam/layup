# T-ebm4 — reconcile the budget and review summaries

Issue [#451](https://github.com/pharzam/layup/issues/451); refs #447, #449 and #441. Serves F-0003#48; goal classes **2**, reconciliation and the Operator-assigned time to overrule, as confirmed by the plan review.
Supplied ordered plan and independent plan review: Claude Opus 5.5 on OpenCode, `approve-with-conditions`; all three conditions are met under amendment 1.
Maximum **400 added-plus-removed lines / 10 files**, against current main at landing, close-out inside under ADR-0028; tolerance **80/4**, limit **480/14**; cycle cap **1**, no panel.
Evidence: [test runs](../../runs/T-ebm4/test-runs.md), [conditions](../../runs/T-ebm4/conditions.md), [occurrence audit](../../runs/T-ebm4/mentions.tsv), [author semantic comparison](../../runs/T-ebm4/semantic.md) and [task check](../../runs/T-ebm4/checks.sh). Only the four declared Markdown/TSV evidence records leave the count; `checks.sh` remains counted.

## Delivered and boundaries

- [x] The reviewed, test-first plan precedes the red check. R12, gate step 5, the named-base pitfall and the Size class, Issue split, Material and Frozen head rows now point to the canonical protocol where their summaries admitted competing readings.
- [x] Records finding, Tolerance, Input, Margin and Coordinated Universal Time have glossary rows; Records finding retains "and changes nothing else", and Cycle cap and Ceiling have concrete examples.
- [x] The canonical section uses straight quotation characters without changing the existing count exclusions, input or margin definitions, end-3 eligibility/exclusions or ends 1, 2 and 4.
- [x] One shared waiting rule requires a new issue comment naming the delegated figure or end; edits post nothing. The posting time selects the window, with both delegations and all four boundary cases compared in the semantic record; end 3 points to the shared rule.
- [x] The full validation section stays unchanged under condition 2. The occurrence audit classifies historical records, other meanings and the README/onboarding/glossary summaries; the existing reference-sweep pitfalls cover the lesson without a duplicate rule.
- [x] Red, green, local gates, clean committed-head K1-K6 and independent round 1 are recorded. The supplied round reports nothing material in scope; both notes are disposed of below.
- Boundaries: no existing check script, hook, CI or product code changes. N4's proposed edits to `runs/T-3ff9/checks.sh` are declined as agreed in the plan review: it records #447's contract; this task's new check covers the replacements. The new waiting rule takes effect only when #451 lands and supplies no authority for its own landing.
- Final hosted checks, posting the landing figure on #451 and any remeasurement if main moves remain for the coordinator. This stage writes only this task record and the completed-log line; no backlog line exists.

## Verification and verdict

**Red (R8).** On unchanged `abdbe996282b33080406223ae8ba19e5a3634360`, `sh runs/T-ebm4/checks.sh` exits **1** at 2026-10-21T07:27:19Z. It reaches every assertion and detects the old summaries, missing rows, missing wait and curly quotation characters; the protected rules pass. This is an assertion failure, not a setup error.

**Green.** The identical check exits **0** at 2026-10-21T07:36:20Z and again at 07:44:17Z. Stage B records passing ADR/PRD/link lints, 138 discipline cases, 10 nested-checkout cases, 85 freeze-check cases, full setup and whitespace. The author comparison prepares the semantic review; it does not substitute for it.

**Pre-freeze and base.** The coordinator merges current main as `1c3ae2f1ac5019adc560e6c455064899727caf31`. Stage C uses `origin/main`, resolved to **`be6a0c7c199c4036c62933b45d38055e16e80681`**, under the plan review's current-main-at-landing condition. On that clean committed merge head, the task check, link lint and full setup pass; `sh docs/tests/freeze-check.sh --task T-ebm4 --base origin/main --budget 400/10 --evidence runs/T-ebm4/conditions.md --evidence runs/T-ebm4/semantic.md --evidence runs/T-ebm4/test-runs.md --evidence runs/T-ebm4/mentions.tsv` runs at 2026-10-21T07:47:47Z–07:47:48Z, exit **0**. K1-K6 pass; K5 reports **232/5**, maximum 400/10, tolerance 80/4, limit 480/14, margin **248/9**. The original-base projection was also 232/5 before the merge. Evidence commit `48dffc7` follows; the run does not claim to check that later commit.

**Round 1.** The supplied independent blind round by Claude Opus 5.5 (`anthropic/claude-opus-5-5`) on OpenCode reports **nothing material in scope** on **`48dffc72c6c27d1e171025b71960d97c8af1d222`**. Two notes, no material findings; notes cause no further round within cap 1. The supplied plan review and round are the authority for their results; this session does not claim a fresh retrieval.

**Note dispositions.** Both are applied within this record, with no correction to a reviewed file.

- **N1, task timing and landing count:** applied. The task file was deferred to close-out, as in T-3ff9, T-4naq and T-jpkc, despite plan step 2 placing its initial record before freeze. K5's 232/5 therefore omitted it. The landing audit below counts every line before `## Verdict` and the task file as touched; it also conservatively counts the completed-log file as touched.
- **N2, old T-3ff9 check:** applied as a disclosure. The N4 decline leaves `runs/T-3ff9/checks.sh` failing: lines 105-109 already fail on main because the named K5 paths changed between `f45fffc` and `be6a0c7`; this task also leaves the curly-quote assertions at lines 54, 66 and 68-71 and the old Cycle cap/Ceiling examples at lines 88-89 unmatched. T-3ff9's promised "matching assertions" are supplied by this task's new check, not by repairing the old script. No fresh run or green result is claimed for that historical check.

**Close-out validation and landing audit.** Against current main `be6a0c7`, the existing counted diff is **232/5**. This record adds **35** counted lines before `## Verdict`, giving **267 lines / 7 files**, with both close-out files counted as touched: margin **133/3** to maximum 400/10 and **213/7** to limit 480/14. The verdict/resource suffix, completed line and four declared evidence records leave the line count. No committed close-out-head check is claimed. `sh runs/T-ebm4/checks.sh` passes at 2026-10-21T08:19:33Z–08:19:33Z, exit 0. `sh .githooks/pre-commit` passes at 08:19:33Z–08:19:57Z, exit 0: ADR/PRD/link lints (3,071 links), 138 discipline cases, adapted/markers, formatting, Go vet and unit tests. `git diff --check` passes at 08:19:57Z–08:19:57Z, exit 0. This evidence sentence is filled afterward; command output stays in scratch to preserve the two-file close-out boundary.

## Verdict

**Nothing material in scope**, the round-1 verdict on `48dffc7`, retained without an implementation fix or another round. Only the two permitted close-out files follow. Merge remains subject to final required checks and posting the audited figure on #451. No commit, push, issue mutation or merge is claimed from this session.

## Resource record

[ADR-0007](../adr/0007-record-task-resource-use.md), **recorded, not budgeted**. UTC intervals and tokens are supplied harness reports; durations are derived from their endpoints. OpenCode tokens include cache and are kept separate; only OpenCode reports USD cost. Effort settings not supplied are `not reported`.

Author assignment: the Operator's supplied build instructions assign GPT-6.1 Sol on Codex CLI to the plan, implementation and close-out, and Claude Opus 5.5 on OpenCode to independent review. Harness exception: ADR-0026 names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13, as earlier task files state. OpenCode needs no harness exception; the reason earlier reviewer harnesses were passed over is not reported.

| Part | Expected tier | Model and harness | Effort | Tokens | Elapsed (UTC) |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | GPT-6.1 Sol on Codex CLI | not reported | 166,780 | 2026-10-21T06:46:25Z–06:55:45Z; 9m 20s |
| Plan review, issue451-plan | reasoning | Claude Opus 5.5 on OpenCode | variant max | 6,831,962 with cache; USD 4.67 | 2026-10-21T07:00:17Z–07:20:36Z; 20m 19s |
| Stage A: check and red evidence | execution | GPT-6.1 Sol on Codex CLI | not reported | 124,192 | 2026-10-21T07:22:27Z–07:32:12Z; 9m 45s |
| Stage B: wording, glossary and green evidence | execution | GPT-6.1 Sol on Codex CLI | not reported | 163,396 | 2026-10-21T07:32:54Z–07:44:57Z; 12m 03s |
| Stage C: merged-main checks and pre-freeze | — | GPT-6.1 Sol on Codex CLI | not reported | 74,049 | 2026-10-21T07:46:20Z–07:49:05Z; 2m 45s |
| Round 1, issue451-r1 | reasoning | Claude Opus 5.5 on OpenCode | not reported | 16,816,344 with cache; USD 7.49 | 2026-10-21T07:50:18Z–08:15:22Z; 25m 04s |
| Stage D: close-out | — | GPT-6.1 Sol on Codex CLI | not reported | not reported | not reported |
| Coordinating session | — | Claude Opus 5.5 on Claude Code | not reported | not kept | not reported |
| **Total reported** | | | | **528,417 Codex; 23,648,306 OpenCode with cache; USD 12.16** | **79m 16s summed measured parts: 33m 53s Codex + 45m 23s OpenCode; Stage D and coordination not reported** |

This stage leaves only the task detail and completed-log line for the coordinator's single commit; the sandbox's read-only `.git` prevents committing here.
