# T-sdtb — the demo of M2d on a test target (uat)

Issue: [#284](https://github.com/pharzam/layup/issues/284), row 59b of the plan: the demo of `M2d`, a uat on a real target, as `T-fnsr`, `T-x7cs` and `T-ekh5` were for `M2a`, `M2b` and `M2c`. Serves `F-0003#42` and `F-0003#43` through REQ-002 and
REQ-003. Base `6c9c314`, then `main` at `461d47a` (the branch took `main` by merges while three fixes landed). Author: Claude Sonnet 5.5 on Claude Code. Evidence: [`runs/T-sdtb/`](../../runs/T-sdtb/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #284. The plan review (GPT-6.1 Sol on Codex CLI) gave `approve-with-conditions`: 900 lines over 24 files, cap 1, no panel. Its six conditions are applied: the coverage map of the
three acceptance criteria; the refusal read independently of the binary, by a transparent wrapper of `git` that keeps each push's exit status and standard error (`git-tee.sh`, tested by `git-tee-test.sh`); the read command, branch and saved
output of each of the five rules; the inputs, the rules and the scenario pushed before any run or any change of a register (step T, `00c9cdf`); D5 rewritten with the failure rule (a failed attempt is kept in the record and told apart from
the passing one; a fix of the check changes it and its specification together; the replay is named in advance); and the acceptance by the Operator on the evidence at a named commit. The plan is reasoning-tier work that an
execution-tier model wrote under the Operator's routing message of 2026-10-13: a routing mismatch by choice. The goal count is 1.

## What was done

Target `pharzam/layup-uat-m2d`, made by Start from the host with a short brief, two registered harnesses (Claude Code and Devin). Runs 0 to 2: Start, the Intake batch of 169 questions, the follow-up of 150; the Operator answered both from
their own account. **Three defects of the build were found, each a task of its own** (the demo's "After a failed attempt"): before any run, GitHub's refusal of a push names the ref and the rules and no ruleset (`T-gkvi`, #323: the probe is
attributed by the ruleset's ID); runs 3 and 4 failed at the step `setup` (S08), because the prompt of the prose session gave none of the rules of the checks that accept its files (`T-92lm`, #326, with a real Opus session as evidence);
and run 6, the first forge check on the real forge, failed `rules:layup-records`, because GitHub accepts `update_allows_fetch_and_merge` `false` and returns the rule with no parameters (`T-j7c5`, #330). Run 5 passed the setup and the
verify with the fixed prompt; the Operator ran the printed `commands.sh` (07:56Z); run 7 (09:20:52Z to 09:51:14Z, 30 minutes of it the wait for the lease that the failed run 6 left) passed the forge check, eight rows `pass`; run 8, a later start (09:51:25Z to 09:51:39Z), passed it again; the Operator accepted the scenario at 09:56Z on the evidence at `45ec2f1`.

**Evidence:** [`runs/T-sdtb/`](../../runs/T-sdtb/): the README (the rules, the scenario, the inputs, the failed attempt, the results), `runs.txt` (each run's table), `forge-check.txt` (each table with its commit and run), `rulesets.txt`,
`probe-and-actors.txt` (GitHub's refusal as the wrapper kept it, the rule-suites read, the rulesets' history, the actors of each write), the brief and the wrapper. **Off the path:** #325 (the review marks `gap` answers open, which made the
follow-up of 150 questions) and #318 (a restart after a failed rebuild), opened as issues. **Known limits:** one target at one time; a failed run leaves the lease held and the next restart waits 30 minutes (`lease.H` is a Start setting);
the passing attempt used the binary built after the three fixes (`8b76fdc`), not the one that row 59a reviewed (`0a104b8`): the README records the diff and the release check.

## Review rounds

Round 1 (`6ef92ae`; GPT-6.1 Sol on Codex CLI, `xhigh`; ADR-0026 does not name it, the Operator allowed it; cycle 0 of the cap of 1; comment 6189585059): `nothing material in scope`, with two notes, both applied in this close-out. The README promised
the Operator's reply about the brief and they gave none, so it now says what they wrote before Start and that the short brief stands; and it says why the integration and e2e levels were not run again here (the branch has no Go change) and where they run
(the CI of the pull request). The reviewer agreed with the classification of the three defects as revealed and off the path of the artifact. The reply is comment 6189601482.

## Verdict

Delivered: the demo of `M2d` passed on the real target `pharzam/layup-uat-m2d`: the five rules hold (the forge check with its eight rows in run 7, a later start in run 8), the Operator accepted the scenario, and the three defects of the build that the demo
found were fixed first, as `T-gkvi`, `T-92lm` and `T-j7c5`. The review ended at `nothing material in scope` in one round, at the cap of 1. This is the last task of `M2d`.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-14 and 15, from the clock and the comments of #284; the tokens of GPT-6.1 Sol are the "tokens used" of Codex CLI; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning (a routing mismatch) | Claude Sonnet 5.5 | `max` | not reported | 11:18 to 11:23 on 2026-10-14 |
| The plan review | reasoning | GPT-6.1 Sol, Codex CLI | `high` | 152,031 | 4 min 12 s, 11:23:16 |
| Step T, the wrapper and its test, the brief, the inputs | execution | Claude Sonnet 5.5 | `max` | not reported | 11:29 to 13:32 on 2026-10-14 |
| The runs 0 to 8, the answers, the evidence, the documents, the acceptance (the three defect tasks have their own records) | execution | Claude Sonnet 5.5 | `max` | not reported | 13:32 on 2026-10-14 to 10:07 on 2026-10-15, with the waits for the Operator, the leases and the fixes |
| The 13 real sessions of the runs (`telemetry.tsv` of the target): probes, Intake, the setup | execution | Claude Opus 5.5 and Sonnet 5.5, Devin SWE-2 | the register's | Claude output 187,043; Devin none reported | 1,803 s in all |
| Round 1 | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 163,106 | 7 min 2 s, 09:59:28 |
| The close-out | execution | Claude Sonnet 5.5 | `max` | not reported | 10:07 to 10:08 |
