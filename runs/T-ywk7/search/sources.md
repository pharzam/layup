# T-ywk7 — the public-solution search of M2b: the sources

The method of [`docs/guardrails.md`](../../../docs/guardrails.md) §2, "A
public-solution search that confirms the design": the search runs before any
decision that fixes the shape of `M2b`; its queries come from the facts that
the requirements of `M2b` cite in `PRD-0001` §12, and from the `M2b` row of the
plan; each curated list is read in full; two searchers on different models,
blind to each other. The author built these sources on 2026-10-09; the two
searchers got [`brief.md`](brief.md), the five lists and
[`gh-search.tsv`](gh-search.tsv) only: not the architecture, the
specification, the survey of `T-w73g`, or the name of any candidate.

## The curated lists

Chosen by eight GitHub repository searches, on 2026-10-09, sorted by stars:
`awesome in:name coding agents`, `awesome in:name cli agents`,
`awesome in:name agent orchestrators`, `awesome in:name claude code`,
`awesome in:name codex`, `awesome in:name llm observability`,
`awesome in:name llmops`, `awesome in:name ai agents stars:>5000`. The rule:
keep the most-starred list of each domain that the problem touches, whose
entries are tools rather than prompts or skills. Each kept list is its
`README.md` at the commit below, given to the searchers in full; the text is not
stored here, and the commit makes it reproducible.

| List | Commit | Lines | Repository links | Why kept |
| ---- | ------ | ----- | ---------------- | -------- |
| `bradAGI/awesome-cli-coding-agents` | `1eb67af8afc7` | 884 | 383 | terminal coding agents and the harnesses that orchestrate them |
| `andyrewlee/awesome-agent-orchestrators` | `fb2eae87aac6` | 324 | 256 | orchestrators of agents |
| `hesreallyhim/awesome-claude-code` | `7eddf4aaebeb` | 770 | 205 | the tools around one harness: usage monitors, runners, hooks |
| `tensorchord/Awesome-LLMOps` | `98def2e0f194` | 801 | 420 | observability and the cost of model calls |
| `RoggeOhta/awesome-codex-cli` | `bd5d0fdfdeac` | 494 | 215 | the tools around a second harness |

Not kept, with the reason:

| List | Stars | Why not |
| ---- | ----- | ------- |
| `Agent-Analytics/awesome-multi-agent-orchestrators` | 155 | the domain of `andyrewlee/awesome-agent-orchestrators`, with fewer entries |
| `ContextJet-ai/awesome-llm-observability` | 40 | the domain of `tensorchord/Awesome-LLMOps`, with fewer entries |
| `VoltAgent/awesome-claude-code-subagents`, `composio-community/awesome-codex-skills`, `VoltAgent/awesome-codex-subagents` | 25,597; 16,816; 6,293 | prompts, subagents and skills, not tools |
| `e2b-dev/awesome-ai-agents` | 30,312 | agents in general, not the running of a harness or its cost; read in full by the search of `T-hbw8` |
| `iSEngLab/Awesome-Self-Evolving-Coding-Agents` and the lists of fewer than 100 stars | — | research papers, or a domain that a kept list covers |

## The GitHub search

[`gh-search.tsv`](gh-search.tsv): 20 repository searches, ten results each by
GitHub's best match, on 2026-10-09; each row has the source of its query (a fact
and its words, or the words of the `M2b` row), the query, the stars, the
repository, the date of its last push and its description. A first run of 18
longer queries found results for only 6 of them (GitHub matches every word of
a query), so the queries are two or three words each.

## The two searchers

| Searcher | Model | Harness | Its file |
| -------- | ----- | ------- | -------- |
| A | GPT-6 Sol (`gpt-6-sol-high`) | Devin 3000.11.3, `devin -p --permission-mode dangerous`, with a `HOME` of its own that holds only Devin's credential and configuration files | [`searcher-A.md`](searcher-A.md) |
| B | Claude Opus 5.5 (`claude-opus-5-5`) | Claude Code 2.1.295, `claude -p --permission-mode dontAsk`, read-only tools and web search, writes only under `out/` | [`searcher-B.md`](searcher-B.md) |

The plan put searcher B on OpenCode (`opencode-go/grok-4.7`). At 06:08 UTC
OpenCode exited 1 after 9 s with `provider.auth`, `Upstream request failed:
Invalid credential` (401); a retry, and the model `opencode-go/grok-4.6`, gave the
same at 06:09. The same model on Devin (`grok-4-7-high`, its own `HOME`) read the
sources and stopped at 06:12 with "Your weekly usage quota has been exhausted",
before it wrote its file; searcher A stopped at the same quota at 06:13, after it
had written its file. By the Operator's O-178 (b) (#147), searcher B is Claude
Opus 5.5 in a fresh session with the same brief and sources: two models on two
harnesses, one of them the author's model, blind to the author's context.
