# T-hbw8 — second public-solution search: the sources

The method of [`../root-cause-missed-solution.md`](../root-cause-missed-solution.md):
queries in the problem's own terms, every curated list read in full, two blind
searchers on different models. Built on 2026-09-28 by the author; the two
searchers got the problem statement, the idea owner's answers, the vision brief
and these sources only — not the architecture, the target design, the deep
check or the name of any candidate.

## Curated lists (GitHub search `awesome agents in:name stars:>500`, sorted by stars)

Kept, and given to the searchers in full (the text is not stored here because of
its size, about 7,800 lines; the commit SHA makes it reproducible):

| List | Commit | Lines | Repository links | Why kept |
| ---- | ------ | ----- | ---------------- | -------- |
| `andyrewlee/awesome-agent-orchestrators` | `a1278087adb6` | 307 | 239 | orchestrators of agents |
| `e2b-dev/awesome-ai-agents` | `9596ab1e69fb` | 5,591 | 127 | general agents, 30,200 stars |
| `bradAGI/awesome-cli-coding-agents` | `bfa968cb5d8a` | 884 | 383 | coding agents and the harnesses that orchestrate them |
| `caramaschiHG/awesome-ai-agents-2026` | `781b695139b8` | 672 | 128 | agents, frameworks and tools of 2026 |
| `kyrolabs/awesome-agents` | `70fd6c84814c` | 224 | 175 | general agents |
| `vijaythecoder/awesome-claude-agents` | `2050f3c60fcf` | 170 | 2 | an agent development team |

Rejected: `fr0gger/Awesome-GPT-Agents` (cybersecurity), `mergisi/awesome-openclaw-agents` (templates for one product), `ANative-Lab/Awesome-Self-Evolving-Agents`, `FoundationAgents/awesome-foundation-agents` and `ysymyth/awesome-language-agents` (research papers), `slavakurilyak/awesome-ai-agents` and `jim-schwoebel/awesome_ai_agents` (general resources, covered by the kept general lists), `kaushikb11/awesome-llm-agents` (agent frameworks, not orchestration of delivery), `steel-dev/awesome-web-agents` (web browsing agents), `AgenticHealthAI/Awesome-AI-Agents-for-Healthcare` (a domain).

## GitHub repository searches ([`gh-search.tsv`](gh-search.tsv), sorted by stars, top 20 or 25 each)

| Source in the problem | Query | Results |
| --------------------- | ----- | ------- |
| PSB §2 (the roles) | multi-agent software development team roles | 10 |
| PSB §2 (the roles) | product owner architect developer QA agents | 0 |
| PSB Problem 1 | human in the loop agents escalation clarifying questions | 0 |
| PSB §6 Human Decision Points | agent orchestration approval budget governance | 1 |
| PSB Problem 5 | agent token cost tracking budget | 20 |
| PSB Problem 4 | agent stall deadlock retry loop detection | 0 |
| PSB Problem 6 | requirements specification from problem statement agents | 1 |
| PSB Problem 7 | multiple coding agents shared rules AGENTS.md | 0 |
| PSB §1 | autonomous software delivery agents | 20 |
| PSB §1 | coding agent orchestrator | 20 |
| PSB In Scope, Role Handoffs | agent handoff structured state | 14 |
| PSB In Scope, Verification | cross agent verification code review different model | 1 |
| GitHub topics | `agent-orchestration`, `multi-agent-systems`, `ai-software-engineer`, `agentic-coding`, `coding-agents`, `multi-agent` | 25 each |

234 rows, 212 distinct repositories. GitHub search needs every word to match,
so several long queries found nothing; the curated lists and the searchers' own
web searches cover that. A fact of this search: the repository search did not
find the candidate that the first search missed; only a curated list held it.
