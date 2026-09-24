# 0012. Route work across harness agents

Date: 2026-09-24

## Status

Proposed

## Context

[ADR-0005](0005-route-work-by-model-tier.md) routes model work by tier, and "Model tiers"
in [`engineering-discipline.md`](../engineering-discipline.md#model-tiers) fills the tiers
with Claude models (Operator decision O-3). No rule names a harness, although resource
records already name one ("Gemini 3.1 Pro on AGY"). The Operator gave a routing policy
over four harnesses (`F-0005`, input for ADRs, not requirements; lines 1–2 are the words of
the tool that generated it). Its 25 clauses C00 to C24 were checked one by one against the
tree: 15 conflict, 3 have no evidence, 6 extend the tree, 1 is not policy, and none only
restates it (`runs/T-w79d/clause-map.md`). The harnesses and models were measured on
2026-09-24 (`runs/T-w79d/harness-inventory.txt`): a smoke run proves that a binding ran,
not its price, quota or fitness. No harness reported how much of a quota was used; one
(OpenCode's Google provider) reported only that a free-tier quota with a limit of 0 was
exceeded.

The constraints: Git is the system of record (`F-0001#1`); no value without evidence
(`F-0001#4`); a check that is not active does not pass (`F-0001#5`); a deterministic check
before an LLM judgement (`F-0001#6`); a harness agent is replaceable (`F-0001#9`); a stall
reaches the Operator only at the limit of the stall procedure (`F-0001#14`); the PSB
pre-registers at least one verification by a harness agent that did not make the change
(`F-0003#66`). The vision brief (`F-0002` §2.2, §3.1, §3.5) was read as input only.

A panel of three members that differ in domain (governance fit; harness operations, cost
and quota evidence; failure and fallback) generated options, one pass each, with no vote
([ADR-0006](0006-convene-a-panel-to-generate-options.md)). A fourth member (portability
and evidence) produced no output on two bindings and was dropped. The compared set and the
recorded disagreements are in the task record of this decision. The Operator decided, on
2026-09-24:

| ID | Decision |
|----|----------|
| O-16 | Keep the boundary: this record decides which harness and which model run each class of gate work, with fallback and quota states; review rounds are routed only as "a binding that Who may review accepts"; "unknown" is never Normal. |
| O-17 | This record routes the gate work of this repository; the routing of the LAYUP product is not decided here. |
| O-18 | In the Operator's words: "for the later panel after current work done, max panel time should be 15 min, and each response timeout should not be more then 5min, blind harness can be same but model shuold be different , for fable 5.1 and ASTRA 6.1 in extream situation before the staleness status should be used" |

## Decision

This record is `Proposed`. If a later task accepts it, it amends ADR-0005 as follows.

### Proposed decisions

- **D1.** A route has two parts: the tier of ADR-0005, and a **binding** — a harness, an
  exact model, an effort setting and an account alias. Determinism still comes first.
  Then, in this order: the review levels (D7), the quota state (D6), the table order (D3).
  If no binding of the step's tier is eligible after these, D5 applies.
- **D2.** (no clause) The bindings are an adopter table in Git, not text in this record,
  so a product change edits the table, not an ADR; where the table lives is O-19. A
  binding enters the table only with a dated inventory entry that names its harness,
  model, effort and account alias; a model with no entry is not routable (`F-0001#4`).
- **D3.** The route is computed by a fixed rule, not chosen by a model. A row is
  **eligible** when its binding had a smoke run in the 24 hours before the dispatch (by UTC
  timestamps), has not failed in this task (D4), and the state of each of its quota pools
  allows the step (D6). The route takes the first eligible row of the step's tier, in
  table order, among the rows not marked `fallback`. The rule prints the binding; a
  harness agent launches it; the resource record names the model and the harness, for
  example "Claude Fable 5.1 on Claude Code".
- **D4.** (no clause) A dispatch fails on a non-zero exit, an empty standard output, a
  standard output of only white space, a standard output without its required shape, or no
  answer within its time limit. Standard error is never the answer. An exit code alone is
  not evidence: a print mode exited 0 with no output on 2026-09-24. After a failure the
  route runs D3, and then D5, again without the failed row.
- **D5.** If D3 finds no row, the route takes the first eligible row marked `fallback`, of
  any tier, in table order. If there is none, it takes the first eligible row of the other
  tier, in table order, among the rows not marked `fallback`, and the resource record
  names the tier it could not reach (ADR-0005). If there is none either, the route stops,
  and the author asks the Operator on the issue (R6). For a review step, each row it takes
  must also meet D7. The model names live in the rows, not here (D2); O-18 gives the
  Operator's proposed first values for the `fallback` rows (see the Context), and a name
  with no inventory binding is not routable (O-21). O-18 puts the fallbacks before any
  stall status; what gives that status is X2.
- **D6.** "Quota" is the vendor's quota, an input to the route, never a budget:
  ADR-0007's "recorded, not budgeted" stands. Each quota pool has one state, set by the
  last vendor report on it: a report that the quota is exceeded gives **Reserve
  exceeded**, with or without figures; a report with both a used amount and a limit gives
  **Normal**, **Constrained** or **Reserve exceeded** by the thresholds the Operator sets;
  no report, or any other report, gives **unknown**, never Normal (`F-0001#5`). Each state
  allows a set of steps: Normal, every step; Constrained, reasoning-tier steps only;
  Reserve exceeded, no step; unknown, the steps that O-22 allows. Each report and each
  change of state is an append-only record in Git (ADR-0011, 2).
- **D7.** (no clause) A plan review, a review round, a judge and a panel member take a
  binding that "Who may review" accepts. A blind reviewer may use the author's harness
  with a different model (O-18). Whether a different harness is required stays X1.
- **D8.** C23's prohibition is not adopted: ADR-0005's routing stands, and the execution
  tier owns execution-tier work. When a reasoning-tier model runs an execution-tier part,
  the resource record names it, and this record makes it no finding (ADR-0007 names only
  the reverse mismatch as one).

### Rejected options

- The Operator's order written into this record (Claude Code, Devin, AGY, OpenCoder): no
  evidence ranks it (`F-0001#4`), two names match no installed command, and a product
  change would need a new ADR (`F-0001#9`).
- An LLM "Meta-Orchestrator" that enforces the route: a machine can compute it (R5,
  `F-0001#6`), and the core engine makes no model call (ADR-0011, 8).
- A harness recorded but not routed: the Operator kept the routing in scope (O-16), and
  no rule would ask for the verification that `F-0003#66` pre-registers.
- Bindings ordered by measured cost: no comparable cost of accepted work exists; a model
  list gives unit prices only. Kept as a possible later order.
- A probe before each dispatch: its delay falls on every dispatch; the smoke run on the day
  of use and D4 detect the same failures.
- A hard stop while a quota is unknown: every harness is unknown today, so all work stops.
- The 70 % and 90 % bands read as a token or money budget: ADR-0007 has one budget, and
  its unit is not tokens.

### Open Operator decisions

- **O-19.** Harness names and table order: is "Davinci CLI (Devin)" the `devin` command
  and "OpenCoder CLI" the `opencode` command? Which criterion sets the order — a recorded
  preference, or a measured one? Where does the table live: beside O-3, or in a file
  under `docs/setup/` with a check?
- **O-20.** Which bindings join which tier: O-3 lists Claude models only. Options: keep
  O-3; add named models outside the current O-3 set, Claude or not, per tier (for example
  Claude Opus 4.8, "latest Haiku", an AGY or a Devin model); or make O-3 bind the
  author's route only and let "Who may review" pick a reviewer's model. Is "Anthropic
  models exclusively" on Claude Code a rule or a description? Which binding counts as
  "free tier": only one whose vendor charges nothing per token, or also one whose rate a
  prepaid allowance or a promotional credit covers? Does a paid fallback need
  authorization?
- **O-21.** Do "Astra 6.1" (O-18) and "Astra GPT / Astra 6" (F-0005 L29, L35) name one
  model or two? For each, which exact model and effort? Devin lists `gpt-6-astra-low` to
  `-max` and no 6.1; AGY and OpenCode list no Astra model. On AGY and OpenCode, is the
  fallback a change of harness to Devin?
- **O-22.** Quota and time figures: the thresholds (F-0005 gives 70 % and 90 %); the
  vendor's quota window, which can differ from the daily and weekly windows in which
  F-0005 L40 reports tokens and cost; the source of each figure and how recent it must
  be; the steps allowed in the state unknown; the time limit of a dispatch other than a
  panel response.
- **O-23.** "Panel arbitration": a reasoning-tier model that selects among a panel's
  options (Solution selection), or a panel that returns one verdict (ADR-0006 rejected it)?
- **O-24.** Prompt preparation, task scoping, option answering and trade-off scoring on
  AGY: does each pass Determinism first, and is "scoring" an input to Solution selection
  only?
- **O-25.** Prompt packages: are they committed in full, as a manifest with references,
  or not at all where they hold sensitive data (Safety limits)?

### Deferred items

- **X1.** Review independence by harness product (F-0005 L80–L82). It would change "Who
  may review". Its decision must reconcile O-18 (same harness, different model) with
  `F-0003#66` (a harness agent that did not make the change) and with F-0005 L21, L27
  and L82: L82 moves a review of Claude Code work to another harness, where L27 allows
  only non-Anthropic models on Devin and OpenCoder, and AGY lists Claude models too.
- **X2.** The stall protocol (F-0005 L52–L79). Where it changes the trigger or the
  hand-off of `F-0001#14`, it needs a PSB revision by the idea owner (`F-0001#23`);
  `F-0003#49`, `#61` and `#73` bound it. X2 also owns the retention of stall-consultation
  notes (F-0005 L84); `F-0003#49` requires at least a stall's diagnosis and outcome in
  Git.

### The clause table

| Clause | F-0005 lines | Relation | Basis | Destination |
| -- | -- | -- | -- | -- |
| C00 | L1–L3 | out of scope | The generator's words and a title state no rule (F-0005 L1, record Notes on capture). | not policy |
| C01 | L5–L6 | conflicts | A model enforces a route that a machine can compute; R5 and `F-0001#6` prefer a script, and no reason is recorded. | D3 |
| C02 | L11–L15 | no evidence | Nothing in Git ranks the harnesses (`F-0001#4`); two names match no installed command. | O-19 |
| C03 | L11 | extends | Reviews and quota go before the order; ADR-0005 already puts independence above routing. | D1 |
| C04 | L21 | extends | O-3 binds models, not a harness; a provider limit per harness is new. | O-20 |
| C05 | L22 | conflicts | "Panel arbitration" read as a verdict contradicts ADR-0006, which rejected a consensus verdict. | O-23 |
| C06 | L23–L24 | conflicts | Opus 4.8 and "latest Haiku" in the execution tier contradict O-3 (Sonnet 5, Haiku 4.5). | O-20 |
| C07 | L25 | extends | A fallback where the tree has none; Fable 5.1 is an O-3 model; O-18 keeps it. | D5 |
| C08 | L27–L28 | conflicts | Non-Anthropic models on Devin and OpenCode are outside O-3; no model is shown free. | O-20 |
| C09 | L29 | conflicts | "Astra GPT / Astra 6" is outside O-3; the inventory lists "GPT-6 Astra" on Devin only. | O-21 |
| C10 | L31–L32 | conflicts | Pruning and formatting on a reasoning model contradict Determinism and the tiers (engineering-discipline.md). | O-24 |
| C11 | L33 | conflicts | "Trade-off scoring" contradicts "considerations, not a scoring formula" (engineering-discipline.md). | O-24 |
| C12 | L34 | no evidence | No AGY execution model or price is in Git (`F-0001#4`). | O-20 |
| C13 | L35 | conflicts | The AGY fallback is outside O-3, and AGY lists no Astra model. | O-21 |
| C14 | L39–L40 | no evidence | Git holds tokens only per task part (ADR-0007), and no cost figure or daily or weekly rollup for any harness (`F-0001#4`). | O-22 |
| C15 | L42–L43 | conflicts | Read as a budget, "Quota / Budget" contradicts ADR-0007 ("recorded, not budgeted"). | D6 |
| C16 | L44–L46 | conflicts | Execution on Devin and AGY moves work to models outside O-3. | O-20 |
| C17 | L47–L48 | conflicts | Failover moves reasoning-tier steps off the O-3 reasoning models. | O-20 |
| C18 | L52–L54 | conflicts | "Requires Human Operator input" is not a stall cause in `F-0001#37`. | X2 |
| C19 | L56 | extends | Restates `F-0001#14` for a stall; new for every other question. | X2 |
| C20 | L57–L60 | conflicts | A poll of four named products makes the procedure depend on them (`F-0001#9`) and uses direct agent queries (R6). | X2 |
| C21 | L61–L79 | conflicts | The package goes to the Operator after every poll; `F-0001#14` sends it only at the limit. | X2 |
| C22 | L80, L82 | conflicts | "Must never" contradicts the ladder that stops and records a limit ("Who may review"); the label matches `F-0003#66`. | X1 |
| C23 | L83 | extends | Makes a prohibition of ADR-0005's routing (the execution tier owns routine edits); ADR-0007 names only the reverse mismatch as a finding. | D8 |
| C24 | L84 | extends | `F-0003#49` requires at least a stall's diagnosis and outcome in Git; prompt packages and the other consultation notes are new. Packages are routing work; notes go to X2. | O-25 |

## Consequences

- While this record is `Proposed`, no rule changes: "Model tiers" and O-3, "Who may
  review", the Human Decision Points and ADR-0007 stand, and ADR-0005's Status stays.
- On acceptance, a later task sets ADR-0005's Status to `Accepted. Amended by ADR-0007
  and ADR-0012`, writes the rule into "Model tiers" and the binding table where O-19
  decides, and updates every summary of both. D1 to D8 wait for the answers to O-19 to O-25.
- Nothing is mechanized. Until a later task writes the route rule as a script or a
  `layup` command, the route has no check (`F-0001#5`), and even then no check can prove
  which harness ran: the resource record stays self-reported (ADR-0007).
- Each harness costs upkeep: a dated inventory entry, a smoke run on the day of use, and
  route records in Git.
- Work sent to an external harness is sent to an external service (Safety limits). The
  Operator allowed it for the task of this decision; acceptance needs a standing rule.
- X1 and X2 each need their own decision; review routing stays conditional until X1.
