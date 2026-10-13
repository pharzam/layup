#!/bin/sh
#
# task-state.sh — print the state of each task, derived, never set by hand.
#
# Issue #175 (T-eep8). A session runs it before it takes a task, so that two
# sessions do not take the same one. It is a command, not a check: it enforces
# nothing, and it needs `gh` and the network to read the forge.
#
# Usage:
#   sh docs/tasks/task-state.sh                       read the forge, then print
#   sh docs/tasks/task-state.sh --read DIR            print from DIR (the fixtures)
#   sh docs/tasks/task-state.sh [--read DIR] --prs FILE   also write to FILE the
#       table of the open pull requests, NUMBER TAB HEAD TAB ISSUES, where ISSUES
#       are the issues that its body closes, and a pull request that closes none
#       has no third field (the one parser of the closing keywords: task-labels.sh
#       reads this table, #303). Give FILE as an absolute path: when the forge is
#       read the command works from the repository root, so a relative FILE
#       lands there
#
# DIR holds `backlog`, `completed` and `plan` (docs/tasks/backlog.md,
# docs/tasks/completed.md and docs/plan/README.md; the fetcher reads them from
# the forge's main, not from the checkout), `forge/branches` (one branch
# of the forge per line), `forge/pr/<number>/head` and `.../body` (one open
# pull request each, its head branch and its body), `now` (the time of the run,
# in seconds since the epoch), `forge/pushes` (`BRANCH TAB SECONDS`, a push, a
# force push or a creation of a branch) and `forge/comments` (`ISSUE TAB
# SECONDS`, a comment created on an issue) and `forge/issues` (`NUMBER TAB
# TITLE`, each open issue).
#
# The tasks: the rows of the LAST task table of the plan (the current
# milestone), in its order, then each task under `## Now` of the backlog that
# has no row there, in the backlog's order; such a task keeps the `After` of its
# row when an earlier task table holds one. Last come the tasks that only an
# open issue names (#303): its title starts with a task ID and a colon
# (`T-xxxx: `, as the issues of the plan do), the task has no `After`, and the
# issues come in the order of their numbers. A task or an issue that a row of
# the plan or a line under `## Now` names makes no second task, so a title that
# carries the ID of such a task, or sits on an issue that one of them names,
# adds none; of two open issues that start with one ID, the lower number is the
# task. A task that only the completed log names is `done`, so an open issue
# whose title starts with its ID is listed `done` and takes `state: merged`
# while it is open, as a reopened issue of a done task of the plan does.
# A task table is one whose header holds the cells `#` and `Task ID`; its
# columns are found by the header, so a table of 13 columns and one of 14 read
# the same.
#
# The state, the first that holds:
#   done       the task ID is the ID field of a line of the completed log (the
#              second bold field of `- **date** — **T-xxxx** — …`; a mention
#              elsewhere in a line does not count)
#   in review  an open pull request has the head branch named the task ID, or
#              its body closes the task's issue: close, closes, closed, fix,
#              fixes, fixed, resolve, resolves or resolved, in any case, an
#              optional colon, then `#N` (the forge's set; R1 names three)
#   in progress  the forge has a branch named exactly the task ID (the claim of
#              docs/engineering-discipline.md, Starting a task), and the task's
#              activity, the latest push to that branch or comment created on
#              its issue, is at most `window` seconds before `now` (#202)
#   stalled    the branch is there, and no such activity is
#   blocked    a row of the task's `After` cell names a task that is not done
#   ready      none of the above
#
# The output, one line per task: TASK TAB #ISSUE TAB STATE TAB DETAIL. The
# detail is the smallest pull request for `in review`, the unfinished
# predecessors for `blocked` (and for `in progress` and `stalled`, as
# `after …`, so a branch started too early shows), else `—`. The same input gives the same bytes.
#
# Exit status: 0 = printed, 1 = an input is missing or malformed, or the forge
# cannot be read, 2 = a bad argument.

set -u

# The window of activity, in seconds: one value, which the reader cuts by and
# the fetcher's two bounds follow (#202).
window=86400

die() { printf 'task-state: %s\n' "$1" >&2; exit 1; }

tmp=$(mktemp -d) || die "cannot make a temporary directory"
trap 'rm -rf "$tmp"' EXIT

# fetch DIR — read the three files of the forge's main into DIR, the clock
# into DIR/now, and the branches, the pull requests, the pushes and the
# comments of the window of the forge into DIR/forge. One progress line
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
	printf 'task-state: reading the open issues\n' >&2
	gh api 'repos/{owner}/{repo}/issues?state=open&per_page=100' --paginate \
		-q '.[] | select(.pull_request | not) | "\(.number)\t\(.title | gsub("[\t\n\r]"; " "))"' > "$_in/forge/issues" ||
		die "gh: cannot read the open issues"
	# The activity of the window. `date +%s` is the clock; gh's jq converts the
	# times (fromdateiso8601, todate), as no option of `date` that does is
	# portable. The forge's period is the smallest of its periods (day, week,
	# month, quarter, year) that covers the window, and
	# the comments' `since` (a time of update, a superset) is now less the
	# window; the reader cuts both by the same value.
	_now=$(date -u +%s) || die "cannot read the clock"
	printf '%s\n' "$_now" > "$_in/now"
	if [ "$window" -le 86400 ]; then _p=day; elif [ "$window" -le 604800 ]; then _p=week
	elif [ "$window" -le 2592000 ]; then _p=month; elif [ "$window" -le 7776000 ]; then _p=quarter
	else _p=year; fi
	printf 'task-state: reading the pushes of the last %s\n' "$_p" >&2
	gh api "repos/{owner}/{repo}/activity?time_period=$_p&per_page=100" --paginate \
		-q '.[] | select(.activity_type == "push" or .activity_type == "force_push" or .activity_type == "branch_creation") | select(.ref | startswith("refs/heads/")) | "\(.ref | ltrimstr("refs/heads/"))\t\(.timestamp | fromdateiso8601)"' \
		> "$_in/forge/pushes" || die "gh: cannot read the activity"
	_since=$(gh api rate_limit -q "$((_now - window)) | todate") || die "gh: cannot convert the time"
	printf 'task-state: reading the comments since %s\n' "$_since" >&2
	gh api "repos/{owner}/{repo}/issues/comments?since=$_since&per_page=100" --paginate \
		-q '.[] | "\(.issue_url | split("/") | last)\t\(.created_at | fromdateiso8601)"' \
		> "$_in/forge/comments" || die "gh: cannot read the comments"
}

in=
prs_out=
while [ $# -gt 0 ]; do
	case $1 in
	--read | --prs)
		[ $# -ge 2 ] || { printf 'task-state: usage: task-state.sh [--read DIR] [--prs FILE]\n' >&2; exit 2; }
		if [ "$1" = --read ]; then in=$2; else prs_out=$2; fi
		shift 2 ;;
	*) printf 'task-state: unknown argument: %s\n' "$1" >&2; exit 2 ;;
	esac
done
if [ -z "$in" ]; then
	in=$tmp/in
	mkdir "$in" || die "cannot make $in"
	fetch "$in"
else
	[ -d "$in" ] || { printf 'task-state: not a directory: %s\n' "$in" >&2; exit 2; }
fi

for f in backlog completed plan now forge/branches forge/pushes forge/comments forge/issues; do
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
_now=$(sed -n 1p "$in/now")
case $_now in
''|*[!0-9]*) die "now is not a number of seconds: $_now" ;;
esac

# The open pull requests, one line each: NUMBER TAB HEAD TAB ISSUES, where
# ISSUES are the issues its body closes, comma-separated, and no tab before an
# empty ISSUES (no line ends in a tab). A keyword counts only
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
	if [ -n "$c" ]; then printf '%s\t%s\t%s\n' "$n" "$h" "$c"; else printf '%s\t%s\n' "$n" "$h"; fi >> "$tmp/prs"
done

[ -z "$prs_out" ] || cp "$tmp/prs" "$prs_out" || die "cannot write $prs_out"

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

# The tasks that only an open issue names (#303): TASK TAB ISSUE, in the order of
# the issue numbers, from a title that starts with a task ID and a colon. An
# issue or a task that the plan or the backlog names makes none, and of two issues
# with one ID the lower number is the task.
awk -F'\t' -v OFS='\t' '
FILENAME == ARGV[1] { named[$3] = 1; named["#" $4] = 1; next }
FILENAME == ARGV[2] { named[$1] = 1; named["#" $2] = 1; next }
$2 ~ /^T-[0-9a-z][0-9a-z][0-9a-z][0-9a-z]: / {
	id = substr($2, 1, 6)
	if (id in named || ("#" $1) in named) next
	if (!(id in low) || $1 + 0 < low[id] + 0) low[id] = $1
}
END { for (id in low) print low[id] + 0, id }' "$tmp/rows" "$tmp/now" "$in/forge/issues" | sort -n | awk -F' ' -v OFS='\t' '{ print $2, $1 }' >> "$tmp/now" || exit 1

# The state of each task.
awk -F'\t' -v OFS='\t' -v dash='—' -v window="$window" '
FILENAME == ARGV[1] { rowtask[$2] = $3; taskrow[$3] = $2; taskafter[$3] = $5; if ($1 > last) last = $1; nrows++; rt[nrows] = $1; rr[nrows] = $2; rk[nrows] = $3; ri[nrows] = $4; ra[nrows] = $5; next }
FILENAME == ARGV[2] { done[$1] = 1; next }
FILENAME == ARGV[3] { nnow++; nt[nnow] = $1; ni[nnow] = $2; next }
FILENAME == ARGV[4] { branch[$0] = 1; next }
FILENAME == ARGV[6] { now = $1 + 0; next }
FILENAME == ARGV[7] { if (!($1 in pushed) || $2 + 0 > pushed[$1]) pushed[$1] = $2 + 0; next }
FILENAME == ARGV[8] { if (!($1 in commented) || $2 + 0 > commented[$1]) commented[$1] = $2 + 0; next }
FILENAME == ARGV[5] {
	if (!($2 in review) || $1 + 0 < review[$2] + 0) review[$2] = $1
	k = split($3, c, ",")
	for (i = 1; i <= k; i++) if (!(("#" c[i]) in review) || $1 + 0 < review["#" c[i]] + 0) review["#" c[i]] = $1
	next
}
function state(t, issue, after,    k, a, i, p, preds, act) {
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
	if (t in branch) {
		act = -1
		if (t in pushed) act = pushed[t]
		if (issue in commented && commented[issue] > act) act = commented[issue]
		return (act >= 0 && act >= now - window ? "in progress" : "stalled") OFS (preds == "" ? dash : "after " preds)
	}
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
}' "$tmp/rows" "$tmp/done" "$tmp/now" "$in/forge/branches" "$tmp/prs" "$in/now" "$in/forge/pushes" "$in/forge/comments" || exit 1
