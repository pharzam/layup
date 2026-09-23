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

CHECKS="pin kit-history"

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
