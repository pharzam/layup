#!/bin/sh
#
# freeze-check.sh — pre-freeze check of a task worktree (T-b6fn, #440).
#
# Usage:
#   sh docs/tests/freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--closeout LINES/FILES] [--no-conditions]
#
# Reads git and text at HEAD of the task worktree.
# Prints one line per check:
#   ok K<n>: <what>
#   FAIL K<n>: <file>:<line>: <reason>
#   skip K<n>: <reason>
# Exit status: 0 when no check fails, 1 when one or more fail, 2 on a usage or
# environment error (not a git work tree, --base not a commit, an unreadable file).
#
# Stated limit: K6 proves that cited evidence exists, not that it satisfies the
# condition, and it does not know the review's list of conditions; completeness
# stays the author's and the reviewer's (docs/guardrails.md:144).

exit 0
