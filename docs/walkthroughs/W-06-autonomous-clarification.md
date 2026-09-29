# W-06 — Autonomous Clarification

**Item.** S6, Autonomous Clarification (`F-0003#46`): "A question that does not
need a human decision gets an accepted answer from the responsible role agent,
without a human." Also Problem 1, Clarification Turnaround (`F-0003#71`) and
`REQ-006`.

**The case.** A developer session of task `T-9` does not know if the invoice
service returns a total with or without tax. It is an interface-contract
question.

Sections are those of [`architecture.md`](../architecture.md).

| # | Actor | Input | Mechanism | Record | Tag | Where |
| - | ----- | ----- | --------- | ------ | --- | ----- |
| 1 | the developer session | its work | open | open | `model` | open |
| 2 | `layup run` | the result | open | open | `code` | open |
| 3 | the smart-if | the question | later: slice E | later: slice E | `model` | — |
| 4 | `layup run` | the kind; the owner map | open | open | `code` | open |
| 5 | the owner session | the question and its context | open | open | `model` | open |
| 6 | the escalation screen | the answer | later: slice E | later: slice E | `code`, `model` | — |
| 7 | `layup run` | the answer | open | open | `code` | open |
| 8 | the developer session | the answer | open | open | `model` | open |
