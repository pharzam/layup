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

Accepted: not yet; the Operator's sign-off, with its date, is written here after the run
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

Before `numbering`, run 1 did:

- `probe`: both harnesses probed and passed;
- `rule-gaps`: 1 rule gap, in `intake/psb.tsv`;
- `meaning`: 21 gaps of meaning, in `intake/meaning.tsv`.

The three refusals share one cause, a sentence of the prompt of the numbering that reads two ways. It is a defect of the build of `M2c`, off the path of this task, fixed by `T-pzgw` (#268, merged at `5dbd7af`). Each `-progress.txt` file holds the step lines of a run, without its ten-second beats.

**The binary of run 4 and later:** built on this branch at `6400637`, after the
merge of `main` at `5dbd7af`, `layup 0.1.0-dev`. `git diff --stat 0b881fa
6400637 -- cmd internal go.mod` names three files. The one non-test file is
`internal/run/numbering.go`, the prompt of `T-pzgw`; the two test files are its
test and the empty `notYetBuilt` list of the block test. So the code that row
49a reviewed is the code of this demo, but the prompt.
