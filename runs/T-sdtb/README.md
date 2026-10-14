# T-sdtb — the demo of M2d (uat)

Row 59b of [the tasks of M2d](../../docs/plan/README.md#the-tasks-of-m2d), issue
[#284](https://github.com/pharzam/layup/issues/284). The plan, its review and the
author's answer are comments on #284; where the answer differs from the plan, the answer is
the plan. This file is step T: its rules, its scenario, the commands that read each rule and
the policy for a failed attempt are written and pushed **before any run and before any
register changes**; the sections below them are added as the work goes on. The rules are not
edited after a result is seen.

## The rules, written before any run

The demo passes when each of the five rules holds on the target, and fails when one is false.
A step of a restart that is `fail` after the batch is posted is a failed **attempt**, not yet
the end of the demo ([After a failed attempt](#after-a-failed-attempt)). The target is
`pharzam/layup-uat-m2d`, public, made by Start from the host. Each command below is run from
the directory `runs/T-sdtb/` of the branch after `git clone
https://github.com/pharzam/layup-uat-m2d.git clone`; the records are `origin/layup-records` of
that clone, never its checkout.

1. **The forge check.** `setup/forge-check.tsv` of the records has the eight rows `tree`,
   `rules:<default branch>`, `rules:layup-probe`, `rules:layup-records`, `bypass:default`,
   `bypass:records`, `probe` and `activity`, each `pass`, and the table's `run` is the run that
   ended with "the forge check passed". Read by `git -C clone log --format='%H %s'
   origin/layup-records -- setup/forge-check.tsv`, then `git -C clone show COMMIT:setup/forge-check.tsv`
   for the commit of each table (the first check and the later start of rule 5); saved in
   `forge-check.txt`, each table with its commit and run ID, byte for byte.
2. **The refusal.** The row `probe` is `pass` by the rule of the specification
   ([`scaffold.md`](../../docs/spec/scaffold.md#the-probe)): the push to `refs/heads/layup-probe`
   exits non-zero and its standard error holds `GH013` and the ruleset's name
   `layup: the default branch`. The words are those of GitHub, kept by the wrapper of
   [the wrapper](#the-wrapper) (`~/layup-host/wrap/push.log`), and after the run `layup-probe` is no
   branch of the forge and no push event of it is in the activity. Read by the entry of the probe in
   the log, `git ls-remote https://github.com/pharzam/layup-uat-m2d.git refs/heads/layup-probe` (no
   output), `gh api 'repos/pharzam/layup-uat-m2d/activity?per_page=100' --paginate` (the push events of
   `refs/heads/layup-probe`: none) and, with the Operator's login, `gh api
   'repos/pharzam/layup-uat-m2d/rulesets/rule-suites?ref=refs/heads/layup-probe'` (which ruleset
   evaluated the push, and its result); saved in `probe-and-actors.txt`.
3. **The apply.** `setup/rulesets.txt` of the records holds both apply outputs, each a JSON object
   whose `source` is `pharzam/layup-uat-m2d`; the forge lists the two rulesets, `active`, named
   `layup: the default branch` and `layup: the records branch`. Read by `git -C clone show
   origin/layup-records:setup/rulesets.txt` (saved as `rulesets.txt`), and, with the Operator's login,
   `gh api repos/pharzam/layup-uat-m2d/rulesets` and `gh api repos/pharzam/layup-uat-m2d/rulesets/ID/history`
   for each of the two (the name, the enforcement and the actor of each apply); saved in
   `probe-and-actors.txt`.
4. **The actors.** Each commit of `layup-records` and of `main` has the App's bot, `layup-agent[bot]`,
   as author and committer; the Intake issue, the control issue and each comment of LAYUP are the bot's;
   the root push, the push of `layup-setup` to the default branch and the applies of the two rulesets are
   the Operator's, `pharzam`. The author of a commit and the actor of its push are read apart: the
   commits that the Operator pushes are LAYUP's authorship. Read by `git -C clone log
   --format='%an <%ae> | %cn <%ce> | %s' origin/layup-records origin/main`, the `user.login` of the two
   issues and of each comment (`gh api repos/pharzam/layup-uat-m2d/issues --paginate`, `gh api
   repos/pharzam/layup-uat-m2d/issues/comments --paginate`) and the actors of the push events of the activity
   read of rule 2; saved in `probe-and-actors.txt`.
5. **A later start.** One more restart, after the first check passed, passes the forge check again: its
   table `run-steps` has `forge-check` between `lease` and `probe`, `done`, and `setup/forge-check.tsv`
   has a new table whose commit is later than the first's. Read in `runs.txt` (the later start) and in the
   second table of `forge-check.txt`.

The row `activity` of LAYUP checks only that one page can be read; it is not evidence of actors here.
No credential is saved in any file of this directory.

## The scenario (uat)

```text
Scenario: The demo of M2d on a test target

Given the target pharzam/layup-uat-m2d, made by Start from the host with a short brief, and two registered harnesses, Claude Code and Devin, routed for the roles of Intake and the Scaffold
When the Operator runs the printed push of the root commit, posts the answers of the batch, and runs the printed commands.sh (the push of layup-setup and the apply of the two rulesets), and layup run restarts
Then the push of the probe to layup-probe is refused by GitHub, naming the ruleset
And setup/forge-check.tsv shows the eight rows as pass
And each write of LAYUP on the target shows the App's bot as its author, and the pushes and the applies are the Operator's

Accepted: pending
Covers REQ-002, REQ-003
```

The scenario rides the e2e path of row 58 (`TestTheStepsOfTheScaffold`,
[`runs/T-eevt/test-runs.md`](../T-eevt/test-runs.md)), on the real forge. The Operator accepts it on reading
this directory at a named commit, as for `M2b` and `M2c`; the `Accepted:` line is written when the answer
arrives, with who, when and where.

## The wrapper

The refusal of the probe is a line of standard error that `layup run` reads and keeps only as the first
line of a failure; a pass has an empty detail. So its words are kept outside the binary, which stays as the
release that row 59a reviewed, with its verdict rule: [`git-tee.sh`](git-tee.sh) is a transparent
wrapper of `git`, installed as `~/layup-host/wrap/git` and put first in the `PATH` of `layup run`. It runs the
real `git` (`/home/layup/.local/bin/git`, 2.56.0) with the same arguments, input, output and environment, and for
a `push` keeps in `~/layup-host/wrap/push.log` the UTC time, the exit status, the last two arguments (the URL and
the refspec: the token reaches `git` in the environment, not in an argument) and the whole standard error. The
probe is the run's own push; no push is made by hand.

[`git-tee-test.sh`](git-tee-test.sh) runs a push through the real `git` and through the wrapper, against a local
bare repository whose pre-receive hook refuses with the text of `GH013` and against one that accepts: 12 checks
(the same exit status, standard output and standard error; the log's two entries, their exit lines, the URL and
refspec, the text of `GH013`; a command that is not a push is passed on and not logged). Written with the wrapper,
before Start. The first run (11:30:13Z) failed one check on a defect of the test (the accepting repository was not
emptied between the two pushes); the second (11:30:21Z) passes all 12. Three mutants of the wrapper, each its own
copy, are caught by it: no replay of the standard error (2 checks fail), the exit status always 0 (1), no standard
error in the log (1).

## The brief

[`brief.md`](brief.md), 49 lines, written for this demo: a Go command that counts the lines, words and
characters of text files, with a table of measured targets and a terms table. `layup psb check brief.md` is
expected to find no gap; its output is recorded with the inputs. It is the short brief, not `PSB-CHAT-001` (1,132
lines): the demo is about the Scaffold, and the Operator was told (on 2026-10-14, 11:20Z) that a request for
`PSB-CHAT-001` before Start would change it. The Operator's reply is recorded with the inputs.

## After a failed attempt

Each finding of the demo is classified (in the path of this task, or revealed and off it), and the classification
is read by the round of the review. A defect of the build that makes a rule false is a **failed attempt**: the demo
is not signed off on it. The attempt stays in this directory with its verdict (its tables, with the failing row),
told apart from the passing attempt; the defect opens an issue, its fix is a linked task, and this task lands after
it. If GitHub's refusal lacks the name of the ruleset, the corrective task changes the check and its specification
together; this task never weakens a rule after it has seen a result, and the Operator's sign-off covers the final
passing attempt only. The replay after a fix is decided from the fix and the persisted state, and named here before it
runs: a fix in `forge-check` is one restart with the new binary (the earlier steps are "done before" by their records,
the table is written again, and the earlier refused probe left no ref); a fix in an earlier step needs its records and
those after it read first, and may need a new target. A later change of the binary is recorded against the reviewed
release (`git diff --stat 0a104b8 COMMIT -- cmd internal go.mod`, and the release check on the new head).

## The inputs (recorded on 2026-10-14, 11:32Z to 11:35Z, before Start)

- **The registers** ([`registers.txt`](registers.txt), as the runs read them, no change of them): the forge
  `github` with the App `layup-agent` (5118100); the harnesses `claude` (cap 5.00, wall 30 minutes) and `devin`
  (wall 15 minutes); the models `claude-sonnet-5-5`, `claude-opus-5-5` and `swe-2-high`; the routing of
  `systems-architect` and `product-owner` at the tier `reasoning` (`claude` `claude-opus-5-5`, then `devin`
  `swe-2-high`) and of `qa-engineer` (`devin`, then `claude`). The credential columns hold paths, not secrets.
- **The harnesses and git:** Claude Code 2.1.299 and Devin 3000.11.3 (9c803229faa4); git 2.56.0 at
  `/home/layup/.local/bin/git`, the real git of the wrapper.
- **The binary:** `go build -o ~/layup-host/bin/layup ./cmd/layup` at `00c9cdf` (this branch after step T);
  `layup version` gives `layup 0.1.0-dev`, its SHA-256 is
  `ddbae124e5de0e63b09bc363888c59a4cc07713b082e081237f2439da0cb7204`, and `git diff --stat 0a104b8 00c9cdf --
  cmd internal go.mod` prints nothing: its Go code is the release that row 59a reviewed. The binary before it
  is kept on the host as `layup.before-T-sdtb`.
- **The brief:** `layup psb check brief.md` on the committed file gives exit 0 and a header row only (11:33Z).
- **The wrapper:** installed as `~/layup-host/wrap/git`, its log `~/layup-host/wrap/push.log` empty, and
  `PATH=$HOME/layup-host/wrap:$PATH` for each `layup run`; `git --version` through it gives 2.56.0.
- **The run settings:** `--operator pharzam --idea-owner pharzam --plan free --intake-cap 50.0,48.0 --lease-h 10
  --watch-t 1`, as `M2c` set them but for the idea owner, who is the Operator's own account here.
