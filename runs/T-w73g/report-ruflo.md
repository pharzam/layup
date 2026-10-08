# Report: ruvnet/ruflo

The report of one read-only agent of task `T-w73g` (#133), word for word below
the line. Model: Claude Opus 5.5, an Explore subagent on Claude Code, 2026-10-07,
29 tool calls, 139.2 s. Repository: <https://github.com/ruvnet/ruflo> at
`9fa701a466159242dfec06379ce1278c849a3e34`, a shallow clone in a scratch
directory outside the checkout. No code from the repository was run.

---

I read the code in the parts you named. **License: MIT** (root `LICENSE`; the packages are `@claude-flow/*` on npm).

**What it is:** a TypeScript/Node harness that sits on top of Claude Code and Codex. It ships an MCP server, Claude Code hooks, about 47 plugins and swarm "consensus". Its state lives in a SQLite/HNSW vector store ("AgentDB") under `.swarm/` and `.claude-flow/`, not in Git. A few small Rust crates handle federation.

Paths below are relative to the clone root, `/private/tmp/claude-501/-Users-farzam-projects-layup/6f4ece82-2504-4e0e-b120-e779fedc77da/scratchpad/repos/ruflo/`.

| Item | Path | LAYUP concern | Verdict | Why |
|---|---|---|---|---|
| Token-usage ledger (Claude + Codex transcripts) | `plugins/ruflo-cost-tracker/scripts/_ledger.mjs` | Per-session cost recording | **Copy pattern** (port to Go) | Deterministic and local. It removes Claude duplicates by `message.id+requestId`. Codex reports running totals, so it takes the difference between totals and treats a drop as a reset. It splits cache reads from 5m/1h cache writes. ~150 lines. |
| Price table | `plugins/ruflo-cost-tracker/data/prices.json` | Milestone budget cap | Copy pattern | One editable file of per-model prices with a source URL for each. Fits a pinned file in the target repo. |
| Budget check | `.../scripts/budget.mjs`, `v3/docs/adr/ADR-164.1-budget-tracker-atomicity.md` | Budget cap | Skip the code, keep the ADR's lesson | The code only warns, never blocks, and calls `npx` into SQLite each time. The ADR's race analysis (check, then spend, then record) is worth reading: reserve before spending. |
| Headless worker spawner | `v3/@claude-flow/codex/src/dual-mode/orchestrator.ts:204-290` | Harness adapters | Copy pattern | Exact flags for `claude -p … --output-format text --max-turns --model` and `codex exec --sandbox read-only\|workspace-write -m`. One worktree per writer, timeout, output cap. Gotcha: `codex exec` hangs unless stdin is closed (`proc.stdin.end()`). |
| Thompson bandit router | `v3/@claude-flow/cli/src/ruvector/model-router.ts` | Routing learned at retrospective | Copy idea, make it deterministic | It keeps Beta(α,β) success/failure counts per complexity bucket and model, with an optional decay factor. For LAYUP, compute these from TSV records at retro and pick by posterior mean. Avoid its random sampling and keyword-based complexity score. |
| ContinueGate | `v3/@claude-flow/guidance/src/continue-gate.ts` | Stall detection | Copy pattern | Pure numeric checks: step limit, ratio of redone steps, cost acceleration (slope over the last 10 steps), forced checkpoints. Decisions are continue, checkpoint, pause or stop, and "pause" means a human looks. Drop its coherence and uncertainty inputs, which are LLM-derived. |
| Enforcement gates | `v3/@claude-flow/guidance/src/gates.ts` | Quality gate / hooks | Copy pattern | Regex checks for destructive commands and secrets, a tool allowlist, and a diff-size limit. All deterministic. |
| Fix-marker check and ratchet | `verification/README.md`, `witness-fixes.json`, `mcp-tool-baseline.json` | Quality gate / CI | Copy pattern | Each fix is pinned as a file plus a substring that must stay present, with history in JSONL. Bad-pattern counts are only allowed to go down. Both are cheap Git-friendly checks. Skip the Ed25519 signing. |
| Event-sourced claims | `v3/@claude-flow/claims/src/` | Single writer / records | Skip | Memory-only event store plus work-stealing. Git already does this job for LAYUP. |
| Raft/Byzantine/Gossip | `v3/@claude-flow/swarm/src/consensus/` | (review rounds?) | Skip | Distributed-systems code applied to agents in one process. The default path changes a fake peer object in memory, and a Byzantine node auto-approves its own proposal (`byzantine.ts:378`). Nothing here is independent review. |
| AgentDB / HNSW / ReasoningBank memory | `v3/@claude-flow/memory/`, `hooks/src/reasoningbank/` | — | Skip | State outside Git; conflicts with ADR on Git as system of record. |
| Plugin hooks | `plugin/hooks/hooks.json` | Hooks | Skip | Almost every hook ends in `\|\| true`, so failures are ignored. `PermissionRequest` auto-approves every `mcp__claude-flow__*` tool. That is the opposite of "agents cannot change the gates". |
| Provider manager | `v3/@claude-flow/providers/` | Harness routing | Skip | Calls model APIs directly with round-robin and fallback. LAYUP routes CLI harnesses, not APIs. |

**README claims the code does not back:**
- **"Intelligent routing (89% accuracy)"**: the prompt hook calls `router.cjs`, which describes itself as "*not* a learned model… a heuristic table". ADR-149 admits the seed scores are "hand-coded, not measured" and that "100% accuracy … is a property of the synthetic corpus".
- **"Learns from every task"**: the bandit only updates when an agent reports an outcome through the MCP tools (`hooks-tools.ts:5179,5265`).
- **"314 MCP tools"**: the repo's own baseline file lists 285.
- **"Queen-led Raft/Byzantine consensus"**: in-process unless a federation transport is wired in (`transport.ts` header says so).
- **"Set budgets"**: the code warns but never blocks.
- The "1.3×–1953×" benchmark is an external gist, not in the repo.

**Top 3 recommendations:**
1. **Port `_ledger.mjs` to Go** for per-session cost rows written to TSV. Pair it with a `prices.json`-style table in the target repo, and reserve budget before each session starts (the ADR-164.1 lesson) rather than warning after.
2. **Use ContinueGate's numeric checks for stall detection** (step limit, redone-work ratio, cost slope). When one trips, hand off to your fresh-context diagnosis or a human, and never ask an LLM whether to stop.
3. **Copy the Beta-count routing table** keyed by (role, complexity bucket, harness/model), rebuilt from records at each retrospective with a deterministic choice. Use the orchestrator's `claude -p` / `codex exec` flags, worktree-per-writer rule and stdin-close fix as the starting point for your harness adapters.

Use the code as reference only, not as a dependency: it is tied to the npm/`npx` CLI, AgentDB and Claude Code plugin internals.
