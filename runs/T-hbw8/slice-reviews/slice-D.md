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
