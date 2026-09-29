# W-12 — Harness-Agent Neutrality

**Item.** S12, Harness-Agent Neutrality (`F-0003#52`): "The rules, the context,
and the task state stay in a form that does not belong to one harness agent. A
second harness agent can do the work, and can check the work of the first
one." Also Invariants 1, 2 and 9, and `NFR-002`.

**The case.** Task `T-7` of a Go target starts on harness H1 (Claude Code). It
stalls, and the Operator gives it to harness H2 (Codex). Later the Operator
stops LAYUP for good, and a human continues the target without it.

## Part 1 — the task moves from H1 to H2

| # | Actor | Input | Mechanism | Record | Tag |
| - | ----- | ----- | --------- | ------ | --- |
| 1 | `layup run` | the task list and the routing table | later: slice D | later: slice D | `code` |
| 2 | `layup run` | the target at the base commit; the harness register row of H1 | open | open | `code` |
| 3 | H1 session | its prompt file and its work tree | open | open | `model` |
| 4 | `layup run` | the stalled attempt | later: slice F | later: slice F | `code` |
| 5 | the Operator | the stall package on the task's issue | a comment "give `T-7` to H2" (Decision Point 5) | open | `human` |
| 6 | `layup run` | that comment, read from the forge | open | open | `code` |
| 7 | `layup run` | the parsed comment | later: slice F | later: slice F | `code` |
| 8 | `layup run` | the target at the head of attempt 1; the records; the register row of H2 | open | open | `code` |
| 9 | H2 session | the same kind of prompt file, built from the records; no memory of H1 | open | open | `model` |
| 10 | `layup run` | the result of H2 | open | open | `code` |
| 11 | a verifier session | the change | later: slice D | later: slice D | `model` |

## Part 2 — a human continues without LAYUP

| # | Actor | Input | Mechanism | Record | Tag |
| - | ----- | ----- | --------- | ------ | --- |
| 12 | the Operator | — | stops `layup run` | open | `human` |
| 13 | the Operator | the target's rulesets | later: slice B | later: slice B | `human` |
| 14 | a human | a fresh clone of the target | open | — | `human` |
| 15 | the target's CI | a pull request of that human | later: slice B | later: slice B | `code` |
| 16 | a human | the open task `T-7` | open | open | `human` |

## Checklist rows

K01 to K07; P04, P09, P13, P21; S12; R07, R08, R11, R13; I1, I2, I9; C2, C6, C7;
D05 to D08, D17; FT3, FT5, FT6.
