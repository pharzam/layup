# Slice C review record

Reviewer: Claude Fable 5.1 (`claude-fable-5-1`), a fresh `claude -p` session in a disposable clone, read-only tools, no forge token, `--setting-sources project,local`. Not a gate round.

## Round 1

2026-09-29, head `5f9f40a`; 5 min 14 s; 23,148 output tokens; USD 3.32 at list price. Word for word:

~~~~text
# Slice C review — round 1
Reviewer: claude-fable-5-1. Head: 5f9f40a. Verdict: `material findings: 7`.

## Material findings

### M1. The line check can pass while it drops text
- Where: architecture.md:515-517 and ADR-0018:36-37: "every non-blank line of the file lies in a fact or in the 'not a fact' list". W-11:18 takes one sentence from inside a paragraph.
- Source: K19 and FT1 (the line check) are claimed; §7:517 says "no text is dropped without a record".
- Why it matters: Suppose a paragraph line reads "Invoices go out by e-mail. Refunds need a manager's approval." The session makes a fact of the first sentence only. If "lies in" means "meets a fact", the line passes and the refund need is lost with no record. If it means "is inside one fact", then the W-11 case fails, because it splits a paragraph. There are three more gaps. The "not a fact" entries are not stated to be byte-exact spans. A text that appears twice is matched by one fact. And §7's fail list for `layup spec check` (540-546) has no overlap or coverage rule, but W-11:19 names `layup spec check` as the actor.
- Fix: Give each fact and each "not a fact" entry byte offsets. The rule becomes: every non-whitespace byte lies in exactly one span. Name the command that runs this check. Say what `layup run` does when the check fails or does not run: no batch is posted.

### M2. The rule for who sets priority is only a name, and the baseline's own CI rejects the PRD
- Where: §7:531-538 (the session writes the PRD, then the idea owner sets priority); W-11:23-24; ADR-0018:45 ("no priority after the first bet").
- Source: K21; F-0003#54; FT3. The pinned baseline's `docs/prd/prd-lint.sh:115-121` fails a row whose MoSCoW is not Must/Should/Could/Won't, and a non-Won't row with no Phase. It runs as a baseline CI job in the target (§5 Scaffold step 4).
- Why it matters: The pull request of W-11 step 6 has no priority, so the target's own `prd-lint` job fails it. If the session fills MoSCoW to pass, a session has set intent. `layup spec check` only tests that a priority is present, so it passes. Nothing says who copies the idea owner's comment into the PRD. Nothing says where "marked out of scope" (§7:543) is read from.
- Fix: Code writes MoSCoW, and "out of scope", from the copied bet comment (from an ID in `approvers.tsv`, one line per requirement ID). `layup spec check` compares the PRD with that copy. Say how the Phase column is filled, or when it is filled.

### M3. The facts that `layup spec check` checks are the copy under review
- Where: §7:541-543 ("a fact is not a byte-exact substring of its file; a `need` fact has no requirement"). The rule-path list at §6:447-451 has no `docs/facts/`.
- Source: FT4; Invariant 3; L-C1 names the idea owner's confirmation as the guard.
- Why it matters: A role session's pull request edits the numbered-facts file. It reclasses a need as `context`, or deletes the fact. On that head the need no longer exists, so `layup spec check` passes. No rule path catches the edit.
- Fix: Check the head against the inventory the idea owner confirmed. Keep it on the records branch with its hash, or read it from the base branch. Or make `docs/facts/` a rule path.

### M4. The output and inputs of `layup spec check` are not defined
- Where: W-11:25: "a check table ... and the status `layup/gates`"; the input is "the target's facts and PRD".
- Source: §6:419-427. `layup/gates` is defined only from `layup gate` gate kinds, and "is a success only when every kind counts as a pass". Lens 1: a `code` step names its input and its rule.
- Why it matters: `layup spec check` is not a gate kind, so no row of the §6 table says how its failure, or its not running, maps to `layup/gates` (FT1). Its rules also read inputs that step 8 does not name: which requirement is "delivered", the tasks, and the idea owner's marks. "A specification section that names it" has no path and no naming rule, so code cannot find it.
- Fix: Add a row to §6, or give it its own required status such as `layup/spec` in the ruleset. List the records it reads. Name the file and the pattern of a specification section, for example a heading that holds the requirement ID.

### M5. A trace to the answers has no ID scheme, and it departs from the PSB wording without a recorded reading
- Where: §7:529-530 ("a requirement may cover an answer's fact ID"); ADR-0018:47.
- Source: K23. F-0003#62 says "a trace to the text of the problem statement". F-0001#39 says "a trace to the problem statement".
- Why it matters: No step numbers the answers or gives them IDs. No line check runs on them. So `covers: <answer ID>` has nothing to resolve against. A requirement that traces only to an answer also fails #62 as written.
- Fix: Code writes the answers fact with one ID per question ID. Each answer keeps the quote from the PSB that its question carries, so the trace reaches the PSB text. Record this as the reading of #62, the way §6 does for #64.

### M6. `constraint` and `measure` facts are classed, but nothing uses them
- Where: §7:518-520 ("One or more requirements for each `need` fact"); 542 (only an uncovered `need` fails); step 4 lists only needs.
- Source: F-0003#51 ("requirements and technical specifications" from the statement). The baseline PRD form carries NFRs, the way LAYUP's own NFR-001 to 006 trace to invariants.
- Why it matters: A legal line such as "invoices kept ten years" is classed `constraint`. It gets no requirement, no criterion and no check, and `layup spec check` passes. A `measure` has no path to the success criteria of §12.
- Fix: Say what each class feeds. For example, a constraint becomes an NFR with `covers` and is checked like a need; a measure feeds §12. Or say why they are `later`.

### M7. P14 (vision 2.1: PDR, PRD, phased plan) is answered only in part
- Where: §14:574 maps "table C 2.1" to W-11 and §7. W-11 has no plan step and no PDR.
- Source: rewrite-checklist P14. Vision 2.1 has "Preliminary Design Reviews" and "Phased Implementation Planning".
- Why it matters: The phased plan is neither answered nor marked `later`. The PDR is never named. The Phase column that prd-lint needs has no author.
- Fix: Name what the PDR is (the specification sections?). Mark the phased plan `later: §8` in the coverage row.

## Notes
- N1. L-C1 needs the first bet to show the idea owner the facts classed `context`/`constraint` with their text. Otherwise a need classed `context` cannot be seen at the confirmation.
- N2. L-C1: "until the idea owner's acceptance (§8)" — acceptance is per delivered requirement, so it cannot find a need that has no requirement.
- N3. §7 step 4 needs a rule for when only one harness is registered: stop, and record "not run". It should not quietly use the same harness.
- N4. W-11 step 6: a session reclasses a fact. Code should check that each class change cites an answer ID.
- N5. "Sets the priorities once, at the first bet" (ADR-0018:57): a requirement added after the first bet (from a later answer or a retrospective) has no path to a priority.
- N6. Nothing says where the classes are stored. The baseline facts template (`docs/facts/template.md`) has no class field.
- N7. The vision brief (`--vision`) is recorded at Intake, but it is not numbered or traced. Say whether it is a trace source.

## Checklist rows
- S11: answered in part (M2, M4 and M6 are still open).
- K19: answered, but the check is ambiguous (M1).
- K20: answered (completeness review on another harness; N3).
- K21: not answered — priority by name only (M2).
- K22: answered in intent; the rule cannot be run as written (M4).
- K23: not answered — no ID scheme (M5).
- P07: answered (§7 components, W-11).
- P14: not answered in part — PDR, phased plan (M7).
- FT1 (the line check): not answered — ambiguous rule, no failure path (M1, M4).
- L-C1: known limit (acceptable), with N1 and N2.

## Existing solutions
- OpenFastTrace (`Covers:`/`Needs:` tags, a deterministic coverage report), Doorstop and StrictDoc all do what `covers` plus `layup spec check` do. Doorstop's "suspect links" (a hash of the parent item) is the pattern that closes M3.
- GitHub Spec Kit's `/clarify` (asks the user questions) and `/analyze` (checks that requirements are covered across the artifacts) repeat steps 4 and 8. Spec Kitty is evaluated in ev06, but that evaluation does not cite them.
- BMAD's PRD splits FR from NFR. This design drops the NFR path (M6).~~~~

### The author's answer to round 1

All seven material findings and notes N1 to N7 are applied in one commit.

| Finding | Fix |
| ------- | --- |
| M1 | Facts and "not a fact" entries are spans with byte offsets; `layup spec check --facts`: every non-white-space byte in exactly one span; on a failure or no run, no batch is posted. |
| M2 | The PRD's MoSCoW and Phase stay empty until the bet; code writes them from the copied bet comment; `layup spec check` compares the PRD with that copy. |
| M3 | The confirmed inventory (facts, classes, marks, hash) lives on the records branch; `docs/facts/` is a rule path. |
| M4 | Own required status `layup/spec` (in the ruleset of §6 and ADR-0017); its inputs are listed; a spec section is a heading with the requirement ID in `docs/spec/`; not run = fail. |
| M5 | The answers fact has one ID per question ID; each question carries its PSB quote; the reading of `F-0003#62` is stated and goes to the approval brief. |
| M6 | A constraint becomes a non-functional requirement, checked like a need; a measure becomes a success criterion row of the PRD. |
| M7 | The PDR is named (the specification with the first bet's architecture); the phased plan is marked later: slice D. |
| N1 to N7 | The bet shows `context` facts with their text; L-C1 corrected; one harness → the review does not run and the run stops; a class change cites an answer ID; a later requirement gets its priority at the next bet; classes live in the confirmed inventory; the vision brief is no trace source. |

## Round 2

2026-09-29, head `30c4629`; 4 min 18 s; 19,985 output tokens; USD 2.45 at list price. Same reviewer and setup. Word for word:

~~~~text
# Slice C review — round 2
Reviewer: claude-fable-5-1. Head: 30c4629. Verdict: `material findings: 3`.

## Round-1 findings
- M1 (the line check): closed. Spans now have byte offsets, the rule is "exactly one span", the command is named, and a failed or skipped check stops the run.
- M2 (priority): closed. Code writes MoSCoW and Phase from the copied bet. `layup spec check` compares them. prd-lint passes only after that write, and the design says so.
- M3 (checking the copy under review): closed. The check reads the confirmed inventory on the records branch, and `docs/facts/` is a rule path.
- M4 (status and inputs): closed. `layup/spec` is its own required status, its inputs are listed, the section rule is stated, and a check that did not run fails.
- M5 (trace to an answer): closed. Each answer has an ID, and the reading of `#62` goes to the Operator.
- M6 (constraint and measure): closed for constraints, only half closed for measures. See new M2.
- M7 (P14): closed. The PDR is named and the phased plan is marked later: slice D.

## Material findings

### M1. The idea owner's "confirmation" has no way to reclass a fact
- Where: architecture.md:548-551: "shows the idea owner the needs, the constraints and the facts classed `context` … The idea owner answers one line per requirement ID: its priority … and its milestone, or 'out of scope' for a fact." W-11:24 has the same answer form. L-C1 (architecture.md:646-650) and ADR-0018:70-71 name "the idea owner's confirmation" as one of the two checks against a wrong class.
- Source: row K20 (Sol-18). L-C1 is accepted as a known limit only because two checks exist. Checklist lens 2 says a check answered by name only is a finding.
- Why it matters: the bet shows the idea owner the context fact "invoices must be sent by e-mail on the first working day", and the owner sees that it is a need. The only answer lines are a priority for a requirement ID, or "out of scope" for a fact. Neither turns a context fact into a need. So the confirmation cannot catch the error it is named for. The confirmed inventory then records the wrong class as confirmed.
- Fix: add an answer line such as "`need` for fact ID". Code then reopens the draft for that fact, as step 6 does after an answer. The spec pull request waits again. Or state that the bet may reject the list, and say what happens next (§8). Update L-C1 to match.

### M2. A measure fact can be dropped, and nothing records it
- Where: architecture.md:525-526: "Each `measure` fact becomes a success criterion row of the target's PRD, which §12 reads." The fail list at 565-570 covers only "a need or constraint". The completeness review at 530-531 lists only "each need or constraint that no requirement covers". The bet at 548-549 shows needs, constraints and context facts, but no measures.
- Source: `F-0003#51`; `F-0003#54` ("what counts as success" is the idea owner's decision); §5:313-314, where the intent form already asks the idea owner for "the success criteria". This is round-1 M6, half closed.
- Why it matters: the statement says "95% of invoices sent within one day", and the numbering session classes it `measure`. The spec session leaves it out of PRD §8. The span check passes, the review does not look at measures, the bet does not show them, and `layup spec check` passes. §12 never learns of the measure. There are also two sources of success criteria, the intent form and the measure rows. Nothing says which one §12 reads, and the session-written rows never reach the idea owner.
- Fix: make each measure fact need a success-criterion row that cites its ID, and have `layup spec check` fail when a measure has no row and no out-of-scope mark. Show the measures at the bet. Say how measure rows and the intent form's success criteria are merged, or mark the merge "later: slice G".

### M3. When the confirmed inventory is written, and how `docs/facts/` stays in step with it, is not defined
- Where: ADR-0018:48-52: "marks a fact out of scope, by one comment at each bet … The facts, classes and marks then form the confirmed inventory". architecture.md:558-559: "With the first bet, the numbered facts, their classes and the out-of-scope marks become the **confirmed inventory**". §5:341-342 puts "the numbered facts of §7" into `docs/facts/` at Scaffold. §6:447-451 makes `docs/facts/` a rule path in the register that is "written at setup".
- Source: an operative contradiction between the ADR and §7. `F-0003#42` (the facts are setup output). Invariant 2.
- Why it matters:
  - The inventory is only written with the first bet. At bet 3, the idea owner marks a fact out of scope, as ADR-0018 allows. Nothing writes that mark into the inventory, so `layup spec check` keeps failing the uncovered fact.
  - The W-11 reclass (step 6) comes after Scaffold has already written the numbered facts. The target's own `docs/facts/` still says `context`, while the inventory says `need`. `docs/facts/` is a rule path, and no step or batch updates it. So with LAYUP absent, the target's record of classes is wrong.
  - §7:558-560 reads as if `docs/facts/` becomes a rule path at the first bet. §6 puts it in the register at setup.
- Fix: say the inventory gets a new version at each bet, from that bet's copy, and store its hash. Say which class list Scaffold writes. Have the reclassed facts go into `docs/facts/` through a rule batch at the first bet, or state that `docs/facts/` holds only the unclassed statement and answers. Use one sentence for when the rule path starts.

## Notes
- N1. Nothing reads the inventory's "hash" (architecture.md:559). Either name the check that compares it, or drop it.
- N2. "Code checks the citation" (step 6) can only prove the answer ID exists, not that it causes the class change. Say that.
- N3. §7:538-541: the trace reading names only "a setup question" as having no quote. A review item "a requirement that no fact supports" also has no quote. Also, nothing makes code check that a requirement which covers an answer "also names that quote's fact".
- N4. The answers are one unclassed fact per question. If an answer states a new need ("also refunds"), no coverage rule applies to it.
- N5. Step 4 needs "a harness other than those of steps 1 and 3". With two harnesses, routing must put steps 1 and 3 on the same harness, and that is not stated.
- N6. "White space" is not defined at the byte level: ASCII only, or also multi-byte Unicode spaces such as U+00A0. The offset convention (start at 0, end exclusive) is also missing. The span check should also check that fact IDs are unique and classes are valid.
- N7. For slice D: the gate-activation batch is also at the first bet. Its `layup/spec` fails ("a need … has no requirement" on the head's PRD) unless the spec pull request merges first. The same holds for any pull request that adds a requirement between bets: prd-lint fails an empty MoSCoW until the next bet.
- N8. "The out-of-scope marks into the PRD" (step 7) has no place in the PRD form. The template's non-goals each need a `Won't` row, but an out-of-scope fact has no requirement.
- N9. W-11 inputs: step 4 also needs the span table (it lists "the fact whose class it doubts"). Step 8 also needs the span table and the step-6 reclasses to build the inventory. The answers that `covers` resolves against: records or head? Not stated.

## Checklist rows
- S11: answered, except measures (M2).
- K19: answered.
- K20: answered in part. The review on another harness is a mechanism; the idea owner's confirmation is by name only (M1).
- K21: answered.
- K22: answered. The set of delivered requirements is later: slice G, §12.
- K23: answered. The reading of `#62` is recorded for the Operator (N3).
- P07: answered.
- P14: answered. The phased plan is later: slice D.
- FT1 (the span check): answered.
- L-C1: known limit (acceptable) once M1 is fixed. Until then, one of its two checks cannot act.

## Existing solutions
- OpenFastTrace, Doorstop and StrictDoc cover `covers` plus `layup spec check`. Doorstop's per-item fingerprint ("suspect links") is the pattern for the inventory hash (N1) and for re-reviewing items after a later bet (M3).
- Standoff annotation is the pattern for byte-offset spans: the W3C Web Annotation `TextPositionSelector` and the brat standoff format. Their offset and whitespace conventions would close N6.
- GitHub Spec Kit's `/clarify` and `/analyze` repeat steps 4 and 9. Its spec template keeps "Success Criteria" as measurable rows next to the requirements, which is the check M2 asks for.~~~~

### The author's answer to round 2

All seven round-1 findings are closed. The three new material findings and notes N1 to N6, N8 and N9 are applied in the next commit; N7 is an input to slice D (the order of the spec pull request, the activation batch and `prd-lint` at the first bet). Under O-98, these fixes are a named item of the whole-design review.

| Finding | Fix |
| ------- | --- |
| M1 | The bet shows every fact with its class and text, and takes one line per fact: "out of scope" or a new class, which reopens the draft for that fact (§7 step 7; W-11 step 7; L-C1). |
| M2 | A measure needs a success-criterion row that cites its ID, or an out-of-scope mark; `layup spec check` fails otherwise; the review covers measures; the merge with the Intake form's success criteria is later: slice G. |
| M3 | A new inventory version, with its SHA-256, at each bet; the numbered facts record lands in `docs/facts/` through the bet's rule batch; the Scaffold writes only the statement and the answers; `docs/facts/` is a rule path from the setup. |
| N1 to N6, N8, N9 | `layup spec check` checks the inventory hash; "code checks that the ID exists"; the quote rule is checked, and review items are named; an answer's new need is classed as a fact; routing keeps steps 1 and 3 on one harness; white space, offsets, unique IDs and valid classes are defined; out-of-scope IDs go under the PRD's non-goals; W-11 inputs named. |
