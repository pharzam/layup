# The brief of the public-solution search of M2d (task T-msrd)

Date: 2026-10-11. You are one of two searchers. The other searcher runs on
another model and another harness, and neither of you sees the other's work.
This is a search, not a design: you find and judge existing public solutions;
you do not design one.

## The problem

A project's software is delivered by role agents. A controller program sets up
each new project repository from a pinned discipline baseline, then protects it
from those agents. The part of the project that this search serves needs:

- the setup of a new repository from a pinned baseline commit, driven by a
  program: each configuration value from an answer, a catalog entry or a fact,
  with its source recorded; a value with no source left as an open gap; the
  setup stopping at each human decision and resuming from its own record;
  free-text parts written by an agent and each checked by code;
- the setup's evidence and answers kept in the project's Git, on a branch that
  only the controller's identity may update;
- the branch rules of the forge (GitHub) set by a human with admin rights from
  a file that the program writes: required checks pinned to the app that posts
  them, pull requests required, no force push, no bypass for the controller;
- the proof from outside that the rules hold: the effective rules read back
  through the forge's API, a probe push by the controller's own identity that
  the forge must refuse for a rule violation (any other failure is not a
  pass), and a read of the repository's activity log to audit who updated
  each branch;
- the repository passing its own gates when the controller is removed;
- a second run of the setup on a repository that an earlier version set up.

These are the facts of the problem statement that this part serves, word for
word. A fact that is a table row keeps its table form:

- F-0003#42. **Reproducible Discipline Setup:** A repeatable, evidence-based setup of the Armature baseline for a new project, for the domain and technology stack of that project.
- F-0003#15. **Root Cause:** Nothing checks a problem statement for gaps before delivery starts. People adapt the discipline template by hand to the domain and technology stack of each project. No record shows the evidence for each value.
- F-0003#37. **Non-Reproducible Project Starts:** Each new project has a different, hand-made discipline setup, so results from one project do not transfer to the next.
- F-0003#43. **Rule Protection:** Protection of the rules from the agents that the rules govern.
- F-0003#63. | **Setup Correctness** (Inv. 4) | Manual setup, not checked. | $100\%$ of new project repositories pass the Armature discipline tests. 0 configuration values without evidence. | Run of the discipline tests of the baseline. Count of values with no citation. |
- F-0003#64. | **Gate Integrity** (Inv. 3) | Not measured. | $100\%$ detection of known-bad commits. 0 agent writes to rule paths. | A set of known-bad commits against the gates. Audit of the writes to rule paths. |
- F-0003#65. | **Independence** (Inv. 2) | Not applicable. | The project repository passes its gates without the automation. | Gate run with the automation removed. |
- F-0001#1. **Git is the system of record.** No project state and no decision is kept only outside the project repository.
- F-0001#2. **The project repository is independent.** A human or a different agent can continue the work without the automation that created the repository.
- F-0001#3. **The agents that do the work cannot change the rules or the gates that check the work.**
- F-0001#4. **No configuration value without evidence.** A value that comes from a guess is a defect.
- F-0001#8. **Armature is used at a pinned, recorded version.**

The idea owner's answer that bounds the solution (F-0004#1): the controller is
written in **Go, with the standard library only, and Git called as the `git`
program**. So a library in another language, or a Go module, cannot be a
dependency; a candidate can still give a pattern, a method, a file format to
follow, an API usage, a separate tool, or a known limit of the forge. The
controller holds a GitHub App's token (contents, issues, pull requests and
commit statuses; no administration and no workflows permission).

The milestone that this search serves states its first demo so: "After the
Operator's pushes and ruleset apply, the push to `layup-probe` is refused, and
`setup/forge-check.tsv` shows each rule and probe as passed."

## Your sources

- `sources/lists/`: the curated lists, each a file of its repository at the
  commit in `sources/lists.tsv` (repository, commit, path, lines, links).
  **Read each list in full.** Do not summarize a list with a filter or a cap.
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
