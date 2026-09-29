# W-05 — Role Handoffs

**Item.** S5, Role Handoffs (`F-0003#45`): "Information that passes between role
agents has a form that a machine can validate, based on … conventions."
Also Inter-Role Communication Format (`F-0003#59`) and `REQ-005`.

**The case.** In milestone 1 of the invoicing target, the Software Architect
session of task `T-4` hands its work to the Software Developer. Its result names
the interface file `docs/spec/billing-api.md`, but the file is not in its
commits.

Sections are those of [`architecture.md`](../architecture.md).

| # | Actor | Input | Mechanism | Record | Tag | Where |
| - | ----- | ----- | --------- | ------ | --- | ----- |
| 1 | the architect session | its task | open | open | `model` | open |
| 2 | `layup run` | the result; the transition table | open | open | `code` | open |
| 3 | `layup run` | the refused handoff | open | open | `code` | open |
| 4 | `layup report` | the handoff records | open | open | `code` | open |
