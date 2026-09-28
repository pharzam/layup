# Evaluation: AI-SDLC

| Field | Value |
| ----- | ----- |
| Repository | https://github.com/ai-sdlc-framework/ai-sdlc |
| Pinned commit | 6e5b5461f129d19f8d0cf410d4e0cbaa9a8b4f86 (2026-09-28T07:45:18Z) |
| License | Apache License 2.0 (from the LICENSE file) |
| Language and needs | TypeScript monorepo (1,376 commits). Node.js 22 or later; npm packages `@ai-sdlc/orchestrator` 0.26.1 and `@ai-sdlc/pipeline-cli` 0.26.1 (159 packages, including the native `better-sqlite3`). A Claude Code plugin (`ai-sdlc-plugin`, version 0.21.1: hooks, 13 agent files, an MCP server `@ai-sdlc/plugin-mcp-server` 0.9.2). An authenticated `gh` to open PRs. Optional: Codex CLI for the Codex reviewers, GitHub Actions, tmux for parallel runs. No server. A local SQLite file `.ai-sdlc/state.db` (gitignored) for the issue path. A Go SDK (`sdk-go`) exists, with non-standard-library dependencies. |
| Evaluated by | Claude Opus 5.5 on Claude Code, 2026-09-28 |
| Elapsed | about 20 minutes (15:31 to 15:50 local), in one session shared with sdlc-gh |
| Agent runs and cost | 3 pipeline runs of `/ai-sdlc:execute TOY-1` on `claude-sonnet-5` with `--max-budget-usd 5`: USD 1.06, 0.38 and 1.26 (total USD 2.70). Plus 3 short probes for the harness setup (USD 0.02; two failed at login, USD 0). |

## Verdict

**Borrow the pattern; reject it as a base or a dependency.** The full pipeline ran on a Go toy and produced a signed, valid attestation. But the cross-harness check proves only that a string says "codex" or "claude-code", and a forged envelope with no review passed the verifier as `independent`. The rule protection is a Claude Code hook that a Bash command bypasses. The good parts are designs: the draft-PR-then-ready sequence, a verifier that never loads code from the checkout it checks, a signer key held only by a protected CI job (opt-in), checkpoint commits and quarantine branches, and several record schemas.

## What it is (from the code)

- **Pipeline.** `/ai-sdlc execute` is a 1,967-line slash-command prompt that the main Claude Code session follows: Steps 0 to 15, worktree, developer subagent, three reviewer subagents, verdicts, v6 attestation, push, draft PR, `gh pr ready`. See https://github.com/ai-sdlc-framework/ai-sdlc/blob/6e5b5461f129d19f8d0cf410d4e0cbaa9a8b4f86/ai-sdlc-plugin/commands/execute.md#L1. The steps are LLM-followed text; the scripts it calls are deterministic.
- **Hard rules.** "Never merge", "never force-push", "never edit `.ai-sdlc/**`" (execute.md#L50 to #L66). The merge path is `cli-merge-if-eligible`.
- **Iteration cap.** At most 2 developer iterations on review failure, then a PR with `[needs-human-attention]`. See https://github.com/ai-sdlc-framework/ai-sdlc/blob/6e5b5461f129d19f8d0cf410d4e0cbaa9a8b4f86/ai-sdlc-plugin/commands/execute.md#L1359.
- **Rule protection.** A PreToolUse hook denies `Write` and `Edit` under `.ai-sdlc/**` and under `blockedPaths`, and a raw `gh pr merge`. It inspects `Bash` only for merge, stash and configured patterns. See https://github.com/ai-sdlc-framework/ai-sdlc/blob/6e5b5461f129d19f8d0cf410d4e0cbaa9a8b4f86/ai-sdlc-plugin/hooks/enforce-blocked-actions.js#L162 and #L705. The file itself says: "A determined agent with Bash access can always obfuscate ... the real backstop is branch protection + humans-merge" (#L190).
- **Branch protection body.** `enforce_admins: false`, `require_code_owner_reviews: false`. See https://github.com/ai-sdlc-framework/ai-sdlc/blob/6e5b5461f129d19f8d0cf410d4e0cbaa9a8b4f86/orchestrator/src/cli/commands/branch-protection-shared.ts#L67.
- **Attestation.** `sign-attestation.mjs` and `cli-attestation sign-v6` sign a Merkle root over reviewer transcript hashes with an ed25519 key from `~/.ai-sdlc/signing-key.pem` (https://github.com/ai-sdlc-framework/ai-sdlc/blob/6e5b5461f129d19f8d0cf410d4e0cbaa9a8b4f86/ai-sdlc-plugin/scripts/init-signing-key.mjs#L24). The verifier reads the public keys from `.ai-sdlc/trusted-reviewers.yaml` in the checkout (https://github.com/ai-sdlc-framework/ai-sdlc/blob/6e5b5461f129d19f8d0cf410d4e0cbaa9a8b4f86/pipeline-cli/attestation-core/verify-core.mjs#L2561). The monorepo's own workflow checks out `main` and reads only the envelopes from the PR, so a PR cannot add its own key; the scaffolded consumer workflow does not do this.
- **Verifier trust root.** The consumer verifier loads its runtime only from outside the checkout it verifies and fails closed otherwise. See https://github.com/ai-sdlc-framework/ai-sdlc/blob/6e5b5461f129d19f8d0cf410d4e0cbaa9a8b4f86/ai-sdlc-plugin/scripts/verify-attestation.mjs#L37.
- **Cross-harness check.** Two parts:
  - The v5 path rejects a codex reviewer only when `predicate.harness.name === 'codex'` (https://github.com/ai-sdlc-framework/ai-sdlc/blob/6e5b5461f129d19f8d0cf410d4e0cbaa9a8b4f86/orchestrator/src/runtime/attestations.ts#L1950). If the field is absent (the Claude Code default), no check runs. The value comes from `--harness-name` or the `CODEX_VERSION` environment variable.
  - The v6 verifier compares no harness at all. `harness` is only logged, "forensic/audit only" (verify-core.mjs#L3199; the runbook says the same). Independence is a `verdictClass`: `independent` when a SubagentStart marker file exists in `.ai-sdlc/subagent-sessions/` near the transcript time (https://github.com/ai-sdlc-framework/ai-sdlc/blob/6e5b5461f129d19f8d0cf410d4e0cbaa9a8b4f86/pipeline-cli/src/attestation/verdict-class.ts#L181). The file calls it "best-effort, single-machine".
  - A third tier, `isolated`, needs a signature by a `ciOnly` key that only a protected CI job holds (`.github/workflows/isolated-review.yml`, RFC-0047). It is opt-in and not scaffolded.
  - At dispatch time, `enforceIndependence` removes the implementer's harness from the reviewer's harness chain (https://github.com/ai-sdlc-framework/ai-sdlc/blob/6e5b5461f129d19f8d0cf410d4e0cbaa9a8b4f86/orchestrator/src/harness/independence.ts#L30). This is a routing rule, not a proof.
- **QualityGate.** Gates with `advisory`, `soft-mandatory` (override by a role with a justification) and `hard-mandatory`. A missing metric fails the rule. See https://github.com/ai-sdlc-framework/ai-sdlc/blob/6e5b5461f129d19f8d0cf410d4e0cbaa9a8b4f86/reference/src/policy/enforcement.ts#L78 and #L324. The scaffolded gates check only "description length" and "has acceptance criteria".
- **Orchestrator.** `cli-orchestrator tick` runs an admission chain (open PR, orphan parent, already in flight, dependencies, blast radius, dispatchability, DoR, external dependencies, blocked, captures), then dispatches. A recoverable abort keeps the worktree and its `wip(checkpoint):` commits for the next tick; an unrecoverable one renames the branch to `quarantine/<id>-<time>` (https://github.com/ai-sdlc-framework/ai-sdlc/blob/6e5b5461f129d19f8d0cf410d4e0cbaa9a8b4f86/pipeline-cli/src/orchestrator/rollback.ts#L308). A stuck counter (over 5 ticks) is in memory (`loop.ts#L19`).
- **Questions.** A headless developer calls `cli-decisions escalate`, which writes a decision record and exits 1; the Operator answers with `cli-decisions answer`, then the task resumes (execute.md#L807).
- **Cost.** A `cost_ledger` table (tokens, USD per agent and run) in `.ai-sdlc/state.db` (https://github.com/ai-sdlc-framework/ai-sdlc/blob/6e5b5461f129d19f8d0cf410d4e0cbaa9a8b4f86/orchestrator/src/state/schema.ts#L109), which `init` adds to `.gitignore`.

## What we ran

All runs used `HOME=EVAL_WORK/ai-sdlc/home` for every tool process, no `GH_TOKEN`, and a fake `gh` that logs its arguments and writes nothing to the network. The target was a Go module (`calc.Add`, one test, a Makefile) with a local bare remote.

1. **Install.** `npm install @ai-sdlc/orchestrator@0.26.1 @ai-sdlc/pipeline-cli@0.26.1` into `EVAL_WORK/ai-sdlc/npm` (159 packages; npm held back the `better-sqlite3` install script). The plugin was copied to `EVAL_WORK/ai-sdlc/plugin`, and `install-runtime-deps.sh` added its three runtime packages (278 packages).
2. **Scaffold.** `ai-sdlc init -y --skip-mcp --with-dor --with-attestation --with-workflows --with-classifier`. It wrote 20 files (`.ai-sdlc/*.yaml`, five workflows, `.husky/pre-push`, a `CLAUDE.md` pointer, `.gitignore` for `state.db`).
   - "Attestation infrastructure was scaffolded in AUDIT-ONLY mode. Nothing currently blocks a merge".
   - "Branch protection was NOT applied: gh repo view failed". With `-y`, init tried to apply branch protection through `gh` although the flag was not given. With the host's `gh` login it would have written to GitHub.
   - The scaffolded gate workflow runs `pnpm install`, `pnpm lint`, `pnpm build`, `pnpm test` (`init-templates.ts#L94`). On the Go toy it has no `package.json`, so the single required check `ai-sdlc/pr-ready` would fail on every PR.
3. **Checks.** `ai-sdlc validate`: six resources `VALID`. `ai-sdlc doctor`: "3 pass, 5 warn", including "attestation-governance: Attestation artifacts present but not enforced (audit-only)". `ai-sdlc health`: "Config: valid; State Store: deferred".
4. **Hook test (I3).** The PreToolUse hook script, fed the JSON that Claude Code sends:

   | Tool input | Result |
   | ---------- | ------ |
   | `Write .ai-sdlc/quality-gate.yaml` | deny |
   | `Edit .AI-SDLC/Quality-Gate.yaml` | deny |
   | `Write .github/workflows/ai-sdlc-gate.yml` | deny |
   | `Bash: echo x >> .ai-sdlc/quality-gate.yaml` | allowed |
   | `Bash: sed -i s/hard-mandatory/advisory/ .ai-sdlc/quality-gate.yaml` | allowed |
   | `Bash: rm .ai-sdlc/trusted-reviewers.yaml` | allowed |
   | `Bash: gh pr merge 5 --squash` | deny |
   | `Bash: gh api -X PUT repos/o/r/pulls/5/merge` | allowed |

5. **Agent run 1** (`claude -p "/ai-sdlc:execute TOY-1" --model claude-sonnet-5 --max-budget-usd 5 --plugin-dir plugin --setting-sources project,local --permission-mode bypassPermissions`, 2 min 7 s, USD 1.06). The developer wrote `Sub` and a four-case test; the three reviewers approved. The run stopped at Step 7b.5: no harness transcripts to persist (session persistence was off) and no signing key. The report said: "`harnessNote`: INDEPENDENCE NOT ENFORCED (codex unavailable, fell back to claude-code)". Two reviewers wrote their own "transcript" files into `.ai-sdlc/transcripts/`; the hook blocked the third.
6. **Agent run 2** (same command, USD 0.38). It found the old branch and stopped: "The command tells me to stop and ask you to run `/ai-sdlc cleanup TOY-1`". The slash-command path does not resume from a checkpoint; it asks the Operator.
7. **Agent run 3** (after the evaluator removed the worktree and branch, and added a key; 2 min 42 s, USD 1.26). The full pipeline ran: commits `c4eada2` (feature), `1857256` (task moved to `backlog/completed/`), `462a32a` (attestation); push to the bare remote; `gh pr create --draft`, `gh pr edit`, `gh pr comment` (an incremental-review marker), `gh pr ready 1`. The agent noted that the husky hook never ran (`core.hooksPath` unset), so it ran `.husky/pre-push` by hand.
8. **Independent verification of run 3** from a fresh clone: `cli-attestation verify --head 462a32a... --base 52065b1...` gave `status=valid`, `reason=ok`, `verdictClass=self-authored`. `inspect-v6`: three leaves, `harness: null`, `harnessTranscriptHash: null`. The real reviewers got the lowest class, because the SubagentStart markers were written to the parent repository's `.ai-sdlc/subagent-sessions/`, and the leaves were made in the worktree.
9. **Forgery test A** (no agent). On a branch with an unreviewed `Mul` function, the evaluator typed three "LGTM" transcript files and one verdict file, ran `cli-attestation emit-leaf --harness codex --model invented-model` for `code-reviewer-codex` and `test-reviewer-codex`, signed with `sign-v6` and the operator key, and verified: `status=valid`, `reason=ok`, `verdictClass=self-authored`.
10. **Forgery test B.** The same, plus three hand-written marker files `{"agentType":"ai-sdlc:<reviewer>","firedAt":"<now>"}` in `.ai-sdlc/subagent-sessions/`. `emit-leaf` printed `verdictClass=independent` three times; the verifier: `status=valid`, `reason=ok`, `verdictClass=independent`. No reviewer ran, and no second harness exists on the host.
11. **Orchestrator.** `cli-orchestrator status` and `tick --dry-run` ran. With `AI_SDLC_AUTONOMOUS_ORCHESTRATOR=experimental` and `--spawner mock`:
    - The first try refused the task: "Already-in-flight check: failed (live subprocess PID 29426)". PID 29426 was another evaluator's unrelated Claude Code session whose command line contained the text "toy-1". The filter scans `ps -ax` for the task ID (`rollback.ts#L128`).
    - With the task renamed `ZQX-7`, all ten filters passed ("Blast-radius overlap check: passed ... no blast-radius files declared — admitted (degrade-open)"), then: "recoverable abort for ZQX-7: worktree preserved ... (0 commits ahead, 0 checkpoints). Next tick will resume." A second tick did the same.
    - Events went to `artifacts/_orchestrator/events-2026-09-28.jsonl`, an untracked directory at the repository root.

Not run: a Codex reviewer (no Codex CLI or account), any GitHub workflow, the `ai-sdlc run` issue path and its SQLite store, the TUI, the `isolated` CI tier, a real (non-mock) tick dispatch.

## Where each piece of state lives

| State | Location | In the target's Git? |
| ----- | -------- | -------------------- |
| Rules and gate config | `.ai-sdlc/*.yaml`, `.github/workflows/*.yml` | yes |
| Task | `backlog/tasks/*.md`, moved to `backlog/completed/` in the PR | yes |
| Transcript hashes and attestation | `.ai-sdlc/transcript-leaves/<patch-id>.jsonl`, `.ai-sdlc/attestations/<patch-id>.v6.dsse.json` | yes |
| Transcripts and verdicts | `.ai-sdlc/transcripts/`, `.ai-sdlc/verdicts/` (left untracked in run 3); source transcripts in `~/.claude/projects/<slug>/.../subagents/` | no |
| Signing key | `~/.ai-sdlc/signing-key.pem` on the agent host | no |
| Independence markers | `<root>/.ai-sdlc/subagent-sessions/*.json` | no (untracked) |
| Worktrees, sentinel | `.worktrees/<id>/`, `.active-task` | no |
| Orchestrator events | `artifacts/_orchestrator/events-<date>.jsonl` | no (untracked) |
| Stuck counter | process memory | no |
| Incremental-review marker | a PR comment | no (GitHub) |
| Cost, audit, handoff, autonomy ledgers | `.ai-sdlc/state.db` (SQLite) | no (gitignored) |
| Decisions | `.ai-sdlc/_decisions/events.jsonl` | yes |
| Dispatch board | `.ai-sdlc/dispatch/{queue,inflight,done,sessions}/` | depends on the adopter |

## The JSON schemas under `spec/schemas/` as LAYUP record schemas

LAYUP is Go standard library only, so it cannot validate JSON Schema 2020-12 with a library. "Use" means: copy a field set into a Go struct and a test fixture, with the Apache-2.0 notice.

| Schema | Fit | Why |
| ------ | --- | --- |
| `decision.v1` | good pattern | Append-only event log in the repository; a decision is a projection. Matches I1 and the question and decision records. |
| `dispatch-verdict.v1` | good pattern | `taskId`, `outcome`, `durationMs`, `iterationsAttempted`, `acceptanceCriteriaMet`, `verifications`, reviewer start and end. LAYUP adds token counts (S10). |
| `orchestrator-events.v1` | good pattern | One JSONL line per event with `ts` and `type`; includes stuck, quarantine and resume events. |
| `attestation-envelope-v6` | pattern | A Merkle root over reviewer transcript hashes, signed once. Useful only when a key the agent cannot reach signs it. |
| `refinement-verdict.v1` | pattern | Readiness verdict with gates and a `questions` list; close to an S1 gap-check output. |
| `agent-role` (`handoffs`) | pattern | A handoff has `target`, `trigger` and a `contract` with a schema and required fields (S5, Author-8). |
| `quality-gate` | pattern | Graduated enforcement and a recorded override (role and justification). |
| `backlog-task.v1`, `capture-record.v1` | minor | Task frontmatter; an issue capture with an audit trail. |
| design system, soul, journey, variant, embedding, vector store, database pool, worktree pool, subscription plan, substrate contract, signal ingestion | no fit | Product features outside the PSB. |

## In-Scope items S1–S12

| Item | Mark | Evidence |
| ---- | ---- | -------- |
| S1 Problem Statement Quality | partly | The DoR gate checks each task and can list questions (`refinement-verdict.v1`); it was a tick filter in run 11. It does not read a problem statement or batch the questions. |
| S2 Reproducible Discipline Setup | partly | `ai-sdlc init` is repeatable and `doctor` is honest (runs 2, 3). No baseline pin, no evidence per value, and the gate is Node-only on a Go target. |
| S3 Rule Protection | partly | Hook denies `Write`/`Edit` to rule paths; `Bash` writes pass (run 4). Branch protection has `enforce_admins: false`. |
| S4 Stack-Dependent Gates | no | The QualityGate checks description and acceptance criteria; the CI template is `pnpm` only. No stack selection. |
| S5 Role Handoffs | partly | Developer and reviewer return JSON envelopes that the pipeline parses; schemas exist for verdicts. The handoff contracts are not validated on each transition. |
| S6 Autonomous Clarification | no | Questions go to the Operator (`cli-decisions escalate`). |
| S7 Verification on Every Change | partly | Three LLM reviewers and an attestation on each change (run 7). No deterministic layout, interface or test-pyramid gate. |
| S8 Human-on-the-Loop | partly | Agents never merge; the merge helper needs green, clean and trusted checks. No escalation rule for business-forking decisions; the Decision Catalog is Draft. |
| S9 Stall Resolution | partly | 2-iteration cap, then `[needs-human-attention]`; recoverable abort and resume; quarantine branches. No fresh-context diagnosis; run 2 stopped and asked the Operator. |
| S10 Cost Visibility | partly | `cost_ledger` has tokens and USD, but in gitignored SQLite. Nothing about cost was committed in run 3. |
| S11 Specification Synthesis | no | README: it covers "contract → shipped", not "idea → contract". |
| S12 Harness-Agent Neutrality | partly | Codex and Copilot spawners and Codex reviewer variants exist. The pipeline and its hooks are Claude Code artifacts, and the second-harness proof is forgeable (runs 9, 10). |

## Invariants I1–I9

| Invariant | Effect | Reason |
| --------- | ------ | ------ |
| I1 Git is the record | conflicts | Cost, audit and handoff ledgers are in gitignored SQLite; events are untracked; a PR comment holds review state; the key and transcripts are on the host. |
| I2 Independent repository | conflicts | The pre-push hook and verifier need the plugin install (`$CLAUDE_PLUGIN_ROOT` or `~/.claude/plugins/cache`); the gate workflow needs `pnpm`. This is against O-76. |
| I3 Agents cannot change rules | conflicts | Bash bypasses the hook (run 4); the key the agent host holds signs any envelope (runs 9, 10); admins bypass branch protection. |
| I4 No value without evidence | conflicts | Defaults with no evidence: `PT30M`, `maxRetries: 2`, gate thresholds of 1, autonomy and calibration numbers. |
| I5 Inactive check is not a pass | conflicts (default config) | The engine fails a missing metric, and `doctor` says "audit-only". But the scaffolded attestation workflow never fails, blast radius is "degrade-open", and a missing second harness gives a note, not a stop. |
| I6 Deterministic first | neutral | The attestation and hooks are deterministic; the review itself is LLM judgement. |
| I7 Stack adds, never removes | neutral | No stack concept. |
| I8 Pinned baseline | neutral | Runtime packages use a range (`>=0.26.1 <1.0.0`), not a pin. No Armature. |
| I9 Replaceable harness | conflicts | The main flow is a Claude Code plugin (slash command, subagents, hooks); the Codex path is attended or a bridge script. |

## Deep-check findings it answers

- **Sol-2 (a PR gate cannot prevent human review).** Pattern: the PR stays a draft until the reviewers and the attestation finish, then `gh pr ready` runs CI once (run 7). A person can still open a draft.
- **Sol-27 (no enforceable PR step).** Pattern: Steps 11 to 13 give branch, draft PR with verdicts and a `References`/`Closes` line, marker comment, ready flip; agents never merge, and one helper merges only on green, clean, trusted checks.
- **Sol-5 and Fable-M3 (commit author is not proof).** Partly: a DSSE signature over the patch-id binds content to a key, and the verifier loads nothing from the checkout. In the default flow the key is on the agent host, so it proves nothing about who acted (runs 9, 10). The `isolated` tier, with a key only a protected CI environment holds, is the part that answers these findings.
- **Fable-M1 (one credential for agents and orchestrator).** Pattern only: `isolated-review.yml` splits an unsigned review job from a signer job in a protected environment.
- **S12 cross-harness verification (Fable-M23 area).** Not answered. The dispatch-time harness filter is useful; the verification side is a string claim. The plugin also adds rules through host-level SessionStart and SubagentStart hooks, which is the second rule source that Fable-M23 warns about.
- **Author-1 (disagreement stall).** Pattern: a hard cap of 2 developer iterations on review failure, then a flagged PR.
- **Author-2 (mid-session questions).** Pattern: `cli-decisions escalate --exit-code 1` writes a decision, the session exits, the answer resumes the task.
- **Author-8 (handoff schema).** Pattern: `agent-role` handoff contracts and `dispatch-verdict.v1`.
- **Fable-M19 (a wrong gate with no exit).** Partly: `soft-mandatory` gates allow an override by a named role with a justification.
- **Fable-M9.** Confirms the risk: the scaffolded required check fails on a Go target (run 2).
- None for Sol-1, Fable-M2, Sol-26 (AI-SDLC puts its gates in its own tooling, the opposite of O-76), or Sol-30.

## Parts and their verdicts

| Part | Verdict | How LAYUP would take it | Reason |
| ---- | ------- | ----------------------- | ------ |
| Cross-harness review check (`harness` field, `verdictClass`) | reject | — | A string and a marker file that the agent side writes; forged to `independent` in run 10. |
| `isolated` tier: signer in a protected CI job with a CI-only key | borrow the pattern | LAYUP's verification record is signed or committed by a credential the working agents never hold (Fable-M1, Sol-5). | The only part that resists the forgery; not tested here. |
| Verifier that never loads code from the checkout it checks, and reads keys from the base branch | borrow the pattern | `layup` checks a target with its own binary and the base branch's rules, never the PR's. | Sound trust boundary (`verify-attestation.mjs#L37`). |
| DSSE and Merkle envelope over transcript hashes | borrow the pattern | Optional evidence bundle for a review record. | Useful binding; needs an agent-proof key. |
| PreToolUse rule-path hook | reject | — | Claude-only and bypassed by Bash (run 4). |
| Draft PR, then `ready` after review and attestation | borrow the pattern | Implement phase: open a draft, verify, then request human review (Sol-2, Sol-27). | Seen in run 7. |
| Iteration cap and `[needs-human-attention]` | borrow the pattern | Parameter N for review rounds (O-82), then a stall record. | Simple and visible. |
| Checkpoint commits, recoverable abort, quarantine branch | borrow the pattern | Git-only resume state for a killed session. | I1-compatible; the resume on a real session was not tested. |
| Admission filter chain | borrow the pattern, corrected | Named filters with a trace per task; no host-wide process scan; no "degrade-open". | Run 11 showed a false positive and a degrade-open pass. |
| `cli-decisions escalate` and the decision event log | borrow the pattern | Headless question: write a record, exit, resume on answer. | Matches issue-centric, Git-first design. |
| Schemas (`decision`, `dispatch-verdict`, `orchestrator-events`, `agent-role` handoffs) | borrow the pattern | Field sets for LAYUP records, in Go structs. | See the schema table. |
| `ai-sdlc init` scaffold | reject | — | Node-only gate on a Go target; tries branch protection without the flag. |
| npm packages and plugin as a runtime dependency | reject | — | Node and SQLite; conflicts with I1, I2, O-76 and F-0004 fact 1. |

## Where the searchers were wrong or incomplete

- Searcher A: "verifier checks reviewer harness differs". Only in the v5 path, and only when the implementer string is `codex`. The v6 verifier logs the harness and checks nothing. The runbook itself says "The harness field is forensic/audit only".
- The README says "Reviewer collusion is mechanically impossible" and the runbook "verified cryptographically, not just by convention". Runs 9 and 10 show a valid envelope with no review, and the `independent` class from two hand-written files.
- Searcher A, S3 "Attestation is required on main": true for the AI-SDLC monorepo. The scaffold for an adopter is audit-only, and `doctor` says so.
- Searcher A, S4 "QualityGate resources": they exist, but the scaffolded gates check PR description and acceptance criteria, not code. No stack-dependent gate.
- Searcher A, I1 "P": state is split; the cost and audit ledgers are in a gitignored SQLite file. That is a conflict.
- Searcher A, S1 "no": the per-task DoR gate with a questions list is a partial answer.
- Searcher A, "Deployment database ... not stated": `better-sqlite3` and `.ai-sdlc/state.db` are the database.
- Summary: "repo-only decision storage not stated". The decision log is in the repository (`.ai-sdlc/_decisions/events.jsonl`); the other ledgers are not.

## Limits of this evaluation

- Codex was not installed, so no run had a second harness. The Codex reviewer path was read, not run.
- No GitHub workflow ran; `gh` was a logging stub. The PR steps were observed as `gh` calls only.
- The Claude Code session needed the host's login, so the harness process used the real `HOME`; every tool process used the isolated `HOME`. User-level `CLAUDE.md` was still loaded (Fable-M23). Run 3 left session transcripts in `~/.claude/projects/-private-tmp-...-eval-work-ai-sdlc-toy/`, which the persist step needed.
- The orchestrator ran only with the mock spawner. A real tick, a real checkpoint resume and a quarantine were read in code, not observed.
- The `ai-sdlc run` issue path, the TUI, the DoR evaluator and the SQLite store were not run.
