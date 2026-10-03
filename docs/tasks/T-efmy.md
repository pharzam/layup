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
