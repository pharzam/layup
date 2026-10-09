#!/bin/sh
#
# task-state.sh — print the state of each task, derived, never set by hand.
#
# Issue #175 (T-eep8). A session runs it before it takes a task, so that two
# sessions do not take the same one. It is a command, not a check: it enforces
# nothing, and it needs `gh` and the network to read the forge.
#
# Usage:
#   sh docs/tasks/task-state.sh             read the forge, then print
#   sh docs/tasks/task-state.sh --read DIR  print from DIR (the fixtures)
#
# DIR holds `backlog`, `completed` and `plan` (docs/tasks/backlog.md,
# docs/tasks/completed.md and docs/plan/README.md; the fetcher reads them from
# the forge's main, not from the checkout), `forge/branches` (one branch
# of the forge per line) and `forge/pr/<number>/head` and `.../body` (one open
# pull request each, its head branch and its body).
#
# The tasks: the rows of the LAST task table of the plan (the current
# milestone), in its order, then each task under `## Now` of the backlog that
# has no row there, in the backlog's order; such a task keeps the `After` of its
# row when an earlier task table holds one. A task table is one whose header
# holds the cells `#` and `Task ID`; its columns are found by the header, so a
# table of 13 columns and one of 14 read the same.
#
# The state, the first that holds:
#   done       the task ID is the ID field of a line of the completed log (the
#              second bold field of `- **date** — **T-xxxx** — …`; a mention
#              elsewhere in a line does not count)
#   in review  an open pull request has the head branch named the task ID, or
#              its body closes the task's issue: close, closes, closed, fix,
#              fixes, fixed, resolve, resolves or resolved, in any case, an
#              optional colon, then `#N` (the forge's set; R1 names three)
#   running    the forge has a branch named exactly the task ID (the claim rule
#              of docs/engineering-discipline.md, Starting a task)
#   blocked    a row of the task's `After` cell names a task that is not done
#   ready      none of the above
#
# The output, one line per task: TASK TAB #ISSUE TAB STATE TAB DETAIL. The
# detail is the smallest pull request for `in review`, the unfinished
# predecessors for `blocked` (and for `running`, as `after …`, so a branch
# started too early shows), else `—`. The same input gives the same bytes.
#
# Exit status: 0 = printed, 1 = an input is missing or malformed, or the forge
# cannot be read, 2 = a bad argument.

set -u

die() { printf 'task-state: %s\n' "$1" >&2; exit 1; }

tmp=$(mktemp -d) || die "cannot make a temporary directory"
trap 'rm -rf "$tmp"' EXIT

# fetch DIR — read the three files of the forge's main into DIR, and the
# branches and pull requests of the forge into DIR/forge. One progress line
# per step: each call to the forge can take seconds
# (docs/engineering-discipline.md, progress indicators).
fetch() {
	_in=$1
	root=$(git rev-parse --show-toplevel 2>/dev/null) || die "not inside a git checkout"
	mkdir -p "$_in/forge/pr"
	cd "$root" || die "cannot enter $root"
	# The three files come from the forge's main, not from this checkout, so a
	# checkout behind main gives the same answer (#182, note 3). The checkout
	# only names the repository, for the {owner}/{repo} of gh.
	printf 'task-state: reading the backlog, the completed log and the plan of main\n' >&2
	for _f in backlog:docs/tasks/backlog.md completed:docs/tasks/completed.md plan:docs/plan/README.md; do
		gh api -H 'Accept: application/vnd.github.raw' "repos/{owner}/{repo}/contents/${_f#*:}?ref=main" > "$_in/${_f%%:*}" ||
			die "gh: cannot read ${_f#*:} of main"
	done
	printf 'task-state: reading the branches of the forge\n' >&2
	gh api 'repos/{owner}/{repo}/branches' --paginate -q '.[].name' > "$_in/forge/branches" ||
		die "gh: cannot read the branches"
	printf 'task-state: reading the open pull requests\n' >&2
	gh pr list --state open --limit 1000 --json number,headRefName,body \
		-q '.[] | "\(.number)\t\(.headRefName)\t\(.body | @base64)"' > "$tmp/prlist" ||
		die "gh: cannot read the open pull requests"
	while IFS="$(printf '\t')" read -r _n _h _b; do
		mkdir "$_in/forge/pr/$_n" || die "cannot write pull request $_n"
		printf '%s\n' "$_h" > "$_in/forge/pr/$_n/head"
		printf '%s' "$_b" | base64 -d > "$_in/forge/pr/$_n/body" || die "cannot decode the body of pull request $_n"
	done < "$tmp/prlist"
}

case $# in
0)
	in=$tmp/in
	mkdir "$in" || die "cannot make $in"
	fetch "$in" ;;
2)
	[ "$1" = --read ] || { printf 'task-state: unknown argument: %s\n' "$1" >&2; exit 2; }
	in=$2
	[ -d "$in" ] || { printf 'task-state: not a directory: %s\n' "$in" >&2; exit 2; } ;;
*)
	printf 'task-state: usage: task-state.sh [--read DIR]\n' >&2
	exit 2 ;;
esac

for f in backlog completed plan forge/branches; do
	[ -f "$in/$f" ] || die "missing input: $f"
done

# Read a copy with each carriage return removed. On a checkout with
# `core.autocrlf=true` every input arrives with them, and `## Now` or a branch
# name would then match nothing: a wrong answer with exit 0 (round 1, finding 3).
mkdir "$tmp/n" && cp -R "$in/." "$tmp/n/" || die "cannot copy the input"
find "$tmp/n" -type f | while IFS= read -r f; do
	tr -d '\r' < "$f" > "$f.lf" && mv "$f.lf" "$f"
done
in=$tmp/n

# The open pull requests, one line each: NUMBER TAB HEAD TAB ISSUES, where
# ISSUES are the issues its body closes, comma-separated. A keyword counts only
# at the start of a word, and a number only when no letter or digit follows it:
# `prefixes #6` and `Closes #9x` close nothing.
: > "$tmp/prs"
for p in "$in"/forge/pr/*/; do
	[ -d "$p" ] || continue
	n=$(basename "$p")
	case $n in
	''|*[!0-9]*) die "a pull request directory is not a number: $n" ;;
	esac
	[ -f "$p/head" ] || die "pull request $n has no head"
	h=$(sed -n 1p "$p/head")
	c=''
	if [ -f "$p/body" ]; then
		c=$(awk '
		{
			s = tolower($0)
			while (match(s, /(^|[^a-z0-9_])(close|closes|closed|fix|fixes|fixed|resolve|resolves|resolved):?[ \t]+#[0-9]+/)) {
				m = substr(s, RSTART, RLENGTH)
				nx = substr(s, RSTART + RLENGTH, 1)
				s = substr(s, RSTART + RLENGTH)
				if (nx ~ /[a-z0-9_]/) continue
				sub(/.*#/, "", m)
				out = out (out == "" ? "" : ",") m
			}
		}
		END { print out }' "$p/body")
	fi
	printf '%s\t%s\t%s\n' "$n" "$h" "$c" >> "$tmp/prs"
done

# The rows of every task table: TABLE TAB ROW TAB TASK TAB ISSUE TAB AFTER.
awk -v OFS='\t' '
function trim(s) { gsub(/^[ \t]+|[ \t]+$/, "", s); return s }
function cells(line, a,    n, i) {
	sub(/^[ \t]*\|/, "", line); sub(/\|[ \t]*$/, "", line)
	n = split(line, a, "|")
	for (i = 1; i <= n; i++) a[i] = trim(a[i])
	return n
}
/^[ \t]*\|/ {
	if (!intable) {
		n = cells($0, h)
		ci = 0; ct = 0; cs = 0; ca = 0
		for (i = 1; i <= n; i++) {
			if (h[i] == "#") ci = i
			else if (h[i] == "Task ID") ct = i
			else if (h[i] == "Issue") cs = i
			else if (h[i] == "After") ca = i
		}
		if (ci && ct) {
			if (!cs) { print "task-state: a task table has no column Issue" > "/dev/stderr"; bad = 1; exit 1 }
			if (!ca) { print "task-state: a task table has no column After" > "/dev/stderr"; bad = 1; exit 1 }
			intable = 1; table++; width = n; sep = 1
		}
		next
	}
	if (sep) { sep = 0; next }
	n = cells($0, r)
	if (n != width) {
		printf "task-state: row %s of a task table has %d cells, its header %d\n", r[ci], n, width > "/dev/stderr"
		bad = 1; exit 1
	}
	t = r[ct]; gsub(/`/, "", t)
	if (t !~ /^T-[0-9a-z][0-9a-z][0-9a-z][0-9a-z]$/) {
		printf "task-state: row %s of a task table has no task ID\n", r[ci] > "/dev/stderr"
		bad = 1; exit 1
	}
	s = r[cs]
	if (!match(s, /#[0-9]+/)) {
		printf "task-state: row %s (%s) has no issue\n", r[ci], t > "/dev/stderr"
		bad = 1; exit 1
	}
	s = substr(s, RSTART + 1, RLENGTH - 1)
	a = r[ca]; gsub(/[ \t]/, "", a)
	print table, r[ci], t, s, a
	next
}
{ intable = 0; sep = 0 }
END {
	if (bad) exit 1
	if (!table) { print "task-state: the plan holds no task table" > "/dev/stderr"; exit 1 }
}' "$in/plan" > "$tmp/rows" || exit 1

# The tasks done: the ID field of each line of the completed log.
awk '
/^- \*\*[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]\*\* — \*\*T-[0-9a-z][0-9a-z][0-9a-z][0-9a-z]\*\*/ {
	s = $0; sub(/^- \*\*[^*]*\*\* — \*\*/, "", s); print substr(s, 1, 6)
}' "$in/completed" > "$tmp/done"

# The open tasks of the backlog: TASK TAB ISSUE, from the section `## Now`.
awk -v OFS='\t' '
/^## / { now = ($0 == "## Now"); next }
now && /^- \*\*T-[0-9a-z][0-9a-z][0-9a-z][0-9a-z]\*\*/ {
	t = substr($0, 5, 6)
	if (!match($0, /\[#[0-9]+\]/)) {
		printf "task-state: the backlog line of %s holds no issue\n", t > "/dev/stderr"
		exit 1
	}
	print t, substr($0, RSTART + 2, RLENGTH - 3)
}' "$in/backlog" > "$tmp/now" || exit 1

# The state of each task.
awk -F'\t' -v OFS='\t' -v dash='—' '
FILENAME == ARGV[1] { rowtask[$2] = $3; taskrow[$3] = $2; taskafter[$3] = $5; if ($1 > last) last = $1; nrows++; rt[nrows] = $1; rr[nrows] = $2; rk[nrows] = $3; ri[nrows] = $4; ra[nrows] = $5; next }
FILENAME == ARGV[2] { done[$1] = 1; next }
FILENAME == ARGV[3] { nnow++; nt[nnow] = $1; ni[nnow] = $2; next }
FILENAME == ARGV[4] { branch[$0] = 1; next }
FILENAME == ARGV[5] {
	if (!($2 in review) || $1 + 0 < review[$2] + 0) review[$2] = $1
	k = split($3, c, ",")
	for (i = 1; i <= k; i++) if (!(("#" c[i]) in review) || $1 + 0 < review["#" c[i]] + 0) review["#" c[i]] = $1
	next
}
function state(t, issue, after,    k, a, i, p, preds) {
	preds = ""
	if (after != "" && after != dash) {
		k = split(after, a, ",")
		for (i = 1; i <= k; i++) {
			if (!(a[i] in rowtask)) {
				printf "task-state: row %s (%s) is after row %s, which no task table holds\n", cur, t, a[i] > "/dev/stderr"
				bad = 1; exit 1
			}
			p = rowtask[a[i]]
			if (!(p in done)) preds = preds (preds == "" ? "" : ",") p
		}
	}
	if (t in done) return "done" OFS dash
	if (t in review && ("#" issue) in review) {
		return "in review" OFS "#" (review[t] + 0 < review["#" issue] + 0 ? review[t] : review["#" issue])
	}
	if (t in review) return "in review" OFS "#" review[t]
	if (("#" issue) in review) return "in review" OFS "#" review["#" issue]
	if (t in branch) return "running" OFS (preds == "" ? dash : "after " preds)
	if (preds != "") return "blocked" OFS preds
	return "ready" OFS dash
}
END {
	if (bad) exit 1
	for (i = 1; i <= nrows; i++) {
		if (rt[i] != last) continue
		cur = rr[i]; inlast[rk[i]] = 1
		line = rk[i] OFS "#" ri[i] OFS state(rk[i], ri[i], ra[i])
		if (bad) exit 1
		print line
	}
	for (i = 1; i <= nnow; i++) {
		if (nt[i] in inlast) continue
		# A task of an earlier task table keeps the After of its row (#182, note 1).
		cur = (nt[i] in taskrow) ? taskrow[nt[i]] : ""
		line = nt[i] OFS "#" ni[i] OFS state(nt[i], ni[i], (nt[i] in taskafter) ? taskafter[nt[i]] : "")
		if (bad) exit 1
		print line
	}
}' "$tmp/rows" "$tmp/done" "$tmp/now" "$in/forge/branches" "$tmp/prs" || exit 1
