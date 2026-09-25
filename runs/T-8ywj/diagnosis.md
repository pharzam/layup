# Diagnosis: why LAYUP does not move, and the pivot

> **Status of this file.** This is the diagnosis as the author presented it to the
> Operator on 2026-09-25, before the decision. §1–§3 are the findings. §4 is the
> author's recommendation, which kept `F-0005`; §5 rejected the revert. The
> Operator chose otherwise: "Bootstrap mode + revert F-0005" (O-41) and the
> revised routing text as an evidence file (O-42). The decision is
> [ADR-0012](../../docs/adr/0012-build-layup-in-bootstrap-mode.md); where §4 or
> §5 differs from it, the ADR governs.

Audit of 2026-09-25, read-only, over `main` at `43a05c8`, branch `T-7wqc`, all 55 issues and their 229 comments, and the three audit reports beside this file (`audit/issues.md`, `audit/facts.md`, `audit/adr-0012-branch.md`).

## 1. The numbers

| Measure | Value |
| --- | --- |
| Calendar span of the repository | 2026-09-23 13:36 → 2026-09-24 10:37 (last merge, PR #50): 21 hours; then 25 hours with no merge |
| Commits on `main` (no merges) / completed tasks | 88 / 18 |
| Go code of the product | 424 lines in 8 files (268 non-test, 156 test); 2 commands: `layup version`, `layup psb check` |
| Product requirements (PRD) | none (`docs/prd/` holds the template only) |
| Completed tasks: product vs process | 5 vs 13; elapsed 386 min vs 774 min (1 : 2) |
| Open issues: product vs process | 4 vs 10; no product build has started since 2026-09-23 20:39 |
| The ADR-0012 line (#46 → #53 → #54, + #47, #51, #52, #55) | 9 plan reviews (5 `reject`), 11 review rounds, 9 frozen heads, 37 findings fixed, 1 panel, 13 clause-map sessions (1,019,459 tokens); about 3.4 h of reviewer runs that answered and 1.4 h that returned nothing; 18 failed external runs; branch `T-7wqc`: 1,416 lines, 0 Go lines, 77 % run evidence; `not mergeable` at the cap twice; paused (O-34); then #55: 2 plan reviews, both `reject`. Merged output of the whole line: `F-0005` (84 lines). |
| Operator decisions consumed by the routing question | 19 decided (O-15–O-18, O-26–O-40) plus 7 open inside the draft (O-19–O-25): 26 of the project's 40 decision numbers. Reviewer bindings changed four times in one day (O-18 → O-27 → O-30 → F-0006 §16 → O-39/O-40). |
| PSB facts that require a model-routing policy before the core engine | none (`facts-audit.md` §A) |
| F-0006 findings that describe the LAYUP product | 0 of 17 (`facts-audit.md` §B) |

## 2. Root causes

**RC1 — Process about process.** F-0005 and F-0006 tell *the agents that build LAYUP* which model reviews which comment. They do not describe LAYUP. ADR-0011 decision 8 says the engine makes no model call, so a routing matrix is not even part of the product's core. The ADR-0012 line spent one full day and 25 Operator decisions on a rule that changes no Go line.

**RC2 — The gate multiplies the work of a small task.** For a 356-line text file the rules need an issue, a sliced plan, an independent plan review by a reasoning model, red and green runs, a frozen head, two review rounds with different lenses and a different model, fixes, a resource record and a verdict: at least 4 LLM sessions and, observed for `F-0005`, 11 recorded runs. Five of the nine plan-review verdicts on the line turned on the R11 goal count; #49, opened on 09-24 to settle that count, has not started. The task of storing `F-0006` then hit R11 twice: two reviewers counted "at least 6" and "at least 4" goal classes where the identical `F-0005` task was counted as 1. The rule is unstable across reviewers, so a plan can be rejected for ever. That is the PSB's Problem 4 (a deadlock with no state and no limit), reproduced inside the process that builds the product that should solve it.

**RC3 — The gate depends on external services at review time.** 18 of the recorded external runs failed: "Upgrade to Pro", "Go usage limit exceeded", no output after 20 minutes, invalid keys. Each failure costs a decision from the Operator (O-27, O-30, O-31, O-39, O-40), which is the "human as message bus" of PSB Problem 1.

**RC4 — The kit is used for the wrong phase.** Armature is a discipline for a team that delivers a product with a spec, a budget and a pilot. LAYUP has no PRD, no pilot and one Operator; the machinery that LAYUP must automate (panels, routing, stall packages) is being run by hand on LAYUP itself, before LAYUP exists. The product path is clear and short — #45 (the idea owner answers 19 gap questions), then #42 (the PDR), then the `layup setup` command of #29 — and nothing on it has moved since 2026-09-23.

**Record defects the audit found** (evidence that the process is now losing its own state): three #46 comments hold only the author's footer — the reviewer records of plan review 3, round 1 and round 3 were never posted, and their findings survive only in replies; O-15 names two different decisions (#42 PRD scope; #46 goal count); `docs/tasks/backlog.md` on `main` is empty while six tasks are open; no resource record exists for #46, #53 or #54; and ADR-0012 at `30ca13f` cites `F-0005` in all 25 clause rows while `F-0006` §2 declares `F-0005` incomplete and §16 reverses O-27.

## 3. What stays, what stops

Keep (they are the PSB's invariants, and they cost little): Git as the record; issue first; test first; the deterministic checks (`setup-check.sh`, the linters, CI); one independent review of every change (PSB §7.1 Harness-Agent Neutrality: "each change gets ≥ 1 verification from a harness agent that did not make the change"); the resource record (model and time; tokens `not reported` when a harness gives none).

Stop, until LAYUP runs its first pilot: routing policies; panels for every ADR; "one pass is never enough"; goal-class counting on documentation tasks; the model-and-harness matrix; capturing operator policy texts as hashed facts.

## 4. The pivot: ADR-0012 "Build LAYUP in bootstrap mode"

One decision, made by the Operator as idea owner (a scope-against-date trade-off is business-forking, PSB `F-0001#13`, `#26`, so no panel selects it; ADR-0006 governs the solutions that agents select). The number 0012 is deliberate: `F-0005`'s record says "the ADR of task T-w79d (planned as ADR-0012) will assess it", and ADR-0012 now says why it does not.

**D1. Scope.** Until the first pilot (LAYUP sets up and drives one target repository), a task in this repository is one of: a PSB In-Scope item (`F-0003#41–#52`), a defect that blocks such a task, or a documentation fix that a task leaves stale. No new process rule, routing rule, or check is added unless a product task needs it.

**D2. The bootstrap gate** for this repository (the CI checks stay exactly as they are; `review-record-lint` still parses the record):
1. Issue with a goal and acceptance criteria (R1, R2, R11: one demo sentence; the goal-class count is not applied to a task whose deliverable is one artifact and its registration).
2. A short plan on the issue (steps, tests, budget). The plan review is one comment by the Operator or by one fresh agent session, with `Verdict`, `Budget maximum`, `Cycle cap`; a `reject` on scope alone is decided by the Operator, whose count is final.
3. Worktree, red then green, commit.
4. One review round (a second round only for Go code that changes a gate, or on request) by any model that is not the author's, on any harness that answers within 5 minutes; a harness that fails is skipped and recorded, never waited for, never decided about. Lens: correctness and acceptance criteria.
5. Close-out: verdict, resource record (model, effort, elapsed; tokens `not reported` when not given), completed-log line, PR, merge.

**D3. Routing, the whole rule.** Author: Claude Code, the best model it lists. Reviewer: any listed model on `claude`, `devin` or `opencode` whose model differs from the author's; first one that answers. Do not use: the models the Operator listed as not to use (F-0006 §14, quoted in the ADR). Nothing else is routed. This amends ADR-0005 (the tiers get concrete names: reasoning = Opus 5.5 / Fable 5.1 / GPT-6 Sol; execution = Sonnet 5 / SWE-2 / Luna) and ADR-0006 (a panel is convened only for an ADR that changes the product architecture).

**D4. The ADR-0012 line.** #46, #53, #54, #51, #52 and #55 are closed as "not needed before the pilot". Branch `T-7wqc` is kept as a tag `archive/adr-0012-draft` and not merged. `F-0005` stays as it is (immutable, `Raw`), with no assessment. `F-0006` is not captured as a fact: the text stays on #55 (verbatim, hash `96aeae82…`) and the lines the ADR uses are quoted in the ADR.

**D5. The product path, in order.** (1) #45: the Operator answers the 19 gap questions in one batch → `F-0004`. (2) #42: `PRD-0001`, one REQ per In-Scope item, traced to `F-0003#41–#52`, with acceptance criteria. (3) #29: `layup setup` (steps.tsv → a target), then `layup gate` (run the kit checks from outside), then telemetry (`runs/`), then the stall record. Each is one task under D2.

**Consequences.** Delivery starts within a day. The full gate returns when LAYUP itself can run it (that is the product). The record of what was cut is this ADR, so nothing is lost silently.

## 5. The alternatives, and why not

- **Revert `main` to before `F-0005` (PR #50).** Possible with `git revert -m 1 43a05c8` (no history rewrite), but it deletes hashed evidence for no gain: `F-0005` is 84 harmless lines, and its checks pass. The cost was never the fact; it was the ADR line built on it. Rejected.
- **Finish ADR-0012 as drafted.** 1,416 lines, two open findings, seven open Operator decisions (O-19 … O-25), and every clause keyed to `F-0005` line numbers; replacing the input by `F-0006` voids the 25-clause table and the 391-line clause map. Rejected.
- **Drop the gate completely.** Breaks Invariants 3, 5 and 6 and the PSB's own success criteria. Rejected.
