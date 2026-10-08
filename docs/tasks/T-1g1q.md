# T-1g1q — row 22b of the plan, the two host registers and the key file

Issue: [#137](https://github.com/pharzam/layup/issues/137), row 22b of [the tasks
of M2a](../plan/README.md#the-tasks-of-m2a), opened by `T-esfe` (#127) after
O-170 (b) split row 22. Serves `F-0003#42` through `NFR-001` and `NFR-007`. Base
`10a5691` (the merge of #138). Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-1g1q/`](../../runs/T-1g1q/).

## Plan and plan review

The plan (R12, comment 6054078662) listed three goal classes. Its review (Claude
Fable 5.1, effort `xhigh`, on the Claude Code CLI with stream output, a fresh
read-only session in a clone at `10a5691`, 12 min 24 s; comment 6054282765) gave
`reject` on the goal count alone (three), with six more conditions; the author's
answer (6054283098) applied them and asked **O-171**. **O-171** (b) (comment
6054413706; the text was written by a Devin session and pasted and adopted by the
Operator, comment 6054453470): the two host-register readers and the key-file
validation are one goal; the forge interface (class 3) moves to row 24, its
delivery scope only, its code staying in `internal/forge`, with conditions 4 to 6
and note 5 of the review carried to #129, open. Budget maximum 1,050 lines added
plus removed over 17 files against `10a5691`, close-out inside; Cycle cap 1; no
panel.

## What was done

1. **`internal/forge`** (new): `ForgeRegisterSchema` and `ReadForgeRegister` (one
   row; each column but `watch_slug` holds a value; an absolute `key_file`; each
   error names its line), and `CheckKeyFile`, which reads the mode, the owner and
   the bytes and calls the pure `checkKey`: mode exactly 0600, the owner the user
   of the run, checked before the read; one PEM block and one `BEGIN` line (`RSA PRIVATE KEY`, PKCS #1, or
   `PRIVATE KEY`, PKCS #8 of the RSA algorithm; version 0), an RSA key that `Validate`
   accepts; it gives the key, so the adapter of row 24 parses nothing. It uses
   `encoding/pem`, `encoding/asn1` and `crypto/rsa`, not `crypto/x509` (rule 5).
2. **`internal/route`** (new, its `M2a` job only): `HarnessRegisterSchema` and
   `ReadHarnesses`, which gives the rows and the IDs in their order; `wall` is 1
   or more; `cap` may be empty; an empty register is allowed.
3. **Both owners** compare their schemas with their blocks; the two names moved
   to `built`; `TestPackageRules` passes with the two packages.
4. **The documents:** `forge.md` says "one PEM block" and its two forms; the Job
   cell of `internal/forge`; the plan rows 22b and 24 (class 3 to row 24 by
   O-171); the traceability rows; the Test cells of `NFR-001` and `NFR-007` and a
   §13 line. [`docs.sh`](../../runs/T-1g1q/docs.sh) checks them.

**Tests:** [`test-runs.md`](../../runs/T-1g1q/test-runs.md): red 1 (no package),
red 2 (a stub; each rule case failed for its reason), green; one comparison shown
to fail with a column changed on a backup copy; `docs.sh` red, then green. The
tests make one 2,048-bit key per test binary at run time: no key file enters the
tree.

**The rejected alternatives:** the key parse with `crypto/x509` (it depends on
`net`); the harness register in `internal/forge`; the interface in this row
(O-171).

## Review rounds

Round 1 (Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI with stream
output, a fresh read-only session in a clone at `5f5d4b8`, 12 min 18 s; comment
6054797508): `material`. Finding 1: a trailing space in `test-runs.md`
(`git diff --check` exit 2). Finding 2: the key-file row of `run.md` "Input
states" had two readings of "naming the file and its mode". Fixed in `3f8682a`
(comment 6054867255): each error names the mode, and `forge.md` gives the one
reading (the mode in each state but a missing file); notes 3, 4, 5, 8 and 9
applied; note 6 no change; note 7 no budget change needed.

Round 2 (the same reviewer set-up, at `3f8682a`, 8 min 16 s; comment
6055024105): `material`, the last round under cap 1. Finding 1: `CheckKeyFile`
read the bytes before the check, so a file that the run cannot read (mode `0000`,
or another user's) gave `permission denied` with no mode; the author confirmed
it. **O-172 (a)** (the Operator's answer in the session to comment 6055024560,
recorded in comment 6055247938): one more cycle, cap 2. Fixed in `4a20bb4`
(comment 6055767408): `checkFile` runs before the read; exactly one `BEGIN` line
(note 3); notes 4, 5 and 6 applied; note 2 no change (`run.md` would be an 18th
file; `forge.md` holds the reading and links the row).

Round 3 (the same set-up, at `4a20bb4`, 5 min 17 s; comment 6055860261):
`nothing material in scope`, two notes, not applied: note 1 (a mode that cannot
be read, when `os.Stat` fails on a parent directory: no state of `run.md`);
note 2 (the read-error branch after `checkFile` has no simple test; `checkFile`
runs twice on a good file, which keeps the pure function complete).

An order slip of the author: the commits of the first fix were pushed before
`gitleaks` ran; it ran right after the push (3 commits, no leaks), and before
each later push.

## Verdict

Delivered: a forge register or a harness register that breaks its schema is
refused with its line, by its owner package (`internal/forge`,
`internal/route`); the key file is checked before any forge call, each refusal
naming the file and its mode, with no package of rule 5. The review ended by
decay at cycle 2 (cap 2 by O-172 a). The diff against `10a5691` is inside 1,050
lines over 17 files. The forge interface stays open in row 24 (#129), with
conditions 4 to 6 and note 5 of comment 6054282765.

Next: row 23 (`T-xhgz`, #128), then row 24 (`T-6bq5`, #129).

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC on 2026-10-08; tokens are the
`result` event of the Claude Code CLI (stream runs); `not reported` otherwise.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | max | not reported | 06:18 to 06:41 |
| Its plan review | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 933,280 (USD 5.91) | 12 min 24 s, to 06:54 |
| The answer; O-171 | reasoning | Claude Opus 5.5 | max | not reported | 06:54 to 07:03 |
| The work, test first; the documents | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 07:03 to 07:08 |
| Round 1 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,225,513 (USD 6.06) | 12 min 18 s, from 07:08 |
| The fixes of round 1 | execution | Claude Opus 5.5 | max | not reported | 07:21 to 07:26 |
| Round 2 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 970,002 (USD 4.25) | 8 min 16 s, from 07:26 |
| O-172; the fixes of round 2 | execution | Claude Opus 5.5 | max | not reported | 07:35 to 08:19 |
| Round 3 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 600,761 (USD 3.32) | 5 min 17 s, from 08:19 |
| The close-out | execution | Claude Opus 5.5 | max | not reported | 08:26 to 08:35 |

From 07:35 to 07:47 the task waited for O-172.
