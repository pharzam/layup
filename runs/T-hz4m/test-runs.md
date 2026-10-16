# T-hz4m — the runs of the tests

UTC, 2026-10-16; times from the clock of the host. Author: Claude Sonnet 5.5 (`claude-sonnet-5-5`) on Claude Code. Plan reviewer: GPT-6 Sol on Devin CLI (comment 6209093081 of #355); answer: comment 6209117686.

## Step T, red (05:59 to 06:00Z)

- **Red 1, with no production edit** (05:59:34 to 05:59:36Z): `TestTheWorkflowOfTheGoEntry` failed on each of the five jobs ("the overlay step at -1 (0 times)") and on the counts of `GATE_BASE` (5, want 10) and `GATE_HEAD_REF` (0, want 5); `TestTheJobsReadTheBase` failed at once, "the job of the kind static has no step that takes the gate files from the base".
- **Red 2, with the step in the workflow and the script as it stood** (05:59:56 to 06:00:45Z), so that each case fails on its own assertion: the script read `--overlay` as a kind, and each case that must pass the step failed with "no row for the kind --overlay" (the six cases of the head's files, the directory path, the six refs of `batch/`); each refusal case failed with a step that did not name its reason (a link, `not a directory`, `in .git`, `is not a path`). Two cases were green by accident and are not red: `a base that is a tree` (the old check of the revisions says "not a commit"; the mutant that drops the check in the new mode catches it) and `a base with no script` (`git show` fails in the step; the mutant `set -u` catches it). The first run of red 2 stopped at a defect of the test (a removal of a file that the entry lacks), fixed; the cases `.git` and `not a path` were first green by accident (the head carried the bad manifest too) and were given a clean head manifest (guardrails §2, Testing pitfalls).

## Green

- 06:01:20 to 06:01:36Z: `TestTheJobsReadTheBase` ok (`sh` is `dash` on this host; `dash` and `bash --posix` run too). 06:01:40 to 06:02:17Z: `go build`, `go vet`, `go test ./internal/catalog/` and with `-tags=integration` ok.
- 06:04:31 to 06:05:16Z: `go test -count=1 ./...` exit 0. 06:05:21 to 06:07:18Z: `go test -count=1 -tags=integration` of `internal/catalog`, `internal/setup`, `internal/verify`, `internal/cli`, `internal/standin` and `cmd/layup`: ok. **The full integration level of the module is not run** (the root schedules it).

## Mutants (06:02:57 to 06:10:46Z; the script is a one-off, not kept)

22 mutants, each applied alone to a copy of the tree, each run against `TestTheJobsReadTheBase` and the two unit tests, each caught by a named case: no overlay of the manifest, of the script, of the `config` paths; a removed path left in the index, left in the tree; the exception for any ref, for `batch/` alone, for `batch/a/b`, for a ref that holds `batch/`; the exception that fails; the checks of `.git`, of a path, of a link of the head, of a file above, of a link or submodule at the base, of the base commit; write before check; the modes not kept; the step after the run line; the step without `set -eu`; the step that runs the head's script; the step without `GATE_HEAD_REF`. One mutant survived the first run (the check of a link at the base: the head inherited the link and its own check refused it); the case was changed so that the head has no link, and the mutant is caught (06:11Z).

## The merge of `main` (row 63b, `d24fae3`) and the full integration run

- 06:19Z: `origin/main` (`d24fae3`, PR #376) merged into the branch as `4e016b7`, no conflict. Shared files that both sides changed: `docs/spec/loop.md` (row 63b at the line 146, this task in "The gate jobs read the base"), `docs/prd/PRD-0001-layup.md` (REQ-005 against REQ-007), `docs/tests/traceability.md` (a row each): different regions, so O-219's limit on shared index lines was not reached and nothing needed a scoped round. Cheap checks on `4e016b7`: `adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests`, `setup-check`, `git diff --check origin/main`: exit 0 each; `go build`, `go vet`, `go test -count=1 ./...`: exit 0.
- Row 64's full integration run (log `c64/integ.log`) ended at 06:23:15Z with exit 0; no other full run was alive, so this one started after it.
- **The full integration level of this task, on `4e016b7ac4f455852718e28410599cee6423ea6c`** (go1.26.9 linux/amd64, `TMPDIR` on disk, `go test -count=1 -tags=integration -timeout 25m ./...`, one run, with a heartbeat each minute): **start 06:25:05Z, end 06:37:46Z (12 min 41 s), exit 0.** 21 packages ok, no `FAIL`: `cmd/layup` 4.3 s, `internal/catalog` 37.9 s, `cli` 18.9 s, `forge` 0.3 s, `forge/github` 0.9 s, `gate` 0.9 s, `git` 8.9 s, `handoff` 0.04 s, `ledger` 0.02 s, `psb` 0.1 s, `records` 0.07 s, `route` 0.03 s, `rules` 0.03 s, `run` 697.7 s, `session` 8.1 s, `setup` 5.8 s, `spec` 0.03 s, `standin` 1.0 s, `tsv` 0.02 s, `verify` 57.6 s, `work` 0.06 s. `TestProcessTimesTheFirstOutput` (#291) did not flake. The log is `~/.cache/layup-session/scratch-hz4m/integ.log` (host-local, not committed).

## The budget (a count of added plus removed lines of `git diff --numstat`, before the close-out)

| Against | Lines | Files |
| ------- | ----- | ----- |
| `d24fae3`, current `main` (this task's own change) | 490 | 11 |
| `2aafa64`, the base of the previous report | 827 | 19 |
| `6bc1051`, the base that the plan review names | 1,391 | 25 |

The plan review's maximum, 650 over 14 files, is read against `6bc1051`, and the branch now holds the changes of rows 62a and 63b through the merges of `main`; the question of the base is #355, comment 6209355485, and the Operator has not answered. The live O-218 is row 63b's only. Nothing here claims it resolved. The close-out would add about 55 lines over 3 files.

## Resources of the implementation invocation (root CLI result)

Author model `claude-sonnet-5-5`, effort high; duration 1,272,732 ms; total cost 2.0120004 USD; model usage: input 126, output 65,893 (of which thinking 30,338, a subset, not added), cache read 7,388,824, cache creation 153,484 tokens. Plan reviewer: GPT-6 Sol on Devin CLI (no figures reported by the harness). The Operator's later direction: Claude Sonnet 5.5 and Devin SWE-2 Max only, reviews included.
