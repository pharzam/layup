# T-dq05 — layup psb check: deterministic gap check of a problem statement

Issue: [#37](https://github.com/pharzam/layup/issues/37) (parent [#29](https://github.com/pharzam/layup/issues/29)). PR: [#40](https://github.com/pharzam/layup/pull/40).

## Verdict

Delivered: `layup psb check FILE` writes the gap questions of a problem statement
as one TSV batch, from five deterministic rules (G1–G5) in `internal/psb`, with
one golden test over four inputs. On LAYUP's own PSB it gives 19 questions, among
them the technology stack (the question the manual setup had to ask). Round 1:
material (G3 with no terms table; a table-state leak; a short row), fixed and
pinned by `edge.md`. Round 2: nothing material. Lesson written back:
`guardrails.md` §2, "The hook refuses a red commit".

## Resource record

Recorded, not budgeted (ADR-0007). Author tokens: `not reported`.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan and its review | reasoning | Claude Opus 5.5; Claude Fable 5.1 | not reported | 259,372 (115,898 + 143,474) | 18 min (6 + 12) |
| The decay review rounds | reasoning | Claude Fable 5.1 | not reported | 236,160 (114,887 + 121,273) | 16 min (9 + 7) |
| Writing the tests and the code | execution | Claude Opus 5.5 (one tier: limit recorded) | not reported | not reported | 40 min |
| Isolate, guardrails, docs, close-out | `—` | Claude Opus 5.5 | not reported | not reported | 10 min |
| **Total** | | | | 495,532 reported + not reported | about 84 min |
