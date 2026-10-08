#!/bin/sh
# check.sh WEB API OWNER/NAME — the check of the demo of M2a (task T-fnsr).
# It reads a target that `layup run --new` started, with a plain `git clone`
# of WEB/OWNER/NAME.git and unauthenticated reads of API (no token), and
# checks each write of LAYUP against docs/spec/run.md and records.md: the
# author and committer of each commit, the files and values of the records
# branch, the pin, and the author of the two issues. It prints `ok` or
# `FAIL: <reason>` per check, and exits 1 when one fails.
# Example: sh runs/T-fnsr/check.sh https://github.com https://api.github.com pharzam/layup-uat
set -u
[ $# -eq 3 ] || { echo "usage: check.sh WEB API OWNER/NAME" >&2; exit 2; }
WEB=$1 API=$2 TARGET=$3
HERE=$(cd "$(dirname "$0")" && pwd)
ROOT=$(cd "$HERE/../.." && pwd)
BRIEF=$HERE/brief.md
SLUG=layup-agent BOT='layup-agent[bot]' OPERATOR=pharzam OPERATOR_ID=1675602 VERSION=0.1.0-dev
T=$(printf '\t')
fails=0
ok() { echo "ok: $1"; }
bad() { echo "FAIL: $1"; fails=$((fails + 1)); }
get() { curl -sfg "$API/$1" | tr -d '\n'; }
field() { sed -n "s/.*\"$1\": *\"\{0,1\}\([^\",}]*\).*/\1/p"; }
sha() { if command -v sha256sum >/dev/null; then sha256sum "$1"; else shasum -a 256 "$1"; fi | cut -d' ' -f1; }

W=$(mktemp -d) || exit 2
trap 'rm -rf "$W"' EXIT
if ! git clone -q --no-tags "$WEB/$TARGET.git" "$W/c" 2>"$W/err"; then
	bad "the clone of $WEB/$TARGET.git: $(head -1 "$W/err")"
	echo "$fails failed"; exit 1
fi
G() { git -C "$W/c" "$@"; }
for b in main layup-records; do
	G rev-parse -q --verify "origin/$b" >/dev/null || bad "no branch $b"
done
[ "$fails" -eq 0 ] || { echo "$fails failed"; exit 1; }
show() { G show "origin/layup-records:$1"; }
row() { show "$1" | awk -F"$T" -v k="$2" '$1 == k { print $2 }'; }

# The bot's ID, by the call that Start uses (GET /users/<slug>[bot]).
id=$(get "users/$SLUG%5Bbot%5D" | field id | head -1)
case $id in '' | *[!0-9]*) bad "the ID of $BOT: no answer from $API"; id=0 ;; *) ok "the ID of $BOT is $id" ;; esac
WHO="$BOT|$id+$BOT@users.noreply.github.com"

# The root commit of main: one commit with no parent, by the bot.
n=$(G rev-list --count origin/main)
[ "$n" -eq 1 ] && ok "main has one commit" || bad "main has $n commits"
c=$(G log -1 --format='%an|%ae|%cn|%ce|%s' origin/main)
[ "${c%|*}" = "$WHO|$WHO" ] && ok "the root commit is by $BOT" || bad "the root commit is by ${c%|*}"
subject=${c##*|}
root_pin=${subject#chore: the unmodified baseline at }
[ "$root_pin" != "$subject" ] && ok "the root commit names the commit $root_pin" || bad "the root commit subject: $subject"

# The commits of layup-records: each by the bot, the seven of Start in order.
G log --reverse --format='%an|%ae|%cn|%ce|%s' origin/layup-records >"$W/log"
others=$(cut -d'|' -f1-4 "$W/log" | grep -cvxF "$WHO|$WHO")
[ "$others" -eq 0 ] && ok "each commit of layup-records is by $BOT" || bad "$others commits of layup-records are not by $BOT"
sed 's/^[^|]*|[^|]*|[^|]*|[^|]*|//' "$W/log" >"$W/subjects"
beats=$(grep -c '^lease: heartbeat [0-9][0-9]*$' "$W/subjects")
[ "$beats" -eq 0 ] || ok "$beats heartbeat commits, by $BOT, apart from the seven"
intake=$(row start/start.tsv issue.intake) control=$(row start/start.tsv issue.control)
printf '%s\n' "records: the Start of $TARGET" "records: issue.intake opening" "records: issue.intake $intake" \
	"records: issue.control opening" "records: issue.control $control" "records: watch not-confirmed" "lease: released" >"$W/want"
grep -v '^lease: heartbeat [0-9][0-9]*$' "$W/subjects" | diff "$W/want" - >"$W/d" &&
	ok "the seven commits of Start, in order" || bad "the commits of Start: $(grep '^[<>]' "$W/d" | tr '\n' ' ')"

# The files of the records branch, and the README of run.md.
G ls-tree -r --name-only origin/layup-records >"$W/files"
printf '%s\n' README.md approvers.tsv lease.tsv start/problem-statement.md start/start.tsv | diff - "$W/files" >/dev/null &&
	ok "the five files, no start/vision.md" || bad "the files: $(tr '\n' ' ' <"$W/files")"
sed -n '/^```text start-readme$/,/^```$/p' "$ROOT/docs/spec/run.md" | sed '1d;$d' >"$W/readme"
show README.md | cmp -s - "$W/readme" && ok "README.md is the text of run.md" || bad "README.md is not the text of run.md"
show start/problem-statement.md | cmp -s - "$BRIEF" && ok "the brief, byte for byte" || bad "start/problem-statement.md is not the brief"

# start.tsv: each value known before the run.
[ "$(show start/start.tsv | head -1)" = "name${T}value${T}source" ] && ok "the header of start.tsv" || bad "the header of start.tsv"
src=$(sed -n 's/^source=//p' "$ROOT/docs/setup/armature.pin")
tree=$(G rev-parse 'origin/main^{tree}')
while IFS='|' read -r k v; do
	got=$(row start/start.tsv "$k")
	[ "$got" = "$v" ] && ok "start.tsv $k = $v" || bad "start.tsv $k is '$got', not '$v'"
done <<EOF
layup.version|$VERSION
psb.sha256|$(sha "$BRIEF")
vision.sha256|—
forge.plan|free
forge.visibility|public
operator.id|$OPERATOR_ID
idea-owner.id|$OPERATOR_ID
intake.cap|5,1
lease.H|60
watch.T|1
pin.source|$src
pin.commit|$root_pin
pin.tree|$tree
watch|not-confirmed
EOF
show start/start.tsv | grep -q '^harness\.' && bad "start.tsv has a harness row" || ok "start.tsv has no harness row"
for k in app.permissions pin.time; do [ -n "$(row start/start.tsv $k)" ] && ok "start.tsv $k is set" || bad "start.tsv $k is empty"; done

# approvers.tsv and lease.tsv.
printf 'id\trole\tlogin\tsince\tsource\n%s\toperator\t%s\n%s\tidea-owner\t%s\n' "$OPERATOR_ID" "$OPERATOR" "$OPERATOR_ID" "$OPERATOR" >"$W/want"
show approvers.tsv | awk -F"$T" -v OFS="$T" 'NR == 1 { print; next } { print $1, $2, $3 ($5 == "start" ? "" : " source " $5) }' | diff "$W/want" - >/dev/null &&
	ok "approvers.tsv: the two rows of Start" || bad "approvers.tsv: $(show approvers.tsv | tr '\t\n' ' |')"
show lease.tsv >"$W/lease"
[ "$(head -1 "$W/lease")" = "run${T}host${T}version${T}started${T}heartbeat${T}state" ] && ok "the header of lease.tsv" || bad "the header of lease.tsv"
lease=$(sed 1d "$W/lease" | awk -F"$T" '{ print NF "|" $3 "|" $6 }')
[ "$lease" = "6|$VERSION|released" ] && ok "lease.tsv: one row, $VERSION, released" || bad "lease.tsv: $lease"

# The two issues, opened by the bot.
for k in intake control; do
	eval "n=\$$k"
	login=$(get "repos/$TARGET/issues/$n" | sed -n 's/.*"user": *{ *"login": *"\([^"]*\)".*/\1/p')
	[ "$login" = "$BOT" ] && ok "the $k issue #$n is by $BOT" || bad "the $k issue #$n is by '$login'"
done

[ "$fails" -eq 0 ] && { echo "all checks pass"; exit 0; }
echo "$fails failed"; exit 1
