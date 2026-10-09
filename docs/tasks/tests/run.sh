#!/bin/sh
#
# run.sh CASE — run docs/tasks/task-state.sh on one fixture case and compare.
#
# The LINTER of the suite `task-state` of docs/tests/run-discipline-tests.sh,
# which runs `sh run.sh CASE_DIR` and asserts the exit code only. task-state.sh
# is a command, not a linter, so this script turns its output into that
# contract:
#
#   good-*  exit 0 when task-state.sh exits 0 and its output equals EXPECT
#   bad-*   exit 1 when task-state.sh exits 1 and its stderr equals EXPECT
#   else    exit 2, with the difference on stderr
#
# A mismatch is 2, never 1: an exit of 1 would let a bad case pass for the
# wrong reason (docs/guardrails.md, a check that cannot fail).
#
# A case holds only the files it changes. The input is base/ with the case's
# files laid over it; a file named ABSENT lists the inputs to remove. The
# inputs carry no extension (backlog, completed, plan), so link-lint and the
# check `adapted` do not read them as documents.

set -u

[ $# -eq 1 ] && [ -d "$1" ] || { printf 'run: usage: run.sh CASE_DIR\n' >&2; exit 2; }
case_dir=${1%/}
name=$(basename "$case_dir")
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
[ -f "$case_dir/EXPECT" ] || { printf 'run: %s has no EXPECT\n' "$name" >&2; exit 2; }

tmp=$(mktemp -d) || exit 2
trap 'rm -rf "$tmp"' EXIT
mkdir "$tmp/in"
cp -R "$here/base/." "$tmp/in/" && cp -R "$case_dir/." "$tmp/in/" || exit 2
rm -f "$tmp/in/EXPECT"
if [ -f "$tmp/in/ABSENT" ]; then
	while IFS= read -r f; do
		[ -n "$f" ] && rm -rf "$tmp/in/$f"
	done < "$tmp/in/ABSENT"
	rm -f "$tmp/in/ABSENT"
fi

sh "$here/../task-state.sh" --read "$tmp/in" > "$tmp/out" 2> "$tmp/err"
got=$?

case $name in
good*)
	if [ "$got" -eq 0 ] && cmp -s "$tmp/out" "$case_dir/EXPECT"; then exit 0; fi
	printf 'run: %s: exit %s; the output and EXPECT:\n' "$name" "$got" >&2
	diff "$case_dir/EXPECT" "$tmp/out" >&2
	cat "$tmp/err" >&2
	exit 2 ;;
bad*)
	if [ "$got" -eq 1 ] && cmp -s "$tmp/err" "$case_dir/EXPECT"; then exit 1; fi
	printf 'run: %s: exit %s; the stderr and EXPECT:\n' "$name" "$got" >&2
	diff "$case_dir/EXPECT" "$tmp/err" >&2
	exit 2 ;;
*)
	printf 'run: %s is neither good* nor bad*\n' "$name" >&2
	exit 2 ;;
esac
