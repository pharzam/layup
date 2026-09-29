# 0018. Derive the specification from numbered source lines

Date: 2026-09-29

## Status

Proposed

## Context

The PSB asks for requirements and technical specifications derived from the
approved problem statement, each with "a trace to the text that states it, and
an acceptance criterion" (`F-0003#51`, `#62`); the intent of a project, which
includes what counts as success, is the idea owner's (`F-0003#54`). The deep check
found that the first draft let code number the "clauses" of prose, which only a
reader of meaning can do (Sol-17, Fable-M21, Author-5); that a byte-exact trace
proves links but finds no missed need (Sol-18); that code set `Must` for each
In-Scope fact (Sol-19, Author-18); that no check asked for a specification per
requirement (Sol-20); and that a requirement traced to an answer failed (Fable-N8).

The evaluation gave the pattern of numbered source lines with a `covers` field
checked byte for byte (BMAD, Spec Kitty), and showed a trace to a line that does
not exist passing, and a gap filled with a guess (BMAD, MetaGPT). Shape Up's
scope hammering splits requirement priority, which is intent, from task priority
([`runs/T-hbw8/evaluation/19-shape-up.md`](../../runs/T-hbw8/evaluation/19-shape-up.md)).
The options compared are in
[`runs/T-hbw8/selection-v2.md`](../../runs/T-hbw8/selection-v2.md), row 0018.

## Decision

We will derive the specification by sessions and check its links by code:

1. **A session numbers** the problem statement into byte-exact facts with an ID
   and a class (`need`, `constraint`, `measure`, `context`), and lists the text
   that is no fact with a reason.
2. **Code checks the lines:** each fact is a substring, no two overlap, and every
   non-blank line lies in a fact or in the "not a fact" list.
3. **A session drafts** the requirements with `covers` and a criterion, and no
   priority; **a session on another harness** reviews the draft for a missed
   need; each gap it finds is a question of the one Intake batch.
4. **The idea owner** confirms the needs and sets each requirement's priority at
   the first bet.
5. **`layup spec check`** fails on an unresolved `covers`, a quote that is not
   byte-exact, an uncovered need that the idea owner did not mark out of scope, a
   requirement with no criterion or no priority after the first bet, a delivered
   requirement with no specification section that names it, and a task with no
   requirement. A trace may point to the problem statement or to the answers.

We reject: code numbering the clauses of prose (it can only split lines); code
setting a priority (intent); the smart-if provider extracting requirements (it
writes no text); a trace check as the proof of completeness.

## Consequences

- The problem statement needs no fixed form: the numbering session reads any.
- The idea owner answers the completeness gaps in the Intake batch, before
  delivery, and sets the priorities once, at the first bet.
- A session can still class a fact wrongly; the completeness review on another
  harness and the idea owner's confirmation are the two checks against it.
- `PRD-0001` REQ-012 keeps its criterion; the meaning of a specification is judged
  by the counterpart verification of each change, not by `layup spec check`.
