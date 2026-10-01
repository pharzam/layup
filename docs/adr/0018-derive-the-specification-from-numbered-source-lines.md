# 0018. Derive the specification from numbered source lines

Date: 2026-09-29

## Status

Accepted

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

1. **A session numbers** the problem statement into spans given by byte offsets:
   facts with an ID and a class (`need`, `constraint`, `measure`, `context`), and
   "not a fact" spans with a reason.
2. **Code checks the spans** (`layup spec check --facts`): every byte that is not
   white space lies in exactly one span, and each span's text matches the file.
   When it fails or does not run, no Intake batch is posted.
3. **A session drafts** a requirement for each need and a non-functional
   requirement for each constraint, with `covers` and a criterion and no
   priority; a measure becomes a success criterion that cites its ID. **A session on another
   harness** reviews the draft for a missed need or constraint; each gap is a
   question of the one Intake batch. With one harness, the review does not run
   and `layup run` stops.
4. **The answers** become a raw fact with one ID per question ID; each question
   carries the quote of the problem statement it asks about, so a trace through
   an answer reaches the problem statement's text.
5. **The idea owner** sets each requirement's priority and milestone, and marks a
   fact out of scope or gives it a new class, by one comment at each bet. Code
   writes the PRD's MoSCoW and Phase columns from the copied comment, and a new
   version of the confirmed inventory with its hash on the records branch; code
   renders the numbered facts record from it, and the record lands in
   `docs/facts/`, a rule path, as its own task.
6. **`layup spec check`** posts the required status `layup/spec`. It checks the
   head against the confirmed inventory and the bet copy: every `covers`
   resolves, every need and constraint is covered, every measure has a success
   criterion, or the fact is out of scope, each
   requirement has a criterion, the PRD's MoSCoW and Phase match the bet, each
   delivered requirement has a section in `docs/spec/` whose heading holds its
   ID, and each task names a requirement. A check that did not run fails.

We reject: code numbering the clauses of prose (it can only split lines); a
session or code setting a priority (intent); the smart-if provider extracting
requirements (it writes no text); a trace check as the proof of completeness;
checking the head's own copy of the facts (FT4).

## Consequences

- The problem statement needs no fixed form: the numbering session reads any.
- The idea owner answers the completeness gaps in the Intake batch, before
  delivery, and sets the priorities at each bet.
- A session can still class a fact wrongly; the completeness review on another
  harness and the idea owner's confirmation are the two checks against it.
- The reading of `F-0003#62` for a trace through an answer goes to the Operator
  in the approval brief.
- `PRD-0001` REQ-012 keeps its criterion; the meaning of a specification is judged
  by the counterpart verification of each change, not by `layup spec check`.
