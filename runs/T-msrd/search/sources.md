# T-msrd — the public-solution search of M2d: the sources

The method of [`docs/guardrails.md`](../../../docs/guardrails.md) §2, "A
public-solution search that confirms the design", as `T-ywk7` and `T-g49f` ran
it: the search runs before any decision that fixes the shape of `M2d`; its
queries come from the facts that the requirements of `M2d` (REQ-002, REQ-003,
NFR-001, NFR-002, NFR-003, NFR-006) cite in `PRD-0001` §12, and from the `M2d`
row of the plan; each curated list is read in full; two searchers on different
models, blind to each other. The author built these sources on 2026-10-11; the
two searchers got [`brief.md`](brief.md), the five lists and
[`gh-search.tsv`](gh-search.tsv) only: not the architecture, the specification,
or the name of any candidate.

## The repository search

[`gh-search.tsv`](gh-search.tsv): GitHub repository searches (the API's best
match, 10 results each where it gave 10), on 2026-10-11 from 14:36Z, each query
with the fact words it comes from. The first 20 queries were long; eight of them
gave no result:

- `repository setup from template with evidence`
- `project scaffolding reproducible configuration`
- `repository template configuration provenance`
- `github rulesets verify effective rules`
- `ai agent cannot modify ci rules`
- `coding agent guardrails branch protection`
- `orphan branch metadata records`
- `ci passes without the generator`

So eight shorter queries in the same fact words were run after them, each marked
"short form" in its `from` column. The file holds 20 query strings with a
result and 150 results. The trailing spaces of seven descriptions are trimmed here for the
repository's whitespace check; the searchers read them untrimmed, and nothing
else differs.

## The curated lists

Chosen on 2026-10-11 from list searches (`awesome policy as code`, `awesome
github security`, `awesome scaffolding`, `awesome project templates`, `awesome
ai coding agents`, `awesome software supply chain security`, `awesome github
apps`, `awesome gitops`, `awesome cookiecutter`, `awesome boilerplate generator`,
`awesome project generator`, `awesome repository template`). The rule: keep the
most-starred list of each domain that the problem touches (policy as code, the
security of a repository's settings and branches, coding agents and their
guardrails, Git as the system of record, the automation of GitHub
repositories), whose entries are tools or methods. No list of project
scaffolding was kept: the searches gave README templates and single generators,
none a list of tools (the best, `trevtrich/awesome-scaffolding`, had 3 stars).
Each kept list is its file at the commit below, given to the searchers in full;
the text is not stored here, and the commit makes it reproducible.

| List | Commit | File | Lines | Links | Why kept |
| ---- | ------ | ---- | ----- | ----- | -------- |
| `hysnsec/awesome-policy-as-code` | `af4e6354067b` | `README.md` | 178 | 63 | the most-starred list of policy as code (211 stars) |
| `bureado/awesome-software-supply-chain-security` | `31252e6810c4` | `README.md` | 847 | 829 | the most-starred list of supply-chain security, which holds the tools that check a repository's branch rules and settings (380 stars) |
| `bradAGI/awesome-cli-coding-agents` | `1eb67af8afc7` | `README.md` | 884 | 407 | the most-starred list of command-line coding agents and their guardrails (1,339 stars) |
| `weaveworks/awesome-gitops` | `82af5b15f445` | `README.md` | 98 | 52 | the most-starred list of GitOps, Git as the system of record (1,731 stars) |
| `sdras/awesome-actions` | `0edccf33480b` | `README.md` | 572 | 408 | the most-starred list of GitHub automation, with the repository settings and the branch rules (28,305 stars) |

## The searchers

| Searcher | Model | Harness and mode | Output |
| -------- | ----- | ---------------- | ------ |
| A | GPT-6 Sol (`gpt-6-sol`, the reasoning tier) | Devin 3000.11.3, `devin -p --permission-mode dangerous`, with a `HOME` of its own that holds only Devin's credential and configuration files | [`searcher-A.md`](searcher-A.md) |
| B | Claude Sonnet 5.5 (`claude-sonnet-5-5`) | Claude Code 2.1.296, `claude -p --permission-mode dontAsk`, read-only tools and web search, writes only under `out/` | [`searcher-B.md`](searcher-B.md) |

The order of tries (S1 of the plan, with note 2 of its review): A on Devin's
paid models, which answered; B on OpenCode first ("Invalid credential"), then
Claude Sonnet 5.5, a model other than the author's. Both tries of B did not
fail, so Claude Opus 5.5 ran no search.
