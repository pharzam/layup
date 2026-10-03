#!/bin/sh
# The deterministic part of the release review of phase 1 (task T-efmy, #96:
# D3 of its plan, with conditions 1 and notes 2 and 4 of its plan review).
#
# Run it from the root of a clone of LAYUP with the commit of the release as
# its argument, from a path outside that clone:
#
#   sh /path/to/release-check.sh COMMIT
#
# It prints the facts of the code that REQ-015 and REQ-017 read, and gives
# exit 1 when one of its checks fails:
#   (1) HEAD is not COMMIT, or the tree has a change;
#   (2) TestPackageRules fails: rules 2 to 5 of NFR-007 in
#       docs/spec/packages.md (the standard library and this module only; no
#       package depends on net, net/http or crypto/tls) and the program of each
#       exec.Command (git in internal/git, sh in internal/gate and
#       internal/verify);
#   (4) a non-test Go file of internal/git holds a string literal that equals
#       a git verb that reaches a remote, other than ls-remote and clone of
#       S02: "push", "send-pack", "fetch", "pull" or "remote". A call of git
#       names its verb so, as in do(dir, "clone", ...).
# The lists of (3) and (5) are for the reviewer, who judges them.

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

echo "== (3) The calls that start a program, in the non-test Go files"
git grep -n -E 'exec\.Command|os\.StartProcess|syscall\.(Exec|ForkExec|StartProcess)' -- 'cmd/*.go' 'internal/*.go' ':!*_test.go' | sed 's/^/  /'

echo "== (4) The git verbs that reach a remote, in the non-test Go files of internal/git"
if hits=$(git grep -n -E '"(push|send-pack|fetch|pull|remote)"' -- 'internal/git/*.go' ':!*_test.go'); then
	echo "$hits" | sed 's/^/  /'
	bad "a git verb that reaches a remote"
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

if [ "$fail" -eq 0 ]; then echo "== PASS"; else echo "== FAIL"; fi
exit "$fail"
