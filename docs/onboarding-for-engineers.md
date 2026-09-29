# Onboarding for engineers

> **Who this is for.** Software engineers joining this project who do **not**
> necessarily know the project's domain. You know how to build software. This
> document teaches you the project's terminology and the one or two facts that
> shape every decision, then hands you off to the deeper docs.
>
> **Read this first.** Every other document is written for someone already fluent
> in the terminology. This one gets you to that point. Budget half an hour.

## 1. Problem statement

The source of truth is the approved LAYUP Problem Statement Brief (PSB),
Revision 6, stored as the raw fact [`F-0001`](facts/F-0001-layup-problem-statement-brief.md)
next to the file itself, [`facts/problem-statement-brief.md`](facts/problem-statement-brief.md).
This section is a plain summary. Where the two disagree, the PSB wins.

Teams now give software work to several autonomous agents, each with one
function, for example a product owner, an architect, an engineer, and a QA
engineer (`F-0001 §2`). The agents write code fast, but the delivery does not
get fast. The PSB names seven causes (`F-0001 §1`):

1. **The human is the message bus.** An agent that meets an unclear point stops
   and asks a human, because it cannot tell who owns the answer.
2. **Discipline decays.** Nothing that a machine enforces keeps the layout, the
   boundaries, and the tests the same across agents and sessions.
3. **Gaps are found late, and setup is manual.** Gaps in the problem statement
   come up one at a time during delivery, and the discipline template of a new
   project is filled by hand, with guesses.
4. **Deadlock has no exit.** When agents disagree or a step repeats, nothing stops
   the loop or diagnoses it.
5. **Cost is invisible.** Nobody knows the tokens, the latency, or the duration
   of a task.
6. **Specifications are written by hand.** Nothing turns the approved problem
   statement into requirements that trace back to its text.
7. **Harness agents are fragmented.** The rules are copied into the format of
   each agent product, and the copies drift.

The PSB states the problem only. It does not choose a solution (`F-0001` header
note). This repository has not yet answered any part of the problem; see
[Where the project stands](#where-the-project-stands).

## 2. Crash course: the domain

The terms below build on each other. PSB §8 defines each one
(`F-0001#15`–`F-0001#39`). The [glossary](glossary.md) takes these terms in task
`T-xgz4` ([#6](https://github.com/pharzam/layup/issues/6)).

### 2.1 Harness agent and role agent

A **harness agent** is an agent product, for example Claude Code or Codex. It
gives the model, the session, the tools, and the context window. A **role agent**
is an autonomous agent with one function. A harness agent runs role agents
(`F-0001#18`, `F-0001#21`).

### 2.2 Discipline system, rule, gate, content

The **discipline system** is the set of rules and gates that a machine enforces
in a project repository. Armature is its baseline (`F-0001#19`). A **rule** says
what the work must satisfy; a **gate** applies rules to a change and gives pass
or fail; **content** is the project-specific text and values that the rules
check or use (`F-0001#34`–`F-0001#36`).

### 2.3 Operator, idea owner, and the Human Decision Points

The **Operator** starts, monitors, and approves. The **idea owner** owns the
problem, the success criteria, and the funding (`F-0001#22`, `F-0001#23`).
Humans monitor at any time, but monitoring is not input. Humans give input only
at five **Human Decision Points**: intent decisions, problem statement answers,
planned approval points, escalations, and stalls (`F-0001#10`–`F-0001#14`).

### 2.4 Task, stall, telemetry

A **task** is one unit of delivery work with one goal (`F-0001#29`). A **stall**
is a task that neither reaches its goal nor fails cleanly (`F-0001#37`).
**Telemetry** is the record of the token count, the latency, and the wall-clock
duration of a task (`F-0001#38`).

## 3. What the system actually does

The PSB states the problem only. This repository holds the discipline system,
the setup record, and the first code: the `layup` command
([ADR-0011](adr/0011-structure-the-core-engine-as-a-go-cli-over-repository-files.md)).
So far it prints its version, and `layup psb check FILE` writes the gap questions
of a problem statement as one batch (five deterministic rules, G1 to G5, in
`internal/psb`).

## 4. Why it is hard

Nine System Invariants bind any solution (`F-0001#1`–`F-0001#9`).
[`guardrails.md`](guardrails.md) takes them in task `T-7ndb`
([#7](https://github.com/pharzam/layup/issues/7)). Two of them stop the obvious designs:

- **Git is the system of record** (`F-0001#1`). A dashboard or a database that
  holds project state is not allowed to be the only copy.
- **The project repository is independent** (`F-0001#2`). A human or a different
  agent must be able to continue the work without the automation that made it.

The success criteria have two layers: checks that must pass, from the
invariants, and numbers that a pilot calibrates (`F-0001 §7`).

## 5. How this project works

### The one cultural thing to understand

**No value without evidence.** A value that comes from a guess is a defect
(`F-0001#4`). When you do not know a value, you do not fill it: you write it down
as an open question for the idea owner. This surprises engineers who are used to
a sensible default.

Consequences you will meet immediately, and which are not negotiable:

- Every setup value has a row with its evidence in
  [`setup/record-T-n1hp.md`](setup/record-T-n1hp.md). A reviewer checks the rows;
  [`setup/setup-check.sh`](setup/setup-check.sh) checks what a script can: the
  pin, the facts, and this document.
- A check that is not active does not count as passed (`F-0001#5`). The record
  says `not active` for a value that no gate runs yet.
- Git is the system of record (`F-0001#1`). A decision made in a chat goes onto
  its issue or into an ADR before anyone acts on it.

### Repo layout

| Path | What |
|------|------|
| [`setup/`](setup/) | The Armature pin, the setup record with the evidence for each value, and `setup-check.sh`, which proves the setup. |
| [`AGENTS.md`](../AGENTS.md) | The agent entry point — the gate in brief and a pointer to the R1–R13 rules, in one short file. [`CLAUDE.md`](../CLAUDE.md) imports it for Claude Code. |
| [`engineering-discipline.md`](engineering-discipline.md) | **How we work**: the quality gate, solution selection, branches, worktrees, commits, tests, reviews, and ADRs. Read before your first commit. |
| [`glossary.md`](glossary.md) | The shared vocabulary. Skim it; come back constantly. |
| [`facts/`](facts/) | Facts stored as-is as immutable evidence: the PSB (`F-0001`, with more numbered facts in `F-0003`), the vision brief (`F-0002`, a solution document, not requirements) and the idea owner's answers to the PSB's gap questions (`F-0004`). Derived requirements cite them by `F-NNNN` ID. |
| [`prd/`](prd/) | Product Requirements Documents, derived from the facts; each `REQ`/`NFR` cites an `F-NNNN` fact. [`PRD-0001`](prd/PRD-0001-layup.md) is LAYUP's own: the full scope of the PSB in four phases (Draft until the PDR records the Operator's acceptance). |
| [`issue-workflow.md`](issue-workflow.md) | The issue-first rules (R1–R13): the ticket policy the gate assumes. |
| [`tasks/backlog.md`](tasks/backlog.md) | What to work on next. |

### What to read next, in order

0. [`AGENTS.md`](../AGENTS.md) — the one-page summary, if you are a coding agent
   (or a human who wants the shape before the detail).
1. [`engineering-discipline.md`](engineering-discipline.md) — how we work.
2. [`issue-workflow.md`](issue-workflow.md) — the issue-first rules the gate assumes.
3. [`glossary.md`](glossary.md) — skim, then reference.
4. [`guardrails.md`](guardrails.md) — the pitfalls and the frozen numbers.
5. [`tasks/backlog.md`](tasks/backlog.md) — what needs doing.

### Where the project stands

The Armature setup of this repository is done
([#1](https://github.com/pharzam/layup/issues/1), closed). The core engine is in
progress under [#29](https://github.com/pharzam/layup/issues/29); that parent
issue lists its children, and [`tasks/completed.md`](tasks/completed.md) lists the
ones that are done, so this section does not repeat the list. The open setup
questions are listed in [`setup/open-gaps.tsv`](setup/open-gaps.tsv). Since
2026-09-25 the repository is in
[bootstrap mode](engineering-discipline.md#bootstrap-mode)
([ADR-0012](adr/0012-build-layup-in-bootstrap-mode.md)): the product path is
[#45](https://github.com/pharzam/layup/issues/45) (the idea owner's answers),
then [#42](https://github.com/pharzam/layup/issues/42) (the PRD), then the
core-engine commands under #29, and the gate runs with one plan-review comment
and one review round, one more after a fix (two for a change to a gate), until the ADR that supersedes
ADR-0012 is accepted (a summary; the rule is the linked section).
