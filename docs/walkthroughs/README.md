# Walkthroughs

Each file follows one concrete case of one In-Scope item of the PSB
(`F-0003#41`–`#52`) through LAYUP, step by step. It is the test of
[`architecture.md`](../architecture.md): a coverage row points to a walkthrough,
and a walkthrough step points to the section and the ADR that give its
mechanism.

A step names its **actor**, its **input**, its **mechanism**, the **record** it
writes, and a **tag**:

- `code` — a deterministic step in `layup`; it names its input file and the rule
  it applies;
- `model` — a role session on a harness, or the smart-if provider;
- `human` — the idea owner or the Operator, at a named decision point.

A step that reads meaning (a gap of meaning, a business fork written in prose, a
clause of a problem statement) is never `code`. A step whose mechanism is not
designed yet says `open`; one whose mechanism belongs to a later slice of task
`T-hbw8` says `later: slice X`.
