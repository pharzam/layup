#!/bin/sh
# gh-stub.sh — a stand-in for `gh` in the applier cases of task-labels.sh. It
# records each call in $GH_LOG, serves the three files and the issues from
# $GH_STUB_DIR, accepts a label create or an issue edit (a write), and fails
# any other call. A missing `forge-issues` fails the issue read.
printf '%s\n' "$*" >> "$GH_LOG"
case "$*" in
*contents/docs/tasks/backlog.md*) cat "$GH_STUB_DIR/backlog" ;;
*contents/docs/tasks/completed.md*) cat "$GH_STUB_DIR/completed" ;;
*contents/docs/plan/README.md*) cat "$GH_STUB_DIR/plan" ;;
*issues\?state=all*) [ -f "$GH_STUB_DIR/forge-issues" ] && cat "$GH_STUB_DIR/forge-issues" || exit 1 ;;
"label create "*|"issue edit "*) : ;;
*) exit 1 ;;
esac
