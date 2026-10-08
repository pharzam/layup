#!/bin/sh
# The check of the documents of task T-6bq5 (#129, row 24 of the plan). A
# one-off evidence command of this task, not a check of the gate: no hook and
# no CI job runs it. Run from the root of the repository; exit 0 when each rule
# holds, 1 otherwise.
fail=0
ok() { printf 'ok   %s\n' "$1"; }
no() { printf 'FAIL %s\n' "$1"; fail=1; }
has() { grep -q -- "$2" "$1" && ok "$1: $3" || no "$1 does not say: $3"; }
# 1. forge.md (conditions 4 and 6 of #137; conditions 1 to 3 and notes 2 and 6 of #129).
has docs/spec/forge.md '`write` holds `read`' "the permissions of M2a, write holds read"
has docs/spec/forge.md 'M2c` specifies the call of `Comment`' "Comment is specified by M2c"
has docs/spec/forge.md 'not exceeding `iat` plus 9 minutes\|`exp` is `iat` plus 9 minutes' "exp is iat plus 9 minutes"
has docs/spec/forge.md 'x-ratelimit-remaining' "the reading of a rate limit"
has docs/spec/forge.md 'the login `ghost`' "a comment of a deleted account"
has docs/spec/forge.md 'that call fails with the call and the status' "the one reading of the test of the adapter"
# 2. run.md: the acceptance row of the adapter keeps the one reading (condition 3).
grep '^| The GitHub adapter |' docs/spec/run.md | grep -q 'forge.md#the-test-of-the-adapter' && ok "run.md: the acceptance row points to forge.md" || no "the acceptance row of the adapter in run.md does not point to forge.md"
# 3. guardrails.md §2: the rate-limit pitfall (note 7).
has docs/guardrails.md 'x-ratelimit-reset' "the rate-limit pitfall"
# 4. packages.md: the Job cell of internal/forge.
grep '^| `internal/forge` |' docs/spec/packages.md | grep -q 'CheckPermissions' && ok "packages.md: the Job cell of internal/forge" || no "the Job cell of internal/forge does not name CheckPermissions"
# 5. docs/tests/traceability.md.
for t in TestMissingNamesEachUndeclaredCapability TestTheJWTHasItsHeaderClaimsAndSignature TestTheAdapterPlaysEachCallOfM2a TestARateLimitWaitsUntilTheResetWithAProgressLineEveryTenSeconds; do
	grep "\`$t\`" docs/tests/traceability.md | grep -q 'T-6bq5' && ok "traceability: $t" || no "traceability has no row of $t with T-6bq5"
done
# 6. The Test cells of NFR-001 and NFR-007.
grep '^| NFR-001 | F-0001#1 ' docs/prd/PRD-0001-layup.md | grep -q 'T-6bq5' && ok "PRD §12, NFR-001" || no "the Test cell of NFR-001 does not name T-6bq5"
grep '^| NFR-007 | F-0004#1 ' docs/prd/PRD-0001-layup.md | grep -q 'T-6bq5' && ok "PRD §12, NFR-007" || no "the Test cell of NFR-007 does not name T-6bq5"
# 7. The endpoint reading (note 5).
has runs/T-6bq5/endpoints.md '2026-10-08' "the endpoint pages and their date"
exit "$fail"
