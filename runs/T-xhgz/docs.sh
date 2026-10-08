#!/bin/sh
# The check of the documents of task T-xhgz (#128, row 23 of the plan). A
# one-off evidence command of this task, not a check of the gate: no hook and
# no CI job runs it. Run from the root of the repository; exit 0 when each rule
# holds, 1 otherwise.
fail=0
ok() { printf 'ok   %s\n' "$1"; }
no() { printf 'FAIL %s\n' "$1"; fail=1; }
has() { grep -q -- "$2" "$1" && ok "$1: $3" || no "$1 does not say: $3"; }
# 1. packages.md: the row of Fetch with --update-head-ok (condition 2).
grep '^| `Fetch` |' docs/spec/packages.md | grep -q -- '--update-head-ok' && ok "packages.md: the row of Fetch" || no "the row of Fetch has no --update-head-ok"
# 2. packages.md: the input rules and the one reading of a refused push (condition 1).
has docs/spec/packages.md 'is a `FailedError` with `Code` 1' "a refused push is Code 1"
has docs/spec/packages.md 'starts with `refs/` and holds no `:`' "the input rules of Fetch and Push"
# 3. packages.md: the fixed list names the three values of a token.
has docs/spec/packages.md '`Fetch` and `Push` add `GIT_CONFIG_COUNT`' "the fixed list and the token"
# 4. forge.md: the value of the header.
has docs/spec/forge.md 'Authorization: Basic' "the value of the header"
# 5. guardrails.md §2: the pitfall of a fetch into the checked-out branch (note 6).
has docs/guardrails.md 'refusing to fetch into branch' "the fetch pitfall"
# 6. docs/tests/traceability.md and the Test cell of NFR-001.
for t in TestPushAndFetchOfABareRepository TestTheTokenReachesTheServerAsAHeaderOnly TestFetchAndPushRefuseTheirInputBeforeGitStarts; do
	grep "\`$t\`" docs/tests/traceability.md | grep -q 'T-xhgz' && ok "traceability: $t" || no "traceability has no row of $t with T-xhgz"
done
grep '^| NFR-001 | F-0001#1 ' docs/prd/PRD-0001-layup.md | grep -q 'T-xhgz' && ok "PRD §12, NFR-001" || no "the Test cell of NFR-001 does not name T-xhgz"
exit "$fail"
