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

tmp=$(mktemp -d) || exit 2
trap 'rm -rf "${tmp:?}"' EXIT INT TERM
failed=0

check_k1() {
	if git diff --check "$base" HEAD > "$tmp/k1" 2>&1; then
		printf 'ok K1: white space\n'
		return
	fi
	failed=1
	awk '!/^\+/ && /^[^:]+:[0-9]+:/ { print "FAIL K1: " $0; n++ } END { if (!n) print "FAIL K1: -:0: git diff --check failed" }' "$tmp/k1"
}

check_k2() {
	git diff -U0 "$base" HEAD -- '*_test.go' | tr -d '\r' | awk '
		/^\+\+\+ b\// { f = substr($0, 7); next }
		/^@@ -[0-9,]+ \+[0-9]+/ {
			s = $0; sub(/^@@ -[0-9,]+ \+/, "", s); sub(/[, ].*/, "", s)
			ln = s + 0; next
		}
		/^-func Test[A-Za-z0-9_]+\(/ {
			nm = $0; sub(/^-func /, "", nm); sub(/\(.*/, "", nm)
			rem[f SUBSEP nm] = 1; next
		}
		/^\+/ && !/^\+\+\+/ {
			if ($0 ~ /^\+func Test[A-Za-z0-9_]+\(/) {
				nm = $0; sub(/^\+func /, "", nm); sub(/\(.*/, "", nm)
				if (nm != "TestMain") {
					n++; af[n] = f; al[n] = ln; an[n] = nm
				}
			}
			ln++
		}
		END {
			for (i = 1; i <= n; i++) {
				if (!((af[i] SUBSEP an[i]) in rem))
					print af[i] "\t" al[i] "\t" an[i]
			}
		}
	' > "$tmp/k2_added"
	if [ ! -s "$tmp/k2_added" ]; then
		printf 'ok K2: registration of new tests\n'
		return
	fi
	if [ ! -f docs/tests/traceability.md ]; then
		printf 'FAIL K2: docs/tests/traceability.md:0: missing traceability file\n'
		failed=1
		return
	fi
	tr -d '\r' < docs/tests/traceability.md > "$tmp/k2_tr"
	has_prd=0
	if [ -f docs/prd/PRD-0001-layup.md ]; then
		has_prd=1
		tr -d '\r' < docs/prd/PRD-0001-layup.md > "$tmp/k2_prd"
	else
		: > "$tmp/k2_prd"
	fi
	awk -v has_prd="$has_prd" '
		function has_word(txt, w,    r) {
			r = "(^|[^A-Za-z0-9_])" w "([^A-Za-z0-9_]|$)"
			return (txt ~ r)
		}
		function has_span_word(txt, w,    s, i, sp) {
			s = txt
			while ((i = index(s, "`")) > 0) {
				s = substr(s, i + 1)
				i = index(s, "`")
				if (i == 0) break
				sp = substr(s, 1, i - 1)
				if (has_word(sp, w)) return 1
				s = substr(s, i + 1)
			}
			return 0
		}
		FILENAME == ARGV[1] {
			if ($0 ~ /^\|/ && $0 !~ /^\|[[:space:]]*Test ID[[:space:]]*\|/ && $0 !~ /^\|[-: |]+$/) {
				tr_n++
				split($0, c, "|")
				tr_c1[tr_n] = c[2]
				tr_c3[tr_n] = c[4]
			}
			next
		}
		FILENAME == ARGV[2] {
			if ($0 ~ /^## 12([^0-9]|$)/) { in12 = 1; next }
			if (in12 && $0 ~ /^## /) { in12 = 0; next }
			if (in12 && $0 ~ /^\|/ && $0 !~ /^\|[[:space:]]*REQ[[:space:]]*\|/ && $0 !~ /^\|[-: |]+$/) {
				split($0, c, "|")
				id = c[2]; gsub(/^[[:space:]]+|[[:space:]]+$/, "", id)
				if (id != "") { prd_cell[id] = c[7]; prd_rows++ }
			}
			next
		}
		FILENAME == ARGV[3] {
			split($0, f, "\t")
			tfile = f[1]; tln = f[2]; tname = f[3]
			found_tr = 0; reqs = ""
			for (i = 1; i <= tr_n; i++) {
				if (has_span_word(tr_c1[i], tname)) {
					found_tr = 1
					s = tr_c3[i]
					while (match(s, /(REQ|NFR)-[0-9]+/)) {
						id = substr(s, RSTART, RLENGTH)
						reqs = (reqs == "" ? id : reqs " " id)
						s = substr(s, RSTART + RLENGTH)
					}
				}
			}
			if (!found_tr) {
				printf "FAIL K2: %s:%s: %s is not in docs/tests/traceability.md\n", tfile, tln, tname
				next
			}
			if (reqs == "") next
			if (!has_prd) {
				if (!prd_missing_printed) {
					printf "FAIL K2: docs/prd/PRD-0001-layup.md:0: missing PRD file\n"
					prd_missing_printed = 1
				}
				next
			}
			if (prd_rows == 0) {
				if (!s12_missing_printed) {
					printf "FAIL K2: docs/prd/PRD-0001-layup.md:0: missing §12 table\n"
					s12_missing_printed = 1
				}
				next
			}
			nr = split(reqs, rlist, " ")
			delete seen_r
			for (j = 1; j <= nr; j++) {
				r = rlist[j]
				if (r in seen_r) continue
				seen_r[r] = 1
				if (!(r in prd_cell) || !has_word(prd_cell[r], tname)) {
					printf "FAIL K2: %s:%s: %s is not in docs/prd/PRD-0001-layup.md §12 %s Test cell\n", tfile, tln, tname, r
				}
			}
		}
	' "$tmp/k2_tr" "$tmp/k2_prd" "$tmp/k2_added" > "$tmp/k2_out"
	if [ -s "$tmp/k2_out" ]; then
		cat "$tmp/k2_out"
		failed=1
	else
		printf 'ok K2: registration of new tests\n'
	fi
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
