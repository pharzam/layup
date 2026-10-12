# T-ekh5 — the demo of M2c (uat)

Row 49b of [the tasks of M2c](../../docs/plan/README.md#the-tasks-of-m2c), issue
[#228](https://github.com/pharzam/layup/issues/228). The inputs are the
Operator's O-201 (comment 6110518260 of #228). The plan, its review and the
author's answer are comments on #228.

## The rules, written before any run

Written on 2026-10-11, before the registers changed and before Start (step T of
the plan, with conditions 3 and 4 of its review). The demo passes when each of
the four conditions holds; it fails when one is false, or when a step of the
restart is `fail` after the batch is posted.

1. The App posts the comment of the batch on the Intake issue, and its first line
   is `LAYUP Intake: the batch of questions`. It holds at least one question of
   the source `rule` and at least one of the source `meaning` (`questions.tsv`
   of the records). A source with no row means the demo is not shown, and the
   record says so (one reading, condition 3).
2. `layup spec check --facts spans.tsv --psb problem-statement.md` exits 0 on
   `spec/spans.tsv` and `start/problem-statement.md` of the records.
3. `spec/check.tsv` of the records has no `fail` row.
4. Each commit of `layup-records`, the Intake issue, the control issue and each
   comment of LAYUP has the App's bot, `layup-agent[bot]`, as its actor; the push
   of the root commit of `main` is the Operator's, by `root-push`. Intake pushes
   no task branch, so no push of one is read.

## The scenario (uat)

```text
Scenario: The demo of M2c on a pilot problem statement

Given the target pharzam/layup-uat-m2c, made by Start from the host with the pilot brief PSB-CHAT-001, and two registered harnesses, Claude Code and Devin, routed for the roles of Intake
When the restarts of layup run run the steps of Intake until the batch waits for its answers
Then one comment of the App on the Intake issue holds the rule gaps and the gaps of meaning of the brief, each with its question and its quote
And layup spec check --facts passes on the numbered spans that the records hold
And each write of LAYUP on the target shows the App's bot as its author

Accepted by the Operator, by Telegram, on 2026-10-12 (19:20Z), recorded in comment 6131023042 of #228
Covers REQ-001, REQ-012
```

The scenario rides the e2e path of row 48 (`TestTheStepsOfIntake`,
[`runs/T-mo9w/test-runs.md`](../T-mo9w/test-runs.md)). The Operator signs it
off, as for `M2b` ([`runs/T-x7cs/README.md`](../T-x7cs/README.md)).

## The inputs

- **The registers** (D1, on 2026-10-11 at 15:38Z): [`registers/`](registers/) holds
  them as the runs read them. The originals are kept on the host in
  `registers.bak-2026-10-11/`.
  - `models.tsv` gains `claude-opus-5-5`, context 1000000, from the models
    overview of the Claude docs.
  - `routing.tsv` gains `systems-architect` and `product-owner` (position 1
    `claude` `claude-opus-5-5`, 2 `devin` `swe-2-high`), and `qa-engineer`
    (1 `devin`, 2 `claude`), all at the tier `reasoning`.
  - The cap of `claude` in `harnesses.tsv` is 5.00 (O-201), and its `wall` 30
    minutes (O-203, at 15:56Z).
  - `forge.tsv` is unchanged; its `watch_slug` is `—`.

  The credential columns hold paths, not secrets.
- **The binary** (D2): `go build -o ~/layup-host/bin/layup ./cmd/layup` on this
  branch at `3f6e910`, `layup 0.1.0-dev`. `git diff --stat 0b881fa 3f6e910 --
  cmd internal go.mod` prints nothing, so its Go code is the release that row
  49a reviewed (`0b881fa`, the merge of row 48).
- **The brief:** `~/layup-host/briefs/psb-chat-001.md`, byte for byte from
  `docs/facts/problem-statement-brief.md` of `pharzam/chat-orchestrator`.
  - SHA-256
    `97c657b6cc38dd2592c84784b0da85af7766a7e31c8cbfe824d67ea794c894fd`, as
    `runs/T-evad/README.md` records it; 1,132 lines.
  - It is not copied here; the records of the target hold it as
    `start/problem-statement.md`.

## The runs

| Run | From | To | Exit | The last row |
| --- | ---- | -- | ---- | ------------ |
| [0](run-0-start.tsv), Start, by the Operator | 2026-10-11 | 20:34Z | 0 | `lease`: released; Intake issue 1, control issue 2 |
| [1](run-1.tsv) | 20:34:32Z | 20:38:09Z | 1 | `numbering`: fail, the file of S-2742cb40 refused at line 2, column `reason` |
| [2](run-2.tsv) | 20:38:53Z | 21:11:18Z | 1 | the same, S-37a7776c, after the takeover of the lease |
| [3](run-3.tsv) | 21:11:48Z | 21:44:13Z | 1 | the same, S-445dee66, whose file was copied before its directory was removed |
| [4](run-4.tsv) | 2026-10-12 05:59:59Z | 06:31:56Z | 1 | `numbering`: fail, the session S-21d94fd8 ended `crash` (the weekly limit of `claude` was reached; comment of 07:07Z on #228) |
| [5](run-5.tsv) | 07:18:21Z | 08:00:54Z | 0 | `answers`: done, waiting for the answers; `numbering` (800 spans, 679 facts), `spans`, `draft` (587 rows), `completeness` (9 items) and `batch` (34 questions) done, after a wait of about 30 minutes (`3 × lease.H`) for the lease that run 4 left held |
| [6](run-6.tsv) | 18:07:38Z | 18:09:06Z | 0 | `follow-up-answers`: waiting; `answers` read 34 answers, `answers-review` reviewed them (4 open), `follow-up` asked 7 questions |
| [7](run-7.tsv) | 18:43:42Z | 18:43:55Z | 0 | `follow-up-answers`: done, 7 answers, `budget.tsv` written, "Intake is done" |

Before `numbering`, run 1 did:

- `probe`: both harnesses probed and passed;
- `rule-gaps`: 1 rule gap, in `intake/psb.tsv`;
- `meaning`: 21 gaps of meaning, in `intake/meaning.tsv`.

The three refusals share one cause, a sentence of the prompt of the numbering that reads two ways. It is a defect of the build of `M2c`, off the path of this task, fixed by `T-pzgw` (#268, merged at `5dbd7af`).

**The binary of run 4 and later:** built on this branch at `6400637`, after the
merge of `main` at `5dbd7af`, `layup 0.1.0-dev`. `git diff --stat 0b881fa
6400637 -- cmd internal go.mod` names three files. The one non-test file is
`internal/run/numbering.go`, the prompt of `T-pzgw`; the two test files are its
test and the empty `notYetBuilt` list of the block test. So the code that row
49a reviewed is the code of this demo, but the prompt.

Each `run-N.tsv` holds the table of the restart, with its times and its exit in the first
line, and not the progress lines and the ten-second beats of standard error. Models are
named by their IDs: `claude-opus-5-5` for the roles of Intake on `claude`, `swe-2-high` for
`devin`, `claude-sonnet-5-5` for the probe of `claude`.

## The result

The four rules written before any run (above), read on 2026-10-12 after run 7.

1. **The batch.** The comment 6120182070 of `layup-agent[bot]` on the Intake issue, posted at
   08:00:52Z by run 5, has the first line `LAYUP Intake: the batch of questions`; its body
   from the forge is [`batch-comment.md`](batch-comment.md). Its 34 questions are the rows of
   the record, [`questions.tsv`](questions.tsv): 1 of the source `rule`, 21 of
   `meaning`, 3 of `intent` and 9 of `coverage`. So the demo holds under the first reading
   of condition 3: at least one question of each of `rule` and `meaning`. The 7 rows of the
   source `follow-up` are those of run 6.
2. **The span check.** [`spec-check.txt`](spec-check.txt): `layup spec check --facts
   spans.tsv --psb problem-statement.md` exits 0 on the copy of `spec/spans.tsv` (800 spans,
   equal to the record byte for byte) and on `start/problem-statement.md`, the brief
   (SHA-256 `97c657b6cc38dd2592c84784b0da85af7766a7e31c8cbfe824d67ea794c894fd`, 1,132
   lines): five rules, five `pass`.
3. **`spec/check.tsv`** of the records has the same five rules, all `pass`
   ([`spec-check.txt`](spec-check.txt)).
4. **The actors,** read from the forge after run 7
   ([`actors-and-telemetry.txt`](actors-and-telemetry.txt)):
   - the 51 commits of `layup-records`: author and committer `layup-agent[bot]`;
   - the one commit of `main`, the baseline `1f0d4de`: author and committer
     `layup-agent[bot]`; `main` was created by the Operator, `pharzam`, at 20:33:47Z on
     2026-10-11, by `root-push` (the `CreateEvent` of the forge);
   - the Intake issue (1) and the control issue (2): opened by `layup-agent[bot]`;
   - the 11 comments of the control issue and the 2 of the Intake issue that are
     LAYUP's (the batch and the follow-up): `layup-agent[bot]`;
   - the other comments of the Intake issue are humans': `pharzam` (2), `yaltunbicak` (2).
   The clause "each push of a task branch" of rule 4 checked nothing: Intake pushes no task
   branch (note 7 of the plan review).

**The demo passed:** each of the four rules holds. The Operator's sign-off of the scenario
of the uat is the last section.

## The telemetry

`telemetry.tsv` of the records, copied at its end in
[`actors-and-telemetry.txt`](actors-and-telemetry.txt), has 11 rows: the probes of `claude`
(`claude-sonnet-5-5`) and of `devin` (`swe-2-high`), and 9 sessions of roles. Six are of
`product-owner` (five of the numbering, one of them the `crash` of run 4, and the draft) and
two of `systems-architect` (the meaning and the review of the answers), all on
`claude-opus-5-5`; one is the completeness review of `qa-engineer` on `swe-2-high`. Their
durations add up to 22 min 25 s. The tokens of the 9 rows of `claude` are `observed`: 162
in, 106,495 out, 4,347,006 of cache; `devin` reports none (`unavailable`). The money of each
row is `unknown`, as the billing is a subscription, so the cap of 50.00 USD of milestone 0
counts no money, by the rule of `intake.md`, The budget of Intake. The hours of the cap,
48.0, ran from the Start (2026-10-11 20:33:58Z, the first commit of the records) to the end
of run 7 (2026-10-12 18:43:55Z): 22.2 hours.

## After the batch (D6, not a criterion)

- **Run 6** read the 34 answers of `pharzam`, the idea owner of record. `answers` did not
  write `budget.tsv`: the three answers of the intent form were `50 USD.`, `100 USD.`
  and `48 hours.`, and the form takes a decimal with a point (`50.0`). `answers-review`
  reviewed the 34 answers and left 4 open, and `follow-up` asked 7 questions in the comment
  6129893541: Q-035 to Q-038 from the open answers, and Q-039 to Q-041, the three of the
  intent form again.
- **The idea owner of the pilot is `yaltunbicak`** (the Operator said so on 2026-10-12, after
  the answers of `pharzam`). The Start wrote `approvers.tsv` with `--idea-owner pharzam`, and the
  spec has no step that changes it, so the comments of `yaltunbicak` (6130283081 at 18:33:57Z
  and 6130399762 at 18:41:00Z) were copied as input and answered nothing. The Operator
  approved the answers of the comment 6130283081 in the session at 18:37Z, and, as the
  approver of record, posted the comment 6130363413 at 18:38:39Z: a quote of those answers
  with `> ` before each line. The parser reads a line only when it starts with `Q-NNN: `, so
  that comment held no answer. The Operator edited it at 18:42:58Z, before any restart had
  copied it, to the answers `Q-035: ` to `Q-041: ` (the last three `Q-039: 50.00`,
  `Q-040: 100.00` and `Q-041: 48.0`) and the word `Approved`. A comment edited before its
  first copy is the decision as copied (`docs/spec/run.md`, A human decision).
- **Run 7** copied the three comments, read the seven answers of the Operator's comment,
  wrote `budget.tsv` (`band.B` 50.00, `band.U` 100.00, `appetite.hours` 48.0) and ended with
  "Intake is done". Q-036 and Q-037 were answered `BLOCKED`: the supplier of the contract of
  the existing KB and the allowlist of roles are inputs that the idea owner has not
  supplied.

## What the demo found

- **#268** (merged by `T-pzgw`): a sentence of the prompt of the numbering that read two ways
  refused the file of runs 1 to 3.
- **The `crash` of run 4** was the weekly limit of `claude`, a limit of the host's
  subscription, and it left the lease held: run 5 waited 1,792 s, the `3 × lease.H` of the
  rule, with `lease.H` of 10 minutes in `start.tsv`, to take it over.
- **Issue #300:** `answers` says "budget.tsv not written (a field of the intent form does
  not parse)" and names neither the field nor the form; the answers `50 USD.`, `100 USD.`
  and `48 hours.` cost a restart and a follow-up.
- **The approvers of a Start cannot change.** The Start named the Operator's account as the
  idea owner, and the real idea owner answered later from their own account; the Operator
  had to confirm the answers as the Operator's own. This is the design of `docs/spec/run.md`
  (the approvers are those of `approvers.tsv`), not a defect that the demo found; it is
  written here so that a later Start names the idea owner.

## The sign-off

The scenario above, "The demo of M2c on a pilot problem statement", rests on the evidence
of this directory. The Operator accepted it on 2026-10-12 at 19:20Z, by Telegram, on reading
this directory at commit `52b3311`; the comment 6131023042 of #228 records the answer. The text
of the scenario is the one written at step T, before any run, and only its `Accepted:`
line has changed since.
