# Audit report — LAYUP at 43a05c8 (2026-09-25)

Clone: `/private/tmp/claude-501/-Users-farzam-projects-layup/26cc4c98-384d-4636-8768-8ceed05e8160/scratchpad/clone-audit`. Read-only; no file changed. Every quote is byte-exact from the tree.

## A. What LAYUP is

**Goal in the PSB's own words.** The PSB gives no goal sentence. Line 8: "This document states the problem only. It does not describe a solution, and it does not depend on a solution document." The problem (line 14): "Current multi-agent software development systems fail to realize expected velocity gains due to seven interdependent structural bottlenecks" (the seven are F-0003#1–#7). The scope is the twelve In-Scope items F-0003#41–#52, for example #41 "Problem Statement Quality: Detection of the gaps in a problem statement before delivery starts, so that the questions for the idea owner come in one batch and not one at a time.", #42 "Reproducible Discipline Setup", #50 "Cost Visibility", #51 "Specification Synthesis", #52 "Harness-Agent Neutrality". The Operator's own definition is ADR-0011 O-12: "The layup project final result is the product that uses armature as a kit to get a new problem statement brief and architectural vision brief from the idea owner and orchestrate the paths to deliver the solution for that problem statement."

**The 9 invariants** (F-0001#1–#9):

1. Git is the system of record.
2. The project repository is independent.
3. The agents that do the work cannot change the rules or the gates that check the work.
4. No configuration value without evidence.
5. A check that is not active does not count as passed.
6. A deterministic check is preferred to an LLM judgement.
7. The project domain changes content, never rules.
8. Armature is used at a pinned, recorded version.
9. A harness agent is replaceable.

**"Step 2 — the core engine."** The phrase exists only in `docs/adr/0011-structure-the-core-engine-as-a-go-cli-over-repository-files.md` (L11: "Step 2 builds the LAYUP core engine with no UI."; L113: "the stack-gate child of the Step 2 parent issue"). The parent issue is #29 (`docs/onboarding-for-engineers.md` L148–149: "The core engine is in progress under #29"). The PSB contains no "Step 2"; grep returns nothing. Completed children of #29: T-edtd (#30, ADR-0011), T-bhsf (#35), T-dq05 (#37), T-mtb9 (#38, carrying #32). ADR-0011 decision 8: "No LLM in the engine. The engine makes no model call; judgement stays with the harness agents that call it."

**PSB facts that require a harness-and-model routing policy before the core engine: none.** "routing" occurs once in the PSB, in F-0003#1 ("routing inquiries to specialist sub-agents", a symptom). "model" occurs three times: §2 unnumbered ("foundation models generate code rapidly in isolation"), F-0003#53 (out of scope: "Modifying base LLM weights or training custom foundation models"), F-0001#18 (a harness "gives the model, the session, the tools, and the context window"). The harness facts — F-0001#9, F-0001#14 ("The Operator can answer, or can give the task to a different harness agent"), F-0003#7, #26–#28, #40, #52, #66 ("The same project rules and gates run under $\ge 2$ harness agents. Each change gets $\ge 1$ verification from a harness agent that did not make the change.") — state neutrality and replaceability, not a routing policy. "Model Routing: Deterministically map tasks to optimal reasoning tiers and harness profiles" is F-0002 §3.1, and the F-0002 record says: "Use it later as input for ADRs, **not as requirements**. No PRD requirement cites `F-0002`." No PRD exists (`docs/prd/` holds README.md, prd-lint.sh, template.md only).

## B. F-0005 and F-0006

**F-0005** (`docs/facts/operator-routing-policy.md`, 84 lines, T-q1x6, 2026-09-24). L5: "You are the Autonomous Meta-Orchestrator for LAYUP." §1 hierarchy: Claude Code CLI > "Davinci CLI (Devin)" > AGY CLI > "OpenCoder CLI". §2A Claude Code: reasoning tier ("Planning, Solution Architecture, Blind Reviews, Panel Arbitration") = "High-reasoning frontier models"; execution tier = "Lowest-cost model satisfying acceptance criteria: Latest Haiku → Sonnet 5 → Opus 4.8"; fallback "Fable 5.1". §2B Devin/OpenCoder: free tier first, then paid; fallback "Astra GPT / Astra 6". §2C AGY: prompt preparation, option answering, execution. §3 quota bands <70% / 70–90% / >90% with daily and weekly tracking. §4 stall protocol: "Query all four registered harnesses" synchronously; a fixed A/B/C package (the fence opened at L64 never closes; L79 is "──────"). §5: "Reciprocal Verification (Invariant 9): A harness must never review its own output."; least-cost execution; commit "All prompt packages, multi-harness consultation notes, telemetry metrics, and operator decisions".

**F-0006** (`F-0006-draft-source.md`, 356 lines, uncommitted). Commands `claude`, `devin`, `opencode`; AGY removed. Exact model ID, effort and account alias per row (`claude-max`, `devin-pro`, `opencode-go`); nine eligibility conditions (§4); first-eligible-row rule with a 5-minute no-response limit (§5). Matrix:

- §6.1 suggest a solution: Opus 5.5 high → Devin gpt-6-sol-high → Devin claude-opus-5-5-high → grok-4.7#high.
- §6.2 plan: Opus 5.5 high → Devin claude-opus-5-5-xhigh → Devin gpt-6-sol-xhigh → grok-4.7#xhigh.
- §6.3 plan review: Devin gpt-6-sol-xhigh → grok-4.7#xhigh → glm-5.3#max.
- §6.4 PR review round 1: Devin gpt-6-sol-xhigh → gpt-6-luna#high → grok-4.7#xhigh → glm-5.3#max → qwen3.8-max#xhigh; round 2 after a fix: Opus 5.5 high (only when the author is not Claude) → Devin claude-opus-5-5-high → the OpenCode rows.
- §6.5 Judge: Fable 5.1 xhigh on `claude` → Devin claude-fable-5-1-xhigh → Devin gpt-6-astra-high; Claude author → gpt-6-astra-high only.
- §7 execution classes: Hard (Devin fusion Opus 5.5 + swe-2), Cheap mechanical (Devin gpt-6-luna-medium), Free (space-bunny-free), Ordinary (Devin swe-2-high → Sonnet 5 high → glm-5.3#high → kimi-k2.7-code).
- §9 limits 5/10/15 min; §10 red-flag issue with labels; §13 `docs/setup/routing-matrix.tsv`; §14 ban list; §16 supersession.

**Differences.** Global hierarchy → per-step ordered rows. Four harnesses → three. Model classes ("Latest Haiku", "Astra GPT / Astra 6") → exact IDs; Haiku and Opus 4.x banned. 70/90% bands with daily/weekly figures → four harness states and "Record an unavailable figure as `not provided`" (the kit's word is `not reported`, ADR-0007). Poll all four → skip an exhausted harness; A/B/C → "It can contain zero, one, two, or more options." Harness-level "Reciprocal Verification" → "A review model must differ from the author's model." (finding 15). Commit all prompt packages → "commit only a redacted manifest". New: Judge tier, execution classes, red-flag issue, TSV storage, benchmark evidence rule.

**Internal contradictions in F-0006.**

1. §8 rule 4 "Use `max` only for high-risk work involving security, data loss, or a costly or irreversible operation." vs §6.3 row 3 "Second lens: `opencode-go/glm-5.3#max`" (an ordinary plan review) and §6.4 round-1 row 4 "`opencode-go/glm-5.3#max`".
2. §6.4 "Use `xhigh` for the first blind round on a frozen head." and §8 rule 3 vs round-1 row 2 "`opencode-go/gpt-6-luna#high`" (and row 4 at `#max`).
3. §9 "In progress but no progress | 10 minutes | Stop the run and ask the Judge whether to continue or select another action." and "No useful time estimate | No unbounded wait | Ask the Judge." vs §6.5 "Judge: the last model step before the human", "Do not use Fable or Astra for ordinary work. They are last-resort Judge models.", and finding 5 "The Judge route is the last model route before human escalation."
4. Finding 1 "An LLM is a worker and does not enforce the route." vs §9, where the Judge (an LLM) may "select another action".
5. §10 requires "the `red-flag` label" and "a severity label", then states "The repository does not currently have these labels."; labels are forge state, outside Git (F-0001#1: "No project state and no decision is kept only outside the project repository").
6. §14 bans "Opus 5" while every Claude row uses `claude-opus-5-5`: two honest readings (R13).

**Product vs process in the 17 findings.** 0 of 17 state a behaviour of the LAYUP product (the `layup` command, its checks, setup, telemetry). 17 of 17 govern the process that builds LAYUP: they bind to the kit's gate steps (§6.2 "Specification plan", §6.4 "Round 2 after a fix", "first blind round on a frozen head"), to "the ADR-0007 resource record" (§8 rule 6), to this repository's labels (§10), to `docs/setup/routing-matrix.tsv` (§13) and to the Operator's accounts (§3; §14 "paid Zen models while that account has no funds"). Three restate PSB concepts without specifying the product: 12 (Stall, F-0001#37), 14 (the stall package, F-0001#14), 15 (corrects the reading of Invariant 9, F-0001#9). F-0005 L5 reads as a system prompt for the product; F-0006 has no such line.

## C. Product vs process ratio

**Span.** First commit d2516fd, 2026-09-23 13:36:52 +0300 ("chore: initialize from Armature kit at a959655"); last commit on main 43a05c8, 2026-09-24 10:37:49 +0300 (merge of PR #50); last non-merge 919e2f0, 2026-09-24 10:24:17. 88 non-merge commits, 18 completed tasks, 0 backlog lines, 21 hours of calendar time.

**Go code.** `cmd/` + `internal/`: 424 lines in 8 `.go` files; 268 non-test (`cmd/layup/main.go` 12, `internal/cli/cli.go` 59, `internal/psb/check.go` 197), 156 test lines. 7 `func Test…` in 4 files (1 e2e, 4 cli unit, 1 psb golden over 4 inputs, 1 psb integration). Root `tests/` holds README.md and .gitkeep only. `git diff --stat <root> main -- cmd internal tests go.mod`: 17 files, +518/−17.

**Classification** (elapsed = each task's resource-record "Total" row, every one written "about"):

| Task | Class | Delivered | Elapsed min | Tokens reported |
| ---- | ----- | --------- | ----------: | --------------: |
| T-edtd | PRODUCT | ADR-0011, core engine architecture | 80 | 881,509 |
| T-mtb9 | PRODUCT | `layup version`, Go gates (carries #32) | 117 | 766,476 |
| T-dq05 | PRODUCT | `layup psb check`, rules G1–G5 | 84 | 495,532 |
| T-fvwj | PRODUCT (borderline) | F-0001 PSB and F-0002 stored; check `facts` | 50 | 242,615 |
| T-ertw | PRODUCT (borderline) | F-0003 numbering "so PRD-0001 can cite them" | 55 | 351,454 |
| T-r7zg | PROCESS | Armature pin; `setup-check.sh` | 80 | 484,171 |
| T-vbwc | PROCESS | kit history removed | 70 | 424,770 |
| T-vpty | PROCESS | onboarding bound to PSB | 46 | 258,443 |
| T-xgz4 | PROCESS | glossary terms | 63 | 410,660 |
| T-7ndb | PROCESS | invariants into guardrails | 56 | 275,330 |
| T-nfh8 | PROCESS | markers filled; ADR-0010 (the stack choice) | 89 | 566,001 |
| T-q344 | PROCESS | CI workflows | 40 | 148,320 |
| T-fvng | PROCESS | lint scripts restored from main | 84 | 557,723 |
| T-afa5 | PROCESS | branch protection | 65 | 376,908 |
| T-6rg3 | PROCESS | README and AGENTS describe LAYUP | 30 | 136,971 |
| T-9mmm | PROCESS | setup procedure, `steps.tsv` (the spec of `layup setup` per ADR-0011) | 78 | 493,858 |
| T-bhsf | PROCESS | Go CI jobs required | 35 | 369,403 |
| T-q1x6 | PROCESS | F-0005 routing policy stored | 38 | 105,752 |

Counts: PRODUCT 5 (3 strict + 2 borderline fact captures of the product's PSB), PROCESS 13. Elapsed: PRODUCT 386 min (strict 281); PROCESS 774 min; total 1,160 min ≈ 19.3 h. Ratio 1:2 (strict 1:2.8). Tokens "reported" (reviewer sessions only; every task writes author tokens `not reported`): PRODUCT 2,737,586; PROCESS 4,608,310; total 7,345,896. The shared plan-review run (195,984 tokens, 12 min, "shared by #1–#12") is counted once, in T-r7zg.

## D. The rules that generate the work

Mandatory steps and artifacts for a small documentation task (store a 356-line text as F-0006), from the cited rules:

1. Issue before any commit (R1; R2 duplicate search).
2. Ordered plan on the issue (R12): DoD-covering, test-first, sliced by domain, with the solution-selection record (R3), one demo sentence and the goal-class count (R11), no vague word (R13).
3. One independent plan review by a fresh session or person (R12), reasoning tier (Model tiers table, ADR-0005), whose confirmation carries `Verdict`, `Budget maximum`, `Cycle cap`, the demo and the class count; CI `review-record-lint` parses it.
4. R7 decision comment before the commit.
5. Worktree `.worktree/<task>` off `origin/main` (step 1); read guardrails and ADRs (step 2).
6. Red then green (R8, step 3): a failing check run watched "for the right reason", then the code. For a fact: the F-NNNN file from `template.md`, the byte-identical source, the SHA-256 line in `facts.sha256`, the index row, check `facts`, a fixture under `docs/setup/tests`. T-q1x6 recorded 11 runs for this (`docs/tasks/T-q1x6.md`, "Test runs").
7. Freeze the head; review rounds (step 5): each a fresh session with Context, Method and Execution independence, plus Model independence for "a change to the checks themselves"; a different lens per round; a ten-field `## Review record — round N` comment; a clause-by-clause semantic pass for summary text; the R13 tripwire each round; a fix re-freezes; at most two cycles; the last round says `nothing material in scope`.
8. R9: a fresh context confirms the tests; then they are frozen.
9. Evidence under `runs/` (step 6); docs current and a guardrails §2 lesson (step 7).
10. Close-out commit (step 8): boxes ticked, verdict, completed-log line, ADR-0007 resource record; only the indexes and the detail file.
11. PR with `Closes #N`, subjects `<type>: T-xxxx …`, 12 required CI checks, plain merge; merge `origin/main` in, never rebase, after the freeze.

**Minimum LLM sessions.** Author 1 + plan reviewer 1 + review round 1 = 3, when round 1 ends `nothing material in scope` (T-6rg3, T-q344 and T-mtb9 closed on round 1). Read literally, "One pass is never enough. Each round catches a different class of error." (`docs/engineering-discipline.md`, Reviewing until findings decay) requires two rounds: 4 sessions. Observed for the 84-line F-0005: 4 sessions (Opus author; Gemini on AGY plan review; Fable round 1; Gemini round 2) plus "a 21 s run that produced no output", 38 min, 11 recorded runs, 6 commits, 1 fix cycle. Under F-0006 §6.3 ("Claude Opus must not review a plan written by Claude Opus") a second harness enters every task. A decision that becomes an ADR adds a panel: "required only for the architecturally-significant or novel decisions — in this project, each decision that becomes a new ADR" (Solution selection; ADR-0006), so ADR-0012 adds at least 3 panel sessions (T-edtd used 3, 284,035 tokens).

**Why R11 multiplies a task with several checkable items.** `docs/issue-workflow.md`, R11, bullet "Count the goal in the DoD, not only the demo":

> A single demo sentence can sit over a Definition of Done that enumerates several independent checks — one verb over many objects, "refuse when *a*, *b* or *c* is missing" — which reads as one goal and is several. The reviewer counts the **goal classes** the DoD names. Two items are the *same* class when one test failing for either fails for both, because they share a failure mode; they are *different* classes when either can fail while the other passes. A DoD of N independent classes is N goals: one child issue per class beyond the first, unless a single mechanism and a single test cover them all. A DoD that asks for a test per item — "a test asserts each *X* when removed" — has already named N failure modes, so the exception does not apply. The plan-review confirmation records the class count beside the demo.

R12 then adds: "**One gate per slice.** Every sliced sub-task passes the same quality gate; slicing is never a shortcut around discipline." Applied to T-q1x6's own verdict — "a changed byte, a removed file, a removed hash line and a removed index row each fail (rows 1–4 and 8)" — that is four failure modes tested one per item, so the exception does not apply and R11 counts four goals.
