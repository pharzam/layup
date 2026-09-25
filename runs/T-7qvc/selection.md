# T-7qvc — the compared set and the selection (solution selection over the panel's options)

Panel members A (CLI and Go architecture; GPT-6 Sol on Devin), B (Git-native state, the forge and CI; Claude Opus 5.5 on Devin) and C (governance, the invariants and the human decision points; Claude Opus 5.5 on `claude -p`, after the OpenCode run failed) gave options, no votes (`panel-A.md`, `panel-B.md`, `panel-C.md`). The Operator answered the panel's questions in one batch (O-52 to O-62 on #66). The considerations are those of Solution selection: Determinism first, then infrastructure, stack fit, and the invariants `F-0001#1`–`#9`. Marks: ✓ serves, ✗ strains, — neutral.

## Q1 — the runner for stack gates → ADR-0013

| Option (members) | Inv-1 record | Inv-2 independent | Inv-3 unforgeable | Inv-5 no-report blocks | Infrastructure | Selected |
| --- | --- | --- | --- | --- | --- | --- |
| A GitHub App that posts required check runs (B 1A, A G1, C 1-B) | ✓ with a row in the target's `runs/` (O-54) | ✗ until the required checks are removed from the protection body (a kit record the target's Operator edits) | ✓ the App identity | ✓ by the forge's own rule | a small receiver, outside the engine | **yes** (O-52) |
| Actions in LAYUP's repository posting a commit status (B 1B, A G2, C 1-A) | ✓ with a row | ✗ same | ✗ any write token posts a status; O-53's identities reduce but do not remove it | ✓ | none | no |
| A local pre-merge step and `layup merge` (B 1C, C 1-C) | ✓ | ✓ | ✗ advisory: `gh pr merge` bypasses it | ✗ | none | no |

Rejected because: the status is forgeable and the local step is advisory (REQ-007 needs the forge to block). Accepted cost: the receiver is a network component; ADR-0011 decision 1 keeps the *engine* free of it, and ADR-0013 places the receiver outside the engine.

## Q2 — the rule-protection control → ADR-0014

| Option (members) | Refuses before merge (REQ-003) | Works with the identities of O-53 | Writes into the target | Human load | Selected |
| --- | --- | --- | --- | --- | --- |
| `CODEOWNERS` + required code-owner review (A P1, B 2A, C 2-A) | ✓ | ✓ | a `CODEOWNERS` file: excluded by O-60 | one approval per rule change | no (O-60) |
| A rule-guard required check by the App, satisfied by the Operator's approving review or signed commit (B 2B, A P2) | ✓ | ✓ (the review comes from the human identity) | nothing | one approval per rule change, a planned approval point (O-61) | **yes** |
| A push-time forge ruleset (C 2-B) | ✓ | ✓ | nothing | none; agents cannot propose a rule change | no: forge-plan dependent, not general across forges (`F-0003#67`) |
| Detection only (A/B/C 2C) | ✗ | — | nothing | an audit | kept as the complement (the §7.1 measure), not as the control |

## Q3 — the role model and the handoff schema → ADR-0015

| Option (members) | Traces to the kit | Harness-neutral (Inv-9) | Role list has evidence (Inv-4) | Adds to the kit copy (O-10) | Selected |
| --- | --- | --- | --- | --- | --- |
| The kit's gate roles; one TSV event table pointing at kit records (B 3A, C 3-A) | ✓ | ✓ | ✓ the kit's steps | no | **yes** (O-55) |
| Four delivery roles (A H1) | — | ✓ | ✗ a hypothesis | no | no |
| Artifact owners (A H2) | — | ✓ | ✗ a hypothesis | no | no |
| A Markdown record per handoff with a new kit linter (B 3B) | ✓ | ✓ | ✓ | a linter in the kit copy | no (O-10) |
| The forge as carrier (B 3C) | ✗ | ✗ | — | no | no (Inv-1) |

Member C's 3-C (the schema is a rule, the role names are content) is kept: a target's roles register is content under ADR-0015.

## Q4 — the escalation rule → ADR-0016

| Option (members) | Deterministic (Inv-6) | Catches a fork no file shows | False positives (TIR) | Selected |
| --- | --- | --- | --- | --- |
| A deterministic floor over the approved-intent files and bounds (B 4A, A E2, C 4-B) | ✓ | ✗ | low | part of the choice |
| A declared class on four axes, enforced by the engine (A E1, B 4B, C 4-A) | partial | partly | low; under-declaration is the risk | part of the choice |
| Both, stacked; the stricter result wins (B "4A and 4B can stack", C 4-C) | floor ✓ | ✓ | higher; each unconfirmed escalation is unplanned input | **yes** (O-56; O-62: when unsure, escalate) |

## Q5 — the stall procedure → ADR-0017

| Option (members) | Deterministic | Harness-neutral | Catches an oscillating loop | Values before the pilot (Inv-4) | Selected |
| --- | --- | --- | --- | --- | --- |
| Fixed limits: N identical-input retries or T minutes without a progress event; disagreement at the cycle cap (A S1, B 5A, C 5-A + 5-B) | ✓ | the time limit measures the harness's speed too (C) | by the time limit | working hypotheses the pilot resets | **yes**, with O-58's values (N = 1, T = 15 min) |
| Pre-registered per-step limits (A S2) | ✓ | ✓ | by the per-step time | each a value to evidence | no: more values to set badly; the ceiling of S2 is kept as a rule |
| A progress metric (B 5B) | ✓ | ✓ | ✓ | — | no: false stalls on a real refactor; its "progress event" definition is kept |
| A cost limit from telemetry (C 5-C) | ✓ | ✓ | — | no baseline exists | no, until the pilot baseline |
| The Operator as examiner (B 5C, C (iii)) | — | — | — | — | no: REQ-010 needs a fresh-context diagnosis before the Operator; the fallback when no examiner is available is "diagnosis pending" in the package |

Examiner: a fresh session that did none of the work, a different model, a different harness where available; the same harness allowed and recorded (O-59).
