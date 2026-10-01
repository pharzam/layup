## Review record — round 2

| Field | Value |
| ----- | ----- |
| Commit reviewed | `32304b953100a28d6838ee3ca01d6bf22b2f2e80` |
| Reviewer | Claude Fable 5.1 on Claude Code |
| Lens | correctness and acceptance criteria; the fixes of round 1 |
| Briefed on | `.review-in/brief.md`; `.review-in/issue-74-body.md`; `.review-in/round-1.md`; `git diff d0634a4 HEAD` in full (9 files); `docs/spec/` at HEAD in full (`README.md`, `packages.md`, `records.md`, `psb-check.md`, `setup.md`, `gate.md`); `docs/architecture.md` §3, §5, §6, the procedure of §11, and L-A6 and L-B1 of §15; `docs/setup/setup-check.sh` (`check_facts`, `check_markers`, `MK_EXEMPT`); ADR-0014; `docs/glossary.md` (the format, and the four new rows); `docs/tasks/T-0drh.md`; `docs/engineering-discipline.md` (the two tests of a material finding); `docs/setup/open-gaps.tsv`, `facts.sha256`, `armature.pin`, `steps.tsv`, `branch-protection.json`; `internal/psb/check.go` (the `Q-NNN` form) and `internal/cli/cli.go` (`version`); the tree of the root commit `d2516fd` (no `docs/setup/`; the rule files); `F-0001` facts 10 to 12. `runs/T-0drh/review-round-1.md` is byte-identical to `.review-in/round-1.md`. The local checks (`adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests`, `setup-check`, `git diff --check 648b37f HEAD`), `go build`, `go vet` and `go test -count=1 ./...` were run at HEAD: each exits 0. The diff against `648b37f` is 17 files, 1,343 lines added plus removed. |
| Barred from | the issue comments of #74 |
| Independence claimed | A fresh session, read-only. The author is Claude Opus 5.5; this reviewer is Claude Fable 5.1. No issue comment was read, and `gh` was not used. No file of the repository was changed; this record is the only file written, in the untracked `.review-in/`. |
| Cycle | 1 |
| Verdict | `not mergeable, findings recorded` |

### Round-1 findings

1. closed — `setup.md:98` writes `git push origin layup-setup:main` before the apply; `setup.md:207-212` gives the order, with the reason.
2. closed — `setup.md:95` names LAYUP's `MK_EXEMPT`, embedded at the engine's version, with the rule for a path a target lacks; `setup.md:239` and `:98` no longer say "the baseline's".
3. closed — `records.md:80-81` has the two rows; `setup.md:179-183` is the schema of `open-gaps.tsv` (no header, key `file` + `marker`); `setup.md:185-188` gives the form of `facts.sha256`; `README.md:86-88` defines `no-header` (note 7 below).
4. closed — `gate.md:71-72` puts `tool not found` before `no product path`, and `gate.md:64` says that the first matching line decides.
5. closed — `setup.md:86` has S01 run `layup psb check`, each gap a question `Q-NNN`; `setup.md:112` lists `Q-NNN`; `setup.md:91` says what a fact holds; `psb-check.md:63-67` makes `id` the question ID in phase 1. The same sentence leaves a second gap for the marker answers: finding 1 below.
6. not closed, no fix — the PRD is not in the diff; a note that the Operator settles (whether O-115's sentence bars the §7.1 change).
7. closed — `setup.md:245` names `--base` and `--head` of both runs and the row for each outcome (note 3 below on the precedence).
8. closed — `setup.md:128` adds `not-active` with exit 1.
9. closed — `packages.md:54-56` says that `internal/cli` runs each step's check from `internal/verify`.
10. closed — `packages.md:89` makes `internal/route` phase 2.
11. closed — `records.md:209` gives the `diagnosis` row a payload with each open unknown and its evidence.
12. closed — `records.md:220` and `setup.md:62` give their reasons (note 5 below on the citation).
13. closed — `records.md:161` names `first_output` as the status column.
14. closed — `gate.md:155-157` says "a connection to a model service", and names the two network uses.
15. closed — `setup.md:89` marks the S04 change **decided here** with §5 Scaffold 6 as the basis.
16. closed — `setup.md:297-300` names the four `layup/` checks, the phase and who applies them, as §5 Scaffold 6 says.
17. closed — `architecture.md:192-194` now allows the Operator's push of the records' first commit by its SHA (O-115); ADR-0014's Consequences line (`0014:95-98`) says the same.
18. closed — `setup.md:233` makes `method` equal the `NFR-006` text with the source and the commit.
19. closed — `setup.md:112` adds `O-<name>`; `setup.md:199` uses `O-verify`.
20. closed — `setup.md:151` marks `computed` **decided here** with its reason.
21. closed — `setup.md:285` makes `fixture` a `text`, relative to the entry's directory.
22. closed — `psb-check.md:39` fixes the case of the word; `psb-check.md:21-24` states the field rule and names the present deviation for #29.
23. closed — `README.md:42-43` exempts `layup version` and gives its line and exit code.
24. closed — `setup.md:166-172` lists the entries; `setup.md:220` adds "No flag." (note 8 below on L-A6).
25. closed — `records.md:82` says that the engine does not read `steps.tsv`.
26. closed — `glossary.md:96-99` adds `TSV`, `UTF-8`, `RFC 3339` and `ISO 4217` (note 6 below).

### Raw findings

Finding 1 is material. Findings 2 to 9 are notes.

**1. `docs/spec/setup.md:91` with `setup.md:57-58`, `:95` and `:235` — the raw fact record cannot hold the marker answers in the documented order, so check `facts` fails on a correct setup.** Cited sentences: `setup.md:91` (S06) "writes the answers as a raw fact record, one fact per question ID of `answers.tsv`, each with the question ID, the question text and the answer"; `setup.md:57-58` "it reads `WORK/out/record.tsv` and skips each step that has a `done` row there"; `setup.md:95` (S10) "A marker with no answer row stops the run; the table lists all of them at once"; `setup.md:235` (check `facts`) "the answers record has one fact per question ID of `answers.tsv`"; `setup.md:247` "a correct setup exits 0". Basis: the ID of a marker's question is a hash that the S10 stop table gives (`setup.md:117-119`), so in the documented flow (`setup.md:196-197`: a stop at each step that needs an input) the Operator adds the `M-<x8>` rows after S10 stopped, when S06 already has its `done` row. The next run skips S06, so the raw fact record holds the `S01-` and `Q-` answers only, while `answers.tsv` now has `M-` rows. Check `facts` then fails ("one fact per question ID of `answers.tsv`"), and `layup setup verify` exits 1 on a correct setup; S15 refuses that `verify.tsv` (`setup.md:201-202`). The round-1 fix made the S06 sentence name "of `answers.tsv`" and left it false for the marker answers; `NFR-001` item 3 (`records.md:105-106`) needs them in the raw fact too. One sentence fixes it: S11 appends the marker answers to the raw fact record (and its index row), or check `facts` counts the IDs that exist when S06 runs and S11's rows apart. `material` (the expected exit code differs from the documented result; a claim left false).

**2. `docs/spec/setup.md:86` against `docs/spec/packages.md:43` — the package that runs `layup psb check` at S01 is not named, and the step runner may not import it.** Cited sentences: `setup.md:86` "runs `layup psb check` on `inputs/briefs/problem-statement.md`"; `packages.md:43` `internal/setup` may import "`internal/tsv`, `internal/git`, `internal/catalog`", starts no program; `packages.md:40` the rules G1 to G5 are `internal/psb`; `packages.md:31` "An import that the table does not allow is a defect". Basis: the table lets `internal/cli` import both (`packages.md:37`), as for the per-step checks (`packages.md:54-56`), but no sentence says that `internal/cli` runs the gap table for S01 and hands it to the step runner. Say which, or add `internal/psb` to the row of `internal/setup`. `note`, classified as round-1 finding 9 was: the table allows the route; the text does not name it.

**3. `docs/spec/setup.md:245` — the three outcome rules of `gate:<kind>` have no stated precedence.** Cited sentence: "either run `not-active` gives `not-active`; a fixture run that is not `fail` gives `fail`, reason `fixture not detected`". Basis: a fixture run that gives `not-active` matches both rules; the exit code is 1 either way, so only the row's result differs. Say "the first rule that matches decides", as `gate.md:64` now does. `note`.

**4. `docs/spec/setup.md:220-221` — "changes neither" with three inputs, and the fixture commit.** Cited sentence: "It reads `WORK/target` at the head of `layup-setup`, `WORK/out/record.tsv` and `WORK/inputs/answers.tsv`, in a scratch work tree, and changes neither." Basis: three things are read; and the fixture run of `:245` needs "a commit of the kind's known-bad fixture applied" that `git rev-parse` resolves in `REPO` (`gate.md:50-51`), so the check writes one commit object, on no ref, into `WORK/target`. Say "changes no ref of `WORK/target` and no file of `WORK/out` or `WORK/inputs`". `note`.

**5. `docs/spec/setup.md:62` — the citation of the work area.** Cited sentence: "decided here: §5 names the host's work area and no layout". Basis: the work area on the host is named in §1 (`architecture.md:26`); §5 names none. `note`.

**6. `docs/glossary.md:96` against `docs/spec/README.md:86-88` — the new `TSV` row says that every record has a header row.** Cited sentence: "LAYUP writes every record in this form, with a header row". Basis: `open-gaps.tsv` is a record with no header row (`setup.md:179`, `no-header`). The row is a summary; add "except where a schema says `no-header`". `note`.

**7. `docs/spec/setup.md:185-188` — `facts.sha256` has a prose form and no named key.** Basis: AC 4 asks for the column order, types, key and writer "in a form that a Go test can read"; the form is `sha256sum`'s, so it is exact, and the writer is in `records.md:80`; the key (one line per raw facts file: the path) is implied, not said. AC 4 is met in substance; one clause names the key. `note`.

**8. `docs/spec/setup.md:171-172` — L-A6 does not cover a rule document that a newer baseline adds.** Cited sentence: "A baseline whose rule files differ is known limit L-A6 of the architecture." Basis: L-A6 (`architecture.md:1522-1525`) says that a changed baseline "can break a step" and that `layup setup verify` "then fails on that step"; a new rule document that the fixed list of `:166-171` does not name breaks no step and fails no check of phase 1 (the `.sh` rule is the only dynamic entry), so `layup/rules` of phase 2 would not protect it. Say that this case is not detected in phase 1, or add it to the limit. `note`.

**9. `docs/spec/psb-check.md:66-67` — "the problem statement does not change during a setup, so the IDs stay" has no check.** Basis: the brief is an input file that the Operator can edit between two runs; S06 copies it byte for byte, but nothing compares the copy with the file that S01 ran the rules on. A `computed` record row with the brief's SHA-256 at S01, checked at S06, makes the claim a check. `note`.

### Acceptance criteria

| AC | Result | Where |
| -- | ------ | ----- |
| 4 | met | `records.md:80-81`, `setup.md:179-188`, `README.md:86-88` (note 7). |
| 8 | not met while finding 1 stands | the section of `REQ-002` contradicts its own order of steps and its own check `facts`; each other round-1 disagreement with §5 is closed. |
| 10 | met | run again at `32304b9`: every check, `go build`, `go vet` and `go test` exit 0; 17 files, 1,343 lines, within the budget of 1,900 over 20. |
| 11 | met | `glossary.md:96-99`, `architecture.md:192-194`, ADR-0014 Consequences; the round-1 record is in `runs/T-0drh/`. |

The other criteria stand as round 1 found them (1, 2, 3, 5, 6, 7, 9: met).
