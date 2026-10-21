#!/bin/sh
# Tests of coverage-check.sh (T-ym0k, #369, S4; conditions 1 and 3 of its plan
# review). It builds a small repository, doctors good records, and expects exit 1
# and a named line from each defect. COVERAGE_CHECK=PATH tests another script.
#   sh runs/T-ym0k/coverage-check-test.sh
here=$(cd "$(dirname "$0")" && pwd)
check=${COVERAGE_CHECK:-$here/coverage-check.sh}
t=$(mktemp -d) || exit 2
trap 'rm -rf "$t"' EXIT
r=$t/r
git init -q "$r" && cd "$r" && git config user.name t && git config user.email t@t || exit 2
mk() { mkdir -p "$(dirname "$1")"; awk -v n="$2" 'BEGIN { for (i = 1; i <= n; i++) print "line", i }' > "$1"; }
mk cmd/app/main.go 3; mk internal/a/a.go 2; mk internal/a/a_test.go 9; mk internal/a/testdata/t.go 7
mk internal/b/b.go 4; mk go.mod 1; mk internal/catalog/go/files/x.tmpl 2; mk internal/catalog/go/kinds.tsv 3; mk README.md 5
git add -A && git commit -q -m A && A=$(git rev-parse HEAD)
mk internal/c/c.go 6; mk internal/b/b.go 5; git add -A && git commit -q -m B && B=$(git rev-parse HEAD)
E=$(git commit-tree -m E "$(git hash-object -t tree /dev/null)") # no file at all
mk internal/d/d.go 2; git add internal/d/d.go # staged, in no commit: HEAD is B
# rec COMMIT: a good record of COMMIT, from git alone.
rec() {
	printf '## Release review\n\n### Coverage\n\n| Path | Lines | Read |\n| ---- | ----- | ---- |\n'
	{ git ls-tree -r --name-only "$1" | grep -E '^(cmd|internal)/.*\.go$' | grep -v -E '_test\.go$|/testdata/'
	  git ls-tree -r --name-only "$1" | grep -x go.mod
	  git ls-tree -r --name-only "$1" | grep '^internal/catalog/go/'; } | while read -r p; do
		printf '| `%s` | %s | `in full` |\n' "$p" "$(git show "$1:$p" | wc -l | tr -d ' ')"
	done
	printf '\nCoverage is complete.\n'
}
rec "$A" > "$t/good"; rec "$B" > "$t/goodB"
doctor() { awk -v mode="$1" -v last="$(grep -c '^| `' "$2")" '
	/^\| `/ { n++
		if (mode == "drop" && n == 2) next
		if (mode == "count" && n == 2) sub(/\| 2 \|/, "| 3 |")
		if (mode == "read" && n == 2) sub(/`in full`/, "`partly`")
		if (mode == "cells" && n == 2) $0 = "| `internal/a/a.go` | 2 |"
		if (mode == "empty") next
		if (mode == "swap" && n == 2) { held = $0; next }
		print
		if (mode == "swap" && n == 3) print held
		if (mode == "dup" && n == 1) print
		if (mode == "extra" && n == last) print "| `internal/z/z.go` | 1 | `in full` |"
		next }
	mode == "nohead" && /^### Coverage/ { next }
	{ print }' "$2"; }
fails=0
run() { # name, want exit, text of the named line ("" for none), commit, record
	out=$(sh "$check" "$4" "$5" 2>&1); rc=$?
	if [ "$rc" -ne "$2" ]; then echo "FAIL $1: exit $rc, wanted $2"; fails=$((fails + 1)); return; fi
	case $out in *"$3"*) echo "ok   $1" ;; *) echo "FAIL $1: no line with: $3"; fails=$((fails + 1)) ;; esac
}
for m in drop count extra swap read dup cells empty nohead; do doctor $m "$t/good" > "$t/$m"; done
run "good record of A" 0 "" "$A" "$t/good"
run "one row dropped" 1 "missing row: internal/a/a.go" "$A" "$t/drop"
run "one count off by one" 1 "count: internal/a/a.go" "$A" "$t/count"
run "one extra path" 1 "extra row: internal/z/z.go" "$A" "$t/extra"
run "two rows swapped" 1 "out of order" "$A" "$t/swap"
run "a Read that is not in full" 1 "not \`in full\`: internal/a/a.go" "$A" "$t/read"
run "a duplicate row" 1 "duplicate row: cmd/app/main.go" "$A" "$t/dup"
run "a row of two cells" 1 "unparsed row" "$A" "$t/cells"
run "an empty record table" 1 "empty record table" "$A" "$t/empty"
run "no Coverage heading" 1 "empty record table" "$A" "$t/nohead"
run "an empty release list and an empty table" 1 "empty release list" "$E" "$t/empty"
run "an empty release list and a table" 1 "empty release list" "$E" "$t/good"
# The index and the tree hold a file that B lacks, and HEAD is B: COMMIT rules.
run "good record of B (HEAD, a staged file out of B)" 0 "" "$B" "$t/goodB"
run "A's record, HEAD at B" 0 "" "$A" "$t/good"
run "B's record for the commit A" 1 "extra row: internal/c/c.go" "$A" "$t/goodB"
run "A's record for the commit B" 1 "missing row: internal/c/c.go" "$B" "$t/good"
run "a name that is not a commit" 2 "not a commit" "nosuchcommit" "$t/good"
[ "$fails" -eq 0 ] && echo "all ok" || echo "$fails failed"
[ "$fails" -eq 0 ]
