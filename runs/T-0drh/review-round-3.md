## Review record — round 3

| Field | Value |
| ----- | ----- |
| Commit reviewed | `f83320e6d82cd49cce30dc5ceb93968ae1a65425` |
| Reviewer | Claude Fable 5.1 on Claude Code |
| Lens | correctness and acceptance criteria; the fixes of round 2 |
| Briefed on | `.review-in/brief.md`; `.review-in/issue-74-body.md`; `.review-in/round-2.md`; `git diff 32304b9 HEAD -- docs/` in full (3 files: `docs/glossary.md`, `docs/spec/setup.md`, `docs/tasks/T-0drh.md`; the commit's fourth file, `runs/T-0drh/review-round-2.md`, is out of scope and is byte-identical to `.review-in/round-2.md`); `docs/spec/setup.md` at HEAD in full; `docs/spec/README.md`, `packages.md`, `psb-check.md` and `records.md` in full; `gate.md:40-80`; `docs/architecture.md` §1, §5 and L-A6 of §15; `docs/setup/setup-check.sh` (`check_facts`); `docs/facts/README.md` (the index, "Correcting a facts document"); `docs/engineering-discipline.md` Bootstrap mode rule 3; `docs/glossary.md` rows `Material` and `TSV`; `docs/tasks/T-0drh.md`. The local checks (`adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests`, `setup-check`, `git diff --check 648b37f HEAD`), `go build`, `go vet` and `go test -count=1 ./...` were run at HEAD: each exits 0. The diff against `648b37f` is 18 files, 1,425 lines added plus removed; the fix commit alone is 4 files, 102 lines. |
| Barred from | the issue comments of #74 |
| Independence claimed | A fresh session, read-only. The author is Claude Opus 5.5; this reviewer is Claude Fable 5.1. No issue comment was read, and `gh` was not used. No file of the repository was changed; this record is the only file written, in the untracked `.review-in/`. |
| Cycle | 2 |
| Verdict | `not mergeable, findings recorded` |

### Round-2 findings

1. closed — `setup.md:96` makes S11 add each marker answer to the raw fact record and update its line in `facts.sha256`; `setup.md:238` splits the IDs by step. The new sentences leave one input state open: finding 1 below.
2. closed — `setup.md:86` names `internal/cli` as the runner of the rules of `internal/psb`, which hands the gap table to `internal/setup`; `packages.md:37` allows both imports.
3. closed — `setup.md:248` says "The first rule that matches decides", as `gate.md:64` does.
4. closed — `setup.md:224` names the three inputs, what the command does not change, and the fixture commit as an object on no ref.
5. closed — `setup.md:62` cites §1 (`architecture.md:26`) for the work area.
6. closed — `glossary.md:96` adds "except where its schema says `no-header`".
7. closed — `setup.md:188` names the key of `facts.sha256` (the path).
8. closed — `setup.md:172-175` records the limit of phase 1: a new rule document of a newer baseline is not in the register.
9. closed — `setup.md:86` adds the record row `brief.sha256` at S01; `setup.md:91` makes S06 refuse a changed problem statement (exit 2).

### Raw findings

Finding 1 is material. Findings 2 to 5 are notes.

**1. `docs/spec/setup.md:96` (S11) and `:238` (check `facts`), with `:91` (S06) — when `answers.tsv` holds the `M-` rows before S06 runs, the sentences give two readings, and one makes `layup setup verify` exit 1 on a correct setup.** Cited sentences: `setup.md:96` "Adds each marker answer to the raw fact record of S06, one fact per `M-` question ID, and updates its row in `facts.sha256`, so the record holds every answer of `answers.tsv`"; `setup.md:238` "the answers record has one fact per question ID of `answers.tsv` (the `S01-` and `Q-` answers from S06, the `M-` answers from S11)"; `setup.md:91` (not changed by the fix) "one fact per question ID of `answers.tsv`"; `setup.md:250` "a correct setup exits 0". Basis: an `M-` ID is the hash of the file and the marker text (`setup.md:117-119`), so it is the same for each target set up from the same baseline commit; the answers file is the Operator's (`setup.md:66`), its `question` column asks only for the ID "as the stop table gives it" (`setup.md:135`), and S10 does not stop when the rows exist (`setup.md:95`). So an Operator who copies the answers of an earlier work area, or who starts a setup again, has the `M-` rows at S06. In the order of `setup.md:199` the state does not arise, because the `M-` IDs come from the S10 stop; the specification does not limit the Operator to that order. In that state S06 writes one fact per `M-` ID, as its sentence says, and S11 "adds each marker answer" again: two facts per `M-` ID, so check `facts` ("one fact per question ID") fails, `layup setup verify` exits 1, and S15 refuses `verify.tsv` (`setup.md:204-205`). The parenthesis of `:238` says the opposite: the `M-` answers come from S11 only, so S06 writes no `M-` fact. Two operative sentences disagree, and the exit code of a correct setup depends on which one the code of #29 follows. Round 2 proposed "S11 appends the marker answers" and left the same gap. One clause fixes it: at S06, "one fact per `S01-` and `Q-` question ID of `answers.tsv`"; or at S11, "each marker answer that the record does not hold yet". `material` (an operative contradiction; the expected exit code differs from the documented result).

**2. `docs/spec/setup.md:86` — the Output cell of S01 does not name the new record row.** Cited sentence: "record rows `stack`, `name`, `visibility`, `baseline`", in the row whose What-it-does cell says "A record row `brief.sha256` (`computed`) holds the SHA-256 of the problem statement that the rules read." Basis: the Output cell of S02 lists each record row of its step; a reader of that column misses `brief.sha256`. Add it to the cell. `note`.

**3. `docs/spec/setup.md:151-152` against `:86` — the rule of the source `computed` does not fit `brief.sha256`.** Cited sentences: `setup.md:151` "`computed`: a `git` command's output"; `setup.md:152` "`computed`: the command". Basis: the SHA-256 of a file is not the output of a `git` command (`git hash-object` gives a SHA-1), and the `ref` of the new row is not named. Check `sources` (`setup.md:246`) reads no `computed` ref, so no exit code changes. Name the ref, for example `sha256 inputs/briefs/problem-statement.md`, and widen the rule to "a command's output, or a hash that the engine computes". `note`.

**4. `docs/spec/setup.md:91` against `:53-56` — the exit code of the refusal, and the way out.** Cited sentences: `setup.md:91` "Refuses a problem statement whose SHA-256 differs from the row `brief.sha256` of S01 (exit 2), so the `Q-NNN` IDs stay true"; `setup.md:54` "2 a usage or input error (an input that does not match its schema)". Basis: a changed brief matches its schema, and the parenthesis is the only definition of an input error in this section (`README.md:55` adds "a missing or unreadable file", not a changed file). After exit 2, a run skips S01 (`setup.md:57-58`: a `done` row), so the Operator cannot run the gap check again on a revised brief in the same work area, and the sentence names no way out. Add "or an input that changed after a step read it" at `:54`, and one sentence at S06: "The Operator restores the file, or starts again in a new work area." `note`.

**5. `docs/spec/setup.md:96` against `docs/facts/README.md:69-72` — S11 changes a raw facts record after S06 committed it.** Cited sentences: `docs/facts/README.md:69` "A raw facts document is immutable, because it is evidence. You do not edit the customer's recorded words after the fact."; `:71` "add a **new** facts document"; `records.md:79` names this convention as the schema of the record. Basis: S06's commit (`setup.md:102-103`) holds the record, and S11 adds facts to it in a later commit. The addition changes no recorded words, both commits are on `layup-setup` before any push, and no check reads the convention, so no exit code changes. The index row that S06 writes is stale after S11 if it states a count, as LAYUP's `F-0004` row does (`docs/facts/README.md:83`), and the Evidence cell of S11 (`setup.md:96`) names checks `markers` and `sources` only, while check `facts` now reads S11's work too. Say that the addition is not a correction and that the index row states no count, or write the marker answers as a second raw fact record with its own index row. `note`.

### Acceptance criteria

| AC | Result | Where |
| -- | ------ | ----- |
| 4 | met | not changed by the fix; `setup.md:188` adds the key of `facts.sha256`. |
| 8 | not met while finding 1 stands | narrower than in round 2: the order of steps of `setup.md:199` now agrees with check `facts`; the input state of finding 1 does not. |
| 10 | met | run again at `f83320e`: each local check, `go build`, `go vet` and `go test` exit 0; 18 files, 1,425 lines, within the budget of 1,900 over 20. |
| 11 | met | `glossary.md:96`; `T-0drh.md:25-28` records O-116; the round-2 record is in `runs/T-0drh/`, byte-identical to `.review-in/round-2.md`. |

The other criteria stand as rounds 1 and 2 found them (1, 2, 3, 5, 6, 7, 9: met).
