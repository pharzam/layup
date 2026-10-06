#!/bin/sh
# T3 of the first pilot (task T-evad, pharzam/layup#97): the fail-fast run of the commands that step S13 of
# `layup setup` handed to the Operator in out/commands.sh.
#
# It runs the five commands of out/commands.sh one at a time, from the work area, and stops at the first
# command that fails. Nothing after a failed command runs. Before the last command (the ruleset), it reads
# back the pushed branches on GitHub; when they are not the approved commits, it stops, so the ruleset is
# never applied after a failed or wrong push. Before the first command, it checks (and changes nothing):
#   - out/commands.sh is the reviewed file (its SHA-256), with exactly five commands in the expected order;
#   - out/ruleset-default.json is the file that the record of S13 names (row ruleset.sha256);
#   - the target is at the approved commits: layup-setup, layup-records, and main = the root commit;
#   - layup-setup is a fast-forward of main, and the target has no remote yet;
#   - the repository on GitHub is empty;
#   - gh is logged in, and a dry run of the pushes works (only in the real run; it changes nothing).
# Each line of the run is also written to the log file (default: evidence/050-t3.txt of the pilot area).
#
# Usage: sh t3.sh [--check] [WORK]
#   --check   run only the checks before the first command, and stop: no command runs.
# The expected values have defaults for the real pilot; a test sets them in the environment:
#   T3_COMMANDS_SHA256, T3_HEAD, T3_RECORDS, T3_ROOT, T3_LOG.
# Exit 0 when the five commands ran (or, with --check, when every check passed); 1 when it stopped.
set -u
check_only=0
if [ "${1:-}" = --check ]; then check_only=1; shift; fi
W=${1:-/Users/farzam/layup-pilot/chat-orchestrator}
EXP_SHA=${T3_COMMANDS_SHA256:-ecbf3b8ae6edc2f50f968f86afe03414f06df9422a316decbbc54041ac683b87}
EXP_HEAD=${T3_HEAD:-cec749a92cc31a07d67cbf7bc74b07dcd85e84ee}
EXP_RECORDS=${T3_RECORDS:-2a339bb375a2e88d9d6c95cf6ca4734cc19920b6}
EXP_ROOT=${T3_ROOT:-242a2059929420875472478d1003b82e7c6130ad}
LOG=${T3_LOG:-/Users/farzam/layup-pilot/evidence/050-t3.txt}

say() { printf '%s\n' "$*"; printf '%s\n' "$*" >> "$LOG"; }
stop() {
	say "T3 STOPPED: $*"
	say "No command after this point ran. Do not run a command again by hand; send this output to the author."
	exit 1
}
sha256() { if command -v sha256sum >/dev/null 2>&1; then sha256sum "$1" | cut -d' ' -f1; else shasum -a 256 "$1" | cut -d' ' -f1; fi; }

: >> "$LOG" || { echo "T3 STOPPED: cannot write the log $LOG" >&2; exit 1; }
say "T3 $(date -u +%Y-%m-%dT%H:%M:%SZ)  work area: $W  mode: $([ $check_only = 1 ] && echo check || echo run)"
cd "$W" 2>/dev/null || stop "no work area $W"

# --- 1. The checks before the first command. They change nothing.
[ -f out/commands.sh ] || stop "no out/commands.sh"
got=$(sha256 out/commands.sh)
[ "$got" = "$EXP_SHA" ] || stop "out/commands.sh is not the reviewed file (SHA-256 $got, want $EXP_SHA)"
cmds=$(grep -v '^[[:space:]]*#' out/commands.sh | grep -v '^[[:space:]]*$')
[ "$(printf '%s\n' "$cmds" | wc -l | tr -d ' ')" = 5 ] || stop "out/commands.sh does not hold exactly five commands"
c1=$(printf '%s\n' "$cmds" | sed -n 1p); c2=$(printf '%s\n' "$cmds" | sed -n 2p); c3=$(printf '%s\n' "$cmds" | sed -n 3p)
c4=$(printf '%s\n' "$cmds" | sed -n 4p); c5=$(printf '%s\n' "$cmds" | sed -n 5p)
url=$(printf '%s\n' "$c1" | sed -n "s/^git -C target remote add origin '\\([^']*\\)'\$/\\1/p")
[ -n "$url" ] || stop "command 1 is not: git -C target remote add origin '<url>'"
[ "$c2" = "git -C target push origin main" ] || stop "command 2 is not: git -C target push origin main"
[ "$c3" = "git -C target push origin layup-setup:main" ] || stop "command 3 is not: git -C target push origin layup-setup:main"
[ "$c4" = "git -C target push origin layup-records" ] || stop "command 4 is not: git -C target push origin layup-records"
case $c5 in
"gh api --method POST 'repos/"*"/rulesets' --input out/ruleset-default.json") : ;;
*) stop "command 5 is not: gh api --method POST 'repos/<owner>/<name>/rulesets' --input out/ruleset-default.json" ;;
esac
want=$(awk -F'\t' '$1 == "S13" && $2 == "ruleset.sha256" { print $3 }' out/record.tsv)
got=$(sha256 out/ruleset-default.json)
[ -n "$want" ] && [ "$got" = "$want" ] || stop "out/ruleset-default.json is not the file that the record of S13 names ($got, want $want)"
h=$(git -C target rev-parse --verify -q layup-setup); r=$(git -C target rev-parse --verify -q layup-records); m=$(git -C target rev-parse --verify -q main)
[ "$h" = "$EXP_HEAD" ] || stop "layup-setup is $h, not the approved commit $EXP_HEAD"
[ "$r" = "$EXP_RECORDS" ] || stop "layup-records is $r, not $EXP_RECORDS"
[ "$m" = "$EXP_ROOT" ] || stop "main is $m, not the root commit $EXP_ROOT"
git -C target merge-base --is-ancestor main layup-setup || stop "layup-setup is not a fast-forward of main"
[ -z "$(git -C target remote)" ] || stop "the target has a remote already: $(git -C target remote | tr '\n' ' ')"
refs=$(git ls-remote "$url" 2>&1) || stop "cannot read $url: $refs"
[ -z "$refs" ] || stop "the repository $url is not empty"
say "checks: all passed (commands.sh $EXP_SHA; layup-setup $EXP_HEAD; layup-records $EXP_RECORDS; main $EXP_ROOT; $url is empty)"
if [ "$check_only" = 1 ]; then say "T3 check only: no command ran"; exit 0; fi
gh auth status >/dev/null 2>&1 || stop "gh is not logged in"
# a dry run of the two pushes: it needs the login of the push and changes nothing on GitHub
dry=$(git -C target push --dry-run "$url" main layup-records 2>&1) || stop "a dry run of the push failed (no change was made): $dry"
say "checks: gh is logged in; a dry run of the push works"

# --- 2. The five commands, one at a time. A command that fails stops the run.
out=$(mktemp) || stop "no temporary file"
trap 'rm -f "$out"' EXIT
n=0
for cmd in "$c1" "$c2" "$c3" "$c4" "$c5"; do
	n=$((n + 1))
	if [ "$n" = 5 ]; then
		# Before the ruleset: read back the pushed branches. A failed or wrong push stops the run here.
		rm_main=$(git ls-remote "$url" refs/heads/main | cut -f1)
		rm_rec=$(git ls-remote "$url" refs/heads/layup-records | cut -f1)
		[ "$rm_main" = "$EXP_HEAD" ] || stop "read-back: main on GitHub is '${rm_main:-none}', not $EXP_HEAD; the ruleset is not applied"
		[ "$rm_rec" = "$EXP_RECORDS" ] || stop "read-back: layup-records on GitHub is '${rm_rec:-none}', not $EXP_RECORDS; the ruleset is not applied"
		say "read-back: main = $rm_main, layup-records = $rm_rec on GitHub"
	fi
	say "[$n/5] $cmd"
	sh -c "$cmd" > "$out" 2>&1 < /dev/null
	rc=$?
	cat "$out"; cat "$out" >> "$LOG"
	[ "$rc" = 0 ] || stop "command $n failed with exit $rc: $cmd"
	say "[$n/5] ok"
done
say "T3 done: the five commands ran. The author reads back main, layup-records and the ruleset next."
exit 0
