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
