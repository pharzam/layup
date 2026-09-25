# ADR-0012 branch audit (T-7wqc @ 30ca13f vs main @ 43a05c8)

Read-only audit of the disposable clone at
`/private/tmp/claude-501/-Users-farzam-projects-layup/26cc4c98-384d-4636-8768-8ceed05e8160/scratchpad/clone-audit`.
Every file was read with `git show T-7wqc:<path>` or `git show main:<path>`; no tracked file was modified.

## 1. Size

`git diff --stat main T-7wqc`: 12 files, +1,416 / -1. `git log --oneline main..T-7wqc`: 24 commits (9 `T-w79d`, 9 `T-9tgn`, 6 `T-7wqc`).

| File | Lines |
|---|---|
| `docs/adr/0012-route-work-across-harness-agents.md` | +204 |
| `runs/T-w79d/clause-map.md` | +391 |
| `runs/T-w79d/clause-table-check.sh` | +297 |
| `runs/T-w79d/test-runs.txt` | +277 |
| `runs/T-w79d/harness-inventory.txt` | +120 |
| `docs/tasks/T-w79d.md` / `T-9tgn.md` / `T-7wqc.md` | +39 / +34 / +29 |
| `docs/guardrails.md` | +18 (two §2 pitfalls) |
| `docs/glossary.md` | +3 rows (AGY CLI, Binding, GPT) |
| `docs/adr/README.md` | +1 index row; "next constitutional ADR is `0013`" |
| `docs/tasks/backlog.md` | +2 |

77 % of the added lines (1,085) are run evidence under `runs/T-w79d/`. `cmd/`, `internal/`, `tests/`, `go.mod`: no change (`git diff --stat main T-7wqc -- cmd internal tests '*.go' go.mod go.sum` is empty).

## 2. The ADR draft

**Status:** `Proposed`. Date 2026-09-24.

**Context:** ADR-0005 routes by tier and O-3 fills the tiers with Claude models; no rule names a harness. The Operator gave F-0005 ("input for ADRs, not requirements; lines 1–2 are the words of the tool that generated it"). "Its 25 clauses C00 to C24 were checked one by one against the tree: 15 conflict, 3 have no evidence, 6 extend the tree, 1 is not policy, and none only restates it." Harnesses were measured on 2026-09-24; "No harness reported how much of a quota was used". A three-member panel generated options; a fourth "produced no output on two bindings and was dropped". Operator decisions: O-16 (keep the boundary; "unknown" is never Normal), O-17 (this repository's gate work only, not the LAYUP product), O-18 (verbatim: "max panel time should be 15 min ... blind harness can be same but model shuold be different , for fable 5.1 and ASTRA 6.1 in extream situation before the staleness status should be used").

**D1–D8, one line each:**
- D1: a route = ADR-0005 tier + a "binding"; order: Determinism, D7, D6, D3, then D5.
- D2 (no clause): bindings live in an adopter table in Git, each with an inventory entry; no entry, not routable.
- D3: a fixed rule; first eligible non-`fallback` row of the step's tier, in table order.
- D4 (no clause): what a failed dispatch is; after a failure rerun D3 then D5 without that row.
- D5: fallback rows of the step's tier, then any row of the other tier (record the tier not reached), else stop and ask the Operator (R6).
- D6: vendor quota is an input, never a budget; states Normal / Constrained / Reserve exceeded / unknown; append-only records.
- D7 (no clause): every review uses a binding "Who may review" accepts whose model differs from the author's; same harness allowed; a different harness is X1.
- D8: C23's prohibition not adopted; a reasoning model on execution work is no finding.

**Exact text of D1, D3, D5, D7** (ADR lines 50–53, 59–65, 72–80, 90–93):

> **D1.** A route has two parts: the tier of ADR-0005, and a **binding** — a harness, an exact model, an effort setting and an account alias. Determinism still comes first. Then, in this order: the review levels (D7), the quota state (D6), the table order (D3). If D3 finds no row, D5 applies.

> **D3.** The route is computed by a fixed rule, not chosen by a model. A row is **eligible** when its binding had a smoke run within the window that O-22 sets before the dispatch (by UTC timestamps), has not failed in this task (D4), and the state of each of its quota pools allows the step (D6). The route takes the first eligible row of the step's tier, in table order, among the rows not marked `fallback`. The rule prints the binding; a harness agent launches it; the resource record names the model and the harness, for example "Claude Fable 5.1 on Claude Code".

> **D5.** If D3 finds no row, the route takes the first eligible row marked `fallback` of the step's tier, in table order. If there is none, it takes the first eligible row of the other tier, in table order, whether it is marked `fallback` or not, and the resource record names the tier it could not reach (ADR-0005). If there is none either, the route stops, and the author asks the Operator on the issue (R6). For a review step, each row it takes must also meet D7. The model names live in the rows, not here (D2); O-18 gives the Operator's proposed first values for the `fallback` rows (see the Context), and a name with no inventory binding is not routable (O-21). O-18 puts the fallbacks before any stall status; what gives that status is X2.

> **D7.** (no clause) A plan review, a review round, a judge and a panel member take a binding that "Who may review" accepts and whose model differs from the author's, for every review, whatever its risk (O-18: "model shuold be different"); it may use the author's harness. Whether a different harness is required stays X1.

**Clause table:** header `| Clause | F-0005 lines | Relation | Basis | Destination |`, 25 rows C00–C24 at ADR lines 162–188. Every row carries a line span of `docs/facts/operator-routing-policy.md` (the 84-line raw F-0005 file, hashed in `facts.sha256`). Relations: 15 conflicts, 3 no evidence, 6 extends, 1 out of scope. Destinations: D1 (C03), D3 (C01), D5 (C07), D6 (C15), D8 (C23); O-19 (C02), O-20 (C04, C06, C08, C12, C16, C17), O-21 (C09, C13), O-22 (C14), O-23 (C05), O-24 (C10, C11), O-25 (C24); X1 (C22); X2 (C18–C21); "not policy" (C00). The spans come from `runs/T-w79d/clause-map.md`: checker Claude Opus 5.5, refuter Claude Fable 5.1, judge Gemini 3.1 Pro on AGY, 2026-09-24.

**Rejected options (7):** the Operator's harness order in the ADR; an LLM "Meta-Orchestrator"; harness recorded but not routed; bindings ordered by measured cost; a probe before each dispatch; a hard stop while quota is unknown ("every harness is unknown today, so all work stops"); the 70 %/90 % bands read as a budget.

**Open Operator decisions:** O-19 (names, table order, where the table lives), O-20 (which bindings join which tier; "free tier"; paid fallback authorization), O-21 (which model is "Astra 6.1" / "Astra GPT / Astra 6"), O-22 (thresholds, quota window, smoke-run window, steps allowed in `unknown`, dispatch time limit), O-23 (panel arbitration), O-24 (AGY prompt-prep roles), O-25 (prompt packages in Git).

**Deferred:** X1 (review independence by harness product, F-0005 L80–L82; #51), X2 (stall protocol, F-0005 L52–L79, L84; #52).

## 3. Task records and run evidence

**Resource records: none.** No `## Resource record` section exists in the three task files (every closed task on main has one, e.g. `docs/tasks/T-q1x6.md:34`). `T-9tgn.md:18` defers it to "the resource record at close-out"; no task closed out. **Tokens: zero figures anywhere on the branch.** `docs/tasks/completed.md` is unchanged.

**Elapsed time**, summed from the per-run seconds in the task tables and the inventory (not a resource record):
- Reviewer runs with output: #46 plan reviews 523 + 605 + 572 s; #46 rounds 1–2 604 + 578 s (round 3 not timed in the branch); clause-map judge 168 s; #53 plan review rev 2 603 s; #53 rounds 517 + 505 + 420 s; #54 plan review 364 s; #54 rounds 685 + 571 s. Sum about 6,715 s, 112 min.
- Runs that gave nothing: Kimi K3 1,200 s + ~600 s + ~1,380 s; Grok 4.7 641 + 191 s; Qwen 300 s; Gemini standalone 240 s; AGY auth 71 s; Sonnet 5 backup ~90 s. Sum about 4,700 s, 79 min.
- Author time (Claude Opus 5.5): not recorded.

**Review rounds:** #46 three, #53 three, #54 two (both `material`); the head commit is "record the final verification before the third freeze", so round 3 of #54 is pending. **8 rounds run**, plus 6 plan reviews (#46 ×3, #53 ×2, #54 ×1), a 3-member panel (+1 dropped), and the checker/refuter/judge workflow. No round on this line ended `nothing material in scope`; #46 and #53 each ended `not mergeable, findings recorded` at the cycle cap and were split.

**Models and harnesses:** author Claude Opus 5.5 on Claude Code; Claude Fable 5.1 on Claude Code (panel A, refuter); GPT-6 Astra xhigh on Devin (#46 plan reviews, panel B); Gemini 3.1 Pro high on AGY (panel C, judge); GPT-6 Sol xhigh on Devin (7 of the 8 recorded rounds — the binding of #46 round 3 is not in the branch files — plus the #53 rev-1 and #54 plan reviews); Claude Sonnet 5 on Claude Code (#53 rev-2 plan review, "an execution-tier model on a reasoning-tier part"); on OpenCode: Kimi K3, Qwen3.8 Max, Grok 4.7, Gemini 3.1 Pro Preview — every OpenCode dispatch of real work produced nothing. O-27 barred Fable 5.1, GPT-6 Astra and AGY for the T-9tgn session.

**Failed runs (18, from `harness-inventory.txt`):**
- "Upgrade to Pro to access this model": 5 (Devin gpt-6-sol-xhigh plan review of #53 plus 4 smoke runs; Devin "worked again at 11:41 UTC").
- Quota or credit: 2 (OpenCode anthropic/claude-opus-4-8 "credit balance is too low"; OpenCode Google "Quota exceeded ... limit: 0" — "the only quota figure a harness reported on this day").
- No output / hang: 9 (agy `--mode plan` exit 0 no output; AGY judge run 1 "authentication timed out"; panel D Kimi K3 stopped at 23 min; Qwen3.8 Max 300 s; Gemini standalone 240 s; Kimi K3 plan review 1,200 s exit 124; Kimi K3 round 1 ~10 min; Grok 4.7 ×2).
- Credential/config: 2 (openai/gpt-5.5 "Incorrect API key"; deepseek-v4-pro "requires Global regions").

**`clause-table-check.sh` (297 lines):** its header says "NOT a gate: no hook or CI job runs it, so it counts as no check (Invariant 5)" and "It checks the PRESENCE and the SHAPE of ADR-0012's clause table ... It does not check MEANING". Seven predicates: P1 one ADR file, header, GFM delimiter row; P2 each C00–C24 exactly once; P3 relation in a closed set; P4 each basis has a citation; P5 destination is D<n>/O-<n>/X<n>/rejected/"not policy" and defined; P6 F-0005 record present and `setup-check.sh --only facts` green; P7 the four Decision subsections have entries and each defined ID is used or says "(no clause)". About 130 lines are a CommonMark/GFM stripper (fences, HTML comments, pipe-less rows, escaped pipes). Self-test: 32 cases (29 mutations + 3 controls), green at d4f3c9c. Six of the eight review-round fix commits fix this script's Markdown parsing.

## 4. Does acceptance change the Go code?

**No.** The branch touches no Go file. O-17: "This record routes the gate work of this repository; the routing of the LAYUP product is not decided here." The Consequences, verbatim (ADR lines 192–204):

> - While this record is `Proposed`, no rule changes: "Model tiers" and O-3, "Who may review", the Human Decision Points and ADR-0007 stand, and ADR-0005's Status stays.
> - On acceptance, a later task sets ADR-0005's Status to `Accepted. Amended by ADR-0007 and ADR-0012`, writes the rule into "Model tiers" and the binding table where O-19 decides, and updates every summary of both. D1 to D8 wait for the answers to O-19 to O-25.
> - Nothing is mechanized. Until a later task writes the route rule as a script or a `layup` command, the route has no check (`F-0001#5`), and even then no check can prove which harness ran: the resource record stays self-reported (ADR-0007).
> - Each harness costs upkeep: an inventory entry with its UTC time, a smoke run within D3's window before each use, and route records in Git.
> - Work sent to an external harness is sent to an external service (Safety limits). The Operator allowed it for the task of this decision; acceptance needs a standing rule.
> - X1 and X2 each need their own decision; review routing stays conditional until X1.

It is entirely a rule for how humans and agents route gate work (plan reviews, review rounds, judges, panels) while building this repository. A `layup` command appears only as a possible later mechanization.

## 5. Dependence on F-0005 line numbers

Yes, throughout. Of the ADR's 204 lines, 37 cite `F-0005`, an `L<n>` span or a `C<nn>` ID: the 27 lines of the clause table (all 25 rows carry an F-0005 span), Context lines 15–16, D8 (C23), O-21, O-22 (2 lines), X1 (3), X2 (2). The check script ties to F-0005 in about 20 of 297 lines: the header string (25, 243), the cite regex `F-0005 L[0-9]+` (145), P6 (170–179: glob `F-0005-*.md`, the `operator-routing-policy.md` hash line), mutations 6a (edits raw line 84 "repository."), 6b, 6c, and cases 5e, k3, 2e, which embed "L5–L6" / "L84" (6 of 32 self-test cases). Glossary rows AGY and GPT cite F-0005 L14 and L29. The clause map (391 lines) is keyed to F-0005 spans in every section.

**If F-0006 replaces F-0005:**
- Same text re-issued (line numbers unchanged): a rename. About 37 ADR lines, 20 script lines, 2 glossary lines, the clause-map header — under 70 lines, mostly `F-0005` to `F-0006` and a new hash line. The 32 self-test cases stay valid (6a's line 84 needs a re-check).
- Different text: the 25 spans are void, so all 25 rows and their relations and bases must be re-derived (the 391-line clause map and its checker/refuter/judge runs); the destinations that feed D1, D3, D5, D6, D8, O-19–O-25, X1 and X2 must be re-checked; the script's P2 (`C00..C24`) and cases 2c/2d/2e/2f ("C25"), 5d (C00 = "not policy"), 5e/k3 change if the clause count changes. Roughly 37 ADR lines, 391 clause-map lines, 25 script lines, 2 glossary rows, and the whole clause-mapping workflow again.

## 6. Is there a 20-line deterministic core?

Yes, and the draft states it in its own words:
- D3: "The route is computed by a fixed rule, not chosen by a model." ... "The route takes the first eligible row of the step's tier, in table order, among the rows not marked `fallback`."
- D2: "The bindings are an adopter table in Git, not text in this record" ... "a model with no entry is not routable".
- D7: "a binding that 'Who may review' accepts and whose model differs from the author's, for every review, whatever its risk".
- D6: "a report that the quota is exceeded gives **Reserve exceeded**, with or without figures" ... "Reserve exceeded, no step".
- D4: "A dispatch fails on a non-zero exit, an empty standard output, a standard output of only white space, a standard output without the shape that the step's brief requires (for a review, its record heading), or no answer within its time limit." ... "After a failure the route runs D3, and then D5, again without the failed row."
- D5: "If D3 finds no row, the route takes the first eligible row marked `fallback` of the step's tier ... If there is none, it takes the first eligible row of the other tier ... If there is none either, the route stops, and the author asks the Operator on the issue (R6)."
- D1: "Determinism still comes first. Then, in this order: the review levels (D7), the quota state (D6), the table order (D3)."

That is one TSV (harness, model, effort, account, tier, fallback flag, last smoke-run UTC, quota state) and one loop: first row of the tier that is not fallback, smoke-run inside the window, not failed this task, quota state allows, model differs from the author's for a review; else first fallback row of the tier; else first row of the other tier and record the mismatch; else stop and ask. On a failed dispatch (bad exit, empty stdout, missing heading, timeout) mark the row and repeat.

Two facts limit that core today. First, every numeric parameter is open: the smoke-run window, the thresholds, the dispatch time limit and the steps allowed in `unknown` are all O-22; where the table lives is O-19; which bindings join which tier is O-20. Second, D6 makes every pool `unknown` today ("No harness reported how much of a quota was used"), and `unknown` allows only "the steps that O-22 allows" — so until O-22 is answered, the rule as written routes nothing, which is the outcome the ADR itself rejects under "A hard stop while a quota is unknown". A 20-line core can stand if the Operator sets one window, one time limit, and "unknown allows every step until a vendor reports otherwise" (D6 already forbids only `unknown = Normal`).

Everything else in the 1,416 lines is either evidence that F-0005 conflicts with the tree (the clause table and map), a Markdown-shape checker with no gate, or the record of the harness failures in §3.

## Sources read

In the clone: `git show T-7wqc:docs/adr/0012-route-work-across-harness-agents.md`, `docs/tasks/T-7wqc.md`, `docs/tasks/T-9tgn.md`, `docs/tasks/T-w79d.md`, `runs/T-w79d/clause-table-check.sh`, `runs/T-w79d/clause-map.md`, `runs/T-w79d/harness-inventory.txt`, `runs/T-w79d/test-runs.txt`; `git diff main T-7wqc` for `docs/adr/README.md`, `docs/glossary.md`, `docs/guardrails.md`, `docs/tasks/backlog.md`; `git show main:docs/facts/operator-routing-policy.md`, `main:docs/facts/F-0005-operator-routing-policy.md`, `main:docs/tasks/T-q1x6.md` (resource-record shape).
