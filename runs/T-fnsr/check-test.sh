#!/bin/sh
# check-test.sh — the offline test of check.sh (task T-fnsr, condition 2 of
# the plan review). It builds a target as Start writes it, as a local bare
# repository (WEB) and saved API responses (API, read by file://), checks that
# check.sh passes on it, then breaks one thing per case and checks that the
# FAIL line of that check comes. No network. Exit 1 when a case fails.
set -u
HERE=$(cd "$(dirname "$0")" && pwd)
ROOT=$(cd "$HERE/../.." && pwd)
W=$(mktemp -d) || exit 2
trap 'rm -rf "$W"' EXIT
T=$(printf '\t')
BOT='layup-agent[bot]' MAIL='335371832+layup-agent[bot]@users.noreply.github.com'
PIN=a95965534b14b0bf14ad74da0c9a45b5f4aedf88
SRC=$(sed -n 's/^source=//p' "$ROOT/docs/setup/armature.pin")
PSB=$(if command -v sha256sum >/dev/null; then sha256sum "$HERE/brief.md"; else shasum -a 256 "$HERE/brief.md"; fi | cut -d' ' -f1)
fails=0

# good writes the files of a correct records branch into $F, and the API.
good() {
	rm -rf "$W/f" "$W/api" "$W/web" && mkdir -p "$W/f/start" "$W/api/users" "$W/api/repos/pharzam/layup-uat/issues" "$W/web/pharzam"
	F=$W/f
	sed -n '/^```text start-readme$/,/^```$/p' "$ROOT/docs/spec/run.md" | sed '1d;$d' >"$F/README.md"
	cp "$HERE/brief.md" "$F/start/problem-statement.md"
	printf 'id\trole\tlogin\tsince\tsource\n1675602\toperator\tpharzam\t2026-10-08T14:00:00Z\tstart\n1675602\tidea-owner\tpharzam\t2026-10-08T14:00:00Z\tstart\n' >"$F/approvers.tsv"
	printf 'run\thost\tversion\tstarted\theartbeat\tstate\n0123456789abcdef\tantarctica.local\t0.1.0-dev\t2026-10-08T14:00:00Z\t0\treleased\n' >"$F/lease.tsv"
	printf '%s\n' "name|value|source" "layup.version|0.1.0-dev|run" "psb.sha256|$PSB|command" "vision.sha256|—|command" \
		"forge.plan|free|command" "forge.visibility|public|forge" "app.permissions|contents:write,issues:write,metadata:read|forge" \
		"operator.id|1675602|forge" "idea-owner.id|1675602|forge" "intake.cap|5,1|command" "lease.H|60|command" "watch.T|1|command" \
		"pin.source|$SRC|run" "pin.commit|$PIN|run" "pin.tree|TREE|run" "pin.time|2026-09-23T00:00:00Z|run" \
		"issue.intake|1|forge" "issue.control|2|forge" "watch|not-confirmed|run" | tr '|' '\t' >"$F/start/start.tsv"
	cp "$HERE/fixtures/users-bot.json" "$W/api/users/$BOT"
	for i in 1 2; do printf '{"number": %s, "user": {"login": "%s", "id": 335371832, "type": "Bot"}}\n' $i "$BOT" >"$W/api/repos/pharzam/layup-uat/issues/$i"; done
	ROOTWHO="$BOT|$MAIL" RECWHO="$BOT|$MAIL" ROOTSUBJ="chore: the unmodified baseline at $PIN" EXTRA= DROP= MAINTWO= NORECORDS=
}

# build makes the bare target from $F and the variables of good.
build() {
	R=$W/r && rm -rf "$R" && git init -q -b main "$R" && echo baseline >"$R/a"
	as() { who=$1; shift; GIT_AUTHOR_NAME=${who%|*} GIT_AUTHOR_EMAIL=${who#*|} GIT_COMMITTER_NAME=${who%|*} GIT_COMMITTER_EMAIL=${who#*|} git -C "$R" commit -q "$@"; }
	git -C "$R" add a && as "$ROOTWHO" -m "$ROOTSUBJ"
	[ -z "$MAINTWO" ] || as "$ROOTWHO" --allow-empty -m second
	tree=$(git -C "$R" rev-parse 'HEAD^{tree}')
	sed "s/	TREE	/	$tree	/" "$F/start/start.tsv" >"$F/start/s" && mv "$F/start/s" "$F/start/start.tsv"
	git -C "$R" checkout -q --orphan layup-records && git -C "$R" rm -rqf . && cp -R "$F/." "$R/"
	git -C "$R" add -A && as "$RECWHO" -m "records: the Start of pharzam/layup-uat"
	for s in "records: issue.intake opening" "records: issue.intake 1" "records: issue.control opening" \
		"records: issue.control 2" "$EXTRA" "records: watch not-confirmed" "lease: released"; do
		[ -z "$s" ] || [ "$s" = "$DROP" ] || as "$RECWHO" --allow-empty -m "$s"
	done
	git clone -q --bare "$R" "$W/web/pharzam/layup-uat.git"
	[ -z "$NORECORDS" ] || git -C "$W/web/pharzam/layup-uat.git" branch -qD layup-records
}

# case NAME WANT: build, run check.sh, and look for the line WANT.
case_() {
	build
	out=$(sh "$HERE/check.sh" "file://$W/web" "file://$W/api" pharzam/layup-uat 2>&1)
	code=$?
	if [ "$2" = pass ]; then
		[ "$code" -eq 0 ] && echo "ok: $1: exit 0" || { echo "FAIL: $1: exit $code"; echo "$out" | grep FAIL; fails=$((fails + 1)); }
	elif [ "$code" -eq 1 ] && echo "$out" | grep -qF "FAIL: $2"; then
		echo "ok: $1: FAIL: $2"
	else
		echo "FAIL: $1: no 'FAIL: $2' (exit $code)"; fails=$((fails + 1))
	fi
}
edit() { sed "$2" "$1" >"$W/e" && mv "$W/e" "$1"; }
tsv() { sed "s/^$1	[^	]*/$1	$2/" "$F/start/start.tsv" >"$W/s" && mv "$W/s" "$F/start/start.tsv"; }

good; case_ "a correct target" pass
good; EXTRA="lease: heartbeat 1"; case_ "a heartbeat commit" pass
good; rm -rf "$W/web" && mkdir -p "$W/web"; out=$(sh "$HERE/check.sh" "file://$W/web" "file://$W/api" pharzam/layup-uat)
echo "$out" | grep -qF "FAIL: the clone of" && echo "ok: no repository: FAIL: the clone of" || { echo "FAIL: no repository"; fails=$((fails + 1)); }
good; NORECORDS=1; case_ "no records branch" "no branch layup-records"
good; rm "$W/api/users/$BOT"; case_ "no bot response" "the ID of $BOT: no answer"
good; MAINTWO=1; case_ "two commits on main" "main has 2 commits"
good; ROOTWHO="pharzam|p@example.com"; case_ "a root commit by the Operator" "the root commit is by pharzam"
good; ROOTSUBJ="chore: other"; case_ "a root subject" "the root commit subject"
good; RECWHO="pharzam|p@example.com"; case_ "records by the Operator" "7 commits of layup-records are not by"
good; DROP="records: issue.control 2"; case_ "a missing commit" "the commits of Start"
good; echo x >"$F/start/vision.md"; case_ "a vision brief" "the files:"
good; echo x >>"$F/README.md"; case_ "another README" "README.md is not the text"
good; echo x >>"$F/start/problem-statement.md"; case_ "another brief" "start/problem-statement.md is not the brief"
good; edit "$F/start/start.tsv" '1s/source/src/'; case_ "the header of start.tsv" "the header of start.tsv"
for k in layup.version psb.sha256 vision.sha256 forge.plan forge.visibility operator.id idea-owner.id intake.cap lease.H watch.T pin.source pin.commit pin.tree watch; do
	good; tsv "$k" x; case_ "start.tsv $k" "start.tsv $k is 'x'"
done
for k in app.permissions pin.time; do good; tsv "$k" ""; case_ "start.tsv $k" "start.tsv $k is empty"; done
good; printf 'harness.claude.cap\t10.0\tregister\n' >>"$F/start/start.tsv"; case_ "a harness row" "start.tsv has a harness row"
good; edit "$F/approvers.tsv" 's/idea-owner/approver/'; case_ "approvers.tsv" "approvers.tsv:"
good; edit "$F/lease.tsv" '1s/host/where/'; case_ "the header of lease.tsv" "the header of lease.tsv"
good; edit "$F/lease.tsv" '2s/released/held/'; case_ "a held lease" "lease.tsv: 6|0.1.0-dev|held"
for n in 1 2; do
	good; edit "$W/api/repos/pharzam/layup-uat/issues/$n" 's/layup-agent\[bot\]/pharzam/'
	k=intake; [ $n -eq 1 ] || k=control; case_ "the $k issue" "the $k issue #$n is by 'pharzam'"
done

[ "$fails" -eq 0 ] && { echo "all cases pass"; exit 0; }
echo "$fails cases failed"; exit 1
