#!/bin/sh
#
# run.sh — fixture self-test for docs/setup/setup-check.sh.
#
# Each case is a directory docs/setup/tests/<check>/<case>/ with:
#   overlay/  files committed on top of the shared kit tree (may be absent)
#   EXPECT    one `key=value` per line:
#               exit=<code>        the exit code setup-check.sh must return
#               line=<text>        a line the output must hold (exact match);
#                                  repeat for more lines
#               mode=root|self|shallow   how the case is run (default: root)
#               only=<check,check>  pass `--only <check,check>`, so the case runs
#                                  only those checks (a fixture tree is not a full
#                                  setup). Default: the case's directory name, so
#                                  pin/<case> runs --only pin. frame/ is not a
#                                  check, so a frame case names its checks.
#               stub-fail=<path>   (mode=self) this stub linter exits 1
#               stub-absent=<path> (mode=self) this stub linter is not written
#
# For each case the runner builds a temp Git repo: the shared tree
# docs/setup/tests/kit/ is the ROOT commit, and overlay/ is a second commit. So
# the root-commit tree of every case is the same, and a pin can name it.
#
#   mode=root     run `sh setup-check.sh <temp repo>` — the checks, no pass-through
#   mode=self     copy setup-check.sh into the temp repo and run it with no
#                 argument — the checks and the pass-through of the kit linters
#                 (the runner writes a stub for each kit linter; each stub
#                 prints `stub <path> OK` and exits 0 unless EXPECT says otherwise)
#   mode=shallow  clone the temp repo with --depth 1, then run as mode=root
#
# Exit status: 0 when every case matches its EXPECT, 1 otherwise. A run that
# found no case fails: a harness that tested nothing is not a pass.

set -u
here=$(cd "$(dirname "$0")" && pwd)
check="$here/../setup-check.sh"
tmp=$(mktemp -d "${TMPDIR:-/tmp}/setup-check-tests.XXXXXX")
trap 'rm -rf "$tmp"' EXIT INT TERM

g() { git -c user.name=fixture -c user.email=fixture@invalid -c core.hooksPath=/dev/null \
	-c commit.gpgsign=false "$@"; }

pass=0 fail=0 n=0
for case in "$here"/*/*/; do
	case=${case%/}
	[ -f "$case/EXPECT" ] || continue
	name=${case#"$here"/}
	n=$((n + 1))
	repo="$tmp/$n/repo"
	mkdir -p "$repo"
	(cd "$repo" && g init -q -b main && cp -R "$here/kit/." . && g add -A && g commit -q -m kit) || { echo "FAIL  $name: cannot build the kit commit"; fail=$((fail + 1)); continue; }
	if [ -d "$case/overlay" ]; then cp -R "$case/overlay/." "$repo/"; fi
	(cd "$repo" && g add -A && g commit -q --allow-empty -m overlay)

	mode=$(sed -n 's/^mode=//p' "$case/EXPECT" | head -1)
	want=$(sed -n 's/^exit=//p' "$case/EXPECT" | head -1)
	only=$(sed -n 's/^only=//p' "$case/EXPECT" | head -1)
	group=${name%%/*}
	[ -z "$only" ] && [ "$group" != frame ] && only=$group
	set --
	[ -n "$only" ] && set -- --only "$only"
	case "${mode:-root}" in
	root) out=$(sh "$check" "$@" "$repo" 2>&1); got=$? ;;
	shallow)
		g clone -q --depth 1 "file://$repo" "$tmp/$n/shallow" 2>/dev/null
		out=$(sh "$check" "$@" "$tmp/$n/shallow" 2>&1); got=$? ;;
	self)
		mkdir -p "$repo/docs/setup" && cp "$check" "$repo/docs/setup/setup-check.sh"
		for f in docs/tests/run-discipline-tests.sh docs/adr/adr-lint.sh docs/prd/prd-lint.sh docs/links/link-lint.sh; do
			grep -Fqx "stub-absent=$f" "$case/EXPECT" && continue
			mkdir -p "$repo/$(dirname "$f")"
			if grep -Fqx "stub-fail=$f" "$case/EXPECT"; then
				printf '#!/bin/sh\necho "stub %s FAIL"\nexit 1\n' "$f" > "$repo/$f"
			else
				printf '#!/bin/sh\necho "stub %s OK"\nexit 0\n' "$f" > "$repo/$f"
			fi
		done
		out=$(sh "$repo/docs/setup/setup-check.sh" "$@" 2>&1); got=$? ;;
	*) echo "FAIL  $name: unknown mode '$mode'"; fail=$((fail + 1)); continue ;;
	esac

	ok=1
	if [ "$got" != "$want" ]; then ok=0; echo "FAIL  $name: exit $got, want $want"; fi
	# Read the wanted lines from the file, not from a pipe, so ok= survives.
	sed -n 's/^line=//p' "$case/EXPECT" > "$tmp/$n/lines"
	while IFS= read -r line; do
		printf '%s\n' "$out" | grep -Fqx -- "$line" || { ok=0; echo "FAIL  $name: no output line: $line"; }
	done < "$tmp/$n/lines"
	if [ "$ok" = 1 ]; then pass=$((pass + 1)); else fail=$((fail + 1)); printf '%s\n' "$out" | sed 's/^/      | /'; fi
done

if [ "$n" = 0 ]; then echo "setup-check tests: FAIL  no case found under $here"; exit 1; fi
echo "setup-check tests: $pass passed, $fail failed"
[ "$fail" = 0 ]
