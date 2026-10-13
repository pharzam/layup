#!/bin/sh
#
# task-labels.sh — set one `state:` label on each open issue, and on each closed
# issue that is a task issue or still carries a `state:` label (#185, T-dm8m;
# every open issue, #303, T-j304).
#
# The workflow .github/workflows/task-state.yml runs it, and is the one writer
# of these labels; a label set by hand is replaced at the next run, and so is a
# color set by hand. Each label's color is the pair in the one list of apply():
# ready blue (1d76db), in progress amber (fbca04), in review purple (5319e7),
# stalled orange (d93f0b), blocked dark red (b60205), merged green (0e8a16, the
# `task` green), rejected slate (57606a), on hold ice (c5def5).
#
# Usage:
#   sh docs/tasks/task-labels.sh                   read, then write the labels
#   sh docs/tasks/task-labels.sh --dry-run         read, then print the changes
#   sh docs/tasks/task-labels.sh --read DIR        print the changes from DIR
#   sh docs/tasks/task-labels.sh --check-workflow FILE   check the workflow
#
# DIR holds `state` (the output of task-state.sh), `prs` (the table of the open
# pull requests that `task-state.sh --prs` writes: NUMBER TAB HEAD TAB ISSUES,
# the issues that the body closes, comma-separated, and no third field when it
# closes none), `plan`, `completed` and
# `backlog` (docs/plan/README.md, docs/tasks/completed.md and
# docs/tasks/backlog.md of main), and `issues` (`NUMBER TAB open|closed TAB
# STATE_REASON TAB LABELS`, the labels comma-separated).
#
# THE RULE. This header is the one home of who gets which label; the documents
# that name the labels summarise it and link here.
#
# A task issue is one that a file or a title names: the `Issue` of each row of
# every task table of the plan, the first `[#N]` of each line of the completed
# log, and of each line under `## Now` of the backlog (the Operator's choice of
# 2026-10-12: a task rejected before it was ever indexed stays unlabelled), and
# an open issue whose title starts with a task ID and a colon, which
# task-state.sh lists (#303). The label of an issue, the first that holds:
#   closed   `state: merged` if its state_reason is `completed`, else
#            `state: rejected` (not_planned, duplicate, empty); for a task issue,
#            and for any other closed issue that still carries a `state:` label
#            (the label of its open time, now stale). Any other closed issue is
#            left as it is.
#   open     the state task-state.sh lists for it: done -> `state: merged`,
#            `in review`, `in progress`, `stalled`, `blocked`, `ready`
#   open     else `state: in review` when an open pull request closes it (the
#            table `prs` names it), else `state: on hold`
# So each open issue has one `state:` label. Every other `state:` label of an
# issue that it labels is removed. Labels outside `state:` are never touched.
#
# The output: `add #N LABEL` and `remove #N LABEL`, in the order of the issue
# numbers, the removals of an issue before its addition; nothing for an issue
# already right. The same input gives the same bytes.
#
# No write before every read has passed: task-state.sh (the states and the
# table of the pull requests), the three files, the issues, and the reader's
# checks. An empty set of task issues, or an issue read that holds none of them,
# exits 1 and writes nothing, so a broken read never strips the labels.
# `--dry-run` writes nothing at all.
#
# LIMIT: a read that succeeds and returns too little is not seen. An empty table
# of pull requests turns each `in review` of an issue outside the plan into
# `on hold`, as an empty list of branches turns each `in progress` of a task into
# `ready`; the next run sets them again.
#
# TASK_STATE_CMD, when set, names the script run in place of task-state.sh; the
# fixtures use it with a stand-in. In use it is not set.
#
# Exit status: 0 = done, 1 = an input is missing or malformed, a read failed, or
# the workflow fails its check, 2 = a bad argument.

set -u

die() { printf 'task-labels: %s\n' "$1" >&2; exit 1; }
say() { printf 'task-labels: %s\n' "$1" >&2; }

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
tmp=$(mktemp -d) || die "cannot make a temporary directory"
trap 'rm -rf "$tmp"' EXIT

# read_dir DIR — print the changes; exit 1 on a missing or malformed input.
read_dir() {
	for f in state plan completed backlog issues prs; do
		[ -f "$1/$f" ] || die "missing input: $f"
		tr -d '\r' < "$1/$f" > "$tmp/r.$f"
	done
	awk -F'\t' -v OFS='\t' '
	function trim(s) { gsub(/^[ \t]+|[ \t]+$/, "", s); return s }
	FILENAME ~ /r\.plan$/ {
		if ($0 !~ /^[ \t]*\|/) { intable = 0; sep = 0; next }
		if (!intable) {
			n = split($0, h, "|"); ci = 0; ct = 0; cs = 0
			for (i = 1; i <= n; i++) { c = trim(h[i]); if (c == "#") ci = i; else if (c == "Task ID") ct = i; else if (c == "Issue") cs = i }
			if (ci && ct && cs) { intable = 1; sep = 1 }
			next
		}
		if (sep) { sep = 0; next }
		split($0, r, "|")
		if (match(r[cs], /#[0-9]+/)) set[substr(r[cs], RSTART + 1, RLENGTH - 1) + 0] = 1
		next
	}
	FILENAME ~ /r\.completed$/ {
		if ($0 ~ /^- \*\*[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]\*\* — \*\*T-/ && match($0, /\[#[0-9]+\]/))
			set[substr($0, RSTART + 2, RLENGTH - 3) + 0] = 1
		next
	}
	FILENAME ~ /r\.backlog$/ {
		if ($0 ~ /^## /) { now = ($0 == "## Now"); next }
		if (now && $0 ~ /^- \*\*T-/ && match($0, /\[#[0-9]+\]/)) set[substr($0, RSTART + 2, RLENGTH - 3) + 0] = 1
		next
	}
	FILENAME ~ /r\.state$/ {
		if ($2 ~ /^#[0-9]+$/) st[substr($2, 2) + 0] = $3
		next
	}
	FILENAME ~ /r\.prs$/ {
		k = split($3, closes, ",")
		for (i = 1; i <= k; i++) if (closes[i] ~ /^[0-9]+$/) pr[closes[i] + 0] = 1
		next
	}
	FILENAME ~ /r\.issues$/ {
		if ($1 !~ /^[0-9]+$/) next
		ni++; num[ni] = $1 + 0; open_[ni] = ($2 == "open"); reason[ni] = $3; labs[ni] = $4
		next
	}
	END {
		for (k in set) nset++
		if (!nset) { print "task-labels: no task issue in the plan, the completed log or the backlog" > "/dev/stderr"; exit 1 }
		label["done"] = "state: merged"; label["in review"] = "state: in review"
		label["in progress"] = "state: in progress"; label["stalled"] = "state: stalled"
		label["blocked"] = "state: blocked"; label["ready"] = "state: ready"
		for (i = 1; i <= ni; i++) {
			n = num[i]
			if (n in set) seen++
			stale = (labs[i] ~ /(^|,)state: /)
			if (!open_[i]) {
				if (!(n in set) && !stale) continue
				want = (reason[i] == "completed") ? "state: merged" : "state: rejected"
			} else if (n in st) {
				if (!(st[n] in label)) { printf "task-labels: #%d has a state the labels do not know: %s\n", n, st[n] > "/dev/stderr"; exit 1 }
				want = label[st[n]]
			} else if (n in pr) want = "state: in review"
			else want = "state: on hold"
			have = 0; out = ""
			k = split(labs[i], l, ",")
			for (j = 1; j <= k; j++) {
				if (l[j] !~ /^state: /) continue
				if (l[j] == want && !have) { have = 1; continue }
				out = out n " 0 remove #" n " " l[j] "\n"
			}
			if (want != "" && !have) out = out n " 1 add #" n " " want "\n"
			printf "%s", out
		}
		if (!seen) { print "task-labels: the issue read holds no task issue" > "/dev/stderr"; exit 1 }
	}' "$tmp/r.plan" "$tmp/r.completed" "$tmp/r.backlog" "$tmp/r.state" "$tmp/r.prs" "$tmp/r.issues" > "$tmp/r.changes" || exit 1
	sort -t' ' -k1,1n -k2,2n -s "$tmp/r.changes" | cut -d' ' -f3-
}

# check_workflow FILE — the rules of the header of task-state.yml.
# LIMIT: the header it reads, not the workflow: a second checkout of another
# `ref:` passes (only the first `ref:` of the default branch sets the flag); an
# event field in a `uses:`, `with:` or `env:` passes (only `run:` lines are read
# for `${{`); a bound written as YAML's explicit key (`? tags`) records no bound
# and passes. Each form needs an edit of the file on `main`, which review of
# `.github/**` covers; the check is a guard against forgetting.
check_workflow() {
	[ -f "$1" ] || die "missing input: the workflow $1"
	tr -d '\r' < "$1" | awk '
	BEGIN { squo = sprintf("%c", 39); dquo = sprintf("%c", 34); keyre = "^ *[" squo dquo "]?[a-z_-]+[" squo dquo "]?[ \t]*:" }
	function ind(l) { match(l, /^ */); return RLENGTH }
	function fail(m) { print "task-labels: the workflow FAIL " m > "/dev/stderr"; bad = 1 }
	function isin(v, w, n,   i) { for (i = 1; i <= n; i++) if (v == w[i]) return 1; return 0 }
	function boundtypes(ev, wants,   i, n, w, m, tv) {
		n = split(wants, w, " ")
		m = split(ord[ev SUBSEP "types"], tv, " ")
		for (i = 1; i <= m; i++) if (tv[i] != "" && !isin(tv[i], w, n)) fail(ev ": types: " tv[i] " is not one of " wants)
		for (i = 1; i <= n; i++) if (!((ev SUBSEP "types" SUBSEP w[i]) in key)) fail(ev ": types: " w[i] " is missing")
	}
	/^[ \t]*#/ || /^[ \t]*$/ { next }
	ind($0) == 0 { top = $0; sub(/:.*/, "", top); inrun = 0 }
	top == "on" && ind($0) == 2 { k = $0; sub(/^ */, "", k); sub(/:.*/, "", k); trig[k] = 1; cur = k; curk = "" }
	top == "on" && ind($0) == 4 && $0 ~ keyre {
		kn = $0; sub(/^ */, "", kn)
		if (substr(kn, 1, 1) == squo || substr(kn, 1, 1) == dquo) kn = substr(kn, 2)
		sub(/:.*/, "", kn); sub(/[ \t]*$/, "", kn)
		if (substr(kn, length(kn), 1) == squo || substr(kn, length(kn), 1) == dquo) kn = substr(kn, 1, length(kn) - 1)
		curk = kn
		if (cur == "push") bound[kn] = 1
		v = $0; sub(keyre "[ \t]*", "", v)
		if (v ~ /^\[/) { sub(/^\[/, "", v); sub(/\][ \t]*$/, "", v)
			m = split(v, a, /,[ \t]*/)
			for (i = 1; i <= m; i++) if (a[i] != "") { key[cur SUBSEP curk SUBSEP a[i]] = 1; ord[cur SUBSEP curk] = ord[cur SUBSEP curk] " " a[i] } }
	}
	top == "on" && ind($0) == 4 && $0 ~ /^ *- / { if ($0 ~ /cron:/) cron[cur] = 1 }
	top == "on" && ind($0) > 4 && $0 ~ /^ *- / && curk != "" {
		v = $0; sub(/^ *- */, "", v); sub(/[ \t]*$/, "", v)
		key[cur SUBSEP curk SUBSEP v] = 1; ord[cur SUBSEP curk] = ord[cur SUBSEP curk] " " v
	}
	top == "permissions" && ind($0) == 2 { k = $0; sub(/^ */, "", k); perm[k] = 1 }
	top == "jobs" && /^ +permissions:/ { fail("permissions: a job sets its own") }
	/uses: actions\/checkout@/ { checkout = 1 }
	/^ +ref: \$\{\{ github\.event\.repository\.default_branch \}\}[ \t]*$/ { if (checkout) ref = 1 }
	/^ +- name: Restore/ { restore = 1; inres = 1; next }
	/^ +- / { inres = 0 }
	inres && /docs\/tasks/ { resdocs = 1 }
	inrun && ind($0) <= runind { inrun = 0 }
	inrun && /\$\{\{/ { fail("run: an expression in a run: step") }
	/^ +(- )?run:/ {
		v = $0; sub(/^ +(- )?run:[ \t]*/, "", v)
		if (v ~ /\$\{\{/) fail("run: an expression in a run: step")
		if (v ~ /^[|>]/) { inrun = 1; runind = ind($0) }
	}
	END {
		n = split("push pull_request_target issues issue_comment schedule", want, " ")
		for (t in trig) { ok = 0; for (i = 1; i <= n; i++) if (t == want[i]) ok = 1; if (!ok) fail("triggers: " t " is not one of push, pull_request_target, issues, issue_comment, schedule") }
		for (i = 1; i <= n; i++) if (!(want[i] in trig)) fail("triggers: " want[i] " is missing")
		n = split("contents: read|issues: write|pull-requests: read", pw, "|")
		for (p in perm) { ok = 0; for (i = 1; i <= n; i++) if (p == pw[i]) ok = 1; if (!ok) fail("permissions: " p " is not one of contents: read, issues: write, pull-requests: read") }
		for (i = 1; i <= n; i++) if (!(pw[i] in perm)) fail("permissions: " pw[i] " is missing")
		for (b in bound) if (b != "branches") fail("push: " b ": a bound of push is branches only")
		m = split(ord["push" SUBSEP "branches"], pb, " ")
		for (i = 1; i <= m; i++) if (pb[i] != "" && pb[i] != "main") fail("push: branches: " pb[i] " is not main")
		if (!(("push" SUBSEP "branches" SUBSEP "main") in key)) fail("push: branches: main is missing")
		boundtypes("pull_request_target", "opened reopened closed edited")
		boundtypes("issues", "opened closed reopened")
		boundtypes("issue_comment", "created")
		if (("schedule" in trig) && !("schedule" in cron)) fail("schedule: no cron entry")
		if (!ref) fail("checkout: no ref: of the default branch")
		if (!restore || !resdocs) fail("restore: no step named Restore that names docs/tasks")
		if (bad) exit 1
		print "task-labels: the workflow OK"
	}'
}

# apply DRY — read everything, then (unless DRY) write.
apply() {
	mkdir "$tmp/in" || die "cannot make $tmp/in"
	say "running task-state.sh"
	: > "$tmp/in/prs"
	if [ -n "${TASK_STATE_CMD:-}" ]; then sh "$TASK_STATE_CMD" --prs "$tmp/in/prs" > "$tmp/in/state"
	else sh "$here/task-state.sh" --prs "$tmp/in/prs" > "$tmp/in/state"; fi || die "task-state.sh failed, so no label is written"
	say "reading the backlog, the completed log and the plan of main"
	for _f in backlog:docs/tasks/backlog.md completed:docs/tasks/completed.md plan:docs/plan/README.md; do
		gh api -H 'Accept: application/vnd.github.raw' "repos/{owner}/{repo}/contents/${_f#*:}?ref=main" > "$tmp/in/${_f%%:*}" ||
			die "gh: cannot read ${_f#*:} of main, so no label is written"
	done
	say "reading the issues"
	gh api 'repos/{owner}/{repo}/issues?state=all&per_page=100' --paginate \
		-q '.[] | select(.pull_request | not) | "\(.number)\t\(.state)\t\(.state_reason // "")\t\([.labels[].name] | join(","))"' \
		> "$tmp/in/issues" || die "gh: cannot read the issues, so no label is written"
	read_dir "$tmp/in" > "$tmp/changes" || exit 1
	if [ "$1" = dry ]; then cat "$tmp/changes"; return 0; fi
	say "making sure the eight labels exist, each in the color of its state"
	# LABEL=COLOR: the pairs the header names, in the order the calls are
	# recorded; a label holds no `=`, so the split is on it.
	for _lc in 'state: ready=1d76db' 'state: in progress=fbca04' 'state: stalled=d93f0b' 'state: in review=5319e7' 'state: blocked=b60205' 'state: merged=0e8a16' 'state: rejected=57606a' 'state: on hold=c5def5'; do
		gh label create "${_lc%%=*}" --force --color "${_lc##*=}" --description 'Set by .github/workflows/task-state.yml; never by hand' > /dev/null ||
			die "gh: cannot create the label ${_lc%%=*}"
	done
	say "writing $(grep -c . "$tmp/changes") changes"
	# One edit per issue: its removals, then its addition; a `-` field is
	# empty (an empty tab field would collapse in `read` and shift a removal
	# into the addition).
	awk '{ n = $2; sub(/^#/, "", n); act = $1; l = $0; sub(/^[a-z]+ #[0-9]+ /, "", l)
		if (n != cur) { if (cur != "") print cur "\t" (add == "" ? "-" : add) "\t" (rem == "" ? "-" : rem); cur = n; add = ""; rem = "" }
		if (act == "add") add = l; else rem = (rem == "" ? l : rem "," l) }
		END { if (cur != "") print cur "\t" (add == "" ? "-" : add) "\t" (rem == "" ? "-" : rem) }' "$tmp/changes" > "$tmp/edits"
	while IFS="$(printf '\t')" read -r _n _add _rem; do
		[ "$_add" = - ] && _add=
		[ "$_rem" = - ] && _rem=
		set -- issue edit "$_n"
		[ -n "$_add" ] && set -- "$@" --add-label "$_add"
		[ -n "$_rem" ] && set -- "$@" --remove-label "$_rem"
		gh "$@" > /dev/null || die "gh: cannot edit the labels of #$_n"
	done < "$tmp/edits"
}

case $# in
0) apply write ;;
1)
	[ "$1" = --dry-run ] || { printf 'task-labels: unknown argument: %s\n' "$1" >&2; exit 2; }
	apply dry ;;
2)
	case $1 in
	--read) [ -d "$2" ] || { printf 'task-labels: not a directory: %s\n' "$2" >&2; exit 2; }; read_dir "$2" ;;
	--check-workflow) check_workflow "$2" ;;
	*) printf 'task-labels: unknown argument: %s\n' "$1" >&2; exit 2 ;;
	esac ;;
*) printf 'task-labels: usage: task-labels.sh [--dry-run | --read DIR | --check-workflow FILE]\n' >&2; exit 2 ;;
esac
