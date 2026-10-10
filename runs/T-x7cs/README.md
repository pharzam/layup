# T-x7cs — the demo of M2b (uat)

The evidence of row 39b of the plan ([#170](https://github.com/pharzam/layup/issues/170)):
on the real GitHub target `pharzam/layup-uat-m2b`, from the host `hetzam`, the
probe of each of two registered harnesses, then one developer session whose
commit lands on `task/T-rmgw/1` with its telemetry row. The uat row:
[`session.md`](../../docs/spec/session.md), The acceptance tests of M2b, The
demo. The home directory is written as `~`. No token and no key content is
here: the credential columns name files, and the session's environment is not
printed.

## The inputs

| Input | Value |
| ----- | ----- |
| Code | `layup` built at `d904465` (the branch `T-x7cs` on `33439ce`, with no change of a non-test Go file) into `~/layup-host/bin/layup`; `layup 0.1.0-dev` |
| Host | `hetzam`, `~/layup-host`: the registers below; no `prices.tsv` (an empty table) |
| Target | `pharzam/layup-uat-m2b`, public and empty, made by the Operator and added to the App's selected repositories (option (a) of comment 6095966713, as #199 refuses a restart of `pharzam/layup-uat`) |
| Harnesses | Claude Code 2.1.296 (the token of `claude setup-token`, by `var:CLAUDE_CODE_OAUTH_TOKEN`); Devin 3000.11.3 (its `credentials.toml`, a static API key, by `file:`), each mode 0600 |

`registers/harnesses.tsv`:

```
harness	cap	wall	command	prompt	version	credential	credential_to	rules	policy	usage	billing	vars
claude	1.00	15	claude -p --model {model} --max-budget-usd {cap} --output-format stream-json --verbose --permission-mode acceptEdits --allowedTools Bash,Edit,Write,Read	stdin	claude --version	~/layup-host/credentials/claude-oauth-token	var:CLAUDE_CODE_OAUTH_TOKEN	CLAUDE.md	—	claude-result	subscription	ANTHROPIC_DEFAULT_HAIKU_MODEL=claude-sonnet-5-5
devin	—	15	devin -p --model {model} --permission-mode dangerous --respect-workspace-trust false --prompt-file {prompt}	file	devin --version	~/.local/share/devin/credentials.toml	file:.local/share/devin/credentials.toml	AGENTS.md	—	none	subscription	—
```

`registers/models.tsv` (Devin publishes the context size of `swe-2-high`,
262K, only in `devin models list`; its `source` is the page of Devin's models,
which states no size: a known gap of the input, written here):

```
harness	model	context	source	date	use	reason
claude	claude-sonnet-5-5	1000000	https://platform.claude.com/docs/en/models/overview	2026-10-10T09:04:43Z	yes	—
devin	swe-2-high	262000	https://docs.devin.ai/cli/models	2026-10-10T09:04:43Z	yes	—
```

`registers/routing.tsv`:

```
role	tier	position	harness	model
developer	execution	1	claude	claude-sonnet-5-5
developer	execution	2	devin	swe-2-high
```

## Start (2026-10-10, 14:52Z to 14:58Z)

`~/layup-host/bin/layup run --new pharzam/layup-uat-m2b --host ~/layup-host --psb runs/T-fnsr/brief.md --operator pharzam --idea-owner pharzam --plan free --intake-cap 5.0,1.0 --lease-h 60 --watch-t 1`,
exit 0. The step `root-push` printed the push of the root commit; the Operator
ran it with the Operator's own login, over SSH. The table:

```
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

## The demo (2026-10-10, 14:58:11Z to 14:59:05Z)

`go test -tags=uat -run TestTheDemoOfM2b -timeout 60m -v ./internal/run -args -host ~/layup-host -target pharzam/layup-uat-m2b`,
exit 0, its output in full:

```
=== RUN   TestTheDemoOfM2b
    demo_uat_test.go:158: step forge: done, the installation token, the six capabilities and the permissions of M2a; a public repository, with layup-records
    demo_uat_test.go:159: step clone: done, start.tsv, approvers.tsv and lease.tsv read
    demo_uat_test.go:160: step version: done, LAYUP 0.1.0-dev
    demo_uat_test.go:161: step lease: done, held by this run
    demo_uat_test.go:171: step probe: done, probed and passed 2, probed and failed 0, skipped 0
    demo_uat_test.go:228: the task T-rmgw, attempt 1, the pair claude claude-sonnet-5-5, the base 657bc68df9205d243e8fd62e97ad9f6ddc936364
    demo_uat_test.go:230: the session S-a0460618: <nil>
    demo_uat_test.go:243: event ["1" "attempt" "1" "" "657bc68df9205d243e8fd62e97ad9f6ddc936364" "" "" "2026-10-10T14:58:47Z"]
    demo_uat_test.go:243: event ["2" "session" "1" "S-a0460618" "" "" "" "2026-10-10T14:58:48Z"]
    demo_uat_test.go:243: event ["3" "result" "1" "S-a0460618" "" "" "done" "2026-10-10T14:58:57Z"]
    demo_uat_test.go:243: event ["4" "push" "1" "S-a0460618" "" "dddf30c76b4a32485c6cb5d3a3c375e71dd25e03" "task/T-rmgw/1" "2026-10-10T14:58:58Z"]
    demo_uat_test.go:243: event ["5" "bound" "1" "S-a0460618" "" "dddf30c76b4a32485c6cb5d3a3c375e71dd25e03" "task/T-rmgw/1" "2026-10-10T14:59:01Z"]
    demo_uat_test.go:250: refs/heads/task/T-rmgw/1 on the target: dddf30c76b4a32485c6cb5d3a3c375e71dd25e03, <nil>
    demo_uat_test.go:261: telemetry ["S-a0460618" "T-rmgw" "" "developer" "claude" "claude-sonnet-5-5" "subscription" "2026-10-10T14:58:49Z" "2026-10-10T14:58:51Z" "2026-10-10T14:58:57Z" "2" "8" "4" "554" "41852" "observed" "" "" "" "unknown" ""]
    demo_uat_test.go:265: step phase: done, each step of Start is done; the lease is released
--- PASS: TestTheDemoOfM2b (52.41s)
PASS
ok  	github.com/pharzam/layup/internal/run	52.421s
```

## The records, read back

With a plain `git clone` of the target, after the run. `harnesses.tsv`:

```
session	harness	version	model	result	reason	files	models	end
S-017560b2	claude	2.1.296 (Claude Code)	claude-sonnet-5-5	passed	—	CLAUDE.md AGENTS.md	claude-sonnet-5-5	2026-10-10T14:58:26Z
S-1363ce3d	devin	devin 3000.11.3 (9c803229faa4)	swe-2-high	passed	—	AGENTS.md CLAUDE.md	—	2026-10-10T14:58:43Z
```

`sessions.tsv`:

```
session	task	attempt	role	harness	version	model	base	records	prompt_bytes	prompt_tokens	context	cap	wall	vars	policy
S-017560b2	T-at3u	1	probe	claude	2.1.296 (Claude Code)	claude-sonnet-5-5	657bc68df9205d243e8fd62e97ad9f6ddc936364	19b7492e9a31c19585714d9965cbd55ec92d11db	453	114	1000000	1.00	15	ANTHROPIC_DEFAULT_HAIKU_MODEL=claude-sonnet-5-5	—
S-1363ce3d	T-17b5	1	probe	devin	devin 3000.11.3 (9c803229faa4)	swe-2-high	657bc68df9205d243e8fd62e97ad9f6ddc936364	a00ab126ff3519cf32d29ef2360ac7d245efa1d0	453	114	262000	—	15	—	—
S-a0460618	T-rmgw	1	developer	claude	2.1.296 (Claude Code)	claude-sonnet-5-5	657bc68df9205d243e8fd62e97ad9f6ddc936364	a815da993bd9035a0750ab6cc5d8b632122fbef1	751	188	1000000	1.00	15	ANTHROPIC_DEFAULT_HAIKU_MODEL=claude-sonnet-5-5	—
```

`telemetry.tsv` (the money is `unknown`: both rows are subscriptions and the host has no `prices.tsv`; Devin reports no usage, so its tokens are `unavailable`):

```
session	task	requirements	role	harness	model	billing	start	first_output	end	latency_s	duration_s	tokens_in	tokens_out	tokens_cache	tokens_status	tokens_reason	money	currency	money_status	price
S-017560b2	T-at3u	—	probe	claude	claude-sonnet-5-5	subscription	2026-10-10T14:58:19Z	2026-10-10T14:58:21Z	2026-10-10T14:58:26Z	2	7	4	344	41319	observed	—	—	—	unknown	—
S-1363ce3d	T-17b5	—	probe	devin	swe-2-high	subscription	2026-10-10T14:58:30Z	2026-10-10T14:58:30Z	2026-10-10T14:58:43Z	0	13	—	—	—	unavailable	the harness reports none	—	—	unknown	—
S-a0460618	T-rmgw	—	developer	claude	claude-sonnet-5-5	subscription	2026-10-10T14:58:49Z	2026-10-10T14:58:51Z	2026-10-10T14:58:57Z	2	8	4	554	41852	observed	—	—	—	unknown	—
```

`tasks/T-rmgw/events.tsv`:

```
n	kind	attempt	session	base	sha	detail	time
1	attempt	1	—	657bc68df9205d243e8fd62e97ad9f6ddc936364	—	—	2026-10-10T14:58:47Z
2	session	1	S-a0460618	—	—	—	2026-10-10T14:58:48Z
3	result	1	S-a0460618	—	—	done	2026-10-10T14:58:57Z
4	push	1	S-a0460618	—	dddf30c76b4a32485c6cb5d3a3c375e71dd25e03	task/T-rmgw/1	2026-10-10T14:58:58Z
5	bound	1	S-a0460618	—	dddf30c76b4a32485c6cb5d3a3c375e71dd25e03	task/T-rmgw/1	2026-10-10T14:59:01Z
```

`tasks/T-rmgw/results/S-a0460618.tsv`:

```
kind	n	value	sha256	reason
status	1	completed	—	the file is written and committed
artifact	1	demo/T-rmgw.txt	1e6377317fea39073120e321cdbe33f93423512717ec8a40c3282f709d37366e	—
```

`rule-paths.tsv` is the register of D2 of the plan, committed by the test
(`21d8800`), as the target has no setup.

## The task branch

```
dddf30c76b4a32485c6cb5d3a3c375e71dd25e03 layup session S-a0460618 <S-a0460618@sessions.layup.invalid>
  docs: T-rmgw the file of the demo
 demo/T-rmgw.txt | 1 +
 1 file changed, 1 insertion(+)

demo/T-rmgw.txt: The demo of LAYUP M2b, task T-rmgw.
its SHA-256: 1e6377317fea39073120e321cdbe33f93423512717ec8a40c3282f709d37366e
```

The commit's author is the session's (`session.md`, The session directory); its
parent is the base `657bc68`, the root commit of `main`.

## The authors (criterion 2)

Each records commit of the target, by `git log`:

```
3091875 layup-agent[bot] | layup-agent[bot] | lease: released
aba7985 layup-agent[bot] | layup-agent[bot] | layup run: S-a0460618 bound to task/T-rmgw/1
6dd31d5 layup-agent[bot] | layup-agent[bot] | layup run: the push of S-a0460618 to task/T-rmgw/1
1d57a90 layup-agent[bot] | layup-agent[bot] | layup run: the end of S-a0460618, done
a5f592d layup-agent[bot] | layup-agent[bot] | layup run: the start of S-a0460618, attempt 1 of T-rmgw
a815da9 layup-agent[bot] | layup-agent[bot] | layup run: attempt 1 of T-rmgw
21d8800 layup-agent[bot] | layup-agent[bot] | layup run: the rule-path register of the demo of M2b
ce20530 layup-agent[bot] | layup-agent[bot] | layup run: the probe S-1363ce3d of devin
26effc3 layup-agent[bot] | layup-agent[bot] | layup run: the start of the probe S-1363ce3d
a00ab12 layup-agent[bot] | layup-agent[bot] | layup run: the probe S-017560b2 of claude
e785298 layup-agent[bot] | layup-agent[bot] | layup run: the start of the probe S-017560b2
19b7492 layup-agent[bot] | layup-agent[bot] | layup run: the routing register
4dfc53c layup-agent[bot] | layup-agent[bot] | lease: the run 18dd32fde4d88d20 takes the released lease from the run 7a86133de2d85c85
94521d5 layup-agent[bot] | layup-agent[bot] | lease: released
93dec9b layup-agent[bot] | layup-agent[bot] | records: watch not-confirmed
d5b0818 layup-agent[bot] | layup-agent[bot] | records: issue.control 2
58ad50f layup-agent[bot] | layup-agent[bot] | records: issue.control opening
25a58de layup-agent[bot] | layup-agent[bot] | records: issue.intake 1
2f5dcb3 layup-agent[bot] | layup-agent[bot] | records: issue.intake opening
da5c4a2 layup-agent[bot] | layup-agent[bot] | records: the Start of pharzam/layup-uat-m2b
```

The repository's activity (`GET /repos/pharzam/layup-uat-m2b/activity`): each
push and branch creation is by `layup-agent[bot]`, the push of
`task/T-rmgw/1` among them, but the push of `main`, which `root-push` gives the
Operator:

```
2026-10-10T14:57:46Z branch_creation pharzam refs/heads/main
2026-10-10T14:57:53Z branch_creation layup-agent[bot] refs/heads/layup-records
2026-10-10T14:57:54Z push layup-agent[bot] refs/heads/layup-records
2026-10-10T14:57:55Z push layup-agent[bot] refs/heads/layup-records
2026-10-10T14:57:57Z push layup-agent[bot] refs/heads/layup-records
2026-10-10T14:57:59Z push layup-agent[bot] refs/heads/layup-records
2026-10-10T14:58:00Z push layup-agent[bot] refs/heads/layup-records
2026-10-10T14:58:01Z push layup-agent[bot] refs/heads/layup-records
2026-10-10T14:58:16Z push layup-agent[bot] refs/heads/layup-records
2026-10-10T14:58:17Z push layup-agent[bot] refs/heads/layup-records
2026-10-10T14:58:19Z push layup-agent[bot] refs/heads/layup-records
2026-10-10T14:58:28Z push layup-agent[bot] refs/heads/layup-records
2026-10-10T14:58:30Z push layup-agent[bot] refs/heads/layup-records
2026-10-10T14:58:44Z push layup-agent[bot] refs/heads/layup-records
2026-10-10T14:58:46Z push layup-agent[bot] refs/heads/layup-records
2026-10-10T14:58:48Z push layup-agent[bot] refs/heads/layup-records
2026-10-10T14:58:49Z push layup-agent[bot] refs/heads/layup-records
2026-10-10T14:58:58Z push layup-agent[bot] refs/heads/layup-records
2026-10-10T14:58:59Z push layup-agent[bot] refs/heads/layup-records
2026-10-10T14:59:00Z branch_creation layup-agent[bot] refs/heads/task/T-rmgw/1
2026-10-10T14:59:02Z push layup-agent[bot] refs/heads/layup-records
2026-10-10T14:59:04Z push layup-agent[bot] refs/heads/layup-records
```

The issues and the comments of the control issue #2, each by
`layup-agent[bot]`:

```
#2 layup-agent[bot] LAYUP control
#1 layup-agent[bot] LAYUP Intake
layup-agent[bot] 2026-10-10T14:58:28Z | S-017560b2: probe of claude 2.1.296 (Claude Code): passed
layup-agent[bot] 2026-10-10T14:58:45Z | S-1363ce3d: probe of devin devin 3000.11.3 (9c803229faa4): passed
layup-agent[bot] 2026-10-10T14:59:02Z | S-a0460618: developer session of T-rmgw, attempt 1: done
```

## The scenario (uat)

```text
Scenario: The demo of M2b on real harnesses

Given the target pharzam/layup-uat-m2b, started from the host hetzam with two registered harnesses, Claude Code and Devin, each with its credential
When the uat test runs the restart's steps and, with the lease held, one developer session of a fixed task
Then each harness is probed at its version and passes
And the developer session's commit lands on task/T-rmgw/1, with its result, its events attempt, session, result, push and bound, and its telemetry row
And each write of LAYUP on the target shows the App's bot as its author

Accepted by: pending, the Operator on #170
Covers REQ-013, REQ-003, REQ-005, REQ-011, NFR-001
```
