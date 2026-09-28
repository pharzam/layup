# T-hbw8 — second search: the result

Two blind searchers ([`searcher-A.md`](searcher-A.md), GPT-6 Sol; [`searcher-B.md`](searcher-B.md), Grok 4.7) screened the same sources ([`sources.md`](sources.md)) and ran their own web searches in the problem's words. Both sessions stopped at Devin's daily quota after they had written complete files.

**Both conclude, independently: no public solution can be LAYUP's base.** Each candidate either keeps the project state outside the project's Git (Invariants 1 and 2) or has none of: the pinned baseline with evidence for each value, rule protection, stack gates that only add rules, the PSB's Human Decision Points and escalation rule, the stall procedure with a fresh diagnosis, per-task telemetry in the repository, and requirements traced to the approved problem statement. Searcher A, "What no verified candidate covers end-to-end"; searcher B, "What no candidate covers".

**Both found Paperclip** (A: candidate H; B: 3.1) and both rate it a pattern or component, not a base: "The server is the system of record" (B); "server-held project state conflicts with Git-only record unless separately exported" (A).

## The candidates worth reuse (union of the two shortlists)

| Candidate | License, language, stars (2026-09-28) | What it gives LAYUP | Why not the base |
| --------- | ------------------------------------- | ------------------- | ---------------- |
| [Paperclip](https://github.com/paperclipai/paperclip) | MIT, TypeScript, about 91,400 | budgets that stop agents, an org chart of roles, approvals, heartbeats, adapters for Claude Code and Codex | a server is the record (Inv. 1, 2); "manage business goals, not pull requests" |
| [Gas Town](https://github.com/gastownhall/gastown) with Beads | MIT, **Go**, about 18,200 | Git-backed work tracking, stuck-agent recovery, a merge queue (B 3.5) | work state in the Beads ledger and a Dolt database (Inv. 1, 2); no baseline pin |
| [Spec Kitty](https://github.com/spec-kitty/spec-kitty) | MIT, Python, about 1,650 | spec, plan, tasks and review over Git worktrees (B 3.4) | no baseline, cost record, stall diagnosis or rule protection |
| [BMAD-METHOD](https://github.com/bmad-code-org/BMAD-METHOD) | license not stated by the API, Python, about 53,600 | specification and role skills (B 3.2) | no discipline baseline |
| [GNAP](https://github.com/farol-team/gnap) | MIT, about 86 | a Git-only task, run and cost protocol: "No servers. No databases." | small; inactive since 2026-03-17; no gates |
| [AI-SDLC](https://github.com/ai-sdlc-framework/ai-sdlc) | Apache-2.0, TypeScript, about 375 | cross-harness review and attestation (A) | repo-only decision storage not stated |
| AgentJury, Loki Mode, isitdone and the second-model review tools (A, B) | various; Loki Mode is BUSL-1.1 | a blind review council, a mechanical done-check, a second-model review | components only |
| [MetaGPT](https://github.com/FoundationAgents/MetaGPT) | — | role SOPs and document output | stale; no invariant support stated |

## What this changes

- The decision to build LAYUP's own orchestrator now has evidence: no base exists that meets Invariants 1 and 2 with the PSB discipline.
- The rewrite of the architecture reuses named parts where they fit, after a hands-on check: the Go candidates first (Gas Town and Beads), then the Git-only protocol (GNAP), the spec workflow (Spec Kitty), the cross-harness review (AI-SDLC), and Paperclip's budget and approval model as a pattern.
