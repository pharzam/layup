#!/bin/sh
# rules-diff.sh: the evidence of REQ-018 for a target that `layup setup` made (task T-evad, #97:
# D8 (b) of its plan and condition 4 of its plan review).
#
# REQ-018: "A target's baseline rules are byte-identical to the pinned baseline's after setup, except
# the adapted values the setup records with evidence."
#
# Usage: sh rules-diff.sh TARGET HEAD RECORD [SETUP-CHECK]
#   TARGET       the git repository of the target (WORK/target); its root commit is the unmodified baseline
#   HEAD         the ref of the setup head (layup-setup)
#   RECORD       WORK/out/record.tsv (columns step, name, value, source, ref)
#   SETUP-CHECK  LAYUP's docs/setup/setup-check.sh, which gives the lines that check `adapted` flags
#                (default: the one of the repository that holds this script)
#
# It compares the union of the paths of the root commit and of HEAD, and fails (exit 1, one `FAIL:` line
# for each problem) unless each difference is accounted for:
#   - each commit after the root is a setup step ("chore: setup Sxx");
#   - a deleted path is a path of the history that step S05 removes (docs/decisions/, docs/audit/,
#     docs/tasks/T-*.md); any other deleted baseline path is a deleted rule;
#   - an added path is listed with the step that added it (it is not a baseline rule);
#   - a task index (docs/tasks/backlog.md, completed.md) changed only by removing lines (S05);
#   - a file that the prose step writes for the target (README.md, AGENTS.md, docs/onboarding-for-engineers.md,
#     docs/glossary.md, docs/guardrails.md) has a record row, and is listed for the audit;
#   - in each other changed file, each changed or removed baseline line is a line that check `adapted`
#     flagged at the root (the kit voice that S14 replaces), or is inside a whole section whose heading is
#     flagged, or holds a marker that the record fills (with its value in the file); lines added to a file
#     are allowed only in an index file (docs/adr/README.md, docs/facts/README.md);
#   - a file with a changed baseline line that is not a filled marker has a record row `file:<path>` whose
#     value is the SHA-256 of the file at HEAD. A record row never makes a changed line allowed.
# Exit 0 and `rules-diff: PASS` only when all hold. The output lists each class for the audit.

set -u
T=${1:?usage: sh rules-diff.sh TARGET HEAD RECORD [SETUP-CHECK]}
HEADREV=${2:?HEAD}
REC=${3:?RECORD}
here=$(cd "$(dirname "$0")" && pwd)
SC=${4:-$here/../../docs/setup/setup-check.sh}
[ -d "$T" ] && [ -f "$REC" ] && [ -f "$SC" ] || { echo "rules-diff: input error: TARGET, RECORD or SETUP-CHECK is missing" >&2; exit 2; }

tmp=$(mktemp -d) || exit 2
trap 'rm -rf "$tmp"' EXIT INT TERM
LQ=$(printf '\342\200\271')   # the marker characters, never typed in this file
export LQ
fail=0
bad() { echo "FAIL: $*"; fail=$((fail + 1)); }
gt() { git -C "$T" -c core.quotepath=off -c core.hooksPath=/dev/null -c maintenance.auto=false "$@"; }
sha256() { if command -v sha256sum >/dev/null 2>&1; then sha256sum "$1" | cut -d' ' -f1; else shasum -a 256 "$1" | cut -d' ' -f1; fi; }

WRITTEN="README.md AGENTS.md docs/onboarding-for-engineers.md docs/glossary.md docs/guardrails.md"
S05INDEX="docs/tasks/backlog.md docs/tasks/completed.md"
INDEXAPP="docs/adr/README.md docs/facts/README.md"
in_list() { case " $2 " in *" $1 "*) return 0 ;; esac; return 1; }

root=$(gt rev-list --max-parents=0 "$HEADREV" | tail -1)
[ -n "$root" ] || { echo "rules-diff: input error: $HEADREV has no root commit" >&2; exit 2; }
total=$(gt ls-tree -r --name-only "$root" | grep -c .)

# 1. Each commit after the root is a setup step.
gt log --format='%h %s' "$root..$HEADREV" > "$tmp/log"
steps=$(grep -c . "$tmp/log")
foreign=$(grep -v -E '^[0-9a-f]+ chore: setup S[0-9][0-9]$' "$tmp/log")
if [ -n "$foreign" ]; then
	printf '%s\n' "$foreign" | while IFS= read -r l; do echo "FAIL: a commit that no step made: $l"; done
	fail=$((fail + $(printf '%s\n' "$foreign" | grep -c .)))
fi

# 2. The lines that check `adapted` flags at the root commit.
git clone -q --no-hardlinks "$T" "$tmp/root" 2>/dev/null && git -C "$tmp/root" -c maintenance.auto=false checkout -q "$root" 2>/dev/null \
	|| { echo "rules-diff: input error: cannot check out the root commit" >&2; exit 2; }
sh "$SC" --only adapted "$tmp/root" > "$tmp/adapted.txt" 2>&1
if grep -q 'cannot list' "$tmp/adapted.txt"; then bad "check adapted could not run on the root commit"; fi
sed -n 's/^setup-check: adapted FAIL [^ ]* .*: \(.*\):\([0-9][0-9]*\)$/\1	\2/p' "$tmp/adapted.txt" | sort -u > "$tmp/flagged.tsv"

# 3. The record: the rows of the files that the prose step copied, and of the markers that S11 filled.
awk -F'\t' '$2 ~ /^file:/ { print substr($2, 6) "\t" $3 "\t" $1 }' "$REC" > "$tmp/filerows.tsv"
awk -F'\t' '$2 ~ /^marker:/ { n = $2; sub(/^marker:/, "", n); sub(/(:[0-9]+)+$/, "", n); print n "\t" $3 }' "$REC" > "$tmp/markerrows.tsv"

# 4. The paths of the union, by status.
gt diff --name-status --no-renames "$root" "$HEADREV" > "$tmp/ns"
n_del=0; n_add=0; n_mod=0; n_adapted=0; n_written=0; n_marker=0; n_index=0; n_s05=0

while IFS='	' read -r st p; do
	case "$p" in
	\"*|*\\*) bad "a path that this check cannot read: $p"; continue ;;
	esac
	case "$st" in
	D)
		n_del=$((n_del + 1))
		if printf '%s\n' "$p" | grep -E -q '^(docs/(decisions|audit)/|docs/tasks/T-[^/]*\.md$)'; then
			echo "DELETED $p: the history that step S05 removes"
		else
			bad "deleted baseline path: $p (no rule of step S05 removes it)"
		fi ;;
	A)
		n_add=$((n_add + 1))
		s=$(gt log --diff-filter=A --format=%s "$root..$HEADREV" -- "$p" | tail -1)
		echo "ADDED $p: ${s:-no commit}" ;;
	M)
		n_mod=$((n_mod + 1))
		gt show "$root:$p" > "$tmp/base" 2>/dev/null; gt show "$HEADREV:$p" > "$tmp/new" 2>/dev/null
		# the hunks, one line "OP FIRST LAST" for each, of the baseline lines that are changed (c), removed (d) or only
		# followed by added lines (a); git diff -U0 keeps an unchanged line between two changes out of both
		gt diff -U0 --minimal --no-renames --no-color "$root" "$HEADREV" -- "$p" | awk '
			/^@@ / {
				l = $2; r = $3; sub(/^-/, "", l); sub(/^\+/, "", r)
				split(l, a, ","); split(r, b, ",")
				ln = (a[2] == "" ? 1 : a[2] + 0); rn = (b[2] == "" ? 1 : b[2] + 0)
				op = (ln == 0 ? "a" : (rn == 0 ? "d" : "c"))
				print op, a[1] + 0, a[1] + (ln > 0 ? ln - 1 : 0), ln, rn
				nrem += ln; nadd += rn
			}
			END { print "TOTAL", nrem + 0, nadd + 0 }' > "$tmp/d"
		nremoved=$(awk '$1 == "TOTAL" { print $2 }' "$tmp/d"); nnew=$(awk '$1 == "TOTAL" { print $3 }' "$tmp/d")
		P=$p IDXAPP=0 W=0 S05=0; in_list "$p" "$INDEXAPP" && IDXAPP=1; in_list "$p" "$WRITTEN" && W=1; in_list "$p" "$S05INDEX" && S05=1
		export P IDXAPP W S05
		LC_ALL=C awk -v BASE="$tmp/base" -v FL="$tmp/flagged.tsv" -v MK="$tmp/markerrows.tsv" -v DIFF="$tmp/d" '
			function lvl(s) { if (match(s, /^#+ /)) return RLENGTH - 1; return 0 }
			function secend(a,   l, e) { l = lvl(B[a]); e = a; while (e < N && !(lvl(B[e + 1]) > 0 && lvl(B[e + 1]) <= l)) e++; return e }
			BEGIN {
				P = ENVIRON["P"]; LQ = ENVIRON["LQ"]; idxapp = ENVIRON["IDXAPP"] + 0; w = ENVIRON["W"] + 0; s05 = ENVIRON["S05"] + 0
				while ((getline line < BASE) > 0) B[++N] = line
				while ((getline line < FL) > 0) { split(line, a, "\t"); if (a[1] == P) F[a[2] + 0] = 1 }
				while ((getline line < MK) > 0) { split(line, a, "\t"); if (a[1] == P) nmk++ }
				# the baseline lines that the diff changes or removes
				while ((getline line < DIFF) > 0) {
					split(line, h, " ")
					if (h[1] == "c" || h[1] == "d") for (n = h[2] + 0; n <= h[3] + 0; n++) C[n] = 1
				}
				# a flagged heading whose whole section (each line that is not blank) is changed or removed is a
				# section replaced as a unit, however the diff splits it into hunks; a heading that is only renamed is not
				for (hd in F) {
					hd += 0
					if (lvl(B[hd]) == 0 || !C[hd]) continue
					e = secend(hd); ok = 1
					for (n = hd; n <= e; n++) if (B[n] != "" && !C[n]) ok = 0
					if (ok) { for (n = hd; n <= e; n++) SEC[n] = 1; sections = sections " " hd "-" e }
				}
			}
			$1 == "TOTAL" { next }
			{
				op = $1; l1 = $2 + 0; l2 = $3 + 0; ln = $4 + 0; rn = $5 + 0
				# a task index (S05) loses lines, or has a marker line replaced by one line: nothing else
				if (s05 && (op == "a" || (op == "c" && rn != ln))) { print "FAIL: task index changed other than by removing lines: " P; bads++; next }
				if (op == "a") { added++; if (!idxapp && !w && !said_a) { print "FAIL: lines added to a baseline rule file: " P; said_a = 1; bads++ } ; next }
				for (n = l1; n <= l2; n++) {
					if (s05 && op == "d") { nrem++ }
					else if (F[n]) { nflag++ }
					else if (SEC[n]) { nsec++ }
					else if (idxapp && B[n] ~ /^\| _none yet_/) { nidx++ }
					else if (index(B[n], LQ) > 0 && nmk > 0) { nmark++ }
					else if (w) { nwr++ }
					else { print "FAIL: unflagged baseline line changed: " P ":" n ": " substr(B[n], 1, 90); bads++ ; nother++ }
				}
			}
			END { print "RESULT placeholder=" nidx + 0 " flagged=" nflag + 0 " section=" nsec + 0 " marker=" nmark + 0 " removed=" nrem + 0 " unflagged=" nother + 0 " added=" added + 0 " bads=" bads + 0 " sections=" sections }
		' "$tmp/d" > "$tmp/awk.out"
		grep '^FAIL:' "$tmp/awk.out" && fail=$((fail + $(grep -c '^FAIL:' "$tmp/awk.out")))
		res=$(grep '^RESULT' "$tmp/awk.out")
		field() { printf '%s' "$res" | sed -n "s/.* $1=\([0-9]*\).*/\1/p; s/^RESULT $1=\([0-9]*\).*/\1/p" | head -1; }
		nfl=$(field flagged); nsec=$(field section); nmk=$(field marker); nph=$(field placeholder); nadd=$(field added); nrm=$(field removed)
		secs=$(printf '%s' "$res" | sed -n 's/.* sections=\(.*\)$/\1/p')
		# the markers that the record fills: each recorded value is in the file
		if [ "$nmk" -gt 0 ]; then
			n_marker=$((n_marker + 1))
			awk -F'\t' -v p="$p" '$1 == p { print $2 }' "$tmp/markerrows.tsv" | while IFS= read -r v; do
				grep -F -q -- "$v" "$tmp/new" || echo "FAIL: recorded marker value is not in the file: $p ($v)"
			done > "$tmp/mk.out"
			grep '^FAIL:' "$tmp/mk.out" && fail=$((fail + $(grep -c '^FAIL:' "$tmp/mk.out")))
			echo "MARKERS $p: $nmk marker lines filled; each recorded value is in the file"
		fi
		need_row=0
		if [ "$S05" = 1 ]; then
			n_s05=$((n_s05 + 1)); echo "S05INDEX $p: $nrm lines removed"
		elif [ "$W" = 1 ]; then
			n_written=$((n_written + 1)); need_row=1
			echo "WRITTEN $p: written for the target by the prose step; $nremoved baseline lines changed or removed, $nnew lines new (listed for the audit)"
		elif [ "$nfl" -gt 0 ] || [ "$nsec" -gt 0 ]; then
			n_adapted=$((n_adapted + 1)); need_row=1
			echo "ADAPTED $p: $nfl flagged lines changed${secs:+, whole sections$secs replaced or removed}${nadd:+, $nadd line groups added}"
		elif [ "$nadd" -gt 0 ] || [ "$nph" -gt 0 ]; then
			n_index=$((n_index + 1))
			echo "INDEX $p: ${nadd:-0} line groups added and $nph placeholder rows replaced by the steps; no other baseline line changed"
		fi
		# a file that the prose step copied has a record row whose value is the SHA-256 of the file at the commit of the
		# step of the row (S11 fills the markers of the file later, and the markers are checked above)
		if [ "$need_row" = 1 ]; then
			rowline=$(awk -F'\t' -v p="$p" '$1 == p { print $2 "\t" $3; exit }' "$tmp/filerows.tsv")
			if [ -z "$rowline" ]; then
				bad "no record row for $p"
			else
				rhash=${rowline%%	*}; rstep=${rowline#*	}
				sc=$(gt log --format=%H --grep="^chore: setup $rstep\$" "$root..$HEADREV" | tail -1)
				if [ -z "$sc" ] || ! gt show "$sc:$p" > "$tmp/atstep" 2>/dev/null || [ "$rhash" != "$(sha256 "$tmp/atstep")" ]; then
					bad "record row differs from the file at its step: $p"
				fi
			fi
		fi ;;
	*) bad "a status that this check does not read: $st $p" ;;
	esac
done < "$tmp/ns"

same=$((total - n_del - n_mod))
echo "rules-diff: $total baseline paths: $same byte-identical, $n_mod changed ($n_adapted adapted, $n_written written, $n_marker with filled markers, $n_s05 task indexes, $n_index index files with added rows), $n_del deleted, $n_add added; $steps setup commits"
if [ "$fail" -eq 0 ]; then echo "rules-diff: PASS"; exit 0; fi
echo "rules-diff: FAIL ($fail problems)"
exit 1
