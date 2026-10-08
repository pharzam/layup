#!/bin/sh
# Check that each source path that survey.md cites exists in its repository at
# the commit that the survey read. A one-off evidence command of task T-w73g,
# not a check of the gate.
#
#   sh runs/T-w73g/cites.sh [SURVEY [CLONES]]
#
# SURVEY is runs/T-w73g/survey.md when it is not given. CLONES is a directory
# that holds one clone of each repository, named as in the table below; when it
# is not given, each repository is fetched at its commit into a temporary
# directory (tree objects only, no file contents). No code of a repository runs.
#
# A source path is a backticked `<repo>:<path>` in the survey, where <repo> is
# one of the five names below; a trailing / is allowed. A directory counts as
# present. Exit 0 when each source path is present; exit 1 when the survey does
# not exist, cites no source path, a repository cannot be read at its commit, or
# a path is missing.

dir=$(dirname "$0")
survey=${1:-$dir/survey.md}
clones=$2

# repo name, directory name, GitHub path, commit. A copy of the commit column
# of the survey's table, by hand: the script does not read that table, so an
# edit of one must be made in the other (round 1, note 8).
repos='agnostic-ai Agnostic-AI ucsandman/Agnostic-AI aab09ff0d9533967045e7f2b04e478ec4e05528e
ruflo ruflo ruvnet/ruflo 9fa701a466159242dfec06379ce1278c849a3e34
cc-multi-cli-plugin cc-multi-cli-plugin greenpolo/cc-multi-cli-plugin 3dcb057c10b16d3ed911e8ce509e91398712c906
ruvnet-brain ruvnet-brain stuinfla/ruvnet-brain b4590d469f52937550d8d89afa86a394a1d72875
self-improvement-loop self-improvement-loop kolezka/self-improvement-loop ba245adbb5be8075fa01044dd82aac636155870f'

if [ ! -f "$survey" ]; then
	echo "cites: no survey: $survey" >&2
	exit 1
fi

tmp=$(mktemp -d) || exit 1
trap 'rm -rf "$tmp"' EXIT

grep -oE '`(agnostic-ai|ruflo|cc-multi-cli-plugin|ruvnet-brain|self-improvement-loop):[^`]+`' "$survey" |
	tr -d '`' | sort -u > "$tmp/cites"
total=$(wc -l < "$tmp/cites" | tr -d ' ')
if [ "$total" -eq 0 ]; then
	echo "cites: the survey cites no source path" >&2
	exit 1
fi

echo "$repos" | while read -r name dname gh commit; do
	if [ -n "$clones" ]; then
		repo=$clones/$dname
	else
		repo=$tmp/$dname
		echo "cites: fetching $gh at $commit" >&2
		git init -q "$repo" &&
			git -C "$repo" fetch -q --depth 1 --filter=blob:none \
				"https://github.com/$gh.git" "$commit" 2>/dev/null
	fi
	if ! git -C "$repo" cat-file -e "$commit^{commit}" 2>/dev/null; then
		echo "cites: cannot read $name at $commit" >&2
		echo 1 >> "$tmp/fail"
		continue
	fi
	grep "^$name:" "$tmp/cites" | while IFS= read -r cite; do
		path=${cite#"$name":}
		path=${path%/}
		if [ -n "$(git -C "$repo" ls-tree "$commit" -- "$path")" ]; then
			echo "ok      $cite"
		else
			echo "MISSING $cite"
			echo 1 >> "$tmp/fail"
		fi
	done
done

if [ -s "$tmp/fail" ]; then
	echo "cites: $(wc -l < "$tmp/fail" | tr -d ' ') problem(s) in $total source path(s)" >&2
	exit 1
fi
echo "cites: $total source path(s), each present at its commit"
