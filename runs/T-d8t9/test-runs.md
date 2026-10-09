# The test runs of T-d8t9

## Red 1: the unit tests before the functions (2026-10-09T18:37Z)

`go vet ./internal/run` does not build the tests, for want of the type they
name (two names of the test that clashed with the package's, `sha` and `at`,
were renamed first):

```
vet: internal/run/session_test.go:48:40: undefined: Sessions
```

## The integration test of the demo (2026-10-09T18:39Z)

`TestTheStartRowIsPushedBeforeTheProcess` was written after `session.go`, so
its red is the mutation `startrow-order` below: the start row committed after
the process starts, and the fake harness does not see it. Its first run failed
on the test itself (a clone whose HEAD was not `main`), fixed with
`git clone -b main`.

## Two cases made to fail on their own rule (2026-10-09T18:40Z)

The case of the next attempt had a session event of attempt 2, which gives 3
by either rule; it is now a refusal of attempt 7. The version check's "no
credential" had no case: the stand-in environment now gives the credential of
`var:` as `CRED`, which the process gets and the version check does not.

## The mutations (2026-10-09T18:40Z to 18:43Z)

A mutation of each rule of `session.go`, one at a time, run by a script that
restores the file from a copy after each; the unit rules run the unit tests,
`startrow-order` the integration test. Each test line in full:

```
== attempt-next: StartAttempt: the attempt is always 1 (exit 1)
session_test.go:72: attempt 2: 1, <nil>
== attempt-kind: StartAttempt: an event of any kind raises the attempt (exit 1)
session_test.go:84: after a refusal of attempt 7: 8, <nil>; want 3
== attempt-read: StartAttempt: an events.tsv that its reader refuses is read as empty (exit 1)
session_test.go:88: an events.tsv that its reader refuses: <nil>, 4 commits; want an error and no commit
== check-attempt: TaskSession: no attempt check (exit 1)
session_test.go:210: attempt: "S-1a2b3c4d", <nil>; want a refusal
session_test.go:213: attempt: the calls ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp" "admit 2.1.295 (Claude Code)" "rules" "process with secret" "end"]; want no process and the directory removed
session_test.go:216: attempt: 1 commits, sessions.tsv "session\ttask\tattempt\trole\tharness\tversion\tmodel\tbase\trecords\tprompt_bytes\tprompt_tokens\tcontext\tcap\twall\tvars\tpolicy\nS-1a2b3c4d\tT-ab12\t2\tdeveloper\tclaude\t2.1.295 (Claude Code)\tclaude-fable-5-1\t0123456789abcdef0123456789abcdef01234567\t0123456789abcdef0123456789abcdef01234567\t9\t3\t1000\t5.0\t30\tA=1\t/etc/claude/policy.json\n"; want one commit and no start row
session_test.go:221: attempt: the event ["2" "session" "2" "S-1a2b3c4d" "" "" "" "2026-10-09T12:00:00Z"]; want ["2" "refused" "2" "S-1a2b3c4d" "" "" "attempt" "2026-10-09T12:00:00Z"]
session_test.go:210: no events: "S-1a2b3c4d", <nil>; want a refusal
session_test.go:213: no events: the calls ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp" "admit 2.1.295 (Claude Code)" "rules" "process with secret" "end"]; want no process and the directory removed
session_test.go:216: no events: 1 commits, sessions.tsv "session\ttask\tattempt\trole\tharness\tversion\tmodel\tbase\trecords\tprompt_bytes\tprompt_tokens\tcontext\tcap\twall\tvars\tpolicy\nS-1a2b3c4d\tT-ab12\t1\tdeveloper\tclaude\t2.1.295 (Claude Code)\tclaude-fable-5-1\t0123456789abcdef0123456789abcdef01234567\t0123456789abcdef0123456789abcdef01234567\t9\t3\t1000\t5.0\t30\tA=1\t/etc/claude/policy.json\n"; want one commit and no start row
session_test.go:221: no events: the event ["1" "session" "1" "S-1a2b3c4d" "" "" "" "2026-10-09T12:00:00Z"]; want ["1" "refused" "1" "S-1a2b3c4d" "" "" "attempt" "2026-10-09T12:00:00Z"]
== version-home: TaskSession: the version check with HOME home/ (exit 1)
session_test.go:160: "S-1a2b3c4d", <nil>, the calls ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/home" "admit 2.1.295 (Claude Code)" "rules" "process with secret" "end"]; want ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp" "admit 2.1.295 (Claude Code)" "rules" "process with secret" "end"]
== version-cred: TaskSession: the version check with the credential (exit 1)
session_test.go:160: "S-1a2b3c4d", <nil>, the calls ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmpsecret" "admit 2.1.295 (Claude Code)" "rules" "process with secret" "end"]; want ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp" "admit 2.1.295 (Claude Code)" "rules" "process with secret" "end"]
== version-admit: TaskSession: no admission of the version read (exit 1)
session_test.go:160: "S-1a2b3c4d", <nil>, the calls ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp" "rules" "process with secret" "end"]; want ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp" "admit 2.1.295 (Claude Code)" "rules" "process with secret" "end"]
session_test.go:210: probe: "S-1a2b3c4d", <nil>; want a refusal
session_test.go:213: probe: the calls ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp" "rules" "process with secret" "end"]; want no process and the directory removed
session_test.go:216: probe: 1 commits, sessions.tsv "session\ttask\tattempt\trole\tharness\tversion\tmodel\tbase\trecords\tprompt_bytes\tprompt_tokens\tcontext\tcap\twall\tvars\tpolicy\nS-1a2b3c4d\tT-ab12\t1\tdeveloper\tclaude\t2.1.295 (Claude Code)\tclaude-fable-5-1\t0123456789abcdef0123456789abcdef01234567\t0123456789abcdef0123456789abcdef01234567\t9\t3\t1000\t5.0\t30\tA=1\t/etc/claude/policy.json\n"; want one commit and no start row
session_test.go:221: probe: the event ["2" "session" "1" "S-1a2b3c4d" "" "" "" "2026-10-09T12:00:00Z"]; want ["2" "refused" "1" "S-1a2b3c4d" "" "" "probe" "2026-10-09T12:00:00Z"]
== check-context: TaskSession: no context check (exit 1)
session_test.go:210: context: "S-1a2b3c4d", column context: the context is never below prompt_tokens; want a refusal
session_test.go:216: context: 0 commits, sessions.tsv ""; want one commit and no start row
session_test.go:221: context: the event ["1" "attempt" "1" "" "0123456789abcdef0123456789abcdef01234567" "" "" "2026-10-09T11:00:00Z"]; want ["1" "refused" "1" "S-1a2b3c4d" "" "" "context 3 2" "2026-10-09T12:00:00Z"]
== check-prompt: TaskSession: no prompt-size check (exit 1)
session_test.go:210: prompt: "S-1a2b3c4d", <nil>; want a refusal
session_test.go:213: prompt: the calls ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp" "admit 2.1.295 (Claude Code)" "rules" "process with secret" "end"]; want no process and the directory removed
session_test.go:216: prompt: 1 commits, sessions.tsv "session\ttask\tattempt\trole\tharness\tversion\tmodel\tbase\trecords\tprompt_bytes\tprompt_tokens\tcontext\tcap\twall\tvars\tpolicy\nS-1a2b3c4d\tT-ab12\t1\tdeveloper\tclaude\t2.1.295 (Claude Code)\tclaude-fable-5-1\t0123456789abcdef0123456789abcdef01234567\t0123456789abcdef0123456789abcdef01234567\t131072\t32768\t1048576\t5.0\t30\tA=1\t/etc/claude/policy.json\n"; want one commit and no start row
session_test.go:221: prompt: the event ["2" "session" "1" "S-1a2b3c4d" "" "" "" "2026-10-09T12:00:00Z"]; want ["2" "refused" "1" "S-1a2b3c4d" "" "" "prompt 131072" "2026-10-09T12:00:00Z"]
== check-rules: TaskSession: no rule-file check (exit 1)
session_test.go:210: rules: "S-1a2b3c4d", <nil>; want a refusal
session_test.go:213: rules: the calls ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp" "admit 2.1.295 (Claude Code)" "process with secret" "end"]; want no process and the directory removed
session_test.go:216: rules: 1 commits, sessions.tsv "session\ttask\tattempt\trole\tharness\tversion\tmodel\tbase\trecords\tprompt_bytes\tprompt_tokens\tcontext\tcap\twall\tvars\tpolicy\nS-1a2b3c4d\tT-ab12\t1\tdeveloper\tclaude\t2.1.295 (Claude Code)\tclaude-fable-5-1\t0123456789abcdef0123456789abcdef01234567\t0123456789abcdef0123456789abcdef01234567\t9\t3\t1000\t5.0\t30\tA=1\t—\n"; want one commit and no start row
session_test.go:221: rules: the event ["2" "session" "1" "S-1a2b3c4d" "" "" "" "2026-10-09T12:00:00Z"]; want ["2" "refused" "1" "S-1a2b3c4d" "" "" "rules /h/CLAUDE.md" "2026-10-09T12:00:00Z"]
== sweep: TaskSession: no sweep before the start (exit 1)
session_test.go:160: "S-1a2b3c4d", <nil>, the calls ["make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp" "admit 2.1.295 (Claude Code)" "rules" "process with secret" "end"]; want ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp" "admit 2.1.295 (Claude Code)" "rules" "process with secret" "end"]
== startrow-one: startRow: the start row and the event session in two commits (exit 1)
session_test.go:164: the commits 2; want one with sessions.tsv and events.tsv
== startrow-fenced: the records writes not through Fenced (exit 1)
session_test.go:160: "S-1a2b3c4d", the records push was refused, the calls ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp" "admit 2.1.295 (Claude Code)" "rules" "remove /h/sessions/S-1a2b3c4d"]; want ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp" "admit 2.1.295 (Claude Code)" "rules" "process with secret" "end"]
session_test.go:234: the records push was refused, the calls ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp" "admit 2.1.295 (Claude Code)" "rules" "remove /h/sessions/S-1a2b3c4d"]; want a LostError and no process
session_test.go:244: a lease of another run: the records push was refused
== startrow-order: TaskSession: the start row committed after the process starts (exit 1)
session_integration_test.go:63: the harness did not see its start row:
session_integration_test.go:66: the harness did not see the event session:
session_integration_test.go:69: the last records commit is "layup run: attempt 1 of T-ab12\n"; want the start row's
== refused-remove: TaskSession: a refused start keeps its directory (exit 1)
session_test.go:213: attempt: the calls ["sweep" "make"]; want no process and the directory removed
session_test.go:213: no events: the calls ["sweep" "make"]; want no process and the directory removed
session_test.go:213: version: the calls ["sweep" "make"]; want no process and the directory removed
session_test.go:213: probe: the calls ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp"]; want no process and the directory removed
session_test.go:213: context: the calls ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp" "admit 2.1.295 (Claude Code)"]; want no process and the directory removed
session_test.go:213: prompt: the calls ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp" "admit 2.1.295 (Claude Code)"]; want no process and the directory removed
session_test.go:213: rules: the calls ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp" "admit 2.1.295 (Claude Code)"]; want no process and the directory removed
session_test.go:237: the calls ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp" "admit 2.1.295 (Claude Code)" "rules"]; want the directory removed
== refused-event: TaskSession: a refused start commits no event (exit 1)
session_test.go:221: attempt: the event ["1" "attempt" "1" "" "0123456789abcdef0123456789abcdef01234567" "" "" "2026-10-09T11:00:00Z"]; want ["1" "refused" "2" "S-1a2b3c4d" "" "" "attempt" "2026-10-09T12:00:00Z"]
panic: runtime error: index out of range [-1] [recovered, repanicked]
/home/layup/projects/layup/.worktree/T-d8t9/internal/run/session_test.go:220 +0xb9a
== refused-detail: TaskSession: the refusal's value left out of the detail (exit 1)
session_test.go:221: context: the event ["2" "refused" "1" "S-1a2b3c4d" "" "" "context" "2026-10-09T12:00:00Z"]; want ["2" "refused" "1" "S-1a2b3c4d" "" "" "context 3 2" "2026-10-09T12:00:00Z"]
session_test.go:221: prompt: the event ["2" "refused" "1" "S-1a2b3c4d" "" "" "prompt" "2026-10-09T12:00:00Z"]; want ["2" "refused" "1" "S-1a2b3c4d" "" "" "prompt 131072" "2026-10-09T12:00:00Z"]
session_test.go:221: rules: the event ["2" "refused" "1" "S-1a2b3c4d" "" "" "rules" "2026-10-09T12:00:00Z"]; want ["2" "refused" "1" "S-1a2b3c4d" "" "" "rules /h/CLAUDE.md" "2026-10-09T12:00:00Z"]
== lost-remove: TaskSession: an error that is no refusal keeps the directory (exit 1)
session_test.go:237: the calls ["sweep" "make" "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp" "admit 2.1.295 (Claude Code)" "rules"]; want the directory removed
```

The first form of `refused-event` did not build, and proved nothing; the one
above commits the event to another file.

## The fix of round 1 (2026-10-09T18:52Z to 18:56Z)

Finding 1: a start whose `Make` failed after it made the directory kept the
directory. `TestAFailedMakeRemovesItsDirectory` gives a `Make` that fails, with
the root there before it (another session's: kept) and made by it (this
start's: removed). The test names the new seam `exists`, so it does not build
on `3db56c2`; the red of each of its two rules is its mutation, in full:

```
== make-partial: TaskSession: a Make that failed part-way keeps what it made (exit 1)
session_test.go:268: a directory there before false: the clone failed, the calls ["sweep" "make"], 0 commits; want an error, removed true, no commit
== make-existed: TaskSession: a directory that was there before removed (exit 1)
session_test.go:268: a directory there before true: the clone failed, the calls ["sweep" "make" "remove /h/sessions/S-1a2b3c4d"], 0 commits; want an error, removed false, no commit
```

## Green (2026-10-09T18:44Z, again at 18:57Z after the fix of round 1)

Each with exit 0: `go build ./...`, `go vet ./...`, `gofmt -l internal cmd`
(empty), `go test ./...`, `go test -tags=integration ./...` (with
`TestPackageRules` and `TestInputRule`); `adr-lint`, `prd-lint`,
`link-lint`, `setup-check`, `run-discipline-tests`, `git diff --check`.
