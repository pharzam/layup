# T-w79d — ADR-0012 (Proposed): route work across harness agents

Issue: [#46](https://github.com/pharzam/layup/issues/46). Child: [#47](https://github.com/pharzam/layup/issues/47) (`F-0005`). Record: [ADR-0012](../adr/0012-route-work-across-harness-agents.md). Evidence: `runs/T-w79d/` (inventory, clause map, check, test runs).

## Operator decisions (word for word; given in the session on 2026-09-24)

- **O-15** (the goal count of this task, after plan review 2). The Operator chose "One goal, keep the check": "Decision note: a decision record is one goal, as in #30, even with a check of its parts in the DoD. A later task makes the R11 text say so (R10). I fix the other conditions in revision 3." (#49 fixes the R11 text.)
- **O-16** (the boundary, after panel members A, B and C). The Operator chose "Keep; reviews conditional (Recommended)": "Keep the boundary as it is. ADR-0012 routes review rounds only as 'a binding that Who may review accepts', so X1 stays open. The quota states stay inside, and 'unknown' is never Normal (Invariant 5). No plan revision is needed."
- **O-17** (the scope). The Operator chose "This repository (Recommended)": "The gate work of this repository (LAYUP's own development). ADR-0012 amends ADR-0005 for it. The kit copy carries the rule to a target, but the product engine is not decided here."
- **O-18** (given when the Operator stopped the wait for member D): "for the later panel after current work done, max panel time should be 15 min, and each response timeout should not be more then 5min, blind harness can be same but model shuold be different , for fable 5.1 and ASTRA 6.1 in extream situation before the staleness status should be used"

## The panel (ADR-0006): the compared set

One pass each, no vote; the full comments are on #46. Member D (portability and evidence) gave no output on Kimi K3 (hung, stopped at 22 min) or on its substitute Qwen3.8 Max (300 s limit), and was dropped; three members remain.

| Option | Member | In short | Result |
| -- | -- | -- | -- |
| A1 | A, Claude Fable 5.1 on Claude Code | Harness axis beside the tier; bindings in an adopter table like O-3 | Selected: D1, D2 |
| A2 | A | F-0005's harness order as the rule in the ADR body | Rejected: Invariants 4 and 9 fail |
| A3 | A | Harness-neutral classes; a deterministic procedure picks the first eligible binding | Selected in part: D3 (the new class list is not taken) |
| A4 | A | The harness is recorded, not routed | Rejected: outside O-16; `F-0003#66` gets no rule |
| A5 | A | The harness is routed by gate role (author, reviewer, judge) | Taken in part: D7 (reviews by "Who may review") |
| B1 | B, GPT-6 Astra on Devin | Recorded priority, no percentage routing | Rejected alone: O-16 keeps the quota states |
| B2 | B | Quota states gated on vendor evidence; `unknown` never Normal | Selected: D6 (thresholds open, O-22) |
| B3 | B | Bindings ordered by measured cost | Rejected now: no comparable cost exists |
| C1 | C, Gemini 3.1 Pro on AGY | Failure detected on exit code and empty output; deterministic fallback | Selected: D4, with the order in the table, not in a script |
| C2 | C | A probe before each dispatch | Rejected: delay on every dispatch |
| C3 | C | Hard stop while a quota is unknown | Rejected: all work stops today |

**Disagreements recorded, not averaged.** (1) The boundary: A found the quota states separable, B found them inside; the Operator decided (O-16). (2) The state `unknown`: A kept the table order with `not reported`, B excluded unknown bindings from quota routing, C proposed "unknown = exhausted" or a hard stop; O-16 rules out "unknown = Normal", and the route in `unknown` is O-22. (3) C offered "waive Invariant 5" as an Operator option; not carried, because the invariants are frozen with the PSB (`guardrails.md` §1.1). (4) C found the boundary wrong because reviews depend on X1; A and B agreed on the dependency and kept the boundary with conditional review routing, which the Operator chose (O-16). No disputed factual claim remained, so no judge ran for the panel.

## Affected documents (step 6)

- `git grep -n -E "0012|next (constitutional )?ADR" -- . ':!docs/adr/' ':!runs/'`: the `F-0005` record (lines 32, 34) and `T-q1x6.md` (line 3) name ADR-0012, which now exists; still true.
- `AGENTS.md` "Model tiers" and the glossary row "Model tier" summarise ADR-0005; a Proposed ADR changes no rule; still true.
- Glossary: rows `AGY`, `Binding`, `GPT` added. Lesson: `guardrails.md` §2, "An agent command exits 0 with no answer, or never ends".
- The PR title check needs a lower-case subject, so a title that starts with the task ID fails (PR #50); it fails loudly, so it is recorded here and not in §2.
