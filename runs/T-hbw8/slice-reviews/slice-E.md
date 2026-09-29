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
