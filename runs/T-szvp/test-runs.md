# T-szvp test runs

## In plain terms

Stage A is red for the expected reasons: the numbered checker exception remains,
and the Restore step differs from its pre-exception form. No workflow or existing
documentation has changed. No green or hosted result is claimed.

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
