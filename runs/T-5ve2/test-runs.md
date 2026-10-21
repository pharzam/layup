# T-5ve2 test runs

## Red (before any .gitattributes change; .gitattributes is main's at 01d7da30)

- Command: `GOROOT=$HOME/.local/share/mise/installs/go/1.26.10 PATH=$GOROOT/bin:$PATH go test -count=1 -tags=integration -run TestTheSharedIndexesMergeWithNoConflictInLandingOrder -v ./internal/git/`
- Start 2026-10-21T11:58:06Z, end 2026-10-21T11:58:07Z (`date -u`).
- Exit status: 1 (`FAIL`).
- First failing assertion: `indexmerge_integration_test.go:109: git merge main into B: exit status 1`, with `CONFLICT (content): Merge conflict in docs/tasks/completed.md` and in `docs/tests/traceability.md`. This is the right reason: no `merge=union` attribute.
- The negative subtest (`check-attr merge` is `unspecified` on backlog.md, the PRD, guardrails.md and blocks_integration_test.go; two edits of one PRD REQ-005 row conflict) passes at red, as the plan review said it may.
