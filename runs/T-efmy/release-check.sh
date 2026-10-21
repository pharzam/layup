#!/bin/sh
# The deterministic part of the release review of phase 1 (task T-efmy, #96:
# D3 of its plan, with conditions 1 and notes 2 and 4 of its plan review),
# adapted for the release of M2b (task T-4tjy, #169: D3 of its plan, with
# conditions 3 and 5 and notes 3, 5 and 6 of its plan review;
# docs/spec/session.md, REQ-015 and REQ-017 — The review of the release of M2b),
# and adapted again for the release of M2d (task T-5xjc, #283: D3 of its plan,
# with conditions 1 and 4 of its plan review; docs/spec/scaffold.md, REQ-015
# and REQ-017 — The review of the release of M2d), and adapted again for the
# release of M2e (task T-ym0k, #369: D2 to D4 of its plan, with conditions 1, 2
# and 7 of its plan review; docs/spec/loop.md, REQ-015 and REQ-017 — The
# review of the release of M2e).
# Its output at 621af09 is runs/T-efmy/release-check.txt, from the script
# before T-4tjy; its output at the release of M2b is
# runs/T-4tjy/release-check.txt; the script is unchanged for the release of M2c
# (runs/T-8byd/release-check.txt); its output at the release of M2d is
# runs/T-5xjc/release-check.txt; its output at the release of M2e is
# runs/T-ym0k/release-check.txt.
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
#       internal/verify/scripts.go; sh with the form exec.Command("sh", script,
#       prdDir) in internal/spec/prdlint.go (the pinned prd-lint, D2 of T-ym0k:
#       the call, not the file, so a "-c" form or a second call fails); the first word of a harness register row's
#       version command (words[0]) or command (s.Words[0]) in
#       internal/session/process.go. The specification's "a new list names
#       each program that the code starts" is read as a check that fails, the
#       chosen reading (note 3 of the plan review of #169);
#   (4) a non-test Go file of internal/git holds a string literal, interpreted
#       ("push") or raw (`push`), that equals a git verb that reaches a remote:
#       push, send-pack, fetch, pull or remote, other than these, each in its
#       own function: fetch in Fetch and FetchSession, push in Push, remote in
#       CloneLocal (remote remove origin, which cuts a session clone's link and
#       reaches no remote) and remote in AddRemote (remote add, which writes a
#       name and a URL into the config of a local clone; by verb, so another
#       subcommand of git remote in these two functions passes). A call of git
#       names its verb so, as in do(dir, "clone", ...). The function of a hit is
#       the last line that starts with "func" above it; a hit after a "}" line at
#       column one, or after a one-line function (a "func" line that ends with
#       "}"), and before the next "func" line is in no function, and fails;
#   (7) a non-test Go file of cmd/ or internal/ out of internal/git and
#       testdata/ refers to git.Push (a call, or a value such as f := git.Push)
#       in a function other than PushHead or commitFiles of
#       internal/run/store.go (the heads of sessions and the records branch) or
#       pushProbe of internal/run/forgecheck.go (the probe of the forge check,
#       which the Operator's ruleset must refuse); or imports internal/git
#       under an explicit name (the selector git.Push would not show it).
#       The function is found as in (4), by its name. The limits: a second call
#       of git.Push in one of the three functions, and a new caller of one of
#       them, pass (7), so the review reads those functions and their callers;
#       (7) reads a line at a time and skips a comment line; the import guard
#       passes when it finds no import at all; a push by a program started with
#       the text "git push" is check (3)'s. The other calls of internal/git
#       that name a remote verb (LsRemote, RemoteHead, Clone, Fetch,
#       CloneLocal, AddRemote and FetchSession) are listed with their callers
#       and fail nothing.
#   (6) lists the calls of the forge adapter, one row for each `.do(` in a
#       non-test Go file of internal/forge/github, and fails when one `.do(`
#       is not in the form a.do(ctx, "Operation", "METHOD", target, ... on one
#       line (a variable operation or method, another receiver, two calls on
#       one line), or when no `.do(` is found: the rows it lists are the
#       parsed ones. The receiver is the name a alone: the byte before a.do
#       must be a space, a tab or one of ( , ; = ! & | { : (or a.do starts the
#       line), and anything else, any non-ASCII byte included (aa, x.a, or a
#       name that starts with a letter outside ASCII), is another receiver and
#       fails. The two expressions run with LC_ALL=C, so bytes are compared,
#       not the characters of a locale. The limit, retained: that the form parses does not show
#       that the operation name or the method fits what the call does, or that
#       do is the only way to the network; the reviewer judges each row.
# Checks (3), (4) and (7) print a list and fail on an entry out of their
# allowance, which the reviewer judges; check (5) only lists.

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
	# The allowed call of prdlint.go is the exact call and nothing else on the
	# line: no text of a comment before it, and no other start of a program
	# (or comment text that names one) around it.
	function allowedcall(l,   call, i, pre, post) {
		call = "exec.Command(\"sh\", script, prdDir)"
		i = index(l, call)
		if (i == 0) return 0
		pre = substr(l, 1, i - 1); post = substr(l, i + length(call))
		sub(/^[^:]*:[0-9]+:/, "", pre)
		if (pre ~ /\/\// || index(post, call) > 0) return 0
		return (pre post) !~ /exec\.Command|os\.StartProcess|syscall\./
	}
	/^internal\/git\/git\.go:[0-9]+:.*exec\.Command\("git", / { next }
	/^internal\/gate\/scratch\.go:[0-9]+:.*exec\.Command\("sh", / { next }
	/^internal\/spec\/prdlint\.go:[0-9]+:/ && allowedcall($0) { next }
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
# git ls-files (2.47) lists no file for one positive pathspec with the
# exclusion ':!*_test.go', so the test files are filtered out by name; a list
# with no file fails, as no file read would pass.
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
		}
		/^func .*}[ \t]*$/ { fn = "" }' "$f"
done)
if [ -n "$hits" ]; then
	echo "$hits" | sed 's/^/  /'
	bad4=$(echo "$hits" | grep -v -E ': (fetch in (Fetch|FetchSession)|push in Push|remote in CloneLocal|remote in AddRemote)$')
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
# The occurrences of ".do(" and the rows that parse: a difference is a call in
# a form the rows cannot show. Two empty counts are not coverage.
dosites=$(git grep -n -o -E '\.do\(' -- 'internal/forge/github/*.go' ':!*_test.go')
dolines=$(git grep -n -E '\.do\(' -- 'internal/forge/github/*.go' ':!*_test.go')
tb=$(printf '\t')
bnd="[ $tb(,;=!&|{:]" # the characters that may stand before the receiver a
dorows=$(echo "$dolines" | LC_ALL=C sed -n -E 's/^([^:]*:[0-9]+):(.*'"$bnd"')?a\.do\(ctx, "([^"]*)", "(GET|POST|PUT|PATCH|DELETE)", ([^,]*),.*/  \1: \3 \4 \5/p')
[ -n "$dosites" ] || bad "no .do( call found in internal/forge/github"
[ -z "$dorows" ] || echo "$dorows"
nsites=$(echo "$dosites" | grep -c .)
nrows=$(echo "$dorows" | grep -c .)
echo "  $nsites .do( occurrences, $nrows rows parsed"
if [ "$nsites" -ne "$nrows" ]; then
	# The loop runs in a subshell of the pipe, so fail is set here.
	echo "$dolines" | LC_ALL=C grep -v -E "$bnd"'a\.do\(ctx, "[^"]*", "(GET|POST|PUT|PATCH|DELETE)", [^,]*,' | while IFS= read -r l; do echo "FAIL: a .do( call that is not in the form a.do(ctx, \"Operation\", \"METHOD\", target, on one line: $l"; done
	echo "$dosites" | cut -d: -f1,2 | sort | uniq -d | while IFS= read -r l; do echo "FAIL: a line with two .do( calls: $l"; done
	bad "the .do( occurrences ($nsites) and the rows parsed ($nrows) differ"
fi

echo "== (7) The calls of internal/git that name a remote verb, and who calls them (non-test Go files of cmd/ and internal/ out of internal/git)"
callfiles=$(git ls-files 'cmd/*.go' 'internal/*.go' | grep -v -E '_test\.go$|/testdata/|^internal/git/')
[ -n "$callfiles" ] || bad "no non-test Go file out of internal/git"
refs=$(for f in $callfiles; do
	awk -v f="$f" '
		/^func / { fn = $0; sub(/^func (\([^)]*\) )?/, "", fn); sub(/[^A-Za-z0-9_].*/, "", fn) }
		/^}/ { fn = "" }
		/^[ \t]*\/\// { next }
		{
			line = $0 " "
			while (match(line, /git\.(LsRemote|RemoteHead|Clone|Fetch|Push|CloneLocal|AddRemote|FetchSession)[^A-Za-z0-9_]/)) {
				call = substr(line, RSTART + 4, RLENGTH - 5)
				print f ":" NR ": " call " in " (fn == "" ? "no function" : fn)
				line = substr(line, RSTART + RLENGTH - 1)
			}
		}
		/^func .*}[ \t]*$/ { fn = "" }' "$f"
done)
if [ -n "$refs" ]; then
	echo "$refs" | sed 's/^/  /'
	badpush=$(echo "$refs" | grep -E ': Push in ' | grep -v -E '^internal/run/store\.go:[0-9]+: Push in (PushHead|commitFiles)$|^internal/run/forgecheck\.go:[0-9]+: Push in pushProbe$')
	if [ -n "$badpush" ]; then
		# The loop runs in a subshell of the pipe, so fail is set here.
		echo "$badpush" | while IFS= read -r l; do echo "FAIL: a reference to git.Push, out of its function: $l"; done
		fail=1
	fi
else
	bad "no reference to a call of internal/git that names a remote verb found"
fi
aliased=$(git grep -n -F '"github.com/pharzam/layup/internal/git"' -- 'cmd/*.go' 'internal/*.go' ':!*_test.go' ':!*/testdata/*' | awk '{ line = $0; sub(/^[^:]*:[0-9]+:/, "", line); sub(/"github\.com\/pharzam\/layup\/internal\/git".*/, "", line); sub(/^[ \t]*(import[ \t]+)?/, "", line); if (line != "") print }')
if [ -n "$aliased" ]; then
	# The loop runs in a subshell of the pipe, so fail is set here.
	echo "$aliased" | while IFS= read -r l; do echo "FAIL: internal/git imported under an explicit name: $l"; done
	fail=1
fi

if [ "$fail" -eq 0 ]; then echo "== PASS"; else echo "== FAIL"; fi
exit "$fail"
