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

CHECKS="pin kit-history facts onboarding glossary guardrails markers ci protection identity procedure"

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
		case "$gr_n" in [1-9]) ;; *) fail guardrails "entry: Inv-$gr_n is not one of Inv-1 to Inv-9"; continue ;; esac
		printf '%s' "$gr_text" | grep -Fq "F-0001#$gr_n\`" || fail guardrails "citation: Inv-$gr_n does not cite F-0001#$gr_n"
		gr_check=$(printf '%s' "$gr_text" | tr '\037' '\n' | sed -n 's/.*Check: //p' | head -1 | sed 's/[[:space:]]*$//')
		case "$gr_check" in
		"no check yet") ;;
		*" ("*")")
			gr_path=${gr_check%% (*}
			gr_gate=${gr_check##* (}; gr_gate=${gr_gate%)}
			[ -n "$gr_path" ] && [ -f "$ROOT/$gr_path" ] || fail guardrails "check: Inv-$gr_n names \"$gr_path\", which is not a file"
			case "$gr_gate" in hook|ci:?*) ;; *) fail guardrails "check: Inv-$gr_n gate \"$gr_gate\" is not hook or ci:<job>" ;; esac ;;
		*) fail guardrails "check: Inv-$gr_n has no valid Check: value (\"$gr_check\")" ;;
		esac
	done < "$tmpdir/gr_entries"
}

# --- markers (Invariants 4 and 5) ----------------------------------------------
# A marker is `‹` plus one or more characters other than `›`, then `›`; one that
# does not close on its line runs to the line end (its key is that first line).
# The literal `‹…›` names the convention and is not a marker, and neither is
# the exact code span `‹` (a backtick on each side: the character named, as in
# "search for `‹`"). Only that one character is skipped. Each marker in a
# git-tracked file must be exempt, allowed as a record-shape example, or listed in
# docs/setup/open-gaps.tsv (`path<TAB>marker<TAB>question`); each listed marker
# must still occur. Key: path plus exact marker text; equal markers in one file
# are one key.
MK_EXEMPT='^(docs/(adr|ci|links|prd|setup)/tests/|\.githooks/tests/|docs/templates/)|^docs/ci/[^/]*\.yml$|^docs/[^/]+/template\.md$|^docs/tests/template-[^/]*\.md$|^docs/tests/traceability-template\.md$|^docs/adr/000[1-8]-[^/]*\.md$|^docs/links/link-lint\.sh$|^docs/prd/prd-lint\.sh$|^docs/setup/setup-check\.sh$|^docs/setup/open-gaps\.tsv$'
# Record-shape examples: a marker that shows the shape of a record, not a value.
mk_allowed() {
	cat <<'MK_ALLOW'
docs/engineering-discipline.md	‹the plan and its review›
docs/engineering-discipline.md	‹the decay review rounds›
docs/engineering-discipline.md	‹writing the tests and the code›
docs/engineering-discipline.md	‹isolate, guardrails, docs, close-out›
docs/engineering-discipline.md	‹model›
docs/engineering-discipline.md	‹model / `not applicable`›
docs/engineering-discipline.md	‹effort›
docs/engineering-discipline.md	‹tokens›
docs/engineering-discipline.md	‹wall-clock›
docs/engineering-discipline.md	‹sum›
docs/engineering-discipline.md	‹task-ID›
docs/glossary.md	‹term›
docs/glossary.md	‹abbr or —›
docs/glossary.md	‹one or two sentences. State what it is and why it matters here.›
docs/glossary.md	‹a concrete instance that makes it real›
docs/tasks/backlog.md	‹ID›
docs/tasks/backlog.md	‹one-sentence summary›
docs/tasks/backlog.md	‹ADR or doc link›
docs/tasks/backlog.md	‹id›
docs/tasks/completed.md	‹ID›
docs/tasks/completed.md	‹one-sentence summary of what the task found or delivered›
docs/tasks/completed.md	‹link›
docs/tasks/completed.md	‹id›
docs/prd/README.md	‹slug›
MK_ALLOW
}
check_markers() {
	git -C "$ROOT" -c core.quotePath=false ls-files > "$tmpdir/mk_files" || { fail markers "git: cannot list the tracked files"; return; }
	grep -Ev "$MK_EXEMPT" "$tmpdir/mk_files" | while IFS= read -r mk_f; do
		[ -f "$ROOT/$mk_f" ] || continue
		awk -v f="$mk_f" '{
			line = $0; prev = ""
			while ((i = index(line, "‹")) > 0) {
				before = (i > 1) ? substr(line, i - 1, 1) : substr(prev, length(prev), 1)
				rest = substr(line, i)
				after = substr(rest, length("‹") + 1, 1)
				# A mention is exactly `‹` in a code span; skip that one character only.
				if (before == "`" && after == "`") { prev = substr(line, 1, i + length("‹") - 1); line = substr(rest, length("‹") + 1); continue }
				j = index(rest, "›")
				if (j > 0) { m = substr(rest, 1, j + length("›") - 1); prev = ""; line = substr(rest, j + length("›")) }
				else { m = rest; line = "" }
				if (m != "‹…›") print f "\t" m
			}
		}' "$ROOT/$mk_f"
	done | sort -u > "$tmpdir/mk_found"
	mk_allowed | sort -u > "$tmpdir/mk_allow"
	mk_gaps="$ROOT/docs/setup/open-gaps.tsv"
	if [ -f "$mk_gaps" ]; then cut -f1,2 "$mk_gaps" | grep . | sort -u > "$tmpdir/mk_listed"; else : > "$tmpdir/mk_listed"; fi
	# Each open gap carries its question; a row without one asks nothing.
	[ -f "$mk_gaps" ] && awk -F'\t' 'NF && $3 == "" { print "setup-check: markers FAIL question: docs/setup/open-gaps.tsv line " NR " has no question" }' "$mk_gaps" > "$tmpdir/mk_q"
	if [ -s "$tmpdir/mk_q" ]; then cat "$tmpdir/mk_q"; cur_fail=1; failed=1; fi
	sort -u "$tmpdir/mk_allow" "$tmpdir/mk_listed" > "$tmpdir/mk_known"
	comm -23 "$tmpdir/mk_found" "$tmpdir/mk_known" | while IFS="$(printf '\t')" read -r mk_p mk_m; do
		printf 'setup-check: markers FAIL unlisted: %s %s\n' "$mk_p" "$mk_m"
	done > "$tmpdir/mk_out"
	comm -13 "$tmpdir/mk_found" "$tmpdir/mk_listed" | while IFS="$(printf '\t')" read -r mk_p mk_m; do
		printf 'setup-check: markers FAIL stale: %s %s is listed in docs/setup/open-gaps.tsv but does not occur\n' "$mk_p" "$mk_m"
	done >> "$tmpdir/mk_out"
	if [ -s "$tmpdir/mk_out" ]; then cat "$tmpdir/mk_out"; cur_fail=1; failed=1; fi
}

# --- ci (the kit's own CI replaced; setup-check runs in CI) --------------------
# No workflow may call itself the Armature repo's CI. ci.yml must run the fixture
# self-test and the setup check, with the full history (the pin tree check reads
# the root commit, and a shallow clone fails it).
check_ci() {
	for ci_f in "$ROOT"/.github/workflows/*.yml "$ROOT"/.github/workflows/*.yaml; do
		[ -f "$ci_f" ] || continue
		grep -Fq 'Armature repo' "$ci_f" \
			&& fail ci "kit-header: .github/workflows/$(basename "$ci_f") says it is the Armature repo's workflow"
	done
	ci_yml="$ROOT/.github/workflows/ci.yml"
	for ci_cmd in 'sh docs/setup/setup-check.sh' 'sh docs/setup/tests/run.sh'; do
		grep -Fq -- "$ci_cmd" "$ci_yml" 2>/dev/null || fail ci "job: .github/workflows/ci.yml does not run $ci_cmd"
	done
	grep -Eq 'fetch-depth:[[:space:]]*0([^0-9]|$)' "$ci_yml" 2>/dev/null \
		|| fail ci "history: .github/workflows/ci.yml has no fetch-depth: 0 (the pin check needs the root commit)"
	# Cause `restore` (the #84 guardrail): per job, each `docs/...sh` path that the
	# job names outside its restore step (a step whose name starts with
	# "Restore"), in any form — `sh`, `bash`, `/bin/sh`, quoted, `./docs/...`,
	# `$GITHUB_WORKSPACE/docs/...` — counts as a script the job runs, and must be
	# named inside the restore step by its path or a parent directory. Comment
	# lines count on neither side; CR line ends are ignored. A job is a key under
	# `jobs:` at the first key's indent. setup-check.sh runs the four kit linters
	# and the discipline runner; the runner tests the two docs/ci lint scripts.
	# nested-checkout-check.sh is exempt: ci.yml documents why it runs from the
	# branch. LIMIT: a script run through `working-directory:` with a relative
	# path is not seen. Both *.yml and *.yaml workflows are read.
	for ci_f in "$ROOT"/.github/workflows/*.yml "$ROOT"/.github/workflows/*.yaml; do
		[ -f "$ci_f" ] || continue
		ci_name=.github/workflows/$(basename "$ci_f")
		awk '
			function ind(l) { match(l, /^ */); return RLENGTH }
			{ sub(/\r$/, "") }
			/^jobs:/ { inj = 1; jind = -1; next }
			inj && /^[^ #]/ { inj = 0 }
			!inj { next }
			/^[ \t]*#/ { next }
			{
				if (jind < 0 && $0 ~ /^ +[A-Za-z0-9_-]+:[ \t]*(#.*)?$/) jind = ind($0)
				if (ind($0) == jind && $0 ~ /^ +[A-Za-z0-9_-]+:[ \t]*(#.*)?$/) {
					job = $0; sub(/^ +/, "", job); sub(/:.*/, "", job); inrs = 0; next
				}
				if (job == "") next
				if ($0 ~ /^ *- /) { sind = ind($0); inrs = ($0 ~ /^ *- name:[ \t]*["\047]?Restore/) }
				else if (inrs && ind($0) <= sind) inrs = 0
				if ($0 ~ /^ *name:[ \t]*["\047]?Restore/) inrs = 1
				line = $0
				while (match(line, /docs\/[A-Za-z0-9_.\/-]+/)) {
					tok = substr(line, RSTART, RLENGTH); pre = substr(line, 1, RSTART - 1)
					sub(/\/$/, "", tok)
					if (inrs) print job "\tmention\t" tok
					else if (tok ~ /\.sh$/) print job "\trun\t" tok
					line = substr(line, RSTART + RLENGTH)
				}
			}' "$ci_f" > "$tmpdir/ci_tok"
		for ci_job in $(cut -f1 "$tmpdir/ci_tok" | sort -u); do
			awk -F'\t' -v j="$ci_job" '$1 == j && $2 == "run" { print $3 }' "$tmpdir/ci_tok" | sort -u > "$tmpdir/ci_need"
			awk -F'\t' -v j="$ci_job" '$1 == j && $2 == "mention" { print $3 }' "$tmpdir/ci_tok" | sort -u > "$tmpdir/ci_named"
			if grep -qx docs/setup/setup-check.sh "$tmpdir/ci_need"; then
				printf '%s\n' docs/adr/adr-lint.sh docs/prd/prd-lint.sh docs/links/link-lint.sh docs/tests/run-discipline-tests.sh >> "$tmpdir/ci_need"
			fi
			if cat "$tmpdir/ci_need" "$tmpdir/ci_named" | grep -qx docs/tests/run-discipline-tests.sh; then
				printf '%s\n' docs/ci/pr-link-lint.sh docs/ci/review-record-lint.sh >> "$tmpdir/ci_need"
			fi
			sort -u "$tmpdir/ci_need" | grep -vx docs/tests/nested-checkout-check.sh | while IFS= read -r ci_s; do
				ci_ok=0 ci_p=$ci_s
				while :; do
					grep -qx -- "$ci_p" "$tmpdir/ci_named" && { ci_ok=1; break; }
					case "$ci_p" in */*) ci_p=${ci_p%/*} ;; *) break ;; esac
				done
				[ "$ci_ok" = 1 ] || printf 'setup-check: ci FAIL restore: %s job %s runs or needs %s but does not restore it\n' "$ci_name" "$ci_job" "$ci_s"
			done
		done
	done > "$tmpdir/ci_restore"
	if [ -s "$tmpdir/ci_restore" ]; then cat "$tmpdir/ci_restore"; cur_fail=1; failed=1; fi
}

# --- protection (branch protection kept in Git) --------------------------------
# docs/setup/branch-protection.json is the body of the `PUT` that protects main
# (Invariant 1: the forge setting also lives in the repository). Its required
# contexts must equal the check-run names of every job in every workflow: the
# job's own `name:` (the first key at the job's child indent), else its id.
# The live setting is compared by hand, with the command in the setup record.
check_protection() {
	pr_json="$ROOT/docs/setup/branch-protection.json"
	if [ ! -f "$pr_json" ]; then fail protection "missing: docs/setup/branch-protection.json is absent"; return; fi
	# The PUT replaces the whole protection object, so a partial body is
	# destructive: the four parameters GitHub requires must be present, and each
	# check must name its app, or any token could post a status under the name.
	for pr_key in required_status_checks enforce_admins required_pull_request_reviews restrictions; do
		grep -Fq "\"$pr_key\"" "$pr_json" || fail protection "body: docs/setup/branch-protection.json has no \"$pr_key\" (a partial body is destructive)"
	done
	pr_nctx=$(grep -o '"context"' "$pr_json" | grep -c .)
	pr_napp=$(grep -Eo '"app_id"[[:space:]]*:[[:space:]]*15368([^0-9]|$)' "$pr_json" | grep -c .)
	[ "$pr_nctx" = "$pr_napp" ] || fail protection "body: $pr_nctx contexts but $pr_napp app_id values of 15368 (pin each check to GitHub Actions)"
	grep -o '"context"[[:space:]]*:[[:space:]]*"[^"]*"' "$pr_json" | sed 's/.*"\([^"]*\)"$/\1/' | sort | uniq -d | while IFS= read -r pr_d; do
		printf 'setup-check: protection FAIL body: context %s appears more than once\n' "$pr_d"; done > "$tmpdir/pr_dup"
	if [ -s "$tmpdir/pr_dup" ]; then cat "$tmpdir/pr_dup"; cur_fail=1; failed=1; fi
	grep -o '"context"[[:space:]]*:[[:space:]]*"[^"]*"' "$pr_json" | sed 's/.*"\([^"]*\)"$/\1/' | sort -u > "$tmpdir/pr_req"
	for pr_f in "$ROOT"/.github/workflows/*.yml "$ROOT"/.github/workflows/*.yaml; do
		[ -f "$pr_f" ] || continue
		awk '
			function ind(l) { match(l, /^ */); return RLENGTH }
			{ sub(/\r$/, "") }
			/^jobs:/ { inj = 1; jind = -1; next }
			inj && /^[^ #]/ { inj = 0 }
			!inj || /^[ \t]*(#.*)?$/ { next }
			{
				if (jind < 0) jind = ind($0)
				if (ind($0) == jind) { if (job != "") print (name != "" ? name : job); job = $0; sub(/^ +/, "", job); sub(/:.*/, "", job); name = ""; cind = -1; next }
				if (cind < 0) cind = ind($0)
				if (ind($0) == cind && $0 ~ /^ +name:/) { name = $0; sub(/^ +name:[ \t]*/, "", name); sub(/[ \t]+#.*$/, "", name); gsub(/^["\047]|["\047]$/, "", name) }
			}
			END { if (job != "") print (name != "" ? name : job) }' "$pr_f"
	done | sort -u > "$tmpdir/pr_jobs"
	comm -23 "$tmpdir/pr_jobs" "$tmpdir/pr_req" | while IFS= read -r pr_n; do
		printf 'setup-check: protection FAIL contexts: %s is a job but not a required context\n' "$pr_n"; done > "$tmpdir/pr_out"
	comm -13 "$tmpdir/pr_jobs" "$tmpdir/pr_req" | while IFS= read -r pr_n; do
		printf 'setup-check: protection FAIL contexts: %s is a required context but not a job\n' "$pr_n"; done >> "$tmpdir/pr_out"
	if [ -s "$tmpdir/pr_out" ]; then cat "$tmpdir/pr_out"; cur_fail=1; failed=1; fi
}

# --- identity (the repository is LAYUP, not the kit) ---------------------------
# The two entry files must not say that this repository is the Armature kit, and
# README.md links the pin that says which Armature version LAYUP was built from.
check_identity() {
	for id_f in README.md AGENTS.md; do
		[ -f "$ROOT/$id_f" ] || continue
		for id_p in 'Agent context for **Armature**' 'This repository is a generic **template**' 'A domain-free **template**, not a product'; do
			grep -Fq -- "$id_p" "$ROOT/$id_f" && fail identity "kit: $id_f says the repository is the Armature kit (\"$id_p\")"
		done
	done
	grep -Fq '](docs/setup/armature.pin)' "$ROOT/README.md" 2>/dev/null \
		|| fail identity "pin: README.md has no link to docs/setup/armature.pin"
}

# --- procedure (the setup procedure, for later automation) ---------------------
# docs/setup/steps.tsv is the machine-readable procedure: a fixed header and six
# tab-separated columns per row (CR line ends ignored). Each step ID has a
# `### <id>` heading in docs/setup/README.md and each `### S<digits>` heading has
# a row; other `###` headings are free text. The kit section "How to
# adapt this kit" is gone: the procedure replaced it.
check_procedure() {
	pc_tsv="$ROOT/docs/setup/steps.tsv"
	pc_md="$ROOT/docs/setup/README.md"
	if [ ! -f "$pc_tsv" ]; then fail procedure "missing: docs/setup/steps.tsv is absent"
	else
		pc_want="$(printf 'id\tinput\taction\toutput\tevidence\thuman_decision')"
		[ "$(head -1 "$pc_tsv" | tr -d '\r')" = "$pc_want" ] \
			|| fail procedure "header: docs/setup/steps.tsv header is not id, input, action, output, evidence, human_decision"
		awk -F'\t' '
			{ sub(/\r$/, "") }
			NR == 1 { next }
			NF != 6 { print "setup-check: procedure FAIL columns: docs/setup/steps.tsv line " NR " has " NF " columns, expected 6" }
			$1 == "" { print "setup-check: procedure FAIL id: docs/setup/steps.tsv line " NR " has an empty id" }
			$1 != "" && seen[$1]++ { print "setup-check: procedure FAIL id: " $1 " appears more than once in docs/setup/steps.tsv" }
			NF >= 6 && $6 != "yes" && $6 != "no" { print "setup-check: procedure FAIL decision: docs/setup/steps.tsv line " NR " human_decision is \"" $6 "\", expected yes or no" }
		' "$pc_tsv" > "$tmpdir/pc_cols"
		if [ -s "$tmpdir/pc_cols" ]; then cat "$tmpdir/pc_cols"; cur_fail=1; failed=1; fi
		awk -F'\t' 'NR > 1 && $1 != "" { print $1 }' "$pc_tsv" | sort -u > "$tmpdir/pc_ids"
		sed -n 's/^### \(S[0-9][0-9]*\)\([^0-9].*\)\{0,1\}$/\1/p' "$pc_md" 2>/dev/null | sort -u > "$tmpdir/pc_heads"
		comm -23 "$tmpdir/pc_ids" "$tmpdir/pc_heads" | while IFS= read -r pc_i; do
			printf 'setup-check: procedure FAIL step: %s is in steps.tsv but has no "### %s" heading in docs/setup/README.md\n' "$pc_i" "$pc_i"; done > "$tmpdir/pc_out"
		comm -13 "$tmpdir/pc_ids" "$tmpdir/pc_heads" | while IFS= read -r pc_i; do
			printf 'setup-check: procedure FAIL step: %s has a heading in docs/setup/README.md but is not in steps.tsv\n' "$pc_i"; done >> "$tmpdir/pc_out"
		if [ -s "$tmpdir/pc_out" ]; then cat "$tmpdir/pc_out"; cur_fail=1; failed=1; fi
	fi
	grep -q '^## How to adapt this kit' "$ROOT/docs/engineering-discipline.md" 2>/dev/null \
		&& fail procedure "kit: docs/engineering-discipline.md still holds the heading \"How to adapt this kit\""
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
