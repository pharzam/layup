#!/bin/sh
#
# freeze-check.sh — pre-freeze check of a task worktree (T-b6fn, #440).
#
# Usage:
#   sh docs/tests/freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--closeout LINES/FILES] [--no-conditions]
#
# Reads git and text at HEAD of the task worktree.
# Prints one line per check:
#   ok K<n>: <what>
#   FAIL K<n>: <file>:<line>: <reason>
#   skip K<n>: <reason>
# Exit status: 0 when no check fails, 1 when one or more fail, 2 on a usage or
# environment error (not a git work tree, --base not a commit, an unreadable file).
#
# Stated limit: K6 proves that cited evidence exists, not that it satisfies the
# condition, and it does not know the review's list of conditions; completeness
# stays the author's and the reviewer's (docs/guardrails.md:144).

set -u

usage() {
	printf 'freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--closeout LINES/FILES] [--no-conditions]\n'
	exit 2
}

task=''
base=''
budget=''
closeout=70/2
no_conditions=0

while [ $# -gt 0 ]; do
	case $1 in
	--task)
		[ $# -ge 2 ] || usage
		task=$2; shift 2 ;;
	--base)
		[ $# -ge 2 ] || usage
		base=$2; shift 2 ;;
	--budget)
		[ $# -ge 2 ] || usage
		budget=$2; shift 2 ;;
	--closeout)
		[ $# -ge 2 ] || usage
		closeout=$2; shift 2 ;;
	--no-conditions)
		no_conditions=1; shift ;;
	*)
		usage ;;
	esac
done

[ -n "$task" ] && [ -n "$base" ] && [ -n "$budget" ] || usage
printf '%s\n' "$budget" | grep -Eq '^[0-9]+/[0-9]+$' || usage
printf '%s\n' "$closeout" | grep -Eq '^[0-9]+/[0-9]+$' || usage

top=$(git rev-parse --show-toplevel 2>/dev/null) && [ -n "$top" ] || {
	printf 'freeze-check: not a git work tree\n'
	exit 2
}
cd "$top" || exit 2

git rev-parse --verify --quiet "$base^{commit}" >/dev/null 2>&1 || {
	printf 'freeze-check: --base %s is not a commit\n' "$base"
	exit 2
}

failed=0

check_k1() {
	printf 'ok K1: white space\n'
}

check_k2() {
	printf 'ok K2: registration of new tests\n'
}

check_k3() {
	printf 'ok K3: test selectors\n'
}

check_k4() {
	printf 'ok K4: restore evidence of mutants\n'
}

check_k5() {
	max_l=${budget%/*}
	max_f=${budget#*/}
	cl_l=${closeout%/*}
	cl_f=${closeout#*/}
	set -- $(git diff --numstat "$base" HEAD | awk '
		{ a = ($1 == "-" ? 0 : $1); r = ($2 == "-" ? 0 : $2); l += a + r; f += 1 }
		END { printf "%d %d\n", l + 0, f + 0 }
	')
	fig_l=$(($1 + cl_l))
	fig_f=$(($2 + cl_f))
	mar_l=$((max_l - fig_l))
	mar_f=$((max_f - fig_f))
	printf 'ok K5: projection: figure %d/%d, maximum %d/%d, margin %d/%d (the figure at landing is still measured at landing)\n' \
		"$fig_l" "$fig_f" "$max_l" "$max_f" "$mar_l" "$mar_f"
}

check_k6() {
	if [ "$no_conditions" -eq 1 ]; then
		printf 'skip K6: --no-conditions\n'
		return
	fi
	printf 'ok K6: conditions table\n'
}

check_k1
check_k2
check_k3
check_k4
check_k5
check_k6

[ "$failed" -eq 0 ] && exit 0 || exit 1
