#!/bin/sh
# The check of S2 of task T-ywk7 (#147): each of the 11 rows of the section
# `M2b` of runs/T-w73g/survey.md has exactly one decision in
# runs/T-ywk7/survey-m2b.tsv (row, pattern, decision, reason, file, heading).
# A one-off evidence command of this task, not a check of the gate.
#
# A row of the table matches the survey's row of the same number when its
# `pattern` is the first four words of that row's Pattern cell. Its decision is
# `take` or `reject`, with a reason. A `take` names the file and the heading of
# docs/spec/ that applies it, and the heading is a whole line of that file; a
# `reject` has `—` in both. Run from the root of the repository; exit 0 when
# each check passes, 1 otherwise.
survey=runs/T-w73g/survey.md
tab=runs/T-ywk7/survey-m2b.tsv
out=$(mktemp) || exit 1
trap 'rm -f "$out"' EXIT
if [ ! -f "$tab" ]; then
	echo "FAIL the table $tab does not exist"
	exit 1
fi
# The first four words of each Pattern cell of the section `M2b`, in order.
awk -F'|' '
	/^### `M2b`/ { in_m2b = 1; next }
	/^### / { in_m2b = 0 }
	in_m2b && /^\| / && $2 !~ /^ *Pattern *$/ && $2 !~ /^ *-+ *$/ {
		n = split($2, w, " ")
		printf "%s %s %s %s\n", w[1], w[2], w[3], w[4]
	}
' "$survey" > "$out"
want=$(wc -l < "$out")
if [ "$want" -ne 11 ]; then
	echo "FAIL the section M2b of $survey has $want rows, not 11"
	exit 1
fi
fail=0
rows=$(tail -n +2 "$tab" | wc -l)
if [ "$rows" -ne 11 ]; then
	echo "FAIL $tab has $rows rows, not 11"
	fail=1
fi
i=0
while IFS= read -r words || [ -n "$words" ]; do
	i=$((i + 1))
	line=$(awk -F'\t' -v r="$i" 'NR > 1 && $1 == r' "$tab")
	count=$(awk -F'\t' -v r="$i" 'NR > 1 && $1 == r' "$tab" | wc -l)
	if [ "$count" -ne 1 ]; then
		echo "FAIL row $i ($words): $count decisions, not 1"
		fail=1
		continue
	fi
	pattern=$(printf '%s\n' "$line" | cut -f2)
	decision=$(printf '%s\n' "$line" | cut -f3)
	reason=$(printf '%s\n' "$line" | cut -f4)
	file=$(printf '%s\n' "$line" | cut -f5)
	heading=$(printf '%s\n' "$line" | cut -f6)
	if [ "$pattern" != "$words" ]; then
		echo "FAIL row $i: the pattern is \"$pattern\", the survey's row is \"$words\""
		fail=1
	fi
	if [ -z "$reason" ] || [ "$reason" = "—" ]; then
		echo "FAIL row $i ($words): no reason"
		fail=1
	fi
	case $decision in
	take)
		if [ -f "$file" ] && grep -Fqx -- "$heading" "$file"; then
			echo "ok   row $i ($words): take, $file: $heading"
		else
			echo "FAIL row $i ($words): take, but $file has no line \"$heading\""
			fail=1
		fi
		;;
	reject)
		if [ "$file" = "—" ] && [ "$heading" = "—" ]; then
			echo "ok   row $i ($words): reject"
		else
			echo "FAIL row $i ($words): reject, with a file or a heading"
			fail=1
		fi
		;;
	*)
		echo "FAIL row $i ($words): the decision \"$decision\" is not take or reject"
		fail=1
		;;
	esac
done < "$out"
exit "$fail"
