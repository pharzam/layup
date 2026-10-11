# T-g49f — the search of M2c, searcher A, word for word

Searcher: SWE-2 High (`swe-2-high`, Devin's free model, O-179; an execution-tier model on a reasoning part, as Devin's paid quota was exhausted, `resource_exhausted`) on Devin 3000.11.3 (`devin -p --permission-mode dangerous --respect-workspace-trust false`), a fresh session with a `HOME` of its own that held only copies of Devin's credential and configuration files (removed after the run), 2026-10-10T17:00:08Z to 2026-10-10T17:12:21Z (12 min 13 s), exit 0. Devin reports no usage. **A first try** (16:56:02Z to 16:59:49Z), with `--permission-mode accept-edits`, ended with three tool calls rejected ("requires confirmation") and no file; that mode was the author's error, and the second try took the mode of `T-ywk7` (`runs/T-ywk7/search/sources.md`). Brief: [`brief.md`](brief.md), with the five lists of [`sources.md`](sources.md) and [`gh-search.tsv`](gh-search.tsv); nothing of the LAYUP architecture or the specification. The output follows word for word inside a fenced block; its paths are relative to the searcher's directory.

~~~~text
# Public-solution search — M2c intake controller (task T-g49f, searcher A)

Scope reminder: the controller is Go + stdlib + `git` as a subprocess; meaning
reading is done by external coding-agent sessions. Therefore no candidate can
be a dependency; keeps are patterns, methods, formats, taxonomies, data
sources, or separate tools.

## Queries

| Query | Words it comes from | What it found |
|---|---|---|
| requirements ambiguity detection tool natural language taxonomy | F-0003#14 terms with two readings; F-0003#41 gaps | mcp-reqcheck (deterministic rule catalog + ISO 29148 citations + question templates), QuARS, SEAMS 13-category scanner, TABASCO, BA-requirement-analyzer |
| requirement smells taxonomy Femmer requirements quality | F-0003#14 vague words, metrics with no method | Femmer et al. JSS 2017 "Rapid quality assurance with Requirements Smells" + Smella prototype; smells derived from ISO 29148 language criteria |
| EARS Easy Approach to Requirements Syntax Mavin template | F-0001#39 requirement format; list mentions | Official EARS guide + RE'09 paper: five sentence patterns (While/When/Where/If-Then/Ubiquitous shall) |
| brat standoff annotation format character offsets text span | F-0003#24 numbered spans; F-0003#62 trace to text | brat standoff spec: `T{id}\t{type} {start} {end}\t{text}`, character offsets + stored text for reference check; bratpy verifier code |
| QuARS quality analyzer requirements specifications tool download | follow-up on ambiguity tools | quars.it (freeware, no source), SEI report 2005, quality model of lexical/syntactic indicators |
| PURE dataset public requirements documents Ferrari | pilot problem statement corpus | Zenodo dataset, 79 public NL requirements documents, 34k sentences, RE 2017 paper |
| ambiguity handbook Berry Kamsties Krieger requirements engineering | F-0003#14 ambiguity taxonomy | Free PDF "From Contract Drafting to Software Specification: Linguistic Sources of Ambiguity" (2003), inspection checklist |
| requirements completeness check coverage tool specification gaps | F-0003#25 cannot be traced; second-agent completeness | flashtrace (markdown spec items, transitive coverage, missing/orphan/duplicate), CodeCharter req check/coverage, ohpauleez/spec-check, Roboticks gap analysis |
| doorstop requirements management text files git | gh-search row, trace to text stored in git | Doorstop: YAML/markdown item files in VCS, UIDs, link validation, publish; readthedocs |
| spec-kit clarify phase underspecified requirements questions | F-0003#41 one batch; F-0001#11 answers in one batch | spec-kit `/speckit-clarify` template: coverage-scan taxonomy, prioritized question queue (cap 5, MC or <=5-word answers), answers folded into spec; `/speckit-analyze`: coverage-gap report |
| architecture decision record ADR markdown git evidence | F-0001#4 no config without evidence; open decision | MADR template, adr-tools, git-adr (ADRs in git notes), Fowler article (records in repo, numbered, superseded not edited) |
| NALABS requirement quality checker bad smells tool | gh-search row | eduardenoiu/NALABS (MIT, C#) + paper arXiv:2202.05641: dictionary-based smell metrics (vague, subjective, optionality, weakness, continuances, imperatives, ARI) |
| Volere requirements specification template snow card quality gateway | F-0003#51 requirement + acceptance criterion + trace | Volere snow card attributes (ID, type, description, rationale, fit criterion, originator, trace to business event/PUC) + quality gateway testpoints (completeness, measurability, ambiguity, traceability...) |
| awslabs duvet requirements traceability markdown | gh-search row | duvet: extracts requirements from spec text, code annotations quote spec section text, reports coverage; TOML store under `.duvet/` |
| rtmx requirements traceability matrix Go agentic | gh-search row, Go candidate | rtmx: requirements register as CSV in git, `rtmx health` lints orphans/cycles/stale refs, status derived from tests, MCP for agents; Go CLI |
| ISO 29148 characteristics of good requirements quality criteria checklist | F-0003#14/#41 gap classes | Nine characteristics of individual requirements (necessary, appropriate, unambiguous, complete, singular, feasible, verifiable, correct, conforming) + set-level; SRS content list incl. assumptions & dependencies |
| requirements vague words list weasel words ambiguity lint tool | F-0003#14 vague words | wyzer-reqlint (INCOSE+EARS word lists), phoenix spec-lint (heuristic spec readiness), speccheck-core Rust rules SC-AMB-*, pi-harness-runtime ambiguity-detector (emits question+evidence per finding), the-startup ambiguity-detection.md |
| cucumber gherkin given when then acceptance criteria specification | F-0003#51 acceptance criterion | Gherkin reference: Feature/Rule/Scenario with Given-When-Then; executable acceptance criteria precedent |
| verify annotation spans cover entire document tiling check text segmentation | mechanical check: every non-blank part in exactly one span | medkit Segment (span+text, length assert), anno-eval validate_ground_truth_entities (bounds, non-overlap, text-match) — span *validity* checkers exist; no full-document partition checker found |
| requirements elicitation open questions assumption log RAID problem statement | risky assumption or open decision | RAID log (Risks, Assumptions, Issues, Dependencies/Decisions register with ID, owner, status) as the gap-class taxonomy for open items |
| mcp-reqcheck github ambiguity rules yaml | follow-up | PyPI mcp-reqcheck 0.1.2: `ambiguity_rules.yaml` versioned catalog, regex detection, MCP surface; license not stated |
| wyzer-reqlint github requirements linter | follow-up | Wyzer-it/reqlint: single-file dependency-free Python lint, MIT; word lists from INCOSE Guide to Writing Requirements + EARS |

## Candidates

| Name (URL) | License | Language | Last activity | Part of problem | State in Git | Deterministic | Keep as + why |
|---|---|---|---|---|---|---|---|
| brat standoff format — https://brat.nlplab.org/standoff , https://github.com/nlplab/brat | MIT | Python/JS | last push 2024-07 (format stable, widely reused) | Numbered fact spans: `type start end` + stored span text; mechanical quote check | `.ann` files beside `.txt` → yes | format only, deterministic | **format**: the span+class+verbatim-text record the facts file needs |
| harttu/annotation_consistency_checker — https://github.com/harttu/annotation_consistency_checker | not stated | Jupyter/Python | 2021-06 | Mechanical check of standoff annotation consistency | n/a | yes | **pattern**: precedent that span integrity is checkable in plain code |
| awslabs/duvet — https://github.com/awslabs/duvet | Apache-2.0 | Rust | 2026-09 | Requirement↔spec-text trace with verbatim quote; coverage + duplicate-annotation checks | annotations in code + `.duvet/*.toml` → yes | yes | **method/tool**: closest existing "trace to the text that states it" + text-match check |
| doorstop — https://github.com/doorstop-dev/doorstop | LGPL-3.0 | Python | 2026-10 | Requirement items as UID'd text files, links, integrity checks, publish | yes, designed for VCS | yes | **format/pattern**: identified requirements with IDs living in project Git |
| rtmx-ai/rtmx — https://github.com/rtmx-ai/rtmx | Apache-2.0 | Go | 2026-10 | Requirements register (CSV in git), `rtmx health` lints orphans/cycles/stale refs | yes (CSV + config in repo) | yes | **format**: single-file requirements register in git; Go CLI precedent |
| goeb/reqflow — https://github.com/goeb/reqflow | GPL (COPYING) | C++ | 2025-11 | Regex-extracts requirement IDs from docs, forward/backward trace, coverage report | docs + `.req` config in repo | yes | **pattern**: deterministic regex trace + coverage errors |
| flashtrace — https://flashtrace.github.io/ | Apache-2.0 | TypeScript | active 2026 | Markdown spec items with IDs + Needs lists; transitive coverage; missing/orphan/duplicate verdicts; JSON + exit codes | spec files in repo | yes | **pattern**: "what no requirement covers" verdict shape for the completeness pass |
| MADR — https://github.com/adr/madr | not stated (repo metadata) | Markdown | 2026-08 | Record of a decision: context, options, decision, consequences | `docs/decisions/*.md` → yes | format | **format**: an owner's answer/decision recorded in git with rationale |
| git-adr (zircote, ex shivasymbl) — https://github.com/zircote/git-adr | MIT | Rust | 2026-09 | ADRs stored in **git notes**, linked to commits, no files | git objects only → yes | yes | **method**: storing raw facts inside Git without touching the tree |
| mcp-reqcheck — https://pypi.org/project/mcp-reqcheck | not stated | Python | v0.1.2, 2026 | Deterministic regex rule catalog (`ambiguity_rules.yaml`), findings cite ISO 29148 clause, per-finding question templates | n/a | yes (explicitly "no model involved") | **pattern + data source**: exactly the "rule gaps → cited questions" shape; catalog content reusable |
| Wyzer-it/reqlint — https://github.com/Wyzer-it/reqlint | MIT | Python (single file, no deps) | v0.1.0, 2026 | Language lint: ambiguous-word, weak-modal, missing-unit, vague-pronoun; INCOSE + EARS word lists | n/a | yes | **data source**: ready-made vague-word/unit rules for deterministic intake checks |
| phoenix-assistant/spec-lint — https://github.com/phoenix-assistant/spec-lint | MIT | TypeScript | 2026 (1 commit) | Heuristic spec lint: clarity, completeness, assumptions, testability, edge cases; line-referenced findings | n/a | yes (zero-LLM) | **pattern**: the rule-check half of the Intake comment |
| pi-harness-runtime requirement-compiler — https://cdn.jsdelivr.net/npm/pi-harness-runtime/ (ambiguity-detector.ts) | not stated | TypeScript | 2026 | Ambiguity detector emitting `{id, question, evidence: sourceRef, category, blocking}` | n/a | heuristic, deterministic | **pattern**: the gap-record shape — question plus file-anchored evidence |
| github/spec-kit — https://github.com/github/spec-kit | MIT | Python + shell/prompt templates | 2026-10 | `/clarify`: coverage-scan taxonomy → prioritized question queue; `/analyze`: coverage gaps, duplicates, ambiguity, constitution violations | spec.md/plan.md/tasks.md in repo | prompts only (agent executes) | **pattern**: scan taxonomy + batched-question protocol + independent coverage pass |
| Fission-AI/OpenSpec — https://github.com/Fission-AI/OpenSpec | MIT | TypeScript | 2026-10 | Specs and change proposals as repo files | yes | CLI scaffolding deterministic | **pattern**: spec-artifact layout in Git for derived requirements |
| SEAMS ambiguity-scanner — https://github.com/ddunnock/claude-plugins (skills/speckit-generator) | MIT | agent spec (Markdown) | 2026-03 | 13 gap categories (SCOPE…ASSUMPTION, STAKEHOLDER, TRACEABILITY) with markers, Impact×Uncertainty ranking, question output format | n/a | agent-driven | **taxonomy**: ready gap classes for the meaning review |
| yon-tyzes/requirements-intake-agent — https://github.com/yon-tyzes/requirements-intake-agent | MIT | agent instructions + Python evals | 2026-10 | Intake gap checks; `[SOURCE]/[INFERRED]/[MISSING]/[CONFLICT]` provenance tags; numbered missing-info questions; grouped clarifying questions | n/a | agent-driven | **method**: intake completeness + numbered questions + evidence tagging per value |
| Tranquility2/quizzing — https://github.com/Tranquility2/quizzing | MIT | JS + agent skill | 2026-08 | All clarifying questions front-loaded into one capped form; answers persisted as JSON contract | answers file → yes | form fixed; questions agent-made | **pattern**: one question batch, structured answer file |
| vshidlovsky/prd-agents-framework — https://github.com/vshidlovsky/prd-agents-framework | MIT | Python + agent specs | 2026-08 | Smell detection (Femmer), perspective-based reading, human-gated PRD pipeline | files in repo | agent-driven | **method**: review-of-meaning roles and smell catalog for agent sessions |
| Agile-V/agile_v_skills — https://github.com/Agile-V/agile_v_skills | CC-BY-SA-4.0 | skill pack (JSON schemas + pytest) | 2026-10 | Independent ambiguity/constraint findings role; frozen requirement baseline; 40 evidence schemas; human gates | `.agile-v/` evidence files → yes | schemas deterministic; agent parts not | **method**: the "second, independent agent" role + evidence records |
| iniesohidham/idea-to-blueprint — https://github.com/iniesohidham/idea-to-blueprint | MIT | SKILL.md + scripts | 2026-09 | Interview → evidence-backed blueprint with epics + Given/When/Then acceptance criteria | bundle files → yes | agent-driven | **method**: evidence-backed requirement writing with acceptance criteria |
| roli-lpci/intent-verify — https://github.com/roli-lpci/intent-verify | Apache-2.0 | Python | 2026-09 | Deterministic zero-LLM CLI: lexical coverage of spec acceptance items vs repo, verified/partial/missing | n/a | yes (explicitly model-free) | **pattern**: the kind of mechanical check `spec check --facts` should be |
| ohpauleez/spec-check — https://github.com/ohpauleez/spec-check | Apache-2.0 + MIT | TypeScript | 2026-08 | Spec defects: ambiguity, contradictions, missing assumptions; requirements traced to source evidence; coverage | reads repo specs | mostly (SMT-backed parts) | **tool/pattern**: defect classes + requirement→source-evidence trace |
| JyMinet/squelette — https://github.com/JyMinet/squelette | MIT | Python | 2026-09 | Human decisions recorded in repo; a task cannot close without evidence verified against its artifacts | yes (pre-commit gate) | yes | **pattern**: enforced "no value without evidence" |
| ISO/IEC/IEEE 29148:2018 — https://www.iso.org/standard/72089.html | paid standard | — | 2018 | Characteristics of individual requirements and of a requirement set; language criteria; required info items incl. assumptions & dependencies | n/a | n/a | **taxonomy**: authoritative gap classes and requirement attributes |
| Femmer et al., Requirements Smells / Smella — https://www.sciencedirect.com/science/article/abs/pii/S0164121216000789 | paper (JSS 2017); Smella prototype | — | 2017 | Smell taxonomy derived from ISO 29148 language criteria; lightweight lexical detection | n/a | lexical detection is deterministic | **taxonomy/method**: the gap-of-meaning checklist basis |
| Berry/Kamsties/Krieger Ambiguity Handbook — https://cs.uwaterloo.ca/~dberry/handbook/ambiguityHandbook.pdf | free tech report | — | 2003 | Linguistic ambiguity sources + inspector checklist (quantifiers, and/or, pronouns, only/also) | n/a | n/a | **taxonomy/data source**: inspection checklist for ambiguity findings |
| QuARS — https://www.quars.it/ + SEI report | freeware, no source | C++/Windows | v4.1 (2007); papers to 2019 | Quality model: vagueness, subjectivity, optionality, weakness, implicitness, multiplicity, underspecification indicators | n/a | yes | **method/taxonomy**: indicator catalog for rule gaps |
| NALABS — https://github.com/eduardenoiu/NALABS | MIT | C# (+ Python variant) | 2025-06 | Dictionary smell metrics: vague phrases, optionality, subjectivity, weakness, continuances, imperatives, readability | n/a | yes | **method/data source**: smell dictionaries + metric formulas |
| TABASCO — https://github.com/a-moharil/TABASCO | MIT | Python | 2023-03 | Intra-domain ambiguity: same term used with different meanings in different contexts | n/a | no (embeddings) | **method**: the "term with two readings" check for agent sessions |
| Ezzini et al., anaphoric-ambiguity — https://gitlab.uni.lu/sabualhaija/anaphoric-ambiguity | GPL-3.0 | Python | ICSE 2022 artifact | Six solutions for anaphoric (pronoun-reference) ambiguity in requirements | n/a | no | **method**: reference-ambiguity gap class |
| SNTSVV/MAANA — https://github.com/SNTSVV/MAANA | MIT | Python | 2022-07 | Coordination and PP-attachment ambiguity using domain-specific corpora | n/a | no | **method**: structural ambiguity classes |
| Volere template + snow card + quality gateway — https://www.volere.org/ | proprietary template; method public in articles | — | current (template maintained) | Atomic requirement attributes: ID, type, description, rationale, **fit criterion**, trace to business event/use case; quality gateway testpoints | n/a | n/a | **format/method**: the identified-requirement record and its acceptance criterion |
| EARS (Mavin et al.) — https://alistairmavin.com/ears/ | paper RE'09; notation free | — | 2009 | Requirement sentence patterns (While/When/Where/If-Then/Ubiquitous + shall) | n/a | n/a | **method/format**: how derived requirements should be phrased |
| Gherkin / Given-When-Then — https://cucumber.io/docs/gherkin/ | MIT | format | maintained | Acceptance criteria as Given/When/Then scenarios | `.feature` files → yes | yes | **format**: acceptance-criterion syntax option |
| RAID log — https://www.atlassian.com/agile/project-management/raid-log (and PMI practice) | public practice | — | — | Risks, Assumptions, Issues, Dependencies/Decisions register: stable ID, owner, status, review date | n/a | n/a | **taxonomy**: classes for risky assumptions and open decisions |
| PURE dataset — https://zenodo.org/records/1414117 | open dataset (curator states no known license restrictions) | — | 2017/2018 | 79 real public NL requirements documents, 34k sentences | n/a | n/a | **data source**: pilot problem statements and eval corpus |
| OpenReqEU/requirements-classifier (ORCS) — https://github.com/OpenReqEU/requirements-classifier | EPL-2.0 | Python | ~2019 (OpenReq project) | Classifies parts of a specification into prose vs actual requirements | n/a | no (ML) | **method**: precedent for classing document parts |
| rabbittsui/ModelCardDiff — https://github.com/rabbittsui/ModelCardDiff | MIT | Python | 2026-09 | Gates a document against a required schema: completeness check + diff + CI fail | yes | yes | **pattern**: mechanical completeness gate on a document |

## The lists

### aabs/awesome-specification-driven-development (README.md, 166 lines, read in full)

- GitHub Spec-Kit: keep: clarify-scan taxonomy + coverage-analyze pattern.
- OpenSpec: keep: specs/changes as files in project Git.
- Kiro-base SDD guide: keep: requirements→design→tasks with EARS — process precedent.
- SpecPulse: reject: generic Python SDD scaffolding, no intake checking.
- Falcon Heavy: reject: OpenAPI-first code generation.
- Concordion: keep: precedent for executable acceptance criteria (weak).
- JGiven / Lettuce / Mamba: reject: BDD test frameworks, downstream of this part.
- Swagger/OpenAPI tools and Pact ecosystem: reject: API contract tooling.
- Specification by Example (Gojko): keep: method for example-based acceptance criteria.
- Augmented Coding Patterns / PDCA article: reject: generic agent workflow advice.
- "Markdown as a programming language" (github.blog): keep: rationale for file-based specs in git (weak).
- AgentScope / MetaGPT / AutoGen / Langroid / Serena / LLM4S / DafnckMachine: reject: agent frameworks, not intake.
- MCP servers, IDE integrations, API-first tools, community links: reject: unrelated.

### Engineering4AI/awesome-spec-driven-development (README.md, 157 lines, read in full)

- "How to write a good spec for AI agents" (O'Reilly): keep: method for what a consumable spec must contain.
- Spec-first development (aiarch.dev): keep: brainstorm→spec→plan workflow with named artifacts.
- "What Is Spec-Driven Development" (felipefontoura): keep: EARS usage + when to skip a spec.
- agents.md: reject: agent-config standard, not intake.
- FPF: reject: heavyweight formal spec framework.
- adversarial-spec: keep: multi-LLM debate = independent review of meaning (weak).
- Archcore / gsd-core / MetaSpec / lean-spec / Cosmosmith / mattpocock-skills / vibe-coding-prompt-template / agent-skills(addyosmani) / openspec-scaffold / spec-driver / Shotgun / DriftWire / compendio-mcp / mcp-server-sdd / @reqlan-mcp / spec-kit-command-cursor / SpecKit Companion / Aide / GraphCode / Plans / Spec Kitty / vibe-kanban / VibeDoc / workkit / YYLO / Okto Pulse / awesome-harness-engineering: reject: spec storage, dashboards, or agent plumbing — not intake checking.
- Cage / reqlan / Kibi: reject: requirement↔code linking; the needed link is requirement↔source text.
- dotdog: keep: deterministic spec compile + completeness validation (weak).
- idea-to-blueprint: keep: interview → evidence-backed bundle with acceptance criteria.
- intent-verify: keep: deterministic spec-coverage CLI.
- MySpec: keep: structured-interview→deterministic 4-file bundle pattern (proprietary SaaS; pattern only).
- OpenSpec: keep.
- quint-code: reject: design-rationale records, not intake.
- Squelette: keep: evidence-verified gates, decisions in repo.
- THROUGHLINE: keep: append-only git decision record across sessions (weak).
- BMAD-METHOD: keep: analyst→PM→architect staged document pipeline (weak).
- BusinessLens: reject: Git-tracked product model, post-intake.
- cc-sdd / Kiro / pi-sdd-kit: keep: three-phase spec process (requirements→design→tasks) with EARS + approval gates.
- Consort / ouroboros / RailWarden / colign / claude-codepro / MoAI-ADK / SpecWeave / superpowers / trailhead / taskmaster: reject: agent orchestration and execution gating; trailhead's decision tickets are nearest but execution-stage.
- fspec: keep: Gherkin specs linked to business rules (weak).
- Notula: reject: spec editor.
- claude-codepro: reject.
- adversarial-spec noted above.

### muchengracedriver/awesome-requirements-writer (README.md, 161 lines, read in full — single-project repo)

- The skill itself (SKILL.md + requirement-patterns-en/zh): keep: requirement-writing defect rules — unambiguous, atomic, testable, named variables for magic numbers, "open questions for missing engineering values instead of invented thresholds"; direct match for missing-intent-value gaps and no-value-without-evidence. MIT.
- Zhihu article link: reject: background essay only.
- npm/nvm install links: reject: packaging.

### JulianFrattini/nlp4re-tools (tools/nlp4re-tools.md, 270 lines, read in full)

- R002 OpenReq-DD: reject: requirement↔requirement dependency detection.
- R009 ConDec: keep: decision-knowledge classification in issues/commits → open decisions, rationale records.
- R012 Key Concepts Extraction: reject: elicitation-session transcripts.
- R021 Ambiguous Words Detector: keep: word2vec domain ambiguity — term with two readings.
- R030 MWE Disambiguator: keep: multiword-term disambiguation (weak).
- R040 PMV classifier: keep: predicts verification method of a requirement — adjacent to "metric that cannot be measured" (weak).
- R046 LadderBot: keep: self-elicitation interview chatbot — question-source precedent (weak).
- R047 Requirements demarcator: keep: separates requirements from other text — fact classification.
- R048 Use Case Smell Detector: keep: smell detection.
- R061 GDPR Completeness Checker / R111 CompAI: keep: completeness-checking method (domain-bound).
- R066 DeepCoref: keep: coreference detection — reference ambiguity (weak).
- R080 Requirements Quality Classifier: keep: quality-violation features.
- R088 MRDQA: keep: document quality analyzer (weak).
- R095 Automated QA for compliance requirements: reject: answers questions, does not find gaps.
- R100 BugBot: keep: chatbot that prompts for missing context — question-source precedent (weak).
- R103 NL2RT Template Converter: keep: NL → EARS/Rupp conversion.
- R118 QAssist: reject: QA over SRS, not gap detection.
- R124 Anaphoric Ambiguity Detector (Ezzini): keep.
- R129 MAANA: keep.
- R139 DEMAR: reject: annotation of forum text.
- R148 ORCS: keep: prose-vs-requirement classification.
- R154 Req2Spec: reject: formalization for Hanfor, not intake.
- R157 Requirements Rewriter: keep: good/bad classification + repair (weak).
- R162 REVV-Light: keep: near-synonym detection — term consistency (weak).
- R163 GuideGen: reject: test-update guidance.
- R176 ScenarioAmigo: keep: omitted-step recommendation — completeness precedent (weak).
- R186 User Story Instability Predictor / R199 linking matrix / R202 VIBE / R107 GPT2SP: reject: prediction/linking/product-line tasks, not gap detection.
- All other entries (extraction, classification, feedback mining, modeling, legal/privacy tools): no line — outside this part.

### iSEngLab/AwesomeLLM4SE (README.md, 1711 lines, read in full)

Only the Software Requirements & Design section touches the problem; per-entry lines:

- Ambiguity 1 — Anaphoric ambiguity multi-solution (ICSE 2022): keep: Ezzini artifact.
- Ambiguity 2 — Requirement contradiction detection via formal logic + LLMs (AUSE 2024): keep: contradiction = conflict-gap class (weak).
- Ambiguity 3 — Intra-domain ambiguity via transformers (NLBSE 2022): keep: term with two readings.
- Ambiguity 4 — TABASCO toolkit (SCP 2022): keep.
- Ambiguity 5 — ChatGPT ubiquitous SE tasks: reject: generic survey.
- Requirement Completeness 1 — Improving requirements completeness with LLMs (RE 2024): keep: automated "what is missing" assistance.
- Requirement Elicitation 1 — Prompts + examples for elicitation (COMPSAC 2024): keep: elicitation question generation (weak).
- Requirement Engineering — MARE multi-agent framework (arXiv 2024): keep: multi-role elicitation/analysis (weak).
- Requirement Engineering — NLI for RE tasks (arXiv 2024): keep: NLI entailment as coverage-verification method (weak).
- Requirement Engineering — others (GenAI assessments, legal compliance, multilingual crowd, prompt patterns, future of LLMs): reject: surveys/patterns, no artifact.
- Requirement Prioritization / Summarization: reject: different tasks.
- Requirement Traceability 1 — NLP for requirements traceability (arXiv 2024): keep: trace-link recovery survey (weak).
- Requirement Traceability 2 — BERT trace links (ICSE 2021): reject: link recovery between artifacts, not spans in text.
- Requirements QA 1 — LLMs for requirements QA (RE 2024): keep: LLM quality checks on requirements.
- Requirements QA 2 — satisfaction assessment: reject: satisfaction, not coverage.
- Requirements QA 3 — high→low level requirements coverage review (MSR 2024): keep: direct precedent for "what no requirement covers".
- Requirements QA 4 — ChatGPT user-story quality: reject: eval only.
- Requirement Classification (SVM/LSTM/ChatGPT experiment, XRC, NoRBERT, NFR, PRCBERT, greening): reject: FR/NFR labels, not document-part classes.
- Specification Generation (SpecGen et al.) / Specifications Repair / Use Case Generation / Software Modeling / Class Diagram / GUI Layouts: reject: formal program specs and models, not problem statements.
- Empirical 25 — ChatGPT prompt patterns incl. requirements elicitation: keep: prompt-pattern source for elicitation (weak).
- Software tool configuration — "Can LLMs Configure Software Tools": reject: tool config, not project intake.
- All development/testing/maintenance/management/benchmark/tuning sections: no line — outside this part.

## The search file

`sources/gh-search.tsv`, 20 queries, 129 data rows, read in full. Lines only for entries that touch the problem.

- `problem statement gap detection` (3 rows): all off-topic (nutrition, labour market, Copilot card) — no lines.
- `requirements ambiguity detection` (10 rows):
  - shujauzzaman/requirements-ambiguity-detection: reject: thin demo web app.
  - muhidabid/SoftReq: reject: notebook prototype, no stated method.
  - arsalaniltaf/Anaphoric-Ambiguity-Detection: keep: anaphora detection implementation (weak).
  - Yousif-FJ/RE-Ambiguity-Detection-Experiment: reject: comparison report only.
  - SNTSVV/taphsir: keep: anaphoric ambiguity resolution in requirements (AGPL-3.0, 2022).
  - profzivkovic/software-requirements-ambiguity-detection-replication: keep: replication package of ambiguity-detection research (weak).
  - CodewithSamiksha-dev/SmartBRD / bse223207: reject: student apps.
  - MetricBlood/srs-ambiguity-detection: keep: rule-based vs spaCy comparison — method (weak).
  - isti-fmt-nlp/AmbiguityToVariability: reject: variability mining dataset.
- `ambiguity detection natural language requirements LLM` (1 row): off-topic — no lines.
- `requirement smells detection` (5 rows):
  - arsalaniltaf/RequirementSmells: keep: multi-label smell-detection framework.
  - LucianoUCV/requirements-analyzer: reject: bachelor thesis notebook.
  - vshidlovsky/prd-agents-framework: keep.
  - Mikun07 ReqSmell backend+frontend: reject: LLM prototype, undocumented.
- `requirements quality checker` (10 rows):
  - eduardenoiu/NALABS: keep.
  - damianhorna / Gho-Ost / czajkowskis scenario-quality-checker: reject: scenario metrics only.
  - SuedaSen / alinaneamtu / maryamhussain28 / Rowell2002 / haydarozden: reject: stubs or LLM notebooks.
  - contextkits/quality-compliance-checker: reject: FDA QSR domain MCP.
- `clarifying questions LLM agent` (10 rows):
  - Tranquility2/quizzing: keep: one capped batch of clarifying questions → JSON answers.
  - anetrebskii/tuplet: reject: agent framework.
  - Rahul20037237/AMBI-agent: reject: when-to-ask MDP, not batching.
  - sanshikasri / siri-clone / Re-Coach / smart-job-match / BurnyCoder / Godot sidebar / kaynquang: reject: assistants, not intake.
- `requirements elicitation LLM` (10 rows):
  - robbrobb/LAT4MRE: keep: LLM-assisted elicitation tool (weak).
  - thomaslavery/requirements_elicitation: keep: LLM fills custom forms — batched questions precedent (weak).
  - eisedevries/Reqify: keep: elicitation + matching + verification pipeline (weak).
  - Ohara124c41/agentic-system-architect: keep: Five-Whys elicitation + deterministic evaluation.
  - VarshiniNeralla/ReqArchitect: keep: elicitation→validation→spec pipeline (weak).
  - jomarthomas NFR / Multi-Stage pipeline / hamdan sketch / requify-crowdsourcing / llm_design_tool: reject: narrow or undocumented.
- `requirements engineering LLM` (10 rows):
  - Agile-V/agile_v_skills: keep.
  - ybeven/4D-ARE: reject: attribution framework, unclear span model.
  - boxdate/ReqPipeline: reject: USDM-specific, C#.
  - garycli / prompt-blueprint / ai-auto-work / LLM4QDARE / srs-with-llm / DAPO prompt-engineering / Software-Engineering-with-LLMs: reject: generation pipelines, no intake checking.
- `requirements extraction from text LLM` (1 row): Robert-Carlton/Signals-Extraction: reject: job-posting extraction.
- `requirements traceability` (10 rows):
  - doorstop-dev/doorstop: keep.
  - awslabs/duvet: keep.
  - bmw-software-engineering/lobster: keep: traceability evidence reporting (weak — code/test links).
  - goeb/reqflow: keep.
  - tobhey/finegrained-traceability: reject: embedding link recovery.
  - ociu/sphinx-traceability-extension: keep: requirement items in docs with link matrices (weak).
  - SAP quality-continuous-traceability-monitor: reject: archived; test linkage.
  - rtmx-ai/rtmx: keep.
  - adrianpuiu/specification-document-generator: keep: six-phase pipeline with mandatory source verification (weak).
  - Ex-Origin/trace-syscall: reject: unrelated.
- `requirements traceability matrix generator` (7 rows): AdamVengroff, Rekkix, llm-test-case-generator, AI-Test-Plan-Agent, test-case-generator-ai, GenAI-RTM: reject: test-plan/Excel RTM generators. Falak359: reject: LLM-to-Excel RTM.
- `spec driven development` (10 rows):
  - Fission-AI/OpenSpec: keep.
  - Engineering4AI/awesome-spec-driven-development: keep: the list itself, mined above.
  - sc-spec-driven-development-files / hello-sdd / marcelsud: reject: courses/articles.
  - get-shit-done / agent-os / zhu1090093659 / mcp-server-sdd / spec-workflow-mcp: reject: prompt/tooling kits without intake gap detection.
- `PRD generator LLM` (10 rows):
  - sankalpdusane/ai-prd-generator-eval-critic: keep: structural SHIP/NO-SHIP gate, critic scores in code, line-level traceable revisions (weak).
  - laymom/ai-prd-generator: keep: idea→PRD through 4 human-gated stages (weak).
  - sonukhobragade/prd-figma-test-generator: reject: reads spec to emit tests, not intake.
  - AungMyoKyaw / cdeust / codemakerss / smalek0 / HR-AR / sumitsuman eval / ranjittilekar: reject: generation only.
- `acceptance criteria generation LLM` (5 rows):
  - SondosAlgyed/qa-agent: keep: reviews stories for missing/ambiguous/untestable acceptance criteria (weak).
  - arunvelmurugan-jpg/story-generator-agent: reject: story generation.
  - farhatparvin/amv-documentation-agents: keep: computed acceptance criteria + recorded traceability (weak, pharma).
  - fanyuuwang / treasury copilot: reject: generation/RAG.
- `standoff annotation format` (8 rows):
  - GateNLP/gateplugin-Format_Brat: keep: reference I/O implementation of the brat format.
  - harttu/annotation_consistency_checker: keep: checks standoff annotation consistency.
  - hatzel/salsa-standoff / brat2ocra / flatten-tei / eaf-to-paula / DiffCompare / brat-kappa: reject: converters/comparators, not the needed format.
- `requirements completeness check` (10 rows):
  - rabbittsui/ModelCardDiff: keep: schema completeness gate.
  - yon-tyzes/requirements-intake-agent: keep: intake gap checks + numbered questions + provenance tags.
  - lrakai / ML-Ops checklist / graduation checkers / hipaa / luhn / Fabi-P: reject: unrelated checklists.
- `awesome requirements engineering` (4 rows):
  - muchengracedriver/awesome-requirements-writer: keep: mined above.
  - jgsystemsconsulting/awesome-requirements-engineering: keep: further curated list — standards, formats, tools (data source, weak).
  - ishandutta2007/Awesome-Engineering-Requirements-Management: reject: thin link dump.
  - Dinovative job posting: reject: unrelated.
- `awesome LLM software engineering` (5 rows):
  - iSEngLab/AwesomeLLM4SE: keep: the list itself, mined above.
  - codefuse-ai/Awesome-Code-LLM / Awesome-Issue-Resolution / awesome-agentic-engineering / awesome-calculators: reject: other domains.

## What no candidate covers

- **The union in one system.** Every piece exists separately: ambiguity/smell
  detection (mcp-reqcheck, NALABS, QuARS, spec-lint), span-anchored facts (brat),
  requirement records with trace and acceptance criterion (Volere, doorstop,
  duvet, rtmx), batched questions (quizzing, spec-kit clarify), evidence
  enforcement (squelette, Agile-V). No public solution composes them into one
  intake pass over a problem statement.
- **Byte-exact gap quotes verified by code.** Agents must quote the file and
  code must verify it. duvet and brat verify stored span text; mcp-reqcheck
  flags terms deterministically — but nothing couples an LLM-produced finding
  to a code-checked quote of a prose document.
- **The tiling check.** "Every non-blank part of the text is in exactly one
  span" is a partition/tiling verification. Span *validity* checkers exist
  (bratpy bounds+text match, anno-eval bounds/overlap/text match,
  harttu consistency checker); none verifies full-document coverage with
  exactly-one-span semantics.
- **Fact classes need/constraint/measure/context.** Classifiers exist for
  requirement-vs-prose (ORCS, demarcator) and FR/NFR (many), none for this
  four-class span taxonomy.
- **Requirement→source-text trace.** Traceability tools link
  requirement→requirement, →code, →test (doorstop, rtmx, reqflow, lobster,
  sphinx-traceability). Only duvet links to *document text* (and then
  requirement→code). Nothing links a derived requirement to numbered fact
  spans.
- **One merged, deduplicated question batch with answer tracking.** quizzing
  and spec-kit produce one question batch from one source each. No candidate
  merges questions from rule checks, meaning review, completeness check, and
  setup into one ID'd, deduplicated batch, checks answers for completeness,
  and re-asks only the missing ones.
- **Answers stored in Git as raw facts.** git notes (git-adr) and ADR files
  store decisions; Agile-V stores evidence records. Nothing stores question-
  answer pairs keyed to gap IDs as immutable Git facts.
- **"No configuration value without evidence" as a mechanical rule.**
  intake-agent's [SOURCE]/[INFERRED] tags and Agile-V/Squelette evidence
  schemas approach it; no tool enforces citation-per-value by code.
- **Independent second-agent coverage review of the *source text*.**
  spec-kit analyze and Agile-V run independent checks of derived artifacts;
  coverage tools check requirement→code. Nothing has an independent agent
  list the *facts* no requirement covers.
- **Go-stdlib implementation.** Nothing usable is written in Go against the
  standard library; rtmx (Go) shows the domain is expressible there, but it
  is a dependency candidate only in shape, not usable under the constraint.

## Limits

- GitHub API rate limit (unauthenticated, 60/h) hit during metadata gathering;
  some last-activity dates come from the provided TSV `pushed` column or repo
  pages rather than fresh API data.
- Licenses not visible in fetched pages are marked "not stated": mcp-reqcheck
  (PyPI page blocked in part), pi-harness-runtime, harttu checker, MADR,
  goeb/reqflow states GPL in COPYING without a version number I could confirm.
- MySpec is a proprietary SaaS — judged from marketing pages only.
- Academic artifacts (TABASCO, MAANA, Ezzini anaphora, ORCS, QuARS, Smella)
  were judged from papers/READMEs; none were executed. Whether their models
  still download is unchecked.
- The Volere template and ISO 29148 are paywalled; judged from public
  articles, previews, and secondary descriptions.
- The other searcher's list coverage may differ; overlapping candidates were
  not deduplicated against their work.
~~~~
