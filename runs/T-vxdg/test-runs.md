# The test runs of T-vxdg

## Red 1: the tests before the package (2026-10-09T13:59Z)

Both test files of `internal/session` were written first; the package had no code, so `go vet ./internal/session/` does not build it:

```
# github.com/pharzam/layup/internal/session
# [github.com/pharzam/layup/internal/session]
vet: internal/session/session_test.go:36:8: undefined: Harness
```

## Red 2: a mutation of each rule, described

On copies of `session.go` and `rules_checker_test.go` (put back after), each rule broken alone; each makes its own cases fail. Three mutations first gave no failure: `refuse` (a second `Make` still failed, at the clone, so the test now asserts `fs.ErrExist`), and `sweepother` and `branch` (their first form did not compile; written again):

```
== lang: environ: no LANG
session_test.go:49: file: gives no variable: ["PATH=/usr/local/bin:/usr/bin" "HOME=/h/sessions/S-1a2b3c4d/home
session_test.go:49: no credential, no variable: ["PATH=/usr/local/bin:/usr/bin" "HOME=/h/sessions/S-1a2b3c4d/h
session_test.go:49: the fixed variables, in order: ["PATH=/usr/local/bin:/usr/bin" "HOME=/h/sessions/S-1a2b3c4
== path: environ: the host's HOME for PATH
session_test.go:49: file: gives no variable: ["PATH=hostile" "LANG=C.UTF-8" "HOME=/h/sessions/S-1a2b3c4d/home"
session_test.go:49: no credential, no variable: ["PATH=hostile" "LANG=C.UTF-8" "HOME=/h/sessions/S-1a2b3c4d/ho
session_test.go:49: the fixed variables, in order: ["PATH=hostile" "LANG=C.UTF-8" "HOME=/h/sessions/S-1a2b3c4d
== trimone: environ: every final line feed and carriage return cut
session_test.go:49: var: a carriage return stays: ["PATH=/usr/local/bin:/usr/bin" "LANG=C.UTF-8" "HOME=/h/sess
session_test.go:49: var: with two, one stays: ["PATH=/usr/local/bin:/usr/bin" "LANG=C.UTF-8" "HOME=/h/sessions
== novar: environ: the credential under another name
session_integration_test.go:137: Environ: ["PATH=/usr/local/go/bin:/home/layup/.local/bin:/home/layup/.opencod
session_test.go:49: var: with one final line feed: ["PATH=/usr/local/bin:/usr/bin" "LANG=C.UTF-8" "HOME=/h/ses
== vars: environ: no fixed variables
session_test.go:49: the fixed variables, in order: ["PATH=/usr/local/bin:/usr/bin" "LANG=C.UTF-8" "HOME=/h/ses
== target: Make: another target written
session_integration_test.go:125: after Sweep: ["S-00000001" "S-00000002"]; want the session of the other targe
== sweepnotarget: Sweep: a directory with no file target kept
session_integration_test.go:125: after Sweep: ["S-00000002" "S-00000003"]; want the session of the other targe
== filecopy: Make: no copy of a file: credential
session_integration_test.go:82: home/ holds [".gitconfig"]; want .gitconfig and the credential of file: only
session_integration_test.go:85: the copied credential: <nil>, stat /tmp/TestMakeASessionDirectory2667830779/00
== mode: Make: the credential copied with mode 0644
session_integration_test.go:85: the copied credential: &{credentials.toml 12 420 {103365047 63927151219 0x708e
== gitconfig: Make: another e-mail in .gitconfig
session_integration_test.go:89: home/.gitconfig: "[user]\n\tname = layup session S-1a2b3c4d\n\temail = S-1a2b3
== allow: TestInputRule: environ of internal/session not allowed
rules_integration_test.go:161: the module breaks the input rule:
rules_test.go:330: the findings:
== refuse: Make: an existing ID taken
session_integration_test.go:103: a second Make of the same ID: git clone --no-local --no-checkout --single-bra
== sweepother: Sweep: the directories of every target
session_integration_test.go:127: after Sweep: []; want the session of the other target only
== branch: Make: the branch of another attempt
session_integration_test.go:72: the refs of repo/: "refs/heads/main\nrefs/heads/task/T-ab12/2"; want main and
session_integration_test.go:75: repo/ is on "refs/heads/task/T-ab12/2", not task/T-ab12/1 at the base
```

## Green (2026-10-09T14:02Z)

Each with exit 0 on the tree of the commit `feat: T-vxdg …`: `go build ./...`, `go vet ./...`, `gofmt -l internal cmd` (empty), `go test ./...`, `go test -tags=integration ./...` (with `TestInputRule` and `TestPackageRules`); `adr-lint`, `prd-lint`, `link-lint`, `setup-check`, `run-discipline-tests`, `git diff --check`.
