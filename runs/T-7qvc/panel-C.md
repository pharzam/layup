## Panel member C — governance, the invariants and the human decision points

Scope: options only (ADR-0006). I do not select. "Test" = a check that can show the argument is wrong. Invariants are `F-0001#1`–`#9`. Human Decision Points are DP1–DP5 (`F-0001#10`–`#14`).

**A condition that applies to all questions.** O-9 records that the agents and the Operator push with the same GitHub account. With one account, no forge control can tell an agent write from an Operator write. So no option below can measure "0 agent writes to rule paths" (`F-0003#64`) or keep an agent from posting its own `pass`. Q1, Q2 and Q4 all depend on this. It is Operator question 1.

---

### Q1. The runner for stack gates

**1-A. A GitHub Actions workflow in LAYUP's own repository.** A `workflow_dispatch` or a schedule starts the job. The job checks out the target PR head (read token), runs `layup gate`, and posts one commit status per gate to the target head SHA. The target's branch protection requires those status contexts.
- For: the gate code stays in LAYUP, so an agent in the target cannot change it (Inv 3, O-11). Test: a target PR that edits any file cannot change the gate's verdict logic. If a diff in the target changes the verdict other than through the checked content, this is false.
- For: GitHub shows a required context that never reports as "expected". That blocks the merge, so `not-active` never passes (Inv 5, NFR-004). Test: stop the runner, open a PR and try to merge. If the merge succeeds, this is false.
- Against: anyone with write access can post a commit status, and the agent has that access under the same account. Test: with the agents' token, post `success` for the context without a gate run. If the forge accepts it, Inv 3 fails.
- Against: a required context that is named for LAYUP stays in the target's protection, so without LAYUP no PR can merge (Inv 2). Test: remove LAYUP and try to merge. If the merge is blocked, Inv 2 is strained until someone removes the context.
- Serves Inv 5, 6, 7 and 8. Strains Inv 2 and Inv 3, and Inv 1 (the verdict is in the forge, not in Git, unless it is also written as a row).
- Main risk: a status that the agent can forge.

**1-B. A LAYUP GitHub App (a check-run producer).** An App identity creates the check runs, and the target's required check names that App as its expected source. The trigger is a webhook or a poll.
- For: GitHub can require that a check comes from a specified App. The agents' account cannot then satisfy it (Inv 3). Test: post a status of the same name with the agents' token. If the protection accepts it, this is false.
- Against: a webhook receiver is a network service, and ADR-0011 decision 1 excludes a network service for the engine. This needs a new ADR, or a poll loop with no service.
- Against: at review time it depends on an external service, the cost that ADR-0012 Context 3 names. Test: count Operator decisions that runner outages cause in the pilot. More than zero is the cost.
- Serves Inv 3, 5 and 6. Strains Inv 2 (as in 1-A) and ADR-0011 decision 1.
- Main risk: one more identity and service for someone to own. Who owns the App key?

**1-C. Invoked by a harness or the Operator, with the verdict in Git first.** A harness agent that is not the author, or the Operator, runs `layup gate` from a LAYUP checkout against the PR head. It appends the verdict to `runs/<task>/gates.tsv` in LAYUP (by SHA and per gate), then posts a status.
- For: Inv 1 holds, because the verdict is a Git row before it is a forge status. Test: delete the forge status. If the verdict cannot be recovered from Git, this is false.
- Against: nothing stops the author's session from invoking the gate itself. Test: the author runs the gate and posts the status. If the merge goes through, Inv 3 and REQ-007 fail.
- Serves Inv 1, 5, 6 and 9 (any harness can run the CLI). Strains Inv 3.
- Main risk: the gate becomes procedural, not structural.

| | Forge blocks when the gate did not run (Inv 5) | Agent cannot forge a pass (Inv 3) | Verdict in Git (Inv 1) | Target merges without LAYUP (Inv 2) | Needs a new ADR |
|---|---|---|---|---|---|
| 1-A | yes | no, with one account | only if a row is added | no | no |
| 1-B | yes | yes | only if a row is added | no | yes (dec. 1) |
| 1-C | yes | no | yes | no | no |

---

### Q2. The rule-protection control (REQ-003, Inv 3)

**2-A. Code owners and a second identity.** Put `CODEOWNERS` on the rule paths (ADR-0011's list) and set branch protection to "require code-owner review". The agents use a machine user or App that is not a code owner, and the Operator's own account approves.
- For: this is the control that ADR-0011 already names. Test (`F-0003#64`): an agent PR that touches `.github/**` cannot merge. If it merges, this is false.
- Against: `CODEOWNERS` is a file in the target. If the kit copy does not have one, adding it strains O-11 ("nothing but the kit and the facts").
- Against: `F-0001#28` lists "a change to a gate" as unplanned input. Test: count rule-path approvals in the pilot. Each one raises the Task Intervention Rate unless it is a planned approval point (DP3).
- Serves Inv 3 and 7. Strains O-11 and the DP3/`#28` accounting.
- Main risk: the Operator becomes the approval queue for every rule change (Problem 1).

**2-B. A push-time refusal.** The agents' identity has a forge ruleset (a push rule that restricts file paths), or a fork-only token with no write to rule paths. The protection body is exported as JSON, in the same way as `docs/setup/branch-protection.json`.
- For: the forge refuses the write before a PR exists, so the audit "0 agent writes" is true by construction. Test: push a commit to a rule path with the agents' token. If the push is accepted, this is false.
- Against: the forge plan decides whether this is available. Test: check the plan's features for the pilot targets. If a pilot cannot enable it, the option fails there.
- Against: a forge setting is outside Git (Inv 1), unless an exported copy is checked for drift.
- Serves Inv 3 and 5. Strains Inv 1 and generality (`F-0003#67`) across forges.
- Main risk: the agents can never propose a rule change, even one that is correct.

**2-C. Detection only: the complement that ADR-0011 names today.** A deterministic list of each rule-path change, made after the merge, plus a `layup gate` rule-path gate. The gate fails a PR that touches a rule path unless the PR links an approved decision record. CI already runs each check from the default branch (the `ci.yml` restore step), so a PR cannot weaken its own check.
- Against: the agent can write the "approved" record under the same account. Test: seed a forged approval record. If the gate passes, REQ-003 ("a control the agents cannot pass by themselves") is not met.
- Against: a check restored from `main` protects only the PR's own run. Test: merge a PR that weakens a check. The next PR runs the weaker check.
- Serves the audit half of `F-0003#64` and Inv 6. Does not meet Inv 3 by itself.
- Main risk: the ADR calls it a control, but it only detects.

| | Refuses an agent write (Inv 3) | Measures "0 agent writes" | Writes into the target | Human load |
|---|---|---|---|---|
| 2-A | yes, with two identities | yes, with two identities | `CODEOWNERS` (O-11?) | one approval per rule change |
| 2-B | yes | yes | nothing in the tree | none; rule changes need a human push |
| 2-C | no | only with two identities | nothing | an audit after the merge |

---

### Q3. The role model and the handoff schema (REQ-005)

**3-A. The roles that the kit's gate already has.** The roles: author (execution tier), plan reviewer, reviewer (a different model), examiner, Operator and idea owner. The record is `runs/<task>/handoff-<n>.md` in the kit's review-record shape. A linter in the manner of `review-record-lint` validates it.
- For: "based on Armature conventions" (REQ-005) is true by construction. Test: the PRD criterion requires that a record which accepts everything does not pass. Feed the linter an empty or garbage record. If it passes, this is false.
- Against: the PSB §2 functions (for example "QA Engineer", `F-0004#3`) have no role here. Test: map each pilot step to a role. A step with no role is a gap.
- Serves Inv 6 and 9. Strains nothing directly.
- Main risk: the roles describe review, not delivery.

**3-B. PSB functional roles and a TSV event table.** The roles: specifier (REQ-012), architect, implementer, tester, reviewer and examiner. The record is `runs/<task>/handoffs.tsv` with the columns `seq, from_role, to_role, harness, model, artifact_path, artifact_sha, req_ids, open_questions, status`. The engine's header and column check validates it (ADR-0011 decision 2), with required fields per `to_role`.
- For: Inv 9. Test (`F-0003#66`): give only the TSV and the files it names to a different harness, with no chat history. If that harness cannot continue, this is false.
- Against: the role list is an architecture choice, not a PSB fact (`F-0001#21`). Six roles is a value that has no evidence yet (Inv 4) until the pilot.
- Serves Inv 1, 4 (header check), 6 and 9. Strains Inv 4 (the size of the list).
- Main risk: the roles and the handoffs are more numerous than the pilot needs.

**3-C. The schema is a rule and the role names are content.** The fields are fixed and protected. Each target lists its roles in a register, `docs/roles.tsv`.
- For: Inv 7 (the domain changes content, never rules). Test: a target that adds a role still fails a record with a missing field.
- Against: without Q2, an agent can edit the register or the schema. Test: an agent PR that edits the schema merges.
- Serves Inv 7. Strains Inv 3 until Q2 is decided.
- Main risk: the roles drift between targets and the pilot cannot compare them.

| | Traces to the kit | Harness-neutral | Role list has evidence | Depends on Q2 |
|---|---|---|---|---|
| 3-A | high | yes | yes (the kit) | low |
| 3-B | medium | yes | no (a hypothesis) | low |
| 3-C | medium | yes | per target | high |

---

### Q4. The escalation rule (REQ-008, DP4)

**4-A. Declared fields.** Each decision record has four required yes/no fields: `budget`, `legal_or_compliance`, `approved_intent`, `trade_off_between_goals`. If any field is `yes`, the engine stops.
- Against: the agent classifies its own decision. Test: seed a business-forking decision with all four fields set to `no`. If the rule misses it, DP4 ("an agent never makes…") fails.
- For: it is cheap and deterministic about the fields (Inv 6).
- Main risk: false negatives that nobody sees.

**4-B. Triggers on intent artifacts.** The rule reads the diff and the telemetry. It stops when a change touches an intent path (`F-0001`, `F-0004`, the PRD success criteria and priorities, a budget file), removes or defers a `Must` requirement, moves a date, changes a dependency licence, or makes telemetry cost pass a budget value.
- For: Inv 6. Test: seed "drop REQ-x to keep the date". If the rule does not fire, this is false.
- Against: a legal position that no file encodes is not selected. Test: seed "store personal data in region Y". It will not fire unless a compliance file exists.
- Main risk: the trigger list is itself intent, so the idea owner must own it.

**4-C. Hybrid, where the stricter result wins.** Use 4-B's triggers together with 4-A's fields. An LLM classifier may only add escalations, never remove them.
- For: fewer false negatives.
- Against: every false positive becomes unplanned input (`F-0001#28`: an answer that the idea owner does not confirm) and raises the Task Intervention Rate. Test: the confirmation ratio in the pilot sample.
- Main risk: the error direction is a business choice.

**Record, the same for all three.** `runs/<task>/escalations.tsv` with the columns `id, trigger, decision_path, stopped_at_sha, answer_path, confirmed_business_forking, answered_at`. The answer is written to Git before the task continues (Inv 1), in the same way as ADR-0011 decision 6. The stop is structural only if an open row makes the gate report `fail` (this depends on Q1).

| | False negatives | False positives (TIR) | Inv 6 | Who owns the triggers |
|---|---|---|---|---|
| 4-A | high | low | partial | the agent |
| 4-B | medium | low to medium | yes | the idea owner |
| 4-C | low | high | partial | the idea owner |

---

### Q5. The stall procedure (REQ-009, REQ-010, DP5)

**The state (the same for all).** An open row in `runs/<task>/stalls.tsv` (`id, kind, limit, limit_value, opened_sha, examiner_harness, examiner_model, session_id, diagnosis_path, outcome, closed_at`). An issue label is outside Git and strains Inv 1.

**5-A. A count limit on an unchanged fingerprint.** "Repeats without progress" means the same gate, the same failing rule and the same input hash N times. "No agreement" means the same finding is contested for more than K rounds.
- For: deterministic and harness-neutral. Test: a seeded loop stops at exactly N, and slow real progress (a different failure each time) is not flagged.
- Against: a harness can change a byte and reset the fingerprint. Test: a seeded loop that only changes whitespace. If it is not flagged, this is false.

**5-B. A wall-clock limit.** There is no change in the set of failing checks for T minutes.
- Against Inv 9: a slower harness stalls sooner. Test: run the same seeded task under two harnesses. If the stall points differ only by speed, the limit measures the harness, not progress.

**5-C. A cost limit from telemetry (REQ-011).** Tokens per step are more than k × the pilot median.
- For: it uses a record that already exists.
- Against: there is no median until the baseline exists (`F-0004#11`). Before that, k is a guess (Inv 4).

**The examiner.** (i) A fresh session on a different harness from the one that did the stalled work, tried in order as ADR-0012 part 3 does: the strongest for Inv 9 and `F-0003#66`. (ii) The same harness, a fresh session and a different model: cheaper, and weaker independence. (iii) The Operator: this is human input, so it works against `F-0003#73` (≥ 90 % of stalls close without human input).

**The records.** `runs/<task>/stall-<id>.md` holds the diagnosis: the kind, the limit that was hit, evidence as SHAs and log paths, the examiner's identity, the recommended action. The package for the Operator is that file plus the Q3 handoff record, so that "give the task to a different harness" (DP5) is possible. The Operator's answer is written to Git.

**All three options:** every value for N, K, T or k chosen now is a value without evidence (Inv 4) unless it is marked as a working hypothesis that the pilot baseline replaces.

| | Inv 6 | Inv 9 | Inv 4 before the pilot | Can be gamed |
|---|---|---|---|---|
| 5-A | yes | yes | hypothesis | fingerprint reset |
| 5-B | yes | no | hypothesis | idling |
| 5-C | yes | yes | no value exists | splitting steps |

---

### Questions for the Operator

1. **Identity split (O-9).** Do the agents get their own GitHub identity (a machine user or an App)? Who holds its key? Q1-B, Q2-A/B and the `F-0003#64` audit all need this.
2. **O-10/O-11 scope.** Does a forge setting on a target (branch protection, required checks, a ruleset) or a `CODEOWNERS` file count as "writing LAYUP into the target"?
3. **Network service.** Does LAYUP accept a GitHub App or webhook runner (Q1-B), which ADR-0011 decision 1 excludes for the engine? Or must the runner stay a job in LAYUP's CI?
4. **Where the verdict lives.** Is the gate verdict written in Git in LAYUP's `runs/`, in the target's `runs/`, or in both (Inv 1 against Inv 2)?
5. **Rule-change approval.** Is it a planned approval point (DP3, recorded before delivery) or unplanned input (`F-0001#28`)?
6. **The escalation triggers.** Must the idea owner approve the trigger list, because it defines "approved intent"? Which error direction is accepted: more false positives (a cost to TIR) or a risk of false negatives (a cost to DP4)?
7. **The stall limits.** Are N, K, T or k set now as marked hypotheses, or left as open gaps until the pilot baseline?
8. **The examiner.** Must the stall examiner be on a different harness, or is a fresh session on the same harness enough?
