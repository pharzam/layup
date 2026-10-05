#!/bin/sh
# tree-equal.sh: condition 2 of O-146 (task T-evad, #97): the tree of the target's main after the merge of the
# target task equals the tree of the head of layup-setup-2, the second setup run, except the four paths of task
# T-a0rt, the target's own change of D7.
#
# Usage: sh tree-equal.sh REPO MAIN SETUP2
#   REPO    a clone of the target
#   MAIN    main after the merge of the target task
#   SETUP2  the head of layup-setup-2
# It lists each path whose entry (mode or object) differs between the two trees, or that one tree has and the
# other has not. Exit 0 and `tree-equal: PASS` when the only such paths are the four of T-a0rt; else exit 1,
# with one `FAIL:` line for each other path.
set -u
R=${1:?usage: sh tree-equal.sh REPO MAIN SETUP2}
A=${2:?MAIN}
B=${3:?SETUP2}
A0RT="docs/onboarding-for-engineers.md docs/tasks/T-a0rt.md runs/T-a0rt/evidence.md docs/tasks/completed.md"
gt() { git -C "$R" -c core.quotepath=off "$@"; }
a=$(gt rev-parse --verify -q "$A^{commit}") && b=$(gt rev-parse --verify -q "$B^{commit}") \
	|| { echo "tree-equal: input error: $A or $B is not a commit of $R" >&2; exit 2; }
echo "tree-equal: main $a against layup-setup-2 $b"
out=$(gt diff --no-renames --name-status "$b" "$a") || { echo "tree-equal: input error: git diff failed" >&2; exit 2; }
fail=0
while IFS='	' read -r st p; do
	[ -n "$p" ] || continue
	case " $A0RT " in
	*" $p "*) echo "T-a0rt: $p ($st)" ;;
	*) echo "FAIL: $p ($st)"; fail=$((fail + 1)) ;;
	esac
done <<EOF
$out
EOF
if [ "$fail" -eq 0 ]; then echo "tree-equal: PASS (only the paths of T-a0rt differ)"; exit 0; fi
echo "tree-equal: FAIL ($fail paths differ that are not the four of T-a0rt)"
exit 1
