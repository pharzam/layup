# The test runs of T-xxr8

Issue: [#457](https://github.com/pharzam/layup/issues/457). Baseline:
`f9d46b0f441202254e8ea6ca16fe4ad8ec860858`. Times are UTC.
The coordinating session supplied the host evidence; the sandbox runs below
were measured in this worktree. No commit is made by this stage.

## Plan amendment 1

Step 3 changes exactly these current-version statements:

- `go.mod`: `go 1.26.9` becomes `go 1.26.10`.
- `docs/spec/setup.md`: the parenthetical becomes `go 1.26.10` since task
  `T-xxr8`; the catalog's minor version stays `1.26`.
- `docs/guardrails.md`: the check's “names the full patch version” example
  becomes `go 1.26.10`.

The history of T-w89c in `docs/guardrails.md` (the advisories found in
`go1.26.8` and fixed in `go1.26.9`) stays unchanged. Every pre-existing
`1.26.9` match under `runs/`, in `docs/tasks/T-w89c.md`, and in the T-w89c
entry of `docs/tasks/completed.md` stays unchanged. There is no application
code change or duplicate CI version setting.

## Baseline security red

The coordinating session ran
`go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./...` on the baseline
with Go 1.26.9 on the host at `2026-10-21T05:42:25Z`, before the directive
and host toolchain changed. The coordinating session supplies shell exit **1**;
govulncheck's exit is **3**, reported by `go run` as `exit status 3`.
End time and elapsed time were not supplied.

Source read for this record:
`/home/layup/.cache/layup-session/m2e-scratchpad/T-xxr8/vuln-red.txt`.
These are excerpts of its result lines (download messages and the example
trace are omitted):

```text
Vulnerability #1: GO-2026-6809
    Distrust compromised intermediate certificates in crypto/x509
  More info: https://pkg.go.dev/vuln/GO-2026-6809
  Standard library
    Found in: crypto/x509@go1.26.9
    Fixed in: crypto/x509@go1.26.10
Your code is affected by 1 vulnerability from the Go standard library.
This scan found no other vulnerabilities in packages you import or modules you
require.
Use '-show verbose' for more details.
exit status 3
```

## Host toolchain step

The coordinating session performed `mise use -g go@1.26.10` on the host.
Source read: `/home/layup/.cache/layup-session/m2e-scratchpad/T-xxr8/host-step.txt`.
Its complete transcript is:

```text
host step at 2026-10-21T05:46:34Z
go version go1.26.9 linux/amd64
go version go1.26.10 linux/amd64
```

The supplied transcript has no command exit, end time or elapsed time; none
is inferred. The sandbox's initial `go version` also printed
`go version go1.26.10 linux/amd64` (exit 0).

## Initial sandbox attempts

Although PATH selected the Go 1.26.10 binary, the inherited `GOROOT` was
`/home/layup/.local/share/mise/installs/go/1.26.9`. All five commands below
failed before tests could execute. Representative output:

```text
# internal/unsafeheader
compile: version "go1.26.9" does not match go tool version "go1.26.10"
```

The build also reported that its default module cache was read-only:

```text
go: writing stat cache: open /home/layup/go/pkg/mod/cache/download/github.com/pharzam/layup/@v/v0.0.0-20261011204429-a4d41cf7b777.info688110345.tmp: read-only file system
```

| Command | Start (UTC) | End (UTC) | Elapsed (s) | Exit |
| --- | --- | --- | --- | --- |
| `go build ./...` | 2026-10-21T05:47:56.951024Z | 2026-10-21T05:47:57.114958Z | 0.164 | 1 |
| `go vet ./...` | 2026-10-21T05:47:57.115276Z | 2026-10-21T05:47:57.229198Z | 0.114 | 1 |
| `go test ./...` | 2026-10-21T05:47:57.229625Z | 2026-10-21T05:47:57.343501Z | 0.114 | 1 |
| `go test -count=1 -tags=integration -timeout 30m ./...` | 2026-10-21T05:47:57.343937Z | 2026-10-21T05:47:57.507962Z | 0.164 | 1 |
| `go test -tags=e2e ./...` | 2026-10-21T05:47:57.508432Z | 2026-10-21T05:47:57.622383Z | 0.114 | 1 |

## Corrected sandbox environment

For the rerun, `GOROOT` is removed from each command's inherited environment
and `GOMODCACHE` is set to
`/home/layup/.cache/layup-session/tmp-green/T-xxr8/modcache`.
No host setting or repository file is changed for this correction.
`env -u GOROOT go env GOVERSION GOROOT GOTOOLDIR` (exit 0) printed:

```text
go1.26.10
/home/layup/.local/share/mise/installs/go/1.26.10
/home/layup/.local/share/mise/installs/go/1.26.10/pkg/tool/linux_amd64
```

The commands run in the order below. Timing uses measured UTC start/end
and monotonic elapsed seconds; exits come directly from the child process,
without a pipeline. A step label and heartbeat every 15 seconds accompany
long commands. Integration uses `-count=1` and `-timeout 30m`, as instructed
for this stage and consistent with the loaded-host guardrail.

## Levels on Go 1.26.10

| Command | Start (UTC) | End (UTC) | Elapsed (s) | Exit |
| --- | --- | --- | --- | --- |
| `go build ./...` | 2026-10-21T05:48:18.751775Z | 2026-10-21T05:48:23.328127Z | 4.576 | 0 |
| `go vet ./...` | 2026-10-21T05:48:23.328430Z | 2026-10-21T05:48:24.494575Z | 1.166 | 0 |
| `go test ./...` | 2026-10-21T05:48:24.494867Z | 2026-10-21T05:48:28.067257Z | 3.572 | 0 |
| `go test -count=1 -tags=integration -timeout 30m ./...` | 2026-10-21T05:48:28.067594Z | 2026-10-21T05:48:52.105465Z | 24.038 | 1 |
| `go test -tags=e2e ./...` | 2026-10-21T05:48:52.105956Z | 2026-10-21T05:48:53.772867Z | 1.667 | 1 |

Build and vet produced no output. The unit level passed all 21 packages.
Integration failed in `internal/forge/github`, `internal/git`, `internal/run`
and `internal/setup`; E2E failed in `cmd/layup`. Each failed package panicked
when `httptest` tried to open a local listening socket. These are failed
runs, not green levels and not timeouts. The sandbox prohibits the sockets;
acceptance criterion 1 still needs successful host integration and E2E runs.

The following excerpts retain all package result lines and failing test names,
plus each panic line; goroutine stacks are omitted. Full local logs for both
attempts are under `/home/layup/.cache/layup-session/tmp-green/T-xxr8/`
(`build.log`, `vet.log`, `unit.log`, `integration.log`, `e2e.log`, and the
corresponding `*-corrected.log` files). Those scratch logs are not committed.

### Unit output

```text
ok  	github.com/pharzam/layup/cmd/layup	0.007s
ok  	github.com/pharzam/layup/internal/catalog	0.004s
ok  	github.com/pharzam/layup/internal/cli	0.752s
ok  	github.com/pharzam/layup/internal/forge	0.039s
ok  	github.com/pharzam/layup/internal/forge/github	0.117s
ok  	github.com/pharzam/layup/internal/gate	0.002s
ok  	github.com/pharzam/layup/internal/git	0.016s
ok  	github.com/pharzam/layup/internal/handoff	0.002s
ok  	github.com/pharzam/layup/internal/ledger	0.002s
ok  	github.com/pharzam/layup/internal/psb	0.003s
ok  	github.com/pharzam/layup/internal/records	0.010s
ok  	github.com/pharzam/layup/internal/route	0.003s
ok  	github.com/pharzam/layup/internal/rules	0.006s
ok  	github.com/pharzam/layup/internal/run	0.174s
ok  	github.com/pharzam/layup/internal/session	0.005s
ok  	github.com/pharzam/layup/internal/setup	0.331s
ok  	github.com/pharzam/layup/internal/spec	0.008s
ok  	github.com/pharzam/layup/internal/standin	0.001s
ok  	github.com/pharzam/layup/internal/tsv	0.003s
ok  	github.com/pharzam/layup/internal/verify	0.011s
ok  	github.com/pharzam/layup/internal/work	0.002s
```

### Integration output

```text
ok  	github.com/pharzam/layup/cmd/layup	0.498s
ok  	github.com/pharzam/layup/internal/catalog	4.607s
ok  	github.com/pharzam/layup/internal/cli	2.520s
ok  	github.com/pharzam/layup/internal/forge	0.029s
--- FAIL: TestTheAdapterPlaysEachCallOfM2a (0.00s)
panic: httptest: failed to listen on a port: listen tcp6 [::1]:0: socket: operation not permitted [recovered, repanicked]
FAIL	github.com/pharzam/layup/internal/forge/github	0.004s
ok  	github.com/pharzam/layup/internal/gate	0.158s
--- FAIL: TestNoCallUsesACredentialOfTheHost (0.00s)
panic: httptest: failed to listen on a port: listen tcp6 [::1]:0: socket: operation not permitted [recovered, repanicked]
FAIL	github.com/pharzam/layup/internal/git	0.220s
ok  	github.com/pharzam/layup/internal/handoff	0.004s
ok  	github.com/pharzam/layup/internal/ledger	0.004s
ok  	github.com/pharzam/layup/internal/psb	0.007s
ok  	github.com/pharzam/layup/internal/records	0.012s
ok  	github.com/pharzam/layup/internal/route	0.005s
ok  	github.com/pharzam/layup/internal/rules	0.836s
--- FAIL: TestTheBetStepIsAtomicToItsLastAct (0.02s)
    --- FAIL: TestTheBetStepIsAtomicToItsLastAct/version (0.02s)
panic: httptest: failed to listen on a port: listen tcp6 [::1]:0: socket: operation not permitted [recovered, repanicked]
FAIL	github.com/pharzam/layup/internal/run	0.028s
ok  	github.com/pharzam/layup/internal/session	7.357s
--- FAIL: TestTheUpdateKeepsWhatIsUnknown (0.03s)
panic: httptest: failed to listen on a port: listen tcp6 [::1]:0: socket: operation not permitted [recovered, repanicked]
FAIL	github.com/pharzam/layup/internal/setup	0.037s
ok  	github.com/pharzam/layup/internal/spec	0.084s
ok  	github.com/pharzam/layup/internal/standin	0.130s
ok  	github.com/pharzam/layup/internal/tsv	0.005s
ok  	github.com/pharzam/layup/internal/verify	19.564s
ok  	github.com/pharzam/layup/internal/work	0.003s
FAIL
```

### E2e output

```text
--- FAIL: TestTheStepsOfIntake (0.28s)
panic: httptest: failed to listen on a port: listen tcp6 [::1]:0: socket: operation not permitted [recovered, repanicked]
FAIL	github.com/pharzam/layup/cmd/layup	1.053s
ok  	github.com/pharzam/layup/internal/catalog	0.006s
ok  	github.com/pharzam/layup/internal/cli	0.667s
ok  	github.com/pharzam/layup/internal/forge	0.047s
ok  	github.com/pharzam/layup/internal/forge/github	0.058s
ok  	github.com/pharzam/layup/internal/gate	0.002s
ok  	github.com/pharzam/layup/internal/git	0.016s
ok  	github.com/pharzam/layup/internal/handoff	0.002s
ok  	github.com/pharzam/layup/internal/ledger	0.004s
ok  	github.com/pharzam/layup/internal/psb	0.003s
ok  	github.com/pharzam/layup/internal/records	0.013s
ok  	github.com/pharzam/layup/internal/route	0.003s
ok  	github.com/pharzam/layup/internal/rules	0.005s
ok  	github.com/pharzam/layup/internal/run	0.188s
ok  	github.com/pharzam/layup/internal/session	0.004s
ok  	github.com/pharzam/layup/internal/setup	0.351s
ok  	github.com/pharzam/layup/internal/spec	0.008s
ok  	github.com/pharzam/layup/internal/standin	0.002s
ok  	github.com/pharzam/layup/internal/tsv	0.003s
ok  	github.com/pharzam/layup/internal/verify	0.015s
ok  	github.com/pharzam/layup/internal/work	0.002s
FAIL
```

## Security green and freeze pending

The coordinating session will run the same govulncheck command on the host
after this stage's commits; this stage claims no security green. Its output,
exit and times, and the PR's green `security` job and run URL at a head before
the freeze, remain to be recorded in the next stage. Host integration and E2E
reruns must also unset the stale `GOROOT` or point it to Go 1.26.10.

The clean committed-tree freeze check, independent review, acceptance
checklist and close-out remain for subsequent stages. They were not performed
in Stage A, where commits are delegated to the coordinating session.

## Discipline checks

These checks ran after the three version edits. Scratch output is in
`/home/layup/.cache/layup-session/tmp-green/T-xxr8/discipline-1.log` through
`discipline-7.log`; TMPDIR was that same scratch directory.

| Command | Start (UTC) | End (UTC) | Elapsed (s) | Exit |
| --- | --- | --- | --- | --- |
| `sh docs/adr/adr-lint.sh` | 2026-10-21T05:49:11.388093Z | 2026-10-21T05:49:12.203315Z | 0.815 | 0 |
| `sh docs/prd/prd-lint.sh` | 2026-10-21T05:49:12.203521Z | 2026-10-21T05:49:12.235344Z | 0.032 | 0 |
| `sh docs/links/link-lint.sh` | 2026-10-21T05:49:12.235512Z | 2026-10-21T05:49:25.982487Z | 13.747 | 0 |
| `sh docs/tests/run-discipline-tests.sh` | 2026-10-21T05:49:25.982776Z | 2026-10-21T05:49:29.104554Z | 3.122 | 0 |
| `sh docs/tests/nested-checkout-check.sh` | 2026-10-21T05:49:29.104807Z | 2026-10-21T05:49:29.419402Z | 0.315 | 0 |
| `sh docs/tests/freeze-check/run-all.sh` | 2026-10-21T05:49:29.419616Z | 2026-10-21T05:49:38.657668Z | 9.238 | 0 |
| `sh docs/setup/setup-check.sh` | 2026-10-21T05:49:38.658006Z | 2026-10-21T05:50:01.542715Z | 22.885 | 0 |

Results: ADR and PRD lint OK; link lint resolved 3029 links; discipline
fixtures 138 passed, 0 failed; nested checkout 10 cases behaved; freeze-check
fixtures 85 passed, 0 failed. The full setup check printed OK for all twelve
checks and its kit linters. The fixture suite is not the task's clean-head
freeze check, which remains pending.

## Final scope verification

A comparison against the baseline confirmed exactly the specified substitution
in each of `go.mod`, `docs/spec/setup.md` and `docs/guardrails.md`; these are
the only changed tracked paths. The only new paths are this record and
`runs/T-xxr8/conditions.md`. All pre-existing records under `runs/`, the
T-w89c task and completed-log entry, and the catalog are unchanged. No
`checks.sh` was created or changed. This comparison exited 0.

## Final record checks

Link lint was rerun after both records were complete. The two no-index
comparisons include the new, untracked records, which ordinary diff omits.

| Command | Start (UTC) | End (UTC) | Elapsed (s) | Exit |
| --- | --- | --- | --- | --- |
| `sh docs/links/link-lint.sh` | 2026-10-21T05:50:51.421736Z | 2026-10-21T05:51:05.072387Z | 13.651 | 0 |
| `git diff --check` | 2026-10-21T05:51:05.072484Z | 2026-10-21T05:51:05.136589Z | 0.064 | 0 |
| `git diff --no-index --check /dev/null runs/T-xxr8/test-runs.md` | 2026-10-21T05:51:05.136665Z | 2026-10-21T05:51:05.140552Z | 0.004 | 1 |
| `git diff --no-index --check /dev/null runs/T-xxr8/conditions.md` | 2026-10-21T05:51:05.140600Z | 2026-10-21T05:51:05.142103Z | 0.002 | 1 |

Link lint printed `link-lint: OK  3029 links resolved`. Ordinary diff
reported no whitespace error. Both no-index comparisons returned 1 with
empty output (the new file differs from the empty file); no whitespace
diagnostic was emitted. A final direct byte check of all five planned files
also passed: no trailing spaces or tabs, no carriage returns, and exactly
one final newline, including this paragraph. That check exited 0.
