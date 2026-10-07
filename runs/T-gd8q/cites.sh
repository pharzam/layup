#!/bin/sh
# Print the text of each line that numbers.md cites, as it stands at the base
# 34dd858, so that a reader can check that the line holds the value. A one-off
# evidence command of task T-gd8q, not a check of the gate. Exit 1 when a cited
# line does not exist.
#
#   sh runs/T-gd8q/cites.sh > runs/T-gd8q/cites.txt

base=34dd858
dir=$(dirname "$0")

# One "path line" pair per citation. A bare `:N` cites the last path named on its
# line; in the table of the pilot's defects (rows that start with "| F-") it cites
# runs/T-evad/findings.md.
awk '
	{
		last = ""
		if ($0 ~ /^\| F-/) last = "runs/T-evad/findings.md"
		s = $0
		while (match(s, /`[^`]+`/)) {
			t = substr(s, RSTART + 1, RLENGTH - 2)
			s = substr(s, RSTART + RLENGTH)
			if (t ~ /^[A-Za-z0-9_.\/-]+\.(md|sh|tsv):[0-9]+$/) {
				n = split(t, p, ":"); last = p[1]; print p[1], p[2]
			} else if (t ~ /^:[0-9]+$/ && last != "") {
				print last, substr(t, 2)
			}
		}
	}' "$dir/numbers.md" | sort -u -k1,1 -k2,2n |
while read -r path line; do
	text=$(git show "$base:$path" 2>/dev/null | sed -n "${line}p")
	if [ -z "$text" ]; then
		printf '%s:%s: NO SUCH LINE\n' "$path" "$line"
		echo x >> "${TMPDIR:-/tmp}/cites.$$"
	else
		printf '%s:%s: %.300s\n' "$path" "$line" "$text"
	fi
done
if [ -f "${TMPDIR:-/tmp}/cites.$$" ]; then rm -f "${TMPDIR:-/tmp}/cites.$$"; exit 1; fi
exit 0
