# T-xhgz — row 23 of the plan, the calls Fetch and Push of internal/git

Issue: [#128](https://github.com/pharzam/layup/issues/128), row 23 of [the tasks
of M2a](../plan/README.md#the-tasks-of-m2a), opened by `T-zwke` (#124). Serves
`F-0003#42` through `NFR-001`. Base `b7c9e90` (the merge of #140). Author:
Claude Opus 5.5 on Claude Code. Started under O-169, the Operator's standing
rule of 2026-10-08 in the session ("after merge start next automatically").
Evidence: [`runs/T-xhgz/`](../../runs/T-xhgz/).

## Plan and plan review

The plan (R12, comment 6056476473). Its review (Claude Fable 5.1, effort
`xhigh`, on the Claude Code CLI with stream output, a fresh read-only session
in a clone at `b7c9e90`, 8 min 40 s; comment 6056634112):
`approve-with-conditions`; two goal classes (the two calls; the token in the
environment) with no split asked, as one table test covers both; condition 1
(the input rules and the one reading of a refused push, in `packages.md`) and
condition 2 (`Init`, then `Fetch` of `main`: `git` refuses it). The author's
answer (6056634444) measured both with `git` 2.54.0 and selected
`--update-head-ok` over a change of `Init` or a second ref. Budget maximum 500
lines added plus removed over 14 files against `b7c9e90`, close-out inside;
Cycle cap 1; no panel.

## What was done

1. **`internal/git`:** `Auth` (the `web` of the forge register and the token);
   `Fetch` (`git fetch --no-tags --update-head-ok -- URL REF:REF`) and `Push`
   (`git push --porcelain -- URL COMMIT:refs/heads/BRANCH`, never `--force`);
   each refuses its input before `git` starts unless the ref starts with
   `refs/` and holds no `:`, the commit is a full object ID and the branch is
   not empty and holds no `:`. With a token, the call's environment gets
   `GIT_CONFIG_COUNT=1`, `GIT_CONFIG_KEY_0=http.<web>/.extraHeader` and
   `GIT_CONFIG_VALUE_0=Authorization: Basic <base64 of x-access-token:TOKEN>`,
   after the fixed list; never in the arguments. A refused push is a
   `FailedError` of `Code` 1.
2. **The documents:** the row of `Fetch` and a bullet of the input rules, of
   `--update-head-ok` and of a refused push in `packages.md`; the fixed list
   names the three values of a token; `forge.md` names the value of the header;
   a pitfall in `guardrails.md` §2; the traceability rows; the Test cell of
   `NFR-001` and a §13 line. [`docs.sh`](../../runs/T-xhgz/docs.sh) checks
   them.

**The header form** (note 1 of the plan review): GitHub reads an installation
token over HTTPS as the user `x-access-token` and the token as the password,
which is HTTP Basic (GitHub Docs, "Authenticating as a GitHub App
installation", read on 2026-10-08).

**Tests:** [`test-runs.md`](../../runs/T-xhgz/test-runs.md): the two
measurements, red 1 and 2, green, and two mutations, each caught; `docs.sh` red,
then green.

**The rejected alternatives:** the token in the URL (arguments, errors, `ps`);
a credential helper (D3 of #79); a third error kind for a refused push; `Init`
on another branch or bare, or a second ref, for the read-back.

**Off this task's path** (note 3 of the plan review): a private target needs
the token for the `Clone` of the restart; written on #130 (row 25).

## Review rounds

Round 1 (Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI with stream
output, a fresh read-only session in a clone at `e7ede10`, 7 min 28 s; comment
6056857046): `material`. Finding 1: `Auth` refused a token whose `web` is empty
or ends in `/`, a rule that no sentence gave. Fixed in `4c28d28` (comment
6056883592): the bullet of `packages.md` gives it, and a rule of `docs.sh`
failed on the old text, then passed; notes 2, 3 and 4 applied (the error keeps
no line `!`; `Code` 1 of a missing commit; two lines wrapped).

Round 2 (the same set-up, at `4c28d28`, 9 min 28 s; comment 6057044577):
`nothing material in scope`, three notes, kept as known limits so the reviewed
text is the text that lands: note 1, the sentence on `Code` 1 of a missing
commit rests on a reading of `git`'s source, not a measurement; note 2,
"refused by the remote" means a `[rejected]` or `[remote rejected]` line, while
an HTTP `403` of the forge exits 128 (row 25 fails the step on both); note 3,
`Clone` also writes the remote `origin`, which `Init` and `Fetch` do not; no
step reads it.

## Verdict

Delivered: `Push` refuses a push to a local bare repository that is not a
fast-forward (`Code` 1, the branch unchanged); `Fetch` reads the branch back
into a repository that `Init` made; with a token, both carry it as a header in
the call's environment only. The review ended by decay at cycle 1. The diff
against `b7c9e90` is inside 500 lines over 14 files.

Next: row 25 (`T-trej`, #130), whose After cell holds rows 21, 22b, 23 and 24,
all merged with this one.

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC on 2026-10-08; tokens are the
`result` event of the Claude Code CLI (stream runs); `not reported` otherwise.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | max | not reported | 08:58 to 09:03 |
| Its plan review | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 620,999 (USD 4.48) | 8 min 40 s, from 09:03 |
| The two measurements; the answer | reasoning | Claude Opus 5.5 | max | not reported | 09:12 to 09:13 |
| The work, test first; the documents | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 09:13 to 09:18 |
| Round 1 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 582,730 (USD 3.58) | 7 min 28 s, from 09:18 |
| The fix of round 1 | execution | Claude Opus 5.5 | max | not reported | 09:26 to 09:28 |
| Round 2 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 698,099 (USD 4.43) | 9 min 28 s, from 09:28 |
| The close-out | execution | Claude Opus 5.5 | max | not reported | 09:38 to 09:42 |
