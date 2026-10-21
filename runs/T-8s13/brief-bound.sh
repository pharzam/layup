#!/bin/sh
# brief-bound.sh CLONE [HEAD]: the lines of the parts of code in the brief of the first bet, from origin/layup-records of a plain clone (T-8s13).
# L = 12 + R + max(M,1) + sum(1+O_d) + Q_io + E + S + A_io + C_op, as internal/run/betbrief.go (assemble) writes it. Without HEAD (G1): L0 with no
# decision and the floor L0+3 (one decision, two options). With HEAD, a commit of CLONE with the verified PRD (G2): M from its success-metric
# table, D and O from spec/options/*.tsv, L and room. Reads git only; brief.lines is LINES (default 100); G2 needs room >= 20.
set -eu
clone=$1; head=${2:-}; lines=${LINES:-100}; r=origin/layup-records; t=$(mktemp -d); trap 'rm -rf "$t"' EXIT
for f in spec/draft.tsv spec/spans.tsv questions.tsv intake/answers.tsv; do git -C "$clone" show "$r:$f" >"$t/$(basename "$f")"; done
eval "$(awk -F'\t' -v D='—' 'FNR==1{f++; next}
f==1{if($2=="requirement"||$2=="non-functional")R++; if($2=="success-criterion")M++}
f==2&&$3!=D{S++; lo[$3]=$1+0; hi[$3]=$2+0; cl[$3]=$4; fo[++nf]=$3}
f==3{n++; id[n]=$1; src[$1]=$2; ref[$1]=$3; blk[$1]=$4; st[$1]=$7; state[$1]=$13; if($2=="follow-up"&&$13=="answered")again[$3]=1}
f==4{ if(blk[$1]!="operator"){A++; next}
  q=$1; for(h=0;h<n&&src[q]=="follow-up"&&st[q]==D;h++) q=ref[q]; c="context"
  if((q in src)&&st[q]!=D) for(i=1;i<=nf;i++) if(st[q]+0>=lo[fo[i]]&&st[q]+0<hi[fo[i]]){c=cl[fo[i]]; break}
  if(!(c in seen)){seen[c]=1; C++}}
END{for(i=1;i<=n;i++){x=id[i]; if(src[x]!="follow-up"&&(state[x]=="open"||(state[x]=="asked"&&!(x in again)))){no++; if(blk[x]=="operator")op=1; else Q++}}
  printf "R=%d;M=%d;S=%d;A=%d;C=%d;Q=%d;E=%d\n", R, M, S, A, C, Q, (op||no==0)}' "$t/draft.tsv" "$t/spans.tsv" "$t/questions.tsv" "$t/answers.tsv")"
mx() { [ "$1" -gt 1 ] && echo "$1" || echo 1; }
L0=$((12 + R + $(mx "$M") + Q + E + S + A + C))
echo "R=$R M=$M S=$S A_io=$A C_op=$C Q_io=$Q E=$E L0=$L0"
echo "G1 floor=$((L0 + 3)) limit=$((lines - 20)) stop=$([ $((L0 + 3)) -gt $((lines - 20)) ] && echo yes || echo no)"
[ -n "$head" ] || exit 0
Mp=$(for p in $(git -C "$clone" ls-tree --name-only "$head" docs/prd/ | grep '/PRD-.*\.md$' | sort); do git -C "$clone" show "$head:$p"; done |
  awk '/^#/{l=match($0,/[^#]/)-1; if(l<=2)m=(l==2&&/Success metrics/); p=""; ok=0; next}
  /^[ \t]*\|/{if($0~/^[ \t|:-]+$/&&p!=""){ok=(m&&p~/\| *Fact *(\||$)/); next} if(ok)n++; p=$0; next} {p=""; ok=0} END{print n+0}')
Dn=0; Ot=0
for p in $(git -C "$clone" ls-tree --name-only "$r" spec/options/ | grep '\.tsv$'); do
  o=$(git -C "$clone" show "$r:$p" | awk -F'\t' 'NR>1&&!s[$1]++{n++} END{print n+0}'); Dn=$((Dn + 1)); Ot=$((Ot + o)); done
L=$((12 + R + $(mx "$Mp") + Dn + Ot + Q + E + S + A + C))
echo "D=$Dn O=$Ot M_prd=$Mp L=$L room=$((lines - L)) G2=$([ $((lines - L)) -ge 20 ] && echo pass || echo fail)"
