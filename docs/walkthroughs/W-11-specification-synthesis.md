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
| 1 | a numbering session | `brief.md` | splits it into facts, each a byte-exact substring with an ID and a class; the sentence "invoices must be sent by e-mail on the first working day" is classed `context`; headings go to the "not a fact" list | a fact table and a "not a fact" list in the session result | `model` | §7; ADR-0018 |
| 2 | `layup spec check` | the fact list; `brief.md` | each fact is a substring of the file; no two overlap; every non-blank line lies in a fact or in the "not a fact" list | a check table, committed by `layup run` | `code` | §7; ADR-0018 |
| 3 | a draft session | the facts | writes requirements for each `need` fact, with `covers`, a statement and a criterion, and no priority; the e-mail sentence gets none, because it is `context` | the draft in the session result | `model` | §7; ADR-0018 |
| 4 | a review session on another harness | `brief.md`; the draft | lists "invoices by e-mail on the first working day" as a need that no requirement covers, and the fact whose class it doubts | a gap row in the session result; it joins the Intake batch (W-01 step 5) | `model` | §7; ADR-0018 |
| 5 | the idea owner | the Intake batch | answers the question: yes, it is a need | the answer comment, copied (§3) | `human` | §5 |
| 6 | a spec session | the draft; the answers | reclasses the fact as `need`, adds a requirement that covers it, and writes the PRD and one specification section per requirement | a pull request of the first bet | `model` | §7; ADR-0018 |
| 7 | the idea owner | the first bet | confirms the needs and sets the priority of each requirement | the bet comment, copied (§3) | `human` | §7, §8 |
| 8 | `layup spec check` | the target's facts and PRD | every `covers` resolves; every need is covered; every requirement has a criterion and a priority; each delivered requirement has a specification section that names it | a check table, committed by `layup run`, and the status `layup/gates` | `code` | §7; ADR-0018 |
| 9 | a verifier session | a delivered requirement and its specification | later: slice D | later: slice D | `model` | — |

## Checklist rows

S11, K19 to K23, P07, P14; FT1 (the line check). Known limit L-C1.
