#!/bin/sh
# prep-test.sh [PREPARER [TEMPLATE]]: named cases of the round preparation (T-wqrk), offline.
# A local bare repository stands for the forge clone URL (git url rewriting) and a stub gh answers the forge reads.
# Default PREPARER and TEMPLATE are the live harness files. Prints PASS or FAIL per case; exit 1 if any FAIL.
S=/home/layup/.cache/layup-session/m2e-scratchpad/rounds
P=${1:-$S/prep-round.py}; export PREP_TEMPLATE=${2:-$S/round-brief.tmpl}
T=$(mktemp -d) || exit 2; trap 'rm -rf "$T"' EXIT
M=$(git rev-parse HEAD); git clone -q --bare "$(git rev-parse --show-toplevel)" "$T/o.git" || exit 2
W="$T/w"; git clone -q "$T/o.git" "$W"; git -C "$W" config user.name t; git -C "$W" config user.email t@example.invalid
git -C "$W" checkout -q -b T-pt01 "$M"; echo first > "$W/pt.txt"; git -C "$W" add pt.txt; git -C "$W" commit -q -m a
A=$(git -C "$W" rev-parse HEAD); echo second >> "$W/pt.txt"; git -C "$W" commit -q -am b
B=$(git -C "$W" rev-parse HEAD); git -C "$W" push -q origin T-pt01; git -C "$T/o.git" update-ref refs/heads/main "$M"
mkdir "$T/bin"; printf '#!/bin/sh\necho "BODY $2"\n' > "$T/bin/gh"; chmod +x "$T/bin/gh"
export PATH="$T/bin:$PATH" GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0="url.$T/o.git.insteadOf" GIT_CONFIG_VALUE_0=https://github.com/pharzam/layup.git
bad=0; n=0
cfg() { # NAME ROUND HEAD EXTRAKEYS
	n=$((n+1)); printf '{"rundir":"%s/r%s","task":"T-pt01","issue":7,"round":%s,"cycle":1,"cap":1,"what":"w","reviewer":"R","harness":"H","exc":"","authors":"A","planextra":"%s","head":"%s","main":"%s","budget":"b","plan":11,"plan_review":12%s}' "$T" "$n" "$2" "${PX:-p}" "$3" "$M" "$4" > "$T/c$n.json"; }
prep() { python3 "$P" "$T/c$n.json" > "$T/out$n" 2>&1; }
D() { echo "$T/r$n/layup/.review-in"; }
ok() { if [ "$2" = 0 ]; then echo "PASS $1"; else echo "FAIL $1"; bad=1; fi; }
cfg x 1 "$A"; prep; ok "round1: prepares, no fix.diff, earlier-review claim kept" "$( [ $? = 0 ] && [ ! -e "$(D)/fix.diff" ] && ! grep -q 'fix.diff' "$(D)/brief.md" && grep -q 'no earlier review' "$(D)/brief.md"; echo $? )"
cfg x 2 "$B" ',"fix_base":"'$A'"'; prep; ok "fix round: fix.diff and the search block" "$( grep -q '^+second' "$(D)/fix.diff" && ! grep -q '^+first' "$(D)/fix.diff" && grep -q 'git grep -n -i -F' "$(D)/brief.md" && grep -q 'Homes of a rule the fix changed' "$(D)/brief.md" && ! grep -q '@' "$(D)/brief.md"; echo $? )"
PX='a@b' cfg x 1 "$A"; prep; ok "unresolved @ in the text is refused" "$( [ $? != 0 ]; echo $? )"
cfg x 2 "$B"; prep; ok "round 2 without fix_base is refused" "$( [ $? != 0 ]; echo $? )"
cfg x 2 "$B" ',"fix_base":"'$A'","fix_exclude":["docs/engineering-discipline.md"]'; prep; ok "fix_exclude outside runs/TASK/NAME.md is refused" "$( [ $? != 0 ]; echo $? )"
cfg x 2 "$B" ',"extra":{"round1.md":13}'; prep; ok "an earlier verdict is refused without earlier_record" "$( [ $? != 0 ]; echo $? )"
cfg x 2 "$B" ',"fix_base":"'$A'","extra":{"round1.md":13},"earlier_record":true'; prep; ok "an earlier verdict handed over is stated in the brief" "$( grep -q "handed an earlier round's record" "$(D)/brief.md"; echo $? )"
exit $bad
