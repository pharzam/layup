# The brief of the public-solution search of M2b (task T-ywk7)

Date: 2026-10-09. You are one of two searchers. The other searcher runs on
another model and another harness, and neither of you sees the other's work.
This is a search, not a design: you find and judge existing public solutions;
you do not design one.

## The problem

A project's software is delivered by role agents. Each role agent runs on a
**harness agent**: a coding-agent program such as a terminal coding agent, run
as a session. A controller program starts each session, gives it its work, and
reads its result. The problem below is the part of the project that needs:

- a session of a harness agent that the controller starts, isolates, limits in
  time and money, stops, and reads the result of;
- a first check of each harness agent that the project registers (its version,
  and what it can do) before it gets work;
- the choice of a harness agent and a model for each piece of work;
- a record of each session's token count, latency, wall-clock duration and
  money, taken from what the harness agent reports;
- the project's rules in a form that no one harness agent owns;
- no change by an agent to the paths of the rules and the gates;
- a form for the information that passes between role agents that a machine
  can validate.

These are the facts of the problem statement that this part serves, word for
word. A fact that is a table row keeps its table form: F-0001#38 is a row of the
glossary (term, definition), and F-0003#59 to #66 are rows of the table of
invariant checks, whose columns are "Check | Today | Requirement | Verification":

- F-0001#1. **Git is the system of record.** No project state and no decision is kept only outside the project repository.
- F-0001#3. **The agents that do the work cannot change the rules or the gates that check the work.**
- F-0001#6. **A deterministic check is preferred to an LLM judgement** where a rule can be checked mechanically (Armature R5).
- F-0001#9. **A harness agent is replaceable.** The rules, the context, and the state of a project do not depend on one harness agent.
- F-0001#38. | **Telemetry** | The record of the token count, the latency, and the wall-clock duration of a task. |
- F-0003#20. **Symptom:** A task consumes an unknown number of tokens, an unknown amount of money, and an unknown quantity of time. Nobody can say which role, task, or retry loop consumes the budget.
- F-0003#21. **Root Cause:** No record holds the token count, the latency, and the wall-clock duration of each task. Cost is not part of the acceptance of a result.
- F-0003#22. **Downstream Friction:** The budget is a business-forking decision (§6), but the idea owner has no number for it. A cheap path and an expensive path cannot be compared. A retry loop can consume a budget before a human sees it.
- F-0003#26. **Symptom:** Work crosses more than one harness agent. Each one has its own rule file, memory, context format, and limits. The project rules are copied into each format, and each copy drifts.
- F-0003#27. **Root Cause:** Rules, context, and task state are kept in the form of one harness agent, and not in one neutral form in the project repository. No procedure gives the same rules to a different harness agent.
- F-0003#28. **Downstream Friction:** An independent harness agent cannot check the result of the first one, so a verification stays inside the product that made the error. The project also loses its independence (Invariant 2).
- F-0003#43. **Rule Protection:** Protection of the rules from the agents that the rules govern.
- F-0003#45. **Role Handoffs:** Information that passes between role agents has a form that a machine can validate, based on Armature conventions.
- F-0003#50. **Cost Visibility:** Each task gets a record of its token count, its latency, and its wall-clock duration in the project repository.
- F-0003#52. **Harness-Agent Neutrality:** The rules, the context, and the task state stay in a form that does not belong to one harness agent. A second harness agent can do the work, and can check the work of the first one.
- F-0003#59. | **Inter-Role Communication Format** | Unstructured chat prompts and markdown notes. | $100\%$ schema-validated state artifacts across all role transitions. | Validated role transitions, divided by all role transitions. |
- F-0003#60. | **Telemetry Completeness** (§6 Cost Visibility) | Absent. | Every task has a token count, a latency, and a wall-clock duration in the project repository. | Tasks with a complete telemetry record, divided by all tasks. |
- F-0003#64. | **Gate Integrity** (Inv. 3) | Not measured. | $100\%$ detection of known-bad commits. 0 agent writes to rule paths. | A set of known-bad commits against the gates. Audit of the writes to rule paths. |
- F-0003#66. | **Harness-Agent Neutrality** (Inv. 9) | Not applicable. | The same project rules and gates run under $\ge 2$ harness agents. Each change gets $\ge 1$ verification from a harness agent that did not make the change. | Count of harness agents that pass the gate run. Changes with an independent verification, divided by all changes. |

The idea owner's answer that bounds the solution (F-0004#1): the controller is
written in **Go, with the standard library only, and Git called as the `git`
program**. So a library in another language, or a Go module, cannot be a
dependency; a candidate can still give a pattern, a protocol or file format to
follow, a data source (for example a price list), or a separate tool.

The milestone that this search serves states its first demo so: "A probe
session on each registered harness, then one developer session whose commit
lands on its task branch, with one complete telemetry row."

## Your sources

- `sources/lists/`: five curated lists, each the README of its repository at the
  commit in `sources/lists.tsv` (repository, commit, path, lines, links). **Read
  each list in full.** Do not summarize a list with a filter or a cap.
- `sources/gh-search.tsv`: 200 results of 20 GitHub repository searches, each
  query with the fact or the milestone words it comes from. Read it in full.
- The web. Run your own searches in the words of the problem above, as well as
  any words of your own. Record each query and the words of the facts it comes
  from.

## What to do

1. Search: your own queries first, then the lists and the search file.
2. For each entry of a list or of the search file that touches the problem,
   give **keep** or **reject** and a reason of a few words. An entry that does
   not touch the problem needs no line.
3. For each candidate that you keep, check its own repository or documentation:
   its license, its language, the date of its last commit or release, which
   part of the problem it solves, where it keeps its state (in the project's
   Git or not), whether it is deterministic, and its documentation. Write "not
   stated" where its source does not say.
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
   language, last activity, the part of the problem, state in Git, deterministic,
   keep as (a dependency, a pattern, a format, a data source, a tool) and why.
3. `## The lists`: per list, one line per entry that touches the problem:
   `name: keep|reject: reason`.
4. `## The search file`: the same, for `sources/gh-search.tsv`.
5. `## What no candidate covers`.
6. `## Limits`: what you could not check.

When the file is complete, end your session.
