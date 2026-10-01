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
| 1 | the plan session | the task | writes the plan; its fixed field `new dependencies` names the paid e-mail service's client module by its manifest identifier, with its monthly cost and the alternative it considered (plain SMTP, the standard library) | the plan | `model` | §8, §10 |
| 2 | `layup run` | the plan's `new dependencies` field | the floor: the module is not on the allowed-dependency list, so it is a candidate | a candidate row | `code` | §10; ADR-0022 |
| 3 | the smart-if (P1) | the plan's text | the four literal questions; "does it change the budget?" answers yes with a high probability | a `decisions.tsv` row | `model` | §10; ADR-0021 |
| 4 | `layup run` | the floor and P1 | OR in code: selected (under `shadow`, the floor alone selects it; the P1 row is only recorded, and the selection row says that P1 was not active) | a `candidates.tsv` row; the task in a wait state | `code` | §10; ADR-0022 |
| 5 | `layup run` | the selected decision | posts one escalation brief on the task's issue: the decision, the numbered options from the plan's field (1: the paid service; 2: plain SMTP), the evidence, the cost | the comment ID | `code` | §10 |
| 6 | the idea owner | the escalation | answers `option 2; business-forking: yes` (Decision Point 4) | the comment | `human` | §10 |
| 7 | `layup run` | the answer | copies it (§3); writes `escalations.tsv`: chosen "plain SMTP", rejected the client module by its manifest identifier; records it as planned input; starts the task's next attempt with the decision in its prompt file | the decision row; an attempt row | `code` | §3, §10, §12 |
| 8 | the developer session | its work | adds the paid service's client module anyway | commits | `model` | §4 |
| 9 | `layup run` | the handoff diff | the floor: a new `require` in `go.mod` is a candidate; its manifest identifier equals that of the rejected option, so code gives the task a finding (no new brief) and refuses the result before any push | a finding row for `T-14` | `code` | §10; ADR-0022 |

## Checklist rows

S8, I4 (in part), I6, K32 to K41, P11, P18 (the smart-if part), C3, C9, C10, D02 (in part), D04, D09, D10, D15, D16 (the bet); FT1, FT3. Known limit L-E1.
