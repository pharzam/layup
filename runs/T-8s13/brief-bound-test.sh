#!/bin/sh
# brief-bound-test.sh: pins the terms and L of brief-bound.sh on small committed records (T-8s13, plan conditions 1 and 2).
# Formula (internal/run/betbrief.go, assemble): L = 12 + R + max(M,1) + sum(1+O_d) + Q_io + E + S + A_io + C_op.
# KEEP=DIR keeps each fixture as DIR/<case> for the cross-check against betbrief.go. Exit 0 when every case matches.
set -u
here=$(cd "$(dirname "$0")" && pwd); bb=$here/brief-bound.sh; w=$(mktemp -d); trap 'rm -rf "$w"' EXIT; bad=0
# fx DIR R M "CLASSES": R requirements, M success criteria, one span fact per class word (fact n at bytes n*100 to n*100+90).
fx() { rm -rf "$1"; mkdir -p "$1/spec/options" "$1/intake" "$1/docs/prd"; git init -q "$1"
  { printf 'id\tkind\tstatement\tcovers\tcriterion\tsession\n'; i=0; while [ $i -lt $2 ]; do i=$((i+1)); printf 'R-%03d\trequirement\ts\t1\tc\tS-00000001\n' $i; done
    i=0; while [ $i -lt $3 ]; do i=$((i+1)); printf 'R-%03d\tsuccess-criterion\ts\t1\tc\tS-00000001\n' $(($2+i)); done; } >"$1/spec/draft.tsv"
  { printf 'start\tend\tfact\tclass\treason\ttext\n'; n=0; for c in $4; do n=$((n+1)); printf '%d\t%d\t%d\t%s\t—\tfact\n' $((n*100)) $((n*100+90)) $n $c; done; } >"$1/spec/spans.tsv"
  printf 'id\tsource\tref\tblock\tphase\tkind\tstart\tend\tquestion\tasked\tasker\tanswered\tstate\n' >"$1/questions.tsv"; printf 'question\tanswer\tcomment\tauthor\ttime\n' >"$1/intake/answers.tsv"; }
q() { k=$3; s=layup\ run; t=—; case $3 in meaning) k=two-readings; s=S-00000001;; setup) s=S-00000001;; esac; [ "$8" != answered ] || t=2026-10-12T09:00:00Z
  printf '%s\t%s\t%s\t%s\tbefore-delivery\t%s\t%s\t%s\ttext\t2026-10-12T09:00:00Z\t%s\t%s\t%s\n' "$2" "$3" "$4" "$5" "$k" "$6" "$7" "$s" "$t" "$8" >>"$1/questions.tsv"; }
a() { printf '%s\tans\t1\t1\t2026-10-12T09:00:00Z\n' "$2" >>"$1/intake/answers.tsv"; }
# opts DIR M: one decision db of three options with two constraints each; a PRD with M success-metric rows.
opts() { { printf 'option\tconstraint\tresult\tevidence\n'; for o in pg sqlite mysql; do printf '%s\tNFR-001\tpass\te\n%s\tNFR-002\tpass\te\n' $o $o; done; } >"$1/spec/options/db.tsv"
  { printf '## 8. Success metrics\n\n| Metric | Fact |\n| --- | --- |\n'; i=0; while [ $i -lt $2 ]; do i=$((i+1)); printf '| m%d | F-1 |\n' $i; done; } >"$1/docs/prd/PRD-0001.md"; }
done_() { git -C "$1" add -A; git -C "$1" -c user.name=t -c user.email=t@t commit -q -m x; git -C "$1" update-ref refs/remotes/origin/layup-records HEAD; [ -z "${KEEP:-}" ] || { mkdir -p "$KEEP"; rm -rf "$KEEP/$(basename "$1")"; cp -r "$1" "$KEEP/"; }; }
chk() { name=$1; shift; got=$(sh "$bb" "$@" 2>&1); if [ "$got" = "$want" ]; then echo "ok $name"; else bad=1; printf 'FAIL %s\nwant:\n%s\ngot:\n%s\n' "$name" "$want" "$got"; fi; }
G1() { printf 'R=%s M=%s S=%s A_io=%s C_op=%s Q_io=%s E=%s L0=%s\nG1 floor=%s limit=80 stop=%s\n' "$@"; }

d=$w/io-open;   fx $d 2 0 "need need"; q $d Q-001 meaning M-001 idea-owner 110 120 open; q $d Q-002 setup M-1a2b3c4d operator — — answered; a $d Q-002; done_ $d
want=$(G1 2 0 2 0 1 1 0 19 22 no); chk "only an idea-owner question open: E=0, Q_io=1" $d
d=$w/none-open; fx $d 2 0 "need need"; q $d Q-001 meaning M-001 idea-owner 110 120 answered; a $d Q-001; q $d Q-002 setup M-1a2b3c4d operator — — answered; a $d Q-002; done_ $d
want=$(G1 2 0 2 1 1 0 1 20 23 no); chk "no question open: E=1 for the line - none" $d
d=$w/op-open;   fx $d 2 0 "need need"; q $d Q-001 meaning M-001 idea-owner 110 120 open; q $d Q-002 setup M-1a2b3c4d operator — — open; done_ $d
want=$(G1 2 0 2 0 0 1 1 19 22 no); chk "an Operator and an idea-owner question open: E=1, Q_io=1" $d
d=$w/followups; fx $d 2 1 "need constraint"; q $d Q-001 meaning M-001 operator 210 220 asked; q $d Q-002 follow-up Q-001 operator — — answered
q $d Q-003 setup M-1a2b3c4d operator — — asked; q $d Q-004 follow-up Q-003 operator — — asked; q $d Q-005 setup M-1a2b3c4d operator — — answered; a $d Q-001; a $d Q-002; a $d Q-005; done_ $d
want=$(G1 2 1 2 0 2 0 1 20 23 no); chk "answered follow-up closes its question, asked one keeps it open; C_op by quote class" $d
d=$w/io-follow; fx $d 2 0 "need"; q $d Q-001 meaning M-001 idea-owner 110 120 asked; q $d Q-002 follow-up Q-001 idea-owner — — asked; done_ $d
want=$(G1 2 0 1 0 0 1 0 17 20 no); chk "idea-owner question with an unanswered follow-up is open" $d
# The 80/81 boundary of G2 (room >= 20 at brief.lines 100) and the floor of G1 (L0 + 3: one decision with two options).
d=$w/g2-80; fx $d 61 0 ""; opts $d 2; done_ $d
want="$(G1 61 0 0 0 0 0 1 75 78 no)
D=1 O=3 M_prd=2 L=80 room=20 G2=pass"; chk "G2 at L=80: room 20 passes" $d HEAD
d=$w/g2-81; fx $d 62 0 ""; opts $d 2; done_ $d
want="$(G1 62 0 0 0 0 0 1 76 79 no)
D=1 O=3 M_prd=2 L=81 room=19 G2=fail"; chk "G2 at L=81: room 19 fails" $d HEAD
d=$w/g1-79; fx $d 62 1 ""; done_ $d
want=$(G1 62 1 0 0 0 0 1 76 79 no); chk "G1 floor 79 (76 plus one decision of two options) passes" $d
d=$w/g1-81; fx $d 64 1 ""; done_ $d
want=$(G1 64 1 0 0 0 0 1 78 81 yes); chk "G1 floor 81 stops" $d
d=$w/g2-81tab; fx $d 62 0 ""; opts $d 2; tab=$(printf '\t'); sed -i "s/^| Metric | Fact |\$/|${tab}Metric${tab}|${tab}Fact${tab}|/" $d/docs/prd/PRD-0001.md; done_ $d
want="$(G1 62 0 0 0 0 0 1 76 79 no)
D=1 O=3 M_prd=2 L=81 room=19 G2=fail"; chk "tabs around the table header cells count as spaces (internal/spec/full.go tableCells)" $d HEAD
# A failed git read must end the script non-zero before any G1 or G2 line that could pass (guardrails, check each read's status).
chkfail() { name=$1; shift; got=$(sh "$bb" "$@" 2>&1); rc=$?
  case $got in *G2=pass*|*stop=no*) rc=0;; esac
  if [ "$rc" -ne 0 ]; then echo "ok $name"; else bad=1; printf 'FAIL %s\nwant: non-zero exit, no G2=pass, no stop=no\ngot:\n%s\n' "$name" "$got"; fi; }
chkfail "an unreadable head fails" $w/g2-81 no-such-verified-head
d=$w/noref; fx $d 62 0 ""; opts $d 2; git -C $d add -A; git -C $d -c user.name=t -c user.email=t@t commit -q -m x
chkfail "a missing records ref fails" $d HEAD
exit $bad
