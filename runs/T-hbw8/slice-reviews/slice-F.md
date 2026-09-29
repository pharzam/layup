# Slice F review record

Reviewer: Claude Fable 5.1 (`claude-fable-5-1`), a fresh `claude -p` session in a disposable clone, read-only tools, no forge token, `--setting-sources project,local`. Not a gate round.

## Round 1

2026-09-29, head `630d9c5`; 5 min 6 s; 23,150 output tokens; USD 3.16 at list price. Word for word:

~~~~text
# Slice F review — round 1
Reviewer: claude-fable-5-1. Head: 630d9c5. Verdict: `material findings: 9`.

## Material findings
### M1. The progress rule counts net totals, so normal review work stalls and an edited test counts as progress
- Where: architecture.md:990-991 "makes progress only when the unknowns go down or the passed tests go up"; :989 "the plan's frozen test list (§8)".
- Source: PSB `F-0001#37` and `F-0003#49` (a stall is "a step that repeats without progress"). K42. §8 never defines a "frozen" list: "frozen" appears only at :989.
- Why it matters: (a) The verifier's first review raises the unknowns from 0 to 1. That is no progress, so with `stall.N` = 1 it is a stall. W-09 step 1 does not treat it as one. (b) The developer fixes finding A and the verifier then finds a new finding B. The count stays at 1, so this is a stall, but nothing repeated. (c) A test is known only by its ID. A session that weakens a failing test so it passes raises "passed tests", and that resets the loop. This is the Sol-22 class again.
- Fix: count progress by set, not by total: at least one unknown from the last round is closed, or a test newly passes. Let the first review set the baseline. Say where §8 freezes the list, and identify each test by ID plus content hash taken at plan approval.

### M2. §8 and §9 point to §11 limits that §11 does not define
- Where: architecture.md:686-687 "past the cap the round limit of §11 opens a stall"; :742-743 "does not count toward the attempt limits of §11; §11 bounds the questions of a task"; :812-813 "the step fails, and §11 handles it".
- Source: §11 has only the no-progress trigger. It has no attempt limit and no rule for the target's cycle cap. The selection row 0023 cites "round and attempt limits (ev01)".
- Why it matters: a review loop can keep making progress (findings 3, 2, 1, 2 ...) past the target's cycle cap. No trigger in §11 fires, so the loop runs until the milestone cap. A step that fails because its context does not fit has no path either.
- Fix: add an attempt limit and a cycle-cap trigger to §11 as parameters, each with a default and its evidence. Say what a cleanly failed step does.

### M3. The breaker has a model decide uphill or downhill, but §11 says code computes it
- Where: :1055 "An examiner diagnoses each open task as uphill or downhill", against :986 "Progress is computed, never reported" and :993 "A task whose unknowns are zero is downhill".
- Source: D16 (the computed hill), O-90, `evaluation/19-shape-up.md`:103.
- Why it matters: code can count one task as uphill (an open finding) while the examiner calls it downhill. P5's extension depends on which value wins, and the text does not say.
- Fix: code computes the hill position for the extension test. The examiner only gives the cause.

### M4. Nothing says what happens when the examiner fails
- Where: :1025-1026 "No stall goes on without a diagnosis, the Operator's branch included".
- Source: `F-0001#14` (at the limit the Operator gets the package); `F-0003#61`; FT1.
- Why it matters: in W-09 part 2 the harness is the fault. The examiner session can hang, return an invalid form, or find no admitted harness. Its hang is trigger 2, which opens a stall that also needs a diagnosis. The stall never reaches the Operator.
- Fix: an examiner failure (with a limit) writes a "diagnosis failed" row and sends the package to the Operator. Stall Diagnosis counts that row as missing, not as zero.

### M5. W-09 steps 7 and 13 use ladder rules that §11 does not state
- Where: W-09:27 "the retry rung is used up by step 2's round, so the panel"; W-09:38 "with one other harness only, the panel's quorum cannot be met (`insufficient panel`)".
- Source: :1028 "**retry** once with the diagnosis in the prompt". Step 2 happened before any diagnosis existed, so it is not that retry. §11 does not exclude the stalled harness from the panel. H1 plus H2 makes two harnesses. ADR-0023:67 says the rung is "skipped", while §11 says the quorum check returns `insufficient panel`.
- Why it matters: under `shadow` the deterministic ladder gives a retry at step 7, not the panel. At step 13 code needs a rule it does not have.
- Fix: change W-09 step 7 to the retry, or change the ladder rule in §11. State whether the stalled harness may sit on the panel, and whether "skipped" or "insufficient panel" is the rule.

### M6. P3 offers an "examiner" branch that the ladder does not have
- Where: architecture.md:850 "P3 stall action | retry, examiner, panel, or the Operator (§11)", against :1027-1029 (retry, panel, Operator, with the examiner always first).
- Source: ADR-0021 (a closed list of named points and branches).
- Why it matters: at a `delegate` P3, the answer "examiner" has no rung to run.
- Fix: make the P3 options `retry | panel | Operator` after the diagnosis.

### M7. The dead-man job breaks key custody, cannot measure heartbeat age, and can stop running without anyone noticing
- Where: :1066-1068 "with an installation token of the LAYUP App, and opens an issue on a target whose heartbeat is older than `3 × lease.H`"; L-F1 at :1176.
- Source: (a) ADR-0014 decision 3 says the private key is one "which only `layup run` holds". §2:59 says the runner writes "through `layup run` only". Per GitHub docs, an installation token needs a JWT signed with the App's private key and expires after 1 hour. So a scheduled workflow must keep the key as a secret. (b) §2:62-68 makes the heartbeat a counter and says hosts' clocks are never compared, so a stateless job has no age to read. (c) GitHub docs: "In a public repository, scheduled workflows are automatically disabled when no repository activity has occurred in 60 days", and queued jobs "may be dropped" under high load.
- Why it matters: the job "writes nothing else", so a public control repository has no activity. After 60 days the job is disabled, no issue is ever opened, and silence looks like a healthy host (FT1). L-F1 names only "late".
- Fix: allow a second key holder, or use a separate App, and amend ADR-0014. Define the age as the forge's push time of the last records update, taken from the activity API. Record "disabled" and "dropped" in L-F1, or use a push-based heartbeat with an outside monitor.

### M8. The stall answer `set <parameter> <value>` skips the parameter rule of §10
- Where: :1044 "`set <parameter> <value>`" on the task's issue, against :955-961 "a comment there [the control issue], in the form `set <name> <value> because <reason>` ... checks the author, the name, the value and the bound".
- Source: Invariant 4 (`F-0001#4`), Fable-M22.
- Why it matters: a value set from a stall answer has no evidence and no role or bound check. The answer is also on a different issue.
- Fix: reuse the §10 form and its checks, or send `set` to the control issue.

### M9. The orientation allowance goes past the Operator's maximum for T
- Where: :1012-1013 "the first session of a task gets `stall.T` once more".
- Source: O-82 "T = 10 minutes (a maximum)".
- Why it matters: a first session can then run 20 minutes with no output before it counts as a hang, which is twice the maximum the Operator set.
- Fix: keep the allowance inside `stall.T`, or ask the Operator on #72 and record the answer.

## Notes
- N1. `ci.T` has no default and no evidence. §10:971 points to "§11" for "the other stall limits", but §11 gives none.
- N2. `stop <task>` does not say what happens to the task's `Must` requirements. Dropping one is a scope trade (`F-0001#26`), which belongs to the idea owner.
- N3. §11 does not say whose reply counts at the Operator rung (an Operator role in `approvers.tsv`?).
- N4. The hang trigger kills a session whose silent command (a long test suite) runs past `stall.T`. The harness register (§9) has no field that says whether a harness sends hook events.
- N5. For an orchestrator stall (trigger 5), the retry and panel rungs have no clear meaning. Say which rungs apply.
- N6. The dead-man job opens an issue on each scheduled run. Say whether it adds to one open issue or opens duplicates.
- N7. The panel's synthesis is a model session. AgentJury's is code. This is a justified change (options, not votes), but the ADR should say so.

## Checklist rows
- S9: answered (see M1, M4).
- R05: answered (the stall, diagnosis and outcome rows).
- K42: answered, but the answer has a new defect (M1).
- K43: answered (see M9).
- K44: answered (a fresh session, not a third harness); the failure path is open (M4).
- K45: answered.
- K46: not answered, the credential and the heartbeat age (M7).
- K47: answered (the early retrospective's batch is slice H).
- K48: answered (reroute; W-12 step 7).
- K49: answered (§10: `lease.H` "set at Intake with evidence").
- P05: answered (trigger 2).
- P19: answered.
- P23: answered (see M3).
- C8 (the circuit breaker): answered (see M3).
- D13: answered (see M9).
- D16 (the breaker and the hill): answered, but the hill by examiner contradicts it (M3).
- FT1: not answered for examiner failure and a disabled or dropped dead-man job (M4, M7).
- FT6: answered.
- L-F1: known limit, but incomplete (M7).

## Existing solutions
- Temporal activity heartbeat timeouts are the proven form of trigger 2. OpenHands' `StuckDetector` (repeated action and observation) matches trigger 1 by identity, not by count (M1).
- Healthchecks.io, Cronitor or Dead Man's Snitch give a push-based dead-man switch without GitHub's schedule limits (M7).
- The design keeps Shape Up's hill and breaker and AgentJury's quorum and `insufficient_jury` correctly, apart from M3.~~~~

### The author's answer to round 1

All nine material findings and notes N1 to N7 are applied in one commit (§11 rewritten).

| Finding | Fix |
| ------- | --- |
| M1 | Progress by sets: a round closes an unknown the last round left open, or a frozen test newly passes; the plan review freezes each test by ID and source hash (§8 step 2); the first review sets the baseline. |
| M2 | Triggers for rounds past the target's cycle cap, attempts past `stall.attempts`, and a step that cannot go on. |
| M3 | Code computes the position for the breaker; the examiner gives only the cause. |
| M4 | A failed examiner writes "diagnosis failed" and sends the package to the Operator; Stall Diagnosis counts it as missing. |
| M5 | W-09 step 7 uses P3 at `cautious`; the panel rung is skipped when fewer than two admitted harnesses are free of the diagnosed failure (W-09 step 13). |
| M6 | P3's options are `retry`, `panel`, `Operator`. |
| M7 | The dead-man job holds a second App key (ADR-0014 amended), reads the forge's time of the last records update, keeps its schedule alive, and L-F1 names late, dropped and disabled runs. |
| M8 | A parameter change from a stall answer goes to the control issue in the form of §10. |
| M9 | The first session gets no time beyond `stall.T`. |
| N1 to N7 | `ci.T`, `stall.attempts`, `panel.K` set at Intake with evidence; a stopped task's requirements go to the next bet; the Operator role; the register's hook-event field; an orchestrator stall has only the diagnosis and the Operator; one open issue for the notices; the ADR says why the synthesis is a session. |
