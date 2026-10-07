#!/bin/sh
# The checks of task T-gd8q (#120): the acceptance criteria of the ADR that ends
# bootstrap mode. A one-off evidence command of this task, not a check of the gate
# (bootstrap rule 1). Run from the root of the repository. Exit 0 when each check
# passes, 1 when one fails.
#
#   sh runs/T-gd8q/checks.sh

fail=0
ok() { printf 'ok   %s\n' "$1"; }
no() { printf 'FAIL %s\n' "$1"; fail=1; }

adr=$(ls docs/adr/0026-*.md 2>/dev/null | head -n 1)

# C1: ADR-0026 exists and is Accepted.
if [ -n "$adr" ] && sed -n '/^## Status/,/^## Context/p' "$adr" | grep -qx 'Accepted'; then
	ok "C1 ADR-0026 is Accepted ($adr)"
else
	no "C1 ADR-0026 is missing or not Accepted"
fi

# C2: the Status of ADR-0012 names ADR-0026.
if sed -n '/^## Status/,/^## Context/p' docs/adr/0012-build-layup-in-bootstrap-mode.md |
	grep -q '^Superseded by \[ADR-0026\](0026-'; then
	ok "C2 the Status of ADR-0012 is Superseded by ADR-0026"
else
	no "C2 the Status of ADR-0012 does not name ADR-0026"
fi

# C3: the sentence that F-15 found in conflict with review-record-lint is gone.
if grep -q 'One pass is never enough' docs/engineering-discipline.md; then
	no "C3 docs/engineering-discipline.md still says 'One pass is never enough'"
else
	ok "C3 'One pass is never enough' is gone"
fi

# C4: each live mention of the word is classed in mentions.tsv (path, line, class).
# The path groups that the table excludes as a group are its rows with line '*'.
tab=runs/T-gd8q/mentions.tsv
if [ ! -f "$tab" ]; then
	no "C4 $tab is missing"
else
	unclassed=$(git grep -n -i bootstrap -- . ':!runs/' ':!docs/tasks/T-*.md' |
		awk -F: -v tab="$tab" '
			BEGIN {
				while ((getline l < tab) > 0) {
					split(l, f, "\t")
					if (f[1] == "path") continue
					if (f[3] != "history" && f[3] != "link" && f[3] != "other-meaning") continue
					if (f[2] == "*") grp[f[1]] = 1; else row[f[1] ":" f[2]] = 1
				}
			}
			!(($1 ":" $2) in row) && !($1 in grp) { print $1 ":" $2 }')
	if [ -z "$unclassed" ]; then
		ok "C4 each live mention of 'bootstrap' is classed in $tab"
	else
		no "C4 mentions with no class in $tab:"
		printf '%s\n' "$unclassed" | sed 's/^/       /'
	fi
fi

# C5: the checks of the acceptance criteria exit 0.
for c in docs/adr/adr-lint.sh docs/links/link-lint.sh docs/setup/setup-check.sh; do
	if sh "$c" >/dev/null 2>&1; then ok "C5 $c exits 0"; else no "C5 $c exits non-zero"; fi
done

exit "$fail"
