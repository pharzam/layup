# T-529i — test runs

The records of Intake (#212, row 40): the eight Go schemas against their blocks,
and the row rules of the readers.

## Red (23:42Z, 2026-10-10)

The tests before `intake.go`, `go vet ./internal/records/` and with
`-tags=integration`, each:

```
vet: internal/records/intake_test.go:34:63: undefined: ReadMeaning
```

## Mutations (2026-10-10, 23:53Z to 23:54Z, the second run; the first ended at 23:51Z)

Each rule of `intake.go` mutated alone, the empty-value rule for two records (its condition joined with `&& false`,
so each mutant builds), then `go test ./internal/records/`: 30 of 30 caught.
The 30 are each rule of a row or a table but the empty-value rule, which was
mutated for `meaning` and `answers` only; the other six records share the same
helper, `missing`, and each column has its case.
A first run counted mutants that did not build as caught (a condition replaced
by `false` left a variable unused); they were rewritten so each builds. That
run also found two rules no test could catch on its own: the session's form
(the test value `S-zz` was refused by the column's type before the rule; the
cases now use `S-1a2b3c4z`, a value of the type that is not hexadecimal), and
an "end with no start" case that the rule "end above start" covered (the two
were made one rule). Each mutant's first failing line:

```
== meaning-order: caught intake_test.go:144: meaning: an ID out of order: read, want an error in column id
== meaning-start-end: caught intake_test.go:144: meaning: start not below end: line 2, column "quote": the quote is the file's bytes from start to end: 0 bytes, want an error in column end
== meaning-quote-length: caught intake_test.go:144: meaning: offsets that differ from the quote's length: read, want an error in column quote
== meaning-session: caught intake_test.go:144: meaning: a session that is not S- and 8 hexadecimal characters: read, want an error in column session
== meaning-required: caught intake_test.go:144: meaning: the empty value in question: read, want an error in column question
== span-start-end: caught intake_test.go:144: spans: start not below end: read, want an error in column end
== span-fact-1: caught intake_test.go:144: spans: fact 0: read, want an error in column fact
== span-class: caught intake_test.go:144: spans: a fact with no class: read, want an error in column class
== span-reason: caught intake_test.go:144: spans: a fact with a reason: read, want an error in column reason
== draft-order: caught intake_test.go:144: draft: an ID out of order: read, want an error in column id
== draft-covers: caught intake_test.go:144: draft: a covered fact 0: read, want an error in column covers
== draft-session: caught intake_test.go:144: draft: a session that is not S- and 8 hexadecimal characters: read, want an error in column session
== coverage-order: caught intake_test.go:144: coverage: an ID out of order: read, want an error in column id
== coverage-fact: caught intake_test.go:144: coverage: an unsupported item with a fact: read, want an error in column fact
== coverage-fact-1: caught intake_test.go:144: coverage: fact 0: read, want an error in column fact
== coverage-requirement: caught intake_test.go:144: coverage: an uncovered item with a requirement: read, want an error in column requirement
== coverage-session: caught intake_test.go:144: coverage: a session that is not S- and 8 hexadecimal characters: read, want an error in column session
== q-ref: caught intake_test.go:144: questions: a meaning question whose ref is not M-NNN: read, want an error in column ref
== q-kind: caught intake_test.go:144: questions: a meaning question of a kind of no meaning: read, want an error in column kind
== q-asker: caught intake_test.go:144: questions: a rule question not asked by layup psb check: read, want an error in column asker
== q-start-pair: caught intake_test.go:144: questions: an end with no start: read, want an error in column start
== q-noquote: caught intake_test.go:144: questions: an intent question with a quote: read, want an error in column start
== q-below: caught intake_test.go:144: questions: a start with no end: read, want an error in column end
== q-asked: caught intake_test.go:144: questions: an asked question with an answer: read, want an error in column answered
== q-answered: caught intake_test.go:144: questions: an answered question with no answer: read, want an error in column answered
== q-order: caught intake_test.go:144: questions: an ID out of order: read, want an error in column id
== q-followup: caught intake_test.go:144: questions: a follow-up of no question before it: read, want an error in column ref
== review-session: caught intake_test.go:144: answers-review: a session that is not S- and 8 hexadecimal characters: read, want an error in column session
== budget-name: caught intake_test.go:144: budget: a name that is not of the block: read, want an error in column name
== answers-required: caught intake_test.go:144: answers: the empty value in author: read, want an error in column author
30 of 30
```

## Green (2026-10-10, 23:59Z)

Each with exit 0: `go build ./...`, `go vet ./...`, `gofmt -l internal` (empty),
`go test ./...`, `go test -tags=integration ./...` (with
`TestTheSchemasEqualTheirBlocks` and the eight names in `built`),
`go test -tags=e2e ./...`; `adr-lint`, `prd-lint`, `link-lint`, `setup-check`,
`run-discipline-tests`, `git diff --check e689d12`. The first run of the tests
failed to link three test binaries, "no space left on device", as `/tmp` (a
tmpfs of 1.9 GB) was full of the author's scratch clones; with those removed and
`GOTMPDIR` on the disk, the run passed.

## The fix of round 1 (2026-10-11, 00:10Z to 00:14Z; O-197 at 03:50Z)

Notes 2 to 6 of round 1, and its finding 1 (the two rows of `traceability.md`
of other packages restored to the base). The reverse direction of the quote of
`questions` (note 3) was written with its four cases; its two mutants, the rule
of a needed quote and `G1` among the rows with none, are each caught:

```
intake_test.go:148: questions: a gap of meaning with no quote: read, want an error in column start
intake_test.go:72: questions: line 7, column "start": a gap of meaning, an uncovered or class item and a rule gap G2 to G5 have a quote
```

The second mutant's first try did not build (an operator precedence of `&&`
over `||`); it was rewritten. The condition of the quote that cannot fail
(`utf8.ValidString`, as `tsv.Read` refuses bytes that are not UTF-8) is gone.
