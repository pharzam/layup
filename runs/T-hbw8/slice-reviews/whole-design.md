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
