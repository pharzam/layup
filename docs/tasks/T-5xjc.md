# T-5xjc — row 59a of the plan, the review of the release of M2d

Issue: [#283](https://github.com/pharzam/layup/issues/283), row 59a of [the tasks
of M2d](../plan/README.md#the-tasks-of-m2d), opened by `T-bng8` (#267). Serves
`F-0003#43` through REQ-015 and REQ-017. Base `0a104b8` (the merge of #320, row 58,
the release of `M2d`). Author: Claude Sonnet 5.5 on Claude Code. Evidence:
[`runs/T-5xjc/`](../../runs/T-5xjc/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #283. The plan review
(GPT-6.1 Sol on Codex CLI, `high`, 5 min 1 s) gave `approve-with-conditions`: 700 lines over
12 files against `0a104b8`, close-out inside; cap 2 (a check script is a gate); no panel. Its four
conditions are applied: the mutation plan separates the two causes of failure (the script as it
stands fails the clean release too, so a control with the one allowance of `AddRemote` shows which
mutant is red); the brief of the release review carries no plan, no plan review, no answer and no
earlier verdict; the record has a Coverage table that is checked against the commit; and the limits of
the check are stated as they are (check (3) fails too, and check (7) is function-level, so a new push
in `store.go` fails). The plan is reasoning-tier work that the session's model, an execution-tier one,
wrote under the Operator's routing message of 2026-10-13: a routing mismatch by choice. The goal count
is 2, final by O-209 of #267.

## What was done

1. **The release of `M2d`** (D1): the code of `main` at `0a104b8`, the merge of row 58: 81 non-test Go
   files of `cmd/` and `internal/` out of `testdata/` (19,854 lines), `go.mod` and the eight embedded
   files. The paragraph of the release in the plan names the commit.
2. **The check** (D3): `runs/T-efmy/release-check.sh` failed at the release commit on one line, check
   (4), `remote in AddRemote` (`git remote add`). Its check (4) now allows that, and a new check (7) lists
   the 14 references to the calls of `internal/git` that reach a remote and fails on a reference to
   `git.Push` out of `PushHead` and `commitFiles` of `internal/run/store.go` and `pushProbe` of
   `internal/run/forgecheck.go`, and on an import of `internal/git` under another name. Six mutants, each
   its own commit of a scratch clone: against the script as it stands and a control that has the one
   allowance, four pass the control (red) and two fail it on check (4); the final script fails all six, each on
   its own line. Output at the release: [`release-check.txt`](../../runs/T-5xjc/release-check.txt).
3. **The release review** (D4, the uat): [`release-review.md`](../../runs/T-5xjc/release-review.md), posted on
   #283: GPT-6.1 Sol on Codex CLI (`xhigh`, 16 min 35 s, 443,462 tokens), a model that wrote none of the
   release, in a clone with `runs/`, `docs/tasks/`, `docs/prd/` and `docs/plan/` out of its working tree,
   read the whole release (90 paths, each `in full`: the Coverage table matches the commit) and recorded
   `holds` for REQ-015 and REQ-017, with 29 findings and none a defect. Ten of them are notes closed with "No
   change": the stricter reading that counts every write to hosted state (token, issues, comments, the three
   pushes, the probe), the `gh api` commands that LAYUP writes for the Operator, the reach of a harness
   session (L-A1), the baseline scripts, and the limits of the check.
4. **The documents** (D5): the plan, the cells and a §13 line of `PRD-0001`, a uat row of
   `traceability.md`, the paragraph "Decided here" of `scaffold.md`, and two pitfalls of `guardrails.md` §2.

**Tests:** [`test-runs.md`](../../runs/T-5xjc/test-runs.md).

**Known limits:** a second call of `git.Push` inside `PushHead`, `commitFiles` or `pushProbe` passes check
(7), so the review read those three functions in full; check (7) reads the selector of the package, and a push
by a program that is started with the text `git push` is check (3)'s. The record is one session of one model:
its reading (a ref written to the target's own repository changes the repository and not the platform) is the
reading of `T-efmy` reached by its own reasons, and a stricter reading would make REQ-017 fail on the token,
the issues, the comments and each push, the probe included (findings 19 and 20). Two sentences of the record
count one too few references (it names each of them). The probe's real refusal by GitHub (`GH013` and the
ruleset's name) is still unverified: the demo of row 59b reads it.

## Review rounds

Round 1 (`62fd959`, cycle 0 of the cap of 2; comment 6166857401): Claude Opus 5.5 on Claude Code (`max`; a model that
differs from both authors of the change, Claude Sonnet 5.5 and GPT-6.1 Sol for the record, and the first harness of
ADR-0026's order, so no exception) gave `nothing material in scope`, with nine notes. It ran the release check three
times (each `== PASS`, 116 lines), read about seventy citations of the record and the 13 of the script, and checked the
Coverage table. Notes 1 to 5 asked for words and are applied in the header and two strings of the script (its logic
and its 197 lines unchanged; the six mutants and the release run again), the paragraph of `scaffold.md` and a PRD cell:
check (7) holds the push to three function names, so a new caller of one of them passes, and it reads a line at a time;
the import guard passes when it finds no import; check (4) admits any subcommand of `git remote` in `CloneLocal` and
`AddRemote`; and the calls that (7) lists name a remote verb, as `AddRemote` and `CloneLocal` reach no remote. Notes 6
and 7 (the chunk size of the brief, 250 against the answer's 200; the earlier reviews of the same model) are in the
test record, and note 8 (the sparse checkout keeps earlier verdicts out of the working tree only) is in the pitfall.
Note 9, revealed (the table of the acceptance tests of `M2d` has no row for the review of the release), is #322.

## Verdict

Delivered: the review of the release of `M2d`: `holds` for REQ-015 and REQ-017 over the code of `main` at `0a104b8`
(phase 1, `M2a`, `M2b`, `M2c` and `M2d`), by a reviewer of a model that wrote none of it, with the release check
changed for `M2d` (the allowance of `AddRemote`, and check (7) on the push of the probe). The review ended at
`nothing material in scope` at cycle 0 of the cap of 2. The branch changes no Go file. The diff against `0a104b8` is
within the plan review's 700 lines over 12 files, close-out inside. Next: row 59b, `T-sdtb` (#284), the demo, whose
inputs the Operator gives.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-14, from the clock and the comments of #283; the tokens of GPT-6.1 Sol
are the "tokens used" of Codex CLI, those of Claude Opus 5.5 the `modelUsage` of Claude Code with the cache; the
author's are not reported. ADR-0026 does not name the harness Codex; the Operator allowed it on 2026-10-13.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning (a routing mismatch) | Claude Sonnet 5.5 | `max` | not reported | 09:35 to 09:40 |
| The plan review | reasoning | GPT-6.1 Sol, Codex CLI | `high` | 121,953 | 5 min 1 s, 09:40:28 |
| The answer, the mutants, the red slice, check (7), the green run, the brief | execution | Claude Sonnet 5.5 | `max` | not reported | 09:47 to 09:53 |
| The release review (the uat) | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 443,462 | 16 min 35 s, 09:53:22 |
| The documents, the evidence, the gate, the freeze | execution | Claude Sonnet 5.5 | `max` | not reported | 10:10 to 10:31 |
| Round 1 | reasoning | Claude Opus 5.5, Claude Code | `max` | 22,437,426 (USD 10.47) | 26 min 57 s, 10:33:27 |
| The notes and the close-out | execution | Claude Sonnet 5.5 | `max` | not reported | 11:01 to 11:05 |
