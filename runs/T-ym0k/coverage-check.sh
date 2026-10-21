#!/bin/sh
# The coverage comparison of the release review of M2e (task T-ym0k, #369; D5 of
# its plan, with conditions 1 and 3 of its plan review). It compares the table
# `### Coverage` of a review record with the release at COMMIT:
#
#   sh coverage-check.sh COMMIT RECORD
#   sh coverage-check.sh --list COMMIT      (prints the release: path, TAB, lines)
#
# The release is the non-test Go files of cmd/ and internal/ out of testdata/
# (filtered by name), go.mod, and git ls-files internal/catalog/go, in that
# order, each with the wc -l of git show COMMIT:path. Both halves come from
# COMMIT: git ls-files reads a temporary index made by git read-tree COMMIT, so
# the index and the tree of the clone in which it runs are not read.
# It exits 1, with one FAIL line for each, when the release has no file, the
# table has no parsed row, a row is not `| `path` | N | Read |` (Read with or
# without backticks), a row is
# missing, extra or duplicate, a count differs, the rows are in another order,
# or a Read is not `in full`; exit 2 when COMMIT is not a commit.
# Stated limit: it proves that the record names each file with the right count,
# not that the reviewer read it.

usage() { echo "usage: sh coverage-check.sh COMMIT RECORD | --list COMMIT" >&2; exit 2; }
list=0; [ "${1:-}" = "--list" ] && { list=1; shift; }
commit=${1:-}; rec=${2:-}
[ -n "$commit" ] && { [ "$list" -eq 1 ] || [ -f "$rec" ]; } || usage
id=$(git rev-parse --verify --quiet "$commit^{commit}") || { echo "FAIL: not a commit: $commit" >&2; exit 2; }
d=$(mktemp -d) || exit 2
trap 'rm -rf "$d"' EXIT
export GIT_INDEX_FILE=$d/index
git read-tree "$id" || exit 2
{ git ls-files -- 'cmd/*.go' 'internal/*.go' | grep -v -E '_test\.go$|/testdata/'
  git ls-files -- go.mod
  git ls-files -- internal/catalog/go; } | while IFS= read -r p; do
	printf '%s\t%s\n' "$p" "$(git show "$id:$p" | wc -l | tr -d ' ')"
done > "$d/exp"
if [ "$list" -eq 1 ]; then cat "$d/exp"; exit 0; fi
fail=0
bad() { echo "FAIL: $*"; fail=1; }
[ -s "$d/exp" ] || bad "empty release list at $id"
awk '
	/^### Coverage[ \t]*$/ { sec = 1; next }
	sec && /^#/ { sec = 0 }
	sec && /^\|/ {
		if ($0 ~ /^\| *-+ *\|/ || $0 ~ /^\| *Path *\|/) next
		if ($0 !~ /^\| `[^`]+` \| [0-9]+ \| +[^| ]([^|]*[^| ])? +\|[ \t]*$/) { print "U\t" $0; next }
		split($0, c, "|")
		for (i = 2; i <= 4; i++) gsub(/^[ \t]*`?|`?[ \t]*$/, "", c[i])
		print "R\t" c[2] "\t" c[3] "\t" c[4] }' "$rec" > "$d/rows"
grep '^U' "$d/rows" | while IFS= read -r l; do echo "FAIL: unparsed row: ${l#U?}"; done > "$d/u"
[ -s "$d/u" ] && { cat "$d/u"; fail=1; }
grep '^R' "$d/rows" | cut -f2- > "$d/got"
[ -s "$d/got" ] || bad "empty record table in $rec"
if [ -s "$d/exp" ] && [ -s "$d/got" ]; then
	awk -F '\t' '
		FNR == NR { ep[$1] = $2; epath[FNR] = $1; en = FNR; next }
		{ if ($1 in gp) { print "FAIL: duplicate row: " $1; bad = 1; next }
		  gp[$1] = $2; gpath[++gn] = $1
		  if (!($1 in ep)) { print "FAIL: extra row: " $1; bad = 1; next }
		  if ($2 != ep[$1]) { print "FAIL: count: " $1 ": record " $2 ", commit " ep[$1]; bad = 1 }
		  if ($3 != "in full") { print "FAIL: Read is not `in full`: " $1; bad = 1 } }
		END { for (i = 1; i <= en; i++) if (!(epath[i] in gp)) { print "FAIL: missing row: " epath[i]; bad = 1 }
		  if (!bad) for (i = 1; i <= en; i++) if (gpath[i] != epath[i]) {
			print "FAIL: out of order: row " i " is " gpath[i] ", the commit has " epath[i]; bad = 1; break }
		  exit bad }' "$d/exp" "$d/got" || fail=1
fi
[ "$fail" -eq 0 ] && echo "ok: $(wc -l < "$d/got" | tr -d ' ') rows of $rec are the release at $id, each in full"
exit "$fail"
