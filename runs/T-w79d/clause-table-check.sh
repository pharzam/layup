#!/bin/sh
#
# clause-table-check.sh — run evidence of task T-w79d (#46). NOT a gate: no hook
# or CI job runs it, so it counts as no check (Invariant 5). It is the test of
# the task's demo, run by hand, with its runs recorded in test-runs.txt beside it.
#
# It checks the PRESENCE and the SHAPE of ADR-0012's clause table and of the
# parts of its Decision section, outside fenced code blocks, and that the cited
# source F-0005 matches its hash. It does not check MEANING — whether a basis
# supports its relation, or a decision is right. That is the review rounds' work.
#
# Usage:
#   sh runs/T-w79d/clause-table-check.sh [ROOT]              check ROOT (default .)
#   sh runs/T-w79d/clause-table-check.sh --self-test [ROOT]  make each predicate fail once
#
# Exit: 0 = every predicate holds; 1 = at least one fails, one line each:
#   clause-table: FAIL P<n> <reason>

set -u

self_test=0
if [ "${1:-}" = "--self-test" ]; then self_test=1; shift; fi
ROOT=${1:-.}

HEADER='| Clause | F-0005 lines | Relation | Basis | Destination |'

check() {
	root=$1
	failed=0
	fail() { printf 'clause-table: FAIL %s %s\n' "$1" "$2"; failed=1; }

	# P1 — one ADR-0012 file, with the clause table header.
	set -- "$root"/docs/adr/0012-*.md
	if [ "$#" != 1 ] || [ ! -f "$1" ]; then
		fail P1 "expected exactly one docs/adr/0012-*.md"
		return 1
	fi
	tmp=$(mktemp -d "${TMPDIR:-/tmp}/clause-table.XXXXXX")
	# Every predicate reads the ADR without its fenced code blocks and its HTML
	# comment blocks, because neither is Markdown structure (CommonMark). A fence
	# opens on a line with at most three leading spaces and then at least three
	# backticks or three tildes, and closes on a line with at most three leading
	# spaces, the same character at least as many times, and then only spaces. A
	# comment block opens on a line with at most three leading spaces and then
	# `<!--`, and closes on the first line that holds `-->`. An unclosed block of
	# either kind runs to the end of the file.
	awk '
		function fence(s,   i, c, n) {
			i = 1
			while (i <= 4 && substr(s, i, 1) == " ") i++
			if (i > 4) return ""
			c = substr(s, i, 1)
			if (c != "`" && c != "~") return ""
			n = 0
			while (substr(s, i + n, 1) == c) n++
			if (n < 3) return ""
			rest = substr(s, i + n)
			return c n
		}
		function lead(s,   i) {
			i = 1
			while (i <= 4 && substr(s, i, 1) == " ") i++
			return i > 4 ? "" : substr(s, i)
		}
		{
			if (comment) { if (index($0, "-->")) comment = 0; next }
			f = fence($0)
			if (!inside) {
				if (f != "") { inside = 1; fc = substr(f, 1, 1); fn = substr(f, 2) + 0; next }
				l = lead($0)
				if (substr(l, 1, 4) == "<!--") { if (!index(substr(l, 5), "-->")) comment = 1; next }
				print
				next
			}
			if (f != "" && substr(f, 1, 1) == fc && substr(f, 2) + 0 >= fn && rest ~ /^ *$/) inside = 0
		}
	' "$1" > "$tmp/adr"
	adr=$tmp/adr
	if ! grep -Fqx -- "$HEADER" "$adr"; then
		fail P1 "no clause table header: $HEADER"
		rm -rf "$tmp"
		return 1
	fi
	# The line after the header must be a delimiter row (GitHub Flavored Markdown):
	# without up to three leading spaces and an optional outer pipe on each side, it
	# splits on "|" into five cells, each of the form :?-+:? with spaces around.
	if ! awk -v h="$HEADER" '
		$0 == h { getline d; found = 1
			k = 0; while (k < 3 && substr(d, 1, 1) == " ") { d = substr(d, 2); k++ }
			sub(/^\|/, "", d); sub(/\|[ \t]*$/, "", d)
			n = split(d, c, "|"); if (n != 5) exit 1
			for (i = 1; i <= n; i++) if (c[i] !~ /^[ \t]*:?-+:?[ \t]*$/) exit 1
			exit 0 }
		END { if (!found) exit 1 }' "$adr"; then
		fail P1 "no delimiter row under the clause table header"
		rm -rf "$tmp"
		return 1
	fi

	# The table's rows (GitHub Flavored Markdown): the lines after the delimiter row
	# up to the first line that is blank or starts another block — after up to three
	# leading spaces: "#" (a heading), ">" (a block quote), "-", "*" or "+" and a
	# space, or digits and "." or ")" and a space (a list item), a thematic break
	# (only three or more of one of "-", "*", "_", spaces allowed), three or more
	# backticks or tildes (a fence), or "<" (an HTML block). Every other line is a
	# row, read as if it had a leading pipe.
	awk -v h="$HEADER" '
		$0 == h { on = 1; getline; next }
		!on { next }
		{ l = $0; k = 0; while (k < 3 && substr(l, 1, 1) == " ") { l = substr(l, 2); k++ } }
		l ~ /^[ \t]*$/ { exit }
		l ~ /^(#|>|[-*+] |[0-9]+[.)] |```|~~~|<)/ { exit }
		{ t = l; gsub(/[ \t]/, "", t) }
		length(t) >= 3 && (t ~ /^-+$/ || t ~ /^\*+$/ || t ~ /^_+$/) { exit }
		{ if (substr(l, 1, 1) != "|") l = "|" l; print l }
	' "$adr" > "$tmp/rows"
	# One tab-separated line per row: id, relation, basis, destination (cells trimmed,
	# backticks dropped). An escaped pipe "\|" is part of a cell, not a separator (GFM).
	awk '{
		l = $0; gsub(/\\\|/, "\001", l); split(l, f, "|")
		for (i = 2; i <= 6; i++) { c = f[i]; gsub(/\001/, "|", c); gsub(/`/, "", c); sub(/^[ \t]+/, "", c); sub(/[ \t]+$/, "", c); v[i] = c }
		printf "%s\t%s\t%s\t%s\n", v[2], v[4], v[5], v[6]
	}' "$tmp/rows" > "$tmp/cells"

	# P2 — each ID C00..C24 opens exactly one row; no other ID does.
	i=0
	while [ "$i" -le 24 ]; do
		id=$(printf 'C%02d' "$i")
		n=$(awk -F'\t' -v id="$id" '$1 == id' "$tmp/cells" | grep -c .)
		[ "$n" = 1 ] || fail P2 "$id opens $n rows, expected 1"
		i=$((i + 1))
	done
	awk -F'\t' '$1 !~ /^C(0[0-9]|1[0-9]|2[0-4])$/ { print $1 }' "$tmp/cells" | while IFS= read -r x; do
		printf 'clause-table: FAIL P2 a row opens with "%s", not an ID C00..C24\n' "$x"
	done > "$tmp/p2"
	[ -s "$tmp/p2" ] && { cat "$tmp/p2"; failed=1; }

	# P3 — each Relation is one of the five values.
	awk -F'\t' '$2 != "out of scope" && $2 != "conflicts" && $2 != "no evidence" && $2 != "extends" && $2 != "consistent" {
		printf "clause-table: FAIL P3 %s has relation \"%s\"\n", $1, $2 }' "$tmp/cells" > "$tmp/p3"
	[ -s "$tmp/p3" ] && { cat "$tmp/p3"; failed=1; }

	# P4 — each Basis holds at least one citation.
	cite='F-000[0-9]#[0-9]+|F-0005 L[0-9]+|[A-Za-z0-9_./-]+\.(md|sh)|ADR-[0-9][0-9][0-9][0-9]|(^|[^A-Za-z0-9-])R[0-9]+|O-[0-9]+|Inv-[0-9]'
	awk -F'\t' '{ print $1 "\t" $3 }' "$tmp/cells" | while IFS="$(printf '\t')" read -r id basis; do
		printf '%s\n' "$basis" | grep -Eq -- "$cite" || printf 'clause-table: FAIL P4 %s has no citation in its basis\n' "$id"
	done > "$tmp/p4"
	[ -s "$tmp/p4" ] && { cat "$tmp/p4"; failed=1; }

	# P5 — each Destination is exactly one allowed value; D, O and X values are defined.
	awk -F'\t' '{ print $1 "\t" $4 }' "$tmp/cells" | while IFS="$(printf '\t')" read -r id dest; do
		if [ "$id" = C00 ] && [ "$dest" != "not policy" ]; then
			printf 'clause-table: FAIL P5 C00 goes to "%s", but only "not policy" is allowed for C00\n' "$dest"
			continue
		fi
		case $dest in
			"not policy") [ "$id" = C00 ] || printf 'clause-table: FAIL P5 %s uses "not policy", which only C00 may use\n' "$id" ;;
			rejected) : ;;
			*)
				if printf '%s\n' "$dest" | grep -Eqx 'D[0-9]+|O-[0-9]+|X[0-9]+'; then
					grep -Eq "^- \*\*$dest\.\*\*" "$adr" || printf 'clause-table: FAIL P5 %s goes to %s, which the ADR does not define\n' "$id" "$dest"
				else
					printf 'clause-table: FAIL P5 %s has destination "%s", not one of D<n>, O-<n>, X<n>, rejected, not policy\n' "$id" "$dest"
				fi ;;
		esac
	done > "$tmp/p5"
	[ -s "$tmp/p5" ] && { cat "$tmp/p5"; failed=1; }

	# P6 — the cited source is present and matches its hash: one F-0005 record, a
	# hash line for the source, and check facts green (which, since #47, fails on a
	# changed byte, a removed hash line and a removed index row).
	set -- "$root"/docs/facts/F-0005-*.md
	if [ "$#" != 1 ] || [ ! -f "$1" ]; then fail P6 "expected exactly one docs/facts/F-0005-*.md"; fi
	grep -Eq '^[0-9a-f]{64}  docs/facts/operator-routing-policy\.md$' "$root/docs/setup/facts.sha256" 2>/dev/null \
		|| fail P6 "docs/setup/facts.sha256 has no line for docs/facts/operator-routing-policy.md"
	if ! sh "$root/docs/setup/setup-check.sh" --only facts "$root" > "$tmp/p6" 2>&1; then
		fail P6 "check facts fails: $(grep FAIL "$tmp/p6" | head -1)"
	fi

	# P7 — the four subsections inside "## Decision", each with an entry; each defined ID
	# is used or says (no clause).
	for sub in 'Proposed decisions' 'Rejected options' 'Open Operator decisions' 'Deferred items'; do
		n=$(awk -v s="### $sub" '
			/^## / { sec = $0 }
			$0 == s && sec == "## Decision" { on = 1; next }
			on && /^#/ { exit }
			on && /^- / { n++ }
			END { print n + 0 }' "$adr")
		[ "$n" -ge 1 ] || fail P7 "no entry under \"### $sub\" inside \"## Decision\""
	done
	grep -E '^- \*\*(D[0-9]+|O-[0-9]+|X[0-9]+)\.\*\*' "$adr" | while IFS= read -r line; do
		def=$(printf '%s\n' "$line" | sed -E 's/^- \*\*([^.]+)\.\*\*.*/\1/')
		awk -F'\t' -v d="$def" '$4 == d { f = 1 } END { exit !f }' "$tmp/cells" \
			|| printf '%s\n' "$line" | grep -Fq '(no clause)' \
			|| printf 'clause-table: FAIL P7 %s is defined but no clause goes to it and it does not say (no clause)\n' "$def"
	done > "$tmp/p7"
	[ -s "$tmp/p7" ] && { cat "$tmp/p7"; failed=1; }

	rm -rf "$tmp"
	return "$failed"
}

if [ "$self_test" = 0 ]; then
	if check "$ROOT"; then echo "clause-table: OK"; exit 0; else exit 1; fi
fi

# --- self-test: mutations that must fail, each with its own P<n>, and valid
# controls that must pass ---
base=$(mktemp -d "${TMPDIR:-/tmp}/clause-table-self.XXXXXX")
trap 'rm -rf "$base"' EXIT INT TERM
# Copy the tracked files one by one (POSIX: no `xargs -0`, no `tar`). The list is
# read from a file, not a pipe, so a failed copy stops the self-test here.
(cd "$ROOT" && git -c core.quotePath=false ls-files) > "$base.list"
while IFS= read -r f; do
	mkdir -p "$base/$(dirname "$f")" && cp -p "$ROOT/$f" "$base/$f" \
		|| { echo "clause-table self-test: FAIL cannot copy $f"; rm -f "$base.list"; exit 1; }
done < "$base.list"
rm -f "$base.list"
check "$base" > "$base.out" 2>&1 || { echo "clause-table self-test: FAIL the unmutated tree does not pass:"; cat "$base.out"; rm -f "$base.out"; exit 1; }
rm -f "$base.out"
adr_rel=$(cd "$base" && ls docs/adr/0012-*.md)

pass=0 bad=0
# mutate <case> <expected line prefix> <shell code run in a fresh copy of the valid tree>
mutate() {
	case_id=$1 want=$2 code=$3
	m=$(mktemp -d "${TMPDIR:-/tmp}/clause-table-mut.XXXXXX")
	cp -R "$base/." "$m/"
	(cd "$m" && ADR=$adr_rel sh -c "$code")
	out=$(check "$m" 2>&1); rc=$?
	if [ "$rc" = 1 ] && printf '%s\n' "$out" | grep -Fq -- "$want"; then
		printf 'ok    %-3s %s\n' "$case_id" "$(printf '%s\n' "$out" | grep -F -- "$want" | head -1)"
		pass=$((pass + 1))
	else
		printf 'FAIL  %-3s want "%s", got exit %s:\n%s\n' "$case_id" "$want" "$rc" "$out"
		bad=$((bad + 1))
	fi
	rm -rf "$m"
}
# row <ID> <field 4|5|6> <value>: set one cell of a clause row (fields as awk -F'|' counts them)
row='awk -F"|" -v OFS="|" -v id="$1" -v f="$2" -v v="$3" "\$2 ~ (\" \" id \" \") { \$f = \" \" v \" \" } { print }" "$ADR" > x && mv x "$ADR"'
H='| Clause | F-0005 lines | Relation | Basis | Destination |'
# control <case> <what> <shell code>: a valid change; the check must still pass, with no output
control() {
	case_id=$1 what=$2 code=$3
	m=$(mktemp -d "${TMPDIR:-/tmp}/clause-table-mut.XXXXXX")
	cp -R "$base/." "$m/"
	(cd "$m" && ADR=$adr_rel sh -c "$code")
	if cmp -s "$base/$adr_rel" "$m/$adr_rel"; then
		printf 'FAIL  %-3s control: %s, but the change left the ADR unchanged\n' "$case_id" "$what"
		bad=$((bad + 1)); rm -rf "$m"; return
	fi
	out=$(check "$m" 2>&1); rc=$?
	if [ "$rc" = 0 ] && [ -z "$out" ]; then
		printf 'ok    %-3s control: %s -> clause-table: OK\n' "$case_id" "$what"
		pass=$((pass + 1))
	else
		printf 'FAIL  %-3s control: %s, want OK, got exit %s:\n%s\n' "$case_id" "$what" "$rc" "$out"
		bad=$((bad + 1))
	fi
	rm -rf "$m"
}
mutate 1a 'clause-table: FAIL P1 expected exactly one' 'rm "$ADR"'
mutate 1b 'clause-table: FAIL P1 expected exactly one' 'cp "$ADR" docs/adr/0012-copy.md'
mutate 1c 'clause-table: FAIL P1 no clause table header' "grep -Fvx -- '$H' \"\$ADR\" > x && mv x \"\$ADR\""
mutate 2a 'clause-table: FAIL P2 C07 opens 0 rows' 'grep -v "^| C07 |" "$ADR" > x && mv x "$ADR"'
mutate 2b 'clause-table: FAIL P2 C07 opens 2 rows' 'awk "{ print } /^\\| C07 \\|/ { print }" "$ADR" > x && mv x "$ADR"'
mutate 2c 'clause-table: FAIL P2 a row opens with "C25"' 'awk "{ print } /^\\| C24 \\|/ { r = \$0; sub(/C24/, \"C25\", r); print r }" "$ADR" > x && mv x "$ADR"'
mutate 3  'clause-table: FAIL P3 C01 has relation "maybe"' "set -- C01 4 maybe; $row"
mutate 4a 'clause-table: FAIL P4 C01 has no citation' "set -- C01 5 ''; $row"
mutate 4b 'clause-table: FAIL P4 C01 has no citation' "set -- C01 5 'see above'; $row"
mutate 5a 'clause-table: FAIL P5 C01 has destination "D1 O-16"' "set -- C01 6 'D1 O-16'; $row"
mutate 5b 'clause-table: FAIL P5 C01 goes to D99, which the ADR does not define' "set -- C01 6 D99; $row"
mutate 5c 'clause-table: FAIL P5 C01 uses "not policy"' "set -- C01 6 'not policy'; $row"
mutate 6a 'clause-table: FAIL P6 check facts fails: setup-check: facts FAIL hash: docs/facts/operator-routing-policy.md' 'LC_ALL=C sed "84s/repository\\./repository!/" docs/facts/operator-routing-policy.md > x && mv x docs/facts/operator-routing-policy.md'
mutate 6b 'clause-table: FAIL P6 check facts fails: setup-check: facts FAIL listed: docs/facts/operator-routing-policy.md' 'rm docs/facts/operator-routing-policy.md && grep -v operator-routing-policy docs/setup/facts.sha256 > x && mv x docs/setup/facts.sha256'
mutate 6c 'clause-table: FAIL P6 expected exactly one docs/facts/F-0005-*.md' 'rm docs/facts/F-0005-*.md'
mutate 7a 'clause-table: FAIL P7 no entry under "### Deferred items" inside "## Decision"' 'grep -vx "### Deferred items" "$ADR" > x && mv x "$ADR"'
mutate 7b 'clause-table: FAIL P7 no entry under "### Rejected options" inside "## Decision"' 'awk "/^### Rejected options\$/ { on = 1; print; next } /^#/ { on = 0 } on && /^- / { next } { print }" "$ADR" > x && mv x "$ADR"'
mutate 7c 'clause-table: FAIL P7 D98 is defined but no clause goes to it' 'awk "{ print } /^### Proposed decisions\$/ { print \"\"; print \"- **D98.** unused\" }" "$ADR" > x && mv x "$ADR"'
mutate 7d 'clause-table: FAIL P7 no entry under "### Proposed decisions" inside "## Decision"' 'awk "/^## Decision\$/ { print \"## Context continued\"; next } /^### The clause table\$/ { print \"## Decision\"; next } { print }" "$ADR" > x && mv x "$ADR"'
mutate 1d 'clause-table: FAIL P1 no clause table header' 'f=$(printf "\140\140\140"); awk -v f="$f" "/^\\| Clause \\| F-0005 lines \\|/ { print f } { print } /^\\| C24 \\|/ { print f }" "$ADR" > x && mv x "$ADR"'
mutate 7e 'clause-table: FAIL P7 no entry under "### Rejected options" inside "## Decision"' 'f=$(printf "\140\140\140"); awk -v f="$f" "/^### Rejected options\$/ { print; print \"\"; print f; print \"- fenced\"; print f; on = 1; next } /^#/ { on = 0 } on && /^(- |  )/ { next } { print }" "$ADR" > x && mv x "$ADR"'
mutate 1e 'clause-table: FAIL P1 no clause table header' 'awk "/^\\| Clause \\| F-0005 lines \\|/ { print \"<!--\" } { print } /^\\| C24 \\|/ { print \"-->\" }" "$ADR" > x && mv x "$ADR"'
mutate 2d 'clause-table: FAIL P2 a row opens with "C25"' 'awk "{ print } /^\\| C24 \\|/ { r = \$0; sub(/C24/, \"C25\", r); print \" \" r }" "$ADR" > x && mv x "$ADR"'
mutate 5d 'clause-table: FAIL P5 C00 goes to "D1", but only "not policy" is allowed for C00' "set -- C00 6 D1; $row"
mutate 7f 'clause-table: FAIL P7 no entry under "### Rejected options" inside "## Decision"' 'awk "/^### Rejected options\$/ { print; print \"<!--\"; on = 1; next } on && /^#/ { print \"-->\"; on = 0 } { print }" "$ADR" > x && mv x "$ADR"'
mutate 2e 'clause-table: FAIL P2 a row opens with "C25"' 'awk "{ print } /^\\| C24 \\|/ { print \"C25 | L84 | extends | R5 | D3\" }" "$ADR" > x && mv x "$ADR"'
mutate 1f 'clause-table: FAIL P1 no delimiter row under the clause table header' 'awk "/^\\| -- \\| -- \\|/ && !d { d = 1; next } { print }" "$ADR" > x && mv x "$ADR"'
control k1 'a thematic break right after C24' 'awk "{ print } /^\\| C24 \\|/ { print \"---\" }" "$ADR" > x && mv x "$ADR"'
control k2 'every clause row without its outer pipes' 'awk "/^\\| C[0-9][0-9] \\|/ { sub(/^\\| /, \"\"); sub(/ \\|\$/, \"\") } { print }" "$ADR" > x && mv x "$ADR"'
mutate 5e 'clause-table: FAIL P5 C01 goes to D99, which the ADR does not define' 'R="| C01 | L5–L6 | conflicts | R5 \\| D3 | D99 |" awk "/^\\| C01 \\|/ { print ENVIRON[\"R\"]; next } { print }" "$ADR" > x && mv x "$ADR"'
control k3 'an escaped pipe inside a basis' 'R="| C01 | L5–L6 | conflicts | R5 \\| D3 | D3 |" awk "/^\\| C01 \\|/ { print ENVIRON[\"R\"]; next } { print }" "$ADR" > x && mv x "$ADR"'
printf 'clause-table self-test: %s passed, %s failed\n' "$pass" "$bad"
[ "$bad" = 0 ]
