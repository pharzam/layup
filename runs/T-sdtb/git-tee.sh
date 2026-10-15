#!/bin/sh
# git-tee.sh: a transparent wrapper of git for the demo of M2d (task T-sdtb, #284). It is installed as an
# executable file named `git` in a directory that comes first in the PATH of `layup run` (~/layup-host/wrap/).
# `internal/git` gives git the environment of PATH and TMPDIR only, so the two paths below are in the file.
# It runs the real git with the same arguments, standard input, standard output and environment. For a command
# that has the word `push` among its arguments it also keeps in the log the UTC time, the exit status, the last
# two arguments (the URL and the refspec: LAYUP gives the token to git in the environment, not in an argument)
# and the whole standard error, and it gives the caller the same standard error and the same exit status.
# Any other command is passed on with exec. It changes no byte that the caller reads.
real=/home/layup/.local/bin/git
log=/home/layup/layup-host/wrap/push.log
push=0
for a in "$@"; do [ "$a" = push ] && push=1; done
[ "$push" = 1 ] || exec "$real" "$@"
err=$(mktemp "${TMPDIR:-/tmp}/git-tee.XXXXXX") || exec "$real" "$@"
"$real" "$@" 2>"$err"
rc=$?
cat "$err" >&2
prev= last=
for a in "$@"; do prev=$last last=$a; done
{
	printf '== %s exit %s\nargs: %s %s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)" "$rc" "$prev" "$last"
	cat "$err"
} >>"$log"
rm -f "$err"
exit "$rc"
