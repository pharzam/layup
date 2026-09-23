#!/bin/sh
#
# setup-check.sh — check that this repository's Armature setup is complete and
# that each setup value has its evidence (PSB Invariants 4, 5, 8).
#
# Usage:
#   sh docs/setup/setup-check.sh          check THIS repository (the one that holds
#                                         the script), then run the kit linters
#   sh docs/setup/setup-check.sh ROOT     check the repository at ROOT only; the
#                                         kit linters do not run (fixtures use this)
#   --only a,b   (before ROOT) run only the named checks. Fixtures use it, because
#                a fixture tree is not a full setup. It does not change the rule
#                for the kit linters below.
#
# The rule for the pass-through is the argument, not the path: with no argument
# the kit linters run, with an argument they do not — even when ROOT is `.`.
#
# Output: one line per check, `setup-check: <check> OK` or
# `setup-check: <check> FAIL <cause>: <detail>` (one FAIL line per cause found).
# Each pin key must appear exactly once; an absent key fails as `key: <k> appears 0 times`.
# Exit status: 0 when every check passes, 1 otherwise.
#
# Each check is a function `check_<name>` (a `-` in the name is `_`) listed in CHECKS. A later setup task
# adds its check here, with fixtures under docs/setup/tests/<name>/, in the same
# change. Self-test: sh docs/setup/tests/run.sh

set -u

CHECKS="pin kit-history facts onboarding glossary guardrails"

if [ "${1:-}" = --only ]; then
	[ $# -ge 2 ] && [ -n "$2" ] || { echo "setup-check: --only needs a list of checks" >&2; exit 2; }
	for want in $(printf '%s' "$2" | tr ',' ' '); do
		case " $CHECKS " in *" $want "*) ;; *) echo "setup-check: unknown check: $want" >&2; exit 2 ;; esac
	done
	CHECKS=$(printf '%s' "$2" | tr ',' ' ')
	shift 2
fi

if [ $# -gt 0 ]; then
	ROOT=$1
	passthrough=0
else
	ROOT=$(cd "$(dirname "$0")/../.." && pwd)
	passthrough=1
fi

failed=0
cur_fail=0
tmpdir=$(mktemp -d "${TMPDIR:-/tmp}/setup-check.XXXXXX")
trap 'rm -rf "$tmpdir"' EXIT INT TERM
fail() { # fail <check> <cause>: <detail>
	printf 'setup-check: %s FAIL %s\n' "$1" "$2"
	cur_fail=1
	failed=1
}

# --- pin (Invariant 8) --------------------------------------------------------
# The pin names the Armature commit and its tree. The tree must equal the tree of
# this repository's one root commit, which is the unmodified kit copy.
pin_value() { sed -n "s/^$1=//p" "$ROOT/docs/setup/armature.pin" | head -1; }
pin_count() { grep -c "^$1=" "$ROOT/docs/setup/armature.pin"; }
check_pin() {
	pin="$ROOT/docs/setup/armature.pin"
	if [ ! -f "$pin" ]; then fail pin "missing: docs/setup/armature.pin is absent"; return; fi
	# Each key once: with two `commit=` lines the pin has two readings.
	# POSIX sh has no local variables, so these names start with pin_.
	for pin_key in source commit tree method date; do
		pin_n=$(pin_count "$pin_key")
		[ "$pin_n" = 1 ] || fail pin "key: $pin_key appears $pin_n times, expected 1"
	done
	commit=$(pin_value commit)
	tree=$(pin_value tree)
	if ! printf '%s\n' "$commit" | grep -Eqx '[0-9a-f]{40}'; then
		fail pin "commit: not 40 hexadecimal characters: $commit"
	fi
	if [ "$(git -C "$ROOT" rev-parse --is-shallow-repository 2>/dev/null)" = true ]; then
		fail pin "shallow: the clone is shallow, so the root commit is not visible (fetch the full history)"
		return
	fi
	roots=$(git -C "$ROOT" rev-list --max-parents=0 HEAD 2>/dev/null)
	nroots=$(printf '%s\n' "$roots" | grep -c .)
	if [ "$nroots" != 1 ]; then
		fail pin "tree: expected one root commit, found $nroots"
		return
	fi
	root_tree=$(git -C "$ROOT" rev-parse "$roots^{tree}")
	if [ "$tree" != "$root_tree" ]; then
		fail pin "tree: pin names $tree, root commit has $root_tree"
	fi
}

# --- kit-history (kit step 4) -------------------------------------------------
# The kit's own history is not LAYUP state. A task detail file must have a line in
# a task index that holds its ID; no task index may link the kit's own issues.
check_kit_history() {
	for kh_dir in decisions audit; do
		[ -e "$ROOT/docs/$kh_dir" ] && fail kit-history "$kh_dir: docs/$kh_dir/ exists (kit step 4 deletes it)"
	done
	for kh_file in "$ROOT"/docs/tasks/T-*.md; do
		[ -f "$kh_file" ] || continue
		kh_id=$(basename "$kh_file" .md)
		cat "$ROOT/docs/tasks/backlog.md" "$ROOT/docs/tasks/completed.md" 2>/dev/null | grep -Fq -- "$kh_id" \
			|| fail kit-history "orphan: docs/tasks/$kh_id.md has no line with $kh_id in backlog.md or completed.md"
	done
	for kh_index in backlog completed; do
		grep -Fq 'github.com/pharzam/armature/' "$ROOT/docs/tasks/$kh_index.md" 2>/dev/null \
			&& fail kit-history "kit-link: docs/tasks/$kh_index.md links github.com/pharzam/armature/ (a kit task or note)"
	done
}

# --- facts (raw facts, docs/facts/README.md) ----------------------------------
# A raw facts file is evidence: it must match the hash recorded when it was
# collected. Each numbered fact of the F-0001 record must be a byte-exact
# substring of the PSB file (the list number `N. ` removed), and the record must
# hold facts 1 to 39 each exactly once. The index lists F-0001 and F-0002.
sha256_of() {
	if command -v sha256sum >/dev/null 2>&1; then sha256sum "$1" | cut -d' ' -f1
	else shasum -a 256 "$1" | cut -d' ' -f1; fi
}
check_facts() {
	fa_sums="$ROOT/docs/setup/facts.sha256"
	if [ ! -f "$fa_sums" ]; then
		fail facts "hash: docs/setup/facts.sha256 is absent"
	else
		# `|| [ -n ... ]` reads a last line that has no final newline.
		fa_ln=0
		while read -r fa_want fa_path || [ -n "$fa_want" ]; do
			fa_ln=$((fa_ln + 1))
			[ -n "$fa_want" ] || continue
			if [ -z "$fa_path" ]; then fail facts "hash: line $fa_ln of docs/setup/facts.sha256 has no path"; continue; fi
			if [ ! -f "$ROOT/$fa_path" ] || [ "$(sha256_of "$ROOT/$fa_path")" != "$fa_want" ]; then
				fail facts "hash: $fa_path does not match docs/setup/facts.sha256"
			fi
		done < "$fa_sums"
		# The floor: both raw facts files must be hashed, so an empty list is no pass.
		for fa_need in docs/facts/problem-statement-brief.md docs/facts/architectural-vision-brief.md; do
			awk -v p="$fa_need" '$2 == p { found = 1 } END { exit !found }' "$fa_sums" \
				|| fail facts "listed: $fa_need is not in docs/setup/facts.sha256"
		done
	fi
	fa_src="$ROOT/docs/facts/problem-statement-brief.md"
	fa_rec=$(ls "$ROOT"/docs/facts/F-0001-*.md 2>/dev/null)
	if [ -z "$fa_rec" ] || [ "$(printf '%s\n' "$fa_rec" | grep -c .)" != 1 ]; then
		fail facts "record: expected one docs/facts/F-0001-*.md"
	elif [ ! -f "$fa_src" ]; then
		fail facts "source: docs/facts/problem-statement-brief.md is absent"
	else
		# Numbers are read as integers, so `01.` and `1.` are the same fact.
		grep -E '^[0-9]+\. ' "$fa_rec" | sed -E 's/^0*([0-9])/\1/' > "$tmpdir/facts" || true
		while IFS= read -r fa_line; do
			fa_n=${fa_line%%. *}
			fa_text=${fa_line#*. }
			if [ -z "$(printf '%s' "$fa_text" | tr -d ' ')" ]; then fail facts "verbatim: F-0001 fact $fa_n is empty"; continue; fi
			if [ "$fa_n" -lt 1 ] || [ "$fa_n" -gt 39 ]; then fail facts "numbering: F-0001 fact $fa_n is outside 1..39"; fi
			grep -Fq -- "$fa_text" "$fa_src" \
				|| fail facts "verbatim: F-0001 fact $fa_n is not a byte-exact substring of docs/facts/problem-statement-brief.md"
		done < "$tmpdir/facts"
		fa_distinct=$(sed 's/\..*//' "$tmpdir/facts" | awk '$1 >= 1 && $1 <= 39' | sort -un | grep -c .)
		[ "$fa_distinct" = 39 ] || fail facts "numbering: F-0001 holds $fa_distinct distinct fact numbers in 1..39, expected 39"
		sed 's/\..*//' "$tmpdir/facts" | sort -n | uniq -c | while read -r fa_c fa_num; do
			[ "$fa_c" = 1 ] || printf 'setup-check: facts FAIL numbering: F-0001 fact %s appears %s times\n' "$fa_num" "$fa_c"
		done > "$tmpdir/repeats"
		if [ -s "$tmpdir/repeats" ]; then cat "$tmpdir/repeats"; cur_fail=1; failed=1; fi
	fi
	for fa_id in F-0001 F-0002; do
		grep -Eq "^\|.*$fa_id" "$ROOT/docs/facts/README.md" 2>/dev/null \
			|| fail facts "index: docs/facts/README.md has no row for $fa_id"
	done
}

# --- onboarding (bound to the PSB) ---------------------------------------------
# The first door has no unfilled marker, links the PSB file, and cites only facts
# that the F-0001 record holds. A section citation such as `F-0001 §6` is not
# checked; the pattern is `F-0001#<digits>`.
check_onboarding() {
	on_doc="$ROOT/docs/onboarding-for-engineers.md"
	if [ ! -f "$on_doc" ]; then fail onboarding "missing: docs/onboarding-for-engineers.md is absent"; return; fi
	grep -q '‹' "$on_doc" && fail onboarding "marker: docs/onboarding-for-engineers.md holds a ‹ character"
	grep -Fq '](facts/problem-statement-brief.md)' "$on_doc" \
		|| fail onboarding "link: docs/onboarding-for-engineers.md has no link to facts/problem-statement-brief.md"
	on_rec=$(ls "$ROOT"/docs/facts/F-0001-*.md 2>/dev/null | head -1)
	grep -Eo 'F-0001#[0-9]+' "$on_doc" | sort -u > "$tmpdir/cited"
	while IFS= read -r on_id; do
		on_n=${on_id#F-0001#}
		if [ -z "$on_rec" ] || ! grep -Eq "^0*$on_n\. " "$on_rec"; then
			fail onboarding "fact: $on_id is not a fact of the F-0001 record"
		fi
	done < "$tmpdir/cited"
}

# --- glossary (PSB §8 terms) ---------------------------------------------------
# The section `## 1. LAYUP domain (PSB §8)` holds the 25 rows of PSB §8, and its
# rows cite F-0001#15 to F-0001#39 each exactly once. The section ends at the next
# `## ` heading. A row is a table line that is not the header or the separator.
check_glossary() {
	gl_doc="$ROOT/docs/glossary.md"
	gl_head='## 1. LAYUP domain (PSB §8)'
	if ! grep -Fxq "$gl_head" "$gl_doc" 2>/dev/null; then
		fail glossary "section: docs/glossary.md has no heading \"$gl_head\""; return
	fi
	awk -v h="$gl_head" '$0 == h { on = 1; next } on && /^## / { on = 0 } on' "$gl_doc" \
		| grep '^|' | grep -v '^| Term |' | grep -v '^|[-|: ]*$' > "$tmpdir/gl_rows"
	gl_n=$(grep -c . "$tmpdir/gl_rows")
	[ "$gl_n" = 25 ] || fail glossary "rows: the LAYUP domain section holds $gl_n rows, expected 25"
	gl_i=15
	while [ "$gl_i" -le 39 ]; do
		gl_c=$(grep -o "F-0001#$gl_i\`" "$tmpdir/gl_rows" | grep -c .)
		[ "$gl_c" = 1 ] || fail glossary "citation: F-0001#$gl_i is cited $gl_c times, expected 1"
		gl_i=$((gl_i + 1))
	done
}

# --- guardrails (the 9 System Invariants) --------------------------------------
# Under `### 1.1 LAYUP System Invariants (PSB §6)`, an entry starts with a line
# `- **Inv-N**` and runs to the next entry or heading. Inv-1 to Inv-9 appear once
# each; entry N cites F-0001#N; its `Check:` value is `no check yet` or
# `<existing path> (<gate>)` with gate `hook` or `ci:<job>`. Whether the gate
# really runs the path is a review judgement, not this check.
check_guardrails() {
	gr_doc="$ROOT/docs/guardrails.md"
	gr_head='### 1.1 LAYUP System Invariants (PSB §6)'
	if ! grep -Fxq "$gr_head" "$gr_doc" 2>/dev/null; then
		fail guardrails "section: docs/guardrails.md has no heading \"$gr_head\""; return
	fi
	# One line per entry: `N<TAB>entry lines joined by the unit separator (octal 037)`,
	# so a `Check:` value can end where its own line ends.
	awk -v h="$gr_head" '
		$0 == h { on = 1; next }
		on && /^#/ { on = 0 }
		!on { next }
		/^- \*\*Inv-[0-9]+\*\*/ { if (cur != "") print cur; n = $0; sub(/^- \*\*Inv-/, "", n); sub(/\*\*.*/, "", n); cur = n "\t" $0; next }
		cur != "" { cur = cur "\037" $0 }
		END { if (cur != "") print cur }' "$gr_doc" > "$tmpdir/gr_entries"
	gr_i=1
	while [ "$gr_i" -le 9 ]; do
		gr_c=$(cut -f1 "$tmpdir/gr_entries" | grep -cx "$gr_i")
		[ "$gr_c" = 1 ] || fail guardrails "entry: Inv-$gr_i appears $gr_c times, expected 1"
		gr_i=$((gr_i + 1))
	done
	while IFS="$(printf '\t')" read -r gr_n gr_text; do
		printf '%s' "$gr_text" | grep -Fq "F-0001#$gr_n\`" || fail guardrails "citation: Inv-$gr_n does not cite F-0001#$gr_n"
		gr_check=$(printf '%s' "$gr_text" | tr '\037' '\n' | sed -n 's/.*Check: //p' | head -1 | sed 's/[[:space:]]*$//')
		case "$gr_check" in
		"no check yet") ;;
		*" ("*")")
			gr_path=${gr_check%% (*}
			gr_gate=${gr_check##* (}; gr_gate=${gr_gate%)}
			[ -e "$ROOT/$gr_path" ] || fail guardrails "check: Inv-$gr_n names $gr_path, which does not exist"
			case "$gr_gate" in hook|ci:?*) ;; *) fail guardrails "check: Inv-$gr_n gate \"$gr_gate\" is not hook or ci:<job>" ;; esac ;;
		*) fail guardrails "check: Inv-$gr_n has no valid Check: value (\"$gr_check\")" ;;
		esac
	done < "$tmpdir/gr_entries"
}

for check_name in $CHECKS; do
	cur_fail=0
	"check_$(printf %s "$check_name" | tr - _)"
	[ "$cur_fail" = 0 ] && printf 'setup-check: %s OK\n' "$check_name"
done

# --- kit linters (pass-through, no argument only) -----------------------------
# The parent demo: one command proves the setup AND the kit's own discipline
# checks. A linter that is absent fails; it does not pass in silence (Invariant 5).
if [ "$passthrough" = 1 ]; then
	cur_fail=0
	for f in docs/tests/run-discipline-tests.sh docs/adr/adr-lint.sh docs/prd/prd-lint.sh docs/links/link-lint.sh; do
		if [ ! -f "$ROOT/$f" ]; then fail kit-linters "absent: $f"; continue; fi
		sh "$ROOT/$f" || fail kit-linters "failed: $f"
	done
	[ "$cur_fail" = 0 ] && printf 'setup-check: kit-linters OK\n'
fi

exit "$failed"
