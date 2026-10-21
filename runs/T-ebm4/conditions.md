# T-ebm4 plan-review conditions

All three plan-review conditions are implemented under amendment 1 on #451.
The unchanged Stage A check ran green in Stage B at 2026-10-21T07:36:20Z.
The complete command, measured UTC boundaries, exit status and assertion output
are in the run record. Independent review and the clean committed-head K1-K6
run belong to later stages; neither is claimed by this table.

| Condition | Evidence |
| --- | --- |
| 1. Require a new issue comment naming the figure or end; edits post nothing. The posting time selects the window. Complete. | `docs/engineering-discipline.md` contains the exact amended shared wait and an end-3 pointer. `runs/T-ebm4/checks.sh` passed C1-new-comment-wait, C1-retired-wait and C1-one-shared-wait; `runs/T-ebm4/semantic.md` walks both boundaries and the edited-comment case; `runs/T-ebm4/test-runs.md` records the green run. |
| 2. Drop decision 8 and preserve the complete validation section. Complete. | `docs/guardrails.md` retains the validation rule. `runs/T-ebm4/checks.sh` passed C2-validation-preserved against the baseline; `runs/T-ebm4/mentions.tsv` classifies its baseline line 1476 as validation; `runs/T-ebm4/test-runs.md` records the green run. |
| 3. Bound Records finding by "and changes nothing else" and add a complete UTC row. Complete. | `docs/glossary.md` contains the bounded Records finding and Coordinated Universal Time rows. `runs/T-ebm4/checks.sh` passed C3-records-bound and every C3-UTC assertion; `runs/T-ebm4/semantic.md` compares the clauses and examples; `runs/T-ebm4/test-runs.md` records the green run. |

Exactly four evidence records are declared: mentions.tsv, semantic.md,
conditions.md and test-runs.md, all under runs/T-ebm4/. Their consumers are
documented in the run and semantic records. The task check reads none as
configuration or expected output. K6 remains enabled and will audit committed
paths and this table's form at the freeze.
