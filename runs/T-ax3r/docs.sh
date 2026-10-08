#!/bin/sh
# The check of the documents of task T-ax3r (#142, row 25b of the plan). A
# one-off evidence command of this task, not a check of the gate: no hook and
# no CI job runs it. Run from the root of the repository; exit 0 when each rule
# holds, 1 otherwise.
fail=0
ok() { printf 'ok   %s\n' "$1"; }
no() { printf 'FAIL %s\n' "$1"; fail=1; }
has() { grep -q -- "$2" "$1" && ok "$1: $3" || no "$1 does not say: $3"; }
# 1. run.md: the decisions of the plan and the conditions of its review.
has docs/spec/run.md 'the baseline is the source of LAYUP.s own pin' "decision 1, the baseline"
has docs/spec/run.md 'git -C DIR/roots/OWNER/NAME push -- <web>/OWNER/NAME.git main' "decision 2, the root commit and its command"
has docs/spec/run.md 'the author of each records commit' "decision 3 and condition 3, the bot's ID"
has docs/spec/run.md 'the run never builds on that branch' "decision 4 and condition 1, the base of a commit"
has docs/spec/run.md 'a lost beat stops the step that runs' "condition 2, the heartbeat"
has docs/spec/run.md 'every ten seconds,$' "condition 5, the period of step 8"
has docs/spec/run.md 'an earlier Start that stopped before step 6' "condition 6, a row of the input states"
has docs/spec/run.md 'leaves the lease `held`' "note 7, the known limit"
has docs/spec/run.md '```text intake-issue' "note 8, the block of the Intake issue"
has docs/spec/run.md '```text control-issue' "note 8, the block of the control issue"
has docs/spec/run.md '`watch.T` of the restart are those of `start.tsv`' "finding 1 of round 1, watch.T of a restart"
# 2. forge.md and packages.md: condition 4, decisions 5 and 6.
grep '^| `Repository` |' docs/spec/forge.md | grep -q 'the names of the branches' && ok "forge.md: Repository gives the branch names" || no "forge.md: Repository does not give the branch names"
has docs/spec/forge.md 'follows no redirect' "decision 6, the client"
grep '^| `Clone` |' docs/spec/packages.md | grep -q 'with a token' && ok "packages.md: Clone with a token" || no "packages.md: Clone takes no token"
# 3. The traceability and the Test cells.
for t in TestStartMakesTheFirstRecordsCommit TestARestartAfterARunThatStoppedAtOpening TestALostBeatEndsTheWatch TestTheForgeStepOfStart; do
	grep "\`$t\`" docs/tests/traceability.md | grep -q 'T-ax3r' && ok "traceability: $t" || no "traceability has no row of $t with T-ax3r"
done
for r in 'NFR-001 | F-0001#1 ' 'NFR-002 | ' 'NFR-006 | ' 'REQ-002 | '; do
	grep "^| $r" docs/prd/PRD-0001-layup.md | grep -q 'T-ax3r' && ok "PRD §12: $r" || no "the Test cell of $r does not name T-ax3r"
done
exit "$fail"
