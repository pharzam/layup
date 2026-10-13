# T-jpkc — test runs

The forge calls of `M2d` (#271, row 51a): `EffectiveRules` and `Activity` of the forge interface
and of the GitHub adapter, and the head of each branch in `Repository`. All times are UTC,
2026-10-13, read from the clock after each step.

## The order of the gate

The plan was posted at 09:11:53Z, its review (GPT-6.1 Sol on Codex CLI, effort `high`) ran from
09:11:59Z to 09:15:57Z and was posted at 09:19:19Z, and the author's answer at 09:19:37Z. The tests
were written after that, and run red at 09:24:39Z.

## The base (10:19:16Z to 10:34:29Z)

At `9144e4e`, each with exit 0 (`go test` of `internal/forge/...`, then with `-tags=integration`):

```
go build ./...: exit 0 (09:11:41Z)
go vet ./...: exit 0 (09:11:44Z)
go test ./internal/forge/...: exit 0 (09:11:46Z)
go test -tags=integration ./internal/forge/...: exit 0 (09:11:48Z)
```

## The source of the calls

Read on 2026-10-13, as `forge.md` asks ("read in GitHub's documentation at the date of the build
task"), after the plan review found that the plan's sentence on the permission was wrong:

- The pages: `https://docs.github.com/en/rest/repos/rules?apiVersion=2022-11-28#get-rules-for-a-branch`
  ("Fine-grained access tokens for "Get rules for a branch"": the permission set `"Metadata"
  repository permissions (read)`) and
  `https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-repository-activities`
  ("Fine-grained access tokens for "List repository activities"": the permission set `"Contents"
  repository permissions (read)`). `runs/T-hbw8/slice-reviews/slice-A.md:257` says the same.
- The REST description: `github/rest-api-description`, `descriptions/api.github.com/api.github.com.json`,
  commit `9492527246942bb85a2e05ad167881076c27835a` of 2026-10-13T02:18:32Z ("Update OpenAPI 3.0
  Descriptions"). `GET /repos/{owner}/{repo}/rules/branches/{branch}` (`repos/get-branch-rules`) takes
  `per_page` and `page`; its `branch` is "The name of the branch. Cannot contain wildcard
  characters", with `x-multi-segment: true`; it gives an array of rules, each with `type`, an
  optional `parameters`, and the fields of its ruleset `ruleset_source_type`, `ruleset_source` and
  `ruleset_id`. `GET /repos/{owner}/{repo}/activity` (`repos/list-activities`) takes `ref`
  ("The `ref` for a branch can be formatted either as `refs/heads/BRANCH_NAME` or `BRANCH_NAME`"),
  `per_page` (default 30, at most 100) and others; each item has `id`, `node_id`, `before`, `after`,
  `ref`, `timestamp`, `activity_type` (`push`, `force_push`, `branch_deletion`, `branch_creation`,
  `pr_merge`, `merge_queue_merge`) and `actor`, a simple user that may be null.

## Red (09:24:39Z to 09:24:41Z)

The tests before the code: `go vet` of `internal/forge/...` with and without the tag `integration`
stops at the first error, and `go test` lists the others, as the packages do not build.

```
== go vet ./internal/forge/...
# github.com/pharzam/layup/internal/forge
# [github.com/pharzam/layup/internal/forge]
vet: internal/forge/forge_test.go:25:59: undefined: Rule
== go vet -tags=integration ./internal/forge/...
# github.com/pharzam/layup/internal/forge
# [github.com/pharzam/layup/internal/forge]
vet: internal/forge/forge_test.go:25:59: undefined: Rule
# github.com/pharzam/layup/internal/forge/github
# [github.com/pharzam/layup/internal/forge/github]
vet: internal/forge/github/github_integration_test.go:267:151: unknown field Heads in struct literal of type forge.Repository
== go test -tags=integration ./internal/forge/github/ (the compile errors)
# github.com/pharzam/layup/internal/forge/github [github.com/pharzam/layup/internal/forge/github.test]
internal/forge/github/github_integration_test.go:267:151: unknown field Heads in struct literal of type forge.Repository
internal/forge/github/github_integration_test.go:544:16: a.EffectiveRules undefined (type *Adapter has no field or method EffectiveRules)
internal/forge/github/github_integration_test.go:545:18: undefined: forge.Rule
internal/forge/github/github_integration_test.go:564:16: a.EffectiveRules undefined (type *Adapter has no field or method EffectiveRules)
internal/forge/github/github_integration_test.go:574:17: a.EffectiveRules undefined (type *Adapter has no field or method EffectiveRules)
internal/forge/github/github_integration_test.go:588:19: a.EffectiveRules undefined (type *Adapter has no field or method EffectiveRules)
internal/forge/github/github_integration_test.go:591:19: a.Activity undefined (type *Adapter has no field or method Activity)
internal/forge/github/github_integration_test.go:602:49: a.EffectiveRules undefined (type *Adapter has no field or method EffectiveRules)
internal/forge/github/github_integration_test.go:603:49: a.Activity undefined (type *Adapter has no field or method Activity)
internal/forge/github/github_integration_test.go:618:14: a.EffectiveRules undefined (type *Adapter has no field or method EffectiveRules)
internal/forge/github/github_integration_test.go:618:14: too many errors
FAIL	github.com/pharzam/layup/internal/forge/github [build failed]
FAIL
== go vet ./internal/run/ (the stand-in)
# github.com/pharzam/layup/internal/run
# [github.com/pharzam/layup/internal/run]
vet: internal/run/steps_test.go:60:73: undefined: forge.Rule
```

## Mutations

Each rule of the three calls mutated alone, on a copy of the tree and never on the worktree:
a path, a query, a page size, a tag, a status, a name of a call, a condition or a field changed,
so each mutant builds; then `go test -tags=integration ./internal/forge/...`. A mutant that does
not build would be counted as not run, never as caught.

- **The first run** (09:31:37Z to 09:36:06Z): 35 mutants, 33 caught, **1 survived** and **1 did not
  build**. `rules-no-escape` survived because the branch of the test, `feature/x y`, is escaped by
  the HTTP client itself whatever the adapter does; the test now uses `feature/x#y z`, whose `#`
  the client would read as a fragment. `rules-page-loop` left a variable unused and did not
  build; its text now keeps it.
- **The final run** (09:36:23Z to 09:37:55Z): 35 of 35 caught, each with its first failing line
  (each line cut at 170 characters).

```
== rules-path: caught github_integration_test.go:553: EffectiveRules over two pages = [], EffectiveRules: status 404: Not Found
== rules-per-page: caught github_integration_test.go:553: EffectiveRules over two pages = [], EffectiveRules: status 400: per_page is not 100
== rules-no-escape: caught github_integration_test.go:575: EffectiveRules: status 400: per_page is not 100
== rules-join: caught github_integration_test.go:578: the path of the request is ["/repos/acme/target/rules/branches/feature%2Fx%23y%20z"]; want ["/repos/acme/target/rules/branches/feature/x%2
== rules-call-name: caught github_integration_test.go:607: EffectiveRules with the body "{\"message\":\"x\"}": EffectiveRule: status 200: the body cannot be read: json: cannot unmarshal object into
== rules-status: caught github_integration_test.go:553: EffectiveRules over two pages = [], EffectiveRules: status 200: [{"ruleset_id":11,"ruleset_source":"acme/target","ruleset_source_type":"Re
== rules-type-tag: caught github_integration_test.go:553: EffectiveRules over two pages = [{Type: Parameters:map[]} {Type: Parameters:map[allowed_merge_methods:[merge] required_approving_review_co
== rules-parameters-tag: caught github_integration_test.go:553: EffectiveRules over two pages = [{Type:creation Parameters:map[]} {Type:pull_request Parameters:map[]} {Type:deletion Parameters:map[]} {T
== rules-empty-parameters: caught github_integration_test.go:553: EffectiveRules over two pages = [{Type:creation Parameters:map[]} {Type:pull_request Parameters:map[allowed_merge_methods:[merge] required
== rules-type-kept: caught github_integration_test.go:553: EffectiveRules over two pages = [{Type:x Parameters:map[]} {Type:x Parameters:map[allowed_merge_methods:[merge] required_approving_review_
== rules-page-loop: caught github_integration_test.go:553: EffectiveRules over two pages = [{Type:creation Parameters:map[]} {Type:pull_request Parameters:map[allowed_merge_methods:[merge] required
== rules-host-check: caught github_integration_test.go:620: EffectiveRules with a next page elsewhere: EffectiveRules: Get "http://elsewhere.invalid/repos/acme/target/rules/branches/main?page=2": di
== rules-host-error-call: caught github_integration_test.go:620: EffectiveRules with a next page elsewhere: Repository: status 200: the next page is not at the api of the forge register; want a forge.Err
== rules-host-error-status: caught github_integration_test.go:620: EffectiveRules with a next page elsewhere: EffectiveRules: the next page is not at the api of the forge register; want a forge.Error of Ef
== activity-path: caught github_integration_test.go:592: Activity with a body of null = [], Activity: status 404: Not Found; want none
== activity-per-page: caught github_integration_test.go:646: the queries are ["per_page=100&ref=main"]; want one request, with the ref and no other parameter
== activity-no-escape: caught github_integration_test.go:654: Activity: status 400: 400 Bad Request
== activity-call-name: caught github_integration_test.go:607: Activity with the body "{\"message\":\"x\"}": Activities: status 200: the body cannot be read: json: cannot unmarshal object into Go value
== activity-status: caught github_integration_test.go:592: Activity with a body of null = [], Activity: status 200: null; want none
== activity-actor-tag: caught github_integration_test.go:643: Activity = [{ActorID:0 ActorLogin:ghost Ref:refs/heads/main Before:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa After:bbbbbbbbbbbbbbbbbbbbbbbb
== activity-ref-tag: caught github_integration_test.go:643: Activity = [{ActorID:9001 ActorLogin:acme-layup[bot] Ref: Before:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa After:bbbbbbbbbbbbbbbbbbbbbbbbbb
== activity-before-tag: caught github_integration_test.go:643: Activity = [{ActorID:9001 ActorLogin:acme-layup[bot] Ref:refs/heads/main Before: After:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb Timestamp:
== activity-after-tag: caught github_integration_test.go:643: Activity = [{ActorID:9001 ActorLogin:acme-layup[bot] Ref:refs/heads/main Before:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa After: Timestamp:
== activity-timestamp-tag: caught github_integration_test.go:643: Activity = [{ActorID:9001 ActorLogin:acme-layup[bot] Ref:refs/heads/main Before:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa After:bbbbbbbbbbb
== activity-type-tag: caught github_integration_test.go:643: Activity = [{ActorID:9001 ActorLogin:acme-layup[bot] Ref:refs/heads/main Before:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa After:bbbbbbbbbbb
== activity-ghost: caught github_integration_test.go:643: Activity = [{ActorID:9001 ActorLogin:acme-layup[bot] Ref:refs/heads/main Before:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa After:bbbbbbbbbbb
== activity-actor-id: caught github_integration_test.go:643: Activity = [{ActorID:0 ActorLogin:acme-layup[bot] Ref:refs/heads/main Before:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa After:bbbbbbbbbbbbbb
== activity-actor-login: caught github_integration_test.go:643: Activity = [{ActorID:9001 ActorLogin:ghost Ref:refs/heads/main Before:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa After:bbbbbbbbbbbbbbbbbbbbb
== activity-type-kept: caught github_integration_test.go:643: Activity = [{ActorID:9001 ActorLogin:acme-layup[bot] Ref:refs/heads/main Before:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa After:bbbbbbbbbbb
== activity-before-kept: caught github_integration_test.go:643: Activity = [{ActorID:9001 ActorLogin:acme-layup[bot] Ref:refs/heads/main Before: After:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb Timestamp:
== heads-sha-tag: caught github_integration_test.go:268: Repository = {DefaultBranch:main Visibility:public HasCommit:true Branches:[main] Heads:map[main:]}, <nil>
== heads-sha-field: caught github_integration_test.go:268: Repository = {DefaultBranch:main Visibility:public HasCommit:true Branches:[main] Heads:map[main:]}, <nil>
== heads-value: caught github_integration_test.go:268: Repository = {DefaultBranch:main Visibility:public HasCommit:true Branches:[main] Heads:map[main:main]}, <nil>
== heads-returned: caught github_integration_test.go:268: Repository = {DefaultBranch:main Visibility:public HasCommit:true Branches:[main] Heads:map[]}, <nil>
== heads-names-kept: caught github_integration_test.go:268: Repository = {DefaultBranch:main Visibility:public HasCommit:false Branches:[] Heads:map[main:0000000000000000000000000000000000000001]},
35 caught, 0 survived, 0 not run, of 35
```

## Green (09:42:13Z to 09:52:56Z, on `db0ad97`)

Each with exit 0, in this order; the time of a line is when its command ended. The three runs of
`go test` are with `-count=1`, then with `-tags=integration` and with `-tags=e2e`, `TMPDIR` is on the
disk, and `git diff --check` reads the commits against `origin/main`:

```
go build ./...: exit 0 (09:42:15Z)
go vet ./...: exit 0 (09:42:15Z)
gofmt -l internal cmd (empty): exit 0 (09:42:16Z)
go test ./...: exit 0 (09:42:23Z)
go test -tags=integration ./...: exit 0 (09:50:39Z)
go test -tags=e2e ./...: exit 0 (09:52:15Z)
adr-lint: exit 0 (09:52:16Z)
prd-lint: exit 0 (09:52:16Z)
link-lint: exit 0 (09:52:27Z)
run-discipline-tests: exit 0 (09:52:33Z)
nested-checkout-check: exit 0 (09:52:34Z)
setup-check: exit 0 (09:52:56Z)
git diff --check origin/main HEAD: exit 0 (09:52:56Z)
```

## After round 1 (10:10:09Z to 11:29:04Z)

Round 1 (comment 6143289886 of #271) gave `material` on two findings. **The sentence on a branch that does not exist yet** was wrong: the endpoint answers by
the name of the branch, so `layup-probe`, before its probe, gets the rules of the rulesets that name
it, and a branch with no rule gets an empty list. `forge.md` and the comment of `EffectiveRules` say
so, and `TestEffectiveRulesOfABranchThatDoesNotExistYetGivesTheRulesThatNameIt` asks a branch that no
branch of the repository has, keeps each rule that the server gives and counts no request of the
list of branches; the test of a branch with no rule is its own. It passes at the code, which
forwards what the server gives, so its check is a mutant on a copy of the tree (10:10:09Z to
10:10:17Z): `EffectiveRules` that first asks `Repository` and returns nothing for a branch that is
not in the list. The new test fails with it; `TestEffectiveRulesGivesEachActiveRuleOfABranchOverEveryPage`
passes, as its branch is in the list; and `TestEffectiveRulesEscapesTheBranchSegmentBySegment` fails
too, as its branch is not in the list either:

```
EffectiveRulesOfABranchThatDoesNotExistYetGivesTheRulesThatNameIt: --- FAIL: TestEffectiveRulesOfABranchThatDoesNotExistYetGivesTheRulesThatNameIt (0.39s) FAIL
EffectiveRulesGivesEachActiveRule: ok  	github.com/pharzam/layup/internal/forge/github	0.228s
EffectiveRulesEscapesTheBranch: --- FAIL: TestEffectiveRulesEscapesTheBranchSegmentBySegment (0.17s) FAIL
```

**The base of the budget.** The plan review set 850 lines over 14 files against `9144e4e`. The
branch took `main` by a merge four times (`38ee78e`, the merge of #306; `6734a23`, the merge of
#307; `b28171b`, the merge of #311; and `5e40952`, the merge of #310), so the diff against `9144e4e` counts those tasks. Against the last of them the
change of this task is the figure of the close-out; the comments "The base of the budget" on #271
give the figures.

**The checks that read the tree.** A first run on `b001c03` (the head after the fix, the first merge of
`main` and the pitfall of the budget; 10:19:16Z to 10:34:29Z) had every line at exit 0 but the
integration level, which exited 1: the package `internal/run` ran past the 10 minutes of its timeout
(`FAIL github.com/pharzam/layup/internal/run 600.079s`) while the host also ran two review sessions
and the tests of another task (a load average of 8 to 10). The level run again alone with the
timeout of CI, `go test -count=1 -tags=integration -timeout 10m ./...` (10:34:44Z to 10:41:04Z), exited 0, and
`internal/run` took 355 s of the 600 that its timeout allows. CI's own `tests` job runs the three
levels in less than 3 minutes. Then the branch took `main` twice more, and the whole run was made
again after each. On `1f896ec` (10:54:09Z to 11:06:04Z, after the merge of `b28171b`) each line was
at exit 0, the integration level in 9 min 14 s of the 10 that its timeout allows, as the host ran a
review session and the tests of another task. On `6200944` (11:20:41Z to 11:29:04Z, after the merge of
`5e40952`), the merged head that round 2 reads with this record, each line was at exit 0, the integration level in 5 min 55 s. The
three runs of `go test` are with `-count=1`, then `-tags=integration` and `-tags=e2e`:

```
go build ./...: exit 0 (11:20:43Z)
go vet ./...: exit 0 (11:20:43Z)
gofmt -l internal cmd (empty): exit 0 (11:20:43Z)
go test ./...: exit 0 (11:20:49Z)
go test -tags=integration ./...: exit 0 (11:26:44Z)
go test -tags=e2e ./...: exit 0 (11:28:13Z)
adr-lint: exit 0 (11:28:15Z)
prd-lint: exit 0 (11:28:15Z)
link-lint: exit 0 (11:28:29Z)
run-discipline-tests: exit 0 (11:28:36Z)
nested-checkout-check: exit 0 (11:28:37Z)
setup-check: exit 0 (11:29:04Z)
git diff --check origin/main HEAD: exit 0 (11:29:04Z)
```
