# T-g49f — the public-solution search of M2c: the sources

The method of [`docs/guardrails.md`](../../../docs/guardrails.md) §2, "A
public-solution search that confirms the design", as `T-ywk7` ran it for `M2b`:
the search runs before any decision that fixes the shape of `M2c`; its queries
come from the facts that the requirements of `M2c` (REQ-001, REQ-012, NFR-003)
cite in `PRD-0001` §12, and from the `M2c` row of the plan; each curated list is
read in full; two searchers on different models, blind to each other. The
author built these sources on 2026-10-10; the two searchers got
[`brief.md`](brief.md), the five lists and [`gh-search.tsv`](gh-search.tsv)
only: not the architecture, the specification, or the name of any candidate.

## The repository search

[`gh-search.tsv`](gh-search.tsv): 20 GitHub repository searches (the API's best
match, 10 results each where it gave 10), on 2026-10-10 at 16:53Z, each query
with the fact words it comes from: 135 results.

## The curated lists

Chosen from the two list queries of the search file (`awesome requirements
engineering`, `awesome LLM software engineering`) and three more, on 2026-10-10:
`NLP4RE`, `awesome spec driven development`, `awesome requirements`. The rule:
keep the most-starred list of each domain that the problem touches (the
derivation of specifications, LLMs in software engineering, the writing of
requirements, NLP for requirements engineering), whose entries are tools,
methods or papers rather than prompts. Each kept list is its file at the commit
below, given to the searchers in full; the text is not stored here, and the
commit makes it reproducible.

| List | Commit | File | Lines | Links | Why kept |
| ---- | ------ | ---- | ----- | ----- | -------- |
| `Engineering4AI/awesome-spec-driven-development` | `ec93c1ecc372` | `README.md` | 157 | 130 | the most-starred list of spec-driven development (293 stars) |
| `iSEngLab/AwesomeLLM4SE` | `f486410a27c5` | `README.md` | 1711 | 70 | the most-starred list of LLMs in software engineering with a requirements section (348 stars) |
| `muchengracedriver/awesome-requirements-writer` | `158411013a3a` | `README.md` | 161 | 6 | the most-starred list of requirement writing (104 stars) |
| `JulianFrattini/nlp4re-tools` | `ac754cdbe470` | `tools/nlp4re-tools.md` | 270 | 102 | the classified list of NLP4RE tools (Frattini et al., the NLP4RE handbook); its list is tools/nlp4re-tools.md, not its README |
| `aabs/awesome-specification-driven-development` | `4fa40ba418f5` | `README.md` | 166 | 102 | a second list of specification-driven development, of its methods |
