# The test runs of T-5pxd

## Red 1: the tests before the functions (2026-10-09T17:31Z)

`go vet` does not build the tests, for want of the functions they name:

```
vet: internal/session/process_test.go:56:55: undefined: Steps
vet: internal/session/process_integration_test.go:48:66: undefined: Spec
```

## Red 2: a fake that could not see SIGINT (2026-10-09T17:33Z)

With the functions written, `TestProcessStopsTheWholeGroup` failed:

```
process_integration_test.go:114: the child saw ["TERM"], want INT and TERM
```

The child of the fake was a job of `&`, which a shell starts with `SIGINT`
ignored, and a shell cannot trap a signal ignored at its entry: a defect of the
fake, not of `Process`. The child now runs in the foreground, under a leader
whose traps are `:` (a lesson in `docs/guardrails.md` §2). It passed then.

## The mutations (2026-10-09T17:35Z to 17:38Z)

A mutation of each rule of `process.go`, one at a time, run by a script that
restores the file from a copy after each. The unit rules run the unit tests;
the others `go test -tags=integration -timeout 40s`. A mutation that removes a
stop hangs the fake that ignores its signals, so the test times out: that is
how it is caught. Each test line in full; of a timeout, the panic line and the
line of the test that waits:

```
== words-cap: Words: {cap} not replaced (exit 1)
process_test.go:23: Words("claude -p --model {model} --max-budget-usd {cap} {prompt}"): ["claude" "-p" "--model" "opus" "--max-budget-usd" "{cap}" "a b {model}"], want ["claude" "-p" "--model" "opus" "--max-budget-usd" "5" "a b {model}"]
== words-onepass: Words: {prompt} replaced first, then the others, in turn (exit 1)
process_test.go:23: Words("claude -p --model {model} --max-budget-usd {cap} {prompt}"): ["claude" "-p" "--model" "opus" "--max-budget-usd" "5" "a b opus"], want ["claude" "-p" "--model" "opus" "--max-budget-usd" "5" "a b {model}"]
== prompt-stdin: Prompt: the text for stdin too (exit 1)
process_test.go:32: Prompt(stdin): "the text", want "/h/sessions/S-1a2b3c4d/prompt.md"
== spec-stdin: NewSpec: prompt.md as the input of every mode (exit 1)
process_test.go:46: file: {Words:[h] Env:[PATH=/bin] Dir:/h/sessions/S-1a2b3c4d/repo Stdin:/h/sessions/S-1a2b3c4d/prompt.md Stdout:/h/sessions/S-1a2b3c4d/stdout Stderr:/h/sessions/S-1a2b3c4d/stderr Wall:30m0s IntWait:10s TermWait:10s StdoutCap:67108864 StderrCap:67108864 LineCap:8388608}, want {Words:[h] Env:[PATH=/bin] Dir:/h/sessions/S-1a2b3c4d/repo Stdin: Stdout:/h/sessions/S-1a2b3c4d/stdout Stderr:/h/sessions/S-1a2b3c4d/stderr Wall:30m0s IntWait:10s TermWait:10s StdoutCap:67108864 StderrCap:67108864 LineCap:8388608}
process_test.go:46: arg: {Words:[h] Env:[PATH=/bin] Dir:/h/sessions/S-1a2b3c4d/repo Stdin:/h/sessions/S-1a2b3c4d/prompt.md Stdout:/h/sessions/S-1a2b3c4d/stdout Stderr:/h/sessions/S-1a2b3c4d/stderr Wall:30m0s IntWait:10s TermWait:10s StdoutCap:67108864 StderrCap:67108864 LineCap:8388608}, want {Words:[h] Env:[PATH=/bin] Dir:/h/sessions/S-1a2b3c4d/repo Stdin: Stdout:/h/sessions/S-1a2b3c4d/stdout Stderr:/h/sessions/S-1a2b3c4d/stderr Wall:30m0s IntWait:10s TermWait:10s StdoutCap:67108864 StderrCap:67108864 LineCap:8388608}
== spec-linecap: NewSpec: a line of stdout cut at 8 KiB (exit 1)
process_test.go:46: file: {Words:[h] Env:[PATH=/bin] Dir:/h/sessions/S-1a2b3c4d/repo Stdin: Stdout:/h/sessions/S-1a2b3c4d/stdout Stderr:/h/sessions/S-1a2b3c4d/stderr Wall:30m0s IntWait:10s TermWait:10s StdoutCap:67108864 StderrCap:67108864 LineCap:8192}, want {Words:[h] Env:[PATH=/bin] Dir:/h/sessions/S-1a2b3c4d/repo Stdin: Stdout:/h/sessions/S-1a2b3c4d/stdout Stderr:/h/sessions/S-1a2b3c4d/stderr Wall:30m0s IntWait:10s TermWait:10s StdoutCap:67108864 StderrCap:67108864 LineCap:8388608}
process_test.go:46: arg: {Words:[h] Env:[PATH=/bin] Dir:/h/sessions/S-1a2b3c4d/repo Stdin: Stdout:/h/sessions/S-1a2b3c4d/stdout Stderr:/h/sessions/S-1a2b3c4d/stderr Wall:30m0s IntWait:10s TermWait:10s StdoutCap:67108864 StderrCap:67108864 LineCap:8192}, want {Words:[h] Env:[PATH=/bin] Dir:/h/sessions/S-1a2b3c4d/repo Stdin: Stdout:/h/sessions/S-1a2b3c4d/stdout Stderr:/h/sessions/S-1a2b3c4d/stderr Wall:30m0s IntWait:10s TermWait:10s StdoutCap:67108864 StderrCap:67108864 LineCap:8388608}
process_test.go:46: stdin: {Words:[h] Env:[PATH=/bin] Dir:/h/sessions/S-1a2b3c4d/repo Stdin:/h/sessions/S-1a2b3c4d/prompt.md Stdout:/h/sessions/S-1a2b3c4d/stdout Stderr:/h/sessions/S-1a2b3c4d/stderr Wall:30m0s IntWait:10s TermWait:10s StdoutCap:67108864 StderrCap:67108864 LineCap:8192}, want {Words:[h] Env:[PATH=/bin] Dir:/h/sessions/S-1a2b3c4d/repo Stdin:/h/sessions/S-1a2b3c4d/prompt.md Stdout:/h/sessions/S-1a2b3c4d/stdout Stderr:/h/sessions/S-1a2b3c4d/stderr Wall:30m0s IntWait:10s TermWait:10s StdoutCap:67108864 StderrCap:67108864 LineCap:8388608}
== call-order: Call: the context check before the version check (exit 1)
process_test.go:99: a task session: "S-1a2b3c4d", <nil>, the steps ["id" "attempt" "context" "version" "prompt" "rules" "start" "process" "end wall" "push"]; want ["id" "attempt" "version" "context" "prompt" "rules" "start" "process" "end wall" "push"]
process_test.go:104: a probe: <nil>, the steps ["id" "context" "version" "prompt" "rules" "start" "process" "end wall"]; want ["id" "version" "context" "prompt" "rules" "start" "process" "end wall"]
process_test.go:116: 2, a refusal at version: refused: version, the steps ["id" "attempt" "context" "version"]; want ["id" "attempt" "version"]
process_test.go:116: 3, a refusal at context: refused: context, the steps ["id" "attempt" "context"]; want ["id" "attempt" "version" "context"]
process_test.go:116: 4, a refusal at prompt: refused: prompt, the steps ["id" "attempt" "context" "version" "prompt"]; want ["id" "attempt" "version" "context" "prompt"]
process_test.go:116: 5, a refusal at rules: refused: rules, the steps ["id" "attempt" "context" "version" "prompt" "rules"]; want ["id" "attempt" "version" "context" "prompt" "rules"]
process_test.go:116: 6, a refusal at start: refused: start, the steps ["id" "attempt" "context" "version" "prompt" "rules" "start"]; want ["id" "attempt" "version" "context" "prompt" "rules" "start"]
process_test.go:116: 7, a refusal at end: a records push refused, the steps ["id" "attempt" "context" "version" "prompt" "rules" "start" "process" "end wall"]; want ["id" "attempt" "version" "context" "prompt" "rules" "start" "process" "end wall"]
== call-refusal: Call: a refusal of a check does not end the call (exit 1)
process_test.go:116: 1, a refusal at attempt: <nil>, the steps ["id" "attempt" "version" "context" "prompt" "rules" "start" "process" "end wall" "push"]; want ["id" "attempt"]
process_test.go:123: a refusal at attempt: <nil>, want its Refusal
process_test.go:116: 2, a refusal at version: <nil>, the steps ["id" "attempt" "version" "context" "prompt" "rules" "start" "process" "end wall" "push"]; want ["id" "attempt" "version"]
process_test.go:123: a refusal at version: <nil>, want its Refusal
process_test.go:116: 3, a refusal at context: <nil>, the steps ["id" "attempt" "version" "context" "prompt" "rules" "start" "process" "end wall" "push"]; want ["id" "attempt" "version" "context"]
process_test.go:123: a refusal at context: <nil>, want its Refusal
process_test.go:116: 4, a refusal at prompt: <nil>, the steps ["id" "attempt" "version" "context" "prompt" "rules" "start" "process" "end wall" "push"]; want ["id" "attempt" "version" "context" "prompt"]
process_test.go:123: a refusal at prompt: <nil>, want its Refusal
process_test.go:116: 5, a refusal at rules: <nil>, the steps ["id" "attempt" "version" "context" "prompt" "rules" "start" "process" "end wall" "push"]; want ["id" "attempt" "version" "context" "prompt" "rules"]
process_test.go:123: a refusal at rules: <nil>, want its Refusal
process_test.go:116: 6, a refusal at start: <nil>, the steps ["id" "attempt" "version" "context" "prompt" "rules" "start" "process" "end wall" "push"]; want ["id" "attempt" "version" "context" "prompt" "rules" "start"]
process_test.go:123: a refusal at start: <nil>, want its Refusal
== call-end: Call: an end that fails does not stop the push (exit 1)
process_test.go:116: 7, a refusal at end: <nil>, the steps ["id" "attempt" "version" "context" "prompt" "rules" "start" "process" "end wall" "push"]; want ["id" "attempt" "version" "context" "prompt" "rules" "start" "process" "end wall"]
== call-starterr: Call: the error of the process not given to End (exit 1)
process_test.go:138: <nil>; End got <nil>, want ErrStart
== version-trim: Version: the first line not trimmed (exit 1)
process_integration_test.go:34: the first line, trimmed, with the environment and in dir: "  2.1.295 (Claude Code) 001 ", <nil>; want "2.1.295 (Claude Code) 001"
process_integration_test.go:42: a blank line: <nil>, want a Refusal version
== version-exit: Version: a non-zero exit not refused (exit 1)
process_integration_test.go:42: a non-zero exit: <nil>, want a Refusal version
== version-empty: Version: an empty first line not refused (exit 1)
process_integration_test.go:42: no output: <nil>, want a Refusal version
process_integration_test.go:42: a blank line: <nil>, want a Refusal version
== process-stdin: Process: no standard input from Stdin (exit 1)
process_integration_test.go:69: stdin "/tmp/TestProcessGivesTheInputOfEachMode3715959100/001/prompt.md": {Start:2026-10-09 17:34:58.840058355 +0000 UTC m=+0.010054339 First:0001-01-01 00:00:00 +0000 UTC End:2026-10-09 17:34:58.842600219 +0000 UTC m=+0.012596193 Exit:0 Signal:signal 0 StoppedBy: Sent:[]}, <nil>, the input ""; want an end at once, with "the task\n"
== process-group: Process: the signals to the leader only, not its group (exit 1)
process_integration_test.go:116: the child saw [], want INT and TERM
process_integration_test.go:122: the child 1262375 lives after SIGKILL
== process-wall: Process: no stop at wall (exit 1)
panic: test timed out after 40s
/home/layup/projects/layup/.worktree/T-5pxd/internal/session/process_integration_test.go:88 +0x493
== process-ended: Process: each wait runs to its end, so the next signal is sent after the end of the process (exit 1)
process_integration_test.go:90: a fake that ends on SIGINT: {Start:2026-10-09 17:38:00.88635046 +0000 UTC m=+0.019976455 First:0001-01-01 00:00:00 +0000 UTC End:2026-10-09 17:38:01.988620684 +0000 UTC m=+1.122246695 Exit:3 Signal:signal 0 StoppedBy:wall Sent:[interrupt terminated killed]}, <nil>; want stopped at wall by [interrupt], exit 3, signal signal 0
process_integration_test.go:90: a fake that ignores SIGINT: {Start:2026-10-09 17:38:02.490513311 +0000 UTC m=+1.624139254 First:0001-01-01 00:00:00 +0000 UTC End:2026-10-09 17:38:03.592473775 +0000 UTC m=+2.726099719 Exit:4 Signal:signal 0 StoppedBy:wall Sent:[interrupt terminated killed]}, <nil>; want stopped at wall by [interrupt terminated], exit 4, signal signal 0
== process-kill: Process: no SIGKILL (exit 1)
panic: test timed out after 40s
/home/layup/projects/layup/.worktree/T-5pxd/internal/session/process_integration_test.go:88 +0x493
== process-filecap: Process: an output not cut at its cap (exit 1)
process_integration_test.go:144: stdout past its cap: {Start:2026-10-09 17:36:47.823357011 +0000 UTC m=+4.791047028 First:2026-10-09 17:36:47.824611564 +0000 UTC m=+4.792301584 End:2026-10-09 17:36:57.824426495 +0000 UTC m=+14.792116485 Exit:-1 Signal:interrupt StoppedBy:wall Sent:[interrupt]}, <nil> <nil>; want stopped by "output", stdout of 1000 bytes
process_integration_test.go:144: stderr past its cap: {Start:2026-10-09 17:36:57.824791025 +0000 UTC m=+14.792481025 First:0001-01-01 00:00:00 +0000 UTC End:2026-10-09 17:37:07.826139439 +0000 UTC m=+24.793829441 Exit:-1 Signal:interrupt StoppedBy:wall Sent:[interrupt]}, <nil> <nil>; want stopped by "output", stderr of 1000 bytes
== process-linecap: Process: no line cap (exit 1)
process_integration_test.go:144: a line of stdout past its cap: {Start:2026-10-09 17:37:13.87793424 +0000 UTC m=+4.751293760 First:2026-10-09 17:37:13.879698362 +0000 UTC m=+4.753057858 End:2026-10-09 17:37:13.880942693 +0000 UTC m=+4.754302189 Exit:-1 Signal:interrupt StoppedBy:output Sent:[interrupt]}, <nil> <nil>; want stopped by "output", stdout of 500 bytes
== process-lineeq: Process: a line of exactly the cap cut (exit 1)
process_integration_test.go:144: a line of stdout past its cap: {Start:2026-10-09 17:37:19.94025885 +0000 UTC m=+4.762095914 First:2026-10-09 17:37:19.941879568 +0000 UTC m=+4.763716632 End:2026-10-09 17:37:19.94238827 +0000 UTC m=+4.764225317 Exit:-1 Signal:interrupt StoppedBy:output Sent:[interrupt]}, <nil> <nil>; want stopped by "output", stdout of 500 bytes
process_integration_test.go:144: a line at its cap: {Start:2026-10-09 17:37:19.947562128 +0000 UTC m=+4.769399218 First:2026-10-09 17:37:19.950050127 +0000 UTC m=+4.771887192 End:2026-10-09 17:37:19.950495999 +0000 UTC m=+4.772333040 Exit:-1 Signal:interrupt StoppedBy:output Sent:[interrupt]}, <nil> <nil>; want stopped by "", stdout of 503 bytes
== process-stderrcap: Process: stderr cut at a cap 100 times its own (exit 1)
process_integration_test.go:144: stderr past its cap: {Start:2026-10-09 17:37:26.001976254 +0000 UTC m=+4.746655662 First:0001-01-01 00:00:00 +0000 UTC End:2026-10-09 17:37:26.107727427 +0000 UTC m=+4.852406662 Exit:-1 Signal:interrupt StoppedBy:output Sent:[interrupt]}, <nil> <nil>; want stopped by "output", stderr of 1000 bytes
== process-first: Process: the first output timed from stderr too (exit 1)
process_integration_test.go:152: {Start:2026-10-09 17:37:32.213105789 +0000 UTC m=+4.774968624 First:2026-10-09 17:37:32.214849287 +0000 UTC m=+4.776712138 End:2026-10-09 17:37:32.719463321 +0000 UTC m=+5.281326123 Exit:0 Signal:signal 0 StoppedBy: Sent:[]}, <nil>; want the first byte of stdout after 0.3 s, and the end 0.2 s later
process_integration_test.go:156: {Start:2026-10-09 17:37:32.720159317 +0000 UTC m=+5.282022140 First:2026-10-09 17:37:32.721853458 +0000 UTC m=+5.283716309 End:2026-10-09 17:37:32.722114014 +0000 UTC m=+5.283976870 Exit:0 Signal:signal 0 StoppedBy: Sent:[]}, <nil>; want no first output
== process-errstart: Process: a start that fails not ErrStart (exit 1)
process_integration_test.go:163: the program cannot start: fork/exec /tmp/TestProcessOfAMissingProgram731866720/001/none: no such file or directory, want ErrStart
```

The first run of `process-ended` (the check of the end before each signal
removed) was not caught: the wait before it checks the end too. The check
before each signal was dropped as a duplicate, and the mutation above removes
the check in the wait.

## The fix of round 1 (2026-10-09T17:55Z to 17:58Z)

Round 1's finding 1: the close of the outputs after `SIGKILL` had no case.
`TestProcessEndsWhenAProgramLeftTheGroup` starts a program that leaves the
group with `setsid` and holds `stdout`, and a leader that exits 0. Its second
assertion is note 2 of round 1: at 0.6 s, inside the stop, the program that
left looks for the leader in `/proc`. On the code of `179646f`, which reaped
the leader at its exit, it failed:

```
process_integration_test.go:189: the leader was reaped before the stop, so its group's ID was free: stat /tmp/TestProcessEndsWhenAProgramLeftTheGroup3707326542/002/held: no such file or directory
```

`Process` now waits for the leader's exit with `waitid` and `WNOWAIT`, and
reaps it only at the end; the test passes. The two new rules' mutations, each
line in full:

```
== process-close: Process: the outputs not closed after SIGKILL (exit 1)
process_integration_test.go:186: {Start:2026-10-09 17:57:22.754709454 +0000 UTC m=+5.278635980 First:0001-01-01 00:00:00 +0000 UTC End:2026-10-09 17:57:53.366985344 +0000 UTC m=+35.890911787 Exit:0 Signal:signal 0 StoppedBy:wall Sent:[interrupt terminated killed]}, <nil>, after 30.612894443s; want stopped at wall by [interrupt terminated killed], exit 0, at once
== process-reap: Process: the leader reaped at its exit, before the stop (exit 1)
process_integration_test.go:189: the leader was reaped before the stop, so its group's ID was free: stat /tmp/TestProcessEndsWhenAProgramLeftTheGroup751069979/002/held: no such file or directory
```

## Green (2026-10-09T17:40Z, again at 18:00Z after the fix of round 1)

Each with exit 0: `go build ./...`, `go vet ./...`, `gofmt -l internal cmd`
(empty), `go test ./...`, `go test -tags=integration ./...` (with
`TestPackageRules` and `TestInputRule`), the tests of the process five times
over (`-count=5`); `adr-lint`, `prd-lint`, `link-lint`, `setup-check`,
`run-discipline-tests`, `git diff --check`.
