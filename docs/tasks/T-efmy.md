# T-efmy — the release review of phase 1 for the `Won't` rows

Issue: [#96](https://github.com/pharzam/layup/issues/96), row 19 of the
[implementation plan](../plan/README.md); child of the core engine (#29).
Serves `REQ-015` and `REQ-017`. Base `621af09` (the merge of row 18, #118).
Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-efmy/`](../../runs/T-efmy/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #96. The
reviewer order is the Operator's (Devin, then OpenCode, then Claude). Devin
(its usage quota) and OpenCode (Grok 4.7: no output in five minutes) gave no
record. The plan review (Claude Fable 5.1, effort `xhigh`, on the Claude Code
CLI, a fresh read-only session in a clone at `621af09`) gave
`approve-with-conditions`: Budget maximum 700 lines added plus removed over 12
files against the base, close-out inside; Cycle cap 1. Its two conditions are
applied: the check of the `git` verbs matches a Go string literal that equals a
verb that reaches a remote, other than `ls-remote` and `clone` of S02 (1); the
scenario names the accepting party by its model, its harness and the date (2).
Its notes are applied, except note 10 (a glossary line for "release": the word
is not an abbreviation, and the sentence of D1 is its one home).

## What was done

1. **The release of phase 1** (D1): the code of `main` that the first pilot
   runs (the non-test Go files of `cmd/` and `internal/`, `go.mod`, and the
   files that the binary embeds), at the commit that the release review reads,
   `621af09`. It is not a build, as phase 1 ships no binary; a later milestone
   can define a versioned release. A task of phase 1 that changes that code
   after the review checks its own diff in its review round, with the check on
   its head. Both sentences are in `docs/plan/README.md`, under the paragraph
   of the first pilot.
2. **The reading of the two rows** (D2): a code path breaks `REQ-015` when it
   holds a machine-learning library, a call to a model or a training service, a
   program that trains, or a write of model files; it breaks `REQ-017` when
   LAYUP itself opens a connection to a cloud provider or a hosting platform, or
   starts a program that changes its infrastructure, its services or its
   settings. The commands of `commands.sh` (the pushes of S03, S13 and S15, and
   the ruleset of S13) are the Operator's, run with the Operator's login on the
   Operator's repository; `layup` writes them and never runs them. Reason:
   `PRD-0001` and the architecture were approved together (O-129), and the
   architecture gives `layup run` a forge adapter from phase 2 (§1, the six
   capabilities), which a reading that forbids each write to the forge would
   make a `Won't` row. The specification task of `M2a` can cite this reading.
3. **The release check** (D3): [`release-check.sh`](../../runs/T-efmy/release-check.sh),
   the deterministic part of the claim (R5). It takes the commit as its
   argument, and gives exit 1 on another `HEAD` or a change in the tree, on a
   failure of `TestPackageRules`, or on a `git` verb that reaches a remote. It
   lists the calls that start a program and the commands that `sh` runs. Three
   mutation runs fail, each for its own reason; at `621af09` it passes
   ([`release-check.txt`](../../runs/T-efmy/release-check.txt)).
4. **The release review** (D4, the uat): Claude Fable 5.1, a fresh read-only
   session in a clone at `621af09`, recorded `holds` for `REQ-015` and for
   `REQ-017`, with ten notes, on #96 and in
   [`release-review.md`](../../runs/T-efmy/release-review.md). A fresh session
   of a model other than the authors' stands for the person of
   `docs/tests/template-uat.md`, by row 19 of the plan, the external input of
   `gov-release-review` ("A reviewer whose model differs from the authors'")
   and Bootstrap mode rule 3.
5. **The documents** (D5): the §12 Test cells of `REQ-015` and `REQ-017` in
   `PRD-0001` with a §13 row, the traceability row of the release review
   (`green`), and the sentences of D1 in the plan.

**The rejected alternatives:** a release as a tagged build (no source names
one, and phase 1 ships no binary); a new CI check for the two rows
(`TestPackageRules` already holds rules 2 to 5 of `NFR-007` on each change); a
heading `## Review record` for the release review (`review-record-lint` would
read it as a round of this task and refuse it); the Operator as the person of
the scenario (row 19 and the inventory give the review to a reviewer whose
model differs from the authors').

**Known limits:** the release review covers the code at `621af09` only; a later
change of that code in phase 1 is checked by the sentence of D1, and each later
release by the specification task of its milestone. The notes of the release
review need no change: `internal/records` imports `net/url`, which pulls in
`net/netip`; both parse text and open no connection, and rule 5 of `NFR-007`
names `net`, `net/http` and `crypto/tls` (note 5). Notes 6 to 10 name the paths
that a stricter reading would catch (the reads of S02, the commands of
`commands.sh`, the workflow that S12 writes, the gate command and the baseline
scripts, the strings that name GitHub); under the reading of D2 none is a call
that modifies a platform.

**Lessons:** none new for `guardrails.md` §2.

## Review round 1

The first try gave no record. Devin (its usage quota) and OpenCode (no output
in five minutes) were skipped, and Claude Fable 5.1 gave no record in fifteen
minutes: its answer was cut off when it began the record, and the API retried
until the run was stopped (Bootstrap mode rule 4). In a second try, Devin and
OpenCode were skipped again, and Claude Fable 5.1 (effort `xhigh`, on the
Claude Code CLI, a fresh read-only session in a fresh clone at `5eb1794`, its
record at 6 min 44 s; the record is on #96) gave `nothing material in scope`,
with seven notes. CI of PR #119 on `5eb1794` passed its job `tests` before the
round. The round ran the check and ten mutation and edge-case runs of its own,
and checked 33 `file:line` citations of the release review at `621af09`.

The notes. Note 1 is applied: check (4) matched only an interpreted string
literal (`"push"`), and now matches a raw one too, with a fourth mutation run.
Note 5 is applied: the header of the check and the sentence of D1 name the full
commit ID, as a name such as `HEAD` passes check (1) on any commit. Note 3 is
checked: the copy of the release review is the text of its comment on #96.
Note 6 is the close-out. Notes 2, 4 and 7 need no change: the row of the
release review gives the fact of the first requirement, the rule of the table;
check (5) lists each entry, and the `go:embed` lines say which the binary
embeds; and the hits on the safe side are those that the answer to condition 1
gives.

## Verdict

Delivered: the release review of phase 1, which records `holds` for `REQ-015`
and `REQ-017` for the code of `main` at `621af09` (Claude Fable 5.1, a fresh
session of a model other than the authors'), with its deterministic part, the
release check, and the phase-1 release defined in the plan. The plan review
(Claude Fable 5.1) gave `approve-with-conditions`, with two conditions,
applied; round 1 (`5eb1794`) gave `nothing material in scope`, in its second
try. The records are on #96; `review-record-lint` reads the round record in
CI. At `5eb1794` the ladder passes (18 steps, each exit 0) and CI passed its
job `tests`. The close-out commit applies notes 1 and 5 of round 1 to the check
and to the sentence of D1, and changes no Go code. The diff against
`origin/main` is 486 lines over 10 files, inside 700 over 12.

Next: row 20 (`T-evad`, #97), the first pilot, whose After cell (rows 16 to
19) is now met. It needs the Operator: its inputs, a problem statement with a
Go stack and an empty public repository of the Operator's account, come at its
start (O-122).

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC, 2026-10-03; the token count
is the `result` event of the Claude Code CLI; `not reported` otherwise.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan (D1 to D6) | reasoning | Claude Opus 5.5 | max | not reported | 04:03 to 04:09 |
| The plan review, first and second harness: skipped | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `xhigh` | not reported | 04:09:45 to 04:16:04 |
| The plan review | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 578,316 (USD 4.15) | 6 min 42 s, from 04:16:10 |
| The answer; the release check, test first (three mutation runs), and its run at `621af09` | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 04:23 to 04:26 |
| The release review, first and second harness: skipped | reasoning | the same | `xhigh` | not reported | 04:26:30 to 04:32:49 |
| The release review (the uat) | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 1,504,470 (USD 6.39) | 4 min 42 s, from 04:32:54 |
| The documents, the evidence, the task record; the freeze, the ladder and the pull request | execution | Claude Opus 5.5 | max | not reported | 04:33 to 04:39 |
| Review round 1, first try: each harness skipped | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI; Claude Fable 5.1 on the Claude Code CLI | `xhigh` | not reported | 04:39:46 to 05:01:18 |
| Review round 1, second try, first and second harness: skipped | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `xhigh` | not reported | 15:57:26 to 16:03:50 |
| Review round 1 | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 688,658 (USD 3.29) | 7 min 5 s, 16:03:58 to 16:11:05 |
| The close-out, with notes 1 and 5 of round 1 | execution | Claude Opus 5.5 | max | not reported | 16:11 to 16:16 |

From 05:03 to 15:56 the task did no work: the background run of the first try
of round 1 was stopped at its time limit, and the work went on at 15:56.
