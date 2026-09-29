# W-08 — Human-on-the-Loop

**Item.** S8, Human-on-the-Loop (`F-0003#48`): "Humans monitor delivery and give
input only at the Human Decision Points below. A rule that a machine can apply
(the escalation rule) selects the business-forking decisions, for example
budget, legal compliance, and strategic trade-offs." Also Decision Point 4
(`F-0001#13`), Missed Escalations (`F-0003#57`), `REQ-008`, and O-67, O-68, O-71,
O-78, O-79, O-84.

**The case.** Build task `T-14` of the invoicing target must send invoices by
e-mail. Its plan session proposes a paid e-mail service. Neither the problem
statement nor the band names such a cost.

Sections are those of [`architecture.md`](../architecture.md).

| # | Actor | Input | Mechanism | Record | Tag | Where |
| - | ----- | ----- | --------- | ------ | --- | ----- |
| 1 | the plan session | the task | writes the plan | the plan | `model` | open |
| 2 | `layup run` | the plan | open | open | `code` | open |
| 3 | the smart-if | the plan | open | open | `model` | open |
| 4 | `layup run` | the two results | open | open | `code` | open |
| 5 | `layup run` | the selected decision | open | open | `code` | open |
| 6 | the idea owner | the escalation | open | open | `human` | open |
| 7 | `layup run` | the answer | open | open | `code` | open |
| 8 | the developer session | its work | adds the paid service's client module anyway | commits | `model` | §4 |
| 9 | `layup run` | the handoff diff | open | open | `code` | open |
