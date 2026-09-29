# W-01 — Problem Statement Quality

**Item.** S1, Problem Statement Quality (`F-0003#41`): "Detection of the gaps in
a problem statement before delivery starts, so that the questions for the idea
owner come in one batch and not one at a time." Also Decision Points 1 and 2
(`F-0001#10`, `#11`) and `REQ-001`.

**The case.** A problem statement for a small invoicing service. Its §7 table has
the metric "Invoice accuracy ≥ 99 %" with no measurement method, and its text
uses "customer" for both the payer and the account holder. It names no budget.

Sections are those of [`architecture.md`](../architecture.md).

| # | Actor | Input | Mechanism | Record | Tag | Where |
| - | ----- | ----- | --------- | ------ | --- | ----- |
| 1 | the Operator | an empty repository on the forge | installs the LAYUP App on it; runs `layup run --new OWNER/NAME --psb brief.md` | — | `human` | open |
| 2 | `layup run` | the repository; `brief.md` | open | open | `code` | open |
| 3 | `layup psb check` | `brief.md` | rules G1 to G5 (`internal/psb`); here G2 fires on the metric row | a gap table on standard output | `code` | open |
| 4 | a review session | `brief.md` | open | open | `model` | open |
| 5 | spec sessions | `brief.md` | later: slice C | later: slice C | `model` | — |
| 6 | `layup run` | the gap tables of steps 3 to 5; the questions of the pinned baseline's setup | open | open | `code` | open |
| 7 | `layup run` | the batch | open | open | `code` | open |
| 8 | the idea owner, the Operator | the Intake comment | open | open | `human` | open |
| 9 | `layup run` | the answer comments | open | open | `code` | open |
