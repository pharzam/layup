# Evaluation: Loki Mode

| Field | Value |
| ----- | ----- |
| Repository | https://github.com/asklokesh/loki-mode |
| Pinned commit | b60ca0ef35273b3553ca99ec530011ae0d33413f (2026-09-28T00:46:29-04:00, v10.1.0) |
| License | Business Source License 1.1 (LICENSE file). Licensor Autonomi, Inc. Additional Use Grant narrowed by COMMERCIAL-TERMS.md. Change Date March 19, 2030; Change License Apache-2.0. Not an open-source license. |
| Language and needs | Bash engine (`autonomy/run.sh` 29,513 lines, `autonomy/loki` 37,789 lines), inline Python 3, `jq`, `git`, Node/npx (optional bootstrap steps), a coding-agent CLI (Claude Code by default; also Codex, Cline, aider, OpenCode). Optional Bun, Docker, dashboard server on `localhost:57374`. State in `.loki/` in the target (ignored by `.loki/.gitignore`) and in `~/.loki/`. Telemetry to PostHog on by default for interactive use. |
| Evaluated by | Claude Opus 5.5 on Claude Code, 2026-09-28 |
| Elapsed | about 40 minutes of the session (shared with three other candidates) |
| Agent runs and cost | 3 `loki start` runs (Claude Code, forced to `claude-haiku-4-5-20251001` by a wrapper, spend cap USD 0.60 per call, `LOKI_BUDGET_LIMIT=3.00`). Run 1 stopped by us after 2 minutes (see the incident). Run 2 stopped by our guard after 7 s (a false alarm). Run 3 completed: main agent USD 0.11 (`.loki/metrics/result-cost-1.json`); reviewer and council calls not costed by Loki. |

## Verdict

Reject. The license forbids LAYUP's use case: "compete" includes "multi-agent AI development orchestration platforms or frameworks" and anything that "substantially replicate[s] ... the Completion Council peer review system", and the licensor claims that even an independent reimplementation from its documentation is under the license. The run also changed the host's global Claude Code setup by default (a plugin install), and several gates pass when they did not run.

## What it is (from the code)

An autonomous build loop (RARV: reason, act, reflect, verify) that runs a coding-agent CLI in iterations, runs quality gates after each iteration, and asks a completion council whether to stop.

- **License terms.** BUSL-1.1 grants non-production use ([LICENSE#L42](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/LICENSE#L42)). Production use is granted only as COMMERCIAL-TERMS.md allows ([LICENSE#L13](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/LICENSE#L13)). Each version becomes Apache-2.0 on March 19, 2030 or on the fourth anniversary of its first release, whichever comes first ([LICENSE#L21-L22](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/LICENSE#L21-L22), [#L46](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/LICENSE#L46)); for v10.1.0 that is March 19, 2030. COMMERCIAL-TERMS.md permits individual, internal-tool, academic use and "(iv) evaluation and testing in non-production environments" ([#L38-L41](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/COMMERCIAL-TERMS.md#L38-L41)). It defines "compete" to include "(a) multi-agent AI development orchestration platforms or frameworks" and "(d) ... substantially replicate ... the RARV ... cycle, Completion Council peer review system", "whether for a fee or free of charge" ([#L21-L32](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/COMMERCIAL-TERMS.md#L21-L32)). Its IP notice says "independent reimplementation based on the Licensed Work's documentation" is "subject to the terms of the License" ([#L45-L54](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/COMMERCIAL-TERMS.md#L45-L54)).
- **The 8 gates.** The gate document lists 8 blocking default-on gates, 3 advisory, 1 opt-in, each with an opt-out flag; gate 8 (Magic Modules debate) is itself advisory by default ([skills/quality-gates.md#L5-L35](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/skills/quality-gates.md#L5-L35)). So 7 gates block by default.
- **Test-integrity gates 5 and 6** scan JavaScript, TypeScript, Python and shell files only ([tests/detect-test-mutations.sh#L139-L146](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/tests/detect-test-mutations.sh#L139-L146)). On a Go target the mutation gate prints PASS ([autonomy/run.sh#L15356](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/autonomy/run.sh#L15356)) with nothing scanned.
- **Blind code review (gate 3).** Selected reviewers run as separate CLI processes in parallel ([autonomy/run.sh#L17513-L17516](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/autonomy/run.sh#L17513-L17516)). If no reviewer returns a usable verdict, the review is inconclusive and blocks ([autonomy/run.sh#L18258-L18264](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/autonomy/run.sh#L18258-L18264)). One reviewer's Critical or High finding blocks.
- **Completion council, default path.** One `claude --agents <json> --json-schema` call; the parent session runs the three declared voter agents and returns one JSON with all findings ([autonomy/lib/voter-agents.sh#L290](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/autonomy/lib/voter-agents.sh#L290), [#L379-L380](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/autonomy/lib/voter-agents.sh#L379-L380)). The members may be isolated from each other as subagents, but one session writes all three votes.
- **Council v2 (opt-in, `LOKI_COUNCIL_VERSION=2`).** Each member gets a copy of the evidence in its own temporary directory and runs in its own process ([autonomy/council-v2.sh#L83-L112](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/autonomy/council-v2.sh#L83-L112)). A member with no parseable verdict is INCONCLUSIVE, never REJECT ([#L138-L160](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/autonomy/council-v2.sh#L138-L160)); approval needs ceil(2/3) of the council size, so a failed member is never an approval ([#L241-L253](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/autonomy/council-v2.sh#L241-L253)). A unanimous approval adds a devil's-advocate member.
- **Evidence gate pass-through.** When the council's evidence is inconclusive, the gate logs a warning and passes, including "completion not backed by test evidence" ([autonomy/completion-council.sh#L2730-L2741](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/autonomy/completion-council.sh#L2730-L2741)). The checklist hard gate fails closed when Python is absent ([#L1397-L1420](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/autonomy/completion-council.sh#L1397-L1420)).
- **Stuck detection.** Escalates when 2 of 3 proxies are hot for N rounds: no change, diff-hash oscillation, and a persistent council split ([skills/quality-gates.md#L502-L586](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/skills/quality-gates.md#L502-L586)). The action is a handoff file, a notification and a PAUSE file. In the default `perpetual` mode the PAUSE is cleared automatically, so it is "notify-only".
- **Host side effects.** The caveman bootstrap is on by default for the Claude provider. It runs an upstream installer with `npx` from GitHub, which "adds a SessionStart hook ... that affects EVERY Claude Code session on this machine" ([autonomy/lib/claude-flags.sh#L772-L835](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/autonomy/lib/claude-flags.sh#L772-L835)). `loki` also writes the trust flag into `~/.claude.json` when the CLI fails on workspace trust ([autonomy/loki#L27854-L27870](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/autonomy/loki#L27854-L27870)). Telemetry posts to PostHog unless opted out or non-interactive ([autonomy/telemetry.sh#L1-L12](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/autonomy/telemetry.sh#L1-L12), [#L55-L56](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/autonomy/telemetry.sh#L55-L56)).

## What we ran

Local evaluation only, which clause (iv) permits. No install: the CLI ran from the clone (`autonomy/loki`), with `HOME` set to the work directory. A wrapper named `claude` on `PATH` gave the CLI the host login (real `HOME`), forced the haiku model and a spend cap, and logged each call. Telemetry: our first commands ran without a TTY, which turns telemetry off; the runs set `DO_NOT_TRACK=1`. No telemetry id file was created. The toy target: Go module `calc` with a local bare remote, and a branch `flawed` with a wrong `Abs`, a changed `Div`, and a deleted assertion (tests still pass).

| Step | Command | Result |
| ---- | ------- | ------ |
| Doctor | `loki doctor` | `17 passed, 0 failed, 14 warnings` |
| Deterministic verify on `flawed` | `loki verify main --no-llm --explain` | `build pass`, `tests pass`, `static_analysis skipped (no lint-capable changed files)`, `nomock skipped`, `secret_scan pass`, `dependency_audit skipped`, `spec_drift skipped` -> `VERDICT: VERIFIED`, exit 0. The wrong code and the deleted assertion pass. `go vet` did not run. |
| Review on `flawed` | `loki review --since main` | `No findings. Code looks clean.`, exit 0 ("No AI needed": pattern checks only) |
| Gates 5 and 6 on the toy | `LOKI_SCAN_DIR=toy detect-test-mutations.sh --block-high`; `detect-mock-problems.sh` | `Results: 0 finding(s)`, `All tests pass mutation detection gate.`, exit 0 for both. No Go file was scanned. |
| Run 1 | `loki start prd.md --simple --no-dashboard --provider claude` (PRD: add `Abs` with a table test, keep `Div`) | Stopped by us: the caveman bootstrap installed a plugin into the host's real `~/.claude` (see below). |
| Run 2 | same, `LOKI_CAVEMAN=0 LOKI_CAVEMAN_AUTO_BOOTSTRAP=0` | Stopped by our guard at 15:49:26: `known_marketplaces.json` changed. The change was only a `lastUpdated` time stamp that Claude Code writes at start. |
| Run 3 | same, guard limited to settings and plugin sets | Completed in 4 minutes, 11 CLI calls. The agent committed `Abs` and a table test. Loop gates: `Static analysis ... recording as not run, not as a pass`, `Test suite gate: go-test passed`, `Mock integrity gate: no applicable tests found; gate did not run`, `Mutation integrity gate: PASS`, `LSP diagnostics ... gate did not run`. Code review: 4 reviewers dispatched as 4 separate calls, `3/4 PASS, 1 FAIL - no blocking issues`. Council: five evidence checks `not confirmed ... Pass-through`, then `round-1.json`: 3 of 3 `COMPLETE` (`"source": "voter-agents-dispatch"`, one call). `COMPLETION COUNCIL: PROJECT APPROVED`. Loki then committed `HANDOFF.md` on branch `loki/session-...`. |

`.loki/metrics/budget.json` read `"budget_used": 0.0` while `result-cost-1.json` read `"total_cost_usd": 0.11`: the budget guard did not count the spend of this run.

**Incident in run 1 (host change, reverted).** At 15:46:45-15:46:52 the caveman bootstrap ran `claude plugin marketplace add` and `claude plugin install` through our wrapper, which used the real `HOME` for the login. It added the `caveman` marketplace and the plugin `caveman@caveman` v2.7.0 at user scope, enabled it in `~/.claude/settings.json`, and the plugin's hook created `~/.claude/.caveman-sessions` and `~/.claude/.caveman-mode-log.jsonl`. Loki printed "bootstrap unavailable ... run proceeds uncompressed" although the install had succeeded. We stopped the run at about 15:48, ran `claude plugin uninstall caveman@caveman --scope user` and `claude plugin marketplace remove caveman`, and moved the plugin cache and the two state files out of `~/.claude` into the work directory. A diff of `settings.json` before and after the revert shows only the two caveman entries removed. Two Claude Code sessions started between 15:46:52 and the revert (session ids in the moved `.caveman-sessions`) may have run with the caveman hook.

## In-Scope items S1-S12

| Item | Mark | Evidence |
| ---- | ---- | -------- |
| S1 Problem Statement Quality | partly | Run 3: `PRD analysis complete: score=3.3/10`, a "spec interrogation" step that records ambiguities as assumptions. It does not ask the idea owner one batch of questions; it assumes. |
| S2 Reproducible Discipline Setup | no | No baseline; setup writes `.loki/` and host-level files. |
| S3 Rule Protection | no | Each gate has an opt-out flag; project config `.loki/config.yaml` is read from the workspace ([autonomy/run.sh#L283-L284](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/autonomy/run.sh#L283-L284)); the agent runs with `--dangerously-skip-permissions`. |
| S4 Stack-Dependent Gates | partly | Build and test runners by stack (`go build`, `go test` ran). Gates 5, 6 and static analysis do not cover Go. |
| S5 Role Handoffs | partly | Structured handoff JSON (`.loki/memory/handoffs/*.json`), review and vote JSON. No schema validation between roles was seen. |
| S6 Autonomous Clarification | no | Ambiguities become recorded assumptions, not answers from a role. |
| S7 Verification on Every Change | partly | Gates run after every iteration. On the toy, `loki verify` returned VERIFIED on a wrong change; the Go-blind gates printed PASS. |
| S8 Human-on-the-Loop | no | "do NOT ask the user questions" is in the agent's system prompt; escalation is notify-only in perpetual mode. |
| S9 Stall Resolution | partly | Stuck detection with 3 proxies, including a persistent council split. No fresh diagnosis; the handoff is written by the loop. |
| S10 Cost Visibility | partly | Per-iteration cost file for the main agent; reviewer and council calls not costed; budget counter read 0.0. Files are ignored by Git. |
| S11 Specification Synthesis | partly | A PRD checklist with verification checks is generated; no trace to a problem-statement text. |
| S12 Harness-Agent Neutrality | partly | Five providers; the council `--agents` path is Claude-only; the code review is "parallel on Claude Code, sequential on other providers" (README). |

## Invariants I1-I9

| Invariant | Effect | Reason |
| --------- | ------ | ------ |
| I1 Git is the system of record | conflicts | `.loki/.gitignore` contains `*`: votes, reviews, costs and state stay out of Git. Learnings go to `~/.loki/learnings`. |
| I2 Independent repository | neutral | The target code is plain Go; Loki adds `USAGE.md`, `HANDOFF.md` and its own commits. |
| I3 Agents cannot change the gates | conflicts | Opt-out flags and a workspace config file; the agent has unrestricted permissions in the same repository. |
| I4 No value without evidence | neutral | Many defaults carry a measured rationale in comments (for example gate 8 made advisory after "3 of 4 personas returned block"), others do not. |
| I5 An inactive check is not a pass | conflicts | Mutation gate PASS with no Go file scanned; evidence gate "Pass-through" on missing test evidence; `loki verify` VERIFIED with 4 of 7 gates skipped. The code review and council v2 do fail closed. |
| I6 Deterministic first | supports | Deterministic gates run before the LLM review and council. |
| I7 Stack adds, never removes | neutral | Gates are selected by file type; not a baseline-rule system. |
| I8 Pinned Armature | neutral | Not related. |
| I9 Replaceable harness | conflicts (partly) | Several providers, but the default council and caveman are Claude-specific, and the run changed the host's Claude Code setup. |

## Deep-check findings it answers

These are mechanisms seen in the code. The license (see the verdict) forbids LAYUP to reuse them, and its IP notice claims reimplementation from the documentation too. The list therefore shows only what the mechanism does.

- **Fable-M18, Sol-24** (panel composition and execution path): council v2 has a size, fixed member roles, per-member isolated evidence directories, parallel processes, an inconclusive-is-not-reject rule, a 2/3 threshold, and a devil's advocate on unanimity. The default council path is one parent session, so it is weaker than it reads.
- **Fable-M4, Author-1** (a disagreement between roles is never a stall): proxy 3, "persistent council split", counts consecutive split verdicts as a stuck signal.
- **Sol-22** (a new record resets the stall clock): proxy 2, diff-hash oscillation at distance, catches A-B-A cycles that each produce a new artifact.
- **Fable-M6** (stall with no diagnosis): not answered; the escalation writes a loop-made handoff and, by default, does not stop.
- **Sol-2, Fable-M19**: not answered. Loki does not open PRs in the runs we saw; a stuck gate ends as `gate_stuck_*` ("Stopped (mutation integrity gate would not clear)", [autonomy/run.sh#L4444](https://github.com/asklokesh/loki-mode/blob/b60ca0ef35273b3553ca99ec530011ae0d33413f/autonomy/run.sh#L4444)), not as a gate review.
- **S7, S12**: partly, above.

## Parts and their verdicts

| Part | Verdict | How LAYUP would take it | Reason |
| ---- | ------- | ----------------------- | ------ |
| The whole orchestrator | reject | none | LAYUP is a multi-agent development orchestrator: the "compete" clause (a) applies once it is offered to anyone, even free. I1, I3, I5 conflicts. |
| Blind code review, council v2 | reject | Take the same idea from MIT sources (AgentJury, no_human) | Clause (d) and the IP notice name the "Completion Council peer review system". |
| Test-integrity gates 5 and 6 | reject | isitdone's scanner covers Go (MIT) | No Go support; PASS when nothing is scanned. |
| `loki verify` | reject | `layup gate` | VERIFIED on a wrong Go change; static analysis skipped for Go. |
| Stuck-detection proxies | reject as code; the idea is generic | Write LAYUP's own rule from the deep-check fixes (Fable-M4, Sol-22) | Same license limit; the proxies are common engineering practice, but do not derive them from Loki's documents. |
| Overall | reject | none | License, host side effects, I5. |

Could LAYUP use it at all? Only as a tool for individual, internal or evaluation use, not inside LAYUP, and not as a model for LAYUP's design while the IP notice stands. On March 19, 2030 this version becomes Apache-2.0.

## Where the searchers were wrong or incomplete

- Searcher B: "8 quality gates" is from the README. The gate document itself says gate 8 is advisory by default, so 7 block by default. Gates 5 and 6 do not scan Go.
- Searcher B: "Audit logs are in the output repo (Invariant 1, partly)". The audit data is in `.loki/`, which `.loki/.gitignore` excludes from Git. We read this as a conflict with I1, not partly.
- Searcher A: "keep: blind completion council and gates" and B's "Use: a pattern for a completion gate and a blind review council". The license text also covers "independent reimplementation based on the Licensed Work's documentation", so "borrow the pattern" from Loki's documents is not safe. Neither searcher read COMMERCIAL-TERMS.md.
- Searcher B: "shows the real cost ... before spending anything". Run 3 printed an estimate (`~$0.46`), but the budget file then read `budget_used 0.0` after USD 0.11 of spend.
- Neither searcher noted the default global plugin install into `~/.claude` or the default-on telemetry.

## Limits of this evaluation

- The council v2 path (`LOKI_COUNCIL_VERSION=2`) was read, not run. The default council path was run once.
- The blind code review ran on a correct change only. We did not get the `flawed` branch through the loop's review, because the loop starts from a PRD, not from a given diff.
- One provider (Claude Code, forced to haiku). Loki's own model routing (sonnet, opus) was overridden by our wrapper.
- `.loki/` files, including the run logs, are kept in the work directory (`eval-work/loki-mode/`) as evidence.
- This is not legal advice. The license reading is ours; a lawyer should confirm it before any use.
