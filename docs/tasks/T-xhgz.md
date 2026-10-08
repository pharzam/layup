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
