# T-ywk7 — the public-solution search of M2b: the result

Two blind searchers screened the same sources ([`sources.md`](sources.md)) and ran
their own web searches in the words of the problem:
[`searcher-A.md`](searcher-A.md), GPT-6 Sol on Devin, and
[`searcher-B.md`](searcher-B.md), Claude Opus 5.5 on Claude Code (O-178 (b)).
Their brief ([`brief.md`](brief.md)) held the facts, the milestone row and the
sources: no row of the survey of `T-w73g`, no name of a candidate, nothing of the
architecture. This file is the author's reading of the two files, written after
both ended. The author read the survey's cited files while the searchers ran;
what keeps the survey out of the search is the brief, which holds none of it,
and the commit of this search before the commit of the survey table (note 6 of
the plan review).

**Both conclude: no public solution can be the base of `M2b`.** No candidate can
be a dependency of a controller in Go with the standard library only
(`F-0004#1`): each is a whole program to run beside it, or a pattern, a format or
a data source. A: "No candidate above is approved as a Go
module dependency". B: "Every candidate is a whole program or a pattern; none can
be vendored into a stdlib-only controller."

## What the search gives each part of M2b

| Part | Candidates (searcher) | What they bear on in `M2b` (the decisions come in S3) |
| ---- | --------------------- | ----------------------------------------------------- |
| The start, the limit and the stop of a session | Bernstein, scion, sortie, NEEDLE, Superharness (B, A) | No harness has a wall-clock limit of its own (B: "`codex exec` and Gemini CLI have no time flag; Claude Code has `--max-turns` only"), so the controller owns the timer and stops the process group, as the architecture's `wall` does |
| The probe | ACP `initialize` (A, B); NEEDLE `test-agent`, the probe spec of mco, gate4agent, the `system/init` event of Claude Code (B) | A version check with no model call (NEEDLE, mco), and a version in the harness's first event (Claude Code). No candidate runs a small real task on each harness and commits a dated probe record (B). ACP is a handshake over JSON-RPC; ADR-0015 starts a harness by its command line and reads a result file, so a protocol would change that ADR |
| The choice of a harness and a model | ordewell, Galley, sortie, Coder Eval, scion (B); Galley, Orchestrate (A) | None chooses from recorded cost and outcome (B, "What no candidate covers") |
| The usage and the money | Claude Code headless, Codex `exec --json`, Gemini CLI (A, B); OpenCode `run` (A; B's query 7 only); Headless CLI (A); MartinLoop, agentacct, LiteLLM's price file, models.dev (B) | Each harness reports usage in its own form, and only Claude Code reports money, as a client estimate (B); a value that a harness does not report stays missing, never 0 (Headless CLI, agentacct); LiteLLM's file (MIT) and models.dev (MIT) publish prices and context sizes per model |
| The record in Git | aGiTrack (A, B), shift-log, copilot-session-usage (B) | They put tokens into commits or notes, with no cost and no latency (B) |
| The rules in one neutral form | AGENTS.md (A, B); agentsync (A); Rulesync, AIWG, agents-md-cookbook (B) | `AGENTS.md` is read by Codex, OpenCode and Gemini CLI, and by Claude Code natively from v2.1.277 when no `CLAUDE.md` exists (B); agents-md-cookbook and AIWG keep tables of each harness's rule-file names |
| No agent change to a rule path | Chock (A, B); gh-aw, CC Safety Net, loopgate, sandbox-runtime, Codex's sandbox, GitHub push rulesets that restrict file paths (B) | Each blocks or unstages a write to a protected path, and none counts the writes or ships known-bad commits (B). An OS sandbox that denies writes, and a push ruleset that restricts file paths (private repositories on Team or Enterprise plans), bear on the rule protection of `M2f` too |
| The handoff | Agent Context Store (A, B); AI-SDLC (A); JSON Schema, AgentPlane, A2A (B) | Go's standard library has no JSON Schema validator (B); none validates every role transition (B) |

## What this task observed while it ran the search

- **The usage of a Claude Code session is in `modelUsage`, not in `usage`.** In
  the plan review (no subagent, no web tool) the two are equal. In searcher B's
  run, `usage` covers only the main loop (53,309 output tokens), while
  `modelUsage` covers each model (433,333 output tokens of `claude-opus-5-5` and
  478,125 of `claude-haiku-5-5`), and its costs add up to `total_cost_usd`. Both
  runs are Claude Code 2.1.295.
- **A harness runs models of its own.** Claude Code's `WebFetch` reads a page with
  `claude-haiku-5-5`, a model of the list not to use, unless
  `ANTHROPIC_DEFAULT_HAIKU_MODEL` names another; a test run showed both cases.
- **A credential can fail for every model.** OpenCode returned 401 for each model,
  its free one included, and Devin's weekly quota stopped two sessions.

## The limits of the search

- Searcher A rejects many list entries as "primary not checked": a rejection by
  scope, not a checked fact. Searcher B rejects list entries by the list text
  alone. A candidate that neither opened can still hold a pattern.
- Searcher B's checks on GitHub's API got 403; its dates come from release feeds or
  a mirror. It knows OpenCode's `step_finish` tokens and cost only from
  third-party pages.
- One searcher is the author's model (O-178 (b)), blind to the author's context.
- Searcher B's file is 703 lines, over the brief's 300 to 450.
