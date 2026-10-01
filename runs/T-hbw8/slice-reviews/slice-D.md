# Slice D review record

Reviewer: Claude Fable 5.1 (`claude-fable-5-1`), a fresh `claude -p` session in a disposable clone, read-only tools, no forge token, `--setting-sources project,local`. Not a gate round.

## Round 1

2026-09-29, head `d7c31bc`; 6 min 36 s; 27,836 output tokens; USD 4.67 at list price. Word for word:

~~~~text
# Slice D review — round 1
Reviewer: claude-fable-5-1. Head: `d7c31bc`. Verdict: `material findings: 10`.

## Material findings
### M1. The verification record goes to the wrong place, and nothing re-runs the check that reads it
- Where: `docs/architecture.md:655-657`: "`layup run` posts its record as a pull-request comment in the target's review-record form". Also W-07 step 8.
- Source: the target's rules in `docs/engineering-discipline.md:416-417`: "Each round is one comment on the issue". The required check `review-record` (`.github/workflows/review-record.yml`) reads the comments of the linked **issue**, and it runs only on `opened, edited, synchronize, reopened`. GitHub docs (events that trigger workflows): an issue comment does not trigger a `pull_request` workflow.
- Why it matters: the draft opens before any review record exists, so `review-record` fails. The record then goes on the pull request, not on the issue, and nothing re-runs the job. §8 step 7 "every required check is green" never happens. K27 stays open.
- Fix: post each round on the task's issue. Name the event that re-runs the job after the record, for example an edit of the pull-request body. Use the target's verdict values (`nothing material in scope`).

### M2. The close-out comes after the verdict and moves the verified head
- Where: `architecture.md:645-646`: the developer writes "the task file and close-out" before verification. `:657`: `layup/verify` is set "at that head SHA". `:661`: the merge is "at that head SHA".
- Source: `engineering-discipline.md:290-296`: the close-out (the verdict, the completed-log line) "lands as a head no round names" and "cannot exist before the rounds finish". AGENTS.md gate step 8.
- Why it matters: if the verdict goes in, the head moves and loses `layup/verify`, so the change cannot merge. If it stays out, the target's gate step 8 has no actor.
- Fix: name an actor and a step for the close-out commit after the verdict. Carry `layup/verify` to that head only when the diff is limited to the bookkeeping files.

### M3. The spec pull request and the rule batch cannot get `layup/verify` or `review-record`
- Where: `architecture.md:476-478` makes `layup/verify` a required check on every pull request, with no bypass. `:622-624`: "Merge order after a bet: the specification pull request … then the rule batch, then the tasks". Only the task loop (`:653-658`) sets `layup/verify`.
- Source: the ruleset of §6. `docs/setup/branch-protection.json` also requires `review-record` and `pr-link`, which need a linked issue with a plan, a plan review and a review record.
- Why it matters: at the first bet, the spec pull request never merges, so the batch and every task wait forever. C8 is not closed.
- Fix: for the non-task pull requests (spec, rule batch, retrospective batch), state the issue, the plan and review records, the verifier and who sets `layup/verify`.

### M4. The bet approves a rule batch whose head moves because of the bet
- Where: ADR-0019:40-42: "The bet approves … its rule batch". `architecture.md:571-574`: the numbered facts record, with the confirmed classes, lands "through that bet's rule batch". `:461-462`: "code refuses the approval when the batch head has moved since the request".
- Source: §6 and §7 (slices B and C), and W-04 steps 4 and 5.
- Why it matters: the idea owner reclasses one fact in the bet comment. The facts file in the batch must change, so the head moves and code refuses the approval. The design gives no second approval. The idea owner's bet and the batch "approver" of §6 may also be different people.
- Fix: build the facts record from the bet copy first, then ask for a separate batch approval by the named approver. Or keep the facts record out of the batch head that is approved.

### M5. The verifier can be a harness that wrote part of the change
- Where: W-12:28: "H1 is admitted again, because the author of attempt 2 is H2". `architecture.md:727-728`: "its harness differs from the author's harness". The author is singular.
- Source: `F-0003#66`: "a verification from a harness agent that did not make the change". Invariant 9.
- Why it matters: attempt 2 starts from the head of attempt 1 (W-12 step 8), so H1's commits are in the diff. H1 then verifies its own work.
- Fix: the authors are every harness whose session is bound to a commit in the diff from the base (from the ledger). Admit only a harness outside that set; otherwise the verification is `not-active`. Fix W-12 step 11 to match.

### M6. The handoff schema is not based on Armature, and the transition table has no content
- Where: `architecture.md:668-669`: "the fields of Spec Kitty's handoff packet v1". `:669-673`: the transition table is given by one example only. W-05:7-10 and :17 use an "architect → developer" pair, but no step of §8 or row of §9 has that pair.
- Source: `F-0003#45` ("based on Armature conventions"). The acceptance criterion of `REQ-005`: "the schema encodes … the record kinds and fields of Armature, so that a schema that accepts everything does not pass". Invariant 4.
- Why it matters: with an empty or unstated table, every handoff passes, and `#59` reads 100% on nothing. K24 is answered by name.
- Fix: give the default rows for each step pair of §8, built from the target's own record kinds (plan, plan review, review record, task file). Fail a pair that has no row.

### M7. "Does not ask the same question again" is a judgement of meaning, tagged `code`
- Where: `architecture.md:686-688`. W-06:23 is tagged "`model`, then `code`".
- Source: the lens: a step that reads meaning is never `code`. `F-0003#71` measures up to "the accepted answer".
- Why it matters: the next attempt rewords the question. Code cannot tell that it is the same question, so it records a false acceptance time.
- Fix: a mechanical rule, for example: status `completed`, the answer ID is cited, and no question of that attempt cites the answer ID. Or make sameness a `model` step and state who does it.

### M8. The tier rule depends on a computation that is not defined and not marked as later
- Where: `architecture.md:734-737`: "the plan session gives each task a size class, and the computed position of §11 (a task with open unknowns is uphill)". ADR-0020:45-46.
- Source: K29 and Sol-23 (no complexity input). Vision 2.2.
- Why it matters: no rule uses the size class. §11 does not exist yet, and "open unknowns" is not defined when a task starts. So every task falls to one tier, or the tier is chosen by chance.
- Fix: state which input sets the tier at the start and who writes it. Mark the hill computation `later: slice F`, or define it here.

### M9. The bet brief has no author
- Where: `architecture.md:614-618`: "`layup run` posts one brief … the rabbit holes, and the no-gos". `:606-607`: back to Shape "only when its bet brief proposes a change of the architecture".
- Source: C8 and Fable-M16 (an actor for each step). The step table (`:700-708`) has no row for this.
- Why it matters: rabbit holes, no-gos and "a change of architecture" are prose judgements that code cannot write. Nobody writes the brief of milestone 2.
- Fix: add a step-table row (role, tier, context) for the bet brief, and say which session proposes a return to Shape.

### M10. Draft pull requests do not exist in a private repository on GitHub Pro
- Where: `architecture.md:649-650` opens a "**draft** pull request". `:280-283` only checks that "the plan does not enforce rulesets".
- Source: docs.github.com, About pull requests: drafts are "in public repositories with GitHub Free … GitHub Pro … and in public and private repositories with GitHub Team and GitHub Enterprise Cloud". Rulesets work on private repositories with Pro.
- Why it matters: a private target on Pro passes the setup check, then fails at the first task.
- Fix: check that drafts are available at setup (§5 step 3), or state what the loop does without drafts.

## Notes
- N1. `branch-protection.json` has `"strict": true`. With tasks in parallel (`:628-629`), each later pull request must update its branch after a merge. That gives a new head, which needs `layup/verify` again. State who updates the branch and whether the verification repeats. A merge queue is one option.
- N2. Marking a pull request ready is a GraphQL mutation (`markPullRequestReadyForReview`). The REST update has no draft field. Add it to the forge calls of §1.
- N3. An attempt that a question ends "does not count toward the attempt limits" (`:688-689`). Until slice F bounds it, a question loop has no limit. Slice F should say so.
- N4. The model's context size and "bytes divided by four" decide whether a start is refused, but they live on the host. Record the values used in the session start row (FT6), and give the context size its evidence (Invariant 4).
- N5. The harness probe runs only at Intake. A harness upgrade during delivery is not probed again.
- N6. A "delivered requirement" (Accept, `:630`) is not defined. For example: every task of it in the task register has merged.
- N7. K28 reads "reaches human review" as a review request or a ready state, while a draft is still visible to anyone. Put this reading in the approval brief, like the readings of §6 and §7.
- N8. "Context too large, the plan session splits the task" does not fit steps that are not tasks: plan review, verification, the specification.
- N9. `:670` gives the transition table to O-81, but O-81 names the role matrix and the owner map, not the transition table.
- N10. W-07 step 2 is tagged `model` only, but it includes the posting by `layup run`, which is `code`.

## Checklist rows
- S5: not answered in part (M6).
- S6: answered, except the acceptance rule (M7).
- S7: answered in part (M1, M2, M5).
- R04: answered (M10 caveat).
- I5 (in part): answered.
- I9 (in part): not answered for a change made on two harnesses (M5).
- K24: not answered in part (M6).
- K25: answered, except M7.
- K26: answered.
- K27: not answered in part (M1, M2, M3).
- K28: answered (N7).
- K29: not answered in part (M8, M9).
- K30: answered; L-D2 is a known limit (acceptable).
- K31: answered.
- P03 (in part): answered for the role path; the smart-if is later (slice E).
- P08: answered, except M5.
- P15: answered, except M8.
- P16: answered, except M5.
- P18 (in part): answered for model and context routing.
- C4: answered.
- C8: answered in part (M3, M9).
- D01 (in part): answered.
- D12: answered.
- FT1: answered.
- FT3: answered.
- FT4: answered.
- L-D1: known limit (acceptable).

## Existing solutions
- GitHub Copilot coding agent: it opens a draft pull request, works, then marks it ready. This is the same pattern; compare how it handles checks that must re-run after comments (M1).
- Shape Up: in Basecamp the hill position is set by the team by hand, not computed. LAYUP's "computed position" is new and needs its own definition (M8).
- GitHub merge queue: it handles strict up-to-date merges of parallel pull requests (N1).
- Spec Kitty's handoff packet and AgentPlane's transition table are reused. Armature's own review-record form is the convention that `F-0003#45` names (M6).

Sources: [About pull requests](https://docs.github.com/articles/about-pull-requests), [Events that trigger workflows](https://docs.github.com/en/actions/writing-workflows/choosing-when-your-workflow-runs/events-that-trigger-workflows), [REST pulls](https://docs.github.com/en/rest/pulls/pulls).~~~~

### The author's answer to round 1

All ten material findings and notes N1 to N10 are applied in one commit (§8 and §9 rewritten).

| Finding | Fix |
| ------- | --- |
| M1 | The review record goes on the task's issue with the target's verdict values; `layup run` then edits the pull request's body, which re-runs `review-record`. |
| M2 | A close-out session commits the close-out after the verdict; code limits that commit to the task file and the completed log and sets `layup/verify` on the new head. |
| M3 | Every pull request is a task (specification, rule batches, facts record, build tasks), so each gets an issue, plan, review records and `layup/verify`. |
| M4 | The bet approves verified batches by head SHA and does not change them; what the bet changes, code renders as its own tasks, checked against the new inventory version. |
| M5 | The authors are every harness bound to a commit of the change; a new attempt starts from the base; W-12 steps 8 and 11 follow. |
| M6 | The transition table has default rows for each step pair, in the target's own record kinds; a pair with no row is not valid; W-05's case is now a real pair. |
| M7 | Acceptance is mechanical: `completed`, cites the answer ID, and no question cites it. |
| M8 | The tier starts from the plan's size class, open question or new interface; the progress position is later: slice F. |
| M9 | A Product Owner session writes the bet brief; code checks the headings, a line limit and links; the step table has its row. |
| M10 | The Start checks that the plan offers draft pull requests. |
| N1 to N10 | Tasks merge one at a time and re-verify after the base moves; the GraphQL ready call is a forge call; §11 bounds a task's questions; context size with its source, and both numbers in the start row; a new harness version is probed again; "delivered" defined; the reading of `F-0003#58` goes to the brief; a non-task step that does not fit fails to §11; the transition table is a parameter; W-07 step 2 tagged `model`, `code`. |

## Round 2

2026-09-29, head `34d7754`; 5 min 7 s; 23,725 output tokens; USD 2.92 at list price. Same reviewer and setup. Word for word:

~~~~text
# Slice D review — round 2
Reviewer: claude-fable-5-1. Head: `34d7754`. Verdict: `material findings: 6`.

## Material findings
### M1. A second verification after a base merge breaks the target's `review-record` check
- Where: `architecture.md:645-647`: "after a merge, `layup run` merges the base into the next task's branch, and the gates and the verification run again on its new head". Step 6 (`:671-674`) posts every verification as a review record on the issue.
- Source: `docs/ci/review-record-lint.sh` RR8: "round N is followed by another round, so its verdict must be `material`". RR7: "the k-th record carries cycle k-1", and no cycle may exceed the cap. `engineering-discipline.md:309`: "A merge of `main` is the one other thing that may land", so the reviewed SHA survives. Close-out "cannot exist before the rounds finish" (`:290-296`).
- Why it matters: task T-2 already has round 1 with `nothing material in scope` and its close-out commit. T-1 merges, and T-2 gets a second record. RR8 fails on round 1, the required check goes red, and T-2 never merges. If the second round is material, the verdict in the close-out is wrong.
- Fix: follow the target's rule. After a clean merge of the base, carry `layup/verify` over to the merge commit, with no new round. Or state the form of the second record, its `Cycle`, and what happens to the close-out.

### M2. The only row into the verifier fits build tasks, so the other tasks cannot pass the handoff check
- Where: `:612-614`: "The specification, each rule batch, the facts record and each build task are tasks, and each goes through the task loop". `:703`: "developer | verifier | … the plan's tests fail at the base and pass at the head". `:711`: "a pair with no row is not valid".
- Source: invariant 4 and `F-0003#45`. Round 1 M3 and M6.
- Why it matters: a specification task, a gate-activation batch, a milestone-plan task and a facts record have no test that fails at the base. None of them has a row of its own, so step 4 refuses each one, and the first bet never gets its batches or its specification. For the facts record and the MoSCoW task that code renders (`:635-637`), steps 2 and 3 have no session to plan or write, so the loop has no actor there.
- Fix: add default rows for each kind of task (specification, batch, milestone plan, rendered record) into the verifier. For a rendered task, state which steps run and who does the plan and plan review, or exempt it by name.

### M3. The code does not enforce a plan review on another harness
- Where: `:659-660`: "a plan-review session on another harness reviews it". `:771-775`: code admits a plan reviewer when "its harness wrote no commit of the change: the **authors** are every harness whose session is bound … to a commit in the diff from the base".
- Source: the step table (`:750`), "QA Engineer, on another harness". ADR-0020 decision 4. The target's reviewer independence rule.
- Why it matters: a plan is an issue comment, not a commit. At plan review, the author set is empty, so code admits the harness that wrote the plan, and the plan reviews itself.
- Fix: for a plan review, the authors are the harnesses of the plan sessions of the task, from the ledger.

### M4. The diff of attempt 1, given as a payload, lets H1 verify its own code
- Where: `:729-731`: "A new attempt starts from the base commit … with the last attempt's diff and findings in its prompt file as payloads, so that its commits have one author harness." W-12 steps 8 and 11: "its only author is H2, and H1 is admitted".
- Source: `F-0003#66`: "a verification from a harness agent that did not make the change". Invariant 9. Round 1 M5.
- Why it matters: H2 applies H1's diff as it is and commits. The ledger names only H2, so H1 verifies code that H1 wrote. This is the same case as round-1 M5.
- Fix: count every harness whose diff went into the attempt's prompt as an author. Or do not pass a diff from another harness.

### M5. The special pass for the facts record contradicts §6 and ADR-0017, and the merge order breaks the batch hash
- Where: `:574-576` and `:636-638`: the facts record's "`layup/rules` passes when the file equals the rendering". `:638-639`: "the specification, the facts record, the rule batches, then the build tasks". `:647`: "the branch rule is 'up to date'".
- Source: `:483-484`: "`layup/rules` fails a pull request that changes a rule path unless it is a batch whose head's rule-file hash equals the hash recorded with its approval". ADR-0017:52 says the same. `docs/facts/` is a rule path (`:454`). The reading of `F-0003#64` (`:487-489`) counts a rule-path change outside an approved batch.
- Why it matters: §6 and ADR-0017 fail the facts record, and the reading of `#64` counts it as an agent write. The facts record then merges first. The rule batch must merge the base in, and if "the tree hash of its rule files" covers `docs/facts/`, the hash no longer matches the approval. The batch cannot merge, and "no other approval inside the milestone" gives it no second one.
- Fix: amend §6, ADR-0017 and the reading of `#64` for a record that code renders. Define the batch hash over the batch's own changed files, or merge the batches before the facts record.

### M6. W-11 step 8 still describes the old design
- Where: W-11:25: "writes the PRD's MoSCoW and Phase columns from it … the numbered facts record goes into the bet's rule batch | a commit on the pull request".
- Source: `architecture.md:574-576`, ADR-0018 decision 5 ("as its own task"), `:635-637` ("does not touch the approved batches … as their own tasks").
- Why it matters: this is the round-1 M4 defect. A reader who follows the walkthrough moves the head of the approved batch and of the verified specification.
- Fix: rewrite W-11 step 8 to match the separate rendered tasks.

## Notes
- N1. `:751` "QA Engineer, on another harness" and W-11 step 10 "on another harness": another than which harness? State "outside the authors" in both.
- N2. The target requires `conventional-title` (`branch-protection.json`). Step 4 does not give the pull request's title form.
- N3. §6 (`:465-466`, O-93): a batch that changes `.github/workflows/` is pushed and merged by its approver. Steps 4 and 8 name `layup run` for every task. Say so in the loop.
- N4. ADR-0019 decision 3 says "the merge … at the verified head SHA". After close-out, the merged head is not the verified one. Match `:681-684`.
- N5. ADR-0019 decision 6 says "The idea owner accepts". §8 (`:651`) says "the approver". Use one term.
- N6. The tier rule (`:781-783`) and "a proposed change of it, which sends the milestone back to Shape" (`:625`) need fixed fields in the plan form and the brief form, or code cannot read them.
- N7. Acceptance (`:723-725`) is a mechanical stand-in: an attempt that rewords the question without citing the answer ID still counts as acceptance. Record this as a limit of the `#71` measure.
- N8. W-04 steps 3 and 4 do not show the batch as a task (issue, plan, verification, close-out) or as approved inside the bet comment. §8 now requires both.

## Checklist rows
- S5: not answered in part (M2).
- S6: answered.
- S7: not answered in part (M1).
- R04: answered.
- I5 (in part): answered.
- I9 (in part): not answered in part (M4).
- K24: not answered in part (M2).
- K25: answered (N7).
- K26: answered.
- K27: not answered in part (M1, M2).
- K28: answered; the reading of `#58` goes to the Operator.
- K29: answered (N6).
- K30: answered; L-D2 is a known limit (acceptable).
- K31: answered.
- P03 (in part): answered for the role path; the smart-if is later (slice E).
- P08: not answered in part (M3, M4).
- P15: answered; the progress position is later (slice F).
- P16: not answered in part (M4).
- P18 (in part): answered.
- C4: answered.
- C8: not answered in part (M2, M5).
- D01 (in part): answered.
- D12: answered.
- FT1: answered.
- FT3: answered: the authors come from the ledger (M4 is a gap in what the ledger counts, not text as the actor).
- FT4: answered.
- L-D1, L-D2: known limits (acceptable).

Round-1 findings: M1, M7, M8, M9 and M10 are closed. M2 is closed, apart from new finding M1. M3 and M6 are closed in name, but M2 above is still open. M4 is closed for the batch head, but new findings M5 and M6 remain. M5 is not closed (new finding M4).

## Existing solutions
- GitHub merge queue tests the merged result once, and does not re-verify each branch after every merge. It answers M1 without extra rounds.
- Armature's own rule, "merge `origin/main` in, so the reviewed SHA survives", is already in the target and the design ignores it (M1).
- The GitHub Copilot coding agent and the AI-SDLC flow use the same draft-until-verified pattern. Spec Kitty's handoff packet and AgentPlane's transition table are reused properly for build tasks.

GitHub facts checked: drafts on private repositories need Team or Enterprise Cloud (docs.github.com, About pull requests). An edit made with an App installation token triggers `pull_request: edited`; an edit made with `GITHUB_TOKEN` would not. Marking a pull request ready is the GraphQL `markPullRequestReadyForReview`.~~~~

### The author's answer to round 2

The six new material findings and notes N1 to N8 are applied in the next commit. Under O-98, these fixes are a named item of the whole-design review.

| Finding | Fix |
| ------- | --- |
| M1 | A clean merge of the base whose gates pass carries `layup/verify` to the merge commit with no new round (the target's own rule); a conflict or a failed gate starts a new attempt. |
| M2 | Rows into the verifier for each kind of task (specification, rule batch, milestone plan, rendered task); a rendered task has a plan and change by code, and a plan review and a verification by sessions. |
| M3 | For a plan review, the authors are the harnesses of the task's plan sessions. |
| M4 | A new attempt gets the findings and the diagnosis, not the diff; if a diff is passed, its harness is an author. |
| M5 | `layup/rules`, ADR-0017 and the reading of `F-0003#64` accept a rendered record of an approved bet; a batch's hash covers only the rule files it changes. |
| M6 | W-11 step 8 renders the columns and the facts record as their own tasks. |
| N1 to N8 | "outside the authors" in the step table and W-11; the pull-request title form; the workflows batch in step 8; ADR-0019 d3 matches the close-out; "the approver" in ADR-0019; fixed fields in the plan and the brief; L-D3; W-04 shows the batch as a task inside the bet. |
