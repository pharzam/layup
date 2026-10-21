# T-5ve2 test runs

## Red (before any .gitattributes change; .gitattributes is main's at 01d7da30)

- Command: `GOROOT=$HOME/.local/share/mise/installs/go/1.26.10 PATH=$GOROOT/bin:$PATH go test -count=1 -tags=integration -run 'TestTheSharedIndexesMergeWithNoConflictInLandingOrder' -v ./internal/git/`
- Start 2026-10-21T11:58:06Z, end 2026-10-21T11:58:07Z (`date -u`).
- Exit status: 1 (`FAIL`).
- First failing assertion: `indexmerge_integration_test.go:109: git merge main into B: exit status 1`, with `CONFLICT (content): Merge conflict in docs/tasks/completed.md` and in `docs/tests/traceability.md`. This is the right reason: no `merge=union` attribute.
- The negative subtest (`check-attr merge` is `unspecified` on backlog.md, the PRD, guardrails.md and blocks_integration_test.go; two edits of one PRD REQ-005 row conflict) passes at red, as the plan review said it may.

## Green (after the .gitattributes block: completed.md and traceability.md merge=union)

- Same command as the red. Read at 2026-10-21T11:58:56Z (`date -u`, taken after the run).
- Exit status 0: `--- PASS: TestTheSharedIndexesMergeWithNoConflictInLandingOrder (0.14s)`, both subtests pass.
- `sh docs/tests/run-discipline-tests.sh`: `138 passed, 0 failed`, so the eol rule (`run-discipline-tests.sh:172`, which reads only `eol=crlf` lines) is unaffected.
