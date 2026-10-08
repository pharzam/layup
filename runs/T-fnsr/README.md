# T-fnsr — the demo of M2a (uat)

The evidence of row 27 of the plan ([#132](https://github.com/pharzam/layup/issues/132)):
on the real GitHub target `pharzam/layup-uat`, the Operator ran Start and the
printed root push, and read the records branch with a plain `git clone`;
[`check.sh`](check.sh) checked each write of LAYUP. The uat row:
[`run.md`](../../docs/spec/run.md), "The acceptance tests of M2a".

**What this file holds.** The table and the progress lines of the run, the
Operator's read, and the test runs. The home directory is written as `~` (the
Operator's choice on #132, comment 6061142551); the 55 wait lines of step 4
are shortened to one line, and in the read the text of the README and of the
brief, which `check.sh` compared byte for byte; each shortening is marked. Published on the public target and here,
with the Operator's consent (the same comment): the host name
`antarctica.local` (in `lease.tsv`), the login `pharzam` and its ID 1675602,
and the key path of the forge register. No token and no key content.

## The inputs

| Input | Value |
| ----- | ----- |
| Binary | `go build -o ~/layup-host/bin/layup ./cmd/layup` at `76d52a1` (no `-ldflags`, so LAYUP's own pin is `armature.pin`); `layup 0.1.0-dev` |
| Host | `~/layup-host`: `registers/forge.tsv` (`github`, App 5118100, `layup-agent`, `~/.config/layup-agent/app.pem` mode 0600, `—`, `https://api.github.com`, `https://github.com`); `registers/harnesses.tsv` (the header only) |
| Target | `pharzam/layup-uat`, public, empty before the run; the App's installation 166044048 holds it |
| Brief | [`brief.md`](brief.md); no vision brief |

## The run

The Operator, in the Operator's own terminal, 2026-10-08, in the worktree of
the task (the first run, with `--intake-cap 5,1`, is test run 3 below):

```text
~/layup-host/bin/layup run --new pharzam/layup-uat --host ~/layup-host --psb <worktree>/runs/T-fnsr/brief.md --operator pharzam --idea-owner pharzam --plan free --intake-cap 5.0,1.0 --lease-h 60 --watch-t 1
```

Exit 0. Standard output, the table `run-steps`:

```tsv
step	result	detail
forge	done	the installation token, the six capabilities and the permissions of M2a; a public repository, empty
plan	done	the plan free on a public repository
baseline	done	the baseline at a95965534b14b0bf14ad74da0c9a45b5f4aedf88
root-push	done	the default branch main exists
read-back	done	one branch, one root commit, root tree = pin.tree
records	done	the first records commit: the pin, the briefs, approvers.tsv, start.tsv, the lease
issues	done	the Intake issue 1, the control issue 2
watch	done	not-confirmed
lease	done	released
```

Standard error:

```text
layup run: [1/9] forge
layup run: [2/9] plan
layup run: [3/9] baseline
layup run: [4/9] root-push
layup run: [4/9] root-push: root-push: the root commit is the unmodified baseline at a95965534b14b0bf14ad74da0c9a45b5f4aedf88, the commit of LAYUP's own pin a959655; push it with your own login: git -C ~/layup-host/roots/pharzam/layup-uat push -- https://github.com/pharzam/layup-uat.git main
layup run: [4/9] root-push: root-push: waiting for the push of the root commit   [55 such lines, one every ten seconds, shortened here]
layup run: [4/9] root-push: 150 s
layup run: [4/9] root-push: 300 s
layup run: [4/9] root-push: 450 s
layup run: [5/9] read-back
layup run: [6/9] records
layup run: [7/9] issues
layup run: [8/9] watch
layup run: [9/9] lease
```

The Operator pushed the root commit with the printed command and the
Operator's own login (`gh` as the credential helper, with the `workflow`
scope) about nine minutes after the line of the push.

**Two observations of the progress**, known since row 26, not defects of this demo: a wait line names its step twice
(`root-push: root-push:`); and a beat line (`150 s`, `300 s`, `450 s`) comes
every 150 seconds between the wait lines, so at those times a wait shows two
lines in ten seconds, where `run.md` says one. Both are known limits of
[`T-mqty.md`](../../docs/tasks/T-mqty.md) (round 1, note 2; round 3, note 1).

## The Operator's read

The Operator, in the Operator's own terminal, with `git` and `cat` only:

```text
rm -rf /tmp/layup-uat-read; git clone --branch layup-records https://github.com/pharzam/layup-uat /tmp/layup-uat-read; cd /tmp/layup-uat-read
begin; git log --format='%an <%ae> | %cn <%ce> | %s' layup-records; echo ---; git log --format='%an <%ae> | %cn <%ce> | %s' origin/main; echo ---; ls -R; for f in README.md start/start.tsv start/problem-statement.md approvers.tsv lease.tsv; echo "== $f"; cat $f; end; end
```

Output:

```text
layup-agent[bot] <335371832+layup-agent[bot]@users.noreply.github.com> | layup-agent[bot] <335371832+layup-agent[bot]@users.noreply.github.com> | lease: released
layup-agent[bot] <335371832+layup-agent[bot]@users.noreply.github.com> | layup-agent[bot] <335371832+layup-agent[bot]@users.noreply.github.com> | records: watch not-confirmed
layup-agent[bot] <335371832+layup-agent[bot]@users.noreply.github.com> | layup-agent[bot] <335371832+layup-agent[bot]@users.noreply.github.com> | records: issue.control 2
layup-agent[bot] <335371832+layup-agent[bot]@users.noreply.github.com> | layup-agent[bot] <335371832+layup-agent[bot]@users.noreply.github.com> | records: issue.control opening
layup-agent[bot] <335371832+layup-agent[bot]@users.noreply.github.com> | layup-agent[bot] <335371832+layup-agent[bot]@users.noreply.github.com> | records: issue.intake 1
layup-agent[bot] <335371832+layup-agent[bot]@users.noreply.github.com> | layup-agent[bot] <335371832+layup-agent[bot]@users.noreply.github.com> | records: issue.intake opening
layup-agent[bot] <335371832+layup-agent[bot]@users.noreply.github.com> | layup-agent[bot] <335371832+layup-agent[bot]@users.noreply.github.com> | records: the Start of pharzam/layup-uat
---
layup-agent[bot] <335371832+layup-agent[bot]@users.noreply.github.com> | layup-agent[bot] <335371832+layup-agent[bot]@users.noreply.github.com> | chore: the unmodified baseline at a95965534b14b0bf14ad74da0c9a45b5f4aedf88
---
README.md
approvers.tsv
lease.tsv
start

./start:
problem-statement.md
start.tsv
== README.md
[19 lines: the block start-readme of docs/spec/run.md, byte for byte (check.sh); shortened here]
== start/start.tsv
name	value	source
layup.version	0.1.0-dev	run
psb.sha256	f963794dfacdaead5c899b75b1234e6964cc0cedaef77303f6f28f8f05e266fa	command
vision.sha256	—	command
forge.plan	free	command
forge.visibility	public	forge
app.permissions	contents:write issues:write metadata:read pull_requests:write	forge
operator.id	1675602	forge
idea-owner.id	1675602	forge
intake.cap	5.0,1.0	command
lease.H	60	command
watch.T	1	command
pin.source	https://github.com/pharzam/armature	run
pin.commit	a95965534b14b0bf14ad74da0c9a45b5f4aedf88	run
pin.tree	8ffb250afd584da8b418bc220fb6d72e802924ce	run
pin.time	2026-10-08T17:34:57Z	run
issue.intake	1	forge
issue.control	2	forge
watch	not-confirmed	run
== start/problem-statement.md
[16 lines: brief.md, byte for byte (check.sh); shortened here]
== approvers.tsv
id	role	login	since	source
1675602	operator	pharzam	2026-10-08T17:44:51Z	start
1675602	idea-owner	pharzam	2026-10-08T17:44:51Z	start
== lease.tsv
run	host	version	started	heartbeat	state
8fcca82c0a5c8f54	antarctica.local	0.1.0-dev	2026-10-08T17:44:51Z	0	released
```

The issues page of the target shows #1 `LAYUP Intake` and #2 `LAYUP control`,
each opened by `layup-agent[bot]`.

So each write of LAYUP on the target is by the App's bot
`layup-agent[bot]` (ID 335371832): the root commit of `main`, the seven
records commits, and the two issues. The push of the root commit is the
Operator's, as `run.md` step 4 says.

## The test runs

Each run is on this host (macOS, `sh` of the system), in the worktree of the
task. The tables are the output of the scripts, shortened where marked.

### 1. `check.sh` on the empty target (red)

`sh runs/T-fnsr/check.sh https://github.com https://api.github.com pharzam/layup-uat`,
2026-10-08, before the run: exit 1.

```text
FAIL: no branch main
FAIL: no branch layup-records
2 failed
```

### 2. `check-test.sh`, the offline test of `check.sh` (condition 2)

First run: exit 1. Three faults of the test (the commit helper lost its `-m`;
`good` reused the loop variable `n`) and one of `check.sh`: `[bot]` in the
`grep` pattern of the authors is a bracket expression, so the check of the
records commits could not pass on a correct target
(`FAIL: 7 commits of layup-records are not by layup-agent[bot]`). Fixed with a
fixed-string comparison (`grep -cvxF`).

Run after the fixes: exit 0, 36 cases. A correct target and one with a
heartbeat commit pass; each other case gives its `FAIL` line once:

```text
ok: a correct target: exit 0
ok: a heartbeat commit: exit 0
ok: no repository: FAIL: the clone of
ok: no records branch: FAIL: no branch layup-records
ok: no bot response: FAIL: the ID of layup-agent[bot]: no answer
ok: two commits on main: FAIL: main has 2 commits
ok: a root commit by the Operator: FAIL: the root commit is by pharzam
ok: a root subject: FAIL: the root commit subject
ok: records by the Operator: FAIL: 7 commits of layup-records are not by
ok: a missing commit: FAIL: the commits of Start
ok: a vision brief: FAIL: the files:
ok: another README: FAIL: README.md is not the text
ok: another brief: FAIL: start/problem-statement.md is not the brief
ok: the header of start.tsv: FAIL: the header of start.tsv
ok: start.tsv <name>: FAIL: start.tsv <name> is 'x'   (14 cases, one per value)
ok: start.tsv app.permissions: FAIL: start.tsv app.permissions is empty
ok: start.tsv pin.time: FAIL: start.tsv pin.time is empty
ok: a harness row: FAIL: start.tsv has a harness row
ok: approvers.tsv: FAIL: approvers.tsv:
ok: the header of lease.tsv: FAIL: the header of lease.tsv
ok: a held lease: FAIL: lease.tsv: 6|0.1.0-dev|held
ok: the intake issue: FAIL: the intake issue #1 is by 'pharzam'
ok: the control issue: FAIL: the control issue #2 is by 'pharzam'
all cases pass
```

The response of `GET /users/layup-agent[bot]` is saved in
[`fixtures/users-bot.json`](fixtures/users-bot.json) (2026-10-08); the two
issue responses are written by the test.

### 3. The first Start: exit 2 before the first step

The Operator's first run, with `--intake-cap 5,1` (the author's proposal in
comment 6060791998), 2026-10-08: exit 2, no table, nothing written to the
target (`git ls-remote` gives no ref):

```text
layup: --intake-cap: "5,1" is not a value of the form that the block start gives intake.cap (MONEY,HOURS, two decimals)
```

This is the documented result (`run.md`, exit code 2 for a flag). The type
`decimal` needs a point (`internal/tsv/types.go`), so the value is `5.0,1.0`:
the same 5 USD and 1 hour. `check.sh` and `check-test.sh` expect `5.0,1.0`;
`check-test.sh` passes again (36 cases).

### 4. `check.sh` on the target after the run (green)

`sh runs/T-fnsr/check.sh https://github.com https://api.github.com pharzam/layup-uat`,
2026-10-08, after the run: exit 0, 33 lines `ok`, `all checks pass`. The lines
that carry the values of the run:

```text
ok: the ID of layup-agent[bot] is 335371832
ok: the root commit names the commit a95965534b14b0bf14ad74da0c9a45b5f4aedf88
ok: each commit of layup-records is by layup-agent[bot]
ok: the seven commits of Start, in order
ok: start.tsv psb.sha256 = f963794dfacdaead5c899b75b1234e6964cc0cedaef77303f6f28f8f05e266fa
ok: start.tsv pin.commit = a95965534b14b0bf14ad74da0c9a45b5f4aedf88
ok: start.tsv pin.tree = 8ffb250afd584da8b418bc220fb6d72e802924ce
ok: the intake issue #1 is by layup-agent[bot]
ok: the control issue #2 is by layup-agent[bot]
all checks pass
```

No heartbeat commit came (`lease.H` is 60 minutes), so the branch holds the
seven commits of Start only.
