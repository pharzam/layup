# T-8s13: the demo of M2e (uat)

Row 68b of [the tasks of M2e](../../docs/plan/README.md#the-tasks-of-m2e), issue
[#370](https://github.com/pharzam/layup/issues/370). The plan (comment 6321274493), its review (6321481585, six conditions) and the
Operator's decisions (6321600725) are comments on #370; where a decision differs from the plan, the decision is the plan. This file is step T: the
scenario, the gates, the rules, the commands that read them, the failure policy, the inputs and the Operator's checklist. It is written and
pushed **before any run, before any register changes and before any write to the target**, and it is **never edited after a result is seen**: a
result goes into a new file or a new dated section below the line "Results", never into a rule.

The target is `pharzam/layup-uat-m2e` (public, empty at the survey), set up by the current release's Scaffold (the Operator, 6318463226 and
6318644376). The binary is the release of row 68a plus `T-36xw` (#471), built at the head of this branch after step T. The brief is
[`runs/T-sdtb/brief.md`](../T-sdtb/brief.md), and `brief.lines` stays 100. Each command below is run from a plain clone `C` of the target outside this
repository (`git clone https://github.com/pharzam/layup-uat-m2e.git C`); the records are `origin/layup-records` of `C`, never its checkout.
No credential is saved in any file of this directory. The agent runs no `gh` write and no push to a target.

## The scenario

```text
Scenario: The demo of M2e on a test target

Given the target pharzam/layup-uat-m2e, set up by the Scaffold, a brief that fits brief.lines 100 with room for the text, and two admitted harnesses with the roles routed
When the idea owner answers the questions of Intake, the Operator applies the setup and sets brief.lines after gate G2, and layup run goes through Shape and the loop
Then a specification task goes through the loop, Shape runs with its filter, the idea owner answers the brief with one comment, the rendered commit is added
And the task merges at its head SHA with the four layup/ statuses success
And each write of LAYUP on the target shows the App's bot as its author, the Operator's apply is the Operator's, and no review occurs on any pull request

Accepted: (the Operator's one word, on the evidence at a named commit, before the review round)
Covers REQ-005, REQ-007, REQ-012
```

### The steps and the gates

| Step | Run | Writes (actor) | Evidence |
| ---- | --- | -------------- | -------- |
| T | this file, `brief-bound.sh` and its test, pushed before any run | none | this file, `test-runs.md` |
| R0 | `layup run --new pharzam/layup-uat-m2e …` (the Operator's checklist, item 3) | Start's records and issues #1 and #2: bot; the root push: Operator | `runs.txt` |
| R1 to R3 | `layup run pharzam/layup-uat-m2e --host /home/layup/layup-host` to `scaffold-push` | records, the batch, the follow-up, session comments: bot; the answers: Operator | `runs.txt` |
| **G1** | `sh runs/T-8s13/brief-bound.sh C` after R3 | none | `brief-bound.txt` |
| item 6 | the Operator's `commands.sh` (the `main` push and the two applies) | Operator | `rulesets.txt` |
| R4 | restart; the forge check passes | records, the refused probe push: bot | `forge-check.txt` |
| R5 | `shape` to the two verified heads; `bet` waits `brief.lines is not set` | issues, branch pushes, draft PRs, statuses, PR body edits, comments: bot | `records.txt`, `forge-reads.txt` |
| **G2** | `sh runs/T-8s13/brief-bound.sh C HEAD`, then the Operator's `set brief.lines 100` | the `set`: Operator | `brief-bound.txt` |
| R5b | restart: the brief is posted and waits for the answer; the idea owner answers once | the brief and the annex commit: bot; the answer: Operator | `records.txt` |
| R6 | `bet`, then `loop` merges both tasks | rendered commit, statuses, ready mark, merges: bot | `forge-reads.txt` |
| R7 | repeat restart: `loop done done before` | no API write (the counts of comments, issues, pull requests, statuses, marks and merges are unchanged); lease and forge-check pushes only | `runs.txt` |
| item 8 | rule 5: the Operator's two pull requests and the merge attempt | Operator | `rule5.txt` |

**G1** (after R3, before any costly session of Shape). `brief-bound.sh C` prints the terms and `L0`, the lines of the parts of code with no decision, and the floor
`L0 + 3` (the brief needs at least one decision, and a decision has at least two options: three lines). The formula, as `internal/run/betbrief.go` writes the
brief (`assemble`, `requirementsPart`, `solutionPart`, `briefOpenPart`, `briefFactsPart`):

`L = 12 + R + max(M,1) + sum over decisions d of (1 + O_d) + Q_io + E + S + A_io + C_op`

R requirements (rows `requirement` and `non-functional` of the draft), M success metrics, O_d the options of decision d, Q_io the questions still open outside
the Operator's block, **E = 1 when an Operator question is still open (one range line) or no question is open (the line `- none`), and 0 otherwise**, S span facts, A_io the
answers outside the Operator's block, C_op the classes among the Operator's answers (one range line each). **Still open** is `stillOpen`
(`internal/run/shapequestions.go:25`): a question that is not a follow-up and whose state is `open`, or is `asked` with no answered follow-up that asks it again.
The test `brief-bound-test.sh` pins these terms, the 80/81 boundary and the floor (`test-runs.md`). **If the floor is above 80, the demo stops: no Shape session runs,
and the shortfall is reported to the Operator.** On the M2d records the floor is 79.

**G2** (after R5, before the Operator's `set`). `brief-bound.sh C HEAD`, with HEAD the verified head of the specification task, prints the exact `L` with the options
of `spec/options/*.tsv` and the success metrics of the PRD at HEAD. **The room, `brief.lines - L`, must be at least 20 (L at most 80)**; only then does the Operator post
the `set` (checklist item 5). Otherwise the demo stops and the Operator decides; the shortfall is a task of its own.

## The rules

The demo passes when each of the five rules holds on the target and fails when one is false. Rules 1, 3 and 4 apply to **every task that `loop` merges**; rule 2 to the
rendered commit. Times are UTC from the forge. **Ties (condition 5):** a status or a check counts before an event **only if its timestamp is strictly before the
event's**. Equal timestamps (the forge's resolution is one second) are **inconclusive, recorded as such, and the rule fails** unless independent ordering evidence
exists. Independent ordering evidence is a source other than those two timestamps that fixes their order, at finer resolution or by construction: a time-stamped log of
the run's forge requests with the request that made the status complete before the request that caused the event. `layup run` keeps none of API calls (the wrapper's
`push.log` of `T-sdtb`, if installed, orders pushes only), so unless such a log is produced and cited, a tie fails.

1. **The handoffs (condition 6).** For each merged task `T` (the tasks with an event `merged` in `tasks/T/events.tsv`), the actual transitions are reconciled with its
   events and with `transitions.tsv` (the table: `loop.md` "The transition table", whose default rows are the twelve of `architecture.md` §8 and the row of the target's
   architecture task; the step table `session.StepTable` gives each session's role). Read by `git -C C show origin/layup-records:tasks/T/events.tsv`,
   `…:sessions.tsv` and `…:transitions.tsv`. **Every handoff has a row, and every row is `valid`:**
   (a) every session of `T` in `sessions.tsv` (not a probe), except the session of the role `software-developer` that closes out (the close-out has no handoff of its own: the
   verifier's handoff to `close-out` is its row), has at least one event `handoff` naming it in the column `session`;
   (b) `T` has at least the four handoffs of its path (task plan to plan review, plan review to developer, developer to verifier, verifier to close-out), each `valid`, in this
   order, by the `n` of the events;
   (c) every `handoff` event has the detail `valid` (or `valid comment <ID>-<seen>`); an `invalid <item>` row, a missing row or a session with none fails the rule;
   (d) every pair that the handoffs used is a row of `transitions.tsv` (a pair with no row is `invalid pair`, which (c) already fails);
   the rendered commit's handoff, written by code, has the session `—` and is `valid` too.
   A task with a session and no handoff **fails**; so does a task with an invalid handoff.
2. **The rendered commit.** `git -C C diff --name-status R^ R` (R the rendered commit of the first bet, from `bets.tsv` and `git log --first-parent`) is exactly
   `A docs/facts/F-NNNN-numbered-facts.md` and `M docs/prd/PRD-*.md`; the changed lines of the PRD are only MoSCoW and Phase cells and the Non-goals paragraph
   (`git diff R^ R -- docs/prd` compared on bytes).
3. **The gates before human review and before the merge (condition 4 of the first plan review, tie rule above).** For each merged task, PR `N`, head `H`, merge time `merged_at`:
   - the human-review events: `gh api repos/pharzam/layup-uat-m2e/issues/$N/timeline --paginate --jq '.[]|select(.event=="review_requested" or .event=="ready_for_review" or .event=="merged")|[.event,.actor.login,(.requested_reviewer.login//""),.created_at]|@tsv'`;
     `t_hr` is the earliest `review_requested` or `ready_for_review` (the Operator's reading of `F-0003#58`);
   - each required native check: `…/commits/$H/check-runs?filter=latest&per_page=100 --jq '.check_runs[]|[.name,.head_sha,.app.id,.status,.conclusion,.completed_at]|@tsv'`; `head_sha` is `H`,
     the app 15368, `completed`/`success`, and `completed_at` **strictly before** `t_hr` and `merged_at`;
   - the four statuses: `…/commits/$H/statuses?per_page=100 --paginate --jq '.[]|[.context,.state,.creator.login,.created_at]|@tsv'`; among the statuses with `created_at` **strictly before** `t_hr`,
     and again among those strictly before `merged_at`, the newest of each of `layup/gates`, `layup/spec`, `layup/rules`, `layup/verify` is `success` by the bot;
   - every status at `H` created after `merged_at` is listed and counts for nothing (the late-status control of `cmd/layup/loop_e2e_test.go:588`);
   - **fails** on a status or check missing at either time, a review event before the gates were green, or a tie without independent evidence.
4. **The merge.** The pull request's `merged`, `merge_commit_sha` M and `head.sha` H; `git -C C rev-parse M^2` equals H; the first-parent merges of `main` are in the order of the tasks; the
   `merged` event of `events.tsv` holds M and N. The request body is not observable: its `sha` is the reviewed `Merge` call (`internal/run/ready.go:213`, `internal/forge/github/github.go:530`).
5. **A weakening pull request is refused.** `<G>` (branch `demo/weaken-gate`, edits `docs/gates.tsv`, `test` from `active` to `pending`) is `REVIEW_REQUIRED` and `BLOCKED`
   (`gh pr view <G> --json reviewDecision,mergeStateStatus`); the control `<C>` (`demo/control`, edits `README.md`) is not `REVIEW_REQUIRED`; the merge attempt
   `gh pr merge <G> -R pharzam/layup-uat-m2e --merge` exits non-zero and its text is kept.

### Shape (S1 to S4), for each decision d

- **S1, the applicable constraints: the specification's set before the first bet** (`loop.md` "The constraints that apply"; the Operator, 6321600725 (2)). Computed from the records and the target at the
  base: each `NFR-NNN` of the draft; each `F-<n>` span fact of class `constraint` (`spec/spans.tsv`); each `Inv-<n>` of `docs/guardrails.md` at the base; each accepted `ADR-NNNN` of
  `docs/adr/` at the base. Answer-derived `Q-NNN` constraint facts are **not** checked (the code does not read them before the first bet; nothing is repaired here). Listed in `shape.txt`.
- **S2.** Each option of `spec/options/<d>.tsv` has a row for every constraint of S1, or it is dropped as "missing".
- **S3.** The survivors and the drops by the filter: a drop for a `fail`, a missing row, or a `pass` with `—` in evidence. At least one survivor; no drop is required.
- **S4.** The `### Decision <d>` lines of the posted brief and the `survivor`/`dropped` rows of `bets/<m>/basis.tsv` agree with S3, option by option.

The brief: `bets/<m>/brief.md` is at most 100 lines with five headings; `bets.tsv` has one row with the verdict `bet`, approver 1675602, and one answer comment.

### AC2: the actors (all from authenticated forge answers, not from commit attribution)

- **Bot (`layup-agent[bot]`):** the creator of each pull request (`.user.login` of `pulls/$N`); each issue and comment of LAYUP (`user.login`, `performed_via_github_app.slug`); each status creator;
  `ready_for_review`, `merged_by` and the push actors (`/activity`); the author and committer of each records commit.
- **Operator (`pharzam`):** the root push and the `main` push; the two applies (`rulesets/<id>/history`); the answers, the `set` and the bet answer; the two pull requests of rule 5.
- **Reviews: none occur (the Operator, 6321600725 (1)).** No pull request of LAYUP touches a code-owner path, and the Operator is the only code owner and the author of `<G>`, so GitHub lets no review
  count. `gh api repos/pharzam/layup-uat-m2e/pulls/$N/reviews` is `[]` for **every** pull request of the demo; a non-empty answer fails the rule it touches and is reported.

### The failure policy

A false rule is a **failed attempt**, kept in this directory with its output, told apart from the passing one. The defect it shows is **a task of its own** (`loop.md`, "a defect that the demo finds is a task of its own");
no repair is made on this branch, and the e2e's repair rule does not apply. The row waits for that task; the replay is named in advance: a fix in a step after the last "done before" step is one restart with the new
binary; a fix in an earlier step needs its records and those after it read first, and may need a new target. A finding is classified only with the independent agreement that `engineering-discipline.md` requires.
A later change of the binary is recorded against the release (`git diff --stat 01d7da3082333efadb227afda14fd98715edf5c4 COMMIT -- cmd internal go.mod` and the release check on the built commit).
G1 and G2 stop before the costly sessions.

## The inputs

Recorded at the survey before Start (file `inputs.txt`), none of them guessed here:

- **The release:** `git diff --stat 01d7da3082333efadb227afda14fd98715edf5c4 76dfb684 -- cmd internal go.mod ':!*_test.go'` (names #471's changes); the binary is rebuilt with
  `go build -o /home/layup/layup-host/bin/layup ./cmd/layup` at the step-T head (the old binary kept as `layup.before-T-8s13`), then `sh runs/T-efmy/release-check.sh <full SHA>`.
- **The target:** `gh api repos/pharzam/layup-uat-m2e --jq '{visibility,size}'` and `…/branches` (public, empty, `[]`).
- **The host:** `cat /home/layup/layup-host/registers/*.tsv`; `codex --version`, `devin --version`, `claude --version`; `ls /home/layup/.local/bin/git` (absent, so `~/layup-host/wrap` stays off `PATH`); `gh auth status` (the host's `gh` is the Operator's login).
- **The brief:** `layup psb check runs/T-sdtb/brief.md` exits 0.
- **Open input, the codex context size:** the `context` value of the `codex` / `gpt-6.1-sol` row of `models.tsv` and its source URL are **PENDING from the Operator** (6321600725 (3)); they are not guessed here, and `codex` is not
  admitted until the Operator gives them. The flag `--dangerously-bypass-approvals-and-sandbox` is accepted by the Operator (a session writes `../result/` outside its clone).
- **Run settings:** `--operator pharzam --idea-owner pharzam --plan free --intake-cap 100.0,336.0 --lease-h 10 --watch-t 1`. The milestone-0 clock of the cap is checked at the entry of a phase session (`internal/run/phase.go:99`), of a Shape task
  session (`internal/run/shapeloop.go:556`) and through the close-out's `startAs` (`internal/run/finish.go:275`); `Probe` calls `TaskSession` without the check (`internal/run/probe.go:224`). `--lease-h 10` gives a takeover after 3 times H (30 minutes of an
  unchanged heartbeat) and a `HeldError` after 30 minutes of waiting (`internal/run/lease.go` 55 to 81).

## The Operator's checklist (own identity; the host is `LayupZAM`, user `layup`)

1. **Registers.** Back up: `ssh layup@LayupZAM 'cp -a /home/layup/layup-host/registers /home/layup/layup-host/registers.bak-2026-10-21'`.
   - Harness row (`harnesses.tsv`): `codex	—	30	codex exec --model {model} --dangerously-bypass-approvals-and-sandbox --skip-git-repo-check --ephemeral -	stdin	codex --version	/home/layup/.codex/auth.json	file:.codex/auth.json	AGENTS.md	—	none	subscription	—`
   - Models row (`models.tsv`): `codex	gpt-6.1-sol	PENDING-OPERATOR (context)	PENDING-OPERATOR (source URL)	PENDING-OPERATOR (date)	yes	—`
   - Routing: `claude` `claude-opus-5-5` (reasoning) for `software-engineer`, `software-developer`, `product-owner`, `systems-architect`, `domain-expert`; `claude` `claude-sonnet-5-5` for `software-developer` (execution);
     `qa-engineer` (reasoning): position 1 `codex` `gpt-6.1-sol`, position 2 `devin` `swe-2-high`. GPT runs only on Codex; Devin only with `swe-2-high`. The old `developer` routing rows stay.
2. **The binary (agent, host only):** build, keep the old one, run the release check (the inputs above).
3. **Start (agent):** `layup run --new pharzam/layup-uat-m2e --host /home/layup/layup-host --psb runs/T-sdtb/brief.md --operator pharzam --idea-owner pharzam --plan free --intake-cap 100.0,336.0 --lease-h 10 --watch-t 1`.
   **The Operator runs the printed root push.**
4. **Answers** on Intake issue #1, one comment per block, with `Q-NNN: answer` lines, the markers given values, not `gap`; answer the follow-up if one comes.
5. **After G2 only:** on control issue #2, `set brief.lines 100 because the Operator chose a 100-line brief on pharzam/layup#370, comment 6210244305`.
6. **Setup apply:** read, then run `ssh layup@LayupZAM 'sh /home/layup/layup-host/work/pharzam/layup-uat-m2e/out/commands.sh'` (the push of `layup-setup` and the POSTs of the two rulesets).
7. **The bet answer** on #1, one comment: a line for each requirement, `option <d>: <o>` for each decision with more than one survivor, and `bet`.
8. **After R7, rule 5:** open `<G>` (branch `demo/weaken-gate`) and `<C>` (branch `demo/control`), then `gh pr merge <G> -R pharzam/layup-uat-m2e --merge` (expected to fail).
9. **The one word `accept`**, on the evidence at a named commit, before the review round.
10. **Cleanup:** `gh pr close <G> --delete-branch` and `gh pr close <C> --delete-branch`. Keep the target.

## The register rows (what the Operator sets; the agent edits none of them)

| Register | Row | State |
| -------- | --- | ----- |
| `harnesses.tsv` | `codex`, wall 30 minutes, stdin, `--dangerously-bypass-approvals-and-sandbox` | to be added by the Operator (checklist 1) |
| `models.tsv` | `codex` `gpt-6.1-sol` | **context size, source and date PENDING from the Operator**; no value is invented |
| `routing.tsv` | the roles of checklist 1 | to be set by the Operator |

## Results

(Added after the runs, in dated sections below this line. The rules above are not edited.)
