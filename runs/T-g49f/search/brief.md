# The brief of the public-solution search of M2c (task T-g49f)

Date: 2026-10-10. You are one of two searchers. The other searcher runs on
another model and another harness, and neither of you sees the other's work.
This is a search, not a design: you find and judge existing public solutions;
you do not design one.

## The problem

A project's software is delivered by role agents, from a problem statement that
an idea owner writes in prose. Before delivery starts, a controller program
must turn that problem statement into checked input. The part of the project
that this search serves needs:

- the gaps of meaning of a problem statement found before delivery: a term
  with two readings, a metric that cannot be measured as written, a missing
  intent value, a risky assumption or open decision; each with a byte-exact
  quote of the text it concerns, and a question; the quote checked against the
  file by code;
- the problem statement split into numbered facts by their exact positions in
  the file, each classed (for example a need, a constraint, a measure, or
  context), with a mechanical check that every non-blank part of the text is
  in exactly one span and each span's text equals the file's;
- identified requirements derived from those facts, each with a trace to the
  facts that state it and an acceptance criterion; a second, independent agent
  then lists what no requirement covers;
- the questions of every source (rule checks, the review of meaning, the
  completeness check, setup questions) merged into one batch with an ID each,
  duplicates removed, for the idea owner and the operator; the answers checked
  for completeness, a follow-up only for the missing or open ones, and the
  answers stored in the project's Git as a raw fact;
- no configuration value without evidence.

These are the facts of the problem statement that this part serves, word for
word. A fact that is a table row keeps its table form:

- F-0003#41. **Problem Statement Quality:** Detection of the gaps in a problem statement before delivery starts, so that the questions for the idea owner come in one batch and not one at a time.
- F-0003#14. **Symptom:** A project starts from a problem statement that has gaps: metrics with no measurement method, terms with two readings, no named technology stack. Each gap becomes a blocking question during delivery. In parallel, the discipline system of the new project is copied and configured by hand.
- F-0003#15. **Root Cause:** Nothing checks a problem statement for gaps before delivery starts. People adapt the discipline template by hand to the domain and technology stack of each project. No record shows the evidence for each value.
- F-0003#51. **Specification Synthesis:** Derivation of identified requirements and technical specifications from the approved problem statement. Each requirement has a trace to the text that states it, and an acceptance criterion.
- F-0003#23. **Symptom:** After the approval of the problem statement, humans write the requirement documents and the technical specifications by hand. Delivery waits for these documents.
- F-0003#24. **Root Cause:** A problem statement is prose. Nothing converts it into identified requirements that a machine can validate and trace back to the text that states them.
- F-0003#25. **Downstream Friction:** Each step states a requirement again in different words, so the intent degrades. A delivered result cannot be traced to a need, so the idea owner cannot accept the work requirement by requirement.
- F-0003#62. | **Specification Traceability** | Manual, not measured. | Every delivered requirement has an identifier, a trace to the text of the problem statement, and an acceptance criterion. | Requirements with a complete trace, divided by all delivered requirements. |
- F-0003#63. | **Setup Correctness** (Inv. 4) | Manual setup, not checked. | $100\%$ of new project repositories pass the Armature discipline tests. 0 configuration values without evidence. | Run of the discipline tests of the baseline. Count of values with no citation. |
- F-0001#4. **No configuration value without evidence.** A value that comes from a guess is a defect.
- F-0001#10. **Intent decisions.** Before delivery starts, the idea owner sets the problem, the success criteria, and the funding.
- F-0001#11. **Problem statement answers.** Before delivery starts, the idea owner answers the gap questions about the problem statement, in one batch.
- F-0001#39. | **Specification** | The requirement and design documents that a project derives from the approved problem statement. Each requirement in them has an identifier, a trace to the problem statement, and an acceptance criterion. |

The idea owner's answer that bounds the solution (F-0004#1): the controller is
written in **Go, with the standard library only, and Git called as the `git`
program**. So a library in another language, or a Go module, cannot be a
dependency; a candidate can still give a pattern, a method, a file format to
follow, a data source, a taxonomy (of ambiguity, of requirement smells), or a
separate tool. The agents that read meaning are coding-agent programs run as
sessions; the controller itself calls no model.

The milestone that this search serves states its first demo so: "On a pilot
problem statement, one Intake comment holds the rule gaps and the gaps of
meaning, and `layup spec check --facts` passes on the numbered spans."

## Your sources

- `sources/lists/`: five curated lists, each a file of its repository at the
  commit in `sources/lists.tsv` (repository, commit, path, lines, links).
  **Read each list in full.** Do not summarize a list with a filter or a cap.
- `sources/gh-search.tsv`: the results of 20 GitHub repository searches, each
  query with the fact words it comes from. Read it in full.
- The web. Run your own searches in the words of the problem above, as well as
  any words of your own. Record each query and the words it comes from.

## What to do

1. Search: your own queries first, then the lists and the search file.
2. For each entry of a list or of the search file that touches the problem,
   give **keep** or **reject** and a reason of a few words. An entry that does
   not touch the problem needs no line.
3. For each candidate that you keep, check its own repository, documentation
   or paper: its license, its language, the date of its last commit or release
   (or its year, for a paper), which part of the problem it solves, where it
   keeps its state (in the project's Git or not), whether it is deterministic,
   and its documentation. Write "not stated" where its source does not say.
4. Say what no candidate covers.

Judge only by what a source says. Run no code of a candidate, install nothing,
and write nothing outside `out/` in this directory. Read no file on this
computer outside this directory.

## Your output

Write one Markdown file, `out/search.md`, with these sections, about 300 to 450
lines in all:

1. `## Queries`: a table, one row per query: the query, the words it comes from,
   what it found.
2. `## Candidates`: a table of the kept candidates: name and URL, license,
   language, last activity, the part of the problem, state in Git,
   deterministic, keep as (a dependency, a pattern, a method, a format, a
   taxonomy, a data source, a tool) and why.
3. `## The lists`: per list, one line per entry that touches the problem:
   `name: keep|reject: reason`.
4. `## The search file`: the same, for `sources/gh-search.tsv`.
5. `## What no candidate covers`.
6. `## Limits`: what you could not check.

When the file is complete, end your session.
