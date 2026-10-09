# T-z5dj — row 31 of the plan, the calls of M2b of internal/git

Issue: [#157](https://github.com/pharzam/layup/issues/157), row 31 of [the tasks
of M2b](../plan/README.md#the-tasks-of-m2b), opened by `T-fdaq` (#152). Serves
`F-0003#43` through `REQ-003` and `NFR-001`. Base `8d2fa2e` (the merge of #171).
Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-z5dj/`](../../runs/T-z5dj/).

## Plan and plan review

The plan (R12, comment 6080323845), its review (6080575151) and the author's
answer (6080575807) are comments on #157. The plan review (Claude Fable 5.1,
effort `xhigh`, on the Claude Code CLI with stream output, a fresh read-only
session in a clone at `8d2fa2e`, 12 min 27 s) gave `approve-with-conditions`:
Budget maximum 750 lines added plus removed over 12 files against `8d2fa2e`,
close-out inside; Cycle cap 1; no panel. Its four conditions are applied: one
error shape (`ErrNotACommit` in a `*FailedError`); `core.alternateRefsCommand`
in the list and a control that shows each live key; both Test cells; the
scratch directories removed on each return. The goal count is 3, final by O-187
of #152.

## What was done

1. **D1:** `CloneLocal`, `IsAncestor`, `DiffFile` and `DiffBinary` in
   `internal/git/git.go`, each with its input checked before `git` starts.
2. **D2:** `FetchSession` and `ErrNotACommit`: a scratch bare repository whose
   alternate is the session's object directory, `cat-file -t`, `update-ref`, and
   the fetch with an empty hooks directory; no `git` runs in the session's clone.
   `TestAHostileSessionRunsNothing` writes 22 program-starting keys of
   `git-config(1)` (git 2.47.3) into the session's configuration, each with a
   marker named for it; a control shows, by plain reads in a copy, that twelve of
   them fire on this host, and `FetchSession` fires none.
3. **D3:** `packages.md` ("The calls of `M2b`": `ErrNotACommit`, the scratch
   directories, the directories of the commands, the refusal of a branch, the
   known limit of `--update-head-ok`); `traceability.md`; the Test cells of
   `REQ-003` and `NFR-001`, with a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-z5dj/test-runs.md): red before the
calls (no compile); the control red three times until each key had its own
driver; red with a `FetchSession` that ran `git` in the session (the hostile
test found `core.fsmonitor` and `filter.y.clean`); green after.

**Known limits.** The list of keys is that of `git-config(1)` of 2.47.3; a key
that a later `git` adds is not in it. Ten keys (`core.sshCommand`,
`core.askPass`, `core.editor`, `core.pager`, `credential.helper`,
`uploadpack.packObjectsHook`, `sequence.editor`, `gpg.program`,
`remote.origin.receivepack`, `merge.x.driver`) are asserted, not shown live. A
`DST` that is the branch of `HEAD` of `dir` is refused by `git`.

**The rejected alternatives:** a fetch from the session's clone (its
`upload-pack` reads the session's configuration); `git bundle` (it runs `git` in
the session's clone); copying the objects by hand.
