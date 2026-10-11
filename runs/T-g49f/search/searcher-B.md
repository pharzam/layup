# T-g49f — the search of M2c, searcher B, word for word

Searcher: Claude Sonnet 5.5 (`claude-sonnet-5-5`) on Claude Code 2.1.296 (`claude -p --permission-mode dontAsk`, read-only tools, web search and fetch, writes only under `out/`; `--setting-sources project,local`, `--no-session-persistence`, `ANTHROPIC_DEFAULT_HAIKU_MODEL=claude-sonnet-5-5`), a fresh session, 2026-10-10T16:59:29Z to 2026-10-10T17:02:57Z (3 min 26 s by its result event); USD 2.08; `modelUsage`: `claude-sonnet-5-5` only, 2,562,338 tokens; no permission denied. A model other than the author's: OpenCode, tried first, failed with "Invalid credential" (2026-10-10 at 16:51Z). **A first try** (16:56:05Z to 16:59:14Z, USD 1.47) ran the search but could not write its file: the allow rule `Write(out/**)` did not match the path, the author's error; the second try named the directory by its absolute path. Brief: [`brief.md`](brief.md), with the five lists of [`sources.md`](sources.md) and [`gh-search.tsv`](gh-search.tsv). The output follows word for word inside a fenced block; its paths are relative to the searcher's directory.

~~~~text
# Public-solution search for M2c (T-g49f), searcher B

Date: 2026-10-10. Judged only by what a source says. No candidate code was run or installed. "Not stated" means the page I read did not say. Last-activity dates marked (gh) come from the `pushed` column of `sources/gh-search.tsv`. The GitHub pages I fetched often show no dates.

Summary in five lines:

1. Nothing found does the whole job, and nothing is a Go dependency.
2. The best patterns:
   - **W3C Web Annotation selectors** for the quote-and-span format.
   - **Duvet** for a quote checked against a source text.
   - **Squelette** for a stdlib-only controller whose decisions and evidence live in Git.
   - **rtmx / doorstop / Kibi / Agile-V** for IDs, trace and coverage checks.
   - **Smella, NALABS and the Berry ambiguity handbook** for rule checks and the ambiguity taxonomy.
3. The question merge, the answer-completeness check and the "no configuration value without evidence" check have no match.
4. The gap-of-meaning work is model-based everywhere; only the checks after it are deterministic.
5. Several "keep" entries are pattern-only because of license (AGPL, CC BY-SA, LGPL) or because the source is a paper.

## Queries

My own queries (web search, standard mode) and then the search-file queries. "From" is the fact words the query comes from.

| # | Query | From | Found |
|---|---|---|---|
| W1 | awslabs duvet requirements traceability specification quote annotation citation | F-0003#62, F-0003#51 trace to text | Duvet docs and repo: a spec URL line, then quoted requirement text lines in code comments. Bidirectional spec-to-code link. Apache-2.0, Rust. |
| W2 | W3C Web Annotation TextQuoteSelector TextPositionSelector robust anchoring | F-0003#41 byte-exact quote; the facts split into spans | W3C selectors (exact/prefix/suffix; start/end), Hypothesis fuzzy anchoring post, W3C robust-anchoring slides, Apache Annotator, Anchorpoint (Python). |
| W3 | EARS Easy Approach to Requirements Syntax Mavin INCOSE guide to writing requirements rules | F-0003#51 acceptance criterion, requirement form | EARS (Rolls-Royce, RE'09): five patterns, keyword ruleset. Nothing on the INCOSE guide came back. |
| W4 | Femmer requirements smells Smella ISO 29148 automatic detection paper | F-0003#14 metrics with no method, vague terms | Nine Smella smells (subjective language, ambiguous adverbs/adjectives, loopholes, non-verifiable terms, superlatives, comparatives, negatives, vague pronouns, incomplete references) via an arXiv 2305.07097 summary. Precision 59%, recall 82%. |
| W5 | spec-kit /speckit.clarify NEEDS CLARIFICATION marker spec template | F-0001#11 gap questions; F-0003#41 one batch | Third-party mirrors of spec-kit templates: `[NEEDS CLARIFICATION]` markers capped at 3, clarify asks up to 5 questions, answers appended as `Q → A` under `## Clarifications` / `### Session date`. |
| W6 | Femmer "Rapid quality assurance with Requirements Smells" arXiv license Smella open source | F-0003#14 | Paper arXiv:1611.08847 (2016, JSS). Open access in the Stuttgart repository. No statement on the tool's source or license. |
| W7 | Berry Kamsties Krieger ambiguity handbook | F-0003#14 terms with two readings | Handbook (Nov 2003, v1.0): linguistic sources of ambiguity. PDF at the Waterloo site. |
| W8 | Improving requirements completeness LLM Luitel Hassani Sabetzadeh | F-0003#41, F-0001#10 missing intent value | arXiv 2308.03784 (2023, v2 Feb 2024; RE 2024 per the AwesomeLLM4SE list). BERT masked-LM predictions find missing terms. 40 PURE specs. |
| W9 | StrictDoc requirements tool traceability text files git SDoc license | F-0003#62, F-0001#39 | StrictDoc: SDoc text files, Git storage, export to HTML/ReqIF/JSON, Apache-2.0. |
| W10 | LLM clarifying questions underspecified requirements ambiguity detection coding agents benchmark 2025 | F-0003#41 | Ambig-SWE (2502.13069), HumanEvalComm, ClarifyCodeBench (2607.00711), ClarEval (2603.00187), ClarifyCoder (2504.16331), "Ask or Assume?" (2603.26233). All papers; none a tool. |
| F1 | (fetch) github.com/awslabs/duvet | F-0003#62 | License, annotation syntax. Date not stated. |
| F2 | (fetch) github.com/rtmx-ai/rtmx | F-0003#62 | Go, Apache-2.0, CSV in Git. |
| F3 | (fetch) doorstop, NALABS, spec-kit, OpenSpec, ouroboros, Squelette | F-0003#51, #14, F-0001#11 | See Candidates. |
| F4 | (fetch) requirements-intake-agent, quizzing, intent-verify, Kibi, Agile-V, lobster | F-0003#41, #51, #62 | See Candidates. |
| F5 | (fetch) w3.org/TR/annotation-model | F-0003#41 | W3C Recommendation, 23 Feb 2017. First 100000 of 136619 characters read. |
| G1 | problem statement gap detection | F-0003#41 | Off-target: student projects named "Problem Statement". |
| G2 | requirements ambiguity detection | F-0003#14 | Student/replication repos; one taphsir (AGPL). No usable tool. |
| G3 | ambiguity detection natural language requirements LLM | F-0003#14, F-0001#11 | One unrelated repo. |
| G4 | requirement smells detection | F-0003#14 | Student repos, prd-agents-framework (MIT). |
| G5 | requirements quality checker | F-0003#15 | NALABS (MIT) and scenario-quality-checker variants. |
| G6 | clarifying questions LLM agent | F-0003#41, F-0001#11 | quizzing (MIT), tuplet (MIT); the rest off-target. |
| G7 | requirements elicitation LLM | F-0001#11 | Research artefacts; no controller-grade tool. |
| G8 | requirements engineering LLM | F-0003#51 | 4D-ARE, Agile-V skills, ReqPipeline; mostly agent frameworks. |
| G9 | requirements extraction from text LLM | F-0003#24 | One off-target repo. |
| G10 | requirements traceability | F-0003#62 | doorstop, duvet, lobster, reqflow, rtmx, SAP monitor. |
| G11 | requirements traceability matrix generator | F-0003#62 | Test-plan RTM generators; nothing text-anchored. |
| G12 | spec driven development | F-0003#51, F-0003#23 | OpenSpec, GSD, Agent OS, spec-workflow-mcp, and others. |
| G13 | PRD generator LLM | F-0003#23 | PRD generators; one has a code-scored critic gate. |
| G14 | acceptance criteria generation LLM | F-0003#51 | qa-agent reviews stories for untestable criteria; otherwise nothing. |
| G15 | standoff annotation format | F-0003#24, F-0003#62 | brat standoff tools; format idea only. |
| G16 | requirements completeness check | F-0003#25 | Mostly off-target; requirements-intake-agent is on-target. |
| G17 | awesome requirements engineering | the curated lists | Two more awesome lists. |
| G18 | awesome LLM software engineering | the curated lists | Survey lists; off-target. |
| G19 | configuration value evidence provenance | F-0001#4, F-0003#63 | One off-target governance repo. |

Rows G1-G19 are the 19 distinct query strings in `sources/gh-search.tsv`. The brief says 20; I found 19. Rows W and F are mine.

## Candidates

All "state in Git" and "deterministic" values are as the source says. "Not stated" if the page I read said nothing.

| Name and URL | License | Language | Last activity | Part of the problem | State in Git | Deterministic | Keep as, and why |
|---|---|---|---|---|---|---|---|
| W3C Web Annotation Data Model (TextQuoteSelector, TextPositionSelector) https://www.w3.org/TR/annotation-model/ | W3C document use terms (copyright W3C); a spec, not code | Not code | Recommendation 23 Feb 2017 | Facts as spans with exact positions, and a quote with context. `exact`+`prefix`+`suffix`; `start`/`end` half-open, first position 0. | Not applicable (a data model) | Yes, a data model | **Format.** Closest public match to "span by position plus exact quote". Adapt: the spec counts Unicode code points and normalizes text; the brief needs byte-exact equality, so use byte offsets and no normalization. Position-only is brittle when the text changes, which fits a check by code. |
| Hypothesis "Fuzzy Anchoring" post, W3C Robust Anchoring slides https://web.hypothes.is/blog/fuzzy-anchoring/ | Not stated | Not code | 2014 slides; blog date not stated | Re-anchoring order: position, then quote with context, then fuzzy | Not applicable | Not stated | **Method.** Use a strict version only: exact quote must equal the file's bytes at the span, else fail. Do not copy the fuzzy fallback; the brief wants a mechanical check. |
| Duvet (awslabs) https://github.com/awslabs/duvet | Apache-2.0 | Rust (needs Rust toolchain; annotations are language-agnostic comments) | Repo page: 126 commits; 2026-09-13 (gh); release date not stated | Trace from requirement text quoted in comments to a spec section; coverage report | Annotations in the repo's source comments; `.duvet` and `specs` dirs. Page does not say outright. | Not stated (core coverage model formally verified with Verus, per the page) | **Pattern (and separate tool at most).** The quote-line convention (`//=` URL line, then `//#` quote lines) is the nearest to "each requirement quotes its fact text and a tool checks the quote". Not usable as a dependency (Rust). Its input is an RFC-style spec, not a problem-statement file. |
| RTMX https://github.com/rtmx-ai/rtmx | Apache-2.0 | Go (`CGO_ENABLED=0`); two outside modules (Cobra, a YAML parser) | 2026-10-04 (gh); 373 commits; no releases listed | Requirement ID, spec, linked tests; status from test results; CSV database | CSV in Git (`.rtmx/database.csv`) plus `rtmx.yaml` | Not stated. Status is derived from tests. | **Pattern.** Same language and Git-first. Cannot be a dependency (Cobra, YAML parser, outside the stdlib-only rule). Shows a row-per-requirement store whose diffs read well. The Python predecessor is end-of-life 2026-09-25. Benchmark gain not significant (p=0.206), so no evidence it helps agents. |
| Doorstop https://github.com/doorstop-dev/doorstop | README says LGPLv3; the gh file says NOASSERTION | Python 3.10+ | 2026-10-10 (gh); release dates not stated | Identified items, links between documents, integrity check of the link tree | One YAML file per item in the repo; needs a VCS | Not stated | **Pattern.** Item-per-file with parent links. "Suspect link" idea was not on the page I read; do not assume. Python, so not a dependency. |
| StrictDoc https://github.com/strictdoc-project/strictdoc | Apache-2.0 (per docs) | Python | Not stated | Requirements as text (SDoc), traceability to source (experimental), exports incl. ReqIF/JSON | SDoc text files stored in Git; says it shall support Git | Not stated | **Format idea.** Text-file requirements with typed fields. Not a dependency. Read only through a search snippet; I did not fetch the repo. |
| Lobster (BMW) https://github.com/bmw-software-engineering/lobster | AGPL-3.0 | Python (inferred) | 2026-10-09 (gh); 695 commits | Common interchange files from many sources into one report against a tracing policy | Not stated | Not stated | **Pattern only** (AGPL). Idea: each source writes one interchange format; one policy checks the whole. Fits "questions from every source merged". |
| Kibi https://github.com/Looted/kibi | AGPL-3.0-or-later | TypeScript/Bun plus bundled SWI-Prolog | Not stated; 2,826 commits; beta | Requirement-to-scenario-to-test links; unresolved clauses reported, not passed | `.kb/` in the repo (implied, not stated outright) | Checking layer yes (Prolog), authoring by LLM | **Pattern only** (AGPL). "Probabilistic interpretation with deterministic verification" is the same split as the brief: agents read meaning, code checks. |
| Agile-V skills https://github.com/Agile-V/agile_v_skills | CC BY-SA 4.0 | Markdown skills, JSON schemas, Python tests | 2026-10-05 (gh) | `REQ-XXXX` IDs with typed links; separate verifier that never checks its own work; human gate approves the requirement baseline | `.agile-v/` evidence files | Schemas and tools described as deterministic | **Pattern only** (share-alike license). Matches "second, independent agent". 40 JSON evidence schemas worth reading for field choices. |
| Squelette https://github.com/JyMinet/squelette | MIT | Python, standard library only (per the page) | Not stated; 6 commits | A controller that refuses: human decisions recorded in `docs/governance/HUMAN_DECISIONS.md`; DONE needs evidence; reports checked against Git refs and digests | Yes: decisions, evidence refs, one branch per work item | Rule-based, fail-closed; "deterministic" not claimed | **Pattern.** Closest match in shape: stdlib-only controller, human decisions as files in Git, evidence required. Six commits, two stars: immature. Does nothing on problem-statement meaning. |
| intent-verify https://github.com/hermes-labs-ai/intent-verify | Apache-2.0 | Python | v0.2.1; 49 commits; date not stated | Lexical check that acceptance items in a Markdown spec appear in a repo; `coverage-map.v1` JSON; exit codes 0/1/2 | Not stated | Yes (no model, no network, token overlap and two thresholds) | **Pattern.** Exit-code and JSON contract style. Its check is lexical overlap, not meaning; too weak for the brief's coverage check, useful for the second "uncovered" pass only as a cheap pre-check. |
| NALABS https://github.com/eduardenoiu/NALABS | MIT | C# GUI and Python (spaCy, TextBlob, textstat, pandas) | 2025-06-18 (gh) | Rule counts: vague phrases, conjunctions, references, optionality, subjectivity, weakness, imperatives, continuances, readability | Not applicable (Excel in/out) | Not stated; formula-based measures suggest yes | **Taxonomy / data source.** The metric list and the word lists give rule checks for the controller to write in Go. Not a dependency. |
| Smella and "Rapid quality assurance with Requirements Smells" https://arxiv.org/pdf/1611.08847 | Paper: arXiv/open access; tool source and license not stated | Java (not stated by the sources I read) | 2016 paper | Nine smells tied to ISO 29148 | Not applicable | Dictionary and POS based per the summary | **Taxonomy.** Nine named smells, each with a rule. Reported precision 59%, recall 82% means false positives; fits "question", not "block". |
| Berry, Kamsties, Krieger, ambiguity handbook https://cs.uwaterloo.ca/~dberry/handbook/ambiguityHandbook.pdf | Not stated | Paper | 2003 (v1.0) | Taxonomy of linguistic ambiguity: quantifiers, and/or, only/also, pronouns, that/which | Not applicable | Not applicable | **Taxonomy.** Source for "term with two readings" categories to put in the review agent's brief. Not read in full; only the search summary. |
| EARS (Mavin et al., RE'09) https://en.wikipedia.org/wiki/Easy_Approach_to_Requirements_Syntax | Not applicable | Not code | 2009 | Five clause patterns with keyword rules; acceptance of a requirement's form | Not applicable | A fixed syntax can be checked by code | **Format.** Cheap mechanical check of requirement form. Read through Wikipedia and a Medium article only. |
| OpenSpec https://github.com/Fission-AI/OpenSpec | MIT | TypeScript (inferred) | 2026-10-09 (gh) | Requirement as `SHALL` text plus WHEN/THEN scenarios; changes as folders (proposal, specs, design, tasks) | `openspec/` in the repo | Not stated; no validate command on the page I read | **Format.** Scenario form doubles as an acceptance criterion. No trace to source text in the format as described. |
| GitHub spec-kit https://github.com/github/spec-kit | MIT | Python 3.11+ (inferred) | Not stated | `[NEEDS CLARIFICATION]` markers capped at 3, a clarify step with up to 5 questions, answers written to the spec | `.specify/` in the repo | Agent output not deterministic | **Pattern.** Shows the cost of one-at-a-time and capped questions, which the brief wants to avoid (one batch, no cap on rule gaps). The answer record (`Q → A`, dated session) is a plain model for stored answers. Template text read only through mirrors. |
| requirements-intake-agent https://github.com/yon-tyzes/requirements-intake-agent | MIT | Python 3.10+, YAML config, JSON Schema | 2026-10-08 (gh); 1 commit, 0 stars | Gap checks from `config/quality-checks.yaml`; clarifying topics from `config/clarifying-topics.yaml`; lines tagged `[SOURCE]`, `[INFERRED]`, `[PRIOR]`, `[MISSING]`, `[CONFLICT]` | Not stated | Build step yes; drafting model-driven | **Pattern.** The five provenance tags map onto "no value without evidence". Config-driven gap list is a good shape. Very immature. |
| Quizzing https://github.com/Tranquility2/quizzing | MIT | Node.js (inferred) | 2026-08-01 (gh) | One capped batch of questions in a local form; answers saved as one JSON file; follow-up round from deferred/open items | JSON file; location not stated | Question generation no; form checks yes | **Pattern.** The nearest to "one batch, answers back as a file, follow-up only for open ones". Early development. |
| Ouroboros https://github.com/Q00/ouroboros | MIT | Python 3.12+, some Rust | Not stated; 2,330 commits | Interview scored per dimension (goal, constraints, success criteria); seed blocked above 0.2 ambiguity | SQLite event store, not Git | Scoring is an LLM at temperature 0.1: no | **Pattern, cautiously.** Dimensions "goal / constraints / success criteria" match the brief's "missing intent value". The weights (40/30/30) have no stated evidence; do not copy them. State is not in Git. |
| Luitel, Hassani, Sabetzadeh, completeness with masked LMs https://arxiv.org/abs/2308.03784 | arXiv perpetual non-exclusive; paper | Not stated | 2023 (v2 2024) | Finds missing terms in requirements | Not applicable | Model-based | **Method idea only.** Needs a BERT model; the controller calls no model. A coding-agent session could do the equivalent. |
| Clarification benchmarks: Ambig-SWE 2502.13069, HumanEvalComm, ClarifyCodeBench 2607.00711, ClarEval 2603.00187 | arXiv terms | Papers | 2025-2026 | Evidence that agents guess rather than ask; scoring ideas for question quality | Not applicable | Not applicable | **Data source / method** for judging the review agent later. Not read beyond the search summaries. |
| SAP quality-continuous-traceability-monitor https://github.com/SAP-archive/quality-continuous-traceability-monitor | Apache-2.0 | Go | 2025-07-18 (gh); archived by name | Links requirements to automated tests and results | Not stated | Not stated | **Pattern, weak.** Go, but archived and test-centred. Not opened beyond the search row. |

## The lists

### aabs/awesome-specification-driven-development (commit 4fa40ba)

- GitHub Spec-Kit: keep: question/clarify pattern; see Candidates
- OpenSpec: keep: requirement and scenario format
- SpecPulse: reject: Python SDD framework; no gap or trace check stated
- Kiro-base SDD: reject: guide for requirements/design/tasks; nothing checkable
- Specification by Example: reject: method book; no tool
- Contract-Driven Development / Specmatic / Pact entries: reject: API contract testing, not prose requirements
- Swagger/OpenAPI/oapi-codegen entries: reject: API specs, not problem statements
- OpenAPI Diff: reject: diff of API specs
- Concordion, JGiven, Lettuce, Mamba: reject: BDD test frameworks in other languages
- Augmented Coding Patterns: reject: patterns for coding with AI, not intake
- AgentScope, MetaGPT, AutoGen, Langroid, Serena, LLM4S: reject: agent frameworks; controller calls no model
- DafnckMachine AGF 2.0: reject: orchestration framework
- MCP server entries (Awesome MCP lists, Context7, Qdrant, Excel, Exa, ArXiv, etc.): reject: unrelated servers
- Blog posts (GitHub SDD article, limits of SDD): reject: opinion; not read
- Shriek-fx: reject: .NET DDD framework

### Engineering4AI/awesome-spec-driven-development (commit ec93c1e)

- adversarial-spec: reject: LLM debate refinement; model-driven, no state in Git stated
- Archcore: reject: context for coding agents; no gap/trace check stated in the entry
- Cage: reject: TypeScript; deterministic contract harness for code, not prose
- dotdog: reject: its own spec language and graph; not seen to anchor to a prose file
- DriftWire: reject: API drift
- gsd-core: reject: meta-prompting system
- idea-to-blueprint: reject: agent skill generating a spec from an idea; no quote check stated
- intent-verify: keep: exit-code/JSON contract; lexical check
- Kibi: keep: pattern only (AGPL)
- lean-spec, MetaSpec, MySpec, openspec-scaffold: reject: spec templating
- quint-code: reject: reasoning framework; entry says nothing checkable
- reqlan: reject: own requirements language and LSP; entry gives no anchoring to prose; not opened
- Shotgun, spec-driver, Spec Kit: Spec Kit keep (pattern); others reject: codebase-aware spec writing
- Squelette: keep: stdlib-only controller; decisions and evidence in Git
- vibe-coding-prompt-template: reject: prompt templates
- agent-skills, BMAD-METHOD, BusinessLens, cc-sdd, claude-codepro, colign, Consort, Cosmosmith, fspec, MoAI-ADK, pi-sdd-kit, skills, SpecKit Companion, SpecWeave: reject: workflows around coding agents, not problem-statement checks
- ouroboros: keep: interview ambiguity dimensions; state not in Git
- THROUGHLINE: reject: personas debate; append-only record is interesting but unread
- Kiro, Notula, spec-kit-command-cursor: reject: editors/IDEs
- MCP entries (@myspec, reqlan MCP, compendio-mcp, mcp-server-spec-driven-development): reject: server wrappers
- Aide, GraphCode, Okto Pulse, Plans, RailWarden, Spec Kitty, taskmaster, trailhead, vibe-kanban, VibeDoc, workkit, YYLO: reject: task/workflow management after the spec
- Introductions (O'Reilly spec article, SDD book, spec-first workflow, practitioner guide with EARS): reject: reading, not a solution; EARS covered separately
- agents.md: reject: agent instruction file
- FPF: reject: framework; nothing checkable stated
- Awesome Harness Engineering: reject: another list

### iSEngLab/AwesomeLLM4SE (commit f486410, 1711 lines, read in full)

Almost all entries are code models, code generation, testing, repair and vulnerability work; none touch the problem. Those that do (all papers; none given a repository in the list):

- Ambiguity detection: Anaphoric Ambiguity (ICSE 2022): keep: taxonomy item (anaphora); GPL-3.0 tool per the NLP4RE list
- Ambiguity detection: Requirement contradiction detection via formal logic and LLMs (2024): reject: formal logic from prose; not read
- Ambiguity detection: Intra-domain ambiguity with transformers; TABASCO: reject: model-based; no source seen
- Requirement Completeness Detection: Improving requirements completeness (RE 2024): keep: method idea only
- Requirement Elicitation: Combining prompts with examples (COMPSAC 2024): reject: prompt study
- Requirement Engineering: MARE multi-agent (2024), NLI lessons, prompt patterns: reject: agent studies; not read
- Requirements Quality Assurance: Leveraging LLMs for QA of requirements (RE 2024): keep: possible source of a smell list for the review agent; not read
- Requirements Quality Assurance: High-level to low-level coverage reviewing with LLMs (MSR 2024): keep: matches "a second agent lists what no requirement covers"; not read
- Requirements Quality Assurance: Requirements satisfaction assessment; story quality with ChatGPT: reject: model-based, not read
- Requirement Traceability: NLP for Requirements Traceability (2024): reject: link recovery between artifacts; the brief has links by construction
- Traceability Link Recovery: Enhancing recovery with unlabeled data: reject: same
- Specification Generation (SpecGen etc.): reject: formal program specs
- Requirement Classification (NoRBERT, PRCBERT, others): reject: classes are FR/NFR; the brief's classes (need, constraint, measure, context) differ
- Software tool configuration: Can LLMs configure software tools (2023): reject: opposite of "no value without evidence"
- Class diagram, GUI layout, summarization, prioritization, use-case, modelling entries: reject: off-target

### JulianFrattini/nlp4re-tools (commit ac754cd, 126 tools, read in full)

- R124 Anaphoric Ambiguity Detector: keep: taxonomy; GPL-3.0, tool is model based
- R129 MAANA: keep: taxonomy (coordination and prepositional-phrase ambiguity); MIT, Python (per the NLP4RE list's license; language not in the list)
- R021 Ambiguous Words Detector, R030 MWE Disambiguator: reject: no source link; word-embedding based
- R162 REVV-Light: keep: taxonomy; terminological ambiguity and near-synonyms in user stories; no license in the list
- R048 Use Case Smell Detector, R080 Requirements Quality Classifier, R157 Requirements Rewriter: reject: no source link in the list
- R073 NERO: reject: no link; feature requests
- R047 Requirements demarcator: reject: classifies requirement sentences vs prose with ML; the brief splits by exact span
- R148 ORCS (OpenReq): reject: classifies prose vs requirement; EPL, ML
- R103 Template Converter (EARS/Rupp): reject: translation tool, no stated license; EARS itself kept
- R040 Potential verification method classifier: reject: ML; GPL-3.0
- R149 CiRA, R028 CEREC, R140, R031, R036: reject: causality extraction; could help acceptance criteria but model-based
- R071 BTC EmbeddedPlatform, R081 ARF, R125 DeepSTL, R154 Req2Spec, R017: reject: formalisation to temporal logic
- R002 OpenReq-DD, R059, R062, R199: reject: dependency detection between requirements
- R024 TT-RecS: reject: taxonomy-based trace recommender; Apache-2.0 but a Zenodo artefact, not read
- R009 ConDec: reject: Jira decision knowledge; MIT
- R118 QAssist: reject: question answering over SRS, model-based
- R061 GDPR Completeness Checker, R111 CompAI: reject: privacy policy completeness
- R176 ScenarioAmigo: reject: omitted use-case steps; no deterministic check described
- R088 MRDQA: reject: classifies document quality with a model
- R100 BugBot, R046 LadderBot: reject: chatbots asking missing context; one-at-a-time, model-based
- All other entries (user review mining, classification, trace link recovery, app store, tweets, privacy, GUI, story points, security classifiers): reject: off-target

### muchengracedriver/awesome-requirements-writer (commit 1584110, v0.2.0)

- awesome-requirements-writer: keep: pattern only; MIT, Node.js installer, `SKILL.md` plus templates. Principles list (atomic, no magic numbers, directly testable, separate requirement from information, traceable, open questions for missing values instead of invented thresholds) is a usable rule list for the review brief. The skill is for agents and the controller cannot use it directly. Language of the checks is natural language; no deterministic check described.

## The search file

`sources/gh-search.tsv`, 135 result rows, read in full. Entries that touch the problem:

- OpenSpec (Fission-AI): keep: format; MIT, 2026-10-09
- rtmx-ai/rtmx: keep: pattern; Go, Apache-2.0
- doorstop-dev/doorstop: keep: pattern; item-per-file, license flag NOASSERTION vs README LGPL
- awslabs/duvet: keep: pattern; quote lines
- bmw-software-engineering/lobster: keep: pattern only (AGPL)
- eduardenoiu/NALABS: keep: taxonomy/data source; MIT
- Agile-V/agile_v_skills: keep: pattern only (CC BY-SA)
- yon-tyzes/requirements-intake-agent: keep: gap tags and config shape; immature
- Tranquility2/quizzing: keep: one-batch question and answer-file shape
- SAP-archive/quality-continuous-traceability-monitor: keep: weak; Go, archived name
- vshidlovsky/prd-agents-framework: reject: Claude Code PRD agents; says "smell detection" but no checkable output seen; MIT
- sankalpdusane/ai-prd-generator-eval-critic: reject: critic scores "in code" but for generated PRDs; not the brief's direction
- SondosAlgyed/qa-agent: reject: reviews Jira stories for untestable criteria; TypeScript; no license
- Pimzino/spec-workflow-mcp: reject: MCP dashboard; GPL-3.0
- gsd-build/get-shit-done, buildermethods/agent-os, zhu1090093659/spec_driven_develop, mouredev/hello-sdd, marcelsud/spec-driven-agentic-development, deeplearning-ai course files: reject: coding-agent workflows or courses
- formulahendry/mcp-server-spec-driven-development: reject: prompt server
- Engineering4AI/awesome-spec-driven-development, muchengracedriver/awesome-requirements-writer, iSEngLab/AwesomeLLM4SE: already read as lists
- jgsystemsconsulting/awesome-requirements-engineering (CC0, pushed 2026-10-10): reject: not in my lists; a further list I could not open in this run, worth a look by the next pass
- ishandutta2007/Awesome-Engineering-Requirements-Management, codefuse-ai/Awesome-Code-LLM, fatihkc/awesome-agentic-engineering, other awesome lists: reject: lists, not tools
- SNTSVV/taphsir: reject: anaphoric ambiguity, AGPL-3.0, model based
- arsalaniltaf/RequirementSmells, MetricBlood/srs-ambiguity-detection, haydarozden/Requirements-Quality-Checker, other student ambiguity repos: reject: coursework, no license or Python ML
- profzivkovic replication package (MIT): reject: replication of a paper, not a tool
- damianhorna / Gho-Ost / czajkowskis scenario-quality-checker: reject: use-case scenario analysers
- goeb/reqflow (C++), tobhey/finegrained-traceability, ociu/sphinx-traceability-extension: reject: other languages or link recovery
- adrianpuiu/specification-document-generator: reject: six-phase LLM prompts; no license
- AungMyoKyaw/prd-creator, cdeust/ai-prd-generator, and other PRD generators: reject: generation only
- santkumar333/llm-test-case-generator, Falak359 RTM generator, and others: reject: test-case traceability
- GateNLP brat plugin, brat tools: reject: Java/Python tools; keep only the idea of standoff format, already covered by W3C selectors
- yon-tyzes: counted above
- rabbittsui/ModelCardDiff: reject: model cards; "requirements schema" not prose
- ybeven/4D-ARE: reject: agent requirements for LLM agents; different problem
- GoldstonianConcordanceBible/402: reject: governance text; not read in depth
- Rahul20037237/AMBI-agent: reject: entropy and MDP choosing when to ask; no license
- anetrebskii/tuplet: reject: agent library with clarifying questions; TypeScript
- All other rows (nutrition, job match, voice assistants, graduation checklists, Godot, HIPAA lists): reject: off-target

## What no candidate covers

1. **Span split with a mechanical cover check.** I found no tool that splits a prose file into numbered spans, classes each (need, constraint, measure, context), and checks that every non-blank byte lies in exactly one span and each span equals the file text. W3C selectors give the record shape only. The cover check (no overlap, no gap except blank) has no public match.
2. **Quote check for gaps of meaning.** Duvet checks quoted lines against a spec in its own format; nothing checks a model's quote against the file byte for byte as a gate on gap findings.
3. **Gap taxonomy built for problem statements.** The smell and ambiguity taxonomies (Smella, NALABS, Berry) are for requirement sentences. Missing intent value, risky assumption and open decision have no public taxonomy that I found; "metric without a measurement method" is closest to the Smella "non-verifiable term" smell.
4. **Trace from requirement to fact IDs with a coverage check by an independent agent.** Agile-V, Kibi and rtmx trace requirements to tests or code, not to numbered facts of a problem statement. The "second agent lists what no requirement covers" step has no tool; the LLM4SE paper on coverage reviewing is the nearest and was not read.
5. **Merging questions from many sources.** Nothing found merges rule-check, meaning-review, completeness and setup questions into one batch with IDs and removes duplicates. Quizzing has the single-batch idea for one source only. Duplicate removal across sources has no match.
6. **Answer completeness and follow-up for missing or open answers only.** Quizzing builds a follow-up round from deferred and open items (page says so); the completeness check as a mechanical rule is not described. Spec-kit appends `Q → A` but caps and asks sequentially.
7. **Answers stored as a raw fact in Git.** Squelette's decisions file and spec-kit's clarification section are the nearest; neither stores a raw, unedited answer as a separate fact with a digest.
8. **No configuration value without evidence.** I found no tool that checks that each configuration value cites its source. The requirements-intake-agent `[SOURCE]/[INFERRED]/[MISSING]` tags are the nearest idea; Squelette's evidence-per-gate rule applies to work items, not config values.
9. **Go and stdlib only.** The only Go hits are RTMX (with two outside modules) and an archived SAP monitor. No candidate can be a dependency under the rule.

## Limits

- **Dates.** GitHub page fetches did not show commit or release dates. I used the `pushed` dates from the search file where the repo was in it, and "not stated" for the rest (spec-kit, Kibi, Squelette, ouroboros, StrictDoc, intent-verify, Smella).
- **Language.** The fetched pages had empty language lists. Languages marked "inferred" come from file names and install commands; I did not see a language breakdown.
- **Search tool.** WebSearch is US-only, so non-English or non-US sources may be missed. Search results returned snippets and summaries, not the full pages. I ran 10 web searches plus fetches, not an exhaustive sweep.
- **Not opened.** I did not open reqlan, dotdog, Cage, Archcore, THROUGHLINE, adversarial-spec, MySpec, jgsystemsconsulting/awesome-requirements-engineering, the Berry handbook, the Luitel paper or the clarification benchmark papers; entries on them rest on the list line or a search summary.
- **W3C spec.** I read the first 100000 of 136619 characters of the Web Annotation Model; the selector definitions were in that part. Not checked: whether any W3C-compatible Go code exists.
- **Smella.** The 59%/82% precision and recall come from a 2023 arXiv paper's description of the 2016 tool, not from the original paper.
- **Query count.** The brief says the search file has 20 queries; I count 19 distinct query strings in the file, and list those 19. I may have miscounted one; the file is otherwise read in full.
- **Licenses.** For Doorstop the README (LGPLv3) and the GitHub metadata (NOASSERTION) disagree; verify before using even as a pattern. For Duvet, Squelette, Kibi, Agile-V and others the license is the one on the repository page I fetched, not verified in the LICENSE file.
- **Determinism.** I did not run anything, so every "deterministic" value is the source's own claim or "not stated".
- **AwesomeLLM4SE.** Read in full, but its entries are paper titles with no repository links for the requirements section, so none of those entries could be checked for license or code.
~~~~
