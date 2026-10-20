# 0028. Delegate two bounded landing decisions

Date: 2026-10-20

## Status

Accepted

## Context

The Operator decides how each build task lands, one decision at a time. [Reviewing until findings decay](../engineering-discipline.md#reviewing-until-findings-decay) makes two of those decisions the Operator's alone:

- **The landing figure.** An overrun of the plan review's `Budget maximum` blocks the merge until the Operator approves "one number, named once, … the figure at landing". No growth is approved in advance (the **Ceiling** of the [glossary](../glossary.md) is the thing the gate does not have).
- **The end at the cap.** With a material finding open at the cycle cap, only the Operator's comment permits end 3, a fix with no further round.

M2e measured the cost of these rules. The figures are from the investigation of [#438](https://github.com/pharzam/layup/issues/438) and its three data comments. They are counts the author made from the issue threads and the task records, not independently recomputed:

| Measure | Value in the 33 merged tasks of M2e |
| --- | --- |
| Operator decisions | about 55: about 8 landing figures, 10 ends at the cap, 9 goal counts, 4 merges of `main`, others |
| Waiting for the Operator | about 72 task-hours, mostly overnight; daytime answers in 3 to 75 min |
| Share of the diff | tests 46 %, `runs/` evidence 15 %, Go product code 25 % |
| A budget finding repeated | in #351, #352, #358, #360, #365 (rounds 1 to 4), #413, #415; #362 and #363 blocked by the budget alone |
| Ends at the cap | end 3 five times, end 1 four times, end 2 once; the end-3 findings were mostly index and record corrections |

Two forces meet here. The budget is R12's bound, and it is what caught #76, which grew without limit. And every fix after a round adds test and evidence lines, so a task near its maximum overruns again at each round. The Operator then names a figure that differs from the last one by a few lines, or permits a fix to a traceability row that no check reads.

## Decision

We will delegate two bounded landing decisions to the author: a small overrun of the budget, and end 3 at the cap for a records finding. Each is recorded on the issue, and the Operator can overrule it there before the merge. Both are clauses of one choice: decisions the Operator has delegated in advance, each with a fixed bound.

1. **What the budget counts.** The unit stays lines added plus lines removed on the whole branch diff, plus files touched, against a named base. Two kinds of change leave the count:
   - **Evidence records** under `runs/<task>/`, meaning the Markdown, TSV and text files that record runs, times and outputs. A file under `runs/` that a check, a test or a script executes or reads as its input (for example `runs/T-gd8q/checks.sh`, or a file that such a script reads) still counts, lines and file.
   - **The close-out bookkeeping** that the gate already names: the completed-log line, the backlog line it removes, and the verdict and the resource record of the task's own file. In a file that holds both bookkeeping and other changes, only the bookkeeping lines leave the count, and the file still counts as touched.

   Executable tests and their inputs stay in the count, because tests are where the size of a task shows.
2. **The tolerance.** A landing figure over the `Budget maximum` by no more than the tolerance is approved in advance. The tolerance is 20 % of the maximum's lines, at least 60 lines and at most 500, plus 4 files. The author posts the figure, the maximum and the margin on the issue before the merge. A figure past the tolerance stays the Operator's, one number once, as now.

   **This amends the rule that no growth is approved in advance** (the **Ceiling** of the glossary). The tolerance is one bound, fixed before the work and the same for every task. It is not a multiple, not a reach that grows with the findings, and not a second approval. The floor keeps a 200-line task from stopping on 41 lines. The ceiling keeps a 3,000-line task from gaining 600 lines, a whole small task, unasked.
3. **The default end at the cap.** At the cap, when **every** material finding still open is a records finding, end 3 applies without an Operator comment. A records finding is one whose fix restores a named source or a record's form, and changes nothing else. Examples: a traceability row or a Test cell made to name a test that exists, an index line, a record's columns made to match their schema, trailing white space.
   - A fix that changes an exit code, an assertion, a verdict, a check's input, a test, product code, the specification, or a sentence that grants or limits an authority is **not** a records finding, and its end stays the Operator's.
   - The last verdict stays `not mergeable, findings recorded`. The task record names the fix that no round read, as end 3 already asks. What this clause replaces is the requirement of an individual Operator comment: the author records "end 3 by ADR-0028" with the findings and their fixes on the issue.
   - A disputed classification is the Operator's, as the material test already says.

**Rejected:**
- **No budget.** R12's bound is what stopped #76.
- **A percentage with no bounds.** 20 % of 3,000 lines is 600 lines, a whole small task, and 20 % of 150 lines is 30 lines, about one close-out.
- **Leaving test lines out of the count.** Tests are where the size shows.
- **An automatic end 1.** It adds a round, where the cost measured is the wait for a decision.
- **Two ADRs.** The two clauses share one context and one effect, the decisions at landing, as ADR-0026 kept several clauses of one gate in one record.

The author proposed 10 %, at least 30 and at most 250 lines, and 2 files, read from the overruns of M2e. The Operator accepted the record with each figure doubled, the floor confirmed as 60 ([#441](https://github.com/pharzam/layup/issues/441), 2026-10-20).

## Consequences

- **Fewer Operator decisions at landing.** In M2e this would have covered most of the landing figures and the records-only ends at the cap. The Operator still sees each figure and each end on the issue and can overrule it.
- **The budget reads closer to the work.** Evidence and bookkeeping no longer push a task over its maximum at close-out.
- **A new judgement for the author.** Whether a finding is a records finding. The bound is by effect, and a dispute goes to the Operator, so a wrong call is visible on the issue before the merge.
- **The rules take effect when the operative text changes**, in a later task. Until it lands, the gate reads as written. That task changes: [Reviewing until findings decay](../engineering-discipline.md#reviewing-until-findings-decay) (the budget and the end at the cap), the glossary's **Ceiling** and **Cycle cap** entries, the budget sentence of [ADR-0007](0007-record-task-resource-use.md), and the Status of [ADR-0026](0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md) to `Accepted. Amended by ADR-0028`. That task also decides whether [`review-record-lint`](../ci/review-record-lint.sh) needs a change. Today it checks the verdict's vocabulary and position, not the size of a change, so no change may be needed.
- **Open:** whether the figures fit M2f and later milestones, to be read again at the milestone's retrospective.
