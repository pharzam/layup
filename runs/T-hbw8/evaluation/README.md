# T-hbw8 — hands-on evaluation of the reuse candidates (O-86, step 1)

The Operator's decision O-86 ([`../operator-decisions.md`](../operator-decisions.md)):
before the rewrite plan, each reuse candidate gets a thorough, hands-on
evaluation with no time bound, and ends with **use**, **borrow the pattern** or
**reject**, with the evidence. This directory is that record. The result is
[`summary.md`](summary.md).

## The candidates

The candidate set is the **union of the two searchers' shortlists**
([`../search-v2/searcher-A.md`](../search-v2/searcher-A.md) candidates A to J,
[`../search-v2/searcher-B.md`](../search-v2/searcher-B.md) §3.1 to §3.8 and its
verification pair). The table of [`../search-v2/summary.md`](../search-v2/summary.md)
says it is that union, but it is not: it leaves out searcher A's candidates B
(Agent Orchestrator, Go, about 12,500 stars, "closest to a Go multi-harness
project control plane"), C (Omnigent), E (Springer), F (sdlc-gh) and G (Agent
Context Store), and searcher B's §3.8 (AgentPlane) and no_human. The handoff
list on #72 copied that table. This evaluation adds them, so that the defect of
the first search (a list reduced before it was read) does not repeat.

| No. | Candidate | Found by | Group |
| --- | --------- | -------- | ----- |
| 01 | Gas Town | B §3.5 | orchestrator, Go |
| 02 | Beads | B §3.5, A list | work tracking, Go |
| 03 | Paperclip | A H, B §3.1 | orchestrator with budgets and approvals |
| 04 | Agent Orchestrator | A B | orchestrator, Go |
| 05 | Omnigent | A C | cross-harness work |
| 06 | Spec Kitty | B §3.4 | specification workflow |
| 07 | BMAD-METHOD | B §3.2 | specification and role skills |
| 08 | MetaGPT | A I, B §3.3 | role SOPs and documents |
| 09 | Springer | A E | role artifacts and human checkpoints |
| 10 | GNAP | A D, B §3.6 | records in Git |
| 11 | AgentPlane | B §3.8 | records in Git |
| 12 | Agent Context Store | A G | handoff records |
| 13 | AI-SDLC | A A | governed delivery, cross-harness review |
| 14 | sdlc-gh | A F | stack-selected CI and rule protection |
| 15 | AgentJury | A J | blind review council |
| 16 | Loki Mode | B §3.7 | completion gates and blind council |
| 17 | isitdone | B §3.8 | deterministic done-check |
| 18 | no_human | B §3.8 | second-model review |
| 19 | Shape Up | the Operator, O-90 | a method for milestones, budgets and progress (a book, not software) |

**Not evaluated hands-on:** the other entries that the searchers marked `keep`
in their full-list passes (about 60, each with a one-line reason in the two
searcher files). They were not shortlisted by either searcher. Where the result
of this evaluation leaves a part with no candidate, `summary.md` names the kept
entries that could fill it.

## The method

- **Brief.** Each evaluator got the same brief (below) and a copy of the inputs
  in a scratch directory: the problem statement, the vision brief, `F-0004`,
  ADR-0010, the Operator's decisions, the three deep-check files and the two
  searcher files. No evaluator read or wrote the LAYUP checkout.
- **Hands-on.** Each candidate was cloned at a pinned commit into a disposable
  directory, its code read, installed and run on a small Go test target with a
  local bare remote. Where a candidate drives a coding agent, Claude Code ran in
  print mode with a small model and a spend cap. No evaluator wrote to GitHub or
  any other network service. Each evaluator stopped every process it started.
- **Criteria.** The twelve In-Scope items S1 to S12 and the nine invariants I1
  to I9 of the problem statement (§6); the deep-check findings
  (`Sol-N`, `Fable-MN`, `Author-N`) that a mechanism of the candidate answers;
  the fit with LAYUP's constraints (Go standard library only, `F-0004` fact 1;
  an external binary with no dependency in the target, O-76).
- **Verdicts.** `use` (run it as a separate program next to LAYUP, or depend on
  its file format or protocol, or vendor code under its license); `borrow the
  pattern` (LAYUP writes its own code after the design of the candidate);
  `reject`. Each part of a candidate gets a verdict, and each candidate one
  overall verdict.
- **Evaluators.** Seven Claude Opus 5.5 sessions (the author's model) on Claude
  Code, one per group of candidates, started by the author's session; each
  file names its evaluator, its date and its elapsed time. The author read each
  file, checked it against the brief, and wrote `summary.md`.

The brief, word for word: [`brief.md`](brief.md).
