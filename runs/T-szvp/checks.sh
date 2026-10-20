#!/bin/sh
set -eu

job() {
    awk '/^  nested-checkout-check:$/ { on=1; next }
        on && /^  [A-Za-z0-9_-]+:/ { exit }
        on { print }'
}
restore() {
    awk '
        /^      - name: Restore the checks from the default branch$/ {
            on=1; print; next
        }
        on && /^      - / { exit }
        on { print }
    '
}

failed=0
current=$(job < .github/workflows/ci.yml)
test -n "$current"

if printf '%s\n' "$current" |
    grep -nE '(^|[^0-9])453([^0-9]|$)|kept as written for PR|PR_NUMBER|EVENT_NAME|PR[[:space:]]+[0-9]+'
then
    printf 'FAIL: numbered checker exception remains\n'
    failed=1
else
    rc=$?
    [ "$rc" -eq 1 ] || exit "$rc"
fi

baseline=$(git show d27b33180a1aa82c23f3ab66e757e06de356eb35:.github/workflows/ci.yml)
expected=$(printf '%s\n' "$baseline" | job | restore)
actual=$(printf '%s\n' "$current" | restore)
test -n "$expected"
if [ "$actual" != "$expected" ]; then
    printf 'FAIL: restore step differs from pre-exception baseline\n'
    failed=1
fi
exit "$failed"
