# 0026. Keep the bootstrap review rules as the standing gate

Date: 2026-10-07

## Status

Accepted

## Context

[ADR-0012](0012-build-layup-in-bootstrap-mode.md) put this repository in bootstrap mode on 2026-09-25: one plan-review comment, one review round by a different model and one more after a fix, a stated test of "material", one home per rule, and a scope limited to the PSB's In-Scope items. Its part 6 ends the mode with the ADR that reads the first pilot's defect and stall numbers and re-decides the gate. The first pilot (`T-evad`, [#97](https://github.com/pharzam/layup/issues/97)) ended phase 1 on 2026-10-06. This record is that ADR (task `T-gd8q`, [#120](https://github.com/pharzam/layup/issues/120)).

The count is in [`runs/T-gd8q/numbers.md`](../../runs/T-gd8q/numbers.md), with its rule and a `file:line` for each value. **The rule was written after the numbers were known**, because the plan of the pilot did not pre-register a threshold ([`guardrails.md`](../guardrails.md) §1 asks for one first). So the count is not a blind test; the Operator accepted the rule on that basis (O-159 part 1). The numbers:

| Measure | Value |
| --- | --- |
| Phase 1: the 20 rows of the plan | 20 plan reviews, all `approve-with-conditions`; 35 rounds (9 tasks with 1, 9 with 2, 1 with 3, 1 with 5); 28 material findings; 18 ends by decay, 2 by the Operator raising the cap (rows 5 and 7: 4 `not mergeable` verdicts, 4 raises) |
| The other 10 tasks with a record under bootstrap mode | 20 rounds; 5 ends by decay, 5 by a decision of the Operator: 2 splits (`T-8ywj`, `T-wjq4`), 2 findings ruled notes (`T-7sbn`, `T-zmj6`), 1 fix with no further round after two raises (`T-0drh`) |
| Completed tasks, product : process | 31 : 3 from 2026-09-25; 5 : 12 before it |
| The pilot's defects in phase-1 code that passed its review | 16 findings in rows 6, 7, 8, 9, 10, 11, 13, 14 and 15; each in at least one row with two or more rounds (row 7 had five); row 11, with one round, shares two with row 13 |
| The full gate in the pilot's target task `T-vu2j` | 3 plan reviews (2 `reject`), 5 rounds with 4, 8, 3, 2 and 1 material findings, no decay |

Three readings of these numbers decide the gate:

1. **One round per frozen head converged.** 18 of the 20 reviews of phase 1 ended by decay. Each other end of a review under bootstrap mode was a decision of the Operator at the cap, of one of four kinds: raise the cap, rule the open findings notes or known limits, close out a fix with no further round, or split. The text of the gate named only the split.
2. **More rounds did not find the defects that the pilot found.** Each defect of the pilot is in at least one plan row with two or more review rounds, and row 7, with five rounds, holds one; of the nine rows with one round, only row 11 shares two defects, with row 13. Most of those defects are of integration or of the specification: a marker over more lines, a second setup run, a parameter that GitHub adds to a ruleset, a gate that a pull request can weaken. This shows a limit of a review of one diff, not that one round is enough: the rows with more rounds are also the larger rows that the pilot used most.
3. **The full gate did not decay in the pilot's target.** The target task `T-vu2j` ran under the full gate of its own repository, as written (one round with four lenses on each frozen head, a workaround for F-15): 3 plan reviews (`reject`, `reject`, `approve-with-conditions`) and 5 rounds with 4, 8, 3, 2 and 1 material findings; only three decisions of the Operator ended it (O-156 to O-158; F-38). That is root cause 5 of ADR-0012, measured again.

F-15 of the pilot is also an input: the written rule "One pass is never enough" and [`review-record-lint`](../ci/review-record-lint.sh), which reads one round per cycle, disagree, and the rule must not come back while they do.

## Decision

We will **keep the rules of bootstrap mode as the standing gate**, without its end date, and add the end at the cap as the Operator used it (O-159 part 2, option C). Bootstrap mode ends with this record. Each rule has one home; the [Bootstrap mode](../engineering-discipline.md#bootstrap-mode) section of `engineering-discipline.md` stays as history and lists the homes.

1. **Scope** (rule 1, without its end date). A task is a PSB In-Scope item (`F-0003#41`–`#52`) or a child of one, a defect that blocks such a task, or a documentation fix that a task leaves stale; its plan names the In-Scope fact it serves; no new process rule, routing rule, check script or policy capture starts on its own. Home: [Working a task under the quality gate](../engineering-discipline.md#working-a-task-under-the-quality-gate).
2. **The plan review** (rule 2) is one comment, by the Operator or by one fresh agent session, with `Verdict`, `Budget maximum` and `Cycle cap`. The goal-class count of R11 is not applied to a task whose deliverable is one artifact and its registration, and a `reject` on the count alone goes to the Operator, whose count is final. Homes: [R12](../issue-workflow.md#r12--slice-and-prioritize) and [R11](../issue-workflow.md#r11--single-goal-issues), which now answers [#49](https://github.com/pharzam/layup/issues/49) in the reading of O-15.
3. **The review** (rule 3) is one round on each frozen head, lens correctness and acceptance criteria, by a fresh session whose model differs from every author's; a second round runs only after a fix. The cycle cap is 1, and 2 when the change touches a gate (a check script, a hook, CI or branch protection). The material test is the Operator's text of O-43, word for word, and a note causes no round. "One pass is never enough" is removed, so the written rule and `review-record-lint` agree (F-15); the check does not change. Home: [Reviewing until findings decay](../engineering-discipline.md#reviewing-until-findings-decay).
4. **The end at the cap.** With a material finding open at the cap, the author posts the open findings on the issue, and the Operator chooses: one more cycle (a `## Plan review` comment records the new cap as a table row, and the last verdict is edited to `material` with the reason); a known limit or a finding ruled a note (the last verdict stays `not mergeable, findings recorded`, and the Operator's comment permits the merge); a fix with no further round (the same verdict, and the task record names the fix that no round read); or a split. Home: the same section.
5. **Reviewer selection** (rule 4, and part 3 of ADR-0012): the author and reviewer rule, the order `claude`, `devin`, `opencode`, and the skip after five minutes with no output or fifteen with no record. Home: [Who may review](../engineering-discipline.md#who-may-review), whose Model level is now required for every review. The tier names and the models not to use: [Model tiers](../engineering-discipline.md#model-tiers).
6. **A panel** (rule 5) is convened only for an ADR that changes the product architecture. Home: [Solution selection](../engineering-discipline.md#solution-selection).
7. **The resource record** (rule 6) and **one home per rule** (rule 7) stay as written. Homes: [Completing a task](../engineering-discipline.md#completing-a-task) and [One reading, not two](../engineering-discipline.md#one-reading-not-two).

Issue first, red then green, the frozen head, the record fields, evidence under `runs/`, the close-out, `review-record-lint`, the hooks, branch protection and CI do not change. The gate of this record applies to each plan reviewed after its merge; task `T-gd8q` ran under bootstrap mode.

We rejected, so that none is reopened without new information: **restore the full gate** as it was before ADR-0012 — the target task `T-vu2j` and the routing line of ADR-0012 show rounds that do not decay and end only by a decision at the cap, and F-15 would come back; **keep rules 2 to 7 as they are** — the split stays the only written end at the cap, while five of the seven stall tasks of bootstrap mode ended by another decision of the Operator; **add a review that looks for integration defects** (for example, a pilot per milestone) — the phase loop of [ADR-0019](0019-run-the-lifecycle-as-a-phase-loop-with-one-bet-per-milestone.md) already ends each milestone with a demo on a target, and the specification task of `M2a` is its place; **a panel for this record** — rule 5 allows a panel only for an ADR that changes the product architecture, and this record changes the process.

## Consequences

- The gate that delivered phase 1 stays, and its one undocumented part, the end at the cap, is now written. Fewer rounds still find fewer defects of one diff; the milestone demos, not more rounds, are where integration defects are found.
- The findings that the pilot did not fix get hosts in the [defect register](../plan/README.md#the-findings-of-the-first-pilot-that-are-not-fixed) of the plan; no code changes in this task.
- ADR-0012's Status is `Superseded by ADR-0026`. ADR-0005's Status names this record (the tier names and the models not to use move into `engineering-discipline.md`), and so does ADR-0006's (the panel rule). Nothing else in either record changes.
- Each summary of bootstrap mode (`AGENTS.md`, `README.md`, the onboarding guide, the glossary, the plan) is now history or a link to a home of this record; [`runs/T-gd8q/mentions.tsv`](../../runs/T-gd8q/mentions.tsv) classes each remaining mention.
- The count rule was written after the numbers. The next pilot pre-registers its thresholds first, as `guardrails.md` §1 asks.
