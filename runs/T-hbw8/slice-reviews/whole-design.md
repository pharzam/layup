# Whole-design review record (plan v2, step 11)

Reviewer: Claude Fable 5.1 (`claude-fable-5-1`), chosen by rule 4's order (`claude` first; not an author, O-87), a fresh `claude -p` session in a disposable clone, read-only tools, no forge token, `--setting-sources project,local`. Lens: two cases that cross slices, the last-round fixes of every slice (O-98), contradictions, the six failure tests, existing solutions, coverage. Not a gate round.

## Round 1

2026-09-29, head `4ad186c`; 13 min 35 s; 65,196 output tokens; USD 23.50 at list price. Word for word:

~~~~text
# Whole-design review — round 1
Reviewer: claude-fable-5-1. Head: 4ad186c. Verdict: `material findings: 19`.

## Case 1
1. Start: `--new`, the root commit, and the records branch with the lease (§5 Start, §2, §3; ADR-0013, ADR-0014). Mostly OK. `lease.H` has no value yet: M8.
2. Probe sessions and Intake sessions before the answers arrive (§9, §5, §7 steps 1–4). They have no cap, no wall-clock limit and no band: M8.
3. The Intake batch and the answers (§5 gap check; §7; ADR-0018). OK.
4. Scaffold, `layup setup verify` and the rulesets (§5, §6; ADR-0016, ADR-0017). The checks `ci` and `procedure` fail on a correct target: M7.
5. Shape tasks: the specification, the architecture and the activation batch (§8, §7 step 6, §6). None of them can pass `layup/spec` before the bet (M1). The batch's gate run is not consistent (M3). The batch may be refused by layer 2 (M4).
6. The first bet (§8 Bet; ADR-0019). The merge order starts with the specification, and the specification cannot merge: M1. A class change blocks every merge: M2.
7. The rendered MoSCoW and facts-record tasks (§7 step 8, §8). Refused by layer 2: M4.
8. The milestone plan (§8 Build). OK.
9. The build task loop: plan, review, test first, draft PR, gates, verifier, close-out, merge (§8, §9; ADR-0019, ADR-0020). The first attempt is OK. A rework attempt breaks the frozen list (M11). A base merge carries the verification wrongly (M12).
10. The budget check before each session, inside a milestone (§12; ADR-0024). OK.
11. The acceptance of the requirement (§8 Accept, §12 `acceptance.tsv`). The mechanism is OK. No walkthrough walks it: M19.
12. The audit (§12; an empty sample is "not measured"). OK.
13. `layup learn` and the routing update (§13; ADR-0025). OK. A route with one unknown-money task freezes: M16.
14. The lesson and its rule batch (§13, §6). A stale known-bad patch blocks the batch (M6). A batch that changes workflows cannot run the loop (M5).

## Case 2
1. The developer session ends `needs_context` (§8; ADR-0019 d5). OK.
2. `layup run` records the question, posts it and asks P2 (§8, §10). With `stall.N` = 1 the question round may already be a stall: M9.
3. An owner-role session is routed (§9 owner map; ADR-0020). OK.
4. The owner session hangs. Trigger 3 fires and the session is killed (§11; ADR-0023 d2). OK.
5. The package and the diagnosis (§11 steps 1–2). OK.
6. The retry rung (§11 step 3). It does not say which session is retried: M10.
7. A second hang moves one rung up, and the panel is skipped (§11 step 3). OK.
8. The Operator answers `answer:` or `reroute` (§11 step 5). Neither form acts on the open question or on the owner role: M10.
9. The outcome, the question row and the measures (§11 step 6, §12). The question row is never closed: M10.

## The last-round fixes
- A round 3 `b63e090`: closed (notes only).
- O-95 `f32a732`: closed (note N1).
- B round 2 `ce87f71`: M7; also M19 (K11).
- C round 2 `00a3f4d`: M2, M4.
- D round 2 `2b4953a`: M1, M3, M12, M13.
- E round 2 `f17bfcd`: M14.
- F round 2 `8e3ba47`: M9, M11, M15.
- G round 2 `d3a348e`: closed (note N2).
- H round 2 `098d565`: M3, M6, M16.
- Step 10 `4ad186c`: M17, M18.

## Material findings
### M1. No Shape task can pass `layup/spec` before the first bet, and the specification cannot merge after it
- Where: architecture.md:588 "It reads the latest confirmed inventory". :578 "At each bet, code writes a new version of it". :595 fails when "a MoSCoW or Phase value in the PRD differs from the bet copy", "or it did not run". :681-682 "A failure goes back to the developer … in a new attempt". :648 "Merge order after a bet: the specification, the facts record …". W-11:25 "the verified specification's head is not touched". :571-574 "only then does the pull request pass the baseline's own `prd-lint`".
- Source: FT1; ADR-0019 d2 (the bet approves verified batches).
- Why it matters: before bet 1 there is no inventory and no bet copy. So the specification, architecture and batch tasks fail `layup/spec`, loop through attempts and never reach a verifier. After the bet, the specification head still has empty MoSCoW columns, so `prd-lint` and `layup/spec` fail. The rendered MoSCoW task edits a PRD that only the unmerged specification contains.
- Fix: say what `layup spec check` checks before the first inventory exists. Put the rendered columns on the specification head, or merge them as one change.

### M2. A new class at a bet blocks every pull request until the next bet
- Where: :569-570 "A new class reopens the draft … and its requirement waits for the next bet line". :592-594 fails when "a need or constraint has no requirement … and it has no out-of-scope mark".
- Source: F-0003#51; the merge order (:648).
- Why it matters: bet 1 reclasses F-12 as `need`. Every head of milestone 1 now fails the required `layup/spec`. The PR that adds the requirement fails `prd-lint`, because its MoSCoW stays empty until bet 2. Bet 2 needs milestone 1 accepted, so nothing merges.
- Fix: exempt a need whose requirement waits for a bet line, or let the same bet set its priority.

### M3. The rule batch's gate run contradicts itself in three places
- Where: :720 "developer (rule batch) | verifier | … a known-bad patch per kind it activates; `layup gate` passes on its head and fails on each patch". :429-430 it "never runs the head's gate files … with one exception: an approved rule batch". ADR-0016:51-52 "except those of an approved activation batch". :435-438 "one new known-bad patch per kind it touches".
- Source: FT4; F-0003#64.
- Why it matters: the handoff to the verifier comes before approval. With the base's files, a `pending` kind cannot fail on its patch, so the activation batch is never "verified" for the bet. If code runs the head's files anyway, FT4 breaks. A retrospective batch that changes a kind (without activating it) needs no new patch at the handoff. The ADR would forbid running its gate files at all.
- Fix: one rule, stated in §6, §8 and ADR-0016: which gate files run before and after approval, and "per kind it touches".

### M4. Layer 2 refuses the batch and the rendered tasks that §8 makes tasks
- Where: :475-478 "Before `layup run` pushes a task branch … A change to a rule path refuses the result". :582 the facts record "lands in `docs/facts/` as its own task". :618-619 "each rule batch, each rendered record … are tasks". Only `layup/rules` (:486-490) has the exceptions.
- Source: ADR-0017 d2, d4.
- Why it matters: at bet 1, the facts-record task and the activation batch are both refused before the push.
- Fix: give layer 2 the same two exceptions as `layup/rules`: an approved batch and a rendered record.

### M5. A batch that changes workflows cannot go through the task loop before its approval
- Where: :463-464 "`layup run` pushes it to a branch `batch/<point>` and opens its pull request, so the approver sees the change". :1274-1275 "its sessions write, review and verify it before the brief". :243-246 the push is refused for `.github/workflows/`. W-13:25 shows "the batch's pull request" for W-03's CI-job change.
- Source: O-92, O-93.
- Why it matters: there is no PR, no CI run and no `layup/verify` status (a status needs a commit on the forge). So the brief cannot name a verified head.
- Fix: state the order for a workflows batch: who pushes it and when, and when CI and verification run.

### M6. A known-bad patch that no longer applies blocks every later batch of its kind
- Where: :438-439 "a patch that passes, or that no longer applies, refuses the merge". :1283-1284 "and the batch goes back as a task".
- Source: F-0003#64; §11:1088-1090 (the early retrospective fixes a wrong gate).
- Why it matters: milestone 3 renames the package that the boundary patch edits. From then on, every batch that touches the boundary kind is refused, including the batch that fixes a wrong gate. Nothing retires or re-bases a recorded patch.
- Fix: a way to replace a patch in an approved batch, with its detection run again.

### M7. `layup setup verify` fails `ci` and `procedure` on a correct target
- Where: :331-334 "does each check of `setup-check.sh` from outside, for this target (its name, its facts, its pin), and the evidence of steps S04 to S15 becomes "`layup setup verify <check>` OK"". :351 "and the others of `setup-check.sh`".
- Source: setup-check.sh:448-449 (check `ci` needs `ci.yml` to run `setup-check.sh`); :582 (check `procedure` needs `docs/setup/steps.tsv`); :346-347 "Step S12 adds no `setup-check` job"; FT5.
- Why it matters: the exit code differs from the documented "OK". To make S15 pass, `steps.tsv` would have to go into the target, which breaks FT5.
- Fix: list the checks that are adapted or dropped for a target, and change S15.

### M8. No budget or limit acts before the first bet
- Where: :1149 "A milestone's clock starts at its bet". :1151-1152 "The project total is the sum over all milestones". :1154 "code sums the milestone's spend". :979-980 `lease.H` and `harness.<id>.cap`/`.wall` are "set at Intake".
- Source: FT2; O-68; Problem 5 ("a retry loop can consume a budget before a human sees it").
- Why it matters: the Intake sessions and every Shape task run with no milestone, no cap and no wall-clock limit. Their spend is left out of the project total that is compared with `B` and `U`. The lease row is written at Start, before `lease.H` has a value.
- Fix: count pre-bet spend in the project total, and give Intake and Shape a cap. Give the pre-answer parameters defaults with evidence.

### M9. The stall rules for questions contradict §8 and W-06
- Where: :746-747 "An attempt that a question ends does not count toward the attempt limits of §11". W-06:22 "an attempt row that does not count toward §11's limits". ADR-0023:76. All three against :1014 "(attempts that a question ended included)". Also :1002-1004 and :1012 "A question round counts the same way".
- Source: F round 2 M2; O-82 (`stall.N` = 1).
- Why it matters: an implementer who follows §8 leaves question loops unbounded. Only the first review sets a baseline, so a round that ends with a question closes nothing. At `stall.N` = 1, W-06's first question is already a stall.
- Fix: one rule for question attempts, and a stated baseline for a question round.

### M10. A stall of an owner-role session has no action that fits it
- Where: :1049-1050 "a **retry** with the diagnosis in the prompt". :1063-1064 "goes into the task's next attempt". :1068-1073 `answer: <text>` has no stated effect. "A reroute writes a harness override for the task, and the next attempt starts from the base".
- Source: F-0001#14; F-0003#46, #71, #75.
- Why it matters (Case 2): it is not said whether the retry starts a new owner session or a new developer attempt, which would re-ask the question. `reroute` moves the developer, not the owner role. The Operator's `answer:` never closes the question row, so it is never answered or accepted. Its place in Turnaround and in Early Question Share is undefined.
- Fix: say which session each rung and each answer form acts on. Make `answer:` the answer of the open question, with its measures.

### M11. The frozen test hash breaks on every later attempt
- Where: :999-1001 "a test whose source hash differs is no longer counted as passed, and the change is a finding; an amended plan with a new plan review freezes the list again". :749-750 "A new attempt starts from the base commit … not its diff".
- Source: F round 2 M1; K42.
- Why it matters: after a verifier finding, attempt 2 rewrites the test from the base, so its hash differs. That opens a finding that only a re-plan closes. No step-table row amends a plan, and the new plan review comes before the new test exists.
- Fix: keep the frozen tests as a payload that each attempt starts from, or define who re-plans and when the list is re-frozen.

### M12. A clean base merge carries `layup/verify` against the target's own rule
- Where: :656-659 "A merge with no conflict whose gates pass carries `layup/verify` … as the target's own rule lets a merge of the base land".
- Source: engineering-discipline.md:313-315: no new round only "while it is clean and changes no file the branch touched". review-record-lint.sh RR8.
- Why it matters: T-1 and T-2 edit `total.go` in different hunks. Git finds no conflict, and no round runs, which the rule forbids. On a conflict, the new attempt's round follows a `nothing material in scope` round, so RR8 fails.
- Fix: follow the target's rule (a scoped round when a touched file changes). Say how the record continues after a conflict.

### M13. W-12 lets H1 verify a change built from its own diff
- Where: W-12:25 "with the diff of attempt 1 as a payload". W-12:28 "its only author is H2, and H1 is admitted".
- Source: :749-752 and :797-798 (a harness whose diff went into a prompt is an author); Invariant 9.
- Fix: remove the diff from step 8, or make H1 an author in step 11.

### M14. An escalation decision that changes a priority cannot be rendered
- Where: :944-949 "carries lines in the bet's form (`REQ-7 Won't`) … written at once, as a rendered task".
- Source: :578 (only a bet writes an inventory version); :722 (a rendering "of the named inventory version"); :595 (MoSCoW must match "the bet copy"). The bet form (:640-642) has no band line.
- Why it matters: the rendered PR fails `layup/spec` and never merges. A band change has no form at all.
- Fix: let a decision write an inventory version and be a copy that `layup spec check` reads. Define the band line.

### M15. Nothing installs `layup-watch` on a target
- Where: :1092-1100 the job "reads … the repository activity" and posts on "the target's control issue". :271 installs only "the LAYUP App".
- Why it matters: with no install and no list entry, a dead host gives no notice. L-F1 names only late or dropped runs.
- Fix: add the install and the list entry to Start, and read both back.

### M16. One task of unknown money freezes a route for good
- Where: :1244-1246 "A route with any task of unknown money gets no step at all". :1226-1227 "the records since that route's last adopted update (a route's records carry over while it gets no change)".
- Why it matters: a token-reporting route has one crashed session with no usage output. That task stays in every later window, so the route never learns. L-H2 names only harnesses with no token report.
- Fix: bound the window, or leave the unknown task out of that route's own mean.

### M17. `clear` for a `pending` kind breaks NFR-004 and REQ-004
- Where: :425 "`pending`, and the head changes no path in the product's scope | `clear`: counts as a pass".
- Source: PRD-0001:135 NFR-004 "report an inactive check as not passed". :117 REQ-004 "a gate that did not run reports `not-active`, which never counts as `pass`". F-0001#5.
- Fix: amend the criteria, with a "Reading of" line for the Operator, or report `not-active` and leave the kind out of that PR's required set.

### M18. Two PRD criteria that step 10 left behind fail by design
- Where: PRD-0001:132 NFR-001 "every decision an issue cites is an ADR, a task record or a decision note in the tree". PRD-0001:179 "`layup psb check` (REQ-001, delivered)".
- Why it matters: the decisions live on the orphan `layup-records` branch, not in the tree (§3). REQ-001 now requires a review session for gaps of meaning, and no code does that.
- Fix: amend NFR-001 to name the records branch. Mark REQ-001 as not delivered.

### M19. Coverage rows fail §14's own rule
- Where: :1305-1307 "a row with none of the three fails". :1335 "the routing and context rules of §9". :1339 "the bets and the Build plan of §8". :1344 "the table of measures in §12". :1343 B3 → "W-10", but W-10 has no acceptance step.
- Why it matters: table C 2.1 (the phased plan) and 3.1, #69 B3, K53–K57, and K11 (the baseline and start values, which has no row) each point to a section, not to a walkthrough, a check or a limit. No walkthrough walks an acceptance, a context refusal, or Delivery Lead Time.
- Fix: add walkthrough steps (the Accept step, the baseline batch), or name checks.

## Notes
- N1. :16-17 still says "O-66 to O-92".
- N2. W-10:21 "a budget escalation is planned input" contradicts :1176 (planned only when it is confirmed as business-forking).
- N3. W-13:23 uses the old "no earlier run read" wording. ADR-0025 d3 leaves out "the middle of the bounds".
- N4. ADR-0020 d6 leaves out `new dependencies` from the tier rule (:804-807).
- N5. No step of the task loop posts `layup/rules`. After the close-out commit, only `layup/verify` is re-posted on the new head.
- N6. :145 "at each start … stops when a rule is missing": a restart during Intake always stops, because the rulesets come only at Scaffold.
- N7. `gofmt` and `go vet` belong to none of the four gate kinds (:408; W-02:19).
- N8. setup/README.md:35 "nothing from LAYUP goes into a target" contradicts §1. It also still describes a TUI.
- N9. REQ-011 says `not reported`; §12 says `unavailable`.
- N10. The version check (:1295-1297) needs the old binary on the host. Where it lives is not stated.
- N11. :751 "the Operator or a parameter passes the diff": there is no parameter row and no answer form for it.
- N12. §3:112 "it opens an issue" contradicts §11 (a notice on the control issue).
- N13. :539-540 "which routing puts on one harness": §9 has no such rule.
- N14. FT5: the target's issues and the `layup-probe` ref are not among the three things. FT6: the harness register and `prices.tsv` stay on the host.
- N15. The checklist statuses of K20, K23, K42, K47, K59, P01, I5, D09 and FT4 point to "§1–4; W-12", which is the wrong section.
- N16. ADR-0015:45-46 and §4:227-230 describe the rule-file scan differently. `tasks/<task>/events.tsv` breaks "one file per record kind".
- N17. §5:283 reads the owner's plan, which an installation token cannot read.
- N18. The setup push from the S04 clone meets the `pre-push` hook on `main`.

## Coverage
- `F-0003#41`–`#52`: each has a walkthrough (W-01 to W-12).
- Phase 1 (REQ-001, 002, 004, 007, 009, 011; NFR-001 to NFR-007): each has a row. The criteria of NFR-001, NFR-004 and REQ-004 fail by design (M17, M18). REQ-001 is marked delivered but is not (M18).
- #69: A1–A3, B1, B2, B4–B9 and the forge question are covered. B3 has no walkthrough step (M19).
- Table C: 8 of 10 rows point to walkthroughs. 2.1 (the phased plan) and 3.1 point to sections (M19).
- §15: L-A1 to L-H2 match their sections. M8, M15 and M16 are gaps that no limit names.

## Existing solutions
- GitHub merge queue (`merge_group`) runs the required checks on the combined commit. It would replace the serial base merge and M12. The audit already expects `merge_queue_merge`.
- SWE-bench's FAIL_TO_PASS / PASS_TO_PASS patch, applied from outside, is the tested way to freeze tests (M11).
- Mutation testing (go-mutesting, Stryker) regenerates known-bad changes on each run, where recorded patches go stale (M6).
- Doorstop, OpenFastTrace and sphinx-needs already do requirement traces with coverage checks. They report gaps rather than block every change (M1, M2).
- A spend-capped proxy (LiteLLM budgets) gives a hard ceiling for API-key harnesses (L-G3, M8).
- Temporal heartbeats, the Kubernetes lease, Shape Up and AgentJury are reused correctly. Healthchecks.io (a push-based monitor) would close L-F1.~~~~

### The author's answer to round 1

All nineteen material findings and notes N1 to N18 are applied in one commit.

| Finding | Fix |
| ------- | --- |
| M1 | Before the first bet, `layup spec check` checks spec-changing tasks against the Intake's facts without the MoSCoW rules, and other tasks are `clear`; at the first bet, code adds one rendered commit (MoSCoW, Phase, non-goals, facts record) to the verified specification's head and carries `layup/verify`. |
| M2 | A requirement marked "waits for a bet" does not fail `layup/spec`. |
| M3 | Before approval, a batch's own gate files run only in a scratch tree, as evidence; after it, they give the status; §6, §8 and ADR-0016 say the same. |
| M4 | Layer 2 has the same two exceptions as `layup/rules`. |
| M5 | A workflows batch: its sessions write it, `layup run` posts its diff, the Operator pushes the branch, then CI and the verifier run. |
| M6 | A batch may replace a recorded patch that no longer applies; the old one is retired with the approval. |
| M7 | `layup setup verify` runs only the checks that apply to a target; `ci` and `procedure` do not; S15 writes the setup record to the records branch. |
| M8 | `--intake-cap` and `--lease-h` at Start: Intake and Shape run as milestone 0 with a cap; pre-answer values have the Start command as evidence. |
| M9 | A question attempt counts toward `stall.attempts` everywhere (§8, §11, W-06, ADR-0019, ADR-0023); a question round is a wait, not a round. |
| M10 | A stalled owner session is a stall of that step (the retry is a new owner session); `answer:` closes the question as a human answer; `reroute` can name a role. |
| M11 | Code applies the frozen test files to each new attempt; a change to them needs an amended plan and a new plan review. |
| M12 | A base merge carries the verification only when it is clean and changes no file the branch touches; otherwise a scoped round. |
| M13 | W-12 step 8 passes findings and the diagnosis, no diff. |
| M14 | A decision's copy writes an inventory version or `budget.tsv`; the band line form is `band B <amount> U <amount>`. |
| M15 | At Start, the Operator installs `layup-watch` and lists the target; `layup run` waits for the job's first notice. |
| M16 | A task of unknown money is left out of the money term, not the whole route. |
| M17 | The reading of `clear` for a pending kind is stated in §6 and in `PRD-0001` NFR-004's criterion, for the Operator to confirm. |
| M18 | NFR-001's criterion names the records branch; §9 says REQ-001's review of meaning is not delivered. |
| M19 | W-11 steps 11 and 12 (acceptance, lead time), W-02 step 8 (the pilot baseline); the coverage rows point to steps or to named checks. |
| N1 to N18 | O-99 in the header; W-10's planned input; W-13 and ADR-0025 wording; ADR-0020 d6; `layup/rules` on each head; rules read only after Scaffold; `gofmt`/`go vet` are a fifth kind; setup README; `unavailable` = `not reported`; the host keeps each recorded version's binary; the diff passes only by a stall answer; §3's dead-man row; one harness for numbering and draft; forge objects and host values stated for FT5 and FT6; checklist statuses corrected; ADR-0015 matches §4, per-task events tables; the plan is the Operator's Start value; the setup push has no hooks. |

## Round 2

2026-09-29, head `552d68d`; 7 min 13 s; 34,026 output tokens; USD 4.90 at list price. Same reviewer and setup. Word for word:

~~~~text
# Whole-design review — round 2
Reviewer: claude-fable-5-1. Head: 552d68d. Verdict: `material findings: 10`.

## Case 1
1. Start: `--new`, root commit, records branch, lease (§5 Start, §2, §3; ADR-0013, ADR-0014). `lease.H` now has a value. The plan value has no source: M10.
2. Intake and probe sessions before the answers (§5 Start step 1, §12; ADR-0024). They run under a cap that §12 never applies: M5.
3. The Intake batch and the answers (§5, §7 steps 1–5; ADR-0018). OK.
4. Scaffold, `layup setup verify`, rulesets (§5, §6; ADR-0016, ADR-0017). The check list does not match the steps' evidence: M4.
5. Shape: the specification (§7 step 6, §8 task loop step 5). Its PR fails the target's `prd-lint` in CI before the bet, so it never reaches the verifier: M1.
6. Shape: the activation batch (§6, §8 step 5; ADR-0016). It fails `layup/rules`, and `layup/gates` for the pending kinds, before its approval: M2.
7. The first bet and the rendered commit (§7 steps 7–8, §8 Bet). The mechanism is OK, but it depends on a verified spec head (M1). A reclass that "waits for a bet" fails `prd-lint`: M1.
8. Merge order, then the milestone plan (§8). OK.
9. A build task, first attempt: plan, test first, draft PR, gates, verifier, close-out, merge (§8; ADR-0019, ADR-0020). OK.
10. A rework attempt with frozen tests (§8 Attempts, §11). This adds an author that §9 does not count: M7.
11. The base merge into the next task (§8 Build). The scoped round fails the target's `review-record` check (RR8): M6.
12. The budget check before each session in milestone 1 (§12; ADR-0024). OK.
13. Acceptance, `acceptance.tsv`, lead time (§8 Accept, §12; W-11 steps 11–12). OK.
14. The audit (§12). OK.
15. `layup learn` and the routing update (§13; ADR-0025). OK. M16 is closed.
16. The lesson and its rule batch (§13, §6). Replacing a stale patch is OK (M6 of round 1 closed). A batch fails step 5 (M2). A workflows batch has two contradictory paths: M3.

## Case 2
1. The developer session ends `needs_context` (§8; ADR-0019 d5). OK.
2. `layup run` records the question, posts it and asks P2 (§8, §10). A question round is a wait, not a round (§11:1046-1048). OK.
3. The owner-role session is routed (§9 owner map; ADR-0020). OK.
4. The owner session hangs. Trigger 3 fires and the session is killed (§11; ADR-0023). OK.
5. The package and the diagnosis (§11 steps 1–2). OK.
6. The retry is a new owner session, and the asker waits (§11:1096-1097). OK.
7. A second hang goes to the panel rung. It is not said which session the panel's path goes to: note N5.
8. The Operator answers `answer:` (it closes the question as a human answer) or `reroute T <role> to H` (§11 step 5). OK.
9. The outcome row and the question row. Early Question Share counts the human question, and Turnaround leaves it out (§12). OK.

## The last-round fixes
- A round 3 `b63e090`: closed.
- O-95 `f32a732`: closed (round-1 N1 fixed).
- B round 2 `ce87f71`: round-1 M19 closed; round-1 M7 still open: M4.
- C round 2 `00a3f4d`: round-1 M4 closed; round-1 M2 closed for `layup/spec`, open for `prd-lint`: M1.
- D round 2 `2b4953a`: round-1 M13 closed; round-1 M1 open (M1), M3 open (M2), M12 open (M6).
- E round 2 `f17bfcd`: round-1 M14 closed (notes N3, N4).
- F round 2 `8e3ba47`: round-1 M9, M10 and M15 closed; the M11 fix adds M7.
- G round 2 `d3a348e`: closed.
- H round 2 `098d565`: round-1 M3 (M2 here), M6 and M16 closed.
- Step 10 `4ad186c`: round-1 M18 closed; round-1 M17 partly open: M8.
- Round-1 fix `552d68d`: round-1 M5 open (M3); round-1 M8 open (M5); new: M9, M10.

## Material findings
### M1. The specification's PR fails the target's `prd-lint` before the bet, and so does a requirement that "waits for a bet"
- Where: architecture.md:590-591 "The PRD's MoSCoW and Phase columns stay empty; the pull request waits." :602-603 "only then does the pull request pass the baseline's own `prd-lint`". :596-597 "its requirement is marked 'waits for a bet' and gets its line at the next bet". :717-719 "A failure goes back to the developer as a finding, in a new attempt." :720 "When they pass, a verifier session reviews the head". :599 "the verified specification's head".
- Source: `prd-lint.sh:115-121` (an empty MoSCoW or Phase fails), run by `ci.yml`. Round-1 M1 and M2.
- Why it matters: the spec PR's CI is red until the bet, so step 5 sends it back each attempt. It hits `stall.attempts` and never gets verified. The rendered commit has no verified head to go on. After bet 1, a reclassed need's new row has an empty MoSCoW, so `prd-lint` fails on that head until bet 2.
- Fix: add a wait state in the task loop for a spec head that waits for a bet (red `prd-lint` only), with verification before the bet. Or keep a "waits for a bet" requirement out of the PRD table until its bet.

### M2. A rule batch fails step 5 before its approval, so it never reaches its verifier
- Where: :717-719 "`layup gate`, `layup spec check` and `layup/rules` report, on each new head … A failure goes back to the developer". :513-515 "`layup/rules` fails a pull request that changes a rule path unless it is an approved batch". :442 "`pending`, and the head changes such a path | failure". W-04:19 "runs its verification before the bet".
- Source: ADR-0019 d2 (the bet approves verified batches); round-1 M3 and N5.
- Why it matters: every batch changes rule paths, so `layup/rules` is red before the approval. The activation batch adds a layout test and contract tests as Go files, so the base manifest's pending kinds fail `layup/gates`. The batch loops through attempts, and the brief can never name a verified head.
- Fix: before approval, exempt a batch's `layup/rules` and pending-kind results from step 5, and use the scratch run as its gate evidence. Say which statuses it must pass.

### M3. A batch that changes workflows has two contradictory paths
- Where: :244-248 "it refuses a branch whose diff from the base touches `.github/workflows/` … the result fails with that reason". :485-489 "its sessions write it, `layup run` posts its diff as a payload, and the Operator pushes the batch branch". :713-714 step 4 "the rule-path and workflow checks (§4, §6); then `layup run` pushes the branch". :250-251 "binds a commit SHA to the session only after the forge accepted the push".
- Source: O-92, O-93; FT3.
- Why it matters: under §4, the batch session's result fails and counts as an attempt. Under §6, it is accepted. Nothing checks that the head the Operator pushed equals the session's SHA before binding. W-13 step 5 walks W-03's CI-job change with no Operator push step.
- Fix: add a batch exception to §4 and to step 4. Make code compare the pushed head with the recorded SHA before binding. Add the push step to W-13.

### M4. `layup setup verify` leaves out checks that the steps cite
- Where: :341-347 "(`pin`, `markers`, `adapted`, `links`, and `identity` and `facts` …) … The evidence of steps S04 to S14 becomes '`layup setup verify <check>` OK'".
- Source: `setup-check.sh:29` (the checks are pin, kit-history, facts, onboarding, glossary, guardrails, markers, adapted, ci, protection, identity, procedure). `steps.tsv` evidence: S05 `kit-history`, S07 `onboarding`, S08 `glossary`, S09 `guardrails`, S12 `ci`.
- Why it matters: five steps get no check or a check that "does not apply", so their evidence cannot be "OK". There is no check called `links`: it is `link-lint.sh`.
- Fix: list every check with "applies" or "dropped, because", and name S12's new evidence.

### M5. Milestone 0 has a cap, but no rule applies it
- Where: :277-278 "the cap of Intake and Shape (milestone 0, §12)". §12:1197 "A milestone's clock starts at its bet". :1216-1217 (rows compare the project total with `B`). :1131-1133 "otherwise the open work goes to the next bet". ADR-0024 d3 "A cap per milestone … fixed by its bet".
- Source: FT2; round-1 M8.
- Why it matters: §12 never mentions milestone 0. Before the answers there is no `B`, so no row of the table can be evaluated. At the intake cap, the breaker sends the work "to the next bet", but no bet can exist before Shape ends. Nothing says who can raise the cap.
- Fix: in §12 and ADR-0024, define milestone 0: when its clock starts, the table rows before `B` exists, and who raises its cap.

### M6. The scoped base-merge round fails RR8 of the target's `review-record` check
- Where: :692-694 "recorded as the next round in the target's form (a round after `nothing material in scope` is `material` only if it finds something)".
- Source: `review-record-lint.sh:276-277` "round N is followed by another round, so its verdict must be `material`"; `review-record.yml` is a required job. Round-1 M12.
- Why it matters: the earlier round's `nothing material in scope` becomes an intermediate round, and the check turns red. The PR cannot merge. The same happens after a conflict's new attempt.
- Fix: record the scoped round in a form that RR8 accepts (for example a separate record kind), or say what the record holds after it.

### M7. Frozen test files carried into a later attempt add an author that §9 does not count
- Where: :788-789 "code applies the frozen test files to each new attempt's start". :836-839 "the authors are … every harness whose session is bound, in the ledger, to a commit in the diff from the base, and every harness whose diff went into an attempt's prompt".
- Source: `F-0003#66`; Invariant 9; round-1 M13.
- Why it matters: H1 writes the tests, and the list freezes. After a stall, `reroute T to H2`. Attempt 3 carries H1's tests in a commit by code, so H1 is not an author and may verify a change whose tests it wrote.
- Fix: count the harness that wrote the frozen tests as an author of every later attempt.

### M8. REQ-004 and REQ-007 still contradict `clear` for a pending kind
- Where: PRD-0001:117 REQ-004 "a gate that did not run reports `not-active`, which never counts as `pass`". :120 REQ-007 "no PR … merges with a selected gate that did not run". architecture.md:440 "(the kind's own gate did not run …)".
- Source: `F-0001#5`; round-1 M17. Only NFR-004 got the reading.
- Fix: add the same reading to REQ-004 and REQ-007, for the Operator to confirm.

### M9. The start values' rendered task cannot pass its handoff
- Where: :1255-1256 "renders the start values into the pilot's PRD as a rendered task (§8)". W-02 step 8. :656-657 a rendered task is "a later bet's or an escalation decision's change". :759 "the change equals code's rendering of the named inventory version".
- Why it matters: start values are no inventory version, so the handoff is never valid and the task never merges. §14's K11 row rests on this step.
- Fix: add start values to the rendered kinds and to the transition row, with the copy they render from.

### M10. Start reads a plan that the command does not give
- Where: :289-291 "records … the plan that the Operator names in the command". :275-276, the command has no plan argument.
- Fix: add the argument (for example `--plan`) or say where the value comes from.

## Notes
- N1. :757 "a known-bad patch per kind it activates" against :453-454 "per kind it touches".
- N2. §10:1021-1022 says `lease.H` and `harness.<id>.cap`/`.wall` are "set at Intake". §5 sets them at Start.
- N3. §7:609 says a band decision writes an inventory version. §10:988 says `budget.tsv`.
- N4. :516 `layup/rules` passes only the rendering of "an inventory version that a bet approved". A decision's facts-record rendering (§7:609-613) fails it.
- N5. For an owner-role stall, the panel's path "goes into the task's next attempt" (:1108-1109). The same gap applies to `external:`.
- N6. §5 step 3 waits for "watch started" with no time limit. A dropped scheduled run (L-F1) holds up Start.
- N7. :692 "Otherwise a verifier session runs one more round" also covers failed gates. Step 5-6 run the verifier only on green checks. Whether this round counts toward trigger 2's cycle cap is not said.
- N8. The fifth kind "static checks" (:423) is not in the catalog list (:399) or ADR-0016. W-02:19 and W-04:17 still name `gofmt` and `go vet` as separate kinds.
- N9. Formatting: stray indentation at :297-298, :369, :1054, :1113, :1406, :1462, :1484, and ADR-0019:60. W-12:25 repeats "with a prompt file". :1096 sits unindented inside item 3.

## Coverage
- `F-0003#41`–`#52`: each has a walkthrough (W-01 to W-12).
- Phase 1: every row is covered. REQ-004 and REQ-007 fail as written (M8). The REQ-001 status is fixed.
- #69: all rows are covered. B3 is now walked in W-11 steps 11–12.
- Table C: 2.1 and 3.1 now name checks and steps. OK.
- K11: W-02 step 8 cannot complete (M9).
- §15: exhausting the milestone-0 cap (M5) and the workflows-batch binding (M3) have no known limit.

## Existing solutions
- GitHub merge queue (`merge_group`): the checks run on the merged commit. That would replace the scoped base-merge round and RR8 (M6).
- SWE-bench FAIL_TO_PASS tests applied from outside, attributed to a separate author (M7).
- Doorstop and sphinx-needs report trace gaps as warnings, not as a red CI. That fits a spec that waits for its bet (M1).
- Healthchecks.io (push-based) would remove the wait for a first scheduled run (N6, L-F1).
- LiteLLM budgets give a hard spend cap for Intake sessions on API keys (M5).~~~~

### The author's answer to round 2

All ten material findings and notes N1 to N8 are applied in the next commit; N9 (formatting) in part. This was the second round, the limit of plan v2 for step 11 (condition 5), so the author asks the Operator on #72 before the approval brief.

| Finding | Fix |
| ------- | --- |
| M1 | A requirement that waits for a bet is written under its own PRD section, not as a table row; a specification that waits for its bet may be red only on `prd-lint`, and is verified anyway (§8 steps 5, 6). |
| M2 | A rule batch before its approval may be red only on `layup/rules` and the pending kinds; the scratch run is its gate evidence (§8 step 5). |
| M3 | A rule batch task's branch with a workflow change is not refused; the Operator pushes it, and code compares the pushed tree with the session's before binding; W-13 step 5 has the push. |
| M4 | `layup setup verify` lists every check that applies to a target, the three that do not and why, and the evidence of S12 and S13. |
| M5 | Milestone 0 (Intake and Shape) in §12 and ADR-0024: its clock starts at Start, its cap is the Start value, and at the cap the Operator raises it or stops. |
| M6 | A task is verified only when it is next to merge and up to date with the base, so no base merge comes between its verdict and its merge; the scoped round is gone. |
| M7 | The harness that wrote the frozen tests is an author of every later attempt. |
| M8 | REQ-004 and REQ-007 carry the reading of `clear`. |
| M9 | Rendered tasks include the start values; the transition row names the copy each renders from. |
| M10 | `--plan` at Start. |
| N1 to N8 | "per kind it touches"; the Start values in §10; a band decision writes `budget.tsv`; `layup/rules` passes a decision's rendering; a panel's path or `external:` goes to the owner's next session; `watch.T` bounds the wait; the scoped round is gone; the static-checks kind is named in the catalog and the walkthroughs. |
