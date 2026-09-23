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
#
# The rule for the pass-through is the argument, not the path: with no argument
# the kit linters run, with an argument they do not — even when ROOT is `.`.
#
# Output: one line per check, `setup-check: <check> OK` or
# `setup-check: <check> FAIL <cause>: <detail>` (one FAIL line per cause found).
# Exit status: 0 when every check passes, 1 otherwise.
#
# Each check is a function `check_<name>` listed in CHECKS. A later setup task
# adds its check here, with fixtures under docs/setup/tests/<name>/, in the same
# change. Self-test: sh docs/setup/tests/run.sh

set -u

CHECKS="pin"

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
check_pin() {
	pin="$ROOT/docs/setup/armature.pin"
	if [ ! -f "$pin" ]; then fail pin "missing: docs/setup/armature.pin is absent"; return; fi
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

for c in $CHECKS; do
	cur_fail=0
	"check_$c"
	[ "$cur_fail" = 0 ] && printf 'setup-check: %s OK\n' "$c"
done

# --- kit linters (pass-through, no argument only) -----------------------------
# The parent demo: one command proves the setup AND the kit's own discipline
# checks. A linter that is absent fails; it does not pass in silence (Invariant 5).
if [ "$passthrough" = 1 ]; then
	cur_fail=0
	for f in docs/tests/run-discipline-tests.sh docs/adr/adr-lint.sh docs/prd/prd-lint.sh docs/links/link-lint.sh; do
		if [ ! -f "$ROOT/$f" ]; then fail kit-linters "$f is absent"; continue; fi
		sh "$ROOT/$f" || fail kit-linters "$f"
	done
	[ "$cur_fail" = 0 ] && printf 'setup-check: kit-linters OK\n'
fi

exit "$failed"
