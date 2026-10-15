#!/bin/sh
# git-tee-test.sh: the test of git-tee.sh, run before Start (task T-sdtb, #284). A bare repository whose
# pre-receive hook refuses with the text of GH013 and one that accepts: the wrapped push gives the same exit
# status, standard output and standard error as the real git, the log holds the entry, and a command that is not
# a push is passed on. It prints `ok` or `FAIL: <reason>` for each check, and exits 1 when one fails.
set -u
here=$(cd "$(dirname "$0")" && pwd)
real=$(command -v git)
fails=0
ok() { echo "ok: $1"; }
bad() { echo "FAIL: $1"; fails=$((fails + 1)); }
W=$(mktemp -d) || exit 2
trap 'rm -rf "$W"' EXIT
git init -q --bare "$W/refuse.git" && git init -q --bare "$W/accept.git" || exit 2
printf '%s\n' '#!/bin/sh' 'echo "error: GH013: Repository rule violations found for refs/heads/layup-probe." >&2' \
	'echo "- Cannot create ref due to creations being restricted." >&2' 'exit 1' >"$W/refuse.git/hooks/pre-receive"
chmod +x "$W/refuse.git/hooks/pre-receive"
git init -q "$W/work" && git -C "$W/work" -c user.name=t -c user.email=t@example.invalid commit -q --allow-empty -m probe || exit 2
mkdir "$W/wrap"
sed "s|^real=.*|real=$real|; s|^log=.*|log=$W/push.log|" "$here/git-tee.sh" >"$W/wrap/git" && chmod +x "$W/wrap/git"
for t in refuse accept; do
	( cd "$W/work" && "$real" push "$W/$t.git" HEAD:refs/heads/layup-probe ) >"$W/$t.out1" 2>"$W/$t.err1"; rc1=$?
	git -C "$W/$t.git" update-ref -d refs/heads/layup-probe 2>/dev/null
	( cd "$W/work" && PATH="$W/wrap:$PATH" git push "$W/$t.git" HEAD:refs/heads/layup-probe ) >"$W/$t.out2" 2>"$W/$t.err2"; rc2=$?
	[ "$rc1" = "$rc2" ] && ok "$t: the same exit status ($rc1)" || bad "$t: exit $rc1 and $rc2"
	cmp -s "$W/$t.out1" "$W/$t.out2" && ok "$t: the same standard output" || bad "$t: the standard output differs"
	cmp -s "$W/$t.err1" "$W/$t.err2" && ok "$t: the same standard error" || bad "$t: the standard error differs"
	git -C "$W/$t.git" update-ref -d refs/heads/layup-probe 2>/dev/null
done
[ "$(grep -c '^== ' "$W/push.log")" = 2 ] && ok "the log has two entries" || bad "the log has $(grep -c '^== ' "$W/push.log") entries"
grep -q '^== [0-9T:Z-]* exit 1$' "$W/push.log" && grep -q '^== [0-9T:Z-]* exit 0$' "$W/push.log" && ok "the log has exit 1 and exit 0" || bad "the exit lines of the log"
grep -qxF "args: $W/refuse.git HEAD:refs/heads/layup-probe" "$W/push.log" && ok "the log has the URL and the refspec" || bad "the args line of the log"
grep -q 'GH013: Repository rule violations found for refs/heads/layup-probe' "$W/push.log" && ok "the log holds the text of GH013" || bad "the log lacks the text"
[ "$(cd "$W/work" && PATH="$W/wrap:$PATH" git rev-parse --git-dir)" = "$(cd "$W/work" && "$real" rev-parse --git-dir)" ] &&
	ok "a command that is not a push is passed on" || bad "rev-parse differs"
[ "$(grep -c '^== ' "$W/push.log")" = 2 ] && ok "a command that is not a push is not logged" || bad "rev-parse was logged"
[ "$fails" -eq 0 ] && { echo "all checks pass"; exit 0; }
echo "$fails failed"; exit 1
