#!/bin/sh
#
# run-all.sh — run every freeze-check fixture case and enforce the coverage
# and CRLF floors.
#
# Exits 0 only when every case matched its EXPECT, each check K1..K6 has at
# least one good-K<n>-* and at least one bad-K<n>-* case, and every file under
# a *crlf* case keeps its carriage returns.

set -u

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
pass=0
fail=0
g1=0; b1=0; g2=0; b2=0; g3=0; b3=0; g4=0; b4=0; g5=0; b5=0; g6=0; b6=0
crlf_cases=0

for case_dir in "$here"/*/; do
	[ -d "$case_dir" ] || continue
	name=$(basename "$case_dir")
	[ "$name" = base ] && continue
	case $name in
	good-K1-*) g1=$((g1 + 1)) ;; bad-K1-*) b1=$((b1 + 1)) ;;
	good-K2-*) g2=$((g2 + 1)) ;; bad-K2-*) b2=$((b2 + 1)) ;;
	good-K3-*) g3=$((g3 + 1)) ;; bad-K3-*) b3=$((b3 + 1)) ;;
	good-K4-*) g4=$((g4 + 1)) ;; bad-K4-*) b4=$((b4 + 1)) ;;
	good-K5-*) g5=$((g5 + 1)) ;; bad-K5-*) b5=$((b5 + 1)) ;;
	good-K6-*) g6=$((g6 + 1)) ;; bad-K6-*) b6=$((b6 + 1)) ;;
	esac
	case $name in
	*crlf*)
		crlf_cases=$((crlf_cases + 1))
		while IFS= read -r cf; do
			[ -n "$cf" ] || continue
			if ! tr -dc '\r' < "$cf" | grep -q '.'; then
				printf 'FAIL  crlf: %s holds no carriage return\n' "$cf"
				fail=$((fail + 1))
			fi
		done <<EOF
$(find "$case_dir" -type f | sort)
EOF
		;;
	esac
	out=$(sh "$here/run.sh" "$case_dir" 2>&1)
	rc=$?
	if [ "$rc" -eq 0 ]; then
		pass=$((pass + 1))
		printf 'ok    %s\n' "$name"
	else
		fail=$((fail + 1))
		printf 'FAIL  %s\n' "$name"
		printf '%s\n' "$out" | sed 's/^/      /'
	fi
done

if [ "$((pass + fail))" -eq 0 ]; then
	printf 'FAIL  floor: no case ran\n'
	fail=$((fail + 1))
fi

check_k_floor() {
	if [ "$2" -eq 0 ] || [ "$3" -eq 0 ]; then
		printf 'FAIL  floor: %s needs >=1 good-%s-* and >=1 bad-%s-* case (have %d good, %d bad)\n' "$1" "$1" "$1" "$2" "$3"
		fail=$((fail + 1))
	fi
}
check_k_floor K1 "$g1" "$b1"
check_k_floor K2 "$g2" "$b2"
check_k_floor K3 "$g3" "$b3"
check_k_floor K4 "$g4" "$b4"
check_k_floor K5 "$g5" "$b5"
check_k_floor K6 "$g6" "$b6"

if [ "$crlf_cases" -eq 0 ]; then
	printf 'FAIL  crlf: no *crlf* fixture case found\n'
	fail=$((fail + 1))
fi

printf 'freeze-check: %d passed, %d failed\n' "$pass" "$fail"
[ "$fail" -eq 0 ] && exit 0 || exit 1
