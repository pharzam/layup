# W-05 — Role Handoffs

**Item.** S5, Role Handoffs (`F-0003#45`): "Information that passes between role
agents has a form that a machine can validate, based on … conventions."
Also Inter-Role Communication Format (`F-0003#59`) and `REQ-005`.

**The case.** In milestone 1 of the invoicing target, the specification task
`T-4` hands its work to the milestone plan. Its result names the section
`docs/spec/billing-api.md` for requirement `REQ-003`, but the file is not in its
commits.

Sections are those of [`architecture.md`](../architecture.md).

| # | Actor | Input | Mechanism | Record | Tag | Where |
| - | ----- | ----- | --------- | ------ | --- | ----- |
| 1 | the specification session | its task | ends with a handoff result: status `completed`, and an artifact `docs/spec/billing-api.md` with a SHA-256 | the result file | `model` | §8 |
| 2 | `layup run` | the result; the transition table row "specification → milestone plan: a section in `docs/spec/` whose heading holds each requirement ID"; the commits | checks the schema; looks for each named artifact in the commits: `docs/spec/billing-api.md` is missing, so no section holds `REQ-003`, and the handoff is invalid, whatever its status field says | a handoff row with `valid = no` and the reason | `code` | §8; ADR-0019 |
| 3 | `layup run` | the refused handoff | gives the specification task a new attempt with the finding; the milestone plan does not start | an attempt row | `code` | §8 |
| 4 | `layup report` | the handoff rows | valid handoffs over all handoffs (`F-0003#59`); this one counts as not valid | the report | `code` | §8, §12 |

## Checklist rows

S5, R04 (in part), K24, P03 (in part); FT3.
