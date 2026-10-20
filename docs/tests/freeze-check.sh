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
	: > "$tmp/k3_tests"
	git ls-files '*_test.go' | while IFS= read -r tf; do
		[ -f "$tf" ] || continue
		tdir=$(dirname -- "$tf")
		tr -d '\r' < "$tf" | awk -v d="$tdir" '
			!in_pkg && /^\/\/go:build / {
				t = $0; sub(/^\/\/go:build[[:space:]]+/, "", t); gsub(/[[:space:]]+$/, "", t)
				tag = t; next
			}
			!in_pkg && /^\/\/ \+build / && tag == "" {
				t = $0; sub(/^\/\/ \+build[[:space:]]+/, "", t); gsub(/[[:space:]]+$/, "", t)
				tag = t; next
			}
			/^package / { in_pkg = 1 }
			/^func Test[A-Za-z0-9_]+\(/ {
				nm = $0; sub(/^func /, "", nm); sub(/\(.*/, "", nm)
				if (nm != "TestMain")
					printf "%s\t%s\t%s\n", d, (tag == "" ? "-" : tag), nm
			}
		'
	done > "$tmp/k3_tests"
	: > "$tmp/k3_out"
	for mf in "runs/$task"/*.md; do
		[ -f "$mf" ] || continue
		[ -r "$mf" ] || { printf 'freeze-check: cannot read %s\n' "$mf"; exit 2; }
		tr -d '\r' < "$mf" > "$tmp/k3_md"
		awk -v md="$mf" '
			FILENAME == ARGV[1] {
				split($0, a, "\t")
				tn++; t_dir[tn] = a[1]; t_tag[tn] = a[2]; t_fn[tn] = a[3]
				next
			}
			function pkg_match(p, d,    pre) {
				sub(/^\.\//, "", p); sub(/\/$/, "", p)
				if (p == "") p = "."
				if (p == "...") return 1
				if (p ~ /\/\.\.\.$/) {
					pre = substr(p, 1, length(p) - 4)
					return (d == pre || substr(d, 1, length(pre) + 1) == pre "/")
				}
				return (d == p)
			}
			function tag_match(tags, ft) {
				if (ft == "-") return 1
				return (index("," tags ",", "," ft ",") > 0)
			}
			function check_cmd(cmd, ln,    rpos, after, qend, sel, pre, post, rest, nt, toks, i, tok, tags, np, pkgs, anchored, s, na, alts, a, j, k, ok) {
				if (!match(cmd, /(^|[[:space:]])go[[:space:]]+test([[:space:]]|$)/)) return
				if (!match(cmd, /[[:space:]]-run([[:space:]=]|$)/)) return
				rpos = RSTART
				pre = substr(cmd, 1, rpos - 1)
				after = substr(cmd, rpos + 1)
				sub(/^-run[[:space:]]*=?/, "", after)
				if (substr(after, 1, 1) == "\"") {
					printf "FAIL K3: %s:%d: quote the selector with single quotes\n", md, ln
					return
				}
				if (substr(after, 1, 1) != "\047") {
					printf "FAIL K3: %s:%d: quote the selector with single quotes\n", md, ln
					return
				}
				after = substr(after, 2)
				qend = index(after, "\047")
				if (qend == 0) {
					printf "FAIL K3: %s:%d: unsupported selector form\n", md, ln
					return
				}
				sel = substr(after, 1, qend - 1)
				post = substr(after, qend + 1)
				if (post != "" && substr(post, 1, 1) !~ /[[:space:]]/) {
					printf "FAIL K3: %s:%d: unsupported selector form\n", md, ln
					return
				}
				rest = pre " " post
				if (!match(rest, /(^|[[:space:]])go[[:space:]]+test([[:space:]]|$)/)) {
					printf "FAIL K3: %s:%d: unsupported selector form\n", md, ln
					return
				}
				rest = substr(rest, RSTART + RLENGTH)
				nt = split(rest, toks, /[[:space:]]+/)
				tags = ""; np = 0
				for (i = 1; i <= nt; i++) {
					tok = toks[i]
					if (tok == "") continue
					if (tok ~ /^-tags=/) { tags = substr(tok, 7); continue }
					if (tok == "-tags") { if (i == nt) { printf "FAIL K3: %s:%d: unsupported selector form\n", md, ln; return }; tags = toks[++i]; continue }
					if (tok ~ /^-(count|timeout)=/ || tok == "-v" || tok == "-short" || tok == "-race") continue
					if (tok == "-count" || tok == "-timeout") { if (i == nt) { printf "FAIL K3: %s:%d: unsupported selector form\n", md, ln; return }; i++; continue }
					if (tok == "." || tok == "..." || tok ~ /^\.\//) { pkgs[++np] = tok; continue }
					printf "FAIL K3: %s:%d: unsupported selector form\n", md, ln
					return
				}
				if (np == 0) {
					printf "FAIL K3: %s:%d: unsupported selector form\n", md, ln
					return
				}
				if (sel ~ /^[A-Za-z0-9_]+(\|[A-Za-z0-9_]+)*$/) {
					anchored = 0; s = sel
				} else if (sel ~ /^\^[A-Za-z0-9_]+\$(\|\^[A-Za-z0-9_]+\$)*$/) {
					anchored = 1; s = sel; gsub(/\^|\$/, "", s)
				} else if (sel ~ /^\^\([A-Za-z0-9_]+(\|[A-Za-z0-9_]+)*\)\$$/) {
					anchored = 1; s = substr(sel, 3, length(sel) - 4)
				} else {
					printf "FAIL K3: %s:%d: unsupported selector form: \047%s\047\n", md, ln, sel
					return
				}
				na = split(s, alts, "|")
				for (a = 1; a <= na; a++) {
					ok = 0
					for (j = 1; j <= tn; j++) {
						if (!tag_match(tags, t_tag[j])) continue
						for (k = 1; k <= np; k++) {
							if (pkg_match(pkgs[k], t_dir[j])) {
								if ((anchored && t_fn[j] == alts[a]) || (!anchored && index(t_fn[j], alts[a]) > 0)) {
									ok = 1; break
								}
							}
						}
						if (ok) break
					}
					if (!ok) {
						printf "FAIL K3: %s:%d: selector %s matches no test in the named packages\n", md, ln, alts[a]
					}
				}
			}
			{
				line = $0; has_span = 0
				while ((p1 = index(line, "`")) > 0) {
					line = substr(line, p1 + 1)
					p2 = index(line, "`")
					if (p2 == 0) break
					sp = substr(line, 1, p2 - 1)
					if (sp ~ /(^|[[:space:]])go[[:space:]]+test([[:space:]]|$)/ && sp ~ /[[:space:]]-run([[:space:]=]|$)/) {
						has_span = 1
						check_cmd(sp, FNR)
					}
					line = substr(line, p2 + 1)
				}
				if (!has_span && $0 !~ /`/ && $0 ~ /(^|[[:space:]])go[[:space:]]+test([[:space:]]|$)/ && $0 ~ /[[:space:]]-run([[:space:]=]|$)/) {
					check_cmd($0, FNR)
				}
			}
		' "$tmp/k3_tests" "$tmp/k3_md" >> "$tmp/k3_out"
	done
	if [ -s "$tmp/k3_out" ]; then
		cat "$tmp/k3_out"
		failed=1
	else
		printf 'ok K3: test selectors\n'
	fi
}

check_k4() {
	: > "$tmp/k4_skip"
	: > "$tmp/k4_fail"
	for mf in "runs/$task"/*.md; do
		[ -f "$mf" ] || continue
		[ -r "$mf" ] || { printf 'freeze-check: cannot read %s\n' "$mf"; exit 2; }
		tr -d '\r' < "$mf" | awk -v md="$mf" -v fskip="$tmp/k4_skip" -v ffail="$tmp/k4_fail" '
			function clean_val(v) {
				gsub(/^[[:space:]`]+|[[:space:]`]+$/, "", v)
				return v
			}
			function is_hex(v) {
				return (length(v) >= 12 && v ~ /^[0-9a-fA-F]+$/)
			}
			!/^\|/ {
				in_tbl = 0; c_before = 0; c_during = 0; c_after = 0; c_restore = 0
				next
			}
			!in_tbl {
				in_tbl = 1; c_before = 0; c_during = 0; c_after = 0; c_restore = 0
				nc = split($0, c, "|")
				for (i = 2; i <= nc; i++) {
					h = tolower(c[i]); gsub(/^[[:space:]]+|[[:space:]]+$/, "", h)
					if (h == "") continue
					if (h ~ /before[[:space:]]*=[[:space:]]*after/ || h ~ /^restore/ || (h ~ /restore/ && h !~ /after/)) c_restore = i
					else if (h ~ /before/) c_before = i
					else if (h ~ /during/) c_during = i
					else if (h ~ /after/) c_after = i
				}
				next
			}
			/^\|[-: |]+$/ { next }
			{
				nc = split($0, c, "|")
				id = clean_val(c[2])
				if (id !~ /^S?M[0-9]+[A-Za-z]?$/) next
				if (!(id in seen_id)) {
					seen_id[id] = 1
					id_order[++nid] = id
					id_line[id] = FNR
				}
				vb = (c_before > 0 && c_before <= nc ? clean_val(c[c_before]) : "")
				vd = (c_during > 0 && c_during <= nc ? clean_val(c[c_during]) : "")
				va = (c_after > 0 && c_after <= nc ? clean_val(c[c_after]) : "")
				vr = (c_restore > 0 && c_restore <= nc ? clean_val(c[c_restore]) : "")
				if (is_hex(vb) && is_hex(va) && vb == va) {
					id_ok[id] = 1
				} else if (is_hex(vb) && is_hex(vd) && vb != vd && (vr == "True" || vr == "equal" || vr == "restored")) {
					id_ok[id] = 1
				}
				for (i = 3; i <= nc; i++) {
					low = tolower(c[i])
					if (match(low, /(^|[^a-z0-9_])(not run|dropped|superseded)([^a-z0-9_]|$)/)) {
						w = substr(low, RSTART, RLENGTH)
						gsub(/^[^a-z]+|[^a-z]+$/, "", w)
						id_skip[id] = w
						id_skip_line[id] = FNR
					}
				}
			}
			END {
				for (k = 1; k <= nid; k++) {
					id = id_order[k]
					if (id_ok[id]) continue
					if (id in id_skip) {
						printf "skip K4: %s:%d: %s (%s)\n", md, id_skip_line[id], id, id_skip[id] >> fskip
					} else {
						printf "FAIL K4: %s:%d: %s has no restore evidence\n", md, id_line[id], id >> ffail
					}
				}
			}
		'
	done
	[ -s "$tmp/k4_skip" ] && cat "$tmp/k4_skip"
	if [ -s "$tmp/k4_fail" ]; then
		cat "$tmp/k4_fail"
		failed=1
	else
		printf 'ok K4: restore evidence of mutants\n'
	fi
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
	if [ "$fig_l" -gt "$max_l" ] || [ "$fig_f" -gt "$max_f" ]; then
		printf 'FAIL K5: -:0: projection: figure %d/%d, maximum %d/%d, margin %d/%d (the figure at landing is still measured at landing)\n' \
			"$fig_l" "$fig_f" "$max_l" "$max_f" "$mar_l" "$mar_f"
		failed=1
		return
	fi
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
