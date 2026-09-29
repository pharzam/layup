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
| 1 | a numbering session | `brief.md` | splits it into spans given by byte offsets: facts with an ID and a class, and "not a fact" spans with a reason; the sentence "invoices must be sent by e-mail on the first working day" is classed `context`; headings are "not a fact" spans | a span table in the session result | `model` | §7; ADR-0018 |
| 2 | `layup spec check --facts` | the span table; `brief.md` | every byte that is not white space lies in exactly one span, and each span's text matches the file; on a failure, no batch is posted | a check table, committed by `layup run` | `code` | §7; ADR-0018 |
| 3 | a draft session | the facts | writes requirements for each `need` fact, with `covers`, a statement and a criterion, and no priority; the e-mail sentence gets none, because it is `context` | the draft in the session result | `model` | §7; ADR-0018 |
| 4 | a review session on another harness | `brief.md`; the span table; the draft | lists "invoices by e-mail on the first working day" as a need that no requirement covers, and the fact whose class it doubts | a gap row in the session result; it joins the Intake batch (W-01 step 5) | `model` | §7; ADR-0018 |
| 5 | the idea owner | the Intake batch | answers the question: yes, it is a need | the answer comment, copied (§3) | `human` | §5 |
| 6 | a spec session | the draft; the answers | reclasses the fact as `need`, citing the answer ID (code checks the citation); adds a requirement that covers it and the answer's fact; writes the PRD with its MoSCoW and Phase columns empty, and one section per requirement in `docs/spec/` | a pull request that waits for the first bet | `model` | §7; ADR-0018 |
| 7 | the idea owner | the first bet, which shows every fact with its class and text | answers one line per requirement ID: its priority and milestone (had step 4 missed the sentence, a line "`need`" for its fact ID would reopen the draft) | the bet comment, copied (§3) | `human` | §7, §8 |
| 8 | `layup run` | the copied bet comment; the span table; the reclasses of step 6 | writes the confirmed inventory version and its SHA-256 on the records branch; adds one rendered commit to the verified specification's head: the MoSCoW and Phase columns, the non-goals, and the numbered facts record; checks that it changes only those and carries `layup/verify` to it | the inventory version; the rendered commit | `code` | §7; ADR-0018 |
| 9 | `layup spec check` | the latest confirmed inventory, the answers fact and the bet copy (records); the PRD and `docs/spec/` (the head); the task register; the delivered requirements | every `covers` resolves; every need and constraint is covered; each requirement has a criterion; MoSCoW and Phase match the bet; each delivered requirement has a section whose heading holds its ID | a check table, committed by `layup run`, and the status `layup/spec` | `code` | §7; ADR-0018 |
| 10 | a verifier session on a harness outside the authors | each change that delivers the requirement, with its specification section | judges that the change and the section agree with the requirement's text (W-07 step 7) | its record and `layup/verify` | `model` | §8, §9; ADR-0019, ADR-0020 |
| 11 | the approver | the delivered requirement (every task that names it has merged), its criterion, its changes, and its gate and verification results (§8 Accept) | accepts it by one comment (Decision Point 3) | the comment, copied (§3) | `human` | §8 |
| 12 | `layup run` | the acceptance | writes an `acceptance.tsv` row: review number 1 (no earlier row), the time; Delivery Lead Time is this time minus the start of the requirement's first task | the row | `code` | §12; ADR-0024 |

## Checklist rows

S11, K19 to K23, P07, P14; FT1 (the line check). Known limit L-C1.
