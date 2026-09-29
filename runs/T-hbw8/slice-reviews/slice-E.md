# Slice E review record

Reviewer: Claude Fable 5.1 (`claude-fable-5-1`), a fresh `claude -p` session in a disposable clone, read-only tools, no forge token, `--setting-sources project,local`. Not a gate round.

## Round 1

2026-09-29, head `6395888`; 5 min 22 s; 24,201 output tokens; USD 3.36 at list price. Word for word:

~~~~text
# Slice E review — round 1
Reviewer: claude-fable-5-1. Head: 6395888. Verdict: `material findings: 13`.

## Material findings
### M1. A threshold on a "probability of yes" has no defined meaning
- Where: architecture.md:861 "an answer below the threshold at a `delegate` point goes to a human"; :853 "an answer at or above the point's threshold decides".
- Source: O-71; jev-sources.md `confidence`: "Noul answers carry no separate confidence"; "A confidence threshold is not one number".
- Why it matters: P1 and P2 return p(yes). Take P1 at `delegate` with p(yes)=0.02. Read one way, that is a confident "no". Read the other way, it is "below the threshold", so every "no" goes to the idea owner. The text also does not say whether a confident P1 "no" at `delegate` can drop the session's own `decision_needed`. The OR rule at :888 says it cannot; "decides" at :853 says it can.
- Fix: define a band per point (yes at p≥t, no at p≤1−t, a human in between) and a threshold for a Choice or Score answer. State that at every level P1 only adds to the floor and to the declarations.

### M2. It is unclear which levels send a provider failure to a human
- Where: :861 "an error, a time-out, a `429` after its backoff, or an answer below the threshold at a `delegate` point goes to a human".
- Source: K34 (Sol-7); `F-0001#24`.
- Why it matters: if "at a `delegate` point" covers only the last item, an error on a `shadow` P4 call sends a routine routing choice to the Operator. That is the K34 defect again.
- Fix: say that under `off`, `shadow` and `cautious` a failure takes the deterministic branch and is recorded, and that only `delegate` sends it to a human.

### M3. Under the `shadow` level that Intake offers, the prose screen is inactive, and nothing records that
- Where: :889 "At `cautious`, the provider can only add a candidate"; :1017 L-E1 "in prose that the four questions miss"; W-08 step 4 "under `shadow`, the floor alone selects it".
- Source: Invariant 5 (`F-0001#5`), FT1; K32, K35 (Fable-M7, M8); `F-0003#57`.
- Why it matters: under `shadow` or `off`, P1 never adds a candidate and P2 always trusts the asker's label. So a scope question labelled "domain" is answered by an agent. Only that agent's own declaration stops it, which is the PSB's root cause. The screen still records "not selected" as if the prose had been screened. L-E1 describes the gap as a miss by the questions. It is really every undeclared prose decision.
- Fix: record in the selection row that P1 or P2 was not active. Extend L-E1 to the `shadow` and `off` case. Mark K32 and K35 as answered only at `cautious` or `delegate`.

### M4. P2's branch "needs a human = yes" goes nowhere
- Where: :845 P2 "A safer branch: 'needs a human' = yes"; §8:736 names only "A question that needs no human".
- Source: K35 (Fable-M8 fix: "Add a branch to the idea owner"); `F-0001#13`.
- Why it matters: at `cautious`, a question that P2 sends to a human has no recipient, no record and no wait state. On an error, :862 sends it to the Operator, even when the question is business-forking.
- Fix: route a "yes" through the escalation screen to the idea owner, and state where every other case goes.

### M5. The calibration record measures agreement with the default, not correctness
- Where: :857 "the share of `shadow` rows where the provider's answer agreed with the branch that a check or a human took".
- Source: K38 (Fable-N7: "promotion has no ground truth"); Invariant 4.
- Why it matters: under `shadow`, P2's branch is the asker's label (which the PSB says is unreliable), and P4's is the table order. At P1, the provider scores low exactly when it finds what the floor missed. The evidence then points the wrong way.
- Fix: name a ground truth per point: the idea owner's confirmations and the Missed Escalations sample for P1, the owner session's disposition for P2, the verification outcome for P4. Otherwise record K38 as a known limit.

### M6. The floor selects every specification task and every rendered task
- Where: :883 the floor reads a diff that "touches ... the PRD's requirement rows, their priorities or criteria".
- Source: §7 step 6 (a session writes the PRD); §8:643 (code renders MoSCoW from the bet copy); `F-0001#28`.
- Why it matters: every spec task and every rendered MoSCoW task escalates to the idea owner. Each "no" then counts as unplanned input.
- Fix: exempt a diff that equals code's rendering of a bet copy, and the Shape-phase specification that the bet approves. Keep the floor for build tasks.

### M7. W-08 step 5: code cannot write the options
- Where: W-08:22 `code` "the options (the paid service; plain SMTP)".
- Source: lens 1; :827 "It writes no text"; the candidate came from the floor and P1, not from a `decision_needed` with options.
- Why it matters: no input to step 5 contains "plain SMTP". Only a session can propose an alternative.
- Fix: add a model step (a session writes the options), or make the options a fixed plan field that the floor reads.

### M8. W-08 step 9 has two outcomes
- Where: W-08:26 "a candidate again; ... the idea owner's earlier decision is in the brief | a candidate row; a finding for the task".
- Source: lens 1; `F-0001#13`, `#28`.
- Why it matters: it is unclear whether the idea owner is asked again (a second human input), or whether the developer gets a finding. Saying that the diff breaks the decision is a reading of meaning, so it cannot be done under a `code` tag.
- Fix: choose one path. For example: code matches the dependency name against the options the decision rejected, stored as a fixed field, and returns a finding; any other candidate is a new brief.

### M9. An escalation decision cannot change a requirement or the band
- Where: :901 "The task starts its next attempt with the decision in its prompt file"; §8:647 "A requirement changes only at a bet".
- Source: `F-0001#26` (scope against date, budget, intent).
- Why it matters: if the idea owner decides "drop REQ-x to keep the date", the requirement stays until the next bet. `layup spec check` and the task register still demand it.
- Fix: state how a confirmed decision updates the inventory or the band, for example as an unscheduled bet line, or name the delay as a limit.

### M10. The frequency parameter conflicts with "the floor always runs"
- Where: :863 "the floor of P1 always runs"; :893 "When (O-79, the frequency is a parameter)".
- Source: O-79 ("before each task only"), O-84.
- Why it matters: if the frequency is set to plans only, it is unclear whether the handoff diff floor still runs. W-08 step 9 depends on it.
- Fix: say that the parameter governs only P1 calls, and that the floor runs on every plan and every handoff diff.

### M11. Some parameter changes are recorded as project input when they change a gate
- Where: :911 "A change is recorded as project-level Operator input (§12), not as input to a task."
- Source: `F-0001#28` ("a change to a gate" is unplanned input to a task); O-69.
- Why it matters: a mid-run change to the transition table (a check) or to a point's authority changes open tasks, yet it would not appear in the Task Intervention Rate.
- Fix: count a change that takes effect on an open task as input to that task, or leave the accounting to slice G without deciding it here.

### M12. The parameter comment has no defined issue and no role check
- Where: :909 "a comment on the target's control issue"; :910 "checks the name, the value and the bound".
- Source: FT3; §3 (approvers.tsv also holds bet approvers that the idea owner names).
- Why it matters: no slice opens a "control issue". Any approver, not only the Operator, could change P1's authority.
- Fix: name the issue and who opens it, and check the author's role for each parameter row.

### M13. The tier rule says "either list" when there are now three lists
- Where: :803 "the reasoning tier when `size` is `large` or either list is not empty" (the line this slice changed).
- Source: operative ambiguity (Bootstrap rule 3).
- Fix: "or any of the three lists is not empty", or exclude `new dependencies` on purpose.

## Notes
- N1. The `decisions.tsv` row has no tokens or price. Jev charges per input token, and a self-hosted Laya has an unknown cost. FT2 needs the row to carry them.
- N2. Pin the model version per threshold. Jev's docs say to pin a versioned ID; a new version should reset `delegate` to `shadow`.
- N3. Sources at :887 include "each pull request description", but the list at :893 omits it. §8 step 4 does not place the screen before the push.
- N4. There is no size check of the state text against Jev's 32k limit for the state plus the question. An oversized state becomes an error, which goes to a human.
- N5. The rows "candidate", "selection" and "decision" name no table. The answer form ("option N; business-forking: yes/no") is not fixed.
- N6. The text does not say where Laya runs (a host or a service) or how its endpoint is registered.
- N7. The `new dependencies` floor escalates every dependency. A deterministic rule on licence and cost would remove many false candidates.

## Checklist rows
- S8: answered, except for M3, M7 and M8.
- I4 (in part): answered, except for M5.
- I6: answered.
- K32: known limit (L-E1), but it understates the case (M3).
- K33: answered.
- K34: not answered (M2).
- K35: not answered (M3, M4).
- K36: answered.
- K37: answered, except for M12.
- K38: not answered (M5).
- K39: answered.
- K40: answered.
- K41: answered (§8).
- P11: answered, except for M10.
- P18 (the smart-if part): answered.
- C3: answered.
- C9: answered.
- C10: answered, except for M1 and M2.
- D02 (in part): answered.
- D04: answered, except for M2.
- D09: answered.
- D10: answered, except for M10.
- D15: answered, except for M11.
- D16 (the bet): answered.
- FT1: not answered (M3).
- FT3: answered, except for M12.

## Existing solutions
- Repeated: the deterministic floor that a model may only raise (Governed APA), the per-rule switch between monitoring and enforcing (numbat), and fail-closed thresholds (jevals, JevLoop). Shadow-then-promote matches ML champion/challenger rollout, which scores against labelled outcomes, not against the default (see M5).
- Ignored: GitHub's dependency-review compare API (`/dependency-graph/compare/{basehead}`, checked on docs.github.com). It returns added dependencies with ecosystem, manifest and licence for any stack. It works on public repositories, and on private ones only with Advanced Security. It needs pushed commits, so it fits after the push, as a second floor.~~~~

### The author's answer to round 1

All thirteen material findings and notes N1 to N7 are applied in one commit (§10 rewritten).

| Finding | Fix |
| ------- | --- |
| M1 | A yes-or-no answer decides at p ≥ t or p ≤ 1 − t; a choice or score at top ≥ t; P1 never removes a candidate. |
| M2 | Failures take the deterministic branch under `off`, `shadow`, `cautious`; only `delegate` sends a failure or an undecided answer to a human. |
| M3 | Each selection row says whether P1 and P2 were active; the Intake form says what `shadow` leaves open; L-E1 covers it; K32 and K35 are answered only at `cautious` or `delegate`. |
| M4 | A question that needs a human is a candidate of the escalation screen, to the idea owner. |
| M5 | Each point names its ground truth; calibration is agreement with it, never with the default. |
| M6 | The floor exempts a rendering of an approved bet or decision and the specification a bet approves. |
| M7 | The options come from the declaration or the plan's `alternatives`, or a Product Owner session writes them. |
| M8 | Decisions record the chosen and rejected options by name; a later candidate that is a rejected option gives a finding, not a new brief. |
| M9 | A decision that changes a requirement, a priority or the band is rendered at once, like a bet line (§8 updated). |
| M10 | The frequency governs P1 only; the floor runs on every plan and every handoff diff. |
| M11 | A parameter change that reaches an open task is input to that task. |
| M12 | `layup run` opens the control issue at Start; each parameter row names the role that may change it. |
| M13 | "any of the three lists". |
| N1 to N7 | Tokens and price in each decision row; the model version pinned, a new one resets `delegate`; the pull-request description dropped from the sources; the state size checked; `candidates.tsv` and `escalations.tsv` named with the answer form; Laya's endpoint in the register; the allowed-licence rule. The dependency-review API is a later option as a second floor after the push (public repositories). |

## Round 2

2026-09-29, head `5551fbe`; 4 min 49 s; 22,512 output tokens; USD 2.80 at list price. Same reviewer and setup. Word for word:

~~~~text
# Slice E review — round 2
Reviewer: claude-fable-5-1. Head: 5551fbe. Verdict: `material findings: 5`.

Round-1 status: M1, M2, M4, M6, M7, M10, M11, M12 and M13 are closed. M3 is closed in text, with a residual in N1. M5 is closed for P1, P2 and P5 only (see M5 below). M8 is closed in intent, but its rule has no key (M2). M9 is not closed (M3).

## Material findings
### M1. A model's own claim can now remove a floor candidate
- Where: architecture.md:905-906, "A dependency whose licence is on the Intake's allowed list and whose cost is zero is not a candidate." ADR-0022:37 says the same.
- Source: §10:914, "a deterministic floor that a model may raise and never lower". FT2 (unknown cost is never zero). I6. PSB root cause (`F-0001#13`).
- Why it matters: the licence and cost of a plan dependency are text that the plan session writes. A `require` in `go.mod` carries no licence and no cost at all. Take W-08's case. The SendGrid Go SDK is MIT and free, and the paid part is the service. A session that writes "MIT, cost 0" for the SDK removes the only floor candidate. The indirect `require` lines that a new module adds have no stated cost, and the text does not say whether they are candidates.
- Fix: code reads the licence from the module itself (as `go-licenses` does), not from the plan. A cost that no record states is unknown, and an unknown cost is a candidate. A manifest dependency that the plan's field does not list is always a candidate.

### M2. The rejected-option match has no defined key
- Where: architecture.md:932-934, "A later floor candidate whose name is an option that a decision rejected gives the task a finding". W-08:24 records the rejected option as "the paid service". W-08:26 (`code`) says "its module name is the option that the decision rejected".
- Source: lens 1 (a `code` step names its input and its rule). `F-0001#13`.
- Why it matters: the floor candidate at step 9 is a module path (`github.com/...`). The recorded option is prose. An exact match fails, so "any other candidate is a new brief" asks the idea owner again. A loose match would be reading meaning under a `code` tag.
- Fix: `new dependencies` names each dependency by its manifest identifier. An option that is a dependency carries that identifier in `escalations.tsv`, and the match is exact equality on it. An option written in prose never matches, so it gets a new brief.

### M3. A decision that changes a requirement or the band has no input that code can render
- Where: architecture.md:930-931, "A decision that changes a requirement, its priority or the band is written at once, as a rendered task (§8), like a bet line". Against that: the answer form at :925-926 (`option <N>; business-forking: yes` or `no`, then free text) and :943-944 ("the idea owner for the band and the intent rows").
- Source: §8:639 (a bet gives "a line per requirement", which code renders from). §8:648 ("A requirement changes only at a bet, or by an escalation decision"). K37, D02. "An operative rule has one home" (Bootstrap mode).
- Why it matters: the idea owner answers "option 2: drop REQ-7 to keep the date". Nothing structured says that option 2 changes REQ-7, or to which priority. Code cannot render a new inventory version from that prose. The band now has three homes: the bet's appetite, an escalation decision, and a `parameters.tsv` row. The parameter table (:950-958) has no band row and no intent row.
- Fix: an option that changes the inventory carries bet-form lines (`REQ-7 Won't`, `band ...`), written by the Product Owner session and checked by code. Then name one home for the band, and add its row to the table.

### M4. The "needs a human" flag of a session is in no schema
- Where: architecture.md:917-918, "or a session sends a question with that flag, the question becomes a candidate". Against that: :849 (P2's deterministic branch is "needs a human = no"), §8:733 (a question is "its text, its own label for the kind, and the records it concerns") and §8:723 (the transition row asks for no flag).
- Source: K35. Lens 1.
- Why it matters: under `shadow`, the default the Intake form offers, a session's flag is the only way a question can reach the idea owner. But the handoff schema has no field for it, and the P2 table says the deterministic branch ignores it. Code cannot read a field that does not exist.
- Fix: add the flag to the question's fixed fields in §8. Make P2's deterministic branch "the asker's label and flag". Or state that the flag is the status `decision_needed`.

### M5. Calibration against the ground truth is undefined for P3 and P4
- Where: architecture.md:851, P4 ground truth "the verification's first-round verdict". :850, P3 "the stall's outcome". :876-878, "the share where the provider's decided answer agreed with that point's ground truth".
- Source: K38 ("promotion has no ground truth"). Invariant 4.
- Why it matters: under `shadow`, the routing order runs pair Y and the provider picked pair X. Y's verdict says nothing about X, and a pass/fail verdict cannot "agree" with a pair choice. P3 has the same problem for a rung that was not taken. So a `delegate` threshold for P3 or P4 has evidence only by name. K38 is not answered for these two points.
- Fix: count only the rows where the provider's answer equals the branch taken, and record the bias this causes. Or record as a known limit that P3 and P4 thresholds rest only on the Operator's comment.

## Notes
- N1. A `candidates.tsv` row is written only "When one is selected" (:920-921). So "Each selection row records whether P1 and P2 were active" (:887) says nothing about a screen that selects nothing. Those are the screens where an inactive P1 matters, and the Missed Escalations audit needs them. Write one row per screen.
- N2. At a `delegate` P2, a failure goes to the Operator (:871). A business-forking question can then reach the Operator, although :873 says it "always goes to the idea owner". Say the Operator only classifies it, and a "needs a human" question goes on to the screen.
- N3. "the tokens and price that the provider reports" (:841). Jev's docs give a price per input token, but the sources do not say its response reports tokens. A self-hosted Laya has no price. Record "unknown", never zero (FT2).
- N4. "the Shape specification that a bet approves" (:907-908). The spec's handoff comes before the bet. Say "a specification task of the Shape phase", so code can apply the exemption at the handoff.
- N5. The control issue is opened "at Start" (:941), but §5's Start steps open only the Intake issue.
- N6. §3:126 gives `approvers.tsv` no role column, but §10:943 now reads a role from it. Since the Operator and the idea owner can share one ID, the file needs a role per row.
- N7. ADR-0022:62 says "under `off` or `shadow` at P1". §10:885 and L-E1 say P1 and P2. Align them.
- N8. It is unclear whether a rejected option binds only its own task or the whole target. If it binds the target, a later task's real need gets a finding, not a brief, unless the session declares it.
- N9. "A change that takes effect on an open task is input to that task" (:946) does not say whether that input is planned or unplanned. `F-0001#28` counts "a change to a gate" as unplanned. Leave the rest to §12.

## Checklist rows
- S8: answered, except M1 and M2.
- I4 (in part): answered, except M5.
- I6: answered, except M1.
- K32: known limit (L-E1, acceptable).
- K33: answered.
- K34: answered.
- K35: known limit (acceptable), except M4.
- K36: answered.
- K37: answered, except M3 (the band and intent rows).
- K38: not answered for P3 and P4 (M5).
- K39: answered.
- K40: answered.
- K41: answered (§8).
- P11: answered.
- P18 (the smart-if part): answered.
- C3: answered.
- C9: answered.
- C10: answered.
- D02 (in part): answered, except M3.
- D04: answered.
- D09: answered.
- D10: answered.
- D15: answered.
- D16 (the bet): answered.
- FT1: answered with L-E1; see N1.
- FT3: answered; see N6.
- L-E1: recorded, and now covers `off` and `shadow`.

## Existing solutions
- Repeated: Governed APA's "raise, never lower" floor; per-rule monitor and enforce modes; Jev's advice to pin a versioned model.
- Ignored for M1: `actions/dependency-review-action` (`allow-licenses`) and Google's `go-licenses` read a dependency's licence from the module itself, not from an agent's claim.
- Ignored for M5: off-policy evaluation (logged propensities, as in contextual bandits) is the standard answer when a shadow policy's pick was never run. Champion/challenger rollouts score only the traffic that the challenger actually served.~~~~

### The author's answer to round 2

The five new material findings and notes N1 to N8 are applied in the next commit; N9 is left to slice G, as the reviewer suggests. Under O-98, these fixes are a named item of the whole-design review.

| Finding | Fix |
| ------- | --- |
| M1 | Every new dependency is a candidate unless the idea owner put its identifier on the allowed-dependency list; nothing a session writes removes a candidate; code reads the licence from the module for the brief. |
| M2 | Dependencies and dependency options carry their manifest identifier; the match is exact; a prose option never matches. |
| M3 | An option that changes the inventory carries bet-form lines that code checks; the band has one home, `budget.tsv` (§12). |
| M4 | A question has the fixed field `needs_human`; P2's deterministic branch reads it. |
| M5 | Known limit L-E2: thresholds for P3 and P4 rest only on the Operator's comment. |
| N1 to N8 | One `screens.tsv` row per screen; the Operator only classifies the kind on a P2 failure; unknown prices are `unknown`; "a specification task of the Shape phase"; the Start opens the control issue; `approvers.tsv` has a role per row; ADR-0022 names P1 and P2; a rejected option binds the whole target. |
