# T-hbw8 — root cause: a missed public solution and about 50 design defects

Asked by the Operator on 2026-09-28 after the deep check
([`deep-check-sol.md`](deep-check-sol.md), [`deep-check-fable.md`](deep-check-fable.md),
[`deep-check-author.md`](deep-check-author.md)): why did the author and the
harness agents miss a competitor such as Paperclip, and how is it fixed.
Author of this analysis: Claude Opus 5.5 (the author of the architecture).

## The evidence

The public-solution search of this task is section 1 of
[`selection.md`](selection.md). It fetched the list
`https://github.com/andyrewlee/awesome-agent-orchestrators` (about 247
entries). Paperclip is line 163 of that list: "Self-hosted platform where agents
wake on heartbeats to claim tickets, governed by org charts, budgets, and
approval gates." The author did not read the list. The author passed it through
a web-fetch summarizer with this instruction: "List the entries that
orchestrate several coding-agent CLIs (Claude Code, Codex, Gemini, etc.) with a
workflow of phases, cross-agent review, or model routing. … Max 15 entries." The
author saw 15 of about 247 entries; Paperclip was not among them.

## Root causes

1. **Anchoring.** The target design of O-67 (a Go orchestrator `layup run`, the
   phases, Jev) was the author's proposal of 2026-09-27, accepted before any
   search; ADR-0010 (Go) and ADR-0011 (a CLI with no service) pointed the same
   way. The search then confirmed a design; it did not look for one.
2. **The vocabulary of the solution, not of the problem.** The queries were
   "coding agents harnesses git worktrees", "durable workflow engine Go
   git-backed", "LLM model routing bandit". The PSB speaks of multi-role agent
   teams, the human as a message bus, budgets, approvals and stalls; Paperclip
   speaks of org charts, budgets, approval gates and heartbeats.
3. **A filtered and capped source.** The one list that held the answer was
   reduced by a summarizer with the author's filter and a cap of 15. What a
   filter removes, nobody sees.
4. **One searcher, and no check of the search.** One session searched in about
   five minutes. Plan-review condition 7 asked for a search but set no test of a
   sufficient one; the panel brief asked for options, not for existing
   products; the deep-check brief did not ask about public solutions.
5. **Coverage by name.** The coverage tables of `docs/architecture.md` matched
   each PSB item to a component name. That looks complete without proof that
   the component can work; the reviewers found the defects by walking one
   concrete case from end to end.
6. **No check of what a mechanism can do.** For "checks", "finds", "detects",
   nobody asked: what is the input, and can code, a model or only a human do
   it? The `layup psb check` overclaim, the fact numbering and the escalation
   floor are this one error.
7. **The human before the independent review.** Plan-review condition 4 put the
   Operator's approval before the freeze, and the author followed that order:
   the approval request reached the Operator before any independent reviewer
   read the design. The Operator was the first reviewer of about 1,000 lines
   (PSB Problem 1).
8. **One author, one pass.** Eight ADRs and the architecture came from one
   session in about 40 minutes, with no scenario test between them; the panel's
   options were compared per question, and no step checked the questions
   together (for example, Intake before Scaffold).

## The fixes

| Root cause | Fix | Checked by |
| ---------- | --- | ---------- |
| 1 | The public-solution search runs before any decision that fixes a solution shape | the search record is dated before the decision that uses it |
| 2 | Queries come from the problem statement's own terms as well as the solution's | the record lists each query with its source in the problem text |
| 3 | Each curated list is read in full, never summarized with a filter or a cap; each entry that touches the problem gets a keep or reject reason | the record holds the full candidate list with a reason per entry |
| 4 | Two independent searchers on different models and harnesses, blind to each other and to the author's conclusion; the panel brief and each review lens include "name the existing solutions this design repeats" | two search outputs under `runs/`; the lens line |
| 5, 6 | A scenario walkthrough per In-Scope item: actor, input, mechanism, record, each step tagged `code`, `model` or `human`; a coverage row points to its walkthrough | a coverage row with no walkthrough fails |
| 7 | An independent deep review runs before a human is asked to approve | the approval request cites a round that ended `nothing material` |
| 8 | The architecture is written in slices, each walked and reviewed before the next | the plan lists the slices |

The first application is the second search of this task, under
`runs/T-hbw8/search-v2/`. The lesson goes into `docs/guardrails.md` §2.
