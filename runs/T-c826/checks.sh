#!/bin/sh
# S1 assertions for T-c826 (#449), amendment 1; no restore replay.
# Run: sh runs/T-c826/checks.sh
set -u
fail=0
assert() {
	if [ "$2" -eq 0 ]; then printf 'ok: %s\n' "$1"
	else printf 'FAIL: %s\n' "$1"; fail=$((fail + 1)); fi
}
margin=$(sed -n '/^  A landing figure above/,/^  The tolerance is one fixed/p' docs/engineering-discipline.md |
	tr '\n' ' ' | sed 's/[[:space:]][[:space:]]*/ /g')
printf '%s\n' "$margin" | grep -Fq 'The margin is the maximum plus the tolerance, minus the landing figure, calculated separately in lines and in files.'
assert 'canonical margin sentence in the tolerance section' "$?"
for runner in docs/tests/freeze-check/run.sh docs/tests/freeze-check/run-all.sh; do
	git show "d27b33180a1aa82c23f3ab66e757e06de356eb35:$runner" | cmp -s - "$runner"
	assert "unchanged runner: $runner" "$?"
done
# Select only this job and its Restore step; assertions cannot match another
# job, task-check messages, or explanatory comments in place of shell code.
job=$(awk '/^  nested-checkout-check:$/ { on=1; next }
	on && /^  [A-Za-z0-9_-]+:/ { exit } on { print }' .github/workflows/ci.yml)
restore=$(printf '%s\n' "$job" | awk '
	/      - name: Restore the checks from the default branch$/ { on=1; next }
	on && /^      - / { exit } on { print }')
code=$(printf '%s\n' "$restore" | awk '
	/^        run: \|$/ { on=1; next }
	on && /^          / {
		sub(/^[[:space:]]+/, ""); if (/^#/ || /^$/) next
		sub(/[[:space:]]+#.*/, ""); gsub(/[[:space:]]+/, " ")
		if (sub(/\\$/, "")) { pending=pending $0; next }
		print pending $0; pending=""
	}')
loop=$(printf '%s\n' "$code" | awk '/^for f in / { on=1 } on { print } on && /^done$/ { exit }')
list=$(printf '%s\n' "$loop" | sed -n '1p')
[ "$list" = 'for f in docs/adr/adr-lint.sh docs/prd/prd-lint.sh docs/links/link-lint.sh docs/tests/run-discipline-tests.sh docs/ci/pr-link-lint.sh docs/ci/review-record-lint.sh docs/tests/freeze-check.sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/run-all.sh; do' ]
assert 'all nine baseline paths stay in the restore list' "$?"
# Also preserve the source and in-place copy, rather than accepting a loop
# over an untrusted directory with otherwise identical path names.
printf '%s\n' "$code" | grep -Fxq 'git clone --quiet --depth 1 --branch "$DEFAULT_BRANCH" "https://github.com/${GITHUB_REPOSITORY}.git" .checks-from-default'
assert 'restoration still reads the default branch of this repository' "$?"
printf '%s\n' "$restore" | grep -Fxq '          DEFAULT_BRANCH: ${{ github.event.repository.default_branch }}'
assert 'default branch still comes from the repository event' "$?"
[ "$fail" -eq 0 ]
