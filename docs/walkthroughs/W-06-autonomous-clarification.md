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
| 1 | the developer session | its work | ends with status `needs_context`: the question, its own label "interface contract", and the records it concerns | the result file | `model` | §8 |
| 2 | `layup run` | the result | records the question with its time; posts it on the issue of `T-9` | a row of the project question table | `code` | §8; ADR-0019 |
| 3 | the smart-if (P2) | the question and its records | the kind "interface contract" and "needs a human: no" (under `shadow`, the asker's label decides, and the row is recorded) | a `decisions.tsv` row | `model` | §10; ADR-0021 |
| 4 | `layup run` | the kind; the confirmed owner map | "interface contract" → Software Architect; routes a session of that role (§9) | a routing row | `code` | §9; ADR-0020 |
| 5 | the Software Architect session | the question; the records it concerns; the specification | answers "the total includes tax", citing the specification section | its result | `model` | §8 |
| 6 | the escalation screen | the answer | the floor finds nothing (no diff); P1's four questions answer no | a `decisions.tsv` row | `code`, `model` | §10; ADR-0022 |
| 7 | `layup run` | the answer | records it with its time; posts it on the issue; starts the developer's next attempt with the answer in its prompt file | the answer row; an attempt row that does not count toward §11's limits | `code` | §8 |
| 8 | the developer session | the answer | ends `completed`, its result cites the answer ID, and none of its questions cites it: the answer is accepted, by this session, at this time | the accepted time and actor in the question row | `model`, then `code` | §8; ADR-0019 |

## Checklist rows

S6, K25, K26, P03 (in part); the turnaround: known limit L-D1.
