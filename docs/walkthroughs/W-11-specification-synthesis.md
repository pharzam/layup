# W-11 — Specification Synthesis

**Item.** S11, Specification Synthesis (`F-0003#51`): "Derivation of identified
requirements and technical specifications from the approved problem statement.
Each requirement has a trace to the text that states it, and an acceptance
criterion." Also Specification Traceability (`F-0003#62`), `REQ-012`, and vision
2.1.

**The case.** The invoicing statement of W-01. Its §2 says, inside a paragraph
about context, "invoices must be sent by e-mail on the first working day". The
numbering session classes that sentence as `context`, so the draft has no
requirement for it.

Sections are those of [`architecture.md`](../architecture.md).

| # | Actor | Input | Mechanism | Record | Tag | Where |
| - | ----- | ----- | --------- | ------ | --- | ----- |
| 1 | a numbering session | `brief.md` | open | open | `model` | open |
| 2 | `layup spec check` | the fact list; `brief.md` | open | open | `code` | open |
| 3 | a draft session | the facts | open | open | `model` | open |
| 4 | a review session on another harness | `brief.md`; the draft | open | open | `model` | open |
| 5 | the idea owner | the Intake batch | open | open | `human` | open |
| 6 | a spec session | the draft; the answers | open | open | `model` | open |
| 7 | the idea owner | the first bet | open | open | `human` | open |
| 8 | `layup spec check` | the target's facts and PRD | open | open | `code` | open |
| 9 | a verifier session | a delivered requirement and its specification | later: slice D | later: slice D | `model` | — |
