# The runs of `cites.sh`

Task `T-w73g` (#133), 2026-10-07, on macOS with `/bin/sh`. The plan review asked
for three runs (condition 3): red for the right reason, a mutation, and green.

## 1. Red: the survey does not exist

Before `survey.md` was written:

```text
$ sh runs/T-w73g/cites.sh
cites: no survey: runs/T-w73g/survey.md
exit 1
```

## 2. Mutation: one path that does not exist at its commit

A copy of the survey, outside the checkout, with one more line that cites
`ruflo:v3/@claude-flow/guidance/src/no-such-gate.ts`; run on the five local
clones:

```text
$ sh runs/T-w73g/cites.sh <copy of survey.md> <directory of the clones>
MISSING ruflo:v3/@claude-flow/guidance/src/no-such-gate.ts
cites: 1 problem(s) in 67 source path(s)
exit 1
```

The other 66 lines of the output were `ok`.

## 3. Green

On the five local clones:

```text
$ sh runs/T-w73g/cites.sh runs/T-w73g/survey.md <directory of the clones>
cites: 66 source path(s), each present at its commit
exit 0
```

With no directory of clones, so that each repository is fetched at its commit
(6.8 s):

```text
$ sh runs/T-w73g/cites.sh
cites: fetching ucsandman/Agnostic-AI at aab09ff0d9533967045e7f2b04e478ec4e05528e
cites: fetching ruvnet/ruflo at 9fa701a466159242dfec06379ce1278c849a3e34
cites: fetching greenpolo/cc-multi-cli-plugin at 3dcb057c10b16d3ed911e8ce509e91398712c906
cites: fetching stuinfla/ruvnet-brain at b4590d469f52937550d8d89afa86a394a1d72875
cites: fetching kolezka/self-improvement-loop at ba245adbb5be8075fa01044dd82aac636155870f
cites: 66 source path(s), each present at its commit
exit 0
```

The 66 `ok` lines are left out above.
