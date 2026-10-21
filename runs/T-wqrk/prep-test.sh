#!/bin/sh
# prep-test.sh [PREPARER [TEMPLATE]]: named cases of the round preparation (T-wqrk), offline.
# A local bare repository stands for the clone URL (git url rewriting) and a stub gh answers the forge reads.
# Every refusal case is an otherwise valid config and asserts its own message. Default files are the live harness files.
# Prints PASS or FAIL per case; exit 1 if any FAIL.
S=/home/layup/.cache/layup-session/m2e-scratchpad/rounds
P=${1:-$S/prep-round.py}; export PREP_TEMPLATE=${2:-$S/round-brief.tmpl}
T=$(mktemp -d) || exit 2; trap 'rm -rf "$T"' EXIT
M=$(git rev-parse HEAD); git clone -q --bare "$(git rev-parse --show-toplevel)" "$T/o.git" || exit 2
W="$T/w"; git clone -q "$T/o.git" "$W"; git -C "$W" config user.name t; git -C "$W" config user.email t@example.invalid
git -C "$W" checkout -q -b T-pt01 "$M"; mkdir "$W/runs/T-pt01"; cd "$W/runs/T-pt01" || exit 2
for f in rec.md rule.md data.md full.md; do echo first > $f; done; echo 'cat runs/T-pt01/full.md' > ../../pt-full.sh; chmod +x rule.md; echo 'cat data.md' > read.sh; echo first > ../../pt.txt
git -C "$W" add -A; git -C "$W" commit -q -m a; A=$(git -C "$W" rev-parse HEAD)
for f in rec.md rule.md data.md full.md ../../pt.txt; do echo second >> $f; done; git -C "$W" commit -q -am b
B=$(git -C "$W" rev-parse HEAD); git -C "$W" push -q origin T-pt01; git -C "$T/o.git" update-ref refs/heads/main "$M"; cd "$W" || exit 2
mkdir "$T/bin"; printf '#!/bin/sh\necho "BODY $2"\n' > "$T/bin/gh"; chmod +x "$T/bin/gh"
export PATH="$T/bin:$PATH" GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0="url.$T/o.git.insteadOf" GIT_CONFIG_VALUE_0=https://github.com/pharzam/layup.git
bad=0; n=0
cfg() { # ROUND HEAD EXTRAKEYS; PX overrides the planextra
	n=$((n+1)); printf '{"rundir":"%s/r%s","task":"T-pt01","issue":7,"round":%s,"cycle":1,"cap":1,"what":"w","reviewer":"R","harness":"H","exc":"","authors":"A","planextra":"%s","head":"%s","main":"%s","budget":"b","plan":11,"plan_review":12%s}' "$T" "$n" "$1" "${PX:-p}" "$2" "$M" "$3" > "$T/c$n.json"; }
prep() { python3 "$P" "$T/c$n.json" > "$T/out$n" 2>&1; }
D() { echo "$T/r$n/layup/.review-in"; }
ok() { if [ "$2" = 0 ]; then echo "PASS $1"; else echo "FAIL $1"; bad=1; fi; }
refused() { prep; ok "$1" "$( [ $? != 0 ] && grep -q -- "$2" "$T/out$n"; echo $? )"; } # nonzero exit and the message
V=',"fix_base":"'$A'"'
cfg 1 "$A"; prep; ok "round1: prepares, no fix.diff, earlier-review claim kept" "$( [ $? = 0 ] && [ ! -e "$(D)/fix.diff" ] && ! grep -q 'fix.diff' "$(D)/brief.md" && grep -q 'no earlier review' "$(D)/brief.md"; echo $? )"
cfg 2 "$B" "$V"; prep; ok "fix round: fix.diff, the search block, nothing left out" "$( grep -q '^+second' "$(D)/fix.diff" && grep -q 'rule.md' "$(D)/fix.diff" && ! grep -q '^+first' "$(D)/fix.diff" && grep -q "git grep -n -i -F -e TERM -- \\." "$(D)/brief.md" && grep -q 'Homes of a rule the fix changed' "$(D)/brief.md" && ! grep -q '@' "$(D)/brief.md"; echo $? )"
PX='a@b' cfg 2 "$B" "$V"; refused "an unresolved @ in the text is refused" AssertionError
cfg 2 "$B"; refused "round 2 without fix_base is refused" "needs fix_base"
cfg 2 "$B" ',"fix_base":"abc"'; refused "a malformed fix_base is refused" "full SHA"
cfg 2 "$B" ',"fix_base":"'$B'"'; refused "a fix_base equal to the head is refused" "the head itself"
cfg 2 "$A" ',"fix_base":"'$B'"'; refused "a fix_base that is not an ancestor is refused" "not an ancestor"
cfg 2 "$B" "$V"',"extra":{"round1.md":13}'; refused "an earlier verdict in extra is refused, with a valid fix_base" "earlier verdict"
cfg 2 "$B" "$V"',"extra":{"round1.md":13},"earlier_record":true'; refused "the earlier_record opt-in no longer lets a verdict in" "earlier verdict"
cfg 2 "$B" "$V"',"passive":["runs/T-pt01/rec.md"]'; prep; ok "a named passive record leaves the diff and the search, and is named" "$( ! grep -q 'rec.md' "$(D)/fix.diff" && grep -q 'rule.md' "$(D)/fix.diff" && grep -q 'data.md' "$(D)/fix.diff" && grep -q "':!runs/T-pt01/rec.md'" "$(D)/brief.md"; echo $? )"
cfg 2 "$B" "$V"',"passive":["docs/engineering-discipline.md"]'; refused "a passive path outside runs/TASK/NAME.md is refused" "passive takes only"
cfg 2 "$B" "$V"',"passive":["runs/T-pt01/rule.md"]'; refused "an executable record cannot be passive" "executable"
cfg 2 "$B" "$V"',"passive":["runs/T-pt01/data.md"]'; refused "a record that a script reads cannot be passive" "consumed"
cfg 2 "$B" "$V"',"passive":["runs/T-pt01/full.md"]'; refused "a record that a script reads by its full path cannot be passive" "consumed"
exit $bad
