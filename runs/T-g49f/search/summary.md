# T-g49f — the public-solution search of M2c: the result

Two blind searchers screened the same sources ([`sources.md`](sources.md)) and ran
their own web searches in the words of the problem:
[`searcher-A.md`](searcher-A.md), SWE-2 High on Devin, and
[`searcher-B.md`](searcher-B.md), Claude Sonnet 5.5 on Claude Code. Their brief
([`brief.md`](brief.md)) held the facts, the milestone row and the sources: no
name of a candidate, nothing of the architecture or the specification. This file
is the author's reading of the two files, written after both ended, and
committed before the decisions of S3.

**Both conclude: no public solution can be the base of `M2c`.** No candidate can
be a dependency of a controller in Go with the standard library only
(`F-0004#1`); each is a program, a pattern, a format, a method or a taxonomy. A:
"No public solution composes them into one intake pass over a problem
statement." B: "No candidate can be a dependency under the rule."

## What the search gives each part of M2c

| Part | Candidates (searcher) | What they bear on in `M2c` (the decisions come in S3) |
| ---- | --------------------- | ----------------------------------------------------- |
| The facts as spans of the file | the brat standoff format (A); the W3C Web Annotation Model, `TextPositionSelector` and `TextQuoteSelector` (B); Hypothesis's anchoring order (B) | A span is a start and an end offset with its text stored beside it, so code can check the text against the file; the W3C model gives both a position and a quote for one target. Neither has the tiling rule: span checkers check bounds, overlap and text (A: bratpy, anno-eval, the consistency checker of harttu), and none checks that each non-blank byte is in exactly one span (A, B) |
| The classes of a fact | ORCS, demarcator (A) | They class a document's parts as requirement or prose, or as functional and non-functional; no candidate has the four classes `need`, `constraint`, `measure`, `context` (A) |
| The gaps of meaning and their kinds | Smella and the smells of ISO/IEC/IEEE 29148 (A, B); the ambiguity handbook of Berry, Kamsties and Krieger (A, B); NALABS (A, B); QuARS (A); TABASCO, the anaphora work of Ezzini et al., MAANA (A); the 13 categories of the SEAMS ambiguity scanner (A); a RAID log (A) | Taxonomies of ambiguity and of requirement smells for sentences; "a metric that cannot be measured as written" is closest to the smell "non-verifiable term" (B). No public taxonomy names a missing intent value, a risky assumption or an open decision of a problem statement (B); a RAID log names assumptions and decisions as register classes (A) |
| The quote of a gap, checked by code | Duvet (A, B); `mcp-reqcheck` and its findings with a question each (A); the gap record `{id, question, evidence, category}` of `pi-harness-runtime` (A) | Duvet checks a quoted requirement against the text of a specification; no candidate checks a model's quote against the file byte for byte as the gate of a finding (A, B) |
| The requirements, their trace and acceptance criterion | Doorstop, RTMX (Go, with two outside modules), StrictDoc, Lobster, Kibi, reqflow, flashtrace (A, B); Volere's fit criterion, EARS, Gherkin (A, B); OpenSpec, spec-kit (A, B) | Identified items in files in Git, with links and an integrity check; each traces a requirement to another, to code or to a test, not to numbered facts of a problem statement (A, B); flashtrace's verdicts (missing, orphan, duplicate) are a shape for a coverage result (A) |
| The completeness review by a second agent | spec-kit `/analyze`, Agile-V's verifier that never checks its own work (A, B); the completeness work of Luitel et al. (B) | An independent check of derived artifacts exists; none has a second agent list the facts of the source text that no requirement covers (A, B) |
| One batch of questions, the answers, the follow-up | Quizzing (A, B); spec-kit `/clarify` (A, B); requirements-intake-agent (A, B); Ouroboros (B) | Quizzing puts all questions of one source in one form and builds a follow-up from the open items; spec-kit caps and asks in turn; no candidate merges the questions of several sources with IDs, removes repeats, and asks again only for the missing or open answers (A, B) |
| The answers in Git as a raw fact | git-adr (decisions in git notes), MADR (A); Squelette's decisions file, spec-kit's clarification section (B) | Decisions and answers recorded in the repository; none stores each answer, unedited, as a fact keyed to its question ID (A, B) |
| No value without evidence | requirements-intake-agent's tags `[SOURCE]`, `[INFERRED]`, `[MISSING]` (A, B); Squelette, Agile-V's evidence schemas (A, B) | Evidence per value as an agent's tag or a work item's rule; no tool checks by code that each configuration value cites its source (A, B) |
| A corpus of problem statements | the PURE dataset (A) | 79 public requirements documents: a source of trial input, not of the pilot's problem statement, which the idea owner chooses |

## The limits of the search

- Searcher A ran on Devin's free model, an execution-tier model on a reasoning
  part, as Devin's paid quota was exhausted; searcher B on Claude Sonnet 5.5,
  as OpenCode's login failed. The reasoning tier of the author's own family ran
  no search.
- Searcher B's run took 3 min 26 s and A's 12 min 13 s; each says it read every
  list in full (B: AwesomeLLM4SE "read in full, but its entries are paper
  titles"; A: "1711 lines, read in full"). A's failed first try had "skimmed"
  the middle of that list; its second try's log says it read the rest "to
  fulfill the 'read in full' requirement". A file's claim of a full read is the
  searcher's own.
- Both searchers judged by pages, READMEs and papers; neither ran a candidate.
  Dates, licenses and languages are partly "not stated" or inferred (each
  file's Limits).
