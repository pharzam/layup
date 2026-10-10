# T-x7cs — row 39b of the plan, the demo of M2b (uat)

Issue: [#170](https://github.com/pharzam/layup/issues/170), row 39b of [the
tasks of M2b](../plan/README.md#the-tasks-of-m2b), opened by `T-fdaq` (#152).
Serves `F-0003#52` and `F-0003#50` through REQ-013, REQ-003, REQ-005, REQ-011
and NFR-001. Base `33439ce` (the merge of #198, row 39a); the budget is read
against it. Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-x7cs/`](../../runs/T-x7cs/).

## Plan and plan review

The Operator's go-ahead (Telegram, 2026-10-10 08:47Z) let the author draft the
inputs. The plan (R12), its review and the author's answer are comments on
#170. The plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only
session in a clone at `33439ce`, 12 min 5 s) gave `approve-with-conditions`:
Budget maximum 900 lines added plus removed over 12 files, close-out inside;
Cycle cap 1; no panel. Its four conditions are applied: Devin's file holds a
static API key and no refresh token (read by its key names); the registers are
drafted before the first run; the full model IDs; the reading of criterion 2
went to the Operator (comment 6095937983), and the Operator's acceptance of the scenario settles it (comment 6099346566). The goal
count is 1, final by O-187.

## What was done

1. **The uat test** (D1): `TestTheDemoOfM2b` in
   [`internal/run/demo_uat_test.go`](../../internal/run/demo_uat_test.go),
   under the build tag `uat`, which no CI job runs. It reads the host's
   registers with the readers of `internal/cli`, runs the restart's steps on
   one run (`forge`, `clone`, `version`, `lease`, `probe`), commits the
   rule-path register when the records hold none (D2), runs one developer
   session of a fixed task with the attempt, the pair, `Admit` and the push of
   rows 30b to 38 (D3), asserts the probes, the events, the branch and the
   telemetry row (D4), then runs `phase`; a stop before `phase` still
   releases the lease. No non-test Go file changes.
2. **The inputs** (D5, D6): the host `hetzam`, `~/layup-host`, with two
   harnesses, Claude Code (`claude-sonnet-5-5`, the token of `claude
   setup-token`) and Devin (`swe-2-high`, its `credentials.toml`); the
   Operator made both credential files and the target.
3. **The defect #199**, revealed: the restart refuses a target whose harness
   register gained a row after Start, so `pharzam/layup-uat` of `M2a` could not
   be the target. The Operator chose a new target, `pharzam/layup-uat-m2b`,
   started from `hetzam` (option (a) of comment 6095966713).
4. **The demo**, green on 2026-10-10 at 14:58Z to 14:59Z: both probes passed;
   the developer session `S-a0460618` of `T-rmgw` ended `done`, its commit
   `dddf30c` landed on `task/T-rmgw/1` with the events `attempt`, `session`,
   `result`, `push` and `bound` and its telemetry row; each records commit,
   comment, issue and push of LAYUP is the bot's, the push of `main` the
   Operator's by `root-push`
   ([`README.md`](../../runs/T-x7cs/README.md)).
5. **The documents** (D7): the §12 cells of the five rows and a §13 line of
   `PRD-0001`; the uat row of `traceability.md`; the tag `uat` in
   `test-levels.md` and `engineering-discipline.md`.

**The milestone sentence as read:** `docs/plan/README.md` asks "one complete
telemetry row"; the issue, the acted-on text, says "with its telemetry row".
The developer session's row has its tokens `observed` and its money `unknown`
(a subscription, no `prices.tsv`).

**Tests:** [`test-runs.md`](../../runs/T-x7cs/test-runs.md): run 1 refused at
`clone` (#199); run 2 green. **A deviation:** the red of condition 2 (each
probe failed with no credential) was not run, for the reasons the record
gives.

**The rejected alternatives:** a hook in `Config` between `probe` and `phase`
(product code for a demo, and the seam of `M2e`); copying
`~/.claude/.credentials.json` (a refresh token that rotates).

## Review rounds

The record is a comment on #170. Round 1 (Claude Fable 5.1, effort `xhigh`, a
fresh read-only session in a clone at `77daba6`, cycle 0, the cap):
`nothing material in scope`, seven notes. Applied in the close-out: note 1 (the
record of the deviation says that a red before the green was possible); note
3 (the evidence names L-A1 for both credentials, and that the token stays the
host's until revoked); note 4 (the Operator's acceptance, below); note 5 (a §2
lesson of #199, and #199 in the open issues of the plan); note 6 (the limit of a
test that no job compiles, in `test-levels.md`); note 7 (this close-out).
Declined: note 2 (the test compares the branch's head with the SHA of the
event `push`): a change of the test after the last round is read by no round;
the evidence checks the head, its parent and its file by hand, from the
records.

## Verdict

Delivered: the demo of `M2b` on real harnesses, by the uat test of the tag
`uat`: the probe of Claude Code and Devin, then one developer session whose
commit lands on `task/T-rmgw/1` with its telemetry row, each write of LAYUP by
the App's bot. The Operator's acceptance: on 2026-10-10 at 15:52Z, by Telegram, recorded in comment 6099346566 of #170; it settles condition 4 (the reading of criterion 2) and rules the red of condition 2, not run, acceptable for this uat. The review ended by decay
at cycle 0 of cap 1. The diff against `33439ce`, the branch's base, is inside
900 lines over 12 files. This is the last row of `M2b`; #199 stays open for a
later milestone.

## Resource record

Recorded, not budgeted (ADR-0007). UTC; reviewers' tokens are `modelUsage`;
the author's are not reported. The demo's own sessions are in its telemetry
rows ([`README.md`](../../runs/T-x7cs/README.md)).

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 2026-10-10 08:48 to 08:50 |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,689,783 (USD 8.71) | 12 min 5 s |
| The registers and the uat test; run 1; #199 | execution | Claude Opus 5.5 | not reported | not reported | 09:03 to 09:10 |
| The inputs with the Operator; Start; the demo | execution | Claude Opus 5.5 | not reported | not reported | 14:43 to 14:59 |
| The evidence and the documents | execution | Claude Opus 5.5 | not reported | not reported | 14:59 to 15:03 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 1,548,845 (USD 5.69) | 11 min 28 s |
| The close-out, with notes 1 and 3 to 7 | execution | Claude Opus 5.5 | not reported | not reported | 15:16 to 15:20; 15:53 |
