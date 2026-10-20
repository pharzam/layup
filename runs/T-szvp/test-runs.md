# T-szvp test runs

## In plain terms

Stage A recorded the expected red. Stage B removes the numbered checker
exception, and the identical check is green. All local checks below pass.
Stage C records the supplied successful hosted run and an exit-0 pre-freeze
check on the clean committed Stage B head. Independent review remains pending.

## Stage A — removal check, red (R8)

Binding inputs: the supplied plan and approve-with-conditions plan review for
issue #452. Condition 1 keeps the complete Restore-step assertion solely in
`runs/T-szvp/checks.sh`; T-c826's exception assertions will be retired in a later
stage without adding a second restore-body pin. The coordinator posts Plan
amendment 1 before committing this stage.

The worktree began clean at `f9d46b0f441202254e8ea6ca16fe4ad8ec860858`.
The check is the plan's proposed script, unchanged. Its deterministic reference
is `.github/workflows/ci.yml` at pre-exception commit
`d27b33180a1aa82c23f3ab66e757e06de356eb35`.

Times below are actual UTC values from `date -u +%FT%TZ`, captured immediately
before and after each validation command. Standard output and standard error of
the red run were captured together under
`/home/layup/.cache/layup-session/tmp-green/T-szvp-stage-a-red.out`.

| Command | UTC start | UTC end | Exit | Result |
| --- | --- | --- | --- | --- |
| `sh runs/T-szvp/checks.sh` | 2026-10-20T23:39:59Z | 2026-10-20T23:39:59Z | 1 | Expected red: both named failures, with the existing exception lines shown below. |
| `git diff --exit-code HEAD -- .github/workflows/ci.yml docs runs/T-c826` | 2026-10-20T23:40:20Z | 2026-10-20T23:40:20Z | 0 | No output; workflow, documentation and prior task records remain unchanged. |
| `git diff --check` | 2026-10-20T23:40:20Z | 2026-10-20T23:40:20Z | 0 | No output; no tracked diff. The two new Stage A files are untracked pending the coordinator's commit. |

Complete red output (line numbers are relative to the extracted job):

```text
19:        # R4: T-c826 (#449), PR 453; remove in T-szvp (#452).
23:          EVENT_NAME: ${{ github.event_name }}
24:          PR_NUMBER: ${{ github.event.pull_request.number }}
37:            if [ "$EVENT_NAME" = pull_request ] && [ "$PR_NUMBER" = 453 ] && [ "$f" = docs/tests/freeze-check.sh ]; then
38:              printf 'kept as written for PR %s: %s\n' "$PR_NUMBER" "$f"
66:      # commits, which no hook step makes. Both runners always restore; PR 453
FAIL: numbered checker exception remains
FAIL: restore step differs from pre-exception baseline
```

The failure is not a missing job, missing reference or shell error: execution
reached both diagnostic assertions and returned their accumulated failure status.
The identical command must exit 0 after the later workflow removal.

## Stage B — removal, current documentation and green

Started from committed Stage A, `66fa80f`. Plan amendment 1 is supplied as posted
comment 6304288973. The workflow removes the authorization comments, event/number
environment bindings and checker-only keep/continue branch. All nine restoration
paths, default-branch source, in-place copies, missing-file fallback, both suite
invocations and the separate nested-checkout-driver exception are preserved.
The CI README, guardrails and test README now describe restoration for every
event. The migration lesson remains in the guardrails.

Condition 1 removes only T-c826's PR interface and exception assertions, with no
replacement body pin. Its six remaining assertions pass. T-szvp's check is
unchanged from Stage A and holds the sole complete Restore-step comparison.
`conditions.md` records the one condition.

Times are actual UTC values from `date -u +%FT%TZ` immediately before and after
each command. Standard output and error were captured together under
`/home/layup/.cache/layup-session/tmp-green`, in `T-szvp-stage-b-NN.out` (01–13,
in table order). `TMPDIR` was set to that directory for these runs, including the
disposable git repositories made by the fixture suite. The runner printed each
command and a still-running message at ten seconds when needed.

| Command | UTC start | UTC end | Exit | Result |
| --- | --- | --- | --- | --- |
| `sh runs/T-szvp/checks.sh` | 2026-10-20T23:44:56Z | 2026-10-20T23:44:56Z | 0 | Green, no output; identical check to Stage A. |
| `sh runs/T-c826/checks.sh` | 2026-10-20T23:44:56Z | 2026-10-20T23:44:56Z | 0 | All six remaining assertions pass without a PR argument. |
| `sh docs/adr/adr-lint.sh` | 2026-10-20T23:44:56Z | 2026-10-20T23:44:57Z | 0 | adr-lint: OK. |
| `sh docs/prd/prd-lint.sh` | 2026-10-20T23:44:57Z | 2026-10-20T23:44:57Z | 0 | prd-lint: OK. |
| `sh docs/tests/run-discipline-tests.sh` | 2026-10-20T23:44:57Z | 2026-10-20T23:45:00Z | 0 | 138 passed, 0 failed. |
| `sh docs/setup/setup-check.sh` | 2026-10-20T23:45:00Z | 2026-10-20T23:45:22Z | 0 | All checks OK, including ci; 138 discipline cases and 3,029 links pass. |
| `sh docs/tests/freeze-check/run-all.sh` | 2026-10-20T23:45:22Z | 2026-10-20T23:45:32Z | 0 | 85 passed, 0 failed. |
| `sh docs/links/link-lint.sh` | 2026-10-20T23:45:32Z | 2026-10-20T23:45:47Z | 0 | OK, 3,029 links resolved. |
| `git diff --check` | 2026-10-20T23:45:47Z | 2026-10-20T23:45:47Z | 0 | No whitespace errors. |
| `git diff d27b331 HEAD -- .github/workflows/ci.yml` | 2026-10-20T23:45:47Z | 2026-10-20T23:45:47Z | 0 | HEAD is still Stage A; four exception hunks remain there. Not a post-commit removal assertion. |
| `git diff d27b331 -- .github/workflows/ci.yml` | 2026-10-20T23:45:47Z | 2026-10-20T23:45:47Z | 0 | Working tree has only the planned suite-comment hunk, shown below. |
| `git diff --exit-code f9d46b0f441202254e8ea6ca16fe4ad8ec860858 -- docs/tests/freeze-check.sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/run-all.sh` | 2026-10-20T23:45:47Z | 2026-10-20T23:45:47Z | 0 | No output; checker and both runners unchanged against f9d46b0. |
| `git diff --exit-code HEAD -- runs/T-szvp/checks.sh` | 2026-10-20T23:45:47Z | 2026-10-20T23:45:47Z | 0 | No output; Stage A check unchanged. |

The HEAD comparison reads the committed Stage A workflow, because this session
cannot commit. Its exit 0 is a display-command status, not proof of an empty diff.
The working-tree comparison reads the actual removal and has exactly this
remaining hunk against `d27b331`:

```diff
diff --git c/.github/workflows/ci.yml w/.github/workflows/ci.yml
index 37360a5..f58a0d9 100644
--- c/.github/workflows/ci.yml
+++ w/.github/workflows/ci.yml
@@ -256,8 +256,8 @@ jobs:
       - run: sh docs/tests/nested-checkout-check.sh
       # T-b6fn (#440) — runs the fixture suite of freeze-check.sh. Each case
       # builds a disposable git repository in mktemp -d with `git init` and two
-      # commits, which no hook step makes. All three scripts are in the restore
-      # loop above.
+      # commits, which no hook step makes. The checker and both runners restore
+      # from the default branch for every event.
       - run: sh docs/tests/freeze-check/run-all.sh

   setup-check:
```

After the coordinator's `ci:` commit, rerun the planned HEAD comparison to confirm
that same single hunk at the committed endpoint.

Task-check output (T-szvp is silent on success; T-c826 prints six assertions):

```text
ok: canonical margin sentence in the tolerance section
ok: unchanged runner: docs/tests/freeze-check/run.sh
ok: unchanged runner: docs/tests/freeze-check/run-all.sh
ok: all nine baseline paths stay in the restore list
ok: restoration still reads the default branch of this repository
ok: default branch still comes from the repository event
```

Complete setup-check output:

```text
setup-check: pin OK
setup-check: kit-history OK
setup-check: facts OK
setup-check: onboarding OK
setup-check: glossary OK
setup-check: guardrails OK
setup-check: markers OK
setup-check: adapted OK
setup-check: ci OK
setup-check: protection OK
setup-check: identity OK
setup-check: procedure OK

run-discipline-tests: 138 passed, 0 failed
adr-lint: OK
prd-lint: OK
link-lint: OK  3029 links resolved
setup-check: kit-linters OK
```

Complete freeze-check fixture output:

```text
ok    bad-K1-blank-eof
ok    bad-K1-trailing-space
ok    bad-K2-missing-file
ok    bad-K2-missing-prd
ok    bad-K2-missing-traceability
ok    bad-K2-no-prd-frame-test
ok    bad-K2-no-prd-mapped-test
ok    bad-K2-no-s12-frame-test
ok    bad-K2-no-s12-mapped-test
ok    bad-K2-prd-other-file
ok    bad-K2-two-missing
ok    bad-K3-backtick-outside
ok    bad-K3-double-quoted
ok    bad-K3-no-match
ok    bad-K3-other-package
ok    bad-K3-other-tag
ok    bad-K3-span-across-lines
ok    bad-K3-subtest
ok    bad-K3-unclosed-span-eof
ok    bad-K3-unclosed-span-paragraph
ok    bad-K3-unquoted
ok    bad-K3-unsupported-form
ok    bad-K4-dropped-in-description
ok    bad-K4-form-a-different
ok    bad-K4-form-b-equal-hashes
ok    bad-K4-result-only
ok    bad-K4-two-hashes-no-restore
ok    bad-K5-over-files
ok    bad-K5-over-lines
ok    bad-K6-empty-table
ok    bad-K6-line-past-end
ok    bad-K6-line-zero
ok    bad-K6-missing-commit
ok    bad-K6-missing-file
ok    bad-K6-missing-path
ok    bad-K6-missing-test
ok    bad-K6-mixed-valid-invalid
ok    bad-K6-no-span
ok    bad-K6-range-zero
ok    bad-base-not-a-commit
ok    bad-crlf-conditions
ok    bad-env-K2-unreadable-prd
ok    bad-env-K2-unreadable
ok    bad-env-K3-unreadable-test
ok    bad-env-dirty-worktree
ok    bad-env-untracked-evidence
ok    bad-not-a-worktree
ok    bad-usage-bad-budget
ok    bad-usage-closeout
ok    bad-usage-evidence-directory-base
ok    bad-usage-evidence-directory-head
ok    bad-usage-evidence-missing
ok    bad-usage-evidence-no-value
ok    bad-usage-evidence-other-task
ok    bad-usage-evidence-suffix
ok    bad-usage-evidence-task-prefix
ok    bad-usage-evidence-trailing-slash
ok    bad-usage-evidence-traversal
ok    bad-usage-missing-budget
ok    bad-usage-no-args
ok    bad-usage-unknown-flag
ok    good-K1-clean
ok    good-K2-prd-by-file
ok    good-K2-registered
ok    good-K3-alternatives
ok    good-K3-span-across-lines
ok    good-K4-forms
ok    good-K5-at-max
ok    good-K5-ceiling
ok    good-K5-evidence-added
ok    good-K5-evidence-deleted
ok    good-K5-evidence-duplicate
ok    good-K5-evidence-inputs
ok    good-K5-evidence-kinds
ok    good-K5-evidence-literal
ok    good-K5-evidence-modified
ok    good-K5-evidence-renamed
ok    good-K5-file-limit
ok    good-K5-line-limit
ok    good-K5-no-closeout
ok    good-K5-percentage
ok    good-K5-truncation
ok    good-K6-no-conditions
ok    good-K6-valid
ok    good-crlf-all
freeze-check: 85 passed, 0 failed
```

## Stage B — operative-reference and history audit

The first command verifies that the six named T-c826 historical records remain
unchanged. Search output is scratch-only in `T-szvp-stage-b-audit-03.out` in the
same scratch directory. The repository-wide search includes hidden files and
excludes `.git`; `-i` also finds lowercase event identifiers. An initial search
without a path returned no matches when run from a captured subprocess; it read
stdin rather than the repository. That result is not used as inventory evidence.
The corrected command explicitly names `.`:

```sh
rg -n -i '(^|[^0-9])453([^0-9]|$)|PR_NUMBER|EVENT_NAME|kept as written for PR|T-szvp|checker.only.exception|temporary.exception' --hidden -g '!.git' .
```

| Command | UTC start | UTC end | Exit | Result |
| --- | --- | --- | --- | --- |
| `git diff --exit-code HEAD -- runs/T-c826/ci-exception.md runs/T-c826/conditions.md runs/T-c826/test-runs.md runs/T-c826/mutants.md docs/tasks/T-c826.md docs/tasks/completed.md` | 2026-10-20T23:46:21Z | 2026-10-20T23:46:21Z | 0 | No output; historical records preserved. |
| Reference search above, initially without final `.` | 2026-10-20T23:46:21Z | 2026-10-20T23:46:21Z | 1 | No matches; unusable inventory because stdin was searched. |
| Reference search above, with final `.` | 2026-10-20T23:46:53Z | 2026-10-20T23:46:53Z | 0 | 112 matching lines; classified below. No operative claim that the exception is active. |

| Remaining matches | Classification |
| --- | --- |
| `runs/T-c826/ci-exception.md`, `conditions.md`, `test-runs.md`, `mutants.md`, `docs/tasks/T-c826.md`, T-c826's completed-log line | Historical approvals, stage observations, commands, output and removal checklist for #449. Preserved as directed. |
| `runs/T-szvp/checks.sh`, `conditions.md`, `test-runs.md` | This task's absence assertion, condition evidence and historical red/current green records. |
| `docs/tests/README.md` | Current commands and removal assertions; states no exception remains. |
| `.github/workflows/ci.yml:31` | Existing lowercase `github.event_name` in concurrency cancellation; unrelated to checker restoration. |
| Other task records under `runs/` | Unrelated numeric hits in hashes, source line references, timings, telemetry and prior evaluation data. |

This stage adds no task detail or completed-log entry. The coordinator makes the
two commits named in the returned commit plan. Pre-freeze validation, current
`origin/main` resolution, hosted restore logs from the removal's own PR, independent
review and close-out remain later stages. Local green does not satisfy hosted
acceptance; no `ci-restore.md` is fabricated here.

Final whitespace check after writing the evidence and condition records:

| Command | UTC start | UTC end | Exit | Result |
| --- | --- | --- | --- | --- |
| `git diff --check` | 2026-10-20T23:47:15Z | 2026-10-20T23:47:15Z | 0 | No output; evidence context lines have trailing whitespace removed. |

## Stage C — hosted evidence and clean-head pre-freeze check

The worktree began clean at `8903e77b298059a7bdc159d36afbc0fbbfbdcd47`,
after the coordinator committed and pushed Stage B (`68931a1`, `8903e77`).
`origin/main` resolved locally to
`f9d46b0f441202254e8ea6ca16fe4ad8ec860858`. The clean-head and reference
inspection ran from 2026-10-20T23:57:29Z to 2026-10-20T23:57:29Z, exit 0:
`git rev-parse HEAD origin/main` printed those two full hashes, and
`git status --porcelain --untracked-files=all` printed nothing.

### Hosted result supplied by the coordinator

The binding Stage C input reports that [PR #456](https://github.com/pharzam/layup/pull/456)'s
[CI run 39517390208](https://github.com/pharzam/layup/actions/runs/39517390208)
succeeded. Its `nested-checkout-check` log contains:

```text
restored from main: docs/tests/freeze-check.sh
restored from main: docs/tests/freeze-check/run.sh
restored from main: docs/tests/freeze-check/run-all.sh
freeze-check: 85 passed, 0 failed
```

The coordinator also reports no "kept as written" line. These are supplied
hosted observations, not a local fixture run or a fresh log retrieval by this
session. Hosted start/end times, the tested merge SHA, job URL and the other six
restore lines were not supplied; no values for them are inferred.

An attempted read-only retrieval first used an unsupported JSON field, then
failed to connect after that field was corrected. Both failures are recorded:

| Command | UTC start | UTC end | Exit | Result |
| --- | --- | --- | --- | --- |
| `gh run view 39517390208 --repo pharzam/layup --json url,headSha,event,status,conclusion,createdAt,startedTime,updatedAt,jobs` | 2026-10-20T23:57:42Z | 2026-10-20T23:57:43Z | 1 | Unknown JSON field: startedTime; no hosted data read. |
| `gh run view 39517390208 --repo pharzam/layup --json url,headSha,event,status,conclusion,createdAt,startedAt,updatedAt,jobs` | 2026-10-20T23:57:46Z | 2026-10-20T23:57:46Z | 1 | Error connecting to api.github.com; no hosted data read. |

### Pre-freeze check

Before editing this record, ran the exact Stage C invocation on the clean head:

```sh
sh docs/tests/freeze-check.sh --task T-szvp --base origin/main --budget 250/8 --evidence runs/T-szvp/conditions.md --evidence runs/T-szvp/test-runs.md > /home/layup/.cache/layup-session/tmp-green/szvp-fc.txt 2>&1
```

`TMPDIR` was set to `/home/layup/.cache/layup-session/tmp-green`, so the check's
temporary files also stayed outside the worktree. Start and end times were
captured with `date -u +%FT%TZ` immediately around the command.

| Command | UTC start | UTC end | Exit | Result |
| --- | --- | --- | --- | --- |
| Exact pre-freeze invocation above | 2026-10-20T23:57:35Z | 2026-10-20T23:57:35Z | 0 | K1–K6 pass; complete output below. |
| `git diff d27b331 HEAD -- .github/workflows/ci.yml` | 2026-10-20T23:58:00Z | 2026-10-20T23:58:00Z | 0 | Committed Stage B endpoint has only the planned suite-comment hunk already shown in Stage B. Restore step matches the pre-exception baseline. |
| `git diff --check` | 2026-10-20T23:58:38Z | 2026-10-20T23:58:38Z | 0 | No output; Stage C evidence has no whitespace errors. Status shows only this evidence file modified. |

Complete pre-freeze output:

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
ok K5: projection: figure 124/6, maximum 250/8, tolerance 60/4, limit 310/12, margin 186/6 (the figure at landing is still measured at landing)
ok K6: conditions table
```

The two declared evidence records are excluded under ADR-0028; the executable
task checks remain counted. The projection is for this committed head, before
close-out. This stage changes only this evidence record, adds no task detail or
completed-log entry, and names no frozen SHA. After the coordinator commits the
record, the resulting clean head still needs its pre-freeze rerun before freeze
and independent review, within cycle cap 2.
