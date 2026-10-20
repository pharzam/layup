#!/bin/sh
#
# run.sh CASE_DIR — run docs/tests/freeze-check.sh on one fixture case.
#
# Builds a disposable git repository in mktemp -d from base/ (commit 1, branch
# `base`) and the case's head/ overlay (commit 2, `HEAD`), stripping the `.in`
# suffix from each stored input file. Runs freeze-check.sh with the case's ARGS
# and compares exit status + stdout with EXPECT (first line the expected exit
# status). Exits 0 on match, 2 on mismatch or usage error.

set -u

[ $# -eq 1 ] && [ -d "$1" ] || { printf 'run: usage: run.sh CASE_DIR\n' >&2; exit 2; }
case_dir=${1%/}
name=$(basename "$case_dir")
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
check=$here/../freeze-check.sh
[ -f "$case_dir/EXPECT" ] && [ -f "$case_dir/ARGS" ] || { printf 'run: %s lacks ARGS or EXPECT\n' "$name" >&2; exit 2; }

want_st=$(tr -d '\r' < "$case_dir/EXPECT" | sed -n '1p')
case $name in
good-*)
	[ "$want_st" = 0 ] || { printf 'run: %s: good-* EXPECT must start with 0\n' "$name" >&2; exit 2; } ;;
bad-usage-*|bad-not-a-worktree*|bad-base-not-a-commit*)
	[ "$want_st" = 2 ] || { printf 'run: %s: frame bad-* EXPECT must start with 2\n' "$name" >&2; exit 2; } ;;
bad-*)
	[ "$want_st" = 1 ] || { printf 'run: %s: bad-* EXPECT must start with 1\n' "$name" >&2; exit 2; } ;;
*)
	printf 'run: %s is neither good-* nor bad-*\n' "$name" >&2; exit 2 ;;
esac

tmp=$(mktemp -d) || exit 2
trap 'cd / || :; rm -rf "${tmp:?}"' EXIT
export GIT_CEILING_DIRECTORIES="$tmp"
repo=$tmp/repo
mkdir -p "$repo" || exit 2

overlay() {
	_src=$1
	[ -d "$_src" ] || return 0
	find "$_src" -type f -name '*.in' | sort | while IFS= read -r _f; do
		_rel=${_f#"$_src"/}
		if [ "$_rel" = ABSENT.in ]; then
			tr -d '\r' < "$_f" | while IFS= read -r _rm; do
				[ -n "$_rm" ] && rm -rf "${repo:?}/$_rm"
			done
		else
			_dst=$repo/${_rel%.in}
			mkdir -p "$(dirname -- "$_dst")" && cp "$_f" "$_dst" || exit 2
		fi
	done
}

git -C "$repo" init -q && git -C "$repo" symbolic-ref HEAD refs/heads/main || exit 2
git -C "$repo" config maintenance.auto false || exit 2
git -C "$repo" config core.autocrlf false || exit 2
case $name in
*crlf*) git -C "$repo" config core.whitespace cr-at-eol || exit 2 ;;
esac

overlay "$here/base" || exit 2
git -C "$repo" add -A \
	&& GIT_AUTHOR_DATE=2026-01-01T00:00:00Z GIT_COMMITTER_DATE=2026-01-01T00:00:00Z \
		git -C "$repo" -c maintenance.auto=false -c user.name=test -c user.email=test@test.invalid commit -qm base \
	&& git -C "$repo" branch base || exit 2

overlay "$case_dir/head" || exit 2
git -C "$repo" add -A \
	&& GIT_AUTHOR_DATE=2026-01-01T00:01:00Z GIT_COMMITTER_DATE=2026-01-01T00:01:00Z \
		git -C "$repo" -c maintenance.auto=false -c user.name=test -c user.email=test@test.invalid commit --allow-empty -qm head || exit 2

case $name in
bad-not-a-worktree*) rm -rf "${repo:?}/.git" ;;
esac

args=$(tr -d '\r' < "$case_dir/ARGS")
# shellcheck disable=SC2086
(cd "$repo" && sh "$check" $args) > "$tmp/out" 2>&1
got=$?
{ printf '%s\n' "$got"; cat "$tmp/out"; } > "$tmp/got"
tr -d '\r' < "$case_dir/EXPECT" > "$tmp/expect"

if cmp -s "$tmp/got" "$tmp/expect"; then
	exit 0
fi
printf 'run: %s: exit %s (wanted %s); diff (-EXPECT +got):\n' "$name" "$got" "$want_st" >&2
diff -u "$tmp/expect" "$tmp/got" >&2
exit 2
