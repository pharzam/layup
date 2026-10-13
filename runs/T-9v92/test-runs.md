# T-9v92 — test runs

The rulesets of S13 (#274, row 52b): `rulesetBody` and `rulesetRecordsBody`, the step S13 in
phase 2 with the App's ID of `Phase2`, and the row `without-layup`, in `internal/setup`.
All times are UTC, 2026-10-13, read from the clock after each step.

## The order of the gate

The plan was posted at 08:05:06Z. The plan review ran from 08:05:09Z to 08:14:06Z, and was
posted at 08:14:30Z, with the author's answer at 08:14:47Z. The tests were written and run
red from 08:06Z to 08:07:44Z, in the time of the review, which R12 does not allow ("one round
of independent review and a reviewer's confirmation before building begins"): a breach of the
author's, found by round 1 (comment 6141901257 of #274). The confirmation asked for changes of
documents and none of a test, so no test was rewritten because of it. The code (the green
step) was written after the answer, from about 08:17Z, and so was each document. The red
below is the red of tests written before the confirmation, and it stays dated as it was; a
later passing run cannot move it. `guardrails.md` §2 gains the lesson.

## The base (08:03:33Z to 08:04:25Z)

At `9144e4e`, each with exit 0 (`go test` of `internal/setup`, `internal/verify` and
`internal/work`, then with `-tags=integration` of `internal/setup` and `internal/verify`):

```
go build ./...: exit 0 (08:03:39Z)
go vet ./...: exit 0 (08:03:42Z)
go test ./internal/setup/ ./internal/verify/ ./internal/work/: exit 0 (08:03:45Z)
go test -tags=integration ./internal/setup/ ./internal/verify/: exit 0 (08:04:25Z)
```

## Red (08:07:37Z to 08:07:44Z)

The tests before the code: `go vet` of `internal/setup` with and without the tag
`integration` stops at the first error, and `go test` lists the others, as the package
does not build.

```
# github.com/pharzam/layup/internal/setup
# [github.com/pharzam/layup/internal/setup]
vet: internal/setup/phase2_test.go:24:165: unknown field AppID in struct literal of type Phase2
---
# github.com/pharzam/layup/internal/setup
# [github.com/pharzam/layup/internal/setup]
vet: internal/setup/phase2_test.go:24:165: unknown field AppID in struct literal of type Phase2
# github.com/pharzam/layup/internal/setup [github.com/pharzam/layup/internal/setup.test]
internal/setup/phase2_test.go:24:165: unknown field AppID in struct literal of type Phase2
internal/setup/rulesets_test.go:134:56: too many arguments in call to rulesetBody
	have ([]string, number)
	want ([]string)
internal/setup/rulesets_test.go:142:14: undefined: rulesetRecordsBody
internal/setup/rulesets_test.go:182:58: too many arguments in call to rulesetBody
	have ([]string, int64)
	want ([]string)
internal/setup/rulesets_test.go:226:15: undefined: rulesetRecordsBody
internal/setup/rulesets_test.go:266:4: p.AppID undefined (type Phase2 has no field or method AppID)
FAIL	github.com/pharzam/layup/internal/setup [build failed]
FAIL
```

**The golden of phase 1 alone** (a compile failure shows that the tests need code, not that
the content of a file fails): on a copy of the tree with the code of the base (`final.go`,
`phase2.go`) and the changed `final_test.go`, `go test -run 'TestS13$' ./internal/setup/`
(08:40:49Z to 08:40:53Z) exits 1, and the first lines of its output, each cut at 150
characters, are:

```
--- FAIL: TestS13 (0.00s)
    final_test.go:208: S13: done "the ruleset file, and its commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72
    final_test.go:214: ruleset-default.json:
```

The file that the base writes has `pull_request` first, and no `creation`, no merge method and
no second parameter, which the new golden names.

Two cases guard what the change must not change, and have no red of their own:
`TestPhase2S13KeepsItsCommandsAndHandOff` and `TestS13OfPhase1WritesNoRecordsRulesetAndNoRow`
(phase 1 writes no records file and no row). They guard the body that two callers share, as
note 5 of the plan review asks; the mutants `step-s13-commands-dropped` and
`phase1-writes-records` below are their red.

## Mutations (08:19:04Z to 08:22:11Z)

Each rule of the two ruleset bodies, of `runS13With` and `runS13Phase2` and of the wiring in
`phase2.go` mutated alone, on a copy of the tree and never on the worktree: a condition
joined with `&& false`, or a name, a value, a key, an ID or a position changed, so each mutant
builds; then `go test ./internal/setup/` (the unit level). A mutant that does not build
would be counted as not run, never as caught. One run: 61 mutants, all caught, none survived,
none failed to build; each line gives the first failing line of the test (cut at 170
characters).

```
== layup-check-gates: caught rulesets_test.go:136: the file of the default branch (err <nil>):
== layup-check-spec: caught rulesets_test.go:136: the file of the default branch (err <nil>):
== layup-check-verify: caught rulesets_test.go:136: the file of the default branch (err <nil>):
== layup-check-rules: caught rulesets_test.go:136: the file of the default branch (err <nil>):
== layup-check-dropped: caught rulesets_test.go:136: the file of the default branch (err <nil>):
== layup-check-order: caught rulesets_test.go:136: the file of the default branch (err <nil>):
== phase2-no-layup-checks: caught rulesets_test.go:136: the file of the default branch (err <nil>):
== phase1-layup-checks: caught final_test.go:208: S13: done "the ruleset file, and its commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0ee3a309cbb06297
== layup-check-pin: caught rulesets_test.go:136: the file of the default branch (err <nil>):
== gate-job-pin: caught final_test.go:208: S13: done "the ruleset file, and its commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0ee3a309cbb06297
== rule-creation-dropped: caught final_test.go:208: S13: done "the ruleset file, and its commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0ee3a309cbb06297
== rule-order-creation-deletion: caught final_test.go:208: S13: done "the ruleset file, and its commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0ee3a309cbb06297
== rule-deletion-type: caught final_test.go:208: S13: done "the ruleset file, and its commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0ee3a309cbb06297
== rule-nff-type: caught final_test.go:208: S13: done "the ruleset file, and its commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0ee3a309cbb06297
== pr-resolution: caught final_test.go:208: S13: done "the ruleset file, and its commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0ee3a309cbb06297
== pr-merge-methods-squash: caught final_test.go:208: S13: done "the ruleset file, and its commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0ee3a309cbb06297
== pr-merge-methods-other: caught final_test.go:208: S13: done "the ruleset file, and its commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0ee3a309cbb06297
== pr-extra-approval: caught final_test.go:208: S13: done "the ruleset file, and its commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0ee3a309cbb06297
== pr-merge-methods-key: caught final_test.go:208: S13: done "the ruleset file, and its commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0ee3a309cbb06297
== pr-extra-key: caught final_test.go:208: S13: done "the ruleset file, and its commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0ee3a309cbb06297
== checks-strict: caught final_test.go:208: S13: done "the ruleset file, and its commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0ee3a309cbb06297
== default-refs-probe: caught final_test.go:208: S13: done "the ruleset file, and its commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0ee3a309cbb06297
== default-bypass: caught final_test.go:208: S13: done "the ruleset file, and its commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0ee3a309cbb06297
== records-ref: caught rulesets_test.go:144: the file of the records branch (err <nil>):
== records-name: caught rulesets_test.go:144: the file of the records branch (err <nil>):
== records-name-text: caught rulesets_test.go:144: the file of the records branch (err <nil>):
== records-enforcement: caught rulesets_test.go:144: the file of the records branch (err <nil>):
== records-target: caught rulesets_test.go:144: the file of the records branch (err <nil>):
== records-update-dropped: caught rulesets_test.go:144: the file of the records branch (err <nil>):
== records-update-param-true: caught rulesets_test.go:144: the file of the records branch (err <nil>):
== records-update-param-key: caught rulesets_test.go:144: the file of the records branch (err <nil>):
== records-deletion-dropped: caught rulesets_test.go:144: the file of the records branch (err <nil>):
== records-nff-dropped: caught rulesets_test.go:144: the file of the records branch (err <nil>):
== records-rule-order: caught rulesets_test.go:144: the file of the records branch (err <nil>):
== bypass-id: caught rulesets_test.go:144: the file of the records branch (err <nil>):
== bypass-type: caught rulesets_test.go:144: the file of the records branch (err <nil>):
== bypass-mode: caught rulesets_test.go:144: the file of the records branch (err <nil>):
== bypass-two: caught rulesets_test.go:144: the file of the records branch (err <nil>):
== bypass-none: caught rulesets_test.go:144: the file of the records branch (err <nil>):
== phase1-writes-records: caught final_test.go:208: S13: done "the ruleset files, and their commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0ee3a309cbb06
== phase2-writes-no-records: caught rulesets_test.go:281: S13 in phase 2: done "the ruleset file, and its commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0e
== records-path: caught rulesets_test.go:286: w/out/ruleset-records.json:
== row-name: caught rulesets_test.go:281: S13 in phase 2: done "the ruleset files, and their commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90
== row-source: caught rulesets_test.go:281: S13 in phase 2: done "the ruleset files, and their commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90
== row-ref: caught rulesets_test.go:281: S13 in phase 2: done "the ruleset files, and their commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90
== row-value: caught rulesets_test.go:281: S13 in phase 2: done "the ruleset files, and their commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90
== row-dropped: caught rulesets_test.go:281: S13 in phase 2: done "the ruleset files, and their commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90
== evidence-phase2: caught rulesets_test.go:281: S13 in phase 2: done "the ruleset file, and its commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0e
== evidence-phase1: caught final_test.go:208: S13: done "the ruleset files, and their commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0ee3a309cbb06
== commit-phase2: caught rulesets_test.go:281: S13 in phase 2: done "the ruleset files, and their commands in commands.sh" false, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe9
== commit-phase1: caught final_test.go:208: S13: done "the ruleset file, and its commands in commands.sh" false, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0ee3a309cbb0629
== appid-zero-accepted: caught rulesets_test.go:326: an App of ID 0: done "the ruleset file, and its commands in commands.sh", 2 files; want fail "S13: the app_id of the forge register is 0, not a posi
== appid-negative-accepted: caught rulesets_test.go:326: an App of a negative ID: done "the ruleset files, and their commands in commands.sh", 3 files; want fail "S13: the app_id of the forge register is -
== appid-reason: caught rulesets_test.go:326: an App of ID 0: fail "S13: the app_id of the forge register is 0, not a number", 0 files; want fail "S13: the app_id of the forge register is 0, not
== appid-reason-prefix: caught rulesets_test.go:326: an App of ID 0: fail "the app_id of the forge register is 0, not a positive number", 0 files; want fail "S13: the app_id of the forge register is 0,
== phase2-app-ignored: caught rulesets_test.go:281: S13 in phase 2: done "the ruleset files, and their commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90
== step-s13-not-replaced: caught phase2_test.go:252: S13: same code true, actor "layup-setup", hand-off true, reads []; want the code of phase 1 exactly when it is not one of ["S01" "S02" "S03" "S04" "S1
== step-s13-phase1-code: caught rulesets_test.go:281: S13 in phase 2: done "the ruleset file, and its commands in commands.sh" true, values [["branch-protection.sha256" "d2ef4e7d66a41995174e72c56fe90c0e
== step-s13-commands-dropped: caught rulesets_test.go:339: S13 in phase 2: actor "layup-setup", hand-off true, commands the same as phase 1's false; want the actor "layup-setup", the hand-off, and the comman
== step-s13-handoff-dropped: caught phase2_test.go:252: S13: same code false, actor "layup-setup", hand-off false, reads []; want the code of phase 1 exactly when it is not one of ["S01" "S02" "S03" "S04" "
== step-s13-actor: caught phase2_test.go:252: S13: same code false, actor "layup-run", hand-off true, reads []; want the code of phase 1 exactly when it is not one of ["S01" "S02" "S03" "S04" "S13
61 caught, 0 survived, 0 not run, of 61
```

## Green (08:32:35Z to 08:40:34Z, on `3235c26`)

Each with exit 0, in this order; the time of a line is when its command ended. The three
runs of `go test` are with `-count=1`, then with `-tags=integration` and with `-tags=e2e`,
`TMPDIR` is on the disk (`/tmp` is a tmpfs of 1.9 GB), and `git diff --check` reads the
commits against `origin/main`:

```
go build ./...: exit 0 (08:32:37Z)
go vet ./...: exit 0 (08:32:38Z)
gofmt -l internal cmd (empty): exit 0 (08:32:38Z)
go test ./...: exit 0 (08:32:47Z)
go test -tags=integration ./...: exit 0 (08:38:21Z)
go test -tags=e2e ./...: exit 0 (08:39:47Z)
adr-lint: exit 0 (08:39:48Z)
prd-lint: exit 0 (08:39:48Z)
link-lint: exit 0 (08:40:00Z)
run-discipline-tests: exit 0 (08:40:08Z)
nested-checkout-check: exit 0 (08:40:08Z)
setup-check: exit 0 (08:40:34Z)
git diff --check origin/main HEAD: exit 0 (08:40:34Z)
```
