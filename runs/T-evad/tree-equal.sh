#!/bin/sh
# tree-equal.sh: condition 2 of O-146 (task T-evad, #97): the tree of the target's main after the merge of the
# target task equals the tree of the head of layup-setup-2, the second setup run, except the four paths of task
# T-a0rt, the target's own change of D7, and the target task's own file (O-147 = a).
#
# Usage: sh tree-equal.sh REPO MAIN SETUP2 [--task-file docs/tasks/T-xxxx.md]
#   REPO    a clone of the target
#   MAIN    main after the merge of the target task
#   SETUP2  the head of layup-setup-2
#   --task-file  the file of the target task, which its close-out adds under the target's gate step 8
# It lists each path whose entry (mode or object) differs between the two trees, or that one tree has and the
# other has not. Exit 0 and `tree-equal: PASS` when the only such paths are the four of T-a0rt and, when named,
# the task file as a new file of MAIN; else exit 1, with one `FAIL:` line for each other path, and for a named
# task file that MAIN does not add. Exit 2 on an input error.
set -u
R=${1:?usage: sh tree-equal.sh REPO MAIN SETUP2 [--task-file docs/tasks/T-xxxx.md]}
A=${2:?MAIN}
B=${3:?SETUP2}
shift 3
A0RT="docs/onboarding-for-engineers.md docs/tasks/T-a0rt.md runs/T-a0rt/evidence.md docs/tasks/completed.md"
TASK=
while [ $# -gt 0 ]; do
	case $1 in
	--task-file)
		[ $# -ge 2 ] || { echo "tree-equal: input error: --task-file needs a path" >&2; exit 2; }
		TASK=$2; shift 2 ;;
	*) echo "tree-equal: input error: an unknown argument: $1" >&2; exit 2 ;;
	esac
done
if [ -n "$TASK" ]; then
	case " $A0RT " in *" $TASK "*) echo "tree-equal: input error: the task file is one of the four paths of T-a0rt: $TASK" >&2; exit 2 ;; esac
	case $TASK in
	docs/tasks/T-[0-9a-z][0-9a-z][0-9a-z][0-9a-z].md) ;;
	*) echo "tree-equal: input error: the task file must be docs/tasks/T-xxxx.md: $TASK" >&2; exit 2 ;;
	esac
fi
gt() { git -C "$R" -c core.quotepath=off "$@"; }
a=$(gt rev-parse --verify -q "$A^{commit}") && b=$(gt rev-parse --verify -q "$B^{commit}") \
	|| { echo "tree-equal: input error: $A or $B is not a commit of $R" >&2; exit 2; }
echo "tree-equal: main $a against layup-setup-2 $b"
[ -z "$TASK" ] || echo "tree-equal: the task file is $TASK"
out=$(gt diff --no-renames --name-status "$b" "$a") || { echo "tree-equal: input error: git diff failed" >&2; exit 2; }
fail=0 seen=0
while IFS='	' read -r st p; do
	[ -n "$p" ] || continue
	if [ -n "$TASK" ] && [ "$p" = "$TASK" ]; then
		seen=1
		if [ "$st" = A ]; then echo "task file: $p ($st)"; else echo "FAIL: $p ($st): the task file is not a new file of main"; fail=$((fail + 1)); fi
		continue
	fi
	case " $A0RT " in
	*" $p "*) echo "T-a0rt: $p ($st)" ;;
	*) echo "FAIL: $p ($st)"; fail=$((fail + 1)) ;;
	esac
done <<EOF
$out
EOF
if [ -n "$TASK" ] && [ "$seen" -eq 0 ]; then echo "FAIL: the task file $TASK is not in main"; fail=$((fail + 1)); fi
if [ "$fail" -eq 0 ]; then
	if [ -n "$TASK" ]; then echo "tree-equal: PASS (only the paths of T-a0rt and the task file $TASK differ)"
	else echo "tree-equal: PASS (only the paths of T-a0rt differ)"; fi
	exit 0
fi
echo "tree-equal: FAIL ($fail findings: a path that is not one of the four of T-a0rt or the task file, or a task file that main does not add)"
exit 1
