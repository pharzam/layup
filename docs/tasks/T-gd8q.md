# T-gd8q — the ADR that ends bootstrap mode

Issue: [#120](https://github.com/pharzam/layup/issues/120), opened by row 20 of the
[implementation plan](../plan/README.md) (`T-evad`, O-121). Serves the In-Scope work
of `F-0003#41`–`#52` that phase 2 continues. Base `34dd858` (the merge of #121).
Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-gd8q/`](../../runs/T-gd8q/).

## Plan and plan review

The plan (R12, comment 6025499388), its review and the author's answer
(6025660425) are comments on #120. The plan review (Claude Fable 5.1, effort
`xhigh`, on the Claude Code CLI, a fresh read-only session in a clone at `34dd858`,
8 min 20 s; comment 6025640705) gave `approve-with-conditions`: Budget maximum 900
lines added plus removed over 26 files against `34dd858`, close-out inside; Cycle
cap 1; no panel. Its ten conditions and ten notes are applied. A first start of the
review gave no record: the author's command line read the prompt as a tool name.

**O-159** (2026-10-07, the Operator's answer "1. A , 2. A" in the author's session,
copied to #120 as comment 6031681243): part 1 (a), the count rule is accepted,
knowing that it was written after the numbers; part 2 (a), ADR-0026 sets option C.

## What was done

1. **The count** (D1): [`numbers.md`](../../runs/T-gd8q/numbers.md), with its rule
   and two populations (the 20 rows of phase 1, and the other 10 tasks with a record
   under bootstrap mode), each value with a `file:line`;
   [`cites.sh`](../../runs/T-gd8q/cites.sh) prints each cited line at `34dd858`
   (116 lines, none missing, exit 0). One claim of
   the plan was wrong and the count corrects it: row 11, with one round, shares two
   pilot findings (F-16, F-25) with row 13.
2. **ADR-0026** (D2): [Keep the bootstrap review rules as the standing
   gate](../adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md),
   Accepted; ADR-0012 `Superseded by ADR-0026`; the Status lines of ADR-0005 and
   ADR-0006 and the index of `docs/adr/README.md` name it; the numbering line no
   longer names `0013`.
3. **The homes** (D3): in `docs/engineering-discipline.md`, the scope under the gate
   steps, the panel in Solution selection, the models not to use in Model tiers,
   one round per frozen head with cap 1 (2 for a gate), the four ends at the cap and
   the material test of O-43 in Reviewing until findings decay, the reviewer rule
   and the Model level for every review in Who may review, one home per rule under
   One reading, not two; `## Bootstrap mode` stays as history with a link to each
   home. In `docs/issue-workflow.md`, R11 holds "A task of one artifact" (the
   answer of #49) and R12 holds "One comment".
4. **The summaries** (D4): `AGENTS.md`, `README.md`, the onboarding guide, the
   glossary (Review lens, Cycle cap, Issue split, Ceiling, Material, First pilot,
   Bootstrap mode), `docs/guardrails.md`, `docs/spec/README.md`,
   `docs/spec/setup.md`, and the plan, which also hosts the findings that the pilot
   did not fix. [`mentions.tsv`](../../runs/T-gd8q/mentions.tsv) classes each
   remaining mention of the word: history, a link to ADR-0026, or another meaning.
5. **The checks** (D5): [`checks.sh`](../../runs/T-gd8q/checks.sh) failed at the
   base for C1 to C4 and passes on the head
   ([`test-runs.md`](../../runs/T-gd8q/test-runs.md)); a line with an unclassed
   mention makes C4 fail.

**The rejected alternatives** are in ADR-0026: restore the full gate; keep rules 2
to 7 as they are; add a review for integration defects; a panel.

**Known limits:** the count rule was written after the numbers, so the count is not
a blind test (O-159 part 1). Four tasks of P2 have no task file (`T-stfn`,
`T-b97r`, `T-vk3k`, `T-meh2`), and the record of `T-8ywj` round 2 does not class
its five findings. The Go comment of `internal/catalog/entries_integration_test.go`
still cites "Bootstrap mode rule 4" (code is out of scope).
