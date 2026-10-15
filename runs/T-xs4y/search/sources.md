# T-xs4y — the public-solution search of M2e: the sources

The method of [`docs/guardrails.md`](../../../docs/guardrails.md) §2, "A
public-solution search that confirms the design", as `T-msrd` ran it: the search
runs before any decision that fixes the shape of `M2e`; its queries come from
the facts that the requirements of `M2e` (REQ-005, REQ-007, REQ-012, REQ-003,
NFR-004, NFR-001) cite in `PRD-0001` §12, and from the `M2e` row of the plan;
each curated list is read in full; two searchers on different models, blind to
each other. The author built the sources on 2026-10-15. The searchers got
[`brief.md`](brief.md), the five lists and [`gh-search.tsv`](gh-search.tsv) only:
not the architecture, the specification, or the name of any candidate.

## The repository search

[`gh-search.tsv`](gh-search.tsv): 24 GitHub repository searches (`gh search
repos`, best match, 10 results each), each with the fact words it comes from.
The first try used long phrases, which the search read as one quoted phrase and
answered with no result (the same as in `T-msrd`); the file holds the second try,
with two to four words per query. All 24 queries gave results.

## The curated lists

Chosen on 2026-10-15 from list searches (`awesome workflow engines`, `awesome ai
agents`, `awesome claude code`, `awesome ai agent orchestration`, `awesome
state machine`, `awesome code review automation`): the most-starred list of each
domain that the problem touches (workflow engines, coding agents and their
orchestration, GitHub automation). The two lists of `T-msrd` that serve the
same domains (`awesome-cli-coding-agents`, `awesome-actions`) are kept at their
current commits. The text is not stored here; the commit makes it reproducible.

| List | Commit | Lines | Links |
| ---- | ------ | ----- | ----- |
| `meirwah/awesome-workflow-engines` | `2d0be88a651d` | 134 | 319 |
| `hesreallyhim/awesome-claude-code` | `30ec3281cb86` | 848 | 524 |
| `bradAGI/awesome-cli-coding-agents` | `0fbf2821f03f` | 884 | 416 |
| `sdras/awesome-actions` | `0edccf33480b` | 572 | 434 |
| `e2b-dev/awesome-ai-agents` | `9596ab1e69fb` | 5591 | 1175 |

## The searchers

| Searcher | Model | Harness and mode | Output |
| -------- | ----- | ---------------- | ------ |
| A | GPT-6 Sol (`gpt-6-sol`, the reasoning tier) | Devin 3000.11.92, `devin --model gpt-6-sol -p --permission-mode dangerous`, 11:51:08 to 11:54:03Z (2 min 55 s) | [`searcher-A.md`](searcher-A.md) |
| B | Claude Sonnet 5.5 (`claude-sonnet-5-5`) | Claude Code, `claude -p --permission-mode dontAsk`, web search and fetch, 11:26 to 11:30:42Z | [`searcher-B.md`](searcher-B.md) |

Two tries of A failed first and are recorded: a start with a workspace-trust
refusal (`Refusing to run in an untrusted workspace`; fixed with
`--respect-workspace-trust false`), and a run on Devin's default model that sat
idle for 13 minutes with no output (2 s of CPU; killed). B's write of `out/search.md`
was denied by its allow rule (an absolute path against `Write(out/*)`); the
author saved the content from B's denied call, unchanged, as `searcher-B.md`. B
says it briefed sub-agents to read the five lists (its Limits).
