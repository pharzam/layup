#!/bin/sh
# The deterministic part of the release review of phase 1 (task T-efmy, #96:
# D3 of its plan, with conditions 1 and notes 2 and 4 of its plan review),
# adapted for the release of M2b (task T-4tjy, #169: D3 of its plan, with
# conditions 3 and 5 and notes 3, 5 and 6 of its plan review;
# docs/spec/session.md, REQ-015 and REQ-017 — The review of the release of M2b).
# Its output at 621af09 is runs/T-efmy/release-check.txt, from the script
# before T-4tjy; its output at the release of M2b is
# runs/T-4tjy/release-check.txt.
#
# Run it from the root of a clone of LAYUP with the commit of the release as
# its argument, from a path outside that clone:
#
#   sh /path/to/release-check.sh COMMIT
#
# COMMIT is the full ID of the commit of the release; a name such as HEAD
# passes check (1) on any commit.
#
# It prints the facts of the code that REQ-015 and REQ-017 read, and gives
# exit 1 when one of its checks fails:
#   (1) HEAD is not COMMIT, or the tree has a change;
#   (2) TestPackageRules fails: rules 2 to 5 of NFR-007 in
#       docs/spec/packages.md (the standard library and this module only; no
#       package depends on net, net/http or crypto/tls unless its row says so:
#       internal/forge/github, cmd/layup and internal/cli may) and the program
#       of each exec.Command of each package's row;
#   (3) a call that starts a program, in a non-test Go file of cmd/ or
#       internal/ out of testdata/, is not one of this list of a file and its
#       program: git in internal/git/git.go; sh in internal/gate/scratch.go and
#       internal/verify/scripts.go; the first word of a harness register row's
#       version command (words[0]) or command (s.Words[0]) in
#       internal/session/process.go. The specification's "a new list names
#       each program that the code starts" is read as a check that fails, the
#       chosen reading (note 3 of the plan review of #169);
#   (4) a non-test Go file of internal/git holds a string literal, interpreted
#       ("push") or raw (`push`), that equals a git verb that reaches a remote:
#       push, send-pack, fetch, pull or remote, other than these, each in its
#       own function: fetch in Fetch and FetchSession, push in Push, and
#       remote in CloneLocal (remote remove origin, which cuts a session
#       clone's link and reaches no remote). A call of git names its verb so,
#       as in do(dir, "clone", ...). The function of a hit is the last line
#       that starts with "func" above it; a hit after a line "}" at column one
#       and before the next such line is in no function, and fails.
# The lists of (3), (5) and (6) are for the reviewer, who judges them.

commit=${1:?usage: sh release-check.sh COMMIT}
fail=0
bad() { echo "FAIL: $*"; fail=1; }
gocmd() { GOFLAGS= GOENV=off GOWORK=off GOTOOLCHAIN=local GOPROXY=off go "$@"; }

echo "== (1) The commit and the tree"
head=$(git rev-parse HEAD) || exit 1
want=$(git rev-parse --verify --quiet "$commit^{commit}")
echo "HEAD $head"
[ "$head" = "$want" ] || bad "HEAD is not the commit $commit"
changes=$(git status --porcelain --untracked-files=normal)
[ -z "$changes" ] || bad "the tree has a change: $changes"

echo "== (2) TestPackageRules (cmd/layup, integration)"
if out=$(gocmd test -count=1 -tags=integration -run '^TestPackageRules$' ./cmd/layup/ 2>&1); then
	echo "$out"
else
	echo "$out"
	bad "TestPackageRules"
fi
echo "The packages of go list -deps ./... that are not in the standard library:"
gocmd list -deps -f '{{if not .Standard}}{{.ImportPath}}{{end}}' ./... | sed 's/^/  /'
echo "net, net/http and crypto/tls in go list -deps ./...: [$(gocmd list -deps ./... | grep -x -E 'net|net/http|crypto/tls' | tr '\n' ' ')]"
echo "The packages that import net, net/http or crypto/tls themselves:"
gocmd list -f '{{.ImportPath}}:{{range .Imports}} {{.}}{{end}}' ./... | awk '{ for (i = 2; i <= NF; i++) if ($i == "net" || $i == "net/http" || $i == "crypto/tls") { print "  " $1 " " $i } }'

echo "== (3) The calls that start a program, in the non-test Go files"
sites=$(git grep -n -E 'exec\.Command|os\.StartProcess|syscall\.(Exec|ForkExec|StartProcess)' -- 'cmd/*.go' 'internal/*.go' ':!*_test.go' ':!*/testdata/*')
[ -n "$sites" ] || bad "no call that starts a program found"
echo "$sites" | sed 's/^/  /'
bad3=$(echo "$sites" | awk '
	/^internal\/git\/git\.go:[0-9]+:.*exec\.Command\("git", / { next }
	/^internal\/gate\/scratch\.go:[0-9]+:.*exec\.Command\("sh", / { next }
	/^internal\/verify\/scripts\.go:[0-9]+:.*exec\.Command\("sh", / { next }
	/^internal\/session\/process\.go:[0-9]+:.*exec\.CommandContext\(ctx, words\[0\], / { next }
	/^internal\/session\/process\.go:[0-9]+:.*exec\.Command\(s\.Words\[0\], / { next }
	NF { print }')
if [ -n "$bad3" ]; then
	# The loop runs in a subshell of the pipe, so fail is set here.
	echo "$bad3" | while IFS= read -r l; do echo "FAIL: a call that starts a program, out of the list: $l"; done
	fail=1
fi

echo "== (4) The git verbs that reach a remote, in the non-test Go files of internal/git"
# git ls-files drops every file with the exclusion ':!*_test.go', so the test
# files are filtered out by name; a list with no file fails, as no file read
# would pass.
gitfiles=$(git ls-files 'internal/git/*.go' | grep -v '_test\.go$')
[ -n "$gitfiles" ] || bad "no non-test Go file in internal/git"
hits=$(for f in $gitfiles; do
	awk -v f="$f" '
		/^func / { fn = $0; sub(/^func (\([^)]*\) )?/, "", fn); sub(/[^A-Za-z0-9_].*/, "", fn) }
		/^}/ { fn = "" }
		{
			line = $0
			while (match(line, /["`](push|send-pack|fetch|pull|remote)["`]/)) {
				verb = substr(line, RSTART + 1, RLENGTH - 2)
				print f ":" NR ": " verb " in " (fn == "" ? "no function" : fn)
				line = substr(line, RSTART + RLENGTH)
			}
		}' "$f"
done)
if [ -n "$hits" ]; then
	echo "$hits" | sed 's/^/  /'
	bad4=$(echo "$hits" | grep -v -E ': (fetch in (Fetch|FetchSession)|push in Push|remote in CloneLocal)$')
	if [ -n "$bad4" ]; then
		# The loop runs in a subshell of the pipe, so fail is set here.
		echo "$bad4" | while IFS= read -r l; do echo "FAIL: a git verb that reaches a remote, out of its function: $l"; done
		fail=1
	fi
else
	echo "  none"
fi
echo "The exported calls of internal/git:"
git grep -n -E '^func [A-Z]' -- 'internal/git/*.go' ':!*_test.go' | sed 's/^/  /'

echo "== (5) The commands that sh runs"
echo "The entries that the binary embeds:"
git grep -n 'go:embed' -- 'internal/catalog/*.go' ':!*_test.go' | sed 's/^/  /'
echo "The command of each kind of each entry, out of testdata (written into a target at S12):"
for f in $(git ls-files 'internal/catalog/*/kinds.tsv' ':!internal/catalog/testdata/*'); do
	LC_ALL=C awk -F '\t' -v f="$f" 'NR > 1 { print "  " f ": " $1 " (" $2 "): " $5 }' "$f"
done
echo "The scripts of the baseline that internal/verify runs:"
git grep -n -E '"docs/[^"]*\.sh"' -- 'internal/verify/*.go' ':!*_test.go' | sed 's/^/  /'

echo "== (6) The calls of the forge adapter (operation, method), in the non-test Go files of internal/forge/github"
git grep -n -E 'a\.do\(ctx, "' -- 'internal/forge/github/*.go' ':!*_test.go' | sed -E 's/^([^:]*:[0-9]+):.*a\.do\(ctx, "([^"]*)", "([^"]*)", ([^,]*),.*/  \1: \2 \3 \4/'

if [ "$fail" -eq 0 ]; then echo "== PASS"; else echo "== FAIL"; fi
exit "$fail"
