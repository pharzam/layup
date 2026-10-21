# T-szvp — remove the numbered checker restore exception

Issue [#452](https://github.com/pharzam/layup/issues/452), removal of the R4 workaround in [T-c826, #449](https://github.com/pharzam/layup/issues/449); refs #84. Serves F-0003#47, Verification on Every Change; goal count **1**.
Claim [6304156228](https://github.com/pharzam/layup/issues/452#issuecomment-6304156228); supplied plan and independent plan review (Claude Opus 5.5 on OpenCode, approve-with-conditions); amendment 1 [6304288973](https://github.com/pharzam/layup/issues/452#issuecomment-6304288973). The one condition is met.
Maximum **250 added-plus-removed lines / 8 files**, against current main at landing, close-out inside under ADR-0028; tolerance **60/4**, limit **310/12**; cycle cap **2**, no panel.
Evidence: [test runs](../../runs/T-szvp/test-runs.md), [conditions](../../runs/T-szvp/conditions.md), [task check](../../runs/T-szvp/checks.sh). Only the two declared Markdown records leave the evidence count; executable checks remain counted.

## Delivered and boundaries

- [x] The ordered, independently reviewed plan precedes the removal check's red; amendment 1 keeps only one complete Restore-step assertion.
- [x] `nested-checkout-check` restores all nine listed scripts from the default branch for every event. Its authorization comments, PR/event bindings and checker-only keep/continue branch are removed; its comments name no PR number.
- [x] The Restore step matches `d27b33180a1aa82c23f3ab66e757e06de356eb35` byte for byte. Default-branch source, in-place copies, missing-file fallback and both suite invocations remain.
- [x] CI and test READMEs, guardrails and T-c826's migration check describe the removal. T-c826 loses only exception-specific assertions; its six remaining assertions pass without an argument. Its historical evidence is preserved.
- [x] The removal's own [PR #456](https://github.com/pharzam/layup/pull/456) has supplied hosted restoration evidence and a successful suite result, with the provenance limits below. Local gates and pre-freeze K1–K6 pass; round 1 finds nothing material in scope.
- Outside delivery: checker behavior, fixtures and runners, the separate nested-checkout-driver exception, and strengthening setup-check's restore analysis. No product code changes.
- Final required checks on the close-out/landing head, posting the landing figure and opening the revealed-limit issue remain for the coordinator before merge; this record claims none of those actions.

## Verification and verdict

**Red (R8).** On unchanged baseline `f9d46b0`, `sh runs/T-szvp/checks.sh` exits **1** at 2026-10-20T23:39:59Z: both `numbered checker exception remains` and `restore step differs from pre-exception baseline`. The output exposes the actual numbered exception; this is the expected assertion failure.

**Green.** The identical check exits **0** at 2026-10-20T23:44:56Z after removal. T-c826's six retained assertions pass. ADR/PRD/link lints, 138 discipline cases, full setup (including `ci OK`), 85 freeze-check fixtures and whitespace all exit 0 in Stage B. The checker and both runners are unchanged against `f9d46b0`. Stage C's committed comparison against `d27b331` leaves only the planned suite-comment hunk.

**Hosted evidence.** The coordinator supplies successful [run 39517390208](https://github.com/pharzam/layup/actions/runs/39517390208) on PR #456, with `restored from main: docs/tests/freeze-check.sh`, both corresponding runner restore lines, `freeze-check: 85 passed, 0 failed`, and no "kept as written" line. These are supplied observations, not a fresh retrieval. The run's head SHA, tested merge SHA, job URL, timestamps and other six restore lines were not supplied; no `ci-restore.md` exists. Stage C's checkout was `8903e77b298059a7bdc159d36afbc0fbbfbdcd47`, containing removal commit `68931a1876cae7ace9c4cdb80ee88aa64993f039`; that does **not** establish the hosted run's head. The local byte comparison proves removal independently. A further `gh run view 39517390208 --repo pharzam/layup --json url,headSha,event,status,conclusion,createdAt,startedAt,updatedAt,jobs` attempt at 2026-10-21T01:06:59Z–01:06:59Z exits 1, unable to connect to api.github.com.

**Pre-freeze.** On clean `8903e77`, Stage C runs `sh docs/tests/freeze-check.sh --task T-szvp --base origin/main --budget 250/8 --evidence runs/T-szvp/conditions.md --evidence runs/T-szvp/test-runs.md` at 2026-10-20T23:57:35Z, exit **0**. `origin/main` is `f9d46b0f441202254e8ea6ca16fe4ad8ec860858`. All K1–K6 pass; K5 reports **124/6**, maximum 250/8, tolerance 60/4, limit 310/12, margin **186/6**. The coordinator commits the evidence as `ee408ba`; the supplied round reads that head. No additional clean-head rerun is present in the supplied record.

**Round 1.** Claude Opus 5.5 (`anthropic/claude-opus-5-5`) on OpenCode, variant high, **nothing material in scope** on `ee408baa94a3454e4b0544e2075ddf4a159039f2`. Five notes, no material findings; notes cause no further round within cap 2. The first attempt produced no output and was stopped at the 3,600-second timeout; only the successful attempt supplies the verdict.

**Note dispositions.** No reviewed file is corrected during this close-out.

- **N1, hosted provenance:** declined as a completed fix: the tested head is unavailable, and assigning the local SHA to the hosted run would invent evidence. This task record applies the disclosure, but cannot supply the requested SHA or cite an as-yet unreported final run. The coordinator must obtain that identity (containing `68931a1`) or cite the final successful required run on the landing head. The supplied acceptance evidence and its narrower provenance remain explicit.
- **N2, T-c826's stale comment:** declined because changing `runs/T-c826/checks.sh` touches a reviewed file outside this close-out. Its "source and in-place copy" comment overstates the remaining local assertion; T-szvp's byte comparison covers the in-place copy. No assertion or exit status changes here.
- **N3, intermediate red:** applied as a historical disclosure: at `68931a1`, before `8903e77` retires the exception assertions, T-c826's check would fail. CI does not run that check. No fix is requested and published history remains intact.
- **N4, task timing and budget:** applied. The task file arrives with its completed line at close-out, as in T-3ff9 and T-c826, rather than creating an unindexed task before freeze. The final audit below counts this record's non-bookkeeping lines and both touched close-out files conservatively.
- **N5, revealed restore-analysis limit:** declined in this branch as off the removal's DoD. Setup-check `ci` checks that paths are named inside Restore, so it misses a conditional `continue`; its documented scope remains accurate. Proposed follow-up: "setup-check `ci`: a conditional skip inside the Restore loop is not seen"; record the limit in its LIMIT comment and guardrails §2, or detect the skip. Either remedy touches other files. No new issue is claimed; the coordinator must search for duplicates and open or link it.

**Close-out validation and landing audit.** `sh .githooks/pre-commit` passes at 2026-10-21T01:08:51Z–01:09:14Z, exit 0: ADR/PRD/link lints (3,034 links), 138 discipline cases, adapted/markers, formatting, Go vet and unit tests. `git diff --check` passes at 01:09:14Z–01:09:14Z, exit 0. This evidence paragraph is filled afterward. Against `origin/main` at `f9d46b0`, the existing counted diff is **124/6**; this task's **39** lines before `## Verdict` make **163 lines**. Conservatively counting both close-out files gives **8 files**: margin **87/0** to maximum 250/8 and **147/4** to limit 310/12. The verdict/resource suffix, completed line and two declared evidence records leave the line count. No backlog line exists. The coordinator must post this figure and remeasure if main moves; no check is claimed on a future committed close-out head.

## Verdict

**Nothing material in scope**, the round-1 verdict on `ee408ba`, retained without an implementation fix or another round. Only the two permitted close-out files follow. Merge remains subject to final hosted checks, the provenance/follow-up handoffs above and posting the audited landing figure on #452. No commit, push, issue mutation or merge is claimed from this session.

## Resource record

[ADR-0007](../adr/0007-record-task-resource-use.md), **recorded, not budgeted**. UTC intervals and token figures are supplied harness reports; durations are derived from the endpoints. OpenCode tokens include cache and are kept separate; only OpenCode reports USD cost. Effort settings not supplied are `not reported`.

Harness exception: ADR-0026 names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13, as earlier task files state. GPT-6.1 Sol on Codex CLI authors the plan, implementation and close-out. OpenCode needs no harness exception; the reason earlier reviewer harnesses were passed over is not reported.

| Part | Expected tier | Model and harness | Effort | Tokens | Elapsed (UTC) |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | GPT-6.1 Sol on Codex CLI | not reported | 123,813 | 2026-10-20T23:31:21Z–23:34:40Z; 3m 19s |
| Plan review | reasoning | Claude Opus 5.5 on OpenCode | not reported | 1,385,671 with cache; USD 0.97 | 2026-10-20T23:35:02Z–23:38:57Z; 3m 55s |
| Stage A: check and red evidence | execution | GPT-6.1 Sol on Codex CLI | not reported | 63,332 | 2026-10-20T23:39:15Z–23:40:45Z; 1m 30s |
| Stage B: removal, docs and green | execution | GPT-6.1 Sol on Codex CLI | not reported | 99,251 | 2026-10-20T23:41:36Z–23:47:21Z; 5m 45s |
| Stage C: hosted evidence and pre-freeze | — | GPT-6.1 Sol on Codex CLI | not reported | 55,657 | 2026-10-20T23:57:13Z–23:58:49Z; 1m 36s |
| Round 1, first attempt: no output, timeout | reasoning | Claude Opus 5.5 on OpenCode | not reported | 0 | 2026-10-20T23:59:35Z–2026-10-21T00:59:35Z; 60m 00s |
| Round 1, successful attempt | reasoning | Claude Opus 5.5 on OpenCode | variant high | 2,172,600 with cache; USD 1.14 | 2026-10-21T01:01:10Z–01:06:18Z; 5m 08s |
| Stage D: close-out | — | GPT-6.1 Sol on Codex CLI | not reported | not reported | not reported |
| Coordinating session: forge and commit-plan commits | — | not reported | not reported | not kept | not reported |
| **Total reported** | | | | **342,053 Codex; 3,558,271 OpenCode with cache; USD 2.11** | **81m 13s summed measured parts: 12m 10s Codex + 69m 03s OpenCode; Stage D and coordination not reported** |

This stage leaves only the task detail and completed-log line for the coordinator's single commit; the sandbox's read-only `.git` prevents committing here.
