#!/bin/sh
# The check of D2 of task T-zwke (#124): each row of "The tasks of M2a" links the
# issue of runs/T-zwke/issues.tsv, and the App's bot opened each issue. It reads
# the forge with gh, so it runs by hand, not in a hook or CI; a one-off evidence
# command. Run from the root of the repository; exit 0 when each rule holds.
tab=runs/T-zwke/issues.tsv
plan=docs/plan/README.md
[ -f "$tab" ] || { echo "FAIL $tab is missing: no issue is open"; exit 1; }
fail=0
tail -n +2 "$tab" | while IFS='	' read -r row task num; do
	line=$(awk -v r="$row" '/^## The tasks of M2a/{t=1; next} t && /^## /{t=0} t && $2==r' "$plan")
	case "$line" in
	*"[#$num](https://github.com/pharzam/layup/issues/$num)"*"\`$task\`"*|*"\`$task\`"*"[#$num](https://github.com/pharzam/layup/issues/$num)"*) ;;
	*) echo "FAIL row $row of The tasks of M2a does not link #$num with $task"; continue ;;
	esac
	who=$(gh api "repos/pharzam/layup/issues/$num" --jq .user.login)
	if [ "$who" = "layup-agent[bot]" ]; then echo "ok   row $row: #$num ($task) by $who"; else echo "FAIL row $row: #$num opened by $who"; fi
done > /tmp/issues.$$
cat /tmp/issues.$$
grep -q '^FAIL' /tmp/issues.$$ && fail=1
rm -f /tmp/issues.$$
exit "$fail"
