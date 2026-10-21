# T-xxr8 — Go 1.26.10

Issue [#457](https://github.com/pharzam/layup/issues/457), the standard-library vulnerability blocking the required `security` job. Serves F-0003#47, Verification on Every Change; goal count **1**.
Claim [6308981059](https://github.com/pharzam/layup/issues/457#issuecomment-6308981059); plan [6309002166](https://github.com/pharzam/layup/issues/457#issuecomment-6309002166); independent plan review [6309035420](https://github.com/pharzam/layup/issues/457#issuecomment-6309035420) (Claude Opus 5.5 on OpenCode, approve-with-conditions); amendment 1 [6309104054](https://github.com/pharzam/layup/issues/457#issuecomment-6309104054). The one condition is met.
Maximum **100 added-plus-removed lines / 6 files**, against current main at landing (baseline `f9d46b0f441202254e8ea6ca16fe4ad8ec860858`); tolerance **60/4**, limit **160/10**; cycle cap **1**, no panel. Evidence and close-out bookkeeping leave the count under ADR-0028.
Evidence: [test runs](../../runs/T-xxr8/test-runs.md), [conditions](../../runs/T-xxr8/conditions.md). The task adds no executable check.

## Delivered and boundaries

- [x] The reviewed plan and amendment name the three version edits and preserve the historical records.
- [x] `go.mod` names `go 1.26.10`; CI continues to read that single source through `go-version-file: go.mod`.
- [x] `docs/spec/setup.md` names Go 1.26.10 since T-xxr8; the guardrail check's full-patch example follows it. The catalog's minor line remains `1.26`, and the T-w89c advisory history and pre-existing records are preserved.
- [x] **Host step:** the coordinating session set the global mise pin with `mise use -g go@1.26.10`, recording `go version` before and after. This is a host change, outside the repository.
- [x] Build, vet and all three test levels pass on Go 1.26.10 in the supplied host runs; govulncheck is clean and PR #459's hosted `security` job passes.
- **PATH and GOROOT trap:** PATH selected Go 1.26.10 while inherited `GOROOT` still named the 1.26.9 install. `go version` alone looked current; compilation and the host commit hook failed with a compiler/tool version mismatch. The sandbox rerun removed inherited `GOROOT` and used a writable module cache; the successful host reruns set both PATH and GOROOT to the 1.26.10 install.
- Outside delivery: application code, catalog changes and additional CI version settings. Final required checks and the landing-figure post remain for the coordinator before merge.

## Verification and verdict

**Red (R8).** Before the directive or host toolchain changed, the coordinator ran `go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./...` on the baseline with Go 1.26.9 at 2026-10-21T05:42:25Z. It found GO-2026-6809 in `crypto/x509@go1.26.9`, fixed in `go1.26.10`. The supplied shell exit is **1**; govulncheck's exit **3** is reported by `go run` as `exit status 3`. No red end time was supplied.

**Host green.** On 2026-10-21, 05:53:12Z–05:57:47Z, the coordinator's Go 1.26.10 runs passed build, vet, uncached unit, integration and E2E levels, then the same govulncheck command exited **0**, reporting `No vulnerabilities found.` Integration used `-count=1 -tags=integration -timeout 30m` and took 202 s; E2E used `-count=1 -tags=e2e -timeout 10m` and took 66 s. These are supplied host results. Earlier sandbox integration and E2E attempts failed because local listening sockets were prohibited; they remain failures in the evidence, separate from the host greens.

**Hosted security.** For [PR #459](https://github.com/pharzam/layup/pull/459), [run 39552815927](https://github.com/pharzam/layup/actions/runs/39552815927), [job 119065633629](https://github.com/pharzam/layup/actions/runs/39552815927/job/119065633629) (`security (govulncheck v1.8.0, go vet, and gitleaks)`) is completed with conclusion **success** on `cb02a1a6211e4a227dc8e993876c948d7c73bdc4`. It ran 06:00:45Z–06:01:47Z; dependency, static-analysis and secret-scan steps succeeded. Stage C fetched the jobs API through the GitHub connector. The failed local `gh pr checks` supplied no result; this record claims no whole-run conclusion.

**Gates and pre-freeze.** Stage A passed ADR/PRD/link lints, 138 discipline cases, 10 nested-checkout cases, 85 freeze-check fixtures, full setup and whitespace checks. Stage C ran the task's clean-head freeze check on `cb02a1a` at 06:07:36.183049Z–06:07:36.897998Z, exit **0**, with `--base origin/main --budget 100/6` and both task records declared as evidence. K1–K6 passed; K5 projected **8/3**, maximum 100/6, tolerance 60/4, limit 160/10, margin **152/7**. The complete invocation and output are in the test-runs record. This projection does not replace measurement at landing.

**Round 1** ([6309308173](https://github.com/pharzam/layup/issues/457#issuecomment-6309308173)): Claude Opus 5.5 on OpenCode, variant high, cycle **0**, **material** on `cb02a1a`. F1 was a records finding: the frozen head lacked the pre-freeze-check record. Fix `8212543` records the check and hosted security result in the two evidence files.

**Round 2** ([6309391860](https://github.com/pharzam/layup/issues/457#issuecomment-6309391860)): a fresh Claude Opus 5.5 session on OpenCode, variant high, cycle **1**, **nothing material in scope** on `82125436dc5a4f264a055abbd3bc7ef9702211e7`. F1 is settled; no raw finding remains. The reviewer reran the clean-head freeze check at that head: exit **0**, all K1–K6 passed with the same output, satisfying the plan's clean rerun. The reviewer also reports green build, vet, unit, setup, freeze-check fixtures and whitespace checks; host integration/E2E and the scan remain attributed to the supplied evidence.

**Note dispositions.** Round 1's hosted-security note is applied by `8212543`. Round 2's clean-rerun note is addressed by its own recorded run above. The optional guardrails addition for PATH/GOROOT (round 1 N2, round 2 N2) is declined under the guardrails §2 filter: the mismatch fails loudly and depends on the host environment; the evidence and this task file retain the lesson. This close-out changes no reviewed guardrail. The task file and completed line satisfy round 2 N3; the dated evidence-section order (N4) is retained as cosmetic.

**Close-out validation and landing audit.** Full setup passes all twelve checks and its kit linters, including ADR/PRD/link lints (3,033 links) and 138 discipline cases. The 10 nested-checkout cases, 85 freeze-check fixtures and `git diff --check` pass. A direct check confirms only the two planned files changed, one completed entry was appended, both files have exactly one final newline with no trailing whitespace or carriage returns, and resource sums and durations agree. Against `origin/main` at `f9d46b0`, the existing counted change is **8/3**; conservatively counting all **35** task lines before `## Verdict` and both close-out files gives **43 lines / 5 files**, margin **57/1** to maximum 100/6 and **117/5** to limit 160/10. The verdict/resource suffix, completed line and two declared evidence records leave the line count. No backlog line exists. The coordinator must post the figure and remeasure if main moves; no check on a future committed close-out head is claimed.

## Verdict

**Nothing material in scope**, the round-2 verdict on `8212543`, retained at cap 1. The directive selects Go 1.26.10, the host levels and scan pass, and the required hosted security job passes on the earlier recorded head. Only close-out bookkeeping follows; no unread implementation fix is delivered. No commit, push or merge is claimed from this session.

## Resource record

[ADR-0007](../adr/0007-record-task-resource-use.md), **recorded, not budgeted**. All intervals are on 2026-10-21, UTC; durations are derived from the supplied endpoints. Tokens are supplied harness reports; OpenCode includes cache and is kept separate. Only OpenCode reports USD cost. Unreported effort and this stage's figures remain `not reported`.

Harness exception: ADR-0026 names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13, as earlier task files state. GPT-6.1 Sol on Codex CLI authors the plan, implementation and close-out. Its use on execution parts is a reasoning-tier model carrying out the fixed plan; no execution-tier model is reported. OpenCode needs no harness exception; the reason earlier reviewer harnesses were passed over is not reported.

| Part | Expected tier | Model and harness | Effort | Tokens | Elapsed (UTC) |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | GPT-6.1 Sol on Codex CLI | not reported | 50,386 | 05:42:25Z–05:43:42Z; 1m 17s |
| Plan review (`issue457-plan`) | reasoning | Claude Opus 5.5 on OpenCode | variant high | 749,920 with cache; USD 0.61 | 05:44:00Z–05:46:21Z; 2m 21s |
| Stage A: version edits and sandbox verification | execution | GPT-6.1 Sol on Codex CLI | not reported | 131,868 | 05:47:07Z–05:51:32Z; 4m 25s |
| Stage B: host-green records | execution | GPT-6.1 Sol on Codex CLI | not reported | 45,840 | 05:58:01Z–05:59:50Z; 1m 49s |
| Round 1 (`issue457-r1`) | reasoning | Claude Opus 5.5 on OpenCode | variant high | 1,805,640 with cache; USD 1.04 | 06:01:40Z–06:06:37Z; 4m 57s |
| Stage C: F1 records fix, freeze check and hosted security | execution | GPT-6.1 Sol on Codex CLI | not reported | 75,023 | 06:06:57Z–06:09:30Z; 2m 33s |
| Round 2 (`issue457-r2`) | reasoning | Claude Opus 5.5 on OpenCode | variant high | 415,063 with cache; USD 0.38 | 06:10:23Z–06:12:35Z; 2m 12s |
| Stage D: close-out | — | GPT-6.1 Sol on Codex CLI | not reported | not reported | not reported |
| Coordinating session: host step, forge and commit-plan commits | — | Claude Opus 5.5 on Claude Code | not reported | not kept | not reported |
| **Total reported** | | | | **303,117 Codex; 2,970,623 OpenCode with cache; USD 2.03** | **19m 34s summed measured parts: 10m 04s Codex + 9m 30s OpenCode; Stage D and coordination not reported** |

The coordinating session wrote none of the reviewed change. This stage leaves only the task detail and completed-log line for its single commit; the sandbox's read-only `.git` prevents committing here.
