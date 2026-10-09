# The test runs of T-6sbe

## Red 1: the tests before the functions (2026-10-09T17:03Z)

`go vet ./internal/session/` does not build the tests:

```
# github.com/pharzam/layup/internal/session
# [github.com/pharzam/layup/internal/session]
vet: internal/session/checks_test.go:17:8: undefined: Refusal
```

## Red 2: a mutation of each rule, described

On a copy of `checks.go` (put back after), each of the fifteen rules broken alone; each makes its own cases fail. Each mutation shows up to three failing cases, cut at 110 characters. Two mutations first gave no failure, as another rule refused the same case with the same reason (`.git` that is not a directory; a loose ref that is a link); their cases now assert the value of the refusal, which names its rule, and they fail then (the last two groups).

```
== idform: NewID: 6 characters
checks_test.go:30: NewID: "S-a34e8e" "S-69bffa" (<nil> <nil>); want two IDs of the form S-xxxxxxxx
== roundup: CheckContext: no rounding up
checks_test.go:41: 401 bytes: estimate 100, want 101
checks_test.go:41: 5 bytes: estimate 1, want 2
checks_test.go:47: an estimate over the context: <nil>, want a Refusal context
== contextover: CheckContext: one over passes
checks_test.go:47: an estimate over the context: <nil>, want a Refusal context
== promptarg: CheckPromptSize: 131,072 passes
checks_test.go:59: 131,072 bytes as one argument: <nil>, want a Refusal prompt
== promptmode: CheckPromptSize: every mode limited
checks_test.go:62: 131,072 bytes by file: refused: prompt 131072, want no limit
checks_test.go:62: 131,072 bytes by stdin: refused: prompt 131072, want no limit
== rulewalk: CheckRuleFiles: every second directory
checks_integration_test.go:31: a rule file at /tmp/TestCheckRuleFiles2860999600/001/a/host/sessions: <nil>, wa
checks_integration_test.go:31: a rule file at /tmp/TestCheckRuleFiles2860999600/001/a: <nil>, want a Refusal r
== ruletop: CheckRuleFiles: no rule file refuses
checks_integration_test.go:31: a rule file at /tmp/TestCheckRuleFiles270827945/001/a/host/sessions: <nil>, wan
checks_integration_test.go:31: a rule file at /tmp/TestCheckRuleFiles270827945/001/a: <nil>, want a Refusal ru
== policy: CheckRuleFiles: every policy path given
checks_integration_test.go:25: no rule file above, a file inside repo/: ["/tmp/TestCheckRuleFiles8200618/001/p
== looseform: the loose ref's form: any spaces trimmed
checks_test.go:145: no line feed: <nil>, want a Refusal branch
checks_test.go:145: two line feeds: <nil>, want a Refusal branch
== symref: a symbolic loose ref taken
checks_test.go:145: a symbolic loose ref: <nil>, want a Refusal branch
== packedfield: packed-refs: a ref named by its prefix
checks_test.go:128: a packed ref beside attempt 10: "f123456789abcdef0123456789abcdef01234567", <nil>; want 01
checks_test.go:145: only the ref of attempt 10: <nil>, want a Refusal branch
== packedsha: packed-refs: a line of another form taken
checks_test.go:145: a packed line of another form: <nil>, want a Refusal branch
== packednone: no ref: no refusal
checks_test.go:145: only the ref of attempt 10: <nil>, want a Refusal branch
== gitdir (with the value of its refusal asserted)
checks_test.go:152: .git that is a file: refused "branch" "no ref refs/heads/task/T-ab12/1", want "branch" "re
checks_test.go:152: .git that is a link: refused "branch" "no ref refs/heads/task/T-ab12/1", want "branch" "re
== looselink (with the value of its refusal asserted)
checks_test.go:152: a loose ref that is a link: refused "branch" "refs/heads/task/T-ab12/1 is not a SHA and a
```

## The close-out: round 1's note 1 (2026-10-09T17:24Z)

The five malformed loose refs now have a packed line of the ref beside them,
and assert the value `… is not a SHA and a line feed`. The mutation: line 134
of `checks.go`, the refusal of a malformed loose ref, made `break`, so the read
falls through to `packed-refs`. Each of the five is caught, in full:

```
checks_test.go:158: a symbolic loose ref: <nil>, want a Refusal branch
checks_test.go:158: a short SHA: <nil>, want a Refusal branch
checks_test.go:158: no line feed: <nil>, want a Refusal branch
checks_test.go:158: upper case: <nil>, want a Refusal branch
checks_test.go:158: two line feeds: <nil>, want a Refusal branch
```

Note 2: `TestCheckRuleFiles` puts a rule file in the session directory itself
too. The mutation: the walk starts at the parent of `d.Root`. Caught, the line
cut at 110 characters:

```
checks_integration_test.go:31: a rule file at /tmp/TestCheckRuleFiles2033462972/001/a/host/sessions/S-1a2b
```

## Green (2026-10-09T17:05Z, again at 17:26Z after the close-out)

Each with exit 0 on the tree of the commit `feat: T-6sbe …`, and again on the close-out's: `go build ./...`, `go vet ./...`, `gofmt -l internal` (empty), `go test ./...`, `go test -tags=integration ./...` (with `TestPackageRules` and `TestInputRule`); `adr-lint`, `prd-lint`, `link-lint`, `setup-check`, `run-discipline-tests`, `git diff --check`.
