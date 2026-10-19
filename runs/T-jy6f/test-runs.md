# T-jy6f - the test runs of row 66r

## S0 - the reconciliation of the plan (6250184301) with the landed code

Read at `origin/main` = `5a290c9` (the branch had no work, so the merge of `origin/main` was a fast-forward). Each name of the plan
was checked in the tree, not in a plan or a review (note 4 of 6250262119).

| Plan name | Landed form | Note |
| --------- | ----------- | ---- |
| 66f: `PostStatuses`, `HeadInput`, `PostStatus`, `Status`, `VerifyAt`, `VerifyInput`, `GatesTable`, `Sessions.Bot` | `internal/run/statuses.go:202`, `:176`, `:63`, `:33`; `verifystatus.go:33`, `:21`; `gates.go:167`; `session.go:65` | `PostStatuses` posts the four contexts, `layup/verify` included (a `not-active` failure too, `verifystatus.go:44-46`) |
| 66f: `holdsLease` | `handoffpost.go:418` | a method of `Sessions`; a `state` reaches it through `state.sessions()` |
| 66g: `OpenPullRequest`, `PullInput`, `Sessions.Web` | `pullrequest.go:104`, `:24`; `session.go:66` | the title is `Type` and `Description` (`checkTitle`), no task ID |
| 66q: the body edit | `(*Sessions).LinkReviewRecord(ctx, task, verifier) (bool, error)`, `reviewlink.go:30` | the plan had no name; it reads the latest `valid comment` handoff of the verifier |
| 66s: `PlanFile`, `PlanReviewFile`, `ReviewRecordFile`, `PromptIn.Options` | `instructions.go:22-24`; `promptfile.go:41-49` | landed |
| 66s: `PromptContext.base` (and `task`, `step`) | `promptfile.go:51-64`, set at `:176` | landed, so the guard needs no new field; `readOnlyStep` is not in `internal/session` (the plan cited `context.go:59-63`, which is the step table): it is to be decided from the private `PromptContext.step` |
| 66s: a reusable per-kind Definition of Done | none | 66s landed no exported text of it; the `shape-issue` text carries its own (two blocks, one for each kind) |
| 66v: `RepairTask` | `repair.go:30` | the driver calls it first |
| the base guard | none | `loop.md:469` leaves it to this row; `TaskSession` (`session.go:208`) has none |
| `probeTaskID` | `probe.go:162`, reads `sessions.tsv` only, `crypto/rand` direct | the seam `Sessions.random` and `drawTaskID` are to be made |
| `registerTask`, `registerFiles`, `parameterFiles`, `loopPermissionsStep` | `register.go:267`, `:213`, `:254`; `parameters.go` | as planned; `state.store` is the concrete `*recordsStore`, so the opening tests run in the real-git world (condition 5) |
| the pinned linter call | `spec.WaitingSpecAccepted(script, prdDir)`, `internal/spec/prdlint.go:61` | for the scope addition; the script reads the facts from `<dirname prdDir>/facts` |

### A landed form that the plan does not settle (STOP)

The verifier session of the plan starts at the bound head (`TaskSession{Base: bound}`, plan §2.5). `session.Make` clones **the
default branch only** (`git.CloneLocal`, `--single-branch`, `internal/git/git.go:507`; `session.md:52-58`) and then
`SwitchCreate(task/<task>/<attempt>, base)`. A developer's head is in the run's clone only as `refs/layup/sessions/<id>` (and, after
a restart, as `origin/task/<task>/<attempt>`), so it is not in the history of the default branch, and the start of the verifier
fails. The existing test world moves the default branch onto the head by hand (`reviewlink_integration_test.go:61`, "as the loop would
before a verifier"); no production code does. The close-out of row 66k's children needs the same.
