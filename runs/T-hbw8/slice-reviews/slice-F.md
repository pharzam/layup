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

## Round 2

2026-09-29, head `7398bdc`; 3 min 57 s; 17,715 output tokens; USD 2.39 at list price. Same reviewer and setup. Word for word:

~~~~text
# Slice F review — round 2
Reviewer: claude-fable-5-1. Head: 7398bdc. Verdict: `material findings: 5`.

## Round-1 findings
- M1: the progress rule by sets is fixed. The fix brings a new defect: the test list cannot be frozen at the point it names (M1 below).
- M2, M3, M4, M5, M6, M8, M9: closed.
- M7: half closed. The heartbeat age now comes from the activity API, and L-F1 names late, dropped and disabled runs. The key custody is still wrong (M5 below).

## Material findings
### M1. The test list is frozen before the tests exist, and nothing can freeze it again
- Where: architecture.md:670-672 "The approved plan freezes its test list: each test by ID and the SHA-256 of its source". Also :994-996 "a test whose source changes is no longer on the list as passed, and the change is a finding".
- Source: §8:673 says "A developer session writes the failing test first", which is after the plan review. The transition row at :716 says the plan's tests "fail at the base and pass at the head". §8:747 says "A new attempt starts from the base commit". Checklist row K42.
- Why it matters: at plan approval a new test has no source, so there is nothing to hash. Each attempt starts from the base and writes the test again, so its hash changes every attempt. In W-09 step 8 the panel says "correct the test to the specification". By §11 that change is a finding that nothing closes, so W-09 step 10 ("the verifier passes; the stall closes") cannot happen as written.
- Fix: take the hash when the first handoff passes the transition check (the test fails at the base and passes at the head). Add a rule that freezes the list again, for example an amended plan with a new plan review, and say that it closes the "test changed" finding.

### M2. With progress counted by sets, nothing bounds the questions of a task
- Where: architecture.md:744-745 "An attempt that a question ends does not count toward the attempt limits of §11; §11 bounds the questions of a task". Also :1007 "A question round counts the same way". Also ADR-0023:73 "A question loop and a review loop stop after `stall.N` rounds without progress".
- Source: §11:997-998. A round makes progress when "at least one unknown that was open at the end of the last round is closed".
- Why it matters: a developer asks Q1, which is answered, then Q2, which is answered, and so on. Each round closes one unknown, so each round is progress, and trigger 1 never fires. Trigger 2 skips these attempts because §8 excludes them. Only the milestone cap stops the loop. Before the fix, a net count stopped it.
- Fix: count question-ended attempts toward `stall.attempts`, or add a question limit per task (a parameter with evidence).

### M3. Trigger 5 calls a clean failure a stall
- Where: architecture.md:1016-1017 "A step that fails cleanly but cannot go on, for example a context that does not fit (§9)."
- Source: `F-0001#37`, quoted at §11:984-986: a stall is a task that "does not fail cleanly". The Stall Rate (`F-0003#73`) counts tasks with a stall.
- Why it matters: a task refused for context size is a clean failure. Here it becomes a stall and raises the Stall Rate, which has a start value of 5% or less. The circuit breaker (trigger 6) has the same problem. `evaluation/19-shape-up.md`:73 calls it "a budget stop, not a stall detector".
- Fix: give a clean failure its own row and path (the task fails and the Operator gets a notice), outside the stall count. Or say in the design, with a reason, why the PSB term is widened.

### M4. A task that has not started counts as downhill, so the breaker can grant an extension it should not
- Where: architecture.md:1000-1001 "A task with no open unknown is **downhill**". Also :1076-1077 "When every open task is downhill ... P5 (§10) may grant **one** extension".
- Source: D16 and the selected mechanism in `evaluation/19-shape-up.md`:80: "uphill while an open question, an unresolved material finding or an unfrozen test list exists". Shape Up's extension test needs the remaining work to be execution only.
- Why it matters: a milestone reaches its cap with three build tasks not yet planned. They have no open unknowns, so all three are "downhill", and P5 may extend. The same rule also sends those tasks to the execution tier (§9:805).
- Fix: count a task as uphill until its test list is frozen, as the evaluation row says.

### M5. The second App key is not limited, and §2 still says the runner writes only through `layup run`
- Where: ADR-0014:50-51 "a second key, which only reads and opens an issue, is the dead-man job's". Also architecture.md:1086-1087. Also §2:59, the runner row: "through `layup run` only".
- Source: GitHub docs, "Managing private keys for GitHub Apps": up to 25 keys per App, and no key can be scoped. Any key signs a JWT for the whole App. ADR-0014 decision 3 makes the App "the only bypass actor of the records branch".
- Why it matters: anyone who can edit the control repository's workflow can use that secret to mint a token that writes and rewrites the records of every target, as LAYUP. The words "only reads and opens an issue" describe the job's script, not what the key can do. L-A1 (:1144-1145) mentions where the key lives but not this new exposure. The §2 row contradicts §11, where the job writes a notice and a commit by itself.
- Fix: use a separate App that has only Issues write and Metadata read. Or state the full power of the key in ADR-0014 and L-A1, and have the job request a down-scoped installation token (the `permissions` and `repositories` fields of the access-token call). Update §2:59 either way.

## Notes
- N1. A finding's identity (file, line, rule) breaks when an edit shifts lines, or when a fresh verifier words the rule differently. Either way it counts as progress. Trigger 2 bounds this, but trigger 1 is weaker than the text suggests.
- N2. W-12 step 4 is tagged `code`, but it includes the examiner's diagnosis, which is `model`. Split the step, or point only to W-09.
- N3. §11:1042-1043 "a new trigger in the same task moves one rung up" does not say whether the new trigger needs its own package and diagnosis. W-09 step 13 skips them, but §11:1039 says "No stall goes on without one of the two rows".
- N4. §11:1091 "It writes nothing else" contradicts :1094-1095 "also commits its last run time to its own repository".
- N5. ADR-0023:74 "with one, its rung is skipped" is narrower than decision 4 ("free of the diagnosed failure").
- N6. Say which "one open issue" on the target gets the dead-man notice (the control issue?).
- N7. Trigger 7 (the orchestrator) has no task, yet the Stall Rate counts per task. Say how it counts.
- N8. W-12 step 7 checks that H2 is admitted, but not that the author holds the Operator role (§11:1059).
- N9. Stray indentation at architecture.md:669 and :1144, and at ADR-0023:48 and :52.

## Checklist rows
- S9: answered (see M1, M2).
- R05: answered.
- K42: answered, but the freeze has a new defect (M1).
- K43: answered (the allowance is gone).
- K44: answered (a failed examiner writes "diagnosis failed" and goes to the Operator).
- K45: answered.
- K46: answered in mechanism (run list, credential, write path). The credential's authority is misstated (M5).
- K47: answered.
- K48: answered (W-12 step 7).
- K49: answered.
- P05: answered (trigger 3).
- P19: answered.
- P23: answered; the hill definition is in M4.
- C8 (the circuit breaker): answered; see M4.
- D13: answered (§10 row: 10 minutes as a maximum, and 1).
- D16: answered in form, but the hill rule departs from the evaluation (M4).
- FT1: answered (the examiner failure row; a silent job is L-F1).
- FT6: answered.
- L-F1: known limit (acceptable).

## Existing solutions
- Temporal's activity heartbeat timeout matches trigger 3. OpenHands' `StuckDetector` matches repeated actions by identity; its approach bears on N1.
- `actions/create-github-app-token` takes `permission-*` inputs for a down-scoped token. A separate narrow App is the usual pattern (M5).
- The `keepalive-workflow` action keeps a schedule alive through the API, without commits. Healthchecks.io and Dead Man's Snitch give a push-based switch; L-F1 already names that path.
- The panel follows AgentJury's quorum and `insufficient` result; the synthesis session is a justified change. Shape Up's extension test needs "unfrozen = uphill" (M4).~~~~

### The author's answer to round 2

The five new material findings and notes N2 to N8 are applied in the next commit; N1 stays as the reviewer states it (trigger 2 bounds it); N9: the indentation is the list continuation. Under O-98, these fixes are a named item of the whole-design review.

| Finding | Fix |
| ------- | --- |
| M1 | The list is frozen at the first valid handoff to the verifier; an amended plan with a new plan review freezes it again and closes the "test changed" finding (W-09 step 10). |
| M2 | Attempts that a question ended count toward `stall.attempts`. |
| M3 | A clean failure and the circuit breaker have their own rows and are not stalls; five triggers remain. |
| M4 | A task is uphill until its test list is frozen. |
| M5 | The dead-man job is a separate App, `layup-watch`, with only read and issue permissions; ADR-0014 keeps one key holder; §2 has its own row. |
| N2 to N8 | W-12 step 4 tagged `code`, `model`; a new trigger gets its own stall row and diagnosis; "nothing else on a target"; ADR-0023's consequence matches; the notice goes to the control issue; an orchestrator stall counts per project; W-12 step 7 checks the Operator role. |
