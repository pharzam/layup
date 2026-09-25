# F-0004. The idea owner's answers to the gap batch of the PSB

| Field | Value |
| ------------ | ----- |
| Fact ID | `F-0004` |
| Source | The idea owner of LAYUP, who is the Operator (decision O-14 on [#42](https://github.com/pharzam/layup/issues/42)). Nineteen answers, one per question of the gap batch, accepted as one batch (Decision Point 2, `F-0001#11`). |
| Collected by | Claude Fable 5.1 (agent), for the Operator, task `T-zmj6` |
| Date collected | 2026-09-25 |
| Origin | The gap batch is [`internal/psb/testdata/psb.tsv`](../../internal/psb/testdata/psb.tsv), the output of `layup psb check docs/facts/problem-statement-brief.md` on `F-0001` (19 questions, `Q-001` to `Q-019`), pinned byte for byte by `TestGoldenRealPSB`. Fact N below answers the Nth question of the batch: the question IDs have three digits, so fact 1 answers `Q-001` and fact 19 answers `Q-019`. The answers were given in a Claude Code session on 2026-09-25: the agent drafted the nineteen answers with a rationale each and asked the idea owner to accept or edit them in one batch; the idea owner accepted all nineteen (see Notes on capture). |
| Status | `Raw` |

## Facts as collected

> One numbered fact per question of the batch, in the idea owner's accepted words.
> No source file exists for the answers, so this record is the evidence; check
> `facts` in `docs/setup/setup-check.sh` holds it to facts 1 to 19, each once.

1. Go, with the standard library only, and Git called as the `git` program. The PSB names no stack on purpose: LAYUP itself is built in Go (ADR-0010, ADR-0011), and each target project selects its own stack, which selects its stack-dependent gates (§6 In Scope, Invariant 7).
2. A short document that states one subject completely enough to act on. "Problem Statement Brief" is the document kind that §8 defines as PSB.
3. Quality Assurance: the engineering function that verifies that a delivered result meets its requirements and its tests. In §2 it names the role "QA Engineer".
4. Application Programming Interface: the contract by which one software component calls another. In Problem 2 "API contracts" are interface contracts between components.
5. Product Owner: the role that owns the product's requirements and priorities. In §4 "PO/PM" names the idea-owner side of the table.
6. Product Manager: the role that manages the product's scope and delivery from the business side. Together with PO it stands for the idea-owner side in §4.
7. Cost of Inaction: what the organization loses if the problem stays unsolved. §5 is the list of those costs.
8. Large Language Model: the kind of model that a harness runs. LAYUP does not modify or train one (§6 Out of Scope).
9. Pull Request: the forge's unit of a proposed change to the project repository, through which every change lands (Armature R1).
10. Rule R5 of Armature's issue workflow, "Deterministic over LLM-based": prefer a script, linter, type check or CI gate to an LLM judgement wherever a rule can be checked by a machine.
11. The rule for all §7.2 targets: the value printed in the table is the working hypothesis. The first pilot measures the baseline with the current process; then the idea owner sets each start value in one batch, records it in the pilot's PRD, and the pilot's report compares the measured trend with it.
12. Start value: median ≤ 50 % of the pilot baseline, as printed. Set by the idea owner after the pilot's baseline measurement (the rule of Q-011).
13. Start value: ≥ 90 % of delivered requirements accepted at the first review, as printed. Set by the idea owner after the baseline (Q-011).
14. Start value: ≤ 10 % of tasks with unplanned human input, as printed. Set by the idea owner after the baseline (Q-011).
15. Start value: ≤ 120 seconds at the 95th percentile, as printed. Set by the idea owner after the baseline (Q-011).
16. Start value: < 5 % of audited agent answers overturned within 30 days, as printed. Set by the idea owner after the baseline (Q-011).
17. Start values: ≤ 5 % of tasks stall, and ≥ 90 % of stalls close without human input, as printed. Set by the idea owner after the baseline (Q-011).
18. Start value: the median token cost per requirement of the pilot baseline itself (the target is "≤ the baseline"). Set by the idea owner after the baseline (Q-011).
19. Start value: ≥ 80 % of all human questions asked before delivery starts, as printed. Set by the idea owner after the baseline (Q-011).

## Notes on capture (optional)

- **How the answers were given.** The agent drafted the nineteen answers from
  the PSB (`F-0001`, `F-0003`), the vision brief (`F-0002`), ADR-0010 and
  ADR-0011, and Armature's issue workflow, with a rationale for each (kept in
  the task's evidence, `runs/T-zmj6/drafts.md`). It then asked the idea owner in
  the session on 2026-09-25: "Decision Point 2: the 19 drafted answers above
  (and in f-0004-drafts.md). Are they your answers to the gap batch?" The idea
  owner chose "Accept all 19 as my answers (Recommended)", whose text was: "I
  store them word for word as F-0004, facts 1–19, with the Origin row naming
  this session and the batch file at the merge SHA." This is the form ADR-0012
  part 5 sets: the author drafts all 19 answers, and the idea owner accepts or
  edits every answer explicitly in one batch before the answers become operative.
- **What the record claims.** The facts are the accepted words. The rationale
  of each draft is the agent's, not the idea owner's, and is not a fact.
- **Open gaps.** None: every question has an answer. Facts 11 to 19 defer the
  numeric start values to the first pilot's baseline measurement, as the PSB
  itself says (`F-0001 §7`); the pilot's PRD records the numbers when set.
- **Status.** Raw facts, the idea owner's words. A derived requirement may cite
  one by `F-0004#N`; which requirement cites which fact is decided where the
  requirements are written, not here.
