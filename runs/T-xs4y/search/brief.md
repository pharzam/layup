# The brief of the public-solution search of M2e (task T-xs4y)

Date: 2026-10-15. You are one of two searchers. The other searcher runs on
another model and another harness, and neither of you sees the other's work.
This is a search, not a design: you find and judge existing public solutions;
you do not design one.

## The problem

A project's software is delivered by role agents (an architect, developers,
a verifier, an approver), each a session of a language-model harness. A
controller program drives them. The part of the project that this search serves
needs:

- a task loop: each change is a task with an issue, a plan, a plan review by a
  different model than the plan's author, a developer session, a draft pull
  request, automated checks, a verifier on a different harness than the
  author's, a close-out commit, and the merge of the exact head commit that
  was checked;
- a typed handoff between roles: each piece of information that passes from one
  role to the next has a form that code validates; a transition that the
  table of allowed transitions does not list is invalid; what an agent writes
  never makes a handoff valid by itself; an artifact that exists only as an
  issue comment still has a place and a hash that code can check;
- the forge's commit statuses posted by the controller, one for each class of
  check (gates, specification, rules, verification), required by the branch
  rules; a pull request is not marked ready and does not merge while a check
  did not run or is not active;
- a filter of architecture options by code, so that an option with a failed or
  missing row, or a pass with no evidence, is dropped; one survivor is chosen;
- a brief to the idea owner, checked by code for its parts and length, and
  answered by one comment that the controller copies into the project's Git;
- identified requirements derived from a problem statement, each with a trace
  to the text that states it and an acceptance criterion, and a checker that
  fails on an unresolved ID, an uncovered need, a requirement with no
  criterion, or a delivered requirement with no specification section;
- the controller's state (task register, transitions, events, gate results,
  parameters) kept as records in the project's Git, on a branch that only the
  controller's identity may update.

These are the facts of the problem statement that this part serves, word for
word:

- F-0003#45. **Role Handoffs:** Information that passes between role agents has a form that a machine can validate, based on Armature conventions.
- F-0003#59. | **Inter-Role Communication Format** | Unstructured chat prompts and markdown notes. | $100\%$ schema-validated state artifacts across all role transitions. | Validated role transitions, divided by all role transitions. |
- F-0003#47. **Verification on Every Change:** Each change passes the gates for repository layout, interface boundaries, and testing pyramids before it reaches human review or merges. The gates are deterministic where a rule can be checked mechanically (Invariant 6).
- F-0003#58. | **Structural Conformance** (Inv. 5) | Variable; manually checked during PR review. | No PR with a failed gate reaches human review or merges. | Count of PRs with a failed gate that reach human review or merge. The first-attempt pass rate is monitored only, so that no one weakens a rule to get a better value. |
- F-0003#51. **Specification Synthesis:** Derivation of identified requirements and technical specifications from the approved problem statement. Each requirement has a trace to the text that states it, and an acceptance criterion.
- F-0003#23. **Symptom:** After the approval of the problem statement, humans write the requirement documents and the technical specifications by hand. Delivery waits for these documents.
- F-0003#24. **Root Cause:** A problem statement is prose. Nothing converts it into identified requirements that a machine can validate and trace back to the text that states them.
- F-0003#25. **Downstream Friction:** Each step states a requirement again in different words, so the intent degrades. A delivered result cannot be traced to a need, so the idea owner cannot accept the work requirement by requirement.
- F-0003#62. | **Specification Traceability** | Manual, not measured. | Every delivered requirement has an identifier, a trace to the text of the problem statement, and an acceptance criterion. | Requirements with a complete trace, divided by all delivered requirements. |
- F-0003#43. **Rule Protection:** Protection of the rules from the agents that the rules govern.
- F-0001#1. **Git is the system of record.** No project state and no decision is kept only outside the project repository.
- F-0001#3. **The agents that do the work cannot change the rules or the gates that check the work.**
- F-0001#5. **A check that is not active does not count as passed.**
- F-0001#39. | **Specification** | The requirement and design documents that a project derives from the approved problem statement. Each requirement in them has an identifier, a trace to the problem statement, and an acceptance criterion. |

The idea owner's answer that bounds the solution (F-0004#1): the controller is
written in **Go, with the standard library only, and Git called as the `git`
program**. So a library in another language, or a Go module, cannot be a
dependency; a candidate can still give a pattern, a method, a file format to
follow, an API usage, a separate tool, or a known limit of the forge. The
controller holds a GitHub App's token (contents, issues, pull requests and
commit statuses; no administration and no workflows permission).

The milestone that this search serves states its first demo so: "The
specification task of a target merges after its first bet: each handoff is
valid, the rendered commit holds the MoSCoW and Phase columns, and the four
`layup/` checks are green."

## Your sources

- `sources/lists/`: the curated lists, each a file of its repository at the
  commit in `sources/lists.tsv`. **Read each list in full.** Do not summarize a
  list with a filter or a cap.
- `sources/gh-search.tsv`: the results of the GitHub repository searches, each
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
   deterministic, keep as (a dependency, a pattern, a method, a format, an API
   usage, a tool, a known limit) and why.
3. `## The lists`: per list, one line per entry that touches the problem:
   `name: keep|reject: reason`.
4. `## The search file`: the same, for `sources/gh-search.tsv`.
5. `## What no candidate covers`.
6. `## Limits`: what you could not check.

When the file is complete, end your session.
