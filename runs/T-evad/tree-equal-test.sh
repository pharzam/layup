#!/bin/sh
# Tests of tree-equal.sh (task T-evad, #97: condition 2 of O-146) on fixture repositories: the tree of main
# after the target task equals the tree of the second setup head, except the four paths of task T-a0rt.
#
# Usage: sh tree-equal-test.sh [TREE-EQUAL-SCRIPT]
# Exit 0 when each case behaves as written, else 1.
here=$(cd "$(dirname "$0")" && pwd)
TE=${1:-$here/tree-equal.sh}
T=$(mktemp -d) || exit 2
trap 'rm -rf "$T"' EXIT INT TERM
g() { git -C "$repo" -c user.name=fixture -c user.email=fixture@example.invalid -c commit.gpgsign=false \
	-c core.hooksPath=/dev/null -c maintenance.auto=false "$@"; }

# build NAME VARIANT: the second setup head (a file, a rule and the two files that T-a0rt changes), then main:
# the same tree, with the change of T-a0rt and the one change that VARIANT names.
build() {
	repo=$T/$1; v=$2
	mkdir -p "$repo/docs/tasks" && git init -q -b main "$repo"
	printf 'a rule\n' > "$repo/docs/rule.md"
	printf '# onboarding\n| gate | no |\n' > "$repo/docs/onboarding-for-engineers.md"
	printf '# completed\n' > "$repo/docs/tasks/completed.md"
	g add -A && g commit -q -m "chore: setup S14"
	g branch -q layup-setup-2
	# the change of T-a0rt: two files changed, two files added
	printf '# onboarding\n| gate | yes |\n' > "$repo/docs/onboarding-for-engineers.md"
	printf '# completed\n- T-a0rt\n' > "$repo/docs/tasks/completed.md"
	mkdir -p "$repo/runs/T-a0rt" && printf 'evidence\n' > "$repo/runs/T-a0rt/evidence.md"
	printf '# T-a0rt\n' > "$repo/docs/tasks/T-a0rt.md"
	case "$v" in
	other-file) printf 'a changed rule\n' > "$repo/docs/rule.md" ;;
	extra-file) printf 'more\n' > "$repo/docs/extra.md" ;;
	lost-file) git -C "$repo" rm -q docs/rule.md ;;
	mode) chmod +x "$repo/docs/rule.md" ;;
	esac
	g add -A && g commit -q -m "chore: T-xxxx main"
}

n=0; bad=0
# run NAME VARIANT WANT_EXIT WANT_MESSAGE
run() {
	n=$((n + 1)); build "$1" "$2"
	out=$(sh "$TE" "$repo" main layup-setup-2 2>&1); code=$?
	if [ "$code" = "$3" ] && printf '%s' "$out" | grep -Fq -- "$4"; then echo "ok   $n $1"; else
		bad=$((bad + 1)); echo "FAIL $n $1: exit $code (want $3); want \"$4\""; printf '%s\n' "$out" | sed 's/^/       | /' | head -10
	fi
}

run equal       none        0 'tree-equal: PASS'
run other-file  other-file  1 'FAIL: docs/rule.md (M)'
run extra-file  extra-file  1 'FAIL: docs/extra.md (A)'
run lost-file   lost-file   1 'FAIL: docs/rule.md (D)'
run mode        mode        1 'FAIL: docs/rule.md (M)'

echo "tree-equal-test: $((n - bad)) of $n cases behave as written"
[ "$bad" = 0 ]
